# Phase 7 Escalations (archived)

> Resolved escalations moved here at Phase 7 wrap. Open items live in
> `../../current/tl_escalations_current.md`.

## ESC-005 — AC verification latest ordering nondeterministic within one second (RESOLVED)

- **Raised:** 2026-06-14 by QA Reviewer / Phase 7 Review 1 (via TRG-001).
- **Trigger:** QA Reviewer triggers 5 and 6 — spans later `ac list` / workflow
  consumers and is a plausible material spec deviation.
- **Context:** Task 2 stored `ac_verifications.verified_at` at `time.RFC3339`
  second precision while `latestVerificationsByRole` selected latest by
  `verified_at DESC, id DESC`. Two same-role verifications inside the same second
  tie on `verified_at` and break by **random** `ACV-` content-hash id, so a fast
  fail→pass could report the wrong latest status; VAL-004 was only passing because
  tests slept past a one-second boundary.
- **Resolution:** **DEC-021** (human-confirmed 2026-06-14). Latest-per-role is
  the verification appended last in journal order (ADR-0001), resolved by SQLite
  implicit `rowid` DESC — deterministic and rebuild-stable, independent of
  `verified_at` second resolution and the random id. `verified_at` stays
  display-only. No schema migration; no CLI shape change. Implemented + de-sleep
  tests landed as **Task 2b**; VAL-004 strengthened to forbid timestamp-boundary
  proof and re-review (Review 1 fixup) cleared it. Spec sync:
  `03_TRD_SCHEMA.md` `ac_verifications` comment now states journal-append-order
  resolution.
- **References:** `decisions.md` DEC-021; `triage-log.md` TRG-001;
  `qa/review-1.md`, `qa/review-1-fixup.md`; `tasks/task-2b.md`; VAL-004/006/008;
  `internal/issues/detail.go` (`rowid DESC`).
