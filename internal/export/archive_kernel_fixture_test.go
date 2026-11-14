package export

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/evidence"
	"github.com/mbell-nzl/fid/internal/issues"
	fidproject "github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/workflow"
)

// addKernelExportFacts creates exportable source and artifact
// facts through the normal journal writers.  In particular, the artifact is a
// retained gate-run receipt artifact, never a row inserted into the derived
// SQLite projection.
func addKernelExportFacts(t *testing.T, workspace string) string {
	t.Helper()
	migrationsDir := testMigrationsDir(t)
	project, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(workspace), "config.toml"))
	if err != nil {
		t.Fatalf("read project config: %v", err)
	}

	if _, err := workflow.Define(workflow.DefineOptions{
		WorkspaceRoot: workspace, MigrationsDir: migrationsDir, Actor: "agent:fixture", IdempotencyKey: "fixture-workflow-define",
		Request: workflow.DefinitionRequest{Name: "fixture-export", AppliesTo: "any", States: []string{"unstarted", "started", "completed", "cancelled"}, Categories: []string{"unstarted", "started", "completed", "cancelled"}, Transitions: []workflow.Transition{{From: "unstarted", To: "started"}, {From: "started", To: "completed"}}},
	}); err != nil {
		t.Fatalf("define fixture workflow: %v", err)
	}
	if _, err := workflow.Assign(workflow.AssignOptions{
		WorkspaceRoot: workspace, MigrationsDir: migrationsDir, Actor: "agent:fixture", IdempotencyKey: "fixture-workflow-assign",
		Request: workflow.AssignmentRequest{Name: "fixture-export", Project: true},
	}); err != nil {
		t.Fatalf("assign fixture workflow: %v", err)
	}
	issue, err := issues.Create(issues.CreateOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, Kind: "task", Title: "Archive fixture", Actor: "agent:fixture", IdempotencyKey: "fixture-issue"})
	if err != nil {
		t.Fatalf("create fixture issue: %v", err)
	}

	activation := activateExportFixture(t, workspace, migrationsDir)
	repository := fixtureRepository(t, workspace, project.ProjectID)
	registration, err := evidence.Register(evidence.RegisterOptions{
		WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "actor-builder", IdempotencyKey: "fixture-gate-register",
		Request: evidence.RegisterRequest{DefinitionID: "fixture-gate", DefinitionVersion: 1, WorktreeID: repository.WorktreeID, Scope: evidence.GateScope{IssueID: issue.ID}, Repository: evidenceRepository(repository)},
	})
	if err != nil {
		t.Fatalf("register fixture gate for activation %q: %v", activation.ID, err)
	}

	source := fixtureSource(t, workspace, repository)
	grant := workflow.AuthorityGrant{Actions: fixtureCapabilities()}
	grant.ActorSelector.ActorID, grant.ActorSelector.RequiredCapabilities = "actor-builder", fixtureCapabilities()
	grant.Scope.IssueID = issue.ID
	obligation := workflow.ObligationInput{ID: "OBL-fixture", Type: "gate-evidence", Applicability: "required", GuardedActions: []string{"envelope.complete", "issue.transition-terminal-success"}, AcceptedEvidenceTypes: []string{"receipt"}, AcceptedTrustSources: []string{"fid-executed"}, FreshnessTarget: "repository-tree", ArtifactPolicy: "required-available", Independence: "fid-executed", Source: source}
	obligation.Subject.Type, obligation.Subject.ID = "issue", issue.ID
	obligation.Registration.ID, obligation.Registration.Version = registration.Registration.ID, registration.Registration.Version
	harnessGrant := workflow.AuthorityGrant{Actions: fixtureHarnessCapabilities()}
	harnessGrant.ActorSelector.ActorID, harnessGrant.ActorSelector.RequiredCapabilities = "actor-harness", fixtureHarnessCapabilities()
	harnessGrant.Scope.IssueID = issue.ID
	request := workflow.EnvelopeCreateRequest{Outcome: "Archive semantic fixture", NonGoals: []string{}, RepositoryScope: repository, Sources: []workflow.Source{source}, PrimaryIssueID: issue.ID, RelatedIssueIDs: []string{}, AuthorityGrants: []workflow.AuthorityGrant{grant, harnessGrant}, Obligations: []workflow.ObligationInput{obligation}}
	request.RiskResult.Source, request.RiskResult.ResultDigest = source, "sha256:"+strings.Repeat("4", 64)
	created, err := workflow.CreateEnvelope(workflow.EnvelopeOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "actor-builder", IdempotencyKey: "fixture-envelope-create"}, request)
	if err != nil {
		t.Fatalf("create fixture envelope: %v", err)
	}
	if _, err := workflow.ActivateEnvelope(workflow.EnvelopeOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "actor-builder", IdempotencyKey: "fixture-envelope-activate"}, workflow.EnvelopeActivateRequest{EnvelopeID: created.Envelope.ID, ExpectedRevision: 1}); err != nil {
		t.Fatalf("activate fixture envelope: %v", err)
	}
	claim, err := claims.AcquireGeneric(claims.MutationOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "actor-builder", IdempotencyKey: "fixture-claim"}, claims.AcquireRequest{IssueID: issue.ID, EnvelopeID: created.Envelope.ID, EnvelopeRevision: 1, WorktreeID: repository.WorktreeID, Scope: claims.Scope{Paths: []string{"fixture-artifact.txt", "fixture-work.txt"}}, TTLSeconds: 3600})
	if err != nil {
		t.Fatalf("acquire fixture claim: %v", err)
	}
	run, err := workflow.StartRun(workflow.RunOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "actor-builder", IdempotencyKey: "fixture-run"}, workflow.RunStartRequest{IssueID: issue.ID, EnvelopeID: created.Envelope.ID, EnvelopeRevision: 1, ClaimID: claim.Claim.ID, ClaimVersion: claim.Claim.Version, Executor: workflow.Executor{ID: "executor-fixture", Source: "activated-harness"}, ExpectedRepository: fidproject.ExpectedRepository{HeadOID: repository.HeadOID, Dirty: repository.Dirty}})
	if err != nil {
		t.Fatalf("start fixture run: %v", err)
	}

	artifactData := []byte("artifact bytes\n")
	if err := os.WriteFile(filepath.Join(workspace, "fixture-artifact.txt"), artifactData, 0o600); err != nil {
		t.Fatalf("write fixture artifact input: %v", err)
	}
	gate, err := evidence.RunGate(context.Background(), evidence.GateRunOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "actor-builder", IdempotencyKey: "fixture-gate-run", Request: evidence.GateRunRequest{RegistrationID: registration.Registration.ID, RegistrationVersion: registration.Registration.Version, ObligationID: obligation.ID, RunID: run.Run.ID, ExpectedRepository: fidproject.ExpectedRepository{HeadOID: repository.HeadOID, Dirty: repository.Dirty}}})
	if err != nil {
		t.Fatalf("run fixture gate: %v", err)
	}
	if gate.Receipt == nil || len(gate.Receipt.Artifacts) != 1 {
		t.Fatalf("fixture receipt artifacts = %#v, want one retained artifact", gate.Receipt)
	}
	artifact := gate.Receipt.Artifacts[0]
	if artifact.Digest != digest(artifactData) || artifact.SizeBytes != int64(len(artifactData)) {
		t.Fatalf("fixture artifact = %#v, want digest %q and size %d", artifact, digest(artifactData), len(artifactData))
	}
	return "artifacts/" + strings.TrimPrefix(artifact.Digest, "sha256:")
}

