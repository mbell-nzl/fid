package issues

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func BenchmarkList10kIssuesWithRegistry(b *testing.B) {
	workspace := b.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: benchmarkIssueMigrationsDir(b),
	}); err != nil {
		b.Fatalf("InitializeWorkspace() error = %v", err)
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		b.Fatalf("open projection: %v", err)
	}
	for index := 0; index < 10000; index++ {
		if _, err := conn.Exec(`
			INSERT INTO issues (id, kind, title, workflow_id, status, created_at, updated_at)
			VALUES (?, 'task', ?, 'task-default', 'todo', '2026-07-19T00:00:00Z', '2026-07-19T00:00:00Z')
		`, fmt.Sprintf("FID-registry-read-%05d", index), fmt.Sprintf("Registry read %05d", index)); err != nil {
			_ = conn.Close()
			b.Fatalf("seed issue %d: %v", index, err)
		}
	}
	if err := conn.Close(); err != nil {
		b.Fatalf("close seeded projection: %v", err)
	}

	options := ListOptions{WorkspaceRoot: workspace, MigrationsDir: benchmarkIssueMigrationsDir(b)}
	b.ReportAllocs()
	b.ResetTimer()
	var total time.Duration
	for iteration := 0; iteration < b.N; iteration++ {
		startedAt := time.Now()
		items, err := List(options)
		total += time.Since(startedAt)
		if err != nil {
			b.Fatalf("List() error = %v", err)
		}
		if len(items) != 10000 {
			b.Fatalf("List() count = %d, want 10000", len(items))
		}
	}
	b.StopTimer()
	if b.N > 0 {
		average := total / time.Duration(b.N)
		b.ReportMetric(float64(average.Microseconds())/1000, "ms/registry-read")
		if average >= 500*time.Millisecond {
			b.Fatalf("10k registry-backed read = %s, want <500ms", average)
		}
	}
}

func benchmarkIssueMigrationsDir(tb testing.TB) string {
	tb.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		tb.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}
