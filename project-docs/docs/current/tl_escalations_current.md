# TL Escalations — Current

## Open escalations

None.

## Recently resolved

Phase 24, archived in [Phase 24 escalations](../phases/phase-24/tl-escalations.md):
- ESC-24-010: phase-end Human acceptance for VAL-053..057 and the
  F1/F2 report receipts. Resolved 2026-10-29 by DEC-293: the checklist is
  recorded (VAL-053 FAIL; HUMAN-05 superseded), its findings carry forward
  to the fid reset, F1/F2 are waived and Phase 24 closes.
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
