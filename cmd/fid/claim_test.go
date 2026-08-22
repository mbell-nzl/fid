package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestClaimListDefaultsToJSONAndListsLiveClaims(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI claim list", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer", "--note", "journal append work")

	stdout, stderr := runFIDCommand(t, workspace, "claim", "list")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	var claims []struct {
		RefType    string  `json:"ref_type"`
		RefID      string  `json:"ref_id"`
		ClaimedBy  string  `json:"claimed_by"`
		ClaimedAt  string  `json:"claimed_at"`
		TTLSeconds int     `json:"ttl_seconds"`
		ExpiresAt  string  `json:"expires_at"`
		IsStale    bool    `json:"is_stale"`
		Note       *string `json:"note"`
	}
	if err := json.Unmarshal([]byte(stdout), &claims); err != nil {
		t.Fatalf("Unmarshal(claim list) error = %v", err)
	}

	if len(claims) != 1 {
		t.Fatalf("claims len = %d, want 1", len(claims))
	}
	if claims[0].RefType != "issue" {
		t.Fatalf("ref_type = %q, want issue", claims[0].RefType)
	}
	if claims[0].RefID != issueID {
		t.Fatalf("ref_id = %q, want %q", claims[0].RefID, issueID)
	}
	if claims[0].ClaimedBy != "developer" {
		t.Fatalf("claimed_by = %q, want developer", claims[0].ClaimedBy)
	}
	if claims[0].ClaimedAt == "" {
		t.Fatal("claimed_at = empty, want timestamp")
	}
	if claims[0].TTLSeconds != 1800 {
		t.Fatalf("ttl_seconds = %d, want 1800", claims[0].TTLSeconds)
	}
	if claims[0].ExpiresAt == "" {
		t.Fatal("expires_at = empty, want timestamp")
	}
	if claims[0].IsStale {
		t.Fatal("is_stale = true, want false")
	}
	if claims[0].Note == nil || *claims[0].Note != "journal append work" {
		t.Fatalf("note = %#v, want journal append work", claims[0].Note)
	}
}

func TestClaimListUsesDirectLegacyJSONAndFiltersByStoredRole(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	staleDeveloperID := createClaimListIssue(t, workspace, "Stale developer claim")
	staleQAID := createClaimListIssue(t, workspace, "Stale QA claim")
	liveDeveloperID := createClaimListIssue(t, workspace, "Live developer claim")
	liveHumanID := createClaimListIssue(t, workspace, "Live human claim")

	runFIDCommand(t, workspace, "issue", "claim", staleDeveloperID, "--by", "developer", "--ttl", "1s")
	runFIDCommand(t, workspace, "issue", "claim", staleQAID, "--by", "qa_reviewer", "--ttl", "1s")
	ageCommandClaims(t, workspace, staleDeveloperID, staleQAID)
	runFIDCommand(t, workspace, "issue", "claim", liveDeveloperID, "--by", "developer", "--note", "live developer claim")
	runFIDCommand(t, workspace, "issue", "claim", liveHumanID, "--by", "human", "--note", "human held claim")

	stdout, stderr := runFIDCommand(t, workspace, "claim", "list")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	want := decodeCommandJSONStructure(t, stdout)
	wantRows := want.([]any)
	if got := commandClaimListIDs(t, wantRows); !equalIssueIDs(got, []string{liveDeveloperID, liveHumanID}) {
		t.Fatalf("live claim ids = %#v, want [%q, %q]", got, liveDeveloperID, liveHumanID)
	}
	if got := commandClaimListField(t, wantRows[1], "claimed_by"); got != "human" {
		t.Fatalf("human claimed_by = %q, want human", got)
	}

	humanStdout, humanStderr := runFIDCommand(t, workspace, "claim", "list", "--by", "human")
	if humanStderr != "" {
		t.Fatalf("human stderr = %q, want empty", humanStderr)
	}
	humanRows := decodeCommandJSONStructure(t, humanStdout).([]any)
	if got := commandClaimListIDs(t, humanRows); !equalIssueIDs(got, []string{liveHumanID}) {
		t.Fatalf("human claim ids = %#v, want [%q]", got, liveHumanID)
	}

	unknownStdout, unknownStderr := runFIDCommand(t, workspace, "claim", "list", "--by", "intern")
	if unknownStderr != "" {
		t.Fatalf("unknown-role stderr = %q, want empty", unknownStderr)
	}
	if got := decodeCommandJSONStructure(t, unknownStdout); !commandStructuresEqual([]any{}, got) {
		t.Fatalf("unknown-role output = %#v, want empty array", got)
	}
}

