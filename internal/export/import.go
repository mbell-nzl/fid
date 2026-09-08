package export

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"golang.org/x/sys/unix"
)

type ImportStage string

const (
	ImportBeforePublication ImportStage = "before-publication"
	ImportAfterPublication  ImportStage = "after-publication"
	importMarkerName                    = "import.marker.json"
)

var ErrImportFault = errors.New("import fault injected")

type ImportTargetConflictError struct{ TargetRoot, Reason string }

func (err *ImportTargetConflictError) Error() string {
	return fmt.Sprintf("import target %q conflicts: %s", err.TargetRoot, err.Reason)
}

type ImportOptions struct {
	ArchivePath, RepositoryRoot, MigrationsDir, IdempotencyKey string
	Fault                                                      func(ImportStage) error
}

type ImportResult struct {
	ProjectID         string `json:"project_id"`
	ArchiveDigest     string `json:"archive_digest"`
	EventCount        int    `json:"event_count"`
	TerminalEventHash string `json:"terminal_event_hash"`
	ProjectionSchema  int    `json:"projection_schema"`
	TargetRoot        string `json:"target_root"`
	Replayed          bool   `json:"replayed"`
}

type importMarker struct {
	config.ImportBinding
	Result ImportResult `json:"result"`
}

// Import validates, stages, rebuilds and atomically installs a portable archive.
// Publication is delivery metadata and does not append a domain Event.
func Import(options ImportOptions) (ImportResult, error) {
	if strings.TrimSpace(options.ArchivePath) == "" || strings.TrimSpace(options.RepositoryRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || strings.TrimSpace(options.IdempotencyKey) == "" {
		return ImportResult{}, fmt.Errorf("import: archive, repository, migrations and idempotency key are required")
	}
	target, err := safeImportTarget(options.RepositoryRoot)
	if err != nil {
		return ImportResult{}, err
	}
	digestBefore, err := fileDigest(options.ArchivePath)
	if err != nil {
		return ImportResult{}, err
	}
	snapshot, err := snapshotImportInput(target, options.ArchivePath, digestBefore)
	if err != nil {
		return ImportResult{}, err
	}
	defer func() { _ = os.Remove(snapshot) }()
	validation, err := Validate(ValidateOptions{ArchivePath: snapshot, MigrationsDir: options.MigrationsDir})
	if err != nil {
		return ImportResult{}, err
	}
	if validation.ManifestSchema != manifestSchema {
		entry := "manifest.json"
		return ImportResult{}, &InvalidArchiveError{ArchivePath: options.ArchivePath, Checks: []Check{{Code: "schema-unsupported", Entry: &entry, Detail: "import requires recovery archive schema"}}}
	}
	if validation.ArchiveDigest != digestBefore {
		return ImportResult{}, fmt.Errorf("import: archive changed during validation")
	}
	manifest, importedEventIDs, err := importArchive(snapshot)
	if err != nil {
		return ImportResult{}, err
	}
	if manifest.Schema != manifestSchema || manifest.ProjectConfig == nil {
		return ImportResult{}, fmt.Errorf("import: archive schema is unsupported")
	}
	projectConfig, err := manifest.ProjectConfig.projectConfig(manifest.ProjectID)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import: project configuration: %w", err)
	}
	result := ImportResult{ProjectID: manifest.ProjectID, ArchiveDigest: digestBefore, EventCount: manifest.EventCount, TerminalEventHash: manifest.TerminalEventHash, ProjectionSchema: manifest.ProjectionSchema, TargetRoot: target}
	err = withImportLock(target, func() error {
		if replay, found, err := installedImport(target, digestBefore, options.IdempotencyKey); err != nil {
			return err
		} else if found {
			result = replay
			result.Replayed = true
			return nil
		}
		if err := ensureEmptyImportTarget(target); err != nil {
			return err
		}
		stage, err := os.MkdirTemp(target, ".fid-import-stage-")
		if err != nil {
			return fmt.Errorf("import: create stage: %w", err)
		}
		defer func() { _ = os.RemoveAll(stage) }()
		if _, err := config.WriteBootstrapFiles(stage, projectConfig, config.DefaultPricingConfig()); err != nil {
			return fmt.Errorf("import: create staged configuration: %w", err)
		}
		if err := stageArchive(stage, projectConfig, manifest, snapshot); err != nil {
			return err
		}
		if _, err := journal.RebuildWorkspace(stage, options.MigrationsDir); err != nil {
			return fmt.Errorf("import: rebuild staged projection: %w", err)
		}
		binding := config.ImportBinding{ArchiveDigest: digestBefore, IdempotencyKey: options.IdempotencyKey, ProjectID: result.ProjectID, TargetRoot: target, ImportedEventIDs: importedEventIDs}
		if err := writeImportMarker(filepath.Join(stage, config.WorkspaceDirName, importMarkerName), importMarker{ImportBinding: binding, Result: result}); err != nil {
			return err
		}
		if err := syncTree(filepath.Join(stage, config.WorkspaceDirName)); err != nil {
			return err
		}
		if err := callImportFault(options.Fault, ImportBeforePublication); err != nil {
			return err
		}
		if err := publishImport(filepath.Join(stage, config.WorkspaceDirName), filepath.Join(target, config.WorkspaceDirName)); err != nil {
			return err
		}
		if err := syncPath(target); err != nil {
			return err
		}
		return callImportFault(options.Fault, ImportAfterPublication)
	})
	return result, err
}

