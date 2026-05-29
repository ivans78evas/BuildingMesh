package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"construction-ar-backend/internal/models"
	"github.com/jung-kurt/gofpdf"
	"construction-ar-backend/internal/repository"
	"construction-ar-backend/internal/queue"
	"construction-ar-backend/internal/cache"
	"construction-ar-backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"time"
	"os"
	"io"
	"path/filepath"
)

type Handler struct {
	repo       *repository.Repository
	queue      *queue.Queue
	cache      *cache.Cache
	inspectSvc *service.InspectionService
	projectSvc *service.ProjectService
}

func NewHandler(repo *repository.Repository, q *queue.Queue, c *cache.Cache, pSvc *service.ProjectService) *Handler {
	return &Handler{
		repo:       repo,
		queue:      q,
		cache:      c,
		inspectSvc: service.NewInspectionService(repo),
		projectSvc: pSvc,
	}
}

// @Summary Get all organizations
func (h *Handler) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	orgs, err := h.repo.GetOrganizations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orgs)
}

// @Summary Create organization
func (h *Handler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var o models.Organization
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	o.ID = uuid.New().String()
	o.CreatedAt = time.Now()
	if o.Plan == "" { o.Plan = "Free" }
	if o.Status == "" { o.Status = "Active" }

	if err := h.repo.CreateOrganization(o); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(o)
}

// @Summary Create a new project
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()

	// Use the ProjectService which handles Twenty CRM syncing
	createdProject, err := h.projectSvc.CreateProject(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProject)
}

// ... (Other handlers kept as they were) ...

func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.repo.GetProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(projects)
}

func (h *Handler) GetOrgUsage(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	total, _ := h.repo.GetTotalUsage(orgID)
	json.NewEncoder(w).Encode(map[string]interface{}{"org": orgID, "sqm": total})
}

func (h *Handler) GetBIMElements(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	elements, _ := h.repo.GetBIMElements(projectID)
	json.NewEncoder(w).Encode(elements)
}

func (h *Handler) GetTemporalLayers(w http.ResponseWriter, r *http.Request) {
	elementID := chi.URLParam(r, "elementID")
	layers, _ := h.repo.GetTemporalLayers(elementID)
	json.NewEncoder(w).Encode(layers)
}

func (h *Handler) CreateTemporalLayer(w http.ResponseWriter, r *http.Request) {
	elementID := chi.URLParam(r, "elementID")
	var l models.TemporalLayer
	json.NewDecoder(r.Body).Decode(&l)
	l.ID = uuid.New().String()
	l.BIMElementID = elementID
	l.CreatedAt = time.Now()
	h.repo.CreateTemporalLayer(l)
	json.NewEncoder(w).Encode(l)
}

func (h *Handler) CreateInspectionCommit(w http.ResponseWriter, r *http.Request) {
	var c models.InspectionCommit
	json.NewDecoder(r.Body).Decode(&c)
	c.ID = uuid.New().String()
	final, _ := h.inspectSvc.SubmitCommit(c)
	json.NewEncoder(w).Encode(final)
}

func (h *Handler) UpdateInspectionCommit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var c models.InspectionCommit
	json.NewDecoder(r.Body).Decode(&c)
	c.ID = id
	h.repo.UpdateInspectionCommit(c)
	json.NewEncoder(w).Encode(c)
}

func (h *Handler) GetInspectionCommits(w http.ResponseWriter, r *http.Request) {
	elementID := chi.URLParam(r, "elementID")
	commits, _ := h.repo.GetInspectionCommits(elementID)
	json.NewEncoder(w).Encode(commits)
}

func (h *Handler) GetWallComparison(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "OK", "deviation": 2.4})
}

func (h *Handler) UploadSplat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetIssues(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	issues, _ := h.repo.GetIssues(projectID)
	json.NewEncoder(w).Encode(issues)
}

func (h *Handler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	var i models.Issue
	json.NewDecoder(r.Body).Decode(&i)
	i.ID = uuid.New().String()
	i.ProjectID = projectID
	h.repo.CreateIssue(i)
	json.NewEncoder(w).Encode(i)
}

func (h *Handler) PatchIssueStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ExportProjectPDF(w http.ResponseWriter, r *http.Request) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "BuildingMesh Report")
	pdf.Output(w)
}

func (h *Handler) GetIoTStats(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(models.IoTStats{RPS: 124.5})
}

func (h *Handler) GetSystemLogs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]models.SystemLog{})
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	users, _ := h.repo.GetUsersByOrg(orgID)
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	json.NewDecoder(r.Body).Decode(&u)
	u.ID = uuid.New().String()
	h.repo.CreateUser(u)
	json.NewEncoder(w).Encode(u)
}

func (h *Handler) HealthLive(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *Handler) HealthReady(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("READY"))
}
