# Phase 18 TL Escalations (archived)

> Archived at wrap (2026-07-25, DEC-073). Phase 18 raised one `ESC-NNN`; it is
> resolved and closed. Open items, if any, live in
> `../../current/tl_escalations_current.md`.

## Summary

| ID | Raised | Source | Status | Resolution |
|---|---|---|---|---|
| ESC-001 | 2026-07-20 | QA Lead final review BLOCK; TRG-006 | **Resolved / closed** | DEC-072; verified by the 2026-07-23 final-review rerun (FRR-E01..E04, FRR-E11) |

Phase 18's other five triages (TRG-001..TRG-005) produced `BLOCK` verdicts that
were adjudicated by the Tech Lead as DEC-067..DEC-071 **without** minting an
`ESC-NNN`; that routing was correct in substance (see the Triage-Log Audit in
`decisions.md`), but in five cases TL-Lite reported an escalation it never
created. That defect class is CF-023, promoted to a required pre-Phase-19
workflow-tooling pass by DEC-074.

## ESC-001 — Phase 18 QA-blocking maintenance test deadlock and stale specification synchronization

- **Status:** **Closed (2026-07-25).** Resolved by DEC-072; remediation landed at
  `795a7df`, report finalized at `37fb202`, QA Lead final-review rerun **PASS**
  at `77a4560`, human checklist PASS at `0edcfc8`.
- **Raised by:** TL-Lite triage
- **Source:** Phase 18 QA final review; TRG-006 (filed as TRG-001; relabelled per
  DEC-072 (10) to clear the in-phase collision with the DEC-067 TRG-001)
- **Escalation triggers:** quality-gate failure; public/specification
  synchronization across subsystems
- **Issue:** `workflow.Load` deadlocked in
  `TestCheckWorkflowStatesReadsProjectedWorkflowTables` when the SQLite fixture
  used one connection. The required reproduction command did not terminate on
  the then-HEAD: `go test -race -count=1 -v -timeout 5s ./internal/maintenance
  -run '^TestCheckWorkflowStatesReadsProjectedWorkflowTables$'`. QA-30 also
  required specification synchronization, but `02_TRD.md` §§4.3/4.6,
  `03_TRD_SCHEMA.md` §14.11, and the current phase status documents were stale.
- **Why escalated:** An unresolved deterministic quality-gate failure cannot be
  RESUMEd. Resolving it also required specification and current-context updates
  spanning more than one subsystem and canonical documents.
- **Evidence at raise:** `qa/final-review.md` (2026-07-20); verified review
  commit `46d8013`; earlier Task 5 commits `37b3fb5` and `d7f1199`.
- **Resolution (DEC-072, 2026-07-23):** Root cause verified — `workflow.Load`
  (`internal/workflow/registry.go`) issued nested per-workflow queries while the
  outer `workflows` result set still held the connection, so it blocked forever
  on a pool capped at one connection. Production `db.Open` sets no cap, which is
  why the real entrypoint worked and the fixture hung; the loader was still
  wrong and the fixture stayed. Disposition: **fix the loader, not the fixture** —
  materialize the workflow IDs and release the outer rows before any nested
  query, add a bounded-timeout single-connection regression test, rerun the full
  seven-command gate to completion, and correct the over-claimed command records
  in `task-5-done.md`. No assertion, contract, fixture, or gate weakening. The
  QA-30 specification synchronization half was applied under Tech Lead authority
  (`02_TRD.md` §§4.4/4.6, `03_TRD_SCHEMA.md` §§3/14.11/15, the canonical-suite
  verdict table, and the current docs).
- **Closure evidence (2026-07-25, DEC-073):** the fresh final review records the
  formerly hung command passing in `1.235s` and the new bounded regression
  `TestLoadReleasesWorkflowRowsBeforeNestedQueries` passing under
  `SetMaxOpenConns(1)` in `0.38s` (FRR-E04); all seven gate commands terminate
  and pass in order (FRR-E01); QA-30 verified committed at `4c029fb` (FRR-E11).
  Registry contents and ordering are unchanged (VAL-014 PASS).
- **Known ID collision:** this `ESC-001` collides with Phase 17's `ESC-001`. Left
  as recorded — renumbering published escalation IDs is a retrospective decision,
  not a mid-phase or wrap edit. Tracked under CF-023 / DEC-074 (2b).
