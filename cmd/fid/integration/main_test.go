//go:build integration

package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var integrationBinary string

func TestMain(m *testing.M) {
	exitCode := runIntegrationSuite(m)
	os.Exit(exitCode)
}

func runIntegrationSuite(m interface{ Run() int }) int {
	binaryRoot, err := os.MkdirTemp("", "fid-integration-binary-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create integration binary directory: %v\n", err)
		return 1
	}

	binaryName := "fid"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	integrationBinary = filepath.Join(binaryRoot, binaryName)
	command := exec.Command("go", "build", "-o", integrationBinary, "./cmd/fid")
	command.Dir = integrationRepositoryRoot()
	if output, err := command.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build integration fid binary: %v\n%s", err, output)
		_ = os.RemoveAll(binaryRoot)
		return 1
	}

	exitCode := m.Run()
	if err := os.RemoveAll(binaryRoot); err != nil {
		fmt.Fprintf(os.Stderr, "remove integration binary directory: %v\n", err)
		if exitCode == 0 {
			exitCode = 1
		}
	}
	return exitCode
}

func integrationRepositoryRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", ".."))
}

func integrationFIDBinary(t *testing.T) string {
	t.Helper()
	if integrationBinary == "" {
		t.Fatal("integration binary was not initialized")
	}
	if _, err := os.Stat(integrationBinary); err != nil {
		t.Fatalf("stat integration binary: %v", err)
	}
	return integrationBinary
}
