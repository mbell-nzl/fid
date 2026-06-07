package journal

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

type ReplayStats struct {
	EventsReplayed int
	FilesReplayed  int
	UpgradedEvents int
}

type ReplayOptions struct {
	StartFile   string
	StartOffset int64
}

type CorruptJournalError struct {
	File string
	Line int
	Err  error
}

func (err *CorruptJournalError) Error() string {
	if err == nil {
		return "<nil>"
	}

	return fmt.Sprintf("malformed journal entry at %s:%d: %v", err.File, err.Line, err.Err)
}

func (err *CorruptJournalError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

type EnsureProjectionResult struct {
	Decision    fiddb.ProjectionFreshness
	ReplayStats ReplayStats
}

type journalCheckpoint struct {
	File      string
	Offset    int64
	EventID   string
	EventHash string
	HasEvents bool
}

type journalInspection struct {
	tail           journalCheckpoint
	checkpointSeen bool
	checkpoint     journalCheckpoint
}

func ReplayJournal(conn *sql.DB, journalDir string, options *ReplayOptions) (ReplayStats, error) {
	if conn == nil {
		return ReplayStats{}, fmt.Errorf("replay journal: connection is nil")
	}

	files, err := orderedJournalFiles(journalDir)
	if err != nil {
		return ReplayStats{}, err
	}

	var stats ReplayStats
	for _, fileName := range files {
		fileReplayed, err := replayJournalFile(conn, journalDir, fileName, options, &stats)
		if err != nil {
			return ReplayStats{}, err
		}
		if fileReplayed {
			stats.FilesReplayed++
		}
	}

	return stats, nil
}

func EnsureProjectionFresh(workspaceRoot, migrationsDir string) (EnsureProjectionResult, error) {
	dbPath := filepath.Join(workspaceRoot, ".fid", "fid.sqlite")
	journalDir := filepath.Join(workspaceRoot, ".fid", "journal")

	if _, err := os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			stats, rebuildErr := rebuildWorkspaceWithOwner(workspaceRoot, migrationsDir, LockOwner{
				Command: "fid catch-up",
			})
			if rebuildErr != nil {
				return EnsureProjectionResult{}, rebuildErr
			}
			return EnsureProjectionResult{
				Decision:    fiddb.ProjectionRebuild,
				ReplayStats: stats,
			}, nil
		}
		return EnsureProjectionResult{}, fmt.Errorf("stat sqlite database %q: %w", dbPath, err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		return EnsureProjectionResult{}, err
	}
	defer closeReplayDB(conn)

	meta, metaExists, err := fiddb.ReadProjectionMeta(conn)
	if err != nil {
		return EnsureProjectionResult{}, err
	}

	inspection, err := inspectJournal(journalDir, meta.JournalFile, meta.JournalOffset)
	if err != nil {
		return EnsureProjectionResult{}, err
	}

	decision := fiddb.DecideProjectionFreshness(fiddb.ProjectionFreshnessInput{
		MetaExists:                    metaExists,
		ReplayStatus:                  meta.ReplayStatus,
		StoredProjectionSchemaVersion: meta.ProjectionSchemaVersion,
		CurrentProjectionSchema:       fiddb.ProjectionSchemaVersion,
		CheckpointExists:              inspection.checkpointSeen,
		CheckpointMatches: metaExists &&
			inspection.checkpointSeen &&
			meta.LastEventID == inspection.checkpoint.EventID &&
			meta.LastEventHash == inspection.checkpoint.EventHash,
		JournalHasNewerEvents: inspection.tail.HasEvents &&
			(!metaExists ||
				meta.JournalFile != inspection.tail.File ||
				meta.JournalOffset != inspection.tail.Offset),
	})

	switch decision {
	case fiddb.ProjectionFresh:
		return EnsureProjectionResult{Decision: decision}, nil
	case fiddb.ProjectionCatchUp:
		stats, err := catchUpWorkspace(workspaceRoot, migrationsDir, meta)
		if err != nil {
			return EnsureProjectionResult{}, err
		}
		return EnsureProjectionResult{Decision: decision, ReplayStats: stats}, nil
	case fiddb.ProjectionRebuild:
		stats, err := rebuildWorkspaceWithOwner(workspaceRoot, migrationsDir, LockOwner{
			Command: "fid catch-up",
		})
		if err != nil {
			return EnsureProjectionResult{}, err
		}
		return EnsureProjectionResult{Decision: decision, ReplayStats: stats}, nil
	default:
		return EnsureProjectionResult{}, fmt.Errorf("ensure projection fresh: unsupported decision %q", decision)
	}
}

