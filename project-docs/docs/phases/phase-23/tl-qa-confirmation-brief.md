# Phase 23 - Tech Lead Brief For The QA Lead Confirmation Pass

> Author: Tech Lead, 2026-09-19. Audience: a fresh QA Lead session.
> Subject: confirm the edits applied under **DEC-232** in response to R-17 and
> R-18 (`qa/qa-plan.md` section 14.5).
> Base: branch `phase-23`; the working tree was clean at commit.
> This supersedes the previous brief, which covered DEC-231. That pass is
> closed: section 14 confirmed DEC-231 on all four reachability paths and
> every mechanical check.

## 1. What This Pass Is

The fourth confirmation cycle. Sections 12, 13, and 14 of the QA plan all stand
in full; **nothing in any of them is reopened**. DEC-231 is confirmed correct
and unchanged in substance - the unfencing task is still the owner and the
executability test is still the test. DEC-232 generalises it and corrects two
of its sub-items.

**No product file changed.** Every edited file is under `project-docs/docs/`.
Do not re-measure the DEC-227 budgets, the DEC-223 migration chain, the DEC-224
raising sequences, the failure-code sets, or the contract row counts. Those
results are fresh and unaffected. The contract is still **79** rows.

Record the outcome as a new section 15 in `qa/qa-plan.md` with an explicit
`READY` or `NOT READY` verdict. If `READY`, the Phase Runner's readiness check
(condition 4) runs next and task 1 may be dispatched.

## 2. The Dispositions, In One Line Each

- **R-17 accepted; owner is task 4 alone.** Derived independently rather than
  adopted, as you asked, and confirmed.
- **R-18 accepted and strengthened.** `Subject` is not widened by any task; a
  separate `JudgmentSubject` is registered instead.
- **I-30, I-31 applied.** I-31 grew from four registries to six.
- **Three further defects of the same family** were found while dispositioning
  the pair, and are folded into DEC-232 as R-17b, R-17c, and R-18b.

## 3. The Three Things I Found, Which You Should Verify Hardest

These are mine, not yours, so they have had no independent review. Two of them
were found because R-17 asked me to derive the task set myself and I-31 asked
for an inventory. Treat them as this pass's primary subject.

**R-17b - `validateCurrentUpdate` (`internal/journal/current_event_updates.go`)
is a second unowned hard-failure registry.** Its `switch effect.Table` at `:14`
has `default: return "current Event update table not admitted"` and lists no
Phase 23 table. Two families emit an update:
`concern.disposed`'s **required** `judgment_concerns/update` and
`requirement.referenced`'s optional `requirement_refs/update`. The file is named
in zero phase documents. Assigned to tasks 8 and 12, renamed by task 15. My
reproduction: both return that error; `attention_items/update` is admitted.
Check that the changed-field and guard sets I directed each task to derive are
actually derivable from the accepted DDL, and that I have not missed a third
Phase 23 update Effect.

**R-17c - the `validateCurrentJSON` column switch
(`current_event_json.go:14`) is a second closed registry inside a file DEC-231
already assigned.** It returns `unregistered current JSON column <t>.<c>`. I
judged this *not* a new ownership gap, because the file is already owned by
every affected task and the failure is loud and in-task; I described the columns
instead of reassigning anything. Judge that call. If you think it deserves the
same treatment as R-17b, say so. I enumerated eleven columns across
`review_records`, `review_findings`, `judgment_concerns`, `deviations`,
`decision_refs`, `human_outcomes`, and `requirement_refs`, and concluded
`evidence_records.target_repository_json` needs nothing because `evidence.linked`
is retained. Verify the enumeration is complete against `schema.md` section 2.

**R-18b - `schema.md` 6.2's I-17 sentence misclassified `envelope.created`.**
It said the first of the four amended rows is in `currentEventRules`; all four
are in `retainedCurrentRules` (`current_event_catalog.go:82-83,87-88`). This
matters because a Developer following it would have moved a live family onto the
strict path. It also happens to be why the 16-entry `protected_decisions` input
needs no cardinality entry. Verify the classification directly in source - if I
have this backwards, several of my R-17 conclusions move with it.

## 4. The Judgement Call In R-18, Which You Should Test

