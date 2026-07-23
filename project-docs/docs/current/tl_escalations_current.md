# TL Escalations - Current

> **Purpose.** Open items that exceed local-decision authority and need full
> Tech Lead or human input.
>
> **Scope.** Open items only. Resolved phase items move to the phase archive.
>
> **Size target:** under 80 lines.

## Open Escalations

- **ESC-001 (Phase 18)** — resolved by **DEC-072** (2026-07-23); stays listed
  until the remediation lands and QA re-review passes, then archives to
  `../phases/phase-18/tl-escalations.md`. See the full entry below.

## Recent Archive Notes

- Phase 17 wrapped 2026-07-18. It raised five escalations, all resolved and
  archived to `../phases/phase-17/tl-escalations.md`: ESC-004 (Task 1
  report-lint), ESC-005 (Task 2 report-lint, backfilled), ESC-009 (QA model
  qualification), ESC-010 (rebuild parity; resolved by DEC-056 Branch A), and
  ESC-011 (Task 7 report-lint, backfilled with a new ID). That file also records
  the WRAP-F1/F2 register reconciliation. The two open retrospective items -
  recurring report-lint commit-metadata friction and **ESC-NNN IDs not being
  globally unique** (ESC-001/004/005 collide across phases) - are now tracked as
  **CF-023** (DEC-058) and flagged for the deep retrospective before v0.2a.
- Phase 16 raised no `ESC-NNN` items. Local unblock and workflow decisions
  DEC-043..048 plus TRG-046 are archived in
  `../phases/phase-16/tl-escalations.md` and `../phases/phase-16/decisions.md`.
- Phase 16 VAL-013 is not an escalation; the completed human checklist records
  PASS for HUMAN-16a..d.

## Standing Escalation Triggers

TL-Lite escalates instead of deciding locally when an issue:

1. Changes a public interface, CLI contract, output shape, journal event schema,
   data model, persistence contract, dependency, tool, quality gate, or workflow.
2. Contradicts `decisions_active.md`, `decisions_principles.md`, `AGENTS.md`, or
   a recorded phase decision.
3. Affects more than one task/subsystem, exposes product-intent ambiguity, or
   changes testing strategy / coverage policy.
4. Resumes context-exhausted or partial work spanning schema/data model,
   dependencies, public interfaces, canonical docs, or multiple Developer
   sessions.

QA findings escalate here if they have architecture implications, ambiguous
acceptance criteria, unclear security/privacy severity, contradict a standing
decision, span tasks/subsystems beyond the review window, plausibly indicate a
material spec deviation, or expose a quality-gate anomaly not anticipated by the
QA plan.
## ESC-001 — Phase 18 QA-blocking maintenance test deadlock and stale specification synchronization

- **Status:** **Resolved by DEC-072 (2026-07-23)** — remediation and QA
  re-review pending; archives at wrap.
- **Raised by:** TL-Lite triage
- **Source:** Phase 18 QA final review; TRG-006 (filed as TRG-001; relabelled
  per DEC-072 (10) to clear the in-phase collision with the DEC-067 TRG-001)
- **Escalation triggers:** quality-gate failure; public/specification synchronization across subsystems
- **Issue:** `workflow.Load` deadlocks in `TestCheckWorkflowStatesReadsProjectedWorkflowTables` when the SQLite fixture uses one connection. The required reproduction command does not terminate on HEAD: `go test -race -count=1 -v -timeout 5s ./internal/maintenance -run '^TestCheckWorkflowStatesReadsProjectedWorkflowTables$'`. QA-30 also requires specification synchronization, but `02_TRD.md` §§4.3/4.6, `03_TRD_SCHEMA.md` §14.11, and current phase status documents remain stale.
- **Why escalated:** The unresolved deterministic quality-gate failure cannot be RESUMEd. Resolving it also requires specification and current-context updates spanning more than one subsystem and canonical documents.
- **Evidence:** `project-docs/docs/phases/phase-18/qa/final-review.md`; verified review commit `46d8013abc46d44ed1694f8290c163d612de5f3a` (`docs(qa): complete phase 18 final review`). Earlier Task 5 commits `37b3fb524bc5ba3d2774f6c14a12e78cbc028e1f` and `d7f11995c77d8ed028b4ecaaa2cbd61e8688522c` were verified earlier, but the final review records that the required maintenance command now hangs.
- **Next:** Full Tech Lead to resolve the root cause, coordinate the necessary specification synchronization, and obtain fresh deterministic gate evidence before resumption.
- **Resolution (DEC-072, 2026-07-23):** Root cause verified — `workflow.Load`
  (`internal/workflow/registry.go:23`) issues nested per-workflow queries while
  the outer `workflows` result set still holds the connection, so it blocks
  forever on a pool capped at one connection. Production `db.Open` sets no cap,
  which is why the real entrypoint works and the fixture hangs; the loader is
  still wrong and the fixture stays. Disposition: **do-now remediation within
  Task 5** — materialize the workflow IDs and release the outer rows before any
  nested query, add a bounded-timeout single-connection regression test, rerun
  the full seven-command gate to completion, and correct the over-claimed command
  records in `task-5-done.md`. No assertion, contract, fixture, or gate
  weakening. The QA-30 specification synchronization half is **applied**
  (`02_TRD.md` §§4.4/4.6, `03_TRD_SCHEMA.md` §§3/14.11/15, canonical-suite
  verdict table, current docs). ID-collision note: this ESC-001 collides with
  Phase 17's ESC-001 — the known CF-023 defect — and is left as-is; renumbering
  published escalation IDs is a retrospective decision, not a mid-phase edit.
