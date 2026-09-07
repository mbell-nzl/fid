# Phase 22 - Independent QA Lead Checkpoint After Task 12

Reviewer: independent QA Lead `checkpoint12_final`
Date: 2026-09-08
Status: **PASS**
Scope: the additional task-12 checkpoint required by QA plan sections 8 and 11.4, before task 13

```text
CHECKPOINT STATUS: PASS
Must-fix findings: 0
Blocker candidates: 0
Continuing task: task-13
```

The required fresh product checks terminated successfully after the orchestrator
verified task 12 COMPLETE and after Review4's independent tagged gate finished.
This checkpoint permits archive implementation to start. It is not the final
phase review and does not record Human judgments.

## Verified hand-back and independence

- Product: `c2d80980b9bfa9c1c20723cc7751526eb8baf2bf`,
  `feat(workflow): reconcile recovered execution targets`.
- Final completion report: `f57b4aa307fe5d3bb700ac04e8229070422b2d0e`.
- Developer's own final-gate addendum: `833b4fc2cefcdade9fd36c74e9b14573caa682fa`.
- Review4: `3addf7d6886b9c2b94ed8e5baeeff8c2dcbe1834`,
  `docs(qa): pass phase 22 review 4`; PASS, zero Must-fix and zero blockers.
  Its one improvement concerns report readability and is deferrable.

This reviewer has written no product code, tests, specifications, task reports
or workflow controls. The checkpoint product commands below were executed in
this reviewer's own session after COMPLETE; both use `-count=1`. Task-12
Developer/root gate receipts and Review4's full tagged receipt retain their
separate authorship. They are not re-attested as this reviewer's executions.

The task report and its addendum contain the exact five required commands in
the task specification's order, each with a closing output token and observed
termination. Earlier failures, superseded passes and the original checkpoint
remain attributed in the report/addenda. The disclosed test-order deviation is
retained; later behavioral tests are not represented as retrospective TDD.

## Own command receipts

All commands ran from `/home/mark/projects/fid`. Required task gates were
reconciled, not narrowed: these fresh focused commands implement this additional
checkpoint's named behavioral matrix. QA plan section 11.4 does not prescribe
a separate byte-exact command list for the checkpoint. Review4 independently
ran the unchanged full tagged gate before these commands.

### CP-01 — fresh ordinary recovery, disposition, budget and projection matrix

```bash
go test -race -count=1 ./internal/workflow ./internal/evidence ./cmd/fid -run '^(TestRecoverRun|TestAttentionResolution|TestEquivalentStops|TestRaiseStop|TestRoutineStatus|TestStopAndAttention|TestObserveRunBudget|TestFinalAllowedAttempt|TestBudget|TestDeriveBudget|TestResume|TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState|TestExecutionBound|TestTerminalSuccess|TestFIDBudget|TestArtifactLossHonorsRequiredMaterialAndPreservesProofIdentity|TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild|TestStopPredicates)'
```

PASS. Complete package output:

```text
ok  	github.com/mbell-nzl/fid/internal/workflow	41.065s
ok  	github.com/mbell-nzl/fid/internal/evidence	1.017s
ok  	github.com/mbell-nzl/fid/cmd/fid	192.410s
```

Observed termination: `exit_code: 0` (session 78472, final result chunk
259e6e). Earlier polls that printed package lines while the command was still
running were not treated as command completion.

### CP-02 — fresh built-binary recovery and activation journeys

```bash
go test -tags integration -race -count=1 ./cmd/fid/integration/... -run '^(TestP22|TestActivationRecoveryAtEveryDurableBoundary|TestActivationPrepareResponseLossReplaysDurableResult|TestPhase21RealProcessJourneyStartsOneActivatedEnvelope)'
```

PASS — `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	45.880s`
(terminated `exit_code: 0`; session 27575, final result chunk d90a4f).
The selector includes all Phase-22 recovery/activation/resume tables, the actual
activation interruption journey and the retained real-process final-budget
boundary journey. It exercises real child processes and disposable Git/SQLite
projects; none of these results are cache receipts.

