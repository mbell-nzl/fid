# Phase 15 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes, and
> wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-037: `fid dev *` composites are one-event multi-effect writes in a new `internal/dev` boundary

- **Date:** 2026-06-29
- **Context:** Phase 15 ships the first **composite** write commands
  (`fid dev start/gate/finish/block`). Each orchestrates several existing
  single-effect writes (claim, transition, completion, `agent_runs`,
  `gate_runs`). Planning interpretations PI-1..PI-6 (confirmed by the human
  2026-06-29) fixed the boundary, the one-event rule, and the per-command
  effect semantics.
- **Decision:**
  - **PI-1:** the four composites live in a new inward-pointing domain package
    `internal/dev` that orchestrates `internal/claims`, `internal/issues`,
    `internal/runs`, and `internal/workflow` effect builders. No new
    persistence contract; recorded in `02_TRD.md` sections 2-3.
  - **PI-2:** each composite emits exactly one journal event
    (`dev.started` / `dev.gate.recorded` / `dev.finished` / `dev.blocked`)
    carrying all effects, validate-before-append (CF-007): any rejection path
    exits non-zero with an actionable error and appends zero journal lines. The
    effect-driven projector needs no upgrader for the new event types.
  - **PI-3:** `dev start` acquires the caller claim and transitions
    `todo → in_progress` in the same event; the claim acquired in-event
    satisfies PRD-CLAIM-005, so no pre-existing claim is required. A live claim
    held by another role is rejected (no `--force` in v0.1 `dev start`).
  - **PI-4:** `dev finish` satisfies `has_completion_report` by the completion
    created in the same event and rejects before append only on unmet
    `all_acs_verified[developer]`.
  - **PI-5:** `dev finish` (→ `dev_complete`) and `dev block` (→ `failed_gate`)
    both release the caller's live claim as part of the composite, matching the
    `06_DESIGN_PLAN_CLI.md` "claim released" output, even though those target
    states stay started-category; plain `fid issue transition` semantics are
    unchanged.
  - **PI-6:** `dev finish` closes the issue's open `agent_runs` row when one
    exists; `dev block` leaves runs alone.
  - CF-013 supplied the shared completion-effect builder (`BuildCompletionEffects`)
    that `dev finish`/`dev block` consume.
- **Consequences:** No schema migration and no `project_effects.go` change; the
  generic effect projector applies each composite's effects. CF-011 broad
  `activity_log` reconciliation stays Phase 16 — the composites only project the
  per-command `activity_log` rows their transition/claim writes already emit.
- **Refs:** PI-1..PI-6; phase-15/task-2; phase-15/task-3; VAL-003, VAL-004,
  VAL-007, VAL-008, VAL-009, VAL-010, VAL-011; PRD-DEV-008..010; PRD-CLAIM-005;
  PRD-WF-007..010; TRD-ARCH-010; `02_TRD.md` sections 2-3; `02_TRD.md` §4.7.

### DEC-038: `fid dev gate` records a reported gate result; it does not execute the command (DN-A)

- **Date:** 2026-06-29
- **Context:** `06_DESIGN_PLAN_CLI.md` showed `fid dev gate --cmd <command> --by`
  with confirmation output (`exit_code`, cumulative `gates_passed`/`gates_failed`)
  but no source for the exit code. DN-A (confirmed by the human 2026-06-29)
  resolved whether fid runs the command.
- **Decision:** fid **records** a reported gate result and does **not** run the
  command (GLOSSARY: "Developers run gates during a task attempt; standalone
  `fid dev gate` writes append-only `gate_runs` rows"; ADR-0001 deterministic
  recorder; avoids a command-execution security surface). The `fid dev gate`
  flag table gains `--exit-code <int>` (required) and `--output-excerpt <string>`
  (optional). One `dev.gate.recorded` event appends one `gate_runs` row on the
  derived active attempt (`MAX(task_completions.attempt_number)+1`, 0 when none),
  recording the reported `command`/`exit_code`/`output_excerpt`/`ran_at`. Reject
  before append when the issue is not in a started-category state or the caller
  does not hold the live claim. `gates_passed`/`gates_failed` stay derived at
  finish/block, not stored on `gate_runs`.
