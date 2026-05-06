package repository

import (
	"database/sql"
	"time"
)

type scanner interface{ Scan(...any) error }

const (
	dateFormat      = "2006-01-02"
	timestampFormat = time.RFC3339
)

func formatDatePtr(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format(dateFormat)
}

func formatTimestamp(t time.Time) string {
	return t.UTC().Format(timestampFormat)
}

func scanDatePtr(raw sql.NullString) (*time.Time, error) {
	if !raw.Valid || raw.String == "" {
		return nil, nil
	}
	parsed, err := time.Parse(dateFormat, raw.String)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func scanTimestamp(raw string) (time.Time, error) {
	return time.Parse(timestampFormat, raw)
}
