//go:build integration

package integration

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func loopAllowanceContract(t *testing.T, c *loopContract, w loopWorker, task string) {
	t.Helper()
	w.call("rework", "allow", task, "--rounds", "1", "--actor", "human:mark", "--reason", "Fix the retry case")
	latest := w.call("rework", "allow", task, "--rounds", "2", "--actor", "human:mark", "--reason", "Cover both retry edges")
	decision := p21PathString(t, latest, "result", "decision_id")
	other := loopWorker{t, w.binary, gittest.AddWorktree(t, w.root), "agent:other"}
	other.wantNext("", true, "none", "")
	var hold string
	for _, test := range []struct {
		name, kind string
		claim      bool
	}{
		{"next-allowed-rework", "rework", false},
		{"next-claim-allowed-rework", "rework", true},
		{"next-allowed-resume", "resume", false},
		{"next-claim-allowed-resume", "resume", true},
	} {
		got := w.wantNext("", test.claim, test.kind, task)
		allowance := p21Path(t, got, "result", "rework_allowance").(map[string]any)
		if allowance["decision_id"] != decision || allowance["rounds"] != float64(2) || allowance["reason"] != "Cover both retry edges" || allowance["actor"] != "human:mark" || allowance["at"] != p21Path(t, latest, "event", "occurred_at") {
			t.Fatalf("latest guidance lost or rounds became the total allowance/cap: %v", allowance)
		}
		if test.claim {
			id := p21PathString(t, got, "result", "hold", "id")
			if hold != "" && hold != id {
				t.Fatal("allowance resume replaced the live hold")
			}
			hold = id
		}
		c.observe(test.name, got, 0)
	}
}

func loopUnacceptedContract(t *testing.T, c *loopContract) {
	t.Helper()
	for _, reservation := range []string{"build", "review"} {
		root, binary := handoffJourneyRepo(t, true)
		builder := loopWorker{t, binary, root, "agent:builder"}
		task := builder.create("task", "Reserved before acceptance was required", "")
		owner, initial, kind := builder, "new", "resume"
		if reservation == "review" {
			builder.build("", task, "unaccepted-review")
			owner = loopWorker{t, binary, gittest.AddWorktree(t, root), "agent:reviewer"}
			initial, kind = "review", "review"
		}
		taken := owner.wantNext("", true, initial, task)
		hold := p21PathString(t, taken, "result", "hold", "id")
		// Recreate an agent-created reservation predating 2.1b. Current public
		// commands correctly refuse to reserve unaccepted work in the first place.
		conn := p22FaultDB(t, root)
		_, err := conn.Exec(`UPDATE issue_creators SET actor='agent:planner' WHERE issue_id=?`, task)
		closeErr := conn.Close()
		p22FaultOK(t, err)
		p22FaultOK(t, closeErr)
		unrelated := builder.create("task", "Accepted unrelated work", "")
		for _, claim := range []bool{false, true} {
			got := owner.wantNext("", claim, "none", "")
			if p21Path(t, got, "result", "waiting") != nil {
				t.Fatalf("acceptance invented a waiting kind: %v", got)
			}
			name := "next-unaccepted-" + reservation
			if claim {
				name = "next-claim-unaccepted-" + reservation
			}
			c.observe(name, got, 0)
		}
		other := loopWorker{t, binary, gittest.AddWorktree(t, root), "agent:other"}
		other.wantNext("", true, "new", unrelated)
		builder.call("accept", task, "--actor", "human:mark")
		resumed := owner.wantNext("", true, kind, task)
		if p21PathString(t, resumed, "result", "hold", "id") != hold {
			t.Fatal("waiting for acceptance released the reservation")
		}
	}
}
