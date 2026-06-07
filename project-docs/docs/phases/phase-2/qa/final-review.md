# Phase 2 Final Review - Widen The Schema Vertically

Reviewer: QA Lead
Date: 2026-06-07

## Verdict

PASS - deterministic phase evidence is complete for VAL-001..VAL-017, with no
blocking findings and no open escalations. Human verification remains pending in
`project-docs/docs/phases/phase-2/qa/human-checklist.md`.

This review consolidates existing evidence from task completion reports and the
mid-phase review. Per the requested scope, it does not introduce new command
output or alter task reports.

## Scope And Evidence

- Phase brief: `project-docs/docs/current/phase_brief.md`
- Current decisions: `project-docs/docs/current/decisions_active.md`
- Current progress: `project-docs/docs/current/progress_current.md`
- Current escalations: `project-docs/docs/current/tl_escalations_current.md`
- Quality plan: `project-docs/specs/08_QUALITY_PLAN.md`
- Validation contract:
  `project-docs/docs/phases/phase-2/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-2/qa/qa-plan.md`
- Mid-phase review: `project-docs/docs/phases/phase-2/qa/review-3.md`
- Task specs and completion reports:
  `project-docs/docs/phases/phase-2/tasks/task-1.md`
  through `task-4.md`, and `task-1-completion.md` through
  `task-4-completion.md`

## Phase-Level Gate Results

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | Task 4 completion records `PASS gofmt -l cmd internal`; Tasks 1-3 also recorded PASS. |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Task 4 completion records `PASS golangci-lint run ./...`; Tasks 1-3 also recorded PASS. |
| Unit/race coverage | `go test -race -coverprofile=coverage.out ./...` | PASS | Task 4 completion records PASS; this covers the phase unit/race package suite with coverage. |
| Coverage report | `go tool cover -func=coverage.out` | PASS | Task 4 completion records `internal/journal` 80.6%, `internal/issues` 85.1%, `cmd/fid` 70.5%, total 72.6%. |
| Integration | `go test -tags integration -race ./...` | PASS | Task 4 completion records PASS as required phase-end evidence. Tasks 1-3 recorded N/A/NOT REQUIRED only where no integration-tag tests applied. |
| End-to-end | N/A | N/A | `AGENTS.md` records end-to-end tests as N/A until an accepted harness exists. Held-out scenarios were not read. |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | Task 4 completion records PASS as the phase-end performance spot-check. |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Task 4 completion records PASS; Task 2 also recorded PASS for its user-facing CLI change. |
| Rebuild smoke | `/tmp/fid init`; `/tmp/fid rebuild --by developer` | PASS | Task 4 completion records both commands passed in `/tmp/fid-task4-smoke.WhoMby`. |
| Adapter import graph | `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/...` | PASS | Task 4 completion records no `cmd/fid` import under `internal/...`. |
| Adapter source scan | `rg -n "github.com/mbell-nzl/fid/cmd/fid\|cmd/fid" internal` | PASS | Task 4 completion records no matches. |

