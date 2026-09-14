# Phase 22 — Task 20: Shorten PR feedback — Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: Local gates, independent review and hosted PR verification complete; see additive hosted closure below.
- Recovery attempts: 0
- Next action: Human owns merge/push; no further performance repair required.
- Human attention required: no

## Summary of what was implemented

- Stop matrix uses relative process markers and the existing definition-keyed fixture cache. Two cloned workspaces execute independently; a third clone proves the source remains untouched.
- Five export/import matrices prepare immutable historical input once and physically copy it per case. Real archive validation, import/rebuild, faults and corruptions remain. Changed-digest validation still builds a genuinely different source archive.
- CI runs one complete ordinary race/coverage suite across four workers. Discovery includes Test, Example and Fuzz seed targets; exact execution reconciliation rejects missing/extra runs. All packages, including packages without tests, remain covered.
- Required quality check fails on failed/cancelled dependencies and verifies every shard's source, inventory, completed targets and full atomic coverage block set before combining counts. Format, lint and build remain.
- No production Go, timeout, test assertion, archive limit, benchmark budget, migration or runtime dependency changed. Makefile local targets remain available.

## Files created or modified

- cmd/fid/stop_operation_matrix_test.go.
- internal/export/export_source_fixture_test.go; export_test.go; archive_v2_test.go; export_boundary_test.go; import_boundary_test.go.
- .github/workflows/ci.yml; project-docs/tools/ci_tests.py; project-docs/tools/test_ci_tests.py; .gitignore (only Python bytecode).

## Documentation updated

- Task20, DEC-221, quality-plan CI cadence, tools README and PR description.
- Original supplied failure preserved; diagnosis, escalation archive/current index, phase/current progress, ledger and final-acceptance follow-up.
- Independent qa/task-20-review.md: PASS, no must-fix. Reviewer-owned seven Python unit tests passed with exit_code 0; root's receipts below are separate.

## Validation Contract Coverage

- VAL-036..039/075 retained: Stop matrix and real correct-root execution regression passed.
- VAL-054..065 retained: all affected archive/import matrices passed, including corruption, inventory, special entry, publication fault and changed-digest/installed-history cases.
- No product VAL changed. Source/clone isolation and exhaustive CI/coverage reconciliation are evidence-control regressions.
- Matching uncached logs contain all 2,653 original executed tests/subtests, with exactly two new fixture regressions and zero removals; comparison.txt records exact added identities.

## Tests run

All root-owned raw captures below are under .cache/phase22-pr-performance/.
Actual nested process exit statuses were observed in the executing session.

- `go test -race -p 1 ./cmd/fid ./internal/export -run '^(TestStopOperationFixtureExecutesInEachClonedWorkspace|TestExportSourceFixtureCopiesKeepDatabaseAndArchiveMutationsIsolated)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  fixture-red.log: wrong-root marker and shared mutable fixture paths rejected.
- `go test -race -p 1 ./cmd/fid ./internal/export -run '^(TestStopOperationFixtureExecutesInEachClonedWorkspace|TestExportSourceFixtureCopiesKeepDatabaseAndArchiveMutationsIsolated)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	3.739s` (terminated `exit_code: 0`)
  fixture-green.log: CLI9.218s, export3.739s.
- `go test -race -p 1 ./cmd/fid ./internal/export -run '^(TestStopPredicatesAcrossRunStartAndGateRunCommands|TestValidateRejectsEveryArchiveInvalidCheckWithoutMutation|TestValidateRecoveryManifestClosedConfigurationAndInventory|TestArchiveRejectsSpecialEntriesAndExpandedMember|TestImportProcessExitAtPublicationBoundaries|TestImportRejectsChangedDigestAndInstalledHistory)$' -count=1 -json` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	60.570s` (terminated `exit_code: 0`)
  matrices-after.jsonl: all six original affected matrices; CLI20.604s/export60.570s. These serial focused timings are not directly comparable to a contended full-suite run.
- `python3 -m unittest discover -s project-docs/tools -p test_ci_tests.py -v` - FAIL - `FAILED (failures=7, errors=17)` (terminated `exit_code: 1`)
  shard-red.log: deliberately incomplete controller fails selection, execution and coverage regressions.
- `python3 -m unittest discover -s project-docs/tools -p test_ci_tests.py -v` - PASS - `OK` (terminated `exit_code: 0`)
  shard-green.log: seven unit tests. Final expanded eight-test run below adds an actual tiny Go module with four race/coverage workers, skipped test, subtest, example, fuzz seeds and a package with uncovered code.

## Exact quality-gate results

Working directory: /home/mark/projects/fid. All seven allocated commands completed.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -count=1 -json ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	503.550s` (terminated `exit_code: 0`)
  before.jsonl: uncached baseline, original sources; CLI541.652s/export368.615s/workflow503.550s.
