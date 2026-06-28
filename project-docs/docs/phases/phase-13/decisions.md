# Phase 13 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes, and
> wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-032: `fid label add` rejects duplicate names before append and stores optional metadata as given

- **Date:** 2026-06-27
- **Context:** Phase 13 Task 1 introduces the first label-creation surface.
  The validation contract's PI-1 fixes the v0.1 semantics that must hold
  without a schema migration or projector change: `labels.name` is the primary
  key, `colour` and `description` are optional, and label creation must not
  project to `activity_log`.
- **Decision:** `fid label add <name> [--colour <hex>] [--description <text>]
  --by <role>` validates and authorizes before append, refreshes the
  projection with the held-lock helper, and appends exactly one `label.added`
  event containing exactly one `labels` insert effect. A second add of the same
  label name is rejected before append with an actionable error naming the
  existing label. `--colour` and `--description` are optional and, when
  provided, are stored verbatim after CLI whitespace trimming; no additional
  colour-format validation is introduced in v0.1. Label creation projects no
  `activity_log` row.
- **Consequences:** The existing `labels` table and generic insert projector are
  sufficient, so Task 1 adds no schema migration and does not modify
  `internal/journal/project_effects.go`. `fid label list` can read the stored
  metadata directly from the `labels` table across json, toon, markdown, and
  yaml. No section-13 mismatch was found: the in-code `label.add` matrix row
  matches `03_TRD_SCHEMA.md` section 13 (`developer`, `tl`,
  `workflow_runner`, and `human` only).
- **Refs:** PI-1; phase-13/task-1; VAL-001; VAL-002; VAL-003; VAL-004;
  VAL-007; `03_TRD_SCHEMA.md` sections 9.6 and 13; `06_DESIGN_PLAN_CLI.md`
  section 9 `fid label add` / `fid label list`.

### DEC-033: `fid issue label add` / `remove` are single-row deltas with per-command `activity_log` projection

- **Date:** 2026-06-27
- **Context:** Phase 13 Task 2 implements PRD-ISSUE-023's additive label
  assignment surface on top of the existing `labels` and `issue_labels` tables.
  The validation contract's PI-2 fixes the v0.1 semantics that must hold
  without a schema migration or projector change: each command changes exactly
  one issue-label edge, validates before append, and projects the delta to
  `activity_log`. PI-3 also requires the new CLI surface to live in a new
  `cmd/fid/issue_label.go` file rather than growing `cmd/fid/issue.go`.
- **Decision:** `fid issue label add <issue-id> <label> --by <role>` requires
  the target issue to exist, the label to already exist in `labels`, and the
  issue not to already hold that label; `fid issue label remove <issue-id>
  <label> --by <role>` requires the issue to currently hold that label.
  Rejections happen before append with actionable errors and zero new journal
  lines. Each successful add/remove appends exactly one journal event:
  `issue.label.added` inserts one `issue_labels` row and `issue.label.removed`
  deletes one `issue_labels` row. Both also project exactly one `activity_log`
  row using an `ACT-...` id generated under the held lock, with `event`
  `issue_label_added` or `issue_label_removed` and `details` containing the
  changed `label` plus `direction` (`added` or `removed`).
- **Consequences:** Task 2 adds no schema migration and does not modify
  `internal/journal/project_effects.go`; the existing generic insert/delete
  projector is sufficient. The in-code permissions matrix now contains exactly
  one `issue.label.add` row and one `issue.label.remove` row, each matching
  `03_TRD_SCHEMA.md` section 13 (`developer`, `tl`, `workflow_runner`, and
  `human` only). The CLI surface lives in `cmd/fid/issue_label.go`, and
  `cmd/fid/issue.go` changes only by registering the new `label` subcommand.
- **Refs:** PI-2; PI-3; phase-13/task-2; VAL-001; VAL-005; VAL-006; VAL-007;
  VAL-008; `03_TRD_SCHEMA.md` sections 9.6 and 13; `06_DESIGN_PLAN_CLI.md`
  section 2 `fid issue label add` / `fid issue label remove`.

## Carry-Forward Review At Wrap (Tech Lead, 2026-06-28)

