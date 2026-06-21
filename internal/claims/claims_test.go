package claims

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestAcquireCreatesLiveIssueClaimWithDefaultTTL(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Claim target")
	journalPath := claimsJournalPathForWorkspace(t, workspace)
	linesBefore := countClaimsJournalLines(t, journalPath)

	result, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Acquire() error = %v", err)
	}

	if result.RefType != "issue" {
		t.Fatalf("Acquire().RefType = %q, want issue", result.RefType)
	}
	if result.RefID != created.ID {
		t.Fatalf("Acquire().RefID = %q, want %q", result.RefID, created.ID)
	}
	if result.ClaimedBy != "developer" {
		t.Fatalf("Acquire().ClaimedBy = %q, want developer", result.ClaimedBy)
	}
	if result.TTLSeconds != 1800 {
		t.Fatalf("Acquire().TTLSeconds = %d, want 1800", result.TTLSeconds)
	}
	if result.Note != nil {
		t.Fatalf("Acquire().Note = %#v, want nil", result.Note)
	}

	claimedAt, err := time.Parse(time.RFC3339, result.ClaimedAt)
	if err != nil {
		t.Fatalf("time.Parse(claimed_at) error = %v", err)
	}
	expiresAt, err := time.Parse(time.RFC3339, result.ExpiresAt)
	if err != nil {
		t.Fatalf("time.Parse(expires_at) error = %v", err)
	}
	if got := expiresAt.Sub(claimedAt); got != 30*time.Minute {
		t.Fatalf("expires_at - claimed_at = %s, want 30m0s", got)
	}

	if got := countClaimsJournalLines(t, journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after acquire = %d, want %d", got, linesBefore+1)
	}

	stored := readClaimStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	if last.EventType != "issue.claim.acquired" {
		t.Fatalf("last event type = %q, want issue.claim.acquired", last.EventType)
	}

	conn := openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeClaimsDB(t, conn)

	var (
		refType     string
		refID       string
		claimedBy   string
		claimedAtDB string
		ttlSeconds  int
		note        sql.NullString
	)
	if err := conn.QueryRow(`
		SELECT ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, note
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
	`, created.ID).Scan(&refType, &refID, &claimedBy, &claimedAtDB, &ttlSeconds, &note); err != nil {
		t.Fatalf("query claims row: %v", err)
	}
	if refType != "issue" || refID != created.ID || claimedBy != "developer" || ttlSeconds != 1800 {
		t.Fatalf("claims row = (%q, %q, %q, %d), want (%q, %q, %q, %d)", refType, refID, claimedBy, ttlSeconds, "issue", created.ID, "developer", 1800)
	}
	if claimedAtDB != result.ClaimedAt {
		t.Fatalf("claims row claimed_at = %q, want %q", claimedAtDB, result.ClaimedAt)
	}
	if note.Valid {
		t.Fatalf("claims row note = %#v, want NULL", note)
	}
}

func TestAcquireSupportsCustomTTLAndNote(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Custom claim target")

	result, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		TTLSeconds:    intPointer(600),
		Note:          stringPointer("Implementing journal append"),
	})
	if err != nil {
		t.Fatalf("Acquire(custom ttl/note) error = %v", err)
	}

	if result.TTLSeconds != 600 {
		t.Fatalf("Acquire(custom ttl/note).TTLSeconds = %d, want 600", result.TTLSeconds)
	}
	if result.Note == nil || *result.Note != "Implementing journal append" {
		t.Fatalf("Acquire(custom ttl/note).Note = %#v, want Implementing journal append", result.Note)
	}

	conn := openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeClaimsDB(t, conn)

	var (
		ttl  int
		note sql.NullString
	)
	if err := conn.QueryRow(`
		SELECT ttl_seconds, note
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
	`, created.ID).Scan(&ttl, &note); err != nil {
		t.Fatalf("query custom claim row: %v", err)
	}
	if ttl != 600 {
		t.Fatalf("claims ttl_seconds = %d, want 600", ttl)
	}
	if !note.Valid || note.String != "Implementing journal append" {
		t.Fatalf("claims note = %#v, want Implementing journal append", note)
	}
}

func TestAcquireRejectsLiveConflictBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Claim conflict target")

	initial, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Acquire(initial) error = %v", err)
	}

	journalPath := claimsJournalPathForWorkspace(t, workspace)
	linesBefore := countClaimsJournalLines(t, journalPath)

	_, err = Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "qa_reviewer",
	})
	if err == nil {
		t.Fatal("Acquire(conflict) error = nil, want live conflict")
	}

	var conflictErr *LiveClaimConflictError
	if !errors.As(err, &conflictErr) {
		t.Fatalf("Acquire(conflict) error = %T, want *LiveClaimConflictError", err)
	}
	if conflictErr.RefID != created.ID {
		t.Fatalf("LiveClaimConflictError.RefID = %q, want %q", conflictErr.RefID, created.ID)
	}
	if conflictErr.ClaimedBy != "developer" {
		t.Fatalf("LiveClaimConflictError.ClaimedBy = %q, want developer", conflictErr.ClaimedBy)
	}
	if conflictErr.ExpiresAt != initial.ExpiresAt {
		t.Fatalf("LiveClaimConflictError.ExpiresAt = %q, want %q", conflictErr.ExpiresAt, initial.ExpiresAt)
	}
	for _, want := range []string{created.ID, "developer", initial.ExpiresAt} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Acquire(conflict) error = %q, want substring %q", err.Error(), want)
		}
	}

	if got := countClaimsJournalLines(t, journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected acquire = %d, want unchanged %d", got, linesBefore)
	}

	conn := openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeClaimsDB(t, conn)

	var count int
	if err := conn.QueryRow(`
		SELECT COUNT(*)
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
	`, created.ID).Scan(&count); err != nil {
		t.Fatalf("count claims rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("claims row count = %d, want 1", count)
	}
}

func TestAcquireTakesOverStaleClaimWithoutForcePermission(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Stale claim target")

	_, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Acquire(initial) error = %v", err)
	}

	staleClaimedAt := time.Now().UTC().Add(-2 * time.Hour).Format(time.RFC3339)
	staleClaimedAtTime, err := time.Parse(time.RFC3339, staleClaimedAt)
	if err != nil {
		t.Fatalf("time.Parse(stale claimed_at) error = %v", err)
	}
	staleExpiresAt := staleClaimedAtTime.Add(30 * time.Minute).Format(time.RFC3339)
	conn := openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	if _, err := conn.Exec(`
		UPDATE claims
		SET claimed_at = ?
		WHERE ref_type = 'issue' AND ref_id = ?
	`, staleClaimedAt, created.ID); err != nil {
		t.Fatalf("expire claim row: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	journalPath := claimsJournalPathForWorkspace(t, workspace)
	linesBefore := countClaimsJournalLines(t, journalPath)

	result, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "qa_reviewer",
	})
	if err != nil {
		t.Fatalf("Acquire(stale takeover) error = %v", err)
	}
	if result.ClaimedBy != "qa_reviewer" {
		t.Fatalf("Acquire(stale takeover).ClaimedBy = %q, want qa_reviewer", result.ClaimedBy)
	}

	if got := countClaimsJournalLines(t, journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after stale takeover = %d, want %d", got, linesBefore+1)
	}

	stored := readClaimStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	var (
		updateEffect   *journal.Effect
		activityEffect *journal.Effect
	)
	for index := range last.Effects {
		effect := &last.Effects[index]
		switch {
		case effect.Table == "claims" && effect.Op == "update":
			updateEffect = effect
		case effect.Table == "activity_log" && effect.Op == "activity":
			activityEffect = effect
		}
	}
	if updateEffect == nil {
		t.Fatal("missing claims update effect")
	}
	if activityEffect == nil {
		t.Fatal("missing activity takeover effect")
	}
	if got := updateEffect.Preconditions["claimed_by"]; got != "developer" {
		t.Fatalf("takeover precondition claimed_by = %#v, want developer", got)
	}
	if got := updateEffect.Values["claimed_by"]; got != "qa_reviewer" {
		t.Fatalf("takeover value claimed_by = %#v, want qa_reviewer", got)
	}
	if got := activityEffect.Values["event"]; got != "claim_superseded" {
		t.Fatalf("activity event = %#v, want claim_superseded", got)
	}
	details, ok := activityEffect.Values["details"].(map[string]any)
	if !ok {
		t.Fatalf("activity details = %#v, want map", activityEffect.Values["details"])
	}
	if got := details["previous_claimed_by"]; got != "developer" {
		t.Fatalf("activity details previous_claimed_by = %#v, want developer", got)
	}
	if got := details["previous_expires_at"]; got != staleExpiresAt {
		t.Fatalf("activity details previous_expires_at = %#v, want %q", got, staleExpiresAt)
	}

	conn = openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeClaimsDB(t, conn)

	var claimedBy string
	if err := conn.QueryRow(`
		SELECT claimed_by
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
	`, created.ID).Scan(&claimedBy); err != nil {
		t.Fatalf("query current claim holder: %v", err)
	}
	if claimedBy != "qa_reviewer" {
		t.Fatalf("current claim holder = %q, want qa_reviewer", claimedBy)
	}
}

