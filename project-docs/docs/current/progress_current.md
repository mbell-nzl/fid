# Progress - Current

> **Purpose.** In-flight phase status only. Loaded by default in every agent
> session. Historical phase progress lives under `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 6 (Ready Queue + v0.1b Closure) is **CLOSED** on branch
`phase-6`. QA final review: **DETERMINISTIC PASS** (0 blocking / 0 major / 1
improvement); human dogfood checklist 3/3 PASS (HV-012a/b/c). VAL-001..VAL-012
all **PASS**. This closes **v0.1b**. Phase 5 and earlier are archived under
`../phases/phase-N/`.

Phase 6 delivered:

- `fid issue update` one-field-per-line confirmation output (Task 1, **CF-009
  closed**, VAL-009).
- File-size split: `transition.go` 262, `update.go` 260, `project.go` 137 lines,
  behavior-preserving + rebuild parity (Task 2, **CF-010 closed**, VAL-010).
- `internal/issues` ready selection engine — sequential/prioritised/unattached/
  project-wide with status-category readiness (Task 3, DEC-017, VAL-002..006).
- `fid ready` CLI: `--epic`, `--format` dispatch, frozen `{"task": null}` empty
  state + stderr `note:` (Task 4, DEC-018/DEC-019, VAL-001/007/008).
- End-to-end v0.1b smoke over create→list→get→update→transition→ready with
  delete-DB + `fid rebuild` parity (Task 5, VAL-011) + held-out scenario seed
  recommendation.

Five tasks, one BLOCK→resolved (TRG-001/ESC-004 → DEC-019), zero spec
deviations, zero dependency changes. DEC-017/DEC-018 formally recorded at wrap.

Coverage held: `cmd/fid` 78.3%, `internal/issues` 82.5%, `internal/journal`
80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5% (COV-001/002 satisfied).

Open carry-forward register: CF-003 (role corrupt-tail, before v0.1d), CF-008
(resolution type/reason, v0.2). CF-009/CF-010 closed this phase.

## Next Action

- **Phase 7 is PLANNED** (roadmap §v0.1c, first slice — acceptance-criteria
  surface only). Validation contract (VAL-001..010) and 4 ordered Developer tasks
  drafted under `../phases/phase-7/`; brief refreshed. v0.1c is split into Phase 7
  (ACs), Phase 8 (completions), Phase 9 (concerns + derived flag) per the roadmap
  §v0.1c slicing note. Proposed **DEC-020** (AC writes don't emit `activity_log`,
  consistent with `fid issue create`) awaits human confirmation.
- Next: QA Lead refines the Phase 7 QA plan and runs the dev/QA loop from Task 1.
- Standing commitment before release-readiness: decide where held-out scenarios
  live and how implementation agents are prevented from loading them
  (`held-out-recommendation.md` recommends the v0.1b smoke as the seed).
