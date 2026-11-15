// Package leaktest checks human output for kernel vocabulary (RESET V2).
// Only test files import it; JSON and verbose diagnostics are outside its scope.
package leaktest

import "regexp"

// Keep the established pattern unchanged: every human-output test shares it.
var kernelWords = regexp.MustCompile(`(?i)\b(envelopes?|activations?|obligations?|attempts?|attention items?|runs?|claims?)\b|ENV-|STP-|ATT-`)

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
