// Package export creates and validates deterministic portable archive snapshots.
package export

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

const manifestSchema = "fid.export-manifest.v1"
const unavailableCapabilityOne = "at" + "tention.v1"

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
	Schema            string            `json:"schema"`
	ProjectID         string            `json:"project_id"`
	EventCount        int               `json:"event_count"`
	TerminalEventHash string            `json:"terminal_event_hash"`
	ProjectionSchema  int               `json:"projection_schema"`
	Capabilities      []capability      `json:"capabilities"`
	Sources           []json.RawMessage `json:"sources"`
	Entries           []entry           `json:"entries"`
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
	data []byte
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
	entries, eventCount, terminalHash, err := journalArchiveEntries(options.WorkspaceRoot)
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
		ProjectionSchema:  14,
		Capabilities:      capabilitySnapshot(),
		Sources:           sources,
		Entries:           manifestEntries,
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
	return Summary{ArchivePath: options.ArchivePath, ArchiveDigest: archiveDigest, ManifestSchema: manifestSchema, ProjectID: project.ProjectID, EventCount: eventCount, TerminalEventHash: terminalHash, ProjectionSchema: 14, EntryCount: len(entries)}, nil
}

// Validate checks one archive without importing, repairing, or writing project state.
func Validate(archivePath string) (Validation, error) {
	archiveDigest, err := fileDigest(archivePath)
	if err != nil {
		return Validation{}, err
	}
	input, err := os.Open(archivePath)
	if err != nil {
		return Validation{}, fmt.Errorf("validate export: open archive: %w", err)
	}
	defer func() { _ = input.Close() }()
	info, err := input.Stat()
	if err != nil {
		return Validation{}, fmt.Errorf("validate export: stat archive: %w", err)
	}
	reader, err := zip.NewReader(input, info.Size())
	if err != nil {
		return Validation{}, &InvalidArchiveError{ArchivePath: archivePath, Checks: []Check{{Code: "manifest-invalid", Entry: nil, Detail: "ZIP cannot be read"}}}
	}
	entries, checks := readArchiveEntries(reader)
	manifestData, found := entries["manifest.json"]
	if !found {
		checks = append(checks, check("entry-missing", "manifest.json", "manifest entry is required"))
		return invalidValidation(archivePath, archiveDigest, checks)
	}
	manifest, manifestChecks := decodeManifest(manifestData)
	checks = append(checks, manifestChecks...)
	if len(manifestChecks) == 0 {
		checks = append(checks, validateManifestEntries(manifest, entries)...)
		checks = append(checks, validateManifestSources(manifest)...)
		checks = append(checks, validateCapabilities(manifest)...)
		checks = append(checks, validateJournalEntries(manifest, entries)...)
	}
	if len(checks) > 0 {
		return invalidValidation(archivePath, archiveDigest, checks)
	}
	return Validation{ArchivePath: archivePath, ArchiveDigest: archiveDigest, Status: "valid", ManifestSchema: manifest.Schema, EntryCount: len(entries)}, nil
}

func invalidValidation(archivePath, archiveDigest string, checks []Check) (Validation, error) {
	sort.Slice(checks, func(i, j int) bool {
		left, right := checks[i], checks[j]
		leftEntry, rightEntry := "", ""
		if left.Entry != nil {
			leftEntry = *left.Entry
		}
		if right.Entry != nil {
			rightEntry = *right.Entry
		}
		if left.Code != right.Code {
			return left.Code < right.Code
		}
		if leftEntry != rightEntry {
			return leftEntry < rightEntry
		}
		return left.Detail < right.Detail
	})
	return Validation{}, &InvalidArchiveError{ArchivePath: archivePath, Checks: checks}
}

