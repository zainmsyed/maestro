package importer

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var headerAliases = map[string]string{
	"parent":         "parent",
	"parent id":      "parent",
	"parentid":       "parent",
	"id":             "id",
	"work item id":   "id",
	"workitemid":     "id",
	"work item type": "work_item_type",
	"workitemtype":   "work_item_type",
	"type":           "work_item_type",
	"state":          "state",
	"assigned to":    "assigned_to",
	"assignedto":     "assigned_to",
	"assigned":       "assigned_to",
	"owner":          "assigned_to",
	"iteration path": "iteration_path",
	"iterationpath":  "iteration_path",
	"sprint":         "iteration_path",
	"iteration":      "iteration_path",
	"story points":   "story_points",
	"storypoints":    "story_points",
	"effort":         "story_points",
	"size":           "story_points",
	"points":         "story_points",
	"target date":    "target_date",
	"targetdate":     "target_date",
	"due date":       "target_date",
	"finish date":    "target_date",
	"duedate":        "target_date",
	"area path":      "area_path",
	"areapath":       "area_path",
	"area":           "area_path",
}

var requiredFields = []string{"parent", "id", "work_item_type"}

var titlePattern = regexp.MustCompile(`^title_(\d+)$`)

func NormalizeHeader(header string) string {
	normalized := strings.ToLower(strings.TrimSpace(header))
	normalized = strings.Join(strings.Fields(normalized), " ")
	compact := strings.NewReplacer("-", " ", "_", " ", ".", " ").Replace(normalized)
	compact = strings.Join(strings.Fields(compact), " ")
	if canonical, ok := headerAliases[compact]; ok {
		return canonical
	}
	return strings.ReplaceAll(compact, " ", "_")
}

func NormalizeHeaders(headers []string) (map[string]int, []string, error) {
	indices := make(map[string]int, len(headers))
	for idx, header := range headers {
		indices[NormalizeHeader(header)] = idx
	}
	if err := ValidateRequiredColumns(indices); err != nil {
		return nil, nil, err
	}
	titleCols := detectTitleColumns(indices)
	return indices, titleCols, nil
}

func ValidateRequiredColumns(indices map[string]int) error {
	missing := make([]string, 0)
	for _, field := range requiredFields {
		if _, ok := indices[field]; !ok {
			missing = append(missing, field)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	sort.Strings(missing)
	return fmt.Errorf("missing required columns: %s", strings.Join(missing, ", "))
}

func detectTitleColumns(indices map[string]int) []string {
	type colInfo struct {
		name   string
		number int
	}
	var cols []colInfo
	for key := range indices {
		if match := titlePattern.FindStringSubmatch(key); match != nil {
			if n, err := strconv.Atoi(match[1]); err == nil {
				cols = append(cols, colInfo{name: key, number: n})
			}
		}
	}
	sort.Slice(cols, func(i, j int) bool {
		return cols[i].number < cols[j].number
	})
	result := make([]string, len(cols))
	for i, c := range cols {
		result[i] = c.name
	}
	return result
}
