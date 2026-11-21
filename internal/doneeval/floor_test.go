package doneeval

import (
	"testing"
)

func TestDoneRefusesHandoffBelowItsRecordedPolicyFloor(t *testing.T) {
	t.Parallel()
	for _, path := range []string{".fid.toml", ".github/workflows/check.yml", "Makefile", "internal/trust/auth.go"} {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("standard", false)
			s.Handoff.ChangedPaths = []string{path}
			lane := s.Policy.Lanes["assured"]
			lane.Paths = []string{"internal/trust/**"}
			s.Policy.Lanes["assured"] = lane
			check := s.Policy.Checks["test"]
			check.DefinedBy = []string{"Makefile"}
			s.Policy.Checks["test"] = check
			got := Evaluate(s)
			if got.Status != Refused || got.Gates[PolicyGate].Reason != PolicyMismatch {
				t.Fatalf("below-floor hand-off accepted: %+v", got)
			}
		})
	}
}

func TestHistoricalHandoffDoesNotAcquireLaterFloorSemantics(t *testing.T) {
	t.Parallel()
	s := completeSnapshot("standard", false)
	s.Policy.RulesVersion = 2
	s.Handoff.ChangedPaths = []string{".fid.toml"}
	if got := Evaluate(s); got.Status != Complete {
		t.Fatalf("later floor changed historical gates: %+v", got)
	}
}

func TestDoneRefusesAnUncomputablePolicyFloor(t *testing.T) {
	t.Parallel()
	for _, source := range []string{"hard paths", "assured paths", "routine paths", "check definition"} {
		t.Run(source, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("assured", false)
			switch source {
			case "hard paths":
				s.Policy.HardAssured = []string{"***"}
			case "assured paths", "routine paths":
				name := "assured"
				if source == "routine paths" {
					name = "routine"
				}
				lane := s.Policy.Lanes[name]
				lane.Paths = []string{"***"}
				s.Policy.Lanes[name] = lane
			case "check definition":
				check := s.Policy.Checks["test"]
				check.DefinedBy = []string{"***"}
				s.Policy.Checks["test"] = check
			}
			got := Evaluate(s)
			if got.Status != Refused || got.Gates[PolicyGate].Reason != PolicyMismatch {
				t.Fatalf("uncomputable floor accepted: %+v", got)
			}
		})
	}
}
