# Phase 21 - TL-Lite Triage Log

> Append-only allocation source for `TRG-21-NNN` and `ESC-21-NNN`.
> Read the highest recorded number and add one; start at `001`.
> Entries use Recovery directive schema v2 from
> `project-docs/templates/triage-log-template.md`.

No Phase 21 triage event has been allocated. DEC-159 records Ready QA and Human
execution approval; task 1 is ready for Phase Runner dispatch.

---

## TRG-21-001 - Role-free bootstrap conflicts with retained V1 mutation contract

- **Trigger:** Phase 21 task 2 required quality gates fail after the accepted role-free Event-v2 bootstrap, exposing a specification/contract contradiction with retained V1 role-authorized mutations.
- **Originating task:** task-2
- **Trigger signature:** `task-2:specification-contradiction:task-2-done.md:role-free-bootstrap-retained-v1`
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** `project-docs/docs/phases/phase-21/tasks/task-2-done.md`
- **Candidate:** blocker candidate
- **Attempt count:** 1
- **Checkpoint:** task implementation `94e752d25809a0e54f827b1872172f8038c0a741` and final report `4efabd292d366c66a6680a9da712813c47b8c315`; prompt-attached commit verification and two-argument report-lint evidence
- **Verdict:** BLOCK
- **Reasoning:** DEC-157 requires `project.initialized` to be Event v2 with zero Effects and to never write `roles`; task 2 implements that accepted fixture. Retained V1 mutation tests still require the former seeded role identities after fresh init, and the required focused and repository coverage gates fail for that exact reason. Restoring roles, changing V1 authorization behavior, changing the bootstrap fixture, or altering the retained command/test contract would change protected specification, schema/event, or public-behavior authority.
- **Local decisions made:** none
- **Resolution evidence:** The attached byte-identical focused gate and repository coverage gate both terminated `exit_code: 1`; the report-lint terminated `exit_code: 0` with literal `OUTCOME BLOCKED`. No accepted decision selects between the role-free bootstrap and the retained V1 role-authorized mutation contract.
- **If REPAIR:** N/A
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Do not authorize product, test, migration, spec, task, QA, decision, report, or workflow-control changes under this triage.
- **Exact verification:** `go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...` in `/home/mark/projects/fid`; observed `FAIL` (terminated `exit_code: 1`).
- **Originating report:** `project-docs/docs/phases/phase-21/tasks/task-2-done.md`
- **Control repair base:** `4efabd292d366c66a6680a9da712813c47b8c315`
- **If BLOCK:** category specification contradiction; ESC-21-001 raised in `tl_escalations_current.md`; resolution condition a full Tech Lead/Human decision reconciles the accepted Event-v2 role-free bootstrap fixture with the retained V1 role-authorized mutation and gate contract.
- **Human attention required:** yes
- **References:** `project-docs/docs/phases/phase-21/decisions.md` DEC-157; `project-docs/docs/phases/phase-21/contracts.md` section 14; `project-docs/docs/phases/phase-21/v1-convergence-inventory.md` sections 1, 3, and 5; `project-docs/specs/09_AI_WORKFLOW.md` sections 2.3 and 9; `project-docs/docs/phases/phase-21/tasks/task-2-done.md`

---

## TRG-21-002 - Workflow-sequence external TRIAGE binding cannot satisfy task-prefix identity

- **Trigger:** the Phase Runner cannot append the already-written schema-v2 ledger TRIAGE for TRG-21-001 because the canonical workflow-sequence-lint verifier has mutually incompatible checks for the external trigger field's Markdown formatting and task-prefix identity.
- **Originating task:** task-2
- **Trigger signature:** phase-runner:control-conformance:workflow-sequence-lint:external-triage-trigger-prefix
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-21/triage-log.md and project-docs/tools/workflow-sequence-lint.sh
- **Candidate:** blocker candidate
- **Attempt count:** 1
- **Checkpoint:** phase HEAD 4efabd292d366c66a6680a9da712813c47b8c315 verified by the prompt-attached git rev-parse --verify result; task implementation 94e752d25809a0e54f827b1872172f8038c0a741 and report commit 4efabd292d366c66a6680a9da712813c47b8c315 were previously verified. Both supplied ledger candidates terminated exit_code: 1 and left the ledger byte-for-byte unchanged.
- **Verdict:** BLOCK
- **Reasoning:** The accepted schema-v2 TRIAGE contract requires both exact external-record binding and a task-prefixed trigger signature. TRG-21-001's rendered external Trigger signature contains Markdown backticks, so byte-exact copying presents a leading backtick to the task-prefix check, while an unquoted ledger value fails exact binding. Resolving the contradiction requires changing the canonical verifier's accepted parsing/control definition or the schema-v2 external-record contract; ordinary task repair cannot do either, and workflow-sequence-lint.sh is expressly excluded from normal control-conformance self-repair.
- **Local decisions made:** none
- **Resolution evidence:** The supplied unquoted candidate fails INVALID schema v2 TRIAGE identity does not match its unique external record, and the byte-for-byte copied candidate fails INVALID schema v2 TRIAGE Originating task does not match Trigger signature; each terminated exit_code: 1. No attached accepted decision authorizes a workflow-control definition change or a self-repair of project-docs/tools/workflow-sequence-lint.sh.
- **If REPAIR:** N/A
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** project-docs/tools/workflow-sequence-lint.sh; all prior run-ledger entries; all prior triage-log and escalation entries; no product, task, report, QA, spec, decision, or generated-file changes; no command, assertion, status, semantics, scope, threshold, budget, tolerance, or policy changes.
- **Exact verification:** project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-21/run-ledger.md --append-stdin with the candidate TRIAGE Markdown in /home/mark/projects/fid; the two reproduced candidates both terminated exit_code: 1.
- **Originating report:** project-docs/docs/phases/phase-21/tasks/task-2-done.md
- **Control repair base:** 4efabd292d366c66a6680a9da712813c47b8c315
- **If BLOCK:** category workflow policy / protected verifier control; ESC-21-002 raised in tl_escalations_current.md; resolution condition a full Tech Lead/Human decision must reconcile the canonical schema-v2 external Trigger signature representation with workflow-sequence-lint.sh's exact-binding and task-prefix identity checks, then record the authorized verifier or contract change and recovery route.
- **Human attention required:** yes
- **References:** project-docs/specs/09_AI_WORKFLOW.md sections 2.3, 2.4, and 9; project-docs/templates/triage-log-template.md; project-docs/docs/phases/phase-21/triage-log.md TRG-21-001; project-docs/tools/workflow-sequence-lint.sh; project-docs/docs/phases/phase-21/tasks/task-2-done.md.

---

## TRG-21-003 - Corrected external TRIAGE binding remains blocked by protected workflow verifier

- **Trigger:** Control-conformance triage continuation. The previous control triage record TRG-21-002 is itself not mechanically bindable because it uses Originating task task-2 with Trigger signature phase-runner:control-conformance:workflow-sequence-lint:external-triage-trigger-prefix; the schema-v2 verifier requires the external Trigger signature to begin with the exact Originating task. A ledger TRIAGE candidate copied byte-for-byte therefore failed `INVALID schema v2 TRIAGE Originating task does not match Trigger signature`, terminated exit_code: 1, leaving the ledger unchanged.
- **Originating task:** task-2
- **Trigger signature:** task-2:control-conformance:workflow-sequence-lint:external-triage-trigger-prefix
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-21/triage-log.md and project-docs/tools/workflow-sequence-lint.sh
- **Candidate:** blocker candidate
- **Attempt count:** 2
- **Checkpoint:** Control repair base 4efabd292d366c66a6680a9da712813c47b8c315; supplied candidate terminated exit_code: 1 and left the run ledger unchanged.
- **Verdict:** BLOCK
- **Reasoning:** The corrected task-prefixed external trigger can be mechanically bound only by changing or reinterpreting the canonical workflow-sequence verifier/control contract. The verifier is protected from normal control-conformance self-repair, and no accepted decision authorizes a protected workflow-control definition or policy change.
- **Local decisions made:** none
- **Resolution evidence:** The byte-for-byte ledger candidate for TRG-21-002 terminated exit_code: 1 with `INVALID schema v2 TRIAGE Originating task does not match Trigger signature` and left the ledger unchanged. The failure demonstrates a malformed prior external control record; it does not authorize rewriting the prior record, ledger, or protected verifier.
- **If REPAIR:** N/A
- **Control trigger signature:** task-2:control-conformance:workflow-sequence-lint:external-triage-trigger-prefix
- **Control attempt count:** 2
- **Control contract:** schema-v2 TRIAGE must byte-bind Originating task, Trigger signature, Attempt count, Repair category, and Repair authority to the external TRG while accepting a task-prefixed bound trigger signature.
- **Control mismatch:** TRG-21-002's phase-runner-prefixed trigger cannot satisfy the verifier's task-prefix identity check for Originating task task-2; the byte-for-byte ledger candidate terminated exit_code: 1 and left the ledger unchanged.
- **Allowed files:** N/A - workflow-sequence-lint.sh is protected from normal control-conformance self-repair.
- **Forbidden changes:** project-docs/tools/workflow-sequence-lint.sh; all prior run-ledger entries; all prior triage-log and escalation entries; no product, task, report, QA, spec, decision, or generated-file changes.
- **Exact verification:** project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-21/run-ledger.md --append-stdin with the TRIAGE Markdown candidate in /home/mark/projects/fid; the candidate terminated exit_code: 1.
- **Originating report:** project-docs/docs/phases/phase-21/tasks/task-2-done.md
- **Control repair base:** 4efabd292d366c66a6680a9da712813c47b8c315
- **If BLOCK:** category workflow policy / protected verifier control; ESC-21-003 raised in `tl_escalations_current.md`; resolution condition a full Tech Lead/Human decision must authorize the accepted canonical repair or replacement route for the malformed TRG-21-002 binding without rewriting historical evidence.
- **Human attention required:** yes
- **References:** project-docs/specs/09_AI_WORKFLOW.md sections 2.3, 2.4, and 9; project-docs/templates/triage-log-template.md; project-docs/docs/phases/phase-21/triage-log.md TRG-21-002; project-docs/tools/workflow-sequence-lint.sh; project-docs/docs/phases/phase-21/tasks/task-2-done.md.

