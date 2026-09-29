# Phase 23 - Task 15: Migration 0018 and V1 judgment removal - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: DEC-264 resolves the 0018-down contradiction and the rollback tests pass. The retained Issue tests still compile against V1 verification and completion APIs removed by this task; the full V1 convergence, activity inventory, schema probes, real-process probes, and nine-command gate remain unfinished.
- Recovery attempts: 2
- Next action: Continue task 15 from this committed tree. Converge the retained Issue, workflow, DB, journal, export, maintenance, CLI, and tagged fixtures against the accepted removal; add the VAL-004/055/061/062/063 probes; run focused red/green and the touched-package race tier; then run all nine byte-identical task commands in order on the final tree. Record every observed result, package wall times, the full per-row inventory, and the real CLI interaction before closing.
- Human attention required: no

## Summary of what was implemented

- Preserved the 04:21 interrupted implementation checkpoint. DEC-264 corrected the protected SQL authority after the 04:34 stop. Copied the nine ordered 0018-down steps from `schema.md` section 5 and its named migration definitions; the focused rollback now passes.
- Added the four DEC-264 rollback proofs in `TestMigration0018DownMatchesFreshSchema17`, `TestMigration0018DownPreservesAgentRunAndFullRollback`, and `TestMigration0018RoundTripIntegrity`, while retaining `TestMigration0018DownRestoresEmptyV1Tables`.
- Removed the V1 Issue mutation files, CLI leaves, workflow evaluators, and initial Issue detail readers. Converged initial maintenance reads and enabled archive semantic validation for projection schema 18. Added DEC-262 activity Effects for requirement, decision, and requirement-subject Concern/Deviation producers. These changes are partial and still need their full behavioral tests and fixture convergence.

## Files created or modified

- `migrations/0018_p23_remove_v1_judgment.up.sql`
- `internal/db/migration0018_test.go`, `internal/db/projection_meta.go`
- `internal/export/archive_semantics_catalog.go`
- `internal/journal/current_event_catalog.go`, `current_event_columns.go`, `current_event_json.go`, `current_event_result_projection.go`, `current_event_updates.go`, `current_event_validation.go`
- `internal/review/concern.go`, `concern_list.go`
- `internal/workflow/resume.go`, `resume_actions.go`
- `project-docs/docs/phases/phase-23/tasks/task-15-done.md`
- `migrations/0018_p23_remove_v1_judgment.down.sql`
- `internal/export/archive_semantics.go`; `internal/review/requirement.go`, `decision.go`, `deviation.go`, `concern_task8_test.go`
- `internal/issues/ac_list.go`, `detail.go`, `get.go`, `issue.go`, `transition_rules.go`; removed `ac_verify.go`, `capability_unavailable.go`, `concern_raise.go`, `concern_address.go`, `concern_list.go`
- `internal/workflow/rules.go`, `registry.go`, `workflow.go` and matching retained workflow tests
- `internal/maintenance/errors.go`, `warnings.go`, `operational.go`; `cmd/fid/issue_ac.go`, `concern.go`

## Documentation updated

- This report retains the attempt-1 failed evidence, the 04:21 interruption, the 04:34 protected-contract stop, DEC-264's resolution, and this incomplete continuation.

## Validation Contract Coverage

- VAL-004: INCOMPLETE - attempt-1 `TestMigration0018DownRestoresEmptyV1Tables` (`internal/db/migration0018_test.go:30`) failed because 0018 down was absent. The DEC-264 focused test now passes, including `TestMigration0018RoundTripIntegrity`; post-rebuild pragmas remain unproven.
- VAL-061: INCOMPLETE - `TestLatestSchemaHasNoV1JudgmentTables` (`internal/db/migration0018_test.go:8`) passes for five tables. The full per-row schema/CLI inventory and real-process absence probes remain pending; no source search is claimed as their substitute.
- VAL-062: INCOMPLETE - migration 0018 up remains preserved, but no seeded 17-to-18 deletion or complete trigger probe has terminated in this session.
- VAL-063: INCOMPLETE - the attempt-1 touched-package test reported `check dependent row kinds: query task_completions: SQL logic error: no such table: task_completions (1)` from `internal/maintenance`. This continuation removed that query and the stale-verification, telemetry-gap, and frozen-V1 checks; a populated real-process maintenance probe remains pending.
- VAL-055: INCOMPLETE - DEC-262 producers now append `activity_log/activity` Effects with requirement, decision, and requirement-subject references in the accepted code paths. Their tagged ten-name inventory and rejection probes remain pending.
- VAL-079: INCOMPLETE - the ordinary and tagged required gates and fresh-process benchmark have not run on a completed task-15 tree; `cmd/fid/integration`, `cmd/fid`, and `internal/workflow` wall times cannot be claimed.

