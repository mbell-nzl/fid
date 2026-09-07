package main

import (
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestActivationFreezesFirstAndUpgradeBindingsAcrossStagingAndRecovery(t *testing.T) {
	t.Parallel()
	for _, upgrade := range []bool{false, true} {
		name := "first"
		if upgrade {
			name = "upgrade"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			workspace := clonePristineWorkspace(t)
			var prior workflow.Activation
			var id string
			if upgrade {
				workspace, _, _ = buildFinalizedEnvelopeWorkspace(t, validGateDefinitions())
				prior = currentFinalizedActivationForTest(t, workspace)
				id = prepareUpgradeActivationForTest(t, workspace, prior)
			} else {
				id = prepareActivationForTest(t, workspace)
			}
			prepared := activationBindingStatus(t, workspace, id)
			if prepared.Version != 1 || len(prepared.Resources) != 5 {
				t.Fatalf("prepare bindings = %#v", prepared)
			}
			if upgrade {
				if prepared.PriorActivationID == nil || *prepared.PriorActivationID != prior.ID {
					t.Fatalf("upgrade predecessor = %#v", prepared.PriorActivationID)
				}
			} else if prepared.PriorActivationID != nil {
				t.Fatalf("first predecessor = %#v", prepared.PriorActivationID)
			}
			for _, resource := range prepared.Resources {
				if resource.State != "pending" || resource.Version != 1 || resource.ObservedIdentity != nil {
					t.Fatalf("initial resource = %#v", resource)
				}
				if upgrade {
					old := activationBindingResource(t, prior, resource.Component).NewIdentity
					if resource.OldIdentity == nil || *resource.OldIdentity != old || resource.NewIdentity != old+":upgrade" {
						t.Fatalf("upgrade resource = %#v, prior=%q", resource, old)
					}
				} else if resource.OldIdentity != nil || resource.NewIdentity != resource.Component+":new" {
					t.Fatalf("first resource = %#v", resource)
				}
			}
			stageActivationBindingComponents(t, workspace, id)
			assertActivationFrozenBindings(t, prepared, activationBindingStatus(t, workspace, id))
			observeActivationResourcesForTest(t, workspace, id, "old")
			observed := activationBindingStatus(t, workspace, id)
			assertActivationFrozenBindings(t, prepared, observed)
			options := workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), BootstrapActor: "binding-test", IdempotencyKey: "binding-select"}
			selected, err := workflow.RecoverActivation(options, workflow.ActivationRecoveryRequest{ActivationID: id, ExpectedVersion: observed.Version, ExpectedStatus: observed.Status, Direction: "restore-prior", Reason: "binding preservation"})
			if err != nil {
				t.Fatal(err)
			}
			assertActivationFrozenBindings(t, prepared, selected.Activation)
			versions := make([]workflow.ActivationResourceVersion, 0, 5)
			for _, resource := range selected.Activation.Resources {
				versions = append(versions, workflow.ActivationResourceVersion{Component: resource.Component, Version: resource.Version})
			}
			options.IdempotencyKey = "binding-complete"
			completed, err := workflow.CompleteActivationRecovery(options, workflow.ActivationRecoveryCompleteRequest{ActivationID: id, ExpectedVersion: selected.Activation.Version, SelectedEventID: selected.EventID, ResourceVersions: versions})
			if err != nil || completed.Activation.Status != "aborted" {
				t.Fatalf("restore = (%#v, %v)", completed, err)
			}
			assertActivationFrozenBindings(t, prepared, completed.Activation)
			if upgrade && !reflect.DeepEqual(prior, activationBindingStatus(t, workspace, prior.ID)) {
				t.Fatal("restore changed immutable predecessor")
			}
		})
	}
}

