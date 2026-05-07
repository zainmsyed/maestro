package api

import (
	"net/http"
	"time"

	"maestro/internal/models"
)

type auditResponse struct {
	ID         int64      `json:"id"`
	EntityType string     `json:"entity_type"`
	EntityID   string     `json:"entity_id"`
	ChangedBy  string     `json:"changed_by"`
	OldDate    *time.Time `json:"old_date"`
	NewDate    *time.Time `json:"new_date"`
	DeltaDays  int        `json:"delta_days"`
	Reason     *string    `json:"reason"`
	ChangedAt  time.Time  `json:"changed_at"`
	DateSource string     `json:"date_source,omitempty"`
}

func toAuditResponse(a models.DateAuditLog, dateSource string) auditResponse {
	return auditResponse{
		ID:         a.ID,
		EntityType: a.EntityType,
		EntityID:   a.EntityID,
		ChangedBy:  a.ChangedBy,
		OldDate:    a.OldDate,
		NewDate:    a.NewDate,
		DeltaDays:  a.DeltaDays,
		Reason:     a.Reason,
		ChangedAt:  a.ChangedAt,
		DateSource: dateSource,
	}
}

func (s *Server) listAudit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs, err := s.repos.Audits.List(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Batch-fetch date_source values to avoid N+1 queries
	featureIDs := make(map[string]bool)
	storyIDs := make(map[string]bool)
	for _, log := range logs {
		switch log.EntityType {
		case "feature":
			featureIDs[log.EntityID] = true
		case "story":
			storyIDs[log.EntityID] = true
		}
	}

	featureMap := make(map[string]string)
	for id := range featureIDs {
		f, err := s.repos.Features.GetByID(ctx, id)
		if err == nil {
			featureMap[id] = f.DateSource
		}
	}
	storyMap := make(map[string]string)
	for id := range storyIDs {
		st, err := s.repos.Stories.GetByID(ctx, id)
		if err == nil {
			storyMap[id] = st.DateSource
		}
	}

	resp := make([]auditResponse, len(logs))
	for i, log := range logs {
		var ds string
		switch log.EntityType {
		case "feature":
			ds = featureMap[log.EntityID]
		case "story":
			ds = storyMap[log.EntityID]
		}
		resp[i] = toAuditResponse(log, ds)
	}
	respondJSON(w, http.StatusOK, resp)
}
