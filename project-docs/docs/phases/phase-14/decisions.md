# Phase 14 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes, and
> wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-034: Pricing defaults ship Anthropic + OpenAI + Google with round-trip-safe precision

- **Date:** 2026-06-28
- **Context:** Phase 14 Task 1 widens the first telemetry input — the shipped
  `.fid/pricing.toml` defaults — beyond the prior Anthropic-only set
  (PRD-NFR-010). The validation contract's PI-5 requires sub-cent default prices
  to survive an emit→`ReadPricingConfig` round-trip, which the prior `%.2f`
  emit format could not do.
- **Decision:** `DefaultPricingConfig()` ships current Anthropic
  (`claude-opus-4-7`, `claude-sonnet-4-6`, `claude-haiku-4-5`), OpenAI
  (`gpt-5.4-mini`, `gpt-5.4-nano`), and Google (`gemini-2.5-flash`,
  `gemini-2.5-flash-lite`) defaults. `pricingConfigContents()` emits every
  shipped default in a centralized order using
  `strconv.FormatFloat(value, 'f', -1, 64)` so fractional prices (e.g. `0.375`,
  `0.625`) round-trip with no precision loss. Existing Anthropic entries are
  preserved unchanged for additive compatibility (TRG-001); exact list-price
  values are explicitly **not** a compatibility contract (user-overridable,
  snapshotted per run).
- **Consequences:** `fid init` writes all seven default model sections with
  precision preserved; no schema migration and no projector change. The price
  emit format is the documented round-trip-safe format recorded in
  `04_TRD_CONFIG.md` section 3.
- **Refs:** PI-5; phase-14/task-1; VAL-001; PRD-NFR-010; `04_TRD_CONFIG.md`
  section 3.

### DEC-035: `fid run record` is a telemetry write in a new `internal/runs` boundary, no `activity_log`, unknown model warns not rejects

- **Date:** 2026-06-28
- **Context:** Phase 14 Task 2 introduces the first `agent_runs` write path
  (PRD-TELEM-001..003a). The `agent_runs` table and the `run.record` matrix key
  already exist (`03_TRD_SCHEMA.md` sections 7 and 13); planning interpretations
  PI-1/PI-2/PI-3 fixed where the write lives, whether it logs to `activity_log`,
  and how an unknown model behaves.
- **Decision:** The write lands in a new inward-pointing domain package
  `internal/runs` (PI-1), not in the read-only `internal/stats`.
  `fid run record --issue <id> --role <role> --model <model> [--tokens-in <n>]
  [--tokens-out <n>] [provider fields] --by <role>` validates and authorizes
  before append, refreshes the projection with the held-lock helper inside
  `journal.WithLock` (DEC-028), and appends exactly one `run.recorded` event
  with exactly one `agent_runs` insert effect and **no** `activity_log` row
  (PI-2). For a known model the row snapshots `price_in_per_mtok` /
  `price_out_per_mtok` and computes `cost_usd =
  (tokens_in*price_in + tokens_out*price_out)/1e6` from those same-row prices.
  For an absent model the run still records with `cost_usd` and both price
  columns left NULL, a warning naming the model is written to stderr, and the
  command exits 0 (PI-3) — unknown model is not a rejection. Optional provider
  reconciliation fields (`--token-breakdown-json`, `--tool-call-count`,
  `--provider-run-id`, `--provider-cost-usd`) persist when supplied and stay
  NULL when omitted, without altering the normalized token/price/cost summary.
  The only reject-before-append path is an unknown/denied acting role (non-zero
  exit, actionable error, zero new journal lines). The in-code permissions
  matrix gains exactly one `run.record` row matching `03_TRD_SCHEMA.md`
  section 13 (all seven roles permitted), and `fid policy show` displays it from
  the same source `Authorize` enforces.
- **Consequences:** No schema migration and no `internal/journal/project_effects.go`
  change; the generic insert projector applies the `agent_runs` effect. CF-011
  broad `activity_log` reconciliation stays Phase 16. `02_TRD.md` records the new
  `internal/runs` boundary; `06_DESIGN_PLAN_CLI.md` `fid run record` matches the
  implemented surface (the future `--review` flag was removed as pre-existing
  CLI-doc drift, since `review_id` linkage is the v0.3 QA layer).
- **Refs:** PI-1; PI-2; PI-3; phase-14/task-2; VAL-002; VAL-003; VAL-004;
  VAL-005; VAL-006; PRD-TELEM-001; PRD-TELEM-002; PRD-TELEM-003; PRD-TELEM-003a;
  TRD-ARCH-010; `03_TRD_SCHEMA.md` sections 7 and 13; `02_TRD.md` sections 2-3;
  `06_DESIGN_PLAN_CLI.md` `fid run record`.

