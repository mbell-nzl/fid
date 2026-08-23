# Phase 21 - Task 2: Event V2, Idempotency, And Role-Free Bootstrap - Completion Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: specification-contradiction
- Condition: DEC-160's allowed `internal/policy` bridge restores retained V1 role lookup after the exact role-free Event-v2 bootstrap, but the required package and repository gates still contain retained tests that require the prohibited `fid rebuild --by` flag and persisted `roles` rows. The Phase 21 interface requires `project.rebuild` to accept no leaf flags, and contracts section 14 requires the bootstrap to have empty Effects and a role-free projection.
- Recovery attempts: 3
- Next action: Tech Lead must provide protected authority to reconcile the retained V1 rebuild/role-row test contract with the accepted role-free bootstrap and no-flag rebuild interface. The authority must not reintroduce bootstrap role Effects or change the Event-v2 fixture without an explicit contract decision.
- Human attention required: yes

The prior BLOCKED outcome remains part of this report. Its original gate failures
and all earlier red/green evidence are retained below.

## Summary of what was implemented

- Added the DEC-160 temporary, nonpersistent V1 policy vocabulary bridge inside `internal/policy` only. It recognizes the seven V1 matrix roles after an exact role-free Event-v2 bootstrap and only when no legacy `roles` Effect exists.
- Kept explicit legacy role Effects authoritative, including an inactive role fact that leaves the profile empty rather than supplementing it.
- Made the lookup fail closed for empty and corrupt journals and unknown names; the lookup changes no journal bytes, Effect, projection row, config, or export field.
- Preserved the existing Event-v2, idempotency, role-free bootstrap, and rebuild implementation from the original task attempt.

## Files created or modified

- `internal/policy/policy.go`
- `internal/policy/policy_test.go`
- `internal/policy/authorize_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-2-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-2-done.md` - preserves the previous BLOCKED evidence and records the DEC-160 bridge red/green, required gate reruns, and protected contradiction.

## Validation Contract Coverage

- VAL-003: PASS (focused evidence retained) - `TestBuildProjectInitializedEventMatchesV2Golden` and `TestEventV2RejectsUnknownTopLevelKey` prove the closed Event-v2 bootstrap shape and unknown-key rejection.
- VAL-004: PASS (focused evidence retained) - `TestAppendEventRejectsV2AggregateSequenceGap` and `TestReplayJournalRejectsOutOfOrderV2AggregateSequence` prove sequence-one enforcement and replay rejection of a gap.
- VAL-005: PASS (focused evidence retained) - `TestCommitIdempotentReplaysStoredResultAfterRebuild` proves same-key replay retains Event identity, changes only `replayed`, and survives delete/rebuild.
- VAL-006: PASS (focused evidence retained) - `TestCommitIdempotentRejectsSameKeyWithDifferentDigest` proves a digest conflict leaves one journal Event.
- VAL-007: PASS (focused evidence retained) - `TestCommitIdempotentSerializesConcurrentSameKeyCalls` proves two concurrent calls return one Event identity with exactly one replay.
- VAL-008: PASS (focused evidence retained) - `TestBuildProjectInitializedEventMatchesV2Golden` and `TestInitCommandSeedsExactlyOneRoleFreeBootstrapEvent` prove the service Actor, project aggregate sequence 1, null repository/idempotency, and empty Effects.
- Intermediate compatibility invariant: PASS (focused) - `TestKnownRolesUsesV1CompatibilityVocabularyForRoleFreeV2Bootstrap` verifies exactly the seven names and unchanged journal bytes; `TestKnownRolesKeepsExplicitRoleFactsAuthoritative`, `TestKnownRolesRejectsCompatibilityForUninitializedJournal`, `TestKnownRolesRejectsCompatibilityForCorruptJournal`, and `TestAuthorizeRejectsUnknownCompatibilityRole` prove precedence and fail-closed behavior.

## Tests run

