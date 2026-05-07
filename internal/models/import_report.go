package models

import "time"

type ImportReport struct {
	EpicCount                int                       `json:"epic_count"`
	FeatureCount             int                       `json:"feature_count"`
	StoryCount               int                       `json:"story_count"`
	ExistingSkipped          int                       `json:"existing_skipped"`
	SprintsDetected          []string                  `json:"sprints_detected"`
	MissingDatesCount        int                       `json:"missing_dates_count"`
	MissingSprintCount       int                       `json:"missing_sprint_count"`
	OrphanedFeatures         int                       `json:"orphaned_features"`
	OrphanedStories          int                       `json:"orphaned_stories"`
	SkippedRows              int                       `json:"skipped_rows"`
	DetectedDateFormat       string                    `json:"detected_date_format"`
	DateAssignmentCandidates []DateAssignmentCandidate `json:"date_assignment_candidates"`
	AmbiguousDates           []AmbiguousDateCandidate  `json:"ambiguous_dates"`
	Warnings                 []string                  `json:"warnings"`
	SyntheticStoryIDs        []string                  `json:"synthetic_story_ids"`
}

type DateAssignmentCandidate struct {
	RowNumber     int    `json:"row_number"`
	WorkItemType  string `json:"work_item_type"`
	ID            string `json:"id"`
	Title         string `json:"title"`
	AssignedOwner string `json:"assigned_owner"`
}

type AmbiguousDateCandidate struct {
	RowNumber    int       `json:"row_number"`
	WorkItemType string    `json:"work_item_type"`
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	RawDate      string    `json:"raw_date"`
	ParsedDate   time.Time `json:"parsed_date"`
}
