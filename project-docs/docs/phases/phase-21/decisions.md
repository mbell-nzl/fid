# Phase 21 Decisions

> Purpose: phase-local decision history, planning decisions, spec sync
> verdicts, workflow notes, and eventual wrap-up decisions. DEC numbering is
> global and never resets.

## Decisions

### DEC-152: Phase 20 retrospective changes Phase 21 delivery controls

- **Date:** 2026-08-22
- **Context:** Phase 20 delivered a clean product slice but generated 23 triage
  identifiers, 19 escalations, multiple control repairs, and a 22-report final
  lint sweep. The Human required the recommended deep retrospective before the
  Phase 21 task pack.
- **Decision:** Adopt the seven planning decisions in
  `../../reviews/2026-08-22-phase-20-retrospective.md`: bounded causal task
  ownership, fresh feasibility for every prescribed verifier, task-local
  report-lint, three-task review windows, a mandatory independent Group 1
  checkpoint, retention of strict typed recovery, and no product-slice tooling
  task for CF-029/CF-030.
- **Consequences:** Phase 21 uses two internal groups inside one Heavy phase.
  Group 2 cannot start on implementation momentum alone; an independent PASS
  must establish Group 1 contract closure and standard-library executor
  feasibility. QA-plan creation remains a QA Lead responsibility.
- **Refs:** Phase 20 DEC-151; Phase 20 `triage-log.md` and
  `tl-escalations.md`; `09_AI_WORKFLOW.md` §§4, 6, 10.

### DEC-153: Human accepts Gate D3.2 with all required amendments

- **Date:** 2026-08-22
- **Context:** The Human approved the core one-Issue walking-skeleton proposal
  and supplied binding answers DP-1..DP-7 plus ordering, authority,
  capability, execution, removal, export, and terminology amendments.
- **Decision:** Close Gate D3.2 exactly as recorded in
  `gate-d3-2-proposal.md`. Plan one Heavy Phase 21 with two internal task
  groups. Defer Attention, Stops, budget consumption/exhaustion, recovery,
  automatic respawn, and export import to Phase 22; implement one generic
  gate-evidence Obligation; advertise no incomplete semantic capability; and
  claim only standard-library-enforceable executor controls.
- **Consequences:** ADR-0012 and the roadmap retain one numeric vertical slice
  while naming the internal checkpoint. The exact capability matrix, JSON,
  DDL/Event, trust, idempotency, output, Artifact, removal, and deferral
  contracts must be accepted before implementation. `Routine` is external kit
  policy language, never fid state.
- **Refs:** Human Gate D3.2 direction 2026-08-22; ADR-0011; ADR-0012;
  `10_IMPLEMENTATION_ROADMAP.md` §§5, 10; DEC-152.

### DEC-154: Fix the exact Heavy Phase 21 slice contract

- **Date:** 2026-08-22
- **Context:** DEC-153 requires exact capability, authority, JSON, DDL/Event,
  executor, evidence, idempotency, removal, and deferral boundaries before a
  Developer task may infer implementation behavior.
- **Decision:** Adopt `contracts.md` and `03_TRD_SCHEMA.md` section 8.2 as the
  planning contract presented for Human execution approval. Capability
  availability is binary and conservative; activated profiles are the only
  kit Actor/Gate trust source; Gate Registration precedes Envelope creation;
  Event schema v2 and operation-scoped durable idempotency apply to converted
  mutations; executor policy is limited to portable standard-library controls;
  and evaluation returns every independent reason in fixed order. Use
  migrations 0013/0014 so generic Claims can pass before V1 `claims` is
  removed. Freeze `task_completions` without writers until Phase 23 and
  `agent_runs` until Phase 25 because existing foreign-key/read ownership makes
  immediate deletion a cross-phase change; neither is authoritative Phase 21
  completion or Run state.
- **Consequences:** One Heavy phase uses two internal task groups and an
  independent checkpoint. Group 2 cannot begin until Group 1 contract closure
  and executor/verifier feasibility independently pass. The slice exposes
  export create/validate but keeps `export.portable.v1` unavailable. No
  unsupported containment, Attention, Stop, recovery, respawn, budget
  consumption, import, review, or Human authority behavior may be inferred.
- **Refs:** DEC-152; DEC-153; `contracts.md`; `03_TRD_SCHEMA.md` §8.2;
  `02_TRD.md` §3.1; ADR-0011; ADR-0012.

### DEC-155: Dispose active carry-forwards against the Phase 21 task pack

- **Date:** 2026-08-22
- **Context:** Phase planning must schedule every carry-forward whose trigger
  the slice fires and explicitly re-defer the rest.
- **Decision:** Schedule CF-020 and CF-027 in task 8: Event v2
  sequence replaces synthetic future occurrence time as the ordering mechanism,
  and reopened `internal/issues/create.go` is decomposed within the recorded
  size targets without behavior change. Schedule CF-032 in task 12 before
  `BuildIssueRevisionEffects` changes. Schedule CF-024's 100k-Event rebuild
  target in task 21. Re-defer CF-008, CF-018, CF-026a, CF-026b, and
  CF-028 because their product/file triggers do not fire. Re-defer CF-023,
  CF-029, and CF-030 as temporary V1 delivery-tool improvements rather than
  mixing them into product tasks. CF-017/022/031 remain cancelled/closed.
- **Consequences:** Tasks 8, 12, and 21 carry explicit tests/evidence. No
  unscheduled triggered debt is hidden in the phase brief, and no workflow
  tooling task dilutes the one-Issue product slice.
- **Refs:** `decisions_active.md`; DEC-152; task 8; task 12; task 21.

### DEC-156: Gate evidence guards the terminal-success Issue transition

- **Date:** 2026-08-22
- **Context:** The accepted DEC-153 walking skeleton evaluated gate evidence,
  but its journey ordering allowed an execution-bound Issue transition before
  Receipt linking and did not distinguish a nonterminal transition from the
  terminal-success transition. The Human supplied protected direction that
  the one gate-evidence Obligation guards terminal success.
- **Decision:** For an Envelope-bound Issue, the terminal-success transition
  is permitted only after fid executes the registered Gate, records and links
  its qualifying Receipt, and the gate-evidence Obligation evaluates
  `satisfied`. Envelope completion is a later operation that re-evaluates the
  same Obligation and also requires the selected Issue to be terminal-success.
  Both operations call one shared Obligation evaluator and return the same
  deterministic Obligation reason values; their separate action predicates
  remain outside that evaluator. A nonterminal execution-bound transition may occur earlier
  under the accepted Actor/Envelope/Claim/Run predicates. DEC-153 remains an
  accurate historical decision and is amended, not rewritten.
- **Consequences:** Task 12 builds the shared fail-closed evaluator boundary
  and the nonterminal transition. Task 17 supplies trusted Evidence inputs.
  Task 18 owns the positive terminal transition and completion reuse. The
  real-process journey mutates once, runs the Gate, links/evaluates Evidence,
  transitions terminal-success, completes, then rebuilds and repeats reads
  and negative checks without creating a second Envelope revision.
- **Refs:** Human Phase 21 planning correction 2026-08-22; PRD-WF-004;
  PRD-ENV-008; `gate-d3-2-proposal.md`; `contracts.md` §§6, 10, 16.

### DEC-157: Complete Phase 21 interface and convergence sequencing

- **Date:** 2026-08-22
- **Context:** The first planning pack left operation inputs/results/errors,
  live frame payloads, bootstrap convergence, V1 caller dispositions, and
  several multi-observable validation rows incomplete. Its checkpoint also
  required final V1 absence before the removal task.
- **Decision:** Adopt `public-interface.md` as the complete final-binary wire
  contract and `v1-convergence-inventory.md` as the live caller/reader/writer
  disposition. Closed public enum members use lowercase kebab-case; opaque
  IDs remain identifiers. `project.initialized` is Event v2 with service Actor,
  aggregate sequence 1, null repository/idempotency, and zero domain Effects;
  it never writes `roles`. Migration 0014 and V1 removal are task 20, before
  task 21 final-schema maintenance/init/rebuild/export parity. The Group 1
  checkpoint proves isolation from explicitly listed temporary V1 state, not
  its premature absence. Atomic validation assertions replace bundled rows.
