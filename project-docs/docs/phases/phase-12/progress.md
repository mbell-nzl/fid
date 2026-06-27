# Phase 12 Progress

> Phase: 12 - Requirements + links (roadmap v0.1d, slice 3 of 4)
> Status: COMPLETE - PASS. All 4 tasks done; final review PASS; human checklist
> PASS; wrapped 2026-06-27.

## Plan

Phase 12 delivers requirement pointers and the typed-link surface:

- `fid requirement sync` (minimal): parse canonical PRD/TRD IDs from a source
  doc, upsert pointer rows with `last_synced_hash`, report created/updated/
  unchanged counts. No requirement get/list or drift (v0.2).
- `fid link add`: typed links with auto-inverse projection, relation/endpoint
  compatibility validation, duplicate and self-link rejection.
- `fid link remove`: declared + inverse delete in one event with `activity_log`
  removal evidence.
- `fid link list`: outgoing and incoming, all four formats, resolved titles.
- `requirement.sync`, `link.add`, `link.remove` rows added to the Phase 10
  in-code permissions matrix.

No schema migration; the `links`/`requirements` tables and triggers already
exist (v0.1a, `migrations/0002`/`0005`).

## Task Index

| Task | Title | Status |
|---|---|---|
| 1 | Requirement pointers + `fid requirement sync` + matrix row | Complete; Review 1 CONTINUE |
| 2 | Links domain + matrix rows + `fid link add` | Complete; Review 1 CONTINUE |
| 3 | `fid link remove` | Complete; Review 2 CONTINUE |
| 4 | `fid link list` + end-to-end flow | Complete; Review 2 CONTINUE |

## Validation

Validation contract at `qa/validation-contract.md` includes VAL-001..012.
VAL-012 is the phase-end human dogfood of sync/add/list/remove legibility. QA
plan at `qa/qa-plan.md`; two mid-phase reviews (`qa/review-1.md`,
`qa/review-2.md`) over a 2+2 task window.

## Progress History

- Kickoff: drafted VAL-001..012 and Tasks 1-4; recorded PI-1/PI-2 to confirm;
  re-deferred CF-008/011/013/014.
- Task 1 complete: `internal/requirements`, `fid requirement sync`, and the
  `requirement.sync` matrix row (DEC-029).
- Task 2 complete: `internal/links` closed relation/inverse registry,
  `fid link add`, and the `link.add`/`link.remove` matrix rows (DEC-030).
- Review 1 (Tasks 1-2): CONTINUE; two report-quality improvement notes
  (TRG-001 RESUME).
- Task 3 complete: `fid link remove` (declared + inverse in one event +
  `link_removed` activity evidence; nonexistent-link rejection).
- Task 4 complete: `fid link list` (four-format parity, resolved titles) +
  `TestV01DLinkFlowAndRebuildParity` integration test.
- Review 2 (Tasks 3-4): CONTINUE; two report-quality notes plus the
  PRD-LINK-006a one-vs-two activity-row wording observation
  (TRG-002 RESUME, accepted under the validation contract).
- Final review: PASS WITH CARRY-FORWARD. VAL-001..011 PASS; VAL-012
  CARRY-FORWARD to the human checklist.
- Human checklist (HV-12a..HV-12d): complete, overall PASS (2026-06-27) ->
  VAL-012 final outcome PASS.
- Wrap (2026-06-27): VAL-012 dispositioned PASS; DEC-031 closes the
  PRD-LINK-006a wording observation; spec-sync verdicts recorded
  (all verified-no-change or N/A); TRG-001/TRG-002 audited (both RESUME
  correct); no new carry-forward opened; workflow notes captured. Phase
  COMPLETE.

## Carry-Forward

- CF-008 re-deferred to v0.2 (resolution-reason field).
- CF-011 broad reconciliation re-deferred to v0.1e; Phase 12 link removal now
  projects to `activity_log` per-command (DEC-031). Link add does not project,
  matching the `issue create` precedent.
- CF-013 re-deferred to v0.1e (`complete.go` orchestration cleanup).
- CF-014 re-deferred to v0.1e (`cmd/fid/issue.go` / `replay.go` file-size
  split).

## Escalations

- None raised in Phase 12. No open escalations at wrap. See
  `tl-escalations.md`.
