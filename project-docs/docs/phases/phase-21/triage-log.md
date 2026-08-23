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
