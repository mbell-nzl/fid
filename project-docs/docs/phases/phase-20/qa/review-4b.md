# Phase 20 - Mid-Phase Review 4b (task 11)

Reviewer: qa-reviewer
Date: 2026-08-15
Tasks covered: phase-20/task-11
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-20/qa/review-4.md

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-017, VAL-023, VAL-024, VAL-025, VAL-044, VAL-045, VAL-046

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Fresh evidence |
|---|---|---|---|
| VAL-017 | task-11 | ON TRACK | `TestLinkMutationsJSONMatchContract` independently passed and decodes exact add/remove JSON envelopes, non-null Event identity, and `link.list` remains a null-Event read (`cmd/fid/link_test.go:573`). |
| VAL-023 | task-11 | ON TRACK | The same mutation test checks closed Link/revision result keys and the committed list test covers typed inverse projection and deterministic ordering (`cmd/fid/link_test.go:448`, `:573`). |
| VAL-024 | task-11 | ON TRACK | The independent traversal-matrix run passed. The committed matrix covers two-node, three-node, pre-existing cyclic storage with terminating traversal, and an acyclic diamond; self dependency is covered separately (`internal/links/cycle_test.go:49`, `:76`). |
| VAL-025 | task-11 | ON TRACK | The independent focused mutation/failure run passed. It decodes exact `invalid-argument` self-link and `link-not-found` subject/details envelopes, asserting zero journal append for every case (`cmd/fid/link_test.go:593`). |
| VAL-044 | task-11 | ON TRACK | The same failure matrix covers omitted add/remove endpoint-type flags and omitted list type as exact pre-effect JSON failures (`cmd/fid/link_test.go:593`). |
| VAL-045 | task-11 | ON TRACK | The independent at/over-limit CLI run passed; the committed domain and CLI tests respectively prove 10,000 succeeds and 10,001 returns typed `result-limit-exceeded` without truncation or mutation (`internal/links/list.go:14`, `:93-119`; `internal/links/list_test.go:111`; `cmd/fid/link_limit_test.go:12`). |
| VAL-046 | task-11 | ON TRACK | The independent JSON matrix passed. It proves a same-text Issue/requirement ID follows supplied endpoint types and verifies distinct typed missing Issue/requirement codes and zero append (`cmd/fid/link_test.go:533`, `:593`). |

## Task 11 - Link contract and dependency invariants

### Must-fix

- None.

### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-01 — incomplete typed-link evidence matrix | task-11 | RESOLVED | Reconciled `cf096da6391c020c24e3392996548158271bc080` (`fix(links): complete typed link evidence matrix`) and its report-finalization/correction commits `d98a8593`, `463f2e37`, and `42856fd3`; all hashes resolve to commits. The repair adds exact canonical add/remove JSON Event and sorted-revision assertions, typed ambiguous/missing endpoint and atomic rejection cases, full bounded cycle traversal, and at/over-limit list behavior. Independent reproductions passed: focused mutation/type matrix `ok  github.com/mbell-nzl/fid/cmd/fid  0.099s`; traversal matrix `ok  github.com/mbell-nzl/fid/internal/links  0.261s`; cycle/limit CLI matrix `ok  github.com/mbell-nzl/fid/cmd/fid  1.972s`; integration `ok  github.com/mbell-nzl/fid/cmd/fid/integration  4.636s` (each terminated `exit_code: 0`). |

## Cross-cutting observations

- The implementation repair is bounded to Task 11: supplied endpoint types now distinguish same-ID Issue/requirement pairs, dependency self-links produce the required cycle path, and returned mutation timestamps retain nanosecond precision for journal-tail identity (`internal/links/add.go:49-55`, `:130`; `internal/links/remove.go:47`, `:135`).
- `task-11-done.md` records the unchanged six commands required by the task spec and QA-plan Task 11 gate set, each `PASS` with observed `exit_code: 0` (`task-11-done.md:64-69`). Independent mechanical verification passed: `PASS report-lint completed with no failures` (terminated `exit_code: 0`).
- Review 4 MF-02 is owned by Task 12 and remains pending. It is intentionally neither reviewed nor dispositioned by this Task 11-only re-review.

## Blocker Candidates

- None.

## Existing Escalation References

- None.

## Final Lines

Review written: project-docs/docs/phases/phase-20/qa/review-4b.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
