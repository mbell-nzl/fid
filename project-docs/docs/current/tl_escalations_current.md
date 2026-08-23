# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None.

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
