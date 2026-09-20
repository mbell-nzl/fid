package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCurrentEventFamiliesRejectMalformedEffectsLiveAndReplay(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Event)
	}{
		{"relabeled-recovery", func(e *Event) { e.EventType = "run.recovered"; e.CommandName = "run.recover" }},
		{"relabeled-budget", func(e *Event) { e.EventType = "run.budget-observed"; e.CommandName = "run.observe-budget" }},
		{"relabeled-blocked", func(e *Event) { e.EventType = "run.execution-blocked"; e.CommandName = "gate.run" }},
		{"unknown-family", func(e *Event) { e.EventType = "unknown.recovery-event" }},
		{"missing-observation", func(e *Event) { e.Effects = e.Effects[1:] }},
		{"missing-runtime", func(e *Event) { e.Effects = e.Effects[:1] }},
		{"reordered", func(e *Event) { e.Effects[0], e.Effects[1] = e.Effects[1], e.Effects[0] }},
		{"wrong-aggregate", func(e *Event) { e.Aggregate.ID = "RUN-unrelated" }},
		{"wrong-aggregate-type", func(e *Event) { e.Aggregate.Type = "activation" }},
		{"wrong-command", func(e *Event) { e.CommandName = "run.close" }},
		{"wrong-observation-kind", func(e *Event) { e.Effects[0].Values["kind"] = "liveness" }},
		{"wrong-created-event", func(e *Event) { e.Effects[0].Values["created_event_id"] = "EVT-other" }},
		{"coerced-text", func(e *Event) { e.Effects[0].Values["actor_id"] = 42 }},
		{"missing-version-precondition", func(e *Event) { e.Effects[1].Preconditions = nil }},
		{"version-gap", func(e *Event) { e.Effects[1].Values["version"] = 4 }},
		{"stale-version", func(e *Event) { e.Effects[1].Preconditions["version"] = 99 }},
		{"duplicate-immutable", func(e *Event) { e.Effects = append(e.Effects[:1], append([]Effect{e.Effects[0]}, e.Effects[1:]...)...) }},
		{"unknown-field", func(e *Event) { e.Effects[0].Values["surprise"] = true }},
		{"bad-reference", func(e *Event) { e.Effects[0].Values["run_id"] = "RUN-missing" }},
		{"missing-idempotency", func(e *Event) { e.Effects = e.Effects[:len(e.Effects)-1] }},
		{"wrong-result-event", func(e *Event) { e.Effects[len(e.Effects)-1].Values["result_event_id"] = "EVT-wrong" }},
		{"wrong-idempotency-digest", func(e *Event) { e.Effects[len(e.Effects)-1].Values["request_digest"] = "sha256:" }},
		{"wrong-result-entity", func(e *Event) { e.Effects[len(e.Effects)-1].Values["result_entity_id"] = "OBS-wrong" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, mode := range []string{"live", "replay"} {
				t.Run(mode, func(t *testing.T) {
					event := checkpointedProjectionEvent()
					tc.mutate(&event)
					conn := openReplayDB(t, filepath.Join(t.TempDir(), "projection.sqlite"))
					defer closeProjectDB(t, conn)
					seedCheckpointProjectionRun(t, conn)
					before := currentProjectionSnapshot(t, conn)
					var err error
					if mode == "live" {
						var hash string
						hash, err = HashEvent(event)
						if err != nil {
							t.Fatal(err)
						}
						err = ProjectEvent(conn, ProjectionInput{Event: event, EventHash: hash, File: "2026-09-06.jsonl", Line: 1, Offset: 1})
					} else {
						workspace := newReplayWorkspace(t)
						writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), event)
						_, err = ReplayJournal(conn, filepath.Join(workspace, ".fid", "journal"), nil)
					}
					if err == nil {
						t.Fatal("malformed current-family Event accepted")
					}
					if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
						t.Fatal("rejected Event partially projected")
					}
				})
			}
		})
	}
}

func currentProjectionSnapshot(t *testing.T, conn *sql.DB) map[string][]string {
	t.Helper()
	names, err := conn.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		t.Fatal(err)
	}
	var tables []string
	for names.Next() {
		var name string
		if err := names.Scan(&name); err != nil {
			t.Fatal(err)
		}
		tables = append(tables, name)
	}
	if err := names.Err(); err != nil {
		t.Fatal(err)
	}
	if err := names.Close(); err != nil {
		t.Fatal(err)
	}
	result := map[string][]string{}
	for _, table := range tables {
		rows, err := conn.Query(`SELECT * FROM "` + table + `" ORDER BY rowid`)
		if err != nil {
			t.Fatal(err)
		}
		columns, err := rows.Columns()
		if err != nil {
			t.Fatal(err)
		}
		result[table] = []string{}
		for rows.Next() {
			values := make([]any, len(columns))
			pointers := make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			if err := rows.Scan(pointers...); err != nil {
				t.Fatal(err)
			}
			raw, err := json.Marshal(values)
			if err != nil {
				t.Fatal(err)
			}
			result[table] = append(result[table], string(raw))
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
	}
	return result
}

