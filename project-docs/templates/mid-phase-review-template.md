# Phase <N> - Mid-Phase Review <R> (Tasks <X>-<Z>)

Reviewer: QA Reviewer
Date: <YYYY-MM-DD>
Tasks covered: phase-<N>/task-<X> through phase-<N>/task-<Z>

## Verdict

CONTINUE / BLOCK

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

- Blocking findings: <count>
- Improvement findings: <count>
- Escalations raised: <ESC-NNN list or none>
- Covered validation assertions: <VAL-NNN list or none>

## Validation Contract Coverage

List each `VAL-NNN` in the reviewed task window and whether the evidence
tracks toward the pass criterion. Missing evidence is Blocking unless the
validation contract or QA plan explicitly assigns the assertion to a later
window.

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-NNN | task-<M> | ON TRACK / NOT ON TRACK / DEFERRED | <test output, file reference, commit, screenshot, or rationale> |

## Gate And Evidence Checks

Compare each reviewed completion report against the task spec's exact commands
and the QA plan's required commands. Missing required gates, silently skipped
gates without a `BLOCKED` reason, or narrower substituted commands are Blocking
findings.

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| <format/lint/test/source-state/manual artifact> | `<command>` or `<artifact>` | PASS/FAIL/NOT RUN | <reason or output excerpt> |

## Per-Task Findings

### Task <M> - <Title>

#### Blocking

- None.

#### Improvement

- None.

Repeat for each task in the review window. Findings should include severity,
file or artifact reference where possible, what is wrong, and the expected fix
direction.

## Escalations

- None.

If a finding crosses task boundaries, has architecture/product/security
implications, contradicts a standing decision, exposes ambiguous acceptance
criteria, or is a quality-gate anomaly not anticipated by the QA plan, append an
`ESC-NNN` entry to `project-docs/docs/current/tl_escalations_current.md` and
reference it here.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-<N>/qa/review-<R>.md
Blocking: <count>  Improvement: <count>  Escalations: <list or none>
VERDICT: CONTINUE | BLOCK
```
