package config

import (
	"path/filepath"
	"testing"
)

func TestWriteBootstrapFilesRoundTripsDefaultConfig(t *testing.T) {
	t.Parallel()

	workspace := t.TempDir()
	paths, err := WriteBootstrapFiles(workspace, DefaultProjectConfig("FID"), DefaultPricingConfig())
	if err != nil {
		t.Fatalf("WriteBootstrapFiles() error = %v", err)
	}

	projectConfig, err := ReadProjectConfig(paths.ConfigPath)
	if err != nil {
		t.Fatalf("ReadProjectConfig() error = %v", err)
	}
	if projectConfig.Prefix != "FID" {
		t.Fatalf("project prefix = %q, want %q", projectConfig.Prefix, "FID")
	}

	pricingConfig, err := ReadPricingConfig(paths.PricingPath)
	if err != nil {
		t.Fatalf("ReadPricingConfig() error = %v", err)
	}
	if len(pricingConfig.Models) != 3 {
		t.Fatalf("pricing models = %d, want %d", len(pricingConfig.Models), 3)
	}
	if got := pricingConfig.Models["claude-sonnet-4-6"]; got.InPerMTok != 3.0 || got.OutPerMTok != 15.0 {
		t.Fatalf("claude-sonnet-4-6 pricing = %+v, want {InPerMTok:3 OutPerMTok:15}", got)
	}

	if paths.DatabasePath != filepath.Join(workspace, WorkspaceDirName, DatabaseFileName) {
		t.Fatalf("database path = %q, want %q", paths.DatabasePath, filepath.Join(workspace, WorkspaceDirName, DatabaseFileName))
	}
}
