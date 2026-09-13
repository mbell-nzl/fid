# Phase 22 current context before wrap

Archived 2026-09-13. This is a historical snapshot of the pending-retest
context; DEC-220 and qa/final-acceptance.md supersede its pending statuses.

## Snapshot: phase_brief.md

```markdown
# Phase Brief - Phase 22 (Approved)

## Phase Goal

Deliver the one-Issue failure-and-recovery journey: interrupt execution,
reconcile repository state, start a linked successor Run, enforce budgets and
typed Stops, dispose Attention, recover activation, stale old proof,
import/export the same state, rebuild, resume without contradictory facts.

## Status

**Ready for Human retest (DEC-219/task-19 complete).** Branch `phase-22` from merged main
`89d964a84dbdd297fdf472d6f78b001e2efac702`; Phase 21 is wrapped (DEC-179) and
merged. Contracts and task pack are refined against QA R-01..14 and all kit
comparison items (DEC-184..187); Human approved the refined pack on 2026-09-05.
See `../phases/phase-22/runner-handoff.md`. Status is per-task below.

All 18 planned tasks, reviews 1b/2b/3/4/5/6b and the independent task-12
checkpoint are complete. Final-review-2 closed FQ-02 through FQ-07.
Final-review-3 identified omitted force-recovery intervention accounting; root
completed the accepted behavior under DEC-218 in the same Task 1 continuation.
Source c361e9e/671e4e4 and report ab3c255 retain all failed evidence and all twelve
passing task gates, including integration 557.259s and 100k rebuild 27490ms.

Independent final-review-4 closed FQ-01 and qualified the original technical
phase. Human attempt1 then recorded 1 PASS, 2 FAIL and 3 PARTIAL outcomes.
DEC-219/task-19 repairs recovery guidance with one independent reviewer;
MF-01 and both wording improvements are closed on product source5232d9c.
All seven task gates pass, including race/coverage65.5% and integration563.279s.
Fresh clean-source ee7a5ca fixtures and all twelve inspections are verified.
The current [checklist attempt2](../phases/phase-22/qa/human-checklist-2.md) has
blank Human responses; lease expiry is **13 September2026,22:18:47NZST**.
See [task19 report](../phases/phase-22/tasks/task-19-done.md),
[independent review](../phases/phase-22/qa/task-19-review.md), and
[run ledger](../phases/phase-22/run-ledger.md). Phase acceptance and post-Human
wrap-up remain pending; original Human failures are preserved. No open ESC.

## Operating Mode

**Heavy (approved, DEC-188).** Persistence, public JSON, executor authority, and
import publication change; full artifacts and fresh QA feasibility are required
under `09_AI_WORKFLOW.md` section 4. Mode and ordered task list are approved.

## Task Index

All specs are under `../phases/phase-22/tasks/`; reports use `task-N-done.md`.
Tasks are sequential; reviews follow every three tasks.

| Task | Title | Status | Dependency |
|---|---|---|---|
| [1](../phases/phase-22/tasks/task-1.md) | Recovery schema and replay foundation | Complete; FQ-01 independently closed | Runner QA readiness refresh |
| [2](../phases/phase-22/tasks/task-2.md) | Typed Stops and Attention domain | Complete | 1 |
| [3](../phases/phase-22/tasks/task-3.md) | Lease lifecycle and executor observations | Complete | 2 |
| [4](../phases/phase-22/tasks/task-4.md) | Read-only repository reconciliation | Complete | 3 |
| [5](../phases/phase-22/tasks/task-5.md) | Atomic linked Run recovery | Complete | 4 |
| [6](../phases/phase-22/tasks/task-6.md) | Recovery CLI and early process journey | Complete (DEC-194) | 5 |
| [7](../phases/phase-22/tasks/task-7.md) | Durable budget accounting and execution admission | Complete | 6 |
| [8](../phases/phase-22/tasks/task-8.md) | Interrupted Attempts and derived Run closure | Complete (DEC-195) | 7 |
| [9](../phases/phase-22/tasks/task-9.md) | Stop, budget, and closure commands | Complete (DEC-196/197) | 8 |
| [10](../phases/phase-22/tasks/task-10.md) | Activation recovery domain | Complete (DEC-198) | 9 |
| [11](../phases/phase-22/tasks/task-11.md) | Activation CLI and config boundary | Complete (DEC-199/200) | 10 |
| [12](../phases/phase-22/tasks/task-12.md) | Recovery freshness and actionable resume | Complete | 11 |
| [13](../phases/phase-22/tasks/task-13.md) | Recovery archive format and validation | Complete | 12 + independent checkpoint |
| [14](../phases/phase-22/tasks/task-14.md) | Atomic import into an empty target | Complete | 13 |
| [15](../phases/phase-22/tasks/task-15.md) | Import CLI and three-way recovery parity | Complete | 14 |
| [16](../phases/phase-22/tasks/task-16.md) | Retire escalation state and duplicate recovery truth | Complete | 15 |
| [17](../phases/phase-22/tasks/task-17.md) | Maintenance, capability qualification, and rebuild budget | Complete | 16 |
| [18](../phases/phase-22/tasks/task-18.md) | Final recovery conformance and Human validation fixture | Complete | 17 |

## Validation Contract

`../phases/phase-22/qa/validation-contract.md`: 80 Ready assertions, each cited
by an owning task. QA Lead owns 77; Human owns VAL-076..078.

| Assertions | Scope / tasks | Final validator |
|---|---|---|
| VAL-001..024,079 | Foundation/recovery, tasks 1-6 | QA Lead |
| VAL-025..039 | Budgets/Stops/closure, tasks 7-9 | QA Lead |
| VAL-040..052 | Activation/freshness/resume, tasks 10-12 | QA Lead |
| VAL-054..065 | Portable recovery, tasks 13-15 | QA Lead |
| VAL-066..072 | Removal/maintenance/performance, tasks 16-17 | QA Lead |
| VAL-053,073..075,080 / VAL-076..078 | Final conformance, task 18 | QA Lead / Human |

## QA Artifacts

| Artifact | Path under phase-22/ | Status |
|---|---|---|
| Handoff / refinement / plan / consistency | qa/qa-lead-handoff.md, qa/tl-refinement.md, qa/qa-plan.md, qa/consistency-review.md | Accepted DEC-188; plan refreshed before task 1 |
| Mid-phase reviews | qa/review-1b.md, qa/review-2b.md, qa/review-3.md, qa/review-4.md, qa/review-5.md, qa/review-6b.md | PASS |
| Independent checkpoint | qa/checkpoint-after-task-12.md | PASS |
| Final review / Human checklist | qa/final-review-4.md / qa/human-checklist.md | Technical PASS; fresh fixtures verified; six Human judgments pending |

## Active Constraints

- Accepted semantic/wire/DDL authority: `../phases/phase-22/contracts.md`,
  `public-interface.md`, and `schema.md`. DEC-188 ratifies DEC-180/182/184..187.
- Exact removal and executable ownership: `v1-convergence-inventory.md` and
  `execution-scope.md` in that phase directory; DEC-181.
- Every task runs the ordinary full task gate AND tagged integration gate.
  Cases enter with each wire surface; task 6 is the early real-process journey.
- Terminated/live/unknown liveness is closed; unknown always blocks. Known-live
  takeover needs an exact current Human force grant; no automatic respawn. No
  recovery Git mutation or budget reset. Adoption is not trusted evidence.
- Claims capability stays false until Phase 24 parallel semantics qualify.
- CF-024 unchanged three-fresh-sample <30 s rebuild rule, early samples at 1/16;
  R-14 synchronized fixture repair at task 8 (DEC-186). CF-034 coverage tasks
  2/7; CF-035 config task 11; CF-033 recorded; CF-018 fires on sqlquery edits.
- DEC-192: each fenced recovery Event in `internal/journal/project_effects.go` is
  unfenced by its `schema.md` section 5 owner (4,5,7,8,10); task 10 removes the
  empty fence; VAL-003 convergence is scheduled for tasks 5 and 10.
- DEC-193: no test pins an absolute date on a writer-appended Event in a
  live-clock workspace; that failure is DEC-189 nondeterminism, not convergence.
- DEC-194: task 6 delivered task 4's `run_{checkpoint,inspect}.go` domain scope
  with task-4 provenance; QA window 2 audits VAL-013's checkpoint clause with it.
- DEC-195: `schema.md` section 1 admits `cancelled` in
  `execution_receipts.terminal_status` plus one composite-CHECK clause. Task 8
  makes those two exact 0013 edits with task-1 provenance, converges the
  `internal/db` golden DDL string, and opens the closed receipt outcome filter.
  VAL-031 stands as written; never substitute an admitted status.
- Canonical digests: contracts 3.1/3.2; closed 23-action set: interface 2.1.
- Kit dispositions: `kit-contract-mapping.md`; CF-037 retains adapter
  obligations for Phase 25.
- No TL product code/tests, gates, or agent activation; Runner lint unchanged.
- In scope: existing workflow, claims, evidence, validation, project, journal,
  db, export, maintenance, CLI boundaries; no new dependency.

## Parity Matrix

| Axis | Live | Rebuild | Import / fake | Evidence |
|---|---|---|---|---|
| Identity/history | Original entity/Event IDs | Same IDs/order | Same IDs/order | VAL-058,063,073 |
| Recovery | Run/Claim/Stop/Attention facts | Same predicates | Same lineage and unresolved work | VAL-017..021,051/052,063 |
| Evidence/budgets | Observed target and counters | Same derived dimensions | Same proof/counters; local observation explicit | VAL-027..029,048..050,063 |
| Machine/Human | Closed JSON projection | Same JSON | Same semantic projection | VAL-053,075 |
| Archive | Manifest v2 | Canonical equivalent | Re-export equivalent | VAL-054,064 |

## Known Risks

- Limited rebuild headroom; retain reference and per-sample failure semantics.
- Force recovery fences fid authority, not a live process's direct writes.
- Archive validation/publication and cross-resource agreement need fault proof;
  DEC-187 records every pinned-kit difference and no production adapter qualifies.
- Neither a task file scope (DEC-194) nor a contract clause (DEC-195) is
  self-enforcing: a named path with no failing criterion, or a closed-domain
  member with no schema row, fails nothing until code reaches it.

## Open Questions / Success / Exit Criteria

No open questions: all 14 QA refinements are dispositioned; runtime evidence and
the three Human assertion outcomes remain required. Exit criteria: all task gates
and 80 assertions have evidence-backed outcomes, required reviews and Human
checks pass, owned legacy state is removed, every exception is dispositioned, and
Tech Lead records spec sync and sign-off. No release authorization is implied.

## Next Planning Checkpoint

Review 3 passed tasks 7-9, including DEC-196 semantic fake and actual-operation
Stop evidence. Tasks 10-12 lead to review 4 and the independent task-12
checkpoint, which must pass before task 13. Final
deterministic QA and runnable Human fixtures precede tomorrow's Human judgments;
Human outcomes and full phase sign-off remain pending.

```

