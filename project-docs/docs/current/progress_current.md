# Progress - Current

> **Purpose.** In-flight phase status only. Loaded by default in every agent
> session. Historical phase progress lives under `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 9 (Concerns + derived flag + v0.1c closure) is **CLOSED** on
branch `phase-9`. **v0.1c is complete.** QA final review: **BLOCK** (2 blocking /
0 major / 3 improvement); both blocking findings fixed by the Task 4 fixup and
re-verified by QA Review 2 (**CONTINUE**). Human dogfood 5/5 PASS (HUMAN-011a..e).
**VAL-001..VAL-011 all PASS** at wrap. Phases 0-8 archived under
`../phases/phase-N/`.

Phase 9 delivered:

- `fid concern raise` — validate-before-append (unknown role / missing issue /
  missing completion / cross-issue completion) + derived-flag round-trip; read
  back via `fid issue get --section concerns` (Task 1, VAL-001/002/003).
- `fid concern address` — in-place `open`→`addressed`, flag clearing, and the
  `concerns_addressed` rule-evaluation-seam proof with real writes (Task 2,
  VAL-004/005/006/007).
- `fid concern list` — issue/status filters, all four formats, empty=`[]`;
  delete-DB + `fid rebuild` parity; end-to-end v0.1c closure smoke (Task 3,
  VAL-008/009/010).
- Task 4 fixup (ESC-006): `addressed_in_completion_id` exposed on both concern
  read surfaces (DEC-023) + deterministic `rowid DESC` ordering (DEC-024); cleared
  the benchmark gate and resolved VAL-004.

Planning decisions DEC-A..E (human-confirmed pre-execution) all honoured. In-phase
DEC-023/DEC-024 issued to resolve ESC-006. See `../phases/phase-9/decisions.md`.

Coverage held: `cmd/fid` 77.8%, `internal/issues` 81.8%, `internal/journal`
80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5% (COV-001/002 satisfied).

Open carry-forward register: CF-003 (role corrupt-tail, **now due — schedule at
Phase 10 / v0.1d planning**), CF-008 (resolution type/reason, v0.2), CF-011
(consistent `activity_log`, v0.1e), CF-013 (`complete.go` orchestration size,
v0.1e).

## Next Action

- **Phase 10 is the next planned cut** (roadmap §v0.1d — Claims, Permissions,
  Links: `fid issue claim`/`unclaim`/`claim list` with TTL + atomic acquire +
  human claims, claim-transition integration, the in-code permissions matrix +
  `fid policy show`, `fid requirement sync` minimal, `fid link add/remove/list`
  with auto-inverse, and label CRUD). Not yet planned — Tech Lead phase start
  creates the validation contract and task list, and **slices v0.1d into coherent
  phases** (it is large; roadmap has no pre-baked slice for it).
- **At Phase 10 planning, schedule CF-003** (role-validation corrupt-tail
  tolerance) — it is explicitly targeted "before v0.1d role/permission work."
- Carry into the next QA plan the two standing checks Phase 9 surfaced
  (`decisions.md` deep-retro rec): cross-surface read-back parity and read
  ordering determinism (journal append order, not `created_at`).
- Standing commitment before release-readiness: decide where held-out scenarios
  live and how implementation agents are prevented from loading them.
</content>
