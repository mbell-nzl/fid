# Phase 21 - Mid-Phase Review 2b (tasks 4-6)

Reviewer: qa-reviewer
Date: 2026-08-23
Tasks covered: phase-21/task-4, phase-21/task-5, phase-21/task-6
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-21/qa/review-2.md
Re-review of: review-2

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-011, VAL-012, VAL-013, VAL-014, VAL-015, VAL-016, VAL-017, VAL-018, VAL-019, VAL-020, VAL-021, VAL-022

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Fresh or retained evidence |
|---|---|---|---|
| VAL-011 | task-4 | ON TRACK | `cmd/fid/activation_test.go:200` rejects every listed skipped or extra lifecycle operation from prepared, staged, committed, finalized, and aborted states; every row asserts unchanged journal Event count. |
| VAL-012 | task-4 | ON TRACK | `cmd/fid/activation_test.go:233` proves successful abort from prepared and staged, exactly one abort Event, and unavailable `ResolveTrustedActor`; lines 209 and 213 retain committed/finalized abort rejection without an Event. |
| VAL-013 | task-4 | ON TRACK | The committed report retains the declared recovery-source absence PASS with expected `exit_code: 1`; `cmd/fid/activation_test.go:209` and `:213` observe post-commit `activation-recovery-unavailable`. |
| VAL-014 | task-4 | ON TRACK | `cmd/fid/activation_test.go:397` rejects staged and committed profile resolution and resolves exact actor type, producer, provenance, and capabilities only after finalization. |
| VAL-015 | task-4 | ON TRACK | `cmd/fid/activation_test.go:105` exercises both `caller_capabilities` and `caller_provenance_ref`, observes typed `invalid-argument`, and asserts unchanged journal Event count for each. |
| VAL-016 | task-5 | ON TRACK | Prior review-2 evidence at `cmd/fid/gate_test.go:26` covers rejection before finalization and successful finalized Activation-owned Definition registration. |
| VAL-017 | task-5 | ON TRACK | Prior review-2 evidence at `cmd/fid/gate_test.go:71` covers each listed unsupported policy field, exact `executor-policy-unsupported`, and zero Registration persistence. |
| VAL-018 | task-5 | ON TRACK | Prior review-2 evidence at `cmd/fid/gate_test.go:112` covers same-binding reuse, changed-binding version 2, and immutable version-1 resolution. |
| VAL-019 | task-6 | ON TRACK | Prior review-2 evidence at `cmd/fid/envelope_test.go:13` covers the exact one-revision boundary and wider-input rejection. |
| VAL-020 | task-6 | ON TRACK | Prior review-2 evidence at `cmd/fid/envelope_test.go:89` covers the exact gate-evidence dual guards and closed-value rejection rows. |
| VAL-021 | task-6 | ON TRACK | Prior review-2 evidence at `cmd/fid/envelope_test.go:97` covers null immutable budget values across create, get, and activate. |
| VAL-022 | task-6 | ON TRACK | Prior review-2 evidence at `cmd/fid/envelope_test.go:42` covers deferred-value typed rejection with no Event. |

## Gate And Evidence Checks

The two commands below ran fresh from `/home/mark/projects/fid`. The fresh re-review requires no further QA-plan command beyond these controls; task-5 and task-6 evidence remains the committed-report and review-2 evidence required by the linked re-review scope.

