# Phase 17 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes,
> and wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-049: Phase 17 runs in Heavy Mode

- **Date:** 2026-07-11
- **Context:** Phase 17 adds public CLI and config contracts, performs an
  explicit derived-DB mutation under the projection lock, and closes v0.1.
  `09_AI_WORKFLOW.md` section 4.3 classifies public-interface, persistence, and
  release-adjacent work as Heavy Mode.
- **Decision:** The Human approved **Heavy Mode**. Required artifacts are the
  validation contract, nine task specs, a QA plan with fresh command and fixture
  feasibility validation, reviews after Tasks 2, 5, and 7, final review, human
  dogfood, and full-scope spec sync at wrap.
- **Consequences:** No Developer task starts until the QA plan is Ready. Phase
  17 is v0.1 early-dogfood closure, not the charter's v0.3 public release;
  public-release threat model, scans, pre-release checklist, and held-out
  scenarios remain deferred under the quality plan. DEC-054 later supersedes
  only the nine-task count and review cadence after QA entry findings.
- **Refs:** `09_AI_WORKFLOW.md` section 4.3; `08_QUALITY_PLAN.md` sections
  5.8, 7, and 9; phase brief; VAL-001..025.

### DEC-050: Maintenance warning definitions and grace config are fixed

- **Date:** 2026-07-11
- **Context:** PRD-OPS-008 required a configurable started-without-claim grace,
  while the roadmap also named stale verifications and telemetry gaps. The
  canonical specs did not define the key, boundary, time source, or linkage
  semantics. Historical `scaffold/initial-spec.md` section 9.8.1 retained the
  original five-minute and superseded-attempt intent.
- **Decision:** Add
  `[maintenance].started_without_claim_grace_seconds`, default 300 seconds;
  zero is valid and negative/malformed values error. Missing config defaults to
  300 without rewriting existing files. Age uses the newest append-ordered
  `status_changed` activity, falling back to `issues.created_at`; non-status
  updates do not reset it. Stale verification means non-NULL `completion_id`
  referencing a superseded completion. Telemetry gap means a completion without
  an `agent_runs.completion_id` reference; issue-only run linkage is not enough.
- **Consequences:** Post-development handoffs can warn after five minutes when
  a started-category issue has no claimant. A completion without full linked
  telemetry can warn even when an issue-level dev run exists. These remain
  warnings, not structural errors.
- **Refs:** PRD-OPS-007/008; `04_TRD_CONFIG.md` section 2; DC-A;
  VAL-002/007/008/009/010.

### DEC-051: Maintenance results use pass, warn, and fail semantics

- **Date:** 2026-07-11
- **Context:** The design spec showed only a failing JSON example and did not
  define warning-only status, warning promotion, output ordering, or how a
  structured result survives a non-zero exit.
- **Decision:** No findings returns `pass`/zero; warnings-only returns
  `warn`/zero; warnings-only with `--warn-as-error` returns `fail`/non-zero; any
  structural error returns `fail`/non-zero. Promotion changes only status and
  exit: rows stay in `warnings`. Render the complete result to stdout before an
  audit-outcome non-zero return. Fatal config/DB/encoding failures use the
  normal stderr error path. Findings order by check, entity ID, then message.
  The command supports json/toon/markdown/yaml.
- **Consequences:** CI can parse structured non-zero audit output without
  confusing promoted operational warnings with data corruption. `warn` becomes
  a public result enum and all four formats follow the standing read contract.
- **Refs:** PRD-OPS-005..008; PRD-OUT-001..004; DC-B;
  `05_DESIGN_PLAN.md` section 6.5; `06_DESIGN_PLAN_CLI.md` section 19;
  VAL-011/012/013.

### DEC-052: The v0.1 governance stub guard covers three complete families

- **Date:** 2026-07-11
- **Context:** The v0.1 schema pre-creates decision, escalation, and triage
  tables for references/views, while their command behavior is deferred. The
  roadmap's `decision create`, `escalation create`, "etc." wording left the
  exact guarded set and validation precedence ambiguous.
