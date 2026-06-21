# Phase 11 - Task 1: Claim Matrix Rows And Domain Claim Semantics

## Objective

Create the internal claims domain layer and append Phase 11 claim action rows to
the in-code permissions matrix. This task proves the core write semantics before
public CLI wiring: acquire, conflict rejection, stale takeover, force-acquire,
renewal, preconditioned journaling/projection, and atomic concurrent acquire.

## Requirements

- PRD refs: PRD-CLAIM-001, PRD-CLAIM-002, PRD-CLAIM-003, PRD-CLAIM-004,
  PRD-CLAIM-008a, PRD-ERR-003
- TRD refs: TRD-ARCH-010, TRD-ARCH-011; `03_TRD_SCHEMA.md` sections 11.1 and 13
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid issue claim` semantics
- Quality refs: `08_QUALITY_PLAN.md` COV-001/COV-003; CF-007 validate-before-append
- Decisions: Phase 10 DEC-B (matrix grows with implemented actions); Phase 11
  DEC-026

## Validation Contract Coverage

- VAL-001: claim action rows exist in the in-code matrix and match schema section 13.
- VAL-002: normal acquire creates one live claim with default/custom TTL and note.
- VAL-003: live-claim conflict rejects before append and names holder/expiry.
- VAL-004: stale claim allows normal acquire takeover and records superseded holder.
- VAL-005: permitted force-acquire succeeds; denied force rejects before append.
- VAL-006: holder renewal updates time/TTL and projects an activity row.
- VAL-013: concurrent acquire is atomic with exactly one winner.

## Scope

### Files To Create Or Modify

- `internal/claims/` (new package) - domain service/types/errors/tests for issue
  claim acquire, renew, force-acquire, stale takeover, and helper queries.
- `internal/policy/matrix.go` and `internal/policy/matrix_test.go` - append the
  four `issue.claim.*` rows from `03_TRD_SCHEMA.md` section 13 and update
  implemented-action completeness tests.
- `internal/journal/*_test.go` only if claim event projection needs focused
  coverage beyond the claims package tests.

### Files To Read For Context

- `internal/issues/transition.go` and `internal/issues/update.go` for the local
  journal write/projection pattern.
- `internal/journal/project_effects.go` for insert/update/delete/activity
  effect semantics and precondition enforcement.
- `internal/policy/authorize.go`, `internal/policy/matrix.go`, and Phase 10 Task 2/3 specs.
- `migrations/0002_v01_tables.up.sql` and `migrations/0005_v01_ref_triggers.up.sql`
  for current `claims` table and v0.1 issue-only scope.

### Out Of Scope

- Public `fid issue claim` / `fid issue unclaim` command wiring (Task 2).
- `fid claim list` read surface (Task 3).
- Transition enforcement and auto-release (Task 4).
- Schema migrations, QA-review claims, named stale-claims reports, or composite
  `fid dev *` commands.

## Acceptance Criteria

- [ ] `internal/claims` owns claim business rules; `cmd/fid` does not contain
  claim state logic (VAL-002..006, VAL-013).
- [ ] The in-code matrix contains exactly one row each for
  `issue.claim.acquire`, `issue.claim.renew`, `issue.claim.force_acquire`, and
  `issue.claim.release`, matching `03_TRD_SCHEMA.md` section 13 (VAL-001).
- [ ] Normal acquire authorizes `issue.claim.acquire`, uses TTL 1800 seconds by
  default, accepts a custom TTL and note, and creates exactly one current issue
  claim row (VAL-002).
- [ ] A live non-expired conflict returns a domain error carrying holder and
  expiry, and dry-run/projection leaves journal and DB unchanged (VAL-003).
- [ ] Stale takeover succeeds under `issue.claim.acquire` without force
  permission, records superseded holder/expiry in the event/activity evidence,
  and makes the new holder current (VAL-004).
- [ ] `--force` semantics authorize `issue.claim.force_acquire`: a permitted role
  supersedes a live claim; a denied role gets `PermissionDeniedError` before
  append (VAL-005).
- [ ] Renewal authorizes `issue.claim.renew`, permits only the current live
  holder, preserves TTL when none is supplied, uses a supplied TTL when present,
  and writes an `activity_log` renewal effect (VAL-006).
- [ ] Concurrent acquire attempts against one issue have exactly one successful
  result and one current claim row (VAL-013).

## TDD Plan

- First failing test: normal acquire on an issue produces one claim row with
  `ttl_seconds = 1800`, expected holder, expiry, and one journal event.
- Additional tests: custom TTL/note; live conflict with zero append; stale
  takeover; permitted and denied force-acquire; holder and non-holder renewal;
  matrix completeness; concurrent acquire with exactly one winner.
- Regression risks: claim writes must use journal effects and preconditions, not
  direct DB mutation; rejected paths must validate before append; activity-log
  effect discipline must use `Op: "activity"` for renewal/takeover evidence.

## Documentation Updates

- `project-docs/docs/phases/phase-11/decisions.md` - record if implementation
  finds any mismatch between `03_TRD_SCHEMA.md` section 13 and the in-code rows.
- `project-docs/docs/phases/phase-11/tasks/task-1-done.md` - completion report
  with command results, coverage, matrix row summary, and commit hash/subject.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: policy show vs `Authorize` parity for claim rows
  starts here and is finalized by Task 2 CLI proof.
- Audit or release-readiness impact: TRD-ARCH-011 requires an integration-style
  concurrent acquire proof.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
```

The full gate above is required before task close because this task introduces a
new domain package, journaled writes, and the TRD-ARCH-011 concurrency proof. A
focused test command may be used during the red/green loop, but it does not
substitute for the full gate at close. Every required command must be recorded
in the completion report as `PASS`, or as `BLOCKED` with reason in a blocker
report.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
go test -bench . -benchmem ./...
```

- `go build -o /tmp/fid ./cmd/fid` - required if any command registration or
  CLI adapter file is touched unexpectedly; otherwise record N/A with reason.
- `go test -bench . -benchmem ./...` - required if benchmark fixtures or journal
  scale helpers are touched; otherwise record N/A with reason.

## Notes For Developer

- Use red/green TDD. Do not add a schema migration unless a spec-backed blocker
  is recorded first.
- Treat stale takeover as normal acquire permission (`issue.claim.acquire`), not
  force permission. Only explicit `--force` uses `issue.claim.force_acquire`.
- Keep claim logic in `internal/claims`; Task 4 can call helper functions from
  that package for transition integration.
- Completion report commit fields must be filled from `git rev-parse HEAD` and
  committed in the task commit or an immediate `docs(task): reconcile ...`
  follow-up. Do not leave working-tree-only commit fields.
