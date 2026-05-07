package api

import (
	"encoding/json"
	"net/http"
	"time"

	"maestro/internal/models"
)

type storyResponse struct {
	ID               string     `json:"id"`
	FeatureID        string     `json:"feature_id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`
	Owner            string     `json:"owner"`
	Sprint           string     `json:"sprint"`
	StoryPoints      *int       `json:"story_points"`
	OriginalEndDate  *time.Time `json:"original_end_date"`
	CommittedEndDate *time.Time `json:"committed_end_date"`
	ActualEndDate    *time.Time `json:"actual_end_date"`
	DateSource       string     `json:"date_source"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func toStoryResponse(s models.Story) storyResponse {
	return storyResponse{
		ID:               s.ID,
		FeatureID:        s.FeatureID,
		Title:            s.Title,
		Description:      s.Description,
		Status:           s.Status,
		Owner:            s.Owner,
		Sprint:           s.Sprint,
		StoryPoints:      s.StoryPoints,
		OriginalEndDate:  s.OriginalEndDate,
		CommittedEndDate: s.CommittedEndDate,
		ActualEndDate:    s.ActualEndDate,
		DateSource:       s.DateSource,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}

func (s *Server) listStories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stories, err := s.repos.Stories.List(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := make([]storyResponse, len(stories))
	for i, st := range stories {
		resp[i] = toStoryResponse(st)
	}
	respondJSON(w, http.StatusOK, resp)
}

func (s *Server) getStory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := extractID(r.URL.Path)
	story, err := s.repos.Stories.GetByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "story not found")
		return
	}
	respondJSON(w, http.StatusOK, toStoryResponse(*story))
}

type createStoryRequest struct {
	ID               string `json:"id"`
	FeatureID        string `json:"feature_id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	Owner            string `json:"owner"`
	Sprint           string `json:"sprint"`
	StoryPoints      *int   `json:"story_points"`
	OriginalEndDate  string `json:"original_end_date"`
	CommittedEndDate string `json:"committed_end_date"`
}

func (s *Server) createStory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createStoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.ID == "" || req.FeatureID == "" || req.Title == "" {
		respondError(w, http.StatusBadRequest, "id, feature_id, and title are required")
		return
	}

	original, _ := parseDate(req.OriginalEndDate)
	committed, _ := parseDate(req.CommittedEndDate)

	story := &models.Story{
		ID:               req.ID,
		FeatureID:        req.FeatureID,
		Title:            req.Title,
		Description:      req.Description,
		Status:           req.Status,
		Owner:            req.Owner,
		Sprint:           req.Sprint,
		StoryPoints:      req.StoryPoints,
		OriginalEndDate:  original,
		CommittedEndDate: committed,
		DateSource:       "manual",
	}
	if err := s.repos.Stories.Create(ctx, story); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toStoryResponse(*story))
}

type patchDateRequest struct {
	CommittedEndDate string  `json:"committed_end_date"`
	Reason           *string `json:"reason"`
	ChangedBy        string  `json:"changed_by"`
}

func changedByOrDefault(v string) string {
	if v == "" {
		return "system"
	}
	return v
}

func computeDatePatch(originalEndDate, committedEndDate, newDate *time.Time) (original, committed *time.Time, deltaDays int) {
	if originalEndDate == nil {
		return newDate, newDate, 0
	}
	if committedEndDate != nil {
		return originalEndDate, newDate, int(newDate.Sub(*committedEndDate).Hours() / 24)
	}
	return originalEndDate, newDate, int(newDate.Sub(*originalEndDate).Hours() / 24)
}

func (s *Server) patchStoryDate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := extractID(r.URL.Path)
	story, err := s.repos.Stories.GetByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "story not found")
		return
	}

	var req patchDateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	newDate, err := parseDate(req.CommittedEndDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid date format")
		return
	}
	if newDate == nil {
		respondError(w, http.StatusBadRequest, "committed_end_date is required")
		return
	}

	changedBy := changedByOrDefault(req.ChangedBy)
	original, committed, deltaDays := computeDatePatch(story.OriginalEndDate, story.CommittedEndDate, newDate)

	if err := s.repos.Stories.UpdateDate(ctx, id, original, committed, "pm_assigned"); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	audit := &models.DateAuditLog{
		EntityType: "story",
		EntityID:   id,
		ChangedBy:  changedBy,
		OldDate:    story.CommittedEndDate,
		NewDate:    newDate,
		DeltaDays:  deltaDays,
		Reason:     req.Reason,
	}
	if err := s.repos.Audits.Create(ctx, audit); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	story, err = s.repos.Stories.GetByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toStoryResponse(*story))
}

type patchStoryFeatureRequest struct {
	FeatureID string `json:"feature_id"`
}

func (s *Server) patchStoryFeature(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := extractID(r.URL.Path)
	_, err := s.repos.Stories.GetByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "story not found")
		return
	}

	var req patchStoryFeatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.FeatureID == "" {
		respondError(w, http.StatusBadRequest, "feature_id is required")
		return
	}

	if err := s.repos.Stories.UpdateFeatureID(ctx, id, req.FeatureID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	story, err := s.repos.Stories.GetByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toStoryResponse(*story))
}