## Resolution Note For TRG-21-001 Through TRG-21-003

- **Decision authority:** DEC-160 and DEC-161, 2026-08-23.
- **TRG-21-001:** sustained; DEC-160 preserves the role-free Event-v2
  bootstrap and authorizes the exact temporary V1 policy vocabulary bridge as
  bounded task-2 recovery. ESC-21-001 is resolved.
- **TRG-21-002:** retained as malformed historical evidence. Its Markdown
  code-span value and non-task-prefixed correction are not governing recovery
  records. ESC-21-002 is resolved by DEC-161 without rewriting the entry.
- **TRG-21-003:** sustained and mechanically valid. DEC-161 selects it as the
  sole governing STOP record and authorizes RUN-RESUME attempt 3 into task-2
  task-scope recovery. ESC-21-003 is resolved.
- **Verifier disposition:** no verifier or control-conformance repair is
  authorized or required. Future schema-v2 identity fields are plain,
  task-prefixed scalars under `09_AI_WORKFLOW.md`.
- **Historical integrity:** prior TRG, failed append, STOP, and ledger evidence
  remain immutable. The Phase Runner appends the DEC-161 RUN-RESUME; the Tech
  Lead does not edit `run-ledger.md`.

---

## TRG-21-004 - Retained rebuild and role-row tests contradict accepted Event-v2 contracts

- **Trigger:** after DEC-160's authorized temporary V1 policy vocabulary bridge and focused policy tests passed, the exact byte-identical required package and repository coverage gates still fail because retained tests require the prohibited `fid rebuild --by` flag and persisted `roles` rows.
- **Originating task:** task-2
- **Trigger signature:** task-2:specification-contradiction:task-2-done.md:retained-v1-rebuild-role-contract
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-21/tasks/task-2-done.md
- **Candidate:** blocker candidate
- **Attempt count:** 5
- **Checkpoint:** implementation commit 861e90da8b9218915ddbd77515326921ff3c0cf4 and final report repair commit 0e231bcb1a86a1142dd67c3137703da53c394e3c; prompt-attached git rev-parse --verify, subject-match, and two-argument report-lint evidence.
- **Verdict:** BLOCK
- **Reasoning:** DEC-160 authorizes only a nonpersistent V1 policy vocabulary bridge and expressly preserves the exact role-free Event-v2 bootstrap. The retained tests require two protected contracts that the accepted authorities prohibit: a `fid rebuild --by` leaf flag despite the no-leaf-flags public interface, and persisted `roles` rows despite the empty-Effects, role-free projection contract. Resolving the gate failure would require a public-interface, data-model/persistence, test-contract, or specification decision; neither ordinary task repair nor control-conformance authority applies.
- **Local decisions made:** none
- **Resolution evidence:** The DEC-160 focused policy command passed without state writes. The exact required package command and byte-identical repository coverage command each printed `FAIL` and terminated `exit_code: 1`; the supplied report-lint result terminated `exit_code: 0` with literal `OUTCOME BLOCKED`. No test, contract, flag, role row, migration, or protected assertion was weakened.
- **If REPAIR:** N/A
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** no test weakening; no rebuild flag restoration; no bootstrap role Effects; no persisted compatibility rows; no migration 0014; no task, QA, spec, decision, report, workflow-control, command, assertion, status, semantics, scope, threshold, budget, tolerance, or policy changes under this triage.
- **Exact verification:** `go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...` and `go test -race -coverprofile=coverage.out ./...` in `/home/mark/projects/fid`; both observed `FAIL` (terminated `exit_code: 1`).
- **Originating report:** project-docs/docs/phases/phase-21/tasks/task-2-done.md
- **Control repair base:** 0e231bcb1a86a1142dd67c3137703da53c394e3c
- **If BLOCK:** category specification contradiction, public interface, and data model/persistence contract; ESC-21-004 raised in `tl_escalations_current.md`; resolution condition a full Tech Lead/Human decision reconciles the retained V1 rebuild flag and persisted role-row test contract with DEC-157/DEC-160's role-free Event-v2 bootstrap, no-leaf-flags rebuild interface, and role-free projection.
- **Human attention required:** yes
- **References:** project-docs/docs/phases/phase-21/decisions.md DEC-157 and DEC-160; project-docs/docs/phases/phase-21/public-interface.md:330; project-docs/docs/phases/phase-21/contracts.md:377-383; project-docs/docs/phases/phase-21/tasks/task-2-done.md; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9.

## Resolution Note For TRG-21-004

- **Decision authority:** DEC-162, 2026-08-23.
- **Verdict audit:** the BLOCK was correct on the authority available. The
  Developer had no scope to touch the retained tests and refused to weaken them,
  which is the required behaviour. The recorded protected category was wrong in
  hindsight: the accepted contracts do not contradict each other. The residual
  failures are retained-test collateral of an accepted interface change that
  task 2's file scope omitted.
- **Resolution:** DEC-162 keeps the no-leaf-flags `project.rebuild` interface and
  the role-free Event-v2 bootstrap exactly as accepted, and rules that retained
  V1 tests converge inside the task whose accepted change breaks them. Task 2's
  scope now covers the affected test files, with ordered adaptation rules and an
  explicit no-weakening boundary in `v1-convergence-inventory.md` section 1.2.
- **Recovery:** one more same-task, task-scope `in-scope-implementation` attempt
  on task 2 under this trigger signature, attempt count 6. Required commands are
  unchanged and rerun byte-identically. ESC-21-004 is resolved.
- **Historical integrity:** no prior TRG, ledger, report, or failed-gate evidence
  is rewritten. The Phase Runner appends the RUN-RESUME; the Tech Lead does not
  edit `run-ledger.md`.

---

## TRG-21-005 - Multi-owner QA review cannot bind one schema-v2 recovery owner

- **Trigger:** The accepted QA Review 1 artifact has two bounded Must-fix findings owned by task-1 and task-3, but the schema-v2 run ledger REVIEW event has one Owning task field and the canonical workflow-sequence validator binds the active QA recovery to one exact owner. The accepted REVIEW event was recorded as Owning task: task-1, task-3. A valid task-1 recovery candidate was rejected, leaving the valid ledger state REPAIR_PENDING; the append-only CORRECTION was accepted but cannot alter validator state.
- **Originating task:** task-1
- **Trigger signature:** phase-runner:orchestration-integrity:review-1:multi-owner-ledger-recovery-binding
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-21/qa/review-1.md
- **Candidate:** blocker candidate
- **Attempt count:** 1
- **Checkpoint:** review commit 17d4e22d15b6fdc0024a8fd8ff539ae6f7ceeb5b verified by prompt-attached git rev-parse --verify result and subject match `docs(qa): review phase 21 tasks 1-3`; phase branch phase-21; runner-recorded control repair base 17d4e22d15b6fdc0024a8fd8ff539ae6f7ceeb5b; only the runner-owned ledger is dirty and hidden index flags are absent.
- **Verdict:** BLOCK
- **Reasoning:** Review 1 has two separately owned, bounded Must-fix findings, but the accepted schema-v2 REVIEW record declares one `Owning task` field and the canonical validator binds active QA recovery to that exact value. The supplied task-1 candidate is rejected because `task-1, task-3` is not one recovery task; the accepted CORRECTION preserves the task-3 finding but cannot alter the validator's active owner. Selecting sequential ownership, admitting a list, or teaching the validator to honor CORRECTION changes canonical workflow sequence semantics and policy.
- **Local decisions made:** none
- **Resolution evidence:** The supplied exact append control `workflow-sequence-lint.sh project-docs/docs/phases/phase-21/run-ledger.md --append-stdin` with a TASK-RECOVERY task-1 candidate required a valid task-1 QA-rework recovery and observed `INVALID QA repair owner task-1, task-3 does not match recovery task task-1`, terminated exit_code: 1. Standalone validation after the accepted CORRECTION remains `VALID state=REPAIR_PENDING`; no validator-valid route exists for either owner.
- **If REPAIR:** N/A
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** project-docs/tools/workflow-sequence-lint.sh; project-docs/docs/phases/phase-21/run-ledger.md; project-docs/docs/phases/phase-21/qa/review-1.md; task specs; product files; tests; decisions; specs; prompts; templates; dependencies; QA contracts; reports; generated files; and no command, assertion, status, semantics, scope, threshold, budget, tolerance, or policy changes.
- **Exact verification:** `workflow-sequence-lint.sh project-docs/docs/phases/phase-21/run-ledger.md --append-stdin` with a TASK-RECOVERY task-1 candidate in `/home/mark/projects/fid`; required a valid task-1 QA-rework recovery; observed `INVALID QA repair owner task-1, task-3 does not match recovery task task-1` (terminated exit_code: 1).
- **Originating report:** project-docs/docs/phases/phase-21/qa/review-1.md
- **Control repair base:** 17d4e22d15b6fdc0024a8fd8ff539ae6f7ceeb5b
- **If BLOCK:** category workflow policy / canonical verifier control; ESC-21-005 raised in `tl_escalations_current.md`; resolution condition a full Tech Lead/Human decision must select and record the accepted multi-owner QA-review recovery representation and canonical sequence behavior, including whether a corrected current owner may bind sequential same-review recoveries, before any validator or workflow-definition change is authorized.
- **Human attention required:** yes
- **References:** project-docs/docs/phases/phase-21/qa/review-1.md; project-docs/docs/phases/phase-21/run-ledger.md:386; project-docs/tools/workflow-sequence-lint.sh; project-docs/specs/09_AI_WORKFLOW.md sections 2.3, 2.4, and 9; project-docs/templates/triage-log-template.md.

