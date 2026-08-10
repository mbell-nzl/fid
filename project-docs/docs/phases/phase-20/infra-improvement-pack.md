# Phase 20 Infrastructure Improvement Pack

> Authoritative execution anchor for the Human-approved pre-Task-4 campaign.
> Read this file first after every compaction or fresh session. Chat history and
> session plan displays are secondary.

## Authorization And Boundary

- Approved by the Human on 2026-08-10; recorded in DEC-133.
- Owner: recovery super-agent across Tasks 3e-3l, with separate task commits,
  completion reports, and gates.
- Sequence: Review 1c -> Tasks 3e-3l -> independent infra review -> Task 4.
- Review 1c: COMPLETE and NON-BLOCK at `cb5f1bfb13048c41463ef772757022ce0237681f`.
- No product contract, VAL, schema, journal format, CLI behavior, durability
  guarantee, or dependency change is authorized.
- Production initialization optimization is out of scope; benchmark only.

## Current Checkpoint

- Pack status: `COMPLETE`
- Branch: `phase-20`
- Last verified commit: `cd9e6bcd599ab5c9914ac76731c095dcc3b5a0c7`
- Current task: none; infrastructure pack complete
- Frozen review range: `81e1576169d3d798b6ee0391f82955210dfb6823^..e80cb3bd48925596b713a3a6f3baec19c87256f9`
- Frozen correction range: `75c856a^..e96f263`
- Next exact action: return control to the normal Phase Runner at Task 4; do
  not treat this recovery authority as continuing into product implementation
- Open blockers: none; ESC-20-004 resolved by correction `75c856a` and the
  same reviewer's NON-BLOCK commit `cd9e6bc`
- Active command: none

When a command is live across a tool yield, replace `Active command: none` with
the exact command, working directory, nested `session_id`, last poll state, and
the fact that no `exit_code` exists yet. Replace it with terminal evidence only
after `exit_code` or an explicit signal is observed. Transient session IDs are
removed before the final pack commit.

## Baseline And Target

Baseline from the committed 2026-08-10 diagnosis, same host:

| Command | Result | `cmd/fid` time |
|---|---|---:|
| `go test -race -coverprofile=coverage.out ./...` | PASS | 131.228s |
| `go test -tags integration -race ./...` | PASS | 162.472s |
| Combined | PASS | approximately 293.700s |

Final measurement uses three `-count=1` same-host samples for the ordinary
race/coverage and accepted dedicated integration surfaces. Target: at least 40%
lower median combined wall time. Timing is not mined for green; all samples and
host variance are recorded. Correctness, race freedom, and exact gates are hard
requirements regardless of timing.

Final same-host samples, all exit 0 and recorded without retry:

| Sample | Ordinary race/coverage `cmd/fid` | Dedicated integration | Combined |
|---|---:|---:|---:|
| 1 | 70.124s | 4.525s | 74.649s |
| 2 | 70.903s | 4.371s | 75.274s |
| 3 | 71.332s | 4.339s | 75.671s |
| Median | 70.903s | 4.371s | 75.274s |

The median combined result is 218.426s lower than the 293.700s baseline, a
74.4% reduction. Ordinary sample range is 1.208s (1.7% of median); dedicated
range is 0.186s (4.3% of median). The 40% target is satisfied without changing
correctness gates.

Task 3e production initialization benchmark, three consecutive samples on the
same host: 34.063 ms/op, 34.544 ms/op, and 35.501 ms/op (median 34.544 ms/op),
approximately 3.84 MB and 5,243 allocations per operation. Rebuild verification
and schema initialization dominate the observed intervals. DEC-133 leaves both
production stages unchanged; later tasks reuse real initialized test state.

Task 3f before/after integration inventory:

- Before: 14 tagged files, one helper and 16 top-level tests in `cmd/fid`.
- Dedicated real-process journeys moved to `cmd/fid/integration`:
  `TestV01ClosureJourneyAndRebuildParity`, `TestV01ParityGoldenOutputs`,
  `TestV02aWorkflowAuthoringExitScenarioAndRebuildParity`, and
  `TestMaintenanceCheckReportsIncompatibleWorkflowAssignment`.
