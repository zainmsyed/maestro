package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"maestro/internal/models"
)

type EpicRepository interface {
	Create(context.Context, *models.Epic) error
	GetByID(context.Context, string) (*models.Epic, error)
	List(context.Context) ([]models.Epic, error)
}

type SQLiteEpicRepository struct{ db *sql.DB }

func NewEpicRepository(db *sql.DB) *SQLiteEpicRepository { return &SQLiteEpicRepository{db: db} }

func (r *SQLiteEpicRepository) Create(ctx context.Context, epic *models.Epic) error {
	now := time.Now().UTC()
	if epic.CreatedAt.IsZero() {
		epic.CreatedAt = now
	}
	if epic.UpdatedAt.IsZero() {
		epic.UpdatedAt = now
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO epics (
			id, title, description, status, owner, sprint_start, sprint_end,
			original_end_date, committed_end_date, actual_end_date,
			is_synthetic, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, epic.ID, epic.Title, epic.Description, epic.Status, epic.Owner, epic.SprintStart, epic.SprintEnd,
		formatDatePtr(epic.OriginalEndDate), formatDatePtr(epic.CommittedEndDate), formatDatePtr(epic.ActualEndDate),
		epic.IsSynthetic, formatTimestamp(epic.CreatedAt), formatTimestamp(epic.UpdatedAt))
	if err != nil {
		return fmt.Errorf("insert epic: %w", err)
	}
	return nil
}

func (r *SQLiteEpicRepository) GetByID(ctx context.Context, id string) (*models.Epic, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, title, description, status, owner, sprint_start, sprint_end,
		       original_end_date, committed_end_date, actual_end_date,
		       is_synthetic, created_at, updated_at
		FROM epics WHERE id = ?
	`, id)
	return scanEpic(row)
}

func (r *SQLiteEpicRepository) List(ctx context.Context) ([]models.Epic, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, description, status, owner, sprint_start, sprint_end,
		       original_end_date, committed_end_date, actual_end_date,
		       is_synthetic, created_at, updated_at
		FROM epics ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("list epics: %w", err)
	}
	defer rows.Close()

	var epics []models.Epic
	for rows.Next() {
		epic, err := scanEpic(rows)
		if err != nil {
			return nil, err
		}
		epics = append(epics, *epic)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate epics: %w", err)
	}
	return epics, nil
}

func scanEpic(s scanner) (*models.Epic, error) {
	var epic models.Epic
	var original, committed, actual sql.NullString
	var createdAt, updatedAt string
	if err := s.Scan(&epic.ID, &epic.Title, &epic.Description, &epic.Status, &epic.Owner, &epic.SprintStart, &epic.SprintEnd,
		&original, &committed, &actual, &epic.IsSynthetic, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("scan epic: %w", err)
	}
	var err error
	if epic.OriginalEndDate, err = scanDatePtr(original); err != nil {
		return nil, fmt.Errorf("scan epic original_end_date: %w", err)
	}
	if epic.CommittedEndDate, err = scanDatePtr(committed); err != nil {
		return nil, fmt.Errorf("scan epic committed_end_date: %w", err)
	}
	if epic.ActualEndDate, err = scanDatePtr(actual); err != nil {
		return nil, fmt.Errorf("scan epic actual_end_date: %w", err)
	}
	if epic.CreatedAt, err = scanTimestamp(createdAt); err != nil {
		return nil, fmt.Errorf("scan epic created_at: %w", err)
	}
	if epic.UpdatedAt, err = scanTimestamp(updatedAt); err != nil {
		return nil, fmt.Errorf("scan epic updated_at: %w", err)
	}
	return &epic, nil
}
