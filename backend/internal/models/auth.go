package models

import "time"

type Organization struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	GoogleAPIKey string    `json:"google_api_key"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserRole string

const (
	RoleOwner      UserRole = "owner"
	RolePM         UserRole = "pm"
	RoleArchitect  UserRole = "architect"
	RoleForeman    UserRole = "foreman"
	RoleSpecialist UserRole = "specialist"
	RoleWorker     UserRole = "worker"
	RoleClient     UserRole = "client"
)

type User struct {
	ID    string   `json:"id"`
	OrgID string   `json:"org_id"`
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
}

type SyncChange struct {
	Table     string `json:"table"`
	Data      string `json:"data"` // JSON string
	Timestamp int64  `json:"timestamp"`
}