// addRecoveredKernelExportFacts extends the ordinary archive fixture through
// a genuine terminated observation, released predecessor Claim, read-only
// reconciliation, and ordinary Run recovery.  It leaves the complete Event
// reference chain available to archive semantic tests.
func addRecoveredKernelExportFacts(t *testing.T, workspace string) string {
	t.Helper()
	artifactEntry := addKernelExportFacts(t, workspace)
	migrationsDir := testMigrationsDir(t)
	runs, err := workflow.ListRuns(workflow.RunListOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, Status: "active", Limit: 1})
	if err != nil || len(runs) != 1 {
		t.Fatalf("list fixture active Run: %#v, %v", runs, err)
	}
	prior := runs[0]
	observedAt := time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)
	observation, err := workflow.ObserveRun(workflow.RunOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "actor-harness", IdempotencyKey: "fixture-run-observe"}, workflow.RunObserveRequest{RunID: prior.ID, ExpectedVersion: prior.Version, Liveness: workflow.Liveness{ExecutorID: prior.Executor.ID, State: "terminated", SourceRef: "fixture:terminated", ObservedAt: observedAt}})
	if err != nil {
		t.Fatalf("observe fixture Run termination: %v", err)
	}
	claimVersion := fixtureClaimVersion(t, workspace, prior.ClaimID)
	if _, err := claims.ReleaseGeneric(claims.MutationOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "actor-builder", IdempotencyKey: "fixture-claim-release"}, claims.ReleaseRequest{ClaimID: prior.ClaimID, ExpectedVersion: claimVersion, Reason: "operator-released"}); err != nil {
		t.Fatalf("release fixture predecessor Claim: %v", err)
	}
	// Artifact material remains ignored. Recovery adopts an actual tracked
	// work change rather than using generated material as repository proof.
	if err := os.WriteFile(filepath.Join(workspace, "fixture-work.txt"), []byte("work to adopt\n"), 0o600); err != nil {
		t.Fatalf("edit fixture recovery work: %v", err)
	}
	inspection, err := workflow.InspectRun(workspace, migrationsDir, prior.ID, workspace)
	if err != nil {
		t.Fatalf("inspect fixture Run for recovery: %v", err)
	}
	if _, err := workflow.RecoverRun(workflow.RunOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "actor-builder", IdempotencyKey: "fixture-run-recover"}, workflow.RunRecoveryRequest{PriorRunID: prior.ID, ExpectedRunVersion: observation.Run.Version, ExpectedClaimVersion: claimVersion + 1, ObservationID: observation.Observation.ID, ReconciliationDigest: inspection.Reconciliation.Digest, AdoptedPaths: []string{"fixture-work.txt"}, Authority: workflow.RecoveryAuthority{Mode: "ordinary", ActorID: "actor-builder", Reason: "trusted termination"}, Executor: workflow.Executor{ID: "executor-recovered", Source: "activated-harness"}, TTLSeconds: 60}); err != nil {
		t.Fatalf("recover fixture Run: %v", err)
	}
	return artifactEntry
}

