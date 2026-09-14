// Package export creates and validates deterministic portable archive snapshots.
package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

const manifestSchema = "fid.export-manifest.v2"
const legacyManifestSchema = "fid.export-manifest.v1"

var zipEpoch = time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC)

// CreateOptions identifies the existing project facts and destination archive.
type CreateOptions struct {
	WorkspaceRoot string
	ArchivePath   string
}

// Summary is the stable result exposed by export create.
type Summary struct {
	ArchivePath       string `json:"archive_path"`
	ArchiveDigest     string `json:"archive_digest"`
	ManifestSchema    string `json:"manifest_schema"`
	ProjectID         string `json:"project_id"`
	EventCount        int    `json:"event_count"`
	TerminalEventHash string `json:"terminal_event_hash"`
	ProjectionSchema  int    `json:"projection_schema"`
	EntryCount        int    `json:"entry_count"`
}

// Validation is the stable result exposed by export validate.
type Validation struct {
	ArchivePath    string `json:"archive_path"`
	ArchiveDigest  string `json:"archive_digest"`
	Status         string `json:"status"`
	ManifestSchema string `json:"manifest_schema"`
	EntryCount     int    `json:"entry_count"`
}

// Check identifies one deterministic archive validation failure.
type Check struct {
	Code   string  `json:"code"`
	Entry  *string `json:"entry"`
	Detail string  `json:"detail"`
}

// InvalidArchiveError is returned when an archive fails integrity validation.
type InvalidArchiveError struct {
	ArchivePath string
	Checks      []Check
}

func (err *InvalidArchiveError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("archive %q is invalid", err.ArchivePath)
}

type manifest struct {
	Schema            string                `json:"schema"`
	ProjectID         string                `json:"project_id"`
	EventCount        int                   `json:"event_count"`
	TerminalEventHash string                `json:"terminal_event_hash"`
	ProjectionSchema  int                   `json:"projection_schema"`
	Capabilities      []capability          `json:"capabilities"`
	Sources           []json.RawMessage     `json:"sources"`
	Entries           []entry               `json:"entries"`
	ProjectConfig     *archiveProjectConfig `json:"project_config,omitempty"`
}

type capability struct {
	Available    bool           `json:"available"`
	CapabilityID string         `json:"capability_id"`
	Limits       map[string]any `json:"limits"`
	SchemaMax    int            `json:"schema_max"`
	SchemaMin    int            `json:"schema_min"`
	TrustSource  string         `json:"trust_source"`
}

type entry struct {
	Path      string `json:"path"`
	Digest    string `json:"digest"`
	SizeBytes int64  `json:"size_bytes"`
	MediaType string `json:"media_type"`
}

type archiveEntry struct {
	entry
	data     []byte
	filePath string
}

// Create writes a byte-stable archive from the current journal facts.
func Create(options CreateOptions) (Summary, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.ArchivePath) == "" {
		return Summary{}, fmt.Errorf("create export: workspace root and archive path are required")
	}
	project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
	if err != nil {
		return Summary{}, fmt.Errorf("create export: read project config: %w", err)
	}
	scratch, err := os.MkdirTemp(filepath.Dir(options.ArchivePath), ".fid-export-chunks-")
	if err != nil {
		return Summary{}, fmt.Errorf("create export: allocate chunks: %w", err)
	}
	defer func() { _ = os.RemoveAll(scratch) }()
	entries, eventCount, terminalHash, err := journalArchiveEntries(options.WorkspaceRoot, scratch)
	if err != nil {
		return Summary{}, err
	}
	sources, artifactEntries, err := supplementalArchiveFacts(options.WorkspaceRoot, project.ArtifactRoot)
	if err != nil {
		return Summary{}, err
	}
	entries = append(entries, artifactEntries...)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	manifestEntries := make([]entry, 0, len(entries))
	for _, item := range entries {
		manifestEntries = append(manifestEntries, item.entry)
	}
	encodedManifest, err := journal.MarshalCanonicalJSON(manifest{
		Schema:            manifestSchema,
		ProjectID:         project.ProjectID,
		EventCount:        eventCount,
		TerminalEventHash: terminalHash,
		ProjectionSchema:  fiddb.ActiveMigrationVersion,
		Capabilities:      capabilitySnapshot(),
		Sources:           sources,
		Entries:           manifestEntries,
		ProjectConfig:     archiveConfig(project),
	})
	if err != nil {
		return Summary{}, fmt.Errorf("create export: encode manifest: %w", err)
	}
	entries = append(entries, archiveEntry{entry: entry{Path: "manifest.json", Digest: digest(encodedManifest), SizeBytes: int64(len(encodedManifest)), MediaType: "application/json"}, data: encodedManifest})
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	if err := writeArchive(options.ArchivePath, entries); err != nil {
		return Summary{}, err
	}
	archiveDigest, err := fileDigest(options.ArchivePath)
	if err != nil {
		return Summary{}, err
	}
	return Summary{ArchivePath: options.ArchivePath, ArchiveDigest: archiveDigest, ManifestSchema: manifestSchema, ProjectID: project.ProjectID, EventCount: eventCount, TerminalEventHash: terminalHash, ProjectionSchema: fiddb.ActiveMigrationVersion, EntryCount: len(entries)}, nil
}
