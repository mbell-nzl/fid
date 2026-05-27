# Phase <N> QA Plan - <Phase Title>

> Phase: <N> - <Phase Title>
> Refreshed: <YYYY-MM-DD>
> QA Lead: <name or tool>
>
> Create this file at
> `project-docs/docs/phases/phase-<N>/qa/qa-plan.md`.

## 1. Scope And Source Documents

Summarize the phase goal and list the source artifacts used to build this QA
plan.

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/docs/phases/phase-<N>/qa/validation-contract.md`
- `project-docs/docs/phases/phase-<N>/tasks/task-<M>.md`

## 2. Validation Contract Refinement Check

Confirm the validation contract is ready to be the verification target before
task execution begins. Gaps are QA-plan defects.

- All in-scope requirements referenced by at least one `VAL-NNN`: YES/NO.
- Every `VAL-NNN` has one validator role: YES/NO.
- Every `VAL-NNN` has a concrete evidence method: YES/NO.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES/NO.
- Revisions requested from Tech Lead: <none or list>.

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task <M> | <task focus> | VAL-NNN | <low/medium/high> | <coverage or review note> |

## 4. Required Evidence Strategy

Each row states how QA expects evidence to prove the validation contract. Rows
that do not map to a validation assertion use `no VAL`.

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | <test output, command output, file-state proof, screenshot, API trace, or human observation> | VAL-NNN | <Developer / QA Reviewer / QA Lead / Human> | <task/window/final review> |

## 5. Deterministic Gates

Use concrete project commands after stack selection. If any required command is
still a placeholder when execution is requested, mark the QA plan BLOCKED.

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `<FORMAT_COMMAND>` | PASS | Every task |
| Lint/static analysis | `<LINT_COMMAND>` | PASS | Every task |
| Unit/default tests | `<TEST_COMMAND>` | PASS | Task-level or focused checks |
| Coverage/confidence | `<COVERAGE_COMMAND>` | PASS | Task-level gate |
| Integration | `<INTEGRATION_COMMAND>` | PASS or N/A | Phase-end when applicable |
| End-to-end | `<E2E_COMMAND>` | PASS or N/A | Surface/user-journey phases |
| Build/package | `<BUILD_COMMAND>` | PASS | Phase-end or release candidate |
| Full phase check | `<CHECK_COMMAND>` | PASS | Phase-end gate |

## 6. Command Syntax Validation

Validate every referenced command using help output, dry-run mode, project
documentation, existing scripts, or equivalent proof. Invalid or unvalidated
commands block the QA plan.

| Command | Validation method | Result | Evidence |
|---|---|---|---|
| `<FORMAT_COMMAND>` | <help/dry-run/docs/script inspection> | PASS/FAIL/BLOCKED | <output excerpt or file reference> |

## 7. Deterministic Checklist

Items here are executed or verified by QA Lead and must end with `(qa-lead)`.
Every `QA-NN` cites covered `VAL-NNN` IDs or `no VAL`. PASS without evidence is
invalid.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | <deterministic check> `(qa-lead)` | VAL-NNN | <command output, file-state proof, log, API trace, screenshot, or source reference> |

## 8. Mid-Phase Review Cadence

- Default cadence: QA Reviewer after every 3 completed tasks.
- Exceptions: <none or named checkpoints>.
- Review template:
  `project-docs/templates/mid-phase-review-template.md`.
- Each mid-phase review verifies in-window `VAL-NNN` assertions and records
  gate evidence or why a gate was not run.

## 9. Human Verification Scope

Items here require human judgment and must end with `(human)`. Include only
perception, product-owner judgment, dogfooding, tone, visual hierarchy, or other
human-only checks.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-01 | <manual check> `(human)` | VAL-NNN or no VAL | <observation, decision, screenshot, or notes> |

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- Required deterministic gates pass or are recorded BLOCKED with escalation.
- Every `VAL-NNN` has planned evidence and a final-review outcome slot.
- Every `(qa-lead)` item has reproducible evidence.
- Human checks are human-judgment-only and mapped to the human checklist.
- No unresolved blocking defects remain.
- Required documentation updates are named.
