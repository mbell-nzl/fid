# Phase 18 - Mid-Phase Review 1 (Tasks 1-3)

Reviewer: qa-reviewer
Date: 2026-07-19
Tasks covered: phase-18/task-1 through phase-18/task-3

## Verdict

BLOCK

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-18/tasks/task-1.md`
  - `project-docs/docs/phases/phase-18/tasks/task-2.md`
  - `project-docs/docs/phases/phase-18/tasks/task-3.md`
- Completion reports:
  - `project-docs/docs/phases/phase-18/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-18/tasks/task-2-done.md`
  - `project-docs/docs/phases/phase-18/tasks/task-3-done.md`
- Validation contract and QA plan:
  - `project-docs/docs/phases/phase-18/qa/validation-contract.md`
  - `project-docs/docs/phases/phase-18/qa/qa-plan.md`
- Implementation commits:
  - `127f26806bd4a0baae622b277a508e51dde1aaf5` — `feat(workflow): add seeded workflow tables`
  - `799acf0943be467b9f00424f7cb0ee8f63267dd7` — `feat(workflow): load workflow registry from database`
  - `cf68153de4058f77d06451a1ceed3f2191658a77` — `feat(issues): resolve transitions from workflow registry`

## Summary

- Blocking findings: 1
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-008, VAL-009, VAL-012, VAL-014, VAL-015, VAL-016, VAL-020

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | `TestMigration0007CreatesWorkflowSchema` inspects all four tables, the three declared PKs, composite requirement FK, and the intentional absence of a requirement-table PK/UNIQUE. Migration 0007 matches those structural requirements. |
| VAL-002 | task-1, advanced by task-2 | ON TRACK | Task 1 deep-compares seeded state/category tuples to `FrozenDefinitions`; Task 2 independently deep-compares the loaded registry. |
| VAL-003 | task-1, advanced by task-2 | ON TRACK | Seed and registry parity tests cover all allowed and absent state pairs, including explicit non-cancelled-to-cancelled rows and no cancelled-to-cancelled row. |
| VAL-004 | task-1, advanced by task-2 | ON TRACK | Seed and registry parity tests compare edge rules, ordered roles, and parameter JSON against the frozen reference. |
| VAL-005 | task-1 | ON TRACK | `TestMigration0007ReservesGateObligationShapeWithoutSeedingReservedRules` inserts, loads, parses, and marshals the canonical gate JSON without loss. |
| VAL-008 | task-3 | ON TRACK | `TestTransitionRegistryMatchesFrozenV01OrderedStateMatrix` covers 81 task, 25 epic, and 16 subtask ordered pairs; focused source tests prove DB edge/rule/cancel/category use. Existing rule fixtures cover both blocked and satisfied branches for the developer, QA-lead, and epic rules. The report records built-CLI allowed, invalid, blocked, cancel, and re-cancel output. |
| VAL-009 | task-3 | ON TRACK | Missing, wrong-holder, stale, and live-caller entry cases are covered; the DB-category mutation test proves claim gating reads the loaded registry rather than the frozen map. |
| VAL-012 | task-1 | ON TRACK | Freshness decision test covers prior-version rebuild and current-version freshness; `TestEnsureProjectionFreshRebuildsStoredPriorProjectionSchemaVersion` independently reproduced PASS. |
| VAL-014 | task-2 | ON TRACK | `TestRegistryEqualsV01Definitions` deep-compares all loaded definitions with the frozen v0.1 definitions. |
| VAL-015 | task-2 | ON TRACK | `TestLoadSurfacesReservedRulesAsUnimplementedWithoutWriting` inserts both reserved names, enables SQLite `query_only`, and verifies structural unimplemented results. |
| VAL-016 | task-1 | ON TRACK | The gate-seam test queries each reserved name and requires a zero seeded-row count. |
| VAL-020 | task-3 | NOT ON TRACK | The same-event delete effect and immediate projected absence are tested, but the required rebuilt-projection absence is not tested for the released-claim fixture. See the Blocking finding. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 report hygiene | `project-docs/tools/report-lint.sh task-1-done.md task-1.md` | PASS | Required headings, all five exact commands, hash, subject, and committed/clean report verified. |
| Task 2 report hygiene | `project-docs/tools/report-lint.sh task-2-done.md task-2.md` | PASS | Required headings, all five exact commands, hash, subject, and committed/clean report verified. |
| Task 3 report hygiene | `project-docs/tools/report-lint.sh task-3-done.md task-3.md` | PASS | Required headings, all five exact commands, hash, subject, and committed/clean report verified. |
| Exact command coverage | Task-spec and QA-plan gate comparison | PASS | Each Task 1-3 report records `gofmt -l cmd internal`, `golangci-lint run ./...`, `go test -race -coverprofile=coverage.out ./...`, `go tool cover -func=coverage.out`, and `go test -tags integration -race ./...` as PASS. The active QA plan adds no further task-window command for Tasks 1-3. |
| Commit metadata and diff integrity | `git rev-parse --verify <hash>^{commit}`; `git show`; `git show --check` | PASS with one Improvement | All three implementation hashes resolve, subjects match reports, include their task `Refs` footer, and diffs pass whitespace checks. Task 2's separate report-finalization commit lacks its `Refs` footer. |
| Reproduced required command | `go test -tags integration -race ./...` | PASS | Exit 0; all packages passed (cached output). This is an exact required command for every reviewed task. |
| Reproduced focused evidence | `go test -race -count=1 ./internal/db ./internal/workflow ./internal/issues -run '^(TestMigration0007CreatesWorkflowSchema|TestMigration0007SeedsFrozenV01WorkflowDefinitions|TestMigration0007ReservesGateObligationShapeWithoutSeedingReservedRules|TestDecideProjectionFreshnessRequiresPhase18SchemaVersion|TestRegistryEqualsV01Definitions|TestLoadSurfacesReservedRulesAsUnimplementedWithoutWriting|TestTransitionRegistryMatchesFrozenV01OrderedStateMatrix|TestTransitionUsesDBRuleReferences|TestTransitionUsesDBCategoriesForClaimGating|TestTransitionRejectsCancelledEdgeMissingFromDBRegistry|TestTransitionIntoStartedRequiresLiveCallerClaim|TestTransitionIntoStartedRejectsWrongHolderAndStaleClaims|TestTransitionOutOfStartedReleasesCallerClaimInSameEvent|TestTransitionRequiresDeveloperVerificationsAndCompletionReport|TestTransitionRequiresQALeadVerificationForQAPass|TestTransitionRequiresDoneOrCancelledChildrenToWrapEpic)$'` | PASS | Non-cached run passed: `internal/db` 1.294s, `internal/workflow` 1.268s, `internal/issues` 4.084s. It confirms the available schema, seed, registry, transition-matrix, rule, and claim tests, while exposing the missing rebuild assertion described below. |
| Reproduced VAL-012 rebuild branch | `go test -race -count=1 ./internal/journal -run '^TestEnsureProjectionFreshRebuildsStoredPriorProjectionSchemaVersion$'` | PASS | Non-cached run passed in 1.361s. |

## Per-Task Findings

### Task 1 - Workflow tables migration + seed defaults + gate-obligation seam

#### Blocking

- None.

#### Improvement

- None.

#### Observations

- The additive migration is correctly isolated from `issues`: its schema test checks the exact four-table contract, while the seed parity test binds the static data to the frozen v0.1 reference. The projection freshness evidence includes both stale and already-current states.

### Task 2 - DB-backed workflow registry loader + v0.1 parity fixtures

#### Blocking

- None.

#### Improvement

- [commit `7edfff6c8a06631ade7691bff089395f2a20480f`] The `docs(task): finalize completion report for phase-18 task-2` commit has no `Refs: phase-18/task-2` footer, unlike the implementation and the other two report-finalization commits. Fix direction: include the phase/task reference on future phase-work documentation commits; preserve the existing task report content and implementation hash.

#### Observations

- `workflow.Load` is read-only and builds a detached registry. The `query_only` test injects both reserved names and verifies they are surfaced as unimplemented rather than evaluated or written.

### Task 3 - Transition-path data-source swap (behavior-preserving)

#### Blocking

- [internal/issues/transition_test.go:534] VAL-020 requires the caller claim to be absent both immediately and after rebuilding the projection. `TestTransitionOutOfStartedReleasesCallerClaimInSameEvent` proves the event's delete effect and immediate `Get()` result through line 648, but neither it nor another test rebuilds this released-claim fixture. The general rebuild test at line 17 rebuilds an entering-started transition and therefore cannot prove persistence of the release. The completion report nevertheless marks VAL-020 PASS. Fix direction: add a characterization test that exits a started state with a live caller claim, rebuilds the exact workspace/projection, and asserts the claim remains absent; then correct the completion evidence and rerun the Task 3 required gates before Task 4 starts.

#### Improvement

- None.

#### Observations

- The transition path loads the registry once inside the existing write lock, resolves edge/rule/category decisions from it, retains app-layer rule evaluation, and preserves the transition event's status/workflow preconditions and activity effect. The DB-mutation source tests and 122-pair structural matrix make the data-source swap observable.

## Escalations

- None. The VAL-020 gap is localized to Task 3's explicit acceptance criterion and has a direct, in-window test fix; it does not require a cross-task, architectural, product, security, or validation-contract decision.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-18/qa/review-1.md
Blocking: 1  Improvement: 1  Escalations: none
VERDICT: BLOCK
```