### CP-03 — committed report reconciliation

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-22/tasks/task-12-done.md project-docs/docs/phases/phase-22/tasks/task-12.md
```

PASS — `PASS report-lint completed with no failures`
(terminated `exit_code: 0`; chunk 7005fa). Output included
`OUTCOME COMPLETE`, all five exact command matches, resolved product commit,
matching subject, and committed/clean report. Direct reading additionally
verified the gate addendum and preserved continuation history.

### CP-04 — source identity, isolation and CF-018 trigger

```bash
python3 - <<'PY'
import subprocess
product='c2d80980b9bfa9c1c20723cc7751526eb8baf2bf'
for ref in (product,'f57b4aa307fe5d3bb700ac04e8229070422b2d0e','833b4fc2cefcdade9fd36c74e9b14573caa682fa','3addf7d6886b9c2b94ed8e5baeeff8c2dcbe1834'):
 assert subprocess.check_output(['git','rev-parse','--verify',ref+'^{commit}'],text=True).strip()==ref
paths=['cmd','internal','migrations','go.mod','go.sum']
assert not subprocess.check_output(['git','diff','--name-only',product,'HEAD','--',*paths],text=True)
assert not subprocess.check_output(['git','status','--porcelain','--',*paths],text=True)
assert not subprocess.check_output(['git','diff','--name-only','89d964a84dbdd297fdf472d6f78b001e2efac702','HEAD','--','internal/sqlquery/sqlquery.go'],text=True)
print('CHECKPOINT IDENTITY PASS: four commits resolve; product unchanged and clean; CF-018 production-file trigger absent')
PY
```

PASS —
`CHECKPOINT IDENTITY PASS: four commits resolve; product unchanged and clean; CF-018 production-file trigger absent`
(terminated `exit_code: 0`; chunk a5631b). Source remained identical to the
verified task-12 product commit through both fresh runs. Only orchestrator-owned
documentation was dirty. CF-018's production-file-change condition did not fire;
this is not a new deferral.

## Required checkpoint matrix

| Check | Outcome | Own reproducible evidence and both-state coverage |
|---|---|---|
| Ordinary recovery | PASS | CP-01 `TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse` and CP-02 `TestP22DualAdapterRecoveryConformance` / real failure corpus: valid terminated recovery creates one successor, same-key response-loss replay preserves identity, stale/conflicting requests reject. |
| Known-live Human force recovery | PASS | CP-01 `TestRecoverRunForceRequiresCurrentHumanAuthorityAndFencesPriorRun`; CP-02 shared corpus and `TestP22RecoveryHumanResultMatchesDurableJSON`: current Human authority permits known-live takeover, ordinary/unauthorized paths reject, prior Run is fenced, result preserves force risk and provenance. |
| Unknown liveness | PASS | CP-01 `TestRecoverRunRejectsUnknownLivenessForOrdinaryAndForce` and CP-02 real failure corpus: both ordinary and force reject uncertainty without successor/effect creation. |
| Stop and Attention disposition | PASS | CP-01 `TestRecoverRunAppliesStopsBeforeOrdinaryOrForceMutation` covers applicable, unrelated-subject, unrelated-action and resolved states for ordinary/force recovery; `TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable` requires the exact separate disposition set. `TestResumeKeepsOpenAttentionAfterAllLinkedStopsAreResolved` proves zero unresolved Stops does not hide independently open Attention. Read-only and shared-Attention identity tests also execute. |
| Budget atomicity | PASS | CP-01 external-sample append/replay/rejection, unknown/exhausted admission, first-exhaustion boundary and committed-denial tests; CP-02 `TestPhase21RealProcessJourneyStartsOneActivatedEnvelope` executes the final permitted Attempt with its durable Budget/Stop/Attention Effects and rejects the next execution before spawn. The pure accounting helper is supplemental, not the sole atomicity proof. |
| Activation recovery, both directions | PASS | CP-02 `TestP22DualAdapterActivationConformance`, `TestActivationRecoveryAtEveryDurableBoundary` and `TestActivationPrepareResponseLossReplaysDurableResult`: pre-commit restore-prior, post-commit complete-forward, wrong direction/missing agreement rejection, actual response loss at prepare/first-stage/commit/finalize, same identity on retry. |
| Observed freshness, recovery and rebuild | PASS | CP-02 `TestP22EarlyRecoveryJourney` clean and dirty cases: actual canonical Git observation, unchanged proof through recovery/branch rename, fresh successor Gate/Receipt/Evidence, full Resume equality across actual rebuild, and one-byte same-HEAD staleness retaining original Evidence identity. |
| Import fault/environment feasibility | PASS — feasibility only | Own CP-ENV-01/02 receipts retained below: disk-backed sibling stage/target, real device identity, file/stage/parent fsync, complete rename publication, four synchronized SIGTERM boundaries and injected missing-stage/closed-fsync/occupied-target faults. These current-session preparation probes remain valid because the environment did not change; they do not claim task-14 product import behavior. |

## Scoped validation outcomes

| Assertion | Outcome | Evidence and practical boundary |
|---|---|---|
| VAL-048 / QA-49 | PASS | CP-02 clean/dirty `TestP22EarlyRecoveryJourney` observes the actual changed target and retains old Evidence ID. CP-01 `TestResumeProjectsExactActiveAndCompletedDurableStateWithoutRecovery/completed` proves actual terminal Issue transition and Envelope completion reject newly stale observed proof without Events, then succeed when the qualifying target is restored. |
| VAL-049 / QA-50 | PASS | CP-02 actual clean and dirty unchanged targets stay fresh through recovery, branch rename and rebuild. The A-to-B adopted-target cases additionally create new successor Gate/Receipt/Evidence within the same original Registration/Obligation and complete terminal operations; adoption does not refresh prior proof. CP-01 rejects foreign Project/worktree bindings. |
| VAL-050 / QA-51 | PASS | CP-01 `TestArtifactLossHonorsRequiredMaterialAndPreservesProofIdentity` runs metadata-sufficient and required-material cases using actual Artifact deletion and verification. Both retain complete Evidence/Receipt/Artifact metadata except availability; only required material produces the missing-Artifact policy rejection. Independent validity/freshness/multi-failure cases also execute. |
| VAL-051 / QA-52 | PASS | CP-02 fresh processes after interruption and complete Resume equality before/after real rebuild reconstruct bounded context. CP-01 exact active/completed projections and independently open Attention cover durable reads; CP-02 elapsed Claim test preserves recorded status/expiry while deriving effective status. Findings remains explicitly unavailable, with no invented completed Review support or stored next-action state. |
| VAL-052 / QA-53 | PASS | CP-02 `TestP22ResumeReadyCandidatesExecuteTheirActualOperations` actually executes run.inspect, run.observe, gate.run, claim.renew and claim.release at matching fixture revisions. The untrusted-producer case advertises a reason and rejects actual observation without Events; incomplete activation advertises gate denial and actually rejects Gate execution. CP-01 scoped Stop/authority/Run/Claim/revision/repository matrices and CP-02 adopted-target terminal journey retain qualifying and rejecting predicates. |

DEC-201 is reconciled by CP-01's derived-state, immutable-policy stripping,
closed authoring input (including empty/null state), unknown counters and
historical Budget projection/replay tests, together with CP-02's four public
Run replay paths. Available/warning/exhausted/unknown output states are derived
without rewriting policy, request digests or Event Effects; replay returns the
stored historical snapshot. The elapsed Claim case projects only the accepted
`effective_status` read member.

The lease matrix follows existing operation predicates: Claim renewal,
Issue transition, Run start and recovery retain their accepted guards. No new
Claim-expiry policy is imposed on `gate.run`; its accepted predicates are
active Run/Actor/Registration/repository/grant, Budget and Stop checks. The
advertised candidate set is verified against actual operations, not assumed to
be a fixed public list.

This is a scoped checkpoint outcome for VAL-048–052 plus supporting regression
evidence. It does not assign final outcomes to all phase rows. Final QA and
task-18 assembly remain due; VAL-076–078 require the Human's six separate
judgments. The neutral human-checklist draft has blank evidence/outcomes and
awaits the task-18 retained fixture manifest.

## Command and fixture feasibility

The report's exact commands and final result/termination fields were read and
matched independently (CP-03). Named test functions were located in the committed
ordinary/tagged files before selection; CP-01 and CP-02 produced actual package
results, not a no-tests result. The tagged harness builds and uses its binary
during each test process. Disposable test state is sufficient for this
deterministic checkpoint; task 18 separately owns durable Human fixtures and
post-generator inspection.

Archive feasibility and semantic-validation preflight below retain their
separate limitations. Task 13 still owes actual measured production ZIP limits,
valid complete controls, semantic rejection and no-effect checks; task 14 still
owes real staged publication/reconstruction/fault/idempotency proof. The accepted
boundary matrix's isolated limiter proof plus separate full-valid controls
interpretation is unchanged.

## Retained preparation history

Everything below this heading is historical preparation. Its PENDING labels
describe the state before task 12 was complete and are not this checkpoint's
current outcome. Prior-root-delegate receipts, their failed scratch attempts and
the inapplicable report-lint probe retain their original attribution and status.
Own preparation receipts remain identified by their independent reviewer
sections. No historical result has been upgraded to product PASS.

### Original preparation heading and status

Reviewer: qa-lead
Date: 2026-09-07
Status: PENDING task-12 completion
Checkpoint type: preparation evidence only; no post-task-12 product evidence yet

## Disposition

```text
CHECKPOINT STATUS: PENDING
Reason: task 12 is not complete, so the required independent recovery,
liveness, Stop, budget, activation, freshness, and built-binary checks have
not yet been run. Task 13 must not start until this artifact is updated with
fresh post-task-12 evidence and CHECKPOINT STATUS: PASS.
```

This preparation record is limited to the task-13/task-14 import fault-injection
and environment feasibility refresh requested by the QA plan. It does not claim
any `VAL-NNN` outcome. The later checkpoint update must retain these receipts,
then run the required fresh product checks from the current post-task-12 tree.

## Required Entry Condition Still Outstanding

The QA plan requires an independent checkpoint after task 12 and before task
13. The required post-task-12 checks are:

- ordinary recovery, known-live force recovery, and unknown-liveness recovery;
- Stop disposition;
- budget atomicity;
- activation recovery in both directions;
- freshness, including unchanged and one-byte changed repository targets;
- import fault injection and environment feasibility refresh.

The first five groups and task-12 report evidence are intentionally absent from
this preparation artifact. They must be run after task 12's completion report
and commit are independently verified.

## Preparation Evidence: Disk-Backed Import Environment

The repository's ignored `.cache` directory was used for disposable scratch
state. It is on the disk-backed `/dev/sdd` filesystem, rather than `/tmp`'s
tmpfs. This satisfies the storage placement prerequisite for an import staging
directory and destination that must remain siblings on one filesystem.

### Command receipts

| Check | Exact command | Result and final token | Termination |
|---|---|---|---|
| Filesystem and free-space inspection | `mkdir -p .cache && stat -f -c 'filesystem=%T device=%d block=%S' .cache && df -P .cache` | `filesystem=ext2/ext3 device=66604382 block=4096`; `/dev/sdd` reports `952852388` KiB available | `exit_code: 0` |
| Publication primitive and injected filesystem faults | `python3 -c "import errno, os, pathlib, tempfile; exec(\"def fails(fn, expected):\\n    try:\\n        fn()\\n    except OSError as e:\\n        return e.errno in expected if isinstance(expected, tuple) else e.errno == expected\\n    return False\"); root=pathlib.Path(tempfile.mkdtemp(prefix='qa-import-feasibility-', dir='.cache')); stage=root/'stage'; destination=root/'published'; stage.mkdir(); (stage/'state').write_bytes(b'complete-state'); fd=os.open(stage/'state', os.O_RDONLY); os.fsync(fd); os.close(fd); fd=os.open(stage, os.O_RDONLY|os.O_DIRECTORY); os.fsync(fd); os.close(fd); assert os.stat(stage).st_dev == os.stat(root).st_dev; assert not destination.exists(); os.rename(stage, destination); fd=os.open(root, os.O_RDONLY|os.O_DIRECTORY); os.fsync(fd); os.close(fd); assert (destination/'state').read_bytes() == b'complete-state'; missing=root/'missing-stage'; try_rename_missing=fails(lambda: os.rename(missing, root/'missing-published'), errno.ENOENT); fd=os.open(destination/'state', os.O_RDONLY); os.close(fd); fsync_closed=fails(lambda: os.fsync(fd), errno.EBADF); occupied=root/'occupied'; occupied.mkdir(); (occupied/'sentinel').write_bytes(b'keep'); stage2=root/'stage2'; stage2.mkdir(); (stage2/'state').write_bytes(b'new'); rename_occupied=fails(lambda: os.rename(stage2, occupied), (errno.EEXIST, errno.ENOTEMPTY)); assert try_rename_missing and fsync_closed and rename_occupied; print(f'import publication feasibility PASS: same-device sibling stage/target; file/stage/parent fsync succeeded; atomic rename yielded complete target; injected missing-stage errno={errno.ENOENT}, closed-fsync errno={errno.EBADF}, occupied-target preserved sentinel (terminated exit_code: 0); scratch={root}')"` | `import publication feasibility PASS: same-device sibling stage/target; file/stage/parent fsync succeeded; atomic rename yielded complete target; injected missing-stage errno=2, closed-fsync errno=9, occupied-target preserved sentinel (terminated exit_code: 0); scratch=/home/mark/projects/fid/.cache/qa-import-feasibility-wmf0ntgu` | `exit_code: 0` |
| Synchronized process fault injection | `python3 -c "import os, pathlib, signal, subprocess, tempfile, time; root=pathlib.Path(tempfile.mkdtemp(prefix='qa-import-fault-', dir='.cache')); stage=root/'stage'; destination=root/'published'; marker=stage/'ready'; stage.mkdir(); child=subprocess.Popen(['python3','-c',\"import pathlib,time,sys; p=pathlib.Path(sys.argv[1]); p.write_text('ready'); time.sleep(60)\", str(marker)]); deadline=time.monotonic()+5; observed=False; while_code=''; exec(\"while time.monotonic() < deadline:\\n    if marker.exists():\\n        observed=True\\n        break\\n    time.sleep(0.01)\"); assert observed; child.send_signal(signal.SIGTERM); status=child.wait(timeout=5); assert status < 0 and -status == signal.SIGTERM; assert not destination.exists() and marker.exists(); print(f'import fault injection feasibility PASS: synchronized readiness marker observed; SIGTERM status={status}; destination absent before publication; staged marker retained; scratch={root} (terminated exit_code: 0)')"` | `import fault injection feasibility PASS: synchronized readiness marker observed; SIGTERM status=-15; destination absent before publication; staged marker retained; scratch=/home/mark/projects/fid/.cache/qa-import-fault-dgsiv7fv (terminated exit_code: 0)` | `exit_code: 0` |
| Tracked-state isolation after probes | `git status --short --branch && find .cache -maxdepth 1 -type d -name 'qa-import-*' -printf '%f\\n' | sort | tail -n 10` | Only the active task-11 product paths were listed as tracked/untracked changes; scratch directories were under ignored `.cache`: `qa-import-fault-dgsiv7fv`, `qa-import-feasibility-wmf0ntgu` | `exit_code: 0` |

The publication probe proved the environment and primitive behavior only. It
did not exercise a future `export.import` implementation, journal reconstruction,
publication marker, target lock, or response-loss idempotency. Those remain
task-14 evidence obligations.

## Preserved Scratch Failures

These failures occurred before the corrected feasibility probes. They are
retained as preparation evidence and are not treated as product failures.

### Attempt 1: missing scratch parent

Exact command:

```text
stat -f -c 'filesystem=%T device=%d block=%S' .cache && df -P .cache
```

Observed output:

```text
stat: cannot read file system information for '.cache': No such file or directory
```

Observed termination: `exit_code: 1`.

### Attempt 2: invalid one-line Python exception syntax

Exact command:

```text
python3 -c 'import errno, os, pathlib, tempfile; root=pathlib.Path(tempfile.mkdtemp(prefix="qa-import-feasibility-", dir=".cache")); stage=root/"stage"; destination=root/"published"; stage.mkdir(); (stage/"state").write_bytes(b"complete-state"); fd=os.open(stage/"state", os.O_RDONLY); os.fsync(fd); os.close(fd); fd=os.open(stage, os.O_RDONLY|os.O_DIRECTORY); os.fsync(fd); os.close(fd); assert os.stat(stage).st_dev == os.stat(root).st_dev; assert not destination.exists(); os.rename(stage, destination); fd=os.open(root, os.O_RDONLY|os.O_DIRECTORY); os.fsync(fd); os.close(fd); assert (destination/"state").read_bytes() == b"complete-state"; missing=root/"missing-stage"; try_rename_missing=False; try: os.rename(missing, root/"missing-published"); except OSError as e: try_rename_missing=(e.errno == errno.ENOENT); fd=os.open(destination/"state", os.O_RDONLY); os.close(fd); fsync_closed=False; try: os.fsync(fd); except OSError as e: fsync_closed=(e.errno == errno.EBADF); occupied=root/"occupied"; occupied.mkdir(); (occupied/"sentinel").write_bytes(b"keep"); stage2=root/"stage2"; stage2.mkdir(); (stage2/"state").write_bytes(b"new"); rename_occupied=False; try: os.rename(stage2, occupied); except OSError as e: rename_occupied=(e.errno in (errno.EEXIST, errno.ENOTEMPTY)); assert (occupied/"sentinel").read_bytes() == b"keep"; assert try_rename_missing and fsync_closed and rename_occupied; print(f"import publication feasibility PASS: same-device sibling stage/target; file/stage/parent fsync succeeded; atomic rename yielded complete target; injected missing-stage errno={errno.ENOENT}, closed-fsync errno={errno.EBADF}, occupied-target preserved sentinel (terminated exit_code: 0); scratch={root}")'
```

Observed result: Python raised `SyntaxError: invalid syntax` at the inline
`try:` statement before the fault cases ran. Observed termination:
`exit_code: 1`.

### Attempt 3: invalid one-line helper definition

Exact command shape: the first corrected helper version ran `python3 -c` with
the publication probe followed by a literal `def fails(fn, expected):` after a
semicolon; the helper body contained escaped newlines and the same staging,
fsync, rename, and fault assertions.

Observed result: Python raised `SyntaxError: invalid syntax` at `def fails`.
Observed termination: `exit_code: 1`.

The corrected command receipts above use `exec()` for the multiline helper and
completed with observed `exit_code: 0`. No source, test, QA plan, task report,
specification, or ledger file was changed by the probes.

## Non-Applicable Control Probe

After committing this preparation artifact, I ran the Runner report-lint
control against it while checking the artifact boundary:

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-22/qa/checkpoint-after-task-12.md project-docs/docs/phases/phase-22/tasks/task-12.md
```

