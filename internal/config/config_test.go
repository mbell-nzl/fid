package config

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
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

func TestParseProjectConfigBoundaryDomains(t *testing.T) {
	t.Parallel()
	valid := strictProjectConfig(ProjectConfig{
		ProjectID:                    "opaque-project-identity",
		IssuePrefix:                  "FID",
		ArtifactRoot:                 "artifacts",
		MaxOutputBytes:               10,
		MaxExcerptBytes:              10,
		AbandonedAttemptGraceSeconds: 0,
	})

	for _, test := range []struct {
		name     string
		contents string
		want     ProjectConfig
		valid    bool
	}{
		{name: "exact schema", contents: valid, want: ProjectConfig{"opaque-project-identity", "FID", "artifacts", 10, 10, 0}, valid: true},
		{name: "missing version", contents: strings.Replace(valid, "config_version = 1\n\n", "", 1)},
		{name: "unsupported version", contents: strings.Replace(valid, "config_version = 1", "config_version = 2", 1)},
		{name: "non-integer version", contents: strings.Replace(valid, "config_version = 1", "config_version = 1.0", 1)},
		{name: "unknown top level key", contents: strings.Replace(valid, "config_version = 1", "config_version = 1\nunknown = 1", 1)},
		{name: "unknown section", contents: valid + "\n[unknown]\nvalue = 1\n"},
		{name: "unknown project key", contents: strings.Replace(valid, "issue_prefix = \"FID\"", "issue_prefix = \"FID\"\nunknown = \"value\"", 1)},
		{name: "unknown storage key", contents: strings.Replace(valid, "artifact_root = \"artifacts\"", "artifact_root = \"artifacts\"\nunknown = \"value\"", 1)},
		{name: "unknown executor limits key", contents: strings.Replace(valid, "max_excerpt_bytes = 10", "max_excerpt_bytes = 10\nunknown = 1", 1)},
		{name: "unknown maintenance key", contents: strings.Replace(valid, "abandoned_attempt_grace_seconds = 0", "abandoned_attempt_grace_seconds = 0\nunknown = 1", 1)},
		{name: "duplicate key", contents: strings.Replace(valid, "issue_prefix = \"FID\"", "issue_prefix = \"FID\"\nissue_prefix = \"ABC\"", 1)},
		{name: "duplicate section", contents: valid + "\n[project]\nid = \"another\"\nissue_prefix = \"ABC\"\n"},
		{name: "missing required section", contents: strings.Replace(valid, "\n[maintenance]\nabandoned_attempt_grace_seconds = 0\n", "\n", 1)},
		{name: "missing required key", contents: strings.Replace(valid, "id = \"opaque-project-identity\"\n", "", 1)},
		{name: "project id one byte", contents: strings.Replace(valid, "opaque-project-identity", "x", 1), want: ProjectConfig{"x", "FID", "artifacts", 10, 10, 0}, valid: true},
		{name: "project id 128 bytes", contents: strings.Replace(valid, "opaque-project-identity", strings.Repeat("x", 128), 1), want: ProjectConfig{strings.Repeat("x", 128), "FID", "artifacts", 10, 10, 0}, valid: true},
		{name: "project id empty", contents: strings.Replace(valid, "opaque-project-identity", "", 1)},
		{name: "project id 129 bytes", contents: strings.Replace(valid, "opaque-project-identity", strings.Repeat("x", 129), 1)},
		{name: "project id non-ASCII", contents: strings.Replace(valid, "opaque-project-identity", "é", 1)},
		{name: "project id control character", contents: strings.Replace(valid, "opaque-project-identity", "id\x01", 1)},
		{name: "prefix one character", contents: strings.Replace(valid, "\"FID\"", "\"A\"", 1), want: ProjectConfig{"opaque-project-identity", "A", "artifacts", 10, 10, 0}, valid: true},
		{name: "prefix 16 characters", contents: strings.Replace(valid, "\"FID\"", "\"A123456789012345\"", 1), want: ProjectConfig{"opaque-project-identity", "A123456789012345", "artifacts", 10, 10, 0}, valid: true},
		{name: "prefix empty", contents: strings.Replace(valid, "\"FID\"", "\"\"", 1)},
		{name: "prefix 17 characters", contents: strings.Replace(valid, "\"FID\"", "\"A1234567890123456\"", 1)},
		{name: "prefix lowercase", contents: strings.Replace(valid, "\"FID\"", "\"fid\"", 1)},
		{name: "prefix leading digit", contents: strings.Replace(valid, "\"FID\"", "\"1FID\"", 1)},
		{name: "prefix hyphen", contents: strings.Replace(valid, "\"FID\"", "\"FI-D\"", 1)},
		{name: "prefix underscore", contents: strings.Replace(valid, "\"FID\"", "\"FI_D\"", 1)},
		{name: "prefix non-ASCII", contents: strings.Replace(valid, "\"FID\"", "\"FÍD\"", 1)},
		{name: "legacy prefix key", contents: strings.Replace(valid, "issue_prefix", "prefix", 1)},
		{name: "artifact root one byte", contents: strings.Replace(valid, "\"artifacts\"", "\"a\"", 1), want: ProjectConfig{"opaque-project-identity", "FID", "a", 10, 10, 0}, valid: true},
		{name: "artifact root 4096 bytes", contents: strings.Replace(valid, "\"artifacts\"", "\""+maxArtifactRoot()+"\"", 1), want: ProjectConfig{"opaque-project-identity", "FID", maxArtifactRoot(), 10, 10, 0}, valid: true},
		{name: "artifact root empty", contents: strings.Replace(valid, "\"artifacts\"", "\"\"", 1)},
		{name: "artifact root 4097 bytes", contents: strings.Replace(valid, "\"artifacts\"", "\""+maxArtifactRoot()+"a\"", 1)},
		{name: "artifact root absolute", contents: strings.Replace(valid, "\"artifacts\"", "\"/artifacts\"", 1)},
		{name: "artifact root dot", contents: strings.Replace(valid, "\"artifacts\"", "\".\"", 1)},
		{name: "artifact root lexical escape", contents: strings.Replace(valid, "\"artifacts\"", "\"../artifacts\"", 1)},
		{name: "artifact root invalid UTF-8", contents: strings.Replace(valid, "\"artifacts\"", "\"\xff\"", 1)},
		{name: "executor limits minimum", contents: strings.Replace(strings.Replace(valid, "max_output_bytes = 10", "max_output_bytes = 1", 1), "max_excerpt_bytes = 10", "max_excerpt_bytes = 1", 1), want: ProjectConfig{"opaque-project-identity", "FID", "artifacts", 1, 1, 0}, valid: true},
		{name: "executor limits maximum", contents: strings.Replace(strings.Replace(valid, "max_output_bytes = 10", "max_output_bytes = 9223372036854775807", 1), "max_excerpt_bytes = 10", "max_excerpt_bytes = 9223372036854775807", 1), want: ProjectConfig{"opaque-project-identity", "FID", "artifacts", math.MaxInt64, math.MaxInt64, 0}, valid: true},
		{name: "executor output zero", contents: strings.Replace(valid, "max_output_bytes = 10", "max_output_bytes = 0", 1)},
		{name: "executor output negative", contents: strings.Replace(valid, "max_output_bytes = 10", "max_output_bytes = -1", 1)},
		{name: "executor output overflow", contents: strings.Replace(valid, "max_output_bytes = 10", "max_output_bytes = 9223372036854775808", 1)},
		{name: "executor output float", contents: strings.Replace(valid, "max_output_bytes = 10", "max_output_bytes = 1.0", 1)},
		{name: "executor output quoted", contents: strings.Replace(valid, "max_output_bytes = 10", "max_output_bytes = \"1\"", 1)},
		{name: "executor excerpt zero", contents: strings.Replace(valid, "max_excerpt_bytes = 10", "max_excerpt_bytes = 0", 1)},
		{name: "executor excerpt negative", contents: strings.Replace(valid, "max_excerpt_bytes = 10", "max_excerpt_bytes = -1", 1)},
		{name: "executor excerpt overflow", contents: strings.Replace(valid, "max_excerpt_bytes = 10", "max_excerpt_bytes = 9223372036854775808", 1)},
		{name: "executor excerpt float", contents: strings.Replace(valid, "max_excerpt_bytes = 10", "max_excerpt_bytes = 1.0", 1)},
		{name: "executor excerpt quoted", contents: strings.Replace(valid, "max_excerpt_bytes = 10", "max_excerpt_bytes = \"1\"", 1)},
		{name: "executor excerpt exceeds output", contents: strings.Replace(valid, "max_excerpt_bytes = 10", "max_excerpt_bytes = 11", 1)},
		{name: "maintenance minimum", contents: valid, want: ProjectConfig{"opaque-project-identity", "FID", "artifacts", 10, 10, 0}, valid: true},
		{name: "maintenance maximum", contents: strings.Replace(valid, "abandoned_attempt_grace_seconds = 0", "abandoned_attempt_grace_seconds = 9223372036854775807", 1), want: ProjectConfig{"opaque-project-identity", "FID", "artifacts", 10, 10, math.MaxInt64}, valid: true},
		{name: "maintenance negative", contents: strings.Replace(valid, "abandoned_attempt_grace_seconds = 0", "abandoned_attempt_grace_seconds = -1", 1)},
		{name: "maintenance overflow", contents: strings.Replace(valid, "abandoned_attempt_grace_seconds = 0", "abandoned_attempt_grace_seconds = 9223372036854775808", 1)},
		{name: "maintenance non-integer", contents: strings.Replace(valid, "abandoned_attempt_grace_seconds = 0", "abandoned_attempt_grace_seconds = 1.0", 1)},
		{name: "retired maintenance key", contents: strings.Replace(valid, "abandoned_attempt_grace_seconds", "started_without_claim_grace_seconds", 1)},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ConfigFileName)
			if err := os.WriteFile(path, []byte(test.contents), 0o644); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}
			got, err := ReadProjectConfig(path)
			if test.valid {
				if err != nil {
					t.Fatalf("ReadProjectConfig() error = %v", err)
				}
				if got != test.want {
					t.Fatalf("ReadProjectConfig() = %#v, want %#v", got, test.want)
				}
				return
			}
			if err == nil {
				t.Fatalf("ReadProjectConfig() = %#v, want rejection", got)
			}
		})
	}
}

