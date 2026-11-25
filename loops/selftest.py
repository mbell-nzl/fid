"""V9: exercise the process boundary with isolated git, fid, gh and agents."""
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

SOURCE = Path(__file__).resolve().parent


def stub(name):
    root = Path(os.environ["STUB_ROOT"])
    scenario = os.environ["STUB_SCENARIO"]
    args = sys.argv[2:]
    with (root / "calls").open("a") as log:
        log.write(json.dumps([name, str(Path.cwd()), args]) + "\n")
    if name == "agent":
        import time
        time.sleep(0.2)
        Path(os.environ["LOOP_REPORT"]).write_text("Report\n")
        return
    if name == "gh":
        print("https://github.com/test/project/pull/1")
        return
    assert args.pop(0) == "--json", args
    command = args[0]
    operation = {"next": "next.claim" if "--claim" in args else "next.get",
                 "done": "task.done", "start": "task.start",
                 "heartbeat": "hold.heartbeat", "park": "hold.park",
                 "show": "issue.show", "backup": "project.backup",
                 "review": "review.list"}[command]
    result = {}
    if command == "next":
        kind = "new" if scenario == "heartbeat" else "none"
        result = {"kind": kind, "task": {"id": "FID-123abc"} if kind == "new" else None,
                  "hold": {"id": "HLD-123abc"} if kind == "new" and "--claim" in args else None,
                  "waiting": {"kind": "waiting-for-you"} if kind == "none" else None, "answer": None}
        if scenario == "unknown":
            result["kind"] = "future-work"
        if scenario == "unknown-wait":
            result["waiting"] = {"kind": "future-wait"}
        if scenario == "missing":
            del result["answer"]
    if command == "done":
        result = {"outcomes": [{"task_id": "FID-123abc", "status": "waiting", "refusal": {"code": "done-not-merged"}},
                               {"task_id": "FID-456abc", "status": "completed", "completion": {"task_id": "FID-456abc"}},
                               {"task_id": "FID-789abc", "status": "refused", "refusal": {"code": "done-open-major"}},
                               {"task_id": "FID-abcdef", "status": "github-unavailable", "refusal": {"code": "github-unavailable"}}]}
    if command == "show":
        result = {"issue": {"id": "FID-123abc", "title": "Test task", "spec_md_path": None, "acceptance_criteria": []}}
    if command == "review":
        result = {"task_id": "FID-123abc", "rounds": []}
    if command == "start":
        result = {"task_id": "FID-123abc", "hold_id": "HLD-123abc", "status": "in_progress"}
    if command == "park":
        result = {"task_id": "FID-123abc", "hold_id": "HLD-123abc", "mode": "interrupted", "kept_for_worktree": True, "reason": "failure"}
    if command == "backup":
        result = {"backup": {"outcome": "pushed"}}
    if command == "heartbeat":
        count = root / "heartbeats"
        if count.exists():
            print(json.dumps({"contract_version": "fid-cli.v1", "ok": False, "operation": operation,
                              "error": {"code": "hold-not-current", "retryable": False}}), file=sys.stderr)
            sys.exit(4)
        count.touch()
        result = {"task_id": "FID-123abc", "hold_id": "HLD-123abc", "last_heartbeat_at": "2026-11-25T00:00:00Z",
                  "expires_at": "2026-11-25T01:00:00Z", "heartbeat_interval_seconds": 0.01}
    if command == "next" and scenario in ("empty", "malformed"):
        print("" if scenario == "empty" else "{bad-json")
        return
    print(json.dumps({"contract_version": "fid-cli.future" if scenario == "version" else "fid-cli.v1",
                      "ok": True, "operation": operation, "result": result}))


class Selftest(unittest.TestCase):
    def run_loop(self, scenario):
        with tempfile.TemporaryDirectory(prefix="fid-loop-selftest-") as tmp:
            root = Path(tmp)
            subprocess.run(["git", "init", "-q", "--template=", "-b", "main", tmp], check=True)
            subprocess.run(["git", "-C", tmp, "-c", "user.name=fid test", "-c", "user.email=fid@example.test",
                            "commit", "-qm", "fixture", "--allow-empty"], check=True)
            subprocess.run(["git", "-C", tmp, "update-ref", "refs/remotes/origin/main", "HEAD"], check=True)
            subprocess.run(["git", "-C", tmp, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main"], check=True)
            bin_dir = root / ".git" / "stubs"
            bin_dir.mkdir()
            for name in ("fid", "gh", "agent", "codex", "claude"):
                script = bin_dir / name
                script.write_text(f'#!/bin/sh\nexec python3 "{Path(__file__).resolve()}" {name} "$@"\n')
                script.chmod(0o700)
            env = {**os.environ, "PATH": str(bin_dir) + os.pathsep + os.environ["PATH"],
                   "STUB_ROOT": tmp, "STUB_SCENARIO": scenario, "LOOP_FID": str(bin_dir / "fid"),
                   "LOOP_BUILDER_COMMAND": str(bin_dir / "agent"), "LOOP_REVIEWER_COMMAND": str(bin_dir / "agent"),
                   "LOOP_GH_REPO": "test/project", "XDG_CONFIG_HOME": str(root / "config")}
            result = subprocess.run(["python3", str(SOURCE / "loop.py")], cwd=tmp, env=env,
                                    capture_output=True, text=True, timeout=15)
            calls = [json.loads(line) for line in (root / "calls").read_text().splitlines()] if (root / "calls").exists() else []
            return result, calls

    def test_guards(self):
        for scenario, diagnostic in (("version", "fid-cli.future"), ("unknown", "future-work"),
                                     ("unknown-wait", "future-wait"), ("empty", "empty"),
                                     ("malformed", "{bad-json"), ("missing", "answer"),
                                     ("heartbeat", "hold-not-current")):
            with self.subTest(scenario=scenario):
                result, calls = self.run_loop(scenario)
                self.assertNotEqual(result.returncode, 0, result.stdout)
                self.assertIn(diagnostic, result.stderr)
                if scenario == "heartbeat":
                    self.assertGreaterEqual(sum(c[0] == "fid" and "heartbeat" in c[2] for c in calls), 2)
                else:
                    self.assertFalse(any(c[0] == "agent" for c in calls))

    def test_none_and_every_done_outcome(self):
        result, calls = self.run_loop("none")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("waiting-for-you", result.stdout)
        for status in ("waiting", "completed", "refused", "github-unavailable"):
            self.assertIn(status, result.stdout)
        self.assertEqual(sum(c[0] == "fid" and "done" in c[2] for c in calls), 2)
        self.assertFalse(any(c[0] == "agent" for c in calls))


if __name__ == "__main__":
    if len(sys.argv) > 1 and sys.argv[1] in ("fid", "gh", "agent", "codex", "claude"):
        stub(sys.argv[1])
    else:
        unittest.main()
