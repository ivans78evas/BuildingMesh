package repository

import (
	"construction-ar-backend/internal/models"
	"database/sql"
	"errors"
	"time"

	_ "modernc.org/sqlite"
)

var (
	ErrApprovedCommitReadOnly = errors.New("cannot modify an approved inspection commit")
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

// --- Multi-tenant Filtered Methods ---

func (r *Repository) GetProjects(orgID string) ([]models.Project, error) {
	rows, err := r.db.Query("SELECT id, organization_id, name, address, created_at FROM projects WHERE organization_id = ?", orgID)
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

func (r *Repository) GetBIMElements(orgID, projectID string) ([]models.BIMElement, error) {
	rows, err := r.db.Query("SELECT id, organization_id, project_id, external_guid, name, element_type, level, position_json, created_at FROM bim_elements WHERE organization_id = ? AND project_id = ?", orgID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elements []models.BIMElement
	for rows.Next() {
		var e models.BIMElement
		if err := rows.Scan(&e.ID, &e.OrganizationID, &e.ProjectID, &e.ExternalGUID, &e.Name, &e.ElementType, &e.Level, &e.PositionJSON, &e.CreatedAt); err != nil {
			return nil, err
		}
		elements = append(elements, e)
	}
	return elements, nil
}

func (r *Repository) GetInspectionCommits(orgID, elementID string) ([]models.InspectionCommit, error) {
	rows, err := r.db.Query("SELECT id, organization_id, bim_element_id, temporal_layer_id, health_score, status, comments, inspector_id, approved_at, created_at FROM inspection_commits WHERE organization_id = ? AND bim_element_id = ?", orgID, elementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commits []models.InspectionCommit
	for rows.Next() {
		var c models.InspectionCommit
		if err := rows.Scan(&c.ID, &c.OrganizationID, &c.BIMElementID, &c.TemporalLayerID, &c.HealthScore, &c.Status, &c.Comments, &c.InspectorID, &c.ApprovedAt, &c.CreatedAt); err != nil {
			return nil, err
		}
		commits = append(commits, c)
	}
	return commits, nil
}

func (r *Repository) CreateProject(p models.Project) error {
	_, err := r.db.Exec("INSERT INTO projects (id, organization_id, name, address, created_at) VALUES (?, ?, ?, ?, ?)",
		p.ID, p.OrganizationID, p.Name, p.Address, p.CreatedAt)
	return err
}

func (r *Repository) CreateOrganization(o models.Organization) error {
	_, err := r.db.Exec(`INSERT INTO organizations (id, name, domain, plan, status, max_projects, max_users, max_iot_hubs, primary_color, logo_url, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.ID, o.Name, o.Domain, o.Plan, o.Status, o.MaxProjects, o.MaxUsers, o.MaxIoTHubs, o.PrimaryColor, o.LogoURL, o.CreatedAt)
	return err
}

func (r *Repository) GetOrganizations() ([]models.Organization, error) {
	rows, err := r.db.Query("SELECT id, name, domain, plan, status, max_projects, max_users, max_iot_hubs, primary_color, logo_url, created_at FROM organizations")
	if err != nil { return nil, err }
	defer rows.Close()
	var orgs []models.Organization
	for rows.Next() {
		var o models.Organization
		rows.Scan(&o.ID, &o.Name, &o.Domain, &o.Plan, &o.Status, &o.MaxProjects, &o.MaxUsers, &o.MaxIoTHubs, &o.PrimaryColor, &o.LogoURL, &o.CreatedAt)
		orgs = append(orgs, o)
	}
	return orgs, nil
}

func (r *Repository) CreateBIMElement(e models.BIMElement) error {
	_, err := r.db.Exec(`INSERT INTO bim_elements (id, organization_id, project_id, external_guid, name, element_type, level, position_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.OrganizationID, e.ProjectID, e.ExternalGUID, e.Name, e.ElementType, e.Level, e.PositionJSON, e.CreatedAt)
	return err
}

func (r *Repository) CreateTemporalLayer(l models.TemporalLayer) error {
	_, err := r.db.Exec(`INSERT INTO temporal_layers (id, organization_id, bim_element_id, name, scan_date, type, s3_path, metadata_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		l.ID, l.OrganizationID, l.BIMElementID, l.Name, l.ScanDate, l.Type, l.S3Path, l.MetadataJSON, l.CreatedAt)
	return err
}

func (r *Repository) GetTemporalLayers(orgID, elementID string) ([]models.TemporalLayer, error) {
	rows, err := r.db.Query("SELECT id, organization_id, bim_element_id, name, scan_date, type, s3_path, metadata_json, created_at FROM temporal_layers WHERE organization_id = ? AND bim_element_id = ? ORDER BY scan_date ASC", orgID, elementID)
	if err != nil { return nil, err }
	defer rows.Close()
	var layers []models.TemporalLayer
	for rows.Next() {
		var l models.TemporalLayer
		rows.Scan(&l.ID, &l.OrganizationID, &l.BIMElementID, &l.Name, &l.ScanDate, &l.Type, &l.S3Path, &l.MetadataJSON, &l.CreatedAt)
		layers = append(layers, l)
	}
	return layers, nil
}

func (r *Repository) CreateInspectionCommit(c models.InspectionCommit) error {
	_, err := r.db.Exec(`INSERT INTO inspection_commits (id, organization_id, bim_element_id, temporal_layer_id, health_score, status, comments, inspector_id, approved_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.OrganizationID, c.BIMElementID, c.TemporalLayerID, c.HealthScore, c.Status, c.Comments, c.InspectorID, c.ApprovedAt, c.CreatedAt)
	return err
}

func (r *Repository) UpdateInspectionCommit(c models.InspectionCommit) error {
	var approvedAt *time.Time
	err := r.db.QueryRow("SELECT approved_at FROM inspection_commits WHERE id = ? AND organization_id = ?", c.ID, c.OrganizationID).Scan(&approvedAt)
	if err != nil { return err }
	if approvedAt != nil { return ErrApprovedCommitReadOnly }
	_, err = r.db.Exec(`UPDATE inspection_commits SET health_score = ?, status = ?, comments = ?, inspector_id = ?, approved_at = ? WHERE id = ? AND organization_id = ?`,
		c.HealthScore, c.Status, c.Comments, c.InspectorID, c.ApprovedAt, c.ID, c.OrganizationID)
	return err
}

func (r *Repository) CreateUsageLog(l models.UsageLog) error {
	_, err := r.db.Exec("INSERT INTO usage_logs (id, organization_id, project_id, captured_sq_m, scan_type, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		l.ID, l.OrganizationID, l.ProjectID, l.CapturedSqM, l.ScanType, l.CreatedAt)
	return err
}

func (r *Repository) GetTotalUsage(orgID string) (float64, error) {
	var total float64
	err := r.db.QueryRow("SELECT COALESCE(SUM(captured_sq_m), 0) FROM usage_logs WHERE organization_id = ?", orgID).Scan(&total)
	return total, err
}

func (r *Repository) CreateIssue(i models.Issue) error {
	_, err := r.db.Exec("INSERT INTO issues (id, organization_id, project_id, x, y, z, status, priority, description, creator, assignee, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		i.ID, i.OrganizationID, i.ProjectID, i.X, i.Y, i.Z, i.Status, i.Priority, i.Description, i.Creator, i.Assignee, i.CreatedAt)
	return err
}

func (r *Repository) GetIssues(orgID, projectID string) ([]models.Issue, error) {
	rows, err := r.db.Query("SELECT id, organization_id, project_id, x, y, z, status, priority, description, creator, assignee, created_at FROM issues WHERE organization_id = ? AND project_id = ?", orgID, projectID)
	if err != nil { return nil, err }
	defer rows.Close()
	var issues []models.Issue
	for rows.Next() {
		var i models.Issue
		rows.Scan(&i.ID, &i.OrganizationID, &i.ProjectID, &i.X, &i.Y, &i.Z, &i.Status, &i.Priority, &i.Description, &i.Creator, &i.Assignee, &i.CreatedAt)
		issues = append(issues, i)
	}
	return issues, nil
}

func (r *Repository) CreateSplat(s models.Splat) error {
	_, err := r.db.Exec("INSERT INTO splats (id, organization_id, project_id, name, file_path, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		s.ID, s.OrganizationID, s.ProjectID, s.Name, s.FilePath, s.CreatedAt)
	return err
}

func (r *Repository) CreateUser(u models.User) error {
	_, err := r.db.Exec("INSERT INTO users (id, organization_id, email, full_name, role, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		u.ID, u.OrganizationID, u.Email, u.FullName, u.Role, u.CreatedAt)
	return err
}

func (r *Repository) GetUsersByOrg(orgID string) ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, organization_id, email, full_name, role, created_at FROM users WHERE organization_id = ?", orgID)
	if err != nil { return nil, err }
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.OrganizationID, &u.Email, &u.FullName, &u.Role, &u.CreatedAt)
		users = append(users, u)
	}
	return users, nil
}

func (r *Repository) Ping() error { return r.db.Ping() }
func (r *Repository) Close() error { return r.db.Close() }
