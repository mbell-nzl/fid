package evidence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

type ReceiptNotFoundError struct{ ID string }

func (err *ReceiptNotFoundError) Error() string {
	return fmt.Sprintf("Receipt %q was not found", err.ID)
}

type ReceiptListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	RunID         string
	ObligationID  string
	Outcome       string
	Limit         int
}

func GetReceipt(workspaceRoot, migrationsDir, id string) (Receipt, error) {
	if strings.TrimSpace(workspaceRoot) == "" || strings.TrimSpace(migrationsDir) == "" || strings.TrimSpace(id) == "" {
		return Receipt{}, fmt.Errorf("receipt get input is invalid")
	}
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return Receipt{}, err
	}
	conn, err := openEvidenceDB(workspaceRoot)
	if err != nil {
		return Receipt{}, err
	}
	defer func() { _ = conn.Close() }()
	return loadReceipt(conn, strings.TrimSpace(id))
}

func ListReceipts(options ReceiptListOptions) ([]Receipt, error) {
	if !validReceiptListOptions(options) {
		return nil, fmt.Errorf("receipt list input is invalid")
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}
	conn, err := openEvidenceDB(options.WorkspaceRoot)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()

	query, args := `SELECT r.id FROM execution_receipts r JOIN execution_attempts a ON a.id = r.attempt_id WHERE 1 = 1`, []any{}
	if options.RunID != "" {
		query, args = query+` AND a.run_id = ?`, append(args, options.RunID)
	}
	if options.ObligationID != "" {
		query, args = query+` AND a.obligation_id = ?`, append(args, options.ObligationID)
	}
	if options.Outcome != "" {
		query, args = query+` AND r.terminal_status = ?`, append(args, options.Outcome)
	}
	rows, err := conn.Query(query+` ORDER BY r.started_at, r.id LIMIT ?`, append(args, options.Limit)...)
	if err != nil {
		return nil, fmt.Errorf("list Receipts: %w", err)
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("list Receipts: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("list Receipts: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("list Receipts: %w", err)
	}

	receipts := make([]Receipt, 0, len(ids))
	for _, id := range ids {
		receipt, err := loadReceipt(conn, id)
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}
	return receipts, nil
}

func validReceiptListOptions(options ReceiptListOptions) bool {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || options.Limit < 1 || options.Limit > 10000 {
		return false
	}
	switch options.Outcome {
	case "", "failed", "signaled", "spawn-failed", "succeeded", "timed-out":
		return true
	default:
		return false
	}
}

func loadReceipt(conn *sql.DB, id string) (Receipt, error) {
	var receipt Receipt
	var repositoryJSON, toolchainJSON, rejectedJSON string
	err := conn.QueryRow(`SELECT r.id, r.attempt_id, a.registration_id, a.registration_version, a.obligation_id, a.run_id, a.actor_id, r.trust_source, r.terminal_status, r.exit_code, r.signal, r.timed_out, r.started_at, r.ended_at, r.terminal_repository_json, r.toolchain_json, r.output_excerpt, r.output_digest, r.output_bytes, r.output_truncated, r.redaction_applied, r.rejected_artifacts_json FROM execution_receipts r JOIN execution_attempts a ON a.id = r.attempt_id WHERE r.id = ?`, id).Scan(&receipt.ID, &receipt.AttemptID, &receipt.RegistrationID, &receipt.RegistrationVersion, &receipt.ObligationID, &receipt.RunID, &receipt.ActorID, &receipt.TrustSource, &receipt.Outcome, &receipt.ExitCode, &receipt.Signal, &receipt.TimedOut, &receipt.StartedAt, &receipt.EndedAt, &repositoryJSON, &toolchainJSON, &receipt.OutputExcerpt, &receipt.OutputDigest, &receipt.OutputBytes, &receipt.OutputTruncated, &receipt.RedactionApplied, &rejectedJSON)
	if err == sql.ErrNoRows {
		return Receipt{}, &ReceiptNotFoundError{ID: id}
	}
	if err != nil {
		return Receipt{}, fmt.Errorf("load Receipt %q: %w", id, err)
	}
	if err := json.Unmarshal([]byte(repositoryJSON), &receipt.TerminalRepository); err != nil {
		return Receipt{}, fmt.Errorf("load Receipt repository: %w", err)
	}
	if err := json.Unmarshal([]byte(toolchainJSON), &receipt.Toolchain); err != nil {
		return Receipt{}, fmt.Errorf("load Receipt toolchain: %w", err)
	}
	if err := json.Unmarshal([]byte(rejectedJSON), &receipt.RejectedArtifacts); err != nil || receipt.RejectedArtifacts == nil {
		return Receipt{}, fmt.Errorf("load Receipt rejected Artifacts: invalid persisted array")
	}
	sortArtifactRejections(receipt.RejectedArtifacts)
	receipt.Artifacts, err = loadReceiptArtifacts(conn, id)
	if err != nil {
		return Receipt{}, err
	}
	return receipt, nil
}

func sortArtifactRejections(values []ArtifactRejection) {
	sort.Slice(values, func(i, j int) bool {
		if values[i].Path == values[j].Path {
			return values[i].Reason < values[j].Reason
		}
		return values[i].Path < values[j].Path
	})
}
