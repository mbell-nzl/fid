# Phase 5 - Task 5: `fid issue transition` with workflow + rule enforcement

## Objective

Implement `fid issue transition <id> --to <state> --by <role>`, wiring the
Task 4 `internal/workflow` package to the live projection. The command validates
the requested edge against the issue's default workflow, evaluates the edge's
transition rules against current DB state, rejects invalid edges and unmet rules
with rich errors (naming states / unmet conditions / entity IDs), and on success
journals one event whose effects update `issues.status` (with status +
`workflow_id` preconditions) and emit an `activity` effect projected to
`activity_log`. Routes through the Task 1 validate-before-append seam so rejected
transitions append zero journal lines.

Transition-rule evidence may seed AC/verification/completion/concern/child rows
through **journal fixtures** (the v0.1c write commands have not shipped).

## Requirements

- PRD refs: PRD-WF-001 (follow allowed edges), PRD-WF-002 (reject invalid edge,
  name states), PRD-WF-003 (evaluate rules before transition), PRD-WF-004
  (report unmet conditions + entity IDs), PRD-WF-007/008/010 (rules),
  PRD-WF-013 (cancel from any non-cancelled state)
- TRD refs: TRD §4.3 (transition update effect with status + workflow_id
  preconditions; `activity` effect), TRD-ARCH-004, TRD-ARCH-005, TRD-ARCH-009
- Design refs: `06_DESIGN_PLAN_CLI.md` §`fid issue transition`
- Quality refs: `08_QUALITY_PLAN.md` §4.3 COV-001 (`internal/issues`,
  `internal/workflow` ≥ 80%), COV-002 (`cmd/fid` ≥ 70%)

## Validation Contract Coverage

- VAL-006: allowed edge transitions; read-back shows new status + category.
- VAL-007: invalid edge rejected naming current + requested state; zero journal
  lines.
- VAL-008: `in_progress → dev_complete` enforces `all_acs_verified{developer}` +
  `has_completion_report`; unmet names conditions + AC IDs.
- VAL-009: `qa_review → qa_pass` enforces `all_acs_verified{qa_lead}`.
- VAL-010: epic `in_progress → wrapped` enforces
  `all_child_tasks_done_or_cancelled`; unmet names blocking child IDs.
- VAL-011: cancel allowed from non-cancelled states; rejected from cancelled.
- VAL-005 (advances): rejected transition appends zero journal lines.

## Scope

### Files To Create Or Modify

- `internal/issues/transition.go` (new) - `Transition` domain service: load
  current issue + workflow, gather rule inputs from the projection
  (verifications, completions, open concerns, child task states), call the
  Task 4 workflow edge/rule checks, build the status update effect (with
  `status` + `workflow_id` preconditions) and the `activity` effect, route
  through the validate-before-append seam.
- `internal/issues/transition_test.go` (new) - package tests with journal
  fixtures seeding ACs/verifications/completions/concerns/child tasks.
- `cmd/fid/issue.go` - wire the `transition` subcommand (thin adapter).
- `cmd/fid/issue_test.go` (or new `issue_transition_test.go`) - CLI tests.

### Files To Read For Context

- `internal/workflow/workflow.go`, `internal/workflow/rules.go` (Task 4)
- `internal/issues/issue.go`, `internal/issues/detail.go` (concern/completion
  readers already exist)
- `internal/journal/project.go` (preconditions, activity op, dry-run seam)
- `06_DESIGN_PLAN_CLI.md` §`fid issue transition`

### Out Of Scope

- Claim-gated transitions (started-category requires live claim) — v0.1d
  (PRD-CLAIM-005/006).
- Composite developer commands (`fid dev start/finish/block`) — v0.1e.
- `fid issue complete` / AC / concern write commands — v0.1c.
- Configurable workflows — v0.2.

## Acceptance Criteria

- [ ] An allowed edge transitions the issue; `fid issue get --format json` reads
  back the new status and correct status category for task/epic/subtask.
  (VAL-006)
- [ ] A disallowed edge is rejected with a non-zero error naming current and
  requested state and appends zero journal lines. (VAL-007, VAL-005)
- [ ] `in_progress → dev_complete` is rejected when ACs are not developer-verified
  or no completion report exists, naming the unmet rule(s) and the failing AC
  IDs; it succeeds when both rules are satisfied (seeded via journal fixtures).
  (VAL-008)
- [ ] `qa_review → qa_pass` enforces `all_acs_verified{qa_lead}` (both states).
  (VAL-009)
- [ ] Epic `in_progress → wrapped` is rejected when a child task is not
  done/cancelled (naming blocking child IDs) and succeeds when all children are
  done/cancelled. (VAL-010)
- [ ] Any non-cancelled state transitions to `cancelled`; cancelled→cancelled is
  rejected. (VAL-011)
- [ ] One successful transition = one journal event with a status `update`
  effect (status + workflow_id preconditions) and one `activity` effect;
  `fid rebuild` reproduces identical state. (VAL-006, supports VAL-002)
- [ ] `--to` and `--by` required; `--by` validated against known roles.
- [ ] COV-001 (`internal/issues`, `internal/workflow`) and COV-002 (`cmd/fid`)
  hold.

## TDD Plan

- First failing test: `internal/issues.Transition` on `todo → in_progress`
  succeeds and read-back shows `in_progress`/`started`. Fails (no Transition).
- Additional tests: each VAL edge above, both states; rich-error content
  assertions (state names, unmet rule names, AC/child IDs present); zero-journal
  on rejection; one-event-per-transition + rebuild parity; precondition mismatch
  surfaces a clear error.
- Regression risks: status-category consistency with `internal/issues` read
  model; ensure rule-input gathering reads the same projection the dry-run uses;
  ensure cancel edge bypasses non-existent rule lookups cleanly.

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` - confirm `fid issue transition` documented contract
  (flags, confirmation line, rejection errors) matches the implementation.
- Completion report: cite VAL-006..011; note CF-005 closes (transition is the
  second real activity emitter) and confirm CF-007 seam reused.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: status-category values shared with `fid issue
  get`; rich-error format should be consistent with existing CLI error style.
- Audit or release-readiness impact: Completes the Phase 5 mutation+workflow
  exit criterion; with Task 3 it closes CF-005.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command as
`PASS` (or `BLOCKED` with reason) in the completion report.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

- `go build` if a manual CLI smoke of `fid issue transition` is used as evidence
  for VAL-006/007/008/010/011.

## Notes For Developer

- Build the status update effect with both `status` and `workflow_id`
  preconditions per TRD §4.3 so a merged/diverged journal surfaces a clean
  replay conflict.
- Rich errors must satisfy PRD-WF-004: name the unmet condition(s) and the
  entity IDs involved (AC IDs for `all_acs_verified`, child task IDs for
  `all_child_tasks_done_or_cancelled`).
- Use journal fixtures to seed prerequisite rows (ACs/verifications/completions/
  concerns/child tasks); do not add v0.1c write commands in this task.
- Reuse the Task 1 validate-before-append seam: rejected transitions must append
  zero journal lines (VAL-005).
