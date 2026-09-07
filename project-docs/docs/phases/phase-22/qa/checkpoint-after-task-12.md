# Phase 22 - Independent QA Lead Checkpoint After Task 12

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