func safeImportTarget(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", &ImportTargetConflictError{TargetRoot: root, Reason: "unsafe-path"}
	}
	info, err := os.Lstat(abs)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return "", &ImportTargetConflictError{TargetRoot: abs, Reason: "unsafe-path"}
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil || resolved != abs {
		return "", &ImportTargetConflictError{TargetRoot: abs, Reason: "unsafe-path"}
	}
	return abs, nil
}

func withImportLock(target string, fn func() error) error {
	path := filepath.Join(target, ".fid-import.lock")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("import: open target lock: %w", err)
	}
	defer func() { _ = file.Close() }()
	// The advisory lock is released by the kernel when a process exits. Removing
	// the pathname on a normal return keeps the destination worktree clean;
	// a stale pathname after a crash remains safely reusable.
	defer func() { _ = os.Remove(path) }()
	if err := unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		if errors.Is(err, unix.EWOULDBLOCK) {
			return fmt.Errorf("import: target lock unavailable")
		}
		return fmt.Errorf("import: acquire target lock: %w", err)
	}
	defer func() { _ = unix.Flock(int(file.Fd()), unix.LOCK_UN) }()
	return fn()
}

func snapshotImportInput(target, archivePath, digest string) (string, error) {
	input, err := os.Open(archivePath)
	if err != nil {
		return "", fmt.Errorf("import: open archive: %w", err)
	}
	defer func() { _ = input.Close() }()
	output, err := os.CreateTemp(target, ".fid-import-input-")
	if err != nil {
		return "", fmt.Errorf("import: create archive snapshot: %w", err)
	}
	path := output.Name()
	defer func() {
		if output != nil {
			_ = output.Close()
		}
	}()
	if _, err := io.Copy(output, input); err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("import: snapshot archive: %w", err)
	}
	if err := output.Close(); err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("import: close archive snapshot: %w", err)
	}
	output = nil
	got, err := fileDigest(path)
	if err != nil || got != digest {
		_ = os.Remove(path)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("import: archive changed while snapshotting")
	}
	return path, nil
}

func ensureEmptyImportTarget(target string) error {
	info, err := os.Lstat(filepath.Join(target, config.WorkspaceDirName))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return &ImportTargetConflictError{TargetRoot: target, Reason: "unsafe-path"}
	}
	return &ImportTargetConflictError{TargetRoot: target, Reason: "populated"}
}

func installedImport(target, archiveDigest, key string) (ImportResult, bool, error) {
	data, err := os.ReadFile(filepath.Join(target, config.WorkspaceDirName, config.ImportBindingFileName))
	if os.IsNotExist(err) {
		return ImportResult{}, false, nil
	}
	if err != nil {
		return ImportResult{}, false, &ImportTargetConflictError{TargetRoot: target, Reason: "populated"}
	}
	var marker importMarker
	if err := json.Unmarshal(data, &marker); err != nil || marker.Result.TargetRoot != target || marker.Result.ProjectID != marker.ProjectID || marker.Result.ArchiveDigest != marker.ArchiveDigest {
		return ImportResult{}, false, &ImportTargetConflictError{TargetRoot: target, Reason: "populated"}
	}
	if _, err := config.ReadImportBinding(target); err != nil {
		return ImportResult{}, false, &ImportTargetConflictError{TargetRoot: target, Reason: "populated"}
	}
	if marker.ArchiveDigest != archiveDigest || marker.IdempotencyKey != key {
		return ImportResult{}, false, &ImportTargetConflictError{TargetRoot: target, Reason: "different-import"}
	}
	return marker.Result, true, nil
}

