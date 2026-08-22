# Phase 20 - Task 9: Issue mutation and concurrency contract - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: Fresh final QA re-review.
- Human attention required: no

## Summary of what was implemented

- Migrated canonical `issue create`, `issue update`, and `issue transition` to `fid-cli.v1` Response Envelopes, with Human success on stdout and canonical JSON failures on stderr.
- Added optional `--expected-revision` to update and transition. The current aggregate revision is checked while the write lock is held; a stale value returns typed `state-conflict` with both revisions and appends nothing.
- Returned the appended Event identity and resulting aggregate revision from each mutation, and sorted `issue.update.changes` by field.
- Resolved Issue mutation projects through command-context project discovery/explicit selection, kept V1 `--by` authorization fail-closed, and mapped transition preconditions to typed `transition-not-permitted` failures.
- Removed public create/update owner flags and their Event command serialization. Internal owner defaults/storage remain unchanged, while canonical mutations cannot change owners.
- Updated real-process integration helpers to consume migrated Issue-create Human output from stdout.

## Files created or modified

- Added `cmd/fid/issue_mutation_response.go`.
- Updated Issue create/update/transition adapters, shared Issue error rendering, and related command tests under `cmd/fid/`.
- Updated Issue create/update/transition request/result handling and tests under `internal/issues/`.
- Updated `cmd/fid/integration/e2e_v01_closure_helpers_test.go` for the migrated Human create channel.

## Documentation updated

- Added this completion report: `project-docs/docs/phases/phase-20/tasks/task-9-done.md`.

## Validation Contract Coverage

- VAL-002: PASS - `TestIssueUpdateRejectsStaleExpectedRevision` at `cmd/fid/issue_update_test.go:132` decodes the JSON failure branch and asserts empty stdout; the real `/tmp/fid` stale update emitted `state-conflict` on stderr and terminated `exit_code: 1`.
- VAL-008: PASS - `TestIssueMutationHumanSuccessUsesStdoutOnly` at `cmd/fid/issue_update_test.go:262` proves Human create, update, and transition successes are stdout-only.
- VAL-010: PASS - `internal/issues/create.go:216`, `internal/issues/update.go:217`, and `internal/issues/transition.go:301` return revisions 1 and observed-plus-one respectively; the entrypoint create/update evidence observed revisions 1 then 2.
- VAL-011: PASS - `internal/issues/update.go:159` and `internal/issues/transition.go:233` preserve Task 6's absolute revision Effects, while `internal/issues/create.go:216` returns the Event that inserted revision 1; the Task 9 real entrypoint observed distinct appended Event IDs for create and update.
- VAL-017: PASS - `cmd/fid/issue_mutation_response.go:48`, `:68`, and `:82` construct the three mutation envelopes from the returned appended Event identity; `/tmp/fid --json issue create` emitted `EVT-20260812-000002` and the following update emitted `EVT-20260812-000003`.
- VAL-018: PASS - `TestIssueUpdateRejectsStaleExpectedRevision` at `cmd/fid/issue_update_test.go:132` proves stale atomic rejection and exact values; `TestIssueMutationHumanSuccessUsesStdoutOnly` at `:262` proves an absent revision succeeds; the real entrypoint matching revision `1` proceeded to revision `2`.
- VAL-019: PASS - `TestIssueTransitionJSONRejectsUnsatisfiedClaimPrecondition` at `cmd/fid/issue_transition_test.go:15` asserts `transition-not-permitted`, `todo`/`in_progress`, and `live_claim_held`; `internal/issues/transition.go:157` rejects missing edges before append.
- VAL-040: PASS - `TestIssueCreateRejectsRemovedOwnerFlagsWithJSONFailure` at `cmd/fid/issue_update_test.go:229` proves both public owner flags are unknown; `internal/issues/update_changes.go` has no owner change generation and `internal/issues/update_test.go` asserts owner fields remain unchanged.
- VAL-041: PASS - `TestUpdateAppendsOneEventProjectsActivityAndRebuilds` at `internal/issues/update_test.go:244` asserts the supplied `tl` actor on the projected mutation activity and checks the Event's scoped Issue/activity effects; canonical mutation serializers in `cmd/fid/issue_mutation_response.go` expose no trusted/protected side-effect field.
- VAL-044: PASS - `TestIssueCreateRejectsRemovedOwnerFlagsWithJSONFailure` at `cmd/fid/issue_update_test.go:229`, `TestIssueUpdateRejectsStaleExpectedRevision` at `:132`, and the closed flags in `cmd/fid/issue_create.go`/`issue_update.go` cover removed inputs, exact stale handling, and closed mutation results.
- VAL-045: PASS - `TestIssueUpdateJSONOrdersChangesByField` at `cmd/fid/issue_update_test.go:179` asserts priority, spec path, and title ordering; `internal/issues/update.go:112` performs the deterministic sort.

## Tests run

- RED: `go test ./cmd/fid -run '^TestIssueUpdateRejectsStaleExpectedRevision$' -count=1` - expected FAIL - `FAIL` (terminated `exit_code: 1`).
- GREEN: `go test ./cmd/fid -run '^TestIssueUpdateRejectsStaleExpectedRevision$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.055s` (terminated `exit_code: 0`).
- RED: `go test ./cmd/fid -run '^TestIssueUpdateJSONOrdersChangesByField$' -count=1` - expected FAIL - `FAIL` (terminated `exit_code: 1`).
- GREEN: `go test ./cmd/fid -run '^TestIssueUpdateJSONOrdersChangesByField$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.056s` (terminated `exit_code: 0`).
- RED: `go test ./cmd/fid -run '^TestIssueTransitionJSONRejectsUnsatisfiedClaimPrecondition$' -count=1` - expected FAIL - `FAIL` (terminated `exit_code: 1`).
- GREEN: `go test ./cmd/fid -run '^TestIssueTransitionJSONRejectsUnsatisfiedClaimPrecondition$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.057s` (terminated `exit_code: 0`).
- RED: `go test ./cmd/fid -run '^TestIssueCreateRejectsRemovedOwnerFlagsWithJSONFailure$' -count=1` - expected FAIL - `FAIL` (terminated `exit_code: 1`).
- GREEN: `go test ./cmd/fid -run '^TestIssueCreateRejectsRemovedOwnerFlagsWithJSONFailure$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.048s` (terminated `exit_code: 0`).
- `go test ./cmd/fid ./internal/issues -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/issues\t4.469s` (terminated `exit_code: 0`)
- `go test ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/... -run '^(TestV01ClosureJourneyAndRebuildParity|TestV02aWorkflowAuthoringExitScenarioAndRebuildParity)$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.962s` (terminated `exit_code: 0`)
- Real entrypoint: in `/tmp/fid-task9.ucjhTX`, ran `/tmp/fid init --prefix FID`, `/tmp/fid --json issue create --kind task --title 'Task 9 entrypoint' --by developer`, `/tmp/fid --json issue update FID-548415 --title 'Task 9 updated' --expected-revision 1 --by developer`, then the same update with stale `--expected-revision 1` - PASS - final output was the exact JSON `state-conflict` envelope with `actual_revision:2` and `expected_revision:1` (terminated `exit_code: 1` for the deliberate stale rejection; preceding commands each terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 80.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t77.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t5.416s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

Fresh direct-recovery evidence captured on 2026-08-21:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.915s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `86d927da88f0fc070709a02ce39a3d68e560e2a6`

## Commit subject

- `feat(issue): migrate canonical mutation envelopes`
