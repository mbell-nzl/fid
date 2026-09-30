# Phase Runner handoff after ESC-23-010

DEC-265 resolves ESC-23-010, but task 18 is **not** complete. Keep all of
the following:

- TRG-23-011 and ESC-23-010;
- the 2026-09-30 15:03 STOP;
- the BLOCKED report `333514946` and its corrections `dac84a1` and `a156691`;
- every red, intermediate, and exact-gate result that report records.

Do not allocate a replacement TRG or ESC.

## What DEC-265 decided

1. **No member conflict.** The section 5 table, the section 5.1 producer
   table, `p23FailureContracts`, and `failures.json` all name the same 19
   codes.
2. **Count.** Section 5 **declares 19** codes. **18 are added** by Phase 23.
   `requirement-not-found` is **retained**: it has been in `codes.go` since
   `ab4ac8f` (phase-20/task-1) with its contract unchanged, and it is
   re-declared because Phase 23 pins its new producers, `requirement.get`
   and `requirement.link`.
   - DEC-230's "18" and the section 5 prose counted the added codes.
   - Section 5's lead sentence called every row "added". That wording is
     now corrected.
3. **Unchanged:**
   - every code, category, exit, details shape, producer, and precedence;
   - `decision-not-found` stays removed;
   - VAL text, validators, and owners;
   - task ownership and all five task-18 commands.

## Resume binding

- Governing TRG: TRG-23-011
- Governing ESC: ESC-23-010
- Trigger signature: task-18:protected-contract:failure-code-cardinality-and-semantic-fake-equality
- Attempt count: 2
- Resolution authority: DEC-265
- Continuing task: task-18
- Next recovery: task-18 | in-scope-implementation | task-scope | task-18:protected-contract:failure-code-cardinality-and-semantic-fake-equality | 2

Include `Sequence contract: v2` and `Recovery directive schema: v2`. In the
prior-STOP reference, name TRG-23-011 and ESC-23-010. Record the checkpoint
you actually verified and your own resolution evidence:

- the commit that records DEC-265, and its subject;
- DEC-265 present in phase decisions;
- ESC-23-010 marked RESOLVED;
- this handoff present.

The next event is a TASK-RECOVERY for task-18 with the same trigger and
attempt 2. Its fields are:

- Category: `in-scope-implementation`
- Repair mode: `task-scope`
- Source: `DEC-265 and TRG-23-011; project-docs/docs/phases/phase-23/qa/resume-after-esc-23-010.md`
- Repair authority: `DEC-265 and TRG-23-011`
- Consecutive no-progress attempts: `0/2`. Attempt 1 does not count. It made
  material progress, and it could not certify the code set while the
  contract text read as contradictory.
- Every control-conformance field: `N/A`
- Allowed files: `ordinary task scope`. `task-18.md` Scope names the files:
  `cmd/fid/integration/`, its `testdata/`, audit-only
  `internal/workflow/authority.go`, and the report.

## Bounded task 18 continuation (same-task Developer)

Start from HEAD, and read `task-18-done.md` and its committed test files
first. Complete the task, not only the count:

1. **Failure-code corpus (DEC-265).**
   - Keep all 19 rows.
   - Add an independently pinned literal partition to `failures.json`: a
     per-row `added` or `retained` marker.
   - Make `TestP23FailureCorpusHasEveryAcceptedCode` assert 19 declared =
     18 added + 1 retained, with `requirement-not-found` as the only retained
     member.
   - Do not derive the partition from `codes.go`, a git query, or a parse of
     `public-interface.md`.
2. **VAL-068 semantic parity.** This is the open work.
   - Extend `TestP23SharedLeafConformance` and
     `TestP23SharedAmendedLeafConformance` so that, for all 23 operations in
     `operations.json`, the fake's pinned expectation is compared with the
     **complete** fid result value, not only its field set and replay flag.
   - **Exact equality:** every caller-supplied or deterministically derived
     value, such as verdicts, states, dispositions, relations, evidence
     types, validity, reason lists and their order, counts, and list
     membership and order.
   - **Structural comparison, only for values fid generates** (IDs,
     timestamps, EventRefs, fid-computed digests): assert well-formedness and
     cross-case consistency. For example, the ID a mutation returns equals
     the ID the matching read returns, and the Finding IDs in a Review equal
     those that `finding.list` returns.
     - Declare the normalization rules once, and apply them identically to
       both adapters.
     - Never copy a fid output value into a fake expectation.
   - Author fake expectations from `public-interface.md` sections 3 and 4
     and `contracts.md`.
3. **VAL-068 rejection classes.** Both adapters run the same case list, and
   each rejection compares code, category, exit, subject, full details under
   the same normalization, and an unchanged Event count. For every operation,
   cover each applicable class:
   - unknown field;
   - duplicate delivery: the same key and payload replays the identical
     result, and the same key with a different payload is rejected as the
     retained contract specifies;
   - authority, where section 4.1 or 4 gives a precondition;
   - stale version or binding, where the leaf takes one;
   - not-found, for every leaf that takes an ID.

   Record each class that does not apply to an operation as `N/A` in the
   report, citing the section that makes it inapplicable. A case that
   exists for only one adapter fails.
4. **Missing cases.** If a case belongs to an earlier owning task, you may
   author it here. Record it in the report as a finding against that owner,
   with provenance (binding constraint 3).
5. **VAL-067 provenance.** For each of the seven Phase 23 grant actions,
   record the introducing commit and its `Refs:` trailer in the report,
   showing that the action arrived at or after its owning task.
6. **Divergence.** If a fake/fid divergence exposes a product defect, do
   **not** edit product code, and do not adjust the fake to match fid.
   Return `RECOVERY_REQUIRED` or `BLOCKED` with the divergence, the owning
   task, and its provenance.
7. **Forbidden:**
   - any contract, code, count, or VAL change;
   - product behavior changes;
   - derived fake expectations;
   - rewriting attempt-1 evidence;
   - changing any command.
8. Run the focused command first:

   ```text
   go test -tags integration -run '^(TestP23FailureCorpusHasEveryAcceptedCode|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestGrantActionSetIsExactlyThirtyNames)$' -count=1 ./cmd/fid/integration/
   ```

   Then run all five commands in task 18's `## Exact Commands To Run`,
   unchanged and in order, with observed termination. Record the tagged,
   `cmd/fid`, and `internal/workflow` wall times against VAL-079, using an
   uncached `internal/workflow` time if the gate output is cached.
9. Update `task-18-done.md`:
   - keep attempt 1's evidence and the report-lint failure;
   - disclose DEC-265 under Deviations;
   - replace the per-operation table with the accepted-value and
     rejection-class comparison, one row per operation;
   - state the 19/18/1 partition;
   - update the VAL-067 and VAL-068 coverage lines.

   Commit the tests and the report.

## Runner closure

1. Run report-lint and sequence-lint yourself, and verify the commit hash
   and subject.
2. Record REPAIR task-18 with `Outcome: RESOLVED`, then TASK task-18
   COMPLETE.
3. Dispatch QA window 16-18 (`review-13.md`) before any phase-end work. In
   addition to its ordinary scope, it must confirm three things:
   - that the fake expectations are independent of fid output, and that the
     normalization rules touch only fid-generated values;
   - the 19/18/1 partition against DEC-265;
   - the `N/A` rejection-class citations and the VAL-067 provenance rows.

This handoff grants authority and scope. It is not a verification receipt.
The Tech Lead ran no implementation gate.
