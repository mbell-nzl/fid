# ESC-23-002 causal repair evidence

Authority: Human delegation recorded in DEC-238, 2026-09-20. Working directory
for all commands unless explicitly stated: `/home/mark/projects/fid`.
This records the repair session; earlier Task 1 failed attempts remain in the
original report history and the diagnostic review. No threshold, accepted
assertion, dependency, Event or migration contract was changed.

## Implementation and validation map

- Rebuild owns one pinned SQLite connection and up to 128 prepared SQL shapes.
  Excess shapes execute normally on the same connection. Effects and metadata
  use the same executor inside each Event transaction; validation and reads use
  that transaction. Statements close before the connection is released and
  before the staged DB is installed. Live apply and catch-up retain their path.
- The execution-transition cache stores durable Envelope setup only. Each clone
  acquires its own claim and starts its own Run using the unchanged real command
  path and 60-second TTL. Explicit expiry tests and rejection lists are retained.
- VAL-004: retained schema-17 pragma checks and rebuild parity, plus failure
  rollback/installed-projection preservation. VAL-079: exact ordinary/integration
  gate wall times below. Rebuild timing has no VAL; it is the TRD section 9 /
  DEC-158 benchmark control. Task 1 also retains VAL-001/002/003/005/075 evidence.
- No new package or dependency. The `internal/db` executor seam permits a
  transaction-bound prepared writer without moving SQL metadata ownership.

## Red/green and failed attempts

1. Existing benchmark red, before implementation:
   `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/`
   — FAIL — `100k-event rebuild = 33.049607193s, want <30s`
   (terminated `exit_code: 1`). This reproduces the existing requirement failure.
2. New regression red:
   `go test ./cmd/fid -run '^TestExecutionTransitionFixtureDoesNotInheritExpiredClaim$' -count=1`
   — FAIL — `FAIL` (terminated `exit_code: 1`). The aged cached template produced
   `[claim-expired run-inactive workflow-transition-invalid]` instead of the
   retained single-reason expectation. The test ages the cached projection
   directly, without sleeping or changing production time/expiry semantics.
3. Initial focused race verification below — FAIL —
   `FAIL` (terminated `exit_code: 1`). Its new installed-projection test attempted
   to read `.fid/fid.sqlite` before creating it (`no such file or directory`).
   The helper creates only the workspace directories. Repair: explicitly create
   and close a migrated database before capturing the original bytes. No
   assertion or production behavior changed. The byte-identical rerun passed.
4. Focused verification after that fixture-setup repair — PASS —
   `14.879s` (terminated `exit_code: 0`); cmd/fid `33.064s`, journal `14.879s`.
   Exact command (same for attempts 3 and 4):

   ```text
   go test -race ./cmd/fid ./internal/journal -run 'TestExecutionTransitionFixtureDoesNotInheritExpiredClaim|TestExecutionBoundTransitionPredicateMatrix|TestExecutionBoundTransitionRetainsClaimHolderAndExpiryNoMutationGuards|TestPristineWorkspaceClonePreservesRebuildParityAndIsolation|TestRebuildReplayRollsBackFailedEventAndReleasesConnection|TestRebuildFailurePreservesInstalledProjectionAndJournal|TestRebuildWorkspaceIsDeterministicAtSingleProjectScale' -count=1
   ```

5. Exact benchmark green after repair, same command as attempt 1 — PASS —
   `23338 ms/rebuild` (terminated `exit_code: 0`); package `26.123s`.
   Root cause and repair: repeated preparation of FK-aware Effect writes and
   metadata upserts; prepared statements now remain valid across per-Event
   transactions within one rebuild. The red attempt remains failed evidence.

The preliminary `go test ./internal/journal ./internal/db -run
'TestProjectEvent|TestReplayJournal|TestRebuildWorkspaceIsDeterministicAtSingleProjectScale'
-count=1` also terminated `exit_code: 0`, journal `1.904s`; db explicitly said
`[no tests to run]` and is not counted as DB acceptance evidence.

## Exact task gate

- `gofmt -l cmd internal` — PASS — `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` — PASS — `0 issues.` (terminated `exit_code: 0`).

The terminated final gate results follow in the attempt and final-verification sections below.

## Review of the final implementation

- The cache is allocated per rebuild, not per package or process. Rebuilds of
  different workspaces cannot share prepared statements. The workspace write
  lock still surrounds staging and replacement.
- The pinned connection prevents preparation from borrowing a second pool
  connection while an Event transaction owns the first. The regression test
  sets the pool limit to one and successfully reads the retained checkpoint
  after a failed Event, proving rollback and release rather than just closure
  of a wrapper object.
- A failing second Event rolls back its earlier Effect while retaining the
  first Event's data and checkpoint. The public staged rebuild separately
  preserves the prior installed DB and exact journal bytes on a located
  precondition conflict, and leaves no temporary projection.