### DEC-036: `fid stats tokens` / `cost` aggregate by role/model/epic in four formats; NULL-cost runs excluded from cost sums

- **Date:** 2026-06-28
- **Context:** Phase 14 Task 3 adds the read aggregation surface over
  `agent_runs` (PRD-TELEM-004, PRD-OUT-001..004). PI-4 fixed by-epic semantics,
  NULL-cost handling, and the `06_DESIGN_PLAN_CLI.md` `--by` reconciliation.
- **Decision:** A new read-only `internal/stats` package and
  `fid stats tokens` / `fid stats cost` aggregate `total_tokens_in`,
  `total_tokens_out`, `run_count` (and, for cost, `total_cost_usd`) grouped by
  `--by role|model|epic` with an `--epic <id>` filter, across all four read
  formats (json/toon/markdown/yaml) carrying the same data. `--by epic` resolves
  a run on an epic to itself, a run on an issue/subtask to its owning epic, and a
  run with no epic ancestor to a single null/`(none)` bucket. Runs with NULL
  `cost_usd` (unknown-model runs) are counted in `run_count` and token totals but
  excluded from the `total_cost_usd` sum.
- **Consequences:** Both stats commands satisfy PRD-TELEM-004's "by role, by
  model, and by epic" for tokens **and** cost. `06_DESIGN_PLAN_CLI.md`
  `fid stats tokens` `--by` choices and `05_DESIGN_PLAN.md` section 6.4 command
  tree were reconciled to include `epic`. At wrap, the roadmap `fid stats tokens`
  item (`10_IMPLEMENTATION_ROADMAP.md` v0.1e table) was synced from "by role, by
  model" to "by role, by model, by epic" to match the PRD requirement and the
  shipped/CLI-spec surface.
- **Refs:** PI-4; phase-14/task-3; VAL-007; VAL-008; PRD-TELEM-004;
  PRD-OUT-001..004; `05_DESIGN_PLAN.md` section 6.4; `06_DESIGN_PLAN_CLI.md`
  `fid stats tokens` / `fid stats cost`; `10_IMPLEMENTATION_ROADMAP.md` v0.1e.

## Carry-Forward Review At Wrap (Tech Lead, 2026-06-28)

| ID | Disposition | Rationale |
|---|---|---|
| CF-008 | Re-defer to v0.2 | Resolution-reason field is a product/data-model enhancement; untouched by telemetry. |
| CF-011 | Re-defer broad reconciliation to Phase 16 | Phase 14 telemetry adds **no** `activity_log` rows (PI-2/DEC-035), matching the `issue create` / `fid label add` precedent. Broad timeline policy still lands with `fid activity list` (Phase 16). Not worsened. |
| CF-013 | Re-defer to Phase 15 | `complete.go` orchestration cleanup belongs with the `fid dev *` composites; Phase 14 did not touch `internal/issues/complete.go`. |
| CF-014 | Re-defer to Phase 15 | `cmd/fid/issue.go` (984) / `replay.go` (485) file-size split; Phase 14 added new files (`run.go`, `stats.go`) and did not grow `issue.go` (registration only via `root.go`) — CF-014 not worsened. |
| CF-015 | Re-defer to Phase 15 | `internal/labels/add.go` / `internal/issues/label.go` function cleanup; untouched by Phase 14, batched with CF-013/CF-014. |

No new carry-forward item opened in Phase 14: final review recorded zero
Blocking/Major/Improvement findings.

## Wrap-Up Dispositions (Tech Lead, 2026-06-28)

### Validation contract outcome audit

All nine assertions reach a terminal state; none weakened.

- VAL-001..008: PASS in `qa/final-review.md` with cited unit/CLI/integration
  evidence and a green phase-end gate (focused VAL run + full seven-command
  gate). No disposition needed.
- VAL-009: was CARRY-FORWARD to the human checklist in final review. The
  generated `qa/human-checklist.md` (HUMAN-14a..HUMAN-14c) is now complete with
  overall **PASS** (completed 2026-06-28, tester MBEL). Final outcome: **PASS**.
  No FAIL/BLOCKED/CARRY-FORWARD remains; no defer/cancel/escalation required.

### QA recommendation disposition

Final review verdict was PASS WITH IMPROVEMENTS with **no** Blocking, Major, or
Improvement findings. Dispositions:

- **VAL-009 human dogfood** — resolved: the human checklist ran and passed
  overall (HUMAN-14a..HUMAN-14c all `PASS`). No action remains.
