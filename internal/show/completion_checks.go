package show

import (
	"fmt"
	"path/filepath"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

func (c *Completion) checksBeforeCompletion(root string, records []checkrecord.Record) ([]checkrecord.Record, error) {
	finished := map[string]bool{}
	found := false
	err := journal.Visit(filepath.Join(config.LedgerDir(root), config.JournalDirName), func(event journal.Event) (bool, error) {
		for _, e := range event.Effects {
			if e.Op != "insert" {
				continue
			}
			if e.Table == "completions" && e.PK["id"] == c.ID {
				found = true
				return true, nil
			}
			if e.Table == "check_results" {
				id, _ := e.PK["check_id"].(string)
				finished[id] = true
			}
		}
		return false, nil
	})
	if err != nil {
		return nil, fmt.Errorf("read completed check boundary: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("completion %s missing from journal", c.ID)
	}
	before := []checkrecord.Record{}
	// checkread.Read already supplies finish journal order. Filter by result
	// insertion, never by start insertion: a pending check can outlive done.
	for _, r := range records {
		if finished[r.Start.ID] {
			before = append(before, r)
		}
	}
	return before, nil
}
