# Phase 20 - Mid-Phase Review 2 (Tasks 4-6)

Reviewer: qa-reviewer
Date: 2026-08-11
Tasks covered: phase-20/task-4 through phase-20/task-6

## Verdict

CONTINUE (NON-BLOCK)

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-20/tasks/task-4.md`
  - `project-docs/docs/phases/phase-20/tasks/task-5.md`
  - `project-docs/docs/phases/phase-20/tasks/task-6.md`
- Completion reports:
  - `project-docs/docs/phases/phase-20/tasks/task-4-done.md`
  - `project-docs/docs/phases/phase-20/tasks/task-5-done.md`
  - `project-docs/docs/phases/phase-20/tasks/task-6-done.md`
- Validation contract: `project-docs/docs/phases/phase-20/qa/validation-contract.md`.
- QA plan: `project-docs/docs/phases/phase-20/qa/qa-plan.md`, especially EVID-01, EVID-03, EVID-04, EVID-05 and Review 2 scope.
- Prior phase evidence: `review-1.md`, `review-1b.md`, `review-1c.md`, `infra-review-1.md`, the current records, `run-ledger.md`, DEC-135 through DEC-138, and the resolved Task 5 triage/escalation records.
- Committed implementation and report commits:
  - `e56b3a0ce0433f1044f70b7969c14e19234dfdf8` — `feat(rebuild): map typed authority and lock failures`
  - `028855b807061b1d1c62a3c4d1b642b60d2f8107` — `feat(issue): remove retired issue fields`
  - `29c52469be0277ac9e3e84cd0c5f63cd16a06066` — `feat(issue): add replayable aggregate revisions`

All three reported implementation hashes resolve with `git rev-parse --verify <hash>^{commit}`. Their subjects and `Refs: phase-20/task-4`, `task-5`, and `task-6` footers match. `git diff --check` is clean for each implementation commit. The completion-report commits are present and `report-lint.sh` reports each reviewed report as committed and clean. The review began with a pre-existing modified `run-ledger.md` and an untracked developer-stoppage analysis; neither is part of this review or was changed.

## Summary

- Blocking findings: 0
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002, VAL-006, VAL-007, VAL-008, VAL-009, VAL-010, VAL-011, VAL-012, VAL-017, VAL-044

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | Task 4 | ON TRACK — shared close at Task 18 | `cmd/fid/rebuild_test.go:264-350` asserts both held-lock failure and post-release JSON success channels/key sets; `event:null` success is also reproduced through the focused command. The all-retained-leaf success matrix remains Task 18. |
| VAL-002 | Task 4 | ON TRACK — later migrated leaves remain in Tasks 8-13/Task 18 | `cmd/fid/rebuild_test.go:17-151` covers missing, unknown, and denied authority with unchanged journal bytes; `response_errors.go:13-38` maps only those retained policy errors to their closed semantic codes. This window proves the rebuild row, not the later leaves. |
| VAL-006 | Task 4 | ON TRACK | The focused reproduction executed `TestRebuildJSONMapsMissingAuthorityWithoutMutation` and `TestRebuildJSONMapsUnknownAndDeniedAuthorityWithoutMutation`; both qualifying permitted and non-qualifying missing/unknown/denied states are committed in `cmd/fid/rebuild_test.go:17-151`, with zero journal append on rejection. |
| VAL-007 | Task 4 | ON TRACK | `cmd/fid/rebuild_test.go:264-350` deterministically holds the write lock, asserts `lock-unavailable`/`retryable:true` and unchanged journal, then retries after release successfully. No timing race is used as the proof. |
| VAL-008 | Task 4 | DEFERRED — Task 18 close | This task proves rebuild Human success stdout-only and JSON failures stderr-only (`cmd/fid/rebuild_test.go:264-350`, `:397-436`). The all-migrated-leaf channel matrix is explicitly owned by Task 18. |
| VAL-009 | Task 5 | ON TRACK | Migration test `internal/db/migrate_test.go:10-42` proves all three columns and the dependent CHECK are absent after up and restored on down. CLI/JSON/effect absence and fail-closed replay are exercised by `cmd/fid/issue_section_test.go:13-42`, `:212`, `internal/issues/issue_test.go:851`, and `internal/journal/replay_test.go:31-71`; the focused suite passed. |
| VAL-010 | Task 6 | ON TRACK — Tasks 9-11 own the remaining aggregate families | The exact six currently owned families are the golden sequence `issue.created`, `issue.updated`, `issue.transitioned`, `dev.started`, `dev.finished`, and `dev.blocked` in `internal/journal/replay_test.go:123-176`. `BuildIssueRevisionEffects` sorts result Effects by Issue ID (`internal/issues/revision.go:18-56`); AC, Issue-label, and link mutations are explicitly later Task 10-11 work. |
| VAL-011 | Task 6 | ON TRACK — Tasks 9-11 own the remaining families | The six-family golden proves create revision 1 and each update-like `Values=Prior+1` / matching `Prior` and `Preconditions`; `TestReplayJournalRejectsStaleIssueRevisionBeforeApplyingEffect` (`internal/journal/replay_test.go:283-336`) independently proves the non-qualifying stale-replay state stops before the stale effect applies. |
| VAL-012 | Tasks 5-6 | ON TRACK — QA Lead owns phase-end outcome | Task 5's `TestV01ParityFixtureRebuildPreservesEveryProjectionTable` dumps every non-system projection table before/after rebuild (`cmd/fid/integration/v01_parity_test.go:47-131`). Task 6 separately proves migration 12 up/down and default revision 1 (`internal/db/migrate_test.go:44-71`). The independent dedicated integration gate and all six fixture manifest entries pass. |
| VAL-017 | Task 4 | DEFERRED — Task 18 close | Projection-only rebuild success has `event:null` at `cmd/fid/rebuild_test.go:337-349`; the journal-tail identity comparison for every mutating retained leaf remains the Task 18 closure. |
| VAL-044 | Task 4 | DEFERRED — Task 18 complete request matrix | `TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork` (`cmd/fid/rebuild_test.go:153-219`) proves both invalid positional and invalid JSON-plus-verbose requests reject with zero journal effect; accepted/missing explicit-project states are covered at `:221-262`. The literal all-21-leaf matrix remains Task 18. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 4 required commands and QA-plan applicability | `report-lint.sh task-4-done.md task-4.md` | PASS | Exit code 0. The report records all six spec commands with valid `go test` tokens; the QA plan adds no benchmark gate to Task 4. DEC-136's mandatory deviation records the original lint stop and the three superseded invalid fragments. |
| Task 5 required commands and QA-plan applicability | `report-lint.sh task-5-done.md task-5.md` | PASS | Exit code 0. All seven task/QA-plan commands, including the exact benchmark and dedicated integration gates, are recorded PASS; the report includes the mandatory DEC-137/DEC-138 deviation with first failure token `FAIL`, cause, constrained fix, and byte-identical rerun token. |
| Task 6 required commands and QA-plan applicability | `report-lint.sh task-6-done.md task-6.md` | PASS | Exit code 0. All seven task/QA-plan commands, including benchmark and dedicated integration, are recorded PASS with structurally valid test tokens. |
| Required independent Review 2 minimum | `go test -bench . -benchmem ./...` | PASS | Executed exactly from `/home/mark/projects/fid`; complete process result terminated with `exit_code: 0`. It emitted benchmark results including `BenchmarkRebuildWorkspace100kEvents-6 ... 24021760392 ns/op`, then ended `ok  github.com/mbell-nzl/fid/internal/workflow  0.568s`. |
| Focused Task 4-6 behavioral reproduction | `go test ./cmd/fid ./internal/db ./internal/issues ./internal/journal ./internal/dev -run '^(TestRebuildJSONMapsMissingAuthorityWithoutMutation|TestRebuildJSONMapsUnknownAndDeniedAuthorityWithoutMutation|TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease|TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork|TestIssueCommandsDoNotRegisterRetiredFieldFlags|TestIssueGetSectionsAndRetiredFieldAbsence|TestIssuesTableExcludesRetiredPhase20Fields|TestIssueRevisionMigrationAddsAndRemovesRevision|TestCreateEventOmitsRetiredPhase20Fields|TestBuildIssueRevisionEffectsSortsObservedRevisionsAndDerivesResults|TestUpdateEventRecordsReplayableRevision|TestUpdateRejectsInvalidChangesWithoutAppending|TestUpdateRejectsNoOpWithoutAppending|TestTransitionAllowsWorkflowEdgeProjectsActivityAndRebuilds|TestReplayJournalClassifiesRetiredIssueEffectKeyAsStopped|TestReplayJournalAppliesPhase20IssueRevisionGolden|TestReplayJournalRejectsStaleIssueRevisionBeforeApplyingEffect|TestStartClaimsTransitionsAndOpensRunInOneEvent|TestFinishClosesRunCompletesTransitionsAndReleasesClaimInOneEvent|TestBlockRecordsFailedGateCompletionTransitionsAndReleasesClaim)$' -count=1` | PASS | Exit code 0; current-tree result includes `ok github.com/mbell-nzl/fid/cmd/fid 0.353s`, `internal/db 0.084s`, `internal/issues 0.399s`, `internal/journal 0.077s`, and `internal/dev 0.183s`. |
| Dedicated real-process/rebuild parity | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | Independently executed from repository root with exit code 0: `ok github.com/mbell-nzl/fid/cmd/fid/integration  4.755s`. |
| Fixture integrity | `sha256sum -c SHA256SUMS` from `cmd/fid/testdata/v01_parity` | PASS | Exit code 0; all six entries ended `OK`, including the three regenerated Issue-get fixtures. |
| Task 5 classification-repair scope | `git diff 028855b^ 028855b -- cmd/fid/integration/classification_test.go`; DEC-137/DEC-138; `classification_test.go:15-38` | PASS | The repair changes only verification-apparatus inventory semantics: the four Task 3f journeys remain required, the failure names a missing journey, and the repository-wide tagged-test walk remains. It changes no product assertion, threshold, budget, tolerance, command, or product source. The report's red evidence records a real omitted-required-journey failure (`FAIL`) before restoration; the independent tagged classification/parity tests pass. DEC-137 — not DEC-138 alone — authorizes replacement of the non-contractual exact-growth tripwire; DEC-138's evidence-record and bookkeeping safeguards are satisfied. |

## Per-Task Findings

### Task 4 - Rebuild, authority, and lock failures

#### Blocking

- None.

#### Improvement

- None.

Observation: fail-closed authority is still evaluated before `journal.RebuildWorkspace` (`cmd/fid/rebuild.go:44-56`), and the public lock response derives only the project-local `.fid/write.lock` path (`cmd/fid/response_errors.go:41-50`), preserving the closed lock detail contract without disclosing lock-owner metadata.

### Task 5 - Remove orphaned Issue fields

#### Blocking

- None.

#### Improvement

- None.

Observation: the migration up/down pair removes/restores the retired columns and CHECK, while the independent parity test compares every projection table after rebuild. The DEC-137 repair retained the four migration-era journeys and DEC-134's repository-wide tagged-test inventory; it is not a product-test weakening.

### Task 6 - Replayable Issue aggregate revision

#### Blocking

- None.

#### Improvement

- [Low][internal/issues/revision_test.go:5-37] The helper test asserts only the successful sorted two-Issue result. `BuildIssueRevisionEffects` has observable error branches for an empty input, blank ID, non-positive observed revision, and duplicate ID (`internal/issues/revision.go:20-40`), but none is directly tested. Fix direction: add a table-driven invalid-input test that asserts each error and no returned Effects; retain the current six-family integration and replay evidence.

Observation: the golden/replay test deliberately covers only the six Task 6 families. This is aligned with the task boundary; AC, Issue-label, and link revision families remain explicit owners of Tasks 10-11 and therefore are not misreported as complete here.

## Cross-cutting observations

- The review applies the prior evidence correction discipline: it accepted neither remembered test output nor silent wrapper completion. The benchmark was retained through both the nested `session_id` result and the terminal `exit_code: 0` polling result.
- No dependency changes, unrecorded spec deviation, scope breach, report-hash mismatch, invalid evidence token, or unresolved gate substitution was found. Task 5's report is the required exception record for the DEC-137 verifier-inventory correction and DEC-138 bookkeeping evidence; its mandatory first-failure and final-result tokens are present.
- No new escalation is warranted. The only finding is local test coverage, has an unambiguous fix, and does not change contract semantics or cross the Task 6 boundary.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-20/qa/review-2.md
Blocking: 0  Improvement: 1  Escalations: none
VERDICT: CONTINUE (NON-BLOCK)
```
