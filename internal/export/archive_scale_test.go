//go:build archive_scale

package export

import (
	"archive/zip"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestArchiveExpandedByteBoundaries(t *testing.T) {
	for _, tc := range []struct {
		name     string
		artifact bool
		size     int64
		want     bool
	}{{"nonartifact minus", false, maxArchiveNonArtifactBytes - 1, true}, {"nonartifact equal", false, maxArchiveNonArtifactBytes, true}, {"nonartifact plus", false, maxArchiveNonArtifactBytes + 1, false}, {"artifact minus", true, maxArchiveArtifactBytes - 1, true}, {"artifact equal", true, maxArchiveArtifactBytes, true}, {"artifact plus", true, maxArchiveArtifactBytes + 1, false}} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeZeroZIP(t, []zipLimitMember{{artifact: tc.artifact, size: tc.size}})
			scan, checks := scanZIP(t, path)
			if tc.want && len(checks) != 0 {
				t.Fatalf("checks=%v", checks)
			}
			if tc.want && scan.ExpandedBytes != tc.size {
				t.Fatalf("expanded=%d want %d", scan.ExpandedBytes, tc.size)
			}
			if !tc.want && (len(checks) == 0 || checks[0].Code != "resource-limit") {
				t.Fatalf("checks=%v", checks)
			}
			if !tc.want && scan.ExpandedBytes != tc.size {
				t.Fatalf("expanded=%d want %d", scan.ExpandedBytes, tc.size)
			}
		})
	}
}
func TestArchiveCumulativeBoundary(t *testing.T) {
	for _, size := range []int64{maxArchiveExpandedBytes - 1, maxArchiveExpandedBytes, maxArchiveExpandedBytes + 1} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			members := make([]zipLimitMember, 0, 129)
			left := size
			for left > 0 {
				n := maxArchiveNonArtifactBytes
				if left < n {
					n = left
				}
				members = append(members, zipLimitMember{size: n})
				left -= n
			}
			scan, checks := scanZIP(t, writeZeroZIP(t, members))
			if size <= maxArchiveExpandedBytes && len(checks) != 0 {
				t.Fatal(checks)
			}
			if size > maxArchiveExpandedBytes && (len(checks) == 0 || checks[0].Code != "resource-limit") {
				t.Fatal(checks)
			}
			if scan.ExpandedBytes != size {
				t.Fatal(scan.ExpandedBytes)
			}
		})
	}
}

func TestArchiveEntryCountBoundaries(t *testing.T) {
	if err := os.MkdirAll("../../.cache", 0o755); err != nil {
		t.Fatal(err)
	}
	for _, count := range []int{999999, 1000000, 1000001} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			d, err := os.MkdirTemp("../../.cache", "archive-count-")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.RemoveAll(d) })
			p := filepath.Join(d, "count.zip")
			f, err := os.Create(p)
			if err != nil {
				t.Fatal(err)
			}
			z := zip.NewWriter(f)
			for i := 0; i < count; i++ {
				h := &zip.FileHeader{Name: fmt.Sprintf("journal/%08d.jsonl", i), Method: zip.Store}
				h.UncompressedSize64 = 0
				h.CompressedSize64 = 0
				if _, err := z.CreateRaw(h); err != nil {
					t.Fatal(err)
				}
			}
			if err := z.Close(); err != nil {
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
			runtime.GC()
			f, err = os.Open(p)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = f.Close() }()
			s, err := f.Stat()
			if err != nil {
				t.Fatal(err)
			}
			countErr := inspectArchiveEntryCount(f, s.Size())
			if count <= maxArchiveEntries && countErr != nil {
				t.Fatalf("pre-allocation count: %v", countErr)
			}
			if count > maxArchiveEntries && !errors.Is(countErr, errArchiveResourceLimit) {
				t.Fatalf("pre-allocation count: %v", countErr)
			}
			r, err := zip.NewReader(f, s.Size())
			if err != nil {
				t.Fatal(err)
			}
			_, checks := inspectArchiveDirectory(r.File)
			if count <= maxArchiveEntries && len(checks) != 0 {
				t.Fatal(checks)
			}
			if count > maxArchiveEntries && (len(checks) == 0 || checks[0].Code != "resource-limit") {
				t.Fatal(checks)
			}
		})
	}
}
