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

A terminated `INVALID` rejection is a conclusive failed append. Validate the
unchanged ledger, correct the candidate against its current state, and resubmit;
retain the rejected output and exit status in the corrected event's evidence.
The byte-identical retry rule applies when termination or the receipt is
missing. A rejected candidate with a valid ledger does not require a stop.

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

Every `RUN-START` or `RUN-RESUME` appended after schema-v2 adoption must itself
contain exactly one `Sequence contract: v2` and exactly one
`Recovery directive schema: v2`; both become effective in that event and no
later RUN may downgrade them. A standalone historical ledger that contains
schema-v1 RUN events remains valid and readable, but its next appended RUN uses
v2. General identity, STOP, authority, budget, and state-transition checks still
apply to historical v1 entries.
Every new v2 ledger RUN therefore carries both markers in the same RUN event.

For every event governed by schema v2, every field shown for that event below
is mandatory exactly once. Use the specified `N/A` value when a field does not
apply; never omit or duplicate the field. Field labels may appear in any order
and need not be contiguous. Ordering is not authority.

### <YYYY-MM-DD HH:MM> RUN-START

- Sequence contract: v2
- Recovery directive schema: v2
- Branch: `<phase branch>`
- ESC baseline (Open at run start): <ESC-NNN list or "none">
- Warm-up/autonomous-resume state: <warm-up | autonomous per DEC-NNN>
- Starting task: <task-M>

### <YYYY-MM-DD HH:MM> RUN-RESUME

Append before a replacement or resumed Phase Runner takes another action.

- Sequence contract: v2
- Recovery directive schema: v2
- Prior run/session or STOP: <reference>
- Governing TRG: <TRG-<phase>-<NNN> retained from STOP | N/A>
- Governing ESC: <ESC-<phase>-<NNN> retained from STOP | N/A>
- Trigger signature: `<task/review>:<category>:<stable source identity> | N/A
  for normal resume`
- Attempt count: <next Phase Runner-mediated continuation number for this
  trigger | N/A for normal resume>
- Checkpoint: <commit/tree state, report, review, and last ledger entry>
- Resolution evidence: <state reconciliation or protected-condition resolution>
- Resolution authority: <N/A for normal resume | verified-commit:<40hex> |
  restored-capability:<stable-id> | DEC-<NNN> for STOP resume>
- Continuing task: <task-M>
- Next recovery: <N/A for normal resume | task-M | category | repair mode |
  exact trigger | attempt for STOP resume>

For a normal resume that does not reconcile a STOP, set `Governing TRG`,
`Governing ESC`, `Trigger signature`, `Attempt count`, `Resolution authority`,
and `Next recovery` to `N/A`; `Continuing task` and checkpoint reconciliation
remain concrete. Omit `Repair owners` entirely, including any `N/A` placeholder.
Choose this form from the current ledger state: an older STOP already followed
by an accepted RUN-RESUME does not authorize another STOP resume. Current-context
summaries may lag the ledger.

A normal resume during an active TASK-RECOVERY preserves its task, trigger,
attempt, authority, and review-owner progress. Continue that recovery after
reconciling existing commits, report, and child-session state; do not append a
duplicate TASK-RECOVERY or restart completed owners. Its `REPAIR_PENDING`
receipt permits continuing that repair, not advancing to a later task.

A pending ordinary schema-v2 REPAIR TRIAGE permits this normal resume only
for its same originating task. The resume preserves its pending authority;
it neither consumes the TRIAGE nor clears the active continuation (DEC-190).

