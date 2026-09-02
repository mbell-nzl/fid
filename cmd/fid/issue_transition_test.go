package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestIssueTransitionJSONAllowsStandaloneActorWithoutGenericClaim(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	create := executeCommandForTest(workspace, "--json", "issue", "create", "--kind", "task", "--title", "Claim precondition", "--actor", "agent:developer")
	if create.err != nil {
		t.Fatalf("issue create error = %v, stderr = %q", create.err, create.stderr)
	}
	var created struct {
		Result struct {
			Issue struct {
				ID string `json:"id"`
			} `json:"issue"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(create.stdout), &created); err != nil {
		t.Fatalf("Unmarshal(create stdout) error = %v", err)
	}

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	transition := executeCommandForTest(workspace, "--json", "issue", "transition", created.Result.Issue.ID, "--to", "in_progress", "--actor", "agent:developer")
	if transition.err != nil || transition.stderr != "" {
		t.Fatalf("standalone transition = %#v, want stdout-only success", transition)
	}
	var failure struct {
		Operation string `json:"operation"`
		Result    struct {
			Issue struct {
				ID       string `json:"id"`
				Revision int    `json:"revision"`
			} `json:"issue"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(transition.stdout), &failure); err != nil {
		t.Fatalf("Unmarshal(transition stdout) error = %v; stdout = %q", err, transition.stdout)
	}
	if failure.Operation != "issue.transition" || failure.Result.Issue.ID != created.Result.Issue.ID || failure.Result.Issue.Revision != 2 {
		t.Fatalf("transition result = %+v, want issue.transition/%s/revision 2", failure, created.Result.Issue.ID)
	}
	if linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath))); linesAfter != linesBefore+1 {
		t.Fatalf("journal lines after standalone transition = %d, want %d", linesAfter, linesBefore+1)
	}
}

