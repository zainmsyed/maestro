package api

import (
	"net/http"
	"strings"
)

type metricsResponse struct {
	Epics struct {
		DeadlineHitRate float64 `json:"deadline_hit_rate"`
		ScopeCreepRate  float64 `json:"scope_creep_rate"`
	} `json:"epics"`
	Features struct {
		DeadlineHitRate float64 `json:"deadline_hit_rate"`
		ScopeCreepRate  float64 `json:"scope_creep_rate"`
	} `json:"features"`
	Stories struct {
		DeadlineHitRate float64 `json:"deadline_hit_rate"`
		ScopeCreepRate  float64 `json:"scope_creep_rate"`
	} `json:"stories"`
}

func (s *Server) getMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var resp metricsResponse

	epicDhr, err := s.repos.Metrics.EpicDeadlineHitRate(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	epicScr, err := s.repos.Metrics.EpicScopeCreepRate(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Epics.DeadlineHitRate = epicDhr
	resp.Epics.ScopeCreepRate = epicScr

	featDhr, err := s.repos.Metrics.FeatureDeadlineHitRate(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	featScr, err := s.repos.Metrics.FeatureScopeCreepRate(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Features.DeadlineHitRate = featDhr
	resp.Features.ScopeCreepRate = featScr

	storyDhr, err := s.repos.Metrics.StoryDeadlineHitRate(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	storyScr, err := s.repos.Metrics.StoryScopeCreepRate(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Stories.DeadlineHitRate = storyDhr
	resp.Stories.ScopeCreepRate = storyScr

	respondJSON(w, http.StatusOK, resp)
}

func (s *Server) getSlip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		respondError(w, http.StatusBadRequest, "missing id")
		return
	}
	id := parts[3]

	metrics, err := s.repos.Metrics.SlipByID(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, metrics)
}

func (s *Server) getOrphanedStories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stories, err := s.repos.Metrics.OrphanedStories(ctx)
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