For a STOP resume, use the retained TRG/ESC and concrete trigger, next attempt,
continuing task, and exact recovery tuple. Those fields bind the next
`TASK-RECOVERY`; a different task, trigger, mode, category, or attempt is
invalid. That recovery's `Source` and `Repair authority` also contain the exact
typed `Resolution authority`. `Resolution evidence` contains the
Runner-observed verification.
A STOP resume after a BLOCK re-review may include one optional
`Repair owners: <ordered, distinct task list>` field under DEC-257. Its value
must be byte-identical to the `- Repair owners:` line in the resolution DEC;
the BLOCK must re-review an open review whose owners all completed, and
`Continuing task` and `Next recovery` must name the first new owner.
The resumed route retains the governing TRG/ESC and never allocates replacement
identities for the resolved BLOCK.
`Resolution authority` is exactly one of `verified-commit:<40hex>`,
`restored-capability:<stable-id>`, or `DEC-<NNN>`. A changed or newly resolved
protected decision requires `DEC-<NNN>`; verified commit or restored capability
authority is valid only when the accepted protected contract did not change.
Retain the governing TRG/ESC and do not create a new TRG. Every boundary-valid
DEC citation in an authority field must exist in this phase's `decisions.md`;
embedded strings such as `NOTDEC-123X` are not decision citations.

### <YYYY-MM-DD HH:MM> TASK task-<M> COMPLETE

- Report: `project-docs/docs/phases/phase-<N>/tasks/task-<M>-done.md`
- report-lint: <PASS - OUTCOME COMPLETE>
- Commit: `<hash>` - verified, subject matched
- QA-plan commands beyond the task spec: <recorded PASS | none required>

### <YYYY-MM-DD HH:MM> REVIEW review-<M>

- Path: `project-docs/docs/phases/phase-<N>/qa/review-<M>.md`
- Tasks covered: <task range>
- Outcome: <PASS | REPAIR | BLOCK>
- Owning task: <ordered comma-separated task IDs for REPAIR, e.g. `task-1, task-3` | none>
- Re-review of: <prior review ID for a repair-closing PASS | none>
- Must-fix: <n> / Improvement: <n> / blocker candidates: <n; QA minted no ESC>

Use exactly one `Owning task` field. For `REPAIR`, list every Must-fix owner in
the review's Repair Plan row order; that order is the repair order. Each listed
owner needs its own `TASK-RECOVERY`, `REPAIR`, and `TASK COMPLETE` before a
linked re-review may close or replace the review. `none` is only for `PASS` and
`BLOCK`.

### <YYYY-MM-DD HH:MM> TASK-RECOVERY task-<M>

Append when a report, review, or interrupted session exposes real unfinished
work that remains local, in scope, and decision-free. This event authorizes a
same-task Developer continuation; it does not close the task or finding.

- Source: <report, review finding, command attempt, or session reference; after
  STOP include the exact typed Resolution authority>
- Category: <in-scope-implementation | gate-repair | report-repair | qa-rework |
  session-continuation | local-environment>
- Repair mode: <task-scope | control-conformance>
- Repair authority: <task spec / QA review / TRG-<phase>-<NNN> written
  control-conformance directive; after STOP include the exact typed Resolution
  authority>
- Triage disposition: <N/A for task-scope | exact
  REPAIR / control-conformance>
- Trigger signature: `<task/review>:<category>:<stable source identity>`
- Attempt count: <next total mediated continuation number for this trigger>
- Control contract: <exact directive value for control-conformance | N/A>
- Control mismatch: <exact directive value for control-conformance | N/A>
- Allowed files: <ordinary task scope | exact comma-separated ordered directive tools-only paths with no spaces around commas>
- Forbidden changes: <exact directive protected paths/semantics | N/A>
- Originating report: <exact directive report path | N/A>
- Exact verification: <exact directive command bytes and working directory | N/A>
- Control repair base: <runner base for control/linked report-repair | N/A>
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

Under schema v2, a control-conformance entry is valid only when its
`Repair authority` names the governing TRIAGE. For the initial attempt, its
`Trigger signature` and `Attempt count` map exactly to that triage record's
`Control trigger signature` and `Control attempt count`; `Control contract`,
`Control mismatch`, `Allowed files`, `Forbidden changes`, `Originating report`,
`Exact verification`, and `Control repair base` match their same-named
directive values byte-for-byte. The heading's task must match
`Originating task`. Each directive value is one physical line so copying cannot
normalize wrapping. Each of the ten directive labels occurs exactly once; the
labels may appear in any order and need not be contiguous.