func TestClaimListStaleFilterIsReadOnlyAndUsesRowIDOrder(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	firstStaleID := createClaimListIssue(t, workspace, "First stale claim")
	secondStaleID := createClaimListIssue(t, workspace, "Second stale claim")
	liveID := createClaimListIssue(t, workspace, "Live claim")

	runFIDCommand(t, workspace, "issue", "claim", firstStaleID, "--by", "developer", "--ttl", "1s")
	runFIDCommand(t, workspace, "issue", "claim", secondStaleID, "--by", "qa_reviewer", "--ttl", "1s")
	runFIDCommand(t, workspace, "issue", "claim", liveID, "--by", "human")
	ageCommandClaims(t, workspace, firstStaleID, secondStaleID)

	journalPath := mustJournalPath(t, workspace)
	beforePayload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}
	beforeLines := splitJournalLines(beforePayload)

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	activityBefore := countActivityLogRows(t, conn)

	stdout, stderr := runFIDCommand(t, workspace, "claim", "list", "--stale")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	rows := decodeCommandJSONStructure(t, stdout).([]any)
	if got := commandClaimListIDs(t, rows); !equalIssueIDs(got, []string{firstStaleID, secondStaleID}) {
		t.Fatalf("stale claim ids = %#v, want [%q, %q]", got, firstStaleID, secondStaleID)
	}
	if got := commandClaimListField(t, rows[0], "is_stale"); got != "true" {
		t.Fatalf("first stale is_stale = %q, want true", got)
	}
	if got := commandClaimListField(t, rows[1], "is_stale"); got != "true" {
		t.Fatalf("second stale is_stale = %q, want true", got)
	}

	afterPayload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after list error = %v", journalPath, err)
	}
	afterLines := splitJournalLines(afterPayload)
	if len(afterLines) != len(beforeLines) {
		t.Fatalf("journal lines after claim list = %d, want unchanged %d", len(afterLines), len(beforeLines))
	}

	activityAfter := countActivityLogRows(t, conn)
	if activityAfter != activityBefore {
		t.Fatalf("activity_log rows after claim list = %d, want unchanged %d", activityAfter, activityBefore)
	}
}

func TestClaimListOrderIsStableAcrossRebuild(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	firstID := createClaimListIssue(t, workspace, "First live claim")
	secondID := createClaimListIssue(t, workspace, "Second live claim")
	thirdID := createClaimListIssue(t, workspace, "Third live claim")

	runFIDCommand(t, workspace, "issue", "claim", firstID, "--by", "developer")
	runFIDCommand(t, workspace, "issue", "claim", secondID, "--by", "human")
	runFIDCommand(t, workspace, "issue", "claim", thirdID, "--by", "qa_reviewer")

	beforeStdout, beforeStderr := runFIDCommand(t, workspace, "claim", "list")
	if beforeStderr != "" {
		t.Fatalf("before rebuild stderr = %q, want empty", beforeStderr)
	}
	beforeRows := decodeCommandJSONStructure(t, beforeStdout).([]any)
	if got := commandClaimListIDs(t, beforeRows); !equalIssueIDs(got, []string{firstID, secondID, thirdID}) {
		t.Fatalf("before rebuild ids = %#v, want [%q, %q, %q]", got, firstID, secondID, thirdID)
	}

	rebuildStdout, rebuildStderr := runFIDCommand(t, workspace, "rebuild", "--by", "workflow_runner")
	if !strings.Contains(rebuildStdout, "rebuild complete") {
		t.Fatalf("rebuild stdout = %q, want rebuild confirmation", rebuildStdout)
	}
	if rebuildStderr != "" {
		t.Fatalf("rebuild stderr = %q, want empty", rebuildStderr)
	}

	afterStdout, afterStderr := runFIDCommand(t, workspace, "claim", "list")
	if afterStderr != "" {
		t.Fatalf("after rebuild stderr = %q, want empty", afterStderr)
	}
	afterRows := decodeCommandJSONStructure(t, afterStdout).([]any)
	if got := commandClaimListIDs(t, afterRows); !equalIssueIDs(got, []string{firstID, secondID, thirdID}) {
		t.Fatalf("after rebuild ids = %#v, want [%q, %q, %q]", got, firstID, secondID, thirdID)
	}
}

func createClaimListIssue(t *testing.T, workspace, title string) string {
	t.Helper()

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", title, "--by", "developer")
	return parseCreatedIssueID(t, stderr)
}

func ageCommandClaims(t *testing.T, workspace string, issueIDs ...string) {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	const staleClaimedAt = "2000-01-01T12:00:00Z"
	for _, issueID := range issueIDs {
		if _, err := conn.Exec(`
			UPDATE claims
			SET claimed_at = ?
			WHERE ref_type = 'issue' AND ref_id = ?
		`, staleClaimedAt, issueID); err != nil {
			t.Fatalf("age claim %q: %v", issueID, err)
		}
	}
}

func commandClaimListIDs(t *testing.T, rows []any) []string {
	t.Helper()

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		fields, ok := row.(map[string]any)
		if !ok {
			t.Fatalf("row = %#v, want object", row)
		}
		refID, _ := fields["ref_id"].(string)
		ids = append(ids, refID)
	}
	return ids
}

func commandClaimListField(t *testing.T, row any, field string) string {
	t.Helper()

	fields, ok := row.(map[string]any)
	if !ok {
		t.Fatalf("row = %#v, want object", row)
	}
	return commandClaimListFieldValue(t, fields, field)
}

func commandClaimListFieldValue(t *testing.T, fields map[string]any, field string) string {
	t.Helper()

	value, ok := fields[field]
	if !ok {
		t.Fatalf("field %q missing from %#v", field, fields)
	}

	switch typed := value.(type) {
	case string:
		return typed
	case json.Number:
		return typed.String()
	case bool:
		return strconv.FormatBool(typed)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", typed)
	}
}
