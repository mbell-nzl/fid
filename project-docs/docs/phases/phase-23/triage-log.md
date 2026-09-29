# Phase 23 - TL-Lite Triage Log

> Seeded by Tech Lead from the triage-log template at planning kickoff.
> Append-only allocation source for TRG-23 and ESC-23 identifiers.

No triage events have occurred. The first allocated identifiers will be
TRG-23-001 and ESC-23-001, in their separate series. Read this file before
allocation; never recall, reserve in chat, or reuse an identifier.

Only TL-Lite or Tech Lead allocates an ESC. Write its durable entry and the
current escalation-register row before reporting a BLOCK. Ordinary REPAIR
allocates no ESC and requests no Human attention.

Actual entries use `project-docs/templates/triage-log-template.md` and the
current schema-v2 exact typed fields. Do not paste placeholder TRG sections
into this file. Append corrections instead of editing earlier verdicts.
Runner-owned run-ledger state is created only after execution approval.

Standing Phase 23 reminders carried from the Phase 22 wrap audit: read accepted
ownership before calling it ambiguous (DEC-194); a contract clause that lacks a
schema row needs the missing authority, never a weakened assertion (DEC-195);
external interruptions and mis-scoped directives never consume a no-progress
recovery budget (DEC-170/171).

## TRG-23-001 - Task 1 cannot close with benchmark and cross-scope race-gate failures

- **Trigger:** Phase 23 task 1 cannot close after a material-progress continuation because two exact required gates remain failed: the 100k rebuild benchmark exceeds the accepted 30-second budget, and the repository-wide race gate fails both the accepted 600-second ordinary-package budget and a cmd/fid execution-transition assertion outside task 1 scope.
- **Originating task:** task-1
- **Trigger signature:** `task-1:gate-repair:final-required-gates`
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** `project-docs/docs/phases/phase-23/tasks/task-1-done.md`
- **Candidate:** RECOVERY_REQUIRED
- **Attempt count:** 2
- **Checkpoint:** verified task commit `0ced511b2e6a7e1f252260024176705106e9b793` predates phase base `a0e8e7ea1111e2a4caa6c39bc6d1e63dee82accf`; runner-recorded control repair base `150569adb859cbd28b177ea1950f6dfdde6cd51d`
- **Verdict:** BLOCK
- **Reasoning:** The benchmark remains 32.500159404 s against its accepted strict under-30-second assertion after the permitted performance experiments; no accepted local repair is evidenced. The ordinary race gate separately exceeds the DEC-227/VAL-079 600-second budget (`cmd/fid` 619.168 s; `internal/workflow` 614.196 s) and retains a `cmd/fid` assertion that task 1 neither owns nor is proven to have displaced. DEC-227 requires a pre-existing cause to be escalated rather than absorbed, and changing a budget or assertion would be protected.
- **Local decisions made:** none
- **Resolution evidence:** Exact benchmark command failed with `FAIL github.com/mbell-nzl/fid/internal/journal 35.245s` (terminated `exit_code: 1`) and measured 32.500159404 s/rebuild. Exact repository-wide race command failed (terminated `exit_code: 1`): `cmd/fid` 619.168 s, `internal/workflow` 614.196 s, and `cmd/fid/task12_execution_transition_test.go:217` observed `claim-expired`, `run-inactive`, and `workflow-transition-invalid` where the retained assertion expects only `workflow-transition-invalid`. Report lint returned `OUTCOME RECOVERY_REQUIRED` (terminated `exit_code: 0`).
- **If REPAIR:** N/A.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Do not change any command, timeout, threshold, budget, assertion semantics, accepted status/scope, public interface, task/QA/spec/decision definition, or route the out-of-scope cmd/fid assertion to task 1 without Tech Lead authority.
- **Exact verification:** N/A
- **Originating report:** project-docs/docs/phases/phase-23/tasks/task-1-done.md
- **Control repair base:** 150569adb859cbd28b177ea1950f6dfdde6cd51d
- **If BLOCK:** category protected gate/control authority and cross-task assertion ownership; ESC-23-002 raised in `tl_escalations_current.md`; resolution condition Tech Lead assigns a bounded owner and repair authority for the retained `cmd/fid` assertion and any non-task-1 budget cause, or a Human/accepted decision changes the protected gate definition.
- **Human attention required:** yes
- **References:** `project-docs/docs/phases/phase-23/tasks/task-1.md` Exact Commands To Run; `project-docs/docs/phases/phase-23/tasks/task-1-done.md` Workflow outcome and Exact quality-gate results; `project-docs/docs/phases/phase-23/execution-scope.md:44-45,501-508`; `project-docs/docs/phases/phase-23/decisions.md` DEC-227; `project-docs/docs/phases/phase-23/qa/validation-contract.md` VAL-079; `project-docs/specs/09_AI_WORKFLOW.md` sections 2.3 and 9.

---

## TRG-23-002 - Corrective schema-v2 external record for TRG-23-001

