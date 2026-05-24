# Phase Brief - Phase <N> (<state>)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Replaced at each phase kickoff by the incoming Tech Lead.
>
> **Size target:** under 150 lines.

## Phase Goal

One paragraph: what this phase delivers, why it matters now, and what success
looks like in user-visible or operator-visible terms.

## Status

**<Planning | Active | Closing | Approved>.** One paragraph: current branch,
current state, and next concrete action.

## Task Index

Full task specs and completion reports live under
`../phases/phase-<N>/tasks/`. This is an index; do not duplicate full task
content here.

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | <title> | <Pending/Active/Done> | <ids or None> |

## Validation Contract

The full contract lives at
`../phases/phase-<N>/qa/validation-contract.md`. This brief is only a compact
index for active phase execution.

| VAL ID | Source | Owning Task/Window | Validator | Status |
|---|---|---|---|---|
| VAL-001 | <requirement or decision> | <task/window> | <QA Lead / QA Reviewer / Human> | <Draft/Ready/PASS/FAIL/BLOCKED/CARRY-FORWARD> |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| QA plan | `../phases/phase-<N>/qa/qa-plan.md` | <Missing/Draft/Ready> |
| Mid-phase reviews | `../phases/phase-<N>/qa/review-M.md` | <cadence or N/A> |
| Final review | `../phases/phase-<N>/qa/final-review.md` | <Pending/Done> |
| Human checklist | `../phases/phase-<N>/qa/human-checklist.md` | <Pending/Done> |

## Active Constraints

Phase-specific constraints in force right now. Standing rules live in
`AGENTS.md` and `decisions_principles.md`; do not duplicate them here.

- <constraint>

## Parity Matrix

Use this section only when multiple surfaces or outputs must stay consistent,
or when the phase is a polish/cross-surface audit phase. Otherwise delete it.

| Axis | Surface A | Surface B | Surface C | Evidence Source | Notes |
|---|---|---|---|---|---|
| Terminology | <expected> | <expected> | <expected> | <command/test/artifact> | <notes> |
| State classification | <expected> | <expected> | <expected> | <command/test/artifact> | <notes> |
| Error behavior | <expected> | <expected> | <expected> | <command/test/artifact> | <notes> |
| Layout / shape / schema | <expected> | <expected> | <expected> | <command/test/artifact> | <notes> |

## Interfaces / Subsystems In Scope

- `<path or surface>` - <reason>

## Known Risks

- **<risk>.** <Description and mitigation or trigger.>

## Open Questions

- [ ] <question>

## Success / Exit Criteria

- [ ] All planned tasks complete with green quality gates.
- [ ] Final QA review has no unresolved blocking findings.
- [ ] Human verification accepted or follow-up work planned.
- [ ] Phase sign-off decision recorded.

## Next Planning Checkpoint

<What happens after this phase.>
