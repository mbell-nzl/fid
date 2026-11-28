#!/usr/bin/env python3
import contextlib, json, os, re, shlex, signal, sys, threading, uuid
from pathlib import Path
import subprocess as sp

HERE = Path(__file__).resolve().parent
CONTRACT = json.loads((HERE.parent / "docs/specs/LOOP_CONTRACT.md").read_text().split("```json loop-contract\n")[1].split("\n```")[0])
VERIFIED_EPICS = set()
ACTIVE, REPO = None, os.getenv("LOOP_GH_REPO", "mbell-nzl/fid")
OPERATIONS = dict(next="next.get", start="task.start", heartbeat="hold.heartbeat", park="hold.park", show="issue.show", status="project.status", backup="project.backup", done="task.done", check="check.execute", handoff="task.handoff", usage="usage.record")
def run(cwd, *args, allowed=(0,)):
    p = sp.run(args, cwd=cwd, capture_output=True, text=True)
    if p.returncode not in allowed: raise RuntimeError(f"{args}: exit {p.returncode}; stdout={p.stdout!r}; stderr={p.stderr!r}")
    return p.stdout.strip()
def required(value, path):
    if not path or value is None: return
    key, array, tail = path[0].removesuffix("[]"), path[0].endswith("[]"), path[1:]
    if key.endswith("?") and key[:-1] not in value: return
    child = value[key.removesuffix("?")]
    for item in child if array else [child]: required(item, tail)
def fid(cwd, *args):
    p = sp.run([os.getenv("LOOP_FID", "fid"), "--json", *args], cwd=cwd, capture_output=True, text=True)
    raw = p.stdout if p.returncode == 0 else p.stderr
    try:
        if not raw.strip(): raise ValueError("empty output")
        reply = json.loads(raw)
        expected = "epic.validation.verify" if args[0] == "epic" else args[0] + "." + args[1] if args[0] in ("review", "finding") else "next.claim" if args[0] == "next" and "--claim" in args else OPERATIONS[args[0]]
        if reply["operation"] != expected: raise ValueError("unexpected operation")
        if reply["contract_version"] != "fid-cli.v2" or p.returncode != 0 or reply["ok"] is not True:
            raise ValueError(f"unsupported contract or command failure (exit {p.returncode})")
        result = reply["result"]
        if not isinstance(result, dict): raise ValueError("result must be an object")
        for path in CONTRACT["fields"][reply["operation"]]:
            required(result, path.split("."))
        if args[0] == "next" and (result["kind"] not in CONTRACT["next_kinds"] or result["waiting"] and result["waiting"]["kind"] not in CONTRACT["waiting_kinds"]):
            raise ValueError("unknown next.kind or waiting.kind")
        if args[0] == "next" and ((result["kind"] == "none") != (result["task"] is None) or "--claim" in args and (result["kind"] == "none") != (result["hold"] is None)): raise ValueError("missing task/hold or work on none")
        return result
    except (ValueError, KeyError, TypeError) as e:
        raise RuntimeError(f"fid {args}: {e}; reply={raw!r}; stdout={p.stdout!r}; stderr={p.stderr!r}") from e
def next_work(cwd, actor, model, claim=False):
    scope = ["--epic", os.environ["LOOP_EPIC"]] if os.getenv("LOOP_EPIC") else []
    return fid(cwd, "next", *scope, "--actor", actor, *(["--claim", "--model", model] if claim else []))
def task_stem(task):
    return Path(task["spec_md_path"]).stem if task["spec_md_path"] else task["id"]
def task_branch(task):
    stem = task_stem(task)
    title = task["title"].removeprefix(stem)
    return "reset-" + stem + "-" + re.sub(r"[^a-z0-9]+", "-", title.lower()).strip("-")
