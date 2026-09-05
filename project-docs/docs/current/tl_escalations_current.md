# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

### ESC-22-001 - Protected workflow-sequence control cannot represent DEC-189 supersession

- **Phase:** 22
- **Status:** OPEN
- **Raised by:** TRG-22-006
- **Category:** protected workflow-control definition/policy
- **Summary:** `project-docs/tools/workflow-sequence-lint.sh` rejects the user-authorized TRIAGE TRG-22-005 and bound task-2 TASK-RECOVERY despite the accepted contract permitting a later same-task schema-v2 TRIAGE REPAIR to replace its pending dispatch and retain the review-form trigger identity.
- **Evidence:** From `/home/mark/projects/fid`, the supplied exact append reproducer terminated `exit_code: 1` with `INVALID TRIAGE REPAIR must be followed by its bound TASK-RECOVERY before TRIAGE; INVALID schema v2 TRIAGE Originating task does not match Trigger signature`; the ledger remained unchanged and standalone validation remained `VALID state=REPAIR_PENDING`.
- **Why blocked:** The required corrective edit is to `project-docs/tools/workflow-sequence-lint.sh`. `09_AI_WORKFLOW.md` §§2.3 and 9 prohibit normal control-conformance self-repair of that verifier; changing it requires separate protected workflow-control authority.
- **Required resolution:** Full Tech Lead or authorized protected workflow-control decision must reconcile the verifier implementation with `project-docs/templates/run-ledger-template.md` and DEC-189's same-task supersession and review-form trigger contract, then provide a concrete authorized recovery route. Preserve the valid Phase 22 ledger; do not bypass the validator or authorize product/task recovery first.
- **References:** TRG-22-006; `project-docs/specs/09_AI_WORKFLOW.md:427-458,586-614,703-713`; `project-docs/templates/run-ledger-template.md`; `project-docs/docs/phases/phase-22/decisions.md:222-270` (DEC-189); `project-docs/docs/phases/phase-22/run-ledger.md:158-209`.

## Latest Resolution

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