func RebuildWorkspace(workspaceRoot, migrationsDir, byRole string) (ReplayStats, error) {
	return rebuildWorkspaceWithOwner(workspaceRoot, migrationsDir, LockOwner{
		Command:   "fid rebuild --by " + byRole,
		ActorRole: byRole,
	})
}

func rebuildWorkspaceWithOwner(workspaceRoot, migrationsDir string, owner LockOwner) (ReplayStats, error) {
	var stats ReplayStats
	err := WithLock(workspaceRoot, owner, func(_ *WriteLock) error {
		dbPath := filepath.Join(workspaceRoot, ".fid", "fid.sqlite")
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove sqlite database %q: %w", dbPath, err)
		}
		if err := fiddb.MigrateUp(dbPath, migrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeReplayDB(conn)

		replayed, err := ReplayJournal(conn, filepath.Join(workspaceRoot, ".fid", "journal"), nil)
		if err != nil {
			return err
		}
		stats = replayed
		return nil
	})
	if err != nil {
		return ReplayStats{}, err
	}

	return stats, nil
}

func catchUpWorkspace(workspaceRoot, migrationsDir string, meta fiddb.ProjectionMeta) (ReplayStats, error) {
	var stats ReplayStats
	err := WithLock(workspaceRoot, LockOwner{Command: "fid catch-up"}, func(_ *WriteLock) error {
		dbPath := filepath.Join(workspaceRoot, ".fid", "fid.sqlite")
		if err := fiddb.MigrateUp(dbPath, migrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeReplayDB(conn)

		replayed, err := ReplayJournal(conn, filepath.Join(workspaceRoot, ".fid", "journal"), &ReplayOptions{
			StartFile:   meta.JournalFile,
			StartOffset: meta.JournalOffset,
		})
		if err != nil {
			return err
		}
		stats = replayed
		return nil
	})
	if err != nil {
		return ReplayStats{}, err
	}

	return stats, nil
}

func replayJournalFile(conn *sql.DB, journalDir, fileName string, options *ReplayOptions, stats *ReplayStats) (bool, error) {
	path := filepath.Join(journalDir, fileName)
	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("open journal file %q: %w", path, err)
	}
	defer closeReplayFile(file)

	reader := bufio.NewReader(file)
	var (
		lineNumber   int
		offset       int64
		fileReplayed bool
		startFile    string
		startOffset  int64
	)
	if options != nil {
		startFile = options.StartFile
		startOffset = options.StartOffset
	}

	for {
		lineBytes, readErr := reader.ReadBytes('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return false, fmt.Errorf("read journal file %q: %w", path, readErr)
		}
		if len(lineBytes) == 0 && errors.Is(readErr, io.EOF) {
			break
		}

		lineNumber++
		offset += int64(len(lineBytes))

		trimmed := bytesTrimSpace(lineBytes)
		if len(trimmed) == 0 {
			if errors.Is(readErr, io.EOF) {
				break
			}
			continue
		}

		if startFile != "" {
			if fileName < startFile || (fileName == startFile && offset <= startOffset) {
				if errors.Is(readErr, io.EOF) {
					break
				}
				continue
			}
		}

		var stored Event
		if err := json.Unmarshal(trimmed, &stored); err != nil {
			return false, &CorruptJournalError{
				File: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
				Line: lineNumber,
				Err:  err,
			}
		}

		eventHash, err := HashEvent(stored)
		if err != nil {
			return false, fmt.Errorf("hash event %q from %s:%d: %w", stored.EventID, fileName, lineNumber, err)
		}

		upgraded, err := UpgradeToCurrent(stored)
		if err != nil {
			return false, fmt.Errorf("upgrade event %q from %s:%d: %w", stored.EventID, fileName, lineNumber, err)
		}
		if upgraded.EventSchema != stored.EventSchema {
			stats.UpgradedEvents++
		}

		if err := ProjectEvent(conn, ProjectionInput{
			Event:      upgraded,
			EventHash:  eventHash,
			File:       fileName,
			Line:       lineNumber,
			Offset:     offset,
			ReplayMode: fiddb.ReplayStatusClean,
		}); err != nil {
			return false, err
		}

		stats.EventsReplayed++
		fileReplayed = true

		if errors.Is(readErr, io.EOF) {
			break
		}
	}

	return fileReplayed, nil
}

