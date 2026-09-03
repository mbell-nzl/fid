# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**Phase 21 is active in Group 2 with the ledger STOPPED after TRG-21-022 /
ESC-21-016 at task 22; DEC-175 resolves it as a bounded product repair, giving
task 22 the causal `internal/workflow/activation.go` Event-name fix its own gate
is the only one to observe.** Tasks 1-21, reviews 1/1b through 7b, and the Group 1
checkpoint are mechanically complete. Phase 20 completed/wrapped on 2026-08-22;
its retrospective is in
`../reviews/2026-08-22-phase-20-retrospective.md`.

Phase 21 planning decisions: DEC-153 (Gate D3.2), DEC-154 (slice/spec
contract), DEC-155 (carry-forward), DEC-156 (terminal gate), DEC-157 (complete
interface, bootstrap, inventory, sequencing), DEC-158 (QA refinements, rebuild
target), DEC-159 (QA readiness and Human approval).

Execution decisions: DEC-160 bounds the temporary V1 policy bridge; DEC-161
selects the valid ledger recovery record; DEC-162 converges retained V1 tests
inside the invalidating task; DEC-163 gives QA `REPAIR` an ordered owner list;
DEC-164 keeps schema-v2 repair authority byte-exact; DEC-165/166 make
convergence ownership follow the invalidating change, never a file subject or
historical label; DEC-167 requires a task's gates to exercise every scope it
owns; DEC-168 removes lifetime budgets from QA child sessions; DEC-169 resizes
the checkpoint repairs; DEC-170/171 make externally interrupted attempts
non-consuming and executable in the Runner controls; DEC-172 keeps the
unavailable `attention.v1` row and replaces the zero-match search with a
boundary test.

DEC-173 (2026-09-02) closed ESC-21-014 by amending task 20's scope to the reach
of its own absence controls and adding the tagged integration gate under
DEC-167; `08_QUALITY_PLAN.md` §3 records the converse executable-ownership rule.
Task 20 and task 21 then completed and review window 7/7b passed.

DEC-174 (2026-09-03) closes ESC-21-015. Task 22's real-process journey reached
Gate execution and Receipt Evidence evaluation, then the terminal transition was
correctly rejected with `transition-not-permitted` /
`workflow-transition-invalid` at exit 4. The product is right and the planning
artifact was wrong: `contracts.md` §16 fixed an exact 20-Event prefix with a
direct terminal edge, but `task-default` needs four edges to reach a `completed`
state, `subtask-default` needs a second `issue.created` for its required task
parent, and `epic-default` adds a rule-guarded intermediate state. The journey
now authors its own workflow with the already-accepted `workflow define
--applies-to task` and `workflow assign --kind task` leaves before `issue
create`, and the exact prefix is 22 Events with VAL-086 amended to match. The
observed rejection becomes required VAL-085 negative coverage. No product,
command, public-interface, Event-schema, transition, or frozen-workflow change.
TRG-21-021, ESC-21-015, the committed corpus and journey guard, and all failed
evidence remain historical.

DEC-175 (2026-09-03) closes ESC-21-016. The DEC-174 journey was implemented and
committed at `a6b3a44` (attempt 2, `MATERIAL-PROGRESS`), and the tagged gate then
narrowed to one product defect: `internal/workflow/activation.go:519` derives the
Event type as `operation + "ted"` and overrides only commit and abort, so
finalize emits `activation.finalizeted` where `public-interface.md` §4.1,
`contracts.md` §14, and `contracts.md` §16 all require `activation.finalized`.
TRG-21-022's BLOCK is sustained on its category. The repair belongs to task 22,
not a reopened task 4, because task 22's gates are the only ones that execute the
accepted Event name (DEC-167, DEC-173); its scope gains the bounded derivation
plus a non-tagged regression assertion under an already-required command. No
contract, schema, interface, command, corpus, or VAL change. `08_QUALITY_PLAN.md`
§3 records the first-observing ownership rule for latent production defects.

## Planned Boundary

One Heavy Phase 21 uses 22 bounded Developer tasks in two internal groups. An
independent checkpoint after task 12 proved Group 1 contract closure and fresh
Group 2 executor/verifier feasibility before task 13. The validation contract
contains VAL-001..087. Seven three-task review windows, Phase Runner
two-argument post-commit report lint before each next dispatch, final QA, and
one Human terminal-output parent assertion with four Human subchecks apply.

## Approval State

- Human approval of revised Heavy mode, exact contracts, VAL contract, and task
  pack: **APPROVED 2026-08-23 / DEC-159**.
- QA Lead Heavy QA plan with fresh feasibility: **READY** at `e8a636e`; no
  outstanding Must-fix.
- Product implementation: **AUTHORIZED UNDER THE ORDERED TASK PACK**.
- Dispatch state: **TASKS 1-21 COMPLETE; TASK 22 STOPPED AFTER TRG-21-022 /
  ESC-21-016 WITH RECOVERY AUTHORIZED BY DEC-175**. The Phase Runner appends
  `RUN-RESUME` retaining those IDs, the exact recorded trigger, attempt 3, and
  resolution authority `DEC-175`, then one `TASK-RECOVERY task-22` with
  `in-scope-implementation`, `task-scope`, and no-progress `0/2`.
- Planning verification: **PASS** - deterministic coverage/inventory audit,
  `project-docs/tools/workflow-lint.sh`, and `git diff --check` all terminated
  `exit_code: 0` after the DEC-159 approval-state update.

## Carry-Forward

DEC-155 schedules CF-020/027 in task 8, CF-032 in task 12, and CF-024 in task
21. CF-008/018/023/026a/026b/028/029/030 are explicitly re-deferred; closed or
cancelled rows remain closed/cancelled.

## Next Work

1. Phase Runner appends the DEC-175 recovery sequence: `RUN-RESUME` naming the
   2026-09-03 16:57 STOP and bound to TRG-21-022 / ESC-21-016 with the recorded
   trigger, attempt 3, resolution authority `DEC-175`, and the exact
   next-recovery tuple from DEC-175 item 6; it then appends the matching
   `TASK-RECOVERY task-22` with no-progress `0/2`.
2. The Developer continues task 22 from `a6b3a44` under the amended spec: keeps
   the committed corpus, journey, and negative matrix, writes the
   `cmd/fid/activation_test.go` Event-type regression assertion red, makes the
   bounded `internal/workflow/activation.go` derivation repair, completes the
   post-rebuild repeat, preserves the attempt-1 and attempt-2 failed evidence,
   commits durable work before long commands, and runs the complete
   byte-identical eight-command block.
3. After the report commit and two-argument report-lint PASS, the Runner records
   `REPAIR task-22 RESOLVED` and `TASK task-22 COMPLETE`, then final QA owns
   task-22 evidence, VAL-087, and the Human checklist.
