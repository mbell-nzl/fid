package main

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

type p24StopRow struct {
	code, condition, details, next string
	labels                         []string
}

var p24StopRows = []p24StopRow{
	{"executor-liveness-unknown", "executor-terminated", `{"run_id":"RUN-one","executor_id":"exec-one","observation_id":null}`, "fid run observe", []string{"Run: RUN-one", "Executor: exec-one", "Observation: unknown"}},
	{"recovery-state-conflict", "run-recovered", `{"run_id":"RUN-one","observation_digest":"sha256:abc","paths":["base.txt"],"reasons":["changed"]}`, "fid run recover", []string{"Run: RUN-one", "Observation digest: sha256:abc", "Paths: base.txt", "Reasons: changed"}},
	{"activation-incomplete", "activation-recovered", `{"activation_id":"ACT-one","point_of_no_return":true,"components":["kit"]}`, "fid activation recover", []string{"Activation: ACT-one", "Point of no return: true", "Components: kit"}},
	{"capability-unavailable", "capability-available", `{"capability_id":"gate.run","available":false,"owner_phase":26}`, "wait for the capability gate.run", []string{"Capability: gate.run", "Available: false", "Owner phase: 26"}},
	{"budget-exhausted", "budget-observed", `{"envelope_id":"ENV-one","budget_name":"tokens","unit":"tokens","maximum":10,"consumed":10,"remaining":0}`, "fid run observe-budget", []string{"Envelope: ENV-one", "Budget: tokens", "Unit: tokens", "Maximum: 10", "Consumed: 10", "Remaining: 0"}},
	{"operator-stop", "work-ended", `{"reason":"operator approval","reference":"DEC-one"}`, "fid run close or fid run abandon", []string{"Reason: operator approval", "Reference: DEC-one"}},
	{"review-rework-exhausted", "review-recorded", `{"envelope_id":"ENV-one","budget_name":"reviews","unit":"review-cycles","maximum":1,"consumed":1,"remaining":0,"last_review_id":"REV-one"}`, "fid review record", []string{"Envelope: ENV-one", "Budget: reviews", "Unit: review-cycles", "Maximum: 1", "Consumed: 1", "Remaining: 0", "Last Review: REV-one"}},
	{"budget-data-unavailable", "finding-resolved", `{"envelope_id":"ENV-one","budget_name":"tokens","unit":"tokens","maximum":10,"consumed":null,"remaining":null}`, "fid review record", []string{"Envelope: ENV-one", "Budget: tokens", "Unit: tokens", "Maximum: 10", "Consumed: unknown", "Remaining: unknown"}},
	{"human-outcome-required", "human-outcome-recorded", `{"envelope_id":"ENV-one","obligation_id":"OBL-one","outcome_class":"security"}`, "fid human-outcome record", []string{"Envelope: ENV-one", "Obligation: OBL-one", "Outcome class: security"}},
	{"protected-decision-required", "decision-satisfied", `{"envelope_id":"ENV-one","decision_class":"security","required_capabilities":["decision.record-reference"],"affected_actions":["run.start"]}`, "fid decision record-reference then fid human-outcome record", []string{"Envelope: ENV-one", "Decision class: security", "Required capabilities: decision.record-reference", "Affected actions: run.start"}},
}

