"""V9: exercise the process boundary with isolated git, fid, gh and agents."""
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

SOURCE = Path(__file__).resolve().parent
ANSWER_SCENARIOS = ("answers-missing", "answers-invalid", "answer-refused", "answer-wrong-operation", "answer-missing-field", "answer-before-check", "answer-other-actor")


def stub(name):
    root = Path(os.environ["STUB_ROOT"])
    scenario = os.environ["STUB_SCENARIO"]
    args = sys.argv[2:]
    with (root / "calls").open("a") as log:
        log.write(json.dumps([name, str(Path.cwd()), args]) + "\n")
    if name == "agent":
        if args == ["--version"]:
            if os.getenv("STUB_VERSION_FAILURE") == "1":
                sys.exit(2)
            print("codex-cli 0.172.0")
            return
        print(json.dumps({"type":"turn.completed","usage":{"input_tokens":10,"cached_input_tokens":0,"cache_write_input_tokens":0,"output_tokens":2,"reasoning_output_tokens":0}}), flush=True)
        print("stderr diagnostic", file=sys.stderr)
        if scenario == "judge":
            context=json.loads(Path(os.environ["LOOP_CONTEXT"]).read_text())
            assert os.environ["LOOP_ACTOR"]=="agent:sol-judge"
            assert context["receipts"]["issue"]["id"]=="FID-123abc"
            assert "judge sentinel" not in Path(os.environ["LOOP_CONTEXT"]).read_text()
            Path(os.environ["LOOP_REVIEW"]).write_text(json.dumps({"stream":"judge","verdict":"HOLD"}))
            return
        import time
        if scenario == "slow-heartbeat":
            deadline = time.monotonic() + 3
            while not (root / "renewing").exists() and time.monotonic() < deadline:
                time.sleep(0.01)
            assert (root / "renewing").exists(), "renewal did not start"
            sys.exit(42)
        time.sleep(0.2)
        Path(os.environ["LOOP_REPORT"]).write_text("Report\n")
        if scenario in ANSWER_SCENARIOS and scenario != "answers-missing":
            Path(os.environ["LOOP_ANSWERS"]).write_text("{bad-json" if scenario == "answers-invalid" else
                json.dumps([{"finding_id": "RND-123abc-1", "answer": "Fixed in `abc123`."}]))
        return
    if name == "gh":
        print("https://github.com/test/project/pull/1")
        return
    assert args.pop(0) == "--json", args
    command = args[0]
    operation = {"next": "next.claim" if "--claim" in args else "next.get",
                 "done": "task.done", "start": "task.start",
                 "heartbeat": "hold.heartbeat", "park": "hold.park",
                 "show": "issue.show", "status": "project.status", "backup": "project.backup",
                 "review": "review." + (args[1] if command == "review" else ""), "usage": "usage.record", "finding": "finding." + (args[1] if command == "finding" else ""), "check": "check.execute"}[command]
    result = {}
    if command == "next":
        kind = "new" if scenario in ("heartbeat", "missing-nested", "slow-heartbeat", "claim-no-hold", "work-no-task", "park-failure", "backup-failure", "initial-heartbeat") else "none"
        if scenario in ANSWER_SCENARIOS and Path.cwd().name == "build":
            kind = "rework"
        if scenario == "judge" and Path.cwd().name != "build" and not (root / "judged").exists():
            kind="judge"
        if scenario == "role-change":
            kind = "new" if "--claim" in args or Path.cwd().name == "build" else "review"
        result = {"kind": kind, "task": {"id": "FID-123abc"} if kind in ("new", "rework") else None,
                  "hold": {"id": "HLD-123abc"} if kind in ("new", "rework") and "--claim" in args else None,
                  "waiting": {"kind": "waiting-for-you"} if kind == "none" else None, "answer": None}
        if kind in ("review", "judge"):
            result["task"] = {"id": "FID-123abc"}
            result["review"] = {"handoff": {"id": "HOF-123abc", "head": "unused", "base": "unused", "pr": "unused"},
                                "stream": "judge" if kind=="judge" else "primary", "required_profiles": ["merge"] if kind=="judge" else ["code"]}
            if kind=="judge":
                result["review"]["handoff"]["head"]=subprocess.check_output(["git","rev-parse","HEAD"],text=True).strip()
                if "--claim" in args:result["hold"]={"id":"RCL-123abc"}
        if scenario == "none-task":
            result["task"] = {"id": "FID-123abc"}
        if scenario == "none-hold":
            result["hold"] = {"id": "HLD-123abc"}
        if scenario == "work-no-task":
            result["task"] = None
        if scenario == "claim-no-hold" and "--claim" in args:
            result["hold"] = None
        if scenario == "unknown":
            result["kind"] = "future-work"
        if scenario == "unknown-wait":
            result["waiting"] = {"kind": "future-wait"}
        if scenario == "missing":
            del result["answer"]
        if scenario == "missing-nested":
            result["task"] = {}
    if command == "done":
        result = {"outcomes": [{"task_id": "FID-123abc", "status": "waiting", "refusal": {"code": "done-not-merged"}},
                               {"task_id": "FID-456abc", "status": "completed", "completion": {"task_id": "FID-456abc"}},
                               {"task_id": "FID-789abc", "status": "refused", "refusal": {"code": "done-open-major"}},
                               {"task_id": "FID-abcdef", "status": "github-unavailable", "refusal": {"code": "github-unavailable"}}]}
        if scenario == "unknown-done":
            result["outcomes"][0]["status"] = "future-outcome"
    if command == "show":
        result = {"issue": {"id": "FID-123abc", "title": "Test task", "spec_md_path": None, "acceptance_criteria": []}, "handoff": None}
    if command == "status":
        result = {"needs_you": []}
    if command == "review":
        result = {"task_id": "FID-123abc", "rounds": []}
        if scenario in ANSWER_SCENARIOS:
            finding = {"id": "RND-123abc-1", "profile": "code", "severity": "major", "summary": "Retry", "location": None,
                       "evidence": "Missing retry", "fix": "Cover retry", "open": True}
            if scenario == "answer-other-actor":
                finding["responses"] = [{"actor": "human:mark", "text": "Please fix this."}]
            result["rounds"] = [{"id": "RND-123abc", "handoff_id": "HOF-123abc", "stream": "primary", "profiles": ["code"],
                                 "verdict": "CHANGES", "current": True, "findings": [finding]}]
        if args[1]=="record":
            assert scenario=="judge"
            (root / "judged").touch()
            result={"round":{"id":"RND-123abc","task_id":"FID-123abc","handoff_id":"HOF-123abc","stream":"judge","profiles":["merge"],"verdict":"","findings":[],"sealed":True}}
    if command == "finding":
        result = {"response": {"id": "FRS-123abc", "finding_id": "RND-123abc-1", "task_id": "FID-123abc", "handoff_id": "HOF-123abc",
                               "actor": "agent:sol-builder", "text": "Fixed in `abc123`.", "at": "2026-11-28T00:00:00Z"}}
        if scenario == "answer-wrong-operation":
            operation = "finding.defer"
        if scenario == "answer-missing-field":
            del result["response"]["text"]
    if command == "check" or command == "finding" and scenario == "answer-refused":
        print(json.dumps({"contract_version": "fid-cli.v2", "ok": False, "operation": operation,
                          "error": {"code": "check-probe-stopped" if command == "check" else "invalid-argument", "retryable": False}}), file=sys.stderr)
        sys.exit(4)
    if command == "usage":
        if os.getenv("STUB_USAGE_FAILURE") == "1":
            print(json.dumps({"contract_version":"fid-cli.v2", "operation":operation, "ok":False, "error":{"code":"invalid-argument", "retryable":False}}), file=sys.stderr)
            sys.exit(2)
        stdout = Path(args[args.index("--tool-output") + 1]).read_text()
        assert "stderr diagnostic" not in stdout
        assert json.loads(stdout)["type"] == "turn.completed"
        if os.getenv("STUB_VERSION_FAILURE") != "1":
            assert args[args.index("--tool-version") + 1] == "codex-cli 0.172.0"
        result = {"usage": {"id":"USG-123abc", "task_id":"FID-123abc", "hold_id":"HLD-123abc", "actor":"agent:sol-builder", "model":"openai/sol", "role":"build", "stream":"", "format":"codex", "counts":{"uncached_input":10,"cache_read":0,"cache_write":0,"output":2,"reasoning":0,"total":12}, "models":{}, "reported":[], "input_digest":"sha256:"+"a"*64, "parser_revision":"codex/1", "tool_version":"codex-cli 0.172.0", "version_source":"caller", "recorded_at":"2026-11-25T00:00:00Z"}}
    if command == "start":
        result = {"task_id": "FID-123abc", "hold_id": "HLD-123abc", "status": "in_progress"}
    if command == "park":
        result = {"task_id": "FID-123abc", "hold_id": "HLD-123abc", "mode": "interrupted", "kept_for_worktree": True, "reason": "failure"}
    if command == "backup":
        result = {"backup": {"outcome": "pushed"}}
    if (scenario in ("park-failure", "backup-failure") and command == "heartbeat" or
            scenario == "park-failure" and command == "park" or
            scenario == "backup-failure" and command == "backup"):
        print(json.dumps({"contract_version": "fid-cli.v2", "ok": False, "operation": operation,
                          "error": {"code": command + "-probe-failed", "retryable": False}}), file=sys.stderr)
        sys.exit(4)
    if command == "heartbeat":
        count = root / "heartbeats"
        if count.exists() and scenario == "slow-heartbeat":
            import time
            (root / "renewing").touch()
            time.sleep(0.3)
            with (root / "calls").open("a") as log:
                log.write(json.dumps(["heartbeat-finished", str(Path.cwd()), []]) + "\n")
        elif count.exists():
            print(json.dumps({"contract_version": "fid-cli.v2", "ok": False, "operation": operation,
                              "error": {"code": "hold-not-current", "retryable": False}}), file=sys.stderr)
            sys.exit(4)
        count.touch()
        result = {"task_id": "FID-123abc", "hold_id": "HLD-123abc", "last_heartbeat_at": "2026-11-25T00:00:00Z",
                  "expires_at": "2026-11-25T01:00:00Z", "heartbeat_interval_seconds": 300 if scenario in ANSWER_SCENARIOS or scenario=="judge" else 0.01}
    if command == "next" and scenario in ("empty", "malformed"):
        print("" if scenario == "empty" else "{bad-json")
        return
    if command == "next" and scenario == "wrong-operation" and "--claim" not in args:
        operation = "next.claim"
    if command == "done" and scenario == "empty-done":
        return
    if scenario == "nonzero":
        print(json.dumps({"contract_version": "fid-cli.v2", "ok": False, "operation": operation,
                          "error": {"code": "internal-error", "retryable": False}}), file=sys.stderr)
        sys.exit(7)
    if command == "next" and scenario == "result-null":
        result = None
    reply = {"contract_version": "fid-cli.future" if scenario == "version" else "fid-cli.v2",
             "ok": scenario != "false-ok", "operation": operation, "result": result}
    print(json.dumps(reply), file=sys.stderr if scenario == "nonzero-success" else sys.stdout)
    if scenario == "nonzero-success":
        sys.exit(7)


