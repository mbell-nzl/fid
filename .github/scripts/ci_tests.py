#!/usr/bin/env python3
"""Exhaustive ordinary Go test shards, CI reuse and the quality gate; Python standard library only."""
import argparse
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import urllib.request

NAME = re.compile(r"(?:Test|Example|Fuzz)\w*")
# DEC-259: race without checkptr in modernc SQLite; checkptr stays on for fid code.
# Decision G (task 0.7): only packages whose non-test code starts goroutines run
# under race; race_packages.go finds them by parsing the source.
RACE = ["-race", "-gcflags=modernc.org/...=-d=checkptr=0"]
RACE_HELPER = Path(__file__).with_name("race_packages.go")
# Race implies atomic coverage, so every package uses it and shard profiles merge.
COVER = ["-covermode=atomic"]
# Go's ten-minute default was too short for a shard's cmd/fid tests under race
# (1.1c1). 0.7 dropped race there and balanced shards by measured time; the
# rest of that entry stays on the later list in
# docs/reset/2026-10-31-reset-decisions.md.
TIMEOUT = ["-timeout", "20m"]
# Seconds per test of at least half a second, written by `durations` from
# shard logs (tests-*.jsonl). Shorter or new tests weigh the typical 0.1 s;
# every test weighs something, so each shard gets at least one.
DURATIONS = Path(__file__).with_name("durations.json")
UNMEASURED_SECONDS = 0.1
MIN_SECONDS = 0.05
TREE = re.compile(r"[0-9a-f]{40}|[0-9a-f]{64}")
WORKFLOW = ".github/workflows/ci.yml"


def load_durations(path=DURATIONS):
    return json.loads(path.read_text())


def partitions(inventory, count, durations=None):
    if count < 1 or not inventory:
        raise ValueError("positive shard count and package inventory required")
    weights = {}
    for package, tests in sorted(inventory.items()):
        if not package or len(tests) != len(set(tests)) or any(not NAME.fullmatch(t) for t in tests):
            raise ValueError("invalid or duplicate discovered test")
        measured = (durations or {}).get(package, {})
        for name in tests:
            weights[name] = weights.get(name, 0.0) + max(measured.get(name, UNMEASURED_SECONDS), MIN_SECONDS)
    if len(weights) < count:
        raise ValueError("empty shard selection")
    # -run selects a name in every package, so a name is the unit; heaviest first, each to the lightest shard.
    loads, owners = [0.0] * count, {}
    for name in sorted(weights, key=lambda n: (-weights[n], n)):
        owners[name] = min(range(count), key=lambda i: (loads[i], i))
        loads[owners[name]] += weights[name]
    result = [set() for _ in range(count)]
    for package, tests in inventory.items():
        for name in tests:
            result[owners[name]].add((package, name))
    return result


def collect_durations(events):
    durations = {}
    for event in events:
        name = event.get("Test")
        if name and "/" not in name and event["Action"] in ("pass", "fail", "skip") and "Elapsed" in event:
            durations.setdefault(event["Package"], {})[name] = round(event["Elapsed"], 1)
    return durations


def listed_tests(events):
    inventory = {}
    for event in events:
        package = event.get("Package")
        if not package:
            continue
        tests = inventory.setdefault(package, [])
        if event["Action"] == "output":
            tests.extend(name for name in event.get("Output", "").splitlines() if NAME.fullmatch(name))
    return {p: sorted(tests) for p, tests in sorted(inventory.items())}


def execution_sets(events):
    ran, done, packages = set(), set(), set()
    for event in events:
        package, name, action = event.get("Package"), event.get("Test"), event["Action"]
        if action == "fail":
            raise ValueError(f"test/package failure: {package} {name}")
        if name and "/" not in name:
            if action == "run":
                if (package, name) in ran:
                    raise ValueError("test ran more than once")
                ran.add((package, name))
            if action in ("pass", "skip"):
                done.add((package, name))
        elif not name and action in ("pass", "skip"):
            packages.add(package)
    return ran, done, packages


def verify_execution(inventory, selected, events):
    ran, done, packages = execution_sets(events)
    if ran != selected or done != selected or packages != set(inventory):
        raise ValueError(f"incomplete or extra execution: missing={selected - done}, extra={ran - selected}, packages={packages ^ set(inventory)}")


