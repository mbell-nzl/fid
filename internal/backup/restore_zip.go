package backup

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/mbell-nzl/fid/internal/export"
)

// checkExportIdentity checks the export's own identity as well as the
// backup wrapper's. Import validates the export, but cannot know which
// project the caller intended to restore.
func checkExportIdentity(repository clone, commit, id string) error {
	var payload bytes.Buffer
	found, err := folderBlob(repository, commit, id+"/manifest.json", &payload)
	if err != nil {
		return err
	}
	if !found {
		return damaged(id, folderCheck(CheckArchiveMissing, "%s/manifest.json is missing", id))
	}
	var header folderHeader
	if err := json.Unmarshal(payload.Bytes(), &header); err != nil {
		return damaged(id, folderCheck(CheckManifestInvalid, "%s/manifest.json is not a manifest: %v", id, err))
	}
	if header.ProjectID != id {
		return damaged(id, folderCheck(CheckManifestInvalid, "%s/manifest.json names project %q, not %s", id, header.ProjectID, id))
	}
	return nil
}

// zipFolder mirrors export's deterministic writer: sorted names, deflate,
// 0644 member modes and the ZIP epoch. Only member contents are stored in
// git; restoring them must reproduce the original export's archive digest.
// Members stream directly from git and are never extracted to local paths.
func zipFolder(repository clone, commit, id, expected, scratch string) (string, error) {
	return zipFolderWithLimits(repository, commit, id, expected, scratch, export.ArchiveResourceLimits())
}

func zipFolderWithLimits(repository clone, commit, id, expected, scratch string, limits export.ResourceLimits) (string, error) {
	entries, err := listFolderMembers(repository, commit, id, limits)
	if err != nil {
		return "", err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].name < entries[j].name
	})
	path := filepath.Join(scratch, ArchiveName)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return "", failure(KindInternal, err, "create the restore's archive")
	}
	defer func() { _ = file.Close() }()
	hash := sha256.New()
	writer := zip.NewWriter(io.MultiWriter(file, hash))
	for _, entry := range entries {
		header := &zip.FileHeader{Name: entry.name, Method: zip.Deflate, Modified: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)}
		header.SetMode(0o644)
		output, err := writer.CreateHeader(header)
		if err != nil {
			return "", failure(KindInternal, err, "add an export member to the restore's archive")
		}
		if err := repository.stream(output, "cat-file", "blob", entry.oid); err != nil {
			return "", failure(KindInternal, err, "read an export member from the backup repository")
		}
	}
	if err := writer.Close(); err != nil {
		return "", failure(KindInternal, err, "finish the restore's archive")
	}
	if err := file.Close(); err != nil {
		return "", failure(KindInternal, err, "write the restore's archive")
	}
	if digest := "sha256:" + hex.EncodeToString(hash.Sum(nil)); digest != expected {
		return "", damaged(id, folderCheck(CheckArchiveDigest, "the rebuilt %s/%s has digest %s, but %s names %q", id, ArchiveName, digest, ManifestName, expected))
	}
	return path, nil
}
