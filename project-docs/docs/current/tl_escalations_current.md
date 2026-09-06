# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None. ESC-22-002 was resolved on 2026-09-06 under DEC-192/193 and is archived
in `../phases/phase-22/tl-escalations.md`.

## Latest Resolution

ESC-22-002 resolved on 2026-09-06 under DEC-192 and DEC-193. Task 4's file
scope is amended to delete exactly the `run.checkpointed` key from the task-1
progressive fence in `internal/journal/project_effects.go` and to add a
task-4-owned projection test; DEC-192 pre-authorizes the same key removal for
tasks 5, 7, 8 and 10 and schedules the displaced VAL-003 convergence. The
retained fixed-date Stop fixture is DEC-189 nondeterminism, assigned to the
observing task and repaired at its source, with task 2 keeping subject
ownership. No accepted contract, gate, VAL row, or task allocation changed and
no Human decision was required. ESC-22-001 remains resolved under DEC-190/191.
Both escalations and their resolutions are archived in
`../phases/phase-22/tl-escalations.md`.

## Prior Phase Resolution

Phase 21 closed on 2026-09-04 (DEC-179) with every allocated escalation
resolved. All seventeen Phase 21 escalations and their resolutions are archived
in `../phases/phase-21/tl-escalations.md`; ESC-21-017 (DEC-177, missing accepted
Receipt read leaves) was the last one raised and closed.

Wrap-up note carried forward as a control weakness, not an open item:
ESC-21-007 through ESC-21-011 were allocated in the Phase 21 triage log and
resolved by DEC-165..DEC-169, but were never written into this register while
they were open. They are now archived from their governing decisions. Register
an escalation here at allocation time, not at resolution time.

## Archive Pointers

- Phase 22: `../phases/phase-22/tl-escalations.md`
- Phase 21: `../phases/phase-21/tl-escalations.md`
- Phase 20: `../phases/phase-20/tl-escalations.md`
- Phase 19: `../phases/phase-19/tl-escalations.md`
- Earlier phases: their `../phases/phase-N/tl-escalations.md` archives.

## Protected Escalation Triggers

Escalate only for a required product, architecture, data-model,
public-interface, dependency, security/data-loss, specification, exhausted
recovery-budget, or unavailable external-capability decision. Ordinary red
tests, failed gates, report defects, QA rework, and incomplete in-scope work
remain bounded repair. Before claiming recovery exhaustion, confirm each
counted attempt returned a typed result and its directive permitted causal
repair. Externally interrupted attempts are session-liveness evidence, not
repair-capability evidence.
