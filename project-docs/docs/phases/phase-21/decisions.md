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

### DEC-168: A child QA session's lifetime is not bounded by its artifact's command budget; the Group 1 checkpoint reruns unchanged. Closes ESC-21-010

- **Date:** 2026-08-29
- **Context:** TRG-21-014 BLOCKed the phase as `unavailable external capability /
  bounded QA-session liveness` after three fresh independent QA Lead checkpoint
  sessions each returned `STATUS: INTERRUPTED` with no artifact, commit, verdict,
  or command session. The condition is real - there is no
  `qa/group-1-checkpoint-review.md` - but the recorded cause is wrong.
  - **The budget was applied to the wrong object.** `qa-plan.md` §11 declares
    "each FEAS probe at most 120s wall clock; the whole FEAS set at most 900s;
    the tagged real-process journey gate ... at most 300s". Those bound *command
    execution*. Orchestration applied the 900s FEAS-set figure as the lifetime of
    the entire QA Lead child session, which also has to build a probe harness,
    audit twelve committed completion reports and their report-lint evidence for
    VAL-042, classify and prove isolation of the eleven-item temporary V1
    inventory for VAL-043, and write one artifact. Each attempt was killed at the
    same wall-clock point, before the artifact.
  - **Each retry restarted from zero.** The probes must be session-captured, so
    every fresh attempt re-created scratch, re-authored `child/main.go`, and
    re-ran the early probes with the same time left. Three attempts against an
    unchanged cut-off could only produce the same result.
  - **The capability is available.** Preserved scratch `/tmp/tmp.zQAE0hBrd2`
    contains the built `child-bin` plus materialized outputs for essentially the
    whole probe set: `out.txt`/`err.txt` (FEAS-02), `badutf8.bin` (FEAS-03),
    `flood.out`/`flood.err` (FEAS-07), `canary.txt` (FEAS-09),
    `outside/target.txt` with the escaping symlink (FEAS-10),
    `hash-before.txt`/`hash-after.txt` (FEAS-11),
    `store/objects/55b62c00...` (FEAS-12), a two-commit `repo-tree/` (FEAS-14),
    `restored-object.bak` (FEAS-15), `ar-one.zip`/`ar-two.zip` (FEAS-16),
    `ar.zip`/`ar-warning.zip` (FEAS-17), and `absence/` (FEAS-18). A fresh
    diagnostic run of `go test -tags integration -race ./cmd/fid/integration/...`
    outside any checkpoint session likewise completed. Nothing here was blocked
    by the environment; the session envelope ran out.
  - **None of that material is checkpoint evidence.** The accepted contract
    requires probes executed and observed *in the checkpoint session* by the
    independent reviewer. Preserved scratch is a diagnostic artefact of a killed
    session, and the diagnostic tagged run was not session-captured. Both are
    inadmissible for FEAS or for any PASS. They establish availability for this
    decision only.
  - **The QA Lead behaved correctly.** Returning no artifact rather than partial,
    transcribed, or wrapper-completion evidence is exactly what the contract
    demands. `09_AI_WORKFLOW.md` §9 also makes interrupted sessions `REPAIR` by
    default, and DEC-167 item 8 counts only attempts whose directive permitted a
    causal repair; none of the three did.
  - **The identity churn is a separate defect.** TRG-21-013 recorded
    `Originating task: N/A` with a `phase-21:` trigger signature, which schema v2
    cannot bind because it requires the trigger signature to begin with the
    originating task. TRG-21-014 exists only to re-label the same block onto
    `task-13`. This is the third such correction cycle in Phase 21 after
    TRG-21-002/003 and TRG-21-005/006.
- **Decision:**
  1. **Verdict audit.** TRG-21-013 and TRG-21-014 are **not sustained on their
     recorded category**. The blocking condition was real and correctly refused,
     but it is an orchestration defect in how a bounded command budget was applied
     to an agent session, not an unavailable external capability. The corrected
     category is orchestration session-liveness.
  2. **Binding reading of the declared budgets, which do not change.** 120s bounds
     each individual FEAS probe. 900s bounds the FEAS set, meaning the aggregate
     wall-clock execution of the FEAS-01..FEAS-19 probe commands. 300s
     independently bounds the tagged integration command
     `go test -tags integration -race ./cmd/fid/integration/...` (FEAS-19), and
     that 300s is not carved out of anything else. **No budget bounds the QA Lead
     child session's total lifetime**, which additionally covers harness setup,
     the VAL-042 and VAL-043 audits, and artifact authoring. No edit to
     `qa-plan.md` §11 or `group-1-checkpoint-contract.md` is authorized or
     required; this item records the reading, and the artifacts stand unchanged.
  3. **Nothing is weakened, split, or substituted.** VAL-042, VAL-043,
     FEAS-01..FEAS-19, the three separately decided checkpoint parts, the
     fresh-session evidence rule, observed-termination rules, the exact commands,
     the accepted statuses, the per-probe and per-set budgets, reviewer
     independence, and the single-artifact `PASS`/`REPAIR`/`BLOCK` verdict
     structure are all unchanged. The checkpoint is **not** split across sessions.
     Partial, conditional, transcribed, or undocumented status still never
     dispatches Group 2.
  4. **Fresh scratch only.** `/tmp/tmp.YC4EUO3gMN` and `/tmp/tmp.zQAE0hBrd2` are
     preserved as evidence of the failed attempts and must not be reused, read
     into, or copied from. The attempt allocates new scratch with `mktemp -d`,
     rebuilds `child/main.go` and `child-bin`, and reruns every probe fresh. The
     diagnostic tagged run is not FEAS-19; FEAS-19 runs again in session.
  5. **Orchestration obligations for the attempt.** The Phase Runner dispatches
     one fresh independent QA Lead checkpoint session and **waits for that child to
     terminate on its own**. It must not interrupt, cancel, or time out the child
     at 900 seconds or at any other whole-session figure; the only enforced
     budgets are the per-probe, per-set, and tagged-command budgets of item 2,
     which the QA Lead enforces on its own commands and records with observed
     termination. The QA Lead **writes `qa/group-1-checkpoint-review.md`
     incrementally during the attempt**, appending each probe and audit result as
     it is observed, so that an interruption leaves inspectable partial evidence.
     A partially written artifact is still never a `PASS`.
  6. **Authorization.** One fresh independent checkpoint attempt is authorized as
     attempt 4 under the retained trigger signature. ESC-21-010 is closed. Task 13
     Developer dispatch remains prohibited until the committed checkpoint artifact
     records `PASS` and the Phase Runner has independently verified that commit.
  7. **Exact one-off resume directive.** Schema v2 has no representation for a
     phase-scoped STOP recovery, so this decision binds the existing
     TRG-21-014/ESC-21-010 STOP to `task-13` as a one-off orchestration wrapper.
     The Phase Runner appends, with no other event between them:

     `RUN-RESUME`
     - `Sequence contract: v2`
     - `Recovery directive schema: v2`
     - `Prior run/session or STOP: 2026-08-28 23:21 STOP / TRG-21-014 / ESC-21-010`
     - `Governing TRG: TRG-21-014`
     - `Governing ESC: ESC-21-010`
     - `Trigger signature: task-13:group-1-checkpoint:qa-lead-feasibility-session-interrupted`
     - `Attempt count: 4`
     - `Checkpoint:` runner-verified phase HEAD with matched subject, absent
       `qa/group-1-checkpoint-review.md`, and the preceding STOP
     - `Resolution evidence:` DEC-168 and its Resolution Note are committed, the
       category is corrected to orchestration session-liveness, ESC-21-010 is
       closed, and one fresh attempt is authorized with no contract change
     - `Resolution authority: DEC-168`
     - `Continuing task: task-13`
     - `Next recovery: task-13 | session-continuation | task-scope | task-13:group-1-checkpoint:qa-lead-feasibility-session-interrupted | 4`

     `TASK-RECOVERY task-13`
     - `Source: TRG-21-014 / DEC-168`
     - `Category: session-continuation`
     - `Repair mode: task-scope`
     - `Repair authority: TRG-21-014 / DEC-168`
     - `Triage disposition: N/A`
     - `Trigger signature: task-13:group-1-checkpoint:qa-lead-feasibility-session-interrupted`
     - `Attempt count: 4`
     - `Control contract: N/A`
     - `Control mismatch: N/A`
     - `Allowed files: project-docs/docs/phases/phase-21/qa/group-1-checkpoint-review.md; after a verified checkpoint PASS, ordinary task scope under project-docs/docs/phases/phase-21/tasks/task-13-executor-policy.md`
     - `Forbidden changes:` do not interrupt or time-bound the QA Lead child
       session; do not dispatch the task-13 Developer before a verified committed
       checkpoint `PASS`; do not reuse `/tmp/tmp.YC4EUO3gMN` or
       `/tmp/tmp.zQAE0hBrd2` or any prior scratch; do not accept transcribed,
       partial, promised, or wrapper-completion evidence; do not change
       `qa/group-1-checkpoint-contract.md`, `qa/qa-plan.md`,
       `qa/validation-contract.md`, VAL-042, VAL-043, FEAS-01..FEAS-19, any probe
       command, budget, tolerance, threshold, accepted status, or task
       specification; no product, test, migration, spec, decision, or prior-ledger
       edit under this recovery
     - `Originating report: N/A`
     - `Exact verification: N/A`
     - `Control repair base: N/A`
     - `Consecutive no-progress attempts: 0/2`
     - `Checkpoint:` as above, plus the `RUN-RESUME` attempt 4 as the immediately
       preceding ledger event
     - `Phase worktree:` current branch and dirty-path statement
     - `Temporary repair worktree: N/A for task-scope`
     - `Hidden-index flags: none; assume-unchanged and skip-worktree are forbidden`
     - `Resolution evidence: PENDING -` run one fresh independent checkpoint
       session to its own termination, commit `qa/group-1-checkpoint-review.md`
       with separate VAL-042 and VAL-043 decisions and FEAS-01..FEAS-19 observed
       termination, then dispatch task 13 only on a verified `PASS`
     - `Next action:` dispatch only a fresh independent QA Lead checkpoint session
       under the unchanged contract and QA plan §§10-11; wait for it to finish;
       verify the committed artifact and its verdict
     - `Human attention required: no`
  8. **The wrapper is orchestration, not task-13 work.** The open recovery names
     `task-13` solely because schema v2 requires a task-prefixed identity. The
     Phase Runner dispatches the **QA Lead first**. Only after the committed
     artifact records `PASS` and the Runner has verified the commit may it dispatch
     the task-13 Developer under the unchanged `task-13-executor-policy.md` and its
     seven unchanged required commands. The checkpoint outcome is **not** recorded
     as a `REVIEW` event: ledger review IDs are `review-N`/`infra-review-N`, and
     `review-5` belongs to the tasks 13-15 window. The checkpoint commit is carried
     in the `REPAIR` event instead.
  9. **Closure sequence.** On a verified checkpoint `PASS` followed by a completed
     task 13, the Phase Runner appends `REPAIR task-13` with the same trigger
     signature and `Attempt count: 4`, `Starting checkpoint` naming the resumed
     STOP state, `Ending checkpoint` naming the verified checkpoint-artifact commit
     **and** the verified task-13 implementation and report commits with the seven
     command results and the two-argument report-lint result, `Control repair
     base: N/A`, `Control repair head: N/A`, `Control scope check: N/A for
     task-scope`, `Phase-branch integration: N/A for task-scope`,
     `Outcome: RESOLVED`, resolution evidence, and
     `Next checkpoint: VERIFY/REPORT/REVIEW`. It then appends `TASK task-13
     COMPLETE` with `Report:
     project-docs/docs/phases/phase-21/tasks/task-13-done.md`, `report-lint: PASS -
     OUTCOME COMPLETE`, the verified commits, and the QA-plan-commands line. That
     closes the wrapper and returns the run to the normal loop at the tasks 13-15
     window.
  10. **Branches, so the Runner does not improvise.** (a) If the checkpoint
      returns `REPAIR` or `BLOCK`, append `REPAIR task-13` with
      `Outcome: MATERIAL-PROGRESS` recording the delivered verdict, route the QA
      outcome through `workflow-route.sh`, and expect `TRIAGE` and `STOP`: the
      findings belong to Group 1 tasks 1-12, and a continuing recovery cannot
      change task identity while this wrapper is open, so ownership routing and
      the checkpoint repeat are a Tech Lead decision. (b) If the child is
      interrupted again, append `REPAIR task-13` with `Outcome: NO-PROGRESS`; the
      next `TASK-RECOVERY` is attempt 5 with `1/2`, same task, trigger, and repair
      mode, and a second consecutive `NO-PROGRESS` forces `TRIAGE`. Per DEC-167
      item 8, an interruption caused by orchestration failing to wait is not a
      counting no-progress attempt. (c) If the checkpoint passes but task 13 needs
      ordinary Developer rework, append `REPAIR task-13`
      `Outcome: MATERIAL-PROGRESS` and continue at attempt 5 with the same trigger
      and repair mode under category `in-scope-implementation`; normal task-13
      recovery rules apply from that point.
  11. **Deferred to workflow maintenance, not this phase.** A general schema-v2
      identity for phase-scoped triggers, and a ledger event for an interrupted
      child session, are real gaps - three identity-correction cycles and three
      unrecorded checkpoint attempts prove it - but designing them is workflow
      infrastructure outside the Phase 21 recovery loop. This decision solves the
      instance with the one-off wrapper above and changes no verifier, template,
      or workflow definition.
