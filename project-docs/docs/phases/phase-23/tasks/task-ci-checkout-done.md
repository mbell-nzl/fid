# Phase 23 - CI Checkout Repair - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Test-worker checkouts now use `fetch-depth: 0`, supplying the historical
  planning commit needed by the migration baseline assertion.
- The baseline test requests `--abbrev=7`, matching its accepted diff headers
  independently of Git's adaptive default. Its expected diff is unchanged.
- PR #25 reported `tests (2)` and `quality` as FAILURE; checks and the other
  three shards were SUCCESS. The workflow requires every shard to succeed
  before `quality` can reconcile coverage.

## Files created or modified

- `.github/workflows/ci.yml`
- `internal/db/review_judgment_migration_test.go`
- `project-docs/tools/README.md`
- `project-docs/docs/phases/phase-23/tasks/task-ci-checkout.md`
- `project-docs/docs/phases/phase-23/tasks/task-ci-checkout-done.md`

## Documentation updated

- Repo tools README records full-history and explicit-abbreviation prerequisites.
- Task specification records the scope and required verification commands.
- This report retains both failed attempts and the passing verification.

## Validation Contract Coverage

- VAL-075: PASS - `TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks`
  passes in the full-history clone with the original pinned commit and exact
  expected diff. The migration files and expected diff were not edited.
  The actual shard entrypoint also completes all 208 selected targets.

## Tests run

Working directory: `/tmp/fid-ci-checkout-5mZiTF/repo`, a fresh local clone of
`7af47be2312a00a6a5770a7d28ae48b39d5e78d7`. The clone started at depth one;
`git fetch --quiet --unshallow origin` restored full history from the local
repository (terminated `exit_code: 0`). Only the one-option test patch was
copied into it. Shard output uses the fresh scratch directory named below.

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/db -run '^TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks$' -v` - FAIL - `FAIL` (terminated `exit_code: 1`); initial shallow-checkout attempt.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/db -run '^TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks$' -v` - FAIL - `FAIL` (terminated `exit_code: 1`); full-history attempt before the abbreviation fix.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/db -run '^TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks$' -v` - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	1.026s` (terminated `exit_code: 0`)

## Exact quality-gate results

Shard and Python controls ran in the clone above. Format and lint ran in
`/home/mark/projects/fid` after `gofmt -w` and `goimports -w` on the edited Go
file. Each result was captured in this repair session.

- `python3 project-docs/tools/ci_tests.py run --count 4 --index 2 --out /tmp/fid-ci-checkout-5mZiTF/shard-2` - PASS - `PASS shard 2: 208 selected targets completed; all 21 packages accounted for` (terminated `exit_code: 0`)
- `python3 -m unittest discover -s project-docs/tools -p test_ci_tests.py -v` - PASS - `PASS coverage union: 4 shards, 23 targets, 3 statement blocks` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `git diff --cached --check` - FAIL - `project-docs/docs/phases/phase-23/tasks/task-ci-checkout-done.md:106: new blank line at EOF.` (terminated `exit_code: 2`); additional staged whitespace check.
- `git diff --cached --check` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-ci-checkout-done.md project-docs/docs/phases/phase-23/tasks/task-ci-checkout.md` - FAIL - `FAIL report-lint found 6 failure(s)` (terminated `exit_code: 6`)
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-ci-checkout-done.md project-docs/docs/phases/phase-23/tasks/task-ci-checkout.md` - PASS - `PASS report-lint completed with no failures` (terminated `exit_code: 0`)

All eight Python controls passed, including real Go execution and coverage
reconciliation. The passing staged whitespace check follows removal of the
extra final blank line.

The shard retained raw Go output in
`/tmp/fid-ci-checkout-5mZiTF/shard-2/tests.jsonl`, discovery output,
`coverage.out`, and `manifest.json`. Observed package times include
`cmd/fid` 78.380s, `internal/db` 20.379s, and `internal/workflow` 19.431s.

## Deviations from spec

- First failure: `diff accepted baseline: exit status 128`, with the Go
  command's final token `FAIL` and termination `exit_code: 1`. The underlying
  Git diagnostic was `fatal: bad object 0195897a347dfe1872db0f7be3800c5c3ccdefcf`
  (terminated `exit_code: 128`). The shallow clone lacked the required commit.
  Full history supplies it; the byte-identical Go command later passes.
- The intervening full-history attempt still ended `FAIL` (terminated
  `exit_code: 1`): Git emitted `index 01899a83..d034272c 100644`, while the
  accepted fixture requires `index 01899a8..d034272 100644`. The SQL hunks
  matched. Adding `--abbrev=7` to the test's Git invocation fixes the formatting
  without changing its baseline, expected bytes, or accepted migration edits.
  The byte-identical Go command then passes (terminated `exit_code: 0`).
- The scoped verification exception is recorded in `task-ci-checkout.md`.
  The unchanged failed CI shard, existing shard controls, format, and lint
  replace a full phase gate for this repair. Phase acceptance is not reasserted.
- The initial staged whitespace check also reported
  `project-docs/docs/phases/phase-23/tasks/task-ci-checkout-done.md:106: new blank line at EOF.`
  (terminated `exit_code: 2`). Report generation added an extra final blank
  line; report finalization removed it, and the byte-identical rerun completed
  with `exit 0 (no output)` (terminated `exit_code: 0`).
- Report lint initially ended `FAIL report-lint found 6 failure(s)` (terminated
  `exit_code: 6`) because PASS rows carried punctuation or commentary after
  their termination receipt. The output tokens and observed statuses were
  correct. Each PASS row now ends at its receipt; the byte-identical command
  then ended `PASS report-lint completed with no failures` (terminated
  `exit_code: 0`) against the committed report.

## Concerns, technical debt, or notes for the Tech Lead

- Hosted GitHub Actions has not rerun this repair. Local verification proves
  the previously failed shard succeeds with full history and the test patch.
- The scratch clone's manifest identifies its starting commit; the test patch
  was an explicit uncommitted change during verification.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `d15893d8bda9bfb973103412981e680171546e56`

## Commit subject

- `fix(ci): make migration baseline checks reproducible`
