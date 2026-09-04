# Decisions - Active

> Live carry-forward constraints and accepted product commitments. Historical
> reasoning remains in ADRs, phase archives, and the decision-session record.

## Durable architecture

- [ADR-0001](../../../docs/adr/0001-journal-canonical-db-derived.md): the
  append-only journal is canonical; SQLite is a rebuildable projection.
- [ADR-0004](../../../docs/adr/0004-operational-work-state-boundary.md): the
  project and kit own intent and policy; fid owns Issue and operational facts.
- [ADR-0005](../../../docs/adr/0005-issues-are-the-durable-work-graph.md): Issues
  are the only durable work graph; there is no persistent Work Item entity.
- [ADR-0006](../../../docs/adr/0006-isolate-parallel-builder-mutations.md):
  parallel mutation requires distinct Issues, worktrees, and disjoint scopes.
- [ADR-0007](../../../docs/adr/0007-recover-interruptions-with-linked-runs.md):
  interruption recovery creates a linked new Run after reconciliation.
- [ADR-0010](../../../docs/adr/0010-supersede-phase-20-and-defer-configurable-role-policy.md):
  the old role-policy Phase 20 is superseded and generic authority belongs to
  the port.
- [ADR-0011](../../../docs/adr/0011-realign-fid-around-the-durable-operational-kernel.md):
  retain the Issue tracker and durable kernel, remove the pre-dogfood V1
  compatibility surface, use JSON as the sole machine representation, and gate
  optional features on observed evidence.
- [ADR-0012](../../../docs/adr/0012-deliver-convergence-as-vertical-scenario-phases.md):
  deliver convergence through six phase-sized end-to-end scenarios, beginning
  with a clean standalone Issue journey in a newly planned Phase 20.

## Ratified fid/V2 gates

- **Gate A:** fid is the durable operational backend, not the owner of project
  prose, engineering policy, orchestration, or provider materialization.
- **Gate B1:** the kit compiles approved policy into an Envelope Revision; fid
  enforces its generic integrity, authority, evidence, freshness, and lifecycle
  constraints without recalculating policy.
- **Gate B2:** fid is the primary local executor for registered gates and
  creates Receipts from observed execution. Reported outcomes cannot satisfy
  deterministic Obligations.
- **Gate B3:** the project owns approved Gate Definition values, the kit owns
  schema/validation/activation, and fid owns immutable registrations and
  execution. fid never executes mutable authoring configuration directly.
- **Gate B4:** Issues are the only durable work units. Envelope Revisions select
  Issues; Claims and Runs bind exact Issue, revision, Actor, repository, and
  optional Obligation context.
- **Gate B4a:** one active Mutation Claim per Issue. Parallel Builders need
  separate selected Issues, worktrees, and provably disjoint mutation scopes.
- **Gate B4b:** interrupted Runs are terminal. Continuation requires confirmed
  liveness resolution, non-destructive reconciliation, and a linked new Run.
- **Gate B4c:** Builder rework stays on the original Issue; only a new
  independent Review Record may resolve Findings.
- **Gate B4d:** partial progress is a projection. Selection or material contract
  change creates a new Envelope Revision; unfinished removal needs an explicit
  authorized disposition.
- **Gate C:** implement and qualify all twelve backend-port capabilities in
  bounded slices before any live migration.
- **Gate D1:** do not implement Phase 20's global role-policy model. Generic
  Actor capability and Envelope authority replace it.
- **Gate D2:** apply ADR-0011's durability test; retire duplicate process
  aggregates, role choreography, generic TOON/YAML/Markdown output, and V1
  compatibility. TUI/dashboard and search lead the candidate ledger.
- **Gate D3.1:** use vertical scenario phases rather than horizontal subsystem
  phases. Each phase owns explicit legacy removals and may contain multiple
  bounded Developer tasks under the current V1 delivery process.
- **Gate D3.2:** the one-Issue walking skeleton was delivered by Heavy Phase 21
  in two groups with an independent checkpoint. Phase 22 owns Attention, Stops,
  recovery, and export import; incomplete capabilities remain unavailable.
- **Gate D3.2 terminal amendment (DEC-156, delivered):** gate execution,
  qualifying Receipt linking, and satisfied Evidence precede the
  terminal-success Issue transition. Envelope completion later re-evaluates the
  same Obligation and requires that terminal Issue, reusing one shared
  evaluator.
- **Delivered wire authority (DEC-157/177, standing):**
  `../phases/phase-21/public-interface.md` is the delivered public contract and
  `../phases/phase-21/v1-convergence-inventory.md` the removal record. Public
  enum members are lowercase kebab-case; bootstrap writes no role Effect.
- **Delivery controls (DEC-158, standing):** report lint is a two-argument Phase
  Runner pre-dispatch control, not a Developer task command. The 100k rebuild
  target is three consecutive self-asserted `<30 s` samples on the recorded
  reference instance; one failure fails the assertion (CF-024).
- **Executable ownership (DEC-167, standing):** required commands exercise every
  owned scope, including tagged suites. Delayed convergence belongs to the first
  task whose gate runs it, with per-row provenance. No-progress counts only when
  the directive permitted causal repair. See `08_QUALITY_PLAN.md` §3.
- **Session liveness and budget accounting (DEC-170, standing):** external
  interruption without a typed result does not consume no-progress; repeated
  interruption is session-liveness, not exhaustion. Only declared command
  budgets apply. Commit before long commands and append results as they
  terminate. See `08_QUALITY_PLAN.md` §3 and `09_AI_WORKFLOW.md` step 16.