---

## TRG-21-006 - Corrected external identity retains the multi-owner workflow-policy block

- **Trigger:** Triage correction trigger: The prior TL-Lite BLOCK for this same protected workflow ambiguity was written as TRG-21-005 with Trigger signature phase-runner:orchestration-integrity:review-1:multi-owner-ledger-recovery-binding. The schema-v2 ledger TRIAGE bind correctly rejected that external record because Originating task is task-1 and the trigger signature must carry the task-prefixed identity. The prior external BLOCK and ESC-21-005 are historical evidence and must not be rewritten or duplicated.
- **Originating task:** task-1
- **Trigger signature:** task-1:orchestration-integrity:review-1:multi-owner-ledger-recovery-binding
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** ESC-21-005; project-docs/docs/phases/phase-21/triage-log.md TRG-21-005
- **Candidate:** correction of external schema-v2 identity binding for an existing blocker
- **Attempt count:** 1
- **Checkpoint:** review commit 17d4e22d15b6fdc0024a8fd8ff539ae6f7ceeb5b verified by prompt-attached git rev-parse --verify result and subject match `docs(qa): review phase 21 tasks 1-3`; runner-recorded control repair base 17d4e22d15b6fdc0024a8fd8ff539ae6f7ceeb5b.
- **Verdict:** BLOCK
- **Reasoning:** The correction supplies the task-prefixed external identity required for a schema-v2 TRIAGE bind without rewriting TRG-21-005 or the existing escalation. It does not resolve the protected condition: Review 1 still has task-1 and task-3 Must-fix owners while the canonical REVIEW/recovery control admits one active recovery owner. Selecting sequential ownership, allowing multiple owners, or changing how CORRECTION affects recovery changes canonical workflow sequence semantics and policy.
- **Local decisions made:** none
- **Resolution evidence:** The supplied TASK-RECOVERY task-1 append candidate required a valid task-1 QA-rework recovery but observed `INVALID QA repair owner task-1, task-3 does not match recovery task task-1` (terminated exit_code: 1). Standalone validation after the accepted CORRECTION remains `VALID state=REPAIR_PENDING`. The prior TRIAGE candidate for TRG-21-005 was rejected because its external Trigger signature was not task-prefixed; this entry corrects that external identity only.
- **If REPAIR:** N/A
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** project-docs/tools/workflow-sequence-lint.sh; project-docs/docs/phases/phase-21/run-ledger.md; project-docs/docs/phases/phase-21/qa/review-1.md; task specs; product files; tests; decisions; specs; prompts; templates; dependencies; QA contracts; reports; generated files; and no command, assertion, status, semantics, scope, threshold, budget, tolerance, or policy changes.
- **Exact verification:** `workflow-sequence-lint.sh project-docs/docs/phases/phase-21/run-ledger.md --append-stdin` with a TASK-RECOVERY task-1 candidate in `/home/mark/projects/fid`; required a valid task-1 QA-rework recovery; observed `INVALID QA repair owner task-1, task-3 does not match recovery task task-1` (terminated exit_code: 1).
- **Originating report:** project-docs/docs/phases/phase-21/qa/review-1.md
- **Control repair base:** 17d4e22d15b6fdc0024a8fd8ff539ae6f7ceeb5b
- **If BLOCK:** category workflow policy / canonical verifier control; ESC-21-005 referenced in `tl_escalations_current.md`; resolution condition a full Tech Lead/Human decision must select and record the accepted multi-owner QA-review recovery representation and canonical sequence behavior, including whether a corrected current owner may bind sequential same-review recoveries, before any validator or workflow-definition change is authorized.
- **Human attention required:** yes
- **References:** project-docs/docs/current/tl_escalations_current.md ESC-21-005; project-docs/docs/phases/phase-21/triage-log.md TRG-21-005; project-docs/docs/phases/phase-21/qa/review-1.md; project-docs/docs/phases/phase-21/run-ledger.md:386; project-docs/tools/workflow-sequence-lint.sh; project-docs/specs/09_AI_WORKFLOW.md sections 2.3, 2.4, and 9; project-docs/templates/triage-log-template.md.

## Resolution Note For TRG-21-005 And TRG-21-006

- **Decision authority:** DEC-163, 2026-08-23.
- **TRG-21-005:** sustained. The protected condition was real; the entry is
  retained as historical evidence and its non-task-prefixed external identity is
  not a governing recovery record.
- **TRG-21-006:** sustained and mechanically valid. It is the sole governing
  STOP record. DEC-163 authorizes RUN-RESUME attempt 2 into a task-1 `qa-rework`
  task-scope recovery for review-1 MF-1.
- **Accepted representation:** a QA `REPAIR` review carries an ordered
  comma-separated owner list in one `Owning task` field; every listed owner
  completes TASK-RECOVERY, REPAIR, and TASK complete before a linked re-review
  may close or replace the review; `CORRECTION` remains evidence-only and never
  selects owners.
- **Verifier disposition:** unlike TRG-21-002/003, a verifier change *is*
  required and is authorized. `workflow-sequence-lint.sh` accepted a REVIEW
  append that no later event could recover, so it is corrected under DEC-163
  item 4 with new drills under item 5, outside the phase recovery loop. A
  control-conformance route is impossible because the open multi-owner review
  rejects the TASK-RECOVERY that would authorize it.
- **Historical integrity:** no prior TRG, ledger, review, escalation, or failed
  append evidence is rewritten. The 15:29 CORRECTION stands; its single-owner
  narrowing is superseded by DEC-163 item 1 and its named first owner `task-1`
  matches the accepted list order. The Phase Runner appends the RUN-RESUME; the
  Tech Lead does not edit `run-ledger.md`.

---

## TRG-21-007 - Task 8 retained assertions require accepted-contract convergence

- **Trigger:** Task 8 retained Issue convergence exposes an alleged protected contradiction between accepted Phase 21 Event-v2, kebab-case activity, and attribution-only standalone Actor contracts and retained test assertions; the child also returned a malformed BLOCKED report that report-lint rejected.
- **Originating task:** task-8
- **Trigger signature:** task-8:specification-contradiction:task-8-done.md:retained-v1-rebuild-role-contract
- **Repair category:** in-scope-implementation
- **Repair authority:** task-8 specification and DEC-162
- **Source:** project-docs/docs/phases/phase-21/tasks/task-8-done.md
- **Candidate:** malformed `OUTCOME BLOCKED` report and retained-test gate failure
- **Attempt count:** 1
- **Checkpoint:** implementation commit e712f85d08a0b3259b23badd34d9306bf0a10f9b and report checkpoint 44370af13c1c5c39e9e040065bc70e15e30b8deb; prompt-attached git rev-parse --verify results and subject matches passed. The Phase Runner report-lint observed `OUTCOME BLOCKED` with exit_code 2, including missing `next_action` and invalid category `protected-contract-conflict`.
- **Verdict:** REPAIR
- **Reasoning:** DEC-162 establishes that retained V1 tests and fixtures are not protected contracts, and its standing rule explicitly applies to task 8's `--by` to `--actor` convergence. The Phase 21 public interface and inventory uniquely require typed standalone Actor attribution, keyed Event-v2 Issue mutations, and kebab-case accepted activity values; the retained assertions must converge while preserving their subjects and strength. The failed required gates and malformed report keep task 8 open but do not require a protected contract change.
- **Local decisions made:** retained assertions for Issue Event fields, activity values, and standalone Actor behavior must be migrated to the already accepted Phase 21 contract under task 8; the child-authored `BLOCKED` outcome is not authoritative lifecycle state and must be replaced by a valid recovery/completion report that preserves the failed evidence.
- **Resolution evidence:** `project-docs/docs/phases/phase-21/public-interface.md` requires retained standalone mutations to use `--actor <type>:<id>` as attribution only and never kit authority, and specifies Event-producing mutation idempotency plus kebab-case ActivityEvent values. `v1-convergence-inventory.md` assigns Issue create/update/transition Event-v2 conversion and retained activity normalization to task 8 and applies DEC-162's no-weakening retained-test convergence rule expressly to task 8. The prompt-attached focused/package, lint, and repository failures are recoverable task-scope evidence; no reproduced canonical-control mismatch exists.
- **If REPAIR:** owning task task-8; category in-scope-implementation; next action dispatch a same-task continuation to converge the failing retained Issue/Event/activity/Actor assertions and lint findings in task-8 scope, preserve failed attempts and assertion strength, rerun every unchanged required command, then write and commit a valid report with the Phase Runner's unchanged two-argument report-lint path; no ESC.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** no reintroduction of `--by`, V1 Event fields, snake-case accepted activity values, standalone role authority, public-interface changes, Event-schema or persistence changes, task/QA/spec/decision/run-ledger changes, required-command or assertion-semantics changes, test weakening, `t.Skip`, build tags, `-run` narrowing, threshold/budget/tolerance/policy changes, or held-out scenario access.
- **Exact verification:** `go test -race ./internal/issues ./internal/activity ./internal/journal ./internal/cli ./cmd/fid/...`; `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` in `/home/mark/projects/fid`; after the report commit, `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-8-done.md project-docs/docs/phases/phase-21/tasks/task-8-issue-actor-convergence.md` in `/home/mark/projects/fid`.
- **Originating report:** project-docs/docs/phases/phase-21/tasks/task-8-done.md
- **Control repair base:** 44370af13c1c5c39e9e040065bc70e15e30b8deb
- **If BLOCK:** N/A
- **Human attention required:** no
- **References:** project-docs/docs/phases/phase-21/decisions.md DEC-162; project-docs/docs/current/decisions_active.md DEC-162; project-docs/docs/phases/phase-21/v1-convergence-inventory.md sections 1, 1.2, and 3; project-docs/docs/phases/phase-21/public-interface.md lines 88-98 and 331-350; project-docs/docs/phases/phase-21/tasks/task-8-issue-actor-convergence.md lines 19-28 and 40-72; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9.

