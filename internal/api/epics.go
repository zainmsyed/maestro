package api

import (
	"encoding/json"
	"net/http"
	"time"

	"maestro/internal/models"
)

type epicResponse struct {
	ID               string            `json:"id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Status           string            `json:"status"`
	Owner            string            `json:"owner"`
	SprintStart      string            `json:"sprint_start"`
	SprintEnd        string            `json:"sprint_end"`
	OriginalEndDate  *time.Time        `json:"original_end_date"`
	CommittedEndDate *time.Time        `json:"committed_end_date"`
	ActualEndDate    *time.Time        `json:"actual_end_date"`
	IsSynthetic      bool              `json:"is_synthetic"`
	Features         []featureResponse `json:"features"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

func toEpicResponse(e models.Epic) epicResponse {
	return epicResponse{
		ID:               e.ID,
		Title:            e.Title,
		Description:      e.Description,
		Status:           e.Status,
		Owner:            e.Owner,
		SprintStart:      e.SprintStart,
		SprintEnd:        e.SprintEnd,
		OriginalEndDate:  e.OriginalEndDate,
		CommittedEndDate: e.CommittedEndDate,
		ActualEndDate:    e.ActualEndDate,
		IsSynthetic:      e.IsSynthetic,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}

func (s *Server) listEpics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	epics, err := s.repos.Epics.List(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	allFeatures, err := s.repos.Features.List(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	featureMap := make(map[string][]models.Feature)
	for _, f := range allFeatures {
		if f.EpicID != nil {
			featureMap[*f.EpicID] = append(featureMap[*f.EpicID], f)
		}
	}

	allStories, err := s.repos.Stories.List(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	storyMap := make(map[string][]models.Story)
	for _, st := range allStories {
		storyMap[st.FeatureID] = append(storyMap[st.FeatureID], st)
	}

	resp := make([]epicResponse, len(epics))
	for i, epic := range epics {
		er := toEpicResponse(epic)
		features := featureMap[epic.ID]
		er.Features = make([]featureResponse, len(features))
		for j, f := range features {
			fr := toFeatureResponse(f)
			stories := storyMap[f.ID]
			fr.Stories = make([]storyResponse, len(stories))
			for k, st := range stories {
				fr.Stories[k] = toStoryResponse(st)
			}
			er.Features[j] = fr
		}
		resp[i] = er
	}
	respondJSON(w, http.StatusOK, resp)
}

func (s *Server) getEpic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := extractID(r.URL.Path)
	epic, err := s.repos.Epics.GetByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "epic not found")
		return
	}

	resp := toEpicResponse(*epic)
	features, err := s.repos.Features.ListByEpicID(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Batch-fetch stories for all features in one query-like pass
	allStories, err := s.repos.Stories.List(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	featureIDs := make(map[string]bool, len(features))
	for _, f := range features {
		featureIDs[f.ID] = true
	}
	storyMap := make(map[string][]models.Story)
	for _, st := range allStories {
		if featureIDs[st.FeatureID] {
			storyMap[st.FeatureID] = append(storyMap[st.FeatureID], st)
		}
	}

	resp.Features = make([]featureResponse, len(features))
	for i, f := range features {
		fr := toFeatureResponse(f)
		stories := storyMap[f.ID]
		fr.Stories = make([]storyResponse, len(stories))
		for j, st := range stories {
			fr.Stories[j] = toStoryResponse(st)
		}
		resp.Features[i] = fr
	}
	respondJSON(w, http.StatusOK, resp)
}

type createEpicRequest struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	Owner            string `json:"owner"`
	SprintStart      string `json:"sprint_start"`
	SprintEnd        string `json:"sprint_end"`
	OriginalEndDate  string `json:"original_end_date"`
	CommittedEndDate string `json:"committed_end_date"`
	IsSynthetic      bool   `json:"is_synthetic"`
}

func (s *Server) createEpic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createEpicRequest
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

	epic := &models.Epic{
		ID:               req.ID,
		Title:            req.Title,
		Description:      req.Description,
		Status:           req.Status,
		Owner:            req.Owner,
		SprintStart:      req.SprintStart,
		SprintEnd:        req.SprintEnd,
		OriginalEndDate:  original,
		CommittedEndDate: committed,
		IsSynthetic:      req.IsSynthetic,
	}
	if err := s.repos.Epics.Create(ctx, epic); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toEpicResponse(*epic))
}