class Selftest(unittest.TestCase):
    def test_judge_role(self):
        result,calls=self.run_loop("judge",usage="none")
        self.assertEqual(result.returncode,0,result.stderr)
        claims=[c[2] for c in calls if c[0]=="fid" and "--claim" in c[2]]
        self.assertIn("agent:sol-judge",claims[0])
        self.assertTrue(any(c[0]=="fid" and c[2][1:3]==["review","record"] for c in calls))
        self.assertFalse(any(c[0]=="fid" and c[2][1] in ("check","handoff","start","park") for c in calls))

    def test_judge_markdown_guidance(self):
        prompt = " ".join((SOURCE / "prompts/judge.md").read_text().split())
        for instruction in ("GitHub Markdown", "in every text field", "`backticks`", "whole diff", "MERGE", "HOLD", "shadow", "LOOP_REVIEW"):
            with self.subTest(instruction=instruction):
                self.assertIn(instruction, prompt)

    def test_reviewer_markdown_guidance(self):
        prompt = " ".join((SOURCE / "prompts/reviewer.md").read_text().split())
        for instruction in ("GitHub Markdown", "in every text field", "commands, file paths, test names, flags, identifiers and literal output",
                            "`backticks`", "short Markdown list", "one command per item with its result", "concise evidence"):
            with self.subTest(instruction=instruction):
                self.assertIn(instruction, prompt)

    def run_loop(self, scenario, usage="codex", version_failure=False, usage_failure=False):
        with tempfile.TemporaryDirectory(prefix="fid-loop-selftest-") as tmp:
            root = Path(tmp) / "build"
            root.mkdir()
            subprocess.run(["git", "init", "-q", "--template=", "-b", "main", str(root)], check=True)
            subprocess.run(["git", "-C", str(root), "-c", "user.name=fid test", "-c", "user.email=fid@example.test",
                            "commit", "-qm", "fixture", "--allow-empty"], check=True)
            subprocess.run(["git", "-C", str(root), "update-ref", "refs/remotes/origin/main", "HEAD"], check=True)
            subprocess.run(["git", "-C", str(root), "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main"], check=True)
            if scenario in ANSWER_SCENARIOS:
                subprocess.run(["git", "-C", str(root), "switch", "-qc", "reset-test-task"], check=True)
            bin_dir = root / ".git" / "stubs"
            bin_dir.mkdir()
            for name in ("fid", "gh", "agent", "codex", "claude"):
                script = bin_dir / name
                script.write_text(f'#!/bin/sh\nexec python3 "{Path(__file__).resolve()}" {name} "$@"\n')
                script.chmod(0o700)
            env = {**os.environ, "PATH": str(bin_dir) + os.pathsep + os.environ["PATH"],
                   "STUB_ROOT": str(root / ".git"),
                   "LOOP_JUDGE_ACTOR": "agent:sol-judge", "LOOP_JUDGE_MODEL": "openai/gpt-6.1-sol",
                   "LOOP_BUILDER_USAGE": usage, "LOOP_REVIEWER_USAGE": usage, "LOOP_JUDGE_USAGE": "none",
                   "STUB_VERSION_FAILURE": str(int(version_failure)), "STUB_USAGE_FAILURE": str(int(usage_failure)), "STUB_SCENARIO": scenario, "LOOP_FID": str(bin_dir / "fid"),
                   "LOOP_BUILDER_COMMAND": str(bin_dir / "agent"), "LOOP_REVIEWER_COMMAND": str(bin_dir / "agent"), "LOOP_JUDGE_COMMAND": str(bin_dir / "agent"),
                   "LOOP_GH_REPO": "test/project", "XDG_CONFIG_HOME": str(root / "config")}
            env["LOOP_EPIC"] = ""
            if scenario == "same-actors":
                env["LOOP_REVIEWER_ACTOR"] = env["LOOP_BUILDER_ACTOR"] = "agent:same"
            if scenario == "epic":
                env["LOOP_EPIC"] = "FID-abcdef"
            result = subprocess.run(["python3", str(SOURCE / "loop.py")], cwd=root, env=env,
                                    capture_output=True, text=True, timeout=15)
            calls_path = root / ".git" / "calls"
            calls = [json.loads(line) for line in calls_path.read_text().splitlines()] if calls_path.exists() else []
            return result, calls

    def test_finding_recording_guards(self):
        for scenario, diagnostic in (("answers-missing", "answers.json"), ("answers-invalid", "answers.json"),
                                     ("answer-refused", "invalid-argument"), ("answer-wrong-operation", "unexpected operation"),
                                     ("answer-missing-field", "text"), ("answer-before-check", "check-probe-stopped"),
                                     ("answer-other-actor", "check-probe-stopped")):
            with self.subTest(scenario=scenario):
                result, calls = self.run_loop(scenario, usage="none")
                self.assertNotEqual(result.returncode, 0)
                self.assertIn(diagnostic, result.stderr)
                commands = [c[2][1] for c in calls if c[0] == "fid"]
                if scenario in ("answer-before-check", "answer-other-actor"):
                    self.assertLess(commands.index("finding"), commands.index("check"))
                else:
                    self.assertNotIn("check", commands)
                self.assertEqual(commands[-2:], ["park", "backup"])

    def test_guards(self):
        for scenario, diagnostic in (("version", "fid-cli.future"), ("unknown", "future-work"),
                                     ("unknown-wait", "future-wait"), ("empty", "empty"),
                                     ("malformed", "{bad-json"), ("missing", "answer"),
                                     ("wrong-operation", "operation"),
                                     ("missing-nested", "id"), ("empty-done", "empty"),
                                     ("nonzero", "internal-error"),
                                     ("heartbeat", "hold-not-current")):
            with self.subTest(scenario=scenario):
                result, calls = self.run_loop(scenario)
                self.assertNotEqual(result.returncode, 0, result.stdout)
                self.assertIn(diagnostic, result.stderr)
                if scenario == "heartbeat":
                    self.assertGreaterEqual(sum(c[0] == "fid" and "heartbeat" in c[2] for c in calls), 2)
                else:
                    self.assertFalse(any(c[0] == "agent" for c in calls))
                    if scenario != "empty-done":
                        self.assertEqual([c[2][1] for c in calls if c[0] == "fid"], ["next"])

    def test_usage_is_attempted_after_failed_agent_before_park(self):
        result, calls = self.run_loop("slow-heartbeat")
        self.assertNotEqual(result.returncode, 0)
        commands = [c[2][1] for c in calls if c[0] == "fid"]
        self.assertIn("usage", commands)
        self.assertLess(commands.index("usage"), commands.index("park"))
        usage = next(c[2] for c in calls if c[0] == "fid" and c[2][1] == "usage")
        self.assertEqual(usage[2], "record")
        self.assertIn("usage:", result.stdout)
        self.assertNotIn("usage record failed", result.stderr)
        self.assertEqual(usage[usage.index("--tool-version") + 1], "codex-cli 0.172.0")
        self.assertEqual(usage[usage.index("--hold") + 1], "HLD-123abc")
        self.assertEqual(usage[usage.index("--tool") + 1], "codex")

    def test_usage_none_version_failure_and_record_failure(self):
        result, calls = self.run_loop("slow-heartbeat", usage="none")
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse(any(c[0] == "fid" and c[2][1] == "usage" for c in calls))
        result, calls = self.run_loop("slow-heartbeat", version_failure=True)
        usage = next(c[2] for c in calls if c[0] == "fid" and c[2][1] == "usage")
        self.assertNotIn("--tool-version", usage)
        self.assertIn("usage:", result.stdout)
        result, calls = self.run_loop("slow-heartbeat", usage_failure=True)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("usage record failed", result.stderr)
        self.assertIn("agent failed", result.stderr)
        self.assertEqual([c[2][1] for c in calls if c[0] == "fid"][-2:], ["park", "backup"])

    def test_none_and_every_done_outcome(self):
        result, calls = self.run_loop("none")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("waiting-for-you", result.stdout)
        for status in ("waiting", "completed", "refused", "github-unavailable"):
            self.assertIn(status, result.stdout)
        self.assertEqual(sum(c[0] == "fid" and "done" in c[2] for c in calls), 2)
        self.assertFalse(any(c[0] == "agent" for c in calls))

    def test_transport_and_selection_guards(self):
        for scenario, diagnostic in (("nonzero-success", "exit 7"), ("false-ok", "command failure"),
                                     ("result-null", "result must be an object"),
                                     ("none-task", "missing task/hold"), ("none-hold", "missing task/hold"),
                                     ("work-no-task", "missing task/hold"), ("claim-no-hold", "missing task/hold"),
                                     ("unknown-done", "unknown done outcome"),
                                     ("same-actors", "actors must differ"), ("role-change", "work selection changed roles")):
            with self.subTest(scenario=scenario):
                result, calls = self.run_loop(scenario)
                self.assertNotEqual(result.returncode, 0, result.stdout)
                self.assertIn(diagnostic, result.stderr)
                self.assertFalse(any(c[0] == "agent" for c in calls))
                if scenario in ("nonzero-success", "false-ok", "result-null", "none-task", "none-hold", "work-no-task"):
                    self.assertIn("reply=", result.stderr)

    def test_epic_scopes_every_selection(self):
        result, calls = self.run_loop("epic")
        self.assertEqual(result.returncode, 0, result.stderr)
        selections = [c for c in calls if c[0] == "fid" and c[2][1] == "next"]
        self.assertEqual(len(selections), 4)
        self.assertEqual(len({c[1] for c in selections}), 2)
        for _, _, args in selections:
            self.assertIn("--epic", args)
            self.assertEqual(args[args.index("--epic") + 1], "FID-abcdef")

    def test_heartbeat_finishes_before_parking(self):
        result, calls = self.run_loop("slow-heartbeat")
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("agent failed", result.stderr)
        finished = next((i for i, c in enumerate(calls) if c[0] == "heartbeat-finished"), None)
        park = next((i for i, c in enumerate(calls) if c[0] == "fid" and c[2][1] == "park"), None)
        self.assertIsNotNone(finished, calls)
        self.assertIsNotNone(park, calls)
        self.assertLess(finished, park, calls)

    def test_failure_cleanup_reports_errors_and_attempts_backup(self):
        for scenario in ("park-failure", "backup-failure"):
            with self.subTest(scenario=scenario):
                result, calls = self.run_loop(scenario)
                self.assertNotEqual(result.returncode, 0)
                self.assertIn("park/backup failed", result.stderr)
                self.assertIn(scenario.split("-")[0] + "-probe-failed", result.stderr)
                self.assertEqual([c[2][1] for c in calls if c[0] == "fid"][-2:], ["park", "backup"])

    def test_initial_heartbeat_precedes_agent(self):
        result, calls = self.run_loop("initial-heartbeat")
        names = [c[2][1] if c[0] == "fid" else c[0] for c in calls]
        self.assertIn("heartbeat", names)
        self.assertIn("agent", names)
        self.assertLess(names.index("heartbeat"), names.index("agent"), calls)


if __name__ == "__main__":
    if len(sys.argv) > 1 and sys.argv[1] in ("fid", "gh", "agent", "codex", "claude"):
        stub(sys.argv[1])
    else:
        unittest.main()
