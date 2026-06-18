# Phase Brief - Between phases (Phase 7 CLOSED, Phase 8 to plan)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 7 (Acceptance Criteria Surface, v0.1c first slice) is
> **CLOSED** on branch `phase-7`. QA final review **PASS WITH IMPROVEMENTS**
> (0 blocking / 0 major / 3 improvement); human dogfood 3/3 PASS. VAL-001..010
> all PASS. Phase 8 is **not yet planned** — the next action is a Tech Lead
> phase start. No active task work is in flight.

## Phase 7 outcome (closed)

Shipped the acceptance-criteria write surface on the existing
`acceptance_criteria` / `ac_verifications` schema and replaced the Phase 5
journal-fixture evidence for `all_acs_verified` with real write commands:

- `fid issue ac add` — ordered ACs on tasks/subtasks; epic rejection before
  append (VAL-001/002).
- `fid issue ac verify` — multi-role, append-only, evidence/status enums;
  `completion_id` NULL (VAL-003/005).
- `fid issue ac list` — latest-per-role, four-format dispatch + JSON default
  (VAL-006/007).
- Real-data `all_acs_verified` reject→verify→pass + delete-DB `fid rebuild`
  parity (VAL-008/009).
- Latest-per-role resolved by journal append order (`rowid` DESC), `verified_at`
  display-only (DEC-021).

Archives: `../phases/phase-7/` (`decisions.md`, `progress.md`, `triage-log.md`,
`tl-escalations.md`, `qa/`).

## Next phase (Phase 8 — Task completions, roadmap §v0.1c)

To be planned at Tech Lead phase start. Scope per roadmap:

- `fid issue complete` — append-only `task_completions` row: attempt number,
  tokens, commits, files touched, human-intervention flag/note, prior-attempt
  supersession, and the `gate_runs`→`gate_results` snapshot mechanism.
- `fid dev gate` ships in v0.1e, so `gate_runs` rows are seeded via journal
  fixtures (Phase 5 pattern) to exercise the snapshot until the write exists.
- Wire real completions into the existing `has_completion_report` transition
  rule (replace Phase 5 fixtures). **Exit:** an agent can record a completion
  attempt, observe attempt numbering + supersession, see the gate snapshot
  rolled into the completion, flag a human intervention, and satisfy a real
  `has_completion_report`-gated transition.

At phase start: draft `../phases/phase-8/qa/validation-contract.md`, break into
ordered one-session tasks, and consider folding CF-012 into this slot.

## Carry-Forward Register (open)

| ID | Item | Disposition |
|---|---|---|
| CF-003 | Role-validation corrupt-tail tolerance | Before v0.1d role/permission work. |
| CF-008 | Resolution type / cancellation reason (DEC-015) | v0.2 (product). |
| CF-011 | Consistent `activity_log` timeline population (DEC-020) | v0.1e (with `fid activity list`). |
| CF-012 | `internal/issues/issue.go` 677-line tripwire split (DEC-022) | Next v0.1c slot; behavior-preserving (à la CF-010). |

## Standing commitments to honor

- Coverage (DEC-007): COV-001 `internal/issues|journal|policy|workflow` ≥ 80%;
  COV-002 `cmd/fid` ≥ 70%. Phase 7 ended `cmd/fid` 78.9%, `internal/issues`
  82.0%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow`
  95.5%.
- Validate-before-append (CF-007): all rejections occur before journal append
  (non-zero exit, clear error, zero new lines/rows).
- DEC-020: AC writes (and `fid issue create`) do not emit `activity_log`;
  completions follow the same precedent unless a phase decision says otherwise —
  consistent timeline population is CF-011 (v0.1e).
- Before release-readiness: decide where held-out scenarios live and how
  implementation agents are prevented from loading them
  (`../phases/phase-6/held-out-recommendation.md` recommends the v0.1b smoke
  seed; Phase 7's AC gating smoke may extend it — recommend only).
- Process nudge: require pasted exact command transcripts for any public-output
  change in completion reports — fold into each task spec.
- Task-spec authoring nudge (Phase 7 lesson): scope file-size acceptance criteria
  to the task's allowed write set, not the whole package (avoids DEC-022-style
  scope BLOCKs).

## Next Checkpoint

Tech Lead: start Phase 8 — read roadmap §v0.1c, draft the validation contract,
break into ordered Developer tasks, refresh this brief, and present the contract
and task list to the human.
