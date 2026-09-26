package export

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestArchiveActualZIPPathTypeAndDuplicateMatrix(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, path string
		mode       os.FileMode
		duplicate  bool
		code       string
	}{
		{"safe relative regular", "journal/00000001.jsonl", 0o644, false, ""},
		{"current directory", ".", 0o644, false, "path-unsafe"},
		{"NUL path", "journal/entry\x00", 0o644, false, "path-unsafe"},
		{"parent traversal", "../escape", 0o644, false, "path-unsafe"},
		{"escaping component", "journal/../../escape", 0o644, false, "path-unsafe"},
		{"noncanonical component", "journal/../entry", 0o644, false, "path-unsafe"},
		{"absolute", "/escape", 0o644, false, "path-unsafe"},
		{"drive rooted", "C:/escape", 0o644, false, "path-unsafe"},
		{"drive relative", "C:escape", 0o644, false, "path-unsafe"},
		{"backslash traversal", `..\escape`, 0o644, false, "path-unsafe"},
		{"symlink", "journal/entry", os.ModeSymlink | 0o777, false, "path-unsafe"},
		{"device", "journal/entry", os.ModeDevice | 0o600, false, "path-unsafe"},
		{"character device", "journal/entry", os.ModeDevice | os.ModeCharDevice | 0o600, false, "path-unsafe"},
		{"directory mode", "journal/entry", os.ModeDir | 0o755, false, "path-unsafe"},
		{"directory name", "journal/", 0o644, false, "path-unsafe"},
		{"duplicate", "journal/entry", 0o644, true, "duplicate-entry"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			archive := filepath.Join(t.TempDir(), "paths.zip")
			output, err := os.Create(archive)
			if err != nil {
				t.Fatal(err)
			}
			writer := zip.NewWriter(output)
			count := 1
			if tc.duplicate {
				count = 2
			}
			for i := 0; i < count; i++ {
				header := &zip.FileHeader{Name: tc.path, Method: zip.Store}
				header.SetMode(tc.mode)
				if _, err := writer.CreateHeader(header); err != nil {
					t.Fatal(err)
				}
			}
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			if err := output.Close(); err != nil {
				t.Fatal(err)
			}
			reader, err := zip.OpenReader(archive)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = reader.Close() }()
			consumed := 0
			_, checks := scanArchive(&reader.Reader, func(_ archiveMember, input io.Reader) error {
				consumed++
				_, err := io.Copy(io.Discard, input)
				return err
			})
			if tc.code == "" {
				if len(checks) != 0 || consumed != 1 {
					t.Fatalf("valid ZIP checks=%v consumed=%d", checks, consumed)
				}
				return
			}
			if len(checks) != 1 || checks[0].Code != tc.code || consumed != 0 {
				t.Fatalf("checks=%v consumed=%d want %s before expansion", checks, consumed, tc.code)
			}
		})
	}
}
