package config

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	WorkspaceDirName                    = ".fid"
	JournalDirName                      = "journal"
	ConfigFileName                      = "config.toml"
	PricingFileName                     = "pricing.toml"
	DatabaseFileName                    = "fid.sqlite"
	GitIgnoreName                       = ".gitignore"
	ConfigVersion                       = 1
	DefaultMaxOutputBytes               = int64(10485760)
	DefaultMaxExcerptBytes              = int64(16384)
	DefaultAbandonedAttemptGraceSeconds = int64(300)
)

type ProjectConfig struct {
	ProjectID                    string
	IssuePrefix                  string
	ArtifactRoot                 string
	MaxOutputBytes               int64
	MaxExcerptBytes              int64
	AbandonedAttemptGraceSeconds int64
}

type ModelPricing struct {
	InPerMTok  float64
	OutPerMTok float64
}

type PricingConfig struct {
	Models map[string]ModelPricing
}

var defaultPricingModelOrder = []string{
	"claude-opus-4-7",
	"claude-sonnet-4-6",
	"claude-haiku-4-5",
	"gpt-5.4-mini",
	"gpt-5.4-nano",
	"gemini-2.5-flash",
	"gemini-2.5-flash-lite",
}

type WorkspacePaths struct {
	WorkspaceDir string
	JournalDir   string
	ConfigPath   string
	PricingPath  string
	DatabasePath string
	GitIgnore    string
}

func NewDefaultProjectConfig(prefix string) (ProjectConfig, error) {
	id, err := newProjectID()
	if err != nil {
		return ProjectConfig{}, err
	}
	return ProjectConfig{
		ProjectID:                    id,
		IssuePrefix:                  prefix,
		ArtifactRoot:                 "artifacts",
		MaxOutputBytes:               DefaultMaxOutputBytes,
		MaxExcerptBytes:              DefaultMaxExcerptBytes,
		AbandonedAttemptGraceSeconds: DefaultAbandonedAttemptGraceSeconds,
	}, nil
}

func DefaultPricingConfig() PricingConfig {
	return PricingConfig{
		Models: map[string]ModelPricing{
			"claude-opus-4-7": {
				InPerMTok:  15.0,
				OutPerMTok: 75.0,
			},
			"claude-sonnet-4-6": {
				InPerMTok:  3.0,
				OutPerMTok: 15.0,
			},
			"claude-haiku-4-5": {
				InPerMTok:  0.80,
				OutPerMTok: 4.0,
			},
			"gpt-5.4-mini": {
				InPerMTok:  0.375,
				OutPerMTok: 2.25,
			},
			"gpt-5.4-nano": {
				InPerMTok:  0.10,
				OutPerMTok: 0.625,
			},
			"gemini-2.5-flash": {
				InPerMTok:  0.30,
				OutPerMTok: 2.50,
			},
			"gemini-2.5-flash-lite": {
				InPerMTok:  0.10,
				OutPerMTok: 0.40,
			},
		},
	}
}

func WorkspacePathsForRoot(workspaceRoot string) WorkspacePaths {
	workspaceDir := filepath.Join(workspaceRoot, WorkspaceDirName)
	return WorkspacePaths{
		WorkspaceDir: workspaceDir,
		JournalDir:   filepath.Join(workspaceDir, JournalDirName),
		ConfigPath:   filepath.Join(workspaceDir, ConfigFileName),
		PricingPath:  filepath.Join(workspaceDir, PricingFileName),
		DatabasePath: filepath.Join(workspaceDir, DatabaseFileName),
		GitIgnore:    filepath.Join(workspaceDir, GitIgnoreName),
	}
}

