# Phase 23 - Tech Lead Brief For The QA Lead Confirmation Pass

> Author: Tech Lead, 2026-09-20. Audience: a fresh QA Lead session.
> Subject: confirm the edits applied under **DEC-237** in response to R-28
> (`qa/qa-plan.md` section 20.5), P-01 (section 19.10), and I-42..I-45
> (section 19.9).
> Base: branch `phase-23`. Tasks 1-8 are dispatching in parallel, so the working
> tree carries their product changes; ignore them, and do not treat worktree
> cleanliness as this pass's authorization boundary.
> This supersedes the previous brief, which covered DEC-233. That pass and
> everything after it stand: QA-plan sections 12 to 20 are **not reopened**.

## 1. What This Pass Is

The ninth confirmation cycle, and a deliberately **bounded** one. It exists
because a Tech Lead edited the validation contract, and a contract edit should
not be its own confirmation. R-28's own requested action is exactly this pass.

**Scope is four things and nothing else:** the three affected rows
(VAL-078, VAL-080, VAL-081, VAL-082 - four IDs, three additions and one split),
their task-spec citations in tasks 9, 10, and 11, the row count and revision
log, and the derived pointers in `contract-coverage.md`.

**Do not re-measure** the DEC-227 budgets, the DEC-223 migration chain, the
DEC-224 raising sequences, the failure-code sets, the registry inventory, or the
R-24 / R-25 / R-26 / R-27 transcriptions. Both pass-8 sessions confirmed those
independently and nothing here changes them. **No product file changed**; every
edited file is under `project-docs/docs/`.

**Your verdict does not release dispatch.** Tasks 9, 10, and 11 remain held
until the Phase Runner receives explicit Human clearance, whatever you return.
Record the outcome as the next unused section of `qa/qa-plan.md`.

## 2. What Was Decided, In One Line Each

- **R-28 accepted.** VAL-078 asserted both that the four Condition types resolve
  correctly and that they are encoded correctly. Those fail independently, so it
  is two rows: VAL-078 keeps resolution, the new **VAL-080** carries the
  encoding matrix with every case, the DEC-236 scoping, and the DEC-235 shape
  evidence class intact. Second split of this kind in the phase after R-11.
- **P-01 accepted** as **VAL-081**: the Phase 22 nullable binding must still
  fire from task 9's new dispatch. Your argument is the reason - only detector,
  no other gate, one checkbox otherwise.
- **Extended to VAL-082** by me, on your own reasoning: DEC-235's six-row
  compatibility bar had no assertion either, and it is the larger invariant -
  `omitempty` on the retained fields changes every stored `semantic_key` and
  `request_digest` in existence. Every other Human-approved amendment in this
  phase carries at least one row; DEC-235 carried none.
- **Two routes declined a row, with reasons recorded in DEC-237:** QA-52 (a
  panic fails task 9's own gate, so it has a detector) and QA-49 (a missing
  details entry is caught late by VAL-064's rebuild leg; R-22's direct assertion
  is the right fix for locality, and the contract asserts end state).
- **I-42..I-45 applied.** I-44 is in task 9's table as "canonical means
  type-appropriate". I-45 promoted the relation confirmation into acceptance
  criteria in tasks 10 and 11, citing VAL-081.

## 3. The Four Asks

1. **Is the split lossless?** Compare VAL-078 before `a0e8e7e`'s successor with
   VAL-078 ∪ VAL-080 now. Every type, case, positive, negative, evidence class,
   validator, and owning-task obligation should be present, and nothing should
   be narrower. If any case was dropped or softened, that is a Must-fix.
2. **Is each row now one observable?** VAL-078 and VAL-080 are the point of the
   exercise. **VAL-082 is the one I am least sure of**: I wrote it as one
   property, byte compatibility, measured by three legs (per-type validation,
   pre-change journal replay, unchanged `attentionSemanticKey`), on the argument
   that the latter two follow from the first. If you judge those independently
   failable, say so and I will split it rather than argue.
3. **Are VAL-081 and VAL-082 executable exactly as task 9 already reads?** They
   were written to add **no** implementation work and **no** command: each maps
   onto a criterion the spec already carried. Confirm that, and confirm the two
   evidence classes cannot be confused - VAL-080 is satisfied only by shape
   rejection, VAL-081 only by `condition nullable binding mismatch`. Also
   confirm I-44's clause makes VAL-082's `budget-observed` and
   `capability-available` cases reachable rather than relation rejections.
4. **Coverage and pointers.** Every row cited by at least one acceptance
   criterion; exactly one validator per row; the count 82 = 78 deterministic + 4
   Human in the header, the numbering note, and the brief; the
   `contract-coverage.md` clause rows; and the three I-43 staleness corrections
   plus the I-42 annotation. And the standing question: **does anything here
   reach tasks 1-8?** I believe not - the edits touch the contract, tasks 9-11,
   and derived records only - but check it by name, as both pass-8 sessions did.

## 4. What I Am Not Asking For

Not a re-derivation of the registry inventory, not a re-run of any gate, not a
Human-fixture change, and not a re-opening of DEC-235, DEC-236, or ESC-23-001.
Your own QA-plan mapping - section 12.9's QA-43 replacement and section 20.3's
route table - now maps to two rows instead of one; that file is yours, and
updating it is your call, not a Tech Lead instruction.

If you find a third bundled row anywhere in the contract while you are in there,
raise it. The pattern is a row that asserts both that a thing behaves correctly
and that it is encoded correctly, and this phase has now produced two.
