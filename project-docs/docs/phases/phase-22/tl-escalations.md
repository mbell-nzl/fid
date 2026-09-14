# Phase 22 Escalations

## Retained original escalation

### ESC-22-001 - Protected workflow-sequence control cannot represent DEC-189 supersession

- **Phase:** 22
- **Status:** OPEN
- **Raised by:** TRG-22-006
- **Category:** protected workflow-control definition/policy
- **Summary:** `project-docs/tools/workflow-sequence-lint.sh` rejects the user-authorized TRIAGE TRG-22-005 and bound task-2 TASK-RECOVERY despite the accepted contract permitting a later same-task schema-v2 TRIAGE REPAIR to replace its pending dispatch and retain the review-form trigger identity.
- **Evidence:** From `/home/mark/projects/fid`, the supplied exact append reproducer terminated `exit_code: 1` with `INVALID TRIAGE REPAIR must be followed by its bound TASK-RECOVERY before TRIAGE; INVALID schema v2 TRIAGE Originating task does not match Trigger signature`; the ledger remained unchanged and standalone validation remained `VALID state=REPAIR_PENDING`.
- **Why blocked:** The required corrective edit is to `project-docs/tools/workflow-sequence-lint.sh`. `09_AI_WORKFLOW.md` §§2.3 and 9 prohibit normal control-conformance self-repair of that verifier; changing it requires separate protected workflow-control authority.
- **Required resolution:** Full Tech Lead or authorized protected workflow-control decision must reconcile the verifier implementation with `project-docs/templates/run-ledger-template.md` and DEC-189's same-task supersession and review-form trigger contract, then provide a concrete authorized recovery route. Preserve the valid Phase 22 ledger; do not bypass the validator or authorize product/task recovery first.
- **References:** TRG-22-006; `project-docs/specs/09_AI_WORKFLOW.md:427-458,586-614,703-713`; `project-docs/templates/run-ledger-template.md`; `project-docs/docs/phases/phase-22/decisions.md:222-270` (DEC-189); `project-docs/docs/phases/phase-22/run-ledger.md:158-209`.

## Resolution ESC-22-001

- **Status:** RESOLVED, 2026-09-06.
- **Authority:** DEC-190 and DEC-191; Human-authorized workflow maintenance.
- **Repair commit:** `3bb42b4efd391b216b752db2c31fe793f3538bac`.
- **Evidence:** All maintenance gates passed. The real ledger accepted normal
  RUN-RESUME, TRIAGE TRG-22-005 and the exact task-2 attempt-2 recovery, each
  with `APPEND_CONFIRMED state=REPAIR_PENDING` (terminated `exit_code: 0`).
  Standalone sequence: `VALID state=REPAIR_PENDING` (terminated `exit_code: 0`).
- **History:** The rejected BLOCK and STOP appends never persisted. Their
  original failure evidence is retained above and in TRG-22-006. No fabricated
  STOP/resume pair or control-conformance amendment was introduced.
- **Remaining work:** Task-2 product recovery under DEC-189, task-3 MF-3,
  and fresh independent review-1 re-review. No Human decision is pending.
- **Full receipts:** `../../reviews/2026-09-06-phase-22-workflow-recovery.md`.

## Retained original escalation

### ESC-22-002 - Task 4 requires protected checkpoint projection scope and retained-fixture ownership authority