Observed output ended with `FAIL report-lint found 30 failure(s)` and included
the expected completion-report heading, outcome-field, task-command, and
commit-section failures. Observed termination: `exit_code: 30`.

This control is for Developer completion reports and is not applicable to a QA
checkpoint artifact. The result is preserved as a diagnostic receipt only; it
does not alter the checkpoint disposition or provide task-12 evidence.

## Later Update Contract

After task 12's completion report and commit are independently verified, append
to this artifact (preserving this preparation section) with:

1. exact post-task-12 commands and observed termination/results for each
   required checkpoint check;
2. VAL-048..052 evidence mapping, including both qualifying and rejecting states;
3. task-12 report-lint and commit verification;
4. a final `CHECKPOINT STATUS: PASS`, `REPAIR`, or `BLOCK` disposition;
5. a fresh commit containing only this checkpoint artifact.

Until those updates are complete, task 13 remains gated and no preparation
probe is a substitute for product evidence.

## Final Lines

```text
Preparation artifact written: project-docs/docs/phases/phase-22/qa/checkpoint-after-task-12.md
Task-12 completion: PENDING
Import feasibility: PASS as environment/primitive preparation only
Post-task-12 product evidence: NOT RUN
CHECKPOINT STATUS: PENDING
```


## Independent Import Feasibility Receipts — 2026-09-08

