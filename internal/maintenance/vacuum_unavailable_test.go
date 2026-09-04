package maintenance

import (
	"errors"
	"testing"
)

func TestVacuumIsUnavailableBeforeAnyWriterWork(t *testing.T) {
	_, err := Vacuum("", "", "")
	var unavailable *CapabilityUnavailableError
	if !errors.As(err, &unavailable) {
		t.Fatalf("Vacuum() error = %v, want CapabilityUnavailableError", err)
	}
	if unavailable.CapabilityID != "telemetry.trajectory.v1" || unavailable.OwnerPhase != 25 {
		t.Fatalf("unavailable = %#v, want telemetry.trajectory.v1 owned by phase 25", unavailable)
	}
}
