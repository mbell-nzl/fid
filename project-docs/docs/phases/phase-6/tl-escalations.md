# Phase 6 Escalations (archived)

> Resolved escalations moved here at Phase 6 wrap. Open items live in
> `../../current/tl_escalations_current.md`.

## ESC-004 — `fid ready` project-wide tie-break order (RESOLVED)

- **Raised by:** TL-Lite via TRG-001 (Task 3 completion report).
- **Trigger:** Task 3 implemented a deterministic project-wide / prioritised-epic
  tie-break that the PRD does not pin; Task 4 would make it a user-visible
  scheduling contract (TL-Lite triggers 1, 4, 7).
- **Resolution:** DEC-019 (human-confirmed 2026-06-14, option A). The bare
  `fid ready` selection order is `priority ASC` (no priority ranks last) →
  `created_at ASC` → `id ASC`; sequential epics order by `sequence` and
  contribute only their current actionable task. Documentation only — Task 3's
  implementation stood, no code change. Pinned in `06_DESIGN_PLAN_CLI.md`
  §`fid ready`. TRG-001 closed; Task 4 unblocked; VAL-006 PASS.
- **References:** `decisions.md` DEC-019; `triage-log.md` TRG-001; PRD-ISSUE-008/011/012/015.
