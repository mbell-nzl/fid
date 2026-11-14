//go:build integration

package integration

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestV01ClosureJourneyAndRebuildParity(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	closureRun(t, workspace, "git", "init")
	closureRunFID(t, binary, workspace, "init", "--prefix", "FID")
	assertClosureGitLayout(t, workspace)

	journalPath := closureJournalPath(t, workspace)
	assertClosureJournalNeverTracked(t, workspace)

	closureRunFID(t, binary, workspace, "label", "add", "release", "--description", "Closure journey label", "--by", "developer")
	closureWriteRequirementFixture(t, workspace)
	assertClosureUnknownCommandWithoutJournalAppend(t, binary, workspace, journalPath, []string{"requirement", "sync", "--source-doc", "project-docs/specs/closure.md", "--id", "PRD-LINK-002", "--by", "tl"})

	parentID := closureCreateIssue(t, binary, workspace, "task", "Closure lifecycle parent", nil, nil, nil)
	mainIssueID := closureCreateIssue(t, binary, workspace, "subtask", "Closure lifecycle subtask", nil, &parentID, nil)
	peerIssueID := closureCreateIssue(t, binary, workspace, "task", "Closure linked peer", nil, nil, nil)
	closureRunFID(t, binary, workspace, "issue", "label", "add", mainIssueID, "release", "--by", "developer")

	closureRunFID(t, binary, workspace, "issue", "transition", mainIssueID, "--to", "in_progress", "--actor", "agent:qa_reviewer")
	acID := closureAddAcceptanceCriterion(t, binary, workspace, mainIssueID)
	closureRunFID(t, binary, workspace, "issue", "transition", mainIssueID, "--to", "done", "--actor", "agent:developer")
	assertClosureUnknownCommandWithoutJournalAppend(t, binary, workspace, journalPath,
		[]string{"issue", "ac", "verify"})
	_ = acID
	assertClosureConcernMutationsUnavailable(t, binary, workspace, mainIssueID, "TC-unavailable")
	assertClosureRejectedWithoutJournalAppend(t, binary, workspace, journalPath,
		[]string{"--json", "link", "add", mainIssueID, "implements", "PRD-LINK-002", "--from-type", "issue", "--to-type", "requirement", "--by", "developer"},
		`"code":"invalid-argument"`)
	closureRunFID(t, binary, workspace, "link", "add", peerIssueID, "discovered_from", mainIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
	assertClosureLinkReadback(t, binary, workspace, mainIssueID, peerIssueID)
	assertClosureRejectedWithoutJournalAppend(t, binary, workspace, journalPath,
		[]string{"link", "add", mainIssueID, "implements", peerIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer"},
		"error: link relation is unsupported\nreason: unsupported-relation; allowed=blocks,blocked_by,relates_to,duplicates,duplicated_by,supersedes,superseded_by,caused_by,causes,discovered_from,discovered, argument=relation, value=implements\n")

	seqEpicID := closureCreateIssue(t, binary, workspace, "epic", "Sequential closure epic", stringPointer("sequential"), nil, nil)
	seqCurrentID := closureCreateIssue(t, binary, workspace, "task", "Sequential current", nil, stringPointer(seqEpicID), intPointer(1))
	seqNextID := closureCreateIssue(t, binary, workspace, "task", "Sequential next", nil, stringPointer(seqEpicID), intPointer(2))
	closureRunFID(t, binary, workspace, "issue", "transition", seqCurrentID, "--to", "cancelled", "--reason", "Test cancellation", "--actor", "human:mark")
	assertClosureNextTaskInEpic(t, binary, workspace, seqEpicID, seqNextID)

	priorityEpicID := closureCreateIssue(t, binary, workspace, "epic", "Prioritised closure epic", stringPointer("prioritised"), nil, nil)
	closureCreateIssue(t, binary, workspace, "task", "Lower priority", nil, stringPointer(priorityEpicID), nil, "--priority", "3")
	priorityHighID := closureCreateIssue(t, binary, workspace, "task", "Higher priority", nil, stringPointer(priorityEpicID), nil, "--priority", "0")
	assertClosureNextTaskInEpic(t, binary, workspace, priorityEpicID, priorityHighID)

	assertClosureUnknownCommandWithoutJournalAppend(t, binary, workspace, journalPath,
		[]string{"run", "record", "--issue", mainIssueID, "--completion", "TC-unavailable",
			"--role", "developer", "--model", "claude-sonnet-4-6", "--tokens-in", "1200", "--tokens-out", "3400", "--by", "developer"})

	assertLeanIntegrationJournal(t, workspace)
	journalBeforeNonJournal := closureReadFile(t, journalPath)
	journalFilesBefore := closureJournalFileBytes(t, workspace)
	assertClosureMaintenanceAndReadCommands(t, binary, workspace, mainIssueID, journalBeforeNonJournal)
	maintenanceNow := closureAgedMaintenanceTime(t, workspace)

	dbPath := filepath.Join(config.LedgerDir(workspace), "fid.sqlite")
	if err := fiddb.MigrateDown(dbPath, filepath.Join("..", "..", "..", "migrations")); err != nil {
		t.Fatalf("migrate down for journal immutability boundary: %v", err)
	}
	if err := fiddb.MigrateUp(dbPath, filepath.Join("..", "..", "..", "migrations")); err != nil {
		t.Fatalf("migrate up for journal immutability boundary: %v", err)
	}
	assertClosureJournalFilesEqual(t, closureJournalFileBytes(t, workspace), journalFilesBefore, "migration")

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("open projection for stale-version boundary: %v", err)
	}
	if _, err := conn.Exec(`UPDATE projection_meta SET projection_schema_version = 3 WHERE key = 'default'`); err != nil {
		_ = conn.Close()
		t.Fatalf("set stale projection version: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("close stale projection: %v", err)
	}
	closureRunFID(t, binary, workspace, "--json", "issue", "list")
	assertClosureJournalFilesEqual(t, closureJournalFileBytes(t, workspace), journalFilesBefore, "stale-version refresh")

	beforeTables := closureProjectedTableSnapshots(t, binary, workspace)
	beforeSchema := closureSchemaSnapshot(t, binary, workspace)
	assertClosureForeignKeysClean(t, binary, workspace)
	beforeMaintenance := closureMaintenanceResultAt(t, workspace, maintenanceNow)
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("remove derived DB: %v", err)
	}
	closureRunFID(t, binary, workspace, "rebuild")
	assertClosureJournalBytesEqual(t, journalPath, journalBeforeNonJournal, "rebuild")
	assertClosureJournalFilesEqual(t, closureJournalFileBytes(t, workspace), journalFilesBefore, "delete-and-rebuild")
	assertClosureProjectedTablesEqual(t, closureProjectedTableSnapshots(t, binary, workspace), beforeTables)
	if afterSchema := closureSchemaSnapshot(t, binary, workspace); afterSchema != beforeSchema {
		t.Fatalf("schema after first rebuild differs from before rebuild")
	}
	assertClosureForeignKeysClean(t, binary, workspace)
	afterMaintenance := closureMaintenanceResultAt(t, workspace, maintenanceNow)
	assertClosureMaintenanceStateEqual(t, afterMaintenance, beforeMaintenance)
	assertClosureWorkflowsAfterRebuild(t, binary, workspace)
	assertClosureJournalNeverTracked(t, workspace)

	closureRunFID(t, binary, workspace, "rebuild")
	assertClosureJournalFilesEqual(t, closureJournalFileBytes(t, workspace), journalFilesBefore, "second rebuild")
	assertClosureProjectedTablesEqual(t, closureProjectedTableSnapshots(t, binary, workspace), beforeTables)
	if afterSchema := closureSchemaSnapshot(t, binary, workspace); afterSchema != beforeSchema {
		t.Fatalf("schema after second rebuild differs from before rebuild")
	}
	assertClosureForeignKeysClean(t, binary, workspace)

	if bytes.Equal(journalBeforeNonJournal, nil) {
		t.Fatal("journal snapshot is empty")
	}
}
