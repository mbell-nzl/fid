package export

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestArchiveSmallStreamingBoundaries(t *testing.T) {
	t.Parallel()
	limits := archiveLimits{nonArtifactBytes: 16, artifactBytes: 32, expandedBytes: 48, entries: 10}
	for _, artifact := range []bool{false, true} {
		limit := int64(16)
		if artifact {
			limit = 32
		}
		for _, delta := range []int64{-1, 0, 1} {
			t.Run(fmt.Sprintf("artifact=%t/delta=%d", artifact, delta), func(t *testing.T) {
				p := writeZeroZIP(t, []zipLimitMember{{artifact: artifact, size: limit + delta}})
				r, err := zip.OpenReader(p)
				if err != nil {
					t.Fatal(err)
				}
				defer func() { _ = r.Close() }()
				scan, checks := scanArchiveWithLimits(&r.Reader, limits, func(_ archiveMember, r io.Reader) error { _, err := io.Copy(io.Discard, r); return err })
				if scan.ExpandedBytes != limit+delta {
					t.Fatalf("expanded=%d", scan.ExpandedBytes)
				}
				if delta <= 0 && len(checks) != 0 {
					t.Fatal(checks)
				}
				if delta > 0 && (len(checks) != 1 || checks[0].Code != "resource-limit") {
					t.Fatal(checks)
				}
			})
		}
	}
	for _, delta := range []int64{-1, 0, 1} {
		t.Run(fmt.Sprintf("cumulative/%d", delta), func(t *testing.T) {
			members := []zipLimitMember{{size: 16}, {size: 16}, {size: 15}}
			if delta >= 0 {
				members[2].size = 16
			}
			if delta > 0 {
				members = append(members, zipLimitMember{size: 1})
			}
			p := writeZeroZIP(t, members)
			r, err := zip.OpenReader(p)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = r.Close() }()
			scan, checks := scanArchiveWithLimits(&r.Reader, limits, func(_ archiveMember, r io.Reader) error { _, err := io.Copy(io.Discard, r); return err })
			if scan.ExpandedBytes != 48+delta {
				t.Fatalf("expanded=%d", scan.ExpandedBytes)
			}
			if delta <= 0 && len(checks) != 0 {
				t.Fatal(checks)
			}
			if delta > 0 && (len(checks) != 1 || checks[0].Code != "resource-limit") {
				t.Fatal(checks)
			}
		})
	}
}

func TestArchiveSmallEntryCountBoundaries(t *testing.T) {
	t.Parallel()
	for _, count := range []int{2, 3, 4} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			p := writeZeroZIP(t, make([]zipLimitMember, count))
			f, err := os.Open(p)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = f.Close() }()
			s, err := f.Stat()
			if err != nil {
				t.Fatal(err)
			}
			err = inspectArchiveEntryCountWithLimit(f, s.Size(), 3)
			if count <= 3 && err != nil {
				t.Fatal(err)
			}
			if count > 3 && !errors.Is(err, errArchiveResourceLimit) {
				t.Fatalf("error=%v", err)
			}
			r, err := zip.NewReader(f, s.Size())
			if err != nil {
				t.Fatal(err)
			}
			_, checks := inspectArchiveDirectoryWithLimit(r.File, 3)
			if count <= 3 && len(checks) != 0 {
				t.Fatal(checks)
			}
			if count > 3 && (len(checks) != 1 || checks[0].Code != "resource-limit") {
				t.Fatal(checks)
			}
		})
	}
}

func TestArchiveProductionBudgetUses64BitBoundaryAndStopsImmediately(t *testing.T) {
	t.Parallel()
	scan, checks := scanZIP(t, writeZeroZIP(t, []zipLimitMember{{size: 3}}))
	if scan.ExpandedBytes != 3 || len(checks) != 0 {
		t.Fatalf("production scan=%+v checks=%v", scan, checks)
	}
	budget := newArchiveByteBudget()
	if budget.limits != (archiveLimits{64 << 20, 1 << 30, 8 << 30, 1000000}) {
		t.Fatalf("production limits=%+v", budget.limits)
	}
	budget.used = (8 << 30) - 1
	input := bytes.NewReader([]byte("abcd"))
	r := archiveLimitReader{Reader: input, limit: 1 << 30, budget: budget}
	n, err := r.Read(make([]byte, 10))
	if n != 2 || !errors.Is(err, errArchiveResourceLimit) || budget.used != (8<<30)+1 || input.Len() != 2 {
		t.Fatalf("n=%d err=%v used=%d unread=%d", n, err, budget.used, input.Len())
	}
	n, err = r.Read(make([]byte, 10))
	if n != 0 || !errors.Is(err, errArchiveResourceLimit) || input.Len() != 2 {
		t.Fatalf("second n=%d err=%v unread=%d", n, err, input.Len())
	}
}