- Reclassified into the ordinary `cmd/fid` suite:
  `TestV01BSmokeCreateListGetUpdateTransitionReadyAndRebuildParity`,
  `TestV01BReadyReturnsEmptyAfterAllTasksFinishedEndToEnd`,
  `TestV01CAllACsVerifiedGateAndRebuildParity`,
  `TestV01CHasCompletionReportGateAndRebuildParity`,
  `TestV01CConcernRebuildParity`, `TestV01CClosureSmoke`,
  `TestV01DClosureLifecycleSmokeAndRebuildParity`,
  `TestV01DLinkFlowAndRebuildParity`,
  `TestV01EActivityLifecycleCoverageAndRebuildParity`,
  `TestV01EDevCompositeSmokeAndRebuildParity`,
  `TestV01ReadCommandsSupportAllFormatsWithStructuralParity`, and
  `TestSQLCommandRejectsMutationAttemptsWithoutChangingProjectedCounts`.
- After: five tagged files in the dedicated package (four moved files plus the
  classification guard), exactly four journey tests, and zero tagged parent
  files. No test or assertion was removed.
- Transitional old gate: `cmd/fid` 152.981s and dedicated package 10.266s.
  The canonical dedicated command passed independently and no longer reruns
  the ordinary package.

Task 3g binary inventory: `main_test.go` contains the sole real
`exec.Command("go", "build", ...)` call; the AST harness guard proves no
`buildClosureFIDBinary` function remains. All four journeys still allocate a
distinct `t.TempDir()` workspace. The final fresh dedicated race run took
8.228s after one package-owned build, versus Task 3f's 10.266s four-build run.

Task 3h fixture inventory: package `TestMain` creates one closed template via
the complete production `InitializeWorkspace` path. The regression proves two
clones are byte-and-mode equivalent, mutation-isolated, and byte-equivalent
after rebuild. 116 exact default-init setups in 25 non-protected command-test
files now clone it. Fifty real-init calls remain across the protected init,
rebuild, journal, maintenance durability, golden/rebuild scenario, and
compatibility-heavy Issue/AC owners. Directly comparable uncached package race
time fell from 150.540s to 88.748s (41.0%); the exact coverage run reported
88.773s for `cmd/fid`.

Task 3i working-directory inventory: production `os.Getwd` exists once, in
`newRootCommand`; three test references prove capture and before/after process
state. `rg -n 'os\.Chdir|mustChdir' cmd/fid --glob '*.go'` returns the passing
no-match state `exit 1 (no output)`. All handlers consume the root context's
explicit directory, including relative `--project` resolution. The exact
five-count race stability gate stayed silent through nested session polling and
then passed after 431.955s.

Task 3j isolation inventory: 104 top-level tests in the 25 cloned-workspace
files plus the fixture owner use `t.Parallel`; the package's protected real-init,
journal, rebuild, maintenance durability, Issue/AC compatibility, and
claim-list rebuild owners remain serial. No selected file uses `t.Setenv`,
environment mutation, signals, or mutable package state. The Go scheduler is
bounded by the host's six available CPUs (`GOMAXPROCS` unset). A concurrent
two-workspace command regression is race-clean and proves journal/template
isolation. Uncached package race time fell from 88.748s to 72.130s (18.7%);
five-count time fell from 431.955s to 353.465s.

Task 3k time inventory: acquire, renew, release, list, and shared live-claim
validation resolve one optional operation timestamp; zero values use current UTC
time. Fixed tests cover just-before, exact, and just-after expiry. All three
claim command sleeps were replaced by deterministic temporary projection state.
The initial exact ten-count gate reached Go's default timeout at 10m0s. Six
isolated validation/list tests were reclassified onto the real initialized
clone and bounded parallel scheduler without changing assertions; one-count
`cmd/fid` time fell from 68.797s to 57.903s, and the byte-identical gate rerun
passed at 578.869s for `cmd/fid`.

Task 3l maintenance inventory: the domain's explicit 300-second grace test now
covers 299.999999999s, exactly 300s, and 300.000000001s. The real-binary
closure journey derives one maintenance evaluation time 24 hours after its
latest journal event and supplies it through existing `CheckOptions.Now` both
before and after rebuild; every finding and info counter must match exactly.
The four-second sleep and clock-drift allowance are gone. Dedicated integration
time fell from 8.440s to 4.331s (48.7%). Remaining `time.Sleep` references are
the production journal-lock backoff seam, its contention/cancellation test,
and the closure source guard's diagnostic string.

