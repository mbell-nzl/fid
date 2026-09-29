# Phase 23 - Task 15: Migration 0018 and V1 judgment removal - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 3
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Migration 0018 removes the V1 judgment tables, their dependent column and role rules, renames the generic Concern projection, and rebuilds the Issue view and reference triggers. The down migration follows the nine ordered steps in the DEC-264 correction to `schema.md` section 5; tests compare its schema-17 result with a fresh schema 17, preserve an agent Run, exercise a full rollback, and check round-trip integrity.
- Removed V1 Issue verification, completion, and Concern readers and CLI leaves. Retained Acceptance Criterion add/list, Issue detail/list, maintenance, journal, and historical schema-17 tests now assert the generic replacement or the explicit absence of removed behavior. Projection schema is 18.
- Requirement and Decision reference writers, plus requirement-subject Concern and Deviation writers, emit DEC-262 ActivityEffects against 0018's new trigger. A built-binary corpus produces the ten accepted activity names with their declared subjects and no extra name among its producing Events.
- Converged retained process fixtures, including the three removed leaves in `TestProcessExitContract`, to exit-2 rejection with no output or journal change. The full ordinary and tagged gates pass.

## Files created or modified

- `migrations/0018_p23_remove_v1_judgment.up.sql`, `migrations/0018_p23_remove_v1_judgment.down.sql`.
- `internal/db/migration0018_test.go`, `schema17_fixture_test.go`, `projection_meta.go`, and matching DB migration, view, trigger, and workflow-seed tests.
- `internal/issues/` Issue readers and retained tests; removed `ac_verify.go`, `capability_unavailable.go`, `concern_raise.go`, `concern_address.go`, and `concern_list.go`.
- `internal/workflow/` rules, registry, workflow, resume, and matching tests; `internal/maintenance/` checks and fixtures; `internal/stats/stats.go`.
- `internal/journal/` current Event catalogs, column/update/JSON validation, result projection, and retained replay tests; `internal/export/` archive semantics and projection-schema handling.
- `internal/review/` requirement, decision, Concern, Deviation producers and retained tests.
- `cmd/fid/` AC and Concern command trees and retained contract tests; `cmd/fid/integration/` closure, reference, judgment, process-exit, and new `p23_activity_inventory_test.go` journeys.
- `project-docs/docs/phases/phase-23/tasks/task-15-done.md`.

## Documentation updated

- This report records the accepted DEC-264 correction, exact removal inventory, retained-test provenance, all required gate results, and measured VAL-079 times. No governing spec, decision, QA, ledger, or control file was edited by this Developer.

## Validation Contract Coverage

- VAL-004: PASS - `TestMigration0018RoundTripIntegrity` runs `foreign_key_check` and `integrity_check` after 17→18→17→18; `TestMigration0018DeletesOnlyLegacyNonIssueReferences` and `TestMigration0018CurrentReferenceGuards` check populated schema 18; `TestP23ConcernDeviationRealProcessJourney` calls `p23AssertRebuiltProjectionIntegrity` after a real CLI rebuild and a superseding Concern. Schema-17 checks are also in `internal/db/review_judgment_migration_test.go`.
- VAL-055: PASS - `TestP23TenNameActivityInventoryFromProducingLeaves` executes the ten producers through the built binary, checks each Event's exact subject/count, then queries their `activity_log` names and compares the sorted distinct set to exactly ten names. `TestP23RequirementReferenceAndLinkRealProcess`, `TestP23DecisionReferenceProtectedLifecycleRealProcess`, and `TestP23ConcernDeviationRealProcessJourney` separately check requirement and decision subjects, including requirement-subject Concern disposal and Deviation.
- VAL-061: PASS - `TestLatestSchemaHasNoV1JudgmentTables` inspects removed tables, indexes, the two implicitly dropped kind triggers, columns, rules, and new objects in schema 18. `TestMigration0018IssuesEnrichedColumnsAndSubjectCount` compares the measured `issues_enriched` column list before and after 0018; it is identical, so the removed view-column set is empty. `TestProcessExitContract` and `TestV01ClosureJourneyAndRebuildParity` exercise removed leaves through the built binary; `TestLiveLeafCommandsMatchContractInventories` checks the command tree. The closed Event maps at `internal/journal/current_event_catalog.go:23` and `internal/export/archive_semantics_catalog.go:32` contain the replacement families; the three declared absence searches below terminated exit 1 with no matches.
- VAL-062: PASS - `TestMigration0018DeletesOnlyLegacyNonIssueReferences` seeds both Issue and non-Issue links and activity at schema 17, then checks the exact surviving Issue rows at 18. `TestMigration0018IssuesEnrichedColumnsAndSubjectCount` proves only open, Issue-subject generic Concerns count. `TestMigration0018CurrentReferenceGuards` runs valid Issue/Requirement/Decision activity and Issue links plus invalid endpoint-type, missing-source, missing-target, invalid-relation, missing-reference, and invalid-type inserts against all five rebuilt triggers.
- VAL-063: PASS - `TestV01ClosureJourneyAndRebuildParity` calls `fid maintenance check` on a populated workspace before and after rebuild and compares the results. `TestCheckErrorsReportsAllSeededStructuralViolationFamilies` names every retained check family; `TestCheckAcceptsPopulatedGenericReferencesWithoutStructuralErrors` verifies valid generic references produce no structural errors.
- VAL-079: PASS - final tagged integration: `240.222s` (limit 900s); final `cmd/fid` ordinary package: `198.174s` (limit 600s); `internal/workflow`: `101.597s` on the uncached first ordinary pass, then cached on the final pass (limit 600s). The first `cmd/fid` pass was `307.101s`. No required gate timed out; none reached 80% of its budget.

