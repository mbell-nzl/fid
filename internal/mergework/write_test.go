package mergework

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestCIIDsAvoidPersistedAndPendingBatchCollisions(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name       string
		persisted  string
		candidates []string
		want       []string
	}{
		{"pending", "", []string{"100001", "100001", "100002"}, []string{"CIR-100001", "CIR-100002"}},
		{"persisted and pending", "CIR-100001", []string{"100001", "100002", "100002", "100003"}, []string{"CIR-100002", "CIR-100003"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			conn, err := fiddb.Open(filepath.Join(t.TempDir(), "ids.sqlite"))
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = conn.Close() }()
			if _, err := conn.Exec(`CREATE TABLE ci_results (id TEXT PRIMARY KEY)`); err != nil {
				t.Fatal(err)
			}
			if row.persisted != "" {
				if _, err := conn.Exec(`INSERT INTO ci_results (id) VALUES (?)`, row.persisted); err != nil {
					t.Fatal(err)
				}
			}
			next := 0
			generate := func(prefix string, exists journal.IDExistsFunc) (string, error) {
				for next < len(row.candidates) {
					candidate := prefix + "-" + row.candidates[next]
					next++
					taken, err := exists(candidate)
					if err != nil {
						return "", err
					}
					if !taken {
						return candidate, nil
					}
				}
				return "", fmt.Errorf("no unused ID")
			}
			ctx := writeContext{conn: conn, generateID: generate}
			// Allocate both IDs before inserting either row, just as RecordCI does.
			var ids []string
			for range row.want {
				id, err := ctx.id("CIR", "ci_results")
				if err != nil {
					t.Fatal(err)
				}
				ids = append(ids, id)
			}
			if !reflect.DeepEqual(ids, row.want) {
				t.Errorf("allocated %v, want %v", ids, row.want)
			}
			for _, id := range ids {
				if _, err := conn.Exec(`INSERT INTO ci_results (id) VALUES (?)`, id); err != nil {
					t.Fatalf("insert allocated ID %s: %v", id, err)
				}
			}
		})
	}
}
