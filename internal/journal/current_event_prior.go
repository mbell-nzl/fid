package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// Capture only producer-side missing prior values under the existing write
// lock. ProjectEvent never manufactures a prior for received journal data.
func captureCurrentEventPrior(conn *sql.DB, event *Event) error {
	if !event.IsV2() || event.EventType != "artifact.availability-observed" {
		return nil
	}
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	for index := range event.Effects {
		effect := &event.Effects[index]
		if err := validateCurrentEffect(*event, *effect); err != nil {
			return err
		}
		if effect.Op != "update" || effect.Prior != nil {
			continue
		}
		prior, found, err := readCurrentEventValues(tx, effect.Table, effect.PK, sortedMapKeys(effect.Values))
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("capture current Event prior: missing %s row", effect.Table)
		}
		for column, value := range prior {
			if currentEventTables[effect.Table].columns[column] == "b" {
				prior[column] = currentValuesEqual(value, 1)
			}
		}
		effect.Prior = prior
	}
	return nil
}

func validateCurrentPrior(tx *sql.Tx, event Event, effect Effect) error {
	if effect.Op != "update" {
		return nil
	}
	if len(effect.Prior) != len(effect.Values) {
		return fmt.Errorf("current Event %s: update requires exact captured prior", event.EventType)
	}
	for key := range effect.Values {
		if _, found := effect.Prior[key]; !found {
			return fmt.Errorf("current Event %s: missing prior %s", event.EventType, key)
		}
	}
	actual, found, err := readCurrentEventValues(tx, effect.Table, effect.PK, sortedMapKeys(effect.Prior))
	if err != nil {
		return err
	}
	for key, expected := range effect.Prior {
		if !found || !currentValuesEqual(expected, actual[key]) {
			return &PreconditionConflictError{EventID: event.EventID, EventType: event.EventType, Table: effect.Table, PK: effect.PK, Expected: effect.Prior, Actual: actual}
		}
	}
	return nil
}

// SQLite's driver decodes DATETIME as time.Time. CAST preserves the exact
// stored text (including subsecond precision) in the Event's captured prior.
func readCurrentEventValues(tx *sql.Tx, table string, pk map[string]any, columns []string) (map[string]any, bool, error) {
	definition, found := currentEventTables[table]
	if !found {
		return nil, false, fmt.Errorf("unknown current table")
	}
	selected := make([]string, 0, len(columns))
	for _, column := range columns {
		kind, found := definition.columns[column]
		if !found {
			return nil, false, fmt.Errorf("unknown current column")
		}
		name := `"` + column + `"`
		if strings.HasPrefix(kind, "s") {
			name = "CAST(" + name + " AS TEXT)"
		}
		selected = append(selected, name)
	}
	where, args, err := buildWhereClause(pk)
	if err != nil {
		return nil, false, err
	}
	values, pointers := make([]any, len(columns)), make([]any, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}
	err = tx.QueryRow("SELECT "+strings.Join(selected, ",")+" FROM "+table+" WHERE "+where, args...).Scan(pointers...)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	result := map[string]any{}
	for i, column := range columns {
		result[column] = values[i]
	}
	return result, true, nil
}

func currentValuesEqual(left, right any) bool {
	a, err := json.Marshal(normalizeValue(currentEventValue(left)))
	if err != nil {
		return false
	}
	b, err := json.Marshal(normalizeValue(currentEventValue(right)))
	return err == nil && string(a) == string(b)
}
