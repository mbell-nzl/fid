# Phase 22 Decisions

> Tech Lead planning, 2026-09-05. Product proposals await Human acceptance.
> Global DEC numbering continues after Phase 21 DEC-179.

## Decisions

### DEC-180: Propose the exact Heavy failure-and-recovery slice

- **Date:** 2026-09-05.
- **Status:** PROPOSED; Human acceptance is pending. No execution authority.
- **Context:** Roadmap section 5 assigns the recovery journey and explicit V1
  removals to Phase 22. Phase 21 is wrapped under DEC-179 and merged as
  89d964a84dbdd297fdf472d6f78b001e2efac702; its current-context branch note was
  stale. Planning started from that clean main HEAD on new branch `phase-22`.
- **Proposal:** Heavy mode, 18 ordered tasks, 80 validation assertions, six
  review windows, an early task-6 subprocess recovery journey, and an independent
  checkpoint after task 12 before archive/import work. Accept `contracts.md`,
  `public-interface.md`, and `schema.md` together before implementation.
- **Product choices presented for approval:** trusted executor observations with
  closed terminated/live/unknown liveness; unknown always blocks; exact Human
  force grant required for live takeover; explicit adoption with no Git mutation;
  one new linked successor; event-derived counters/Stops/Attention; measured
  unknown values stay null; derived close/intentional abandon replace caller
  finish outcomes; activation restore-prior before commit and complete-forward
  after; bounded manifest-v2 import into an empty target with atomic publication.
- **Capability boundary:** attention, activation, and portable export qualify
  only after their complete slice proofs. Claims recovery is implemented but
  `work.claims.v1` stays false until Phase 24 parallel-conflict semantics pass.
  No partial availability or unsupported Human Outcome resolution is implied.
- **Consequences:** Explicit development-baseline amendment of 0013 and new
  0015/0016 migrations; no old-data migration, new runtime dependency, scheduler,
  respawn, Review subsystem, or multi-Issue contract. The roadmap permits
  horizontal schema/logic/adapter/removal tasks within this coherent scenario;
  this is the explicit phase-boundary exception to otherwise isolated changes.
- **Refs:** roadmap sections 5/7; PRD sections 6/7/9/10; TRD sections 2-9;
  schema section 8.3; ADR-0007; `qa/validation-contract.md`.

### DEC-181: Make every task's ownership executable from its first gate

- **Date:** 2026-09-05.
- **Status:** Tech Lead planning decision; applies when the pack is approved.
- **Context:** Phase 21 workflow notes identify unexecutable retained-test
  ownership and a late first end-to-end gate as repeated planning defects.
- **Decision:** Every task requires the ordinary full task gate plus the exact
  tagged integration gate. `execution-scope.md` binds common minimum causal
  convergence scope, per-row introducing/invalidating provenance, and the
  source-inspected tag/migration/benchmark reach. Add dual-adapter conformance
  cases alongside each wire surface; task 18 audits the already-running corpus.
- **Consequences:** No gate waits until the final task to execute inherited
  assertions. No final-task ownership inference from historical names. A newly
  exposed production defect still needs minimum causal scope per DEC-175.
  QA Lead performs fresh feasibility; TL planning records no PASS evidence.
- **Refs:** DEC-162/165/166/167/173/175; quality section 3; task specs 1..18.

### DEC-182: Dispose every carry-forward against this phase's actual scope

- **Date:** 2026-09-05.
- **Status:** Planning dispositions; CF-035 policy proposal awaits approval.
- **Context:** The wrap-up explicitly carried CF-024/033/034/035 and re-targeted
  CF-018. Trigger-based items must be checked against the incoming file scopes.
- **Decision:** Use the following complete disposition, and recheck triggers
  against actual diffs at the task-12 checkpoint and wrap-up.