| ID | Disposition | Rationale |
|---|---|---|
| CF-008 | Re-defer to v0.2 | Resolution-reason field is a product/data-model enhancement; unrelated to labels/closure. |
| CF-011 | Broad reconciliation re-deferred to v0.1e; per-command projection delivered now | PRD-ISSUE-023 requires issue label add/remove to project to `activity_log`; Phase 13 delivered that per-command (DEC-033, VAL-007) exactly like Phase 12 did for link removal. `fid label add` does not project, matching the `issue create` precedent. The broad timeline policy still lands with `fid activity list` (v0.1e). |
| CF-013 | Re-defer to v0.1e | `complete.go` orchestration cleanup belongs with the `fid dev *` composites. |
| CF-014 | Re-defer to v0.1e | `cmd/fid/issue.go` / `replay.go` file-size split; Phase 13 honored PI-3 and did not grow `issue.go` (only registered the new `label` subcommand) — CF-014 not worsened. |
| **CF-015 (new)** | **Open; schedule v0.1e cleanup batch (with CF-013/CF-014)** | Phase 13 Review 1 + final review Improvement: `internal/labels/add.go` `Add` (~98 lines) and the `internal/issues/label.go` add/remove write flows (long, partially duplicated held-lock/prepare/dry-run/append/project flow) exceed the 60-line non-test function target. Not a blocker — behavior covered, gates green, no file over the 400-line target. Candidate for behavior-preserving function extraction (shared validation/options/event-construction helpers). |

## Wrap-Up Dispositions (Tech Lead, 2026-06-28)

### Validation contract outcome audit

All ten assertions reach a terminal state; none weakened.

- VAL-001..009: PASS in `qa/final-review.md` with cited test/CLI/integration
  evidence and a green phase-end gate. No disposition needed.
- VAL-010: was CARRY-FORWARD to the human checklist in final review. The
  generated `qa/human-checklist.md` (HV-13a..HV-13d) is now complete with
  overall **PASS** (completed 2026-06-28). Final outcome: **PASS**. No
  FAIL/BLOCKED/CARRY-FORWARD remains; no defer/cancel/escalation required.

### QA recommendation disposition

Final review had no Blocking or Major findings. Dispositions:

- **VAL-010 human dogfood** — resolved: the human checklist ran and passed
  overall (HV-13a..HV-13d all `pass`). No action remains.
- **Improvement: `internal/labels/add.go` `Add` over the 60-line target** —
  **defer** to v0.1e via **CF-015**. Behavior is covered and gates are green; a
  behavior-preserving extraction is the right home alongside the CF-013/CF-014
  cleanup batch.
- **Improvement: `internal/issues/label.go` add/remove long, partially
  duplicated write flows** — **defer** to v0.1e via **CF-015**. Same rationale;
  factor the shared held-lock/prepare/dry-run/append/project flow when the
  `fid dev *` composites touch the issue paths.

### DEC / deviation spec-impact walk

- DEC-032 (`fid label add` rejects duplicate names before append; optional
  metadata stored as given; no `activity_log` projection): spec-impacting? No.
  `03_TRD_SCHEMA.md` §§9.6/13 already expose the `labels` columns and the
  `label.add` §13 row; the in-code matrix matches (VAL-001/004 PASS), and
  `06_DESIGN_PLAN_CLI.md` §9 matches the implemented `fid label add` /
  `fid label list` (VAL-002/003 PASS) — `verified-no-change`.
- DEC-033 (`fid issue label add`/`remove` single-row deltas with per-command
  `activity_log` projection; new CLI file): spec-impacting? No. It uses the
  existing `issue_labels` table and generic insert/delete projector with no
  migration; in-code `issue.label.add`/`issue.label.remove` rows match §13
  (VAL-001/005..008 PASS); `06_DESIGN_PLAN_CLI.md` add/remove sections match
  (VAL-005/006 PASS) — `verified-no-change`.
- Developer deviations: none reported across Tasks 1-3.
- Validation-contract deviations: none; all ten VAL IDs were satisfied as
  written.