- **Consequences:** `contracts.md`, canonical specs, DDL checks, task specs,
  phase brief, and QA handoff all cite the two new authorities. Task 1 remains
  unauthorized until Human approval and a Ready QA plan are accepted. No
  runtime dependency or OS-containment claim is added.
- **Refs:** Human Phase 21 planning correction 2026-08-22;
  `public-interface.md`; `v1-convergence-inventory.md`;
  `qa/validation-contract.md`; `qa/group-1-checkpoint-contract.md`.

### DEC-158: Resolve QA-plan refinements without weakening the rebuild target

- **Date:** 2026-08-23
- **Context:** The QA Lead committed the Heavy plan at `f5ccbb6` with complete
  VAL mapping and feasibility, but found two blockers: every task specified an
  invalid/self-referential one-argument report-lint command, and the retained
  100k rebuild control produced one 33.46s failure followed by 27.11s and
  24.97s passes on the only verification machine. It also proposed bounded
  evidence refinements R-04..R-19.
- **Decision:** Remove report lint from Developer command gates. The Phase
  Runner runs the exact two-argument control against the committed report and
  task spec before the next dispatch. Keep the TRD `<30 s` target and the
  benchmark's per-sample `b.Fatalf`; do not use best/median tolerance. The
  Phase 21 reference instance is Linux/amd64 on an AMD Ryzen 7 5700G with six
  exposed logical CPUs and 8,133,236 kB recorded memory, using Go 1.26.2;
  kernel, toolchain, CPU count/model, and memory are captured with the gate.
  VAL-083 requires three consecutive fresh `-benchtime 1x` processes, each
  self-asserting below 30 seconds with parity. Any one failure fails the
  assertion and task 21 owns causal optimization, not threshold relaxation.
- **Decision:** Accept R-04..R-08 and R-10..R-18. Human VAL-036 remains one
  parent perception assertion and fans out to HUMAN-01a..01d as required by
  workflow section 6.1. R-19 remains deferred until Human execution approval.
  Correct R-09's status premise: at final task-20 HEAD the role-policy and
  `--by` recipes must terminate `exit_code: 1`, while the broad retained-reader
  and Event recipes terminate 0; every actual match still must classify and
  the completion report records `UNCLASSIFIED=0`.
- **Consequences:** Task 1 remains unauthorized. QA Lead must re-review the
  applied refinements, correct its task-20 recipe expectations, and change the
  QA plan to `Ready` before Human execution approval. The 30-second target,
  Event schema, product surface, phase shape, and task order do not change.
- **Refs:** QA plan BC-01/BC-02 and R-01..R-19 at `f5ccbb6`; DEC-071;
  CF-024; TRD section 9; `qa/validation-contract.md`; tasks 1-22.

### DEC-159: Accept the Heavy Phase 21 pack and authorize execution

- **Date:** 2026-08-23
- **Context:** The QA Lead independently re-reviewed DEC-158 at `b367a3b`,
  committed a `READY` Heavy QA plan at `e8a636e`, closed BC-01/BC-02, accepted
  every R-01..R-19 disposition, and reported no outstanding Must-fix. The Human
  then explicitly approved execution.
- **Decision:** Accept the Phase 21 slice contract, public interface, V1
  convergence inventory, VAL-001..087 validation contract, Group 1 checkpoint,
  Ready QA plan, and ordered 22-task pack for execution. Complete R-19 by
  changing every validation row from `Draft` to `Accepted`. Authorize the Phase
  Runner to dispatch task 1; later tasks remain gated by their recorded
  dependencies, review windows, report-lint controls, and Group 1 checkpoint.
- **Consequences:** Phase 21 moves from Planning to Active. This approval changes
  no product contract, threshold, task scope, dependency, validator, or gate.
  Any protected change discovered during execution still requires its normal
  decision authority; approval is not permission to weaken an assertion.
- **Refs:** Human execution approval 2026-08-23; QA plan readiness commit
  `e8a636e`; DEC-152..158; `qa/qa-plan.md`; `qa/validation-contract.md`;
  `current/phase_brief.md`.

### DEC-160: Preserve role-free bootstrap with a bounded V1 policy vocabulary bridge. Closes ESC-21-001

- **Date:** 2026-08-23
- **Context:** Task 2 correctly implemented the DEC-157 Event-v2 bootstrap:
  `project.initialized` has the service Actor, aggregate sequence 1, null
  repository/idempotency, and zero Effects. Its focused VAL-003..008 tests and
  role-free init/rebuild journey pass. The required package and repository
  gates then fail because `internal/policy.KnownRoles` derives the V1
  permissions vocabulary only from bootstrap role Effects, while inventoried
  V1 mutations intentionally remain until tasks 8-20. Planning removed their
  only intermediate identity source before their replacement owners ran.
- **Decision:** Keep the exact role-free bootstrap and final Phase 21 model.
  Authorize task 2 to add one temporary, nonpersistent V1 compatibility lookup
  inside `internal/policy` under these closed rules:
  1. It is eligible only when the journal contains the exact role-free
     Event-v2 bootstrap and contains no legacy role Effect.
  2. It recognizes exactly `developer`, `qa_reviewer`, `qa_lead`, `tl_lite`,
     `tl`, `workflow_runner`, and `human` for the unchanged V1 permissions
     matrix and enumerated legacy callers.
  3. If any legacy role Effects exist, those facts remain authoritative and
     the fixed vocabulary does not add to them. If neither the exact bootstrap
     nor role Effects exist, lookup stays empty. Unknown names always fail.
  4. It writes no Event, Effect, row, config, export field, or reconstruction
     fact. It is not a kit Actor, Activation, capability, Grant, provenance,
     Evidence, or Human Outcome source and cannot satisfy a generic operation,
     guarded transition, Obligation, or Envelope completion.
  5. No new Phase 21 generic package may call it. Task 20 deletes it with V1
     policy and proves the existing `KnownRoles|policy.Authorize` absence
     control; task 21 still proves final fresh init and rebuild without roles.
  Amend task 2 to own `internal/policy`, focused red/green compatibility tests,
  and its existing byte-identical failed-command reruns. Do not change the V1
  matrix, Event fixture, VAL semantics, required failed commands, or final
  removal boundary.
- **Consequences:** Task 2 has a bounded causal repair that restores the
  repository gate without reintroducing persisted role authority. The bridge
  is an execution-order device, not a new public compatibility promise. The
  Group 1 checkpoint must classify and isolate it; task 20 must remove it.
  ESC-21-001 is resolved and archived after this decision is committed.
- **Refs:** TRG-21-001; ESC-21-001; DEC-157; `contracts.md` section 14;
  `v1-convergence-inventory.md` sections 1, 4-7;
  `tasks/task-2-event-idempotency.md`; `tasks/task-20-v1-removal.md`;
  `tasks/task-2-done.md`; VAL-008, VAL-024, VAL-075, VAL-076, VAL-077.

### DEC-161: Use the valid TRG-21-003 boundary and resume task 2 without verifier repair. Closes ESC-21-002 and ESC-21-003

- **Date:** 2026-08-23
- **Context:** TRG-21-001 wrapped its schema-v2 `Trigger signature` in
  Markdown backticks, so exact external binding and the task-prefix check could
  not both pass. TRG-21-002 removed the backticks but used a `phase-runner:`
  prefix for originating task `task-2`, so it also could not bind. Both records
  are immutable failed evidence. TRG-21-003 uses the required plain,
  task-prefixed scalar, was bound into the ledger, and the final STOP ledger
  validated as `state=STOPPED`. Its BLOCK condition was absence of protected
  resolution authority, not a remaining verifier defect.
