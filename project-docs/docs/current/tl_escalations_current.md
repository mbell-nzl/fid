# TL Escalations - Current

> **Purpose.** Open items that exceed current repair authority and need full
> Tech Lead or Human input.
>
> **Scope.** Open items only. Resolved detail lives in each phase archive.
>
> **Size target:** under 80 lines.

## Open Escalations

None.

`ESC-20-016` is resolved by DEC-147 (2026-08-18) and archived in
`project-docs/docs/phases/phase-20/tl-escalations.md`. Task 3k's first required
command is now
`go test -race -count=10 -timeout 30m ./internal/claims ./cmd/fid`; the count,
race detector, package scope, and pass/fail semantics are unchanged. One
bounded report-scope Task 3k continuation is authorized to rerun the seven
commands and reach `OUTCOME COMPLETE`, after which MF-01's fresh independent QA
re-review proceeds. `ESC-20-015` is resolved by DEC-146 (2026-08-16) and
archived in the same file.

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
