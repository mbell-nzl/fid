# Phase 23 - Tech Lead Brief For The QA Lead Confirmation Pass

> Author: Tech Lead, 2026-09-19. Audience: a fresh QA Lead session.
> Subject: confirm the edits applied under DEC-229 and DEC-230 in response to
> the readiness refresh (`qa/qa-plan.md` section 12).
> Base: branch `phase-23` at `5ff523a`; the working tree was clean at commit.

## 1. What This Pass Is

Readiness condition 3 is already executed. Section 12 of the QA plan is a
fresh, independent refresh and **its measured results stand** except where an
edit below invalidates a stated fact. This pass is narrower: confirm that the
R-14 and R-15 dispositions are correctly applied, and independently judge the
two places where I did not do what the refresh asked and the one finding the
refresh did not catch.

Nothing in this pass requires re-running a gate. **No product file changed** -
all 19 edited files are under `project-docs/docs/`. Do not re-measure the
DEC-227 budgets, re-prove DEC-223 on the migration chain, or re-verify the
DEC-224 raising sequences; those results are fresh and unaffected.

Record the outcome as a new section 13 in `qa/qa-plan.md` with an explicit
`READY` or `NOT READY` verdict. If `READY`, the Phase Runner's readiness check
(condition 4) runs next and task 1 may then be dispatched.

## 2. What Changed, And Where

All of it is in commit `5ff523a`. Full reasoning is in
`decisions.md` under **DEC-229**, **DEC-230**, and the
"QA readiness-refresh disposition table".

| Revision | Disposition | Applied in |
|---|---|---|
| R-14 (a) `internal/workflow/stop.go` | Accept | tasks 6, 9, 10, 11; `execution-scope.md` 6 |
| R-14 (b) `current_event_value_catalog.go` | Accept | tasks 6, 9, 10, 11; `execution-scope.md` 6 |
| R-14 (c) `current_event_columns.go` | Accept the gap, **change the owner** | task 1 (all ten tables), task 15 (rename); `execution-scope.md` 6 |
| R-14 (d) `archive_semantics_catalog.go` | Accept, with a standing mirror rule | tasks 1, 9, 10, 15 (12 already had it); `execution-scope.md` 6 |
| R-14, `current_event_catalog.go` at tasks 4 and 5 | **Reject** | reason in DEC-229 and `execution-scope.md` 6 |
| R-15 CF-042 cause and mitigation | Accept, option (i) | task 1; `execution-scope.md` 7; `v1-convergence-inventory.md`; CF-042 in `decisions.md` and `decisions_active.md` |
| I-17..I-26 | All applied | see the disposition table |
| (new) `decision-not-found` | **DEC-230**, removed | `public-interface.md` 5; task 9 |

## 3. The Four Things To Judge, Not Just Check

These are where an independent read matters. The mechanical re-checks in
section 4 are secondary.

**3.1 The Condition-type sequencing inside R-14.** The refresh proposed adding
the two Stop registries to tasks 4, 5, 6, 8, 9, 10, 11, 12 "as each one's own
additions require", without saying which task registers which Condition type.
I fixed the split explicitly: a Condition type is registered by the task that
first **raises** a Stop carrying it, because `validateStopRaise` rejects an
unregistered type, so task 9 registers `decision-satisfied` and task 10
registers `human-outcome-recorded`. Task 11 registers only `review-recorded`
and `finding-resolved` and owns the resolution behavior of all four. I also
corrected task 11's objective and binding constraint, which read as though it
added all four. Confirm the split is complete and that no task raises a Stop
whose code or condition type is registered later than its own raising Event.

**3.2 Task 1 as the owner of `current_event_columns.go`.** The refresh proposed
five domain tasks. I assigned task 1 instead. My reasoning: task 1's binding
constraints already register all nine Event types with their complete
`schema.md` 6.2 effect rules; the ten effect tables those rules name belong
with the rules; the 0017 DDL that fixes the column domains ships in the same
task; and `TestCurrentEventColumnDomainsMatchAcceptedDDL` then proves the
registration against that DDL in the task that shipped it rather than five
tasks editing one shared map. Judge whether that executable proof genuinely
lands at task 1, and whether any effect on a new table can occur before task 1
registers it.

**3.3 The rejection at tasks 4 and 5.** The refresh said tasks 4 and 5 "own the
new `review.recorded` and `finding.addressed` rules and their
`currentEventCommands` entries" and so need `current_event_catalog.go`. Task
1's binding constraints already register all nine Event types, **their
operation mappings**, and their complete effect rules. Verify that directly
against `task-1.md`. If task 1's registration is in fact incomplete for those
two families, my rejection is wrong and this is a Must-fix; say so.

