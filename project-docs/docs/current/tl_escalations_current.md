# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None. DEC-160 and DEC-161 resolved ESC-21-001..003 on 2026-08-23. The full
conditions, evidence, and recovery dispositions are archived in
`../phases/phase-21/tl-escalations.md`; the stopped ledger remains unchanged
until the Phase Runner appends DEC-161's exact RUN-RESUME.

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