- **Decision:** Schema-v2 external identity fields are plain single-line
  scalars, not Markdown code spans; `09_AI_WORKFLOW.md` now states this
  explicitly. Preserve TRG-21-001/TRG-21-002 and all failed append evidence.
  Use TRG-21-003 as the sole governing recovery record. No edit to
  `workflow-sequence-lint.sh`, no control-conformance directive amendment, and
  no historical ledger/triage rewrite is authorized or required.

  The Phase Runner may append exactly one schema-v2 RUN-RESUME after this
  decision commit, with these binding values:
  - `Prior run/session or STOP`: `2026-08-23 12:37 STOP / TRG-21-003 / ESC-21-003`
  - `Governing TRG`: `TRG-21-003`
  - `Governing ESC`: `ESC-21-003`
  - `Trigger signature`: `task-2:control-conformance:workflow-sequence-lint:external-triage-trigger-prefix`
  - `Attempt count`: `3`
  - `Resolution authority`: `DEC-161`
  - `Continuing task`: `task-2`
  - `Next recovery`: `task-2 | in-scope-implementation | task-scope | task-2:control-conformance:workflow-sequence-lint:external-triage-trigger-prefix | 3`

  Its checkpoint and resolution evidence must cite the committed DEC-160/161
  HEAD, the already-valid stopped ledger, and DEC-160's task-2 scope amendment.
  The following TASK-RECOVERY uses category `in-scope-implementation`, repair
  mode `task-scope`, source `TRG-21-003 / DEC-161`, and repair authority
  `DEC-161`. This is not a tools-only control repair. It dispatches only the
  task-2 continuation under DEC-160; task 3 remains dependency-gated.
- **Consequences:** ESC-21-002 and ESC-21-003 are resolved without weakening or
  changing the canonical verifier. During Phase 21, every new external
  schema-v2 identity value must be written without Markdown delimiters and
  begin with its exact originating task. The valid replacement record plus
  append-only history is the standing recovery pattern, so this formatting
  mistake does not require another protected verifier escalation.
- **Refs:** TRG-21-002; TRG-21-003; ESC-21-002; ESC-21-003; DEC-160;
  `09_AI_WORKFLOW.md` sections 2.3-2.4 and 9;
  `project-docs/tools/workflow-sequence-lint.sh` schema-v2 TRIAGE and
  RUN-RESUME checks; `run-ledger.md` final STOP.

### DEC-162: Retained V1 tests converge inside the task whose accepted change breaks them. Closes ESC-21-004

- **Date:** 2026-08-23
- **Context:** After the DEC-160 bridge, task 2's focused evidence passes and the
  policy vocabulary is restored, but the required package and repository gates
  still fail. Every residual failure is a retained V1 *test*, not a product
  defect: about twenty-three call sites invoke `fid rebuild --by <role>` or
  `journal.RebuildWorkspace(..., role)`; a few tests assert V1 seeded `roles`
  rows in a freshly bootstrapped workspace (`cmd/fid/sql_test.go`,
  `internal/sqlquery/sqlquery_test.go`, the closure helpers); and the third
  subtest of `TestKnownRolesTRG002CorruptTailTolerance` uses a schema-2 fixture
  carrying a V1 `actor.role`, which the accepted Event-v2 rules no longer treat
  as un-upgradeable. DEC-157 fixed `project.rebuild` with no leaf flags and a
  zero-Effect role-free bootstrap; DEC-159 approved both. Task 2 listed only
  `cmd/fid/init_test.go` and `cmd/fid/rebuild_test.go` in scope, so the
  Developer had no authority to adapt the other call sites and correctly stopped
  instead of weakening them. The defect is a planning scope omission, not a
  contract contradiction.
- **Decision:**
  1. The accepted contracts stand and are not reopened. No `--by` on
     `fid rebuild`, no bootstrap role Effect, no persisted compatibility row, no
     early migration 0014, and no change to the V1 permissions matrix, the
     Event-v2 fixture, VAL semantics, or task 2's required commands.
  2. Retained V1 tests and test fixtures are not protected contracts. When an
     accepted interface change invalidates a retained test, the task that made
     the change adapts that test inside the same task. Adaptation that preserves
     the test's subject and assertion strength is convergence, not weakening,
     and needs no further escalation.
  3. Task 2's file scope extends to retained Go test files and test fixtures
     repository-wide. Production scope is unchanged except item 6.
  4. Adaptation rules, in priority order:
     - Rebuild call sites: delete the `--by <role>` argument pair from `fid
       rebuild` invocations and the role argument from `journal.RebuildWorkspace`
       calls. Every other argument, expectation, and assertion stays identical.
     - A test whose subject *was* removed V1 rebuild authority (for example
       `cmd/fid/issue_test.go` expecting `rebuild --by intern` to be rejected for
       an unauthorized role) is converted to the accepted behaviour, or deleted
       only when `cmd/fid/rebuild_test.go` already owns an equivalent accepted
       assertion. The report names each deleted assertion and its accepted owner.
     - A test that reads `roles` rows as incidental data repoints the read to
       state the accepted bootstrap or the test's own journey still produces
       (migration-seeded `workflows`/`workflow_states`, `issues`,
       `activity_log`). Only where the subject genuinely requires a legacy role
       row may the fixture append an explicit legacy `roles` Effect; the report
       names each such fixture so task 20 removes it.
     - A test whose subject *was* bootstrap role seeding asserts the accepted
       role-free result instead, as `cmd/fid/init_test.go` already does.
     - A stale Event fixture is repaired to a shape that is genuinely invalid or
       un-upgradeable under the accepted Event-v2 rules while keeping the
       existing error, file, and line assertions.
  5. Forbidden in this recovery: `t.Skip`, build tags, `-run` narrowing,
     deleting or loosening any assertion without a named accepted owner, and any
     change to coverage thresholds, timeouts, required commands, public
     interface, migrations, or production behaviour beyond item 6.
  6. `journal.RebuildWorkspace` drops its now-vestigial `byRole` parameter and
     acquires the lock with the exact owner command `fid rebuild`. The current
     code writes `fid rebuild --by ` into lock metadata for a flag that no longer
     exists, so this is a truthfulness repair, not new behaviour.
     `LockOwner.ActorRole`, `internal/journal/lock.go` role metadata, and V1
     lock-owner strings inside unrelated tests remain inventoried
     replacement-bound state owned by tasks 20-21.
  7. Recovery route: one more same-task, task-scope `in-scope-implementation`
     attempt on task 2 under the TRG-21-004 trigger signature. Attempts 3-5 each
     recorded material progress, so the two-attempt no-progress budget is not
     exhausted. No Human decision is required: DEC-157 and DEC-159 already
     accepted these contracts and this decision changes none of them.
- **Consequences:** Task 2 can close its unchanged required gates without any
  protected contract change. Item 2 becomes the standing convergence rule, which
  matters immediately: task 8's `--by` to `--actor` rename and the task 9/20
  removals invalidate the same retained suites, and those tasks now adapt their
  own collateral rather than raising a new escalation. The rule is recorded in
  `v1-convergence-inventory.md` section 1.2 and `08_QUALITY_PLAN.md` section 3.
  Task 20's removal inventory gains any explicit legacy role Effect fixture that
  item 4 permits. Group 1 checkpoint, task 20, and task 21 obligations are
  unchanged. ESC-21-004 is resolved once this decision is committed.
- **Refs:** TRG-21-004; ESC-21-004; DEC-157; DEC-159; DEC-160;
  `public-interface.md` section 4.1 `project.rebuild`; `contracts.md` section 14;
  `v1-convergence-inventory.md` sections 1, 1.2, 4-6;
  `tasks/task-2-event-idempotency.md`; `tasks/task-2-done.md`;
  `08_QUALITY_PLAN.md` section 3; VAL-003..008, VAL-076, VAL-077.

### DEC-163: A QA REPAIR review carries an ordered owner list; the sequence verifier is corrected to honour it. Closes ESC-21-005

- **Date:** 2026-08-23
- **Context:** Review 1 covers tasks 1-3 and raises two bounded, decision-free
  Must-fix findings: MF-1 owned by task-1 (VAL-001 exhaustive 0013 inventory)
  and MF-2 owned by task-3 (VAL-010 named deferred-surface probes). The
  mid-phase review template already requires one Repair Plan row per Must-fix
  with its own owning task, so a multi-owner REPAIR review is an accepted and
  expected artifact. The ledger REVIEW event and `workflow-sequence-lint.sh`
  model exactly one owner: the runner appended `Owning task: task-1, task-3`,
  the verifier accepted that append, and every subsequent recovery route is now
  rejected. `TASK-RECOVERY task-1` fails
  `QA repair owner task-1, task-3 does not match recovery task task-1`; a
  list-valued recovery task fails the task-ID shape check; a replacement or
  closing REVIEW fails because `review_task_complete` can never become true.
  The append-only CORRECTION at 15:29 preserves the task-3 finding but has no
  verifier semantics. TRG-21-005 and TRG-21-006 correctly returned BLOCK.

  This is a verifier defect independent of the policy chosen: the sequence
  validator admitted an append that puts the ledger into an unrecoverable
  state. A correct validator either rejects a multi-owner REVIEW at append time
  or supports it. It must therefore change under `09_AI_WORKFLOW.md` section
  2.3, which reserves `workflow-sequence-lint.sh` changes to separate protected
  decision authority. This decision is that authority.
