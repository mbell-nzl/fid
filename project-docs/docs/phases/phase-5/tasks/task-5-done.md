# Phase 5 - Task 5: `fid issue transition` with workflow + rule enforcement - Completion Report

## Summary of what was implemented

- Added `internal/issues.Transition`, the live transition service for
  `fid issue transition`, reusing the Task 1 validate-before-append seam and
  the Task 2 activity discipline.
- Wired Task 4 workflow graphs and rule functions to projection-backed reads so
  invalid edges, unmet AC verification/completion requirements, unmet qa_lead
  verification requirements, unmet epic child-task requirements, and cancelled
  re-transitions reject with rich errors before any journal append.
- Successful transitions now append exactly one `issue.transitioned` event with
  an `issues` status update effect guarded by `status` + `workflow_id`
  preconditions and one `activity` effect targeting `activity_log`.
- Added the `fid issue transition <id> --to <state> --by <role>` CLI subcommand
  as a thin adapter over the new domain service.
- This closes CF-005 for a second real activity emitter and reuses the CF-007
  seam so rejected transitions append zero journal lines.

## Files created or modified

- `internal/issues/transition.go`
- `internal/issues/transition_test.go`
- `cmd/fid/issue.go`
- `cmd/fid/issue_transition_test.go`
- `project-docs/docs/phases/phase-5/tasks/task-5-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - verified no change. The existing
  `fid issue transition` contract already matched the implemented flags and
  confirmation line.
- `project-docs/docs/phases/phase-5/tasks/task-5-done.md` - added this
  completion report with VAL evidence and gate results.

## Validation Contract Coverage

- `VAL-006` - PASS. `TestTransitionAllowsWorkflowEdgeProjectsActivityAndRebuilds`
  proves task `todo -> in_progress` transitions, projects activity, and rebuilds
  cleanly; `TestTransitionRequiresDoneOrCancelledChildrenToWrapEpic` proves an
  allowed epic transition succeeds and reads back `wrapped` / `started`;
  `TestTransitionSupportsSubtaskWorkflowAndCancellationRules` proves the
  subtask workflow advances to `in_progress` / `started`.
- `VAL-007` - PASS. `TestTransitionRejectsInvalidEdgeWithoutAppending` and
  `TestIssueTransitionRejectsInvalidEdgeAndLeavesStateUnchanged` prove a
  disallowed edge names the current/requested states and appends zero journal
  lines.
- `VAL-008` - PASS. `TestTransitionRequiresDeveloperVerificationsAndCompletionReport`
  proves `in_progress -> dev_complete` rejects until developer AC verifications
  and a completion report are present, names `all_acs_verified`,
  `has_completion_report`, and the AC IDs, then succeeds once fixtures are
  seeded through journal events.
- `VAL-009` - PASS. `TestTransitionRequiresQALeadVerificationForQAPass` proves
  `qa_review -> qa_pass` rejects until `qa_lead` verification exists and names
  the missing role plus AC ID, then succeeds.
- `VAL-010` - PASS. `TestTransitionRequiresDoneOrCancelledChildrenToWrapEpic`
  proves epic `in_progress -> wrapped` rejects with blocking child IDs and
  succeeds once all children are `done` or `cancelled`.
- `VAL-011` - PASS. `TestTransitionSupportsSubtaskWorkflowAndCancellationRules`
  proves cancellation succeeds from `todo` and `in_progress`, and rejects from
  `cancelled`.
- `VAL-005` - ADVANCED. The invalid-edge and unmet-rule transition tests above
  assert rejected transitions leave the canonical journal line count unchanged
  while later reads still succeed through `Get` and CLI `issue get`.

## Tests run

- `go test ./internal/issues -run TestTransitionAllowsWorkflowEdgeAndUpdatesReadModel -count=1`
  - FAIL before implementation (`undefined: Transition`, `undefined: TransitionOptions`)
- `go test ./internal/issues -run TestTransitionAllowsWorkflowEdgeAndUpdatesReadModel -count=1`
  - PASS
- `go test ./cmd/fid -run 'TestIssueTransition(ChangesStateAndPrintsConfirmation|RejectsInvalidEdgeAndLeavesStateUnchanged|RendersRuleFailures|RequiresTargetStateAndKnownByRole)$' -count=1`
  - FAIL before CLI wiring (`unknown flag: --to`)
- `go test ./cmd/fid -run 'TestIssueTransition(ChangesStateAndPrintsConfirmation|RejectsInvalidEdgeAndLeavesStateUnchanged|RendersRuleFailures|RequiresTargetStateAndKnownByRole)$' -count=1`
  - PASS
- `go test ./internal/issues ./cmd/fid`
  - PASS
- `go test -race -coverprofile=coverage.out ./...`
  - PASS (`cmd/fid` 79.0%, `internal/issues` 82.8%, `internal/workflow` 95.5%)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS
- `go build -o /tmp/fid ./cmd/fid` - not run; conditional command was not needed because evidence came from automated CLI tests rather than a manual smoke.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit creation.

## Commit subject

- Pending commit creation.
