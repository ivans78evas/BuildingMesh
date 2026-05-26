package models

import "time"


type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Plan      string    `json:"plan"` // Free, Pro, Enterprise
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Email          string    `json:"email"`
	FullName       string    `json:"full_name"`
	Role           string    `json:"role"` // Superadmin, Admin, Engineer, Viewer
	CreatedAt      time.Time `json:"created_at"`
}

type Project struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	CreatedAt      time.Time `json:"created_at"`
}

type Wall struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Thickness float64   `json:"thickness"`
	PositionX float64   `json:"pos_x"`
	PositionY float64   `json:"pos_y"`
	PositionZ float64   `json:"pos_z"`
	CreatedAt time.Time `json:"created_at"`
}

type Layer struct {
	ID        string  `json:"id"`
	WallID    string  `json:"wall_id"`
	Name      string  `json:"name"`
	Material  string  `json:"material"`
	Thickness float64 `json:"thickness"`
	Order     int     `json:"order"`
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
	Scale     float64   `json:"scale"`
	CreatedAt time.Time `json:"created_at"`
}

type Issue struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	X           float64   `json:"x"`
	Y           float64   `json:"y"`
	Z           float64   `json:"z"`
	Status      string    `json:"status"`   // Open, InProgress, Resolved, Closed
	Priority    string    `json:"priority"` // Low, Medium, High
	Description string    `json:"description"`
	Creator     string    `json:"creator"`
	Assignee    string    `json:"assignee"`
	CreatedAt   time.Time `json:"created_at"`
}
