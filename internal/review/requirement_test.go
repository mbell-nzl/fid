package review

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRequirementRecordReferenceSupersedesPriorActiveRow(t *testing.T) {
	t.Parallel()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	if _, err := conn.Exec(`CREATE TABLE requirement_refs (id TEXT PRIMARY KEY, requirement_id TEXT, status TEXT, superseded_by_id TEXT); INSERT INTO requirement_refs VALUES ('RREF-old', 'PRD-1', 'active', NULL)`); err != nil {
		t.Fatal(err)
	}
	effects, err := requirementReferenceEffects(conn, "PRD-1", "RREF-next", "active")
	if err != nil {
		t.Fatal(err)
	}
	if len(effects) != 2 || effects[0].Table != "requirement_refs" || effects[0].Op != "update" || effects[0].PK["id"] != "RREF-old" || effects[0].Values["status"] != "superseded" || effects[0].Values["superseded_by_id"] != "RREF-next" || effects[1].Op != "insert" {
		t.Fatalf("supersession effects = %#v", effects)
	}
}
