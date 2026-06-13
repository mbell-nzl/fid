# Phase 4 - Task 6: internal/issues Coverage Closeout (COV-001)

## Objective

Restore the COV-001 deterministic gate for Phase 4 by raising `internal/issues`
line coverage from `79.5%` to at least `80%`. The phase final review
(`project-docs/docs/phases/phase-4/qa/final-review.md`) records a single
blocking finding: COV-001 fails for `internal/issues`. Every behavioral
assertion (VAL-001..VAL-009) already passed; this task adds package-level tests
for the remaining uncovered branches in already-shipped Phase 4 code. It adds no
new product behavior and changes no public interface, schema, or journal
contract.

## Requirements

- PRD refs: N/A (no new product behavior)
- TRD refs: TRD-ARCH-005 (adapter boundary unchanged)
- Design refs: N/A
- Quality refs: `08_QUALITY_PLAN.md` section 4.3 COV-001 (domain package >= 80%);
  DEC-007 COV-001 carry-in

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-4/qa/validation-contract.md`.

- No new `VAL-NNN` is advanced. This is a coverage-closeout task: the COV-001
  deterministic gate is a precondition for the phase final-review PASS that the
  already-passing VAL-001..VAL-009 evidence depends on. New tests strengthen,
  but do not replace, existing evidence for the `internal/issues` behaviors in
  VAL-003, VAL-006, VAL-007, VAL-008, and VAL-009.

## Scope

### Files To Create Or Modify

- `internal/issues/*_test.go` (add/extend package-level tests only)
- `project-docs/docs/phases/phase-4/tasks/task-5-completion.md` (correct the
  inaccurate COV claim flagged as the final-review Improvement finding)
- `project-docs/docs/phases/phase-4/tasks/task-6-completion.md`

### Files To Read For Context

- `internal/issues/list.go`
- `internal/issues/detail.go`
- `internal/issues/issue.go`
- existing `internal/issues/*_test.go`

### Out Of Scope

- Any change to non-test product code in `internal/issues`, `cmd/fid`, or
  migrations. If reaching `80%` appears to require a product-code change, stop
  and escalate to the Tech Lead instead of altering behavior.
- New CLI surfaces, schema, or journal events.
- Lowering or reinterpreting the COV-001 threshold.

## Acceptance Criteria

- [ ] `internal/issues` line coverage is `>= 80%` reported by
  `go tool cover -func=coverage.out`. (COV-001)
- [ ] No `internal/issues` non-test source file is modified; the diff is
  test-only plus the two documentation files. (no VAL)
- [ ] New tests assert real observable behavior of the previously uncovered
  branches (e.g. `ParsePriorityFilter` empty/non-integer rejection,
  `issueListOrderClause` unknown-order error, `SectionValue` unknown-section
  error, `normalizeDateOnlyPointer` nil/empty/valid/invalid date cases,
  `issueParentTitle` no-row case) rather than calling code for coverage alone.
  (no VAL)
- [ ] `cmd/fid` (COV-002 >= 70%) and the other domain packages remain at or
  above their thresholds after the change. (COV-001, COV-002)

## TDD Plan

- First step: run `go test -race -coverprofile=coverage.out ./internal/issues`
  then `go tool cover -func=coverage.out` (and optionally
  `go tool cover -html=coverage.out -o /tmp/issues-cov.html`) to enumerate the
  exact uncovered lines. Target the highest-value uncovered branches first.
- First failing test: a table test over `ParsePriorityFilter` covering valid
  list, empty-string input, embedded empty value, and non-integer value.
- Additional tests: `issueListOrderClause` unknown-order returns
  `*UnsupportedListOrderError` naming supported orders; `SectionValue` returns
  `*UnsupportedSectionError` for an unknown section and each supported section
  returns its slice; `normalizeDateOnlyPointer` for nil, empty, full-RFC3339,
  date-only, and malformed input; any remaining uncovered scan/error helpers
  reachable without faking the driver.
- Regression risks: tests that assert on error strings drifting from the
  implementation; flaky ordering assertions; accidentally editing product code
  to satisfy coverage.

## Documentation Updates

- `project-docs/docs/phases/phase-4/tasks/task-5-completion.md` - correct the
  VAL-009 evidence sentence that claims `internal/issues 79.5%` "satisfies the
  task's COV references"; it does not satisfy COV-001. Note the correction was
  made under Task 6 coverage closeout. (Final-review Improvement finding.)
- `project-docs/docs/phases/phase-4/tasks/task-6-completion.md` - record exact
  before/after `internal/issues` coverage, the full gate results, and the
  test-only diff scope.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None; do not read held-out scenarios.
- Cross-surface parity impact: None; no output shapes change.
- Audit or release-readiness impact: Restores the COV-001 component of the
  Phase 4 deterministic gate so QA Lead can rerun the full gate to PASS.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each as `PASS`,
or `BLOCKED` with reason in a blocker report. Confirm in the completion report
that `internal/issues >= 80%` (COV-001) and `cmd/fid >= 70%` (COV-002).

Conditional commands:

```text
git diff --stat
```

List any conditional command required by this task and why:

- `git diff --stat` is required to demonstrate the diff is test-only plus the
  two named documentation files (acceptance criterion 2).

## Notes For Developer

- Coverage is the goal, but every new test must assert observable behavior. Do
  not add no-op tests that call functions only to mark lines covered.
- If `80%` cannot be reached without editing non-test product code, stop and
  escalate to the Tech Lead. Do not change behavior, the threshold, or the
  validation contract.
- Do not touch `cmd/fid` product code; this is an `internal/issues` test gap.
