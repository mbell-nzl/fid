# Phase 14 Validation Contract - Telemetry + Pricing

> Phase: 14 - Telemetry + pricing (roadmap v0.1e, slice 1 of 4)
> Author: Tech Lead (draft)
> Last refined: 2026-06-28 by Tech Lead (planning draft)
>
> Location: `project-docs/docs/phases/phase-14/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 14 is complete. The contract specifies **what must be proven**; the QA
plan specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Phase Scoping Notes

Phase 14 is the first v0.1e slice: the telemetry and pricing surface. It is
self-contained — it depends only on the existing journal/projection machinery
and the existing `issues`/`epic` rows — and it leads the version because the
`fid dev *` composites (Phase 15) open and close the `agent_runs` rows this
phase introduces.

Scope:

- `.fid/pricing.toml` default model set: expand the shipped defaults beyond the
  current Anthropic-only set to include current OpenAI and Google models
  (PRD-NFR-010). Values are user-overridable and snapshotted per run, so exact
  list prices are not a compatibility contract.
- `fid run record`: append one `agent_runs` row capturing role, model, tokens
  in/out, start/end timestamps, snapshot pricing, and computed `cost_usd`, with
  optional provider reconciliation fields and a warning when the model is absent
  from pricing config (PRD-TELEM-001..003a).
- `fid stats tokens` / `fid stats cost`: aggregate tokens and cost by role, by
  model, and by epic across `agent_runs`, in all four read formats
  (PRD-TELEM-004, PRD-OUT-001..004).

The `agent_runs` table and all columns this phase writes already exist
(`migrations/0002_v01_tables.up.sql`; `03_TRD_SCHEMA.md` section 7), so **no
schema migration is required**. The projector applies generic insert effects
against validated table identifiers (`internal/journal/project_effects.go`), so
the `agent_runs` insert needs **no projector change**. The `run.record` action
already appears in `03_TRD_SCHEMA.md` section 13 (all seven roles permitted);
Task 2 appends that one row to the in-code matrix that Phase 10 intentionally
kept limited to implemented actions (DEC-026 pattern).

`internal/config` already defines `PricingConfig`, `ReadPricingConfig`, and
`DefaultPricingConfig`, and `fid init` already writes `pricing.toml` from those
defaults — Task 1 widens the default set and addresses price precision rather
than building the loader from scratch.

Planning interpretations (surfaced to the human for confirmation; recorded as
phase decisions during execution):

- **PI-1 (new `internal/runs` boundary).** The `agent_runs` write
  (`fid run record`) lands in a new domain package `internal/runs`, not in
  `internal/stats` (whose responsibility is read aggregation, TRD section 3).
  This is a new inward-pointing internal boundary, so Task 2 records it in
  `02_TRD.md` sections 2-3. The `agent_runs` table and the `run.record` matrix
  key already exist in the specs; no public interface or data model changes.
- **PI-2 (no `activity_log` for `run.record`).** `fid run record` emits one
  `run.recorded` event with one `agent_runs` insert effect and projects **no**
  `activity_log` row. Telemetry capture is not a domain mutation; this matches
  the `issue create` / `fid label add` no-activity precedent, and CF-011 broad
  reconciliation stays Phase 16.
- **PI-3 (unknown-model behavior).** When the model is present in pricing
  config, the run records snapshot `price_in_per_mtok`/`price_out_per_mtok` and
  computed `cost_usd`, with no warning. When the model is absent, the run is
  still recorded with `cost_usd` and both snapshot price columns left NULL, a
  warning naming the model is written to stderr, and the command exits 0
  (PRD-TELEM-003). Unknown model is not a rejection.
- **PI-4 (by-epic aggregation + CLI doc reconciliation).** `--by epic` groups
  each run under its issue's owning epic (a run on an epic groups under itself;
  a run on an issue with no epic ancestor groups under a single null/`(none)`
  bucket). `06_DESIGN_PLAN_CLI.md` `fid stats tokens` `--by` choices are
  reconciled to include `epic` so both stats commands satisfy PRD-TELEM-004
  ("by role, by model, and by epic"). NULL-cost runs are counted in
  `run_count` and token totals but excluded from the `total_cost_usd` sum.
- **PI-5 (pricing precision).** `pricing.toml` is emitted with enough decimal
  precision that sub-cent default prices (e.g. cheap OpenAI/Google models)
  round-trip through `ReadPricingConfig` without loss; the exact emit format is
  recorded in the Task 1 completion report.

Carry-forward review at kickoff:

- CF-008 remains deferred to v0.2.
- CF-011 (broad `activity_log` reconciliation) remains scheduled for Phase 16
  (with `fid activity list`); Phase 14 adds no `activity_log` rows (PI-2).
- CF-013 / CF-014 / CF-015 remain scheduled for Phase 15 (the `fid dev *`
  composites and their orchestration cleanup); none land in Phase 14.
- Lock-ordering lesson (DEC-028): the new write command `fid run record`
  refreshes the projection with the held-lock freshness helper inside
  `journal.WithLock`; no nested `.fid/write.lock` acquisition.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 14, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | PRD-NFR-010; `04_TRD_CONFIG.md` section 3; roadmap v0.1e Phase 14; PI-5 | `fid init` writes `.fid/pricing.toml` containing default model entries for at least one current Anthropic, one OpenAI, and one Google model, and every shipped default entry round-trips through `ReadPricingConfig` with its `in_per_mtok` and `out_per_mtok` preserved, including sub-cent prices. | Domain + CLI: unit test asserts `DefaultPricingConfig` covers all three providers and round-trips through emit→`ReadPricingConfig`; `fid init` in a temp dir then inspect/parse the generated `pricing.toml`. | QA Lead | task-1 | Defaults include all three providers and every entry round-trips with no precision loss. | A provider is missing, an entry is dropped, or a price loses precision on round-trip. | Draft | Exact list-price values are not a contract (user-overridable, snapshotted). |
| VAL-002 | PRD-TELEM-001; `03_TRD_SCHEMA.md` section 7; `06_DESIGN_PLAN_CLI.md` `fid run record`; PI-1/PI-2 | `fid run record --issue <id> --role <role> --model <model> [--tokens-in <n>] [--tokens-out <n>] --by <role>` appends exactly one `run.recorded` journal event and projects exactly one `agent_runs` row capturing role, model, tokens_in, tokens_out, started_at, ended_at, and (for a known model) the snapshot `price_in_per_mtok`/`price_out_per_mtok`; no `activity_log` row is written. | Domain + CLI: record a run for a known model, read the `agent_runs` row and assert all fields, assert exactly one new journal line and zero new `activity_log` rows. | QA Lead | task-2 | One event, one `agent_runs` row with all listed fields populated and snapshot prices set for a known model, zero `activity_log` rows. | A field is missing, more than one event/row is written, snapshot prices are absent for a known model, or an `activity_log` row appears. | Draft | Run git context and run lifecycle (`status`, `fid run list/show`) are v0.2c. |
| VAL-003 | PRD-TELEM-002; `04_TRD_CONFIG.md` section 3 | When the model is present in pricing config, the recorded `agent_runs.cost_usd` equals `(tokens_in * price_in_per_mtok + tokens_out * price_out_per_mtok) / 1e6` computed from the snapshot prices stored on the same row. | Domain + CLI: record a run with known tokens and a known model; assert `cost_usd` equals the formula applied to the row's snapshot prices within float tolerance. | QA Lead | task-2 | `cost_usd` matches the formula over the row's snapshot prices. | `cost_usd` is null for a known model, or diverges from the formula. | Draft | Provider-billed cost reconciliation is a separate optional field (VAL-005). |
| VAL-004 | PRD-TELEM-003; PI-3 | When the model is absent from pricing config, `fid run record` still records the run with `cost_usd` and both snapshot price columns left NULL, emits a warning naming the model to stderr, and exits 0; when the model is present, no such warning is emitted and `cost_usd`/prices are populated. | CLI: record a run with an unknown model (assert exit 0, warning on stderr naming the model, NULL cost/prices in the row) and a run with a known model (assert no warning, populated cost/prices). | QA Lead | task-2 | Unknown-model run records with NULL cost/prices + warning + exit 0; known-model run has populated cost/prices and no warning. | Unknown model is rejected, silently guesses a cost, omits the warning, or known model emits the warning. | Draft | No silent price guessing for unknown models. |
| VAL-005 | PRD-TELEM-003a; `03_TRD_SCHEMA.md` section 7 | Optional provider reconciliation fields (`--token-breakdown-json`, `--tool-call-count`, `--provider-run-id`, `--provider-cost-usd`) are persisted on the `agent_runs` row when supplied and left NULL when omitted, without altering the normalized token/price/cost summary fields. | Domain + CLI: record one run supplying all provider fields and one omitting them; assert the provider columns are populated/NULL respectively and the normalized fields are identical given identical inputs. | QA Lead | task-2 | Provider fields persist when supplied and are NULL when omitted; normalized summary fields are unaffected. | A provider field is dropped when supplied, is non-NULL when omitted, or overwrites a normalized field. | Draft | `review_id` linkage ships with the v0.3 QA layer. |
| VAL-006 | TRD-ARCH-010; `03_TRD_SCHEMA.md` section 13; CF-007 | The in-code permissions matrix contains exactly one `run.record` row with the same role booleans as `03_TRD_SCHEMA.md` section 13 (all seven roles permitted), `fid policy show` displays that row from the same source `Authorize` enforces, every default-profile role may record a run, and an unknown role is rejected before append with a non-zero exit, an error naming role + `run.record`, and zero new journal lines. | Unit + CLI: matrix completeness test updated for `run.record`; record as each permitted role (allowed) and as an unknown role (rejected before append); `fid policy show` output includes the row and matches the matrix object. | QA Lead | task-2 | Exactly one `run.record` row matching section 13; all known roles allowed; unknown role rejected before append naming role + action; show matches enforce. | The row is missing/duplicated or diverges from section 13, a known role is rejected, an unknown role appends, or show/enforce drift. | Draft | Configurable role permissions remain v0.2 (PRD-ROLE-006). |
| VAL-007 | PRD-TELEM-004; PRD-OUT-001..004; `06_DESIGN_PLAN_CLI.md` section 13; PI-4 | `fid stats tokens [--by role\|model\|epic] [--epic <id>] [--format ...]` aggregates `total_tokens_in`, `total_tokens_out`, and `run_count` grouped by the chosen dimension across `agent_runs`, honoring the `--epic` filter, and json/toon/markdown/yaml carry the same data. | CLI: seed several runs across roles/models/epics; run `stats tokens` for each `--by` value and with `--epic`; compare group totals to a hand-computed reference and compare all four formats to the JSON reference. | QA Lead | task-3 | Group totals and run counts match the reference for each `--by` dimension and the `--epic` filter; all four formats expose the same data. | A group total/count is wrong, a dimension or filter is ignored, or a format diverges from JSON. | Draft | Stats expansion (`tasks`/`attempts`/`gates`/`interventions`) is v0.2c. |
| VAL-008 | PRD-TELEM-004; PRD-OUT-001..004; `05_DESIGN_PLAN.md` section 6.4; PI-4 | `fid stats cost [--by role\|model\|epic] [--epic <id>] [--format ...]` aggregates `total_cost_usd` (plus token totals and `run_count`) grouped by the chosen dimension, where runs with NULL `cost_usd` are counted in `run_count` and token totals but excluded from the cost sum, in all four read formats. | CLI: seed runs including at least one NULL-cost (unknown-model) run; run `stats cost` for each `--by` value and `--epic`; assert cost sums exclude the NULL-cost run while counts/tokens include it, and compare all four formats. | QA Lead | task-3 | Cost sums match the reference excluding NULL-cost runs while counts/tokens include them; all four formats agree. | A NULL-cost run is summed or dropped from counts/tokens, a group is wrong, or a format diverges. | Draft | Provider-billed cost reconciliation reporting is not in scope. |
| VAL-009 | `08_QUALITY_PLAN.md` sections 5.2, 5.6; phase goal | A human dogfooding a real `.fid` workspace can tell from CLI output that an agent run was recorded with its computed USD cost, that recording a run with an unknown model still succeeds but warns, and that `fid stats cost`/`tokens` report aggregates by role, model, and epic. | Human checklist with one row per observable; QA Lead generates from this VAL after deterministic QA passes. | Human | phase-end | Human confirms run recording with cost, the unknown-model warning, and cost/token aggregates by role/model/epic from CLI output without reading implementation. | Human cannot identify those outcomes from CLI output, or the checklist is not run. | Draft | QA Lead fans this into per-observable human-checklist rows. |

## Coverage Rules

Confirmed for Phase 14 execution:

- Phase-scope requirements each map to at least one `VAL-NNN`: PRD-NFR-010 to
  VAL-001; PRD-TELEM-001 to VAL-002; PRD-TELEM-002 to VAL-003; PRD-TELEM-003 to
  VAL-004; PRD-TELEM-003a to VAL-005; TRD-ARCH-010 / section 13 to VAL-006;
  PRD-TELEM-004 to VAL-007/008; PRD-OUT-001..004 to VAL-007/008; phase dogfood
  goal to VAL-009.
- Every `VAL-NNN` is referenced by a task spec or explicitly marked `phase-end`
  (VAL-009).
- Binary/conditional assertions name both states: VAL-001, VAL-004, VAL-005,
  VAL-006, VAL-007, VAL-008.
- VAL-001 carries a round-trip conformance proof; VAL-006 carries a
  matrix-vs-section-13 and show-vs-enforce parity proof; VAL-007/008 carry
  four-format parity-vs-JSON proofs.
- Coverage is enforced at the gate, not as a `VAL-NNN`: domain package coverage
  and `cmd/fid` coverage must stay above thresholds. Each task runs the required
  coverage gate and records package percentages.
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase review, or
  `(human)` verification, and fans VAL-009 into per-observable human rows.
- Standing QA-plan checks carry forward: show-vs-enforce parity for the new
  `run.record` matrix row (VAL-006), validate-before-append on the unknown-role
  rejection (VAL-006), no `activity_log` write for telemetry (VAL-002), and no
  stale-projection nested-lock failure on the new write path (VAL-002).

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-28 | VAL-001..009 | created | Phase 14 planning draft | roadmap v0.1e Phase 14; PRD-TELEM-001..004; PRD-NFR-010; TRD-ARCH-010 |
