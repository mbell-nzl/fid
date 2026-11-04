//go:build integration

package integration

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

type p22HumanInspection struct {
	CWD  string   `json:"cwd"`
	Args []string `json:"args"`
}

type p22HumanCheckpoint struct {
	Name           string               `json:"name"`
	Project        string               `json:"project"`
	IDs            map[string]string    `json:"ids"`
	LeaseExpiresAt string               `json:"lease_expires_at,omitempty"`
	Inspections    []p22HumanInspection `json:"inspections"`
}

type p22HumanManifest struct {
	Version      int                  `json:"version"`
	GeneratedAt  string               `json:"generated_at"`
	SourceCommit string               `json:"source_commit"`
	SourceDirty  bool                 `json:"source_dirty"`
	Binary       string               `json:"binary"`
	BinarySHA256 string               `json:"binary_sha256"`
	Checkpoints  []p22HumanCheckpoint `json:"checkpoints"`
}

// Generation happens in a child test process. The parent only inspects once
// that process (including TestMain cleanup) has exited successfully.
func TestP22HumanFixtures(t *testing.T) {
	t.Parallel()
	if root := os.Getenv("FID_P22_HUMAN_FIXTURE_CHILD"); root != "" {
		p22GenerateHumanFixtures(t, root)
		return
	}
	root := os.Getenv("FID_P22_HUMAN_FIXTURE_ROOT")
	if root == "" {
		root = t.TempDir()
	}
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	fresh, err := os.MkdirTemp(root, "checkpoints-")
	if err != nil {
		t.Fatal(err)
	}
	fresh, err = filepath.Abs(fresh)
	if err != nil {
		t.Fatal(err)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	child := exec.Command(executable, "-test.run=^TestP22HumanFixtures$", "-test.v")
	child.Env = append(os.Environ(), "FID_P22_HUMAN_FIXTURE_CHILD="+fresh)
	child.Dir = integrationRepositoryRoot()
	if output, err := child.CombinedOutput(); err != nil {
		t.Fatalf("fixture generator child: %v\n%s", err, output)
	}
	manifestPath := filepath.Join(fresh, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest p22HumanManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	binary, err := os.ReadFile(manifest.Binary)
	if err != nil || fmt.Sprintf("sha256:%x", sha256.Sum256(binary)) != manifest.BinarySHA256 {
		t.Fatalf("retained binary unavailable or changed: %v", err)
	}
	if len(manifest.Checkpoints) != 4 {
		t.Fatalf("retained checkpoint count = %d, want 4", len(manifest.Checkpoints))
	}
	for _, checkpoint := range manifest.Checkpoints {
		if checkpoint.Project == "" || len(checkpoint.IDs) == 0 || len(checkpoint.Inspections) < 2 {
			t.Fatalf("incomplete retained checkpoint: %#v", checkpoint)
		}
		for _, inspection := range checkpoint.Inspections {
			stdout, stderr, err := closureRunCommand(inspection.CWD, manifest.Binary, inspection.Args...)
			if err != nil || stderr != "" || strings.TrimSpace(stdout) == "" {
				t.Fatalf("post-exit inspection %s %v: %v stdout=%q stderr=%q", checkpoint.Name, inspection.Args, err, stdout, stderr)
			}
		}
	}
	t.Logf("Retained Human fixture manifest: %s", manifestPath)
}

func p22GenerateHumanFixtures(t *testing.T, root string) {
	t.Helper()
	existing, err := os.ReadDir(root)
	if err != nil || len(existing) != 0 {
		t.Fatalf("generator requires a fresh empty directory: %v", err)
	}
	contents, err := os.ReadFile(integrationFIDBinary(t))
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(root, "fid")
	if err := os.WriteFile(binary, contents, 0o700); err != nil {
		t.Fatal(err)
	}
	commit, _ := closureRun(t, integrationRepositoryRoot(), "git", "rev-parse", "HEAD")
	status, _ := closureRun(t, integrationRepositoryRoot(), "git", "status", "--porcelain")
	manifest := p22HumanManifest{Version: 1, GeneratedAt: time.Now().UTC().Format(time.RFC3339Nano), SourceCommit: strings.TrimSpace(commit), SourceDirty: status != "", Binary: binary, BinarySHA256: fmt.Sprintf("sha256:%x", sha256.Sum256(contents))}
	workspace := p22HumanDirectory(t, root, "recovery")
	a := p22NewFIDRecoveryAdapterAt(t, binary, workspace, "agent", 86400, []string{}, nil)
	prior := a.runID
	p22Write(t, workspace, "base.txt", "retained adopted work\n")
	p22AcceptedCall(t, a, "fixture-observe", "run.observe", nil)
	p22AcceptedCall(t, a, "fixture-inspect", "run.inspect", nil)
	p22AcceptedCall(t, a, "fixture-checkpoint", "run.checkpoint", nil)
	p22AcceptedCall(t, a, "fixture-expire", "claim.expire", nil)
	recovered := p22AcceptedCall(t, a, "fixture-recover", "run.recover", map[string]any{"adopted_paths": []string{"base.txt"}, "ttl_seconds": 86400})
	p22AcceptedCall(t, a, "fixture-budget", "run.observe-budget", map[string]any{"value": 0})
	resume := p21Run(t, binary, workspace, "--json", "resume", a.issueID, "--repository", workspace)
	human, _ := closureRun(t, workspace, binary, "--project", workspace, "resume", a.issueID, "--repository", workspace)
	p22AssertHumanRecovery(t, human, p21Path(t, recovered, "result", "recovery").(map[string]any), p21Path(t, resume, "result", "resume", "budgets").([]any))
	p22AssertHumanResume(t, human, p21Path(t, resume, "result", "resume").(map[string]any), false)
	checkpoint := p22HumanRunCheckpoint(a, "Recovery checkpoint")
	checkpoint.IDs["prior_run_id"] = prior
	checkpoint.LeaseExpiresAt = p21PathString(t, resume, "result", "resume", "claim", "expires_at")
	manifest.Checkpoints = append(manifest.Checkpoints, checkpoint)
	for _, fixture := range []struct{ name, boundary, direction string }{{"Activation checkpoint A", "prepared", "restore-prior"}, {"Activation checkpoint B", "committed", "complete-forward"}} {
		workspace := p22HumanDirectory(t, root, fixture.boundary)
		_, id := p22ActivationBoundaryFixtureAt(t, binary, workspace, fixture.boundary)
		activation := p21ActivationStatus(t, workspace, id)
		input := p21Input(t, workspace, "fixture-select", map[string]any{"activation_id": id, "expected_version": activation.Version, "expected_status": activation.Status, "direction": fixture.direction, "reason": "resume interrupted installation"})
		p21Run(t, binary, workspace, "--json", "activation", "recover", "--input", input, "--idempotency-key", "fixture-select")
		human, _ := closureRun(t, workspace, binary, "--project", workspace, "activation", "status", id)
		if !strings.Contains(human, fixture.direction) {
			t.Fatalf("activation Human view lacks selected direction: %q", human)
		}
		project, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(workspace), "config.toml"))
		if err != nil {
			t.Fatal(err)
		}
		manifest.Checkpoints = append(manifest.Checkpoints, p22HumanCheckpoint{Name: fixture.name, Project: workspace, IDs: map[string]string{"activation_id": id, "project_id": project.ProjectID}, Inspections: p22HumanInspections(workspace, []string{"activation", "status", id})})
	}
	source := p22HumanDirectory(t, root, "import-source")
	original := p22NewFIDRecoveryAdapterAt(t, binary, source, "agent", 86400, []string{}, nil)
	archive := filepath.Join(root, "import.zip")
	p21Run(t, binary, source, "--json", "export", "create", archive)
	p21Run(t, binary, root, "--json", "export", "validate", archive)
	target := filepath.Join(root, "imported")
	closureRun(t, root, "git", "clone", "-q", source, target)
	p21Run(t, binary, target, "--json", "export", "import", archive, "--repository", target, "--idempotency-key", "fixture-import")
	imported := *original
	imported.workspace = target
	checkpoint = p22HumanRunCheckpoint(&imported, "Imported project checkpoint")
	checkpoint.IDs["source_project_path"] = source
	checkpoint.IDs["archive_path"] = archive
	inspected := p21Run(t, binary, target, "--json", "run", "inspect", imported.runID, "--repository", target)
	if p21Path(t, inspected, "result", "liveness") != nil {
		t.Fatal("import manufactured a destination liveness observation")
	}
	manifest.Checkpoints = append(manifest.Checkpoints, checkpoint)
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "manifest.json")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write(append(data, '\n')); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func p22HumanDirectory(t *testing.T, root, name string) string {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}

func p22HumanRunCheckpoint(a *p22FIDRecoveryAdapter, name string) p22HumanCheckpoint {
	return p22HumanCheckpoint{Name: name, Project: a.workspace, IDs: map[string]string{"project_id": a.projectID, "issue_id": a.issueID, "run_id": a.runID, "claim_id": a.claimID}, Inspections: p22HumanInspections(a.workspace, []string{"resume", a.issueID, "--repository", a.workspace}, []string{"run", "inspect", a.runID, "--repository", a.workspace})}
}

func p22HumanInspections(workspace string, leaves ...[]string) []p22HumanInspection {
	inspections := []p22HumanInspection{}
	for _, leaf := range leaves {
		for _, jsonMode := range []bool{false, true} {
			args := []string{"--project", workspace}
			if jsonMode {
				args = append(args, "--json")
			}
			inspections = append(inspections, p22HumanInspection{CWD: workspace, Args: append(args, leaf...)})
		}
	}
	return inspections
}
