//go:build integration

package integration

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// This model receives archive entry bytes, not a workspace path. It has no
// filesystem, database, fid domain, or journal-applier dependency. Each logical
// Effect changes an owned in-memory row; the resulting facts drive assertions.
type p22PortableModel struct {
	rows   map[string]map[string]map[string]any
	order  map[string][]string
	events []string
}

type p22PortableEffect struct {
	Op            string         `json:"op"`
	Table         string         `json:"table"`
	PK            map[string]any `json:"pk"`
	Values        map[string]any `json:"values"`
	Prior         map[string]any `json:"prior"`
	Preconditions map[string]any `json:"preconditions"`
}

func p22PortableFakeImport(entries map[string][]byte) (*p22PortableModel, error) {
	var manifest struct {
		Schema     string `json:"schema"`
		EventCount int    `json:"event_count"`
		Entries    []struct {
			Path   string `json:"path"`
			Digest string `json:"digest"`
			Size   int    `json:"size_bytes"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(entries["manifest.json"], &manifest); err != nil {
		return nil, err
	}
	if manifest.Schema != "fid.export-manifest.v2" {
		return nil, fmt.Errorf("unsupported archive schema")
	}
	model := &p22PortableModel{rows: map[string]map[string]map[string]any{}, order: map[string][]string{}}
	seen := map[string]bool{}
	for _, entry := range manifest.Entries {
		data, exists := entries[entry.Path]
		if !exists || len(data) != entry.Size || fmt.Sprintf("sha256:%x", sha256.Sum256(data)) != entry.Digest {
			return nil, fmt.Errorf("archive inventory mismatch")
		}
		if !strings.HasPrefix(entry.Path, "journal/") {
			continue
		}
		for _, line := range bytes.Split(bytes.TrimSuffix(data, []byte{'\n'}), []byte{'\n'}) {
			var event struct {
				ID      string              `json:"event_id"`
				Schema  int                 `json:"event_schema"`
				Type    string              `json:"event_type"`
				Effects []p22PortableEffect `json:"effects"`
			}
			if err := json.Unmarshal(line, &event); err != nil {
				return nil, err
			}
			if event.Schema != 2 || event.ID == "" || seen[event.ID] {
				return nil, fmt.Errorf("invalid or duplicate Event identity")
			}
			if !strings.Contains(" "+p22PortableEventTypes+" ", " "+event.Type+" ") {
				return nil, fmt.Errorf("unknown Event type %s", event.Type)
			}
			seen[event.ID] = true
			model.events = append(model.events, event.ID)
			for i, effect := range event.Effects {
				if err := model.apply(effect, fmt.Sprintf("%s:%d", event.ID, i)); err != nil {
					return nil, err
				}
			}
		}
	}
	if len(model.events) != manifest.EventCount {
		return nil, fmt.Errorf("Event count mismatch")
	}
	return model, nil
}

func (model *p22PortableModel) apply(effect p22PortableEffect, eventKey string) error {
	columns, known := p22PortableColumns[effect.Table]
	if !known {
		return fmt.Errorf("unknown Effect table %s", effect.Table)
	}
	for _, fields := range []map[string]any{effect.PK, effect.Values, effect.Prior, effect.Preconditions} {
		for field := range fields {
			if !strings.Contains(" "+columns+" ", " "+field+" ") {
				return fmt.Errorf("unknown Effect column %s.%s", effect.Table, field)
			}
		}
	}
	if model.rows[effect.Table] == nil {
		model.rows[effect.Table] = map[string]map[string]any{}
	}
	encoded, _ := json.Marshal(effect.PK)
	key := string(encoded)
	if effect.Op == "activity" {
		key = eventKey
	}
	row, exists := model.rows[effect.Table][key]
	for _, conditions := range []map[string]any{effect.Prior, effect.Preconditions} {
		for field, value := range conditions {
			if !exists || !reflect.DeepEqual(row[field], value) {
				return fmt.Errorf("Effect precondition differs: %s.%s got=%v want=%v", effect.Table, field, row[field], value)
			}
		}
	}
	switch effect.Op {
	case "insert", "activity":
		if exists {
			return fmt.Errorf("duplicate immutable Effect identity")
		}
		row = map[string]any{}
		for _, field := range strings.Fields(columns) {
			row[field] = nil
		}
		for field, value := range p22PortableDefaults[effect.Table] {
			data, _ := json.Marshal(value)
			_ = json.Unmarshal(data, &value)
			row[field] = value
		}
		for field, value := range effect.PK {
			row[field] = value
		}
		model.rows[effect.Table][key] = row
		model.order[effect.Table] = append(model.order[effect.Table], key)
	case "update":
		if !exists {
			return fmt.Errorf("updated Effect identity missing")
		}
	case "delete":
		if !exists {
			return fmt.Errorf("deleted Effect identity missing")
		}
		delete(model.rows[effect.Table], key)
		return nil
	default:
		return fmt.Errorf("unknown Effect operation %s", effect.Op)
	}
	for field, value := range effect.Values {
		row[field] = value
	}
	return nil
}

func (model *p22PortableModel) table(name string) []map[string]any {
	rows := []map[string]any{}
	for _, key := range model.order[name] {
		if row, exists := model.rows[name][key]; exists {
			rows = append(rows, row)
		}
	}
	return rows
}

func p22PortableFakeSummary(t *testing.T, entries map[string][]byte) map[string]any {
	t.Helper()
	model, err := p22PortableFakeImport(entries)
	if err != nil {
		t.Fatalf("independent portable Effect interpretation: %v", err)
	}
	// A second reconstruction owns new rows and independently reapplies every
	// Effect from the preserved canonical bytes; no copied expected view is used.
	rebuilt, err := p22PortableFakeImport(entries)
	if err != nil || !reflect.DeepEqual(model, rebuilt) {
		t.Fatalf("fake rebuild differs from imported logical state: %v", err)
	}
	statuses := []string{}
	for _, run := range model.table("runs") {
		statuses = append(statuses, run["status"].(string))
	}
	sort.Strings(statuses)
	summary := map[string]any{"run_statuses": statuses, "attention_count": len(model.table("attention_items")), "stop_count": len(model.table("stops")), "attention_link_count": len(model.table("attention_stop_links"))}
	activationStatuses := []string{}
	for _, activation := range model.table("kit_activations") {
		activationStatuses = append(activationStatuses, activation["status"].(string))
	}
	sort.Strings(activationStatuses)
	summary["activation_statuses"] = activationStatuses
	availability := []string{}
	for _, artifact := range model.table("referenced_artifacts") {
		availability = append(availability, artifact["availability"].(string))
	}
	sort.Strings(availability)
	summary["artifact_availability"] = availability
	count := 0
	for path := range entries {
		if strings.HasPrefix(path, "artifacts/") {
			count++
		}
	}
	summary["artifact_body_count"] = count
	known := false
	consumed := float64(0)
	maximum := float64(0)
	for _, policy := range model.table("envelope_budgets") {
		if policy["name"] == "tokens" {
			var budget struct {
				Maximum float64 `json:"maximum"`
			}
			if err := json.Unmarshal([]byte(policy["budget_json"].(string)), &budget); err != nil {
				t.Fatalf("fake budget policy: %v", err)
			}
			maximum = budget.Maximum
		}
	}
	for _, sample := range model.table("budget_samples") {
		if sample["budget_name"] != "tokens" {
			continue
		}
		if sample["value"] == nil {
			known = false
			continue
		}
		if sample["sample_kind"] == "baseline" {
			consumed = sample["value"].(float64)
			known = true
		} else if known {
			consumed += sample["value"].(float64)
		}
	}
	summary["tokens_state"], summary["tokens_consumed"], summary["tokens_remaining"] = "unknown", nil, nil
	if known {
		remaining := max(float64(0), maximum-consumed)
		state := "available"
		if remaining == 0 {
			state = "exhausted"
		}
		summary["tokens_state"], summary["tokens_consumed"], summary["tokens_remaining"] = state, consumed, remaining
	}
	data, _ := json.Marshal(summary)
	_ = json.Unmarshal(data, &summary)
	return summary
}

const p22PortableEventTypes = "project.initialized activation.prepared activation.component-staged activation.committed activation.finalized activation.aborted activation.resource-observed activation.recovery-selected activation.recovery-completed gate.registered envelope.created envelope.activated envelope.completed claim.acquired claim.renewed claim.released claim.expired run.started run.observed run.checkpointed run.recovered run.budget-observed run.execution-blocked run.closed run.abandoned gate.execution-started gate.receipt-recorded evidence.linked artifact.availability-observed stop.raised stop.resolved attention.resolved issue.created issue.updated issue.transitioned issue.ac-added issue.label-added issue.label-removed label.added link.added link.removed workflow.defined workflow.assigned requirement.synced"
