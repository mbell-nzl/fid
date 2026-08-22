# Phase 20 - Infrastructure Review 1 (Tasks 3e-3l)

Reviewer: qa-reviewer
Date: 2026-08-10
Frozen range: `81e1576169d3d798b6ee0391f82955210dfb6823^..e80cb3bd48925596b713a3a6f3baec19c87256f9`
Tasks covered: phase-20/task-3e through phase-20/task-3l

## Findings first

### Blocking

- [internal/issues/transition_rules_test.go:1; cmd/fid/integration/classification_test.go:37] The accepted gate migration does not account for every integration-tagged test. `internal/issues/transition_rules_test.go` remains tagged, but the classification guard only inspects the direct `cmd/fid` parent and the canonical command only runs `./cmd/fid/integration/...`. Therefore the ordinary race/coverage gate excludes the test and the dedicated gate cannot execute it. This contradicts DEC-133 and the pack invariant at `infra-improvement-pack.md:178` that every existing tagged test is either a dedicated real-process journey or an ordinary in-process test. Fix direction: make the test part of an accepted gate (normally remove its integration tag so the ordinary suite runs it) and replace the inventory guard with a repository-wide tagged-test inventory that prevents a recurrence. **ESC-20-004.**

### Improvement

- None.

## Verdict

**BLOCK.** Task 4 must not begin. The retained real concern-write transition-rule test has no accepted canonical gate after Task 3f's migration. This is a gate-definition and cross-subsystem coverage regression, not a product-behavior failure.

## Scope and evidence

- Reviewed DEC-133, the infrastructure pack, QA-plan addendum, validation contract, Tasks 3e-3l and all eight completion reports, the full frozen commit range, and the relevant canonical workflow/quality documentation.
- Each completion report claims no new `VAL-NNN`; this agrees with the validation contract and DEC-133. No product validation assertion is due from this infrastructure-only window.
- Independently ran `report-lint.sh` for every task report/spec pair. All eight passed, including exact-command coverage, report commit verification, subject verification, and committed-clean report paths.
- Verified the eight reported implementation hashes resolve and subjects match through the independent report-lint runs: `af615cb`, `ec72fda`, `1edb689`, `bca799f`, `3e830b4`, `338c5dc`, `bec40c8`, and `16978f5`.

## Gate equivalence reproduction

| Check | Result | Evidence |
|---|---|---|
| Canonical dedicated integration gate | PASS | `go test -tags integration -race ./cmd/fid/integration/...` ended `ok  github.com/mbell-nzl/fid/cmd/fid/integration (cached)`. |
| Retained tagged test, tagged invocation | PASS | `go test -tags integration -race -count=1 -run '^TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites$' ./internal/issues` ended `ok  github.com/mbell-nzl/fid/internal/issues 1.897s`. |
| Same test under ordinary gate semantics | FAIL as coverage proof | `go test -race -count=1 -run '^TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites$' ./internal/issues` ended `ok  github.com/mbell-nzl/fid/internal/issues 1.015s [no tests to run]`; the build tag excludes it. |
| Workspace, directory, and isolation regressions | PASS | Focused race invocation ended `ok  github.com/mbell-nzl/fid/cmd/fid 2.147s` for clone/rebuild parity, concurrent clone isolation, injected-directory behavior, production capture, and relative-project resolution. |
| Claims fixed-time behavior | PASS | Focused race invocation ended `ok  github.com/mbell-nzl/fid/internal/claims 2.844s`, covering zero UTC default, acquire boundary, renew/release timestamps, and list/live validation. |
| Maintenance fixed-time boundary | PASS | Focused race invocation ended `ok  github.com/mbell-nzl/fid/internal/maintenance 1.295s` for before/exact/after 300-second eligibility. |

## Per-task observations

### Task 3e - Production initialization benchmark

#### Blocking

- None.

#### Improvement

- None. `BenchmarkInitializeWorkspace` exercises `InitializeWorkspace` without production-path edits. The three recorded metrics (34.063, 34.544, and 35.501 ms/op) yield the stated 34.544 ms/op median; the report records every required command, including the then-canonical transitional integration command.

### Task 3f - Dedicated integration package and gate migration

#### Blocking

