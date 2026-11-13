package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const policyExample = `[project]
id = "project-0123456789abcdef0123456789abcdef"
[checks.test]
run = ["go", "test", "./..."]
timeout = "15m"
defined_by = ["Makefile", ".github/scripts/**"]
when_changed = ["**/*.go"]
[lanes.routine]
paths = ["docs/**"]
checks = []
rework_cap = 4
token_budget = 2_000_000
auto_merge = true
[lanes.standard]
checks = ["test"]
[lanes.assured]
paths = ["migrations/**"]
checks = ["test"]
[review]
test_adequacy = ["**/*_test.go", "**/testdata/**"]
`

func TestProjectIdentityAcceptsPolicyTables(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ProjectIdentityFileName), []byte(policyExample), 0o600); err != nil {
		t.Fatal(err)
	}
	id, err := ReadProjectIdentity(root)
	if err != nil || id != "project-0123456789abcdef0123456789abcdef" {
		t.Fatalf("identity with policy = %q, %v", id, err)
	}
	if err := EnsureProjectIdentity(root, id); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, ProjectIdentityFileName))
	if err != nil || string(data) != policyExample {
		t.Fatalf("init changed existing policy: %q, %v", data, err)
	}
}

func TestDecodeProjectPolicy(t *testing.T) {
	t.Parallel()
	got, err := DecodeProject([]byte(policyExample))
	if err != nil {
		t.Fatal(err)
	}
	check := got.Checks["test"]
	if strings.Join(check.Run, " ") != "go test ./..." || check.Timeout != 15*time.Minute ||
		len(check.DefinedBy) != 2 || len(check.WhenChanged) != 1 {
		t.Fatalf("check = %#v", check)
	}
	routine := got.Lanes["routine"]
	if !routine.ChecksSet || len(routine.Checks) != 0 || routine.ReworkCap != 4 || routine.TokenBudget != 2_000_000 || !routine.AutoMerge || len(routine.Paths) != 1 {
		t.Fatalf("routine = %#v", routine)
	}
	if got.Lanes["standard"].Checks[0] != "test" || got.Lanes["assured"].Paths[0] != "migrations/**" || len(got.TestAdequacy) != 2 {
		t.Fatalf("policy = %#v", got)
	}
	defaults, err := DecodeProject(nil)
	if err != nil {
		t.Fatal(err)
	}
	for name, lane := range defaults.Lanes {
		if lane.ChecksSet || lane.ReworkCap != 3 || lane.TokenBudget != 0 || lane.AutoMerge || len(lane.Paths) != 0 {
			t.Errorf("%s defaults = %#v", name, lane)
		}
	}
}

func TestDecodeProjectPolicyRefusals(t *testing.T) {
	t.Parallel()
	for _, test := range []struct{ name, contents, problem string }{
		{"unknown root", "version = 1", "unknown"},
		{"unknown table", "[actors]\nagent = true", "unknown"},
		{"unknown nested", "[checks.test]\nrun = ['true']\ntimeout = '1m'\nextra = 1", "unknown"},
		{"wrong type", "[checks.test]\nrun = 'true'\ntimeout = '1m'", "run"},
		{"missing timeout", "[checks.test]\nrun = ['true']", "timeout"},
		{"bad duration", "[checks.test]\nrun = ['true']\ntimeout = 'forever'", "timeout"},
		{"zero duration", "[checks.test]\nrun = ['true']\ntimeout = '0s'", "timeout"},
		{"negative duration", "[checks.test]\nrun = ['true']\ntimeout = '-1m'", "timeout"},
		{"empty run", "[checks.test]\nrun = []\ntimeout = '1m'", "run"},
		{"missing run", "[checks.test]\ntimeout = '1m'", "run"},
		{"empty executable", "[checks.test]\nrun = ['']\ntimeout = '1m'", "run"},
		{"undefined check", "[lanes.routine]\nchecks = ['missing']", "missing"},
		{"standard paths", "[lanes.standard]\npaths = []", "paths"},
		{"assured auto merge", "[lanes.assured]\nauto_merge = true", "auto_merge"},
		{"assured false auto merge", "[lanes.assured]\nauto_merge = false", "auto_merge"},
		{"unknown lane", "[lanes.fast]\nchecks = []", "lane"},
		{"bad defining glob", "[checks.test]\nrun = ['true']\ntimeout = '1m'\ndefined_by = ['../Makefile']", "defined_by"},
		{"bad conditional glob", "[checks.test]\nrun = ['true']\ntimeout = '1m'\nwhen_changed = ['/src/**']", "when_changed"},
		{"bad lane glob", "[lanes.routine]\npaths = ['[x]']", "paths"},
		{"bad review glob", "[review]\ntest_adequacy = ['']", "test_adequacy"},
		{"negative rework", "[lanes.routine]\nrework_cap = -1", "rework_cap"},
		{"negative budget", "[lanes.routine]\ntoken_budget = -1", "token_budget"},
		{"noninteger budget", "[lanes.routine]\ntoken_budget = 2.5", "token_budget"},
		{"invalid identity", "[project]\nid = 'bad'", "project.id"},
		{"missing identity", "[project]", "project.id"},
		{"duplicate key", "[lanes.routine]\nchecks = []\nchecks = []", "checks"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := DecodeProject([]byte(test.contents))
			if err == nil || !strings.Contains(err.Error(), test.problem) {
				t.Fatalf("decode = %v, want %q", err, test.problem)
			}
		})
	}
}

