# Phase 15 Validation Contract - `fid dev *` Composites + Orchestration Cleanup

> Phase: 15 - `fid dev *` composites + CF-013/CF-014/CF-015 cleanup (roadmap
> v0.1e, slice 2 of 4)
> Author: Tech Lead (draft)
> Last refined: 2026-06-28 by Tech Lead (planning draft)
>
> Location: `project-docs/docs/phases/phase-15/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 15 is complete. The contract specifies **what must be proven**; the QA
plan specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Phase Scoping Notes

Phase 15 is the second v0.1e slice: the developer **composite** commands
(`fid dev start/gate/finish/block`) plus the three scheduled carry-forward
cleanups on the same paths (CF-013/CF-014/CF-015). It adds **no new persistence
contract**: the `claims`, `issues`, `agent_runs`, `gate_runs`, and
`task_completions` tables and their triggers already exist
(`migrations/0002_v01_tables.up.sql`; `03_TRD_SCHEMA.md` sections 2, 7). It
introduces the first **multi-effect composite write commands** — each composite
orchestrates several existing single-effect writes inside **one journal event**
(one logical write = one event, DEC-028 / standing constraint).

Scope:

- `fid dev start` (PRD-DEV-008, PRD-CLAIM-005): one event = acquire the caller's
  claim + transition `todo → in_progress` + open one `agent_runs` row
  (`ended_at` NULL).
- `fid dev gate` (PRD-DEV-004, PRD-DEV-010; TRD `02_TRD.md` §4.7): one event =
  append one `gate_runs` row on the derived active attempt, recorded only when
  the issue is in a started-category state and the caller holds the live claim.
- `fid dev finish` (PRD-DEV-008, 009, 010): one event = close the open
  `agent_runs` row + insert a `task_completion` (`dev_complete`, gate snapshot,
  supersede priors) + transition `in_progress → dev_complete` + release the
  caller's claim.
- `fid dev block` (PRD-DEV-008, 010): one event = insert a `task_completion`
  (`failed_gate`, gate snapshot, reason, supersede priors) + transition
  `in_progress → failed_gate` + release the caller's claim.
- The in-code permissions matrix gains `dev.start`, `dev.gate`, `dev.finish`,
  `dev.block` (copied from `03_TRD_SCHEMA.md` section 13: `developer` and
  `human` only).
- **CF-014:** behavior-preserving split of `cmd/fid/issue.go` (984) and
  `internal/journal/replay.go` (485) under the 400-line non-test file target.
- **CF-015:** behavior-preserving function extraction in `internal/labels/add.go`
  (`Add`) and `internal/issues/label.go` (`AddLabel`/`RemoveLabel`) under the
  60-line non-test function target.
- **CF-013:** behavior-preserving extraction of the `internal/issues/complete.go`
  `Complete` orchestration into a reusable completion-effect builder that
  `fid dev finish`/`fid dev block` consume.

## Decisions (confirmed by human 2026-06-29)

Two interpretations add a public-interface surface not fully pinned in the
specs. **Both confirmed by the human on 2026-06-29**, along with PI-1..PI-6;
they are recorded as phase decisions during execution. Developers proceed
without re-escalating.

- **DN-A — `fid dev gate` records a reported result; it does not execute the
  command.** `06_DESIGN_PLAN_CLI.md` shows `fid dev gate --cmd <command> --by`
  and confirmation output with `exit_code` and cumulative `gates_passed`/
  `gates_failed`, but no source for the exit code. Recommended interpretation:
  fid does **not** run the command (GLOSSARY: "Developers run gates during a
  task attempt; standalone `fid dev gate` writes append-only `gate_runs`
  rows"; ADR-0001 deterministic recorder; avoids a new command-execution
  security surface and threat-model obligation under TRD-SEC). Add
  `--exit-code <int>` (required) and `--output-excerpt <string>` (optional) to
  the `fid dev gate` flag table; `gates_passed`/`gates_failed` in the
  confirmation are the cumulative attempt tally (`exit_code == 0` ⇒ passed,
  else failed). **Public-interface addition → confirm at kickoff.** (VAL-005)
- **DN-B — `fid dev start --model <m>` opens the base run.** `agent_runs.model`
  is `NOT NULL` but the CLI spec gives `dev start` no model flag. Recommended
  interpretation: add `--model <string>` (required) to `fid dev start` so the
  opened run records its model; `tokens_in`/`tokens_out`, snapshot prices, and
  `cost_usd` stay NULL on the dev-opened run in v0.1 (the dedicated
  `fid run record` from Phase 14 remains the full-fidelity cost-capture path;
  per-composite cost capture and run lifecycle are v0.2c). `fid dev finish
  --tokens` populates `task_completions.actual_tokens` (per the CLI table), not
  run tokens. **Public-interface addition → confirm at kickoff.** (VAL-003)

## Planning Interpretations (recorded as phase decisions during execution)

- **PI-1 (new `internal/dev` boundary).** The four composites live in a new
  domain package `internal/dev` that orchestrates `internal/claims`,
  `internal/issues`, `internal/runs`, and `internal/workflow` effect builders.
  Record the boundary in `02_TRD.md` sections 2-3. No new persistence contract.
- **PI-2 (one event per composite).** Each composite emits exactly one journal
  event (`dev.started` / `dev.gate.recorded` / `dev.finished` / `dev.blocked`)
  carrying all of its effects, validate-before-append (CF-007): on any rejection
  path the command exits non-zero with an actionable error and appends **zero**
  journal lines. The projector is effect-driven (`project_effects.go`), so the
  new event types need no upgrader; the Developer confirms no event-type
  registry requires the new types.
- **PI-3 (`dev start` claim-in-same-event).** `dev start` acquires the claim and
  transitions into the started category in the same event; PRD-CLAIM-005
  (started-category requires a live caller claim) is satisfied by the claim
  acquired in this composite, so `dev start` does not require a pre-existing
  claim. It targets the `todo → in_progress` edge; re-entry from
  `failed_gate`/`qa_block` is `dev resume` (v0.2). A live claim held by another
  role is rejected (no `--force` in v0.1 `dev start`).
- **PI-4 (`dev finish` rule ordering).** The `in_progress → dev_complete` edge
  requires `all_acs_verified[developer]` and `has_completion_report`. The
  completion is created in the same event, so `has_completion_report` is
  satisfied by construction; `dev finish` evaluates `all_acs_verified[developer]`
  against existing verifications and rejects before append when unmet, and does
  **not** require a pre-existing completion.
- **PI-5 (composite claim release on finish/block).** `dev finish` (→
  `dev_complete`) and `dev block` (→ `failed_gate`) both release the caller's
  live claim as part of the composite, matching the `06_DESIGN_PLAN_CLI.md`
  "claim released" output, even though those target states stay in the started
  category (a plain `fid issue transition` would not auto-release them).
  PRD-CLAIM-006 semantics for plain `fid issue transition` are unchanged.
- **PI-6 (`dev finish` closes the open run; `dev block` leaves runs alone).**
  `dev finish` closes the issue's open `agent_runs` row (set `ended_at`) when one
  exists, matching its "agent run ... closed" output; `dev block` does not manage
  runs (its output shows none). Closing is conditional so finish works even when
  no `dev start` opened a run.

Carry-forward review at kickoff:

- CF-008 remains deferred to v0.2.
- CF-011 (broad `activity_log` reconciliation) remains scheduled for Phase 16.
  The composites project the same per-command `activity_log` rows the underlying
  transition/claim writes already project today (status change; claim takeover);
  they do not widen or backfill CF-011.
- CF-013, CF-014, CF-015 are scheduled here and land in Phase 15 (Task 1 for
  CF-014/CF-015; Task 3 for CF-013 alongside `dev finish`/`dev block`).

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 15, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | CF-014; `08_QUALITY_PLAN.md` size budgets; AGENTS.md | After the split, `cmd/fid/issue.go` and `internal/journal/replay.go` are each ≤400 lines, and every existing `cmd/fid` and `internal/journal` test passes unchanged (no exported-surface or behavior change). | `wc -l` on both files after the split; full task gate green with the pre-existing issue/journal tests unmodified except for moves. | QA Lead | task-1 | Both files ≤400 lines and all pre-existing behavior tests pass with no exported-surface change. | Either file >400 lines, a behavior test is changed/weakened, or any exported symbol/CLI/JSON/journal surface changes. | Draft | Splitting other oversized files is not in scope. |
| VAL-002 | CF-015; `08_QUALITY_PLAN.md` size budgets | After refactor, every non-test function in `internal/labels/add.go` and the `internal/issues/label.go` `AddLabel`/`RemoveLabel` write flows is ≤60 lines, and label add/remove behavior is unchanged (existing label tests pass). | Function line-length inspection of the named functions; full task gate green with pre-existing `internal/labels` and issue-label tests unmodified except for moves. | QA Lead | task-1 | Named functions ≤60 lines and label add/remove behavior unchanged. | A named function exceeds 60 lines or label add/remove behavior/tests change. | Draft | Other oversized functions are out of scope. |
| VAL-003 | PRD-DEV-008; PRD-CLAIM-005; `03_TRD_SCHEMA.md` §2, §7; `06_DESIGN_PLAN_CLI.md` `fid dev start`; DN-B; PI-1/PI-2/PI-3 | `fid dev start <id> --by developer --model <m>` appends exactly one `dev.started` event whose effects (a) insert a live `claims` row for the caller, (b) transition the issue `todo → in_progress`, and (c) insert one open `agent_runs` row (`role`, `model`, `started_at` set, `ended_at` NULL); reading the issue shows `in_progress`, `fid claim list` shows the live caller claim, and the run is open. | CLI + domain: run `dev start` on a fresh `todo` task; assert exactly one journal line, the three effects in the projected DB, issue state, claim list, and open run. | QA Lead | task-2 | One `dev.started` event with claim + `in_progress` transition + one open run; issue/claim/run state correct. | More than one event, a missing/extra effect, the run not opened, or the model unrecorded. | Draft | Git context, run lifecycle `status`, and `dev resume` are v0.2c. |
| VAL-004 | PRD-CLAIM-005; CF-007; PI-2/PI-3 | `fid dev start` rejects **before append** (non-zero exit, actionable error, zero new journal lines) when a live claim is held by another role or the issue is not in a state with a valid `todo → in_progress` edge; on the valid path the claim acquired in the same event satisfies the started-category requirement so no pre-existing claim is needed. | CLI: attempt `dev start` on an issue claimed by another role and on a non-`todo` issue (assert reject + zero new lines); run `dev start` on a fresh unclaimed `todo` task (assert success). | QA Lead | task-2 | Both reject paths exit non-zero with zero new journal lines; the valid unclaimed path succeeds via claim-in-same-event. | A reject path appends a line or returns success, or the valid path demands a pre-existing claim. | Draft | `--force` start is out of scope (v0.1 `dev start` has no force). |
| VAL-005 | PRD-DEV-004; PRD-DEV-010; `02_TRD.md` §4.7; `03_TRD_SCHEMA.md` §7; `06_DESIGN_PLAN_CLI.md` `fid dev gate`; DN-A; PI-2 | `fid dev gate <id> --cmd <c> --exit-code <n> [--output-excerpt <s>] --by developer` appends exactly one `dev.gate.recorded` event with one `gate_runs` row on the derived active attempt (`MAX(task_completions.attempt_number)+1`, 0 when none), recording the reported `command`/`exit_code`/`output_excerpt`/`ran_at`; it rejects before append (zero new lines) when the issue is not in a started-category state or the caller does not hold the live claim. | CLI + domain: start a task, record a gate (assert one `gate_runs` row with the reported fields and correct attempt number); attempt a gate on an unstarted/unclaimed issue (assert reject + zero new lines). | QA Lead | task-2 | One `gate_runs` row on the correct attempt with reported fields; gate rejected before append without a started state + live claim. | fid executes the command, the attempt number is wrong, fields are dropped, or a disallowed gate appends a row. | Draft | `gates_passed`/`gates_failed` are derived at finish/block, not stored on `gate_runs`. |
| VAL-006 | TRD-ARCH-010; `03_TRD_SCHEMA.md` §13; CF-007 | The in-code permissions matrix contains exactly one row each for `dev.start`, `dev.gate`, `dev.finish`, `dev.block` with the same role booleans as section 13 (`developer` and `human` only); `fid policy show` displays those rows from the same source `Authorize` enforces; `developer` is allowed and a denied role (e.g. `qa_reviewer`) is rejected before append with a non-zero exit, an error naming role + action, and zero new journal lines. | Unit + CLI: matrix completeness test for the four rows; `fid policy show` includes them and matches the matrix object; run each composite as `developer` (allowed) and as `qa_reviewer` (rejected before append). | QA Lead | tasks-2-3 | Four rows matching section 13; show matches enforce; developer allowed, denied role rejected before append. | A row is missing/duplicated or diverges from section 13, a permitted role is rejected, a denied role appends, or show/enforce drift. | Draft | Configurable role permissions remain v0.2 (PRD-ROLE-006). |
| VAL-007 | PRD-DEV-008, 009, 010; `03_TRD_SCHEMA.md` §2, §7; `06_DESIGN_PLAN_CLI.md` `fid dev finish`; PI-2/PI-5/PI-6 | `fid dev finish <id> --by developer --completion-md <p> [--tokens --commits --files --human-intervention --intervention-note]` appends exactly one `dev.finished` event whose effects close the issue's open `agent_runs` row (set `ended_at`), insert one `task_completions` row (`dev_complete`, gate snapshot from the active attempt's `gate_runs`, `gates_passed`/`gates_failed`, supersede prior completions), transition `in_progress → dev_complete`, and release the caller's claim. | CLI + domain: after start+verify+gate, run `dev finish`; assert one event, run closed, completion with gate snapshot + counts, issue `dev_complete`, claim released. | QA Lead | task-3 | One `dev.finished` event with run close + completion(+snapshot) + transition + claim release. | More than one event, a missing effect, the run left open, the claim not released, or the gate snapshot absent. | Draft | Optional fields (`--tokens`/`--commits`/`--files`/intervention) populate the completion, not the run. |
| VAL-008 | PRD-WF-007..010; PRD-DEV-008; CF-007; PI-4 | `fid dev finish` enforces the `dev_complete` transition rules: it rejects before append (non-zero exit, zero new journal lines, error naming the unmet rule) when `all_acs_verified[developer]` is unmet; `has_completion_report` is satisfied by the completion created in the same event, so finish succeeds with no pre-existing completion once developer ACs pass. | CLI: run `dev finish` with unverified developer ACs (assert reject naming `all_acs_verified` + zero new lines); verify developer ACs then `dev finish` with no prior completion (assert success). | QA Lead | task-3 | Unmet-AC finish rejected before append naming the rule; verified-AC finish succeeds without a pre-existing completion. | Finish appends despite unverified ACs, or finish demands a pre-existing completion. | Draft | qa_lead-tier AC verification (`qa_pass`) is a later edge, not exercised here. |
| VAL-009 | PRD-DEV-008, 010; `03_TRD_SCHEMA.md` §2, §7; `06_DESIGN_PLAN_CLI.md` `fid dev block`; PI-2/PI-5 | `fid dev block <id> --by developer --reason-md <p>` appends exactly one `dev.blocked` event whose effects insert one `task_completions` row (`failed_gate`, gate snapshot from the active attempt, supersede priors, reason recorded as `notes_md_path`), transition `in_progress → failed_gate`, and release the caller's claim. | CLI + domain: after start (+optional gate), run `dev block`; assert one event, completion(`failed_gate`, snapshot, reason), issue `failed_gate`, claim released. | QA Lead | task-3 | One `dev.blocked` event with completion(`failed_gate`, snapshot, reason) + transition + claim release. | More than one event, a missing effect, wrong completion status, the reason dropped, or the claim not released. | Draft | Concern raising and escalation remain separate commands, not part of `dev block`. |
| VAL-010 | CF-013; `08_QUALITY_PLAN.md` size budgets | The `internal/issues/complete.go` `Complete` orchestration is extracted into a reusable completion-effect builder that `fid dev finish`/`fid dev block` consume; `fid issue complete` behavior is unchanged (existing complete tests pass) and `Complete` plus the extracted helpers each meet the 60-line non-test function target. | Function inspection of `Complete` and the extracted builder; full task gate green with pre-existing `internal/issues` completion tests unmodified except for moves; the builder is referenced by `internal/dev`. | QA Lead | task-3 | Shared builder exists and is reused by `internal/dev`; `fid issue complete` behavior unchanged; functions ≤60 lines. | Behavior of `fid issue complete` changes, the builder is duplicated rather than shared, or functions exceed 60 lines. | Draft | A broader rewrite of the completion path is out of scope; extraction only. |
| VAL-011 | roadmap v0.1e Phase 15 exit; `08_QUALITY_PLAN.md` §5.2 | An end-to-end smoke over a real `.fid` workspace drives `fid dev start → gate → finish` (claim → `in_progress` → open run → gate → completion → `dev_complete` → claim released) and `fid dev start → block` (→ `failed_gate` → completion → claim released), proves each composite emits exactly one event with the correct effects, denies a non-developer role, and rebuilds deterministically from the journal. | Integration/e2e test: full composite lifecycle with permission denial and a delete-and-rebuild equality check. | QA Lead | task-4 | Both composite paths complete with correct one-event effects, denial enforced, and rebuild reproduces the projection. | A composite path fails, emits more than one event, denial is not enforced, or rebuild diverges. | Draft | Activity timeline reads and maintenance audits are Phases 16-17. |
| VAL-012 | `08_QUALITY_PLAN.md` §§5.2, 5.6; phase goal | A human dogfooding a real `.fid` workspace can tell from CLI output that `dev start` claims + transitions + opens a run, `dev gate` records a gate result, `dev finish` completes + closes the run + releases the claim, and `dev block` records a failed gate with a reason. | Human checklist with one row per observable; QA Lead generates from this VAL after deterministic QA passes. | Human | phase-end | Human confirms the four composite outcomes from CLI output without reading implementation. | Human cannot identify those outcomes from CLI output, or the checklist is not run. | Draft | QA Lead fans this into per-observable human-checklist rows. |

## Coverage Rules

Confirmed for Phase 15 execution:

- Phase-scope requirements each map to at least one `VAL-NNN`: PRD-DEV-008 to
  VAL-003/006/007/009/011; PRD-CLAIM-005 to VAL-003/004; PRD-DEV-004 +
  PRD-DEV-010 to VAL-005/007/009; PRD-DEV-009 to VAL-007; PRD-WF-007..010 to
  VAL-008; TRD-ARCH-010 / section 13 to VAL-006; CF-014 to VAL-001; CF-015 to
  VAL-002; CF-013 to VAL-010; phase exit smoke to VAL-011; dogfood goal to
  VAL-012.
- Every `VAL-NNN` is referenced by a task spec or explicitly marked `phase-end`
  (VAL-012).
- Binary/conditional assertions name both states: VAL-001, VAL-002, VAL-004,
  VAL-005, VAL-006, VAL-008.
- VAL-006 carries a matrix-vs-section-13 and show-vs-enforce parity proof;
  VAL-011 carries a delete-and-rebuild determinism proof; VAL-001/002/010 carry
  behavior-preservation proofs (pre-existing tests unmodified except for moves).
- Coverage is enforced at the gate, not as a `VAL-NNN`: `internal/dev`,
  `internal/issues`, `internal/labels`, `internal/policy`, and `cmd/fid`
  coverage must stay above thresholds; each task runs the required coverage gate
  and records package percentages.
- The phase touches `cmd/fid` and adds composite paths exercised by
  integration-tagged tests, so the required integration gate
  (`go test -tags integration -race ./...`) must be recorded PASS (or BLOCKED
  with reason) at every task that adds or moves integration-tagged tests.
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase review, or
  `(human)` verification, and fans VAL-012 into per-observable human rows.
- Standing QA-plan checks carry forward: validate-before-append on every
  composite rejection path (VAL-004/005/006/008), one event per composite
  (VAL-003/005/007/009/011), no stale-projection nested-lock failure on the new
  write paths, and behavior-preservation for the three cleanups (VAL-001/002/010).

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-28 | VAL-001..012 | created | Phase 15 planning draft | roadmap v0.1e Phase 15; PRD-DEV-004/008-010; PRD-CLAIM-005; PRD-WF-007..010; TRD-ARCH-010; CF-013/CF-014/CF-015 |
