# Phase 20 - Mid-Phase Review 3b (Task 9 recovery via Task 9a)

Reviewer: qa-reviewer
Date: 2026-08-15
Tasks covered: phase-20/task-9 recovery via phase-20/task-9a

## Verdict

CONTINUE

## Scope Reviewed

- This is the DEC-140 bounded recovery re-review, not a cadence review. It reviews only the three Blocking findings in `qa/review-3.md` for Task 9 and Task 9a's completion report/gate record. Tasks 7 and 8 were not re-reviewed and Task 9's implementation commit was not reopened.
- Task specs and reports: `tasks/task-9.md`, `tasks/task-9-done.md`, `tasks/task-9a.md`, and `tasks/task-9a-done.md`.
- Contract and QA evidence: `contracts.md` §§2.4, 3.2, 3.2.1, 3.6, and 7; `qa/qa-plan.md`; and VAL-017, VAL-018, VAL-019, VAL-040, VAL-041, and VAL-044 in `qa/validation-contract.md`.
- Prior review and governing decision: `qa/review-3.md` and DEC-140 in `decisions.md`.
- Verified commits: `1f196a94b064554fbc8bae1b0d9238125e2ab1d4` (`fix(issue): recover mutation contract evidence`, `Refs: phase-20/task-9a`) and `764745785c6530e00021ded8f8a8d331ac1cf034` (`docs(task): finalize completion report for phase-20 task-9a`). Both hashes resolve with `git rev-parse --verify <hash>^{commit}`. The completion report is committed and clean.

## Summary

- Blocking findings: 0
- Improvement findings: 0
- Escalations raised: none
- Covered validation assertions: VAL-017, VAL-018, VAL-019, VAL-040, VAL-041, VAL-044

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-017 | task-9a | ON TRACK | `TestIssueTransitionCanonicalOutcomeTable` compares the successful canonical transition response Event ID and `occurred_at` with the appended journal tail at `cmd/fid/issue_transition_test.go:331`. This closes Task 9's transition portion; the all-leaf phase outcome remains Task 18. |
| VAL-018 | task-9a | ON TRACK | `TestIssueTransitionExpectedRevisionMatrix` has six command-boundary cells: absent, matching, and stale for both update and transition (`cmd/fid/issue_transition_test.go:65`). It asserts one append and revision increment on success, and `state-conflict` details, unchanged projection, and zero append on stale rejection (`:142-180`). |
| VAL-019 | task-9a | ON TRACK | The canonical transition table covers legal edge, unknown target, absent known-state edge, stored unsupported precondition, and unsatisfied precondition (`cmd/fid/issue_transition_test.go:186-277`). Every rejection checks exact `transition-not-permitted` details, unchanged status, and zero append (`:338-368`). |
| VAL-040 | task-9a | ON TRACK | The four create/update-by-owner-flag cells assert canonical failure envelopes, empty stdout, and zero append (`cmd/fid/issue_update_test.go:266-345`). The all-mutable-field update verifies owner-free changes and unchanged projected owners (`:350-406`); the named retained default-owner regression is recorded in the report. |
| VAL-041 | task-9a | ON TRACK | Create, update, and transition each assert local Event/activity attribution, permitted Event-effect tables only, and invariant counts for every other projection table (`cmd/fid/issue_transition_test.go:374-461`). This closes Task 9's Issue-mutation portion; link/workflow coverage remains assigned to later tasks and Task 18. |
| VAL-044 | task-9a | ON TRACK | The exact update/transition expected-revision request cells and four removed owner-input cells are command-boundary assertions with no durable effect on rejection. Full retained-leaf closure remains Task 18. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 9a exact gate: format | `gofmt -l cmd internal` | PASS | Report records `exit 0 (no output)` (terminated `exit_code: 0`). |
| Task 9a exact gate: lint | `golangci-lint run ./...` | PASS | Report records `0 issues.` (terminated `exit_code: 0`). |
| Task 9a exact gate: race/coverage | `go test -race -coverprofile=coverage.out ./...` | PASS | Report records `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements` (terminated `exit_code: 0`). |
| Task 9a exact gate: coverage report | `go tool cover -func=coverage.out` | PASS | Report records `total:                                    (statements)                 77.3%` (terminated `exit_code: 0`). |
| Task 9a exact gate: integration | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | Report records `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.100s` (terminated `exit_code: 0`). |
| Task 9a exact gate: build | `go build -o /tmp/fid ./cmd/fid` | PASS | Report records `exit 0 (no output)` (terminated `exit_code: 0`). |
| Mechanical report verification | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-9a-done.md project-docs/docs/phases/phase-20/tasks/task-9a.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). All six Task 9a and QA-plan-required commands are recorded exactly as PASS with a token and observed termination. |
| Required independent reproduction | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | From `/home/mark/projects/fid`: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.174s` (terminated `exit_code: 0`). |
| Focused recovery reproduction | `go test ./cmd/fid -run '^(TestIssueTransitionExpectedRevisionMatrix|TestIssueTransitionCanonicalOutcomeTable|TestIssueMutationCanonicalAttributionAvoidsProtectedEffects|TestIssueCreateAndUpdateRejectRemovedOwnerFlagsWithCanonicalFailures|TestIssueUpdateCanonicalChangesExcludeOwnerAndPreserveProjection)$' -count=1` | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid	0.381s` (terminated `exit_code: 0`). |
| Stored-rule behavior reproduction | `go test ./internal/issues -run '^TestTransitionTreatsReservedGateRuleAsUnsatisfied$' -count=1` | PASS | `ok  	github.com/mbell-nzl/fid/internal/issues	0.055s` (terminated `exit_code: 0`). |

