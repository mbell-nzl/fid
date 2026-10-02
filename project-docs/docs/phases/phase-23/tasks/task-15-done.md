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

- This report records the accepted DEC-264 correction, exact removal inventory, retained-test provenance, all required gate results, and measured VAL-079 times. Review-11 MF-2's report-only repair adds per-row provenance and representative Git history checks. No governing spec, decision, QA, ledger, or control file was edited by this Developer.

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

Retained-test provenance is recorded below per converged assertion. An
introducing commit is the commit that introduced the named test or the
individual case before task 15; the cited `Refs` in that commit supplies its
historical task. Task 15's accepted change is migration 0018's removal of V1
judgment tables, role verification, completion, and Concern leaves, and its
generic Concern and reference replacements. Each row states the particular
part of that change that displaced its old assertion. Paths and lines refer to
the final task-15 source tree. All commands below are the byte-identical task
gate commands already recorded under Exact quality-gate results; this report
repair did not rerun them.

| Retained converged assertion (`file:line`) | Introducing commit and task | Accepted task-15 change that invalidated the old assertion | Executing command |
|---|---|---|---|
| `internal/issues/ac_test.go:587` AC add/list/detail parity and absent verification field, replacing the old verify row | `437d093`, phase-7/task-2 (old verify row); retained add/list subjects `59f02df`, phase-7/task-1 and `1806aa4`, phase-7/task-3 | `ac_verifications` and role verification were removed while AC add/list remain; the row now proves ordinal ordering and no `latest_verifications`. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/ac_test.go:518` AC add Event round trip | `59f02df`, phase-7/task-1 | AC add remains after removal of AC verify and its projection; the retained Event assertion runs against schema 18. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/ac_test.go:665` AC add missing-Issue error | `59f02df`, phase-7/task-1 | AC authoring remains, but the schema-18 fixture no longer has role verification tables. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/issue_ac_test.go:340` verify confirmation case, via helper at line 398 | `437d093`, phase-7/task-2 | The verify leaf was removed; its retained call now rejects as an unknown command with no journal append. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/issue_ac_test.go:348` latest-per-role case, via helper at line 398 | `437d093`, phase-7/task-2 | Removal of role verification replaces latest-per-role readback with unknown-command rejection and no journal append. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/issue_ac_test.go:357` fail-to-pass history case, via helper at line 398 | `437d093`, phase-7/task-2 | Removal of role verification replaces verification history with unknown-command rejection and no journal append. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/issue_ac_test.go:366` rebuild selection case, via helper at line 398 | `69e761d`, phase-7/task-2b | Removed `ac_verifications` cannot supply latest-per-role rebuild state; the retained call rejects without an Event. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/issue_ac_test.go:375` reject-before-append case, via helper at line 398 | `437d093`, phase-7/task-2 | Removal of the verify leaf makes the entire call reject as unknown command, with journal bytes unchanged. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/issue_ac_test.go:388` role-permission case, via helper at line 398 | `5293fd2`, phase-10/task-3 | Actor roles cannot enable the removed verify leaf; each retained call rejects without journal mutation. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/issue_ac_test.go:423` retained AC list Human/JSON case | `16f4f8d`, phase-20/task-10 | AC list remains, but role verification data is gone from its schema-18 read shape. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/retained_actor_request_contract_test.go:207` retained verify-leaf request case | `bc328f9`, phase-21/task-9 | Removal of the verify leaf displaces guarded capability-unavailable behavior; it now rejects as an unknown command without mutation. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/task9_contract_test.go:66` AC verify inventory case | `e51afe3`, phase-21/task-9 | The inventory now records the removed verify leaf instead of a fail-closed advertised leaf. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/task9_contract_test.go:72` Concern raise inventory case | `e51afe3`, phase-21/task-9 | The inventory now records the removed raise leaf instead of a fail-closed advertised leaf. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/detail_test.go:9` centralized Issue detail | `c62c6aa`, phase-4/task-4 | Detail no longer reads V1 completion and role verification; it retains AC and generic Issue-subject Concern fields. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/issue_test.go:963` enriched Concern fields | `00d1820`, phase-2/task-2 | Generic `concerns` replaces V1 Issue Concern rows and their completion links. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/issue_test.go:1010` open-Concern count | `00d1820`, phase-2/task-2 | The view now counts open generic Concerns with Issue subjects, excluding disposed and non-Issue rows. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/issue_test.go:1054` same-second Concern order | `c927de0`, phase-9/task-4 | Generic Concern readback replaces the V1 address/completion ordering fixture. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/list_test.go:71` Issue list and its fixture at line 479 | `2a281f0`, phase-4/task-5 | List retains filters and ordering but seeds generic Issue-subject Concerns instead of V1 completion/Concern rows. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/issue_list_test.go:11` CLI Issue list and its fixture at line 85 | `2a281f0`, phase-4/task-5 | CLI list retains ordering and filters; the fixture now uses generic Concerns because V1 completion and Concern tables are removed. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/issue_test.go:453` Issue get JSON and rebuild | `00d1820`, phase-2/task-2 | JSON/rebuild parity now reads generic Concern state and omits removed completion and AC verification fields. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/concern_list_test.go:9` filtered Concern list | `a172634`, phase-9/task-3 | V1 Concern list writer/reader is removed; retained list behavior uses generic Issue-subject Concerns. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/concern_list_test.go:21` same-second Concern list order | `c927de0`, phase-9/task-4 | Generic Concern rows replace the old address readback while preserving deterministic ordering. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/concern_test.go:20` raise case | `f35b3fe`, phase-9/task-1 | `concern raise` is removed; the retained call rejects without a new journal Event. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/concern_test.go:49` address case | `3ea2d4f`, phase-9/task-2 | `concern address` is removed; the retained call rejects without a new journal Event. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/concern_test.go:70` Issue-scoped V1 list case | `cdcc4f0`, phase-23/task-14 | Schema 18 removes the V1 Concern table; the retained case now exercises the generic list path. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/concern_test.go:116` generic list case | `cdcc4f0`, phase-23/task-14 | The generic table is renamed to `concerns`, and the V1 row source disappears. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/judgment_unavailable_test.go:9` projection absence case | `b441601`, phase-21/task-9 (old unavailable-writer test) | V1 writer tables are removed; the retained assertion checks absent V1 projection and present generic Concerns. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/transition_rules_test.go:8` removed Concern rule case | `3ea2d4f`, phase-9/task-2 (old `concerns_addressed` test) | The `concerns_addressed` rule is removed with the V1 Concern lifecycle; no seeded transition requires it. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/maintenance/check_test.go:15` structural finding names | `d2a087d`, phase-17/task-3 | Dropped V1 tables remove frozen-state and dependent-row findings; generic reference checks remain. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/maintenance/check_test.go:106` valid-reference check | `a690e49`, phase-21/task-21 (old frozen-state case) | Generic Requirement/Decision references are valid and V1 frozen-state presence is no longer reportable. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/maintenance/check_test.go:272` final-schema integrity case | `8435120`, phase-21/task-21 | The `task_completions` table and its frozen-state finding are removed; retained integrity classes remain asserted. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/maintenance/warnings_test.go:11` warning-state case | `4257bce`, phase-17/task-4 | The valid fixture omits V1 completion/Concern rows because schema 18 removes them; current warning states remain tested. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/maintenance/test_helpers_test.go:92` impossible and valid structural fixtures (valid rows at line 132) | `d2a087d`, phase-17/task-3 | Seeded V1 completion/Concern/Decision rows no longer exist at schema 18; the retained maintenance cases now seed generic references and Concerns. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/integration/exit_contract_test.go:25` `AC_verification` process-exit case | `1d89aa4`, phase-21/task-9 (individual old row at line 23) | Removed `issue ac verify` now yields exit 2, empty output, and unchanged journal instead of capability-unavailable exit 5. | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` |
| `cmd/fid/integration/exit_contract_test.go:26` `concern_raise` process-exit case | `1d89aa4`, phase-21/task-9 (individual old row at line 24) | Removed `concern raise` now yields exit 2, empty output, and unchanged journal instead of capability-unavailable exit 5. | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` |
| `cmd/fid/integration/exit_contract_test.go:27` `concern_address` process-exit case | `1d89aa4`, phase-21/task-9 (individual old row at line 25) | Removed `concern address` now yields exit 2, empty output, and unchanged journal instead of capability-unavailable exit 5. | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` |
| `cmd/fid/integration/e2e_v01_closure_test.go:14` closure and rebuild journey, including snapshot inventory at `cmd/fid/integration/e2e_v01_closure_helpers_test.go:437` | `ec72fda`, phase-20/task-3f | Removed V1 leaves reject without append; rebuild snapshots now include generic Concern/Requirement/Decision reference rows instead of removed V1 tables. | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` |
| `cmd/fid/integration/p23_concern_deviation_test.go:17` Concern/Deviation journey | `3a0aed9`, phase-23/task-8 | Generic Concern/Deviation writers now target renamed `concerns`; Requirement-subject activity and rebuilt-projection checks replace V1 assumptions. | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` |
| `cmd/fid/integration/p23_requirement_reference_test.go:10` Requirement reference journey | `0e4bb33`, phase-23/task-12 | New `requirement_refs` and rebuilt activity triggers replace removed V1 Requirement rows and sync Event. | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` |
| `cmd/fid/integration/p23_decision_reference_test.go:16` protected Decision reference journey | `828eb18`, phase-23/task-9 | New `decision_refs` and activity triggers replace removed V1 Decision rows. | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` |
| `cmd/fid/integration/p23_review_record_test.go:172` retained activity assertion helper | `d215e40`, phase-23/task-4 | Activity assertions now resolve Requirement/Decision subjects through the new reference tables and generic judgment subjects after V1 removal. | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` |
| `internal/db/migrate_v01_tables_test.go:15` latest table inventory | `b63a3cd`, phase-21/task-20 | Schema 18 removes five V1 tables, retains generic `concerns`, and drops `agent_runs.completion_id`; the inventory now asserts these facts. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/db/kind_triggers_test.go:281` Issue kind guards | `e38ae45`, phase-2/task-3 | The old Requirement/Decision kind guards are removed at 18; this historical kind-guard assertion is pinned to schema 17. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/db/kind_triggers_test.go:398` Epic reference kind guards | `e38ae45`, phase-2/task-3 | The V1 reference kind guards are removed at 18; this historical assertion is pinned to schema 17. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/db/ref_triggers_test.go:14` polymorphic reference trigger test (schema-17 fixture at line 245) | `457f480`, phase-2/task-4 | Schema-18 reference triggers resolve through generic references; this retained V1 link/activity assertion uses the schema-17 fixture. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/db/ref_triggers_test.go:72` link scope and relation trigger test (schema-17 fixture at line 245) | `457f480`, phase-2/task-4 | Schema-18 links no longer accept V1 Requirement endpoints; this retained V1 relation matrix uses the schema-17 fixture. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/db/ref_triggers_test.go:157` journal-projected link and activity test (schema-17 fixture at line 245) | `457f480`, phase-2/task-4 | Schema-18 link/activity guards resolve through new references; this retained V1 journal projection assertion uses the schema-17 fixture. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/db/workflow_seed_parity_test.go:274` historical workflow requirement fixture | `10d4d24`, phase-19/task-2 | V1 `RequiredRoles` is removed; schema-17 seed parity checks the old requirement shape while schema-18 tests cover the replacement. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/db/issues_enriched_test.go:89` view query-plan case | `00d1820`, phase-2/task-2 | The view counts generic Issue-subject Concerns; the retained plan fixture no longer seeds V1 completion/Concern Events. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/db/escalation_removal_test.go:59` active-schema version row | `644f30e`, phase-22/task-16 | Migration 0018 advances `ActiveMigrationVersion` from 17 to 18; the retained escalation absence case now expects version 18. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/journal/v01_tables_test.go:14` Phase-2 representative Event fixture | `5b1467f`, phase-2/task-1 | Removed current V1 Event effects cannot project at schema 18; this historical fixture remains pinned to schema 2. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/journal/current_event_task8_test.go:5` generic Concern disposal Event | `3a0aed9`, phase-23/task-8 | The generic Concern table is now `concerns`; the retained disposal shape and open guard use that table. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/workflow/registry_test.go:62` reserved-rule registry row | `799acf0`, phase-18/task-2 | Task-15 workflow convergence removes the legacy parameter JSON from this fixture; its two reserved rules remain unimplemented. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/workflow/rules_test.go:5` retired judgment rule row | `cde39fd`, phase-5/task-4 (old three rule tests) | `all_acs_verified`, `has_completion_report`, and `concerns_addressed` are removed; the replacement row asserts all three are unimplemented. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/workflow/workflow_test.go:58` edge-rule matrix | `cde39fd`, phase-5/task-4 | Default workflow edges no longer require the three removed judgment rules; retained edges stay checked. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/workflow/authoring_test.go:71` rule vocabulary and parameters | `10d4d24`, phase-19/task-2 | Removed role/completion/Concern rule names and parameters are no longer accepted in authored definitions. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/workflow/authoring_test.go:100` requirement-edge validation | `10d4d24`, phase-19/task-2 | The old judgment requirements disappear from accepted definition edges. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/workflow/define_test.go:16` atomic workflow definition Event | `50533e2`, phase-19/task-3 | Authored workflow fixture no longer includes removed judgment rules; atomic Event behavior stays asserted. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/issues/transition_test.go:368` reserved-rule transition case | `1f196a9`, phase-20/task-9a | The task-15 registry convergence drops legacy parameter JSON from this reserved-rule fixture; its unsatisfied denial remains asserted. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/contract_surface_test.go:128` live leaf inventory | `37b59d6`, phase-20/task-18 | The three V1 mutation leaves are absent from the command tree; retained leaves still match the contract inventory. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/capabilities_test.go:163` deferred-surface case | `401f609`, phase-21/task-3 | Removed V1 leaves reject as unknown commands while capability data remains the declared separate boundary. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/link_test.go:507` ambiguous same-ID endpoint case | `cf096da`, phase-20/task-11 | Requirement links leave generic Issue-link endpoints; the retained typed Issue-link assertion no longer seeds removed V1 Requirement rows. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/e2e_v01d_link_test.go:170` Requirement snapshot helper used by retained link journey | `e30e40e`, phase-12/task-4 | The removed V1 `requirements` table is replaced by `requirement_refs`; the snapshot now reads the new reference ID, source JSON, and status. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `cmd/fid/task18_completion_resume_test.go:315` complete-envelope case | `2313235`, phase-21/task-18 | Completion guard fixture drops V1 `task_completions`; current Envelope facts and no-mutation denial remain asserted. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/review/concern_task8_test.go:14` carried-forward Concern disposal | `3a0aed9`, phase-23/task-8 | The generic Concern table is renamed to `concerns`; the retained subject/disposition assertion reads it there. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |
| `internal/workflow/authoring_effects_test.go:58` retired requirement rejection | `556b4e6`, phase-21/task-20 | Task 15 removes the `RuleHasCompletionReport` constant with V1 completion rules; the retained rejection passes the literal retired rule name. | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` |

