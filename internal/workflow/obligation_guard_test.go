package workflow

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"
)

const (
	guardActionComplete = "envelope.complete"
	guardActionTerminal = "issue.transition-terminal-success"
)

// obligationGuardCase describes one Envelope for the DEC-249 guarded-action
// probes: which actions each Obligation guards and whether the gate-evidence
// Obligation has qualifying Evidence.
type obligationGuardCase struct {
	name                string
	gateGuards          []string
	reviewGuards        []string
	gateSatisfied       bool
	wantTerminalBlocked []string
	wantCompleteBlocked []string
}

func obligationGuardCases() []obligationGuardCase {
	both := []string{guardActionComplete, guardActionTerminal}
	return []obligationGuardCase{
		{
			// VAL-009 blocking half: a satisfied gate-evidence Obligation and an
			// unsatisfied review Obligation guarding BOTH actions block both.
			name:       "unsatisfied review guarding both actions blocks both",
			gateGuards: both, reviewGuards: both, gateSatisfied: true,
			wantTerminalBlocked: []string{"OBL-review"}, wantCompleteBlocked: []string{"OBL-review"},
		},
		{
			// Contrast: the same unsatisfied review that guards only completion
			// leaves terminal success unblocked and still blocks completion.
			name:       "unsatisfied review guarding only completion leaves terminal success unblocked",
			gateGuards: both, reviewGuards: []string{guardActionComplete}, gateSatisfied: true,
			wantTerminalBlocked: nil, wantCompleteBlocked: []string{"OBL-review"},
		},
		{
			// No Obligation guards terminal success, so it carries no Obligation
			// requirement even though two Obligations are unsatisfied.
			name:       "no Obligation guards terminal success",
			gateGuards: []string{guardActionComplete}, reviewGuards: []string{guardActionComplete}, gateSatisfied: false,
			wantTerminalBlocked: nil, wantCompleteBlocked: []string{"OBL-gate", "OBL-review"},
		},
		{
			// The filter is symmetric: an unsatisfied Obligation guarding only
			// terminal success does not block completion.
			name:       "unsatisfied gate guarding only terminal success does not block completion",
			gateGuards: []string{guardActionTerminal}, reviewGuards: []string{guardActionComplete}, gateSatisfied: false,
			wantTerminalBlocked: []string{"OBL-gate"}, wantCompleteBlocked: []string{"OBL-review"},
		},
		{
			// Every reason from every applicable Obligation is returned.
			name:       "every unsatisfied applicable Obligation is named",
			gateGuards: both, reviewGuards: both, gateSatisfied: false,
			wantTerminalBlocked: []string{"OBL-gate", "OBL-review"}, wantCompleteBlocked: []string{"OBL-gate", "OBL-review"},
		},
	}
}

func TestGuardedActionsEvaluateOnlyObligationsThatGuardThem(t *testing.T) {
	for _, test := range obligationGuardCases() {
		t.Run(test.name, func(t *testing.T) {
			fixture := newObligationGuardFixture(t, test)
			terminal := fixture.terminalReasons(t)
			if got := obligationReasonSubjects(terminal); !reflect.DeepEqual(got, test.wantTerminalBlocked) {
				t.Fatalf("terminal success blocked by Obligations %v, want %v (reasons %#v)", got, test.wantTerminalBlocked, terminal)
			}
			if test.wantTerminalBlocked == nil && len(terminal) != 0 {
				t.Fatalf("terminal success reasons = %#v, want a fully permitted transition", terminal)
			}
			complete := fixture.completeReasons(t)
			if got := obligationReasonSubjects(complete); !reflect.DeepEqual(got, test.wantCompleteBlocked) {
				t.Fatalf("envelope.complete blocked by Obligations %v, want %v (reasons %#v)", got, test.wantCompleteBlocked, complete)
			}

			// Reads keep reporting every Obligation, whatever it guards.
			evaluations, err := currentEnvelopeEvaluations(fixture.conn, fixture.envelope, fixture.repository)
			if err != nil || len(evaluations) != 2 || evaluations[0].ObligationID != "OBL-gate" || evaluations[1].ObligationID != "OBL-review" {
				t.Fatalf("read evaluations = %#v, %v; want both Obligations in ID order", evaluations, err)
			}
			wantGate := "unsatisfied"
			if test.gateSatisfied {
				wantGate = "satisfied"
			}
			if evaluations[0].Status != wantGate || evaluations[1].Status != "unsatisfied" {
				t.Fatalf("read statuses = %s/%s, want %s/unsatisfied regardless of guarded_actions", evaluations[0].Status, evaluations[1].Status, wantGate)
			}
		})
	}
}