func fixtureClaimVersion(t *testing.T, workspace, claimID string) int {
	t.Helper()
	database, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("open fixture projection for Claim: %v", err)
	}
	defer func() { _ = database.Close() }()
	var version int
	if err := database.QueryRow(`SELECT version FROM claims WHERE id = ?`, claimID).Scan(&version); err != nil {
		t.Fatalf("read fixture Claim version: %v", err)
	}
	return version
}

func activateExportFixture(t *testing.T, workspace, migrationsDir string) workflow.Activation {
	t.Helper()
	resources := make([]workflow.ActivationResourceInput, 0, 5)
	for _, component := range []string{"backend", "kit", "overlay", "project", "provider"} {
		resources = append(resources, workflow.ActivationResourceInput{Component: component, NewIdentity: "fixture-" + component})
	}
	prepared, err := workflow.PrepareActivation(workflow.ActivationPrepareOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "fid-bootstrap", IdempotencyKey: "fixture-activation-prepare", Request: workflow.ActivationPrepareRequest{KitIdentity: "fixture", BackendIdentity: "fid", Phase: 21, Direction: "activate", Resources: resources}})
	if err != nil {
		t.Fatalf("prepare fixture activation: %v", err)
	}
	profiles, err := json.Marshal([]workflow.ActorProfile{
		{ActorID: "actor-builder", ActorType: "agent", Producer: "fid", ProvenanceRef: "fixture:activation", Capabilities: fixtureCapabilities()},
		{ActorID: "actor-harness", ActorType: "agent", Producer: "activated-harness", ProvenanceRef: "fixture:activation", Capabilities: fixtureHarnessCapabilities()},
	})
	if err != nil {
		t.Fatalf("encode fixture authority profile: %v", err)
	}
	profile, err := workflow.StageActivationComponent(workflow.ActivationComponentOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "fid-bootstrap", IdempotencyKey: "fixture-activation-profile", Request: workflow.ActivationComponentRequest{ActivationID: prepared.Activation.ID, ComponentType: "authority-profile", ComponentID: "fixture-authority", ComponentVersion: 1, Content: profiles, ExpectedVersion: prepared.Activation.Version}})
	if err != nil {
		t.Fatalf("stage fixture authority profile: %v", err)
	}
	definition := []map[string]any{{"artifact_bytes": 1024, "artifact_paths": []string{"fixture-artifact.txt"}, "argv": []string{"/bin/true"}, "credential_names": []string{}, "cwd": ".", "definition_id": "fixture-gate", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}}
	definitions, err := json.Marshal(definition)
	if err != nil {
		t.Fatalf("encode fixture gate definition: %v", err)
	}
	gates, err := workflow.StageActivationComponent(workflow.ActivationComponentOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "fid-bootstrap", IdempotencyKey: "fixture-activation-gates", Request: workflow.ActivationComponentRequest{ActivationID: prepared.Activation.ID, ComponentType: "gate-definitions", ComponentID: "fixture-gates", ComponentVersion: 1, Content: definitions, ExpectedVersion: profile.Activation.Version}})
	if err != nil {
		t.Fatalf("stage fixture gate definitions: %v", err)
	}
	activation := gates.Activation
	for _, resource := range activation.Resources {
		identity := resource.NewIdentity
		observed, err := workflow.ObserveActivationResource(workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "fid-bootstrap", IdempotencyKey: "fixture-resource-" + resource.Component}, workflow.ActivationResourceObservationRequest{ActivationID: activation.ID, Component: resource.Component, ExpectedActivationVersion: activation.Version, ExpectedVersion: resource.Version, ObservedIdentity: &identity, State: "new", EvidenceRef: "fixture:" + resource.Component})
		if err != nil {
			t.Fatalf("observe fixture resource %q: %v", resource.Component, err)
		}
		activation = observed.Activation
	}
	digests := []string{profile.Component.Digest, gates.Component.Digest}
	sort.Strings(digests)
	committed, err := workflow.CommitActivation(workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "fid-bootstrap", IdempotencyKey: "fixture-activation-commit"}, workflow.ActivationCommitRequest{ActivationID: activation.ID, ExpectedStatus: "staged", ComponentDigests: digests, ExpectedVersion: activation.Version})
	if err != nil {
		t.Fatalf("commit fixture activation: %v", err)
	}
	finalized, err := workflow.FinalizeActivation(workflow.ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, BootstrapActor: "fid-bootstrap", IdempotencyKey: "fixture-activation-finalize"}, workflow.ActivationFinalizeRequest{ActivationID: activation.ID, ExpectedStatus: "committed", ExpectedVersion: committed.Activation.Version})
	if err != nil {
		t.Fatalf("finalize fixture activation: %v", err)
	}
	return finalized.Activation
}

