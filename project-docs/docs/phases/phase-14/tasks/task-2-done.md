## Summary of what was implemented

- Added the new `internal/runs` domain package with `Record`, the first Phase 14 telemetry write path for `agent_runs`.
- Implemented `fid run record` in `cmd/fid`, including known-model cost computation, unknown-model warning with NULL pricing/cost, optional provider reconciliation fields, and stderr confirmation output.
- Added the `run.record` permissions-matrix row and parity coverage so `fid policy show` and `Authorize` use the same source of truth.

## Files created or modified

- `internal/runs/record.go`
- `internal/runs/helpers.go`
- `internal/runs/record_test.go`
- `internal/policy/matrix.go`
- `internal/policy/matrix_test.go`
- `cmd/fid/run.go`
- `cmd/fid/run_test.go`
- `cmd/fid/root.go`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-14/tasks/task-2-done.md`

## Documentation updated

- `project-docs/specs/02_TRD.md`: recorded the new `internal/runs` boundary in the architecture diagram and module-boundary table (PI-1).
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: aligned `fid run record` with the implemented Phase 14 surface, removed the future `--review` flag that belongs to the v0.3 `review_id` schema extension, and updated the stderr confirmation example.

## Validation Contract Coverage

- `VAL-002` PASS: `internal/runs.Record` and `fid run record` append exactly one `run.recorded` event and one `agent_runs` row for a known model, with zero `activity_log` rows. PI-1 first applied by introducing `internal/runs`; PI-2 first applied by emitting only the `agent_runs` insert effect.
- `VAL-003` PASS: known-model `cost_usd` is computed from the same row’s snapshotted `price_in_per_mtok` and `price_out_per_mtok`.
- `VAL-004` PASS: unknown-model runs still record, return exit 0, print a warning naming the model, and persist NULL pricing/cost fields. PI-3 first applied in the unknown-model branch.
- `VAL-005` PASS: `--token-breakdown-json`, `--tool-call-count`, `--provider-run-id`, and `--provider-cost-usd` persist when supplied, stay NULL when omitted, and do not change tokens or normalized cost math.
- `VAL-006` PASS: added exactly one `run.record` matrix row matching `03_TRD_SCHEMA.md` section 13; `fid policy show` exposes the row from the same matrix used by `Authorize`; all seven default-profile actor roles were exercised successfully; an unknown acting role is rejected before append with unchanged journal length.
- `PI-4` not applicable in Task 2.
- `PI-5` not applicable in Task 2.

## Tests run

- `go test ./internal/runs -run TestRecordKnownModelPersistsAgentRunAndComputedCost`
  Red: failed with `undefined: Record` / `undefined: RecordOptions`.
  Green: passed after adding `internal/runs.Record` and the `run.record` matrix row.
- `go test ./cmd/fid -run TestRunRecordCreatesAgentRunAndPrintsConfirmation`
  Red: failed with `unknown command "run" for "fid"`.
  Green: passed after adding `cmd/fid/run.go` and registering `newRunCommand()`.
- `go test ./internal/runs ./internal/policy ./cmd/fid`
  PASS.

## Exact quality-gate results

- PASS `gofmt -l cmd internal`
- PASS `golangci-lint run ./...`
- PASS `go test -race -coverprofile=coverage.out ./...`
- PASS `go tool cover -func=coverage.out`
- PASS `go build -o /tmp/fid ./cmd/fid`
- Conditional command not run: `go test -tags integration -race ./...` (no integration-tagged tests added in Task 2)
- Coverage: `internal/runs` 75.1%, `internal/policy` 84.6%, `cmd/fid` 77.8%
- Total statement coverage from `go tool cover -func=coverage.out`: 76.1%
- Built-binary evidence: `/tmp/fid` recorded `RUN-bb9ab2` on `FID-dc4c61` in temp workspace `/tmp/tmp.z9oSMVzOsg` with stderr confirmation:
  `recorded agent run RUN-bb9ab2 on FID-dc4c61`, `role: developer, model: claude-sonnet-4-6`, `tokens: 1200 in, 3400 out`, `cost: $0.0546`

## Deviations from spec

- None. The only spec-touch beyond the task body was reconciling pre-existing CLI-doc drift by removing the future `--review` flag from the Phase 14 `fid run record` section to match the current schema and VAL-005 note.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- TBD

## Commit subject

- TBD
