# Phase 5 Decisions

> Purpose: phase-local decision history, spec sync verdicts, triage audit,
> workflow notes, and wrap-up decisions. DEC numbering is global and never
> resets (Phase 4 ended at DEC-014).

## Decisions

### DEC-015: `todo → done` rejection is correct; "close without completing" is `cancelled`; resolution-reason deferred (CF-008)

- **Date:** 2026-06-14
- **Source:** Phase 5 human checklist `HUMAN-013c` (`qa/human-checklist.md`) —
  human asked, after a correctly-rejected `todo → done` transition, whether the
  spec supports marking a task done/closed without completing it (e.g. a
  low-priority backlog task made redundant by a later refactor), à la a Jira
  `resolution=wont_complete`.
- **Analysis:** The rejection is **spec-conformant, not a defect.** In the
  `task-default` workflow (TRD §4.6) `done` has status category `completed`;
  there is intentionally no `todo → done` edge. The design already provides a
  "close without completing" terminal: `cancelled`, reachable from any
  non-cancelled state (PRD-WF-013), with its own `cancelled` status category
  distinct from `completed`. The epic-wrap rule
  `all_child_tasks_done_or_cancelled` confirms intent: **done OR cancelled** are
  the two legitimate terminal outcomes for a child. A redundant backlog task
  should transition `todo → cancelled`.
- **Genuine gap (separate concern):** fid has no *resolution type / reason*
  field, so `cancelled` is a single undifferentiated bucket and cannot
  distinguish "redundant", "duplicate", "won't fix", etc. That is a real product
  enhancement but is a **new field + CLI flag (e.g. cancel `--reason`), not a new
  workflow edge**, and it touches the data model and the `activity_log` audit.
- **Decision:** No spec change in this wrap. Phase 5 stays PASS. Capture the
  resolution-reason enhancement as carry-forward **CF-008** for a product-owner
  decision in the v0.2 era (alongside configurable workflows PRD-WF-011/012). Do
  not widen the v0.1 default graphs. Recommend the human dogfood `todo →
  cancelled` to confirm it serves the redundant-backlog case today.
- **Refs:** `qa/human-checklist.md` HUMAN-013c; TRD §4.6 (`task-default`,
  `epic-default`); PRD-WF-006/010/013; CF-008 (`decisions_active.md`).

### DEC-016: `fid issue update` multi-field output formatting — deferred CLI polish (CF-009)

- **Date:** 2026-06-14
- **Source:** Phase 5 human checklist `HUMAN-013a` / `HUMAN-013f`
  (`qa/human-checklist.md`). All update fields worked (PASS), but the confirmation
  line renders all changed fields as one long comma-separated string; the human
  asked for one field per line for readability.
- **Decision:** Accept as a non-blocking presentation improvement. Defer, not
  cancel; fold into the next `fid issue update`-adjacent change rather than
  spawning a standalone task. Capture as carry-forward **CF-009**. The behaviour
  and audit content are correct (VAL-001/002/003 PASS); this is output rendering
  only and touches `06_DESIGN_PLAN_CLI.md` plus the `cmd/fid` update presenter
  when actioned.
- **Refs:** `qa/human-checklist.md` HUMAN-013a/HUMAN-013f; CF-009
  (`decisions_active.md`).

## QA Recommendation Disposition

QA final review (`qa/final-review.md`, VERDICT: **PASS WITH IMPROVEMENTS**)
recorded **0 blocking / 0 major / 1 improvement**.

- **Improvement — non-test file size:** `internal/issues/transition.go` (426),
  `internal/issues/update.go` (551), and `internal/journal/project.go` (524) are
  above the 400-line target in `08_QUALITY_PLAN.md` §4.4; review-4 recommends
  splitting event assembly from append/project orchestration. **Disposition:
  DEFER.** Not a behaviour or contract issue; fold the split into the next
  transition/update/projection-adjacent change rather than refactoring now (Tech
  Lead writes no product code, and the next write-command phases touch these
  files). Captured as carry-forward **CF-010**.

No QA blocking or major finding remains. The two human-raised items are
dispositioned in DEC-015 (no change) and DEC-016 (defer / CF-009).

## Validation Contract Audit

- **VAL-001..VAL-012: PASS** with cited deterministic evidence
  (`qa/final-review.md` Validation Contract Outcomes; reviews 1–4 plus
  review-1-addendum; Task 1/1b/2/3/4/5 reports).
- **VAL-013: PASS.** At final review VAL-013 was `CARRY-FORWARD` (human-only
  judgment not executable by QA Lead). The human dogfood checklist has since
  recorded **6/6 PASS** (`qa/human-checklist.md`, `HUMAN-013a`..`HUMAN-013f`).
  No VAL outcome remains FAIL, BLOCKED, or CARRY-FORWARD.
- **Both-state coverage** held for the binary/conditional assertions: VAL-001
  (valid update / invalid unchanged), VAL-004 (one row inserted / unresolved ref
  rejected / generic writer rejected), VAL-005 (rejected zero-append + recovery /
  accepted +1 line), VAL-007 (allowed / disallowed naming both states),
  VAL-008/009/010 (rule unmet names IDs / satisfied transitions), VAL-011 (cancel
  from non-cancelled / rejected from cancelled), VAL-012 (open-concern unsatisfied
  / no-open satisfied).
