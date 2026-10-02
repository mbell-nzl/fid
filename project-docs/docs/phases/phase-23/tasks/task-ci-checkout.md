# Phase 23 - CI Checkout Repair

## Objective

Make the accepted migration baseline available to PR test workers. PR #25's
`tests (2)` failed because its shallow checkout omitted the pinned planning
commit; the dependent `quality` check then also failed. After restoring full
history, Git's adaptive hash abbreviation also changed the diff headers from
seven to eight characters. Request the expected abbreviation explicitly.

This is a scoped post-phase repair authorized by the reported CI failure.
Phase 23 remains complete; this task does not reopen product implementation.

## Validation Contract Coverage

- VAL-075: migration source changes match the seven accepted baseline hunks.
  Its existing test must run with the required Git history available.

## Scope

- `.github/workflows/ci.yml`: fetch full history for test workers.
- `internal/db/review_judgment_migration_test.go`: request `--abbrev=7` for the
  baseline diff; retain the pinned commit and exact expected diff.
- `project-docs/tools/README.md`: document the Git-history prerequisite.
- This task specification and `task-ci-checkout-done.md`: diagnostic and
  verification evidence.

The accepted migration files, test assertions, shard selection, coverage
reconciliation, required-check semantics, budgets, and dependencies are fixed.

## Acceptance Criteria

- Reproduce the reported failure in a fresh depth-one local clone.
- Run the same migration assertion successfully in that clone after fetching
  full history and specifying the expected hash abbreviation (VAL-075).
- The expected diff and migration files remain byte-identical.
- Run the exact failing shard entrypoint successfully with full history.
- Existing shard-selection and coverage-reconciliation controls pass.

## First Failing Test

`TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks` in the shallow clone.

## Exact Commands To Run

All commands run from the repository root of
`/tmp/fid-ci-checkout-5mZiTF/repo`. It is a fresh local clone of the unchanged
Phase 23 product source at `7af47be`. The one-option test change is then copied
from the main worktree; no product source changes. Shard output goes into a
fresh directory under the scratch root. Format and lint checks run in the main
worktree at `/home/mark/projects/fid`.

```text
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/db -run '^TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks$' -v
python3 project-docs/tools/ci_tests.py run --count 4 --index 2 --out /tmp/fid-ci-checkout-5mZiTF/shard-2
python3 -m unittest discover -s project-docs/tools -p test_ci_tests.py -v
gofmt -l cmd internal
golangci-lint run ./...
```

The first command's initial failure is retained as red evidence; its identical
rerun must pass after full history is available and the one-option test fix
is applied. Preserve the intervening abbreviation failure as a separate
failed attempt. For this CI and test-portability repair, these commands replace
the full task gate: the actual failed CI shard is run
with its unchanged race and coverage flags, and the existing Python controls
exercise selection and coverage merging. No Go implementation changes or
new dependencies require further ordinary-suite, integration, or benchmark runs.
