# Phase 14 - Task 3: `fid stats tokens` and `fid stats cost`

## Objective

Implement the read-only telemetry aggregates `fid stats tokens` and
`fid stats cost` over `agent_runs`, grouped by role, by model, and by epic, with
all four read formats. This completes the Phase 14 telemetry surface: record a
run (Task 2), then aggregate cost and tokens.

## Requirements

- PRD refs: PRD-TELEM-004, PRD-OUT-001, PRD-OUT-002, PRD-OUT-003, PRD-OUT-004
- TRD refs: `02_TRD.md` section 3 (`internal/stats` aggregations)
- Design refs: `06_DESIGN_PLAN_CLI.md` section 13 (`fid stats tokens`/`cost`);
  `05_DESIGN_PLAN.md` section 6.4 (`fid stats cost` JSON)
- Quality refs: `08_QUALITY_PLAN.md` (size budgets)

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-14/qa/validation-contract.md`.

- VAL-007: `fid stats tokens [--by role|model|epic] [--epic <id>]` aggregates
  tokens + run_count by dimension across all four formats.
- VAL-008: `fid stats cost [--by role|model|epic] [--epic <id>]` aggregates
  total_cost_usd (NULL-cost runs excluded from the sum but counted in
  run_count/tokens) by dimension across all four formats.

## Scope

### Files To Create Or Modify

- `internal/stats/stats.go` (new package) - `Tokens(options)` and
  `Cost(options)`: `EnsureProjectionFresh`, open db, run the grouped
  aggregation queries over `agent_runs` (join to issues for `--by epic` and the
  `--epic` filter, PI-4), return typed result slices with `json` tags.
- `internal/stats/stats_test.go` (new) - aggregation correctness per dimension,
  `--epic` filter, and NULL-cost handling.
- `cmd/fid/stats.go` (new) - `newStatsCommand()` with `stats tokens` and
  `stats cost` subcommands; `--by`, `--epic`, `--format` flags; default
  `--by role`; `fidformat.Write` dispatch; `renderIssueError` on failure.
- `cmd/fid/stats_test.go` (new) - CLI tests incl. four-format parity.
- `cmd/fid/root.go` - register `newStatsCommand()`.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - `fid stats tokens`: add `epic`
  to the `--by` choices so both stats commands satisfy PRD-TELEM-004 (PI-4).

### Files To Read For Context

- `internal/claims/list.go` (read pattern: EnsureProjectionFresh → open →
  query → typed items)
- `cmd/fid/claim.go` (read command + `--format` dispatch pattern)
- `internal/format/` (format writer; confirm four-format support)
- `internal/issues/` issue-hierarchy helpers for resolving an issue's epic
  ancestor (for `--by epic` / `--epic`)

### Out Of Scope

- `fid run record` and the write path (Task 2).
- Stats expansion: `fid stats tasks`/`attempts`/`gates`/`interventions`
  (all v0.2c).
- Provider-billed cost reconciliation reporting.

## Acceptance Criteria

- [ ] `fid stats tokens [--by role|model|epic] [--epic <id>]` returns
  `total_tokens_in`, `total_tokens_out`, `run_count` grouped by the chosen
  dimension; group totals match a hand-computed reference; `--epic` filters
  correctly. (VAL-007)
- [ ] `fid stats cost [--by role|model|epic] [--epic <id>]` returns
  `total_cost_usd` plus token totals and `run_count`; NULL-cost runs are
  excluded from the cost sum but counted in `run_count`/tokens. (VAL-008)
- [ ] All four formats (json, toon, markdown, yaml) carry the same data for both
  commands. (VAL-007, VAL-008)
- [ ] `06_DESIGN_PLAN_CLI.md` `fid stats tokens` `--by` includes `epic`. (doc)

## TDD Plan

- First failing test: `internal/stats` `Tokens` grouped by role over seeded
  `agent_runs` returns expected per-role totals — fails red (package absent).
- Additional tests: `--by model`, `--by epic` (incl. the null/`(none)` bucket),
  `--epic` filter; `Cost` NULL-cost exclusion from the sum while counted in
  run_count/tokens; CLI four-format parity-vs-JSON for both commands.
- Regression risks: ensure deterministic ordering of grouped rows for stable
  format output; keep files under 400 lines and functions under 60.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - `fid stats tokens` `--by`
  choices include `epic` (PI-4).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: four-format parity for both stats commands.
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
`internal/stats` and `cmd/fid` coverage percentages.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` is required whenever the package set
  under test contains integration-tagged tests, which it does (`cmd/fid`
  already has integration-tagged tests). Run it and record `PASS` or `BLOCKED`
  even when this task adds no new integration-tagged tests; "not run" is not a
  valid disposition for a required conditional command. (See ESC-001.)
- `go build -o /tmp/fid ./cmd/fid` for the CLI build used by command tests.

## Notes For Developer

- PI-4: `--by epic` groups each run under its issue's owning epic (run on an
  epic → itself; issue with no epic ancestor → a single null/`(none)` bucket).
  NULL-cost runs count in `run_count`/tokens but not in `total_cost_usd`.
- Reuse the existing read/format pattern; do not add a new formatting surface.
- Stats are read-only: no journal writes, no lock acquisition beyond the
  freshness catch-up that `EnsureProjectionFresh` performs.
