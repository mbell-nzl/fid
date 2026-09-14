# Phase 22 stored-result clarification — typed snapshots accepted

**2026-09-12 status:** The Human explicitly approved the recommended
typed-snapshot option: “I approve the recommended typed-snapshot option.”
DEC-216 records that authority. The proposal and its original pending-state
discussion below are retained as decision context and are superseded by this
acceptance. The specialization is now incorporated into schema.md and
public-interface.md. Final implementation verification remains pending;
acceptance is not a passing gate.

The inherited Phase 21 storage rule says `idempotency_results.result_json`
contains the full non-replayed Response Envelope (contracts.md §14, including
its normative example). Existing producers instead store closed operation
snapshots. Phase 22 adds a durable pending Attempt and committed budget denial,
without defining how those states fit the blanket envelope-storage rule.
An unsuccessful Gate must emit a FailureEnvelope with no outer result/Event;
a pending Attempt has no terminal response. Wrapping every stored snapshot in
`ok:true` would therefore be incorrect.

## Recommended option: specify the existing typed snapshots explicitly

For the finite 22 current Phase 22 Event families, amend schema.md §5 to override
the inherited full-Response-Envelope storage sentence. Store the complete,
closed, canonical operation-result snapshot plus `idempotency_key` and
`replayed:false`, with exact IDs/times bound to the containing Event and its
Effects. This is a persistence snapshot, never a claim that execution succeeded.

The Gate states are explicit:

- Pending `gate.execution-started`: `{attempt:Attempt,receipt:null,
  idempotency_key:ID,replayed:false}`. The recorded key fences another spawn.
- Terminal `gate.receipt-recorded`: `{receipt:Receipt,idempotency_key:ID,
  replayed:false}`, for every observed terminal outcome, including cancellation.
- Committed `run.execution-blocked`: `{denial:Denial,idempotency_key:ID,
  replayed:false}`, retaining its original Budget, Stop/Attention IDs and EventRef.
- Ordinary current mutations: their complete existing operation-result body,
  with the same required idempotency metadata and full fact bindings.

CLI adapters still produce the accepted SuccessEnvelope or FailureEnvelope.
Unsuccessful Gate outcomes remain gate-unsuccessful/exit8. Committed budget
denials remain ok:false/exit4. Replay preserves original facts and identities,
changes only the existing replay metadata, and never respawns or reevaluates
current budget state to reconstruct a historical result.

The amendment would explicitly reconcile Phase22 public-interface.md §1 replay
wording and §4 committed-denial storage wording with this snapshot boundary.
It would not change CLI output shapes, exit codes, process policy, Event header
or Effect columns, DDL, accepted terminal outcomes, or validation thresholds.
It would authorize no dual-format reader or rewriting existing journals.
Pinned Phase21 archive grammar remains separately scoped and unchanged.

Verification remains the unchanged Task1/full phase commands, with original
snapshot replay, successful/unsuccessful/pending/denied Gate branches, corrupted
nested value/result rejection, and full live/replay/rebuild/import parity.

## Alternative: implement full response storage plus a pending specialization

Retain full accepted success/failure Response Envelopes for ordinary/terminal
results and define one closed pending-fence snapshot exception. This requires
storage/read adapter changes for success, failure, denial and pending variants,
and corresponding fixture convergence. Public behavior remains unchanged.
The pending exception still needs an accepted specification amendment.

## Authority

Neither representation may be selected from passing tests or historical code.
Human acceptance is needed for this storage specialization under AGENTS.md's
product/data-model decision hard stop. The read-only Tech Lead audit found no
later accepted decision already resolving it. Until accepted, the snapshot
validator is a candidate repair; Phase22 is not technically complete.


## Independent contract audit

A bounded read-only Tech Lead audit confirmed both the ordinary envelope rule
and the missing Gate specialization. Its initial generic-envelope conclusion
was explicitly qualified after checking unsuccessful/pending/denial branches.
No source edits or gates were delegated. Relevant authoritative clauses:

- Phase21 contracts.md §14, normative example and full-envelope storage rule.
- Phase21 public-interface.md §6 terminal replay and §8 gate.run outcomes:
  unsuccessful delivery/replay use FailureEnvelope, not an ok:true wrapper.
- Phase22 public-interface.md §4: committed denial stays ok:false with the
  original EventRef and Stop/Attention IDs in error details.
- Phase22 contracts.md §7: a pending/abandoned Attempt key must not respawn.

The Human clarification request is pending. No normalization or storage-format
change has been made. The current candidate continues to validate typed bodies;
that is not an accepted replacement for the inherited storage rule yet.