func TestAcquireForceSupersedesLiveClaimForPermittedRole(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Force claim target")

	initial, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Acquire(initial) error = %v", err)
	}

	journalPath := claimsJournalPathForWorkspace(t, workspace)
	linesBefore := countClaimsJournalLines(t, journalPath)

	result, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "tl",
		Force:         true,
	})
	if err != nil {
		t.Fatalf("Acquire(force) error = %v", err)
	}
	if result.ClaimedBy != "tl" {
		t.Fatalf("Acquire(force).ClaimedBy = %q, want tl", result.ClaimedBy)
	}

	if got := countClaimsJournalLines(t, journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after force acquire = %d, want %d", got, linesBefore+1)
	}

	stored := readClaimStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	var activityEffect *journal.Effect
	for index := range last.Effects {
		effect := &last.Effects[index]
		if effect.Table == "activity_log" && effect.Op == "activity" {
			activityEffect = effect
		}
	}
	if activityEffect == nil {
		t.Fatal("missing force activity effect")
	}
	details, ok := activityEffect.Values["details"].(map[string]any)
	if !ok {
		t.Fatalf("activity details = %#v, want map", activityEffect.Values["details"])
	}
	if got := details["previous_claimed_by"]; got != "developer" {
		t.Fatalf("activity details previous_claimed_by = %#v, want developer", got)
	}
	if got := details["previous_expires_at"]; got != initial.ExpiresAt {
		t.Fatalf("activity details previous_expires_at = %#v, want %q", got, initial.ExpiresAt)
	}
}

func TestAcquireForceRejectsDeniedRoleBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Denied force claim target")

	_, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "qa_reviewer",
	})
	if err != nil {
		t.Fatalf("Acquire(initial) error = %v", err)
	}

	journalPath := claimsJournalPathForWorkspace(t, workspace)
	linesBefore := countClaimsJournalLines(t, journalPath)

	_, err = Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Force:         true,
	})
	if err == nil {
		t.Fatal("Acquire(force denied) error = nil, want permission denial")
	}

	var deniedErr *policy.PermissionDeniedError
	if !errors.As(err, &deniedErr) {
		t.Fatalf("Acquire(force denied) error = %T, want *policy.PermissionDeniedError", err)
	}
	if deniedErr.Role != "developer" {
		t.Fatalf("PermissionDeniedError.Role = %q, want developer", deniedErr.Role)
	}
	if deniedErr.Action != "issue.claim.force_acquire" {
		t.Fatalf("PermissionDeniedError.Action = %q, want issue.claim.force_acquire", deniedErr.Action)
	}

	if got := countClaimsJournalLines(t, journalPath); got != linesBefore {
		t.Fatalf("journal lines after denied force = %d, want unchanged %d", got, linesBefore)
	}
}

