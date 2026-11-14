package export

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Validate checks one archive without importing, repairing, or writing project state.
type ValidateOptions struct {
	ArchivePath   string
	MigrationsDir string
}

func Validate(options ValidateOptions) (Validation, error) {
	archivePath := options.ArchivePath
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
	if err := inspectArchiveEntryCount(input, info.Size()); err != nil {
		code := "manifest-invalid"
		if errors.Is(err, errArchiveResourceLimit) {
			code = "resource-limit"
		}
		return invalidValidation(archivePath, archiveDigest, []Check{check(code, nil, err.Error())})
	}
	reader, err := zip.NewReader(input, info.Size())
	if err != nil {
		return Validation{}, &InvalidArchiveError{ArchivePath: archivePath, Checks: []Check{{Code: "manifest-invalid", Entry: nil, Detail: "ZIP cannot be read"}}}
	}
	entries := make(map[string]scannedEntry)
	var manifestData []byte
	scan, checks := scanArchive(reader, func(member archiveMember, input io.Reader) error {
		hash := sha256.New()
		var destination io.Writer = hash
		var buffer bytes.Buffer
		if member.path == "manifest.json" {
			destination = io.MultiWriter(hash, &buffer)
		}
		n, err := io.Copy(destination, input)
		if err != nil {
			return err
		}
		entries[member.path] = scannedEntry{Digest: "sha256:" + hex.EncodeToString(hash.Sum(nil)), SizeBytes: n, File: member.file}
		if member.path == "manifest.json" {
			manifestData = buffer.Bytes()
		}
		return nil
	})
	if len(checks) > 0 {
		return invalidValidation(archivePath, archiveDigest, checks)
	}
	if manifestData == nil {
		return invalidValidation(archivePath, archiveDigest, []Check{check("entry-missing", "manifest.json", "manifest entry is required")})
	}
	manifest, manifestChecks := decodeManifest(manifestData)
	checks = append(checks, manifestChecks...)
	if len(manifestChecks) == 0 {
		checks = append(checks, validateManifestEntries(manifest, entries)...)
		profileChecks := append(validateManifestSources(manifest), validateCapabilities(manifest)...)
		if len(checks) == 0 {
			journalChecks := validateJournalEntries(manifest, entries, options.MigrationsDir)
			// An old profile cannot disguise removed history as a mere profile mismatch.
			// Entry integrity still precedes semantic inspection of archive contents.
			unsupported := false
			for _, item := range journalChecks {
				unsupported = unsupported || item.Code == "schema-unsupported"
			}
			if unsupported || len(profileChecks) == 0 {
				checks = journalChecks
			} else {
				checks = profileChecks
			}
		} else {
			checks = append(checks, profileChecks...)
		}

	}
	if len(checks) > 0 {
		return invalidValidation(archivePath, archiveDigest, checks)
	}
	return Validation{ArchivePath: archivePath, ArchiveDigest: archiveDigest, Status: "valid", ManifestSchema: manifest.Schema, EntryCount: scan.EntryCount}, nil
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

func decodeManifest(data []byte) (manifest, []Check) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var decoded manifest
	if err := decoder.Decode(&decoded); err != nil || decoder.More() {
		return manifest{}, []Check{check("manifest-invalid", "manifest.json", "manifest is not a closed JSON object")}
	}
	canonical, err := journal.MarshalCanonicalJSON(decoded)
	if err != nil || !bytes.Equal(canonical, data) || (decoded.Schema != manifestSchema && decoded.Schema != legacyManifestSchema) || decoded.ProjectID == "" || decoded.EventCount < 1 || decoded.TerminalEventHash == "" {
		return manifest{}, []Check{check("manifest-invalid", "manifest.json", "manifest fields are invalid or noncanonical")}
	}
	if decoded.Schema != manifestSchema || decoded.ProjectionSchema != fiddb.ActiveMigrationVersion {
		return manifest{}, []Check{check("schema-unsupported", "manifest.json", "unsupported archive schema or projection")}
	}
	if _, err := decoded.ProjectConfig.projectConfig(decoded.ProjectID); err != nil {
		return manifest{}, []Check{check("configuration-invalid", "manifest.json", err.Error())}
	}
	return decoded, nil
}