- **Consequences:** ESC-21-010 is closed with no change to any accepted contract,
  assertion, command, budget, or task specification, and with no split of the
  checkpoint. The next attempt is the fourth attempt at the same unchanged
  checkpoint; its cost is a full independent session, and the Runner is expected
  to wait for it. If it passes, Group 2 opens. The workflow carries one recorded
  one-off wrapper whose only purpose is to make a phase-scoped resume expressible.
- **Verification evidence:** `qa-plan.md:686-689` states the three budgets and
  attaches them to "the checkpoint session"; `qa-plan.md:657-660` requires probes
  executed in the checkpoint session and rejects transcribed evidence;
  `group-1-checkpoint-contract.md:72-80` defines the three verdicts;
  `workflow-sequence-lint.sh:577` accepts `session-continuation` as a recovery
  category; `:903-904` requires the schema-v2 trigger signature to begin with the
  originating task, which is why TRG-21-013 could not bind; `:413` limits REVIEW
  IDs to `review-N`/`infra-review-N`; `:657-659` defines the exact `Next recovery`
  tuple; `:1023-1034` and `:1043-1053` define the `REPAIR`/`TASK-COMPLETE` closure
  used above. `run-ledger.md:2321-2339` holds the governing TRIAGE and STOP, and
  contains no event for the three interrupted attempts.
- **Refs:** TRG-21-013; TRG-21-014; ESC-21-010; DEC-157; DEC-158; DEC-167;
  `qa/group-1-checkpoint-contract.md`; `qa/qa-plan.md` §§10-11;
  `qa/validation-contract.md` VAL-042 and VAL-043;
  `tasks/task-13-executor-policy.md`; `09_AI_WORKFLOW.md` §§2.4 and 9.

### DEC-169: Task 4 owns the Activation replay ordering; task 9 owns the whole section-7 exit contract. Closes ESC-21-011

- **Date:** 2026-08-29
- **Context:** The fresh Group 1 checkpoint ran to completion under DEC-168 and
  returned `VAL-042 PASS`, `VAL-043 PASS`, `FEAS-01..FEAS-19 PASS`, and an overall
  `REPAIR` with two bounded Must-fix findings. TRG-21-015 then stopped the run
  exactly as DEC-168 item 10(a) instructs, because the open task-13 wrapper cannot
  change owner identity. **The stop is correct and by design**; this decision
  supplies the owner routing it asks for.
  - **MF-1 is real and is task 4's.** `internal/workflow/activation.go:176` calls
    `assertActivationAvailable` before `journal.CommitIdempotent` at `:213`, so the
    `ActivationConflictError` at `:125-132` preempts an exact same-key,
    same-digest replay that the idempotency layer would have served. The same
    ordering exists in the component path (`:305-330`) and in
    `transitionActivation`, where `build(conn, activation)` evaluates preconditions
    at `:471` before the commit at `:502`. This contradicts `public-interface.md`
    §6, which requires a same-key/same-digest replay to return the original Event
    and domain fields with only `replayed:true` and no new Event. The generic
    layer is fine - the checkpoint separately observed `claim.acquire` replaying
    correctly - so the defect is the Activation caller path, which task 4 authored
    and owns in its file scope. VAL-005 is task 2's journal-layer observable and
    stays green; no Group 1 acceptance criterion asserted caller-path replay, which
    is why twelve tasks passed over it.
  - **One evidence detail must be pinned down during repair.** The review recorded
    that the second call emitted no stdout or stderr, but
    `cmd/fid/activation.go:273-313` maps `ActivationConflictError` to an
    `activation-conflict` failure envelope on stderr. Either the second call did
    emit that envelope and the review's silence detail is imprecise, or a second
    output-channel defect exists. The red test records both streams verbatim; if
    the process really is silent, that silence is part of the same repair.
  - **MF-2 is real but is not a task-9 value choice.** `cmd/fid/main.go:5-9` holds
    the only `os.Exit` in the product and always exits 1. `internal/cli` has no
    exit concept at all, and 28 `cli.WriteFailure` call sites across 16 files
    return untyped errors afterwards. The entire section 7 class table - `2`
    invalid input, `3` not found, `4` conflict/precondition, `5`
    authority/capability, `6` integrity, `7` internal, `8` observed gate
    non-success, plus unknown-command `2` - is unimplemented. Task 9's four
    fail-closed leaves emit exactly the right envelope; they are simply the first
    real-process observation of a phase-wide gap.
  - **Why nobody saw it.** No Group 1 required command runs a real built binary,
    and `cmd/fid/integration/` contains no exit-status assertion whatsoever. This is
    the fifth occurrence of the DEC-167 item 7 pattern: an accepted contract with
    no gate that can observe it. It is already scheduled to bite again - QA-69
    expects `fid export import` unknown-command exit 2 at task 19, and VAL-072
    verifies every leaf's exit contract at task 22.
  - **The fix is small.** `internal/cli/codes.go:10-19` and `:164+` already attach a
    closed `Category` to every code, so the mapping is a seven-entry
    category-to-exit function plus a typed error, not a per-code rewrite.
- **Decision:**
  1. **Verdict audit.** TRG-21-015 is **sustained**. The Phase Runner and TL-Lite
     followed DEC-168 item 10(a) precisely, and the checkpoint artifact is high
     quality: it separated the three decided parts, kept VAL-042/VAL-043 passing on
     their own evidence, and classified both findings as decision-free and bounded
     with no blocker candidate. The only correction is scope: MF-2 named the right
     owning task for the wrong-sized defect.
  2. **MF-1 owner: task 4.** Repair boundary: no lifecycle precondition may preempt
     an exact same-key, same-digest idempotent replay anywhere in the Activation
     family - `prepare`, `component`, and the `commit`/`abort`/`finalize`
     transitions. A replay returns the original Event, entity, and domain fields
     with only `replayed:true` and appends no Event. New-key lifecycle conflicts and
     same-key/different-digest conflicts still fail exactly as they do now, with
     byte-identical envelopes. Every error path emits its typed failure envelope on
     the accepted channel; a silent nonzero exit is a defect. No Event schema,
     persistence, public-interface, result-shape, or command-grammar change.
  3. **MF-2 owner: task 9, widened to the whole contract.** Task 9 implements the
     complete section 7 mapping, not the single `capability-unavailable` value:
     a category-to-exit function in `internal/cli` covering `invalid-input`→2,
     `not-found`→3, `conflict`→4, `precondition`→4, `authority`→5, `integrity`→6,
     `internal`→7, and `gate-unsuccessful`→8; a typed error that carries the
     failure's code/category out of the leaf; `cmd/fid/main.go` wired to it;
     Cobra unknown-command, unknown-flag, and usage errors mapped to 2; and every
     existing `cli.WriteFailure` call site converted to return it. Human-mode
     failures take the same exit class as JSON. Streaming `gate run` terminal-frame
     exit semantics are untouched and stay with tasks 14-15. **Nothing about any
     envelope changes**: no JSON field, code, category, message, subject, details,
     channel, command grammar, capability availability, or writer absence moves.
     Success stays exit 0.
  4. **Why task 9 owns the sweep.** This is a bounded, recorded departure from
     DEC-162's default owner, on the DEC-167 precedent: task 9 owns `internal/cli/`
     and `cmd/fid/` in its file scope, its VAL-030 observable is the reproduced
     failure, and it is the first task obligated to make the observable green. The
     alternative - a one-code special case now and sixteen files of conversion
     discovered at task 19 and task 22 - is the accumulation this project has
     already paid for four times. The conversion is mechanical and envelope
     preserving, and no existing test asserts a process exit, so the blast radius
     is the contract itself.
  5. **Gate coverage (DEC-167 item 7 compliance).** Both tasks are amended to name
     the `cmd/fid/integration/` package in file scope and to carry
     `go test -tags integration -race ./cmd/fid/integration/...` in `## Exact
     Commands To Run`. Task 4's real-process regression proves the two-call
     Activation result and zero Event delta; task 9's proves the four fail-closed
     leaves at exit 5, one representative each of exit 2, 3, and 4, and
     unknown-command exit 2. Without those gate additions neither scope is
     assignable.
  6. **No contract change.** VAL-030 with its accepted QA-29 exit-5 evidence and
     VAL-072 already own these observables. No VAL, QA plan, checkpoint contract,
     public interface, Event schema, migration, budget, threshold, or accepted
     status changes. Standing rule for the rest of the phase: every leaf failure
     returns the typed exit error, and any leaf added by tasks 13-22 uses the same
     mapping rather than reintroducing an untyped error.
  7. **Order and closure.** Task 4 first, then task 9. Each returns its own
     committed report with fresh two-argument report-lint. Then one fresh
     independent QA Lead reviewer repeats the **complete, unsplit** Group 1
     checkpoint, all FEAS probes fresh in newly allocated scratch, and appends the
     second complete review to `qa/group-1-checkpoint-review.md`, preserving the
     first verdict verbatim. The task-13 Developer is dispatched only on that
     rerun's verified `PASS`.
  8. **Exact ledger route.** The task-13 wrapper is open at needs-continuation and
     only a `RUN-RESUME` clears it, while a resume is bound to its governing
     TRIAGE's task-prefixed trigger signature. TRG-21-015 is `task-13`-prefixed, so
     reaching task 4 requires an identity-only successor triage, exactly as
     TRG-21-014 succeeded TRG-21-013. The Phase Runner appends, in this order:
     1. `REVIEW infra-review-1` - `Path:
        project-docs/docs/phases/phase-21/qa/group-1-checkpoint-review.md`,
        `Tasks covered: task-1 through task-12 Group 1 independent checkpoint`,
        `Outcome: REPAIR`, `Owning task: task-4, task-9`, `Re-review of: none`,
        `Must-fix: 2 / Improvement: 0 / blocker candidates: 0`. This records the
        already-committed artifact under the accepted DEC-163 multi-owner
        representation and mechanically forbids a checkpoint rerun before both
        owners complete. It is appended before the successor triage because a
        `TRIAGE BLOCK` must be immediately followed by `STOP`, and because owner
        state is only captured for recoveries appended after the review. This
        supersedes DEC-168 item 8's no-REVIEW instruction; `infra-review-N` is a
        non-window review ID and does not consume `review-5`.
     2. `TRIAGE TRG-21-016` - identity-only successor of TRG-21-015, allocated by
        TL-Lite in the triage log under this decision, retaining ESC-21-011 rather
        than allocating a new escalation: `Originating task: task-4`,
        `Trigger signature:
        task-4:checkpoint-repair:group-1-checkpoint-review.md:mf-1-activation-idempotent-replay`,
        `Attempt count: 1`, `Repair category: N/A`, `Repair authority: N/A`,
        `Verdict: BLOCK (escalated as ESC-21-011)`. Its triage-log section must
        contain `- **Verdict:** BLOCK` and reference `ESC-21-011`.
     3. `STOP` - required immediately after the TRIAGE BLOCK, naming TRG-21-016 and
        ESC-21-011.
     4. `RUN-RESUME` - `Prior run/session or STOP: <that STOP> / TRG-21-016 /
        ESC-21-011`, `Governing TRG: TRG-21-016`, `Governing ESC: ESC-21-011`,
        `Trigger signature:
        task-4:checkpoint-repair:group-1-checkpoint-review.md:mf-1-activation-idempotent-replay`,
        `Attempt count: 2`, `Resolution authority: DEC-169`,
        `Continuing task: task-4`, `Next recovery: task-4 | qa-rework |
        task-scope |
        task-4:checkpoint-repair:group-1-checkpoint-review.md:mf-1-activation-idempotent-replay
        | 2`.
     5. `TASK-RECOVERY task-4` - `Source: TRG-21-016 / DEC-169`, `Repair authority:
        TRG-21-016 / DEC-169`, `Category: qa-rework`, `Repair mode: task-scope`,
        the same trigger signature, `Attempt count: 2`, `Consecutive no-progress
        attempts: 0/2`, allowed files the amended task-4 scope, forbidden changes
        per item 2. Then `REPAIR task-4` and, on `Outcome: RESOLVED`,
        `TASK task-4 COMPLETE` with `Report:
        project-docs/docs/phases/phase-21/tasks/task-4-done.md` and `report-lint:
        PASS - OUTCOME COMPLETE`.
     6. `TASK-RECOVERY task-9` - authorized by the open review, no further resume:
        `Source: group-1-checkpoint-review.md MF-2 / DEC-169`, `Repair authority:
        DEC-169`, `Category: qa-rework`, `Repair mode: task-scope`,
        `Trigger signature:
        task-9:checkpoint-repair:group-1-checkpoint-review.md:mf-2-failure-exit-class`,
        `Attempt count: 1`, `Consecutive no-progress attempts: 0/2`. Then
        `REPAIR task-9` and `TASK task-9 COMPLETE` the same way.
     7. `REVIEW infra-review-1b` - same `Path`, `Outcome: PASS`, `Owning task:
        none`, `Re-review of: infra-review-1`, `Must-fix: 0 / Improvement: 0 /
        blocker candidates: 0`, appended only after the rerun checkpoint's own
        `PASS` and only after both owners are complete.
     8. The ledger is then `ADVANCEABLE` and task 13 dispatches normally. The
        task-13 wrapper is dissolved by the `RUN-RESUME` at step 4 and is not
        reopened; task 13 completes as an ordinary task.
  9. **Branches.** A failed task-4 or task-9 repair is ordinary same-task recovery
     under its own trigger signature, with the two-attempt no-progress budget and
     DEC-167 item 8 accounting. A `REPAIR` verdict on the checkpoint rerun returns
     to this decision's owner-routing route with a fresh review ID; a `BLOCK`
     verdict escalates. The Phase Runner still never improvises cross-task work
     inside an open recovery.
  10. **Deferred, again.** This is the second instance of the DEC-168 item 11 gap:
      schema v2 can express neither a phase-scoped trigger nor a STOP whose repair
      belongs to a different task than the stopped one, so both instances have cost
      an identity-only successor triage. The general fix stays with workflow
      maintenance, outside this phase; this decision does not change any verifier,
      template, or workflow definition.
