# Phase 23 - Tech Lead Brief For The QA Lead Confirmation Pass

> Author: Tech Lead, 2026-09-19. Audience: a fresh QA Lead session.
> Subject: confirm the edits applied under **DEC-231** in response to R-16
> (`qa/qa-plan.md` section 13.5).
> Base: branch `phase-23`; the working tree was clean at commit.
> This supersedes the previous brief, which covered DEC-229 / DEC-230. That
> pass is closed: section 13 confirmed every one of its edits.

## 1. What This Pass Is

The third confirmation cycle, and the narrowest. Sections 12 and 13 of the QA
plan both stand in full; **nothing in either is reopened**. Section 13's own
verdict was that DEC-229 and DEC-230 are sound and faithfully applied, and I
have changed nothing they touched. The only open item was R-16, and this pass
confirms its disposition.

**No product file changed.** Every edited file is under `project-docs/docs/`.
Do not re-measure the DEC-227 budgets, the DEC-223 migration chain, the DEC-224
raising sequences, the failure-code sets, or the contract row counts. Those
results are fresh and unaffected.

Record the outcome as a new section 14 in `qa/qa-plan.md` with an explicit
`READY` or `NOT READY` verdict. If `READY`, the Phase Runner's readiness check
(condition 4) runs next and task 1 may be dispatched.

## 2. The Disposition, In One Line

**R-16 accepted. Owner is the task that deletes the family's fence key, not
task 1.** Recorded as DEC-231 with a standing rule in `execution-scope.md`
section 6.

I verified R-16's mechanism independently in source before dispositioning it -
`validateCurrentAppliedResults` gating on `currentEventRules` membership, the
`fields` map returning `unregistered current result family`, the
`currentExpectedResult` default returning `unknown current result`, and
`currentProjectedValue` reading `currentValueShapes`. The finding is correct and
is the largest of the five registry gaps this phase has caught.

## 3. The One Thing To Judge, Not Just Check

R-16 offered two owners and asked me to pick. **I picked the one QA listed
second, and against the apparent precedent of DEC-229.** That choice is the
substance of this pass; everything in section 4 is mechanical.

My reasoning, which you should attack:

- The test is whether the owning task can *execute* what it registers. DEC-229
  put `current_event_columns.go` at task 1 because
  `TestCurrentEventColumnDomainsMatchAcceptedDDL` introspects the 0017 DDL that
  task 1 itself ships, needs no Event appended, and fails closed - your own
  section 13.1 confirmed exactly that.
- These three registries have no equivalent. They execute *only* on append of an
  Event of the family, and `unavailableJudgmentEvents` makes that impossible for
  all nine until the unfencing task. Putting all nine at task 1 would write
  roughly nine grammars, nine switch branches, fourteen value shapes, and three
  fill helpers behind a fence guaranteeing that none is ever executed in the
  task that wrote them - and would hand task 4 a silent inherited defect instead
  of its own red test.
- Derivability differs as well. Column domains follow from the DDL alone.
  `currentProjectedValue` maps a projected public row plus child-table queries
  onto a field list, which is the domain the unfencing task builds.
- I read this as the *consistent* application of DEC-229, not a departure:
  DEC-229 already ruled that a Condition type is registered by the task that
  first raises a Stop carrying it. Same principle - registration lands where it
  first executes.
- The shared-map objection does not bind. The eighteen tasks are strictly
  sequential, so there is no concurrent edit; and I searched the
  `internal/journal` and `internal/export` test files for any test iterating
  `currentEventRules` to cross-check the other three registries and found none,
  so a partial registration cannot be observed before the fence opens.

**Judge whether that reasoning holds.** Specifically: is there in fact any proof
that would land at task 1 and that I have missed? And is there any path by which
one of these three registries is reached before its family's fence key is
deleted - a retained family projecting a new value, an export or import path, a
rebuild, or a replay - that would make the unfencing task too late? If such a
path exists, my disposition is wrong and this is a Must-fix; say so.

The second half of the ownership test also needs your check: for each of tasks
3, 4, 5, 6, 8, 9, 10, 12, confirm the task's own required commands actually
execute the files it just gained. I checked that all eight carry the byte-
identical `go test -race -timeout 30m -coverprofile=coverage.out ./...` and that
`internal/journal` has in-package tests, but confirm it rather than take it from
me.

