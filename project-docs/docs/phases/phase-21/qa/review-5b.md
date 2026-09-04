# Phase 21 - Mid-Phase Review 5b (Tasks 13-15)

Reviewer: qa-reviewer
Date: 2026-08-30
Tasks covered: phase-21/task-13, phase-21/task-14, phase-21/task-15
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-21/qa/review-5.md

## Verdict

PASS

- `PASS`: no Must-fix finding remains; Improvements may remain.
- `REPAIR`: valid Must-fix work has bounded, decision-free ownership in an existing reviewed task. The window does not advance.
- `BLOCK`: protected authority is required, the repair exceeds the reviewed task/window scope, or the bounded no-progress recovery budget is exhausted.

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-21/tasks/task-13-executor-policy.md`
  - `project-docs/docs/phases/phase-21/tasks/task-14-live-output.md`
  - `project-docs/docs/phases/phase-21/tasks/task-15-attempt-receipt.md`
- Completion reports:
  - `project-docs/docs/phases/phase-21/tasks/task-13-done.md`
  - `project-docs/docs/phases/phase-21/tasks/task-14-done.md`
  - `project-docs/docs/phases/phase-21/tasks/task-15-done.md`
- Validation contract: `project-docs/docs/phases/phase-21/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-21/qa/qa-plan.md`
- Repair and report commits: `da6970b4ac5562e20c62f6c72ed8568c10d36f08`, `fdd87aff6e2b4a7f21fff77f5b66423bcbb8b7b1`
- Prior task-15 report commits: `79e3a933cd48568db0866e09255e001497245cca`, `8491cdeed312908e01eb898a3d5795d6000bf746`

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-21-012 (closed by DEC-170; provenance context only)
- Covered validation assertions: VAL-036, VAL-044, VAL-045, VAL-046, VAL-047, VAL-048, VAL-049, VAL-050, VAL-051, VAL-052, VAL-053, VAL-054

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-036 | task-14 | DEFERRED (Human) | `cmd/fid/live_output_test.go` byte-compares the committed machine and Human fixtures. The required HUMAN-01a..01d judgment remains correctly reserved for the Human. |
| VAL-044 | task-13 | ON TRACK | `internal/evidence/executor_test.go` observes structured argv, exact cwd, closed environment, credentials, limits, redaction inputs, and safe Artifact declarations; the exact task-13 package command reproduced PASS. |
| VAL-045 | task-13 | ON TRACK | `cmd/fid/gate_test.go` provides the one-field unsupported-policy/no-Registration matrix; the task-13 report records all required gate results with captured termination. |
| VAL-046 | task-13 | ON TRACK | The two declared containment inventories are recorded as expected `exit_code: 1` no-match proofs, while the accepted explicit values remain `sandbox:none`, `network:allowed`, and `write:allowed`. |
| VAL-047 | task-14 | ON TRACK | `internal/evidence/live_test.go` and the committed `fid.live.v1` golden cover interleaved globally sequenced stdout/stderr data frames. |
| VAL-048 | task-14 | ON TRACK | The live test and golden cover exactly `output-truncated`, `redaction-applied`, and `artifact-rejected` controls in required order. |
| VAL-049 | task-14 | ON TRACK | The fixture tests reject both redaction canaries from machine and Human surfaces and retain the typed redaction control. |
| VAL-050 | task-14 | ON TRACK | The high-output two-pipe test asserts exactly one truncation control, bounded output, drain completion, and observed exit 24. |
| VAL-051 | task-15 | ON TRACK | `TestGateRunAttemptAppendFailureDoesNotSpawn` injects an Attempt append failure, observes no `/bin/touch` side effect, and leaves zero Attempt rows; the success path preserves Event ordering. |
| VAL-052 | task-15 | ON TRACK | The complete binding assertion verifies Receipt identity, execution, terminal, repository, toolchain, output, and redaction fields. A projected update is rejected and every field remains byte-identical. |
| VAL-053 | task-15 | ON TRACK | Failed, spawn-failed, timed-out, and signaled outcomes each replay the original Receipt in one replayed terminal frame, retain exit class 8, and leave one Attempt/Receipt/idempotency row. The handled `fid-error` path emits exit 7 with no Receipt or idempotency result. |
| VAL-054 | task-15 | ON TRACK | A normal run has zero Evidence rows; the injected caller-authored `producer:'caller'`/`trust_source:'fid-executed'` Evidence insert is rejected by the fail-closed schema and the count remains zero. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 13 report versus required commands | `task-13-done.md` against task-13 spec | PASS | All seven required commands are present as PASS with captured termination. The two declared containment controls retain expected and observed `exit_code: 1`. Fresh two-argument report lint ended `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 14 report versus required commands | `task-14-done.md` against task-14 spec | PASS | All five required commands are present as PASS with captured termination. Fresh two-argument report lint ended `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 15 report versus required commands | `task-15-done.md` against task-15 spec | PASS | The unchanged five required commands are all recorded as recovery attempt 6 PASS rows with tokens and session-captured termination. Fresh two-argument report lint ended `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Task 13 reproduced required command | `go test -race ./internal/evidence ./internal/config ./internal/cli ./cmd/fid/...` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 14 reproduced required command | `go test -race ./internal/evidence ./internal/cli ./cmd/fid/...` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 15 reproduced required command | `go test -race ./internal/evidence ./internal/journal ./internal/cli ./cmd/fid/...` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task-15 fresh negative and replay evidence | `go test -race ./cmd/fid -run '^(TestGateRunAttemptAppendFailureDoesNotSpawn|TestGateRunReceiptAppendFailureEmitsFrameOnlyFidError|TestGateRunRecordsExactNonSuccessTerminalOutcomes|TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame)$' -count=1` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 21.100s` (terminated `exit_code: 0`). This includes every required negative element and the four non-success replay cases. |
| Repair provenance and interface boundary | `git show`, `git diff`, `git rev-parse --verify` for `da6970b`, `fdd87af`, `79e3a93`, and `8491cde` | PASS | All hashes resolve. `da6970b` changes only task-15 implementation/tests; `fdd87af` changes only the final task-15 report. Neither changes an accepted contract, public interface, Event schema, policy, migration, validation contract, QA artifact, or task spec. |
| Forward-only report provenance | task-15 history and attempt-6 ledger evidence | PASS | `79e3a93` and `8491cde` remain ancestors and are not rewritten. Their missing footer is truthfully retained in the final report under DEC-170; `da6970b` and `fdd87af` carry `Refs: phase-21/task-15`. The attempt-6 ledger records the exact five-command session evidence and passing report lint. |

## Per-Task Findings

### Task 13 - Portable Executor Policy

#### Must-fix

- None.

#### Improvement

- None. Regression evidence remains aligned with the accepted permissive-policy boundary and the recorded negative containment controls.

### Task 14 - Live Framing, Redaction, And Controls

#### Must-fix

- None.

#### Improvement

- None. The committed golden fixtures remain covered by deterministic byte comparison; the separate Human decision for VAL-036 is not claimed as QA evidence.

### Task 15 - Execution Attempts And Receipts

#### Must-fix

- None.

#### Improvement

- None. The attempt-6 report retains prior failed and unverified history, distinguishes that history from its session-captured recovery evidence, and records the forward-only DEC-170 provenance disposition.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-15-01: non-success Receipts were committed outside idempotency and therefore could spawn/persist again on same-key replay. | task-15 | RESOLVED | `da6970b` removes the non-success `CommitEvent` branch. `internal/evidence/execution.go` now builds every observed Receipt through `CommitIdempotent` and replay restores the persisted Receipt outcome. Fresh focused test PASS proves failed, spawn-failed, timed-out, and signaled requests each produce one replayed terminal with the original Receipt, exit class 8, and exactly one Attempt/Receipt/idempotency row. |
| MF-15-02: missing negative evidence for Attempt append failure, immutable Receipt fields/mutation, handled `fid-error`, and caller proof; prior report commits lacked `Refs` and the prior repair plan required their amendment. | task-15 | RESOLVED | `da6970b` adds all four negative cases and the focused fresh command PASS covers them. The final report correctly records those assertions, attempt-6 captured five-command results, and the prior missing footers as forward-only historical facts. `79e3a93` and `8491cde` remain unrewritten ancestors; DEC-170 establishes that rewriting them would destroy the protected evidence chain. Attempt-6 commits `da6970b` and `fdd87af` each carry `Refs: phase-21/task-15`. |

## Cross-cutting Observations

- The re-review inspected task-13 and task-14 regression evidence as well as the task-15 repair. No reviewed task has a missing or narrower required command, and no QA-plan command beyond those task-spec blocks applies to this window.
- The task-15 repair remains inside its existing task scope. It changes no accepted product or protected delivery artifact; the forward-only exception resolves the prior history requirement without rewriting any commit at or before `09a2a7b`.

## Blocker Candidates

- None.

## Existing Escalation References

- ESC-21-012, closed by DEC-170. It is cited only for the forward-only provenance boundary that resolves the prior MF-15-02 commit-message route; it is not an active blocker candidate.

## Final Lines

Review written: project-docs/docs/phases/phase-21/qa/review-5b.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
