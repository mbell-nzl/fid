package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/evidence"
)

func TestLiveOutputGoldenFixtures(t *testing.T) {
	result := runLiveGoldenFixture(t)
	machine := strings.Join(result.JSONLines(), "\n") + "\n"
	human := renderLiveHuman(result, "unsatisfied")
	assertGoldenFixture(t, "live-frames.jsonl", machine)
	assertGoldenFixture(t, "human-session.txt", human)
	if !strings.Contains(machine, `"stream":"stdout"`) || !strings.Contains(machine, `"stream":"stderr"`) {
		t.Fatal("machine fixture does not contain both process streams")
	}
	for _, kind := range []string{"redaction-applied", "artifact-rejected", "output-truncated"} {
		if !strings.Contains(machine, kind) || !strings.Contains(human, kind) {
			t.Fatalf("fixtures omit control kind %q", kind)
		}
	}
	if !strings.Contains(machine, `"outcome":"failed"`) || !strings.Contains(human, "outcome=failed") {
		t.Fatal("fixtures do not contain non-success terminal outcome")
	}
	for _, source := range []string{"CANARY-val-036", "literal-val-036"} {
		if strings.Contains(machine, source) || strings.Contains(human, source) {
			t.Fatalf("fixtures leak redaction source %q", source)
		}
	}
	for _, escaped := range []string{`\u001b`, `\r`, `\u0007`} {
		if !strings.Contains(machine, escaped) || !strings.Contains(human, escaped) {
			t.Fatalf("fixtures do not safely render %q", escaped)
		}
	}
}

func runLiveGoldenFixture(t *testing.T) evidence.LiveResult {
	t.Helper()
	workspace := t.TempDir()
	t.Setenv("FID_CREDENTIAL_GATE_TOKEN", "CANARY-val-036")
	literal := "literal-val-036"
	definition := evidence.GateDefinition{
		ArtifactBytes:     1024,
		ArtifactPaths:     []string{"missing-artifact.txt"},
		Argv:              []string{os.Args[0], "-test.run=^TestLiveGoldenChild$", "--", "--live-golden-child"},
		CredentialNames:   []string{"GATE_TOKEN"},
		CWD:               ".",
		DefinitionID:      "val-036-live-output",
		DefinitionVersion: 1,
		Environment: []evidence.EnvironmentEntry{
			{Name: "GATE_TOKEN", Source: "credential"},
			{Name: "LITERAL", Source: "literal", Value: &literal},
		},
		ExcerptBytes:     96,
		Network:          "allowed",
		OutputBytes:      128,
		RedactLiterals:   []string{"literal-val-036"},
		Sandbox:          "none",
		Shell:            false,
		SuccessExitCodes: []int{0},
		TimeoutMS:        5000,
		Write:            "allowed",
	}
	result, err := evidence.RunLive(context.Background(), workspace, definition, "EXE-val-036")
	if err != nil {
		t.Fatalf("RunLive() error = %v", err)
	}
	if err := result.AppendTerminal(evidence.LiveTerminal{
		AttemptID: "ATT-val-036",
		Outcome:   result.Outcome,
		ExitCode:  result.ExitCode,
		TimedOut:  result.TimedOut,
		Response: map[string]any{
			"contract_version": "fid-cli.v1",
			"error": map[string]any{
				"category": "gate-unsuccessful", "code": "gate-unsuccessful", "details": map[string]any{"attempt_id": "ATT-val-036", "exit_code": 23, "outcome": "failed", "receipt_id": nil, "signal": nil, "timed_out": false},
				"message": "gate process was unsuccessful", "retryable": false, "subject": map[string]any{"id": "", "type": "receipt"},
			},
			"ok": false, "operation": "gate.run",
		},
	}); err != nil {
		t.Fatalf("AppendTerminal() error = %v", err)
	}
	return result
}

func TestLiveGoldenChild(t *testing.T) {
	if !strings.Contains(strings.Join(os.Args, "\x00"), "--live-golden-child") {
		return
	}
	_, _ = os.Stdout.WriteString("first stdout CANARY-val-036\x1b[31m\r")
	time.Sleep(20 * time.Millisecond)
	_, _ = os.Stderr.WriteString("first stderr literal-val-036\a\n")
	time.Sleep(20 * time.Millisecond)
	_, _ = os.Stdout.WriteString("second stdout")
	time.Sleep(20 * time.Millisecond)
	_, _ = os.Stderr.WriteString("second stderr")
	time.Sleep(100 * time.Millisecond)
	_, _ = os.Stdout.WriteString(strings.Repeat("x", 120))
	os.Exit(23)
}

func assertGoldenFixture(t *testing.T, name, want string) {
	t.Helper()
	path := filepath.Join("testdata", "val-036", name)
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v; generated=%q", path, err, want)
	}
	if string(got) != want {
		t.Fatalf("golden %s mismatch\nwant: %q\ngot:  %q", path, want, got)
	}
}