---

## TRG-21-008 - Schema-v2 ordinary repair authority cannot bind TASK-RECOVERY

- **Trigger:** The accepted ordinary REPAIR authority in TRG-21-007 cannot bind to a schema-v2 TASK-RECOVERY because workflow-sequence-lint.sh requires the byte-identical Repair authority value to also reference TRG-21-007, while the external TRG's exact Repair authority is 'task-8 specification and DEC-162'. The exact TASK-RECOVERY candidate is rejected; adding the TRG ID would violate external TRG byte binding.
- **Originating task:** task-8
- **Trigger signature:** task-8:orchestration-integrity:workflow-sequence-lint:ordinary-triage-authority-binding
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** TRG-21-007 / project-docs/docs/phases/phase-21/run-ledger.md append candidate
- **Candidate:** schema-v2 TASK-RECOVERY task-8 append candidate
- **Attempt count:** 1
- **Checkpoint:** runner-recorded control repair base 44370af13c1c5c39e9e040065bc70e15e30b8deb; prompt-attached commit verification resolves it as `docs(task): checkpoint phase-21 task-8 conflict`, with full-hash verification and subject match passed.
- **Verdict:** BLOCK
- **Reasoning:** The accepted schema-v2 contract requires ordinary TRIAGE REPAIR authority to bind byte-for-byte from the external TRG into TASK-RECOVERY. The canonical verifier additionally requires the bound TASK-RECOVERY Repair authority to reference TRG-21-007, but the binding authority is exactly `task-8 specification and DEC-162` and contains no TRG ID. Changing that verifier requirement or admitting an authority representation that differs from the external TRG changes protected canonical workflow-control definition and semantics. `workflow-sequence-lint.sh` is expressly excluded from normal control-conformance self-repair.
- **Local decisions made:** none
- **Resolution evidence:** The attached exact append reproducer required a valid bound schema-v2 TASK-RECOVERY using task-8, the governing TRG-21-007 trigger and attempt, category `in-scope-implementation`, mode `task-scope`, and exact external Repair authority. It observed `INVALID TASK-RECOVERY does not preserve its governing schema v2 TRIAGE REPAIR authority; INVALID overlapping TASK-RECOVERY precedes required REPAIR` (terminated exit_code: 1). Prefixing the authority with `TRG-21-007 / ` would violate the accepted byte-identical external TRIAGE binding.
- **If REPAIR:** N/A
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** project-docs/tools/workflow-sequence-lint.sh; project-docs/tools/workflow-control-scope-lint.sh; project-docs/docs/phases/phase-21/run-ledger.md; triage-log history; task, spec, QA, decision, public-interface, product, test, dependency, report, prompt, template, and generated paths; and no command, assertion, status, semantics, scope, threshold, budget, tolerance, or policy changes.
- **Exact verification:** `project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-21/run-ledger.md --append-stdin` in `/home/mark/projects/fid` using a byte-identical TASK-RECOVERY candidate copied from TRG-21-007.
- **Originating report:** project-docs/docs/phases/phase-21/tasks/task-8-done.md
- **Control repair base:** 44370af13c1c5c39e9e040065bc70e15e30b8deb
- **If BLOCK:** category workflow policy / protected verifier definition; ESC-21-006 raised in `tl_escalations_current.md`; resolution condition a full Tech Lead/Human decision reconciles the schema-v2 ordinary TRIAGE-to-TASK-RECOVERY byte-binding contract with the canonical verifier's additional TRG-reference requirement, then authorizes any verifier or contract-definition change and recovery route.
- **Human attention required:** yes
- **References:** TRG-21-007; project-docs/tools/workflow-sequence-lint.sh lines 619-637; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9; project-docs/templates/triage-log-template.md; project-docs/docs/phases/phase-21/decisions.md DEC-162.

---

## TRG-21-009 - Task 9 gate closure requires Task 10 Claim/Run fixture authority

- **Trigger:** Task 9's exact required package and coverage gates remain red on retained cmd/fid Event-v2/result/activity and AC/concern no-writer fixtures plus explicit replacement-bound Task 10 Claim/Run authority fixtures. The task-9 report says the Task 10 fixtures remain outside task-9 scope.
- **Originating task:** task-9
- **Trigger signature:** task-9:scope-boundary:task-9-done.md:replacement-bound-task10-claim-run-fixtures
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-21/tasks/task-9-done.md
- **Candidate:** `OUTCOME RECOVERY_REQUIRED` report and retained required-gate failures
- **Attempt count:** 8
- **Checkpoint:** report commit 8eabb7c4e3ec54352bed9cf04fb7224280b9415a; prompt-attached git rev-parse --verify returned the same full commit and subject matched `docs(task): finalize recovery report for phase-21 task-9`; two-argument report-lint returned `OUTCOME RECOVERY_REQUIRED` (terminated `exit_code: 0`).
- **Verdict:** BLOCK
- **Reasoning:** DEC-162 authorizes retained-test convergence only inside the task whose accepted change invalidates the test, preserving subject and assertion strength. Task 9 owns retained attribution and judgment-freeze convergence, but expressly excludes generic Claims/Runs; the inventory and task 10 specification assign Claim/Run replacement authority and its contract to task 10. Moving explicit Task 10 authority fixtures into task 9, changing their assertion subject, or advancing around their red gates requires a protected task-boundary re-slice or cross-task authority decision. The prompt supplies no reproduced canonical-control mismatch, so control-conformance repair is unavailable.
- **Local decisions made:** none
- **Resolution evidence:** The task-9 report records `RECOVERY_REQUIRED` and identifies the remaining explicit Task 10/Claim/Run authority fixtures as outside task-9 scope. Its exact required package command and repository coverage command are recorded `BLOCKED` with observed `exit_code: 1`; report-lint confirms recovery remains required. Task 9's Out Of Scope section excludes Generic Claims/Runs and final V1 removal, while task 10 owns generic Claim operations and the V1 inventory assigns Claim replacement to task 10 and Run replacement to task 11.
- **If REPAIR:** N/A
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** no task-10 Claim/Run implementation, no final V1 removal, no widening of task-9 scope, no command/assertion/status/threshold/policy changes, no test weakening, no Task 10 or later product/spec/task/QA/decision edits without their authority.
- **Exact verification:** the task-9 exact package command and repository coverage command are recorded in task-9-done.md as BLOCKED with observed exit_code: 1; the report identifies remaining explicit replacement-bound Task 10 Claim/Run authority fixtures; report-lint project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-9-done.md project-docs/docs/phases/phase-21/tasks/task-9-retained-actor-policy-removal.md returned OUTCOME RECOVERY_REQUIRED with observed exit_code: 0 in /home/mark/projects/fid.
- **Originating report:** project-docs/docs/phases/phase-21/tasks/task-9-done.md
- **Control repair base:** 8eabb7c4e3ec54352bed9cf04fb7224280b9415a
- **If BLOCK:** category protected task-boundary / cross-task authority; ESC-21-007 raised in `tl_escalations_current.md`; resolution condition a full Tech Lead decision selects a bounded ownership and sequencing route for the named Task 10 Claim/Run fixtures while retaining all accepted public authority contracts, task gates, and assertion strength.
- **Human attention required:** yes
- **References:** project-docs/docs/phases/phase-21/tasks/task-9-done.md; project-docs/docs/phases/phase-21/tasks/task-9-retained-actor-policy-removal.md Scope and Out Of Scope; project-docs/docs/phases/phase-21/tasks/task-10-generic-claims.md Scope and Acceptance Criteria; project-docs/docs/phases/phase-21/v1-convergence-inventory.md sections 1, 1.2, and 5; project-docs/docs/phases/phase-21/decisions.md DEC-162; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9.

