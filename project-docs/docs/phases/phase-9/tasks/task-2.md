# Phase 9 - Task 2: `fid concern address` + `concerns_addressed` rule proof

## Objective

Add the `fid concern address` write command: transition a concern in place from
`open` to `addressed`, recording the addressing completion, with
validate-before-append rejection. Prove the derived flag clears and that the
existing `concerns_addressed` transition rule reports correctly against real
concern data (DEC-A: at the rule-evaluation seam; no default workflow edge gates
on it).

## Requirements

- PRD refs: PRD-DEV-006 (concern lifecycle to `addressed`), PRD-DEV-007 (derived
  concern flag clears), PRD-WF-009 (`concerns_addressed` rule).
- TRD refs: §5.1 `concerns` schema (`status`, `addressed_in_completion_id`);
  §4.3 write-path seam; §4.6 default workflow graphs (no `concerns_addressed`
  edge — DEC-A).
- Design refs: `06_DESIGN_PLAN_CLI.md` §6 `fid concern address`.
- Quality refs: `08_QUALITY_PLAN.md` (TDD, coverage gate; DEC-007 COV-001/002).

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-9/qa/validation-contract.md`.

- VAL-004: `fid concern address` updates the concern `open`→`addressed`, sets
  `addressed_in_completion_id` + `updated_at`, retains the row.
- VAL-005: address rejected before any append for unknown role, nonexistent
  concern, concern not `open`, nonexistent `--in-completion`; valid address
  applies exactly one status mutation (CF-007 / DEC-C/DEC-D).
- VAL-006: addressing the last open concern clears `concerns_raised` to false and
  `open_concerns_count` to 0; a remaining open concern keeps the flag true.
- VAL-007: the `concerns_addressed` rule reports unsatisfied (naming open concern
  IDs) while an open concern exists and satisfied once all are addressed, against
  real concern data via the rule-evaluation seam (DEC-A).

## Scope

### Files To Create Or Modify

- `internal/issues/concern_address.go` (new) - `AddressConcern` domain function:
  `RequireKnownRole`, load the concern, reject when not `open`, validate
  `--in-completion` exists (and belongs to the concern's issue), build an
  `update` effect (`status='addressed'`, `addressed_in_completion_id`,
  `updated_at`) under `.fid/write.lock`, dry-run, append, project.
- `internal/issues/concern_address_test.go` (new).
- `internal/issues/transition_rules_test.go` (new or extend) - integration test
  exercising `evaluateTransitionRule(RuleConcernsAddressed)` /`loadConcernStatus`
  against real raise/address writes for VAL-007.
- `cmd/fid/concern.go` - add the `address` subcommand (`--in-completion`, `--by`)
  + design §6 stderr confirmation.
- `cmd/fid/concern_test.go` - extend with address command tests.

### Files To Read For Context

- `internal/issues/concern_raise.go` (Task 1 — mirror its structure for update).
- `internal/issues/complete.go` (the prior-attempt `update` effect pattern).
- `internal/issues/transition_rules.go` (`evaluateTransitionRule`,
  `loadConcernStatus`).
- `internal/workflow/rules.go` (`ConcernsAddressed`, `UnmetCondition`).
- `internal/issues/detail.go` (`issueConcerns` read-back of `status`).

### Out Of Scope

- `escalated` / `wontfix` transitions (deferred v0.1e+; DEC-D).
- Any new `fid issue transition` edge gated on `concerns_addressed` (DEC-A:
  out of scope; would change the TRD §4.6 default graph).
- `fid concern list` top-level command (Task 3).
- `activity_log` rows (DEC-E: none).

## Acceptance Criteria

- [ ] `fid concern address <CON> --in-completion <TC> --by <role>` updates the
  concern `open`→`addressed`, sets `addressed_in_completion_id` + `updated_at`,
  retains the row; read-back shows `status='addressed'` (VAL-004)
- [ ] Rejected before append (non-zero exit, clear error, zero new lines/rows)
  for: unknown role; nonexistent concern; concern not `open` (e.g. already
  addressed); nonexistent `--in-completion` (VAL-005)
- [ ] Addressing the last open concern clears `concerns_raised`→false /
  `open_concerns_count`→0; a remaining open concern keeps it true (VAL-006)
- [ ] The `concerns_addressed` rule reports unsatisfied naming open concern IDs
  while open and satisfied once all addressed, via the rule seam with real data
  (VAL-007)

## TDD Plan

- First failing test: `AddressConcern` updates `status` to `addressed` and sets
  `addressed_in_completion_id`, read back from the DB.
- Additional tests: each reject path asserts zero new journal lines + error
  (esp. re-addressing an already-addressed concern); partial vs full derived-flag
  clear (VAL-006); `evaluateTransitionRule(RuleConcernsAddressed)` returns
  unsatisfied+named IDs while open and satisfied after address (VAL-007);
  command-level address tests.
- Regression risks: in-place update vs accidental insert (assert row count
  unchanged); journal-poisoning on reject paths (CF-007).

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` §6 `fid concern address` - confirm flags/behavior.
- Completion report - record DEC-A (rule proven at the seam, no edge) and DEC-D
  (raise→address only) as implemented; note any deviation.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: confirm the derived flag clears identically on
  `fid issue get` and `fid issue list`.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Coverage must show
`internal/issues` ≥ 80% and `cmd/fid` ≥ 70% (DEC-007); record percentages in the
completion report.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

- Run `go build -o /tmp/fid ./cmd/fid` (command surface change).
- Run the integration suite if any test is added under the `integration` tag.

## Notes For Developer

- The address effect is an `update` on the existing `concerns` row keyed by `id`
  (mirror the `superseded` update in `complete.go`), not an insert.
- VAL-007 is proven against the **rule-evaluation seam** (`evaluateTransitionRule`
  / `loadConcernStatus`), NOT via `fid issue transition` — no default workflow
  edge gates on `concerns_addressed` (TRD §4.6). Do not add such an edge.
- `loadConcernStatus` already returns open concern IDs sorted; assert the
  `UnmetCondition.EntityIDs` carry the open `CON-` IDs.
- Keep new functions ≤60 lines and new files ≤400 lines.