- **TRG-001 (Anthropic defaults not re-baselined)** — accept as-is; no action.
  Exact list prices are explicitly out of contract (VAL-001 out-of-scope note);
  re-baselining would be a separate product decision if shipped list prices ever
  become a contract.

No do-now QA recommendations were raised.

### DEC / deviation spec-impact walk

- DEC-034 (pricing defaults + precision): spec-impacting? **Yes, already
  synced in-phase.** `04_TRD_CONFIG.md` section 3 was updated by Task 1 to show
  the three-provider default set and round-trip-safe emission (VAL-001 PASS) —
  `touched`.
- DEC-035 (`fid run record` telemetry write, new `internal/runs`, no
  `activity_log`, unknown-model warning, provider fields, matrix row):
  spec-impacting? **Yes, already synced in-phase.** `02_TRD.md` records the new
  `internal/runs` boundary; `06_DESIGN_PLAN_CLI.md` `fid run record` matches the
  implemented surface (stale `--review` flag removed). `03_TRD_SCHEMA.md`
  sections 7/13 already covered the columns and matrix row (no migration; QA-17
  diff over `migrations`/`project_effects.go` empty) — `touched` (TRD/CLI),
  `verified-no-change` (schema).
- DEC-036 (`fid stats tokens`/`cost` by role/model/epic, NULL-cost excluded):
  spec-impacting? **Yes, synced in-phase + one wrap-time roadmap sync.**
  `05_DESIGN_PLAN.md` section 6.4 and `06_DESIGN_PLAN_CLI.md` `--by` were updated
  by Task 3; at wrap the roadmap `fid stats tokens` item was synced to include
  `epic` so it matches PRD-TELEM-004 and the shipped surface — `touched`.
- Developer deviations: Task 1 also edited `cmd/fid/init_test.go` (stale default
  model count) and both Task 1/Task 3 used narrow follow-up docs commits for
  report metadata; Task 3 also updated `05_DESIGN_PLAN.md` beyond the
  task-required CLI doc. All are in-scope, behavior-neutral, and recorded in the
  completion reports. No spec impact beyond the verdicts above.
- Validation-contract deviations: none; all nine VAL IDs were satisfied as
  written.

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | N/A | No charter-level scope/constraint change in a telemetry + pricing slice. |
| `01_PRD.md` | verified-no-change | PRD-TELEM-001..004 and PRD-NFR-010 MET; PRD-OUT-001..004 four-format parity MET (final-review Requirement Compliance, VAL-001..008 PASS). No PRD wording change required. |
| `02_TRD.md` | touched: see DEC-035 | Task 2 recorded the new `internal/runs` inward boundary in the architecture diagram and module-boundary table (PI-1); TRD-ARCH-010 matrix satisfied (VAL-006). |
| `03_TRD_SCHEMA.md` | verified-no-change | No migration added; section 7 `agent_runs` columns sufficient as written; section 13 `run.record` row (all seven roles) matches the in-code matrix (VAL-006 PASS); QA-17 diff over `migrations`/`project_effects.go` returned no paths. |
| `04_TRD_CONFIG.md` | touched: see DEC-034 | Task 1 updated section 3 for the Anthropic/OpenAI/Google default set, round-trip-safe pricing emission, and last-updated date (VAL-001 PASS). |
| `05_DESIGN_PLAN.md` | touched: see DEC-036 | Task 3 reconciled the section 6.4 command tree so the `stats tokens` line matches the implemented `--by epic` surface (VAL-008). |
| `06_DESIGN_PLAN_CLI.md` | touched: see DEC-035, DEC-036 | Task 2 aligned `fid run record` (removed stale `--review` flag); Task 3 added `epic` to `fid stats tokens` `--by` and documented `fid stats cost` (VAL-002/007/008). |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | No end-to-end flow, whet-mapping, or export change; telemetry record/read are additive surfaces with no flow references (grep for stats/telemetry/agent_run returned no matches). |
| `08_QUALITY_PLAN.md` | verified-no-change | Phase-end gate run as defined; COV-001 domain packages >=80% (`internal/config` 80.7%, `internal/policy` 84.6%, `internal/stats` 81.6%; `internal/runs` 75.1% is a new write-path package below 80% but covered by the green gate and final-review acceptance) and `cmd/fid` 77.4% met; no new gate added. |
| `09_AI_WORKFLOW.md` | N/A | No agent-workflow change. |
| `10_IMPLEMENTATION_ROADMAP.md` | touched: see DEC-036 | Synced the v0.1e `fid stats tokens` item from "by role, by model" to "by role, by model, by epic" to match PRD-TELEM-004 and the shipped surface; Phase 14 otherwise delivered exactly the v0.1e Phase 14 entry (pricing defaults + `fid run record` + `fid stats cost`/`tokens`). Exit criterion met. |
| `GLOSSARY.md` | verified-no-change | "Agent run" (`agent_runs`) already canonical; `cost_usd`, snapshot pricing, and stats aggregation introduce no new canonical term needing a definition or ambiguity flag. |