func TestRenewUpdatesClaimTimePreservesOrOverridesTTLAndProjectsActivity(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Renew claim target")

	initial, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		TTLSeconds:    intPointer(900),
		Note:          stringPointer("focus mode"),
	})
	if err != nil {
		t.Fatalf("Acquire(initial) error = %v", err)
	}

	journalPath := claimsJournalPathForWorkspace(t, workspace)
	linesBeforeFirstRenew := countClaimsJournalLines(t, journalPath)

	firstRenewed, err := Renew(RenewOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Renew(preserve ttl) error = %v", err)
	}
	if firstRenewed.TTLSeconds != 900 {
		t.Fatalf("Renew(preserve ttl).TTLSeconds = %d, want 900", firstRenewed.TTLSeconds)
	}
	if firstRenewed.Note == nil || *firstRenewed.Note != "focus mode" {
		t.Fatalf("Renew(preserve ttl).Note = %#v, want focus mode", firstRenewed.Note)
	}
	if firstRenewed.ClaimedAt == initial.ClaimedAt {
		t.Fatalf("Renew(preserve ttl).ClaimedAt = %q, want a later timestamp than %q", firstRenewed.ClaimedAt, initial.ClaimedAt)
	}
	if got := countClaimsJournalLines(t, journalPath); got != linesBeforeFirstRenew+1 {
		t.Fatalf("journal lines after first renew = %d, want %d", got, linesBeforeFirstRenew+1)
	}

	secondRenewed, err := Renew(RenewOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		TTLSeconds:    intPointer(1200),
	})
	if err != nil {
		t.Fatalf("Renew(override ttl) error = %v", err)
	}
	if secondRenewed.TTLSeconds != 1200 {
		t.Fatalf("Renew(override ttl).TTLSeconds = %d, want 1200", secondRenewed.TTLSeconds)
	}

	stored := readClaimStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	if last.EventType != "issue.claim.renewed" {
		t.Fatalf("last event type = %q, want issue.claim.renewed", last.EventType)
	}
	var (
		updateEffect   *journal.Effect
		activityEffect *journal.Effect
	)
	for index := range last.Effects {
		effect := &last.Effects[index]
		switch {
		case effect.Table == "claims" && effect.Op == "update":
			updateEffect = effect
		case effect.Table == "activity_log" && effect.Op == "activity":
			activityEffect = effect
		}
	}
	if updateEffect == nil {
		t.Fatal("missing renew claims update effect")
	}
	if activityEffect == nil {
		t.Fatal("missing renew activity effect")
	}
	if got := updateEffect.Preconditions["claimed_by"]; got != "developer" {
		t.Fatalf("renew precondition claimed_by = %#v, want developer", got)
	}
	if got := updateEffect.Values["ttl_seconds"]; got != float64(1200) {
		t.Fatalf("renew value ttl_seconds = %#v, want 1200", got)
	}
	if got := activityEffect.Values["event"]; got != "claim_renewed" {
		t.Fatalf("activity event = %#v, want claim_renewed", got)
	}

	conn := openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeClaimsDB(t, conn)

	var (
		ttl  int
		note sql.NullString
	)
	if err := conn.QueryRow(`
		SELECT ttl_seconds, note
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
	`, created.ID).Scan(&ttl, &note); err != nil {
		t.Fatalf("query renewed claim row: %v", err)
	}
	if ttl != 1200 {
		t.Fatalf("renewed claim ttl_seconds = %d, want 1200", ttl)
	}
	if !note.Valid || note.String != "focus mode" {
		t.Fatalf("renewed claim note = %#v, want focus mode", note)
	}
}

func TestRenewRejectsNonHolderBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Renew wrong holder target")

	_, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Acquire(initial) error = %v", err)
	}

	journalPath := claimsJournalPathForWorkspace(t, workspace)
	linesBefore := countClaimsJournalLines(t, journalPath)

	_, err = Renew(RenewOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "qa_reviewer",
	})
	if err == nil {
		t.Fatal("Renew(non-holder) error = nil, want rejection")
	}

	var notHeldErr *ClaimNotHeldError
	if !errors.As(err, &notHeldErr) {
		t.Fatalf("Renew(non-holder) error = %T, want *ClaimNotHeldError", err)
	}
	if notHeldErr.RefID != created.ID {
		t.Fatalf("ClaimNotHeldError.RefID = %q, want %q", notHeldErr.RefID, created.ID)
	}
	if notHeldErr.ClaimedBy != "developer" {
		t.Fatalf("ClaimNotHeldError.ClaimedBy = %q, want developer", notHeldErr.ClaimedBy)
	}

	if got := countClaimsJournalLines(t, journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected renew = %d, want unchanged %d", got, linesBefore)
	}
}

