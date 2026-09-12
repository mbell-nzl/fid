# Decisions - Active

> Current commitments and carry-forwards. Historical reasoning stays in ADRs
> and phase decisions; standing delivery rules also live in decisions_principles.md.

## Durable Architecture

- ADR-0001: journal canonical; SQLite is a rebuildable projection.
- ADR-0004/0005: project/kit own intent and policy; fid owns operational facts;
  Issues are the only durable work graph.
- ADR-0006/0007: parallel Builders need distinct Issues/worktrees/disjoint scope;
  interruption recovery uses a new linked Run and non-destructive reconciliation.
- ADR-0010/0011: generic activated Actor/Envelope authority replaces global role
  policy; no pre-dogfood compatibility surface or old fid-data migration;
  JSON is the sole machine format; optional products remain evidence-gated.
- ADR-0012: six vertical scenario phases, 20..25, with explicit removal owners.
  All ADRs live under `../../../docs/adr/`.

## Ratified Product Commitments

- Gates A/B1: fid is the durable operational backend; kit compiles approved
  policy into immutable Envelope Revisions and fid enforces generic constraints.
- Gates B2/B3: immutable activated Gate Registrations drive fid execution.
  Receipts come from observed termination; reported results cannot satisfy proof.
- Gates B4/B4a: Issues are durable work units; exact Issue/Envelope/Actor/
  repository/optional-Obligation context binds Claims and Runs; one mutating
  Claim per Issue and isolated disjoint parallel work.
- Gates B4b/B4c: interrupted Runs never reopen; recovery records liveness and
  reconciliation. Builder rework stays on the same Issue; independent Review
  alone resolves Findings.
- Gate B4d: progress is a projection; material contract/selection changes create
  a revision; removing unfinished work needs an authorized disposition.
- Gates C/D1/D2: all twelve port capabilities qualify before live adoption;
  no global role-policy model; no duplicate process aggregate or generic format.
- Gates D3.1/D3.2: vertical delivery; Phase 21 delivered the one-Issue skeleton;
  Phase 22 owns failure/recovery and conditional capability completion.
- DEC-156: qualifying Evidence guards terminal-success Issue transition;
  Envelope completion reuses the same evaluator after the Issue is terminal.
- DEC-157/177: Phase 21 delivered wire authority is
  `../phases/phase-21/public-interface.md`; its removal inventory is
  `v1-convergence-inventory.md`. Public enums are kebab-case; bootstrap has no role.
- DEC-179: Phase 21 wrapped with 87/87 VAL and 4/4 Human PASS. Its local
  DEC-152..178 decisions are archived in `../phases/phase-21/decisions.md`.

## Standing Delivery Constraints

- DEC-158: report lint is a two-argument Runner pre-dispatch control, not a
  Developer task command; 100k rebuild needs three fresh <30 s samples (CF-024).
- DEC-162/165/166: retained-test convergence follows the accepted invalidating
  change; establish introducing commit provenance, never filename/task labels.
- DEC-167/173: a task's required gates execute every owned scope, and its scope
  covers every path its gates walk. Amend incomplete assignments before dispatch.
- DEC-175: a latent implementation defect first observable under a new gate
  belongs to that task with minimum causal source/regression scope.
- DEC-189: a nondeterministic retained test that first fails under an existing
  required gate is neither convergence nor DEC-175. Assign it to the observing
  task with a scope amendment, repair the nondeterminism at its source, and
  never re-baseline a golden to a load- or instrumentation-dependent value.
- DEC-170/171: external interruptions do not consume no-progress budget; polling
  timeout means wait. Only declared command budgets apply; preserve real outcomes.
- DEC-172: an unavailable capability row is required negative metadata, not an
  implementation of that capability. Replace only when owning semantics qualify.
- DEC-178: shipped migration SQL and canonical slice DDL must agree exactly.
- DEC-002: autonomous TL-Lite resume remains enabled under the current workflow.
- DEC-190: separate Human-authorized verifier maintenance resolves ESC-22-001;
  ordinary pending-triage correction and normal restart retain task/attempt
  authority, and active owning-review triggers work across recovery boundaries.
