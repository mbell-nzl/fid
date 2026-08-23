# Phase 21 - Task 3: Capability Discovery - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: The Phase Runner runs the two-argument report-lint control against this committed report.
- Human attention required: no

## Summary of what was implemented

- Added the migration-derived binary capability projection, which re-materializes the fixed Phase 21 matrix into `system_capabilities` and reads it back in capability-ID order.
- Added `fid capabilities` / `capabilities.get` with the exact 12-row JSON result, no-argument typed rejection, and no event.
- Updated the shared CLI envelope serializer to emit lexicographically ordered JSON object keys required by the public transport contract.

## Files created or modified

- `internal/workflow/capabilities.go`
- `cmd/fid/capabilities.go`
- `cmd/fid/root.go`
- `internal/cli/envelope.go`
- `cmd/fid/capabilities_test.go`
- `cmd/fid/contract_surface_test.go`
- `internal/cli/envelope_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-3-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-3-done.md` - records task evidence, the immutable initial gate failure, and the required final gate results.

## Validation Contract Coverage

- VAL-009: PASS - `TestCapabilitiesJSONGoldenAndProjectionRows` in `cmd/fid/capabilities_test.go:16` compares the complete literal `fid capabilities --json` result at `:22` and reads the `system_capabilities` projection in capability-ID order at `:28`; the focused green run ended `ok  	github.com/mbell-nzl/fid/cmd/fid	2.407s` (terminated `exit_code: 0`)
- VAL-010: PASS - `TestCapabilitiesRejectsOverridesAndKeepsDeferredRowsUnavailable` in `cmd/fid/capabilities_test.go:55` asserts the exact typed no-argument rejection at `:61` and sets each of the seven named deferred rows in `:76` to an optimistic `available=1` before confirming discovery restores `false` at `:85`; the required package rerun ended `ok  	github.com/mbell-nzl/fid/cmd/fid	131.062s` (terminated `exit_code: 0`)

## Tests run

- `go test -race ./cmd/fid -run '^TestCapabilities' -count=1` - RED - `FAIL` (terminated `exit_code: 1`)
  Before implementation, `fid capabilities` was not registered and both capability tests received an execution error.
- `go test -race ./cmd/fid -run '^TestCapabilities' -count=1` - GREEN - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.407s` (terminated `exit_code: 0`)
- `go run ./cmd/fid --project /tmp/fid-task-3-entrypoint.H7F4xk --json capabilities` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"capabilities.get","result":{"items":[{"available":false,"capability_id":"attention.v1","limits":{"state":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":true,"capability_id":"evidence.freshness.v1","limits":{"branch_labels":false,"target":"repository-tree"},"schema_max":1,"schema_min":1,"trust_source":"fid-git-observation"},{"available":true,"capability_id":"evidence.gate-definitions.v1","limits":{"executor_policy":"standard-library","network":"allowed","sandbox":"none","write":"allowed"},"schema_max":1,"schema_min":1,"trust_source":"active-kit-activation"},{"available":true,"capability_id":"evidence.receipts.v1","limits":{"attestation":false,"producer":"fid","terminal_receipts":true},"schema_max":1,"schema_min":1,"trust_source":"fid-executed"},{"available":false,"capability_id":"export.portable.v1","limits":{"create":true,"import":false,"validate":true},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":false,"capability_id":"human.outcomes.v1","limits":{"record":false},"schema_max":1,"schema_min":1,"trust_source":"phase-23-unavailable"},{"available":false,"capability_id":"kit.activation.v1","limits":{"abort_before_commit":true,"automatic_respawn":false,"finalize":true,"recover":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":false,"capability_id":"reviews.coverage.v1","limits":{"record":false},"schema_max":1,"schema_min":1,"trust_source":"phase-23-unavailable"},{"available":false,"capability_id":"telemetry.trajectory.v1","limits":{"ingest":false},"schema_max":1,"schema_min":1,"trust_source":"phase-25-unavailable"},{"available":false,"capability_id":"work.claims.v1","limits":{"max_active_per_issue":1,"parallelism":false,"recover":false},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"},{"available":true,"capability_id":"workflow.state.v1","limits":{"max_envelope_revisions":1,"max_primary_issues":1,"related_issues":false},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"},{"available":true,"capability_id":"workflow.transitions.v1","limits":{"execution_bound":true,"terminal_evidence_guard":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"}]}}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/workflow ./internal/cli ./cmd/fid/...` - FAIL - `FAIL` (terminated `exit_code: 1`).
- `go test -race ./internal/workflow ./internal/cli ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	131.062s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	25.341s	coverage: 75.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.7%` (terminated `exit_code: 0`)

## Deviations from spec

- The first required package gate ended `FAIL` (terminated `exit_code: 1`)
  `TestWriteFailureUsesExactContractKeys` still expected the prior non-lexicographic `subject` object order. The root cause was test collateral predating the public transport's lexicographic-key rule. The repair changed only that expected JSON order to `id,type`. No test was skipped, narrowed, or weakened. The byte-identical package gate rerun ended `ok  	github.com/mbell-nzl/fid/cmd/fid	131.062s` (terminated `exit_code: 0`)

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `6ff90152fdcdd373c14a434ec736e557144a1617`

## Commit subject

- `feat(capabilities): add capability discovery transport`
