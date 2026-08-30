package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
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
	Mode          string  `json:"mode"`
	Version       int     `json:"version"`
	Status        string  `json:"status"`
	AcquiredAt    string  `json:"acquired_at"`
	ExpiresAt     string  `json:"expires_at"`
	ReleasedAt    *string `json:"released_at"`
	ReleaseReason *string `json:"release_reason"`
}

type ResumeCompletion struct {
	Status  string             `json:"status"`
	Reasons []EvaluationReason `json:"reasons"`
}

type Resume struct {
	Envelope    Envelope               `json:"envelope"`
	Issue       ResumeIssue            `json:"issue"`
	Claim       *ResumeClaim           `json:"claim"`
	Run         *Run                   `json:"run"`
	Repository  Repository             `json:"repository"`
	Obligations []ObligationEvaluation `json:"obligations"`
	Completion  ResumeCompletion       `json:"completion"`
	NextActions []string               `json:"next_actions"`
}

func ResumeIssueState(workspaceRoot, migrationsDir, issueID string) (Resume, error) {
	if strings.TrimSpace(issueID) == "" {
		return Resume{}, &EnvelopeInvalidError{Reason: "resume issue id is required"}
	}
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return Resume{}, err
	}
	conn, err := openWorkflowDB(workspaceRoot)
	if err != nil {
		return Resume{}, err
	}
	defer closeWorkflowDB(conn)
	return resumeIssueState(conn, issueID)
}

func resumeIssueState(conn *sql.DB, issueID string) (Resume, error) {
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
	evaluation, err := currentEnvelopeEvaluation(conn, envelope, envelope.RepositoryScope)
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
	completion := ResumeCompletion{Status: "incomplete", Reasons: append([]EvaluationReason(nil), evaluation.Reasons...)}
	if !terminal {
		completion.Reasons = append(completion.Reasons, evaluationReason("issue-not-terminal-success", "completed workflow category", issue.ID, "issue"))
	}
	completion.Reasons = sortEvaluationReasons(completion.Reasons)
	if envelope.Status == "completed" && len(completion.Reasons) == 0 {
		completion.Status = "complete"
	}
	return Resume{Envelope: envelope, Issue: issue, Claim: claim, Run: run, Repository: envelope.RepositoryScope, Obligations: []ObligationEvaluation{evaluation}, Completion: completion, NextActions: resumeNextActions(envelope, terminal, evaluation, claim, run)}, nil
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

func resumeNextActions(envelope Envelope, terminal bool, evaluation ObligationEvaluation, claim *ResumeClaim, run *Run) []string {
	if envelope.Status == "completed" {
		return []string{}
	}
	actions := []string{}
	if len(evaluation.Reasons) > 0 {
		actions = append(actions, "evaluate-obligations", "run-gate", "link-evidence")
	}
	if !terminal {
		actions = append(actions, "transition-issue")
	}
	if run != nil && run.Status == "active" {
		actions = append(actions, "finish-run")
	}
	if claim != nil {
		actions = append(actions, "release-claim")
	}
	if terminal && len(evaluation.Reasons) == 0 {
		actions = append(actions, "complete-envelope")
	}
	return actions
}
