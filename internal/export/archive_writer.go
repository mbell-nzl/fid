package export

import (
	"archive/zip"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

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
	var sourceBytes int64 = 2
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
		sourceBytes += int64(len(canonical)) + 1
		if sourceBytes > maxArchiveNonArtifactBytes {
			return nil, errArchiveResourceLimit
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
		path := filepath.Join(config.LedgerDir(workspaceRoot), artifactRoot, name)
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() {
			return nil, fmt.Errorf("create export: available Artifact %q is unavailable", expected.Digest)
		}
		if info.Size() > maxArchiveArtifactBytes || expected.SizeBytes > maxArchiveArtifactBytes {
			return nil, errArchiveResourceLimit
		}
		actualDigest, err := fileDigest(path)
		if err != nil {
			return nil, fmt.Errorf("create export: read Artifact %q: %w", expected.Digest, err)
		}
		if info.Size() != expected.SizeBytes || actualDigest != expected.Digest {
			return nil, fmt.Errorf("create export: available Artifact %q does not match metadata", expected.Digest)
		}
		expected.Path = "artifacts/" + name
		entries = append(entries, archiveEntry{entry: expected, filePath: path})
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

func writeArchive(path string, entries []archiveEntry) error {
	if err := validateWriterInventory(entries); err != nil {
		return err
	}
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
		if err := writeArchiveEntry(output, item); err != nil {
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

func writeArchiveEntry(output io.Writer, item archiveEntry) error {
	if item.filePath == "" {
		_, err := output.Write(item.data)
		return err
	}
	input, err := os.Open(item.filePath)
	if err != nil {
		return err
	}
	defer func() { _ = input.Close() }()
	hash := sha256.New()
	n, err := io.Copy(io.MultiWriter(output, hash), io.LimitReader(input, item.SizeBytes+1))
	if err != nil {
		return err
	}
	if n != item.SizeBytes || "sha256:"+hex.EncodeToString(hash.Sum(nil)) != item.Digest {
		return fmt.Errorf("archive input changed while being read")
	}
	return nil
}

func validateWriterInventory(entries []archiveEntry) error {
	if len(entries) > maxArchiveEntries {
		return errArchiveResourceLimit
	}
	var total int64
	for _, item := range entries {
		limit := maxArchiveNonArtifactBytes
		if strings.HasPrefix(item.Path, "artifacts/") {
			limit = maxArchiveArtifactBytes
		}
		if item.SizeBytes < 0 || item.SizeBytes > limit || item.SizeBytes > maxArchiveExpandedBytes-total {
			return errArchiveResourceLimit
		}
		total += item.SizeBytes
	}
	return nil
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
