package journal

import (
	"database/sql"
	"fmt"
)

// All qualifying Phase 22 Stop facts have a durable original idempotency
// result (observations, recoveries, terminal Runs, or activation lifecycle).
// The Attempt-to-Receipt replacement is not a qualifying Stop fact.
func validateCurrentStopReferences(tx *sql.Tx, event Event) error {
	for _, effect := range event.Effects {
		if effect.Table == "stops" {
			table := map[any]string{"run": "runs", "envelope": "work_envelopes", "activation": "kit_activations"}[effect.Values["subject_type"]]
			if table == "" {
				return fmt.Errorf("current Stop: unsupported subject")
			}
			id := effect.Values["subject_id"]
			inserted := false
			for _, candidate := range event.Effects {
				if candidate.Op == "insert" && candidate.Table == table && candidate.PK["id"] == id {
					inserted = true
				}
			}
			if !inserted {
				var count int
				if err := tx.QueryRow("SELECT COUNT(*) FROM "+table+" WHERE id=?", id).Scan(&count); err != nil {
					return err
				}
				if count != 1 {
					return fmt.Errorf("current Stop: missing subject reference")
				}
			}
		}
		if effect.Table == "stop_resolutions" && event.EventType == "stop.resolved" {
			var count int
			if err := tx.QueryRow("SELECT COUNT(*) FROM idempotency_results WHERE result_event_id=?", effect.Values["fact_event_id"]).Scan(&count); err != nil {
				return err
			}
			if count == 0 {
				return fmt.Errorf("current Stop: missing qualifying Event reference")
			}
		}
	}
	return nil
}