- **Decision:**
  1. **Accepted representation.** A REVIEW event with outcome `REPAIR` carries
     exactly one `Owning task` field whose value is an ordered, comma-separated
     list of one or more distinct task IDs, copied in the same order as the
     review artifact's Repair Plan rows. Multiple owners are valid. `none` is
     valid only for `PASS` and `BLOCK`. Splitting one review window into
     several review identities to dodge the field is not the accepted answer.
  2. **Accepted sequencing.** Every listed owner must complete
     `TASK-RECOVERY` -> `REPAIR` -> `TASK COMPLETE` before the open review can
     be closed by a linked `PASS` or replaced by a linked `REPAIR`. Owners are
     repaired one at a time; recovery-task exclusivity already forbids
     interleaving. The list order is the required repair order. Ownership
     transfer to a task outside the list still requires the whole list to
     complete first.
  3. **CORRECTION stays evidence-only.** A `CORRECTION` never mutates review
     ownership, recovery state, or any other verifier state. Sequential owner
     selection through `CORRECTION` is rejected as the accepted mechanism. The
     15:29 CORRECTION remains immutable historical evidence; its narrowing to a
     single current owner is superseded by item 1, and its named first owner
     `task-1` coincides with the list order, so nothing in the ledger is
     contradicted.
  4. **Authorized verifier change.** `project-docs/tools/workflow-sequence-lint.sh`
     is corrected to implement items 1-3 exactly:
     - `Owning task` remains exactly one field per REVIEW event; the existing
       duplicate-field rejection is unchanged.
     - Parse the value leniently by splitting on commas and trimming
       surrounding spaces, then validate each token strictly against
       `^task-[0-9]+[a-z]*$`. Reject an empty token, a duplicate token, and any
       non-conforming token when the outcome is `REPAIR`.
     - A `TASK-RECOVERY` while a review is open is valid when its task is a
       member of the open owner list; the existing rejection message shape
       `QA repair owner <list> does not match recovery task <task>` is
       preserved for a non-member.
     - `REPAIR` and `TASK COMPLETE` mark per-owner progress; `TASK COMPLETE`
       for a member still requires that member's own `TASK-RECOVERY` and
       `REPAIR`, with the existing message preserved.
     - The review-complete gate that guards a linked `PASS`, a linked
       replacement `REPAIR`, and `RUN-END` becomes true only when every listed
       owner is complete; existing messages are preserved.
     - Single-owner behaviour is unchanged. Every existing drill in
       `project-docs/tools/drills/run-drills.sh` must still pass unmodified.
  5. **Required new drills.** Add coverage for: a multi-owner REPAIR review
     driven through both owners to a linked re-review `PASS`; a linked `PASS`
     rejected while only the first owner is complete; a `TASK-RECOVERY` for a
     non-member task rejected; a duplicate owner in the list rejected; a
     malformed owner token rejected; `Owning task: none` with outcome `REPAIR`
     still rejected; and a duplicate `Owning task` field still rejected.
  6. **Execution route.** This change is workflow infrastructure, not phase-21
     product work and not a control-conformance recovery. A control-conformance
     route is impossible here: the ledger's open multi-owner review rejects the
     very `TASK-RECOVERY` that would authorize the repair. The Human, or a
     Human-directed session outside the phase loop, applies items 4-5, runs
     `project-docs/tools/drills/run-drills.sh` to PASS, confirms
     `project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-21/run-ledger.md`
     reports `VALID state=STOPPED`, and commits on `phase-21`. No historical
     ledger, triage, review, or escalation record is edited. The Phase Runner
     does not perform this change and does not resume before it lands.
  7. **Phase-21 recovery.** Once item 6 lands, the Phase Runner appends exactly
     one RUN-RESUME with these binding values:
     - `Prior run/session or STOP`: `2026-08-23 15:29 STOP / TRG-21-006 / ESC-21-005`
     - `Governing TRG`: `TRG-21-006`
     - `Governing ESC`: `ESC-21-005`
     - `Trigger signature`: `task-1:orchestration-integrity:review-1:multi-owner-ledger-recovery-binding`
     - `Attempt count`: `2`
     - `Resolution authority`: `DEC-163`
     - `Continuing task`: `task-1`
     - `Next recovery`: `task-1 | qa-rework | task-scope | task-1:orchestration-integrity:review-1:multi-owner-ledger-recovery-binding | 2`

     The bound `TASK-RECOVERY task-1` uses category `qa-rework`, repair mode
     `task-scope`, source and repair authority that both name `TRG-21-006` and
     `DEC-163`, the same trigger signature and attempt count, and
     `Consecutive no-progress attempts: 0/2`. It dispatches only review-1 MF-1.
     Review-1 MF-2 follows as an ordinary `task-3` `qa-rework` recovery under
     trigger signature `task-3:qa-rework:review-1` at attempt 1 after task 1
     completes. Task 4 stays unauthorized until a fresh independent re-review
     `review-1b` closes review-1.
  8. **No finding is weakened.** MF-1 and MF-2 keep their exact repair
     boundaries, required evidence, and re-review scope as written in
     `qa/review-1.md`. The validation contract, VAL-001, VAL-010, QA-09, task
     scopes, and required commands are unchanged.
- **Consequences:** A QA Reviewer covering a multi-task window no longer has to
  hide, defer, or artificially split findings to keep the ledger recoverable,
  and the runner copies the Repair Plan owners verbatim. `09_AI_WORKFLOW.md`
  sections 2.3-2.4 and `run-ledger-template.md` now state the ordered-owner
  rule, and `mid-phase-review-template.md` states that Repair Plan row order is
  the repair order. Phase 21 has six more multi-task review windows plus a
  checkpoint and a final review, so this recurs; without the fix each one is a
  potential deadlock. ESC-21-005 is resolved once item 6 lands. Phase 21's
  product scope, task pack, and validation contract are untouched.
- **Refs:** TRG-21-005; TRG-21-006; ESC-21-005; `qa/review-1.md` Repair Plan
  MF-1 and MF-2; `run-ledger.md` REVIEW review-1, CORRECTION, TRIAGE TRG-21-006,
  and STOP; `project-docs/tools/workflow-sequence-lint.sh` REVIEW/TASK-RECOVERY/
  REPAIR/TASK-COMPLETE state machine; `project-docs/tools/drills/run-drills.sh`
  review-schema drills; `09_AI_WORKFLOW.md` sections 2.3, 2.4, and 9;
  `project-docs/templates/run-ledger-template.md`;
  `project-docs/templates/mid-phase-review-template.md`; VAL-001; VAL-010.

### DEC-164: Ordinary schema-v2 TRIAGE recovery preserves exact authority and replaces its pending same-task dispatch. Closes ESC-21-006

- **Date:** 2026-08-24
- **Context:** TRG-21-007 classified task 8's retained Issue/Event/activity/
  Actor failures as ordinary in-scope repair under the accepted task-8
  specification and DEC-162. Its schema-v2 external and ledger TRIAGE records
  therefore bind `Repair authority` exactly as `task-8 specification and
  DEC-162`. `workflow-sequence-lint.sh` then requires the following
  TASK-RECOVERY to both equal that exact value and reference `TRG-21-007` in
  the same field. No string can satisfy both predicates. The byte-identical
  candidate failed `INVALID TASK-RECOVERY does not preserve its governing
  schema v2 TRIAGE REPAIR authority`; adding the TRG ID would instead violate
  the accepted external binding. It simultaneously failed `INVALID overlapping
  TASK-RECOVERY precedes required REPAIR`: TRG-21-007 was accepted while the
  prior task-8 dispatch awaited classification, but direct TRIAGE `REPAIR`
  requires its newly bound TASK-RECOVERY immediately, so no intervening REPAIR
  can close the old dispatch. TRG-21-008 and ESC-21-006 correctly stopped
  the proposed protected verifier change. The ledger itself never accepted
  that second TRIAGE or a STOP and remains valid at `REPAIR_PENDING` after
  TRG-21-007.
