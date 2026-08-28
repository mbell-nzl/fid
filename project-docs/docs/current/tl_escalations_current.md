# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None.

---

DEC-168 resolved ESC-21-010 on 2026-08-29. TRG-21-013 and TRG-21-014 are not
sustained on their recorded category. The blocking condition was real and the
refusal to dispatch task 13 on partial evidence was correct, but the QA
capability is available: preserved scratch `/tmp/tmp.zQAE0hBrd2` holds a built
`child-bin` and materialized outputs for essentially the whole probe set, and a
fresh diagnostic tagged integration run outside any checkpoint session completed.
The cause is orchestration applying the 900-second FEAS-set command budget as the
QA Lead child session's total lifetime and restarting each attempt from zero.
DEC-168 records the binding reading - 120s per probe, 900s for the FEAS set, 300s
independently for the tagged integration command, and no budget on the child
session's lifetime - with no edit to the QA plan or checkpoint contract. VAL-042,
VAL-043, FEAS-01..FEAS-19, fresh-session evidence, observed termination,
commands, accepted statuses, budgets, and reviewer independence are unchanged,
and the checkpoint is not split. Preserved scratch and the diagnostic run are
inadmissible as checkpoint evidence. One fresh independent attempt is authorized:
the Phase Runner appends RUN-RESUME then one `TASK-RECOVERY task-13`, category
`session-continuation`, mode `task-scope`, attempt 4, waits for the QA Lead child
to finish without interrupting it at 900 seconds, and dispatches the task-13
Developer only after the committed checkpoint artifact records a verified PASS.
The general phase-scoped trigger identity and interrupted-child ledger event are
deferred to workflow maintenance.

DEC-167 resolved ESC-21-009 on 2026-08-28. TRG-21-012 is not sustained. Both
owner routings are wrong and the recovery-budget category does not hold. The
`internal/dev` failure is task 12's own collateral from removing the V1 Claim
precondition, not task-20 work; `v1-convergence-inventory.md:54` assigns task 20
only the source-package removal. The `issue ac verify` failure is not a `--by`
transport gap: VAL-030 makes the leaf unconditionally `capability-unavailable`,
so no flag change makes the step succeed, and `concern raise`/`address` are the
same fail-closed family and fail next. The budget was consumed by run-ledger
:1904, which forbade task 12 from repairing the only failing rows. Root cause is
gate coverage: `cmd/fid/integration/` is build-tagged and only tasks 12, 18, 21,
and 22 run it, so DEC-165 gave task 9 ownership it could not execute. DEC-167
amends task 12's scope to name `internal/dev/` and `cmd/fid/integration/`,
assigns it the accumulated convergence with per-row provenance, and adds two
standing rules to `08_QUALITY_PLAN.md` §3: required commands must exercise owned
scope, and a no-progress budget counts only attempts whose directive permitted
causal repair. Phase Runner appends RUN-RESUME then one same-task
`TASK-RECOVERY task-12` at attempt 6.

DEC-166 resolved ESC-21-008 on 2026-08-26. The accepted contracts do not
conflict. `cmd/fid/task10_contract_test.go` and its actor-request assertions
were introduced by Phase 20 task 10; current Phase 21 task 10 is generic Claims
and owns none of those inventory leaves. The filename-based ownership inference
was wrong, and repeating it did not create a protected category. Task 9 owns
subject-preserving convergence to the accepted `--actor` plus idempotency and
current request-validation contract. The Phase Runner resumes TRG-21-010 at
attempt 12 under DEC-166; no `--by`, old envelope, production behavior, gate,
VAL, or current task-10 semantics are restored or changed.

DEC-165 resolved ESC-21-007 on 2026-08-25. The escalation's premise did not hold:
the remaining red fixtures are not Task 10 Claim/Run authority fixtures. No
generic Claim or Run fixture exists in the tree and task 10 has not been
dispatched; every red assertion is an incidental read of a surface task 9 itself
removed or renamed - public `fid policy show`, stats `--by role`, and fail-closed
`maintenance vacuum`. DEC-165 states DEC-162's rule explicitly - convergence
ownership follows the invalidating accepted change, not the fixture file's
primary subject - amends task 9 to name `cmd/fid/integration/` and narrow Out Of
Scope to generic Claim/Run implementation and semantics, and requires the V1
claim/run/dev behavioural assertions to stay byte-identical. No re-slicing,
sequencing change, gate change, or task-10 work. Phase Runner appends RUN-RESUME
then one same-task `TASK-RECOVERY task-9` at attempt 9.

DEC-164 resolved ESC-21-006 on 2026-08-24. Ordinary schema-v2 TRIAGE recovery
names its governing TRG in `Source` and preserves the external `Repair
authority` byte-for-byte; the protected verifier correction and regression
drills also let that bound recovery replace the pending same-task dispatch.
The ledger remains `REPAIR_PENDING` after TRG-21-007, so Phase Runner resumes
by appending that bound task-8 recovery, not by appending RUN-RESUME or binding
TRG-21-008.

DEC-163 resolved ESC-21-005 on 2026-08-23. A QA `REPAIR` review carries an
ordered comma-separated owner list in exactly one `Owning task` field, every
listed owner completes TASK-RECOVERY/REPAIR/TASK complete before a linked
re-review may close or replace the review, and `CORRECTION` stays evidence-only.
DEC-163 also authorizes the corresponding `workflow-sequence-lint.sh` correction
and its drills as workflow infrastructure outside the phase recovery loop. That
change must land before the Phase Runner appends the DEC-163 RUN-RESUME.

DEC-160, DEC-161, and DEC-162 resolved ESC-21-001..004 on 2026-08-23. The full
conditions, evidence, and recovery dispositions are archived in
`../phases/phase-21/tl-escalations.md`.

DEC-162 additionally records a standing rule: retained V1 tests are not
protected contracts, and the task whose accepted change invalidates one adapts
it inside that task, preserving subject and assertion strength. Tasks 8, 9, and
20 inherit it through `../phases/phase-21/v1-convergence-inventory.md` section
1.2, so equivalent retained-test collateral is not an escalation trigger.

Phase 21 QA blocker candidates BC-01 and BC-02 were planning-contract defects,
not open Human escalations. DEC-158 resolved them without changing the accepted
rebuild threshold or product slice; the Ready QA re-review and Human execution
approval completed under DEC-159.

Phase 20 completed on 2026-08-22. Its allocated escalations are resolved and
archived in `../phases/phase-20/tl-escalations.md`; DEC-151 records the final
triage and escalation audit.

## Archive Pointers

- Phase 20: `../phases/phase-20/tl-escalations.md`
- Phase 21 in-flight resolved items: `../phases/phase-21/tl-escalations.md`
- Phase 19: `../phases/phase-19/tl-escalations.md`
- Phase 18: `../phases/phase-18/tl-escalations.md`
- Phase 17: `../phases/phase-17/tl-escalations.md`

## Protected Escalation Triggers

Escalate only for a required product, architecture, data-model,
public-interface, dependency, security/data-loss, specification, exhausted
recovery-budget, or unavailable external-capability decision. Ordinary red
tests, failed gates, report defects, QA rework, and incomplete in-scope work
remain bounded repair.
