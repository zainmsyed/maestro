package importer

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ParseTargetDate(raw string) (parsedDate, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return parsedDate{}, nil
	}

	for _, layout := range []struct {
		layout string
		name   string
	}{
		{"2006-01-02T15:04:05Z07:00", "ISO with timezone"},
		{"2006-01-02", "ISO 8601"},
		{"2006/01/02", "ISO slashes"},
		{"01/02/2006 3:04:05 PM", "US datetime"},
		{"01/02/2006 15:04:05", "US datetime 24h"},
		{"January 2, 2006", "long form"},
		{"2-Jan-2006", "abbreviated"},
	} {
		if parsed, err := time.Parse(layout.layout, raw); err == nil {
			normalized := normalizeDate(parsed)
			return parsedDate{Time: &normalized, Format: layout.name}, nil
		}
	}

	if strings.Count(raw, "/") == 2 {
		parts := strings.Split(raw, "/")
		if len(parts) == 3 {
			first, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			second, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
			year, err3 := strconv.Atoi(strings.TrimSpace(parts[2]))
			if err1 == nil && err2 == nil && err3 == nil {
				ambiguous := first <= 12 && second <= 12
				layout := "01/02/2006"
				formatName := "MM/DD/YYYY"
				if first > 12 {
					layout = "02/01/2006"
					formatName = "DD/MM/YYYY"
				}
				if parsed, err := time.Parse(layout, fmt.Sprintf("%02d/%02d/%04d", first, second, year)); err == nil {
					normalized := normalizeDate(parsed)
					return parsedDate{Time: &normalized, Format: formatName, Ambiguous: ambiguous}, nil
				}
			}
		}
	}

	return parsedDate{}, fmt.Errorf("unsupported target date format: %q", raw)
}

func normalizeDate(value time.Time) time.Time {
	utc := value.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}

func summarizeDateFormats(counts map[string]int) string {
	if len(counts) == 0 {
		return ""
	}
	if len(counts) == 1 {
		for format := range counts {
			return format
		}
	}
	return "mixed"
}

func sortedKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
