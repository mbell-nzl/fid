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

### DEC-043: `dev gate` is part of CF-011 Task 3 activity reconciliation

- **Date:** 2026-07-05
- **Context:** Task 3 stopped because VAL-007 requires a full lifecycle with
  `dev gate`, but `dev.gate.recorded` currently appends only a `gate_runs` row
  and no `activity_log` row. This creates a silent timeline gap in the approved
  lifecycle proof. The Phase 16 roadmap requires CF-011 to backfill "any other
  write that does not yet emit a row" so the timeline is consistent across
  every write command, and `fid dev gate` appends an issue-scoped gate fact.
- **Decision:** Broaden Phase 16 Task 3 to include `fid dev gate` activity
  logging. A successful `dev.gate.recorded` event emits exactly one additional
  issue-scoped `activity_log` row named `gate_recorded` in the same journal
  event, with details sufficient to identify the gate run (`gate_run_id`),
  attempt number, command, exit code, and optional agent run. Rejected
  `dev gate` attempts still append zero journal/activity rows.
- **Consequences:** VAL-006 and VAL-007 remain valid; the Task 3 writer set now
  includes concern raise/address, unclaim, link add, and dev gate. This is an
  additive, rebuild-safe projection change with no event-schema bump. The
  blocker report at `tasks/task-3-done.md` is superseded by this decision; the
  developer should resume Task 3 from the partial implementation and produce a
  normal completion report after the full Task 3 gate passes.
- **Refs:** CF-011; PRD-OUT-017; PRD-DEV-004; PRD-DEV-010; `02_TRD.md`
  section 4.7; `03_TRD_SCHEMA.md` sections 7, 9, 14.2, 14.6;
  phase-16/tasks/task-3.md; VAL-006/007.

### DEC-044: Task 4 may remediate the pre-existing activity CLI fixture blocker

- **Date:** 2026-07-05
- **Context:** Task 4 implemented `fid sql` and produced focused VAL-008/009
  evidence, but the required global unit and integration gates stop in
  `cmd/fid/activity_test.go`: the test seeds a `2026-07-04` journal event after
  `fid init` has already created a journal file dated `2026-07-05`. The journal
  writer correctly rejects this as PRD-ERR-014 clock skew. `08_QUALITY_PLAN.md`
  TDD-006 also requires tests not to depend on system clock behavior.
- **Decision:** Treat the failure as a do-now deterministic fixture remediation
  in the Task 4 unblock window. The resumed Developer may edit only the
  activity CLI test fixture/helper needed to make seeded activity events
  monotonic with the workspace journal while preserving VAL-001..004 filter and
  format assertions. This does not change `fid sql` scope, product behavior, or
  the required Task 4 gates.
- **Consequences:** Do not weaken PRD-ERR-014, do not bypass the journal writer
  in production code, and do not narrow the required global gates. After the
  fixture remediation, rerun the exact Task 4 command list, replace the Task 4
  blocker report with a normal completion report, and record the remediation as
  a spec-backed deviation from the original Task 4 file list.
- **Refs:** PRD-ERR-014; PRD-OUT-017; `08_QUALITY_PLAN.md` TDD-006;
  phase-16/tasks/task-1.md VAL-001..004; phase-16/tasks/task-4.md VAL-008/009;
  `cmd/fid/activity_test.go`; `project-docs/docs/phases/phase-16/tasks/task-4-done.md`.

### DEC-045: Phase 16 must repair the missed mid-phase QA cadence before continuing

- **Date:** 2026-07-05
- **Context:** The Phase 16 QA plan set a stricter cadence than the generic
  every-third-task rule: Review 1 after Tasks 1-2 and Review 2 after Tasks 3-4.
  The run ledger records Tasks 1-3 complete and no REVIEW entry; the QA
  directory contains no `review-1.md` or `review-2.md`. The workflow has
  therefore not been followed even though Task 3 is complete.
- **Decision:** Before any further Developer work, the Phase Runner must spawn
  `qa_reviewer` for the overdue Review 1 covering Tasks 1-2 and append a REVIEW
  entry to the run ledger. If Review 1 returns Blocking findings or opens an
  escalation, the normal TL-Lite/STOP rules apply. If Review 1 clears or
  resumes, the runner may continue Task 4 from the DEC-044 blocker point. After
  Task 4 is finalized, the runner must spawn Review 2 covering Tasks 3-4 before
  Task 5 starts.
