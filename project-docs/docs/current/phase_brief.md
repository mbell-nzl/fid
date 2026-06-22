# Phase Brief - Phase 11 COMPLETE / Phase 12 planning pending

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.

## Status

**Phase 11 (Claims layer, v0.1d slice 2/4) is COMPLETE - PASS**, wrapped
2026-06-22 on branch `phase-11`. No phase is currently active. The next action
is **Tech Lead planning of Phase 12 (Requirements + links)**; this brief is
replaced when Phase 12 planning begins.

## Phase 11 Outcome (for reference)

Phase 11 delivered the claims layer: issue claim acquire/renew/force/release,
`fid claim list`, transition integration (live-claim gate on started entry,
auto-release on exit), human claims, and the claim rows in the Phase 10
permissions matrix. No schema migration was needed.

- All five tasks complete with recorded gates; Reviews 1-5 CONTINUE.
- Final review PASS; full deterministic gate green (coverage total 76.2%,
  domain packages >=80%, binary < 30 MB).
- VAL-001..013 + VAL-015 PASS; VAL-014 PASS via human checklist (overall PASS).
- DEC-026 (no migration; matrix copies §13), DEC-027 (`claim list --by` is a
  read filter), DEC-028 (lock-aware projection freshness, no nested write-lock).
- ESC-008 resolved by DEC-028 + Task 5/VAL-015; TRG-001 audited (BLOCK correct).
- Spec-sync verdicts recorded for all specs; no spec edits required.
- Archives: `../phases/phase-11/` (tasks, qa, decisions, progress, triage-log,
  tl-escalations).

## Carry-Forward Open At Wrap

| ID | Item | Target |
|---|---|---|
| CF-008 | Resolution-reason field + cancel `--reason` | v0.2 |
| CF-011 | Broad `activity_log` reconciliation (renewal slice done) | v0.1e |
| CF-013 | `complete.go` orchestration cleanup | v0.1e |
| CF-014 | Split `cmd/fid/issue.go` (984) + `replay.go` (485) under 400-line target | v0.1e |

## Next Phase Target - Phase 12 (Requirements + links)

Scope from roadmap v0.1d Phase 12 (plan in detail at phase start):

- `fid requirement sync` (minimal): parse canonical PRD/TRD IDs from a source
  doc, upsert pointer rows (`id`, type, summary, source, `last_synced_hash`).
  No requirement get/list or drift report (v0.2). Requirements first so
  issue -> requirement `implements` links have valid endpoints.
- `fid link add`: typed links with auto-inverse projection; validate
  relation/endpoint compatibility; reject links that duplicate dedicated schema.
- `fid link remove`: removes declared + inverse link in one event.
- `fid link list`: outgoing and incoming, all four formats.
- Scope refs: PRD-LINK-001..006a, PRD-LINK-008, PRD-GOV-006; TRD link/schema
  sections. Schema/triggers for `links` and `requirements` already exist
  (v0.1a); Phase 12 adds CLI surface only.
- v0.1d explicit deferrals still apply: requirement get/list + drift (v0.2,
  PRD-GOV-007); link endpoints to decisions/escalations (v0.2, PRD-LINK-007).

## Planning Reminders For Phase 12

- Create `../phases/phase-12/qa/validation-contract.md` with `VAL-NNN`
  assertions before any task execution.
- Review the Carry-Forward Register (CF-008/011/013/014) at kickoff and
  schedule, re-defer, or cancel each.
- Carry forward the Phase 11 lock-ordering lesson: any new write command that
  refreshes the projection under `.fid/write.lock` must use the held-lock
  freshness helper from the start (DEC-028).
