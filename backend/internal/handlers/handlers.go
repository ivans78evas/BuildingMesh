package handlers

import (
	"encoding/json"
	"net/http"
	"construction-ar-backend/internal/models"
	"construction-ar-backend/internal/repository"
	"construction-ar-backend/internal/queue"
	"construction-ar-backend/internal/cache"
	"construction-ar-backend/internal/service"
	"construction-ar-backend/internal/middleware"
	"construction-ar-backend/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"time"
	"strings"
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

func getClaims(r *http.Request) *auth.Claims {
	if claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims); ok {
		return claims
	}
	return nil
}

func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	claims := getClaims(r)
	projects, _ := h.repo.GetProjects(claims.OrganizationID)
	json.NewEncoder(w).Encode(projects)
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	claims := getClaims(r)
	var p models.Project
	json.NewDecoder(r.Body).Decode(&p)
	p.ID = uuid.New().String()
	p.OrganizationID = claims.OrganizationID
	p.CreatedAt = time.Now()
	res, _ := h.projectSvc.CreateProject(p)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetBIMElements(w http.ResponseWriter, r *http.Request) {
	claims := getClaims(r)
	projectID := chi.URLParam(r, "projectID")
	e, _ := h.repo.GetBIMElements(claims.OrganizationID, projectID)
	json.NewEncoder(w).Encode(e)
}

func (h *Handler) CreateInspectionCommit(w http.ResponseWriter, r *http.Request) {
	claims := getClaims(r)
	var c models.InspectionCommit
	json.NewDecoder(r.Body).Decode(&c)
	c.ID = uuid.New().String()
	c.OrganizationID = claims.OrganizationID
	c.InspectorID = claims.UserID
	res, _ := h.inspectSvc.SubmitCommit(c)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) UploadSplat(w http.ResponseWriter, r *http.Request) {
	claims := getClaims(r)
	projectID := chi.URLParam(r, "projectID")
	r.ParseMultipartForm(32 << 20)
	file, header, _ := r.FormFile("file")
	defer file.Close()
	fileID := uuid.New().String()
	s3Path, _ := h.storageSvc.Upload(file, fileID+"_"+header.Filename)
	s := models.Splat{
		ID: fileID, OrganizationID: claims.OrganizationID, ProjectID: projectID,
		Name: r.FormValue("name"), FilePath: s3Path, CreatedAt: time.Now(),
	}
	h.repo.CreateSplat(s)
	h.queue.Publish(r.Context(), s)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(s)
}

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

func (h *Handler) GetOrgUsage(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	total, _ := h.repo.GetTotalUsage(orgID)
	json.NewEncoder(w).Encode(map[string]interface{}{"org": orgID, "sqm": total})
}

func (h *Handler) HealthLive(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) }
func (h *Handler) HealthReady(w http.ResponseWriter, r *http.Request) { w.Write([]byte("READY")) }

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

func (h *Handler) GetIoTStats(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(models.IoTStats{RPS: 124.5}) }
func (h *Handler) ExportProjectPDF(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PDF DATA")) }
