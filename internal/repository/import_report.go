package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"maestro/internal/models"
)

type ImportReportRepository interface {
	Save(context.Context, *models.ImportReport) error
	GetLast(context.Context) (*models.ImportReport, error)
}

type SQLiteImportReportRepository struct{ db *sql.DB }

func NewImportReportRepository(db *sql.DB) *SQLiteImportReportRepository {
	return &SQLiteImportReportRepository{db: db}
}

func (r *SQLiteImportReportRepository) Save(ctx context.Context, report *models.ImportReport) error {
	data, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("marshal import report: %w", err)
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO import_reports (report_json) VALUES (?)
	`, string(data))
	if err != nil {
		return fmt.Errorf("insert import report: %w", err)
	}
	return nil
}

func (r *SQLiteImportReportRepository) GetLast(ctx context.Context) (*models.ImportReport, error) {
	var raw string
	err := r.db.QueryRowContext(ctx, `
		SELECT report_json FROM import_reports ORDER BY id DESC LIMIT 1
	`).Scan(&raw)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no import reports found")
		}
		return nil, fmt.Errorf("get last import report: %w", err)
	}
	var report models.ImportReport
	if err := json.Unmarshal([]byte(raw), &report); err != nil {
		return nil, fmt.Errorf("unmarshal import report: %w", err)
	}
	return &report, nil
}