- **Consequences:** This is a workflow correction, not a validation-contract
  change. Do not combine Review 1 and Review 2 into one artifact; Review 1 stays
  scoped to the planned Tasks 1-2 window, and Review 2 stays scoped to Tasks
  3-4 after Task 4 has a normal completion report.
- **Refs:** `project-docs/docs/phases/phase-16/qa/qa-plan.md` sections 3, 5,
  and 8; `08_QUALITY_PLAN.md` §5.2; `09_AI_WORKFLOW.md` §4.3;
  `project-docs/docs/phases/phase-16/run-ledger.md`.

### DEC-046: Task 5 may remediate the internal activity fixture blocker

- **Date:** 2026-07-05
- **Context:** Task 5 implemented `fid journal inspect` and produced focused
  VAL-010/011 evidence, but the required repo-wide gate stops in
  `internal/activity/activity_test.go` (`TestListFiltersByOnSinceAndEvent`).
  The fixture seeds a `2026-07-04` journal event after workspace initialization
  has created a later `2026-07-05` journal file, so the PRD-ERR-014
  out-of-order guard correctly rejects the append. This is the same
  deterministic fixture family as DEC-044, now in the domain package test.
- **Decision:** Treat the failure as a do-now deterministic fixture remediation
  in the Task 5 unblock window. The resumed Developer may edit only
  `internal/activity/activity_test.go` and directly local test helper code
  needed to derive seeded activity event IDs/timestamps from the workspace's
  latest journal date, preserving VAL-001..004 behavior. This does not change
  `fid journal inspect` scope, production journal ordering, product behavior, or
  the required Task 5 gates.
- **Consequences:** Do not weaken PRD-ERR-014, do not bypass production journal
  append paths, and do not narrow the required repo-wide gate. After the fixture
  remediation, rerun the exact Task 5 command list, replace the Task 5 blocker
  report with a normal completion report, and record the remediation as a
  DEC-046-backed deviation from the original Task 5 file list.
- **Refs:** PRD-ERR-014; PRD-OUT-017; PRD-OPS-013; `08_QUALITY_PLAN.md`
  TDD-006; phase-16/tasks/task-1.md VAL-001..004; phase-16/tasks/task-5.md
  VAL-010/011; `internal/activity/activity_test.go`;
  `project-docs/docs/phases/phase-16/tasks/task-5-done.md`.

### DEC-047: Blocking task outcomes require TL-Lite triage before STOP

- **Date:** 2026-07-05
- **Context:** The Phase 16 run ledger records STOP entries for task blockers,
  but no `triage-log.md` exists and no TRIAGE entries were appended. The active
  Phase Runner prompt currently says to stop on blocker reports or quality-gate
  failures before invoking TL-Lite, while TL-Lite triage is only described for
  soft blockers, concerns, review findings, or new open escalations. This
  regressed the intended first-step triage path for blocking task outcomes.
- **Decision:** Restore TL-Lite as the first classifier for task-level blocking
  outcomes. When a Developer returns a blocker report, a required quality gate
  fails, report-lint fails, required command evidence is missing/skipped/narrower
  than specified, or an ad hoc gate substitution is detected, Phase Runner must
  spawn `tl_lite` before appending a final STOP. TL-Lite records a TRG entry and
  returns `VERDICT: RESUME` or `VERDICT: BLOCK`. `RESUME` is allowed only when
  the issue is misclassified, already resolved by verified evidence, or covered
  by an accepted recorded decision. `BLOCK` always stops the run and the STOP
  entry must reference the TRG/ESC.
- **Consequences:** The hard-stop rules still stand: the runner may not continue
  past an untriaged failing gate or a `BLOCK` verdict. Existing STOP entries are
  left append-only; the next resume must first triage the Task 5 STOP under
  DEC-047, with DEC-046 available as the accepted recorded decision that can
  justify a `RESUME`.
- **Refs:** DEC-002; `project-docs/docs/phases/phase-16/run-ledger.md`;
  `project-docs/templates/phase-runner-prompt.md`;
  `.codex/agents/phase_runner.toml`; `.claude/agents/phase-runner.md`.

### DEC-048: Phase 16 wraps with all validation assertions terminal

- **Date:** 2026-07-05
- **Context:** Final review returned PASS WITH IMPROVEMENTS with VAL-001..012
  PASS and VAL-013 carried forward to the human checklist. The human checklist
  now records HUMAN-16a..d PASS, summary totals 4/4 pass, and overall result
  PASS on 2026-07-05.
