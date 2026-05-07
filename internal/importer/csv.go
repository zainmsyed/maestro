package importer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"maestro/internal/models"
)

func (i *Importer) ImportCSV(ctx context.Context, reader io.Reader) (*ImportReport, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("csv is empty")
	}

	headerIndex, titleCols, err := NormalizeHeaders(records[0])
	if err != nil {
		return nil, err
	}

	state := newImportState()
	entities := make(map[string]*rawEntity)

	// Pass 1: Parse all rows into entity map (last occurrence of each ID wins)
	for rowIndex, record := range records[1:] {
		rowNumber := rowIndex + 2
		id := normalizeID(field(record, headerIndex, "id"))
		itemType, ok := NormalizeWorkItemType(field(record, headerIndex, "work_item_type"))
		if !ok {
			state.report.SkippedRows++
			state.report.Warnings = append(state.report.Warnings,
				fmt.Sprintf("row %d: unsupported work item type %q", rowNumber, field(record, headerIndex, "work_item_type")))
			continue
		}
		if id == "" {
			id = state.nextSyntheticID(itemType)
			if itemType == "story" {
				state.report.SyntheticStoryIDs = append(state.report.SyntheticStoryIDs, id)
			}
		}
		title := extractTitle(record, headerIndex, titleCols)
		if title == "" {
			title = fmt.Sprintf("Untitled %s %s", itemType, id)
		}
		parentID := normalizeID(field(record, headerIndex, "parent"))

		assigned := ParseAssignedTo(field(record, headerIndex, "assigned_to"))
		sprint, isScheduled := ParseIterationPath(field(record, headerIndex, "iteration_path"))
		if !isScheduled {
			state.report.MissingSprintCount++
		} else {
			state.addSprint(sprint)
		}

		targetDate, err := ParseTargetDate(field(record, headerIndex, "target_date"))
		if err != nil {
			state.report.SkippedRows++
			state.report.Warnings = append(state.report.Warnings, fmt.Sprintf("row %d: %v", rowNumber, err))
			continue
		}
		if targetDate.Time == nil {
			state.report.MissingDatesCount++
			state.report.DateAssignmentCandidates = append(state.report.DateAssignmentCandidates, DateAssignmentCandidate{
				RowNumber:     rowNumber,
				WorkItemType:  itemType,
				ID:            id,
				Title:         title,
				AssignedOwner: assigned.DisplayName,
			})
		} else {
			state.dateFormats[targetDate.Format]++
			if targetDate.Ambiguous {
				state.report.AmbiguousDates = append(state.report.AmbiguousDates, AmbiguousDateCandidate{
					RowNumber:    rowNumber,
					WorkItemType: itemType,
					ID:           id,
					Title:        title,
					RawDate:      strings.TrimSpace(field(record, headerIndex, "target_date")),
					ParsedDate:   *targetDate.Time,
				})
			}
		}

		storyPoints, warning := ParseStoryPoints(field(record, headerIndex, "story_points"))
		if warning != "" {
			state.report.Warnings = append(state.report.Warnings, fmt.Sprintf("row %d: %s", rowNumber, warning))
		}

		entities[id] = &rawEntity{
			id:          id,
			parentID:    parentID,
			itemType:    itemType,
			title:       title,
			state:       field(record, headerIndex, "state"),
			owner:       assigned.DisplayName,
			sprint:      sprint,
			storyPoints: storyPoints,
			targetDate:  targetDate.Time,
			rowNumber:   rowNumber,
		}
	}

	// Pass 2: Link and persist in type order (epic → feature → story)
	for _, e := range entities {
		if e.itemType == "epic" {
			if err := state.createEpic(ctx, i, e); err != nil {
				return nil, err
			}
		}
	}
	for _, e := range entities {
		if e.itemType == "feature" {
			if err := state.createFeature(ctx, i, e, entities); err != nil {
				return nil, err
			}
		}
	}
	for _, e := range entities {
		if e.itemType == "story" {
			if err := state.createStory(ctx, i, e, entities); err != nil {
				return nil, err
			}
		}
	}

	return state.finalizeReport(), nil
}

func field(record []string, headerIndex map[string]int, name string) string {
	idx, ok := headerIndex[name]
	if !ok || idx >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[idx])
}

func normalizeID(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.Contains(raw, ".") {
		raw = strings.Split(raw, ".")[0]
	}
	return raw
}

func extractTitle(record []string, headerIndex map[string]int, titleCols []string) string {
	last := ""
	for _, col := range titleCols {
		if val := strings.TrimSpace(field(record, headerIndex, col)); val != "" {
			last = val
		}
	}
	return last
}

func (s *importState) ensureUnassignedEpic(ctx context.Context, importer *Importer) (string, error) {
	if s.createdEpics[syntheticUnassignedEpicID] {
		return syntheticUnassignedEpicID, nil
	}
	if _, err := importer.repos.Epics.GetByID(ctx, syntheticUnassignedEpicID); err == nil {
		s.createdEpics[syntheticUnassignedEpicID] = true
		s.report.ExistingSkipped++
		return syntheticUnassignedEpicID, nil
	}
	epic := &models.Epic{
		ID:          syntheticUnassignedEpicID,
		Title:       "Unassigned Epic",
		Description: "Synthetic epic for orphaned imported records",
		Status:      "Imported",
		IsSynthetic: true,
	}
	if err := importer.repos.Epics.Create(ctx, epic); err != nil {
		return "", fmt.Errorf("create synthetic unassigned epic: %w", err)
	}
	s.createdEpics[syntheticUnassignedEpicID] = true
	s.report.EpicCount++
	return syntheticUnassignedEpicID, nil
}

