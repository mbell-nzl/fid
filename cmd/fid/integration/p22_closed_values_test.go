//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Independent wire oracle: neither the documentation parser, inventory nor
// product allowlists supply these field-specific accepted values at runtime.
// This checks enum grammar only; existing lifecycle/admission tests determine
// which otherwise well-formed values are executable in a particular state.
var p22ClosedFieldGrammar = map[string]map[string]string{
	"Executor":           {"source": "activated-harness"},
	"Claim (amend)":      {"effective_status": "active expired released superseded"},
	"Liveness":           {"state": "live terminated unknown"},
	"PathChange":         {"classification": "in-scope outside-scope protected"},
	"Reconciliation":     {"operation": "none merge rebase cherry-pick revert bisect", "ancestry": "unchanged descendant conflict unknown"},
	"RecoveryAuthority":  {"mode": "ordinary force"},
	"Condition":          {"type": "executor-terminated run-recovered activation-recovered capability-available budget-observed work-ended"},
	"Attention":          {"status": "open resolved"},
	"BudgetSample":       {"source": "fid activated-harness", "sample_kind": "delta baseline"},
	"Budget (amend)":     {"state": "unknown available warning exhausted"},
	"ActivationResource": {"component": "kit overlay provider project backend", "state": "pending old new error"},
	"ActivationRecovery": {"direction": "restore-prior complete-forward"},
}

func p22ValidateClosedFields(domain string, value map[string]any) error {
	fields, exists := p22ClosedFieldGrammar[domain]
	if !exists {
		return fmt.Errorf("unknown closed domain %s", domain)
	}
	for field, grammar := range fields {
		member, ok := value[field].(string)
		if !ok {
			return fmt.Errorf("%s.%s is missing or not a string", domain, field)
		}
		found := false
		for _, allowed := range strings.Fields(grammar) {
			if member == allowed {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%s.%s rejects %q", domain, field, member)
		}
	}
	return nil
}

func TestP22ClosedValueSchemaCases(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(p22CorpusDir, "wire-inventory.json"))
	if err != nil {
		t.Fatal(err)
	}
	var inventory p22WireInventory
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	p22SameMembers(t, "closed field domains", p22ClosedDomainKeys(inventory.ClosedFields), p22ClosedDomainKeys(p22ClosedFieldGrammar))
	accepted, rejected := 0, 0
	for domain, fields := range inventory.ClosedFields {
		t.Run(domain, func(t *testing.T) {
			base := map[string]any{}
			var union []string
			for field, values := range fields {
				if len(values) == 0 {
					t.Fatalf("%s.%s has no qualifying cases", domain, field)
				}
				base[field] = values[0]
				union = append(union, values...)
			}
			p22SameMembers(t, "field fixture membership", union, inventory.ClosedValues[domain])
			for field, values := range fields {
				for _, member := range values {
					t.Run(field+"/"+member, func(t *testing.T) {
						candidate := p22CloneMap(base)
						candidate[field] = member
						if err := p22ValidateClosedFields(domain, candidate); err != nil {
							t.Fatal(err)
						}
						accepted++
					})
				}
				for _, invalid := range []any{"not-an-accepted-member", "", nil, 1, true} {
					candidate := p22CloneMap(base)
					candidate[field] = invalid
					if err := p22ValidateClosedFields(domain, candidate); err == nil {
						t.Fatalf("%s.%s admitted invalid member %#v", domain, field, invalid)
					}
					rejected++
				}
				candidate := p22CloneMap(base)
				delete(candidate, field)
				if err := p22ValidateClosedFields(domain, candidate); err == nil {
					t.Fatalf("%s.%s admitted missing member", domain, field)
				}
				rejected++
			}
		})
	}
	if accepted != 50 || rejected != 90 {
		t.Fatalf("closed cases = %d accepted/%d rejected, want50/90", accepted, rejected)
	}
	t.Logf("Closed wire fields: %d accepted member cases, %d invalid/missing cases across12 domains", accepted, rejected)
}

func p22ClosedDomainKeys[T any](domains map[string]T) []string {
	keys := []string{}
	for key := range domains {
		keys = append(keys, key)
	}
	return keys
}

// Actual public responses bind the grammar oracle to the product. Mutating a
// captured field proves that an unaccepted output fails this same conformance
// assertion; it is not a claim that a caller can submit derived output fields.
func TestP22ClosedValuesObserveRealResponses(t *testing.T) {
	a := p22NewFIDRecoveryAdapterWithLease(t, "human", 3600)
	run := p21Run(t, a.binary, a.workspace, "--json", "run", "get", a.runID)
	// DEC-201 explicitly adds effective_status to the resume Claim projection;
	// other Claim payloads retain their existing fields.
	claim := p21Run(t, a.binary, a.workspace, "--json", "resume", a.issueID, "--repository", a.workspace)
	observed := p22AcceptedCall(t, a, "closed-observe", "run.observe", nil)
	stopped := p22AcceptedCall(t, a, "closed-stop", "stop.raise", nil)
	budget := p22AcceptedCall(t, a, "closed-budget", "run.observe-budget", map[string]any{"value": 0})
	p22Write(t, a.workspace, "base.txt", "closed-value dirty work\n")
	inspected := p22AcceptedCall(t, a, "closed-inspect", "run.inspect", nil)
	p22AcceptedCall(t, a, "closed-expire", "claim.expire", nil)
	recovered := p22AcceptedCall(t, a, "closed-recover", "run.recover", map[string]any{"adopted_paths": []string{"base.txt"}, "ttl_seconds": 3600})
	activation := p21Run(t, a.binary, a.workspace, "--json", "activation", "status", a.activation["id"].(string))
	b := p22NewFIDActivationAdapter(t, "prepared").(*p22FIDRecoveryAdapter)
	selected := p22AcceptedCall(t, b, "closed-activation-recover", "activation.recover", nil)
	samples := map[string]any{
		"Executor":           p22ResultPath(run, "result.run.executor"),
		"Claim (amend)":      p22ResultPath(claim, "result.resume.claim"),
		"Liveness":           p22ResultPath(observed, "result.observation.liveness"),
		"PathChange":         p22ResultPath(inspected, "result.reconciliation.paths.0"),
		"Reconciliation":     p22ResultPath(inspected, "result.reconciliation"),
		"RecoveryAuthority":  p22ResultPath(recovered, "result.recovery.authority"),
		"Condition":          p22ResultPath(stopped, "result.stop.resolution_condition"),
		"Attention":          p22ResultPath(stopped, "result.attention"),
		"BudgetSample":       p22ResultPath(budget, "result.sample"),
		"Budget (amend)":     p22ResultPath(run, "result.run.budgets.0"),
		"ActivationResource": p22ResultPath(activation, "result.activation.resources.0"),
		"ActivationRecovery": p22ResultPath(selected, "result.activation.recovery"),
	}
	p22SameMembers(t, "actual response domains", p22ClosedDomainKeys(samples), p22ClosedDomainKeys(p22ClosedFieldGrammar))
	for domain, sample := range samples {
		t.Run(domain, func(t *testing.T) {
			value, ok := sample.(map[string]any)
			if !ok {
				t.Fatalf("actual %s object missing: %#v", domain, sample)
			}
			if err := p22ValidateClosedFields(domain, value); err != nil {
				t.Fatal(err)
			}
			for field := range p22ClosedFieldGrammar[domain] {
				changed := p22CloneMap(value)
				changed[field] = "not-an-accepted-member"
				if err := p22ValidateClosedFields(domain, changed); err == nil {
					t.Fatalf("mutated actual %s.%s was accepted", domain, field)
				}
			}
		})
	}
}
