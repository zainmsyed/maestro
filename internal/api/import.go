package api

import (
	"net/http"
)

func (s *Server) getImportReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report, err := s.repos.ImportReports.GetLast(ctx)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, report)
}