## 4. Mechanical Re-checks

| Claim | Where |
|---|---|
| Nine fence keys map one-to-one onto six unfencing tasks: 4 `review.recorded`; 5 `finding.addressed`; 8 `concern.recorded`/`concern.disposed`/`deviation.recorded`; 9 `decision.referenced`; 10 `human-outcome.recorded`; 12 `requirement.referenced`/`requirement.linked` | task specs 4, 5, 8, 9, 10, 12; `execution-scope.md` 6 |
| `current_event_json.go` and `current_event_result_projection.go` are named by exactly those six tasks and no others | task specs |
| Value shapes are registered once and reused: task 4 `Review`, `Coverage`, `CoverageItem`, `CoverageExclusion`, `Finding`, `FindingDisposition`; task 8 `Concern`, `ConcernDisposition`, `Deviation`; task 9 `DecisionRef`, `ProtectedDecision`; task 10 `HumanOutcome`; task 12 `RequirementRef`, `RequirementLink`. Tasks 5 and 10 explicitly register nothing new and are told to escalate rather than edit a shape they do not own | task specs 4, 5, 8, 9, 10, 12 |
| Composite fill helpers are assigned: `Review.findings[]`/`coverage` and `Finding.dispositions[]` at task 4, `Concern.disposition` at task 8 | tasks 4, 8 |
| Task 4 owns `FindingDisposition` although task 5 is the disposition task, because `review.recorded`'s own accepted rule inserts `finding_dispositions` rows for `resolves_finding_ids` | `schema.md` 6.2; task 4 |
| Sub-items: `Claim.mode` widened at task 3; `Budget`/`ImmutableBudget` `unit` and `exhaustion_action` widened at task 6; `Subject` widened to the full accepted SubjectType set at task 4 in one monotonic edit | tasks 3, 4, 6 |
| Task 1 is unchanged and gains nothing from R-16 | `task-1.md` |
| The DEC-229 rejection at tasks 4 and 5 still stands: they gain these three files, not `current_event_catalog.go` | `execution-scope.md` 6 |
| The contract is untouched at **79** rows; no `VAL-NNN` added, removed, renumbered, reworded, or weakened. Only the header provenance block changed (I-27) | `qa/validation-contract.md` |
| I-28 applied: the vacuous `internal/journal/` bullets in tasks 6 and 11 are replaced with explicit "no catalog edit; escalate instead" instructions | tasks 6, 11 |
| I-29 applied: the `envelope.activated` owner cell now reads "rows added by task 9; also raised from this family by task 10" | `schema.md` 6.2 |

## 5. Two Things I Recorded That You Did Not Raise

**5.1** `idempotency_results` is encoded as an *optional* stage in
`currentEventRules` for nearly every existing family, yet
`validateCurrentAppliedResults` requires the result row unconditionally once a
family is registered. The Phase 23 rows in `schema.md` 6.2 mirror that existing
encoding exactly, so this is pre-existing and not a new gap - but a Developer
reading "optional" could reasonably omit the typed snapshot that DEC-216 and
`schema.md` 6.3 make mandatory. I recorded the distinction in DEC-231 and in
`execution-scope.md` 6. Confirm it is in fact pre-existing and uniform, and that
no Phase 23 family is designed to omit the row; if one is, that is a real gap I
have mis-recorded as a note.

**5.2** The `Subject` widening was a sub-item in your R-16 conditioned on "if
any new shape is written as `subject:Subject`". Four are - `Finding`, `Concern`,
`Deviation`, `HumanOutcome` - so I made it unconditional and assigned the whole
widening to task 4 in one edit rather than letting six tasks each widen their
own slice. Judge that: it is monotonic and the full set is already accepted in
`public-interface.md` section 2, but it does mean the journal accepts a
`concern` subject type from task 4, four tasks before Concerns exist. My
reasoning is that the fence makes that unreachable and that six chances to
forget is the larger risk. Disagree if you see a reachable path.

## 6. What Would Make This NOT READY

A fourth Must-fix. Each of the last three passes found the next gap inside the
previous fix, so I am not treating "confirmed" as the expected answer. If a
fourth closed registry on the strict path has no owning task, say so plainly and
I will disposition it the same way; the confirmation gate stays until a pass
adds nothing. Do not weaken an assertion or narrow a `VAL-NNN` to reach Ready.
