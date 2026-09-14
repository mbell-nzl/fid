package export

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/config"
)

type archiveProjectConfig struct {
	ConfigVersion int `json:"config_version"`
	Project       struct {
		ID          string `json:"id"`
		IssuePrefix string `json:"issue_prefix"`
	} `json:"project"`
	Storage struct {
		ArtifactRoot string `json:"artifact_root"`
	} `json:"storage"`
	ExecutorLimits struct {
		MaxOutputBytes  int64 `json:"max_output_bytes"`
		MaxExcerptBytes int64 `json:"max_excerpt_bytes"`
	} `json:"executor_limits"`
	Maintenance struct {
		AbandonedAttemptGraceSeconds int64 `json:"abandoned_attempt_grace_seconds"`
	} `json:"maintenance"`
}

func archiveConfig(project config.ProjectConfig) *archiveProjectConfig {
	c := &archiveProjectConfig{ConfigVersion: config.ConfigVersion}
	c.Project.ID, c.Project.IssuePrefix = project.ProjectID, project.IssuePrefix
	c.Storage.ArtifactRoot = project.ArtifactRoot
	c.ExecutorLimits.MaxOutputBytes, c.ExecutorLimits.MaxExcerptBytes = project.MaxOutputBytes, project.MaxExcerptBytes
	c.Maintenance.AbandonedAttemptGraceSeconds = project.AbandonedAttemptGraceSeconds
	return c
}

func (c *archiveProjectConfig) projectConfig(projectID string) (config.ProjectConfig, error) {
	if c == nil {
		return config.ProjectConfig{}, fmt.Errorf("project configuration is required")
	}
	p := config.ProjectConfig{ProjectID: c.Project.ID, IssuePrefix: c.Project.IssuePrefix, ArtifactRoot: c.Storage.ArtifactRoot, MaxOutputBytes: c.ExecutorLimits.MaxOutputBytes, MaxExcerptBytes: c.ExecutorLimits.MaxExcerptBytes, AbandonedAttemptGraceSeconds: c.Maintenance.AbandonedAttemptGraceSeconds}
	if c.ConfigVersion != config.ConfigVersion || p.ProjectID != projectID || !safeArchivePath(p.ArtifactRoot) {
		return config.ProjectConfig{}, fmt.Errorf("configuration identity, version or Artifact root is invalid")
	}
	if err := config.ValidatePortableProjectConfig(p); err != nil {
		return config.ProjectConfig{}, err
	}
	return p, nil
}
