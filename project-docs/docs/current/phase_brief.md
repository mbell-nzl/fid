# Phase Brief - Between Phases (Phase 6 CLOSED → Phase 7 to plan)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 6 (Ready Queue + v0.1b Closure) is **CLOSED** — QA
> DETERMINISTIC PASS; human dogfood 3/3; VAL-001..012 all PASS; **v0.1b complete**.
> Archives under `../phases/phase-6/`. **Phase 7 is not yet planned** — the Tech
> Lead must draft the validation contract and ordered task list before execution.

## Phase 6 outcome (closed)

`fid ready` + v0.1b closure shipped. CF-009 and CF-010 closed. DEC-017
(status-category readiness), DEC-018 (frozen `{"task": null}` empty state), and
DEC-019 (`priority ASC` nil-last → `created_at` → `id` selection order) are
recorded in `../phases/phase-6/decisions.md`. Coverage held (COV-001/002).

## Next phase target (roadmap §v0.1c — ACs, Completions, Concerns)

Phase 7 implements the v0.1c write surface:

- `fid issue ac add` / `ac verify` / `ac list` (PRD-AC-001..006) — ordered
  acceptance criteria, multi-role append-only verification, latest-status list.
- `fid issue complete` (PRD-DEV-001..004, 009..010) — append-only
  `task_completions` row with attempt number, tokens, commits, files touched,
  `gate_runs` snapshot, human-intervention flag; supersede prior attempts.
- `fid concern raise` / `address` / `list` (PRD-DEV-005..007).
- `issues_enriched` derived `concerns_raised` / `open_concerns_count` (PRD-DEV-007).

## Carry-Forward Register (open)

| ID | Item | Disposition |
|---|---|---|
| CF-003 | Role-validation corrupt-tail tolerance | Before v0.1d role/permission work. |
| CF-008 | Resolution type / cancellation reason (DEC-015) | v0.2 (product). |

CF-009 and CF-010 closed at Phase 6 wrap.

## Standing commitments to honor at Phase 7 planning

- Create `../phases/phase-7/qa/validation-contract.md` with `VAL-NNN` assertions
  before task execution begins (one observable each; binary/conditional
  assertions name both states; parity/"matches reference" needs deterministic
  conformance evidence).
- DEC-007 carry-in: COV-001 (`internal/issues`, `internal/journal`,
  `internal/policy`, `internal/workflow` ≥ 80%); COV-002 `cmd/fid` ≥ 70%.
  Phase 6 ended: `cmd/fid` 78.3%, `internal/issues` 82.5%, `internal/journal`
  80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5%.
- Before release-readiness work, decide where held-out scenarios live and how
  implementation agents are prevented from loading them
  (`../phases/phase-6/held-out-recommendation.md` recommends the v0.1b smoke seed).
- Process nudge (Phase 6 workflow note): require pasted exact command transcripts
  for any public-output change in completion reports — fold into task specs to
  avoid repeat evidence-format findings.

## Next Checkpoint

Tech Lead: read the roadmap §v0.1c entry and relevant PRD/TRD/design sections,
review the carry-forward register, then draft the Phase 7 validation contract
and ordered Developer task list. Present both to the human before execution.
