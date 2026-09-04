package maintenance

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type CheckOptions struct {
	DB                              *sql.DB
	JournalDir                      string
	Now                             time.Time
	StartedWithoutClaimGraceSeconds int
}

type Result struct {
	Errors   []Finding `json:"errors"`
	Warnings []Finding `json:"warnings"`
	Info     Info      `json:"info"`
}

type Info struct {
	IssueCount     int `json:"issue_count"`
	ClaimCountLive int `json:"claim_count_live"`
	RecentActivity int `json:"recent_activity_7d"`
}

type Finding struct {
	Check    string `json:"check"`
	Message  string `json:"message"`
	EntityID string `json:"entity_id"`
}

const MaxFindingsPerSeverity = 1000

type ResultLimitExceededError struct {
	Collection string
	Limit      int
	Observed   int
}

func (err *ResultLimitExceededError) Error() string {
	return fmt.Sprintf("maintenance result %s has %d items, maximum is %d", err.Collection, err.Observed, err.Limit)
}

type structuralCheck func(*sql.DB, string) ([]Finding, error)

func Check(options CheckOptions) (Result, error) {
	if options.DB == nil {
		return Result{}, fmt.Errorf("maintenance check: database is required")
	}
	if strings.TrimSpace(options.JournalDir) == "" {
		return Result{}, fmt.Errorf("maintenance check: journal dir is required")
	}
	if options.StartedWithoutClaimGraceSeconds < 0 {
		return Result{}, fmt.Errorf("maintenance check: started-without-claim grace must be non-negative")
	}

	now := options.Now.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}

	checks := []structuralCheck{
		checkKindHierarchy,
		checkDependentRowKinds,
		checkSequenceScheduling,
		checkLinkReferences,
		checkActivityReferences,
		checkForeignKeys,
		checkWorkflowStates,
		checkWorkflowAssignments,
		checkOperationalIntegrity,
	}
	result := Result{Errors: make([]Finding, 0), Warnings: make([]Finding, 0)}
	for _, check := range checks {
		findings, err := check(options.DB, options.JournalDir)
		if err != nil {
			return Result{}, err
		}
		result.Errors = append(result.Errors, findings...)
	}
	warnings, err := checkWarnings(options.DB, now, options.StartedWithoutClaimGraceSeconds)
	if err != nil {
		return Result{}, err
	}
	info, err := checkInfo(options.DB, now)
	if err != nil {
		return Result{}, err
	}
	result.Warnings = append(result.Warnings, warnings...)
	result.Info = info
	sortFindings(result.Errors)
	sortFindings(result.Warnings)
	if err := enforceFindingLimit("maintenance.errors", result.Errors); err != nil {
		return Result{}, err
	}
	if err := enforceFindingLimit("maintenance.warnings", result.Warnings); err != nil {
		return Result{}, err
	}
	return result, nil
}

func enforceFindingLimit(collection string, findings []Finding) error {
	if len(findings) <= MaxFindingsPerSeverity {
		return nil
	}
	return &ResultLimitExceededError{Collection: collection, Limit: MaxFindingsPerSeverity, Observed: len(findings)}
}

func sortFindings(findings []Finding) {
	sort.Slice(findings, func(left, right int) bool {
		if findings[left].Check != findings[right].Check {
			return findings[left].Check < findings[right].Check
		}
		if findings[left].EntityID != findings[right].EntityID {
			return findings[left].EntityID < findings[right].EntityID
		}
		return findings[left].Message < findings[right].Message
	})
}

func checkForeignKeys(db *sql.DB, _ string) ([]Finding, error) {
	rows, err := db.Query(`PRAGMA foreign_key_check`)
	if err != nil {
		return nil, fmt.Errorf("check foreign key integrity: query foreign_key_check: %w", err)
	}

	findings := make([]Finding, 0)
	for rows.Next() {
		var table, parent string
		var rowID, foreignKeyID int64
		if err := rows.Scan(&table, &rowID, &parent, &foreignKeyID); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("check foreign key integrity: scan foreign_key_check: %w", err)
		}
		entityID := fmt.Sprintf("%s:%d", table, rowID)
		findings = append(findings, Finding{
			Check:    "foreign_key_integrity",
			EntityID: entityID,
			Message:  fmt.Sprintf("foreign key %d on %s row %d references missing %s row", foreignKeyID, table, rowID, parent),
		})
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("check foreign key integrity: iterate foreign_key_check: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("check foreign key integrity: close foreign_key_check: %w", err)
	}
	return findings, nil
}

