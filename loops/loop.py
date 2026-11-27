#!/usr/bin/env python3
import contextlib, json, os, re, shlex, signal, sys, threading, uuid
from pathlib import Path
import subprocess as sp

HERE = Path(__file__).resolve().parent
CONTRACT = json.loads((HERE.parent / "docs/specs/LOOP_CONTRACT.md").read_text().split("```json loop-contract\n")[1].split("\n```")[0])
ACTIVE, REPO = None, os.getenv("LOOP_GH_REPO", "mbell-nzl/fid")
OPERATIONS = dict(next="next.get", start="task.start", heartbeat="hold.heartbeat", park="hold.park", show="issue.show", status="project.status", backup="project.backup", done="task.done", check="check.execute", handoff="task.handoff")
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
        expected = "review." + args[1] if args[0] == "review" else "next.claim" if args[0] == "next" and "--claim" in args else OPERATIONS[args[0]]
        if reply["operation"] != expected: raise ValueError("unexpected operation")
        if reply["contract_version"] != "fid-cli.v1" or p.returncode != 0 or reply["ok"] is not True:
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
    return fid(cwd, "next", *scope, *(["--claim", "--actor", actor, "--model", model] if claim else []))
def done(cwd, actor):
    for outcome in fid(cwd, "done", "--actor", actor)["outcomes"]:
        if outcome["status"] not in ("completed", "waiting", "refused", "github-unavailable"): raise RuntimeError(f"unknown done outcome: {outcome!r}")
        print("done:", json.dumps(outcome), flush=True)
    for item in fid(cwd, "status")["needs_you"]:
        if item["kind"] in ("merge", "approval"):
            pr = fid(cwd, "show", item["subject"]["id"])["handoff"]["pr"]
            if json.loads(run(cwd, "gh", "pr", "view", pr, "--repo", REPO, "--json", "isDraft"))["isDraft"]:
                run(cwd, "gh", "pr", "ready", pr, "--repo", REPO)
def agent(cwd, role, actor, model, work, output):
    hold, session = work["hold"]["id"], output.name
    context = {"work": work, "task": fid(cwd, "show", work["task"]["id"])["issue"], "reviews": fid(cwd, "review", "list", work["task"]["id"])}
    if work.get("review", {}).get("stream") == "separate": context["reviews"]["rounds"] = [r for r in context["reviews"]["rounds"] if r["stream"] == "separate"]
    (output / "context.json").write_text(json.dumps(context))
    spec = context["task"]["spec_md_path"]
    prompt = (HERE / "prompts" / (role + ".md")).read_text() + "\n" + json.dumps(context) + "\n" + ((Path(cwd) / spec).read_text() if spec else "")
    (output / "prompt.md").write_text(prompt)
    env = {**os.environ, **{"LOOP_" + k: str(v) for k, v in {"ACTOR": actor, "MODEL": model, "HOLD": hold, "SESSION": session, "CONTEXT": output / "context.json", "REPORT": output / "report.md", "REVIEW": output / "review.json"}.items()}}
    interval = fid(cwd, "heartbeat", "--hold", hold, "--actor", actor)["heartbeat_interval_seconds"]
    stop, errors = threading.Event(), []
    def renew():
        while not stop.wait(interval) and not errors:
            try:
                fid(cwd, "heartbeat", "--hold", hold, "--actor", actor)
            except Exception as e:
                errors.append(e)
    thread = threading.Thread(target=renew, daemon=True)
    command = os.getenv("LOOP_" + role.upper() + "_COMMAND", "codex exec -m {model} --add-dir {git_dir} --add-dir {common_dir} --add-dir {output_dir} -")
    command = [s.format(model=model.split("/", 1)[1], git_dir=run(cwd, "git", "rev-parse", "--absolute-git-dir"), common_dir=run(cwd, "git", "rev-parse", "--path-format=absolute", "--git-common-dir"), output_dir=str(output)) for s in shlex.split(command)]
    with (output / "prompt.md").open() as inp, (output / "agent.log").open("w") as log:
        child = sp.Popen(command, cwd=cwd, env=env, stdin=inp, stdout=log, stderr=sp.STDOUT, start_new_session=True)
        thread.start()
        try:
            while child.poll() is None and not errors:
                stop.wait(0.05)
            if errors or child.returncode != 0: raise RuntimeError(f"agent failed: {errors or child.returncode}; see {output / 'agent.log'}")
        finally:
            stop.set()
            with contextlib.suppress(ProcessLookupError):
                os.killpg(child.pid, signal.SIGKILL)
            child.wait()
            thread.join()
