package config

import (
	"math"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

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
	paths, err := WriteBootstrapFiles(workspace, DefaultProjectConfig("FID"), defaultPricing)
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
