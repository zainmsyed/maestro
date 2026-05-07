package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"maestro/internal/api"
	maestrodb "maestro/internal/db"
	"maestro/internal/models"
	"maestro/internal/repository"
)

func newTestServer(t *testing.T) (*api.Server, repository.Repositories) {
	t.Helper()
	db, err := maestrodb.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	repos := repository.New(db)
	return api.New(repos), repos
}

func strPtr(s string) *string { return &s }
func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func TestStoryCRUD(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	// POST
	body := map[string]any{"id": "S-1", "feature_id": "F-1", "title": "Test Story", "status": "New"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/stories", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var created map[string]any
	json.Unmarshal(rec.Body.Bytes(), &created)
	if created["id"] != "S-1" || created["date_source"] != "manual" {
		t.Fatalf("unexpected created story: %v", created)
	}

	// GET list
	req = httptest.NewRequest(http.MethodGet, "/api/stories", nil)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	var list []map[string]any
	json.Unmarshal(rec.Body.Bytes(), &list)
	if len(list) != 1 {
		t.Fatalf("expected 1 story, got %d", len(list))
	}

	// GET by id
	req = httptest.NewRequest(http.MethodGet, "/api/stories/S-1", nil)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	var got map[string]any
	json.Unmarshal(rec.Body.Bytes(), &got)
	if got["title"] != "Test Story" {
		t.Fatalf("unexpected title: %v", got["title"])
	}
}

func TestPatchStoryDate_FirstAssignment(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New"})

	body := map[string]any{"committed_end_date": "2026-05-20", "changed_by": "Zain"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/stories/S-1/date", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var updated map[string]any
	json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated["date_source"] != "pm_assigned" || updated["original_end_date"] != "2026-05-20T00:00:00Z" {
		t.Fatalf("unexpected update: %v", updated)
	}

	audits, _ := repos.Audits.List(ctx)
	if len(audits) != 1 || audits[0].DeltaDays != 0 {
		t.Fatalf("unexpected audit: %+v", audits)
	}
}

func TestPatchStoryDate_SecondAssignment(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New", OriginalEndDate: &original, CommittedEndDate: &committed})

	body := map[string]any{"committed_end_date": "2026-05-20", "changed_by": "Zain", "reason": "Slip"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/stories/S-1/date", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	audits, _ := repos.Audits.List(ctx)
	if len(audits) != 1 || audits[0].DeltaDays != 5 {
		t.Fatalf("unexpected audit: %+v", audits)
	}
}

func TestPatchStoryFeature_Reassign(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "F1", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-2", EpicID: strPtr("E-1"), Title: "F2", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New"})

	body := map[string]any{"feature_id": "F-2"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/stories/S-1/feature", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var updated map[string]any
	json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated["feature_id"] != "F-2" {
		t.Fatalf("unexpected feature_id: %v", updated["feature_id"])
	}
}

func TestFeatureGetByID_WithNestedStories(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active", DateSource: "imported"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story1", Description: "", Status: "New"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-2", FeatureID: "F-1", Title: "Story2", Description: "", Status: "New"})

	req := httptest.NewRequest(http.MethodGet, "/api/features/F-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var feature map[string]any
	json.Unmarshal(rec.Body.Bytes(), &feature)
	if feature["date_source"] != "imported" {
		t.Fatalf("unexpected date_source: %v", feature["date_source"])
	}
	stories := feature["stories"].([]any)
	if len(stories) != 2 {
		t.Fatalf("expected 2 stories, got %d", len(stories))
	}
}

func TestPatchFeatureDate_FirstAssignment(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	body := map[string]any{"committed_end_date": "2026-05-20", "changed_by": "Zain"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/features/F-1/date", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var updated map[string]any
	json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated["date_source"] != "pm_assigned" {
		t.Fatalf("unexpected date_source: %v", updated["date_source"])
	}

	audits, _ := repos.Audits.List(ctx)
	if len(audits) != 1 || audits[0].EntityType != "feature" || audits[0].DeltaDays != 0 {
		t.Fatalf("unexpected audit: %+v", audits)
	}
}

func TestEpicGetByID_WithNestedFeaturesAndStories(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New"})

	req := httptest.NewRequest(http.MethodGet, "/api/epics/E-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var epic map[string]any
	json.Unmarshal(rec.Body.Bytes(), &epic)
	features := epic["features"].([]any)
	if len(features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(features))
	}
	f0 := features[0].(map[string]any)
	stories := f0["stories"].([]any)
	if len(stories) != 1 {
		t.Fatalf("expected 1 story, got %d", len(stories))
	}
}

func TestListEpics_ThreeLevelNesting(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic1", Description: "", Status: "Active"})
	repos.Epics.Create(ctx, &models.Epic{ID: "E-2", Title: "Epic2", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-2", FeatureID: "F-1", Title: "Story2", Description: "", Status: "New"})

	req := httptest.NewRequest(http.MethodGet, "/api/epics", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var epics []map[string]any
	json.Unmarshal(rec.Body.Bytes(), &epics)
	if len(epics) != 2 {
		t.Fatalf("expected 2 epics, got %d", len(epics))
	}
	for _, epic := range epics {
		if epic["id"] == "E-1" {
			features := epic["features"].([]any)
			if len(features) != 1 {
				t.Fatalf("expected 1 feature in E-1, got %d", len(features))
			}
			f0 := features[0].(map[string]any)
			stories := f0["stories"].([]any)
			if len(stories) != 2 {
				t.Fatalf("expected 2 stories in F-1, got %d", len(stories))
			}
		}
	}
}

func TestMetrics_Endpoint(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	actual := date(2026, 5, 14)
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story1", Description: "", Status: "Done", OriginalEndDate: &original, CommittedEndDate: &committed, ActualEndDate: &actual})
	committed2 := date(2026, 5, 20)
	actual2 := date(2026, 5, 25)
	repos.Stories.Create(ctx, &models.Story{ID: "S-2", FeatureID: "F-1", Title: "Story2", Description: "", Status: "Done", OriginalEndDate: &original, CommittedEndDate: &committed2, ActualEndDate: &actual2})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var metrics map[string]any
	json.Unmarshal(rec.Body.Bytes(), &metrics)
	stories := metrics["stories"].(map[string]any)
	if stories["deadline_hit_rate"] != 0.5 {
		t.Fatalf("expected deadline_hit_rate=0.5, got %v", stories["deadline_hit_rate"])
	}
}

func TestMetricsSlip_ExcludesInherited(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature1", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "imported"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-2", EpicID: strPtr("E-1"), Title: "Feature2", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "inherited"})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip/E-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var slip map[string]any
	json.Unmarshal(rec.Body.Bytes(), &slip)
	if slip["item_count"] != float64(1) {
		t.Fatalf("expected item_count=1, got %v", slip["item_count"])
	}
}

func TestMetricsOrphanedStories(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "feature-unassigned", EpicID: strPtr("E-1"), Title: "Unassigned", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Normal", Description: "", Status: "New"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-2", FeatureID: "feature-unassigned", Title: "Orphan", Description: "", Status: "New"})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/orphaned-stories", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var stories []map[string]any
	json.Unmarshal(rec.Body.Bytes(), &stories)
	if len(stories) != 1 || stories[0]["id"] != "S-2" {
		t.Fatalf("unexpected orphans: %v", stories)
	}
}

func TestAudit_WithDateSource(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active", DateSource: "manual"})
	oldDate := date(2026, 5, 1)
	newDate := date(2026, 5, 8)
	repos.Audits.Create(ctx, &models.DateAuditLog{EntityType: "feature", EntityID: "F-1", ChangedBy: "Zain", OldDate: &oldDate, NewDate: &newDate, DeltaDays: 7})

	req := httptest.NewRequest(http.MethodGet, "/api/audit", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var audits []map[string]any
	json.Unmarshal(rec.Body.Bytes(), &audits)
	if len(audits) != 1 || audits[0]["date_source"] != "manual" {
		t.Fatalf("unexpected audits: %v", audits)
	}
}

func TestImportReport_Endpoint(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	report := &models.ImportReport{EpicCount: 2, FeatureCount: 3, StoryCount: 5}
	repos.ImportReports.Save(ctx, report)

	req := httptest.NewRequest(http.MethodGet, "/api/import/report", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var got map[string]any
	json.Unmarshal(rec.Body.Bytes(), &got)
	if got["epic_count"] != float64(2) || got["story_count"] != float64(5) {
		t.Fatalf("unexpected report: %v", got)
	}
}

func TestFeatureEpicReassign(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic1", Description: "", Status: "Active"})
	repos.Epics.Create(ctx, &models.Epic{ID: "E-2", Title: "Epic2", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	body := map[string]any{"epic_id": "E-2"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/features/F-1/epic", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var updated map[string]any
	json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated["epic_id"] != "E-2" {
		t.Fatalf("unexpected epic_id: %v", updated["epic_id"])
	}
}

func TestStoryGetByID_NotFound(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/stories/S-NOT-FOUND", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestFeatureGetByID_NotFound(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/features/F-NOT-FOUND", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestEpicGetByID_NotFound(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/epics/E-NOT-FOUND", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestImportReport_NotFound(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/import/report", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestSlip_NotFound(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip/UNKNOWN", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestMetricsScopeCreep(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story1", Description: "", Status: "New", OriginalEndDate: &original, CommittedEndDate: &committed})
	committed2 := date(2026, 5, 10)
	repos.Stories.Create(ctx, &models.Story{ID: "S-2", FeatureID: "F-1", Title: "Story2", Description: "", Status: "New", OriginalEndDate: &original, CommittedEndDate: &committed2})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var metrics map[string]any
	json.Unmarshal(rec.Body.Bytes(), &metrics)
	stories := metrics["stories"].(map[string]any)
	if stories["scope_creep_rate"] != 0.5 {
		t.Fatalf("expected scope_creep_rate=0.5, got %v", stories["scope_creep_rate"])
	}
}

func TestSlip_FeatureLevel(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story1", Description: "", Status: "New", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "imported"})
	committed2 := date(2026, 5, 20)
	repos.Stories.Create(ctx, &models.Story{ID: "S-2", FeatureID: "F-1", Title: "Story2", Description: "", Status: "New", OriginalEndDate: &original, CommittedEndDate: &committed2, DateSource: "imported"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-3", FeatureID: "F-1", Title: "Story3", Description: "", Status: "New", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "inherited"})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip/F-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var slip map[string]any
	json.Unmarshal(rec.Body.Bytes(), &slip)
	if slip["item_count"] != float64(2) {
		t.Fatalf("expected item_count=2, got %v", slip["item_count"])
	}
	if slip["entity_type"] != "feature" {
		t.Fatalf("expected entity_type=feature, got %v", slip["entity_type"])
	}
	avg := slip["average_slip_days"].(float64)
	if avg != 7.5 {
		t.Fatalf("expected average_slip_days=7.5, got %v", avg)
	}
}

func TestPatchFeatureDate_SecondAssignment(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed})

	body := map[string]any{"committed_end_date": "2026-05-20", "changed_by": "Zain", "reason": "Scope increase"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/features/F-1/date", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	audits, _ := repos.Audits.List(ctx)
	if len(audits) != 1 || audits[0].DeltaDays != 5 || *audits[0].Reason != "Scope increase" {
		t.Fatalf("unexpected audit: %+v", audits)
	}
}

func TestGetEpic_WithNoFeatures(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})

	req := httptest.NewRequest(http.MethodGet, "/api/epics/E-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var epic map[string]any
	json.Unmarshal(rec.Body.Bytes(), &epic)
	features := epic["features"].([]any)
	if len(features) != 0 {
		t.Fatalf("expected 0 features, got %d", len(features))
	}
}

func TestGetFeature_WithNoStories(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	req := httptest.NewRequest(http.MethodGet, "/api/features/F-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var feature map[string]any
	json.Unmarshal(rec.Body.Bytes(), &feature)
	stories := feature["stories"].([]any)
	if len(stories) != 0 {
		t.Fatalf("expected 0 stories, got %d", len(stories))
	}
}

func TestEpicMetrics_WithData(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	actual := date(2026, 5, 14)
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed, ActualEndDate: &actual})
	committed2 := date(2026, 5, 20)
	actual2 := date(2026, 5, 25)
	repos.Epics.Create(ctx, &models.Epic{ID: "E-2", Title: "Epic2", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed2, ActualEndDate: &actual2})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var metrics map[string]any
	json.Unmarshal(rec.Body.Bytes(), &metrics)
	epics := metrics["epics"].(map[string]any)
	if epics["deadline_hit_rate"] != 0.5 {
		t.Fatalf("expected epic deadline_hit_rate=0.5, got %v", epics["deadline_hit_rate"])
	}
}

func TestFeatureMetrics_WithData(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	actual := date(2026, 5, 14)
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature1", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed, ActualEndDate: &actual})
	committed2 := date(2026, 5, 20)
	actual2 := date(2026, 5, 25)
	repos.Features.Create(ctx, &models.Feature{ID: "F-2", EpicID: strPtr("E-1"), Title: "Feature2", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed2, ActualEndDate: &actual2})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var metrics map[string]any
	json.Unmarshal(rec.Body.Bytes(), &metrics)
	features := metrics["features"].(map[string]any)
	if features["deadline_hit_rate"] != 0.5 {
		t.Fatalf("expected feature deadline_hit_rate=0.5, got %v", features["deadline_hit_rate"])
	}
}

func TestImportReport_MultipleSaves(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.ImportReports.Save(ctx, &models.ImportReport{EpicCount: 1, FeatureCount: 1, StoryCount: 1})
	repos.ImportReports.Save(ctx, &models.ImportReport{EpicCount: 2, FeatureCount: 3, StoryCount: 5})

	req := httptest.NewRequest(http.MethodGet, "/api/import/report", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var got map[string]any
	json.Unmarshal(rec.Body.Bytes(), &got)
	if got["epic_count"] != float64(2) {
		t.Fatalf("expected latest epic_count=2, got %v", got["epic_count"])
	}
}

func TestPatchStoryDate_NegativeDelta(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	original := date(2026, 5, 20)
	committed := date(2026, 5, 15)
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New", OriginalEndDate: &original, CommittedEndDate: &committed})

	body := map[string]any{"committed_end_date": "2026-05-10", "changed_by": "Zain"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/stories/S-1/date", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	audits, _ := repos.Audits.List(ctx)
	if len(audits) != 1 || audits[0].DeltaDays != -5 {
		t.Fatalf("unexpected audit: %+v", audits)
	}
}

func TestMetricsSlip_EpicWithMixedSources(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature1", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "imported"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-2", EpicID: strPtr("E-1"), Title: "Feature2", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "inherited"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-3", EpicID: strPtr("E-1"), Title: "Feature3", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "pm_assigned"})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip/E-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var slip map[string]any
	json.Unmarshal(rec.Body.Bytes(), &slip)
	if slip["item_count"] != float64(2) {
		t.Fatalf("expected item_count=2, got %v", slip["item_count"])
	}
}

func TestOrphanedStories_TrueOrphan(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Normal", Description: "", Status: "New"})
	// Create synthetic unassigned feature and story under it
	repos.Features.Create(ctx, &models.Feature{ID: "feature-unassigned", EpicID: strPtr("E-1"), Title: "Unassigned", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-ORPHAN", FeatureID: "feature-unassigned", Title: "Orphan", Description: "", Status: "New"})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/orphaned-stories", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var stories []map[string]any
	json.Unmarshal(rec.Body.Bytes(), &stories)
	if len(stories) != 1 || stories[0]["id"] != "S-ORPHAN" {
		t.Fatalf("unexpected orphans: %v", stories)
	}
}

func TestCreateStory_WithDates(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	body := map[string]any{"id": "S-DATED", "feature_id": "F-1", "title": "Dated Story", "status": "New", "original_end_date": "2026-05-10", "committed_end_date": "2026-05-15"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/stories", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	story, _ := repos.Stories.GetByID(ctx, "S-DATED")
	if story.OriginalEndDate == nil || !story.OriginalEndDate.Equal(date(2026, 5, 10)) {
		t.Fatalf("unexpected original_end_date: %v", story.OriginalEndDate)
	}
}

func TestCreateFeature_WithDates(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})

	body := map[string]any{"id": "F-DATED", "epic_id": "E-1", "title": "Dated Feature", "status": "New", "original_end_date": "2026-05-10", "committed_end_date": "2026-05-15"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/features", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	feature, _ := repos.Features.GetByID(ctx, "F-DATED")
	if feature.OriginalEndDate == nil || !feature.OriginalEndDate.Equal(date(2026, 5, 10)) {
		t.Fatalf("unexpected original_end_date: %v", feature.OriginalEndDate)
	}
}

func TestCreateEpic_WithDates(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()

	body := map[string]any{"id": "E-DATED", "title": "Dated Epic", "status": "Active", "original_end_date": "2026-05-10", "committed_end_date": "2026-05-15"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/epics", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	epic, _ := repos.Epics.GetByID(ctx, "E-DATED")
	if epic.OriginalEndDate == nil || !epic.OriginalEndDate.Equal(date(2026, 5, 10)) {
		t.Fatalf("unexpected original_end_date: %v", epic.OriginalEndDate)
	}
}

func TestMetricsSlip_EpicWithFeaturesAndStories(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature1", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "imported"})
	committed2 := date(2026, 5, 20)
	repos.Features.Create(ctx, &models.Feature{ID: "F-2", EpicID: strPtr("E-1"), Title: "Feature2", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed2, DateSource: "imported"})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip/E-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var slip map[string]any
	json.Unmarshal(rec.Body.Bytes(), &slip)
	if slip["item_count"] != float64(2) {
		t.Fatalf("expected item_count=2, got %v", slip["item_count"])
	}
	avg := slip["average_slip_days"].(float64)
	if avg != 7.5 {
		t.Fatalf("expected average_slip_days=7.5, got %v", avg)
	}
}

func TestFeatureGetByID_NestedStoryOrder(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-B", FeatureID: "F-1", Title: "Story B", Description: "", Status: "New"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-A", FeatureID: "F-1", Title: "Story A", Description: "", Status: "New"})

	req := httptest.NewRequest(http.MethodGet, "/api/features/F-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var feature map[string]any
	json.Unmarshal(rec.Body.Bytes(), &feature)
	stories := feature["stories"].([]any)
	if len(stories) != 2 {
		t.Fatalf("expected 2 stories, got %d", len(stories))
	}
	s0 := stories[0].(map[string]any)
	s1 := stories[1].(map[string]any)
	if s0["id"] != "S-A" || s1["id"] != "S-B" {
		t.Fatalf("expected sorted order S-A, S-B, got %v, %v", s0["id"], s1["id"])
	}
}

func TestListEpics_EpicOrder(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-B", Title: "Epic B", Description: "", Status: "Active"})
	repos.Epics.Create(ctx, &models.Epic{ID: "E-A", Title: "Epic A", Description: "", Status: "Active"})

	req := httptest.NewRequest(http.MethodGet, "/api/epics", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var epics []map[string]any
	json.Unmarshal(rec.Body.Bytes(), &epics)
	if len(epics) != 2 {
		t.Fatalf("expected 2 epics, got %d", len(epics))
	}
	if epics[0]["id"] != "E-A" || epics[1]["id"] != "E-B" {
		t.Fatalf("expected sorted order E-A, E-B, got %v, %v", epics[0]["id"], epics[1]["id"])
	}
}

func TestMetricsSlip_FeatureLevelExcludesInherited(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story1", Description: "", Status: "New", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "imported"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-2", FeatureID: "F-1", Title: "Story2", Description: "", Status: "New", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "inherited"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-3", FeatureID: "F-1", Title: "Story3", Description: "", Status: "New", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "inherited"})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip/F-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var slip map[string]any
	json.Unmarshal(rec.Body.Bytes(), &slip)
	if slip["item_count"] != float64(1) {
		t.Fatalf("expected item_count=1, got %v", slip["item_count"])
	}
}

func TestGetFeature_WithEpicID(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	req := httptest.NewRequest(http.MethodGet, "/api/features/F-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var feature map[string]any
	json.Unmarshal(rec.Body.Bytes(), &feature)
	if feature["epic_id"] != "E-1" {
		t.Fatalf("expected epic_id=E-1, got %v", feature["epic_id"])
	}
}

func TestGetStory_WithFeatureID(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New"})

	req := httptest.NewRequest(http.MethodGet, "/api/stories/S-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var story map[string]any
	json.Unmarshal(rec.Body.Bytes(), &story)
	if story["feature_id"] != "F-1" {
		t.Fatalf("expected feature_id=F-1, got %v", story["feature_id"])
	}
}

func TestGetEpic_NestedStoryDateSource(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New", DateSource: "pm_assigned"})

	req := httptest.NewRequest(http.MethodGet, "/api/epics/E-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var epic map[string]any
	json.Unmarshal(rec.Body.Bytes(), &epic)
	features := epic["features"].([]any)
	f0 := features[0].(map[string]any)
	stories := f0["stories"].([]any)
	s0 := stories[0].(map[string]any)
	if s0["date_source"] != "pm_assigned" {
		t.Fatalf("expected story date_source=pm_assigned, got %v", s0["date_source"])
	}
}

func TestGetEpic_NestedFeatureDateSource(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active", DateSource: "pm_assigned"})

	req := httptest.NewRequest(http.MethodGet, "/api/epics/E-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var epic map[string]any
	json.Unmarshal(rec.Body.Bytes(), &epic)
	features := epic["features"].([]any)
	f0 := features[0].(map[string]any)
	if f0["date_source"] != "pm_assigned" {
		t.Fatalf("expected feature date_source=pm_assigned, got %v", f0["date_source"])
	}
}

func TestAudit_WithEpicEntityType(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	oldDate := date(2026, 5, 1)
	newDate := date(2026, 5, 8)
	repos.Audits.Create(ctx, &models.DateAuditLog{EntityType: "epic", EntityID: "E-1", ChangedBy: "Zain", OldDate: &oldDate, NewDate: &newDate, DeltaDays: 7})

	req := httptest.NewRequest(http.MethodGet, "/api/audit", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var audits []map[string]any
	json.Unmarshal(rec.Body.Bytes(), &audits)
	if len(audits) != 1 {
		t.Fatalf("unexpected audits: %v", audits)
	}
	// Epics don't have date_source, so the key may be missing (omitempty)
	ds, ok := audits[0]["date_source"]
	if ok && ds != "" {
		t.Fatalf("expected empty date_source for epic, got %v", ds)
	}
}

func TestMetricsSlip_FeatureWithNullDates(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New", DateSource: "imported"})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip/F-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var slip map[string]any
	json.Unmarshal(rec.Body.Bytes(), &slip)
	if slip["item_count"] != float64(0) {
		t.Fatalf("expected item_count=0, got %v", slip["item_count"])
	}
}

func TestMetricsSlip_EpicWithNullDates(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active", DateSource: "imported"})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip/E-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var slip map[string]any
	json.Unmarshal(rec.Body.Bytes(), &slip)
	if slip["item_count"] != float64(0) {
		t.Fatalf("expected item_count=0, got %v", slip["item_count"])
	}
}

func TestPatchStoryFeature_MissingFeatureID(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New"})

	body := map[string]any{}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/stories/S-1/feature", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestPatchFeatureEpic_MissingEpicID(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	body := map[string]any{}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/features/F-1/epic", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 (clears epic), got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestGetStory_WithDateSource(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New", DateSource: "inherited"})

	req := httptest.NewRequest(http.MethodGet, "/api/stories/S-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var story map[string]any
	json.Unmarshal(rec.Body.Bytes(), &story)
	if story["date_source"] != "inherited" {
		t.Fatalf("expected date_source=inherited, got %v", story["date_source"])
	}
}

func TestListEpics_FeatureWithoutStories(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	req := httptest.NewRequest(http.MethodGet, "/api/epics", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var epics []map[string]any
	json.Unmarshal(rec.Body.Bytes(), &epics)
	if len(epics) != 1 {
		t.Fatalf("expected 1 epic, got %d", len(epics))
	}
	features := epics[0]["features"].([]any)
	if len(features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(features))
	}
	f0 := features[0].(map[string]any)
	stories := f0["stories"].([]any)
	if len(stories) != 0 {
		t.Fatalf("expected 0 stories, got %d", len(stories))
	}
}

func TestMetricsSlip_EpicAllInherited(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	original := date(2026, 5, 10)
	committed := date(2026, 5, 15)
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active", OriginalEndDate: &original, CommittedEndDate: &committed, DateSource: "inherited"})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip/E-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var slip map[string]any
	json.Unmarshal(rec.Body.Bytes(), &slip)
	if slip["item_count"] != float64(0) {
		t.Fatalf("expected item_count=0, got %v", slip["item_count"])
	}
	if slip["average_slip_days"] != 0.0 {
		t.Fatalf("expected average_slip_days=0, got %v", slip["average_slip_days"])
	}
}

func TestPatchFeatureEpic_MissingFeature(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	body := map[string]any{"epic_id": "E-1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/features/F-NOT-FOUND/epic", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestPatchStoryFeature_MissingStory(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	body := map[string]any{"feature_id": "F-1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/stories/S-NOT-FOUND/feature", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestFeaturePatch_WithInvalidEpicID(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	body := map[string]any{"epic_id": "E-NOT-FOUND"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/features/F-1/epic", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 due to FK violation, got %d", rec.Code)
	}
}

func TestCreateStory_DuplicateID(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New"})

	body := map[string]any{"id": "S-1", "feature_id": "F-1", "title": "Duplicate", "status": "New"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/stories", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for duplicate id, got %d", rec.Code)
	}
}

func TestMetrics_NoData(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var metrics map[string]any
	json.Unmarshal(rec.Body.Bytes(), &metrics)
	stories := metrics["stories"].(map[string]any)
	if stories["deadline_hit_rate"] != 0.0 || stories["scope_creep_rate"] != 0.0 {
		t.Fatalf("expected 0 rates with no data, got %v", stories)
	}
}

func TestMetricsSlip_FeatureNoStories(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip/F-1", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var slip map[string]any
	json.Unmarshal(rec.Body.Bytes(), &slip)
	if slip["item_count"] != float64(0) {
		t.Fatalf("expected item_count=0, got %v", slip["item_count"])
	}
}

func TestPatchStoryDate_MissingDate(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New"})

	body := map[string]any{"changed_by": "Zain"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/stories/S-1/date", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestPatchFeatureDate_InvalidDate(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	body := map[string]any{"committed_end_date": "not-a-date"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/features/F-1/date", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestUnknownEndpoint(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestPatchStoryDate_InvalidJSON(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New"})
	req := httptest.NewRequest(http.MethodPatch, "/api/stories/S-1/date", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestPatchFeatureEpic_InvalidJSON(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	req := httptest.NewRequest(http.MethodPatch, "/api/features/F-1/epic", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateStory_InvalidJSON(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/api/stories", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateEpic_MissingTitle(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	body := map[string]any{"id": "E-1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/epics", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateFeature_MissingTitle(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	body := map[string]any{"id": "F-1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/features", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateStory_MissingTitle(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	body := map[string]any{"id": "S-1", "feature_id": "F-1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/stories", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestPatchStoryDate_NotFound(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	body := map[string]any{"committed_end_date": "2026-05-20"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/stories/S-NOT-FOUND/date", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestPatchFeatureDate_NotFound(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	body := map[string]any{"committed_end_date": "2026-05-20"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/features/F-NOT-FOUND/date", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestAuditList_Empty(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/audit", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var audits []map[string]any
	json.Unmarshal(rec.Body.Bytes(), &audits)
	if len(audits) != 0 {
		t.Fatalf("expected 0 audits, got %d", len(audits))
	}
}

func TestMetricsOrphanedStories_Empty(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/orphaned-stories", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var stories []map[string]any
	json.Unmarshal(rec.Body.Bytes(), &stories)
	if len(stories) != 0 {
		t.Fatalf("expected 0 stories, got %d", len(stories))
	}
}

func TestMetricsSlip_EpicNotFound(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip/E-NOT-FOUND", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetStoryDate_SubResource(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New"})

	req := httptest.NewRequest(http.MethodGet, "/api/stories/S-1/date", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for GET /date, got %d", rec.Code)
	}
}

func TestGetStoryFeature_SubResource(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})
	repos.Stories.Create(ctx, &models.Story{ID: "S-1", FeatureID: "F-1", Title: "Story", Description: "", Status: "New"})

	req := httptest.NewRequest(http.MethodGet, "/api/stories/S-1/feature", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for GET /feature, got %d", rec.Code)
	}
}

func TestGetFeatureDate_SubResource(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	req := httptest.NewRequest(http.MethodGet, "/api/features/F-1/date", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for GET /date, got %d", rec.Code)
	}
}

func TestGetFeatureEpic_SubResource(t *testing.T) {
	t.Parallel()
	server, repos := newTestServer(t)
	ctx := t.Context()
	repos.Epics.Create(ctx, &models.Epic{ID: "E-1", Title: "Epic", Description: "", Status: "Active"})
	repos.Features.Create(ctx, &models.Feature{ID: "F-1", EpicID: strPtr("E-1"), Title: "Feature", Description: "", Status: "Active"})

	req := httptest.NewRequest(http.MethodGet, "/api/features/F-1/epic", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for GET /epic, got %d", rec.Code)
	}
}

func TestGetMetricsSlip_MissingID(t *testing.T) {
	t.Parallel()
	server, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/slip", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for GET /metrics/slip without id, got %d", rec.Code)
	}
}
