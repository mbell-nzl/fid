package review

import (
	"database/sql"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func openJudgmentMutation(lock *journal.WriteLock, options JudgmentOptions) (*sql.DB, string, workflow.ActorProfile, error) {
	if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
		return nil, "", workflow.ActorProfile{}, err
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
	if err != nil {
		return nil, "", workflow.ActorProfile{}, err
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		_ = conn.Close()
		return nil, "", workflow.ActorProfile{}, err
	}
	actor, err := workflow.ResolveTrustedActor(options.WorkspaceRoot, options.MigrationsDir, options.ActorID)
	if err != nil {
		_ = conn.Close()
		return nil, "", workflow.ActorProfile{}, err
	}
	return conn, projectConfig.ProjectID, actor, nil
}

func authorizeJudgmentSubject(options JudgmentOptions, conn *sql.DB, actor workflow.ActorProfile, action string, subject Subject) error {
	if !contains(actor.Capabilities, action) {
		return &policy.EnvelopeAuthorityDeniedError{ActorID: options.ActorID, Action: action, Capability: action}
	}
	binding, found, err := judgmentEnvelopeBinding(conn, subject)
	if err != nil || !found {
		return err
	}
	return workflow.AuthorizeEnvelopeGrant(options.WorkspaceRoot, options.MigrationsDir, workflow.EnvelopeAuthorityRequest{ActorID: options.ActorID, Action: action, Capability: action, IssueID: binding.IssueID, EnvelopeID: binding.EnvelopeID, Revision: binding.Revision})
}

type judgmentBinding struct {
	IssueID, EnvelopeID string
	Revision            int
}

func judgmentEnvelopeBinding(conn *sql.DB, subject Subject) (judgmentBinding, bool, error) {
	queries := map[string]string{
		"issue":    `SELECT s.issue_id,s.envelope_id,s.revision FROM envelope_issue_selection s JOIN work_envelopes e ON e.id=s.envelope_id AND e.current_revision=s.revision WHERE s.issue_id=? AND e.status IN ('draft','active') ORDER BY CASE e.status WHEN 'active' THEN 0 ELSE 1 END LIMIT 1`,
		"envelope": `SELECT s.issue_id,s.envelope_id,s.revision FROM envelope_issue_selection s JOIN work_envelopes e ON e.id=s.envelope_id AND e.current_revision=s.revision WHERE s.envelope_id=? LIMIT 1`,
		"run":      `SELECT issue_id,envelope_id,envelope_revision FROM runs WHERE id=?`,
		"review":   `SELECT issue_id,envelope_id,envelope_revision FROM review_records WHERE id=?`,
		"finding":  `SELECT r.issue_id,r.envelope_id,r.envelope_revision FROM review_findings f JOIN review_records r ON r.id=f.review_id WHERE f.id=?`,
	}
	query, ok := queries[subject.Type]
	if !ok {
		return judgmentBinding{}, false, nil
	}
	var value judgmentBinding
	err := conn.QueryRow(query, subject.ID).Scan(&value.IssueID, &value.EnvelopeID, &value.Revision)
	if err == sql.ErrNoRows {
		return judgmentBinding{}, false, nil
	}
	return value, err == nil, err
}

func setJudgmentEventIdentity(effects []journal.Effect, event journal.Event) {
	occurredAt := event.OccurredAt.UTC().Format(time.RFC3339Nano)
	for index := range effects {
		for key, value := range effects[index].Values {
			switch value {
			case "<event-id>":
				effects[index].Values[key] = event.EventID
			case "<occurred-at>":
				effects[index].Values[key] = occurredAt
			case "<actor-id>":
				effects[index].Values[key] = event.ActorV2.ActorID
			}
		}
	}
}

func judgmentActor(actor workflow.ActorProfile) journal.EventActor {
	return journal.EventActor{ActorID: actor.ActorID, ActorType: actor.ActorType, Producer: actor.Producer, ProvenanceRef: actor.ProvenanceRef, SessionRef: actor.SessionRef}
}

func rowExists(conn *sql.DB, table, id string) (bool, error) {
	var found int
	err := conn.QueryRow("SELECT 1 FROM "+table+" WHERE id=?", id).Scan(&found)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func subjectFromValue(value map[string]any) Subject {
	subject := value["subject"].(map[string]any)
	return Subject{Type: subject["type"].(string), ID: subject["id"].(string)}
}

func nullableValue(value sql.NullString) any {
	if value.Valid {
		return value.String
	}
	return nil
}

func validJudgmentOptions(options JudgmentOptions) bool {
	return strings.TrimSpace(options.WorkspaceRoot) != "" && strings.TrimSpace(options.MigrationsDir) != "" && strings.TrimSpace(options.ActorID) != "" && strings.TrimSpace(options.IdempotencyKey) != ""
}

func validJudgmentSubject(subject Subject) bool {
	return strings.TrimSpace(subject.ID) != "" && oneOf(subject.Type, "issue", "envelope", "run", "review", "finding", "requirement")
}

func validScope(scope Scope) bool {
	return len(scope.Paths) <= 256 && sort.StringsAreSorted(scope.Paths) && uniqueStrings(scope.Paths)
}

func validText(value string) bool { return strings.TrimSpace(value) != "" && len(value) <= 4096 }

func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func uniqueStrings(values []string) bool {
	for index := 1; index < len(values); index++ {
		if values[index-1] >= values[index] {
			return false
		}
	}
	return true
}
