package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"maestro/internal/models"
)

type AuditRepository interface {
	Create(context.Context, *models.DateAuditLog) error
	List(context.Context) ([]models.DateAuditLog, error)
}

type SQLiteAuditRepository struct{ db *sql.DB }

func NewAuditRepository(db *sql.DB) *SQLiteAuditRepository { return &SQLiteAuditRepository{db: db} }

func (r *SQLiteAuditRepository) Create(ctx context.Context, audit *models.DateAuditLog) error {
	if audit.ChangedAt.IsZero() {
		audit.ChangedAt = time.Now().UTC()
	}
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO date_audit_logs (
			entity_type, entity_id, changed_by, old_date, new_date, delta_days, reason, changed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, audit.EntityType, audit.EntityID, audit.ChangedBy, formatDatePtr(audit.OldDate), formatDatePtr(audit.NewDate), audit.DeltaDays, audit.Reason, formatTimestamp(audit.ChangedAt))
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("read audit log id: %w", err)
	}
	audit.ID = id
	return nil
}

func (r *SQLiteAuditRepository) List(ctx context.Context) ([]models.DateAuditLog, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, entity_type, entity_id, changed_by, old_date, new_date, delta_days, reason, changed_at
		FROM date_audit_logs ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	defer rows.Close()

	var audits []models.DateAuditLog
	for rows.Next() {
		audit, err := scanAudit(rows)
		if err != nil {
			return nil, err
		}
		audits = append(audits, *audit)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate audit logs: %w", err)
	}
	return audits, nil
}

func scanAudit(s scanner) (*models.DateAuditLog, error) {
	var audit models.DateAuditLog
	var oldDate, newDate, reason sql.NullString
	var changedAt string
	if err := s.Scan(&audit.ID, &audit.EntityType, &audit.EntityID, &audit.ChangedBy, &oldDate, &newDate, &audit.DeltaDays, &reason, &changedAt); err != nil {
		return nil, fmt.Errorf("scan audit log: %w", err)
	}
	if reason.Valid {
		audit.Reason = &reason.String
	}
	var err error
	if audit.OldDate, err = scanDatePtr(oldDate); err != nil {
		return nil, fmt.Errorf("scan audit old_date: %w", err)
	}
	if audit.NewDate, err = scanDatePtr(newDate); err != nil {
		return nil, fmt.Errorf("scan audit new_date: %w", err)
	}
	if audit.ChangedAt, err = scanTimestamp(changedAt); err != nil {
		return nil, fmt.Errorf("scan audit changed_at: %w", err)
	}
	return &audit, nil
}