- **Consequences:** Public-interface addition recorded in
  `06_DESIGN_PLAN_CLI.md` (Task 2) — `--exit-code`, `--output-excerpt`, and a
  "records, does not execute" clarification. No threat-model obligation added
  (no command execution). No schema change (`gate_runs` columns sufficient).
- **Refs:** DN-A; phase-15/task-2; VAL-005; PRD-DEV-004; PRD-DEV-010;
  `02_TRD.md` §4.7; `03_TRD_SCHEMA.md` §7; `06_DESIGN_PLAN_CLI.md` `fid dev gate`;
  GLOSSARY; ADR-0001.

### DEC-039: `fid dev start --model <m>` opens the base run; dev-opened run leaves cost columns NULL (DN-B)

- **Date:** 2026-06-29
- **Context:** `agent_runs.model` is `NOT NULL`, but the CLI spec gave
  `dev start` no model flag. DN-B (confirmed by the human 2026-06-29) resolved
  how the opened run records its model.
- **Decision:** `fid dev start` gains `--model <string>` (required) so the
  opened `agent_runs` row records its model. `tokens_in`/`tokens_out`, snapshot
  prices, and `cost_usd` stay NULL on the dev-opened run in v0.1; the dedicated
  `fid run record` (Phase 14) remains the full-fidelity cost-capture path
  (per-composite cost capture and run lifecycle are v0.2c). `fid dev finish
  --tokens` populates `task_completions.actual_tokens`, not run tokens.
- **Consequences:** Public-interface addition recorded in
  `06_DESIGN_PLAN_CLI.md` (Task 2). No schema change.
- **Refs:** DN-B; phase-15/task-2; VAL-003; PRD-DEV-008; PRD-CLAIM-005;
  `03_TRD_SCHEMA.md` §2, §7; `06_DESIGN_PLAN_CLI.md` `fid dev start`.

## Carry-Forward Review At Wrap (Tech Lead, 2026-06-29)

| ID | Disposition | Rationale |
|---|---|---|
| CF-008 | Re-defer to v0.2 | Resolution-reason field is a product/data-model enhancement; untouched by the composites. |
| CF-011 | Re-defer broad reconciliation to Phase 16 | The composites project only the per-command `activity_log` rows their transition/claim writes already emit; they do not widen or backfill CF-011. Broad timeline policy lands with `fid activity list` (Phase 16). Not worsened. |
| CF-013 | **CLOSED this phase** | Shared completion-effect builder extracted (`issues.BuildCompletionEffects`); `fid issue complete` behavior unchanged, `internal/dev` reuses it (VAL-010 PASS). |
| CF-014 | **CLOSED this phase** | `cmd/fid/issue.go` 984→22 (largest split sibling 228) and `internal/journal/replay.go` 485→205; all ≤400 lines, behavior-preserving (VAL-001 PASS). |
| CF-015 | **CLOSED this phase** | `internal/labels/add.go` and `internal/issues/label.go` add/remove write flows extracted to ≤60-line functions, behavior-preserving (VAL-002 PASS). |

Two new carry-forward items opened (below): **CF-016** (first-class issue-scoped
gate read surface, surfaced by the human at wrap) and **CF-017** (`internal/dev`
`Finish`/`Block` function-size cleanup, the recurring function-size debt noted by
TRG-002 and the final review).

## Wrap-Up Dispositions (Tech Lead, 2026-06-29)

### Validation contract outcome audit

All twelve assertions reach a terminal state; none weakened.

- VAL-001..011: PASS in `qa/final-review.md` with cited unit/CLI/integration
  evidence and a green nine-line phase-end gate (focused e2e VAL-011 run +
  full phase command list). No disposition needed.
- VAL-012: was CARRY-FORWARD to the human checklist in final review. The
  generated `qa/human-checklist.md` (HV-15a..HV-15d) is now complete with
  overall **PASS** (completed 2026-06-29, tester MBEL): start claims +
  transitions + opens a run, gate records (not executes), finish completes +
  closes run + releases claim, block records a failed gate with reason — all
  observable from CLI output. Final outcome: **PASS**. No FAIL/BLOCKED/
  CARRY-FORWARD remains; no defer/cancel/escalation required.

