# Phase Brief - Phase 13 (Planning)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Replaced at each phase kickoff by the incoming Tech Lead.
>
> **Size target:** under 150 lines.

## Phase Goal

Phase 13 closes v0.1d (and v0.1's headline lifecycle). It ships the label
surface — `fid label add` / `fid label list` (additive label CRUD, the first
label-creation surface) and `fid issue label add` / `fid issue label remove`
(single-label assignment that does **not** replace the full set) — and then
proves the v0.1 headline exit criterion with an end-to-end closure smoke over
the full gated lifecycle: create -> claim -> transition -> complete -> unclaim
with permission enforcement, links, and labels. Success means an agent can
create and list labels, add/remove one label on an issue without disturbing the
rest, and run the whole gated lifecycle end-to-end with rebuild-stable state.

## Status

**Planning.** No phase active. Branch will be `phase-13` off `main` (Phase 12
wrapped 2026-06-27, PR #13 merged). Next concrete action: QA Lead refines the
evidence strategy into `qa/qa-plan.md`, then the Phase Runner starts Task 1.
Labels/`issue_labels` tables and the section-13 matrix rows already exist, so
this phase expects CLI surface + in-code matrix rows only — **no schema
migration and no projector change** (the projector applies insert/delete effects
generically).

## Task Index

Full task specs and completion reports live under `../phases/phase-13/tasks/`.

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Label CRUD (`fid label add` / `list`) + `label.add` matrix row | Pending | None |
| 2 | Per-issue label assignment (`fid issue label add` / `remove`) + matrix rows + activity_log | Pending | Task 1 |
| 3 | v0.1d closure smoke (end-to-end lifecycle, integration test) | Pending | Tasks 1-2 |

## Validation Contract

Full contract: `../phases/phase-13/qa/validation-contract.md`.

| VAL ID | Source | Owning Task/Window | Validator | Status |
|---|---|---|---|---|
| VAL-001 | TRD-ARCH-010; schema §13 | tasks-1-2 | QA Lead | Draft |
| VAL-002 | PRD-ISSUE-015/023; design §9 | task-1 | QA Lead | Draft |
| VAL-003 | PRD-OUT-001..004; design §9 | task-1 | QA Lead | Draft |
| VAL-004 | TRD-ARCH-010; schema §13; CF-007 | task-1 | QA Lead | Draft |
| VAL-005 | PRD-ISSUE-023; PI-2 | task-2 | QA Lead | Draft |
| VAL-006 | PRD-ISSUE-023; PI-2 | task-2 | QA Lead | Draft |
| VAL-007 | PRD-ISSUE-023; CF-011; DEC-031 | task-2 | QA Lead | Draft |
| VAL-008 | TRD-ARCH-010; schema §13; CF-007 | task-2 | QA Lead | Draft |
| VAL-009 | Roadmap v0.1d Phase 13 exit; PRD-NFR-008; DEC-028 | task-3 | QA Lead | Draft |
| VAL-010 | `08_QUALITY_PLAN.md` §5.2/5.6; phase goal | phase-end | Human | Draft |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| QA plan | `../phases/phase-13/qa/qa-plan.md` | Missing |
| Mid-phase reviews | `../phases/phase-13/qa/review-M.md` | QA Lead sets cadence (3 tasks) |
| Final review | `../phases/phase-13/qa/final-review.md` | Pending |
| Human checklist | `../phases/phase-13/qa/human-checklist.md` | Pending |

## Active Constraints

- No schema migration unless implementation records a spec-backed blocker first;
  `labels`/`issue_labels` and the section-13 rows already exist.
- Every new write command refreshes the projection with the held-lock freshness
  helper inside `journal.WithLock` (DEC-028); one logical write = one event.
- Validate-before-append on every rejection path (CF-007): non-zero exit,
  actionable error, zero new journal lines.
- Append matrix rows per-implemented-action (DEC-026); copy `03_TRD_SCHEMA.md`
  §13 booleans verbatim, edit the spec only if implementation proves it wrong.
- New `fid issue label` subcommands go in `cmd/fid/issue_label.go`, **not** the
  984-line `cmd/fid/issue.go` (PI-3); do not worsen CF-014.

## Planning Interpretations (confirm during execution, record as DEC-NNN)

- **PI-1** (`fid label add`): duplicate name rejected before append;
  colour/description optional, stored as given (validation deferred, behaviour
  recorded in Task 1 report); label creation writes no `activity_log` row.
- **PI-2** (`fid issue label add`/`remove`): label must already exist; add of an
  already-present label and remove of an absent label reject before append; each
  success mutates one `issue_labels` row + one `activity_log` row (PRD-ISSUE-023,
  mirroring DEC-031); `fid issue update --labels` replace-all unchanged.
- **PI-3** (file placement): label subcommands in a new file; `issue.go` grows by
  registration only.

## Carry-Forward Disposition (this kickoff)

| ID | Disposition |
|---|---|
| CF-008 | Re-defer v0.2 (resolution-reason field). |
| CF-011 | Re-defer broad reconciliation to v0.1e; implement per-command `activity_log` for issue label add/remove now (PRD-ISSUE-023, VAL-007). |
| CF-013 | Re-defer v0.1e (`complete.go` cleanup). |
| CF-014 | Re-defer v0.1e; Phase 13 must not grow `issue.go` (PI-3). |

## Known Risks

- **Conflating add/remove with replace-all.** The single-label delta must not
  reuse the `--labels` replace-all path; trigger: VAL-005/006 assert other labels
  survive. Mitigation: distinct `issue_labels` insert/delete effect on one row.
- **Closure smoke flake.** TTL/time-dependent claim assertions could flake.
  Mitigation: assert journal-line counts and final state, not wall-clock timing.

## Open Questions

- [ ] None blocking. PI-1..PI-3 are recorded as planning interpretations to be
  confirmed as DEC-NNN during execution.

## Success / Exit Criteria

- [ ] All three tasks complete with green quality gates.
- [ ] VAL-001..009 PASS deterministically; VAL-010 PASS via human checklist.
- [ ] Final QA review has no unresolved blocking findings.
- [ ] v0.1d closed: full gated lifecycle runs end-to-end and is rebuild-stable.
- [ ] Phase sign-off decision recorded.

## Next Planning Checkpoint

v0.1d complete after Phase 13. Next is **v0.1e** (Telemetry, Composites,
Maintenance): `fid run record`, `fid stats`, the `fid dev *` composites,
`fid sql`, `fid journal inspect`, `fid maintenance check/vacuum`, activity log
read (`fid activity list` — the CF-011 broad reconciliation home), and the stub
table CLI guard. CF-013 and CF-014 batch with the `fid dev *` composites.
