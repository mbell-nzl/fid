# Phase Runner handoff after ESC-23-008

DEC-263 resolves ESC-23-008, but task 13 is **not** complete. Keep all of
the following:

- TRG-23-008 and ESC-23-008;
- the 2026-09-28 21:57 STOP;
- the attempt-2 BLOCKED report (`348a538`) and its implementation commit
  (`d1c016d`);
- the failing `TestP23DecisionReferenceProtectedLifecycleRealProcess` runs;
- the failed tagged gate;
- every earlier DEC-262 and attempt-1 record.

Do not allocate a replacement TRG or ESC.

## What DEC-263 decided

1. **Cause.**
   - Section 5 of the accepted `public-interface.md` copied the
     `protected-decision-required` and `human-outcome-required` StopDetails
     verbatim as the failure details.
   - It never added the denial fields that DEC-226 part 4 (Human-approved),
     section 5.2, and VAL-057 require.
   - This is a drafting defect, not a task-13 defect.
2. **Shape.** Both codes' details become their StopDetails followed by
   `stop_ids`, `attention_ids`, `event` (always null), and `replayed`
   (always `false`). This is the same form as `review-rework-exhausted` and
   Phase 22 `budget-exhausted`.
3. **Unchanged:**
   - the `fid-cli.v1` envelope;
   - StopDetails and Stop projections;
   - categories, exits, subjects, producers, and precedence;
   - VAL-056's and VAL-057's pass and fail text;
   - all gates.
4. **Specs aligned:** `public-interface.md` sections 5 and 5.2,
   `06_DESIGN_PLAN_CLI.md` section 5, and the validation-contract change
   log (a note only). Task 13 carries the "Planning Amendment - DEC-263".

## Resume binding

- Governing TRG: TRG-23-008
- Governing ESC: ESC-23-008
- Trigger signature: task-13:protected-contract:protected-decision-denial-envelope-val-057
- Attempt count: 2
- Resolution authority: DEC-263
- Continuing task: task-13
- Next recovery: task-13 | in-scope-implementation | task-scope | task-13:protected-contract:protected-decision-denial-envelope-val-057 | 2

Include `Sequence contract: v2` and `Recovery directive schema: v2`. In the
prior-STOP reference, name TRG-23-008 and ESC-23-008. Record the checkpoint
you actually verified and your own resolution evidence:

- the commit that records DEC-263, and its subject;
- DEC-263 present in phase decisions;
- ESC-23-008 marked RESOLVED.

The next event is a TASK-RECOVERY for task-13 with the same trigger and
attempt 2. Its fields are:

- Category: `in-scope-implementation`
- Repair mode: `task-scope`
- Source: `DEC-263 and TRG-23-008; project-docs/docs/phases/phase-23/qa/resume-after-esc-23-008.md`
- Repair authority: `DEC-263 and TRG-23-008`
- Consecutive no-progress attempts: `0/2`. Attempt 1 made material
  progress: VAL-055's DEC-262 subset and VAL-056 passed. It was also not
  authorized to widen the closed schema that caused its only remaining
  failure.
- Every control-conformance field: `N/A`
- Allowed files: `ordinary task scope`. The amended task-13 spec names the
  adapters, `internal/cli/codes.go`, and the retained test rows.

## Bounded task 13 continuation (same-task Developer)

Follow `tasks/task-13.md` "Planning Amendment - DEC-263". In short:

1. Extend the `protected-decision-required` and `human-outcome-required`
   definitions in `internal/cli/codes.go` with the existing
   `stop_ids`/`attention_ids`/`event`/`replayed` field helpers.
2. In `protectedDecisionFailure` (`cmd/fid/decision.go`) and
   `humanOutcomeRequiredFailure` (`cmd/fid/human_outcome.go`), add the four
   fields from the `StopBlockedError` the adapters already receive:
   - `stop_ids` from `Stops`, sorted by ID;
   - `attention_ids` from `AttentionIDs`;
   - `event: null`;
   - `replayed: false`.

   Touch `cmd/fid/issue_mutation_response.go` or `cmd/fid/envelope.go` only
   if a caller must pass that error through unchanged.
3. Converge the retained exact-details rows from tasks 9, 10, 11, and 13
   (for example `assertP23ProtectedDecisionDenial`):
   - keep every existing check;
   - add the four fields, with the IDs equal to the durable Stop and its
     Attention;
   - record per-row provenance and the executing command.
4. Strengthen the VAL-057 protected-Decision repeat probe so it asserts
   `stop_ids` and `attention_ids` equal the durable blocking Stop and its
   Attention, as well as byte identity.
5. **Forbidden:**
   - editing StopDetails, Stop rows, the journal value catalog, or the
     `internal/workflow/stop.go` StopDetails keys;
   - changing any category, exit, subject, producer, precedence, or the
     envelope;
   - weakening any assertion;
   - changing any command.

   If any of these is needed, return `RECOVERY_REQUIRED`.
6. Run the focused command first:

   ```text
   go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23ReviewCyclesRealProcessJourney|TestP23FailurePrecedenceWithProtectedDecisionAndObligation|TestP23FailurePrecedenceWithReadinessStopAndObligation|TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess)$' -count=1 ./cmd/fid/integration
   ```

   Then run all five commands in task 13's `## Exact Commands To Run`,
   unchanged. Record the `cmd/fid`, `internal/workflow`, and tagged wall
   times against VAL-079.
7. Update `task-13-done.md`:
   - keep every earlier failure, including attempt 2's protected-denial
     red and its failed tagged gate;
   - disclose DEC-263 under Deviations;
   - keep reporting the four task-15 activity rows under VAL-055 as
     "deferred to task 15 (DEC-262)".

   Commit the implementation and the report.

## Runner closure

1. Run report-lint and sequence-lint yourself, and verify the
   implementation hash and subject.
2. Record REPAIR task-13 with `Outcome: RESOLVED`, then TASK task-13
   COMPLETE.
3. Continue to task 14. QA window 13-15 reviews after task 15. It must
   confirm three things:
   - DEC-263's details for both codes;
   - the converged rows' provenance;
   - that task 15 closed VAL-055's ten-name inventory.

## Simulation check

The binding was checked on a scratch copy of `phase-23`, with DEC-263
recorded. A RUN-RESUME and a TASK-RECOVERY, both labelled STRUCTURE
SIMULATION ONLY, were appended to the retained ledger.

- The unmodified ledger returned `VALID state=STOPPED`.
- The simulated ledger returned `VALID state=REPAIR_PENDING`.

Each run terminated `exit_code: 0`. This proves only the event binding;
replace the simulation values with real checkpoints and receipts. The Tech
Lead ran no implementation gate. This handoff grants authority and scope.
It is not a verification receipt.
