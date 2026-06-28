package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	WorkspaceDirName = ".fid"
	JournalDirName   = "journal"
	ConfigFileName   = "config.toml"
	PricingFileName  = "pricing.toml"
	DatabaseFileName = "fid.sqlite"
	GitIgnoreName    = ".gitignore"
)

type ProjectConfig struct {
	Prefix string
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

func DefaultProjectConfig(prefix string) ProjectConfig {
	return ProjectConfig{Prefix: prefix}
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
	if strings.TrimSpace(project.Prefix) == "" {
		return WorkspacePaths{}, fmt.Errorf("write bootstrap files: project prefix is required")
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
	prefix := ""
	for scanner.Scan() {
		line := normalizeLine(scanner.Text())
		if line == "" {
			continue
		}
		if isSection(line) {
			section = strings.Trim(line, "[]")
			continue
		}
		key, value, ok := splitAssignment(line)
		if !ok {
			continue
		}
		if section == "project" && key == "prefix" {
			parsed, err := parseQuotedString(value)
			if err != nil {
				return ProjectConfig{}, fmt.Errorf("parse project config %q: %w", path, err)
			}
			prefix = parsed
		}
	}
	if err := scanner.Err(); err != nil {
		return ProjectConfig{}, fmt.Errorf("scan project config %q: %w", path, err)
	}
	if prefix == "" {
		return ProjectConfig{}, fmt.Errorf("parse project config %q: missing [project].prefix", path)
	}

	return ProjectConfig{Prefix: prefix}, nil
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
	return fmt.Sprintf(`# Project identity
[project]
prefix = %q
`, project.Prefix)
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
	if comment := strings.Index(trimmed, "#"); comment >= 0 {
		trimmed = strings.TrimSpace(trimmed[:comment])
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
