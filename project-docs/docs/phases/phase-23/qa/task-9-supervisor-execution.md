# Task 9 supervisor execution note

2026-09-21. Root owns this record; these are root-observed results, not
Developer-session gate attestations or an independent QA review.

## Sequence-control repair

`project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-23/run-ledger.md`
initially terminated `exit_code: 1` with
`INVALID unrecognized or malformed event heading` because root included `UTC`
in two event headings. Root used the accepted local timestamp format and reran
the byte-identical command: `VALID state=REPAIR_PENDING` (terminated
`exit_code: 0`). No command or control semantics changed.

## Unnecessary duplicate ordinary gate

Working directory: `/home/mark/projects/fid`.

Root incorrectly inferred a stale execution session from empty output and an
empty process inventory. Tool receipts only established that the Developer's
command was pending. Root then launched the same heavy command while the
Developer's original command was still pending, contrary to the serial-heavy-
gate rule. This was an orchestration error, not a Developer decision. Root
corrected the inference and instructed the Developer to retain the original
session and observe termination.

- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Root session `14205`, launch receipt `6fa558`, final receipt `5107ce`.

The root execution used the version before the claims fixture repair. Its
observed failure was:

```text
--- FAIL: TestClaimAcquireRejectsOpenReviewReworkStop (0.05s)
    generic_test.go:49: claimAcquireStopBlock() = (*workflow.StopBlockedError)(nil), evaluate Stops: SQL logic error: no such table: envelope_protected_decisions (1)
FAIL
```

The final result also recorded `cmd/fid 802.292s` (over the 600s ordinary-package
budget), `internal/workflow 431.786s`, and `internal/export 491.699s`. Overlap
confounds attribution; this run does not establish a task-caused regression.
The full result ended `FAIL`. Root sent a Ctrl-C character on the final poll,
but the tool reported `exit_code: 1` and complete package results, with no signal
or interrupted outcome; no different termination is asserted.

No result from this execution is used as Developer-session PASS evidence. The
duplicate is closed. The Developer owns the causal fixture repair and the
next byte-identical required gate, run without another root gate in parallel.
Prior failures and timing observations remain retained regardless of the later
result. No threshold, test assertion, or required command has been relaxed.