func TestAcquireIsAtomicAcrossConcurrentAttempts(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Concurrent claim target")

	type acquireResult struct {
		role  string
		claim Claim
		err   error
	}

	start := make(chan struct{})
	results := make(chan acquireResult, 2)
	var wg sync.WaitGroup

	for _, role := range []string{"developer", "qa_reviewer"} {
		wg.Add(1)
		go func(role string) {
			defer wg.Done()
			<-start
			claim, err := Acquire(AcquireOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testClaimsMigrationsDir(t),
				IssueID:       created.ID,
				ActorRole:     role,
			})
			results <- acquireResult{role: role, claim: claim, err: err}
		}(role)
	}

	close(start)
	wg.Wait()
	close(results)

	var (
		successes []acquireResult
		conflicts []acquireResult
	)
	for result := range results {
		if result.err == nil {
			successes = append(successes, result)
			continue
		}
		conflicts = append(conflicts, result)
	}

	if len(successes) != 1 {
		t.Fatalf("successful acquires = %d, want 1", len(successes))
	}
	if len(conflicts) != 1 {
		t.Fatalf("conflicting acquires = %d, want 1", len(conflicts))
	}

	winner := successes[0]
	loser := conflicts[0]

	var conflictErr *LiveClaimConflictError
	if !errors.As(loser.err, &conflictErr) {
		t.Fatalf("loser error = %T, want *LiveClaimConflictError", loser.err)
	}
	if conflictErr.ClaimedBy != winner.claim.ClaimedBy {
		t.Fatalf("loser conflict claimed_by = %q, want %q", conflictErr.ClaimedBy, winner.claim.ClaimedBy)
	}
	if conflictErr.ExpiresAt != winner.claim.ExpiresAt {
		t.Fatalf("loser conflict expires_at = %q, want %q", conflictErr.ExpiresAt, winner.claim.ExpiresAt)
	}

	conn := openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeClaimsDB(t, conn)

	var (
		count     int
		claimedBy string
	)
	if err := conn.QueryRow(`
		SELECT COUNT(*), MIN(claimed_by)
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
	`, created.ID).Scan(&count, &claimedBy); err != nil {
		t.Fatalf("query concurrent claims state: %v", err)
	}
	if count != 1 {
		t.Fatalf("current claims row count = %d, want 1", count)
	}
	if claimedBy != winner.claim.ClaimedBy {
		t.Fatalf("current claim holder = %q, want %q", claimedBy, winner.claim.ClaimedBy)
	}
}

func testClaimsMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func newClaimsWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}

	return workspace
}

func createClaimIssue(t *testing.T, workspace, title string) issues.CreateResult {
	t.Helper()

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		Kind:          "task",
		Title:         title,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	return created
}

func claimsJournalPathForWorkspace(t *testing.T, workspace string) string {
	t.Helper()

	journalDir := filepath.Join(workspace, ".fid", "journal")
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", journalDir, err)
	}
	if len(entries) != 1 {
		t.Fatalf("journal entries = %d, want 1", len(entries))
	}

	return filepath.Join(journalDir, entries[0].Name())
}

func countClaimsJournalLines(t *testing.T, path string) int {
	t.Helper()

	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	trimmed := strings.TrimSpace(string(payload))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func readClaimStoredEvents(t *testing.T, path string) []journal.Event {
	t.Helper()

	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}

	lines := strings.Split(strings.TrimSpace(string(payload)), "\n")
	events := make([]journal.Event, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var event journal.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		events = append(events, event)
	}
	return events
}

func openClaimsDB(t *testing.T, path string) *sql.DB {
	t.Helper()

	conn, err := fiddb.Open(path)
	if err != nil {
		t.Fatalf("Open(%q) error = %v", path, err)
	}
	return conn
}

func closeClaimsDB(t *testing.T, conn interface{ Close() error }) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func stringPointer(value string) *string {
	return &value
}

func intPointer(value int) *int {
	return &value
}
