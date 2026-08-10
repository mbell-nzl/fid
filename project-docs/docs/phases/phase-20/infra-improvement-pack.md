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

- Pack status: `APPROVED`
- Branch: `phase-20`
- Last verified commit: `1edb689b2df560a1f7ce60875458a3e639e60910`
- Current task: Task 3g - reuse one integration binary
- Next exact action: commit Task 3g's finalized report, then run its exact
  report-lint command
- Open blockers: none
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
| Task 3g shared binary | IN_PROGRESS | - | `tasks/task-3g-done.md` | - | - |
| Task 3h workspace fixture | PENDING | - | `tasks/task-3h-done.md` | - | - |
| Task 3i remove `os.Chdir` | PENDING | - | `tasks/task-3i-done.md` | - | - |
| Task 3j bounded parallelism | PENDING | - | `tasks/task-3j-done.md` | - | - |
| Task 3k claims clock | PENDING | - | `tasks/task-3k-done.md` | - | - |
| Task 3l maintenance clock use | PENDING | - | `tasks/task-3l-done.md` | - | - |
| Final measurements/gates | PENDING | - | this file | - | - |
| Independent infra review | PENDING | - | `qa/infra-review-1.md` | - | - |

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
