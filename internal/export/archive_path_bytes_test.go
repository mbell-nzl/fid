package export

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

// Set the frozen fixture policy at the admission boundary, leaving all prior
// Events authentic. In particular, checks before the one under test must pass
// even when an overlay mutates that consumer of path bytes.
func pathBytesPolicy(t *testing.T, v *semanticValidator, edit func(*policy.Effective)) policy.Effective {
	t.Helper()
	var raw string
	if err := v.conn.QueryRow(`SELECT effective_json FROM policy_snapshots ORDER BY sequence DESC LIMIT 1`).Scan(&raw); err != nil {
		t.Fatal(err)
	}
	var p policy.Effective
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatal(err)
	}
	edit(&p)
	encoded, err := p.Canonical()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := v.conn.Exec(`UPDATE policy_snapshots SET effective_json=?,digest=? WHERE sequence=(SELECT max(sequence) FROM policy_snapshots)`, string(encoded), p.Digest()); err != nil {
		t.Fatal(err)
	}
	return p
}

func beforePathBytesEvent(t *testing.T, lines [][]byte, mode, family string) (*semanticValidator, journal.Event) {
	t.Helper()
	v, err := newSemanticValidator(testMigrationsFS(t), 26)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = v.close() })
	for _, raw := range lines {
		var event journal.Event
		if err := json.Unmarshal(raw, &event); err != nil {
			t.Fatal(err)
		}
		if event.EventType == family {
			return v, event
		}
		if err := applyLeanEmptyEvent(v, mode, event); err != nil {
			t.Fatalf("fixture %s: %v", event.EventType, err)
		}
	}
	t.Fatalf("fixture has no %s", family)
	return nil, journal.Event{}
}

func admitPathBytesEvent(t *testing.T, v *semanticValidator, mode string, event journal.Event) error {
	t.Helper()
	if mode == "live" {
		// Writers use this exact admission call before appending. A dry run must
		// leave no Effects behind, whether it accepts or refuses the receipt.
		before := handoffProjectionSnapshot(t, v)
		err := journal.DryRunProjectEvent(v.conn, event)
		if before != handoffProjectionSnapshot(t, v) {
			t.Fatal("live admission changed projection")
		}
		if err != nil {
			return err
		}
	}
	return applyLeanEmptyEvent(v, mode, event)
}

func setCheckPathBytes(t *testing.T, event *journal.Event, path string) {
	t.Helper()
	if event.EventType == "check.started" {
		for i := range event.Effects {
			if event.Effects[i].Table == "check_starts" {
				var start checkrecord.Start
				if err := json.Unmarshal([]byte(event.Effects[i].Values["start_json"].(string)), &start); err != nil {
					t.Fatal(err)
				}
				start.ChangedPathBytes = policy.PathBytes([]string{path})
				event.Effects[i].Values["start_json"] = encodedLean(start)
			}
		}
		return
	}
	var response checkrecord.Response
	if err := json.Unmarshal([]byte(event.Effects[0].Values["response_json"].(string)), &response); err != nil {
		t.Fatal(err)
	}
	response.Lane, response.ChangedPathBytes = "assured", policy.PathBytes([]string{path})
	event.Effects[0].Values["response_json"] = encodedLean(response)
	event.Effects[1].Values["result_json"] = encodedLean(response)
}

func TestCheckFloorAdmissionPreservesInvalidUTF8Paths(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, addLeanExportFacts)
	lines := leanJournalLines(t, fixture.workspace)
	for _, mode := range []string{"live", "replay", "import"} {
		for _, family := range []string{"check.started", "check.completed"} {
			for _, tc := range []struct {
				name, path string
				valid      bool
			}{
				{"invalid UTF-8", "notes\xff.md", true},
				{"readable replacement", "notes\ufffd.md", false},
			} {
				t.Run(mode+"/"+family+"/"+tc.name, func(t *testing.T) {
					t.Parallel()
					v, event := beforePathBytesEvent(t, lines, mode, family)
					pathBytesPolicy(t, v, func(p *policy.Effective) {
						lane := p.Lanes["standard"]
						lane.Checks, lane.ChecksSet = nil, false
						p.Lanes["standard"] = lane
						lane.Checks, lane.ChecksSet = []string{}, true
						p.Lanes["assured"] = lane
					})
					setCheckPathBytes(t, &event, tc.path)
					before := handoffProjectionSnapshot(t, v)
					err := admitPathBytesEvent(t, v, mode, event)
					if tc.valid {
						if err != nil {
							t.Fatalf("exact-byte floor refused: %v", err)
						}
					} else {
						want := "check definition differs from trusted policy"
						if family == "check.completed" {
							want = "empty check requires a started task, lane and current hold"
						}
						if err == nil || !strings.Contains(err.Error(), want) || before != handoffProjectionSnapshot(t, v) {
							t.Fatalf("readable name must not justify an Assured check: %v", err)
						}
					}
				})
			}
		}
	}
}

func TestHandoffAdmissionPreservesInvalidUTF8Paths(t *testing.T) {
	t.Parallel()
	lines, _ := checkedHandoffArchiveFixture(t)
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name, path  string
			conditional bool
		}{
			{"invalid UTF-8", "notes\xff.md", true},
			{"readable replacement", "notes\ufffd.md", false},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, event := beforePathBytesEvent(t, lines, mode, "task.handed-off")
				// Use the same lane for both names so required-check recomputation
				// cannot hide behind a lane mismatch. The recorded floor still differs.
				if _, err := v.conn.Exec(`UPDATE issues SET lane='assured'`); err != nil {
					t.Fatal(err)
				}
				p := pathBytesPolicy(t, v, func(p *policy.Effective) {
					p.Lanes["assured"] = p.Lanes["standard"]
					for name, check := range p.Checks {
						check.WhenChanged = []string{"**/*.go"}
						p.Checks[name] = check
					}
				})
				var receipt handoffrecord.Record
				if err := json.Unmarshal([]byte(event.Effects[0].Values["record_json"].(string)), &receipt); err != nil {
					t.Fatal(err)
				}
				receipt.Lane = "assured"
				receipt.SetPaths([]string{tc.path})
				floor, err := policy.ComputeFloor(p, []string{tc.path})
				if err != nil {
					t.Fatal(err)
				}
				receipt.Floor = &floor
				if !tc.conditional {
					receipt.RequiredChecks, receipt.ResultIDs = []string{}, []string{}
					receipt.DefinitionDigests = map[string]string{}
				}
				event.Effects[0].Values["record_json"] = encodedLean(receipt)
				var response map[string]any
				if err := json.Unmarshal([]byte(event.Effects[3].Values["result_json"].(string)), &response); err != nil {
					t.Fatal(err)
				}
				response["fields"] = receipt.Fields()
				event.Effects[3].Values["result_json"] = encodedLean(response)
				if err := admitPathBytesEvent(t, v, mode, event); err != nil {
					t.Fatalf("serialized hand-off refused: %v", err)
				}
				assertCheckedHandoffApplied(t, v, receipt, event)
			})
		}
	}
}