## Snapshot: decisions_active.md

```markdown
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

```

## Snapshot: progress_current.md

```markdown
# Progress - Current

## Project State

Phase 22 is **ready for Human retest (task-19 complete)** on `phase-22` under DEC-188. The planning base is
merged main `89d964a84dbdd297fdf472d6f78b001e2efac702`. Phase 21 completed
under DEC-179 on 2026-09-04 (87/87 VAL, 4/4 Human), and its branch has since
been merged.

Supplemental task-19 repairs the failed Human guidance under DEC-219. Product
source 5232d9c and independent task-19-review closure PASS are complete. Final
race/coverage passed (CLI425.789s; overall65.5%), integration563.279s, build and
fresh retained-fixture generation18.972s. Root independently verified all twelve
inspection commands, the binary hash and clean source ee7a5ca. Attempt 2 is
`../phases/phase-22/qa/human-checklist-2.md`; its lease expires 13 September2026
at22:18:47NZST. Attempt 1 remains FAIL. No Human assertion is promoted to PASS.

All 18 planned tasks, six review windows and the independent task-12
checkpoint are complete. Original completion reports and causal allocations
remain in the phase archive and run ledger. Final-review-2 independently closed
FQ-02 through FQ-07. Final-review-3 found the remaining accepted force-recovery
intervention accounting missing; DEC-218 allocated its completion to Task 1.

Root completed that repair at c361e9e and fixture follow-up 671e4e4, with report
ab3c255 and clean review checkpoint 3832df1. All twelve unchanged task gates
passed: race CLI 414.292s / export 270.522s, coverage 65.4%, integration 557.259s, and
100000-event rebuild 27490ms below 30000ms. Every failed attempt is preserved.
The first integration run timed out at 600.834s; isolated copies of an immutable
valid prefix reduced the new force test group from 89.800s to 35.396s without removing
cases or rollback comparisons. Independent final-review-4 closes FQ-01 and
passes every technical phase check.
The overall review remains BLOCK solely for the three Human-owned assertions.
Human attempt 1 recorded 1 PASS, 2 FAIL and 3 PARTIAL outcomes. DEC-219 and
task-19 allocate one serialized guidance repair and independent review before
a fresh Human retest. No escalation is open.

The Human delegated completion of tasks 8-18, required reviews, final
deterministic QA, and runnable Human checklist preparation on 2026-09-07.
See `../phases/phase-22/delegate-handoff.md`. Human attempt 1 failed; a retest is required;
this delegation does not claim phase sign-off or release readiness.

Earlier recovery history: task 6 entered a same-task
continuation: its real built-binary journey reached the absent `fid run inspect`
leaf and returned `RECOVERY_REQUIRED`, raising TRG-22-008 / ESC-22-003. DEC-194
resolves it. Task 4 closed COMPLETE without writing the `run.checkpoint`
producer or `run.inspect` read API that its own file scope, `schema.md` section
5, `public-interface.md` section 3, and `contracts.md` section 3.1 all assign to
it, so the owner was never ambiguous and this is undelivered in-scope work, not
a protected boundary. Task 4 is closed with two successors on top and the
Runner's only recovery route is a same-task continuation, so DEC-194 moves
execution of exactly four `internal/workflow` files into task 6 with task-4
provenance recorded. Task 4 stays COMPLETE and unedited. ESC-22-001 (DEC-190/
191) and ESC-22-002 (DEC-192/193) closed earlier. No escalation is open and no
product/spec decision is pending.

## Phase 22 Pack

Heavy mode approved; 18 sequential tasks and 80 Ready validation assertions.
Exact slice/wire/DDL/Event contracts, executable-scope matrix, V1 convergence
inventory, QA Lead handoff, decisions and triage log are in
`../phases/phase-22/`. `phase_brief.md` indexes every task.

Execution started on 2026-09-05 under DEC-188. QA's plan and consistency
assessment remain in `qa/qa-plan.md` and `qa/consistency-review.md`; the
refinement authority is DEC-184..187 and `qa/tl-refinement.md`. The ledger
records completed tasks 1..18 and the task-8 STOP resolved by DEC-195. Product QA
receipts retain their original session attribution; the separate maintenance
report records freshly executed workflow checks only.

## Approval State

- Mode, contracts, task list and refinements: Human approved under DEC-188.
- QA plan and cross-product review: assessment recorded; TL refinements applied.
- Task statuses: original 1..18 complete; all six windows/checkpoint passed.
  The seven initial final-QA findings are independently closed: FQ-02..07
  by final-review-2 and FQ-01 by final-review-4.
- Next: the Human runs the [six checklist checks](../phases/phase-22/qa/human-checklist.md)
  for VAL-076..078, then QA/Human results are recorded before phase acceptance
  and Tech Lead wrap-up. [Final review](../phases/phase-22/qa/final-review-4.md)
  records technical qualification and the remaining Human-only block.
- Fresh fixtures: checkpoints-4132086272, clean source 3832df1, all twelve
  inspection commands independently verified. Successor lease expires
  **13 September 2026 at 12:19:09 NZST**; regenerate and requalify if needed.
- DEC-216 records the Human-approved typed snapshot contract. DEC-217 keeps
  archive ceilings unchanged while splitting small routine checks from the
  explicit production-scale phase gate. DEC-218 completes accepted force
  intervention accounting; it introduces no new product policy or threshold.
- Implementation and independent technical review are complete. Human judgments
  remain pending; no product decision or implementation repair is outstanding.

## Carry-Forward

The register below retains its pre-wrap-up dispositions. The complete trigger
audit and closure/re-targeting follow Human verification at Tech Lead wrap-up.

CF-024: task 17/final QA, unchanged reference rebuild rule; early samples at 1/16.
CF-034: direct workflow/evidence branch tests and before/after coverage in
tasks 2/7. CF-035: config-consuming-path proposal in task 11.
CF-033: kickoff inventory recorded; AGENTS edit awaits separate maintenance.
CF-018: exact production sqlquery.go trigger retained.
Other CFs re-deferred with rationale in DEC-182; CF-036 schedules separate
workflow identity/registration maintenance planning from DEC-183.
CF-037: external kit contract/adapter duties required for Phase 25 qualification.
R-14: task 8 owns process fixture repair; earlier-failure route is DEC-186.
CF-038: live-output golden hardening and post-truncation `output_bytes`
wording, both from DEC-189; neither blocks this phase.
CF-039: a task file scope is not self-enforcing (DEC-194). Planning must give
every named production path an acceptance criterion or VAL citation that fails
when it is absent; wrap-up should decide whether report-lint or QA review should
also diff a report's delivered files against its accepted file scope.

## References

- Phase 22 planning/history: `../phases/phase-22/decisions.md`, `progress.md`.
- Phase 21 closure: `../phases/phase-21/decisions.md` DEC-179 and
  `qa/final-review-2.md`.

```

