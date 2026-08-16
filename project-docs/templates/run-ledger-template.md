# Phase <N> Run Ledger

> Maintainer: Phase Runner
> Location: `project-docs/docs/phases/phase-<N>/run-ledger.md`
>
> Append-only. The Phase Runner appends one entry per event through
> `project-docs/tools/workflow-sequence-lint.sh <run-ledger> --append-stdin`
> and never edits or deletes prior accepted entries. The command validates and
> atomically appends at EOF, returning `APPEND_CONFIRMED state=<state>`; an
> `INVALID` result leaves the ledger unchanged. If the receipt or termination is
> missing, retry the byte-identical candidate to confirm the exact EOF event
> without duplicating it. This file, not chat notes, is the durable run state: a
> resumed or replacement runner session reads it first and reconciles it
> against completion reports and git history before continuing.

Entry types, in the order they typically occur. Timestamps are
`YYYY-MM-DD HH:MM` local time. Append new entries at the end of the file, in
wall-clock order, one entry per event.

Before appending, read the last entry. If the new entry would duplicate it or
carry an earlier timestamp, you have lost track of run state: reconcile against
completion reports and git history first. Submit the event Markdown beginning
with its single `###` heading to the transactional append command; do not use
`apply_patch`, redirection, or another editor. Never edit or delete a prior
accepted entry to fix an ordering or duplication mistake — append a `CORRECTION`
entry naming the entries it supersedes and why. The ledger is evidence; its
errors are part of the record.

### <YYYY-MM-DD HH:MM> RUN-START

- Sequence contract: v2
- Recovery directive schema: v1
- Branch: `<phase branch>`
- ESC baseline (Open at run start): <ESC-NNN list or "none">
- Warm-up/autonomous-resume state: <warm-up | autonomous per DEC-NNN>
- Starting task: <task-M>

### <YYYY-MM-DD HH:MM> RUN-RESUME

Append before a replacement or resumed Phase Runner takes another action.

- Sequence contract: v2
- Recovery directive schema: v1
- Prior run/session or STOP: <reference>
- Trigger signature: `<task/review>:<category>:<stable source identity>`
- Attempt count: <total Phase Runner-mediated continuations for this signature>
- Checkpoint: <commit/tree state, report, review, and last ledger entry>
- Resolution evidence: <state reconciliation or protected-condition resolution>
- Continuing task: <task-M>

### <YYYY-MM-DD HH:MM> TASK task-<M> COMPLETE

- Report: `project-docs/docs/phases/phase-<N>/tasks/task-<M>-done.md`
- report-lint: <PASS - OUTCOME COMPLETE>
- Commit: `<hash>` - verified, subject matched
- QA-plan commands beyond the task spec: <recorded PASS | none required>

### <YYYY-MM-DD HH:MM> REVIEW review-<M>

- Path: `project-docs/docs/phases/phase-<N>/qa/review-<M>.md`
- Tasks covered: <task range>
- Outcome: <PASS | REPAIR | BLOCK>
- Owning task: <task-M for REPAIR | none>
- Re-review of: <prior review ID for a repair-closing PASS | none>
- Must-fix: <n> / Improvement: <n> / blocker candidates: <n; QA minted no ESC>

### <YYYY-MM-DD HH:MM> TASK-RECOVERY task-<M>

Append when a report, review, or interrupted session exposes real unfinished
work that remains local, in scope, and decision-free. This event authorizes a
same-task Developer continuation; it does not close the task or finding.

- Source: <report, review finding, command attempt, or session reference>
- Category: <in-scope-implementation | gate-repair | report-repair | qa-rework |
  session-continuation | local-environment>
- Repair mode: <task-scope | control-conformance>
- Repair authority: <task spec / QA review / TRG-<phase>-<NNN> written
  control-conformance directive>
- Control contract: <for control-conformance, binding source and uniquely
  required behavior; otherwise N/A>
- Control mismatch: <for control-conformance, exact reproducer, expected and
  observed result, and termination; otherwise N/A>
