//go:build integration

package integration

// These assertions cover retiring kernel/historical behavior; retain until its owning cut.
import (
	"testing"
)

func p24SingleLifecycle(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	p24LifecycleCommands(t, a.binary, a.workspace)
}