### QA recommendation disposition

Final review verdict was PASS WITH IMPROVEMENTS with **no** Blocking or Major
findings and two Improvement notes (both already triaged RESUME). Dispositions:

- **VAL-012 human dogfood** — resolved: the human checklist ran and passed
  overall (HV-15a..HV-15d all `PASS`). No action remains.
- **TRG-001 (Review-1 improvements: completion-report H1 consistency; Task 2
  function-size debt)** — accept. The H1 omission is an artifact-hygiene note,
  not a defect; carry the lesson into Phase 16 task specs (developer completion
  reports include the template H1). Task 2 structure debt folds into CF-017.
- **TRG-002 / final-review unresolved-risk note (`internal/dev` `Finish` 109 /
  `Block` 97 lines above the 60-line target)** — **defer via new CF-017**.
  Behavior-preserving function extraction, gates green, not a blocker. This is
  the same function-size-budget family as CF-015; tracking it as a carry-forward
  keeps it from being lost. Target v0.2c (the next version that reopens
  `internal/dev` for `fid dev resume`/`abandon`), low priority — mirroring how
  CF-013 was scheduled alongside the work that next touched that path.

No do-now QA recommendations were raised.

### Human wrap review note disposition (gate read surface) → CF-016

The human's checklist TL wrap review notes flag a forward observability gap:
`fid dev gate` writes append-only `gate_runs`, but there is **no first-class
issue-scoped read surface** for those recorded gates while a task is still in
progress. `fid issue get --section` supports `acceptance`, `concerns`,
`completions`, `links`, `claims` — not `gates`. After `dev finish`/`dev block`
the gate evidence is visible only indirectly through the completion snapshot
(`task_completions.gate_results`), not as raw append-only gate history.

Disposition: **not a Phase 15 blocker** (out of the composite-write scope;
VAL-005 only requires the `gate_runs` append/reject mechanics, which PASS; the
dogfood checklist HV-15a..d all PASS). Open as **CF-016** with a recommended
shape and target, and surface the scheduling for human confirmation because it
changes a public interface and touches product intent (TL-Lite triggers 1 and 7).

Tech Lead recommendation recorded in CF-016: prefer the
`fid issue get <id> --section gates` shape (fits the existing issue-section read
model and keeps raw append-only gate history distinct from the completion
snapshot) over a new `fid dev gate list` verb (verb-shaped reads cluster with the
v0.2c `fid run show` / `fid stats gates` family). Recommended target **Phase 16**
(the read-surface phase: `fid activity list` + `fid sql` + `fid journal inspect`),
because it is a small additive read over the existing `gate_runs` table with no
schema change and keeps the dogfood-surfaced gap inside v0.1. Note the **interim
mitigation lands in Phase 16 regardless**: `fid sql` (read-only SELECT over
`gate_runs`) and `fid journal inspect <event_id>` (the `dev.gate.recorded` event
effects) both provide a path to inspect recorded gates, so the gap is not
unmitigated. Fallback target if Phase 16 dependency analysis has no room: v0.2c
alongside `fid run show` / `fid stats gates`.

### DEC / deviation spec-impact walk

- **DEC-037** (`internal/dev` boundary + one-event composites, PI-1..PI-6):
  spec-impacting? **Yes, synced in-phase.** `02_TRD.md` sections 2-3 record the
  new `internal/dev` boundary (Task 2). TRD-ARCH-010 matrix satisfied (VAL-006).
  No schema change. — `touched`.
- **DEC-038** (DN-A: `dev gate` records, not executes; `--exit-code`/
  `--output-excerpt`): spec-impacting? **Yes, synced in-phase.**
  `06_DESIGN_PLAN_CLI.md` `fid dev gate` updated by Task 2 with the two flags and
  the "records, does not execute" clarification (VAL-005). — `touched`.
- **DEC-039** (DN-B: `dev start --model`): spec-impacting? **Yes, synced
  in-phase.** `06_DESIGN_PLAN_CLI.md` `fid dev start` updated by Task 2 with the
  required `--model` flag (VAL-003). — `touched`.