- **Trigger:** schema-v2 correction for malformed external TRG-23-001; the prior BLOCK record wrote Markdown backticks into the typed Trigger signature, so the Runner's exact TRIAGE event is mechanically impossible to append while preserving both the external byte identity and the required task-binding check.
- **Originating task:** task-1
- **Trigger signature:** task-1:gate-repair:final-required-gates
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** TRG-23-001 and ESC-23-002
- **Candidate:** schema-v2 external-record correction
- **Attempt count:** 2
- **Checkpoint:** runner-recorded control repair base 150569adb859cbd28b177ea1950f6dfdde6cd51d; ledger unchanged
- **Verdict:** BLOCK
- **Reasoning:** TRG-23-001 preserves the original protected gate/control and retained-assertion ownership BLOCK, but its Markdown-delimited Trigger signature cannot satisfy the schema-v2 plain-scalar task binding and external-record byte-identity requirements at the same time. This append-only corrective record supplies the one mechanically valid external identity with the unchanged originating task, trigger signature, and attempt count; it neither changes gate definition nor assigns ownership.
- **Local decisions made:** allocated this corrective TRG only; retained the original BLOCK and ESC-23-002 without allocating a duplicate ESC.
- **Resolution evidence:** Reproduced mismatch: a byte-for-byte ledger candidate was rejected with `INVALID schema v2 TRIAGE Originating task does not match Trigger signature`; a plain-scalar candidate was rejected with `INVALID schema v2 TRIAGE identity does not match its unique external record`. The ledger remains unchanged.
- **If REPAIR:** N/A.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Do not edit or delete TRG-23-001, ESC-23-002, the run ledger, product code, task/QA/spec/decision files, or any protected gate/assertion. Preserve the original BLOCK record and escalation.
- **Exact verification:** project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-23/run-ledger.md after a corrective external triage record exists
- **Originating report:** project-docs/docs/phases/phase-23/tasks/task-1-done.md
- **Control repair base:** 150569adb859cbd28b177ea1950f6dfdde6cd51d
- **If BLOCK:** category protected gate/control authority and cross-task retained-assertion ownership; ESC-23-002 referenced in `tl_escalations_current.md`; resolution condition full Tech Lead assigns bounded causal repair ownership and authority, or accepted protected authority changes the gate definition.
- **Human attention required:** yes
- **References:** TRG-23-001 in this log; ESC-23-002 in `project-docs/docs/current/tl_escalations_current.md`; `project-docs/docs/phases/phase-23/tasks/task-1.md` Exact Commands To Run; `project-docs/docs/phases/phase-23/tasks/task-1-done.md` Workflow outcome and Exact quality-gate results; `project-docs/specs/09_AI_WORKFLOW.md` sections 2.3 and 9.

---

## Resolution of ESC-23-002 under DEC-238 — 2026-09-20

The Human explicitly delegated bounded authority above the fixed roles to
complete the causal repair, close the escalation, commit, and prepare resume.
DEC-238 assigns that ownership; all accepted gate definitions and assertions
remain unchanged. The original TRG-23-001 and corrective TRG-23-002 BLOCK
records above remain intact. No new TRG or escalation is allocated.

ESC-23-002 is **RESOLVED 2026-09-20 under DEC-238**. The Human delegated
bounded repair authority above the fixed agent roles. Rebuild statement reuse
and fresh/isolated test setup repair the causes without changing any gate or
assertion. Fresh exact gates pass: CLI **527.099s**, workflow **231.010s**,
tagged integration **669.586s**; rebuild samples **22.728s / 22.623s / 22.489s**.
Integration binary directories stay **0 → 0**. Task 1 is complete in its report;
the Runner must reconcile the retained STOP before advancing. Tasks 9-11 stay
held. See [repair evidence](qa/esc-23-002-repair.md) and
[resume handoff](qa/resume-after-esc-23-002.md).

Resume uses the existing TRG-23-002 / ESC-23-002, authority DEC-238, and next
attempt 3 for `task-1:gate-repair:final-required-gates`. This resolves the
protected ownership condition into the authorized task-scope recovery; it is
not a control-conformance directive or a rewrite of either historical verdict.

---

## TRG-23-003 - VAL-015 real CLI Run-holder-only evidence contradicts accepted lifecycle