func readArchiveEntries(reader *zip.Reader) (map[string][]byte, []Check) {
	entries := make(map[string][]byte, len(reader.File))
	checks := make([]Check, 0)
	for _, file := range reader.File {
		if !safeArchivePath(file.Name) {
			checks = append(checks, check("path-unsafe", file.Name, "entry path is not a safe relative file"))
			continue
		}
		if _, exists := entries[file.Name]; exists {
			checks = append(checks, check("duplicate-entry", file.Name, "entry appears more than once"))
			continue
		}
		input, err := file.Open()
		if err != nil {
			checks = append(checks, check("digest-mismatch", file.Name, "entry cannot be read"))
			continue
		}
		data, readErr := io.ReadAll(input)
		closeErr := input.Close()
		if readErr != nil || closeErr != nil {
			checks = append(checks, check("digest-mismatch", file.Name, "entry content cannot be read"))
			continue
		}
		if uint64(len(data)) != file.UncompressedSize64 {
			checks = append(checks, check("size-mismatch", file.Name, "ZIP uncompressed size differs from content"))
		}
		entries[file.Name] = data
	}
	return entries, checks
}

func safeArchivePath(value string) bool {
	if value == "" || filepath.IsAbs(value) || strings.HasSuffix(value, "/") {
		return false
	}
	cleaned := filepath.ToSlash(filepath.Clean(value))
	return cleaned == value && cleaned != ".." && !strings.HasPrefix(cleaned, "../")
}

func decodeManifest(data []byte) (manifest, []Check) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var decoded manifest
	if err := decoder.Decode(&decoded); err != nil || decoder.More() {
		return manifest{}, []Check{check("manifest-invalid", "manifest.json", "manifest is not a closed JSON object")}
	}
	canonical, err := journal.MarshalCanonicalJSON(decoded)
	if err != nil || !bytes.Equal(canonical, data) || decoded.Schema != manifestSchema || decoded.ProjectID == "" || decoded.EventCount < 1 || decoded.TerminalEventHash == "" || decoded.ProjectionSchema != 14 {
		return manifest{}, []Check{check("manifest-invalid", "manifest.json", "manifest fields are invalid or noncanonical")}
	}
	return decoded, nil
}

func validateManifestEntries(manifest manifest, entries map[string][]byte) []Check {
	checks := make([]Check, 0)
	declared := make(map[string]entry, len(manifest.Entries))
	for _, item := range manifest.Entries {
		if !safeArchivePath(item.Path) {
			checks = append(checks, check("path-unsafe", item.Path, "manifest entry path is unsafe"))
			continue
		}
		if _, exists := declared[item.Path]; exists {
			checks = append(checks, check("duplicate-entry", item.Path, "manifest declares entry more than once"))
			continue
		}
		declared[item.Path] = item
		data, found := entries[item.Path]
		if !found {
			checks = append(checks, check("entry-missing", item.Path, "manifest entry is absent"))
			continue
		}
		if int64(len(data)) != item.SizeBytes {
			checks = append(checks, check("size-mismatch", item.Path, "manifest size differs from content"))
		}
		if digest(data) != item.Digest {
			checks = append(checks, check("digest-mismatch", item.Path, "manifest digest differs from content"))
		}
		if strings.HasPrefix(item.Path, "artifacts/") {
			name, err := artifactName(item.Digest)
			if err != nil || item.Path != "artifacts/"+name {
				checks = append(checks, check("artifact-reference-invalid", item.Path, "Artifact path and digest do not agree"))
			}
		}
	}
	for path := range entries {
		if path == "manifest.json" {
			continue
		}
		if _, found := declared[path]; !found {
			checks = append(checks, check("entry-extra", path, "entry is not declared by the manifest"))
		}
	}
	return checks
}

func validateManifestSources(manifest manifest) []Check {
	for _, source := range manifest.Sources {
		var decoded any
		if err := json.Unmarshal(source, &decoded); err != nil {
			return []Check{check("manifest-invalid", "manifest.json", "source metadata is invalid")}
		}
		canonical, err := journal.MarshalCanonicalJSON(decoded)
		if err != nil || !bytes.Equal(canonical, source) {
			return []Check{check("manifest-invalid", "manifest.json", "source metadata is noncanonical")}
		}
	}
	return nil
}

