# Phase 14 - Task 1: Pricing defaults (Anthropic + OpenAI + Google) + precision

## Objective

Widen the shipped `.fid/pricing.toml` default model set so that `fid init`
seeds defaults for current Anthropic, OpenAI, and Google models (today only
Anthropic is shipped), and ensure the emitted file preserves sub-cent price
precision so every default entry round-trips through `ReadPricingConfig`. This
de-risks the cost computation that Task 2's `fid run record` depends on.

## Requirements

- PRD refs: PRD-NFR-010
- TRD refs: `02_TRD.md` section 8 (`.fid/pricing.toml` cost computation)
- Design refs: N/A
- Quality refs: `08_QUALITY_PLAN.md` (file/function size budgets)

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-14/qa/validation-contract.md`.

- VAL-001: `fid init` writes `pricing.toml` with defaults for all three
  providers and every entry round-trips through `ReadPricingConfig` without
  precision loss (PI-5).

## Scope

### Files To Create Or Modify

- `internal/config/config.go` - extend `DefaultPricingConfig()` with current
  OpenAI and Google models alongside the existing Anthropic entries; extend the
  ordered list in `pricingConfigContents`; fix the price emit format so sub-cent
  values round-trip (current `%.2f` loses precision).
- `internal/config/config_test.go` - tests for three-provider coverage and
  emit→`ReadPricingConfig` round-trip including a sub-cent price.
- `project-docs/specs/04_TRD_CONFIG.md` - section 3: reflect that the shipped
  default set spans Anthropic, OpenAI, and Google (the inline example is
  illustrative), and note the precision behavior. Update "Last updated".

### Files To Read For Context

- `internal/config/config.go` (existing `DefaultPricingConfig`,
  `pricingConfigContents`, `ReadPricingConfig`)
- `project-docs/specs/04_TRD_CONFIG.md` section 3
- `project-docs/specs/03_TRD_SCHEMA.md` section 7 (pricing source note)

### Out Of Scope

- `fid run record`, cost computation, and stats (Tasks 2-3).
- Any `pricing.toml` reload/override CLI; users edit the file directly.
- Locking exact list-price values as a contract; values are user-overridable.

## Acceptance Criteria

- [ ] `DefaultPricingConfig()` returns entries for at least one current
  Anthropic, one OpenAI, and one Google model. (VAL-001)
- [ ] `pricingConfigContents` emits every default entry with enough precision
  that a sub-cent price round-trips through `ReadPricingConfig` unchanged.
  (VAL-001)
- [ ] `fid init` in a temp workspace writes a `pricing.toml` whose parsed models
  match `DefaultPricingConfig()`. (VAL-001)
- [ ] `04_TRD_CONFIG.md` section 3 reflects the three-provider default set and
  the precision behavior. (doc)

## TDD Plan

- First failing test: assert `DefaultPricingConfig()` contains an OpenAI and a
  Google model key (currently absent) — fails red.
- Additional tests: emit→parse round-trip equality for the full default set
  including a deliberately sub-cent price; existing Anthropic entries unchanged.
- Regression risks: changing the emit format must not break `ReadPricingConfig`
  parsing or the existing `init` golden expectations; keep `config.go` under the
  400-line file budget and functions under 60 lines.

## Documentation Updates

- `project-docs/specs/04_TRD_CONFIG.md` - section 3 default-set + precision note;
  bump "Last updated".

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: None (no read-format surface in this task).
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command in
the completion report as `PASS`, or `BLOCKED` with reason. Record the
`internal/config` package coverage percentage.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` then `fid init` in a temp dir to inspect the
  generated `pricing.toml` for the round-trip evidence (VAL-001).

## Notes For Developer

- Record the chosen emit precision format in the completion report (PI-5).
- The exact current prices are external facts; pick representative current
  values per provider and cite the source in the completion report. The VAL
  asserts presence + round-trip, not specific prices.
- Keep the change additive: do not remove or rename existing Anthropic entries
  (compatibility — `pricing.toml` is a committed config surface).
