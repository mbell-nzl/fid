//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestIntegrationJournalAssertionsAcrossUTCMidnight(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	binary := integrationFIDBinary(t)
	init := p21Run(t, binary, workspace, "--json", "init", "--prefix", "FID")
	path := p21PathString(t, init, "result", "journal_path")
	var genesis map[string]any
	if err := json.Unmarshal(readP21File(t, path), &genesis); err != nil {
		t.Fatal(err)
	}
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	genesis["occurred_at"] = yesterday.Format(time.RFC3339Nano)
	genesis["event_id"] = "EVT-" + yesterday.Format("20060102") + "-000001"
	encoded, err := json.Marshal(genesis)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	path = filepath.Join(filepath.Dir(path), yesterday.Format("2006-01-02")+".jsonl")
	if err := os.WriteFile(path, append(encoded, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName)); err != nil {
		t.Fatal(err)
	}
	p21Run(t, binary, workspace, "--json", "rebuild")
	created := p21Run(t, binary, workspace, "--json", "issue", "create", "--kind", "task", "--title", "After UTC midnight", "--actor", "agent:developer", "--idempotency-key", "midnight-issue")
	p21AssertJournalTypes(t, workspace, []string{"project.initialized", "issue.created"})
	t.Run("closure journal", func(t *testing.T) {
		assertClosureJournalNeverTracked(t, workspace)
	})
	t.Run("event attribution", func(t *testing.T) {
		p20AssertJournalEventAttribution(t, workspace, p21PathString(t, created, "event", "id"), "developer")
	})
}
