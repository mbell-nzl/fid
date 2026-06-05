# Phase <N> Validation Contract - <Phase Title>

> Phase: <N> - <Phase Title>
> Author: Tech Lead (draft)
> Last refined: <YYYY-MM-DD> by QA Lead
>
> Create this file at
> `project-docs/docs/phases/phase-<N>/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
the phase is complete. The validation contract is the specification of **what
must be proven**; the QA plan is the strategy for **how to prove it**.

Tech Lead drafts this contract before or alongside task specs. QA Lead refines
the evidence strategy before implementation starts. Developers cite the
`VAL-NNN` IDs their task advances, QA Reviewer verifies the IDs in its task
window, and QA Lead records final outcomes at phase end.

If a required validation assertion cannot be satisfied, stop and escalate
instead of weakening the assertion or silently adapting the spec.

## Numbering

- IDs use `VAL-NNN`, scoped to one phase and sequential from `VAL-001`.
- Carry-forward assertions are re-issued in the receiving phase with a new
  `VAL-NNN` and a source note that points to the earlier assertion.
- Retired or superseded assertions stay in the table with status and rationale;
  do not reuse IDs.

## Assertions

Each row is one testable behavior, invariant, or release-readiness assertion.
Prefer observable behavior over implementation phrasing. Do not bundle several
independent observables into one `VAL-NNN`.

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | <PRD/TRD/design/quality section, DEC-NNN, or phase brief item> | <observable behavior or invariant> | <test, command, source-state proof, screenshot, human check, API trace, or other reproducible evidence> | <QA Lead / QA Reviewer / Human> | <task-M / tasks-M-N / phase-end> | <specific condition that proves the assertion> | <specific condition that disproves or blocks the assertion> | <Draft / Ready / PASS / FAIL / BLOCKED / CARRY-FORWARD> | <related behavior intentionally excluded> |

## Field Guidance

- **Source requirement:** Cite the durable source that justifies the assertion.
  Phase-only assertions may cite the phase brief.
- **Behavior or invariant:** State one thing a user, caller, operator, file, or
  system boundary must observe.
- **Evidence method:** Name the smallest reproducible proof available for the
  assertion. QA Lead turns this into concrete QA-plan steps.
- **Validator:** Assign exactly one role: `QA Lead`, `QA Reviewer`, or `Human`.
  Use `Human` only for perception, product-owner judgment, dogfooding, tone, or
  flow coherence.
- **Owning task/window:** Name the task or task range expected to satisfy the
  assertion. Use `phase-end` only when no smaller window is honest.
- **Pass / fail criterion:** Write criteria specific enough that independent
  reviewers should reach the same result. Conditional or binary assertions must
  name both qualifying and non-qualifying states; positive-only evidence is not
  sufficient.
- **Status:** Use `Draft` during planning, `Ready` before execution, and
  `PASS`, `FAIL`, `BLOCKED`, or `CARRY-FORWARD` at final review.
- **Out-of-scope notes:** Disambiguate adjacent behavior that this phase is not
  promising to prove.

## Coverage Rules

Before task execution begins:

- Every PRD/TRD/design/quality requirement in phase scope is referenced by at
  least one `VAL-NNN`, or the phase brief records why no assertion applies.
- Every `VAL-NNN` is referenced by at least one task spec or explicitly marked
  `phase-end`.
- Every Developer task spec has a `Validation Contract Coverage` section that
  lists the `VAL-NNN` IDs it advances or states why none apply.
- Every `VAL-NNN` has a planned evidence method, one validator role, and an
  owning task/window.
- QA plan steps map every `VAL-NNN` to deterministic `(qa-lead)`,
  mid-phase-review, or `(human)` verification.
- Human `VAL-NNN` assertions are one observable where practical. If one must
  span several observables, the row notes that QA Lead must fan it out into one
  human-checklist row per observable, using sub-IDs that cite the parent
  `VAL-NNN`.
- Parity or "matches reference" assertions include a structural conformance or
  parity proof whenever the reference can be checked deterministically.

Gaps are QA-plan defects and block implementation until resolved.

## Outcomes

Phase-end final review records every `VAL-NNN` outcome:

- `PASS` - assertion satisfied with cited evidence.
- `FAIL` - assertion not satisfied; blocking unless Tech Lead records an
  accepted carry-forward or scope decision.
- `BLOCKED` - assertion could not be evaluated because prerequisite evidence,
  environment, or configuration was missing.
- `CARRY-FORWARD` - assertion intentionally deferred with Tech Lead rationale
  and a receiving phase or backlog target.

`PASS` without evidence is invalid. If evidence cannot be shown, use `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| <YYYY-MM-DD> | VAL-NNN | <created/refined/superseded> | <why> | <DEC-NNN / ESC-NNN / N/A> |