func validateManifestEntries(manifest manifest, entries map[string]scannedEntry) []Check {
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
		if !strings.HasPrefix(item.Path, "journal/") && !strings.HasPrefix(item.Path, "artifacts/") {
			checks = append(checks, check("entry-extra", item.Path, "entry is outside accepted journal and Artifact families"))
		}
		if item.SizeBytes < 0 || item.MediaType == "" {
			checks = append(checks, check("manifest-invalid", item.Path, "entry metadata is invalid"))
		}
		data, found := entries[item.Path]
		if !found {
			checks = append(checks, check("entry-missing", item.Path, "manifest entry is absent"))
			continue
		}
		if data.SizeBytes != item.SizeBytes {
			checks = append(checks, check("size-mismatch", item.Path, "manifest size differs from content"))
		}
		if data.Digest != item.Digest {
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
	if manifest.Sources == nil || len(manifest.Sources) != 0 {
		return []Check{check("reference-invalid", "manifest.json", "Source inventory must be empty")}
	}

	return nil
}

func validateCapabilities(manifest manifest) []Check {
	profile := capabilitySnapshot()
	want, err := journal.MarshalCanonicalJSON(profile)
	if err != nil {
		return []Check{check("capability-snapshot-mismatch", "manifest.json", "current capability snapshot cannot be encoded")}
	}
	got, err := journal.MarshalCanonicalJSON(manifest.Capabilities)
	if err != nil || !bytes.Equal(got, want) {
		return []Check{check("capability-snapshot-mismatch", "manifest.json", "capability snapshot differs from this binary")}
	}
	return nil
}

func validateJournalEntries(manifest manifest, entries map[string]scannedEntry, migrationsDir string) []Check {
	paths := make([]string, 0)
	for path := range entries {
		if strings.HasPrefix(path, "journal/") {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	if manifest.Schema == manifestSchema {
		for index, path := range paths {
			if path != fmt.Sprintf("journal/%08d.jsonl", index+1) {
				return []Check{check("journal-invalid", path, "journal chunk ordinals must be contiguous from 00000001")}
			}
		}
	}
	var count int
	var terminalHash string
	tracker := newJournalTracker()
	var semantic *semanticValidator
	if manifest.Schema == manifestSchema {
		var err error
		semantic, err = newSemanticValidator(migrationsDir, manifest.ProjectionSchema)
		if err != nil {
			return []Check{check("schema-unsupported", nil, err.Error())}
		}
	}
	defer func() { _ = semantic.close() }()
	for _, path := range paths {
		_, err := streamArchiveMember(archiveMember{file: entries[path].File, path: path}, newArchiveByteBudget(), func(input io.Reader) error {
			return readEventLines(input, func(line []byte) error {
				hash, err := validateJournalLine(line, tracker)
				if err != nil {
					return err
				}
				count++
				terminalHash = hash
				if semantic != nil {
					return semantic.accept(line)
				}
				return nil
			})
		})
		if err != nil {
			code := "journal-invalid"
			if errors.Is(err, errArchiveResourceLimit) {
				code = "resource-limit"
			} else if semanticCode := semanticCheckCode(err); semanticCode != "" {
				code = semanticCode
			}
			return []Check{check(code, path, err.Error())}
		}
	}
	checks := make([]Check, 0)
	if count != manifest.EventCount {
		checks = append(checks, check("journal-invalid", nil, "manifest Event count differs from journal"))
	}
	if terminalHash != manifest.TerminalEventHash {
		checks = append(checks, check("terminal-hash-mismatch", nil, "manifest terminal Event hash differs from journal"))
	}
	if semantic != nil {
		facts, err := semantic.finish()
		if err != nil {
			checks = append(checks, check("reference-invalid", nil, err.Error()))
		} else {
			checks = append(checks, validateSemanticManifest(manifest, entries, facts)...)
		}
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

type scannedEntry struct {
	Digest    string
	SizeBytes int64
	File      *zip.File
}
