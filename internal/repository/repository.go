package repository

import "database/sql"

type Repositories struct {
	Epics    EpicRepository
	Features FeatureRepository
	Stories  StoryRepository
	Sprints  SprintRepository
	Audits   AuditRepository
}

func New(db *sql.DB) Repositories {
	return Repositories{
		Epics:    NewEpicRepository(db),
		Features: NewFeatureRepository(db),
		Stories:  NewStoryRepository(db),
		Sprints:  NewSprintRepository(db),
		Audits:   NewAuditRepository(db),
	}
}
