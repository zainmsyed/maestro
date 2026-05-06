package importer

import "strings"

func NormalizeWorkItemType(raw string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "epic":
		return "epic", true
	case "feature":
		return "feature", true
	case "user story", "product backlog item", "requirement":
		return "story", true
	default:
		return "", false
	}
}