def run_go(arguments, log):
    events = []
    with log.open("w") as output:
        with subprocess.Popen(["go", "test", *arguments], stdout=subprocess.PIPE, text=True) as process:
            for line in process.stdout:
                output.write(line)
                event = json.loads(line)
                events.append(event)
                if event["Action"] == "output":
                    message = event.get("Output", "")
                    if not NAME.fullmatch(message.strip()) and not message.lstrip().startswith(("===", "--- PASS", "--- SKIP")):
                        print(message, end="", flush=True)
            code = process.wait()
    if code:
        raise subprocess.CalledProcessError(code, process.args)
    return events


def discover(groups, output):
    """Lists each group's tests with the flags it runs with, so tests behind //go:build race are found."""
    inventory = {}
    for flags, packages in groups:
        inventory.update(listed_tests(run_go([*flags, *COVER, "-list", ".", "-json", *packages], output / f"discovery-{label(flags)}.jsonl")))
    if sorted(inventory) != sorted(p for _, packages in groups for p in packages):
        raise ValueError("discovery did not list exactly the packages it was given")
    return inventory


def list_packages():
    return subprocess.check_output(["go", "list", "./..."], text=True).split()


def race_packages(cwd=None):
    return sorted(set(subprocess.check_output(["go", "run", str(RACE_HELPER), "./..."], cwd=cwd, text=True).split()))


def package_groups(packages, race):
    if not set(race) <= set(packages):
        raise ValueError(f"race packages outside the package list: {sorted(set(race) - set(packages))}")
    plain = [p for p in packages if p not in race]
    return [(flags, group) for flags, group in ((RACE, sorted(race)), ([], plain)) if group]


def label(flags):
    return "race" if flags else "plain"


def run_shard(index, count, output):
    if not 0 <= index < count:
        raise ValueError("shard index outside declared count")
    output.mkdir(parents=True, exist_ok=False)
    race = race_packages()
    groups = package_groups(list_packages(), race)
    inventory = discover(groups, output)
    selected = partitions(inventory, count, load_durations())[index]
    pattern = "^(" + "|".join(re.escape(t) for t in sorted({t for _, t in selected})) + ")$"
    print(f"Shard {index + 1}/{count}: {len(selected)} tests, examples and fuzz seed targets", flush=True)
    print(f"Race detector on {', '.join(race) or 'no packages'}: their non-test code starts goroutines", flush=True)
    events, profile = [], ["mode: atomic"]
    for flags, packages in groups:
        coverage = output / f"coverage-{label(flags)}.out"
        events += run_go([*flags, *COVER, *TIMEOUT, "-count=1", "-coverprofile=" + str(coverage), "-json", "-run", pattern, *packages], output / f"tests-{label(flags)}.jsonl")
        lines = coverage.read_text().splitlines()
        if not lines or lines[0] != "mode: atomic":
            raise ValueError("coverage must use atomic mode")
        profile += lines[1:]
    (output / "coverage.out").write_text("\n".join(profile) + "\n")
    coverage_blocks(output / "coverage.out")
    verify_execution(inventory, selected, events)
    ran, done, packages = execution_sets(events)
    source = subprocess.check_output(["git", "rev-parse", "HEAD"], text=True).strip()
    manifest = {"index": index, "count": count, "source": source, "inventory": inventory, "race": race,
                "selected": sorted(selected), "executed": sorted(ran), "completed": sorted(done), "packages": sorted(packages)}
    (output / "manifest.json").write_text(json.dumps(manifest, indent=2) + "\n")
    print(f"PASS shard {index}: {len(done)} selected targets completed; all {len(packages)} packages accounted for")


def coverage_blocks(path):
    lines = path.read_text().splitlines()
    if not lines or lines[0] != "mode: atomic":
        raise ValueError("coverage must use atomic mode")
    blocks = {}
    for line in lines[1:]:
        coordinate, statements, count = line.split()
        statements, count = int(statements), int(count)
        if not re.fullmatch(r".+:\d+\.\d+,\d+\.\d+", coordinate) or statements < 0 or count < 0 or coordinate in blocks:
            raise ValueError("invalid or duplicate coverage block")
        blocks[coordinate] = (statements, count)
    if not blocks:
        raise ValueError("empty coverage profile")
    return blocks


