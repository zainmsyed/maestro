package importer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sprintPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^sprint\s+\d+`),
	regexp.MustCompile(`(?i)^FY\d{2}\s+Q\d`),
	regexp.MustCompile(`(?i)^Q\d{1,2}$`),
	regexp.MustCompile(`(?i)^iteration\s+\d+`),
	regexp.MustCompile(`(?i)^\d{4}\s+Q\d`),
}

var unscheduledKeywords = []string{"backlog", "archive", "queue"}

func ParseIterationPath(path string) (string, bool) {
	segments := strings.Split(path, `\\`)

	for _, seg := range segments {
		for _, kw := range unscheduledKeywords {
			if strings.EqualFold(strings.TrimSpace(seg), kw) {
				return "", false
			}
		}
	}

	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		for _, pattern := range sprintPatterns {
			if pattern.MatchString(seg) {
				return seg, true
			}
		}
	}

	return "", false
}

func ParseAssignedTo(raw string) assignedTo {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return assignedTo{}
	}
	if left := strings.Index(raw, "<"); left >= 0 && strings.HasSuffix(raw, ">") {
		email := strings.TrimSuffix(strings.TrimSpace(raw[left+1:]), ">")
		name := strings.TrimSpace(raw[:left])
		if name == "" {
			name = email
		}
		return assignedTo{DisplayName: name, Email: cloneStringPtr(email)}
	}
	if left := strings.LastIndex(raw, "("); left >= 0 && strings.HasSuffix(raw, ")") {
		email := strings.TrimSuffix(strings.TrimSpace(raw[left+1:]), ")")
		name := strings.TrimSpace(raw[:left])
		if strings.Contains(email, "@") {
			return assignedTo{DisplayName: name, Email: cloneStringPtr(email)}
		}
	}
	if strings.Contains(raw, "@") && !strings.Contains(raw, " ") {
		return assignedTo{DisplayName: raw, Email: cloneStringPtr(raw)}
	}
	return assignedTo{DisplayName: raw}
}

func ParseStoryPoints(raw string) (*int, string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, ""
	}
	if f, err := strconv.ParseFloat(raw, 64); err == nil {
		rounded := int(f + 0.5)
		return &rounded, ""
	}
	return nil, fmt.Sprintf("non-numeric story points %q", raw)
}
