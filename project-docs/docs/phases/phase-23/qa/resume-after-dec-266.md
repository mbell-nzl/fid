# Phase Runner handoff after DEC-266 (ESC-23-010, second disposition)

DEC-266 resolves the `claim.acquire` authority-details stop, but task 18 is
**not** complete. Keep all of the following:

- TRG-23-011, including its second disposition, and ESC-23-010;
- the 2026-09-30 18:58 STOP;
- task-18 attempts 2-6, their checkpoints, and the report through `117625d`
  (implementation checkpoint `72113dd`);
- every failed focused and tagged run.

Do not allocate a replacement TRG or ESC.

## What DEC-266 decided

1. **The contract governs.** Phase 21 `public-interface.md` section 7 is
   retained by Phase 23 section 1. All four `authority-*` codes return
   exactly `{operation, actor_id (nullable), required_actions,
   required_capabilities}` with a null subject (Phase 23 section 2 /
   DEC-233).
2. **Cause.** The registry (`internal/cli/codes.go:286-288`,
   phase-20/task-1) was never converged for `authority-missing`,
   `authority-unknown-actor`, and `authority-denied`. As a result:
   - eight producers emit the Phase 20 `{operation, actor}` shape;
   - `gate.go:325/327` emit the contract shape, which the registry rejects,
     so those paths can never return a typed failure;
   - the fake's `MissingCapability` and `human-outcome.record`
     `authority-missing` rows expect the non-conforming shape.
3. **Unchanged:**
   - contract text, codes, categories, exits, and subjects;
   - VAL text, validators, and owners;
   - all five task-18 commands;
   - DEC-265.
4. **Scope.** Task 18 gains a bounded product convergence under
   `tasks/task-18.md` "Planning Amendment - DEC-266".

## Resume binding

**Binding correction (2026-09-30).** The first issue of this handoff
prescribed attempt 7, taken from TL-Lite's narrative count. The sequence
validator binds a post-STOP RUN-RESUME to its governing ledger TRIAGE plus
one. The 18:58 STOP governs through the 15:02 TRIAGE TRG-23-011, whose
`Attempt count` is 1, so the ledger attempt is **2**. The Runner correctly
refused attempt 7: `INVALID RUN-RESUME Attempt count must increment its
governing TRIAGE` (terminated `exit_code: 1`), and the ledger stayed
`VALID state=STOPPED`. Ledger attempts 2-6 before the 18:58 STOP remain as
recorded. The ledger restarts its count at each post-STOP resume, and the
event timestamps disambiguate the repeats. In the task-18 report, call this
continuation "report attempt 7 (ledger attempt 2 after the 18:58 STOP)".

- Governing TRG: TRG-23-011
- Governing ESC: ESC-23-010
- Trigger signature: task-18:protected-contract:failure-code-cardinality-and-semantic-fake-equality
- Attempt count: 2
- Resolution authority: DEC-266
- Continuing task: task-18
- Next recovery: task-18 | in-scope-implementation | task-scope | task-18:protected-contract:failure-code-cardinality-and-semantic-fake-equality | 2

Include `Sequence contract: v2` and `Recovery directive schema: v2`. In the
prior-STOP reference, name the 18:58 STOP, TRG-23-011, and ESC-23-010.
Record the checkpoint you actually verified and your own resolution
evidence:

- the commit that records DEC-266, and its subject;
- DEC-266 present in phase decisions;
- the task-18 Planning Amendment present;
- the ESC-23-010 second disposition marked RESOLVED;
- this handoff present.

The next event is a TASK-RECOVERY for task-18 that repeats the resume tuple
exactly: the same trigger, attempt 2, and a `Checkpoint` byte-identical to
the RUN-RESUME `Checkpoint`. The validator rejects a TASK-RECOVERY that
changes either. Its fields are:

- Category: `in-scope-implementation`
- Repair mode: `task-scope`
- Source: `DEC-266 and TRG-23-011; project-docs/docs/phases/phase-23/qa/resume-after-dec-266.md`
- Repair authority: `DEC-266 and TRG-23-011`
- Consecutive no-progress attempts: `0/2`. Attempts 5 and 6 made material
  progress, and both were forbidden from editing the producer that caused
  the failure.
