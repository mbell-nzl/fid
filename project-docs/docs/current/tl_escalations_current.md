# TL Escalations — Current

## Open escalations

### ESC-24-010 - Phase-end Human acceptance for VAL-053 through VAL-057

- **Status:** OPEN (2026-10-11). Dispositioned 2026-10-11 by DEC-292: the
  F1/F2 routes, the ledger route and the final-review-1 scope are decided.
  B1-H still waits for the Human's committed checklist and DEC-292 clause 9.
- **Governing triage:** TRG-24-019 (corrected task-bound identity; the STOP and
  DEC-292 bind to it). TRG-24-018 is the historical phase-level record.
- **Source:** project-docs/docs/phases/phase-24/qa/final-review.md; QA artifacts at 6d633351774d804a80b90a8c888d2a93cc0d615b, with invoking Runner's attached commit-resolution and exact-subject verification evidence.
- **Originating task:** task-14 (identity anchor only; F1/F2 keep their owners)
- **Trigger signature:** task-14:qa-block:phase-end-final-review (TRG-24-018 used
  phase-24:qa-block:final-review-b1h-human-judgment-and-report-receipts)
- **Attempt count:** 1
- **Protected category:** product-owner judgment / unresolved Human-only validation (B1-H).
- **Condition:** Five phase-end Human assertions VAL-053 through VAL-057 remain BLOCKED; the phase-end checklist has no recorded outcomes. The Task 11 walkthrough was interim and cannot substitute. Workflow section 2.6 requires BLOCK for unresolved Human-only validation. ESC-24-007 resolved only the earlier Window D2 checkpoint and is not reopened or reused.
- **Retained deterministic evidence:** Final review records 63 deterministic VAL PASS and passing required product gates, full-checkptr, archive-scale and three rebuild samples. These results remain intact; they do not supply Human acceptance. QA-36 remains FAIL on F1/F2.
- **Preserved bounded F1:** Owner phase-24/task-1; only project-docs/docs/phases/phase-24/tasks/task-1-done.md. Fresh report-author execution of rendered R4 `ls migrations; rg -n "ProjectionSchemaVersion|projection_schema" internal -g '!*_test.go'` at /home/mark/projects/fid, own-session output/status/observed termination and missing-receipt note, preserving history. Independent exact two-argument report lint against task-1.md must return OUTCOME COMPLETE/exit 0. Fresh independent QA Lead reruns R4 and report lint and reviews only the report delta/history and QA-36.
- **Preserved bounded F2:** Owner phase-24/task-10; only project-docs/docs/phases/phase-24/tasks/task-10-done.md. Fresh report-author execution of R10 `rg -n '"error: ' cmd internal -g '*_test.go'` at /home/mark/projects/fid, own-session nonempty output/exit 0 with observed termination and missing-receipt note, preserving failures and history. Independent exact two-argument report lint against task-10.md must return OUTCOME COMPLETE/exit 0. Fresh independent QA Lead reruns R10 and report lint and reviews the report delta/history/identity and QA-36. Exact lint commands and full boundaries are retained in TRG-24-018.
- **Requested resolution:** Human completes project-docs/docs/phases/phase-24/qa/human-checklist.md using the prepared public fixtures and permitted guidance, supplying actual outcomes/reasons and required attempts/inputs for all five assertions. Full Tech Lead arranges the two separate bounded report continuations and classifies any Human rejection/partial result without weakening the accepted contract. Fresh independent QA Lead then writes the next unused final-review-R.md, citing the unchanged initial final review and closing B1-H, F1/F2 and QA-36 with evidence.
- **Resume condition:** Protected Human condition resolved through recorded actual judgment and appropriate authority, both report findings independently closed, and fresh final-review acceptance. A report fix or verified commit containing the original BLOCK review alone is insufficient. Runner retains TRG-24-018 / ESC-24-010 for STOP resolution; no phase advance or wrap-up while open.
- **Next:** The Human records every row of `qa/human-checklist.md` (DEC-292
  clause 1; live workspaces expire about 11:50 NZDT 2026-10-12, after which
  the fixtures are regenerated first) and commits it. Then invoke the Tech
  Lead with "Close ESC-24-010 under DEC-292 clause 9." The Runner does not
  resume before clause 9 (DEC-292 clause 5).
- **References:** Phase 24 validation contract VAL-053 through VAL-057; QA plan sections 8-9; human-checklist.md instructions/Summary; final-review.md Findings/Repair Plan/Blocker Candidates; input-inventory.md rendered-command header/R4/R10; 09_AI_WORKFLOW.md sections 2.3, 2.6 and 9.

## Recently resolved

