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

None.

Phase 8 (Task Completions) raised **no escalations**: both mid-phase reviews
returned CONTINUE, with zero Developer BLOCKs and zero TL-Lite triage events. See
`project-docs/docs/phases/phase-8/tl-escalations.md`.

ESC-005 (Phase 7 AC verification latest-ordering nondeterminism) was resolved by
DEC-021 (human-confirmed 2026-06-14); Task 2b implemented `rowid`-DESC
append-order latest-per-role and the Review 1 fixup cleared VAL-004. Archived to
`project-docs/docs/phases/phase-7/tl-escalations.md` at the Phase 7 wrap.

ESC-004 (Phase 6 `fid ready` project-wide tie-break order) was resolved by
DEC-019 (human-confirmed 2026-06-14) and archived to
`project-docs/docs/phases/phase-6/tl-escalations.md` at the Phase 6 wrap.

ESC-003 (Phase 5 Task 1 VAL-005 rebuild-recovery evidence gap) was
resolved by Task 1b + `review-1-addendum.md` and archived to
`project-docs/docs/phases/phase-5/tl-escalations.md` at the Phase 5 wrap.

---

The Phase 4 journal-poisoning defect (rejected `fid issue create` appends to the
journal and poisons the tail) is tracked as carry-forward **CF-007**
(`decisions_active.md`) / **DEC-014** (`phases/phase-4/decisions.md`), not as an
open escalation — it is a well-diagnosed, scheduled Phase 5 fixup.

ESC-002 (Phase 4 due-date read-back gap) was resolved by Task 4 and archived to
`project-docs/docs/phases/phase-4/tl-escalations.md`.

ESC-001 (SQLite pooled-connection pragmas) was resolved by DEC-006 / Phase 1
Task 7 and archived to `project-docs/docs/phases/phase-1/tl-escalations.md`.

<!--
Add new escalations above this comment using the template at
project-docs/templates/tl-escalations-template.md. Newest at top.
Use global ESC-NNN sequential numbering that never resets.
-->
