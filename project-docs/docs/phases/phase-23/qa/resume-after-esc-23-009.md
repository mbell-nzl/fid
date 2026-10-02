# Phase Runner handoff after ESC-23-009

DEC-264 resolves ESC-23-009, but task 15 is **not** complete. Keep all of
the following:

- TRG-23-010 and ESC-23-009;
- the 2026-09-30 04:34 STOP;
- the BLOCKED report (`7884725`, corrected at `b1f0028`) and the checkpoint
  commit (`fa6355a`);
- the failing `TestMigration0018DownRestoresEmptyV1Tables` run and the failed
  touched-package run;
- the 04:21 interruption and its REPAIR event.

Do not allocate a replacement TRG or ESC.

## What DEC-264 decided

1. **Cause.**
   - `schema.md` section 5 restored objects as they were first created
     (0002/0003), not as they stood at schema 17.
   - This is a drafting defect with three instances:
     - V1 `concerns` kept `escalated_to_esc_id`, which 0016 down re-adds;
     - `issues_enriched` cited 0003, but the live text is 0011's;
     - the two 0004 kind triggers were missing.
2. **Rule.** 0018 down restores each object's fresh-schema-17 definition. The
   only accepted difference is SQLite's quoting of a renamed table's name. No
   earlier migration changes, and 0016 down keeps its `ADD COLUMN`.
3. **Exact SQL.** Corrected `schema.md` section 5 lists nine ordered steps.
   The Tech Lead verified them on a scratch database through the full down
   chain. That check is not a gate.
4. **Unchanged:**
   - migrations 0001-0017;
   - section 4 (0018 up);
   - VAL text and validators;
   - task ownership;
   - all nine task-15 commands.

## Resume binding

- Governing TRG: TRG-23-010
- Governing ESC: ESC-23-009
- Trigger signature: task-15:protected-contract:0018-down-concerns-column-duplication
- Attempt count: 2
- Resolution authority: DEC-264
- Continuing task: task-15
- Next recovery: task-15 | in-scope-implementation | task-scope | task-15:protected-contract:0018-down-concerns-column-duplication | 2

Include `Sequence contract: v2` and `Recovery directive schema: v2`. In the
prior-STOP reference, name TRG-23-010 and ESC-23-009. Record the checkpoint
you actually verified and your own resolution evidence:

- the commit that records DEC-264, and its subject;
- DEC-264 present in phase decisions;
- ESC-23-009 marked RESOLVED;
- this handoff present.

The next event is a TASK-RECOVERY for task-15 with the same trigger and
attempt 2. Its fields are:

- Category: `in-scope-implementation`
- Repair mode: `task-scope`
- Source: `DEC-264 and TRG-23-010; project-docs/docs/phases/phase-23/qa/resume-after-esc-23-009.md`
- Repair authority: `DEC-264 and TRG-23-010`
- Consecutive no-progress attempts: `0/2`. Attempt 1 does not count: the
  definition it needed was contradictory, so no attempt could have authored
  the down file. The 04:21 interruption consumed no budget.
- Every control-conformance field: `N/A`
- Allowed files: `ordinary task scope`. The task-15 spec, including both
  planning amendments, names the files.

## Bounded task 15 continuation (same-task Developer)

Start from the preserved partial tree at `fa6355a`, and read that diff
first. Then complete the whole task, not only the down file:

1. Write `migrations/0018_p23_remove_v1_judgment.down.sql` exactly as
   `schema.md` section 5 lists it. Copy the byte-identical passages from the
   cited migration files; do not retype them.
2. Extend `internal/db/migration0018_test.go` with the four section 5 proof
   obligations. The task-15 "Planning Amendment - DEC-264" lists them.
   Keep `TestMigration0018DownRestoresEmptyV1Tables`.
3. Complete the rest of task 15 as specified:
   - the V1 reader, maintenance, and CLI-leaf removal;
   - the retained-row convergence, with per-row provenance;
   - the DEC-262 `requirement`/`decision` activity producers and the
     ten-name inventory (VAL-055);
   - VAL-004 and VAL-061..063;
   - projection schema 18.
4. **Forbidden:**
   - editing migrations 0001-0017;
   - any 0018 down statement not listed in section 5;
   - weakening any assertion;
   - changing any command.

   If any of these is needed, return `RECOVERY_REQUIRED`.
5. Run the focused command first:

   ```text
   go test -run '^(TestLatestSchemaHasNoV1JudgmentTables|TestMigration0018Down.*)$' -count=1 ./internal/db
   ```

   Then run all nine commands in task 15's `## Exact Commands To Run`,
   unchanged, including the three expected-non-zero searches with their
   observed termination. Record the tagged, `cmd/fid`, and
   `internal/workflow` wall times against VAL-079.
6. Update `task-15-done.md`:
   - keep attempt 1's failures;
   - disclose DEC-264 under Deviations;
   - record per-row absence evidence for the whole inventory.

   Commit the implementation and the report.

## Runner closure

1. Run report-lint and sequence-lint yourself, and verify the
   implementation hash and subject.
2. Record REPAIR task-15 with `Outcome: RESOLVED`, then TASK task-15
   COMPLETE.
3. Dispatch QA window 13-15 before task 16. It must confirm four things:
   - DEC-263's details for both codes;
   - task 15's VAL-055 ten-name inventory;
   - the 0018 down fidelity proofs against `schema.md` section 5;
   - the converged rows' provenance.

## Simulation check

The binding was checked on a scratch copy of the `phase-23` docs with DEC-264
recorded (`532032b`). A RUN-RESUME and a TASK-RECOVERY, both labelled
STRUCTURE SIMULATION ONLY, were appended one at a time through
`workflow-sequence-lint.sh --append-stdin`.

- The unmodified ledger returned `VALID state=STOPPED`.
- Each append returned `APPEND_CONFIRMED state=REPAIR_PENDING`.
- The simulated ledger returned `VALID state=REPAIR_PENDING`.

Each run terminated `exit_code: 0`. This proves only the event binding;
replace the simulation values with real checkpoints and receipts. The Tech
Lead ran no implementation gate. This handoff grants authority and scope.
It is not a verification receipt.
