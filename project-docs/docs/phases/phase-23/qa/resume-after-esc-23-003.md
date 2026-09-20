# Phase Runner handoff after ESC-23-003

DEC-239 resolves the protected contradiction; Task 4 is not complete. Retain
TRG-23-003, ESC-23-003, the original STOP, blocked report and failed attempts.
Verify the committed decision and actual checkpoint before appending a schema-v2
RUN-RESUME. Do not allocate a replacement TRG or ESC.

## Resume binding

- Governing TRG: TRG-23-003
- Governing ESC: ESC-23-003
- Trigger signature: task-4:in-scope-implementation:review-record-independence-tracer
- Attempt count: 9
- Resolution authority: DEC-239
- Continuing task: task-4
- Next recovery: task-4 | in-scope-implementation | task-scope | task-4:in-scope-implementation:review-record-independence-tracer | 9

Include `Sequence contract: v2` and `Recovery directive schema: v2`, actual
resolution verification evidence and the prior STOP reference. The next
TASK-RECOVERY binds exactly the same task, trigger, attempt, category and mode;
its `Source` and `Repair authority` are exactly `DEC-239`. Control-conformance
fields are `N/A`; this is ordinary task-scope implementation. Follow the ledger
template for the other mandatory fields and verify sequence state yourself.

## One bounded Task 4 closure

1. Preserve the blocked report's original evidence and disclose the DEC-239
   evidence-route amendment. The Run-holder-only fixture proves the exact
   conflicting Run and Actor using the real domain query over isolated stored
   history; it is not evidence that this state can pass `run.start`, journal
   validation, replay, or import.
2. Restore accepted repeat-Reviewer behavior. `contracts.md` section 3 and
   DEC-225 permit a later Review from the same independent Actor to resolve
   that Actor's earlier Finding. The present `creatorActorID == actorID`
   rejection in `internal/review/review.go` violates that contract. Use the
   existing task scope and red/green regression; preserve the per-Review
   prohibition and same Issue/Envelope revision and non-terminal checks.
   Evidence a real CLI Builder rework followed by the same Reviewer's second
   successful Review; different Reviewers do not prove this requirement.
   Carry the accepted EVID-12/13 real-process evidence through the built binary
   in `cmd/fid/integration/`; existing in-process `executeCommandForTest`
   coverage is insufficient for the already-required real-process assertion.
   No `review.record` journey currently appears in that tagged suite; use its
   existing Phase 23 Claim-mode fixture patterns and the task-owned gate.
3. Prove the Task 4 acceptance criterion for resolving more than one earlier
   Finding atomically. The current fixture names only `FND-prior`; retain the
   multi-Finding creation case and add the multiple-resolution observation.
4. Close the report's VAL-003 UNVERIFIED fence/parity row. Preserve its original
   subject and record the existing convergence or causal task-scope correction
   with provenance and executing command; do not mark it passed from memory.
5. Run every unchanged command in Task 4's exact command block in the recording
   session, observe termination, and record all results and package budgets.
   Preserve previous failure tokens and report the repairs as deviations.
   Commit the complete report only once all acceptance criteria have evidence.

The Runner then performs report-lint and sequence-lint, verifies the cited
implementation identity, records the linked RESOLVED recovery and Task 4
completion, and advances through Tasks 5-9 with reviews after Tasks 6 and 9.
Tech Lead has not run those controls or implementation gates; this handoff is
authority and scope, not a verification receipt.

## Task 9 readiness and stopping boundary

The current Human instruction supplies Task 9 dispatch clearance. Before Task
9, obtain the one bounded DEC-237 QA confirmation: validation rows, task
citations, row count, revision log and preservation of each evidence case.
R-24 through R-27 need no repeat review. Record any lossless granularity repair
under TL authority and confirm its final diff, without restarting the full
readiness sequence. Tasks 10 and 11 stay held. Stop after Task 9's window 7-9
review and its causal repairs; do not wrap Phase 23.
