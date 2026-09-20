# Phase 23 - Mid-Phase Review 1 (Tasks 1-3)

Reviewer: qa-reviewer
Date: 2026-09-20
Tasks covered: phase-23/task-1, phase-23/task-2, phase-23/task-3
Review type: INITIAL
Prior review: none

## Verdict

REPAIR

The Task 2 implementation and its named tests pass, but its tests do not cover
the complete accepted authoring and evaluation matrices. This is bounded,
decision-free test-and-report work in Task 2. Task 4 must not be dispatched
until a fresh independent re-review closes MF-01.

## Scope Reviewed

- Task specs: `project-docs/docs/phases/phase-23/tasks/task-{1,2,3}.md`.
- Completion reports: `project-docs/docs/phases/phase-23/tasks/task-{1,2,3}-done.md`.
- Validation contract: `project-docs/docs/phases/phase-23/qa/validation-contract.md`.
- QA plan: section 4 evidence rows EVID-01 through EVID-11, section 7
  QA-03 through QA-11, section 8, and QA-34/QA-39.
- Accepted authorities: `contracts.md` sections 2-3, `public-interface.md`
  sections 2, 4-5, and `schema.md` sections 1-3 and 6.
- Commits: `730128cd8842e2943c9ab2b21f8f2b154148c10a`,
  `710db23c946891133f1967e4ec317a17d43a6038`, and
  `35944c72aba0f2c0b918bdb494b9fa9487d27142`. Each resolves to a commit,
  carries its stated `Refs: phase-23/task-N` footer, and has a clean diff check.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 1
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-23-002 (resolved; Task 1 historical recovery only)
- Covered validation assertions: VAL-001 through VAL-013

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | Fresh `internal/db` named migration suite passed; it compares the exact object set and accepted DDL. |
| VAL-002 | task-1 | ON TRACK | The same fresh suite exercises both rejected discriminator inserts and both accepted inserts. |
| VAL-003 | task-1 | ON TRACK | Fresh journal fence test rejects every currently unavailable Event and compares projection state. The qualifying live-owner parity state is explicitly deferred to tasks 4, 5, 8, 9, 10, and 12. |
| VAL-004 | task-1 | ON TRACK | Named migration test checks clean foreign-key and integrity pragmas at schema 17; Task 1's retained local rebuild smoke records clean rebuild pragmas. Schema-18 evidence remains with tasks 15-16. |
| VAL-005 | task-1 | ON TRACK | Fresh named migration test preserves seeded receipt Evidence and its Artifact link, with null new references and a clean foreign-key check. |
| VAL-006 | task-2 | NOT ON TRACK | MF-01: only a three-Obligation positive sample and three negatives are tested; 1/16 acceptance, 17 rejection, all field combinations, and the 0/16/17 protected-Decision matrix are absent. |
| VAL-007 | task-2 | NOT ON TRACK | MF-01: the named negative test covers only duplicate `review_profile` and one unguarded decision action, not all five whole-request negatives with unchanged persistence. |
| VAL-008 | task-2 | NOT ON TRACK | MF-01: the evaluator test supplies only unsatisfied/missing, invalid/stale, and absent states; it has no qualifying valid/fresh/available state. |
| VAL-009 | task-2 | ON TRACK / DEFERRED | The blocking half is exercised through multi-Obligation activation. The all-satisfied gate + review + Human state is contractually owned by task 10, so it is not a Task 2 defect. |
| VAL-010 | task-2 | NOT ON TRACK | MF-01: `not-applicable` is exercised for the new types and `missing` for absent gate Evidence, but no available gate-Evidence state proves the third value does not substitute for `available`. |
| VAL-011 | task-3 | ON TRACK | Fresh CLI named tests cover review acquisition, competing mutation rejection, and acquisition after release in both directions. |
| VAL-012 | task-3 | ON TRACK | Fresh domain and real-process tests cover review dirty-adoption rejection for both recovery and abandonment, no Event, clean review handling, and the mutating helper path. Existing Phase 22 real-process recovery coverage remains part of the passing tagged gate. |
| VAL-013 | task-3 | ON TRACK | Fresh real-process lifecycle test covers close, clean abandon, and clean recovery with the successor Claim still in review mode. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 report/spec coverage | `report-lint.sh task-1-done.md task-1.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). All six exact Task 1 gates, including tagged integration and the rebuild benchmark, are recorded with tokens and observed terminations. Earlier failed attempts and DEC-238 repair evidence remain recorded. |
| Task 2 report/spec coverage | `report-lint.sh task-2-done.md task-2.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). All six required commands, including the required pre/post `internal/workflow` command, are recorded with tokens and observed terminations. This mechanical result does not make the incomplete VAL matrix sufficient. |
| Task 3 report/spec coverage | `report-lint.sh task-3-done.md task-3.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). All five exact gates are recorded with tokens and observed terminations. |
| Task 1 named tests | `go test -race -count=1 -run '^(TestMigration0017CreatesExactlyTheAcceptedObjects|TestMigration0017ObligationRegistrationDiscriminator|TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks|TestMigration0017PreservesReceiptEvidenceAndArtifactLinksInPlace)$' ./internal/db` | PASS | `ok github.com/mbell-nzl/fid/internal/db 4.541s` (terminated `exit_code: 0`). |
| Task 1 fence test | `go test -race -count=1 -run '^TestProjectEventFenceRejectsEveryUnavailablePhase23EventWithoutProjectionChange$' ./internal/journal` | PASS | `ok github.com/mbell-nzl/fid/internal/journal 2.206s` (terminated `exit_code: 0`). |
| Task 2 domain tests | `go test -race -count=1 -run '^(TestEnvelopeCreateAcceptsTypedObligationSet|TestEnvelopeCreateRejectsInvalidTypedObligationAndProtectedDecisionAuthoring|TestEnvelopeCreatePersistsEveryTypedObligationAndProtectedDecision|TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults|TestEnvelopeActivationEvaluationIncludesEveryObligation)$' ./internal/workflow` | PASS | `ok github.com/mbell-nzl/fid/internal/workflow 3.383s` (terminated `exit_code: 0`). Source inspection of this passing set establishes MF-01. |
| Task 2 command surface | `go test -race -count=1 -run '^TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions$' ./cmd/fid` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 6.417s` (terminated `exit_code: 0`). |
| Task 2 portable integration | `go test -tags integration -race -count=1 -run '^TestP22(CurrentEventFamiliesLiveReplayValidation|CurrentEventClosedValuesAndConditionalEffects|PortableFakeAndFIDRejectUnknownEffects)$' ./cmd/fid/integration` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid/integration 325.047s` (terminated `exit_code: 0`). |
| Task 3 CLI tests | `go test -race -count=1 -run '^(TestClaimAcquireAcceptsReviewMode|TestClaimAcquireKeepsOneActiveClaimAcrossModes)$' ./cmd/fid` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 11.570s` (terminated `exit_code: 0`). |
| Task 3 domain test | `go test -race -count=1 -run '^TestReviewClaimRejectsDirtyPathAdoption$' ./internal/workflow` | PASS | `ok github.com/mbell-nzl/fid/internal/workflow 3.362s` (terminated `exit_code: 0`). |
| Task 3 real-process tests | `go test -tags integration -race -count=1 -run '^TestP23Review(RunRetainsPhase22Lifecycle|ClaimRejectsDirtyPathAdoption)$' ./cmd/fid/integration` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid/integration 8.431s` (terminated `exit_code: 0`). |
| Wall-time and local-evidence check | Task reports, QA-39, CF-041 | PASS | Task 1 records `cmd/fid` 527.099s, `internal/workflow` 231.010s, and tagged integration 669.586s; Task 2 records 440.185s, 193.138s, and 650.206s; Task 3 records 438.942s, 237.598s, and 702.407s. All meet the 600s ordinary-package and 900s tagged-integration budgets. The reports cite locally observed commands and local logs only; none uses hosted CI as phase evidence. |

## Per-Task Findings

### Task 1 - Schema baseline

#### Must-fix

- None.

#### Improvement

- None.

### Task 2 - Typed multi-Obligation Envelope contract

#### Must-fix

- **MF-01** [internal/workflow/envelope_test.go:20; cmd/fid/envelope_test.go:94] The passing named tests cover only three rejection cases and one three-Obligation/one-Decision success. They do not prove the contract's required 0/1/16/17 Obligation counts; each accepted and forbidden type-field combination; all five request-level negatives; 0/16/17 protected Decisions and its duplicate-class, duplicate-ordinal, and unguarded-action rejections; unchanged journal/projection for each rejection; or the qualifying `valid`/`fresh`/`available` evaluation state. The report's VAL-006 through VAL-010 claims therefore exceed the observed tests. Owning task: task-2. Repair boundary: add only Task 2 table-driven domain and command-surface tests in the existing Envelope test files, plus correct the Task 2 report's test/validation entries; do not change contracts, schema, product behavior, or gate definitions. Required evidence: a named test result proving every listed accepted and rejected authoring/evaluation state, each command-level rejection with unchanged journal count and projection digest, followed by the byte-identical Task 2 exact command list with observed terminations and retained prior failed attempts. Re-review scope: fresh independent reviewer inspects the new cases and report, reruns their named tests and the Task 2 portable integration command, and rechecks VAL-006, VAL-007, VAL-008, and VAL-010.

#### Improvement

- None.

### Task 3 - Non-mutating review Claim mode and Reviewer Run

#### Must-fix

- None.

#### Improvement

- None.

## Prior Finding Resolution

Not applicable: this is the initial review.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-01 | task-2 | VAL-006, VAL-007, VAL-008, VAL-010 | `internal/workflow/envelope_test.go`, `cmd/fid/envelope_test.go`, and `task-2-done.md` only; add complete authoring/evaluation matrix evidence and accurately record it. | Named matrix test with qualifying and non-qualifying states; no-Event/projection-invariance assertions for every reject; full byte-identical Task 2 gate list with command tokens and observed termination. | Inspect changed test cases and report; rerun the new named matrix test, Task 2 portable integration command, and report lint; verify all original MF-01 clauses. |

## Blocker Candidates

- None. No expected-class closed-registry rejection had a cause outside the failing task. MF-01 is a missing Task 2 evidence matrix, not an external registry or authority failure.

## Existing Escalation References

- ESC-23-002 is resolved under DEC-238 and pertains to Task 1's preserved historical performance recovery. It does not cover MF-01.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-23/qa/review-1.md
Must-fix: 1  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
