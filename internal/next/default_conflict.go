package next

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/mbell-nzl/fid/internal/defaultconflict"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func returnDefaultConflict(lock *journal.WriteLock, conn *sql.DB, options Options, result *Result, planned order) error {
	id, err := journal.GenerateID("DFC", func(id string) (bool, error) {
		var n int
		err := conn.QueryRow(`SELECT count(*) FROM default_conflicts WHERE id=?`, id).Scan(&n)
		return n != 0, err
	})
	if err != nil {
		return err
	}
	effects, err := journal.DefaultConflictEffects(conn, result.Task.ID, planned.own.ID, id, *planned.reviews[result.Task.ID].conflict, result.Now)
	if err != nil {
		return err
	}
	kind, actor, _ := strings.Cut(options.Actor, ":")
	_, err = journal.CommitEvent(lock, conn, journal.Event{EventSchema: journal.EventSchemaV2, EventType: "task.default-conflict", CommandName: "next.claim", OccurredAt: result.Now, ActorV2: journal.EventActor{ActorType: kind, ActorID: actor, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: id, Type: "default-conflict", Sequence: 1}, Effects: effects})
	if err == nil {
		result.Task.Status = "in_progress"
		result.Queue[result.Served].Status = "in_progress"
	}
	return err
}

// A process can stop between the atomic return and taking the linked hold.
func readDefaultConflict(conn reviewrecord.Queryer, hold string) (*defaultconflict.Evidence, error) {
	var raw string
	err := conn.QueryRow(`SELECT record_json FROM default_conflicts WHERE hold_id=?`, hold).Scan(&raw)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var record defaultconflict.Record
	if err := json.Unmarshal([]byte(raw), &record); err != nil {
		return nil, err
	}
	return &record.Evidence, nil
}
