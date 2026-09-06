# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None. ESC-22-001, ESC-22-002, and ESC-22-003 are all resolved and archived in
`../phases/phase-22/tl-escalations.md`.

## Latest Resolution

ESC-22-003 resolved on 2026-09-06 under DEC-194. The owner was never ambiguous:
`tasks/task-4.md`, `schema.md` section 5, `public-interface.md` section 3, and
`contracts.md` sections 3/3.1 all assign the `run.checkpoint` producer and the
`run.inspect` read API to task 4. Task 4 closed COMPLETE without writing
`internal/workflow/run_checkpoint.go` (which has never existed in any commit) or
any inspection read API, so this is undelivered accepted in-scope work, not a
protected task boundary. Because task 4 is closed with two successors on top and
the Runner's only recovery route is a same-task continuation, DEC-194 moves
execution of exactly four `internal/workflow` files into task 6 with task-4
provenance recorded. Task 4 stays COMPLETE and unedited. No contract, schema,
Event, VAL row, validator, gate, task allocation, or workflow control changed;
no Human decision was required. Task 6's existing required commands already
execute the assigned files.

ESC-22-002 was resolved earlier the same day under DEC-192/193, and ESC-22-001
under DEC-190/191. All three and their resolutions are archived in
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

Before classifying an owner as ambiguous, read the accepted ownership records:
the candidate task specs' file scopes, `schema.md` section 5 owner column,
`public-interface.md` section 3 task column, and `contracts.md` subject
statements. Behavior that is missing is not the same as behavior that is
unassigned; undelivered accepted in-scope work is bounded repair (DEC-194).

## ESC-22-004 - Task 8 cancelled Receipt requirement contradicts canonical terminal-status DDL

- **Phase:** 22
- **Status:** OPEN
- **Raised by:** TRG-22-009
- **Category:** specification-contradiction
- **Summary:** Task 8's accepted VAL-031/contracts behavior requires a fid-observed cancellation to persist a durable Receipt with `terminal_status` `cancelled`, but canonical `execution_receipts.terminal_status` DDL excludes that value.
- **Evidence:** `project-docs/docs/phases/phase-22/contracts.md` section 7 states that cancellation adds terminal outcome `cancelled`; `migrations/0013_p21_operational_kernel.up.sql:250-268` permits only `succeeded`, `failed`, `timed-out`, `signaled`, and `spawn-failed`. The committed task-8 report preserves focused passing cancellation/workflow evidence and records the required gates `BLOCKED` because VAL-031 cannot be completed without changing the protected DDL or accepted requirement. Runner-supplied verification confirms implementation `b965b24748f57a0f5e42b8dd3054dfe6fd5dbb74` and report `764396118ce7d1cbb9a068ca7179ad204afb1566`; independent report-lint returned `OUTCOME BLOCKED` and `PASS report file committed and clean` (terminated `exit_code: 0`).
- **Why blocked:** The current task forbids DDL, wire-shape, and contract changes. Mapping cancellation to an existing permitted status violates the accepted cancellation Receipt behavior, while modifying the DDL, schema/task scope, VAL-031, or contracts requirement is protected authority outside TL-Lite and task 8.
- **Required resolution:** Full Tech Lead/Human authority must reconcile the accepted durable `cancelled` Receipt requirement with the canonical `execution_receipts.terminal_status` CHECK. It must record whether the canonical DDL and bounded owning scope are amended to admit `cancelled`, or the accepted VAL-031/contracts behavior is amended, and then provide the authorized same-task continuation and verification path. Preserve all failed receipts and do not substitute an existing terminal status.
- **Links:** TRG-22-009; `project-docs/docs/phases/phase-22/tasks/task-8-done.md`; `project-docs/docs/phases/phase-22/tasks/task-8.md`; `project-docs/docs/phases/phase-22/contracts.md` section 7; `migrations/0013_p21_operational_kernel.up.sql:250-268`; `project-docs/specs/09_AI_WORKFLOW.md` §§2.3, 9.
- **Human attention required:** yes