def pr_title(cwd, task):
    base = run(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD")
    subjects = run(cwd, "git", "log", "--reverse", "--format=%s", base + "..HEAD").splitlines()
    if not subjects: raise RuntimeError("PR title requires a commit on the task branch")
    subject = next((s for s in subjects if re.match(r"^feat(?:\([^)]+\))?!?:", s)), subjects[0])
    return subject + " (RESET " + task_stem(task) + ")"
def print_pending(cwd):
    status = fid(cwd, "status")
    for warning in status.get("warnings", []):
        if warning["kind"] == "review-post-pending":
            print(f"review-post-pending: {warning['summary']}; {warning['next_action']}", flush=True)
    return status
def record_answers(cwd, actor, task_id, path):
    try:
        entries = json.loads(path.read_text())
        if not isinstance(entries, list): raise ValueError("expected a JSON list")
        findings = {f["id"]: f for r in fid(cwd, "review", "list", task_id)["rounds"] for f in r["findings"]}
        seen, pending = set(), []
        # Validate the entire file before recording any entry.
        for entry in entries:
            if not isinstance(entry, dict) or set(entry) not in ({"finding_id", "answer"}, {"finding_id", "defer"}):
                raise ValueError("each entry needs finding_id and exactly one of answer or defer")
            verb = "answer" if "answer" in entry else "defer"
            finding_id, text = entry["finding_id"], entry[verb]
            if not isinstance(finding_id, str) or finding_id not in findings or finding_id in seen:
                raise ValueError(f"unknown or repeated finding_id: {finding_id!r}")
            if not isinstance(text, str) or not text.strip(): raise ValueError(f"blank or invalid {verb} for {finding_id}")
            seen.add(finding_id)
            finding = findings[finding_id]
            if any(r["actor"] == actor for r in finding.get("responses", [])): continue
            if not finding["open"]: raise ValueError(f"finding is closed: {finding_id}")
            if verb == "defer" and finding["severity"] not in ("minor", "nit"):
                raise ValueError(f"only a minor or nit can be deferred: {finding_id}")
            pending.append((verb, finding_id, text))
    except (OSError, ValueError) as error:
        raise RuntimeError(f"invalid answers file {path}: {error}") from error
    for verb, finding_id, text in pending:
        response = fid(cwd, "finding", verb, finding_id, "--actor", actor, "--text" if verb == "answer" else "--title", text)
        print("finding:", json.dumps(response), flush=True)
    print_pending(cwd)
def verify_epics(cwd, actor, status):
    for item in status["needs_you"]:
        if item["kind"] != "epic-validation-verify": continue
        epic = item["subject"]["id"]
        if os.getenv("LOOP_EPIC") and epic != os.environ["LOOP_EPIC"]: continue
        if epic in VERIFIED_EPICS: continue
        VERIFIED_EPICS.add(epic)
        result = fid(cwd, "epic", "validation", "verify", epic, "--actor", actor, "--idempotency-key", "loop-validation-" + uuid.uuid4().hex)
        print("epic validation:", json.dumps(result), flush=True)

def done(cwd, actor):
    for outcome in fid(cwd, "done", "--actor", actor)["outcomes"]:
        if outcome["status"] not in ("completed", "waiting", "refused", "github-unavailable"): raise RuntimeError(f"unknown done outcome: {outcome!r}")
        print("done:", json.dumps(outcome), flush=True)
    status = print_pending(cwd)
    verify_epics(cwd, actor, status)
    for item in status["needs_you"]:
        if item["kind"] in ("merge", "approval"):
            pr = fid(cwd, "show", item["subject"]["id"])["handoff"]["pr"]
            if json.loads(run(cwd, "gh", "pr", "view", pr, "--repo", REPO, "--json", "isDraft"))["isDraft"]:
                run(cwd, "gh", "pr", "ready", pr, "--repo", REPO)
def agent(cwd, role, actor, model, work, output):
    hold, session = work["hold"]["id"], output.name
    shown = fid(cwd, "show", work["task"]["id"])
    context = {"work": work, "task": shown["issue"], "reviews": fid(cwd, "review", "list", work["task"]["id"])}
    if role == "judge": context["receipts"] = shown
    if work.get("review", {}).get("stream") == "separate": context["reviews"]["rounds"] = [r for r in context["reviews"]["rounds"] if r["stream"] == "separate"]
    (output / "context.json").write_text(json.dumps(context))
    spec = context["task"]["spec_md_path"]
    prompt = (HERE / "prompts" / (role + ".md")).read_text() + "\n" + json.dumps(context) + "\n" + ((Path(cwd) / spec).read_text() if spec else "")
    (output / "prompt.md").write_text(prompt)
    env = {**os.environ, **{"LOOP_" + k: str(v) for k, v in {"ACTOR": actor, "MODEL": model, "HOLD": hold, "SESSION": session, "CONTEXT": output / "context.json", "REPORT": output / "report.md", "REVIEW": output / "review.json", "ANSWERS": output / "answers.json"}.items()}}
    interval = fid(cwd, "heartbeat", "--hold", hold, "--actor", actor)["heartbeat_interval_seconds"]
    stop, errors = threading.Event(), []
    def renew():
        while not stop.wait(interval) and not errors:
            try:
                fid(cwd, "heartbeat", "--hold", hold, "--actor", actor)
            except Exception as e:
                errors.append(e)
    thread = threading.Thread(target=renew, daemon=True)
    command = os.getenv("LOOP_" + role.upper() + "_COMMAND", "codex exec --json -m {model} --add-dir {git_dir} --add-dir {common_dir} --add-dir {output_dir} -")
    command = [s.format(model=model.split("/", 1)[1], git_dir=run(cwd, "git", "rev-parse", "--absolute-git-dir"), common_dir=run(cwd, "git", "rev-parse", "--path-format=absolute", "--git-common-dir"), output_dir=str(output)) for s in shlex.split(command)]
    usage_format = os.getenv("LOOP_" + role.upper() + "_USAGE", "codex")
    if usage_format not in ("codex", "claude", "none"):
        raise ValueError(f"invalid {role} usage format: {usage_format}")
    version = ""
    try:
        probe = sp.run([command[0], "--version"], cwd=cwd, env=env, capture_output=True, text=True, timeout=5)
        if probe.returncode == 0:
            version = probe.stdout.strip()
    except (OSError, sp.TimeoutExpired):
        pass
    try:
        with (output / "prompt.md").open() as inp, (output / "agent.stdout").open("w") as stdout, (output / "agent.stderr").open("w") as stderr:
            child = sp.Popen(command, cwd=cwd, env=env, stdin=inp, stdout=stdout, stderr=stderr, start_new_session=True)
            thread.start()
            try:
                while child.poll() is None and not errors:
                    stop.wait(0.05)
                if errors or child.returncode != 0:
                    raise RuntimeError(f"agent failed: {errors or child.returncode}; see {output / 'agent.stderr'}")
            finally:
                stop.set()
                with contextlib.suppress(ProcessLookupError):
                    os.killpg(child.pid, signal.SIGKILL)
                child.wait()
                thread.join()
    finally:
        if usage_format != "none":
            try:
                version_flag = ["--tool-version", version] if version else []
                recorded = fid(cwd, "usage", "record", "--hold", hold, "--actor", actor,
                               "--tool", usage_format, "--tool-output", str(output / "agent.stdout"),
                               *version_flag, "--idempotency-key", session + "-usage")
                print("usage:", json.dumps(recorded), flush=True)
            except Exception as error:
                print(f"usage record failed: {error}", file=sys.stderr, flush=True)
    return context

def main():
    global ACTIVE
    build = run(Path.cwd(), "git", "rev-parse", "--show-toplevel")
    actors = [os.getenv("LOOP_" + r + "_ACTOR", "agent:sol-" + r.lower()) for r in ("BUILDER", "REVIEWER", "JUDGE")]
    models = [os.getenv("LOOP_" + r + "_MODEL", "openai/gpt-6.1-sol") for r in ("BUILDER", "REVIEWER", "JUDGE")]
    if len(set(actors)) != 3: raise ValueError("builder, reviewer and judge actors must differ")
    next_work(build, actors[0], models[0])
    done(build, actors[0])
    review = os.getenv("LOOP_REVIEW_WORKTREE", build + "-review")
    if not Path(review).exists(): run(build, "git", "worktree", "add", "--detach", review, "HEAD")
    while True:
        peek = next_work(review, actors[1], models[1])
        if peek["kind"] != "review":
            judge_peek = next_work(review, actors[2], models[2])
            if judge_peek["kind"] == "judge": peek = judge_peek
        is_review = peek["kind"] in ("review", "judge")
        role_index = 2 if peek["kind"] == "judge" else 1 if is_review else 0
        role = ("builder", "reviewer", "judge")[role_index]
        cwd, actor, model = review if is_review else build, actors[role_index], models[role_index]
        work = next_work(cwd, actor, model, True)
        if work["kind"] == "none":
            print("none:", json.dumps(work["waiting"]), flush=True)
            done(build, actors[0])
            return
        if (work["kind"] in ("review", "judge")) != is_review or is_review and work["kind"] != peek["kind"]:
            raise RuntimeError(f"work selection changed roles: {work!r}")
        hold = ["--hold", work["hold"]["id"], "--actor", actor]
        ACTIVE = (cwd, hold, is_review)
        output = Path(run(build, "git", "rev-parse", "--absolute-git-dir")) / "fid-loop/runs" / uuid.uuid4().hex
        output.mkdir(parents=True)
        print(f"{work['kind']}: {work['task']['id']} ({output})", flush=True)
        if is_review:
            if run(cwd, "git", "status", "--porcelain"): raise RuntimeError("review worktree has uncommitted changes")
            run(cwd, "git", "checkout", "--detach", work["review"]["handoff"]["head"])
        else:
            fid(cwd, "start", *hold, "--idempotency-key", output.name + "-start")
            if work["kind"] == "new":
                if run(cwd, "git", "status", "--porcelain"): raise RuntimeError("new work needs a clean worktree")
                task = fid(cwd, "show", work["task"]["id"])["issue"]
                run(cwd, "git", "switch", "-c", task_branch(task), run(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD"))
            branch = run(cwd, "git", "branch", "--show-current")
            if branch in ("", "main", run(cwd, "git", "symbolic-ref", "--short", "refs/remotes/origin/HEAD").removeprefix("origin/")): raise RuntimeError("build requires a task branch")
            if work.get("rework", {}).get("cause") == "default-conflict":
                run(cwd, "git", "merge", "--no-edit", work["rework"]["default_ref"], allowed=(0, 1))
        context = agent(cwd, role, actor, model, work, output)
        if is_review:
            if run(cwd, "git", "rev-parse", "HEAD") != work["review"]["handoff"]["head"] or run(cwd, "git", "status", "--porcelain"):
                raise RuntimeError("review agent changed the head or left uncommitted changes")
            review_input = json.loads((output / "review.json").read_text())
            if review_input.get("stream") != work["review"]["stream"]:
                raise RuntimeError("review output stream differs from requested work")
            fid(cwd, "review", "record", work["task"]["id"], "--input", str(output / "review.json"), "--actor", actor, "--idempotency-key", output.name)
        elif next_work(cwd, actor, model)["kind"] != "none":
            if run(cwd, "git", "status", "--porcelain"): raise RuntimeError(f"agent left uncommitted changes; see {output}")
            if work["kind"] == "rework" or any(f["open"] for r in context["reviews"]["rounds"] for f in r["findings"]):
                record_answers(cwd, actor, work["task"]["id"], output / "answers.json")
            fid(cwd, "check", *hold, "--idempotency-key", output.name + "-check")
            if run(cwd, "git", "branch", "--show-current") != branch: raise RuntimeError("agent changed the task branch")
            run(cwd, "git", "push", "-u", "origin", branch)
            prs = json.loads(run(cwd, "gh", "pr", "list", "--repo", REPO, "--head", branch, "--state", "open", "--json", "url"))
            pr = prs[0]["url"] if prs else run(cwd, "gh", "pr", "create", "--repo", REPO, "--draft", "--head", branch, "--title", pr_title(cwd, context["task"]), "--body-file", str(output / "report.md"))
            if prs: run(cwd, "gh", "pr", "edit", pr, "--repo", REPO, "--body-file", str(output / "report.md"))
            fid(cwd, "handoff", *hold, "--pr", pr, "--idempotency-key", output.name + "-handoff")
        ACTIVE = None
        print("backup:", json.dumps(fid(cwd, "backup")), flush=True)

try: main()
except (Exception, KeyboardInterrupt) as error:
    print(f"loop stopped: {str(error) or 'Ctrl-C'}", file=sys.stderr)
    if ACTIVE:
        cwd, hold, reviewing = ACTIVE
        for args in ([] if reviewing else [("park", *hold, "--interrupted", "--reason", str(error)[:500] or "Ctrl-C", "--idempotency-key", uuid.uuid4().hex)]) + [("backup",)]:
            try: fid(cwd, *args)
            except Exception as park_error: print(f"park/backup failed: {park_error}", file=sys.stderr)
    sys.exit(1)