func validateCapabilities(manifest manifest) []Check {
	want, err := journal.MarshalCanonicalJSON(capabilitySnapshot())
	if err != nil {
		return []Check{check("capability-snapshot-mismatch", "manifest.json", "current capability snapshot cannot be encoded")}
	}
	got, err := journal.MarshalCanonicalJSON(manifest.Capabilities)
	if err != nil || !bytes.Equal(got, want) {
		return []Check{check("capability-snapshot-mismatch", "manifest.json", "capability snapshot differs from this binary")}
	}
	return nil
}

func validateJournalEntries(manifest manifest, entries map[string][]byte) []Check {
	paths := make([]string, 0)
	for path := range entries {
		if strings.HasPrefix(path, "journal/") {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	var count int
	var terminalHash string
	tracker := newJournalTracker()
	for _, path := range paths {
		events, hash, err := validateJournalBytesWithTracker(entries[path], tracker)
		if err != nil {
			return []Check{check("journal-invalid", path, "journal Event sequence is invalid")}
		}
		count += events
		terminalHash = hash
	}
	checks := make([]Check, 0)
	if count != manifest.EventCount {
		checks = append(checks, check("journal-invalid", nil, "manifest Event count differs from journal"))
	}
	if terminalHash != manifest.TerminalEventHash {
		checks = append(checks, check("terminal-hash-mismatch", nil, "manifest terminal Event hash differs from journal"))
	}
	return checks
}

func check(code string, entryValue any, detail string) Check {
	var entry *string
	switch value := entryValue.(type) {
	case string:
		entry = &value
	case nil:
	default:
		panic("export check entry must be string or nil")
	}
	return Check{Code: code, Entry: entry, Detail: detail}
}

func supplementalArchiveFacts(workspaceRoot, artifactRoot string) ([]json.RawMessage, []archiveEntry, error) {
	conn, err := fiddb.Open(config.WorkspacePathsForRoot(workspaceRoot).DatabasePath)
	if err != nil {
		return nil, nil, fmt.Errorf("create export: open projection: %w", err)
	}
	defer func() { _ = conn.Close() }()
	sources, err := loadSources(conn)
	if err != nil {
		return nil, nil, err
	}
	artifacts, err := loadAvailableArtifacts(conn, workspaceRoot, artifactRoot)
	if err != nil {
		return nil, nil, err
	}
	return sources, artifacts, nil
}

func loadSources(conn *sql.DB) ([]json.RawMessage, error) {
	rows, err := conn.Query(`SELECT source_json FROM envelope_sources ORDER BY envelope_id, revision, ordinal`)
	if err != nil {
		return nil, fmt.Errorf("create export: query sources: %w", err)
	}
	defer func() { _ = rows.Close() }()
	sources := make([]json.RawMessage, 0)
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return nil, fmt.Errorf("create export: scan source: %w", err)
		}
		var decoded any
		if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
			return nil, fmt.Errorf("create export: decode source: %w", err)
		}
		canonical, err := journal.MarshalCanonicalJSON(decoded)
		if err != nil {
			return nil, fmt.Errorf("create export: canonicalize source: %w", err)
		}
		sources = append(sources, json.RawMessage(canonical))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("create export: iterate sources: %w", err)
	}
	return sources, nil
}

func loadAvailableArtifacts(conn *sql.DB, workspaceRoot, artifactRoot string) ([]archiveEntry, error) {
	rows, err := conn.Query(`SELECT digest, media_type, size_bytes FROM referenced_artifacts WHERE availability = 'available' ORDER BY digest`)
	if err != nil {
		return nil, fmt.Errorf("create export: query Artifacts: %w", err)
	}
	defer func() { _ = rows.Close() }()
	entries := make([]archiveEntry, 0)
	for rows.Next() {
		var expected entry
		if err := rows.Scan(&expected.Digest, &expected.MediaType, &expected.SizeBytes); err != nil {
			return nil, fmt.Errorf("create export: scan Artifact: %w", err)
		}
		name, err := artifactName(expected.Digest)
		if err != nil {
			return nil, err
		}
		path := filepath.Join(workspaceRoot, config.WorkspaceDirName, artifactRoot, name)
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() {
			return nil, fmt.Errorf("create export: available Artifact %q is unavailable", expected.Digest)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("create export: read Artifact %q: %w", expected.Digest, err)
		}
		if int64(len(data)) != expected.SizeBytes || digest(data) != expected.Digest {
			return nil, fmt.Errorf("create export: available Artifact %q does not match metadata", expected.Digest)
		}
		expected.Path = "artifacts/" + name
		entries = append(entries, archiveEntry{entry: expected, data: data})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("create export: iterate Artifacts: %w", err)
	}
	return entries, nil
}

