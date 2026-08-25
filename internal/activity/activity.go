package activity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type ListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	On            string
	Since         string
	Event         string
}

type Entry struct {
	ID             string `json:"id"`
	JournalEventID string `json:"journal_event_id"`
	RefType        string `json:"ref_type"`
	RefID          string `json:"ref_id"`
	Event          string `json:"event"`
	Actor          string `json:"actor"`
	OccurredAt     string `json:"occurred_at"`
	Details        any    `json:"details"`
}

func List(options ListOptions) ([]Entry, error) {
	workspaceRoot, migrationsDir, on, since, err := validateListOptions(options)
	if err != nil {
		return nil, err
	}

	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return nil, err
	}

	conn, err := fiddb.Open(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer closeActivityConn(conn)

	query, args, err := activityListQuery(on, since, "")
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list activity: %w", err)
	}
	defer closeActivityRows(rows)

	items := []Entry{}
	wantEvent := normalizeActivityEvent(options.Event)
	for rows.Next() {
		item, err := scanEntry(rows)
		if err != nil {
			return nil, err
		}
		if wantEvent == "" || item.Event == wantEvent {
			items = append(items, item)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate activity rows: %w", err)
	}

	return items, nil
}

func validateListOptions(options ListOptions) (string, string, onFilter, string, error) {
	workspaceRoot := strings.TrimSpace(options.WorkspaceRoot)
	if workspaceRoot == "" {
		return "", "", onFilter{}, "", fmt.Errorf("activity list: workspace root is required")
	}

	migrationsDir := strings.TrimSpace(options.MigrationsDir)
	if migrationsDir == "" {
		return "", "", onFilter{}, "", fmt.Errorf("activity list: migrations dir is required")
	}

	on, err := parseOnFilter(options.On)
	if err != nil {
		return "", "", onFilter{}, "", err
	}

	since, err := normalizeSince(options.Since)
	if err != nil {
		return "", "", onFilter{}, "", err
	}

	return workspaceRoot, migrationsDir, on, since, nil
}

type onFilter struct {
	refType string
	refID   string
}

func parseOnFilter(value string) (onFilter, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return onFilter{}, nil
	}

	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return onFilter{}, fmt.Errorf("activity list: --on must be <type>/<id>")
	}

	return onFilter{
		refType: strings.TrimSpace(parts[0]),
		refID:   strings.TrimSpace(parts[1]),
	}, nil
}

func normalizeSince(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}

	if parsed, err := time.Parse("2006-01-02", trimmed); err == nil {
		return parsed.UTC().Format(time.RFC3339), nil
	}

	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return "", fmt.Errorf("activity list: --since must be YYYY-MM-DD or RFC3339")
	}

	return parsed.Format(time.RFC3339), nil
}

func activityListQuery(on onFilter, since, event string) (string, []any, error) {
	query := strings.Builder{}
	query.WriteString(`
SELECT
	id,
	journal_event_id,
	ref_type,
	ref_id,
	event,
	actor,
	occurred_at,
	details
FROM activity_log`)

	whereParts := []string{}
	args := []any{}
	if on.refType != "" {
		whereParts = append(whereParts, "ref_type = ?", "ref_id = ?")
		args = append(args, on.refType, on.refID)
	}
	if since != "" {
		whereParts = append(whereParts, "occurred_at >= ?")
		args = append(args, since)
	}
	if event != "" {
		whereParts = append(whereParts, "event = ?")
		args = append(args, event)
	}
	if len(whereParts) > 0 {
		query.WriteString("\nWHERE ")
		query.WriteString(strings.Join(whereParts, "\n  AND "))
	}
	query.WriteString("\nORDER BY occurred_at ASC, id ASC")
	return query.String(), args, nil
}

func scanEntry(rows *sql.Rows) (Entry, error) {
	var (
		entry       Entry
		actor       sql.NullString
		occurredAt  sql.NullString
		detailsText sql.NullString
	)

	if err := rows.Scan(
		&entry.ID,
		&entry.JournalEventID,
		&entry.RefType,
		&entry.RefID,
		&entry.Event,
		&actor,
		&occurredAt,
		&detailsText,
	); err != nil {
		return Entry{}, fmt.Errorf("scan activity row: %w", err)
	}

	entry.Actor = nullStringValue(actor)
	entry.OccurredAt = nullStringValue(occurredAt)

	details, err := parseDetails(detailsText)
	if err != nil {
		return Entry{}, err
	}
	entry.Details = details
	entry.Event = normalizeActivityEvent(entry.Event)

	return entry, nil
}

func normalizeActivityEvent(value string) string {
	return strings.ReplaceAll(strings.TrimSpace(value), "_", "-")
}

func parseDetails(value sql.NullString) (any, error) {
	if !value.Valid || strings.TrimSpace(value.String) == "" {
		return nil, nil
	}

	decoder := json.NewDecoder(strings.NewReader(value.String))
	decoder.UseNumber()

	var details any
	if err := decoder.Decode(&details); err != nil {
		return nil, fmt.Errorf("decode activity details: %w", err)
	}
	return details, nil
}

func nullStringValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

func closeActivityConn(conn interface{ Close() error }) {
	_ = conn.Close()
}

func closeActivityRows(rows interface{ Close() error }) {
	_ = rows.Close()
}
