package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"maestro/internal/models"
)

type StoryRepository interface {
	Create(context.Context, *models.Story) error
	GetByID(context.Context, string) (*models.Story, error)
	List(context.Context) ([]models.Story, error)
	ListByFeatureID(context.Context, string) ([]models.Story, error)
	UpdateDate(context.Context, string, *time.Time, *time.Time, string) error
	UpdateFeatureID(context.Context, string, string) error
}

type SQLiteStoryRepository struct{ db *sql.DB }

func NewStoryRepository(db *sql.DB) *SQLiteStoryRepository {
	return &SQLiteStoryRepository{db: db}
}

func (r *SQLiteStoryRepository) Create(ctx context.Context, story *models.Story) error {
	now := time.Now().UTC()
	if story.CreatedAt.IsZero() {
		story.CreatedAt = now
	}
	if story.UpdatedAt.IsZero() {
		story.UpdatedAt = now
	}
	if story.DateSource == "" {
		story.DateSource = "imported"
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO stories (
			id, feature_id, title, description, status, owner, sprint,
			original_end_date, committed_end_date, actual_end_date,
			story_points, date_source, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, story.ID, story.FeatureID, story.Title, story.Description, story.Status, story.Owner, story.Sprint,
		formatDatePtr(story.OriginalEndDate), formatDatePtr(story.CommittedEndDate), formatDatePtr(story.ActualEndDate),
		story.StoryPoints, story.DateSource, formatTimestamp(story.CreatedAt), formatTimestamp(story.UpdatedAt))
	if err != nil {
		return fmt.Errorf("insert story: %w", err)
	}
	return nil
}

func (r *SQLiteStoryRepository) GetByID(ctx context.Context, id string) (*models.Story, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, feature_id, title, description, status, owner, sprint,
		       original_end_date, committed_end_date, actual_end_date,
		       story_points, date_source, created_at, updated_at
		FROM stories WHERE id = ?
	`, id)
	return scanStory(row)
}

func (r *SQLiteStoryRepository) List(ctx context.Context) ([]models.Story, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, feature_id, title, description, status, owner, sprint,
		       original_end_date, committed_end_date, actual_end_date,
		       story_points, date_source, created_at, updated_at
		FROM stories ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("list stories: %w", err)
	}
	defer rows.Close()

	var stories []models.Story
	for rows.Next() {
		story, err := scanStory(rows)
		if err != nil {
			return nil, err
		}
		stories = append(stories, *story)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stories: %w", err)
	}
	return stories, nil
}

func (r *SQLiteStoryRepository) ListByFeatureID(ctx context.Context, featureID string) ([]models.Story, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, feature_id, title, description, status, owner, sprint,
		       original_end_date, committed_end_date, actual_end_date,
		       story_points, date_source, created_at, updated_at
		FROM stories WHERE feature_id = ? ORDER BY id
	`, featureID)
	if err != nil {
		return nil, fmt.Errorf("list stories by feature: %w", err)
	}
	defer rows.Close()

	var stories []models.Story
	for rows.Next() {
		story, err := scanStory(rows)
		if err != nil {
			return nil, err
		}
		stories = append(stories, *story)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stories by feature: %w", err)
	}
	return stories, nil
}

func (r *SQLiteStoryRepository) UpdateDate(ctx context.Context, id string, original, committed *time.Time, dateSource string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE stories SET
			original_end_date = ?,
			committed_end_date = ?,
			date_source = ?,
			updated_at = ?
		WHERE id = ?
	`, formatDatePtr(original), formatDatePtr(committed), dateSource, formatTimestamp(time.Now().UTC()), id)
	if err != nil {
		return fmt.Errorf("update story date: %w", err)
	}
	return nil
}

func (r *SQLiteStoryRepository) UpdateFeatureID(ctx context.Context, id, featureID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE stories SET feature_id = ?, updated_at = ? WHERE id = ?
	`, featureID, formatTimestamp(time.Now().UTC()), id)
	if err != nil {
		return fmt.Errorf("update story feature: %w", err)
	}
	return nil
}

func scanStory(s scanner) (*models.Story, error) {
	var story models.Story
	var original, committed, actual sql.NullString
	var storyPoints sql.NullInt64
	var createdAt, updatedAt string
	if err := s.Scan(&story.ID, &story.FeatureID, &story.Title, &story.Description, &story.Status, &story.Owner, &story.Sprint,
		&original, &committed, &actual, &storyPoints, &story.DateSource, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("scan story: %w", err)
	}
	if storyPoints.Valid {
		value := int(storyPoints.Int64)
		story.StoryPoints = &value
	}
	var err error
	if story.OriginalEndDate, err = scanDatePtr(original); err != nil {
		return nil, fmt.Errorf("scan story original_end_date: %w", err)
	}
	if story.CommittedEndDate, err = scanDatePtr(committed); err != nil {
		return nil, fmt.Errorf("scan story committed_end_date: %w", err)
	}
	if story.ActualEndDate, err = scanDatePtr(actual); err != nil {
		return nil, fmt.Errorf("scan story actual_end_date: %w", err)
	}
	if story.CreatedAt, err = scanTimestamp(createdAt); err != nil {
		return nil, fmt.Errorf("scan story created_at: %w", err)
	}
	if story.UpdatedAt, err = scanTimestamp(updatedAt); err != nil {
		return nil, fmt.Errorf("scan story updated_at: %w", err)
	}
	return &story, nil
}
