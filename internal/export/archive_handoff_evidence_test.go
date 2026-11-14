package export

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

func checkedHandoffArchiveFixture(t *testing.T) ([][]byte, string) {
	t.Helper()
	var olderResult string
	fixture := newExportSourceFixture(t, func(t *testing.T, root string) string {
		path := filepath.Join(root, ".fid.toml")
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		text := string(raw) + "\n[lanes.standard]\nchecks=['lint','test']\n[review]\ntest_adequacy=[]\n"
		// Distinct definitions make swapped digests a meaningful forgery.
		text += "\n[checks.lint]\nrun=['true','lint']\ntimeout='1m'\n[checks.test]\nrun=['true','test']\ntimeout='1m'\n"
		if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
			t.Fatal(err)
		}
		fixtureGit(t, root, "add", ".")
		fixtureGit(t, root, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-m", "Required checks policy")
		fixtureGit(t, root, "remote", "add", "origin", root)
		fixtureGit(t, root, "update-ref", "refs/remotes/origin/trunk", "HEAD")
		fixtureGit(t, root, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
		lane := "standard"
		task, err := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), Kind: "task", Lane: &lane, Title: "Checked hand-off", Actor: "human:mark", IdempotencyKey: "task"})
		if err != nil {
			t.Fatal(err)
		}
		held, err := holdwork.Takeover(holdwork.Options{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), TaskID: task.ID, Actor: "human:mark", Reason: "Build receipt fixture"})
		if err != nil {
			t.Fatal(err)
		}
		options := holdwork.Options{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), HoldID: held.Fields["hold_id"].(string), Actor: "human:mark", IdempotencyKey: "start"}
		if _, err := holdwork.Start(options); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"older-checks", "latest-checks"} {
			response, err := checkrun.Run(context.Background(), checkrun.Options{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), HoldID: options.HoldID, Actor: options.Actor, IdempotencyKey: key, Version: "test"})
			if err != nil || len(response.Checks) != 2 {
				t.Fatalf("fixture checks = %#v / %v", response, err)
			}
			for _, check := range response.Checks {
				if !check.Counts || check.Outcome != "passed" {
					t.Fatalf("fixture result does not count: %#v", check)
				}
			}
			if key == "older-checks" {
				olderResult = response.Checks[0].ID
			}
		}
		options.IdempotencyKey, options.PR = "handoff", "https://github.com/owner/repo/pull/66"
		if _, err := holdwork.Handoff(options); err != nil {
			t.Fatal(err)
		}
		return task.ID
	})
	return leanJournalLines(t, fixture.workspace), olderResult
}

func TestHandoffReplayAndImportValidateCountingEvidence(t *testing.T) {
	t.Parallel()
	lines, olderResult := checkedHandoffArchiveFixture(t)
	for _, mode := range []string{"replay", "import"} {
		for _, row := range []struct {
			name   string
			mutate func(*handoffrecord.Record)
		}{
			{"valid", nil},
			{"unknown result", func(r *handoffrecord.Record) { r.ResultIDs[0] = "CHK-unrecorded" }},
			{"older passing result", func(r *handoffrecord.Record) { r.ResultIDs[0] = olderResult }},
			{"swapped results", func(r *handoffrecord.Record) { r.ResultIDs[0], r.ResultIDs[1] = r.ResultIDs[1], r.ResultIDs[0] }},
			{"wrong definition digest", func(r *handoffrecord.Record) { r.DefinitionDigests["lint"] = "sha256:" + strings.Repeat("0", 64) }},
			{"swapped definition digests", func(r *handoffrecord.Record) {
				r.DefinitionDigests["lint"], r.DefinitionDigests["test"] = r.DefinitionDigests["test"], r.DefinitionDigests["lint"]
			}},
		} {
			t.Run(mode+"/"+row.name, func(t *testing.T) {
				t.Parallel()
				validator, err := newSemanticValidator(testMigrationsDir(t), 26)
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = validator.close() })
				for _, raw := range lines {
					var event journal.Event
					if err := json.Unmarshal(raw, &event); err != nil {
						t.Fatal(err)
					}
					if event.EventType != "task.handed-off" {
						if err := applyLeanEmptyEvent(validator, mode, event); err != nil {
							t.Fatalf("fixture %s: %v", event.EventType, err)
						}
						continue
					}
					var receipt handoffrecord.Record
					if err := json.Unmarshal([]byte(event.Effects[0].Values["record_json"].(string)), &receipt); err != nil {
						t.Fatal(err)
					}
					if !slices.Equal(receipt.RequiredChecks, []string{"lint", "test"}) || len(receipt.ResultIDs) != 2 || receipt.ResultIDs[0] == olderResult {
						t.Fatalf("fixture receipt = %#v", receipt)
					}
					if row.mutate != nil {
						row.mutate(&receipt)
					}
					// Keep the response consistent so only evidence validation can
					// reject an otherwise well-formed forged receipt.
					event.Effects[0].Values["record_json"] = encodedLean(receipt)
					var response map[string]any
					if err := json.Unmarshal([]byte(event.Effects[3].Values["result_json"].(string)), &response); err != nil {
						t.Fatal(err)
					}
					response["fields"] = receipt.Fields()
					event.Effects[3].Values["result_json"] = encodedLean(response)
					before := handoffProjectionSnapshot(t, validator)
					err := applyLeanEmptyEvent(validator, mode, event)
					if row.mutate == nil {
						if err != nil {
							t.Fatalf("valid checked hand-off: %v", err)
						}
						assertCheckedHandoffApplied(t, validator, receipt, event)
					} else {
						if err == nil || !strings.Contains(err.Error(), "hand-off does not name the counting result and definition") {
							t.Errorf("forgery accepted or rejected for another reason: %v", err)
						}
						if after := handoffProjectionSnapshot(t, validator); after != before {
							t.Error("rejected receipt changed projection state")
						}
					}
					return
				}
				t.Fatal("fixture has no hand-off")
			})
		}
	}
}

func assertCheckedHandoffApplied(t *testing.T, validator *semanticValidator, receipt handoffrecord.Record, event journal.Event) {
	t.Helper()
	var state, reason, raw, retry string
	var kept bool
	err := validator.conn.QueryRow(`SELECT i.status,h.end_reason,h.kept,r.record_json,k.result_json FROM issues i JOIN holds h ON h.task_id=i.id JOIN handoffs r ON r.hold_id=h.id JOIN idempotency_results k ON k.result_event_id=? WHERE i.id=?`, event.EventID, receipt.TaskID).Scan(&state, &reason, &kept, &raw, &retry)
	if err != nil {
		t.Fatal(err)
	}
	if state != "in_review" || reason != "review" || !kept || raw != encodedLean(receipt) || retry != event.Effects[3].Values["result_json"] {
		t.Fatalf("incomplete checked hand-off = %s / %s / %v / %s / %s", state, reason, kept, raw, retry)
	}
}

// Snapshot every domain table, including check history and retry bindings.
func handoffProjectionSnapshot(t *testing.T, validator *semanticValidator) string {
	t.Helper()
	snapshot := map[string][]string{}
	for table := range validator.tables {
		rows, err := validator.conn.Query(`SELECT * FROM "` + table + `"`)
		if err != nil {
			t.Fatal(err)
		}
		columns, err := rows.Columns()
		if err != nil {
			t.Fatal(err)
		}
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		for rows.Next() {
			if err := rows.Scan(pointers...); err != nil {
				t.Fatal(err)
			}
			snapshot[table] = append(snapshot[table], encodedLean(values))
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
		slices.Sort(snapshot[table])
	}
	return encodedLean(snapshot)
}
