package policy

import (
	"bytes"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestTrustByActorType(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		actor          string
		work, decision bool
	}{
		{"agent:builder", true, false}, {"human:mark", true, true}, {"service:ci", false, false},
		{"agent:", false, false}, {"human:bad id", false, false}, {"robot:one", false, false}, {"agent:a:b", false, false},
	} {
		t.Run(test.actor, func(t *testing.T) {
			if Allows(test.actor, WorkSet) != test.work || Allows(test.actor, PersonDecisions) != test.decision || Allows(test.actor, ActionSet("unknown")) {
				t.Fatalf("unexpected action sets for %q", test.actor)
			}
		})
	}
}

func TestEffectiveDigest(t *testing.T) {
	t.Parallel()
	decode := func(text string) Effective {
		t.Helper()
		file, err := config.DecodeProject([]byte(text))
		if err != nil {
			t.Fatal(err)
		}
		return Resolve(file)
	}
	first := decode("[checks.test]\nrun = ['go', 'test']\ntimeout = '1m'\ndefined_by = ['b', 'a']\n[lanes.standard]\nchecks = ['test']")
	second := decode("# reformatted and reordered\n[lanes.standard]\nchecks = [\"test\"]\n[checks.test]\ndefined_by = ['a', 'b']\ntimeout = '60s'\nrun = [\"go\", \"test\"]")
	a, err := first.Canonical()
	if err != nil {
		t.Fatal(err)
	}
	b, err := second.Canonical()
	if err != nil {
		t.Fatal(err)
	}
	if first.Digest() != second.Digest() || !bytes.Equal(a, b) {
		t.Fatalf("equivalent policies differ: %s vs %s", a, b)
	}
	for name, changed := range map[string]Effective{
		"command":               decode("[checks.test]\nrun = ['go', 'test', './...']\ntimeout = '1m'"),
		"unset vs empty checks": decode("[lanes.standard]\nchecks = []"),
		"default":               decode(""),
	} {
		if changed.Digest() == first.Digest() {
			t.Errorf("%s did not change digest", name)
		}
	}
	unset, empty := decode(""), decode("[lanes.standard]\nchecks = []")
	if unset.Digest() == empty.Digest() {
		t.Fatal("missing and explicit empty checks share a digest")
	}
	upgrade := unset
	upgrade.RulesVersion++
	if upgrade.Digest() == unset.Digest() {
		t.Fatal("built-in rule upgrade did not change digest")
	}
	if !bytes.Contains(a, []byte(".github/workflows/**")) || !bytes.Contains(a, []byte(".fid.toml")) || len(first.ActorRules["human"]) != 2 || len(first.ActorRules["service"]) != 0 {
		t.Fatalf("effective policy misses built-in rules: %s", a)
	}
	withIdentity := decode("[project]\nid = 'project-0123456789abcdef0123456789abcdef'")
	if withIdentity.Digest() != unset.Digest() {
		t.Fatal("project identity changed policy digest")
	}
}