- **Decision:**
  1. **Ordinary authority representation.** For a schema-v2 ordinary TRIAGE
     `REPAIR`, TASK-RECOVERY continues to copy `Originating task`, `Trigger
     signature`, `Attempt count`, `Repair category`/`Category`, and `Repair
     authority` byte-for-byte from the authoritative external TRG through the
     ledger TRIAGE. Its `Source` must reference the governing TRG. Its exact
     `Repair authority` is not prefixed, suffixed, or otherwise changed merely
     to repeat that TRG ID.
  2. **Pending same-task handoff.** When schema-v2 direct TRIAGE `REPAIR`
     classifies a blocker candidate returned by an active recovery whose stage
     is `needs-repair`, its immediately bound TASK-RECOVERY replaces that
     returned dispatch only when the task identity is unchanged. The governing
     TRG supplies the new trigger, attempt, category, and authority. A different
     task, a missing TRIAGE bind, any other recovery stage, or any ordinary
     overlapping recovery remains invalid.
  3. **Control authority remains distinct.** A control-conformance
     TASK-RECOVERY still requires `Repair authority` to name the governing
     TRIAGE, and all typed-directive, decision, allowlist, protected-verifier,
     retry, report-repair, and scope checks remain unchanged. Legacy pre-v2
     recovery and STOP/RUN-RESUME authority checks also remain unchanged.
  4. **Authorized verifier correction.** In the schema-v2
     `triage_repair_pending` comparison only,
     `project-docs/tools/workflow-sequence-lint.sh` removes the redundant
     predicate that `event_repair_authority` reference `latest_triage_id`.
     It retains exact equality with `latest_triage_authority`, the Source TRG
     reference, and the other four exact identity comparisons. The separate
     control-conformance check continues to require the authority's TRG ID. In
     the `needs-repair` overlap check only, an exact schema-v2 pending TRIAGE
     recovery for the same task is admitted as the item-2 replacement; all
     other overlap checks remain intact.
  5. **Required drills.** The existing ordinary schema-v2 TRIAGE fixture uses
     an external `Repair authority` that does not contain its TRG ID and must
     validate when TASK-RECOVERY copies it exactly and names the TRG in
     `Source`. A companion negative fixture removes the TRG from `Source` and
     must remain rejected. A third fixture begins with an active same-task
     recovery, records direct TRIAGE `REPAIR`, and proves its bound replacement
     recovery reaches `REPAIR_PENDING`. Every existing drill, including
     unrelated overlap rejection, must still pass.
  6. **Execution route.** This is protected workflow infrastructure, not task-8
     product work and not a control-conformance recovery. A Human-directed
     session outside the phase loop applies items 4-5, updates the canonical
     workflow docs/prompts, runs `project-docs/tools/drills/run-drills.sh` and
     `project-docs/tools/workflow-lint.sh`, proves the current phase ledger is
     `VALID state=REPAIR_PENDING`, proves an exact TRG-21-007 recovery candidate
     appends to a scratch copy, and commits on `phase-21`. No historical ledger,
     triage, task report, test assertion, or escalation record is rewritten.
  7. **Phase-21 recovery.** No RUN-RESUME is appended: the ledger never entered
     STOPPED. After item 6 lands, the Phase Runner validates the ledger and
     appends the pending `TASK-RECOVERY task-8` governed by TRG-21-007 with
     category `in-scope-implementation`, repair mode `task-scope`, trigger
     `task-8:specification-contradiction:task-8-done.md:retained-v1-rebuild-role-contract`,
     attempt 1, `Source` naming `TRG-21-007`, exact `Repair authority` of
     `task-8 specification and DEC-162`, and no-progress `0/2`. It then
     dispatches only the task-8 continuation defined by TRG-21-007. TRG-21-008
     and ESC-21-006 are resolved historical evidence, not ledger events to bind.
  8. **No product or gate change.** Task 8, VAL-025..027, its five exact
     commands, retained assertion subjects/strength, public Issue contract,
     Event schema, activity vocabulary, Actor semantics, and report-lint
     contract are unchanged.
- **Consequences:** Ordinary schema-v2 repair now has one authority source and
  one provenance reference instead of an impossible double encoding, and a
  direct REPAIR verdict can replace the same dispatch whose blocker candidate
  it classified. Exact external binding remains fail-closed,
  control-conformance remains separately guarded, and task 8 can resume without
  rewriting accepted history or weakening any product assertion. ESC-21-006 is
  resolved once item 6 lands.
- **Execution evidence:** The pre-change verifier blob from
  `44370af13c1c5c39e9e040065bc70e15e30b8deb` rejected the exact scratch
  task-8 candidate with both `INVALID TASK-RECOVERY does not preserve its
  governing schema v2 TRIAGE REPAIR authority` and `INVALID overlapping
  TASK-RECOVERY precedes required REPAIR` (terminated `exit_code: 1`). After
  the correction, `project-docs/tools/drills/run-drills.sh` ended `DRILLS PASS
  mechanical drills completed with no misses` (terminated `exit_code: 0`),
  the live ledger ended `VALID state=REPAIR_PENDING` (terminated `exit_code:
  0`), and the byte-identical candidate against a scratch phase directory
  ended `APPEND_CONFIRMED state=REPAIR_PENDING` (terminated `exit_code: 0`).
  `project-docs/tools/workflow-lint.sh` ended `PASS workflow-lint completed
  with no failures`, provider sync check passed, and `git diff --check` was
  silent success (each terminated `exit_code: 0`).
- **Refs:** TRG-21-007; TRG-21-008; ESC-21-006; DEC-162;
  `project-docs/tools/workflow-sequence-lint.sh` schema-v2
  `triage_repair_pending` and control-conformance checks;
  `project-docs/tools/drills/run-drills.sh` ordinary-TRIAGE drills;
  `09_AI_WORKFLOW.md` sections 2.3, 2.4, and 9;
  `project-docs/templates/run-ledger-template.md`;
  `project-docs/templates/triage-log-template.md`.

### DEC-165: Retained-test convergence ownership follows the invalidating change, not the fixture's primary subject. Closes ESC-21-007

- **Date:** 2026-08-25
- **Context:** TRG-21-009 BLOCKed task 9 after eight attempts, reporting that the
  unchanged required package and coverage gates stay red on "explicit
  replacement-bound Task 10 Claim/Run authority fixtures" that task 9's Out Of
  Scope forbids it to touch. The live tree does not support that premise. Every
  remaining red fixture named as Task 10/11/20 property fails on an incidental
  read of a surface **task 9 itself removed or renamed**, not on Claim, Run, or
  dev semantics:
  - `cmd/fid/issue_claim_test.go:459`
    `TestIssueClaimCommandMatchesPolicyShowForForceAcquire` fails at line 464 on
    `runFIDCommand(t, workspace, "policy", "show")`. Its Claim assertions -
    `developer` denied `issue.claim.force_acquire`, `tl` allowed, and the
    force-acquire confirmation string - are live V1 behaviour that still passes.
  - `cmd/fid/run_test.go:179` `TestRunRecordPolicyShowMatchesAllowedRoles` fails
    at line 184 on the same removed command; its `run record` role matrix still
    passes.
  - `cmd/fid/dev_test.go:129` `TestDevCommandsPolicyShowMatchesRowsAndDeniedRoles`
    fails at line 134 on the same removed command; its `dev.start` denial and
    journal-unchanged assertions still pass.
  - `cmd/fid/label_test.go:336`, `cmd/fid/link_test.go:413`,
    `cmd/fid/issue_label_test.go:517`, and `cmd/fid/requirement_test.go:306` are
    the same removed read inside families task 9 already owns outright.
  - `cmd/fid/integration/e2e_v01_closure_helpers_test.go:316` and `:327-329` read
    `maintenance vacuum --by tl`, `stats cost --by role`, `stats tokens --by
    role`, and `policy show` - the fail-closed vacuum, the `--group-by` rename,
    and the public policy removal, all three of them task-9 changes.

  No generic Claim or Run fixture exists in the tree. The only claim sources are
  V1 `cmd/fid/claim.go`, `cmd/fid/issue_claim.go`, `internal/claims/`, and their
  V1 tests; task 10 has not been dispatched and owns no fixture yet. The
  reported boundary is therefore a misattribution, not a protected task-boundary
  or cross-task authority condition.
