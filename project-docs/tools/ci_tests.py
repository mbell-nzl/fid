#!/usr/bin/env python3
"""Exhaustive ordinary Go test shards; Python standard library only."""
import argparse
import json
from pathlib import Path
import re
import subprocess

NAME = re.compile(r"(?:Test|Example|Fuzz)\w*")


def partitions(inventory, count):
    if count < 1 or not inventory:
        raise ValueError("positive shard count and package inventory required")
    names = set()
    for package, tests in inventory.items():
        if not package or len(tests) != len(set(tests)) or any(not NAME.fullmatch(t) for t in tests):
            raise ValueError("invalid or duplicate discovered test")
        names.update(tests)
    if len(names) < count:
        raise ValueError("empty shard selection")
    owners = {name: i % count for i, name in enumerate(sorted(names))}
    result = [set() for _ in range(count)]
    for package, tests in inventory.items():
        for name in tests:
            result[owners[name]].add((package, name))
    return result


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


def discover(log):
    return listed_tests(run_go(["-race", "-cover", "-list", ".", "-json", "./..."], log))


def run_shard(index, count, output):
    if not 0 <= index < count:
        raise ValueError("shard index outside declared count")
    output.mkdir(parents=True, exist_ok=False)
    inventory = discover(output / "discovery.jsonl")
    selected = partitions(inventory, count)[index]
    pattern = "^(" + "|".join(re.escape(t) for t in sorted({t for _, t in selected})) + ")$"
    print(f"Shard {index + 1}/{count}: {len(selected)} tests, examples and fuzz seed targets", flush=True)
    events = run_go(["-race", "-count=1", "-coverprofile=" + str(output / "coverage.out"), "-json", "-run", pattern, "./..."], output / "tests.jsonl")
    verify_execution(inventory, selected, events)
    ran, done, packages = execution_sets(events)
    source = subprocess.check_output(["git", "rev-parse", "HEAD"], text=True).strip()
    manifest = {"index": index, "count": count, "source": source, "inventory": inventory,
                "selected": sorted(selected), "executed": sorted(ran), "completed": sorted(done), "packages": sorted(packages)}
    (output / "manifest.json").write_text(json.dumps(manifest, indent=2) + "\n")
    print(f"PASS shard {index}: {len(done)} selected targets completed; all {len(packages)} packages accounted for")


def coverage_blocks(path):
    lines = path.read_text().splitlines()
    if not lines or lines[0] != "mode: atomic":
        raise ValueError("race coverage must use atomic mode")
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
    inventory, source = manifests[0]["inventory"], manifests[0]["source"]
    if not source or {m["index"] for m in manifests} != set(range(count)):
        raise ValueError("missing or duplicate shard indices/source")
    plan = partitions(inventory, count)
    combined = {}
    for path, manifest in zip(paths, manifests):
        if manifest["inventory"] != inventory or manifest["source"] != source or manifest["count"] != count:
            raise ValueError("inconsistent shard inventory/source/count")
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
    print(f"PASS coverage union: {count} shards, {sum(map(len, inventory.values()))} targets, {len(combined)} statement blocks")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("command", choices=("plan", "run", "merge"))
    parser.add_argument("--count", type=int, required=True)
    parser.add_argument("--index", type=int)
    parser.add_argument("--out", type=Path)
    parser.add_argument("--root", type=Path)
    args = parser.parse_args()
    if args.command == "run":
        if args.index is None or args.out is None:
            parser.error("run requires --index and --out (fresh directory)")
        run_shard(args.index, args.count, args.out.resolve())
    elif args.command == "merge":
        if args.root is None or args.out is None:
            parser.error("merge requires --root and --out")
        merge(args.root, args.count, args.out)
    else:
        import tempfile
        with tempfile.TemporaryDirectory(prefix="fid-ci-plan-") as scratch:
            inventory = discover(Path(scratch) / "discovery.jsonl")
        for i, part in enumerate(partitions(inventory, args.count)):
            print(f"Shard {i}: {len(part)} targets")
        print(f"PASS exhaustive plan: {len(inventory)} packages, {sum(map(len, inventory.values()))} targets")


if __name__ == "__main__":
    main()