## Resolution Note For TRG-21-009

- **Decision authority:** DEC-165, 2026-08-25.
- **Verdict audit:** **not sustained on its recorded category.** The refusal to
  widen scope or weaken the Claim assertions on child authority was correct
  behaviour. The classification was wrong: the named fixtures are not Task 10
  Claim/Run authority fixtures. No generic Claim or Run fixture exists in the
  tree, task 10 has not been dispatched, and every red assertion is an incidental
  read of a surface task 9 itself removed or renamed - public `fid policy show`
  (`cmd/fid/issue_claim_test.go:464`, `cmd/fid/run_test.go:184`,
  `cmd/fid/dev_test.go:134`, `cmd/fid/label_test.go:336`,
  `cmd/fid/link_test.go:413`, `cmd/fid/issue_label_test.go:517`,
  `cmd/fid/requirement_test.go:306`,
  `cmd/fid/integration/e2e_v01_closure_helpers_test.go:329`), stats `--by role`
  (`:327-328`), and fail-closed `maintenance vacuum --by tl` (`:316`). The V1
  claim, run, and dev behavioural assertions in those same files still pass.
- **Resolution:** DEC-165 states DEC-162's rule explicitly - convergence
  ownership follows the invalidating accepted change, not the fixture file's
  primary subject - and makes the bounded task-9 amendment: `cmd/fid/integration/`
  is named in file scope, Out Of Scope is narrowed to generic Claim/Run
  implementation and semantics, and the incidental `policy show` reads are
  repointed onto the retained unregistered legacy matrix helper with every
  role/action assertion byte-identical. No task-10/11 work, no re-slicing, no
  sequencing change, no gate/VAL/command/contract change.
- **Recovery:** RUN-RESUME from the STOP, then one same-task `TASK-RECOVERY
  task-9`, category `in-scope-implementation`, mode `task-scope`, trigger
  signature `task-9:scope-boundary:task-9-done.md:replacement-bound-task10-claim-run-fixtures`,
  attempt 9, `Source` and `Repair authority` `TRG-21-009 / DEC-165`,
  no-progress `0/2`. All five required commands rerun byte-identically.
  ESC-21-007 is resolved.
- **Recurrence note:** this is the second occurrence of the TRG-21-004 pattern -
  residual failures that are collateral of an accepted change the escalating
  task itself made, misclassified as protected by file subject. `08_QUALITY_PLAN.md`
  §3 now carries the generalised ownership rule.
- **Historical integrity:** no prior TRG, ledger, report, escalation, or
  failed-gate evidence is rewritten. The Phase Runner appends the RUN-RESUME; the
  Tech Lead does not edit `run-ledger.md`.

---

## TRG-21-010 - Task 9 convergence now conflicts with retained Task 10 contract assertions

- **Trigger:** Task 9's unchanged required gates remain red on two assertions in cmd/fid/task10_contract_test.go that exercise Task 9-converted inventory leaves and still require legacy --by behavior and old authority-missing envelopes.
- **Originating task:** task-9
- **Trigger signature:** task-9:scope-boundary:task-9-done.md:task9-converted-leaf-actor-contract-retained-legacy-by-assertions
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-21/tasks/task-9-done.md
- **Candidate:** protected-blocker candidate after Task 9 recovery attempt 10
- **Attempt count:** 11
- **Checkpoint:** prompt-attached commit verification identifies 77575e419dd20e341fa04a22344063187281db41 with subject `docs(task): finalize completion report for phase-21 task-9`; two-argument report-lint returned `OUTCOME BLOCKED` and `PASS report-lint completed with no failures` (terminated `exit_code: 0`). The exact required package command and repository coverage command remain FAIL (terminated `exit_code: 1`).
- **Verdict:** BLOCK
- **Reasoning:** DEC-165 binds retained-test convergence to the invalidating accepted Task 9 change rather than a fixture filename, but it also preserves the Task 9 `--actor` plus idempotency interface and expressly forbids restoring `--by` on a Task 9 family. The retained assertions now require both the prohibited legacy flag and old authority-missing envelopes for Task 9-converted leaves. Repairing them requires selecting a public CLI/authority contract or changing retained assertion semantics, neither of which Task 9 can decide. This is the third recurrence of the TRG-21-004 file-subject misclassification pattern: TRG-21-009's Resolution Note records itself as the second; the supplied recurrence rule requires escalation rather than another local recovery.
- **Local decisions made:** none
- **Resolution evidence:** the reproduced command `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` required PASS but observed FAIL (terminated `exit_code: 1`) in `TestTask10JSONMissingActorsUseAuthorityEnvelope` and `TestTask10RequiredActorHumanErrorsRemainProse`. Those assertions retain legacy `--by` behavior for Task 9-converted `issue ac add`, `label add`, and `issue label add/remove` leaves, while Task 9's accepted interface is `--actor` plus idempotency. No canonical workflow-control mismatch is reproduced.
- **If REPAIR:** N/A
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A - no control-conformance candidate
- **Control mismatch:** N/A - no canonical control mismatch reproduced
- **Allowed files:** N/A
- **Forbidden changes:** do not restore --by on a task-9 family; do not modify cmd/fid/task10_contract_test.go or Task 10 implementation/semantics; do not weaken or narrow the exact gates; no task-spec, QA artifact, validation-contract, public-interface, Event-schema, required-command, workflow-verifier, t.Skip, build-tag, -run, threshold, budget, tolerance, or policy change.
- **Exact verification:** N/A - no control-conformance candidate
- **Originating report:** project-docs/docs/phases/phase-21/tasks/task-9-done.md
- **Control repair base:** 77575e419dd20e341fa04a22344063187281db41
- **If BLOCK:** category public-interface / specification contradiction and repeated post-DEC-165 task-boundary classification; ESC-21-008 raised in `tl_escalations_current.md`; resolution condition a full Tech Lead must determine the compatible accepted contract and ownership route for the retained Task 10-named assertions without restoring a forbidden Task 9 legacy interface or weakening assertion semantics.
- **Human attention required:** yes
- **References:** project-docs/docs/phases/phase-21/decisions.md DEC-165 items 1, 2, 4, 5, and 6; project-docs/docs/current/phase_brief.md DEC-165 and Dispatch Controls; project-docs/docs/phases/phase-21/triage-log.md TRG-21-004, TRG-21-009, and Resolution Note For TRG-21-009; project-docs/docs/phases/phase-21/tasks/task-9-done.md Recovery attempt 10 and Concerns; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9.

## Resolution Note For TRG-21-010

- **Decision authority:** DEC-166, 2026-08-26.
- **Verdict audit:** **not sustained on its recorded category.** The product
  interface is not contradictory. The red file and functions were introduced
  under `Refs: phase-20/task-10`; current Phase 21 task 10 is generic Claims and
  owns none of the affected inventory leaves. The test filename was treated as
  ownership evidence even after DEC-165 established invalidating-change
  ownership. Repetition of that inference did not create a protected category.
- **Resolution:** Task 9 converges the stale Phase 20 request rows to the exact
  accepted `--actor` plus idempotency contract, renames historical `Task10`
  fixture labels to subject-based names, and preserves the actual subjects:
  missing/malformed attribution rejects before mutation, JSON and Human errors
  match the current public contract, and journal state is unchanged. Unrelated
  assertions remain byte-identical. No production behavior or accepted
  contract changes.
- **Recovery:** RUN-RESUME from the STOP retains TRG-21-010 and ESC-21-008,
  uses resolution authority DEC-166, and binds one same-task TASK-RECOVERY for
  task 9 with trigger signature
  `task-9:scope-boundary:task-9-done.md:task9-converted-leaf-actor-contract-retained-legacy-by-assertions`,
  category `in-scope-implementation`, mode `task-scope`, attempt 12, `Source`
  and `Repair authority` `TRG-21-010 / DEC-166`, and no-progress `0/2`. All
  five task commands rerun byte-identically.
- **Historical integrity:** TRG-21-010, ESC-21-008, their ledger TRIAGE/STOP,
  the BLOCKED report, and all failed command evidence remain unchanged. The
  Phase Runner appends the recovery events; this note does not edit the ledger.

---

## TRG-21-011 - Task 10 supplementary closure failure is Task 12 transition integration

