# Phase 17 - Mid-Phase Review 1 (Tasks 1-2)

Reviewer: qa-reviewer
Date: 2026-07-15
Tasks covered: phase-17/task-1, phase-17/task-2

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-17/tasks/task-1.md`
  - `project-docs/docs/phases/phase-17/tasks/task-2.md`
- Completion reports:
  - `project-docs/docs/phases/phase-17/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-17/tasks/task-2-done.md`
- Validation contract: `project-docs/docs/phases/phase-17/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-17/qa/qa-plan.md`
- Task commits:
  - `65d2f94587310668b252ebfd4141ed9ef26c8ec0` — `refactor(format): split markdown renderer`
  - `d9c12ecd1b59a26a8239f4c608aca62003f7acdb` — `feat(config): add maintenance grace contract`
- Finalized report commits:
  - `2240785f426c3e3360f3f57981807ffa53bd2b7b` — `docs(task): finalize completion report for phase-17 task-1`
  - `13facdf9616e705db5959634a08c0aa23a23ec0f` — `docs(task): correct completion report metadata for phase-17 task-2`

## Summary

- Blocking findings: 0
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | The committed diff only relocates private renderer helpers; the dispatcher remains in `internal/format/markdown.go`. Reproduced exact-output formatter tests pass, including populated/empty issue, populated/empty gates, and pointer `ReadyView`; all four formatter files are below the 400-line limit. |
| VAL-002 | task-2 | ON TRACK | Reproduced tests pass for legacy `[project]`-only default, bootstrap default section/value, zero and positive round trips, plus negative and malformed rejection. The implementation's default, serialization, and contextual error match DEC-050 and `04_TRD_CONFIG.md` section 2. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 report contract | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-1-done.md project-docs/docs/phases/phase-17/tasks/task-1.md` | PASS | All five exact task commands are recorded as PASS; recorded hash/subject resolve and match; report is committed and clean. |
| Task 2 report contract | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-2-done.md project-docs/docs/phases/phase-17/tasks/task-2.md` | PASS | All four exact task commands are recorded as PASS; recorded hash/subject resolve and match; report is committed and clean. |
| VAL-001 file-size gate | `wc -l internal/format/markdown.go internal/format/markdown_issue.go internal/format/markdown_lists.go internal/format/markdown_stats.go` | PASS | 88, 164, 321, and 121 lines respectively; each is at or below 400. This is the Task 1 exact required command reproduced independently. |
| VAL-001 byte-parity evidence | `go test -race ./internal/format -run '^(TestWriteMarkdownMatchesGoldenFixtureForPopulatedIssue|TestWriteMarkdownOmitsEmptySectionsForEmptyIssue|TestWriteMarkdownRendersIssueGateLists|TestWriteMarkdownRendersReadyViewPointer)$' -count=1 -v` | PASS | Golden populated issue, empty issue, both gate-list states, and exact pointer `ReadyView` output passed. `git diff 65d2f94^ 65d2f94 -- internal/format` shows private-helper extraction only; no dispatch, API, or non-markdown-format change. |
| VAL-002 legacy, round-trip, and invalid-state matrix | `go test -race ./internal/config -run '^(TestReadProjectConfigDefaultsMissingMaintenanceGraceToFiveMinutes|TestProjectConfigMaintenanceGraceRoundTrips|TestReadProjectConfigRejectsInvalidMaintenanceGrace|TestWriteBootstrapFilesRoundTripsDefaultConfig)$' -count=1 -v` | PASS | Legacy missing state reads 300; bootstrap default round-trips; zero and positive 45 round-trip; negative and malformed states both fail with key context. |
| VAL-002 bootstrap section/value | `go test -race ./cmd/fid -run '^TestInitCommandCreatesWorkspaceFiles$' -count=1 -v` | PASS | Fresh `fid init` workspace contains `[maintenance]` with `started_without_claim_grace_seconds = 300`, then reads the default back. |
| Task and report hash/subject verification | `git rev-parse --verify '65d2f94587310668b252ebfd4141ed9ef26c8ec0^{commit}'`; `git rev-parse --verify 'd9c12ecd1b59a26a8239f4c608aca62003f7acdb^{commit}'`; `git show -s --format='%H%n%s%n%b' <hash>`; `git log -1 --format='%H%n%s%n%b' -- <report>` | PASS | Both task hashes resolve with their reported Conventional subjects and `Refs` footers. The finalized report files are separately committed with the subjects listed in Scope Reviewed. |
| Entry blocker disposition | QA plan status note, QA-BLK-01/02 table, and committed `qa/review-entry.md` | PASS | QA-BLK-01 is closed by Task 0's accepted entry review; QA-BLK-02 is closed by committed model-audit `2d5e536`. Neither blocker applies to the Task 1-2 evidence window. Terminal VAL-024 remains QA Lead-owned at phase end and receives no terminal outcome here. |

## Task 1 - Close CF-019 by Splitting the Markdown Formatter

### Blocking

- None.

### Improvement

- [internal/format/markdown_test.go:12] Direct byte-exact assertions cover populated/empty issue, gate-list, and pointer-ready rendering, while the other existing dispatch cases rely on the source-only relocation check. Fix direction: on the next formatter change, add a table-driven characterization matrix with exact expected output for the remaining dispatched value families.

### Observation

- [internal/format/markdown.go:23] The private dispatch surface remains unchanged and routes each pre-existing type to the extracted helpers; the reproduced golden and exact-output tests confirm the representative rendered branches remain byte-identical.

## Task 2 - Add the Maintenance Grace Configuration Contract

### Blocking

- None.

### Improvement

- None.

### Observation

- [internal/config/config.go:145] Missing maintenance configuration starts from the named 300-second default, explicit values serialize under `[maintenance]`, and the parser rejects negative or non-integer input rather than coercing it. The reproduced matrix exercises both qualifying and non-qualifying states.

## Cross-cutting observations

- No deviation from DEC-050, no new dependency, and no change outside the reviewed formatter/config boundaries was found in the task commits.
- The unrelated existing worktree modification to `project-docs/docs/phases/phase-17/run-ledger.md` was preserved and is not part of this review.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-17/qa/review-1.md
Blocking: 0  Improvement: 1  Escalations: none
VERDICT: CONTINUE
```