- **Decision:**
  1. **Ownership rule.** DEC-162 §1.2 already governs and is now stated
     explicitly: retained-test convergence ownership follows the *invalidating
     accepted change*, not the fixture file's primary subject. A fixture whose
     headline subject belongs to a later task does not transfer, defer, or
     bound the obligation of the task whose change broke it. Task 9's public
     `fid policy` removal, `--by` to `--group-by` stats rename, and fail-closed
     vacuum/AC/concern conversions are task-9 changes, so task 9 converges every
     assertion they invalidated wherever it lives.
  2. **Bounded task-9 scope amendment.** Task 9's Files To Create Or Modify
     covers `cmd/fid/` including the `cmd/fid/integration/` package. Its Out Of
     Scope is clarified to exclude generic Claim/Run *implementation, semantics,
     table, CLI, and contract* - it does not exclude converging an incidental
     read of a task-9-removed surface that happens to sit inside a V1 claim,
     run, dev, or closure fixture. No other task-9 scope changes.
  3. **Convergence route for the removed policy inspection.** Apply DEC-162's
     preference order. The unregistered legacy matrix helper is deliberately
     retained for enumerated replacement-bound callers, so each incidental
     `policy show` row read is repointed onto that retained internal helper with
     every role/action assertion byte-identical; this is the "repoint incidental
     reads" case. The enumerated legacy test-caller list is extended to name
     these callers, and the completion report names them for task-20 removal.
     `fid policy` is not re-registered. Deleting a parity block is permitted only
     where `TestPublicPolicyCommandIsAbsent` plus a named surviving denial
     assertion already own the accepted equivalent, and the report names both.
  4. **Untouchable assertions.** The V1 Claim, Run, dev, and closure
     *behavioural* assertions stay byte-identical: `issue.claim.force_acquire`
     denial and `tl` force-acquire, one-active/expiry/wrong-holder claim paths,
     the `run.record` allowed-role matrix, `dev.start|gate|finish|block` denial
     and journal-unchanged checks, and closure journal-byte equality. They keep
     full strength until task 10, 11, or 20 replaces them under its own
     authority.
  5. **Forbidden.** No task-10 or task-11 implementation; no final V1 removal;
     no `fid policy` re-registration; no restored writer for AC verification,
     concerns, or vacuum; no `--by` restoration on a task-9 family; no
     `t.Skip`, build tags, `-run` narrowing, threshold/budget/tolerance change;
     no task-spec, QA-artifact, validation-contract, run-ledger, public-interface,
     Event-schema, or required-command change.
  6. **No gate, VAL, or contract change.** Task 9's five exact commands,
     VAL-028..031, the public interface, the Event schema, the activity
     vocabulary, and the task 10/11/20 specifications are unchanged. Task 9
     closes on its own unchanged gate.
  7. **Recovery route.** The ledger is STOPPED after TRIAGE TRG-21-009. The
     Phase Runner appends `RUN-RESUME` naming the STOP, governing `TRG-21-009`,
     governing `ESC-21-007`, trigger signature
     `task-9:scope-boundary:task-9-done.md:replacement-bound-task10-claim-run-fixtures`,
     attempt 9, and resolution authority `DEC-165`; then one same-task
     `TASK-RECOVERY task-9` with `Source` and `Repair authority`
     `TRG-21-009 / DEC-165`, category `in-scope-implementation`, repair mode
     `task-scope`, the same trigger signature, attempt 9, and no-progress `0/2`.
     It then dispatches only the task-9 Developer continuation. The Tech Lead
     does not edit `run-ledger.md`.
  8. **Triage audit.** TRG-21-009's BLOCK is *not* sustained on its recorded
     category. The Developer and TL-Lite were right to refuse to widen scope on
     their own authority, and refusing to weaken the Claim assertions was
     correct. The recorded protected category was wrong in hindsight for the
     same reason as TRG-21-004: residual failures were collateral of an accepted
     change the escalating task itself made, and the escalation classified them
     by file subject instead of by invalidating change. This is the second
     occurrence; item 1 and the `08_QUALITY_PLAN.md` §3 amendment exist to stop a
     third.
- **Consequences:** Task 9 can close on its unchanged gate without any
  re-slicing, sequencing change, or cross-task authority grant. Task 10 keeps
  every accepted Claim semantic and inherits its V1 fixtures with full assertion
  strength. The generalised ownership rule removes the recurring
  "another task owns this file" deadlock. ESC-21-007 is resolved.
- **Verification evidence:** `cmd/fid/policy_test.go:150` still defines
  `commandPolicyCellAllowed`, and `cmd/fid/root.go` no longer registers
  `newPolicyCommand`, so the retained helper is reachable and the public command
  is absent; `cmd/fid/workspace_fixture_test.go:238`
  `TestTask9FixtureAdapterLeavesJudgmentAndPolicyCommandsUntouched` confirms the
  fixture adapter deliberately does not rewrite `policy show`; `cmd/fid/claim.go`,
  `cmd/fid/issue_claim.go`, `cmd/fid/claim_test.go`, and
  `cmd/fid/issue_claim_test.go` are the only claim surfaces present, all V1.
- **Refs:** TRG-21-009; ESC-21-007; DEC-162;
  `v1-convergence-inventory.md` §§1, 1.2, 3, 5;
  `tasks/task-9-retained-actor-policy-removal.md` Scope and Out Of Scope;
  `tasks/task-10-generic-claims.md`; `tasks/task-9-done.md`;
  `run-ledger.md` TRIAGE TRG-21-009 and STOP; `08_QUALITY_PLAN.md` §3;
  `09_AI_WORKFLOW.md` sections 2.3 and 9; VAL-028; VAL-029; VAL-030; VAL-031.

### DEC-166: Historical task labels do not establish current-phase test ownership. Closes ESC-21-008

- **Date:** 2026-08-26
- **Context:** TRG-21-010 BLOCKed task 9 because two red functions in
  `cmd/fid/task10_contract_test.go` still expect the Phase 20 `--by` transport
  and its old missing-role failures for `issue ac add`, `label add`, and
  `issue label add/remove`. The escalation treated the filename and `TestTask10`
  prefix as Phase 21 task-10 ownership, then treated recurrence of that
  classification as evidence of a protected contradiction. Repository
  provenance disproves both premises:
  - commit `16f4f8d2d2b5870d1502e6713300f5671dcfdf4e` introduced the file under
    `Refs: phase-20/task-10`;
  - commits `62a3e46ebb69b76b93b5a559dc82d9bc9c0ccd3e` and
    `627ca9dea9bdef451260022718b7bfcfcae9186a` added the red request/actor
    assertions under the same Phase 20 task;
  - Phase 21 task 10 is generic Claim acquisition/renewal/release and owns none
    of these inventory leaves; and
  - Phase 21 `public-interface.md` section 5 requires `--actor` plus an
    idempotency key for these leaves, while section 1 makes a missing required
    flag or malformed flag value `invalid-argument` request validation.
  The implementation's rejection of `--by` is therefore correct. The red state
  is an obsolete retained fixture, not a product defect or an accepted-contract
  contradiction.