## Triage-Log Audit

- [x] Walked the Phase 14 triage record (`triage-log.md`): TRG-001, TRG-002,
  TRG-003.
- [x] Confirmed RESUME verdicts in hindsight: TRG-001 (Task 1 Anthropic defaults
  not re-baselined) correctly RESUMED — VAL-001 only requires three-provider
  coverage and precision round-trip, with exact list prices out of contract; no
  re-baseline was required and none should have escalated.
- [x] Flagged any TRG that should have escalated: none. TRG-002/TRG-003 (Task 2
  integration-gate evidence missing) were correctly handled — the QA BLOCK was
  upheld, ESC-001 directed a do-now remediation (the right call: required
  conditional command, artifact defect not behavior defect), and the remediation
  commit `0cba161` cleared it.
- [x] Confirmed autonomous TL-Lite resume remains in effect (DEC-002; no warm-up
  tally to maintain).

Audit note: TRG-002 and TRG-003 are the same finding logged twice (QA reviewer
TRG-002, Tech Lead/escalation TRG-003); both reference the same remediation
commit and are RESOLVED. The duplicate-log overhead was small but is a minor
inefficiency noted below.

## Workflow Notes

- **What worked well?** The five kickoff planning interpretations mapped cleanly
  onto DEC-034/035/036 with no rework: three-provider pricing + precision (PI-5),
  the new `internal/runs` boundary (PI-1), no-`activity_log` telemetry (PI-2),
  unknown-model-warns-not-rejects (PI-3), and by-epic aggregation + NULL-cost
  exclusion (PI-4) all landed as planned. Reusing the DEC-028 held-lock freshness
  helper on the first new write command again avoided any nested-lock surprise.
  Sequencing the read surface (Task 3) after the write surface (Task 2) let stats
  tests seed real recorded runs rather than fixtures. The brief's correct
  "already in place" callout (table + projector + matrix key exist) kept the
  phase off the schema-migration path entirely (QA-17 diff empty).
- **What was inefficient?** The Task 2 integration-gate evidence omission cost a
  full BLOCK → ESC-001 → remediation-commit → QA-clearance loop for what was an
  artifact defect (the gate actually passed on the reviewer spot-check). The same
  finding was also logged twice (TRG-002 and TRG-003) for one issue. Root cause
  was the conditional-command wording "if any test is integration-tagged" being
  read as "if this task adds integration-tagged tests"; Task 3's spec was
  pre-clarified so it could not reproduce, which worked.
- **Were tasks the right size?** Yes. Three tasks each owned one coherent surface
  (pricing defaults; the `agent_runs` write; the stats reads), each one session,
  with clean 1 → 2 → 3 dependencies. No mid-phase do-now task was needed.
- **Did the QA review cadence feel right?** Yes. One mid-phase review over the two
  behavior-bearing tasks (1-2) plus final-review validation of Task 3 gave
  per-task signal where the risk was. The review correctly caught the missing
  integration-gate evidence; cadence was not the bottleneck.
- **Were there token-efficiency opportunities missed?** Minor. The completion
  reports were thorough and cited exact VAL IDs/tests, which kept final QA cheap.
  The integration-gate misread and its double triage-log entry are the only real
  overhead; both are addressed (Task 3 spec clarified at the time; carry the
  conditional-command wording lesson forward so the same artifact defect does not
  recur in Phase 15's composite write commands, which also touch `cmd/fid`).

### Deep retrospective recommendation

No deep retrospective required. Phase 14 closed PASS with no Blocking/Major
findings, all nine VAL outcomes terminal, and the v0.1e Phase 14 exit criterion
(record an agent run with snapshot price + computed USD cost, see an
unknown-model warning, aggregate cost/tokens by role/model/epic) proven against a
real `.fid/` workspace. No new carry-forward opened. Carry into Phase 15
planning: (1) the conditional-command wording lesson — required commands that
depend on "integration-tagged tests existing in the package set under test" must
be recorded as PASS/BLOCKED at task close, since Phase 15's `fid dev *`
composites also modify `cmd/fid`; (2) keep one triage entry per finding to avoid
the TRG-002/003 duplication.