- [internal/issues/transition_rules_test.go:1; cmd/fid/integration/classification_test.go:37] See the finding above. The migration's direct-parent-only inventory silently omits a pre-existing tagged in-process test, so old/new transition evidence did not establish gate equivalence. Fix direction: ensure all tagged tests are accounted for repository-wide and run by exactly one accepted gate. **ESC-20-004.**

#### Improvement

- None. Aside from the omission, the four named real-process journeys, canonical command surfaces, and test-only adapter documentation are coherent.

### Task 3g - Reuse one integration binary

#### Blocking

- None.

#### Improvement

- None. `cmd/fid/integration/main_test.go` owns one `TestMain` build, uses a package-private temporary binary path, reports build/cleanup failure, and each journey continues to allocate `t.TempDir()` workspaces. The guard rejects the retired per-test builder.

### Task 3h - Reusable pristine workspace fixture

#### Blocking

- None.

#### Improvement

- None. `workspace_fixture_test.go` constructs the template through real `InitializeWorkspace`, preserves modes while cloning, tests clone/template/sibling isolation, and tests byte-equivalent projection state after rebuild. Protected real-init owners remain separately evidenced in the pack.

### Task 3i - Explicit command working directory

#### Blocking

- None.

#### Improvement

- None. Production `newRootCommand` captures `os.Getwd` once; handlers obtain the injected context directory. The focused reproduction proves a command uses its supplied root without changing process CWD, and `rg` inventory evidence records no command-test `os.Chdir`.

### Task 3j - Bounded command-test parallelism

#### Blocking

- None.

#### Improvement

- None. Parallel tests use the existing Go scheduler rather than a custom pool. The concurrent-clone regression is race-clean and asserts both journal separation and pristine-template immutability; protected real-init/durability categories remain serial.

### Task 3k - Deterministic claims time

#### Blocking

- None.

#### Improvement

- None. Claims options resolve a supplied instant once per operation and `operationTime(time.Time{})` uses real UTC. Tests demonstrate qualifying and non-qualifying expiry states at just-before, exact, and just-after boundaries, including list and live-claim validation; command fixtures use persisted state instead of expiry sleeps.

### Task 3l - Deterministic maintenance time

#### Blocking

- None.

#### Improvement

- None. Maintenance uses the pre-existing `CheckOptions.Now` boundary, retains the CLI zero-value path, and tests both qualifying and non-qualifying 300-second states. The closure scenario supplies one derived fixed instant before and after rebuild and compares the complete maintenance result.

## Documentation and timing audit

- Canonical current surfaces are synchronized on `go test -tags integration -race ./cmd/fid/integration/...`: `AGENTS.md`, `08_QUALITY_PLAN.md`, `09_AI_WORKFLOW.md`, the QA plan, remaining Phase 20 task specs, and QA templates. Old-command occurrences are limited to immutable historical artifacts, the Task 3f transitional command/evidence, and the baseline diagnosis.
- The reported final timing arithmetic is correct: median ordinary 70.903s plus dedicated 4.371s equals 75.274s; 293.700s minus 75.274s equals 218.426s; the reduction is 74.4% to one decimal. The ordinary and dedicated ranges are 1.208s and 0.186s, matching the stated 1.7% and 4.3% of their medians. These results do not cure the gate-equivalence defect.
- The final evidence tokens and task report tokens are structurally valid: the independent report lints accepted all eight reports, and the focused reproductions above reached terminal exit codes.

## Escalations

- **ESC-20-004** (cross-task gate-definition regression): Task 3f's canonical integration migration leaves `TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites` outside both the ordinary and dedicated accepted gates.

## Residual risks

- Until ESC-20-004 is resolved, later changes to concern-address transition behavior can pass every active gate while this real persistence test is not executed.
- The recorded performance target has been exceeded, but its timing figures should be treated as non-authorizing while the corresponding test inventory is incomplete; any repair must preserve the canonical no-duplication design and rerun the affected exact gates.

