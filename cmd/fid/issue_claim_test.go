package main

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestIssueClaimCommandCreatesLiveClaimWithDefaultTTL(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI claim", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	stdout, stderr := runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer")
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "claimed "+issueID+" by developer (expires ") {
		t.Fatalf("stderr = %q, want claim confirmation", stderr)
	}

	claimsStdout, claimsStderr := runFIDCommand(t, workspace, "claim", "list", "--by", "developer")
	if claimsStderr != "" {
		t.Fatalf("claims stderr = %q, want empty", claimsStderr)
	}

	var claims []struct {
		ClaimedBy  string  `json:"claimed_by"`
		TTLSeconds int     `json:"ttl_seconds"`
		ExpiresAt  *string `json:"expires_at"`
		Note       *string `json:"note"`
	}
	if err := json.Unmarshal([]byte(claimsStdout), &claims); err != nil {
		t.Fatalf("Unmarshal(claims) error = %v", err)
	}
	if len(claims) != 1 {
		t.Fatalf("claims len = %d, want 1", len(claims))
	}
	if claims[0].ClaimedBy != "developer" {
		t.Fatalf("claimed_by = %q, want developer", claims[0].ClaimedBy)
	}
	if claims[0].TTLSeconds != 1800 {
		t.Fatalf("ttl_seconds = %d, want 1800", claims[0].TTLSeconds)
	}
	if claims[0].ExpiresAt == nil || *claims[0].ExpiresAt == "" {
		t.Fatalf("expires_at = %#v, want populated value", claims[0].ExpiresAt)
	}
	if claims[0].Note != nil {
		t.Fatalf("note = %#v, want nil", claims[0].Note)
	}
}

func TestIssueClaimCommandRenewsHolderClaimAndPreservesTTL(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI claim renew", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	_, claimStderr := runFIDCommand(
		t,
		workspace,
		"issue", "claim", issueID,
		"--by", "developer",
		"--ttl", "15m",
		"--note", "focus mode",
	)
	if !strings.Contains(claimStderr, "claimed "+issueID+" by developer (expires ") {
		t.Fatalf("claim stderr = %q, want initial claim confirmation", claimStderr)
	}

	beforeRenew := mustReadIssueClaimsSection(t, workspace, issueID)

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	activityBefore := countActivityLogRows(t, conn)

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "claim", issueID, "--by", "developer", "--renew")
	if err != nil {
		t.Fatalf("Execute(renew) error = %v, stderr = %q", err, stderr)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "renewed claim on "+issueID+" by developer (expires ") {
		t.Fatalf("stderr = %q, want renew confirmation", stderr)
	}

	afterRenew := mustReadIssueClaimsSection(t, workspace, issueID)
	if len(afterRenew) != 1 {
		t.Fatalf("claims len after renew = %d, want 1", len(afterRenew))
	}
	if afterRenew[0].ClaimedBy != "developer" {
		t.Fatalf("claimed_by after renew = %q, want developer", afterRenew[0].ClaimedBy)
	}
	if afterRenew[0].TTLSeconds != 900 {
		t.Fatalf("ttl_seconds after renew = %d, want 900", afterRenew[0].TTLSeconds)
	}
	if afterRenew[0].Note == nil || *afterRenew[0].Note != "focus mode" {
		t.Fatalf("note after renew = %#v, want focus mode", afterRenew[0].Note)
	}
	if afterRenew[0].ClaimedAt == beforeRenew[0].ClaimedAt {
		t.Fatalf("claimed_at after renew = %q, want change from %q", afterRenew[0].ClaimedAt, beforeRenew[0].ClaimedAt)
	}
	if afterRenew[0].ExpiresAt == nil || beforeRenew[0].ExpiresAt == nil || *afterRenew[0].ExpiresAt == *beforeRenew[0].ExpiresAt {
		t.Fatalf("expires_at after renew = %#v, want change from %#v", afterRenew[0].ExpiresAt, beforeRenew[0].ExpiresAt)
	}

	activityAfter := countActivityLogRows(t, conn)
	if activityAfter != activityBefore+1 {
		t.Fatalf("activity_log rows after renew = %d, want %d", activityAfter, activityBefore+1)
	}
}

func TestIssueUnclaimCommandReleasesHolderClaim(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI unclaim", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer")

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr := runFIDCommand(t, workspace, "issue", "unclaim", issueID, "--by", "developer")
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "released claim on "+issueID+" by developer") {
		t.Fatalf("stderr = %q, want release confirmation", stderr)
	}

	if got := len(mustReadIssueClaimsSection(t, workspace, issueID)); got != 0 {
		t.Fatalf("claims len after unclaim = %d, want 0", got)
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore+1 {
		t.Fatalf("journal lines after unclaim = %d, want %d", linesAfter, linesBefore+1)
	}
}

