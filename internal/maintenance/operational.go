package maintenance

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/mbell-nzl/fid/internal/journal"
)

func checkOperationalIntegrity(db *sql.DB, journalDir string) ([]Finding, error) {
	events, findings, err := inspectOperationalEvents(journalDir)
	if err != nil {
		return nil, err
	}
	bindings, err := eventBindingFindings(db, events)
	if err != nil {
		return nil, err
	}
	attempts, err := attemptReceiptFindings(db)
	if err != nil {
		return nil, err
	}
	artifacts, err := artifactFindings(db)
	if err != nil {
		return nil, err
	}
	recovery, err := recoveryIntegrityFindings(db)
	if err != nil {
		return nil, err
	}
	frozen, err := frozenStateFindings(db)
	if err != nil {
		return nil, err
	}
	residue, err := removedStateResidueFindings(db)
	if err != nil {
		return nil, err
	}
	findings = append(findings, bindings...)
	findings = append(findings, attempts...)
	findings = append(findings, artifacts...)
	findings = append(findings, recovery...)
	findings = append(findings, frozen...)
	return append(findings, residue...), nil
}

func inspectOperationalEvents(journalDir string) (map[string]journal.Event, []Finding, error) {
	files, err := journalFiles(journalDir)
	if err != nil {
		return nil, nil, err
	}
	events := make(map[string]journal.Event)
	sequences := make(map[string]int)
	findings := make([]Finding, 0)
	for _, file := range files {
		path := filepath.Join(journalDir, file)
		if err := readOperationalEvents(path, file, events, sequences, &findings); err != nil {
			return nil, nil, err
		}
	}
	return events, findings, nil
}

func journalFiles(journalDir string) ([]string, error) {
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		return nil, fmt.Errorf("inspect operational journal: read directory: %w", err)
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.Type().IsRegular() && filepath.Ext(entry.Name()) == ".jsonl" {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)
	return files, nil
}

func readOperationalEvents(path, file string, events map[string]journal.Event, sequences map[string]int, findings *[]Finding) error {
	input, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("inspect operational journal: open %q: %w", path, err)
	}
	defer func() { _ = input.Close() }()
	reader := bufio.NewReader(input)
	for line := 1; ; line++ {
		bytes, readErr := reader.ReadBytes('\n')
		if len(bytes) > 0 {
			inspectOperationalEvent(bytes, fmt.Sprintf("%s:%d", file, line), events, sequences, findings)
		}
		if readErr == nil {
			continue
		}
		if readErr == io.EOF {
			return nil
		}
		return fmt.Errorf("inspect operational journal: read %q: %w", path, readErr)
	}
}

func inspectOperationalEvent(data []byte, location string, events map[string]journal.Event, sequences map[string]int, findings *[]Finding) {
	var event journal.Event
	if err := json.Unmarshal(data, &event); err != nil {
		*findings = append(*findings, Finding{Check: "event-sequence", EntityID: location, Message: "journal Event is invalid"})
		return
	}
	if !event.IsV2() {
		return
	}
	if _, exists := events[event.EventID]; exists {
		*findings = append(*findings, Finding{Check: "event-sequence", EntityID: event.EventID, Message: "Event ID occurs more than once"})
	}
	events[event.EventID] = event
	key := event.Aggregate.Type + ":" + event.Aggregate.ID
	if prior, exists := sequences[key]; exists && event.Aggregate.Sequence != prior+1 {
		*findings = append(*findings, Finding{Check: "aggregate-sequence", EntityID: key, Message: fmt.Sprintf("aggregate sequence %d follows %d, want %d", event.Aggregate.Sequence, prior, prior+1)})
	}
	sequences[key] = event.Aggregate.Sequence
}

func eventBindingFindings(db *sql.DB, events map[string]journal.Event) ([]Finding, error) {
	rows, err := db.Query(`SELECT operation_scope, idempotency_key, result_event_id, aggregate_sequence FROM idempotency_results ORDER BY operation_scope, idempotency_key`)
	if err != nil {
		return nil, fmt.Errorf("check Event bindings: query idempotency results: %w", err)
	}
	defer func() { _ = rows.Close() }()
	findings := make([]Finding, 0)
	for rows.Next() {
		var scope, key, eventID string
		var sequence int
		if err := rows.Scan(&scope, &key, &eventID, &sequence); err != nil {
			return nil, fmt.Errorf("check Event bindings: scan result: %w", err)
		}
		event, exists := events[eventID]
		if !exists || event.Aggregate.Sequence != sequence {
			findings = append(findings, Finding{Check: "event-binding", EntityID: scope + ":" + key, Message: "idempotency result does not bind its journal Event"})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("check Event bindings: iterate results: %w", err)
	}
	return findings, nil
}

func attemptReceiptFindings(db *sql.DB) ([]Finding, error) {
	rows, err := db.Query(`SELECT attempt.id FROM execution_attempts AS attempt LEFT JOIN execution_receipts AS receipt ON receipt.attempt_id = attempt.id WHERE receipt.id IS NULL ORDER BY attempt.id`)
	if err != nil {
		return nil, fmt.Errorf("check Attempt/Receipt integrity: query attempts: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return simpleFindings(rows, "attempt-receipt-integrity", "Attempt has no terminal Receipt")
}

func artifactFindings(db *sql.DB) ([]Finding, error) {
	rows, err := db.Query(`SELECT id FROM referenced_artifacts WHERE availability = 'missing' ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("check Artifact integrity: query Artifacts: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return simpleFindings(rows, "artifact-integrity", "Artifact is unavailable")
}

func simpleFindings(rows *sql.Rows, check, message string) ([]Finding, error) {
	findings := make([]Finding, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("check %s: scan row: %w", check, err)
		}
		findings = append(findings, Finding{Check: check, EntityID: id, Message: message})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("check %s: iterate rows: %w", check, err)
	}
	return findings, nil
}

func frozenStateFindings(db *sql.DB) ([]Finding, error) {
	return tableCountFindings(db, []string{"agent_runs"}, "frozen-state-present", "frozen state is non-authoritative")
}

func removedStateResidueFindings(db *sql.DB) ([]Finding, error) {
	return tablePresenceFindings(db, []string{"gate_runs", "roles", "work_claims"})
}

func tableCountFindings(db *sql.DB, tables []string, check, message string) ([]Finding, error) {
	findings := make([]Finding, 0)
	for _, table := range tables {
		var count int
		if err := db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&count); err != nil {
			return nil, fmt.Errorf("check %s: count %s: %w", check, table, err)
		}
		if count > 0 {
			findings = append(findings, Finding{Check: check, EntityID: table, Message: fmt.Sprintf("%s has %d row(s); %s", table, count, message)})
		}
	}
	return findings, nil
}

func tablePresenceFindings(db *sql.DB, tables []string) ([]Finding, error) {
	findings := make([]Finding, 0)
	for _, table := range tables {
		var count int
		if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_schema WHERE type = 'table' AND name = ?`, table).Scan(&count); err != nil {
			return nil, fmt.Errorf("check removed state residue: inspect %s: %w", table, err)
		}
		if count > 0 {
			findings = append(findings, Finding{Check: "removed-state-residue", EntityID: table, Message: "removed table is present"})
		}
	}
	return findings, nil
}
