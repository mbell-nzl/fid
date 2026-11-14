package mergework

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

type Options struct {
	WorkspaceRoot, MigrationsDir, TaskID, Actor string
	Clock                                       func() time.Time
}

type writeContext struct {
	conn    *sql.DB
	options Options
	at      string
}

// withWrite holds the ledger lock only for refreshing, validating and appending
// the observation. No network or local check process runs here.
func withWrite(options Options, build func(writeContext) (journal.Event, error)) error {
	if !policy.Allows(options.Actor, policy.WorkSet) {
		return fmt.Errorf("merge/CI recording requires a work actor")
	}
	return journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid record merge/CI"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		now := time.Now()
		if options.Clock != nil {
			now = options.Clock()
		}
		now = now.UTC().Truncate(time.Second)
		event, err := build(writeContext{conn, options, now.Format(time.RFC3339)})
		if err != nil {
			return err
		}
		if event.EventType == "" {
			return nil
		}
		actorType, actorID, _ := strings.Cut(options.Actor, ":")
		event.EventSchema = journal.EventSchemaV2
		event.OccurredAt = now
		event.ActorV2 = journal.EventActor{ActorID: actorID, ActorType: actorType, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}
		event.Aggregate.Sequence, err = journal.NextAggregateSequence(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.JournalDirName), event.Aggregate)
		if err != nil {
			return err
		}
		_, err = journal.CommitEvent(lock, conn, event)
		return err
	})
}

func (ctx writeContext) id(prefix, table string) (string, error) {
	return journal.GenerateID(prefix, func(id string) (bool, error) {
		var count int
		err := ctx.conn.QueryRow("SELECT count(*) FROM "+table+" WHERE id=?", id).Scan(&count)
		return count > 0, err
	})
}
