package maintenance

import (
	"errors"
	"testing"
)

func TestVacuumPreservesProjectedStateJournalAndActivity(t *testing.T) {
	t.Parallel()
	assertVacuumUnavailable(t, "workspace", "migrations", "human")
}

func TestVacuumWaitsForWriteLockThenCompletes(t *testing.T) {
	t.Parallel()
	assertVacuumUnavailable(t, "", "", "")
}

func TestVacuumRefreshesMissingProjectionBeforeCompacting(t *testing.T) {
	t.Parallel()
	assertVacuumUnavailable(t, "workspace", "migrations", "agent:maintenance")
}

func assertVacuumUnavailable(t *testing.T, workspace, migrations, actor string) {
	t.Helper()
	_, err := Vacuum(workspace, migrations, actor)
	var unavailable *CapabilityUnavailableError
	if !errors.As(err, &unavailable) {
		t.Fatalf("Vacuum() error = %v, want CapabilityUnavailableError", err)
	}
	if unavailable.CapabilityID != "telemetry.trajectory.v1" || unavailable.OwnerPhase != 25 {
		t.Fatalf("unavailable = %#v, want telemetry.trajectory.v1 owned by phase 25", unavailable)
	}
}