- **Decision:** Register and version-gate `decision create/get/list`,
  `escalation create/resolve/get/list`, and `triage record/start/verdict`. Each
  returns `error: commands for <entity> ship in v0.2` non-zero before argument,
  flag, actor, workspace, DB, or journal handling. It appends nothing. Keep the
  shipped `requirement sync` unchanged; requirement pointers are not a
  schema-only stub. Use public version `v0.2`, not internal slice `v0.2b`.
- **Consequences:** Root help can expose honest future command families without
  partial behavior. Actual governance CRUD/reads and their permissions remain
  v0.2 work.
- **Refs:** PRD-ERR-012; DC-C; `06_DESIGN_PLAN_CLI.md` section 24;
  `10_IMPLEMENTATION_ROADMAP.md` Phase 17; VAL-016.

### DEC-053: Codex roles use evidence-matched GPT-5.6 routing

- **Date:** 2026-07-11
- **Context:** Phase 17 introduces GPT-5.6-class models to the Codex role
  routing. Model selection must match each role's scope, preserve the proven
  Phase Runner configuration, and avoid adopting higher-cost tiers without
  role-specific evidence.
- **Decision:** Use the following Codex routing:

  | Role | Model | Reasoning |
  |---|---|---:|
  | Tech Lead | `gpt-5.6-sol` | xhigh |
  | Phase Runner | `gpt-5.4-mini` | xhigh |
  | Developer | `gpt-5.6-terra` | high |
  | QA Lead | `gpt-5.6-sol` | xhigh |
  | QA Reviewer | `gpt-5.6-terra` | high |
  | TL-Lite | `gpt-5.6-terra` | medium |

  Sol is reserved for phase-wide planning and QA synthesis. Terra handles
  scoped implementation, review, and tactical triage. Phase Runner remains on
  `gpt-5.4-mini` at xhigh because that routing is proven in this workflow and
  the Luna cost increase lacks role-specific evidence. No max or ultra routing
  is adopted.
- **Consequences:** Phase 17 must not trust the new QA Lead or QA Reviewer
  routing until model-audit evidence confirms that seeded defects are detected,
  deterministic evidence is not replaced by prose inference, `VAL-NNN`
  failures and blockers are dispositioned correctly, QA Reviewer escalates
  cross-cutting findings, and QA Lead preserves hard stops without weakening
  assertions. Phase Runner needs no new orchestrator-model qualification
  because its model and effort are unchanged; the prompt-change regression
  drill still applies. Existing Codex sessions must be replaced with fresh
  sessions before Phase 17 execution relies on this routing.
- **Refs:** `09_AI_WORKFLOW.md` section 7; `project-docs/agents/`;
  `project-docs/tools/drills/run-drills.sh`.

### DEC-054: Phase 17 uses a tests-only entry remediation, not a coverage waiver

- **Date:** 2026-07-11
- **Context:** The QA Lead's fresh feasibility run found `internal/journal` at
  78.9%, below the existing COV-001 80% threshold. The approved Tasks 1-9 do
  not authorize journal coverage work, so VAL-024 had no path to PASS.
- **Decision:** Add Task 0 before product work. It may add focused
  `internal/journal` behavioral tests only and must prove at least 80% in both
  the focused and unchanged full-suite coverage runs. Production-code changes,
  threshold changes, and waivers are prohibited. A QA Reviewer entry review
  must accept Task 0's test quality and evidence before Task 1 starts; QA Lead
  retains terminal VAL-024 ownership.
- **Consequences:** Phase 17 now has ten tasks and reviews after Tasks 0, 2, 5,
  and 7. The phase remains blocked on ESC-009 and on a QA Lead revision that
  authorizes Task 0 and reconciles the QA plan's inventory, evidence, cadence,
  report-lint commands, and entry criteria. The Phase Runner's append-only
  ledger remains untouched and is reconciled when the run resumes.