- **Consequences:** ESC-21-011 closes with no contract change and no re-slicing.
  Phase 21 gains the exit-class implementation it had specified but never built,
  thirteen tasks before the conformance corpus would have found it, and the
  Activation family gains the replay behaviour the interface always required.
  Reviewers should expect two bounded repairs, then a complete checkpoint rerun
  whose FEAS set runs fresh again.
- **Verification evidence:** `cmd/fid/main.go:5-9` is the only product `os.Exit`
  and is unconditional `os.Exit(1)`; a repository search for `os.Exit`,
  `ExitStatus`, or an exit mapping in `internal/cli` returns nothing outside test
  `TestMain` helpers; `cli.WriteFailure` has 28 call sites across 16 files;
  `internal/cli/codes.go:10-19` defines the closed `Category` set the mapping keys
  on; `cmd/fid/integration/` contains no `ExitCode` or `ProcessState` assertion;
  `internal/workflow/activation.go:176` precedes `:213`; `public-interface.md:402`
  requires exit 5 for the four fail-closed leaves and `:441-445` defines the class
  table; `qa-plan.md:428` (QA-29) and `:496` (QA-69) already expect exit 5 and exit
  2 respectively.
- **Refs:** TRG-21-015; ESC-21-011; DEC-168 items 8 and 10; DEC-167 items 3, 7,
  and 8; DEC-163; DEC-162; `qa/group-1-checkpoint-review.md` MF-1 and MF-2;
  `qa/group-1-checkpoint-contract.md`; `qa/qa-plan.md` §§10-11 and QA-29;
  `qa/validation-contract.md` VAL-005, VAL-030, VAL-072;
  `public-interface.md` §§6-7; `tasks/task-4-activation-trust.md`;
  `tasks/task-9-retained-actor-policy-removal.md`.

---

### DEC-170: An externally interrupted attempt never consumes the no-progress budget; task-15 Window 5 repair resumes at attempt 6. Closes ESC-21-012

- **Date:** 2026-08-30
- **Context:** Review 5 returned `REPAIR` with two bounded task-15 Must-fixes.
  Five runner-mediated continuations followed; TRG-21-017 and its identity
  successor TRG-21-018 triaged attempts 4 and 5 as consecutive no-progress and
  escalated `recovery-budget-exhausted` as ESC-21-012. Four facts decide it.
  - **Not one of the five attempts returned a typed result.** The ledger records
    the same ending for every attempt: "Developer continuation **was shut down**"
    while "in running state" - `run-ledger.md:2730` (attempt 1),
    `:2773` (~30 min), `:2816` (~40 min), `:2859` (several hours), and `:2902`
    (~20 min). No child self-terminated, produced a report, or reported command
    termination. A `NO-PROGRESS` verdict over that population measures the
    orchestrator's shutdown policy, not the Developer's repair capability. This
    is **exactly the DEC-168 defect on the Developer side**: a whole-session
    lifetime cap applied where only per-command budgets exist.
  - **The durable tree contradicts exhaustion.** MF-15-01's causal fix is
    present: `internal/evidence/execution.go:339` now routes *every* observed
    process outcome through one `journal.CommitIdempotent` path, `:120` replays
    the stored outcome, and `internal/journal/project_effects.go` rejects
    non-insert Receipt effects. MF-15-02's four negative cases are present:
    `cmd/fid/gate_run_test.go:141` (Attempt-append failure, no spawn), `:170`
    (frame-only handled `fid-error`), `:210` (non-success replay matrix), plus
    `assertCompleteReceiptBinding:328` and
    `assertReceiptMutationRejectedAndImmutable:344`. The modified
    `task-15-done.md` records all five required commands PASS with concrete
    durations. The repair is at verification-and-closure, not stalled.
  - **One honest verification pass cannot fit in the observed session window.**
    Required command 1 alone terminated at 219-235 s in this task's own records,
    and the repository coverage command runs the whole tree under `-race`; a
    single byte-identical pass is roughly 15-25 minutes of pure command
    wall-clock. Attempts 2, 3, and 5 were shut down at ~30, ~40, and ~20
    minutes. No attempt could both author and verify. Because nothing was ever
    committed, each shutdown discarded the *evidence* of progress and
    re-presented an unchanged tree to the next attempt, which is how genuine
    progress kept reading as no-progress.
  - **MF-15-02's provenance element is unexecutable as written.** `79e3a93` and
    `8491cde` are ancestors of `09a2a7b`. Amending their messages rewrites the
    review-5 commit itself and every hash cited as phase HEAD or control repair
    base by TRG-21-015 through TRG-21-018, the two STOPs, ESC-21-011,
    ESC-21-012, and DEC-168/DEC-169. `task-15-done.md:79` already asserts the
    rewrite happened; the branch says otherwise. An instruction that can only be
    obeyed by destroying the phase's immutable evidence chain is a defective
    repair route, not work the Developer failed to do.
