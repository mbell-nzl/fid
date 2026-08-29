# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**Phase 21 is active at the Group 1 independent checkpoint with the ledger
STOPPED after TRG-21-015; the checkpoint ran complete and returned REPAIR, and
DEC-169 resolves ESC-21-011 by routing MF-1 to task 4 and MF-2 to task 9.** Tasks
1-12 and reviews 1/1b, 2/2b, 3/3b, and 4/4b are mechanically complete. Phase 20
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
gate; no VAL, QA plan, checkpoint contract, or public contract changes.

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
- Dispatch state: **TASKS 1-12 COMPLETE; CHECKPOINT REPAIR, STOPPED AFTER
  TRG-21-015**. Under DEC-169 the Phase Runner appends `REVIEW infra-review-1`
  (REPAIR, owners `task-4, task-9`), `TRIAGE TRG-21-016`, `STOP`, and a
  `RUN-RESUME` continuing task-4 at attempt 2, then dispatches task 4 and task 9
  in that order. Task 13 waits for a PASS on the complete checkpoint rerun.
- Planning verification: **PASS** - deterministic coverage/inventory audit,
  `project-docs/tools/workflow-lint.sh`, and `git diff --check` all terminated
  `exit_code: 0` after the DEC-159 approval-state update.

## Carry-Forward

DEC-155 schedules CF-020/027 in task 8, CF-032 in task 12, and CF-024 in task
21. CF-008/018/023/026a/026b/028/029/030 are explicitly re-deferred; closed or
cancelled rows remain closed/cancelled.

## Next Work

1. Phase Runner appends the DEC-169 item 8 sequence: `REVIEW infra-review-1`,
   `TRIAGE TRG-21-016` (TL-Lite allocates the identity-only successor of
   TRG-21-015, retaining ESC-21-011), `STOP`, and the `RUN-RESUME` continuing
   task-4 at attempt 2 with its bound `TASK-RECOVERY task-4`.
2. Task 4 repairs MF-1: no lifecycle precondition preempts an exact
   same-key/same-digest replay anywhere in the Activation family; conflicts stay
   byte-identical; every error path emits its typed envelope. Then task 9 repairs
   MF-2: the complete section 7 category-to-exit mapping, typed error, `main.go`
   wiring, unknown-command exit 2, and every `cli.WriteFailure` site converted,
   with no envelope change. Both run the tagged integration gate and return a
   committed report with fresh two-argument report-lint.
3. A fresh independent reviewer repeats the complete unsplit checkpoint with all
   FEAS probes fresh, appends the second review to
   `qa/group-1-checkpoint-review.md`, and the Runner records `REVIEW
   infra-review-1b` PASS. Only then does task 13 dispatch, followed by QA Review
   Window 5 for tasks 13-15. The reviewer still treats the first green tagged
   closure journey as new evidence, not continuity (DEC-167 QA note).