func TestP24StopHumanCorpus(t *testing.T) {
	t.Parallel()
	for _, row := range p24StopRows {
		t.Run(row.code, func(t *testing.T) {
			t.Parallel()
			condition := workflow.Condition{Type: row.condition, Subject: workflow.Subject{Type: "run", ID: "RUN-one"}}
			switch row.condition {
			case "budget-observed":
				condition.BudgetName = stringPointer("tokens")
			case "capability-available":
				condition.CapabilityID = stringPointer("gate.run")
			case "finding-resolved":
				condition.FindingID = stringPointer("FND-one")
			case "human-outcome-recorded":
				condition.OutcomeClass = stringPointer("security")
			case "decision-satisfied":
				condition.DecisionClass = stringPointer("security")
			}
			for _, resolved := range []bool{false, true} {
				stop := workflow.Stop{ID: "STP-one", Subject: condition.Subject, Code: row.code, BlockedActions: []string{"run.start"}, Details: json.RawMessage(row.details), ResolutionCondition: condition, ResolverCapability: "stop.resolve"}
				attention := workflow.Attention{ID: "ATT-one", Status: "open", Code: row.code, StopIDs: []string{stop.ID}}
				if resolved {
					stop.Resolution = &workflow.StopResolution{FactEventID: "EVT-fact", EventID: "EVT-resolved"}
					attention.Status = "resolved"
					attention.ResolutionEventID = stringPointer("EVT-attention")
				}
				for _, view := range []string{"stop.list", "attention.list", "resume"} {
					var out bytes.Buffer
					cmd := &cobra.Command{}
					cmd.SetOut(&out)
					var err error
					switch view {
					case "stop.list":
						err = writeStopHuman(cmd, view, map[string]any{"items": []workflow.Stop{stop}})
					case "attention.list":
						text, e := attentionHumanDescription(cmd, attention, map[string]workflow.Stop{stop.ID: stop})
						err = e
						out.WriteString(text)
					case "resume":
						err = writeResumeHuman(cmd, workflow.Resume{Stops: []workflow.Stop{stop}, Attention: []workflow.Attention{attention}}, workflow.RecoveryReadContext{})
					}
					if err != nil {
						t.Fatal(err)
					}
					text := out.String()
					wants := append([]string{"Subject: run:RUN-one", "Code: " + row.code, "Blocked actions: run.start", "Resolution condition: " + row.condition + " for run:RUN-one", "Resolver capability: stop.resolve"}, row.labels...)
					if resolved {
						wants = append(wants, "Status: resolved", "Resolution Event: EVT-resolved", "Fact Event: EVT-fact")
					} else {
						wants = append(wants, "Status: open", "Next: "+row.next+"; after the fact: fid stop resolve, then fid attention resolve")
					}
					for _, want := range wants {
						if !strings.Contains(text, want) {
							t.Errorf("%s resolved=%t missing %q:\n%s", view, resolved, want, text)
						}
					}
					if strings.Contains(text, "{") || strings.Contains(text, "<nil>") || resolved && strings.Contains(text, "Next:") {
						t.Errorf("%s leaked JSON/nil or resolved guidance: %s", view, text)
					}
					if view != "stop.list" {
						for _, want := range []string{"Attention ATT-one", "Code: " + row.code, "Linked Stops: STP-one"} {
							if !strings.Contains(text, want) {
								t.Errorf("missing Attention %q", want)
							}
						}
						if resolved && !strings.Contains(text, "Resolution Event: EVT-attention") {
							t.Error("missing Attention resolution Event")
						}
					}
				}
			}
		})
	}
}

func TestP24StopHumanCorpusCompleteness(t *testing.T) {
	t.Parallel()
	file, err := parser.ParseFile(token.NewFileSet(), "../../internal/workflow/stop.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, registry := range []string{"stopCodes", "conditionTypes"} {
		seen := map[string]bool{}
		for _, row := range p24StopRows {
			key := row.code
			if registry == "conditionTypes" {
				key = row.condition
			}
			if seen[key] {
				t.Fatalf("duplicate %s row %s", registry, key)
			}
			seen[key] = true
		}
		ast.Inspect(file, func(n ast.Node) bool {
			v, ok := n.(*ast.ValueSpec)
			if !ok || len(v.Names) != 1 || v.Names[0].Name != registry {
				return true
			}
			literal := v.Values[0].(*ast.CompositeLit)
			for _, member := range literal.Elts {
				key, _ := strconv.Unquote(member.(*ast.KeyValueExpr).Key.(*ast.BasicLit).Value)
				if !seen[key] {
					t.Errorf("missing %s row %s", registry, key)
				}
				delete(seen, key)
			}
			return false
		})
		for key := range seen {
			t.Errorf("stale %s row %s", registry, key)
		}
	}
}
