package leaktest_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/leaktest"
)

func TestAssertionRejectsKernelVocabulary(t *testing.T) {
	t.Parallel()
	for _, word := range []string{
		"Envelope", "envelopes", "Activation", "activations", "Obligation", "obligations",
		"Attention Item", "attention items",
		"ENV-123", "STP-123", "ATT-123", "env-123", "stp-123", "att-123",
	} {
		t.Run(word, func(t *testing.T) {
			t.Parallel()
			output := "error: " + word + " is unavailable\n"
			recorder := &testRecorder{}
			leaktest.AssertNoKernelWords(recorder, output)
			if !strings.Contains(recorder.failure, "human output names kernel words") || !strings.Contains(recorder.failure, output) {
				t.Fatalf("assertion failure = %q, want vocabulary diagnostic and original output", recorder.failure)
			}
		})
	}
}

func TestAssertionAcceptsPlainWordsAndWordBoundaries(t *testing.T) {
	t.Parallel()
	for _, output := range []string{
		"", "Task held here; check passed; ready for hand-off", "running reclaimed enveloped activationless", "HLD-123 CHK-123 EVT-123",
		// Ordinary fid English since the kernel went (Mark, 2026-11-23): a real flag, a check that runs, a retry.
		"Take the task in its builder's worktree: fid next --claim", "Claim it with fid next --claim; claims are holds",
		"fid check runs the lane's checks; 2 runs", "Run fid status", "two attempts failed; attempt again",
	} {
		t.Run(output, func(t *testing.T) {
			t.Parallel()
			recorder := &testRecorder{}
			leaktest.AssertNoKernelWords(recorder, output)
			if recorder.failure != "" {
				t.Fatal(recorder.failure)
			}
		})
	}
}

type testRecorder struct{ failure string }

func (*testRecorder) Helper() {}
func (r *testRecorder) Fatalf(format string, args ...any) {
	r.failure = fmt.Sprintf(format, args...)
}
