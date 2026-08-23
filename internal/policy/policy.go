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

	scan := roleScan{roles: make(map[string]bool)}
	for _, name := range files {
		if err := collectRoles(filepath.Join(journalDir, name), &scan); err != nil {
			return nil, err
		}
	}

	if scan.hasExactRoleFreeBootstrap && !scan.hasLegacyRoleEffect {
		knownRoles := MatrixRoleNames()
		sort.Strings(knownRoles)
		return knownRoles, nil
	}

	knownRoles := make([]string, 0, len(scan.roles))
	for role, active := range scan.roles {
		if active {
			knownRoles = append(knownRoles, role)
		}
	}
	sort.Strings(knownRoles)
	return knownRoles, nil
}

type roleScan struct {
	roles                     map[string]bool
	hasExactRoleFreeBootstrap bool
	hasLegacyRoleEffect       bool
}

func collectRoles(path string, scan *roleScan) error {
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
		if isRoleFreeV2Bootstrap(stored) {
			scan.hasExactRoleFreeBootstrap = true
		}

		for _, effect := range upgraded.Effects {
			if effect.Table != "roles" {
				continue
			}
			scan.hasLegacyRoleEffect = true

			roleName, _ := effect.PK["name"].(string)
			if roleName == "" {
				continue
			}

			switch effect.Op {
			case "insert":
				scan.roles[roleName] = activeValue(effect.Values)
			case "update":
				scan.roles[roleName] = activeValue(effect.Values)
			case "delete":
				delete(scan.roles, roleName)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan journal file %q: %w", path, err)
	}

	return nil
}

func isRoleFreeV2Bootstrap(event journal.Event) bool {
	return event.EventSchema == journal.EventSchemaV2 &&
		event.EventType == "project.initialized" &&
		event.ActorV2.ActorID == "fid-bootstrap" &&
		event.ActorV2.ActorType == "service" &&
		event.ActorV2.Producer == "fid" &&
		event.ActorV2.ProvenanceRef == "fid:init" &&
		event.ActorV2.SessionRef == nil &&
		event.Aggregate.ID != "" &&
		event.Aggregate.Sequence == 1 &&
		event.Aggregate.Type == "project" &&
		event.CausationID == nil &&
		event.CommandName == "project.init" &&
		event.CorrelationID == nil &&
		len(event.Effects) == 0 &&
		event.Idempotency == nil &&
		event.Repository == nil
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
