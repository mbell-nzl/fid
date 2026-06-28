## Summary of what was implemented

- Added the new `internal/stats` read-only domain package for telemetry aggregates over `agent_runs`, including grouping by role, model, and owning epic plus `--epic` filtering.
- Implemented `fid stats tokens` and `fid stats cost` in `cmd/fid`, wired them into the root command tree, and reused the existing JSON/TOON/Markdown/YAML output surface.
- Extended markdown formatting so both stats commands render human-readable grouped tables while keeping TOON and YAML structurally aligned with JSON.

## Files created or modified

- `internal/stats/stats.go`
- `internal/stats/stats_test.go`
- `cmd/fid/stats.go`
- `cmd/fid/stats_test.go`
- `cmd/fid/root.go`
- `internal/format/markdown.go`
- `project-docs/specs/05_DESIGN_PLAN.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-14/tasks/task-3-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: updated `fid stats tokens` so the documented `--by` choices include `epic`.
- `project-docs/specs/05_DESIGN_PLAN.md`: reconciled the high-level command tree so the `stats tokens` line matches the implemented `--by epic` surface.
- `project-docs/docs/phases/phase-14/tasks/task-3-done.md`: recorded Task 3 implementation, VAL evidence, and gate results.

## Validation Contract Coverage

- `VAL-007` PASS: `internal/stats` tests cover token aggregation by role, model, and epic plus the owning-epic filter; `cmd/fid` tests prove JSON/TOON/Markdown/YAML parity for `fid stats tokens`.
- `VAL-008` PASS: `internal/stats` tests prove `total_cost_usd` excludes NULL-cost runs while token totals and `run_count` still include them; `cmd/fid` tests prove the same behavior and four-format parity for `fid stats cost`.
- `PI-4` applied: `--by epic` resolves a run on an epic to itself, a run on a task/subtask to its owning epic, and a run with no epic ancestor to a single null/`(none)` bucket.

## Tests run

- `go test ./internal/stats -run TestTokensGroupsByRole`
  Red: failed with `undefined: Tokens` / `undefined: TokensOptions` / `undefined: GroupByRole` / `undefined: TokensRow`.
  Green: passed after adding the initial `internal/stats.Tokens` implementation.
- `go test ./internal/stats`
  Red: failed with `undefined: Cost` / `undefined: CostOptions` / `undefined: CostRow`.
  Green: passed after adding `internal/stats.Cost` and NULL-cost aggregation behavior.
- `go test ./cmd/fid -run TestStatsTokensDefaultsToJSONByRole`
  Red: failed with `unknown command "stats" for "fid"`.
  Green: passed after adding `cmd/fid/stats.go` and registering `newStatsCommand()`.
- `go test ./cmd/fid -run 'TestStatsTokensSupportsFormatsAndEpicFilter|TestStatsCostSupportsFormatsAndExcludesNullCostsFromSum'`
  Red: failed with `unsupported markdown value type []stats.TokensRow` and `unsupported markdown value type []stats.CostRow`.
  Green: passed after extending `internal/format/markdown.go` for stats tables.
- `go test ./internal/stats ./cmd/fid`
  PASS.

## Exact quality-gate results

- PASS `gofmt -l cmd internal`
- PASS `golangci-lint run ./...`
- PASS `go test -race -coverprofile=coverage.out ./...`
- PASS `go tool cover -func=coverage.out`
- PASS `go test -tags integration -race ./...`
- PASS `go build -o /tmp/fid ./cmd/fid`
- Coverage: `internal/stats` `81.6%`, `cmd/fid` `77.4%`
- Total statement coverage from `go tool cover -func=coverage.out`: `75.7%`

## Deviations from spec

- Updated `project-docs/specs/05_DESIGN_PLAN.md` in addition to the task-required `06_DESIGN_PLAN_CLI.md` change so the high-level command tree stays consistent with the implemented public CLI surface.
- Report metadata is committed in a narrow follow-up docs commit because a committed report cannot contain its own final Git object ID until after the implementation commit exists.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending follow-up metadata commit.

## Commit subject

- Pending follow-up metadata commit.
