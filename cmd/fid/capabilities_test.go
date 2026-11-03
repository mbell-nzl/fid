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

const legacyCapabilitiesGolden = `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"capabilities.get","result":{"items":[{"available":false,"capability_id":"attention.v1","limits":{"state":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":true,"capability_id":"evidence.freshness.v1","limits":{"branch_labels":false,"target":"repository-tree"},"schema_max":1,"schema_min":1,"trust_source":"fid-git-observation"},{"available":true,"capability_id":"evidence.gate-definitions.v1","limits":{"executor_policy":"standard-library","network":"allowed","sandbox":"none","write":"allowed"},"schema_max":1,"schema_min":1,"trust_source":"active-kit-activation"},{"available":true,"capability_id":"evidence.receipts.v1","limits":{"attestation":false,"producer":"fid","terminal_receipts":true},"schema_max":1,"schema_min":1,"trust_source":"fid-executed"},{"available":false,"capability_id":"export.portable.v1","limits":{"create":true,"import":false,"validate":true},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":false,"capability_id":"human.outcomes.v1","limits":{"record":false},"schema_max":1,"schema_min":1,"trust_source":"phase-23-unavailable"},{"available":false,"capability_id":"kit.activation.v1","limits":{"abort_before_commit":true,"automatic_respawn":false,"finalize":true,"recover":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":false,"capability_id":"reviews.coverage.v1","limits":{"record":false},"schema_max":1,"schema_min":1,"trust_source":"phase-23-unavailable"},{"available":false,"capability_id":"telemetry.trajectory.v1","limits":{"ingest":false},"schema_max":1,"schema_min":1,"trust_source":"phase-25-unavailable"},{"available":false,"capability_id":"work.claims.v1","limits":{"max_active_per_issue":1,"parallelism":false,"recover":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":true,"capability_id":"workflow.state.v1","limits":{"max_envelope_revisions":1,"max_primary_issues":1,"related_issues":false},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"},{"available":true,"capability_id":"workflow.transitions.v1","limits":{"execution_bound":true,"terminal_evidence_guard":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"}]}}`

var capabilitiesGolden = strings.NewReplacer(
	`{"available":false,"capability_id":"attention.v1","limits":{"state":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"}`,
	`{"available":true,"capability_id":"attention.v1","limits":{"human_outcome_resolution":true,"state":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"}`,
	`{"available":false,"capability_id":"human.outcomes.v1","limits":{"record":false},"schema_max":1,"schema_min":1,"trust_source":"phase-23-unavailable"}`,
	`{"available":true,"capability_id":"human.outcomes.v1","limits":{"provenance":"active-kit-activation","record":true},"schema_max":1,"schema_min":1,"trust_source":"active-kit-activation"}`,
	`{"available":false,"capability_id":"reviews.coverage.v1","limits":{"record":false},"schema_max":1,"schema_min":1,"trust_source":"phase-23-unavailable"}`,
	`{"available":true,"capability_id":"reviews.coverage.v1","limits":{"findings":true,"independence":"distinct-actor","record":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"}`,
	`{"available":false,"capability_id":"export.portable.v1","limits":{"create":true,"import":false,"validate":true},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"}`,
	`{"available":true,"capability_id":"export.portable.v1","limits":{"create":true,"import":true,"import_mode":"empty-target","validate":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"}`,
	`{"available":false,"capability_id":"kit.activation.v1","limits":{"abort_before_commit":true,"automatic_respawn":false,"finalize":true,"recover":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"}`,
	`{"available":true,"capability_id":"kit.activation.v1","limits":{"abort_before_commit":true,"automatic_respawn":false,"finalize":true,"recover":true},"schema_max":1,"schema_min":1,"trust_source":"active-kit-activation"}`,
	`{"available":false,"capability_id":"work.claims.v1","limits":{"max_active_per_issue":1,"parallelism":false,"recover":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"}`,
	`{"available":false,"capability_id":"work.claims.v1","limits":{"max_active_per_issue":1,"parallelism":false,"recover":true},"schema_max":1,"schema_min":1,"trust_source":"phase-24-unavailable"}`,
).Replace(legacyCapabilitiesGolden)

