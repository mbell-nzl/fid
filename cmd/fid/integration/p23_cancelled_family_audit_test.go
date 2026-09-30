//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

type p23ProductSurface struct {
	Tables   []string `json:"tables"`
	Commands []string `json:"commands"`
}

// The removed delivery-workflow concepts have no product namespace. Inspect
// the built command tree and migrated SQLite catalog, including command aliases.
func TestP23CancelledFamiliesAbsentFromProduct(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapter(t)
	actual := p23ProductSurface{}
	for name := range p22FaultProjection(t, a.workspace) {
		actual.Tables = append(actual.Tables, name)
		if p23CancelledFamilyName(name) {
			t.Fatalf("cancelled product table %q", name)
		}
	}
	sort.Strings(actual.Tables)
	queue := [][]string{{}}
	seen := map[string]bool{}
	commands := 0
	for len(queue) != 0 {
		path := queue[0]
		queue = queue[1:]
		key := strings.Join(path, " ")
		if seen[key] {
			continue
		}
		seen[key] = true
		args := append(append([]string{}, path...), "--help")
		stdout, stderr, err := closureRunCommand(a.workspace, a.binary, args...)
		if err != nil || stderr != "" {
			t.Fatalf("fid %v --help: %v %q", path, err, stderr)
		}
		for _, name := range p23HelpCommandNames(stdout) {
			if p23CancelledFamilyName(name) {
				t.Fatalf("cancelled product command %q beneath %q", name, key)
			}
			queue = append(queue, append(append([]string{}, path...), name))
			actual.Commands = append(actual.Commands, strings.TrimSpace(key+" "+name))
			commands++
		}
		for _, line := range strings.Split(stdout, "\n") {
			if !strings.HasPrefix(line, "Aliases:") {
				continue
			}
			for _, alias := range strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "Aliases:")), ",") {
				if p23CancelledFamilyName(strings.TrimSpace(alias)) {
					t.Fatalf("cancelled product alias %q beneath %q", alias, key)
				}
			}
		}
	}
	if commands < 18 {
		t.Fatalf("real command tree traversal saw only %d commands", commands)
	}
	sort.Strings(actual.Commands)
	data, err := os.ReadFile(filepath.Join("testdata", "p23_conformance", "product-surface.json"))
	if err != nil {
		t.Logf("observed product surface: %+v", actual)
		t.Fatal(err)
	}
	var pinned p23ProductSurface
	if err := json.Unmarshal(data, &pinned); err != nil {
		t.Fatal(err)
	}
	p22SameMembers(t, "latest product table names", actual.Tables, pinned.Tables)
	p22SameMembers(t, "real command-tree paths", actual.Commands, pinned.Commands)
	t.Logf("inspected %d real commands and %d product tables", commands, len(p22FaultProjection(t, a.workspace)))
}

func p23HelpCommandNames(help string) []string {
	section := strings.SplitN(help, "Available Commands:\n", 2)
	if len(section) != 2 {
		return nil
	}
	section[1] = strings.SplitN(section[1], "\n\n", 2)[0]
	var names []string
	for _, match := range regexp.MustCompile(`(?m)^  ([a-z][a-z-]+)(?:\s|$)`).FindAllStringSubmatch(section[1], -1) {
		if match[1] != "help" {
			names = append(names, match[1])
		}
	}
	return names
}

func p23CancelledFamilyName(name string) bool {
	name = strings.ToLower(strings.ReplaceAll(name, "_", "-"))
	for _, family := range []string{"escalation", "escalations", "qa-plan", "qa-plans", "validation-contract", "validation-contracts", "template", "templates", "retrospective", "retrospectives", "carry-forward", "carry-forwards", "carryforward"} {
		if name == family || strings.HasPrefix(name, family+"-") {
			return true
		}
	}
	return false
}
