# TL Escalations - Current

> **Purpose.** Open items that exceed local-decision authority and need full
> Tech Lead or human input.
>
> **Scope.** Open items only. Resolved phase items move to the phase archive.
>
> **Size target:** under 80 lines.

## Open Escalations

### ESC-004 — Task 1 completion report has uncommitted post-commit edits

- **Raised by:** TRG-004
- **Phase/task:** Phase 17 / Task 1
- **Trigger:** Quality-gate failure / report-lint evidence invalidation
- **Issue:** `project-docs/docs/phases/phase-17/tasks/task-1-done.md` passes its structural and recorded-evidence parser checks, but report-lint fails because metadata-wrap changes remain uncommitted after commit `65d2f94587310668b252ebfd4141ed9ef26c8ec0` (`refactor(format): split markdown renderer`).
- **Required resolution:** Commit the current report changes using the required task documentation commit, then rerun report-lint and attach fresh passing evidence. Do not treat the earlier code commit as verification of the modified report.
- **Tech Lead disposition (2026-07-15):** Triage upheld. TRG-004 was correct to BLOCK — an uncommitted completion report has invalid report-lint evidence by definition, and the code commit `65d2f94` does not verify the report's current modifications. Resolution is mechanical and carries no spec, interface, contract, journal-schema, or architecture impact: the report body (backtick-wrapped metadata per TRG-003, plus VAL-001 coverage and the five recorded gates) is already correct and complete. No DEC-NNN required. Approved path: create the required `docs(task)` commit for the current `task-1-done.md`, rerun report-lint against the committed file, attach the fresh passing evidence, then resume Phase 17.
- **Status:** Resolution approved — closes automatically once the `docs(task)` commit lands and report-lint passes on the committed file. Phase Runner confirms closure on resume.

## Recent Archive Notes

- ESC-009 (Phase 17 QA route model-qualification) **resolved 2026-07-14** by
  TRG-001: the committed model audit (`phase-17/qa/model-audit.md`, commit
  `2d5e536`) records `PASS - QA routing qualified` with no MISS, closing
  QA-BLK-02 and satisfying DEC-053/055. Task 0 and the Entry Review proceeded
  on that basis. Full archival to `phase-17/tl-escalations.md` happens at
  phase wrap-up.
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
