package export

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/flate"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

type zipLimitMember struct {
	artifact bool
	size     int64
}

func writeZeroZIP(t *testing.T, members []zipLimitMember) string {
	t.Helper()
	if err := os.MkdirAll("../../.cache", 0o755); err != nil {
		t.Fatal(err)
	}
	d, err := os.MkdirTemp("../../.cache", "archive-limit-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(d) })
	p := filepath.Join(d, "a.zip")
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	buffer := bufio.NewWriter(f)
	z := zip.NewWriter(buffer)
	zero := make([]byte, 64<<10)
	block := zeroDeflateBlock(t, zero, false)
	for i, m := range members {
		name := fmt.Sprintf("journal/%08d.jsonl", i+1)
		if m.artifact {
			name = fmt.Sprintf("artifacts/%064d", i)
		}
		full, remainder := m.size/int64(len(zero)), m.size%int64(len(zero))
		final := zeroDeflateBlock(t, zero[:remainder], true)
		var checksum uint32
		for j := int64(0); j < full; j++ {
			checksum = crc32.Update(checksum, crc32.IEEETable, zero)
		}
		checksum = crc32.Update(checksum, crc32.IEEETable, zero[:remainder])
		header := &zip.FileHeader{Name: name, Method: zip.Deflate, CRC32: checksum, UncompressedSize64: uint64(m.size), CompressedSize64: uint64(full*int64(len(block)) + int64(len(final)))}
		w, err := z.CreateRaw(header)
		if err != nil {
			t.Fatal(err)
		}
		for j := int64(0); j < full; j++ {
			if _, err := w.Write(block); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := w.Write(final); err != nil {
			t.Fatal(err)
		}
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	if err := buffer.Flush(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return p
}

// Flush ends on a byte boundary without ending the DEFLATE stream. Each block
// is generated from a fresh dictionary, so repeating it needs no external data.
// ZIP readers still decompress every byte and verify the complete member CRC.
func zeroDeflateBlock(t *testing.T, zero []byte, final bool) []byte {
	t.Helper()
	var data bytes.Buffer
	w, err := flate.NewWriter(&data, flate.BestSpeed)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write(zero); err != nil {
		t.Fatal(err)
	}
	if final {
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
		return data.Bytes()
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
	block := bytes.Clone(data.Bytes())
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return block
}

func scanZIP(t *testing.T, path string) (archiveScan, []Check) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	s, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	r, err := zip.NewReader(f, s.Size())
	if err != nil {
		t.Fatal(err)
	}
	return scanArchive(r, func(_ archiveMember, in io.Reader) error { _, err := io.Copy(io.Discard, in); return err })
}

func TestArchiveStopsAtFirstExpandedOverflow(t *testing.T) {
	path := writeZeroZIP(t, []zipLimitMember{{size: maxArchiveNonArtifactBytes + 2}, {size: 1}})
	reader, err := zip.OpenReader(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = reader.Close() }()
	consumed := 0
	scan, checks := scanArchive(&reader.Reader, func(_ archiveMember, input io.Reader) error {
		consumed++
		_, err := io.Copy(io.Discard, input)
		return err
	})
	if consumed != 1 || scan.ExpandedBytes != maxArchiveNonArtifactBytes+1 || len(checks) != 1 || checks[0].Code != "resource-limit" {
		t.Fatalf("consumed=%d scan=%+v checks=%v", consumed, scan, checks)
	}
}

func TestReadEventLinesUsesPerLineLimitAcrossFragments(t *testing.T) {
	if err := readEventLines(&fragmentReader{data: []byte(strings.Repeat("x", (1<<20)-1) + "\n"), width: 257}, func(value []byte) error {
		if len(value) != (1<<20)-1 {
			t.Fatal(len(value))
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	line := strings.Repeat("x", 1<<20)
	var got int
	if err := readEventLines(&fragmentReader{data: []byte(line + "\n" + line + "\n"), width: 257}, func(value []byte) error { got += len(value); return nil }); err != nil || got != 2*(1<<20) {
		t.Fatalf("legal fragmented lines: got=%d err=%v", got, err)
	}
	if err := readEventLines(&fragmentReader{data: []byte(strings.Repeat("x", (1<<20)+1)), width: 257}, func([]byte) error { return nil }); err == nil {
		t.Fatal("over-limit fragmented line accepted")
	}
}

type fragmentReader struct {
	data  []byte
	width int
}

func (r *fragmentReader) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	n := r.width
	if n > len(p) {
		n = len(p)
	}
	if n > len(r.data) {
		n = len(r.data)
	}
	copy(p, r.data[:n])
	r.data = r.data[n:]
	return n, nil
}

func TestSafeArchivePathRejectsPlatformEscapes(t *testing.T) {
	for _, value := range []string{"../x", "/x", "C:/x", "C:x", `x\\y`, "x/../y", "x//y"} {
		if safeArchivePath(value) {
			t.Fatalf("accepted unsafe %q", value)
		}
	}
	if !safeArchivePath("journal/00000001.jsonl") {
		t.Fatal("rejected safe path")
	}
}
