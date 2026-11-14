package journal

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestCurrentIdempotencyFieldDigestsReplayFailClosed(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		value any
		valid bool
	}{
		{"present", `{"--actor":"sha256:` + strings.Repeat("a", 64) + `"}`, true},
		{"missing", nil, false},
		{"array", `[]`, false},
		{"non-digest", `{"--actor":"plain"}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			event := leanStartProjectionEvent()
			if tc.value == nil {
				delete(event.Effects[len(event.Effects)-1].Values, "request_field_digests_json")
			} else {
				event.Effects[len(event.Effects)-1].Values["request_field_digests_json"] = tc.value
			}
			conn := openReplayDB(t, filepath.Join(t.TempDir(), "projection.sqlite"))
			defer closeProjectDB(t, conn)
			seedLeanStartProjection(t, conn)
			before := currentProjectionSnapshot(t, conn)
			workspace := newReplayWorkspace(t)
			writeReplayEvents(t, filepath.Join(config.LedgerDir(workspace), "journal", "2026-09-06.jsonl"), event)
			_, err := ReplayJournal(conn, filepath.Join(config.LedgerDir(workspace), "journal"), nil)
			if tc.valid {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			var corrupt *CorruptJournalError
			if !errors.As(err, &corrupt) || corrupt.Line != 1 || corrupt.File != filepath.ToSlash(filepath.Join(config.LedgerDir(workspace), "journal", "2026-09-06.jsonl")) {
				t.Fatalf("replay = %v, want line-1 journal corruption with file evidence", err)
			}
			if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
				t.Fatal("rejected map changed projection")
			}
			var count int
			if err := conn.QueryRow(`SELECT count(*) FROM idempotency_results`).Scan(&count); err != nil || count != 0 {
				t.Fatalf("rejected journal projected %d idempotency rows: %v", count, err)
			}
		})
	}
}

func TestCurrentEventFamiliesRejectMalformedEffectsLiveAndReplay(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		mutate func(*Event)
	}{
		{"unknown-family", func(e *Event) { e.EventType = "unknown.task-event" }},
		{"unknown-operation", func(e *Event) { e.Effects[0].Op = "unknown" }},
		{"missing-state-change", func(e *Event) { e.Effects = e.Effects[1:] }},
		{"missing-result", func(e *Event) { e.Effects = e.Effects[:1] }},
		{"reordered", func(e *Event) { e.Effects[0], e.Effects[1] = e.Effects[1], e.Effects[0] }},
		{"wrong-aggregate", func(e *Event) { e.Aggregate.ID = "FID-other" }},
		{"wrong-aggregate-type", func(e *Event) { e.Aggregate.Type = "issue" }},
		{"wrong-command", func(e *Event) { e.CommandName = "issue.transition" }},
		{"wrong-state", func(e *Event) { e.Effects[0].Values["status"] = "done" }},
		{"coerced-text", func(e *Event) { e.Effects[0].Values["status"] = 42 }},
		{"missing-revision-precondition", func(e *Event) { delete(e.Effects[0].Preconditions, "revision") }},
		{"revision-gap", func(e *Event) { e.Effects[0].Values["revision"] = 4 }},
		{"stale-revision", func(e *Event) { e.Effects[0].Preconditions["revision"] = 99 }},
		{"duplicate-effect", func(e *Event) { e.Effects = append([]Effect{e.Effects[0]}, e.Effects...) }},
		{"unknown-field", func(e *Event) { e.Effects[0].Values["surprise"] = true }},
		{"wrong-result-event", func(e *Event) { e.Effects[1].Values["result_event_id"] = "EVT-other" }},
		{"wrong-request-digest", func(e *Event) { e.Effects[1].Values["request_digest"] = "sha256:" }},
		{"wrong-result-entity", func(e *Event) { e.Effects[1].Values["result_entity_id"] = "FID-other" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, mode := range []string{"live", "replay"} {
				t.Run(mode, func(t *testing.T) {
					event := leanStartProjectionEvent()
					tc.mutate(&event)
					conn := openReplayDB(t, filepath.Join(t.TempDir(), "projection.sqlite"))
					defer closeProjectDB(t, conn)
					seedLeanStartProjection(t, conn)
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
						writeReplayEvents(t, filepath.Join(config.LedgerDir(workspace), "journal", "2026-09-06.jsonl"), event)
						_, err = ReplayJournal(conn, filepath.Join(config.LedgerDir(workspace), "journal"), nil)
					}
					if err == nil {
						t.Fatal("malformed current-family Event accepted")
					}
					if tc.name == "unknown-operation" && !strings.Contains(err.Error(), "unsupported effect op") {
						t.Fatalf("wrong operation refusal: %v", err)
					}
					if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
						t.Fatal("rejected Event partially projected")
					}
				})
			}
		})
	}
}

func TestCurrentEventRequiresExactCapturedPrior(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"valid", "missing", "wrong", "extra"} {
		t.Run(kind, func(t *testing.T) {
			event := leanStartProjectionEvent()
			event.Effects[0].Preconditions = map[string]any{"status": "todo", "revision": 1}
			switch kind {
			case "missing":
				event.Effects[0].Preconditions = nil
			case "wrong":
				event.Effects[0].Preconditions["revision"] = 99
			case "extra":
				event.Effects[0].Preconditions["title"] = "Fixture task"
			}
			conn := openReplayDB(t, filepath.Join(t.TempDir(), "prior.sqlite"))
			defer closeProjectDB(t, conn)
			seedLeanStartProjection(t, conn)
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
				t.Fatal("invalid prior state state accepted")
			}
			if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
				t.Fatal("invalid prior state partially projected")
			}
		})
	}
}

func TestCurrentEventColumnDomainsMatchAcceptedDDL(t *testing.T) {
	t.Parallel()
	if _, ok := currentEventTables["activity_log"]; !ok {
		t.Fatal("activity_log is missing from the strict current Event table registry")
	}
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