- **Trigger:** Phase 23 Task 4 VAL-015 requires a real CLI Run-holder-only independence case that the accepted Run/Claim lifecycle cannot represent.
- **Originating task:** task-4
- **Trigger signature:** task-4:in-scope-implementation:review-record-independence-tracer
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** `project-docs/docs/phases/phase-23/tasks/task-4-done.md`
- **Candidate:** specification-contradiction
- **Attempt count:** 8
- **Checkpoint:** implementation commit `a1ae1bb61363a591bb46c9ddfd4b8df36462687b` verified by attached `git rev-parse --verify` evidence (terminated `exit_code: 0`) with subject `test(review): cover independent reviewer rework`; attached report-lint evidence verified the committed report's typed `OUTCOME BLOCKED` (terminated `exit_code: 0`); runner-recorded control repair base `eee11a2a1964e7802371e33cdcbf872c233c41c7`.
- **Verdict:** BLOCK
- **Reasoning:** The accepted lifecycle requires `claim.ActorID == options.ActorID` when a Run is started. Durable mutating Claim history remains authoritative for independence and is checked before the Run can be named, so the required real CLI case of an Actor holding only a Run under another Actor's mutating Claim cannot be constructed. The remaining choices either change the accepted Run/Claim actor-binding and durable-history semantics or change VAL-015's accepted evidence/pass criterion; both are protected data-model/public-interface or validation-contract decisions.
- **Local decisions made:** none.
- **Resolution evidence:** The attached reproduced mismatch and the completion report establish that the direct durable-history query covers the intended Run-holder row, while the admissible real CLI fixture proves the two-Reviewer rework journey and all five exact task gates pass. This does not clear VAL-015's required real CLI observation. The cited commit/report evidence confirms the attempted implementation and its typed BLOCKED outcome, not a resolution.
- **If REPAIR:** N/A.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** changing the accepted Run/Claim actor-binding semantics, durable Claim-history precedence, VAL-015 acceptance, public/data-model contract, or task scope without Tech Lead/Human authority.
- **Exact verification:** `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordIndependentReworkAndReceiptBinding$' -count=1` from `/home/mark/projects/fid`; attached observed final token `ok github.com/mbell-nzl/fid/cmd/fid 25.479s` (terminated `exit_code: 0`), with the direct Run-history query proof in `internal/review` and the real CLI Run-holder-only row remaining unrepresentable.
- **Originating report:** project-docs/docs/phases/phase-23/tasks/task-4-done.md
- **Control repair base:** eee11a2a1964e7802371e33cdcbf872c233c41c7
- **If BLOCK:** category specification contradiction requiring protected data-model/public-interface and validation-contract authority; ESC-23-003 raised in `tl_escalations_current.md`; resolution condition a Tech Lead/Human decision either authorizes a compatible Run/Claim lifecycle change or amends VAL-015's accepted evidence requirement while preserving the intended independence guarantee.
- **Human attention required:** yes
- **References:** `project-docs/docs/phases/phase-23/tasks/task-4.md:16-19,153-154`; `project-docs/docs/phases/phase-23/tasks/task-4-done.md:5-10,24-25,44-50`; `project-docs/docs/phases/phase-23/contracts.md:112-130`; `project-docs/docs/phases/phase-23/qa/validation-contract.md` VAL-015 and DEC-225 refinement; `project-docs/specs/09_AI_WORKFLOW.md` sections 2.3 and 9.

---

## TRG-23-004 - Post-task-11 checkpoint re-review BLOCK (BC-1, BC-2, BC-3, MF-4, MF-5)