def merge(root, count, output):
    paths = sorted(root.glob("shard-*/manifest.json"))
    if count < 1 or len(paths) != count:
        raise ValueError("missing or extra shard manifests")
    manifests = [json.loads(path.read_text()) for path in paths]
    inventory, source, race = manifests[0]["inventory"], manifests[0]["source"], manifests[0].get("race")
    if not source or {m["index"] for m in manifests} != set(range(count)):
        raise ValueError("missing or duplicate shard indices/source")
    if not isinstance(race, list) or not set(race) <= set(inventory):
        raise ValueError("race packages missing or outside the inventory")
    plan = partitions(inventory, count, load_durations())
    combined = {}
    for path, manifest in zip(paths, manifests):
        if manifest["inventory"] != inventory or manifest["source"] != source or manifest["count"] != count or manifest.get("race") != race:
            raise ValueError("inconsistent shard inventory/source/count/race packages")
        expected = sorted(map(list, plan[manifest["index"]]))
        if any(manifest[key] != expected for key in ("selected", "executed", "completed")) or manifest["packages"] != sorted(inventory):
            raise ValueError("shard did not account for every selected test/package")
        blocks = coverage_blocks(path.parent / "coverage.out")
        if combined and blocks.keys() != combined.keys():
            raise ValueError("shard coverage omits or adds statement blocks")
        for coordinate, (statements, hits) in blocks.items():
            previous_statements, previous_hits = combined.get(coordinate, (statements, 0))
            if previous_statements != statements:
                raise ValueError("inconsistent coverage statement counts")
            combined[coordinate] = statements, previous_hits + hits
    output.write_text("mode: atomic\n" + "".join(f"{coordinate} {statements} {hits}\n" for coordinate, (statements, hits) in sorted(combined.items())))
    print(f"PASS coverage union: {count} shards, {sum(map(len, inventory.values()))} targets, {len(combined)} statement blocks; race detector on {', '.join(race) or 'no packages'}")


def find_reuse(tree, get):
    """The successful pull_request run of this workflow that recorded tree, or None."""
    if not TREE.fullmatch(tree):
        raise ValueError(f"not a tree hash: {tree!r}")
    name = f"ci-tree-{tree}"
    for artifact in get(f"actions/artifacts?name={name}&per_page=100")["artifacts"]:
        if artifact["name"] != name or artifact["expired"]:
            continue
        run = get(f"actions/runs/{int(artifact['workflow_run']['id'])}")
        if (run["event"], run["status"], run["conclusion"], run["path"]) == ("pull_request", "completed", "success", WORKFLOW):
            return run
    return None


def decide_reuse(tree, get):
    if not TREE.fullmatch(tree):
        raise ValueError(f"not a tree hash: {tree!r}")
    try:
        run = find_reuse(tree, get)
    except (OSError, ValueError, KeyError, TypeError) as error:
        return {"run": "", "message": f"Full run: looking up CI runs for tree {tree} failed ({error})."}
    if run is None:
        return {"run": "", "message": f"Full run: no successful pull_request CI run recorded tree {tree}."}
    return {"run": run["html_url"], "message": f"Reusing {run['html_url']}: that pull_request run passed on tree {tree}, the tree of this commit."}


def github_get(path):
    request = urllib.request.Request(f"{os.environ['GITHUB_API_URL']}/repos/{os.environ['GITHUB_REPOSITORY']}/{path}", headers={
        "Authorization": f"Bearer {os.environ['GITHUB_TOKEN']}", "Accept": "application/vnd.github+json", "X-GitHub-Api-Version": "2022-11-28"})
    with urllib.request.urlopen(request, timeout=30) as response:
        return json.load(response)


def gate(event, reuse, run, checks, tests):
    """quality passes on a full run whose jobs all succeeded, or on a push that verified a reuse."""
    if checks == tests == "success":
        return "full"
    if event == "push" and reuse == "success" and run and checks == tests == "skipped":
        return "reuse"
    raise ValueError(f"CI did not pass: event={event} reuse={reuse} reused run={run or 'none'} checks={checks} tests={tests}")