- **Trigger:** Task 10 continuation returned a typed BLOCKED report after the supplementary retained-integration gate failed at a Claim-dependent closure transition before reaching the removed legacy Claim command.
- **Originating task:** task-10
- **Trigger signature:** task-10:integration-fixture:task-12-done.md:retained-claim-transport-collateral
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-21/tasks/task-10-done.md
- **Candidate:** typed `OUTCOME BLOCKED` report after supplementary tagged integration failure
- **Attempt count:** 2
- **Checkpoint:** report commit 0a0fe828d9f51cd86252e4ed50163f46fcf4bec8; prompt-attached `git rev-parse --verify 0a0fe828d9f51cd86252e4ed50163f46fcf4bec8^{commit}` returned the exact hash, subject matched `docs(task): record phase-21 task-10 blocker`, and two-argument report-lint returned `OUTCOME BLOCKED` with `PASS report-lint completed with no failures` (terminated `exit_code: 0`).
- **Verdict:** RESUME
- **Reasoning:** This is not an unresolved Task 10 task-scope defect or a protected contract conflict. Task 10's exact five-command gate passed, and its accepted generic Claim transport correctly writes `work_claims` while keeping legacy issue-claim aliases absent. The reproduced supplementary journey instead fails at the retained transition's V1 `claims.RequireLiveIssueClaimHeld` read before it reaches the removed legacy command; the convergence inventory assigns V1 `claims` readers to tasks 10, 12, and 20, and Task 12 expressly owns the execution-bound transition integration. DEC-162 requires assertion-preserving convergence, not a Task 10 workaround that weakens the closure assertion or changes Task 12 behavior.
- **Local decisions made:** Task 10 may advance from this misclassified BLOCKED disposition; preserve its failed supplementary gate evidence. Route the actual repair through the already-open Task 12 `RECOVERY_REQUIRED` continuation: integrate the retained transition's Claim observation with the accepted generic Claim contract, preserve the closure journey's assertion strength, and rerun Task 12's unchanged six commands including the tagged integration gate. This TRG grants no Task 10 product edit and no change to commands, assertions, public interface, persistence contract, task specifications, QA artifacts, or validation contract.
- **Resolution evidence:** `task-10-generic-claims.md` limits Task 10's required gate to five commands, all recorded PASS in the source report. `task-12-execution-bound-transition.md` owns `internal/issues/`, `internal/workflow/`, and the tagged integration gate; `v1-convergence-inventory.md` assigns the V1 `claims` reader boundary to tasks 10, 12, and 20. The attached reproducer fails before the removed `issue claim` call, at `cmd/fid/integration/e2e_v01_closure_test.go:33`, because `internal/issues/transition.go:221` still requires a live V1 Claim. No canonical workflow-control mismatch, protected definition change, or unresolved product ambiguity is reproduced.
- **If REPAIR:** N/A; the appropriate next route is the existing task-12 `RECOVERY_REQUIRED` same-task continuation under the task-12 specification and DEC-162, not a new task-10 recovery.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** no workflow verifier repair; no command, assertion, public-interface, specification, QA-artifact, validation-contract, persistence, or task-definition change; no weakening of retained V1 closure coverage; no reintroduction of issue claim/unclaim aliases or `--by`.
- **Exact verification:** N/A
- **Originating report:** project-docs/docs/phases/phase-21/tasks/task-10-done.md
- **Control repair base:** N/A
- **If BLOCK:** N/A
- **Human attention required:** no
- **References:** project-docs/docs/phases/phase-21/tasks/task-10-generic-claims.md Exact Commands To Run and Scope; project-docs/docs/phases/phase-21/tasks/task-12-execution-bound-transition.md Objective, Scope, and Exact Commands To Run; project-docs/docs/phases/phase-21/tasks/task-10-done.md lines 5-10 and 72-80; project-docs/docs/phases/phase-21/tasks/task-12-done.md lines 5-9 and 13-14; project-docs/docs/phases/phase-21/v1-convergence-inventory.md sections 1, 1.1, and 4; project-docs/docs/phases/phase-21/decisions.md DEC-162; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9.
---

## TRG-21-012 - Task 12 recovery budget exhausted on owner-routed gate failures

- **Trigger:** Task 12 recovery budget exhausted after two consecutive runner-mediated continuations failed to clear the unchanged required gate; the remaining failures are owner-routed to Task 9 and future Task 20.
- **Originating task:** task-12
- **Trigger signature:** task-12:gate-repair:task-12-done.md:tagged-integration-retained-claim-transition-collateral
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-21/tasks/task-12-done.md
- **Candidate:** recovery-budget-exhausted after Task 12 continuation attempt 5
- **Attempt count:** 5
- **Checkpoint:** af4c9ecf9d1be5668725d8287fde5bc531ea2fe3; prompt-attached `git rev-parse --verify` succeeded and subject matched `docs(task): record task-12 attempt-5 gate evidence`.
- **Verdict:** BLOCK
- **Reasoning:** The same trigger reached two consecutive Phase Runner-mediated no-progress continuations. Attempt 5 reproduced the unchanged Task 9-owned tagged-integration failure and future Task 20-owned repository-coverage failure; neither is within Task 12 authority, and Task 12 has completed its accepted implementation obligations. Further advance would require selecting cross-task owner/order handling, while restoring `--by`, the obsolete V1 Claim behavior, or changing a required command or assertion is prohibited. No canonical-control mismatch or control-conformance proposal is presented.
- **Local decisions made:** none
- **Resolution evidence:** `go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...` passed with `ok github.com/mbell-nzl/fid/cmd/fid 173.260s` (terminated `exit_code: 0`); `go test -tags integration -race ./cmd/fid/integration/...` failed with `FAIL` (terminated `exit_code: 1`) because the retained fixture reaches Task 9's `issue ac verify --status pass --by qa_lead` and receives `error: unknown flag: --status`; `gofmt -l cmd internal` passed with `exit 0 (no output)` (terminated `exit_code: 0`); `golangci-lint run ./...` passed with `0 issues.` (terminated `exit_code: 0`); `go test -race -coverprofile=coverage.out ./...` failed with `FAIL` (terminated `exit_code: 1`) because future Task 20-owned `internal/dev/TestStartRejectsInvalidWorkflowEdgeBeforeAppend` expects the removed V1 Claim precondition; and `go tool cover -func=coverage.out` passed with `total: (statements) 65.3%` (terminated `exit_code: 0`). The valid truthful Task 12 report remains `OUTCOME RECOVERY_REQUIRED`.
- **If REPAIR:** N/A
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Task 12 must not edit Task 9 or Task 20 fixtures; restore `--by`; restore obsolete V1 Claim behavior; weaken assertions; change required commands, command definitions, accepted status, scope, threshold, budget, tolerance, policy, or any protected contract; no control-conformance directive.
- **Exact verification:** N/A
- **Originating report:** project-docs/docs/phases/phase-21/tasks/task-12-done.md
- **Control repair base:** N/A
- **If BLOCK:** category recovery-budget-exhausted / protected task-boundary and phase-order authority; ESC-21-009 raised in `tl_escalations_current.md`; resolution condition a Human/full Tech Lead decision selects the owner and order handling for the remaining Task 9 and future Task 20 gate failures, after which the phase may resume without changing accepted contracts.
- **Human attention required:** yes
- **References:** project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9; project-docs/docs/phases/phase-21/run-ledger.md lines 1968-2002; project-docs/docs/phases/phase-21/tasks/task-12-execution-bound-transition.md Scope and Exact Commands To Run; project-docs/docs/phases/phase-21/tasks/task-12-done.md Workflow outcome and Exact quality-gate results; project-docs/docs/phases/phase-21/tasks/task-9-retained-actor-policy-removal.md Scope and Out Of Scope; project-docs/docs/current/phase_brief.md Group 2 task order.

## Resolution Note For TRG-21-012

- **Decision authority:** DEC-167, 2026-08-28.
- **Verdict audit:** **not sustained.** The Developer's refusal to widen scope or
  weaken assertions on child authority was correct, and its recorded evidence is
  accurate. Every other element of the disposition fails:
  - The `internal/dev` failure is task 12's own collateral, not task 20's. Task
    12 removed the V1 Claim precondition from `internal/issues/transition.go`;
    `internal/dev/start_test.go:225-233` still requires it.
    `v1-convergence-inventory.md:54` assigns task 20 the *source-package
    removal*, which is not test-convergence ownership.
  - The `issue ac verify` failure is not a `--by` to `--actor` transport gap.
    VAL-030 makes the leaf unconditionally `capability-unavailable`
    (`cmd/fid/issue_ac.go:109-111`), so no flag change can make the journey step
    succeed. Task 9 correctly excluded it from its conversion family.
  - It is a chain, not a line: `concern raise` and `concern address`
    (`e2e_v01_closure_helpers_test.go:263`, `:266`) are the same VAL-030
    fail-closed family and fail next.
  - `recovery-budget-exhausted` does not hold. Run-ledger :1904 forbade task 12
    from repairing the only two failing rows; attempt 4 was report-only and
    attempt 5 changed nothing. Both NO-PROGRESS outcomes were guaranteed by
    directive.
- **Root cause:** gate coverage, not ownership. `cmd/fid/integration/` is
  `//go:build integration`, and only tasks 12, 18, 21, and 22 carry a tagged
  command. DEC-165 gave task 9 ownership of that package without a command that
  runs it. Tasks 8, 9, and 10 each observed the red package and correctly walked
  away; task 12 is the first task obligated to make it green.
- **Resolution:** DEC-167 amends task 12's file scope to name `internal/dev/` and
  the `cmd/fid/integration/` package, assigns it the accumulated tagged-journey
  convergence with per-row provenance recording, and keeps the no-weakening
  boundary. It adds two standing rules to `08_QUALITY_PLAN.md` §3: required
  commands must exercise owned scope, and a no-progress budget counts only
  attempts whose directive permitted causal repair.
- **Recovery:** RUN-RESUME from the STOP retains TRG-21-012 and ESC-21-009, uses
  resolution authority DEC-167, and binds one same-task TASK-RECOVERY for task 12
  with trigger signature
  `task-12:gate-repair:task-12-done.md:tagged-integration-retained-claim-transition-collateral`,
  category `in-scope-implementation`, mode `task-scope`, attempt 6, `Source` and
  `Repair authority` `TRG-21-012 / DEC-167`, and no-progress `0/2`. All six
  task-12 commands rerun byte-identically. Reaching a later journey step resets
  the no-progress count. ESC-21-009 is resolved.
