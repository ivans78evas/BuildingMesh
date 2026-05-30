package models

import (
	"time"
)

// Organization represents a tenant in the SaaS system
type Organization struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Domain       string    `json:"domain"`
	Plan         string    `json:"plan"`   // Free, Pro, Enterprise
	Status       string    `json:"status"` // Active, Suspended, Trial
	MaxProjects  int       `json:"max_projects"`
	MaxUsers     int       `json:"max_users"`
	MaxIoTHubs   int       `json:"max_iot_hubs"`
	PrimaryColor string    `json:"primary_color"`
	LogoURL      string    `json:"logo_url"`
	CreatedAt    time.Time `json:"created_at"`
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

type BIMElement struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"` // Strict isolation
	ProjectID      string    `json:"project_id"`
	ExternalGUID   string    `json:"external_guid"`
	Name           string    `json:"name"`
	ElementType    string    `json:"element_type"`
	Level          string    `json:"level"`
	PositionJSON   string    `json:"position_json"`
	CreatedAt      time.Time `json:"created_at"`
}

type DesignLayer struct {
	ID             string  `json:"id"`
	OrganizationID string  `json:"organization_id"`
	BIMElementID   string  `json:"bim_element_id"`
	Name           string  `json:"name"`
	Material       string  `json:"material"`
	Thickness      float64 `json:"thickness"`
	Order          int     `json:"order"`
}

type TemporalLayer struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	BIMElementID   string    `json:"bim_element_id"`
	Name           string    `json:"name"`
	ScanDate       time.Time `json:"scan_date"`
	Type           string    `json:"type"`    // LiDAR, CSI
	S3Path         string    `json:"s3_path"`
	MetadataJSON   string    `json:"metadata_json"`
	CreatedAt      time.Time `json:"created_at"`
}

type InspectionCommit struct {
	ID              string     `json:"id"`
	OrganizationID  string     `json:"organization_id"`
	BIMElementID    string     `json:"bim_element_id"`
	TemporalLayerID string     `json:"temporal_layer_id"`
	HealthScore     float64    `json:"health_score"`
	Status          string     `json:"status"` // Draft, Approved, Read-Only
	Comments        string     `json:"comments"`
	InspectorID     string     `json:"inspector_id"`
	ApprovedAt      *time.Time `json:"approved_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type UsageLog struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	ProjectID      string    `json:"project_id"`
	CapturedSqM    float64   `json:"captured_sq_m"`
	ScanType       string    `json:"scan_type"`
	CreatedAt      time.Time `json:"created_at"`
}

type Splat struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	ProjectID      string    `json:"project_id"`
	Name           string    `json:"name"`
	FilePath       string    `json:"file_path"`
	CreatedAt      time.Time `json:"created_at"`
}

type Issue struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	ProjectID      string    `json:"project_id"`
	X, Y, Z        float64   `json:"x"`
	Status         string    `json:"status"`
	Priority       string    `json:"priority"`
	Description    string    `json:"description"`
	Creator        string    `json:"creator"`
	Assignee       string    `json:"assignee"`
	CreatedAt      time.Time `json:"created_at"`
}

// Legacy/Compatibility
type Wall struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	ProjectID      string    `json:"project_id"`
	Name           string    `json:"name"`
	Thickness      float64   `json:"thickness"`
	CreatedAt      time.Time `json:"created_at"`
}

type IoTStats struct {
	RPS           float64 `json:"rps"`
	ActiveDevices int     `json:"active_devices"`
	TotalEvents   int64   `json:"total_events"`
	ErrorCount    int     `json:"error_count"`
}

type SystemLog struct {
	ID        string    `json:"id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}
