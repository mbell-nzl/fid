# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None. ESC-21-015 was resolved by DEC-174 on 2026-09-03.

## Latest Resolution

DEC-174 resolved ESC-21-015 on 2026-09-03. TRG-21-021 correctly stopped task
22, but its category was wider than the fact. The product is correct: the
transition engine resolves edges and the terminal-success `completed` category
from the projected workflow tables, and its exit-4 rejection of an undeclared
edge is the guard working. The defective artifact is `contracts.md` §16, which
fixed an exact 20-Event prefix with a direct terminal edge that no workflow
available to kind `task` provides, and never named the journey's workflow.
`task-default` needs four edges to reach a `completed` state,
`subtask-default` needs a second `issue.created` for the required task parent,
and `epic-default` adds a rule-guarded intermediate state.

The journey now authors its own workflow with the already-accepted
`workflow define --applies-to task` and `workflow assign --kind task` leaves
before `issue create`: one `unstarted`, one `started`, one `completed`, one
PRD-WF-013 `cancelled` state, and exactly the edges `unstarted -> started` and
`started -> completed`. §16's prefix is exactly 22 Events, inserting
`workflow.defined` and `workflow.assigned` after `project.initialized`, and
VAL-086's counts are amended to match. VAL-084/085/086 keep their identity,
subject, and strength; the observed rejection becomes required VAL-085 negative
coverage. No product, command, public-interface, Event-schema, transition, or
frozen-workflow change. Recovery is RUN-RESUME under TRG-21-021 / ESC-21-015 /
DEC-174, task-22 attempt 2, `in-scope-implementation` in task scope.

DEC-173 resolved ESC-21-014 on 2026-09-02: task 20's scope was amended to cover
every path its absence controls walk, and it gained the tagged integration gate
under DEC-167. Detail is in the Phase 21 archive.

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