**3.4 DEC-230, the finding the refresh passed.** Section 12.6 checked producer
coverage at category level and recorded PASS. Pinning each producer under I-25
showed that `decision-not-found` has no producing leaf at all: section 4 gives
Decision References `record-reference` and `list` with no `get`, no mutation
takes a Decision Reference ID, and `linked_to_decision` is a transition rule
that reports satisfied or unsatisfied rather than a not-found error. VAL-056
requires one reproduction per added code through its named producer, so that
row was unsatisfiable for exactly that code. I removed the code rather than
adding a `decision.get` leaf, because four other Phase 23 entities already ship
`list` with no `get` and their not-found codes are reachable only through other
leaves. Verify the premise independently. If you find a real producer, the
removal is wrong. If you agree with the premise but prefer the leaf, say so -
it is reversible in one line and the Human has been told it is a live choice.

## 4. Mechanical Re-checks Worth Redoing

- **`public-interface.md` 5 and 5.1.** The grouped `*-not-found` row is
  replaced by six pinned rows. Confirm each pinned producer exists as a leaf or
  a named input field, that the table is well formed, and that section 5's code
  list and section 5.1's producer list are still the same set.
- **The new-code count.** Section 12.6 records 19 new codes. With
  `decision-not-found` removed it is **18** (19 declared in section 5, of which
  `requirement-not-found` already exists). Section 13 should restate that, as
  section 12.8 restated 6.5's "20".
- **Section 12.7's file-scope table.** Its R-08/R-09 rows are unaffected. The
  four rows this pass adds are new; re-derive the owning task for each of the
  four registries from the task specs, not from this brief.
- **The validation contract.** Still 79 rows. Two rows were edited and no row
  was added, removed, renumbered, or weakened: **VAL-006** gains Obligation
  counts 0, 1, 16, 17 in its evidence method and pass criterion (I-22);
  **VAL-034** makes the DEC-228 unresolvable-subject proof deliberate (I-23).
  Re-run your usual row count, one-validator-per-row, and citation checks.
- **Gate reach.** Every task that gained a file must have a required command
  that executes it. All four registries live in `internal/workflow`,
  `internal/journal`, and `internal/export`, covered by
  `go test -race -timeout 30m -coverprofile=coverage.out ./...`, which every
  task spec carries; `main_test.go` is covered by task 1's tagged gate. Confirm
  rather than assume.
- **Spec edits from the improvements.** `schema.md` 6.2 now names
  `retainedCurrentRules` for the three amended families (I-17);
  `contracts.md` 2 states that only the relative reason order is accepted
  (I-18); `public-interface.md` 5.1 step 1 excludes `obligation-unsatisfied`
  (I-19). Task 13 no longer renames `reasonCodes` in `cmd/fid/evidence_test.go`
  (I-20), task 4 names `idempotency-conflict` (I-21), task 10 records the blob
  remove/restore fixture technique (I-24), task 1 records the
  no-single-step-down fact and the dirty `schema_migrations` failure state
  (I-26).

## 5. R-15 Specifically

The recorded cause is corrected in all four places: `task-1.md`,
`execution-scope.md` 7, `v1-convergence-inventory.md`, and CF-042 in both
`decisions.md` and `decisions_active.md`. Your measured figures - 20
directories and 306 MB from one clean run, ~26 tagged runs against a 3.9 GB
`/tmp` - are quoted as the recorded cause.

I took option (i). Two points to check. First, the change stays inside
`cmd/fid/integration/main_test.go`: all five re-executing call sites already
spawn with `append(os.Environ(), ...)`, so a variable published by the parent
is inherited without editing the fault tests. Second, I rejected option (ii) -
pruning everything but the current process's directory would delete a
concurrent run's live binary, so age remains the only safe bound and the
24-hour prune survives only as the residual net for a parent killed mid-run.
Task 1 must now record the `fid-integration-binary-*` directory count and
`/tmp` free space either side of its tagged gate, and a clean run must leak
none. That overlaps QA-45 at phase end by design.

If you think the reuse fix belongs to a task other than 1, or that it is too
much for the schema baseline task, say so - it is an ownership question, not a
contract question.

## 6. Bounds On This Pass

- Confirm the applied edits and the reasoning behind them. Do not re-open
  R-01..R-13 or DEC-223..228; they are Human-approved and section 12 already
  re-proved them against fresh evidence.
- If you find a further Must-fix, return it to the Tech Lead. Do not weaken,
  narrow, or re-scope a `VAL-NNN` row to reach `READY`.
- Anything that needs a new contract amendment, a protected command change, a
  threshold or budget change, or new public surface is a blocker candidate and
  goes to the Human, not into a disposition.
- `decision.get` is the one live Human-facing choice; DEC-230 records both
  options and my recommendation.
