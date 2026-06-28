package issues

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

type AddLabelOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	Label         string
	ActorRole     string
}

type AddLabelResult struct {
	IssueID string
	Label   string
}

type RemoveLabelOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	Label         string
	ActorRole     string
}

type RemoveLabelResult struct {
	IssueID string
	Label   string
}

func AddLabel(options AddLabelOptions) (AddLabelResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return AddLabelResult{}, fmt.Errorf("issue label add: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return AddLabelResult{}, fmt.Errorf("issue label add: migrations dir is required")
	}

	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return AddLabelResult{}, fmt.Errorf("issue label add: issue id is required")
	}
	label := strings.TrimSpace(options.Label)
	if label == "" {
		return AddLabelResult{}, fmt.Errorf("issue label add: label is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionIssueLabelAdd); err != nil {
		return AddLabelResult{}, err
	}

	commandArgv := issueLabelAddArgv(issueID, label, options.ActorRole)
	var result AddLabelResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		if err := assertIssueLabelAddable(conn, issueID, label); err != nil {
			return err
		}

		activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
			return activityIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "issue.label.added",
			Actor: journal.Actor{
				Role: options.ActorRole,
			},
			Command: journal.Command{
				Argv:       commandArgv,
				CWD:        ".",
				FIDVersion: fidVersion,
			},
		})
		if err != nil {
			return err
		}

		timestamp := event.OccurredAt.UTC().Format(time.RFC3339)
		event.Effects = []journal.Effect{
			{
				Op:    "insert",
				Table: "issue_labels",
				PK: map[string]any{
					"issue_id": issueID,
					"label":    label,
				},
				Values: map[string]any{},
			},
			issueLabelActivityEffect(activityID, event.EventID, options.ActorRole, timestamp, issueID, label, "added", "issue_label_added"),
		}

		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
		}

		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}
		if err := projectIssueEvent(conn, options.WorkspaceRoot, appended); err != nil {
			return err
		}

		result = AddLabelResult{IssueID: issueID, Label: label}
		return nil
	})
	if err != nil {
		return AddLabelResult{}, err
	}

	return result, nil
}

func RemoveLabel(options RemoveLabelOptions) (RemoveLabelResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return RemoveLabelResult{}, fmt.Errorf("issue label remove: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return RemoveLabelResult{}, fmt.Errorf("issue label remove: migrations dir is required")
	}

	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return RemoveLabelResult{}, fmt.Errorf("issue label remove: issue id is required")
	}
	label := strings.TrimSpace(options.Label)
	if label == "" {
		return RemoveLabelResult{}, fmt.Errorf("issue label remove: label is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionIssueLabelRemove); err != nil {
		return RemoveLabelResult{}, err
	}

	commandArgv := issueLabelRemoveArgv(issueID, label, options.ActorRole)
	var result RemoveLabelResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		if err := assertIssueLabelRemovable(conn, issueID, label); err != nil {
			return err
		}

		activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
			return activityIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "issue.label.removed",
			Actor: journal.Actor{
				Role: options.ActorRole,
			},
			Command: journal.Command{
				Argv:       commandArgv,
				CWD:        ".",
				FIDVersion: fidVersion,
			},
		})
		if err != nil {
			return err
		}

		timestamp := event.OccurredAt.UTC().Format(time.RFC3339)
		event.Effects = []journal.Effect{
			{
				Op:    "delete",
				Table: "issue_labels",
				PK: map[string]any{
					"issue_id": issueID,
					"label":    label,
				},
			},
			issueLabelActivityEffect(activityID, event.EventID, options.ActorRole, timestamp, issueID, label, "removed", "issue_label_removed"),
		}

		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
		}

		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}
		if err := projectIssueEvent(conn, options.WorkspaceRoot, appended); err != nil {
			return err
		}

		result = RemoveLabelResult{IssueID: issueID, Label: label}
		return nil
	})
	if err != nil {
		return RemoveLabelResult{}, err
	}

	return result, nil
}

func assertIssueLabelAddable(conn *sql.DB, issueID, label string) error {
	if _, found, err := issueKindForID(conn, issueID); err != nil {
		return err
	} else if !found {
		return &NotFoundError{IssueID: issueID}
	}

	if err := requireExistingLabels(conn, []string{label}); err != nil {
		return fmt.Errorf("%s", strings.Replace(err.Error(), "create issue:", "issue label add:", 1))
	}

	labels, err := issueLabels(conn, issueID)
	if err != nil {
		return err
	}
	for _, existing := range labels {
		if existing == label {
			return fmt.Errorf("issue label add: issue %q already has label %q", issueID, label)
		}
	}
	return nil
}

func assertIssueLabelRemovable(conn *sql.DB, issueID, label string) error {
	if _, found, err := issueKindForID(conn, issueID); err != nil {
		return err
	} else if !found {
		return &NotFoundError{IssueID: issueID}
	}

	labels, err := issueLabels(conn, issueID)
	if err != nil {
		return err
	}
	for _, existing := range labels {
		if existing == label {
			return nil
		}
	}
	return fmt.Errorf("issue label remove: issue %q does not have label %q", issueID, label)
}

func issueLabelAddArgv(issueID, label, actorRole string) []string {
	return []string{"fid", "issue", "label", "add", issueID, label, "--by", actorRole}
}

func issueLabelRemoveArgv(issueID, label, actorRole string) []string {
	return []string{"fid", "issue", "label", "remove", issueID, label, "--by", actorRole}
}

func issueLabelActivityEffect(activityID, eventID, actorRole, occurredAt, issueID, label, direction, eventName string) journal.Effect {
	return journal.Effect{
		Op:    "activity",
		Table: "activity_log",
		PK:    map[string]any{"id": activityID},
		Values: map[string]any{
			"journal_event_id": eventID,
			"ref_type":         "issue",
			"ref_id":           issueID,
			"event":            eventName,
			"details": map[string]any{
				"label":     label,
				"direction": direction,
			},
			"actor":       actorRole,
			"occurred_at": occurredAt,
		},
	}
}
