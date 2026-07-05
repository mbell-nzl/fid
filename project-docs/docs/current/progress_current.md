# Progress - Current

> **Purpose.** In-flight phase status only. Historical phase progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 16 is **WRAPPED**; no implementation task is active. Phase 17
kickoff is next.

Phase 16 outcome:

- Deterministic final review returned PASS WITH IMPROVEMENTS.
- VAL-001..013 PASS after human checklist completion.
- CF-011 and CF-016 closed.
- No `ESC-NNN` items are open; TRG-046 RESUME was audited correct at wrap.

## Carry-Forward

- CF-008 (v0.2): resolution-reason field/flag.
- CF-017 (v0.2c): behavior-preserving `internal/dev` `Finish`/`Block`
  function-size extraction.
- CF-018: split `internal/sqlquery/sqlquery.go` when SQL package work next
  reopens.
- CF-019: split `internal/format/markdown.go` when formatter-heavy work next
  reopens.

## Next Action

Tech Lead should kick off Phase 17 from the roadmap entry "Maintenance + stub
guard + v0.1 closure."

Do not resume Phase 16 Developer work. Use the archive for Phase 16 details:

- `project-docs/docs/phases/phase-16/decisions.md`
- `project-docs/docs/phases/phase-16/progress.md`
- `project-docs/docs/phases/phase-16/qa/final-review.md`
- `project-docs/docs/phases/phase-16/qa/human-checklist.md`