The Task 9a report's first-run deviations are within DEC-140's authorized recovery clause. The outcome-table failures revealed contract-visible product defects: declared target-state validation now precedes authorization, and an unknown stored rule is returned as an unmet precondition. Both changes are confined to `internal/issues/transition.go:152-182` and `internal/issues/transition_rules.go:50-78`, preserve the closed `transition-not-permitted` contract, and were re-run with the same focused command. The remaining decoder and fixture corrections repaired new test harness setup without weakening, deleting, narrowing, or re-pointing an assertion, threshold, or required matrix cell.

## Per-Task Findings

### Task 9 recovery evidence - Task 9a: Mutation concurrency, transition, and attribution evidence

#### Blocking

- None.

#### Improvement

- None.

Disposition of Review 3 Blocking findings:

1. **Closed — expected revision states.** The six-cell matrix proves absent, matching, and stale `--expected-revision` behavior for both Issue update and Issue transition, including exact stale details and zero-append/unchanged-projection rejection evidence.
2. **Closed — canonical transition outcomes and Event identity.** The complete legal/unknown-target/absent-edge/unsupported-precondition/unsatisfied-precondition table uses decoded canonical JSON, checks the exact details and zero append for every rejection, and compares the successful envelope Event identity with the journal tail.
3. **Closed — owner boundary and attribution negative state.** Both removed owner flags reject for both create and update; a successful all-field update omits owner changes and leaves owners unchanged; and create/update/transition attribution proves local activity attribution with no protected Event effect or projection-table mutation.

Observation: the corrective ordering in `internal/issues/transition.go:152-169` makes an unknown target return the contract's `transition-not-permitted` result before actor authorization, while `internal/issues/transition_rules.go:76-77` fail-closes a stored unknown rule as an unmet precondition. The focused reproductions pass these behavior assertions.

## Escalations

- None. DEC-140 resolves the evidence-recovery scope, and the implemented behavior matches the accepted contract without a new cross-task, architecture, security, or acceptance-criteria ambiguity.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-20/qa/review-3b.md
Blocking: 0  Improvement: 0  Escalations: none
VERDICT: CONTINUE
```