## Validation Contract Outcomes

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead | Task 1 completion records migration `0002` introspection for all remaining v0.1 base/stub tables, representative columns, FKs, and static CHECKs; review-3 carried it ON TRACK. |
| VAL-002 | PASS | QA Lead | Task 1 completion records `TestMigration0002DownLeavesPhaseOneSchemaAndUpRestoresPhaseTwoTables`, preserving `projection_meta`/`roles`/`issues` and restoring Task 1 tables; review-3 carried it ON TRACK. |
| VAL-003 | PASS | QA Lead | Task 1 completion records `TestProjectEventPopulatesPhaseTwoRepresentativeTablesViaJournalEffects` in `internal/journal`, covering journal-effect population and read-back equality for representative new table groups; review-3 spot check passed the targeted test. |
| VAL-004 | PASS | QA Lead | Task 1 completion records `ProjectionSchemaVersion=2`, stale version `ProjectionRebuild`, current version `ProjectionFresh`, and `TestEnsureProjectionFreshRebuildsStoredPriorProjectionSchemaVersion`; review-3 carried it ON TRACK. |
| VAL-005 | PASS | QA Lead | Task 2 completion records `internal/db/issues_enriched_test.go` proving open-concern count/flag for open, addressed-only, and absent-concern states; review-3 carried it ON TRACK. |
| VAL-006 | PASS | QA Lead | Task 2 completion records latest `RESUME`, latest `BLOCK`, and no-triage coverage in `internal/db` and `internal/issues`; review-3 carried it ON TRACK. |
| VAL-007 | PASS | QA Lead | Task 2 completion records package-level `internal/issues` tests plus `cmd/fid` JSON/rebuild parity for enriched fields, zero/nonzero states, and delete-DB rebuild parity; review-3 spot check passed targeted `internal/issues` tests. |
| VAL-008 | PASS | QA Lead | Task 2 completion records `EXPLAIN QUERY PLAN` assertions for `idx_concerns_issue_status` and `idx_triage_issue_created`, with no base-table scan fallback; review-3 carried it ON TRACK. |
| VAL-009 | PASS | QA Lead | Task 3 completion records `TestIssueParentKindTriggerRejectsInvalidHierarchiesAndAllowsValidRows` covering invalid hierarchy RAISE cases and valid inserts; review-3 carried it ON TRACK. |
| VAL-010 | PASS | QA Lead | Task 3 completion records `TestIssueKindGuardTriggersRejectWrongKindsAndAllowValidRows`, including `task_completions` task-or-subtask resolution and valid/invalid states; review-3 carried it ON TRACK. |
| VAL-011 | PASS | QA Lead | Task 3 completion records `TestEpicReferenceKindGuardTriggersRejectWrongKindsAndAllowValidRows`, including non-epic rejection, valid epic refs, and `decisions.epic_id = NULL`; review-3 carried it ON TRACK. |
| VAL-012 | PASS | QA Lead | Task 3 completion records `TestIssueKindImmutabilityTriggerRejectsKindChangesAndAllowsOtherUpdates`; review-3 carried it ON TRACK. |
| VAL-013 | PASS | QA Lead | Task 4 completion records raw-SQL dangling-ref rejection and resolvable acceptance for `claims`, `links` endpoints, and `activity_log`, plus a valid journal-projection control with triggers active. |
| VAL-014 | PASS | QA Lead | Task 4 completion records `TestLinkScopeAndRelationTriggersRejectInvalidCombosAndAllowValidRows`, rejecting `decision`/`escalation` endpoints and incompatible `implements issue->issue`, while accepting `implements`, `blocks`, and `relates_to` valid combinations. |
| VAL-015 | PASS | QA Lead | Task 4 completion records `TestClaimScopeTriggersRejectUnsupportedRefTypeAndAllowIssueRef`, rejecting `qa_review` and accepting a resolvable issue claim. |
| VAL-016 | PASS | QA Lead | Task 4 completion records `internal/journal/phase2_rebuild_test.go` replaying journal events that touch every new Phase 2 table, snapshotting schema objects plus projected rows, deleting the DB, rebuilding, and asserting exact parity; `/tmp/fid init` and `/tmp/fid rebuild --by developer` smoke also passed. |
| VAL-017 | PASS | QA Lead | Task 4 completion records import-graph and source-scan evidence with no `internal/...` to `cmd/fid` dependency, and notes enriched-read logic remains in `internal/issues.Get` while `cmd/fid/issue.go` delegates. |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| PRD-OPS-003 / PRD-OPS-004 | MET | VAL-004 and VAL-016 prove stale projection rebuild and deterministic delete-DB rebuild over the widened schema. |
| PRD-ISSUE-006 | MET | VAL-007 proves `fid issue get --format json` returns enriched issue state. |
| PRD-DEV-007 | MET | VAL-005 and VAL-007 prove concern-derived state is computed and surfaced. |
| PRD-OUT-001 / PRD-OUT-004 | MET | VAL-007 proves JSON output preserves enriched fields before and after rebuild. |
| TRD §4.1 / schema inventory | MET | VAL-001 proves the remaining v0.1 base and stub tables; VAL-005..006 prove `issues_enriched`. |
| TRD §4.3 / projection freshness | MET | VAL-004 proves prior-version rebuild and current-version freshness. |
| TRD §4.5 / SQL backstops | MET | VAL-009..015 prove raw-SQL invalid writes RAISE and valid writes succeed. |
| TRD §5 / numbered migrations | MET | VAL-001..002 prove migration `0002`; Tasks 2-4 report migrations `0003`..`0005`. |
| TRD SCHEMA §14.1, §14.2, §14.3, §14.4, §14.7 | MET | VAL-009..015 prove kind, polymorphic-ref, link, claim, and kind-immutability triggers. |
| TRD-ARCH-001 / TRD-ARCH-002 | MET | VAL-003 and VAL-016 prove journal/projection reachability and rebuild determinism. |
| TRD-ARCH-005 | MET | VAL-017 proves adapter boundary. |
| TRD-ARCH-006 / TRD-SEC-004 | MET | VAL-009..015 prove trigger backstops through raw-SQL reject/accept evidence. |
| TRD-ARCH-008 / TRD-PERF-004 | MET | VAL-008 proves required enriched-view indexes are used; Task 4 records benchmark command PASS as phase-end spot-check. |
| Design `fid issue get` / `05_DESIGN_PLAN.md` §6.1 / `06_DESIGN_PLAN_CLI.md` | MET | VAL-007 proves enriched JSON fields are populated from `issues_enriched` and survive rebuild. |
| Deferred Phase 2 scope | CARRY-FORWARD | Sequence-vs-scheduling trigger and `activity_log` projection-only guard / `activity` op remain recorded in `decisions_active.md` as CF-004 and CF-005. |

