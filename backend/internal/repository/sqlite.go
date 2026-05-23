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

	repo := &Repository{db: db}
	if err := repo.initSchema(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		name TEXT,
		address TEXT,
		created_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS walls (
		id TEXT PRIMARY KEY,
		project_id TEXT,
		name TEXT,
		thickness REAL,
		pos_x REAL,
		pos_y REAL,
		pos_z REAL,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);

	CREATE TABLE IF NOT EXISTS layers (
		id TEXT PRIMARY KEY,
		wall_id TEXT,
		name TEXT,
		material TEXT,
		thickness REAL,
		"order" INTEGER,
		FOREIGN KEY(wall_id) REFERENCES walls(id)
	);

	CREATE TABLE IF NOT EXISTS installations (
		id TEXT PRIMARY KEY,
		project_id TEXT,
		type TEXT,
		name TEXT,
		glb_path TEXT,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);

	CREATE TABLE IF NOT EXISTS splats (
		id TEXT PRIMARY KEY,
		project_id TEXT,
		name TEXT,
		file_path TEXT,
		created_at DATETIME,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);

	CREATE TABLE IF NOT EXISTS floor_plans (
		id TEXT PRIMARY KEY,
		project_id TEXT,
		name TEXT,
		image_path TEXT,
		scale REAL,
		created_at DATETIME,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);

	CREATE TABLE IF NOT EXISTS sync_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		client_id TEXT,
		table_name TEXT,
		record_id TEXT,
		operation TEXT, -- 'INSERT', 'UPDATE', 'DELETE'
		data TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := r.db.Exec(schema)
	return err
}

func (r *Repository) GetProjects() ([]models.Project, error) {
	rows, err := r.db.Query("SELECT id, name, address, created_at FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *Repository) CreateProject(p models.Project) error {
	_, err := r.db.Exec("INSERT INTO projects (id, name, address, created_at) VALUES (?, ?, ?, ?)",
		p.ID, p.Name, p.Address, p.CreatedAt)
	return err
}

func (r *Repository) GetWallDetails(wallID string) (*models.Wall, []models.Layer, error) {
	var w models.Wall
	err := r.db.QueryRow("SELECT id, project_id, name, thickness, pos_x, pos_y, pos_z FROM walls WHERE id = ?", wallID).
		Scan(&w.ID, &w.ProjectID, &w.Name, &w.Thickness, &w.PositionX, &w.PositionY, &w.PositionZ)
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

func (r *Repository) CreateWall(w models.Wall) error {
	_, err := r.db.Exec("INSERT INTO walls (id, project_id, name, thickness, pos_x, pos_y, pos_z) VALUES (?, ?, ?, ?, ?, ?, ?)",
		w.ID, w.ProjectID, w.Name, w.Thickness, w.PositionX, w.PositionY, w.PositionZ)
	return err
}

func (r *Repository) CreateLayer(l models.Layer) error {
	_, err := r.db.Exec("INSERT INTO layers (id, wall_id, name, material, thickness, \"order\") VALUES (?, ?, ?, ?, ?, ?)",
		l.ID, l.WallID, l.Name, l.Material, l.Thickness, l.Order)
	return err
}

func (r *Repository) GetInstallations(projectID string) ([]models.Installation, error) {
	rows, err := r.db.Query("SELECT id, project_id, type, name, glb_path FROM installations WHERE project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var installations []models.Installation
	for rows.Next() {
		var i models.Installation
		if err := rows.Scan(&i.ID, &i.ProjectID, &i.Type, &i.Name, &i.GLBPath); err != nil {
			return nil, err
		}
		installations = append(installations, i)
	}
	return installations, nil
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
