//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"github.com/mbell-nzl/fid/internal/config"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func p24InstallLeanPolicy(t *testing.T, root string, configured bool) {
	t.Helper()
	file, err := os.OpenFile(filepath.Join(root, ".fid.toml"), os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	text := "\n[checks.pass]\nrun=['true']\ntimeout='15m'\n[checks.fail]\nrun=['false']\ntimeout='15m'\n[checks.slow]\nrun=['sleep','60']\ntimeout='15m'\n"
	if configured {
		text += "[lanes.standard]\nchecks=[]\n"
	}
	if _, err := file.WriteString(text); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-q", "--allow-empty", "-m", "Lean check policy")
	closureRun(t, root, "git", "remote", "add", "origin", root)
	closureRun(t, root, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
	closureRun(t, root, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
}

func p24LeanPair(t *testing.T, code string) p24CodePair {
	t.Helper()
	root, binary := gittest.NewWorkspace(t), integrationFIDBinary(t)
	p21Run(t, binary, root, "init", "--prefix", "FID", "--json")
	p24InstallLeanPolicy(t, root, code != "checks-not-configured")
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, root, append([]string{"--json"}, args...)...)
	}
	created := call("issue", "create", "--kind", "task", "--title", "Check process pair", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "create")
	taken := call("takeover", p21PathString(t, created, "result", "issue", "id"), "--actor", "human:mark", "--reason", "Build process pair")
	hold := p21PathString(t, taken, "result", "hold_id")
	call("start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "start")
	args := []string{"check", "pass", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "check"}
	switch code {
	case "comparison-unavailable":
		closureRun(t, root, "git", "remote", "remove", "origin")
		args = append(args[:1], args[2:]...)
	case "check-unsuccessful":
		args[1] = "fail"
	case "check-not-qualified":
		if err := os.WriteFile(filepath.Join(root, "dirty"), []byte("uncommitted"), 0o600); err != nil {
			t.Fatal(err)
		}
	case "check-incomplete", "check-still-running":
		slow := append([]string{"--json"}, args...)
		slow[2] = "slow"
		command := exec.Command(binary, slow...)
		command.Dir = root
		var stdout, stderr bytes.Buffer
		command.Stdout, command.Stderr = &stdout, &stderr
		if err := command.Start(); err != nil {
			t.Fatal(err)
		}
		var pid int
		conn := p22FaultDB(t, root)
		deadline := time.Now().Add(10 * time.Second)
		for time.Now().Before(deadline) {
			var raw string
			if err := conn.QueryRow(`SELECT process_json FROM check_processes LIMIT 1`).Scan(&raw); err == nil {
				var process struct {
					PID int `json:"pid"`
				}
				if err := json.Unmarshal([]byte(raw), &process); err != nil {
					t.Fatal(err)
				}
				if _, err := os.Stat(filepath.Join(config.LedgerDir(root), "write.lock")); os.IsNotExist(err) {
					observed, err := exec.Command("ps", "-p", strconv.Itoa(process.PID), "-o", "args=").Output()
					if err == nil && strings.HasPrefix(strings.TrimSpace(string(observed)), "sleep ") {
						pid = process.PID
						break
					}
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
		if err := command.Process.Kill(); err != nil {
			t.Fatal(err)
		}
		_ = command.Wait()
		if pid == 0 {
			t.Fatalf("no child identity: %s %s", stdout.String(), stderr.String())
		}
		t.Cleanup(func() { _ = syscall.Kill(-pid, syscall.SIGKILL) })
		if code == "check-incomplete" {
			args[1] = "slow"
		} else {
			args[len(args)-1] = "new-key"
		}
	}
	return p24CodePair{root, "check.execute", args}
}

func p24LeanCodePair(t *testing.T, code string) {
	t.Helper()
	pair := p24LeanPair(t, code)
	binary := integrationFIDBinary(t)
	stdout, stderr, err := closureRunCommand(pair.root, binary, append([]string{"--json"}, pair.args...)...)
	if err == nil || stdout != "" {
		t.Fatalf("JSON %s = %v/%s/%s", code, err, stdout, stderr)
	}
	envelope := p23Envelope(t, stderr)
	failure := envelope["error"].(map[string]any)
	if failure["code"] != code || envelope["operation"] != "check.execute" || envelope["result"] != nil {
		t.Fatalf("failure=%#v", envelope)
	}
	status := processExitCode(t, err)
	wantStatus := 4
	if code == "check-unsuccessful" {
		wantStatus = 8
	}
	if status != wantStatus {
		t.Fatalf("status=%d", status)
	}
	stdout, stderr, err = closureRunCommand(pair.root, binary, pair.args...)
	if err == nil || processExitCode(t, err) != status || stderr != p24ExpectedFailureText(t, envelope, pair.operation, false) {
		t.Fatalf("text %s = %v/%q/%q, want %q", code, err, stdout, stderr, p24ExpectedFailureText(t, envelope, pair.operation, false))
	}
	if code == "check-unsuccessful" && !strings.Contains(stdout, "✗ fail failed") {
		t.Fatalf("missing failure summary: %s", stdout)
	}
	if code == "check-not-qualified" && !strings.Contains(stdout, "the tree had uncommitted changes") {
		t.Fatalf("missing qualification summary: %s", stdout)
	}
	if code == "check-still-running" {
		details := failure["details"].(map[string]any)
		pid := int(details["pid"].(float64))
		_ = syscall.Kill(-pid, syscall.SIGKILL)
		// A zombie has stopped too; the next check must warn and execute.
		time.Sleep(20 * time.Millisecond)
		args := append([]string{}, pair.args...)
		args[len(args)-1] = "after-stop"
		out, errs, err := closureRunCommand(pair.root, binary, args...)
		if err != nil || errs != "" || !strings.Contains(out, "never finished; its process has stopped") || !strings.Contains(out, "✓ pass passed") {
			t.Fatalf("stopped process=%v/%s/%s, pid=%d", err, out, errs, pid)
		}
	}
}

func TestLeanCheckFailureCodePairs(t *testing.T) {
	t.Parallel()
	for _, code := range []string{"comparison-unavailable", "checks-not-configured", "check-unsuccessful", "check-not-qualified", "check-still-running", "check-incomplete"} {
		t.Run(code, func(t *testing.T) { t.Parallel(); p24LeanCodePair(t, code) })
	}
}

func p24LeanPairName(code string) func(*testing.T) {
	return func(t *testing.T) { t.Helper(); p24LeanCodePair(t, code) }
}