- **Phase:** 22
- **Status:** OPEN (as raised)
- **Raised by:** TRG-22-007
- **Category:** protected task boundary and ambiguous retained-test ownership
- **Why blocked:** Task 4 is contractually responsible for `run.checkpoint` and the `run.checkpointed` Event, but its recorded source boundary excludes `internal/journal/project_effects.go`, whose `unavailableRecoveryEvents` currently rejects the Event before its accepted `run_observations` and `run_runtime` effects. The same required full gate fails in a fixed-date Task 2 Stop fixture; supplied provenance identifies its introducing commit, while no accepted task-4 change is shown to have displaced the assertion. Gate reach alone cannot transfer its ownership.
- **Checkpoint:** Runner-attested implementation commit `ac5e6f9bec80e06bfe5b7b33e3611a894d16038a` and current completion-report HEAD `79c4bd10f84ba7fa43c467b10952055dcc1c03d9`, with the supplied subject/commit-verification evidence.
- **Evidence:** The task-4 report records that `run.checkpointed` cannot apply while the projection rejects it and that `go test -race -coverprofile=coverage.out ./...` terminated `exit_code: 1` in `TestEquivalentStopsShareAttentionButKeepDistinctFacts` with `system clock is behind latest journal date 2026-09-06; refusing to append out of order`. The fixture's fixed Event date is 2026-09-05 and its supplied introducing commit is `ac827959f8359473cbd6546ac6bea1ff1fd20751` (phase-22/task-2).
- **Question / proposed direction:** Decide whether to grant a bounded task-4 amendment for the canonical journal projection that implements the already accepted `run.checkpointed` effects, and separately establish the causal owner and safe bounded repair scope for the retained fixed-date fixture. Preserve the accepted event contract, required gate bytes, pass/fail semantics, and task identities unless a recorded decision changes them.
- **Resolution condition:** Full Tech Lead records concrete authority for each obligation, including exact allowed files, retained-fixture owner/provenance, required verification, and any fresh independent review; a Human decision is required if the authority changes the accepted journal/event contract, task allocation, or test semantics.
- **Links:** TRG-22-007; `project-docs/docs/phases/phase-22/tasks/task-4.md`; `project-docs/docs/phases/phase-22/tasks/task-4-done.md`; `project-docs/docs/phases/phase-22/schema.md`; `project-docs/docs/phases/phase-22/public-interface.md`; `project-docs/docs/phases/phase-22/tasks/task-2.md`; `project-docs/specs/09_AI_WORKFLOW.md` §2.3.
- **Human attention required:** yes

## Resolution ESC-22-002

- **Status:** RESOLVED, 2026-09-06.
- **Authority:** DEC-192 and DEC-193. No Human decision was required: nothing
  accepted changed, and neither obligation altered the journal/event contract,
  task allocation, gate definition, or test semantics.
- **Obligation 1, checkpoint projection.** `schema.md` section 5 already owns
  `run.checkpointed` to task 4, and `unavailableRecoveryEvents` is task 1's
  progressive availability fence whose per-key removal is how the owning task
  implements its Event. Tasks 5, 7, 8 and 10 already carry `internal/journal/`;
  task 4's file list omitted it. That is a Tech Lead planning omission under
  DEC-167/173, corrected by a bounded amendment: delete exactly the
  `"run.checkpointed"` key and add a new task-4-owned projection test. DEC-192
  pre-authorizes the same removal for tasks 5, 7, 8 and 10 and schedules task
  1's displaced VAL-003 fence assertion as convergence for tasks 5 and 10.
- **Obligation 2, retained fixed-date fixture.** Introducing provenance is
  `ac827959f8359473cbd6546ac6bea1ff1fd20751` (phase-22/task-2). The mechanism is
  calendar-only: the fixture pins 2026-09-05 while the workspace initializer
  writes a live-clock Event, so the writer's out-of-order guard correctly
  refuses the backdated append from 2026-09-06. Task 4's commit touches nothing
  on that path, so it is not task-4 convergence, exactly as TRG-22-007 held. It
  is the already accepted DEC-189 class of nondeterministic retained test:
  assigned to the observing task by bounded amendment, repaired at its source,
  with the introducing task keeping subject ownership. DEC-193 records the
  standing ban on absolute calendar dates for writer-appended Events.
- **Executability:** Task 4's unchanged full gate walks both amended packages;
  DEC-193 adds one repeat command proving stability rather than date luck.
- **TL-Lite audit:** The BLOCK was correct. TRG-22-007 refused to extend a task
  boundary and refused to infer ownership from gate reach, which is exactly the
  standing rule. The only gap was that DEC-189 already governed obligation 2;
  applying an existing decision to a new instance is Tech Lead work, so the
  escalation was still the right call.
