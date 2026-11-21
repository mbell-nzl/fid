//go:build integration

package integration

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestWorkflowHardFloorPublicReceipt(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, path string }{
		{"newline", ".github/workflows/odd\n.yml"},
		{"backslash", ".github/workflows/odd\\.yml"},
		{"invalid UTF-8", ".github/workflows/odd\xff.yml"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, binary := gittest.NewWorkspace(t), integrationFIDBinary(t)
			call := func(args ...string) map[string]any {
				t.Helper()
				return p21Run(t, binary, root, append([]string{"--json"}, args...)...)
			}
			call("init", "--prefix", "FID")
			p22Write(t, root, ".fid.toml", string(readP21File(t, filepath.Join(root, ".fid.toml")))+`
[lanes.routine]
checks=[]
[lanes.assured]
checks=['pass']
[checks.pass]
run=['true']
timeout='1m'
[review]
test_adequacy=[]
`)
			closureRun(t, root, "git", "add", ".")
			closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Policy base")
			closureRun(t, root, "git", "remote", "add", "origin", root)
			closureRun(t, root, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
			closureRun(t, root, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
			task := p21PathString(t, call("issue", "create", "--kind", "task", "--title", "Workflow hard floor", "--lane", "routine", "--actor", "human:mark", "--idempotency-key", "create"), "result", "issue", "id")
			hold := p21PathString(t, call("takeover", task, "--actor", "human:mark", "--reason", "Build the journey"), "result", "hold_id")
			call("start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "start")
			p22Write(t, root, tc.path, "name: workflow\n")
			closureRun(t, root, "git", "add", ".")
			closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Workflow change")
			checked := call("check", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "check")
			checks := p21Path(t, checked, "result", "checks").([]any)
			if len(checks) != 1 || checks[0].(map[string]any)["name"] != "pass" || checks[0].(map[string]any)["counts"] != true {
				t.Fatalf("floor checks=%v", checks)
			}
			if lane := p21PathString(t, call("show", task), "result", "issue", "lane"); lane != "routine" {
				t.Fatalf("check changed lane: %s", lane)
			}
			call("handoff", "--hold", hold, "--actor", "human:mark", "--pr", "https://github.com/owner/repo/pull/1", "--idempotency-key", "handoff")
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					call("rebuild")
				}
				shown := call("show", task)
				if lane := p21PathString(t, shown, "result", "issue", "lane"); lane != "assured" {
					t.Errorf("rebuild=%v: task lane=%s", rebuild, lane)
				}
				want := []string{string([]rune(tc.path))}
				for _, path := range [][]string{{"result", "floor"}, {"result", "handoff", "floor"}} {
					floor := p21Path(t, shown, path...).(map[string]any)
					if floor["lane"] != "assured" || floor["hard_floor"] != "assured" || !slices.Equal(p21StringSlice(t, floor["paths"], "paths"), want) || !slices.Equal(p21StringSlice(t, floor["hard_paths"], "hard_paths"), want) {
						t.Errorf("rebuild=%v: %s=%v, want Assured floor and hard floor with %q", rebuild, strings.Join(path, "."), floor, want)
					}
				}
			}
		})
	}
}
