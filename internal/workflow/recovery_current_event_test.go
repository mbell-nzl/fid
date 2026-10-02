package workflow

import (
	"fmt"
	"testing"
	"time"
)

// Exercise the real recovery producer for every durable predecessor/Claim
// branch. The journal tests independently apply/replay the resulting grammar.
func TestRecoverRunCurrentEventStateBranches(t *testing.T) {
	t.Parallel()
	for _, terminal := range []bool{false, true} {
		for _, activeClaim := range []bool{false, true} {
			t.Run(fmt.Sprintf("terminal=%v/active-claim=%v", terminal, activeClaim), func(t *testing.T) {
				options, request := recoveryFixture(t)
				conn, err := openWorkflowDB(options.WorkspaceRoot)
				if err != nil {
					t.Fatal(err)
				}
				if terminal {
					_, err = conn.Exec(`UPDATE runs SET status='failed',terminal_repository_json=start_repository_json,ended_at=? WHERE id='RUN-prior'`, time.Now().UTC().Format(time.RFC3339Nano))
					if err != nil {
						t.Fatal(err)
					}
				}
				if activeClaim {
					_, err = conn.Exec(`UPDATE claims SET status='active',expires_at=? WHERE id='CLM-prior'`, time.Now().UTC().Add(time.Hour).Format(time.RFC3339Nano))
					if err != nil {
						t.Fatal(err)
					}
				}
				closeWorkflowDB(conn)
				if activeClaim {
					configureForceRecoveryAuthority(t, options.WorkspaceRoot)
					setRecoveryObservationState(t, options.WorkspaceRoot, "live")
					reference := "HUM-state-branches"
					request.Authority = RecoveryAuthority{Mode: "force", ActorID: options.ActorID, Reason: "exact state branch", Reference: &reference, GrantRef: &RecoveryGrantRef{EnvelopeID: "ENV-recovery", Revision: 1, Ordinal: 0}}
				}
				result, err := RecoverRun(options, request)
				if err != nil {
					t.Fatal(err)
				}
				prior, err := GetRun(options.WorkspaceRoot, options.MigrationsDir, request.PriorRunID)
				if err != nil {
					t.Fatal(err)
				}
				want, version := "interrupted", 2
				if terminal {
					want, version = "failed", 1
				}
				if prior.Status != want || prior.Version != version || result.Run.Status != "active" {
					t.Fatalf("wrong predecessor/successor branch: %#v / %#v", prior, result.Run)
				}
				conn, err = openWorkflowDB(options.WorkspaceRoot)
				if err != nil {
					t.Fatal(err)
				}
				defer closeWorkflowDB(conn)
				claim, err := loadRecoveryClaim(conn, "CLM-prior")
				if err != nil {
					t.Fatal(err)
				}
				claimWant := "expired"
				if activeClaim {
					claimWant = "superseded"
				}
				if claim.Status != claimWant {
					t.Fatalf("prior Claim %q, want %q", claim.Status, claimWant)
				}
			})
		}
	}
}
