# Phase 20 - Task 4: Rebuild, authority, and lock failures

## Recovery Directive (DEC-136; replaces the withdrawn DEC-135 directive)

Task 4's implementation is **done and committed** as
`e56b3a0ce0433f1044f70b7969c14e19234dfdf8`
(`feat(rebuild): map typed authority and lock failures`). Scope, acceptance
criteria, VAL coverage, and the command list below are unchanged, and **no code
change is in scope**. This session repairs the completion report's evidence,
which failed `report-lint.sh` on four checks.

Follow your own role prompt's Report And Commit Sequence exactly. This spec does
not restate or override it (DEC-136 decision 2). This session:

1. Verifies the evidence will apply to the reviewed tree: `git rev-parse HEAD`
   is `e56b3a0`, and `git status --porcelain` shows no source modification. If
   the tree differs, stop with a blocker instead of recording gate results
   against a different state.
2. Re-runs `## Exact Commands To Run` below, in order and complete, plus the
   focused test command recorded under `## Tests run`. Do **not** reconstruct,
   expand, or reuse the previously recorded fragments `0.009s`,
   `coverage: 80.5% of statements`, or `4.625s`; recalling evidence from an
   earlier session is `UNVERIFIED`, never `PASS`. Real values from this run will
   differ from those fragments — record what your terminal actually printed and
   never reconcile it to the old numbers.
3. Records each result with one short verbatim token copied from that command's
   own last printed line. For any `go test` command the line must carry its real
   marker (`ok`, `FAIL`, `PASS`, `---`, or `no test files`) with its package
   path and duration, not a bare duration or a bare `coverage:` fragment; for
   `go tool cover -func=coverage.out` copy the whole final `total:` line.
   `exit 0 (no output)` remains valid for `gofmt` and `go build`, never for
   `go test`. The race and dedicated integration gates print nothing for
   minutes before finishing; wait for them and capture `exit_code` or an
   explicit termination signal.
4. Rewrites `task-4-done.md` with those tokens, re-records every VAL row from
   this session's own evidence, and names `e56b3a0ce0433f1044f70b7969c14e19234dfdf8`
   under `## Commit hash` with its exact subject under `## Commit subject`.
5. Keeps the existing deviation note about the first session's lint stop and
   adds one entry that quotes the three superseded fragments verbatim, states
   they were replaced by a fresh capture at the same tree state, and cites
   DEC-136.
6. Commits the finalized report as
   `docs(task): finalize completion report for phase-20 task-4`, then runs
   `project-docs/tools/report-lint.sh` against the committed file and iterates
   with a further `docs(task): ...` commit if it still fails. Do not amend,
   revert, or force-push `e56b3a0`, and do not commit `run-ledger.md`,
   `triage-log.md`, or `tl_escalations_current.md`.

If a re-run gate now fails, that is a genuine finding: stop and file a blocker
report rather than repairing code under a report-repair task.

## Objective

Migrate `fid rebuild` to the project-selection and Response Envelope contract,
map retained fail-closed V1 authority errors, and make write-lock contention a
typed retryable failure.

## Validation Contract Coverage

- VAL-001, VAL-002, VAL-006, VAL-007, VAL-008, VAL-017, and the rebuild row of VAL-044.

## Binding Constraints

- V1 permission checks remain fail-closed and execute before journal/projection work.
- Map existing missing, unknown, and denied cases to the semantic authority codes;
  do not weaken checks or promise those general codes will be deleted in Phase 21.
- Rebuild is projection-only and returns `event: null`.
- Lock tests acquire and release the lock deterministically; timing races are invalid evidence.

## Scope

### Files To Create Or Modify

- `internal/journal/lock.go` and tests.
- `cmd/fid/rebuild.go`, shared error mapping, and rebuild/root tests.

### Files To Read For Context

- `internal/policy/authorize.go`; `../contracts.md` sections 2.2-2.5, 3.2, and 3.5.

### Out Of Scope

- Generic Actor/Envelope authority, Claims/Runs, idempotency, and Issue schema.

## Acceptance Criteria

- [ ] Rebuild accepts only its closed request and returns its exact result with `event:null`. (VAL-017, VAL-044)
- [ ] Missing, unknown, and denied V1 authority map to exact codes/details and append nothing. (VAL-002, VAL-006)
- [ ] Held write lock returns `lock-unavailable`, `retryable:true`; release allows retry. (VAL-007)
- [ ] Human success is stdout-only; Human and JSON failures are stderr-only. (VAL-001, VAL-008)

## First Failing Test

`cmd/fid/rebuild_test.go` - `TestRebuildJSONMapsMissingAuthorityWithoutMutation`.

## Documentation Updates

- None expected.

## Exact Commands To Run

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go build -o /tmp/fid ./cmd/fid
```

## Notes For Developer

- Completion report: `task-4-done.md`.
