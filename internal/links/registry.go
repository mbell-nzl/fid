package links

import (
	"regexp"
	"slices"
	"strings"
)

var requirementIDPattern = regexp.MustCompile(`^(PRD|TRD)-[A-Z0-9]+-[0-9]+[a-z]?$`)

var relationInverses = map[string]string{
	"blocks":          "blocked_by",
	"blocked_by":      "blocks",
	"relates_to":      "relates_to",
	"duplicates":      "duplicated_by",
	"duplicated_by":   "duplicates",
	"supersedes":      "superseded_by",
	"superseded_by":   "supersedes",
	"implements":      "implemented_by",
	"implemented_by":  "implements",
	"caused_by":       "causes",
	"causes":          "caused_by",
	"discovered_from": "discovered",
	"discovered":      "discovered_from",
}

var relationNames = []string{
	"blocks",
	"blocked_by",
	"relates_to",
	"duplicates",
	"duplicated_by",
	"supersedes",
	"superseded_by",
	"implements",
	"implemented_by",
	"caused_by",
	"causes",
	"discovered_from",
	"discovered",
}

func RelationInverseMap() map[string]string {
	cloned := make(map[string]string, len(relationInverses))
	for relation, inverse := range relationInverses {
		cloned[relation] = inverse
	}
	return cloned
}

func SupportedRelations() []string {
	return slices.Clone(relationNames)
}

func inverseRelation(relation string) (string, bool) {
	inverse, ok := relationInverses[strings.TrimSpace(relation)]
	return inverse, ok
}

func endpointTypeForID(id string) string {
	if requirementIDPattern.MatchString(strings.TrimSpace(id)) {
		return "requirement"
	}
	return "issue"
}

func relationAllowedForTypes(relation, fromType, toType string) bool {
	switch relation {
	case "relates_to":
		return isV01EndpointType(fromType) && isV01EndpointType(toType)
	case "implements":
		return fromType == "issue" && toType == "requirement"
	case "implemented_by":
		return fromType == "requirement" && toType == "issue"
	case "blocks", "blocked_by", "duplicates", "duplicated_by", "supersedes", "superseded_by", "caused_by", "causes", "discovered_from", "discovered":
		return fromType == "issue" && toType == "issue"
	default:
		return false
	}
}

func isV01EndpointType(value string) bool {
	switch value {
	case "issue", "requirement":
		return true
	default:
		return false
	}
}
