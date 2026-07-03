# Phase <N> Run Ledger

> Maintainer: Phase Runner
> Location: `project-docs/docs/phases/phase-<N>/run-ledger.md`
>
> Append-only. The Phase Runner appends one entry per event and never edits or
> deletes prior entries. This file, not chat notes, is the durable run state:
> a resumed or replacement runner session reads it first and reconciles it
> against completion reports and git history before continuing.

Entry types, in the order they typically occur. Timestamps are
`YYYY-MM-DD HH:MM` local time. Append new entries at the end of the file.

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

### <YYYY-MM-DD HH:MM> TRIAGE TRG-<NNN>

- Trigger: <one line>
- Verdict: <RESUME | BLOCK (escalated as ESC-NNN)>

### <YYYY-MM-DD HH:MM> STOP

- Cause: <task/review/TRG/ESC reference and one-line reason>
- Next human action: <one line>

### <YYYY-MM-DD HH:MM> RUN-END

- All tasks complete. Final review: `project-docs/docs/phases/phase-<N>/qa/final-review.md`
- Human checklist: `project-docs/docs/phases/phase-<N>/qa/human-checklist.md`
- Phase ready for Tech Lead wrap-up.
