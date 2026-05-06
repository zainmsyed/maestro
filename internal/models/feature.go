package models

import "time"

type Feature struct {
	ID               string
	EpicID           *string
	Title            string
	Description      string
	Status           string
	Owner            string
	Sprint           string
	OriginalEndDate  *time.Time
	CommittedEndDate *time.Time
	ActualEndDate    *time.Time
	StoryPoints      *int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
