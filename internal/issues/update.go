package issues

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

type UpdateOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	ID            string
	ActorRole     string
	Title         *string
	Priority      *int
	Labels        *[]string
	SpecMDPath    *string
	EstTokens     *int
	DueDate       *string
	OwnerRole     *string
	OwnerActor    *string
}

type UpdateResult struct {
	ID      string
	Changes []FieldChange
}

type FieldChange struct {
	Field string
	From  any
	To    any
}

type issueUpdateState struct {
	ID         string
	Title      string
	Priority   *int
	Labels     []string
	SpecMDPath *string
	EstTokens  *int
	DueDate    *string
	OwnerRole  *string
	OwnerActor *string
	UpdatedAt  string
}

func Update(options UpdateOptions) (UpdateResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return UpdateResult{}, fmt.Errorf("update issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return UpdateResult{}, fmt.Errorf("update issue: migrations dir is required")
	}
	if strings.TrimSpace(options.ID) == "" {
		return UpdateResult{}, fmt.Errorf("update issue: issue id is required")
	}
	if err := policy.RequireKnownRole(options.WorkspaceRoot, options.ActorRole); err != nil {
		return UpdateResult{}, err
	}

	request, err := normalizeUpdateRequest(options)
	if err != nil {
		return UpdateResult{}, err
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return UpdateResult{}, err
	}

	commandArgv := issueUpdateArgv(options.ID, options.ActorRole, request)
	var result UpdateResult
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		dbPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName)
		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		current, err := loadIssueUpdateState(conn, options.ID)
		if err != nil {
			return err
		}
		if request.labels != nil {
			if err := requireExistingLabels(conn, *request.labels); err != nil {
				return fmt.Errorf("%s", strings.Replace(err.Error(), "create issue:", "update issue:", 1))
			}
		}

		scalarValues, scalarPreconditions, scalarPrior, changes := buildIssueUpdateChanges(current, request)
		labelEffects, labelChanged, labelChange := buildLabelEffects(current.ID, current.Labels, request.labels)
		if labelChanged {
			changes = append(changes, labelChange)
		}
		if len(changes) == 0 {
			return fmt.Errorf("update issue: at least one field must change")
		}

		activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
			return activityIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}
		occurredAt := nextIssueUpdateTime(current.UpdatedAt, time.Now().UTC())

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "issue.updated",
			OccurredAt:  occurredAt,
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
		scalarValues["updated_at"] = timestamp
		scalarPreconditions["updated_at"] = current.UpdatedAt
		scalarPrior["updated_at"] = current.UpdatedAt

		effects := []journal.Effect{
			{
				Op:            "update",
				Table:         "issues",
				PK:            map[string]any{"id": current.ID},
				Values:        scalarValues,
				Prior:         scalarPrior,
				Preconditions: scalarPreconditions,
			},
		}
		effects = append(effects, labelEffects...)
		effects = append(effects, journal.Effect{
			Op:    "activity",
			Table: "activity_log",
			PK:    map[string]any{"id": activityID},
			Values: map[string]any{
				"journal_event_id": event.EventID,
				"ref_type":         "issue",
				"ref_id":           current.ID,
				"event":            "issue_updated",
				"details": map[string]any{
					"changes": activityDetailChanges(changes),
				},
				"actor":       options.ActorRole,
				"occurred_at": timestamp,
			},
		})
		event.Effects = effects

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

		result = UpdateResult{
			ID:      current.ID,
			Changes: changes,
		}
		return nil
	})
	if err != nil {
		return UpdateResult{}, err
	}

	return result, nil
}

type normalizedUpdateRequest struct {
	title      *string
	priority   *int
	labels     *[]string
	specMDPath *string
	estTokens  *int
	dueDate    *string
	ownerRole  *string
	ownerActor *string
}

