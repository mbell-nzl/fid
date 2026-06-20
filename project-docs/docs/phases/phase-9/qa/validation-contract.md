# Phase 9 Validation Contract - Concerns + v0.1c Closure

> Phase: 9 - Concerns + derived flag + v0.1c closure (roadmap §v0.1c)
> Author: Tech Lead (draft)
> Last refined: 2026-06-20 by Tech Lead (awaiting QA Lead evidence refinement)
>
> Location: `project-docs/docs/phases/phase-9/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 9 is complete. The contract specifies **what must be proven**; the QA plan
specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Phase scoping notes

Phase 9 ships the **concern write surface** (`fid concern raise`,
`fid concern address`, `fid concern list`) on top of the existing `concerns`
schema, the `issues_enriched` derived columns, and the concern read path — all of
which already shipped in earlier phases:

- The `concerns` table, `idx_concerns_issue_status` index, and the
  `issues_enriched` view's `concerns_raised` / `open_concerns_count` columns
  shipped in Phase 2 (v0.1a schema-widening).
- `concerns_raised` and `open_concerns_count` **already surface today** on
  `fid issue get` (`internal/issues/get.go`) and `fid issue list`
  (`internal/issues/list.go`, including the `--has-concerns` filter), and
  `fid issue get --section concerns` already reads concern rows
  (`internal/issues/detail.go` `issueConcerns`).
- The `concerns_addressed` rule function (`internal/workflow/rules.go`
  `ConcernsAddressed`) and its DB loader (`internal/issues/transition_rules.go`
  `loadConcernStatus` + the `RuleConcernsAddressed` dispatch case) already exist.

So Phase 9 adds **only the CLI write/read commands** and proves the existing
read path + rule with **real writes**, replacing the absence of a real concern
producer. This mirrors Phase 7 (ACs) and Phase 8 (completions) exactly.

Planning decisions confirmed by the human (2026-06-20): DEC-A..E.

- **DEC-A (confirmed) — PRD-WF-009 proven at the rule-evaluation seam, no default
  edge.** No v0.1 default workflow edge gates on `concerns_addressed` (TRD §4.6;
  `internal/workflow/workflow.go` matches the spec graph verbatim). The rule is
  *supported* (function + dispatch wired, pure-function unit-tested in Phase 5)
  but unreachable through any `fid issue transition` edge, and there were never
  journal fixtures for it. Phase 9 therefore proves PRD-WF-009 by exercising the
  rule-evaluation seam (`evaluateTransitionRule` → `loadConcernStatus` →
  `ConcernsAddressed`) with **real** concern data: an open concern yields
  *unsatisfied* naming the concern ID(s); addressing all concerns yields
  *satisfied*. Attaching `concerns_addressed` to a workflow edge would change the
  TRD §4.6 default graph (public workflow contract) and is **out of scope** —
  deferred to v0.2 configurable workflows or a separate product decision.
- **DEC-B (confirmed) — concern domain lives in `internal/issues`; no new
  boundary.** Concern writes/reads are implemented in `internal/issues`
  (new files, e.g. `concern_raise.go`, `concern_address.go`, `concern_list.go`)
  with CLI wiring in `cmd/fid/concern*.go`. The concern read model
  (`issueConcerns`) and the rule loader (`loadConcernStatus`) already live in
  `internal/issues`, so this adds **no new architecture boundary** (AGENTS.md /
  TRD §2) and needs no boundary decision. The CLI noun `fid concern` is a
  command namespace only.
- **DEC-C (confirmed) — validate-before-append linkage rules (CF-007).**
  `fid concern raise` rejects, before any journal append, when: the role is
  unknown; the `--on` issue does not exist; the `--in-completion` completion does
  not exist; or the completion does not belong to the `--on` issue (the FK only
  enforces existence, but a cross-issue link is nonsensical). `fid concern
  address` rejects, before any append, when: the role is unknown; the concern
  does not exist; the concern is not in `open` status; or the `--in-completion`
  completion does not exist (and, if present, does not belong to the concern's
  issue). Every rejection is non-zero exit, a clear error, and zero new journal
  lines/rows.
- **DEC-D (confirmed) — address is a status mutation; raise→address only.**
  `fid concern address` updates the same concern row in place
  (`status` `open`→`addressed`, sets `addressed_in_completion_id` and
  `updated_at`); the row is retained. Only the `raise`→`address` lifecycle is
  modelled in Phase 9. The `escalated` / `wontfix` statuses named in PRD-DEV-006
  are schema-valid but have **no Phase 9 write command** (deferred to v0.1e /
  later per the roadmap).
- **DEC-E (confirmed) — no `activity_log` for concern writes (DEC-020 / CF-011
  precedent).** Concern raise/address follow the `fid issue create` / AC-write /
  completion precedent: no `activity_log` row. Consistent timeline population is
  deferred to v0.1e via **CF-011**. No activity assertion in this contract by
  design.

Permissions are `RequireKnownRole` only (the permission matrix is v0.1d),
matching Phase 7/8.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 9, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | PRD-DEV-005; design §6 `fid concern raise` | `fid concern raise --on <task/subtask> --by <role> --in-completion <TC> --summary <text> [--body-md <path>]` creates a `CON-<hash>` `concerns` row storing `issue_id`, `raised_in_completion_id`, `raised_by_role`, `summary`, `body_md_path`, and `status='open'`; read-back via `fid concern list` and `fid issue get --section concerns --format json` shows the concern with those fields. | CLI: complete a task, raise a concern with all flags; read back via `fid concern list --format json` and `fid issue get --section concerns`; assert `CON-` ID shape and each stored field. | QA Lead | task-1 | Concern row persists with `CON-` ID, `status='open'`, and all named fields stored/read-back correctly. | A named field is dropped/mis-stored, wrong ID shape, wrong initial status, or no row persisted. | Draft | `escalated`/`wontfix` raise paths and triage linkage (`triaged_in_event_id`) are out of scope (v0.1e/v0.2). |
| VAL-002 | PRD-DEV-005; CF-007; DEC-C | `fid concern raise` is rejected **before any journal append** (non-zero exit, clear error, zero new journal lines/rows) for each of: unknown role; nonexistent `--on` issue; nonexistent `--in-completion`; a completion that does not belong to `--on`. A fully valid raise succeeds with exactly one new concern row and one new journal line. | CLI over each invalid input asserting zero new journal lines + non-zero exit + clear error; plus the valid success case (one row, one line). | QA Lead | task-1 | Every invalid input is rejected with a clear error and zero append; the valid case appends exactly one event/row. | Any invalid input appends a journal line/row, succeeds, or fails opaquely; or the valid case is rejected. | Draft | Schema FK/kind triggers are a backstop; the CLI must reject first (no journal poisoning). |
| VAL-003 | PRD-DEV-007 | Raising an open concern sets `concerns_raised=true` and increments `open_concerns_count` for the target issue on **both** `fid issue get` and `fid issue list`; before any concern exists the same reads show `concerns_raised=false` and `open_concerns_count=0`. | CLI: read the issue before (false/0) and after (true/count) a real `fid concern raise`, via `fid issue get --format json` and `fid issue list --format json`. | QA Lead | task-1 | Both reads show false/0 before and true/incremented after a real raise. | Either read fails to flip to true / increment, or the pre-raise state is not false/0. | Draft | The derived columns already exist (Phase 2/4); this proves the round-trip with a real writer. |
| VAL-004 | PRD-DEV-006; design §6 `fid concern address`; DEC-D | `fid concern address <CON> --in-completion <TC> --by <role>` updates the concern in place `status` `open`→`addressed`, sets `addressed_in_completion_id` and `updated_at`, and retains the row; read-back via `fid concern list` / `fid issue get --section concerns` shows `status='addressed'` with the addressing completion recorded. | CLI: raise then address a concern; read back and assert `status='addressed'`, `addressed_in_completion_id` set, row still present. | QA Lead | task-2 | Concern row updated to `addressed` with addressing completion recorded and row retained. | Status not changed, addressing completion dropped, row deleted/duplicated, or a new row created instead of in-place update. | Draft | `escalated`/`wontfix` transitions are out of scope (v0.1e+). |
| VAL-005 | PRD-DEV-006; CF-007; DEC-C/DEC-D | `fid concern address` is rejected **before any journal append** (non-zero exit, clear error, zero new lines/rows) for each of: unknown role; nonexistent concern; concern not in `open` status (e.g. already addressed); nonexistent `--in-completion`. A valid address succeeds with exactly the single status mutation. | CLI over each invalid input asserting zero new journal lines + non-zero exit + clear error; plus the valid success case. | QA Lead | task-2 | Every invalid input is rejected with a clear error and zero append; the valid case applies exactly one status mutation. | Any invalid input mutates state/appends, or the valid case is rejected; re-addressing an addressed concern succeeds. | Draft | Idempotent re-address is intentionally rejected (only `open`→`addressed`). |
| VAL-006 | PRD-DEV-007 | Addressing the **last** open concern clears `concerns_raised` to `false` and `open_concerns_count` to `0` on `fid issue get` and `fid issue list`; when another open concern remains, the flag stays `true` and the count reflects the remaining open concerns. | CLI: with two open concerns, address one (flag stays true, count drops by one) then the other (flag false, count 0), asserting both reads each time. | QA Lead | task-2 | Both the partial state (true / remaining count) and the cleared state (false / 0) are observed. | The flag/count does not drop on partial address, or does not clear to false/0 when the last open concern is addressed. | Draft | Proves the derived flag clears via a real writer, not a fixture. |
| VAL-007 | PRD-WF-009; DEC-A | The `concerns_addressed` rule, evaluated against **real** concern data through the transition rule-evaluation seam (`evaluateTransitionRule` → `loadConcernStatus` → `ConcernsAddressed`), reports **unsatisfied and names the open concern ID(s)** while any open concern exists and **satisfied** once all concerns are addressed — driven by real `fid concern raise`/`address` writes, no journal fixtures. | Integration test: seed concerns via real raise/address commands, evaluate the rule via the seam in both states; assert unsatisfied+named IDs (open) and satisfied (all addressed). | QA Lead | task-2 | Rule reports unsatisfied naming open concern IDs while open, and satisfied after all addressed, with real concern data. | Rule reports satisfied while an open concern exists, omits the concern ID, or reports unsatisfied after all are addressed. | Draft | No v0.1 default workflow edge gates on this rule (TRD §4.6); attaching an edge is out of scope (v0.2 configurable workflows). |
| VAL-008 | PRD-DEV-006/007; design §6 `fid concern list` | `fid concern list` returns concerns filtered by `--issue` and by `--status`, with the filters composing (issue + status together), across all read formats (`json` default, `toon`, `markdown`, `yaml`); a filter that matches nothing returns an empty collection (not an error / non-zero exit). | CLI: seed concerns across two issues/statuses; assert `--issue`, `--status`, and combined filters return the right subset in each format; assert empty-match returns an empty collection with exit 0. | QA Lead | task-3 | Each filter (and the combination) returns the correct subset in all four formats; empty match is an empty collection, exit 0. | A filter returns the wrong subset, a format diverges/omits fields, or empty match errors / non-zero exit. | Draft | `escalated`/`wontfix` status values are listable if present but have no Phase 9 producer. |
| VAL-009 | ADR-0001; roadmap §v0.1c exit | Deleting the SQLite DB and running `fid rebuild` reproduces identical `concerns` state (`status`, `raised_in_completion_id`, `addressed_in_completion_id`) and identical derived `concerns_raised` / `open_concerns_count` on issue reads, replayed from the journal. | Integration: snapshot concern + derived-flag state, delete DB, `fid rebuild`, assert identical projected state. | QA Lead | task-3 | Post-rebuild concern rows and derived flags are identical to pre-delete state. | Any divergence in `concerns` rows, addressing linkage, or derived flags after rebuild. | Draft | Held-out scenario harness location remains a separate scheduled commitment. |
| VAL-010 | roadmap §v0.1c exit; phase brief | An end-to-end v0.1c smoke exercises the full slice on one real `.fid/` workspace — `ac add` → `ac verify` → `ac list` → `issue complete` → `concern raise` (derived flag lights up) → `concern address` (derived flag clears) → `concern list` — with each step's observable outcome asserted. | Integration e2e test driving the real CLI/command path end-to-end over the v0.1c slice. | QA Lead | task-3 | The full v0.1c slice runs end-to-end with every asserted step observable and correct. | Any step in the slice fails, or the derived flag does not light up/clear at the raise/address steps. | Draft | Extends the existing `cmd/fid/e2e_v01c_*_test.go` coverage to close the slice. |
| VAL-011 | phase brief; roadmap §v0.1c exit | Human dogfood confirms the concern lifecycle reads and behaves as expected on a real `.fid/` workspace: raise a concern against a completion, see the derived flag light up on issue reads, address it and watch the flag clear, and list concerns by issue/status. | Human checklist against a real `.fid/` workspace; record outcomes per observable. | Human | phase-end | Human confirms raise creates a visible open concern, the derived flag lights up, address clears it, and list filters read as expected. | Human finds a concern not created/visible, a flag that does not light/clear, or list behavior that contradicts expectation. | Draft | QA Lead fans this into one human-checklist row per observable (raise+visible / flag lights / address clears flag / list filters). |

## Coverage Rules

Confirmed before task execution begins:

- Phase-scope requirements each map to ≥1 `VAL-NNN`: PRD-DEV-005→VAL-001/002;
  PRD-DEV-006→VAL-004/005/008; PRD-DEV-007→VAL-003/006/008; PRD-WF-009→VAL-007;
  journal/rebuild canonical (ADR-0001)→VAL-009; v0.1c end-to-end closure→VAL-010;
  human concern-lifecycle feel→VAL-011. CF-007 validate-before-append→VAL-002/005.
- Every `VAL-NNN` is referenced by a task spec or marked `phase-end` (VAL-011).
- Binary/conditional assertions name both states — VAL-002, 003, 005, 006, 007, 008.
- DEC-E (no `activity_log` for concern writes; matches `create`/AC/completion) is
  human-confirmed; consistent timeline population is deferred to v0.1e via CF-011.
  No activity assertion in this contract by design.
- DEC-A (PRD-WF-009 proven at the rule-evaluation seam; no default edge gates on
  `concerns_addressed`) is human-confirmed; attaching a workflow edge is out of
  scope (v0.2 configurable workflows).
- `escalated`/`wontfix` concern lifecycle (PRD-DEV-006 statuses beyond
  raise→address) has **no Phase 9 assertion by design** — deferred to v0.1e/later
  per the roadmap §v0.1c phase-slicing note.
- Coverage (DEC-007) is enforced at the gate, not as a `VAL-NNN`: COV-001
  `internal/issues|journal|policy|workflow` ≥ 80%; COV-002 `cmd/fid` ≥ 70%. Each
  task runs the coverage gate and the completion report records the percentages.
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase-review, or
  `(human)` verification and fans VAL-011 into per-observable human rows.

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-20 | VAL-001..011 | created | Phase 9 planning draft | phase brief; roadmap §v0.1c; PRD-DEV-005..007; PRD-WF-009; ADR-0001; CF-007 |
| 2026-06-20 | n/a | DEC-A..E confirmed | Human accepted Tech Lead recommendations pre-execution | this contract §Phase scoping notes |
