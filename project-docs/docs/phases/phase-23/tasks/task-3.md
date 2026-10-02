# Phase 23 - Task 3: Non-mutating review Claim mode and Reviewer Run

## Objective

Add the `review` Claim mode, keep exactly one active Claim per Issue across
both modes, and fence repository adoption so a review Claim can never authorize
a dirty-path recovery or abandonment.

## Validation Contract Coverage

- VAL-011: `claim.acquire` accepts both modes and still permits at most one
  active Claim per Issue.
- VAL-012: A review Claim cannot authorize repository adoption.
- VAL-013: A review Run starts, closes, abandons, and recovers under the
  existing Phase 22 Run rules.

## Binding Constraints

- `contracts.md` section 3 is the exact authority. The existing partial unique
  index on active Claims stays unchanged; mode-aware concurrency is Phase 24.
- `claim-mode-invalid` is a new closed failure code with the exact details in
  `public-interface.md` section 5. A rejected request creates no Event, Claim,
  Run, Stop, or Attention.
- fid records authority, never filesystem enforcement. The fence is on adoption
  of dirty paths in `run.recover` and `run.abandon`, not on file access.
- A review Run keeps every other Phase 22 Run rule, including executor identity
  uniqueness, liveness observation, budgets, closure derivation, and recovery.
- Do not add a review-specific Run status, a second active-Run rule, or any
  Review content; Review recording is task 4.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/claims/generic.go` and its package files.
- `internal/workflow/recovery.go`, `run.go`, `run_read.go`, `authority.go`.
- `internal/cli/codes.go` - the one new code.
- `internal/journal/current_event_value_catalog.go` - widen
  `currentValueShapes["Claim"]`'s `mode:=mutation` to `mode:=mutation,review`
  (R-16, DEC-231). This is reachable in this task and not deferrable: a
  review-mode Claim carries a Run, and `run.abandoned` and `run.recovered` both
  project the `Claim` value through `currentProjectedValue`, so without the
  widening the first abandon or recover of a review-mode Run is rejected at
  append time with its cause outside the failing task. `public-interface.md`
  section 2 fixes the amended domain; add nothing else to this file.
- `cmd/fid/claim.go`, `run.go`, `run_recovery.go` - request and result only.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-3-done.md`.

### Files To Read For Context

- `contracts.md` section 3; `public-interface.md` sections 2, 4, 5;
  `schema.md` section 1.
- `project-docs/docs/phases/phase-22/contracts.md` sections 2 and 4 for the
  retained lease and recovery rules.
- `task-2-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Review Records, Findings, and every judgment fact.
- Parallel or cross-Issue Claims and multi-revision Envelopes.
- Any change to the active-Claim uniqueness index.

## Acceptance Criteria

- [ ] Competing acquires in either mode reject with the existing conflict code,
  and post-release acquires succeed. (VAL-011)
- [ ] Review-mode `run.recover` and `run.abandon` with adopted paths reject
  `claim-mode-invalid` with no Event, while mutating-mode requests are
  unchanged. (VAL-012)
- [ ] A review Run completes the full lifecycle identically to a mutating Run
  except adoption. (VAL-013)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestClaimAcquireAcceptsReviewMode` fails against the current mutation-only
mode validator.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-3-done.md` in the full format.
- Flag spec impact; only Tech Lead accepts a contract or schema change.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -timeout 30m ./cmd/fid/integration/...
```

Two gates run for minutes and are explicitly bounded (DEC-227): the `-timeout
30m` above replaces Go's implicit 10-minute per-binary default, which the tagged
suite already approaches at 560 s. Record in the completion report the wall time
`go test` reports for `cmd/fid/integration`, `cmd/fid`, and `internal/workflow`,
so a regression is attributed to the task that caused it. Budgets: tagged
integration under 900 s and no ordinary package over 600 s (VAL-079). A breach is
a bounded speed repair owned by this task when this task's own new journeys
caused it; otherwise stop and escalate. Never record a timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. No exception to the full task
gate. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination. An
interrupted or hung command is never `PASS`; preserve all failed attempts and
rerun the byte-identical command after an authorized causal repair. Report lint
is a Runner control after the report commit, not a Developer gate.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- This task belongs to QA Review window 1.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-3`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
