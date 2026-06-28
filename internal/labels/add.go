package labels

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

const (
	actionLabelAdd = "label.add"
	fidVersion     = "0.1.0"
)

type DuplicateLabelError struct {
	Name string
}

func (err *DuplicateLabelError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("label add: label %q already exists", err.Name)
}

func Add(options AddOptions) (AddResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return AddResult{}, fmt.Errorf("label add: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return AddResult{}, fmt.Errorf("label add: migrations dir is required")
	}

	name := strings.TrimSpace(options.Name)
	if name == "" {
		return AddResult{}, fmt.Errorf("label add: label name is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionLabelAdd); err != nil {
		return AddResult{}, err
	}

	colour := trimmedOptionalString(options.Colour)
	description := trimmedOptionalString(options.Description)
	commandArgv := labelAddArgv(name, colour, description, options.ActorRole)

	var result AddResult
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
		defer closeLabelsConn(conn)

		if err := requireMissingLabel(conn, name); err != nil {
			return err
		}

		effectValues := map[string]any{}
		if colour != nil {
			effectValues["colour"] = *colour
		}
		if description != nil {
			effectValues["description"] = *description
		}

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "label.added",
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
					Table:  "labels",
					PK:     map[string]any{"name": name},
					Values: effectValues,
				},
			},
		})
		if err != nil {
			return err
		}

		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
		}

		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}
		if err := projectLabelsEvent(options.WorkspaceRoot, conn, appended); err != nil {
			return err
		}

		result = AddResult{
			Label: Label{
				Name:        name,
				Colour:      colour,
				Description: description,
			},
		}
		return nil
	})
	if err != nil {
		return AddResult{}, err
	}

	return result, nil
}

func requireMissingLabel(conn *sql.DB, name string) error {
	var existing string
	err := conn.QueryRow(`SELECT name FROM labels WHERE name = ? LIMIT 1`, name).Scan(&existing)
	if err == nil {
		return &DuplicateLabelError{Name: name}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return fmt.Errorf("query label %q: %w", name, err)
}

func labelAddArgv(name string, colour, description *string, actorRole string) []string {
	argv := []string{"fid", "label", "add", name}
	if colour != nil {
		argv = append(argv, "--colour", *colour)
	}
	if description != nil {
		argv = append(argv, "--description", *description)
	}
	return append(argv, "--by", actorRole)
}

func projectLabelsEvent(workspaceRoot string, conn *sql.DB, event journal.Event) error {
	eventHash, err := journal.HashEvent(event)
	if err != nil {
		return fmt.Errorf("hash label event %q: %w", event.EventID, err)
	}

	journalFile := event.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
	journalPath := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
	info, err := os.Stat(journalPath)
	if err != nil {
		return fmt.Errorf("stat label journal %q: %w", journalPath, err)
	}

	return journal.ProjectEvent(conn, journal.ProjectionInput{
		Event:     event,
		EventHash: eventHash,
		File:      journalFile,
		Line:      countJournalLines(journalPath),
		Offset:    info.Size(),
	})
}

func countJournalLines(path string) int {
	payload, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	trimmed := strings.TrimSpace(string(payload))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func trimmedOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func closeLabelsConn(conn *sql.DB) {
	_ = conn.Close()
}
