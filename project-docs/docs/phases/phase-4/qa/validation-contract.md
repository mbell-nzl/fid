# Phase 4 Validation Contract - Issue Authoring And Discovery

> Phase: 4 - v0.1b Issue authoring and discovery
> Author: Tech Lead (draft)
> Last refined: pending QA Lead refinement
>
> Lives at `project-docs/docs/phases/phase-4/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 4 is complete. The validation contract specifies **what must be proven**;
the QA plan is the strategy for **how to prove it**.

Phase 4 is the first v0.1b slice. It makes issue authoring and discovery usable
end to end: create epics/tasks/subtasks with scheduling and owner metadata,
enforce the scheduling backstop carried from Phase 2, centralise the enriched
issue read model carried from Phase 3, add structured `fid issue get --section`,
and add `fid issue list` filters/order with the full read-format surface.

If a required validation assertion cannot be satisfied, stop and escalate
instead of weakening the assertion or silently adapting the spec.

## Phase 4 Planning Assumptions

- `fid issue update`, `fid issue transition`, `fid ready`, activity-log
  population, and workflow transition-rule enforcement are explicitly deferred
  to later v0.1b phases.
- `fid ready --explain` remains v0.2 (PRD-ISSUE-024), and task reparenting /
  inbox promotion remains v0.2 (PRD-ISSUE-025).
- Label CRUD remains out of scope. Phase 4 `--label` assignment and label
  filtering operate only against existing `labels` rows; no task may silently
  add auto-create semantics for labels.
- Phase 4 epic metadata covers scheduling only. `epic_metadata.brief_md_path`,
  `target_start`, and `target_end` are deferred until a CLI flag/update surface
  is recorded.
- If a task is inserted under an epic that lacks an `epic_metadata` row, Phase 4
  treats that as invalid. The sequence rule cannot be evaluated without the
  epic scheduling mode.
- For Phase 4, `fid issue get --section` means the CLI's named structured
  sections (`acceptance`, `concerns`, `completions`, `links`, `claims`). Full
  markdown-body heading extraction from `spec_md_path` is not implemented unless
  the human requests a spec update before execution.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 4 and sequential from `VAL-001`.
- Carry-forward assertions are re-issued in the receiving phase with a new
  `VAL-NNN` and a source note that points to the earlier assertion.
- Retired or superseded assertions stay in the table with status and rationale;
  do not reuse IDs.

## Coverage-Evidence Carry-In

Where a `VAL-NNN` behavior lives in a COV-001 domain package
(`internal/issues`, `internal/journal`, `internal/policy`; thresholds in
`08_QUALITY_PLAN.md` section 4.3), its evidence method names package-level
tests in that package. Hold COV-001 >= 80% and COV-002 (`cmd/fid`) >= 70%.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | CF-004; PRD-ISSUE-013, PRD-ISSUE-014; `03_TRD_SCHEMA.md` section 14.1 | The database backstop enforces task sequence semantics against `epic_metadata`: tasks under sequential epics require positive sequence, tasks under prioritised epics or no epic must not have sequence, tasks under epics missing metadata are invalid, and duplicate sibling sequence values are rejected | Migration/trigger test in `internal/db` for raw SQL reject/accept states plus journal projection test in `internal/journal` proving the trigger is active during event projection; `cmd/fid` create test proves invalid scheduling/sequence states render actionable CLI errors | QA Lead | tasks 1, 3 | Sequential epic task with sequence succeeds; sequential epic task without sequence fails; prioritised epic task with sequence fails; unattached task with sequence fails; task under metadata-less epic fails; duplicate sibling sequence fails; valid journal-projected rows replay cleanly; CLI errors name the violated constraint | Any invalid state is accepted, any valid state is rejected, the rule is only tested through direct SQL without projection evidence, or CLI create hides the constraint behind an opaque error | Draft | `fid ready` scheduling behavior is Phase 6; this VAL covers the persistence/create backstop only |
| VAL-002 | PRD-ISSUE-001, PRD-ISSUE-002, PRD-ISSUE-003, PRD-ISSUE-004, PRD-WF-005, PRD-WF-006; `02_TRD.md` section 4.6 | `fid issue create` creates epic, task, and subtask issues with generated IDs, legal hierarchy, default workflow IDs, initial `todo` status, and correct status categories when read back | Package-level `internal/issues` create/read tests plus `cmd/fid` CLI smoke for all three kinds; invalid hierarchy table tests cover both qualifying and non-qualifying states | QA Lead | tasks 2-3 | Epic without parent, task with optional epic parent, and subtask with task parent all create and read back with `*-default` workflow, `todo`, and `unstarted`; invalid parent/kind combinations exit non-zero with an error naming the constraint | Unsupported kind remains the only create path, invalid hierarchy is accepted, valid hierarchy is rejected, or read-back defaults are wrong | Draft | Transition edges are Phase 5 |
| VAL-003 | PRD-ISSUE-015, PRD-ISSUE-016, PRD-ISSUE-019, PRD-ISSUE-020; `03_TRD_SCHEMA.md` section 2.1 | `fid issue create` persists supported authoring fields: priority, sequence, spec markdown path, estimate tokens, due date, labels from existing rows, inbox flag, owner role default/override, and owner actor | Package-level `internal/issues` tests over populated and omitted fields plus CLI JSON read-back after create | QA Lead | tasks 2-3 | Populated fields survive journal append, projection, read-back, and delete-DB rebuild; omitted fields render as null/false/empty as specified; owner defaults are `tl` for epics and `developer` for tasks/subtasks unless overridden; unknown labels are rejected without auto-create | Populated fields are dropped or mis-typed, omitted fields produce unstable shapes, owner defaults are wrong, labels are auto-created silently, or rebuild changes the read shape | Draft | Cycle assignment is v0.5; task reparenting / inbox promotion is v0.2; label CRUD ships later |
| VAL-004 | Roadmap v0.1b Epic metadata; PRD-ISSUE-001; `03_TRD_SCHEMA.md` section 2.1 | Creating an epic creates one matching `epic_metadata` row with scheduling (`sequential` or `prioritised`); non-epics cannot receive epic metadata | `internal/issues` tests and CLI smoke for sequential/prioritised epics; raw SQL or journal projection rejection for task/subtask metadata | QA Lead | tasks 2-3 | Sequential and prioritised epics each create exactly one metadata row and read back scheduling-related fields where exposed; task/subtask metadata insert is rejected | Epic metadata is missing/duplicated, invalid scheduling is accepted, or non-epic metadata is accepted | Draft | Ready-queue interpretation of scheduling is Phase 6; brief path and target dates need a recorded CLI/update surface |
| VAL-005 | PRD-OPS-001, PRD-OPS-003, PRD-OPS-004; TRD-ARCH-001, TRD-ARCH-002 | Each issue create command is one logical journal event whose projection is deterministic and rebuildable across epic/task/subtask plus metadata and label-assignment states | `internal/journal` or `cmd/fid` rebuild parity test: create representative issues, remove DB, run `fid rebuild`, compare JSON/list outputs and journal line counts | QA Lead | tasks 2-3 | Each create appends exactly one event; delete-DB + rebuild reconstructs the same issues, metadata, and assignments; malformed/conflict behavior from prior phases remains unchanged | Create writes multiple events for one issue, projection diverges after rebuild, or the read path serves stale data after journal advance | Draft | Performance benchmark refresh is not required unless the implementation changes replay complexity materially |
| VAL-006 | CF-006; PRD-ISSUE-006; TRD-ARCH-005; `05_DESIGN_PLAN.md` section 6.1 | The enriched issue detail read model is centralised in `internal/issues` and includes labels, acceptance criteria with latest verification per role, open concern count, latest triage status, and links without SQL assembly in `cmd/fid` | Package-level `internal/issues` tests over populated and empty states; `cmd/fid` get tests for JSON/markdown parity; import/source scan showing `cmd/fid` routes but does not assemble SQL-backed detail sections | QA Lead | task-4 | Empty issue returns stable empty arrays/nulls; populated issue returns labels, AC latest verifications, concern/triage fields, and links; `cmd/fid` contains no direct read-model SQL assembly for those sections | Detail fields are missing, empty states are unstable, latest verification selection is wrong, or read-model assembly remains in `cmd/fid` | Draft | Claims and activity timeline details are not enriched until their write/read commands ship |
| VAL-007 | PRD-ISSUE-018; `06_DESIGN_PLAN_CLI.md` `fid issue get --section` | `fid issue get <id> --section <name>` returns only the requested structured section for supported names and rejects unknown sections with an actionable error | `cmd/fid` table tests for `acceptance`, `concerns`, `completions`, `links`, `claims`, and an unknown section; seeded empty and populated section states | QA Lead | task-4 | Each supported section returns the matching section data only, preserving format behavior where applicable; an unknown section exits non-zero and names the supported sections | Supported section emits the whole issue unintentionally, unknown section silently falls back, or empty section states error unexpectedly | Draft | Full markdown-body section extraction from `spec_md_path` is a planning assumption, not Phase 4 scope |
| VAL-008 | PRD-ISSUE-007; `05_DESIGN_PLAN.md` section 6.2; `06_DESIGN_PLAN_CLI.md` `fid issue list` | `fid issue list` filters issues by kind, exact status, status category, epic, existing label, priority, inbox state, and open-concern flag | Package-level `internal/issues` list tests plus `cmd/fid` CLI tests over seeded mixed issue set; each filter has include and exclude evidence | QA Lead | task-5 | Every filter returns the expected subset and excludes at least one non-matching issue; empty result returns an empty list at exit 0 | Any filter is missing, includes non-matching rows, excludes matching rows, or empty results fail | Draft | Label CRUD and task reparenting are out of scope; label filter tests may seed label rows |
| VAL-009 | PRD-ISSUE-007, PRD-OUT-004; `05_DESIGN_PLAN.md` section 6.2; `06_DESIGN_PLAN_CLI.md` `fid issue list` | `fid issue list` supports deterministic ordering (`sequence`, `priority`, `created`, `updated`) and dispatches JSON, TOON, markdown, and YAML through the format engine; unknown formats are rejected | `internal/issues` ordering tests plus `cmd/fid` CLI tests for all formats and unknown format; markdown list renderer has a golden or structural assertion | QA Lead | task-5 | Ordering is stable with tie-breaks; JSON/TOON/YAML carry equivalent list fields; markdown renders a readable list; `--format xml` exits non-zero naming known formats | Ordering is nondeterministic, a supported format fails, fields drift across structured formats, or unknown format silently falls back | Draft | `fid ready` format dispatch is Phase 6 |
| VAL-010 | `08_QUALITY_PLAN.md` section 5.3; v0.1 dogfood goal in roadmap | A human can dogfood the Phase 4 authoring/discovery loop and judge the command wording/output useful enough for the next v0.1b phase | Human checklist rows for create epic/task/subtask, get, section, list, and one invalid hierarchy error | Human | phase-end | Human checklist has no FAIL outcomes; any PARTIAL outcome has an accepted follow-up or carry-forward target | Human records FAIL without accepted disposition, or checklist cannot be executed from the built binary | Draft | This is product judgment only; deterministic mechanics are VAL-001..VAL-009 |

## Coverage Rules Check (Tech Lead draft)

- In-scope Phase 4 requirements are covered: PRD-ISSUE-001..007,
  PRD-ISSUE-013..016, PRD-ISSUE-018..020, PRD-WF-005..006, PRD-OUT-004,
  CF-004, and CF-006.
- Explicitly deferred requirements are out of Phase 4 scope: PRD-ISSUE-008..012
  (`fid ready`, Phase 6), PRD-ISSUE-021..023 (`issue update` / label mutation,
  Phase 5 or later), PRD-WF-001..004/007..010/013 (`issue transition`, Phase 5),
  PRD-ISSUE-024..025 (v0.2), and CF-005 (Phase 5).
- Every `VAL-NNN` is owned by a task or phase-end: YES.
- Conditional/binary assertions name both states: YES, especially VAL-001,
  VAL-002, VAL-004, VAL-007, VAL-008, and VAL-009.
- Parity/reference assertions carry deterministic proof: VAL-005 (rebuild
  parity), VAL-006 (centralised read shape and adapter-boundary scan), VAL-009
  (format parity / markdown list check).
- COV-001 carry-in applied: `internal/issues` and `internal/journal` behaviors
  require package-level tests; `cmd/fid` coverage remains COV-002.

## Outcomes

Phase-end final review records every `VAL-NNN` as `PASS`, `FAIL`, `BLOCKED`, or
`CARRY-FORWARD`. `PASS` without cited evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-07 | VAL-001..010 | created (Tech Lead draft) | Phase 4 planning | N/A |