| Check | Command or evidence | Result | Verbatim token and observed termination |
|---|---|---|---|
| Task 4 representative gate | `go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...` | PASS | `github.com/mbell-nzl/fid/cmd/fid` (terminated `exit_code: 0`) |
| Task 4 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-4-done.md project-docs/docs/phases/phase-21/tasks/task-4-activation-trust.md` | PASS | `OUTCOME COMPLETE` (terminated `exit_code: 0`) |
| Task 4 repair scope and contract preservation | `git diff --name-status 8ca2f8a... 5606254...` | PASS | only `cmd/fid/activation_test.go` and `project-docs/docs/phases/phase-21/tasks/task-4-done.md` changed; no task spec, QA plan, validation contract, contracts, public interface, decision, or ledger path changed. |
| Task reports and implementation commits | `git rev-parse --verify` on `ad9098f...`, `0c4c04e...`, `56dad8f...`, `cd83fe9...`, and `5606254...` | PASS | every named hash resolved as a commit (terminated `exit_code: 0`). |
| Task 5 and task 6 retained evidence | committed reports plus review-2 | PASS | both reports retain all required task-spec gate rows, deviations, canonical `Category: none`, and previous two-argument report-lint PASS evidence. |
| Ledger prerequisite | `run-ledger.md:645-737` | PASS | records review-2 TASK-RECOVERY, REPAIR attempts 1 and 2, and final task-4 TASK COMPLETE before this linked review. |
| Task 7 authorization | `rg` search of the current ledger | PASS | no task-7 recovery, completion, or PASS REVIEW record (`exit 1 (no output)`); task 7 remains unauthorized pending the Phase Runner's post-review PASS event. |

## Task 4 - Activation Trust Lifecycle

### Must-fix

- None.

### Improvement

- None.

### Evidence

- Review-2 MF-1 is resolved by the committed `cd83fe9f12aee24bfc1fb3e73145bdaffaf731ed` test change. The lifecycle table at [activation_test.go](/home/mark/projects/fid/cmd/fid/activation_test.go:200) covers every rejected skipped/extra transition in the repaired scope, including post-commit abort rejection; its Event count is unchanged on each rejection.
- The pre-commit table at [activation_test.go](/home/mark/projects/fid/cmd/fid/activation_test.go:233) qualifies both allowed abort states and proves no active trusted actor remains. The caller-augmentation table at [activation_test.go](/home/mark/projects/fid/cmd/fid/activation_test.go:105) covers capability and provenance variants with typed rejection and no Event.
- The scope diff contains only the permitted test and report files. The accepted lifecycle and point-of-no-return contract in `contracts.md` section 4, the public error contract in `public-interface.md` section 7, and the task-4 spec are unchanged. `56062540e901553ad5ba1ca7ce208774c05f979f` changes the report category to canonical `none` and preserves the earlier failed, red, green, and passing evidence, command tokens, and observed terminations. Fresh report lint returned `OUTCOME COMPLETE`.

## Task 5 - Activated Gate Registration

### Must-fix

- None.

### Improvement

- None.

### Evidence

- Task 5 remains intact at implementation commit `0c4c04ea0c7d0786e2e6cf88a7573d8a3105aef5` and committed report repair `b794ffe066b0ef16d455948175dd2d117cb4cb3a`. Review-2's qualifying/non-qualifying evidence remains sufficient for VAL-016 through VAL-018, and its report shows every required task-spec command with token and observed termination.

## Task 6 - One-Revision Envelope Values

### Must-fix

- None.

### Improvement

- None.

### Evidence

- Task 6 remains intact at implementation commit `56dad8fa0dcb9b79f7df2a31b356f1b2eb8a694e` and committed report `10fd45f8a58d3c304c56eb0b1349cae5aaeb73b7`. Review-2's qualifying/non-qualifying evidence remains sufficient for VAL-019 through VAL-022, and its report shows every required task-spec command with token and observed termination.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-1: incomplete negative lifecycle, pre-commit abort/no-trust, and caller-provenance/no-Event evidence for VAL-011, VAL-012, and VAL-015 | phase-21/task-4 | RESOLVED | `cmd/fid/activation_test.go:105`, `:200`, and `:233`; fresh representative gate PASS; fresh two-argument report-lint `OUTCOME COMPLETE`; bounded repair diff and canonical report category verified. |

## Cross-cutting observations

- The run ledger is runner-owned and remains its only pre-existing worktree modification. It contains the linked recovery sequence and no Phase Runner PASS REVIEW event yet; this review does not authorize task 7.
- Task-5 and task-6 were not modified by the task-4 repair. Their committed reports and review-2 evidence remain valid for this linked re-review.

## Blocker Candidates

- None.

## Existing Escalation References

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-21/qa/review-2b.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
```
