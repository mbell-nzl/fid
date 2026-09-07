# Task 9 — delegate-owned semantic corpus evidence

Root collected these receipts during the task-9 session in
`/home/mark/projects/fid`. This is separately attributed evidence for DEC-196;
the Developer runs and records the required task gates independently. It does
not attest to those gates or Human validation.

## Change and provenance

The task-6 predecessor is `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`.
Its response-only fake has been replaced with independently owned in-memory
Run, Claim, observation, recovery, budget, Stop and Attention state. The fake
consumes built operation requests rather than expected fixture results. Shared
request construction is IO-free; real execution uses the built fid binary.
Exact requests are cached for replay, including timestamps and versions.

The 37 shared cases assert typed values, independently specified versions and
counters, recovery lineage, immutable result/Event identity on same-key replay,
and state after rejected operations. Stop/Attention cases cover disposition
links and stale/nonqualifying resolution. External budget cases preserve
unknown values and derive a baseline plus delta without admitting a rejected
stale sample. A committed unknown-budget denial asserts one terminal frame,
no Attempt/Receipt/start/output, stable Stop/Attention/Event identities and
both replay indicators. This fake gate operation is explicitly limited to
unknown-budget denial; it does not simulate arbitrary process execution.

Root owns the four semantic adapter/helper test files and the shared fixture
directory. No runtime dependency changed. Product repairs exposed by this
corpus were assigned to the Developer under DEC-197, with predecessor
provenance retained there and in the task report.

## Focused command receipts

All commands below ran from `/home/mark/projects/fid`. Each termination was
observed from the command result in this same root session. Repeated commands
are listed in chronological order; failures remain evidence.

Command A:
`go test -tags integration ./cmd/fid/integration/... -run '^TestP22DualAdapterRecoveryConformance$/fake$' -count=1`

Command B:
`go test -tags integration ./cmd/fid/integration/... -run '^TestP22DualAdapterRecoveryConformance$' -count=1`

| Attempt | Command | Observed token | Observed termination | Finding or repair |
|---|---|---|---|---|
| 1 | A | `FAIL` | `exit_code: 1` | Strong typed assertion rejected `run.observe first delivery replayed = "replayed-fake", want boolean false`; run inspection likewise returned a string instead of state. |
| 2 | A | `ok  github.com/mbell-nzl/fid/cmd/fid/integration 0.836s` | `exit_code: 0` | Independent memory state passed the initial ten cases. |
| 3 | A | `FAIL` | `exit_code: 1` | Replay assertion found regenerated request timestamps/versions: `idempotency-conflict`, exit 4. |
| 4 | B | `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.167s` | `exit_code: 0` | Exact request caching and replay handling passed both adapters. |
| 5 | B | `FAIL` | `exit_code: 1` | Expanded 32-case corpus exposed missing fake operations, omitted required Stop reference and fixture directory lifetime tied to a completed subtest. Implemented independent operations, supplied the reference and retained fixtures at parent-test lifetime. |
| 6 | B | `FAIL` | `exit_code: 1` | Real stop.resolve and attention.resolve first deliveries succeeded but identical replays returned invalid-argument/2. DEC-197 assigned the causal persisted replay repair. |
| 7 | A | `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.736s` | `exit_code: 0` | Expanded independent fake passed. |
| 8 | B | `FAIL` | `exit_code: 1` | With wrong-fact and stale Attention cases added, real stop.raise failed during an incomplete typed-not-found transition. Developer repaired loadOpenAttentionByKey's typed absence handling. |
| 9 | B | `FAIL` | `exit_code: 1` | At 37 cases the sole remaining failure was committed denial replay: top-level data.replayed false while nested error.details.replayed true. Developer corrected the pre-start committed-denial frame builder. |
| 10 | A | `2.729s` | `exit_code: 0` | All final fake cases passed. |
| 11 | B | `2.568s` | `exit_code: 0` | All final shared cases passed against both adapters. |

The earlier green rows reproduce the package/result/time content with spacing
normalized; the final receipts use a short verbatim time token. No failed
attempt was omitted to obtain a green result or used to narrow a test.

Final combined regression command:

- `go test -tags integration ./cmd/fid/integration/... -run '^(TestP22DualAdapterRecoveryConformance|TestP22RealBinaryRecoveryFailureCorpus)$' -count=1`
  - PASS — `4.476s` (terminated `exit_code: 0`).

The four root Go files were formatted with gofmt and goimports; each command
terminated `exit_code: 0` with no output. These focused receipts supplement
the unchanged five-command task gate; they do not substitute for it.
