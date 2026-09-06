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

	startCharges, err := fidBudgetCharges(budgets, attempt, nil)
	if err != nil {
		t.Fatalf("fidBudgetCharges(start) error = %v", err)
	}
	if len(startCharges) != 1 || startCharges[0].BudgetName != "attempts" || startCharges[0].Value == nil || *startCharges[0].Value != 1 {
		t.Fatalf("fidBudgetCharges(start) = %#v, want one observed attempt charge", startCharges)
	}
	charges, err := fidBudgetCharges(budgets, attempt, &receipt)
	if err != nil {
		t.Fatalf("fidBudgetCharges(receipt) error = %v", err)
	}
	want := map[string]int64{"failures": 1, "milliseconds": 125}
	if len(charges) != len(want) {
		t.Fatalf("fidBudgetCharges() = %#v, want %d observed charges", charges, len(want))
	}
	for _, charge := range charges {
		if got, ok := want[charge.BudgetName]; !ok || charge.Value == nil || *charge.Value != got || charge.Source != "fid" || charge.SampleKind != "delta" {
			t.Fatalf("charge = %#v, want fid observed charge in %#v", charge, want)
		}
	}
}

func TestBudgetAdmissionDeniesUnknownAndExhaustedChainsBeforeAttempt(t *testing.T) {
	known := budgetChain{Budget: workflow.Budget{Name: "attempts", Unit: "attempts", Maximum: 1}}
	if blocks, err := budgetAdmissionBlocks([]budgetChain{known}); err != nil || len(blocks) != 0 {
		t.Fatalf("budgetAdmissionBlocks() known empty chain = %#v, %v; want admitted", blocks, err)
	}
	unknown := budgetChain{Budget: workflow.Budget{Name: "input", Unit: "input-tokens", Maximum: 10}}
	exhausted := budgetChain{Budget: workflow.Budget{Name: "attempts", Unit: "attempts", Maximum: 1}, Samples: []workflow.BudgetSample{{ID: "ATT-1", Source: "fid", SourceSampleID: "ATT-1", SampleKind: "delta", Value: int64Pointer(1)}}}

	blocks, err := budgetAdmissionBlocks([]budgetChain{unknown, exhausted})
	if err != nil {
		t.Fatalf("budgetAdmissionBlocks() error = %v", err)
	}
	if len(blocks) != 2 || blocks[0].Code != "budget-data-unavailable" || blocks[1].Code != "budget-exhausted" {
		t.Fatalf("budgetAdmissionBlocks() = %#v, want unknown and exhausted denial blocks", blocks)
	}
}

func TestBudgetBoundaryBlocksOnlyOnFirstExhaustion(t *testing.T) {
	available := workflow.Budget{Name: "attempts", Unit: "attempts", Maximum: 2, Consumed: int64Pointer(1), Remaining: int64Pointer(1)}
	exhausted := available
	exhausted.Consumed, exhausted.Remaining = int64Pointer(2), int64Pointer(0)

	if blocks := budgetBoundaryBlocks(available, exhausted); len(blocks) != 1 || blocks[0].Code != "budget-exhausted" {
		t.Fatalf("budgetBoundaryBlocks() = %#v, want one exhaustion block", blocks)
	}
	if blocks := budgetBoundaryBlocks(exhausted, exhausted); len(blocks) != 0 {
		t.Fatalf("budgetBoundaryBlocks() after existing exhaustion = %#v, want no duplicate Stop", blocks)
	}
}

func int64Pointer(value int64) *int64 { return &value }