- Do-now QA recommendations: none (both improvements deferred via CF-015).

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | N/A | No charter-level scope/constraint change in a labels + lifecycle-closure slice. |
| `01_PRD.md` | verified-no-change | PRD-ISSUE-015 (label add/list) and PRD-ISSUE-023 (additive issue-label add/remove + activity projection) MET (final-review Requirement Compliance); PRD-OUT-001..004 four-format parity MET; PRD-NFR-008 one-event-per-write held (VAL-009). |
| `02_TRD.md` | verified-no-change | TRD-ARCH-010 matrix satisfied (VAL-001/004/008); held-lock write + rebuild parity held (VAL-009); no TRD change. |
| `03_TRD_SCHEMA.md` | verified-no-change | DEC-032/033; no migration added; §13 rows for `label.add`/`issue.label.add`/`issue.label.remove` match the in-code matrix (VAL-001 PASS); §9.6 `labels`/`issue_labels` tables sufficient as written (QA-12: diff over `migrations`/`project_effects.go` returned no paths). |
| `04_TRD_CONFIG.md` | N/A | No `.fid/` directory or config change. |
| `05_DESIGN_PLAN.md` | verified-no-change | Label fields and the create→claim→transition→complete→unclaim lifecycle surface consistent across reads (VAL-009/VAL-010 PASS); §3.1 CLI command tree includes the label commands (review-1 parity check). |
| `06_DESIGN_PLAN_CLI.md` | verified-no-change | `fid label add`/`list` (§9) and `fid issue label add`/`remove` flag, confirmation, and four-format shapes match implementation (VAL-002/003/005/006); Task 1-3 reports record no CLI spec change required. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | End-to-end v0.1d closure flow (create→claim→transition→complete→unclaim with link + label, rebuild-stable) matches documented behaviour (VAL-009 PASS). |
| `08_QUALITY_PLAN.md` | verified-no-change | Phase-end gate run as defined; COV-001 domain packages >=80% (`internal/labels` 81.1%, `internal/issues` 81.3%, `internal/policy` 84.6%) and COV-002 `cmd/fid` 78.1% met; no new gate added. |
| `09_AI_WORKFLOW.md` | N/A | No agent-workflow change. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 13 delivered exactly the v0.1d Phase 13 entry (label add/list + issue label add/remove + closure smoke); v0.1d now closed, no re-slice. |
| `GLOSSARY.md` | verified-no-change | "Label" / "Issue label" usage already covered; single-label delta vs replace-all is behaviour described in `06_DESIGN_PLAN_CLI.md`, not new canonical language. |

## Triage-Log Audit

- [x] Walked the Phase 13 triage record (`triage-log.md`).
- [x] Confirmed no RESUME verdict needs second-guessing in hindsight.
- [x] Flagged any TRG that should have escalated: none.
- [x] Confirmed autonomous TL-Lite resume remains in effect (DEC-002; no warm-up tally).

Audit note: Phase 13 recorded **no TL-Lite triage events**. The single
mid-phase review window (`qa/review-1.md`, Tasks 1-2) returned CONTINUE with
zero blocking findings and raised no escalation, so the Phase Runner resumed
without a triage decision; Task 3 was validated at final review. There are no
TRG verdicts to confirm or flag, and no RESUME that should have been a BLOCK.
`triage-log.md` records this for the audit trail.

## Workflow Notes

- **What worked well?** The kickoff planning interpretations (PI-1/PI-2/PI-3)
  mapped cleanly onto DEC-032/DEC-033 with no rework: duplicate/absent/already-
  present rejections, single-row deltas, per-command activity projection, and
  the `cmd/fid/issue_label.go` placement all landed as planned. The DEC-028
  held-lock freshness helper was reused from the first new write command, so —
  as in Phase 12 — no nested-lock surprise appeared. Sequencing the closure
  smoke last (Task 3) over already-shipped claim/transition/complete primitives
  made VAL-009 a pure integration assertion with no new behavior risk.
- **What was inefficient?** Minor: the function-size improvement notes recurred
  from Review 1 into the final review without a concrete carry-forward home;
  opening CF-015 at wrap closes that loop. No `triage-log.md` was maintained
  during the run, so the triage audit had to be reconstructed from `review-1.md`
  at wrap — cheap here (one CONTINUE review, no escalation) but a gap to keep an
  eye on for multi-review phases.
- **Were tasks the right size?** Yes. Three tasks each owned one coherent
  surface (label CRUD; per-issue label delta; closure smoke), each one session,
  with clean dependencies (1 → 2 → 3). No mid-phase do-now task was needed.
- **Did the QA review cadence feel right?** Yes. One mid-phase review over a
  three-task phase (a 2+1 window with Task 3 closed at final review) gave
  per-task signal on the two behavior-bearing tasks; the test-only closure smoke
  was appropriately validated at the phase gate. CONTINUE with improvement notes
  only.
- **Were there token-efficiency opportunities missed?** Minor. Task reports were
  thorough and cited VAL IDs and exact tests, which kept final QA cheap. The one
  small opportunity: pin recurring improvement notes to a carry-forward ID the
  first time they appear, so they do not get re-described each review round.

### Deep retrospective recommendation

No deep retrospective required. Phase 13 closed PASS with no Blocking/Major
findings, all ten VAL outcomes terminal, and v0.1d (plus the v0.1 headline
lifecycle exit criterion) closed end-to-end and rebuild-stable. One low-priority
carry-forward (CF-015) was opened for function-level cleanup, batched with the
existing v0.1e CF-013/CF-014 cleanup work. Carry into v0.1e planning: maintain
`triage-log.md` per review window even on CONTINUE, and pin recurring improvement
notes to a carry-forward ID on first appearance.