- Every control-conformance field: `N/A`
- Allowed files: `ordinary task scope`. The task-18 spec, including the
  DEC-266 Planning Amendment, names the files.

## Bounded task 18 continuation (same-task Developer)

Start from HEAD, and read `task-18-done.md` and the DEC-266 amendment first.

1. **Inventory.** Rerun the amendment's `grep`, and record every non-test
   producer. Expect ten lines: the eight producers named in DEC-266 plus
   `gate.go:325/327`.
2. **Red first.** Correct the fake's `MissingCapability` and
   `human-outcome.record` `authority-missing` expectations to the four
   contract fields. Add the two gate regression tests. Run the focused
   command below, and record the red.
3. **Registry.** In `codes.go`, give the three codes the
   `authority-unavailable` four-field definition, and converge
   `codes_test.go`'s `authorityDetails()`.
4. **Producers.** Fill the four fields at the eight producer lines from the
   existing errors:
   - `EnvelopeAuthorityDeniedError`: `ActorID`, [`Action`], [`Capability`];
   - `review.AuthorityMissingError`: `ActorID`, [`Operation`],
     [`Operation`].
5. **Retained tests.** Any retained test that fails only on the shape
   converges with per-row provenance. Any other failure is
   `RECOVERY_REQUIRED`.
6. **Forbidden:**
   - any other product change;
   - adding a `SubjectType`;
   - any contract, code, VAL, or command change;
   - fake expectations derived from fid output;
   - rewriting the evidence from attempts 1-6.
7. **Commands.** Run the focused command first:

   ```text
   go test -tags integration -run '^(TestP23FailureCorpusHasEveryAcceptedCode|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestGrantActionSetIsExactlyThirtyNames)$' -count=1 ./cmd/fid/integration/
   ```

   Then run the new gate regression tests by name. Then run all five task-18
   `## Exact Commands To Run`, unchanged and in order, with observed
   termination. Record the tagged, `cmd/fid`, and uncached
   `internal/workflow` wall times against VAL-079.
8. **Report.** Update `task-18-done.md`:
   - keep the evidence from attempts 1-6;
   - disclose DEC-266 under Deviations;
   - record each changed producer as a finding against its introducing
     task, with commit and `Refs:`;
   - update VAL-068 and the per-operation authority column.

   Commit the implementation and the report.

## Runner closure

1. Run report-lint and sequence-lint yourself, and verify the commit hash
   and subject.
2. Record REPAIR task-18 with `Outcome: RESOLVED`, then TASK task-18
   COMPLETE.
3. Dispatch QA window 16-18 (`review-13.md`) before phase-end QA. Besides
   the three checks in `qa/resume-after-esc-23-010.md`, it must confirm
   three things:
   - that the authority-producer inventory is exhaustive, and that all four
     `authority-*` codes carry the Phase 21 four-field details;
   - that the gate regression tests exercise `gate.go:325/327`;
   - that every authority fake row matches the contract and none was
     derived from fid.

## Simulation check

The corrected binding was checked on a scratch copy of the `phase-23` docs at
`e439d37` plus this correction, through
`workflow-sequence-lint.sh --append-stdin`:

- The unmodified ledger returned `VALID state=STOPPED`.
- A RUN-RESUME at attempt 2 with the tuple above returned
  `APPEND_CONFIRMED state=REPAIR_PENDING`.
- A TASK-RECOVERY at attempt 7 was rejected:
  `INVALID TASK-RECOVERY does not match the exact schema v2 RUN-RESUME
  recovery tuple`. So was a changed checkpoint:
  `INVALID TASK-RECOVERY does not preserve its RUN-RESUME checkpoint`.
- A TASK-RECOVERY at attempt 2 with the identical checkpoint returned
  `APPEND_CONFIRMED state=REPAIR_PENDING`, and the ledger then returned
  `VALID state=REPAIR_PENDING`.

Each passing run terminated `exit_code: 0`, and each rejection
`exit_code: 1`. This proves only the event binding; use real checkpoints
and receipts.

This handoff grants authority and scope. It is not a verification receipt.
The Tech Lead ran no implementation gate; the registry probe cited in
DEC-266 was a scratch overlay only.