func checkActivityReferences(db *sql.DB, journalDir string) ([]Finding, error) {
	rows, err := db.Query(`SELECT id, journal_event_id, ref_type, ref_id FROM activity_log ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("check activity references: query activity log: %w", err)
	}

	activities := make([]activityRow, 0)
	for rows.Next() {
		var activity activityRow
		if err := rows.Scan(&activity.id, &activity.eventID, &activity.refType, &activity.refID); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("check activity references: scan activity row: %w", err)
		}
		activities = append(activities, activity)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("check activity references: iterate activity log: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("check activity references: close activity log: %w", err)
	}

	findings := make([]Finding, 0)
	knownEvents := make(map[string]bool)
	missingEvents := make(map[string]bool)
	for _, activity := range activities {
		if !validActivityType(activity.refType) {
			findings = append(findings, Finding{Check: "activity_refs", EntityID: activity.id, Message: fmt.Sprintf("activity %s has unsupported ref type %s", activity.id, activity.refType)})
		} else if exists, err := referencedEntityExists(db, activity.refType, activity.refID); err != nil {
			return nil, fmt.Errorf("check activity references: query activity %s ref: %w", activity.id, err)
		} else if !exists {
			findings = append(findings, Finding{Check: "activity_refs", EntityID: activity.id, Message: fmt.Sprintf("activity %s references missing %s:%s", activity.id, activity.refType, activity.refID)})
		}
		if missingEvents[activity.eventID] {
			findings = append(findings, missingJournalEventFinding(activity))
			continue
		}
		if !knownEvents[activity.eventID] {
			if err := validateActivityJournalEvent(journalDir, activity.eventID); err != nil {
				var notFound *journal.EventNotFoundError
				if !errors.As(err, &notFound) {
					return nil, fmt.Errorf("check activity references: inspect journal event %s: %w", activity.eventID, err)
				}
				findings = append(findings, missingJournalEventFinding(activity))
				missingEvents[activity.eventID] = true
				continue
			}
			knownEvents[activity.eventID] = true
		}
	}
	return findings, nil
}

type activityRow struct{ id, eventID, refType, refID string }

func missingJournalEventFinding(activity activityRow) Finding {
	return Finding{
		Check:    "activity_journal_events",
		EntityID: activity.id,
		Message:  fmt.Sprintf("activity %s references missing journal event %s", activity.id, activity.eventID),
	}
}

func validActivityType(refType string) bool {
	return refType == "issue" || refType == "requirement" || refType == "decision" || refType == "escalation"
}

func validateActivityJournalEvent(journalDir, eventID string) error {
	_, err := journal.FindEvent(journalDir, eventID)
	return err
}

func checkWorkflowStates(db *sql.DB, _ string) ([]Finding, error) {
	registry, err := workflow.Load(db)
	if err != nil {
		return nil, fmt.Errorf("check workflow states: load registry: %w", err)
	}

	rows, err := db.Query(`SELECT id, workflow_id, status FROM issues ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("check workflow states: query issues: %w", err)
	}

	findings := make([]Finding, 0)
	for rows.Next() {
		var id, workflowID, state string
		if err := rows.Scan(&id, &workflowID, &state); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("check workflow states: scan issue: %w", err)
		}
		if _, ok := registry.CategoryForState(workflowID, state); !ok {
			findings = append(findings, Finding{Check: "workflow_state", EntityID: id, Message: fmt.Sprintf("issue %s has unknown workflow/state pair %s/%s", id, workflowID, state)})
		}
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("check workflow states: iterate issues: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("check workflow states: close issues: %w", err)
	}
	return findings, nil
}
