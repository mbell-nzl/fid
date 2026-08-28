# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**Phase 21 is active at task 12 with the ledger STOPPED after TRG-21-012;
DEC-167 resolves ESC-21-009 and authorizes the RUN-RESUME and same-task task-12
recovery route.** Tasks 1-11 and reviews 1/1b, 2/2b, and 3/3b are mechanically
complete. Phase 20
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
ESC-21-009 without a product-contract change.

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
- Dispatch state: **TASKS 1-11 COMPLETE; TASK 12 STOPPED AFTER TRG-21-012**. The
  Phase Runner appends a DEC-167 `RUN-RESUME` from that STOP, then one same-task
  `TASK-RECOVERY task-12` at attempt 6, and dispatches only the task-12
  continuation under the amended spec.
- Planning verification: **PASS** - deterministic coverage/inventory audit,
  `project-docs/tools/workflow-lint.sh`, and `git diff --check` all terminated
  `exit_code: 0` after the DEC-159 approval-state update.

## Carry-Forward

DEC-155 schedules CF-020/027 in task 8, CF-032 in task 12, and CF-024 in task
21. CF-008/018/023/026a/026b/028/029/030 are explicitly re-deferred; closed or
cancelled rows remain closed/cancelled.

## Next Work

1. Phase Runner appends the DEC-167 `RUN-RESUME` and the bound same-task
   `TASK-RECOVERY task-12` at attempt 6, then dispatches only task 12.
2. Developer preserves all failed evidence, converges the stale `internal/dev`
   claim precondition and the accumulated retained V1 rows in the build-tagged
   `cmd/fid/integration/` package to their accepted fail-closed and `--actor`
   contracts, records per-row provenance, leaves unrelated assertions
   byte-identical, reruns all six exact commands, and returns a truthful report.
   Expect several continuations; each later journey step reached is progress.
3. Phase Runner runs two-argument report lint, closes task 12, enters QA Review
   Window 4 for tasks 10-12, then the independent Group 1 checkpoint. The
   checkpoint reviewer treats the first green tagged closure journey as new
   evidence, not continuity (DEC-167 QA note).
