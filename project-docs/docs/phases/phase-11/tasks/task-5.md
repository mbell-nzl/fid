# Phase 11 - Task 5: Claim Write Projection Freshness Lock Ordering

## Objective

Resolve ESC-008 by refactoring claim write paths so acquire, renew, and release
validate against a fresh SQLite projection without attempting to acquire
`.fid/write.lock` recursively when the journal is ahead of projection.

## Requirements

- PRD refs: PRD-OPS-003, PRD-NFR-008
- TRD refs: `02_TRD.md` §4.3 Write path, §4.9 Concurrency; TRD-ARCH-011
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid issue claim` and
  `fid issue unclaim`
- Quality refs: `08_QUALITY_PLAN.md` TDD-001, COV-001, COV-002, COV-003;
  DEC-028; ESC-008

## Validation Contract Coverage

- VAL-015: claim write paths catch up stale projections before validation
  without nested-lock failures and without changing public claim behavior.
- VAL-013: concurrent acquire atomicity remains green after the lock-ordering
  refactor.

## Scope

### Files To Create Or Modify

- `internal/journal/` - add a lock-aware projection freshness path or equivalent
  narrow refactor that avoids nested `journal.WithLock` while preserving normal
  `EnsureProjectionFresh` behavior for reads and callers without a held lock.
- `internal/claims/acquire.go`, `internal/claims/renew.go`,
  `internal/claims/release.go` - use the corrected projection freshness path
  before reading claim state.
- `internal/claims/claims_test.go` and, if needed,
  `internal/journal/*_test.go` - stale-projection and regression coverage.
- `cmd/fid/issue_claim_test.go` only if CLI-level stale-projection proof is
  needed to make VAL-015 deterministic.
- `project-docs/docs/phases/phase-11/tasks/task-5-done.md` - completion report.

### Files To Read For Context

- `internal/journal/replay.go`
- `internal/journal/lock.go`
- `internal/journal/writer.go`
- `internal/claims/acquire.go`
- `internal/claims/renew.go`
- `internal/claims/release.go`
- `internal/claims/claims_test.go`
- `project-docs/docs/phases/phase-11/tasks/task-4-done.md`
- `project-docs/docs/phases/phase-11/decisions.md` DEC-028
- `project-docs/docs/phases/phase-11/qa/validation-contract.md` VAL-013/015

### Out Of Scope

- Changing claim CLI flags, stdout/stderr wording, JSON shapes, or exit-code
  semantics.
- Adding or changing schema migrations.
- Normalizing every non-claim write path unless the minimal claim fix proves an
  unavoidable shared journal contract change.
- Changing workflow transition claim-gate precedence from Task 4.

## Acceptance Criteria

- [ ] Stale-projection acquire fixture proves a live claim present only in the
  journal is caught up before validation, so a non-holder acquire rejects with
  the current holder/expiry and appends nothing; the error must not be a nested
  `.fid/write.lock` timeout or held-lock error. (VAL-015)
- [ ] Stale-projection renew fixture proves the current holder can renew a claim
  that exists in the journal tail but not yet in SQLite, producing one renew
  event and one projected current claim row. (VAL-015)
- [ ] Stale-projection release fixture proves the current holder can release a
  claim that exists in the journal tail but not yet in SQLite, producing one
  release event and no current claim row. (VAL-015)
- [ ] Fresh-projection acquire, renew, release, stale takeover, force-acquire,
  and rejection behavior remains unchanged. (VAL-015)
- [ ] Concurrent acquire still produces exactly one winner, one current claim
  row, and holder/expiry loser errors. (VAL-013)
- [ ] No schema migration is added and no claim command public output contract
  changes. (VAL-015)

## TDD Plan

- First failing test: claim acquire on a journal-ahead-of-projection workspace
  currently fails with nested lock contention or validates against stale state;
  assert the expected holder conflict after catch-up.
- Additional tests: stale-projection renew, stale-projection release, existing
  concurrent acquire, and focused fresh-projection claim behavior regression.
- Regression risks: reintroducing the concurrent-acquire race Task 2 avoided;
  accidentally serving claim validation from stale SQLite; changing read-command
  `EnsureProjectionFresh`; widening the journal API without tests.

## Documentation Updates

- `project-docs/docs/phases/phase-11/tasks/task-5-done.md` - completion report
  with red/green evidence, VAL-013/015 coverage, exact command results, and
  commit hash/subject.
- `project-docs/docs/phases/phase-11/decisions.md` - update only if
  implementation deviates from DEC-028.
- `project-docs/specs/02_TRD.md`, `03_TRD_SCHEMA.md`, and
  `06_DESIGN_PLAN_CLI.md` - edit only if implementation proves the current
  specs wrong; otherwise record verified-no-change in the completion report.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: claim CLI public output should remain byte-for-
  byte compatible with Tasks 2-4 tests unless a documented bug forces
  escalation.
- Audit or release-readiness impact: final QA must explicitly verify ESC-008 is
  closed by VAL-015 evidence before phase-end review proceeds.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

The full gate above is required before task close because this task touches
journal/projection freshness and claim write concurrency. A focused test
command may be used during the red/green loop, but it does not substitute for
the full gate at close. Every required command must be recorded in the
completion report as `PASS`, or as `BLOCKED` with reason in a blocker report.

Conditional commands:

```text
N/A
```

- None.

## Notes For Developer

- Prefer a minimal lock-aware journal helper over moving claim freshness checks
  fully outside the write lock; the latter can validate against stale SQLite if
  another write appends after the pre-lock freshness check.
- Preserve `journal.EnsureProjectionFresh` behavior for existing read paths and
  callers that do not already hold a write lock.
- If the fix requires changing the documented TRD write path, stop and
  escalate instead of silently adapting the spec.
- Pair completion-report commit-hash filling with the task commit or immediate
  reconciliation commit.
