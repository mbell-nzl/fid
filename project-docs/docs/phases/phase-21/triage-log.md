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
