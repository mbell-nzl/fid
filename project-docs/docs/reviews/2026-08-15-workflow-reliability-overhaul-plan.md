# Workflow Reliability Overhaul: Implementation Plan And Progress

Date: 2026-08-15  
Owner: recovery super-agent acting as principal engineer/architect  
Status: IN PROGRESS  
Branch strategy: preserve the Phase 20 Task 10 worktree, implement the workflow
changes on an isolated branch/worktree, then recover Task 10 under the accepted
workflow.

## Objective

Make routine implementation and verification failures recoverable without human
intervention. Human-visible stops should represent product, architecture,
public-interface, dependency, security, data-loss, or genuinely external
infrastructure concerns, not ordinary red/green work, incomplete in-scope
migrations, evidence bookkeeping, report repair, or long-running command
handling.

## Trigger And Evidence

Phase 20 Task 10 ended after the Developer voluntarily wrote a blocker report.
The child session recorded a normal `task_complete` after about 10 minutes and
no abort, cancellation, interrupt, timeout, or kill. The Developer made an
in-scope test migration, did not rerun after its final patch, and classified the
remaining in-scope work as ESC-20-010. The Phase Runner then accepted that
classification and stopped.

Primary evidence:

- Developer session:
  `/home/mark/.codex/sessions/2026/08/15/rollout-2026-08-15T08-56-09-01a0020f-8f14-7770-ae5e-d5010bd1e7a9.jsonl`
- Task report:
  `project-docs/docs/phases/phase-20/tasks/task-10-done.md`
- Existing recovery instruction:
  `project-docs/specs/09_AI_WORKFLOW.md` section 5.3, step 7.

## Scope

1. Audit recent stoppages and distinguish true decisions from recoverable work,
   evidence defects, orchestration defects, and external infrastructure faults.
2. Define a binding task/review recovery state machine and blocker taxonomy.
3. Patch the v1 workflow spec, canonical agent prompts, templates, lints, and
   deterministic drills.
4. Add reliable orchestration rules for child blocker classification, bounded
   recovery, fresh-session continuation, QA rework, and long-running gates.
5. Obtain independent review of the design and final implementation.
6. Recover and complete Phase 20 Task 10 without broadening its product scope.
7. Record the v2 starter-kit adoption work without editing the separate starter
   kit during this repository recovery unless an explicit, isolated follow-up is
   justified by the completed audit.
8. Produce an HTML summary suitable for non-terminal review.

## Non-Goals

- No redesign of fid product contracts or Phase 20 acceptance criteria.
- No weakening, narrowing, or silent substitution of deterministic gates.
- No broad test-performance refactor during Task 10 recovery.
- No erasure of failed attempts or conversion of unverified evidence to PASS.
- No reliance on chat history as the execution record.

## Target State Machine

Normal task states:

`ACTIVE -> RED -> REPAIR -> VERIFY -> REPORT -> COMPLETE`

Recovery may loop from `RED`, `VERIFY`, or `REPORT` back to `REPAIR`. A fresh
Developer session may continue the same task from a committed or otherwise
durable checkpoint without converting recovery into a human escalation.

A task may enter `BLOCKED` only when the blocking condition is outside the
task's authorized scope or capabilities and matches a protected category. A
blocker report must name the condition, category, evidence, attempted recovery,
and why further work cannot proceed within scope.

## Protected Blocker Categories

- Product, architecture, data-model, public-interface, or terminology decision.
- Dependency or tool addition without an accepted budget and approval path.
- Security concern, possible data loss, or destructive migration uncertainty.
- External service, permission, or infrastructure failure that cannot be
  repaired locally.
- Repeated deterministic failure after bounded causal recovery attempts, with
  evidence that continuing would only repeat the same action.
- Required specification or validation contract that is absent or internally
  contradictory.

The following are recovery states, not blockers by themselves:

- Task implementation is incomplete.
- A test fails because the assigned change invalidated an old expectation.
- A focused or full gate exposes an in-scope defect.
- A completion report, inventory, generated provider copy, or evidence receipt
  needs repair.
- A command is silent while still running, or a live session must be polled.
- Context is becoming constrained and a fresh task continuation is needed.
- A bounded QA finding has an unambiguous, decision-free fix.

