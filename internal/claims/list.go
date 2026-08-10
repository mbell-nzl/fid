package claims

import (
	"database/sql"
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
	ClaimedBy     string
	StaleOnly     bool
	Now           time.Time
}

type ListItem struct {
	RefType    string  `json:"ref_type"`
	RefID      string  `json:"ref_id"`
	ClaimedBy  string  `json:"claimed_by"`
	ClaimedAt  string  `json:"claimed_at"`
	TTLSeconds int     `json:"ttl_seconds"`
	ExpiresAt  string  `json:"expires_at"`
	IsStale    bool    `json:"is_stale"`
	Note       *string `json:"note"`
}

func List(options ListOptions) ([]ListItem, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return nil, fmt.Errorf("list claims: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return nil, fmt.Errorf("list claims: migrations dir is required")
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}

	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer closeClaimsConn(conn)

	return queryClaimList(conn, options, operationTime(options.Now))
}

func queryClaimList(conn *sql.DB, options ListOptions, now time.Time) ([]ListItem, error) {
	query := `
		SELECT ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, note
		FROM claims
	`
	args := []any{}
	if claimedBy := strings.TrimSpace(options.ClaimedBy); claimedBy != "" {
		query += "\nWHERE claimed_by = ?"
		args = append(args, claimedBy)
	}
	query += "\nORDER BY rowid ASC"

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list claims: %w", err)
	}
	defer closeClaimsRows(rows)

	items := []ListItem{}
	for rows.Next() {
		var (
			item ListItem
			note sql.NullString
		)
		if err := rows.Scan(&item.RefType, &item.RefID, &item.ClaimedBy, &item.ClaimedAt, &item.TTLSeconds, &note); err != nil {
			return nil, fmt.Errorf("scan claim list row: %w", err)
		}

		claimedAt, err := time.Parse(time.RFC3339, item.ClaimedAt)
		if err != nil {
			return nil, fmt.Errorf("parse claim time for %q: %w", item.RefID, err)
		}

		current := currentClaim{
			ClaimedBy:  item.ClaimedBy,
			ClaimedAt:  claimedAt.UTC(),
			TTLSeconds: item.TTLSeconds,
		}
		item.ExpiresAt = current.expiresAt().Format(time.RFC3339)
		item.IsStale = !current.isLiveAt(now.UTC())
		item.Note = nullStringPointer(note)

		if options.StaleOnly && !item.IsStale {
			continue
		}
		if !options.StaleOnly && item.IsStale {
			continue
		}

		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate claim list rows: %w", err)
	}

	return items, nil
}
