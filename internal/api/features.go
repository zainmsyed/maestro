package api

import (
	"encoding/json"
	"net/http"
	"time"

	"maestro/internal/models"
)

type featureResponse struct {
	ID               string          `json:"id"`
	EpicID           *string         `json:"epic_id"`
	Title            string          `json:"title"`
	Description      string          `json:"description"`
	Status           string          `json:"status"`
	Owner            string          `json:"owner"`
	Sprint           string          `json:"sprint"`
	StoryPoints      *int            `json:"story_points"`
	OriginalEndDate  *time.Time      `json:"original_end_date"`
	CommittedEndDate *time.Time      `json:"committed_end_date"`
	ActualEndDate    *time.Time      `json:"actual_end_date"`
	DateSource       string          `json:"date_source"`
	Stories          []storyResponse `json:"stories"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

func toFeatureResponse(f models.Feature) featureResponse {
	return featureResponse{
		ID:               f.ID,
		EpicID:           f.EpicID,
		Title:            f.Title,
		Description:      f.Description,
		Status:           f.Status,
		Owner:            f.Owner,
		Sprint:           f.Sprint,
		StoryPoints:      f.StoryPoints,
		OriginalEndDate:  f.OriginalEndDate,
		CommittedEndDate: f.CommittedEndDate,
		ActualEndDate:    f.ActualEndDate,
		DateSource:       f.DateSource,
		CreatedAt:        f.CreatedAt,
		UpdatedAt:        f.UpdatedAt,
	}
}

func (s *Server) listFeatures(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	features, err := s.repos.Features.List(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := make([]featureResponse, len(features))
	for i, f := range features {
		resp[i] = toFeatureResponse(f)
	}
	respondJSON(w, http.StatusOK, resp)
}

func (s *Server) getFeature(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := extractID(r.URL.Path)
	feature, err := s.repos.Features.GetByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "feature not found")
		return
	}

	stories, err := s.repos.Stories.ListByFeatureID(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := toFeatureResponse(*feature)
	resp.Stories = make([]storyResponse, len(stories))
	for i, st := range stories {
		resp.Stories[i] = toStoryResponse(st)
	}
	respondJSON(w, http.StatusOK, resp)
}

type createFeatureRequest struct {
	ID               string `json:"id"`
	EpicID           string `json:"epic_id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	Owner            string `json:"owner"`
	Sprint           string `json:"sprint"`
	StoryPoints      *int   `json:"story_points"`
	OriginalEndDate  string `json:"original_end_date"`
	CommittedEndDate string `json:"committed_end_date"`
}

func (s *Server) createFeature(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createFeatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.ID == "" || req.Title == "" {
		respondError(w, http.StatusBadRequest, "id and title are required")
		return
	}

	original, _ := parseDate(req.OriginalEndDate)
	committed, _ := parseDate(req.CommittedEndDate)

	var epicID *string
	if req.EpicID != "" {
		epicID = &req.EpicID
	}

	feature := &models.Feature{
		ID:               req.ID,
		EpicID:           epicID,
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
	if err := s.repos.Features.Create(ctx, feature); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toFeatureResponse(*feature))
}

func (s *Server) patchFeatureDate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := extractID(r.URL.Path)
	feature, err := s.repos.Features.GetByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "feature not found")
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
	original, committed, deltaDays := computeDatePatch(feature.OriginalEndDate, feature.CommittedEndDate, newDate)

	if err := s.repos.Features.UpdateDate(ctx, id, original, committed, "pm_assigned"); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	audit := &models.DateAuditLog{
		EntityType: "feature",
		EntityID:   id,
		ChangedBy:  changedBy,
		OldDate:    feature.CommittedEndDate,
		NewDate:    newDate,
		DeltaDays:  deltaDays,
		Reason:     req.Reason,
	}
	if err := s.repos.Audits.Create(ctx, audit); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feature, err = s.repos.Features.GetByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toFeatureResponse(*feature))
}

type patchFeatureEpicRequest struct {
	EpicID string `json:"epic_id"`
}

func (s *Server) patchFeatureEpic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := extractID(r.URL.Path)
	_, err := s.repos.Features.GetByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "feature not found")
		return
	}

	var req patchFeatureEpicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	var epicID *string
	if req.EpicID != "" {
		epicID = &req.EpicID
	}

	if err := s.repos.Features.UpdateEpicID(ctx, id, epicID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feature, err := s.repos.Features.GetByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toFeatureResponse(*feature))
}
