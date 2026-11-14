package checkrun

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/journal"
)

type store struct {
	options Options
	lock    *journal.WriteLock
	conn    *sql.DB
	now     time.Time
}

func (options Options) write(action func(store) error) error {
	return journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid check"}, func(lock *journal.WriteLock) error {
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
		return action(store{options, lock, conn, now.UTC()})
	})
}

func (db store) event(family, id string, effects []journal.Effect) (journal.Event, error) {
	kind, actor, _ := strings.Cut(db.options.Actor, ":")
	sequence, err := journal.NextAggregateSequence(filepath.Join(config.LedgerDir(db.options.WorkspaceRoot), "journal"), journal.Aggregate{ID: id, Type: "check"})
	if err != nil {
		return journal.Event{}, err
	}
	return journal.Event{EventSchema: journal.EventSchemaV2, EventType: family, OccurredAt: db.now, ActorV2: journal.EventActor{ActorID: actor, ActorType: kind, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: id, Type: "check", Sequence: sequence}, CommandName: "check.execute", Effects: effects}, nil
}

func (db store) current() error {
	_, err := holdwork.Current(db.conn, holdwork.Options{WorkspaceRoot: db.options.WorkspaceRoot, HoldID: db.options.HoldID, Actor: db.options.Actor}, db.now)
	return err
}

func (db store) request() (journal.IdempotencyRequest, error) {
	cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(db.options.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return journal.IdempotencyRequest{}, err
	}
	digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": db.options.Actor, "operation": "check.execute", "project_id": cfg.ProjectID, "request": map[string]any{"hold": db.options.HoldID, "names": db.options.Names}}, db.options.Actor, map[string]string{"hold": "--hold", "names": "names"})
	return journal.IdempotencyRequest{OperationScope: "check.execute:" + cfg.ProjectID, IdempotencyKey: db.options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}, err
}

func encoded(value any) string { raw, _ := json.Marshal(value); return string(raw) }
