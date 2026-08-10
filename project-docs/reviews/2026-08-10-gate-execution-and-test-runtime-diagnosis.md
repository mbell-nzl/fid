# Gate Execution And Test Runtime Diagnosis

**Date:** 2026-08-10
**Branch:** `phase-20`
**Reviewed product baseline:** `7a09d98` (`f978b6d` product change plus its
completion-report commit)
**Scope:** Phase Runner and Developer gate interruptions involving
`go test -race -coverprofile=coverage.out ./...` and
`go test -tags integration -race ./...`, plus `cmd/fid` test runtime.

## Executive Summary

Two different failures occurred during Phase 20 and must remain distinct:

1. The earlier integration gate had a real product-test failure. Three
   binary-driven journeys initialized temporary directories whose derived
   prefix was invalid under the accepted project-prefix contract. Task 3c
   commit `f978b6d` repaired those stale fixtures by supplying an explicit
   valid prefix.
2. The later "terminated without output" result was not a Go test failure or
   a demonstrated timeout. The Developer discarded the live nested command's
   `session_id`, observed only the completion of its JavaScript wrapper, and
   incorrectly reported that the Go process had terminated.

Both exact gates pass on the Task 3c product baseline. They are slow because
the race detector instruments a large pure-Go SQLite workload, all `cmd/fid`
tests run serially, workspace initialization repeatedly migrates and rebuilds
the projection, and the integration gate reruns the complete untagged suite.

## Exact Gate Reproduction

The commands below were run literally from `/home/mark/projects/fid`. The
long-running sessions were polled until an exit code was observed.

| Command | Result | `cmd/fid` package time |
|---|---|---:|
| `go test -race -coverprofile=coverage.out ./...` | PASS | `131.228s` |
| `go test -tags integration -race ./...` | PASS | `162.472s` |
| `go test -count=1 ./cmd/fid` | PASS | `13.056s` |
| `go test -tags integration -count=1 ./cmd/fid` | PASS | `22.778s` |

The exact race-enabled commands emitted no bytes while `cmd/fid` was running.
Their package results arrived only after that first package completed. A silent
intermediate yield is therefore expected for more than two minutes on a cold
or invalidated cache; it is not evidence of termination.

## Developer Session Evidence

The Task 3d Developer session is recorded at:

`/home/mark/.codex/sessions/2026/08/09/rollout-2026-08-09T20-44-47-019fe5b2-2fa1-70f1-9214-fc5ed2f35e0f.jsonl`

The relevant event sequence is:

1. Line 173 invoked the coverage gate with `yield_time_ms: 30000`, then emitted
   only `text(r.output)` from the nested `exec_command` result.
2. Line 174 returned `Script running with cell ID 25`. This meant that the
   outer JavaScript execution cell was still active.
3. Line 178 correctly waited on cell 25.
4. Line 179 reported that the outer script completed with empty output. The
   nested `exec_command` had yielded a result containing a live `session_id`,
   but the wrapper printed only its empty `output` field and discarded all
   other result fields.
5. Lines 187-188 incorrectly concluded that the coverage command itself had
   terminated without output and triggered the required hard stop.

The same `text(r.output)` wrapper pattern appears in several earlier Phase 20
Developer sessions. This explains the intermittent behavior:

- A cached or focused command that finishes inside 30 seconds returns an exit
  code and terminal output through the initial call.
- A cold gate that exceeds 30 seconds returns a live command session. The
  wrapper discards its `session_id`, so the agent cannot poll it.
- An agent may then start another gate or retry while the first command is
  still running, creating avoidable resource contention until the abandoned
  command or agent session is cleaned up.

The correct long-command protocol is:

1. Serialize the complete nested tool result, not only `r.output`.
2. Wait when the outer tool returns a `cell_id`.
3. When the nested result returns a `session_id`, poll it with `write_stdin`.
4. Continue until an `exit_code` or signal is observed.
5. Treat empty intermediate output only as "no new bytes yet". It never proves
   success, failure, termination, or timeout.

