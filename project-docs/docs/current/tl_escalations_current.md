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

### ESC-005: AC verification latest ordering can be nondeterministic within one second

- **Raised:** 2026-06-14
- **By:** QA Reviewer / Phase 7 Review 1
- **Trigger:** QA Reviewer trigger 5 and 6 — spans later `ac list` / workflow consumers and is a plausible material spec deviation.
- **Context:** Task 2 stores `ac_verifications.verified_at` with `time.RFC3339` second precision while `latestVerificationsByRole` selects latest by `verified_at DESC, id DESC`. Two same-role verifications inside the same second can therefore select by random `ACV-` id rather than append order, so fail→pass may not reliably produce latest=pass.
- **Question / proposed direction:** Decide the canonical ordering for same-timestamp AC verifications and require implementation/tests that prove append-only latest-per-role semantics without sleeping past a timestamp boundary.
- **Resolution direction (DEC-021):** Latest-per-role is the verification appended last in journal order (per ADR-0001), resolved by SQLite implicit `rowid` DESC — deterministic and rebuild-stable, independent of `verified_at` second resolution and the random `ACV-` id. `verified_at` stays display-only. No schema migration. Implementation + de-sleep tests scheduled as **Task 2b**; VAL-004 strengthened to forbid timestamp-boundary proof.
- **Links:** `project-docs/docs/phases/phase-7/qa/review-1.md`; `project-docs/docs/phases/phase-7/decisions.md` (DEC-021); `project-docs/docs/phases/phase-7/tasks/task-2b.md`; `internal/issues/ac.go`; `internal/issues/detail.go`; VAL-004, VAL-006, VAL-008.
- **Status:** Decided (DEC-021) — awaiting human confirmation of the persistence-contract semantics; closes when Task 2b lands and re-review clears VAL-004.

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
