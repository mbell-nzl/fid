// Package linkrelation defines the retained link relation inverses shared by
// the link commands and journal admission.
package linkrelation

var inverses = map[string]string{
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

func Inverse(relation string) string { return inverses[relation] }

func Inverses() map[string]string {
	result := make(map[string]string, len(inverses))
	for relation, inverse := range inverses {
		result[relation] = inverse
	}
	return result
}
