package config

import (
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestProjectConfigRoundTripsStrictTargetSchema(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), ConfigFileName)
	project := ProjectConfig{
		ProjectID:                    "opaque-project-identity",
		IssuePrefix:                  "FID",
		ArtifactRoot:                 "artifacts",
		MaxOutputBytes:               10485760,
		MaxExcerptBytes:              16384,
		AbandonedAttemptGraceSeconds: 300,
	}
	contents := projectConfigContents(project)
	for _, want := range []string{
		"config_version = 1",
		"[project]\nid = \"opaque-project-identity\"\nissue_prefix = \"FID\"",
		"[storage]\nartifact_root = \"artifacts\"",
		"[executor_limits]\nmax_output_bytes = 10485760\nmax_excerpt_bytes = 16384",
		"[maintenance]\nabandoned_attempt_grace_seconds = 300",
	} {
		if !strings.Contains(contents, want) {
			t.Fatalf("project config contents = %q, want %q", contents, want)
		}
	}
	if strings.Contains(contents, "\nprefix =") || strings.Contains(contents, "started_without_claim") {
		t.Fatalf("project config contents = %q, contains retired key", contents)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	roundTripped, err := ReadProjectConfig(path)
	if err != nil {
		t.Fatalf("ReadProjectConfig() error = %v", err)
	}
	if roundTripped != project {
		t.Fatalf("round-tripped project = %#v, want %#v", roundTripped, project)
	}
}

func TestReadProjectConfigRejectsUnknownDuplicateAliasAndInvalidInputs(t *testing.T) {
	t.Parallel()

	valid := `config_version = 1

[project]
id = "opaque-project-identity"
issue_prefix = "FID"

[storage]
artifact_root = "artifacts"

[executor_limits]
max_output_bytes = 10
max_excerpt_bytes = 10

[maintenance]
abandoned_attempt_grace_seconds = 0
`
	for name, contents := range map[string]string{
		"unknown key":         strings.Replace(valid, "issue_prefix", "prefix", 1),
		"duplicate key":       strings.Replace(valid, "issue_prefix = \"FID\"", "issue_prefix = \"FID\"\nissue_prefix = \"ABC\"", 1),
		"missing version":     strings.Replace(valid, "config_version = 1\n\n", "", 1),
		"unsupported version": strings.Replace(valid, "config_version = 1", "config_version = 2", 1),
		"escaping artifact":   strings.Replace(valid, "artifact_root = \"artifacts\"", "artifact_root = \"../artifacts\"", 1),
		"invalid prefix":      strings.Replace(valid, "issue_prefix = \"FID\"", "issue_prefix = \"fid\"", 1),
		"invalid limits":      strings.Replace(valid, "max_excerpt_bytes = 10", "max_excerpt_bytes = 11", 1),
		"negative grace":      strings.Replace(valid, "abandoned_attempt_grace_seconds = 0", "abandoned_attempt_grace_seconds = -1", 1),
	} {
		name, contents := name, contents
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), ConfigFileName)
			if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			_, err := ReadProjectConfig(path)
			if err == nil {
				t.Fatal("ReadProjectConfig() error = nil, want strict config rejection")
			}
		})
	}
}

func TestValidateIssuePrefixRejectsInvalidBasename(t *testing.T) {
	if err := ValidateIssuePrefix("INVALID_PREFIX"); err == nil {
		t.Fatal("ValidateIssuePrefix() error = nil, want invalid underscore rejection")
	}
}

func TestReadProjectConfigRejectsArtifactRootSymlinkEscape(t *testing.T) {
	fidRoot := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(fidRoot, "artifacts")); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}
	project := ProjectConfig{
		ProjectID:                    "opaque-project-identity",
		IssuePrefix:                  "FID",
		ArtifactRoot:                 "artifacts",
		MaxOutputBytes:               10,
		MaxExcerptBytes:              10,
		AbandonedAttemptGraceSeconds: 0,
	}
	path := filepath.Join(fidRoot, ConfigFileName)
	if err := os.WriteFile(path, []byte(projectConfigContents(project)), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := ReadProjectConfig(path); err == nil {
		t.Fatal("ReadProjectConfig() error = nil, want artifact-root symlink escape rejection")
	}
}

