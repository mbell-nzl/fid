# Phase Runner session-resume repair — 2026-09-27

## Finding

The reported rejection is correct for the submitted event shape. The live
Phase 23 ledger validates as `VALID state=REPAIR_PENDING`. Its STOP was already
consumed by the DEC-256 RUN-RESUME at 2026-09-26 17:38. Tasks 6, 2, and 10's
repairs subsequently completed; TASK-RECOVERY task-11 at 19:19 remains active.
A replacement Runner needs a normal session resume for task 11.

Three instruction problems contributed to the stop:

1. Current-context summaries still instructed the Runner to resume the STOP
   and install the DEC-256 owners. They lagged the append-only ledger.
2. The role/spec/template said not to submit a different event after describing
   ambiguous append receipts, without clearly limiting that prohibition to the
   ambiguous case. Later text called an INVALID candidate recoverable.
3. The resume prompt told the Runner to stop whenever reports and ledger
   disagreed. An active recovery can legitimately have newer implementation
   commits than its unfinished report.

## Reproduction

Starting revision: `6e08357`. Working directory for all commands below:
`/home/mark/projects/fid`.

The original rejected candidate was unavailable. A reconstructed normal-session
candidate carrying the old TRG/ESC, DEC-256 authority, and Repair owners
reproduced both diagnostics against a fresh copy of the real ledger and its
adjacent decisions/triage files:

```text
INVALID normal schema v2 RUN-RESUME must use N/A for STOP-only fields
INVALID Repair owners are valid only on a RUN-RESUME after STOP
```

The nested append command terminated `exit_code: 1`, as expected for that
negative probe, and a byte comparison confirmed rejection preserved the copy.
The Python probe ran each command with `subprocess.run` and checked its actual
return code. With all six STOP-only fields set to `N/A` and Repair owners
omitted, the append returned `APPEND_CONFIRMED state=REPAIR_PENDING`
(terminated `exit_code: 0`). A byte-identical retry returned the same receipt
(terminated `exit_code: 0`) without changing the bytes. Standalone validation
of the resulting copy returned `VALID state=REPAIR_PENDING`
(terminated `exit_code: 0`). The probe also verified the live ledger was unchanged.

The durable regression is in `project-docs/tools/drills/review-block-owners.sh`:
restart during the second owner's active recovery, reject repeated STOP
authority and owners without mutation, accept the corrected normal resume,
confirm idempotent retry, reject a duplicate recovery and premature PASS,
then complete the existing recovery and close through a linked re-review.
This regression exercises existing behavior; the sequence verifier was unchanged.

## Changes

- Clarified current-state selection and rejected-candidate correction in the
  canonical Runner, workflow spec, templates, and operator playbook. Generated
  Claude and Codex Runner prompts were synchronized.
- Reconciled all four current-context summaries to the task-11 recovery.
  Archived the earlier progress snapshot in Phase 23 `progress.md`, clearly
  marked historical. Condensed active decisions while retaining the CF register
  and links to authoritative phase decisions.
- Preserved the user's `gpt-6-sol` selections for Developer, QA Reviewer, and
  TL-Lite from commit `acf1869c730cfd489221ee85cca92129129bfac6` in their canonical
  sources. Previously only generated Codex files held them, so sync reverted
  them. Their generated files now retain exactly the committed selections.

Scope: Human-delegated workflow maintenance, outside the active product task.
No product code, task command, VAL assertion, validator semantics, or live ledger
entry changed. No dependencies added. Product VAL coverage is not applicable.

## Verification

Captured in this maintenance session:

- `sh project-docs/tools/drills/review-block-owners.sh` — PASS —
  `REVIEW BLOCK OWNERS PASS all checks` (terminated `exit_code: 0`).
- `project-docs/tools/sync-agents.sh` — PASS —
  `Generated .codex/agents/tl_lite.toml` (terminated `exit_code: 0`).
- `project-docs/tools/drills/run-drills.sh` — PASS —
  `DRILLS PASS mechanical drills completed with no misses`
  (terminated `exit_code: 0`). This run includes the final exact duplicate-recovery
  diagnostic assertion and all existing mechanical drills.
- `project-docs/tools/workflow-lint.sh` — PASS, attempt 4 —
  `PASS workflow-lint completed with no failures` (terminated `exit_code: 0`).
- `project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-23/run-ledger.md`
  — PASS — `VALID state=REPAIR_PENDING` (terminated `exit_code: 0`).
- `git diff --check` — PASS — `exit 0 (no output)` (terminated `exit_code: 0`).

Retained failed attempts and causal repair:

- Workflow lint attempts 1 and 2 each ended with
  `FAIL workflow-lint found 2 failure(s)` (terminated `exit_code: 2`). Diagnostics:
  active decisions 252 lines against 200; current progress 199 against 133.
- Archived historical progress and condensed stale summaries. Attempt 3 ended
  with `FAIL workflow-lint found 1 failure(s)` (terminated `exit_code: 1`):
  active decisions still had 209 lines. Further condensation left 200 lines;
  current progress has 63. Attempt 4 reran the byte-identical command and passed.
- No gate definition, command, or line threshold was changed. Earlier failed
  attempts remain failed evidence. Historical product-test claims in the
  archived snapshot are not verification of the current task-11 repair.

## Runner handoff

Read the ledger first and run standalone sequence validation. At this checkpoint,
append a normal RUN-RESUME with `Continuing task: task-11`; set `Governing TRG`,
`Governing ESC`, `Trigger signature`, `Attempt count`, `Resolution authority`,
and `Next recovery` to `N/A`; omit `Repair owners`. Use a fresh timestamp and
concrete current HEAD/report/ledger reconciliation. Record the rejected append's
diagnostics and observed exit in that reconciliation evidence.

Continue the existing `review-7:qa-rework:task-11-BC-3` attempt 1. Implementation
commit `466745f2dc4505d60eb67ac50193e004532b7239` is present; the current report
still names the earlier `c5bf0f8` implementation. Reconcile any surviving child
session, finish the exact task-11 verification and report, and record its real
outcome before completion. Do not duplicate TASK-RECOVERY, replay the STOP, or
rerun completed owners. After task 11 closes, obtain the fresh independent
`checkpoint-after-task-11-3.md` re-review linked to review-7. Task 12 remains
held until that review passes.
