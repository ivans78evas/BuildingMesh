package models

import "time"

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	CreatedAt   time.Time `json:"created_at"`
}

type Wall struct {
	ID        string  `json:"id"`
	ProjectID string  `json:"project_id"`
	Name      string  `json:"name"`
	Thickness float64 `json:"thickness"`
	// Координаты относительно опорной точки (щитка)
	PositionX float64 `json:"pos_x"`
	PositionY float64 `json:"pos_y"`
	PositionZ float64 `json:"pos_z"`
}

type Layer struct {
	ID        string  `json:"id"`
	WallID    string  `json:"wall_id"`
	Name      string  `json:"name"`
	Material  string  `json:"material"`
	Thickness float64 `json:"thickness"`
	Order     int     `json:"order"` // Порядок слоя внутри стены
}

type Installation struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Type      string `json:"type"` // e.g., "pipe", "cable"
	Name      string `json:"name"`
	// Геометрия в упрощенном виде или ссылка на GLB
	GLBPath string `json:"glb_path"`
}

type Splat struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}

type FloorPlan struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	ImagePath string    `json:"image_path"`
	Scale     float64   `json:"scale"` // пикселей на метр
	CreatedAt time.Time `json:"created_at"`
}