func WriteBootstrapFiles(workspaceRoot string, project ProjectConfig, pricing PricingConfig) (WorkspacePaths, error) {
	if strings.TrimSpace(workspaceRoot) == "" {
		return WorkspacePaths{}, fmt.Errorf("write bootstrap files: workspace root is required")
	}
	if err := validateProjectConfig(project, pathsForConfigRoot(workspaceRoot)); err != nil {
		return WorkspacePaths{}, fmt.Errorf("write bootstrap files: %w", err)
	}

	paths := WorkspacePathsForRoot(workspaceRoot)
	if err := os.MkdirAll(paths.JournalDir, 0o755); err != nil {
		return WorkspacePaths{}, fmt.Errorf("create journal directory %q: %w", paths.JournalDir, err)
	}

	if err := os.WriteFile(paths.GitIgnore, []byte(gitIgnoreContents()), 0o644); err != nil {
		return WorkspacePaths{}, fmt.Errorf("write gitignore %q: %w", paths.GitIgnore, err)
	}
	if err := os.WriteFile(paths.ConfigPath, []byte(projectConfigContents(project)), 0o644); err != nil {
		return WorkspacePaths{}, fmt.Errorf("write config %q: %w", paths.ConfigPath, err)
	}
	if err := os.WriteFile(paths.PricingPath, []byte(pricingConfigContents(pricing)), 0o644); err != nil {
		return WorkspacePaths{}, fmt.Errorf("write pricing %q: %w", paths.PricingPath, err)
	}

	return paths, nil
}

func ReadProjectConfig(path string) (ProjectConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return ProjectConfig{}, fmt.Errorf("open project config %q: %w", path, err)
	}
	defer closeConfigFile(file)

	scanner := bufio.NewScanner(file)
	section := ""
	values := make(map[string]string)
	seenSections := make(map[string]bool)
	for scanner.Scan() {
		line := normalizeLine(scanner.Text())
		if line == "" {
			continue
		}
		if isSection(line) {
			section = strings.Trim(line, "[]")
			if !isAllowedSection(section) || seenSections[section] {
				return ProjectConfig{}, fmt.Errorf("parse project config %q: unsupported or duplicate section %q", path, section)
			}
			seenSections[section] = true
			continue
		}
		key, value, ok := splitAssignment(line)
		if !ok {
			return ProjectConfig{}, fmt.Errorf("parse project config %q: invalid line %q", path, line)
		}
		name := key
		if section != "" {
			name = section + "." + key
		}
		if !isAllowedConfigKey(section, key) {
			return ProjectConfig{}, fmt.Errorf("parse project config %q: unsupported key %q", path, name)
		}
		if _, duplicate := values[name]; duplicate {
			return ProjectConfig{}, fmt.Errorf("parse project config %q: duplicate key %q", path, name)
		}
		values[name] = value
	}
	if err := scanner.Err(); err != nil {
		return ProjectConfig{}, fmt.Errorf("scan project config %q: %w", path, err)
	}
	project, err := parseProjectConfigValues(values)
	if err != nil {
		return ProjectConfig{}, fmt.Errorf("parse project config %q: %w", path, err)
	}
	if err := validateProjectConfig(project, filepath.Dir(path)); err != nil {
		return ProjectConfig{}, fmt.Errorf("parse project config %q: %w", path, err)
	}
	return project, nil
}

func ReadPricingConfig(path string) (PricingConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return PricingConfig{}, fmt.Errorf("open pricing config %q: %w", path, err)
	}
	defer closeConfigFile(file)

	pricing := PricingConfig{Models: make(map[string]ModelPricing)}
	scanner := bufio.NewScanner(file)
	currentModel := ""
	for scanner.Scan() {
		line := normalizeLine(scanner.Text())
		if line == "" {
			continue
		}
		if isSection(line) {
			section := strings.Trim(line, "[]")
			currentModel = ""
			if strings.HasPrefix(section, "models.") {
				currentModel = strings.TrimPrefix(section, "models.")
			}
			continue
		}
		if currentModel == "" {
			continue
		}
		key, value, ok := splitAssignment(line)
		if !ok {
			continue
		}
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return PricingConfig{}, fmt.Errorf("parse pricing config %q model %q key %q: %w", path, currentModel, key, err)
		}
		model := pricing.Models[currentModel]
		switch key {
		case "in_per_mtok":
			model.InPerMTok = parsed
		case "out_per_mtok":
			model.OutPerMTok = parsed
		}
		pricing.Models[currentModel] = model
	}
	if err := scanner.Err(); err != nil {
		return PricingConfig{}, fmt.Errorf("scan pricing config %q: %w", path, err)
	}
	if len(pricing.Models) == 0 {
		return PricingConfig{}, fmt.Errorf("parse pricing config %q: no models found", path)
	}

	return pricing, nil
}

