package export

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	migratesqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	migratefile "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mbell-nzl/fid/internal/journal"
	_ "modernc.org/sqlite"
)

// semanticValidator replays a bounded archive's already-integrity-checked
// Event lines against only the installed migration profile selected by its
// manifest. It never reads archive-supplied SQL or writes a project database.
type semanticValidator struct {
	conn         *sql.DB
	profile      int
	tables       map[string]semanticTable
	eventIDs     map[string]struct{}
	projectIDs   map[string]struct{}
	bootstrapID  string
	eventCount   int
	terminalHash string
	journalLine  int
}

type semanticTable struct {
	columns map[string]string
	primary map[string]struct{}
}

type semanticValidationError struct {
	Code string
	Err  error
}

func (err *semanticValidationError) Error() string { return err.Err.Error() }

func (err *semanticValidationError) Unwrap() error { return err.Err }

type semanticArtifact struct {
	ID           string
	Digest       string
	MediaType    string
	SizeBytes    int64
	Availability string
}

type archiveSemanticFacts struct {
	ProjectIDs        map[string]struct{}
	Sources           []json.RawMessage
	Artifacts         []semanticArtifact
	EventCount        int
	TerminalEventHash string
}

func newSemanticValidator(migrationsDir string, projectionSchema int) (*semanticValidator, error) {
	if projectionSchema != 14 && projectionSchema != 17 && projectionSchema != 18 && projectionSchema != 19 && projectionSchema != 20 && projectionSchema != 21 && projectionSchema != 22 && projectionSchema != 23 && projectionSchema != 24 {
		return nil, fmt.Errorf("semantic archive validation: unsupported projection schema %d", projectionSchema)
	}
	if strings.TrimSpace(migrationsDir) == "" {
		return nil, fmt.Errorf("semantic archive validation: migrations directory is required")
	}
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("semantic archive validation: open private projection: %w", err)
	}
	conn.SetMaxOpenConns(1)
	if _, err := conn.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("semantic archive validation: enable foreign keys: %w", err)
	}
	if _, err := conn.Exec(`PRAGMA recursive_triggers = ON`); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("semantic archive validation: enable recursive triggers: %w", err)
	}
	if err := migrateSemanticProjection(conn, migrationsDir, projectionSchema); err != nil {
		_ = conn.Close()
		return nil, err
	}
	tables, err := loadSemanticTables(conn)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &semanticValidator{
		conn:       conn,
		profile:    projectionSchema,
		tables:     tables,
		eventIDs:   make(map[string]struct{}),
		projectIDs: make(map[string]struct{}),
	}, nil
}

func migrateSemanticProjection(conn *sql.DB, migrationsDir string, projectionSchema int) error {
	absMigrations, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("semantic archive validation: resolve migrations: %w", err)
	}
	driver, err := migratesqlite.WithInstance(conn, &migratesqlite.Config{})
	if err != nil {
		return fmt.Errorf("semantic archive validation: create migration driver: %w", err)
	}
	source, err := (&migratefile.File{}).Open("file://" + absMigrations)
	if err != nil {
		return fmt.Errorf("semantic archive validation: open trusted migration source: %w", err)
	}
	defer func() { _ = source.Close() }()
	migrator, err := migrate.NewWithInstance("trusted-migrations", source, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("semantic archive validation: open trusted migrations: %w", err)
	}
	if err := migrator.Migrate(uint(projectionSchema)); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("semantic archive validation: apply trusted migrations: %w", err)
	}
	return nil
}

func loadSemanticTables(conn *sql.DB) (map[string]semanticTable, error) {
	rows, err := conn.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND name NOT IN ('schema_migrations', 'projection_meta') ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("semantic archive validation: list trusted tables: %w", err)
	}
	defer func() { _ = rows.Close() }()
	names := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("semantic archive validation: scan trusted table: %w", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("semantic archive validation: iterate trusted tables: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("semantic archive validation: close trusted table list: %w", err)
	}
	tables := make(map[string]semanticTable, len(names))
	for _, name := range names {
		columns, primary, err := loadSemanticColumns(conn, name)
		if err != nil {
			return nil, err
		}
		tables[name] = semanticTable{columns: columns, primary: primary}
	}
	return tables, nil
}

