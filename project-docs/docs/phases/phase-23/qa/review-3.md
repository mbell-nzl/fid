# Phase 23 - Mid-Phase Review 3 (Task 2)

Reviewer: qa-reviewer
Date: 2026-09-20
Tasks covered: phase-23/task-2
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: `project-docs/docs/phases/phase-23/qa/review-2.md`

## Verdict

PASS

Review-2's remaining MF-01 clause is resolved. The repaired command-surface
matrix sends every one of the 24 forbidden typed-field combinations through
`envelope create`; each case requires a typed `invalid-argument` error, empty
stdout, and byte-identical journal and SQLite projection snapshots. Fresh
domain, command-surface, portable-integration, and report-lint evidence all
passed. No Must-fix finding remains in this Task 2 re-review window.

## Scope Reviewed

- Task spec: `project-docs/docs/phases/phase-23/tasks/task-2.md`.
- Completion report: `project-docs/docs/phases/phase-23/tasks/task-2-done.md`.
- Validation contract: `project-docs/docs/phases/phase-23/qa/validation-contract.md`.
- QA plan: sections 4 (EVID-07 through EVID-09), 7 (QA-08 through QA-10,
  QA-34, and QA-39), and 8.
- Prior reviews: `project-docs/docs/phases/phase-23/qa/review-1.md` and
  `project-docs/docs/phases/phase-23/qa/review-2.md`.
- Repair commits: `d71618aae4d8f3f56fee98468d2400a896a7a760`,
  `39317283fa9fc6b0087e4ad84d4df7733f064ff7`,
  `13cf819c0b96bb3a540654fa718f61fc5f6621f3`,
  `c9afe9016eeb5ff783f3aabe5fa68847b3f69be2`, and
  `29fff2599cb27d1352fb4d7baf057d7aa4afc6b5`.
- Changed test and report files: `internal/workflow/envelope_test.go`,
  `cmd/fid/envelope_test.go`, and `task-2-done.md`. For the review-2 repair,
  `c9afe90` changes only `cmd/fid/envelope_test.go` and `task-2-done.md`, and
  `29fff25` changes only the report. `git diff --check 1a4db05..HEAD` passed.
  The only worktree dirt is the pre-existing Runner-owned
  `project-docs/docs/phases/phase-23/run-ledger.md`.

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-006, VAL-007, VAL-008, VAL-009, VAL-010

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-006 | task-2 | ON TRACK | Fresh domain matrix covers Obligation 0/1/16/17, all three accepted typed shapes, 24 forbidden typed-field shapes, and protected-Decision 0/16/17 plus duplicate class, duplicate ordinal, and unguarded action. Fresh command matrix proves accepted 1/16 Obligation and 0/16 Decision bounds, while its rejection cases prove no durable effect for 0/17 and every forbidden typed-field/protected-Decision case. |
| VAL-007 | task-2 | ON TRACK | Fresh command matrix exercises duplicate Obligation ID, second gate-evidence, duplicate review profile, duplicate outcome class, and non-primary subject. Each returns `invalid-argument`, writes no stdout, and retains byte-identical journal/database snapshots. |
| VAL-008 | task-2 | ON TRACK | The fresh domain matrix includes ordered mixed-state evaluation with independent absence, validity/freshness, and Artifact-availability dimensions and all reasons; it also includes the qualifying valid/fresh/available state. |
| VAL-009 | task-2 | ON TRACK / DEFERRED | `TestEnvelopeActivationEvaluationIncludesEveryObligation` remains in the passing domain matrix for this task's blocking half. The all-satisfied gate/review/Human state remains expressly owned by task 10. |
| VAL-010 | task-2 | ON TRACK | The qualifying state proves gate Evidence is `available` and metadata-sufficient review/Human rows are exactly `not-applicable`; the mixed-state test preserves the gate-Evidence `missing` state. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 2 required-gate record | `task-2-done.md` against `task-2.md` | PASS | The report retains all six exact task commands, including the mandatory tagged integration gate and pre/post dedicated workflow gate, with output tokens and observed terminations. The review-2 exact rerun is recorded, as are both prior focused red attempts and their byte-identical green reruns. Reported local package times remain within the 600s ordinary-package and 900s tagged-integration budgets; hosted CI is not cited as evidence. |
| Domain authoring/evaluation matrix | `go test -race -timeout 30m ./internal/workflow/ -run '^(TestEnvelopeCreateAcceptsTypedObligationSet|TestEnvelopeCreateRejectsInvalidTypedObligationAndProtectedDecisionAuthoring|TestEnvelopeCreatePersistsEveryTypedObligationAndProtectedDecision|TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults|TestEnvelopeActivationEvaluationIncludesEveryObligation|TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix|TestEnvelopeAuthoringAcceptedAndForbiddenTypedFieldCombinations|TestEnvelopeAuthoringProtectedDecisionMatrix|TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/internal/workflow 3.738s` (terminated `exit_code: 0`). |
| Command authoring/persistence matrix | `go test -race -timeout 30m ./cmd/fid/ -run '^(TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions|TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries|TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange|TestEnvelopeCreateCommandRejectsEveryForbiddenTypedFieldCombinationWithoutPersistenceChange)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid 7.795s` (terminated `exit_code: 0`). The latter test table has all 24 forbidden typed-field cases and compares raw SQLite and journal bytes before/after every command rejection. |
| Task 2 portable integration | `go test -tags integration -race -timeout 30m ./cmd/fid/integration/ -run '^TestP22(CurrentEventFamiliesLiveReplayValidation|CurrentEventClosedValuesAndConditionalEffects|PortableFakeAndFIDRejectUnknownEffects)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration 325.990s` (terminated `exit_code: 0`). |
| Report/spec mechanical control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-2-done.md project-docs/docs/phases/phase-23/tasks/task-2.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). It confirms all required exact commands, the resolving `c9afe90` hash and subject, and a clean committed report. |
| Repair-commit identities | `git rev-parse --verify <reported-hash>^{commit}` | PASS | All five reported repair/report hashes resolve to commits (terminated `exit_code: 0`). The two behavior commits carry `Refs: phase-23/task-2`; their subjects are `test(workflow): complete envelope QA evidence matrix` and `test(workflow): complete envelope field rejection matrix`. |

## Per-Task Findings

### Task 2 - Typed multi-Obligation Envelope contract

#### Must-fix

- None.

#### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-01: complete authoring/evaluation matrix with no-Event/projection invariance | task-2 | RESOLVED | `cmd/fid/envelope_test.go:224-319` now runs all 24 forbidden gate-evidence, review, and human-outcome typed-field cases through `envelope create`. Every case checks non-nil command error, empty stdout, `"code":"invalid-argument"`, and `snapshotWorkspacePersistence` equality. That snapshot reads the raw database and journal bytes (`cmd/fid/workspace_fixture_test.go:250-264`). The fresh named command matrix, domain matrix, portable integration, and report lint all passed above. |

## Repair Plan

Not required: no Must-fix finding remains.

## Blocker Candidates

- None. No expected-class closed-registry rejection occurred in the fresh matrices. Consequently there is no out-of-Task-2 cause to record for escalation and no protected authority is implicated.

## Existing Escalation References

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-23/qa/review-3.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
```
