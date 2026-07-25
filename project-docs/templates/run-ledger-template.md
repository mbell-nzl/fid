# Phase <N> Run Ledger

> Maintainer: Phase Runner
> Location: `project-docs/docs/phases/phase-<N>/run-ledger.md`
>
> Append-only. The Phase Runner appends one entry per event and never edits or
> deletes prior entries. This file, not chat notes, is the durable run state:
> a resumed or replacement runner session reads it first and reconciles it
> against completion reports and git history before continuing.

Entry types, in the order they typically occur. Timestamps are
`YYYY-MM-DD HH:MM` local time. Append new entries at the end of the file, in
wall-clock order, one entry per event.

Before appending, read the last entry. If the new entry would duplicate it or
carry an earlier timestamp, you have lost track of run state: reconcile against
completion reports and git history first. Never edit or delete a prior entry to
fix an ordering or duplication mistake — append a `CORRECTION` entry naming the
entries it supersedes and why. The ledger is evidence; its errors are part of
the record.

### <YYYY-MM-DD HH:MM> RUN-START

- Branch: `<phase branch>`
- ESC baseline (Open at run start): <ESC-NNN list or "none">
- Warm-up/autonomous-resume state: <warm-up | autonomous per DEC-NNN>
- Starting task: <task-M>

### <YYYY-MM-DD HH:MM> TASK task-<M> COMPLETE

- Report: `project-docs/docs/phases/phase-<N>/tasks/task-<M>-done.md`
- report-lint: <PASS | FAIL (n failures) - hard stop>
- Commit: `<hash>` - verified, subject matched
- QA-plan commands beyond the task spec: <recorded PASS | none required>

### <YYYY-MM-DD HH:MM> REVIEW review-<M>

- Path: `project-docs/docs/phases/phase-<N>/qa/review-<M>.md`
- Tasks covered: <task range>
- Blocking: <n> / Improvement: <n> / ESC raised: <list or "none">

### <YYYY-MM-DD HH:MM> TRIAGE TRG-<phase>-<NNN>

- Trigger: <one line>
- Verdict: <RESUME | BLOCK (escalated as ESC-<phase>-<NNN>)>
- Triage log: <`.../triage-log.md` entry confirmed | "MISSING - TL-Lite
  returned a verdict without writing the entry">

Record the concrete IDs TL-Lite wrote, never the literal placeholders
`TRG-NNN` / `ESC-NNN`. If TL-Lite returns a verdict carrying a placeholder or
writes no triage-log entry, say so verbatim in the `Triage log:` field. That
field is the durable signal for the Tech Lead's wrap-up triage audit; do not
leave the gap only in prose.

### <YYYY-MM-DD HH:MM> STOP

- Cause: <task/review/TRG/ESC reference and one-line reason>
- Triage: <TRG-<phase>-<NNN> (BLOCK, escalated as ESC-<phase>-<NNN>) |
  "immediate stop - no triage required">. Triage-gated causes (report-lint
  failure, missing report, blocker report, gate failure,
  placeholder/missing/skipped/narrower/ad hoc command) must carry a TRG
  reference; a STOP for those causes without one means TL-Lite was skipped.
- Next human action: <one line>

### <YYYY-MM-DD HH:MM> RUN-END

- All tasks complete. Final review: `project-docs/docs/phases/phase-<N>/qa/final-review.md`
- Human checklist: `project-docs/docs/phases/phase-<N>/qa/human-checklist.md`
- Phase ready for Tech Lead wrap-up.
