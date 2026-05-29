package service

import (
	"construction-ar-backend/internal/models"
	"construction-ar-backend/internal/repository"
	"github.com/google/uuid"
	"time"
	"errors"
)

type BIMService struct {
	repo *repository.Repository
}

func NewBIMService(repo *repository.Repository) *BIMService {
	return &BIMService{repo: repo}
}

func (s *BIMService) ProcessIFC(projectID string, filePath string) error {
	if projectID == "" || filePath == "" {
		return errors.New("missing project ID or file path")
	}

	// Mock Extraction Logic
	elementID := uuid.New().String()
	err := s.repo.CreateBIMElement(models.BIMElement{
		ID:           elementID,
		ProjectID:    projectID,
		ExternalGUID: "IFC-" + uuid.New().String()[:8],
		Name:         "Structural Wall B2",
		ElementType:  "Wall",
		Level:        "Floor 2",
		PositionJSON: `{"x": 0, "y": 1.5, "z": 0}`,
		CreatedAt:    time.Now(),
	})

	return err
}
