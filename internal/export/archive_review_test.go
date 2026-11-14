package export

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func reviewArchiveFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	v, h := handoffArchiveFixture(t)
	if err := migrateSemanticProjection(v.conn, testMigrationsDir(t), 30); err != nil {
		t.Fatal(err)
	}
	var err error
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	v.profile = 30
	acceptMergeArchiveEvent(t, v, h)
	r := reviewrecord.Round{ID: "RND-123456", TaskID: "FID-a", HandoffID: "HOF-100001", Reviewer: reviewrecord.Reviewer{Actor: "agent:reviewer", Model: "claude/opus-5.5", SessionID: "session-1"}, Stream: "primary", Profiles: []string{"code"}, Verdict: "CHANGES", Claim: "Persist review rounds", Verification: "go test passed", Binding: reviewrecord.Binding{Algorithm: fingerprint.Algorithm, Fingerprint: strings.Repeat("f", 40), Base: strings.Repeat("b", 40), Head: strings.Repeat("a", 40)}, RecordedBy: "human:courier", At: "2026-11-10T09:06:00Z", Findings: []reviewrecord.Finding{}}
	for i, severity := range []string{"major", "minor", "nit"} {
		r.Findings = append(r.Findings, reviewrecord.Finding{ID: []string{"FND-1", "FND-2", "FND-3"}[i], FindingInput: reviewrecord.FindingInput{Profile: "code", Severity: severity, Summary: "Problem", Evidence: "Observed behaviour", Fix: "Add regression coverage"}})
	}
	e := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261110-000004", EventType: "review.round-recorded", CommandName: "review.record", OccurredAt: time.Date(2026, 11, 10, 9, 6, 0, 0, time.UTC), ActorV2: journal.EventActor{ActorType: "human", ActorID: "courier", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{Type: "review", ID: r.ID, Sequence: 1}}
	setReviewArchiveRound(t, &e, r)
	return v, e
}

func setReviewArchiveRound(t *testing.T, e *journal.Event, r reviewrecord.Round) {
	t.Helper()
	var err error
	e.Effects, err = reviewwork.Effects(r)
	if err != nil {
		t.Fatal(err)
	}
	takeoverArchiveRetry(e)
	e.Idempotency.Scope = "review.record:project-test"
	retry := &e.Effects[len(e.Effects)-1]
	retry.PK["operation_scope"] = e.Idempotency.Scope
	retry.Values["created_at"] = r.At
	retry.Values["result_json"] = encodedLean(reviewwork.Result{Round: r, At: r.At})
}

func TestReviewArchiveReplayAndLiveAdmission(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name   string
			change func(*journal.Event, *reviewrecord.Round)
			want   bool
		}{
			{name: "complete", want: true},
			{"ended builder", func(_ *journal.Event, r *reviewrecord.Round) { r.Reviewer.Actor = "agent:builder" }, false},
			{"no required profile", func(_ *journal.Event, r *reviewrecord.Round) {
				r.Profiles = []string{}
				r.Findings = []reviewrecord.Finding{}
			}, false},
			{"wrong handoff", func(_ *journal.Event, r *reviewrecord.Round) { r.HandoffID = "HOF-other" }, false},
			{"wrong base", func(_ *journal.Event, r *reviewrecord.Round) { r.Binding.Base = strings.Repeat("d", 40) }, false},
			{"wrong head", func(_ *journal.Event, r *reviewrecord.Round) { r.Binding.Head = strings.Repeat("d", 40) }, false},
			{"invalid fingerprint", func(_ *journal.Event, r *reviewrecord.Round) { r.Binding.Fingerprint = "empty" }, false},
			{"unknown algorithm", func(_ *journal.Event, r *reviewrecord.Round) { r.Binding.Algorithm = "patch-id" }, false},
			{"bad model", func(_ *journal.Event, r *reviewrecord.Round) { r.Reviewer.Model = "claude" }, false},
			{"bad severity", func(_ *journal.Event, r *reviewrecord.Round) { r.Findings[0].Severity = "blocker" }, false},
			{"verdict mismatch", func(_ *journal.Event, r *reviewrecord.Round) { r.Verdict = "APPROVE" }, false},
			{"separate standard", func(_ *journal.Event, r *reviewrecord.Round) { r.Stream = "separate" }, false},
			{"service recorder", func(e *journal.Event, r *reviewrecord.Round) {
				e.ActorV2.ActorType = "service"
				r.RecordedBy = "service:courier"
			}, false},
			{"wrong command", func(e *journal.Event, _ *reviewrecord.Round) { e.CommandName = "review.list" }, false},
			{"wrong actor", func(_ *journal.Event, r *reviewrecord.Round) { r.RecordedBy = "agent:other" }, false},
			{"wrong time", func(_ *journal.Event, r *reviewrecord.Round) { r.At = "2026-11-10T09:07:00Z" }, false},
			{"duplicate finding ID", func(_ *journal.Event, r *reviewrecord.Round) { r.Findings[1].ID = r.Findings[0].ID }, false},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, e := reviewArchiveFixture(t)
				var r reviewrecord.Round
				if err := json.Unmarshal([]byte(e.Effects[0].Values["record_json"].(string)), &r); err != nil {
					t.Fatal(err)
				}
				if tc.change != nil {
					tc.change(&e, &r)
					setReviewArchiveRound(t, &e, r)
				}
				before := handoffProjectionSnapshot(t, v)
				// No git repository exists here: every mode validates recorded facts only.
				err := applyLeanEmptyEvent(v, mode, e)
				if (err == nil) != tc.want {
					t.Fatalf("admit=%v want=%t", err, tc.want)
				}
				if !tc.want && handoffProjectionSnapshot(t, v) != before {
					t.Fatal("refusal changed projection")
				}
				if tc.want {
					var n int
					if err := v.conn.QueryRow(`SELECT count(*) FROM round_findings`).Scan(&n); err != nil || n != 3 {
						t.Fatalf("findings=%d/%v", n, err)
					}
				}
			})
		}
	}
}