// TestCompleteEnvelopeDeniedByGuardingObligationCommitsNothing drives the real
// mutation: the review Obligation guards completion, so envelope.complete is
// denied naming it and leaves the journal and Envelope status unchanged. The
// terminal-success half of the same denial is asserted by the table above.
func TestCompleteEnvelopeDeniedByGuardingObligationCommitsNothing(t *testing.T) {
	fixture := newObligationGuardFixture(t, obligationGuardCases()[0])
	before := journalSnapshot(t, fixture.options.WorkspaceRoot)
	_, err := CompleteEnvelope(EnvelopeOptions{WorkspaceRoot: fixture.options.WorkspaceRoot, MigrationsDir: fixture.options.MigrationsDir, ActorID: fixture.options.ActorID, IdempotencyKey: "guard-complete-denied"}, EnvelopeCompleteRequest{EnvelopeID: "ENV-recovery", ExpectedRevision: 1, Repository: fixture.repository})
	var denial *EnvelopeCompletionError
	if !errors.As(err, &denial) {
		t.Fatalf("CompleteEnvelope() error = %v, want EnvelopeCompletionError", err)
	}
	if got := obligationReasonSubjects(denial.Reasons); !reflect.DeepEqual(got, []string{"OBL-review"}) {
		t.Fatalf("completion denial names Obligations %v, want [OBL-review]", got)
	}
	if after := journalSnapshot(t, fixture.options.WorkspaceRoot); !reflect.DeepEqual(before, after) {
		t.Fatalf("denied completion changed the journal: %v -> %v", before, after)
	}
	var status string
	if err := fixture.conn.QueryRow(`SELECT status FROM work_envelopes WHERE id = 'ENV-recovery'`).Scan(&status); err != nil || status != "active" {
		t.Fatalf("Envelope status after denied completion = %q, %v; want active", status, err)
	}
}

type obligationGuardFixture struct {
	options    RunOptions
	conn       *sql.DB
	envelope   Envelope
	repository Repository
}

func newObligationGuardFixture(t *testing.T, test obligationGuardCase) obligationGuardFixture {
	t.Helper()
	options, _ := recoveryFixture(t)
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	t.Cleanup(func() { closeWorkflowDB(conn) })
	repository, err := resumeObservedRepository(options.WorkspaceRoot, options.WorkspaceRoot)
	if err != nil {
		t.Fatalf("observe repository: %v", err)
	}
	gate := testObligation("OBL-gate", "gate-evidence", "ISS-recovery")
	gate.GuardedActions = test.gateGuards
	review := testObligation("OBL-review", "review", "ISS-recovery")
	review.ReviewProfile = stringPointer("review-profile")
	review.GuardedActions = test.reviewGuards
	now := time.Now().UTC()
	timestamp := now.Format(time.RFC3339Nano)
	repositoryJSON := string(mustRunJSON(repository))
	statements := []struct {
		query string
		args  []any
	}{
		{`UPDATE activation_actor_profiles SET capabilities_json = '["envelope.complete","issue.transition"]' WHERE actor_id = 'agent:recovery'`, nil},
		{`UPDATE envelope_authority_grants SET actor_selector_json = '{"actor_id":"agent:recovery","required_capabilities":["envelope.complete","issue.transition"]}', actions_json = '["envelope.complete","issue.transition"]' WHERE envelope_id = 'ENV-recovery'`, nil},
		{`UPDATE claims SET status = 'active', expires_at = ? WHERE id = 'CLM-prior'`, []any{now.Add(time.Hour).Format(time.RFC3339Nano)}},
		{`INSERT INTO envelope_obligations (id, envelope_id, revision, obligation_type, registration_id, registration_version, review_profile, outcome_class, obligation_json) VALUES ('OBL-gate', 'ENV-recovery', 1, 'gate-evidence', 'REG-test', 1, NULL, NULL, ?)`, []any{string(mustRunJSON(gate))}},
		{`INSERT INTO envelope_obligations (id, envelope_id, revision, obligation_type, registration_id, registration_version, review_profile, outcome_class, obligation_json) VALUES ('OBL-review', 'ENV-recovery', 1, 'review', NULL, NULL, 'review-profile', NULL, ?)`, []any{string(mustRunJSON(review))}},
	}
	if test.gateSatisfied {
		statements = append(statements,
			struct {
				query string
				args  []any
			}{`INSERT INTO execution_receipts (id, attempt_id, trust_source, terminal_status, exit_code, signal, timed_out, started_at, ended_at, terminal_repository_json, toolchain_json, output_excerpt, output_digest, output_bytes, output_truncated, redaction_applied, rejected_artifacts_json) VALUES ('RCP-gate', 'ATT-gate', 'fid-executed', 'succeeded', 0, NULL, 0, ?, ?, ?, '{}', '', 'sha256:0', 0, 0, 0, '[]')`, []any{timestamp, timestamp, repositoryJSON}},
			struct {
				query string
				args  []any
			}{`INSERT INTO evidence_records (id, obligation_id, evidence_type, receipt_id, review_id, human_outcome_id, producer, trust_source, target_repository_json, validity, created_at) VALUES ('EVD-gate', 'OBL-gate', 'receipt', 'RCP-gate', NULL, NULL, 'fid', 'fid-executed', ?, 'valid', ?)`, []any{repositoryJSON, timestamp}},
		)
	}
	if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("disable fixture foreign keys: %v", err)
	}
	for index, statement := range statements {
		if _, err := conn.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed guard fixture statement %d: %v", index, err)
		}
	}
	if _, err := conn.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		t.Fatalf("enable fixture foreign keys: %v", err)
	}
	envelope, err := loadEnvelope(conn, "ENV-recovery")
	if err != nil {
		t.Fatalf("load guard fixture Envelope: %v", err)
	}
	if len(envelope.Obligations) != 2 {
		t.Fatalf("guard fixture Envelope Obligations = %#v, want 2", envelope.Obligations)
	}
	return obligationGuardFixture{options: options, conn: conn, envelope: envelope, repository: repository}
}

