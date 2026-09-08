//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func p22AssertActivationFaultFacts(t *testing.T, workspace string, prior, original workflow.Activation, operation, direction string, before map[string][]string, result map[string]any) {
	t.Helper()
	current := p21ActivationStatus(t, workspace, original.ID)
	if current.PriorActivationID == nil || *current.PriorActivationID != prior.ID || current.Recovery == nil || current.Recovery.Direction != direction || current.Recovery.ActorAttribution != "fq7-bootstrap" || current.Recovery.Reason != "recover interrupted installation" || !reflect.DeepEqual(current.ComponentDigests, original.ComponentDigests) {
		t.Fatal("recovery changed immutable activation/component binding")
	}
	if len(current.Resources) != 5 {
		t.Fatal("recovery lost resource identities")
	}
	for i, resource := range original.Resources {
		actual := current.Resources[i]
		if actual.Component != resource.Component || !reflect.DeepEqual(actual.OldIdentity, resource.OldIdentity) || actual.NewIdentity != resource.NewIdentity {
			t.Fatal("recovery rewrote frozen resource binding")
		}
	}
	eventID := p21PathString(t, result, "event", "id")
	after := p22FaultProjection(t, workspace)
	allowed := map[string]bool{"kit_activations": true, "activation_recoveries": true, "idempotency_results": true, "applied_events": true, "projection_meta": true}
	if operation == "recover" {
		allowed["stops"], allowed["attention_items"], allowed["attention_stop_links"] = true, true, true
		if current.Status != original.Status || current.Version != original.Version+1 || current.Recovery.SelectedEventID != eventID || current.Recovery.CompletedEventID != nil {
			t.Fatal("selection lost durable identity/version")
		}
	} else {
		allowed["activation_components"], allowed["stop_resolutions"] = true, true
		status := "aborted"
		if direction == "complete-forward" {
			status = "finalized"
		}
		if current.Status != status || current.Recovery.CompletedEventID == nil || *current.Recovery.CompletedEventID != eventID {
			t.Fatal("completion lost terminal identity")
		}
		if current.Recovery.SelectedEventID == eventID {
			t.Fatal("completion replaced selection Event")
		}
	}
	for table, rows := range before {
		if !allowed[table] && !reflect.DeepEqual(rows, after[table]) {
			t.Fatalf("recovery changed preserved %s", table)
		}
	}
	if operation == "recover-complete" {
		expected := make([]string, 0, len(before["activation_components"]))
		for _, encoded := range before["activation_components"] {
			var row []any
			p22FaultOK(t, json.Unmarshal([]byte(encoded), &row))
			if row[0] == original.ID {
				row[6] = current.Status
			}
			value, err := json.Marshal(row)
			p22FaultOK(t, err)
			expected = append(expected, string(value))
		}
		if !reflect.DeepEqual(expected, after["activation_components"]) {
			t.Fatal("completion changed immutable component content/bindings")
		}
	}
	gotPrior := p21ActivationStatus(t, workspace, prior.ID)
	if operation == "recover-complete" && direction == "complete-forward" {
		if gotPrior.Status != "superseded" || gotPrior.Version != prior.Version+1 {
			t.Fatal("forward completion did not atomically supersede prior")
		}
		gotPrior.Status, gotPrior.Version, gotPrior.UpdatedAt = prior.Status, prior.Version, prior.UpdatedAt
	}
	if !reflect.DeepEqual(gotPrior, prior) {
		t.Fatal("recovery changed prior trust outside forward supersession")
	}
	actor, err := workflow.ResolveTrustedActor(workspace, filepath.Join(integrationRepositoryRoot(), "migrations"), "actor-activation")
	if operation == "recover" {
		if err == nil {
			t.Fatal("incomplete recovery exposed normal Actor trust")
		}
	} else {
		p22FaultOK(t, err)
		provenance := "activation:corpus"
		if direction == "complete-forward" {
			provenance = "activation:fq7-new"
		}
		if actor.ProvenanceRef != provenance {
			t.Fatal("recovery selected mixed/wrong profile")
		}
	}
	p22AssertActivationStopFacts(t, workspace, current)
}