For example, the wrapper must preserve metadata:

```javascript
const result = await tools.exec_command({
  cmd: "go test -race -coverprofile=coverage.out ./...",
  workdir: "/home/mark/projects/fid",
  yield_time_ms: 30000,
  max_output_tokens: 30000,
});
text(JSON.stringify(result));
```

If that result includes `session_id`, subsequent calls must use
`tools.write_stdin` and must likewise serialize the complete result.

## Runtime Decomposition

With the `integration` tag, `cmd/fid` contains:

- 164 top-level tests;
- 92 subtests;
- no `t.Parallel` calls;
- approximately 483 in-process CLI helper invocations;
- 16 integration-tagged top-level tests; and
- four independent integration-test builds of the `fid` binary.

A race-enabled JSON timing trace gave `159.77s` of summed top-level test time,
closely matching the package's `162.472s` result. This is accumulated serial
work, not one hung test.

The 16 integration-tagged tests accounted for approximately `29.9s`. The 148
ordinary tests accounted for approximately `129.9s` and were rerun because
`go test -tags integration -race ./...` includes both tagged and untagged test
files. About 81% of the integration gate therefore duplicates the ordinary
suite already exercised by the preceding race-and-coverage gate.

The largest race-enabled top-level test times were:

| Test | Elapsed |
|---|---:|
| `TestWorkflowDefineRejectsWithoutMutatingJournalOrWorkflowTables` | `7.20s` |
| `TestV01ClosureJourneyAndRebuildParity` | `6.12s` |
| `TestWorkflowDefineNormalizesArrowsAndAcceptsAppliesToVocabulary` | `5.13s` |
| `TestV01ReadCommandsSupportAllFormatsWithStructuralParity` | `3.95s` |
| `TestLinkAddImplementsRequirementAndRejectsInvalidInputBeforeAppend` | `3.84s` |
| `TestIssueACVerifyCommandRejectsBeforeAppend` | `3.84s` |
| `TestIssueCreateRejectsInvalidFlagsAndConstraintViolations` | `3.44s` |

There are four fixed sleeps totaling 7.3 seconds: three 1.1-second claim-expiry
sleeps in the ordinary suite and one four-second maintenance sleep in the
closure integration journey.

## CPU Profile

The diagnostic command was:

```bash
go test -count=1 -cpuprofile=/tmp/fid-cmd.cpu ./cmd/fid
go tool pprof -top -cum -focus='github.com/mbell-nzl/fid' /tmp/fid-cmd.cpu
```

The non-race profile reported these cumulative paths. Percentages overlap and
must not be added:

| Path | Cumulative CPU share |
|---|---:|
| Workspace initialization | `62.60%` |
| Database migration | `55.87%` |
| SQLite VM execution | `44.22%` |
| Staged rebuild verification | `33.00%` |
| Linux syscalls | `29.60%` flat |

`InitializeWorkspace` creates a staging workspace, applies all migrations,
initializes journaled state, performs a real rebuild, compares the projection
bytes, syncs every staged file and directory, renames the staged workspace,
and syncs the workspace root. This durability path is appropriate production
behavior, but ordinary command tests pay for it repeatedly.

The test helper also changes the process-wide working directory around every
in-process CLI invocation. That global mutation prevents safe parallel tests
even though each test otherwise owns a temporary workspace.

## Root Causes

Ranked by immediate operational impact:

1. **Lost command-session metadata.** Agent wrappers print only `r.output`, so
   they cannot continue a command that yields a `session_id`.
2. **Duplicate phase-gate work.** The integration command reruns the entire
   untagged suite after the race-and-coverage command.
3. **Race amplification of pure-Go SQLite.** `modernc.org/sqlite`, migrations,
   rebuilds, and filesystem-heavy initialization are instrumented throughout.
