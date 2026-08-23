//go:build integration

package integration

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestV01ClosureJourneyAndRebuildParity(t *testing.T) {
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	closureRun(t, workspace, "git", "init")
	closureRunFID(t, binary, workspace, "init", "--prefix", "FID")
	assertClosureGitLayout(t, workspace)

	journalPath := closureJournalPath(t, workspace)
	closureRun(t, workspace, "git", "add", filepath.ToSlash(filepath.Join(".fid", "journal")))
	assertClosureJournalTracked(t, workspace)

	closureRunFID(t, binary, workspace, "label", "add", "release", "--description", "Closure journey label", "--by", "developer")
	closureWriteRequirementFixture(t, workspace)
	closureRunFID(t, binary, workspace, "requirement", "sync", "--source-doc", "project-docs/specs/closure.md", "--id", "PRD-LINK-002", "--by", "tl")

	mainIssueID := closureCreateIssue(t, binary, workspace, "task", "Closure lifecycle task", nil, nil, nil)
	peerIssueID := closureCreateIssue(t, binary, workspace, "task", "Closure linked peer", nil, nil, nil)
	closureRunFID(t, binary, workspace, "issue", "label", "add", mainIssueID, "release", "--by", "developer")

	assertClosureRejectedWithoutJournalAppend(t, binary, workspace, journalPath,
		[]string{"issue", "transition", mainIssueID, "--to", "in_progress", "--actor", "agent:qa_reviewer"},
		`role "qa_reviewer" not permitted`)
	assertClosureRejectedWithoutJournalAppend(t, binary, workspace, journalPath,
		[]string{"issue", "transition", mainIssueID, "--to", "in_progress", "--actor", "agent:developer"},
		`has no live claim held by "developer"`)

	closureRunFID(t, binary, workspace, "issue", "claim", mainIssueID, "--by", "developer")
	closureRunFID(t, binary, workspace, "issue", "transition", mainIssueID, "--to", "in_progress", "--actor", "agent:developer")
	acID := closureAddAcceptanceCriterion(t, binary, workspace, mainIssueID)
	assertClosureRejectedWithoutJournalAppend(t, binary, workspace, journalPath,
		[]string{"issue", "transition", mainIssueID, "--to", "dev_complete", "--actor", "agent:developer"},
		"has_completion_report")
	closureRunFID(t, binary, workspace, "issue", "ac", "verify", mainIssueID, acID, "--status", "pass", "--by", "qa_lead")
	firstCompletionID := closureCompleteIssue(t, binary, workspace, mainIssueID)
	concernID, addressCompletionID := closureRaiseAndAddressConcern(t, binary, workspace, mainIssueID, firstCompletionID)
	if concernID == "" {
		t.Fatal("addressed concern id = empty")
	}
	closureRunFID(t, binary, workspace, "issue", "unclaim", mainIssueID, "--by", "developer")

	closureRunFID(t, binary, workspace, "link", "add", mainIssueID, "implements", "PRD-LINK-002", "--from-type", "issue", "--to-type", "requirement", "--by", "developer")
	closureRunFID(t, binary, workspace, "link", "add", peerIssueID, "discovered_from", mainIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
	assertClosureLinkReadback(t, binary, workspace, mainIssueID, peerIssueID)
	assertClosureRejectedWithoutJournalAppend(t, binary, workspace, journalPath,
		[]string{"link", "add", mainIssueID, "implements", peerIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer"},
		"relation")

	seqEpicID := closureCreateIssue(t, binary, workspace, "epic", "Sequential closure epic", stringPointer("sequential"), nil, nil)
	seqCurrentID := closureCreateIssue(t, binary, workspace, "task", "Sequential current", nil, stringPointer(seqEpicID), intPointer(1))
	seqNextID := closureCreateIssue(t, binary, workspace, "task", "Sequential next", nil, stringPointer(seqEpicID), intPointer(2))
	closureRunFID(t, binary, workspace, "issue", "transition", seqCurrentID, "--to", "cancelled", "--actor", "agent:tl")
	assertClosureReadyTask(t, binary, workspace, seqEpicID, seqNextID)

	priorityEpicID := closureCreateIssue(t, binary, workspace, "epic", "Prioritised closure epic", stringPointer("prioritised"), nil, nil)
	closureCreateIssue(t, binary, workspace, "task", "Lower priority", nil, stringPointer(priorityEpicID), nil, "--priority", "3")
	priorityHighID := closureCreateIssue(t, binary, workspace, "task", "Higher priority", nil, stringPointer(priorityEpicID), nil, "--priority", "0")
	assertClosureReadyTask(t, binary, workspace, priorityEpicID, priorityHighID)

	closureRunFID(t, binary, workspace, "run", "record", "--issue", mainIssueID, "--completion", firstCompletionID,
		"--role", "developer", "--model", "claude-sonnet-4-6", "--tokens-in", "1200", "--tokens-out", "3400", "--by", "developer")
	closureRunFID(t, binary, workspace, "run", "record", "--issue", mainIssueID, "--completion", addressCompletionID,
		"--role", "developer", "--model", "claude-sonnet-4-6", "--tokens-in", "0", "--tokens-out", "0", "--by", "developer")
	assertClosureTelemetry(t, binary, workspace, mainIssueID, firstCompletionID)

	journalBeforeNonJournal := closureReadFile(t, journalPath)
	journalFilesBefore := closureJournalFileBytes(t, workspace)
	assertClosureMaintenanceAndReadCommands(t, binary, workspace, mainIssueID, journalBeforeNonJournal)
	maintenanceNow := closureAgedMaintenanceTime(t, workspace)

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
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
	assertClosureJournalTracked(t, workspace)

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
