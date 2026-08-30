# Phase 21 - Mid-Phase Review 6c (tasks 16-18)

Reviewer: qa-reviewer
Date: 2026-08-30
Tasks covered: phase-21/task-16, phase-21/task-17, phase-21/task-18
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-21/qa/review-6b.md

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-21-013, closed by DEC-172; provenance only
- Covered validation assertions: VAL-055, VAL-056, VAL-057, VAL-058, VAL-059, VAL-060, VAL-061, VAL-062, VAL-063, VAL-064, VAL-065, VAL-066, VAL-067, VAL-068

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-055 | task-16 | ON TRACK | `TestGateRunPublishesSanitizedArtifactByDigest` asserts sanitized content, SHA-256 identity, exact metadata, and no staging file. |
| VAL-056 | task-16 | ON TRACK | `cmd/fid/artifact_test.go:74` runs changed-during-copy, directory, missing, non-regular, outside-CWD, oversized, and symlink rows; each asserts its `artifact-rejected` control, no Receipt Artifact link, and unsatisfied evaluation. |
| VAL-057 | task-16 | ON TRACK | `cmd/fid/artifact_test.go:252` links Evidence before loss/restoration and compares Artifact, Receipt, and Evidence identity snapshots while availability alone changes. |
| VAL-058 | task-17 | ON TRACK | `cmd/fid/evidence_test.go:180` rejects caller proof, and `:203` rejects persisted Obligation, Registration, Run, Issue, and target mismatches while asserting no Event or `evidence_records` mutation. |
| VAL-059 | task-17 | ON TRACK | The validity pair at `cmd/fid/evidence_test.go:18` holds freshness and Artifact availability fixed and changes only the `evidence-invalid` reason. |
| VAL-060 | task-17 | ON TRACK | The freshness pair at `cmd/fid/evidence_test.go:25` remains fresh after only `head_oid` changes and becomes stale only after `tree_oid` changes. |
| VAL-061 | task-17 | ON TRACK | The Artifact pair at `cmd/fid/evidence_test.go:38` changes only availability and the `artifact-unavailable` reason. |
| VAL-062 | task-17 | ON TRACK | `cmd/fid/evidence_test.go:46` byte-compares the full ordered Evaluation on repeat and after `rebuild`. |
| VAL-063 | task-17 | ON TRACK | `TestPhase21AttentionBoundary` requires exactly the unavailable `attention.v1` discovery row, rejects the unavailable leaf without journal/projection mutation, and accepts exactly one non-test production match. |
| VAL-064 | task-18 | ON TRACK | `TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease` compares the evaluator result and observes exactly one completion Event. |
| VAL-065 | task-18 | ON TRACK | `cmd/fid/task18_completion_resume_test.go:16` and `:54` cover unsuccessful, missing, invalid, stale, and Artifact-missing evidence, each with its exact sole reason and zero Issue/Event mutation. |
| VAL-066 | task-18 | ON TRACK | The positive completion test reuses the qualified terminal Evaluation after finishing the Run and releasing the Claim. |
| VAL-067 | task-18 | ON TRACK | `TestEnvelopeCompleteReturnsAllCurrentReasonsWithoutMutation` compares Envelope, Issue, Evidence, `task_completions`, and Event state. |
| VAL-068 | task-18 | ON TRACK | `cmd/fid/task18_completion_resume_test.go:201` asserts active and completed IDs, obligations, completion, Claim/Run selection, ordered next actions, no recovery/respawn, and zero Event deltas. |

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Task 16 exact command audit | `task-16-done.md` against `task-16-artifacts.md` | PASS | Every required command is recorded PASS with its own token and observed termination. Fresh exact scoped `go test -race ./internal/evidence ./internal/config ./internal/journal ./internal/cli ./cmd/fid/...` ended `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 17 exact command audit | `task-17-done.md` against `task-17-evidence-evaluation.md` | PASS | All six DEC-172 commands are recorded PASS; the superseded broad search remains a distinct historical BLOCKED attempt. Fresh exact scoped `go test -race ./internal/evidence ./internal/validation ./internal/workflow ./internal/cli ./cmd/fid/...` ended `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 18 exact command audit | `task-18-done.md` against `task-18-completion-resume.md` | PASS | All six required commands, including tagged integration, are recorded PASS with observed termination. Fresh exact scoped `go test -race ./internal/workflow ./internal/issues ./internal/evidence ./internal/validation ./internal/cli ./cmd/fid/...` ended `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 16 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-16-done.md project-docs/docs/phases/phase-21/tasks/task-16-artifacts.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 17 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-17-done.md project-docs/docs/phases/phase-21/tasks/task-17-evidence-evaluation.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 18 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-18-done.md project-docs/docs/phases/phase-21/tasks/task-18-completion-resume.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Reproduced Window 6 behavior | `go test -race ./cmd/fid -run '^(TestGateRunUnsafeArtifactsEmitControlsAndRemainRequiredUnavailable|TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation|TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild|TestEvidenceLinkRejectsEveryCallerAuthoredProof|TestEvidenceLinkRejectsPersistedReceiptBindingMismatches|TestPhase21AttentionBoundary|TestTerminalSuccessTransitionRejectsUnsuccessfulReceiptEvidence|TestTerminalSuccessTransitionRejectsEveryOtherNonqualifyingEvidenceState|TestResumeProjectsExactActiveAndCompletedDurableStateWithoutRecovery)$' -count=1` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 89.516s` (terminated `exit_code: 0`). This independently exercises every original Window 6 finding. |
| DEC-172 static boundary | `rg -n -i --glob '*.go' --glob '!**/*_test.go' --glob '*.sql' 'attention|attention_items' cmd/fid internal migrations` | PASS | Sole production match: `internal/workflow/capabilities.go:105: {"attention.v1", false, `{"state":false}`, "phase-22-unavailable"},` (terminated `exit_code: 0`). No Phase 22 Attention surface is present. |
| Commit and report provenance | `git rev-parse --verify` for `457366d`, `200b807`, `1a8ecaa`, `83706d`, and `74e1efc`; current diffs | PASS | All reported commits resolve. `83706d` is `test(evidence): prove independent evaluation dimensions` with `Refs: phase-21/task-17`; `74e1efc` corrects the task-17 report. `git diff --check 22ad81d..HEAD` terminated `exit_code: 0`. |

