package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Capability is one binary-derived capability availability row.
type Capability struct {
	Available    bool           `json:"available"`
	CapabilityID string         `json:"capability_id"`
	Limits       map[string]any `json:"limits"`
	SchemaMax    int            `json:"schema_max"`
	SchemaMin    int            `json:"schema_min"`
	TrustSource  string         `json:"trust_source"`
}

// DiscoverCapabilities reads the current binary capability matrix from the fresh projection.
func DiscoverCapabilities(workspaceRoot, migrationsDir string) ([]Capability, error) {
	if strings.TrimSpace(workspaceRoot) == "" {
		return nil, fmt.Errorf("discover capabilities: workspace root is required")
	}
	if strings.TrimSpace(migrationsDir) == "" {
		return nil, fmt.Errorf("discover capabilities: migrations dir is required")
	}
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return nil, fmt.Errorf("discover capabilities: ensure projection: %w", err)
	}

	conn, err := fiddb.Open(config.WorkspacePathsForRoot(workspaceRoot).DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("discover capabilities: open projection: %w", err)
	}
	defer closeWorkflowDB(conn)
	if err := materializeCapabilities(conn); err != nil {
		return nil, err
	}
	return readCapabilities(conn)
}

func materializeCapabilities(conn *sql.DB) error {
	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("materialize capabilities: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.Exec(`DELETE FROM system_capabilities`); err != nil {
		return fmt.Errorf("materialize capabilities: clear rows: %w", err)
	}
	for _, row := range capabilityMatrix {
		if _, err := tx.Exec(`
			INSERT INTO system_capabilities (capability_id, schema_min, schema_max, available, limits_json, trust_source)
			VALUES (?, 1, 1, ?, ?, ?)`, row.capabilityID, row.available, row.limitsJSON, row.trustSource); err != nil {
			return fmt.Errorf("materialize capabilities: insert %q: %w", row.capabilityID, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("materialize capabilities: commit: %w", err)
	}
	return nil
}

func readCapabilities(conn *sql.DB) ([]Capability, error) {
	rows, err := conn.Query(`
		SELECT capability_id, schema_min, schema_max, available, limits_json, trust_source
		FROM system_capabilities
		ORDER BY capability_id`)
	if err != nil {
		return nil, fmt.Errorf("discover capabilities: query rows: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items := make([]Capability, 0, len(capabilityMatrix))
	for rows.Next() {
		var item Capability
		var limitsJSON string
		if err := rows.Scan(&item.CapabilityID, &item.SchemaMin, &item.SchemaMax, &item.Available, &limitsJSON, &item.TrustSource); err != nil {
			return nil, fmt.Errorf("discover capabilities: scan row: %w", err)
		}
		if err := json.Unmarshal([]byte(limitsJSON), &item.Limits); err != nil {
			return nil, fmt.Errorf("discover capabilities: decode limits for %q: %w", item.CapabilityID, err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("discover capabilities: iterate rows: %w", err)
	}
	return items, nil
}

type capabilityMatrixRow struct {
	capabilityID string
	available    bool
	limitsJSON   string
	trustSource  string
}

var capabilityMatrix = []capabilityMatrixRow{
	{"attention.v1", true, `{"human_outcome_resolution":false,"state":true}`, "fid-journal"},
	{"evidence.freshness.v1", true, `{"branch_labels":false,"target":"repository-tree"}`, "fid-git-observation"},
	{"evidence.gate-definitions.v1", true, `{"executor_policy":"standard-library","network":"allowed","sandbox":"none","write":"allowed"}`, "active-kit-activation"},
	{"evidence.receipts.v1", true, `{"attestation":false,"producer":"fid","terminal_receipts":true}`, "fid-executed"},
	{"export.portable.v1", true, `{"create":true,"import":true,"import_mode":"empty-target","validate":true}`, "fid-journal"},
	{"human.outcomes.v1", false, `{"record":false}`, "phase-23-unavailable"},
	{"kit.activation.v1", true, `{"abort_before_commit":true,"automatic_respawn":false,"finalize":true,"recover":true}`, "active-kit-activation"},
	{"reviews.coverage.v1", false, `{"record":false}`, "phase-23-unavailable"},
	{"telemetry.trajectory.v1", false, `{"ingest":false}`, "phase-25-unavailable"},
	{"work.claims.v1", false, `{"max_active_per_issue":1,"parallelism":false,"recover":true}`, "phase-24-unavailable"},
	{"workflow.state.v1", true, `{"max_envelope_revisions":1,"max_primary_issues":1,"related_issues":false}`, "fid-journal"},
	{"workflow.transitions.v1", true, `{"execution_bound":true,"terminal_evidence_guard":true}`, "fid-journal"},
}
