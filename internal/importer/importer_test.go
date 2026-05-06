package importer

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	maestrodb "maestro/internal/db"
	"maestro/internal/repository"
)

func TestNormalizeHeaders(t *testing.T) {
	t.Parallel()
	headers := []string{"Parent", "ID", "Work Item Type", "Title 1", "Title 2", "Title 3", "State", "Assigned To", "Iteration Path", "Story Points", "Target Date", "Area Path"}
	indices, titleCols, err := NormalizeHeaders(headers)
	if err != nil {
		t.Fatalf("normalize headers: %v", err)
	}
	want := map[string]int{
		"parent":         0,
		"id":             1,
		"work_item_type": 2,
	}
	for key, index := range want {
		if indices[key] != index {
			t.Fatalf("header %s index = %d, want %d", key, indices[key], index)
		}
	}
	if len(titleCols) != 3 {
		t.Fatalf("expected 3 title columns, got %d", len(titleCols))
	}
}

func TestNormalizeHeadersMissingRequired(t *testing.T) {
	t.Parallel()
	_, _, err := NormalizeHeaders([]string{"ID", "Work Item Type"})
	if err == nil {
		t.Fatalf("expected error for missing Parent column")
	}
}

func TestNormalizeWorkItemType(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"Epic":                 "epic",
		"Feature":              "feature",
		"User Story":           "story",
		"Product Backlog Item": "story",
		"Requirement":          "story",
	}
	for input, want := range cases {
		got, ok := NormalizeWorkItemType(input)
		if !ok || got != want {
			t.Fatalf("NormalizeWorkItemType(%q) = (%q, %v), want (%q, true)", input, got, ok, want)
		}
	}
	if got, ok := NormalizeWorkItemType("Bug"); ok || got != "" {
		t.Fatalf("expected Bug to be unsupported, got (%q, %v)", got, ok)
	}
}

