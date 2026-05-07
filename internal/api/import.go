package api

import (
	"fmt"
	"net/http"
	"strings"

	"maestro/internal/importer"
)

const maxImportUploadBytes = 20 << 20 // 20 MiB

func (s *Server) postImport(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxImportUploadBytes)
	if err := r.ParseMultipartForm(maxImportUploadBytes); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("parse upload: %v", err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "missing CSV file field \"file\"")
		return
	}
	defer file.Close()

	if header != nil && !strings.HasSuffix(strings.ToLower(header.Filename), ".csv") {
		respondError(w, http.StatusBadRequest, "unsupported format: upload a CSV file")
		return
	}

	report, err := importer.New(s.repos).ImportCSV(r.Context(), file)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.repos.ImportReports.Save(r.Context(), report); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, report)
}

func (s *Server) getImportReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report, err := s.repos.ImportReports.GetLast(ctx)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, report)
}
