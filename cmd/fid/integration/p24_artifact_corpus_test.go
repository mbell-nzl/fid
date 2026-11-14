//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

var p24ArtifactFields = map[string]string{"artifact.verify": "artifact_id"}

func leanArtifactFixture(t *testing.T, binary, workspace string, bodies ...string) string {
	t.Helper()
	leanFlagFixture(t, binary, workspace)
	p24InstallLeanPolicy(t, workspace, true)
	path := filepath.Join(workspace, ".fid.toml")
	raw := readP21File(t, path)
	raw = []byte(strings.Replace(string(raw), "[checks.pass]\nrun=['true']", "[checks.pass]\nartifacts=['artifact.txt']\nrun=['true']", 1))
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	body := "phase 24 artifact\n"
	if len(bodies) > 0 {
		body = bodies[0]
	}
	if err := os.WriteFile(filepath.Join(workspace, "artifact.txt"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	closureRun(t, workspace, "git", "add", ".")
	closureRun(t, workspace, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Artifact fixture")
	closureRun(t, workspace, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, workspace, append([]string{"--json"}, args...)...)
	}
	created := call("issue", "create", "--kind", "task", "--title", "Artifact task", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "artifact-task")
	taken := call("takeover", p21PathString(t, created, "result", "issue", "id"), "--actor", "human:mark", "--reason", "Retain file")
	hold := p21PathString(t, taken, "result", "hold_id")
	call("start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "artifact-start")
	checked := call("check", "pass", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "artifact-check")
	ids := p21Path(t, checked, "result", "checks").([]any)[0].(map[string]any)["artifact_ids"].([]any)
	if len(ids) != 1 {
		t.Fatalf("artifact IDs=%v", ids)
	}
	return ids[0].(string)
}

func TestP24ArtifactVerifyKeyedLeafCorpus(t *testing.T) {
	t.Parallel()
	binary, workspace := integrationFIDBinary(t), t.TempDir()
	id := leanArtifactFixture(t, binary, workspace)
	conn := p22FaultDB(t, workspace)
	var digest string
	p22FaultOK(t, conn.QueryRow("SELECT digest FROM referenced_artifacts WHERE id=?", id).Scan(&digest))
	p22FaultOK(t, conn.Close())
	if err := os.Remove(filepath.Join(config.LedgerDir(workspace), "artifacts", strings.TrimPrefix(digest, "sha256:"))); err != nil {
		t.Fatal(err)
	}
	p24StartKeyedCapture(workspace)
	input := p21Input(t, workspace, "verify", map[string]any{"artifact_id": id})
	p21Run(t, binary, workspace, "--project", workspace, "--json", "artifact", "verify", "--input", input, "--actor", "agent:builder", "--idempotency-key", "artifact-verify")
	entry, found := p24EndKeyedCapture(workspace)["artifact.verify"]
	if !found {
		t.Fatal("artifact.verify was not captured")
	}
	p24ExerciseCapturedLeaf(t, workspace, binary, "artifact.verify", p24ArtifactFields["artifact.verify"], entry)
	p24AssertFamilyConflict(t, workspace, binary, "artifact.verify", entry, []string{"artifact_id", "--actor"})
	p24AssertFamilyConflict(t, workspace, binary, "artifact.verify", entry, []string{"--actor"})
	assertLeanIntegrationJournal(t, workspace)
	p24AssertRebuildMapParity(t, workspace, binary)
}