func p22AssertActivationStopFacts(t *testing.T, workspace string, activation workflow.Activation) {
	t.Helper()
	conn := p22FaultDB(t, workspace)
	defer conn.Close()
	var count int
	err := conn.QueryRow(`SELECT COUNT(*) FROM stops s JOIN attention_stop_links l ON l.stop_id=s.id JOIN attention_items a ON a.id=l.attention_id WHERE s.subject_type='activation' AND s.subject_id=? AND s.code='activation-incomplete' AND s.actor_id='fq7-bootstrap' AND s.created_event_id=? AND a.subject_type=s.subject_type AND a.subject_id=s.subject_id AND a.code=s.code AND a.status='open' AND a.version=1 AND a.resolution_event_id IS NULL AND s.resolver_capability='stop.resolve'`, activation.ID, activation.Recovery.SelectedEventID).Scan(&count)
	p22FaultOK(t, err)
	if count != 1 {
		t.Fatal("recovery missing exact Stop/Attention/link transaction")
	}
	var resolutions int
	p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*) FROM stop_resolutions r JOIN stops s ON s.id=r.stop_id WHERE s.subject_id=?`, activation.ID).Scan(&resolutions))
	if activation.Recovery.CompletedEventID == nil {
		if resolutions != 0 {
			t.Fatal("selection prematurely disposed Stop")
		}
		return
	}
	if resolutions != 1 {
		t.Fatal("completion missing unique Stop disposition")
	}
	p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*) FROM stop_resolutions r JOIN stops s ON s.id=r.stop_id WHERE s.subject_id=? AND r.actor_id='fq7-bootstrap' AND r.fact_event_id=? AND r.created_event_id=? AND r.resolved_at<>''`, activation.ID, *activation.Recovery.CompletedEventID, *activation.Recovery.CompletedEventID).Scan(&count))
	if count != 1 {
		t.Fatal("completion lost disposition fact/Event provenance")
	}
}

func p22RejectActivationFaultConflicts(t *testing.T, binary string, input p22ActivationCrashInput, original any) {
	t.Helper()
	root := input.Options.WorkspaceRoot
	before, history := p22FaultProjection(t, root), p21JournalBytes(t, root)
	kinds := []string{"changed-request", "new-key"}
	if input.Operation == "recover" {
		kinds = append(kinds, "opposite-direction")
	}
	for _, kind := range kinds {
		key, request := input.Options.IdempotencyKey, original
		if kind == "changed-request" {
			if input.Operation == "recover" {
				changed := input.Selection
				changed.Reason += " changed"
				request = changed
			} else {
				changed := input.Completion
				changed.SelectedEventID = "EVT-wrong-selection"
				request = changed
			}
		} else {
			key += "-competing"
			current := p21ActivationStatus(t, root, input.Selection.ActivationID)
			if input.Operation == "recover" {
				changed := input.Selection
				changed.ExpectedVersion, changed.ExpectedStatus = current.Version, current.Status
				if kind == "opposite-direction" {
					if changed.Direction == "restore-prior" {
						changed.Direction = "complete-forward"
					} else {
						changed.Direction = "restore-prior"
					}
				}
				request = changed
			} else {
				changed := input.Completion
				changed.ExpectedVersion = current.Version
				request = changed
			}
		}
		encoded, err := json.Marshal(request)
		p22FaultOK(t, err)
		path := filepath.Join(t.TempDir(), "conflict.json")
		p22FaultOK(t, os.WriteFile(path, encoded, 0o600))
		stdout, stderr, err := closureRunCommand(root, binary, "--json", "activation", input.Operation, "--actor", input.Options.BootstrapActor, "--input", path, "--idempotency-key", key)
		wantCode, wantExit := "invalid-argument", 2
		if kind == "opposite-direction" {
			wantCode, wantExit = "activation-direction-conflict", 4
		}
		assertP22ProcessFailure(t, "activation."+input.Operation, wantCode, wantExit, p22ConformanceResponse{Stdout: stdout, Stderr: stderr, ExitCode: processExitCode(t, err)})
		if !reflect.DeepEqual(before, p22FaultProjection(t, root)) || string(history) != string(p21JournalBytes(t, root)) {
			t.Fatal("rejected activation conflict changed facts")
		}
	}
}
