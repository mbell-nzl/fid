# Phase 1 - Task 8: Domain-Package Unit Coverage To Meet COV-001

## Objective

Close the Phase 1 final-review blocking finding (QA-11 / Quality Plan COV-001).
The domain packages `internal/issues/`, `internal/policy/`, and
`internal/journal/` are below the 80% package-level coverage threshold:
`internal/issues` 0.0%, `internal/policy` 0.0%, `internal/journal` 68.3%
(`go tool cover -func=coverage.out`, final review). The behaviors are largely
exercised through `cmd/fid` integration tests, but Go attributes coverage to a
package's own test binary, so these packages have **no `_test.go` files of their
own** and several real domain branches are unproven at any layer (notably
`internal/policy` corrupt-tail tolerance — see TRG-002 — and `internal/issues`
not-found / input-guard paths).

Add **package-level unit tests** that exercise the untested domain branches and
bring each of the three packages to or above the COV-001 80% threshold. This is
a test-only task: do not change production behavior.

## Requirements

- PRD refs: PRD-OPS-001, PRD-OPS-012, PRD-ROLE-001, PRD-ROLE-002
- TRD refs: TRD-ARCH-001..003 (replay/freshness), TRD §4.3 (events, replay,
  corruption, preconditions, checkpoint)
- Design refs: N/A
- Quality refs: Quality Plan §4.3 COV-001 (80% across `internal/issues/`,
  `internal/policy/`, `internal/journal/`), COV-003; QA plan QA-11
- Triage refs: TRG-002 (role-validation corrupt-tail tolerance)

## Validation Contract Coverage

This task introduces no new `VAL-NNN` and changes no asserted behavior. It
strengthens the evidence base for already-PASS assertions whose domain logic
lives in these packages: VAL-006, VAL-008, VAL-009, VAL-010, VAL-011, VAL-012,
VAL-016. The binding driver for this task is **Quality Plan COV-001**, a quality
requirement rather than a `VAL-NNN`.

## Scope

### Files To Create Or Modify

- `internal/issues/issue_test.go` (new; package-level unit tests)
- `internal/policy/policy_test.go` (new; package-level unit tests)
- `internal/journal/*_test.go` (extend existing tests for uncovered branches)
- `project-docs/docs/phases/phase-1/tasks/task-8-completion.md`

### Files To Read For Context

- `internal/issues/issue.go`
- `internal/policy/policy.go`
- `internal/journal/replay.go`, `internal/journal/project.go`,
  `internal/journal/upgrade.go`, and existing `internal/journal/*_test.go`
- `cmd/fid/issue_test.go` (to avoid duplicating CLI-level coverage)
- `project-docs/specs/08_QUALITY_PLAN.md` §4.3 (COV-001, COV-003)
- `project-docs/docs/phases/phase-1/qa/final-review.md` (QA-11 finding)
- `project-docs/docs/phases/phase-1/triage-log.md` (TRG-002)

### Out Of Scope

- Any change to production behavior in `internal/`, `cmd/fid/`, schema, or
  migrations. This is test-only.
- Adding TOON/markdown/YAML formats, new kinds, or any new feature.
- Refactoring production code to make it testable, unless a tiny, behavior-
  preserving seam is unavoidable — if so, record it as a deviation in the
  completion report and keep it minimal.
- Reaching for 100%; the target is the COV-001 threshold proven by meaningful
  assertions, not coverage padding.

## Acceptance Criteria

- [ ] `internal/issues` package-level unit tests cover at minimum: `Get`
  not-found returns `*NotFoundError`; `Create`/`Get` input guards (empty
  workspace root, empty migrations dir, empty title, empty id); unsupported
  kind returns `*UnsupportedKindError`; `statusCategory` mapping for each
  category (unstarted/started/completed/cancelled/default); and the
  populated-value paths of `nullStringPointer`/`nullIntPointer`. (COV-001)
- [ ] `internal/policy` package-level unit tests cover at minimum:
  `RequireKnownRole` missing role (`*MissingRoleError`), unknown role
  (`*UnknownRoleError` naming the role), and accepted seeded role;
  `KnownRoles` empty-workspace-root error and missing-journal-dir error;
  role-event `insert`/`update`/`delete` and `active=false` handling; and the
  corrupt-tail tolerance branch (a malformed trailing line after at least one
  valid role event is tolerated, while a malformed line with no prior roles
  surfaces an error). This branch is the TRG-002 behavior and must be asserted
  in both directions. (COV-001)
