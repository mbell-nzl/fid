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
	name, colour, description, commandArgv, err := prepareLabelAdd(options)
	if err != nil {
		return AddResult{}, err
	}
	var result AddResult
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		return addLabelLocked(lock, options, name, colour, description, commandArgv, &result)
	})
	if err != nil {
		return AddResult{}, err
	}

	return result, nil
}

func prepareLabelAdd(options AddOptions) (string, *string, *string, []string, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return "", nil, nil, nil, fmt.Errorf("label add: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return "", nil, nil, nil, fmt.Errorf("label add: migrations dir is required")
	}

	name := strings.TrimSpace(options.Name)
	if name == "" {
		return "", nil, nil, nil, fmt.Errorf("label add: label name is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionLabelAdd); err != nil {
		return "", nil, nil, nil, err
	}

	colour := trimmedOptionalString(options.Colour)
	description := trimmedOptionalString(options.Description)
	return name, colour, description, labelAddArgv(name, colour, description, options.ActorRole), nil
}

func addLabelLocked(lock *journal.WriteLock, options AddOptions, name string, colour, description *string, commandArgv []string, result *AddResult) error {
	if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
		return err
	}

	conn, err := openLabelsProjectDB(options.WorkspaceRoot)
	if err != nil {
		return err
	}
	defer closeLabelsConn(conn)

	if err := requireMissingLabel(conn, name); err != nil {
		return err
	}

	event, err := prepareLabelAddEvent(lock, name, colour, description, options.ActorRole, commandArgv)
	if err != nil {
		return err
	}
	if err := appendAndProjectLabelEvent(options.WorkspaceRoot, lock, conn, event); err != nil {
		return err
	}

	*result = AddResult{Label: Label{Name: name, Colour: colour, Description: description}}
	return nil
}

func openLabelsProjectDB(workspaceRoot string) (*sql.DB, error) {
	return fiddb.Open(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
}

func prepareLabelAddEvent(lock *journal.WriteLock, name string, colour, description *string, actorRole string, commandArgv []string) (journal.Event, error) {
	return journal.PrepareEvent(lock, journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   "label.added",
		Actor: journal.Actor{
			Role: actorRole,
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
				Values: labelAddEffectValues(colour, description),
			},
		},
	})
}

func labelAddEffectValues(colour, description *string) map[string]any {
	values := map[string]any{}
	if colour != nil {
		values["colour"] = *colour
	}
	if description != nil {
		values["description"] = *description
	}
	return values
}

func appendAndProjectLabelEvent(workspaceRoot string, lock *journal.WriteLock, conn *sql.DB, event journal.Event) error {
	if err := journal.DryRunProjectEvent(conn, event); err != nil {
		return err
	}

	appended, err := journal.AppendEvent(lock, event)
	if err != nil {
		return err
	}
	return projectLabelsEvent(workspaceRoot, conn, appended)
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
