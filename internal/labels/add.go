package labels

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/validation"
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
	name, colour, description, actor, commandArgv, err := prepareLabelAdd(options)
	if err != nil {
		return AddResult{}, err
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return AddResult{}, err
	}
	request := map[string]any{"name": options.Name}
	if options.Colour != nil {
		request["colour"] = *options.Colour
	}
	if options.Description != nil {
		request["description"] = *options.Description
	}
	digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": actor, "operation": "label.add", "project_id": projectConfig.ProjectID, "request": request}, options.Actor, map[string]string{"name": "name", "colour": "--colour", "description": "--description"})
	if err != nil {
		return AddResult{}, err
	}
	var result AddResult
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command: strings.Join(commandArgv, " "),
	}, func(lock *journal.WriteLock) error {
		return addLabelLocked(lock, options, actor, name, colour, description, digest, fields, &result)
	})
	if err != nil {
		return AddResult{}, err
	}

	return result, nil
}

func prepareLabelAdd(options AddOptions) (string, *string, *string, standaloneActor, []string, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return "", nil, nil, standaloneActor{}, nil, fmt.Errorf("label add: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return "", nil, nil, standaloneActor{}, nil, fmt.Errorf("label add: migrations dir is required")
	}

	if strings.TrimSpace(options.Name) == "" {
		return "", nil, nil, standaloneActor{}, nil, fmt.Errorf("label add: label name is required")
	}
	name, err := validation.NormalizeText("label name", options.Name, 255)
	if err != nil {
		return "", nil, nil, standaloneActor{}, nil, err
	}
	actor, err := standaloneActorForOptions(options.Actor)
	if err != nil {
		return "", nil, nil, standaloneActor{}, nil, err
	}

	colour, err := normalizeOptionalLabelText("label colour", options.Colour, 64)
	if err != nil {
		return "", nil, nil, standaloneActor{}, nil, err
	}
	description, err := normalizeOptionalLabelText("label description", options.Description, 4096)
	if err != nil {
		return "", nil, nil, standaloneActor{}, nil, err
	}
	return name, colour, description, actor, labelAddArgv(name, colour, description, options.Actor, options.IdempotencyKey), nil
}

func addLabelLocked(lock *journal.WriteLock, options AddOptions, actor standaloneActor, name string, colour, description *string, digest string, fields map[string]string, result *AddResult) error {
	if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
		return err
	}

	conn, err := openLabelsProjectDB(options.WorkspaceRoot)
	if err != nil {
		return err
	}
	defer closeLabelsConn(conn)

	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return err
	}
	idempotencyRequest := journal.IdempotencyRequest{OperationScope: "label.add:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
	if outcome, found, err := journal.ReplayIdempotent(conn, idempotencyRequest); err != nil {
		return err
	} else if found {
		decoded, err := labelResultFromJSON(outcome.ResultJSON, outcome.EventID, true)
		if err != nil {
			return err
		}
		*result = decoded
		return nil
	}
	if err := requireMissingLabel(conn, name); err != nil {
		return err
	}
	sequence, err := journal.NextAggregateSequence(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.JournalDirName), journal.Aggregate{ID: name, Type: "label"})
	if err != nil {
		return err
	}
	event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "label.added", ActorV2: actor.eventActor(), Aggregate: journal.Aggregate{ID: name, Sequence: sequence, Type: "label"}, CommandName: "label.add", Effects: []journal.Effect{{Op: "insert", Table: "labels", PK: map[string]any{"name": name}, Values: labelAddEffectValues(colour, description)}}}
	seed := AddResult{Label: Label{Name: name, Colour: colour, Description: description}, IdempotencyKey: options.IdempotencyKey}
	outcome, err := journal.CommitIdempotent(lock, conn, idempotencyRequest, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
		seed.OccurredAt = prepared.OccurredAt.UTC().Format(time.RFC3339Nano)
		encoded, err := labelResultJSON(seed)
		if err != nil {
			return journal.IdempotencyMutation{}, err
		}
		return journal.IdempotencyMutation{EntityID: name, ResultJSON: encoded}, nil
	})
	if err != nil {
		return err
	}
	decoded, err := labelResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
	if err != nil {
		return err
	}
	*result = decoded
	return nil
}

func openLabelsProjectDB(workspaceRoot string) (*sql.DB, error) {
	return fiddb.Open(filepath.Join(config.LedgerDir(workspaceRoot), config.DatabaseFileName))
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

func labelAddArgv(name string, colour, description *string, actor, key string) []string {
	argv := []string{"fid", "label", "add", name}
	if colour != nil {
		argv = append(argv, "--colour", *colour)
	}
	if description != nil {
		argv = append(argv, "--description", *description)
	}
	return append(argv, "--actor", actor, "--idempotency-key", key)
}

func normalizeOptionalLabelText(field string, value *string, maximumBytes int) (*string, error) {
	if value == nil {
		return nil, nil
	}
	normalized, err := validation.NormalizeText(field, *value, maximumBytes)
	if err != nil {
		return nil, err
	}
	return &normalized, nil
}

func closeLabelsConn(conn *sql.DB) {
	_ = conn.Close()
}

func labelResultJSON(result AddResult) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"event_occurred_at": result.OccurredAt, "idempotency_key": result.IdempotencyKey, "label": result.Label, "replayed": false})
	return string(payload), err
}

func labelResultFromJSON(raw, eventID string, replayed bool) (AddResult, error) {
	var value struct {
		OccurredAt     string `json:"event_occurred_at"`
		IdempotencyKey string `json:"idempotency_key"`
		Label          Label  `json:"label"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return AddResult{}, fmt.Errorf("decode stored label add result: %w", err)
	}
	return AddResult{Label: value.Label, EventID: eventID, OccurredAt: value.OccurredAt, IdempotencyKey: value.IdempotencyKey, Replayed: replayed}, nil
}