VAL-061 per-row inventory status at this checkpoint (the full acceptance proof is pending):

| Inventory section 1 row | Current absence evidence |
|---|---|
| Envelope Obligation/`Obligations[0]` | Replacement owned by task 2; task-15 final projection/process check pending. |
| Baseline CHECK domains | Replacement owned by task 1; task-15 final schema check pending. |
| Claim `review` mode | Replacement owned by task 3; task-15 final process check pending. |
| Evidence dimensions | Replacement owned by tasks 6/10; task-15 final process check pending. |
| Review-cycle budget | Replacement owned by task 6; task-15 final process check pending. |
| Stop/Condition vocabulary | Replacement owned by tasks 6/9/10/11; task-15 final process check pending. |
| AC verify/role projection | CLI leaf and producer source removed in this checkpoint; retained-test convergence and built-binary rejection proof pending. |
| V1 Concern raise/address | CLI leaves and producers removed in this checkpoint; built-binary rejection proof pending. |
| Issue detail/list V1 readers | Initial detail reader removed in this checkpoint; retained fixture and projection evidence pending. |
| Role-bound workflow rules | Three evaluators and roles field removed in this checkpoint; registry/retained assertion proof pending. |
| Requirement sync package/Event | Prior task-12 replacement; required exact absence command and projection/process proof pending. |
| Requirement link endpoint model | Prior task-12 replacement; 0018 seeded link deletion probe pending. |
| Maintenance V1 checks | Initial V1 queries removed in this checkpoint; populated real-process check pending. |
| `agent_runs.completion_id` | 0018 up removes it; post-0018 column inventory pending. |
| View and ref triggers | 0018 up rebuilds them; column-list measurement and valid/invalid trigger probes pending. |
| Latest-schema DB expectations | Five removed tables proven by `TestLatestSchemaHasNoV1JudgmentTables`; all other DB objects/columns pending. |
| Journal catalogs and replay fixtures | Initial rename present from checkpoint; full replay and Event absence proof pending. |
| Four pinned fail-closed conformance fixtures | Retained fixture convergence and tagged process proof pending. |
| V1 judgment unavailable guard | Source removed in this checkpoint; retained-test and built-process proof pending. |
| Capability advertisement | Task 17 retained boundary; task-15 checks replacement readiness, final matrix remains pending. |
| Resume findings placeholder | Task-14 replacement; task-15 final process check pending. |
| Export projection schema | Manifest emits active schema 18, semantic validator now admits 18; archive validation regression pending. |
| Removed ActivityEvent names | Closed catalog/source removal check and ten-name tagged inventory pending. |

The inventory rows above are recovery tracking, not absence PASS claims. The final report must replace every pending item with a schema query, Event-catalog query, or real-process probe and its observed output.

## Tests run

