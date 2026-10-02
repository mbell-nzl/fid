package links

import (
	"slices"
	"strings"
)

var relationInverses = map[string]string{
	"blocks":          "blocked_by",
	"blocked_by":      "blocks",
	"relates_to":      "relates_to",
	"duplicates":      "duplicated_by",
	"duplicated_by":   "duplicates",
	"supersedes":      "superseded_by",
	"superseded_by":   "supersedes",
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

func endpointType(value, argument, operation string) (string, error) {
	endpointType := strings.TrimSpace(value)
	if endpointType == "" {
		return "", nil
	}
	if isV01EndpointType(endpointType) {
		return endpointType, nil
	}
	return "", &InvalidEndpointTypeError{Operation: operation, Argument: argument, Value: endpointType}
}

func relationAllowedForTypes(relation, fromType, toType string) bool {
	switch relation {
	case "relates_to":
		return fromType == "issue" && toType == "issue"
	case "blocks", "blocked_by", "duplicates", "duplicated_by", "supersedes", "superseded_by", "caused_by", "causes", "discovered_from", "discovered":
		return fromType == "issue" && toType == "issue"
	default:
		return false
	}
}

func isV01EndpointType(value string) bool {
	switch value {
	case "issue":
		return true
	default:
		return false
	}
}