- **Decision:** Archive Phase 16 with VAL-001..013 terminal PASS. Close
  **CF-011** and **CF-016** because their deterministic validation assertions
  passed. Defer the non-blocking file-size findings as **CF-018**
  (`internal/sqlquery/sqlquery.go`) and **CF-019**
  (`internal/format/markdown.go`).
- **Consequences:** Phase 17 does not need a Phase 16 human-dogfood
  carry-forward item. Keep the retrospective lesson: a run sheet is not enough;
  the checklist outcomes themselves must be filled before wrap closes.
- **Refs:** `qa/final-review.md`; `qa/human-checklist.md`;
  `qa/human-checklist-run-sheet.md`; VAL-013; `08_QUALITY_PLAN.md` §5.3;
  `09_AI_WORKFLOW.md` §4.3.

## Carry-Forward Review At Kickoff

| ID | Disposition | Rationale |
|---|---|---|
| CF-008 | Re-defer to v0.2 | Resolution-reason field is a product/data-model enhancement unrelated to this read-surface phase. |
| CF-011 | Schedule in Phase 16 | Addressed by Tasks 2-3 after DEC-041 and DEC-043. Remains open until VAL-005/006/007 pass. |
| CF-016 | Schedule in Phase 16 | Addressed by Task 6 after DEC-042. Remains open until VAL-012 passes. |
| CF-017 | Re-defer to v0.2c | Behavior-preserving `internal/dev` function-size cleanup belongs with the next `internal/dev` work. |

## Phase-End Validation Audit

| VAL ID | Final outcome | Tech Lead disposition |
|---|---|---|
| VAL-001..004 | PASS | Accepted. Task 1, Review 1, and final phase gate cover `fid activity list` shape, filters, ordering, and formats. |
| VAL-005 | PASS | Accepted. Task 2 and Review 1 cover create/ac add/ac verify/complete activity rows and rejection zero-append behavior. |
| VAL-006..007 | PASS | Accepted. Task 3, Review 2, and final lifecycle rerun cover remaining CF-011 writers, `gate_recorded`, and rebuild parity. |
| VAL-008..009 | PASS | Accepted. Task 4, Review 2, and final SQL reruns cover read forms, mutation rejection, and the read-only connection backstop. |
| VAL-010..011 | PASS | Accepted. Task 5 and final journal-focused reruns cover inspect fields, `HashEvent` parity, not-found behavior, and revert verdicts. |
| VAL-012 | PASS | Accepted. Task 6 and final reruns cover raw `gate_runs` history, empty state, formats, and separation from completion snapshots. |
| VAL-013 | PASS | Accepted. `human-checklist.md` records HUMAN-16a..d PASS, 4/4 total pass, and overall result PASS on 2026-07-05. |

## QA Recommendation Disposition

| Source | Recommendation | Disposition | Rationale |
|---|---|---|---|
| Review 1 | Fix `cmd/fid/activity_test.go` monotonic-date fixture issue. | Do now, complete | DEC-044 authorized the narrow fixture remediation; Task 4 report and Review 2 record the rerun PASS. |
| Review 1 | Future CF-011 reports should enumerate each rejection fixture explicitly. | Cancel as no-op | Final review confirmed VAL-005 evidence and report-lint passed; keep as a reviewer lesson, not carry-forward debt. |
| Review 2 / Final review | Split `internal/sqlquery/sqlquery.go` (595 lines) when reopened. | Defer as CF-018 | Non-blocking quality-budget follow-up; behavior and gates passed. |
| Final review | Split `internal/format/markdown.go` (661 lines) when formatter work reopens it. | Defer as CF-019 | Non-blocking quality-budget follow-up; behavior and gates passed. |
| Final review | Record human dogfood outcomes for VAL-013. | Do now, complete | `human-checklist.md` now records all four HUMAN-16 rows PASS and overall PASS. |

## Carry-Forward Review At Wrap

