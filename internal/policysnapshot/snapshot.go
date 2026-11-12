// Package policysnapshot records policy changes under the caller's write lock.
// Policy resolution and actor authority live in policy; this package owns
// only the journal/projection boundary, shared by policy-consuming writers.
package policysnapshot

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

// Ensure writes one snapshot if the effective digest changed, including when
// no task can be taken or a live hold is returned. Both Events use the same
// caller-held lock; a crash after the snapshot leaves a complete policy fact.
func Ensure(lock *journal.WriteLock, conn *sql.DB, loaded policy.Loaded, now time.Time) error {
	var sequence int
	var previous string
	err := conn.QueryRow(`SELECT sequence, digest FROM policy_snapshots ORDER BY sequence DESC LIMIT 1`).Scan(&sequence, &previous)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("read latest policy snapshot: %w", err)
	}
	digest := loaded.Effective.Digest()
	if previous == digest {
		return nil
	}
	encoded, err := loaded.Effective.Canonical()
	if err != nil {
		return fmt.Errorf("encode effective policy: %w", err)
	}
	var commit any
	if loaded.Source.Commit != "" {
		commit = loaded.Source.Commit
	}
	sequence++
	_, err = journal.CommitEvent(lock, conn, journal.Event{EventSchema: journal.EventSchemaV2, EventType: "policy.snapshotted", OccurredAt: now,
		ActorV2:   journal.EventActor{ActorID: "fid", ActorType: "service", Producer: "fid-standalone", ProvenanceRef: "built-in-policy"},
		Aggregate: journal.Aggregate{ID: "policy", Sequence: sequence, Type: "policy"}, CommandName: "policy.snapshot",
		Effects: []journal.Effect{{Op: "insert", Table: "policy_snapshots", PK: map[string]any{"sequence": sequence}, Values: map[string]any{
			"source_ref": loaded.Source.Ref, "source_commit": commit, "digest": digest, "effective_json": string(encoded),
		}}}})
	if err != nil {
		return fmt.Errorf("record policy snapshot: %w", err)
	}
	return nil
}
