package importer

import (
	"fmt"
	"strings"
	"time"

	"maestro/internal/repository"
)

const (
	syntheticUnassignedEpicID    = "epic-unassigned"
	syntheticUnassignedFeatureID = "feature-unassigned"
)

type Importer struct {
	repos repository.Repositories
}

type ImportReport struct {
	EpicCount                int
	FeatureCount             int
	StoryCount               int
	SprintsDetected          []string
	MissingDatesCount        int
	MissingSprintCount       int
	OrphanedFeatures         int
	OrphanedStories          int
	SkippedRows              int
	DetectedDateFormat       string
	DateAssignmentCandidates []DateAssignmentCandidate
	AmbiguousDates           []AmbiguousDateCandidate
	Warnings                 []string
	SyntheticStoryIDs        []string
}

type DateAssignmentCandidate struct {
	RowNumber     int
	WorkItemType  string
	ID            string
	Title         string
	AssignedOwner string
}

type AmbiguousDateCandidate struct {
	RowNumber    int
	WorkItemType string
	ID           string
	Title        string
	RawDate      string
	ParsedDate   time.Time
}

type assignedTo struct {
	DisplayName string
	Email       *string
}

type parsedDate struct {
	Time      *time.Time
	Format    string
	Ambiguous bool
}

type rawEntity struct {
	id          string
	parentID    string
	itemType    string
	title       string
	state       string
	owner       string
	sprint      string
	storyPoints *int
	targetDate  *time.Time
	rowNumber   int
}

type importState struct {
	report          *ImportReport
	createdEpics    map[string]bool
	createdFeatures map[string]bool
	createdStories  map[string]bool
	sprints         map[string]struct{}
	dateFormats     map[string]int
	syntheticIDSeq  int
}

func New(repos repository.Repositories) *Importer {
	return &Importer{repos: repos}
}

func newImportState() *importState {
	return &importState{
		report:          &ImportReport{},
		createdEpics:    map[string]bool{},
		createdFeatures: map[string]bool{},
		createdStories:  map[string]bool{},
		sprints:         map[string]struct{}{},
		dateFormats:     map[string]int{},
	}
}

func (s *importState) addSprint(name string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	s.sprints[name] = struct{}{}
}

func (s *importState) finalizeReport() *ImportReport {
	s.report.SprintsDetected = sortedKeys(s.sprints)
	s.report.DetectedDateFormat = summarizeDateFormats(s.dateFormats)
	return s.report
}

func (s *importState) nextSyntheticID(itemType string) string {
	s.syntheticIDSeq++
	return fmt.Sprintf("%s-auto-%d", itemType, s.syntheticIDSeq)
}

func cloneStringPtr(value string) *string {
	copied := value
	return &copied
}

func withDateLocked(date *time.Time) (*time.Time, *time.Time) {
	if date == nil {
		return nil, nil
	}
	original := *date
	committed := *date
	return &original, &committed
}
