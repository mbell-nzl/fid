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
- **Gate D3.2:** one Heavy Phase 21 delivers the one-Issue walking skeleton in
  two groups with an independent checkpoint. Phase 22 owns Attention, Stops,
  recovery, and export import; incomplete capabilities remain unavailable.
- **Gate D3.2 terminal amendment (DEC-156):** gate execution, qualifying
  Receipt linking, and satisfied Evidence precede the terminal-success Issue
  transition. Envelope completion later re-evaluates the same Obligation and
  requires that terminal Issue while reusing one shared Obligation evaluator.
- **Phase 21 planning correction (DEC-157):** the complete wire contract and
  V1 convergence inventory are protected planning authorities; public enums
  are kebab-case; final bootstrap has no role Effect; task 20 owns removal/0014
  before task 21 final-schema parity.
- **Phase 21 QA refinement (DEC-158):** report lint is a two-argument Phase
  Runner pre-dispatch control, not a Developer task command. The unchanged
  100k rebuild target is three consecutive self-asserted `<30 s` samples on
  the recorded Phase 21 reference instance; one failure fails VAL-083.
- **Phase 21 execution approval (DEC-159):** QA Ready commit `e8a636e` and the
  2026-08-23 Human approval accept VAL-001..087 and the ordered Heavy task pack.
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
- **The operational journey authors its own workflow (DEC-174):** the
  Phase 21 §16 journey defines and assigns its task workflow before Issue
  creation; the exact mutation prefix is 22 Events and VAL-086 matches.
- **First-observing ownership of a latent production defect (DEC-175,
  standing):** where an accepted contract was implemented wrongly and no earlier
  required gate could observe it, the defect belongs to the first unit of work
  whose gate observes it. That unit receives minimum causal source and
  regression scope; the originating unit is not reopened.
- **Task 22 append-only attempt identity (DEC-176):** TRG-21-022's immutable
  TRIAGE already carries attempt 3, so its schema-v2 STOP resume and matching
  recovery use attempt 4. DEC-175 remains the product/scope authority; DEC-176
  changes only recovery identity and leaves no-progress at `0/2`.
- **Task 22 Receipt completion (DEC-177):** Task 22 owns accepted read-only
  Receipt get/list plus the minimum pre-baseline migration/schema/golden change
  that durably stores canonical `rejected_artifacts`. It corrects only stale
  Resume prose to `fid resume`; accepted operations, wire shapes, filters,
  errors, Events, 22-Event budget, VALs, and commands stay fixed.

The detailed scenarios and reasoning are in
`../reviews/2026-08-08-fid-v2-direction-decision-session.md` and the current
specification suite.

## Open decision

None for Phase 21. QA readiness and Human execution approval are complete under
DEC-159; protected changes discovered during execution still require their
normal authority.

## Carry-forward disposition

| ID | Disposition after ADR-0011 |
|---|---|
| CF-008 | Candidate Issue-tracker vocabulary improvement; schedule only with a retained cancellation use case. |
| CF-017 | Cancel. `internal/dev` is displaced role choreography. |
| CF-018 | Keep as trigger-based engineering debt if `internal/sqlquery` is reopened; `fid sql` remains diagnostic only. |
| CF-020 | Scheduled Phase 21 task 8 / VAL-026: Event sequence replaces synthetic future occurrence time when the retained Issue/Event path converts to v2. |
| CF-022 | Cancel. Generic Markdown formatter characterization is retired. |
| CF-023 | Reassess only as temporary repository-delivery tooling while the V1 development procedure remains active; it is not product scope. |
| CF-024 | Scheduled Phase 21 task 21 / VAL-083: exercise the accepted 100k-Event rebuild target under the evidence event load. |
| CF-026a | Candidate tracker/workflow inspection projection; no committed slice. |
| CF-026b | Candidate Issue-workflow lifecycle design; requires a decision before scheduling. |
| CF-027 | Scheduled Phase 21 task 8 / VAL-027: decompose reopened `internal/issues/create.go`; the `writeMarkdown` half remains cancelled. |
| CF-028 | Split `internal/bootstrap/init.go` staging/fsync helpers only when a later change reopens that file; the trigger did not fire after Task 3b. |
| CF-029 | Extend evidence-token validation to QA review artifacts while the V1 delivery workflow remains active. |
| CF-030 | Improve `report-lint.sh` rejected-token diagnostics so they state the accepted token shape. |
| CF-031 | Closed in Phase 20: PASS termination validation and expected-nonzero parsing are implemented and drilled. |
| CF-032 | Scheduled Phase 21 task 12: add direct invalid-input tests before `internal/issues/revision.go` changes. |

Optional product ideas are not carry-forwards. Their evidence triggers live only
in ADR-0011's concise candidate ledger.

## Temporary delivery constraints

fid still uses the repository's V1 seven-role development procedure until the
V2 kit and backend port are ready. `AGENTS.md` remains binding for current
implementation work, including TDD, evidence, command, phase, and hard-stop
rules. Those requirements are temporary engineering scaffolding, not forward
fid product requirements.

No implementation phase starts from the historical version roadmap. Phase 20
completed and wrapped on 2026-08-22: all tasks and reviews closed, all 46 VAL
assertions PASS, the Human checklist is 3/3 PASS, and no escalation remains.
Its accepted and delivered slice contract is
`../phases/phase-20/contracts.md`; Gate D3.2 is accepted in Phase 21 DEC-153.

Two Phase 20 constraints bind execution: fail-closed V1 authority stays until
Phase 21 proves generic authority (DEC-113), and DEC-117's Gate D3.2 condition
is now satisfied only through the exact DEC-153/154 boundary; it creates no fid
Routine entity.
Canonical workflow authoring exposes states/transitions only, never the V1
role/assurance `--require` grammar; retained legacy `dev` status mutations still
advance complete Issue aggregate revision until Phase 21 deletes them (DEC-128).

## Environment

Go and tool caches remain redirected into gitignored repository directories by
`.codex/config.toml`. A fresh agent session is needed after changing that
configuration.
