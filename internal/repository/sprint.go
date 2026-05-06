package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"maestro/internal/models"
)

type SprintRepository interface {
	Create(context.Context, *models.Sprint) error
	GetByID(context.Context, string) (*models.Sprint, error)
	List(context.Context) ([]models.Sprint, error)
}

type SQLiteSprintRepository struct{ db *sql.DB }

func NewSprintRepository(db *sql.DB) *SQLiteSprintRepository { return &SQLiteSprintRepository{db: db} }

func (r *SQLiteSprintRepository) Create(ctx context.Context, sprint *models.Sprint) error {
	now := time.Now().UTC()
	if sprint.CreatedAt.IsZero() {
		sprint.CreatedAt = now
	}
	if sprint.UpdatedAt.IsZero() {
		sprint.UpdatedAt = now
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO sprints (
			id, name, start_date, end_date, team, capacity, source, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, sprint.ID, sprint.Name, formatDatePtr(sprint.StartDate), formatDatePtr(sprint.EndDate), sprint.Team,
		sprint.Capacity, sprint.Source, formatTimestamp(sprint.CreatedAt), formatTimestamp(sprint.UpdatedAt))
	if err != nil {
		return fmt.Errorf("insert sprint: %w", err)
	}
	return nil
}

func (r *SQLiteSprintRepository) GetByID(ctx context.Context, id string) (*models.Sprint, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, start_date, end_date, team, capacity, source, created_at, updated_at
		FROM sprints WHERE id = ?
	`, id)
	return scanSprint(row)
}

func (r *SQLiteSprintRepository) List(ctx context.Context) ([]models.Sprint, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, start_date, end_date, team, capacity, source, created_at, updated_at
		FROM sprints ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("list sprints: %w", err)
	}
	defer rows.Close()

	var sprints []models.Sprint
	for rows.Next() {
		sprint, err := scanSprint(rows)
		if err != nil {
			return nil, err
		}
		sprints = append(sprints, *sprint)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sprints: %w", err)
	}
	return sprints, nil
}

func scanSprint(s scanner) (*models.Sprint, error) {
	var sprint models.Sprint
	var startDate, endDate sql.NullString
	var capacity sql.NullInt64
	var createdAt, updatedAt string
	if err := s.Scan(&sprint.ID, &sprint.Name, &startDate, &endDate, &sprint.Team, &capacity, &sprint.Source, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("scan sprint: %w", err)
	}
	if capacity.Valid {
		value := int(capacity.Int64)
		sprint.Capacity = &value
	}
	var err error
	if sprint.StartDate, err = scanDatePtr(startDate); err != nil {
		return nil, fmt.Errorf("scan sprint start_date: %w", err)
	}
	if sprint.EndDate, err = scanDatePtr(endDate); err != nil {
		return nil, fmt.Errorf("scan sprint end_date: %w", err)
	}
	if sprint.CreatedAt, err = scanTimestamp(createdAt); err != nil {
		return nil, fmt.Errorf("scan sprint created_at: %w", err)
	}
	if sprint.UpdatedAt, err = scanTimestamp(updatedAt); err != nil {
		return nil, fmt.Errorf("scan sprint updated_at: %w", err)
	}
	return &sprint, nil
}