## Test Quality Review

- Red/green evidence is recorded in all task completion reports.
- Binary and conditional assertions include both states: stale/fresh
  projection, concern/no-concern, triage states, trigger reject/accept controls,
  and valid journal-applied trigger controls.
- VAL-003 has package-level `internal/journal` evidence, satisfying the
  Phase-1 DEC-007 carry-in.
- VAL-007 has package-level `internal/issues` evidence plus `cmd/fid` CLI JSON
  evidence, satisfying the COV-001/COV-002 carry-in.
- Final reported coverage remains at or above the relevant thresholds:
  `internal/journal` 80.6%, `internal/issues` 85.1%, `cmd/fid` 70.5%.

## Deterministic Checklist Execution

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | PASS | Task 1 completion records exact required gates as PASS and conditional integration/build disposition. |
| QA-02 | PASS | Task 1 completion records all required Task 1 tables, representative columns, FKs, and static CHECK evidence. |
| QA-03 | PASS | Task 1 completion records `TestMigration0002DownLeavesPhaseOneSchemaAndUpRestoresPhaseTwoTables`. |
| QA-04 | PASS | Task 1 completion records `internal/journal` journal-effect population/read-back evidence; review-3 spot check passed. |
| QA-05 | PASS | Task 1 completion records stale/current freshness evidence and prior-version rebuild. |
| QA-06 | PASS | Task 2 completion records exact required gates and required build command as PASS. |
| QA-07 | PASS | Task 2 completion records concern count evidence for open, addressed-only, and absent states. |
| QA-08 | PASS | Task 2 completion records latest `RESUME`, latest `BLOCK`, and no triage. |
| QA-09 | PASS | Task 2 completion records `internal/issues` tests, CLI JSON evidence, zero/nonzero states, and rebuild parity. |
| QA-10 | PASS | Task 2 completion records query-plan index assertions for both required indexes. |
| QA-11 | PASS | Task 3 completion records exact gates and the `task_completions` task-or-subtask resolution. |
| QA-12 | PASS | Task 3 completion records invalid and valid parent/kind hierarchy trigger evidence. |
| QA-13 | PASS | Task 3 completion records wrong-kind rejection and valid-kind insertion for ACs, completions, and epic metadata. |
| QA-14 | PASS | Task 3 completion records non-epic rejection, valid epic refs, and NULL `decisions.epic_id` acceptance. |
| QA-15 | PASS | Task 3 completion records rejected kind update and accepted non-kind update. |
| QA-16 | PASS | `review-3.md` exists, verdict CONTINUE, VAL-001..012 ON TRACK, no blockers. |
| QA-17 | PASS | Task 4 completion records exact required gates plus integration, benchmark, and build as PASS. |
| QA-18 | PASS | Task 4 completion records dangling rejection and resolvable acceptance for links, claims, and activity_log. |
| QA-19 | PASS | Task 4 completion records endpoint-scope/relation reject and accept cases. |
| QA-20 | PASS | Task 4 completion records claim `issue` acceptance and unsupported ref-type rejection. |
| QA-21 | PASS | Phase-level command evidence is represented by Task 4 completion: format, lint, race coverage, coverage report, integration, benchmark, build, and rebuild smoke all PASS. |
| QA-22 | PASS | Task 4 coverage report records `internal/journal` 80.6%, `internal/issues` 85.1%, and `cmd/fid` 70.5%. |
| QA-23 | PASS | Task 4 completion records full-schema journal rebuild determinism with all triggers active and delete-DB parity. |
| QA-24 | PASS | Task 4 completion records import graph no `cmd/fid` import and `rg` source scan no matches. |
| QA-25 | PASS | This final review contains requirement compliance, security review, performance spot-check, documentation review, and risk review sections. |
| QA-26 | PASS | This final review records explicit PASS outcomes for VAL-001..VAL-017 with cited evidence. |

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: YES.
  `qa-plan.md` section 6 records validation for every referenced command, file,
  fixture, and human checklist assumption.
