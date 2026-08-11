# Phase 20 - Task 4: Rebuild, authority, and lock failures

## Recovery Directive (DEC-135, supersedes nothing below)

The first Task 4 session implemented this task but stopped on the required lint
gate with `internal/journal/lock.go:185:6: func describeLockHolder is unused
(unused)`, produced no commit, and left the work uncommitted in the tree. Scope,
acceptance criteria, VAL coverage, and the command list below are unchanged.
This session:

1. Inherits the existing uncommitted working tree. Do not `stash`, `checkout`,
   `reset`, or otherwise discard it, and do not re-implement from scratch.
2. Deletes `describeLockHolder` from `internal/journal/lock.go`. It is dead code:
   `contracts.md` §2.4/§2.5 close `lock-unavailable` to the exact detail object
   `{lock_path: string}` and forbid environment values in details, so holder
   PID/hostname/command/actor has no contracted destination. Do not suppress the
   finding with `//nolint`, a blank assignment, or any other fake reference, and
   do not add holder identity to the envelope. Keep the `owner.json` metadata
   write and `readLockMetadata`. If the deletion cascades to further unexported
   symbols inside `lock.go`, delete those as well; if it reaches exported API or
   another package, stop and file a blocker report.
3. Before running gates, re-reads the inherited tests and confirms each
   acceptance criterion below is actually asserted. Add the missing test first,
   red then green, if any criterion is unproven.
4. Runs `## Exact Commands To Run` in the listed order, complete, from a clean
   start. Record each with one verbatim token from that command's own output
   under the DEC-131/DEC-132 rules: `exit 0 (no output)` is never valid for a
   `go test` command, and a command with no captured `exit_code` or explicit
   termination signal is `UNVERIFIED`/`BLOCKED`, never `PASS`. The race and
   dedicated integration gates can run silently for minutes; wait for them.
5. Produces exactly one commit containing the code, tests, and `task-4-done.md`.
   Do not commit `run-ledger.md`, `triage-log.md`, or
   `tl_escalations_current.md`; those belong to the Phase Runner and Tech Lead.
6. Rewrites `task-4-done.md` in place as the full completion report, replacing
   the blocker report. Its `## Deviations from spec` section must quote the
   first session's failing token verbatim, name the fix, and state that the
   repair and passing gates happened in a later session. Every VAL row below
   must be re-recorded from this session's own evidence; the earlier
   focused-run tokens do not carry over.

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
