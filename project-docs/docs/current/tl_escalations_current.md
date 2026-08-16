# TL Escalations - Current

> **Purpose.** Open items that exceed current repair authority and need full
> Tech Lead or Human input.
>
> **Scope.** Open items only. Resolved detail lives in each phase archive.
>
> **Size target:** under 80 lines.

## Open Escalations

### ESC-20-015 — report-lint corrupts a truthful exact-command comparison

- **Raised:** 2026-08-16
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-17-done.md` and
  attached independent `project-docs/tools/report-lint.sh` output.
- **Protected trigger:** quality-gate definition / workflow authority.
- **Condition:** Task 17's declared expected-non-zero absence command and its
  committed report row are literally byte-identical (104 bytes), including the
  required `v0\.2` regex escape. `report-lint.sh` passes the command through
  `awk -v`, normalizes `\.` to `.`, then falsely reports both an undeclared
  expected-non-zero command and a missing required command. Rewording the row
  to satisfy that normalization would violate exact-command coverage.
- **Why blocked:** Task 17 has no authority to modify `report-lint.sh`, and
  changing the canonical control's literal-command comparison is a protected
  quality-gate/workflow change. The supplied product, report-finalization, and
  report-recovery commit verification establishes the committed state, but the
  attached report-lint result remains `RECOVERY_REQUIRED` (terminated
  `exit_code: 2`), so it does not support RESUME.
- **Resolution required:** Full Tech Lead must authorize and direct a
  control-level repair that preserves byte-exact command comparison and all
  report evidence, then dispatch the appropriate recovery and rerun the
  unchanged report-lint control. Do not alter Task 17's declared or recorded
  command spelling.

## Archive Pointers

- Phase 20: `project-docs/docs/phases/phase-20/tl-escalations.md`
- Phase 19: `project-docs/docs/phases/phase-19/tl-escalations.md`
- Phase 18: `project-docs/docs/phases/phase-18/tl-escalations.md`
- Phase 17: `project-docs/docs/phases/phase-17/tl-escalations.md`

## Protected Escalation Triggers

TL-Lite or the full Tech Lead allocates an ESC only when a condition requires:

1. Product, architecture, data-model, or public-interface authority.
2. Dependency, quality-gate definition, testing-policy, or workflow authority.
3. Security, privacy, journal-integrity, or possible data-loss judgment.
4. Resolution of a contradictory or missing binding specification.
5. External capability or infrastructure that remains unavailable after
   bounded recovery.
6. Human judgment explicitly required by a validation contract or checklist.
7. Resolution after two consecutive same-trigger continuation attempts make no
   material progress (`recovery-budget-exhausted`).

A red test, failed gate, incomplete in-scope implementation, report defect,
bounded QA Must-fix finding, long-running command, or session continuation is
REPAIR, not an escalation, unless a separate protected trigger above is
demonstrated.

Developer and QA roles emit blocker candidates and may cite an existing
`ESC-<phase>-<NNN>`. They never allocate an ESC or edit this register.
