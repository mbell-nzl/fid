package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

const capabilitiesGolden = `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"capabilities.get","result":{"items":[{"available":false,"capability_id":"attention.v1","limits":{"state":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":true,"capability_id":"evidence.freshness.v1","limits":{"branch_labels":false,"target":"repository-tree"},"schema_max":1,"schema_min":1,"trust_source":"fid-git-observation"},{"available":true,"capability_id":"evidence.gate-definitions.v1","limits":{"executor_policy":"standard-library","network":"allowed","sandbox":"none","write":"allowed"},"schema_max":1,"schema_min":1,"trust_source":"active-kit-activation"},{"available":true,"capability_id":"evidence.receipts.v1","limits":{"attestation":false,"producer":"fid","terminal_receipts":true},"schema_max":1,"schema_min":1,"trust_source":"fid-executed"},{"available":false,"capability_id":"export.portable.v1","limits":{"create":true,"import":false,"validate":true},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":false,"capability_id":"human.outcomes.v1","limits":{"record":false},"schema_max":1,"schema_min":1,"trust_source":"phase-23-unavailable"},{"available":false,"capability_id":"kit.activation.v1","limits":{"abort_before_commit":true,"automatic_respawn":false,"finalize":true,"recover":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":false,"capability_id":"reviews.coverage.v1","limits":{"record":false},"schema_max":1,"schema_min":1,"trust_source":"phase-23-unavailable"},{"available":false,"capability_id":"telemetry.trajectory.v1","limits":{"ingest":false},"schema_max":1,"schema_min":1,"trust_source":"phase-25-unavailable"},{"available":false,"capability_id":"work.claims.v1","limits":{"max_active_per_issue":1,"parallelism":false,"recover":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":true,"capability_id":"workflow.state.v1","limits":{"max_envelope_revisions":1,"max_primary_issues":1,"related_issues":false},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"},{"available":true,"capability_id":"workflow.transitions.v1","limits":{"execution_bound":true,"terminal_evidence_guard":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"}]}}`

func TestCapabilitiesJSONGoldenAndProjectionRows(t *testing.T) {
	workspace := clonePristineWorkspace(t)
	result := executeCommandForTest(workspace, "--json", "capabilities")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("capabilities = %#v", result)
	}
	if got, want := result.stdout, capabilitiesGolden+"\n"; got != want {
		t.Fatalf("capabilities JSON = %s\nwant %s", got, want)
	}

	conn := openCapabilitiesDB(t, workspace)
	defer closeCommandDB(t, conn)
	rows, err := conn.Query(`
		SELECT capability_id, schema_min, schema_max, available, limits_json, trust_source
		FROM system_capabilities
		ORDER BY capability_id`)
	if err != nil {
		t.Fatalf("query capability projection = %v", err)
	}
	defer func() { _ = rows.Close() }()

	var projected []string
	for rows.Next() {
		var id, limits, trust string
		var schemaMin, schemaMax int
		var available bool
		if err := rows.Scan(&id, &schemaMin, &schemaMax, &available, &limits, &trust); err != nil {
			t.Fatalf("scan capability projection = %v", err)
		}
		projected = append(projected, fmt.Sprintf("%s|%d|%d|%t|%s|%s", id, schemaMin, schemaMax, available, limits, trust))
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate capability projection = %v", err)
	}
	if got, want := len(projected), 12; got != want {
		t.Fatalf("projected capability rows = %d, want %d: %v", got, want, projected)
	}
}

func TestCapabilitiesRejectsOverridesAndKeepsDeferredRowsUnavailable(t *testing.T) {
	workspace := clonePristineWorkspace(t)
	if result := executeCommandForTest(workspace, "--json", "capabilities"); result.err != nil {
		t.Fatalf("initialize capability projection = %#v", result)
	}

	invalid := executeCommandForTest(workspace, "--json", "capabilities", "available=true")
	if invalid.err == nil || invalid.stdout != "" {
		t.Fatalf("capabilities override = %#v", invalid)
	}
	const wantInvalid = `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["no arguments"],"argument":"arguments","value":"available=true"},"message":"capabilities accepts no arguments","retryable":false,"subject":null},"ok":false,"operation":"capabilities.get"}` + "\n"
	if invalid.stderr != wantInvalid {
		t.Fatalf("capabilities override error = %s\nwant %s", invalid.stderr, wantInvalid)
	}

	conn := openCapabilitiesDB(t, workspace)
	defer closeCommandDB(t, conn)
	for _, fixture := range []struct {
		capabilityID string
		surfaces     string
	}{
		{"kit.activation.v1", "Activation lifecycle"},
		{"work.claims.v1", "Claim lifecycle and list"},
		{"export.portable.v1", "export create and validate"},
		{"reviews.coverage.v1", "fail-closed AC and Concern judgment"},
		{"human.outcomes.v1", "fail-closed AC and Concern judgment"},
		{"telemetry.trajectory.v1", "historical stats"},
		{"attention.v1", "Attention-absence surfaces"},
	} {
		t.Run(fixture.capabilityID+" "+fixture.surfaces, func(t *testing.T) {
			if _, err := conn.Exec(`UPDATE system_capabilities SET available = 1 WHERE capability_id = ?`, fixture.capabilityID); err != nil {
				t.Fatalf("optimistically update %s = %v", fixture.capabilityID, err)
			}
			result := executeCommandForTest(workspace, "--json", "capabilities")
			if result.err != nil || result.stderr != "" {
				t.Fatalf("capabilities after %s override = %#v", fixture.capabilityID, result)
			}
			if available := capabilityAvailability(t, result.stdout, fixture.capabilityID); available {
				t.Fatalf("%s availability = true after optimistic projection write, want false for %s", fixture.capabilityID, fixture.surfaces)
			}
		})
	}
}