## Task 16 - Safe Referenced Artifacts

### Must-fix

- None.

### Improvement

- None.

## Task 17 - Trusted Evidence Evaluation

### Must-fix

- None.

### Improvement

- None.

## Task 18 - Guarded Terminal Transition, Completion, And Resume

### Must-fix

- None.

### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-16-01 | task-16 | RESOLVED | Current unsafe-candidate matrix contains all seven named rejection/control/evaluation rows, and the Evidence-linked loss/restore lifecycle preserves Artifact, Receipt, and Evidence identity. The reproduced focused suite passed with `exit_code: 0`. |
| MF-17-01 | task-17 | RESOLVED | Current binding matrix covers Obligation, Registration, Run, Issue, repository target, and caller proof; every rejected row asserts zero Evidence Event and `evidence_records` deltas. The reproduced focused suite passed with `exit_code: 0`. |
| MF-17-02 | task-17 | RESOLVED | Commit `83706d` adds independent paired validity, freshness, and Artifact fixtures; the `head_oid` control remains fresh, `tree_oid` becomes stale, and the full ordered-reason Evaluation is byte-identical on repeat and after rebuild. The reproduced focused suite passed with `exit_code: 0`. |
| MF-18-01 | task-18 | RESOLVED | Current terminal-transition tests cover unsuccessful, missing, invalid, stale, and Artifact-missing evidence, each with its exact reason and zero Issue/Event mutation. The reproduced focused suite passed with `exit_code: 0`. |
| MF-18-02 | task-18 | RESOLVED | Current active and completed resume fixtures assert durable IDs, obligation/completion state, Claim/Run selection, ordered actions, no recovery/respawn, and zero Event deltas. The reproduced focused suite passed with `exit_code: 0`. |

## Cross-cutting observations

- DEC-172 remains exact: discovery exposes the sole unavailable `attention.v1` row as negative capability metadata; the boundary test and current production scan reject any additional Attention table, Event, Effect, leaf, state, or operation surface.

## Blocker Candidates

- None.

## Existing Escalation References

- ESC-21-013, closed by DEC-172; provenance only, not a blocker candidate.

## Final Lines

Review written: project-docs/docs/phases/phase-21/qa/review-6c.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