- **Recurrence note:** fourth occurrence of the TRG-21-004 pattern after
  TRG-21-009 and TRG-21-010. DEC-165 and DEC-166 fixed the ownership rule; the
  surviving generator was unexecutable ownership, which DEC-167 item 7 closes.
- **Historical integrity:** TRG-21-012, ESC-21-009, their ledger TRIAGE/STOP, the
  attempt 1-5 evidence, and all failed command results remain unchanged. The
  Phase Runner appends the recovery events; this note does not edit the ledger.

---

## TRG-21-013 - Group 1 independent checkpoint QA capability unavailable after bounded sessions

- **Trigger:** the required Phase 21 Group 1 independent checkpoint has no usable QA artifact after three fresh bounded QA Lead sessions; each exceeded the declared 900-second FEAS-set budget and was explicitly interrupted before returning a checkpoint verdict. This is a possible unavailable external capability / bounded QA-session liveness blocker, not a product finding.
- **Originating task:** N/A
- **Trigger signature:** phase-21:group-1-checkpoint:qa-lead-feasibility-session-interrupted
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** Phase 21 Group 1 checkpoint QA session; preserved fresh scratch paths `/tmp/tmp.YC4EUO3gMN` and `/tmp/tmp.zQAE0hBrd2`.
- **Candidate:** unavailable external capability / bounded QA-session liveness blocker
- **Attempt count:** 3
- **Checkpoint:** phase HEAD `6caaa6f7dc7cca0c5e9e4306f81ae53228f25d25`, verified by prompt-attached `git rev-parse --verify HEAD^{commit}` result and subject match `docs(qa): pass phase 21 window 4 re-review`; required artifact `project-docs/docs/phases/phase-21/qa/group-1-checkpoint-review.md` is absent.
- **Verdict:** BLOCK
- **Reasoning:** The checkpoint is an accepted Group 2 dispatch prerequisite and requires one fresh independent artifact that separately decides VAL-042 and VAL-043 and records FEAS-01 through FEAS-19 with observed termination. Three fresh bounded QA Lead attempts each exhausted the declared 900-second FEAS-set budget and returned `STATUS: INTERRUPTED`, without a verdict, artifact, commit, or command session. No evidence establishes that a fourth identical bounded QA-session retry can complete, while task 13 cannot advance on partial, undocumented, or transcribed evidence.
- **Local decisions made:** none
- **Resolution evidence:** Required result: one fresh independent Group 1 checkpoint artifact, separately deciding VAL-042 and VAL-043 and recording FEAS-01..FEAS-19 with observed termination, committed before task 13 dispatch. Observed result: attempts 1, 2, and 3 each returned `STATUS: INTERRUPTED` after the 900-second budget, with no artifact, commit, or command session; attempt 3 reached child build and fresh FEAS setup/materials but still produced no artifact. The liveness route returned `ACTION REPAIR` for each terminated interruption, but the bounded recovery has now produced no usable checkpoint artifact after all three fresh attempts.
- **If REPAIR:** N/A
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** do not edit product code, tests, task specs, QA plan/contract, validation contract, decisions, triage/escalations, or run ledger; do not weaken checkpoint requirements, FEAS budgets, observed-termination rules, or any protected contract.
- **Exact verification:** N/A; Phase Runner will verify any later checkpoint artifact and commit independently.
- **Originating report:** N/A
- **Control repair base:** 6caaa6f7dc7cca0c5e9e4306f81ae53228f25d25
- **If BLOCK:** category unavailable external capability / bounded QA-session liveness; ESC-21-010 raised in `tl_escalations_current.md`; resolution condition a full Tech Lead/Human restores or provides an independent QA capability that can execute the unchanged checkpoint and return one complete artifact, or records an authorized protected resolution for the checkpoint/capability route.
- **Human attention required:** yes
- **References:** project-docs/docs/current/phase_brief.md Group 1 and Dispatch Controls; project-docs/docs/phases/phase-21/qa/group-1-checkpoint-contract.md Verdict; project-docs/docs/phases/phase-21/qa/qa-plan.md sections 10-11; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9; project-docs/templates/triage-log-template.md.

---

## TRG-21-014 - Corrected task identity for the blocked Group 1 checkpoint

- **Trigger:** correction of malformed TRG-21-013 identity for the already-blocked Phase 21 Group 1 independent checkpoint.
- **Originating task:** task-13
- **Trigger signature:** task-13:group-1-checkpoint:qa-lead-feasibility-session-interrupted
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** correction source TRG-21-013 / existing ESC-21-010; Phase Runner schema-v2 append rejection because TRG-21-013 recorded `Originating task: N/A` against a concrete checkpoint trigger.
- **Candidate:** already-blocked unavailable external capability / bounded QA-session liveness condition; orchestration-record identity correction for ledger consumption.
- **Attempt count:** 3
- **Checkpoint:** current phase HEAD `6caaa6f7dc7cca0c5e9e4306f81ae53228f25d25`, verified by prompt-attached `git rev-parse --verify HEAD^{commit}` result and subject match `docs(qa): pass phase 21 window 4 re-review`; required artifact `project-docs/docs/phases/phase-21/qa/group-1-checkpoint-review.md` remains absent.
- **Verdict:** BLOCK
- **Reasoning:** The independent checkpoint remains an accepted task-13 dispatch prerequisite. It requires one fresh artifact that separately decides VAL-042 and VAL-043 and records FEAS-01 through FEAS-19 with observed termination. Three fresh independent QA Lead sessions each returned `STATUS: INTERRUPTED` after the declared 900-second FEAS-set budget with no artifact, commit, verdict, or command session; attempt 3 reached fresh child build and FEAS setup but still produced no artifact. This corrects only the malformed external identity so the schema-v2 ledger can bind the existing stop; it does not create a new capability result or escalation.
- **Local decisions made:** corrected the external triage identity to task-13 and task-13:group-1-checkpoint:qa-lead-feasibility-session-interrupted for ledger consumption; retained the existing BLOCK and ESC-21-010 without dispatching task 13.
- **Resolution evidence:** none. The existing protected unavailable-QA-capability condition remains: a fresh independent checkpoint artifact must separately decide VAL-042 and VAL-043 and record FEAS-01..FEAS-19 with observed termination before task 13 dispatch; no such artifact exists.
- **If REPAIR:** N/A
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** do not weaken or bypass the accepted checkpoint contract, VAL-042/VAL-043, FEAS-01..FEAS-19, per-probe/whole-set budgets, or observed-termination requirements; do not dispatch task 13; do not edit the run ledger or protected project artifacts.
- **Exact verification:** N/A; Phase Runner will verify any later checkpoint artifact and commit independently.
- **Originating report:** N/A
- **Control repair base:** 6caaa6f7dc7cca0c5e9e4306f81ae53228f25d25
- **If BLOCK:** category unavailable external capability / bounded QA-session liveness; ESC-21-010 referenced in `tl_escalations_current.md`; resolution condition a full Tech Lead/Human restores or provides an independent QA capability that can execute the unchanged checkpoint and return one complete artifact, or records an authorized protected resolution for the checkpoint/capability route.
- **Human attention required:** yes
- **References:** TRG-21-013; ESC-21-010 in project-docs/docs/current/tl_escalations_current.md; project-docs/docs/current/phase_brief.md Group 1 and Dispatch Controls; project-docs/docs/phases/phase-21/qa/group-1-checkpoint-contract.md Verdict; project-docs/docs/phases/phase-21/qa/qa-plan.md sections 10-11; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9; project-docs/templates/triage-log-template.md.

## Resolution Note For TRG-21-013 And TRG-21-014

- **Decision authority:** DEC-168, 2026-08-29.
- **Verdict audit:** **not sustained on the recorded category.** The blocking
  condition was real - no `qa/group-1-checkpoint-review.md` exists - and refusing
  to dispatch task 13 on partial or transcribed evidence was correct. The
  category was wrong. The QA capability is available: preserved scratch
  `/tmp/tmp.zQAE0hBrd2` holds a built `child-bin` and materialized outputs for
  essentially the whole probe set, and a fresh diagnostic tagged integration run
  outside any checkpoint session completed. The real cause is orchestration
  applying the 900-second FEAS-set command budget as the QA Lead child session's
  total lifetime, then restarting each attempt from zero.
- **TRG-21-013:** retained as historical evidence. Its `Originating task: N/A`
  with a `phase-21:` trigger signature cannot bind under schema v2 and is not a
  governing recovery record.
- **TRG-21-014:** retained as the sole governing STOP record; its identity
  correction is mechanically valid and its BLOCK is superseded by DEC-168, not
  rewritten.
- **Resolution:** DEC-168 records the binding reading of the declared budgets -
  120s per probe, 900s for the FEAS set, 300s independently for the tagged
  integration command, and no budget on the child session's lifetime - with no
  edit to `qa-plan.md` or `group-1-checkpoint-contract.md`. VAL-042, VAL-043,
  FEAS-01..FEAS-19, fresh-session evidence, observed termination, commands,
  accepted statuses, and reviewer independence are unchanged, and the checkpoint
  is not split. Preserved scratch and the diagnostic run are inadmissible as
  checkpoint evidence; all probes rerun fresh in newly allocated scratch, and the
  QA Lead writes the artifact incrementally during the attempt.
