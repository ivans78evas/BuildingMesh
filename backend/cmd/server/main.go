package main

import (
	"context"
	_ "construction-ar-backend/docs"
	"construction-ar-backend/internal/cache"
	"construction-ar-backend/internal/config"
	"construction-ar-backend/internal/handlers"
	"construction-ar-backend/internal/logger"
	"construction-ar-backend/internal/middleware"
	"construction-ar-backend/internal/queue"
	"construction-ar-backend/internal/repository"
	"construction-ar-backend/internal/service"
	"construction-ar-backend/internal/twenty"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.Environment)
	defer logger.Log.Sync()

	logger.Log.Info("Starting BuildingMesh Backend", zap.String("env", cfg.Environment))

	repo, _ := repository.NewRepository(cfg.DBPath)
	defer repo.Close()
	runMigrations(repo, "migrations")

	twentyClient := twenty.NewClient(os.Getenv("TWENTY_API_URL"), os.Getenv("TWENTY_API_KEY"))
	projectSvc := service.NewProjectService(repo, twentyClient)
	storageSvc := service.NewLocalFileStorage("./uploads")

	redisCache, _ := cache.NewCache(cfg.RedisURL)
	rabbitQueue, _ := queue.NewQueue(cfg.RabbitMQURL)
	rabbitQueue.StartWorker(context.Background())

	h := handlers.NewHandler(repo, rabbitQueue, redisCache, projectSvc, storageSvc)
	r := chi.NewRouter()

	r.Use(middleware.Recovery)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	limiter := middleware.NewIPRateLimiter(1, 5)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Get("/health/live", h.HealthLive)
		r.Get("/health/ready", h.HealthReady)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Get("/me", h.Me)

			// Admin/SaaS
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("Superadmin"))
				r.Get("/organizations", h.GetOrganizations)
				r.Post("/organizations", h.CreateOrganization)
				r.Get("/saas/iot/stats", h.GetIoTStats)
			})

			// Business Logic
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("Superadmin", "Admin", "Engineer"))
				r.Get("/projects", h.GetProjects)
				r.Post("/projects", h.CreateProject)
				r.Get("/organizations/{orgID}/usage", h.GetOrgUsage)
				r.Get("/projects/{projectID}/elements", h.GetBIMElements)
				r.Post("/commits", h.CreateInspectionCommit)
				r.Post("/projects/{projectID}/splats", h.UploadSplat)
				r.Get("/projects/{projectID}/export", h.ExportProjectPDF)
			})
		})
	})

	webDir := filepath.Join(".", "backend/cmd/server/web")
	if _, err := os.Stat(webDir); os.IsNotExist(err) { webDir = filepath.Join(".", "cmd/server/web") }

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, ".") {
			http.FileServer(http.Dir(webDir)).ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})

	srv := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		logger.Log.Info("Server listening on :8080")
		srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	srv.Shutdown(context.Background())
}

func runMigrations(repo *repository.Repository, migrationDir string) {
	driver, _ := sqlite.WithInstance(repo.DB(), &sqlite.Config{})
	m, _ := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationDir), "sqlite", driver)
	m.Up()
}