func TestIssueClaimCommandForceAcquireReportsPreviousHolder(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI force claim", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer")
	beforeForce := mustReadIssueClaimsSection(t, workspace, issueID)
	if len(beforeForce) != 1 || beforeForce[0].ExpiresAt == nil {
		t.Fatalf("claims before force = %#v, want one claim with expires_at", beforeForce)
	}

	stdout, stderr := runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "tl", "--force")
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	want := "force-acquired claim on " + issueID + " by tl (previously held by developer until " + *beforeForce[0].ExpiresAt + ")"
	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr = %q, want force confirmation %q", stderr, want)
	}

	afterForce := mustReadIssueClaimsSection(t, workspace, issueID)
	if len(afterForce) != 1 {
		t.Fatalf("claims len after force = %d, want 1", len(afterForce))
	}
	if afterForce[0].ClaimedBy != "tl" {
		t.Fatalf("claimed_by after force = %q, want tl", afterForce[0].ClaimedBy)
	}
}

func TestIssueClaimCommandSupportsCustomTTLNoteAndHumanReadBack(t *testing.T) {
	t.Parallel()
	t.Run("custom ttl and note", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI custom claim", "--actor", "agent:developer")
		issueID := parseCreatedIssueID(t, createStderr)

		_, stderr := runFIDCommand(
			t,
			workspace,
			"issue", "claim", issueID,
			"--by", "developer",
			"--ttl", "10m",
			"--note", "journal append work",
		)
		if !strings.Contains(stderr, "claimed "+issueID+" by developer (expires ") {
			t.Fatalf("stderr = %q, want claim confirmation", stderr)
		}

		claims := mustReadIssueClaimsSection(t, workspace, issueID)
		if len(claims) != 1 {
			t.Fatalf("claims len = %d, want 1", len(claims))
		}
		if claims[0].TTLSeconds != 600 {
			t.Fatalf("ttl_seconds = %d, want 600", claims[0].TTLSeconds)
		}
		if claims[0].Note == nil || *claims[0].Note != "journal append work" {
			t.Fatalf("note = %#v, want journal append work", claims[0].Note)
		}
	})

	t.Run("human role", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI human claim", "--actor", "agent:developer")
		issueID := parseCreatedIssueID(t, createStderr)

		_, stderr := runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "human")
		if !strings.Contains(stderr, "claimed "+issueID+" by human (expires ") {
			t.Fatalf("stderr = %q, want human claim confirmation", stderr)
		}

		claims := mustReadIssueClaimsSection(t, workspace, issueID)
		if len(claims) != 1 {
			t.Fatalf("claims len = %d, want 1", len(claims))
		}
		if claims[0].ClaimedBy != "human" {
			t.Fatalf("claimed_by = %q, want human", claims[0].ClaimedBy)
		}
	})
}

func TestIssueClaimCommandRejectsInvalidTTLBeforeAppend(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI invalid ttl", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "claim", issueID, "--by", "developer", "--ttl", "tomorrow")
	if err == nil {
		t.Fatal("Execute(invalid ttl) error = nil, want rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, `error: parse --ttl:`) {
		t.Fatalf("stderr = %q, want ttl parse error", stderr)
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore {
		t.Fatalf("journal lines after invalid ttl = %d, want %d", linesAfter, linesBefore)
	}
}

func TestIssueClaimCommandRejectsLiveConflictBeforeAppend(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI conflict claim", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer")
	initial := mustReadIssueClaimsSection(t, workspace, issueID)
	if len(initial) != 1 || initial[0].ExpiresAt == nil {
		t.Fatalf("initial claims = %#v, want one live claim", initial)
	}

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "claim", issueID, "--by", "qa_reviewer")
	if err == nil {
		t.Fatal("Execute(conflict) error = nil, want rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	for _, want := range []string{"error: issue", issueID, "developer", *initial[0].ExpiresAt} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr = %q, want substring %q", stderr, want)
		}
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore {
		t.Fatalf("journal lines after conflict = %d, want %d", linesAfter, linesBefore)
	}
}