func TestActivationCommitRejectsMissingMismatchedAndMixedComponentDigests(t *testing.T) {
	t.Parallel()
	workspace, _, _ := buildFinalizedEnvelopeWorkspace(t, validGateDefinitions())
	prior := currentFinalizedActivationForTest(t, workspace)
	id := prepareUpgradeActivationForTest(t, workspace, prior)
	digests := stageActivationBindingComponents(t, workspace, id)
	observeActivationResourcesForTest(t, workspace, id, "new")
	before := activationBindingStatus(t, workspace, id)
	events := journalEventCount(t, workspace)
	for name, values := range map[string][]string{
		"missing":    {digests[0]},
		"mismatched": {digests[0], "sha256:" + strings.Repeat("0", 64)},
		"mixed":      {digests[0], prior.ComponentDigests[0]},
		"prior-only": prior.ComponentDigests,
	} {
		t.Run(name, func(t *testing.T) {
			values = append([]string(nil), values...)
			sort.Strings(values)
			path := filepath.Join(workspace, "binding-"+name+".json")
			writeActivationInput(t, path, activationCommitInput(t, workspace, id, values))
			result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", path, "--idempotency-key", "binding-"+name)
			if result.err == nil || result.stdout != "" {
				t.Fatalf("%s digest selection = %#v", name, result)
			}
			if !reflect.DeepEqual(before, activationBindingStatus(t, workspace, id)) || journalEventCount(t, workspace) != events {
				t.Fatal("rejected digest selection changed activation or Events")
			}
		})
	}
	fixture := activationLifecycleFixture{workspace: workspace, activationID: id, componentDigests: digests}
	for _, operation := range []string{"commit", "finalize"} {
		if result := executeActivationAgreementTransition(t, fixture, operation); result.err != nil {
			t.Fatalf("exact new set %s = %#v", operation, result)
		}
	}
	finalized := activationBindingStatus(t, workspace, id)
	if !reflect.DeepEqual(finalized.ComponentDigests, digests) || activationBindingStatus(t, workspace, prior.ID).Status != "superseded" {
		t.Fatalf("selected final set = %#v", finalized)
	}
	actor, err := workflow.ResolveTrustedActor(workspace, issueMigrationsDir(), "actor-builder")
	if err != nil || actor.ProvenanceRef != "activation:profile/binding-new" {
		t.Fatalf("selected actor = (%#v, %v)", actor, err)
	}
}

