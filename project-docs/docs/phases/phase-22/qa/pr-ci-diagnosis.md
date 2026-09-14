# Phase 22 PR CI test timeout — diagnosis and proposed disposition

Date: 2026-09-13. Status: **CI failure remains open; proposed gate-budget change
awaits Human approval (ESC-22-005).** No product, test, Makefile or CI workflow edit performed.
Human checklist acceptance remains valid; successful PR CI is still outstanding.

## Failure and expected behavior

The Human supplied [pr-CI-test-failure-log.md](pr-CI-test-failure-log.md).
`make test` runs `go test -race ./...`. The supplied output records
`panic: test timed out after 10m0s`, export package failure at `600.033s`,
Makefile failure and Actions process exit 2. These are attributed CI receipts,
not a fresh local execution of the whole gate. The excerpt also shows the CLI
package timing out while its Stop matrix is running. It contains no failing
assertion or demonstrated deadlock; omitted parts of the complete job log have
not been inspected.

The timeout applies to each package's test binary, not to each named test.
The Stop matrix's displayed 33s and current subtest's 7s are not evidence of a
single ten-minute hang. Go 1.26.2's local `go help testflag` explicitly states
that `-timeout` defaults to 10 minutes. Neither Makefile target sets it, and
.github/workflows/ci.yml invokes those targets without an override. This is
Go's default timeout, not an explicit Actions job deadline.

## Hypotheses and observations

1. **Aggregate race-suite CPU/fixture cost exceeds CI's package budget:**
   supported. The export stack is runnable inside SQLite migration parsing;
   other packages continue to finish. The CLI's parallel tests waiting in
   testing.T.Parallel are normal scheduling evidence, not a demonstrated lock
   cycle. Focused profiles show considerable SQLite schema parsing/setup work.
2. **One named test hangs or fails its contract:** not reproduced. Both named
   tests complete with all their existing assertions under race instrumentation.
   Isolated success cannot prove that the full package never hangs.
3. **The 8GiB archive fixture still runs routinely:** excluded by the source.
   archive_scale_test.go has the archive_scale build constraint; make test
   supplies no such tag. DEC-217's real-scale matrix remains separately gated.

The local machine and hosted CI environment differ; no equivalence of their
CPU allocation or contention is assumed. These runs diagnose the visible paths;
they do not reproduce or clear the complete remote failure.

## Own diagnostic evidence

Working directory: /home/mark/projects/fid. Source: f323479, with only the
Human-supplied untracked failure log present before diagnosis. No source changes.
Profiles and test binaries are retained in .cache/phase22-pr-ci/ locally.

- `go test -race ./cmd/fid -run '^TestStopPredicatesAcrossRunStartAndGateRunCommands$' -count=1 -cpuprofile .cache/phase22-pr-ci/stop-before.cpu -o .cache/phase22-pr-ci/cmd.test` — PASS — `43.531s` (terminated `exit_code: 0`)
- `go test -race ./internal/export -run '^TestImportAfterPublicationFaultReplaysInstalledResult$' -count=1 -cpuprofile .cache/phase22-pr-ci/export-before.cpu -o .cache/phase22-pr-ci/export.test` — PASS — `11.066s` (terminated `exit_code: 0`)

Both commands ran serially. The CLI profile contains 41.58s sampled CPU:
16.98s (40.84%) under sqlite3Prepare, 15.48s under sqlite3Init, and 16.11s under
cloneReusableWorkspaceFixture. These are overlapping cumulative costs, not
additive. The export profile contains 10.11s sampled CPU: 3.80s (37.59%) under
sqlite3RunParser and 3.19s under migrate.runMigrations, also overlapping.
These are small operational fixtures; schema/trigger preparation and repeated
real workflow operations are costly even without a large Issue dataset.

The CLI matrix embeds a different absolute marker path in each Gate Definition,
so those cases cannot reuse the same existing definition-keyed fixture. Export
source fixtures also repeatedly bootstrap; validation/import deliberately exercise
real reconstruction. Optimizing these paths is possible but broader than a proven
single-test defect and needs isolation/provenance/rebuild guarantees retained.
A fixture rewrite is not yet shown sufficient to put the complete hosted suite
below ten minutes, so it is not presented as a verified CI fix.

## Concrete recommended change — not applied

Use an explicit **20-minute per-package timeout** for both Makefile targets
executed by CI. This doubles the inherited default as an initial operational
budget; it does not establish a measured optimal timeout or promise a CI PASS.
Keep the full package selection, every test, race detection and coverage intact.

