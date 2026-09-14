package workflow

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

type StopListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Subject       string
	Status        string
	Limit         int
}

// ListStops returns immutable Stop facts without appending a journal Event.
func ListStops(options StopListOptions) ([]Stop, error) {
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		return nil, err
	}
	defer closeWorkflowDB(conn)
	query, args, err := stopListQuery(options)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list Stops: %w", err)
	}
	defer func() { _ = rows.Close() }()
	items := []Stop{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		item, err := loadStop(conn, id)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// ListAttention returns linked Stop identities without appending a journal Event.
func ListAttention(options StopListOptions) ([]Attention, error) {
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		return nil, err
	}
	defer closeWorkflowDB(conn)
	query, args, err := attentionListQuery(options)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list Attention: %w", err)
	}
	defer func() { _ = rows.Close() }()
	items := []Attention{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		item, err := loadAttentionByQuery(conn, `SELECT id, semantic_key, subject_type, subject_id, code, status, version, raised_at, resolution_event_id FROM attention_items WHERE id = ?`, id)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func stopListQuery(options StopListOptions) (string, []any, error) {
	where, args, err := stopListFilters(options, "s.subject_type", "s.subject_id", "CASE WHEN r.stop_id IS NULL THEN 'open' ELSE 'resolved' END")
	return "SELECT s.id FROM stops s LEFT JOIN stop_resolutions r ON r.stop_id = s.id" + where + " ORDER BY s.raised_at, s.id LIMIT ?", append(args, listLimit(options.Limit)), err
}

func attentionListQuery(options StopListOptions) (string, []any, error) {
	where, args, err := stopListFilters(options, "a.subject_type", "a.subject_id", "a.status")
	return "SELECT a.id FROM attention_items a" + where + " ORDER BY a.raised_at, a.id LIMIT ?", append(args, listLimit(options.Limit)), err
}

func stopListFilters(options StopListOptions, typeColumn, idColumn, statusColumn string) (string, []any, error) {
	where := []string{}
	args := []any{}
	if options.Subject != "" {
		parts := strings.Split(options.Subject, ":")
		if len(parts) != 2 || !validSubject(Subject{Type: parts[0], ID: parts[1]}) {
			return "", nil, fmt.Errorf("Stop subject filter is invalid")
		}
		where, args = append(where, typeColumn+" = ?", idColumn+" = ?"), append(args, parts[0], parts[1])
	}
	if options.Status != "" {
		if options.Status != "open" && options.Status != "resolved" {
			return "", nil, fmt.Errorf("Stop status filter is invalid")
		}
		where, args = append(where, statusColumn+" = ?"), append(args, options.Status)
	}
	if len(where) == 0 {
		return "", args, nil
	}
	return " WHERE " + strings.Join(where, " AND "), args, nil
}

func listLimit(value int) int {
	if value <= 0 {
		return 100
	}
	return value
}
