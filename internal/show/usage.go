package show

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/mbell-nzl/fid/internal/reviewwork"
	"github.com/mbell-nzl/fid/internal/usagerecord"
)

type Usage struct {
	Total    int64            `json:"total"`
	Sessions int              `json:"sessions"`
	ByRole   map[string]int64 `json:"by_role"`
	ByModel  map[string]int64 `json:"by_model"`
}

func (result *Result) readUsage(root string) error {
	if result.Issue.Kind != "task" {
		return nil
	}
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	rows, err := conn.Query(`SELECT record_json FROM usage_sessions WHERE task_id=? ORDER BY rowid`, result.Issue.ID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	u := Usage{ByRole: map[string]int64{}, ByModel: map[string]int64{}}
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return err
		}
		var r usagerecord.Record
		if err := json.Unmarshal([]byte(raw), &r); err != nil {
			return err
		}
		if err := r.Validate(); err != nil {
			return err
		}
		if u.Total > math.MaxInt64-r.Counts.Total {
			return fmt.Errorf("task usage total exceeds signed 64-bit range")
		}
		u.Total += r.Counts.Total
		u.Sessions++
		role := r.Role
		if r.Stream != "" {
			role += "." + r.Stream
		}
		u.ByRole[role] += r.Counts.Total
		for model, counts := range r.Models {
			u.ByModel[model] += counts.Total
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if u.Sessions != 0 {
		result.Usage = &u
	}
	return nil
}
func (text textWriter) usage(u *Usage) []string {
	if u == nil {
		return nil
	}
	sessionWord := "sessions"
	if u.Sessions == 1 {
		sessionWord = "session"
	}
	parts := []string{fmt.Sprintf("%d tokens · %d %s", u.Total, u.Sessions, sessionWord)}
	for _, group := range []map[string]int64{u.ByRole, u.ByModel} {
		keys := make([]string, 0, len(group))
		for key := range group {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			label := key
			if label == "" {
				label = "unknown model"
			}
			parts = append(parts, fmt.Sprintf("%s %d", label, group[key]))
		}
	}
	return text.value("Usage", strings.Join(parts, " · "))
}
