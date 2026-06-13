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

### ESC-002: Phase 4 due-date create field lacks read-back surface

- **Raised:** 2026-06-07
- **By:** qa-reviewer / Phase 4 mid-phase review 3
- **Trigger:** QA Reviewer trigger 5 and 6 - spans Tasks 2-3 and is a plausible material spec deviation.
- **Context:** Task 2 and Task 3 completion reports claim VAL-003 PASS for due-date persistence, but both prove `due_date` by direct SQLite query. `internal/issues.Issue` has no `due_date` field, `internal/issues.Get` does not select it from `issues_enriched`, and `fid issue get --format json` therefore cannot read it back through the public read path.
- **Question / proposed direction:** Decide whether Phase 4 read-back requires `due_date` in `internal/issues.Get` / CLI JSON now, or whether VAL-003 and the CLI acceptance criteria should be narrowed with an explicit recorded carry-forward.
- **Links:** `project-docs/docs/phases/phase-4/qa/review-3.md`; `project-docs/docs/phases/phase-4/qa/validation-contract.md` VAL-003; `project-docs/docs/phases/phase-4/tasks/task-2-completion.md`; `project-docs/docs/phases/phase-4/tasks/task-3-completion.md`; `internal/issues/issue.go`
- **TL disposition (2026-06-09):** DO-NOW, do not narrow VAL-003. The contract
  names `due_date` from PRD-ISSUE-015/016/019/020 and the read-back clause stands.
  Fix is additive and folded into Task 4 (centralised read model): add `DueDate`
  to `internal/issues.Issue`, select `due_date` in `Get`, surface it in
  `fid issue get --format json`, and assert read-back in `internal/issues` and
  `cmd/fid` tests. No schema/journal change (`issues_enriched` projects `i.*`;
  `issues.due_date` exists). Tasks 2-3 stand; their VAL-003 claim is persistence
  only, with read-back closed by Task 4.
- **Status:** Resolved-pending — closes when Task 4 lands due-date read-back
  evidence; flip to Resolved at phase wrap-up and archive.

ESC-001 (SQLite pooled-connection pragmas) was resolved by DEC-006 / Task 7 and
archived at Phase 1 wrap-up to
`project-docs/docs/phases/phase-1/tl-escalations.md`.

<!--
Add new escalations above this comment using the template at
project-docs/templates/tl-escalations-template.md. Newest at top.
Use global ESC-NNN sequential numbering that never resets.
-->