## Commands run

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3e-done.md project-docs/docs/phases/phase-20/tasks/task-3e.md
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3f-done.md project-docs/docs/phases/phase-20/tasks/task-3f.md
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3g-done.md project-docs/docs/phases/phase-20/tasks/task-3g.md
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3h-done.md project-docs/docs/phases/phase-20/tasks/task-3h.md
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3i-done.md project-docs/docs/phases/phase-20/tasks/task-3i.md
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3j-done.md project-docs/docs/phases/phase-20/tasks/task-3j.md
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3k-done.md project-docs/docs/phases/phase-20/tasks/task-3k.md
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3l-done.md project-docs/docs/phases/phase-20/tasks/task-3l.md
go test -tags integration -race ./cmd/fid/integration/...
go test -tags integration -race -count=1 -run '^TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites$' ./internal/issues
go test -race -count=1 -run '^TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites$' ./internal/issues
go test -race -count=1 -run '^(TestPristineWorkspaceClonePreservesRebuildParityAndIsolation|TestConcurrentClonedWorkspaceCommandsRemainIsolated|TestRootCommand(UsesInjectedWorkingDirectoryWithoutChangingProcess|CapturesCurrentWorkingDirectory|ResolvesRelativeProjectFromInjectedWorkingDirectory))$' ./cmd/fid
go test -race -count=1 -run '^(TestOperationTimeUsesSuppliedUTCAndDefaultsToCurrentUTC|TestAcquireUsesSuppliedTimeAtExpiryBoundary|TestRenewUsesSuppliedTimeForValidationAndTimestamp|TestReleaseUsesSuppliedTimeForValidationAndTimestamp|TestListAndLiveValidationUseSuppliedTimeAtExpiryBoundary)$' ./internal/claims
go test -race -count=1 -run '^TestCheckStartedWithoutClaimGraceBoundary$' ./internal/maintenance
```

```text
Review written: project-docs/docs/phases/phase-20/qa/infra-review-1.md
Blocking: 1  Improvement: 0  Escalations: ESC-20-004
VERDICT: BLOCK
```

## Re-review - ESC-20-004 correction (75c856a through e96f263)

### Findings first

#### Blocking

- None. ESC-20-004 is corrected: the retained concern-write transition-rule
  test is now an ordinary test, and the repository-wide inventory prevents an
  integration-tagged test outside the dedicated package from being silently
  omitted again.

#### Improvement

- None.

### Correction audit

- Verified commits and subjects: `75c856a` (`test(integration): enforce
  repository-wide classification`), `c4131a7` (`docs(phase): record
  ESC-20-004 correction evidence`), and `e96f263` (`docs(qa): return
  ESC-20-004 for re-review`) all resolve as commits on `phase-20`.
- `75c856a` removes only the two-line `//go:build integration` constraint from
  `internal/issues/transition_rules_test.go`; the complete persistence and
  transition-rule assertion body is unchanged. It changes no accepted gate
  command, product behavior, validation assertion, or test dependency.
- `cmd/fid/integration/classification_test.go:18-39` now resolves the
  repository root and requires no tagged test outside the dedicated package.
  `taggedTestFilesOutside` walks repository test files and returns a relative
  path for every outside file beginning with `//go:build integration`; it skips
  the dedicated package, whose four real-process journeys are separately
  enumerated. Current source inventory finds exactly six tagged test files, all
  in `cmd/fid/integration`.

### Independent command evidence

| Check | Result | Terminal evidence |
|---|---|---|
| Ordinary ownership | PASS | `go test -race -count=1 -run '^TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites$' ./internal/issues` ended `ok  github.com/mbell-nzl/fid/internal/issues 1.904s`; the test is compiled and executed without the integration tag. |
| Repository-wide classifier | PASS | `go test -tags integration -race -count=1 -run '^TestIntegrationClassification$' ./cmd/fid/integration/...` ended `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.596s`. |
| Canonical ordinary race/coverage gate | PASS | `go test -race -coverprofile=coverage.out ./...` exited 0; output includes `ok  github.com/mbell-nzl/fid/internal/issues (cached) coverage: 82.1% of statements` and completed through `internal/workflow`. |
| Canonical dedicated integration gate | PASS | `go test -tags integration -race ./cmd/fid/integration/...` exited 0 with `ok  github.com/mbell-nzl/fid/cmd/fid/integration 4.307s`. |

### Re-review verdict

**NON-BLOCK.** ESC-20-004 is resolved. The correction preserves the original
assertion, restores it to the ordinary accepted gate, and prevents the same
class of repository-wide classification omission without weakening a gate.

```text
Re-review appended: project-docs/docs/phases/phase-20/qa/infra-review-1.md
Blocking: 0  Improvement: 0  Escalations: none (ESC-20-004 resolved)
VERDICT: NON-BLOCK
```
