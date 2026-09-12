package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func TestHumanRecoveryGuidance(t *testing.T) {
	resume := workflow.Resume{
		Issue:       workflow.ResumeIssue{ID: "ISS-example", Status: "open"},
		Completion:  workflow.ResumeCompletion{Status: "incomplete"},
		Recoveries:  []workflow.Recovery{{PriorRunID: "RUN-old", SuccessorRunID: "RUN-new", AdoptedPaths: []string{"base.txt"}}},
		NextActions: []workflow.NextAction{{Operation: "gate.run", Subject: workflow.Subject{Type: "run", ID: "RUN-new"}}},
	}
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)
	if err := writeResumeHuman(cmd, resume, workflow.RecoveryReadContext{}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Next steps", "Completion requirements", "existing uncommitted changes", "tested", "reviewed", "base.txt", "gate.run"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("Human recovery missing %q: %s", want, out.String())
		}
	}
	if strings.Index(out.String(), "Next steps") > strings.Index(out.String(), "Completion requirements") {
		t.Error("completion requirements obscure immediate next steps")
	}
}

func TestHumanGuidanceDoesNotRecommendBlockedGate(t *testing.T) {
	resume := workflow.Resume{NextActions: []workflow.NextAction{{Operation: "gate.run", Reasons: []workflow.EvaluationReason{{Code: "budget-exhausted"}}}}}
	human := strings.Join(resumeNextStepLines(resume), "\n")
	if strings.Contains(human, "To collect proof") || !strings.Contains(human, "blocked") || !strings.Contains(human, "budget-exhausted") {
		t.Fatalf("blocked Gate was recommended or lost its reason: %s", human)
	}
}

func TestHumanActivationTerminalStateDoesNotOfferRecovery(t *testing.T) {
	for _, status := range []string{"finalized", "aborted"} {
		var out bytes.Buffer
		cmd := &cobra.Command{}
		cmd.SetOut(&out)
		if err := writeActivationStatusHuman(cmd, workflow.Activation{Status: status}); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(out.String(), "Recovery direction:") || !strings.Contains(out.String(), "terminal") {
			t.Fatalf("terminal activation offers recovery: %s", out.String())
		}
	}
}

func TestHumanImportedFreshObservationDoesNotRequestItAgain(t *testing.T) {
	context := workflow.RecoveryReadContext{Import: &workflow.ImportedHistoryContext{}, Liveness: &workflow.Liveness{State: "terminated"}}
	human := strings.Join(recoveryLocalContextLines(context, true), "\n")
	if strings.Contains(human, "establish fresh destination liveness") {
		t.Fatalf("existing destination observation described as missing: %s", human)
	}
}

func TestHumanImportedInactiveRunDoesNotRequestObservation(t *testing.T) {
	context := workflow.RecoveryReadContext{Import: &workflow.ImportedHistoryContext{}}
	for _, run := range []*workflow.Run{nil, {Status: "succeeded"}, {Status: "failed"}, {Status: "abandoned"}} {
		var out bytes.Buffer
		cmd := &cobra.Command{}
		cmd.SetOut(&out)
		if err := writeResumeHuman(cmd, workflow.Resume{Run: run}, context); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(out.String(), "Before continuing work here") || strings.Contains(out.String(), "establish fresh destination") {
			t.Fatalf("inactive or absent Run directs recovery: %s", out.String())
		}
	}
}