func TestProjectConfigWriterRoundTripsByteStableAtInt64Boundaries(t *testing.T) {
	t.Parallel()
	for _, project := range []ProjectConfig{
		{ProjectID: "opaque-project-identity", IssuePrefix: "FID", ArtifactRoot: "artifacts", MaxOutputBytes: 10, MaxExcerptBytes: 10, AbandonedAttemptGraceSeconds: 0},
		{ProjectID: strings.Repeat("x", 128), IssuePrefix: "A123456789012345", ArtifactRoot: maxArtifactRoot(), MaxOutputBytes: math.MaxInt64, MaxExcerptBytes: math.MaxInt64, AbandonedAttemptGraceSeconds: math.MaxInt64},
	} {
		project := project
		t.Run(project.IssuePrefix, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ConfigFileName)
			want := strictProjectConfig(project)
			if err := os.WriteFile(path, []byte(want), 0o644); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}
			got, err := ReadProjectConfig(path)
			if err != nil {
				t.Fatalf("ReadProjectConfig() error = %v", err)
			}
			if got != project {
				t.Fatalf("ReadProjectConfig() = %#v, want %#v", got, project)
			}
			if roundTrip := strictProjectConfig(got); roundTrip != want {
				t.Fatalf("round-trip contents differ\ngot:  %q\nwant: %q", roundTrip, want)
			}
		})
	}
}

