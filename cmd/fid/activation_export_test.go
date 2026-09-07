package main

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestIncompleteActivationAllowsExportWithoutChangingActivation(t *testing.T) {
	t.Parallel()
	workspace, _, _ := buildFinalizedEnvelopeWorkspace(t, validGateDefinitions())
	prior := currentFinalizedActivationForTest(t, workspace)
	upgradeID := prepareUpgradeActivationForTest(t, workspace, prior)
	before, err := workflow.ActivationStatus(workspace, issueMigrationsDir(), upgradeID)
	if err != nil || before.Status != "prepared" {
		t.Fatalf("incomplete activation = (%#v, %v), want prepared upgrade", before, err)
	}
	eventsBefore := journalEventCount(t, workspace)
	archive := filepath.Join(t.TempDir(), "incomplete-activation.zip")
	for _, action := range []string{"create", "validate"} {
		result := executeCommandForTest(workspace, "--project", workspace, "--json", "export", action, archive)
		if result.err != nil || !strings.Contains(result.stdout, `"operation":"export.`+action+`"`) {
			t.Fatalf("export %s during incomplete activation = %#v, want success", action, result)
		}
		after, err := workflow.ActivationStatus(workspace, issueMigrationsDir(), upgradeID)
		if err != nil || !reflect.DeepEqual(after, before) {
			t.Fatalf("activation after export %s = (%#v, %v), want unchanged %#v", action, after, err, before)
		}
		if count := journalEventCount(t, workspace); count != eventsBefore {
			t.Fatalf("Event count after export %s = %d, want unchanged %d", action, count, eventsBefore)
		}
	}
}
