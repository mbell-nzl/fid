# Phase 8 - Task 2: Gate snapshot (`gate_runs` → `task_completions.gate_results`)

## Objective

Implement the gate-snapshot mechanism on `fid issue complete`: at completion
time, snapshot the `gate_runs` rows matching the completion's `(issue_id,
attempt_number)` into `task_completions.gate_results` as a JSON array of
`{cmd, exit_code, ran_at, output_excerpt}` and derive `gates_passed` (exit_code
0) / `gates_failed` (exit_code ≠ 0). Source `gate_runs` rows are preserved as
audit rows. Because `fid dev gate` (the `gate_runs` writer) is v0.1e, exercise
the mechanism by seeding `gate_runs` through a journal fixture.

## Requirements

- PRD refs: PRD-DEV-004 (gate results as structured data on completions),
  PRD-DEV-010 (gate runs append-only and rolled into the completion for the same
  issue and attempt)
- TRD refs: TRD_SCHEMA §2 `gate_runs` + `task_completions.gate_results`
  (incl. the schema comment: "snapshot the gate_runs for the same issue and
  attempt into task_completions.gate_results and derive gates_passed/gates_failed")
- Design refs: `06_DESIGN_PLAN_CLI.md` §3 `fid dev gate` (snapshot note)
- Quality refs: `08_QUALITY_PLAN.md` (TDD, DEC-007 coverage)

## Validation Contract Coverage

- VAL-005: Completing an attempt with matching `gate_runs` snapshots them into
  `gate_results` with derived `gates_passed`/`gates_failed`, preserves source
  rows; a no-gate completion yields `gate_results = []` and 0/0.

## Scope

### Files To Create Or Modify

- `internal/issues/complete.go` - read matching `gate_runs` under the lock,
  build the `gate_results` JSON + counts, and include them in the completion
  insert effect values.
- `internal/issues/complete_test.go` - add gate-snapshot cases.
- A journal-fixture test helper for seeding `gate_runs` (mirror the Phase 5/7
  fixture pattern) if one does not already exist; keep it test-only.

### Files To Read For Context

- `internal/issues/detail.go` `issueCompletions` (how `gate_results` JSON is
  decoded for read-back — keep shapes consistent).
- `internal/journal/project_effects.go` `applyInsert` (generic insert path the
  fixture uses to seed `gate_runs`).
- Phase 5 task that seeded gate/completion fixtures (`internal/issues` tests
  referencing `task_completions`/`gate_runs`) for the fixture convention.
- TRD_SCHEMA §2 comment block under `gate_runs`.

### Out Of Scope

- The real `fid dev gate` writer (v0.1e); seed `gate_runs` via fixture only.
- `has_completion_report` wiring and rebuild parity: Task 3.
- Changing attempt numbering / supersession / core fields from Task 1.
- `internal/issues/issue.go` size split: Task 4.

## Acceptance Criteria

- [ ] With `gate_runs` seeded for `(issue, attempt)` containing one exit_code 0
  and one non-zero row, completing the attempt snapshots both into
  `gate_results` as `[{cmd, exit_code, ran_at, output_excerpt}]` and derives
  `gates_passed = 1`, `gates_failed = 1` (VAL-005).
- [ ] The source `gate_runs` rows are unchanged after the snapshot (still present,
  unmodified) (VAL-005).
- [ ] Completing an attempt with **no** matching `gate_runs` yields
  `gate_results = []` and `gates_passed = gates_failed = 0` (VAL-005).
- [ ] Only `gate_runs` whose `attempt_number` matches the completion's derived
  attempt are snapshotted (rows for a different attempt are excluded) (VAL-005).
- [ ] `fid issue get --format json` reads back the `gate_results` array and the
  derived counts consistent with the stored row.

## TDD Plan

- First failing test: seed pass+fail `gate_runs` for attempt 1, complete, assert
  `gate_results` array contents and `gates_passed`/`gates_failed`.
- Additional tests: no-gate completion → empty/0/0; gates from a non-matching
  attempt excluded; source `gate_runs` preserved after snapshot.
- Regression risks: keep the `gate_results` JSON field names aligned with the
  read-back decoder in `detail.go`; ensure the snapshot read runs against the
  fresh in-lock projection so seeded fixtures are visible.

## Documentation Updates

- `project-docs/docs/phases/phase-8/` completion report - paste a transcript /
  test output showing the snapshot array and derived counts, and note that the
  `gate_runs` source was fixture-seeded because `fid dev gate` is v0.1e (DEC-C).
- `06_DESIGN_PLAN_CLI.md` / TRD_SCHEMA - "verified, no change" unless the
  implemented snapshot shape diverges from the spec comment (spec-sync input).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: `gate_results` shape must match the `detail.go`
  read-back decoder and the schema comment.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Coverage targets (DEC-007): `internal/issues` ≥ 80%.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` only if a gate-snapshot test is placed
  behind the `integration` build tag.
- `go build -o /tmp/fid ./cmd/fid` if pasting CLI transcripts.

## Notes For Developer

- The gate writer does not exist yet (v0.1e); the production path of
  `fid issue complete` will simply find no `gate_runs` and snapshot an empty set.
  Prove the mechanism by seeding `gate_runs` via a journal fixture (DEC-C),
  mirroring the Phase 5/7 fixture approach.
- Match on `(issue_id, attempt_number)` only. Derive counts from `exit_code`
  (0 = passed, non-zero = failed); do not infer from any other field.
- Keep the snapshot a pure read of `gate_runs` plus values folded into the same
  completion insert effect — do not emit update effects against `gate_runs`.