func TestDefaultPricingConfigIncludesAnthropicOpenAIAndGoogleModels(t *testing.T) {
	t.Parallel()

	pricing := DefaultPricingConfig()

	requiredModels := []string{
		"claude-sonnet-4-6",
		"gpt-5.4-mini",
		"gemini-2.5-flash-lite",
	}
	for _, model := range requiredModels {
		if _, ok := pricing.Models[model]; !ok {
			t.Fatalf("DefaultPricingConfig() missing model %q", model)
		}
	}
}

func TestPricingConfigContentsRoundTripsDefaultPricingWithoutPrecisionLoss(t *testing.T) {
	t.Parallel()

	pricing := DefaultPricingConfig()

	contents := pricingConfigContents(pricing)
	path := filepath.Join(t.TempDir(), PricingFileName)
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	roundTripped, err := ReadPricingConfig(path)
	if err != nil {
		t.Fatalf("ReadPricingConfig() error = %v", err)
	}

	if !reflect.DeepEqual(roundTripped.Models, pricing.Models) {
		t.Fatalf("round-tripped pricing = %#v, want %#v", roundTripped.Models, pricing.Models)
	}

	got := roundTripped.Models["gpt-5.4-mini"]
	if math.Float64bits(got.InPerMTok) != math.Float64bits(0.375) || math.Float64bits(got.OutPerMTok) != math.Float64bits(2.25) {
		t.Fatalf("gpt-5.4-mini pricing = %+v, want {InPerMTok:0.375 OutPerMTok:2.25}", got)
	}
}

func TestWriteBootstrapFilesRoundTripsDefaultConfig(t *testing.T) {
	t.Parallel()

	defaultPricing := DefaultPricingConfig()
	workspace := t.TempDir()
	project, err := NewDefaultProjectConfig("FID")
	if err != nil {
		t.Fatalf("NewDefaultProjectConfig() error = %v", err)
	}
	paths, err := WriteBootstrapFiles(workspace, project, defaultPricing)
	if err != nil {
		t.Fatalf("WriteBootstrapFiles() error = %v", err)
	}

	projectConfig, err := ReadProjectConfig(paths.ConfigPath)
	if err != nil {
		t.Fatalf("ReadProjectConfig() error = %v", err)
	}
	if projectConfig.IssuePrefix != "FID" {
		t.Fatalf("project issue prefix = %q, want %q", projectConfig.IssuePrefix, "FID")
	}

	pricingConfig, err := ReadPricingConfig(paths.PricingPath)
	if err != nil {
		t.Fatalf("ReadPricingConfig() error = %v", err)
	}
	if len(pricingConfig.Models) != len(defaultPricing.Models) {
		t.Fatalf("pricing models = %d, want %d", len(pricingConfig.Models), len(defaultPricing.Models))
	}
	if got := pricingConfig.Models["claude-sonnet-4-6"]; got.InPerMTok != 3.0 || got.OutPerMTok != 15.0 {
		t.Fatalf("claude-sonnet-4-6 pricing = %+v, want {InPerMTok:3 OutPerMTok:15}", got)
	}
	if got := pricingConfig.Models["gpt-5.4-mini"]; got.InPerMTok != 0.375 || got.OutPerMTok != 2.25 {
		t.Fatalf("gpt-5.4-mini pricing = %+v, want {InPerMTok:0.375 OutPerMTok:2.25}", got)
	}

	if paths.DatabasePath != filepath.Join(workspace, WorkspaceDirName, DatabaseFileName) {
		t.Fatalf("database path = %q, want %q", paths.DatabasePath, filepath.Join(workspace, WorkspaceDirName, DatabaseFileName))
	}
}