- **Carry-forward closures proven by contract:** CF-005 (`activity` op + sole-
  writer discipline) closed via VAL-004 (Task 2) and the Task 3/5 emitters;
  CF-007 (rejected-write journal poisoning) closed via VAL-005 (Task 1b /
  review-1-addendum) and reinforced by VAL-007.

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not a default skip.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | Phase 5 delivered the in-scope v0.1b mutation/transition slice; no identity/scope/constraint change. |
| `01_PRD.md` | verified-no-change | PRD-ISSUE-015/021 and PRD-WF-001..010/013 MET as written (final review Requirement Compliance). DEC-015 confirms `cancelled` already satisfies "close without completing" (PRD-WF-013); resolution-reason is a future enhancement (CF-008), not a v0.1 requirement reword. |
| `02_TRD.md` | verified-no-change | TRD §4.3 validate-before-append seam realized by Task 1/1b (CF-007 closed); §4.6 default workflow graphs/rules implemented exactly as specified (VAL-006..012). No dependency change this phase. |
| `03_TRD_SCHEMA.md` | verified-no-change | §4.3 `activity` op + §14.2 polymorphic-ref backstop + §14.6 sole-writer discipline implemented app-layer per DEC-006 (SQL session guard stays deferred); no schema/trigger/view reword (VAL-004). |
| `04_TRD_CONFIG.md` | N/A | Phase 5 added no `.fid/` layout or config surface. |
| `05_DESIGN_PLAN.md` | verified-no-change | Mutation and workflow-transition flows realized as described; `internal/workflow` hard-coded graphs match the documented design. |
| `06_DESIGN_PLAN_CLI.md` | updated | Task 3 documented `fid issue update`; Task 5 documented `fid issue transition` (final review Documentation Review). DEC-016/CF-009 records a deferred update-output formatting refinement against this doc; not actioned in-wrap. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | update and transition follow the documented journal→projection→read flow with validate-before-append commit point; rejected writes leave the journal clean (VAL-005/007). |
| `08_QUALITY_PLAN.md` | verified-no-change | COV-001 (`internal/issues` 82.8%, `internal/journal` 80.3%, `internal/policy` 85.0%, new `internal/workflow` 95.5%) and COV-002 (`cmd/fid` 79.0%) held; DEC-007 carry-in applied. §4.4 file-size improvement recorded as CF-010, not a plan reword. |
| `09_AI_WORKFLOW.md` | verified-no-change | Workflow ran as specified (contract → QA plan → dev/QA loop with four mid-phase reviews + addendum → final review → human checklist). The one mid-phase BLOCK (ESC-003) routed and cleared correctly. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 5 delivered the v0.1b "mutation + workflow transitions" cut per §v0.1b slicing note; Phase 6 (ready queue + v0.1b closure) cut unchanged. |

## Triage-Log Audit

- **TRG-001** (Task 1 VAL-005 rebuild-recovery evidence gap) — verdict **BLOCK**.
  Confirmed correct in hindsight: VAL-005 is a hard contract assertion and Tasks
  3/5 build on the seam, so a missing rebuild-after-reject proof could not RESUME.
  Routed to **ESC-003**, dispositioned DO-NOW (test-and-report-only Task 1b), and
  cleared by Task 1b + `review-1-addendum.md`. Appropriate handling; no duplicate
  escalation.
- No TRG should have escalated that did not, and no RESUME verdicts were recorded
  this phase. Autonomous TL-Lite resume remains permanently ENABLED (DEC-002); no
  warm-up tally to maintain.

## ESC-003 Disposition

**Resolved and archived.** ESC-003 (Task 1 VAL-005 rebuild-recovery evidence gap)
was dispositioned DO-NOW on 2026-06-13 (test-and-report-only remediation in
`tasks/task-1b.md`, no contract weakening). Task 1b and `review-1-addendum.md`
prove rejected-create zero append, later read success, **rebuild recovery after
the rejected create**, accepted-create +1 line, and dry-run rollback — the stated
closure condition. VAL-005 is PASS. Moved to
`project-docs/docs/phases/phase-5/tl-escalations.md`.

## Workflow Notes

Answer in 1-3 sentences each. Completed at phase wrap-up.

- **What worked well?** The Phase 4 deep-retro lesson paid off immediately:
  baking rejected-write journal hygiene into VAL-005/VAL-007 (line-count-unchanged
  + read-recovery) made CF-007 a first-class, provable closure rather than a
  latent defect, and the human dogfood again earned its keep by surfacing a real
  product question (DEC-015) that automated gates could never raise.
- **What was inefficient?** The Task 1 VAL-005 evidence gap forced an ESC-003 →
  Task 1b → review-1-addendum loop before Task 2 could start; the proof of
  rebuild-after-reject should have been part of Task 1's original acceptance
  evidence rather than discovered at review.
- **Were tasks the right size?** Yes. The five-task cut (seam → activity op →
  update → workflow package → transition) was well-ordered and one-session-sized;
  Task 1b was a small, correctly-scoped remediation, not a sign of mis-slicing.
- **Did the QA review cadence feel right?** Yes. Four mid-phase reviews plus the
  addendum tracked a phase with real foundational risk (the write-path seam and a
  new package); the BLOCK landed at the right moment and the final review re-ran
  the full gate.
- **Were there token-efficiency opportunities missed?** Minor. The Task 1 →
  ESC-003 → Task 1b cycle spent tokens that tighter first-pass acceptance evidence
  on the rebuild-recovery state would have saved; otherwise the phase ran lean.

### Deep Retrospective Recommendation

**No standalone deep retro warranted; carry one planning checklist item forward.**
Phase 5 validated the Phase 4 lesson cleanly, so the durable takeaway is narrow:
when an acceptance criterion asserts *recovery* (rebuild succeeds after a rejected
write), the task's original evidence — not a later review addendum — must
demonstrate the recovery state, because "no bad append" and "system still rebuilds"
are distinct observables (the ESC-003 gap). Fold this into the write-command
slicing checklist already established in Phase 4. The resolution-reason product
question (DEC-015 / CF-008) is a product-roadmap input, not a process retro item.
