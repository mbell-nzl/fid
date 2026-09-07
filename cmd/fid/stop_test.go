package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func TestStopSuccessHumanNamesBlockingReasonAndQualifyingResolution(t *testing.T) {
	condition := workflow.Condition{Type: "activation-recovered", Subject: workflow.Subject{Type: "activation", ID: "ACT-22"}}
	stop := workflow.Stop{
		ID: "STOP-22", Code: "operator-stop", BlockedActions: []string{"run.start"},
		Details: json.RawMessage(`{"reason":"operator approval","reference":"DEC-188"}`), ResolutionCondition: condition,
		Resolution: &workflow.StopResolution{FactEventID: "EVT-fact", EventID: "EVT-stop-resolved"},
	}
	attention := workflow.Attention{ID: "ATT-22", Code: "operator-stop", StopIDs: []string{"STOP-22"}, ResolutionEventID: stringPointer("EVT-attention-resolved")}

	var output bytes.Buffer
	cmd := &cobra.Command{}
	root := &cobra.Command{}
	root.PersistentFlags().Bool("json", false, "")
	root.AddCommand(cmd)
	cmd.SetOut(&output)
	if err := writeStopSuccess(cmd, "stop.raise", map[string]any{"stop": stop, "attention": attention}, "EVT-stop-raised", "2026-09-07T00:00:00Z"); err != nil {
		t.Fatalf("write Stop Human result: %v", err)
	}
	for _, want := range []string{"operator approval", "activation-recovered", "activation:ACT-22", "Attention ATT-22", "STOP-22"} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("Stop Human result = %q, want %q", output.String(), want)
		}
	}

	output.Reset()
	if err := writeStopSuccess(cmd, "attention.resolve", map[string]any{"attention": attention}, "EVT-attention-resolved", "2026-09-07T00:00:01Z"); err != nil {
		t.Fatalf("write Attention Human result: %v", err)
	}
	for _, want := range []string{"Attention ATT-22", "operator-stop", "EVT-attention-resolved"} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("Attention Human result = %q, want %q", output.String(), want)
		}
	}
}

func TestStopHumanListsShowOpenResolvedAndEmptyState(t *testing.T) {
	open := workflow.Stop{ID: "STOP-open", Code: "operator-stop", BlockedActions: []string{"run.start"}, Details: json.RawMessage(`{"reason":"operator approval","reference":"DEC-188"}`), ResolutionCondition: workflow.Condition{Type: "work-ended", Subject: workflow.Subject{Type: "run", ID: "RUN-22"}}}
	resolved := open
	resolved.ID = "STOP-resolved"
	resolved.Resolution = &workflow.StopResolution{FactEventID: "EVT-fact", EventID: "EVT-disposition"}
	for _, test := range []struct {
		name      string
		operation string
		result    map[string]any
		want      []string
		reject    string
	}{
		{"open", "stop.list", map[string]any{"items": []workflow.Stop{open}}, []string{"STOP-open", "operator approval", "run.start", "work-ended", "run:RUN-22"}, ""},
		{"resolved", "stop.list", map[string]any{"items": []workflow.Stop{resolved}}, []string{"STOP-resolved", "resolved", "EVT-fact", "EVT-disposition"}, "STOP-resolved blocks"},
		{"empty-stops", "stop.list", map[string]any{"items": []workflow.Stop{}}, []string{"No Stops"}, ""},
		{"empty-attention", "attention.list", map[string]any{"items": []workflow.Attention{}}, []string{"No Attention"}, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			cmd := &cobra.Command{}
			root := &cobra.Command{}
			root.PersistentFlags().Bool("json", false, "")
			root.AddCommand(cmd)
			cmd.SetOut(&output)
			if err := writeStopSuccess(cmd, test.operation, test.result, "", ""); err != nil {
				t.Fatal(err)
			}
			for _, want := range test.want {
				if !strings.Contains(output.String(), want) {
					t.Fatalf("Human output %q missing %q", output.String(), want)
				}
			}
			if test.reject != "" && strings.Contains(output.String(), test.reject) {
				t.Fatalf("Human output incorrectly says %q", output.String())
			}
		})
	}
}
