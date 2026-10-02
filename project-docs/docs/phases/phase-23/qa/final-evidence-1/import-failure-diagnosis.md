# Full-checkptr import failure — fresh QA diagnosis

Reviewed HEAD: `32a55a94efb2a89dbf50b29accb434b9c753418c`.

## Observation and reproduction

From `/home/mark/projects/fid`, the required command
`go test -race -timeout 30m ./...` terminated with `exit_code: 1`.
The failed test was `TestImportRequiresFreshDestinationObservationBeforeRecovery`:

```text
--- FAIL: TestImportRequiresFreshDestinationObservationBeforeRecovery (109.65s)
    import_test.go:263: release imported Claim: claim expired
FAIL
FAIL	github.com/mbell-nzl/fid/internal/export	357.899s
```

Complete result: [full-checkptr.json](full-checkptr.json). The command was not
retried. No source or test was changed. The `cmd/fid` package passed in
575.300s; both later live-golden selections passed, so this is separate from
the DEC-267 golden-child cause.

## Expected behavior and narrow causal path

The test must reject imported historical liveness for close/abandon/recovery
and reject a Gate bound to the source repository. After a fresh destination
observation and a valid Claim release, it must recover with destination
bindings, execute a successor Gate, and preserve idempotent import replay.

`internal/export/import_test.go:189` obtains its active Claim from
`addJournalBackedExportFixtureFacts`. The helper creates it with
`TTLSeconds: 60` in `internal/export/archive_fixture_test.go:81`.
The test exports and imports that history, performs several checks, then calls
`claims.ReleaseGeneric` at `internal/export/import_test.go:262`. That call
reported expiration; the positive destination recovery leg was not reached.

## Ranked hypotheses and evidence limits

1. **Incidental fixture lease expires under the full-checkptr workload.**
   Supported by the observed expiration, the real 60-second lease, and the
   109.65-second failing test. A Developer should establish this deterministically
   with a controlled fixture-age regression and preserve all import/expiry
   behavior. Exact lease age was not traced in this QA run.
2. **Product recovery or import semantics are wrong.** The failure is before
   the positive recovery request, in ordinary Claim expiry enforcement. The
   seven-command gate passed this test and fresh Phase-23 round-trip checks
   are separately required. This run does not establish a product defect or
   justify changing expiry rules.

The inherited test came from `cc73b092bf5ecd3fdba62a1362062a6d543697ab`
(phase-22/task-14); the 60-second source fixture lease came from
`1be8e882ef689f15705507c448e3b087e23d5f65` (phase-22/task-13).
Both predate Phase-23 base `0195897a347dfe1872db0f7be3800c5c3ccdefcf`, as
shown by [provenance](import-failure-provenance.json). Inheritance does not
waive the current gate or turn bounded work into a protected blocker.

## Ownership and repair boundary

Phase-23 task 16 owns import binding and matching in-package tests. Limit
repair to this import fixture's incidental live-Claim premise in
`internal/export/import_test.go` and the narrowly needed
`internal/export/archive_fixture_test.go` helper, matching regression evidence,
and `task-16-done.md`. Preserve the negative imported-history checks, positive
destination recovery, repository binding, successor Gate and retry assertions.
Preserve product expiry enforcement, test parallelism, exact commands,
race/checkptr flags, timeouts and budgets. Do not mine an isolated green.

This is **REPAIR**, no VAL, mandatory QA plan §22 control. A fresh independent
QA Lead must rerun the failed byte-identical full command and the affected
phase evidence after causal repair, then write the next unused final review.
The full-gate prerequisite for the original final review's MF-1 remains open;
its specific DEC-267 code repair is verified. No new TRG or ESC is requested.