- **Refs:** QA-BLK-01; COV-001; VAL-024; `tasks/task-0.md`; ESC-009.

### DEC-055: Phase 17 must qualify the approved QA routing

- **Date:** 2026-07-11
- **Context:** ESC-009 offered qualification of the DEC-053 QA routing or a
  return to the Phase 16 QA models as possible resolutions to QA-BLK-02.
- **Decision:** The Human rejected the routing fallback under all
  circumstances. Phase 17 must qualify `gpt-5.6-sol` xhigh for QA Lead and
  `gpt-5.6-terra` high for QA Reviewer through the seeded-defect model audit.
- **Consequences:** ESC-009 has one valid resolution: committed, Human-scored
  evidence satisfying every DEC-053 obligation with no misses. A model change,
  waiver, inferred pass, or reliance on the blocked QA plan alone cannot unblock
  Developer execution.
- **Refs:** DEC-053; ESC-009; QA-BLK-02;
  `project-docs/tools/drills/README.md`.

### DEC-056: Rebuild parity covers derived state, not clock-relative reads

- **Date:** 2026-07-17
- **Context:** ESC-010 (TRG-008) reported a shipped rebuild-parity defect: after
  deleting the derived DB, `fid rebuild --by workflow_runner` left the journal
  byte-identical but `fid maintenance check` reported `recent_activity_7d` as 17
  before rebuild and 20 after. Task 9 stopped and VAL-022 was recorded BLOCKED.
- **Decision:** The canonical parity contract is defined as follows.
  1. **The journal is the sole source of truth for derived state.** After
     delete-and-rebuild, every projected table must be exactly equal to its
     pre-delete content. This is the substance of VAL-022 and it is not relaxed.
  2. **Command output is not derived state.** A read whose value is a function of
     `(projected state, wall-clock now)` is not a rebuild-parity observable,
     because `now` is an input to the read, not something the journal reproduces.
     Rebuild cannot and must not be expected to hold such a value constant across
     two invocations separated by real elapsed time.
  3. In `maintenance check`, `info.recent_activity_7d` and `info.claim_count_live`
     are both clock-relative by construction (`countRecentActivity` filters on a
     `now`-anchored window; `countLiveClaims` evaluates TTL liveness at `now`).
     They are excluded from equality parity. `errors`, `warnings`, `status`, and
     `info.issue_count` are state-derived and remain under exact parity.
- **Root-cause finding (to be confirmed by the DEC-056 diagnostic):** The
  reported 17 -> 20 movement is consistent with a clock artifact, not lost or
  duplicated projection. Issue writes stamp `OccurredAt` via
  `nextIssueUpdateTime` (`internal/issues/write_helpers.go`), which preserves
  monotonicity by returning `prior + 1s` whenever the wall clock has not advanced
  past the previous `updated_at`. A fast journey performs many issue mutations
  inside the same second, so the last few activity rows carry `occurred_at`
  values slightly **ahead** of the wall clock. `countRecentActivity` excludes
  future rows (`!at.After(now)`). The pre-rebuild check therefore omits the
  future-stamped rows; the post-rebuild check runs seconds later, by which time
  the clock has caught up and the same rows count. Same rows, different `now`.
- **Consequences:**
  - ESC-010's premise is not accepted as established. No product change to
    projection, rebuild, or maintenance is authorized on the current evidence.
  - Task 9 resumes as an audit-only task and must first run the discriminating
    diagnostic in `tasks/task-9.md` ("Required Diagnostic"). The diagnostic, not
    a hypothesis, decides the branch.
  - Branch A (control reproduces the movement with no rebuild): the defect is in
    the assertion, not the product. Adopt the amended VAL-022 method and close.
  - Branch B (projected `activity_log` rows genuinely differ after rebuild):
    that is a real persistence defect. Stop again, do not amend VAL-022, and
    re-escalate for a scoped remediation task.
  - The future-stamped-`occurred_at` interaction is recorded as CF-020 for v0.2
    rather than fixed inside a v0.1 closure audit task; it is a deliberate
    monotonicity trade-off with a bounded, self-correcting effect on one info
    counter.