QA's R-18 said the Subject widening is "not reachable as a defect today". I
agree, and I went further: I concluded no widening is possible at all, because
the two domains are not nested. `stops`/`attention_items` CHECK
`('run','envelope','activation')` (`0015:56,82`, unchanged by 0017); the four
judgment tables CHECK
`('issue','envelope','run','review','finding','requirement')`. Stop has
`activation` and judgment does not; judgment has `issue`, `review`, `finding`,
`requirement` and Stop does not.

I also concluded that DEC-231's cited source set was the wrong vocabulary
entirely: `public-interface.md` section 2's `SubjectType` row is the CLI **error**
subject vocabulary (`internal/cli/codes.go:88`), which is why `requirement` can
be "retained" there and why the six added values are exactly the six new entity
kinds that appear as error subjects. That reading is load-bearing - I renamed the
row `ErrorSubjectType` on the strength of it. **If that reading is wrong, say so
plainly; the rename is reversible in one line.** The alternative readings I
rejected are in DEC-232.

The remedy is `JudgmentSubject`, a new **name** for a domain already accepted
under DEC-222 in `contracts.md` section 8 and four identical DDL CHECKs. Confirm
that no accepted value is added, removed, or changed anywhere by it.

## 5. What Changed, Exactly

| File | Change |
|---|---|
| `execution-scope.md` 6 | The DEC-231 rule restated generally; a six-row registry inventory with file, line, rejection message, owners; two registries recorded as deliberately needing nothing; the Subject bullet replaced; I-30's 15-of-22 count |
| `public-interface.md` 1, 2, 3, 4 | Section 2 row renamed `ErrorSubjectType`; `JudgmentSubject` defined in section 3 with a four-vocabulary disambiguation table; the four judgment shapes and three leaf input rows reference it; section 1 pointer |
| `schema.md` 6.2 | The I-17 sentence corrected; why the retained stages are still required |
| `contracts.md` 8 | The already-accepted closed set linked to its new name |
| `tasks/task-4.md` | Cardinality scope and its bound; the JSON column switch; `JudgmentSubject` replaces the withdrawn Subject widening; two acceptance criteria citing VAL-017 and VAL-020 |
| `tasks/task-5.md` | Records the checked negative: this family needs no entry in the other three registries, and must not widen the disposition bound |
| `tasks/task-8.md` | `current_event_updates.go` added with its field and guard sets; JSON columns; the Subject reference corrected to `JudgmentSubject` |
| `tasks/task-9.md` | JSON columns; the retained-map correction for both amended families |
| `tasks/task-10.md` | JSON column; reuses `JudgmentSubject`; must not edit `Subject` |
| `tasks/task-12.md` | JSON column; `current_event_updates.go`; the end-of-phase audit widened from three registries to six |
| `tasks/task-15.md` | The `judgment_concerns` rename carried through five catalogs, not three |
| `decisions.md` | DEC-232 and its disposition table; status header |
| `../../current/phase_brief.md` | Status, active constraints, known risks; trimmed to 149 lines |

## 6. What To Check, In Order

1. **The three Tech Lead findings in section 3.** Independent source
   verification, not agreement with my prose.
2. **The `ErrorSubjectType` reading in section 4.** The one item I would most
   like a second opinion on.
3. **Inventory completeness.** `execution-scope.md` 6 now claims six registries
   on the strict path and two more that deliberately need nothing. Walk
   `validateCurrentEvent` and `validateCurrentAppliedResults` yourself and tell
   me if there is a seventh. This is the fifth pass to find a registry inside
   the previous fix, so the honest prior is that there may be.
4. **Per-family completeness.** For each of the nine families, does its owning
   task now name every registry entry it needs, and does task 5's "needs
   nothing" hold?
5. **Mechanical.** 79 contract rows unchanged; no VAL row renumbered or
   reworded; no product file touched; no protected definition, command,
   threshold, budget, or public value changed.

## 7. What Would Make This NOT READY

A seventh unowned registry; a Phase 23 family whose owning task still lacks an
entry it needs; the `ErrorSubjectType` reading being wrong; `JudgmentSubject`
changing an accepted value rather than naming one; or any of my three findings
failing source verification. Anything smaller is an improvement, and improvements
do not block the readiness check.