- **Late execution refinements (DEC-171/172):** polling timeout routes to WAIT;
  externally interrupted children retain no-progress state. The unavailable
  `attention.v1` discovery row is required negative metadata, while a paired
  boundary test proves that no Attention behavior exists.
- **A control's reach bounds the owning scope (DEC-173, standing):** the
  converse of DEC-167. Task scope covers every path its required controls walk;
  an incomplete assignment is amended before dispatch, never narrowed.
- **First-observing ownership of a latent production defect (DEC-175,
  standing):** where an accepted contract was implemented wrongly and no earlier
  required gate could observe it, the defect belongs to the first unit of work
  whose gate observes it. That unit receives minimum causal source and
  regression scope; the originating unit is not reopened.
- **Migration authority (DEC-178, standing):** the shipped migration files are
  the exact authority for their own SQL; `03_TRD_SCHEMA.md` §8.2 matches
  `0013`/`0014` byte for byte.
- **Phase 21 sign-off (DEC-179):** Phase 21 is complete and wrapped with 87
  PASS / 0 FAIL / 0 BLOCKED / 0 CARRY-FORWARD and a 4/4 Human checklist. The
  delivered one-Issue operational walking skeleton and its owned V1 removal are
  accepted; incomplete capabilities remain unavailable until their owning phase
  completes them. Phase-local execution decisions DEC-152..178 are archived in
  `../phases/phase-21/decisions.md`.

The detailed scenarios and reasoning are in
`../reviews/2026-08-08-fid-v2-direction-decision-session.md` and the current
specification suite.

## Open decision

None. Phase 21 is signed off under DEC-179 and no phase is active. Phase 22
planning is the next unit of work and requires its own brief, validation
contract, Ready QA plan, and Human execution approval before implementation.

## Carry-forward disposition

Audited against the Phase 21 diff at wrap-up (DEC-179 item 6). A trigger-based
row states the exact file or condition that fires it.

| ID | Disposition after ADR-0011 |
|---|---|
| CF-008 | Candidate Issue-tracker vocabulary improvement; schedule only with a retained cancellation use case. Trigger did not fire in Phase 21. |
| CF-018 | **Re-targeted 2026-09-04.** Behavior-preserving split of `internal/sqlquery/sqlquery.go` (596 lines, over the 400-line target) when that production file next changes. Phase 21 reopened only `sqlquery_test.go` under DEC-162 convergence, which does not fire this row; `fid sql` remains diagnostic only. |
| CF-023 | Reassess only as temporary repository-delivery tooling while the V1 development procedure remains active; it is not product scope. |
| CF-024 | **Active into Phase 22.** Phase 21 VAL-083 passed at 27.3-27.5 s against the 30 s target, roughly 2.5-3.0 seconds of headroom. Re-issue the same three-consecutive-fresh-sample rule while recovery and import extend the replay path. |
| CF-026a | Candidate tracker/workflow inspection projection; no committed slice. |
| CF-026b | Candidate Issue-workflow lifecycle design; requires a decision before scheduling. |
| CF-029 | Extend evidence-token validation to QA review artifacts while the V1 delivery workflow remains active. Trigger did not fire in Phase 21. |
| CF-030 | Improve `report-lint.sh` rejected-token diagnostics so they state the accepted token shape. Phase 21's one rejection (`COMPLETE category must be none: qa-rework`) already stated the accepted shape, so the row stays deferred. |
| CF-033 | At Phase 22 kickoff, reconcile the `AGENTS.md` Architecture package inventory against the retained live tree and the accepted Phase 22 boundaries. |
| CF-034 | The first Phase 22 tasks reopening `internal/evidence` and `internal/workflow` add direct in-package recovery-branch tests and record package coverage before/after; no numeric threshold. |
| CF-035 | At Phase 22 planning, decide whether unknown project-config keys must reject on every command or only on paths that consume config, before broadening config validation. |

Closed or cancelled, retained here only so an ID is never reused: CF-017 and
CF-022 cancelled; CF-031 closed in Phase 20; CF-020, CF-027, CF-028, and CF-032
closed in Phase 21. Their closing evidence is in
`../phases/phase-21/decisions.md` DEC-179 item 6 and
`../phases/phase-21/progress.md`.

Optional product ideas are not carry-forwards. Their evidence triggers live only
in ADR-0011's concise candidate ledger.

## Temporary delivery constraints

fid still uses the repository's V1 seven-role development procedure until the
V2 kit and backend port are ready. `AGENTS.md` remains binding for current
implementation work, including TDD, evidence, command, phase, and hard-stop
rules. Those requirements are temporary engineering scaffolding, not forward
fid product requirements.

No implementation phase starts from the historical version roadmap. Phase 20
completed and wrapped on 2026-08-22 (46/46 VAL PASS, 3/3 Human) and Phase 21
completed and wrapped on 2026-09-04 (87/87 VAL PASS, 4/4 Human). Their accepted
and delivered slice contracts are `../phases/phase-20/contracts.md` and
`../phases/phase-21/contracts.md` with
`../phases/phase-21/public-interface.md` as the delivered wire contract.

The two remaining Phase 20 execution constraints are now discharged by Phase 21
delivery: DEC-113's fail-closed V1 authority is replaced by generic activated
Actor and Envelope authority, and DEC-128's retained legacy `dev` status
mutations are deleted with the rest of the V1 role choreography. Canonical
workflow authoring still exposes states and transitions only, never the V1
role/assurance `--require` grammar, and Gate D3.2 creates no fid Routine entity.

## Environment

Go and tool caches remain redirected into gitignored repository directories by
`.codex/config.toml`. A fresh agent session is needed after changing that
configuration.