func TestNormalizeID(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input string
		want  string
	}{
		{"500588.0", "500588"},
		{"500000", "500000"},
		{"  500001  ", "500001"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := normalizeID(tc.input); got != tc.want {
			t.Fatalf("normalizeID(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestExtractTitle(t *testing.T) {
	t.Parallel()
	headers := []string{"Title 1", "Title 2", "Title 3", "Title 4", "Title 5"}
	indices, titleCols, err := NormalizeHeaders(append([]string{"Parent", "ID", "Work Item Type"}, headers...))
	if err != nil {
		t.Fatalf("normalize headers: %v", err)
	}

	cases := []struct {
		name   string
		record []string
		want   string
	}{
		{
			name:   "title in title3",
			record: []string{"", "", "", "Epic", "Feature", "Story", "", "", "", "", "", "", "", ""},
			want:   "Story",
		},
		{
			name:   "title in title5 only",
			record: []string{"", "", "", "", "", "", "", "Deep", "", "", "", "", "", ""},
			want:   "Deep",
		},
		{
			name:   "no titles",
			record: []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""},
			want:   "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := extractTitle(tc.record, indices, titleCols); got != tc.want {
				t.Fatalf("extractTitle() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseTargetDate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input     string
		wantFmt   string
		wantDate  time.Time
		ambiguous bool
	}{
		{input: "2026-05-20", wantFmt: "ISO 8601", wantDate: date(2026, 5, 20)},
		{input: "2026/05/10", wantFmt: "ISO slashes", wantDate: date(2026, 5, 10)},
		{input: "05/11/2026", wantFmt: "MM/DD/YYYY", wantDate: date(2026, 5, 11), ambiguous: true},
		{input: "14/05/2026", wantFmt: "DD/MM/YYYY", wantDate: date(2026, 5, 14)},
		{input: "May 22, 2026", wantFmt: "long form", wantDate: date(2026, 5, 22)},
		{input: "2026-05-20T14:30:00Z", wantFmt: "ISO with timezone", wantDate: date(2026, 5, 20)},
		{input: "02/13/2026 12:00:00 AM", wantFmt: "US datetime", wantDate: date(2026, 2, 13)},
		{input: "02/13/2026 15:04:05", wantFmt: "US datetime 24h", wantDate: date(2026, 2, 13)},
		{input: "2-Jan-2026", wantFmt: "abbreviated", wantDate: date(2026, 1, 2)},
	}
	for _, tc := range cases {
		got, err := ParseTargetDate(tc.input)
		if err != nil {
			t.Fatalf("ParseTargetDate(%q): %v", tc.input, err)
		}
		if got.Time == nil || !got.Time.Equal(tc.wantDate) {
			t.Fatalf("ParseTargetDate(%q) date = %#v, want %v", tc.input, got.Time, tc.wantDate)
		}
		if got.Format != tc.wantFmt {
			t.Fatalf("ParseTargetDate(%q) format = %q, want %q", tc.input, got.Format, tc.wantFmt)
		}
		if got.Ambiguous != tc.ambiguous {
			t.Fatalf("ParseTargetDate(%q) ambiguous = %v, want %v", tc.input, got.Ambiguous, tc.ambiguous)
		}
	}
}

func TestParseIterationPath(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input         string
		wantSprint    string
		wantScheduled bool
	}{
		{input: "FinDash\\\\Archive\\\\FY26 Q3", wantSprint: "", wantScheduled: false},
		{input: "FinDash\\\\Backlog\\\\Queue", wantSprint: "", wantScheduled: false},
		{input: "Program\\\\Sprint 10", wantSprint: "Sprint 10", wantScheduled: true},
		{input: "Program\\\\Delivery\\\\FY26 Q3\\\\FY26 Q3.1", wantSprint: "FY26 Q3", wantScheduled: true},
		{input: "", wantSprint: "", wantScheduled: false},
	}
	for _, tc := range cases {
		sprint, scheduled := ParseIterationPath(tc.input)
		if sprint != tc.wantSprint || scheduled != tc.wantScheduled {
			t.Fatalf("ParseIterationPath(%q) = (%q, %v), want (%q, %v)", tc.input, sprint, scheduled, tc.wantSprint, tc.wantScheduled)
		}
	}
}

func TestParseAssignedTo(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		raw         string
		wantName    string
		wantEmail   *string
	}{
		{name: "angle bracket", raw: "Kline, Zoe \u003czoe.kline70@example.com\u003e", wantName: "Kline, Zoe", wantEmail: stringPtr("zoe.kline70@example.com")},
		{name: "parenthesis", raw: "Casey Coach (casey@example.com)", wantName: "Casey Coach", wantEmail: stringPtr("casey@example.com")},
		{name: "plain email", raw: "robin@example.com", wantName: "robin@example.com", wantEmail: stringPtr("robin@example.com")},
		{name: "empty", raw: "", wantName: "", wantEmail: nil},
		{name: "no email", raw: "Alex Admin", wantName: "Alex Admin", wantEmail: nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ParseAssignedTo(tc.raw)
			if got.DisplayName != tc.wantName {
				t.Fatalf("DisplayName = %q, want %q", got.DisplayName, tc.wantName)
			}
			if (got.Email == nil) != (tc.wantEmail == nil) {
				t.Fatalf("Email nil mismatch")
			}
			if got.Email != nil && tc.wantEmail != nil && *got.Email != *tc.wantEmail {
				t.Fatalf("Email = %q, want %q", *got.Email, *tc.wantEmail)
			}
		})
	}
}

func TestParseStoryPoints(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input     string
		wantValue *int
		wantWarn  string
	}{
		{input: "", wantValue: nil, wantWarn: ""},
		{input: "5", wantValue: intPtr(5), wantWarn: ""},
		{input: "5.0", wantValue: intPtr(5), wantWarn: ""},
		{input: "5.7", wantValue: intPtr(6), wantWarn: ""},
		{input: "tbd", wantValue: nil, wantWarn: "non-numeric story points \"tbd\""},
	}
	for _, tc := range cases {
		got, warn := ParseStoryPoints(tc.input)
		if (got == nil) != (tc.wantValue == nil) || (got != nil && tc.wantValue != nil && *got != *tc.wantValue) {
			var gotVal, wantVal int
			if got != nil {
				gotVal = *got
			}
			if tc.wantValue != nil {
				wantVal = *tc.wantValue
			}
			t.Fatalf("ParseStoryPoints(%q) value = %d, want %d", tc.input, gotVal, wantVal)
		}
		if warn != tc.wantWarn {
			t.Fatalf("ParseStoryPoints(%q) warning = %q, want %q", tc.input, warn, tc.wantWarn)
		}
	}
}