- **References:** TRG-22-007 and its resolution; `../phase-22/decisions.md`
  DEC-192/193; `../phase-22/execution-scope.md` section 3;
  `../phase-22/tasks/task-4.md`; `../phase-22/tasks/task-{5,7,8,10}.md`.

## Retained original escalation

### ESC-22-003 - Task 6 lacks canonical checkpoint and inspection domain behavior

- **Phase:** 22
- **Status:** OPEN
- **Raised by:** TRG-22-008
- **Category:** protected task boundary and ambiguous owner for required canonical workflow behavior
- **Summary:** Task 6's real built-binary journey reaches absent `fid run inspect` and fails with `exit_code: 1` / `FAIL`. Its accepted source scope is CLI and integration only, while the required canonical `internal/workflow` checkpoint mutation/Event producer and inspection read API are absent.
- **Required resolution:** Identify the accepted owner for those domain operations and record a bounded task-scope amendment or other accepted authority that preserves the existing public-interface, schema/Event, VAL, gate, and workflow-control contracts, naming exact permissible paths, required verification, and review route before a recovery dispatch.
- **Evidence:** `tasks/task-6-done.md` records the focused real-process journey's absent inspect leaf and the missing workflow APIs; Runner verified task implementation commit `c69456fc9bf42f811e4999f1968a48c4a82d3cc5` and report HEAD `7f91325ee2c96d50205dca9564806080dc4e3059`; independent report-lint returned `OUTCOME RECOVERY_REQUIRED` (terminated `exit_code: 0`).
- **Human attention required:** yes

## Resolution ESC-22-003

- **Status:** RESOLVED, 2026-09-06.
- **Authority:** DEC-194. No Human decision was required: an accepted allocation
  already existed, nothing accepted changed, and no contract, VAL row, gate,
  validator, or workflow control was altered.
- **The owner was never ambiguous.** Four accepted records agree that task 4
  owns both domain operations: `tasks/task-4.md` names
  `internal/workflow/run_checkpoint.go` and its test in its file scope;
  `schema.md` section 5 assigns the `run.checkpointed` producer to task 4;
  `public-interface.md` section 3 lists `run.inspect` and `run.checkpoint` as
  `4,6` on the same domain-then-CLI split used by `run.observe` (3,6) and
  `run.recover` (5,6); `contracts.md` sections 3/3.1 place both subjects and the
  observer in task 4. Task 4 delivered the observer, reconciliation, the
  `run.checkpoint` capability, and the DEC-192 fence removal, then closed
  COMPLETE. `internal/workflow/run_checkpoint.go` has never existed in any
  commit on any branch and no inspection read API was written. That is
  undelivered accepted in-scope work, which is bounded repair, not a protected
  task boundary and not a Human decision.
- **The enabling defect is Tech Lead planning.** Task 4's acceptance criteria
  demanded VAL-013..016, `run.checkpointed` projection parity, and the DEC-193
  row, but never the producer or read API its own file scope named. A truthful
  COMPLETE was therefore reachable with two named production files unwritten,
  and the gap stayed invisible until the first process journey two tasks later.
  A file scope is not self-enforcing; carried to wrap-up.
- **Resolution route.** Reopening task 4 is rejected: it is closed, tasks 5 and
  6 ran on top of it, and the Runner's only recovery route is a same-task
  continuation of the active task. DEC-194 therefore moves execution of exactly
  four `internal/workflow` files into task 6, the first task whose accepted
  validation cannot pass without them (VAL-022 names `inspect`; VAL-023 covers
  every section-3 recovery leaf), with task-4 provenance recorded rather than
  relabelled. Task 4 stays COMPLETE and its report is preserved unedited.
- **Executability:** task 6's existing required commands already compile and run
  `internal/workflow` (`go test -race -coverprofile=coverage.out ./...`) and
  build the binary exposing the leaves (tagged integration). No gate, tag, or
  environment switch is added, so the assignment is executable where it lands.
