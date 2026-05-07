package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"maestro/internal/models"
)

type FeatureRepository interface {
	Create(context.Context, *models.Feature) error
	GetByID(context.Context, string) (*models.Feature, error)
	List(context.Context) ([]models.Feature, error)
	ListByEpicID(context.Context, string) ([]models.Feature, error)
	UpdateDate(context.Context, string, *time.Time, *time.Time, string) error
	UpdateEpicID(context.Context, string, *string) error
}

type SQLiteFeatureRepository struct{ db *sql.DB }

func NewFeatureRepository(db *sql.DB) *SQLiteFeatureRepository {
	return &SQLiteFeatureRepository{db: db}
}

func (r *SQLiteFeatureRepository) Create(ctx context.Context, feature *models.Feature) error {
	now := time.Now().UTC()
	if feature.CreatedAt.IsZero() {
		feature.CreatedAt = now
	}
	if feature.UpdatedAt.IsZero() {
		feature.UpdatedAt = now
	}
	if feature.DateSource == "" {
		feature.DateSource = "imported"
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO features (
			id, epic_id, title, description, status, owner, sprint,
			original_end_date, committed_end_date, actual_end_date,
			story_points, date_source, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, feature.ID, feature.EpicID, feature.Title, feature.Description, feature.Status, feature.Owner, feature.Sprint,
		formatDatePtr(feature.OriginalEndDate), formatDatePtr(feature.CommittedEndDate), formatDatePtr(feature.ActualEndDate),
		feature.StoryPoints, feature.DateSource, formatTimestamp(feature.CreatedAt), formatTimestamp(feature.UpdatedAt))
	if err != nil {
		return fmt.Errorf("insert feature: %w", err)
	}
	return nil
}

func (r *SQLiteFeatureRepository) GetByID(ctx context.Context, id string) (*models.Feature, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, epic_id, title, description, status, owner, sprint,
		       original_end_date, committed_end_date, actual_end_date,
		       story_points, date_source, created_at, updated_at
		FROM features WHERE id = ?
	`, id)
	return scanFeature(row)
}

func (r *SQLiteFeatureRepository) List(ctx context.Context) ([]models.Feature, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, epic_id, title, description, status, owner, sprint,
		       original_end_date, committed_end_date, actual_end_date,
		       story_points, date_source, created_at, updated_at
		FROM features ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("list features: %w", err)
	}
	defer rows.Close()

	var features []models.Feature
	for rows.Next() {
		feature, err := scanFeature(rows)
		if err != nil {
			return nil, err
		}
		features = append(features, *feature)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate features: %w", err)
	}
	return features, nil
}

func (r *SQLiteFeatureRepository) ListByEpicID(ctx context.Context, epicID string) ([]models.Feature, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, epic_id, title, description, status, owner, sprint,
		       original_end_date, committed_end_date, actual_end_date,
		       story_points, date_source, created_at, updated_at
		FROM features WHERE epic_id = ? ORDER BY id
	`, epicID)
	if err != nil {
		return nil, fmt.Errorf("list features by epic: %w", err)
	}
	defer rows.Close()

	var features []models.Feature
	for rows.Next() {
		feature, err := scanFeature(rows)
		if err != nil {
			return nil, err
		}
		features = append(features, *feature)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate features by epic: %w", err)
	}
	return features, nil
}

func (r *SQLiteFeatureRepository) UpdateDate(ctx context.Context, id string, original, committed *time.Time, dateSource string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE features SET
			original_end_date = ?,
			committed_end_date = ?,
			date_source = ?,
			updated_at = ?
		WHERE id = ?
	`, formatDatePtr(original), formatDatePtr(committed), dateSource, formatTimestamp(time.Now().UTC()), id)
	if err != nil {
		return fmt.Errorf("update feature date: %w", err)
	}
	return nil
}

func (r *SQLiteFeatureRepository) UpdateEpicID(ctx context.Context, id string, epicID *string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE features SET epic_id = ?, updated_at = ? WHERE id = ?
	`, epicID, formatTimestamp(time.Now().UTC()), id)
	if err != nil {
		return fmt.Errorf("update feature epic: %w", err)
	}
	return nil
}

func scanFeature(s scanner) (*models.Feature, error) {
	var feature models.Feature
	var epicID sql.NullString
	var original, committed, actual sql.NullString
	var storyPoints sql.NullInt64
	var createdAt, updatedAt string
	if err := s.Scan(&feature.ID, &epicID, &feature.Title, &feature.Description, &feature.Status, &feature.Owner, &feature.Sprint,
		&original, &committed, &actual, &storyPoints, &feature.DateSource, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("scan feature: %w", err)
	}
	if epicID.Valid {
		feature.EpicID = &epicID.String
	}
	if storyPoints.Valid {
		value := int(storyPoints.Int64)
		feature.StoryPoints = &value
	}
	var err error
	if feature.OriginalEndDate, err = scanDatePtr(original); err != nil {
		return nil, fmt.Errorf("scan feature original_end_date: %w", err)
	}
	if feature.CommittedEndDate, err = scanDatePtr(committed); err != nil {
		return nil, fmt.Errorf("scan feature committed_end_date: %w", err)
	}
	if feature.ActualEndDate, err = scanDatePtr(actual); err != nil {
		return nil, fmt.Errorf("scan feature actual_end_date: %w", err)
	}
	if feature.CreatedAt, err = scanTimestamp(createdAt); err != nil {
		return nil, fmt.Errorf("scan feature created_at: %w", err)
	}
	if feature.UpdatedAt, err = scanTimestamp(updatedAt); err != nil {
		return nil, fmt.Errorf("scan feature updated_at: %w", err)
	}
	return &feature, nil
}
