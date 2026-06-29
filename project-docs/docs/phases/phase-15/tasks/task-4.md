# Phase 15 - Task 4: End-to-end composite smoke + dogfood prep

## Objective

Prove the `fid dev *` composites work end-to-end over a real `.fid` workspace
and prepare the human dogfood. One integration/e2e test drives both composite
paths, asserts one event per composite with the correct effects, enforces
permission denial, and checks delete-and-rebuild determinism. This task adds no
new product behavior.

## Requirements

- PRD refs: PRD-DEV-004, 008, 009, 010, PRD-CLAIM-005
- TRD refs: `02_TRD.md` §4.3 (journal-canonical, rebuild); ADR-0001
- Design refs: `07_DESIGN_PLAN_FLOWS.md` §1.1 (developer execution flow)
- Quality refs: `08_QUALITY_PLAN.md` §5.2 (phase gate), §5.6 (human checklist)

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-15/qa/validation-contract.md`.

- VAL-011: e2e smoke — `dev start → gate → finish` and `dev start → block` over
  a real workspace; one event per composite; permission denial enforced; rebuild
  deterministic.
- VAL-012 is phase-end (Human); this task only prepares the workspace/run-sheet
  inputs QA Lead needs (no human run in this task).

## Scope

### Files To Create Or Modify

- `cmd/fid/e2e_v01e_dev_test.go` (new) - integration-tagged end-to-end test:
  1. `fid init` a temp workspace; create an epic + a task with an acceptance
     criterion.
  2. Path A: `dev start` (assert claim live, issue `in_progress`, run open, one
     journal line) → `ac verify --by developer --status pass` → `dev gate`
     (assert one `gate_runs` row, attempt 1) → `dev finish` (assert run closed,
     completion `dev_complete` with gate snapshot, issue `dev_complete`, claim
     released, one journal line).
  3. Path B: a second task → `dev start` → `dev block` (assert completion
     `failed_gate` with reason, issue `failed_gate`, claim released, one journal
     line).
  4. Permission denial: a composite invoked by `qa_reviewer` is rejected before
     append (non-zero exit, zero new journal lines).
  5. Determinism: count journal lines, `fid rebuild`, and assert the projected
     state (issue states, claims, completions, runs, gate_runs) is unchanged.
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` - only if the implemented
  composite output/flow diverges from the documented §1.1 example; otherwise
  leave unchanged and record `verified-no-change` in the completion report.

### Files To Read For Context

- `cmd/fid/e2e_v01d_closure_test.go` (existing end-to-end smoke pattern:
  init → create → claim → transition → complete → unclaim with rebuild check)
- `internal/dev/start.go`, `gate.go`, `finish.go`, `block.go` (Tasks 2-3)
- `internal/journal/projection.go` / `replay.go` (rebuild entry points after the
  Task 1 split)

### Out Of Scope

- New composite behavior (Tasks 2-3 own it).
- Activity timeline reads, maintenance audits, stub guard (Phases 16-17).
- The human run itself (VAL-012 is executed at phase-end by QA Lead + human).

## Acceptance Criteria

- [ ] An integration-tagged e2e test drives `dev start → gate → finish` and
  `dev start → block` over a real `.fid` workspace, asserting one journal event
  per composite with the correct claim/transition/run/gate/completion effects. (VAL-011)
- [ ] The test asserts a non-developer role is rejected before append on a
  composite (non-zero exit, zero new journal lines). (VAL-011)
- [ ] The test asserts `fid rebuild` reproduces the projected state after the
  composite lifecycle (delete-and-rebuild determinism). (VAL-011)

## TDD Plan

- First failing test: the e2e scenario fails red until the full composite
  lifecycle and rebuild equality hold.
- Additional tests: permission-denial branch; block path; rebuild equality.
- Regression risks: test flakiness from timestamps/ordering — assert on
  projected state and journal-line counts, not wall-clock values; keep the test
  hermetic (temp workspace, no network).

## Documentation Updates

- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` - reconcile §1.1 only if the
  implementation diverges; otherwise record `verified-no-change`.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None (uses a temp workspace; do not read held-out
  scenario material).
- Cross-surface parity impact: the smoke exercises the composite paths against
  the same projection the single-command paths use.
- Audit or release-readiness impact: this is the v0.1e slice-2 exit smoke; the
  full v0.1 exit-criteria smoke remains Phase 17.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

The full task gate above is required before task close, **including the
integration gate and build** (this task adds an integration-tagged e2e test).
Record each command in the completion report as `PASS`, or `BLOCKED` with
reason. Record the `cmd/fid` coverage percentage.

Conditional commands:

```text
go test -bench . -benchmem ./...
```

List any conditional command required by this task and why:

- `go test -bench . -benchmem ./...` — only if a benchmark is added or touched;
  otherwise not required at task level (it is part of the phase gate).

## Notes For Developer

- Model the test on `cmd/fid/e2e_v01d_closure_test.go`; reuse its
  init/create/rebuild scaffolding.
- Assert observable projected state and journal-line deltas, not internal
  function calls.
- This task is the phase exit smoke for the composites; if any composite cannot
  be driven end-to-end, stop and escalate rather than weakening the assertions.
