# Phase 14 Progress

> Phase: 14 - Telemetry + pricing (roadmap v0.1e, slice 1 of 4)
> Status: COMPLETE - PASS. All 3 tasks done; final review PASS WITH IMPROVEMENTS
> (no findings); human checklist PASS; wrapped 2026-06-28.

## Plan

Phase 14 ships the first v0.1e slice: the self-contained telemetry + pricing
surface that the Phase 15 `fid dev *` composites later open and close.

- `.fid/pricing.toml` defaults: expand beyond Anthropic-only to current
  Anthropic/OpenAI/Google models, with sub-cent precision preserved on round-trip
  (PRD-NFR-010).
- `fid run record`: append one `agent_runs` row (role, model, tokens in/out,
  snapshot pricing, computed `cost_usd`, optional provider reconciliation
  fields); warn-not-reject on an unknown model (PRD-TELEM-001..003a).
- `fid stats tokens` / `fid stats cost`: aggregate by role/model/epic in all four
  read formats, with NULL-cost runs excluded from cost sums (PRD-TELEM-004,
  PRD-OUT-001..004).

No schema migration and no projector change: the `agent_runs` table/columns and
the `run.record` section-13 matrix key already exist; the generic insert
projector applies the new effect.

## Task Index

| Task | Title | Status |
|---|---|---|
| 1 | Pricing defaults (3 providers) + round-trip-safe precision | Complete; Review 1 CONTINUE |
| 2 | `fid run record` (new `internal/runs`, cost + unknown-model warning + provider fields + matrix row) | Complete; Review 1 CONTINUE (integration-gate evidence remediated) |
| 3 | `fid stats tokens` / `cost` (new `internal/stats`, by role/model/epic, four formats) | Complete; validated at final review |

## Validation

Validation contract at `qa/validation-contract.md` includes VAL-001..009.
VAL-009 is the phase-end human dogfood of run recording with cost, the
unknown-model warning, and cost/token aggregates by role/model/epic. QA plan at
`qa/qa-plan.md`; one mid-phase review (`qa/review-1.md`) over a 2+1 task window
(Task 3 closed at final review).

## Progress History

- Kickoff: drafted VAL-001..009 and Tasks 1-3; surfaced PI-1..PI-5 to confirm;
  re-deferred CF-008/011/013/014/015 (none land in Phase 14).
- Task 1 complete: three-provider `DefaultPricingConfig`, round-trip-safe emit
  (`FormatFloat(v, 'f', -1, 64)`), `04_TRD_CONFIG.md` section 3 (DEC-034).
- Task 2 complete: `internal/runs.Record`, `fid run record`, cost computation,
  unknown-model warning, provider fields, single `run.record` matrix row,
  `02_TRD.md` boundary, `06_DESIGN_PLAN_CLI.md` (DEC-035).
- Review 1 (Tasks 1-2): CONTINUE; one Blocking — Task 2 omitted the required
  integration-gate evidence from its completion report (artifact defect; reviewer
  spot-check passed). Escalated as ESC-001, dispositioned do-now; remediation
  commit `0cba161` recorded `PASS go test -tags integration -race ./...` and the
  review addendum cleared the block (TRG-002/TRG-003 RESOLVED).
- Task 3 complete: `internal/stats`, `fid stats tokens`/`cost`, `--by epic`,
  four-format parity, NULL-cost exclusion, `05_DESIGN_PLAN.md` /
  `06_DESIGN_PLAN_CLI.md` (DEC-036).
- Final review: PASS WITH IMPROVEMENTS, no Blocking/Major/Improvement findings.
  VAL-001..008 PASS; VAL-009 CARRY-FORWARD to the human checklist; green
  seven-command phase gate (total coverage 75.7%).
- Human checklist (HUMAN-14a..HUMAN-14c): complete, overall PASS (2026-06-28) →
  VAL-009 final outcome PASS.
- Wrap (2026-06-28): VAL-009 dispositioned PASS; PI-1..PI-5 recorded as
  DEC-034/035/036; roadmap `fid stats tokens` item synced to include `epic`;
  spec-sync verdicts recorded for all specs; TRG-001/002/003 audited; ESC-001
  archived; no new carry-forward. Phase COMPLETE.

## Carry-Forward

- CF-008 re-deferred to v0.2 (resolution-reason field).
- CF-011 re-deferred to Phase 16 (broad `activity_log` reconciliation); Phase 14
  telemetry adds no `activity_log` rows (PI-2), matching the `issue create`
  precedent. Not worsened.
- CF-013 re-deferred to Phase 15 (`complete.go` orchestration cleanup); untouched.
- CF-014 re-deferred to Phase 15 (`cmd/fid/issue.go` / `replay.go` file split);
  not worsened — new surfaces live in `run.go` / `stats.go`, not `issue.go`.
- CF-015 re-deferred to Phase 15 (label write-path function cleanup); untouched.
- No new carry-forward opened in Phase 14.

## Escalations

- ESC-001 (Task 2 missing integration-gate evidence) resolved 2026-06-28 via
  do-now remediation; archived in `tl-escalations.md`. No open escalations at
  wrap.
