package journal

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/usagerecord"
)

// UsageEffect is the complete immutable session row shared by producer and
// admission. SQL exposes categories directly; record_json retains model detail.
func UsageEffect(r usagerecord.Record) (Effect, error) {
	raw, err := MarshalCanonicalJSON(r)
	if err != nil {
		return Effect{}, err
	}
	return Effect{Op: "insert", Table: "usage_sessions", PK: map[string]any{"id": r.ID}, Values: map[string]any{
		"task_id": r.TaskID, "hold_id": r.HoldID, "actor": r.Actor, "model": r.Model, "role": r.Role, "stream": r.Stream,
		"format": r.Format, "input_digest": r.InputDigest, "parser_revision": r.ParserRevision, "tool_version": r.ToolVersion, "version_source": r.VersionSource,
		"uncached_input": r.Counts.Input, "cache_read": r.Counts.CacheRead, "cache_write": r.Counts.CacheWrite, "output": r.Counts.Output, "reasoning": r.Counts.Reasoning, "total": r.Counts.Total,
		"record_json": string(raw), "recorded_at": r.RecordedAt}}, nil
}

func validateUsage(tx *sql.Tx, e Event) error {
	if e.CommandName != "usage.record" || e.Aggregate.Sequence != 1 || e.Repository != nil || len(e.Effects) != 2 {
		return fmt.Errorf("invalid usage Event")
	}
	if err := validateCurrentIdempotency(e); err != nil {
		return err
	}
	effect := e.Effects[0]
	if err := validateCurrentEffect(e, effect); err != nil {
		return err
	}
	raw, _ := effect.Values["record_json"].(string)
	var r usagerecord.Record
	if err := unmarshalCurrentJSON(raw, &r); err != nil {
		return err
	}
	object, err := decodeCurrentJSON(raw)
	if err != nil {
		return err
	}
	if !usageEqual(object, r) {
		return fmt.Errorf("usage record must contain exactly its complete fields")
	}
	if err := r.Validate(); err != nil {
		return err
	}
	if !currentDigest(r.InputDigest) || r.ID != e.Aggregate.ID || r.Actor != e.ActorV2.ActorType+":"+e.ActorV2.ActorID || r.RecordedAt != e.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("usage Event attribution mismatch")
	}
	session, err := usagerecord.ReadSession(tx, r.HoldID)
	if err != nil {
		return fmt.Errorf("usage session: %w", err)
	}
	if session != (usagerecord.Session{TaskID: r.TaskID, Actor: r.Actor, Model: r.Model, Role: r.Role, Stream: r.Stream}) {
		return fmt.Errorf("usage differs from recorded session")
	}
	want, err := UsageEffect(r)
	if err != nil {
		return err
	}
	if !currentValuesEqual(effect, want) {
		return fmt.Errorf("usage projection differs from record")
	}
	retry := e.Effects[1]
	result, _ := retry.Values["result_json"].(string)
	stored, err := decodeCurrentJSON(result)
	if err != nil {
		return err
	}
	if !usageEqual(stored, map[string]any{"usage": r, "replayed": false}) {
		return fmt.Errorf("usage retry differs from record")
	}
	var count int
	if err := tx.QueryRow(`SELECT count(*) FROM usage_sessions WHERE hold_id=? AND input_digest=?`, r.HoldID, r.InputDigest).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("duplicate usage output for session")
	}
	return nil
}

// ValidateUsageRetry binds a stored response to the immutable row, even when
// the projection was already fresh and its Event did not need replay.
func ValidateUsageRetry(conn *sql.DB, raw, id string) error {
	var result struct {
		Usage usagerecord.Record `json:"usage"`
	}
	if err := unmarshalCurrentJSON(raw, &result); err != nil {
		return err
	}
	var record string
	if err := conn.QueryRow(`SELECT record_json FROM usage_sessions WHERE id=?`, id).Scan(&record); err != nil {
		return err
	}
	var original usagerecord.Record
	if err := unmarshalCurrentJSON(record, &original); err != nil {
		return err
	}
	if original.ID != id || !currentValuesEqual(result.Usage, original) {
		return fmt.Errorf("usage retry differs from immutable record")
	}
	return original.Validate()
}

func usageEqual(a, b any) bool {
	left, err := MarshalCanonicalJSON(a)
	if err != nil {
		return false
	}
	right, err := MarshalCanonicalJSON(b)
	return err == nil && string(left) == string(right)
}
