# Phase 21 - Mid-Phase Review 5 (Tasks 13-15)

Reviewer: qa-reviewer
Date: 2026-08-30
Tasks covered: phase-21/task-13 through phase-21/task-15
Review type: INITIAL
Prior review: none

## Verdict

REPAIR

- `PASS`: no Must-fix finding remains; Improvements may remain.
- `REPAIR`: valid Must-fix work has bounded, decision-free ownership in an
  existing reviewed task. The window does not advance.
- `BLOCK`: protected authority is required, the repair exceeds the reviewed
  task/window scope, or the bounded no-progress recovery budget is exhausted.

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-21/tasks/task-13-executor-policy.md`
  - `project-docs/docs/phases/phase-21/tasks/task-14-live-output.md`
  - `project-docs/docs/phases/phase-21/tasks/task-15-attempt-receipt.md`
- Completion reports:
  - `project-docs/docs/phases/phase-21/tasks/task-13-done.md`
  - `project-docs/docs/phases/phase-21/tasks/task-14-done.md`
  - `project-docs/docs/phases/phase-21/tasks/task-15-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-21/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-21/qa/qa-plan.md`
- Commits or changed files:
  - `d2feab0be5c53f46553cc18b79e605e6586da434`, `f0eec233bba8695c2b33b5a74237fc0eea591ff7`, and `b69499116e70d4c88423a30daa2010a4ff221f14`
  - committed task-report follow-ups `f139ecf`, `cabd0d4`, `79e3a93`, and `8491cde`

## Summary

- Review outcome: REPAIR
- Must-fix findings: 2
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-044, VAL-045, VAL-046, VAL-047, VAL-048, VAL-049, VAL-050, VAL-051, VAL-052, VAL-053, VAL-054

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-044 | task-13 | ON TRACK | `TestExecuteAppliesAcceptedGatePolicy` observes a real child argv, cwd, and closed environment; the independently rerun exact package gate ended `exit_code: 0`. |
| VAL-045 | task-13 | ON TRACK | `cmd/fid/gate_test.go:71` covers one-field rejected policy values and checks `executor-policy-unsupported` plus zero Registration rows. |
| VAL-046 | task-13 | ON TRACK | Both declared containment inventories were independently reproduced as `exit 1 (no output)` with observed `exit_code: 1`; the policy test uses only `sandbox:none`, `network:allowed`, and `write:allowed`. |
| VAL-047 | task-14 | ON TRACK | `internal/evidence/live_test.go:11` and the committed `live-frames.jsonl` fixture show ordered stdout/stderr frames, escaped controls, and sequential frame numbers. |
| VAL-048 | task-14 | ON TRACK | The committed fixture and `TestRunLiveFramesInterleavedOutputControlsAndFailure` contain only `output-truncated`, `redaction-applied`, and `artifact-rejected` in the required order. |
| VAL-049 | task-14 | ON TRACK | Live and Human fixture canaries are absent and a typed redaction control is present; `internal/evidence/execution.go:359` derives retained output from those redacted data frames. Artifact and export writers are not yet implemented in this window. |
| VAL-050 | task-14 | ON TRACK | `TestRunLiveDrainsBothPipesAfterTruncationAndPreservesOutcome` uses a 256 KiB two-pipe process and asserts one control, bounded retained data, and observed exit 24. |
| VAL-051 | task-15 | NOT ON TRACK | The success path checks journal ordering, but no injected Attempt-append failure proves that no spawn occurs; QA-49 requires both states. |
| VAL-052 | task-15 | NOT ON TRACK | The success/replay test counts rows, but it does not attempt a Receipt mutation or assert the complete immutable binding and terminal fields required by QA-50. |
| VAL-053 | task-15 | NOT ON TRACK | The four process outcomes have a first-run matrix, but non-success Receipts are not idempotently replayable; the handled `fid-error` test is frame-only and does not exercise CLI exit 7. |
| VAL-054 | task-15 | NOT ON TRACK | The test proves a successful gate run creates no Evidence row, but it supplies no caller-reported proof and does not observe the required rejected or unsatisfied result. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 13 exact-command coverage and report structure | `task-13-done.md` against `task-13-executor-policy.md` | PASS | All seven exact commands are recorded with tokens and observed termination. Both declared containment commands record expected and observed `exit_code: 1`; independent report lint ended `PASS report-lint completed with no failures` (`exit_code: 0`). |
| Task 14 exact-command coverage and report structure | `task-14-done.md` against `task-14-live-output.md` | PASS | All five exact commands are recorded with tokens and observed termination; independent report lint ended `PASS report-lint completed with no failures` (`exit_code: 0`). |
| Task 15 exact-command coverage and report structure | `task-15-done.md` against `task-15-attempt-receipt.md` | PASS | All five exact commands are recorded with tokens and observed termination; the extra tagged integration result does not substitute for a required command; independent report lint ended `PASS report-lint completed with no failures` (`exit_code: 0`). |
| Task 13 reproduced required command | `go test -race ./internal/evidence ./internal/config ./internal/cli ./cmd/fid/...` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 14 reproduced required command | `go test -race ./internal/evidence ./internal/cli ./cmd/fid/...` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 15 reproduced required command | `go test -race ./internal/evidence ./internal/journal ./internal/cli ./cmd/fid/...` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid (cached)` (terminated `exit_code: 0`). |
| Task 13 reproduced containment controls | the two verbatim `rg -n -i` commands from task 13 | PASS | Each produced `exit 1 (no output)` and terminated `exit_code: 1`, matching the declared expected non-zero condition. |
| QA-42 | executor positive process probe | PASS | `executor_test.go:16` compares literal argv, resolved cwd, the complete child environment, credential injection, and exit 0. |
| QA-43 | unsupported-policy registration matrix | PASS | `gate_test.go:71` covers argv, shell, sandbox, network, write, environment, limits, credentials, redaction, Artifact values, and cwd without persistence. |
| QA-44 | containment inventory plus permissive policy values | PASS | Reproduced source and contract no-match controls above; `policy.go` admits only the three explicit permissive values. |
| QA-45 | committed live-frame golden | PASS | `f0eec23:cmd/fid/testdata/val-036/live-frames.jsonl` has `fid.live.v1`, one execution ID, sequences 1 through 9, interleaved streams, and JSON-escaped ESC/CR/BEL. |
| QA-46 | control-frame golden | PASS | The same committed fixture has exactly one `output-truncated`, `redaction-applied`, and `artifact-rejected` control with contract-shaped payloads; `live.go` creates controls internally from the closed union. |
| QA-47 | pre-surface redaction | PASS | The committed machine and Human fixtures contain `[REDACTED]` and no source canary; current stored/hashing flow consumes the already-redacted live data. No Artifact or export writer exists at this review HEAD. |
| QA-48 | bounded two-pipe drain | PASS | `live_test.go:98` asserts a single truncation frame, at-most-limit emitted text, drain completion, and observed non-success exit 24. |
| QA-49 | Attempt-before-spawn plus append-failure negative | FAIL | `gate_run_test.go:69` establishes positive Event ordering only. There is no injected append failure or spawn counter, as required by QA-49. |
| QA-50 | immutable Receipt | FAIL | `gate_run_test.go:78` counts one Receipt on successful replay but has no mutation attempt or complete binding/terminal-field assertion. |
| QA-51 | five terminal outcomes | FAIL | Focused matrix rerun passed: `ok github.com/mbell-nzl/fid/cmd/fid 12.565s` (`exit_code: 0`), but it contains no non-success replay and `live_test.go:139` is not a CLI exit-7 test for handled `fid-error`. |
| QA-52 | Attempt-only and caller-proof rejection | FAIL | `gate_run_test.go:85` only counts zero Evidence rows after a normal run; no caller-reported proof is constructed or rejected. |

