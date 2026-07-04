# Phase 16 - Task 6: CF-016 - `fid issue get --section gates`

> **Confirmed task.** The human confirmed DC-B at kickoff on 2026-07-05
> (workspace timezone), so CF-016 is scheduled in Phase 16.

## Objective

Add a first-class issue-scoped read surface for recorded gates (CF-016): expose
the append-only `gate_runs` history through `fid issue get <id> --section gates`,
fitting the existing issue-section read model and keeping the raw history
distinct from the completion snapshot (`task_completions.gate_results`).

## Requirements

- PRD refs: PRD-DEV-004, PRD-DEV-010 (gate runs); PRD-ISSUE-006/018 (issue get
  `--section`)
- TRD refs: `03_TRD_SCHEMA.md` §7 (`gate_runs` columns)
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid issue get --section` (add `gates`)
- Quality refs: `08_QUALITY_PLAN.md` (read freshness, size budgets)

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-16/qa/validation-contract.md`.

- VAL-012: `fid issue get <id> --section gates` returns the `gate_runs` history
  for in-progress and finished tasks (empty when none), all formats, distinct
  from the completion snapshot.

## Scope

### Files To Create Or Modify

- `internal/issues/detail.go` - add `gates` to `supportedIssueSections`; add an
  `IssueGate` struct (`id`, `attempt_number`, `agent_run_id`, `recorded_by_role`,
  `command`, `exit_code`, `output_excerpt`, `ran_at`); add `issueGates(conn,
  issueID)` querying `gate_runs WHERE issue_id = ? ORDER BY attempt_number,
  ran_at, id`; wire `gates` into `SectionValue` and the `Issue` detail assembly.
- `internal/issues/detail_test.go` - gates section: history for an issue with
  gates (ordered), empty list when none, fields populated incl. nullable
  `output_excerpt`/`agent_run_id`.
- `cmd/fid/issue_section_test.go` - CLI `--section gates` in all four formats +
  empty case.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - `fid issue get --section`: add
  `gates` to the `--section` enum and document the row shape (mirror the
  `--section concerns` style; note it is the raw `gate_runs` history, distinct
  from `task_completions.gate_results`).

### Files To Read For Context

- `internal/issues/detail.go` (existing `--section` model: `SupportedSections`,
  `SectionValue`, per-section query helpers like `issueConcerns`/`issueClaims`)
- `internal/issues/get.go` / `internal/issues/issue.go` (the `Issue` detail
  struct + assembly where sections are populated)
- `cmd/fid/issue_section_test.go`, `cmd/fid/issue_get_list.go` (CLI `--section`
  wiring)
- `03_TRD_SCHEMA.md` §7 `gate_runs`; `internal/dev/gate.go` (gate row writer)

### Out Of Scope

- A new `fid dev gate list` verb (rejected in favor of `--section gates`).
- `fid stats gates` / `fid run show` (v0.2c).
- Any write path or permissions row (read-only, PI-1).

## Acceptance Criteria

- [ ] `fid issue get <id> --section gates` returns the issue's `gate_runs`
  history (`command`, `exit_code`, `output_excerpt`, `ran_at`, `attempt_number`,
  `recorded_by_role`) ordered by attempt then time, for both in-progress and
  finished tasks, in all four formats. (VAL-012)
- [ ] An issue with no recorded gates returns an empty list (not null/error).
  (VAL-012)
- [ ] The section reads `gate_runs` directly and is distinct from the completion
  snapshot `task_completions.gate_results`. (VAL-012)
- [ ] `06_DESIGN_PLAN_CLI.md` lists `gates` in the `--section` enum with its row
  shape. (doc)

## TDD Plan

- First failing test: `SectionValue(issue, "gates")` returns the gate history —
  fails red (section unsupported today).
- Additional tests: ordered multi-attempt history; empty list when none; nullable
  `output_excerpt`/`agent_run_id`; CLI all-format render; `--section gates` after
  `dev gate` (in-progress) and after `dev finish` (finished) both return history.
- Regression risks: do not change the other sections' behavior or the
  `UnsupportedSectionError` message format beyond adding `gates`; keep the
  raw-history vs snapshot distinction (read `gate_runs`, not
  `task_completions`); functions ≤60 lines.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - add `gates` to `fid issue get
  --section` enum + row shape.
- Completion report: note CF-016 closed (or its disposition).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: `gates` section row shape vs the `gate_runs`
  schema (§7) and vs the completion snapshot (must stay distinct).
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
`internal/issues` and `cmd/fid` coverage percentages.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` — required if any new test is
  integration-tagged (package set includes `cmd/fid`); record PASS/BLOCKED.
- `go build -o /tmp/fid ./cmd/fid` — required: CLI build used by the section test.

## Notes For Developer

- DC-B is human-confirmed for Phase 16; do not reopen it unless implementation
  reveals a new product or public-interface contradiction.
- Mirror an existing section helper (`issueConcerns`/`issueClaims`) for the query
  + assembly so the read shape and freshness contract match.
- Keep the raw `gate_runs` history distinct from `task_completions.gate_results`
  (the completion snapshot) — read `gate_runs` directly.
- This is a read command — no `--by`, no matrix row (PI-1).
