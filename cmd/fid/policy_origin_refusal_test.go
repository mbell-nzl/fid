package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInvalidOriginPolicyKeySpellingNeverWrites(t *testing.T) {
	t.Parallel()
	for _, contents := range []string{
		"[lanes.assured]\nAUTO_MERGE = true\n",
		"[lanes.standard]\nPATHS = ['src/**']\n",
		"[lanes.standard]\nCHECKS = []\n",
		"[lanes.routine]\ntoken_budget = 1\nTOKEN_BUDGET = 2\n",
	} {
		t.Run(contents, func(t *testing.T) {
			t.Parallel()
			root, _ := policyWorkspace(t)
			path := filepath.Join(root, ".fid.toml")
			valid, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			appendPolicy(t, root, "\n"+contents)
			policyOrigin(t, root)
			commit := strings.TrimSpace(identityGit(t, root, "rev-parse", "origin/trunk"))
			if err := os.WriteFile(path, valid, 0o600); err != nil {
				t.Fatal(err)
			}
			before := readAllPolicyJournal(t, root)
			result := executeCommandForTest(root, "--json", "next", "--claim", "--actor", "agent:builder")
			if result.err == nil || !strings.Contains(result.stderr, `"code":"invalid-argument"`) || !strings.Contains(result.stderr, commit) {
				t.Errorf("invalid origin policy refusal = %#v, want commit %s", result, commit)
			}
			if after := readAllPolicyJournal(t, root); before != after {
				t.Error("invalid origin policy wrote a snapshot or hold to the journal")
			}
			if got := policySnapshots(t, root); len(got) != 0 {
				t.Errorf("invalid origin policy projected snapshots: %#v", got)
			}
		})
	}
}