These commands ran in this independent checkpoint reviewer's session from
`/home/mark/projects/fid` while task 12 was still implementing. They prove only
environment and primitive feasibility. The orchestrator confirmed that these
fresh receipts need not be repeated after task 12 unless the environment changes.

### CP-ENV-01 — disk-backed scratch and same-filesystem identity

Exact command:

```bash
stat -f -c 'filesystem=%T block=%S' .cache /tmp && stat -c 'path=%n st_dev=%d' .cache /tmp && df -Pk .cache /tmp
```

PASS — observed `filesystem=ext2/ext3 block=4096`,
`path=.cache st_dev=2096`; `/tmp` is tmpfs with `st_dev=72`.
Disk-backed scratch had `952266204` KiB available; tmpfs had `1917552` KiB.
Final verbatim output token:
`tmpfs              4066620  2149068   1917552      53% /tmp`
(terminated `exit_code: 0`).

The filesystem-format `stat -f %d` field used in an earlier preparation command
counts free file nodes; it is not a device identifier. This independent receipt
uses ordinary `stat -c %d` and the probe below uses `os.stat().st_dev` for
actual device equality. The preceding historical receipt remains unmodified.

### CP-ENV-02 — synchronized interruption and filesystem faults

Exact command:

```bash
python3 - <<'PY'
import errno
import os
import pathlib
import selectors
import signal
import subprocess
import sys
import tempfile

root = pathlib.Path(tempfile.mkdtemp(prefix="qa-checkpoint12-import-", dir=".cache")).resolve()
child_source = """
import os, pathlib, sys
root, boundary = pathlib.Path(sys.argv[1]), sys.argv[2]
stage, target = root / "stage", root / "published"
stage.mkdir()
def barrier(name):
    if name == boundary:
        print(name, flush=True)
        os.read(0, 1)
(stage / "state").write_bytes(b"complete-state")
fd = os.open(stage / "state", os.O_RDONLY)
os.fsync(fd)
os.close(fd)
barrier("file-fsync")
fd = os.open(stage, os.O_RDONLY | os.O_DIRECTORY)
os.fsync(fd)
os.close(fd)
barrier("stage-fsync")
os.rename(stage, target)
barrier("published")
fd = os.open(root, os.O_RDONLY | os.O_DIRECTORY)
os.fsync(fd)
os.close(fd)
barrier("parent-fsync")
"""
for boundary in ("file-fsync", "stage-fsync", "published", "parent-fsync"):
    case = root / boundary
    case.mkdir()
    child = subprocess.Popen([sys.executable, "-c", child_source, str(case), boundary],
                             stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    ready = selectors.DefaultSelector()
    ready.register(child.stdout, selectors.EVENT_READ)
    assert ready.select(timeout=10), boundary
    assert child.stdout.readline().decode().strip() == boundary, boundary
    child.send_signal(signal.SIGTERM)
    assert child.wait(timeout=5) == -signal.SIGTERM, boundary
    assert child.stderr.read() == b"", boundary
    published = boundary in ("published", "parent-fsync")
    assert (case / "published").exists() == published, boundary
    materialized = case / ("published" if published else "stage")
    assert (materialized / "state").read_bytes() == b"complete-state", boundary
    assert materialized.stat().st_dev == case.stat().st_dev
    ready.close()
    child.stdin.close()
    child.stdout.close()
    child.stderr.close()
    print(f"{boundary}: synchronized SIGTERM observed; child returncode=-15; published={published}; complete bytes; same device")

def expect_error(operation, expected):
    try:
        operation()
    except OSError as failure:
        assert failure.errno in expected, failure
        return failure.errno
    raise AssertionError("fault did not reject")

missing = expect_error(lambda: os.rename(root / "missing", root / "missing-target"), (errno.ENOENT,))
fd = os.open(root, os.O_RDONLY | os.O_DIRECTORY)
os.close(fd)
closed = expect_error(lambda: os.fsync(fd), (errno.EBADF,))
occupied, stage = root / "occupied", root / "new-stage"
occupied.mkdir()
(occupied / "sentinel").write_bytes(b"keep")
stage.mkdir()
(stage / "state").write_bytes(b"new")
collision = expect_error(lambda: os.rename(stage, occupied), (errno.EEXIST, errno.ENOTEMPTY))
assert (occupied / "sentinel").read_bytes() == b"keep"
assert (stage / "state").read_bytes() == b"new"
print(f"fault errnos: missing-stage={missing}; closed-fsync={closed}; occupied-target={collision}; sentinel and stage retained")
print(f"IMPORT FEASIBILITY PASS: four synchronized interruption boundaries; scratch={root}")
PY
```

