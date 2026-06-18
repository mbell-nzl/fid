# Progress - Current

> **Purpose.** In-flight phase status only. Loaded by default in every agent
> session. Historical phase progress lives under `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 7 (Acceptance Criteria Surface, v0.1c first slice) is
**CLOSED** on branch `phase-7`. QA final review: **PASS WITH IMPROVEMENTS**
(0 blocking / 0 major / 3 improvement); human dogfood checklist 3/3 PASS
(HV-010a/b/c). VAL-001..VAL-010 all **PASS** (VAL-010 resolved by human
dogfood). Phase 6 closed v0.1b; phases 0-6 archived under `../phases/phase-N/`.

Phase 7 delivered:

- `fid issue ac add` — ordered ACs on tasks/subtasks, epic rejection before
  append (Task 1, VAL-001/002).
- `fid issue ac verify` — multi-role append-only verification with evidence/
  status enums, `completion_id` NULL (Task 2, VAL-003/005).
- Deterministic latest-per-role by journal append order (`rowid` DESC) +
  `ac.go` 190 / `ac_verify.go` 255 split (Task 2b, DEC-021/DEC-022, VAL-004/009).
- `fid issue ac list` — latest-per-role, four-format dispatch + JSON default
  (Task 3, VAL-006/007).
- Real-data `all_acs_verified` reject→verify→pass + delete-DB `fid rebuild`
  parity (Task 4, VAL-008/009).

Five task units, one BLOCK→resolved (TRG-001/ESC-005 → DEC-021), one scope
correction (DEC-022 → CF-012), zero dependency changes. DEC-020/021/022 recorded
at wrap.

Coverage held: `cmd/fid` 78.9%, `internal/issues` 82.0%, `internal/journal`
80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5% (COV-001/002 satisfied).

Open carry-forward register: CF-003 (role corrupt-tail, before v0.1d), CF-008
(resolution type/reason, v0.2), CF-011 (consistent `activity_log`, v0.1e),
CF-012 (`issue.go` 677-line split, next v0.1c slot).

## Next Action

- **Phase 8 is the next planned cut** (roadmap §v0.1c — task completions:
  `fid issue complete`, append-only `task_completions`, attempt numbering +
  supersession, `gate_runs`→`gate_results` snapshot via journal fixtures, real
  `has_completion_report`-gated transition). Not yet planned — Tech Lead phase
  start creates the validation contract and task list.
- Consider folding **CF-012** (`internal/issues/issue.go` behavior-preserving
  split) into the Phase 8 slot per its disposition.
- Standing commitment before release-readiness: decide where held-out scenarios
  live and how implementation agents are prevented from loading them.
