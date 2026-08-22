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
- **Decision:** Schedule CF-020 and CF-027 in task 8: Event v2
  sequence replaces synthetic future occurrence time as the ordering mechanism,
  and reopened `internal/issues/create.go` is decomposed within the recorded
  size targets without behavior change. Schedule CF-032 in task 12 before
  `BuildIssueRevisionEffects` changes. Schedule CF-024's 100k-Event rebuild
  target in task 21. Re-defer CF-008, CF-018, CF-026a, CF-026b, and
  CF-028 because their product/file triggers do not fire. Re-defer CF-023,
  CF-029, and CF-030 as temporary V1 delivery-tool improvements rather than
  mixing them into product tasks. CF-017/022/031 remain cancelled/closed.
- **Consequences:** Tasks 8, 12, and 21 carry explicit tests/evidence. No
  unscheduled triggered debt is hidden in the phase brief, and no workflow
  tooling task dilutes the one-Issue product slice.
- **Refs:** `decisions_active.md`; DEC-152; task 8; task 12; task 21.

### DEC-156: Gate evidence guards the terminal-success Issue transition

- **Date:** 2026-08-22
- **Context:** The accepted DEC-153 walking skeleton evaluated gate evidence,
  but its journey ordering allowed an execution-bound Issue transition before
  Receipt linking and did not distinguish a nonterminal transition from the
  terminal-success transition. The Human supplied protected direction that
  the one gate-evidence Obligation guards terminal success.
- **Decision:** For an Envelope-bound Issue, the terminal-success transition
  is permitted only after fid executes the registered Gate, records and links
  its qualifying Receipt, and the gate-evidence Obligation evaluates
  `satisfied`. Envelope completion is a later operation that re-evaluates the
  same Obligation and also requires the selected Issue to be terminal-success.
  Both operations call one shared Obligation evaluator and return the same
  deterministic Obligation reason values; their separate action predicates
  remain outside that evaluator. A nonterminal execution-bound transition may occur earlier
  under the accepted Actor/Envelope/Claim/Run predicates. DEC-153 remains an
  accurate historical decision and is amended, not rewritten.
- **Consequences:** Task 12 builds the shared fail-closed evaluator boundary
  and the nonterminal transition. Task 17 supplies trusted Evidence inputs.
  Task 18 owns the positive terminal transition and completion reuse. The
  real-process journey mutates once, runs the Gate, links/evaluates Evidence,
  transitions terminal-success, completes, then rebuilds and repeats reads
  and negative checks without creating a second Envelope revision.
- **Refs:** Human Phase 21 planning correction 2026-08-22; PRD-WF-004;
  PRD-ENV-008; `gate-d3-2-proposal.md`; `contracts.md` §§6, 10, 16.

### DEC-157: Complete Phase 21 interface and convergence sequencing

- **Date:** 2026-08-22
- **Context:** The first planning pack left operation inputs/results/errors,
  live frame payloads, bootstrap convergence, V1 caller dispositions, and
  several multi-observable validation rows incomplete. Its checkpoint also
  required final V1 absence before the removal task.
- **Decision:** Adopt `public-interface.md` as the complete final-binary wire
  contract and `v1-convergence-inventory.md` as the live caller/reader/writer
  disposition. Closed public enum members use lowercase kebab-case; opaque
  IDs remain identifiers. `project.initialized` is Event v2 with service Actor,
  aggregate sequence 1, null repository/idempotency, and zero domain Effects;
  it never writes `roles`. Migration 0014 and V1 removal are task 20, before
  task 21 final-schema maintenance/init/rebuild/export parity. The Group 1
  checkpoint proves isolation from explicitly listed temporary V1 state, not
  its premature absence. Atomic validation assertions replace bundled rows.
- **Consequences:** `contracts.md`, canonical specs, DDL checks, task specs,
  phase brief, and QA handoff all cite the two new authorities. Task 1 remains
  unauthorized until Human approval and a Ready QA plan are accepted. No
  runtime dependency or OS-containment claim is added.
- **Refs:** Human Phase 21 planning correction 2026-08-22;
  `public-interface.md`; `v1-convergence-inventory.md`;
  `qa/validation-contract.md`; `qa/group-1-checkpoint-contract.md`.

### DEC-158: Resolve QA-plan refinements without weakening the rebuild target

