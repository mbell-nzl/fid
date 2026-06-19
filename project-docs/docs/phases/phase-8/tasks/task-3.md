# Phase 8 - Task 3: `has_completion_report` real-data gating + rebuild parity

## Objective

Prove the Phase 8 write surface drives the existing `has_completion_report`
transition rule with **real** `fid issue complete` data (replacing the Phase 5
journal fixtures), and that completions + gate snapshots survive a delete-DB +
`fid rebuild` with identical projected state. This is the Phase 8 end-to-end
closure task.

## Requirements

- PRD refs: PRD-DEV-001/010 (completions), PRD-WF-008 (`has_completion_report`)
- TRD refs: TRD_SCHEMA §2; ADR-0001 (journal canonical, DB derived)
- Design refs: `06_DESIGN_PLAN_CLI.md` §2 (`issue complete`, `issue transition`)
- Quality refs: `08_QUALITY_PLAN.md` integration/rebuild-parity; DEC-007 coverage

## Validation Contract Coverage

- VAL-006: A real `has_completion_report`-gated transition rejects when no
  completion exists (naming the issue, zero state change) and succeeds once a
  real `fid issue complete` row exists.
- VAL-007: Delete DB + `fid rebuild` reproduces identical `task_completions`
  (attempt, superseded status, gate snapshot, counts, human-intervention) and
  `gate_runs` state.

## Scope

### Files To Create Or Modify

- `cmd/fid/e2e_v01c_complete_test.go` (new) - end-to-end test over create →
  transition-reject (no completion) → `fid issue complete` → transition-pass →
  delete-DB → `fid rebuild` → parity. Mirror `e2e_v01c_ac_test.go` /
  `e2e_v01b_test.go` structure.
- Do **not** modify the rule functions or the Task 1/2 write path except to fix a
  defect this task surfaces (raise via triage if the default workflow does not
  actually require `has_completion_report` on the transition under test).

### Files To Read For Context

- `internal/issues/transition.go`, `internal/issues/transition_rules.go`
  (`loadCompletionReportStatus` + `has_completion_report` rule — already wired to
  count `task_completions`).
- `internal/workflow/rules.go` (`HasCompletionReport`).
- `cmd/fid/e2e_v01c_ac_test.go` (Phase 7 closure pattern + rebuild-parity asserts).
- Phase 5 `task-4`/`task-5` notes on the fixture-based `has_completion_report`
  evidence being replaced.

### Out Of Scope

- `all_acs_verified` (Phase 7) and `concerns_addressed` (Phase 9) gating.
- New ranking/selection or workflow-graph changes.
- `internal/issues/issue.go` size split: Task 4.
- Held-out scenario harness relocation (separate scheduled commitment; recommend
  only).

## Acceptance Criteria

- [ ] End-to-end: create a task whose workflow transition requires
  `has_completion_report`, attempt that transition with no completion → **rejected**
  with an error naming the unmet issue, and the issue status is unchanged (VAL-006).
- [ ] After `fid issue complete` records a real completion via the real command,
  the same transition **succeeds** (VAL-006).
- [ ] The reject→pass evidence uses real `fid issue complete` / `fid issue
  transition` commands, not journal fixtures (VAL-006).
- [ ] The rejected transition appends **no** journal event (validate-before-append
  holds; assert journal line count unchanged across the rejected attempt).
- [ ] After exercising completions (including a second attempt → supersession and a
  fixture-seeded gate snapshot), deleting the SQLite DB and running `fid rebuild`
  reproduces identical `task_completions` (attempt_number, `superseded` status,
  `gate_results`, `gates_passed`/`gates_failed`, `human_intervention`) and
  `gate_runs` rows (VAL-007).

## TDD Plan

- First failing test: gated transition before any completion is rejected with the
  unmet issue named and zero state change.
- Additional tests: transition succeeds after a real completion; rejected
  transition adds no journal line; post-rebuild completion + gate parity.
- Regression risks: confirm the default task workflow's transition under test
  actually lists `has_completion_report`; reuse the fixture-seeded gate snapshot
  from Task 2 in the rebuild-parity assertion so the snapshot is exercised through
  rebuild.

## Documentation Updates

- `project-docs/docs/phases/phase-8/` completion report - paste full end-to-end
  transcript (reject, complete, pass, rebuild parity).
- Note in the completion report that the Phase 5 fixture-based
  `has_completion_report` gating evidence is now superseded by real commands
  (spec-sync input for the TL wrap audit).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: this smoke may extend the held-out scenario library;
  recommend placement only — do not paste hidden scenario content.
- Cross-surface parity impact: completion + gate state parity across `fid issue
  get` JSON and the transition-rule loader before vs after rebuild.
- Audit or release-readiness impact: confirms Phase 5 fixture evidence is replaced
  by real write commands (spec-sync input for wrap).

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Coverage targets (DEC-007): `cmd/fid` ≥ 70%, `internal/issues`/`internal/workflow` ≥ 80%.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` if the end-to-end test is behind the
  `integration` build tag (match the existing e2e test's tagging convention).
- `go build -o /tmp/fid ./cmd/fid` for pasting real CLI transcripts.

## Notes For Developer

- The `has_completion_report` rule and its loader already exist (count of
  `task_completions` for the issue). This task proves it with real `fid issue
  complete` data and replaces the Phase 5 fixture evidence; do not change the rule
  unless a defect requires it (raise via triage first).
- The gated transition must reject **without** appending a journal event
  (CF-007 / VAL-006); assert journal line count is unchanged across the rejected
  attempt.
- If the default task workflow does not require `has_completion_report` on the
  transition under test, stop and raise it (triage) rather than altering the
  workflow graph to make the test pass.
