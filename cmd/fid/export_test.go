package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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

func TestExportValidateRendersInvalidArchiveChecks(t *testing.T) {
	for _, namedEntry := range []bool{false, true} {
		t.Run(fmt.Sprint(namedEntry), func(t *testing.T) {
			archivePath := filepath.Join(t.TempDir(), "invalid.zip")
			data := []byte("not a ZIP")
			if namedEntry {
				var buffer bytes.Buffer
				writer := zip.NewWriter(&buffer)
				entry, err := writer.Create("manifest.json")
				if err != nil {
					t.Fatal(err)
				}
				if _, err := entry.Write([]byte("{}")); err != nil {
					t.Fatal(err)
				}
				if err := writer.Close(); err != nil {
					t.Fatal(err)
				}
				data = buffer.Bytes()
			}
			if err := os.WriteFile(archivePath, data, 0o600); err != nil {
				t.Fatal(err)
			}
			result := executeCommandForTest(t.TempDir(), "export", "validate", archivePath, "--json")
			var payload map[string]any
			if result.err == nil || result.stdout != "" || json.Unmarshal([]byte(result.stderr), &payload) != nil {
				t.Fatalf("invalid archive did not render its JSON failure: %#v", result)
			}
			failure := payload["error"].(map[string]any)
			if failure["code"] != "archive-invalid" {
				t.Fatalf("failure=%v", failure)
			}
			checks := failure["details"].(map[string]any)["checks"].([]any)
			entry := checks[0].(map[string]any)["entry"]
			if namedEntry && entry != "manifest.json" || !namedEntry && entry != nil {
				t.Fatalf("entry=%v", entry)
			}
		})
	}
}
