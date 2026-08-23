# Phase 21 - Task 4: Activation Trust Lifecycle - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added the durable Activation lifecycle: prepare, component staging, commit, pre-commit abort, finalize, and read-only status.
- Activation components are content-addressed, Event-v2/idempotent mutations; commit selects both staged component digests and becomes the point of no return.
- Added finalized-only Actor resolution. Bootstrap attribution and caller-supplied augmentation fields cannot create trusted authority.
- Kept recovery and respawn unavailable, with post-commit abort returning `activation-recovery-unavailable`.
- Review-2 MF-1 adds CLI-level lifecycle rejection, pre-commit abort, and caller-augmentation tables without changing production behavior or accepted command contracts.

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

- VAL-011: PASS - `TestActivationLifecycleRejectsSkippedAndExtraTransitionsWithoutEvents` in `cmd/fid/activation_test.go` exercises every skipped or repeated CLI lifecycle action from `prepared`, `staged`, `committed`, `finalized`, and `aborted`; every table row receives a typed rejection and preserves its journal Event count.
- VAL-012: PASS - `TestActivationAbortPreCommitLeavesNoTrustedActor` proves abort succeeds from both `prepared` and `staged`, emits exactly one abort Event, and leaves `ResolveTrustedActor` unavailable; `TestActivationLifecycleRejectsSkippedAndExtraTransitionsWithoutEvents` preserves the committed/finalized post-commit abort rejection with no Event.
- VAL-013: PASS - `TestActivationCommitIsPointOfNoReturn` observes `activation-recovery-unavailable` after commit; the declared source-absence command passed with no recovery command/path.
- VAL-014: PASS - `TestFinalizedActivationIsTheOnlyTrustedActorSource` proves the stored profile is unavailable while staged and committed, then verifies exact type, producer, provenance, and capabilities after finalization.
- VAL-015: PASS - `TestActivationComponentStagesValidatedProfilesAndRejectsCallerAugmentation` table-injects both `caller_capabilities` and `caller_provenance_ref`, observes `invalid-argument` for each, and asserts the journal Event count remains unchanged.

## Tests run

- `go test -race ./cmd/fid -run '^TestActivationPrepareCreatesPreparedActivation$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t1.823s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestActivationComponentStagesValidatedProfilesAndRejectsCallerAugmentation$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t1.903s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestActivationCommitIsPointOfNoReturn$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t2.121s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestActivationFinalizeAndStatusExposeOnlyFinalLifecycleState$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t2.194s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestFinalizedActivationIsTheOnlyTrustedActorSource$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t2.314s` (terminated `exit_code: 0`)
- `/tmp/fid-task4 --project /tmp/FIDTASK4b9FjVd --json activation prepare --input - --idempotency-key smoke-activation-prepare` - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260823-000002","occurred_at":"2026-08-23T05:30:33.707873055Z"},"ok":true,"operation":"activation.prepare","result":{"activation":{"id":"ACT-ebdf0d","status":"prepared","kit_identity":"kit","backend_identity":"fid","phase":21,"direction":"activate","point_of_no_return":false,"component_digests":[],"created_at":"2026-08-23T05:30:33.707873055Z","updated_at":"2026-08-23T05:30:33.707873055Z"},"idempotency_key":"smoke-activation-prepare","replayed":false}}` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestActivationLifecycleRejectsSkippedAndExtraTransitionsWithoutEvents$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Red: the lifecycle test referenced its table fixture before that test helper existed.
- `go test -race ./cmd/fid -run '^(TestActivationComponentStagesValidatedProfilesAndRejectsCallerAugmentation|TestActivationLifecycleRejectsSkippedAndExtraTransitionsWithoutEvents|TestActivationAbortPreCommitLeavesNoTrustedActor)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Red: the first complete table run exposed the fixture ordering error and verified the CLI uses the existing typed `invalid-argument` envelope for rejected non-recovery transitions.
- Focused green evidence for the byte-identical lifecycle table run retained token `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t6.723s` with observed termination `exit_code: 0`; all new MF-1 rows pass.

## Exact quality-gate results

- `go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t(cached)` (terminated `exit_code: 0`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"recover|newActivationRecover|activation\.recover' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t25.205s\tcoverage: 49.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t74.7%` (terminated `exit_code: 0`)

### Review-2 MF-1 recovery rerun

- `go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t149.185s` (terminated `exit_code: 0`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"recover|newActivationRecover|activation\.recover' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 38.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t70.7%` (terminated `exit_code: 0`)

## Deviations from spec

- Required focused gate attempt 1: `go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` - FAIL - `FAIL` (terminated `exit_code: 1`). Root cause: task-3 inventory tests still treated the new accepted activation leaves as absent and the CLI registry expected its pre-task code count. Repair: updated the exact task-4 contract inventories and schema expectations without changing the accepted capability availability. Later byte-identical result: `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t(cached)` (terminated `exit_code: 0`).
- Required lint attempt 1: `golangci-lint run ./...` - FAIL - `2 issues:` (terminated `exit_code: 1`). Root cause: two unused activation helpers. Repair: removed the unused declarations. Later byte-identical result: `0 issues.` (terminated `exit_code: 0`).
- Review-2 MF-1 added no specification deviation. The two new focused red attempts above were test construction evidence: first the absent table fixture, then fixture ordering and the CLI's already-established typed rejection envelope. The final green run adds only executable assertions and preserves the accepted lifecycle and error behavior.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `ad9098fac66f2b530a566d87df087fd72fbb584d`

## Commit subject

- `feat(activation): add trusted activation lifecycle`