- `go test -race -count=1 -json ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	464.948s` (terminated `exit_code: 0`)
  after.jsonl: matching full-suite comparison; CLI483.994s/export264.146s/workflow464.948s. Workflow source unchanged; less contention is a plausible explanation, not a measured causal allocation. First-to-last JSON event span616.574s ->574.741s, excluding pre-event startup/compilation.
- `make coverage` - PASS - `65.5%` (terminated `exit_code: 0`)
  coverage.log: full race coverage and go tool cover report; unchanged total65.5%.
- `make build` - PASS - `go build -o /tmp/fid ./cmd/fid` (terminated `exit_code: 0`)
  build.log.
- `python3 -m unittest discover -s project-docs/tools -p test_ci_tests.py -v` - PASS - `OK` (terminated `exit_code: 0`)
  shard-integration.log: eight tests,11.631s; actual four-shard tiny-module coverage union passes.
- `python3 project-docs/tools/ci_tests.py plan --count 4` - PASS - `PASS exhaustive plan: 21 packages, 724 targets` (terminated `exit_code: 0`)
  plan.log: four disjoint selections,181 top-level targets each. This is discovery evidence; full real fid hosted shard execution remains pending.

## Deviations from spec

- DEC-221 explicitly changes CI scheduling/cadence while preserving scope, race, timeout and coverage semantics. Root used the Human-approved serial repair with one independent review.
- Fixture red failures (verbatim final token FAIL, exit_code1) exposed absolute marker binding and shared mutable fixture paths. Relative argv plus independent physical copies repaired their causes; the byte-identical command passed, retaining all existing assertions.
- Controller red failure (FAILED (failures=7, errors=17), exit_code1) preceded implemented partitioning/reconciliation; the byte-identical command passed after implementation. Expanded real-Go integration then passed without reducing the earlier checks.
- Required local gates did not fail. The original hosted timeout remains immutable failed evidence in qa/pr-CI-test-failure-log.md and original-ci.log. A hosted success is required before closing that failure.
- No extra tagged integration, benchmark or archive_scale rerun: task20 explicitly allocates ordinary-fixture/CI-only qualification; production and tagged fixtures are unchanged. Earlier phase evidence remains attributed to its original executors.

Administrative checkpoint: first report-lint terminated exit_code5 with
`FAIL report-lint found 5 failure(s)`: four Go tokens were too short/JSON
fragments for its output recognizer and the draft report was uncommitted.
Copied complete existing package-output lines and committed the report; no
command, result, termination or verification control changed.

## Concerns, technical debt, or notes for the Tech Lead

- Local fixture savings: CLI10.6%, export28.3%; total observed JSON event span6.8%. These do not predict hosted wall time.
- Original hosted Test step31m33s, quality job34m35s. Failed package work was followed by remaining queued work; it was not all initial compilation. SetupGo12s, lint install106s and lint57s are separate.
- Four workers consume concurrent runner capacity and may repeat compilation/package setup. This trades parallel resources for feedback time; no aggregate runner-minute saving is claimed.
- Hosted action execution, artifact transport, full fid shard reconciliation and hosted runtime are UNVERIFIED at this checkpoint. No new Human checklist is needed for unchanged product behavior, but merge still awaits CI.

## Dependency-budget note

- Python3.10+ standard library only; no pip/Go/runtime addition. Existing hosted/local Python is sufficient. Alternative shell JSON/coverage parsing is more fragile.
- Official artifact actions avoid custom authenticated API transfer and a duplicate reporting-suite run. Upload v7.0.1 and download v8.0.1 are MIT-licensed, verified from exact tagged sources: [upload release](https://github.com/actions/upload-artifact/releases/tag/v7.0.1), [upload license](https://raw.githubusercontent.com/actions/upload-artifact/v7.0.1/LICENSE), [download release](https://github.com/actions/download-artifact/releases/tag/v8.0.1), [download license](https://raw.githubusercontent.com/actions/download-artifact/v8.0.1/LICENSE).

## Commit hash

- `ecb2722b65db8eba1a8db096b3ebd0dd9bfef854`

## Commit subject

- `test(phase-22): shorten CI with isolated fixtures and exhaustive shards`

## Hosted closure — 2026-09-14

The earlier pending statements above describe the pre-push checkpoint and are
superseded by [hosted qualification](../qa/task-20-hosted-ci.md). GitHub run
34805256766 on6b8fe50bcd47148572ba841d561b23348f041dad succeeded in13m53s. All
six jobs succeeded;724 targets and12967 statement blocks reconciled;65.5%
coverage. Root observed exit_code0 for each gh retrieval, with full metadata
and logs retained. The hosted CI failure is resolved; task20 COMPLETE.
Combined execution47m48s, approximately50 per-job rounded runner minutes; no
account billing receipt or cost-saving claim. These closure records change
documentation only; no extra product gate run is claimed.
