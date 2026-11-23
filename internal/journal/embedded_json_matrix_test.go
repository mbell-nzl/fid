package journal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// EmbeddedJSONProbeForTest crosses into the external differential test, which
// can import export without creating a journal -> export -> journal cycle.
type EmbeddedJSONProbeForTest struct {
	Name, Placement, Column, Reason, MigrationsDir string
	Variant, Raw                                   string
	Event                                          Event
}

// EmbeddedJSONMatrixForTest derives columns from SQLite, and table/operation
// placements from every retained family's catalog. These grammar-shaped
// carriers isolate JSON admission before family state, keys or SQL can mask it.
func EmbeddedJSONMatrixForTest(t *testing.T) []EmbeddedJSONProbeForTest {
	t.Helper()
	conn := openReplayDB(t, filepath.Join(t.TempDir(), "schema.sqlite"))
	defer closeProjectDB(t, conn)
	var probes []EmbeddedJSONProbeForTest
	for _, family := range slices.Sorted(maps.Keys(retainedCurrentRules)) {
		rule := retainedCurrentRules[family]
		seen := map[string]bool{}
		for _, stage := range rule.stages {
			for _, operation := range slices.Sorted(maps.Keys(stage.allowed)) {
				if seen[operation] {
					continue
				}
				seen[operation] = true
				table, op, _ := strings.Cut(operation, "/")
				rows, err := conn.Query(`PRAGMA table_info("` + table + `")`)
				if err != nil {
					t.Fatal(err)
				}
				var columns []string
				for rows.Next() {
					var index, notnull, pk int
					var name, kind string
					var fallback any
					if err := rows.Scan(&index, &name, &kind, &notnull, &fallback, &pk); err != nil {
						t.Fatal(err)
					}
					if strings.HasSuffix(name, "_json") {
						columns = append(columns, name)
					}
				}
				if err := rows.Err(); err != nil {
					t.Fatal(err)
				}
				if err := rows.Close(); err != nil {
					t.Fatal(err)
				}
				for _, column := range columns {
					// These two destinations are JSON maps/slices of scalars. The other
					// retained JSON columns have struct readers (binding_json is typed
					// for approvals; other person decisions don't consume a binding).
					variants := map[string]string{
						"exact":    `{"probe":1,"probe":2}`,
						"overflow": `{"probe":1e999}`,
						"depth":    `{"deep":` + strings.Repeat("[", 10_001) + "0" + strings.Repeat("]", 10_001) + "}",
					}
					if column != "names_json" && column != "request_field_digests_json" && (table != "person_decisions" || family == "task.approved") {
						variants["fold"] = `{"probe":1,"PROBE":2}`
					}
					for _, placement := range []string{"pk", "values", "prior", "preconditions"} {
						for _, variant := range slices.Sorted(maps.Keys(variants)) {
							e := admissionEvent("EVT-20261110-000004", family, rule.aggregate, "JSON-probe", 1)
							e.ActorV2.ActorType = strings.Fields(rule.actors)[0]
							// Import checks grammar before projection. Include every mandatory
							// stage and the target optional stage in their catalog order.
							target := -1
							for _, s := range rule.stages {
								choice := slices.Sorted(maps.Keys(s.allowed))[0]
								if s.allowed[operation] && target == -1 {
									choice = operation
									target = len(e.Effects)
								} else if s.minimum == 0 {
									continue
								}
								tab, action, _ := strings.Cut(choice, "/")
								e.Effects = append(e.Effects, Effect{Table: tab, Op: action, PK: map[string]any{"id": "JSON-probe"}, Values: map[string]any{}})
							}
							row := &e.Effects[target]
							if row.Table != table || row.Op != op {
								t.Fatal("matrix target drift")
							}
							maps := map[string]*map[string]any{"pk": &row.PK, "values": &row.Values, "prior": &row.Prior, "preconditions": &row.Preconditions}
							*maps[placement] = map[string]any{column: variants[variant]}
							reason := "duplicate key"
							switch variant {
							case "overflow":
								reason = "cannot unmarshal number"
							case "depth":
								reason = "exceeded max depth"
							}
							probes = append(probes, EmbeddedJSONProbeForTest{Name: family + "/" + operation + "/" + column + "/" + placement + "/" + variant, Placement: placement, Column: column, Reason: reason, Variant: variant, Raw: variants[variant], Event: e})
						}
					}
				}
			}
		}
	}
	return probes
}

