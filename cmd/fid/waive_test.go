package main

import (
	"strings"
	"testing"
)

func TestWaiveRequiresInteractivePerson(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	before := reviewProjection(t, f)
	got := f.run("--json", "waive", "FID-9f13aa", "--lane", "routine", "--reason", "Reviewed risk", "--actor", "human:mark", "--idempotency-key", "waiver")
	if got.err == nil || !strings.Contains(got.stderr, "standard input and standard error must be terminals") || reviewProjection(t, f) != before {
		t.Fatalf("waive refusal = %#v", got)
	}
}