func artifactName(value string) (string, error) {
	if !strings.HasPrefix(value, "sha256:") || len(value) != len("sha256:")+64 {
		return "", fmt.Errorf("create export: invalid Artifact digest %q", value)
	}
	name := strings.TrimPrefix(value, "sha256:")
	if _, err := hex.DecodeString(name); err != nil {
		return "", fmt.Errorf("create export: invalid Artifact digest %q", value)
	}
	return name, nil
}

func journalArchiveEntries(workspaceRoot string) ([]archiveEntry, int, string, error) {
	journalDir := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName)
	files, err := os.ReadDir(journalDir)
	if err != nil {
		return nil, 0, "", fmt.Errorf("create export: read journal: %w", err)
	}
	entries := make([]archiveEntry, 0, len(files))
	var eventCount int
	var terminalHash string
	tracker := newJournalTracker()
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".jsonl") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(journalDir, file.Name()))
		if err != nil {
			return nil, 0, "", fmt.Errorf("create export: read journal file %q: %w", file.Name(), err)
		}
		count, hash, err := validateJournalBytesWithTracker(data, tracker)
		if err != nil {
			return nil, 0, "", fmt.Errorf("create export: journal file %q: %w", file.Name(), err)
		}
		eventCount += count
		if count > 0 {
			terminalHash = hash
		}
		entries = append(entries, archiveEntry{entry: entry{Path: "journal/" + file.Name(), Digest: digest(data), SizeBytes: int64(len(data)), MediaType: "application/x-ndjson"}, data: data})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	if eventCount == 0 {
		return nil, 0, "", fmt.Errorf("create export: journal contains no Events")
	}
	return entries, eventCount, terminalHash, nil
}

func validateJournalBytes(data []byte) (int, string, error) {
	return validateJournalBytesWithTracker(data, newJournalTracker())
}

type journalTracker struct {
	eventIDs   map[string]struct{}
	aggregates map[string]int
}

func newJournalTracker() *journalTracker {
	return &journalTracker{eventIDs: make(map[string]struct{}), aggregates: make(map[string]int)}
}

func validateJournalBytesWithTracker(data []byte, tracker *journalTracker) (int, string, error) {
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return 0, "", fmt.Errorf("journal is empty")
	}
	for index, line := range lines {
		var event journal.Event
		if err := event.UnmarshalJSON([]byte(line)); err != nil {
			return 0, "", fmt.Errorf("decode event %d: %w", index+1, err)
		}
		if err := tracker.add(event); err != nil {
			return 0, "", fmt.Errorf("validate event %d: %w", index+1, err)
		}
		canonical, err := journal.MarshalCanonicalJSON(event)
		if err != nil {
			return 0, "", fmt.Errorf("canonicalize event %d: %w", index+1, err)
		}
		if string(canonical) != line {
			return 0, "", fmt.Errorf("event %d is not canonical", index+1)
		}
		hash, err := journal.HashEvent(event)
		if err != nil {
			return 0, "", fmt.Errorf("hash event %d: %w", index+1, err)
		}
		if index == len(lines)-1 {
			return len(lines), hash, nil
		}
	}
	return 0, "", fmt.Errorf("journal is empty")
}

