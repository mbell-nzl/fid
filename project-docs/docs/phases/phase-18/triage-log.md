# Phase 18 Triage Log

> **ID correction (Tech Lead, 2026-07-23, DEC-072).** This entry was originally
> filed as `TRG-001`, colliding with the Phase 18 TRG-001 already recorded in
> `decisions.md` (blocked QA plan, resolved by DEC-067). It is relabelled
> **TRG-006** — the next free ID in this phase's sequence — so the resume is
> unambiguous. Content is otherwise unchanged. The collision is another instance
> of the CF-023 ID-uniqueness defect flagged for the deep retrospective. The
> append-only `run-ledger.md` is not modified.

## TRG-006 (originally filed as TRG-001)

- **Date:** 2026-07-20
- **Trigger:** QA Lead final review returned BLOCK with two findings: (1) `workflow.Load` deadlocks in `TestCheckWorkflowStatesReadsProjectedWorkflowTables` when the SQLite fixture uses one connection, so required quality commands do not terminate on HEAD despite Task 5 reporting all-green; reproduction: `go test -race -count=1 -v -timeout 5s ./internal/maintenance -run '^TestCheckWorkflowStatesReadsProjectedWorkflowTables$'`. (2) QA-30 requires Tech Lead specification synchronization, but `02_TRD.md` §§4.3/4.6, `03_TRD_SCHEMA.md` §14.11, and current phase status documents remain stale.
- **Source:** `project-docs/docs/phases/phase-18/qa/final-review.md`
- **Disposition:** BLOCK
- **Reason:** QA Blocking finding (1) is an unresolved required quality-gate failure. The supplied verification confirms the final-review and earlier Task 5 commits, but does not verify a fix for the current hang. Finding (2) requires canonical specification/current-context synchronization across documents and subsystems.
- **Escalation:** ESC-001
- **References:** `project-docs/docs/phases/phase-18/qa/final-review.md`; `project-docs/specs/02_TRD.md` §§4.3, 4.6; `project-docs/specs/03_TRD_SCHEMA.md` §14.11; `project-docs/specs/09_AI_WORKFLOW.md` §§2.3, 9.
- **Resolution (2026-07-23):** **RESOLVED by DEC-072.** Root cause verified as a
  connection-lifetime defect in `workflow.Load` (`internal/workflow/registry.go`)
  — nested per-workflow queries issued while the outer `workflows` result set
  still holds the connection. Disposition: do-now remediation within Task 5 (fix
  the loader, add a bounded-timeout single-connection regression test, do not
  touch the fixtures or the gate), correct the Task 5 completion report's
  over-claimed command records, and apply the QA-30 spec sync under Tech Lead
  authority. ESC-001 is closed by DEC-072.
