//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// p23AssertStopBlockedDetails compares the response IDs with the durable open
// Stop and its linked Attention, independently of the response under test.
func p23AssertStopBlockedDetails(t *testing.T, a *p22FIDRecoveryAdapter, response p22ConformanceResponse, code, detailField, detailValue string) {
	t.Helper()
	envelope := p23Envelope(t, response.Stderr)
	details := p21Path(t, envelope, "error", "details").(map[string]any)
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	rows, err := conn.Query(`SELECT s.id, l.attention_id FROM stops s
		JOIN attention_stop_links l ON l.stop_id=s.id
		LEFT JOIN stop_resolutions r ON r.stop_id=s.id
		WHERE s.code=? AND s.subject_id=? AND json_extract(s.details_json, '$.'||?)=? AND r.stop_id IS NULL
		ORDER BY s.id, l.attention_id`, code, p23EnvelopeID(t, a), detailField, detailValue)
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	stopIDs, attentionIDs := []any{}, []any{}
	for rows.Next() {
		var stopID, attentionID string
		p22FaultOK(t, rows.Scan(&stopID, &attentionID))
		if len(stopIDs) == 0 || stopIDs[len(stopIDs)-1] != stopID {
			stopIDs = append(stopIDs, stopID)
		}
		attentionIDs = append(attentionIDs, attentionID)
	}
	p22FaultOK(t, rows.Err())
	if len(stopIDs) != 1 || !reflect.DeepEqual(details["stop_ids"], stopIDs) || !reflect.DeepEqual(details["attention_ids"], attentionIDs) || details["event"] != nil || details["replayed"] != false {
		t.Fatalf("%s denial details = %#v, want durable Stops %v, Attention %v, event:null, replayed:false", code, details, stopIDs, attentionIDs)
	}
}

func p23AssertNoncommittingDenialRetry(t *testing.T, a *p22FIDRecoveryAdapter, first p22ConformanceResponse, retry func() p22ConformanceResponse) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	count := func(table string) int {
		t.Helper()
		var n int
		p22FaultOK(t, conn.QueryRow("SELECT COUNT(*) FROM "+table).Scan(&n))
		return n
	}
	beforeEvents := p21JournalEventCount(t, a.workspace)
	beforeKeys, beforeStops, beforeAttention := count("idempotency_results"), count("stops"), count("attention_items")
	firstEnvelope := p23Envelope(t, first.Stderr)
	code := p21PathString(t, firstEnvelope, "error", "code")
	if code == "protected-decision-required" {
		class := p21PathString(t, firstEnvelope, "error", "details", "decision_class")
		p23AssertStopBlockedDetails(t, a, first, code, "decision_class", class)
	}
	second := retry()
	if second != first || !p23HasJSON(second.Stderr, `"event":null`) || !p23HasJSON(second.Stderr, `"replayed":false`) {
		t.Fatalf("noncommitting denial retry changed its response: first=%#v second=%#v", first, second)
	}
	if p21JournalEventCount(t, a.workspace) != beforeEvents || count("idempotency_results") != beforeKeys || count("stops") != beforeStops || count("attention_items") != beforeAttention {
		t.Fatal("noncommitting denial retry changed Event, idempotency, Stop, or Attention state")
	}
	if code == "protected-decision-required" {
		class := p21PathString(t, firstEnvelope, "error", "details", "decision_class")
		p23AssertStopBlockedDetails(t, a, second, code, "decision_class", class)
	}
}

// TestP23AmendedEnvelopeReplayThroughBuiltBinary covers the two Envelope
// mutations amended by Phase 23 after the standard fixture has committed
// their original typed snapshots.
func TestP23AmendedEnvelopeReplayThroughBuiltBinary(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, p23ObligationGuardActions())
	cases := []struct{ command, file, key string }{
		{"create", "p22-corpus-envelope.json", "p22-corpus-envelope-create"},
		{"activate", "p22-corpus-envelope-activate.json", "p22-corpus-envelope-activate"},
	}
	for _, tc := range cases {
		t.Run(tc.command, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join(a.workspace, tc.file))
			p22FaultOK(t, err)
			var input map[string]any
			p22FaultOK(t, json.Unmarshal(raw, &input))
			before := p21JournalEventCount(t, a.workspace)
			replay := p23RunInput(t, a, "envelope", tc.command, a.actor, tc.key, input)
			p23AssertReviewSuccess(t, "Envelope "+tc.command+" replay", replay)
			if p21Path(t, p23Envelope(t, replay.Stdout), "result", "replayed") != true || p21JournalEventCount(t, a.workspace) != before {
				t.Fatalf("Envelope %s same-key replay changed the journal or lost its replay marker", tc.command)
			}
			changed := p23Clone(input)
			actor := a.actor
			if tc.command == "create" {
				changed["outcome"] = input["outcome"].(string) + " with changed digest"
			} else {
				actor = "actor-denied"
			}
			conflict := p23RunInput(t, a, "envelope", tc.command, actor, tc.key, changed)
			assertP22ProcessFailure(t, "envelope."+tc.command, "idempotency-conflict", 4, conflict)
			if p21JournalEventCount(t, a.workspace) != before {
				t.Fatalf("Envelope %s digest conflict appended an Event", tc.command)
			}
		})
	}
}
