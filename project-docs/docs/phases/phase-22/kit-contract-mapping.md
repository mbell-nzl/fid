# Phase 22 Kit Contract Dispositions

> Tech Lead, 2026-09-05; DEC-187, accepted by Human under DEC-188.
> Independent source: `qa/consistency-review.md`, authored by QA Lead.
> Kit pin: `c5229617220f97ff9228df261e24718b1d1d6d89`.

This is the disposition of that review, not new independent QA evidence.
Phase 22 qualifies the documented fid binary surface. Production kit adapter
qualification still requires all twelve capabilities and the Phase 25 proof.
No external kit file is changed or adapter implementation claimed here.

## Decisions Requiring Explicit Disposition

### C-06b: Coordinator-Owned Recovery Export Binding

Accept an explicit boundary difference. The kit coordinator owns the durable
activation transaction binding its recovery export to the fid activation ID.
fid records component agreement, versions, and direction; it does not store an
archive reference in ActivationRecovery. `export.portable.v1` provides archive
creation/validation/import, but availability alone does not satisfy the binding.

Before destructive materialization, the coordinator must durably bind the
validated archive's identity and digest, the prior activation/version, the new
fid activation ID, and intended old/new component identities. Crash recovery
must select that same export and activation; missing/mismatched export binding
must prevent destructive continuation. First activation records null prior
activation identity but still binds an export of the existing standalone fid
project. Absence of a prior activation does not imply absence of project data
or waive its recovery export.

This coordinator proof belongs to the kit V2-3 adapter and fid Phase 25
qualification (CF-037), with present/missing/mismatched binding and crash/retry
cases. It is an accepted Phase 22 exclusion, not completed capability evidence.
Task 10 implements the fid boundary; task 18 and final QA verify that exports
are not advertised as proving the coordinator obligation. No Phase 22 DDL,
Event, dependency, or activation input extension is needed for this disposition.

### C-08: Derived Terminal State

Retain `run.close` and `run.abandon`, and remove `run.finish` including its
grant action. A kit request for succeeded/failed maps to `run.close`; its
`terminalState` is only an expected result checked against fid's derived truth.
It never supplies an outcome to fid. A mismatch is an adapter contract failure,
not permission to rewrite or hide the durable result. Replay returns that same
result. Intentional abandoned maps to `run.abandon` only with its required
cessation/liveness/authority inputs; a terminalState string alone grants none.
Interrupted is produced by linked recovery, not a caller finish request.

The kit V2-3 adapter must reject unsupported mappings and bind its idempotency
handling to the derived result. Phase 25 qualification must exercise matching,
mismatching, missing-proof, abandonment, and replay cases (CF-037). Task 8 owns
fid derivation; task 9 removes caller finish; task 18 proves the fid surface
and records this boundary without claiming a production adapter exists.

## Complete Comparison Register

| Review item | Disposition / owner |
|---|---|
| C-01 | AGREE: terminated/live/unknown and lease-expiry distinction retained. |
| C-02 | RECORD: exact Human force grant remains pinned to Envelope revision 1 in this slice. Phase 24 owns multi-revision generalization. |
| C-03 | RECORD: Review `derived`/`observed` charging and `review_cycles` are excluded here; Phase 23 plans their exact budget vocabulary and evidence. |
| C-04 | RECORD: only `budget-exhausted` is accepted now. Phase 23 owns `review-rework-exhausted`; Phase 24 owns generic `retry-exhausted`. Kit `circuit-breaker` is a coordinator classification, not a fid Attention field. No mapping is inferred from message text. |
| C-05 | AGREE: immutable history and IDs retained across recovery/replay/import. |
| C-06a | RECORD: kit recoverActivation maps to resource observations, recovery selection, then recovery completion under the same fid ID. V2-3 adapter owns this multi-call mapping. |
| C-06b | DECIDE: coordinator-owned export binding, qualification obligation above; no fid archive reference claimed. |
| C-06c | RECORD: activation-version-conflict and activation-direction-conflict both map to kit activation-conflict with the original fid reason preserved; competing prepare already uses that code. V2-3 adapter owns the mapping. |
| C-07 | RECORD: kit validate-only maps to non-mutating `export validate`; import remains empty-target. Project-wide current export only; scope/asOf selection requires a later accepted port decision before qualification, not an invented Phase 22 flag. |
| C-08 | DECIDE: terminalState is an expectation against derived truth, with separate intentional abandonment as specified above. |
| C-09 | RECORD: fid derives semantic Attention identity. Adapter computes the same canonical identity and rejects a conflicting supplied deduplicationKey; it cannot override fid's key. Kit signature/invariant inconsistency is a kit contract follow-up under CF-037. |
| C-10 | RECORD: fid resolves Attention from all linked Stop disposition Events; Human Outcome references remain unavailable until Phase 23. `human_outcome_resolution:false` stays explicit. |
| C-11 | RECORD: harness reconciliation is an assertion; fid re-observes and verifies it using the canonical digest. No supplied reconciliation is trusted as an observed fact. |
| C-12 | RECORD: recovery count derives from immutable run_recoveries; only force charges interventions. No stored resumptions budget unit. Phase 25 owns trajectory telemetry. |
| C-13 | RECORD: use the mapping below. Missing kit code mappings are qualification work, not a generic fallback. |

| fid Stop code | Kit mapping / qualification condition |
|---|---|
| executor-liveness-unknown | Same code |
| recovery-state-conflict | Same code |
| budget-exhausted | Same code; coordinator may classify Attention as circuit-breaker |
| activation-incomplete | Same code |
| capability-unavailable | unsupported-capability; owner_phase is opaque metadata |
| budget-data-unavailable | No accepted kit equivalent. Require an explicit kit code/representation decision before V2-3 qualification; never silently map data absence to unsupported-capability. |
| operator-stop | No accepted kit equivalent. Require an explicit kit code/representation decision before V2-3 qualification; preserve the authorized intervention fact. |

## Final QA Obligation

VAL-080 remains QA Lead-owned. Recheck source currency and the as-built corpus
against every row here and the independent review. An accepted scoped exclusion
is not evidence that a kit adapter obligation has been implemented. Any new
unresolved difference returns to Tech Lead. CF-037 and the Phase 23/24 roadmap
notes keep these receiving obligations visible after Phase 22 wrap-up.