func TestCapabilitiesJSONGoldenAndProjectionRows(t *testing.T) {
	t.Parallel()
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

func TestCapabilitiesAdvertiseOnlyFinalQualifiedRows(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	result := executeCommandForTest(workspace, "--json", "capabilities")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("capabilities = %#v", result)
	}
	var envelope struct {
		Result struct {
			Items []struct {
				Available    bool           `json:"available"`
				CapabilityID string         `json:"capability_id"`
				Limits       map[string]any `json:"limits"`
				TrustSource  string         `json:"trust_source"`
			} `json:"items"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("decode capability discovery: %v", err)
	}
	got := make(map[string]struct {
		available bool
		limits    map[string]any
		trust     string
	})
	for _, item := range envelope.Result.Items {
		got[item.CapabilityID] = struct {
			available bool
			limits    map[string]any
			trust     string
		}{item.Available, item.Limits, item.TrustSource}
	}
	for id, want := range map[string]struct {
		available bool
		limits    map[string]any
		trust     string
	}{
		"attention.v1":        {true, map[string]any{"state": true, "human_outcome_resolution": true}, "fid-journal"},
		"human.outcomes.v1":   {true, map[string]any{"record": true, "provenance": "active-kit-activation"}, "active-kit-activation"},
		"reviews.coverage.v1": {true, map[string]any{"record": true, "independence": "distinct-actor", "findings": true}, "fid-journal"},
		"kit.activation.v1":   {true, map[string]any{"abort_before_commit": true, "automatic_respawn": false, "finalize": true, "recover": true}, "active-kit-activation"},
		"export.portable.v1":  {true, map[string]any{"create": true, "import": true, "validate": true, "import_mode": "empty-target"}, "fid-journal"},
		"work.claims.v1":      {false, map[string]any{"max_active_per_issue": float64(1), "parallelism": false, "recover": true}, "phase-24-unavailable"},
	} {
		if actual, ok := got[id]; !ok || actual.available != want.available || actual.trust != want.trust || !reflect.DeepEqual(actual.limits, want.limits) {
			t.Fatalf("%s = %#v, want %#v", id, actual, want)
		}
	}
}

func TestCapabilitiesRejectsOverridesAndKeepsDeferredRowsUnavailable(t *testing.T) {
	t.Parallel()
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
		{"reviews.coverage.v1", "Review coverage and Findings"},
		{"human.outcomes.v1", "Human Outcome judgment"},
		{"telemetry.trajectory.v1", "historical stats"},
		{"attention.v1", "Attention state surfaces"},
	} {
		t.Run(fixture.capabilityID+" "+fixture.surfaces, func(t *testing.T) {
			if _, err := conn.Exec(`UPDATE system_capabilities SET available = 1 WHERE capability_id = ?`, fixture.capabilityID); err != nil {
				t.Fatalf("optimistically update %s = %v", fixture.capabilityID, err)
			}
			result := executeCommandForTest(workspace, "--json", "capabilities")
			if result.err != nil || result.stderr != "" {
				t.Fatalf("capabilities after %s override = %#v", fixture.capabilityID, result)
			}
			if available := capabilityAvailability(t, result.stdout, fixture.capabilityID); available != (fixture.capabilityID != "work.claims.v1" && fixture.capabilityID != "telemetry.trajectory.v1") {
				t.Fatalf("%s availability = %t after optimistic projection write, want binary-qualified state for %s", fixture.capabilityID, available, fixture.surfaces)
			}
		})
	}
}

func TestRemovedLeavesRejectWithoutSideEffectsAndDeferredRowsStayUnavailable(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	initial := executeCommandForTest(workspace, "--json", "capabilities")
	if initial.err != nil || initial.stderr != "" {
		t.Fatalf("initialize capability projection = %#v", initial)
	}

	for _, state := range []struct {
		capabilityID string
		available    bool
	}{
		{"attention.v1", true},
		{"evidence.freshness.v1", true},
		{"evidence.gate-definitions.v1", true},
		{"evidence.receipts.v1", true},
		{"export.portable.v1", true},
		{"human.outcomes.v1", true},
		{"kit.activation.v1", true},
		{"reviews.coverage.v1", true},
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
			surface:      "export import required arguments",
			args:         []string{"--json", "export", "import"},
			errorText:    `"code":"invalid-argument"`,
		},
		{
			capabilityID: "reviews.coverage.v1",
			surface:      "removed AC verification leaf",
			args:         []string{"issue", "ac", "verify"},
			errorText:    `unknown command "verify"`,
		},
		{
			capabilityID: "human.outcomes.v1",
			surface:      "removed Concern raise leaf",
			args:         []string{"concern", "raise"},
			errorText:    `unknown command "raise"`,
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
			if available := capabilityAvailability(t, result.stdout, probe.capabilityID); !available {
				t.Fatalf("%s availability = %t after %s probe, want binary-qualified state", probe.capabilityID, available, probe.surface)
			}
		})
	}
}

func openCapabilitiesDB(t *testing.T, workspace string) *sql.DB {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
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
