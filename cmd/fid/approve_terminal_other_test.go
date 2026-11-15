//go:build !linux

package main

import (
	"testing"
	"time"
)

// Keep V8's full state comparison on other platforms too. The Linux V6 journey
// verifies the real terminal boundary with the existing Linux PTY machinery.
func approvePublicForTest(t *testing.T, root, task, answer string, inputTTY, errorTTY, _ bool, extra ...string) commandExecutionResult {
	t.Helper()
	return approveForTest(&statusFixture{t: t, workspace: root, clock: time.Now}, task, answer, inputTTY && errorTTY, extra...)
}