func TestImportCSV(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)
	fixturePath := filepath.Join("..", "..", "testdata", "tree-export.csv")
	file, err := os.Open(fixturePath)
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer file.Close()

	imp := New(repos)
	report, err := imp.ImportCSV(ctx, file)
	if err != nil {
		t.Fatalf("import csv: %v", err)
	}

	epics, err := repos.Epics.List(ctx)
	if err != nil {
		t.Fatalf("list epics: %v", err)
	}
	features, err := repos.Features.List(ctx)
	if err != nil {
		t.Fatalf("list features: %v", err)
	}
	stories, err := repos.Stories.List(ctx)
	if err != nil {
		t.Fatalf("list stories: %v", err)
	}

	if len(epics) != 2 || report.EpicCount != 2 {
		t.Fatalf("expected 2 epics, got list=%d report=%d", len(epics), report.EpicCount)
	}
	if len(features) != 3 || report.FeatureCount != 3 {
		t.Fatalf("expected 3 features, got list=%d report=%d", len(features), report.FeatureCount)
	}
	if len(stories) != 5 || report.StoryCount != 5 {
		t.Fatalf("expected 5 stories, got list=%d report=%d", len(stories), report.StoryCount)
	}
	if len(report.SyntheticStoryIDs) != 2 {
		t.Fatalf("expected 2 synthetic story ids, got %d", len(report.SyntheticStoryIDs))
	}
	if report.MissingDatesCount != 1 || len(report.DateAssignmentCandidates) != 1 {
		t.Fatalf("expected 1 missing date candidate, got missing=%d candidates=%d", report.MissingDatesCount, len(report.DateAssignmentCandidates))
	}
	if report.OrphanedFeatures != 1 {
		t.Fatalf("expected 1 orphaned feature, got %d", report.OrphanedFeatures)
	}
	if report.OrphanedStories != 1 {
		t.Fatalf("expected 1 orphaned story, got %d", report.OrphanedStories)
	}
	if report.MissingSprintCount != 1 {
		t.Fatalf("expected 1 missing sprint, got %d", report.MissingSprintCount)
	}
	if report.DetectedDateFormat != "mixed" {
		t.Fatalf("expected mixed date format, got %q", report.DetectedDateFormat)
	}
	if len(report.AmbiguousDates) != 1 {
		t.Fatalf("expected 1 ambiguous date, got %d", len(report.AmbiguousDates))
	}
	if len(report.Warnings) == 0 {
		t.Fatalf("expected warnings for non-numeric story points")
	}

	var syntheticEpicFound, syntheticFeatureFound bool
	for _, epic := range epics {
		if epic.ID == syntheticUnassignedEpicID {
			syntheticEpicFound = epic.IsSynthetic
		}
	}
	for _, feature := range features {
		if feature.ID == syntheticUnassignedFeatureID {
			syntheticFeatureFound = true
			if feature.EpicID == nil || *feature.EpicID != syntheticUnassignedEpicID {
				t.Fatalf("synthetic feature not attached to synthetic epic: %#v", feature.EpicID)
			}
		}
		if feature.ID == "F-ORPH" {
			if feature.EpicID == nil || *feature.EpicID != syntheticUnassignedEpicID {
				t.Fatalf("orphan feature not attached to synthetic epic: %#v", feature.EpicID)
			}
		}
	}
	if !syntheticEpicFound {
		t.Fatalf("expected synthetic unassigned epic to exist")
	}
	if !syntheticFeatureFound {
		t.Fatalf("expected synthetic unassigned feature to exist")
	}

	orphanStoryFound := false
	deepTitleFound := false
	for _, story := range stories {
		if story.Title == "Orphan login story" {
			orphanStoryFound = true
			if story.FeatureID != syntheticUnassignedFeatureID {
				t.Fatalf("orphan story attached to %q, want %q", story.FeatureID, syntheticUnassignedFeatureID)
			}
			if story.OriginalEndDate != nil || story.CommittedEndDate != nil {
				t.Fatalf("orphan story dates should be nil when target date missing")
			}
		}
		if story.Title == "Deep nested title" {
			deepTitleFound = true
		}
		if story.DateSource != "imported" {
			t.Fatalf("story %s date_source = %q, want imported", story.ID, story.DateSource)
		}
	}
	if !orphanStoryFound {
		t.Fatalf("expected orphan story to be imported")
	}
	if !deepTitleFound {
		t.Fatalf("expected deep nested title story to be imported")
	}
}

