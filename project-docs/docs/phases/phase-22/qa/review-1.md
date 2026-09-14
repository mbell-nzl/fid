# Phase 22 - Mid-Phase Review 1 (tasks 1-3)

Reviewer: qa-reviewer
Date: 2026-09-05
Tasks covered: phase-22/task-1, phase-22/task-2, phase-22/task-3
Review type: INITIAL
Prior review: none

## Verdict

REPAIR

The window cannot advance until the bounded task-2 and task-3 evidence repairs below receive a fresh independent re-review. No protected decision, gate definition, public-interface change, dependency, or escalation is needed.

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-22/tasks/task-1.md`
  - `project-docs/docs/phases/phase-22/tasks/task-2.md`
  - `project-docs/docs/phases/phase-22/tasks/task-3.md`
- Completion reports:
  - `project-docs/docs/phases/phase-22/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-22/tasks/task-2-done.md`
  - `project-docs/docs/phases/phase-22/tasks/task-3-done.md`
- Validation contract and QA plan:
  - `project-docs/docs/phases/phase-22/qa/validation-contract.md`
  - `project-docs/docs/phases/phase-22/qa/qa-plan.md`
  - `project-docs/docs/phases/phase-22/qa/tl-refinement.md`
- Commits and changed files:
  - `72c5974f6698fa08a4fe1a24099b3a93b3882f1d` — `feat(recovery): add schema and replay foundation`
  - `d8a546988ea4f27872c2e87352332b02cb349a22` — `docs(task): finalize completion report for phase-22 task-1`
  - `ac827959f8359473cbd6546ac6bea1ff1fd20751` — `feat(workflow): add typed Stop and Attention facts`
  - `f3936a041cb4a55a782c6145049b5091db0af539` — `docs(task): finalize completion report for phase-22 task-2`
  - `7e419bfd8c233f933680d4a3b3957814c29dbd5e` — `feat(workflow): add lease expiry and executor observations`
  - `6715c73f3c478c2a3244dbe6c41c91aa933b7901` — `docs(task): finalize completion report for phase-22 task-3`

All six supplied hashes resolve with `git rev-parse --verify <hash>^{commit}`. Their subjects and `Refs: phase-22/task-N` footers match the reported task. `git diff --check` is clean. The sole pre-existing worktree item is the Runner-owned untracked `project-docs/docs/phases/phase-22/run-ledger.md`, which this review did not modify.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 3
- Improvement findings: 1
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-006, VAL-007, VAL-008, VAL-009, VAL-010, VAL-011, VAL-012, VAL-079

The review cadence is correct: window 1 covers tasks 1-3 and produces `review-1.md`; the next window is tasks 4-6 and produces review 2 after task 6. This agrees with `execution-scope.md` and QA-plan section 11.4.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | `internal/db/recovery_migration_test.go` checks schema-15 tables, non-null recovery observation, accepted and invalid Stop/Attention codes; the tagged binary test checks fresh init and rebuild. |
| VAL-002 | task-1 | ON TRACK | The same migration test executes down to zero and up again, then checks the restored empty recovery tables. |
| VAL-003 | task-1 | ON TRACK | `internal/journal/recovery_schema_test.go` exercises both live apply and replay rejection for unavailable recovery effects and corrupt Stop codes. |
| VAL-004 | task-2 | ON TRACK | `TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable` retains the original Stop code and details while recording resolution separately. |
| VAL-005 | task-2 | ON TRACK | Direct workflow tests show two distinct Stops share one open Attention identity and preserve both links. |
| VAL-006 | task-2 | ON TRACK | The direct resolution test rejects a missing disposition set without a journal change and accepts the exact sorted set. |
| VAL-007 | task-2 | NOT ON TRACK | The report proves explicit Stop creation and a rejected `GetStop` read, but not the QA-plan-required routine status and warning paths. The non-qualifying evidence is incomplete. |
| VAL-079 | task-2 | ON TRACK | The report preserves pre-edit 27.4% and final 34.7% `internal/workflow` coverage results with direct positive and rejection tests. |
| VAL-008 | task-3 | ON TRACK | Clock-controlled before/equal/after expiry testing and the explicit expiry/reacquisition fence cover qualifying and rejection states. |
| VAL-009 | task-3 | ON TRACK | The tagged `TestPhase21RealProcessJourneyStartsOneActivatedEnvelope` starts two `fid` processes and asserts one winner, a typed conflict loser, and one Event. |
| VAL-010 | task-3 | ON TRACK | The observer test accepts the bound harness and rejects both an ordinary Builder and a forged executor without state change. |
| VAL-011 | task-3 | NOT ON TRACK | The cited observation test only fences stale **Run** versions. The existing renewal test fences a Claim version but creates no Run. Neither proves renewal retains Run lineage while only the current Claim version authorizes the relevant operation. |
| VAL-012 | task-3 | ON TRACK | Start-binding tests cover valid binding, mismatched Issue/Envelope/revision/actor/repository cases, optional Obligation, and the exact reused-executor duplicate-entity envelope. |

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Task 1 report lint | `project-docs/tools/report-lint.sh task-1-done.md task-1.md` | PASS | All headings, required command receipts, reported implementation hash/subject, and clean report commit verified. |
| Task 2 report lint | `project-docs/tools/report-lint.sh task-2-done.md task-2.md` | PASS | All headings, baseline/final coverage receipts, required command receipts, implementation hash/subject, and clean report commit verified. |
| Task 3 report lint | `project-docs/tools/report-lint.sh task-3-done.md task-3.md` | PASS | All headings, required command receipts, reported implementation hash/subject, and clean report commit verified. Mechanical lint does not validate the semantic VAL-011 attribution. |
| Task 1 exact commands | format, lint, race/coverage, cover, tagged integration, environment inventory, and schema-15 rebuild benchmark | PASS | Every task-spec command is recorded with a pass token and observed `exit_code: 0`; the earlier failed race attempt remains preserved in the report. |
| Task 2 exact commands | workflow coverage baseline/final, format, lint, race/coverage, cover, tagged integration | PASS | Every command is present. The report preserves initial lint and full-gate failures and later byte-identical passing attempts. |
| Task 3 exact commands | format, lint, race/coverage, cover, tagged integration | PASS | Every task-spec command is present with a pass token and observed `exit_code: 0`; the diagnostic non-gate failure remains recorded as failed evidence. |
| Representative reproduced evidence | `go test -race -cover ./internal/workflow` | PASS | `ok  github.com/mbell-nzl/fid/internal/workflow  31.674s  coverage: 32.9% of statements` (terminated `exit_code: 0`). This independently rerun required task-2 command completed normally. |
| Retained-test provenance | task-1 and task-2 report citations plus `git log -S` / `git blame` | PASS | The cited Phase 21 introducing commits resolve and predate the phase-22 base; reports preserve rather than relabel earlier failed attempts. |

## Per-Task Findings

### Task 1 - Recovery schema and replay foundation

#### Must-fix

- None.

#### Improvement

- None. The migration/replay tests assert observable schema constraints and both live/replay rejection paths; changed rebuild/export documentation behavior agrees with the active schema constant.

### Task 2 - Typed Stops and Attention domain

#### Must-fix

- [task-2.md:47] No tagged `cmd/fid/integration/` contract fixture covers the new Stop/Attention behavior. The report's changed-file inventory contains no integration fixture, and the integration package search finds only task-1's schema mention of the tables. Owning task: task-2. Repair boundary: add a tagged integration regression fixture within `cmd/fid/integration/` that exercises the accepted external boundary without adding the deferred public Stop/Attention leaves. Required evidence: focused tagged test output plus the unchanged required task-2 gates and an updated completion report. Re-review scope: inspect the fixture's asserted public boundary and rerun it and the representative affected gate.
- [task-2-done.md:38] VAL-007 is reported PASS from an explicit `RaiseStop` case and rejected `GetStop` read, but neither the report nor its cited test covers the required routine status and warning paths. Owning task: task-2. Repair boundary: add direct observable tests for routine status and warning behavior showing unchanged Attention count/journal state, alongside the qualifying real-Stop case; correct the report's VAL-007 evidence. Required evidence: focused test output proving both non-qualifying paths and the real-Stop path, then all existing task-2 exact commands. Re-review scope: inspect the cases, report correction, and rerun the focused test and `go test -race -cover ./internal/workflow`.

#### Improvement

- [internal/workflow/stop.go:1] `stop.go` is 668 lines and contains both Stop and Attention types, persistence, validation, and operations, exceeding the repository's 400-line non-test file target. Fix direction: split Attention-specific behavior into the already named `internal/workflow/attention.go` when a subsequent in-scope change makes the extraction economical.

### Task 3 - Lease lifecycle and executor observations

#### Must-fix

- [task-3-done.md:49] The VAL-011 evidence cites `TestRunObserveRequiresTrustedBoundHarnessAndCurrentVersion`, whose stale check is for `Run.version`; its setup never renews a Claim. `TestClaimRenewReleaseVersionAndReplayAreExact` renews and fences Claim versions but never creates or inspects a Run. Owning task: task-3. Repair boundary: add a single end-to-end task-owned scenario that starts a Run, renews its Claim, proves the Run lineage/binding remains immutable, and proves current versus stale Claim-version admission as the accepted contract requires; update the VAL-011 report attribution. Required evidence: focused red/green test receipt, final exact task-3 gates, and corrected completion-report evidence. Re-review scope: inspect the new scenario and report correction; rerun it, the tagged integration suite, and a representative required task command.

#### Improvement

- None. The expiry, real-process race, observer authority, executor identity, error-envelope, and retained-fixture behavior otherwise align with the task-3 scope and contract.

## Prior Finding Resolution

Not applicable: this is the initial review.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-1 | task-2 | no VAL | Add the scope-required tagged Stop/Attention boundary fixture in `cmd/fid/integration/`; do not add deferred public leaves. | Focused tagged receipt, unchanged exact task-2 gates, corrected report. | Fixture assertion, report, and representative affected command. |
| MF-2 | task-2 | VAL-007 | Test routine status and warning paths with unchanged Attention/journal state, paired with explicit Stop creation; correct the report. | Focused two-state receipt and all task-2 exact commands. | Cases, report attribution, and `go test -race -cover ./internal/workflow`. |
| MF-3 | task-3 | VAL-011 | Test Claim renewal after Run creation, immutable Run lineage/binding, and current/stale Claim-version admission; correct the report. | Focused red/green receipt, all task-3 exact commands, corrected report. | New test, report attribution, tagged suite, and representative required command. |

## Blocker Candidates

- None. Each finding has decision-free ownership inside an existing reviewed task and leaves every accepted command, assertion semantic, interface, schema, dependency, and policy unchanged.

## Existing Escalation References

- None. `project-docs/docs/current/tl_escalations_current.md` was reread immediately before this review was written and still lists no open ESCs.

## Cross-cutting observations

- The completion reports retain original red and diagnostic failed attempts with observed exit codes rather than presenting later passes as clean histories.
- Task reports and the actual diffs show no dependency additions. The changed domain behavior remains in the workflow/claims/journal layers; CLI changes parse and render the accepted existing surfaces.
- The phase remains on the planned six-window cadence. A fresh independent reviewer, rather than this reviewer, must assess closure of every Must-fix item after the Developer continuations.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-22/qa/review-1.md
Must-fix: 3  Improvement: 1  Blocker candidates: 0
VERDICT: REPAIR
```