PASS — each child reached the named pipe barrier before SIGTERM; every child
terminated with the observed `returncode=-15`. The destination remained absent
after file/stage fsync and contained complete bytes after publication/parent
fsync. Each stage/target shared its parent device. Missing-stage rename returned
errno 2, closed-descriptor fsync errno 9, and occupied-target rename errno 39;
the occupied sentinel and rejected staging bytes remained intact.

Final verbatim token:
`IMPORT FEASIBILITY PASS: four synchronized interruption boundaries; scratch=/home/mark/projects/fid/.cache/qa-checkpoint12-import-7j7ms836`
(terminated `exit_code: 0`).

This exercises Python/OS primitives, not a fid import operation. Task 14 still
owns the real target lock, accepted Event reconstruction, semantic verification,
durability/publication faults, populated-target rejection, and same-key
response-loss retry. No project source or Git state was changed by the probe;
scratch is retained under ignored `.cache`.

### Semantic-validation feasibility for task 13

Task 13 expressly calls for existing Event appliers and no new dependency.
Contracts section 10 prohibits SQLite/SQL in the archive and mutation of the
source/destination project; it does not prohibit private in-memory SQLite during
validation. The existing exported applier is
`journal.ProjectEvent(*sql.DB, ProjectionInput)`. Sequential application can
retain references between Events; `DryRunProjectEvent` rolls back each Event
and cannot by itself validate a complete cross-Event history.

