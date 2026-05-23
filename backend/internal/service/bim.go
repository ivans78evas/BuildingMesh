package service

import (
	"construction-ar-backend/internal/models"
	"construction-ar-backend/internal/repository"
	"github.com/google/uuid"
)

type BIMService struct {
	repo *repository.Repository
}

func NewBIMService(repo *repository.Repository) *BIMService {
	return &BIMService{repo: repo}
}

// ProcessIFC имитирует парсинг IFC файла и сохранение данных в БД.
// В реальном приложении здесь был бы вызов ifcopenshell или аналогичной библиотеки.
func (s *BIMService) ProcessIFC(projectID string, filePath string) error {
	// 1. Извлечение геометрии стен
	// 2. Извлечение слоев
	// 3. Извлечение инсталляций

	// Пример добавления тестовых данных после "парсинга"
	wallID := uuid.New().String()
	err := s.repo.CreateWall(models.Wall{
		ID:        wallID,
		ProjectID: projectID,
		Name:      "Wall-001",
		Thickness: 0.3,
		PositionX: 0,
		PositionY: 0,
		PositionZ: 0,
	})
	if err != nil {
		return err
	}

	layers := []models.Layer{
		{ID: uuid.New().String(), WallID: wallID, Name: "Plaster", Material: "Gypsum", Thickness: 0.02, Order: 1},
		{ID: uuid.New().String(), WallID: wallID, Name: "Brick", Material: "Red Brick", Thickness: 0.25, Order: 2},
		{ID: uuid.New().String(), WallID: wallID, Name: "Finish", Material: "Paint", Thickness: 0.01, Order: 3},
	}

	for _, l := range layers {
		if err := s.repo.CreateLayer(l); err != nil {
			return err
		}
	}

	return nil
}
