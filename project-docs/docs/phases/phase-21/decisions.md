# Phase 21 Decisions

> Purpose: phase-local decision history, planning decisions, spec sync
> verdicts, workflow notes, and eventual wrap-up decisions. DEC numbering is
> global and never resets.

## Decisions

### DEC-152: Phase 20 retrospective changes Phase 21 delivery controls

- **Date:** 2026-08-22
- **Context:** Phase 20 delivered a clean product slice but generated 23 triage
  identifiers, 19 escalations, multiple control repairs, and a 22-report final
  lint sweep. The Human required the recommended deep retrospective before the
  Phase 21 task pack.
- **Decision:** Adopt the seven planning decisions in
  `../../reviews/2026-08-22-phase-20-retrospective.md`: bounded causal task
  ownership, fresh feasibility for every prescribed verifier, task-local
  report-lint, three-task review windows, a mandatory independent Group 1
  checkpoint, retention of strict typed recovery, and no product-slice tooling
  task for CF-029/CF-030.
- **Consequences:** Phase 21 uses two internal groups inside one Heavy phase.
  Group 2 cannot start on implementation momentum alone; an independent PASS
  must establish Group 1 contract closure and standard-library executor
  feasibility. QA-plan creation remains a QA Lead responsibility.
- **Refs:** Phase 20 DEC-151; Phase 20 `triage-log.md` and
  `tl-escalations.md`; `09_AI_WORKFLOW.md` §§4, 6, 10.

### DEC-153: Human accepts Gate D3.2 with all required amendments

- **Date:** 2026-08-22
- **Context:** The Human approved the core one-Issue walking-skeleton proposal
  and supplied binding answers DP-1..DP-7 plus ordering, authority,
  capability, execution, removal, export, and terminology amendments.
- **Decision:** Close Gate D3.2 exactly as recorded in
  `gate-d3-2-proposal.md`. Plan one Heavy Phase 21 with two internal task
  groups. Defer Attention, Stops, budget consumption/exhaustion, recovery,
  automatic respawn, and export import to Phase 22; implement one generic
  gate-evidence Obligation; advertise no incomplete semantic capability; and
  claim only standard-library-enforceable executor controls.
- **Consequences:** ADR-0012 and the roadmap retain one numeric vertical slice
  while naming the internal checkpoint. The exact capability matrix, JSON,
  DDL/Event, trust, idempotency, output, Artifact, removal, and deferral
  contracts must be accepted before implementation. `Routine` is external kit
  policy language, never fid state.
- **Refs:** Human Gate D3.2 direction 2026-08-22; ADR-0011; ADR-0012;
  `10_IMPLEMENTATION_ROADMAP.md` §§5, 10; DEC-152.

## Spec Sync Verdict

Planning-time changes are recorded with their owning decisions. Full positive
verdicts are completed at phase wrap-up.

| Spec | Planning verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending wrap audit | - |
| `01_PRD.md` | verified-no-planning-change | Existing requirements remain target baseline requirements; DEC-153 changes slice ownership/availability, not requirement semantics. |
| `02_TRD.md` | pending slice contract sync | DEC-153/154 |
| `03_TRD_SCHEMA.md` | pending slice DDL/Event sync | DEC-154 |
| `04_TRD_CONFIG.md` | pending executor/export sync | DEC-154 |
| `05_DESIGN_PLAN.md` | pending live-output sync | DEC-154 |
| `06_DESIGN_PLAN_CLI.md` | pending slice leaf sync | DEC-154 |
| `07_DESIGN_PLAN_FLOWS.md` | touched: see DEC-153 | Gate Registration now precedes Envelope creation; Phase 21 export is create/validate only. |
| `08_QUALITY_PLAN.md` | pending evidence sync | DEC-152/154 |
| `09_AI_WORKFLOW.md` | verified-no-planning-change | DEC-152 applies existing Heavy-mode and recovery rules without changing them. |
| `10_IMPLEMENTATION_ROADMAP.md` | touched: see DEC-153 | Gate closed; one Heavy Phase 21/two groups; Attention/import ownership corrected. |
| `GLOSSARY.md` | touched: see DEC-153 | Routine and capability availability are now unambiguous. |

## Triage-Log Audit

Phase not started; no Phase 21 TRG exists.

## Workflow Notes

To be completed at phase wrap-up.

### Deep Retrospective Recommendation

To be completed at phase wrap-up.
