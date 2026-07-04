# Phase 16 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes, and
> wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-040: Phase 16 runs in Heavy Mode

- **Date:** 2026-07-05
- **Context:** `09_AI_WORKFLOW.md` section 4.3 says heavy mode applies when a
  phase includes public-interface changes or a security/privacy surface. Phase
  16 adds four public CLI read surfaces and includes TRD-SEC-002 for the
  read-only `fid sql` escape hatch.
- **Decision:** Phase 16 runs in **Heavy Mode**. Required artifacts are the
  validation contract, QA plan with fresh command/feasibility validation,
  mid-phase review cadence, final review, human checklist, and full-scope spec
  sync at wrap.
- **Consequences:** The runner must not treat Phase 16 as a lightweight read
  phase. `fid sql` security evidence, integration gates, and human dogfood are
  required before wrap.
- **Refs:** `09_AI_WORKFLOW.md` section 4.3; PRD-OUT-008; TRD-SEC-002;
  phase-16/qa/validation-contract.md VAL-008/009/013.

### DEC-041: CF-011 activity timeline scope and vocabulary are approved

- **Date:** 2026-07-05
- **Context:** CF-011 requires consistent population of the human-facing
  `activity_log` projection. The schema already constrains valid
  `activity_log.ref_type` values, so the remaining question is which writers
  should add issue-scoped activity rows and what event names they use.
- **Decision:** Every write command that mutates or appends an issue fact, or an
  issue-to-requirement link, emits exactly one `activity_log` row inside its
  single journal event. Backfilled writers are `issue create`, `issue ac add`,
  `issue ac verify`, `issue complete`, `concern raise`, `concern address`,
  `unclaim`, and `link add`. Event names are `issue_created`, `ac_added`,
  `ac_verified`, `task_completed`, `concern_raised`, `concern_addressed`,
  `claim_released`, and `link_added`.
- **Consequences:** `fid label add` remains excluded because label creation is
  not a valid activity reference; issue-label add/remove already log against
  the issue. `fid run record` remains telemetry owned by stats. `fid requirement
  sync` remains excluded as noisy bulk pointer sync; requirement reads/drift are
  v0.2. Existing tests that expected zero activity rows for the backfilled
  writers are updated as a decision-backed behavior change, not weakened.
- **Refs:** CF-011; PRD-OUT-017; `03_TRD_SCHEMA.md` sections 9, 14.2, 14.6;
  phase-16/tasks/task-2.md; phase-16/tasks/task-3.md; VAL-005/006/007.

### DEC-042: CF-016 gate history read surface ships as `issue get --section gates`

- **Date:** 2026-07-05
- **Context:** Phase 15 added append-only `gate_runs`, but raw gate history has
  no first-class issue-scoped read surface. After `dev finish` or `dev block`,
  gate data is visible only indirectly through the completion snapshot.
- **Decision:** Phase 16 includes CF-016 as Task 6: `fid issue get <id>
  --section gates` reads the raw `gate_runs` history for the issue. This is
  preferred over a new `fid dev gate list` verb because it fits the existing
  issue section model and keeps raw history distinct from
  `task_completions.gate_results`.
- **Consequences:** Task 6 and VAL-012 are in scope. The command is read-only:
  no `--by` flag, no permissions-matrix row, and no schema change.
- **Refs:** CF-016; PRD-DEV-004; PRD-DEV-010; `03_TRD_SCHEMA.md` section 7;
  phase-16/tasks/task-6.md; VAL-012.

## Carry-Forward Review At Kickoff

| ID | Disposition | Rationale |
|---|---|---|
| CF-008 | Re-defer to v0.2 | Resolution-reason field is a product/data-model enhancement unrelated to this read-surface phase. |
| CF-011 | Schedule in Phase 16 | Addressed by Tasks 2-3 after DEC-041. Remains open until VAL-005/006/007 pass. |
| CF-016 | Schedule in Phase 16 | Addressed by Task 6 after DEC-042. Remains open until VAL-012 passes. |
| CF-017 | Re-defer to v0.2c | Behavior-preserving `internal/dev` function-size cleanup belongs with the next `internal/dev` work. |

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending | phase wrap |
| `01_PRD.md` | pending | phase wrap |
| `02_TRD.md` | pending | phase wrap |
| `03_TRD_SCHEMA.md` | pending | phase wrap |
| `04_TRD_CONFIG.md` | pending | phase wrap |
| `05_DESIGN_PLAN.md` | pending | phase wrap |
| `06_DESIGN_PLAN_CLI.md` | pending | phase wrap |
| `07_DESIGN_PLAN_FLOWS.md` | pending | phase wrap |
| `08_QUALITY_PLAN.md` | pending | phase wrap |
| `09_AI_WORKFLOW.md` | pending | phase wrap |
| `10_IMPLEMENTATION_ROADMAP.md` | pending | phase wrap |
| `GLOSSARY.md` | pending | phase wrap |

## Triage-Log Audit

- [ ] Walked every TRG-NNN in `triage-log.md`.
- [ ] Confirmed each RESUME verdict was correct in hindsight.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed autonomous TL-Lite resume remains in effect (DEC-002; no warm-up tally).

## Workflow Notes

Answer in 1-3 sentences each at wrap.

- **What worked well?** pending
- **What was inefficient?** pending
- **Were tasks the right size?** pending
- **Did the QA review cadence feel right?** pending
- **Were there token-efficiency opportunities missed?** pending

### Deep Retrospective Recommendation

Pending phase wrap.
