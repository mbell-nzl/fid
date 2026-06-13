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

### ESC-003: Task 1 VAL-005 rebuild recovery evidence gap

- **Raised:** 2026-06-13
- **By:** Tech Lead Lite triage (`TRG-001`, Phase 5)
- **Trigger:** QA blocking triage found an uncovered acceptance criterion / validation-evidence gap.
- **Context:** Phase 5 Task 1 completion claims VAL-005 PASS for rejected-create journal hygiene. Review-1 found the cited evidence proves zero journal append, later `Get`, dry-run rollback, and unchanged `projection_meta`, but does not prove rebuild succeeds after the rejected-create repro.
- **Question / proposed direction:** Require Task 1 follow-up evidence that the rejected-create fixture can rebuild successfully after the failed write, and explicitly cite accepted-create one-line append evidence before VAL-005 is treated as on track.
- **TL resolution (2026-06-13):** Confirmed evidence gap, not an implementation defect — the seam is correct (reviewer re-ran focused tests PASS). Directed a do-now, test-and-report-only remediation: `project-docs/docs/phases/phase-5/tasks/task-1b.md`. Cannot defer (VAL-005 is a hard contract assertion and Tasks 3/5 build on the seam) and must not weaken. After Task 1b closes, Phase Runner re-runs Review 1 (or a review-1 addendum) to clear the BLOCK before Task 2.
- **Links:** `project-docs/docs/phases/phase-5/qa/review-1.md`; `project-docs/docs/phases/phase-5/tasks/task-1.md`; `project-docs/docs/phases/phase-5/tasks/task-1-done.md`; `project-docs/docs/phases/phase-5/tasks/task-1b.md`; `project-docs/docs/phases/phase-5/qa/validation-contract.md`
- **Status:** Open — resolution directed; closes when Task 1b lands and Review 1 re-passes.

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