Representative history checks: `git blame` at `1d89aa4` identifies the three
old process-exit rows at lines 23-25; the task-15 rows are now at lines 25-27.
`git blame` before task 15 identifies the old AC verify test at
`internal/issues/ac_test.go:586` as `437d093`. `git show -s` confirms their
`Refs: phase-21/task-9` and `Refs: phase-7/task-2` trailers, respectively;
`git show 457f480:internal/db/ref_triggers_test.go` confirms the three
historical trigger tests; `git show -s 457f480` confirms its
`Refs: phase-2/task-4` trailer.
The three changed test files without retained-row provenance are new task-15
tests; the changed closure helper is named in its journey row above.

These entries use introducing-commit provenance rather than historical test
names as ownership. The schema-17 and schema-2 fixtures retain their original
subjects; current-schema assertions separately prove the replacement.

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
- `go test -run '^TestMigration0018(DeletesOnlyLegacyNonIssueReferences|CurrentReferenceGuards|IssuesEnrichedColumnsAndSubjectCount)$' -count=1 ./internal/db` - FAIL - `FAIL` (terminated `exit_code: 1`), first new fixture reused `created_event_id` and hit trigger-order precedence; a second identical run - FAIL - `FAIL` (terminated `exit_code: 1`) on endpoint precedence. Corrected only fixture IDs/input isolation; the unchanged command then - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	0.167s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23ConcernDeviationRealProcessJourney$' -count=1 ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`) with an invalid fixture activation after adding an unsupported action; removing that extra action yielded PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.377s` (terminated `exit_code: 0`). The earlier tagged Concern journey failure on unresolved requirement ID was repaired by recording the Requirement reference through the CLI.
- `go test -tags integration -run '^TestP23(RequirementReferenceAndLinkRealProcess|DecisionReferenceProtectedLifecycleRealProcess|ConcernDeviationRealProcessJourney)$' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.622s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23TenNameActivityInventoryFromProducingLeaves$' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.310s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/db ./internal/issues ./internal/maintenance ./internal/journal ./internal/workflow ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	254.927s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP23TenNameActivityInventoryFromProducingLeaves|TestP23ConcernDeviationRealProcessJourney|TestP23RequirementReferenceAndLinkRealProcess|TestP23DecisionReferenceProtectedLifecycleRealProcess|TestV01ClosureJourneyAndRebuildParity)$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	9.855s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestProcessExitContract$' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.758s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid-task15 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- Direct `/tmp/fid-task15 --json issue ac verify` - expected rejection - `exit 2 (no output)` (terminated `exit_code: 2`)
- Review-11 MF-2 report-only attempt 1 `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-15-done.md project-docs/docs/phases/phase-23/tasks/task-15.md` - FAIL - `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`); its sole finding was that the edited report had not yet been committed. The same command must run again after the report commit.

The built-binary `TestProcessExitContract` additionally checks no journal change for each removed leaf.

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. These are the final-tree, byte-identical task commands in task-spec order. Earlier failed attempts remain below and under Deviations.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)			60.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	240.222s` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	25.273s` (terminated `exit_code: 0`)

Benchmark output also recorded `BenchmarkRebuildWorkspace100kEvents-6 1 22558290358 ns/op 22558 ms/rebuild`.

- `rg -n --glob '!*_test.go' "task_completions|ac_verifications" cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n --glob '!*_test.go' "has_completion_report|concerns_addressed" cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)

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

- `0dd4c57849f8be7779e1c2e43cf7b14093c945f0`

## Commit subject

- `feat(task): complete phase 23 task 15 judgment removal`
