# Phase 7 Validation Contract - Acceptance Criteria Surface

> Phase: 7 - Acceptance criteria surface (v0.1c)
> Author: Tech Lead (draft)
> Last refined: 2026-06-14 by Tech Lead (awaiting QA Lead evidence refinement)
>
> Location: `project-docs/docs/phases/phase-7/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 7 is complete. The contract specifies **what must be proven**; the QA
plan specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Phase scoping notes

Phase 7 ships the **acceptance-criteria write surface** (`fid issue ac add`,
`fid issue ac verify`, `fid issue ac list`) on top of the existing
`acceptance_criteria` / `ac_verifications` schema and read path. It replaces the
Phase 5 journal-fixture evidence for the `all_acs_verified` transition rule with
real write commands. Decisions to confirm with the human at planning:

- **DEC-020 (confirmed by human) — AC writes do not emit `activity_log` in Phase
  7.** AC add and AC verify are append-only fact creations and follow the
  `fid issue create` precedent (creations do not emit a row in the derived
  `activity_log` timeline table; only `update`/`transition` mutations do, per
  Phase 5 CF-005). The journal still records every AC write — rebuild is
  unaffected (VAL-009). Consistent timeline population across all write commands
  is deferred to v0.1e via **CF-011** (with `fid activity list`, PRD-OUT-017).
  No activity assertion in this contract by design.
- **Validate-before-append discipline (CF-007).** Rejections (AC-on-epic, unknown
  role, invalid `--status`, AC/issue mismatch) must occur **before** any journal
  append: non-zero exit, clear error, and zero new journal lines / projected rows.
- **`ac_verifications.completion_id` stays NULL in Phase 7.** Completion linkage
  is Phase 8 (`fid issue complete`); `fid issue ac verify` has no completion flag.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 7, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | PRD-AC-001; design §3 `fid issue ac add` | `fid issue ac add <task-or-subtask> --text --ordinal --by` creates an `AC-<hash>` row in `acceptance_criteria` with the given text and ordinal; read-back orders ACs by ordinal (ties by id). | CLI: add ≥2 ACs with distinct ordinals to a task, read back via `fid issue ac list` / `fid issue get --section acceptance`; assert text, `AC-` ID shape, and ordinal ordering. | QA Lead | task-1 | ACs persist with correct text + ordinal and read back in ordinal order with `AC-` IDs. | AC not persisted, ordinal lost/ignored, wrong ID shape, or wrong ordering. | Draft | Ordinal uniqueness is not enforced (no schema constraint); ties break by id. |
| VAL-002 | PRD-AC-002; design §3; CF-007 | `fid issue ac add` on an **epic** is rejected before any journal append (non-zero exit, clear error, no new event/row); the same command on a **task/subtask** succeeds. | CLI over two targets: (a) epic → non-zero exit + error naming the kind constraint + zero new journal lines; (b) task → success + one new AC. | QA Lead | task-1 | Epic case rejected with clear error and zero append; task/subtask case succeeds. | Epic case appends a journal line/row, succeeds, or fails with an opaque error; or task case is rejected. | Draft | Epic kind-guard trigger is a backstop; the CLI must reject first (no journal poisoning). |
| VAL-003 | PRD-AC-003; design §3 `fid issue ac verify` | Two different verifier roles can each record an independent verification on the same AC, and both latest statuses are retained independently (one role's outcome never overwrites another's). | CLI: verify one AC as `developer` (pass) and `qa_lead` (fail); assert latest-per-role shows developer=pass AND qa_lead=fail. | QA Lead | task-2 | Both roles' latest verifications are retained independently and correct. | A second role's verification overwrites/clears the first role's, or only one role is retained. | Draft | Permissions matrix (which roles may verify) is v0.1d; Phase 7 only requires `RequireKnownRole`. |
| VAL-004 | PRD-AC-005; design §3 | Verifications are append-only: a `fail` from a role followed by a `pass` from the **same** role yields two retained `ac_verifications` rows, and latest-per-role reflects the most recent (`pass`). | CLI/unit: verify same AC+role fail then pass; assert two rows persist and latest-per-role = pass; after only the fail, latest = fail. | QA Lead | task-2 | After fail→pass: two rows retained, latest=pass; after fail only: latest=fail. | A re-verification overwrites the prior row (only one row), or latest does not reflect the most recent. | Draft | History pruning/compaction is out of scope; all rows are preserved. |
| VAL-005 | PRD-AC-004; design §3 | `fid issue ac verify` records verifier_role, status, evidence_type, evidence_ref, notes, and verified_at; an invalid `--status` (not pass/fail/blocked/not_applicable) is rejected before append. | CLI: (a) valid verify with evidence flags → row stores all fields; (b) `--status bogus` → non-zero exit + error + zero append. | QA Lead | task-2 | Valid verify persists all named fields; invalid status rejected with zero append. | A named field is dropped/mis-stored, or an invalid status is accepted/appended. | Draft | Evidence-ref resolution/integrity (`fid evidence check`) is v0.4; Phase 7 stores values verbatim. |
| VAL-006 | PRD-AC-006; design §3 ac-list JSON | `fid issue ac list <issue-id>` returns each AC with `latest_verifications` keyed by verifier role, showing the most recent status per role (and not raw history). | CLI: after multi-role + repeated verifications, `fid issue ac list --format json` shows one latest entry per verifying role matching the design §3 shape. | QA Lead | task-3 | List shows latest-per-role (one entry per verifying role) in the documented shape; superseded statuses are not shown as separate entries. | List shows stale/duplicate per-role entries, omits a verifying role, or diverges from the documented shape. | Draft | `ac_verifications.completion_id` linkage and full history listing are out of scope. |
| VAL-007 | PRD-OUT-001–004; design §3 | `fid issue ac list --format json\|toon\|markdown\|yaml` renders without error; JSON is the default when `--format` is omitted. | CLI: run all four formats against a populated issue; assert each renders successfully and omitting `--format` equals JSON. | QA Lead | task-3 | All four formats render; default (no flag) is JSON. | A format errors/malforms, or the default is not JSON. | Draft | The format engine itself is Phase 3; this proves `ac list` dispatches through it. |
| VAL-008 | PRD-WF-007; roadmap §v0.1c; Phase 5 fixture replacement | A real `all_acs_verified`-gated transition is **rejected** when required-role pass verifications are missing (error names AC IDs / missing roles) and **succeeds** once `fid issue ac verify` records the required pass(es) — driven by real write commands, no journal fixtures. | Integration/CLI: create task → add ACs → attempt gated transition (reject, naming unmet ACs/roles, zero state change) → `ac verify` required role(s) pass → re-attempt (success). | QA Lead | task-4 | Reject path names unmet AC IDs/roles with no state change; success path transitions after real pass verifications. | Reject path omits IDs/roles or mutates state, or the transition fails after required passes are recorded. | Draft | Other transition rules (`has_completion_report`, `concerns_addressed`) are exercised in Phases 8/9. |
| VAL-009 | ADR-0001; roadmap §v0.1c exit | Deleting the SQLite DB and running `fid rebuild` reproduces identical `acceptance_criteria` and `ac_verifications` projected state from the journal (AC adds + verifications are journaled and deterministically replayed). | Integration test/scripted CLI: snapshot AC + verification state, delete DB, `fid rebuild`, assert identical projected state. | QA Lead | task-4 | Post-rebuild AC + verification state is identical to pre-delete state. | Any divergence in AC rows, verification rows, or latest-per-role after rebuild. | Draft | Held-out scenario harness location remains a separate scheduled commitment. |
| VAL-010 | Phase brief; roadmap §v0.1c | Human dogfood confirms the AC lifecycle (add ordered ACs to a task, multi-role verify, list latest-per-role, epic rejection) reads and behaves as expected in a real `.fid/` workspace. | Human checklist against a real `.fid/` workspace; record outcomes per observable. | Human | phase-end | Human confirms ordered add+list, multi-role latest-per-role, and a clear epic rejection all match expectation. | Human finds incorrect ordering/latest-per-role, an unclear epic rejection, or behavior that contradicts expectation. | Draft | QA Lead fans this into one human-checklist row per observable (add+list / multi-role verify / epic rejection). |

## Coverage Rules

Confirmed before task execution begins:

- Phase-scope requirements each map to ≥1 `VAL-NNN`: PRD-AC-001→001, AC-002→002,
  AC-003→003, AC-004→005, AC-005→004, AC-006→006; PRD-OUT-001–004 (formats on
  `ac list`)→007; PRD-WF-007 real-data gating→008; journal/rebuild canonical
  (ADR-0001)→009; human AC-lifecycle feel→010.
- PRD-DEV-001–010 (completions, concerns) have **no Phase 7 assertion by design**
  — they are Phases 8–9; the phase brief records the boundary.
- Every `VAL-NNN` is referenced by a task spec or marked `phase-end` (VAL-010).
- Binary/conditional assertions name both states — VAL-002, 003, 004, 005, 008.
- DEC-020 (no `activity_log` for AC writes; matches `create`) is confirmed by the
  human; consistent timeline population is deferred to v0.1e via CF-011. No
  activity assertion in this contract by design.
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase-review, or
  `(human)` verification and fans VAL-010 into per-observable human rows.

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-14 | VAL-001..010 | created | Phase 7 planning draft | phase brief; roadmap §v0.1c; PRD-AC-001..006; PRD-WF-007 |
</content>
</invoke>