## Implementation Workstreams

### A. Preservation And Isolation

- Commit this plan as the durable anchor.
- Capture the complete dirty-tree inventory and stash all tracked and untracked
  Task 10/Phase Runner work with a unique recovery label.
- Record the stash object and base commit here.
- Create an isolated `workflow-reliability-overhaul` branch/worktree.

### B. Audit And Decision Record

- Audit Phase 17-20 stoppage artifacts and relevant agent session evidence.
- Classify each stoppage by root cause and human-intervention necessity.
- Record the recovery state machine and escalation policy as a formal decision.
- Identify v2 starter-kit adoption deltas.

### C. V1 Workflow Implementation

- Update `09_AI_WORKFLOW.md` and `AGENTS.md` source where applicable.
- Update canonical Developer, Phase Runner, QA Reviewer, QA Lead, and TL-Lite
  prompts while preserving role boundaries.
- Update Phase Runner, run-ledger, escalation, completion-report, and review
  templates.
- Add machine-checkable blocker classification and recovery rules to workflow
  lint/report lint where they can be deterministic.
- Sync generated provider prompts and skills.

### D. Deterministic Drills

Cover at least:

- In-scope focused test failure followed by repair and byte-identical rerun.
- Full gate exposing a stale in-scope assertion.
- Completion-report or inventory repair.
- Long silent command with a live session handle.
- Developer returning an invalid “incomplete work” blocker.
- Fresh-session continuation of the same task.
- Bounded, decision-free QA rework and re-review.
- Genuine protected decision requiring a human stop.

### E. Review And Integration

- Request an early independent architecture review after the state machine and
  policy are drafted.
- Request a final independent review after implementation and gates.
- Resolve findings or record a justified disposition.
- Merge/cherry-pick the workflow commits into the Phase 20 line.
- Restore Task 10 from the recorded checkpoint and finish it under the new
  workflow.

### F. Verification And Reporting

- Run prompt/skill generation and workflow lint.
- Run report-lint self-tests and workflow drills.
- Run all exact Task 10 product gates in the session recording their results.
- Update Phase 20 records to close or supersede ESC-20-010 accurately.
- Produce an HTML report and copy it to the Windows Downloads folder.

## Progress Ledger

| Time (NZST) | State | Evidence / result | Next action |
|---|---|---|---|
| 2026-08-15 start | IN PROGRESS | User authorized autonomous overhaul, independent review, regular commits, and final HTML report. | Commit this plan and preserve the dirty tree. |
| 2026-08-15 checkpoint 1 | COMPLETE | Plan committed as `6517c81`; all remaining tracked and untracked changes preserved in stash `779beab1` (`stash@{0}` at capture), covering 28 files with 1,052 insertions and 230 deletions. Original worktree verified clean. | Create the isolated overhaul branch/worktree. |
| 2026-08-15 checkpoint 2 | IN PROGRESS | Created branch `workflow-reliability-overhaul` and worktree `/tmp/fid-workflow-overhaul` from `8951011`. Started independent audits of stoppage history, v1 workflow mechanics, and v2 adoption deltas. | Draft the audit and binding workflow decision. |

## Checkpoints

| Checkpoint | Reference | Status |
|---|---|---|
| Phase 20 base commit | `db36832` | RECORDED |
| Durable plan commit | `6517c81` | COMPLETE |
| Dirty-tree recovery stash | `779beab1` (`recovery/phase-20-task-10-and-runner-2026-08-15`) | COMPLETE |
| Overhaul branch/worktree | `workflow-reliability-overhaul` at `/tmp/fid-workflow-overhaul` | COMPLETE |
| Audit and decision commit | pending | NOT STARTED |
| V1 implementation commits | pending | NOT STARTED |
| Independent review | pending | NOT STARTED |
| Task 10 recovery commit(s) | pending | NOT STARTED |
| Final HTML report | pending | NOT STARTED |

## Continuation Instructions

On any compaction or fresh session:

1. Read this document first.
2. Inspect the latest `Progress Ledger` and `Checkpoints` entries.
3. Verify branch, worktree, stash, and commit state before editing.
4. Continue the first incomplete workstream; do not restart completed analysis.
5. Update this ledger and commit at each material boundary.
