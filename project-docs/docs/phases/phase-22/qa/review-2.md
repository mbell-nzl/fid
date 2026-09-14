# Phase 22 - Mid-Phase Review 2 (Tasks 4-6)

Reviewer: qa-reviewer
Date: 2026-09-06
Tasks covered: phase-22/task-4, phase-22/task-5, phase-22/task-6
Review type: INITIAL
Prior review: none

## Verdict

REPAIR

## Scope Reviewed

- Task specs: `project-docs/docs/phases/phase-22/tasks/task-4.md`,
  `task-5.md`, and `task-6.md`.
- Completion reports: `project-docs/docs/phases/phase-22/tasks/task-4-done.md`,
  `task-5-done.md`, and `task-6-done.md`.
- Validation contract and QA plan: `qa/validation-contract.md` and
  `qa/qa-plan.md`, including QA-14..25, QA-84, and the window-2 entry
  condition.
- DEC-192, DEC-193, and DEC-194; `triage-log.md`; `execution-scope.md`; task
  commits `59c62b9e1d82c5f698d1e173908446d35b64f539`,
  `b7dce8e747a3f2d720002faee3db7f95b3cb4019`, and
  `69941fc1b6262fdd7efaaeecfe72760390dbd511`.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 2
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-22-003 (resolved by DEC-194)
- Covered validation assertions: VAL-013, VAL-014, VAL-015, VAL-016, VAL-017, VAL-018, VAL-019, VAL-020, VAL-021, VAL-022, VAL-023, VAL-024

Task 6 supplies the required real built-binary positive journey: its tagged
test uses the binary created by integration `TestMain`, and the fresh exact
journey rerun passed. The DEC-194 transfer is also correctly recorded: each
transferred file names phase-22/task-4 as accepted owner, says it was absent
from every earlier commit, cites DEC-194, and names the unchanged whole-suite
gate that executes it. `git log -S` confirms the four paths first appear in
`69941fc`; task-4's report has no diff after the task-6 recovery checkpoint
`c69456f`, so task 4 was neither reopened nor edited by the transfer.

This satisfies the required VAL-013/VAL-022 provenance audit. The clean,
stable inspect and checkpoint-base behavior is exercised by
`TestRunInspectIsStableAndCheckpointBecomesComparisonBase`; repository-byte
inert inspection is exercised by `TestRunInspectLeavesRepositoryByteInert`.
Those task-6 tests execute the task-4-owned checkpoint/inspection behavior
under DEC-194. The absence of that producer/read-API evidence from the
task-4 report is therefore not a task-6 finding.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-013 | task-4; DEC-194 execution in task-6 | ON TRACK | Fresh selected project/workflow reproduction passed; inspect is stable, checkpoint-aware, and byte-inert under task-4 provenance. |
| VAL-014 | task-4 | ON TRACK | `TestObserveRepositoryPinsCanonicalDirtyVectorsAndExclusions` and exact-adoption coverage passed in fresh selected reproduction. |
| VAL-015 | task-4 | ON TRACK | Task-4 reconciliation tests and report evidence cover unsafe/conflicting state rejection; required gate receipt is present. |
| VAL-016 | task-4 | ON TRACK | Observer and transferred inspect code are read-only; `TestRunInspectLeavesRepositoryByteInert` passed in fresh selected reproduction. |
| VAL-017 | task-5 | ON TRACK | Fresh recovery reproduction passed; the test asserts distinct linked successor and interrupted predecessor. |
| VAL-018 | task-5 | ON TRACK | Fresh recovery reproduction covers same-key replay/conflicting reuse; report records projection parity. |
| VAL-019 | task-5 | ON TRACK | Fresh ordinary-and-force unknown-liveness matrix passed. |
| VAL-020 | task-5 | ON TRACK | Fresh force-authority/fencing test passed. |
| VAL-021 | task-5 | ON TRACK | Fresh recovery test compares the Issue and immutable Envelope before and after recovery. |
| VAL-022 | task-6 | ON TRACK | The real tagged subprocess journey passed independently; it initializes, activates, starts, observes, inspects, checkpoints, expires, recovers, and inspects the successor. |
| VAL-023 | task-6 | NOT ON TRACK | No required per-leaf accepted/unknown-field conformance corpus exists; see MF-02. |
| VAL-024 | task-6 | NOT ON TRACK | The claimed test invokes an in-process error renderer and proves neither a process exit nor the required real-process failure corpus; see MF-02. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task-4 report control | `project-docs/tools/report-lint.sh task-4-done.md task-4.md` | PASS | The three-report lint transcript ended `PASS report-lint completed with no failures` (aggregate terminated `exit_code: 0`); task 4 records every required gate including the DEC-193 repeat command. |
| Task-5 report control | `project-docs/tools/report-lint.sh task-5-done.md task-5.md` | PASS | The same transcript reports all five exact gates and `OUTCOME COMPLETE`. |
| Task-6 report control | `project-docs/tools/report-lint.sh task-6-done.md task-6.md` | PASS | The report preserves the prior red/`RECOVERY_REQUIRED` history, the report-repair commits, the five final gates, and the byte-identical preserved-red rerun. |
| Commit provenance | `git rev-parse --verify <reported-hash>^{commit}` | PASS | All three reported implementation hashes resolved; their subjects and `Refs: phase-22/task-N` footers match their reports. |
| Independent task-6 tagged gate | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration  9.145s` (terminated `exit_code: 0`). This is the unchanged required task-6 integration gate and builds/runs the binary. |
| Independent DEC-194 rerun | `go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration  1.162s` (terminated `exit_code: 0`). This is byte-identical to the preserved red command. |
| Representative task-4 evidence | `go test ./internal/project -run '^(TestRecoveryObservationDetectsChangedUntrackedFileWithNULDelimitedPaths|TestObserveRepositoryPinsCanonicalDirtyVectorsAndExclusions)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/internal/project  0.297s` (terminated `exit_code: 0`). |
| Representative task-4/5/6 workflow evidence | selected recovery, checkpoint, and inspect tests | PASS | `ok  github.com/mbell-nzl/fid/internal/workflow  1.664s` (terminated `exit_code: 0`); covers recovery links/liveness/force authority and transferred checkpoint/inspect behavior. |
| Window-2 dual-adapter entry condition | `test ! -e cmd/fid/integration/p22_fake_test.go`; `test ! -d cmd/fid/integration/testdata/p22_conformance` | FAIL | Both absence checks returned `exit 0 (no output)` (each terminated `exit_code: 0`): neither required path exists. |
| VAL-024 evidence placement | `cmd/fid/run_recovery_test.go:10-40` | FAIL | The sole cited test calls `renderRunRecoveryError` directly; it does not execute a built binary or assert process exit status. |

