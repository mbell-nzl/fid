#!/usr/bin/env python3
"""Agents and GitHub for the first-loop journeys; never contacts a service."""
import json
import os
from pathlib import Path
import subprocess
import sys

state = Path(os.environ["FIRST_LOOP_STATE"])
scenario = os.environ["FIRST_LOOP_SCENARIO"]
args = sys.argv[1:]


def log(**fields):
    with (state / "calls.jsonl").open("a") as out:
        out.write(json.dumps({"cwd": str(Path.cwd()), **fields}) + "\n")


def fid(*argv):
    return json.loads(subprocess.check_output([os.environ["LOOP_FID"], "--json", *argv], text=True))["result"]


def pushed_head():
    branch = subprocess.check_output(["git", "branch", "--show-current"], text=True).strip()
    refs = subprocess.check_output(["git", "ls-remote", "origin", "refs/heads/" + branch], text=True).split()
    return {"head": subprocess.check_output(["git", "rev-parse", "HEAD"], text=True).strip(),
            "remote_head": refs[0] if refs else None}


if Path(sys.argv[0]).name == "gh":
    assert args[:2] != ["pr", "merge"], "loop must never merge"
    if args[:2] == ["pr", "list"]:
        print(json.dumps([{"url": "https://github.com/mbell-nzl/fid/pull/63"}] if (state / "pr").exists() else []))
    elif args[:2] == ["pr", "create"]:
        log(kind="create", draft="--draft" in args, body=Path(args[args.index("--body-file") + 1]).read_text(), **pushed_head())
        (state / "pr").touch()
        print("https://github.com/mbell-nzl/fid/pull/63")
    elif args[:2] == ["pr", "edit"]:
        log(kind="edit", body=Path(args[args.index("--body-file") + 1]).read_text(), **pushed_head())
    elif args[:2] == ["pr", "ready"]:
        if (state / "ready").exists():
            print("PR is already ready", file=sys.stderr)
            sys.exit(1)
        rounds = fid("review", "list", os.environ["FIRST_LOOP_TASK"])["rounds"]
        log(kind="ready", recorded=bool(rounds and rounds[-1]["verdict"] == "APPROVE"))
        (state / "ready").touch()
    elif args[:2] == ["pr", "view"] and args[-1] == "isDraft":
        print(json.dumps({"isDraft": not (state / "ready").exists()}))
    elif args[:2] == ["pr", "comment"]:
        log(kind="comment", body=Path(args[args.index("--body-file") + 1]).read_text())
        print("https://github.com/mbell-nzl/fid/pull/63#issuecomment-123")
    elif args[0] == "api" and any("/comments?" in a for a in args):
        print("[]")
    else:
        print("offline fixture", file=sys.stderr)
        sys.exit(1)
    sys.exit(0)

assert Path(sys.argv[0]).name == "agent", "real agent invocation is forbidden"
# Bound omitted-step mutations that would otherwise claim the same work forever.
calls = state / "calls.jsonl"
if calls.exists() and sum(json.loads(line)["kind"] in ("build", "review") for line in calls.read_text().splitlines()) >= 6:
    print("stub agent session limit: loop repeated work without making progress", file=sys.stderr)
    sys.exit(42)
prompt = sys.stdin.read()
context = json.loads(Path(os.environ["LOOP_CONTEXT"]).read_text())
work = context["work"]
if work["kind"] == "review":
    expected = work["review"]["handoff"]["head"]
    head = subprocess.check_output(["git", "rev-parse", "HEAD"], text=True).strip()
    detached = subprocess.run(["git", "symbolic-ref", "-q", "HEAD"], capture_output=True).returncode != 0
    log(kind="review", prompt=prompt, actor=os.environ["LOOP_ACTOR"], head=head, expected_head=expected, detached=detached,
        stream=work["review"]["stream"], prior_streams=[r["stream"] for r in context["reviews"]["rounds"]])
    count = state / "reviews"
    first = not count.exists()
    count.touch()
    changes = scenario == "changes" and first
    review = {"reviewer": {"actor": os.environ["LOOP_ACTOR"], "model": os.environ["LOOP_MODEL"],
                            "session_id": os.environ["LOOP_SESSION"]}, "stream": work["review"]["stream"],
              "profiles": work["review"]["required_profiles"], "verdict": "CHANGES" if changes else "APPROVE",
              "purpose": "Review first-loop change", "verification": "Inspected retry coverage", "findings": []}
    if changes:
        review["findings"] = [{"profile": "code", "severity": "major", "summary": "Missing retry case",
                               "location": None, "evidence": "Retry is untested", "fix": "Cover the retry"}]
    else:
        review["dispositions"] = [{"finding_id": f["id"], "disposition": "resolved", "reason": "Retry regression passes"}
                                  for r in context["reviews"]["rounds"] for f in r["findings"] if f["open"]]
    Path(os.environ["LOOP_REVIEW"]).write_text(json.dumps(review))
else:
    merging = subprocess.run(["git", "rev-parse", "--verify", "MERGE_HEAD"], capture_output=True).returncode == 0
    log(kind="build", context=context, prompt=prompt, merging=merging)
    first = not (state / "built").exists()
    (state / "built").touch()
    if first and scenario == "waiting":
        Path("unfinished.txt").write_text("Kept until Mark answers.\n")
        fid("park", "--hold", os.environ["LOOP_HOLD"], "--actor", os.environ["LOOP_ACTOR"],
            "--waiting", "Which retries?", "--idempotency-key", "waiting")
        sys.exit(0)
    if first and scenario == "interrupt":
        sys.exit(42)
    if first and scenario == "ctrl-c":
        import time
        time.sleep(30)
    Path("change_test.go").write_text("package example\n")
    Path("change.txt").write_text(os.environ["LOOP_SESSION"] + "\n")
    if first and scenario == "dirty":
        sys.exit(0)
    subprocess.run(["git", "add", "."], check=True)
    subprocess.run(["git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test",
                    "commit", "-qm", "test: stub implementation"], check=True)
    Path(os.environ["LOOP_REPORT"]).write_text("Stub completion report: " + work["kind"] + "\n\n## Rollback\nRevert this change.\n")
