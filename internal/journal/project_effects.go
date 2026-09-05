package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

func applyEventEffects(tx *sql.Tx, input ProjectionInput) error {
	if err := validateProjectedEvent(input.Event); err != nil {
		return err
	}
	for _, effect := range input.Event.Effects {
		if err := applyEffect(tx, input, effect); err != nil {
			return err
		}
	}

	return nil
}

func validateProjectedEvent(event Event) error {
	if event.EventSchema != EventSchemaV2 {
		return nil
	}
	if _, unavailable := unavailableRecoveryEvents[event.EventType]; unavailable {
		return fmt.Errorf("Event-v2 %q does not permit effect application before its recovery domain is available", event.EventType)
	}
	return nil
}

var unavailableRecoveryEvents = map[string]struct{}{
	"run.observed":                  {},
	"run.checkpointed":              {},
	"run.recovered":                 {},
	"run.budget-observed":           {},
	"run.execution-blocked":         {},
	"run.closed":                    {},
	"run.abandoned":                 {},
	"activation.resource-observed":  {},
	"activation.recovery-selected":  {},
	"activation.recovery-completed": {},
}

func applyEffect(tx *sql.Tx, input ProjectionInput, effect Effect) error {
	table, err := validatedIdentifier(effect.Table)
	if err != nil {
		return fmt.Errorf("project event %q: %w", input.Event.EventID, err)
	}

	if table == "activity_log" && effect.Op != "activity" {
		return fmt.Errorf("project event %q: activity_log writes require effect op %q", input.Event.EventID, "activity")
	}
	if table == "execution_receipts" && effect.Op != "insert" {
		return fmt.Errorf("project event %q: execution_receipts are immutable", input.Event.EventID)
	}

	switch effect.Op {
	case "insert":
		return applyInsert(tx, input.Event.EventID, table, effect)
	case "update":
		return applyUpdate(tx, input, table, effect)
	case "delete":
		return applyDelete(tx, input, table, effect)
	case "activity":
		return applyActivity(tx, input.Event.EventID, table, effect)
	default:
		return fmt.Errorf("project event %q: unsupported effect op %q", input.Event.EventID, effect.Op)
	}
}

func applyActivity(tx *sql.Tx, eventID, table string, effect Effect) error {
	if table != "activity_log" {
		return fmt.Errorf("project event %q: activity effect must target %q, got %q", eventID, "activity_log", table)
	}

	normalizedValues, err := normalizeActivityValues(effect.Values)
	if err != nil {
		return fmt.Errorf("project event %q: activity %s: %w", eventID, table, err)
	}

	return applyInsert(tx, eventID, table, Effect{
		PK:     cloneMap(effect.PK),
		Values: normalizedValues,
	})
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

func normalizeActivityValues(values map[string]any) (map[string]any, error) {
	normalized := cloneMap(values)

	details, ok := normalized["details"]
	if !ok {
		return normalized, nil
	}

	serialized, err := serializeActivityDetails(details)
	if err != nil {
		return nil, err
	}
	normalized["details"] = serialized

	return normalized, nil
}

func serializeActivityDetails(details any) (any, error) {
	switch typed := details.(type) {
	case nil:
		return nil, nil
	case string:
		return typed, nil
	case []byte:
		return string(typed), nil
	default:
		payload, err := json.Marshal(typed)
		if err != nil {
			return nil, fmt.Errorf("serialize details: %w", err)
		}
		return string(payload), nil
	}
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
