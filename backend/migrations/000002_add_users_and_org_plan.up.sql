-- 000002_add_users_and_org_plan.up.sql
ALTER TABLE organizations ADD COLUMN plan TEXT DEFAULT 'Free';

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    full_name TEXT,
    role TEXT NOT NULL, -- Superadmin, Admin, Engineer, Viewer
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations(id)
);
