package evidence

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestSameGateRepositoryMatchesEqualOptionalValues(t *testing.T) {
	baseOID := "base-oid"
	dirtyDigest := "sha256:dirty"
	left := Repository{
		BaseOID:       &baseOID,
		CaptureSource: "git",
		Dirty:         true,
		DirtyDigest:   &dirtyDigest,
		HeadOID:       "head-oid",
		ProjectID:     "PRJ-test",
		TreeOID:       "tree-oid",
		WorktreeID:    "WT-test",
	}
	baseOIDCopy := string([]byte(baseOID))
	dirtyDigestCopy := string([]byte(dirtyDigest))
	right := left
	right.BaseOID = &baseOIDCopy
	right.DirtyDigest = &dirtyDigestCopy

	if !sameGateRepository(left, right) {
		t.Fatal("sameGateRepository rejected equal optional repository values")
	}
}

func TestFIDBudgetChargesDeriveOnlyObservedAttemptAndReceiptFacts(t *testing.T) {
	budgets := []workflow.Budget{
		{Name: "attempts", Unit: "attempts"},
		{Name: "failures", Unit: "failures"},
		{Name: "milliseconds", Unit: "milliseconds"},
		{Name: "tool-failures", Unit: "tool-failures"},
		{Name: "tokens", Unit: "input-tokens"},
	}
	attempt := Attempt{ID: "ATT-1", StartedAt: "2026-09-06T00:00:00Z"}
	receipt := Receipt{ID: "RCP-1", AttemptID: attempt.ID, Outcome: "failed", StartedAt: attempt.StartedAt, EndedAt: "2026-09-06T00:00:00.125Z"}

	charges, err := fidBudgetCharges(budgets, attempt, &receipt)
	if err != nil {
		t.Fatalf("fidBudgetCharges() error = %v", err)
	}
	want := map[string]int64{"attempts": 1, "failures": 1, "milliseconds": 125}
	if len(charges) != len(want) {
		t.Fatalf("fidBudgetCharges() = %#v, want %d observed charges", charges, len(want))
	}
	for _, charge := range charges {
		if got, ok := want[charge.BudgetName]; !ok || charge.Value == nil || *charge.Value != got || charge.Source != "fid" || charge.SampleKind != "delta" {
			t.Fatalf("charge = %#v, want fid observed charge in %#v", charge, want)
		}
	}
}