- `go test -run '^TestLatestSchemaHasNoV1JudgmentTables$' -count=1 ./internal/db` - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	0.080s` (terminated `exit_code: 0`)
- `go test -run '^TestMigration0018DownRestoresEmptyV1Tables$' -count=1 ./internal/db` - FAIL - `FAIL` (terminated `exit_code: 1`); all five restored tables were absent because no 0018 down file exists.
- `go test -count=1 ./internal/db ./internal/issues ./internal/workflow ./internal/maintenance ./internal/review` - FAIL - `FAIL` (terminated `exit_code: 1`); the initial failure exposed the missing rollback and retained V1 readers/tests against the latest schema. This is immutable failed-attempt evidence, not a task gate.
- Prior externally interrupted child: no typed result, no observed command termination, and no accepted green result.
- Attempt 2 `go test -run '^(TestLatestSchemaHasNoV1JudgmentTables|TestMigration0018Down.*)$' -count=1 ./internal/db` before 0018 down - FAIL - `FAIL` (terminated `exit_code: 1`); `TestMigration0018DownRestoresEmptyV1Tables` found all five V1 tables missing.
- Attempt 2 same focused command after copying 0018 down with a truncated trigger statement - FAIL - `FAIL` (terminated `exit_code: 1`); the diagnostic was `SQL logic error: near "CREATE": syntax error (1)`. The copy extractor stopped at a semicolon inside the trigger body. It was repaired to copy through `END;` without changing the SQL authority.
- Attempt 2 same focused command after the extractor repair - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/db\t0.136s` (terminated `exit_code: 0`).
- `go test -run '^(TestLatestSchemaHasNoV1JudgmentTables|TestMigration0018Down.*|TestMigration0018RoundTripIntegrity)$' -count=1 ./internal/db` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/db\t0.458s` (terminated `exit_code: 0`); this includes the three new DEC-264 proof tests.
- `go test -count=1 ./internal/db ./internal/issues ./internal/workflow ./internal/maintenance ./internal/review ./internal/journal ./internal/export ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); the final command line was `FAIL` after DB latest-schema, Issue V1 reader, maintenance, and schema-18 archive semantic failures. This was a diagnostic touched-package attempt before the later source edits.
- `go test -run '^$' ./internal/issues ./internal/workflow ./internal/maintenance ./internal/review ./internal/export ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); retained `internal/issues/ac_test.go` still compiles against the removed verification API. This is the current in-scope recovery condition.
- `go test -run '^TestConcernSubjectAndDispositionMatrix$' -count=1 ./internal/review` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/review\t0.006s` (terminated `exit_code: 0`); the fixture now uses the renamed generic `concerns` table.

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. No required gate was started on a final tree because in-scope implementation and retained-test convergence remain incomplete. Each unrun command is recorded as blocked without a fabricated output token or exit status. The ordinary and tagged gate times required by VAL-079 remain unobserved.

- `gofmt -l cmd internal` - BLOCKED - not run on a final tree; Issue tests still fail to compile.
- `golangci-lint run ./...` - BLOCKED - not run on a final tree; Issue tests still fail to compile.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - not run on a final tree; Issue tests still fail to compile.
- `go tool cover -func=coverage.out` - BLOCKED - not run; no final-tree coverage profile exists.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - not run on a final tree; Issue tests still fail to compile.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - BLOCKED - not run on a final tree; replay corpus convergence remains pending.
- `rg -n --glob '!*_test.go' "task_completions|ac_verifications" cmd internal` - BLOCKED - exact command not yet run as a task gate on the final tree.
- `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` - BLOCKED - exact command not yet run as a task gate on the final tree.
- `rg -n --glob '!*_test.go' "has_completion_report|concerns_addressed" cmd internal` - BLOCKED - exact command not yet run as a task gate on the final tree.

## Deviations from spec

- DEC-264 resolved the 04:34 protected-contract stop by correcting `schema.md` section 5. This continuation used its nine ordered steps and added the four required proof obligations. The earlier contradiction is retained as historical failed evidence; no task-scope SQL definition was changed.
- The first focused rollback attempt ended `FAIL` (terminated `exit_code: 1`), with `restored V1 table "judgment_concerns" is missing` as a verbatim diagnostic. The root cause was the absent 0018 down file; repair was withheld until DEC-264 resolved the governing definition. The later byte-identical focused command passed with `ok  \tgithub.com/mbell-nzl/fid/internal/db\t0.136s` (terminated `exit_code: 0`).
- The first committed-report lint attempt returned `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`) because the Recovery attempts field included prose instead of an integer. The field is now `1`; the interruption explanation remains under Workflow outcome. The unchanged lint command is rerun against this correction.

## Concerns, technical debt, or notes for the Tech Lead

- DEC-264 supplies the post-0016 definition. No new protected decision is requested.
- The preserved partial migration removes V1 tables before the Issue and maintenance readers, CLI leaves, retained fixtures, and activity producers have converged. The touched-package command exposed these as unfinished task-15 work, not green evidence or another task's assignment.
- The external interruption consumed no no-progress budget per the Runner's TASK-RECOVERY binding. No provenance-based ownership claim is made for the retained tests before their introducing commits are checked during continuation.
- The interrupted child returned no typed result or command termination; none is treated as green evidence.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `6f74c77c6b19ec41f12a7e5a4bcfb75773fb43e0`

## Commit subject

- `chore(task): checkpoint phase 23 task 15 removal convergence`