func projectConfigContents(project ProjectConfig) string {
	return fmt.Sprintf(`config_version = 1

[project]
id = %q
issue_prefix = %q

[storage]
artifact_root = %q

[executor_limits]
max_output_bytes = %d
max_excerpt_bytes = %d

[maintenance]
abandoned_attempt_grace_seconds = %d
`, project.ProjectID, project.IssuePrefix, project.ArtifactRoot, project.MaxOutputBytes, project.MaxExcerptBytes, project.AbandonedAttemptGraceSeconds)
}

func pricingConfigContents(pricing PricingConfig) string {
	var builder strings.Builder
	builder.WriteString("# Default model pricing in USD per million tokens.\n")
	builder.WriteString("# cost_usd = (tokens_in * in_per_mtok + tokens_out * out_per_mtok) / 1e6\n\n")

	for _, name := range defaultPricingModelOrder {
		model, ok := pricing.Models[name]
		if !ok {
			continue
		}
		_, _ = fmt.Fprintf(&builder, "[models.%s]\n", name)
		_, _ = fmt.Fprintf(&builder, "in_per_mtok = %s\n", formatPricingValue(model.InPerMTok))
		_, _ = fmt.Fprintf(&builder, "out_per_mtok = %s\n\n", formatPricingValue(model.OutPerMTok))
	}

	return builder.String()
}

func formatPricingValue(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
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

func newProjectID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate project identity: %w", err)
	}
	return "project-" + hex.EncodeToString(bytes), nil
}

func isAllowedSection(section string) bool {
	return section == "project" || section == "storage" || section == "executor_limits" || section == "maintenance"
}

func isAllowedConfigKey(section, key string) bool {
	switch section {
	case "":
		return key == "config_version"
	case "project":
		return key == "id" || key == "issue_prefix"
	case "storage":
		return key == "artifact_root"
	case "executor_limits":
		return key == "max_output_bytes" || key == "max_excerpt_bytes"
	case "maintenance":
		return key == "abandoned_attempt_grace_seconds"
	default:
		return false
	}
}

func parseProjectConfigValues(values map[string]string) (ProjectConfig, error) {
	version, err := parseRequiredInt(values, "config_version", 1, 1)
	if err != nil || version != ConfigVersion {
		return ProjectConfig{}, fmt.Errorf("config_version must be the integer %d", ConfigVersion)
	}
	id, err := parseRequiredString(values, "project.id")
	if err != nil {
		return ProjectConfig{}, err
	}
	prefix, err := parseRequiredString(values, "project.issue_prefix")
	if err != nil {
		return ProjectConfig{}, err
	}
	artifactRoot, err := parseRequiredString(values, "storage.artifact_root")
	if err != nil {
		return ProjectConfig{}, err
	}
	maxOutput, err := parseRequiredInt(values, "executor_limits.max_output_bytes", 1, math.MaxInt64)
	if err != nil {
		return ProjectConfig{}, err
	}
	maxExcerpt, err := parseRequiredInt(values, "executor_limits.max_excerpt_bytes", 1, math.MaxInt64)
	if err != nil {
		return ProjectConfig{}, err
	}
	grace, err := parseRequiredInt(values, "maintenance.abandoned_attempt_grace_seconds", 0, math.MaxInt64)
	if err != nil {
		return ProjectConfig{}, err
	}
	return ProjectConfig{id, prefix, artifactRoot, maxOutput, maxExcerpt, grace}, nil
}