func TestIssueTransitionExpectedRevisionMatrix(t *testing.T) {
	tests := []struct {
		operation        string
		expectedRevision string
		wantSuccess      bool
	}{
		{operation: "update", expectedRevision: "", wantSuccess: true},
		{operation: "update", expectedRevision: "1", wantSuccess: true},
		{operation: "update", expectedRevision: "2", wantSuccess: false},
		{operation: "transition", expectedRevision: "", wantSuccess: true},
		{operation: "transition", expectedRevision: "1", wantSuccess: true},
		{operation: "transition", expectedRevision: "2", wantSuccess: false},
	}

	for _, tt := range tests {
		name := tt.operation + "/"
		if tt.expectedRevision == "" {
			name += "absent"
		} else if tt.wantSuccess {
			name += "matching"
		} else {
			name += "stale"
		}
		t.Run(name, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)

			created := executeCommandForTest(workspace, "--json", "issue", "create", "--kind", "task", "--title", "Expected revision "+name, "--actor", "agent:developer")
			if created.err != nil || created.stderr != "" {
				t.Fatalf("issue create = %#v, want stdout-only success", created)
			}
			var createEnvelope struct {
				Result struct {
					Issue struct {
						ID       string `json:"id"`
						Revision int    `json:"revision"`
					} `json:"issue"`
				} `json:"result"`
			}
			if err := json.Unmarshal([]byte(created.stdout), &createEnvelope); err != nil {
				t.Fatalf("Unmarshal(create stdout) error = %v", err)
			}
			issueID := createEnvelope.Result.Issue.ID
			before := canonicalIssueForTest(t, workspace, issueID)
			journalPath := mustJournalPath(t, workspace)
			linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

			args := []string{"--json", "issue", tt.operation, issueID}
			if tt.operation == "update" {
				args = append(args, "--title", "Updated "+name)
			} else {
				args = append(args, "--to", "in_progress")
			}
			if tt.expectedRevision != "" {
				args = append(args, "--expected-revision", tt.expectedRevision)
			}
			args = append(args, "--actor", "agent:developer")
			result := executeCommandForTest(workspace, args...)

			if tt.wantSuccess {
				if result.err != nil || result.stderr != "" {
					t.Fatalf("%s = %#v, want stdout-only success", tt.operation, result)
				}
				var success struct {
					Result struct {
						Issue struct {
							Revision int `json:"revision"`
						} `json:"issue"`
					} `json:"result"`
				}
				if err := json.Unmarshal([]byte(result.stdout), &success); err != nil {
					t.Fatalf("Unmarshal(%s stdout) error = %v", tt.operation, err)
				}
				if got, want := success.Result.Issue.Revision, before.Revision+1; got != want {
					t.Fatalf("%s result revision = %d, want %d", tt.operation, got, want)
				}
				if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore+1 {
					t.Fatalf("journal lines after %s = %d, want %d", tt.operation, got, linesBefore+1)
				}
				return
			}

			if result.err == nil {
				t.Fatalf("%s error = nil, want stale revision rejection", tt.operation)
			}
			if result.stdout != "" {
				t.Fatalf("stdout = %q, want empty", result.stdout)
			}
			var failure struct {
				Operation string `json:"operation"`
				Error     struct {
					Code    string `json:"code"`
					Details struct {
						Expected int `json:"expected_revision"`
						Actual   int `json:"actual_revision"`
					} `json:"details"`
				} `json:"error"`
			}
			if err := json.Unmarshal([]byte(result.stderr), &failure); err != nil {
				t.Fatalf("Unmarshal(%s stderr) error = %v; stderr = %q", tt.operation, err, result.stderr)
			}
			if failure.Operation != "issue."+tt.operation || failure.Error.Code != "state-conflict" {
				t.Fatalf("failure = %+v, want issue.%s state-conflict", failure, tt.operation)
			}
			if failure.Error.Details.Expected != 2 || failure.Error.Details.Actual != before.Revision {
				t.Fatalf("state-conflict details = %+v, want expected 2 actual %d", failure.Error.Details, before.Revision)
			}
			if after := canonicalIssueForTest(t, workspace, issueID); after != before {
				t.Fatalf("Issue after stale %s = %+v, want unchanged %+v", tt.operation, after, before)
			}
			if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
				t.Fatalf("journal lines after stale %s = %d, want %d", tt.operation, got, linesBefore)
			}
		})
	}
}

