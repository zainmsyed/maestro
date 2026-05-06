package repository

import "database/sql"

type Repositories struct {
	Epics    EpicRepository
	Features FeatureRepository
	Sprints  SprintRepository
	Audits   AuditRepository
}

func New(db *sql.DB) Repositories {
	return Repositories{
		Epics:    NewEpicRepository(db),
		Features: NewFeatureRepository(db),
		Sprints:  NewSprintRepository(db),
		Audits:   NewAuditRepository(db),
	}
}