func (s *importState) ensureUnassignedFeature(ctx context.Context, importer *Importer) (string, error) {
	if s.createdFeatures[syntheticUnassignedFeatureID] {
		return syntheticUnassignedFeatureID, nil
	}
	if _, err := importer.repos.Features.GetByID(ctx, syntheticUnassignedFeatureID); err == nil {
		s.createdFeatures[syntheticUnassignedFeatureID] = true
		s.report.ExistingSkipped++
		return syntheticUnassignedFeatureID, nil
	}
	epicID, err := s.ensureUnassignedEpic(ctx, importer)
	if err != nil {
		return "", err
	}
	feature := &models.Feature{
		ID:          syntheticUnassignedFeatureID,
		EpicID:      cloneStringPtr(epicID),
		Title:       "Unassigned Feature",
		Description: "Synthetic feature for orphaned imported stories",
		Status:      "Imported",
		DateSource:  "imported",
	}
	if err := importer.repos.Features.Create(ctx, feature); err != nil {
		return "", fmt.Errorf("create synthetic unassigned feature: %w", err)
	}
	s.createdFeatures[syntheticUnassignedFeatureID] = true
	s.report.FeatureCount++
	return syntheticUnassignedFeatureID, nil
}

func (s *importState) createEpic(ctx context.Context, importer *Importer, e *rawEntity) error {
	if s.createdEpics[e.id] {
		return nil
	}
	if _, err := importer.repos.Epics.GetByID(ctx, e.id); err == nil {
		s.createdEpics[e.id] = true
		s.report.ExistingSkipped++
		return nil
	}
	original, committed := withDateLocked(e.targetDate)
	epic := &models.Epic{
		ID:               e.id,
		Title:            e.title,
		Description:      "",
		Status:           e.state,
		Owner:            e.owner,
		SprintEnd:        e.sprint,
		OriginalEndDate:  original,
		CommittedEndDate: committed,
		IsSynthetic:      false,
	}
	if err := importer.repos.Epics.Create(ctx, epic); err != nil {
		return fmt.Errorf("create epic %s: %w", e.id, err)
	}
	s.createdEpics[e.id] = true
	s.report.EpicCount++
	return nil
}

func (s *importState) createFeature(ctx context.Context, importer *Importer, e *rawEntity, entities map[string]*rawEntity) error {
	if s.createdFeatures[e.id] {
		return nil
	}
	if _, err := importer.repos.Features.GetByID(ctx, e.id); err == nil {
		s.createdFeatures[e.id] = true
		s.report.ExistingSkipped++
		return nil
	}
	parentEpicID := ""
	if e.parentID != "" {
		if parent, ok := entities[e.parentID]; ok && parent.itemType == "epic" {
			parentEpicID = e.parentID
		}
	}
	if parentEpicID == "" {
		unassignedEpicID, err := s.ensureUnassignedEpic(ctx, importer)
		if err != nil {
			return err
		}
		parentEpicID = unassignedEpicID
		s.report.OrphanedFeatures++
	}
	original, committed := withDateLocked(e.targetDate)
	feature := &models.Feature{
		ID:               e.id,
		EpicID:           cloneStringPtr(parentEpicID),
		Title:            e.title,
		Description:      "",
		Status:           e.state,
		Owner:            e.owner,
		Sprint:           e.sprint,
		OriginalEndDate:  original,
		CommittedEndDate: committed,
		StoryPoints:      e.storyPoints,
		DateSource:       "imported",
	}
	if err := importer.repos.Features.Create(ctx, feature); err != nil {
		return fmt.Errorf("create feature %s: %w", e.id, err)
	}
	s.createdFeatures[e.id] = true
	s.report.FeatureCount++
	return nil
}

func (s *importState) createStory(ctx context.Context, importer *Importer, e *rawEntity, entities map[string]*rawEntity) error {
	if s.createdStories[e.id] {
		return nil
	}
	if _, err := importer.repos.Stories.GetByID(ctx, e.id); err == nil {
		s.createdStories[e.id] = true
		s.report.ExistingSkipped++
		return nil
	}
	parentFeatureID := ""
	if e.parentID != "" {
		if parent, ok := entities[e.parentID]; ok && parent.itemType == "feature" {
			parentFeatureID = e.parentID
		}
	}
	if parentFeatureID == "" {
		unassignedFeatureID, err := s.ensureUnassignedFeature(ctx, importer)
		if err != nil {
			return err
		}
		parentFeatureID = unassignedFeatureID
		s.report.OrphanedStories++
	}
	original, committed := withDateLocked(e.targetDate)
	story := &models.Story{
		ID:               e.id,
		FeatureID:        parentFeatureID,
		Title:            e.title,
		Description:      "",
		Status:           e.state,
		Owner:            e.owner,
		Sprint:           e.sprint,
		StoryPoints:      e.storyPoints,
		OriginalEndDate:  original,
		CommittedEndDate: committed,
		DateSource:       "imported",
	}
	if err := importer.repos.Stories.Create(ctx, story); err != nil {
		return fmt.Errorf("create story %s: %w", e.id, err)
	}
	s.createdStories[e.id] = true
	s.report.StoryCount++
	return nil
}