func TestIssueTransitionCanonicalOutcomeTable(t *testing.T) {
	tests := []struct {
		name            string
		toState         string
		actor           string
		setup           func(t *testing.T, workspace, issueID string)
		wantSuccess     bool
		wantFrom        string
		wantUnsatisfied []string
	}{
		{
			name:        "legal edge",
			toState:     "in_progress",
			actor:       "developer",
			wantSuccess: true,
			wantFrom:    "todo",
		},
		{
			name:            "unknown target state",
			toState:         "not_a_state",
			actor:           "developer",
			wantFrom:        "todo",
			wantUnsatisfied: []string{},
		},
		{
			name:            "absent edge between known states",
			toState:         "cancelled",
			actor:           "tl",
			wantFrom:        "todo",
			wantUnsatisfied: []string{},
			setup: func(t *testing.T, workspace, _ string) {
				appendCommandIssueEvent(t, workspace, journal.Event{
					EventSchema: journal.CurrentEventSchema,
					EventType:   "workflow.transition.removed",
					Actor:       journal.Actor{Role: "tl"},
					Command:     journal.Command{Argv: []string{"fid", "maintenance", "seed"}, CWD: ".", FIDVersion: "0.1.0"},
					Effects: []journal.Effect{{
						Op:    "delete",
						Table: "workflow_transitions",
						PK: map[string]any{
							"workflow_id": "task-default",
							"from_state":  "todo",
							"to_state":    "cancelled",
						},
					}},
				})
			},
		},
		{
			name:            "stored unsupported precondition",
			toState:         "in_progress",
			actor:           "developer",
			wantFrom:        "todo",
			wantUnsatisfied: []string{"unknown_rule"},
			setup: func(t *testing.T, workspace, _ string) {
				appendCommandIssueEvent(t, workspace, journal.Event{
					EventSchema: journal.CurrentEventSchema,
					EventType:   "workflow.precondition.seeded",
					Actor:       journal.Actor{Role: "tl"},
					Command:     journal.Command{Argv: []string{"fid", "maintenance", "seed"}, CWD: ".", FIDVersion: "0.1.0"},
					Effects: []journal.Effect{{
						Op:    "insert",
						Table: "workflow_transition_requires",
						PK: map[string]any{
							"workflow_id": "task-default",
							"from_state":  "todo",
							"to_state":    "in_progress",
							"rule":        "unknown_rule",
						},
						Values: map[string]any{"params_json": nil},
					}},
				})
			},
		},
		{
			name:        "retired completion precondition is absent",
			toState:     "dev_complete",
			actor:       "developer",
			wantSuccess: true,
			wantFrom:    "in_progress",
			setup: func(t *testing.T, workspace, issueID string) {
				started := executeCommandForTest(workspace, "--json", "issue", "transition", issueID, "--to", "in_progress", "--actor", "agent:developer")
				if started.err != nil || started.stderr != "" {
					t.Fatalf("setup transition = %#v, want stdout-only success", started)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)
			created := executeCommandForTest(workspace, "--json", "issue", "create", "--kind", "task", "--title", "Transition "+tt.name, "--actor", "agent:developer")
			if created.err != nil || created.stderr != "" {
				t.Fatalf("issue create = %#v, want stdout-only success", created)
			}
			var createEnvelope struct {
				Result struct {
					Issue struct {
						ID string `json:"id"`
					} `json:"issue"`
				} `json:"result"`
			}
			if err := json.Unmarshal([]byte(created.stdout), &createEnvelope); err != nil {
				t.Fatalf("Unmarshal(create stdout) error = %v", err)
			}
			issueID := createEnvelope.Result.Issue.ID
			if tt.setup != nil {
				tt.setup(t, workspace, issueID)
			}

			before := canonicalIssueForTest(t, workspace, issueID)
			journalPath := mustJournalPath(t, workspace)
			linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
			result := executeCommandForTest(workspace, "--json", "issue", "transition", issueID, "--to", tt.toState, "--actor", "agent:"+tt.actor)

			if tt.wantSuccess {
				if result.err != nil || result.stderr != "" {
					t.Fatalf("transition = %#v, want stdout-only success", result)
				}
				var success struct {
					Result struct {
						Issue struct {
							Revision int `json:"revision"`
						} `json:"issue"`
					} `json:"result"`
					Event struct {
						ID         string `json:"id"`
						OccurredAt string `json:"occurred_at"`
					} `json:"event"`
				}
				if err := json.Unmarshal([]byte(result.stdout), &success); err != nil {
					t.Fatalf("Unmarshal(transition stdout) error = %v", err)
				}
				if got, want := success.Result.Issue.Revision, before.Revision+1; got != want {
					t.Fatalf("transition revision = %d, want %d", got, want)
				}
				if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore+1 {
					t.Fatalf("journal lines after legal transition = %d, want %d", got, linesBefore+1)
				}
				journalTail := lastCommandJournalEvent(t, journalPath)
				if success.Event.ID != journalTail.EventID || success.Event.OccurredAt != journalTail.OccurredAt.UTC().Format(time.RFC3339Nano) {
					t.Fatalf("transition event = (%q, %s), want journal tail (%q, %s)", success.Event.ID, success.Event.OccurredAt, journalTail.EventID, journalTail.OccurredAt)
				}
				return
			}

			if result.err == nil {
				t.Fatal("transition error = nil, want rejection")
			}
			if result.stdout != "" {
				t.Fatalf("stdout = %q, want empty", result.stdout)
			}
			var failure struct {
				Operation string `json:"operation"`
				Error     struct {
					Code    string `json:"code"`
					Details struct {
						FromState   string   `json:"from_state"`
						ToState     string   `json:"to_state"`
						Unsatisfied []string `json:"unsatisfied"`
					} `json:"details"`
				} `json:"error"`
			}
			if err := json.Unmarshal([]byte(result.stderr), &failure); err != nil {
				t.Fatalf("Unmarshal(transition stderr) error = %v; stderr = %q", err, result.stderr)
			}
			if failure.Operation != "issue.transition" || failure.Error.Code != "transition-not-permitted" {
				t.Fatalf("failure = %+v, want issue.transition transition-not-permitted", failure)
			}
			if failure.Error.Details.FromState != tt.wantFrom || failure.Error.Details.ToState != tt.toState || strings.Join(failure.Error.Details.Unsatisfied, ",") != strings.Join(tt.wantUnsatisfied, ",") {
				t.Fatalf("transition details = %+v, want from=%q to=%q unsatisfied=%#v", failure.Error.Details, tt.wantFrom, tt.toState, tt.wantUnsatisfied)
			}
			if after := canonicalIssueForTest(t, workspace, issueID); after.Status != before.Status {
				t.Fatalf("Issue status after rejected transition = %q, want unchanged %q", after.Status, before.Status)
			}
			if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
				t.Fatalf("journal lines after rejected transition = %d, want %d", got, linesBefore)
			}
		})
	}
}

