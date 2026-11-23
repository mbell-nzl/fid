package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

// FindingResponseEffects is also used at admission: every dependent row is
// determined by the immutable response and existing reviewed task.
func FindingResponseEffects(conn reviewrecord.Queryer, r reviewrecord.Response) ([]Effect, error) {
	if err := r.Admit(conn); err != nil {
		return nil, err
	}
	effects := []Effect{}
	var taskID any
	if d := r.DeferredTask; d != nil {
		taskID = d.ID
		var parent any
		if d.EpicID != nil {
			parent = *d.EpicID
		}
		effects = append(effects,
			Effect{Op: "insert", Table: "issues", PK: map[string]any{"id": d.ID}, Values: map[string]any{"kind": "task", "title": d.Title, "parent_id": parent, "workflow_id": "task-default", "status": "todo", "revision": 1, "owner_role": "developer", "lane": d.Lane, "created_at": r.At, "updated_at": r.At}},
			Effect{Op: "insert", Table: "acceptance_criteria", PK: map[string]any{"id": r.ID + "-ac"}, Values: map[string]any{"issue_id": d.ID, "ordinal": 1, "text": r.Text}},
		)
		if d.EpicID != nil {
			var sequence sql.NullInt64
			if err := conn.QueryRow(`SELECT CASE WHEN scheduling='sequential' THEN (SELECT COALESCE(MAX(sequence),0)+1 FROM issues WHERE parent_id=?) END FROM epic_metadata WHERE issue_id=?`, *d.EpicID, *d.EpicID).Scan(&sequence); err != nil {
				return nil, fmt.Errorf("assign deferred sequence: %w", err)
			}
			if sequence.Valid {
				effects[0].Values["sequence"] = sequence.Int64
			}
		}
		for _, link := range []struct{ from, relation, to string }{{d.ID, "discovered_from", r.TaskID}, {r.TaskID, "discovered", d.ID}} {
			effects = append(effects, Effect{Op: "insert", Table: "links", PK: map[string]any{"from_type": "issue", "from_id": link.from, "relation": link.relation, "to_type": "issue", "to_id": link.to}, Values: map[string]any{"created_at": r.At}})
		}
		var revision int
		if err := conn.QueryRow(`SELECT revision FROM issues WHERE id=?`, r.TaskID).Scan(&revision); err != nil {
			return nil, err
		}
		effects = append(effects, Effect{Op: "update", Table: "issues", PK: map[string]any{"id": r.TaskID}, Values: map[string]any{"revision": revision + 1}, Prior: map[string]any{"revision": revision}, Preconditions: map[string]any{"revision": revision}})
	}
	raw, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return append(effects, Effect{Op: "insert", Table: "finding_responses", PK: map[string]any{"id": r.ID}, Values: map[string]any{"finding_id": r.FindingID, "deferred_task_id": taskID, "record_json": string(raw)}}), nil
}

func validateFindingResponse(conn *sql.Tx, e Event) error {
	verb := "answer"
	if e.EventType == "finding.deferred" {
		verb = "defer"
	}
	if e.CommandName != "finding."+verb || e.Aggregate.Sequence != 1 || e.Repository != nil || e.Idempotency != nil || len(e.Effects) == 0 {
		return fmt.Errorf("invalid finding response command or aggregate")
	}
	last := e.Effects[len(e.Effects)-1]
	if err := validateCurrentEffect(e, last); err != nil {
		return err
	}
	var r reviewrecord.Response
	if err := decodeCheckJSON(last.Values["record_json"], &r); err != nil {
		return err
	}
	if r.ID != e.Aggregate.ID || r.Actor != e.ActorV2.ActorType+":"+e.ActorV2.ActorID || r.At != e.OccurredAt.UTC().Format(time.RFC3339) || (r.DeferredTask != nil) != (verb == "defer") {
		return fmt.Errorf("finding response differs from Event")
	}
	want, err := FindingResponseEffects(conn, r)
	if err != nil {
		return err
	}
	if !currentValuesEqual(e.Effects, want) {
		return fmt.Errorf("finding response effects differ from canonical record")
	}
	return nil
}