func TestReadProjectConfigRejectsArtifactRootSymlinkComponentEscape(t *testing.T) {
	t.Parallel()
	fidRoot := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(fidRoot, "nested")); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}
	path := filepath.Join(fidRoot, ConfigFileName)
	if err := os.WriteFile(path, []byte(strictProjectConfig(ProjectConfig{
		ProjectID: "opaque-project-identity", IssuePrefix: "FID", ArtifactRoot: "nested/artifacts", MaxOutputBytes: 1, MaxExcerptBytes: 1,
	})), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if _, err := ReadProjectConfig(path); err == nil {
		t.Fatal("ReadProjectConfig() error = nil, want symlink escape rejection")
	}
}

func strictProjectConfig(project ProjectConfig) string {
	return projectConfigContents(project)
}

func maxArtifactRoot() string {
	return strings.Repeat(strings.Repeat("a", 240)+"/", 16) + strings.Repeat("a", 240)
}

func TestValidateIssuePrefixRejectsInvalidBasename(t *testing.T) {
	t.Parallel()
	if err := ValidateIssuePrefix("INVALID_PREFIX"); err == nil {
		t.Fatal("ValidateIssuePrefix() error = nil, want invalid underscore rejection")
	}
}

func TestReadProjectConfigRejectsArtifactRootSymlinkEscape(t *testing.T) {
	t.Parallel()
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

func TestWriteLedgerFilesRoundTripsDefaultConfig(t *testing.T) {
	t.Parallel()

	workspace := gittest.NewWorkspace(t)
	project, err := NewDefaultProjectConfig("FID")
	if err != nil {
		t.Fatalf("NewDefaultProjectConfig() error = %v", err)
	}
	paths, err := WriteLedgerFiles(LedgerDir(workspace), project)
	if err != nil {
		t.Fatalf("WriteLedgerFiles() error = %v", err)
	}

	projectConfig, err := ReadProjectConfig(paths.ConfigPath)
	if err != nil {
		t.Fatalf("ReadProjectConfig() error = %v", err)
	}
	if projectConfig.IssuePrefix != "FID" {
		t.Fatalf("project issue prefix = %q, want %q", projectConfig.IssuePrefix, "FID")
	}

	if paths.DatabasePath != filepath.Join(workspace, ".git", LedgerDirName, DatabaseFileName) {
		t.Fatalf("database path = %q, want %q", paths.DatabasePath, filepath.Join(workspace, ".git", LedgerDirName, DatabaseFileName))
	}
}