An implementation must create a real in-memory connection with the existing
modernc driver and accepted trusted schema, a controlled connection count,
foreign keys and recursive triggers. `internal/db.Open(":memory:")` is not an
in-memory constructor: its DSN helper applies `filepath.Abs`, yielding a disk
path. Never obtain executable SQL or schema definitions from an archive.
A private export-owned constructor can stay in task 13's current scope; adding
a new internal/db API needs an explicit bounded scope allocation.

ProjectEvent enforces SQL constraints, update preconditions and selected
immutability rules, but it is not the complete archive semantic validator.
`validateV2Event` currently accepts a nonempty Event type and `applyEffect`
validates identifiers/operations; task 13 must additionally enforce the accepted
closed Event/Effect grammar from Phase 22 schema section 5 and the retained
Phase 21 contract, including references not represented by foreign keys,
duplicate keys/IDs, ordering, and allowed tables/columns. A valid hash is not
semantic acceptance.

The current export validator compares capabilities to its binary's
`capabilitySnapshot()`; the accepted replacement is an explicit archive-version
profile (contracts section 10, public-interface section 5). Retain the Phase 21
v1 profile for read-only v1 validation and the accepted developmental v2 profile
during schema-15 work. Task 17 owns the final qualification update and exact
schema-16 capability profile. Never infer historical archive acceptance from
current discovery; intermediate pre-adoption work does not create a new legacy
compatibility promise.

