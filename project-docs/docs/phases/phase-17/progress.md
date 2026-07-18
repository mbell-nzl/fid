# Phase 17 Progress

> Phase: 17 - Maintenance + stub guard + v0.1 closure (roadmap v0.1e,
> slice 4 of 4)
> Status: COMPLETE / WRAPPED (2026-07-18). All ten tasks closed, QA final review
> PASS, Human checklist PASS, VAL-001..025 resolved. v0.1 is closed for early
> dogfood (DEC-057). No Developer task is active.

## Plan

Phase 17 delivers `fid maintenance check`, locked `fid maintenance vacuum`,
explicit v0.2 governance command guards, full read-format conformance, and the
v0.1 closure journey/rebuild/traceability evidence.

## Task Index

| Task | Title | Status |
|---|---|---|
| 0 | Restore `internal/journal` COV-001 coverage | Pending; entry prerequisite |
| 1 | CF-019 markdown formatter split | Pending |
| 2 | Maintenance grace config contract | Pending |
| 3 | Maintenance structural errors | Pending |
| 4 | Maintenance warnings and info | Pending |
| 5 | `fid maintenance check` CLI + formats | Pending |
| 6 | `fid maintenance vacuum` | Pending |
| 7 | Governance stub-table guards | Pending |
| 8 | v0.1 read-format conformance | Pending |
| 9 | Full v0.1 closure journey | Complete (DEC-056 Branch A) |

## Validation

`qa/validation-contract.md` VAL-001..025 all resolved: VAL-001..024 PASS in the
QA Lead final review, VAL-025 PASS via the Human checklist (HUMAN-25a..e PASS).
VAL-022 was narrowed under DEC-056 (Branch A confirmed) and VAL-023 unblocked
with it. DEC-054's Task 0 restored `internal/journal` to 81.0% with no waiver.

## Progress History

- Planning draft (2026-07-10): Tech Lead drafted the Heavy Mode contract and
  nine-task pack; implementation blocked on Human approval.
- Kickoff approval (2026-07-11): Human approved Heavy Mode and DC-A/DC-B/DC-C.
  DEC-049..052 record the mode, maintenance warning/config contract,
  result/exit semantics, and exact governance guard matrix. Canonical affected
  specs were synchronized. No Developer task started.
- Entry re-plan (2026-07-11): QA-BLK-01 accepted. DEC-054 adds tests-only Task
  0 plus an entry review; no coverage waiver. QA-BLK-02 is tracked as ESC-009.
- Entry unblock (2026-07-14): model audit qualified the DEC-053 QA routing
  (closes ESC-009/QA-BLK-02); TRG-001 resumed for Task 0. Tasks 0-8 completed
  with reviews after Tasks 0, 2, 5, 7. Recurring report-lint metadata friction
  drove TRG-002/003/004/005/007 and ESC-004/005/011.
- VAL-022 escalation (2026-07-15): Task 9 exposed a rebuild-parity anomaly;
  TRG-008 BLOCKed and raised ESC-010.
- ESC-010 resolution (2026-07-17): DEC-056 defined the rebuild-parity contract
  (journal is sole source of derived state; clock-relative reads excluded from
  equality), amended VAL-022, and required a diagnostic. Task 9 resumed
  audit-only.
- Closure (2026-07-18): Task 9 completed on Branch A (control reproduced the
  17->20 movement; no product change). QA final review PASS; Human checklist
  PASS. Tech Lead wrapped: DEC-057 closes v0.1, DEC-058/CF-023 disposition the
  report-lint friction, full spec-sync verdicts recorded.

## Carry-Forward (final)

- CF-019 CLOSED — satisfied by VAL-001 (markdown formatter split into four files
  <=400 lines with byte-identical behavior).
- CF-020 NEW (DEC-056) — future-`occurred_at` under burst writes vs the
  future-excluding `recent_activity_7d` window; re-defer to v0.2.
- CF-021 NEW — Go doc comments on `maintenance.VacuumResult`/`Vacuum`; next
  maintenance reopen.
- CF-022 NEW — broader exact-output markdown characterization tests; next
  `internal/format` markdown reopen.
- CF-023 NEW (DEC-058) — workflow-tooling hardening (report-lint /
  completion-report metadata friction; ESC-NNN global uniqueness); flagged for
  deep retrospective.
- CF-008 remains v0.2; CF-017 remains v0.2c; CF-018 remains at the next
  `internal/sqlquery` modification (Task 8 exercised but did not modify it).

## Escalations

All Phase 17 escalations resolved and archived in `tl-escalations.md`:
ESC-004 (Task 1 report-lint), ESC-005 (Task 2 report-lint, backfilled),
ESC-009 (QA model qualification), ESC-010 (rebuild parity; DEC-056 Branch A),
ESC-011 (Task 7 report-lint, backfilled). Register-vs-ledger reconciliation is
recorded as WRAP-F1/F2.
