# TL Escalations — Current

## Open escalations

None newly minted. ESC-23-004 is resolved in substance (below). The sequence
verifier had no valid route out of the STOP: review-6 stayed open with no next
owner after the review-7 BLOCK. DEC-256 (Tech Lead) reshapes the repairs to
four distinct owners. DEC-257 (Human, Option A) lets the RUN-RESUME carry them.
DEC-257 implemented; Runner resumes the DEC-256 repairs. No new TRG/ESC: it is
the same STOP and the same authority question.

## Recently resolved

### ESC-23-004 - Post-task-11 checkpoint re-review: three protected blocker candidates (BC-1, BC-2, BC-3)

- **Raised by:** TRG-23-004
- **Status:** RESOLVED 2026-09-26.
  - BC-1 by DEC-249 and BC-3 by DEC-250 (Tech Lead).
  - BC-2 by DEC-251 (Human, Option B).
  - Carry items: I-1 by DEC-254 (Part 2 Human-approved), I-2 by DEC-253, I-5 by DEC-252.
  - DEC-255 releases the ordered repairs: MF-4, task-2 BC-1, MF-5+I-6, task-11 BC-3, task-6 I-5.
  - Next: a fresh `checkpoint-after-task-11-3.md` re-review.
  - Task 12 stays held until that re-review returns PASS with Must-fix 0 (qa-plan section 8).
- **Originating task:** task-10
- **Trigger signature:** task-10:qa-rework:review-7-BC-1-BC-2-BC-3
- **Source:** `project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11-2.md` (review-7, commit 9cb9420516c461ddd9b4d22d1944759a416f7bb7)
- **Category:** specification, product semantics and public interface. Each sub-item is separately resolvable.
- **Sub-item BC-1 (product semantics, VAL-009):** contracts.md:91-92 says a guarded action requires "every applicable Obligation" to be satisfied, but "applicable" is undefined. contracts.md:60-63 and :408-409 imply per-Obligation guarded_actions, while the evaluator (internal/workflow/evaluation.go:436,496-512; envelope.go:521,546-556) blocks terminal success on any unsatisfied Obligation regardless of its guarded_actions. Decision needed: does an Obligation block only its declared guarded_actions? Option A (yes): a decision-free task-2 evaluator conformance repair, then VAL-009 (both halves), the journey stage 1-3 `qa_pass` assertions and task 13 VAL-056 are re-derived. Option B (no, global): record the reading in a DEC and a GLOSSARY entry, including what guarded_actions governs. The task-10 stage 3 assertion then stands and MF-5's `qa_pass` expectation is fixed. Affects tasks 2, 10, 13, 14 (VAL-059, VAL-071).
- **Sub-item BC-2 (public interface / authorization, VAL-043):** public-interface.md:243 (section 4.1) requires an Envelope grant for `human-outcome.record`. public-interface.md:222-224 and contracts.md:391-396 do not, and the product and task-9 journey need none (internal/review/human_outcome.go:79-89). Decision needed: which text governs. Option A (grant required): a task-10 authorization repair plus a VAL-043 negative case, and the task-9 journey fixture gains a grant for `actor-human-outcome`. Option B (grant not required): correct the section 4.1 cell under a DEC. Human ruling needed if the interface changes (DEC-222). Affects tasks 9, 10, 14 (VAL-071) and 18 (VAL-067/068).
- **Sub-item BC-3 (resolution semantics, VAL-048, VAL-050, VAL-078):** internal/workflow/stop.go `eventMatchesCondition` and `qualifyingFact` have no rule for `human-outcome-recorded` or `decision-satisfied`, so a Stop carrying such a Condition raised by a non-producer code can never be resolved. contracts.md:453-455 and :461-462 fix that a `human-outcome-recorded` Stop resolves with a class-and-result-matching Outcome Event and that matching Stops resolve. They do not fix the `decision-satisfied` qualifying fact (which fact of the pair, Envelope scope), subject scoping for `human-outcome-recorded`, or whether recording resolves every Condition-matching Stop regardless of raising code. Decision needed: the Tech Lead fixes both qualifying rules. Task 11 then implements them in its stop.go scope and evidences the VAL-048 reject clause and VAL-078 per type. Affects task 11 and task 17 (`human_outcome_resolution:true`).
- **Independent bounded Must-fix items (not blocked by the candidates unless noted):** MF-4 (task-6, VAL-027 linked-recovery leg, test-only, independent of all three candidates, may dispatch once the Tech Lead releases the stop). MF-5 (task-10, VAL-009 evidence-method state plus I-6 citations, test/report-only, its `qa_pass` assertion waits on BC-1). Each has an owning task, bounded scope, required evidence and a fresh independent re-review in the review's Repair Plan.
- **Carry items:** I-1 `cmd/fid` headroom, I-2 task-11 scope ratification, I-5 multi-grant `gate.run` owner, I-6 report citations. No action here.
- **Constraints:** Do not choose any reading in code, tests or reports before the decision is recorded. Do not narrow VAL-009, VAL-043, VAL-048, VAL-050 or VAL-078. A fresh independent checkpoint re-review must PASS before task 12.
- **References:** TRG-23-004 in `project-docs/docs/phases/phase-23/triage-log.md`; review-7 lines 661-776.

