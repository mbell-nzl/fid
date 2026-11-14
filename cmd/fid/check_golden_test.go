package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkrun"
)

var checkGoldenID = regexp.MustCompile(`\b(?:CHK|ART)-[0-9a-f]{6}\b`)

func checkGolden(t *testing.T, f *statusFixture, name, output string) {
	t.Helper()
	output = regexp.MustCompile(`\b[0-9a-f]{40}\b`).ReplaceAllString(output, "<commit-or-tree>")
	output = checkGoldenID.ReplaceAllStringFunc(output, func(id string) string { return id[:3] + "-<new>" })
	output = regexp.MustCompile(`pid [0-9]+|pid=[0-9]+|"pid":[0-9]+`).ReplaceAllStringFunc(output, func(value string) string { return regexp.MustCompile(`[0-9]+`).ReplaceAllString(value, "10001") })
	lifecycleGolden(t, f, filepath.Join("..", "check", name), output)
}

func TestLeanShowGolden(t *testing.T) {
	t.Parallel()
	f := leanJourneyFixture(t, leanHelperDefinition(t, "pass", "1m", []string{"coverage.out"}))
	first := f.run("--json", "check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "complete")
	if first.err != nil {
		t.Fatalf("check=%#v", first)
	}
	options := leanOptions(f, "unfinished")
	options.Before = func(string) error { return errors.New("interrupted") }
	if _, err := checkrun.Run(context.Background(), options); err == nil {
		t.Fatal("expected interruption")
	}
	for _, jsonMode := range []bool{false, true} {
		args := []string{"show", "FID-9f13aa"}
		ext := ".txt"
		if jsonMode {
			args = append([]string{"--json"}, args...)
			ext = ".json"
		}
		got := f.run(args...)
		if got.err != nil {
			t.Fatalf("show=%#v", got)
		}
		output := got.stdout
		binary, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		output = strings.ReplaceAll(output, binary, "<helper>")
		output = regexp.MustCompile(`sha256:[0-9a-f]{64}`).ReplaceAllString(output, "sha256:<digest>")
		output = regexp.MustCompile(`"[0-9a-f]{40}"`).ReplaceAllString(output, `"<commit-or-tree>"`)
		for _, field := range []string{"machine", "boot_id", "start_time", "os", "arch"} {
			output = regexp.MustCompile(`"`+field+`":"[^"]*"`).ReplaceAllString(output, `"`+field+`":"<`+field+`>"`)
		}
		checkGolden(t, f, "show"+ext, output)
	}
}

func TestLeanCheckGolden(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, mode, timeout string
		dirty               bool
	}{
		{"passed", "pass", "15m", false}, {"failed", "fail", "15m", false}, {"timed-out", "wait", "100ms", false}, {"signaled", "signal", "15m", false}, {"dirty", "pass", "15m", true}, {"changed", "change", "15m", false},
		{"spawn-failed", "spawn-failed", "15m", false},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			definition := leanHelperDefinition(t, row.mode, row.timeout, []string{"coverage.out"})
			if row.mode == "spawn-failed" {
				definition = "\n[checks.lean]\nrun=['fid-check-missing-executable']\ntimeout='15m'\nartifacts=['coverage.out']\n"
			}
			f := leanJourneyFixture(t, definition)
			if row.dirty {
				writeBackupTestFile(t, filepath.Join(f.workspace, "dirty"), "dirty")
			}
			args := []string{"check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "golden"}
			for _, jsonMode := range []bool{false, true} {
				invoke := args
				ext := ".txt"
				if jsonMode {
					invoke = append([]string{"--json"}, args...)
					ext = ".json"
				}
				got := f.run(invoke...)
				output := got.stdout
				if jsonMode && got.err != nil {
					output = got.stderr
				}
				if !jsonMode {
					output += got.stderr
				}
				checkGolden(t, f, row.name+ext, output)
			}
		})
	}
}

func TestLeanCheckRefusalGolden(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"unknown", "checks-not-configured", "check-incomplete", "check-still-running", "missing-key", "not-started", "policy-unavailable", "service-actor"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f := leanJourneyFixture(t, leanHelperDefinition(t, "pass", "1m", nil))
			args := []string{"check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "refusal"}
			switch name {
			case "service-actor":
				args[5] = "service:runner"
			case "unknown":
				args[1] = "unknown"
			case "missing-key":
				args = args[:len(args)-2]
			case "not-started":
				startSQL(t, f, `UPDATE issues SET status='todo' WHERE id='FID-9f13aa'`)
			case "checks-not-configured", "policy-unavailable":
				path := filepath.Join(f.workspace, ".fid.toml")
				raw, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				text := strings.Replace(string(raw), "[lanes.standard]\nchecks = []", "[lanes.standard]", 1)
				if name == "policy-unavailable" {
					text += "\nunknown=1\n"
				}
				if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
					t.Fatal(err)
				}
				policyCommit(t, f.workspace)
				identityGit(t, f.workspace, "push", "-q", "origin", "trunk")
			case "check-incomplete", "check-still-running":
				options := leanOptions(f, "refusal")
				options.Before = func(string) error { return errors.New("interrupted before spawn") }
				if _, err := checkrun.Run(context.Background(), options); err == nil {
					t.Fatal("expected interruption")
				}
				if name == "check-still-running" {
					args[len(args)-1] = "another"
				}
			}
			for _, jsonMode := range []bool{false, true} {
				invoke := args
				ext := ".txt"
				if jsonMode {
					invoke = append([]string{"--json"}, args...)
					ext = ".json"
				}
				got := f.run(invoke...)
				if got.err == nil || got.stdout != "" {
					t.Fatalf("refusal=%#v", got)
				}
				checkGolden(t, f, name+ext, got.stderr)
			}
		})
	}
}
