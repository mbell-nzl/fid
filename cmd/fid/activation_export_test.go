package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestMalformedConfigRejectsActivationMutationButNotHelpOrExportValidation(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	archive := filepath.Join(t.TempDir(), "config-independent.zip")
	created := executeCommandForTest(workspace, "--project", workspace, "--json", "export", "create", archive)
	if created.err != nil {
		t.Fatalf("create archive before malformed config: %#v", created)
	}

	configPath := filepath.Join(workspace, ".fid", "config.toml")
	contents, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, append(contents, []byte("\nunknown_key = true\n")...), 0o600); err != nil {
		t.Fatal(err)
	}

	before := journalEventCount(t, workspace)
	inputPath := filepath.Join(workspace, "malformed-config-activation.json")
	writeActivationInput(t, inputPath, activationPrepareInput("config-boundary"))
	consuming := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "prepare", "--input", inputPath, "--idempotency-key", "malformed-config-activation")
	if consuming.err == nil || consuming.stdout != "" || !strings.Contains(consuming.stderr, `"code":"invalid-argument"`) {
		t.Fatalf("activation prepare with unknown config key = %#v, want typed consuming-path rejection", consuming)
	}
	if after := journalEventCount(t, workspace); after != before {
		t.Fatalf("malformed-config activation prepare Event count = %d, want unchanged %d", after, before)
	}

	help := executeCommandForTest(workspace, "--project", workspace, "activation", "--help")
	if help.err != nil || !strings.Contains(help.stdout, "Manage kit activation trust") {
		t.Fatalf("activation help with malformed config = %#v, want available help", help)
	}
	validated := executeCommandForTest(workspace, "--project", workspace, "--json", "export", "validate", archive)
	if validated.err != nil || !strings.Contains(validated.stdout, `"operation":"export.validate"`) {
		t.Fatalf("export validate with malformed config = %#v, want config-independent success", validated)
	}
}

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
