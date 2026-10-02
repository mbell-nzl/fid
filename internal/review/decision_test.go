package review

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDecisionReferenceSourceMustExistAndMatchDigest(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "decision.md"), []byte("approved decision\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte("approved decision\n"))
	path := "decision.md"
	source := Source{SourceType: "project-file", DisplayLabel: "Decision", Path: &path, Digest: "sha256:" + hex.EncodeToString(sum[:])}
	if err := verifyDecisionSource(root, source); err != nil {
		t.Fatalf("matching Source rejected: %v", err)
	}

	source.Digest = "sha256:" + strings.Repeat("0", 64)
	var mismatch *SourceDigestMismatchError
	if err := verifyDecisionSource(root, source); !errors.As(err, &mismatch) {
		t.Fatalf("digest mismatch error = %T %v", err, err)
	}
	missing := "missing.md"
	source.Path = &missing
	if err := verifyDecisionSource(root, source); err == nil {
		t.Fatal("missing Source file accepted")
	}
}
