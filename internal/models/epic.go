package models

import "time"

type Epic struct {
	ID               string
	Title            string
	Description      string
	Status           string
	Owner            string
	SprintStart      string
	SprintEnd        string
	OriginalEndDate  *time.Time
	CommittedEndDate *time.Time
	ActualEndDate    *time.Time
	IsSynthetic      bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
