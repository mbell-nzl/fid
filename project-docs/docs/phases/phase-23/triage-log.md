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