- **Trigger:** QA BLOCK - the fresh independent QA Lead re-review of the post-task-11 checkpoint (ledger review-7, commit 9cb9420516c461ddd9b4d22d1944759a416f7bb7) returned VERDICT: BLOCK with three protected blocker candidates (BC-1, BC-2, BC-3) and two bounded Must-fix evidence gaps (MF-4, MF-5). qa-plan section 8 requires the checkpoint to PASS before task 12 is dispatched, so task 12 is held.
- **Originating task:** task-10
- **Trigger signature:** task-10:qa-rework:review-7-BC-1-BC-2-BC-3
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** `project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11-2.md` (re-review of `qa/checkpoint-after-task-11.md`)
- **Candidate:** QA BLOCK / protected blocker candidates
- **Attempt count:** 1
- **Checkpoint:** QA review commit 9cb9420516c461ddd9b4d22d1944759a416f7bb7 verified by attached evidence (rev-parse exit 0, subject match, one added file, final lines end `VERDICT: BLOCK`); the task-10 repair (test af0cbb2c1448844ec7024cd823631533de655945, reports 90d921b8084e58f2415c812bf789128d9a343195 and 908080d10c46f4ca9736f2a914035a3191a4ed8e) is verified. The QA Lead's built-binary probe outputs were not verified by the Runner or by me; my dispositions rest on the accepted contract text read directly, and I take the QA code citations as reported.
- **Verdict:** BLOCK
- **Reasoning:** Part (a) is verified, so the task-10 repair MF-1..MF-3, I-4 stands. Each candidate was judged against the accepted text. BC-1 is class (ii). contracts.md:91-92 says "A guarded action requires every applicable Obligation to be `satisfied`". contracts.md:60-63 constrains which Obligations may guard which action, and contracts.md:408-409 has the readiness Stop block "exactly that Obligation's guarded actions". "Applicable" is undefined, so the text supports both a per-Obligation reading (blocks only its declared guarded_actions) and a global reading (any unsatisfied Obligation blocks every guarded action). The shipped evaluator uses the global reading (internal/workflow/evaluation.go:436,496-512, per QA). The accepted text does not uniquely fix which reading governs, so this is a product-semantics and acceptance-criteria ambiguity (VAL-009) that neither I nor the Developer may choose. BC-2 is class (ii). public-interface.md:243 (section 4.1) requires "plus Envelope grant" for `human-outcome.record`. public-interface.md:222-224 (section 4) and contracts.md:391-396 (section 10) list only actor_type human, capability, and Activation provenance. This is a direct text contradiction inside a Human-accepted public interface (DEC-222). It is a public-interface and authorization/security rule, so the governing text must be ruled by the Tech Lead or Human. BC-3 is class (ii), partly determined. contracts.md:453-455 fixes that a `human-outcome-recorded` Stop resolves only with the Event ID of a Human Outcome whose class and result match, and contracts.md:461-462 says only Stops whose recorded condition matches are resolved. That indicates the shipped absence of any qualifying rule for the type is a defect. Two things stay unfixed by the text. First, the qualifying fact for `decision-satisfied` (which fact of the pair, and Envelope scope) is not stated anywhere I could find. Second, subject scoping for `human-outcome-recorded`, and whether recording resolves every matching-Condition Stop regardless of raising code, is not stated. Because a resolution-semantics rule and the closed registry of qualifying facts are missing (the phase brief's recurring registry family, which is escalated), this needs a rule fixed before task 11 implements it. VAL-048's "reject" clause and the VAL-078 per-type route wait on it. MF-4 and MF-5 each name an owning task, bounded scope, required evidence and a fresh independent re-review in the review's Repair Plan table (lines 656-659) and are decision-free in themselves. MF-4 (task-6, test-only, VAL-027 linked-recovery leg) is independent of all three candidates. MF-5 (task-10, test/report-only, VAL-009 evidence-method state plus I-6 citations) has its `issue.transition-terminal-success` `qa_pass` expectation waiting on BC-1. The BLOCK cannot be waived by the Runner or QA: qa-plan section 8 keeps task 12 held until the checkpoint PASSes.
- **Local decisions made:** none. I decided no protected question. Recorded: MF-4 and MF-5 are decision-free; I-1, I-2, I-5 and I-6 are carry items with no action here.
- **Resolution evidence:** Not resolved. The reading of the contract text above, plus the review's Blocker Candidates section (lines 661-776), is the evidence. Resolution needs an explicit Tech Lead decision for each of BC-1..BC-3 (Human where it changes the interface or a product rule), recorded as DECs.
- **If REPAIR:** N/A.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Do not dispatch task 12. Do not choose the BC-1 reading, the BC-2 governing text, or the BC-3 qualifying rules in code, tests, or reports before a decision is recorded. Do not narrow VAL-009, VAL-043, VAL-048, VAL-050 or VAL-078. Do not edit the review, contracts, public-interface, validation-contract, or DEC files without Tech Lead authority. Do not add a workaround grant or qualifying rule.
- **Exact verification:** N/A
- **Originating report:** project-docs/docs/phases/phase-23/tasks/task-10-done.md
- **Control repair base:** 9cb9420516c461ddd9b4d22d1944759a416f7bb7
- **If BLOCK:** umbrella ESC-23-004 raised in `tl_escalations_current.md` with three separately resolvable sub-items. BC-1 asks for a Tech Lead ruling (Human if a product rule) on whether an Obligation blocks only its declared guarded_actions. If yes, a decision-free task-2 evaluator conformance repair follows, and VAL-009, the journey stage 1-3 `qa_pass` assertions and task 13 VAL-056 are re-derived. If no, the reading goes into a DEC and a GLOSSARY entry, and the task-10 stage 3 assertion stands. BC-2 asks the Tech Lead (Human if the interface changes) which text governs. If the grant is required, a task-10 authorization repair and a VAL-043 negative follow. Otherwise the section 4.1 cell is corrected under a DEC. BC-3 asks the Tech Lead to fix the `stop resolve` qualifying rules for `human-outcome-recorded` and `decision-satisfied`, with subject scoping and the match-regardless-of-code question. Task 11 then implements them in its stop.go scope and evidences the VAL-048 reject clause and VAL-078 per type. Forward tasks affected: 13 (VAL-056), 14 (VAL-059, VAL-071), 17 (`human_outcome_resolution:true`), 18 (VAL-067/068). MF-4 is dispatchable independently once the Tech Lead releases the stop, and MF-5 completes after BC-1. A fresh independent checkpoint re-review is required afterwards.
- **Human attention required:** yes
- **Carry items (no action here):** I-1 `cmd/fid` at 567.662 s of the 600 s budget with a 39.3 s same-code spread, for a Tech Lead decision before tasks 13 and 14 (budget unchanged). I-2 ratify task-11's stop_conditions.go scope reading in decisions.md. I-5 multi-grant `gate.run` (internal/evidence/execution.go:381), a latent Phase 21 conformance defect needing a Tech Lead owner before task 14 VAL-071 or task 18. I-6 report citation drift, folded into the MF-5 continuation.
- **References:** `project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11-2.md` lines 410-780; `project-docs/docs/phases/phase-23/contracts.md:60-63,91-92,391-396,407-411,453-462`; `project-docs/docs/phases/phase-23/public-interface.md:220-243`; qa-plan section 8; DEC-222, DEC-224, DEC-227; `project-docs/specs/09_AI_WORKFLOW.md` sections 2.3 and 9.

---

## TRG-23-005 - Task-11 ordinary gate exceeds accepted package-time budget without a bounded task-11 cause

- **Trigger:** task-11's unchanged exact ordinary race/coverage gate terminates exit_code 0 with tests passing but breaches the accepted cmd/fid package-time budget. Attempt 1 measured cmd/fid 686.386s; attempt 2 measured cmd/fid 630.927s after clearing the test cache; a fresh same-session pre-repair parent measured 614.601s. The exact tagged integration gate passed at 766.070s under its 900s budget. Task-11's accepted changes add no ordinary cmd/fid test file; focused TestStop timing was 109.999s on the parent and 112.324s current. Unrelated packages also rose. The completion report concludes no task-11 causal speed repair is established and requests bounded disposition of the remaining scope gap.
- **Originating task:** task-11
- **Trigger signature:** review-7:qa-rework:task-11-BC-3
- **Attempt count:** 4
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** `project-docs/docs/phases/phase-23/tasks/task-11-done.md`
- **Candidate:** RECOVERY_REQUIRED / protected gate-repair scope gap
- **Checkpoint:** attached verification of implementation commit `466745f2dc4505d60eb67ac50193e004532b7239` (rev-parse and subject match); attached report-lint `OUTCOME RECOVERY_REQUIRED` and `PASS report-lint completed with no failures` (terminated `exit_code: 0`); report commit `c1588fc4389c7dc8210247fc48ede2618c783799` supplied by Runner. These verify the checkpoint and report state, not the package-time budget.
- **Verdict:** BLOCK
- **Reasoning:** The accepted exact ordinary gate still breaches DEC-227 and VAL-079 despite test-process exit 0. The fresh same-session parent also exceeds 600 s; task 11 added no ordinary `cmd/fid` test file, its focused Stop timing increased only 2.325 s, and other packages slowed. This evidence does not identify a task-11 causal speed repair. Task-11's spec assigns speed repair to this task only when its own new journeys caused the breach and otherwise requires escalation. DEC-254 already authorizes a separate task-13 headroom step, but it does not clear task-11's failed gate or authorize TL-Lite to transfer its closure obligation. The gate and budget definition remain protected.
- **Local decisions made:** none; no new owner, budget, threshold, command, or pass criterion assigned.
- **Resolution evidence:** Not cleared. Exact `go test -race -timeout 30m -coverprofile=coverage.out ./...` produced `cmd/fid` 686.386 s and 630.927 s (both terminated `exit_code: 0`); the same-session parent produced 614.601 s (terminated `exit_code: 0`). The tagged integration gate produced 766.070 s, below its separate 900 s budget. Full Tech Lead disposition is needed for the out-of-task-11 speed cause and task-11 gate closure; any changed gate definition needs its protected decision authority. Preserve every failed timing attempt.
- **If REPAIR:** N/A.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Do not change the exact command, accepted 600-second ordinary or 900-second tagged budgets, assertion semantics, or policy; do not weaken tests, claim a later flaky green as clearing preserved failures, or assign a new task owner locally.
- **Exact verification:** N/A
- **Originating report:** project-docs/docs/phases/phase-23/tasks/task-11-done.md
- **Control repair base:** N/A
- **If BLOCK:** category protected gate-repair scope and ownership authority; ESC-23-005 raised in `tl_escalations_current.md`; resolution condition full Tech Lead records a bounded owner and closure route under the unchanged DEC-227/VAL-079 contract, or the authority owning the gate definition accepts an explicit change. Task 12 remains held pending the required fresh checkpoint PASS.
- **Human attention required:** yes
- **References:** `project-docs/docs/phases/phase-23/tasks/task-11.md` Exact Commands To Run; `project-docs/docs/phases/phase-23/tasks/task-11-done.md` Workflow outcome, Exact quality-gate results, and Concerns; `project-docs/docs/phases/phase-23/decisions.md` DEC-227 and DEC-254; `project-docs/docs/phases/phase-23/qa/validation-contract.md` VAL-079; `project-docs/specs/09_AI_WORKFLOW.md` sections 2.3 and 9.

## TRG-23-006 - Task 12 requirement authority and report-lint contract conflict

- **Trigger:** Phase 23 task 12 cannot reach its real-process requirement reference journey because the accepted capability contract and closed activation registry disagree; the task spec's expected-nonzero declaration also fails the canonical report-lint grammar.
- **Originating task:** task-12
- **Trigger signature:** task-12:protected-contract:requirement-capability-registry-and-report-lint-grammar
- **Attempt count:** 1
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-23/tasks/task-12-done.md
- **Candidate:** RECOVERY_REQUIRED / protected contract contradiction
- **Checkpoint:** Runner-supplied verification of implementation commit 0e4bb33cbe32fa9d4d83955b3121a4135edf78cd and report commit 32848694d6f46cf8f8bd7ed17d1e5425efb11a17 by git rev-parse --verify and subject match; independently observed report-lint exit_code: 1 with `FAIL expected non-zero declaration has invalid grammar`. Runner-recorded phase HEAD: 32848694d6f46cf8f8bd7ed17d1e5425efb11a17.
- **Verdict:** BLOCK
- **Reasoning:** Public-interface section 4 requires matching Actor-profile capabilities for the two requirement writes, while section 4.1, execution-scope section 6, and the accepted R-19 disposition explicitly exclude both names from the closed 30-name capability set. Adding them to activation would change that accepted public authority contract; removing the Actor-profile requirement would also change it. Neither is a uniquely determined control implementation repair. Separately, task-12's two-line expected-nonzero declaration contradicts the canonical report-lint parser's same-bullet grammar; changing the task declaration or accepted control grammar needs its owning authority. Incomplete task work remains for same-task continuation after resolution, but does not itself justify BLOCK.
- **Local decisions made:** none; no capability list, authority semantics, task-spec declaration, or report-lint grammar authorized for change.
- **Resolution evidence:** Not cleared. `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` failed twice at Activation staging with `invalid activation input` (each terminated exit_code: 1), before the requirement writer. `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-12-done.md project-docs/docs/phases/phase-23/tasks/task-12.md` failed with `FAIL expected non-zero declaration has invalid grammar` (terminated exit_code: 1). Attached verified commits establish the checkpoint, not resolution of either mismatch.
- **If REPAIR:** N/A.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A; protected contract contradiction requires a decision.
- **Control mismatch:** N/A; both reproduced failures are recorded above.
- **Allowed files:** N/A
- **Forbidden changes:** Do not change internal/workflow/activation.go, the closed capability set, task-12.md, project-docs/tools/report-lint.sh, product/spec/QA/ledger/decision files, or accepted authority and gate semantics without the owning protected authority. Preserve failed attempts.
- **Exact verification:** N/A pending protected authority.
- **Originating report:** project-docs/docs/phases/phase-23/tasks/task-12-done.md
- **Control repair base:** N/A; supplied phase HEAD is recorded in Checkpoint, and no control directive is granted.
- **If BLOCK:** category public-interface/specification contradiction and protected task/control definition authority; ESC-23-006 raised in tl_escalations_current.md. Resolution condition: full Tech Lead obtains the necessary Human/contract-owner ruling on requirement write authority and the closed capability set, plus the task-spec or report-lint grammar authority needed for the expected-nonzero declaration. Then Runner resumes task 12 under the recorded resolution and reruns the real-process and report-lint commands before ordinary completion gates.
- **Human attention required:** yes
- **References:** project-docs/docs/phases/phase-23/public-interface.md sections 4 and 4.1; project-docs/docs/phases/phase-23/execution-scope.md section 6; project-docs/docs/phases/phase-23/decisions.md R-19; project-docs/docs/phases/phase-23/tasks/task-12.md Scope, Expected Non-Zero Commands; project-docs/tools/report-lint.sh:104; project-docs/docs/phases/phase-23/tasks/task-12-done.md Tests run and Deviations; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9.

## TRG-23-007 - Task 13 ActivityEvent producers cross the journal and migration boundary

- **Trigger:** Task 13's accepted real-process ActivityEvent behavior is blocked by journal schema/effect ownership: activity_log is rejected as an unknown table, and the requirement.record-reference producer cannot reference the requirement_refs row under the current trigger schema.
- **Originating task:** task-13
- **Trigger signature:** task-13:scope-gap:review-list-and-activity-producers
- **Attempt count:** 1
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-23/tasks/task-13-done.md
- **Candidate:** RECOVERY_REQUIRED / protected persistence and task-order scope gap
- **Checkpoint:** Runner-supplied commit-verification evidence: `git rev-parse HEAD^{commit}` returned 007c9487863f5eb3cf037a0f0affd290742fed56, subject matched `docs(task): correct phase-23 task-13 blocker evidence format`; report-lint returned `OUTCOME BLOCKED` (terminated `exit_code: 0`). Task-13 recovery attempt 1 added `review.list`; the ActivityEvent probe and full tagged gate remain failed.
- **Verdict:** BLOCK
- **Reasoning:** VAL-055 and schema section 6.2 require ActivityEvent effects from their producing leaves. The focused real-process probe failed (terminated `exit_code: 1`): the journal validator rejects `activity_log/activity` as `unknown table`. The requirement-reference producer also needs a reference to its new `requirement_refs` row, while the current activity trigger resolves the V1 `requirements` table and the accepted replacement trigger is in task 15's migration 0018. Repairing that persistence contract before task 13 closes would change migration/schema and cross-task ownership; no task-13 scope or accepted decision authorizes that move. A bounded journal registry fix may be decision-free after the protected ordering is resolved, but it does not clear the trigger by itself. Incomplete producer implementation remains same-task work after resolution. The attached package times are below VAL-079 budgets, so timing is not a blocker.
- **Local decisions made:** none; no schema, migration order, ActivityEvent set, or validation criterion changed.
- **Resolution evidence:** Not cleared. The focused `go test -tags integration -run '^TestP23ReviewRecordRealProcessJourney$' -count=1 ./cmd/fid/integration` returned `FAIL` (terminated `exit_code: 1`); the full exact tagged task-13 gate also returned `FAIL` (terminated `exit_code: 1`). The accepted 0018 activity trigger references `requirement_refs`; the current trigger does not. The verified checkpoint and passing report-lint establish durable blocker evidence, not a passing VAL-055 result.
- **If REPAIR:** N/A pending protected authority. After resolution, task 13 continues under the authorized scope, preserves failed evidence, reruns the exact required gates, and obtains any affected independent review.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Do not change the journal table registry, migration or trigger schema, task order/scope, accepted Event effects or ActivityEvent names, public interface, tests, gates, thresholds, validation contract, or Runner ledger under this triage alone. Do not use a temporary V1 requirement row or invented Issue reference to bypass the accepted model.
- **Exact verification:** N/A pending protected authority.
- **Originating report:** project-docs/docs/phases/phase-23/tasks/task-13-done.md
- **Control repair base:** N/A; Runner supplied 007c9487863f5eb3cf037a0f0affd290742fed56, but no control-conformance directive is proposed.
- **If BLOCK:** category data-model/persistence-contract and protected cross-task migration ownership; ESC-23-007 raised in `tl_escalations_current.md`; resolution condition full Tech Lead obtains the required schema and task-order authority for the accepted ActivityEvent effects and requirement-reference integrity before task 13 resumes.
- **Human attention required:** yes
- **References:** project-docs/docs/phases/phase-23/tasks/task-13.md Objective, VAL-055, Files To Create Or Modify; project-docs/docs/phases/phase-23/tasks/task-13-done.md Validation Contract Coverage, Exact quality-gate results, Concerns; project-docs/docs/phases/phase-23/schema.md sections 4 and 6.2; project-docs/docs/phases/phase-23/contract-coverage.md public-interface 2 ActivityEvent row; project-docs/docs/phases/phase-23/qa/validation-contract.md VAL-055 and VAL-079; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9.

## TRG-23-008 - Task 13 protected-decision denial conflicts with closed failure schema

- **Trigger:** Task 13's VAL-057 protected-decision denial replay probe contradicts the accepted public failure schema.
- **Originating task:** task-13
- **Trigger signature:** task-13:protected-contract:protected-decision-denial-envelope-val-057
- **Attempt count:** 1
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-23/tasks/task-13-done.md
- **Candidate:** BLOCKED / protected public-interface and acceptance-contract contradiction
- **Checkpoint:** Runner-supplied verification: `git rev-parse HEAD^{commit}` returned 348a538167f94bd1d5a5f9a9ef798a90698565e7; subject matched `docs(task): finalize completion report for phase-23 task-13`; report-lint returned `OUTCOME BLOCKED` (terminated `exit_code: 0`).
- **Verdict:** BLOCK
- **Reasoning:** The accepted public-interface section 5 closes `protected-decision-required` details to `envelope_id`, `decision_class`, `required_capabilities`, and `affected_actions`; section 1 retains the closed `fid-cli.v1` error envelope. Section 5.2, DEC-226, and VAL-057 also demand immutable Stop/Attention IDs, null EventRef, and `replayed:false` in the denial response. The built binary returns the declared closed failure but omits those additional fields. Satisfying both readings needs protected authority over the public response contract or validation assertion semantics; this is not a verification-control implementation mismatch. The focused and full tagged gates failed (each terminated `exit_code: 1`). DEC-262 resolved the separate activity boundary and does not resolve this interface contradiction. DEC-227 timings remain within budget.
- **Local decisions made:** none; retain the failing assertion and original failed gate evidence.
- **Resolution evidence:** Not cleared. The focused `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23ReviewCyclesRealProcessJourney)$' -count=1 ./cmd/fid/integration` returned `FAIL` (terminated `exit_code: 1`), and the full exact tagged task-13 gate returned `FAIL` (terminated `exit_code: 1`). The accepted sources do not identify an unambiguous closed response shape satisfying both requirements.
- **If REPAIR:** N/A pending protected authority. After resolution, task 13 continues under the recorded contract, preserves the red evidence, and reruns the byte-identical required commands and any affected independent review.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Do not alter the public failure envelope, closed details, VAL-057 or VAL-056 assertion semantics, tests, gates, task-15 ownership, or Runner ledger under this triage alone.
- **Exact verification:** N/A
- **Originating report:** N/A
- **Control repair base:** N/A
- **If BLOCK:** category public-interface/specification contradiction; ESC-23-008 raised in `tl_escalations_current.md`; resolution condition full Tech Lead obtains the Human or accepted contract owner's ruling on the protected-decision denial response shape and aligns the public-interface contract, validation assertion, and task-13 verification before the Runner resumes.
- **Human attention required:** yes
- **References:** project-docs/docs/phases/phase-23/public-interface.md sections 1, 5, and 5.2; project-docs/docs/phases/phase-23/qa/validation-contract.md VAL-056 and VAL-057; project-docs/docs/phases/phase-23/decisions.md DEC-226 and DEC-262; project-docs/docs/phases/phase-23/tasks/task-13-done.md Workflow outcome, Validation Contract Coverage, Exact quality-gate results; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9.

## TRG-23-009 - Task 14 generic Concern list read path needs bounded convergence

- **Trigger:** Task 14's required tagged integration gate fails because the new Human projection test exercises the accepted generic concern.list --subject leaf, but cmd/fid/concern.go still routes through the V1 Issue-only reader.
- **Originating task:** task-14
- **Trigger signature:** task-14:scope-gap:generic-concern-list-reader
- **Attempt count:** 1
- **Repair category:** gate-repair
- **Repair authority:** project-docs/docs/phases/phase-23/tasks/task-14.md Objective and Binding Constraints; project-docs/docs/phases/phase-23/execution-scope.md section 1
- **Source:** project-docs/docs/phases/phase-23/tasks/task-14-done.md
- **Candidate:** RECOVERY_REQUIRED / bounded product reader scope gap
- **Checkpoint:** Runner-attached verification names implementation commit 59d2653da7cfe14699dc4d1cd85ef032bec26a5c and report commit eb6f4f110b30a72b4203d30605217b80613df776 with matching subjects; report-lint returned OUTCOME RECOVERY_REQUIRED (terminated exit_code: 0). Phase HEAD is eb6f4f110b30a72b4203d30605217b80613df776; only the Runner ledger is dirty.
- **Verdict:** REPAIR
- **Reasoning:** Task 14 owns every Phase 23 leaf's Human projection and introduced the tagged integration assertion for generic Concern list. The accepted public-interface section 2 already fixes its read flags and `{items:Concern[]}` result. The V1 Issue-only adapter is an undelivered accepted read path that this task's required gate now exposes; execution-scope section 1 and the standing DEC-175/194 convergence rules support minimum causal same-task repair. The public-interface row's earlier task 8/13 allocation does not move this failing task-14 gate to another task. No accepted definition or persistence contract needs to change.
- **Local decisions made:** Authorize only the generic Concern list reader and CLI routing needed for the existing task-14 every-leaf assertion; preserve task 15's separate V1 removal and table rename.
- **Resolution evidence:** Trigger remains real. The focused `TestP23ConcernListHumanProjection` command and the exact tagged integration gate each returned `FAIL` (terminated `exit_code: 1`): `concern list --subject` exits 2 through `issues.ListConcerns`. The ordinary race gate, format, lint, and coverage commands passed, but they do not clear this required tagged gate. The verified checkpoint and report-lint result establish durable recovery evidence, not task completion.
- **If REPAIR:** Owning task task-14; category gate-repair; dispatch same-task Developer continuation. Implement the accepted generic `--subject`, `--status`, `--severity`, and `--limit` read through `internal/review`, return the accepted Concern list projection, and retain the red Human/JSON integration assertion. Preserve the first failed results, rerun the focused red/green test and all five byte-identical required task commands, update the originating report with deviation and termination evidence, then repeat any affected independent window 13-15 re-review. Do not advance before resolution; no ESC.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A; this is an ordinary product task repair.
- **Control mismatch:** N/A; the reproduced product gate failure is recorded above.
- **Allowed files:** N/A for ordinary task-scope triage. Minimum causal recovery scope: `cmd/fid/concern.go`, `internal/review/concern.go` or a new `internal/review/concern_list.go` and matching in-package tests, `cmd/fid/integration/p23_concern_list_human_test.go`, and `project-docs/docs/phases/phase-23/tasks/task-14-done.md`.
- **Forbidden changes:** Do not change the accepted public-interface contract, validation assertions, exact commands, task ownership, persistence schema or journal contracts. Do not weaken or remove the retained red test or gate; do not perform task 15's V1 reader removal or table rename.
- **Exact verification:** N/A for ordinary task-scope triage; rerun the five exact commands in task-14.md from `/home/mark/projects/fid` (or the Runner-recorded same-task worktree root) and the focused test recorded in task-14-done.md.
- **Originating report:** project-docs/docs/phases/phase-23/tasks/task-14-done.md
- **Control repair base:** N/A
- **If BLOCK:** N/A; no protected category or exhausted no-progress budget is established.
- **Human attention required:** no
- **References:** project-docs/docs/phases/phase-23/tasks/task-14.md Objective, Binding Constraints, Scope, Exact Commands To Run; project-docs/docs/phases/phase-23/tasks/task-14-done.md Workflow outcome, Tests run, Exact quality-gate results, Deviations; project-docs/docs/phases/phase-23/public-interface.md section 2 concern.list row; project-docs/docs/phases/phase-23/execution-scope.md section 1; project-docs/docs/phases/phase-23/v1-convergence-inventory.md sections 1-2; project-docs/docs/current/decisions_active.md DEC-175/194; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9.

## TRG-23-010 - Task 15 migration 0018 down contradicts the accepted rollback chain

- **Trigger:** Task 15 cannot author migration 0018 down without choosing between two contradictory accepted schema authorities.
- **Originating task:** task-15
- **Trigger signature:** task-15:protected-contract:0018-down-concerns-column-duplication
- **Attempt count:** 1
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-23/tasks/task-15-done.md
- **Candidate:** BLOCKED / protected schema and persistence-contract contradiction
- **Checkpoint:** Runner-attached verification identifies checkpoint commit fa6355a0e221d97d7785b1e2718e134bb5609509 and report commit b1f00284de12be4c9c3403b8ad9433fa9f6ff4c8 with matching subjects; report-lint returned OUTCOME BLOCKED and PASS report-lint completed with no failures (terminated exit_code: 0). Runner reports phase HEAD b1f00284de12be4c9c3403b8ad9433fa9f6ff4c8.
- **Verdict:** BLOCK
- **Reasoning:** Accepted schema.md section 5 requires the 0018 down concerns definition to be byte-identical to migration 0002, which includes escalated_to_esc_id. The next rollback through migration 0016 down adds escalated_to_esc_id to concerns, so the literal section 5 definition makes that rollback fail with a duplicate column. Omitting the column revises the accepted byte-identity requirement. No accepted authority uniquely chooses a repair; this is a protected schema/persistence-contract contradiction, not a workflow-control defect.
- **Local decisions made:** none; no schema definition or migration order was selected.
- **Resolution evidence:** Not cleared. The focused `go test -run '^TestMigration0018DownRestoresEmptyV1Tables$' -count=1 ./internal/db` returned FAIL (terminated exit_code: 1) because 0018 down is absent. The test failure alone is ordinary incomplete work; the incompatible accepted definitions require a protected decision. The nine required task commands remain BLOCKED/unrun, with no green gate inferred from report-lint.
- **If REPAIR:** N/A pending protected schema authority. After resolution, resume task-15 from the preserved partial tree and verify the full down chain and all nine unchanged required commands; preserve failed evidence and obtain affected independent re-review.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A for ordinary task-scope triage.
- **Forbidden changes:** Do not silently choose or alter schema authority, migration ordering, accepted SQL byte identity, task ownership, validation assertions, exact commands, or the Runner ledger; do not claim blocked commands passed.
- **Exact verification:** N/A for ordinary task-scope triage.
- **Originating report:** project-docs/docs/phases/phase-23/tasks/task-15-done.md
- **Control repair base:** N/A
- **If BLOCK:** category protected schema/persistence-contract and specification contradiction; ESC-23-009 raised in `tl_escalations_current.md`; resolution condition the Tech Lead or delegated Human decision authority specifies the exact 0018-down concerns definition and keeps migration 0016 down and the full rollback chain executable, then authorizes same-task continuation.
- **Human attention required:** yes
- **References:** project-docs/docs/phases/phase-23/schema.md section 5; migrations/0002_v01_tables.up.sql:125-138; migrations/0016_p22_remove_escalations.down.sql:12; project-docs/docs/phases/phase-23/tasks/task-15.md Objective and Exact Commands To Run; project-docs/docs/phases/phase-23/tasks/task-15-done.md Workflow outcome, Tests run, Exact quality-gate results, Deviations; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9.
