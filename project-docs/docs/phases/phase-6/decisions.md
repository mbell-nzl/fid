# Phase 6 Decisions

> Per-phase decision log. Closed phase decision history lives here; live
> carry-forward constraints live in `../../current/decisions_active.md`.
>
> DEC-017 (readiness semantics) and DEC-018 (empty-state contract) were
> human-confirmed at planning (2026-06-14) and are formally recorded at phase
> wrap per the phase brief.

## DEC-017 — Readiness semantics from workflow status category (v0.1b)

- **Status:** Human-confirmed at planning (2026-06-14); formally recorded at wrap.
- **Context:** v0.1b ships **no dependency (`blocked_by`) links and no claims**
  (both v0.1d). Readiness therefore had to be defined from workflow status
  category alone for Phase 6.
- **Decision:** Actionable / ready = status category `unstarted` (`todo`).
  Unfinished = not `completed` and not `cancelled`. A sequential epic's
  lowest-`sequence` unfinished task sitting in a `started` category (e.g.
  `in_progress`) is **not actionable**; the epic then yields no ready task and
  the empty result explains that the current task is already in flight.
- **Scope:** Defines VAL-002/003/004/005/006 semantics. Dependency-link and claim
  blocking remain out of scope until v0.1d.
- **Outcome:** Realized by Task 3 engine + Task 4 CLI; VAL-002..006 PASS; human
  dogfood HV-012a/b confirmed expected sequential / prioritised selection.
- **References:** validation-contract VAL-002..006; PRD-ISSUE-009/010/011/012;
  GLOSSARY "Sequential epic" / "Prioritised epic".

## DEC-018 — Empty-state contract for `fid ready` (frozen)

- **Status:** Human-confirmed at planning (2026-06-14); formally recorded at wrap.
- **Context:** The Workflow Runner depends on a stable `fid ready` answer when
  nothing is actionable. "Nothing ready" is a valid answer, not an error.
- **Decision:** A bare `fid ready` finding nothing ready exits **0**, emits
  `{"task": null}` on stdout (**frozen shape**), and writes a `note:`
  explanation to stderr. A non-empty result returns the task object with no
  `note`.
- **Scope:** The empty stdout shape is frozen; changing it requires a new
  recorded decision. Defines VAL-007.
- **Outcome:** Realized by Task 4; VAL-007 PASS; human dogfood HV-012c confirmed
  the empty `{"task":null}` + clear stderr `note:` in a real workspace.
- **References:** validation-contract VAL-007; PRD-ISSUE-008/010;
  `06_DESIGN_PLAN_CLI.md` §`fid ready`.

## DEC-019 — `fid ready` deterministic selection order (resolves ESC-004 / TRG-001)

- **Status:** Confirmed by human 2026-06-14 (option A).
- **Context:** Task 3 implemented the project-wide / prioritised-epic tie-break
  as `priority ASC` (nil last) → `created_at ASC` → `id ASC`. The PRD pins the
  behavior (PRD-ISSUE-008/012 deterministic; PRD-ISSUE-011 highest-priority) but
  did not spell out the full ordering tuple or nil-priority handling. TRG-001
  BLOCKed and raised ESC-004 because Task 4 makes this a user-visible scheduling
  contract.
- **Decision:** The bare `fid ready` selection order is, in precedence:
  1. `priority` ascending (PRD-ISSUE-015: `0` critical → lowest wins); a task
     with **no priority ranks last** (human-confirmed option A).
  2. `created_at` ascending.
  3. `id` ascending — final deterministic tie-break only, no product meaning.
  Sequential epics continue to order by `sequence` (PRD-ISSUE-009) and
  contribute only their current actionable task project-wide (PRD-ISSUE-012).
- **Rationale:** This matches the selection factors the design plan already
  documented inside the v0.2 `--explain` example (`priority`, `created_at`);
  `id` only guarantees a single deterministic winner. Nil-last treats an
  unprioritised task as least urgent ("if you didn't prioritise it, it waits").
- **Scope of change:** Documentation only. Task 3's implementation stands — no
  code change. Task 4 may proceed and wire the CLI surface against this order.
- **Spec sync:** `06_DESIGN_PLAN_CLI.md` §`fid ready` now states the bare
  selection order explicitly (lifted out of the `--explain`-only block).
  `05_DESIGN_PLAN.md` §6.3 documents the output *shape* only (not ordering):
  verified-no-change. No other spec touched.
- **References:** ESC-004, TRG-001, `tasks/task-3-done.md`, VAL-006; PRD-ISSUE-008,
  011, 012, 015.

## Validation Contract Outcome Audit

Final review records VAL-001..011 PASS with reproducible QA evidence; the human
checklist (HV-012a/b/c, 3/3 PASS, 2026-06-14) resolves VAL-012. All twelve
assertions now stand PASS — no FAIL / BLOCKED / CARRY-FORWARD remains.

| VAL | Outcome | Disposition |
|---|---|---|
| VAL-001..008 | PASS | `fid ready` engine + CLI (Tasks 3/4); QA focused reruns. |
| VAL-009 | PASS | CF-009 update one-field-per-line (Task 1); **CF-009 closed**. |
| VAL-010 | PASS | CF-010 file split 262/260/137 lines + rebuild parity (Task 2); **CF-010 closed**. |
| VAL-011 | PASS | v0.1b cross-command smoke + rebuild parity (Task 5). |
| VAL-012 | PASS | Human dogfood HV-012a/b/c all PASS in a real `.fid/` workspace. |

No assertion was weakened in-flight. DEC-017/DEC-018 (the only planning
decisions that could have forced a VAL re-draft, per the contract's Coverage
Rules) were human-confirmed as written; VAL-003 / VAL-007 stood unchanged.

## Spec Sync Verdict

