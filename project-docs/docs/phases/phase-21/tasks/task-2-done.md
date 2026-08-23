# Phase 21 - Task 2: Event V2, Idempotency, And Role-Free Bootstrap - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: N/A
- Condition: DEC-162 authorized retained V1 test and fixture convergence while preserving the accepted no-leaf-flags `project.rebuild` interface and role-free Event-v2 bootstrap.
- Recovery attempts: 6
- Next action: The Phase Runner runs the two-argument report-lint control against this committed report.
- Human attention required: no

The prior BLOCKED outcome and every failed required-command result remain immutable evidence below. DEC-162 resolved the omitted retained-test scope; it did not alter the accepted product contract.

## Summary of what was implemented

- Preserved the Event-v2, durable idempotency, role-free bootstrap, and DEC-160 nonpersistent seven-name policy vocabulary bridge from the original Task 2 implementation.
- Removed the vestigial role parameter from `journal.RebuildWorkspace`; role-free rebuilds now acquire the lock as exactly `fid rebuild`.
- Converged retained rebuild journeys by removing only their obsolete `--by <role>` argument pairs. Their existing projection parity, journal-byte, and read-shape assertions remain intact.
- Repointed incidental bootstrap-role SQL assertions to migration-seeded `workflows` rows, asserted the role-free vacuum projection, and replaced the invalid Event-v2 corrupt-tail fixture with a future-schema corrupt event while retaining its error, file, and line assertions.

## Files created or modified

- Task-owned implementation: `internal/journal/projection.go`, `internal/bootstrap/init.go`, and `cmd/fid/rebuild.go`.
- Retained rebuild test callers: `cmd/fid/claim_test.go`, `cmd/fid/e2e_v01b_test.go`, `cmd/fid/e2e_v01c_ac_test.go`, `cmd/fid/e2e_v01c_complete_test.go`, `cmd/fid/e2e_v01c_concern_test.go`, `cmd/fid/e2e_v01d_closure_test.go`, `cmd/fid/e2e_v01d_link_test.go`, `cmd/fid/e2e_v01e_activity_test.go`, `cmd/fid/e2e_v01e_dev_test.go`, `cmd/fid/issue_ac_test.go`, `cmd/fid/issue_test.go`, `cmd/fid/issue_transition_test.go`, `cmd/fid/root_flags_test.go`, and `cmd/fid/workspace_fixture_test.go`.
- Retained integration callers: `cmd/fid/integration/e2e_p20_issue_journey_test.go`, `cmd/fid/integration/e2e_v01_closure_helpers_test.go`, `cmd/fid/integration/e2e_v01_closure_test.go`, `cmd/fid/integration/e2e_v02a_workflow_authoring_test.go`, and `cmd/fid/integration/v01_parity_test.go`.
- Retained internal test callers and fixtures: `internal/issues/ac_test.go`, `internal/issues/issue_test.go`, `internal/issues/transition_test.go`, `internal/issues/update_test.go`, `internal/journal/idempotency_test.go`, `internal/journal/phase2_rebuild_test.go`, `internal/journal/project_test.go`, `internal/journal/rebuild_bench_test.go`, `internal/journal/rebuild_scale_test.go`, `internal/journal/replay_test.go`, `internal/journal/v01_tables_test.go`, `internal/maintenance/vacuum_test.go`, `internal/policy/policy_test.go`, `internal/sqlquery/sqlquery_test.go`, and `internal/workflow/assign_test.go`.
- Retained SQL CLI assertion: `cmd/fid/sql_test.go`.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-2-done.md`.

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-2-done.md` - preserves historical failed evidence and records DEC-162 convergence, red/green evidence, and final quality-gate results.

## Validation Contract Coverage

- VAL-003: PASS - `internal/journal/event_v2_test.go:54` (`TestEventV2RejectsUnknownTopLevelKey`) and `internal/bootstrap/init_test.go:15` (`TestBuildProjectInitializedEventMatchesV2Golden`) prove the closed Event-v2 bootstrap fixture and key rejection.
- VAL-004: PASS - `internal/journal/event_v2_test.go:13` and `internal/journal/event_v2_test.go:31` prove aggregate sequence-gap rejection during append and replay.
- VAL-005: PASS - `internal/journal/idempotency_test.go:15` proves same-key replay preserves Event identity across delete/rebuild and changes only `replayed`.
- VAL-006: PASS - `internal/journal/idempotency_test.go:58` proves a same-key digest conflict appends no second Event.
- VAL-007: PASS - `internal/journal/idempotency_test.go:82` proves two concurrent same-key calls return one Event identity with exactly one replay.
- VAL-008: PASS - `internal/bootstrap/init_test.go:15` and `cmd/fid/init_test.go:428` prove the exact service Actor, project sequence 1, null repository/idempotency, and empty Effects bootstrap.
- Intermediate compatibility invariant: PASS - the required `internal/policy` command proves the closed seven-name bridge, explicit-role precedence, fail-closed empty/corrupt behavior, and unknown-role rejection without persisted compatibility state.

## Tests run

