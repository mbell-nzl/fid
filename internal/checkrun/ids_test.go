package checkrun

import (
	"fmt"
	"path/filepath"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestCheckAllocationNeverReusesARequestOrSecondaryStartID(t *testing.T) {
	t.Parallel()
	conn, err := fiddb.Open(filepath.Join(t.TempDir(), "ids.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	for _, query := range []string{
		`CREATE TABLE check_requests (id TEXT PRIMARY KEY)`,
		`CREATE TABLE check_starts (id TEXT PRIMARY KEY)`,
		`INSERT INTO check_requests VALUES ('CHK-100001')`,
		`INSERT INTO check_starts VALUES ('CHK-100001'),('CHK-100002')`,
	} {
		if _, err := conn.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	// Supply collisions in order, including a secondary start with no request
	// of its own. The observable allocation must skip both occupied IDs.
	generate := func(prefix string, exists journal.IDExistsFunc) (string, error) {
		for _, suffix := range []string{"100001", "100002", "100003"} {
			candidate := prefix + "-" + suffix
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
	id, err := (store{conn: conn}).newCheckID(generate)
	if err != nil || id != "CHK-100003" {
		t.Fatalf("allocated %q: %v; want unused CHK-100003", id, err)
	}
}
