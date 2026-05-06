package repository_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	maestrodb "maestro/internal/db"
	"maestro/internal/models"
	"maestro/internal/repository"
)

func TestEpicCRUD(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	original := date(2026, 5, 20)
	committed := date(2026, 5, 27)
	epic := &models.Epic{
		ID:               "E-1",
		Title:            "Platform Roadmap",
		Description:      "Top-level epic",
		Status:           "Active",
		Owner:            "Zain",
		SprintStart:      "Sprint 12",
		SprintEnd:        "Sprint 13",
		OriginalEndDate:  &original,
		CommittedEndDate: &committed,
		IsSynthetic:      false,
	}
	if err := repos.Epics.Create(ctx, epic); err != nil {
		t.Fatalf("create epic: %v", err)
	}

	got, err := repos.Epics.GetByID(ctx, epic.ID)
	if err != nil {
		t.Fatalf("get epic: %v", err)
	}
	if got.Title != epic.Title || got.Status != epic.Status || got.SprintEnd != epic.SprintEnd {
		t.Fatalf("unexpected epic: %#v", got)
	}
	if got.OriginalEndDate == nil || !got.OriginalEndDate.Equal(original) {
		t.Fatalf("unexpected original date: %#v", got.OriginalEndDate)
	}

	list, err := repos.Epics.List(ctx)
	if err != nil {
		t.Fatalf("list epics: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 epic, got %d", len(list))
	}
}

