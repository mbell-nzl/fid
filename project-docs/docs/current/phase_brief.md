# Phase Brief - No Active Phase

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session.
>
> **Size target:** under 150 lines.

## Status

**No phase is active.** Phase 19 is **COMPLETE / WRAPPED** (2026-08-08,
DEC-106). Do not start implementation work against this file. The next phase
begins with a Tech Lead planning session that replaces this brief entirely.

## Last Completed Phase

**Phase 19 — `fid workflow define/assign/list` + escalation default workflow**
(v0.2a slice 2 of 3), heavy mode, six tasks, 70-assertion validation contract.

- **Outcome:** 70 PASS / 0 FAIL / 0 BLOCKED / 0 CARRY-FORWARD. QA Lead
  subsequent final review **PASS WITH IMPROVEMENTS**; Human checklist **9/9
  PASS**. **Zero product defects.**
- **Shipped:** `fid workflow define` / `assign` / `list`, the
  `workflow_assignments` projected table with kind→project→built-in precedence,
  the `escalation-default` seed (PRD-WF-005a), and assignment-aware issue
  creation.
- **Archive:** `../phases/phase-19/` — `progress.md`, `decisions.md`,
  `triage-log.md`, `tl-escalations.md`, `tasks/`, `qa/`.
- Branch: `phase-19-workflow-authoring`.

## Next Planning Checkpoint

**Phase 20 — `role_permissions` migration + `fid role create` /
`fid policy set` + v0.2a closure.** Roadmap entry:
`../../specs/10_IMPLEMENTATION_ROADMAP.md` (v0.2a section).

Inputs the planning session must take:

- **Deep retrospective first (recommended, DEC-104/DEC-106).** Phase 19 spent
  eight escalations and twelve decisions to ship zero product defects. The brief
  is in `../phases/phase-19/decisions.md` → *Deep retrospective recommendation*.
  Its output should shape Phase 20's evidence design before the contract is
  drafted.
- **DEC-079 unblock.** Custom workflow *state names* are not transitionable
  until the permissions matrix is data-backed. Phase 20's `fid policy set` is
  the unblock, and the v0.2a exit smoke depends on it.
- **CF-026a is a scheduling candidate** — `fid workflow show`, a bounded
  additive read surface over already-projected tables, raised by the Phase 19
  human dogfood. **CF-026b is not schedulable** until its lifecycle model is
  recorded as a decision.
- **CF-027 is now safe to take** — CF-025's parity test is the regression guard
  it was sequenced behind. It reopens `internal/issues/create.go`, which Phase
  20 may not otherwise touch; decide explicitly rather than letting the trigger
  lapse.
- **CF-023** carries five accumulated drill cases plus mechanical `Refs:`
  enforcement and has not landed across two phases. Decide whether it is
  scheduled work or a wish list.

## Standing Constraints Carried Out Of Phase 19

Binding on all future phases. Full text in
`../phases/phase-19/decisions.md`; each is also now stated in a canonical spec.

- **Uncaptured gate results are `UNVERIFIED`, never `PASS`** (DEC-102). A
  result attested from memory, transcribed from another agent's report, or
  promised against a named future control is not a passing state.
  **"Pre-existing" is a provenance claim**: cite the introducing commit and show
  it predates phase start. `git show HEAD` proves a line is committed, not
  inherited.
- **A required command is executed by its author before publication**
  (DEC-103). A command in a task spec or QA plan that nobody has ever run is a
  guess. Working-directory sensitivity is declared.
- **Evidence is an automated test** unless the validation contract records why
  it cannot be (DEC-104). Phase 19's headline lesson: four escalations landed in
  a manual ceremony that a 60-line golden test replaced.
- **Cleanup is never evidence**; scratch state is unique per attempt; task specs
  declare required environment capabilities and escalate shortfalls at step zero
  (DEC-101 items 1-3).
- **A task-authored evidence harness that fails while being authored is a
  deviation, not a failed gate** (DEC-100), under its four conditions.
- **No-match inventories terminate `exit 1 (no output)`** (DEC-094).
  `exit 0 (no output)` there is fabricated evidence.
- **The Tech Lead seeds `phases/phase-N/triage-log.md` at planning time**
  (DEC-095), and IDs are phase-scoped `TRG-<phase>-<NNN>` / `ESC-<phase>-<NNN>`
  (DEC-084), written **before** the verdict (DEC-085).
- **Spec-sync obligations land in the owning task's `## Documentation Updates`**,
  not only in the phase sync table (DEC-106 item 5). An obligation recorded only
  in a document its owner does not read is not assigned.
- **`internal/workflow` owns workflow authoring, assignment, and resolution**
  (DEC-089) and may use the journal write API (DEC-105); it must not create
  issues, verify ACs, own migrations, or format output. No in-package
  (`package db` / `package journal`) test file may import `internal/workflow`
  (DEC-096).
- **Custom state names are not transitionable until the Phase 20 matrix work**
  (DEC-079). Do not auto-register matrix rows to work around it.
- **The DEC-066 gate-obligation seam stays behavior-free until v0.2b**, and no
  production Go path reads the frozen v0.1 `definitions` map.
- **Do not weaken TRD-PERF-003 (`<30 s`) or TRD-PERF-001 (`<500 ms`).**
  Current: 100k rebuild `24.749 s`, 10k registry read `44.29 ms` (CF-024).
