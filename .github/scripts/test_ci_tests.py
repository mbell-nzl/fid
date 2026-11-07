"""Behavioral checks for complete CI selection and coverage reconciliation."""
import importlib.util
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest
from unittest import mock

spec = importlib.util.spec_from_file_location("ci_tests", Path(__file__).with_name("ci_tests.py"))
ci = importlib.util.module_from_spec(spec)
spec.loader.exec_module(ci)
SCRIPT = str(Path(__file__).with_name("ci_tests.py").resolve())


class ShardTests(unittest.TestCase):
    def setUp(self):
        self.inventory = {"example/a": [f"TestCase{i}" for i in range(20)] + ["ExampleOutput", "FuzzSeeds"], "example/b": ["TestCase0"], "example/empty": []}

    def test_every_test_example_and_fuzz_seed_is_selected_exactly_once(self):
        parts = ci.partitions(self.inventory, 4)
        expected = {(p, t) for p, tests in self.inventory.items() for t in tests}
        actual = [item for part in parts for item in part]
        self.assertEqual(set(actual), expected)
        self.assertEqual(len(actual), len(set(actual)))
        self.assertTrue(all(parts))
        for i, part in enumerate(parts):
            if ("example/a", "TestCase0") in part:
                self.assertIn(("example/b", "TestCase0"), part)

    def test_new_tests_are_automatically_included(self):
        self.inventory["example/a"].append("TestBrandNew")
        self.assertEqual(sum(("example/a", "TestBrandNew") in x for x in ci.partitions(self.inventory, 4)), 1)

    def test_measured_durations_balance_shards(self):
        tests = ["TestA"] + [f"TestB{i}" for i in range(10)]
        inventory = {"p": tests}
        durations = {"p": {"TestA": 10.0, **{f"TestB{i}": 1.0 for i in range(10)}}}
        parts = ci.partitions(inventory, 2, durations)
        self.assertIn({("p", "TestA")}, parts)
        self.assertEqual(sorted(map(len, parts)), [1, 10])

    def test_a_name_shared_by_packages_carries_their_combined_weight(self):
        inventory = {"p": ["TestShared", "TestX"], "q": ["TestShared", "TestY", "TestZ"]}
        durations = {"p": {"TestShared": 3.0, "TestX": 2.0}, "q": {"TestShared": 3.0, "TestY": 2.0, "TestZ": 2.0}}
        parts = ci.partitions(inventory, 2, durations)
        self.assertIn({("p", "TestShared"), ("q", "TestShared")}, parts)

    def test_unmeasured_and_zero_duration_tests_still_fill_every_shard_exactly_once(self):
        inventory = {"p": ["TestA", "TestB", "TestC", "TestD"]}
        for durations in [None, {}, {"p": {"TestA": 0.0, "TestB": 0.0, "TestC": 0.0, "TestD": 0.0}}]:
            with self.subTest(durations=durations):
                parts = ci.partitions(inventory, 4, durations)
                self.assertEqual(sorted(map(len, parts)), [1, 1, 1, 1])

    def test_durations_are_collected_from_top_level_test_events(self):
        events = [{"Action": "run", "Package": "p", "Test": "TestA"}, {"Action": "pass", "Package": "p", "Test": "TestA/sub", "Elapsed": 9.0},
                  {"Action": "pass", "Package": "p", "Test": "TestA", "Elapsed": 1.26}, {"Action": "skip", "Package": "p", "Test": "TestB", "Elapsed": 0},
                  {"Action": "pass", "Package": "p", "Elapsed": 12.0}, {"Action": "output", "Package": "p", "Test": "TestA", "Output": "x"}]
        self.assertEqual(ci.collect_durations(events), {"p": {"TestA": 1.3, "TestB": 0.0}})

    def test_checked_in_durations_are_well_formed(self):
        durations = ci.load_durations()
        self.assertTrue(durations)
        for package, tests in durations.items():
            self.assertTrue(package.startswith("github.com/mbell-nzl/fid/"), package)
            for name, seconds in tests.items():
                self.assertRegex(name, ci.NAME)
                self.assertIsInstance(seconds, (int, float))
                self.assertGreaterEqual(seconds, 0)

    def test_race_packages_run_with_the_race_detector_and_the_rest_without(self):
        groups = ci.package_groups(["p/a", "p/race", "p/b"], ["p/race"])
        self.assertEqual(groups, [(ci.RACE, ["p/race"]), ([], ["p/a", "p/b"])])
        self.assertIn("-gcflags=modernc.org/...=-d=checkptr=0", ci.RACE)
        self.assertEqual(ci.package_groups(["p/a"], []), [([], ["p/a"])])
        with self.assertRaises(ValueError):
            ci.package_groups(["p/a"], ["p/unknown"])

    def test_invalid_or_empty_inventory_fails(self):
        for inventory, count in [({}, 4), ({"a": []}, 4), (self.inventory, 0), ({"a": ["TestOne", "TestOne"]}, 1), ({"a": ["BenchmarkOnly"]}, 1)]:
            with self.subTest(inventory=inventory, count=count), self.assertRaises(ValueError):
                ci.partitions(inventory, count)

    def test_listing_includes_packages_without_tests_but_not_benchmarks(self):
        events = [{"Action": "start", "Package": "p"}, {"Action": "output", "Package": "p", "Output": "TestOne\nExampleOne\nFuzzOne\nBenchmarkOne\n"}, {"Action": "start", "Package": "empty"}]
        self.assertEqual(ci.listed_tests(events), {"p": ["ExampleOne", "FuzzOne", "TestOne"], "empty": []})

    def test_execution_requires_real_run_completion_and_all_packages(self):
        inv = {"p": ["TestOne"], "empty": []}
        selected = {("p", "TestOne")}
        events = [{"Action": "run", "Package": "p", "Test": "TestOne"}, {"Action": "pass", "Package": "p", "Test": "TestOne"}, {"Action": "pass", "Package": "p"}, {"Action": "skip", "Package": "empty"}]
        ci.verify_execution(inv, selected, events)
        for bad in [events[1:], events[:-1], events + [{"Action": "run", "Package": "p", "Test": "TestOther"}], [dict(e, Action="fail") if e.get("Test") else e for e in events]]:
            with self.subTest(events=bad), self.assertRaises(ValueError):
                ci.verify_execution(inv, selected, bad)

    def write_shards(self, root, count=4, flat=False, names=None):
        """Shard evidence as download-artifact leaves it: one artifact flat in root, several in a directory each."""
        parts = ci.partitions(self.inventory, count)
        for i, part in enumerate(parts):
            path = root if flat else root / (names[i] if names else f"shard-{i}")
            path.mkdir(parents=True, exist_ok=True)
            manifest = {"index": i, "count": count, "source": "same-source", "inventory": self.inventory, "race": ["example/b"], "selected": sorted(part), "executed": sorted(part), "completed": sorted(part), "packages": sorted(self.inventory)}
            (path / "manifest.json").write_text(json.dumps(manifest))
            (path / "coverage.out").write_text(f"mode: atomic\nexample/a/a.go:1.1,3.2 1 {int(i == 0)}\nexample/a/a.go:4.1,6.2 1 {int(i == 1)}\nexample/empty/e.go:1.1,3.2 1 0\n")
            for log in ("discovery-plain.jsonl", "discovery-race.jsonl", "tests-plain.jsonl", "tests-race.jsonl", "coverage-plain.out", "coverage-race.out"):
                (path / log).write_text("")

    def test_one_artifact_downloaded_flat_into_the_root_merges(self):
        # With one matching artifact, actions/download-artifact extracts it straight into path (run 42320491923).
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory) / "ci-results"
            self.write_shards(root, count=1, flat=True)
            ci.merge(root, 1, Path(directory) / "combined.out")
            self.assertEqual((Path(directory) / "combined.out").read_text(), "mode: atomic\nexample/a/a.go:1.1,3.2 1 1\nexample/a/a.go:4.1,6.2 1 0\nexample/empty/e.go:1.1,3.2 1 0\n")

    def test_shards_are_identified_by_their_manifests_not_their_directory_names(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory) / "ci-results"
            self.write_shards(root, names=["zeta", "alpha", "shard-0", "mid"])
            ci.merge(root, 4, Path(directory) / "combined.out")
            self.assertEqual((Path(directory) / "combined.out").read_text(), "mode: atomic\nexample/a/a.go:1.1,3.2 1 1\nexample/a/a.go:4.1,6.2 1 1\nexample/empty/e.go:1.1,3.2 1 0\n")

    def test_duplicate_stray_or_mismatched_manifests_fail_in_either_layout(self):
        def same_index_flat_and_directory(root):
            self.write_shards(root, count=2, names=["a", "b"])
            (root / "a" / "manifest.json").rename(root / "manifest.json")
            (root / "a" / "coverage.out").rename(root / "coverage.out")
            changed = json.loads((root / "b" / "manifest.json").read_text())
            changed["index"] = 0
            (root / "b" / "manifest.json").write_text(json.dumps(changed))

        def same_index_two_directories(root):
            self.write_shards(root, count=2)
            changed = json.loads((root / "shard-1" / "manifest.json").read_text())
            changed["index"] = 0
            (root / "shard-1" / "manifest.json").write_text(json.dumps(changed))

        def stray_beside_flat(root):
            self.write_shards(root, count=1, flat=True)
            (root / "old").mkdir()
            for name in ("manifest.json", "coverage.out"):
                (root / "old" / name).write_text((root / name).read_text())

        def stray_nested_in_a_shard(root):
            self.write_shards(root, count=2)
            (root / "shard-0" / "nested").mkdir()
            for name in ("manifest.json", "coverage.out"):
                (root / "shard-0" / "nested" / name).write_text((root / "shard-0" / name).read_text())

        def flat_with_another_count(root):
            self.write_shards(root, count=1, flat=True)
            changed = json.loads((root / "manifest.json").read_text())
            changed["count"] = 2
            (root / "manifest.json").write_text(json.dumps(changed))

        def flat_with_another_index(root):
            self.write_shards(root, count=1, flat=True)
            changed = json.loads((root / "manifest.json").read_text())
            changed["index"] = 1
            (root / "manifest.json").write_text(json.dumps(changed))

        cases = [(same_index_flat_and_directory, 2), (same_index_two_directories, 2), (stray_beside_flat, 1), (stray_nested_in_a_shard, 2),
                 (flat_with_another_count, 1), (flat_with_another_index, 1), (lambda root: root.mkdir(parents=True), 1)]
        for layout, count in cases:
            with self.subTest(layout=getattr(layout, "__name__", "empty")), tempfile.TemporaryDirectory() as directory:
                root = Path(directory) / "ci-results"
                layout(root)
                with self.assertRaises(ValueError):
                    ci.merge(root, count, Path(directory) / "combined.out")
                self.assertFalse((Path(directory) / "combined.out").exists())

    def test_coverage_union_retains_all_blocks_and_sums_counts(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            self.write_shards(root)
            ci.merge(root, 4, root / "combined.out")
            self.assertEqual((root / "combined.out").read_text(), "mode: atomic\nexample/a/a.go:1.1,3.2 1 1\nexample/a/a.go:4.1,6.2 1 1\nexample/empty/e.go:1.1,3.2 1 0\n")

    def test_missing_duplicate_or_inconsistent_shards_fail_closed(self):
        for mutation in ["missing", "duplicate", "source", "inventory", "selected", "executed", "completed", "packages", "count", "race", "race-unknown", "empty-profile", "missing-block", "wrong-mode", "statements", "negative-count"]:
            with self.subTest(mutation=mutation), tempfile.TemporaryDirectory() as directory:
                root = Path(directory)
                self.write_shards(root)
                path = root / "shard-1"
                mpath = path / "manifest.json"
                m = json.loads(mpath.read_text())
                if mutation == "missing":
                    mpath.unlink()
                elif mutation == "duplicate":
                    m["index"] = 0
                elif mutation == "source":
                    m["source"] = "other-source"
                elif mutation == "inventory":
                    m["inventory"]["example/a"].append("TestUnaccounted")
                elif mutation in ("selected", "executed", "completed", "packages"):
                    m[mutation] = []
                elif mutation == "count":
                    m["count"] = 3
                elif mutation == "race":
                    m["race"] = []
                elif mutation == "race-unknown":
                    for shard in range(4):
                        other = root / f"shard-{shard}" / "manifest.json"
                        changed = json.loads(other.read_text())
                        changed["race"] = ["example/unknown"]
                        other.write_text(json.dumps(changed))
                    m["race"] = ["example/unknown"]
                else:
                    profile = path / "coverage.out"
                    text = profile.read_text()
                    text = {"empty-profile": "mode: atomic\n", "missing-block": "\n".join(text.splitlines()[:-1]) + "\n", "wrong-mode": text.replace("atomic", "set"), "statements": text.replace("3.2 1", "3.2 2"), "negative-count": text.replace("1 0", "1 -1")}[mutation]
                    profile.write_text(text)
                if mutation != "missing":
                    mpath.write_text(json.dumps(m))
                with self.assertRaises(ValueError):
                    ci.merge(root, 4, root / "combined.out")


class ShardTimeoutTests(unittest.TestCase):
    def test_shard_run_gives_go_test_more_than_the_default_ten_minutes(self):
        calls = []

        def fake_run_go(arguments, log):
            calls.append(arguments)
            Path(next(a for a in arguments if a.startswith("-coverprofile="))[len("-coverprofile="):]).write_text("mode: atomic\nexample/a/a.go:1.1,2.2 1 1\n")
            return [{"Package": "example/a", "Test": "TestA", "Action": "run"},
                    {"Package": "example/a", "Test": "TestA", "Action": "pass"},
                    {"Package": "example/a", "Action": "pass"}]

        with tempfile.TemporaryDirectory(prefix="fid-ci-timeout-") as directory, \
                mock.patch.object(ci, "discover", return_value={"example/a": ["TestA"]}), \
                mock.patch.object(ci, "classify", return_value=(["example/a"], [])), \
                mock.patch.object(ci, "run_go", side_effect=fake_run_go), \
                mock.patch.object(ci.subprocess, "check_output", return_value="source\n"):
            ci.run_shard(0, 1, Path(directory) / "shard-0")
        (arguments,) = calls
        self.assertIn("-timeout", arguments)
        self.assertEqual(arguments[arguments.index("-timeout") + 1], "20m")

    def test_shard_runs_race_packages_under_race_and_the_rest_without_in_one_coverage_mode(self):
        inventory = {"example/a": ["TestA"], "example/r": ["TestR"]}
        calls = []

        def fake_run_go(arguments, log):
            calls.append(arguments)
            packages = [a for a in arguments if a in inventory]
            profile = Path(next(a for a in arguments if a.startswith("-coverprofile="))[len("-coverprofile="):])
            profile.write_text("mode: atomic\n" + "".join(f"{p}/x.go:1.1,2.2 1 1\n" for p in packages))
            return [event for p in packages for event in ({"Package": p, "Test": inventory[p][0], "Action": "run"}, {"Package": p, "Test": inventory[p][0], "Action": "pass"}, {"Package": p, "Action": "pass"})]

        with tempfile.TemporaryDirectory(prefix="fid-ci-groups-") as directory, \
                mock.patch.object(ci, "discover", return_value=inventory), \
                mock.patch.object(ci, "classify", return_value=(sorted(inventory), ["example/r"])), \
                mock.patch.object(ci, "run_go", side_effect=fake_run_go), \
                mock.patch.object(ci.subprocess, "check_output", return_value="source\n"):
            output = Path(directory) / "shard-0"
            ci.run_shard(0, 1, output)
            manifest = json.loads((output / "manifest.json").read_text())
            blocks = ci.coverage_blocks(output / "coverage.out")
        race, plain = calls
        self.assertEqual([a for a in race if a in inventory], ["example/r"])
        self.assertEqual([a for a in plain if a in inventory], ["example/a"])
        self.assertIn("-race", race)
        self.assertIn("-gcflags=modernc.org/...=-d=checkptr=0", race)
        self.assertNotIn("-race", plain)
        for arguments in calls:
            self.assertIn("-covermode=atomic", arguments)
            self.assertEqual(arguments[arguments.index("-timeout") + 1], "20m")
        self.assertEqual(manifest["race"], ["example/r"])
        self.assertEqual(set(blocks), {"example/a/x.go:1.1,2.2", "example/r/x.go:1.1,2.2"})


class ReuseTests(unittest.TestCase):
    TREE = "0123456789abcdef0123456789abcdef01234567"

    def run_record(self, run_id, **changes):
        return {"id": run_id, "event": "pull_request", "status": "completed", "conclusion": "success", "path": ".github/workflows/ci.yml",
                "html_url": f"https://github.com/o/r/actions/runs/{run_id}", **changes}

    def artifact(self, run_id, **changes):
        return {"id": run_id * 10, "name": f"ci-tree-{self.TREE}", "expired": False, "workflow_run": {"id": run_id}, **changes}

    def api(self, artifacts, runs):
        requests = []

        def get(path):
            requests.append(path)
            if path.startswith("actions/artifacts?"):
                return {"total_count": len(artifacts), "artifacts": artifacts}
            return runs[int(path.rsplit("/", 1)[1])]
        return get, requests

    def test_a_successful_pull_request_run_that_tested_the_tree_is_reused(self):
        get, requests = self.api([self.artifact(7)], {7: self.run_record(7)})
        self.assertEqual(ci.find_reuse(self.TREE, get)["html_url"], "https://github.com/o/r/actions/runs/7")
        self.assertIn(f"name=ci-tree-{self.TREE}", requests[0])

    def test_only_a_completed_successful_ci_pull_request_run_with_a_live_marker_counts(self):
        cases = {
            "no marker": ([], {}),
            "expired marker": ([self.artifact(7, expired=True)], {7: self.run_record(7)}),
            "other tree": ([self.artifact(7, name="ci-tree-" + "f" * 40)], {7: self.run_record(7)}),
            "failed run": ([self.artifact(7)], {7: self.run_record(7, conclusion="failure")}),
            "cancelled run": ([self.artifact(7)], {7: self.run_record(7, conclusion="cancelled")}),
            "running": ([self.artifact(7)], {7: self.run_record(7, status="in_progress", conclusion=None)}),
            "push run": ([self.artifact(7)], {7: self.run_record(7, event="push")}),
            "other workflow": ([self.artifact(7)], {7: self.run_record(7, path=".github/workflows/other.yml")}),
        }
        for label, (artifacts, runs) in cases.items():
            with self.subTest(label):
                get, _ = self.api(artifacts, runs)
                self.assertIsNone(ci.find_reuse(self.TREE, get))

    def test_any_qualifying_run_among_several_markers_is_reused(self):
        get, _ = self.api([self.artifact(7), self.artifact(8)], {7: self.run_record(7, conclusion="failure"), 8: self.run_record(8)})
        self.assertEqual(ci.find_reuse(self.TREE, get)["id"], 8)

    def test_a_malformed_tree_is_refused_before_any_lookup(self):
        for tree in ["", "HEAD", "0123", self.TREE.upper(), self.TREE + "\n"]:
            with self.subTest(tree=tree), self.assertRaises(ValueError):
                ci.find_reuse(tree, lambda path: self.fail("looked up a malformed tree"))

    def test_a_failed_lookup_means_a_full_run(self):
        def unreachable(path):
            raise OSError("network down")
        decision = ci.decide_reuse(self.TREE, unreachable)
        self.assertEqual(decision["run"], "")
        self.assertIn("network down", decision["message"])
        get, _ = self.api([self.artifact(7)], {7: self.run_record(7)})
        decision = ci.decide_reuse(self.TREE, get)
        self.assertEqual(decision["run"], "https://github.com/o/r/actions/runs/7")
        self.assertIn(self.TREE, decision["message"])


RUN_URL = "https://github.com/o/r/actions/runs/7"
TREE = "0123456789abcdef0123456789abcdef01234567"
# (event, reuse result, reused run, checks result, tests result) that quality must refuse.
REFUSED = [
    ("pull_request", "skipped", "", "failure", "success"),
    ("pull_request", "skipped", "", "success", "cancelled"),
    ("pull_request", "skipped", "", "skipped", "skipped"),
    ("pull_request", "skipped", RUN_URL, "skipped", "skipped"),
    ("push", "success", "", "skipped", "skipped"),
    ("push", "failure", "", "skipped", "skipped"),
    ("push", "failure", RUN_URL, "skipped", "skipped"),
    ("push", "success", RUN_URL, "success", "skipped"),
    ("push", "success", RUN_URL, "skipped", "failure"),
    ("push", "cancelled", "", "skipped", "skipped"),
]


class GateTests(unittest.TestCase):
    URL = RUN_URL

    def test_a_full_run_passes_only_when_checks_and_tests_succeeded(self):
        self.assertEqual(ci.gate("pull_request", "skipped", "", "success", "success"), "full")
        self.assertEqual(ci.gate("push", "success", "", "success", "success"), "full")

    def test_a_push_passes_without_tests_only_on_a_verified_reuse(self):
        self.assertEqual(ci.gate("push", "success", self.URL, "skipped", "skipped"), "reuse")

    def test_anything_else_fails(self):
        for case in REFUSED:
            with self.subTest(case=case), self.assertRaises(ValueError):
                ci.gate(*case)


class GateCommandTests(unittest.TestCase):
    """The gate command exactly as the quality job runs it."""

    def gate_command(self, event, reuse, run, checks, tests):
        with tempfile.TemporaryDirectory(prefix="fid-ci-gate-") as directory:
            output, summary = Path(directory) / "output", Path(directory) / "summary"
            env = {**os.environ, "GITHUB_OUTPUT": str(output), "GITHUB_STEP_SUMMARY": str(summary)}
            result = subprocess.run([sys.executable, SCRIPT, "gate", "--event", event, "--reuse", reuse, "--reused-run", run, "--tree", TREE if run else "",
                                     "--checks", checks, "--tests", tests], env=env, capture_output=True, text=True)
            return result, output.read_text() if output.exists() else "", summary.read_text() if summary.exists() else ""

    def test_a_full_run_exits_zero_and_outputs_mode_full(self):
        result, output, _ = self.gate_command("pull_request", "skipped", "", "success", "success")
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertEqual(output, "mode=full\n")

    def test_a_verified_reuse_exits_zero_outputs_mode_reuse_and_names_the_run(self):
        result, output, summary = self.gate_command("push", "success", RUN_URL, "skipped", "skipped")
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertEqual(output, "mode=reuse\n")
        for text in (result.stdout, summary):
            self.assertIn(RUN_URL, text)
            self.assertIn(TREE, text)

    def test_every_refused_combination_exits_nonzero_and_outputs_no_mode(self):
        for case in REFUSED:
            with self.subTest(case=case):
                result, output, summary = self.gate_command(*case)
                self.assertNotEqual(result.returncode, 0, result.stdout)
                self.assertIn("FAIL: CI did not pass", result.stderr)
                self.assertNotIn("mode=", output)
                self.assertNotIn("PASS", summary)


class ReuseCommandTests(unittest.TestCase):
    """The reuse command's outputs, which decide whether checks and tests run."""

    def reuse_command(self, get):
        with tempfile.TemporaryDirectory(prefix="fid-ci-reuse-") as directory:
            output, summary = Path(directory) / "output", Path(directory) / "summary"
            with mock.patch.dict(os.environ, {"GITHUB_OUTPUT": str(output), "GITHUB_STEP_SUMMARY": str(summary)}), \
                    mock.patch.object(ci, "github_get", side_effect=get), \
                    mock.patch.object(sys, "argv", ["ci_tests.py", "reuse", "--tree", TREE]), \
                    mock.patch("builtins.print"):
                ci.main()
            return output.read_text(), summary.read_text()

    def test_a_found_run_is_output_for_the_test_jobs_to_skip(self):
        run = {"id": 7, "event": "pull_request", "status": "completed", "conclusion": "success", "path": ci.WORKFLOW, "html_url": RUN_URL}
        get = lambda path: {"artifacts": [{"name": f"ci-tree-{TREE}", "expired": False, "workflow_run": {"id": 7}}]} if path.startswith("actions/artifacts?") else run
        output, summary = self.reuse_command(get)
        self.assertEqual(output, f"run={RUN_URL}\ntree={TREE}\n")
        self.assertIn(RUN_URL, summary)

    def test_no_run_or_a_failed_lookup_outputs_an_empty_run_so_everything_runs(self):
        def unreachable(path):
            raise OSError("network down")
        for get in (lambda path: {"artifacts": []}, unreachable):
            with self.subTest(get=get):
                output, summary = self.reuse_command(get)
                self.assertEqual(output, f"run=\ntree={TREE}\n")
                self.assertIn("Full run", summary)


def write_files(root, files):
    for name, text in files.items():
        path = root / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(text)


GO_MOD = "module fixture\n\ngo 1.26.2\n"


def spawning(package, tag=None):
    """Product code that starts a goroutine, optionally behind a build constraint."""
    return (f"//go:build {tag}\n" if tag else "") + f"package {package}\nfunc Twice(x int) int {{ done := make(chan int); go func() {{ done <- 2 * x }}(); return <-done }}\n"


def quiet(package, tag=None):
    return (f"//go:build {tag}\n" if tag else "") + f"package {package}\nfunc One() int {{ return 1 }}\n"


def go_test(package, name, body="", tag=None):
    return (f"//go:build {tag}\n" if tag else "") + f"package {package}\nimport \"testing\"\nfunc {name}(t *testing.T) {{ {body} }}\n"


def run_helper(root):
    return subprocess.run(["go", "run", str(ci.RACE_HELPER), "./..."], cwd=root, capture_output=True, text=True)


class RacePackagesTests(unittest.TestCase):
    def test_packages_whose_non_test_go_files_start_goroutines_are_found_by_parsing(self):
        with tempfile.TemporaryDirectory(prefix="fid-ci-race-") as directory:
            root = Path(directory)
            write_files(root, {
                "go.mod": "module fixture\n\ngo 1.26.2\n",
                "spawn/spawn.go": "package spawn\n\nfunc Start(f func()) { go f() }\n",
                "method/method.go": "package method\n\ntype W struct{}\n\nfunc (w *W) Run() {\n\tdone := make(chan int)\n\tgo func() { done <- 1 }()\n\t<-done\n}\n",
                "nested/nested.go": "package nested\n\nvar Run = func() {\n\tfunc() {\n\t\tgo println()\n\t}()\n}\n",
                "platform/platform.go": "package platform\n\nfunc F() {}\n",
                "platform/platform_plan9.go": "//go:build plan9\n\npackage platform\n\nfunc G() { go F() }\n",
                "testonly/testonly.go": "package testonly\n\nfunc F() int { return 1 }\n",
                "testonly/testonly_test.go": "package testonly\n\nimport \"testing\"\n\nfunc TestF(t *testing.T) { go F() }\n",
                "testonly/external_test.go": "package testonly_test\n\nimport \"testing\"\n\nfunc TestG(t *testing.T) { go t.Log() }\n",
                "decoy/decoy.go": "package decoy\n\n// go Start()\n/*\n\tgo func() {}()\n*/\n\nconst text = \"go Start()\\n\\tgo func() {}()\"\n\nvar raw = `\n\tgo Start()\n`\n\nfunc Loop() int {\n\tn := 0\nretry:\n\tn++\n\tif n < 3 {\n\t\tgoto retry\n\t}\n\treturn n\n}\n",
            })
            self.assertEqual(ci.race_packages(root), ["fixture/method", "fixture/nested", "fixture/platform", "fixture/spawn"])

    def test_packages_that_build_in_only_one_mode_are_found_and_run_in_that_mode(self):
        with tempfile.TemporaryDirectory(prefix="fid-ci-modes-") as directory:
            root = Path(directory)
            write_files(root, {
                "go.mod": GO_MOD,
                "plain/plain.go": quiet("plain"),
                "raceonly/raceonly.go": spawning("raceonly", tag="race"),
                "norace/norace.go": quiet("norace", tag="!race"),
            })
            self.assertEqual(ci.race_packages(root), ["fixture/raceonly"])
            self.assertEqual(ci.classify(root), (["fixture/norace", "fixture/plain", "fixture/raceonly"], ["fixture/raceonly"]))

    def test_packages_and_tests_that_could_never_run_in_their_mode_are_refused(self):
        cases = {
            "a race-only test in a package that runs without race": ({"p/p.go": quiet("p"), "p/p_test.go": go_test("p", "TestOnlyUnderRace", tag="race")}, "TestOnlyUnderRace"),
            "a non-race test in a package that runs under race": ({"s/s.go": spawning("s"), "s/s_test.go": go_test("s", "TestOnlyWithoutRace", tag="!race")}, "TestOnlyWithoutRace"),
            "an external race-only test": ({"x/x.go": quiet("x"), "x/x_test.go": go_test("x_test", "TestExternalUnderRace", tag="race")}, "TestExternalUnderRace"),
            "a race-only package that starts no goroutines": ({"r/r.go": quiet("r", tag="race")}, "fixture/r"),
            "a non-race package that starts goroutines": ({"n/n.go": spawning("n", tag="!race")}, "fixture/n"),
        }
        for label, (files, named) in cases.items():
            with self.subTest(label), tempfile.TemporaryDirectory(prefix="fid-ci-refuse-") as directory:
                root = Path(directory)
                write_files(root, {"go.mod": GO_MOD, **files})
                result = run_helper(root)
                self.assertNotEqual(result.returncode, 0, result.stdout)
                self.assertIn(named, result.stderr)

    def test_a_test_name_in_both_modes_is_accepted(self):
        with tempfile.TemporaryDirectory(prefix="fid-ci-paired-") as directory:
            root = Path(directory)
            write_files(root, {"go.mod": GO_MOD, "s/s.go": spawning("s"), "s/on_test.go": go_test("s", "TestBuildMode", tag="race"),
                               "s/off_test.go": go_test("s", "TestBuildMode", 't.Fatal("built without race")', tag="!race")})
            result = run_helper(root)
            self.assertEqual(result.returncode, 0, result.stderr)

    def test_unparseable_go_fails_closed(self):
        with tempfile.TemporaryDirectory(prefix="fid-ci-race-bad-") as directory:
            root = Path(directory)
            write_files(root, {"go.mod": "module fixture\n\ngo 1.26.2\n", "bad/bad.go": "package bad\n\nfunc F( {\n"})
            with self.assertRaises(ci.subprocess.CalledProcessError):
                ci.race_packages(root)

    def test_fids_race_set_is_derived_from_its_own_source(self):
        repository = Path(__file__).resolve().parents[2]
        self.assertIn("github.com/mbell-nzl/fid/internal/evidence", ci.race_packages(repository))


class LocalTestCommandTests(unittest.TestCase):
    """`make test` runs `ci_tests.py test`; its exit status is the local check's verdict."""

    def local_test(self, failing=None):
        def body(package):
            return f't.Fatal("{package} failed")' if package == failing else ""
        with tempfile.TemporaryDirectory(prefix="fid-ci-local-") as directory:
            root = Path(directory)
            write_files(root, {
                "go.mod": GO_MOD,
                "plain/plain.go": quiet("plain"), "plain/plain_test.go": go_test("plain", "TestPlain", body("plain")),
                "plain/on_test.go": go_test("plain", "TestBuildMode", 't.Fatal("plain ran under race")', tag="race"),
                "plain/off_test.go": go_test("plain", "TestBuildMode", tag="!race"),
                "spawn/spawn.go": spawning("spawn"), "spawn/spawn_test.go": go_test("spawn", "TestSpawn", body("spawn")),
                "spawn/on_test.go": go_test("spawn", "TestBuildMode", tag="race"),
                "spawn/off_test.go": go_test("spawn", "TestBuildMode", 't.Fatal("spawn ran without race")', tag="!race"),
                "raceonly/raceonly.go": spawning("raceonly", tag="race"),
                "raceonly/raceonly_test.go": go_test("raceonly", "TestRaceOnly", body("raceonly"), tag="race"),
            })
            return subprocess.run([sys.executable, SCRIPT, "test"], cwd=root, capture_output=True, text=True)

    def test_it_passes_when_every_package_passes_in_its_mode(self):
        result = self.local_test()
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        for package in ("plain", "spawn", "raceonly"):
            self.assertRegex(result.stdout, rf"ok\s+fixture/{package}\b")

    def test_a_failure_in_either_group_fails_it(self):
        for package in ("plain", "spawn", "raceonly"):
            with self.subTest(failing=package):
                result = self.local_test(failing=package)
                self.assertNotEqual(result.returncode, 0, result.stdout + result.stderr)
                self.assertIn(f"{package} failed", result.stdout)


def commit_module(root):
    for args in [["git", "init", "-q"], ["git", "add", "."], ["git", "-c", "user.name=CI test", "-c", "user.email=ci@example.test", "commit", "-qm", "fixture"]]:
        subprocess.run(args, cwd=root, check=True, capture_output=True, text=True)


class RealGoShardTests(unittest.TestCase):
    def test_a_failing_test_built_only_under_race_fails_the_shard(self):
        with tempfile.TemporaryDirectory(prefix="fid-ci-race-only-") as directory:
            root = Path(directory)
            write_files(root, {
                "go.mod": "module fixture\n\ngo 1.26.2\n",
                "spawn/spawn.go": "package spawn\nfunc Twice(x int) int { done := make(chan int); go func() { done <- 2 * x }(); return <-done }\n",
                "spawn/spawn_test.go": "package spawn\nimport \"testing\"\nfunc TestPlain(t *testing.T) { if Twice(2) != 4 { t.Fatal(\"twice\") } }\n",
                "spawn/race_only_test.go": "//go:build race\npackage spawn\nimport \"testing\"\nfunc TestRaceOnlyFails(t *testing.T) { t.Fatal(\"race-only test ran and failed\") }\n",
            })
            commit_module(root)
            result = subprocess.run([sys.executable, SCRIPT, "run", "--count", "1", "--index", "0", "--out", str(root / "shard-0")], cwd=root, capture_output=True, text=True)
            self.assertNotEqual(result.returncode, 0, result.stdout + result.stderr)
            self.assertIn("race-only test ran and failed", result.stdout + result.stderr)
            self.assertFalse((root / "shard-0" / "manifest.json").exists())

    def test_a_failing_package_built_only_under_race_fails_the_shard(self):
        with tempfile.TemporaryDirectory(prefix="fid-ci-race-package-") as directory:
            root = Path(directory)
            write_files(root, {"go.mod": GO_MOD, "plain/plain.go": quiet("plain"), "plain/plain_test.go": go_test("plain", "TestPlain"),
                               "raceonly/raceonly.go": spawning("raceonly", tag="race"),
                               "raceonly/raceonly_test.go": go_test("raceonly", "TestRacePackageFails", 't.Fatal("race-only package ran and failed")', tag="race")})
            commit_module(root)
            result = subprocess.run([sys.executable, SCRIPT, "run", "--count", "1", "--index", "0", "--out", str(root / "shard-0")], cwd=root, capture_output=True, text=True)
            self.assertNotEqual(result.returncode, 0, result.stdout + result.stderr)
            self.assertIn("race-only package ran and failed", result.stdout + result.stderr)
            self.assertFalse((root / "shard-0" / "manifest.json").exists())

    def test_one_shard_runs_and_merges_every_package_including_single_mode_ones(self):
        with tempfile.TemporaryDirectory(prefix="fid-ci-one-shard-") as directory:
            root = Path(directory)
            write_files(root, {"go.mod": GO_MOD, "plain/plain.go": quiet("plain"), "plain/plain_test.go": go_test("plain", "TestPlain", "if One() != 1 { t.Fatal(1) }"),
                               "raceonly/raceonly.go": spawning("raceonly", tag="race"),
                               "raceonly/raceonly_test.go": go_test("raceonly", "TestRaceOnlyPackage", "if Twice(2) != 4 { t.Fatal(4) }", tag="race"),
                               "norace/norace.go": quiet("norace", tag="!race"),
                               "norace/norace_test.go": go_test("norace", "TestNoRacePackage", "if One() != 1 { t.Fatal(1) }", tag="!race")})
            commit_module(root)
            result = subprocess.run([sys.executable, SCRIPT, "run", "--count", "1", "--index", "0", "--out", str(root / "results" / "shard-0")], cwd=root, capture_output=True, text=True)
            self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
            manifest = json.loads((root / "results" / "shard-0" / "manifest.json").read_text())
            self.assertEqual(manifest["race"], ["fixture/raceonly"])
            self.assertEqual(sorted(map(tuple, manifest["completed"])), [("fixture/norace", "TestNoRacePackage"), ("fixture/plain", "TestPlain"), ("fixture/raceonly", "TestRaceOnlyPackage")])
            merged = subprocess.run([sys.executable, SCRIPT, "merge", "--count", "1", "--root", str(root / "results"), "--out", str(root / "combined.out")], cwd=root, capture_output=True, text=True)
            self.assertEqual(merged.returncode, 0, merged.stdout + merged.stderr)
            self.assertIn("PASS coverage union: 1 shards, 3 targets", merged.stdout)
            blocks = ci.coverage_blocks(root / "combined.out")
            for package in ("plain", "raceonly", "norace"):
                self.assertTrue(any(key.startswith(f"fixture/{package}/") and hits > 0 for key, (_, hits) in blocks.items()), package)
            flat = subprocess.run([sys.executable, SCRIPT, "merge", "--count", "1", "--root", str(root / "results" / "shard-0"), "--out", str(root / "flat.out")], cwd=root, capture_output=True, text=True)
            self.assertEqual(flat.returncode, 0, flat.stdout + flat.stderr)
            self.assertEqual((root / "flat.out").read_text(), (root / "combined.out").read_text())
            wrong = subprocess.run([sys.executable, SCRIPT, "merge", "--count", "2", "--root", str(root / "results"), "--out", str(root / "wrong.out")], cwd=root, capture_output=True, text=True)
            self.assertNotEqual(wrong.returncode, 0)
            self.assertFalse((root / "wrong.out").exists())

    def test_real_go_tests_examples_fuzz_seeds_and_empty_package_are_reconciled(self):
        with tempfile.TemporaryDirectory(prefix="fid-ci-smoke-") as directory:
            root = Path(directory)
            (root / "go.mod").write_text("module fixture\n\ngo 1.26.2\n")
            (root / "fixture.go").write_text('package fixture\nfunc Value(b bool) int { if b { return 1 }; return 0 }\n')
            (root / "fixture_test.go").write_text('''package fixture
import ("fmt"; "testing")
func TestA(t *testing.T) { if Value(true) != 1 { t.Fatal("true branch") } }
func TestB(t *testing.T) { if Value(false) != 0 { t.Fatal("false branch") } }
func TestC(t *testing.T) { t.Skip("existing conditional skip remains visible") }
func TestD(t *testing.T) { t.Run("nested", func(t *testing.T) { if Value(true) != 1 { t.Fatal("nested") } }) }
func ExampleValue() { fmt.Println(Value(true)); // Output: 1
}
func FuzzValue(f *testing.F) { f.Add(true); f.Fuzz(func(t *testing.T, b bool) { v := Value(b); if v < 0 || v > 1 { t.Fatal(v) } }) }
''')
            (root / "empty").mkdir()
            (root / "empty" / "empty.go").write_text("package empty\nfunc Uncovered() int { return 7 }\n")
            mode = "//go:build {tag}\npackage {package}\nimport \"testing\"\nfunc TestBuildMode(t *testing.T) {{ {body} }}\n"
            write_files(root, {
                "spawn/spawn.go": "package spawn\nfunc Twice(x int) int { done := make(chan int); go func() { done <- 2 * x }(); return <-done }\n",
                "spawn/spawn_test.go": "package spawn\nimport \"testing\"\nfunc TestTwice(t *testing.T) { if Twice(2) != 4 { t.Fatal(\"twice\") } }\n",
                "spawn/race_on_test.go": mode.format(tag="race", package="spawn", body=""),
                "spawn/race_off_test.go": mode.format(tag="!race", package="spawn", body='t.Fatal("spawn starts goroutines but ran without the race detector")'),
                "race_on_test.go": mode.format(tag="race", package="fixture", body='t.Fatal("fixture starts no goroutines but ran under the race detector")'),
                "race_off_test.go": mode.format(tag="!race", package="fixture", body=""),
                # Distinct names that exist only in the build mode each package runs in.
                "spawn/race_only_test.go": "//go:build race\npackage spawn\nimport \"testing\"\nfunc TestOnlyUnderRace(t *testing.T) {}\n",
                "plain_only_test.go": "//go:build !race\npackage fixture\nimport \"testing\"\nfunc TestOnlyWithoutRace(t *testing.T) {}\n",
            })
            commit_module(root)
            for i in range(4):
                result = subprocess.run([sys.executable, SCRIPT, "run", "--count", "4", "--index", str(i), "--out", str(root / f"shard-{i}")], cwd=root, capture_output=True, text=True)
                self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
            ci.merge(root, 4, root / "combined.out")
            manifests = [json.loads((root / f"shard-{i}" / "manifest.json").read_text()) for i in range(4)]
            self.assertEqual(manifests[0]["race"], ["fixture/spawn"])
            completed = [tuple(item) for manifest in manifests for item in manifest["completed"]]
            self.assertIn(("fixture/spawn", "TestOnlyUnderRace"), completed)
            self.assertIn(("fixture", "TestOnlyWithoutRace"), completed)
            self.assertEqual(completed.count(("fixture/spawn", "TestOnlyUnderRace")), 1)
            blocks = ci.coverage_blocks(root / "combined.out")
            self.assertTrue(any(key.startswith("fixture/spawn/") and hits > 0 for key, (_, hits) in blocks.items()))
            self.assertTrue(any(key.startswith("fixture/empty/") and hits == 0 for key, (_, hits) in blocks.items()))
            self.assertTrue(all(hits > 0 for key, (_, hits) in blocks.items() if key.startswith("fixture/fixture.go:")))
            result = subprocess.run(["go", "tool", "cover", "-func", str(root / "combined.out")], cwd=root, capture_output=True, text=True)
            self.assertEqual(result.returncode, 0, result.stdout + result.stderr)


if __name__ == "__main__":
    unittest.main()