func loadSemanticColumns(conn *sql.DB, table string) (map[string]string, map[string]struct{}, error) {
	rows, err := conn.Query(`SELECT name, type, pk FROM pragma_table_info(?)`, table)
	if err != nil {
		return nil, nil, fmt.Errorf("semantic archive validation: list columns for %q: %w", table, err)
	}
	defer func() { _ = rows.Close() }()
	columns := make(map[string]string)
	primary := make(map[string]struct{})
	for rows.Next() {
		var name, kind string
		var primaryKey int
		if err := rows.Scan(&name, &kind, &primaryKey); err != nil {
			return nil, nil, fmt.Errorf("semantic archive validation: scan column for %q: %w", table, err)
		}
		columns[name] = strings.ToUpper(kind)
		if primaryKey > 0 {
			primary[name] = struct{}{}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("semantic archive validation: iterate columns for %q: %w", table, err)
	}
	return columns, primary, nil
}

func (validator *semanticValidator) close() error {
	if validator == nil || validator.conn == nil {
		return nil
	}
	err := validator.conn.Close()
	validator.conn = nil
	return err
}

func (validator *semanticValidator) accept(raw []byte) error {
	if validator == nil || validator.conn == nil {
		return fmt.Errorf("semantic archive validation: validator is closed")
	}
	if err := rejectDuplicateJSONKeys(raw); err != nil {
		return semanticFailure("effect-invalid", "Event has duplicate JSON key", err)
	}
	var event journal.Event
	if err := event.UnmarshalJSON(raw); err != nil {
		return semanticFailure("effect-invalid", "decode Event", err)
	}
	if _, exists := validator.eventIDs[event.EventID]; exists {
		return semanticFailure("effect-invalid", fmt.Sprintf("duplicate Event ID %q", event.EventID), nil)
	}
	if err := validator.validateEvent(event); err != nil {
		return semanticFailure("effect-invalid", "validate Event family", err)
	}
	if event.EventType == "project.initialized" {
		if validator.bootstrapID != "" {
			return semanticFailure("effect-invalid", "validate project bootstrap", fmt.Errorf("event %q repeats project initialization", event.EventID))
		}
		validator.bootstrapID = event.Aggregate.ID
		validator.projectIDs[event.Aggregate.ID] = struct{}{}
	}
	if event.Repository != nil && event.Repository.ProjectID != validator.bootstrapID {
		return semanticFailure("reference-invalid", "validate repository project identity", fmt.Errorf("event %q repository project %q does not match initialized project %q", event.EventID, event.Repository.ProjectID, validator.bootstrapID))
	}
	if err := validator.validateEffects(event); err != nil {
		return semanticFailure("effect-invalid", "validate Event Effects", err)
	}
	if err := validator.validateReferences(event); err != nil {
		return semanticFailure("reference-invalid", "validate Event references", err)
	}
	hash, err := journal.HashEvent(event)
	if err != nil {
		return fmt.Errorf("semantic archive validation: hash Event %q: %w", event.EventID, err)
	}
	validator.journalLine++
	if err := journal.ProjectEvent(validator.conn, journal.ProjectionInput{
		Event: event, EventHash: hash, File: "archive", Line: validator.journalLine, Offset: int64(validator.journalLine),
	}); err != nil {
		return semanticFailure("reference-invalid", fmt.Sprintf("apply Event %q", event.EventID), err)
	}
	if err := validator.validateProjectedReferences(event); err != nil {
		return semanticFailure("reference-invalid", "validate projected references", err)
	}
	validator.eventIDs[event.EventID] = struct{}{}
	if event.Repository != nil && event.Repository.ProjectID != "" {
		validator.projectIDs[event.Repository.ProjectID] = struct{}{}
	}
	validator.eventCount++
	validator.terminalHash = hash
	return nil
}

func semanticFailure(code, operation string, err error) error {
	if err == nil {
		err = errors.New(operation)
	} else {
		err = fmt.Errorf("%s: %w", operation, err)
	}
	return &semanticValidationError{Code: code, Err: fmt.Errorf("semantic archive validation: %w", err)}
}

func (validator *semanticValidator) finish() (archiveSemanticFacts, error) {
	if validator == nil || validator.conn == nil {
		return archiveSemanticFacts{}, fmt.Errorf("semantic archive validation: validator is closed")
	}
	sources, err := semanticSources(validator.conn)
	if err != nil {
		return archiveSemanticFacts{}, err
	}
	artifacts, err := semanticArtifacts(validator.conn)
	if err != nil {
		return archiveSemanticFacts{}, err
	}
	return archiveSemanticFacts{
		ProjectIDs:        cloneSemanticIDs(validator.projectIDs),
		Sources:           sources,
		Artifacts:         artifacts,
		EventCount:        validator.eventCount,
		TerminalEventHash: validator.terminalHash,
	}, nil
}

func cloneSemanticIDs(ids map[string]struct{}) map[string]struct{} {
	cloned := make(map[string]struct{}, len(ids))
	for id := range ids {
		cloned[id] = struct{}{}
	}
	return cloned
}

func semanticSources(conn *sql.DB) ([]json.RawMessage, error) {
	rows, err := conn.Query(`SELECT source_json FROM envelope_sources ORDER BY envelope_id, revision, ordinal`)
	if err != nil {
		return nil, fmt.Errorf("semantic archive validation: read Sources: %w", err)
	}
	defer func() { _ = rows.Close() }()
	sources := []json.RawMessage{}
	for rows.Next() {
		var source string
		if err := rows.Scan(&source); err != nil {
			return nil, fmt.Errorf("semantic archive validation: scan Source: %w", err)
		}
		sources = append(sources, json.RawMessage(source))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("semantic archive validation: iterate Sources: %w", err)
	}
	return sources, nil
}

func semanticArtifacts(conn *sql.DB) ([]semanticArtifact, error) {
	rows, err := conn.Query(`SELECT id, digest, media_type, size_bytes, availability FROM referenced_artifacts ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("semantic archive validation: read Artifacts: %w", err)
	}
	defer func() { _ = rows.Close() }()
	artifacts := []semanticArtifact{}
	for rows.Next() {
		var artifact semanticArtifact
		if err := rows.Scan(&artifact.ID, &artifact.Digest, &artifact.MediaType, &artifact.SizeBytes, &artifact.Availability); err != nil {
			return nil, fmt.Errorf("semantic archive validation: scan Artifact: %w", err)
		}
		artifacts = append(artifacts, artifact)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("semantic archive validation: iterate Artifacts: %w", err)
	}
	return artifacts, nil
}

func (validator *semanticValidator) validateEvent(event journal.Event) error {
	if event.EventSchema != journal.EventSchemaV2 && (validator.profile != 14 || event.EventSchema != journal.EventSchemaV1) {
		return fmt.Errorf("semantic archive validation: Event %q has unsupported schema %d for projection schema %d", event.EventID, event.EventSchema, validator.profile)
	}
	rule, found := semanticEventRules[event.EventType]
	if !found || (validator.profile == 14 && semanticPost14Events[event.EventType]) || validator.profile < semanticFirstProfile[event.EventType] {
		return fmt.Errorf("semantic archive validation: Event %q has unsupported type %q", event.EventID, event.EventType)
	}
	if event.Aggregate.Type != rule.aggregate {
		return fmt.Errorf("semantic archive validation: Event %q aggregate %q does not match type %q", event.EventID, event.Aggregate.Type, event.EventType)
	}
	return nil
}

func (validator *semanticValidator) validateEffects(event journal.Event) error {
	rule := semanticEventRules[event.EventType]
	if len(event.Effects) == 0 {
		if event.EventType == "project.initialized" {
			return nil
		}
		return fmt.Errorf("semantic archive validation: Event %q has no Effects", event.EventID)
	}
	stageIndex := 0
	stageCounts := make([]int, len(rule.stages))
	for index, effect := range event.Effects {
		if err := validateSemanticUpdate(effect); err != nil {
			return err
		}
		for stageIndex < len(rule.stages) && !rule.stages[stageIndex].accepts(effect.Table, effect.Op) {
			if stageCounts[stageIndex] < rule.stages[stageIndex].minimum {
				return fmt.Errorf("event %q Effect %d %s/%s is out of order; expected stage %d", event.EventID, index+1, effect.Table, effect.Op, stageIndex+1)
			}
			stageIndex++
		}
		if stageIndex == len(rule.stages) {
			return fmt.Errorf("event %q Effect %d table/operation %s/%s is not accepted for %q", event.EventID, index+1, effect.Table, effect.Op, event.EventType)
		}
		stageCounts[stageIndex]++
		if maximum := rule.stages[stageIndex].maximum; maximum >= 0 && stageCounts[stageIndex] > maximum {
			return fmt.Errorf("event %q Effect %d exceeds accepted repetitions for %s/%s", event.EventID, index+1, effect.Table, effect.Op)
		}
		table, found := validator.tables[effect.Table]
		if !found {
			return fmt.Errorf("semantic archive validation: Event %q Effect %d uses unavailable table %q", event.EventID, index+1, effect.Table)
		}
		if err := validateSemanticMap(effect.Table, table, effect.PK, true); err != nil {
			return fmt.Errorf("semantic archive validation: Event %q Effect %d primary key: %w", event.EventID, index+1, err)
		}
		for name, values := range map[string]map[string]any{"values": effect.Values, "prior": effect.Prior, "preconditions": effect.Preconditions} {
			if err := validateSemanticMap(effect.Table, table, values, false); err != nil {
				return fmt.Errorf("semantic archive validation: Event %q Effect %d %s: %w", event.EventID, index+1, name, err)
			}
		}
		if created, found := effect.Values["created_event_id"]; found && created != event.EventID {
			return fmt.Errorf("event %q Effect %d created_event_id does not match Event identity", event.EventID, index+1)
		}
	}
	for index := stageIndex; index < len(rule.stages); index++ {
		if stageCounts[index] < rule.stages[index].minimum {
			return fmt.Errorf("event %q omits required Effect stage %d", event.EventID, index+1)
		}
	}
	return nil
}
