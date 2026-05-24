package main

import (
	"construction-ar-backend/internal/handlers"
	"construction-ar-backend/internal/repository"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	repo, err := repository.NewRepository("construction.db")
	if err != nil {
		log.Fatal(err)
	}

	h := handlers.NewHandler(repo)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API эндпоинты
	r.Route("/api/v1", func(r chi.Router) {
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

	// Раздача статики для веб-морды (SPA)
	workDir, _ := os.Getwd()
	// Determine correct web directory path
	webDir := filepath.Join(workDir, "cmd/server/web")
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		webDir = filepath.Join(workDir, "backend/cmd/server/web")
	}

	filesDir := http.Dir(webDir)

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Если запрашивается файл с расширением (js, css и т.д.), отдаем его
		if strings.Contains(r.URL.Path, ".") {
			http.FileServer(filesDir).ServeHTTP(w, r)
			return
		}
		// Для всех остальных путей (SPA роутинг) отдаем index.html
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})

	log.Println("Starting server on :8080...")
	log.Printf("Web directory: %s", webDir)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
