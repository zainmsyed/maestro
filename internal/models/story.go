package models

import "time"

type Story struct {
	ID               string
	FeatureID        string
	Title            string
	Description      string
	Status           string
	Owner            string
	Sprint           string
	StoryPoints      *int
	OriginalEndDate  *time.Time
	CommittedEndDate *time.Time
	ActualEndDate    *time.Time
	DateSource       string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