// CheckEmbeddedJSONPathsForTest checks bytes, all projected tables, refusal
// classification, and recovery through a valid write on every journal path.
func CheckEmbeddedJSONPathsForTest(t *testing.T, probe EmbeddedJSONProbeForTest) {
	t.Helper()
	CheckEmbeddedJSONMatrixPathsForTest(t, []EmbeddedJSONProbeForTest{probe})
}

// Reuse a bounded valid prefix per path. Every refusal is followed by a real
// commit, so the next probe also proves an advancing ledger remains usable.
// Batching avoids thousands of redundant git init/migration operations.
func CheckEmbeddedJSONMatrixPathsForTest(t *testing.T, probes []EmbeddedJSONProbeForTest) {
	t.Helper()
	checkEmbeddedJSONModesForTest(t, probes, admissionModes(admissionLaneRaise()))
}

func checkEmbeddedJSONModesForTest(t *testing.T, probes []EmbeddedJSONProbeForTest, modes []string) {
	t.Helper()
	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			modeTest := t
			var w jsonAdmissionWorkspace
			count := 0
			for _, probe := range probes {
				t.Run(probe.Name, func(t *testing.T) {
					n := count % 16
					if n == 0 {
						if w.conn != nil {
							closeProjectDB(t, w.conn)
						}
						migrations := probe.MigrationsDir
						if migrations == "" {
							migrations = testProjectMigrationsDir(modeTest)
						}
						w = newJSONAdmissionWorkspace(modeTest, mode, migrations)
					}
					count++
					probe.Event.EventID = fmt.Sprintf("EVT-20261110-%06d", n+4)
					before := w.snapshot(t)
					err := w.admit(t, mode, probe.Event)
					reason := probe.Reason
					if reason == "" {
						reason = "duplicate key"
					}
					if err == nil || !strings.Contains(err.Error(), reason) {
						t.Fatalf("embedded JSON %s in %s = %v, want %s refusal", probe.Column, probe.Placement, err, reason)
					}
					if probe.Variant != "" {
						checkEmbeddedJSONRowRefusal(t, probe, err)
					}
					live := mode == "commit" || mode == "dry-run" || mode == "project"
					if !live {
						var corrupt *CorruptJournalError
						if !errors.As(err, &corrupt) {
							t.Fatalf("replay = %v, want CorruptJournalError", err)
						}
					}
					if !reflect.DeepEqual(before, w.snapshot(t)) {
						t.Fatal("refusal changed projection")
					}
					want := w.prefix
					if !live {
						want = w.journalWith(t, probe.Event)
					}
					got, err := os.ReadFile(w.file)
					if err != nil || !bytes.Equal(got, want) {
						t.Fatalf("refusal changed journal: %v", err)
					}
					if err := os.WriteFile(w.file, w.prefix, 0o644); err != nil {
						t.Fatal(err)
					}
					if _, err := EnsureProjectionFresh(w.root, w.migrations); err != nil {
						t.Fatal(err)
					}
					valid := admissionLaneRaise()
					valid.EventID = probe.Event.EventID
					valid.Aggregate.Sequence = n + 2
					valid.Effects[0].Preconditions = nil
					valid.Effects[1].Values["revision"] = n + 2
					valid.Effects[1].Prior["revision"] = n + 1
					valid.Effects[1].Preconditions["revision"] = n + 1
					valid.Effects[2].PK["id"] = fmt.Sprintf("ACT-valid-%d", n)
					valid.Effects[2].Values["journal_event_id"] = valid.EventID
					if err := w.admit(t, "commit", valid); err != nil {
						t.Fatalf("valid follow-up: %v", err)
					}
					w.prefix, err = os.ReadFile(w.file)
					if err != nil {
						t.Fatal(err)
					}
				})
				if t.Failed() {
					return
				} // A failed probe must not spoil later fixtures.
			}
		})
	}
}

func TestEmbeddedJSONMatrixCoverage(t *testing.T) {
	t.Parallel()
	probes := EmbeddedJSONMatrixForTest(t)
	counts := map[string]int{}
	for _, p := range probes {
		counts[p.Placement]++
	}
	for _, placement := range []string{"pk", "values", "prior", "preconditions"} {
		if counts[placement] == 0 {
			t.Fatal("empty placement: " + placement)
		}
	}
	t.Logf("generated %d probes: %v", len(probes), counts)
}

