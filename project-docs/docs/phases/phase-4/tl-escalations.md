# TL Escalations - Phase 4 (Archived)

> Resolved escalations from Phase 4, moved here at phase wrap-up. Open items
> live in `project-docs/docs/current/tl_escalations_current.md`.

## ESC-002: Phase 4 due-date create field lacked a read-back surface

- **Raised:** 2026-06-07
- **By:** qa-reviewer / Phase 4 mid-phase review 3
- **Trigger:** QA Reviewer trigger 5 and 6 - spans Tasks 2-3 and is a plausible material spec deviation.
- **Context:** Task 2 and Task 3 completion reports claimed VAL-003 PASS for due-date persistence, but both proved `due_date` only by direct SQLite query. `internal/issues.Issue` had no `due_date` field, `internal/issues.Get` did not select it from `issues_enriched`, and `fid issue get --format json` therefore could not read it back through the public read path.
- **TL disposition (2026-06-09):** DO-NOW, do not narrow VAL-003. The contract names `due_date` from PRD-ISSUE-015/016/019/020 and the read-back clause stands. Fix is additive and folded into Task 4 (centralised read model): add `DueDate` to `internal/issues.Issue`, select `due_date` in `Get`, surface it in `fid issue get --format json`, and assert read-back in `internal/issues` and `cmd/fid` tests. No schema/journal change (`issues_enriched` projects `i.*`; `issues.due_date` exists). Tasks 2-3 stand; their VAL-003 claim is persistence only, with read-back closed by Task 4.
- **Resolution (2026-06-13, phase wrap):** RESOLVED. Task 4 evidence proves `due_date` reads back through `internal/issues.Get` (`2026-06-30` populated, `nil` omitted) and `fid issue get --format json` (`"due_date":"2026-06-30"` populated, `null` omitted). VAL-003 PASS at final review; no carry-forward remains. See `qa/final-review.md` ESC-002 Disposition and `decisions.md` DEC-014 context.
- **Status:** Resolved.

---

ESC-001 (SQLite pooled-connection pragmas) was resolved by DEC-006 / Phase 1
Task 7 and is archived under `project-docs/docs/phases/phase-1/tl-escalations.md`.
