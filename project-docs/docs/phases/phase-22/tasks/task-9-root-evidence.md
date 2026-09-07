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
| 2 | A | `0.836s` | `exit_code: 0` | Independent memory state passed the initial ten cases. |
| 3 | A | `FAIL` | `exit_code: 1` | Replay assertion found regenerated request timestamps/versions: `idempotency-conflict`, exit 4. |
| 4 | B | `1.167s` | `exit_code: 0` | Exact request caching and replay handling passed both adapters. |
| 5 | B | `FAIL` | `exit_code: 1` | Expanded 32-case corpus exposed missing fake operations, omitted required Stop reference and fixture directory lifetime tied to a completed subtest. Implemented independent operations, supplied the reference and retained fixtures at parent-test lifetime. |
| 6 | B | `FAIL` | `exit_code: 1` | Real stop.resolve and attention.resolve first deliveries succeeded but identical replays returned invalid-argument/2. DEC-197 assigned the causal persisted replay repair. |
| 7 | A | `1.736s` | `exit_code: 0` | Expanded independent fake passed. |
| 8 | B | `FAIL` | `exit_code: 1` | With wrong-fact and stale Attention cases added, real stop.raise failed during an incomplete typed-not-found transition. Developer repaired loadOpenAttentionByKey's typed absence handling. |
| 9 | B | `FAIL` | `exit_code: 1` | At 37 cases the sole remaining failure was committed denial replay: top-level data.replayed false while nested error.details.replayed true. Developer corrected the pre-start committed-denial frame builder. |
| 10 | A | `2.729s` | `exit_code: 0` | All final fake cases passed. |
| 11 | B | `2.568s` | `exit_code: 0` | All final shared cases passed against both adapters. |

The green receipts use a short verbatim time token. No failed
attempt was omitted to obtain a green result or used to narrow a test.

Final combined regression command:

- `go test -tags integration ./cmd/fid/integration/... -run '^(TestP22DualAdapterRecoveryConformance|TestP22RealBinaryRecoveryFailureCorpus)$' -count=1`
  - PASS — `4.476s` (terminated `exit_code: 0`).

The four root Go files were formatted with gofmt and goimports; each command
terminated `exit_code: 0` with no output. These focused receipts supplement
the unchanged five-command task gate; they do not substitute for it.

## Human rendering completion — root-owned follow-through

Before task closure root inspected the default Stop/Attention output against
public-interface section 6. The first Developer renderer handled singular
results only; list results still printed only their operation. Root took
exclusive ownership of the renderer and matching tests while the Developer
completed cross-surface Stop enforcement, then handed every file back before
the final five-command gate. The Developer's product commit includes this
source; these focused receipts remain root-attributed.

- `go test ./cmd/fid -run '^TestStopHumanListsShowOpenResolvedAndEmptyState$' -count=1`
  — first attempt `FAIL` (terminated `exit_code: 1`), because the new test
  declared json as a local flag rather than the required root persistent flag.
  Root corrected only that harness setup and reran the identical command.
- The identical command then produced the intended red `FAIL` (terminated
  `exit_code: 1`): open/resolved Stop lists and empty Stop/Attention lists
  printed only their operation, missing IDs/state.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22StopAttentionHumanListsNameReasonAndResolution$' -count=1`
  — intended red `FAIL` (terminated `exit_code: 1`): the built binary's
  `stop list` printed `stop.list` without the actual Stop ID.

Root added `cmd/fid/stop_human.go`, delegated Human output from
`writeStopSuccess`, and extended `cmd/fid/stop_test.go` plus the real-process
`cmd/fid/integration/p22_stop_human_test.go`. Lists now report open/resolved
state and empty results. Open Attention reads linked immutable Stop facts to
name the actual reason and qualifying condition. Resolved Stop text names its
fact and disposition without claiming it still blocks; Attention resolution
is correctly labelled as its own disposition Event. JSON remains unchanged.

- `go test ./cmd/fid -run '^TestStop(SuccessHumanNamesBlockingReasonAndQualifyingResolution|HumanListsShowOpenResolvedAndEmptyState)$' -count=1`
  — PASS — `0.076s` (terminated `exit_code: 0`).
- `go test ./cmd/fid -run '^TestStopHumanListsShowOpenResolvedAndEmptyState$' -count=1`
  — byte-identical rerun PASS — `0.076s` (terminated `exit_code: 0`).
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22StopAttentionHumanListsNameReasonAndResolution$' -count=1`
  — PASS — `1.194s` (terminated `exit_code: 0`). The real binary checks both
  Stop and Attention lists for the actual Stop/Run identities, intervention
  reason and `work-ended` condition.
- gofmt and goimports on the renderer and matching two test files terminated
  `exit_code: 0` with no output. No dependency changed.

Root explicitly required the Developer to rerun the unchanged five-command
gate on the final source after handback; earlier race receipts predating final
source changes remain historical evidence, not substitutes for the final gate.