func TestIssueClaimCommandRejectsDeniedForceBeforeAppend(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI denied force", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "qa_reviewer")

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "claim", issueID, "--by", "developer", "--force")
	if err == nil {
		t.Fatal("Execute(denied force) error = nil, want rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		`error: role "developer" not permitted for action "issue.claim.force_acquire"`,
		"check `fid policy show` for permitted role/action combinations",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr = %q, want substring %q", stderr, want)
		}
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore {
		t.Fatalf("journal lines after denied force = %d, want %d", linesAfter, linesBefore)
	}
}

func TestIssueClaimCommandRenewsWithOverrideTTLAndRejectsNonHolder(t *testing.T) {
	t.Parallel()
	t.Run("override ttl", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI renew ttl override", "--actor", "agent:developer")
		issueID := parseCreatedIssueID(t, createStderr)

		runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer", "--ttl", "15m", "--note", "focus mode")

		_, stderr := runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer", "--renew", "--ttl", "20m")
		if !strings.Contains(stderr, "renewed claim on "+issueID+" by developer (expires ") {
			t.Fatalf("stderr = %q, want renew confirmation", stderr)
		}

		claims := mustReadIssueClaimsSection(t, workspace, issueID)
		if len(claims) != 1 {
			t.Fatalf("claims len = %d, want 1", len(claims))
		}
		if claims[0].TTLSeconds != 1200 {
			t.Fatalf("ttl_seconds = %d, want 1200", claims[0].TTLSeconds)
		}
		if claims[0].Note == nil || *claims[0].Note != "focus mode" {
			t.Fatalf("note = %#v, want focus mode", claims[0].Note)
		}
	})

	t.Run("non-holder", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI renew wrong holder", "--actor", "agent:developer")
		issueID := parseCreatedIssueID(t, createStderr)

		runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer")

		journalPath := mustJournalPath(t, workspace)
		linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "claim", issueID, "--by", "qa_reviewer", "--renew")
		if err == nil {
			t.Fatal("Execute(non-holder renew) error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		for _, want := range []string{"error: issue", issueID, "qa_reviewer", "developer"} {
			if !strings.Contains(stderr, want) {
				t.Fatalf("stderr = %q, want substring %q", stderr, want)
			}
		}

		linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
		if linesAfter != linesBefore {
			t.Fatalf("journal lines after rejected renew = %d, want %d", linesAfter, linesBefore)
		}
	})
}

func TestIssueClaimCommandRejectsForceAndRenewTogetherBeforeAppend(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI force renew conflict", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "claim", issueID, "--by", "tl", "--force", "--renew")
	if err == nil {
		t.Fatal("Execute(force+renew) error = nil, want rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "error: --force and --renew are mutually exclusive") {
		t.Fatalf("stderr = %q, want mutually exclusive error", stderr)
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore {
		t.Fatalf("journal lines after force+renew rejection = %d, want %d", linesAfter, linesBefore)
	}
}

func TestIssueClaimCommandTakesOverStaleClaimWithoutForce(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI stale claim", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer", "--ttl", "1s")
	ageCommandClaims(t, workspace, issueID)

	stdout, stderr := runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "qa_reviewer")
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "claimed "+issueID+" by qa_reviewer (expires ") {
		t.Fatalf("stderr = %q, want stale-takeover confirmation", stderr)
	}

	claims := mustReadIssueClaimsSection(t, workspace, issueID)
	if len(claims) != 1 {
		t.Fatalf("claims len = %d, want 1", len(claims))
	}
	if claims[0].ClaimedBy != "qa_reviewer" {
		t.Fatalf("claimed_by = %q, want qa_reviewer", claims[0].ClaimedBy)
	}
}

func TestIssueClaimCommandMatchesPolicyShowForForceAcquire(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	rows := legacyPolicyRows(t)
	if got := commandPolicyCellAllowed(t, rows, "issue.claim.force_acquire", "developer"); got {
		t.Fatal("policy show developer/issue.claim.force_acquire = true, want false")
	}
	if got := commandPolicyCellAllowed(t, rows, "issue.claim.force_acquire", "tl"); !got {
		t.Fatal("policy show tl/issue.claim.force_acquire = false, want true")
	}

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI policy parity claim", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)
	runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer")

	deniedStdout, deniedStderr, err := runFIDCommandExpectError(t, workspace, "issue", "claim", issueID, "--by", "developer", "--force")
	if err == nil {
		t.Fatal("developer force-acquire error = nil, want permission denial")
	}
	if deniedStdout != "" {
		t.Fatalf("developer force-acquire stdout = %q, want empty", deniedStdout)
	}
	if !strings.Contains(deniedStderr, `error: role "developer" not permitted for action "issue.claim.force_acquire"`) {
		t.Fatalf("developer force-acquire stderr = %q, want permission denial", deniedStderr)
	}

	allowedStdout, allowedStderr := runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "tl", "--force")
	if allowedStdout != "" {
		t.Fatalf("tl force-acquire stdout = %q, want empty", allowedStdout)
	}
	if !strings.Contains(allowedStderr, "force-acquired claim on "+issueID+" by tl") {
		t.Fatalf("tl force-acquire stderr = %q, want force-acquire confirmation", allowedStderr)
	}
}

