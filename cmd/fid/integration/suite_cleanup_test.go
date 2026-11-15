//go:build integration

package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestIntegrationGoTelemetryIsOff(t *testing.T) {
	t.Parallel()
	output, err := exec.Command("go", "env", "GOTELEMETRY").CombinedOutput()
	if err != nil || strings.TrimSpace(string(output)) != "off" {
		t.Fatalf("go env GOTELEMETRY = %q, %v; want off in the suite's isolated config", output, err)
	}
}

func TestIntegrationSuiteRemovesTemporaryDirectories(t *testing.T) {
	t.Parallel()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	scratch := t.TempDir()
	command := exec.CommandContext(t.Context(), binary, "-test.run=^TestIntegrationBinaryOwnedByPackage$")
	// Clear the inherited binary so this also exercises the suite's go build.
	command.Env = append(os.Environ(), "TMPDIR="+scratch, integrationBinaryEnv+"=", "GOCOVERDIR="+t.TempDir())
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("child integration suite: %v\n%s", err, output)
	}
	deadline := time.Now().Add(250 * time.Millisecond)
	for {
		entries, err := os.ReadDir(scratch)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 0 {
			t.Fatalf("temporary directory has %d leaked entries after suite exit: %v", len(entries), entries)
		}
		if time.Now().After(deadline) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}
