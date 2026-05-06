package models

import "time"

type Sprint struct {
	ID        string
	Name      string
	StartDate *time.Time
	EndDate   *time.Time
	Team      string
	Capacity  *int
	Source    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
