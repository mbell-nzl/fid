"""Behavioral checks for complete CI selection and coverage reconciliation."""
import importlib.util
import json
from pathlib import Path
import tempfile
import unittest

spec = importlib.util.spec_from_file_location("ci_tests", Path(__file__).with_name("ci_tests.py"))
ci = importlib.util.module_from_spec(spec)
spec.loader.exec_module(ci)


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

    def write_shards(self, root):
        parts = ci.partitions(self.inventory, 4)
        for i, part in enumerate(parts):
            path = root / f"shard-{i}"
            path.mkdir()
            manifest = {"index": i, "count": 4, "source": "same-source", "inventory": self.inventory, "selected": sorted(part), "executed": sorted(part), "completed": sorted(part), "packages": sorted(self.inventory)}
            (path / "manifest.json").write_text(json.dumps(manifest))
            (path / "coverage.out").write_text(f"mode: atomic\nexample/a/a.go:1.1,3.2 1 {int(i == 0)}\nexample/a/a.go:4.1,6.2 1 {int(i == 1)}\nexample/empty/e.go:1.1,3.2 1 0\n")

    def test_coverage_union_retains_all_blocks_and_sums_counts(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            self.write_shards(root)
            ci.merge(root, 4, root / "combined.out")
            self.assertEqual((root / "combined.out").read_text(), "mode: atomic\nexample/a/a.go:1.1,3.2 1 1\nexample/a/a.go:4.1,6.2 1 1\nexample/empty/e.go:1.1,3.2 1 0\n")

    def test_missing_duplicate_or_inconsistent_shards_fail_closed(self):
        for mutation in ["missing", "duplicate", "source", "inventory", "selected", "executed", "completed", "packages", "count", "empty-profile", "missing-block", "wrong-mode", "statements", "negative-count"]:
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
                else:
                    profile = path / "coverage.out"
                    text = profile.read_text()
                    text = {"empty-profile": "mode: atomic\n", "missing-block": "\n".join(text.splitlines()[:-1]) + "\n", "wrong-mode": text.replace("atomic", "set"), "statements": text.replace("3.2 1", "3.2 2"), "negative-count": text.replace("1 0", "1 -1")}[mutation]
                    profile.write_text(text)
                if mutation != "missing":
                    mpath.write_text(json.dumps(m))
                with self.assertRaises(ValueError):
                    ci.merge(root, 4, root / "combined.out")


class RealGoShardTests(unittest.TestCase):
    def test_real_go_tests_examples_fuzz_seeds_and_empty_package_are_reconciled(self):
        import subprocess
        import sys
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
            for args in [["git", "init", "-q"], ["git", "add", "."], ["git", "-c", "user.name=CI test", "-c", "user.email=ci@example.test", "commit", "-qm", "fixture"]]:
                subprocess.run(args, cwd=root, check=True, capture_output=True, text=True)
            script = str(Path(__file__).with_name("ci_tests.py").resolve())
            for i in range(4):
                result = subprocess.run([sys.executable, script, "run", "--count", "4", "--index", str(i), "--out", str(root / f"shard-{i}")], cwd=root, capture_output=True, text=True)
                self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
            ci.merge(root, 4, root / "combined.out")
            blocks = ci.coverage_blocks(root / "combined.out")
            self.assertTrue(any(key.startswith("fixture/empty/") and hits == 0 for key, (_, hits) in blocks.items()))
            self.assertTrue(all(hits > 0 for key, (_, hits) in blocks.items() if key.startswith("fixture/fixture.go:")))
            result = subprocess.run(["go", "tool", "cover", "-func", str(root / "combined.out")], cwd=root, capture_output=True, text=True)
            self.assertEqual(result.returncode, 0, result.stdout + result.stderr)


if __name__ == "__main__":
    unittest.main()