- **TL-Lite / Runner audit:** The BLOCK was correct in effect and wrong in one
  reason. Refusing to fabricate domain behavior inside a CLI-scoped task and
  refusing to widen a task boundary without authority is exactly the standing
  rule, and the escalation was the right call. But "ambiguous owner" was not
  established: the task-4 spec, `schema.md` section 5, `public-interface.md`
  section 3, and `contracts.md` section 3.1 all name the owner, and none was
  consulted before classification. Missing behavior is not the same as
  unassigned behavior. Check the accepted ownership records before classifying
  an owner as ambiguous.
- **References:** TRG-22-008; `../phase-22/decisions.md` DEC-194;
  `../phase-22/tasks/task-4.md`; `../phase-22/tasks/task-4-done.md`;
  `../phase-22/tasks/task-6.md`; `../phase-22/schema.md` section 5;
  `../phase-22/public-interface.md` section 3; `../phase-22/contracts.md`
  sections 3/3.1/3.2.

## Retained original escalation

### ESC-22-004 - Task 8 cancelled Receipt requirement contradicts canonical terminal-status DDL

- **Phase:** 22
- **Status:** OPEN (as raised)
- **Raised by:** TRG-22-009
- **Category:** specification-contradiction
- **Summary:** Task 8's accepted VAL-031/contracts behavior requires a fid-observed cancellation to persist a durable Receipt with `terminal_status` `cancelled`, but canonical `execution_receipts.terminal_status` DDL excludes that value.
- **Evidence:** `contracts.md` section 7 states that cancellation adds terminal outcome `cancelled`; `migrations/0013_p21_operational_kernel.up.sql:250-268` permits only `succeeded`, `failed`, `timed-out`, `signaled`, and `spawn-failed`. The committed task-8 report preserves focused passing cancellation/workflow evidence and records the required gates `BLOCKED`. Runner-supplied verification confirms implementation `b965b24` and report `7643961`; independent report-lint returned `OUTCOME BLOCKED` and `PASS report file committed and clean` (terminated `exit_code: 0`).
- **Why blocked:** The current task forbids DDL, wire-shape, and contract changes. Mapping cancellation to an existing permitted status violates the accepted cancellation Receipt behavior, while modifying the DDL, schema/task scope, VAL-031, or contracts requirement is protected authority outside TL-Lite and task 8.
- **Required resolution:** Full Tech Lead/Human authority must reconcile the accepted durable `cancelled` Receipt requirement with the canonical `execution_receipts.terminal_status` CHECK, record which side is amended, and provide the authorized same-task continuation and verification path. Preserve all failed receipts and do not substitute an existing terminal status.
- **Human attention required:** yes

## Resolution ESC-22-004

- **Status:** RESOLVED, 2026-09-07. Governing decision: DEC-195.
- **The contradiction was real, and it was a Tech Lead planning omission.** Every
  other accepted record agrees that cancellation is its own terminal outcome:
  `contracts.md` section 7, `public-interface.md` section 4's
  `gate-unsuccessful`/8 with outcome `cancelled`, VAL-031's mislabelled-outcome
  failure condition and its "exact signal/timeout/cancelled outcomes" note, the
  v1 convergence inventory, and QA plan R-14, which states that task 8's accepted
  change adds the `cancelled` terminal outcome. Only `schema.md` section 1 is
  silent: it lists five 0013 CHECK member-list replacements for Phase 22 and
  omits the one its own accepted contract requires. Migration 0015 never touches
  the table, and `03_TRD_SCHEMA.md` section 8.3 incorporates `schema.md` by
  reference, so the omission propagated into the canonical record. Unlike
  ESC-22-003, no accepted record assigned this work to anyone, so there was
  nothing to deliver and nothing to repair below full Tech Lead authority.
