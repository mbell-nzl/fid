# Phase 21 - Mid-Phase Review 6 (Tasks 16-18)

Reviewer: qa-reviewer
Date: 2026-08-30
Tasks covered: phase-21/task-16, phase-21/task-17, phase-21/task-18
Review type: INITIAL
Prior review: none; `project-docs/docs/phases/phase-21/qa/review-5b.md` covered tasks 13-15.

## Verdict

REPAIR

- `PASS`: no Must-fix finding remains; Improvements may remain.
- `REPAIR`: valid Must-fix work has bounded, decision-free ownership in an existing reviewed task. The window does not advance.
- `BLOCK`: protected authority is required, the repair exceeds the reviewed task/window scope, or the bounded no-progress recovery budget is exhausted.

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-21/tasks/task-16-artifacts.md`
  - `project-docs/docs/phases/phase-21/tasks/task-17-evidence-evaluation.md`
  - `project-docs/docs/phases/phase-21/tasks/task-18-completion-resume.md`
- Completion reports:
  - `project-docs/docs/phases/phase-21/tasks/task-16-done.md`
  - `project-docs/docs/phases/phase-21/tasks/task-17-done.md`
  - `project-docs/docs/phases/phase-21/tasks/task-18-done.md`
- Validation contract: `project-docs/docs/phases/phase-21/qa/validation-contract.md` (VAL-055 through VAL-068).
- QA plan: `project-docs/docs/phases/phase-21/qa/qa-plan.md` (Window 6, QA-53 through QA-66).
- Contracts and authority: `contracts.md` sections 2, 9, 10, 15, and 16; `public-interface.md` sections 4 and 5; DEC-172; the Window 6 run-ledger entries.
- Implementation and report commits: `89edce34bca87308117a1fe4bafad56c6e9bc292`, `cc890a8ed796d19e48daee78d58656e2b4c7eff7`, `7277139ee5d5dde4cd3627335d39b14f616f46a1`, `2519799d58ff27fd78f17035d21e81cb42588077`, `2313235ba463060538f2f0d4bcfa61f039fcdc13`, and `a91331362b4b0b1268a4782d03312bd963f6bd03`.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 5
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-21-013, closed by DEC-172; provenance context only
- Covered validation assertions: VAL-055, VAL-056, VAL-057, VAL-058, VAL-059, VAL-060, VAL-061, VAL-062, VAL-063, VAL-064, VAL-065, VAL-066, VAL-067, VAL-068

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-055 | task-16 | ON TRACK | `TestGateRunPublishesSanitizedArtifactByDigest` passed fresh and observes sanitized bytes, digest, metadata, and no staging file. |
| VAL-056 | task-16 | NOT ON TRACK | The current unit-only candidate tests assert reasons, but do not produce each required `artifact-rejected` control or prove the corresponding `required-available` evaluation remains unsatisfied (MF-16-01). |
| VAL-057 | task-16 | NOT ON TRACK | The lifecycle test preserves Artifact and Receipt IDs but creates no Evidence row and therefore cannot prove Evidence identity preservation (MF-16-01). |
| VAL-058 | task-17 | NOT ON TRACK | A qualifying link and caller-authored fields are tested, but no mismatch row exists for Obligation, Registration, Run, Issue, or target (MF-17-01). |
| VAL-059 | task-17 | NOT ON TRACK | The only non-qualifying evaluation changes validity, tree, and Artifact together, not validity alone (MF-17-02). |
| VAL-060 | task-17 | NOT ON TRACK | The stale-tree observation is coupled to invalidity and missing Artifact state; it does not isolate the tree dimension or prove branch labels are immaterial (MF-17-02). |
| VAL-061 | task-17 | NOT ON TRACK | The missing-Artifact observation is coupled to invalidity and tree staleness, rather than showing the isolated Artifact reason/status diff (MF-17-02). |
| VAL-062 | task-17 | NOT ON TRACK | One three-reason assertion exists, but the QA-required repeat and post-rebuild byte-identical golden evidence is absent (MF-17-02). |
| VAL-063 | task-17 | ON TRACK | Fresh `TestPhase21AttentionBoundary` PASS requires exactly one canonical unavailable row, verifies the unavailable leaf has no journal/projection mutation, and inventories the sole production match. |
| VAL-064 | task-18 | ON TRACK | The successful Receipt link, terminal transition, same evaluation object, and exactly one completion Event are asserted by the fresh focused test. |
| VAL-065 | task-18 | NOT ON TRACK | Only the unsuccessful-Receipt case is exercised; missing, invalid, stale, and Artifact-missing evidence lack the required atomic negative matrix (MF-18-01). |
| VAL-066 | task-18 | ON TRACK | `TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease` compares `evidence.evaluate` and completion evaluations and observes one completion Event after terminal success. |
| VAL-067 | task-18 | ON TRACK | The failed completion test asserts both current reasons and unchanged Envelope, Issue, Evidence, legacy `task_completions`, and Event state. |
| VAL-068 | task-18 | NOT ON TRACK | The resume test proves returned IDs and zero Event delta only; it does not assert the exact bounded projection and closed `next_actions`, nor the unavailable recovery/respawn counterparts (MF-18-02). |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 16 report command audit | `task-16-done.md` against `task-16-artifacts.md` | PASS | All five exact required commands are recorded PASS with verbatim tokens and observed termination. Fresh `report-lint.sh` ended `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 17 report command audit | `task-17-done.md` against `task-17-evidence-evaluation.md` | PASS | All six revised DEC-172 commands are recorded PASS with captured termination. The pre-DEC-172 broad `rg` remains a BLOCKED historical attempt with observed `exit_code: 0`, distinct from the accepted boundary-test command. Fresh report lint ended `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 18 report command audit | `task-18-done.md` against `task-18-completion-resume.md` | PASS | All six exact required commands, including the tagged integration gate, are recorded PASS with captured termination. Fresh report lint ended `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| QA-53 through QA-55 reproduced evidence | `go test -race ./cmd/fid ./internal/evidence -run '^(TestGateRunPublishesSanitizedArtifactByDigest|TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation|TestArtifactCandidatePathRejectsEveryUnsafeCandidate|TestPublishArtifactCandidateRejectsChangedDuringCopy|TestPublishArtifactCandidateRejectsOversizedContent)$' -count=1` | PASS | `ok github.com/mbell-nzl/fid/internal/evidence 1.031s` (terminated `exit_code: 0`). The test execution exposes the QA-54/55 matrix omissions in MF-16-01. |
| QA-56 through QA-61 reproduced evidence | `go test -race ./cmd/fid -run '^(TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestEvidenceLinkRejectsEveryCallerAuthoredProof|TestPhase21AttentionBoundary)$' -count=1` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 7.684s` (terminated `exit_code: 0`). The source assertions expose the missing QA-56 through QA-60 matrix coverage in MF-17-01 and MF-17-02. |
| QA-62 through QA-66 reproduced evidence | `go test -race ./cmd/fid -run '^(TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease|TestEnvelopeCompleteReturnsAllCurrentReasonsWithoutMutation|TestResumeReadsDurableFactsWithoutMutation|TestTerminalSuccessTransitionRejectsUnsuccessfulReceiptEvidence)$' -count=1` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 15.678s` (terminated `exit_code: 0`). The exercised test set contains only one QA-63 negative state and no QA-66 exact-action assertion. |
| DEC-172 static boundary check | `rg -n -i --glob '*.go' --glob '!**/*_test.go' --glob '*.sql' '\\battention\\b|attention_items' cmd/fid internal migrations` | PASS | Exactly one production match, `internal/workflow/capabilities.go:105: {"attention.v1", false, `{\"state\":false}`, "phase-22-unavailable"},` (terminated `exit_code: 0`; one canonical match is the DEC-172 pass condition). |
| Commit and diff evidence | `git rev-parse --verify` for each implementation hash; `git diff --check 89edce3^ 2313235` | PASS | All three reported implementation hashes resolved (each terminated `exit_code: 0`); their subjects and `Refs: phase-21/task-N` footers match. `git diff --check` produced `exit 0 (no output)` (terminated `exit_code: 0`). |

## Per-Task Findings

### Task 16 - Safe Referenced Artifacts

#### Must-fix

- [internal/evidence/artifact_test.go:10; cmd/fid/artifact_test.go:71] QA-54 and QA-55 evidence is narrower than VAL-056/057. Unsafe-path tests only inspect internal rejection reasons, not the required per-case `artifact-rejected` control and unsatisfied `required-available` evaluation; the availability lifecycle snapshots omit Evidence identity entirely. Owning task: task-16. Repair boundary: add observable CLI/integration coverage for every named unsafe candidate and link Evidence before the loss/restore lifecycle; repair product code only if the new tests reveal a behavioral defect. Required evidence: each named rejection/control/evaluation row, a lifecycle snapshot that includes Artifact/Receipt/Evidence IDs, the exact task-16 command block, updated report, and fresh two-argument report lint. Re-review scope: inspect the new matrix and rerun its focused test command plus report lint.

#### Improvement

- None.

### Task 17 - Trusted Evidence Evaluation

#### Must-fix

- [cmd/fid/evidence_test.go:83] QA-56/VAL-058 requires one rejected row for each mismatched Obligation, Registration, Run, Issue, and target, in addition to caller-proof rejection. The test exercises only caller-supplied fields against an otherwise exact persisted chain; it does not assert any required binding mismatch. Owning task: task-17. Repair boundary: add an observable link matrix that mutates or selects each named persisted binding dimension and proves typed rejection/no Evidence Event; repair implementation only if a matrix row succeeds. Required evidence: one qualifying row, all five mismatch rows, caller-proof row, exact task-17 command block, report update, and fresh two-argument report lint. Re-review scope: inspect the persisted-chain matrix and rerun the focused evidence/Attention tests and report lint.

- [cmd/fid/evidence_test.go:67] QA-57 through QA-60 and VAL-059 through VAL-062 are not proven independently. The test changes validity, Artifact availability, and `tree_oid` before one evaluation at line 80, so it cannot show an isolated reason/dimension change; it also has no repeated or post-rebuild byte-identical multi-failure golden. Owning task: task-17. Repair boundary: add paired valid/invalid, fresh/stale, and available/missing fixtures that alter one dimension at a time, prove tree rather than branch-label freshness, and add repeat/rebuild byte-comparison of the full ordered reason object. Required evidence: qualifying and non-qualifying output for all three dimensions, repeat/rebuild golden output, exact task-17 command block, report update, and fresh two-argument report lint. Re-review scope: inspect the independent fixtures and rerun the focused evidence/Attention tests and report lint.

#### Improvement

- None. DEC-172 handling is correct: the required unavailable row is retained as the sole production match, `TestPhase21AttentionBoundary` passed fresh, the original failed broad search remains historical evidence, and no Phase 22 Attention surface was found.

### Task 18 - Guarded Terminal Transition, Completion, And Resume

#### Must-fix

- [cmd/fid/task18_completion_resume_test.go:11] QA-63/VAL-065 requires an atomic negative matrix for missing, unsuccessful, invalid, stale, and Artifact-missing Evidence. The current test has only the unsuccessful Receipt case. Owning task: task-18. Repair boundary: add the four missing cases, each asserting the exact reasons and zero Event/Issue mutation; repair implementation only if a case advances the terminal transition. Required evidence: all five negative rows, exact task-18 command block including tagged integration, report update, and fresh two-argument report lint. Re-review scope: inspect the matrix and rerun the task-18 focused test command, tagged integration command, and report lint.

- [cmd/fid/task18_completion_resume_test.go:121] QA-66/VAL-068 requires the exact bounded resume projection and closed next actions, with no recovery or respawn. The test asserts only Envelope/Issue IDs and a zero Event delta; it never checks `obligations`, `completion`, `next_actions`, Claim/Run selection, or unavailable recovery/respawn behavior. Owning task: task-18. Repair boundary: add a fresh-command projection assertion covering the exact durable fields and ordered next actions for a representative active state, plus a terminal/completed state proving no recovery or respawn action; repair implementation only if those reads differ. Required evidence: before/after Event counts and byte/field-level projection assertions for both states, exact task-18 command block including tagged integration, report update, and fresh two-argument report lint. Re-review scope: inspect the read-only cases and rerun the task-18 focused test command, tagged integration command, and report lint.

#### Improvement

- None.

## Prior Finding Resolution

Not applicable. This is the initial review for Window 6; `review-5b.md` is a resolved re-review of the preceding tasks 13-15.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-16-01 | task-16 | VAL-056, VAL-057 | Add only the missing observable unsafe-candidate/control/evaluation matrix and Evidence-identity lifecycle coverage; repair product code only if the test exposes a defect. | Named adverse rows, lifecycle IDs, exact task-16 command block, report, report lint. | New task-16 tests, focused Artifact command, report lint. |
| MF-17-01 | task-17 | VAL-058 | Add only the exact Receipt-binding mismatch matrix and no-Evidence-Event assertions. | One qualifying plus Obligation/Registration/Run/Issue/target/caller-proof rejection rows, exact task-17 command block, report, report lint. | New evidence matrix, focused Evidence/Attention command, report lint. |
| MF-17-02 | task-17 | VAL-059, VAL-060, VAL-061, VAL-062 | Add only independent paired dimension fixtures and repeat/rebuild ordered-reason golden coverage. | Paired state outputs, branch-label control, byte-identical repeat/rebuild result, exact task-17 command block, report, report lint. | New dimension/golden fixtures, focused Evidence/Attention command, report lint. |
| MF-18-01 | task-18 | VAL-065 | Add only the four omitted terminal-transition negative Evidence cases with existing atomicity assertions. | Five-case matrix, exact task-18 command block including tagged integration, report, report lint. | New task-18 negative matrix, focused task-18 command, tagged integration, report lint. |
| MF-18-02 | task-18 | VAL-068 | Add only exact active and terminal/completed resume projection/action assertions with no recovery/respawn. | Durable projection/action fields and zero-Event proof for both states, exact task-18 command block including tagged integration, report, report lint. | New resume cases, focused task-18 command, tagged integration, report lint. |

## Blocker Candidates

- None. Every finding is decision-free test/evidence completion inside an existing reviewed task. DEC-172 already supplies the protected Attention authority; the current implementation conforms to it.

## Existing Escalation References

- ESC-21-013, closed by DEC-172. It is retained only as provenance for the corrected Attention evidence method and is not an active blocker candidate.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-21/qa/review-6.md
Must-fix: 5  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
