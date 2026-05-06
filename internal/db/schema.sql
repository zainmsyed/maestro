CREATE TABLE IF NOT EXISTS epics (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    owner TEXT NOT NULL DEFAULT '',
    sprint_start TEXT NOT NULL DEFAULT '',
    sprint_end TEXT NOT NULL DEFAULT '',
    original_end_date TEXT,
    committed_end_date TEXT,
    actual_end_date TEXT,
    is_synthetic INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS features (
    id TEXT PRIMARY KEY,
    epic_id TEXT,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    owner TEXT NOT NULL DEFAULT '',
    sprint TEXT NOT NULL DEFAULT '',
    original_end_date TEXT,
    committed_end_date TEXT,
    actual_end_date TEXT,
    story_points INTEGER,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (epic_id) REFERENCES epics(id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS sprints (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    start_date TEXT,
    end_date TEXT,
    team TEXT NOT NULL DEFAULT '',
    capacity INTEGER,
    source TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS date_audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    changed_by TEXT NOT NULL,
    old_date TEXT,
    new_date TEXT,
    delta_days INTEGER NOT NULL,
    reason TEXT,
    changed_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_features_epic_id ON features(epic_id);
CREATE INDEX IF NOT EXISTS idx_features_sprint ON features(sprint);
CREATE INDEX IF NOT EXISTS idx_epics_sprint_start ON epics(sprint_start);
CREATE INDEX IF NOT EXISTS idx_epics_sprint_end ON epics(sprint_end);
CREATE INDEX IF NOT EXISTS idx_sprints_name ON sprints(name);
CREATE INDEX IF NOT EXISTS idx_date_audit_logs_entity ON date_audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_date_audit_logs_changed_at ON date_audit_logs(changed_at);
