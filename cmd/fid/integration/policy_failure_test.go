//go:build integration

package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

// Policy refusals must have process-level text/JSON coverage in the
// registered-code corpus, including their closed details and repair steps.
func TestPolicyUnavailableHumanAndJSON(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"unknown-default", "invalid-file", "deleted-default"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root, binary := gittest.NewWorkspace(t), integrationFIDBinary(t)
			p21Run(t, binary, root, "init", "--prefix", "POL", "--json")
			git := func(args ...string) string {
				t.Helper()
				out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput()
				if err != nil {
					t.Fatalf("git %v: %v: %s", args, err, out)
				}
				return strings.TrimSpace(string(out))
			}
			git("remote", "add", "origin", filepath.Join(t.TempDir(), "unused.git"))
			problem, action := "origin's default branch is unknown", "git remote set-head origin <branch>"
			var ref, commit any
			if name == "invalid-file" {
				file, err := os.OpenFile(filepath.Join(root, ".fid.toml"), os.O_APPEND|os.O_WRONLY, 0)
				if err != nil {
					t.Fatal(err)
				}
				_, err = file.WriteString("\n[lanes.standard]\npaths = []\n")
				closeErr := file.Close()
				if err != nil || closeErr != nil {
					t.Fatalf("write invalid policy: %v, close: %v", err, closeErr)
				}
				git("add", ".fid.toml")
				git("-c", "user.name=fid test", "-c", "user.email=fid@example.test", "-c", "commit.gpgsign=false", "commit", "-q", "--no-verify", "-m", "invalid policy")
				git("update-ref", "refs/remotes/origin/trunk", "HEAD")
				git("symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
				ref, commit = "origin/trunk", git("rev-parse", "HEAD")
				problem, action = "invalid .fid.toml: lanes.standard.paths is refused: other paths are Standard", "fix .fid.toml on the default branch, merge it and git fetch"
			}
			if name == "deleted-default" {
				git("add", ".fid.toml")
				git("-c", "user.name=fid test", "-c", "user.email=fid@example.test", "-c", "commit.gpgsign=false", "commit", "-q", "--no-verify", "-m", "valid policy")
				git("update-ref", "refs/remotes/origin/trunk", "HEAD")
				git("symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
				git("update-ref", "-d", "refs/remotes/origin/trunk")
				ref, problem = "origin/trunk", "origin's default branch origin/trunk does not exist"
			}
			args := []string{"next", "--claim", "--actor", "agent:builder"}
			out, stderr, err := closureRunCommand(root, binary, append([]string{"--json"}, args...)...)
			if err == nil || processExitCode(t, err) != 4 || out != "" {
				t.Fatalf("JSON policy refusal: %v %q %q", err, out, stderr)
			}
			reply := p23Envelope(t, stderr)
			failure := reply["error"].(map[string]any)
			wantDetails := map[string]any{"problem": problem, "ref": ref, "commit": commit}
			if reply["operation"] != "next.claim" || failure["code"] != "policy-unavailable" || failure["category"] != "precondition" || failure["retryable"] != false || failure["subject"] != nil || failure["message"] != problem || !reflect.DeepEqual(failure["details"], wantDetails) {
				t.Fatalf("policy refusal = %#v, want closed details %#v", reply, wantDetails)
			}
			sourceRef, sourceCommit := "null", "null"
			if ref != nil {
				sourceRef = ref.(string)
			}
			if commit != nil {
				sourceCommit = commit.(string)
			}
			reason := fmt.Sprintf("reason: policy-unavailable; commit=%s, problem=%s, ref=%s", sourceCommit, problem, sourceRef)
			want := "error: " + problem + "\n" + p24PairReasonCut(reason) + "\nnext: " + action + "\n"
			out, stderr, err = closureRunCommand(root, binary, args...)
			if err == nil || processExitCode(t, err) != 4 || out != "" || stderr != want {
				t.Fatalf("human policy refusal: %v stdout=%q stderr=%q want=%q", err, out, stderr, want)
			}
		})
	}
}