def main():
    global ACTIVE
    build = run(Path.cwd(), "git", "rev-parse", "--show-toplevel")
    actors = [os.getenv("LOOP_" + r + "_ACTOR", "agent:sol-" + r.lower()) for r in ("BUILDER", "REVIEWER")]
    models = [os.getenv("LOOP_" + r + "_MODEL", "openai/gpt-6.1-sol") for r in ("BUILDER", "REVIEWER")]
    if actors[0] == actors[1]: raise ValueError("builder and reviewer actors must differ")
    next_work(build, actors[0], models[0])
    done(build, actors[0])
    review = os.getenv("LOOP_REVIEW_WORKTREE", build + "-review")
    if not Path(review).exists(): run(build, "git", "worktree", "add", "--detach", review, "HEAD")
    while True:
        is_review = next_work(review, actors[1], models[1])["kind"] == "review"
        cwd, actor, model = (review, actors[1], models[1]) if is_review else (build, actors[0], models[0])
        work = next_work(cwd, actor, model, True)
        if work["kind"] == "none":
            print("none:", json.dumps(work["waiting"]), flush=True)
            done(build, actors[0])
            return
        if (work["kind"] == "review") != is_review:
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
                name = Path(task["spec_md_path"]).stem if task["spec_md_path"] else task["id"]
                run(cwd, "git", "switch", "-c", "reset-" + name + "-" + re.sub(r"[^a-z0-9]+", "-", task["title"].lower()).strip("-"), run(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD"))
            branch = run(cwd, "git", "branch", "--show-current")
            if branch in ("", "main", run(cwd, "git", "symbolic-ref", "--short", "refs/remotes/origin/HEAD").removeprefix("origin/")): raise RuntimeError("build requires a task branch")
            if work.get("rework", {}).get("cause") == "default-conflict":
                run(cwd, "git", "merge", "--no-edit", work["rework"]["default_ref"], allowed=(0, 1))
        agent(cwd, "reviewer" if is_review else "builder", actor, model, work, output)
        if is_review:
            if run(cwd, "git", "rev-parse", "HEAD") != work["review"]["handoff"]["head"] or run(cwd, "git", "status", "--porcelain"):
                raise RuntimeError("review agent changed the head or left uncommitted changes")
            fid(cwd, "review", "record", work["task"]["id"], "--input", str(output / "review.json"), "--actor", actor, "--idempotency-key", output.name)
        elif next_work(cwd, actor, model)["kind"] != "none":
            if run(cwd, "git", "status", "--porcelain"): raise RuntimeError(f"agent left uncommitted changes; see {output}")
            fid(cwd, "check", *hold, "--idempotency-key", output.name + "-check")
            if run(cwd, "git", "branch", "--show-current") != branch: raise RuntimeError("agent changed the task branch")
            run(cwd, "git", "push", "-u", "origin", branch)
            prs = json.loads(run(cwd, "gh", "pr", "list", "--repo", REPO, "--head", branch, "--state", "open", "--json", "url"))
            pr = prs[0]["url"] if prs else run(cwd, "gh", "pr", "create", "--repo", REPO, "--draft", "--head", branch, "--title", fid(cwd, "show", work["task"]["id"])["issue"]["title"], "--body-file", str(output / "report.md"))
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
