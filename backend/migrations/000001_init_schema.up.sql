-- 000001_init_schema.up.sql
CREATE TABLE IF NOT EXISTS organizations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL,
    name TEXT NOT NULL,
    address TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

CREATE TABLE IF NOT EXISTS walls (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT,
    status TEXT,
    thickness REAL,
    pos_x REAL,
    pos_y REAL,
    pos_z REAL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE TABLE IF NOT EXISTS layers (
    id TEXT PRIMARY KEY,
    wall_id TEXT NOT NULL,
    name TEXT NOT NULL,
    material TEXT,
    thickness REAL,
    "order" INTEGER,
    FOREIGN KEY (wall_id) REFERENCES walls(id)
);

CREATE TABLE IF NOT EXISTS splats (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT,
    file_path TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE TABLE IF NOT EXISTS floor_plans (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT,
    image_path TEXT,
    scale REAL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE TABLE IF NOT EXISTS issues (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    x REAL NOT NULL,
    y REAL NOT NULL,
    z REAL NOT NULL,
    status TEXT NOT NULL,
    priority TEXT NOT NULL,
    description TEXT,
    creator TEXT,
    assignee TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);