func fixtureCapabilities() []string {
	return []string{"artifact.verify", "claim.acquire", "claim.release", "claim.renew", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "issue.transition", "run.close", "run.recover", "run.start"}
}

func fixtureHarnessCapabilities() []string { return []string{"run.observe"} }

func fixtureRepository(t *testing.T, workspace, projectID string) workflow.Repository {
	t.Helper()
	fixtureGit(t, workspace, "init", "-q")
	fixtureGit(t, workspace, "config", "user.email", "fixture@example.test")
	fixtureGit(t, workspace, "config", "user.name", "Archive fixture")
	if err := os.MkdirAll(filepath.Join(workspace, "project-docs", "gates"), 0o755); err != nil {
		t.Fatalf("create fixture Source directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workspace, ".gitignore"), []byte("fixture-artifact.txt\n"), 0o644); err != nil {
		t.Fatalf("write fixture gitignore: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "project-docs", "gates", "fixture.json"), []byte("{\"fixture\":true}\n"), 0o644); err != nil {
		t.Fatalf("write fixture Source prose: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "fixture-work.txt"), []byte("original work\n"), 0o600); err != nil {
		t.Fatalf("write fixture recovery work: %v", err)
	}
	fixtureGit(t, workspace, "add", ".gitignore", config.ProjectIdentityFileName, "project-docs/gates/fixture.json", "fixture-work.txt")
	fixtureGit(t, workspace, "commit", "-qm", "fixture source")
	observed, err := fidproject.ObserveRepository(workspace)
	if err != nil {
		t.Fatalf("capture fixture repository: %v", err)
	}
	base := observed.HeadOID
	return workflow.Repository{BaseOID: &base, CaptureSource: "git", Dirty: observed.Dirty, DirtyDigest: observed.DirtyDigest, HeadOID: observed.HeadOID, ProjectID: projectID, TreeOID: observed.TreeOID, WorktreeID: observed.WorktreeID}
}

func fixtureSource(t *testing.T, workspace string, repository workflow.Repository) workflow.Source {
	path := "project-docs/gates/fixture.json"
	t.Helper()
	data, err := os.ReadFile(filepath.Join(workspace, path))
	if err != nil {
		t.Fatal(err)
	}
	return workflow.Source{Digest: digest(data), DisplayLabel: "Fixture gate policy", Path: &path, Repository: repository, SourceType: "project-file"}
}

func evidenceRepository(repository workflow.Repository) evidence.Repository {
	return evidence.Repository{BaseOID: repository.BaseOID, CaptureSource: repository.CaptureSource, Dirty: repository.Dirty, DirtyDigest: repository.DirtyDigest, HeadOID: repository.HeadOID, ProjectID: repository.ProjectID, TreeOID: repository.TreeOID, WorktreeID: repository.WorktreeID}
}
