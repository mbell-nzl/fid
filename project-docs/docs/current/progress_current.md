# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**Phase 21 is active in Group 2 with the ledger STOPPED after TRG-21-019 /
ESC-21-013 at task 17; DEC-172 resolves the protected contradiction between
VAL-063's broad absence recipe and the required unavailable `attention.v1`
capability row.** Tasks 1-16, reviews 1/1b, 2/2b, 3/3b, 4/4b, 5/5b, and the Group
1 checkpoint rerun `infra-review-1b` are mechanically complete. Phase 20
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
changing DEC-170's attempt-6 route. Task 15 then completed with a fresh
`review-5b` PASS, and task 16 completed. DEC-172 retains the exact public
unavailable capability row and replaces the defective zero-match search with a
named boundary test that requires that row while rejecting any Attention state,
effect, CLI leaf, or additional production match. The prior failed command and
STOP remain historical evidence.

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
- Dispatch state: **TASKS 1-16 COMPLETE; TASK 17 STOPPED AFTER TRG-21-019 /
  ESC-21-013 WITH RECOVERY AUTHORIZED BY DEC-172**. The Phase Runner appends
  `RUN-RESUME` retaining those IDs, the exact recorded trigger, attempt 3, and
  resolution authority `DEC-172`, then one `TASK-RECOVERY task-17` with
  `gate-repair`, `task-scope`, and no-progress `0/2`.
- Planning verification: **PASS** - deterministic coverage/inventory audit,
  `project-docs/tools/workflow-lint.sh`, and `git diff --check` all terminated
  `exit_code: 0` after the DEC-159 approval-state update.

## Carry-Forward

DEC-155 schedules CF-020/027 in task 8, CF-032 in task 12, and CF-024 in task
21. CF-008/018/023/026a/026b/028/029/030 are explicitly re-deferred; closed or
cancelled rows remain closed/cancelled.

## Next Work

1. Phase Runner appends the DEC-172 recovery sequence: `RUN-RESUME` bound to
   TRG-21-019 / ESC-21-013 with the recorded trigger, attempt 3, resolution
   authority `DEC-172`, and the exact next-recovery tuple from DEC-172; it then
   appends the matching `TASK-RECOVERY task-17` with no-progress `0/2`.
2. The Developer adds `TestPhase21AttentionBoundary`, commits the preserved
   task-17 implementation and test with `Refs: phase-21/task-17` before long
   commands, and runs the complete revised six-command block. The old failed
   `rg` result remains in the report as historical evidence.
3. After the report commit and two-argument report-lint PASS, the Runner records
   `REPAIR task-17 RESOLVED` and `TASK task-17 COMPLETE`, then dispatches task
   18. Window 6 review remains gated on completion of tasks 16-18.
