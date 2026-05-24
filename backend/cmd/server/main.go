package main

import (
	"context"
	"construction-ar-backend/internal/cache"
	"construction-ar-backend/internal/config"
	"construction-ar-backend/internal/handlers"
	"construction-ar-backend/internal/logger"
	"construction-ar-backend/internal/middleware"
	"construction-ar-backend/internal/queue"
	"construction-ar-backend/internal/repository"
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
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.Environment)
	defer logger.Log.Sync()

	logger.Log.Info("Starting BuildingMesh Backend", zap.String("env", cfg.Environment))

	// Database & Migrations
	repo, err := repository.NewRepository(cfg.DBPath)
	if err != nil {
		logger.Log.Fatal("Failed to connect to DB", zap.Error(err))
	}
	defer repo.Close()

	runMigrations(repo, "migrations")

	// Cache (Redis)
	redisCache, err := cache.NewCache(cfg.RedisURL)
	if err != nil {
		logger.Log.Warn("Failed to connect to Redis", zap.Error(err))
	} else {
		defer redisCache.Close()
		logger.Log.Info("Connected to Redis")
	}

	// Queue (RabbitMQ)
	rabbitQueue, err := queue.NewQueue(cfg.RabbitMQURL)
	if err != nil {
		logger.Log.Warn("Failed to connect to RabbitMQ", zap.Error(err))
	} else {
		defer rabbitQueue.Close()
		logger.Log.Info("Connected to RabbitMQ")
		rabbitQueue.StartWorker(context.Background())
	}

	h := handlers.NewHandler(repo, rabbitQueue, redisCache)
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Recovery)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	// Rate limiting for uploads
	limiter := middleware.NewIPRateLimiter(1, 5) // 1 request per second, burst 5

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/projects", h.GetProjects)
		r.Post("/projects", h.CreateProject)
		r.Get("/walls/{wallID}", h.GetWall)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RateLimit(limiter))
			r.Post("/projects/{projectID}/splats", h.UploadSplat)
		})

		r.Get("/projects/{projectID}/issues", h.GetIssues)
		r.Post("/projects/{projectID}/issues", h.CreateIssue)
		r.Patch("/issues/{id}/status", h.PatchIssueStatus)

		r.Get("/health/live", h.HealthLive)
		r.Get("/health/ready", h.HealthReady)
	})

	// Swagger
	r.Get("/api/docs/*", httpSwagger.WrapHandler)

	// Static & SPA
	workDir, _ := os.Getwd()
	webDir := filepath.Join(workDir, "cmd/server/web")
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		webDir = filepath.Join(workDir, "backend/cmd/server/web")
	}
	filesDir := http.Dir(webDir)

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, ".") {
			http.FileServer(filesDir).ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Graceful Shutdown
	go func() {
		logger.Log.Info("Server listening", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Listen error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("Server exited gracefully")
}

func runMigrations(repo *repository.Repository, migrationDir string) {
	driver, err := sqlite3.WithInstance(repo.DB(), &sqlite3.Config{})
	if err != nil {
		logger.Log.Fatal("Migration driver error", zap.Error(err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		"sqlite3", driver)
	if err != nil {
		logger.Log.Fatal("Migration init error", zap.Error(err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Fatal("Migration failed", zap.Error(err))
	}
	logger.Log.Info("Migrations applied successfully")
}