4. **Serial test architecture.** Global working-directory changes prevent
   `t.Parallel`, so 164 top-level tests accumulate directly into wall time.
5. **Repeated expensive fixtures.** Tests not concerned with initialization
   repeatedly execute the full production initialization and rebuild path.
6. **Integration harness repetition.** Four tests build the same binary
   independently and integration scenarios execute many short-lived processes.
7. **Wall-clock sleeps.** Expiry and abandoned-attempt behavior waits 7.3
   seconds instead of using a controllable clock.

## Recommended Remediation

### 1. Correct Long-Command Handling

Update the canonical Developer and every gate-running QA role prompt with the
complete-result and `session_id` polling protocol. Regenerate provider prompts
and run workflow lint. This is the immediate fix for false blockers and does
not change any product or gate contract.

The Phase Runner itself should also reject a blocker claim that says a command
"terminated" unless the attached evidence contains an exit code or signal.
Completion of an outer execution cell is not command termination evidence.

### 2. Separate Integration Execution

Move integration journeys to a dedicated Go package or directory and revise
the accepted integration gate to target that package. The race-and-coverage
gate should continue to own all ordinary packages; the integration gate should
own only the tagged real-command journeys.

This changes the recorded testing strategy and exact quality command, so it
requires an accepted decision plus updates to `AGENTS.md`, the Quality Plan,
AI Workflow, templates, and applicable QA plans/task specs. It must not be
introduced as an ad hoc narrower substitute during an active task.

### 3. Introduce A Reusable Workspace Fixture

Initialize and migrate a pristine workspace once, then clone it into each
temporary test workspace for command tests whose subject is not initialization.
Keep the full production path in `fid init`, bootstrap, rebuild, durability,
and compatibility tests. The fixture must preserve journal/DB parity and must
not become a mock of accepted persistence behavior.

### 4. Remove Process-Wide Working-Directory Mutation

Inject the workspace/project-location boundary into the root command or its
test harness so in-process command tests do not call `os.Chdir`. Once global
state is removed, add bounded `t.Parallel` coverage and verify that SQLite and
filesystem load does not introduce contention or flakiness.

### 5. Reuse The Integration Binary

Build `fid` once for the integration package, using `TestMain` or an equivalent
package-owned fixture, and share the immutable binary across isolated
workspaces. Preserve real-process execution and all scenario assertions.

### 6. Replace Sleeps With A Clock Boundary

Inject a clock into claim expiry and maintenance-age calculations. Unit and
command tests should advance deterministic time; retain one appropriately
bounded real-clock integration check only if the contract needs it.

### 7. Benchmark Production Initialization Separately

Add a focused benchmark for `InitializeWorkspace` and record migration,
rebuild-verification, and sync costs. Optimize the production durability path
only against its accepted atomicity and crash-safety contract. Test-gate speed
alone is not authority to remove rebuild verification or filesystem syncs.

## Delivery Constraints

These remediations should not be combined into one Developer session:

- The prompt/session fix is delivery tooling and generated-agent scope.
- The integration-package and exact-command change is a quality-strategy and
  specification decision.
- The reusable fixture and clock are testability changes in separate owning
  boundaries.
- Removing global `chdir` affects the CLI/project-location boundary and needs
  focused regression evidence.
- Production initialization optimization, if pursued, is a distinct durability
  and performance task.

The active Phase 20 Task 3d blocker and partial working-tree changes must be
resolved before normal Phase Runner execution resumes. The remediation plan
must preserve those changes and follow the recorded escalation disposition.

## Conclusion

The gates are currently green but operationally fragile. The immediate stop is
caused by incorrect long-command session handling, while the long runtime is a
real test-architecture cost dominated by repeated initialization, pure-Go
SQLite race instrumentation, serial execution, and duplicate integration-gate
coverage. The session-handling repair is small and urgent. Runtime work should
then proceed as several measured tasks without weakening durability or evidence
contracts.