- **Date:** 2026-08-23
- **Context:** The QA Lead committed the Heavy plan at `f5ccbb6` with complete
  VAL mapping and feasibility, but found two blockers: every task specified an
  invalid/self-referential one-argument report-lint command, and the retained
  100k rebuild control produced one 33.46s failure followed by 27.11s and
  24.97s passes on the only verification machine. It also proposed bounded
  evidence refinements R-04..R-19.
- **Decision:** Remove report lint from Developer command gates. The Phase
  Runner runs the exact two-argument control against the committed report and
  task spec before the next dispatch. Keep the TRD `<30 s` target and the
  benchmark's per-sample `b.Fatalf`; do not use best/median tolerance. The
  Phase 21 reference instance is Linux/amd64 on an AMD Ryzen 7 5700G with six
  exposed logical CPUs and 8,133,236 kB recorded memory, using Go 1.26.2;
  kernel, toolchain, CPU count/model, and memory are captured with the gate.
  VAL-083 requires three consecutive fresh `-benchtime 1x` processes, each
  self-asserting below 30 seconds with parity. Any one failure fails the
  assertion and task 21 owns causal optimization, not threshold relaxation.
- **Decision:** Accept R-04..R-08 and R-10..R-18. Human VAL-036 remains one
  parent perception assertion and fans out to HUMAN-01a..01d as required by
  workflow section 6.1. R-19 remains deferred until Human execution approval.
  Correct R-09's status premise: at final task-20 HEAD the role-policy and
  `--by` recipes must terminate `exit_code: 1`, while the broad retained-reader
  and Event recipes terminate 0; every actual match still must classify and
  the completion report records `UNCLASSIFIED=0`.
- **Consequences:** Task 1 remains unauthorized. QA Lead must re-review the
  applied refinements, correct its task-20 recipe expectations, and change the
  QA plan to `Ready` before Human execution approval. The 30-second target,
  Event schema, product surface, phase shape, and task order do not change.
- **Refs:** QA plan BC-01/BC-02 and R-01..R-19 at `f5ccbb6`; DEC-071;
  CF-024; TRD section 9; `qa/validation-contract.md`; tasks 1-22.

## Spec Sync Verdict

Planning-time changes are recorded with their owning decisions. Full positive
verdicts are completed at phase wrap-up.

| Spec | Planning verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending wrap audit | - |
| `01_PRD.md` | touched: see DEC-156 | Phase 21 qualification makes the terminal-success Issue transition and later Envelope completion independently evaluate the same gate-evidence Obligation. |
| `02_TRD.md` | touched: see DEC-154, DEC-156, DEC-157, DEC-158 | Capability availability, shared evaluator, bootstrap, public enums, trust, idempotency, executor boundary, and the unchanged deterministic rebuild target are explicit. |
| `03_TRD_SCHEMA.md` | touched: see DEC-157 | Exact migrations, kebab checks, displaced rule deletion, and role-free bootstrap fixture are accepted for review. |
| `04_TRD_CONFIG.md` | touched: see DEC-157 | Exact public executor bounds agree with the closed input contract. |
| `05_DESIGN_PLAN.md` | touched: see DEC-157 | Exact live NDJSON payload union and Human rendering authority is the public-interface contract. |
| `06_DESIGN_PLAN_CLI.md` | touched: see DEC-157 | Complete leaves, flags, inputs, results, errors, exits, and enum spelling are explicit. |
| `07_DESIGN_PLAN_FLOWS.md` | touched: see DEC-156, DEC-157 | Gate/evidence precede terminal transition; rebuild repeats reads and checks without a second Envelope. |
| `08_QUALITY_PLAN.md` | touched: see DEC-152, DEC-157, DEC-158 | Checkpoint isolation, removal-before-parity, cadence, feasibility, and committed-report lint ownership are explicit. |
| `09_AI_WORKFLOW.md` | verified-no-planning-change | DEC-152 applies existing Heavy-mode and recovery rules without changing them. |
| `10_IMPLEMENTATION_ROADMAP.md` | touched: see DEC-153, DEC-156, DEC-157 | Gate closed; one Heavy Phase 21/two groups; terminal guard and tasks 19-22 ordering fixed. |
| `GLOSSARY.md` | touched: see DEC-153, DEC-156, DEC-157 | Routine, capability availability, public enum, attribution, and terminal-success are unambiguous. |

## Triage-Log Audit

Phase not started; no Phase 21 TRG exists.

## Workflow Notes

To be completed at phase wrap-up.

### Deep Retrospective Recommendation

To be completed at phase wrap-up.
