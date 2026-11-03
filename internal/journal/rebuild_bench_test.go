package journal

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

func BenchmarkRebuildWorkspace1kEvents(b *testing.B) {
	benchmarkRebuildWorkspace(b, 1000)
}

func BenchmarkRebuildWorkspace10kEvents(b *testing.B) {
	benchmarkRebuildWorkspace(b, 10000)
}

func BenchmarkRebuildWorkspace100kEvents(b *testing.B) {
	benchmarkRebuildWorkspace(b, 100000)
}

func benchmarkRebuildWorkspace(b *testing.B, targetEvents int) {
	workspace := b.TempDir()
	journalDir := filepath.Join(config.LedgerDir(workspace), "journal")
	if err := os.MkdirAll(journalDir, 0o755); err != nil {
		b.Fatalf("osMkdirAll(%q) error = %v", journalDir, err)
	}

	seed := seedScaleReplayWorkspace(b, workspace, targetEvents)
	migrationsDir := testProjectMigrationsDirTB(b)

	b.ReportAllocs()
	b.ResetTimer()

	var total time.Duration
	for iteration := 0; iteration < b.N; iteration++ {
		start := time.Now()
		stats, err := RebuildWorkspace(workspace, migrationsDir)
		total += time.Since(start)
		if err != nil {
			b.Fatalf("RebuildWorkspace() error = %v", err)
		}
		if stats.EventsReplayed != seed.EventCount {
			b.Fatalf("rebuilt events = %d, want %d", stats.EventsReplayed, seed.EventCount)
		}
	}

	b.StopTimer()
	if b.N > 0 {
		average := total / time.Duration(b.N)
		b.ReportMetric(float64(average.Microseconds())/1000, "ms/rebuild")
		if limit, ok := rebuildTarget(targetEvents); ok && average >= limit {
			b.Fatalf("%dk-event rebuild = %s, want <%s", targetEvents/1000, average, limit)
		}
	}
}

func rebuildTarget(eventCount int) (time.Duration, bool) {
	switch eventCount {
	case 10000:
		return 5 * time.Second, true
	case 100000:
		return 30 * time.Second, true
	default:
		return 0, false
	}
}

func testProjectMigrationsDirTB(tb testing.TB) string {
	tb.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		tb.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}