## Snapshot: tl_escalations_current.md

```markdown
# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None. ESC-22-001 through ESC-22-004 are all resolved and archived in
`../phases/phase-22/tl-escalations.md`.

## Latest Resolution

ESC-22-004 resolved on 2026-09-07 under DEC-195. The contradiction was real and
was a Tech Lead planning omission, not a Developer or TL-Lite defect. Every
accepted record except one agrees that cancellation is its own terminal outcome:
`contracts.md` section 7, `public-interface.md` section 4, VAL-031's
mislabelled-outcome failure condition, the convergence inventory, and QA plan
R-14. Only `schema.md` section 1 was silent, listing five 0013 CHECK member-list
replacements and omitting the one its own contract required. DEC-195 amends the
schema rather than the assertion: `'cancelled'` joins the
`execution_receipts.terminal_status` member list and the composite CHECK gains
one `cancelled` clause admitting signal-only, exit-code-only, or neither, never
both, and never `timed_out = 1`. Execution is an authorized task-8 same-task
continuation with task-1 provenance recorded; task 1 stays COMPLETE. Two bounded
convergence surfaces travel with it: the retained golden DDL string in
`internal/db/operational_kernel_migration_test.go` and the closed outcome filter
in `internal/evidence/receipt_read.go`. Task 8's existing commands already
execute all of them. TRG-22-009's BLOCK, its category, and its refusal to
substitute an admitted status were correct; no no-progress attempt is counted,
because the blocker was a missing authority, not exhausted repair capability.

ESC-22-003 was resolved on 2026-09-06 under DEC-194 (task 4's undelivered
accepted checkpoint/inspection domain scope moved into task 6), ESC-22-002 under
DEC-192/193, and ESC-22-001 under DEC-190/191. All four and their resolutions are
archived in `../phases/phase-22/tl-escalations.md`.

## Prior Phase Resolution

Phase 21 closed on 2026-09-04 (DEC-179) with every allocated escalation
resolved. All seventeen Phase 21 escalations and their resolutions are archived
in `../phases/phase-21/tl-escalations.md`; ESC-21-017 (DEC-177, missing accepted
Receipt read leaves) was the last one raised and closed.

Wrap-up note carried forward as a control weakness, not an open item:
ESC-21-007 through ESC-21-011 were allocated in the Phase 21 triage log and
resolved by DEC-165..DEC-169, but were never written into this register while
they were open. They are now archived from their governing decisions. Register
an escalation here at allocation time, not at resolution time.

## Archive Pointers

- Phase 22: `../phases/phase-22/tl-escalations.md`
- Phase 21: `../phases/phase-21/tl-escalations.md`
- Phase 20: `../phases/phase-20/tl-escalations.md`
- Phase 19: `../phases/phase-19/tl-escalations.md`
- Earlier phases: their `../phases/phase-N/tl-escalations.md` archives.

## Protected Escalation Triggers

Escalate only for a required product, architecture, data-model,
public-interface, dependency, security/data-loss, specification, exhausted
recovery-budget, or unavailable external-capability decision. Ordinary red
tests, failed gates, report defects, QA rework, and incomplete in-scope work
remain bounded repair. Before claiming recovery exhaustion, confirm each
counted attempt returned a typed result and its directive permitted causal
repair. Externally interrupted attempts are session-liveness evidence, not
repair-capability evidence.

Before classifying an owner as ambiguous, read the accepted ownership records:
the candidate task specs' file scopes, `schema.md` section 5 owner column,
`public-interface.md` section 3 task column, and `contracts.md` subject
statements. Behavior that is missing is not the same as behavior that is
unassigned; undelivered accepted in-scope work is bounded repair (DEC-194).

A genuine specification contradiction is different from both, and it does escalate
(DEC-195): when an accepted contract clause names a closed-domain value that the
canonical schema rejects, and no accepted record assigns anyone the amendment,
there is nothing to deliver and nothing to repair. Block, preserve the evidence,
and never substitute an admitted value to make a gate pass. A blocker of that kind
is a missing authority, not exhausted repair capability, so it consumes no
no-progress budget.

```
