package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"maestro/internal/models"
)

type SlipMetrics struct {
	EntityID        string  `json:"entity_id"`
	EntityType      string  `json:"entity_type"`
	AverageSlipDays float64 `json:"average_slip_days"`
	ItemCount       int     `json:"item_count"`
	TotalSlipDays   float64 `json:"total_slip_days"`
}

type MetricsRepository interface {
	StoryDeadlineHitRate(ctx context.Context) (float64, error)
	StoryScopeCreepRate(ctx context.Context) (float64, error)
	FeatureDeadlineHitRate(ctx context.Context) (float64, error)
	FeatureScopeCreepRate(ctx context.Context) (float64, error)
	EpicDeadlineHitRate(ctx context.Context) (float64, error)
	EpicScopeCreepRate(ctx context.Context) (float64, error)
	SlipByID(ctx context.Context, id string) (*SlipMetrics, error)
	OrphanedStories(ctx context.Context) ([]models.Story, error)
}

type SQLiteMetricsRepository struct{ db *sql.DB }

func NewMetricsRepository(db *sql.DB) *SQLiteMetricsRepository {
	return &SQLiteMetricsRepository{db: db}
}

func (r *SQLiteMetricsRepository) StoryDeadlineHitRate(ctx context.Context) (float64, error) {
	return r.deadlineHitRate(ctx, "stories")
}

func (r *SQLiteMetricsRepository) StoryScopeCreepRate(ctx context.Context) (float64, error) {
	return r.scopeCreepRate(ctx, "stories")
}

func (r *SQLiteMetricsRepository) FeatureDeadlineHitRate(ctx context.Context) (float64, error) {
	return r.deadlineHitRate(ctx, "features")
}

func (r *SQLiteMetricsRepository) FeatureScopeCreepRate(ctx context.Context) (float64, error) {
	return r.scopeCreepRate(ctx, "features")
}

func (r *SQLiteMetricsRepository) EpicDeadlineHitRate(ctx context.Context) (float64, error) {
	return r.deadlineHitRate(ctx, "epics")
}

func (r *SQLiteMetricsRepository) EpicScopeCreepRate(ctx context.Context) (float64, error) {
	return r.scopeCreepRate(ctx, "epics")
}

var allowedMetricTables = map[string]bool{
	"stories":  true,
	"features": true,
	"epics":    true,
}

func validateMetricTable(table string) error {
	if !allowedMetricTables[table] {
		return fmt.Errorf("invalid metric table: %s", table)
	}
	return nil
}

func (r *SQLiteMetricsRepository) deadlineHitRate(ctx context.Context, table string) (float64, error) {
	if err := validateMetricTable(table); err != nil {
		return 0, err
	}
	var count, hits sql.NullInt64
	query := fmt.Sprintf(`
		SELECT COUNT(*), COUNT(CASE WHEN actual_end_date IS NOT NULL AND committed_end_date IS NOT NULL AND actual_end_date <= committed_end_date THEN 1 END)
		FROM %s
		WHERE actual_end_date IS NOT NULL AND committed_end_date IS NOT NULL
	`, table)
	err := r.db.QueryRowContext(ctx, query).Scan(&count, &hits)
	if err != nil {
		return 0, fmt.Errorf("deadline hit rate for %s: %w", table, err)
	}
	if !count.Valid || count.Int64 == 0 {
		return 0, nil
	}
	return float64(hits.Int64) / float64(count.Int64), nil
}

func (r *SQLiteMetricsRepository) scopeCreepRate(ctx context.Context, table string) (float64, error) {
	if err := validateMetricTable(table); err != nil {
		return 0, err
	}
	var count, creeps sql.NullInt64
	query := fmt.Sprintf(`
		SELECT COUNT(*), COUNT(CASE WHEN original_end_date IS NOT NULL AND committed_end_date IS NOT NULL AND committed_end_date > original_end_date THEN 1 END)
		FROM %s
		WHERE original_end_date IS NOT NULL AND committed_end_date IS NOT NULL
	`, table)
	err := r.db.QueryRowContext(ctx, query).Scan(&count, &creeps)
	if err != nil {
		return 0, fmt.Errorf("scope creep rate for %s: %w", table, err)
	}
	if !count.Valid || count.Int64 == 0 {
		return 0, nil
	}
	return float64(creeps.Int64) / float64(count.Int64), nil
}

func (r *SQLiteMetricsRepository) SlipByID(ctx context.Context, id string) (*SlipMetrics, error) {
	// Determine entity type
	var entityType string
	var exists int
	if err := r.db.QueryRowContext(ctx, `SELECT 1 FROM epics WHERE id = ?`, id).Scan(&exists); err == nil {
		entityType = "epic"
	} else if err := r.db.QueryRowContext(ctx, `SELECT 1 FROM features WHERE id = ?`, id).Scan(&exists); err == nil {
		entityType = "feature"
	} else {
		return nil, fmt.Errorf("entity not found: %s", id)
	}

	var query string
	if entityType == "epic" {
		query = `
			SELECT COUNT(*), COALESCE(SUM(julianday(committed_end_date) - julianday(original_end_date)), 0)
			FROM features
			WHERE epic_id = ? AND date_source != 'inherited'
			  AND original_end_date IS NOT NULL AND committed_end_date IS NOT NULL
		`
	} else {
		query = `
			SELECT COUNT(*), COALESCE(SUM(julianday(committed_end_date) - julianday(original_end_date)), 0)
			FROM stories
			WHERE feature_id = ? AND date_source != 'inherited'
			  AND original_end_date IS NOT NULL AND committed_end_date IS NOT NULL
		`
	}

	var count sql.NullInt64
	var total sql.NullFloat64
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&count, &total); err != nil {
		return nil, fmt.Errorf("slip calculation: %w", err)
	}

	itemCount := int(count.Int64)
	var avg float64
	if itemCount > 0 {
		avg = total.Float64 / float64(itemCount)
	}

	return &SlipMetrics{
		EntityID:        id,
		EntityType:      entityType,
		AverageSlipDays: math.Round(avg*100) / 100,
		ItemCount:       itemCount,
		TotalSlipDays:   math.Round(total.Float64*100) / 100,
	}, nil
}

func (r *SQLiteMetricsRepository) OrphanedStories(ctx context.Context) ([]models.Story, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.feature_id, s.title, s.description, s.status, s.owner, s.sprint,
		       s.original_end_date, s.committed_end_date, s.actual_end_date,
		       s.story_points, s.date_source, s.created_at, s.updated_at
		FROM stories s
		LEFT JOIN features f ON s.feature_id = f.id
		WHERE s.feature_id = 'feature-unassigned' OR f.id IS NULL
		ORDER BY s.id
	`)
	if err != nil {
		return nil, fmt.Errorf("orphaned stories: %w", err)
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
		return nil, fmt.Errorf("iterate orphaned stories: %w", err)
	}
	return stories, nil
}


