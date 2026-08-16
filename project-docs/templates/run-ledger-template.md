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

Every v2 ledger must contain a `Recovery directive schema: v1` marker. After
that marker, recovery entries require the complete typed authority, trigger,
attempt, no-progress budget, verification, and disposition fields below.

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
- Resolution authority: <existing DEC-<NNN> when a decision resolves the prior
  BLOCK | verified commit/restored capability/other recorded external authority
  | N/A when no prior BLOCK is being resolved>
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
- Triage disposition: <N/A for task-scope | exact
  REPAIR / control-conformance>
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
- Phase worktree: <path and branch; for control-conformance, status proving the
  run ledger is the only dirty path>
- Temporary repair worktree: <N/A for task-scope | dedicated path, temporary
  branch, and starting commit supplied to the Developer>
- Hidden-index flags: <none; assume-unchanged and skip-worktree are forbidden>
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
  termination from the repair worktree, with Developer-recorded base/head
  commits, from>
  `project-docs/tools/workflow-control-scope-lint.sh <repo> <base> <head> <allowed...>`
- Phase-branch integration: <N/A for task-scope | fast-forward-only from base
  to repair HEAD after scope PASS, with phase worktree still ledger-only dirty>
- Outcome: <RESOLVED | MATERIAL-PROGRESS | NO-PROGRESS>
- Resolution evidence: <relevant in-scope change tied to the trigger that
  satisfies/removes an obligation; narrower failure; or byte-identical rerun
  result with observed termination. Formatting/metadata/unrelated churn is not
  progress>
- Next checkpoint: <VERIFY/REPORT/REVIEW | another TASK-RECOVERY | triage as
  recovery-budget-exhausted>

For `REPAIR / control-conformance`, leave the dirty run ledger in the phase
worktree, with no other dirt, and create a dedicated clean temporary repair
branch/worktree from the recorded starting commit. The Developer works only in
that supplied repair worktree. The ending checkpoint is one isolated tools-only
control commit. Run the scope check there; it requires that commit to be
checked-out HEAD, a clean index/worktree, and regular-file allowlisted
artifacts, and rejects later commits, dirty state, symlinks, and gitlinks. Never
use `assume-unchanged`, `skip-worktree`, or an equivalent index flag to conceal
dirt. Only after scope PASS, fast-forward the phase branch to repair HEAD,
verify the phase worktree remains ledger-only dirty, and append this REPAIR.
Clean up the temporary worktree/branch only after the append receipt.

`RESOLVED` only permits an ordinary continuation under the same
task/trigger/budget to rerun the originating exact verification and update its
report; it does not complete the task or authorize phase advance. Attempt count
is monotonic for a trigger. Material progress may reset only consecutive
no-progress to `0/2`; a second consecutive `NO-PROGRESS` exhausts the budget
and returns to triage. Necessary provider sync is separate and is not part of
the isolated control commit.

### <YYYY-MM-DD HH:MM> TRIAGE TRG-<phase>-<NNN>

- Trigger: <one line>
- Trigger signature: `<stable signature>`
- Attempt count: <total mediated continuations for this signature>
- Checkpoint: <current durable state>
- Verdict: <RESUME | REPAIR | REPAIR / control-conformance | BLOCK
  (escalated as ESC-<phase>-<NNN>)>
- Resolution evidence: <required for RESUME; repair boundary for REPAIR;
  protected/exhausted condition for BLOCK>
- Triage log: <`.../triage-log.md` entry confirmed | "MISSING - TL-Lite
  returned a verdict without writing the entry">

Record the concrete IDs TL-Lite wrote, never the literal placeholders
`TRG-NNN` / `ESC-NNN`. If TL-Lite returns a verdict carrying a placeholder or
writes no triage-log entry, say so verbatim in the `Triage log:` field. That
field is the durable signal for the Tech Lead's wrap-up triage audit; do not
leave the gap only in prose.

If this TRIAGE resolves a prior BLOCK, reference its phase-scoped
`ESC-<phase>-<NNN>` and the recorded resolution authority. If that authority is
a `DEC-<NNN>`, it must exist. A control-conformance repair after a
decision-bound protected BLOCK requires that existing decision; a new repair
directive is not authority to erase or reinterpret a BLOCK.

### <YYYY-MM-DD HH:MM> STOP

- Cause: <task/review/TRG/ESC reference and one-line reason>
- Triage: <TRG-<phase>-<NNN> (BLOCK, escalated as or referencing
  ESC-<phase>-<NNN>)>. Every STOP carries a triage and ESC reference. Ordinary
  incomplete work, report defects, and gate failures are `REPAIR`, not STOP
  causes.
- Protected category: <closed blocker category>
- Next Human or external-authority action: <one line>

The sequence validator rejects a STOP without the immediately governing
TRIAGE BLOCK and concrete ESC reference. A later RUN-RESUME records the
resolution; it does not retroactively legitimize an untriaged STOP.

### <YYYY-MM-DD HH:MM> RUN-END

- All tasks complete. Final review: `project-docs/docs/phases/phase-<N>/qa/final-review.md`
- Human checklist: `project-docs/docs/phases/phase-<N>/qa/human-checklist.md`
- Phase ready for Tech Lead wrap-up.