func TestDeferredCapabilitySurfacesStayUnavailableWithoutSideEffects(t *testing.T) {
	workspace := clonePristineWorkspace(t)
	initial := executeCommandForTest(workspace, "--json", "capabilities")
	if initial.err != nil || initial.stderr != "" {
		t.Fatalf("initialize capability projection = %#v", initial)
	}

	for _, state := range []struct {
		capabilityID string
		available    bool
	}{
		{"attention.v1", false},
		{"evidence.freshness.v1", true},
		{"evidence.gate-definitions.v1", true},
		{"evidence.receipts.v1", true},
		{"export.portable.v1", false},
		{"human.outcomes.v1", false},
		{"kit.activation.v1", false},
		{"reviews.coverage.v1", false},
		{"telemetry.trajectory.v1", false},
		{"work.claims.v1", false},
		{"workflow.state.v1", true},
		{"workflow.transitions.v1", true},
	} {
		t.Run("state "+state.capabilityID, func(t *testing.T) {
			if got := capabilityAvailability(t, initial.stdout, state.capabilityID); got != state.available {
				t.Fatalf("%s availability = %t, want %t", state.capabilityID, got, state.available)
			}
		})
	}

	journalPath := mustJournalPath(t, workspace)
	for _, probe := range []struct {
		capabilityID string
		surface      string
		args         []string
		errorText    string
	}{
		{
			capabilityID: "export.portable.v1",
			surface:      "export create and validate",
			args:         []string{"--json", "export", "create"},
			errorText:    `unknown command "export"`,
		},
		{
			capabilityID: "reviews.coverage.v1",
			surface:      "fail-closed AC judgment",
			args:         []string{"--json", "issue", "ac", "verify", "ISS-1", "AC-1"},
			errorText:    `"code":"capability-unavailable"`,
		},
		{
			capabilityID: "human.outcomes.v1",
			surface:      "fail-closed Concern judgment",
			args:         []string{"--json", "concern", "raise"},
			errorText:    `"code":"capability-unavailable"`,
		},
		{
			capabilityID: "attention.v1",
			surface:      "Attention absence",
			args:         []string{"--json", "attention", "list"},
			errorText:    `unknown command "attention"`,
		},
	} {
		probe := probe
		t.Run(probe.capabilityID+" "+probe.surface, func(t *testing.T) {
			journalBefore := mustReadCommandFile(t, journalPath)
			projectionBefore := commandProjectionTableCounts(t, workspace)

			result := executeCommandForTest(workspace, probe.args...)
			if result.err == nil || result.stdout != "" {
				t.Fatalf("%s probe = %#v, want unavailable without stdout", probe.surface, result)
			}
			failureText := result.stderr + "\n" + result.err.Error()
			if !strings.Contains(failureText, probe.errorText) {
				t.Fatalf("%s failure = %q, want unavailable surface error containing %q", probe.surface, failureText, probe.errorText)
			}
			if journalAfter := mustReadCommandFile(t, journalPath); !bytes.Equal(journalAfter, journalBefore) {
				t.Fatalf("%s probe changed journal bytes", probe.surface)
			}
			if projectionAfter := commandProjectionTableCounts(t, workspace); !reflect.DeepEqual(projectionAfter, projectionBefore) {
				t.Fatalf("%s probe changed projection counts: got %#v, want %#v", probe.surface, projectionAfter, projectionBefore)
			}

			result = executeCommandForTest(workspace, "--json", "capabilities")
			if result.err != nil || result.stderr != "" {
				t.Fatalf("capability discovery after %s probe = %#v", probe.surface, result)
			}
			if capabilityAvailability(t, result.stdout, probe.capabilityID) {
				t.Fatalf("%s became available after %s probe", probe.capabilityID, probe.surface)
			}
		})
	}
}

func openCapabilitiesDB(t *testing.T, workspace string) *sql.DB {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open capability projection = %v", err)
	}
	return conn
}

func capabilityAvailability(t *testing.T, output, capabilityID string) bool {
	t.Helper()
	var envelope struct {
		Result struct {
			Items []struct {
				CapabilityID string `json:"capability_id"`
				Available    bool   `json:"available"`
			} `json:"items"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(output), &envelope); err != nil {
		t.Fatalf("decode capabilities result = %v", err)
	}
	for _, item := range envelope.Result.Items {
		if item.CapabilityID == capabilityID {
			return item.Available
		}
	}
	t.Fatalf("capability %s absent from discovery result", capabilityID)
	return false
}
