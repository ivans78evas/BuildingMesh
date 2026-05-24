package handlers

import (
	"encoding/json"
	"net/http"
	"construction-ar-backend/internal/models"
	"construction-ar-backend/internal/repository"
	"construction-ar-backend/internal/queue"
	"construction-ar-backend/internal/cache"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"time"
	"os"
	"io"
	"path/filepath"
)

type Handler struct {
	repo  *repository.Repository
	queue *queue.Queue
	cache *cache.Cache
}

func NewHandler(repo *repository.Repository, q *queue.Queue, c *cache.Cache) *Handler {
	return &Handler{repo: repo, queue: q, cache: c}
}

// @Summary Get all projects
// @Description Get list of all projects
// @Tags Projects
// @Produce json
// @Success 200 {array} models.Project
// @Router /projects [get]
func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.repo.GetProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(projects)
}

// @Summary Create a new project
// @Description Create a new construction project
// @Tags Projects
// @Accept json
// @Produce json
// @Param project body models.Project true "Project object"
// @Success 201 {object} models.Project
// @Router /projects [post]
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

// @Summary Get wall details
// @Description Get details of a specific wall including layers
// @Tags Walls
// @Produce json
// @Param wallID path string true "Wall ID"
// @Success 200 {object} object
// @Router /walls/{wallID} [get]
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

// @Summary Upload point cloud (.ply)
// @Description Upload a 3D scan for a project. Returns 202 if accepted for background processing.
// @Tags 3D
// @Accept multipart/form-data
// @Param projectID path string true "Project ID"
// @Param file formData file true "Point cloud file"
// @Param name formData string true "Splat name"
// @Success 202 {object} models.Splat
// @Router /projects/{projectID}/splats [post]
func (h *Handler) UploadSplat(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")

	err := r.ParseMultipartForm(100 << 20) // 100 MB limit
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "./uploads/splats"
	os.MkdirAll(uploadDir, os.ModePerm)

	fileID := uuid.New().String()
	filePath := filepath.Join(uploadDir, fileID + filepath.Ext(header.Filename))

	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := models.Splat{
		ID:        fileID,
		ProjectID: projectID,
		Name:      r.FormValue("name"),
		FilePath:  filePath,
		CreatedAt: time.Now(),
	}

	if err := h.repo.CreateSplat(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Push to RabbitMQ for processing
	h.queue.Publish(r.Context(), s)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(s)
}

// @Summary Create a spatial issue
// @Description Mark a defect or issue in 3D space
// @Tags Issues
// @Accept json
// @Produce json
// @Param projectID path string true "Project ID"
// @Param issue body models.Issue true "Issue object"
// @Success 201 {object} models.Issue
// @Router /projects/{projectID}/issues [post]
func (h *Handler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	var i models.Issue
	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	i.ID = uuid.New().String()
	i.ProjectID = projectID
	i.CreatedAt = time.Now()

	if err := h.repo.CreateIssue(i); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(i)
}

// @Summary Get project issues
// @Description Get all spatial issues for a project
// @Tags Issues
// @Produce json
// @Param projectID path string true "Project ID"
// @Success 200 {array} models.Issue
// @Router /projects/{projectID}/issues [get]
func (h *Handler) GetIssues(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	issues, err := h.repo.GetIssues(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(issues)
}

// @Summary Update issue status
// @Description Quick status change for an issue
// @Tags Issues
// @Param id path string true "Issue ID"
// @Param status query string true "New Status"
// @Success 200
// @Router /issues/{id}/status [patch]
func (h *Handler) PatchIssueStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	status := r.URL.Query().Get("status")
	if err := h.repo.UpdateIssueStatus(id, status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HealthLive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) HealthReady(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.Ping(); err != nil {
		http.Error(w, "DB not ready", http.StatusServiceUnavailable)
		return
	}
	if err := h.cache.Ping(r.Context()); err != nil {
		http.Error(w, "Redis not ready", http.StatusServiceUnavailable)
		return
	}
	if err := h.queue.Ping(); err != nil {
		http.Error(w, "RabbitMQ not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("READY"))
}
