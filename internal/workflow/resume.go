package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
)

type ResumeIssue struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	StatusCategory string `json:"status_category"`
	Revision       int    `json:"revision"`
}

type ResumeClaim struct {
	ID               string `json:"id"`
	IssueID          string `json:"issue_id"`
	EnvelopeID       string `json:"envelope_id"`
	EnvelopeRevision int    `json:"envelope_revision"`
	ActorID          string `json:"actor_id"`
	WorktreeID       string `json:"worktree_id"`
	Scope            struct {
		Paths []string `json:"paths"`
	} `json:"scope"`
	Mode            string  `json:"mode"`
	Version         int     `json:"version"`
	Status          string  `json:"status"`
	AcquiredAt      string  `json:"acquired_at"`
	ExpiresAt       string  `json:"expires_at"`
	ReleasedAt      *string `json:"released_at"`
	ReleaseReason   *string `json:"release_reason"`
	EffectiveStatus string  `json:"effective_status"`
}

type ResumeCompletion struct {
	Status  string             `json:"status"`
	Reasons []EvaluationReason `json:"reasons"`
}

// NextAction projects one available read or possible mutation.
// Empty Reasons means the current durable state permits the operation once its
// ordinary request values are supplied.
type NextAction struct {
	Operation            string             `json:"operation"`
	Subject              Subject            `json:"subject"`
	RequiredCapabilities []string           `json:"required_capabilities"`
	Reasons              []EvaluationReason `json:"reasons"`
}

type Resume struct {
	Envelope           Envelope               `json:"envelope"`
	Issue              ResumeIssue            `json:"issue"`
	Claim              *ResumeClaim           `json:"claim"`
	Run                *Run                   `json:"run"`
	Repository         Repository             `json:"repository"`
	Obligations        []ObligationEvaluation `json:"obligations"`
	Completion         ResumeCompletion       `json:"completion"`
	Stops              []Stop                 `json:"stops"`
	Attention          []Attention            `json:"attention"`
	Budgets            []Budget               `json:"budgets"`
	Recoveries         []Recovery             `json:"recoveries"`
	Findings           []map[string]any       `json:"findings"`
	Reviews            []map[string]any       `json:"reviews"`
	Concerns           []map[string]any       `json:"concerns"`
	Deviations         []map[string]any       `json:"deviations"`
	Decisions          []map[string]any       `json:"decisions"`
	HumanOutcomes      []map[string]any       `json:"human_outcomes"`
	ProtectedDecisions []ProtectedDecision    `json:"protected_decisions"`
	NextActions        []NextAction           `json:"next_actions"`
}

func ResumeIssueState(workspaceRoot, migrationsDir, issueID, repositoryRoot string) (Resume, error) {
	if strings.TrimSpace(issueID) == "" {
		return Resume{}, &EnvelopeInvalidError{Reason: "resume issue id is required"}
	}
	if strings.TrimSpace(repositoryRoot) == "" {
		return Resume{}, &EnvelopeInvalidError{Reason: "resume repository is required"}
	}
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return Resume{}, err
	}
	repository, err := resumeObservedRepository(workspaceRoot, repositoryRoot)
	if err != nil {
		return Resume{}, err
	}
	conn, err := openWorkflowDB(workspaceRoot)
	if err != nil {
		return Resume{}, err
	}
	defer closeWorkflowDB(conn)
	return resumeIssueState(conn, issueID, repository)
}