| ID | Disposition | Rationale |
|---|---|---|
| CF-008 | Re-defer to v0.2 | Resolution-reason field remains a product/data-model enhancement outside Phase 16. |
| CF-011 | Closed | VAL-005/006/007 PASS; final review marks the broad `activity_log` reconciliation MET. |
| CF-016 | Closed | VAL-012 PASS; `fid issue get --section gates` shipped with raw `gate_runs` evidence. |
| CF-017 | Re-defer to v0.2c | `internal/dev` `Finish`/`Block` function-size cleanup remains tied to the next `internal/dev` reopen. |
| CF-018 | Open, target next SQL package reopen | Behavior-preserving split of `internal/sqlquery/sqlquery.go` under the 400-line file target. |
| CF-019 | Open, target next formatter-heavy task | Behavior-preserving split of `internal/format/markdown.go` under the 400-line file target. |

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | Phase 16 shipped planned local-first CLI read surfaces and did not alter vision, scope, or constraints. |
| `01_PRD.md` | verified-no-change | Existing PRD-OUT-017, PRD-OUT-008, PRD-OPS-013, PRD-DEV-004/010, and PRD-ERR-014 covered the delivered behavior. |
| `02_TRD.md` | touched: see DEC-048 | Added `internal/activity/` and `internal/sqlquery/` package boundaries; existing journal inspect and gate-run sections remained consistent. |
| `03_TRD_SCHEMA.md` | touched: see DEC-041 / DEC-043 / DEC-048 | Updated `activity_log.event` vocabulary for the CF-011 backfill, including `gate_recorded`; no migration or event-schema bump. |
| `04_TRD_CONFIG.md` | verified-no-change | Phase 16 did not change `.fid/` layout, config files, pricing config, or config loading. |
| `05_DESIGN_PLAN.md` | verified-no-change | Command-tree and SQL escape-hatch interaction pattern already matched the delivered read surfaces. |
| `06_DESIGN_PLAN_CLI.md` | touched: see DEC-042 / DEC-048 | Updated `fid journal inspect` verdict buckets and `fid issue get --section gates` contract. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | No end-to-end flow contract changed; Phase 16 added read surfaces and projection coverage. |
| `08_QUALITY_PLAN.md` | verified-no-change | TDD-006, phase gates, validation-outcome rules, and human checklist rules already covered the fixture and VAL-013 outcomes. |
| `09_AI_WORKFLOW.md` | touched: see DEC-047 | Workflow correction restored TL-Lite triage before final STOP/resume for task-level blockers and report-lint/gate failures. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 16 implementation matched the roadmap item; Phase 17 remains the v0.1 closure phase. |
| `GLOSSARY.md` | verified-no-change | No new ambiguous project language surfaced; existing read-model, journal, validation, and SQL terms were sufficient. |

## Triage-Log Audit

- Walked every TRG-NNN in `triage-log.md`: **complete**.
- TRG-046 RESUME was correct in hindsight. DEC-046 already authorized the
  narrow `internal/activity/activity_test.go` fixture remediation, and Task 5
  later recorded the full required gate PASS.
- No TRG should have escalated. The blocker was a deterministic test-fixture
  monotonicity issue, not a product, persistence, or public-interface decision.
- Autonomous TL-Lite resume remains in effect per DEC-002; no warm-up tally is
  maintained.

## Workflow Notes

Answer in 1-3 sentences each at wrap.

- **What worked well?** The validation contract kept the large read-surface
  phase testable: each public command had focused VAL coverage, and final review
  could rerun narrow representative tests before the full gate. DEC-043 also
  showed the contract doing its job by exposing the missing `dev gate` timeline
  row before the phase closed.
- **What was inefficient?** The fixed-date activity fixtures caused two
  cross-task blockers after they collided with PRD-ERR-014. The fixture pattern
  should have been audited once when the first monotonicity failure appeared.
- **Were tasks the right size?** Mostly yes. Task 3 was near the upper bound
  because it combined several writer backfills with lifecycle rebuild parity,
  but the DEC-043 split point kept the resume narrow.
- **Did the QA review cadence feel right?** The cadence was right after DEC-045
  repaired the missed Review 1 checkpoint. Reviews after Tasks 2 and 4 caught
  the highest-risk windows without delaying the lower-risk read tasks.
- **Were there token-efficiency opportunities missed?** Yes. Current docs stayed
  stale after Task 5/6 completion, so later agents had to reconcile the archive
  against default-loaded context. Phase Runner should refresh `current/` as soon
  as a STOP is resolved, not only at wrap.

### Deep Retrospective Recommendation

Make human checklist completion a ledger-gated phase-end event. A run sheet is
useful, but the phase should not reach `RUN-END` or Tech Lead wrap without either
recorded human outcomes in `human-checklist.md` or an explicit DEC/CF
disposition that says the human-only VAL remains carried forward.