- **The escalation was correct.** TRG-22-009's BLOCK, its
  `specification-contradiction` category, and its refusal to map cancellation to
  `signaled`, `timed-out`, or `failed` were all right. Substituting an admitted
  status would have falsified an accepted, Human-approved behavior and defeated
  the exact failure condition VAL-031 exists to catch. The preserved BLOCKED gate
  rows and failed receipts are the correct artifact. No no-progress attempt is
  counted against the three-attempt history: the blocker was a missing authority,
  not exhausted repair capability.
- **Resolution route - amend the schema, keep the assertion.** DEC-188 already
  accepted the cancellation behavior by Human approval; only its persistence
  derivation was missing, and `schema.md` section 1 is already an accepted
  pre-adoption in-place baseline replacement of 0013 that task 1 executed five
  times. Section 1 now carries two further `execution_receipts` rows: `'cancelled'`
  joins the `terminal_status` member list, and the composite terminal CHECK gains
  exactly `OR (terminal_status = 'cancelled' AND timed_out = 0 AND NOT (exit_code
  IS NOT NULL AND signal IS NOT NULL))`. The clause follows observed behavior:
  `RunLive` classifies context cancellation as `cancelled` and then records
  whichever of signal or exit code the process state carried, so signal-only,
  exit-code-only, and neither must all be admitted, while the impossible pair is
  refused and `timed_out = 1` stays the `timed-out` discriminator.
- **Execution assignment:** task 8, as an authorized same-task continuation with
  task-1 provenance recorded; task 1 stays COMPLETE and unedited. Its scope now
  names the 0013 file limited to those two exact replacements, the retained
  golden DDL string in `internal/db/operational_kernel_migration_test.go` as
  subject-preserving convergence, and `validReceiptListOptions` in
  `internal/evidence/receipt_read.go`, because a stored status that its own read
  API rejects is an incoherent closed domain. `internal/workflow/evaluation.go`
  needs nothing: it derives success from an `outcome == "succeeded"` allowlist.
- **Executability:** task 8's unchanged eight-command gate already executes every
  assigned file. `go test -race -coverprofile=coverage.out ./...` compiles and
  runs `internal/db` against the real migration directory plus `internal/evidence`
  and `cmd/fid`; the tagged integration gate builds the binary and initializes
  fresh workspaces from the same SQL. No gate, tag, or environment switch added.
- **What did not change:** no Event, Effect, table, column, index, wire shape,
  error code, exit code, VAL row or wording, validator, gate command, threshold,
  task allocation, operating mode, dependency, or workflow control. VAL-031 is
  not weakened; it becomes satisfiable as written.
- **Carried to wrap-up:** DEC-194 recorded that a task file scope is not
  self-enforcing. This is the converse: an accepted contract clause is not
  self-enforcing either. A closed-domain member was added in prose with no schema
  row, no task scope, and no acceptance criterion behind it, and nothing failed
  until an implementation reached the CHECK.
- **References:** TRG-22-009; `decisions.md` DEC-195; `schema.md` section 1;
  `contracts.md` section 7; `public-interface.md` section 4;
  `qa/validation-contract.md` VAL-031; `tasks/task-8.md`; `tasks/task-8-done.md`;
  `qa/qa-plan.md` R-14; `v1-convergence-inventory.md`.


## ESC-22-005 — PR test runtime budget / performance disposition

- Raised: 2026-09-13 by root after the Human supplied the PR CI timeout log.
  No TL-Lite TRG was issued for this post-wrap diagnosis.
- Original request: explicit 20-minute Makefile test/coverage timeout, a protected
  gate-budget change. It was never approved or applied.
- Human response: reduce runtime without compromising quality; subsequently
  approved the bounded performance-repair task.
- Resolution: DEC-221/task 20 authorizes fixture optimization, one complete race/
  coverage execution, and exhaustive sharding if justified by measurements.
  Default package timeout, all assertions and product thresholds remain intact.
- Status: decision request RESOLVED. Local/hosted implementation verification is
  tracked in task-20-done.md; this disposition alone does not claim CI PASS.
- Evidence: qa/pr-ci-diagnosis.md retains original failure/proposal and the revised
  recommendation. The archived Human checklist remains accepted.
