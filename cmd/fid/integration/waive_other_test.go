//go:build integration && !linux

package integration

import "testing"

func waiverCodePair(string) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		t.Fatal("the interactive waiver integration corpus requires Linux pseudo-terminal ioctls")
	}
}
