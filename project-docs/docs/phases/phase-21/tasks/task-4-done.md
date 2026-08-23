# Phase 21 - Task 4: Activation Trust Lifecycle - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added the durable Activation lifecycle: prepare, component staging, commit, pre-commit abort, finalize, and read-only status.
- Activation components are content-addressed, Event-v2/idempotent mutations; commit selects both staged component digests and becomes the point of no return.
- Added finalized-only Actor resolution. Bootstrap attribution and caller-supplied augmentation fields cannot create trusted authority.
- Kept recovery and respawn unavailable, with post-commit abort returning `activation-recovery-unavailable`.

## Files created or modified

- `internal/workflow/activation.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `cmd/fid/activation.go`
- `cmd/fid/activation_test.go`
- `cmd/fid/root.go`
- `cmd/fid/capabilities_test.go`
- `cmd/fid/contract_surface_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-4-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-4-done.md`

## Validation Contract Coverage

- VAL-011: PASS - `TestActivationCommitIsPointOfNoReturn` in `cmd/fid/activation_test.go` proves staged components commit together and exposes `committed` with `point_of_no_return:true`; `TestActivationFinalizeAndStatusExposeOnlyFinalLifecycleState` proves the only forward terminal state is `finalized`.
- VAL-012: PASS - `TestActivationCommitIsPointOfNoReturn` asserts the exact post-commit abort rejection; `AbortActivation` permits only `prepared|staged` in `internal/workflow/activation.go`.
- VAL-013: PASS - `TestActivationCommitIsPointOfNoReturn` observes `activation-recovery-unavailable` after commit; the declared source-absence command passed with no recovery command/path.
- VAL-014: PASS - `TestFinalizedActivationIsTheOnlyTrustedActorSource` proves the stored profile is unavailable while staged and committed, then verifies exact type, producer, provenance, and capabilities after finalization.
- VAL-015: PASS - `TestActivationComponentStagesValidatedProfilesAndRejectsCallerAugmentation` injects `caller_capabilities`, observes `invalid-argument`, and asserts no successful output.

## Tests run

- `go test -race ./cmd/fid -run '^TestActivationPrepareCreatesPreparedActivation$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t1.823s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestActivationComponentStagesValidatedProfilesAndRejectsCallerAugmentation$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t1.903s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestActivationCommitIsPointOfNoReturn$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t2.121s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestActivationFinalizeAndStatusExposeOnlyFinalLifecycleState$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t2.194s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestFinalizedActivationIsTheOnlyTrustedActorSource$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t2.314s` (terminated `exit_code: 0`)
- `/tmp/fid-task4 --project /tmp/FIDTASK4b9FjVd --json activation prepare --input - --idempotency-key smoke-activation-prepare` - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260823-000002","occurred_at":"2026-08-23T05:30:33.707873055Z"},"ok":true,"operation":"activation.prepare","result":{"activation":{"id":"ACT-ebdf0d","status":"prepared","kit_identity":"kit","backend_identity":"fid","phase":21,"direction":"activate","point_of_no_return":false,"component_digests":[],"created_at":"2026-08-23T05:30:33.707873055Z","updated_at":"2026-08-23T05:30:33.707873055Z"},"idempotency_key":"smoke-activation-prepare","replayed":false}}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t(cached)` (terminated `exit_code: 0`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"recover|newActivationRecover|activation\.recover' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t25.205s\tcoverage: 49.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t74.7%` (terminated `exit_code: 0`)

## Deviations from spec

- Required focused gate attempt 1: `go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` - FAIL - `FAIL` (terminated `exit_code: 1`). Root cause: task-3 inventory tests still treated the new accepted activation leaves as absent and the CLI registry expected its pre-task code count. Repair: updated the exact task-4 contract inventories and schema expectations without changing the accepted capability availability. Later byte-identical result: `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t(cached)` (terminated `exit_code: 0`).
- Required lint attempt 1: `golangci-lint run ./...` - FAIL - `2 issues:` (terminated `exit_code: 1`). Root cause: two unused activation helpers. Repair: removed the unused declarations. Later byte-identical result: `0 issues.` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `ad9098fac66f2b530a566d87df087fd72fbb584d`

## Commit subject

- `feat(activation): add trusted activation lifecycle`