These are source/contract feasibility findings, not executed product acceptance.
No declaration above replaces task 13's red/green semantic archive tests, full
valid controls or measured-byte boundary evidence. No contradiction requiring
new product policy was identified by this bounded preflight.


## Independent Checkpoint Preparation — 2026-09-08

Reviewer: independent QA Lead `checkpoint12_final`. The preceding preparation
receipts belong to the prior non-implementing root delegate's session. They
remain intact, including its failed scratch constructions and inapplicable
report-lint probe; they are not fresh product PASS evidence for this reviewer.
Task 12 is still implementing at this preparation point. Final product checks
will begin only after the orchestrator verifies its COMPLETE hand-back.

### Archive boundary interpretation for task 13

The accepted validation contract already distinguishes boundary admission from
whole-archive validity. Its `VAL-056 Boundary Acceptance Matrix` introduction
states: "Permitted means that limit check accepts; full archive acceptance still
requires every semantic check. Use valid complete controls as well". Task 13's
acceptance criteria likewise require full valid controls to accompany the
matrix. No contract amendment or relaxed limit is needed.

Consequently 056-a may use real ZIP files with 999999, 1000000 and 1000001 entries
to exercise the actual production count limiter, together with separate valid
complete archives through public validation. The large count fixtures need not
pretend to be semantically valid exports: manifest metadata can exceed its own
64 MiB bound before an otherwise valid million-entry export exists. Evidence
must identify the actual count-check acceptance/rejection, rather than infer it
from an unrelated error or repeat the arithmetic in test-only code. All three
numeric states remain required, and large fixtures remain generated per run.

For 056-d the measured production limiter must process actual decompressed bytes
at 8 GiB minus one, exactly 8 GiB and 8 GiB plus one, with every individual entry
within its own bound. A full-archive fixture includes manifest and other entry
bytes in that total; a focused limiter fixture can isolate the measurement and
use separate complete valid controls. Metadata-only rejection or an invalid ZIP
that never yields bytes is insufficient under the matrix's final paragraph.

This is a contract/feasibility interpretation, not evidence that the task-13
validator exists or passes. Product boundary evidence remains with task 13.