VAL-061 inventory, section 1 row by row (a retained replacement row names its executing process/projection test; the task-15 removed part names the absence proof):

| Inventory row | Final evidence |
|---|---|
| Envelope Obligation and `Obligations[0]` | `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess` reads four typed Obligation evaluations from one Envelope; no single-slot result survives. |
| Baseline CHECK domains | `TestMigration0018RoundTripIntegrity` and `TestMigration0018DownMatchesFreshSchema17` exercise the accepted migration chain and exact schema-17 rollback. |
| Claim review mode | `TestP23FindingDispositionRealProcessJourney` acquires review and mutation Claims and rejects addressing under review mode. |
| Evidence dimensions | `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess` uses gate, Review, and Human Outcome Evidence in one real-process workspace. |
| Review-cycle budget | `TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild` checks the replacement budget projection across rebuild. |
| Stop and Condition vocabulary | `TestP23DecisionReferenceProtectedLifecycleRealProcess` exercises protected Decision denial and resolution against new Stop facts. |
| AC verify and role projection | `TestProcessExitContract/AC_verification` returns exit 2 with empty output and unchanged journal; `TestAcceptanceCriteriaRemainOrderedWithoutRoleVerifications` retains add/list/detail parity and proves `ac_verifications` absent. |
| V1 Concern raise/address | `TestProcessExitContract/concern_raise` and `/concern_address` reject with no output or Event; `TestP23ConcernDeviationRealProcessJourney` records and disposes generic Concerns. |
| Issue detail/list V1 readers | `TestGetReturnsCentralizedDetailShape`, `TestGetReturnsEnrichedConcernFields`, and `TestListFiltersAndOrdering` read generic Issue-subject Concerns without V1 completion/verification fields. |
| Role-bound rules and roles parameter | `TestLatestSchemaHasNoV1JudgmentTables` counts zero rows for all three retired rules; `TestRetiredJudgmentRulesAreNotImplemented` and `TestRegistryEqualsV01Definitions` check the runtime registry. |
| Requirement sync package/Event | `TestV01ClosureJourneyAndRebuildParity` rejects `requirement sync` through the binary without journal change; the second declared absence search proves no production `requirement.synced` or `last_synced_hash` reference. |
| Requirement link endpoints | `TestMigration0018DeletesOnlyLegacyNonIssueReferences` removes seeded non-Issue link rows; `TestMigration0018CurrentReferenceGuards` rejects a new non-Issue endpoint; `TestP23RequirementReferenceAndLinkRealProcess` uses `requirement_issue_links`. |
| V1 maintenance checks | `TestCheckErrorsReportsAllSeededStructuralViolationFamilies` pins the retained check-name set; `TestV01ClosureJourneyAndRebuildParity` executes maintenance on the latest projection. |
| `agent_runs.completion_id` | `TestLatestSchemaHasNoV1JudgmentTables` finds no column; `TestMigration0018DownPreservesAgentRunAndFullRollback` restores a null column with the rest of the Run intact on down. |
| View and reference triggers | `TestMigration0018IssuesEnrichedColumnsAndSubjectCount` measures identical view columns and Issue-only count; `TestMigration0018CurrentReferenceGuards` probes all five rebuilt triggers in both directions. |
| Latest-schema DB expectations | `TestLatestSchemaHasNoV1JudgmentTables` enumerates five removed V1 tables, old indexes, both implicit kind guards, three removed rules, removed columns, and the replacement objects. |
| Journal catalogs and replay | `TestV01ClosureJourneyAndRebuildParity` snapshots and compares replayed projection tables; `TestProjectEventPopulatesPhaseTwoRepresentativeTablesViaJournalEffects` remains pinned to schema 2; the current and archive catalog maps cited under VAL-061 contain the renamed generic Concern effects. |
| Four pinned fail-closed conformance fixtures | `TestProcessExitContract`, `TestV01ClosureJourneyAndRebuildParity`, `TestP23RequirementReferenceAndLinkRealProcess`, and `TestP23ConcernDeviationRealProcessJourney` reject removed calls without state change and exercise replacements. |
| V1 judgment unavailable guard | `TestV1JudgmentProjectionIsAbsentAfterReplacement` proves removed writer tables absent and generic `concerns` present; the binary removed-leaf rows prove rejection. |
| Capability advertisement | Remains task 17's declared matrix boundary; `TestLiveLeafCommandsMatchContractInventories` checks the task-15 live tree without changing advertisement. |
| Resume findings placeholder | `TestP23FindingDispositionRealProcessJourney` reads real Finding states and next actions from `fid resume`. |
| Export projection schema | `TestCreateIsByteIdenticalWithCompleteDeterministicInventory` and the ordinary export suite pass against projection schema 18; `internal/db/projection_meta.go` supplies version 18. |
| Removed ActivityEvent names | `TestP23TenNameActivityInventoryFromProducingLeaves` compares exactly ten replacement names from their Events; the closed catalogs at the two VAL-061 paths omit the V1 activity producers. |

