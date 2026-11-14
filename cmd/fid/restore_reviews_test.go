package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func recordV8Reviews(t *testing.T, root, task string) {
	t.Helper()
	for i, stream := range []string{"primary", "separate"} {
		in := reviewInput()
		in["stream"] = stream
		in["verdict"] = "CHANGES"
		in["reviewer"].(map[string]any)["session_id"] = "v8-" + stream
		if i == 0 {
			in["findings"] = []any{reviewFinding("major"), reviewFinding("minor"), reviewFinding("nit")}
		}
		raw, err := json.Marshal(in)
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(t.TempDir(), "review.json")
		if err := os.WriteFile(path, raw, 0o600); err != nil {
			t.Fatal(err)
		}
		lifecycleResult(t, executeCommandForTest(root, "--json", "review", "record", task, "--input", path, "--actor", "human:courier", "--idempotency-key", "v8-"+stream), "review.record")
	}
}

func assertV8Reviews(t *testing.T, workspace string) {
	t.Helper()
	conn, err := reviewwork.OpenProjection(workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var rounds, findings, models, manual int
	for _, q := range []struct {
		sql string
		to  *int
	}{
		{`SELECT count(*) FROM review_rounds`, &rounds}, {`SELECT count(*) FROM round_findings`, &findings},
		{`SELECT count(*) FROM holds WHERE model<>''`, &models}, {`SELECT count(*) FROM holds WHERE model=''`, &manual},
	} {
		if err := conn.QueryRow(q.sql).Scan(q.to); err != nil {
			t.Fatal(err)
		}
	}
	if rounds != 2 || findings != 3 || models == 0 || manual == 0 {
		t.Fatalf("V8 reviews=%d findings=%d modeled holds=%d manual holds=%d", rounds, findings, models, manual)
	}
	for _, severity := range []string{"major", "minor", "nit"} {
		var n int
		if err := conn.QueryRow(`SELECT count(*) FROM round_findings WHERE json_extract(record_json,'$.severity')=?`, severity).Scan(&n); err != nil || n != 1 {
			t.Fatalf("V8 %s=%d/%v", severity, n, err)
		}
	}
	for _, stream := range []string{"primary", "separate"} {
		var n int
		if err := conn.QueryRow(`SELECT count(*) FROM review_rounds WHERE json_extract(record_json,'$.stream')=?`, stream).Scan(&n); err != nil || n != 1 {
			t.Fatalf("V8 %s=%d/%v", stream, n, err)
		}
	}
}
