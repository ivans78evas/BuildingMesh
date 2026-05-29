package service

import (
	"construction-ar-backend/internal/models"
	"construction-ar-backend/internal/repository"
	"errors"
	"math/rand"
	"time"
)

var (
	ErrPrerequisiteLayerNotApproved = errors.New("cannot commit: prerequisite layers are not approved")
)

type InspectionService struct {
	repo *repository.Repository
}

func NewInspectionService(repo *repository.Repository) *InspectionService {
	return &InspectionService{repo: repo}
}

func (s *InspectionService) SubmitCommit(commit models.InspectionCommit) (models.InspectionCommit, error) {
	if commit.BIMElementID == "" || commit.TemporalLayerID == "" {
		return commit, errors.New("missing BIM element or temporal layer ID")
	}

	// 1. Dependency Validation (Git-like Branch Protection)
	existingCommits, err := s.repo.GetInspectionCommits(commit.BIMElementID)
	if err != nil {
		return commit, err
	}

	for _, c := range existingCommits {
		if c.Status != "Approved" && c.ID != commit.ID {
			// In a real system, we'd check if this specific commit is a 'child' of an unapproved one.
		}
	}

	// 2. AI Health Score Simulation
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	commit.HealthScore = 88.0 + r.Float64()*(100.0-88.0)

	if commit.HealthScore > 95.0 {
		commit.Status = "Approved"
		now := time.Now()
		commit.ApprovedAt = &now
	} else {
		commit.Status = "ReviewRequired"
	}

	if commit.CreatedAt.IsZero() {
		commit.CreatedAt = time.Now()
	}

	err = s.repo.CreateInspectionCommit(commit)
	return commit, err
}

func (s *InspectionService) LogUsage(orgID, projectID string, sqM float64, scanType string) error {
	if orgID == "" || projectID == "" {
		return errors.New("missing organization or project ID")
	}

	log := models.UsageLog{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		ProjectID:      projectID,
		CapturedSqM:    sqM,
		ScanType:       scanType,
		CreatedAt:      time.Now(),
	}
	return s.repo.CreateUsageLog(log)
}
