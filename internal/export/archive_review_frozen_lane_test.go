package export

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func TestReviewArchiveAdmissionUsesFrozenLaneAfterRaise(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"replay", "import"} {
		for _, tc := range []struct {
			name, frozen, current, stream, refusal string
			profiles                               []string
		}{
			{"standard separate", "standard", "assured", "separate", "separate review requires an Assured hand-off", []string{"code"}},
			{"standard missing code", "standard", "assured", "primary", "missing required profile code", []string{}},
			{"routine raised to standard", "routine", "standard", "primary", "", []string{}},
			{"routine raised to assured", "routine", "assured", "primary", "", []string{}},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, e := reviewArchiveFixtureWithLane(t, tc.frozen)
				// The lane raise follows the hand-off in the canonical event order.
				raise := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261110-000005", EventType: "issue.updated", CommandName: "issue.update", OccurredAt: e.OccurredAt, ActorV2: e.ActorV2, Aggregate: journal.Aggregate{Type: "issue", ID: "FID-a", Sequence: 1}, Effects: []journal.Effect{
					{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-a"}, Values: map[string]any{"lane": tc.current}, Preconditions: map[string]any{"lane": tc.frozen}},
					{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "raise-lane"}, Values: map[string]any{"journal_event_id": "EVT-20261110-000005", "ref_type": "issue", "ref_id": "FID-a", "event": "issue-updated"}},
				}}
				raise = writerUpdateFixture(t, raise, 2)
				if err := applyLeanEmptyEvent(v, mode, raise); err != nil {
					t.Fatal(err)
				}
				h, err := reviewrecord.LatestHandoff(v.conn, "FID-a")
				if err != nil {
					t.Fatal(err)
				}
				var lane string
				if err := v.conn.QueryRow(`SELECT lane FROM issues WHERE id='FID-a'`).Scan(&lane); err != nil {
					t.Fatal(err)
				}
				if h.Lane != tc.frozen || lane != tc.current || lane == h.Lane {
					t.Fatalf("hand-off lane=%s task lane=%s, want %s/%s", h.Lane, lane, tc.frozen, tc.current)
				}
				var r reviewrecord.Round
				if err := json.Unmarshal([]byte(e.Effects[0].Values["record_json"].(string)), &r); err != nil {
					t.Fatal(err)
				}
				r.Stream, r.Profiles = tc.stream, tc.profiles
				r.Findings, r.Verdict = []reviewrecord.Finding{}, "APPROVE"
				e.EventID = "EVT-20261110-000006"
				setReviewArchiveRound(t, &e, r)
				before := handoffProjectionSnapshot(t, v)
				// No git or repository objects: admission must use recorded facts.
				err = applyLeanEmptyEvent(v, mode, e)
				if tc.refusal != "" {
					if err == nil || !strings.Contains(err.Error(), tc.refusal) {
						t.Fatalf("admission=%v, want %s", err, tc.refusal)
					}
					if handoffProjectionSnapshot(t, v) != before {
						t.Fatal("refusal changed projection")
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				rounds, err := reviewrecord.Read(v.conn, r.TaskID)
				if err != nil || len(rounds) != 1 || encodedLean(rounds[0]) != encodedLean(r) {
					t.Fatalf("stored rounds=%v/%v, want %v", rounds, err, r)
				}
			})
		}
	}
}
