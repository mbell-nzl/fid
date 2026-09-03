# Decisions - Active

> Live carry-forward constraints and accepted product commitments. Historical
> reasoning remains in ADRs, phase archives, and the decision-session record.

## Durable architecture

- [ADR-0001](../../../docs/adr/0001-journal-canonical-db-derived.md): the
  append-only journal is canonical; SQLite is a rebuildable projection.
- [ADR-0002](../../../docs/adr/0002-configurable-roles-default-profile.md):
  superseded for forward product scope by ADR-0010/0011.
- [ADR-0003](../../../docs/adr/0003-three-tier-qa-model.md): superseded for
  forward product scope by ADR-0011.
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
  two internal task groups separated by an independent checkpoint. Attention,
  Stops, budget consumption/exhaustion, automatic respawn/recovery, and export
  import belong to Phase 22; one gate-evidence Obligation and only enforceable
  executor policy values land in Phase 21. Incomplete capabilities remain
  unavailable.
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
- **Phase 21 execution approval (DEC-159):** the QA Lead's Ready plan at
  `e8a636e` and the Human's 2026-08-23 approval accept VAL-001..087 and the
  ordered Heavy task pack. Task 1 is authorized for Phase Runner dispatch;
  every later task remains dependency- and review-gated.
- **Executable ownership (DEC-167, standing):** a task's required commands must
  exercise every file scope it owns. A build-tagged or otherwise gated suite
  covering owned files must have its gate named in that same task, or the scope
  is not assignable to it. Where a gated suite first becomes required several
  tasks after the changes that invalidated its assertions, the accumulated
  convergence belongs to the first task whose gate runs it, recorded with
  per-row provenance. A no-progress recovery budget counts only attempts whose
  directive permitted causal repair of the failing rows. Canonical text is
  `08_QUALITY_PLAN.md` §3.
- **Session liveness and budget accounting (DEC-170, standing):** an attempt
  whose child session is interrupted, cancelled, timed out, or shut down from
  outside, and which never returns a typed result, does not consume the
  no-progress budget. No budget bounds a child session's lifetime; only declared
  per-command budgets are enforced and the orchestrator waits for the session to
  end on its own. Repeated external interruption is triaged as orchestration
  session-liveness, never as `recovery-budget-exhausted`. Where verification cost
  approaches the session window, commit durable work before the long commands and
  append each command result as it terminates. Canonical text is
  `08_QUALITY_PLAN.md` §3 and `09_AI_WORKFLOW.md` step 16.
- **Executable session-liveness controls (DEC-171, standing):** a parent polling
  wait timeout routes to `WAIT`; an external child interruption is recorded as
  typed `INTERRUPTED` evidence and preserves the prior no-progress count. The
  first interruption repairs, the second requires orchestration
  session-liveness triage. Canonical Phase Runner prompts, routing, ledger
  validation, workflow lint, and drill 19 enforce the rule without changing
  product or task-15 contracts.
- **Unavailable capability metadata versus absent behavior (DEC-172,
  standing):** a required `available:false` discovery row is not implementation
  of that capability. Phase 21 retains the exact unavailable `attention.v1`
  row and proves Attention deferral with a paired boundary test that requires
  that row as the sole production match while rejecting any table, Event,
  Effect, leaf, state, or operation surface. A zero-match search is invalid
  because it can hide required negative capability metadata.
- **A control's reach bounds the owning scope (DEC-173, standing):** the
  converse of DEC-167. A task that owns a required absence or inventory control
  owns every path that control walks, and its file scope must cover every
  search root and every package its gates compile. A control reaching beyond
  the listed scope is an incomplete planning assignment, corrected by amending
  the scope before dispatch, never by narrowing the control. Canonical text is
  `08_QUALITY_PLAN.md` §3.
- **The operational journey authors its own workflow (DEC-174):** the
  migration-seeded `task-default`, `subtask-default`, and `epic-default`
  reference workflows are retained V1 choreography and are not the V2 journey's
  authority. The Phase 21 §16 journey defines one workflow
  (`--applies-to task`, one `unstarted`/`started`/`completed`/`cancelled` state,
  exactly the edges `unstarted -> started` and `started -> completed`) and
  assigns it to kind `task` before `issue create`; the transition engine already
  resolves edges and the terminal-success `completed` category from the
  projected workflow tables, so this is project-owned policy under Gate B1/B3,
  not a fid behaviour change. The exact mutation prefix is 22 Events and
  VAL-086 matches. No product, command, public-interface, Event-schema,
  transition, or frozen-workflow change.

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
