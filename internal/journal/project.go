package journal

import (
	"database/sql"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type ProjectionInput struct {
	Event      Event
	EventHash  string
	File       string
	Line       int
	Offset     int64
	ReplayMode string
}

type PreconditionConflictError struct {
	EventID   string
	EventType string
	Table     string
	PK        map[string]any
	Expected  map[string]any
	Actual    map[string]any
	File      string
	Line      int
}

func (err *PreconditionConflictError) Error() string {
	if err == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"replay precondition conflict at %s:%d for event %s (%s) on %s pk=%v expected=%v actual=%v",
		err.File,
		err.Line,
		err.EventID,
		err.EventType,
		err.Table,
		err.PK,
		err.Expected,
		err.Actual,
	)
}

func DryRunProjectEvent(conn *sql.DB, event Event) error {
	if err := validateProjectionConnection(conn); err != nil {
		return err
	}

	input := ProjectionInput{
		Event: event,
	}
	if input.Event.EventID == "" {
		input.Event.EventID = "EVT-dry-run"
	}

	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("project event %q: begin transaction: %w", input.Event.EventID, err)
	}

	if err := applyEventEffects(tx, input); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("project event %q: rollback transaction: %w", input.Event.EventID, err)
	}

	return nil
}

func ProjectEvent(conn *sql.DB, input ProjectionInput) error {
	if err := validateProjectionConnection(conn); err != nil {
		return err
	}
	if err := validateProjectionInput(input); err != nil {
		return err
	}

	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("project event %q: begin transaction: %w", input.Event.EventID, err)
	}

	if err := applyEventEffects(tx, input); err != nil {
		_ = tx.Rollback()
		return err
	}

	replayStatus := input.ReplayMode
	if replayStatus == "" {
		replayStatus = fiddb.ReplayStatusClean
	}
	if err := fiddb.WriteProjectionMeta(tx, fiddb.ProjectionMeta{
		JournalFile:             input.File,
		JournalOffset:           input.Offset,
		LastEventID:             input.Event.EventID,
		LastEventHash:           input.EventHash,
		EventSchema:             input.Event.EventSchema,
		ProjectionSchemaVersion: fiddb.ProjectionSchemaVersion,
		ReplayStatus:            replayStatus,
		UpdatedAt:               input.Event.OccurredAt.UTC(),
	}); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("project event %q: %w", input.Event.EventID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("project event %q: commit transaction: %w", input.Event.EventID, err)
	}

	return nil
}

func validateProjectionConnection(conn *sql.DB) error {
	if conn == nil {
		return fmt.Errorf("project event: connection is nil")
	}
	return nil
}

func validateProjectionInput(input ProjectionInput) error {
	if input.Event.EventID == "" {
		return fmt.Errorf("project event: event_id is required")
	}
	if input.EventHash == "" {
		return fmt.Errorf("project event %q: event hash is required", input.Event.EventID)
	}
	return nil
}

func applyEventEffects(tx *sql.Tx, input ProjectionInput) error {
	for _, effect := range input.Event.Effects {
		if err := applyEffect(tx, input, effect); err != nil {
			return err
		}
	}
	return nil
}

func applyEffect(tx *sql.Tx, input ProjectionInput, effect Effect) error {
	table, err := validatedIdentifier(effect.Table)
	if err != nil {
		return fmt.Errorf("project event %q: %w", input.Event.EventID, err)
	}

	switch effect.Op {
	case "insert":
		return applyInsert(tx, input.Event.EventID, table, effect)
	case "update":
		return applyUpdate(tx, input, table, effect)
	case "delete":
		return applyDelete(tx, input, table, effect)
	default:
		return fmt.Errorf("project event %q: unsupported effect op %q", input.Event.EventID, effect.Op)
	}
}

func applyInsert(tx *sql.Tx, eventID, table string, effect Effect) error {
	row, err := mergeEffectValues(effect.PK, effect.Values)
	if err != nil {
		return fmt.Errorf("project event %q: insert %s: %w", eventID, table, err)
	}

	columns := sortedMapKeys(row)
	placeholders := make([]string, 0, len(columns))
	args := make([]any, 0, len(columns))
	for _, column := range columns {
		validatedColumn, err := validatedIdentifier(column)
		if err != nil {
			return fmt.Errorf("project event %q: insert %s: %w", eventID, table, err)
		}
		columns = replaceValue(columns, column, validatedColumn)
		placeholders = append(placeholders, "?")
		args = append(args, row[column])
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)
	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("project event %q: insert into %s: %w", eventID, table, err)
	}

	return nil
}

func applyUpdate(tx *sql.Tx, input ProjectionInput, table string, effect Effect) error {
	if err := enforcePreconditions(tx, input, table, effect); err != nil {
		return err
	}

	setColumns := sortedMapKeys(effect.Values)
	assignments := make([]string, 0, len(setColumns))
	args := make([]any, 0, len(setColumns)+len(effect.PK))
	for _, column := range setColumns {
		validatedColumn, err := validatedIdentifier(column)
		if err != nil {
			return fmt.Errorf("project event %q: update %s: %w", input.Event.EventID, table, err)
		}
		assignments = append(assignments, validatedColumn+" = ?")
		args = append(args, effect.Values[column])
	}

	whereClause, whereArgs, err := buildWhereClause(effect.PK)
	if err != nil {
		return fmt.Errorf("project event %q: update %s: %w", input.Event.EventID, table, err)
	}
	args = append(args, whereArgs...)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(assignments, ", "), whereClause)
	result, err := tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("project event %q: update %s: %w", input.Event.EventID, table, err)
	}
	return ensureSingleRowAffected(result, input.Event.EventID, "update", table)
}