## Per-Task Findings

### Task 4 - Read-only repository reconciliation

#### Must-fix

- None.

#### Improvement

- None.

Task-4's exact gate receipts, retained failed attempts/deviations, report-lint
result, and implementation commit are present. DEC-194 leaves its accepted
ownership intact and requires no task-4 repair here.

### Task 5 - Atomic linked Run recovery

#### Must-fix

- None.

#### Improvement

- None.

Task-5's exact gate receipts and report-lint result are complete. Its current
tests assert the linked successor, terminal predecessor, replay, liveness,
force-authority, and immutable-contract behaviors claimed for VAL-017..021.

### Task 6 - Recovery CLI and early process journey

#### Must-fix

- [cmd/fid/integration/p22_fake_test.go: absent; cmd/fid/integration/testdata/p22_conformance/: absent] The reusable dual-adapter case structure required by task 6 and the QA-plan window-2 entry condition does not exist. The passing `TestP22EarlyRecoveryJourney` is a valid real subprocess test, but it is one direct binary journey and supplies no test-only fake or shared cases run against both adapters. Owning task: task-6. Repair boundary: add only the task-6-scoped reusable semantic case runner, `p22_fake_test.go`, and `testdata/p22_conformance/`, then execute the implemented recovery cases against fake and fid without changing accepted contracts. Required evidence: committed files plus test output showing each shared case runs against both adapters and the unchanged required tagged integration gate. Re-review scope: inspect the runner/fake/fixtures and rerun the exact tagged task-6 integration gate and affected cases. (no VAL; QA-plan window-2 entry condition)

- [cmd/fid/run_recovery_test.go:10] VAL-023 and VAL-024 are claimed without the required process-level evidence. `TestRunRecoveryErrorsUseClosedRecoveryEnvelopes` directly calls `renderRunRecoveryError`; it does not invoke any recovery leaf, exercise accepted/unknown-field request/result fixtures for every leaf, or observe a process exit. The only tagged binary test is a positive journey. Owning task: task-6. Repair boundary: in the new task-6 conformance corpus, add real-binary cases for each section-3 recovery leaf's accepted and unknown/invalid input/result behavior, and failure cases for not-found, authority, liveness, revision, and conflict with exact JSON code and nonzero exit status; retain the existing production behavior unless a case exposes an in-scope defect. Required evidence: the corpus transcript identifies every case, adapter, JSON code, and observed binary exit; the unchanged five-command task-6 gate passes; the preserved-red journey command is rerun byte-identically. Re-review scope: inspect the corpus and error assertions, verify no in-process substitute is presented as process evidence, and rerun the tagged gate and affected binary tests. (VAL-023, VAL-024)

#### Improvement

- None.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-01 | task-6 | no VAL | Add the scoped p22 semantic runner, fake adapter, and fixtures only; run common implemented recovery cases against fake and fid. | Committed runner/fake/corpus; case output for both adapters; exact tagged integration gate. | Inspect new artifacts and rerun tagged gate plus shared cases. |
| MF-02 | task-6 | VAL-023, VAL-024 | Add per-leaf closed-schema and real-binary failure corpus cases within task-6 CLI/integration scope; correct only a defect exposed by those cases. | Case transcript with exact JSON code and process exit; unchanged five-command gate; byte-identical preserved-red rerun. | Inspect all recovery-leaf fixtures and rerun affected binary cases and tagged gate. |

## Blocker Candidates

- None. Both findings have decision-free, bounded ownership in task 6's accepted
  test/CLI/integration scope; no contract, gate, or task-4 ownership change is
  required.

## Existing Escalation References

- ESC-22-003 is resolved by DEC-194. It supplies provenance for the transferred
  task-4 checkpoint/inspect implementation; it does not cover these task-6
  conformance repairs.

## Final Lines

Review written: project-docs/docs/phases/phase-22/qa/review-2.md
Must-fix: 2  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