- **Decision:**
  1. **Provenance before ownership.** A filename, Go test name, task number, or
     prose label is not current-phase ownership evidence. For a retained test,
     establish the introducing phase/task from commit provenance and compare
     the test with the current phase task pack, convergence inventory, and
     accepted contract. Historical task labels remain historical unless the
     current task pack explicitly adopts them.
  2. **Task-9 convergence.** Task 9 owns the exact stale rows because its
     accepted change displaced their `--by` transport. It may modify and rename
     `cmd/fid/task10_contract_test.go` and its `TestTask10*` functions to
     subject-based inventory/request-contract names. This is retained-fixture
     convergence under DEC-162/165, not task-10 implementation or a public
     contract change.
  3. **Preserved subjects and strength.** The migrated rows still prove every
     missing and malformed acting-attribution request fails before mutation,
     JSON failures use the exact current closed `invalid-argument` envelope,
     Human failures remain concise prose naming `--actor` or the accepted Actor
     grammar, and journal length stays unchanged. Selected-project success rows
     use `--actor agent:developer` plus idempotency. Unrelated request, Claim,
     Run, dev, closure, and journal assertions remain byte-identical.
  4. **No legacy restoration.** Do not restore `--by`, old missing-role
     envelopes, or any writer; do not change production behavior merely to make
     the retained fixture green. No generic Claim/Run implementation, final V1
     removal, gate/VAL/command/public-interface/Event-schema change, skip,
     narrowed test selection, or assertion loss is authorized.
  5. **Recurrence semantics.** Repetition of a classification is not evidence
     that the classification is correct and does not create a protected
     category. Recovery-budget exhaustion counts only consecutive no-progress
     attempts for one valid causal trigger. A narrowed failure or newly proven
     provenance is material diagnostic progress; a repeated ownership inference
     from the same filename is not an independent protected condition.
  6. **Recovery route.** The ledger remains STOPPED after TRG-21-010 and
     ESC-21-008. The Phase Runner appends `RUN-RESUME` retaining those IDs,
     trigger signature
     `task-9:scope-boundary:task-9-done.md:task9-converted-leaf-actor-contract-retained-legacy-by-assertions`,
     continuing task 9, next attempt 12, and resolution authority `DEC-166`.
     It then appends one same-task `TASK-RECOVERY task-9` with `Source` and
     `Repair authority` `TRG-21-010 / DEC-166`, category
     `in-scope-implementation`, mode `task-scope`, the same trigger and attempt,
     ordinary task scope, and consecutive no-progress `0/2`. The Developer
     converges the historical fixture, preserves all failed evidence, runs all
     five task-9 commands byte-identically, and returns a truthful report.
- **Consequences:** ESC-21-008 is resolved without changing the accepted
  product interface, task ordering, validation contract, or gate. The workflow
  now distinguishes phase provenance from coincidental task numbering before
  it treats a retained fixture as protected cross-task property.
- **Verification evidence:** A focused reproduction of the two functions ended
  `FAIL` (terminated `exit_code: 1`) solely because the fixture expects
  `authority-missing`/`--by` while the commands return current
  `invalid-argument`/`--actor` failures. `git log --follow` and the three commits
  above bind the file to Phase 20 task 10; Phase 21 task 10's specification is
  generic Claims. No production file is dirty at this decision point.
- **Refs:** TRG-21-010; ESC-21-008; DEC-162; DEC-165;
  `public-interface.md` sections 1, 5, and 7;
  `v1-convergence-inventory.md` sections 1 and 1.2;
  `tasks/task-9-retained-actor-policy-removal.md`;
  `tasks/task-10-generic-claims.md`; `08_QUALITY_PLAN.md` section 3.

### DEC-167: A required gate must exercise the scope its task owns; task 12 owns the tagged integration convergence. Closes ESC-21-009

- **Date:** 2026-08-28
- **Context:** TRG-21-012 BLOCKed task 12 as `recovery-budget-exhausted` after
  attempts 4 and 5 returned NO-PROGRESS, routing its two remaining gate failures
  to task 9 and future task 20. Inspection disproves every element of that
  disposition.
  - **The `internal/dev` routing is wrong.** `internal/dev/start_test.go:225-233`
    requires `issues.Transition(... "in_progress")` to fail with *"want missing
    claim rejection"*. Task 12's own accepted change removed that precondition:
    `internal/issues/transition.go` no longer calls
    `claims.RequireLiveIssueClaimHeld` and retains only
    `claims.MaybeBuildIssueReleaseEffect` at :268, exactly as `task-12-done.md`
    line 15 describes. The routing came from `v1-convergence-inventory.md:54`,
    which assigns **source-package removal** of `internal/dev` to task 20. That
    is not test-convergence ownership, and reading it as such is the file-subject
    inference DEC-165 and DEC-166 already rejected. The test's actual subject -
    `Start()` rejects an invalid workflow edge before append, :255-270 - is
    untouched by convergence.
  - **The `issue ac verify` routing is misframed.** It is not a `--by` to
    `--actor` transport gap task 9 left undone. `cmd/fid/issue_ac.go:109-111`
    makes the leaf unconditionally `capability-unavailable`, which is accepted
    VAL-030 owned by task 9. Converting the flags would still fail. Task 9
    deliberately and correctly excluded verify from its conversion family
    (`e2e_v01_closure_helpers_test.go:131-142`) and asserted the fail-closed
    contract in `task9_contract_test.go:75-79` and
    `TestIssueACVerifyRemainsGuardedLegacyLeaf`. The journey line must become a
    fail-closed rejection assertion, not a success call.
  - **It is a chain, not a line.** The same VAL-030 change also invalidates
    `concern raise` (`e2e_v01_closure_helpers_test.go:263`) and `concern address`
    (:266); both are `capability-unavailable` at `cmd/fid/concern.go:101` and
    :83. Task 9 already converged the same family in the same file for
    `maintenance vacuum` (:389-393) and `policy show` (:409-413), because DEC-165
    named those sites; it stopped at the named list.
  - **Root cause is gate coverage, not ownership.** `cmd/fid/integration/` is
    `//go:build integration`, and `go test ./...` never compiles tagged files.
    Only tasks 12, 18, 21, and 22 carry `go test -tags integration` in `## Exact
    Commands To Run`. DEC-165 gave task 9 ownership of `cmd/fid/integration/`
    without giving it a command that runs the package. Every prior owner saw the
    package red and correctly walked away: `task-8-done.md:71-72` recorded it
    BLOCKED; `task-9-done.md:120-122` ran it voluntarily, recorded FAIL, and
    named preserved task-8 Event-actor and `issue create --by` assertions in
    three journeys; `task-10-done.md:56` failed at `:33` and routed to task 12
    through TRG-21-011. Task 12 is simply the first task obligated to make the
    gate green, so the collateral of tasks 2, 8, 9, 10, and 12 arrives at once.
  - **The budget was consumed by directive, not by capability.** Run-ledger
    :1904 already forbade task 12 from repairing *"Task 9-owned AC --status or
    later Task 20-owned internal/dev fixtures"*. Attempt 4 was report-only;
    attempt 5 *"made no implementation or assertion change"*
    (`task-12-done.md:112`). Both NO-PROGRESS outcomes were structurally
    guaranteed by their own directives.
- **Decision:**
  1. **Verdict audit.** TRG-21-012 is **not sustained**. The Developer's refusal
     to widen scope or weaken assertions on child authority was correct
     behaviour, and its evidence is accurate. Both recorded owner routings are
     wrong, and `recovery-budget-exhausted` does not hold, because neither
     no-progress attempt was permitted to address the failing rows. This is the
     fourth recurrence of the TRG-21-004 pattern; DEC-165 and DEC-166 fixed the
     ownership rule, and this decision fixes the generator behind it.
  2. **`internal/dev` convergence belongs to task 12.** Task 12's file scope is
     amended to name `internal/dev/`. It converges only the stale claim
     precondition expectation its own change invalidated. The `Start()` subject,
     the `*issues.InvalidTransitionError` assertion, the `CurrentState` check,
     and the journal-unchanged check remain byte-identical. Every other
     `internal/dev` assertion stays untouched for task 20.
  3. **Task 12 is the designated convergence owner for `cmd/fid/integration/`.**
     Its file scope is amended to name the `cmd/fid/integration/` package
     explicitly, and it converges every retained V1 row in that package that an
     accepted Phase 21 change invalidated, until
     `go test -tags integration -race ./cmd/fid/integration/...` is green. This
     is a bounded, recorded departure from DEC-162's default owner: the earlier
     invalidating tasks are closed, their required gates never executed this
     package, and reopening tasks 8 and 9 with their closed review windows costs
     more than the audit trail it would preserve. The audit trail is preserved
     instead by item 4.
  4. **Per-row provenance is required.** `task-12-done.md` records, for each
     converged row, the file and line, the invalidating accepted change, and the
     task that made it. A row whose provenance cannot be established as an
     accepted Phase 21 change is not converged; it is reported as a concern.
  5. **No weakening.** Fail-closed rows become explicit rejection-plus-no-journal
     -append assertions using the current typed envelope. Do not restore `--by`,
     `--status`, V1 Claim behaviour, `fid policy`, role authority, or any writer;
     do not change production behaviour to make a retained fixture green. No
     `t.Skip`, build tags, `-run` narrowing, threshold change, or silent
     assertion loss. Deleting an assertion requires a named surviving test that
     already owns the accepted equivalent, named in the report.
  6. **Progress semantics for a chained journey.** Advancing the tagged journey
     to a later step is material diagnostic and causal progress and resets
     consecutive no-progress to `0/2`. A sequential end-to-end fixture reveals
     one blocking step at a time by construction; treating each newly reached
     step as no-progress would exhaust any budget regardless of the work done.
  7. **Gate coverage must match owned scope (standing planning rule).** A task's
     `## Exact Commands To Run` must exercise every file scope the task owns.
     Where a build-tagged, environment-gated, or otherwise excluded suite covers
     files in a task's scope, the Tech Lead names the gate that runs it in that
     task, or the scope is not assignable to that task. A scope amendment that
     adds files without adding the gate that exercises them is incomplete. This
     rule is recorded in `08_QUALITY_PLAN.md` §3.
  8. **Recovery-budget accounting (standing rule).** A no-progress attempt whose
     own recovery directive forbade repairing the only rows causing the failure
     does not count against the two-attempt budget. Before recording a second
     consecutive NO-PROGRESS, the Phase Runner confirms the directive actually
     permitted causal repair. A directive must not encode an owner routing that
     no accepted decision or verified provenance supports.
  9. **Recovery route.** The ledger remains STOPPED after TRG-21-012 and
     ESC-21-009. The Phase Runner appends `RUN-RESUME` retaining those IDs,
     trigger signature
     `task-12:gate-repair:task-12-done.md:tagged-integration-retained-claim-transition-collateral`,
     continuing task 12, next attempt 6, and resolution authority `DEC-167`. It
     then appends one same-task `TASK-RECOVERY task-12` with `Source` and
     `Repair authority` `TRG-21-012 / DEC-167`, category
     `in-scope-implementation`, repair mode `task-scope`, the same trigger and
     attempt, ordinary task scope under the amended specification, and
     consecutive no-progress `0/2`. All six task-12 commands rerun
     byte-identically. No gate, VAL, command, public-interface, Event-schema, or
     accepted product contract changes.
