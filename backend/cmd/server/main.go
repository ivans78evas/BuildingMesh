package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"construction-ar-backend/internal/handlers"
	"construction-ar-backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Поиск папки статики
	workDir, _ := os.Getwd()
	possiblePaths := []string{
		filepath.Join(workDir, "static", "web"),
		filepath.Join(workDir, "backend", "static", "web"),
		"/workspaces/BuildingMesh/backend/static/web",
	}

	staticDir := ""
	for _, p := range possiblePaths {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			staticDir = p
			break
		}
	}

	log.Printf("Starting server...")
	log.Printf("Working Directory: %s", workDir)
	log.Printf("Static Directory: %s", staticDir)

	repo, err := repository.NewRepository("construction.db")
	if err != nil {
		log.Fatal(err)
	}

	h := handlers.NewHandler(repo)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API Маршруты
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/organizations", h.GetOrganizations)
		r.Post("/organizations", h.CreateOrganization)
		r.Get("/projects", h.GetProjects)
		r.Post("/projects", h.CreateProject)
		r.Get("/walls/{wallID}", h.GetWall)
		r.Get("/projects/{projectID}/splats", h.GetSplats)
		r.Post("/projects/{projectID}/splats", h.UploadSplat)
		r.Get("/projects/{projectID}/plans", h.GetFloorPlans)
		r.Post("/projects/{projectID}/plans", h.UploadFloorPlan)
		r.Get("/sync", h.GetSyncChanges)
		r.Post("/sync", h.PushSyncChanges)
	})

	// Статика
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1") {
			return
		}

		path := filepath.Join(staticDir, r.URL.Path)
		if r.URL.Path == "/" || r.URL.Path == "" {
			path = filepath.Join(staticDir, "index.html")
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		http.ServeFile(w, r, path)
	})

	log.Println("Server running on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
