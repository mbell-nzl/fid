# Phase 18 - Task 3: Transition-path data-source swap (behavior-preserving) - Completion Report

## Summary of what was implemented

- Repointed `issues.Transition` to load the DB-backed workflow registry once through the existing locked database connection and resolve edge legality, ordered rules, and started-category checks from it.
- Added an explicit registry edge index so the transition path consumes the seeded cancellation rows directly; the frozen Go definitions and all non-transition consumers remain unchanged.
- Preserved the existing `issue.transitioned` effects, status/workflow preconditions, activity effect, v0.1 invalid/blocked error text, started-claim enforcement, and same-event release effect.

## Files created or modified

- `internal/issues/transition.go`
- `internal/issues/transition_test.go`
- `internal/workflow/registry.go`
- `project-docs/docs/phases/phase-18/tasks/task-3-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-18/tasks/task-3-done.md`

## Validation Contract Coverage

- VAL-008: PASS — `TestTransitionRegistryMatchesFrozenV01OrderedStateMatrix` exhaustively compares every ordered state pair across `task-default` (81), `epic-default` (25), and `subtask-default` (16) against the frozen v0.1 edge/rule reference, including explicit cancel edges. Existing transition and CLI tests cover the categories; isolated built-CLI runs recorded exact allowed, invalid-edge, blocked-rule, cancel, and re-cancel output. The developer-rule failure assertion also locks the full v0.1 `BlockedTransitionError` text.
- VAL-009: PASS — existing missing, wrong-holder, stale, and live-caller transition cases remain green; `TestTransitionUsesDBCategoriesForClaimGating` proves started-category claim gating reads the loaded registry category rather than the Go map.
- VAL-020: PASS — `TestTransitionOutOfStartedReleasesCallerClaimInSameEvent` verifies the transition event includes the claims delete effect and the projected claim is absent. Its existing rebuild check remains covered by the transition suite.
- QA-13 / EVID-08: PASS — exhaustive 122-pair registry matrix, DB-controlled edge/rule/cancel source tests, exact errors, and CLI transition coverage passed.
- QA-14 / QA-15 / EVID-09: PASS — started-entry claim matrix and same-event out-of-started release test passed.

## Tests run

- `go test ./internal/issues -run '^TestTransitionRejectsEdgeMissingFromDBRegistry$' -count=1` - expected RED: the Go-map path allowed the DB-removed edge.
- `go test ./internal/issues -run '^TestTransitionRejectsEdgeMissingFromDBRegistry$' -count=1` - PASS after the registry swap.
- `go test ./internal/issues -run '^(TestTransitionRejectsEdgeMissingFromDBRegistry|TestTransitionRegistryMatchesFrozenV01OrderedStateMatrix|TestTransitionUsesDBRuleReferences|TestTransitionUsesDBCategoriesForClaimGating|TestTransitionRejectsCancelledEdgeMissingFromDBRegistry|TestTransitionRequiresDeveloperVerificationsAndCompletionReport)$' -count=1` - PASS.
- `go test ./internal/issues -count=1` - PASS.
- `go test ./internal/workflow -count=1` - PASS.
- Runtime entrypoint exercise: built `/tmp/fid-phase18-task3` and ran `fid issue transition` in isolated workspaces - PASS. Observed exact category output: allowed `transitioned <id>: todo → in_progress`; invalid `error: transition issue "<id>" from "todo" to "done" is not allowed`; blocked `error: transition issue "<id>" from "in_progress" to "dev_complete" blocked by workflow rules:` followed by the `all_acs_verified` (AC ID, missing `developer`) and `has_completion_report` (issue ID) lines; cancel `transitioned <id>: todo → cancelled`; and re-cancel `error: transition issue "<id>" from "cancelled" to "cancelled" is not allowed`.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS (exit 0; no files listed).
- `golangci-lint run ./...` - PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` - PASS (exit 0; coverage profile produced).
- `go tool cover -func=coverage.out` - PASS (exit 0; total statement coverage: 76.3%).
- `go test -tags integration -race ./...` - PASS (exit 0).

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `cf68153de4058f77d06451a1ceed3f2191658a77`

## Commit subject

- `feat(issues): resolve transitions from workflow registry`
