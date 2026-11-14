package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkrun"
)

func TestLeanCheckFreshKeyAfterPriorBootKeepsTheUnfinishedStart(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "linux" {
		t.Skip("Linux boot identity")
	}
	f := leanJourneyFixture(t, leanHelperDefinition(t, "pass", "1m", nil))
	options := leanOptions(f, "prior-boot")
	cut := errors.New("interrupted before spawning")
	options.Before = func(string) error { return cut }
	if _, err := checkrun.Run(context.Background(), options); !errors.Is(err, cut) {
		t.Fatalf("interruption=%v", err)
	}
	record := latestLeanRecord(t, f)
	raw, err := json.Marshal(record.Start)
	if err != nil {
		t.Fatal(err)
	}
	var start map[string]any
	if err := json.Unmarshal(raw, &start); err != nil {
		t.Fatal(err)
	}
	owner := start["owner"].(map[string]any)
	machine, err := os.ReadFile("/etc/machine-id")
	if err == nil && strings.TrimSpace(string(machine)) != "" {
		owner["machine"] = "machine-id:" + strings.TrimSpace(string(machine))
	} else {
		host, err := os.Hostname()
		if err != nil {
			t.Fatal(err)
		}
		owner["machine"] = "hostname:" + host
	}
	owner["boot_id"] = "00000000-0000-0000-0000-000000000000"
	raw, err = json.Marshal(start)
	if err != nil {
		t.Fatal(err)
	}
	// Model an immutable earlier-boot observation in this disposable projection.
	startSQL(t, f, "UPDATE check_starts SET start_json='"+strings.ReplaceAll(string(raw), "'", "''")+"' WHERE id='"+record.Start.ID+"'")
	got := f.run("check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "after-reboot")
	if got.err != nil || !strings.Contains(got.stdout, "⚠ "+record.Start.ID+" never finished; its process has stopped") || !strings.Contains(got.stdout, "✓ lean passed") {
		t.Fatalf("fresh check after reboot=%#v", got)
	}
	shown := f.run("--json", "show", "FID-9f13aa")
	if shown.err != nil || !strings.Contains(shown.stdout, record.Start.ID) || !strings.Contains(shown.stdout, `"result":null`) {
		t.Fatalf("unfinished start lost=%#v", shown)
	}
	retry := f.run("--json", "check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "prior-boot")
	assertJSONFailureCode(t, retry.stderr, "check.execute", "check-incomplete")
}