func normalizeUpdateRequest(options UpdateOptions) (normalizedUpdateRequest, error) {
	request := normalizedUpdateRequest{
		priority:  options.Priority,
		estTokens: options.EstTokens,
	}

	if options.Title != nil {
		title := strings.TrimSpace(*options.Title)
		if title == "" {
			return normalizedUpdateRequest{}, fmt.Errorf("update issue: title must not be blank")
		}
		request.title = &title
	}
	if options.Priority != nil && (*options.Priority < 0 || *options.Priority > 5) {
		return normalizedUpdateRequest{}, fmt.Errorf("update issue: priority must be between 0 and 5")
	}
	if options.Labels != nil {
		labels, err := normalizedLabels(*options.Labels)
		if err != nil {
			return normalizedUpdateRequest{}, fmt.Errorf("%s", strings.Replace(err.Error(), "create issue:", "update issue:", 1))
		}
		sort.Strings(labels)
		request.labels = &labels
	}
	if options.SpecMDPath != nil {
		specMDPath := strings.TrimSpace(*options.SpecMDPath)
		if specMDPath == "" {
			return normalizedUpdateRequest{}, fmt.Errorf("update issue: spec_md_path must not be blank")
		}
		request.specMDPath = &specMDPath
	}
	if options.DueDate != nil {
		dueDate := strings.TrimSpace(*options.DueDate)
		if dueDate == "" {
			return normalizedUpdateRequest{}, fmt.Errorf("update issue: due_date must use YYYY-MM-DD")
		}
		if _, err := time.Parse("2006-01-02", dueDate); err != nil {
			return normalizedUpdateRequest{}, fmt.Errorf("update issue: due_date must use YYYY-MM-DD")
		}
		request.dueDate = &dueDate
	}
	if options.OwnerRole != nil {
		ownerRole := strings.TrimSpace(*options.OwnerRole)
		if ownerRole == "" {
			return normalizedUpdateRequest{}, fmt.Errorf("update issue: owner_role must not be blank")
		}
		if err := policy.RequireKnownRole(options.WorkspaceRoot, ownerRole); err != nil {
			return normalizedUpdateRequest{}, err
		}
		request.ownerRole = &ownerRole
	}
	if options.OwnerActor != nil {
		ownerActor := strings.TrimSpace(*options.OwnerActor)
		if ownerActor == "" {
			return normalizedUpdateRequest{}, fmt.Errorf("update issue: owner_actor must not be blank")
		}
		request.ownerActor = &ownerActor
	}

	return request, nil
}

func buildIssueUpdateChanges(current issueUpdateState, request normalizedUpdateRequest) (map[string]any, map[string]any, map[string]any, []FieldChange) {
	values := make(map[string]any)
	preconditions := make(map[string]any)
	prior := make(map[string]any)
	changes := make([]FieldChange, 0, 8)

	if request.title != nil && current.Title != *request.title {
		values["title"] = *request.title
		preconditions["title"] = current.Title
		prior["title"] = current.Title
		changes = append(changes, FieldChange{Field: "title", From: current.Title, To: *request.title})
	}
	if request.priority != nil && !sameOptionalInt(current.Priority, request.priority) {
		values["priority"] = *request.priority
		preconditions["priority"] = intPointerValue(current.Priority)
		prior["priority"] = intPointerValue(current.Priority)
		changes = append(changes, FieldChange{Field: "priority", From: intPointerValue(current.Priority), To: *request.priority})
	}
	if request.specMDPath != nil && !sameOptionalString(current.SpecMDPath, request.specMDPath) {
		values["spec_md_path"] = *request.specMDPath
		preconditions["spec_md_path"] = stringPointerValue(current.SpecMDPath)
		prior["spec_md_path"] = stringPointerValue(current.SpecMDPath)
		changes = append(changes, FieldChange{Field: "spec_md_path", From: stringPointerValue(current.SpecMDPath), To: *request.specMDPath})
	}
	if request.estTokens != nil && !sameOptionalInt(current.EstTokens, request.estTokens) {
		values["est_tokens"] = *request.estTokens
		preconditions["est_tokens"] = intPointerValue(current.EstTokens)
		prior["est_tokens"] = intPointerValue(current.EstTokens)
		changes = append(changes, FieldChange{Field: "est_tokens", From: intPointerValue(current.EstTokens), To: *request.estTokens})
	}
	if request.dueDate != nil && !sameOptionalString(current.DueDate, request.dueDate) {
		values["due_date"] = *request.dueDate
		preconditions["due_date"] = stringPointerValue(current.DueDate)
		prior["due_date"] = stringPointerValue(current.DueDate)
		changes = append(changes, FieldChange{Field: "due_date", From: stringPointerValue(current.DueDate), To: *request.dueDate})
	}
	if request.ownerRole != nil && !sameOptionalString(current.OwnerRole, request.ownerRole) {
		values["owner_role"] = *request.ownerRole
		preconditions["owner_role"] = stringPointerValue(current.OwnerRole)
		prior["owner_role"] = stringPointerValue(current.OwnerRole)
		changes = append(changes, FieldChange{Field: "owner_role", From: stringPointerValue(current.OwnerRole), To: *request.ownerRole})
	}
	if request.ownerActor != nil && !sameOptionalString(current.OwnerActor, request.ownerActor) {
		values["owner_actor"] = *request.ownerActor
		preconditions["owner_actor"] = stringPointerValue(current.OwnerActor)
		prior["owner_actor"] = stringPointerValue(current.OwnerActor)
		changes = append(changes, FieldChange{Field: "owner_actor", From: stringPointerValue(current.OwnerActor), To: *request.ownerActor})
	}

	return values, preconditions, prior, changes
}

func nextIssueUpdateTime(currentUpdatedAt string, now time.Time) time.Time {
	parsed, err := time.Parse(time.RFC3339, currentUpdatedAt)
	if err != nil {
		return now.UTC()
	}

	candidate := now.UTC()
	if candidate.Format(time.RFC3339) <= parsed.UTC().Format(time.RFC3339) {
		return parsed.UTC().Add(time.Second)
	}
	return candidate
}