- DEC-191: canonical prompts retain the committed Codex model choices and
  the TL Bash capability; generated provider copies are synchronized.
- DEC-192: a progressive availability fence is removed key-by-key by the task
  that owns the fenced Event; the owning table, not the fence file's author,
  fixes ownership, and the last owner removes the empty fence and its guard.
  Fence removal is implementation of an accepted contract, never a boundary
  change, and the displaced fence assertions are scheduled convergence.
- DEC-193: a calendar-dependent retained test is DEC-189 nondeterminism, not
  convergence and not DEC-175. Assign it to the observing task, repair the date
  binding at its source with recorded introducing-commit provenance, and never
  pin an absolute date on an Event appended through the journal writer into a
  live-clock workspace.
- DEC-194: an accepted allocation that exists but was never delivered is bounded
  repair, not an ambiguous owner. Read the task spec, schema owner column,
  public-interface task column, and contract subjects before classifying.
- DEC-195: an accepted contract clause that adds a closed-domain member requires
  a matching schema row. Where the row is missing, amend the schema and keep the
  assertion; never substitute an admitted value or weaken the assertion. Assign
  execution to the first task that cannot pass without it, with provenance to the
  original owner recorded rather than relabelled.
- Current V1 delivery scaffolding is binding engineering procedure, not forward
  fid product scope. No failed gate authorizes a protected definition change.

## Phase 22 Approval

DEC-196 makes accepted semantic fake behavior, Human resume rendering and
retained Human fixture generation explicit in tasks 9/12/18 under the Human's
delegated execution authority. Existing gates and VAL ownership remain;
task 18 adds a concrete retained-fixture generation command. See phase decisions.

DEC-197 assigns task 9 the omitted domain/read/frame support required by its
accepted leaves and Stop guards, plus causal replay/error/retirement cleanup.
The task spec names each bounded path; prior task provenance and failed
receipts remain intact, with no protected contract or gate change.

`../phases/phase-22/decisions.md` records:
DEC-180 proposed Heavy mode/18-task/80-VAL product contract;
DEC-181 executable ownership and early journey;
DEC-182 complete CF disposition including proposed CF-035 behavior;
DEC-183 pinned kit comparison and separate workflow-maintenance recommendation.
DEC-184 closed wire/authority/observation/Stop refinements;
DEC-185 canonical dirty/reconciliation digest;
DEC-186 boundary evidence, early rebuild samples and bounded R-14 repair;
DEC-187 complete kit comparison dispositions and CF-037 qualification work.
DEC-188 Human acceptance of Heavy mode, the full refined pack and execution.

Human approval on 2026-09-05 ratifies DEC-180/182/184..187. All 18 planned
tasks, six review windows and the independent task-12 checkpoint are complete.
Final-review-2 closed FQ-02 through FQ-07; root completed final-review-3's force
accounting finding under DEC-218. Independent final-review-4 qualifies technical
readiness and fresh Human fixtures; FQ-01 is closed. Human attempt 1 recorded one PASS, two FAIL and three PARTIAL outcomes;
DEC-219/task-19 repairs the guidance before a fresh VAL-076..078 retest. Delegation
covers implementation, QA and checklist preparation, not Human acceptance or
release. See `../phases/phase-22/delegate-handoff.md`. Branch is `phase-22`;
Phase 21 is merged.

DEC-198 resolves contradictory prepare-time prose in favor of the already
closed activation protocol: prepare freezes the predecessor and five resources;
component staging supplies immutable internal content; commit selects the exact
two digests with resource agreement. No new wire/DDL/Event or weakened VAL row.

DEC-199 gives task 11 bounded domain error enrichment for the accepted closed
activation-incomplete result, preserving task 10's predicates. It also makes
the existing Human activation direction rendering explicit. No command or VAL
definition changed.

