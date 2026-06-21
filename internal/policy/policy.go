package policy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

type MissingRoleError struct{}

func (err *MissingRoleError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return "--by is required"
}

type UnknownRoleError struct {
	Role       string
	KnownRoles []string
}

func (err *UnknownRoleError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("role %q not in project profile", err.Role)
}

func RequireKnownRole(workspaceRoot, role string) error {
	trimmedRole := strings.TrimSpace(role)
	if trimmedRole == "" {
		return &MissingRoleError{}
	}

	knownRoles, err := KnownRoles(workspaceRoot)
	if err != nil {
		return err
	}
	for _, known := range knownRoles {
		if known == trimmedRole {
			return nil
		}
	}

	return &UnknownRoleError{
		Role:       trimmedRole,
		KnownRoles: knownRoles,
	}
}

func KnownRoles(workspaceRoot string) ([]string, error) {
	if strings.TrimSpace(workspaceRoot) == "" {
		return nil, fmt.Errorf("known roles: workspace root is required")
	}

	journalDir := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName)
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		return nil, fmt.Errorf("read journal directory %q: %w", journalDir, err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)

	roles := make(map[string]bool)
	for _, name := range files {
		if err := collectRoles(filepath.Join(journalDir, name), roles); err != nil {
			return nil, err
		}
	}

	knownRoles := make([]string, 0, len(roles))
	for role, active := range roles {
		if active {
			knownRoles = append(knownRoles, role)
		}
	}
	sort.Strings(knownRoles)
	return knownRoles, nil
}

func collectRoles(path string, roles map[string]bool) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open journal file %q: %w", path, err)
	}
	defer closeRoleFile(file)

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var stored journal.Event
		if err := json.Unmarshal([]byte(line), &stored); err != nil {
			return &journal.CorruptJournalError{
				File: roleJournalErrorPath(path),
				Line: lineNumber,
				Err:  err,
			}
		}

		upgraded, err := journal.UpgradeToCurrent(stored)
		if err != nil {
			return fmt.Errorf("upgrade event %q from %s:%d: %w", stored.EventID, roleJournalErrorPath(path), lineNumber, err)
		}

		for _, effect := range upgraded.Effects {
			if effect.Table != "roles" {
				continue
			}

			roleName, _ := effect.PK["name"].(string)
			if roleName == "" {
				continue
			}

			switch effect.Op {
			case "insert":
				roles[roleName] = activeValue(effect.Values)
			case "update":
				roles[roleName] = activeValue(effect.Values)
			case "delete":
				delete(roles, roleName)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan journal file %q: %w", path, err)
	}

	return nil
}

func roleJournalErrorPath(path string) string {
	return filepath.ToSlash(filepath.Join(".fid", "journal", filepath.Base(path)))
}

func activeValue(values map[string]any) bool {
	if values == nil {
		return true
	}

	active, ok := values["active"]
	if !ok {
		return true
	}

	switch typed := active.(type) {
	case bool:
		return typed
	default:
		return true
	}
}

func closeRoleFile(closer io.Closer) {
	_ = closer.Close()
}