- **Decision:**
  1. **Verdict audit.** TRG-21-017 and TRG-21-018 are **not sustained on their
     recorded category**. The behaviour around them was right: the Runner
     preserved both findings, kept task-15 ownership, refused to advance, and
     TRG-21-018's identity correction to the last actual TASK-RECOVERY was
     correct. Only the classification is wrong. The corrected category is
     **orchestration session-liveness**, the same correction DEC-168 item 1
     made. `recovery-budget-exhausted` does not hold, so ESC-21-012's protected
     premise fails and no protected authority is required to continue.
  2. **Standing accounting rule.** An attempt whose child session is
     interrupted, cancelled, timed out, or shut down by the orchestrator while
     running, and which therefore never returns a typed result, **does not
     consume the two-attempt no-progress budget**. The budget exists to detect a
     repair that cannot converge, so it counts only attempts that actually
     exercised repair capability and reported an outcome. This generalises
     DEC-167 item 8: a directive-suppressed attempt could not have made
     progress, and an interrupted attempt was never allowed to finish making it.
     Under this rule the consecutive count for
     `task-15:qa-rework:review-5:MF-15-01-MF-15-02` is **0/2**.
  3. **No session-lifetime budget, restated for Developer children.** DEC-168
     item 5 already binds this for QA Lead children; it applies identically
     here. The Phase Runner dispatches the continuation and **waits for the
     child to terminate on its own**. Only per-command budgets are enforced, and
     the task-15 block's budgets are: 900 s for
     `go test -race ./internal/evidence ./internal/journal ./internal/cli ./cmd/fid/...`,
     60 s for `gofmt -l cmd internal`, 300 s for `golangci-lint run ./...`,
     1800 s for `go test -race -coverprofile=coverage.out ./...`, and 60 s for
     `go tool cover -func=coverage.out`. These bound commands, not the session,
     and **change no command, threshold, or accepted status**.
  4. **Loop guard.** Two consecutive externally interrupted attempts under item
     2 are themselves a condition: TL-Lite triages them as
     `orchestration session-liveness`, naming the observed interruption cause,
     the elapsed session duration, and which command was in flight, and
     escalates for an execution-shape decision. That condition never re-enters
     `recovery-budget-exhausted`, because its remedy is the session shape, not
     the repair.
  5. **Durability first.** Attempt 6's first action, before any long command, is
     to commit the existing durable repair - `internal/evidence/execution.go`,
     `internal/journal/idempotency.go`, `internal/journal/project_effects.go`,
     and `cmd/fid/gate_run_test.go` - as a checkpoint commit carrying
     `Refs: phase-21/task-15`. Five sessions and roughly five hours have already
     produced zero commits; no further interruption may be allowed to discard
     this work. Each required command's result is then appended to
     `task-15-done.md` with its observed termination **as that command
     terminates**, so an interruption always leaves inspectable partial
     evidence, exactly as DEC-168 item 5 requires of the QA Lead artifact.
  6. **The recorded attempt-5 gate rows are not evidence.** Every command row
     currently in `task-15-done.md` for the repair is unverified: no
     runner-observed termination and no commit exist for it. Each is rerun
     **byte-identically in session** and replaced with the in-session observed
     value; prior rows may be retained only if labelled as historical attempt
     evidence and never presented as the gate result. The header's recovery
     attempt count is corrected to 6.
  7. **MF-15-02 provenance: forward-only, no history rewrite.** No commit at or
     before `09a2a7b` is amended, rebased, or replaced. The finding's observable
     - task-15's commit history is traceable to its phase and task, and the
     omission is on the record - is satisfied forward: every task-15 commit from
     attempt 6 onward carries `Refs: phase-21/task-15`, and the report's
     Deviations section names `79e3a93` and `8491cde`, states that they omit the
     footer required by `AGENTS.md`, and records that they are not rewritten
     because they are ancestors of the control repair base cited throughout the
     ledger, triage log, and ESC-21-011/012. The false claim at
     `task-15-done.md:79` is replaced by that statement. This is a Tech Lead
     disposition of a QA finding's **repair route**, taken because the reviewer's
     route collides with the phase's historical-integrity rule; the assertion is
     not weakened, dropped, or re-pointed, and the re-reviewer still verifies
     provenance with `git show -s --format=fuller`.
  8. **Nothing else moves.** Task-15 ownership, MF-15-01, MF-15-02's four test
     elements, the five required commands byte-for-byte, all failed evidence and
     prior attempts, VAL-051..054, QA-49..QA-52, the QA plan, validation
     contract, public interface, Event schema, migrations, exit classes,
     thresholds, and budgets are unchanged. No cross-task rerouting: MF-15-01 is
     `internal/evidence/` and MF-15-02 is task-15 tests and report, both inside
     task-15's accepted scope. A **fresh independent Window 5 re-review** still
     gates task 16.
  9. **Exact ledger route.** TRG-21-018's trigger signature is already
     `task-15`-prefixed and its continuing task is task-15, so no identity-only
     successor triage is needed this time. The Phase Runner appends, in order:
     1. `RUN-RESUME` - `Prior run/session or STOP: 2026-08-30 08:36 STOP /
        TRG-21-018 / ESC-21-012`, `Governing TRG: TRG-21-018`,
        `Governing ESC: ESC-21-012`, `Trigger signature:
        task-15:qa-rework:review-5:MF-15-01-MF-15-02`, `Attempt count: 6`,
        `Resolution authority: DEC-170`, `Continuing task: task-15`,
        `Next recovery: task-15 | session-continuation | task-scope |
        task-15:qa-rework:review-5:MF-15-01-MF-15-02 | 6`.
     2. `TASK-RECOVERY task-15` - `Source: TRG-21-018 / DEC-170`,
        `Repair authority: TRG-21-018 / DEC-170`,
        `Category: session-continuation`, `Repair mode: task-scope`, the same
        trigger signature, `Attempt count: 6`,
        `Consecutive no-progress attempts: 0/2` under item 2, allowed files the
        ordinary task-15 scope plus `tasks/task-15-done.md`, forbidden changes
        as in attempt 5 **plus** no rewrite, amend, rebase, or replacement of any
        commit at or before `09a2a7b`.
     3. `REPAIR task-15`; on `Outcome: RESOLVED`, `TASK task-15 COMPLETE` with
        the two-argument report-lint recorded `PASS`.
     4. A fresh independent QA Reviewer performs the linked Window 5 re-review
        and the Runner records `REVIEW review-5b` with
        `Re-review of: review-5`. Task 16 dispatches only on its `PASS`.
  10. **Branches.** If the reran gate is red, that is ordinary in-scope task-15
      repair under this same trigger and directive, and item 2's accounting
      applies to any further attempt. A red result that requires a protected
      contract change is a new triage on its own category, not a return to this
      one.
  11. **Spec sync.** `08_QUALITY_PLAN.md` §3 and `09_AI_WORKFLOW.md` step 16
      record item 2 and item 3 as standing rules. No other specification,
      template, verifier, or workflow definition changes.
- **Consequences:** ESC-21-012 closes with no contract change, no re-slicing, no
  weakened assertion, and no rewritten history. Phase 21 stops converting
  orchestration interruptions into protected escalations - the third such
  conversion after ESC-21-010 - and the recovery budget once again measures only
  what it was designed to measure. Reviewers should expect one commit-first
  task-15 continuation, byte-identical gate reruns, and a report whose provenance
  deviation is stated rather than fabricated.
- **Verification evidence:** `run-ledger.md:2730`, `:2773`, `:2816`, `:2859`, and
  `:2902` each record an externally shut-down child in running state with no
  returned result; `:2865` and `:2908` are the two `NO-PROGRESS` entries;
  `:2913` and `:2928` are the TRIAGE and STOP. `internal/evidence/execution.go`
  contains exactly one `CommitEvent` (`:298`, the Attempt) and one
  `CommitIdempotent` (`:339`, every Receipt outcome). `cmd/fid/gate_run_test.go`
  defines the four `TestGateRun*` cases and the two `assert*` helpers at the
  lines the report cites. `git log` orders `09a2a7b` after `8491cde` after
  `79e3a93`, so the requested amendment rewrites the review commit.
  `task-15-done.md:79` claims an amendment that phase HEAD contradicts.
  `AGENTS.md:65` is the `Refs: phase-N/task-M` convention MF-15-02 relies on.
- **Refs:** TRG-21-017; TRG-21-018; ESC-21-012; DEC-168 items 1, 2, and 5;
  DEC-167 items 7 and 8; DEC-162; `qa/review-5.md` MF-15-01, MF-15-02, and the
  Repair Plan; `qa/qa-plan.md` QA-49..QA-52; `qa/validation-contract.md`
  VAL-051..054; `tasks/task-15-attempt-receipt.md`; `tasks/task-15-done.md`;
  `08_QUALITY_PLAN.md` §3; `09_AI_WORKFLOW.md` §§2.3, 9, and step 16.

### DEC-171: Make DEC-170 session-liveness accounting executable in the workflow controls

- **Date:** 2026-08-30
- **Context:** DEC-170 accepted the standing distinction between a parent
  polling wait, a bounded command, and an externally interrupted child, but its
  item 11 deliberately changed only the specifications. The canonical Phase
  Runner prompt still contained the older two-attempt rule, the router had no
  child-session input, and the ledger verifier rejected any outcome other than
  `RESOLVED`, `MATERIAL-PROGRESS`, or `NO-PROGRESS`. A resumed Runner could
  therefore repeat the same misclassification. The Human directed the super
  agent to make all necessary workflow updates before starting a new Runner.
- **Decision:**
  1. `workflow-route.sh` gains typed child-session outcomes. `WAIT_TIMEOUT` and
     `LIVE` return `ACTION WAIT`; `RETURNED` returns `ACTION EVALUATE`; the
     first consecutive `INTERRUPTED` returns `ACTION REPAIR` and the second
     returns `ACTION TRIAGE`.
  2. `workflow-sequence-lint.sh` accepts `REPAIR Outcome: INTERRUPTED` only
     with exactly one interruption cause, elapsed session duration, and command
     in flight. The following recovery increments total attempt count while
     preserving the prior no-progress count. A second consecutive interruption
     requires TRIAGE before another recovery; a TRIAGE repair advances the next
     total attempt, preserves no-progress, and starts a new interruption streak.
  3. The canonical Phase Runner source, generated provider prompts, operator
     prompt template, run-ledger template, playbook, glossary, quality/workflow
     specs, and tools documentation carry the same rule. A parent polling wait
     timeout never authorizes child interruption.
  4. Mechanical drill 19 proves all routes, required interruption context,
     no-progress preservation, and the two-interruption triage guard.
     `workflow-lint.sh` checks the canonical and generated prompt surfaces so
     provider sync cannot remove the rule.
  5. This is protected workflow-control implementation under the Human's
     explicit authority. It supersedes DEC-170 item 11 only to the extent that
     verifier, template, prompt, and drill changes are required to enforce
     DEC-170 items 2-4. It does not change task-15 ownership, product code,
     product Event schema, VAL-051..054, QA-49..52, the five task commands or
     their budgets, TRG/ESC identity, or DEC-170's exact attempt-6 ledger route.
- **Consequences:** A wait poll can expire repeatedly without killing a healthy
  child. Actual external interruption is durable typed evidence, cannot be
  scored as Developer no-progress, and cannot loop indefinitely without a
  distinct orchestration session-liveness triage after the second occurrence.
  The Phase Runner still appends DEC-170's RUN-RESUME and TASK-RECOVERY before
  dispatching task-15 attempt 6.
- **Verification evidence:** the first red drill-19 run returned nine mechanical
  misses (terminated `exit_code: 9`). A later red schema-v2 case exposed the
  post-triage attempt-number defect (terminated `exit_code: 1`). After repair,
  `project-docs/tools/drills/run-drills.sh` passed with `DRILLS PASS mechanical
  drills completed with no misses` (terminated `exit_code: 0`),
  `project-docs/tools/workflow-lint.sh` passed with `PASS workflow-lint
  completed with no failures` (terminated `exit_code: 0`), generated-agent
  sync check terminated `exit_code: 0`, and the unchanged phase ledger returned
  `VALID state=STOPPED` (terminated `exit_code: 0`).
- **Refs:** DEC-170; DEC-168; ESC-21-012; TRG-21-017; TRG-21-018;
  `08_QUALITY_PLAN.md` §3; `09_AI_WORKFLOW.md` steps 10 and 16.

### DEC-172: Keep the unavailable Attention capability row and replace the defective task-17 absence recipe. Closes ESC-21-013

- **Date:** 2026-08-30
- **Context:** Task 17 correctly stopped after its required production-source
  search matched `internal/workflow/capabilities.go`'s exact accepted row
  `{"attention.v1", false, {"state":false}, "phase-22-unavailable"}`. The
  command expected no matches and `exit_code: 1`, but terminated `exit_code: 0`.
  TRG-21-019 correctly classified this as a protected specification/public-
  interface contradiction. The authorities and chronology resolve it:
  `contracts.md` section 2 and DEC-153/154 require every incomplete capability
  to remain discoverable as `available:false`; the row entered the accepted
  contract before QA refinement commit `b367a3b` later added the broad search.
  VAL-063's behavior is no Attention state/effect, not erasure of negative
  capability metadata. Existing capability tests already require the false row
  and reject an `attention` leaf without mutation.
- **Decision:**
  1. **Verdict audit.** TRG-21-019 and the STOP were correct under available
     authority. Preserve them, ESC-21-013, the failed command, the blocker
     report, and all dirty task-17 implementation. The contradiction is real,
     but its defective side is the later evidence recipe.
  2. **Canonical boundary.** Keep the exact public `attention.v1` row with
     schema 1..1, `available:false`, limits `{"state":false}`, and trust source
     `phase-22-unavailable`. It is required discovery metadata, not Attention
     state or an Effect. Actual Attention tables, Events, Effects, CLI leaves,
     and state/operation surfaces remain deferred to Phase 22.
  3. **Protected evidence-method correction.** Replace only task 17's broad
     expected-nonzero search with:
     `go test -race ./cmd/fid -run '^TestPhase21AttentionBoundary$' -count=1`.
     This Human-authorized change alters an accepted gate definition while
     preserving and strengthening VAL-063's behavior: zero matches could hide
     the required false row, whereas the new boundary requires that exact row
     and rejects any additional production Attention surface.
  4. **Boundary-test contract.** `TestPhase21AttentionBoundary` uses only the
     existing Go toolchain/stdlib and must: require exactly one canonical
     unavailable discovery row; prove `fid --json attention list` is unavailable
     with no journal/projection mutation; inventory non-test production Go under
     `cmd/fid` and `internal` plus SQL migrations; and fail if the canonical row
     is missing/mutated or any other case-insensitive `attention` or
     `attention_items` production match exists. Test-only fixture text is not a
     product surface.
  5. **Synchronized assertion.** VAL-063 and QA-61 use the paired positive and
     negative proof: the task-17 evaluation returns status/reasons; discovery
     retains the exact false row; and the boundary test proves no additional
     state/table/Event/Effect/leaf surface or mutation. The capability matrix,
     product JSON, Phase 22 ownership, task-17 feature scope, VAL identity, and
     validator remain unchanged.
  6. **Failed evidence stays failed.** `task-17-done.md` retains the original
     `rg` row as the pre-DEC-172 failed attempt with observed `exit_code: 0`.
     The continuation adds a distinct revised-gate row and reruns the complete
     six-command task block in session. No result is relabelled, deleted, or
     transcribed from the stopped session.
  7. **Exact recovery route.** The Phase Runner appends `RUN-RESUME` retaining
     TRG-21-019 and ESC-21-013, trigger signature
     `task-17:specification-contradiction:task-17-done.md:attention-absence-vs-unavailable-capability`,
     `Attempt count: 3`, `Resolution authority: DEC-172`, `Continuing task:
     task-17`, and `Next recovery: task-17 | gate-repair |
     task-scope | task-17:specification-contradiction:task-17-done.md:attention-absence-vs-unavailable-capability | 3`.
     The following TASK-RECOVERY uses `Source` and `Repair authority`
     `TRG-21-019 / DEC-172`, category `gate-repair`, task scope,
     attempt 3, and consecutive no-progress `0/2`.
  8. **Continuation and advance.** The task-17 Developer adds the boundary test,
     preserves the existing implementation, commits implementation/test with
     `Refs: phase-21/task-17` before long commands, runs the revised six-command
     block, updates/commits the report, and returns typed completion. After
     two-argument report lint and `TASK task-17 COMPLETE`, task 18 may dispatch.
     Window 6 review remains after tasks 16-18; no special re-review or task-18
     scope change is introduced.