Infrastructure Review 1 found one Task 3f inventory omission: the direct-parent
guard could not see `internal/issues/transition_rules_test.go`, leaving that
tagged in-process test outside both accepted gates. Correction `75c856a`
reproduced the finding with a repository-wide guard, removed only the misplaced
build constraint, and retained the complete persistence assertion in the
ordinary suite. The corrected inventory permits integration tags only under
`cmd/fid/integration`; the ordinary targeted race invocation visibly ran the
retained test, and the dedicated classifier passed.

## Final Gate Results

- All eight Task 3e-3l report lints: PASS; each ended with
  `PASS report-lint completed with no failures` on a committed clean path.
- `gofmt -l cmd internal`: PASS - `exit 0 (no output)`.
- `golangci-lint run ./...`: PASS - `0 issues.`.
- `go test -race -coverprofile=coverage.out ./...`: PASS -
  `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements`.
- `go tool cover -func=coverage.out`: PASS -
  `total:										(statements)					77.7%`.
- `go test -tags integration -race ./cmd/fid/integration/...`: PASS -
  `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	(cached)`.
- `go test -bench . -benchmem ./...`: PASS - 100k rebuild measured
  24.383s/op and the last package token was
  `ok  	github.com/mbell-nzl/fid/internal/workflow	0.439s`.
- `go build -o /tmp/fid ./cmd/fid`: PASS - `exit 0 (no output)`.
- `project-docs/tools/workflow-lint.sh`: PASS -
  `PASS workflow-lint completed with no failures`.
- `go test -race -count=3 ./internal/claims ./internal/maintenance ./cmd/fid`:
  PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	172.806s`.
- `go test -tags integration -race -count=3 ./cmd/fid/integration/...`:
  PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.161s`.

Post-review correction gate on commit `75c856a`:

- Focused ordinary ownership: PASS -
  `--- PASS: TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites (0.89s)`.
- Repository-wide classification red: expected FAIL before tag removal -
  `integration-tagged tests outside dedicated package = [internal/issues/transition_rules_test.go], want none`.
- Repository-wide classification green: PASS - dedicated classifier ended
  `ok` in 1.444s.
- `gofmt -l cmd internal`: PASS - `exit 0 (no output)`.
- `golangci-lint run ./...`: PASS - `0 issues.`.
- `go test -race -coverprofile=coverage.out ./...`: PASS - the retained
  package ran for 30.529s at 82.1% coverage and the command reached
  `exit_code: 0`.
- `go tool cover -func=coverage.out`: PASS - total statement coverage 77.9%.
- `go test -tags integration -race ./cmd/fid/integration/...`: PASS -
  dedicated integration ran for 4.470s.
- `go test -bench . -benchmem ./...`: PASS - 100k rebuild measured
  24.430s/op and the final workflow package ran for 0.438s.
- `go build -o /tmp/fid ./cmd/fid`: PASS - `exit 0 (no output)`.
- `project-docs/tools/workflow-lint.sh`: PASS -
  `PASS workflow-lint completed with no failures`.

## Invariants

- Every existing integration-tagged test is accounted for by name: a true
  real-process journey moves to the dedicated package; an in-process test
  becomes an ordinary `cmd/fid` test.
- The old and new integration commands both pass in Task 3f before the old
  command is retired from active gate surfaces.
- The reusable workspace fixture is a copied real initialized workspace, not a
  mock. Rebuild parity and per-test isolation are tested.
- Init/bootstrap/migration/rebuild/durability/compatibility tests retain real
  production initialization.
- No `t.Parallel` lands until the command harness has no process-wide directory
  mutation. Only tests with isolated global state become parallel.
- Claims uses deterministic time only when supplied; zero-value production
  behavior remains real UTC time. Maintenance uses its existing `Now` option.
- No required test, assertion, race detector, coverage command, or evidence
  token is weakened.

## Progress