func buildLabelEffects(issueID string, current []string, desired *[]string) ([]journal.Effect, bool, FieldChange) {
	if desired == nil {
		return nil, false, FieldChange{}
	}
	if sameStringSlice(current, *desired) {
		return nil, false, FieldChange{}
	}

	currentSet := make(map[string]bool, len(current))
	for _, label := range current {
		currentSet[label] = true
	}
	desiredSet := make(map[string]bool, len(*desired))
	for _, label := range *desired {
		desiredSet[label] = true
	}

	effects := make([]journal.Effect, 0, len(current)+len(*desired))
	for _, label := range current {
		if desiredSet[label] {
			continue
		}
		effects = append(effects, journal.Effect{
			Op:    "delete",
			Table: "issue_labels",
			PK: map[string]any{
				"issue_id": issueID,
				"label":    label,
			},
		})
	}
	for _, label := range *desired {
		if currentSet[label] {
			continue
		}
		effects = append(effects, journal.Effect{
			Op:    "insert",
			Table: "issue_labels",
			PK: map[string]any{
				"issue_id": issueID,
				"label":    label,
			},
			Values: map[string]any{},
		})
	}

	return effects, true, FieldChange{
		Field: "labels",
		From:  append([]string(nil), current...),
		To:    append([]string(nil), (*desired)...),
	}
}

func activityDetailChanges(changes []FieldChange) []map[string]any {
	details := make([]map[string]any, 0, len(changes))
	for _, change := range changes {
		details = append(details, map[string]any{
			"field": change.Field,
			"from":  change.From,
			"to":    change.To,
		})
	}
	return details
}

func loadIssueUpdateState(conn *sql.DB, issueID string) (issueUpdateState, error) {
	row := conn.QueryRow(`
		SELECT id, title, priority, spec_md_path, est_tokens, due_date, owner_role, owner_actor, updated_at
		FROM issues
		WHERE id = ?
	`, issueID)

	var (
		state                                      issueUpdateState
		priority, estTokens                        sql.NullInt64
		specMDPath, dueDate, ownerRole, ownerActor sql.NullString
		updatedAt                                  sql.NullString
	)
	if err := row.Scan(
		&state.ID,
		&state.Title,
		&priority,
		&specMDPath,
		&estTokens,
		&dueDate,
		&ownerRole,
		&ownerActor,
		&updatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return issueUpdateState{}, &NotFoundError{IssueID: issueID}
		}
		return issueUpdateState{}, fmt.Errorf("query issue %q: %w", issueID, err)
	}
	state.Priority = nullIntPointer(priority)
	state.SpecMDPath = nullStringPointer(specMDPath)
	state.EstTokens = nullIntPointer(estTokens)
	state.DueDate = normalizeDateOnlyPointer(nullStringPointer(dueDate))
	state.OwnerRole = nullStringPointer(ownerRole)
	state.OwnerActor = nullStringPointer(ownerActor)
	if !updatedAt.Valid || strings.TrimSpace(updatedAt.String) == "" {
		return issueUpdateState{}, fmt.Errorf("query issue %q: updated_at is required", issueID)
	}
	state.UpdatedAt = updatedAt.String

	labels, err := issueLabels(conn, issueID)
	if err != nil {
		return issueUpdateState{}, err
	}
	state.Labels = labels
	return state, nil
}

func issueUpdateArgv(issueID, actorRole string, request normalizedUpdateRequest) []string {
	argv := []string{"fid", "issue", "update", issueID}
	if request.title != nil {
		argv = append(argv, "--title", *request.title)
	}
	if request.priority != nil {
		argv = append(argv, "--priority", fmt.Sprintf("%d", *request.priority))
	}
	if request.labels != nil {
		if len(*request.labels) == 0 {
			argv = append(argv, "--labels", "")
		}
		for _, label := range *request.labels {
			argv = append(argv, "--labels", label)
		}
	}
	if request.ownerRole != nil {
		argv = append(argv, "--owner-role", *request.ownerRole)
	}
	if request.ownerActor != nil {
		argv = append(argv, "--owner-actor", *request.ownerActor)
	}
	if request.specMDPath != nil {
		argv = append(argv, "--spec-md", *request.specMDPath)
	}
	if request.estTokens != nil {
		argv = append(argv, "--est-tokens", fmt.Sprintf("%d", *request.estTokens))
	}
	if request.dueDate != nil {
		argv = append(argv, "--due-date", *request.dueDate)
	}
	argv = append(argv, "--by", actorRole)
	return argv
}

func activityIDExists(conn *sql.DB, activityID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM activity_log WHERE id = ? LIMIT 1`, activityID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("query activity id %q: %w", activityID, err)
}

func stringPointerValue(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func intPointerValue(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}

func sameOptionalString(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func sameOptionalInt(left, right *int) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func sameStringSlice(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