- **Consequences:** ESC-21-013 closes without removing or concealing the public
  unavailable row and without implementing Phase 22 Attention early. The new
  proof is non-vacuous in both directions: it fails if the negative capability
  advertisement disappears and if any real Attention surface appears.
- **Refs:** TRG-21-019; ESC-21-013; DEC-153; DEC-154; DEC-158; VAL-010;
  VAL-063; QA-09; QA-61; `contracts.md` sections 2, 10, and 15;
  `tasks/task-17-evidence-evaluation.md`; `tasks/task-17-done.md`.

### DEC-173: Task 20 owns every path its absence controls walk; bounded source-scope amendment. Closes ESC-21-014

- **Date:** 2026-09-02
- **Context:** TRG-21-020 BLOCKed task 20 at attempt 1 with migration 0014 and
  the first V1 source removals committed (`aef03cf`, report `f189a37`). The
  required absence controls still match production in `internal/bootstrap/`,
  `internal/journal/`, `internal/labels/`, `internal/links/`, and
  `internal/requirements/`, none of which the task's Files To Create Or Modify
  list names. Live inspection of the residues at that checkpoint:
  - `internal/journal/lock.go:20,45,149-150` `LockOwner.ActorRole` and the
    `actor_role` field of `.fid/write.lock/owner.json`;
    `internal/bootstrap/init.go:312` passes `ActorRole: "human"` into it, and
    every retained mutation passes `ActorRole: actor.id`. `04_TRD_CONFIG.md`
    names `write.lock` only as an ignored cross-process lock with no field
    contract; the field is V1 lock attribution, not a public surface.
  - `internal/bootstrap/init.go:18,296,379` `RoleSeed`, `DefaultRoles`, and
    `SeededRoleNames`: a fixed seven-name role registry with no remaining
    production caller once `internal/policy`'s bridge is gone.
  - `internal/labels/{add,attribution,types}.go`,
    `internal/links/{add,remove,attribution,types}.go`,
    `internal/requirements/{sync,attribution,types}.go`, and the in-scope
    `internal/workflow/{define,assign,attribution}.go` and
    `internal/issues/{actor,issue,...}.go`: an `ActorRole` option field whose
    only readers are `standaloneActorForOptions(actor, legacyRole)`, which
    synthesizes `agent:<role>` when `Actor` is empty, and
    `*IdempotencyKey(key, legacyRole)`, which synthesizes
    `legacy-in-process-N` keys. Both are the "fallback reader" that acceptance
    criterion 4 already orders deleted.
  - In-scope residues the report already classifies for removal:
    `internal/issues/{ac_verify,concern_raise,concern_address}.go` writers,
    `internal/maintenance/vacuum.go:29-30`,
    `cmd/fid/issue_mutation_response.go:237`, the retained `internal/dev`,
    `internal/runs`, and `internal/policy` V1 tests, and the
    `issue.claim.seeded` fixtures in `internal/issues/transition_test.go:1408`
    and `cmd/fid/issue_transition_test.go:875`.

  The gate reaches further than the file list in three more places the report
  did not name: (a) `testdata/journal/p20-issue-revision-effects.golden.jsonl`
  carries `dev.started|finished|blocked` and is matched by the required Event
  fixture search whose roots are `cmd/fid internal testdata`; (b) 0014 drops
  `roles`, but `testdata/journal/replay-v1.golden.jsonl`,
  `testdata/journal/corrupt-tail.jsonl`, the `internal/journal` replay,
  projection, phase-2 rebuild, and 100k scale fixtures
  (`rebuild_scale_test.go:191` seeds a `roles` insert), `internal/db/*_test.go`,
  `internal/sqlquery/sqlquery_test.go`, and `cmd/fid/init_test.go` still seed or
  query it, and `go test -race -coverprofile=coverage.out ./...` runs all of
  them; (c) the build-tagged `cmd/fid/integration/` closure journey still drives
  `issue complete`, `run record`, `gate_runs`/`task_completions`/`agent_runs`
  snapshots, and the `has_completion_report` rule, and no task-20 command runs
  that suite, so its convergence would otherwise land on task 21 under
  DEC-167 with task-20 provenance.

  `v1-convergence-inventory.md` §1.1 names task 20 as co-owner of every one of
  the production files above (rows "tasks 2, 20, 21" and "tasks 9, 20"); §6
  assigns task 20 "every mutation `--by`", "the task-2 nonpersistent V1 role
  vocabulary bridge", and "accepted V1 dev/completion/claim/run fixtures";
  DEC-160 §5 states task 20 deletes the bridge and proves the role absence
  control. The inventory, the ten absence declarations, VAL-074/075, 0014, and
  the frozen-table policy are consistent with each other. The task's file list
  is the incomplete artifact: it was written from the removal targets and not
  from the reach of the task's own commands.
- **Decision:**
  1. **Verdict audit.** TRG-21-020's BLOCK was correct under available
     authority: an explicit file list is a task boundary, the inventory's
     co-ownership rows do not override it, and the Developer was right not to
     widen scope alone. Preserve TRG-21-020, ESC-21-014, both checkpoints, and
     every failed command row as immutable evidence.
  2. **Bounded source-scope amendment.** Task 20's Files To Create Or Modify
     additionally covers `internal/bootstrap/`, `internal/journal/`,
     `internal/labels/`, `internal/links/`, `internal/requirements/`,
     `internal/db/`, `internal/activity/` and `internal/stats/` and
     `internal/sqlquery/` (test convergence only), `testdata/journal/`, and
     `cmd/fid/integration/`. Production edits in the added packages are
     limited to the enumerated dispositions in item 3; every other edit there
     is subject-preserving convergence of an assertion that 0014, a V1 leaf
     removal, or the `ActorRole` deletion invalidated, recorded with per-row
     provenance (file, line, invalidating change, accepted replacement).
  3. **Subject-preserving final dispositions.**
     - `internal/journal/lock.go`: delete `LockOwner.ActorRole` and
       `lockMetadata.ActorRole`; `pid`, `hostname`, `started_at`, and `command`
       remain; every `LockOwner{...}` call site drops the field.
     - `internal/bootstrap/init.go`: drop the `ActorRole` lock value and delete
       `RoleSeed`, `DefaultRoles`, and `SeededRoleNames`.
       `BuildProjectInitializedEvent`, the VAL-008 fixture, and init/rebuild
       behaviour are untouched.
     - Retained families (labels, links, requirements, workflow, Issue
       create/update/transition/AC-add/label): delete the `ActorRole` option
       field, the `legacyRole` fallback, and the `legacy-in-process-N` key
       synthesis; `Actor` plus `IdempotencyKey` are the only inputs. Public
       `--actor`/`--idempotency-key` grammar, Event-v2 bytes, digests, replay,
       error envelopes, and every VAL-028..031 assertion keep full strength;
       in-process test callers converge to `Actor: "agent:<role>"` plus an
       explicit key.
     - `testdata/journal/p20-issue-revision-effects.golden.jsonl` and
       `internal/journal/replay_test.go`: the three `dev.*` Events become
       retained Issue families carrying the same revision Effects, and the
       paired `wantTypes`/`priorByType` tables change by the same substitution
       only. The ten-Event revision chain is the preserved subject.
     - `roles` seeds and assertions in journal goldens, `rebuild_scale_test.go`,
       `internal/db`, `internal/sqlquery`, and `cmd/fid` tests converge to a
       surviving table with identical assertion strength and, for the scale
       fixture, identical Event count and shape. Adding a projection branch
       that ignores, skips, or tolerates an Effect on a dropped table is a
       data-model decision and is forbidden; if a fixture's subject genuinely
       requires a legacy `roles` Effect against the final schema, stop and
       escalate.
     - `cmd/fid/integration/`: apply the DEC-162 order to the closure journey.
       Removed-behaviour steps (`issue complete`, `run record`, `policy show`,
       `has_completion_report`, frozen/dropped table snapshots) become their
       accepted behaviour - unknown command, journal unchanged, table absent or
       empty - and every retained assertion stays byte-identical; a step is
       deleted only where a named surviving test owns the accepted equivalent.
  4. **Executable ownership (DEC-167).** Task 20 gains
     `go test -tags integration -race ./cmd/fid/integration/...` as a required
     command under the DEC-168 300 s budget, inserted after its package gate.
     Nothing else in Exact Commands To Run, Expected Non-Zero Commands, the two
     broad inventories, VAL-002/074/075, QA-71/72/93, 0014, the frozen-table
     policy, the public interface, the Event schema, or the QA plan changes.
  5. **Coordination constraints.** Task 9 is complete and its public contract
     is frozen: task 20 may touch task-9 packages only for item 3 and must
     leave every VAL-028..031 test passing at unchanged strength. Task 21
     retains final-schema init/rebuild/maintenance/export parity and the 100k
     benchmark: task 20 may not change `BuildProjectInitializedEvent`,
     projection or rebuild semantics, maintenance findings, export, or the
     benchmark threshold/shape; task 21 must not reintroduce a role lock value,
     role registry, or role fixture, and its spec gains that bullet.
  6. **Standing rule.** `08_QUALITY_PLAN.md` §3 gains the converse of
     DEC-167: a unit of work that owns a required absence or inventory control
     owns every path that control walks, and its file scope must cover every
     search root and every package its required gates compile. A control
     that reaches beyond the listed scope is an incomplete planning
     assignment, not a Developer boundary, and is corrected by amending the
     scope, never by narrowing the control.
  7. **Recovery route.** The Phase Runner appends `RUN-RESUME` naming the
     2026-08-30 19:15 STOP, governing TRG-21-020, governing ESC-21-014,
     trigger signature
     `task-20:specification-contradiction:v1-convergence-scope-vs-absence-gates`,
     `Attempt count: 2`, `Resolution authority: DEC-173`, `Continuing task:
     task-20`, and `Next recovery: task-20 | in-scope-implementation |
     task-scope | task-20:specification-contradiction:v1-convergence-scope-vs-absence-gates | 2`.
     The following `TASK-RECOVERY task-20` uses `Source` and `Repair
     authority` `TRG-21-020 / DEC-173`, category `in-scope-implementation`,
     repair mode `task-scope`, attempt 2, consecutive no-progress `0/2`, and
     allowed files "ordinary task scope under the amended
     task-20-v1-removal.md plus task-20-done.md". The Tech Lead does not edit
     `run-ledger.md`.
  8. **Continuation and advance.** The Developer continues from `aef03cf`,
     commits durable work before the long commands (DEC-170), reruns the
     complete byte-identical amended command block, records `UNCLASSIFIED=0`
     for both broad inventories, and lists per-row provenance for every
     converged row outside the original file list. After the report commit,
     two-argument report lint, `REPAIR task-20 RESOLVED`, and
     `TASK task-20 COMPLETE`, task 21 may dispatch. Window 7 review remains
     after tasks 19-21.
  9. **Triage audit.** This is the fourth escalation whose root cause was a
     task scope narrower than the reach of its own gate (DEC-162, DEC-165,
     DEC-167, DEC-173). Items 2-6 close this instance; the wrap-up
     retrospective should evaluate a planning-time check that compares each
     task's `rg` roots and `go test` package arguments with its file list.
