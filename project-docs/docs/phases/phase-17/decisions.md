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
