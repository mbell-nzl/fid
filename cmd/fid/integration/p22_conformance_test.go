//go:build integration

package integration

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// This list is independent of the document parser and the implementation.
var p22AcceptedOperations = strings.Fields(`claim.expire run.start run.observe run.inspect run.checkpoint run.recover stop.raise stop.resolve stop.list attention.list attention.resolve run.observe-budget run.close run.abandon activation.prepare activation.component activation.commit activation.abort activation.finalize activation.resource activation.recover activation.recover-complete export.import`)

type p22WireInventory struct {
	Operations     map[string][]string            `json:"operations"`
	Failures       map[string][]string            `json:"failures"`
	Frames         map[string][]string            `json:"frames"`
	ClosedValues   map[string][]string            `json:"closed_values"`
	ClosedEvidence map[string][]string            `json:"closed_evidence"`
	ClosedFields   map[string]map[string][]string `json:"closed_fields"`
}

func TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join(p22CorpusDir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Inventory string `json:"wire_inventory"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	var inventory p22WireInventory
	if manifest.Inventory != "" {
		data, err = os.ReadFile(filepath.Join(p22CorpusDir, manifest.Inventory))
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, &inventory); err != nil {
			t.Fatal(err)
		}
	}
	for _, operation := range p22AcceptedOperations {
		if len(inventory.Operations[operation]) == 0 {
			t.Fatalf("amended operation %s has no conformance fixture", operation)
		}
	}
	doc, err := os.ReadFile(filepath.Join("testdata", "p22_public_interface.md"))
	if err != nil {
		t.Fatal(err)
	}
	leaves := p21Section(t, string(doc), "## 3. New And Replaced Leaves", "## 4. Failure And Live Frames")
	var documented []string
	for _, line := range strings.Split(leaves, "\n") {
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		first := strings.Split(line, "|")[1]
		for _, match := range regexp.MustCompile("`([a-z]+\\.[a-z-]+)`").FindAllStringSubmatch(first, -1) {
			documented = append(documented, match[1])
		}
	}
	p22SameMembers(t, "documented amended leaves", documented, p22AcceptedOperations)
	allOperations := append(append([]string{}, p22AcceptedOperations...), strings.Fields("activation.status export.create export.validate projection.resume run.get run.list claim.list gate.run")...)
	p22SameMembers(t, "operation fixtures", p22InventoryKeys(inventory.Operations), allOperations)
	failures := p21Section(t, string(doc), "## 4. Failure And Live Frames", "## 5. Manifest V2 And Availability")
	var codes []string
	for _, line := range strings.Split(failures, "\n") {
		if !strings.HasPrefix(line, "| ") {
			continue
		}
		cells := strings.Split(line, "|")
		if len(cells) < 4 || strings.TrimSpace(cells[1]) == "Code" || !strings.Contains(cells[2], " / ") {
			continue
		}
		for _, code := range strings.Split(strings.TrimSpace(cells[1]), " / ") {
			codes = append(codes, code)
		}
	}
	if len(codes) != 18 {
		t.Fatalf("documented new error count = %d, want18", len(codes))
	}
	codes = append(codes, "duplicate-entity", "activation-conflict", "archive-invalid", "gate-unsuccessful")
	p22SameMembers(t, "failure fixtures", p22InventoryKeys(inventory.Failures), codes)
	p22SameMembers(t, "live frame fixtures", p22InventoryKeys(inventory.Frames), strings.Fields("terminal.succeeded terminal.failed terminal.signaled terminal.timed-out terminal.cancelled fid-error.committed-denial fid-error.uncommitted budget-stop output.stdout-stderr output.truncated-redacted"))
	closed := p21Section(t, string(doc), "## 2. Closed Values", "### 2.1 Activated Capabilities And Grant Actions")
	wantClosed := map[string][]string{}
	for _, row := range regexp.MustCompile(`(?m)^\| ([^|]+) \| (.*) \|$`).FindAllStringSubmatch(closed, -1) {
		members := map[string]bool{}
		for _, match := range regexp.MustCompile(`"([a-z][a-z0-9.-]*)"`).FindAllStringSubmatch(row[2], -1) {
			members[match[1]] = true
		}
		if len(members) == 0 {
			continue
		}
		values := []string{}
		for member := range members {
			values = append(values, member)
		}
		sort.Strings(values)
		wantClosed[row[1]] = values
	}
	// The Phase 22 document remains the historical base. Phase 23 adds these
	// four Condition types; the retained fixture now covers the final wire set.
	wantClosed["Condition"] = append(wantClosed["Condition"],
		"review-recorded", "finding-resolved", "decision-satisfied", "human-outcome-recorded")
	sort.Strings(wantClosed["Condition"])
	if len(wantClosed) != 12 || !reflect.DeepEqual(wantClosed, inventory.ClosedValues) {
		t.Fatalf("closed-value snapshot differs: doc=%#v fixture=%#v", wantClosed, inventory.ClosedValues)
	}
	for domain := range inventory.ClosedValues {
		if len(inventory.ClosedEvidence[domain]) == 0 {
			t.Fatalf("closed domain %s has no qualifying/rejecting behavioral evidence", domain)
		}
	}
	p22SameMembers(t, "closed evidence domains", p22InventoryKeys(inventory.ClosedEvidence), p22InventoryKeys(inventory.ClosedValues))
	// References bind the inventory to real Go test declarations. The required
	// ordinary/tagged gates execute those tests; this check is not itself a
	// substitute for their behavioral assertions or successful gate receipts.
	parsed := map[string]*ast.File{}
	for kind, entries := range map[string]map[string][]string{"operation": inventory.Operations, "failure": inventory.Failures, "frame": inventory.Frames, "closed domain": inventory.ClosedEvidence} {
		for name, refs := range entries {
			if len(refs) == 0 {
				t.Fatalf("%s %s has no behavioral reference", kind, name)
			}
			for _, ref := range refs {
				p22CheckBehavioralReference(t, ref, parsed)
			}
		}
	}
}

func p22InventoryKeys(items map[string][]string) []string {
	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	return keys
}

func p22SameMembers(t *testing.T, name string, got, want []string) {
	t.Helper()
	got, want = append([]string{}, got...), append([]string{}, want...)
	sort.Strings(got)
	sort.Strings(want)
	if len(got) == 0 || !reflect.DeepEqual(got, want) {
		t.Fatalf("%s = %v, want exact %v", name, got, want)
	}
}

func p22CheckBehavioralReference(t *testing.T, ref string, parsed map[string]*ast.File) {
	t.Helper()
	parts := strings.Split(ref, "#")
	if len(parts) != 2 || !strings.HasSuffix(parts[0], "_test.go") || !strings.HasPrefix(parts[1], "Test") || strings.Contains(parts[0], "..") || (!strings.HasPrefix(parts[0], "cmd/") && !strings.HasPrefix(parts[0], "internal/")) {
		t.Fatalf("invalid behavioral reference %q", ref)
	}
	file := parsed[parts[0]]
	if file == nil {
		var err error
		file, err = parser.ParseFile(token.NewFileSet(), filepath.Join(integrationRepositoryRoot(), parts[0]), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		parsed[parts[0]] = file
	}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && function.Name.Name == parts[1] && function.Body != nil && len(function.Body.List) > 0 {
			return
		}
	}
	t.Fatalf("behavioral test is missing or empty: %s", ref)
}
