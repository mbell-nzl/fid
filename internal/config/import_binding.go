package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ImportBinding is private delivery metadata for one atomically installed archive.
// It is intentionally outside the journal and portable archive content.
type ImportBinding struct {
	ArchiveDigest    string               `json:"archive_digest"`
	IdempotencyKey   string               `json:"idempotency_key"`
	ProjectID        string               `json:"project_id"`
	TargetRoot       string               `json:"target_root"`
	ImportedEventIDs []string             `json:"imported_event_ids"`
	JournalEntries   []ImportJournalEntry `json:"journal_entries"`
}

type ImportJournalEntry struct {
	Name      string `json:"name"`
	Digest    string `json:"digest"`
	SizeBytes int64  `json:"size_bytes"`
}

const ImportBindingFileName = "import.marker.json"

func ReadImportBinding(workspaceRoot string) (ImportBinding, error) {
	root, err := filepath.Abs(workspaceRoot)
	if err != nil {
		return ImportBinding{}, fmt.Errorf("read import binding: resolve workspace: %w", err)
	}
	workspace := LedgerDir(root)
	info, err := os.Lstat(workspace)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return ImportBinding{}, fmt.Errorf("read import binding: workspace is unsafe")
	}
	marker := filepath.Join(workspace, ImportBindingFileName)
	info, err = os.Lstat(marker)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return ImportBinding{}, fmt.Errorf("read import binding: delivery metadata is unsafe")
	}
	data, err := os.ReadFile(marker)
	if err != nil {
		return ImportBinding{}, err
	}
	var binding ImportBinding
	if err := json.Unmarshal(data, &binding); err != nil || binding.ProjectID == "" || binding.TargetRoot == "" || binding.ArchiveDigest == "" || binding.IdempotencyKey == "" || len(binding.ImportedEventIDs) == 0 {
		return ImportBinding{}, fmt.Errorf("read import binding: invalid delivery metadata")
	}
	if binding.TargetRoot != root {
		return ImportBinding{}, fmt.Errorf("read import binding: target root does not match workspace")
	}
	project, err := ReadProjectConfig(filepath.Join(workspace, ConfigFileName))
	if err != nil || project.ProjectID != binding.ProjectID {
		return ImportBinding{}, fmt.Errorf("read import binding: project binding is invalid")
	}
	journalInfo, err := os.Lstat(filepath.Join(workspace, JournalDirName))
	if err != nil || !journalInfo.IsDir() || journalInfo.Mode()&os.ModeSymlink != 0 {
		return ImportBinding{}, fmt.Errorf("read import binding: journal is unsafe")
	}
	if err := verifyImportedHistory(workspace, binding); err != nil {
		return ImportBinding{}, err
	}
	return binding, nil
}
