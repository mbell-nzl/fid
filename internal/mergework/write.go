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
	lock                                        *journal.WriteLock
	conn                                        *sql.DB
	WorkspaceRoot, MigrationsDir, TaskID, Actor string
	Clock                                       func() time.Time
}

type writeContext struct {
	conn       *sql.DB
	at         string
	generateID func(string, journal.IDExistsFunc) (string, error)
	reserved   map[string]struct{}
}

// withWrite holds the ledger lock only for refreshing, validating and appending
// the observation. No network or local check process runs here.
func withWrite(options Options, build func(writeContext) (journal.Event, error)) error {
	if !policy.Allows(options.Actor, policy.WorkSet) {
		return fmt.Errorf("merge/CI recording requires a work actor")
	}
	if options.lock != nil {
		return commit(options, options.lock, options.conn, build)
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
		return commit(options, lock, conn, build)
	})
}

func commit(options Options, lock *journal.WriteLock, conn *sql.DB, build func(writeContext) (journal.Event, error)) error {
	now := time.Now()
	if options.Clock != nil {
		now = options.Clock()
	}
	now = now.UTC().Truncate(time.Second)
	event, err := build(writeContext{conn: conn, at: now.Format(time.RFC3339)})
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
}

func (ctx *writeContext) id(prefix, table string) (string, error) {
	generate := ctx.generateID
	if generate == nil {
		generate = journal.GenerateID
	}
	id, err := generate(prefix, func(id string) (bool, error) {
		if _, taken := ctx.reserved[id]; taken {
			return true, nil
		}
		var count int
		err := ctx.conn.QueryRow("SELECT count(*) FROM "+table+" WHERE id=?", id).Scan(&count)
		return count > 0, err
	})
	if err != nil {
		return "", err
	}
	// CI allocates every ID before the batch is projected into SQLite.
	if ctx.reserved == nil {
		ctx.reserved = make(map[string]struct{})
	}
	ctx.reserved[id] = struct{}{}
	return id, nil
}