## Per-Task Findings

### Task 13 - Portable Executor Policy

#### Must-fix

- None.

#### Improvement

- None.

### Task 14 - Live Framing, Redaction, And Controls

#### Must-fix

- None.

#### Improvement

- None. The committed `human-session.txt` and `live-frames.jsonl` are byte-compared by `cmd/fid/live_output_test.go:14`; VAL-036 remains a separate Human-pending assertion.

### Task 15 - Execution Attempts And Receipts

#### Must-fix

- [internal/evidence/execution.go:334] Observed non-success outcomes use `journal.CommitEvent` and return before `CommitIdempotent`; the replay lookup at `execution.go:119` therefore cannot find them. A repeated failed, timed-out, signaled, or spawn-failed request can spawn again and persist another Attempt/Receipt, contrary to the accepted `gate.run` replay contract. Owning task: task-15. Repair boundary: persist the accepted terminal idempotency record for every observed process outcome while retaining the existing failure-envelope/no-outer-Event behavior. Required evidence: a five-outcome process matrix that repeats each same key/digest, asserts no second spawn/Attempt/Receipt, one terminal frame with `replayed:true`, the original Receipt binding, and exit class 8. Re-review scope: inspect `internal/evidence/execution.go`, `cmd/fid/gate.go`, and rerun the task-15 exact package command plus the non-success replay matrix. (VAL-053)
- [cmd/fid/gate_run_test.go:69] The report claims VAL-051 through VAL-054 pass at `task-15-done.md:40`, but its tests omit the required non-qualifying evidence: injected Attempt-append failure with zero spawn, Receipt mutation rejection and exact immutable fields, a CLI-level handled `fid-error` terminal with exit 7/no Receipt/no idempotency, and caller-reported proof rejection or unsatisfied Evidence. In addition, neither task-15 report commit (`79e3a93` and `8491cde`) carries the required `Refs: phase-21/task-15` footer. Owning task: task-15. Repair boundary: add only task-15 tests and required test seams under `internal/evidence/` and `cmd/fid/`, correct the completion report to cite the resulting evidence, and amend the two report commit messages to add the missing footer without changing task scope or accepted contracts. Required evidence: the QA-49 through QA-52 negative and positive cases, the byte-identical task-15 exact command, successful two-argument report lint, and `git show -s --format=fuller` for both report commits showing the footer. Re-review scope: inspect the amended report history and test assertions, then rerun the task-15 exact package gate and focused terminal/error cases. (VAL-051, VAL-052, VAL-053, VAL-054; no VAL for report provenance)

#### Improvement

- None.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-15-01 | task-15 | VAL-053 | In `internal/evidence/execution.go` and task-15 tests, store/replay every observed non-success Receipt without changing the public failure envelope or exit class. | Repeat failed, timed-out, signaled, and spawn-failed same-key requests; prove no second spawn, Attempt, or Receipt and one replayed terminal frame; task-15 exact command passes. | Inspect the non-success idempotency path and rerun the exact task-15 package gate plus the replay matrix. |
| MF-15-02 | task-15 | VAL-051, VAL-052, VAL-053, VAL-054; no VAL for report provenance | Add the missing append-failure, immutable Receipt, handled-fid-error, and caller-proof cases; correct the report claims and amend the two task-15 report-commit messages with `Refs: phase-21/task-15`. | Each QA-49 through QA-52 case has qualifying and non-qualifying proof; task-15 exact command and report lint pass; both report commits show their Refs footer. | Inspect tests, report, and commit metadata; rerun the exact task-15 package gate and focused affected tests. |

## Blocker Candidates

- None.

## Existing Escalation References

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-21/qa/review-5.md
Must-fix: 2  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