func importArchive(path string) (manifest, []string, error) {
	input, err := os.Open(path)
	if err != nil {
		return manifest{}, nil, fmt.Errorf("import: open archive: %w", err)
	}
	defer func() { _ = input.Close() }()
	info, err := input.Stat()
	if err != nil {
		return manifest{}, nil, fmt.Errorf("import: stat archive: %w", err)
	}
	reader, err := zip.NewReader(input, info.Size())
	if err != nil {
		return manifest{}, nil, fmt.Errorf("import: read archive: %w", err)
	}
	members, checks := inspectArchiveDirectory(reader.File)
	if len(checks) != 0 {
		return manifest{}, nil, fmt.Errorf("import: archive directory changed after validation")
	}
	var data []byte
	for _, member := range members {
		if member.path == "manifest.json" {
			_, err := streamArchiveMember(member, newArchiveByteBudget(), func(input io.Reader) error { var readErr error; data, readErr = io.ReadAll(input); return readErr })
			if err != nil {
				return manifest{}, nil, fmt.Errorf("import: read manifest: %w", err)
			}
		}
	}
	decoded, checks := decodeManifest(data)
	if len(checks) != 0 {
		return manifest{}, nil, fmt.Errorf("import: manifest changed after validation")
	}
	ids := make([]string, 0, decoded.EventCount)
	for _, member := range members {
		if !strings.HasPrefix(member.path, "journal/") {
			continue
		}
		_, err := streamArchiveMember(member, newArchiveByteBudget(), func(input io.Reader) error {
			return readEventLines(input, func(line []byte) error {
				var event journal.Event
				if err := event.UnmarshalJSON(line); err != nil {
					return err
				}
				ids = append(ids, event.EventID)
				return nil
			})
		})
		if err != nil {
			return manifest{}, nil, fmt.Errorf("import: read Event identities: %w", err)
		}
	}
	return decoded, ids, nil
}

func stageArchive(stage string, project config.ProjectConfig, manifest manifest, archivePath string) error {
	input, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("import: reopen archive for staging: %w", err)
	}
	defer func() { _ = input.Close() }()
	info, err := input.Stat()
	if err != nil {
		return fmt.Errorf("import: stat archive for staging: %w", err)
	}
	reader, err := zip.NewReader(input, info.Size())
	if err != nil {
		return fmt.Errorf("import: reread archive for staging: %w", err)
	}
	members, checks := inspectArchiveDirectory(reader.File)
	if len(checks) != 0 {
		return fmt.Errorf("import: archive directory changed during staging")
	}
	files := make(map[string]*zip.File, len(members))
	for _, member := range members {
		files[member.path] = member.file
	}
	budget := newArchiveByteBudget()
	for _, item := range manifest.Entries {
		file := files[item.Path]
		if file == nil {
			return fmt.Errorf("import: archive entry %q is absent", item.Path)
		}
		path, err := stagedArchivePath(stage, project, item.Path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("import: create staged entry directory: %w", err)
		}
		output, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err != nil {
			return fmt.Errorf("import: create staged entry %q: %w", item.Path, err)
		}
		hash := sha256.New()
		n, copyErr := streamArchiveMember(archiveMember{file: file, path: item.Path, artifact: strings.HasPrefix(item.Path, "artifacts/")}, budget, func(input io.Reader) error { _, err := io.Copy(io.MultiWriter(output, hash), input); return err })
		closeErr := output.Close()
		if copyErr != nil || closeErr != nil || n != item.SizeBytes || "sha256:"+hex.EncodeToString(hash.Sum(nil)) != item.Digest {
			_ = os.Remove(path)
			return fmt.Errorf("import: staged entry %q does not match its manifest (size %d/%d digest %s/%s copy %v close %v)", item.Path, n, item.SizeBytes, "sha256:"+hex.EncodeToString(hash.Sum(nil)), item.Digest, copyErr, closeErr)
		}
	}
	return nil
}

func stagedArchivePath(stage string, project config.ProjectConfig, entryPath string) (string, error) {
	if strings.HasPrefix(entryPath, "journal/") {
		return filepath.Join(stage, config.WorkspaceDirName, config.JournalDirName, strings.TrimPrefix(entryPath, "journal/")), nil
	}
	if strings.HasPrefix(entryPath, "artifacts/") {
		return filepath.Join(stage, config.WorkspaceDirName, project.ArtifactRoot, strings.TrimPrefix(entryPath, "artifacts/")), nil
	}
	return "", fmt.Errorf("import: unsupported archive entry %q", entryPath)
}

func writeImportMarker(path string, marker importMarker) error {
	payload, err := json.Marshal(marker)
	if err != nil {
		return fmt.Errorf("import: encode publication marker: %w", err)
	}
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		return fmt.Errorf("import: write publication marker: %w", err)
	}
	return nil
}
func publishImport(stage, destination string) error {
	if err := unix.Renameat2(unix.AT_FDCWD, stage, unix.AT_FDCWD, destination, unix.RENAME_NOREPLACE); err != nil {
		if errors.Is(err, unix.EEXIST) {
			return &ImportTargetConflictError{TargetRoot: filepath.Dir(destination), Reason: "populated"}
		}
		return fmt.Errorf("import: atomically publish staged project: %w", err)
	}
	return nil
}
func syncTree(root string) error {
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("import: staged symlink is unsafe")
		}
		if entry.IsDir() || entry.Type().IsRegular() {
			return syncPath(path)
		}
		return fmt.Errorf("import: staged special file is unsafe")
	})
}
func syncPath(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("import: open for sync %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()
	if err := file.Sync(); err != nil {
		return fmt.Errorf("import: sync %q: %w", path, err)
	}
	return nil
}
func callImportFault(fault func(ImportStage) error, stage ImportStage) error {
	if fault == nil {
		return nil
	}
	return fault(stage)
}