func TestStagedActivationContentCannotBeReplacedUnderANewKey(t *testing.T) {
	t.Parallel()
	fixture := activationLifecycleFixtureAt(t, "staged")
	before := activationBindingStatus(t, fixture.workspace, fixture.activationID)
	events := journalEventCount(t, fixture.workspace)
	path := filepath.Join(fixture.workspace, "replace-profile.json")
	writeActivationInput(t, path, map[string]any{"activation_id": fixture.activationID, "component_type": "authority-profile", "component_id": "lifecycle-builders", "component_version": 1, "expected_prior_digest": nil, "expected_version": before.Version, "content": []any{map[string]any{"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "changed:profile", "session_ref": nil, "capabilities": []string{"claim.acquire"}}}})
	result := executeCommandForTest(fixture.workspace, "--project", fixture.workspace, "--json", "activation", "component", "--input", path, "--idempotency-key", "replace-profile-new-key")
	if result.err == nil || result.stdout != "" {
		t.Fatalf("replace content = %#v", result)
	}
	if !reflect.DeepEqual(before, activationBindingStatus(t, fixture.workspace, fixture.activationID)) || journalEventCount(t, fixture.workspace) != events {
		t.Fatal("replacement changed immutable staged content or Events")
	}
}

func TestUpgradeForwardRecoverySelectsOnlyNewFinalTrust(t *testing.T) {
	t.Parallel()
	workspace, _, _ := buildFinalizedEnvelopeWorkspace(t, validGateDefinitions())
	prior := currentFinalizedActivationForTest(t, workspace)
	id := prepareUpgradeActivationForTest(t, workspace, prior)
	digests := stageActivationBindingComponents(t, workspace, id)
	observeActivationResourcesForTest(t, workspace, id, "new")
	commitPath := filepath.Join(workspace, "forward-recovery-commit.json")
	writeActivationInput(t, commitPath, activationCommitInput(t, workspace, id, digests))
	if committed := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "forward-recovery-commit"); committed.err != nil {
		t.Fatalf("commit upgrade: %#v", committed)
	}
	committed := activationBindingStatus(t, workspace, id)
	selected, err := workflow.RecoverActivation(workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), BootstrapActor: "binding-test", IdempotencyKey: "forward-recovery-select"}, workflow.ActivationRecoveryRequest{ActivationID: id, ExpectedVersion: committed.Version, ExpectedStatus: "committed", Direction: "complete-forward", Reason: "complete installed upgrade"})
	if err != nil || selected.Activation.Recovery == nil || selected.Activation.Recovery.Direction != "complete-forward" {
		t.Fatalf("RecoverActivation() = (%#v, %v)", selected, err)
	}
	versions := make([]workflow.ActivationResourceVersion, 0, len(selected.Activation.Resources))
	for _, resource := range selected.Activation.Resources {
		versions = append(versions, workflow.ActivationResourceVersion{Component: resource.Component, Version: resource.Version})
	}
	completed, err := workflow.CompleteActivationRecovery(workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), BootstrapActor: "binding-test", IdempotencyKey: "forward-recovery-complete"}, workflow.ActivationRecoveryCompleteRequest{ActivationID: id, ExpectedVersion: selected.Activation.Version, SelectedEventID: selected.EventID, ResourceVersions: versions})
	if err != nil || completed.Activation.Status != "finalized" || activationBindingStatus(t, workspace, prior.ID).Status != "superseded" {
		t.Fatalf("CompleteActivationRecovery() = (%#v, %v), prior=%#v", completed, err, activationBindingStatus(t, workspace, prior.ID))
	}
	actor, err := workflow.ResolveTrustedActor(workspace, issueMigrationsDir(), "actor-builder")
	if err != nil || actor.ProvenanceRef != "activation:profile/binding-new" {
		t.Fatalf("selected forward actor = (%#v, %v)", actor, err)
	}
}

func stageActivationBindingComponents(t *testing.T, workspace, id string) []string {
	t.Helper()
	profile := stageActivationComponentForTest(t, workspace, id, "authority-profile", "binding-new", []any{map[string]any{"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/binding-new", "session_ref": nil, "capabilities": []string{"claim.acquire", "gate.register"}}})
	definitions := validGateDefinitions()
	definitions[0].(map[string]any)["timeout_ms"] = 2345
	gate := stageActivationComponentForTest(t, workspace, id, "gate-definitions", "binding-new-gates", definitions)
	digests := []string{profile, gate}
	sort.Strings(digests)
	return digests
}

func activationBindingStatus(t *testing.T, workspace, id string) workflow.Activation {
	t.Helper()
	activation, err := workflow.ActivationStatus(workspace, issueMigrationsDir(), id)
	if err != nil {
		t.Fatal(err)
	}
	return activation
}

func activationBindingResource(t *testing.T, activation workflow.Activation, component string) workflow.ActivationResource {
	t.Helper()
	for _, resource := range activation.Resources {
		if resource.Component == component {
			return resource
		}
	}
	t.Fatalf("activation %s missing %s", activation.ID, component)
	return workflow.ActivationResource{}
}

func assertActivationFrozenBindings(t *testing.T, before, after workflow.Activation) {
	t.Helper()
	if before.ID != after.ID || !reflect.DeepEqual(before.PriorActivationID, after.PriorActivationID) || len(before.Resources) != len(after.Resources) {
		t.Fatalf("activation binding changed: %#v -> %#v", before, after)
	}
	for _, resource := range before.Resources {
		actual := activationBindingResource(t, after, resource.Component)
		if !reflect.DeepEqual(resource.OldIdentity, actual.OldIdentity) || resource.NewIdentity != actual.NewIdentity {
			t.Fatalf("frozen resource changed: %#v -> %#v", resource, actual)
		}
	}
}
