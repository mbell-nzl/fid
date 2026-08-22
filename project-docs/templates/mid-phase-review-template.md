# Phase <N> - Mid-Phase Review <R> (Tasks <X>-<Z>)

Reviewer: QA Reviewer
Date: <YYYY-MM-DD>
Tasks covered: phase-<N>/task-<X> through phase-<N>/task-<Z>
Review type: INITIAL / FRESH INDEPENDENT RE-REVIEW
Prior review: <path or none>

## Verdict

PASS / REPAIR / BLOCK

- `PASS`: no Must-fix finding remains; Improvements may remain.
- `REPAIR`: valid Must-fix work has bounded, decision-free ownership in an
  existing reviewed task. The window does not advance.
- `BLOCK`: protected authority is required, the repair exceeds the reviewed
  task/window scope, or the bounded no-progress recovery budget is exhausted.
  A failed gate, incomplete in-scope work, or a valid Must-fix finding is not
  `BLOCK` by itself.

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-<N>/tasks/task-<M>.md`
- Completion reports:
  - `project-docs/docs/phases/phase-<N>/tasks/task-<M>-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-<N>/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-<N>/qa/qa-plan.md`
- Commits or changed files:
  - `<hash or path>`

## Summary

- Review outcome: <PASS / REPAIR / BLOCK>
- Must-fix findings: <count>
- Improvement findings: <count>
- Blocker candidates: <count>
- Existing escalation references: <ESC-<phase>-<NNN> list or none>
- Covered validation assertions: <VAL-NNN list or none>

## Validation Contract Coverage

List each `VAL-NNN` in the reviewed task window and whether the evidence
tracks toward the pass criterion. Missing evidence is Must-fix unless the
validation contract or QA plan explicitly assigns the assertion to a later
window. Classify it as `REPAIR` when the evidence can be produced inside an
owning task; use `BLOCK` only when protected authority is required.

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-NNN | task-<M> | ON TRACK / NOT ON TRACK / DEFERRED | <test output, file reference, commit, screenshot, or rationale> |

## Gate And Evidence Checks

Compare each reviewed completion report against the task spec's exact commands
and the QA plan's required commands. Missing required gates, silently skipped
gates, or narrower substituted commands are Must-fix findings. A report's
`BLOCKED` label does not establish a protected blocker without the matching
category and authority boundary.

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| <format/lint/test/source-state/manual artifact> | `<command>` or `<artifact>` | PASS/FAIL/NOT RUN | <reason or output excerpt> |

## Per-Task Findings

### Task <M> - <Title>

#### Must-fix

- None.

#### Improvement

- None.

Repeat for each task in the review window. Findings should include severity,
file or artifact reference where possible, what is wrong, and the expected fix
direction. Every Must-fix finding must also appear in the Repair Plan or Blocker
Candidates section according to the review outcome.

## Prior Finding Resolution

Required for a fresh independent re-review. Preserve every finding from the
prior review and resolve it only from newly inspected evidence. The re-review
must rerun representative affected verification; do not append corrections to
the prior review.

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| <finding ID> | task-<M> | RESOLVED / STILL OPEN / PROTECTED BLOCKER CANDIDATE | <command result, diff, commit, or artifact> |

## Repair Plan

Required for outcome `REPAIR`. Every Must-fix finding needs one row. The owning
task must already be in the reviewed window, and the repair must not change
accepted product, architecture, data-model, public-interface, dependency,
security, or validation-contract decisions.

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| <finding ID> | task-<M> | VAL-NNN / no VAL | <bounded files/behavior> | <exact test, gate, report, or file-state proof> | <what the new reviewer must inspect and rerun> |

## Blocker Candidates

- None.

If a finding crosses task boundaries, has architecture/product/security
implications, contradicts a standing decision, exposes ambiguous acceptance
criteria, or is a quality-gate anomaly that cannot be repaired without changing
the gate contract, record a blocker candidate here. Include the protected
category, evidence, why it is outside existing repair authority, the authority
or external action required, and any existing `ESC-<phase>-<NNN>` reference. QA
roles do not allocate `ESC` identities or edit escalation registers.

## Existing Escalation References

- None.

Cite an existing `ESC-<phase>-<NNN>` only when it already covers a finding or
blocker candidate. Do not create one.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-<N>/qa/review-<R>.md
Must-fix: <count>  Improvement: <count>  Blocker candidates: <count>
VERDICT: PASS | REPAIR | BLOCK
```
