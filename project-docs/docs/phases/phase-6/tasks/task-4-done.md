## Summary of what was implemented

- Added the `fid ready` Cobra command with `--epic` and `--format`, and registered it in the root CLI.
- Wired the command to the Phase 6 Task 3 ready-selection engine without changing selection logic, so DEC-019's pinned ordering continues to come from `internal/issues/ready.go`.
- Implemented the DEC-018 empty-state contract: exit `0`, stdout payload with exactly the frozen `task` key, and a `note:` explanation on stderr only.
- Routed `fid ready` through the existing format engine for `json`, `toon`, `markdown`, and `yaml`, including empty-state rendering for each format.

## Files created or modified

- Created `cmd/fid/ready.go`.
- Created `cmd/fid/ready_test.go`.
- Modified `cmd/fid/root.go`.
- Created `internal/format/ready.go`.
- Modified `internal/format/markdown.go`.
- Modified `project-docs/specs/06_DESIGN_PLAN_CLI.md`.
- Created `project-docs/docs/phases/phase-6/tasks/task-4-done.md`.

## Documentation updated

- Updated `project-docs/specs/06_DESIGN_PLAN_CLI.md` §`fid ready` to make the v0.1b surface explicit: only `--epic` and `--format` are implemented now, the empty state exits `0`, stdout carries `{"task": null}`, and the explanatory `note:` is stderr-only.

## Validation Contract Coverage

- `VAL-001` PASS via `TestReadyDefaultsToFrozenJSONShapeWithoutNote`, which asserts bare `fid ready` returns one selected task in the frozen `{"task": {...}}` §6.3 JSON shape.
- `VAL-002` PASS via `TestReadyEpicScopesSelectionAndBlockerStates/sequential_actionable_current_task`, which proves `fid ready --epic <id>` returns the lowest-sequence unfinished actionable task for a sequential epic.
- `VAL-003` PASS via `TestReadyEpicScopesSelectionAndBlockerStates/sequential_started_current_task_returns_empty_note`, which proves a sequential epic whose current task is already started returns no task plus a stderr `note:` naming the blocking task and status.
- `VAL-004` PASS via `TestReadyEpicScopesSelectionAndBlockerStates/prioritised_epic_uses_highest_priority_actionable_task`, which proves `fid ready --epic <id>` uses the highest-priority actionable task for a prioritised epic.
- `VAL-005` PASS via `TestReadyProjectWideUsesMixedSourceSelectionOrder`, whose mixed-source fixture includes unattached ready tasks and excludes unattached non-selected cases while proving project-wide eligibility.
- `VAL-006` PASS via `TestReadyProjectWideUsesMixedSourceSelectionOrder`, which proves project-wide `fid ready` returns one deterministic winner across sequential, prioritised, and unattached candidates while sequential epics contribute only their current task. The fixture also covers the QA review improvement by proving a nil-priority unattached task loses to a prioritised candidate.
- `VAL-007` PASS via `TestReadyDefaultsToFrozenJSONShapeWithoutNote`, `TestReadyEpicScopesSelectionAndBlockerStates/sequential_started_current_task_returns_empty_note`, and `TestReadyFormatsRenderSelectedAndEmptyStates`, which together prove non-empty output emits no `note:` and empty output exits successfully with a null payload plus stderr `note:`.
- `VAL-008` PASS via `TestReadyFormatsRenderSelectedAndEmptyStates`, which proves `json`, `toon`, `markdown`, and `yaml` render both selected and empty states, and via `TestReadyDefaultsToFrozenJSONShapeWithoutNote`, which proves omitting `--format` yields the JSON default.

## Tests run

- `go test ./cmd/fid -run '^TestReady'`
  Expected red result before implementation: `unknown command "ready" for "fid"`.
- `go test ./cmd/fid -run '^TestReady'`
  PASS after adding the command, formatter payload, and markdown support.
- `gofmt -l cmd internal`
  PASS.
- `golangci-lint run ./...`
  PASS.
- `go test -race -coverprofile=coverage.out ./...`
  PASS. `cmd/fid` coverage: `78.3%`.
- `go tool cover -func=coverage.out`
  PASS. Total coverage: `76.7%`.
- `go build -o /tmp/fid ./cmd/fid`
  PASS.

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
  Output: `0 issues.`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
  `cmd/fid`: `78.3%`
  `internal/issues`: `82.5%`
  `internal/journal`: `80.3%`
  `internal/policy`: `85.0%`
  `internal/workflow`: `95.5%`
- `PASS` `go tool cover -func=coverage.out`
  Total coverage: `76.7%`
- `PASS` `go build -o /tmp/fid ./cmd/fid`

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending.

## Commit subject

- Pending.
