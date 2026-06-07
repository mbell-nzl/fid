package db

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	ProjectionSchemaVersion = 2

	ReplayStatusClean      = "clean"
	ReplayStatusCatchingUp = "catching_up"
	ReplayStatusRebuilding = "rebuilding"
	ReplayStatusError      = "error"
)

type ProjectionMeta struct {
	JournalFile             string
	JournalOffset           int64
	LastEventID             string
	LastEventHash           string
	EventSchema             int
	ProjectionSchemaVersion int
	ReplayStatus            string
	UpdatedAt               time.Time
}

type ProjectionFreshness string

const (
	ProjectionFresh   ProjectionFreshness = "clean"
	ProjectionCatchUp ProjectionFreshness = "catch_up"
	ProjectionRebuild ProjectionFreshness = "rebuild"
)

type ProjectionFreshnessInput struct {
	MetaExists                    bool
	ReplayStatus                  string
	StoredProjectionSchemaVersion int
	CurrentProjectionSchema       int
	CheckpointExists              bool
	CheckpointMatches             bool
	JournalHasNewerEvents         bool
}

func ReadProjectionMeta(conn *sql.DB) (ProjectionMeta, bool, error) {
	if conn == nil {
		return ProjectionMeta{}, false, fmt.Errorf("read projection meta: connection is nil")
	}

	var meta ProjectionMeta
	var journalFile sql.NullString
	var lastEventID sql.NullString
	var lastEventHash sql.NullString

	err := conn.QueryRow(`
		SELECT journal_file, journal_offset, last_event_id, last_event_hash,
		       event_schema, projection_schema_version, replay_status, updated_at
		FROM projection_meta
		WHERE key = 'default'
	`).Scan(
		&journalFile,
		&meta.JournalOffset,
		&lastEventID,
		&lastEventHash,
		&meta.EventSchema,
		&meta.ProjectionSchemaVersion,
		&meta.ReplayStatus,
		&meta.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ProjectionMeta{}, false, nil
		}
		return ProjectionMeta{}, false, fmt.Errorf("read projection meta row: %w", err)
	}

	if journalFile.Valid {
		meta.JournalFile = journalFile.String
	}
	if lastEventID.Valid {
		meta.LastEventID = lastEventID.String
	}
	if lastEventHash.Valid {
		meta.LastEventHash = lastEventHash.String
	}

	return meta, true, nil
}

func WriteProjectionMeta(tx *sql.Tx, meta ProjectionMeta) error {
	if tx == nil {
		return fmt.Errorf("write projection meta: transaction is nil")
	}

	if meta.ProjectionSchemaVersion == 0 {
		meta.ProjectionSchemaVersion = ProjectionSchemaVersion
	}
	if meta.ReplayStatus == "" {
		meta.ReplayStatus = ReplayStatusClean
	}
	if meta.UpdatedAt.IsZero() {
		return fmt.Errorf("write projection meta: updated_at is required")
	}

	_, err := tx.Exec(`
		INSERT INTO projection_meta (
			key, journal_file, journal_offset, last_event_id, last_event_hash,
			event_schema, projection_schema_version, replay_status, updated_at
		)
		VALUES ('default', ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			journal_file = excluded.journal_file,
			journal_offset = excluded.journal_offset,
			last_event_id = excluded.last_event_id,
			last_event_hash = excluded.last_event_hash,
			event_schema = excluded.event_schema,
			projection_schema_version = excluded.projection_schema_version,
			replay_status = excluded.replay_status,
			updated_at = excluded.updated_at
	`,
		nullIfEmpty(meta.JournalFile),
		meta.JournalOffset,
		nullIfEmpty(meta.LastEventID),
		nullIfEmpty(meta.LastEventHash),
		meta.EventSchema,
		meta.ProjectionSchemaVersion,
		meta.ReplayStatus,
		meta.UpdatedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("write projection meta row: %w", err)
	}

	return nil
}

func DecideProjectionFreshness(input ProjectionFreshnessInput) ProjectionFreshness {
	if !input.MetaExists {
		return ProjectionRebuild
	}
	if input.ReplayStatus != "" && input.ReplayStatus != ReplayStatusClean {
		return ProjectionRebuild
	}
	if input.StoredProjectionSchemaVersion != input.CurrentProjectionSchema {
		return ProjectionRebuild
	}
	if !input.CheckpointExists || !input.CheckpointMatches {
		return ProjectionRebuild
	}
	if input.JournalHasNewerEvents {
		return ProjectionCatchUp
	}
	return ProjectionFresh
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}