| ID | Phase 22 disposition |
|---|---|
| CF-008 | Re-defer to a retained Issue-cancellation vocabulary use case. Run abandonment is an execution fact, not Issue cancellation; no trigger. |
| CF-018 | Keep exact production-file trigger. `internal/sqlquery/sqlquery.go` is not planned for edit; if a causal change opens it, assign its behavior-preserving split to that task with ordinary/tagged gates. Test-only edits do not fire it. |
| CF-023 | Re-defer to dedicated temporary workflow-maintenance work; no product boundary connection. |
| CF-024 | Schedule task 17 and final QA, VAL-071. Preserve three fresh <30 s samples and exact reference; headroom is a risk, never a new tolerance. |
| CF-026a | Re-defer optional tracker/workflow inspection design. Recovery resume adds only accepted operational facts. |
| CF-026b | Re-defer Issue lifecycle design until its own accepted product decision; no new cancellation model here. |
| CF-029 | Re-defer QA artifact token-validator changes to workflow maintenance; current evidence rules remain binding. |
| CF-030 | Re-defer report-lint diagnostic enhancement; no rejected-token diagnostic change is needed by this product scope. |
| CF-033 | Kickoff package inventory completed in execution-scope.md and linked TRD amendment. AGENTS.md still needs a matching inventory-only maintenance edit; it is outside this TL's permitted write paths. Keep that edit open for the next workflow-maintenance session; no new package depends on it. |
| CF-034 | Schedule workflow baseline/direct tests in task 2 (VAL-079) and evidence baseline/direct tests in task 7 (VAL-035), the first respective recovery reopens. Record before/after coverage; introduce no threshold. |
| CF-035 | Propose whole-file validation on config-consuming paths; leave config-independent help/archive validation/diagnostics available. Task 11, VAL-047, pins both states; no global validation middleware. |

- **Consequences:** CF-033 is partially discharged, not silently closed.
  Closed CF IDs from earlier phases stay archived and are never reused.
- **Refs:** current decisions register; contracts section 12; quality section 3.

### DEC-183: Pin cross-product review and defer workflow-control redesign

- **Date:** 2026-09-05.
- **Status:** Tech Lead planning disposition; independent review remains pending.
- **Context:** The prior phase recommended a delivery retrospective covering
  gate reach, phase/cross-task ledger identities, missing ESC registration,
  and an earlier thin journey. This session reviewed the wrap-up calibration
  and normal kit contract material; it did not conduct a separate deep audit.
- **Decision:** Apply existing gate-reach/provenance rules and early journey
  now through DEC-181. Require immediate ESC registration under existing rules.
  Track workflow identity/registration verifier design as new CF-036 for a
  dedicated workflow-maintenance planning session; do not modify canonical
  controls, prompts, or skills inside this product phase.
- **Cross-product input:** kit repository HEAD
  c5229617220f97ff9228df261e24718b1d1d6d89; `trd/02_INTERFACES.md` sections
  3.2-3.7 and `FID_COMPATIBILITY.md`. Planning aligned unknown/live force
  semantics and bound observed budget evidence. QA must independently compare
  all proposed semantics and record any differences before a Ready QA plan.
  The stale compatibility matrix's historical fid progress is not current
  qualification evidence; no production adapter is claimed qualified here.
- **Consequences:** Full workflow retrospective remains recommended in its
  dedicated scope. This phase needs neither workflow-control changes nor
  another agent activation by Tech Lead. Existing STOP/recovery rules remain.
- **Refs:** Phase 21 Workflow Notes/DEC-179; `qa/qa-lead-handoff.md`.

## Planning Spec Impact

| Document | Planning change |
|---|---|
| GLOSSARY.md | Flags proposed executor-observation, abandonment, and budget-sample terminology under DEC-180. |
| 02_TRD.md | Proposed Phase 22 package/contract index and pinned capability boundary. |
| 03_TRD_SCHEMA.md | Section 8.3 incorporates the proposed exact DDL/Event amendment, without rewriting Phase 21 historical SQL. |
| 04_TRD_CONFIG.md | Proposed CF-035 consuming-path interpretation, pending acceptance. |
| 06_DESIGN_PLAN_CLI.md | Proposed Phase 22 wire-contract pointer. |
| 10_IMPLEMENTATION_ROADMAP.md | Phase 22 planning state and conditional Claims qualification owner. |

This is not a phase-end per-spec verdict. The complete canonical suite,
GLOSSARY, every in-phase deviation, and every QA recommendation must be audited
at wrap-up. No `verified-no-change` assertion is being defaulted during kickoff.

## Triage-Log Audit

No TRG or ESC has been allocated for this phase. The seeded triage log is the
allocation source. Audit every actual entry at wrap-up; do not fabricate events
for pending planning approval or ordinary QA preparation.

## Workflow Notes

Reserved for observed phase execution and final wrap-up. Planning calibration
and deferred deep retrospective work are recorded in DEC-181/183 above.
