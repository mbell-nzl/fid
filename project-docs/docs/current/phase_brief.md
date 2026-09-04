# Phase Brief - No Active Phase

> Status: **NO PHASE ACTIVE.** Phase 21 completed and wrapped on 2026-09-04
> (DEC-179). Phase 22 is not planned.
> Branch: `phase-21` (wrapped; merge at the Human's discretion)
> Product implementation: **UNAUTHORIZED** until a Phase 22 pack is approved.

## Last Completed Phase

Phase 21 - one-Issue operational walking skeleton (Heavy, 22 tasks, two internal
groups, one independent checkpoint). Final outcome: **87 PASS / 0 FAIL / 0
BLOCKED / 0 CARRY-FORWARD** across `VAL-001..VAL-087`, Human checklist 4/4 PASS,
`RUN-END` appended, no open escalation.

Delivered: activation-sourced trust, immutable Gate Registration, one Work
Envelope revision with one gate-evidence Obligation and immutable budgets,
generic Actor/Claim/Run authority, fid-executed Attempts and Receipts with live
framed output, redaction and Artifact safety, the guarded terminal-success
transition, Envelope completion, resume, export create/validate, rebuild parity
over an exact 22-Event journey, and the owned V1 removal through migrations
`0013`/`0014`.

Archive: `../phases/phase-21/` - `decisions.md` (DEC-152..179), `progress.md`,
`contracts.md`, `public-interface.md`, `v1-convergence-inventory.md`,
`qa/validation-contract.md`, `qa/final-review-2.md`, `qa/human-checklist.md`,
`triage-log.md`, `tl-escalations.md`, `run-ledger.md`.

## Next Planning Unit

**Phase 22 - failure and recovery**, per `../../specs/10_IMPLEMENTATION_ROADMAP.md`
§5: interrupt or abandon execution, recover liveness, reconcile repository
state, start a linked Run, retry or stop safely, stale prior proof,
import/export recovery state, rebuild, and resume without contradictory truth.
It completes typed Stops, Attention disposition, budget consumption and
exhaustion, lease/liveness rules, explicit recovery without automatic respawn,
activation recovery, freshness failure, and import identity/history
preservation, and advertises `kit.activation.v1`, `work.claims.v1`,
`attention.v1`, and `export.portable.v1` only when they are complete.

Nothing is authorized yet. Phase 22 needs an operating mode, a validation
contract, an ordered task pack, a QA Lead handoff, a Ready QA plan, and Human
execution approval before any implementation task is dispatched.

## Inputs Waiting For Phase 22 Planning

- **CF-024** - keep the three-consecutive-fresh-sample 100k rebuild rule; final
  Phase 21 QA left only about 2.5-3.0 seconds of headroom under the 30 s target
  while recovery and import will extend the same replay path.
- **CF-033** - reconcile the `AGENTS.md` Architecture package inventory at
  kickoff, against the retained live tree and the accepted Phase 22 boundaries.
- **CF-034** - the first tasks reopening `internal/evidence` or
  `internal/workflow` add direct in-package recovery-branch tests and record
  package coverage before and after; no numeric threshold is introduced.
- **CF-035** - decide whether unknown project-config keys must reject on every
  command or only on config-reading paths, before broadening config validation.
- **CF-018** (re-targeted) - split `internal/sqlquery/sqlquery.go` (596 lines)
  if that production file is reopened.
- **Deep retrospective (recommended, DEC-179)** - before planning, review
  scope-versus-reach checking of task packs, ledger identity expressiveness for
  phase-scoped and cross-task stops, fail-closed escalation registration, and
  whether a phase should carry an early thin end-to-end task.

## Standing Rules That Survive The Phase

These are recorded in `decisions_active.md` and `08_QUALITY_PLAN.md` §3; do not
re-derive them during Phase 22 planning.

- **DEC-167** - a required command must exercise every scope its task owns;
  delayed convergence belongs to the first task whose gate runs it.
- **DEC-173** - the converse: a task's scope must cover every path its required
  controls walk. Amend before dispatch; never narrow the control.
- **DEC-175** - a latent production defect no earlier gate could observe belongs
  to the first task whose gate observes it, with minimum causal scope.
- **DEC-162 / DEC-165 / DEC-166** - retained-test convergence follows the
  invalidating change; a filename, test name, or historical task label is not
  ownership evidence.
- **DEC-170 / DEC-171** - an externally interrupted attempt never consumes the
  no-progress budget; only declared command budgets apply.