For an ordinary TRIAGE `REPAIR`, the heading task, `Trigger signature`,
`Attempt count`, `Category`, and `Repair authority` bind byte-for-byte to the
external triage record's `Originating task`, `Trigger signature`, `Attempt
count`, `Repair category`, and `Repair authority`. A triage record is not
merely a reference: it is the authority for those recovery identity fields.
`Source` references the governing TRG. Do not prefix or otherwise change the
ordinary `Repair authority` merely to repeat that TRG ID; control-conformance
authority naming remains governed by the preceding paragraph.
If the TRIAGE `REPAIR` classified a blocker candidate returned by an active
schema-v2 recovery awaiting its REPAIR event, this bound TASK-RECOVERY replaces
that pending dispatch only for the same task. Other overlap and recovery-stage
checks remain unchanged.
An unconsumed ordinary schema-v2 REPAIR TRIAGE may itself be corrected by a
later ordinary REPAIR TRIAGE for the same task and dispatch attempt (DEC-190).
The correcting external TRG becomes the exact recovery authority. Any active
continuation retains its trigger and next attempt; control directives and
recovery budgets cannot be replaced by this correction. Review-form triggers
name the active owning review and use `review-M:qa-rework:task-N-<source>`
(the exact task token may also end the value or be followed by `:`).
After `MATERIAL-PROGRESS` or the first `NO-PROGRESS`, the next
control-conformance `TASK-RECOVERY` increments `Attempt count`, records the
prior integrated HEAD as its new runner-owned `Control repair base`, and keeps
authority, contract, mismatch, allowlist, forbidden set, report, and
verification byte-identical. A Developer never supplies or changes that base.

### <YYYY-MM-DD HH:MM> REPAIR task-<M>

Append after one Phase Runner-mediated continuation. Preserve earlier failed
evidence; this event records what the continuation changed and proved.

- Trigger signature: `<same signature as TASK-RECOVERY>`
- Attempt count: <completed mediated continuation number for this trigger>
- Starting checkpoint: <TASK-RECOVERY checkpoint>
- Ending checkpoint: <commit or exact durable tree/report/review state>
- Control repair base: <same runner-owned TASK-RECOVERY base | N/A>
- Control repair head: <full Developer-produced checked-out HEAD | N/A>
- Control scope check: <N/A for task-scope | PASS control repair scope: <base> -> <head>; allowed=<exact Allowed files> (terminated exit_code: 0) from the trusted-base verifier blob>
- Phase-branch integration: <N/A for task-scope | PASS fast-forward from
  directive base to repair HEAD after scope PASS, phase worktree ledger-only>
- Outcome: <RESOLVED | MATERIAL-PROGRESS | NO-PROGRESS | INTERRUPTED>
- Interruption cause: <required exactly once for INTERRUPTED; otherwise omit>
- Elapsed session duration: <required exactly once for INTERRUPTED; otherwise omit>
- Command in flight: <exact command | none observed | unknown after external
  termination; required exactly once for INTERRUPTED; otherwise omit>
- Resolution evidence: <relevant in-scope change tied to the trigger that
  satisfies/removes an obligation; narrower failure; or byte-identical rerun
  result with observed termination. Formatting/metadata/unrelated churn is not
  progress>
- Next checkpoint: <VERIFY/REPORT/REVIEW | another TASK-RECOVERY | triage as
  recovery-budget-exhausted>

`INTERRUPTED` means the child was ended from outside before returning a typed
result. It is not `NO-PROGRESS`: the next TASK-RECOVERY increments total
`Attempt count` but preserves the prior `Consecutive no-progress attempts`
value. A parent polling wait timeout is not an interruption and produces no
REPAIR event; retain and poll the live child. Two consecutive `INTERRUPTED`
outcomes require orchestration session-liveness TRIAGE before another recovery.
A TRIAGE `REPAIR` advances the next total `Attempt count`, preserves the prior
no-progress count, and starts a new consecutive-interruption streak.

For `REPAIR / control-conformance`, leave the dirty run ledger in the phase
worktree, with no other dirt, and create a dedicated clean temporary repair
branch/worktree from `Control repair base`. The Developer works only in that
supplied worktree and cannot rebaseline it. Run the scope check there from the
immutable directive base to the Developer-produced ending `HEAD`; the repair
worktree must be clean and every allowlisted artifact a regular file. Dirty
state, symlinks, gitlinks, hidden index flags, and any path outside the exact
allowlist fail. Never use `assume-unchanged`, `skip-worktree`, or an equivalent
index flag to conceal dirt. Materialize
`project-docs/tools/workflow-control-scope-lint.sh` from the trusted base commit
into fresh scratch storage and execute that blob, never the repair HEAD's copy.
The scope receipt must repeat the exact comma-separated ordered `Allowed files`
value, and the Runner rejects a missing, reordered, normalized, or different
receipt. Normal control-conformance may not repair
`workflow-control-scope-lint.sh` or `workflow-sequence-lint.sh`; either proposal
is a protected verifier change and returns to BLOCK/decision authority.

After every scope PASS, fast-forward the phase branch to repair HEAD, verify it
remains ledger-only dirty, append REPAIR, and clean up after the append receipt.
After `MATERIAL-PROGRESS` or `NO-PROGRESS`, issue the next control recovery as
specified above, using that integrated HEAD as its runner-owned base.

A control `RESOLVED` never permits direct task completion. It permits exactly
one linked `TASK-RECOVERY` with `Category: report-repair`, `Repair mode:
task-scope`, `Allowed files` equal only the exact `Originating report`, and the
same `Trigger signature`, `Originating report`, `Exact verification`, and
resolved attempt's `Control repair base`. That Developer may edit only that
report and must run the exact verification. Only a
subsequent report-repair `REPAIR` with `Outcome: RESOLVED` permits the task's
completion event. Attempt count remains monotonic. Material progress resets
only consecutive no-progress to `0/2`; a second consecutive `NO-PROGRESS`
exhausts the budget and returns to triage. Necessary provider sync is separate
and is not part of the isolated control history.

### <YYYY-MM-DD HH:MM> TRIAGE TRG-<phase>-<NNN>

- Trigger: <one line>
- Originating task: <task-M>
- Trigger signature: `<stable signature>`
- Attempt count: <next mediated continuation number for this trigger>
- Repair category: <recovery category | N/A for RESUME/BLOCK>
- Repair authority: <task spec / QA review / TRG-<phase>-<NNN> typed directive | N/A for non-REPAIR>
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

For `REPAIR / control-conformance`, the external triage-log entry's ten typed
fields are authoritative. The Phase Runner copies their values byte-for-byte
into the initial `TASK-RECOVERY`; it may not paraphrase, widen, or
self-author them. `Control repair base` is the full phase HEAD recorded by the
Phase Runner before triage and dispatch, never a Developer-selected base.

For every schema-v2 TRIAGE, the ledger's `Originating task`, `Trigger
signature`, `Attempt count`, `Repair category`, and `Repair authority` match the
ordinary fields in its external TRG section byte-for-byte. The control
directive adds its ten independently labelled values; label position does not
affect matching.

When an existing DEC resolves a prior BLOCK into control-conformance repair,
do not append a new TRIAGE/TRG. Append one unique, separate H2 directive
amendment whose heading names the existing governing TRG. It is parsed
independently from the original TRG section and its retained BLOCK verdict; it
is not nested in or treated as a rewrite of that section. Retain the TRG/ESC,
cite the governing DEC, then append RUN-RESUME with those same IDs. Every
boundary-valid DEC citation in an authority field must resolve to this phase's
decisions file; one real citation cannot mask another nonexistent or embedded
pseudo-ID.

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
