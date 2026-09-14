# Phase 22 protected workflow recovery

Authority: Human request of 2026-09-06, recorded as DEC-190. This single
maintenance task resolves ESC-22-001; it does not implement product task 2.
Starting HEAD: `3e06f12fe4ae42e2694958043b95741bbb2c5c8a`, branch `phase-22`.
Existing uncommitted TRG-22-006 and the escalation register are preserved.

## Diagnosis and acceptance

The ledger ends in ordinary TRIAGE TRG-22-003 after task-2 review attempt 1
returned MATERIAL-PROGRESS. Its pending directive has the wrong trigger for
that continuation. DEC-189/TRG-22-005 correct the trigger and scope. The
validator rejects the correcting TRIAGE before checking that it belongs to
the same task, requires task-prefix-only triggers at several boundaries, and
also rejects the required normal Runner restart while that triage is pending.

Prove correction and restart through the real append entrypoint, with negative
tests for cross-task substitution, changed attempts, stale authority, forged
review identity, budget/STOP bypass, and rejected-append immutability. Retain
all existing mechanical drills. No phase VAL applies: these are workflow
controls outside the product validation contract; no product VAL is closed.

## Exact Commands To Run

Working directory for every command: `/home/mark/projects/fid`.

- `sh -n project-docs/tools/workflow-sequence-lint.sh project-docs/tools/drills/review-recovery.sh`
- `sh project-docs/tools/drills/review-recovery.sh`
- `project-docs/tools/drills/run-drills.sh`
- `project-docs/tools/sync-agents.sh` (DEC-191 canonical configuration repair)
- `project-docs/tools/workflow-lint.sh`
- `project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-22/run-ledger.md`
- `git diff --check`

Run each candidate through the unchanged
`project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-22/run-ledger.md --append-stdin`
entrypoint, capturing its own termination and `APPEND_CONFIRMED` receipt.
The normal resume, correcting TRIAGE and exact recovery are sequential
transactions, not a hand-edited ledger or a batch-event bypass.

## Evidence

- Baseline standalone sequence: `VALID state=REPAIR_PENDING`
  (terminated `exit_code: 0`).
- Red focused regression: `REVIEW RECOVERY FAIL 7 checks` (terminated
  `exit_code: 1`). It reproduced `INVALID schema v2 TRIAGE Originating task
  does not match Trigger signature` and `INVALID TRIAGE REPAIR must be
  followed by its bound TASK-RECOVERY before RUN-RESUME` / `before TRIAGE`.
- First green: `REVIEW RECOVERY PASS all checks` (terminated `exit_code: 0`).
  Root cause: unconditional pending-triage barrier and task-prefix-only
  identity checks. Repair: bounded ordinary reconciliation and an owning-review
  identity predicate shared by triage, budget/session, and STOP resume checks.
- Added STOP/session regression attempt: `REVIEW RECOVERY FAIL 7 checks`
  (terminated `exit_code: 1`). Its external BLOCK fixture incorrectly used the
  ledger's decorated verdict instead of the required raw `BLOCK`; the first
  failure was `INVALID TRIAGE verdict does not match the phase triage log:
  TRG-1-002`. Corrected only that fixture's representation and separate ESC
  reference, preserving the contract and assertion. Byte-identical focused
  rerun: `REVIEW RECOVERY PASS all checks` (terminated `exit_code: 0`). All 30
  positive/negative append checks passed. Both failed attempts remain evidence.
- Existing mechanical suite before adding the focused script to its runner:
  `DRILLS PASS mechanical drills completed with no misses` (terminated
  `exit_code: 0`). Final integrated run also passed with the same token and
  observed `exit_code: 0`, including all 30 new regression checks.
