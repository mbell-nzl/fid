# Phase Brief - Phase 12 (Requirements + Links) PLANNED

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.

## Status

**Phase 12 (Requirements + links, v0.1d slice 3/4) is PLANNED**, awaiting
human review of the validation contract and task list. Phase 11 (claims) is
COMPLETE - PASS. Validation contract and four task specs are drafted under
`../phases/phase-12/`.

## Goal

Ship requirement pointers and typed links so the v0.1d link surface works end to
end: `fid requirement sync` (minimal), `fid link add`, `fid link remove`,
`fid link list`. Requirements ship first so issue -> requirement `implements`
links have valid endpoints. All writes are gated by the Phase 10 matrix; new
matrix rows (`requirement.sync`, `link.add`, `link.remove`) are appended this
phase. Schema and triggers for `links`/`requirements` already exist (v0.1a); no
migration is expected.

Scope refs: PRD-LINK-001..006a, PRD-LINK-008, PRD-GOV-006; TRD-ARCH-010;
`03_TRD_SCHEMA.md` sections 5.2, 10, 13, 14.2-14.3. Roadmap v0.1d Phase 12.

## Ordered Tasks

1. **Task 1 - Requirement pointers + `fid requirement sync` + matrix row.**
   New `internal/requirements`; parse canonical PRD/TRD IDs from a source doc;
   upsert pointer rows in one event; created/updated/unchanged counts;
   `requirement.sync` matrix row. Covers VAL-001 (partial), VAL-002..004.
2. **Task 2 - Links domain + matrix rows + `fid link add`.** New
   `internal/links` (relation/inverse registry + validation); declared + inverse
   in one event; `link.add`/`link.remove` matrix rows. Covers VAL-001
   (completes), VAL-005..008. Depends on Task 1.
3. **Task 3 - `fid link remove`.** Declared + inverse delete in one event +
   `activity_log` removal evidence; nonexistent-link rejection. Covers VAL-009.
   Depends on Task 2.
4. **Task 4 - `fid link list` + end-to-end flow.** Outgoing/incoming, four
   formats, resolved titles; integration test + rebuild parity. Covers VAL-010,
   VAL-011. Depends on Tasks 1-3.

## Validation Contract Index

Full contract: `../phases/phase-12/qa/validation-contract.md`.

- VAL-001 matrix rows (`requirement.sync`, `link.add`, `link.remove`) match
  section 13 and policy show (tasks-1-2).
- VAL-002 requirement sync upserts in one event with correct fields/counts;
  unchanged sync appends nothing (task-1).
- VAL-003 requirement sync malformed-input rejection before append (task-1).
- VAL-004 `requirement.sync` permission (tl/human) before append (task-1).
- VAL-005 link add: declared + inverse in one event; duplicate rejected (task-2).
- VAL-006 issue -> requirement `implements` link + `implemented_by` inverse;
  missing requirement rejected (task-2).
- VAL-007 relation/endpoint incompatibility + unresolvable endpoint rejected
  (task-2).
- VAL-008 relation registry equals the documented seven relations + inverses;
  unknown relation / self-link rejected (task-2).
- VAL-009 link remove: declared + inverse in one event + activity log;
  nonexistent rejected (task-3).
- VAL-010 link list: outgoing/incoming, four-format parity, resolved titles
  (task-4).
- VAL-011 end-to-end flow + rebuild parity, one event per write, no nested lock
  (task-4).
- VAL-012 human dogfood of sync/add/list/remove legibility (phase-end).

## Planning Interpretations To Confirm

Two interpretations are recorded in the validation contract scoping notes and
must be confirmed (to be formalized as phase decisions during execution):

- **PI-1 (`requirement.sync` parse contract).** Canonical source = PRD/TRD
  markdown requirement tables (first cell = ID `(PRD|TRD)-[A-Z0-9]+-[0-9]+[a-z]?`,
  second cell = summary). `type` from ID prefix (filterable by `--type`); `--id`
  syncs one ID; `status='active'`; `last_synced_hash` is a deterministic hash
  over the summary; `source_section` best-effort from nearest heading else null.
- **PI-2 (PRD-LINK-008 v0.1 enforcement).** v0.1's closed seven-relation set +
  issue/requirement-only endpoints structurally admit no dedicated-schema
  relation (hierarchy/escalation/resolution). App mirrors the relation trigger,
  rejects unknown relations + self-links. Decision/escalation dedicated-schema
  duplication activates with link endpoints in v0.2 (PRD-LINK-007).

## Carry-Forward Open At Kickoff

| ID | Item | Target |
|---|---|---|
| CF-008 | Resolution-reason field + cancel `--reason` | v0.2 |
| CF-011 | Broad `activity_log` reconciliation (link removal slice done here) | v0.1e |
| CF-013 | `complete.go` orchestration cleanup | v0.1e |
| CF-014 | Split `cmd/fid/issue.go` (984) + `replay.go` (485) under 400 lines | v0.1e |

PRD-LINK-006a requires link removal to project to `activity_log`; Phase 12 does
that per-command now. Link add does **not** project to `activity_log` (matches
the `issue create` precedent). CF-011's broad reconciliation stays v0.1e.

## Planning Reminders

- Schema/triggers exist; do not add a migration unless a spec-backed blocker is
  recorded first.
- Every new write command (`requirement sync`, `link add`, `link remove`)
  refreshes the projection with the held-lock freshness helper inside
  `journal.WithLock` (DEC-028); reads (`link list`) use the unlocked helper.
- One logical write = one event (sync with N requirements = one event with N
  effects; link add/remove = one event with two effects + any activity row).
- v0.1d deferrals still apply: requirement get/list + drift (v0.2, PRD-GOV-007);
  decision/escalation link endpoints (v0.2, PRD-LINK-007); label CRUD +
  lifecycle closure smoke (Phase 13).
