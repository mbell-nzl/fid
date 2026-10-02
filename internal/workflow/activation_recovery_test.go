package workflow

import (
	"errors"
	"testing"
)

func TestActivationResourceVersionConflictsUseActivationVersionConflict(t *testing.T) {
	t.Parallel()
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	prepared, err := PrepareActivation(ActivationPrepareOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-resource-version-prepare", Request: ActivationPrepareRequest{KitIdentity: "kit:test", BackendIdentity: "backend:test", Phase: 21, Direction: "activate", Resources: activationRecoveryTestResources()}})
	if err != nil {
		t.Fatalf("PrepareActivation() error = %v", err)
	}
	resource := prepared.Activation.Resources[0]
	identity := resource.NewIdentity
	staleActivation := ActivationResourceObservationRequest{ActivationID: prepared.Activation.ID, Component: resource.Component, ExpectedActivationVersion: prepared.Activation.Version + 1, ExpectedVersion: resource.Version, ObservedIdentity: &identity, State: "new", EvidenceRef: "test:stale-activation"}
	_, err = ObserveActivationResource(ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-resource-stale-activation"}, staleActivation)
	assertActivationVersionConflict(t, err, staleActivation.ExpectedActivationVersion, prepared.Activation.Version)

	observed, err := ObserveActivationResource(ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-resource-current"}, ActivationResourceObservationRequest{ActivationID: prepared.Activation.ID, Component: resource.Component, ExpectedActivationVersion: prepared.Activation.Version, ExpectedVersion: resource.Version, ObservedIdentity: &identity, State: "new", EvidenceRef: "test:current"})
	if err != nil {
		t.Fatalf("ObserveActivationResource(current) error = %v", err)
	}
	staleResource := ActivationResourceObservationRequest{ActivationID: prepared.Activation.ID, Component: resource.Component, ExpectedActivationVersion: observed.Activation.Version, ExpectedVersion: resource.Version, ObservedIdentity: &identity, State: "new", EvidenceRef: "test:stale-resource"}
	_, err = ObserveActivationResource(ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-resource-stale-resource"}, staleResource)
	assertActivationVersionConflict(t, err, staleResource.ExpectedVersion, resource.Version+1)
}

func assertActivationVersionConflict(t *testing.T, err error, expected, actual int) {
	t.Helper()
	conflict, ok := err.(*ActivationVersionConflictError)
	if !ok || conflict.ExpectedVersion != expected || conflict.ActualVersion != actual {
		t.Fatalf("version conflict = %#v, want expected=%d actual=%d", err, expected, actual)
	}
}

func TestActivationRecoverySelectionRaisesAndCompletionDisposesIncompleteStop(t *testing.T) {
	t.Parallel()
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	prepared, err := PrepareActivation(ActivationPrepareOptions{
		WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-stop-prepare",
		Request: ActivationPrepareRequest{KitIdentity: "kit:test", BackendIdentity: "backend:test", Phase: 21, Direction: "activate", Resources: activationRecoveryTestResources()},
	})
	if err != nil {
		t.Fatalf("PrepareActivation() error = %v", err)
	}
	selectionRequest := ActivationRecoveryRequest{
		ActivationID: prepared.Activation.ID, ExpectedVersion: prepared.Activation.Version, ExpectedStatus: "prepared", Direction: "restore-prior", Reason: "test interruption",
	}
	selectionOptions := ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-stop-select"}
	selected, err := RecoverActivation(selectionOptions, selectionRequest)
	if err != nil {
		t.Fatalf("RecoverActivation() error = %v", err)
	}
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	var stopID string
	if err := conn.QueryRow(`SELECT id FROM stops WHERE subject_type = 'activation' AND subject_id = ? AND code = 'activation-incomplete'`, prepared.Activation.ID).Scan(&stopID); err != nil {
		t.Fatalf("load activation-incomplete Stop: %v", err)
	}
	var attentionCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM attention_stop_links WHERE stop_id = ?`, stopID).Scan(&attentionCount); err != nil || attentionCount != 1 {
		t.Fatalf("activation Stop Attention links = %d, %v, want one", attentionCount, err)
	}
	selected.Activation = observeActivationResourcesForRecoveryTest(t, workspace, migrationsDir, selected.Activation, "old")
	completed, err := CompleteActivationRecovery(ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-stop-complete"}, ActivationRecoveryCompleteRequest{
		ActivationID: prepared.Activation.ID, ExpectedVersion: selected.Activation.Version, SelectedEventID: selected.EventID, ResourceVersions: activationResourceVersions(selected.Activation.Resources),
	})
	if err != nil {
		t.Fatalf("CompleteActivationRecovery() error = %v", err)
	}
	if completed.Activation.Status != "aborted" {
		t.Fatalf("completed Activation status = %q, want aborted", completed.Activation.Status)
	}
	replayed, err := RecoverActivation(selectionOptions, selectionRequest)
	if err != nil || !replayed.Replayed || replayed.EventID != selected.EventID || replayed.Activation.Status != "prepared" || replayed.Activation.Recovery == nil || replayed.Activation.Recovery.SelectedEventID != selected.EventID {
		t.Fatalf("RecoverActivation(replay after completion) = (%#v, %v), want original prepared selection identity", replayed, err)
	}
	var resolutions int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM stop_resolutions WHERE stop_id = ?`, stopID).Scan(&resolutions); err != nil || resolutions != 1 {
		t.Fatalf("activation Stop resolutions = %d, %v, want one", resolutions, err)
	}
	before := readWorkflowTestFile(t, workflowJournalPath(t, workspace))
	_, err = RecoverActivation(ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-terminal-select"}, ActivationRecoveryRequest{ActivationID: prepared.Activation.ID, ExpectedVersion: completed.Activation.Version, ExpectedStatus: "aborted", Direction: "restore-prior", Reason: "terminal retry"})
	var invalid *ActivationInvalidError
	if !errors.As(err, &invalid) {
		t.Fatalf("RecoverActivation(terminal) error = %v, want ActivationInvalidError", err)
	}
	if after := readWorkflowTestFile(t, workflowJournalPath(t, workspace)); string(after) != string(before) {
		t.Fatal("terminal recovery selection appended an Event")
	}
}

func TestCompleteActivationRecoveryRejectsMissingPriorResourceSet(t *testing.T) {
	t.Parallel()
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	prepared, err := PrepareActivation(ActivationPrepareOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-missing-prior-prepare", Request: ActivationPrepareRequest{KitIdentity: "kit:test", BackendIdentity: "backend:test", Phase: 21, Direction: "activate", Resources: activationRecoveryTestResources()}})
	if err != nil {
		t.Fatalf("PrepareActivation() error = %v", err)
	}
	selected, err := RecoverActivation(ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-missing-prior-select"}, ActivationRecoveryRequest{ActivationID: prepared.Activation.ID, ExpectedVersion: prepared.Activation.Version, ExpectedStatus: "prepared", Direction: "restore-prior", Reason: "test missing resource"})
	if err != nil {
		t.Fatalf("RecoverActivation() error = %v", err)
	}
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	if _, err := conn.Exec(`DELETE FROM activation_resources WHERE activation_id = ?`, selected.Activation.ID); err != nil {
		t.Fatalf("delete resource projection rows: %v", err)
	}
	_, err = CompleteActivationRecovery(ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-missing-prior-complete"}, ActivationRecoveryCompleteRequest{ActivationID: selected.Activation.ID, ExpectedVersion: selected.Activation.Version, SelectedEventID: selected.EventID, ResourceVersions: activationResourceVersions(selected.Activation.Resources)})
	var incomplete *ActivationIncompleteError
	if !errors.As(err, &incomplete) || incomplete.ActivationID != selected.Activation.ID || len(incomplete.Components) != len(activationResourceComponents) || len(incomplete.Reasons) != len(activationResourceComponents) {
		t.Fatalf("CompleteActivationRecovery(missing resources) error = %#v, want typed incomplete resources", err)
	}
}

func observeActivationResourcesForRecoveryTest(t *testing.T, workspace, migrationsDir string, activation Activation, state string) Activation {
	t.Helper()
	for _, resource := range activation.Resources {
		result, err := ObserveActivationResource(ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: "activation-resource-" + resource.Component}, ActivationResourceObservationRequest{ActivationID: activation.ID, Component: resource.Component, ExpectedActivationVersion: activation.Version, ExpectedVersion: resource.Version, ObservedIdentity: resource.OldIdentity, State: state, EvidenceRef: "test:" + resource.Component})
		if err != nil {
			t.Fatalf("ObserveActivationResource(%s) error = %v", resource.Component, err)
		}
		activation = result.Activation
	}
	return activation
}

func activationRecoveryTestResources() []ActivationResourceInput {
	return []ActivationResourceInput{
		{Component: "backend", NewIdentity: "backend:test"},
		{Component: "kit", NewIdentity: "kit:test"},
		{Component: "overlay", NewIdentity: "overlay:test"},
		{Component: "project", NewIdentity: "project:test"},
		{Component: "provider", NewIdentity: "provider:test"},
	}
}

func activationResourceVersions(resources []ActivationResource) []ActivationResourceVersion {
	versions := make([]ActivationResourceVersion, 0, len(resources))
	for _, resource := range resources {
		versions = append(versions, ActivationResourceVersion{Component: resource.Component, Version: resource.Version})
	}
	return versions
}
