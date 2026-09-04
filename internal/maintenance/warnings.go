package maintenance

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/workflow"
)

type claimRecord struct {
	id, issueID, actorID, expiresAt, status string
}

func checkWarnings(db *sql.DB, now time.Time, graceSeconds int) ([]Finding, error) {
	claims, err := loadClaims(db)
	if err != nil {
		return nil, err
	}
	findings, err := staleClaimFindings(claims, now)
	if err != nil {
		return nil, err
	}
	started, err := startedWithoutClaimFindings(db, claims, now, graceSeconds)
	if err != nil {
		return nil, err
	}
	verifications, err := staleVerificationFindings(db)
	if err != nil {
		return nil, err
	}
	telemetry, err := telemetryGapFindings(db)
	if err != nil {
		return nil, err
	}
	findings = append(findings, started...)
	findings = append(findings, verifications...)
	return append(findings, telemetry...), nil
}

func loadClaims(db *sql.DB) ([]claimRecord, error) {
	rows, err := db.Query(`SELECT id, issue_id, actor_id, expires_at, status FROM claims ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("check claims: query claims: %w", err)
	}
	defer func() { _ = rows.Close() }()

	claims := make([]claimRecord, 0)
	for rows.Next() {
		var claim claimRecord
		if err := rows.Scan(&claim.id, &claim.issueID, &claim.actorID, &claim.expiresAt, &claim.status); err != nil {
			return nil, fmt.Errorf("check claims: scan claim: %w", err)
		}
		claims = append(claims, claim)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("check claims: iterate claims: %w", err)
	}
	return claims, nil
}

func staleClaimFindings(claims []claimRecord, now time.Time) ([]Finding, error) {
	findings := make([]Finding, 0)
	for _, claim := range claims {
		live, expiresAt, err := claim.liveAt(now)
		if err != nil {
			return nil, err
		}
		if !live {
			findings = append(findings, Finding{
				Check:    "stale_claims",
				EntityID: claim.issueID,
				Message:  fmt.Sprintf("claim on %s by %s expired at %s", claim.issueID, claim.actorID, expiresAt.Format(time.RFC3339)),
			})
		}
	}
	return findings, nil
}

func (claim claimRecord) liveAt(now time.Time) (bool, time.Time, error) {
	expiresAt, err := time.Parse(time.RFC3339, claim.expiresAt)
	if err != nil {
		return false, time.Time{}, fmt.Errorf("check claim %s time: %w", claim.issueID, err)
	}
	expiresAt = expiresAt.UTC()
	return claim.status == "active" && now.Before(expiresAt), expiresAt, nil
}

func startedWithoutClaimFindings(db *sql.DB, claims []claimRecord, now time.Time, graceSeconds int) ([]Finding, error) {
	registry, err := workflow.Load(db)
	if err != nil {
		return nil, fmt.Errorf("check started-without-claim: load registry: %w", err)
	}
	liveClaims, err := liveIssueClaims(claims, now)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`
		SELECT issue.id, issue.workflow_id, issue.status, issue.created_at,
			(
				SELECT occurred_at FROM activity_log
				WHERE ref_type = 'issue' AND ref_id = issue.id AND event = 'status_changed'
				ORDER BY rowid DESC LIMIT 1
			)
		FROM issues AS issue
		ORDER BY issue.id
	`)
	if err != nil {
		return nil, fmt.Errorf("check started-without-claim: query issues: %w", err)
	}
	defer func() { _ = rows.Close() }()

	findings := make([]Finding, 0)
	grace := time.Duration(graceSeconds) * time.Second
	for rows.Next() {
		var id, workflowID, status, createdAt string
		var transitionAt sql.NullString
		if err := rows.Scan(&id, &workflowID, &status, &createdAt, &transitionAt); err != nil {
			return nil, fmt.Errorf("check started-without-claim: scan issue: %w", err)
		}
		category, ok := registry.CategoryForState(workflowID, status)
		if !ok || category != "started" || liveClaims[id] {
			continue
		}
		stateTime := createdAt
		if transitionAt.Valid {
			stateTime = transitionAt.String
		}
		startedAt, err := time.Parse(time.RFC3339, stateTime)
		if err != nil {
			return nil, fmt.Errorf("check started-without-claim issue %s time: %w", id, err)
		}
		if now.Sub(startedAt.UTC()) >= grace {
			findings = append(findings, Finding{
				Check:    "started_without_claim",
				EntityID: id,
				Message:  fmt.Sprintf("started issue %s has no live claim since %s", id, startedAt.UTC().Format(time.RFC3339)),
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("check started-without-claim: iterate issues: %w", err)
	}
	return findings, nil
}

func liveIssueClaims(claims []claimRecord, now time.Time) (map[string]bool, error) {
	live := make(map[string]bool)
	for _, claim := range claims {
		isLive, _, err := claim.liveAt(now)
		if err != nil {
			return nil, err
		}
		if isLive {
			live[claim.issueID] = true
		}
	}
	return live, nil
}

func staleVerificationFindings(db *sql.DB) ([]Finding, error) {
	rows, err := db.Query(`
		SELECT verification.id, verification.completion_id
		FROM ac_verifications AS verification
		JOIN task_completions AS completion ON completion.id = verification.completion_id
		WHERE completion.status = 'superseded'
		ORDER BY verification.id
	`)
	if err != nil {
		return nil, fmt.Errorf("check stale verifications: query verifications: %w", err)
	}
	defer func() { _ = rows.Close() }()

	findings := make([]Finding, 0)
	for rows.Next() {
		var id, completionID string
		if err := rows.Scan(&id, &completionID); err != nil {
			return nil, fmt.Errorf("check stale verifications: scan verification: %w", err)
		}
		findings = append(findings, Finding{Check: "stale_verifications", EntityID: id, Message: fmt.Sprintf("AC verification %s references superseded completion %s", id, completionID)})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("check stale verifications: iterate verifications: %w", err)
	}
	return findings, nil
}

func telemetryGapFindings(db *sql.DB) ([]Finding, error) {
	rows, err := db.Query(`
		SELECT completion.id
		FROM task_completions AS completion
		LEFT JOIN agent_runs AS run ON run.completion_id = completion.id
		WHERE run.id IS NULL
		ORDER BY completion.id
	`)
	if err != nil {
		return nil, fmt.Errorf("check telemetry gaps: query completions: %w", err)
	}
	defer func() { _ = rows.Close() }()

	findings := make([]Finding, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("check telemetry gaps: scan completion: %w", err)
		}
		findings = append(findings, Finding{Check: "telemetry_gaps", EntityID: id, Message: fmt.Sprintf("completion %s has no completion-linked agent run", id)})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("check telemetry gaps: iterate completions: %w", err)
	}
	return findings, nil
}
