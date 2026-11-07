//go:build integration

package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
)

var integrationBinary string

const integrationBinaryEnv = "FID_INTEGRATION_BINARY"

func TestMain(m *testing.M) {
	exitCode := runIntegrationSuite(m)
	os.Exit(exitCode)
}

func runIntegrationSuite(m interface{ Run() int }) int {
	if err := gittest.Guard(); err != nil {
		fmt.Fprintf(os.Stderr, "guard git discovery: %v\n", err)
		return 1
	}
	configHome, err := os.MkdirTemp("", "fid-integration-config-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create isolated user config: %v\n", err)
		return 1
	}
	defer func() { _ = os.RemoveAll(configHome) }()
	if err := os.Setenv("XDG_CONFIG_HOME", configHome); err != nil {
		fmt.Fprintf(os.Stderr, "isolate user config: %v\n", err)
		return 1
	}
	// fid backup keeps its clone in the user cache, so no test may use ~/.cache.
	// Go's build cache defaults to a directory there, so pin it first: tests
	// that run go build keep the shared cache.
	if err := pinGoBuildCache(); err != nil {
		fmt.Fprintf(os.Stderr, "keep the Go build cache: %v\n", err)
		return 1
	}
	cacheHome, err := os.MkdirTemp("", "fid-integration-cache-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create isolated user cache: %v\n", err)
		return 1
	}
	defer func() { _ = os.RemoveAll(cacheHome) }()
	if err := os.Setenv("XDG_CACHE_HOME", cacheHome); err != nil {
		fmt.Fprintf(os.Stderr, "isolate user cache: %v\n", err)
		return 1
	}
	if binary := os.Getenv(integrationBinaryEnv); binary != "" {
		if info, err := os.Stat(binary); err == nil && !info.IsDir() {
			integrationBinary = binary
			return m.Run()
		}
	}
	pruneStaleIntegrationBinaries()
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
	if err := os.Setenv(integrationBinaryEnv, integrationBinary); err != nil {
		fmt.Fprintf(os.Stderr, "publish integration fid binary: %v\n", err)
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

func pruneStaleIntegrationBinaries() {
	entries, err := os.ReadDir(os.TempDir())
	if err != nil {
		return
	}
	cutoff := time.Now().Add(-24 * time.Hour)
	for _, entry := range entries {
		if !entry.IsDir() || len(entry.Name()) < len("fid-integration-binary-") || entry.Name()[:len("fid-integration-binary-")] != "fid-integration-binary-" {
			continue
		}
		info, err := entry.Info()
		if err == nil && info.ModTime().Before(cutoff) {
			_ = os.RemoveAll(filepath.Join(os.TempDir(), entry.Name()))
		}
	}
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

// pinGoBuildCache sets GOCACHE to Go's default under the real user cache,
// unless it is already set, before the suite isolates XDG_CACHE_HOME.
func pinGoBuildCache() error {
	if os.Getenv("GOCACHE") != "" {
		return nil
	}
	userCache, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	return os.Setenv("GOCACHE", filepath.Join(userCache, "go-build"))
}
