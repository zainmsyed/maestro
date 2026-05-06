package db_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	maestrodb "maestro/internal/db"
)

func TestOpenCreatesSchema(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "maestro.db")
	database, err := maestrodb.Open(dbPath)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()

	assertTableExists(t, database, "epics")
	assertTableExists(t, database, "features")
	assertTableExists(t, database, "sprints")
	assertTableExists(t, database, "date_audit_logs")

	assertIndexExists(t, database, "idx_features_epic_id")
	assertIndexExists(t, database, "idx_features_sprint")
	assertIndexExists(t, database, "idx_epics_sprint_start")
	assertIndexExists(t, database, "idx_epics_sprint_end")
	assertIndexExists(t, database, "idx_sprints_name")
	assertIndexExists(t, database, "idx_date_audit_logs_entity")
	assertIndexExists(t, database, "idx_date_audit_logs_changed_at")
}

func assertTableExists(t *testing.T, db *sql.DB, table string) {
	t.Helper()
	var name string
	if err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&name); err != nil {
		t.Fatalf("table %s not found: %v", table, err)
	}
}

func assertIndexExists(t *testing.T, db *sql.DB, index string) {
	t.Helper()
	var name string
	if err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'index' AND name = ?`, index).Scan(&name); err != nil {
		t.Fatalf("index %s not found: %v", index, err)
	}
}
