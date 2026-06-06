# Phase 1 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow
> notes, and wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-006: Enforce SQLite pragmas per physical connection via DSN

- **Date:** 2026-06-07
- **Context:** QA Review 1 / ESC-001 / TRG-003 (BLOCK). `internal/db.Open`
  applies `PRAGMA foreign_keys = ON` and `PRAGMA recursive_triggers = ON` by
  `Exec` on the pooled `*sql.DB` (`internal/db/db.go:23-28`). SQLite pragma
  state is per physical connection, and `database/sql` opens additional
  physical connections on demand (default `MaxOpenConns = 0`, unlimited).
  `foreign_keys` defaults OFF on each new connection, so VAL-014's "both
  pragmas ON on every connection path" is not met. `TestOpenEnablesRequiredPragmas`
  only inspects one connection and almost always reuses the same physical
  connection, so it passed while the guarantee was false. The DB helper is
  consumed by the journal writer, replay, init, and issue tasks, so the gap
  spans subsystems.
- **Decision:** Enforce both pragmas through the connection string using
  `modernc.org/sqlite`'s `_pragma` query parameters, e.g.
  `file:<abs-path>?_pragma=foreign_keys(1)&_pragma=recursive_triggers(1)`,
  so the driver applies them on every physical connection it opens. Keep the
  post-open `Ping`. Add a regression test that forces multiple distinct
  physical connections (via `SetMaxOpenConns(N)` and N concurrently held
  `*sql.Conn`) and asserts both pragmas report ON on each. Implemented as
  remediation Task 7. Rejected alternatives: pinning the pool to a single
  connection (`SetMaxOpenConns(1)`) satisfies VAL-014 only by degenerate
  accident and serializes all DB access; a global `RegisterConnectionHook`
  is heavier than DSN params and was kept only as a fallback.
- **Consequences:** `db.Open` keeps its public signature; this is an internal
  correctness change plus a stronger test. ESC-001 is dispositioned **do-now**
  in Phase 1 and closes when Task 7 lands with passing multi-connection
  evidence. No schema, data model, or other public-interface change.
- **Refs:** ESC-001 (`tl_escalations_current.md`); TRG-003 (`triage-log.md`);
  QA Review 1 Task 1 Blocking (`qa/review-1.md`); VAL-014
  (`qa/validation-contract.md`); QA-01 (`qa/qa-plan.md`); task-7
  (`tasks/task-7.md`); `internal/db/db.go`; TRD-ARCH-007; TRD §4.4.

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is
a positive assertion after review, not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending wrap-up | - |
| `01_PRD.md` | pending wrap-up | - |
| `02_TRD.md` | pending wrap-up | - |
| `03_TRD_SCHEMA.md` | pending wrap-up | - |
| `04_TRD_CONFIG.md` | pending wrap-up | - |
| `05_DESIGN_PLAN.md` | pending wrap-up | - |
| `06_DESIGN_PLAN_CLI.md` | pending wrap-up | - |
| `07_DESIGN_PLAN_FLOWS.md` | pending wrap-up | - |
| `08_QUALITY_PLAN.md` | pending wrap-up | - |
| `09_AI_WORKFLOW.md` | pending wrap-up | - |
| `10_IMPLEMENTATION_ROADMAP.md` | pending wrap-up | - |

## Triage-Log Audit

- [ ] Walked every TRG-NNN in `triage-log.md`.
- [ ] Confirmed each RESUME verdict was correct in hindsight.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed autonomous TL-Lite resume remains in effect (DEC-002; no warm-up tally).

## Workflow Notes

Answer in 1-3 sentences each. Completed at phase wrap-up.

- **What worked well?** <pending wrap-up>
- **What was inefficient?** <pending wrap-up>
- **Were tasks the right size?** <pending wrap-up>
- **Did the QA review cadence feel right?** <pending wrap-up>
- **Were there token-efficiency opportunities missed?** <pending wrap-up>

### Deep Retrospective Recommendation

<pending wrap-up>
