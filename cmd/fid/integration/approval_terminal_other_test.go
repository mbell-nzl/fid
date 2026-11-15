//go:build integration && !linux

package integration

import "testing"

func approvalCorpusCommand(t *testing.T, _, _ string, _ ...string) (string, string, error) {
	t.Helper()
	t.Fatal("the interactive approval integration corpus requires Linux pseudo-terminal ioctls")
	return "", "", nil
}
