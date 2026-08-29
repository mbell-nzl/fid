# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None.

## Latest Resolution

DEC-170 resolved ESC-21-012 on 2026-08-30. TRG-21-017 and TRG-21-018 correctly
preserved the two task-15 Must-fixes, ownership, and failed evidence, but their
`recovery-budget-exhausted` category was not sustained. All five Developer
children were externally shut down before returning a typed result, while the
durable working tree contains the causal repair. The corrected category is
orchestration session-liveness.

The authorized recovery remains one `RUN-RESUME` retaining TRG-21-018 and
ESC-21-012 with `Resolution authority: DEC-170`, followed by task-15
`TASK-RECOVERY` attempt 6, `session-continuation`, `task-scope`, and
no-progress `0/2`. Attempt 6 first commits the existing repair with
`Refs: phase-21/task-15`, reruns the five unchanged commands under their
per-command budgets, records each observed termination, corrects the false
history-rewrite claim, and obtains fresh independent `review-5b`. Task 16 waits
for that review's PASS.

DEC-171 implements DEC-170's standing rule in the workflow controls. A parent
polling wait timeout retains the live child; a real external interruption is
typed `INTERRUPTED` evidence that preserves the prior no-progress count; two
consecutive interruptions require orchestration session-liveness triage.

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