- **Consequences:** ESC-21-014 closes with no change to any command, expected
  status, assertion, migration, or public contract. Task 20 finishes the
  removal the inventory already assigned it, task 21 starts on a green tagged
  suite and a role-free journal fixture set, and the standing rule makes the
  scope-versus-reach check a planning obligation.
- **Refs:** TRG-21-020; ESC-21-014; DEC-160; DEC-162; DEC-165; DEC-166;
  DEC-167; DEC-168; DEC-170; `v1-convergence-inventory.md` §§1, 1.1, 1.2, 3,
  4, 6, 7; `tasks/task-20-v1-removal.md`; `tasks/task-20-done.md`;
  `tasks/task-21-maintenance-rebuild.md`; `04_TRD_CONFIG.md` §1;
  `08_QUALITY_PLAN.md` §3; VAL-002; VAL-028..031; VAL-074; VAL-075; QA-71;
  QA-72; QA-93.

### DEC-174: The journey authors its own workflow; the section 16 prefix is 22 Events. Closes ESC-21-015

- **Date:** 2026-09-03
- **Context:** TRG-21-021 BLOCKed task 22 at attempt 1 with the conformance
  corpus (`a99b427`) and the real-process journey guard (`236540c`) committed.
  The journey reached Gate execution and Receipt Evidence link/evaluation, then
  the terminal `in_progress -> done` transition returned
  `transition-not-permitted` / `workflow-transition-invalid` with process
  `exit_code: 4`. Live inspection of the shipped product at that checkpoint:
  - `internal/workflow/definitions_fixture.go:6-44`: the migration-seeded
    `task-default` reference workflow has states
    `todo/in_progress/dev_complete/qa_review/qa_block/failed_gate/qa_pass/done/
    cancelled` and no `in_progress -> done` edge. Its shortest path from `todo`
    into a `completed`-category state is four edges
    (`todo -> in_progress -> dev_complete -> qa_review -> qa_pass`), and
    `qa_pass` and `done` are both `completed`, so the shipped default would
    produce three surplus transitions and two separate terminal-success guard
    evaluations. The exact-20 prefix is unsatisfiable on it, exactly as the
    report records.
  - `internal/issues/create.go:80-108`: `subtask-default`
    (`todo -> in_progress -> done`) is unreachable without a second
    `issue.created`, because a subtask must have a task parent. `epic-default`
    needs `in_progress -> wrapped -> closed` plus the
    `all_child_tasks_done_or_canceled` rule. No seeded workflow satisfies the
    accepted shape.
  - `internal/workflow/registry.go:23-102` and `internal/issues/transition.go:
    176-200`: the transition engine resolves edges and
    `TerminalSuccess: toCategory == "completed"` from the *projected* workflow
    tables through `workflow.Load`, not from the frozen fixture map. A
    project-authored workflow is therefore first-class to the guard.
  - `internal/workflow/resolve.go:16-50`: `ResolveIssueWorkflow` prefers a
    `kind`-scoped then a project-scoped `workflow_assignments` row and only
    then falls back to the seeded default, so an assignment recorded before
    `issue create` binds the new Issue to the authored workflow.
  - `public-interface.md` §4.1 already accepts `workflow.define` and
    `workflow.assign` as retained leaves emitting `workflow.defined` and
    `workflow.assigned`; `cmd/fid/integration/e2e_p20_issue_journey_test.go:
    22-32` already drives exactly that define-then-assign-then-create sequence
    against the real binary. The route is shipped, accepted, and proven.
  - `internal/workflow/authoring.go:11,112-120`: an Issue-applicable workflow
    must declare exactly one `cancelled` state (PRD-WF-013), and state IDs
    match `^[a-z][a-z0-9_]{0,63}$`, so the authored workflow has four states
    and snake-case IDs.

  The accepted contracts do not contradict each other and the product is
  correct: the guard fail-closed on an undeclared edge is the behaviour
  Phase 21 exists to prove. `contracts.md` §16 was written assuming a direct
  terminal edge that no workflow available to kind `task` provides, and never
  named the journey's workflow. §16 is the defective artifact. No canonical
  spec fixes the count: `07_DESIGN_PLAN_FLOWS.md` §10 and
  `08_QUALITY_PLAN.md` §8 require only "mutates once" and no second Envelope.
- **Decision:**
  1. **Verdict audit.** TRG-21-021's BLOCK was correct. The Developer reached
     the real failure, refused to edit the frozen workflow or the accepted
     contract, and correctly identified in its Concerns that define/assign was
     a mutation outside the enumerated budget - which is precisely why the
     route needed this authority. The recorded category is narrowed in
     hindsight from "specification contradiction, public interface, and
     workflow contract" to a phase planning-authority defect in §16 alone: no
     public interface, Event schema, transition contract, or product workflow
     changes. Preserve TRG-21-021, ESC-21-015, both commits, and every failed
     command row as immutable evidence.
  2. **The journey authors its own workflow.** `contracts.md` §16 now requires
     the journey to define one workflow with `--applies-to task`, exactly one
     `unstarted`, one `started`, one `completed`, and one `cancelled` state,
     and exactly the edges `unstarted -> started` and `started -> completed`,
     then assign it to kind `task`, both before `issue create`. The journey
     never enters the `cancelled` state, which exists only to satisfy
     PRD-WF-013. This is project-owned policy compiled into fid under
     Gate B1/B3, not a product change; the seeded `task-default` chain is
     retained V1 role choreography that Gate D2 retires and is not the V2
     journey's authority.
  3. **Exact Event budget.** The mutation prefix is exactly 22 Events:
     `project.initialized`, `workflow.defined`, `workflow.assigned`, then the
     unchanged prior enumeration in unchanged relative order. The journey
     contains exactly one `envelope.created`, one `workflow.defined`, and one
     `workflow.assigned`; post-rebuild Event delta stays zero; any other count
     is a conformance failure. VAL-086 is amended to the same numbers. This is
     an exact-count assertion replacing an exact-count assertion, with named
     membership on both sides; it is not a weakening, and VAL-084/085/086 keep
     their identity, subject, and strength.
  4. **The failure becomes coverage.** The observed `transition-not-permitted`
     / `workflow-transition-invalid` exit-4 rejection of an undeclared edge is
     required VAL-085 negative-matrix coverage, not discarded evidence. Task
     22's acceptance criteria name it.
  5. **Unchanged.** No change to the task-22 command block, Expected Non-Zero
     Commands, the Event schema, the public interface, the transition or error
     contract, the accepted mutation ordering, the frozen v0.1 reference
     workflows, migration 0014, the frozen-table policy, task-20/21 evidence,
     or any other VAL row. §14's illustrative Event fixture keeps its
     `task-default` precondition value and carries no journey-workflow
     authority. The QA plan's derivation rows for the former count of 20
     (`qa-plan.md` R-17 and the VAL-086 feasibility row) are superseded by this
     decision; Final QA re-derives the enumeration against the amended §16 and
     records the amended count rather than editing the Ready plan's history.
  6. **Recovery route.** The Phase Runner appends `RUN-RESUME` naming the
     2026-09-02 22:26 STOP, governing TRG-21-021, governing ESC-21-015,
     trigger signature
     `task-22:specification-contradiction:contracts-16-event-budget-vs-task-default-transition-chain`,
     `Attempt count: 2`, `Resolution authority: DEC-174`, `Continuing task:
     task-22`, and `Next recovery: task-22 | in-scope-implementation |
     task-scope | task-22:specification-contradiction:contracts-16-event-budget-vs-task-default-transition-chain | 2`.
     The following `TASK-RECOVERY task-22` uses `Source` and `Repair
     authority` `TRG-21-021 / DEC-174`, category `in-scope-implementation`,
     repair mode `task-scope`, attempt 2, consecutive no-progress `0/2`, and
     allowed files "ordinary task scope under the amended
     task-22-walking-skeleton.md plus task-22-done.md". The Tech Lead does not
     edit `run-ledger.md`.
  7. **Continuation and advance.** The Developer continues from `236540c`,
     keeps the committed corpus and journey guard, extends the journey through
     the authored workflow to completion/export and the post-rebuild repeat,
     commits durable work before the long commands (DEC-170), and runs the
     complete byte-identical eight-command block; the eight `not run` rows are
     replaced by observed terminations. After the report commit, two-argument
     report lint, `REPAIR task-22 RESOLVED`, and `TASK task-22 COMPLETE`, the
     phase proceeds to final QA and VAL-087.
  8. **Triage audit.** This is the fifth escalation whose root cause was a
     planning artifact narrower or more optimistic than the shipped surface it
     governs (DEC-162, DEC-165, DEC-167, DEC-173, DEC-174). The four prior
     instances were scope-versus-reach; this one is a hand-derived exact
     expected value that no one executed against the product during planning.
     The wrap-up retrospective should evaluate requiring every exact count,
     enumeration, or ordering assertion in a phase contract to be derived from
     an executed probe before it is accepted, alongside the DEC-173
     scope-versus-reach check.
- **Consequences:** ESC-21-015 closes with no product, command, public-contract,
  Event-schema, or VAL-identity change. Task 22 gains a journey that is
  expressible, that keeps every accepted ordering and guard assertion, and that
  additionally proves two accepted authoring leaves and the undeclared-edge
  negative branch end to end against the real binary.
- **Refs:** TRG-21-021; ESC-21-015; DEC-156; DEC-157; DEC-158; DEC-167;
  DEC-170; DEC-173; `contracts.md` §§14, 16; `public-interface.md` §4.1;
  `tasks/task-22-walking-skeleton.md`; `tasks/task-22-done.md`;
  `qa/validation-contract.md` VAL-084/085/086; `07_DESIGN_PLAN_FLOWS.md` §10;
  `08_QUALITY_PLAN.md` §8; ADR-0011 Gate B1/B3/D2.

### DEC-175: Task 22 owns the bounded activation Event-name repair. Closes ESC-21-016

