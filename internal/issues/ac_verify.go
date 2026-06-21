package issues

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

var acceptanceVerificationStatuses = []string{"pass", "fail", "blocked", "not_applicable"}
var acceptanceVerificationEvidenceTypes = []string{"test", "source", "screenshot", "manual", "command", "commit"}

type VerifyAcceptanceCriterionOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	CriterionID   string
	ActorRole     string
	Status        string
	EvidenceType  *string
	EvidenceRef   *string
	Notes         *string
}

type VerifyAcceptanceCriterionResult struct {
	ID           string
	IssueID      string
	CriterionID  string
	ActorRole    string
	Status       string
	EvidenceType *string
	EvidenceRef  *string
	Notes        *string
	VerifiedAt   string
}

func VerifyAcceptanceCriterion(options VerifyAcceptanceCriterionOptions) (VerifyAcceptanceCriterionResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return VerifyAcceptanceCriterionResult{}, fmt.Errorf("verify acceptance criterion: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return VerifyAcceptanceCriterionResult{}, fmt.Errorf("verify acceptance criterion: migrations dir is required")
	}
	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return VerifyAcceptanceCriterionResult{}, fmt.Errorf("verify acceptance criterion: issue id is required")
	}
	criterionID := strings.TrimSpace(options.CriterionID)
	if criterionID == "" {
		return VerifyAcceptanceCriterionResult{}, fmt.Errorf("verify acceptance criterion: acceptance criterion id is required")
	}
	status := strings.TrimSpace(options.Status)
	if status == "" {
		return VerifyAcceptanceCriterionResult{}, fmt.Errorf("verify acceptance criterion: status is required")
	}
	if !slices.Contains(acceptanceVerificationStatuses, status) {
		return VerifyAcceptanceCriterionResult{}, fmt.Errorf(
			"verify acceptance criterion: status %q must be one of %s",
			status,
			strings.Join(acceptanceVerificationStatuses, ", "),
		)
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, issueACVerifyAction(options.ActorRole)); err != nil {
		return VerifyAcceptanceCriterionResult{}, err
	}

	evidenceType := trimmedOptionalString(options.EvidenceType)
	if evidenceType != nil && !slices.Contains(acceptanceVerificationEvidenceTypes, *evidenceType) {
		return VerifyAcceptanceCriterionResult{}, fmt.Errorf(
			"verify acceptance criterion: evidence type %q must be one of %s",
			*evidenceType,
			strings.Join(acceptanceVerificationEvidenceTypes, ", "),
		)
	}
	evidenceRef := trimmedOptionalString(options.EvidenceRef)
	notes := trimmedOptionalString(options.Notes)

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return VerifyAcceptanceCriterionResult{}, err
	}

	commandArgv := issueACVerifyArgv(issueID, criterionID, options.ActorRole, status, evidenceType, evidenceRef, notes)
	var result VerifyAcceptanceCriterionResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		dbPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName)
		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		if _, found, err := issueKindForID(conn, issueID); err != nil {
			return err
		} else if !found {
			return &NotFoundError{IssueID: issueID}
		}
		if err := requireAcceptanceCriterionBelongsToIssue(conn, issueID, criterionID); err != nil {
			return err
		}

		verificationID, err := journal.GenerateID("ACV", func(candidate string) (bool, error) {
			return acceptanceCriterionVerificationIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		verifiedAt := time.Now().UTC().Format(time.RFC3339)
		values := map[string]any{
			"ac_id":         criterionID,
			"verifier_role": options.ActorRole,
			"status":        status,
			"verified_at":   verifiedAt,
		}
		if evidenceType != nil {
			values["evidence_type"] = *evidenceType
		}
		if evidenceRef != nil {
			values["evidence_ref"] = *evidenceRef
		}
		if notes != nil {
			values["notes"] = *notes
		}

		event := journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "issue.ac_verified",
			Actor: journal.Actor{
				Role: options.ActorRole,
			},
			Command: journal.Command{
				Argv:       commandArgv,
				CWD:        ".",
				FIDVersion: fidVersion,
			},
			Effects: []journal.Effect{
				{
					Op:     "insert",
					Table:  "ac_verifications",
					PK:     map[string]any{"id": verificationID},
					Values: values,
				},
			},
		}

		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
		}

		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}

		eventHash, err := journal.HashEvent(appended)
		if err != nil {
			return fmt.Errorf("hash issue event %q: %w", appended.EventID, err)
		}

		journalFile := appended.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
		journalPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
		info, err := os.Stat(journalPath)
		if err != nil {
			return fmt.Errorf("stat issue journal %q: %w", journalPath, err)
		}

		if err := journal.ProjectEvent(conn, journal.ProjectionInput{
			Event:     appended,
			EventHash: eventHash,
			File:      journalFile,
			Line:      countJournalLines(journalPath),
			Offset:    info.Size(),
		}); err != nil {
			return err
		}

		result = VerifyAcceptanceCriterionResult{
			ID:           verificationID,
			IssueID:      issueID,
			CriterionID:  criterionID,
			ActorRole:    options.ActorRole,
			Status:       status,
			EvidenceType: evidenceType,
			EvidenceRef:  evidenceRef,
			Notes:        notes,
			VerifiedAt:   verifiedAt,
		}
		return nil
	})
	if err != nil {
		return VerifyAcceptanceCriterionResult{}, err
	}

	return result, nil
}

func requireAcceptanceCriterionBelongsToIssue(conn *sql.DB, issueID, criterionID string) error {
	var criterionIssueID string
	err := conn.QueryRow(`SELECT issue_id FROM acceptance_criteria WHERE id = ? LIMIT 1`, criterionID).Scan(&criterionIssueID)
	if err == nil {
		if criterionIssueID == issueID {
			return nil
		}
		return fmt.Errorf(
			"verify acceptance criterion: acceptance criterion %q does not belong to issue %q",
			criterionID,
			issueID,
		)
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("verify acceptance criterion: acceptance criterion %q not found", criterionID)
	}
	return fmt.Errorf("query acceptance criterion %q: %w", criterionID, err)
}

func acceptanceCriterionVerificationIDExists(conn *sql.DB, verificationID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM ac_verifications WHERE id = ? LIMIT 1`, verificationID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("query acceptance criterion verification id %q: %w", verificationID, err)
}

func issueACVerifyArgv(issueID, criterionID, actorRole, status string, evidenceType, evidenceRef, notes *string) []string {
	argv := []string{
		"fid", "issue", "ac", "verify", issueID, criterionID,
		"--by", actorRole,
		"--status", status,
	}
	if evidenceType != nil {
		argv = append(argv, "--evidence-type", *evidenceType)
	}
	if evidenceRef != nil {
		argv = append(argv, "--evidence-ref", *evidenceRef)
	}
	if notes != nil {
		argv = append(argv, "--notes", *notes)
	}
	return argv
}
