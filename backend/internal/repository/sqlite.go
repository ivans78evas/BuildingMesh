package repository

import (
	"database/sql"
	"construction-ar-backend/internal/models"
	_ "modernc.org/sqlite"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) DB() *sql.DB {
	return r.db
}

func (r *Repository) GetProjects() ([]models.Project, error) {
	rows, err := r.db.Query("SELECT id, organization_id, name, address, created_at FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.OrganizationID, &p.Name, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *Repository) CreateProject(p models.Project) error {
	_, err := r.db.Exec("INSERT INTO projects (id, organization_id, name, address, created_at) VALUES (?, ?, ?, ?, ?)",
		p.ID, p.OrganizationID, p.Name, p.Address, p.CreatedAt)
	return err
}

func (r *Repository) CreateOrganization(o models.Organization) error {
	_, err := r.db.Exec("INSERT INTO organizations (id, name, created_at) VALUES (?, ?, ?)",
		o.ID, o.Name, o.CreatedAt)
	return err
}

func (r *Repository) GetOrganizations() ([]models.Organization, error) {
	rows, err := r.db.Query("SELECT id, name, created_at FROM organizations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []models.Organization
	for rows.Next() {
		var o models.Organization
		if err := rows.Scan(&o.ID, &o.Name, &o.CreatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}
	return orgs, nil
}

func (r *Repository) CreateWall(w models.Wall) error {
	_, err := r.db.Exec("INSERT INTO walls (id, project_id, name, type, status, thickness, pos_x, pos_y, pos_z, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		w.ID, w.ProjectID, w.Name, w.Type, w.Status, w.Thickness, w.PositionX, w.PositionY, w.PositionZ, w.CreatedAt)
	return err
}

func (r *Repository) CreateLayer(l models.Layer) error {
	_, err := r.db.Exec("INSERT INTO layers (id, wall_id, name, material, thickness, \"order\") VALUES (?, ?, ?, ?, ?, ?)",
		l.ID, l.WallID, l.Name, l.Material, l.Thickness, l.Order)
	return err
}

func (r *Repository) GetWallDetails(wallID string) (*models.Wall, []models.Layer, error) {
	var w models.Wall
	err := r.db.QueryRow("SELECT id, project_id, name, type, status, thickness, pos_x, pos_y, pos_z, created_at FROM walls WHERE id = ?", wallID).
		Scan(&w.ID, &w.ProjectID, &w.Name, &w.Type, &w.Status, &w.Thickness, &w.PositionX, &w.PositionY, &w.PositionZ, &w.CreatedAt)
	if err != nil {
		return nil, nil, err
	}

	rows, err := r.db.Query("SELECT id, wall_id, name, material, thickness, \"order\" FROM layers WHERE wall_id = ? ORDER BY \"order\"", wallID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var layers []models.Layer
	for rows.Next() {
		var l models.Layer
		if err := rows.Scan(&l.ID, &l.WallID, &l.Name, &l.Material, &l.Thickness, &l.Order); err != nil {
			return nil, nil, err
		}
		layers = append(layers, l)
	}

	return &w, layers, nil
}

func (r *Repository) CreateSplat(s models.Splat) error {
	_, err := r.db.Exec("INSERT INTO splats (id, project_id, name, file_path, created_at) VALUES (?, ?, ?, ?, ?)",
		s.ID, s.ProjectID, s.Name, s.FilePath, s.CreatedAt)
	return err
}

func (r *Repository) GetSplats(projectID string) ([]models.Splat, error) {
	rows, err := r.db.Query("SELECT id, project_id, name, file_path, created_at FROM splats WHERE project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var splats []models.Splat
	for rows.Next() {
		var s models.Splat
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Name, &s.FilePath, &s.CreatedAt); err != nil {
			return nil, err
		}
		splats = append(splats, s)
	}
	return splats, nil
}

func (r *Repository) CreateFloorPlan(fp models.FloorPlan) error {
	_, err := r.db.Exec("INSERT INTO floor_plans (id, project_id, name, image_path, scale, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		fp.ID, fp.ProjectID, fp.Name, fp.ImagePath, fp.Scale, fp.CreatedAt)
	return err
}

func (r *Repository) GetFloorPlans(projectID string) ([]models.FloorPlan, error) {
	rows, err := r.db.Query("SELECT id, project_id, name, image_path, scale, created_at FROM floor_plans WHERE project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.FloorPlan
	for rows.Next() {
		var fp models.FloorPlan
		if err := rows.Scan(&fp.ID, &fp.ProjectID, &fp.Name, &fp.ImagePath, &fp.Scale, &fp.CreatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, fp)
	}
	return plans, nil
}

func (r *Repository) CreateIssue(i models.Issue) error {
	_, err := r.db.Exec("INSERT INTO issues (id, project_id, x, y, z, status, priority, description, creator, assignee, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		i.ID, i.ProjectID, i.X, i.Y, i.Z, i.Status, i.Priority, i.Description, i.Creator, i.Assignee, i.CreatedAt)
	return err
}

func (r *Repository) GetIssues(projectID string) ([]models.Issue, error) {
	rows, err := r.db.Query("SELECT id, project_id, x, y, z, status, priority, description, creator, assignee, created_at FROM issues WHERE project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []models.Issue
	for rows.Next() {
		var i models.Issue
		if err := rows.Scan(&i.ID, &i.ProjectID, &i.X, &i.Y, &i.Z, &i.Status, &i.Priority, &i.Description, &i.Creator, &i.Assignee, &i.CreatedAt); err != nil {
			return nil, err
		}
		issues = append(issues, i)
	}
	return issues, nil
}

func (r *Repository) UpdateIssueStatus(id string, status string) error {
	_, err := r.db.Exec("UPDATE issues SET status = ? WHERE id = ?", status, id)
	return err
}

func (r *Repository) Ping() error {
	return r.db.Ping()
}

func (r *Repository) Close() error {
	return r.db.Close()
}
