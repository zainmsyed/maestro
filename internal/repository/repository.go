package repository

import "database/sql"

type Repositories struct {
	Epics         EpicRepository
	Features      FeatureRepository
	Stories       StoryRepository
	Sprints       SprintRepository
	Audits        AuditRepository
	ImportReports ImportReportRepository
	Metrics       MetricsRepository
}

func New(db *sql.DB) Repositories {
	return Repositories{
		Epics:         NewEpicRepository(db),
		Features:      NewFeatureRepository(db),
		Stories:       NewStoryRepository(db),
		Sprints:       NewSprintRepository(db),
		Audits:        NewAuditRepository(db),
		ImportReports: NewImportReportRepository(db),
		Metrics:       NewMetricsRepository(db),
	}
}
