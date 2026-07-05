# TL Escalations - Current

> **Purpose.** Unresolved escalations that exceed local-decision authority and
> need full Tech Lead or human input. Default-loaded by every agent session.
>
> **Scope.** Open items only. Resolved escalations live in the relevant phase
> archive.
>
> **Size target:** under 80 lines.

## Standing Escalation Triggers

### TL-Lite Triggers

Escalate, do not decide locally, when an issue:

1. Changes a public interface.
2. Changes the data model or persistence contract.
3. Overrides or contradicts a standing decision.
4. Affects more than one task or subsystem.
5. Introduces a dependency, tool, quality gate, or workflow.
6. Materially changes testing strategy or coverage policy.
7. Exposes ambiguity in product intent or acceptance criteria.
8. Resumes context-exhausted or partial work spanning schema/data model,
   dependencies, public interfaces, canonical docs, or multiple Developer
   sessions.

### QA Reviewer Triggers

QA findings escalate here if they have architecture implications, ambiguous
acceptance criteria, unclear security/privacy severity, contradict a standing
decision, span tasks/subsystems beyond the review window, plausibly indicate a
material spec deviation, or expose a quality-gate anomaly not anticipated by the
QA plan.

## Open Escalations

None.

## Recent Resolutions

- Phase 16 TL-Lite first-step blocker triage restored 2026-07-05 by DEC-047:
  task-level blocker reports, failing gates, report-lint failures, missing
  command evidence, and ad hoc gate substitutions must be triaged by TL-Lite
  before a final STOP is appended or resumed. No open escalation remains; the
  next runner must first triage the existing Task 5 STOP.
- Phase 16 Task 5 gate blocker resolved locally 2026-07-05 by DEC-046:
  `internal/activity/activity_test.go` may receive the same narrow deterministic
  fixture remediation pattern as DEC-044 because it violates
  `08_QUALITY_PLAN.md` TDD-006 and trips the PRD-ERR-014 clock-skew guard. No
  open escalation remains; Task 5 still must rerun the full required gate.
- Phase 16 mid-phase QA cadence defect resolved locally 2026-07-05 by DEC-045:
  Phase Runner must run the overdue Review 1 for Tasks 1-2 before any more
  Developer work, then run Review 2 for Tasks 3-4 after Task 4 finalizes. No
  open escalation remains unless Review 1 returns Blocking findings or opens an
  ESC.
- Phase 16 Task 4 gate blocker resolved locally 2026-07-05 by DEC-044:
  `cmd/fid/activity_test.go` may receive a narrow deterministic fixture
  remediation in the Task 4 unblock window because it violates
  `08_QUALITY_PLAN.md` TDD-006 and trips the PRD-ERR-014 clock-skew guard. No
  open escalation remains; Task 4 still must rerun the full required gate.
- Phase 16 Task 3 VAL-007 / `dev gate` scope mismatch resolved 2026-07-05 by
  DEC-043: `dev gate` is part of CF-011 Task 3 activity reconciliation and emits
  `gate_recorded`; no open escalation remains.
- Phase 15 (`fid dev *` composites + cleanup) raised **no escalations**. Two
  TL-Lite triage events resolved locally as RESUME: TRG-001 (Review-1 artifact/
  structure improvements) and TRG-002 (`internal/dev` `Finish`/`Block`
  function-size debt, folded into CF-017). See
  `../phases/phase-15/tl-escalations.md`.
- ESC-001 (Phase 14 Task 2 missing integration-gate evidence) resolved 2026-06-28
  by the Tech Lead as do-now remediation; remediation commit `0cba161` recorded
  the gate PASS and QA cleared the block. Archived at Phase 14 wrap in
  `../phases/phase-14/tl-escalations.md` (TRG-002/TRG-003 RESOLVED).
- Phase 13 (Labels + v0.1d closure) raised no escalations and recorded
  no TL-Lite triage events; the two final-review improvement findings were
  deferred to v0.1e via CF-015 (see `../phases/phase-13/tl-escalations.md`).
- Phase 12 (Requirements + links) raised no escalations; TRG-001 and TRG-002
  both resolved locally as RESUME (see `../phases/phase-12/tl-escalations.md`).

<!--
Add new escalations above this comment using the template at
project-docs/templates/tl-escalations-template.md. Newest at top.
Use global ESC-NNN sequential numbering that never resets.
-->
