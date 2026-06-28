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

- ESC-001 (Phase 14 Task 2 missing integration-gate evidence) resolved 2026-06-28
  by the Tech Lead as do-now remediation; see the ESC-001 Resolution block below.
  Clears on the Task 2 remediation commit + QA confirmation; archive at phase wrap.
- Phase 13 (Labels + v0.1d closure) raised no escalations and recorded
  no TL-Lite triage events; the two final-review improvement findings were
  deferred to v0.1e via CF-015 (see `../phases/phase-13/tl-escalations.md`).
- Phase 12 (Requirements + links) raised no escalations; TRG-001 and TRG-002
  both resolved locally as RESUME (see `../phases/phase-12/tl-escalations.md`).
- ESC-008 (Phase 11 claim-write projection refresh lock ordering) resolved by
  DEC-028 and Task 5/VAL-015, and archived in
  `../phases/phase-11/tl-escalations.md` at Phase 11 wrap.

<!--
Add new escalations above this comment using the template at
project-docs/templates/tl-escalations-template.md. Newest at top.
Use global ESC-NNN sequential numbering that never resets.
-->

## ESC-001 - Task 2 closeout missing required integration evidence

Date: 2026-06-28
Source: Triage trigger from phase-14 QA review and Task 2 completion report.
Status: Resolved 2026-06-28 (Tech Lead) — do-now remediation directed; clears on
the Task 2 remediation commit + QA confirmation. Archive at phase wrap.

The Phase 14 Task 2 completion report omitted the required integration gate from
closeout evidence. The QA reviewer spot-checked `go test -tags integration -race
./...` and it passed, but the task closeout artifact itself still lacks the
required evidence entry. This is a QA-blocking artifact defect because required
commands must be recorded in the completion report as PASS or BLOCKED.

Trigger category: QA blocking / required completion report evidence missing.

### Resolution (Tech Lead, 2026-06-28)

Disposition: **do-now remediation** — not a misclassification, deferral, cancel,
or VAL change. The QA reviewer's BLOCK is correct: integration-tagged tests
exist under `cmd/fid` and Task 2 modified `cmd/fid` (`run.go`, `run_test.go`,
`root.go`), so `go test -tags integration -race ./...` is in scope per QA plan
section 8 (line 156) and section 9 (line 183). Recording it as "not run (no
integration-tagged tests added in Task 2)" is a narrower/incorrect disposition
of a required conditional command (QA plan sections 6 and 11) — an artifact
defect, not a behavior defect (reviewer spot-check passed in 53.385s).

Root cause: the conditional-command note "if any test is integration-tagged" was
read as "if this task adds integration-tagged tests" rather than "if the package
set under test contains integration-tagged tests."

Required action (Developer, via Phase Runner; the Tech Lead does not run gates
or author completion-report evidence):

1. Rerun exactly `go test -tags integration -race ./...`.
2. Update `task-2-done.md` to record it as `PASS` (with the result) or `BLOCKED`
   with a reason, replacing the "Conditional command not run" line.
3. Commit the corrected completion-report evidence.
4. QA reviewer confirms the report now records the gate; the BLOCK clears with no
   re-run of Task 2 product code.

No spec or behavior change; no VAL weakened (VAL-002..006 stand as reviewed).
Recurrence prevention: Task 3's identical conditional wording was clarified in
its spec so the same misread cannot reproduce the block.
