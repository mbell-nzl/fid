//go:build integration

package integration

import (
	"database/sql"
	"reflect"
	"testing"
)

func p24AssertRebuildMapParity(t *testing.T, workspace, binary string) {
	t.Helper()
	live := p24ProjectedFieldMaps(t, workspace)
	if len(live) == 0 {
		t.Fatal("rebuild parity workspace has no keyed mutations")
	}
	stdout, stderr, err := closureRunCommand(workspace, binary, "--project", workspace, "--json", "rebuild")
	if err != nil || stderr != "" {
		t.Fatalf("rebuild parity = %v stdout=%q stderr=%q", err, stdout, stderr)
	}
	if rebuilt := p24ProjectedFieldMaps(t, workspace); !reflect.DeepEqual(rebuilt, live) {
		t.Fatalf("rebuilt field maps differ from live: live=%v rebuilt=%v", live, rebuilt)
	}
}

func p24ProjectedFieldMaps(t *testing.T, workspace string) map[string]string {
	t.Helper()
	conn := p22FaultDB(t, workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	rows, err := conn.Query(`SELECT operation_scope,idempotency_key,request_field_digests_json FROM idempotency_results ORDER BY operation_scope,idempotency_key`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { p22FaultOK(t, rows.Close()) }()
	result := map[string]string{}
	for rows.Next() {
		var scope, key, encoded string
		if err := rows.Scan(&scope, &key, &encoded); err != nil {
			t.Fatal(err)
		}
		if encoded == "" || encoded == "{}" {
			t.Fatalf("%s/%s has no field digests", scope, key)
		}
		result[scope+"\x00"+key] = encoded
	}
	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
	return result
}
