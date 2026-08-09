# Phase 20 - Mid-Phase Review 1 (Tasks 1-3)

Reviewer: qa-reviewer
Date: 2026-08-09
Tasks covered: phase-20/task-1 through phase-20/task-3

## Verdict

BLOCK

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-20/tasks/task-1.md`
  - `project-docs/docs/phases/phase-20/tasks/task-2.md`
  - `project-docs/docs/phases/phase-20/tasks/task-3.md`
- Completion reports:
  - `project-docs/docs/phases/phase-20/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-20/tasks/task-2-done.md`
  - `project-docs/docs/phases/phase-20/tasks/task-3-done.md`
- Validation contract: `project-docs/docs/phases/phase-20/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-20/qa/qa-plan.md`
- Commits and changed files:
  - `ab4ac8f288b3f180b698bab8924b92f375d09eb0` — `feat(cli): add response contract foundation`
  - `cd53c43319a3301d3a8cc0089b6785d9afbb163f` — `feat(project): add strict project selection config`
  - `42f5d05244eb67aaa9b389e2c612a5d667a16636` — `feat(project): stage atomic fid initialization`

All three reported hashes resolve to commits, the subjects and `Refs: phase-20/task-N`
footers match, and `git diff --check 8415ba0..HEAD` is clean. The pre-existing
working-tree modification to `project-docs/docs/phases/phase-20/run-ledger.md`
was not edited by this review.

## Summary

- Blocking findings: 2
- Improvement findings: 2
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-008, VAL-017, VAL-026, VAL-027, VAL-039, VAL-042, VAL-043, VAL-044

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | Tasks 1, 3 | ON TRACK — shared close at Task 18 | `internal/cli/envelope_test.go:9` asserts the five-key/newline success transport; `cmd/fid/init_test.go:125` proves the `project.init` success envelope and stdout-only channel. |
| VAL-002 | Tasks 1, 3 | NOT ON TRACK | Task 3 proves the typed initialized-project failure envelope, but Task 1's registry test does not prove the exact value types/nested schemas for every reason code; see Task 1 Blocking finding. |
| VAL-003 | Task 2 | ON TRACK | `internal/project/resolve_test.go:9`, `:29`, and `:55` cover Git-boundary stop, nearest/direct-root selection, and `.fid` symlink rejection. |
| VAL-004 | Task 3 | ON TRACK | `cmd/fid/init_test.go:125` asserts that `--project` initializes only the exact temporary root and reports that absolute root. |
| VAL-005 | Task 2 | ON TRACK | `cmd/fid/root.go:67` has the single fail-closed root registry; `cmd/fid/root_flags_test.go:9` exercises both global flags against legacy leaves, and Task 3 registers only `init`. |
| VAL-008 | Task 3 | DEFERRED — closes at Task 18 | `cmd/fid/init_test.go:21` and `:425` prove Human success stdout and collision failure stderr for `project.init`; the all-migrated-leaf matrix remains the Task 18 close. |
| VAL-017 | Task 3 | DEFERRED — closes at Task 18 | `cmd/fid/init_test.go:174-183` compares the returned Event ID and RFC3339-nanosecond timestamp with the appended genesis Event. Other leaves remain later work. |
| VAL-026 | Task 2 | DEFERRED — closes at Task 12 | `cmd/fid/workflow_assign_test.go:61` covers `--project-wide` collision rejection and no journal/projection mutation; canonical workflow schemas and precedence remain Task 12. |
| VAL-027 | Task 2 | DEFERRED — closes at Task 13 | `cmd/fid/maintenance_test.go:220` proves the legacy check ignores the target Attempt grace. The closed registry, 299/300 boundary, and both envelope branches remain Task 13. |
| VAL-039 | Tasks 2, 3 | ON TRACK | Strict generated ID/`issue_prefix` key inventory and parser round-trip appear in `internal/config/config_test.go:12` and `cmd/fid/init_test.go:73-110`. |
| VAL-042 | Task 3 | ON TRACK | `internal/bootstrap/init_test.go:12-44` injects every named pre-activation failure and asserts absent target/residue; `cmd/fid/init_test.go:425` snapshots the regular collision. A disposable source-built check also preserved a dangling `.fid` symlink and returned the typed collision failure. |
| VAL-043 | Tasks 2, 3 | NOT ON TRACK | The parser rejects representative malformed inputs and init proves one 16-character valid prefix plus malformed input, but required just-inside/outside domain matrices are incomplete; see Task 2 Blocking finding. |
| VAL-044 | Tasks 2, 3 | DEFERRED — closes at Task 18 | `project.init` has closed root/prefix transport in `cmd/fid/init.go:20-57` and its valid/invalid prefix tests. The literal all-21-leaf request/schema matrix is explicitly Task 18 work. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 exact gates | `gofmt`, `golangci-lint`, race/coverage, cover report, build | PASS | All five exact commands are recorded PASS with terminal tokens in `task-1-done.md`; `report-lint.sh` passed. Integration is not required for Task 1 by QA-plan row 147. |
| Task 2 exact gates | Task 1 set plus `go test -tags integration -race ./...` | PASS | All six commands are recorded PASS with exact text/tokens in `task-2-done.md`; `report-lint.sh` passed. |
| Task 3 exact gates | Task 2 six-command set | PASS | All six commands are recorded PASS with exact text/tokens in `task-3-done.md`; `report-lint.sh` passed. |
| QA-plan applicability | QA-plan rows 143-150 and 250-272 | PASS | No extra Task 1-3 command is required beyond the task-spec sets: integration begins at Task 2; benchmarks begin at Task 5. The phase-end QA checklist is not substituted for task evidence. |
| Independent required reproduction | `go test -tags integration -race ./...` | PASS | Run from `/home/mark/projects/fid` during this review. Exit 0 with no output (`exit 0 (no output)`). |
| Focused changed-surface reproduction | `go test -race -count=1 ./internal/cli ./internal/project ./internal/config ./internal/bootstrap ./cmd/fid` and focused `cmd/fid` init/root-flag selection | PASS | Changed packages passed; focused command result: `ok github.com/mbell-nzl/fid/cmd/fid 2.338s`. |
| Atomic collision risk check | Disposable dangling `.fid` symlink, source-built `go run ./cmd/fid --project <temp> --json init --prefix FID` | PASS | Returned the `project-already-initialized` failure envelope with exit 1; `target=symlink` afterward. Scratch directory was removed after capture; cleanup is not gate evidence. |

## Per-Task Findings

### Task 1 - Response contract foundation

#### Blocking

- [High][internal/cli/codes_test.go:8] `TestReasonCodeRegistryIsClosedAndComplete` compares only exported flattened `DetailFields` names. It cannot detect a changed detail value type, nullable rule, or nested object/object-list schema for any registered code; `TestNewFailureRejectsInvalidSubjectDetailsAndRedaction` at `:69` samples only four invalid cases. This does not substantiate the Task 1 requirement that every code has the exact details schema and leaves VAL-002's public error contract unproven. Fix direction: add an exhaustive table-driven valid/invalid payload matrix for all 23 codes, asserting each exact scalar/null/nested type and closed redaction rejection before serializing representative failures.

#### Improvement

- None.

### Task 2 - Project selection and config namespace

#### Blocking

- [High][internal/config/config_test.go:51] The strict-config rejection table covers representative invalid data but not the contract's required just-inside/outside domain matrix. In particular, no test proves the one-character valid prefix, 128-byte opaque ID boundary, `math.MaxInt64` valid executor/maintenance values, or overflow/out-of-range rejection; Task 3 adds only the 16-character prefix boundary. VAL-043 closes in Tasks 2-3, so source validation alone cannot defer this missing qualifying/non-qualifying evidence. Fix direction: add table-driven parser and init cases for every documented boundary, malformed form, signed-64-bit just-inside/outside value, and required key/section relationship.

#### Improvement

- [Low][internal/config/config.go:1] This Task 2 change leaves the non-test config file at 512 lines, above the project target of 400. Fix direction: move the unrelated pricing parsing/writing or the strict project-config parser into a narrowly named owning file while preserving the package API.

### Task 3 - Atomic project initialization

#### Blocking

- None.

#### Improvement

- [Medium][cmd/fid/init_test.go:425] The collision regression covers an existing directory but not an existing dangling `.fid` symlink. The implementation's `activateStagedWorkspace` Lstat guard handled the independently reproduced symlink case correctly, but a committed test is needed to prevent a future Stat/Lstat regression from replacing a path entry during init. Fix direction: add a disposable dangling-symlink collision fixture that asserts the typed conflict envelope, unchanged symlink, and no staging residue.

## Cross-cutting observations

- The response transport remains correctly isolated: `internal/cli` imports neither Cobra nor domain packages, and root-flag enforcement is centralized in `cmd/fid/root.go:67-86`.
- `internal/bootstrap/init.go` is now 427 lines, above the 400-line non-test target. This is a low-priority follow-up; split staging/fsync helpers only when it does not obscure the intentionally linear activation protocol.
- No dependency changes, spec deviations, report-hash issues, gate substitutions, or documentation ownership breaches were found in this window.

## Escalations

- None. Both Blocking findings have explicit acceptance criteria, are contained within Tasks 1-3, and need no architectural or product decision.

## Final Lines

Blocking: 2  Improvement: 2  Escalations: none
VERDICT: BLOCK