- Invalid or unvalidated command, test scope, fixture, artifact, or human
  checklist step found: none.
- Required command placeholders found: none.

## Mid-Phase Review Disposition

| Review | Verdict | Blocking Fixed? | Improvement Disposition | ESC References |
|---|---|---|---|---|
| `review-3.md` | CONTINUE | N/A - no blockers | N/A - no improvement findings | none |

Review-3 covered Tasks 1-3 and VAL-001..VAL-012. Its ON TRACK outcomes are
carried forward as final PASS outcomes above, and Task 4 completes VAL-013..017.

## Escalation Review

- Phase-start baseline: no open ESC IDs, per `qa-plan.md` and review-3.
- Current state: `tl_escalations_current.md` still records `_None open._`.
- New open ESC relative to baseline: none.

## Security Review

PASS. Phase 2 security-relevant scope is the SQLite backstop layer. Task 3 and
Task 4 completion reports record raw-SQL trigger tests proving invalid writes
RAISE for kind hierarchy, kind guards, epic refs, kind immutability,
polymorphic refs, link scope/relation compatibility, and claim scope. Valid
control writes and valid journal-applied rows are also recorded, reducing the
risk of false-positive trigger rejection.

No dependency changes were recorded in Tasks 1-4.

## Performance Spot-Check

PASS. VAL-008 records deterministic query-plan evidence for index use on
`idx_concerns_issue_status` and `idx_triage_issue_created`. Task 4 records
`go test -bench . -benchmem ./...` PASS as the phase-end benchmark spot-check.
Performance tuning remains out of scope until Phase 3 / later roadmap work.

## Documentation Review

- Completion reports present for all tasks: YES.
- QA plan and mid-phase review present: YES.
- Final review and human checklist generated by this review: YES.
- Task reports record no dependency-budget changes and no untriaged spec
  deviations.
- Spec-sync/carry-forward items for Tech Lead wrap-up:
  - `task_completions` kind guard resolution is task-or-subtask per
    `03_TRD_SCHEMA.md` §14.1, recorded by Task 3.
  - CF-004 sequence-vs-scheduling trigger deferral remains active.
  - CF-005 `activity_log` projection-only guard / `activity` op deferral
    remains active.
  - `phase_brief.md` and `progress_current.md` still describe Phase 2 as
    planned; Tech Lead wrap-up should refresh current-context docs after final
    QA and human verification.

## Human Checklist Readiness

Generated at
`project-docs/docs/phases/phase-2/qa/human-checklist.md`.

- Contains only `(human)` items: YES.
- Uses one observable per row: YES.
- Human-only `VAL-NNN` assertions mapped: N/A. VAL-001..VAL-017 are
  deterministic in this phase.
- Product/perception/dogfooding risks to call out: human should dogfood the
  `/tmp/fid` binary and decide whether enriched `fid issue get --format json`
  is useful for real agent workflow consumption.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | READY | Validation contract, QA plan, task reports, mid-phase review, final review, and human checklist exist after this review. |
| Automated gates | READY | Task 4 records phase-level format, lint, race coverage, coverage report, integration, benchmark, and build PASS. |
| Held-out scenarios | DEFERRED | AGENTS and task specs record held-out scenarios as N/A/not to be read until an accepted harness exists. |
| Release readiness | DEFERRED | Public release readiness scans and threat-model work are not in Phase 2 scope. |

## Findings

### Blocking

- None.

### Major

- None.

### Improvement

- None.

## Unresolved Risks And Carry-Forward Candidates

- Human checklist is pending; deterministic PASS does not replace product-owner
  judgment or dogfooding.
- CF-003 remains active from Phase 1 and is not exercised by Phase 2.
- CF-004 and CF-005 remain active Phase 2 deferrals for future lifecycle and
  activity-log work.
- Current-context docs still say Phase 2 is planned; refresh during Tech Lead
  phase wrap-up after human verification.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-2/qa/final-review.md
Human checklist written: project-docs/docs/phases/phase-2/qa/human-checklist.md
VAL outcomes: 17 PASS, 0 FAIL, 0 BLOCKED, 0 CARRY-FORWARD
Blocking: 0  Major: 0  Improvement: 0  Escalations: none
VERDICT: PASS
```
