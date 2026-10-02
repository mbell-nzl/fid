# Phase Runner handoff after ESC-23-006

DEC-260 resolves both halves of ESC-23-006. Task 12 is **not** complete.
Keep TRG-23-006, ESC-23-006, the 19:01 STOP, the RECOVERY_REQUIRED report,
both failed journey attempts and the failed report-lint result. Do not
allocate a replacement TRG or ESC.

## What DEC-260 decided

1. **Requirement write authority (Option B).** `requirement.record-reference`
   and `requirement.link` need a trusted Actor from the finalized Kit
   Activation, but no Actor-profile capability and no Envelope grant.
   - Only one paragraph changed: the authority paragraph in
     `public-interface.md` section 4.
   - The closed capability set stays at thirty names in both copies.
   - Section 4.1, R-19, execution-scope row 7, VAL-067 and task 18 are
     unchanged.
2. **Report-lint grammar.** The grammar is intended, so it is unchanged.
   Task 12's expected-nonzero declaration had been wrapped over two lines,
   which was an authoring defect. It is now one line, with the same command
   and `exit_code: 1`. Tech Lead checked the fix: report-lint on the retained
   report and the corrected spec returned `OUTCOME RECOVERY_REQUIRED` and
   `PASS report-lint completed with no failures` (terminated `exit_code: 0`).

## Resume binding

- Governing TRG: TRG-23-006
- Governing ESC: ESC-23-006
- Trigger signature: task-12:protected-contract:requirement-capability-registry-and-report-lint-grammar
- Attempt count: 2
- Resolution authority: DEC-260
- Continuing task: task-12
- Next recovery: task-12 | in-scope-implementation | task-scope | task-12:protected-contract:requirement-capability-registry-and-report-lint-grammar | 2

Include `Sequence contract: v2` and `Recovery directive schema: v2`. In the
prior-STOP reference, name TRG-23-006 and ESC-23-006. Record the checkpoint
you actually verified and your own resolution evidence: the commit that
records DEC-260 and its subject, DEC-260 in phase decisions, and ESC-23-006
marked RESOLVED.

The next event is a TASK-RECOVERY for the same task, trigger and attempt,
with category `in-scope-implementation` and mode `task-scope`. `Source` is
`DEC-260 and TRG-23-006; project-docs/docs/phases/phase-23/qa/resume-after-esc-23-006.md`
and `Repair authority` is `DEC-260 and TRG-23-006`. Set no-progress to `0/2`
and every control-conformance field to `N/A`. `Allowed files` is `ordinary
task scope`.

## Bounded task 12 continuation (same-task Developer)

1. In `internal/review/requirement.go`, remove the two
   `contains(actor.Capabilities, "requirement.…")` admission checks. Keep:
   - trusted-Actor resolution;
   - the `Source` digest check under the write lock;
   - idempotency;
   - every retained failure code.
2. In `cmd/fid/integration/p23_requirement_reference_test.go`, remove both
   requirement names from the fixture's Actor-profile action list. The
   journey's successful record, supersede and link then prove that admission
   needs no capability.
   - Complete the journey's VAL-051 digest-mismatch rejection and
     one-active-row observations.
   - Complete VAL-052's three relations and its unknown-requirement and
     unknown-relation rejections.
3. Converge any task-12-introduced test that asserts a requirement capability
   denial (introduced in `0e4bb33`). Record its provenance.
4. Do not touch either allowlist copy (`internal/workflow/activation.go`,
   `internal/journal/current_event_value_relations.go`), `report-lint.sh`, or
   the task spec.
5. Run the focused journey first:

   ```text
   go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration
   ```

   Then run all seven commands in task 12's `## Exact Commands To Run`,
   unchanged. Record `cmd/fid`, `internal/workflow` and tagged wall times
   against VAL-079.
6. Update `task-12-done.md`:
   - keep the two earlier journey FAILs and the earlier BLOCKED gate rows;
   - disclose DEC-260 under Deviations;
   - record the `rg` result in the exact declared form.

   Commit the implementation and the report.

## Runner closure

1. Run report-lint and sequence-lint yourself, and verify the implementation
   hash and subject.
2. Record REPAIR task-12 with `Outcome: RESOLVED`, then TASK task-12 COMPLETE.
3. Task 12 closes QA Review window 10-12. Dispatch that review before
   task 13, as the QA plan requires.

## Simulation check

The binding was checked on a scratch copy of `phase-23` with this sequence
appended to the retained ledger:

- RUN-RESUME;
- TASK-RECOVERY;
- REPAIR RESOLVED;
- TASK COMPLETE.

Each event was labelled STRUCTURE SIMULATION ONLY.
`workflow-sequence-lint.sh` returned `VALID state=REPAIR_PENDING` after the
TASK-RECOVERY and `VALID state=ADVANCEABLE` after TASK COMPLETE (each
terminated `exit_code: 0`). The canonical ledger separately returned `VALID
state=STOPPED` (terminated `exit_code: 0`). This proves only the event
binding. Replace the simulation values with real checkpoints and receipts.
Tech Lead ran no implementation gate; this handoff grants authority and
scope. It is not a verification receipt.