- Prior red evidence for `go test ./internal/bootstrap -run '^TestBuildProjectInitializedEventMatchesV2Golden$'`: `FAIL` (terminated `exit_code: 1`).
- Prior red evidence for `go test ./internal/journal -run '^TestCommitIdempotent'`: `FAIL` (terminated `exit_code: 1`). Its prior green rerun printed `ok  github.com/mbell-nzl/fid/internal/journal  0.120s` (terminated `exit_code: 0`).
- Prior red evidence for `go test ./internal/journal -run '^(TestAppendEventRejectsV2AggregateSequenceGap|TestReplayJournalRejectsOutOfOrderV2AggregateSequence|TestEventV2RejectsUnknownTopLevelKey)$'`: `FAIL` (terminated `exit_code: 1`). Its prior green rerun printed `ok  github.com/mbell-nzl/fid/internal/journal  0.029s` (terminated `exit_code: 0`).
- DEC-160 policy bridge red evidence for the required `internal/policy` command: `FAIL` (terminated `exit_code: 1`). Its prior green rerun printed `ok  github.com/mbell-nzl/fid/internal/policy  1.018s` (terminated `exit_code: 0`).
- DEC-162 API red evidence: `go test ./internal/journal -run '^TestCommitIdempotentReplaysStoredResultAfterRebuild$'` printed `FAIL` because `RebuildWorkspace` still required three arguments (terminated `exit_code: 1`); after the role-free API repair it printed `ok  github.com/mbell-nzl/fid/internal/journal  0.060s` (terminated `exit_code: 0`).
- DEC-162 SQL fixture red evidence: `go test ./internal/sqlquery -run '^TestRunReturnsStructuredRows$'` printed `FAIL` because `escalation-default` is the second sorted migration workflow (terminated `exit_code: 1`); after correcting that migration-backed expectation it printed `ok  github.com/mbell-nzl/fid/internal/sqlquery  0.138s` (terminated `exit_code: 0`).
- `go test -tags integration ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration  3.767s` (terminated `exit_code: 0`).
- `/tmp/fid-task2-convergence --project /tmp/fid-task2-entrypoint.P9MZMC --json init --prefix FID` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-08-23.jsonl","project":{"id":"project-7883296cc2078cd7c69173682537a31d","issue_prefix":"FID","root":"/tmp/fid-task2-entrypoint.P9MZMC"}},"event":{"id":"EVT-20260823-000001","occurred_at":"2026-08-23T02:38:48.90909758Z"}}` (terminated `exit_code: 0`)
- `/tmp/fid-task2-convergence --project /tmp/fid-task2-entrypoint.P9MZMC --json rebuild` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"project.rebuild","result":{"events_replayed":1,"files_replayed":1,"schema_version":1,"upgraded_events":0},"event":null}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/policy -run '^(TestKnownRolesUsesV1CompatibilityVocabularyForRoleFreeV2Bootstrap|TestKnownRolesKeepsExplicitRoleFactsAuthoritative|TestKnownRolesRejectsCompatibilityForUninitializedJournal|TestKnownRolesRejectsCompatibilityForCorruptJournal|TestAuthorizeRejectsUnknownCompatibilityRole)$'` - PASS - `ok  github.com/mbell-nzl/fid/internal/policy  1.018s` (terminated `exit_code: 0`)
- Prior attempt of `go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...` printed `FAIL` (terminated `exit_code: 1`): retained V1 callers required roles after role-free init.
- DEC-160 byte-identical rerun of `go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...` printed `FAIL` (terminated `exit_code: 1`): retained tests still required forbidden `rebuild --by` and persisted role rows.
- DEC-162 byte-identical rerun of `go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  135.014s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- Prior attempt of `go test -race -coverprofile=coverage.out ./...` printed `FAIL` (terminated `exit_code: 1`): retained V1 role-authorized packages failed after role-free init.
- DEC-160 byte-identical rerun of `go test -race -coverprofile=coverage.out ./...` printed `FAIL` (terminated `exit_code: 1`): retained rebuild-flag and role-row assertions remained.
- DEC-162 first byte-identical rerun of `go test -race -coverprofile=coverage.out ./...` printed `FAIL` (terminated `exit_code: 1`): `TestVacuumRefreshesMissingProjectionBeforeCompacting` still expected bootstrap roles, and `TestCommitIdempotentSerializesConcurrentSameKeyCalls` observed `database is locked (5)`.
- DEC-162 later byte-identical rerun of `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow  25.511s  coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:                                  (statements)                          78.2%` (terminated `exit_code: 0`)

## Deviations from spec

- The original package and repository coverage failures remain immutable evidence. Their retained V1 role/flag assumptions conflicted with the accepted role-free bootstrap and no-leaf-flags rebuild surface; DEC-162 authorized adapting the affected test collateral without changing either product contract. The later byte-identical package and repository reruns passed.
- The first DEC-162 repository coverage rerun ended `FAIL` (terminated `exit_code: 1`). Its actionable retained-test root cause was `TestVacuumRefreshesMissingProjectionBeforeCompacting` expecting default bootstrap roles; the repair asserts the accepted empty role projection while keeping rebuild and journal-byte checks. The concurrent idempotency test also observed one `SQLITE_BUSY` result; it passed ten focused race-enabled repetitions and the later byte-identical repository rerun. No test was skipped, narrowed, or weakened.
- Deleted assertion: `TestIssueCreateAndRebuildValidateRoleExistence/rebuild missing by`; accepted owner: `TestRebuildJSONSucceedsWithoutActorAuthority` in `cmd/fid/rebuild_test.go`.
- Deleted assertion: `TestIssueCreateAndRebuildValidateRoleExistence/rebuild unknown role`; accepted owner: `TestRebuildRejectsEveryFormerByValueWithoutMutation` in `cmd/fid/rebuild_test.go`.
- Explicit legacy `roles` Effect fixtures added for convergence: none.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
