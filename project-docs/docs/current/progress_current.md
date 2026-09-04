# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**No phase is active.** Phase 21 completed and wrapped on 2026-09-04 under
DEC-179: 22 tasks, seven review windows, an independent Group 1 checkpoint, and
a fresh independent final review PASS, with `VAL-001..VAL-087` at 87 PASS / 0
FAIL / 0 BLOCKED / 0 CARRY-FORWARD, the Human checklist 4/4 PASS, `RUN-END`
appended, and no open escalation. Phase 20 completed on 2026-08-22.

Phase 21 delivered the one-Issue operational walking skeleton through the built
binary: activation, immutable Gate Registration, one Envelope revision with one
gate-evidence Obligation, generic Actor/Claim/Run authority, fid-executed
Attempts and Receipts with live framed output, redaction and Artifact safety,
guarded terminal-success transition, completion, resume, export create/validate,
and rebuild parity over an exact 22-Event journey. It also completed the owned
V1 removal through migrations `0013`/`0014`: the role matrix, `fid policy show`,
every `--by` attribution, role-bound Claim/Run commands, `fid dev` choreography,
reported `gate_runs`, and completion-report-as-state writers are gone.

Phase 21 detail lives in `../phases/phase-21/`: `decisions.md` (DEC-152..179),
`progress.md`, `qa/final-review-2.md`, `qa/human-checklist.md`, `triage-log.md`,
`tl-escalations.md`, and `run-ledger.md`.

## Planned Boundary

None active. The next planning unit is **Phase 22 - failure and recovery** from
`../../specs/10_IMPLEMENTATION_ROADMAP.md` §5: typed Stops, Attention
disposition, budget consumption and exhaustion, lease/liveness rules, explicit
recovery without automatic respawn, activation recovery, freshness failure,
export import, and recovery-grade export equivalence. `kit.activation.v1`,
`work.claims.v1`, `attention.v1`, and `export.portable.v1` stay unavailable
until they are complete. No implementation is authorized before a Phase 22
brief, validation contract, Ready QA plan, and Human execution approval.

## Approval State

- Phase 21 sign-off: **COMPLETE / WRAPPED 2026-09-04 / DEC-179**.
- Phase 22: **NOT PLANNED**. No task pack, validation contract, QA plan, or
  execution approval exists.
- Product implementation: **UNAUTHORIZED** until a Phase 22 pack is approved.

## Carry-Forward

Active into Phase 22 planning: **CF-024** (100k rebuild headroom; keep the
three-consecutive-fresh-sample rule), **CF-033** (`AGENTS.md` Architecture
package inventory), **CF-034** (in-package tests when `internal/evidence` or
`internal/workflow` is reopened), **CF-035** (uniform unknown-config-key
validation decision), and re-targeted **CF-018**
(`internal/sqlquery/sqlquery.go`). Closed in Phase 21: CF-020, CF-027, CF-028,
CF-032. Re-deferred unchanged: CF-008, CF-023, CF-026a, CF-026b, CF-029,
CF-030.

## Next Work

1. Optional but recommended: a deep delivery-workflow retrospective before
   Phase 22 planning (DEC-179), focused on scope-versus-reach checking, ledger
   identity expressiveness, fail-closed escalation registration, and earlier
   thin end-to-end tasks.
2. Plan Phase 22: operating mode, validation contract, ordered task pack, and
   QA Lead handoff; carry CF-024/033/034/035 and re-targeted CF-018.
3. Merge the `phase-21` branch at the Human's discretion.
