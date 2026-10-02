# Phase Runner handoff after DEC-270 (task 20, TRG-23-017)

DEC-270 resolves the task-20 blocker raised under ESC-23-012. ESC-23-012 itself
stays open: the DEC-269 repairs, the Human re-check, and `final-review-1.md`
are still pending. Task 20 is **not** complete. Keep all of the following:

- TRG-23-017 and the 2026-10-02 07:21 STOP;
- the BLOCKED report (`f739417`) and its 07:08 report-repair recovery;
- the **uncommitted** task-20 partial product and test tree, which is the
  continuation base. Do not stash, reset, or commit it on the Developer's
  behalf;
- the preserved pre-change journal under `/tmp/fid-task20-prechange.NjCRbu/`,
  if it still exists.

Do not allocate a new TRG or ESC.

## What DEC-270 decided

1. **The block was correct.** The pre-change replay criterion could not hold.
   A pre-task-20 `requirement.referenced` stored result lacks
   `requirement_id`, and exact result validation refuses it.
2. **The criterion was wrong.** `requirement.referenced` exists only in
   unreleased Phase 23 binaries: it is absent from `main` and `v0.1.0`.
   ADR-0011 makes no pre-adoption journal promise. Pre-change history is
   unsupported and must **fail closed**.
3. **Replacement criterion** (task 20 Planning Amendment; VAL-089):
   - (a) a committed new-shape round-trip test, run by the tagged gate;
   - (b) one-time evidence that the pre-change journal is refused with
     `journal-replay-stopped`, exit 6, and an unchanged digest.
4. **Unchanged:** DEC-269 Part 3's shape, exact validation, task-20 scope,
   all five commands, ownership, and the other criteria. No validator,
   legacy path, or versioned path may be added.
5. **Task 21:** if a pre-task-21 journal fails to rebuild because the rework
   Stop Condition changed, that is expected fail-closed behaviour under
   DEC-270, not a blocker.

## Resume binding

- Governing TRG: TRG-23-017
- Governing ESC: ESC-23-012
- Trigger signature: task-20:protected-contract:requirement-ref-journal-replay-compatibility
- Attempt count: 2
- Resolution authority: DEC-270
- Continuing task: task-20
- Next recovery: task-20 | in-scope-implementation | task-scope | task-20:protected-contract:requirement-ref-journal-replay-compatibility | 2

Include `Sequence contract: v2` and `Recovery directive schema: v2`. Name
`2026-10-02 07:21 STOP after TRG-23-017 (ESC-23-012)` as the prior STOP. Record
the checkpoint you actually verified and your own resolution evidence:

- the commit that records DEC-270, and its subject;
- DEC-270 present in phase decisions;
- the task-20 Planning Amendment and the amended VAL-089 present;
- this handoff present;
- the dirty paths, which must be only the task-20 partial paths plus the
  Runner-owned ledger and triage log.

The next event is a TASK-RECOVERY for task-20 with the same trigger and
attempt 2. Its fields are:

- Category: `in-scope-implementation`
- Repair mode: `task-scope`
- Source: `DEC-270 and TRG-23-017; project-docs/docs/phases/phase-23/qa/resume-after-dec-270.md`
- Repair authority: `DEC-270 and TRG-23-017`
- Consecutive no-progress attempts: `0/2`. Attempt 1 does not count: the
  criterion it had to meet was unsatisfiable by design, and it correctly
  stopped.
- Every control-conformance field: `N/A`
- Allowed files: `ordinary task scope`. The task-20 spec names the files.
- Originating report: `project-docs/docs/phases/phase-23/tasks/task-20-done.md`

## Bounded task 20 continuation (same-task Developer)

Read `git diff` first. It holds the value-catalog, projection,
`loadRequirementRef`, and test edits from attempt 1. Then:

1. Finish the requirement Human rendering in `cmd/fid/p23_human.go`, as
   task 20's Binding Constraints state.
2. Add the DEC-270 (a) round-trip test in an owned
   `cmd/fid/integration/p23_requirement_*_test.go` file. It records a
   reference, a superseding reference, and links. It captures JSON
   `requirement get` (stable ID and each reference ID) and `requirement
   list`, then runs `rebuild` and checks:
   - the source journal digest is unchanged;
   - the JSON reads are byte-identical;
   - an idempotent `record-reference` retry returns `requirement_id`.
3. Remove `TestP23Task20SeedPrechangeJournal`. No env-gated or skipped test
   is committed.
4. Converge any displaced retained assertion with per-row provenance.
5. Run all five task-20 commands unchanged, and record the tagged,
   `cmd/fid`, and `internal/workflow` wall times (VAL-079).
6. For DEC-270 (b), build the final binary and rerun:

   ```text
   go build -o /tmp/fid-task20-candidate ./cmd/fid
   /tmp/fid-task20-candidate --project /tmp/fid-task20-prechange.NjCRbu/judgment --json rebuild
   sha256sum /tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl
   ```

   Expect `journal-replay-stopped` at the `requirement.referenced` line,
   exit 6, and digest `a8fdbdccfc7fd39e27b80f256811d8b70113a0e30285122a166148455d1ba16e`.

   If `/tmp` was cleared, re-seed from a `d7f960d` worktree with that
   worktree's own binary. Use a temporary local seed test that is never
   committed, and record the new digest before and after.
7. **Forbidden:** relaxing result validation; any legacy, versioned, or
   tolerant path; touching `current_event_json.go` or any other unowned
   registry; rewriting the pre-change journal; changing any command. If any
   of these is needed, return `RECOVERY_REQUIRED`.
8. Update `task-20-done.md` to `COMPLETE`. Keep attempt 1's red and blocked
   evidence, and disclose DEC-270 under Deviations. Commit with
   `Refs: phase-23/task-20`.

## Runner closure

1. Run report-lint and sequence-lint yourself, and verify the commit hash and
   subject.
2. Record REPAIR task-20 with `Outcome: RESOLVED`, then TASK task-20 COMPLETE.
3. Continue with `qa/resume-after-dec-269.md` from task 21 onward: task 21,
   then `review-15` (window 19-21), then `final-review-1.md`. Review-15 must
   confirm DEC-270 (a) as a committed test and (b) as report evidence, and
   check that no env-gated test remains.

## Simulation check

The binding was checked on a scratch copy of `project-docs/docs`, with
DEC-270 recorded and the current dirty ledger (STOP at 07:21). A RUN-RESUME
and a TASK-RECOVERY, both labelled STRUCTURE SIMULATION ONLY, were appended
one at a time through `workflow-sequence-lint.sh --append-stdin`:

- the unmodified ledger returned `VALID state=STOPPED`;
- each append returned `APPEND_CONFIRMED state=REPAIR_PENDING`;
- the result returned `VALID state=REPAIR_PENDING`.

Each run terminated with `exit_code: 0`. This proves only the event binding.
Replace the simulation values with real checkpoints and receipts. The Tech Lead
ran no implementation gate. This handoff grants authority and scope; it is not
a verification receipt.
