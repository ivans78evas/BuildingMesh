package service

import (
	"construction-ar-backend/internal/models"
	"construction-ar-backend/internal/repository"
	"construction-ar-backend/internal/twenty"
)

type ProjectService struct {
	repo         *repository.Repository
	twentyClient *twenty.Client
}

func NewProjectService(repo *repository.Repository, tClient *twenty.Client) *ProjectService {
	return &ProjectService{
		repo:         repo,
		twentyClient: tClient,
	}
}

func (s *ProjectService) CreateProject(p models.Project) (models.Project, error) {
	// 1. Save core proprietary data to our DB
	if err := s.repo.CreateProject(p); err != nil {
		return p, err
	}

	// 2. Provision metadata mirror in Twenty (Enterprise Engine)
	// This makes the project searchable and manageable in the corporate CRM
	if s.twentyClient != nil {
		_, _ = s.twentyClient.CreateProjectInTwenty(p.Name, p.Address)
	}

	return p, nil
}
