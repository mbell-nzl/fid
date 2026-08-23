# Phase 21 - Task 2: Event V2, Idempotency, And Role-Free Bootstrap - Completion Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: specification-contradiction
- Condition: The required repository gates exercise retained V1 role-authorized mutations after a fresh `fid init`, while the accepted Event-v2 bootstrap contract requires empty Effects and a role-free projection. The resulting V1 role profile is empty, so every retained V1 mutation rejects its former `--by` identity.
- Recovery attempts: 1
- Next action: Tech Lead must decide whether temporary V1 role profiles are restored by an accepted migration/projection authority, or whether the retained V1 command/test contract is intentionally unavailable after role-free init. That decision must preserve the exact Event-v2 bootstrap fixture or explicitly amend it.
- Human attention required: yes

## Summary of what was implemented

- Added dual-schema journal support with strict Event-v2 decoding, canonical Event-v2 serialization, and aggregate sequence checks during append and replay.
- Added durable `idempotency_results` commit/replay/conflict handling with replay metadata rewritten only from `false` to `true`, plus same-key serialization through the write lock.
- Converted `project.initialized`, `fid init`, and `fid rebuild` to the role-free bootstrap and projection-only rebuild contracts. Rebuild no longer accepts `--by` or calls V1 policy authorization.

## Files created or modified

- `internal/journal/event.go`
- `internal/journal/aggregate.go`
- `internal/journal/idempotency.go`
- `internal/journal/replay.go`
- `internal/journal/upgrade.go`
- `internal/journal/writer.go`
- `internal/journal/event_v2_test.go`
- `internal/journal/idempotency_test.go`
- `internal/bootstrap/init.go`
- `internal/bootstrap/init_test.go`
- `cmd/fid/init_test.go`
- `cmd/fid/rebuild.go`
- `cmd/fid/rebuild_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-2-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-2-done.md` - records the protected gate conflict, red/green evidence, and exact command outcomes.

## Validation Contract Coverage

- VAL-003: PASS - `TestBuildProjectInitializedEventMatchesV2Golden` and `TestEventV2RejectsUnknownTopLevelKey` prove the exact closed Event-v2 bootstrap shape and unknown-key rejection.
- VAL-004: PASS - `TestAppendEventRejectsV2AggregateSequenceGap` and `TestReplayJournalRejectsOutOfOrderV2AggregateSequence` prove sequence-one enforcement and replay rejection of a gap.
- VAL-005: PASS - `TestCommitIdempotentReplaysStoredResultAfterRebuild` proves same-key replay retains Event identity, changes only `replayed`, and survives delete/rebuild.
- VAL-006: PASS - `TestCommitIdempotentRejectsSameKeyWithDifferentDigest` proves a digest conflict leaves one journal Event.
- VAL-007: PASS - `TestCommitIdempotentSerializesConcurrentSameKeyCalls` proves two concurrent calls return one Event identity with exactly one replay.
- VAL-008: PASS - `TestBuildProjectInitializedEventMatchesV2Golden` and `TestInitCommandSeedsExactlyOneRoleFreeBootstrapEvent` prove the service Actor, project aggregate sequence 1, null repository/idempotency, and empty Effects.

## Tests run

- `go test ./internal/bootstrap -run '^TestBuildProjectInitializedEventMatchesV2Golden$'` - FAIL (red) - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/journal -run '^TestCommitIdempotent'` - FAIL (red) - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/journal -run '^TestCommitIdempotent'` - PASS (green) - `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t0.120s` (terminated `exit_code: 0`)
- `go test ./internal/journal -run '^(TestAppendEventRejectsV2AggregateSequenceGap|TestReplayJournalRejectsOutOfOrderV2AggregateSequence|TestEventV2RejectsUnknownTopLevelKey)$'` - FAIL (red) - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/journal -run '^(TestAppendEventRejectsV2AggregateSequenceGap|TestReplayJournalRejectsOutOfOrderV2AggregateSequence|TestEventV2RejectsUnknownTopLevelKey)$'` - PASS (green) - `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t0.029s` (terminated `exit_code: 0`)
- `go test ./internal/bootstrap ./internal/journal ./cmd/fid -run '^(TestBuildProjectInitializedEventMatchesV2Golden|TestInitCommandSeedsExactlyOneRoleFreeBootstrapEvent|TestInitCommandRebuildRoundTripPreservesRoleFreeBootstrapAndProjectionMeta|TestCommitIdempotent|TestAppendEventRejectsV2AggregateSequenceGap|TestReplayJournalRejectsOutOfOrderV2AggregateSequence|TestEventV2RejectsUnknownTopLevelKey)$'` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.185s` (terminated `exit_code: 0`)
- `/tmp/fid-task2 --project /tmp/fid-task2-cli.xkIUsW --json init --prefix FID` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-08-23.jsonl","project":{"id":"project-114173b20e9b911a65e9497198830679","issue_prefix":"FID","root":"/tmp/fid-task2-cli.xkIUsW"}},"event":{"id":"EVT-20260823-000001","occurred_at":"2026-08-23T00:24:35.443567052Z"}}` (terminated `exit_code: 0`)
- `/tmp/fid-task2 --project /tmp/fid-task2-cli.xkIUsW --json rebuild` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"project.rebuild","result":{"events_replayed":1,"files_replayed":1,"schema_version":1,"upgraded_events":0},"event":null}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; retained V1 command tests require roles that the accepted role-free bootstrap no longer creates)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; retained V1 role-authorized packages fail after role-free init)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t\t44.9%` (terminated `exit_code: 0`)

## Deviations from spec

- The full required gates are blocked by the accepted bootstrap/V1 retention conflict. No out-of-scope role seeding, migration change, policy fallback, or test weakening was applied.
- Initial `golangci-lint run ./...` attempt failed with `1 issues:` (terminated `exit_code: 1`) because `initCommand` became unused after removing Event-v1 command fields. The in-scope repair removed that helper; the byte-identical rerun passed with `0 issues.` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- `contracts.md` section 14 says live bootstrap apply records only projection metadata and produces a role-free projection, while retained V1 command behavior and the required package gate rely on the old seeded role profile. A decision is required before this task can close without violating one of those authorities.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `94e752d25809a0e54f827b1872172f8038c0a741`

## Commit subject

- `feat(journal): add event v2 idempotency foundation`
