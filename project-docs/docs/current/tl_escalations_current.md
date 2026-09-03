# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None. ESC-21-016 was resolved by DEC-175 on 2026-09-03; ESC-21-015 was resolved
by DEC-174 on 2026-09-03.

## Latest Resolution

DEC-175 resolved ESC-21-016 on 2026-09-03. TRG-21-022's BLOCK is sustained on
its recorded category: the Developer reproduced a deterministic failure,
refused to change the Event contract, assertion, command, or test to make the
gate pass, and correctly reported the causal writer as outside its named source
scope, and the tools-only control-conformance route genuinely cannot reach a
production workflow path.

The cause is a product defect, not a contract question.
`internal/workflow/activation.go:519` builds the Event type as
`operation + "ted"` and overrides only `activation.commit` and
`activation.abort`, so `activation.finalize` emits `activation.finalizeted`.
`public-interface.md` §4.1, `contracts.md` §14, and `contracts.md` §16 all
require `activation.finalized`; only the product disagrees. The misspelling
exists in no source, test, fixture, corpus, migration, or golden file, and
rebuild is Effect-driven, so the repair changes no other observable. Every other
family in the §16 sequence was verified present verbatim in a product writer.

The repair belongs to task 22, not a reopened task 4: task 22's required
commands are the only gates that execute the accepted activation Event-type
name, and the assertion that executes it is task 22's own (DEC-167, DEC-173).
Task 22 gains `internal/workflow/activation.go` bounded to one exhaustive
operation-to-Event-type mapping, and `cmd/fid/activation_test.go` for a
non-tagged regression assertion that runs under an already-required command. No
Event-schema, public-interface, contract, command, corpus, or VAL change.
Recovery is RUN-RESUME under TRG-21-022 / ESC-21-016 / DEC-175, task-22 attempt
3, `in-scope-implementation` in task scope, no-progress `0/2`.

DEC-174 resolved ESC-21-015 on 2026-09-03: the §16 journey authors its own
workflow before `issue create` and its exact prefix is 22 Events, with no
product, command, public-interface, Event-schema, transition, or frozen-workflow
change. DEC-173 resolved ESC-21-014 on 2026-09-02: task 20's scope was amended
to cover every path its absence controls walk. Detail for both is in the Phase
21 archive.

## Archive Pointers

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