- [ ] `internal/journal` tests are extended to cover the remaining uncovered
  branches identified from the coverage report (e.g. `ReplayJournal` nil
  connection; `orderedJournalFiles` skipping non-`.jsonl` / non-date files;
  `EnsureProjectionFresh` catch-up vs rebuild branches not already covered).
  (COV-001)
- [ ] `go tool cover -func=coverage.out` reports **≥ 80%** for each of
  `internal/issues`, `internal/policy`, and `internal/journal`. Record the
  per-package percentages in the completion report. (COV-001)
- [ ] `cmd/fid` coverage stays ≥ 70% (COV-002) and no existing test regresses.
  (no VAL; regression guard)
- [ ] Production code behavior is unchanged. If a new test reveals a genuine
  defect, stop and report it as a finding in the completion report (and raise
  an escalation) rather than altering behavior to make a test pass. (no VAL)

## TDD Plan

- First failing test: add `internal/issues/issue_test.go` with a `Get` against
  a freshly initialized workspace for a non-existent id, asserting the returned
  error is `*NotFoundError`. It fails until the test file exists and exercises
  the path (the production code already returns it).
- Additional tests: the branch sets listed in Acceptance Criteria for all three
  packages. Use small fixtures/temp workspaces; for `internal/policy`
  corrupt-tail cases, write a journal file with a valid role event followed by a
  malformed trailing line, and a separate file whose first line is malformed.
- Use the coverage report to find residual gaps: run
  `go test -coverprofile=coverage.out ./internal/issues/ ./internal/policy/ ./internal/journal/`
  then `go tool cover -func=coverage.out` and target the lines still at 0.
- Regression risks: tests that assert on time-dependent fields or absolute
  paths; do not assert behavior the production code does not actually have.
  Coverage must come from real assertions, not bare execution.

## Documentation Updates

- `project-docs/docs/phases/phase-1/tasks/task-8-completion.md` - completion
  report with: the new/extended test files and the branches each covers;
  per-package coverage percentages for `internal/issues`, `internal/policy`,
  `internal/journal`, and `cmd/fid` from `go tool cover -func`; full task-gate
  results; confirmation that no production behavior changed (or a deviation note
  if an unavoidable test seam was added); and explicit reference to QA-11,
  COV-001, and TRG-002.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: None.
- Audit or release-readiness impact: Resolves QA-11 / COV-001 blocking finding
  for Phase 1 wrap. No dependency change expected.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. A focused test command
(e.g. `go test -coverprofile=coverage.out ./internal/issues/ ./internal/policy/
./internal/journal/`) may be used during the red/green loop, but it does not
substitute for the full task gate at close. Every required command must be
recorded in the completion report as `PASS`, or as `BLOCKED` with reason in a
blocker report. Missing, silently skipped, or narrower substituted commands are
failed gates.

QA-plan command the Developer must also confirm at close (QA-11):

- After the gate, read `go tool cover -func=coverage.out` and confirm each of
  `internal/issues`, `internal/policy`, and `internal/journal` is ≥ 80% and
  `cmd/fid` is ≥ 70%. Record the numbers in the completion report.

Conditional commands:

```text
go test -tags integration -race ./...
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` only if any new test is placed behind
  the `integration` build tag. Prefer plain unit tests where possible.

## Notes For Developer

- This is a test-only remediation. The production code in these packages was
  reviewed and is believed correct; your job is to prove the untested branches,
  not to redesign anything.
- Do not raise package coverage by deleting or weakening assertions elsewhere,
  by adding tests that execute code without asserting outcomes, or by removing
  hard-to-cover code. Cover real behavior.
- `internal/policy` corrupt-tail tolerance (the `if len(roles) > 0 break`
  branches) is the TRG-002 behavior. Assert both states explicitly: tolerated
  trailing corruption after a valid role event, and a surfaced error when the
  first line is malformed. If you judge the current tolerance to be wrong,
  report it as a finding rather than changing it here.
- Keep CLI-level behavior in `cmd/fid` tests; add package-level tests here so
  Go attributes the coverage to the domain packages COV-001 names.
