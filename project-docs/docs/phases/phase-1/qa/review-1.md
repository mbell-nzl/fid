# Phase 1 - Mid-Phase Review 1 (Tasks 1-3)

Reviewer: qa-reviewer
Date: 2026-06-06
Tasks covered: phase-1/task-1 through phase-1/task-3

## Verdict

BLOCK

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-1/tasks/task-1.md`
  - `project-docs/docs/phases/phase-1/tasks/task-2.md`
  - `project-docs/docs/phases/phase-1/tasks/task-3.md`
- Completion reports:
  - `project-docs/docs/phases/phase-1/tasks/task-1-completion.md`
  - `project-docs/docs/phases/phase-1/tasks/task-2-completion.md`
  - `project-docs/docs/phases/phase-1/tasks/task-3-completion.md`
- Validation contract:
  - `project-docs/docs/phases/phase-1/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-1/qa/qa-plan.md`
- Commits or changed files:
  - `afbd61b54f9cb0eb0bc301e397727298a72d46c5` - `feat(db): add sqlite migration foundation`
  - `735b2dcfc09369f405d0793acb8c25a8c6fcf11b` - `feat(journal): add canonical event model`
  - `b0f36534963ae4768629d9850aca36e3f46c89ff` - `feat(journal): add write lock and append writer`
  - Current working tree also contains later Phase 1 task changes; this review only assessed Tasks 1-3 except for cross-window architecture/source-state checks.

## Summary

- Blocking findings: 1
- Improvement findings: 0
- Escalations raised: ESC-001
- Covered validation assertions: VAL-014, VAL-007, VAL-009, VAL-013, VAL-015

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-014 | task-1 | NOT ON TRACK | Completion report cites `TestOpenEnablesRequiredPragmas`, and `go test -race ./internal/db ./internal/journal` passed in this review. Source review found `internal/db.Open` applies pragmas by `Exec` on `*sql.DB`, which does not prove future physical connections opened by the pool also have connection-local pragmas enabled. See Blocking finding and ESC-001. |
| VAL-007 | task-2, task-3 | ON TRACK | Task 2 report cites canonical v1 envelope and golden round-trip tests. Task 3 report cites one-line append, clock-skew no-append, and concurrent append tests. Review spot-check `go test -race ./internal/db ./internal/journal` passed. Full logical CLI write evidence remains assigned to Task 6. |
| VAL-009 | task-2 | ON TRACK | Task 2 report cites `TestUpgradeToCurrentLeavesV1EventUnchanged`, `TestUpgradeChainDispatchesByEventSchema`, and committed fixture `testdata/journal/event-v1.golden.json`; commit `735b2dc` verified. Full replay golden evidence remains assigned to Task 4. |
| VAL-013 | task-3 | ON TRACK | Task 3 report cites owner metadata, holder-named timeout error, and race-enabled concurrent write tests. Review spot-check `go test -race ./internal/db ./internal/journal` passed. |
| VAL-015 | task-2, task-3 | ON TRACK | Task 2 report cites 6/8 hex issue ID shape, collision retry, and event ID formatting tests. Task 3 report cites first/next per-day sequence allocation and concurrent monotonic sequences under lock. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 required command | `gofmt -l cmd internal` | PASS | Completion report records exact command as PASS. Review rerun produced no output. |
| Task 1 required command | `golangci-lint run ./...` | PASS | Completion report records exact command as PASS. |
| Task 1 required command | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report records exact command as PASS. |
| Task 1 required command | `go tool cover -func=coverage.out` | PASS | Completion report records exact command as PASS. |
| Task 1 conditional command | `go test -tags integration -race ./...` | NOT RUN | Completion report explicitly says not required because migration tests were not behind the `integration` build tag. |
| Task 1 dependency budget | `modernc.org/sqlite`; `github.com/golang-migrate/migrate/v4` | PASS | Completion report records need, alternative, and license for both accepted dependencies. |
| Task 2 required command | `gofmt -l cmd internal` | PASS | Completion report records exact command as PASS. Review rerun produced no output. |
| Task 2 required command | `golangci-lint run ./...` | PASS | Completion report records exact command as PASS. |
| Task 2 required command | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report records exact command as PASS. |
| Task 2 required command | `go tool cover -func=coverage.out` | PASS | Completion report records exact command as PASS. |
| Task 2 conditional command | `go test -tags integration -race ./...` | NOT RUN | Task spec says none expected; completion report does not claim it. |
| Task 2 fixture evidence | `testdata/journal/event-v1.golden.json` | PASS | Fixture exists in the current tree; commit `735b2dc` includes it. |
| Task 3 required command | `gofmt -l cmd internal` | PASS | Completion report records exact command as PASS. Review rerun produced no output. |
| Task 3 required command | `golangci-lint run ./...` | PASS | Completion report records exact command as PASS. |
| Task 3 required command | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report records exact command as PASS. |
| Task 3 required command | `go tool cover -func=coverage.out` | PASS | Completion report records exact command as PASS. |
| Task 3 conditional command | `go test -tags integration -race ./...` | NOT RUN | Completion report explicitly says not triggered because lock/concurrency tests are in the default suite. |
| Task 3 race/concurrency spot-check | `go test -race ./internal/db ./internal/journal` | PASS | Review rerun passed: `ok github.com/mbell-nzl/fid/internal/db`; `ok github.com/mbell-nzl/fid/internal/journal`. |
| Commit verification | `git rev-parse --verify <hash>^{commit}` | PASS | All three reported hashes resolve to commits with the reported subjects. |
| Early architecture boundary | `rg -n "github.com/mbell-nzl/fid/cmd/fid|cmd/fid" internal` | PASS | Exit 1 with no output, meaning no matches. Task 2 commit is pure `internal/journal` data code; Task 3 I/O is contained in the journal adapter. |

## Per-Task Findings

### Task 1 - SQLite Connection, Pragmas, And Minimal Migration Harness

#### Blocking

- `internal/db/db.go:17` `Open` returns a pooled `*sql.DB` but applies `PRAGMA foreign_keys = ON` and `PRAGMA recursive_triggers = ON` by executing statements once on the pool at `internal/db/db.go:23`. SQLite pragma state is per physical connection, so later connections opened by `database/sql` are not proven to inherit those settings; `internal/db/db_test.go:10` only checks the currently used connection and does not force multiple physical connections. This means VAL-014's "every connection opens with both pragmas ON" pass criterion is not met. Fix direction: enforce the pragmas for each physical SQLite connection, for example through driver/DSN connector configuration or a deliberate pool constraint, and add a test that forces multiple underlying connections and verifies both pragmas on each. Escalated as ESC-001 because the DB helper is shared by later replay/init/issue work.

#### Improvement

- None.

### Task 2 - Journal Event v1 Model, ID Generation, Canonical JSON, Upgrader Chain

#### Blocking

- None.

#### Improvement

- None.

### Task 3 - Write Lock And Journal Writer (Append + Flush)

#### Blocking

- None.

#### Improvement

- None.

## Cross-cutting observations

- Command coverage for Tasks 1-3 is otherwise clean: every required exact command from the task specs and QA plan is recorded as PASS, and every applicable conditional command is explicitly evaluated.
- The current working tree contains later Phase 1 files and modified later-task reports. Those were not reviewed as Task 1-3 implementation evidence, but the import scan still found no `internal` import of `cmd/fid`.

## Escalations

- ESC-001 (QA Reviewer trigger 5 - spans tasks or subsystems beyond the review window): `internal/db.Open` does not yet prove pragma enforcement for every physical SQLite connection created through the shared DB pool.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-1/qa/review-1.md
Blocking: 1  Improvement: 0  Escalations: ESC-001
VERDICT: BLOCK
```