func TestIssueMutationCanonicalAttributionAvoidsProtectedEffects(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, workspace string) string
		args  func(issueID string) []string
	}{
		{
			name: "create",
			setup: func(_ *testing.T, _ string) string {
				return ""
			},
			args: func(_ string) []string {
				return []string{"--json", "issue", "create", "--kind", "task", "--title", "Attribution create", "--actor", "agent:developer"}
			},
		},
		{
			name: "update",
			setup: func(t *testing.T, workspace string) string {
				return createAttributionIssue(t, workspace, "Attribution update target")
			},
			args: func(issueID string) []string {
				return []string{"--json", "issue", "update", issueID, "--title", "Attribution updated", "--actor", "agent:developer"}
			},
		},
		{
			name: "transition",
			setup: func(t *testing.T, workspace string) string {
				issueID := createAttributionIssue(t, workspace, "Attribution transition target")
				return issueID
			},
			args: func(issueID string) []string {
				return []string{"--json", "issue", "transition", issueID, "--to", "in_progress", "--actor", "agent:developer"}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)
			issueID := tt.setup(t, workspace)
			refreshed := executeCommandForTest(workspace, "--json", "rebuild")
			if refreshed.err != nil || refreshed.stderr != "" {
				t.Fatalf("fixture rebuild = %#v, want stdout-only success", refreshed)
			}
			beforeCounts := commandProjectionTableCounts(t, workspace)
			result := executeCommandForTest(workspace, tt.args(issueID)...)
			if result.err != nil || result.stderr != "" {
				t.Fatalf("issue %s = %#v, want stdout-only success", tt.name, result)
			}
			var envelope struct {
				Result struct {
					Issue struct {
						ID string `json:"id"`
					} `json:"issue"`
				} `json:"result"`
				Event struct {
					ID string `json:"id"`
				} `json:"event"`
			}
			if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
				t.Fatalf("Unmarshal(%s stdout) error = %v", tt.name, err)
			}
			if tt.name == "create" {
				issueID = envelope.Result.Issue.ID
			}

			journalTail := lastCommandJournalEvent(t, mustJournalPath(t, workspace))
			if journalTail.EventID != envelope.Event.ID || journalTail.ActorV2.ActorID != "developer" || journalTail.ActorV2.Producer != "fid-standalone" || journalTail.ActorV2.ProvenanceRef != "project-local-attribution" {
				t.Fatalf("journal tail = (%q, actor %#v), want returned event %q with local developer attribution", journalTail.EventID, journalTail.ActorV2, envelope.Event.ID)
			}
			for _, effect := range journalTail.Effects {
				if effect.Table != "issues" && effect.Table != "activity_log" && effect.Table != "idempotency_results" && effect.Table != "issue_revisions" {
					t.Fatalf("%s Event effect table = %q, want retained Issue, activity, revision, or idempotency projection", tt.name, effect.Table)
				}
			}
			assertProjectedIssueActivityAttribution(t, workspace, journalTail.EventID, issueID, "developer")

			afterCounts := commandProjectionTableCounts(t, workspace)
			for table, before := range beforeCounts {
				if table == "issues" || table == "activity_log" || table == "idempotency_results" || table == "issue_revisions" || table == "journal_events" {
					continue
				}
				if after := afterCounts[table]; after != before {
					t.Fatalf("projection table %q count after %s = %d, want unchanged %d", table, tt.name, after, before)
				}
			}
		})
	}
}