func resumeObservedRepository(workspaceRoot, repositoryRoot string) (Repository, error) {
	projectConfig, err := config.ReadProjectConfig(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
	if err != nil {
		return Repository{}, err
	}
	observation, err := project.ObserveRecordedRepository(repositoryRoot, projectConfig.ProjectID)
	return Repository(observation), err
}

func resumeIssueState(conn *sql.DB, issueID string, repository Repository) (Resume, error) {
	var envelopeID string
	err := conn.QueryRow(`
		SELECT envelope_id
		FROM envelope_issue_selection
		WHERE issue_id = ?
		ORDER BY envelope_id DESC
		LIMIT 1
	`, issueID).Scan(&envelopeID)
	if err == sql.ErrNoRows {
		return Resume{}, &EnvelopeNotFoundError{ID: issueID}
	}
	if err != nil {
		return Resume{}, fmt.Errorf("resume envelope: %w", err)
	}
	envelope, err := loadEnvelope(conn, envelopeID)
	if err != nil {
		return Resume{}, err
	}
	issue, terminal, err := resumeIssue(conn, issueID)
	if err != nil {
		return Resume{}, err
	}
	evaluations, err := currentEnvelopeEvaluations(conn, envelope, repository)
	if err != nil {
		return Resume{}, err
	}
	if len(evaluations) == 0 {
		return Resume{}, fmt.Errorf("resume has no Obligation evaluations")
	}
	primaryEvaluation, err := currentEnvelopeEvaluation(conn, envelope, repository)
	if err != nil {
		return Resume{}, err
	}
	claim, err := resumeClaim(conn, issueID)
	if err != nil {
		return Resume{}, err
	}
	run, err := resumeRun(conn, issueID, envelope.ID)
	if err != nil {
		return Resume{}, err
	}
	stops, attention, err := resumeStopsAndAttention(conn, envelope.ID, run)
	if err != nil {
		return Resume{}, err
	}
	recoveries, err := resumeRecoveries(conn, issueID, envelope.ID)
	if err != nil {
		return Resume{}, err
	}
	completion := ResumeCompletion{Status: "incomplete", Reasons: []EvaluationReason{}}
	for _, evaluation := range evaluations {
		completion.Reasons = append(completion.Reasons, evaluation.Reasons...)
	}
	if !terminal {
		completion.Reasons = append(completion.Reasons, evaluationReason("issue-not-terminal-success", "completed workflow category", issue.ID, "issue"))
	}
	completion.Reasons = sortEvaluationReasons(completion.Reasons)
	if envelope.Status == "completed" && len(completion.Reasons) == 0 {
		completion.Status = "complete"
	}
	budgets := append([]Budget(nil), envelope.Budgets...)
	if run != nil {
		budgets = append([]Budget(nil), run.Budgets...)
	}
	nextActions, err := resumeNextActions(conn, issueID, envelope, terminal, primaryEvaluation, claim, run, repository)
	if err != nil {
		return Resume{}, err
	}
	judgment, err := resumeJudgment(conn, issue.ID, envelope)
	if err != nil {
		return Resume{}, err
	}
	return Resume{Envelope: envelope, Issue: issue, Claim: claim, Run: run, Repository: repository, Obligations: evaluations, Completion: completion, Stops: stops, Attention: attention, Budgets: budgets, Recoveries: recoveries, Findings: judgment.Findings, Reviews: judgment.Reviews, Concerns: judgment.Concerns, Deviations: judgment.Deviations, Decisions: judgment.Decisions, HumanOutcomes: judgment.HumanOutcomes, ProtectedDecisions: judgment.ProtectedDecisions, NextActions: nextActions}, nil
}

const resumeJudgmentLimit = 100

type resumeJudgmentState struct {
	Findings, Reviews, Concerns, Deviations, Decisions, HumanOutcomes []map[string]any
	ProtectedDecisions                                                []ProtectedDecision
}

func resumeJudgment(conn *sql.DB, issueID string, envelope Envelope) (resumeJudgmentState, error) {
	state := resumeJudgmentState{ProtectedDecisions: append([]ProtectedDecision{}, envelope.ProtectedDecisions...)}
	var err error
	state.Findings, err = resumeJSONRows(conn, "findings", `
		SELECT json_object('id',f.id,'review_id',f.review_id,'subject',json_object('type',f.subject_type,'id',f.subject_id),
		'scope',json(f.scope_json),'severity',f.severity,'summary',f.summary,'blocking',json(CASE f.blocking WHEN 1 THEN 'true' ELSE 'false' END),
		'created_at',f.created_at,'dispositions',json(COALESCE((SELECT json_group_array(json_object('id',d.id,'action',d.action,
		'actor_id',d.actor_id,'qualifying_review_id',d.qualifying_review_id,'reason',d.reason,'disposed_at',d.disposed_at,
		'event_id',d.created_event_id)) FROM finding_dispositions d WHERE d.finding_id=f.id ORDER BY d.disposed_at,d.id),'[]')),
		'state',COALESCE((SELECT action FROM finding_dispositions d WHERE d.finding_id=f.id ORDER BY disposed_at DESC,id DESC LIMIT 1),'open'))
		FROM review_findings f JOIN review_records r ON r.id=f.review_id WHERE r.issue_id=?
		AND COALESCE((SELECT action FROM finding_dispositions d WHERE d.finding_id=f.id ORDER BY disposed_at DESC,id DESC LIMIT 1),'open')
		IN ('open','addressed') ORDER BY f.id LIMIT ?`, issueID, resumeJudgmentLimit+1)
	if err != nil {
		return state, err
	}
	state.Reviews, err = resumeJSONRows(conn, "reviews", `
		SELECT json_object('id',r.id,'run_id',r.run_id,'issue_id',r.issue_id,'envelope_id',r.envelope_id,
		'envelope_revision',r.envelope_revision,'actor_id',r.actor_id,'profile',r.profile,'repository',json(r.repository_json),
		'coverage',json(r.coverage_json),'verdict',r.verdict,'evidence_refs',json(r.evidence_refs_json),
		'receipt_ids',json(r.receipt_ids_json),'findings',json(COALESCE((SELECT json_group_array(json_object('id',f.id,
		'subject',json_object('type',f.subject_type,'id',f.subject_id),'scope',json(f.scope_json),'severity',f.severity,
		'summary',f.summary,'blocking',json(CASE f.blocking WHEN 1 THEN 'true' ELSE 'false' END),'review_id',r.id,
		'created_at',f.created_at,'dispositions',json(COALESCE((SELECT json_group_array(json_object('id',d.id,'action',d.action,
		'actor_id',d.actor_id,'qualifying_review_id',d.qualifying_review_id,'reason',d.reason,'disposed_at',d.disposed_at,
		'event_id',d.created_event_id)) FROM finding_dispositions d WHERE d.finding_id=f.id),'[]')),
		'state',COALESCE((SELECT action FROM finding_dispositions d WHERE d.finding_id=f.id ORDER BY disposed_at DESC,id DESC LIMIT 1),'open')))
		FROM review_findings f WHERE f.review_id=r.id),'[]')),
		'resolved_finding_ids',json(COALESCE((SELECT json_group_array(finding_id) FROM finding_dispositions d
		WHERE d.qualifying_review_id=r.id AND d.action='resolved'),'[]')),
		'independence',r.independence,'recorded_at',r.recorded_at,'event_id',r.created_event_id)
		FROM review_records r WHERE r.issue_id=? ORDER BY r.id LIMIT ?`, issueID, resumeJudgmentLimit+1)
	if err != nil {
		return state, err
	}
	state.Concerns, err = resumeJSONRows(conn, "concerns", `
		SELECT json_object('id',id,'subject',json_object('type',subject_type,'id',subject_id),'actor_id',actor_id,
		'scope',json(scope_json),'severity',severity,'summary',summary,'reasoning',reasoning,'references',json(references_json),
		'status',status,'disposition',CASE WHEN disposition_action IS NULL THEN NULL ELSE json_object('action',disposition_action,
		'issue_id',disposition_issue_id,'supersedes_concern_id',supersedes_concern_id,'supersedes_finding_id',supersedes_finding_id,
		'reason',disposition_reason,'actor_id',disposed_actor_id,'disposed_at',disposed_at,'event_id',disposed_event_id) END,
		'created_at',created_at,'event_id',created_event_id) FROM concerns c
		WHERE `+resumeJudgmentSubjectFilter("c")+` ORDER BY id LIMIT ?`,
		issueID, issueID, issueID, issueID, issueID, issueID, resumeJudgmentLimit+1)
	if err != nil {
		return state, err
	}
	state.Deviations, err = resumeJSONRows(conn, "deviations", `
		SELECT json_object('id',id,'subject',json_object('type',subject_type,'id',subject_id),'actor_id',actor_id,
		'expected',expected,'observed',observed,'authority',json_object('actor_id',authority_actor_id,'reference',authority_reference),
		'disposition',disposition,'reasoning',reasoning,'references',json(references_json),'supersedes_deviation_id',supersedes_deviation_id,
		'created_at',created_at,'event_id',created_event_id) FROM deviations d
		WHERE `+resumeJudgmentSubjectFilter("d")+` ORDER BY id LIMIT ?`,
		issueID, issueID, issueID, issueID, issueID, issueID, resumeJudgmentLimit+1)
	if err != nil {
		return state, err
	}
	state.Decisions, err = resumeJSONRows(conn, "decisions", `
		SELECT json_object('id',id,'class',class,'source',json(source_json),'authority',json_object('actor_id',authority_actor_id,
		'required_capabilities',json(authority_capabilities_json)),'outcome',outcome,'scope',json_object('issue_id',scope_issue_id),
		'created_at',created_at,'event_id',created_event_id) FROM decision_refs WHERE scope_issue_id=? ORDER BY id LIMIT ?`, issueID, resumeJudgmentLimit+1)
	if err != nil {
		return state, err
	}
	state.HumanOutcomes, err = resumeJSONRows(conn, "human_outcomes", `
		SELECT json_object('id',id,'actor_id',actor_id,'provenance',json_object('activation_id',activation_id,
		'activation_version',activation_version,'profile_ref',profile_ref),'class',class,
		'subject',json_object('type',subject_type,'id',subject_id),'result',result,'scope',json_object('issue_id',scope_issue_id),
		'reasoning',reasoning,'references',json(references_json),'created_at',created_at,'event_id',created_event_id)
		FROM human_outcomes WHERE scope_issue_id=? ORDER BY id LIMIT ?`, issueID, resumeJudgmentLimit+1)
	return state, err
}

func resumeJudgmentSubjectFilter(alias string) string {
	subjectType, subjectID := alias+".subject_type", alias+".subject_id"
	return fmt.Sprintf(`(%[1]s='issue' AND %[2]s=?)
		OR (%[1]s='envelope' AND EXISTS(SELECT 1 FROM envelope_issue_selection s WHERE s.envelope_id=%[2]s AND s.issue_id=?))
		OR (%[1]s='run' AND EXISTS(SELECT 1 FROM runs r WHERE r.id=%[2]s AND r.issue_id=?))
		OR (%[1]s='review' AND EXISTS(SELECT 1 FROM review_records r WHERE r.id=%[2]s AND r.issue_id=?))
		OR (%[1]s='finding' AND EXISTS(SELECT 1 FROM review_findings f JOIN review_records r ON r.id=f.review_id WHERE f.id=%[2]s AND r.issue_id=?))
		OR (%[1]s='requirement' AND EXISTS(SELECT 1 FROM requirement_issue_links l WHERE l.requirement_id=%[2]s AND l.issue_id=?))`, subjectType, subjectID)
}

func resumeJSONRows(conn *sql.DB, collection, query string, args ...any) ([]map[string]any, error) {
	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("resume %s: %w", collection, err)
	}
	defer func() { _ = rows.Close() }()
	items := []map[string]any{}
	for rows.Next() {
		var value []byte
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		var item map[string]any
		if err := json.Unmarshal(value, &item); err != nil {
			return nil, fmt.Errorf("resume %s JSON: %w", collection, err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(items) > resumeJudgmentLimit {
		return nil, &ResultLimitExceededError{Collection: collection, Limit: resumeJudgmentLimit, Observed: len(items)}
	}
	return items, nil
}

func resumeIssue(conn *sql.DB, issueID string) (ResumeIssue, bool, error) {
	var issue ResumeIssue
	var workflowID string
	if err := conn.QueryRow(`SELECT id, workflow_id, status, revision FROM issues WHERE id = ?`, issueID).Scan(&issue.ID, &workflowID, &issue.Status, &issue.Revision); err != nil {
		return ResumeIssue{}, false, fmt.Errorf("resume Issue: %w", err)
	}
	registry, err := Load(conn)
	if err != nil {
		return ResumeIssue{}, false, err
	}
	category, found := registry.CategoryForState(workflowID, issue.Status)
	if !found {
		return ResumeIssue{}, false, fmt.Errorf("resume Issue: workflow state is unknown")
	}
	issue.StatusCategory = category
	return issue, category == "completed", nil
}

func resumeClaim(conn *sql.DB, issueID string) (*ResumeClaim, error) {
	var claim ResumeClaim
	var scope string
	var releasedAt, releaseReason sql.NullString
	err := conn.QueryRow(`
		SELECT id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id, scope_json, mode, version, status, acquired_at, expires_at, released_at, release_reason
		FROM claims WHERE issue_id = ? AND status = 'active' ORDER BY id LIMIT 1
	`, issueID).Scan(&claim.ID, &claim.IssueID, &claim.EnvelopeID, &claim.EnvelopeRevision, &claim.ActorID, &claim.WorktreeID, &scope, &claim.Mode, &claim.Version, &claim.Status, &claim.AcquiredAt, &claim.ExpiresAt, &releasedAt, &releaseReason)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("resume Claim: %w", err)
	}
	if err := json.Unmarshal([]byte(scope), &claim.Scope); err != nil {
		return nil, fmt.Errorf("resume Claim scope: %w", err)
	}
	if releasedAt.Valid {
		claim.ReleasedAt = &releasedAt.String
	}
	if releaseReason.Valid {
		claim.ReleaseReason = &releaseReason.String
	}
	claim.EffectiveStatus = claim.Status
	if claim.Status == "active" {
		expiresAt, err := time.Parse(time.RFC3339Nano, claim.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("resume Claim expiry: %w", err)
		}
		if !time.Now().UTC().Before(expiresAt.UTC()) {
			claim.EffectiveStatus = "expired"
		}
	}
	return &claim, nil
}

func resumeRun(conn *sql.DB, issueID, envelopeID string) (*Run, error) {
	var runID string
	err := conn.QueryRow(`
		SELECT id FROM runs WHERE issue_id = ? AND envelope_id = ?
		ORDER BY CASE status WHEN 'active' THEN 0 ELSE 1 END, started_at DESC, id DESC LIMIT 1
	`, issueID, envelopeID).Scan(&runID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("resume Run: %w", err)
	}
	run, err := loadRun(conn, runID)
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func resumeStopsAndAttention(conn *sql.DB, envelopeID string, run *Run) ([]Stop, []Attention, error) {
	subjects := []Subject{{Type: "envelope", ID: envelopeID}}
	if run != nil {
		subjects = append(subjects, Subject{Type: "run", ID: run.ID})
	}
	activationIDs, err := resumeOpenActivationIDs(conn)
	if err != nil {
		return nil, nil, err
	}
	for _, activationID := range activationIDs {
		subjects = append(subjects, Subject{Type: "activation", ID: activationID})
	}
	stops := []Stop{}
	attention := []Attention{}
	seenAttention := map[string]bool{}
	for _, subject := range subjects {
		rows, err := conn.Query(`
			SELECT s.id FROM stops s
			LEFT JOIN stop_resolutions r ON r.stop_id = s.id
			WHERE s.subject_type = ? AND s.subject_id = ? AND r.stop_id IS NULL
			ORDER BY s.raised_at, s.id
		`, subject.Type, subject.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("resume Stops: %w", err)
		}
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				_ = rows.Close()
				return nil, nil, err
			}
			stop, err := loadStop(conn, id)
			if err != nil {
				_ = rows.Close()
				return nil, nil, err
			}
			stops = append(stops, stop)
			item, err := loadAttentionForStop(conn, id)
			if err != nil {
				_ = rows.Close()
				return nil, nil, err
			}
			if !seenAttention[item.ID] {
				seenAttention[item.ID] = true
				attention = append(attention, item)
			}
		}
		if err := rows.Close(); err != nil {
			return nil, nil, err
		}
		rows, err = conn.Query(`SELECT id FROM attention_items WHERE subject_type = ? AND subject_id = ? AND status = 'open' ORDER BY raised_at, id`, subject.Type, subject.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("resume Attention: %w", err)
		}
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				_ = rows.Close()
				return nil, nil, err
			}
			if seenAttention[id] {
				continue
			}
			item, err := loadAttentionByQuery(conn, `SELECT id, semantic_key, subject_type, subject_id, code, status, version, raised_at, resolution_event_id FROM attention_items WHERE id = ?`, id)
			if err != nil {
				_ = rows.Close()
				return nil, nil, err
			}
			seenAttention[id] = true
			attention = append(attention, item)
		}
		if err := rows.Close(); err != nil {
			return nil, nil, err
		}
	}
	return stops, attention, nil
}