Recorded for every spec at wrap. `verified-no-change` is a positive review
assertion, not a default skip.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | Phase 6 delivered the in-scope `fid ready` + v0.1b-closure slice; no identity / scope / constraint change. |
| `01_PRD.md` | verified-no-change | PRD-ISSUE-008..012 and PRD-OUT-001..004 MET as written (final review Requirement Compliance). `--explain` (PRD-ISSUE-024) and reparent/promote (PRD-ISSUE-025) confirmed out of scope, not reworded. |
| `02_TRD.md` | verified-no-change | `fid ready` is read-only scheduling over the existing read model; no architecture or dependency change. CF-010 split was behavior-preserving. |
| `03_TRD_SCHEMA.md` | verified-no-change | No schema / view / trigger change; readiness derived from existing status-category projection (DEC-017). |
| `04_TRD_CONFIG.md` | N/A | Phase 6 added no `.fid/` layout or config surface. |
| `05_DESIGN_PLAN.md` | verified-no-change | §6.3 documents the frozen `task` output **shape** only (not ordering); DEC-019 confirmed it needs no change. Mutation/projection flows unchanged. |
| `06_DESIGN_PLAN_CLI.md` | touched: see DEC-019 | §`fid ready` now pins the bare selection order (priority ASC, nil last → `created_at` ASC → `id` ASC) and the DEC-018 empty-state stdout/stderr split; §`fid issue update` updated to one-field-per-line output (CF-009 / DEC-016, Task 1). |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | `fid ready` follows the documented read-model query flow; v0.1b smoke confirms journal→projection→read + rebuild parity end-to-end (VAL-011). |
| `08_QUALITY_PLAN.md` | verified-no-change | COV-001 (`internal/issues` 82.5%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5%) and COV-002 (`cmd/fid` 78.3%) held. §4.4 file-size target now satisfied for the three CF-010 files; no plan reword. |
| `09_AI_WORKFLOW.md` | verified-no-change | Workflow ran as specified (contract → QA plan → dev/QA loop with two mid-phase reviews → final review → human checklist). The one BLOCK (TRG-001/ESC-004) routed and cleared via DEC-019. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 6 delivered the §v0.1b "ready queue + closure" cut and the v0.1b exit criterion; next cut (v0.1c — ACs/Completions/Concerns) unchanged. |

## Triage-Log Audit

| TRG | Verdict | Hindsight |
|---|---|---|
| TRG-001 | BLOCK → RESOLVED (ESC-004 / DEC-019) | **Correct.** The project-wide tie-break is an externally visible scheduling contract not pinned in the PRD; escalating rather than resolving locally was right (TL-Lite triggers 1/4/7). Human confirmed option A; implementation needed no change. No TRG should have escalated that did not, and none over-escalated. |

Autonomous TL-Lite resume is permanently enabled (DEC-002); no warm-up tally to
maintain. Phase 6 had one BLOCK and no autonomous RESUME entries.

## Carry-Forward Disposition (at wrap)

- **CF-009 — CLOSED.** One-field-per-line `fid issue update` output delivered by
  Task 1 (VAL-009 PASS); `06_DESIGN_PLAN_CLI.md` §`fid issue update` updated.
- **CF-010 — CLOSED.** `transition.go` (262), `update.go` (260), `project.go`
  (137) all ≤400 lines after behavior-preserving split (Task 2, VAL-010 PASS).
- **CF-003 — re-deferred** to before v0.1d role-permission / mutable-role work.
- **CF-008 — re-deferred** to v0.2 (resolution type/reason product enhancement).

## QA Recommendation Disposition

- **Improvement (evidence format):** completion reports should paste compact
  focused command output / transcripts for public-output changes. **Defer as a
  standing process nudge** — non-blocking; final QA compensated with focused
  reruns. Folded into the workflow notes below rather than reopened as code work.

## Workflow notes

- **What worked well.** Confirming DEC-017/DEC-018 *and* pre-resolving the
  tie-break ambiguity (DEC-019) at planning meant the user-visible scheduling
  contract was locked before Task 4 wired the CLI; the one mid-phase BLOCK
  (TRG-001) was a clean documentation-only resolution with zero code churn. The
  read-only nature of `fid ready` kept the phase free of schema/dependency risk.
- **What was inefficient.** The recurring improvement finding across review-1 and
  review-2 was the same one (completion reports summarise rather than paste exact
  command transcripts), and final QA had to re-run focused tests to compensate.
  A one-line "paste the exact transcript for any public-output change" reminder
  in the task-spec template would have removed three repeat findings.
- **Were tasks the right size?** Yes. Front-loading the two carry-forward fixups
  (CF-009 Task 1, CF-010 Task 2) before the new engine/CLI/smoke work was the
  right order and each was a clean one-session unit. The engine/CLI split
  (Task 3 pure selection + Task 4 CLI dispatch) made VAL evidence easy to attribute.
- **Did QA cadence feel right?** Yes. Two mid-phase reviews (1-3, 4-5) plus final
  review matched a 5-task phase; no review window was starved or overloaded, and
  the human dogfood fanned cleanly into three per-observable rows.
- **Token-efficiency opportunities missed?** Minor. The repeat evidence-format
  finding cost some re-run tokens; otherwise the phase was lean. Reusing the
  v0.1b smoke as the held-out seed (recommendation, not relocation) avoids future
  duplicate scenario authoring.

### Deep retrospective recommendation

Not required. Phase 6 was low-friction: zero blocking findings, zero spec
deviations, zero dependency changes, one clean BLOCK→DEC resolution, and a 3/3
human dogfood. The only systemic signal — repeat "paste exact transcript"
improvement findings — is a cheap template tweak, not a retro-worthy theme. Defer
a deep retro until a phase surfaces a recurring or cross-phase problem.
