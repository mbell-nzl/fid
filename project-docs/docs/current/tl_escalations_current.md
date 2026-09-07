# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None. ESC-22-001 through ESC-22-004 are all resolved and archived in
`../phases/phase-22/tl-escalations.md`.

## Latest Resolution

ESC-22-004 resolved on 2026-09-07 under DEC-195. The contradiction was real and
was a Tech Lead planning omission, not a Developer or TL-Lite defect. Every
accepted record except one agrees that cancellation is its own terminal outcome:
`contracts.md` section 7, `public-interface.md` section 4, VAL-031's
mislabelled-outcome failure condition, the convergence inventory, and QA plan
R-14. Only `schema.md` section 1 was silent, listing five 0013 CHECK member-list
replacements and omitting the one its own contract required. DEC-195 amends the
schema rather than the assertion: `'cancelled'` joins the
`execution_receipts.terminal_status` member list and the composite CHECK gains
one `cancelled` clause admitting signal-only, exit-code-only, or neither, never
both, and never `timed_out = 1`. Execution is an authorized task-8 same-task
continuation with task-1 provenance recorded; task 1 stays COMPLETE. Two bounded
convergence surfaces travel with it: the retained golden DDL string in
`internal/db/operational_kernel_migration_test.go` and the closed outcome filter
in `internal/evidence/receipt_read.go`. Task 8's existing commands already
execute all of them. TRG-22-009's BLOCK, its category, and its refusal to
substitute an admitted status were correct; no no-progress attempt is counted,
because the blocker was a missing authority, not exhausted repair capability.

ESC-22-003 was resolved on 2026-09-06 under DEC-194 (task 4's undelivered
accepted checkpoint/inspection domain scope moved into task 6), ESC-22-002 under
DEC-192/193, and ESC-22-001 under DEC-190/191. All four and their resolutions are
archived in `../phases/phase-22/tl-escalations.md`.

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

Before classifying an owner as ambiguous, read the accepted ownership records:
the candidate task specs' file scopes, `schema.md` section 5 owner column,
`public-interface.md` section 3 task column, and `contracts.md` subject
statements. Behavior that is missing is not the same as behavior that is
unassigned; undelivered accepted in-scope work is bounded repair (DEC-194).

A genuine specification contradiction is different from both, and it does escalate
(DEC-195): when an accepted contract clause names a closed-domain value that the
canonical schema rejects, and no accepted record assigns anyone the amendment,
there is nothing to deliver and nothing to repair. Block, preserve the evidence,
and never substitute an admitted value to make a gate pass. A blocker of that kind
is a missing authority, not exhausted repair capability, so it consumes no
no-progress budget.
