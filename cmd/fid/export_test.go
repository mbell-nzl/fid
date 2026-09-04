package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestExportCreateValidateAndImportRemainsAbsent(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	archivePath := filepath.Join(t.TempDir(), "fid-export.zip")

	created := executeCommandForTest(workspace, "export", "create", archivePath, "--json")
	if created.err != nil || !strings.Contains(created.stdout, `"operation":"export.create"`) {
		t.Fatalf("export create = %#v, want JSON success", created)
	}
	validated := executeCommandForTest(workspace, "export", "validate", archivePath, "--json")
	if validated.err != nil || !strings.Contains(validated.stdout, `"operation":"export.validate"`) {
		t.Fatalf("export validate = %#v, want JSON success", validated)
	}
	imported := executeCommandForTest(workspace, "export", "import", archivePath)
	if imported.err == nil || !strings.Contains(imported.err.Error(), "unknown command") {
		t.Fatalf("export import = %#v, want Cobra unknown-command rejection", imported)
	}
	capabilities := executeCommandForTest(workspace, "capabilities", "--json")
	if capabilities.err != nil || !strings.Contains(capabilities.stdout, `"capability_id":"export.portable.v1","limits":{"create":true,"import":false,"validate":true},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"`) {
		t.Fatalf("export capability = %#v, want exact unavailable row", capabilities)
	}
}