- The existing live/rebuild raw-byte equality assertion passes unchanged. No
  locking-mode/header normalization is used to satisfy it.
- Existing foreign-key and recursive-trigger configuration is unchanged. The
  executor shares the active transaction's physical connection; it does not
  disable constraints or defer validation until the end of replay.
- The cache is bounded at 128 shapes, with an uncached same-connection fallback.
  Statement-close errors are propagated and prevent installing a staged DB.
- The fixture regression retains the old expected rejection and reproduces
  the original pair of unexpected reasons before the repair. Only setup moves:
  claims and Runs are created after cloning, while durable Envelope creation
  stays cached. Production transition/expiry evaluation is untouched.

This is a direct review by the Human-delegated repair agent, not a claim of
independent-agent QA. The normal Phase Runner review windows remain required.

## Ordinary gate attempt 1: correctness passed, wall-time control failed

`go test -race -timeout 30m -coverprofile=coverage.out ./...` terminated
`exit_code: 0`, with final token `62.6% of statements`. The command's
correctness result was green, but the required budget result is FAIL:
`cmd/fid` 628.145 s and `internal/workflow` 623.362 s both exceed 600 s.
The complete output is `esc-23-002-evidence/race-attempt-1.log`.
This attempt is not accepted as task closure. No narrower gate or cached old
package timing will substitute for a fresh rerun after causal speed repair.

Investigation: the workflow `recoveryFixture` still calls the full bootstrap
for every case, even though other workflow fixtures already clone the real
pristine bootstrap created by TestMain. Bootstrap itself builds and rebuilds
schema 17 and validates its projection. Recovery cases do not test bootstrap;
they seed their own repository-bound Recovery/Claim/Run data after setup.

### Causal setup optimization after attempt 1

A representative recovery test was profiled before changing its setup:

`go test -race -cover ./internal/workflow -run '^TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse$' -count=5 -cpuprofile /tmp/fid-recovery-setup-before.cpu -o /tmp/fid-recovery-setup-before.test`

- Before: PASS, `8.7% of statements`, package `21.763s`
  (terminated `exit_code: 0`). Cumulative sampled CPU was 9.70 s in
  `recoveryFixture` and 9.46 s / 52.15% in `initializeWorkspaceWithHooks`.
- Repair: `recoveryFixture` now clones the real pristine bootstrap already
  created by workflow TestMain. Git initialization, live repository observation,
  fresh time, and all Recovery/Claim/Run seed data still happen per case.
  No test cases, assertions, or mutation/validation calls under test are removed.
- Byte-identical diagnostic command after repair: PASS,
  `8.7% of statements`, package `9.591s` (terminated `exit_code: 0`).
  Reusing the profile output paths replaced the before profile; the before
  duration and sampled-CPU figures above were captured directly in this session
  before that replacement, not inferred from the after profile.
- A new isolation test checks distinct workspace/repository bindings, observes
  each actual Git repository, mutates one projection, and verifies the other
  projection and the pristine template bytes are unchanged.

The original full gate remains a failed budget attempt. Test-cache invalidation
before the next full gate is measurement housekeeping, not a passing gate: it
ensures CLI timing is measured again rather than reusing its older cached pass.

The new isolation check passed:
`go test -race ./internal/workflow -run '^TestRecoveryFixturesKeepRepositoryAndProjectionIsolation$' -count=1`
- PASS - `4.078s` (terminated `exit_code: 0`). Format and lint were rerun after
this fixture repair: `gofmt -l cmd internal` returned no output and
`golangci-lint run ./...` printed `0 issues.`; both terminated `exit_code: 0`.

Reference identity re-captured in this repair session: Linux/amd64,
`6.6.114.1-microsoft-standard-WSL2`, `go1.26.2`, six exposed logical CPUs,
AMD Ryzen 7 5700G with Radeon Graphics, `MemTotal: 8133236 kB`.

## Ordinary gate attempt 2: workflow repaired; CLI budget still failed

The exact `go test -race -timeout 30m -coverprofile=coverage.out ./...`
terminated `exit_code: 0` with token `62.6% of statements`; every package
passed correctness. Fresh timings were workflow **233.338 s**, down from
623.362 s, and CLI **604.046 s**, still above 600 s. The required budget result
remains FAIL; this is not accepted as task closure. Raw output is
`esc-23-002-evidence/race-attempt-2.log`.

Further causal setup repair is limited to 25 setup sites in `issue_test.go`,
`issue_ac_test.go`, `issue_get_list_test.go`, `maintenance_vacuum_test.go`, and
`maintenance_integration_test.go`: each calls `fid init` with its return unused
only to establish a default FID workspace before exercising its own command.
They can use the existing real-bootstrap pristine clone, matching the rest of
the command suite. The assertion bodies and command-under-test invocations stay
unchanged. Init, rebuild, journal, and end-to-end test files keep their actual
initialization calls; tagged real-process journeys remain unchanged.

