package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	JournalDirName   = "journal"
	ConfigFileName   = "config.toml"
	DatabaseFileName = "fid.sqlite"
	GitIgnoreName    = ".gitignore"
)

const (
	// LedgerDirName is the ledger's name in the git common directory. Before
	// ADR-0013 the ledger was .fid in the worktree; fid ignores a directory of
	// that name left there.
	LedgerDirName = "fid"
	// LedgerStagingPrefix starts the name of a directory, beside the ledger,
	// that init and import fill before they rename it into place.
	LedgerStagingPrefix = "fid-staging-"
)

type WorkspacePaths struct {
	WorkspaceDir string
	JournalDir   string
	ConfigPath   string
	DatabasePath string
	GitIgnore    string
}

// LedgerDir returns the directory that holds the ledger of the project whose
// worktree top level is root: its journal, projection, config, artifacts and
// write lock. It is the only place that decides where the ledger lives; every
// ledger path starts from it.
//
// The ledger is the fid directory in the git common directory, so every
// worktree of a repository shares it and git never tracks it (ADR-0013).
// LedgerDir finds the common directory by reading root's .git entry as git
// does, without running git: a directory, or a file naming the git
// directory, plus that directory's commondir file. It returns no error: where
// root has no .git entry it can read, it names root/.git/fid, where a main
// worktree keeps the ledger.
func LedgerDir(root string) string {
	return filepath.Join(gitCommonDir(root), LedgerDirName)
}

func gitCommonDir(root string) string {
	gitDir := filepath.Join(root, ".git")
	if info, err := os.Stat(gitDir); err == nil && info.Mode().IsRegular() {
		data, err := os.ReadFile(gitDir)
		named, ok := strings.CutPrefix(strings.TrimRight(string(data), "\r\n"), "gitdir: ")
		if err != nil || !ok || named == "" {
			return gitDir
		}
		if !filepath.IsAbs(named) {
			named = filepath.Join(root, named)
		}
		gitDir = filepath.Clean(named)
	}
	data, err := os.ReadFile(filepath.Join(gitDir, "commondir"))
	common := strings.TrimRight(string(data), "\r\n")
	if err != nil || common == "" {
		return gitDir
	}
	if !filepath.IsAbs(common) {
		common = filepath.Join(gitDir, common)
	}
	return filepath.Clean(common)
}

// NewLedgerStaging creates a staging directory beside ledgerDir, in the git
// common directory, so a ledger staged in it renames into place on the same
// filesystem. It returns the staging directory, which the caller removes, and
// the path inside it to stage the ledger at.
func NewLedgerStaging(ledgerDir string) (staging, stagedLedger string, err error) {
	staging, err = os.MkdirTemp(filepath.Dir(ledgerDir), LedgerStagingPrefix)
	if err != nil {
		return "", "", fmt.Errorf("create staging directory: %w", err)
	}
	return staging, filepath.Join(staging, LedgerDirName), nil
}

func WorkspacePathsForRoot(workspaceRoot string) WorkspacePaths {
	return LedgerPaths(LedgerDir(workspaceRoot))
}

// LedgerPaths names the files of the ledger at ledgerDir. A staging directory
// has no .git, so code that fills one names its files from here.
func LedgerPaths(ledgerDir string) WorkspacePaths {
	return WorkspacePaths{
		WorkspaceDir: ledgerDir,
		JournalDir:   filepath.Join(ledgerDir, JournalDirName),
		ConfigPath:   filepath.Join(ledgerDir, ConfigFileName),
		DatabasePath: filepath.Join(ledgerDir, DatabaseFileName),
		GitIgnore:    filepath.Join(ledgerDir, GitIgnoreName),
	}
}

// WriteLedgerFiles writes a new ledger's config and .gitignore into ledgerDir,
// and creates its journal directory.
func WriteLedgerFiles(ledgerDir string, project ProjectConfig) (WorkspacePaths, error) {
	if strings.TrimSpace(ledgerDir) == "" {
		return WorkspacePaths{}, fmt.Errorf("write bootstrap files: ledger directory is required")
	}
	if err := validateProjectConfig(project, ledgerDir); err != nil {
		return WorkspacePaths{}, fmt.Errorf("write bootstrap files: %w", err)
	}

	paths := LedgerPaths(ledgerDir)
	if err := os.MkdirAll(paths.JournalDir, 0o755); err != nil {
		return WorkspacePaths{}, fmt.Errorf("create journal directory %q: %w", paths.JournalDir, err)
	}

	if err := os.WriteFile(paths.GitIgnore, []byte(gitIgnoreContents()), 0o644); err != nil {
		return WorkspacePaths{}, fmt.Errorf("write gitignore %q: %w", paths.GitIgnore, err)
	}
	if err := os.WriteFile(paths.ConfigPath, []byte(projectConfigContents(project)), 0o644); err != nil {
		return WorkspacePaths{}, fmt.Errorf("write config %q: %w", paths.ConfigPath, err)
	}

	return paths, nil
}

func gitIgnoreContents() string {
	return DatabaseFileName + "\n"
}

func normalizeLine(line string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return ""
	}
	return trimmed
}

func isSection(line string) bool {
	return strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")
}

func splitAssignment(line string) (string, string, bool) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
}

func parseQuotedString(value string) (string, error) {
	if len(value) < 2 || value[0] != '"' || value[len(value)-1] != '"' {
		return "", fmt.Errorf("expected quoted string, got %q", value)
	}
	return value[1 : len(value)-1], nil
}

func closeConfigFile(file *os.File) {
	_ = file.Close()
}
