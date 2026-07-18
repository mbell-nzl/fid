# Phase Brief - No Active Phase (v0.1 closed)

> **Purpose.** Default execution brief for the current phase. There is no active
> phase; Phase 17 wrapped on 2026-07-18 and v0.1 is closed for early dogfood.
>
> **Size target:** under 150 lines.

## Status

**No phase active.** Phase 17 (Maintenance + stub guard + v0.1 closure) is
COMPLETE / WRAPPED. QA final review PASS, Human checklist PASS, VAL-001..025
resolved. **v0.1 is closed for early dogfood (DEC-057).**

The next planning checkpoint slices roadmap **v0.2a - Workflow Engine Migration**
into a new approved phase. This brief is refreshed at that kickoff.

## Phase 17 Result (for reference)

- Delivered `fid maintenance check`, locked `fid maintenance vacuum`, v0.2
  governance command guards, full four-format read conformance, and the v0.1
  closure journey/rebuild/traceability evidence.
- Ten tasks (0-9) closed with verified single-commit implementations under Heavy
  Mode; reviews after Tasks 0, 2, 5, 7; final review and Human dogfood complete.
- VAL-022 narrowed under DEC-056 (journal is the sole source of derived state;
  clock-relative reads excluded from rebuild-equality); Branch A confirmed with
  no product change. VAL-023 unblocked with it.
- Full archive: `../phases/phase-17/` (tasks, qa, decisions, progress,
  run-ledger, triage-log, tl-escalations).

## Open Carry-Forward Into v0.2

- **CF-023 (Medium, deep-retrospective flagged):** workflow-tooling hardening -
  report-lint / completion-report commit-metadata friction and ESC-NNN global
  uniqueness (DEC-058). Recommended before v0.2a execution.
- **CF-020 (v0.2):** future-`occurred_at` under burst writes vs the
  future-excluding `recent_activity_7d` window (DEC-056).
- **CF-021 / CF-022 (Low):** Go doc comments on maintenance exports; broader
  markdown characterization tests. Do at next relevant package reopen.
- **CF-008 (v0.2), CF-017 (v0.2c), CF-018 (next SQL package reopen)** unchanged.

## Deferred Release Maturity

Public-release threat model, SAST/scans, pre-release checklist, and held-out
scenarios remain deferred to v0.3 per DEC-049 and `08_QUALITY_PLAN.md`
sections 5.8/7/9. Phase 17 was v0.1 early-dogfood closure, not the charter's v0.3
public release.

## Next Planning Checkpoint

Plan v0.2a (Workflow Engine Migration) per `10_IMPLEMENTATION_ROADMAP.md`:
workflow tables/migration, `fid workflow define/assign/list`, escalation
workflow, and the `role_permissions` migration. Propose the operating mode,
draft the validation contract, and slice one-session tasks. Do not begin v0.2
work until that phase is Human-approved.