- **Date:** 2026-09-03
- **Context:** TRG-21-022 BLOCKed task 22 at attempt 2 with the DEC-174 authored
  workflow, the exact 22-Event assertion, and the undeclared-edge negative
  implemented and committed at `a6b3a44`. Seven required commands passed; the
  byte-identical tagged integration gate failed at
  `cmd/fid/integration/p21_conformance_test.go:139` with `journal Event 9 =
  "activation.finalizeted", want "activation.finalized"`. Live inspection of the
  shipped product at that checkpoint:
  - `internal/workflow/activation.go:519`: the shared `transitionActivation`
    writer derives its Event type by string suffix, `EventType: operation +
    "ted"`, then overrides two of the three operations it serves -
    `activation.commit` to `activation.committed` at line 521 and
    `activation.abort` to `activation.aborted` at line 524. `activation.finalize`
    (line 431) has no override, so the derivation emits
    `activation.finalize` + `ted`. The defect is a missing third case in a
    per-operation exception list, not a contract question.
  - The accepted name is fixed and unambiguous in three places:
    `public-interface.md` §4.1 line 369 (`activation.finalize` emits
    `activation.finalized`), `contracts.md` §14 line 436, and `contracts.md` §16
    line 533. All three agree; the product disagrees with all three. There is no
    contradiction to resolve, only an implementation to correct.
  - The misspelling exists in no product source, test, fixture, corpus,
    migration, or golden file. A repository-wide search finds `finalizeted` only
    in `run-ledger.md`, `triage-log.md`, `tl_escalations_current.md`, and
    `task-22-done.md` - that is, only in this recovery's own evidence. No test
    asserts it, so nothing locks the defect in and no assertion is weakened by
    removing it.
  - Rebuild is Effect-driven. No reducer, projection, or rebuild path switches on
    an activation Event type (`internal/journal/event.go:190` special-cases only
    `project.initialized`), so correcting the string changes no derived state,
    no export digest, and no other observable.
  - The committed conformance corpus is derived from `public-interface.md`
    operation/command/kind rows, not from emitted Event types
    (`testdata/p21_conformance/manifest.json:373-377`), so VAL-072 and VAL-073
    stay byte-exact and require no regeneration.
  - Every other family in the §16 sequence exists verbatim in a product writer:
    `project.initialized` (`internal/bootstrap/init.go:345`), `workflow.defined`
    (`internal/workflow/define.go:107`), `workflow.assigned`
    (`internal/workflow/assign.go:111`), `issue.created`
    (`internal/issues/create_commit.go:133`), `issue.transitioned`
    (`internal/issues/transition.go:230`), `activation.prepared` and
    `activation.component-staged` (`internal/workflow/activation.go:209,337`),
    `activation.committed`, `gate.registered`
    (`internal/evidence/registration.go:180`), `envelope.created`,
    `envelope.activated`, `envelope.completed`
    (`internal/workflow/envelope.go:234,304,386`), `claim.acquired` and
    `claim.released` (`internal/claims/generic.go:158,228`), `run.started` and
    `run.finished` (`internal/workflow/run.go:119,148`),
    `gate.execution-started` and `gate.receipt-recorded`
    (`internal/evidence/execution.go:297,371`), and `evidence.linked`
    (`internal/evidence/evidence.go:152`). This is the only mismatch in the
    sequence, so the repair is expected to clear the assertion outright rather
    than expose a queue of further name defects.
  - Provenance: the writer was introduced by task 4, which owned
    `internal/workflow/`. Task 4 could not have observed the defect. Its VAL rows
    assert lifecycle ordering (VAL-011), pre-commit abort (VAL-012), recovery
    unavailability (VAL-013), and trusted-profile resolution (VAL-014/015); none
    asserts an emitted Event-type name, and no non-tagged test in the repository
    asserts any activation Event type. The assertion that catches this was
    authored by task 22 and did not exist before it.
- **Decision:**
  1. **Verdict audit.** TRG-21-022's BLOCK is sustained on its recorded category.
     The Developer reproduced a deterministic failure, refused to change the
     Event contract, the assertion, the command, or the test to make the gate
     pass, and correctly reported that the causal writer was outside its named
     source scope. Protected task-boundary authority was genuinely required; this
     is the correct behaviour, and TL-Lite was right that the tools-only
     control-conformance route cannot reach a production workflow source path.
     TRG-21-022, ESC-21-016, the STOP, the report, and every failed command row
     remain immutable evidence.
  2. **Ownership: task 22, not a reopened task 4.** DEC-167 makes ownership
     assignable only where it is executable, and DEC-173 makes a control's reach
     bound the owning scope. Task 22's required commands are the only gates in
     the phase that execute the accepted activation Event-type name, and the
     assertion that executes it is task 22's own. Reopening a closed, reviewed,
     checkpointed task 4 for a defect none of its gates could observe would
     invert the phase order and add sequencing risk without adding coverage.
     Task 22 owns the repair.
  3. **Bounded source-scope amendment.** Task 22's file scope gains
     `internal/workflow/activation.go`, limited to the Event-type derivation in
     `transitionActivation`, and `cmd/fid/activation_test.go` for the regression
     assertion in item 4. The repair replaces the `operation + "ted"` suffix
     derivation and its two ad-hoc overrides with one explicit, exhaustive
     operation-to-Event-type mapping covering `activation.commit` ->
     `activation.committed`, `activation.abort` -> `activation.aborted`, and
     `activation.finalize` -> `activation.finalized`, returning a typed error for
     any unmapped operation rather than silently synthesising a name. The emitted
     values for commit and abort stay byte-identical. This authority extends to
     nothing else in that file and to no other production file.
  4. **The gate hole closes with the defect.** Task 22 adds a non-tagged
     regression assertion in `cmd/fid/activation_test.go` requiring each
     activation transition to emit its `public-interface.md` §4.1 Event type
     verbatim. It runs under the already-required
     `go test -race -coverprofile=coverage.out ./...`, so the amended scope is
     executable under the existing command block and no command is added,
     removed, reordered, or narrowed. Amending scope without amending the gate
     that runs it would repeat the defect DEC-167 exists to prevent.
  5. **Unchanged.** No change to the Event schema, the public interface, the
     contracts, any command, expected exit, error contract, assertion semantics
     or strength, accepted status or scope values, the transition engine, the
     frozen v0.1 reference workflows, migrations, the committed
     `p21_conformance` corpus, task 22's eight-command block, any VAL row's
     identity, subject, or strength, any threshold, budget, tolerance, or
     policy, or any prior task's accepted evidence. No canonical spec names an
     activation Event type, so no canonical spec edit follows from the repair
     itself; item 8's standing rule is the only spec change.
  6. **Recovery route.** The Phase Runner appends `RUN-RESUME` naming the
     2026-09-03 16:57 STOP, governing TRG-21-022, governing ESC-21-016, trigger
     signature
     `task-22:specification-contradiction:contracts-16-event-budget-vs-task-default-transition-chain`,
     `Attempt count: 3`, `Resolution authority: DEC-175`, `Continuing task:
     task-22`, and `Next recovery: task-22 | in-scope-implementation |
     task-scope | task-22:specification-contradiction:contracts-16-event-budget-vs-task-default-transition-chain | 3`.
     The following `TASK-RECOVERY task-22` uses `Source` and `Repair authority`
     `TRG-21-022 / DEC-175`, category `in-scope-implementation`, repair mode
     `task-scope`, attempt 3, consecutive no-progress `0/2`, and allowed files
     "ordinary task scope under the amended task-22-walking-skeleton.md plus
     task-22-done.md". Attempt 2 is recorded `MATERIAL-PROGRESS` and reached
     Gate execution, Receipt evidence, the terminal transition, Run finish,
     Claim release, and Envelope completion, so under DEC-167 and
     `08_QUALITY_PLAN.md` §3 it does not consume the no-progress budget. The
     Tech Lead does not edit `run-ledger.md`.
  7. **Continuation and advance.** The Developer continues from `a6b3a44`, keeps
     the committed corpus, journey, and negative matrix, applies the item 3
     repair with the item 4 regression assertion as its red test, commits durable
     work before the long commands (DEC-170), and runs the complete
     byte-identical eight-command block with observed terminations. The report
     preserves the attempt-1 and attempt-2 failed evidence verbatim and records
     the repair's red/green evidence. After the report commit, two-argument
     report lint, `REPAIR task-22 RESOLVED`, and `TASK task-22 COMPLETE`, the
     phase proceeds to final QA and VAL-087.
  8. **Standing rule.** A latent production defect that no earlier required gate
     could observe is owned by the first unit of work whose required gate
     observes it. That unit receives the minimum source scope needed for the
     causal repair plus a regression assertion executed by a command it already
     runs; it does not receive generic access to the owning package, and the
     originating unit is not reopened. Canonical text is
     `08_QUALITY_PLAN.md` §3.
  9. **Triage audit.** Every prior Phase 21 escalation root-caused to a planning
     artifact narrower or more optimistic than the shipped surface (DEC-162,
     DEC-165, DEC-167, DEC-173, DEC-174). This is the first whose root cause is a
     product defect: a public contract string synthesised by suffix
     concatenation with per-operation exceptions, and never asserted by any gate
     until the phase's final conformance journey. The wrap-up retrospective
     should evaluate a standing conformance check that every accepted Event type
     in the public interface is emitted verbatim by a product writer, alongside
     DEC-174 item 8's rule that exact expected values be derived from an executed
     probe.
- **Consequences:** ESC-21-016 closes with no contract, schema, interface,
  command, corpus, or VAL change. The product converges to the Event name three
  accepted artifacts already require, the accepted name gains permanent
  non-tagged coverage, and task 22 resumes with its assertions intact.
- **Refs:** TRG-21-022; ESC-21-016; DEC-167; DEC-170; DEC-173; DEC-174;
  `internal/workflow/activation.go:431,519-525`; `public-interface.md` §4.1;
  `contracts.md` §§14, 16; `tasks/task-22-walking-skeleton.md`;
  `tasks/task-22-done.md`; `qa/validation-contract.md` VAL-011..015 and
  VAL-084/085/086; `08_QUALITY_PLAN.md` §3.

### DEC-176: Reconcile Task 22 recovery attempt identity append-only

- **Date:** 2026-09-03
- **Context:** DEC-175 correctly authorizes the bounded Task 22 product repair,
  but its recovery route calls the next recovery attempt 3. The immutable
  schema-v2 TRIAGE TRG-21-022 already records `Attempt count: 3`, while the
  preceding TASK-RECOVERY/REPAIR pair records attempt 2 as
  `MATERIAL-PROGRESS`. The Phase Runner's exact DEC-175 RUN-RESUME candidate was
  therefore rejected with `INVALID RUN-RESUME Attempt count must increment its
  governing TRIAGE`. The append-only ledger retains the TRIAGE, STOP, failed
  append evidence, and a CORRECTION that names the mismatch.
- **Decision:**
  1. TRG-21-022, ESC-21-016, the 2026-09-03 16:57 STOP, DEC-175, the failed
     RUN-RESUME candidate, and all Task 22 command/report evidence remain
     immutable. No historical entry is rewritten or reclassified.
  2. The TRIAGE's attempt identity is authoritative for sequence validation.
     The resumed recovery is attempt 4, exactly one greater than TRG-21-022.
     This numbering reconciliation does not create a product attempt, does not
     consume the no-progress budget, and does not change attempt 2's recorded
     `MATERIAL-PROGRESS`; the resumed recovery starts at `0/2` as DEC-175
     requires.
  3. The Phase Runner appends `RUN-RESUME` naming the 2026-09-03 16:57 STOP,
     TRG-21-022, ESC-21-016, the unchanged trigger signature, `Attempt count:
     4`, `Resolution authority: DEC-176`, `Continuing task: task-22`, and `Next
     recovery: task-22 | in-scope-implementation | task-scope |
     task-22:specification-contradiction:contracts-16-event-budget-vs-task-default-transition-chain
     | 4`. The matching TASK-RECOVERY uses `TRG-21-022 / DEC-175 / DEC-176` in
     both `Source` and `Repair authority`, attempt 4, and no-progress `0/2`.
  4. DEC-175 remains the sole product/scope authority. DEC-176 changes only the
     orchestration attempt identity and supplies the typed STOP-resolution
     authority required by schema v2. It changes no task scope, allowed file,
     forbidden change, command, expected status, assertion, contract, public
     interface, Event schema, transition behavior, threshold, budget,
     tolerance, policy, corpus, VAL identity, or prior evidence.
- **Consequences:** The accepted verifier can validate the append-only recovery
  without a protected control-definition change. Task 22 may execute the exact
  DEC-175 repair as recovery attempt 4 and advance to final QA after its eight
  commands, report lint, REPAIR RESOLVED, and TASK COMPLETE all succeed.
- **Refs:** DEC-175; TRG-21-022; ESC-21-016; `run-ledger.md` 2026-09-03 16:56
  TRIAGE, 16:57 STOP, and 17:44 CORRECTION;
  `tasks/task-22-walking-skeleton.md`; `tasks/task-22-done.md`;
  `09_AI_WORKFLOW.md` §2.3.