func (fixture obligationGuardFixture) terminalReasons(t *testing.T) []EvaluationReason {
	t.Helper()
	version := 1
	result, err := EvaluateExecutionTransition(fixture.options.WorkspaceRoot, fixture.options.MigrationsDir, ExecutionTransitionRequest{IssueID: "ISS-recovery", ActorID: fixture.options.ActorID, ExpectedIssueRevision: &version, ActualIssueRevision: version, WorkflowEdgeValid: true, TerminalSuccess: true})
	if err != nil || !result.Bound {
		t.Fatalf("EvaluateExecutionTransition() = %#v, %v; want a bound evaluation", result, err)
	}
	return result.Reasons
}

func (fixture obligationGuardFixture) completeReasons(t *testing.T) []EvaluationReason {
	t.Helper()
	options := EnvelopeOptions{WorkspaceRoot: fixture.options.WorkspaceRoot, MigrationsDir: fixture.options.MigrationsDir, ActorID: fixture.options.ActorID, IdempotencyKey: "guard-complete"}
	reasons, evaluations, err := evaluateEnvelopeCompletion(options, fixture.conn, fixture.envelope, EnvelopeCompleteRequest{EnvelopeID: "ENV-recovery", ExpectedRevision: 1, Repository: fixture.repository})
	if err != nil {
		t.Fatalf("evaluateEnvelopeCompletion() error = %v", err)
	}
	if len(evaluations) != 2 {
		t.Fatalf("completion evaluations = %#v, want every Obligation reported", evaluations)
	}
	return reasons
}

// obligationReasonSubjects returns the sorted, distinct Obligation IDs named by
// Obligation-typed reasons (an Obligation-scoped reason carries its own ID).
func obligationReasonSubjects(reasons []EvaluationReason) []string {
	seen := map[string]bool{}
	var subjects []string
	for _, reason := range reasons {
		if reason.SubjectType == "obligation" && !seen[reason.SubjectID] {
			seen[reason.SubjectID] = true
			subjects = append(subjects, reason.SubjectID)
		}
	}
	sort.Strings(subjects)
	return subjects
}

func journalSnapshot(t *testing.T, workspace string) map[string]int64 {
	t.Helper()
	snapshot := map[string]int64{}
	root := filepath.Join(workspace, ".fid", "journal")
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relative, _ := filepath.Rel(root, path)
			snapshot[relative] = info.Size()
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("snapshot journal: %v", err)
	}
	return snapshot
}
