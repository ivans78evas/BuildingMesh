package main

import (
	"log"
	"net/http"
	"construction-ar-backend/internal/handlers"
	"construction-ar-backend/internal/repository"
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

	// Раздача статических файлов Web-панели (fallback)
	r.Handle("/*", http.FileServer(http.Dir("./static/web")))

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