Retained-test provenance was checked by introducing commits, not filenames: AC `59f02df`; Issue detail `c62c6aa`, list `2a281f0`, Issue fixtures `f851b6a`; maintenance `d2a087d`; workflow registry `799acf0`; DB reference and kind triggers `457f480`/`e38ae45`; historical journal fixture `5b1467f`; tagged process exit `1d89aa4` and closure journey `ec72fda`. The accepted task-15 removal displaced their V1 assertions. Schema-17 and schema-2 fixtures remain pinned to their historical versions; current readers and process rows test the generic replacement. No earlier task number was used as an ownership inference.

Activity name inventory from `TestP23TenNameActivityInventoryFromProducingLeaves`:

| Producing leaf | Activity name | Subject |
|---|---|---|
| `requirement record-reference` | `requirement-referenced` | `requirement:REQ-activity` |
| `requirement link` | `requirement-linked` | `requirement:REQ-activity` |
| `decision record-reference` | `decision-referenced` | `decision:DREF-activity` |
| `review record` | `review-recorded` | Issue |
| `finding address` | `finding-addressed` | Issue |
| `review record` resolving a Finding | `finding-resolved` | Issue |
| `human-outcome record` | `human-outcome-recorded` | Issue |
| `concern record` | `concern-recorded` | Issue, and requirement in the separate subject test |
| `concern dispose` | `concern-disposed` | Issue, and requirement in the separate subject test |
| `deviation record` | `deviation-recorded` | Issue, and requirement in the separate subject test |

## Tests run

