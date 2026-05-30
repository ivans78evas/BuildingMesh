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
)

type Handler struct {
	repo       *repository.Repository
	queue      *queue.Queue
	cache      *cache.Cache
	inspectSvc *service.InspectionService
	projectSvc *service.ProjectService
	storageSvc service.S3Service
}

func NewHandler(repo *repository.Repository, q *queue.Queue, c *cache.Cache, pSvc *service.ProjectService, sSvc service.S3Service) *Handler {
	return &Handler{
		repo:       repo,
		queue:      q,
		cache:      c,
		inspectSvc: service.NewInspectionService(repo),
		projectSvc: pSvc,
		storageSvc: sSvc,
	}
}

// --- High-Load Optimized Ingestion ---

func (h *Handler) UploadSplat(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")

	// Use 32MB buffer for multipart parsing
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 1. Offload heavy binary to storage (S3/Local) IMMEDIATELY
	// This keeps the binary blob out of SQLite.
	fileID := uuid.New().String()
	fileName := fileID + "_" + header.Filename
	s3Path, err := h.storageSvc.Upload(file, fileName)
	if err != nil {
		http.Error(w, "Storage failure", http.StatusInternalServerError)
		return
	}

	// 2. Create lightweight record in SQLite
	s := models.Splat{
		ID:        fileID,
		ProjectID: projectID,
		Name:      r.FormValue("name"),
		FilePath:  s3Path,
		CreatedAt: time.Now(),
	}

	if err := h.repo.CreateSplat(s); err != nil {
		http.Error(w, "DB failure", http.StatusInternalServerError)
		return
	}

	// 3. Trigger async processing via RabbitMQ (decimation, etc.)
	// This dampers the peak load on the CPU.
	h.queue.Publish(r.Context(), s)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(s)
}

// --- Rest of Handlers ---

func (h *Handler) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	orgs, _ := h.repo.GetOrganizations()
	json.NewEncoder(w).Encode(orgs)
}

func (h *Handler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var o models.Organization
	json.NewDecoder(r.Body).Decode(&o)
	o.ID = uuid.New().String()
	o.CreatedAt = time.Now()
	h.repo.CreateOrganization(o)
	json.NewEncoder(w).Encode(o)
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var p models.Project
	json.NewDecoder(r.Body).Decode(&p)
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	res, _ := h.projectSvc.CreateProject(p)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	p, _ := h.repo.GetProjects()
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) GetOrgUsage(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	total, _ := h.repo.GetTotalUsage(orgID)
	json.NewEncoder(w).Encode(map[string]interface{}{"org": orgID, "sqm": total})
}

func (h *Handler) GetBIMElements(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	e, _ := h.repo.GetBIMElements(projectID)
	json.NewEncoder(w).Encode(e)
}

func (h *Handler) CreateInspectionCommit(w http.ResponseWriter, r *http.Request) {
	var c models.InspectionCommit
	json.NewDecoder(r.Body).Decode(&c)
	c.ID = uuid.New().String()
	res, _ := h.inspectSvc.SubmitCommit(c)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) ExportProjectPDF(w http.ResponseWriter, r *http.Request) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "High-Load Audit Report")
	pdf.Output(w)
}

func (h *Handler) GetIoTStats(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(models.IoTStats{RPS: 124.5})
}

func (h *Handler) HealthLive(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *Handler) HealthReady(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("READY"))
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	u, _ := h.repo.GetUsersByOrg(orgID)
	json.NewEncoder(w).Encode(u)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	json.NewDecoder(r.Body).Decode(&u)
	u.ID = uuid.New().String()
	h.repo.CreateUser(u)
	json.NewEncoder(w).Encode(u)
}
