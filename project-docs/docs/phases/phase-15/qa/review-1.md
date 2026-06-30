# Phase 15 - Mid-Phase Review 1 (Tasks 1-2)

Reviewer: QA Reviewer
Date: 2026-06-29
Tasks covered: phase-15/task-1 through phase-15/task-2

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-15/tasks/task-1.md`
  - `project-docs/docs/phases/phase-15/tasks/task-2.md`
- Completion reports:
  - `project-docs/docs/phases/phase-15/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-15/tasks/task-2-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-15/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-15/qa/qa-plan.md`
- Commits or changed files:
  - `0906428adf9bb4d89b41dc94ad5b147b1b7c538f` - `refactor: split issue and replay cleanup paths`
  - `07ac0ae653753f2a12921cf26bf0cde1c12dc6c5` - `feat(dev): add dev start and gate composites`
  - `cmd/fid/issue*.go`, `cmd/fid/dev.go`, `cmd/fid/dev_test.go`
  - `internal/journal/replay.go`, `internal/journal/projection.go`, `internal/journal/inspect.go`
  - `internal/labels/add.go`, `internal/issues/label.go`
  - `internal/dev/*.go`, `internal/policy/matrix.go`, `internal/policy/matrix_test.go`
  - `project-docs/specs/02_TRD.md`, `project-docs/specs/06_DESIGN_PLAN_CLI.md`

## Summary

- Blocking findings: 0
- Improvement findings: 2
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-006 partial

No Blocking findings were found in this review window. The required command evidence is present for both tasks, commit metadata verifies, dependency-budget notes are present, and Task 2's VAL-006 evidence is appropriately partial because Task 3 owns finish/block enforcement.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | Task 1 report records split-file `wc -l` proof for all resulting `cmd/fid/issue*.go` and `internal/journal` files, all <=400 lines. Reviewer spot-check `wc -l` matched the reported counts: `cmd/fid/issue.go` 22, largest split file `cmd/fid/issue_helpers.go` 228, `internal/journal/replay.go` 205, `projection.go` 192, `inspect.go` 128. Required task gates are recorded PASS. |
| VAL-002 | task-1 | ON TRACK | Task 1 report records function-span proof for `internal/labels/add.go` and `internal/issues/label.go`; reviewer spot-check found the targeted functions at or below 60 lines. Required package coverage evidence is recorded for `internal/labels` 82.4%, `internal/issues` 81.7%, `internal/journal` 80.5%, and `cmd/fid` 77.4%. |
| VAL-003 | task-2 | ON TRACK | `internal/dev/start_test.go` covers one `dev.started` event with claim insert, issue transition, activity effect, and open run; `cmd/fid/dev_test.go` covers CLI confirmation and projected `in_progress` state. Reviewer ran `go test ./internal/dev -run 'TestStart|TestGate'`: PASS. |
| VAL-004 | task-2 | ON TRACK | `internal/dev/start_test.go` covers successful unclaimed `todo` start plus zero-append rejection for conflicting live claim and invalid edge. Reviewer focused test above passed. |
| VAL-005 | task-2 | ON TRACK | `internal/dev/gate_test.go` covers one `dev.gate.recorded` event, active attempt derivation, open-run linking, and zero-append rejection without started state or live claim. `cmd/fid/dev_test.go` verifies record-only behavior by checking `touch should-not-exist.txt` is stored as a reported command and not executed. Reviewer focused test above passed. |
| VAL-006 | task-2 | PARTIAL ON TRACK | Task 2 adds all four matrix rows and covers section-13 parity plus `policy show`/enforce behavior for `dev.start` and `dev.gate`. Reviewer ran `go test ./internal/policy -run TestMatrixDevRowsMatchSchemaSection13` and `go test ./cmd/fid -run TestDevCommandsPolicyShowMatchesRowsAndDeniedRoles`: PASS. Per task spec and QA plan, `dev.finish`/`dev.block` enforcement remains for Task 3. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 required gate | `gofmt -l cmd internal` | PASS | Completion report records empty output. |
| Task 1 required gate | `golangci-lint run ./...` | PASS | Completion report records `0 issues.` |
| Task 1 required gate | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report records required package coverage excerpts. |
| Task 1 required gate | `go tool cover -func=coverage.out` | PASS | Completion report records total coverage 75.8%. |
| Task 1 QA-plan gate | `go test -tags integration -race ./...` | PASS | Required by QA plan because Task 1 touched `cmd/fid`; completion report records PASS. |
| Task 1 QA-plan gate | `go build -o /tmp/fid ./cmd/fid` | PASS | Completion report records empty output. |
| Task 1 report headings | Completion-report required sections | PASS WITH NOTE | All required evidence sections are present. The report omits the template H1 title; see Improvement. |
| Task 1 dependency budget | `## Dependency-budget note` | PASS | Records no dependency changes. |
| Task 1 commit verification | `git rev-parse --verify 0906428^{commit}` and `git show -s --format=%H%n%s%n%b 0906428` | PASS | Hash resolves to `0906428adf9bb4d89b41dc94ad5b147b1b7c538f`; subject and `Refs: phase-15/task-1` footer verified. |
| Task 2 required gate | `gofmt -l cmd internal` | PASS | Completion report records empty output. |
| Task 2 required gate | `golangci-lint run ./...` | PASS | Completion report records `0 issues.` |
| Task 2 required gate | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report records `internal/dev` 74.2%, `internal/policy` 84.6%, `cmd/fid` 77.3%. |
| Task 2 required gate | `go tool cover -func=coverage.out` | PASS | Completion report records total coverage 75.8%. |
| Task 2 QA-plan gate | `go test -tags integration -race ./...` | PASS | Required by QA plan because Task 2 touched `cmd/fid`; completion report records PASS. |
| Task 2 QA-plan gate | `go build -o /tmp/fid ./cmd/fid` | PASS | Completion report records empty output. |
| Task 2 focused reviewer check | `go test ./internal/dev -run 'TestStart|TestGate'` | PASS | Confirms reviewed start/gate domain behavior in this window. |
| Task 2 focused reviewer check | `go test ./internal/policy -run TestMatrixDevRowsMatchSchemaSection13` | PASS | Confirms `dev.*` rows match section 13 for the added rows. |
| Task 2 focused reviewer check | `go test ./cmd/fid -run TestDevCommandsPolicyShowMatchesRowsAndDeniedRoles` | PASS | Confirms policy show includes rows and start/gate denial is before append. |
| Task 2 report headings | Completion-report required sections | PASS WITH NOTE | All required evidence sections are present. The report omits the template H1 title; see Improvement. |
| Task 2 dependency budget | `## Dependency-budget note` | PASS | Records no dependency changes. |
| Task 2 commit verification | `git rev-parse --verify 07ac0ae^{commit}` and `git show -s --format=%H%n%s%n%b 07ac0ae` | PASS | Hash resolves to `07ac0ae653753f2a12921cf26bf0cde1c12dc6c5`; subject and `Refs: phase-15/task-2` footer verified. |

## Per-Task Findings

### Task 1 - Behavior-preserving cleanup (CF-014 + CF-015)

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-15/tasks/task-1-done.md:1` The completion report starts at `## Summary of what was implemented` and omits the H1 required by `project-docs/templates/completion-report-template.md`. This did not hide evidence, but it weakens artifact consistency. Fix direction: add the template title heading in future completion reports, e.g. `# Phase 15 - Task 1: Behavior-preserving cleanup (CF-014 + CF-015) - Completion Report`.

### Task 2 - `internal/dev` foundation + `fid dev start` + `fid dev gate`

#### Blocking

- None.

#### Improvement

- `internal/dev/start.go:34` and `internal/dev/gate.go:35` The new exported orchestration functions are materially over the 60-line non-test function target (`Start` spans lines 34-197; `Gate` spans lines 35-162), and `cmd/fid/dev.go:82` also puts the gate command constructor just over the target. Task 2's validation behavior is covered and this is not blocking the current window, but it reintroduces structure debt in a phase that is actively paying down function-size debt. Fix direction: during Task 3 or follow-up cleanup, extract option validation, state validation, event/effect construction, and result assembly helpers while keeping the one-event and validate-before-append behavior intact.
- `project-docs/docs/phases/phase-15/tasks/task-2-done.md:1` The completion report starts at `## Summary of what was implemented` and omits the H1 required by `project-docs/templates/completion-report-template.md`. This did not hide evidence, but it weakens artifact consistency. Fix direction: add the template title heading in future completion reports, e.g. `# Phase 15 - Task 2: internal/dev foundation + fid dev start + fid dev gate - Completion Report`.

## Cross-cutting observations

- Task 2's VAL-006 treatment is correctly partial for this review window: the four rows exist and start/gate show-vs-enforce evidence is present; finish/block enforcement remains assigned to Task 3 by the task spec and QA plan.
- The required integration gate discipline from the Phase 15 QA plan was followed for both tasks.
- No dependency-budget anomaly was found; both reports explicitly record no dependency changes.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-15/qa/review-1.md
Blocking: 0  Improvement: 2  Escalations: none
VERDICT: CONTINUE
```
