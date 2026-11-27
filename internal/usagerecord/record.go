package usagerecord

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
)

type Record struct {
	ID          string `json:"id"`
	TaskID      string `json:"task_id"`
	HoldID      string `json:"hold_id"`
	Actor       string `json:"actor"`
	Model       string `json:"model"`
	Role        string `json:"role"`
	Stream      string `json:"stream"`
	InputDigest string `json:"input_digest"`
	RecordedAt  string `json:"recorded_at"`
	Parsed
}

type Session struct{ TaskID, Actor, Model, Role, Stream string }
type Queryer interface{ QueryRow(string, ...any) *sql.Row }

// ReadSession includes ended holds and review reservations. A fix descends from
// a returned hand-off hold; an interrupted fix keeps that ancestry on resume.
func ReadSession(conn Queryer, id string) (Session, error) {
	var s Session
	err := conn.QueryRow(`SELECT task_id,actor,model,stream FROM review_claims WHERE id=?`, id).Scan(&s.TaskID, &s.Actor, &s.Model, &s.Stream)
	if err == nil {
		s.Role = "review"
		return s, nil
	}
	if err != sql.ErrNoRows {
		return Session{}, fmt.Errorf("read review usage session: %w", err)
	}
	var parent sql.NullString
	err = conn.QueryRow(`SELECT task_id,actor,model,replaces_hold_id FROM holds WHERE id=?`, id).Scan(&s.TaskID, &s.Actor, &s.Model, &parent)
	if err != nil {
		return Session{}, err
	}
	s.Role = "build"
	seen := map[string]bool{id: true}
	for parent.Valid {
		if seen[parent.String] {
			return Session{}, fmt.Errorf("cyclic hold history")
		}
		seen[parent.String] = true
		var rework sql.NullString
		if err := conn.QueryRow(`SELECT rework_at,replaces_hold_id FROM holds WHERE id=?`, parent.String).Scan(&rework, &parent); err != nil {
			return Session{}, err
		}
		if rework.Valid {
			s.Role = "fix"
			break
		}
	}
	return s, nil
}

// Validate re-derives counts from the retained agent report using its frozen
// parser revision. Historical numbers never silently change with new rules.
func (r Record) Validate() error {
	if r.ID == "" || r.TaskID == "" || r.HoldID == "" || r.Actor == "" || r.RecordedAt == "" {
		return fmt.Errorf("usage record identity is incomplete")
	}
	if r.ParserRevision != r.Format+"/1" {
		return fmt.Errorf("unknown usage parser revision")
	}
	if r.Role != "build" && r.Role != "fix" && r.Role != "review" || r.Role == "review" && r.Stream != "primary" && r.Stream != "separate" || r.Role != "review" && r.Stream != "" {
		return fmt.Errorf("invalid usage role or stream")
	}
	if r.ToolVersion == "" && r.VersionSource != "" || r.ToolVersion != "" && r.VersionSource != "output" && r.VersionSource != "caller" || r.VersionSource == "output" && r.Format != "claude" {
		return fmt.Errorf("invalid tool version provenance")
	}
	if len(r.Reported) == 0 {
		return fmt.Errorf("usage has no reported counts")
	}
	var raw []byte
	switch r.Format {
	case "codex":
		for _, report := range r.Reported {
			event, err := json.Marshal(map[string]any{"type": "turn.completed", "usage": report})
			if err != nil {
				return err
			}
			raw = append(raw, event...)
			raw = append(raw, '\n')
		}
	case "claude":
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(r.Reported[0], &fields); err != nil {
			return err
		}
		key := "modelUsage"
		if value, exists := fields["input_tokens"]; exists {
			var count int64
			if json.Unmarshal(value, &count) == nil {
				key = "usage"
			}
		}
		event := map[string]any{"type": "result", key: r.Reported[0]}
		if len(r.Reported) > 2 || key == "usage" && len(r.Reported) != 1 {
			return fmt.Errorf("invalid retained Claude reports")
		}
		if len(r.Reported) == 2 {
			event["usage"] = r.Reported[1]
		}
		var err error
		raw, err = json.Marshal(event)
		if err != nil {
			return err
		}
	}
	parsed, err := Parse(r.Format, raw, r.Model)
	if err != nil {
		return err
	}
	if parsed.Counts != r.Counts || !reflect.DeepEqual(parsed.Models, r.Models) {
		return fmt.Errorf("usage counts differ from reported counts")
	}
	return nil
}