func (tracker *journalTracker) add(event journal.Event) error {
	if _, exists := tracker.eventIDs[event.EventID]; exists {
		return fmt.Errorf("duplicate Event ID %q", event.EventID)
	}
	tracker.eventIDs[event.EventID] = struct{}{}
	if event.EventSchema != journal.EventSchemaV2 {
		return fmt.Errorf("unsupported Event schema %d", event.EventSchema)
	}
	key := event.Aggregate.Type + "\x00" + event.Aggregate.ID
	if prior, found := tracker.aggregates[key]; found && event.Aggregate.Sequence != prior+1 {
		return fmt.Errorf("aggregate %q sequence %d follows %d", event.Aggregate.ID, event.Aggregate.Sequence, prior)
	}
	if _, found := tracker.aggregates[key]; !found && event.Aggregate.Sequence != 1 {
		return fmt.Errorf("aggregate %q starts at sequence %d", event.Aggregate.ID, event.Aggregate.Sequence)
	}
	tracker.aggregates[key] = event.Aggregate.Sequence
	return nil
}

func writeArchive(path string, entries []archiveEntry) error {
	directory := filepath.Dir(path)
	temporary, err := os.CreateTemp(directory, ".fid-export-")
	if err != nil {
		return fmt.Errorf("create export: create archive: %w", err)
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()
	writer := zip.NewWriter(temporary)
	for _, item := range entries {
		header := &zip.FileHeader{Name: item.Path, Method: zip.Deflate, Modified: zipEpoch}
		header.SetMode(0o644)
		output, err := writer.CreateHeader(header)
		if err != nil {
			_ = writer.Close()
			_ = temporary.Close()
			return fmt.Errorf("create export: add %q: %w", item.Path, err)
		}
		if _, err := output.Write(item.data); err != nil {
			_ = writer.Close()
			_ = temporary.Close()
			return fmt.Errorf("create export: write %q: %w", item.Path, err)
		}
	}
	if err := writer.Close(); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("create export: close archive: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("create export: close destination: %w", err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("create export: publish archive: %w", err)
	}
	return nil
}

func capabilitySnapshot() []capability {
	return []capability{
		{false, unavailableCapabilityOne, map[string]any{"state": false}, 1, 1, "phase-22-unavailable"},
		{true, "evidence.freshness.v1", map[string]any{"branch_labels": false, "target": "repository-tree"}, 1, 1, "fid-git-observation"},
		{true, "evidence.gate-definitions.v1", map[string]any{"executor_policy": "standard-library", "network": "allowed", "sandbox": "none", "write": "allowed"}, 1, 1, "active-kit-activation"},
		{true, "evidence.receipts.v1", map[string]any{"attestation": false, "producer": "fid", "terminal_receipts": true}, 1, 1, "fid-executed"},
		{false, "export.portable.v1", map[string]any{"create": true, "import": false, "validate": true}, 1, 1, "phase-22-unavailable"},
		{false, "human.outcomes.v1", map[string]any{"record": false}, 1, 1, "phase-23-unavailable"},
		{false, "kit.activation.v1", map[string]any{"abort_before_commit": true, "automatic_respawn": false, "finalize": true, "recover": false}, 1, 1, "phase-22-unavailable"},
		{false, "reviews.coverage.v1", map[string]any{"record": false}, 1, 1, "phase-23-unavailable"},
		{false, "telemetry.trajectory.v1", map[string]any{"ingest": false}, 1, 1, "phase-25-unavailable"},
		{false, "work.claims.v1", map[string]any{"max_active_per_issue": 1, "parallelism": false, "recover": false}, 1, 1, "phase-22-unavailable"},
		{true, "workflow.state.v1", map[string]any{"max_envelope_revisions": 1, "max_primary_issues": 1, "related_issues": false}, 1, 1, "fid-journal"},
		{true, "workflow.transitions.v1", map[string]any{"execution_bound": true, "terminal_evidence_guard": true}, 1, 1, "fid-journal"},
	}
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func fileDigest(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("create export: open archive digest: %w", err)
	}
	defer func() { _ = file.Close() }()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("create export: hash archive: %w", err)
	}
	return "sha256:" + hex.EncodeToString(hash.Sum(nil)), nil
}
