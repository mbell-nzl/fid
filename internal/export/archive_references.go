package export

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

func (v *semanticValidator) validateReferences(event journal.Event) error {
	for _, effect := range event.Effects {
		for column, value := range effect.Values {
			if value == nil || !strings.HasSuffix(column, "event_id") {
				continue
			}
			id, ok := value.(string)
			_, prior := v.eventIDs[id]
			if !ok || (id != event.EventID && !prior) {
				return fmt.Errorf("%s.%s references absent Event %q", effect.Table, column, id)
			}
		}
		if effect.Table == "run_recoveries" && effect.Values["prior_run_id"] != event.Aggregate.ID {
			return fmt.Errorf("recovery predecessor differs from Event aggregate")
		}
		if event.EventType == "run.started" && effect.Table == "runs" && effect.PK["id"] != event.Aggregate.ID {
			return fmt.Errorf("run identity differs from Event aggregate")
		}
		if event.Aggregate.Type == "run" && event.EventType != "run.recovered" && effect.Table == "run_runtime" && effect.PK["run_id"] != event.Aggregate.ID {
			return fmt.Errorf("runtime identity differs from Run aggregate")
		}
	}
	return nil
}

func (v *semanticValidator) validateProjectedReferences(event journal.Event) error {
	for _, effect := range event.Effects {
		if effect.Table == "stops" || effect.Table == "attention_items" {
			kind, hasSubject := effect.Values["subject_type"].(string)
			if !hasSubject {
				continue
			}
			table := map[string]string{"run": "runs", "envelope": "work_envelopes", "activation": "kit_activations"}[kind]
			if table == "" {
				return fmt.Errorf("unsupported polymorphic subject %q", kind)
			}
			var count int
			if err := v.conn.QueryRow("SELECT COUNT(*) FROM "+table+" WHERE id = ?", effect.Values["subject_id"]).Scan(&count); err != nil {
				return err
			}
			if count != 1 {
				return fmt.Errorf("%s subject does not resolve", effect.Table)
			}
		}
		if effect.Table == "run_recoveries" && effect.Op == "insert" {
			var count int
			err := v.conn.QueryRow(`SELECT COUNT(*) FROM run_recoveries r
JOIN runs p ON p.id=r.prior_run_id JOIN runs s ON s.id=r.successor_run_id
JOIN run_observations o ON o.id=r.observation_id
JOIN claims c ON c.id=r.successor_claim_id
WHERE r.id=? AND p.claim_id=r.prior_claim_id AND s.claim_id=r.successor_claim_id
AND o.run_id=p.id AND s.issue_id=p.issue_id AND s.envelope_id=p.envelope_id
AND s.envelope_revision=p.envelope_revision AND c.issue_id=s.issue_id
AND c.envelope_id=s.envelope_id AND c.envelope_revision=s.envelope_revision`, effect.PK["id"]).Scan(&count)
			if err != nil {
				return err
			}
			if count != 1 {
				return fmt.Errorf("recovery ownership references disagree")
			}
		}
	}
	return nil
}
