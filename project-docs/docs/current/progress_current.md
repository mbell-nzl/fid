# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**Phase 21 is active at the Group 1 independent checkpoint with the ledger
STOPPED after TRG-21-014; DEC-168 resolves ESC-21-010 and authorizes one fresh
checkpoint attempt under a one-off task-13 orchestration wrapper.** Tasks 1-12
and reviews 1/1b, 2/2b, 3/3b, and 4/4b are mechanically complete. Phase 20
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
VAL-042, VAL-043, FEAS-01..FEAS-19, or any budget.

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
- Dispatch state: **TASKS 1-12 COMPLETE; CHECKPOINT STOPPED AFTER TRG-21-014**.
  The Phase Runner appends a DEC-168 `RUN-RESUME` from that STOP, then one
  `TASK-RECOVERY task-13` at attempt 4 with category `session-continuation` and
  mode `task-scope`, and dispatches only a fresh independent QA Lead checkpoint
  session. Task 13 Developer dispatch waits for a verified checkpoint PASS.
- Planning verification: **PASS** - deterministic coverage/inventory audit,
  `project-docs/tools/workflow-lint.sh`, and `git diff --check` all terminated
  `exit_code: 0` after the DEC-159 approval-state update.

## Carry-Forward

DEC-155 schedules CF-020/027 in task 8, CF-032 in task 12, and CF-024 in task
21. CF-008/018/023/026a/026b/028/029/030 are explicitly re-deferred; closed or
cancelled rows remain closed/cancelled.

## Next Work

1. Phase Runner appends the DEC-168 `RUN-RESUME` and the bound
   `TASK-RECOVERY task-13` at attempt 4, then dispatches only a fresh independent
   QA Lead Group 1 checkpoint session and waits for that child to terminate on
   its own. It must not interrupt the session at 900 seconds.
2. QA Lead runs the unchanged checkpoint: separate VAL-042 and VAL-043 decisions,
   FEAS-01..FEAS-19 fresh in newly allocated scratch with observed termination,
   and `qa/group-1-checkpoint-review.md` written incrementally during the
   attempt. Preserved scratch and any diagnostic run are inadmissible. The
   reviewer treats the first green tagged closure journey as new evidence, not
   continuity (DEC-167 QA note).
3. On a verified committed checkpoint `PASS`, Phase Runner dispatches task 13
   under its unchanged spec and seven commands, runs two-argument report lint,
   then appends the DEC-168 item 9 `REPAIR task-13` / `TASK task-13 COMPLETE`
   closure and enters QA Review Window 5 for tasks 13-15. Any non-PASS verdict
   follows DEC-168 item 10 instead.
