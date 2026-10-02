# Phase 23 - Mid-Phase Review 2 (Task 2)

Reviewer: qa-reviewer
Date: 2026-09-20
Tasks covered: phase-23/task-2
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: `project-docs/docs/phases/phase-23/qa/review-1.md`
Re-review of: review-1

## Verdict

REPAIR

MF-01 is only partially resolved. The repair supplies the requested domain
matrix and command-level boundary/request-negative matrix, and the fresh named
tests and portable integration command pass. However, the command-level test
does not send every forbidden per-type field combination through `envelope
create`, so it does not prove unchanged journal and projection bytes for those
rejections. This is bounded, decision-free Task 2 test-and-report work. Task 4
must not dispatch before a fresh independent re-review closes this remaining
MF-01 clause.

## Scope Reviewed

- Task spec: `project-docs/docs/phases/phase-23/tasks/task-2.md`.
- Completion report: `project-docs/docs/phases/phase-23/tasks/task-2-done.md`.
- Validation contract: `project-docs/docs/phases/phase-23/qa/validation-contract.md`.
- QA plan: sections 4 (EVID-07 through EVID-09), 7 (QA-08 through QA-10,
  QA-34, and QA-39), and 8.
- Prior finding: `project-docs/docs/phases/phase-23/qa/review-1.md`, MF-01.
- Repair commits: `d71618aae4d8f3f56fee98468d2400a896a7a760`,
  `39317283fa9fc6b0087e4ad84d4df7733f064ff7`, and
  `13cf819c0b96bb3a540654fa718f61fc5f6621f3`.
- Changed files: `internal/workflow/envelope_test.go`,
  `cmd/fid/envelope_test.go`, and `task-2-done.md` only. `git diff --check
  22ec707..HEAD` was clean. The original Task 2 failures remain recorded in
  the amended report as failed attempts before their byte-identical green
  reruns.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 1
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-006, VAL-007, VAL-008, VAL-009, VAL-010

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-006 | task-2 | NOT ON TRACK | `TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix` covers 0/1/16/17 Obligations and `TestEnvelopeAuthoringProtectedDecisionMatrix` covers 0/16/17 Decisions plus duplicate class, ordinal, and unguarded action. The fresh domain matrix passed, but its forbidden type-field cases do not make a command-level no-Event/projection-invariance assertion. |
| VAL-007 | task-2 | ON TRACK | `TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange` executes all five request-level negatives and compares complete journal and database bytes before and after each. Fresh command-surface matrix passed. |
| VAL-008 | task-2 | ON TRACK | The named domain matrix passed. The mixed state test proves ID ordering, independent dimensions, and complete unsatisfied reasons; the qualifying-state test proves valid/fresh/available. |
| VAL-009 | task-2 | ON TRACK / DEFERRED | `TestEnvelopeActivationEvaluationIncludesEveryObligation` remains green for Task 2's blocking half. The all-satisfied gate/review/Human case is expressly owned by task 10. |
| VAL-010 | task-2 | ON TRACK | The qualifying test proves gate Evidence is `available` and metadata-sufficient review/Human Obligations are exactly `not-applicable`; the mixed-state test retains gate `missing`. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 2 exact gates recorded | `task-2-done.md` against task-2 exact commands | PASS | The report records all six required commands, including tagged integration and pre/post `go test -race -timeout 30m -cover ./internal/workflow/`, each PASS with a token and observed `exit_code: 0`. Its repair rerun records the full same six-command list. The reported package times remain inside the 600s/900s budgets and use local evidence, not hosted CI. |
| Domain named matrix | `go test -race -timeout 30m ./internal/workflow/ -run '^(TestEnvelopeCreateAcceptsTypedObligationSet|TestEnvelopeCreateRejectsInvalidTypedObligationAndProtectedDecisionAuthoring|TestEnvelopeCreatePersistsEveryTypedObligationAndProtectedDecision|TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults|TestEnvelopeActivationEvaluationIncludesEveryObligation|TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix|TestEnvelopeAuthoringAcceptedAndForbiddenTypedFieldCombinations|TestEnvelopeAuthoringProtectedDecisionMatrix|TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/internal/workflow 3.472s` (terminated `exit_code: 0`). |
| Command-surface named matrix | `go test -race -timeout 30m ./cmd/fid/ -run '^(TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions|TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries|TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid 7.163s` (terminated `exit_code: 0`). |
| Portable integration | `go test -tags integration -race -timeout 30m ./cmd/fid/integration/ -run '^TestP22(CurrentEventFamiliesLiveReplayValidation|CurrentEventClosedValuesAndConditionalEffects|PortableFakeAndFIDRejectUnknownEffects)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration 329.380s` (terminated `exit_code: 0`). |
| Report/spec mechanical control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-2-done.md project-docs/docs/phases/phase-23/tasks/task-2.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Reported repair commit | `git rev-parse --verify d71618aae4d8f3f56fee98468d2400a896a7a760^{commit}` | PASS | `d71618aae4d8f3f56fee98468d2400a896a7a760` (terminated `exit_code: 0`); subject and footer are `test(workflow): complete envelope QA evidence matrix` and `Refs: phase-23/task-2`. |

## Per-Task Findings

### Task 2 - Typed multi-Obligation Envelope contract

#### Must-fix

- **MF-01 (still open, narrowed)** [cmd/fid/envelope_test.go:156] The command-level rejection loop snapshots journal and projection bytes for 0/17 Obligation counts, the five VAL-007 negatives, and the protected-Decision negatives, but omits every forbidden per-type field combination exercised only by `internal/workflow/envelope_test.go:159`. Consequently the domain validator proves rejection, but the repair does not prove that each such rejected `envelope create` leaves both durable artifacts unchanged. Owning task: task-2. Repair boundary: extend only `cmd/fid/envelope_test.go` with table-driven `envelope create` cases for every forbidden gate-evidence, review, and human-outcome field combination, using the existing `snapshotWorkspacePersistence` before/after assertion; amend only the matching Task 2 report evidence. Required evidence: each forbidden type-field input returns the typed rejection and has byte-identical journal and database snapshots; rerun the affected named command-surface matrix and record token plus observed termination. Re-review scope: a fresh independent reviewer inspects the expanded command matrix and report, reruns both named matrices, the Task 2 portable integration command, and report lint, and rechecks every original MF-01 clause.

#### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-01: complete authoring/evaluation matrix with no-Event/projection invariance | task-2 | STILL OPEN | Fresh passing domain and command-surface matrices cover the stated boundaries, all five request negatives, protected-Decision negatives, qualifying valid/fresh/available, and whole-workspace invariance for the command cases present. Source inspection shows the 24 forbidden type-field combinations remain direct validator cases only, with no matching command/persistence cases. |

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-01 | task-2 | VAL-006 | `cmd/fid/envelope_test.go` and the matching evidence in `task-2-done.md` only; add command-level persistence-invariance coverage for all forbidden typed fields. | Every added case returns `invalid-argument`, emits no stdout, and has unchanged journal/database bytes; affected named matrix passes with observed termination. | Inspect the changed cases/report; rerun both named matrices, portable integration, and report lint; verify each original MF-01 clause. |

## Blocker Candidates

- None. No closed-registry rejection was encountered. The missing evidence is bounded Task 2 test coverage and requires no protected authority.

## Existing Escalation References

- None. ESC-23-002 is a resolved Task 1 recovery and does not cover MF-01.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-23/qa/review-2.md
Must-fix: 1  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