- Attempt 1 `go test -run '^TestLatestSchemaHasNoV1JudgmentTables$' -count=1 ./internal/db` - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	0.080s` (terminated `exit_code: 0`); `go test -run '^TestMigration0018DownRestoresEmptyV1Tables$' -count=1 ./internal/db` - FAIL - `FAIL` (terminated `exit_code: 1`), because 0018 down was absent.
- Attempt 1 `go test -count=1 ./internal/db ./internal/issues ./internal/workflow ./internal/maintenance ./internal/review` - FAIL - `FAIL` (terminated `exit_code: 1`); maintenance still queried dropped `task_completions`. The 04:21 externally interrupted child had no observed termination and is not PASS evidence.
- Attempt 2 exact focused migration command before down SQL - FAIL - `FAIL` (terminated `exit_code: 1`), missing restored V1 tables; after an incomplete trigger copy - FAIL - `FAIL` (terminated `exit_code: 1`), `near "CREATE": syntax error`. The extractor then copied through `END;`. Attempt 2's report claimed two focused PASS rows, but independent report review could not verify their termination evidence; they are not used to close this task.
- Attempt 2 touched-package command - FAIL - `FAIL` (terminated `exit_code: 1`); it exposed retained Issue, maintenance, DB, and archive-schema readers. Attempt 2 compile probe - FAIL - `FAIL` (terminated `exit_code: 1`) on removed AC verification APIs. `TestConcernSubjectAndDispositionMatrix` passed in its focused run; the attempt ended as material progress at checkpoint `dea2320a2e0c5a0d87f2305f0b726e14ddcacf47`.
- Attempt 3 initial `go test -run '^(TestLatestSchemaHasNoV1JudgmentTables|TestMigration0018Down.*)$' -count=1 ./internal/db` - FAIL - `FAIL` (terminated `exit_code: 1`), retained seed parity still referenced `RequiredRoles`. The unchanged focused command after convergence - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	0.340s` (terminated `exit_code: 0`); later expanded migration focus - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	0.647s` (terminated `exit_code: 0`).
- `go test -run '^TestMigration0018(DeletesOnlyLegacyNonIssueReferences|CurrentReferenceGuards|IssuesEnrichedColumnsAndSubjectCount)$' -count=1 ./internal/db` - FAIL - `FAIL` (terminated `exit_code: 1`), first new fixture reused `created_event_id` and hit trigger-order precedence; a second identical run - FAIL - `FAIL` (terminated `exit_code: 1`) on endpoint precedence. Corrected only fixture IDs/input isolation; the unchanged command then - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	0.167s` (terminated `exit_code: 0`).
- `go test -tags integration -run '^TestP23ConcernDeviationRealProcessJourney$' -count=1 ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`) with an invalid fixture activation after adding an unsupported action; removing that extra action yielded PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.377s` (terminated `exit_code: 0`). The earlier tagged Concern journey failure on unresolved requirement ID was repaired by recording the Requirement reference through the CLI.
- `go test -tags integration -run '^TestP23(RequirementReferenceAndLinkRealProcess|DecisionReferenceProtectedLifecycleRealProcess|ConcernDeviationRealProcessJourney)$' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.622s` (terminated `exit_code: 0`).
- `go test -tags integration -run '^TestP23TenNameActivityInventoryFromProducingLeaves$' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.310s` (terminated `exit_code: 0`).
- Touched ordinary race tier `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/db ./internal/issues ./internal/maintenance ./internal/journal ./internal/workflow ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	254.927s` (terminated `exit_code: 0`); touched tagged focus with `-tags integration -race` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	9.855s` (terminated `exit_code: 0`).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestProcessExitContract$' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.758s` (terminated `exit_code: 0`).
- `go build -o /tmp/fid-task15 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`); direct `/tmp/fid-task15 --json issue ac verify` - expected rejection - `exit 2 (no output)` (terminated `exit_code: 2`). The built-binary `TestProcessExitContract` additionally checks no journal change for each removed leaf.

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. These are the final-tree, byte-identical task commands in task-spec order. Earlier failed attempts remain below and under Deviations.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`).
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)			60.5%` (terminated `exit_code: 0`).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	240.222s` (terminated `exit_code: 0`).
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	25.273s` (terminated `exit_code: 0`); benchmark line: `BenchmarkRebuildWorkspace100kEvents-6         1 22558290358 ns/op 22558 ms/rebuild`.
- `rg -n --glob '!*_test.go' "task_completions|ac_verifications" cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`).
- `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`).
- `rg -n --glob '!*_test.go' "has_completion_report|concerns_addressed" cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`).

Earlier required-command attempts retained: `gofmt -l cmd internal` printed `internal/journal/current_event_columns.go` (terminated `exit_code: 0`, semantic FAIL because the listed file was unformatted); formatting it made the identical command silent. `golangci-lint run ./...` first printed `* unused: 1` (terminated `exit_code: 1`); the five findings were two unchecked row closes, one switch style, one struct conversion, and an unused rule helper. In-scope fixes gave `0 issues.` (terminated `exit_code: 0`). The first full tagged gate printed `FAIL` (terminated `exit_code: 1`) after `TestProcessExitContract` expected unavailable-capability exit 5 for three removed leaves; the repaired, identical gate printed the 240.222s PASS row above.

## Deviations from spec

- DEC-264 resolved the 04:34 protected-contract stop by correcting only `schema.md` section 5 under the governing authority. The original contradiction was the 0018 down duplication of the `concerns` column; task 15 copied the accepted correction and proved schema-17 fidelity, retained Run state, full rollback, and 17→18→17→18 integrity. No 0018 up statement, migration 0001–0017, command, or protected definition changed.
- The 04:21 attempt-1 interruption had no observed child termination; no PASS was inferred. The independent report review after attempt 2 returned `OUTCOME RECOVERY_REQUIRED` (terminated `exit_code: 2`) because two focused PASS rows lacked verifiable `exit_code: 0` termination evidence. This report preserves that finding and uses attempt-3 observed focused and full-gate results for closure. It makes no current report-control claim.
- The final-tree format, lint, and tagged-gate first failures remain immutable. Root causes and exact repairs are recorded under Exact quality-gate results; all three commands were rerun byte-for-byte. The tagged failure came from retained process-exit rows introduced at `1d89aa4`, whose old unavailable-capability assertion was displaced by task-15's accepted leaf removal. The rows now assert exit 2, empty output, and no Event; the unrelated `maintenance vacuum` unavailable row remains unchanged.
- Historical schema-17 trigger tests remain at 17 and the Phase-2 journal representative fixture remains at schema 2. Current-schema tests prove the replacement independently; this preserves the original historical subjects without executing obsolete V1 expectations against schema 18.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