### ESC-23-003 - VAL-015 Run-holder-only CLI evidence is unrepresentable under the accepted lifecycle

- **Raised by:** TRG-23-003
- **Status:** RESOLVED 2026-09-20 under DEC-239. The Human delegated resolution
  through Task 9 and its review. Retain Run/Claim Actor binding and durable
  history; amend only the impossible Run-holder-only CLI evidence route to an
  explicitly isolated query fixture. Same-Reviewer CLI rework and the other
  Task 4 obligations remain required before completion.
- **Originating task:** task-4
- **Category:** specification contradiction; protected data-model/public-interface and validation-contract authority
- **Condition:** VAL-015 requires a real CLI proof that an Actor holding a Run under another Actor's mutating Claim is rejected for non-independence. The accepted `run.start` binding requires the Claim Actor and Run Actor to be equal; durable Claim history is retained and checked before the Run can be named. The state is therefore unconstructable through the accepted CLI lifecycle.
- **Evidence:** `project-docs/docs/phases/phase-23/tasks/task-4-done.md` records the attached real-CLI rework fixture, direct durable-history Run-holder query, and typed `OUTCOME BLOCKED`. The invoking evidence verifies implementation commit `a1ae1bb61363a591bb46c9ddfd4b8df36462687b` and report lint. The five exact task gates pass, but that does not supply the unrepresentable real CLI observation.
- **Original decision request (retained):** Decide whether to authorize a compatible Run/Claim actor-binding or durable-history contract change, or amend VAL-015 to accept the direct durable-history proof for the Run-holder-only direction while retaining the required independence behavior.
- **Constraints:** Do not weaken or silently waive VAL-015, alter actor-binding/history precedence, or expand task scope without the required Tech Lead/Human authority.
- **References:** TRG-23-003; `project-docs/docs/phases/phase-23/contracts.md:112-130`; `project-docs/docs/phases/phase-23/tasks/task-4.md:16-19`; `project-docs/docs/phases/phase-23/qa/validation-contract.md` VAL-015; DEC-225.

Resolution and exact Runner binding:
[DEC-239](../phases/phase-23/decisions.md),
[resume handoff](../phases/phase-23/qa/resume-after-esc-23-003.md).

ESC-23-002 is **RESOLVED 2026-09-20 under DEC-238**. The Human delegated
bounded repair authority above the fixed agent roles. Rebuild statement reuse
and fresh/isolated test setup repair the causes without changing any gate or
assertion. Fresh exact gates pass: CLI **527.099s**, workflow **231.010s**,
tagged integration **669.586s**; rebuild samples **22.728s / 22.623s / 22.489s**.
Integration binary directories stay **0 → 0**. Task 1 is complete in its report;
the Runner must reconcile the retained STOP before advancing. Tasks 9-11 stay
held. See [repair evidence](../phases/phase-23/qa/esc-23-002-repair.md) and
[resume handoff](../phases/phase-23/qa/resume-after-esc-23-002.md).

ESC-23-001 is **resolved**: the Human chose Option A on 2026-09-20 - `omitempty`
on the three new `Condition` fields only, the six Phase 22 types keeping
`budget_name` and `capability_id` present-and-null, with per-type journal shapes
keyed on `type`. Recorded as DEC-235, which amends `public-interface.md` section
2, restates inventory row 8, creates row 10, and rewrites tasks 9, 10, and 11.
The three held revisions (R-22, R-23, I-37) shipped with it. **ESC-23-001
blocks no task**. DEC-239 supplies Human clearance for Task 9, and qa-plan
section 21 records DEC-237/240 bounded QA readiness; Tasks 10-11 remain held. See
[Phase 23 escalations](../phases/phase-23/tl-escalations.md).

QA pass 8 raised one **blocker candidate with no identifier allocated** (R-28,
VAL-078's granularity), which is the correct handling: the repair is a Tech Lead
contract edit, not a Human decision. Dispositioned under DEC-237 by splitting
the row. No escalation was or should have been opened.

ESC-22-005 is resolved by the Human's performance-repair approval under
DEC-221; its proposed timeout increase was not approved or applied. Task 20 is
complete; hosted run34805256766 passed on6b8fe50. See
[Phase 22 escalations](../phases/phase-22/tl-escalations.md) for all five records.
The wrap's per-TRG hindsight audit is in
[wrap-up-audit.md](../phases/phase-22/wrap-up-audit.md).

## Standing handling

Register an escalation when allocated, not after resolution. Ordinary causal
repair and incomplete in-scope work do not themselves require a Human stop.
Read accepted ownership before calling it ambiguous (DEC-194). A real contract/
schema contradiction needs the missing authority, not a weakened assertion
(DEC-195). External interruptions do not consume repair no-progress budget.

Prior archives: [Phase 21](../phases/phase-21/tl-escalations.md),
[Phase 20](../phases/phase-20/tl-escalations.md),
[Phase 19](../phases/phase-19/tl-escalations.md).
