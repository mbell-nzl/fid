# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None.

## Latest Resolution

DEC-172 resolved ESC-21-013 on 2026-08-30. TRG-21-019 correctly stopped task 17:
the broad Attention absence command could not pass while the accepted
capability contract required the public unavailable `attention.v1` row. The
row remains exact and public. It is negative capability metadata, not Attention
state or an Effect.

The defective search is replaced by `TestPhase21AttentionBoundary`, which must
require the false row as the sole production Attention match and reject any
additional table, Event, Effect, CLI leaf, state, operation, or mutation. The
old failed command remains historical evidence. Recovery is RUN-RESUME under
TRG-21-019 / ESC-21-013 / DEC-172, task-17 attempt 3, followed by bounded
`gate-repair`; task 18 remains gated on task-17 completion.

## Archive Pointers

- Phase 20: `../phases/phase-20/tl-escalations.md`
- Phase 21: `../phases/phase-21/tl-escalations.md`
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
