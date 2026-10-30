package project

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveSourceInputProjectFileCases(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	outside := t.TempDir()
	content := []byte("accepted source\n")
	if err := os.WriteFile(filepath.Join(root, "source.md"), content, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("source.md", filepath.Join(root, "inside.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "outside.md"), filepath.Join(root, "outside.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "directory"), 0o700); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(content)
	matching := "sha256:" + hex.EncodeToString(sum[:])
	wrong := "sha256:" + strings.Repeat("0", 64)
	for _, tc := range []struct {
		name, path, wantDigest, wantField string
		digest                            *string
		mismatch, wantObserved            bool
	}{
		{"matching", "source.md", matching, "", &matching, false, false},
		{"omitted", "source.md", matching, "", nil, false, false},
		{"inside symlink", "inside.md", matching, "", nil, false, false},
		{"mismatch", "source.md", "", "", &wrong, true, true},
		{"missing omitted", "missing.md", "", "", nil, true, false},
		{"outside symlink", "outside.md", "", "", nil, true, false},
		{"directory", "directory", "", "", nil, true, false},
		{"absolute", filepath.Join(root, "source.md"), "", "path", nil, false, false},
		{"parent", "../source.md", "", "path", nil, false, false},
		{"dot segment", "./source.md", "", "path", nil, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			source := SourceInput{SourceType: "project-file", DisplayLabel: "Source", Path: &tc.path, Digest: tc.digest}
			got, err := ResolveSourceInput(root, source)
			if tc.wantField != "" {
				var invalid *InvalidSourceInputError
				if !errors.As(err, &invalid) || invalid.Field != tc.wantField {
					t.Fatalf("error = %v, want invalid %s", err, tc.wantField)
				}
				return
			}
			if tc.mismatch {
				var mismatch *SourceDigestMismatchError
				if !errors.As(err, &mismatch) || (mismatch.ObservedDigest != nil) != tc.wantObserved || mismatch.ExpectedDigest != tc.digest {
					t.Fatalf("error = %v, want mismatch observed=%t", err, tc.wantObserved)
				}
				return
			}
			if err != nil || got != tc.wantDigest {
				t.Fatalf("digest = %q, error = %v, want %q", got, err, tc.wantDigest)
			}
		})
	}
}

func TestResolveSourceInputExternalReference(t *testing.T) {
	t.Parallel()
	digest := "sha256:" + strings.Repeat("a", 64)
	id := "external-1"
	source := SourceInput{SourceType: "external-reference", DisplayLabel: "External", ExternalID: &id}
	_, err := ResolveSourceInput(t.TempDir(), source)
	var invalid *InvalidSourceInputError
	if !errors.As(err, &invalid) || invalid.Field != "digest" {
		t.Fatalf("omitted digest = %v", err)
	}
	source.Digest = &digest
	got, err := ResolveSourceInput(t.TempDir(), source)
	if err != nil || got != digest {
		t.Fatalf("supplied external digest = %q, %v", got, err)
	}
}
