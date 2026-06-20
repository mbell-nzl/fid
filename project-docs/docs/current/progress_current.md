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

- **v0.1d sliced into 4 phases** (roadmap §v0.1d phase-slicing note, 2026-06-21):
  **10 permissions → 11 claims → 12 requirements+links → 13 labels+closure.**
- **Phase 10 is PLANNED, awaiting human confirmation of DEC-A..E.** Permissions
  matrix + `policy.Authorize` seam + `fid policy show` + retrofit enforcement on
  every existing write + CF-003. Artifacts under
  `project-docs/docs/phases/phase-10/`: validation contract (VAL-001..009) and 3
  ordered task specs (T1 CF-003, T2 matrix/policy-show, T3 retrofit). No schema
  migration. On confirmation, QA Lead refines evidence + writes the QA plan, then
  the dev/QA loop runs T1→T2→T3.
- **CF-003 scheduled into Phase 10 Task 1** (DEC-D: align role discovery with the
  journal corruption contract). CF-008/CF-011/CF-013 stay deferred.
- QA plan must carry the Phase 9 standing checks: show-vs-enforce cross-surface
  parity (VAL-008) and validate-before-append on every permission denial
  (VAL-004/005). Wrap should promote the append-order read-ordering rule toward a
  standing principle (Phase 9 deep-retro rec).
- Standing commitment before release-readiness: decide where held-out scenarios
  live and how implementation agents are prevented from loading them.
</content>
