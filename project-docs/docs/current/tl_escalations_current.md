# TL Escalations - Current

> **Purpose.** Unresolved escalations that exceed local-decision authority and
> need full Tech Lead or human input. Default-loaded by every agent session.
>
> **Scope.** Open items only. Resolved escalations move to
> `docs/phases/phase-N/tl-escalations.md` at phase wrap-up.
>
> **Size target:** under 80 lines.

## Standing Escalation Triggers

### TL-Lite Triggers

An issue must be escalated, not decided locally by TL-Lite, if it:

1. Changes a public interface.
2. Changes the data model or persistence contract.
3. Overrides or contradicts a standing decision.
4. Affects more than one task or subsystem.
5. Introduces a dependency, tool, quality gate, or workflow.
6. Materially changes testing strategy or coverage policy.
7. Exposes ambiguity in product intent or acceptance criteria.
8. Resumes a context-exhausted or partial task where remaining or completed
   work spans schema/data model, dependency/platform configuration, public
   interfaces, canonical docs, or multiple Developer sessions.

### QA Reviewer Triggers

Most QA findings stay in `review-M.md`. QA Reviewer also escalates here if the
finding:

1. Has architecture implications beyond the three-task review window.
2. Implies acceptance criteria are ambiguous.
3. Has unclear security or privacy severity.
4. Contradicts a standing decision.
5. Spans tasks or subsystems beyond the review window.
6. Is a plausible but material spec deviation.
7. Is a quality-gate anomaly not anticipated by the QA plan.

## Open Escalations

None. ESC-006 (Phase 9) resolved at Phase 9 wrap.

## Resolved (archived in their phase `tl-escalations.md`)

- ESC-006 (Phase 9 final-review blockers: benchmark gate failure + concern
  address read-back gap) → DEC-023 / DEC-024 / Task 4 fixup; QA Review 2 verified,
  all eleven VALs PASS at wrap.
- Phase 8: no escalations raised (both reviews CONTINUE, zero BLOCKs/triage).
- ESC-005 (Phase 7 AC verify ordering nondeterminism) → DEC-021 / Task 2b.
- ESC-004 (Phase 6 `fid ready` tie-break order) → DEC-019.
- ESC-003 (Phase 5 VAL-005 rebuild-recovery evidence gap) → Task 1b.
- ESC-002 (Phase 4 due-date read-back gap) → Task 4.
- ESC-001 (SQLite pooled-connection pragmas) → DEC-006 / Phase 1 Task 7.
- Phase 4 journal-poisoning defect tracked as CF-007 / DEC-014 (scheduled
  fixup), not an open escalation.

<!--
Add new escalations above this comment using the template at
project-docs/templates/tl-escalations-template.md. Newest at top.
Use global ESC-NNN sequential numbering that never resets.
-->
