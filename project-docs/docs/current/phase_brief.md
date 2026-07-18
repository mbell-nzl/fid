# Phase Brief - No Active Phase (v0.1 closed; roadmap re-sliced)

> **Purpose.** Default execution brief for the current phase. There is no active
> phase; Phase 17 wrapped on 2026-07-18 and v0.1 is closed for early dogfood.
>
> **Size target:** under 150 lines.

## Status

**No phase active. No implementation has begun.** Phase 17 (Maintenance + stub
guard + v0.1 closure) is COMPLETE / WRAPPED. QA final review PASS, Human
checklist PASS, VAL-001..025 resolved. **v0.1 is closed for early dogfood
(DEC-057).**

On 2026-07-18 the Tech Lead performed the **trusted-evidence planning
checkpoint** (`../reviews/2026-07-18-fid-roadmap-update-tl-brief.md`),
recording **DEC-059..DEC-066** in
`../reviews/2026-07-18-roadmap-update-decisions.md` and re-slicing the
roadmap (specs 00-10 updated). Headlines, all **pending Human ratification**:

- fid executes registered gates and issues trusted receipts; reported exit
  codes can never satisfy trusted obligations (DEC-059/060).
- Typed evidence records with trust source and freshness; conservative
  whole-tree staleness first (DEC-061).
- Completion prose optional behind structured `has_completion`; completion
  views and the human validation brief are generated (DEC-062/064).
- Work envelopes + deterministic runner primitives; runner stays a separate
  kit component (DEC-063). Extended trajectory economics (DEC-065).
- Roadmap re-slice (DEC-066): v0.2a workflow engine (+ seams) → **v0.2b
  evidence foundation** → v0.2c lean completion → v0.2d governance → v0.2e
  reports/run-ops; work envelope leads v0.3.

## Proposed Next Phases (pending Human approval)

1. **CF-023 workflow-tooling pass** (pre-phase, light mode, doc/tooling only):
   completion-report template affordance, Developer-prompt commit rule,
   report-lint messaging, ESC-NNN uniqueness (DEC-058). Recommended before
   any v0.2 execution.
2. **Phase 18 - v0.2a Workflow Engine Migration (heavy mode):** workflow
   tables/migration seeded verbatim from v0.1 defaults, `fid workflow
   define/assign/list`, escalation workflow, `role_permissions` migration,
   plus the gate-obligation seams from DEC-066. Heavy mode because it
   migrates hard-to-change contracts (schema + workflow semantics). It is
   the correct first dependency: every later evidence obligation hangs off
   workflow rules referencing stable gate definitions.

Do not begin v0.2 work until Phase 18 is planned (validation contract, QA
plan, task slices) and Human-approved.

## Open Carry-Forward Into v0.2

- **CF-023 (Medium, deep-retrospective flagged):** workflow-tooling hardening
  (DEC-058). Recommended before v0.2a execution.
- **CF-020 (v0.2):** future-`occurred_at` under burst writes vs the
  future-excluding `recent_activity_7d` window (DEC-056).
- **CF-021 / CF-022 (Low):** Go doc comments on maintenance exports; broader
  markdown characterization tests. Do at next relevant package reopen.
- **CF-008 (v0.2), CF-017 (now v0.2e - dev resume/abandon relabeled by
  DEC-066), CF-018 (next SQL package reopen)** unchanged in substance.

## Deferred Release Maturity

Public-release threat model, SAST/scans, pre-release checklist, and held-out
scenarios remain deferred to v0.3 per DEC-049 and `08_QUALITY_PLAN.md`
sections 5.8/7/9. The v0.2b registered-execution threat surface adds
TRD-SEC-005..009 and quality-plan §5.10 scenarios to that eventual pass.
