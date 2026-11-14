//go:build integration

package integration

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestP24RetainedReadFailurePairs(t *testing.T) {
	t.Parallel()
	rows := []struct {
		name, operation, message, reason string
		args                             []string
		project                          bool
		exit                             int
	}{
		{"since", "activity.list", "activity list: --since must be YYYY-MM-DD or RFC3339", "invalid-argument; argument=request, allowed=accepted issue read input", []string{"activity", "list", "--since", "invalid"}, true, 2},
		{"on", "activity.list", "activity list: --on must be <type>/<id>", "invalid-argument; argument=request, allowed=accepted issue read input", []string{"activity", "list", "--on", "invalid"}, true, 2},
		{"activity-no-project", "activity.list", "fid project was not found", "project-not-found; path=", []string{"activity", "list"}, false, 3},
		{"ac-no-project", "issue.ac.list", "fid project was not found", "project-not-found; path=", []string{"issue", "ac", "list", "FID-missing"}, false, 3},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			binary, workspace := integrationFIDBinary(t), gittest.NewWorkspace(t)
			if row.project {
				p21Run(t, binary, workspace, "init", "--prefix", "P24", "--json")
			}
			before := p24FailureFilesystem(t, workspace)
			args := append([]string{"--project", workspace}, row.args...)
			out, stderr, err := closureRunCommand(workspace, binary, append([]string{"--json"}, args...)...)
			if err == nil || processExitCode(t, err) != row.exit || out != "" {
				t.Fatalf("JSON refusal: %v stdout=%q stderr=%q", err, out, stderr)
			}
			payload := p23Envelope(t, stderr)
			f := payload["error"].(map[string]any)
			code, category := "project-not-found", "not-found"
			var subject any = map[string]any{"type": "project", "id": ""}
			details := map[string]any{"path": ""}
			if row.project {
				code, category, subject = "invalid-argument", "invalid-input", nil
				details = map[string]any{"argument": "request", "value": "", "allowed": []any{"accepted issue read input"}}
			}
			if payload["operation"] != row.operation || f["message"] != row.message || f["code"] != code || f["category"] != category || f["retryable"] != false || !reflect.DeepEqual(f["subject"], subject) || !reflect.DeepEqual(f["details"], details) {
				t.Fatalf("DEC-285 mapping: %#v", payload)
			}
			if !reflect.DeepEqual(before, p24FailureFilesystem(t, workspace)) {
				t.Fatal("JSON refusal changed filesystem")
			}
			out, stderr, err = closureRunCommand(workspace, binary, args...)
			want := "error: " + f["message"].(string) + "\nreason: " + row.reason + "\n"
			if err == nil || processExitCode(t, err) != row.exit || out != "" || stderr != want {
				t.Fatalf("Human refusal: %v stdout=%q stderr=%q want=%q", err, out, stderr, want)
			}
			if !reflect.DeepEqual(before, p24FailureFilesystem(t, workspace)) {
				t.Fatal("Human refusal changed filesystem")
			}
		})
	}
}

// Includes .fid, the projection, journal, Git metadata, file modes and names.
func p24FailureFilesystem(t *testing.T, root string) map[string]string {
	t.Helper()
	result := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		value := info.Mode().String()
		if info.Mode().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			value += fmt.Sprintf(":%x", sha256.Sum256(data))
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		result[relative] = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

// DEC-286: refusal precedes all workspace and archive access.
func TestP24UnavailableImportRootFlagFailurePair(t *testing.T) {
	t.Parallel()
	const jsonFailure = `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["--idempotency-key","--repository"],"argument":"--project","value":"--project"},"message":"--project is not available for export import","retryable":false,"subject":null},"ok":false,"operation":"export.import"}` + "\n"
	const humanFailure = "error: --project is not available for export import\nreason: invalid-argument; argument=--project, allowed=--idempotency-key|--repository\n"
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprintf("missing-project=%t", missing), func(t *testing.T) {
			t.Parallel()
			binary, root := integrationFIDBinary(t), t.TempDir()
			project := root
			if missing {
				project = filepath.Join(root, "nonexistent")
			}
			before := p24FailureFilesystem(t, root)
			args := []string{"--project", project, "export", "import", "missing.zip", "--repository", root, "--idempotency-key", "root-flag"}
			for _, jsonMode := range []bool{false, true} {
				request, want := append([]string(nil), args...), humanFailure
				if jsonMode {
					request = append([]string{"--json"}, request...)
					want = jsonFailure
				}
				out, stderr, err := closureRunCommand(root, binary, request...)
				if err == nil || processExitCode(t, err) != 2 || out != "" || stderr != want {
					t.Fatalf("json=%t refusal: %v stdout=%q stderr=%q want=%q", jsonMode, err, out, stderr, want)
				}
				if !reflect.DeepEqual(before, p24FailureFilesystem(t, root)) {
					t.Fatal("refusal changed filesystem")
				}
			}
		})
	}
}
