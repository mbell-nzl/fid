package export

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/usagerecord"
)

func usageArchiveFixture(t *testing.T) (*semanticValidator, journal.Event, usagerecord.Record) {
	t.Helper()
	v, e := handoffArchiveFixture(t)
	if err := migrateSemanticProjection(v.conn, testMigrationsFS(t), 43); err != nil {
		t.Fatal(err)
	}
	var err error
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	v.profile = 43
	parsed, err := usagerecord.Parse("codex", []byte(`{"type":"turn.completed","usage":{"input_tokens":10,"cached_input_tokens":2,"cache_write_input_tokens":0,"output_tokens":3,"reasoning_output_tokens":1}}`), "")
	if err != nil {
		t.Fatal(err)
	}
	r := usagerecord.Record{ID: "USG-123456", TaskID: "FID-a", HoldID: "HLD-1", Actor: "agent:builder", Role: "build", InputDigest: "sha256:" + strings.Repeat("a", 64), RecordedAt: e.OccurredAt.UTC().Format("2006-01-02T15:04:05Z"), Parsed: parsed}
	e.EventType, e.CommandName = "usage.recorded", "usage.record"
	e.Aggregate = journal.Aggregate{Type: "usage", ID: r.ID, Sequence: 1}
	setUsageArchive(t, &e, r)
	return v, e, r
}
func setUsageArchive(t *testing.T, e *journal.Event, r usagerecord.Record) {
	t.Helper()
	effect, err := journal.UsageEffect(r)
	if err != nil {
		t.Fatal(err)
	}
	e.Effects = []journal.Effect{effect}
	takeoverArchiveRetry(e)
	e.Idempotency.Scope = "usage.record:project-test"
	retry := &e.Effects[1]
	retry.PK["operation_scope"] = e.Idempotency.Scope
	retry.Values["result_json"] = encodedLean(map[string]any{"usage": r, "replayed": false})
}
func TestUsageAdmission(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, name := range []string{"valid", "actor", "task", "hold", "model", "role", "stream", "counts", "models", "reported", "digest", "parser", "version", "time", "projection", "retry", "command", "aggregate", "extra", "legacy", "update", "alias", "omitted", "oversize"} {
			t.Run(mode+"/"+name, func(t *testing.T) {
				t.Parallel()
				v, e, r := usageArchiveFixture(t)
				switch name {
				case "actor":
					r.Actor = "agent:other"
					e.ActorV2.ActorID = "other"
				case "task":
					r.TaskID = "FID-other"
				case "hold":
					r.HoldID = "HLD-missing"
				case "model":
					r.Model = "other/model"
				case "role":
					r.Role = "fix"
				case "stream":
					r.Stream = "primary"
				case "counts":
					r.Counts.Total++
				case "models":
					r.Models[""] = usagerecord.Counts{}
				case "reported":
					r.Reported = []json.RawMessage{json.RawMessage(`{}`)}
				case "digest":
					r.InputDigest = "bad"
				case "parser":
					r.ParserRevision = "codex/99"
				case "version":
					r.ToolVersion = "v1"
					r.VersionSource = "output"
				case "time":
					r.RecordedAt = "2026-11-10T09:09:00Z"
				case "command":
					e.CommandName = "task.start"
				case "aggregate":
					e.Aggregate.ID = "USG-other"
				case "oversize":
					r.ToolVersion = strings.Repeat("a", 1024*1024)
					r.VersionSource = "caller"
				}
				setUsageArchive(t, &e, r)
				switch name {
				case "projection":
					e.Effects[0].Values["total"] = 100
				case "retry":
					e.Effects[1].Values["result_json"] = `{"usage":{},"replayed":false}`
				case "extra":
					e.Effects = append(e.Effects, e.Effects[0])
				case "legacy":
					e.EventSchema = journal.EventSchemaV1
				case "update":
					e.Effects[0].Op = "update"
				case "alias":
					e.Effects[0].Table = "USAGE_SESSIONS"
				case "omitted":
					e.Effects[0].Values["record_json"] = strings.Replace(e.Effects[0].Values["record_json"].(string), `,"version_source":""`, "", 1)
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyLeanEmptyEvent(v, mode, e)
				if name == "valid" {
					if err != nil {
						t.Fatal(err)
					}
					var count int
					if err := v.conn.QueryRow(`SELECT count(*) FROM usage_sessions WHERE total=13`).Scan(&count); err != nil || count != 1 {
						t.Fatalf("usage: %d %v", count, err)
					}
					return
				}
				if err == nil {
					t.Fatal("forgery admitted")
				}
				if handoffProjectionSnapshot(t, v) != before {
					t.Fatal("refused Event changed projection")
				}
			})
		}
	}
}
func TestUsageArchiveProfile(t *testing.T) {
	t.Parallel()
	v, e, _ := usageArchiveFixture(t)
	v.profile = 42
	if err := applyLeanEmptyEvent(v, "import", e); err == nil {
		t.Fatal("profile 42 admitted usage")
	}
}

func TestUsageFamilyRequiresEventV2(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			v, e, _ := usageArchiveFixture(t)
			e.EventSchema = journal.EventSchemaV1
			e.Effects = []journal.Effect{{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-a"}, Values: map[string]any{"title": "Forged legacy usage"}}}
			before := handoffProjectionSnapshot(t, v)
			if err := applyLeanEmptyEvent(v, mode, e); err == nil {
				t.Fatal("new usage family admitted as legacy Issue history")
			}
			if handoffProjectionSnapshot(t, v) != before {
				t.Fatal("refused legacy usage changed projection")
			}
		})
	}
}
