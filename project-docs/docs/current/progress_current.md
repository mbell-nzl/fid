# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**Phase 21 is active with all 22 tasks complete and is ready for final QA.**
Reviews 1/1b through 7b and the Group 1 checkpoint are complete. Task 22's
attempt-6 recovery, eight-command gate, committed report, authoritative
two-argument report lint, `REPAIR RESOLVED`, and `TASK COMPLETE` are recorded;
the ledger state is `ADVANCEABLE`. Phase 20 completed/wrapped on 2026-08-22;
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

DEC-176 (2026-09-03) preserves the immutable TRG-21-022 attempt-3 TRIAGE and
corrects only the next recovery identity: schema v2 requires the STOP resume to
increment its governing TRIAGE, so Task 22 resumes at attempt 4 with no-progress
`0/2`. DEC-175 remains the product/scope authority; no command, contract,
assertion, budget, or evidence changes.

DEC-177 (2026-09-03) closes ESC-21-017. Task 22 owns the minimum
`internal/evidence`, `cmd/fid`, migration, golden, schema, and test scope needed
for accepted read-only Receipt get/list and durable canonical
`rejected_artifacts`. `contracts.md` §10 changes only to the accepted root
`fid resume` spelling. The two attempt-4 tagged failures remain immutable; the
attempt-6 recovery implemented the accepted surface without changing its wire
contract, Event budget, VALs, or gates. Product commit `ac6270a` and report
commit `5d8204c` close the task with all prior failed evidence preserved.

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
- Dispatch state: **TASKS 1-22 COMPLETE; FINAL QA AUTHORIZED**. No task
  recovery or escalation remains open. The ledger validates as `ADVANCEABLE`.
- Planning verification: **PASS** - deterministic coverage/inventory audit,
  `project-docs/tools/workflow-lint.sh`, and `git diff --check` all terminated
  `exit_code: 0` after the DEC-159 approval-state update.

## Carry-Forward

DEC-155 schedules CF-020/027 in task 8, CF-032 in task 12, and CF-024 in task
21. CF-008/018/023/026a/026b/028/029/030 are explicitly re-deferred; closed or
cancelled rows remain closed/cancelled.

## Next Work

1. QA Lead executes `qa/qa-plan.md` QA-81..93 in a fresh final-QA session,
   writes `qa/final-review.md`, and records every VAL-001..087 outcome with
   session-captured evidence.
2. The Human completes `qa/human-checklist.md` HUMAN-01a..01d for VAL-036; QA
   must not pre-answer or infer these checks.
3. Phase Runner records the final `REVIEW` outcome. Only a passing review and
   completed Human checklist permit `RUN-END` and Tech Lead phase wrap-up.
