## Summary of what was implemented

- Expanded `DefaultPricingConfig()` to ship defaults for all three Phase 14
  providers: the existing Anthropic set plus OpenAI (`gpt-5.4-mini`,
  `gpt-5.4-nano`) and Google (`gemini-2.5-flash`,
  `gemini-2.5-flash-lite`).
- Centralized the emitted default model order in `internal/config` so
  `pricingConfigContents()` writes every shipped default entry instead of only
  the original Anthropic subset.
- Changed pricing emission from fixed `%.2f` formatting to
  `strconv.FormatFloat(value, 'f', -1, 64)`, which preserves exact float64
  round-trip values for prices like `0.375` and `0.625`.
- Added red/green tests proving three-provider coverage, emit-to-parse
  round-trip equality for the full default set, and bootstrap-file parity.
- Verified with a built `/tmp/fid` binary that `fid init` writes a
  `.fid/pricing.toml` containing all seven default model sections and preserves
  fractional-dollar entries such as `0.375` and `0.625`.
- Pricing sources reviewed on 2026-06-28:
  Anthropic pricing overview: <https://claude.com/pricing>
  OpenAI API pricing: <https://developers.openai.com/api/docs/pricing>
  Google Gemini/Agent Platform pricing:
  <https://cloud.google.com/gemini-enterprise-agent-platform/generative-ai/pricing>
  The new Task 1 defaults were taken from the OpenAI and Google pages above;
  Anthropic defaults were intentionally preserved unchanged per the task note's
  additive-compatibility constraint.

## Files created or modified

- `internal/config/config.go`
- `internal/config/config_test.go`
- `cmd/fid/init_test.go`
- `project-docs/specs/04_TRD_CONFIG.md`
- `project-docs/docs/phases/phase-14/tasks/task-1-done.md`

## Documentation updated

- Updated `project-docs/specs/04_TRD_CONFIG.md` section 3 to show the
  Anthropic/OpenAI/Google default set, document round-trip-safe pricing
  emission, and bump the last-updated date to `2026-06-28`.
- Created `project-docs/docs/phases/phase-14/tasks/task-1-done.md`.

## Validation Contract Coverage

- `VAL-001` - PASS. `internal/config/config_test.go` now proves:
  `DefaultPricingConfig()` includes Anthropic, OpenAI, and Google defaults;
  `pricingConfigContents()` round-trips the full default set through
  `ReadPricingConfig()` with no precision loss; and `WriteBootstrapFiles()`
  writes pricing that matches `DefaultPricingConfig()`.
- `VAL-001` CLI evidence - PASS. After `go build -o /tmp/fid ./cmd/fid`, a temp
  workspace run of `/tmp/fid init` wrote `.fid/pricing.toml` with these default
  sections: `claude-opus-4-7`, `claude-sonnet-4-6`, `claude-haiku-4-5`,
  `gpt-5.4-mini`, `gpt-5.4-nano`, `gemini-2.5-flash`, and
  `gemini-2.5-flash-lite`. The generated file preserved fractional-dollar
  values including `in_per_mtok = 0.375` and `out_per_mtok = 0.625`.
- Chosen emit precision format: `strconv.FormatFloat(value, 'f', -1, 64)`.

## Tests run

- `go test ./internal/config -run 'TestDefaultPricingConfigIncludesAnthropicOpenAIAndGoogleModels|TestPricingConfigContentsRoundTripsDefaultPricingWithoutPrecisionLoss|TestWriteBootstrapFilesRoundTripsDefaultConfig'` - FAIL
  Red step: `DefaultPricingConfig()` was missing `gpt-5.4-mini`, and
  `pricingConfigContents()` emitted only the original Anthropic subset.
- `go test ./internal/config -run 'TestDefaultPricingConfigIncludesAnthropicOpenAIAndGoogleModels|TestPricingConfigContentsRoundTripsDefaultPricingWithoutPrecisionLoss|TestWriteBootstrapFilesRoundTripsDefaultConfig'` - PASS
  Green step after widening defaults, extending the emit order, and changing
  the pricing formatter.
- `gofmt -w internal/config/config.go internal/config/config_test.go cmd/fid/init_test.go` - PASS
- `go test ./cmd/fid -run TestInitCommandCreatesWorkspaceFiles` - PASS

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
  Key package coverage from command output:
  `cmd/fid` `78.1%`
  `internal/config` `80.7%`
  `internal/issues` `81.3%`
  `internal/policy` `84.6%`
- `PASS` `go tool cover -func=coverage.out`
  `internal/config/config.go` coverage details:
  `DefaultPricingConfig` `100.0%`
  `ReadPricingConfig` `82.9%`
  `pricingConfigContents` `90.9%`
  `formatPricingValue` `100.0%`
  overall coverage: `76.1%`
- `PASS` `go build -o /tmp/fid ./cmd/fid`
  Conditional `VAL-001` evidence:
  `/tmp/fid init` in a temp workspace emitted a `.fid/pricing.toml` containing
  all seven default model sections and preserved `0.375` / `0.625` in the file.

## Deviations from spec

- The task file list named only `internal/config/*` plus the config spec, but
  I also updated `cmd/fid/init_test.go` because the required all-package gate
  still hardcoded the pre-Task-1 three-model default count. The change was
  limited to that stale expectation and kept `fid init` coverage aligned with
  Task 1 behavior.

## Concerns, technical debt, or notes for the Tech Lead

- The task note required existing Anthropic defaults to remain unchanged for
  compatibility, so Task 1 did not re-baseline those entries against the
  current Anthropic pricing page reviewed on 2026-06-28. If exact shipped
  Anthropic list prices matter, that needs an explicit follow-up task or
  decision because the committed config surface now mixes preserved Anthropic
  defaults with freshly sourced OpenAI/Google defaults.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit.

## Commit subject

- Pending commit.
