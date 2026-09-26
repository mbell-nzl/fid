package evidence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRejectedArtifactsResolvePathsUnderExecutionCWD(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	if err := os.Mkdir(filepath.Join(workspace, "cwd"), 0o755); err != nil {
		t.Fatalf("create execution cwd: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "sibling.txt"), []byte("outside execution cwd"), 0o600); err != nil {
		t.Fatalf("write workspace-only Artifact: %v", err)
	}
	rejections := rejectedArtifacts(workspace, "cwd", []string{"sibling.txt"}, nil)
	if len(rejections) != 1 || rejections[0].Path != "sibling.txt" || rejections[0].Reason != "missing" {
		t.Fatalf("Artifact rejections = %#v, want cwd-relative missing path", rejections)
	}
}