func resumeOpenActivationIDs(conn *sql.DB) ([]string, error) {
	rows, err := conn.Query(`
		SELECT subject_id FROM (
			SELECT s.subject_id
			FROM stops s
			LEFT JOIN stop_resolutions r ON r.stop_id = s.id
			WHERE s.subject_type = 'activation' AND s.code = 'activation-incomplete' AND r.stop_id IS NULL
			UNION
			SELECT subject_id
			FROM attention_items
			WHERE subject_type = 'activation' AND status = 'open'
		)
		ORDER BY subject_id
	`)
	if err != nil {
		return nil, fmt.Errorf("resume activation Stops: %w", err)
	}
	defer func() { _ = rows.Close() }()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func resumeRecoveries(conn *sql.DB, issueID, envelopeID string) ([]Recovery, error) {
	rows, err := conn.Query(`
		SELECT r.id
		FROM run_recoveries r
		JOIN runs prior ON prior.id = r.prior_run_id
		WHERE prior.issue_id = ? AND prior.envelope_id = ?
		ORDER BY r.id
	`, issueID, envelopeID)
	if err != nil {
		return nil, fmt.Errorf("resume recoveries: %w", err)
	}
	defer func() { _ = rows.Close() }()
	recoveries := []Recovery{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		recovery, err := loadRecovery(conn, id)
		if err != nil {
			return nil, err
		}
		recoveries = append(recoveries, recovery)
	}
	return recoveries, rows.Err()
}