Phase 24, archived in [Phase 24 escalations](../phases/phase-24/tl-escalations.md):
- ESC-24-009: Task 13's R13 receipt directive carried the Markdown
  table escape `\|` as command bytes, so `rg` found nothing (exit 1).
  Resolved 2026-10-11 by DEC-291: a table command is its rendered text,
  so R13 is `rg -n "p23CheckJudgmentQualifiers|PRIVATE_" cmd`, unchanged.
  The Runner resumes the report-only Task 13 recovery (attempt 3), then
  review-12 re-reviews F1.
- ESC-24-008: Task 14's journey had no permitted source for the
  terminal-success `--to` state (`qa_pass`) or the Issue revision.
  Resolved 2026-10-11 by DEC-290: `resume` shows `Revision:` and
  `Next states:` (JSON `issue.next_states`), and `issue transition` help
  points there. Task 14 owns the change (VAL-068). The Runner appends the
  clause 8 RUN-RESUME and TASK-RECOVERY (attempt 2).
- ESC-24-007: Window D2 (review-9 BLOCK, B1-H) lacked the interim
  VAL-054/055 Human walkthrough. Resolved 2026-10-10 by DEC-288: the
  Human recorded all nine rows (`a23b9f33`). Five refusals were clear
  and four views were accepted. The Runner resumes through an
  evidence-only Task 11 recovery (attempt 2, `qa-rework`), then a fresh
  QA review-10 closes B1-H.
- ESC-24-006: Task 11/VAL-041 required plain-copy `run.inspect` to
  render `different worktree`, but its retained Phase 22 refusal
  (`recovery-state-conflict`, exit 4) is kept by contract section 2.
  Resolved 2026-10-09 by DEC-287: the refusal stays, unchanged in both
  modes. `resume` proves the plain-copy state, and the one fixture row
  asserts the exact refusal. The Runner appends the clause 7 RUN-RESUME
  and TASK-RECOVERY (attempt 2).
- ESC-24-005: `export import --project` printed raw text and exited 7 in
  both modes, with no `fid-cli.v1` failure. Resolved 2026-10-09 by
  DEC-286: it is a Phase 21 section 2 forbidden flag, so it returns
  `invalid-argument`/2 with the flag spelling and the task 9a allowed
  list. Task 10 gains `root.go` and a `canonical_request.go` extraction.
  The Runner appends the clause 8 RUN-RESUME and TASK-RECOVERY (attempt 2).
- ESC-24-004: `activity.list` had no `fid-cli.v1` failure, so Task 10
  had no accepted mapping to adopt. Resolved 2026-10-07 by DEC-285: it
  adopts the retained-read mapper (Phase 21 sections 1/4.1/7), and Task 10
  gains the stats and AC-list callers. The Runner appends the clause 8
  RUN-RESUME and TASK-RECOVERY (attempt 2).
- ESC-24-003: the sequence verifier rejected the corrected Task 9
  recovery; resolved 2026-10-07 by DEC-284 without a verifier change. The
  Runner appends correcting TRIAGE TRG-24-008 and the bound recovery, then
  runs QA window C after Task 9.
- ESC-24-002: the accepted Task 6 RUN-RESUME bound an unsatisfiable recovery
  tuple; resolved 2026-10-06 by DEC-283. The Runner then does the verifier
  maintenance and the ledger sequence.
- ESC-24-001: Task 6 Developer model capacity; restored and verified by the
  2026-10-05 20:48 RUN-RESUME. Closure recorded in DEC-283.

Phase 23, all resolved and archived in
[Phase 23 escalations](../phases/phase-23/tl-escalations.md):
- ESC-23-012: phase-end final review blockers and Human judgment; resolved
  2026-10-03 by DEC-267..271 (tasks 19-21, final review 2 PASS, Human
  checklist 2 7/7 PASS).
- ESC-23-011 (model capacity), ESC-23-010 (DEC-265/266), ESC-23-009
  (DEC-264), ESC-23-008 (DEC-263), ESC-23-007 (DEC-262), ESC-23-006
  (DEC-260), ESC-23-005 (DEC-258), ESC-23-004 (DEC-249..256), ESC-23-003
  (DEC-239), ESC-23-002 (DEC-238), ESC-23-001 (DEC-235).

## Standing handling

Register an escalation when allocated, not after resolution. Ordinary causal
repair and incomplete in-scope work do not themselves require a Human stop.
Read accepted ownership before calling it ambiguous (DEC-194). A real contract/
schema contradiction needs the missing authority, not a weakened assertion
(DEC-195). External interruptions do not consume repair no-progress budget.

Prior archives: [Phase 22](../phases/phase-22/tl-escalations.md),
[Phase 21](../phases/phase-21/tl-escalations.md),
[Phase 20](../phases/phase-20/tl-escalations.md).