- Prior red evidence for `go test ./internal/bootstrap -run '^TestBuildProjectInitializedEventMatchesV2Golden$'`: `FAIL` (terminated `exit_code: 1`).
- Prior red evidence for `go test ./internal/journal -run '^TestCommitIdempotent'`: `FAIL` (terminated `exit_code: 1`). Its prior green rerun printed `ok  github.com/mbell-nzl/fid/internal/journal  0.120s` (terminated `exit_code: 0`).
- Prior red evidence for `go test ./internal/journal -run '^(TestAppendEventRejectsV2AggregateSequenceGap|TestReplayJournalRejectsOutOfOrderV2AggregateSequence|TestEventV2RejectsUnknownTopLevelKey)$'`: `FAIL` (terminated `exit_code: 1`). Its prior green rerun printed `ok  github.com/mbell-nzl/fid/internal/journal  0.029s` (terminated `exit_code: 0`).
- Prior focused Event-v2 evidence for `go test ./internal/bootstrap ./internal/journal ./cmd/fid -run '^(TestBuildProjectInitializedEventMatchesV2Golden|TestInitCommandSeedsExactlyOneRoleFreeBootstrapEvent|TestInitCommandRebuildRoundTripPreservesRoleFreeBootstrapAndProjectionMeta|TestCommitIdempotent|TestAppendEventRejectsV2AggregateSequenceGap|TestReplayJournalRejectsOutOfOrderV2AggregateSequence|TestEventV2RejectsUnknownTopLevelKey)$'` printed `ok  github.com/mbell-nzl/fid/cmd/fid  0.185s` (terminated `exit_code: 0`).
- DEC-160 policy bridge red evidence for `go test -race ./internal/policy -run '^(TestKnownRolesUsesV1CompatibilityVocabularyForRoleFreeV2Bootstrap|TestKnownRolesKeepsExplicitRoleFactsAuthoritative|TestKnownRolesRejectsCompatibilityForUninitializedJournal|TestKnownRolesRejectsCompatibilityForCorruptJournal|TestAuthorizeRejectsUnknownCompatibilityRole)$'`: `FAIL` (terminated `exit_code: 1`). The subsequent focused green run printed `ok  github.com/mbell-nzl/fid/internal/policy  1.018s` (terminated `exit_code: 0`); its canonical required-command result is recorded below.
- `/tmp/fid-task2-policy-bridge --project /tmp/fid-task2-policy-bridge.zJqnvO --json init --prefix FID` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-08-23.jsonl","project":{"id":"project-2a8cd4edf113ecba40410a5423d1a5f5","issue_prefix":"FID","root":"/tmp/fid-task2-policy-bridge.zJqnvO"}},"event":{"id":"EVT-20260823-000001","occurred_at":"2026-08-23T01:37:21.581377191Z"}}` (terminated `exit_code: 0`)
- `/tmp/fid-task2-policy-bridge --project /tmp/fid-task2-policy-bridge.zJqnvO --json issue create --kind task --title 'Role-free bridge' --by developer` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"issue.create","result":{"issue":{"id":"FID-74707b","kind":"task","revision":1,"status":"todo","title":"Role-free bridge"}},"event":{"id":"EVT-20260823-000002","occurred_at":"2026-08-23T01:37:25.492211168Z"}}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/policy -run '^(TestKnownRolesUsesV1CompatibilityVocabularyForRoleFreeV2Bootstrap|TestKnownRolesKeepsExplicitRoleFactsAuthoritative|TestKnownRolesRejectsCompatibilityForUninitializedJournal|TestKnownRolesRejectsCompatibilityForCorruptJournal|TestAuthorizeRejectsUnknownCompatibilityRole)$'` - PASS - `ok  github.com/mbell-nzl/fid/internal/policy  1.018s` (terminated `exit_code: 0`)
- The prior attempt of `go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...` printed `FAIL` (terminated `exit_code: 1`) because retained V1 command tests required roles after role-free init.
- The DEC-160 byte-identical rerun of `go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...` printed `FAIL` (terminated `exit_code: 1`) because retained tests require forbidden `rebuild --by` and seeded `roles` rows.
- `go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- The prior attempt of `go test -race -coverprofile=coverage.out ./...` printed `FAIL` (terminated `exit_code: 1`) because retained V1 role-authorized packages failed after role-free init.
- The DEC-160 byte-identical rerun of `go test -race -coverprofile=coverage.out ./...` printed `FAIL` (terminated `exit_code: 1`) because retained rebuild flag and role-row expectations conflict with the accepted role-free contracts.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total:                                  (statements)                          78.1%` (terminated `exit_code: 0`)

## Deviations from spec

- Required command failure remains immutable evidence. The original `go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...` failed with `FAIL` (terminated `exit_code: 1`) because `KnownRoles` returned no roles after the required role-free bootstrap. DEC-160 authorized the internal/policy bridge; the byte-identical rerun still ended `FAIL` (terminated `exit_code: 1`) because retained tests require the prohibited `rebuild --by` interface and persisted role rows.
- The original `go test -race -coverprofile=coverage.out ./...` failed with `FAIL` (terminated `exit_code: 1`) on the same empty-role profile. After the bridge, the byte-identical rerun still ended `FAIL` (terminated `exit_code: 1`) on retained rebuild/role-row tests and `TestKnownRolesTRG002CorruptTailTolerance`, whose invalid Event-v2 fixture predates Phase 21 (`716329cc`). No test was weakened, removed, or repointed.
- The initial `golangci-lint run ./...` attempt failed with `1 issues:` (terminated `exit_code: 1`) because `initCommand` became unused after removing Event-v1 command fields. The prior in-scope repair removed that helper; the byte-identical rerun passed with `0 issues.` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- `project-docs/docs/phases/phase-21/public-interface.md:330` requires `fid rebuild` to have no leaf flags, while the current required tests still invoke `rebuild --by workflow_runner` and assert a populated `roles` table. The bridge does not and must not restore either contract.
- The compatibility bridge is intentionally temporary and is identified for deletion by task 20 in `v1-convergence-inventory.md` section 6.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `861e90da8b9218915ddbd77515326921ff3c0cf4`

## Commit subject

- `feat(policy): add temporary role vocabulary bridge`
