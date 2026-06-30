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