func TestDecodeProjectPolicyKeySpelling(t *testing.T) {
	t.Parallel()
	for _, key := range []string{
		"project", "id", "checks", "run", "timeout", "defined_by", "when_changed",
		"lanes", "routine", "standard", "assured", "paths", "rework_cap", "token_budget", "auto_merge",
		"review", "test_adequacy",
	} {
		t.Run(key, func(t *testing.T) {
			contents := strings.ReplaceAll(policyExample, key, strings.ToUpper(key))
			if _, err := DecodeProject([]byte(contents)); err == nil {
				t.Fatalf("accepted schema key %q", strings.ToUpper(key))
			}
		})
	}
	for _, test := range []struct{ name, contents string }{
		{"assured auto merge", "[lanes.assured]\nAUTO_MERGE = true"},
		{"standard paths", "[lanes.standard]\nPATHS = ['src/**']"},
		{"explicit checks", "[lanes.standard]\nCHECKS = []"},
		{"case colliding keys", "[lanes.routine]\ntoken_budget = 1\nTOKEN_BUDGET = 2"},
		{"case colliding tables", "[lanes.routine]\ntoken_budget = 1\n[LANES.routine]\ntoken_budget = 2"},
		{"inline table", "lanes = { assured = { AUTO_MERGE = true } }"},
		{"dotted key", "lanes.assured.AUTO_MERGE = true"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := DecodeProject([]byte(test.contents)); err == nil {
				t.Fatalf("accepted non-schema spelling: %s", test.contents)
			}
		})
	}
}

func TestDecodeProjectPreservesCheckNameCase(t *testing.T) {
	t.Parallel()
	file, err := DecodeProject([]byte("[checks.Unit]\nrun = ['upper']\ntimeout = '1m'\n[checks.unit]\nrun = ['lower']\ntimeout = '2m'\n[lanes.standard]\nchecks = ['Unit', 'unit']"))
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Checks) != 2 || file.Checks["Unit"].Run[0] != "upper" || file.Checks["unit"].Run[0] != "lower" || !file.Lanes["standard"].ChecksSet {
		t.Fatalf("case-sensitive check names = %#v", file)
	}
}

func TestDecodeProjectTestAdequacyPresence(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, text string
		set        bool
		paths      int
	}{
		{"missing", "", false, 0}, {"empty table", "[review]", false, 0},
		{"empty", "[review]\ntest_adequacy = []", true, 0},
		{"paths", "[review]\ntest_adequacy = ['**/*_test.go']", true, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			file, err := DecodeProject([]byte(test.text))
			if err != nil {
				t.Fatal(err)
			}
			if file.TestAdequacySet != test.set || len(file.TestAdequacy) != test.paths {
				t.Fatalf("review policy = %#v", file)
			}
		})
	}
}
