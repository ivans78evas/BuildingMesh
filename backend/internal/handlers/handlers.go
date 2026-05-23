package handlers

import (
	"encoding/json"
	"net/http"
	"construction-ar-backend/internal/models"
	"construction-ar-backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"time"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.repo.GetProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(projects)
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()

	if err := h.repo.CreateProject(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) GetWall(w http.ResponseWriter, r *http.Request) {
	wallID := chi.URLParam(r, "wallID")
	wall, layers, err := h.repo.GetWallDetails(wallID)
	if err != nil {
		http.Error(w, "Wall not found", http.StatusNotFound)
		return
	}

	response := struct {
		Wall   *models.Wall   `json:"wall"`
		Layers []models.Layer `json:"layers"`
	}{
		Wall:   wall,
		Layers: layers,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetSplats(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	splats, err := h.repo.GetSplats(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(splats)
}

func (h *Handler) UploadSplat(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	// В реальной реализации здесь была бы обработка multipart/form-data
	// Для прототипа имитируем сохранение
	s := models.Splat{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Name:      r.URL.Query().Get("name"),
		FilePath:  "/uploads/splats/" + uuid.New().String() + ".ply",
		CreatedAt: time.Now(),
	}

	if err := h.repo.CreateSplat(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(s)
}

func (h *Handler) GetFloorPlans(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	plans, err := h.repo.GetFloorPlans(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(plans)
}

func (h *Handler) UploadFloorPlan(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	fp := models.FloorPlan{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Name:      r.URL.Query().Get("name"),
		ImagePath: "/uploads/plans/" + uuid.New().String() + ".png",
		Scale:     100.0, // Default
		CreatedAt: time.Now(),
	}

	if err := h.repo.CreateFloorPlan(fp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(fp)
}
