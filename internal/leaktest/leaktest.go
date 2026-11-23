// Package leaktest checks human output for kernel vocabulary (RESET V2).
// Only test files import it; JSON and verbose diagnostics are outside its scope.
package leaktest

import "regexp"

// Every human-output test shares this one pattern. It names only words that
// have no meaning in fid since the kernel went; "claim", "run" and "attempt"
// are ordinary fid English (`fid next --claim`, a check that runs) and were
// dropped on Mark's decision of 2026-11-23.
var kernelWords = regexp.MustCompile(`(?i)\b(envelopes?|activations?|obligations?|attention items?)\b|ENV-|STP-|ATT-`)

type testingT interface {
	Helper()
	Fatalf(string, ...any)
}

// AssertNoKernelWords fails the test if output names a kernel noun or ID.
func AssertNoKernelWords(t testingT, output string) {
	t.Helper()
	if found := kernelWords.FindAllString(output, -1); len(found) > 0 {
		t.Fatalf("human output names kernel words %q:\n%s", found, output)
	}
}