- Workflow lint first attempt: `FAIL workflow-lint found 1 failure(s)`
  (terminated `exit_code: 1`). Diagnostic byte-identical second attempt retained
  the same final token and `exit_code: 1`; its detailed output identified
  generated drift in Codex phase_runner, qa_lead, tech_lead and Claude tech-lead.
  Cause and repair authority are DEC-191: canonicalize the committed Human
  configuration and regenerate providers, without reverting the model choices.
  These failures remain immutable evidence. Post-repair byte-identical lint
  passed: `PASS workflow-lint completed with no failures` (terminated
  `exit_code: 0`).

## Required gate receipts

- `sh -n project-docs/tools/workflow-sequence-lint.sh project-docs/tools/drills/review-recovery.sh` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `sh project-docs/tools/drills/review-recovery.sh` - PASS - `REVIEW RECOVERY PASS all checks` (terminated `exit_code: 0`); failed attempts above retained.
- `project-docs/tools/drills/run-drills.sh` - PASS - `DRILLS PASS mechanical drills completed with no misses` (terminated `exit_code: 0`).
- `project-docs/tools/sync-agents.sh` - PASS - `Generated .codex/agents/tl_lite.toml` (terminated `exit_code: 0`).
- `project-docs/tools/workflow-lint.sh` - PASS - `PASS workflow-lint completed with no failures` (terminated `exit_code: 0`); failed attempts above retained.
- `project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-22/run-ledger.md` - PASS - `VALID state=REPAIR_PENDING` (terminated `exit_code: 0`), before appending reconciliation.
- `git diff --check` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).

## Real Ledger Resolution

Maintenance commit: `3bb42b4efd391b216b752db2c31fe793f3538bac` —
`fix(workflow): unblock reviewed task recovery and runner restart`.

From `/home/mark/projects/fid`, each of the following candidates was supplied
separately to
`project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-22/run-ledger.md --append-stdin`:

- Normal RUN-RESUME — PASS — `APPEND_CONFIRMED state=REPAIR_PENDING`
  (terminated `exit_code: 0`).
- Correcting TRIAGE TRG-22-005 — PASS — `APPEND_CONFIRMED state=REPAIR_PENDING`
  (terminated `exit_code: 0`).
- Exact task-2 TASK-RECOVERY attempt 2 — PASS —
  `APPEND_CONFIRMED state=REPAIR_PENDING` (terminated `exit_code: 0`).
- Post-append standalone sequence command — PASS —
  `VALID state=REPAIR_PENDING` (terminated `exit_code: 0`).

Candidate bodies are the final three entries in the committed ledger. The
five authority scalars and exact allowed/forbidden/verification fields were
read directly from TRG-22-005 when preparing the candidates. An independent
byte comparison against a snapshot taken before any append printed
`PASS original ledger prefix preserved; exactly three canonical events appended`
(terminated `exit_code: 0`).

ESC-22-001 is resolved, with its original entry archived and TRG-22-006's
original BLOCK retained. No rejected BLOCK/STOP entry is fabricated.
`runner-handoff.md` now directs the next runner to the already recorded
task-2 recovery. The task-2 child has not yet been dispatched; its product
repair/gates, task-3 MF-3 and independent re-review remain phase work.

The next-runner handoff was also rehearsed against a fresh scratch copy of the
resolved real phase ledger plus its current decision/triage records. Appending
the replacement runner's normal RUN-RESUME through the same validator returned
`APPEND_CONFIRMED state=REPAIR_PENDING` (terminated `exit_code: 0`), preserving
the already bound recovery. This rehearsal did not append another event to
the real ledger.

Final product-scope check:
`git diff --exit-code 3e06f12fe4ae42e2694958043b95741bbb2c5c8a -- cmd internal`
— PASS — `exit 0 (no output)` (terminated `exit_code: 0`). The real ledger
again returned `VALID state=REPAIR_PENDING` (terminated `exit_code: 0`), and
final `git diff --check` returned `exit 0 (no output)` (terminated
`exit_code: 0`).

## Expected Non-Zero Commands

None among the required passing gates. Regression fixtures internally require
invalid appends to exit 1 and preserve the original ledger; the regression
script itself must exit 0. The red execution is retained as failed evidence.