def write_github(variable, text):
    if os.environ.get(variable):
        with open(os.environ[variable], "a") as output:
            output.write(text)


def local_test():
    race = race_packages()
    packages = list_packages()
    print(f"Race detector on {', '.join(race) or 'no packages'}: their non-test code starts goroutines", flush=True)
    codes = [subprocess.run(["go", "test", *flags, *group]).returncode for flags, group in package_groups(packages, race)]
    return 1 if any(codes) else 0


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest="command", required=True)
    plan = commands.add_parser("plan", help="print the shard plan with estimated seconds")
    plan.add_argument("--count", type=int, required=True)
    run = commands.add_parser("run", help="run one shard into a fresh directory")
    run.add_argument("--count", type=int, required=True)
    run.add_argument("--index", type=int, required=True)
    run.add_argument("--out", type=Path, required=True)
    combine = commands.add_parser("merge", help="verify every shard and combine coverage")
    combine.add_argument("--count", type=int, required=True)
    combine.add_argument("--root", type=Path, required=True)
    combine.add_argument("--out", type=Path, required=True)
    durations = commands.add_parser("durations", help="write measured test durations from go test -json logs")
    durations.add_argument("--out", type=Path, default=DURATIONS)
    durations.add_argument("--min", type=float, default=0.5, help="omit tests shorter than this many seconds")
    durations.add_argument("logs", type=Path, nargs="+")
    reuse = commands.add_parser("reuse", help="find a passed pull_request run that tested this tree")
    reuse.add_argument("--tree", required=True)
    quality = commands.add_parser("gate", help="pass only a full successful run or a verified reuse")
    for name in ("event", "reuse", "reused-run", "checks", "tests", "tree"):
        quality.add_argument(f"--{name}", default="")
    commands.add_parser("test", help="local go test, with race only where goroutines start")
    args = parser.parse_args()
    if args.command == "run":
        run_shard(args.index, args.count, args.out.resolve())
    elif args.command == "merge":
        merge(args.root, args.count, args.out)
    elif args.command == "durations":
        events = [json.loads(line) for log in args.logs for line in log.read_text().splitlines() if line.startswith("{")]
        measured = {p: {t: s for t, s in tests.items() if s >= args.min} for p, tests in collect_durations(events).items()}
        args.out.write_text(json.dumps({p: tests for p, tests in measured.items() if tests}, indent=1, sort_keys=True) + "\n")
    elif args.command == "reuse":
        decision = decide_reuse(args.tree, github_get)
        print(decision["message"])
        write_github("GITHUB_OUTPUT", f"run={decision['run']}\ntree={args.tree}\n")
        write_github("GITHUB_STEP_SUMMARY", decision["message"] + "\n")
    elif args.command == "gate":
        try:
            mode = gate(args.event, args.reuse, args.reused_run, args.checks, args.tests)
        except ValueError as error:
            sys.exit(f"FAIL: {error}")
        message = (f"PASS: reused {args.reused_run}, which passed on tree {args.tree}, the tree of this commit; no tests ran here."
                   if mode == "reuse" else "PASS: checks and every test shard succeeded.")
        print(message)
        write_github("GITHUB_OUTPUT", f"mode={mode}\n")
        write_github("GITHUB_STEP_SUMMARY", message + "\n")
    elif args.command == "test":
        sys.exit(local_test())
    else:
        import tempfile
        with tempfile.TemporaryDirectory(prefix="fid-ci-plan-") as scratch:
            inventory = discover(package_groups(list_packages(), race_packages()), Path(scratch))
        measured = load_durations()
        for i, part in enumerate(partitions(inventory, args.count, measured)):
            seconds = sum(max(measured.get(p, {}).get(t, UNMEASURED_SECONDS), MIN_SECONDS) for p, t in part)
            print(f"Shard {i}: {len(part)} targets, about {seconds:.0f} s measured")
        print(f"PASS exhaustive plan: {len(inventory)} packages, {sum(map(len, inventory.values()))} targets")


if __name__ == "__main__":
    main()