func TestImportCSV_SyntheticIDs(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	csv := `Parent,ID,Work Item Type,Title 1,Title 2,Title 3,State,Assigned To,Iteration Path,Story Points,Target Date,Area Path
,,Epic,Synthetic Epic,,,,New,Alice,Program\Sprint 1,,2026-05-01,Area
E-1,,Feature,Synthetic Epic,Synthetic Feature,,,New,Bob,Program\Sprint 1,,2026-05-02,Area
F-1,,User Story,Synthetic Epic,Synthetic Feature,Synthetic Story,,New,Charlie,Program\Sprint 1,,2026-05-03,Area
`

	imp := New(repos)
	report, err := imp.ImportCSV(ctx, strings.NewReader(csv))
	if err != nil {
		t.Fatalf("import csv: %v", err)
	}

	if report.EpicCount != 2 {
		t.Fatalf("expected 2 epics (1 real + 1 synthetic empty-id), got %d", report.EpicCount)
	}
	if report.FeatureCount != 2 {
		t.Fatalf("expected 2 features (1 real + 1 synthetic empty-id), got %d", report.FeatureCount)
	}
	if report.StoryCount != 1 {
		t.Fatalf("expected 1 story (synthetic empty-id), got %d", report.StoryCount)
	}
	if len(report.SyntheticStoryIDs) != 1 {
		t.Fatalf("expected 1 synthetic story ID, got %d", len(report.SyntheticStoryIDs))
	}

	epics, _ := repos.Epics.List(ctx)
	features, _ := repos.Features.List(ctx)
	stories, _ := repos.Stories.List(ctx)

	var emptyEpicFound, emptyFeatureFound bool
	for _, e := range epics {
		if strings.HasPrefix(e.ID, "epic-auto-") {
			emptyEpicFound = true
		}
	}
	for _, f := range features {
		if strings.HasPrefix(f.ID, "feature-auto-") {
			emptyFeatureFound = true
		}
	}
	if !emptyEpicFound {
		t.Fatalf("expected empty epic to receive synthetic epic-auto-* ID")
	}
	if !emptyFeatureFound {
		t.Fatalf("expected empty feature to receive synthetic feature-auto-* ID")
	}
	if len(stories) != 1 || !strings.HasPrefix(stories[0].ID, "story-auto-") {
		t.Fatalf("expected empty story to receive synthetic story-auto-* ID, got %v", stories)
	}
}

func TestImportCSV_WrongTypeParent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	csv := `Parent,ID,Work Item Type,Title 1,Title 2,Title 3,State,Assigned To,Iteration Path,Story Points,Target Date,Area Path
,E-1,Epic,Real Epic,,,Active,Alice,Program\Sprint 1,,2026-05-01,Area
E-1,S-1,User Story,,,,Orphaned by wrong parent,New,Bob,Program\Sprint 1,,2026-05-02,Area
`

	imp := New(repos)
	report, err := imp.ImportCSV(ctx, strings.NewReader(csv))
	if err != nil {
		t.Fatalf("import csv: %v", err)
	}

	if report.EpicCount != 2 {
		t.Fatalf("expected 2 epics (real + synthetic unassigned for orphaned story), got %d", report.EpicCount)
	}
	if report.StoryCount != 1 {
		t.Fatalf("expected 1 story, got %d", report.StoryCount)
	}
	if report.OrphanedStories != 1 {
		t.Fatalf("expected 1 orphaned story (wrong-type parent), got %d", report.OrphanedStories)
	}

	stories, _ := repos.Stories.List(ctx)
	for _, s := range stories {
		if s.Title == "Orphaned by wrong parent" {
			if s.FeatureID != syntheticUnassignedFeatureID {
				t.Fatalf("wrong-parent story assigned to %q, want %q", s.FeatureID, syntheticUnassignedFeatureID)
			}
		}
	}
}

func newRepos(t *testing.T) repository.Repositories {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "maestro.db")
	database, err := maestrodb.Open(dbPath)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	return repository.New(database)
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func intPtr(v int) *int {
	return &v
}

func stringPtr(v string) *string {
	return &v
}
