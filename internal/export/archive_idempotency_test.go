package export

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Recomputed hashes ensure both direct validation and CLI import reach the
// same request/result guard that rebuild checks in the original journal.
func assertIdempotencyArchiveAdmission(t *testing.T, binary string, original map[string][]byte, journalFiles []string, valid bool) {
	t.Helper()
	entries := make(map[string][]byte, len(original))
	for path, data := range original {
		entries[path] = bytes.Clone(data)
	}
	var data []byte
	for _, path := range journalFiles {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		data = append(data, body...)
	}
	lines := bytes.Split(bytes.TrimSuffix(data, []byte{'\n'}), []byte{'\n'})
	var terminal journal.Event
	if err := terminal.UnmarshalJSON(lines[len(lines)-1]); err != nil {
		t.Fatal(err)
	}
	hash, err := journal.HashEvent(terminal)
	if err != nil {
		t.Fatal(err)
	}
	path := archivePathWithPrefix(t, entries, "journal/")
	entries[path] = data
	manifest := decodeManifestMap(t, entries["manifest.json"])
	manifest["terminal_event_hash"] = hash
	updateManifestEntry(t, manifest, path, data)
	entries["manifest.json"] = encodeManifestMap(t, manifest)
	archive := filepath.Join(t.TempDir(), "guard.zip")
	writeTestArchive(t, archive, entries)
	_, err = Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
	if valid && err != nil {
		t.Fatalf("valid check archive=%v", err)
	}
	if !valid {
		var invalid *InvalidArchiveError
		if !errors.As(err, &invalid) {
			t.Fatalf("malformed check archive=%v", err)
		}
	}
	target := gittest.NewWorkspace(t)
	before := workspaceTreeDigest(t, target)
	cmd := exec.Command(binary, "--json", "export", "import", archive, "--repository", target, "--idempotency-key", "guard-import")
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err = cmd.Run()
	if valid {
		if err != nil || stderr.Len() != 0 {
			t.Fatalf("valid check import=%v/%s", err, stderr.String())
		}
		return
	}
	var exit *exec.ExitError
	var failure struct {
		Operation string `json:"operation"`
		Error     struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if !errors.As(err, &exit) || exit.ExitCode() != 6 || stdout.Len() != 0 || json.Unmarshal(stderr.Bytes(), &failure) != nil || failure.Operation != "export.import" || failure.Error.Code != "archive-invalid" || workspaceTreeDigest(t, target) != before {
		t.Fatalf("malformed check import=%v stdout=%s stderr=%s", err, stdout.String(), stderr.String())
	}
}