func TestReviewArchiveRefusesMalformedEffects(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name   string
			change func(*journal.Event)
		}{
			{"finding omitted", func(e *journal.Event) { e.Effects = append(e.Effects[:1], e.Effects[2:]...) }},
			{"finding forged", func(e *journal.Event) { e.Effects[1].Values["round_id"] = "RND-other" }},
			{"finding content forged", func(e *journal.Event) {
				e.Effects[1].Values["record_json"] = strings.Replace(e.Effects[1].Values["record_json"].(string), "major", "minor", 1)
			}},
			{"retry corrupt", func(e *journal.Event) { e.Effects[len(e.Effects)-1].Values["result_json"] = "{}" }},
			{"retry differs", func(e *journal.Event) {
				e.Effects[len(e.Effects)-1].Values["result_json"] = strings.Replace(e.Effects[len(e.Effects)-1].Values["result_json"].(string), "Persist review rounds", "Other claim", 1)
			}},
			{"no retry", func(e *journal.Event) { e.Effects = e.Effects[:len(e.Effects)-1] }},
			{"reordered", func(e *journal.Event) { e.Effects[0], e.Effects[1] = e.Effects[1], e.Effects[0] }},
			{"extra round", func(e *journal.Event) { e.Effects = append([]journal.Effect{e.Effects[0]}, e.Effects...) }},
			{"mutating round", func(e *journal.Event) { e.Effects[0].Op = "update" }},
			{"wrong aggregate", func(e *journal.Event) { e.Aggregate.Type = "task" }},
			{"v1", func(e *journal.Event) { e.EventSchema = 1 }},
			{"other family", func(e *journal.Event) { e.EventType = "issue.updated" }},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, e := reviewArchiveFixture(t)
				tc.change(&e)
				before := handoffProjectionSnapshot(t, v)
				if err := applyLeanEmptyEvent(v, mode, e); err == nil {
					t.Fatal("accepted malformed review")
				}
				if handoffProjectionSnapshot(t, v) != before {
					t.Fatal("refusal changed state")
				}
			})
		}
	}
}

func TestBuildModelArchiveAndReplayValidation(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, model := range []string{"", "codex/gpt-6.1-sol", "invalid", "a/b/c"} {
			t.Run(mode+"/"+model, func(t *testing.T) {
				t.Parallel()
				v, err := newSemanticValidator(testMigrationsDir(t), 30)
				if err != nil {
					t.Fatal(err)
				}
				defer func() { _ = v.close() }()
				if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-a','task','Build','task-default','todo','developer','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`); err != nil {
					t.Fatal(err)
				}
				e := archiveHoldEvent(func(e *journal.Event) { e.Effects[0].Values["model"] = model })
				before := handoffProjectionSnapshot(t, v)
				err = applyLeanEmptyEvent(v, mode, e)
				valid := model == "" || model == "codex/gpt-6.1-sol"
				if (err == nil) != valid {
					t.Fatalf("model %q admission=%v", model, err)
				}
				if !valid && handoffProjectionSnapshot(t, v) != before {
					t.Fatal("malformed model changed projection")
				}
			})
		}
	}
}

func TestBuildModelsRejectLegacyOrUnrelatedWriters(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, schema := range []int{1, 2} {
			t.Run(fmt.Sprintf("%s/v%d", mode, schema), func(t *testing.T) {
				t.Parallel()
				v, err := newSemanticValidator(testMigrationsDir(t), 30)
				if err != nil {
					t.Fatal(err)
				}
				defer func() { _ = v.close() }()
				if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-a','task','Build','task-default','todo','developer','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`); err != nil {
					t.Fatal(err)
				}
				e := archiveHoldEvent(func(e *journal.Event) { e.Effects[0].Values["model"] = "codex/gpt-6.1-sol" })
				e.EventSchema = schema
				if schema == 2 {
					e.EventType = "issue.updated"
					e.Aggregate.Type = "issue"
				}
				before := handoffProjectionSnapshot(t, v)
				err = applyLeanEmptyEvent(v, mode, e)
				if err == nil || handoffProjectionSnapshot(t, v) != before {
					t.Fatalf("forged model writer=%v", err)
				}
			})
		}
	}
}