func parseRequiredString(values map[string]string, key string) (string, error) {
	value, ok := values[key]
	if !ok {
		return "", fmt.Errorf("missing %s", key)
	}
	parsed, err := parseQuotedString(value)
	if err != nil {
		return "", fmt.Errorf("%s: %w", key, err)
	}
	return parsed, nil
}

func parseRequiredInt(values map[string]string, key string, minimum, maximum int64) (int64, error) {
	value, ok := values[key]
	if !ok {
		return 0, fmt.Errorf("missing %s", key)
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed < minimum || parsed > maximum {
		return 0, fmt.Errorf("%s must be an integer from %d to %d", key, minimum, maximum)
	}
	return parsed, nil
}

func validateProjectConfig(project ProjectConfig, fidRoot string) error {
	if !isOpaqueProjectID(project.ProjectID) {
		return fmt.Errorf("project.id must be non-empty ASCII and at most 128 bytes")
	}
	if err := ValidateIssuePrefix(project.IssuePrefix); err != nil {
		return err
	}
	if err := validateArtifactRoot(project.ArtifactRoot, fidRoot); err != nil {
		return err
	}
	if project.MaxOutputBytes < 1 || project.MaxExcerptBytes < 1 || project.MaxExcerptBytes > project.MaxOutputBytes {
		return fmt.Errorf("executor_limits require positive values with max_excerpt_bytes no greater than max_output_bytes")
	}
	if project.AbandonedAttemptGraceSeconds < 0 {
		return fmt.Errorf("maintenance.abandoned_attempt_grace_seconds must be non-negative")
	}
	return nil
}

func isOpaqueProjectID(value string) bool {
	if len(value) == 0 || len(value) > 128 {
		return false
	}
	for _, character := range value {
		if character < 0x21 || character > 0x7e {
			return false
		}
	}
	return true
}

func isIssuePrefix(value string) bool {
	if len(value) < 1 || len(value) > 16 {
		return false
	}
	for index, character := range value {
		if character >= 'A' && character <= 'Z' {
			continue
		}
		if index > 0 && character >= '0' && character <= '9' {
			continue
		}
		return false
	}
	return true
}

func ValidateIssuePrefix(value string) error {
	if !isIssuePrefix(value) {
		return fmt.Errorf("project.issue_prefix must match ^[A-Z][A-Z0-9]{0,15}$")
	}
	return nil
}

func validateArtifactRoot(value, fidRoot string) error {
	if !utf8.ValidString(value) || len(value) == 0 || len(value) > 4096 || filepath.IsAbs(value) || value == "." || filepath.Clean(value) != value {
		return fmt.Errorf("storage.artifact_root must be a normalized project-relative path")
	}
	for _, part := range strings.Split(value, string(filepath.Separator)) {
		if part == ".." {
			return fmt.Errorf("storage.artifact_root must not escape .fid")
		}
	}
	path := fidRoot
	for _, part := range strings.Split(value, string(filepath.Separator)) {
		path = filepath.Join(path, part)
		info, err := os.Lstat(path)
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
			return fmt.Errorf("inspect storage.artifact_root: %w", err)
		}
		if info.Mode()&fs.ModeSymlink == 0 {
			continue
		}
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			return fmt.Errorf("resolve storage.artifact_root: %w", err)
		}
		relative, err := filepath.Rel(fidRoot, resolved)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return fmt.Errorf("storage.artifact_root symlink escapes .fid")
		}
	}
	return nil
}

func pathsForConfigRoot(workspaceRoot string) string {
	return filepath.Join(workspaceRoot, WorkspaceDirName)
}

func closeConfigFile(file *os.File) {
	_ = file.Close()
}
