package main

// These assertions cover retiring kernel/historical behavior; retain until its owning cut.
import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestKernelIssueConcernJSONAndRebuildParity(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Surface enriched issue state", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, stderr)

	appendCommandIssueEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "issue.enriched_state_seeded",
		Actor: journal.Actor{
			Role: "qa_lead",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []journal.Effect{
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-epiccmd",
				},
				Values: map[string]any{
					"kind":        "epic",
					"title":       "CLI enriched epic",
					"workflow_id": "epic-default",
					"status":      "todo",
					"created_at":  "2026-06-07T14:00:00Z",
					"updated_at":  "2026-06-07T14:00:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "concerns",
				PK: map[string]any{
					"id": "CON-CMD-001",
				},
				Values: map[string]any{
					"subject_type": "issue", "subject_id": issueID,
					"actor_id": "developer", "scope_json": `{"paths":[]}`,
					"severity": "major", "summary": "CLI concern", "reasoning": "Needs review",
					"references_json": "[]", "status": "open",
					"created_at": "2026-06-07T14:01:30Z", "created_event_id": "EVT-CON-CMD-001",
				},
			},
			{
				Op:    "insert",
				Table: "triage_events",
				PK: map[string]any{
					"id": "TRG-CMD-001",
				},
				Values: map[string]any{
					"epic_id":       "FID-epiccmd",
					"issue_id":      issueID,
					"triager_role":  "qa_lead",
					"source":        "manual",
					"source_ref_id": "REV-CMD-001",
					"verdict":       "RESUME",
					"reasoning_md":  "project-docs/triage/TRG-CMD-001.md",
					"created_at":    "2026-06-07T14:02:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "triage_events",
				PK: map[string]any{
					"id": "TRG-CMD-002",
				},
				Values: map[string]any{
					"epic_id":       "FID-epiccmd",
					"issue_id":      issueID,
					"triager_role":  "qa_lead",
					"source":        "manual",
					"source_ref_id": "REV-CMD-001",
					"verdict":       "BLOCK",
					"reasoning_md":  "project-docs/triage/TRG-CMD-002.md",
					"created_at":    "2026-06-07T14:03:00Z",
				},
			},
		},
	})

	beforeStdout := mustGetIssueJSON(t, workspace, issueID)
	beforeConcerns := mustGetConcernListJSONForIssue(t, workspace, issueID)
	if len(beforeConcerns) != 1 || beforeConcerns[0].Status != "open" {
		t.Fatalf("before concerns = %#v, want one open concern", beforeConcerns)
	}
	if got := mustIssueTriageStatus(t, workspace, issueID); got == nil || *got != "blocked" {
		t.Fatalf("before triage_status = %#v, want blocked", got)
	}

	dbPath := filepath.Join(config.LedgerDir(workspace), "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	runFIDCommand(t, workspace, "rebuild")

	afterStdout := mustGetIssueJSON(t, workspace, issueID)
	if beforeStdout != afterStdout {
		t.Fatalf("issue JSON after rebuild = %q, want %q", afterStdout, beforeStdout)
	}
	if afterConcerns := mustGetConcernListJSONForIssue(t, workspace, issueID); !reflect.DeepEqual(afterConcerns, beforeConcerns) {
		t.Fatalf("concerns after rebuild = %#v, want %#v", afterConcerns, beforeConcerns)
	}
	if got := mustIssueTriageStatus(t, workspace, issueID); got == nil || *got != "blocked" {
		t.Fatalf("after triage_status = %#v, want blocked", got)
	}
}