func applyDelete(tx *sql.Tx, input ProjectionInput, table string, effect Effect) error {
	if err := enforcePreconditions(tx, input, table, effect); err != nil {
		return err
	}

	whereClause, whereArgs, err := buildWhereClause(effect.PK)
	if err != nil {
		return fmt.Errorf("project event %q: delete %s: %w", input.Event.EventID, table, err)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, whereClause)
	result, err := tx.Exec(query, whereArgs...)
	if err != nil {
		return fmt.Errorf("project event %q: delete %s: %w", input.Event.EventID, table, err)
	}
	return ensureSingleRowAffected(result, input.Event.EventID, "delete", table)
}

func enforcePreconditions(tx *sql.Tx, input ProjectionInput, table string, effect Effect) error {
	if len(effect.Preconditions) == 0 {
		return nil
	}

	actual, found, err := readCurrentValues(tx, table, effect.PK, sortedMapKeys(effect.Preconditions))
	if err != nil {
		return fmt.Errorf("project event %q: read preconditions for %s: %w", input.Event.EventID, table, err)
	}
	if !found {
		actual = make(map[string]any, len(effect.Preconditions))
		for key := range effect.Preconditions {
			actual[key] = nil
		}
	}

	for column, expected := range effect.Preconditions {
		if !valuesEqual(expected, actual[column]) {
			return &PreconditionConflictError{
				EventID:   input.Event.EventID,
				EventType: input.Event.EventType,
				Table:     table,
				PK:        cloneMap(effect.PK),
				Expected:  cloneMap(effect.Preconditions),
				Actual:    cloneMap(actual),
				File:      input.File,
				Line:      input.Line,
			}
		}
	}

	return nil
}

func readCurrentValues(tx *sql.Tx, table string, pk map[string]any, columns []string) (map[string]any, bool, error) {
	selectColumns := make([]string, 0, len(columns))
	for _, column := range columns {
		validatedColumn, err := validatedIdentifier(column)
		if err != nil {
			return nil, false, err
		}
		selectColumns = append(selectColumns, validatedColumn)
	}

	whereClause, whereArgs, err := buildWhereClause(pk)
	if err != nil {
		return nil, false, err
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT 1", strings.Join(selectColumns, ", "), table, whereClause)
	row := tx.QueryRow(query, whereArgs...)

	destinations := make([]any, len(selectColumns))
	values := make([]any, len(selectColumns))
	for index := range destinations {
		destinations[index] = &values[index]
	}

	if err := row.Scan(destinations...); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}

	actual := make(map[string]any, len(selectColumns))
	for index, column := range columns {
		actual[column] = values[index]
	}

	return actual, true, nil
}

func buildWhereClause(pk map[string]any) (string, []any, error) {
	if len(pk) == 0 {
		return "", nil, fmt.Errorf("primary key is required")
	}

	columns := sortedMapKeys(pk)
	clauses := make([]string, 0, len(columns))
	args := make([]any, 0, len(columns))
	for _, column := range columns {
		validatedColumn, err := validatedIdentifier(column)
		if err != nil {
			return "", nil, err
		}
		clauses = append(clauses, validatedColumn+" = ?")
		args = append(args, pk[column])
	}

	return strings.Join(clauses, " AND "), args, nil
}

func mergeEffectValues(pk, values map[string]any) (map[string]any, error) {
	merged := cloneMap(pk)
	for key, value := range values {
		if existing, exists := merged[key]; exists && !valuesEqual(existing, value) {
			return nil, fmt.Errorf("conflicting values for key %q", key)
		}
		merged[key] = value
	}
	return merged, nil
}

func ensureSingleRowAffected(result sql.Result, eventID, op, table string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("project event %q: %s %s rows affected: %w", eventID, op, table, err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("project event %q: %s %s affected %d rows, want 1", eventID, op, table, rowsAffected)
	}
	return nil
}

func validatedIdentifier(identifier string) (string, error) {
	if !identifierPattern.MatchString(identifier) {
		return "", fmt.Errorf("invalid SQL identifier %q", identifier)
	}
	return identifier, nil
}

func sortedMapKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func cloneMap(values map[string]any) map[string]any {
	if len(values) == 0 {
		return map[string]any{}
	}

	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}

func valuesEqual(expected, actual any) bool {
	expected = normalizeValue(expected)
	actual = normalizeValue(actual)

	expectedNumber, expectedIsNumber := asFloat(expected)
	actualNumber, actualIsNumber := asFloat(actual)
	if expectedIsNumber && actualIsNumber {
		return math.Abs(expectedNumber-actualNumber) < 1e-9
	}

	return fmt.Sprintf("%v", expected) == fmt.Sprintf("%v", actual)
}

func normalizeValue(value any) any {
	switch typed := value.(type) {
	case []byte:
		return string(typed)
	case bool:
		if typed {
			return int64(1)
		}
		return int64(0)
	default:
		return value
	}
}

func asFloat(value any) (float64, bool) {
	switch typed := value.(type) {
	case int:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case float32:
		return float64(typed), true
	case float64:
		return typed, true
	default:
		return 0, false
	}
}

func replaceValue(values []string, oldValue, newValue string) []string {
	replaced := make([]string, len(values))
	for index, value := range values {
		if value == oldValue {
			replaced[index] = newValue
			continue
		}
		replaced[index] = value
	}
	return replaced
}