- **Consequences:** ESC-21-009 is resolved with no product-contract change and no
  re-slicing. Task 12 absorbs accumulated tagged-journey convergence that the
  task pack never assigned an executable owner, and the workflow gains the
  planning rule that prevents the same accumulation. Reviewers should expect
  several task-12 continuations, each reaching a later journey step.
- **QA note for the Group 1 checkpoint:**
  `TestV01ClosureJourneyAndRebuildParity` has not been green at any point in
  Phase 21. Its rebuild-parity, journal-immutability, and maintenance-state
  assertions (`e2e_v01_closure_test.go:100-122`) are therefore retained coverage
  that Phase 21 has never actually exercised end to end. This is not a VAL gap -
  VAL-084/085 are separate task-22 journeys - but the checkpoint reviewer should
  treat the first green run as new evidence rather than continuity.
- **Verification evidence:** `internal/issues/transition.go` contains no
  `RequireLiveIssueClaimHeld` call; `cmd/fid/issue_ac.go:110`,
  `cmd/fid/concern.go:83`, and `cmd/fid/concern.go:101` all return
  `writeCapabilityUnavailable`; `cmd/fid/integration/e2e_v01_closure_test.go:1`
  and `e2e_v01_closure_helpers_test.go:1` carry `//go:build integration`; a
  search of `project-docs/docs/phases/phase-21/tasks/` for `tags integration`
  returns it in the required-command block of tasks 12, 18, 21, and 22 only. No
  product file is dirty at this decision point.
- **Refs:** TRG-21-012; ESC-21-009; TRG-21-011; DEC-162; DEC-165; DEC-166;
  `tasks/task-12-execution-bound-transition.md`;
  `tasks/task-9-retained-actor-policy-removal.md`;
  `qa/validation-contract.md` VAL-030 and VAL-038..041;
  `v1-convergence-inventory.md` sections 1, 1.2, and 4;
  `08_QUALITY_PLAN.md` section 3.

## Spec Sync Verdict

Planning-time changes are recorded with their owning decisions. Full positive
verdicts are completed at phase wrap-up.

| Spec | Planning verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending wrap audit | - |
| `01_PRD.md` | touched: see DEC-156 | Phase 21 qualification makes the terminal-success Issue transition and later Envelope completion independently evaluate the same gate-evidence Obligation. |
| `02_TRD.md` | touched: see DEC-154, DEC-156, DEC-157, DEC-158 | Capability availability, shared evaluator, bootstrap, public enums, trust, idempotency, executor boundary, and the unchanged deterministic rebuild target are explicit. |
| `03_TRD_SCHEMA.md` | touched: see DEC-157 | Exact migrations, kebab checks, displaced rule deletion, and role-free bootstrap fixture are accepted for review. |
| `04_TRD_CONFIG.md` | touched: see DEC-157 | Exact public executor bounds agree with the closed input contract. |
| `05_DESIGN_PLAN.md` | touched: see DEC-157 | Exact live NDJSON payload union and Human rendering authority is the public-interface contract. |
| `06_DESIGN_PLAN_CLI.md` | touched: see DEC-157 | Complete leaves, flags, inputs, results, errors, exits, and enum spelling are explicit. |
| `07_DESIGN_PLAN_FLOWS.md` | touched: see DEC-156, DEC-157 | Gate/evidence precede terminal transition; rebuild repeats reads and checks without a second Envelope. |
| `08_QUALITY_PLAN.md` | touched: see DEC-152, DEC-157, DEC-158, DEC-162, DEC-165, DEC-166 | Checkpoint isolation, removal-before-parity, cadence, feasibility, committed-report lint ownership, superseded-contract convergence, invalidating-change ownership, and retained-test phase provenance are explicit. |
| `09_AI_WORKFLOW.md` | touched: see DEC-161, DEC-163, DEC-164, DEC-166 | Schema-v2 identity/binding and QA repair ownership are explicit; retained-test escalation now requires phase provenance before cross-task classification. |
| `10_IMPLEMENTATION_ROADMAP.md` | touched: see DEC-153, DEC-156, DEC-157 | Gate closed; one Heavy Phase 21/two groups; terminal guard and tasks 19-22 ordering fixed. |
| `GLOSSARY.md` | touched: see DEC-153, DEC-156, DEC-157 | Routine, capability availability, public enum, attribution, and terminal-success are unambiguous. |

## Triage-Log Audit

TRG-21-004 correctly refused to weaken retained tests without authority, but its
protected category was wrong in hindsight: the accepted contracts never
conflicted with each other, and the residual failures were retained-test
collateral of an accepted interface change that task 2's file scope omitted.
DEC-162 supplies the bounded convergence authority and the standing rule.

TRG-21-005 and TRG-21-006 were both correct BLOCK outcomes. The condition was a
genuine protected workflow-policy and canonical-verifier question, TL-Lite
refused to invent multi-owner semantics, and the correction in TRG-21-006
supplied the required task-prefixed external identity without rewriting
TRG-21-005. Hindsight verdict: correct, and cheaply reached. The upstream defect
belongs to the verifier, which accepted a REVIEW append that no later event
could recover; DEC-163 fixes that and records the accepted representation.

TRG-21-001 correctly identified the bootstrap/V1 sequencing contradiction;
DEC-160 supplies the bounded task-2 repair. TRG-21-002 is malformed historical
evidence. TRG-21-003 is the valid governing BLOCK and DEC-161 supplies its
exact RUN-RESUME boundary. ESC-21-001..003 are resolved without rewriting any
TRG or ledger entry. Hindsight verdict: all three BLOCK outcomes were correct
under the authority available when written; TRG-21-002 should have used a
task-prefixed plain scalar initially, which would have avoided TRG-21-003.

TRG-21-007 correctly classified task 8's retained-contract collision as
ordinary in-scope convergence under DEC-162. TRG-21-008 correctly escalated the
protected verifier definition change after the exact TASK-RECOVERY candidate
proved the accepted five-field byte binding and the additional authority TRG
reference mutually exclusive. DEC-164 keeps the former, locates provenance in
`Source`, removes the redundant schema-v2 ordinary-repair predicate, and lets
the bound task-8 recovery replace the pending task-8 dispatch that produced the
classified blocker candidate.
Hindsight verdict: both outcomes were correct under the authority available;
no historical record needs correction.

## Workflow Notes

To be completed at phase wrap-up.

### Deep Retrospective Recommendation

To be completed at phase wrap-up.