- **Developer deviations:**
  - Task 1 split `cmd/fid/issue.go` across new sibling files and
    `internal/journal/replay.go` into `replay.go`/`projection.go`/`inspect.go`;
    behavior-preserving, in scope (VAL-001/002). No spec impact.
  - Task 3 deliberately left `06_DESIGN_PLAN_CLI.md` `fid dev finish`/`fid dev
    block` unchanged because the implemented confirmations already match the
    documented examples (flags incl. `--human-intervention`/`--intervention-note`
    and "claim released" output verified present). Confirmed accurate at wrap —
    `06` finish/block sections are `verified-no-change`.
  - Task 4 updated `07_DESIGN_PLAN_FLOWS.md` to reconcile the developer flow
    example with the implemented composite flags/output (`--model`,
    `--exit-code`, `--output-excerpt`, claim-release). In scope; QA-09 verified.
- **Validation-contract deviations:** none; all twelve VAL IDs satisfied as
  written.

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | N/A | No charter-level scope/constraint change in a composite + cleanup slice. |
| `01_PRD.md` | verified-no-change | PRD-DEV-008..010, PRD-CLAIM-005, PRD-WF-007..010 MET (final-review Requirement Compliance, VAL-003..011 PASS); the composites implement existing requirements with no PRD wording change. |
| `02_TRD.md` | touched: see DEC-037 | Task 2 recorded the new `internal/dev` inward boundary in the architecture map and module-boundary table (PI-1); §4.7 gate-recording semantics satisfied (VAL-005). |
| `03_TRD_SCHEMA.md` | verified-no-change | No migration added; section 7 `agent_runs`/`gate_runs`/`task_completions` columns sufficient; section 13 `dev.start/gate/finish/block` rows (developer + human) match the in-code matrix (VAL-006 PASS); QA confirmed no `migrations`/`project_effects.go` change. |
| `04_TRD_CONFIG.md` | verified-no-change | No pricing/config surface touched; the composites add no config keys. |
| `05_DESIGN_PLAN.md` | verified-no-change | Command-tree section already lists the `fid dev *` subcommands; the composites match the documented tree, no edit needed. |
| `06_DESIGN_PLAN_CLI.md` | touched: see DEC-038, DEC-039 | Task 2 added `fid dev start --model`, `fid dev gate --exit-code`/`--output-excerpt`, and the records-not-executes note. `fid dev finish`/`fid dev block` sections verified accurate and left unchanged (Task 3). |
| `07_DESIGN_PLAN_FLOWS.md` | touched: see deviation walk | Task 4 reconciled the developer execution flow example with the implemented composite flags and claim-release output (QA-09 verified, 12 changed lines). |
| `08_QUALITY_PLAN.md` | verified-no-change | Phase-end gate run as defined; COV thresholds met (`internal/issues` 81.7%, `internal/labels` 82.4%, `internal/journal` 80.5%, `internal/policy` 84.6%, `cmd/fid` 76.9%; `internal/dev` 75.4% is a new write-path package below 80% but covered by the green gate and final-review acceptance). No new gate added. The 60-line function target is not a QA-plan gate; `internal/dev` `Finish`/`Block` debt tracked as CF-017. |
| `09_AI_WORKFLOW.md` | N/A | No agent-workflow change. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 15 delivered exactly the v0.1e Phase 15 entry (the four `fid dev *` composites + CF-013/014/015 cleanup); exit criterion met. No roadmap wording change required. |
| `GLOSSARY.md` | verified-no-change | "Composite", "gate run", and "developer composite" language already canonical; the records-not-executes gate semantics (DEC-038) match the existing GLOSSARY entry, introducing no new term needing definition. |

## Triage-Log Audit