- **Refs:** ESC-010; TRG-008; VAL-022/023; `tasks/task-9.md` binding constraint
  "compare structurally without timestamp-string fragility";
  `internal/issues/write_helpers.go`; `internal/maintenance/info.go`;
  `internal/journal/projection.go`; `internal/journal/replay.go`.

## Carry-Forward Review At Kickoff

| ID | Disposition | Rationale |
|---|---|---|
| CF-008 | Re-defer to v0.2 | Resolution reason is an unrelated product/data-model enhancement. |
| CF-017 | Re-defer to v0.2c | Phase 17 does not modify `internal/dev`; Task 9 consumes it only through CLI tests. |
| CF-018 | Re-defer to next SQL package reopen | Task 8 must not modify `internal/sqlquery`; any required change stops for CF-018 re-planning. |
| CF-019 | Schedule in Phase 17 Task 1 | Task 5 reopens markdown rendering, so the behavior-preserving split lands first. |

## Kickoff Spec Sync

| Spec | Kickoff action | Decision reference |
|---|---|---|
| `01_PRD.md` | Clarified PRD-OPS-008 and PRD-ERR-012. | DEC-050/052 |
| `04_TRD_CONFIG.md` | Added the v0.1 maintenance grace contract. | DEC-050 |
| `05_DESIGN_PLAN.md` | Added maintenance result/exit semantics. | DEC-051 |
| `06_DESIGN_PLAN_CLI.md` | Added maintenance truth table and exact v0.1 guard matrix. | DEC-051/052 |
| `10_IMPLEMENTATION_ROADMAP.md` | Replaced the Phase 17 stub `etc.` with the approved command set. | DEC-052 |
| `GLOSSARY.md` | Defined started-without-claim, stale verification, and telemetry gap. | DEC-050 |

Full canonical spec-sync verdicts remain mandatory at phase wrap; this kickoff
table does not substitute for the wrap audit.

## Interim Workflow Notes (for wrap-up retrospective)

### WN-1: Recurring completion-report / report-lint friction

- **Observed:** Three of the first four TL-Lite triages this phase were
  report-lint stumbles — TRG-002 (Task 0 commit metadata not backtick-wrapped),
  TRG-003 (Task 1 commit metadata not backtick-wrapped, same defect a second
  time), and TRG-004 (Task 1 report left uncommitted after the metadata fix,
  invalidating the lint evidence → ESC-004). The metadata-format defect
  recurring across two separate tasks/sessions is a pattern, not a one-off.
- **Likely contributing factor (flagged by Human 2026-07-15):** The subagent
  prompts and the Codex GPT model routing (DEC-053) were both updated recently.
  The Developer role now runs a newly routed model against a revised prompt, so
  the recurrence may be a prompt/model-fit symptom rather than coincidence.
- **Wrap-up action:** Review and tighten the Developer prompt/instructions
  alongside the completion-report template and report-lint contract. Assess
  specifically whether (a) the template needs a clearer worked example of the
  backtick-wrapped `## Commit hash` / `## Commit subject` block, (b) the
  Developer prompt must state explicitly that the completion report is its own
  committed `docs(task)` step and must be clean before report-lint, and
  (c) report-lint's error messaging should call out uncommitted/dirty report
  state distinctly from parse failures. Decide at wrap-up whether this becomes a
  DEC (template/tooling change) or a prompt-source edit under
  `project-docs/agents/` + `sync-agents.sh`.
- **Refs:** TRG-002/003/004; ESC-004; DEC-053; completion-report template;
  `project-docs/tools/report-lint.sh`; `project-docs/agents/developer.md`.
