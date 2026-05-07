package api

import (
	"net/http"
	"strings"
)

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/stories", s.listStories)
	s.mux.HandleFunc("GET /api/stories/", s.handleStory)
	s.mux.HandleFunc("POST /api/stories", s.createStory)
	s.mux.HandleFunc("PATCH /api/stories/", s.handleStoryPatch)

	s.mux.HandleFunc("GET /api/features", s.listFeatures)
	s.mux.HandleFunc("GET /api/features/", s.handleFeature)
	s.mux.HandleFunc("POST /api/features", s.createFeature)
	s.mux.HandleFunc("PATCH /api/features/", s.handleFeaturePatch)

	s.mux.HandleFunc("GET /api/epics", s.listEpics)
	s.mux.HandleFunc("GET /api/epics/", s.handleEpic)
	s.mux.HandleFunc("POST /api/epics", s.createEpic)

	s.mux.HandleFunc("GET /api/metrics", s.getMetrics)
	s.mux.HandleFunc("GET /api/metrics/", s.handleMetrics)

	s.mux.HandleFunc("GET /api/audit", s.listAudit)

	s.mux.HandleFunc("GET /api/import/report", s.getImportReport)
}

func extractID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

func (s *Server) handleStory(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 {
		respondError(w, http.StatusNotFound, "not found")
		return
	}
	s.getStory(w, r)
}

func (s *Server) handleStoryPatch(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/date") {
		s.patchStoryDate(w, r)
		return
	}
	if strings.HasSuffix(r.URL.Path, "/feature") {
		s.patchStoryFeature(w, r)
		return
	}
	respondError(w, http.StatusNotFound, "unknown story patch endpoint")
}

func (s *Server) handleFeature(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 {
		respondError(w, http.StatusNotFound, "not found")
		return
	}
	s.getFeature(w, r)
}

func (s *Server) handleFeaturePatch(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/date") {
		s.patchFeatureDate(w, r)
		return
	}
	if strings.HasSuffix(r.URL.Path, "/epic") {
		s.patchFeatureEpic(w, r)
		return
	}
	respondError(w, http.StatusNotFound, "unknown feature patch endpoint")
}

func (s *Server) handleEpic(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 {
		respondError(w, http.StatusNotFound, "not found")
		return
	}
	s.getEpic(w, r)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 3 && parts[2] == "slip" {
		respondError(w, http.StatusNotFound, "missing id")
		return
	}
	if strings.HasSuffix(r.URL.Path, "/orphaned-stories") {
		s.getOrphanedStories(w, r)
		return
	}
	if strings.Contains(r.URL.Path, "/slip/") {
		s.getSlip(w, r)
		return
	}
	respondError(w, http.StatusNotFound, "unknown metrics endpoint")
}