func createAttributionIssue(t *testing.T, workspace, title string) string {
	t.Helper()
	result := executeCommandForTest(workspace, "--json", "issue", "create", "--kind", "task", "--title", title, "--actor", "agent:developer")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("issue create = %#v, want stdout-only success", result)
	}
	var envelope struct {
		Result struct {
			Issue struct {
				ID string `json:"id"`
			} `json:"issue"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("Unmarshal(create stdout) error = %v", err)
	}
	return envelope.Result.Issue.ID
}

func assertProjectedIssueActivityAttribution(t *testing.T, workspace, eventID, issueID, wantActor string) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projected database: %v", err)
	}
	defer closeTransitionCommandDB(t, conn)

	var actor, refID string
	if err := conn.QueryRow(`SELECT actor, ref_id FROM activity_log WHERE journal_event_id = ?`, eventID).Scan(&actor, &refID); err != nil {
		t.Fatalf("query projected activity: %v", err)
	}
	if actor != wantActor || refID != issueID {
		t.Fatalf("projected activity = (actor %q, ref_id %q), want (%q, %q)", actor, refID, wantActor, issueID)
	}
}

func commandProjectionTableCounts(t *testing.T, workspace string) map[string]int {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projected database: %v", err)
	}
	defer closeTransitionCommandDB(t, conn)

	rows, err := conn.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		t.Fatalf("list projection tables: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			t.Fatalf("close projection tables: %v", err)
		}
	}()

	counts := make(map[string]int)
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			t.Fatalf("scan projection table: %v", err)
		}
		var count int
		if err := conn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %q", table)).Scan(&count); err != nil {
			t.Fatalf("count projection table %q: %v", table, err)
		}
		counts[table] = count
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate projection tables: %v", err)
	}
	return counts
}

func lastCommandJournalEvent(t *testing.T, journalPath string) journal.Event {
	t.Helper()

	lines := splitJournalLines(mustReadCommandFile(t, journalPath))
	if len(lines) == 0 {
		t.Fatal("journal has no events")
	}
	var event journal.Event
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &event); err != nil {
		t.Fatalf("Unmarshal(journal tail) error = %v", err)
	}
	return event
}

type canonicalIssue struct {
	Title    string `json:"title"`
	Status   string `json:"status"`
	Revision int    `json:"revision"`
}

func canonicalIssueForTest(t *testing.T, workspace, issueID string) canonicalIssue {
	t.Helper()

	result := executeCommandForTest(workspace, "--json", "issue", "get", issueID)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("issue get = %#v, want stdout-only success", result)
	}
	var envelope struct {
		Result struct {
			Issue canonicalIssue `json:"issue"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("Unmarshal(issue get stdout) error = %v", err)
	}
	return envelope.Result.Issue
}