The byte-identical diagnostic command
`go test -race -cover ./cmd/fid -run '^TestMaintenanceVacuum' -count=3`
passed before and after this setup repair (both terminated `exit_code: 0`):
25.391 s / `22.3% of statements` before, 3.731 s / `21.4% of statements`
after. Diagnostic coverage differs because redundant initialization is no
longer exercised by these command tests; dedicated initialization tests remain.

## Ordinary gate attempt 3: PASS

After both causal setup repairs, the byte-identical required command
`go test -race -timeout 30m -coverprofile=coverage.out ./...` passed,
with final token `62.6% of statements` (terminated `exit_code: 0`).
Fresh CLI time is **527.099 s**, workflow **231.010 s**; every ordinary package
is below 600 s. Full output: `esc-23-002-evidence/race-attempt-3.log`.
`go tool cover -func=coverage.out` passed with total `65.6%`
(terminated `exit_code: 0`); full output is `esc-23-002-evidence/coverage.log`.
Format and lint after the final code changes passed: no format output and
`0 issues.`, both terminated `exit_code: 0`.

Attempts 1 and 2 remain failed wall-time-control evidence even though their
Go processes returned zero. The repair removes repeated bootstrap work from
fixture construction; it changes no assertion, command, test inventory, or
accepted runtime limit. Test cache was cleared before attempts 2 and 3, so the
CLI and workflow figures are fresh executions rather than cached timings.

### Additional smoke harness failures retained

The first invocation of the exact smoke script command terminated `exit_code: 1`:
its Issue-create child returned 2 because the harness omitted the required
`--idempotency-key`. After adding that key, the second invocation terminated
`exit_code: 1` with `sqlite3.DatabaseError`: Python's bundled SQLite 3.46.1
cannot parse the existing dynamic RAISE expression in the schema trigger.
Neither failure changed product code. The harness now obtains the Issue ID
from the actual create JSON and uses the already-installed SQLite CLI 3.50.6
for pragma checks. The byte-identical script invocation then printed the PASS
line above and terminated `exit_code: 0`. Both failed logs remain as
`rebuild-smoke-attempt-1.log` and `rebuild-smoke-attempt-2.log`.
The smoke harness requires Python 3 and the existing SQLite CLI with dynamic
RAISE-expression support (SQLite 3.47 or later); no package was installed.

## Final integration, benchmark and CLI verification

- Exact tagged command: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	669.586s` (terminated `exit_code: 0`),
  **669.586s < 900s**. Binary directory count **0 before / 0 after**;
  `/tmp` free bytes **3696472064 before /
  3696328704 after**. Raw evidence: `integration.log`
  and `integration-resources.jsonl` under `esc-23-002-evidence/`.
- Three consecutive fresh processes ran the byte-identical required 100k
  benchmark command, after the other heavy gates finished: **22.728s / 22.623s / 22.489s**.
  Every run printed `PASS` and terminated `exit_code: 0`; logs are
  `esc-23-002-evidence/benchmark-1.log` through `benchmark-3.log`.
- `go build -o /tmp/fid ./cmd/fid` returned no output and terminated
  `exit_code: 0`.
- `python3 project-docs/docs/phases/phase-23/qa/esc-23-002-evidence/rebuild-smoke.py`
  printed `PASS real CLI rebuild: read parity, unchanged journal, schema 17, FK/integrity clean`
  and terminated `exit_code: 0`. Its raw log records every actual CLI argv,
  return code and output, and the fresh workspace path.

ESC-23-002 is **RESOLVED 2026-09-20 under DEC-238**. The Human delegated
bounded repair authority above the fixed agent roles. Rebuild statement reuse
and fresh/isolated test setup repair the causes without changing any gate or
assertion. Fresh exact gates pass: CLI **527.099s**, workflow **231.010s**,
tagged integration **669.586s**; rebuild samples **22.728s / 22.623s / 22.489s**.
Integration binary directories stay **0 → 0**. Task 1 is complete in its report;
the Runner must reconcile the retained STOP before advancing. Tasks 9-11 stay
held. See [repair evidence](esc-23-002-repair.md) and
[resume handoff](resume-after-esc-23-002.md).

The 30-second assertion remains a useful reference-machine regression target;
it is not a universal wall-clock guarantee across hardware. This repair restores
headroom rather than accepting the schema-induced repeated preparation cost.
No later task was dispatched and the normal independent QA review is still due.

Evidence formatting note: the three benchmark logs omit only trailing spaces
from the CPU identification line; timings, result tokens and all other content
are retained. The historical escalation copy omits its extra trailing blank
line. This avoids whitespace-only commit defects without changing evidence.
