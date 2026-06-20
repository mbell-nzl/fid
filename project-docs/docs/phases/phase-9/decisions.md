# Phase 9 Decisions

> In-phase decisions (DEC-NNN / DEC-A..E). Per-doc spec-sync verdicts and the
> full decision history are finalized at phase wrap. Newest at top.

## DEC-024: Concern read ordering is deterministic by journal append order (`rowid DESC`)

Source: ESC-006 / final review benchmark-gate failure
(`TestConcernListFiltersByIssueAndStatus/all_concerns`), Phase 9 final review.
Status: Decided 2026-06-20 (Tech Lead); resolves an ESC-006 closeout blocker.

Decision: `fid concern list` (`internal/issues/concern_list.go`) and the
`fid issue get --section concerns` read model (`internal/issues/detail.go`
`issueConcerns`) order concern rows by **journal append order, newest first**
(`ORDER BY rowid DESC`) instead of `ORDER BY created_at DESC, id DESC`.

Rationale:
- `created_at` has second resolution, so two concerns raised in the same second
  tie; the `id DESC` tiebreak is a content hash, which orders rows
  nondeterministically relative to creation order. The phase benchmark gate
  surfaced this as a flaky `cmd/fid` test (`first concern issue_id` mismatch).
- `rowid DESC` is monotonic, append-ordered, and rebuild-stable: `fid rebuild`
  replays the journal in order, so post-rebuild `rowid` order equals original
  append order. VAL-009 rebuild parity is preserved.
- This mirrors the human-confirmed DEC-021 precedent (AC verification
  latest-per-role resolved by journal append order, no timestamp dependence).

Scope: ordering only. No schema change, no column change, no new public flag.
The newest-first contract is unchanged in intent; only the ordering key becomes
deterministic. `02_TRD.md` / `03_TRD_SCHEMA.md`: no schema change (the `concerns`
table is a normal rowid table). `06_DESIGN_PLAN_CLI.md` §6 gains an explicit
ordering note. Spec-sync verdicts finalized at wrap.

## DEC-023: Concern public read models expose `addressed_in_completion_id`

Source: ESC-006 / VAL-004 / QA-12 public read-back failure, Phase 9 final review.
Status: Decided 2026-06-20 (Tech Lead); resolves an ESC-006 closeout blocker.

Decision: both public concern read surfaces expose the addressing completion:
- `ConcernListItem` (`internal/issues/concern_list.go`) gains
  `addressed_in_completion_id` (nullable `*string`, JSON
  `"addressed_in_completion_id"`), selected and scanned from the existing
  `concerns.addressed_in_completion_id` column.
- `IssueConcern` (`internal/issues/detail.go`) gains the same nullable field, so
  `fid issue get --section concerns` and `fid concern list` stay parity-aligned
  (design §6: "row shape mirrors `fid issue get --section concerns` plus
  `issue_id`").

Rationale:
- The write path (`concern_address.go`) already persists the column on
  `issue.concern_addressed`; the lifecycle fact existed in the journal/DB but was
  invisible to public read-back, failing VAL-004 / PRD-DEV-006.
- The value is `null` for `open` concerns and the addressing `TC-` ID once
  addressed — the conditional/binary read-back both states VAL-004 requires.

Scope: read-model field addition only. No schema change (column exists since
`migrations/0002_v01_tables.up.sql`), no write-path change, no workflow-edge
change. `06_DESIGN_PLAN_CLI.md` §6 concern list JSON example updated to include
the field. Spec-sync verdicts finalized at wrap.

## Phase 9 planning decisions (DEC-A..E, human-confirmed 2026-06-20)

Recorded in the validation contract §Phase scoping notes. Summarized here for the
wrap audit; each task spec and QA artifact cites them by DEC-A..E.

- **DEC-A — PRD-WF-009 proven at the rule-evaluation seam, no default edge.**
- **DEC-B — concern domain lives in `internal/issues`; no new boundary.**
- **DEC-C — validate-before-append linkage rules (CF-007).**
- **DEC-D — address is a status mutation; raise→address only.**
- **DEC-E — no `activity_log` for concern writes (DEC-020 / CF-011 precedent).**

(Outcome audit, spec-sync verdicts, TRG audit, and workflow notes are appended at
phase wrap once the Task 4 fixup is QA-verified and VAL-004/benchmark gate clear.)