- Allowed files: <ordinary task scope | exact named control/self-test paths
  under project-docs/tools/>
- Forbidden changes: <for control-conformance, product/spec/prompt/template/
  dependency/task/QA-contract/report files and command/assertion/status/
  semantics/scope/threshold/budget/tolerance/policy changes>
- Trigger signature: `<task/review>:<category>:<stable source identity>`
- Attempt count: <next total mediated continuation number for this signature>
- Consecutive no-progress attempts: <n>/2
- Checkpoint: <commit or exact dirty-tree/report/review state>
- Resolution evidence: <PENDING - exact verification and any independent
  re-review required to prove this trigger resolved>
- Next action: <bounded repair and exact verification/re-review>
- Human attention required: no

### <YYYY-MM-DD HH:MM> REPAIR task-<M>

Append after one Phase Runner-mediated continuation. Preserve earlier failed
evidence; this event records what the continuation changed and proved.

- Trigger signature: `<same signature as TASK-RECOVERY>`
- Attempt count: <completed mediated continuation number for this signature>
- Starting checkpoint: <TASK-RECOVERY checkpoint>
- Ending checkpoint: <commit or exact durable tree/report/review state>
- Control scope check: <N/A for task-scope | PASS token and observed
  termination, with Developer-recorded base/head commits, from>
  `project-docs/tools/workflow-control-scope-lint.sh <repo> <base> <head> <allowed...>`
- Outcome: <RESOLVED | MATERIAL-PROGRESS | NO-PROGRESS>
- Resolution evidence: <relevant in-scope change tied to the trigger that
  satisfies/removes an obligation; narrower failure; or byte-identical rerun
  result with observed termination. Formatting/metadata/unrelated churn is not
  progress>
- Next checkpoint: <VERIFY/REPORT/REVIEW | another TASK-RECOVERY | triage as
  recovery-budget-exhausted>

For control-conformance, the ending checkpoint is one isolated tools-only
control commit. Append this REPAIR only after the scope check passes, then
dispatch an ordinary continuation under the same task/trigger/budget to rerun
the originating exact verification and update its report. Necessary provider
sync is separate and is not part of the isolated control commit.

### <YYYY-MM-DD HH:MM> TRIAGE TRG-<phase>-<NNN>

- Trigger: <one line>
- Trigger signature: `<stable signature>`
- Attempt count: <total mediated continuations for this signature>
- Checkpoint: <current durable state>
- Verdict: <RESUME | REPAIR | BLOCK (escalated as ESC-<phase>-<NNN>)>
- Resolution evidence: <required for RESUME; repair boundary for REPAIR;
  protected/exhausted condition for BLOCK>
- Triage log: <`.../triage-log.md` entry confirmed | "MISSING - TL-Lite
  returned a verdict without writing the entry">

Record the concrete IDs TL-Lite wrote, never the literal placeholders
`TRG-NNN` / `ESC-NNN`. If TL-Lite returns a verdict carrying a placeholder or
writes no triage-log entry, say so verbatim in the `Triage log:` field. That
field is the durable signal for the Tech Lead's wrap-up triage audit; do not
leave the gap only in prose.

### <YYYY-MM-DD HH:MM> STOP

- Cause: <task/review/TRG/ESC reference and one-line reason>
- Triage: <TRG-<phase>-<NNN> (BLOCK, escalated as or referencing
  ESC-<phase>-<NNN>)>. Every STOP carries a triage and ESC reference. Ordinary
  incomplete work, report defects, and gate failures are `REPAIR`, not STOP
  causes.
- Protected category: <closed blocker category>
- Next Human or external-authority action: <one line>

### <YYYY-MM-DD HH:MM> RUN-END

- All tasks complete. Final review: `project-docs/docs/phases/phase-<N>/qa/final-review.md`
- Human checklist: `project-docs/docs/phases/phase-<N>/qa/human-checklist.md`
- Phase ready for Tech Lead wrap-up.
