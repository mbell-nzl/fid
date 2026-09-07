package main

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestActivationRequiresEveryResourceAtCommitAndFinalize(t *testing.T) {
	t.Parallel()
	for _, operation := range []string{"commit", "finalize"} {
		for _, component := range []string{"kit", "overlay", "provider", "project", "backend"} {
			for _, defect := range []string{"missing-row", "missing-observation", "pending", "mismatched-identity"} {
				t.Run(operation+"/"+component+"/"+defect, func(t *testing.T) {
					t.Parallel()
					fixture := activationAgreementFixture(t, operation)
					// Isolate the transition's agreement guard with admitted projection
					// states. These fixtures do not claim that an observation producer
					// can write a mismatched identity or that corruption replays.
					setActivationResourceDefect(t, fixture, component, defect)
					before, err := workflow.ActivationStatus(fixture.workspace, issueMigrationsDir(), fixture.activationID)
					if err != nil {
						t.Fatal(err)
					}
					eventsBefore := journalEventCount(t, fixture.workspace)
					result := executeActivationAgreementTransition(t, fixture, operation)
					var failure struct {
						OK        bool   `json:"ok"`
						Operation string `json:"operation"`
						Error     struct {
							Code    string `json:"code"`
							Details struct {
								ActivationID string   `json:"activation_id"`
								Components   []string `json:"components"`
								Reasons      []struct {
									Code        string `json:"code"`
									Detail      string `json:"detail"`
									SubjectID   string `json:"subject_id"`
									SubjectType string `json:"subject_type"`
								} `json:"reasons"`
							} `json:"details"`
						} `json:"error"`
					}
					if result.err == nil || result.stdout != "" || json.Unmarshal([]byte(result.stderr), &failure) != nil || failure.OK || failure.Operation != "activation."+operation || failure.Error.Code != "activation-incomplete" {
						t.Fatalf("%s with %s %s = %#v, want structured rejection", operation, component, defect, result)
					}
					if failure.Error.Details.ActivationID != fixture.activationID || len(failure.Error.Details.Components) != 1 || failure.Error.Details.Components[0] != component || len(failure.Error.Details.Reasons) != 1 || failure.Error.Details.Reasons[0].Code != "activation-incomplete" || !strings.Contains(failure.Error.Details.Reasons[0].Detail, component) || failure.Error.Details.Reasons[0].SubjectID != fixture.activationID || failure.Error.Details.Reasons[0].SubjectType != "activation" {
						t.Fatalf("%s with %s %s details = %+v, want exact activation and component reason", operation, component, defect, failure.Error.Details)
					}
					after, err := workflow.ActivationStatus(fixture.workspace, issueMigrationsDir(), fixture.activationID)
					if err != nil || !reflect.DeepEqual(after, before) {
						t.Fatalf("rejected transition changed activation: before=%#v after=%#v err=%v", before, after, err)
					}
					if count := journalEventCount(t, fixture.workspace); count != eventsBefore {
						t.Fatalf("rejected transition wrote Events: %d -> %d", eventsBefore, count)
					}
					if actor, err := workflow.ResolveTrustedActor(fixture.workspace, issueMigrationsDir(), "actor-builder"); err == nil {
						t.Fatalf("incomplete activation granted trust: %#v", actor)
					}
				})
			}
		}
	}
}

func TestActivationCompleteResourceAgreementFinalizesTrust(t *testing.T) {
	t.Parallel()
	fixture := activationAgreementFixture(t, "commit")
	for _, operation := range []string{"commit", "finalize"} {
		result := executeActivationAgreementTransition(t, fixture, operation)
		if result.err != nil {
			t.Fatalf("all-agree %s = %#v, want success", operation, result)
		}
	}
	if actor, err := workflow.ResolveTrustedActor(fixture.workspace, issueMigrationsDir(), "actor-builder"); err != nil || actor.ProvenanceRef == "" {
		t.Fatalf("complete activation trust = (%#v, %v)", actor, err)
	}
}

func activationAgreementFixture(t *testing.T, operation string) activationLifecycleFixture {
	t.Helper()
	if operation == "finalize" {
		return activationLifecycleFixtureAt(t, "committed")
	}
	fixture := activationLifecycleFixtureAt(t, "staged")
	observeActivationResourcesForTest(t, fixture.workspace, fixture.activationID, "new")
	return fixture
}

func setActivationResourceDefect(t *testing.T, fixture activationLifecycleFixture, component, defect string) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(fixture.workspace, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	query := map[string]string{
		"missing-row":         "DELETE FROM activation_resources WHERE activation_id = ? AND component = ?",
		"missing-observation": "UPDATE activation_resources SET observed_identity = NULL WHERE activation_id = ? AND component = ?",
		"pending":             "UPDATE activation_resources SET state = 'pending', observed_identity = NULL, evidence_ref = NULL WHERE activation_id = ? AND component = ?",
		"mismatched-identity": "UPDATE activation_resources SET observed_identity = 'different:identity' WHERE activation_id = ? AND component = ?",
	}[defect]
	result, err := conn.Exec(query, fixture.activationID, component)
	if err != nil {
		t.Fatal(err)
	}
	if count, err := result.RowsAffected(); err != nil || count != 1 {
		t.Fatalf("resource fixture affected %d rows: %v", count, err)
	}
}

func executeActivationAgreementTransition(t *testing.T, fixture activationLifecycleFixture, operation string) commandExecutionResult {
	t.Helper()
	input := activationFinalizeInput(t, fixture.workspace, fixture.activationID)
	if operation == "commit" {
		input = activationCommitInput(t, fixture.workspace, fixture.activationID, fixture.componentDigests)
	}
	path := filepath.Join(fixture.workspace, "agreement-"+operation+".json")
	writeActivationInput(t, path, input)
	return executeCommandForTest(fixture.workspace, "--project", fixture.workspace, "--json", "activation", operation, "--input", path, "--idempotency-key", "agreement-"+operation)
}