func TestIssueUnclaimCommandRejectsWrongHolderRepeatAndRoleErrorsBeforeAppend(t *testing.T) {
	t.Parallel()
	t.Run("wrong holder", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI wrong holder unclaim", "--actor", "agent:developer")
		issueID := parseCreatedIssueID(t, createStderr)

		runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer")

		journalPath := mustJournalPath(t, workspace)
		linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "unclaim", issueID, "--by", "qa_reviewer")
		if err == nil {
			t.Fatal("Execute(wrong holder unclaim) error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		for _, want := range []string{"error: issue", issueID, "qa_reviewer", "developer"} {
			if !strings.Contains(stderr, want) {
				t.Fatalf("stderr = %q, want substring %q", stderr, want)
			}
		}

		linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
		if linesAfter != linesBefore {
			t.Fatalf("journal lines after wrong-holder unclaim = %d, want %d", linesAfter, linesBefore)
		}
	})

	t.Run("repeat release", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI repeat unclaim", "--actor", "agent:developer")
		issueID := parseCreatedIssueID(t, createStderr)

		runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer")
		runFIDCommand(t, workspace, "issue", "unclaim", issueID, "--by", "developer")

		journalPath := mustJournalPath(t, workspace)
		linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "unclaim", issueID, "--by", "developer")
		if err == nil {
			t.Fatal("Execute(repeat unclaim) error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: issue "`+issueID+`" has no live claim held by "developer"`) {
			t.Fatalf("stderr = %q, want repeat-release error", stderr)
		}

		linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
		if linesAfter != linesBefore {
			t.Fatalf("journal lines after repeat unclaim = %d, want %d", linesAfter, linesBefore)
		}
	})

	t.Run("missing role", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI missing unclaim role", "--actor", "agent:developer")
		issueID := parseCreatedIssueID(t, createStderr)

		journalPath := mustJournalPath(t, workspace)
		linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "unclaim", issueID)
		if err == nil {
			t.Fatal("Execute(missing role) error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, "error: --by is required") {
			t.Fatalf("stderr = %q, want missing role error", stderr)
		}

		linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
		if linesAfter != linesBefore {
			t.Fatalf("journal lines after missing-role unclaim = %d, want %d", linesAfter, linesBefore)
		}
	})

	t.Run("unknown role", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI unknown unclaim role", "--actor", "agent:developer")
		issueID := parseCreatedIssueID(t, createStderr)

		journalPath := mustJournalPath(t, workspace)
		linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "unclaim", issueID, "--by", "intern")
		if err == nil {
			t.Fatal("Execute(unknown role) error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		for _, want := range []string{
			`error: role "intern" not in project profile`,
			"known roles:",
			"developer",
			"qa_reviewer",
			"qa_lead",
			"tl_lite",
			"tl",
			"workflow_runner",
			"human",
		} {
			if !strings.Contains(stderr, want) {
				t.Fatalf("stderr = %q, want substring %q", stderr, want)
			}
		}

		linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
		if linesAfter != linesBefore {
			t.Fatalf("journal lines after unknown-role unclaim = %d, want %d", linesAfter, linesBefore)
		}
	})
}

type issueClaimSectionEntry struct {
	ClaimedBy  string  `json:"claimed_by"`
	ClaimedAt  string  `json:"claimed_at"`
	TTLSeconds int     `json:"ttl_seconds"`
	ExpiresAt  *string `json:"expires_at"`
	Note       *string `json:"note"`
}

func mustReadIssueClaimsSection(t *testing.T, workspace, issueID string) []issueClaimSectionEntry {
	t.Helper()

	claimsStdout, claimsStderr := runFIDCommand(t, workspace, "claim", "list")
	if claimsStderr != "" {
		t.Fatalf("claims stderr = %q, want empty", claimsStderr)
	}

	var claims []issueClaimSectionEntry
	if err := json.Unmarshal([]byte(claimsStdout), &claims); err != nil {
		t.Fatalf("Unmarshal(claims) error = %v", err)
	}
	return claims
}
