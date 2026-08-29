# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**Phase 21 is active in Group 2 with the ledger STOPPED after TRG-21-018; QA
Review Window 5 returned REPAIR on two bounded task-15 Must-fixes, and DEC-170
resolves ESC-21-012 by correcting the recovery category and authorizing the
task-15 attempt-6 continuation.** Tasks 1-15, reviews 1/1b, 2/2b, 3/3b, 4/4b, and
the Group 1 checkpoint rerun `infra-review-1b` are mechanically complete. Phase 20
completed/wrapped on 2026-08-22. The Phase 20 retrospective is committed in
`../reviews/2026-08-22-phase-20-retrospective.md`; Gate D3.2 is closed by
DEC-153; the initial slice/spec contract is DEC-154; carry-forward disposition
is DEC-155; the terminal gate amendment is DEC-156; and the complete interface,
bootstrap, inventory, sequencing, and validation correction is DEC-157.
QA-plan refinements and the unchanged deterministic rebuild target are DEC-158.
QA readiness and Human execution approval are DEC-159. DEC-160 preserves the
role-free bootstrap while bounding temporary V1 policy compatibility; DEC-161
selects the valid ledger recovery record without a verifier change; DEC-162
converges retained V1 tests inside the task whose accepted change breaks them;
DEC-163 makes a QA `REPAIR` review carry an ordered owner list and authorizes
the matching `workflow-sequence-lint.sh` correction with new drills. DEC-164
keeps ordinary schema-v2 repair authority byte-exact and locates its governing
TRG reference in TASK-RECOVERY `Source`; its bound same-task recovery replaces
the pending dispatch whose blocker candidate the TRIAGE classified. DEC-165
states that retained-test convergence ownership follows the invalidating
accepted change rather than the fixture file's primary subject, and closes
ESC-21-007 with a bounded task-9 scope amendment and no task-10 work. DEC-166
establishes that historical task labels do not determine current-phase test
ownership, resolves the stale Phase 20 request fixture as task-9 convergence,
and closes ESC-21-008 without a product-contract change. DEC-167 requires a
task's required commands to exercise every file scope it owns, assigns task 12
its own `internal/dev` collateral plus the accumulated build-tagged
`cmd/fid/integration/` convergence with per-row provenance, and closes
ESC-21-009 without a product-contract change. DEC-168 records that no budget
bounds a QA Lead child session's lifetime - 120s bounds each FEAS probe, 900s the
FEAS set, 300s the tagged integration command - reruns the unchanged, unsplit
Group 1 checkpoint as attempt 4, and closes ESC-21-010 without weakening
VAL-042, VAL-043, FEAS-01..FEAS-19, or any budget. That checkpoint then passed
VAL-042, VAL-043, and FEAS-01..FEAS-19 and returned REPAIR on two bounded
findings. DEC-169 keeps both owners but resizes both repairs: task 4 fixes
same-key idempotent replay across the whole Activation family, and task 9
implements the complete `public-interface.md` section 7 category-to-exit mapping
rather than one code, because `cmd/fid/main.go` currently exits 1 for every
failure. Both tasks gain `cmd/fid/integration/` scope and the tagged integration
gate; no VAL, QA plan, checkpoint contract, or public contract changes. Both
repairs and the complete checkpoint rerun then passed, and tasks 13-15 completed.
DEC-170 corrects the Window 5 task-15 escalation: all five attempts were shut
down by the orchestrator while running, so the NO-PROGRESS verdicts measured
interruption rather than repair capability. An externally interrupted attempt
never consumes the no-progress budget, no budget bounds a Developer child
session's lifetime, durable work is committed before the long commands, and
MF-15-02's provenance is satisfied forward-only rather than by rewriting commits
at or before `09a2a7b`. No finding, command, VAL, QA artifact, or contract
changes. DEC-171 makes that standing rule executable in the Phase Runner
prompts, session router, ledger verifier, workflow lint, and drill 19 without
changing DEC-170's attempt-6 route.

## Planned Boundary

One Heavy Phase 21 uses 22 bounded Developer tasks in two internal groups. An
independent checkpoint after task 12 must prove Group 1 contract closure and
fresh Group 2 executor/verifier feasibility before task 13. The validation
contract contains VAL-001..087. Seven three-task review windows, Phase Runner
two-argument post-commit report lint before each next dispatch, final QA, and
one Human terminal-output parent assertion with four Human subchecks apply.

## Approval State

- Human approval of revised Heavy mode, exact contracts, VAL contract, and task
  pack: **APPROVED 2026-08-23 / DEC-159**.
- QA Lead Heavy QA plan with fresh feasibility: **READY** at `e8a636e`; no
  outstanding Must-fix.
- Product implementation: **AUTHORIZED UNDER THE ORDERED TASK PACK**.
- Dispatch state: **TASKS 1-15 COMPLETE; WINDOW 5 REPAIR, STOPPED AFTER
  TRG-21-018**. Under DEC-170 the Phase Runner appends `RUN-RESUME` retaining
  TRG-21-018/ESC-21-012 with resolution authority `DEC-170`, then one
  `TASK-RECOVERY task-15`, `session-continuation`, `task-scope`, attempt 6,
  consecutive no-progress `0/2`. Task 16 waits for a fresh independent
  `review-5b` PASS.
- Planning verification: **PASS** - deterministic coverage/inventory audit,
  `project-docs/tools/workflow-lint.sh`, and `git diff --check` all terminated
  `exit_code: 0` after the DEC-159 approval-state update.

## Carry-Forward

DEC-155 schedules CF-020/027 in task 8, CF-032 in task 12, and CF-024 in task
21. CF-008/018/023/026a/026b/028/029/030 are explicitly re-deferred; closed or
cancelled rows remain closed/cancelled.

## Next Work

1. Phase Runner appends the DEC-170 item 9 sequence: `RUN-RESUME` bound to
   TRG-21-018/ESC-21-012 at attempt 6, then `TASK-RECOVERY task-15`
   (`session-continuation`, `task-scope`, attempt 6, no-progress `0/2`), and
   waits for the Developer child to end on its own.
2. Task 15 commits the existing durable repair with its `Refs: phase-21/task-15`
   footer before any long command, then runs the five unchanged commands under
   the DEC-170 per-command budgets, appending each observed termination to the
   report as it terminates. The unverified attempt-5 rows are rerun in session,
   and the report's provenance deviation names `79e3a93` and `8491cde` instead of
   claiming an amendment; no commit at or before `09a2a7b` is rewritten.
3. `REPAIR task-15`, then `TASK task-15 COMPLETE` with two-argument report-lint,
   then a fresh independent QA Reviewer writes the linked Window 5 re-review and
   the Runner records `REVIEW review-5b` (`Re-review of: review-5`). Task 16
   dispatches only on its PASS, followed by Window 6 for tasks 16-18.
