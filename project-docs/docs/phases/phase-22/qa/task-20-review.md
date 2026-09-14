# Phase 22 task 20 — Independent performance-repair review

Reviewer: independent bounded reviewer (`task20_review`)
Date: 2026-09-14
Review type: INITIAL
Authority: task-20.md and DEC-221; one review, no product or control edits.

## Verdict

PASS — no concrete correctness defect or lost original assertion/case found in
the reviewed repair. Must-fix findings: 0. Blocker candidates: 0.
This is a bounded code review, not task closure or hosted-CI acceptance.

## Scope and findings

Reviewed task-20.md, DEC-221, the task's validation assertion rows, the working
tree diff for all changed Go tests and `.github/workflows/ci.yml`, new
`internal/export/export_source_fixture_test.go`, `project-docs/tools/ci_tests.py`
and `test_ci_tests.py`, Makefile, and the existing fixture helpers used by the
changes. Only this review artifact was deliberately edited by the reviewer.

- **Original behavioral coverage retained.** The Stop matrix retains both
  surfaces and all five states, with its existing result/count/persistence/spawn
  assertions. The export/import diffs replace repeated preparation only; their
  corruption rows, publication stages, controls, expected errors, and
  non-mutation assertions remain. No production or migration change occurs in
  the reviewed diff.
- **Stop fixture isolation and execution root.**
  `cmd/fid/stop_operation_matrix_test.go:169` uses a stable relative marker to
  reuse the existing definition-keyed fixture cache. The existing cache copies
  the template before returning the first mutable workspace and physically
  copies later workspaces. The added regression runs two clones, checks each
  marker in its own root, and checks that a later clone is uncontaminated.
- **Export source/history fidelity.**
  `internal/export/export_source_fixture_test.go:23` builds source history with
  real bootstrap/journal writers and exports once. Its clone method writes
  separate regular files and a separate archive, retaining their bytes without
  identity rewriting or a shared mutable database. The regression checks DB
  physical separation, database mutation, archive mutation, sibling/source
  preservation, a later clone, and actual archive validation. These copies are
  only archive inputs; the reviewed call sites still invoke actual validation
  and import. The fixture comment correctly limits use of path-bound history.
- **Changed-digest and publication checks remain non-vacuous.**
  `internal/export/import_boundary_test.go:83` still calls the unchanged
  `boundaryImportFixture` for the second archive in the changed-digest case.
  That archive comes from a separately initialized source, rather than another
  copy of the installed archive. Each import case receives a fresh destination;
  crash-helper execution, retry, and installed-history checks remain real.
- **Exhaustive selection.** `project-docs/tools/ci_tests.py:12` partitions the
  discovered top-level Test/Example/Fuzz names without a maintained allowlist.
  Equal names in different packages share an owner, so the all-package anchored
  `-run` expression cannot accidentally select that name on another shard.
  New names enter the plan automatically. Empty selections and duplicate or
  invalid discovery names are rejected. All packages remain in each Go command,
  including packages with no selected test or no tests at all.
- **Execution reconciliation.** `ci_tests.py:42` requires exact agreement among
  selected, actually started, and completed top-level targets, plus terminal
  accounting for every discovered package. Failed child/subtest/package events,
  duplicate starts, missing targets, and extra targets reject the result. Native
  Go test selection executes subtests and fuzz seeds under the selected parent;
  existing conditional skips remain visible and retain ordinary-suite semantics.
  `run_go` waits for the Go process and rejects nonzero termination before a
  successful manifest can be emitted.
- **Coverage completeness.** `ci_tests.py:106` accepts atomic profiles only.
  Merge requires every shard index, matching source/inventory/count, exact
  selection/execution/completion manifests, and identical block coordinates and
  statement counts. It adds hit counts and retains zero-hit blocks, including
  no-test packages. Missing, extra, malformed, empty, or inconsistent evidence
  fails before a combined profile is written. Atomic here is the Go coverage
  mode; no file-publication atomicity claim is made.
- **Required quality check.** `.github/workflows/ci.yml` retains format, lint,
  build and race/coverage instrumentation, removes duplicate full-suite testing,
  and sets matrix `fail-fast: false`. The required `quality` job runs with
  `always()` and explicitly requires both dependency results to equal success.
  Failed/cancelled/skipped dependencies therefore cannot satisfy that condition;
  subsequent merge/report failures also fail the job. Artifact transfer and
  workflow cancellation behavior were assessed statically, not on GitHub.

## Evidence and limits

Working directory for the review command: `/home/mark/projects/fid`.

- `PYTHONPATH=project-docs/tools python3 -m unittest test_ci_tests.ShardTests -v`
  — PASS — `Ran 7 tests in 0.017s`; `OK`; last printed line
  `PASS coverage union: 4 shards, 23 targets, 3 statement blocks`
  (terminated `exit_code: 0`, observed by this reviewer).
- Inspected root-produced `.cache/phase22-pr-performance/fixture-green.log`
  and `shard-integration.log`. They contain focused Go package successes and
  the real four-shard tiny-Go integration success, respectively. These are
  supporting records only: their execution termination was not observed in
  this reviewer session, so their gate status is UNVERIFIED by this reviewer.
- Root is responsible for the task's exact final commands and full before/after
  timing receipts. No full suite, heavy Go execution, benchmark, tagged journey,
  or hosted CI run was performed by this reviewer. Static balancing by sorted
  name does not guarantee equal shard duration; hosted wall-time improvement
  and resolution of the original PR CI failure require the actual hosted run.

VAL-036..039/075 and VAL-054..065: reviewed for preservation of existing evidence
paths, with no assertion definition change or replacement of historical QA.
This review does not independently requalify those product assertions. Fixture
isolation and CI scheduling/reconciliation checks have **no independent VAL**;
they protect the accepted validation allocation. Prior Human acceptance and
phase QA remain historical evidence, unchanged by this review.
