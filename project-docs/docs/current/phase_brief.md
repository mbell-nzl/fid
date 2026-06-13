# Phase Brief - Phase 4 Closed / Phase 5 Pending Planning

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 4 (v0.1b Issue authoring and discovery) is **CLOSED** (QA
> PASS, human checklist 5/5 PASS). Phase 5 is **NOT YET PLANNED** — a fresh Tech
> Lead planning pass owns the Phase 5 validation contract and task list before
> any implementation begins.

## Phase 4 Outcome (for reference)

Phase 4 made issue authoring and discovery usable end to end: sequence
scheduling backstop, issue authoring domain service, `fid issue create`,
centralised enriched `fid issue get` + `--section`, and `fid issue list`
filters/order/formats. VAL-001..VAL-010 all PASS. Archives under
`../phases/phase-4/` (`decisions.md`, `progress.md`, `qa/final-review.md`,
`qa/human-checklist.md`, `tl-escalations.md`).

Closed at wrap: CF-004 (sequence backstop), CF-006 (read-model centralisation),
ESC-002 (due-date read-back).

## Phase 5 Scope (roadmap, to be sliced at planning)

Per `10_IMPLEMENTATION_ROADMAP.md` v0.1b slicing note, Phase 5 is
**mutation + workflow transitions**:

- `fid issue update` (title, priority, labels, spec_md, est_tokens, due_date,
  owner role/actor) with `activity_log` projection.
- `activity_log` projection discipline / `activity` effect op (**CF-005**).
- `internal/workflow` hard-coded graphs and rule functions
  (`all_acs_verified`, `has_completion_report`, `concerns_addressed`,
  `all_child_tasks_done_or_cancelled`).
- `fid issue transition --to <state> --by <role>` with rich rejection errors.
- **CF-007 write-path fixup** (see below).

## Carry-Forward Into Phase 5

| ID | Item | Priority |
|---|---|---|
| CF-005 | `activity_log` projection guard trigger + `activity` effect op; belongs with `fid issue update`. | Medium |
| CF-007 | Rejected `fid issue create` poisons the journal tail — `internal/issues.Create` appends before projecting, so a trigger-rejected create still commits and breaks later reads. Conformance gap vs TRD §4.3 (validate before append). Fix: validate (or projection-dry-run) before append; add a rejected-write journal-hygiene VAL assertion + regression test (line-count unchanged) + poisoned-tail read-recovery test. **Needs human priority confirmation.** | High |

CF-003 (role corrupt-tail tolerance) remains targeted before v0.1d.

## Planning Reminders For Phase 5

- Apply DEC-007 carry-in: package-level tests for COV-001 packages
  (`internal/issues`, `internal/journal`, `internal/policy`); hold COV-002
  `cmd/fid` >= 70%.
- Phase 4 deep-retro lesson (`../phases/phase-4/decisions.md`): when a write
  path can be rejected by a projection-time trigger, the validation contract
  must assert journal hygiene (line-count-unchanged + read-recovery) for the
  rejected state, not only an actionable CLI error. Bake this into the CF-007
  and transition-rejection assertions.
- Transition-rule evidence may seed AC/completion/concern rows via journal
  fixtures until the v0.1c write commands ship.

## Next Planning Checkpoint

Tech Lead opens Phase 5 planning: read the roadmap entry and spec sections,
review the Carry-Forward Register, create `../phases/phase-5/` directories,
draft `qa/validation-contract.md`, slice ordered tasks, and present the contract
and task list to the human. Confirm CF-007 priority with the human first.