DEC-200 explicitly adds stop and activation to retained ReasonSubjectType;
activation disagreement uses activation-incomplete with the Activation ID.
Component names remain in components/detail, with no private subject/code enum.

DEC-201 assigns the omitted derived Budget output state to task 12 with task-7
provenance, preserving policy/accounting/history. It explicitly adds only
effective_status to the Resume Claim projection so elapsed leases are visible
without rewriting durable Claim status.

DEC-216 accepts complete typed current-family idempotency snapshots, including
pending Attempts, derived terminal Receipts and committed denials. Public
response behavior is preserved. DEC-217 keeps the 64MiB/1GiB/8GiB/1-million-entry archive
ceilings unchanged; small routine fixtures exercise the same limiter and the
original full matrices remain mandatory through the archive_scale phase gate.
DEC-218 completes already accepted force intervention charges, immutable chain,
successor boundary Stops and exact replay/import validation. See the phase
[decisions](../phases/phase-22/decisions.md) and Task 1 report for authority,
causal provenance and preserved runtime-failure evidence.

## Carry-Forward Register

| ID | Current disposition |
|---|---|
| CF-008 | Re-defer until a retained Issue-cancellation vocabulary use case; Run abandonment does not fire it. |
| CF-018 | Split production `internal/sqlquery/sqlquery.go` (596 lines) when that file next changes. No planned edit; test-only convergence does not fire it. |
| CF-023 | Re-defer temporary repository-delivery tooling to dedicated workflow maintenance. |
| CF-024 | Task 17/final QA, VAL-071: unchanged three-fresh-sample <30 s reference rule; task 1/16 add single early samples (DEC-186). |
| CF-026a | Candidate inspection projection; not scheduled beyond accepted recovery resume. |
| CF-026b | Issue lifecycle design requires its own product decision; deferred. |
| CF-029 | QA review evidence-token validator work deferred to workflow maintenance. |
| CF-030 | Report-lint rejected-token diagnostics deferred; no new product need. |
| CF-033 | Kickoff inventory recorded in Phase 22 execution-scope.md/TRD. Separate AGENTS maintenance still owns matching package inventory and correction of its nonexistent quality §5.2 pointer (QA plan incidental finding). |
| CF-034 | Task 2 workflow and task 7 evidence: direct recovery-branch tests plus before/after package coverage; VAL-079/035, no numeric threshold. |
| CF-035 | Accepted under DEC-188: whole-config validation on consuming paths, preserving config-independent diagnostics. Task 11/VAL-047. |
| CF-036 | New from DEC-183: dedicated workflow-maintenance planning for phase/cross-task ledger identities and fail-closed ESC registration; product Phase 22 does not change controls. |
| CF-037 | DEC-187: kit V2-3 adapter/Phase 25 qualification must prove activation export binding, terminalState expectation mapping, scope/asOf/validate-only, dedup-key agreement and explicit unmapped Stop-code representations. Exact obligations: Phase 22 kit-contract-mapping.md. |
| CF-038 | New from DEC-189: harden `cmd/fid/live_output_test.go` against remaining same-stream read coalescing when a task next opens it with cause, and qualify the terminal frame's `output_bytes` as a post-truncation lower bound in the next accepted live-output wire revision. Neither item blocks Phase 22. |

All rows were dispositioned against the Phase 22 task scope in DEC-182..187;
check triggers again at task-12 checkpoint and wrap-up. Closed CF IDs remain in
phase archives and are never reused. Optional product ideas remain in ADR-0011's
candidate ledger rather than this register.

## Environment

Go/tool caches remain redirected into gitignored repository directories by
`.codex/config.toml`; changing that configuration requires a fresh session.
No cache configuration change is planned.

- DEC-219: Human approved one supplemental task-19 guidance repair, serial root
  implementation and one independent review. Imported observation provenance,
  adoption versus proof and activation direction are explained without changing
  JSON or mutation authority. Original Human failures remain; retest required.