| Step | Status | Task commit | Report | Exact gate | Review |
|---|---|---|---|---|---|
| Review 1c | COMPLETE | `cb5f1bf` | `qa/review-1c.md` | PASS | NON-BLOCK |
| Planning / DEC-133 | COMPLETE | `81e1576` | this file | workflow-lint PASS | Human-approved |
| Task 3e benchmark | COMPLETE | `af615cb` | `tasks/task-3e-done.md` | PASS | report-lint PASS |
| Task 3f integration package/gate | COMPLETE | `ec72fda` | `tasks/task-3f-done.md` | PASS | report-lint PASS |
| Task 3g shared binary | COMPLETE | `1edb689` | `tasks/task-3g-done.md` | PASS | report-lint PASS |
| Task 3h workspace fixture | COMPLETE | `bca799f` | `tasks/task-3h-done.md` | PASS | report-lint PASS |
| Task 3i remove `os.Chdir` | COMPLETE | `3e830b4` | `tasks/task-3i-done.md` | PASS | report-lint PASS |
| Task 3j bounded parallelism | COMPLETE | `338c5dc` | `tasks/task-3j-done.md` | PASS | report-lint PASS |
| Task 3k claims clock | COMPLETE | `bec40c8` | `tasks/task-3k-done.md` | PASS | report-lint PASS |
| Task 3l maintenance clock use | COMPLETE | `16978f5` | `tasks/task-3l-done.md` | PASS | report-lint PASS |
| Final measurements/gates | COMPLETE | - | this file | PASS | NON-BLOCK after correction |
| Independent infra review | COMPLETE | `b34367a`, `cd9e6bc` | `qa/infra-review-1.md` | reviewer reproductions PASS | BLOCK, then NON-BLOCK |
| ESC-20-004 correction/re-review | COMPLETE | `75c856a` | this file | focused red/green and full gate PASS | NON-BLOCK |

Only one row may be `IN_PROGRESS`. A task becomes `COMPLETE` only after focused
red/green evidence, every exact command PASS, verified task commit, finalized
committed report, report-lint PASS, and a clean report path. A failure becomes
`BLOCKED` with the exact command/token and next action; it is never silently
retried or represented as incomplete prose.

## Task Index

| Task | Objective | Dependency |
|---|---|---|
| 3e | Benchmark production initialization without changing it | Review 1c |
| 3f | Classify/move integration tests and migrate the exact gate | 3e |
| 3g | Build one immutable integration binary per package run | 3f |
| 3h | Clone a real pristine workspace for non-init command tests | 3g |
| 3i | Inject command working directory and remove harness `os.Chdir` | 3h |
| 3j | Add measured bounded parallelism to isolated command tests | 3i |
| 3k | Replace claim-expiry sleeps with deterministic claim time | 3j |
| 3l | Replace maintenance sleep through existing `CheckOptions.Now` | 3k |

Detailed scope and exact commands live only in `tasks/task-3e.md` through
`tasks/task-3l.md`; this index does not override them.

## Final Evidence Contract

- All eight task reports pass `report-lint.sh` and their paths are clean.
- Format, lint, ordinary race/coverage, cover report, dedicated integration
  race, full benchmarks, build, workflow-lint, and repeated affected-package
  race stability all pass.
- Before/after samples, integration inventory, binary build count, remaining
  `os.Chdir`/`time.Sleep` inventories, fixture parity, and parallel isolation
  evidence are recorded.
- A fresh QA Reviewer receives the frozen pack commit range, reruns important
  gates, and commits `qa/infra-review-1.md` with Blocking/Improvement findings.
- Blocking findings are corrected in separate commits and returned to the same
  reviewer. The pack is `COMPLETE` only after a non-BLOCK outcome and final
  post-fix gates.

## Resume Protocol

1. Read this file first.
2. Verify branch, worktree, `Last verified commit`, reports, and run ledger
   against Git. Stop on disagreement; do not trust chat recollection.
3. If `Active command` contains a session with no `exit_code`, poll it before
   starting any command.
4. Read only the active task spec, its report draft if present, and affected
   code/docs.
5. Continue from `Next exact action`; update this checkpoint after each
   terminal gate or verified commit.
6. Never mark work complete from source inspection or uncommitted state.
