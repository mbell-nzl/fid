# Phase Runner handoff after ESC-23-007

DEC-262 resolves ESC-23-007, but task 13 is **not** complete. Keep all of
the following:

- TRG-23-007 and ESC-23-007;
- the 20:46 STOP;
- the BLOCKED report and the failed activity trial;
- the zero-row journey probes;
- the failed tagged gate;
- the failed report-lint.

Do not allocate a replacement TRG or ESC.

## What DEC-262 decided

1. **Cause.**
   - `activity_log` was never registered in the strict-path
     `currentEventTables`. This is a Tech Lead planning gap, not a task-13
     defect.
   - Separately, the 0016 trigger resolves `requirement` and `decision`
     activity refs against the V1 tables until task 15's 0018 replaces it.
     This affects `decision-referenced` as well as the requirement names.
2. **Split at 0018.** No DDL, migration, order, effect rule, Event, or name
   changes.
   - Task 13 owns the registry entry (`execution-scope.md` section 6 row 11)
     and emits the seven `issue`-subject ActivityEvent names.
   - Task 15 emits the following with 0018's trigger, and records VAL-055's
     final ten-name inventory:
     - `requirement-referenced`;
     - `requirement-linked`;
     - `decision-referenced`;
     - `requirement`-subject Concern/Deviation rows.
3. **Subjects.** The `ref_type`/`ref_id` for each name is now in `schema.md`
   section 6.2.
4. **VAL-055's owner cell** names task 13 and task 15. Its pass and fail
   text is unchanged. Task specs 13 and 15 carry DEC-262 planning
   amendments.

## Resume binding

- Governing TRG: TRG-23-007
- Governing ESC: ESC-23-007
- Trigger signature: task-13:scope-gap:review-list-and-activity-producers
- Attempt count: 2
- Resolution authority: DEC-262
- Continuing task: task-13
- Next recovery: task-13 | in-scope-implementation | task-scope | task-13:scope-gap:review-list-and-activity-producers | 2

Include `Sequence contract: v2` and `Recovery directive schema: v2`. In the
prior-STOP reference, name TRG-23-007 and ESC-23-007. Record the checkpoint
you actually verified and your own resolution evidence:

- the commit that records DEC-262, and its subject;
- DEC-262 present in phase decisions;
- ESC-23-007 marked RESOLVED.

The next event is a TASK-RECOVERY for task-13 with the same trigger and
attempt 2. Its fields are:

- Category: `in-scope-implementation`
- Repair mode: `task-scope`
- Source: `DEC-262 and TRG-23-007; project-docs/docs/phases/phase-23/qa/resume-after-esc-23-007.md`
- Repair authority: `DEC-262 and TRG-23-007`
- Consecutive no-progress attempts: `0/2`. Attempt 1 made material progress
  (`review.list`). It was also not authorized to touch the registry that
  caused its remaining failure.
- Every control-conformance field: `N/A`
- Allowed files: `ordinary task scope`. The amended task-13 spec now names
  the journal registry and the Review producers.

## Bounded task 13 continuation (same-task Developer)

Follow `tasks/task-13.md` "Planning Amendment - DEC-262". In short:

1. Add `activity_log` to `currentEventTables` from `migrations/0002` DDL.
   Then admit the `activity` Effect in any further strict-path check it
   reaches:
   - `review.recorded` multi-row cardinality;
   - the `activity_log.details` JSON column.

   Loosen nothing, and edit no effect rule, stage, or `semanticEventRules`
   row. If a retained writer or a registry outside `internal/journal`
   rejects, return `RECOVERY_REQUIRED`.
2. In `internal/review/review.go`, `finding.go`, `concern.go`,
   `deviation.go`, and `human_outcome.go`, emit the `issue`-subject rows
   with the `schema.md` section 6.2 subjects:
   - `finding-resolved` is written once per `resolved` disposition and never
     for an `inapplicable` one.
   - Do not touch `decision.go` or `requirement.go`.
   - Write no `requirement` or `decision` activity row.
3. Complete VAL-056's producer matrix and both precedence probes, and
   VAL-057's replay corpus.
4. Run the focused journeys first:

   ```text
   go test -tags integration -run '^(TestP23ReviewRecordRealProcessJourney|TestPhase23LeavesRejectUnknownRequestFields)$' -count=1 ./cmd/fid/integration
   ```

   Then run all five commands in task 13's `## Exact Commands To Run`,
   unchanged. Record the `cmd/fid`, `internal/workflow`, and tagged wall
   times against VAL-079.
5. Update `task-13-done.md`:
   - keep every earlier failure;
   - disclose DEC-262 under Deviations;
   - report the four task-15 activity rows under VAL-055 as "deferred to
     task 15 (DEC-262)".

   Commit the implementation and the report.

## Runner closure

1. Run report-lint and sequence-lint yourself, and verify the
   implementation hash and subject.
2. Record REPAIR task-13 with `Outcome: RESOLVED`, then TASK task-13
   COMPLETE.
3. Continue to task 14. QA window 13-15 reviews after task 15, and must
   confirm that task 15 closed VAL-055's ten-name inventory.

## Simulation check

The binding was checked on a scratch copy of `phase-23`, with DEC-262
recorded and ESC-23-007 resolved. A RUN-RESUME and a TASK-RECOVERY, both
labelled STRUCTURE SIMULATION ONLY, were appended to the retained ledger.

- The unmodified ledger returned `VALID state=STOPPED`.
- The simulated ledger returned `VALID state=REPAIR_PENDING`.

Each run terminated `exit_code: 0`. This proves only the event binding;
replace the simulation values with real checkpoints and receipts. The Tech
Lead ran no implementation gate. This handoff grants authority and scope.
It is not a verification receipt.
