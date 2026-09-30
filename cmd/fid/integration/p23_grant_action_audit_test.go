//go:build integration

package integration

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type p23GrantAction struct {
	Name  string `json:"name"`
	Phase int    `json:"phase"`
	Task  int    `json:"task"`
}

func TestGrantActionSetIsExactlyThirtyNames(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "p23_conformance", "grant-actions.json"))
	if err != nil {
		t.Fatal(err)
	}
	var pinned []p23GrantAction
	if err := json.Unmarshal(data, &pinned); err != nil {
		t.Fatal(err)
	}
	if len(pinned) != 30 {
		t.Fatalf("pinned grant actions = %d, want 30", len(pinned))
	}
	want := make([]string, 0, len(pinned))
	for _, action := range pinned {
		if action.Phase < 21 || action.Phase > 23 || action.Phase == 23 && action.Task == 0 {
			t.Fatalf("missing owner for %s", action.Name)
		}
		want = append(want, action.Name)
	}
	profile := p23ProfileActionMap(t)
	journal := p23JournalActionConst(t)
	p22SameMembers(t, "profile actions versus independently pinned corpus", profile, want)
	p22SameMembers(t, "journal actions versus independently pinned corpus", journal, want)
	p22SameMembers(t, "profile actions versus journal actions", profile, journal)
	for _, unknown := range []string{"unknown.action", "requirement.record-reference", "requirement.link", "review.get", "human-outcome.list", "concern.raise", "issue.ac.verify"} {
		for _, actual := range profile {
			if actual == unknown {
				t.Fatalf("profile admitted %s", unknown)
			}
		}
		for _, actual := range journal {
			if actual == unknown {
				t.Fatalf("journal admitted %s", unknown)
			}
		}
	}
	// The real fixture puts the whole independent corpus into both an activated
	// profile and an Envelope Grant. Setup fails if either wire path rejects it.
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, want)
	for _, unknown := range []string{"unknown.action", "requirement.record-reference", "requirement.link", "review.get", "concern.raise"} {
		before := p21JournalEventCount(t, a.workspace)
		grant := map[string]any{"actions": []string{unknown}, "actor_selector": map[string]any{
			"actor_id": a.actor, "required_capabilities": []string{"envelope.create"},
		}, "scope": map[string]any{"issue_id": a.issueID}}
		obligation := map[string]any{"id": "OBL-grant-audit", "type": "gate-evidence", "applicability": "required",
			"guarded_actions": []string{"envelope.complete"}, "subject": map[string]any{"type": "issue", "id": a.issueID},
			"registration":            map[string]any{"id": a.registrationID, "version": a.registrationVersion},
			"accepted_evidence_types": []string{"receipt"}, "accepted_trust_sources": []string{"fid-executed"},
			"freshness_target": "repository-tree", "artifact_policy": "required-available", "independence": "fid-executed", "required_lane": nil, "source": a.measurement}
		request := map[string]any{"outcome": "Audit unknown grant", "non_goals": []string{},
			"risk_result":      map[string]any{"source": a.measurement, "result_digest": "sha256:" + p22Ones(64)},
			"repository_scope": a.repository, "sources": []any{a.measurement}, "primary_issue_id": a.issueID,
			"related_issue_ids": []string{}, "authority_grants": []any{grant}, "obligations": []any{obligation},
			"budgets": []any{map[string]any{"name": "attempts", "scope": "envelope", "unit": "attempts", "maximum": 3,
				"warning_threshold": nil, "exhaustion_action": "budget-exhausted", "consumed": nil, "remaining": nil}}}
		response := p23RunInput(t, a, "envelope", "create", a.actor, "grant-audit-"+strings.ReplaceAll(unknown, ".", "-"), request)
		assertP22ProcessFailure(t, "envelope.create", "invalid-argument", 2, response)
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("unknown grant %s appended Event %d -> %d", unknown, before, after)
		}
	}
}

func p23ProfileActionMap(t *testing.T) []string {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(integrationRepositoryRoot(), "internal/workflow/activation.go"), nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	var actions []string
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Name.Name != "sortedUniqueCapabilities" {
			continue
		}
		ast.Inspect(function.Body, func(node ast.Node) bool {
			assignment, ok := node.(*ast.AssignStmt)
			if !ok || len(assignment.Lhs) != 1 || len(assignment.Rhs) != 1 {
				return true
			}
			name, ok := assignment.Lhs[0].(*ast.Ident)
			if !ok || name.Name != "allowed" {
				return true
			}
			literal, ok := assignment.Rhs[0].(*ast.CompositeLit)
			if !ok {
				t.Fatal("profile action map missing")
			}
			for _, item := range literal.Elts {
				pair, ok := item.(*ast.KeyValueExpr)
				if !ok {
					t.Fatal("profile action map entry missing")
				}
				key, ok := pair.Key.(*ast.BasicLit)
				if !ok {
					t.Fatal("profile action key missing")
				}
				value, err := strconv.Unquote(key.Value)
				if err != nil {
					t.Fatal(err)
				}
				actions = append(actions, value)
			}
			return false
		})
	}
	if len(actions) == 0 {
		t.Fatal("sortedUniqueCapabilities action map missing")
	}
	return actions
}

func p23JournalActionConst(t *testing.T) []string {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(integrationRepositoryRoot(), "internal/journal/current_event_value_relations.go"), nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	var actions []string
	ast.Inspect(file, func(node ast.Node) bool {
		declaration, ok := node.(*ast.ValueSpec)
		if !ok || len(declaration.Names) != 1 || declaration.Names[0].Name != "actions" || len(declaration.Values) != 1 {
			return true
		}
		literal, ok := declaration.Values[0].(*ast.BasicLit)
		if !ok {
			t.Fatal("journal action const missing")
		}
		value, err := strconv.Unquote(literal.Value)
		if err != nil {
			t.Fatal(err)
		}
		actions = strings.Split(value, ",")
		return false
	})
	if len(actions) == 0 {
		t.Fatal("journal action const missing")
	}
	return actions
}
