package journal

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// The catalog/schema matrix checks every destination at the shared boundary,
// before family rules. Only its representative set needs repository journeys.
func TestEmbeddedJSONValidators(t *testing.T) {
	t.Parallel()
	for _, probe := range EmbeddedJSONMatrixForTest(t) {
		t.Run(probe.Name, func(t *testing.T) {
			t.Run("shape", func(t *testing.T) {
				checkEmbeddedJSONRowRefusal(t, probe, ValidateEventShape(probe.Event))
			})
			t.Run("decode", func(t *testing.T) {
				raw, err := json.Marshal(probe.Event)
				if err != nil {
					t.Fatal(err)
				}
				var decoded Event
				checkEmbeddedJSONRowRefusal(t, probe, decoded.UnmarshalJSON(raw))
			})
			t.Run("import-text", func(t *testing.T) {
				// Import's schema validator consumes untyped JSON. Struct aliases
				// are refused by Event decoding, before that schema check runs;
				// case-distinct map keys remain valid at the untyped boundary.
				err := ValidateJSONText(probe.Raw)
				if probe.Variant == "fold" {
					if err != nil {
						t.Fatalf("distinct map keys refused: %v", err)
					}
				} else if err == nil || !strings.Contains(err.Error(), probe.Reason) {
					t.Fatalf("import JSON validation = %v, want %s refusal", err, probe.Reason)
				}
			})
		})
	}
}

func checkEmbeddedJSONRowRefusal(t *testing.T, probe EmbeddedJSONProbeForTest, err error) {
	t.Helper()
	var refused *RowIdentityError
	if !errors.As(err, &refused) || refused.EventID != probe.Event.EventID ||
		!strings.Contains(refused.Reason, probe.Placement+"."+probe.Column) || !strings.Contains(refused.Reason, probe.Reason) {
		t.Fatalf("embedded JSON validation = %v, want RowIdentityError for %s.%s: %s", err, probe.Placement, probe.Column, probe.Reason)
	}
}

// EmbeddedJSONRepresentativesForTest selects a stable carrier for each map and
// malformed variant. Every path must still reach the shared validator before a
// family-specific rule can mask its refusal.
func EmbeddedJSONRepresentativesForTest(probes []EmbeddedJSONProbeForTest) []EmbeddedJSONProbeForTest {
	seen := map[string]bool{}
	var representatives []EmbeddedJSONProbeForTest
	for _, probe := range probes {
		key := probe.Placement + "/" + probe.Variant
		if !seen[key] {
			seen[key] = true
			representatives = append(representatives, probe)
		}
	}
	return representatives
}

func TestEmbeddedJSONRepresentatives(t *testing.T) {
	t.Parallel()
	probes := EmbeddedJSONMatrixForTest(t)
	representatives := EmbeddedJSONRepresentativesForTest(probes)
	counts := map[string]int{}
	for _, probe := range representatives {
		counts[probe.Placement+"/"+probe.Variant]++
	}
	for _, placement := range []string{"pk", "values", "prior", "preconditions"} {
		for _, variant := range []string{"exact", "fold", "overflow", "depth"} {
			key := placement + "/" + variant
			if counts[key] != 1 {
				t.Errorf("representative %s count = %d, want 1", key, counts[key])
			}
		}
	}
	if len(representatives) != 16 {
		t.Fatalf("representative count = %d, want 16", len(representatives))
	}
	t.Logf("%d generated probes, %d path representatives", len(probes), len(representatives))
}