// CheckEmbeddedJSONReadersForTest keeps the serialized Event and direct stored
// readers in the differential suite as well: an already-current projection
// does not run admission before rendering a retry.
func CheckEmbeddedJSONReadersForTest(t *testing.T) {
	t.Helper()
	t.Run("live-serialized", func(t *testing.T) {
		t.Parallel()
		e := admissionLaneRaise()
		e.Effects[2].Values["details"] = json.RawMessage(`{"probe":1,"probe":2}`)
		// Canonical journal encoding collapses raw duplicate keys. Offer the
		// unnormalized producer value to live admission; literal stored-key
		// probes below exercise each replay path without canonicalizing them.
		checkEmbeddedJSONModesForTest(t, []EmbeddedJSONProbeForTest{{Event: e, Column: "details", Placement: "values"}}, []string{"commit", "dry-run", "project"})
	})
	t.Run("live-numbers", func(t *testing.T) {
		t.Parallel()
		e := admissionLaneRaise()
		e.Effects[2].Values["details"] = json.RawMessage(`{"probe":1e999}`)
		CheckEmbeddedJSONPathsForTest(t, EmbeddedJSONProbeForTest{Event: e, Column: "details", Placement: "values", Reason: "cannot unmarshal number"})
	})
	t.Run("stored", TestStoredRetryJSONKeys)
	t.Run("serialized", TestMalformedEventDuplicateJSONKeys)
	t.Run("typed", TestTypedJSONKeysFollowDestination)
	t.Run("policy", TestPolicyPayloadStructAliases)
	t.Run("retry-binding", func(t *testing.T) {
		t.Parallel()
		for _, column := range []string{"result_json", "request_field_digests_json"} {
			t.Run(column, func(t *testing.T) {
				t.Parallel()
				event := admissionLaneRaise()
				request := IdempotencyRequest{OperationScope: "issue.update:project-1", IdempotencyKey: "reader", RequestDigest: "sha256:" + strings.Repeat("b", 64)}
				event.Idempotency = &EventIdempotency{Key: request.IdempotencyKey, Scope: request.OperationScope, RequestDigest: request.RequestDigest}
				row := idempotencyResultEffect(request, event, event.Aggregate.ID, `{"replayed":false}`)
				if column == "result_json" {
					row.Values[column] = `{"replayed":true,"replayed":false}`
				} else {
					row.Values[column] = `{"x":1,"x":2}`
				}
				event.Effects = append(event.Effects, row)
				if _, err := validateIdempotencyBinding(event); err == nil || !strings.Contains(err.Error(), "duplicate key") {
					t.Fatalf("retry reader = %v, want duplicate-key refusal", err)
				}
			})
		}
	})
}

func TestEmbeddedJSONDestinationMaps(t *testing.T) {
	t.Parallel()
	// Unknown nested response members remain untyped. A column-wide superset
	// of every retry shape would accidentally fold maps in other families.
	for _, family := range []string{"issue.created", "issue.updated", "issue.transitioned", "issue.ac-added", "issue.label.added", "label.added", "link.added", "hold.parked", "artifact.availability-observed"} {
		t.Run(family, func(t *testing.T) {
			t.Parallel()
			typed := map[string]string{
				"issue.created": "issue", "issue.updated": "issue changes", "issue.transitioned": "issue", "issue.ac-added": "acceptance_criterion", "label.added": "label", "link.added": "link inverse", "artifact.availability-observed": "artifact",
			}[family]
			raw := `{"extra":{"A":1,"a":2},"fields":{"build":1,"Build":2},"issue_revisions":{"s":1,"ſ":2}}`
			for _, field := range strings.Fields("issue changes acceptance_criterion link inverse label artifact") {
				if strings.Contains(" "+typed+" ", " "+field+" ") {
					continue
				}
				value := `{"A":1,"a":2}`
				if field == "changes" {
					value = "[" + value + "]"
				}
				raw = strings.TrimSuffix(raw, "}") + `,"` + field + `":` + value + "}"
			}

			for _, placement := range []string{"pk", "values", "prior", "preconditions"} {
				t.Run(placement, func(t *testing.T) {
					row := Effect{Table: "idempotency_results"}
					target := map[string]*map[string]any{"pk": &row.PK, "values": &row.Values, "prior": &row.Prior, "preconditions": &row.Preconditions}
					*target[placement] = map[string]any{"result_json": raw}
					e := admissionEvent("EVT-20261110-000004", family, "issue", "FID-a", 2, row)
					if err := validateEmbeddedJSON(e); err != nil {
						t.Fatalf("case-sensitive response maps: %v", err)
					}
				})
			}
		})
	}
}
