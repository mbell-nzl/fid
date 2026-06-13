## Summary of what was implemented

- Added the library-only ready-selection engine in `internal/issues` via `Ready`, `ReadyOptions`, `ReadyResult`, `ReadyTask`, and `ReadyExplanation`.
- Implemented DEC-017 readiness semantics: actionable = status category `unstarted`; unfinished = not `completed` or `cancelled`.
- Added deterministic selection for sequential epics, prioritised epics, unattached tasks, and project-wide selection, returning the frozen §6.3 task field set plus structured empty/blocker explanations.
- Added exhaustive deterministic unit coverage for VAL-002..006, including sequential blocked behavior and an explicit project-wide tie-break test.

## Files created or modified

- `internal/issues/ready.go`
- `internal/issues/ready_test.go`
- `project-docs/docs/phases/phase-6/tasks/task-3-done.md`

## Documentation updated

- Added this completion report at `project-docs/docs/phases/phase-6/tasks/task-3-done.md`.
- No glossary update was needed.

## Validation Contract Coverage

- `VAL-002` PASS: `TestReadySelectsLowestSequenceActionableTaskForSequentialEpic` proves a sequential epic returns the lowest-sequence unfinished actionable task and not a higher-sequence task.
- `VAL-003` PASS: `TestReadyReturnsStructuredBlockerWhenSequentialCurrentTaskIsStarted` proves a sequential epic with a started current task returns no task plus a structured blocker naming the current task.
- `VAL-004` PASS: `TestReadySelectsHighestPriorityActionableTaskForPrioritisedEpic` proves a prioritised epic selects the highest-priority actionable task and excludes a started higher-priority task.
- `VAL-005` PASS: `TestReadyProjectWideIncludesOnlyActionableUnattachedTasks` proves an unattached actionable task is eligible project-wide and started/completed/cancelled unattached tasks are excluded.
- `VAL-006` PASS: `TestReadyProjectWideSelectsDeterministicTopTaskAcrossMixedSources` proves project-wide selection returns one deterministic winner across prioritised epics and unattached tasks while sequential epics contribute only their current actionable task. `TestReadyProjectWideUsesIDAsFinalTieBreak` proves the final deterministic tie-break.
- `VAL-001` support PASS: `TestReadySelectsLowestSequenceActionableTaskForSequentialEpic` asserts the engine returns the frozen §6.3 task fields: `id`, `title`, `status`, `epic_id`, `epic_title`, `sequence`, `priority`, `est_tokens`.

## Tests run

- `go test ./internal/issues -run TestReadySelectsLowestSequenceActionableTaskForSequentialEpic`
  Expected red result before implementation: build failed with `undefined: Ready` and `undefined: ReadyOptions`.
- `go test ./internal/issues -run TestReadySelectsLowestSequenceActionableTaskForSequentialEpic`
  PASS after initial implementation.
- `go test ./internal/issues -run '^TestReady'`
  PASS.
- `go test -race -coverprofile=coverage.out ./...`
  PASS.

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS (no output).
- `golangci-lint run ./...` — PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` — PASS. Key package coverage: `cmd/fid` 79.0%, `internal/issues` 82.5%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5%.
- `go tool cover -func=coverage.out` — PASS. Total coverage 77.8%; `internal/issues/ready.go` functions covered, and `internal/issues` remains above the Phase 6 `COV-001` floor.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- Deterministic project-wide tie-break chosen for Phase 6 Task 3: `priority ASC` (nil last), then `created_at ASC`, then `id ASC`. The PRD requires deterministic selection but does not pin the exact tie-break order, so this choice needs TL confirmation when Task 4 wires the CLI surface.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit.

## Commit subject

- Pending commit.
