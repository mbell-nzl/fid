package issues

import "testing"

func TestBuildIssueRevisionEffectsSortsObservedRevisionsAndDerivesResults(t *testing.T) {
	effects, err := BuildIssueRevisionEffects([]IssueRevision{
		{IssueID: "FID-z", Observed: 7},
		{IssueID: "FID-a", Observed: 1},
	})
	if err != nil {
		t.Fatalf("BuildIssueRevisionEffects() error = %v", err)
	}
	if len(effects) != 2 {
		t.Fatalf("revision effects = %d, want 2", len(effects))
	}
	for index, want := range []struct {
		id          string
		prior, next int
	}{
		{id: "FID-a", prior: 1, next: 2},
		{id: "FID-z", prior: 7, next: 8},
	} {
		effect := effects[index]
		if got := effect.PK["id"]; got != want.id {
			t.Fatalf("effect %d issue id = %#v, want %q", index, got, want.id)
		}
		if got := effect.Values["revision"]; got != want.next {
			t.Fatalf("effect %d revision value = %#v, want %d", index, got, want.next)
		}
		if got := effect.Prior["revision"]; got != want.prior {
			t.Fatalf("effect %d revision prior = %#v, want %d", index, got, want.prior)
		}
		if got := effect.Preconditions["revision"]; got != want.prior {
			t.Fatalf("effect %d revision precondition = %#v, want %d", index, got, want.prior)
		}
	}
}
