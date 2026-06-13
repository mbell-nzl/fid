# Phase 5 - Task 1: Validate-before-append write-path seam (CF-007)

## Objective

Close CF-007 / DEC-014: a write that will be rejected by a projection-time
trigger or constraint must be rejected **before** the journal append, so the
append remains the commit point only for events that project cleanly. Today
`internal/issues.Create` appends the event and only then projects it, so a
trigger-rejected create still commits to the canonical journal and poisons the
tail — every later read replays the bad line and fails.

Introduce a reusable pre-append projection-dry-run seam in `internal/journal`
(project the built event into a transaction that is always rolled back; surface
the same error the real projection would raise) and route `internal/issues`
writes through it so validation precedes the append. This seam is also the
foundation Tasks 3 and 5 build their `update`/`transition` write paths on, and
it is the same validate/build-before-commit separation TRD §4.3 calls out for
the deferred `--dry-run` flag.

## Requirements

- PRD refs: PRD-OPS-001, PRD-OPS-002 (journal append is the commit point; one
  logical write = one clean event)
- TRD refs: TRD §4.3 (write path steps 3–6; dry-run simulation seam),
  TRD-ARCH-001
- Design refs: N/A
- Quality refs: `08_QUALITY_PLAN.md` §4.3 COV-001 (`internal/issues`,
  `internal/journal` ≥ 80%)

## Validation Contract Coverage

- VAL-005: A rejected write appends zero journal lines and later reads still
  succeed; an accepted write appends exactly one line and projects cleanly.

## Scope

### Files To Create Or Modify

- `internal/journal/project.go` - add an exported pre-append dry-run projection
  helper that applies an event's effects inside a transaction and always rolls
  back, returning the projection error (trigger/constraint/precondition) without
  committing.
- `internal/journal/project_test.go` - dry-run helper unit tests.
- `internal/issues/issue.go` - run the dry-run projection (or equivalent
  pre-append validation) before `AppendEvent` in `Create`; only append when the
  event will project cleanly.
- `internal/issues/issue_test.go` - regression tests (see TDD plan).
- Any small shared helper file under `internal/journal/` if extraction is
  cleaner than inlining.

### Files To Read For Context

- `internal/journal/project.go`, `internal/journal/writer.go`
- `internal/issues/issue.go` (`Create`)
- `project-docs/docs/phases/phase-4/decisions.md` (DEC-014 diagnosis)
- `project-docs/specs/02_TRD.md` §4.3

### Out Of Scope

- `fid issue update` / `fid issue transition` (Tasks 3, 5).
- The `activity` effect op (Task 2).
- Any change to journal file format, event schema, or projection_meta shape.

## Acceptance Criteria

- [ ] A create whose projection would fail (e.g. task with `--sequence` under a
  prioritised epic, the Phase 4 repro) is rejected with the existing actionable
  error and appends **zero** lines to the journal. (VAL-005)
- [ ] After such a rejected create, an unrelated valid `fid issue get` (and a
  rebuild) still succeeds — no poisoned tail. (VAL-005)
- [ ] An accepted create still appends exactly one line and projects cleanly;
  behaviour and output are unchanged from Phase 4. (VAL-005)
- [ ] The dry-run projection helper never commits and never leaves
  `projection_meta` advanced. (VAL-005)
- [ ] COV-001 holds for `internal/issues` and `internal/journal` (≥ 80%).

## TDD Plan

- First failing test: in `internal/issues`, assert that after a rejected create
  (projection-failing event) the journal file line count is unchanged and a
  subsequent `Get` of a previously-valid issue succeeds. This fails today
  because `Create` appends before projecting.
- Additional tests: dry-run helper rolls back on success and on failure and
  does not mutate `projection_meta`; accepted create still appends exactly one
  line; poisoned-tail recovery (a workspace that already has a bad tail is out
  of scope to auto-repair, but a freshly rejected create must not create one).
- Regression risks: double-projection (dry-run then real) must be transaction
  isolated; ensure the real projection still runs once after append; ensure the
  write lock still wraps the whole sequence; ensure error types/messages
  surfaced to the CLI are unchanged.

## Documentation Updates

- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` - if the create flow narrative
  needs the validate-before-append ordering made explicit, update it; otherwise
  record `verified-no-change` rationale in the completion report.
- Completion report: cite VAL-005 and the CF-007 closure evidence.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: None (no new CLI surface).
- Audit or release-readiness impact: Closes CF-007 / DEC-014; note the
  carry-forward closure in the completion report for the Tech Lead wrap.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. A focused test command
may be used during the red/green loop but does not substitute for the full
gate. Record each command as `PASS` (or `BLOCKED` with reason) in the
completion report.

Conditional commands:

- None.

## Notes For Developer

- Keep the append as the commit point: the dry-run is a pre-check, the real
  projection after append is unchanged. Do not skip the post-append projection.
- The dry-run must project into the **current** DB state (after
  `EnsureProjectionFresh`) so trigger/precondition outcomes match the real
  projection.
- Do not weaken or remove the Phase 4 sequence-vs-scheduling backstop trigger;
  this task makes the app layer reject before append, with the trigger remaining
  the backstop for non-CLI writes.