func TestIssueTransitionChangesStateAndPrintsConfirmation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI transition", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)
	_, transitionStderr := runFIDCommand(t, workspace, "issue", "transition", issueID, "--to", "in_progress", "--actor", "agent:developer")
	if !strings.Contains(transitionStderr, "transitioned "+issueID+": todo → in_progress") {
		t.Fatalf("stderr = %q, want transition confirmation", transitionStderr)
	}

	stdout := mustGetIssueJSON(t, workspace, issueID)

	var got struct {
		Status         string `json:"status"`
		StatusCategory string `json:"status_category"`
	}
	if err := json.Unmarshal([]byte(stdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.Status != "in_progress" {
		t.Fatalf("status = %q, want in_progress", got.Status)
	}
	if got.StatusCategory != "started" {
		t.Fatalf("status_category = %q, want started", got.StatusCategory)
	}
}

func TestIssueTransitionRejectsInvalidEdgeAndLeavesStateUnchanged(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI invalid transition", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	beforeStdout := mustGetIssueJSON(t, workspace, issueID)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "transition", issueID, "--to", "done", "--actor", "agent:tl")
	if err == nil {
		t.Fatal("Execute() error = nil, want invalid transition rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		"error: transition issue",
		`"todo"`,
		`"done"`,
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr = %q, want substring %q", stderr, want)
		}
	}

	afterStdout := mustGetIssueJSON(t, workspace, issueID)
	if afterStdout != beforeStdout {
		t.Fatalf("issue JSON after rejected transition = %q, want %q", afterStdout, beforeStdout)
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore {
		t.Fatalf("journal lines after rejected transition = %d, want %d", linesAfter, linesBefore)
	}
}

func TestIssueTransitionDoesNotEnforceRetiredV1Rules(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI rule failure", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	runFIDCommand(t, workspace, "issue", "transition", issueID, "--to", "in_progress", "--actor", "agent:developer")

	stdout, stderr := runFIDCommand(t, workspace, "issue", "transition", issueID, "--to", "dev_complete", "--actor", "agent:developer")
	if stdout != "" || !strings.Contains(stderr, "transitioned "+issueID+": in_progress → dev_complete") {
		t.Fatalf("retired-requirement transition stdout/stderr = %q/%q, want successful transition", stdout, stderr)
	}
}

func TestIssueTransitionRequiresTargetStateAndKnownByRole(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI role validation", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	t.Run("missing to", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "transition", issueID, "--actor", "agent:developer")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing --to rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, "error: transition issue: target state is required") {
			t.Fatalf("stderr = %q, want missing --to error", stderr)
		}
	})

	t.Run("missing actor", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "transition", issueID, "--to", "in_progress")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing --actor rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, "error: --actor is required") {
			t.Fatalf("stderr = %q, want missing --actor error", stderr)
		}
	})

	t.Run("unregistered actor is standalone attribution", func(t *testing.T) {
		_, stderr := runFIDCommand(t, workspace, "issue", "transition", issueID, "--to", "in_progress", "--actor", "agent:intern")
		if !strings.Contains(stderr, "transitioned "+issueID+": todo → in_progress") {
			t.Fatalf("stderr = %q, want standalone transition confirmation", stderr)
		}
	})
}

func TestIssueTransitionUsesStandaloneActorAttribution(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI state permissions", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)
	stdout, stderr := runFIDCommand(t, workspace, "issue", "transition", issueID, "--to", "in_progress", "--actor", "agent:qa_reviewer")
	if stdout != "" {
		t.Fatalf("qa_reviewer stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "transitioned "+issueID+": todo → in_progress") {
		t.Fatalf("qa_reviewer stderr = %q, want standalone attribution transition confirmation", stderr)
	}
	last := lastCommandJournalEvent(t, mustJournalPath(t, workspace))
	if last.ActorV2.ActorID != "qa_reviewer" || last.ActorV2.Producer != "fid-standalone" {
		t.Fatalf("transition Event actor = %#v, want qa_reviewer standalone attribution", last.ActorV2)
	}
}

func closeTransitionCommandDB(t *testing.T, conn *sql.DB) {
	t.Helper()
	if err := conn.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
}