```diff
 test:
-\tgo test -race ./...
+\tgo test -race -timeout=20m ./...

 coverage:
-\tgo test -race -coverprofile=coverage.out ./...
+\tgo test -race -timeout=20m -coverprofile=coverage.out ./...
 \tgo tool cover -func=coverage.out
```

The diff uses visible \t for Make recipe tabs. Approval would be recorded in a
phase decision and the quality plan, limited to these Makefile/CI suite budgets.
No archive limit, 100k rebuild performance threshold, product operation timeout,
assertion or test case changes. Historical task commands/receipts remain intact.
The bounded timeout still catches tests that cannot finish; it does not make
individual tests run longer if they already complete.

After approval: apply the two-line Makefile change, record its authority, run
both exact Makefile targets serially with observed termination, and report the
result. Actual hosted CI rerun remains necessary to close the reported failure.
Any further failure requires diagnosis, not another unexamined budget increase.

AGENTS.md's protected gate/control rule (lines 234..238) says a change to a
command or threshold/budget/tolerance needs the authority owning the definition.
The earlier delegated repair authority and DEC-217 do not specifically approve
this different timeout-budget change. Human approval is therefore requested
before applying the concrete change above. No new product decision is involved.

## Revised recommendation after Human runtime concern

The Human questioned the acceptability of 30–40 minute CI feedback and asked
whether runtime can be reduced without compromising quality. Timeout expansion
is not approved. Root recommends reducing duplicated work before reconsidering
the timeout; the earlier proposal above remains historical, not accepted authority.

Concrete findings from Makefile and .github/workflows/ci.yml:

- CI runs make test and make coverage sequentially. Both execute every ordinary
  package test with race detection; coverage adds instrumentation/reporting. One
  race-and-coverage execution retains those checks and removes the duplicate pass.
  This reduces total work but does not by itself cure a single-package timeout.
- The Stop-operation matrix changes an absolute marker path in every Gate
  Definition. The existing reusable fixture key includes that definition, so all
  ten cases rebuild their otherwise equivalent setup. A stable relative marker
  with a verified per-clone working directory is a candidate; blocked/permitted
  subprocess marker assertions and independent fixture mutations must remain.
- Export matrices repeatedly initialize equivalent source workspaces and build
  genuine history. Reusing an immutable source/archive with distinct writable
  copies can remove setup work. Changed-digest cases must still use genuinely
  different archives; destination import/rebuild operations under test stay real.
- Package jobs can reduce contention and overlap lint/build with testing, but
  heavy single-package suites may still need exhaustive test sharding. Any shards
  must prove disjoint selection whose union covers every discovered test, retain
  race instrumentation, and merge coverage before reporting. This can lower wall
  time without lowering total CPU and should follow actual timings, not guessing.

Prior demonstrated evidence: task-1-done.md's integration runtime repair reduced
one force-recovery matrix from 89.800s to 35.396s by cloning an immutable migrated/
replayed prefix. All case-specific writers, independent observers, full snapshots,
faults and real rebuild controls remained. That group-level result supports the
approach; it is not a measured prediction for the full CI suite.

Suggested bounded execution: establish uncached package/test timings, optimize
measured fixture hotspots with isolation regressions, consolidate the duplicated
CI execution under recorded gate authority, then measure full CI and shard only
if necessary. No assertion, race check or ordinary behavioral test is proposed
for deletion or deferral. No changes in this section have been implemented.

## Accepted repair and corrected CI timeline — 2026-09-14

DEC-221/task 20 selects performance repair without increasing the timeout.
GitHub run 34744340253/job 103689338838 reports the Test step from
07:07:57Z to 07:39:30Z (31m33s); the job ran 34m35s. CLI timed out at
07:19:02Z (605.208s), export at 07:21:14Z (600.033s), and workflow finally
passed at 07:39:30Z (602.082s). Queued package work continued after early
failures. The timestamps do not support treating the excess as twenty minutes
of initial compilation. Coverage and build were skipped. Full raw log is
.cache/phase22-pr-performance/original-ci.log.

Local matching full-suite before/after runs passed: CLI 541.652s -> 483.994s,
export 368.615s -> 264.146s, workflow 503.550s -> 464.948s. All 2,653 original
executed cases remain, plus two new isolation checks. Local JSON event spans
616.574s -> 574.741s are not total process durations. Setup savings alone do
not establish sufficient hosted headroom; four exhaustive workers therefore
share the full race/coverage execution under DEC-221. Final local coverage is
65.5%; independent review and all task gates pass. Hosted timing remains
unverified until the revised PR run completes. See ../tasks/task-20-done.md.
