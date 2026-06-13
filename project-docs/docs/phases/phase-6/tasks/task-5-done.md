## Summary of what was implemented

- Added the first integration-tagged end-to-end v0.1b smoke in `cmd/fid/e2e_v01b_test.go`.
- The smoke runs a real temp `.fid/` workspace through `fid init`, `fid issue create`, `fid issue list`, `fid issue get`, `fid issue update`, `fid issue transition`, and `fid ready`, with read-back assertions at each stage.
- The smoke proves post-rebuild parity by deleting `.fid/fid.sqlite`, running `fid rebuild`, and comparing the resulting `issue get`, `issue list`, and `ready` outputs against the final pre-rebuild state.
- Added a second integration-tagged end-to-end check for the empty/all-finished `fid ready` contract using only public commands.
- Added `project-docs/docs/phases/phase-6/held-out-recommendation.md` with a recommend/defer judgment and a concrete isolation mechanism for a future held-out harness.

## Files created or modified

- Created `cmd/fid/e2e_v01b_test.go`.
- Created `project-docs/docs/phases/phase-6/held-out-recommendation.md`.
- Created `project-docs/docs/phases/phase-6/tasks/task-5-done.md`.

## Documentation updated

- Added `project-docs/docs/phases/phase-6/held-out-recommendation.md` to recommend promoting the new smoke to the first held-out scenario seed once the harness location is formally accepted, while keeping the scenario material outside the implementation-agent workspace.

## Validation Contract Coverage

- `VAL-011` PASS via `TestV01BSmokeCreateListGetUpdateTransitionReadyAndRebuildParity`, which proves create -> list -> get -> update -> transition -> ready read-back consistency in a real temp workspace and proves delete-DB + `fid rebuild` parity for the touched issues plus ready/list outputs.
- `VAL-011` PASS via `TestV01BReadyReturnsEmptyAfterAllTasksFinishedEndToEnd`, which exercises the empty/all-finished `fid ready` contract end-to-end through public commands.
- `VAL-012` support only: the new smoke creates deterministic sequential and prioritised epic fixtures that can seed later human dogfood, but this task does not claim the human-only validation row.

## Tests run

- `go test -tags integration ./cmd/fid -run 'TestV01B(SmokeCreateListGetUpdateTransitionReadyAndRebuildParity|ReadyReturnsEmptyAfterAllTasksFinishedEndToEnd)$'`
  Red: failed first on the new smoke assertions while tightening the expected public output and parity snapshot timing.
- `go test -tags integration ./cmd/fid -run 'TestV01B(SmokeCreateListGetUpdateTransitionReadyAndRebuildParity|ReadyReturnsEmptyAfterAllTasksFinishedEndToEnd)$'`
  PASS after fixing the new smoke assertions and final-state parity capture.
- `gofmt -w cmd/fid/e2e_v01b_test.go`
  PASS.
- `go test -tags integration ./cmd/fid -run 'TestV01B(SmokeCreateListGetUpdateTransitionReadyAndRebuildParity|ReadyReturnsEmptyAfterAllTasksFinishedEndToEnd)$'`
  PASS after formatting.

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
- `PASS` `go test -tags integration -race ./...`
- `PASS` `go build -o /tmp/fid ./cmd/fid`

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- `AGENTS.md` should keep `End-to-end tests: N/A` for now. This task only recommends a held-out harness seed and isolation mechanism; it does not accept the harness location or concrete project command yet.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `PENDING`

## Commit subject

- `PENDING`