func orderedJournalFiles(journalDir string) ([]string, error) {
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		return nil, fmt.Errorf("read journal directory %q: %w", journalDir, err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".jsonl") {
			continue
		}
		if _, err := time.Parse("2006-01-02", strings.TrimSuffix(name, ".jsonl")); err != nil {
			continue
		}
		files = append(files, name)
	}
	sort.Strings(files)
	return files, nil
}

func inspectJournal(journalDir, targetFile string, targetOffset int64) (journalInspection, error) {
	files, err := orderedJournalFiles(journalDir)
	if err != nil {
		return journalInspection{}, err
	}

	var inspection journalInspection
	for _, fileName := range files {
		path := filepath.Join(journalDir, fileName)
		file, err := os.Open(path)
		if err != nil {
			return journalInspection{}, fmt.Errorf("open journal file %q: %w", path, err)
		}

		reader := bufio.NewReader(file)
		var (
			offset     int64
			lineNumber int
		)
		for {
			lineBytes, readErr := reader.ReadBytes('\n')
			if readErr != nil && !errors.Is(readErr, io.EOF) {
				_ = file.Close()
				return journalInspection{}, fmt.Errorf("read journal file %q: %w", path, readErr)
			}
			if len(lineBytes) == 0 && errors.Is(readErr, io.EOF) {
				break
			}

			lineNumber++
			offset += int64(len(lineBytes))
			trimmed := bytesTrimSpace(lineBytes)
			if len(trimmed) == 0 {
				if errors.Is(readErr, io.EOF) {
					break
				}
				continue
			}

			var stored Event
			if err := json.Unmarshal(trimmed, &stored); err != nil {
				_ = file.Close()
				return journalInspection{}, &CorruptJournalError{
					File: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
					Line: lineNumber,
					Err:  err,
				}
			}

			hash, err := HashEvent(stored)
			if err != nil {
				_ = file.Close()
				return journalInspection{}, fmt.Errorf("hash event %q from %s: %w", stored.EventID, fileName, err)
			}

			checkpoint := journalCheckpoint{
				File:      fileName,
				Offset:    offset,
				EventID:   stored.EventID,
				EventHash: hash,
				HasEvents: true,
			}
			inspection.tail = checkpoint
			if fileName == targetFile && offset == targetOffset {
				inspection.checkpointSeen = true
				inspection.checkpoint = checkpoint
			}

			if errors.Is(readErr, io.EOF) {
				break
			}
		}

		if err := file.Close(); err != nil {
			return journalInspection{}, fmt.Errorf("close journal file %q: %w", path, err)
		}
	}

	return inspection, nil
}

func bytesTrimSpace(value []byte) []byte {
	return []byte(strings.TrimSpace(string(value)))
}

func closeReplayDB(conn *sql.DB) {
	if conn != nil {
		_ = conn.Close()
	}
}

func closeReplayFile(file *os.File) {
	if file != nil {
		_ = file.Close()
	}
}
