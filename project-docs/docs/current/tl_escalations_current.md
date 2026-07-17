# TL Escalations - Current

> **Purpose.** Open items that exceed local-decision authority and need full
> Tech Lead or human input.
>
> **Scope.** Open items only. Resolved phase items move to the phase archive.
>
> **Size target:** under 80 lines.

## Open Escalations

None.

## Recent Archive Notes

- Phase 17 raised five escalations, all resolved and archived 2026-07-17 to
  `../phases/phase-17/tl-escalations.md`: ESC-004 (Task 1 report-lint), ESC-005
  (Task 2 report-lint, backfilled), ESC-009 (QA model qualification), ESC-010
  (rebuild parity; resolved by DEC-056 Branch A), and ESC-011 (Task 7
  report-lint, backfilled with a new ID). That file also records the WRAP-F1/F2
  register reconciliation and two open retrospective items: recurring
  report-lint commit-metadata friction, and **ESC-NNN IDs are not globally
  unique** despite the template requiring it (ESC-001/004/005 each collide
  across phases).
- Phase 16 raised no `ESC-NNN` items. Local unblock and workflow decisions
  DEC-043..048 plus TRG-046 are archived in
  `../phases/phase-16/tl-escalations.md` and `../phases/phase-16/decisions.md`.
- Phase 16 VAL-013 is not an escalation; the completed human checklist records
  PASS for HUMAN-16a..d.

## Standing Escalation Triggers

TL-Lite escalates instead of deciding locally when an issue:

1. Changes a public interface, CLI contract, output shape, journal event schema,
   data model, persistence contract, dependency, tool, quality gate, or workflow.
2. Contradicts `decisions_active.md`, `decisions_principles.md`, `AGENTS.md`, or
   a recorded phase decision.
3. Affects more than one task/subsystem, exposes product-intent ambiguity, or
   changes testing strategy / coverage policy.
4. Resumes context-exhausted or partial work spanning schema/data model,
   dependencies, public interfaces, canonical docs, or multiple Developer
   sessions.

QA findings escalate here if they have architecture implications, ambiguous
acceptance criteria, unclear security/privacy severity, contradict a standing
decision, span tasks/subsystems beyond the review window, plausibly indicate a
material spec deviation, or expose a quality-gate anomaly not anticipated by the
QA plan.
