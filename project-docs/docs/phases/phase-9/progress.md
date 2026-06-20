# Phase 9 Progress (archived)

> Phase: 9 - Concerns + derived flag + v0.1c closure (roadmap §v0.1c)
> Status: CLOSED on branch `phase-9`, 2026-06-21. **v0.1c complete.**

## Outcome

Phase 9 shipped the **concern write/read surface** (`fid concern raise`,
`fid concern address`, `fid concern list`) on the existing `concerns` schema,
`issues_enriched` derived columns, and concern read path, proved the
`concerns_addressed` rule with real writes, and ran the end-to-end v0.1c closure
smoke — completing the v0.1c slice (ACs, completions, concerns).

QA final review was **BLOCK** (2 blocking / 0 major / 3 improvement); both
blocking findings were fixed by the Task 4 fixup and re-verified by QA Review 2
(**CONTINUE**, 0 blocking / 1 improvement). Human dogfood checklist 5/5 PASS
(HUMAN-011a..e). **VAL-001..VAL-011 all PASS** at wrap.

## Delivered

- **Task 1 — `fid concern raise`** (`f35b3fe`). Write command +
  validate-before-append (unknown role / missing issue / missing completion /
  cross-issue completion) + derived-flag round-trip; read back via
  `fid issue get --section concerns`. (VAL-001/002/003)
- **Task 2 — `fid concern address` + rule proof** (`3ea2d4f`). In-place
  `open`→`addressed` mutation, flag clearing, and `concerns_addressed`
  rule-evaluation-seam proof with real concern data. (VAL-004/005/006/007)
- **Task 3 — `fid concern list`, rebuild parity, v0.1c smoke** (`a172634`).
  Top-level list (issue/status filters, all four formats, empty=`[]`), delete-DB +
  `fid rebuild` parity, end-to-end v0.1c closure smoke. (VAL-008/009/010)
- **Task 4 — fixup (ESC-006)** (`c927de0`). Exposed `addressed_in_completion_id`
  on both concern read surfaces (DEC-023) and made concern read ordering
  deterministic via `ORDER BY rowid DESC` (DEC-024); cleared the benchmark gate
  and resolved VAL-004. (VAL-004 resolved; VAL-008/009 preserved)

## Decisions

- DEC-A..E (planning decisions, human-confirmed 2026-06-20) all honoured.
- DEC-023 (concern read models expose `addressed_in_completion_id`) and DEC-024
  (concern read ordering by journal append order, `rowid DESC`) issued mid-wrap to
  resolve ESC-006. `06_DESIGN_PLAN_CLI.md` §6 carries both. See `decisions.md`.

## Escalations

- **ESC-006 — RESOLVED.** Final-review benchmark-gate failure + VAL-004 read-back
  gap; direction DEC-023/DEC-024 + Task 4; QA Review 2 verified the fix. Archived
  in `tl-escalations.md`.

## Carry-forward

- No CF items closed or opened by Phase 9. CF-003 / CF-008 / CF-011 / CF-013
  re-deferred (CF-003 now due for scheduling at Phase 10 / v0.1d planning).

## Coverage at close

`cmd/fid` 77.8%, `internal/issues` 81.8%, `internal/journal` 80.3%,
`internal/policy` 85.0%, `internal/workflow` 95.5% (COV-001/COV-002 satisfied).
Total statement coverage 76.5%.