### DEC-177: Complete the accepted Receipt read and persistence surface in Task 22

- **Date:** 2026-09-03
- **Context:** After the DEC-175 activation repair, Task 22 recovery attempt 4
  advanced through projection rebuild and then proved that the accepted
  `receipt.get` and `receipt.list` leaves have no registered command or domain
  read API. The public `Receipt` also includes durable `rejected_artifacts`,
  while migration 0013 and `loadReceipt` retain only accepted Artifact links;
  an idempotent replay or post-rebuild read therefore loses rejection facts.
  Separately, `contracts.md` §10 retains the stale prose path `fid projection
  resume`, although the exact wire authority, corpus, Task 18 implementation,
  and root command all agree on `fid resume`.
- **Decision:**
  1. TRG-21-023 and ESC-21-017 are sustained as a protected task-boundary and
     data-model gap. The accepted public interface uniquely determines the
     missing behavior; no new product choice is introduced.
  2. Task 22 owns the minimum causal implementation: read-only Receipt get/list
     APIs in `internal/evidence`; CLI registration/adapters and permanent tests
     in `cmd/fid`; the existing real-process journey; and exact migration,
     schema, golden, and Event-effect updates needed to store
     `rejected_artifacts` canonically as a required JSON array on each Receipt.
     The field sorts by path then reason on write and read. No new table, Event,
     operation, dependency, or mutation leaf is added.
  3. Migration 0013 gains non-null `rejected_artifacts_json` with JSON validity
     and array-type checks. This is an accepted pre-baseline replacement under
     ADR-0011, not a compatibility migration: Phase 21 has not established a
     supported external baseline, task 21's final-schema assertions are updated
     in the same first-observing task, and rebuild remains journal-canonical.
  4. Receipt get returns the exact complete Receipt or
     `receipt-not-found` / not-found / exit 3. Receipt list supports only the
     accepted `--run`, `--obligation`, closed `--outcome`, and `--limit
     <1..10000>` filters, returns exact complete Receipts ordered by
     `started_at,id`, and performs no mutation. Both commands emit no Event.
  5. `contracts.md` §10 is corrected only from `fid projection resume` to the
     already-authoritative `fid resume`. The Receipt Event row is clarified to
     include the canonical rejection JSON in the complete Receipt insert while
     rejected candidates still produce no Artifact or link Effect.
  6. All Task 22 accepted corpus rows, result types, negative matrix, 22-Event
     budget, VAL identities, and eight required commands remain unchanged. New
     focused red/green assertions execute under the existing repository-wide
     unit gate; the tagged journey must pass Receipt get before and after
     rebuild with a zero Event delta.
  7. The Phase Runner resumes the 2026-09-03 18:35 STOP under TRG-21-023 /
     ESC-21-017 / DEC-177 at attempt 6, one greater than the immutable TRIAGE
     attempt 5, then appends one Task 22 `in-scope-implementation` /
     `task-scope` recovery with no-progress `0/2`. Attempt 4 remains
     `MATERIAL-PROGRESS` and its two tagged failures remain immutable.
- **Consequences:** The final conformance journey can exercise every accepted
  Receipt read before and after rebuild without fabricating or dropping public
  Receipt facts. Task 22 reruns the complete byte-identical eight-command block,
  finalizes its report, and then advances to end-of-phase QA.
- **Refs:** TRG-21-023; ESC-21-017; ADR-0011; DEC-167; DEC-173; DEC-175;
  `public-interface.md` §§3, 4.2, 6, 7; `contracts.md` §§10, 14, 16;
  `03_TRD_SCHEMA.md` §5; `tasks/task-22-walking-skeleton.md`;
  `qa/validation-contract.md` VAL-052/072/084/086.

### DEC-178: Reconcile the final-review specification blocker and bound remediation

- **Date:** 2026-09-03
- **Context:** Initial final review `1eee01c` records `VERDICT: BLOCK`. MF-1
  proves that the accepted `03_TRD_SCHEMA.md` §8.2 `0014` SQL contradicts the
  shipped migration, live convergence inventory, task-20 binding constraint,
  and observed final index name, leaving VAL-002 uncertifiable. MF-2 proves
  that two unsupported keys in `04_TRD_CONFIG.md` §2 make its example invalid.
  IMP-1..4 identify stale validation sources, a stale Resume grouping,
  inherited package-inventory drift, and thin in-package coverage. Contrary to
  MF-2's scope note, the accepted task-13 spec held `04_TRD_CONFIG.md` only as
  read context, so a scope and executable-gate amendment is required.
- **Decision:**
  1. Phase 21 remains the approved Heavy phase and enters final-review
     remediation. Phase 22 does not start while VAL-002 or VAL-036 is open.
  2. The shipped `0014_p21_retire_v1_authority_execution` up/down files are the
     exact migration authority. They uniquely agree with
     `v1-convergence-inventory.md` §§4/6 and task 20: remove the three V1
     tables, preserve and rename the replacement Claim index, and delete only
     the displaced role/judgment rules on Issue workflows. §8.2 is corrected
     byte-for-byte to those files. No migration, product behavior, VAL-002
     wording, evidence method, pass/fail criterion, or status is changed.
  3. IMP-1 is fixed now in the Tech Lead-owned validation contract:
     VAL-030 cites PRD-ISSUE-007/PRD-HUM-001; VAL-032 cites
     PRD-RUN-001/002; VAL-033 cites PRD-RUN-001/002/004; and VAL-034 cites
     PRD-RUN-004. Only traceability changes.
  4. IMP-2 is fixed now in `06_DESIGN_PLAN_CLI.md`: `fid resume <issue-id>` is
     a root derived-read leaf, while the remaining named reads stay under
     `fid projection`. The canonical operation remains `projection.resume`.
  5. MF-2 remains Developer-owned QA rework under existing task 13. Amend its
     writable scope to only the §2 example and its report, remove only
     `max_timeout_ms` and `max_artifact_bytes`, and add the exact fresh-workspace
     documented-example command before the unchanged task gate. The Developer
     preserves the initial exit-2 reproduction, reruns the command
     byte-identically to exit 0, changes no product code or test, runs every
     amended task command, commits the report, and obtains Phase Runner
     two-argument report lint.
  6. IMP-3 becomes CF-033 for Phase 22 kickoff: reconcile the `AGENTS.md`
     Architecture package inventory against the live retained tree after the
     Phase 22 boundary is sliced, so it includes current packages and does not
     advertise nonexistent boundaries.
  7. IMP-4 becomes CF-034 and is scheduled, not thresholded: the first Phase 22
     tasks reopening `internal/evidence` and `internal/workflow` must add direct
     in-package tests for their new recovery branches and report package
     coverage before/after. Repository and real-process coverage remain
     required; no numeric coverage threshold is introduced.
  8. Keep CF-024 active in Phase 22 with the existing three-consecutive-fresh
     sample rule because final QA observed only 2.6-3.0 seconds of rebuild
     headroom. Record the inherited path-dependent unknown-config-key question
     as CF-035 for an explicit Phase 22 planning disposition before config
     validation behavior is broadened.
  9. VAL-036 is not carried forward. The Human must answer HUMAN-01a..01d before
     the fresh final review can pass and before the Runner may append RUN-END.
  10. The Phase Runner records the initial final-review BLOCK and routes it
      through TL-Lite as required. This committed decision clears the protected
      MF-1 authority question and authorizes only task-13 QA rework for MF-2.
      After task-13 completion, a fresh independent QA Lead writes a new final
      review, re-decides VAL-002 from fresh migration/spec evidence, verifies
      MF-2 with the documented example and bogus-key control, confirms IMP-1/2,
      and preserves every original finding and failed attempt.
- **Consequences:** No new task or phase is created. The repair order is
  specification reconciliation already committed by the Tech Lead, task-13
  documentation recovery, Human VAL-036, and fresh independent final review.
  IMP-3/4 and the two named risks are explicit Phase 22 planning inputs rather
  than silent Phase 21 scope growth.
- **Refs:** `qa/final-review.md` MF-1, MF-2, IMP-1..4 and BC-F1;
  `qa/human-checklist.md`; `qa/validation-contract.md` VAL-002/030/032/033/034/036;
  `tasks/task-13-executor-policy.md`; `03_TRD_SCHEMA.md` §8.2;
  `04_TRD_CONFIG.md` §2; `06_DESIGN_PLAN_CLI.md`; DEC-154; DEC-157;
  DEC-167; DEC-173; ADR-0011.

## Spec Sync Verdict

Planning-time changes are recorded with their owning decisions. Full positive
verdicts are completed at phase wrap-up.

| Spec | Planning verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending wrap audit | - |
| `01_PRD.md` | touched: see DEC-156 | Phase 21 qualification makes the terminal-success Issue transition and later Envelope completion independently evaluate the same gate-evidence Obligation. |
| `02_TRD.md` | touched: see DEC-154, DEC-156, DEC-157, DEC-158 | Capability availability, shared evaluator, bootstrap, public enums, trust, idempotency, executor boundary, and the unchanged deterministic rebuild target are explicit. |
| `03_TRD_SCHEMA.md` | touched: see DEC-157, DEC-178 | Final-review reconciliation makes the §8.2 `0014` blocks byte-identical to the shipped migration and live removal authority. |
| `04_TRD_CONFIG.md` | remediation pending: see DEC-178 | Task 13 must remove only the two unsupported example keys and prove the documented example is accepted verbatim. |
| `05_DESIGN_PLAN.md` | touched: see DEC-157 | Exact live NDJSON payload union and Human rendering authority is the public-interface contract. |
| `06_DESIGN_PLAN_CLI.md` | touched: see DEC-157, DEC-178 | Complete leaves remain explicit; the target grouping now places `fid resume` at root as the accepted wire contract requires. |
| `07_DESIGN_PLAN_FLOWS.md` | touched: see DEC-156, DEC-157, DEC-174 | Gate/evidence precede terminal transition; rebuild repeats reads and checks without a second Envelope; the operational journey runs on a project-authored workflow, not on a seeded V1 reference chain. |
| `08_QUALITY_PLAN.md` | touched: see DEC-152, DEC-157, DEC-158, DEC-162, DEC-165, DEC-166, DEC-167, DEC-170, DEC-173, DEC-175 | Checkpoint isolation, removal-before-parity, cadence, feasibility, committed-report lint ownership, superseded-contract convergence, invalidating-change ownership, retained-test phase provenance, executable ownership and its converse (a control's reach bounds the owning scope), first-observing ownership of a latent production defect, and interrupted-attempt accounting are explicit in §3. |
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

TRG-21-020 correctly BLOCKed task 20: the explicit file list is a task boundary
and the inventory's co-ownership rows do not override it. In hindsight the
defect was planning-side - the task's absence controls walk `cmd/fid internal
testdata` and its coverage gate compiles `./...`, but its file list was written
from the removal targets rather than from that reach. DEC-173 amends the scope,
adds the tagged gate the removed fixtures need (DEC-167), and records the
converse executable-ownership rule so the next task pack is checked for
scope-versus-reach before dispatch. No historical record needs correction.

TRG-21-021 correctly BLOCKed task 22. The real-process failure was reproduced,
the accepted contract and the frozen workflow were both left intact, and the
Concerns section correctly named define/assign as a route the exact Event budget
forbade. In hindsight the defect was planning-side and narrower than recorded:
`contracts.md` §16 assumed a direct terminal edge that no workflow available to
kind `task` provides and never named the journey's workflow, while the product,
public interface, and Event schema were consistent throughout. DEC-174 amends
§16 and VAL-086, keeps every product and interface contract, and turns the
observed rejection into required negative coverage. No historical record needs
correction.

## Workflow Notes

To be completed at phase wrap-up.

### Deep Retrospective Recommendation

To be completed at phase wrap-up.