func TestFeatureCRUD(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	epicID := "E-1"
	if err := repos.Epics.Create(ctx, &models.Epic{ID: epicID, Title: "Epic", Description: "", Status: "Active"}); err != nil {
		t.Fatalf("seed epic: %v", err)
	}
	points := 8
	original := date(2026, 5, 14)
	feature := &models.Feature{
		ID:               "F-1",
		EpicID:           &epicID,
		Title:            "Import parser",
		Description:      "Build parser",
		Status:           "New",
		Owner:            "Alex",
		Sprint:           "Sprint 12",
		OriginalEndDate:  &original,
		CommittedEndDate: &original,
		StoryPoints:      &points,
	}
	if err := repos.Features.Create(ctx, feature); err != nil {
		t.Fatalf("create feature: %v", err)
	}

	got, err := repos.Features.GetByID(ctx, feature.ID)
	if err != nil {
		t.Fatalf("get feature: %v", err)
	}
	if got.EpicID == nil || *got.EpicID != epicID {
		t.Fatalf("unexpected epic id: %#v", got.EpicID)
	}
	if got.StoryPoints == nil || *got.StoryPoints != points {
		t.Fatalf("unexpected story points: %#v", got.StoryPoints)
	}

	list, err := repos.Features.List(ctx)
	if err != nil {
		t.Fatalf("list features: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(list))
	}
}

func TestSprintCRUD(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	start := date(2026, 4, 7)
	end := date(2026, 4, 21)
	capacity := 34
	sprint := &models.Sprint{
		ID:        "S12",
		Name:      "Sprint 12",
		StartDate: &start,
		EndDate:   &end,
		Team:      "Platform",
		Capacity:  &capacity,
		Source:    "imported",
	}
	if err := repos.Sprints.Create(ctx, sprint); err != nil {
		t.Fatalf("create sprint: %v", err)
	}

	got, err := repos.Sprints.GetByID(ctx, sprint.ID)
	if err != nil {
		t.Fatalf("get sprint: %v", err)
	}
	if got.Capacity == nil || *got.Capacity != capacity || got.Source != "imported" {
		t.Fatalf("unexpected sprint: %#v", got)
	}

	list, err := repos.Sprints.List(ctx)
	if err != nil {
		t.Fatalf("list sprints: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 sprint, got %d", len(list))
	}
}

func TestAuditCRUD(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	oldDate := date(2026, 5, 1)
	newDate := date(2026, 5, 8)
	reason := "Dependency slip"
	audit := &models.DateAuditLog{
		EntityType: "feature",
		EntityID:   "F-1",
		ChangedBy:  "Zain",
		OldDate:    &oldDate,
		NewDate:    &newDate,
		DeltaDays:  7,
		Reason:     &reason,
	}
	if err := repos.Audits.Create(ctx, audit); err != nil {
		t.Fatalf("create audit: %v", err)
	}
	if audit.ID == 0 {
		t.Fatalf("expected audit id to be set")
	}

	list, err := repos.Audits.List(ctx)
	if err != nil {
		t.Fatalf("list audits: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 audit, got %d", len(list))
	}
	if list[0].Reason == nil || *list[0].Reason != reason {
		t.Fatalf("unexpected audit reason: %#v", list[0].Reason)
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

func TestEpicGetByID_NotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	_, err := repos.Epics.GetByID(ctx, "E-NOT-FOUND")
	if err == nil {
		t.Fatalf("expected error for missing epic, got nil")
	}
}

func TestEpicCreate_DuplicateID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	epic := &models.Epic{ID: "E-DUP", Title: "First", Description: "", Status: "Active"}
	if err := repos.Epics.Create(ctx, epic); err != nil {
		t.Fatalf("create first epic: %v", err)
	}

	duplicate := &models.Epic{ID: "E-DUP", Title: "Second", Description: "", Status: "New"}
	if err := repos.Epics.Create(ctx, duplicate); err == nil {
		t.Fatalf("expected error for duplicate epic id, got nil")
	}
}

func TestEpicIsSynthetic_RoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	epic := &models.Epic{ID: "E-SYN", Title: "Unassigned", Description: "", Status: "Active", IsSynthetic: true}
	if err := repos.Epics.Create(ctx, epic); err != nil {
		t.Fatalf("create synthetic epic: %v", err)
	}

	got, err := repos.Epics.GetByID(ctx, epic.ID)
	if err != nil {
		t.Fatalf("get synthetic epic: %v", err)
	}
	if !got.IsSynthetic {
		t.Fatalf("expected IsSynthetic=true, got false")
	}

	regular := &models.Epic{ID: "E-REG", Title: "Regular", Description: "", Status: "Active", IsSynthetic: false}
	if err := repos.Epics.Create(ctx, regular); err != nil {
		t.Fatalf("create regular epic: %v", err)
	}

	got2, err := repos.Epics.GetByID(ctx, regular.ID)
	if err != nil {
		t.Fatalf("get regular epic: %v", err)
	}
	if got2.IsSynthetic {
		t.Fatalf("expected IsSynthetic=false, got true")
	}
}

func TestFeatureGetByID_NotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	_, err := repos.Features.GetByID(ctx, "F-NOT-FOUND")
	if err == nil {
		t.Fatalf("expected error for missing feature, got nil")
	}
}

func TestSprintGetByID_NotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	_, err := repos.Sprints.GetByID(ctx, "S-NOT-FOUND")
	if err == nil {
		t.Fatalf("expected error for missing sprint, got nil")
	}
}

func TestAuditCreate_NilDatesAndReason(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := newRepos(t)

	audit := &models.DateAuditLog{
		EntityType: "feature",
		EntityID:   "F-1",
		ChangedBy:  "Zain",
		OldDate:    nil,
		NewDate:    nil,
		DeltaDays:  0,
		Reason:     nil,
	}
	if err := repos.Audits.Create(ctx, audit); err != nil {
		t.Fatalf("create audit with nil fields: %v", err)
	}
	if audit.ID == 0 {
		t.Fatalf("expected audit id to be set")
	}

	list, err := repos.Audits.List(ctx)
	if err != nil {
		t.Fatalf("list audits: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 audit, got %d", len(list))
	}
	if list[0].OldDate != nil || list[0].NewDate != nil || list[0].Reason != nil {
		t.Fatalf("expected nil fields to round-trip as nil")
	}
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