- [x] Walked the Phase 15 triage record (`triage-log.md`): TRG-001, TRG-002.
- [x] Confirmed RESUME verdicts in hindsight:
  - **TRG-001** (Review-1 improvements: missing template H1 in two task reports;
    Task 2 `internal/dev` / `cmd/fid/dev.go` function-size debt) correctly
    RESUMED — Review-1 recorded zero Blocking findings and `VERDICT: CONTINUE`;
    the findings were artifact-hygiene and non-blocking structure notes, none
    requiring escalation.
  - **TRG-002** (Task 3 `internal/dev` `Finish`/`Block` above the 60-line target)
    correctly RESUMED — all required Task 3 gates PASSED, no spec deviation, and
    the CF-013 shared-builder extraction landed. Folded into CF-017 at wrap so the
    debt is tracked, not lost.
- [x] Flagged any TRG that should have escalated: none. Both are within
  TL-Lite/local-resume authority; neither crossed a public-interface, data-model,
  or standing-decision boundary.
- [x] Confirmed autonomous TL-Lite resume remains in effect (DEC-002; no warm-up
  tally to maintain).

## Workflow Notes

- **What worked well?** The kickoff confirmation pass (DN-A, DN-B, PI-1..PI-6 all
  confirmed by the human before execution) let all four composites land with zero
  re-escalation; they mapped cleanly onto DEC-037/038/039. Batching the three
  carry-forward cleanups (CF-013/014/015) into the same phase that introduced the
  composites paid off: Task 1's behavior-preserving split shrank `issue.go`
  (984→22) and `replay.go` (485→205) before any new code landed, and CF-013's
  shared `BuildCompletionEffects` builder was then directly consumed by
  `dev finish`/`dev block` (VAL-010), so the cleanup and the feature reinforced
  each other rather than competing. Sequencing `dev start`/`gate` (Task 2) before
  `dev finish`/`block` (Task 3) meant the closing composites tested against real
  opened runs and recorded gates. Carrying the Phase 14 conditional-command
  wording lesson forward worked: every task that touched `cmd/fid` recorded the
  integration gate PASS, and no ESC-001-style artifact defect recurred.
- **What was inefficient?** Minor. Two Review-1 task reports omitted the template
  H1 (artifact hygiene, no rework). The `internal/dev` `Finish`/`Block`
  function-size debt was noted twice (Review-1 for Task 2, final review for Task
  3) before being consolidated into CF-017 — a single carry-forward pinned at
  first appearance would have been tidier (this is exactly the standing
  "pin recurring improvement notes to a carry-forward ID on first appearance"
  rule, applied a phase late).
- **Were tasks the right size?** Yes. Four tasks, each one coherent session:
  cleanup (1), the two opening composites (2), the two closing composites + CF-013
  (3), and the e2e smoke + dogfood prep (4). Clean 1→2→3→4 dependencies; no
  mid-phase do-now task needed.
- **Did the QA review cadence feel right?** Yes. One mid-phase review over the
  behavior-bearing opening tasks (1-2) plus final-review validation of the closing
  composites (3) and the e2e smoke (4) put per-task signal where the risk was. The
  composites are the first multi-effect writes, and the one-event/zero-append
  assertions (VAL-003/005/007/009/011) were all checked at the right window.
- **Were there token-efficiency opportunities missed?** Minor. Completion reports
  cited exact VAL IDs and test names, keeping final QA cheap. The only overhead was
  the twice-noted function-size debt and the H1 omissions; both are addressed
  (CF-017 opened; H1 lesson carried into Phase 16 task specs).

### Deep retrospective recommendation

No deep retrospective required. Phase 15 closed PASS with no Blocking/Major
findings, all twelve VAL outcomes terminal, and the v0.1e Phase 15 exit criterion
met (an agent can run `fid dev start → gate → finish` to drive a task from claim
through `dev_complete` in composite steps, `fid dev block` to record a failed gate
with reason, each composite one auditable event) against a real `.fid/` workspace.
Three scheduled carry-forwards (CF-013/014/015) closed; two new ones opened
(CF-016 gate read surface, CF-017 `internal/dev` function-size cleanup). Carry
into Phase 16 planning: (1) the gate read-surface decision (CF-016) — confirm the
`fid issue get --section gates` shape and Phase 16 vs v0.2c target with the human,
and note `fid sql`/`fid journal inspect` are the interim path; (2) pin recurring
improvement notes (function-size debt) to a carry-forward ID on first appearance,
not at wrap; (3) keep the template H1 in developer completion reports.
