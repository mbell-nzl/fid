//go:build integration

package integration

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

//nolint:paralleltest // Each case changes XDG_CONFIG_HOME for the Go subprocess.
func TestIntegrationGoTelemetryConfigIsolation(t *testing.T) {
	for _, row := range []struct {
		name           string
		resolvedConfig string
		resolveErr     error
		wantRefusal    bool
	}{
		{"isolated-config", "fixture", nil, false},
		{"resolution-error", "", errors.New("config resolution failed"), true},
		{"resolution-error-with-path", "fixture", errors.New("config resolution failed"), true},
		{"outside-fixture", "user", nil, true},
	} {
		t.Run(row.name, func(t *testing.T) {
			fixture, user := t.TempDir(), t.TempDir()
			// Seed a simulated user's actual Go telemetry setting. Never point
			// the command at the developer's own config, even if a guard fails.
			modePath := filepath.Join(user, "go", "telemetry", "mode")
			const seeded = "local 2026-11-15"
			if err := os.MkdirAll(filepath.Dir(modePath), 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(modePath, []byte(seeded), 0o600); err != nil {
				t.Fatal(err)
			}
			configHome := fixture
			if row.wantRefusal {
				// Model a platform resolving Go's config outside the fixture.
				// An unguarded go telemetry off would overwrite the seed.
				configHome = user
			}
			t.Setenv("XDG_CONFIG_HOME", configHome)
			resolve := func() (string, error) {
				return map[string]string{"fixture": fixture, "user": user}[row.resolvedConfig], row.resolveErr
			}
			err := disableGoTelemetry(fixture, resolve)
			if row.wantRefusal {
				if err == nil || !strings.Contains(err.Error(), "refuse to change Go telemetry outside isolated config") {
					t.Errorf("telemetry setup = %v, want config isolation refusal", err)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				output, err := exec.Command("go", "telemetry").CombinedOutput()
				if err != nil || strings.TrimSpace(string(output)) != "off" {
					t.Errorf("isolated Go telemetry = %q, %v; want off", output, err)
				}
			}
			data, err := os.ReadFile(modePath)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != seeded {
				t.Errorf("simulated user telemetry changed to %q; want unchanged %q", data, seeded)
			}
		})
	}
}
