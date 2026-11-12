package export

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

func (v *semanticValidator) validateReferences(event journal.Event) error {
	if err := validateAggregateEffectIdentity(event); err != nil {
		return err
	}
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
		if event.EventType == "hold.taken" && effect.Table == "holds" && effect.Op == "insert" && effect.PK["id"] != event.Aggregate.ID {
			return fmt.Errorf("hold identity differs from Event aggregate")
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

var aggregateIdentityTables = map[string]map[string]struct{}{
	"activation":        {"kit_activations": {}},
	"attempt":           {"execution_attempts": {}},
	"claim":             {"claims": {}},
	"envelope":          {"work_envelopes": {}},
	"evidence":          {"evidence_records": {}},
	"gate-registration": {"gate_registrations": {}},
	"issue":             {"issues": {}},
	"label":             {"labels": {}},
	"link":              {"links": {}},
	"project":           {"projects": {}},
	"receipt":           {"execution_receipts": {}},
	"stop":              {"stops": {}},
	"attention":         {"attention_items": {}},
	"workflow":          {"workflows": {}},
}

func validateAggregateEffectIdentity(event journal.Event) error {
	tables := aggregateIdentityTables[event.Aggregate.Type]
	for _, effect := range event.Effects {
		if _, relevant := tables[effect.Table]; !relevant {
			continue
		}
		id, hasID := effect.PK["id"].(string)
		if hasID && id != event.Aggregate.ID {
			return fmt.Errorf("%s identity differs from Event aggregate", effect.Table)
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
	if err := v.validateProjectedObservationJSON(event.EventID); err != nil {
		return err
	}
	return nil
}

func (v *semanticValidator) validateProjectedObservationJSON(eventID string) error {
	rows, err := v.conn.Query(`SELECT id, run_id, kind, actor_id, observed_at, observation_json FROM run_observations WHERE created_event_id = ?`, eventID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var row struct {
			ID         string
			RunID      string
			Kind       string
			ActorID    string
			ObservedAt string
			Raw        string
		}
		if err := rows.Scan(&row.ID, &row.RunID, &row.Kind, &row.ActorID, &row.ObservedAt, &row.Raw); err != nil {
			return err
		}
		if err := rejectDuplicateJSONKeys([]byte(row.Raw)); err != nil {
			return fmt.Errorf("observation %q JSON is invalid: %w", row.ID, err)
		}
		var value struct {
			ID         string `json:"id"`
			RunID      string `json:"run_id"`
			Kind       string `json:"kind"`
			ActorID    string `json:"actor_id"`
			ObservedAt string `json:"observed_at"`
		}
		if err := json.Unmarshal([]byte(row.Raw), &value); err != nil {
			return fmt.Errorf("decode observation %q: %w", row.ID, err)
		}
		if value.ID != row.ID || value.RunID != row.RunID || value.Kind != row.Kind || value.ActorID != row.ActorID || value.ObservedAt != row.ObservedAt {
			return fmt.Errorf("observation %q JSON disagrees with immutable row identity", row.ID)
		}
	}
	return rows.Err()
}
