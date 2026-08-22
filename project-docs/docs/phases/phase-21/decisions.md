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

### DEC-154: Fix the exact Heavy Phase 21 slice contract

- **Date:** 2026-08-22
- **Context:** DEC-153 requires exact capability, authority, JSON, DDL/Event,
  executor, evidence, idempotency, removal, and deferral boundaries before a
  Developer task may infer implementation behavior.
- **Decision:** Adopt `contracts.md` and `03_TRD_SCHEMA.md` section 8.2 as the
  planning contract presented for Human execution approval. Capability
  availability is binary and conservative; activated profiles are the only
  kit Actor/Gate trust source; Gate Registration precedes Envelope creation;
  Event schema v2 and operation-scoped durable idempotency apply to converted
  mutations; executor policy is limited to portable standard-library controls;
  and evaluation returns every independent reason in fixed order. Use
  migrations 0013/0014 so generic Claims can pass before V1 `claims` is
  removed. Freeze `task_completions` without writers until Phase 23 and
  `agent_runs` until Phase 25 because existing foreign-key/read ownership makes
  immediate deletion a cross-phase change; neither is authoritative Phase 21
  completion or Run state.
- **Consequences:** One Heavy phase uses two internal task groups and an
  independent checkpoint. Group 2 cannot begin until Group 1 contract closure
  and executor/verifier feasibility independently pass. The slice exposes
  export create/validate but keeps `export.portable.v1` unavailable. No
  unsupported containment, Attention, Stop, recovery, respawn, budget
  consumption, import, review, or Human authority behavior may be inferred.
- **Refs:** DEC-152; DEC-153; `contracts.md`; `03_TRD_SCHEMA.md` §8.2;
  `02_TRD.md` §3.1; ADR-0011; ADR-0012.

### DEC-155: Dispose active carry-forwards against the Phase 21 task pack

- **Date:** 2026-08-22
- **Context:** Phase planning must schedule every carry-forward whose trigger
  the slice fires and explicitly re-defer the rest.
- **Decision:** Schedule CF-020 and CF-027 in task 8 as VAL-062/063: Event v2
  sequence replaces synthetic future occurrence time as the ordering mechanism,
  and reopened `internal/issues/create.go` is decomposed within the recorded
  size targets without behavior change. Schedule CF-032 in task 12 before
  `BuildIssueRevisionEffects` changes. Schedule CF-024's 100k-Event rebuild
  target in task 20 as VAL-064. Re-defer CF-008, CF-018, CF-026a, CF-026b, and
  CF-028 because their product/file triggers do not fire. Re-defer CF-023,
  CF-029, and CF-030 as temporary V1 delivery-tool improvements rather than
  mixing them into product tasks. CF-017/022/031 remain cancelled/closed.
- **Consequences:** Tasks 8, 12, and 20 carry explicit tests/evidence. No
  unscheduled triggered debt is hidden in the phase brief, and no workflow
  tooling task dilutes the one-Issue product slice.
- **Refs:** `decisions_active.md`; DEC-152; task 8; task 12; task 20;
  VAL-062..064.

## Spec Sync Verdict

Planning-time changes are recorded with their owning decisions. Full positive
verdicts are completed at phase wrap-up.

| Spec | Planning verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending wrap audit | - |
| `01_PRD.md` | verified-no-planning-change | Existing requirements remain target baseline requirements; DEC-153 changes slice ownership/availability, not requirement semantics. |
| `02_TRD.md` | touched: see DEC-154 | Binary capability availability, Activation trust, idempotency, and the portable executor boundary are explicit. |
| `03_TRD_SCHEMA.md` | touched: see DEC-154 | Exact 0013/0014 DDL/down behavior and Event-v2 fixture authority are accepted for review. |
| `04_TRD_CONFIG.md` | touched: see DEC-154 | Timeout ceiling, closed environment, redaction, safe content-addressed Artifact storage, and availability are aligned. |
| `05_DESIGN_PLAN.md` | touched: see DEC-154 | Exact live NDJSON/control/terminal behavior replaces the D3 placeholder. |
| `06_DESIGN_PLAN_CLI.md` | touched: see DEC-154 | Accepted leaves, `--actor`, `--group-by`, input, streaming, and deferred leaf absences are explicit. |
| `07_DESIGN_PLAN_FLOWS.md` | touched: see DEC-153 | Gate Registration now precedes Envelope creation; Phase 21 export is create/validate only. |
| `08_QUALITY_PLAN.md` | touched: see DEC-152 and DEC-154 | Three-task cadence, mandatory independent checkpoint, fresh feasibility, and task-local report lint are explicit. |
| `09_AI_WORKFLOW.md` | verified-no-planning-change | DEC-152 applies existing Heavy-mode and recovery rules without changing them. |
| `10_IMPLEMENTATION_ROADMAP.md` | touched: see DEC-153 | Gate closed; one Heavy Phase 21/two groups; Attention/import ownership corrected. |
| `GLOSSARY.md` | touched: see DEC-153 | Routine and capability availability are now unambiguous. |

## Triage-Log Audit

Phase not started; no Phase 21 TRG exists.

## Workflow Notes

To be completed at phase wrap-up.

### Deep Retrospective Recommendation

To be completed at phase wrap-up.
