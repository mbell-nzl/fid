# Phase 11 Validation Contract - Claims Layer

> Phase: 11 - Claims layer (roadmap v0.1d, slice 2 of 4)
> Author: Tech Lead (draft)
> Last refined: 2026-06-21 by Tech Lead (awaiting QA Lead evidence refinement)
>
> Location: `project-docs/docs/phases/phase-11/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 11 is complete. The contract specifies **what must be proven**; the QA
plan specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Phase Scoping Notes

Phase 11 is the second v0.1d slice: the claims layer. It ships issue claim
acquire, renew, force-acquire, release, list, and transition integration. All
claim writes are gated by the Phase 10 permissions matrix.

The `claims` table, ref-integrity triggers, and v0.1 issue-only claim scope
already exist. Phase 11 should not add a schema migration unless implementation
proves a recorded spec gap. The claim rows already appear in
`03_TRD_SCHEMA.md` section 13; Task 1 appends those rows to the in-code matrix
that Phase 10 intentionally kept limited to implemented actions.

Carry-forward review at kickoff:

- CF-008 remains deferred to v0.2.
- CF-011 remains deferred to v0.1e, but PRD-CLAIM-008a explicitly requires claim
  renewal to project to `activity_log`; Phase 11 must implement that per-command
  projection discipline now.
- CF-013 remains deferred to v0.1e.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 11, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | TRD-ARCH-010; `03_TRD_SCHEMA.md` section 13; roadmap v0.1d Phase 11 | The in-code permissions matrix contains exactly one row for each implemented issue claim action (`issue.claim.acquire`, `issue.claim.renew`, `issue.claim.force_acquire`, `issue.claim.release`) with the same role booleans as `03_TRD_SCHEMA.md`, and `fid policy show` displays those rows from the same source that `Authorize` enforces. | Unit and CLI: matrix completeness test updated for claim actions; `Authorize` allowed/denied pairs for acquire/force/release; `fid policy show` output includes the rows and matches the matrix object. | QA Lead | task-1 | All four claim actions have exactly one in-code row; allowed and denied role/action checks match section 13; policy show includes the same rows. | A claim action is missing/duplicated, row booleans diverge from section 13, or policy show/enforcement drift. | Draft | Link, label, and requirement actions remain future v0.1d slices. |
| VAL-002 | PRD-CLAIM-001; `03_TRD_SCHEMA.md` section 11.1; `06_DESIGN_PLAN_CLI.md` `fid issue claim` | A normal `fid issue claim <id> --by <role>` creates one live issue claim with default TTL 1800 seconds; `--ttl` overrides the TTL; `--note` persists and reads back with the claim. | Domain + CLI: claim command on an issue; DB/read-back or `issue get --section claims`; assert one current `claims` row, `ttl_seconds`, `expires_at`, actor, and optional note. | QA Lead | tasks-1-2 | Default and custom TTL claims persist exactly one current row with the expected holder, expiry, and note. | TTL defaults incorrectly, `--ttl` ignored/misparsed, note is dropped, duplicate current rows appear, or the row is not readable. | Draft | QA review claims are v0.3 and out of scope. |
| VAL-003 | PRD-CLAIM-002; PRD-ERR-003; CF-007 | A non-holder claim attempt against a live non-expired claim is rejected before append with non-zero exit, zero new journal lines/rows, and an error identifying the current holder and expiry time. | CLI: acquire by one role, then acquire by another without `--force`; compare journal line count/DB state before/after and assert stderr contains holder and expiry. | QA Lead | tasks-1-2 | Conflict reject names holder and expiry, exits non-zero, and leaves journal/DB unchanged. | Conflict succeeds, appends/projects anything, or error omits holder/expiry. | Draft | Stale takeover and force-acquire are covered separately. |
| VAL-004 | PRD-CLAIM-003; `03_TRD_SCHEMA.md` section 11.1 | A past-TTL claim no longer blocks acquisition: a new claimant can acquire the issue without `--force`, the prior holder is recorded as superseded in the journal/activity evidence, and the current claim row reflects the new holder/TTL. | Domain + CLI fixture with an expired claim; acquire as a new role; inspect last event/effects and current row. | QA Lead | tasks-1-2 | Stale claim takeover succeeds under `issue.claim.acquire`, records prior holder/expiry as superseded, and current row belongs to the new claimant. | Stale claim blocks normal acquire, requires force permission, loses prior-holder evidence, or leaves stale holder current. | Draft | Human-facing stale reports are covered by claim list. |
| VAL-005 | PRD-CLAIM-004; TRD-ARCH-010; CF-007 | `--force` supersedes a live claim only for a role permitted for `issue.claim.force_acquire`; a denied force attempt is rejected before append with role/action-named permission error. | Domain + CLI: developer live claim; `tl --force` succeeds and records takeover; `developer --force` or another denied role fails with zero append and `issue.claim.force_acquire` in the error. | QA Lead | tasks-1-2 | Permitted force-acquire updates the current holder and records takeover; denied force-acquire appends nothing and names role/action. | Denied force succeeds/appends, permitted force is rejected, or takeover evidence is missing. | Draft | Stale takeover without `--force` is VAL-004. |
| VAL-006 | PRD-CLAIM-008a; CF-011; `03_TRD_SCHEMA.md` section 11.1 | `--renew` succeeds only for the current live holder, updates `claimed_at`, keeps the prior TTL when no new `--ttl` is supplied, uses the supplied TTL when present, and journals/projects an `activity_log` renewal row; non-holder renewal rejects before append. | Domain + CLI: renew by holder with and without TTL; inspect claim row, last event, and `activity_log`; renew by non-holder and assert zero append. | QA Lead | tasks-1-2 | Holder renewal updates time/TTL as specified and adds activity evidence; non-holder renewal rejects with no append. | Renewal loses TTL semantics, omits activity evidence, allows non-holder renewal, or appends on rejection. | Draft | CF-011 broad activity-log reconciliation remains v0.1e. |
| VAL-007 | Roadmap v0.1d Phase 11; `06_DESIGN_PLAN_CLI.md` `fid issue unclaim`; CF-007 | `fid issue unclaim <id> --by <role>` releases only the caller's current live claim; missing/unknown role and wrong-holder release reject before append. | CLI: release by holder succeeds and removes current row; repeat/wrong-holder/missing/unknown role cases fail with zero append and clear errors. | QA Lead | task-2 | Holder release deletes the current claim and appends one event; reject paths append nothing and preserve the claim. | Wrong role releases another holder's claim, missing/unknown role is accepted, rejection appends, or current row remains after holder release. | Draft | Force-release is not a separate command in v0.1. |
| VAL-008 | PRD-CLAIM-007; `06_DESIGN_PLAN_CLI.md` `fid claim list`; `05_DESIGN_PLAN.md` cross-surface fields | `fid claim list` lists current live claims with `ref_type`, `ref_id`, `claimed_by`, `claimed_at`, `ttl_seconds`, `expires_at`, `is_stale`, and `note`; `--by` filters by claiming role; json, toon, markdown, and yaml carry the same data. | CLI: create multiple claims; run list in all four formats with and without `--by`; compare parsed structures to JSON reference. | QA Lead | task-3 | All formats expose the same claims and fields; `--by` returns only claims held by that role. | A format omits/diverges, required fields are missing, or `--by` filters incorrectly. | Draft | `--by` on `claim list` is a read filter, not an actor/permission flag. |
| VAL-009 | PRD-CLAIM-007; decisions principle read-ordering rule | `fid claim list --stale` returns only past-TTL current claims, excludes live claims, and deterministic default ordering follows `claims.rowid` append order unless a spec later names another order. | CLI/DB fixture: seed or create live and stale claims; assert `--stale` output membership and rowid order; assert default list order is stable across rebuild. | QA Lead | task-3 | Stale filter includes only expired current rows; list order is rowid-stable before and after rebuild. | Live claims appear in stale output, stale claims are omitted, or same data returns nondeterministic order. | Draft | Named stale-claims report is v0.2. |
| VAL-010 | PRD-CLAIM-008; roadmap v0.1d Phase 11 | `--by human` can hold an issue claim, reads/lists as `human`, and blocks agents from advancing the issue because they do not hold that live claim. | CLI: human claim; `issue get --section claims` and `claim list` show `claimed_by: human`; developer transition into a started state rejects for missing caller claim. | QA Lead | tasks-2-4 | Human claim persists and displays as human; a non-holder agent cannot advance the issue while the human holds the live claim. | Human claim is rejected/not visible, or an agent can advance despite not holding the human claim. | Draft | Human attention reports are v0.2. |
| VAL-011 | PRD-CLAIM-005; `03_TRD_SCHEMA.md` section 11.1; CF-007 | A transition into a started-category state requires a live claim held by the caller: no claim, wrong holder, and stale holder reject before status append; matching live holder succeeds. | Domain + CLI: transition to `in_progress` for no-claim, wrong-holder, stale-holder, and matching-holder cases; inspect journal lines and status. | QA Lead | task-4 | Only the matching live holder can transition into started; all reject paths append nothing and preserve status. | Started transition succeeds without the caller's live claim, rejects the matching holder, or appends on rejection. | Draft | Workflow edge validity still applies independently before/after the claim gate as implemented. |
| VAL-012 | PRD-CLAIM-006; roadmap v0.1d Phase 11 | A transition out of a started-category state automatically releases the caller's live claim in the same journal event/transaction; after projection, the issue status changes and no current claim remains. | Domain + CLI: claim issue, transition into started, then transition out; inspect last event effects and DB/read-back claim absence. | QA Lead | task-4 | Out-of-started transition appends one event containing status change plus claim release, projects atomically, and leaves no current claim. | Claim remains current, release happens in a separate event, status changes without release, or projection is non-atomic. | Draft | Composite `fid dev finish` is v0.1e. |
| VAL-013 | TRD-ARCH-011 | Concurrent claim acquisition is atomic: simultaneous acquire attempts for the same issue produce exactly one winner; losers see the existing live claim rather than duplicate rows or projection corruption. | Integration test: concurrent acquire attempts against a temp SQLite-backed workspace; assert one successful append/current row and loser errors identify holder/expiry. | QA Lead | task-1 | Exactly one attempt succeeds and exactly one current claim row exists; every loser sees conflict with holder/expiry. | Multiple attempts succeed, duplicate rows occur, loser error is ambiguous, or DB/journal projection corrupts. | Draft | Cross-process locking implementation detail is not prescribed beyond observable atomicity. |
| VAL-014 | `08_QUALITY_PLAN.md` sections 5.2, 5.6; phase goal | Human dogfoods a real `.fid` workspace claim lifecycle and can tell from CLI output who holds a claim, when it expires, when it was released, and why an agent cannot advance a human-held claim. | Human checklist with one row per observable; QA Lead generates from this VAL after deterministic QA passes. | Human | phase-end | Human confirms the claim lifecycle messages and human-held claim behavior are understandable without reading implementation details. | Human cannot identify holder/expiry/release/block reason from CLI output, or checklist is not run. | Draft | QA Lead must fan this into per-observable human checklist rows. |

## Coverage Rules

Confirmed before task execution begins:

- Phase-scope requirements each map to at least one `VAL-NNN`: PRD-CLAIM-001
  to VAL-002; PRD-CLAIM-002 and PRD-ERR-003 to VAL-003; PRD-CLAIM-003 to
  VAL-004; PRD-CLAIM-004 to VAL-005/007; PRD-CLAIM-005 to VAL-011;
  PRD-CLAIM-006 to VAL-012; PRD-CLAIM-007 to VAL-008/009; PRD-CLAIM-008 to
  VAL-010/014; PRD-CLAIM-008a to VAL-006; TRD-ARCH-010 to VAL-001/005;
  TRD-ARCH-011 to VAL-013.
- Every `VAL-NNN` is referenced by a task spec or explicitly marked
  `phase-end` (VAL-014).
- Binary/conditional assertions name both states where applicable:
  VAL-003, VAL-005, VAL-006, VAL-007, VAL-009, VAL-011, and VAL-012.
- Coverage is enforced at the gate, not as a `VAL-NNN`: COV-001 domain package
  coverage and COV-002 `cmd/fid` coverage must stay above thresholds. Each task
  runs the required coverage gate and records package percentages.
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase review,
  or `(human)` verification, and fans VAL-014 into per-observable human rows.
- Standing QA-plan checks carry forward: show-vs-enforce parity for the new claim
  matrix rows (VAL-001) and validate-before-append on every claim rejection
  (VAL-003/005/006/007/011).

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-21 | VAL-001..014 | created | Phase 11 planning draft | roadmap v0.1d Phase 11; PRD-CLAIM-001..008a; TRD-ARCH-010..011 |
