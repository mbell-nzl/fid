# Progress - Current

> **Purpose.** In-flight phase status only. Historical progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 18 **COMPLETE / WRAPPED (2026-07-25)**. **No phase is active.**
v0.1 remains closed for early dogfood (DEC-057); v0.2a slice 1 of 3 is delivered.
Next: the **required CF-023 workflow-tooling pass** (DEC-074) plus the
recommended deep retrospective, then Phase 19 planning.

Phase 18 outcome:

- All five tasks closed with verified commits and a terminating, green
  seven-command heavy-mode gate; coverage total `76.8%`, all enforced package
  thresholds met.
- QA Lead final review (rerun 2026-07-23): **PASS** — 24 PASS / 0 FAIL / 0
  BLOCKED / 0 CARRY-FORWARD, no Blocking or Major finding.
- Human checklist (2026-07-25): **PASS** (HUMAN-01a..d all PASS).
- VAL-001..024 all PASS against fresh evidence. No assertion, threshold,
  fixture, or gate command was weakened at any point in the phase.
- Six triages (TRG-001..006), all `BLOCK`, all correct in hindsight; one
  escalation (ESC-001), resolved by DEC-072 and archived.
- Wrap decisions: **DEC-073** (sign-off, audits, CF-024, glossary sync) and
  **DEC-074** (CF-023 required before Phase 19; roadmap sequencing corrected).

Delivered: workflow structure is now data (migrations `0007`/`0008`), every
consumer resolves through `workflow.Load` from the projected workflow tables, the
frozen Go `definitions` map is test-only, and the DEC-066 gate-obligation seam is
present with zero behavior.

## Carry-Forward

- **CF-024 NEW (DEC-073)** — 100k-event rebuild sits at `26.99–27.66 s` against
  the `<30 s` TRD-PERF-003 control (~8–10% headroom). Reopen with the first
  v0.2b projection/rebuild work; next lever is batching commits per N events.
  Do not weaken the target.
- **CF-023 (Medium, promoted by DEC-074)** — workflow-tooling hardening is now
  **required before Phase 19 execution**, scope expanded from Phase 18 evidence.
- **CF-021 (re-targeted)** — trigger fired in Phase 18 and was missed; now an
  explicit Phase 19 planning item.
- **CF-017 (v0.2e, re-targeted)** — trigger fired in Phase 18; correctly not
  taken during a behavior-preserving migration; now scheduled, not trigger-based.
- **CF-008 (v0.2)** — explicit decision point at Phase 19 planning.
- CF-018, CF-020, CF-022 unchanged; their trigger packages were untouched.

## Escalations

**None open.** Phase 18's ESC-001 (final-review BLOCK — `workflow.Load`
connection deadlock plus the then-absent QA-30 spec sync) is resolved by DEC-072,
verified by the 2026-07-23 rerun, and archived in
`../phases/phase-18/tl-escalations.md`. All Phase 17 escalations remain archived
in `../phases/phase-17/tl-escalations.md`.

## Planning Checkpoint (2026-07-18)

The trusted-evidence planning checkpoint ran per
`../reviews/2026-07-18-fid-roadmap-update-tl-brief.md`. DEC-059..DEC-066
recorded in `../reviews/2026-07-18-roadmap-update-decisions.md`; specs 00-10 and
the roadmap re-sliced (evidence foundation v0.2b before governance/reports; lean
completion v0.2c; work envelope v0.3a). **Pending Human ratification**; Phase 18
executed within the ratified-by-practice v0.2a slice.

## Next Action

Tech Lead planning, in order (each needs Human approval before execution):

1. **CF-023 workflow-tooling pass (required, DEC-074)** — light mode, doc/tooling
   only: mint `ESC-NNN` + triage-log entry on every `BLOCK`; `ESC-NNN`/`TRG-NNN`
   global uniqueness; report-lint/commit-metadata friction; run-ledger
   chronological append without duplicates; wrap-time reopen audit for
   trigger-based carry-forwards; evidence-honesty hardening so a non-terminating
   command cannot be recorded `PASS`.
2. **Deep retrospective (RECOMMENDED, DEC-073)** — CF-023 recurrence plus the
   new evidence-over-claiming theme; pair it with the tooling pass.
3. **Phase 19 planning** — `fid workflow define/assign/list` + escalation default
   workflow (v0.2a slice 2 of 3), including CF-021 and CF-008 disposition.