- **Recovery:** RUN-RESUME from the STOP retains TRG-21-014 and ESC-21-010, uses
  resolution authority DEC-168, continuing task task-13, attempt 4, and next
  recovery tuple
  `task-13 | session-continuation | task-scope | task-13:group-1-checkpoint:qa-lead-feasibility-session-interrupted | 4`,
  then one `TASK-RECOVERY task-13` with `Source` and `Repair authority`
  `TRG-21-014 / DEC-168`, category `session-continuation`, mode `task-scope`, and
  no-progress `0/2`. This is an orchestration wrapper: the Phase Runner dispatches
  the QA Lead first, must not interrupt it at 900 seconds, and must not dispatch
  the task-13 Developer until the committed checkpoint artifact records a verified
  PASS. DEC-168 items 9 and 10 fix the REPAIR/TASK COMPLETE closure and every
  branch. ESC-21-010 is resolved.
- **Recurrence note:** third identity-correction cycle after TRG-21-002/003 and
  TRG-21-005/006, and the first phase-scoped trigger with no valid schema-v2
  identity. The three interrupted attempts also have no ledger event. DEC-168
  item 11 defers both general fixes to workflow maintenance rather than solving
  them inside this phase.
- **Historical integrity:** TRG-21-013, TRG-21-014, ESC-21-010, the ledger TRIAGE
  and STOP, the rejected TRG-21-013 append, and the preserved scratch paths remain
  unchanged. The Phase Runner appends the RUN-RESUME; this note does not edit
  `run-ledger.md`.

## Resolution Note For TRG-21-015

- **Decision authority:** DEC-169, 2026-08-29.
- **Verdict audit:** **sustained.** DEC-168 item 10(a) required exactly this stop,
  and the Phase Runner and TL-Lite followed it precisely. The checkpoint artifact
  is high quality: three separately decided parts, VAL-042 and VAL-043 passing on
  their own evidence, FEAS-01..FEAS-19 fresh, and two decision-free bounded
  findings with no blocker candidate. The only correction is finding scope, not
  triage conduct.
- **MF-1 routing:** correct. `internal/workflow/activation.go:176` evaluates
  `assertActivationAvailable` before `journal.CommitIdempotent` at `:213`, so an
  exact same-key/same-digest replay is rejected by a lifecycle conflict the
  idempotency layer would have served; the same ordering exists in the component
  path and in `transitionActivation`. Task 4 authored and owns that path. DEC-169
  widens the repair to the whole Activation family and requires both streams of
  the pre-repair second call to be recorded, because the review's "no output"
  detail conflicts with the `activation-conflict` envelope mapped at
  `cmd/fid/activation.go:284-294`.
- **MF-2 routing:** right task, wrong-sized defect. `cmd/fid/main.go:5-9` is the
  only product `os.Exit` and is unconditionally `os.Exit(1)`; `internal/cli` has no
  exit concept; 28 `cli.WriteFailure` sites across 16 files return untyped errors.
  The whole section 7 class table is unimplemented, so task 9's four fail-closed
  leaves are the first real-process observation of a phase-wide gap, not a wrong
  value choice. No Group 1 required command runs a real binary and
  `cmd/fid/integration/` asserts no exit status anywhere - the fifth occurrence of
  the DEC-167 item 7 pattern. DEC-169 keeps task 9 as owner and widens the repair
  to the complete category-to-exit mapping, because the alternative accumulates the
  same defect to QA-69 at task 19 and VAL-072 at task 22.
- **Recovery:** DEC-169 item 8 fixes the exact route. Because schema v2 binds a
  resume to its governing TRIAGE's task-prefixed trigger and TRG-21-015 is
  `task-13`-prefixed, reaching task 4 needs an identity-only successor triage,
  exactly as TRG-21-014 succeeded TRG-21-013. TL-Lite allocates **TRG-21-016**
  under DEC-169, retaining ESC-21-011 rather than allocating a new escalation:
  `Originating task: task-4`, `Trigger signature:
  task-4:checkpoint-repair:group-1-checkpoint-review.md:mf-1-activation-idempotent-replay`,
  `Attempt count: 1`, `Repair category: N/A`, `Repair authority: N/A`,
  `Verdict: BLOCK (escalated as ESC-21-011)`. The ledger order is
  `REVIEW infra-review-1` (REPAIR, owners `task-4, task-9`), `TRIAGE TRG-21-016`,
  `STOP`, `RUN-RESUME` (continuing task-4, attempt 2, authority DEC-169), the
  task-4 recovery/repair/complete triple, the task-9 triple, then
  `REVIEW infra-review-1b` PASS after the complete checkpoint rerun. ESC-21-011 is
  resolved.
- **Recurrence note:** second instance of the DEC-168 item 11 identity gap, now
  cross-task rather than phase-scoped: schema v2 cannot express a STOP whose repair
  belongs to a task other than the stopped one. Both instances cost an identity-only
  successor triage. The general fix remains deferred to workflow maintenance.
- **Historical integrity:** TRG-21-015, ESC-21-011, their ledger TRIAGE and STOP,
  the committed checkpoint artifact, and all recorded evidence remain unchanged.
  The Phase Runner appends the new events; this note does not edit
  `run-ledger.md`.

---

## TRG-21-015 - Group 1 checkpoint repair requires completed-task owner routing

- **Trigger:** fresh Group 1 checkpoint evidence has real Must-fix findings after VAL-042, VAL-043, and FEAS-01..FEAS-19 passed; the DEC-168 one-off task-13 wrapper cannot change owner identity while open.
- **Originating task:** task-13
- **Trigger signature:** task-13:group-1-checkpoint:checkpoint-review.md:must-fix-owner-routing
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** project-docs/docs/phases/phase-21/qa/group-1-checkpoint-review.md
- **Candidate:** independent checkpoint REPAIR with MF-1 assigned to completed task-4 and MF-2 assigned to completed task-9.
- **Attempt count:** 1
- **Checkpoint:** committed `9a4ed612abf8430332585555d2b35a8f289c3fe2`; prompt-attached `git rev-parse --verify 9a4ed612abf8430332585555d2b35a8f289c3fe2^{commit}` returned the exact hash and the subject matched `docs(qa): record phase 21 group 1 checkpoint`.
- **Verdict:** BLOCK
- **Reasoning:** The checkpoint is not PASS: its independent contract review returned REPAIR with two reproduced product-conformance defects. MF-1 is bounded to task 4's Activation caller path and MF-2 to task 9's fail-closed process-exit mapping, but DEC-168 item 10(a) expressly prohibits the open task-13 orchestration wrapper from changing to either completed task owner and requires TRIAGE/STOP for ownership routing and the required complete checkpoint repeat. Selecting and sequencing those two task recoveries changes accepted work across task boundaries and cannot be decided locally.
- **Local decisions made:** retained VAL-042 PASS, VAL-043 PASS, and fresh FEAS-01..FEAS-19 PASS as evidence only; they do not convert the checkpoint verdict to PASS. Did not reopen TRG-21-014: its distinct trigger was unavailable checkpoint execution, while this trigger is verified, newly delivered two-owner repair evidence.
- **Resolution evidence:** MF-1 reproduced a byte-identical second `activation prepare` that emitted no output and terminated `exit_code: 1`, rather than returning the original Event and domain fields with `replayed:true` and no appended Event. MF-2 reproduced four `capability-unavailable` envelopes that each terminated `exit_code: 1`, rather than the accepted exit 5, with no journal mutation. Both findings preserve the accepted public contracts and require task-4/task-9 conformance repairs followed by a fresh complete independent checkpoint.
- **If REPAIR:** N/A; the open task-13 wrapper cannot authorize or dispatch completed task-4 and task-9 owner recoveries, and the checkpoint repeat/order requires full Tech Lead authority under DEC-168 item 10(a).
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** do not dispatch task-13 before a verified checkpoint PASS; do not alter the accepted Group 1 checkpoint, VAL-042/VAL-043, FEAS-01..FEAS-19, command budgets/statuses, task specifications, public interface, Event schema, or workflow policy; do not route work by unverified owner inference; preserve the committed checkpoint and all evidence.
- **Exact verification:** N/A
- **Originating report:** N/A
- **Control repair base:** 9a4ed612abf8430332585555d2b35a8f289c3fe2
- **If BLOCK:** category cross-task ownership and phase-order authority under DEC-168 item 10(a); ESC-21-011 raised in `tl_escalations_current.md`; resolution condition a full Tech Lead selects the task-4/task-9 recovery order and exact authority, requires their bounded conformance evidence and reports, then requires the unchanged complete independent Group 1 checkpoint rerun before task-13 dispatch.
- **Human attention required:** yes
- **References:** project-docs/docs/phases/phase-21/decisions.md DEC-168 items 8 and 10; project-docs/docs/current/phase_brief.md Group 1 and Dispatch Controls; project-docs/docs/phases/phase-21/qa/group-1-checkpoint-review.md MF-1, MF-2, and Findings And Verdict; project-docs/docs/phases/phase-21/tasks/task-4-activation-trust.md Scope and Exact Commands To Run; project-docs/docs/phases/phase-21/tasks/task-9-retained-actor-policy-removal.md Scope, Convergence Ownership, and Exact Commands To Run; project-docs/specs/09_AI_WORKFLOW.md sections 2.3 and 9.