func TestCurrentEventRequiresExactCapturedPrior(t *testing.T) {
	for _, kind := range []string{"valid", "missing", "wrong", "extra"} {
		t.Run(kind, func(t *testing.T) {
			event := checkpointedProjectionEvent()
			event.Effects[1].Prior = map[string]any{"version": 1}
			switch kind {
			case "missing":
				event.Effects[1].Prior = nil
			case "wrong":
				event.Effects[1].Prior["version"] = 99
			case "extra":
				event.Effects[1].Prior["executor_id"] = "executor-checkpoint"
			}
			conn := openReplayDB(t, filepath.Join(t.TempDir(), "prior.sqlite"))
			defer closeProjectDB(t, conn)
			seedCheckpointProjectionRun(t, conn)
			before := currentProjectionSnapshot(t, conn)
			hash, err := HashEvent(event)
			if err != nil {
				t.Fatal(err)
			}
			err = ProjectEvent(conn, ProjectionInput{Event: event, EventHash: hash})
			if kind == "valid" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil {
				t.Fatal("invalid captured prior accepted")
			}
			if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
				t.Fatal("invalid prior partially projected")
			}
		})
	}
}

func TestCurrentEventRejectsConditionalBindingCorruption(t *testing.T) {
	cases := []struct {
		name   string
		event  func() Event
		seed   func(*testing.T, *sql.DB)
		mutate func(*Event)
	}{
		{"additional-observation", checkpointedProjectionEvent, seedCheckpointProjectionRun, func(e *Event) {
			raw, _ := json.Marshal(e.Effects[0])
			var extra Effect
			_ = json.Unmarshal(raw, &extra)
			extra.PK["id"] = "OBS-extra"
			extra.Values["observation_json"] = `{"id":"OBS-extra","kind":"checkpoint","run_id":"RUN-checkpoint"}`
			e.Effects = append(e.Effects[:1], append([]Effect{extra}, e.Effects[1:]...)...)
		}},
		{"missing-attention-link", executionBlockedProjectionEvent, seedExecutionBlockedProjection, func(e *Event) { e.Effects = append(e.Effects[:2], e.Effects[3:]...) }},
		{"mismatched-attention", executionBlockedProjectionEvent, seedExecutionBlockedProjection, func(e *Event) { e.Effects[1].Values["code"] = "operator-stop" }},
		{"foreign-recovery-aggregate", recoveredProjectionEvent, seedRecoveredProjection, func(e *Event) { e.Aggregate.ID = "RUN-successor" }},
		{"missing-selected-stop", func() Event { return activationRecoveryProjectionEvent("activation.recovery-selected") }, func(t *testing.T, c *sql.DB) { seedActivationRecoveryProjection(t, c, "activation.recovery-selected") }, func(e *Event) { e.Effects = append(e.Effects[:2], e.Effects[len(e.Effects)-1]) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, mode := range []string{"live", "replay"} {
				t.Run(mode, func(t *testing.T) {
					c := openReplayDB(t, filepath.Join(t.TempDir(), "projection.sqlite"))
					defer closeProjectDB(t, c)
					tc.seed(t, c)
					before := currentProjectionSnapshot(t, c)
					original := tc.event()
					if err := DryRunProjectEvent(c, original); err != nil {
						t.Fatalf("valid control: %v", err)
					}
					tc.mutate(&original)
					var err error
					if mode == "live" {
						hash, hashErr := HashEvent(original)
						if hashErr != nil {
							t.Fatal(hashErr)
						}
						err = ProjectEvent(c, ProjectionInput{Event: original, EventHash: hash})
					} else {
						workspace := newReplayWorkspace(t)
						writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), original)
						_, err = ReplayJournal(c, filepath.Join(workspace, ".fid", "journal"), nil)
					}
					if err == nil {
						t.Fatal("accepted corrupt conditional binding")
					}
					if !reflect.DeepEqual(before, currentProjectionSnapshot(t, c)) {
						t.Fatal("rejection changed projection")
					}
				})
			}
		})
	}
}

func TestCurrentEventColumnDomainsMatchAcceptedDDL(t *testing.T) {
	conn := openReplayDB(t, filepath.Join(t.TempDir(), "columns.sqlite"))
	defer closeProjectDB(t, conn)
	for table, definition := range currentEventTables {
		t.Run(table, func(t *testing.T) {
			rows, err := conn.Query(`PRAGMA table_info("` + table + `")`)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = rows.Close() }()
			found := map[string]bool{}
			primary := map[string]bool{}
			for rows.Next() {
				var index, notnull, pk int
				var name, kind string
				var fallback any
				if err := rows.Scan(&index, &name, &kind, &notnull, &fallback, &pk); err != nil {
					t.Fatal(err)
				}
				want, exists := definition.columns[name]
				if !exists {
					t.Errorf("missing accepted column %s", name)
					continue
				}
				found[name] = true
				if pk > 0 {
					primary[name] = true
				}
				integer := strings.Contains(kind, "INT")
				if (strings.HasPrefix(want, "i") || want == "b") != integer {
					t.Errorf("%s type %s differs from %s", name, kind, want)
				}
				if strings.HasSuffix(want, "?") != (notnull == 0 && pk == 0) {
					t.Errorf("%s nullability differs", name)
				}
			}
			if err := rows.Err(); err != nil {
				t.Fatal(err)
			}
			if len(found) != len(definition.columns) {
				t.Fatal("column catalog includes non-DDL fields")
			}
			if len(primary) != len(definition.primary) {
				t.Fatal("primary-key size differs")
			}
			for _, name := range definition.primary {
				if !primary[name] {
					t.Errorf("%s is not a primary-key column", name)
				}
			}
		})
	}
}

// TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues protects the
// strict journal path reached after Task 4 removes review.recorded's fence.
func TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues(t *testing.T) {
	for _, kind := range []string{
		"Review", "Coverage", "CoverageItem", "CoverageExclusion", "Finding",
		"FindingDisposition", "JudgmentSubject",
	} {
		if _, found := currentValueShapes[kind]; !found {
			t.Fatalf("currentValueShapes[%q] is not registered", kind)
		}
	}
	if currentValueShapes["Subject"] != "type:=run,envelope,activation id:ID" {
		t.Fatalf("Subject = %q, want retained Phase 22 shape", currentValueShapes["Subject"])
	}
}

func TestReviewRecordedAloneAdmitsRepeatedFindingDispositionEffects(t *testing.T) {
	dispositions := []Effect{
		{Op: "insert", Table: "finding_dispositions", PK: map[string]any{"id": "FDP-one"}},
		{Op: "insert", Table: "finding_dispositions", PK: map[string]any{"id": "FDP-two"}},
	}
	if err := validateCurrentCardinality(Event{EventType: "review.recorded", Effects: dispositions}); err != nil {
		t.Fatalf("review.recorded repeated Finding dispositions: %v", err)
	}
	err := validateCurrentCardinality(Event{EventType: "finding.addressed", Effects: dispositions})
	if err == nil || err.Error() != "current Event finding.addressed: too many finding_dispositions/insert Effects" {
		t.Fatalf("finding.addressed repeated Finding dispositions error = %v, want cardinality rejection", err)
	}
}

func TestCurrentRecoveryAdmitsTerminalPredecessorAndOptionalClaimDisposal(t *testing.T) {
	for _, terminal := range []bool{false, true} {
		for _, expired := range []bool{false, true} {
			name := fmt.Sprintf("terminal=%t/expired=%t", terminal, expired)
			t.Run(name, func(t *testing.T) {
				var states []map[string][]string
				for _, mode := range []string{"live", "replay"} {
					c := openReplayDB(t, filepath.Join(t.TempDir(), "projection.sqlite"))
					defer closeProjectDB(t, c)
					seedRecoveredProjection(t, c)
					event := recoveredProjectionEvent()
					var effects []Effect
					for _, effect := range event.Effects {
						if terminal && effect.Op == "update" && (effect.Table == "runs" || effect.Table == "run_runtime") {
							continue
						}
						if expired && effect.Op == "update" && effect.Table == "claims" {
							continue
						}
						effects = append(effects, effect)
					}
					event.Effects = effects
					if terminal {
						if _, err := c.Exec(`UPDATE runs SET status='failed',terminal_repository_json='{}',ended_at='2026-09-06T00:00:30Z' WHERE id='RUN-prior'`); err != nil {
							t.Fatal(err)
						}
					}
					if expired {
						if _, err := c.Exec(`UPDATE claims SET status='expired' WHERE id='CLM-prior'`); err != nil {
							t.Fatal(err)
						}
					}
					raw, err := MarshalCanonicalJSON(event)
					if err != nil {
						t.Fatal(err)
					}
					if mode == "live" {
						hash, err := HashEvent(event)
						if err != nil {
							t.Fatal(err)
						}
						err = ProjectEvent(c, ProjectionInput{Event: event, EventHash: hash, File: "2026-09-06.jsonl", Offset: int64(len(raw) + 1)})
						if err != nil {
							t.Fatal(err)
						}
					} else {
						workspace := newReplayWorkspace(t)
						writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), event)
						if _, err := ReplayJournal(c, filepath.Join(workspace, ".fid", "journal"), nil); err != nil {
							t.Fatal(err)
						}
					}
					var priorStatus, claimStatus string
					if err := c.QueryRow(`SELECT status FROM runs WHERE id='RUN-prior'`).Scan(&priorStatus); err != nil {
						t.Fatal(err)
					}
					if err := c.QueryRow(`SELECT status FROM claims WHERE id='CLM-prior'`).Scan(&claimStatus); err != nil {
						t.Fatal(err)
					}
					wantRun, wantClaim := "interrupted", "superseded"
					if terminal {
						wantRun = "failed"
					}
					if expired {
						wantClaim = "expired"
					}
					if priorStatus != wantRun || claimStatus != wantClaim {
						t.Fatalf("prior Run/Claim = %s/%s", priorStatus, claimStatus)
					}
					states = append(states, currentProjectionSnapshot(t, c))
				}
				if !reflect.DeepEqual(states[0], states[1]) {
					t.Fatal("conditional recovery live/replay mismatch")
				}
			})
		}
	}
}
