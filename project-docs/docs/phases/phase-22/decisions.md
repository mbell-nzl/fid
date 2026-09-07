# Phase 22 Decisions

> Tech Lead planning, 2026-09-05. Human acceptance recorded in DEC-188.
> Global DEC numbering continues after Phase 21 DEC-179.

## Decisions

### DEC-180: Propose the exact Heavy failure-and-recovery slice

- **Date:** 2026-09-05.
- **Status:** PROPOSED; Human acceptance is pending. No execution authority.
- **Context:** Roadmap section 5 assigns the recovery journey and explicit V1
  removals to Phase 22. Phase 21 is wrapped under DEC-179 and merged as
  89d964a84dbdd297fdf472d6f78b001e2efac702; its current-context branch note was
  stale. Planning started from that clean main HEAD on new branch `phase-22`.
- **Proposal:** Heavy mode, 18 ordered tasks, 80 validation assertions, six
  review windows, an early task-6 subprocess recovery journey, and an independent
  checkpoint after task 12 before archive/import work. Accept `contracts.md`,
  `public-interface.md`, and `schema.md` together before implementation.
- **Product choices presented for approval:** trusted executor observations with
  closed terminated/live/unknown liveness; unknown always blocks; exact Human
  force grant required for live takeover; explicit adoption with no Git mutation;
  one new linked successor; event-derived counters/Stops/Attention; measured
  unknown values stay null; derived close/intentional abandon replace caller
  finish outcomes; activation restore-prior before commit and complete-forward
  after; bounded manifest-v2 import into an empty target with atomic publication.
- **Capability boundary:** attention, activation, and portable export qualify
  only after their complete slice proofs. Claims recovery is implemented but
  `work.claims.v1` stays false until Phase 24 parallel-conflict semantics pass.
  No partial availability or unsupported Human Outcome resolution is implied.
- **Consequences:** Explicit development-baseline amendment of 0013 and new
  0015/0016 migrations; no old-data migration, new runtime dependency, scheduler,
  respawn, Review subsystem, or multi-Issue contract. The roadmap permits
  horizontal schema/logic/adapter/removal tasks within this coherent scenario;
  this is the explicit phase-boundary exception to otherwise isolated changes.
- **Refs:** roadmap sections 5/7; PRD sections 6/7/9/10; TRD sections 2-9;
  schema section 8.3; ADR-0007; `qa/validation-contract.md`.

### DEC-181: Make every task's ownership executable from its first gate

- **Date:** 2026-09-05.
- **Status:** Tech Lead planning decision; applies when the pack is approved.
- **Context:** Phase 21 workflow notes identify unexecutable retained-test
  ownership and a late first end-to-end gate as repeated planning defects.
- **Decision:** Every task requires the ordinary full task gate plus the exact
  tagged integration gate. `execution-scope.md` binds common minimum causal
  convergence scope, per-row introducing/invalidating provenance, and the
  source-inspected tag/migration/benchmark reach. Add dual-adapter conformance
  cases alongside each wire surface; task 18 audits the already-running corpus.
- **Consequences:** No gate waits until the final task to execute inherited
  assertions. No final-task ownership inference from historical names. A newly
  exposed production defect still needs minimum causal scope per DEC-175.
  QA Lead performs fresh feasibility; TL planning records no PASS evidence.
- **Refs:** DEC-162/165/166/167/173/175; quality section 3; task specs 1..18.

### DEC-182: Dispose every carry-forward against this phase's actual scope

- **Date:** 2026-09-05.
- **Status:** Planning dispositions; CF-035 policy proposal awaits approval.
- **Context:** The wrap-up explicitly carried CF-024/033/034/035 and re-targeted
  CF-018. Trigger-based items must be checked against the incoming file scopes.
- **Decision:** Use the following complete disposition, and recheck triggers
  against actual diffs at the task-12 checkpoint and wrap-up.

| ID | Phase 22 disposition |
|---|---|
| CF-008 | Re-defer to a retained Issue-cancellation vocabulary use case. Run abandonment is an execution fact, not Issue cancellation; no trigger. |
| CF-018 | Keep exact production-file trigger. `internal/sqlquery/sqlquery.go` is not planned for edit; if a causal change opens it, assign its behavior-preserving split to that task with ordinary/tagged gates. Test-only edits do not fire it. |
| CF-023 | Re-defer to dedicated temporary workflow-maintenance work; no product boundary connection. |
| CF-024 | Schedule task 17 and final QA, VAL-071. Preserve three fresh <30 s samples and exact reference; headroom is a risk, never a new tolerance. |
| CF-026a | Re-defer optional tracker/workflow inspection design. Recovery resume adds only accepted operational facts. |
| CF-026b | Re-defer Issue lifecycle design until its own accepted product decision; no new cancellation model here. |
| CF-029 | Re-defer QA artifact token-validator changes to workflow maintenance; current evidence rules remain binding. |
| CF-030 | Re-defer report-lint diagnostic enhancement; no rejected-token diagnostic change is needed by this product scope. |
| CF-033 | Kickoff package inventory completed in execution-scope.md and linked TRD amendment. AGENTS.md still needs a matching inventory-only maintenance edit; it is outside this TL's permitted write paths. Keep that edit open for the next workflow-maintenance session; no new package depends on it. |
| CF-034 | Schedule workflow baseline/direct tests in task 2 (VAL-079) and evidence baseline/direct tests in task 7 (VAL-035), the first respective recovery reopens. Record before/after coverage; introduce no threshold. |
| CF-035 | Propose whole-file validation on config-consuming paths; leave config-independent help/archive validation/diagnostics available. Task 11, VAL-047, pins both states; no global validation middleware. |

- **Consequences:** CF-033 is partially discharged, not silently closed.
  Closed CF IDs from earlier phases stay archived and are never reused.
- **Refs:** current decisions register; contracts section 12; quality section 3.

### DEC-183: Pin cross-product review and defer workflow-control redesign

- **Date:** 2026-09-05.
- **Status:** Tech Lead planning disposition; independent review remains pending.
- **Context:** The prior phase recommended a delivery retrospective covering
  gate reach, phase/cross-task ledger identities, missing ESC registration,
  and an earlier thin journey. This session reviewed the wrap-up calibration
  and normal kit contract material; it did not conduct a separate deep audit.
- **Decision:** Apply existing gate-reach/provenance rules and early journey
  now through DEC-181. Require immediate ESC registration under existing rules.
  Track workflow identity/registration verifier design as new CF-036 for a
  dedicated workflow-maintenance planning session; do not modify canonical
  controls, prompts, or skills inside this product phase.
- **Cross-product input:** kit repository HEAD
  c5229617220f97ff9228df261e24718b1d1d6d89; `trd/02_INTERFACES.md` sections
  3.2-3.7 and `FID_COMPATIBILITY.md`. Planning aligned unknown/live force
  semantics and bound observed budget evidence. QA must independently compare
  all proposed semantics and record any differences before a Ready QA plan.
  The stale compatibility matrix's historical fid progress is not current
  qualification evidence; no production adapter is claimed qualified here.
- **Consequences:** Full workflow retrospective remains recommended in its
  dedicated scope. This phase needs neither workflow-control changes nor
  another agent activation by Tech Lead. Existing STOP/recovery rules remain.
- **Refs:** Phase 21 Workflow Notes/DEC-179; `qa/qa-lead-handoff.md`.

### DEC-184: Resolve QA's Wire, Authority, And Resolution Ambiguities

- **Date:** 2026-09-05.
- **Status:** PROPOSED contract refinements; Human acceptance with DEC-180.
- **Decision:** Apply R-01/02/04/05/06/07/09/10/11/12 now. Retain `--project`;
  enumerate 23 final activated actions and their first domain owners; remove
  run.finish authority in task 9; retain bootstrap activation without invented
  activation capabilities. Require observation IDs for recovery and abandonment.
  Reused executor identity uses duplicate-entity with an exact run identity
  shape. Attention cites disposition-containing Events, never just underlying
  facts; rejected recovery raises no Stop. Add Attention's seven-code CHECK,
  separate activation-version-conflict, explicit grammar exceptions, and opaque
  owner_phase semantics. Tasks 2..11 and 18 receive matching obligations.
- **Why:** These close contradictory/undefined inputs before implementation.
  The extra abandonment nullability correction follows the same existing
  liveness rule as R-04; it adds no authorized takeover state.
- **Spec impact:** public-interface.md, contracts.md, schema.md, owning task
  criteria, validation rows and canonical CLI/schema/glossary pointers.
- **Evidence limit:** QA's original DDL feasibility applies to its reviewed
  draft. The amended constraints are not freshly executed by this TL; task 1
  and QA must prove the exact accepted DDL. No previous gate is relabelled.

### DEC-185: Define One Canonical Repository Dirty Digest

- **Date:** 2026-09-05.
- **Status:** PROPOSED contract refinement; Human acceptance with DEC-180.
- **Decision:** Apply R-03 through contracts.md sections 3.1/3.2. Define the
  versioned canonical JSON input, per-path HEAD/index/worktree images, modes,
  byte hashing, ordering, exclusions, rejection cases, and separate complete
  reconciliation digest. Caller fields retain their shape and must use the
  same algorithm; legacy synthetic digest strings are not translated.
- **Ownership:** Task 4 implements observation and deterministic vectors;
  task 12 exercises real unchanged dirty proof and same-HEAD byte-change
  staleness. Both own only the fixtures their accepted changes displace, under
  ordinary and tagged gates. Recovery/import preserve history and never refresh
  proof. No runtime dependency or repository mutation is authorized.
- **Spec impact:** contracts.md, task 4/12, VAL-014/048/049, TRD and glossary.

### DEC-186: Make Boundary Evidence And Known Repairs Executable

- **Date:** 2026-09-05.
- **Status:** Planning decision; applies to the approved pack.
- **Decision:** Apply R-08 as an explicit VAL-056 acceptance matrix with every
  declared limit and unsafe-entry case, including cumulative bytes measured
  during decompression. Preserve the exact resource ceilings; require no forged
  ZIP mode rejected before streaming by archive/zip. Apply R-13 with one exact
  fresh reference sample at tasks 1 and 16. Three fresh samples at task 17 and
  final QA and the 30 s threshold remain unchanged.
- **R-14 ownership:** Task 8 explicitly owns cmd/fid/gate_run_test.go and its
  signal/timeout fixture repair alongside cancelled-outcome work. Preserve
  exact expected outcomes, real subprocesses, and failed evidence; synchronize
  signal completion and prove timeout occurs after spawn. No test deletion,
  broadened outcome acceptance, production deadline change, or retries to find
  green. Repeated focused execution and two full benchmark-gate runs are named
  there; task 17 also includes QA-78's two runs.
- **Provenance correction:** git log -S and blame in this TL session identify
  b69499116e70d4c88423a30daa2010a4ff221f14 (phase-21/task-15) as introducing
  both cases, including timeout_ms=1. git show 76ad781727b9e8545e069c4ac0e9f926012e6af9
  (phase-21/task-21) adds 10000 ms to the *signal* case only. QA R-14's claim
  that it touched timeout_ms=1 is incorrect. Both precede planning base
  89d964a84dbdd297fdf472d6f78b001e2efac702; the original QA report is retained.
- **Earlier failure route:** Ordinary gates already execute this test before
  task 8. If it fails in tasks 1..7, the Runner obtains a bounded same-task
  scope amendment for this exact fixture repair before retrying; do not wait
  for task 8, ignore a failure, or spend attempts forbidden to repair it.
  Task 8 then verifies the integrated repair and owns cancellation convergence.
- **Spec impact:** VAL-056, task 1/8/13/16/17/18 evidence obligations and scope
  inventory; quality section 8 distinguishes historical Phase 21 checkpoints
  from proposed Phase 22 cadence. This planning addition does not modify a protected gate's accepted
  outcome, command body, budget, or tolerance.

### DEC-187: Disposition Every Independent Kit Comparison Item

- **Date:** 2026-09-05.
- **Status:** PROPOSED boundary decisions; Human acceptance with DEC-180.
- **Decision:** Adopt kit-contract-mapping.md as the complete C-01..C-13
  disposition register. C-06b assigns recovery export binding to the kit's
  durable activation coordinator, with Phase 25 qualification proof. Export
  availability alone is insufficient. C-08 keeps fid-derived terminal truth;
  the kit terminalState is a checked expectation, not a fid outcome input.
- **Follow-through:** Record every other RECORD item, including revision-1
  scope, deferred Review/retry units, validate-only mapping, Attention identity
  and disposition, repository observation, derived resumption counts, and code
  mappings. CF-037 tracks adapter/kit contract work until Phase 25 qualification;
  Phase 23/24 roadmap notes own their vocabulary extensions. Unmapped Stop codes
  and absent scope/asOf semantics cannot silently qualify via fallback.
- **Authority limit:** This resolves fid's planning choice and records accepted
  exclusions; it does not approve or edit the external kit contract. Production
  adapter qualification remains unclaimed. VAL-080 still requires independent
  QA recheck of the final implementation and pin currency.

### DEC-188: Human Approves The Refined Phase 22 Pack

- **Date:** 2026-09-05.
- **Status:** ACCEPTED by Human: "Approved. Provide handoff prompt to get phase runner going."
- **Decision:** Approve Heavy mode, all 18 ordered task specs, the 80-assertion
  validation contract, exact slice/wire/DDL contracts, CF-035 under DEC-182,
  and refinements DEC-184..187 as committed in
  8dee0a194a494f575b7c5f64caf6ed9ac7da70f8. DEC-180/182/184/185/187 proposals
  are now accepted; DEC-181/183/186 planning dispositions apply. Earlier
  proposed/pending wording in historical entries describes their original state.
- **Execution authority:** The Phase Runner may start Phase 22 on branch
  phase-22, perform its required QA Lead readiness refresh using the existing
  plan and TL addendum, then dispatch task 1 and continue the approved sequence.
  Approval does not waive any gate, independent review, or Human outcome.
- **Evidence state:** All VAL rows become Ready, not PASS. No implementation
  has started and no RUN-START exists at this handoff. QA's original failed
  attempts remain immutable. Tech Lead ran no quality gates or agents.
- **Handoff:** runner-handoff.md records startup inputs, cadence, R-14 handling,
  ledger/report controls and the phase-end Human/Tech Lead boundary. The Runner
  creates its own ledger through the canonical transactional append control.
- **Limits:** No release, adoption, merge, external kit edit, workflow-control
  redesign, or acceptance of Human validation outcomes is implied.

### DEC-189: Resolve the task-2 retained live-output golden failure as a nondeterministic-test defect

- **Date:** 2026-09-06.
- **Status:** ACCEPTED, Tech Lead. Execution authority; no Human product decision
  is required and none is implied. This is the admissible substantive authority
  for the repair proposed in TRG-22-001..004.
- **Context:** The task-2 review-1 continuation reached MATERIAL-PROGRESS with
  both Must-fix findings implemented, then the byte-identical required command
  `go test -race -coverprofile=coverage.out ./...` failed in retained
  `cmd/fid/TestLiveOutputGoldenFixtures` (terminated `exit_code: 1`). Three
  successive triage records corrected only the formal admissibility of the
  disposition (TRG-22-001 cited DEC-175, which has no heading in this file;
  TRG-22-002 used a backticked trigger signature; TRG-22-003 fixed the raw
  scalar) and TRG-22-004 corrected the recovery-identity binding. None of them
  established the defect's mechanism, so the bounded directive risked being
  discharged by re-baselining a golden rather than by repairing a cause.

- **Classification (why the existing rules do not reach this case):**
  - This is **not** subject-preserving convergence. Task 2 changed Stop and
    Attention domain behavior; it did not displace the live-output framing
    subject. DEC-162/165/166 assign convergence to the task whose accepted
    change invalidated the assertion, and no task-2 change did.
  - This is **not** DEC-175. That rule covers a latent *production* defect first
    observable under a *new* gate. This gate is not new — task 1 and two earlier
    task-2 attempts ran the identical command to `exit 0` — and the defect is in
    the retained test harness, not in production code.
  - Provenance is phase-21/task-14 (`tasks/task-14-live-output.md`, phase-21
    VAL-036). The directory name `cmd/fid/testdata/val-036/` refers to that
    archived phase-21 assertion, not to this phase's unrelated VAL-036; no
    ownership follows from the name.
  - The correct class is a **nondeterministic retained test that first failed
    under an existing required gate**. It is a real latent defect, and its prior
    greens were partly luck, so it is neither waivable as pre-existing nor
    dismissible as flake.

- **Mechanism (arithmetic-supported hypothesis, to be confirmed before repair):**
  `internal/evidence.RunLive` gives the child a fully closed environment
  (`executor.go:66`), so under a coverage-instrumented build the spawned child —
  which is the `cmd/fid` test binary itself, re-entered as `os.Args[0]` — has no
  `GOCOVERDIR` and its exit hook writes
  `warning: GOCOVERDIR not set, no coverage data emitted\n` (54 bytes) to the
  observed stderr. `RunLive` consumes both pipes through one channel in arrival
  order and stops accumulating `OutputBytes` once the ceiling fires
  (`live.go:179-181`), so the last accounted read is whichever chunk wins that
  race. Redacted output before the final write is 29+25+13+13 = 80 bytes against
  a 128-byte ceiling. If the 120-byte stdout payload is consumed first the result
  is the committed golden, 48 retained `x` bytes and
  `observed_bytes_lower_bound` 200. If the coverage warning is consumed first the
  result is 80+54 = **134**, which is exactly the value the failed run recorded.
  The independently derived 134 matching the observed 134 is why this mechanism,
  not read-chunk splitting, is the working explanation.

- **Decision:**
  1. Task 2 owns this repair. Amend its scope to add exactly
     `cmd/fid/live_output_test.go`; task 2's required gates already walk
     `cmd/fid`, so the assignment is executable under DEC-167/173.
  2. The repair must remove the nondeterminism at its source by making the
     fixture's child emit only its intended bytes in every build mode. The
     expected minimal form is a `GOCOVERDIR` literal entry in the definition's
     closed environment, inserted in sorted position between `GATE_TOKEN` and
     `LITERAL` as `validEnvironment` requires, pointing at a fresh `t.TempDir()`.
  3. **Both committed goldens must remain byte-identical.** Neither
     `live-frames.jsonl` nor `human-session.txt` is in the allowed file list.
     Re-baselining either one to a load- or instrumentation-dependent value is
     mining a flaky green and is forbidden regardless of how many times it
     passes afterwards.
  4. No product change is authorized. `RunLive` conforms to its accepted
     contract: `05_DESIGN_PLAN.md` and phase-21 `contracts.md` require exactly
     one `output-truncated` control, continued draining, and observed
     termination, and the control field is explicitly named
     `observed_bytes_lower_bound`.
  5. The mechanism must be confirmed by a deterministic diagnostic receipt
     before the repair lands. If the coverage-instrumented child does not
     reproduce the warning, the hypothesis is falsified: stop, return
     `RECOVERY_REQUIRED`, and hand the observation back to Tech Lead rather than
     adjusting a fixture to fit an unexplained failure.
  6. TRG-22-005 carries the exact bounded directive. It supersedes the
     substantive scope of TRG-22-003 and the pending binding of TRG-22-004 for
     the same task, and retains the review-1 recovery identity and attempt 2.
     The earlier records stay unedited as orchestration evidence.

- **Consequences:** The immutable failed receipt in `task-2-done.md` stays, and
  a mandatory deviation note must carry its verbatim token, termination, root
  cause, repair, and later result. Task 2 still cannot close on the strength of
  a later green alone. Fresh independent review-1 re-review remains required for
  MF-1 and MF-2, and task-3's MF-3 is untouched by this decision.
- **Authority limits:** No protected gate command, result, scope, threshold,
  tolerance, assertion semantic, public interface, schema, dependency, policy,
  or workflow control changes. No ESC is allocated and no Human attention is
  requested. Archived phase-21 VAL-036 is strengthened, never weakened: every
  asserted subject — interleaving, redaction, the three control kinds, the
  non-success terminal frame, and safe escape rendering — is preserved with
  unchanged golden bytes.
- **Spec impact:** `05_DESIGN_PLAN.md` gains one clarifying sentence that
  retained byte counts become lower bounds once the ceiling fires and are not
  reproducible across runs. This records existing accepted behavior so no future
  fixture pins an exact post-truncation count; it changes no frame schema, field
  type, or wire shape. Phase-21 archived contracts are not edited.
- **New carry-forward CF-038:** two residual items, neither blocking. First,
  `cmd/fid/live_output_test.go` still depends on five sleep-separated child
  writes arriving as five separate reads, so same-stream coalescing under
  extreme load remains a theoretical golden break; harden it when a task next
  opens that file with cause. Second, the terminal frame's `output_bytes` is
  typed `int` with no lower-bound qualifier in phase-21 `public-interface.md`
  section 8; fold that wording into the next accepted live-output wire revision.
- **Triage audit note for wrap-up:** TRG-22-001 asserted an inadmissible
  authority and an incorrect DEC-175 classification; TRG-22-002 and TRG-22-003
  repaired only its form and carried the classification forward unexamined. All
  four records were correct to keep this as REPAIR and to refuse an ESC. Record
  at wrap-up that formal-admissibility churn substituted for root-cause analysis
  across three consecutive triage cycles.
- **Refs:** `triage-log.md` TRG-22-001..005; `tasks/task-2-done.md:5-9,95-118`;
  `tasks/task-2.md:32-34,94-116`; `qa/review-1.md:97-98,122-123`;
  `execution-scope.md:6-27`; `internal/evidence/live.go:96-115,175-213`;
  `internal/evidence/executor.go:59-80`; `internal/evidence/policy.go:147-174`;
  `cmd/fid/live_output_test.go:43-107`; DEC-162/165/166/167/173/175/181.

### DEC-190: Repair the protected sequence verifier and reconcile ESC-22-001

- **Date:** 2026-09-06.
- **Status:** ACCEPTED under the Human instruction to do what is necessary to
  resolve the repeated Phase 22 stops. This is separate protected workflow
  maintenance authority, not a normal task-2 control-conformance directive.
- **Scope:** `project-docs/tools/workflow-sequence-lint.sh`, its shell regression
  drills, workflow specification/template/tool documentation, and the Phase 22
  decision, escalation, progress, handoff and repair-evidence records. The
  run ledger may change only through the canonical append transaction.
- **Contract:** A schema-v2 ordinary REPAIR triage may correct another
  unconsumed ordinary REPAIR triage for the same originating task at the same
  dispatch attempt. It must match its unique external record. If an active
  continuation exists, the correction retains that continuation's trigger and
  next attempt; its eventual recovery binds the correcting TRG and all five
  authority fields. This does not permit cross-task substitution, control
  directive replacement, budget resets, or skipping recovery/QA completion.
- **Review identity:** Task-prefixed signatures remain accepted. A review-form
  signature is accepted when it names the active owning review, has category
  `qa-rework`, and its source begins with the exact owning task token (followed
  by `-`, `:`, or end of value). The same identity rule applies at triage,
  recovery-budget/session triage, and STOP resume boundaries.
- **Runner restart:** A normal schema-v2 RUN-RESUME may reconcile a pending
  ordinary triage for its same task, with all STOP-only fields `N/A`. It leaves
  the pending authority and continuation untouched. It cannot replace the
  mandatory recovery bound by a STOP resume or skip a required STOP.
- **Resolution route:** Preserve TRG-22-006 and the existing ESC evidence.
  Reproduce the defects, run red/green boundary regressions and the complete
  mechanical workflow drills, then verify the real Phase 22 append sequence.
  Since neither rejected BLOCK nor STOP was persisted, do not invent either
  event or a STOP-resume tuple. Append a normal RUN-RESUME, TRG-22-005, and its
  bound task-2 recovery at attempt 2; close ESC-22-001 only after canonical
  append receipts. DEC-189 remains the product-test repair authority.
- **Limits:** No product code, product tests, goldens, dependencies, accepted
  product gates, task/VAL definitions, or scope verifier changes. Product
  execution and independent re-review remain under the approved phase roles.
- **Evidence and exact maintenance gates:**
  `project-docs/docs/reviews/2026-09-06-phase-22-workflow-recovery.md`.

### DEC-191: Preserve committed role configuration in canonical sources

- **Date:** 2026-09-06.
- **Status:** ACCEPTED under the same Human workflow-recovery instruction.
- **Cause:** The required workflow lint reports generated prompt drift.
  `cbeadd3c8f06adee3a7e4be6fd4920d23cff4077` changed only generated Codex model
  routing, and `3e06f12fe4ae42e2694958043b95741bbb2c5c8a` added Bash only to
  the generated Claude TL tools while leaving Bash in `disallowedTools`.
- **Resolution:** Record the already committed model choices in canonical
  `project-docs/agents/phase-runner.md`, `qa-lead.md`, and `tech-lead.md`.
  Preserve the Human's TL Bash addition in the canonical source and remove its
  contradictory Bash denial. TL's existing role prohibition on product
  implementation remains binding; shell access enables diagnosis/verification.
  Run the normal agent generator, then the byte-identical workflow lint.
- **Scope:** Those three canonical prompts and their generated provider copies,
  plus maintenance evidence. No model migration is selected, no root config is
  changed, and no product/task/gate semantics change. This is an explicit
  maintenance extension, outside normal isolated control-conformance repair.

### DEC-192: The owning task removes its own progressive recovery-Event fence

- **Date:** 2026-09-06.
- **Status:** ACCEPTED. Resolves obligation 1 of ESC-22-002 / TRG-22-007.
- **Cause:** Task 1 introduced `unavailableRecoveryEvents` in
  `internal/journal/project_effects.go` as a progressive availability fence:
  every Event-v2 recovery type is rejected before effects apply until its
  recovery domain exists. `schema.md` section 5 assigns each fenced Event to
  exactly one owning task (`run.checkpointed` -> 4, `run.recovered` -> 5,
  `run.budget-observed`/`run.execution-blocked` -> 7, `run.closed`/
  `run.abandoned` -> 8, the three `activation.*` rows -> 10). No task spec
  listed the fence file, so task 4 could not implement the `run.checkpointed`
  Event it owns and reported RECOVERY_REQUIRED.
- **Finding:** This is a Tech Lead planning omission under DEC-167/173, not a
  contract change and not a protected task-boundary change. The accepted
  contract already requires `run.checkpointed` to become projectable in task 4;
  the fence is the implementation of "not yet", and its per-key removal is the
  implementation of "now". Deleting one key neither changes accepted Event or
  Effect semantics nor alters replay/live-apply parity: `applyEventEffects`
  keeps the identical generic insert/update path for both, and the
  `run_observations.kind` CHECK already admits `'checkpoint'` from task 1's
  migration.
- **Resolution - task 4 scope amendment:** Add to task 4's Files To Create Or
  Modify, restricted exactly as written:
  - `internal/journal/project_effects.go` - delete exactly the
    `"run.checkpointed"` map key. No other edit to that file: the guard, its
    message, `validateProjectedEvent`, and every remaining key stay
    byte-identical.
  - `internal/journal/checkpoint_projection_test.go` (new, task-4 owned) -
    prove `run.checkpointed` applies its accepted `run_observations` insert and
    `run_runtime` version increment identically under live apply and replay,
    and that the eight remaining fenced Events still reject. Do not edit
    `internal/journal/recovery_schema_test.go`; it is task-1 owned and its
    `run.recovered` rejection must stay intact.
- **Resolution - standing protocol for tasks 5, 7, 8, 10:** The same amendment
  is pre-authorized, so this must not be escalated again. Each owning task's
  scope includes `internal/journal/project_effects.go` restricted to deleting
  exactly its own key(s) from `unavailableRecoveryEvents`, plus proof of live/
  replay parity for the newly available Event. Task 10 owns the last three
  keys and therefore also deletes the now-empty map and the
  `validateProjectedEvent` fence branch it feeds.
- **Resolution - VAL-003 convergence, scheduled not discovered:** Task 1's
  `internal/journal/recovery_schema_test.go` pins VAL-003 partly on the
  `run.recovered` rejection, which task 5's accepted change displaces, and on
  the fence generally, which task 10's accepted change displaces. Both are
  subject-preserving convergence owned by the displacing task under
  DEC-162/165/166, and both are executed by the ordinary race/coverage gate.
  VAL-003's subject is live/replay parity for accepted and for rejected
  Effects; converge it by substituting a still-fenced Event (task 5) and then a
  genuinely unknown Event type (task 10), keeping the corrupt Stop-code branch.
  Never satisfy convergence by deleting the parity assertion.
- **Verification:** Task 4's unchanged full gate, including the byte-identical
  `go test -race -coverprofile=coverage.out ./...` and the tagged integration
  command. No gate command, threshold, tolerance, pass/fail semantic, VAL row,
  validator, task allocation, public interface, DDL, or Event contract changes.
- **Human decision required:** No. Nothing accepted is altered; a planning
  omission in one task's file list is corrected against the already accepted
  `schema.md` section 5 ownership table.

### DEC-193: Calendar-dependent retained fixtures are DEC-189 nondeterminism

- **Date:** 2026-09-06.
- **Status:** ACCEPTED. Resolves obligation 2 of ESC-22-002 / TRG-22-007.
- **Cause:** Task 4's byte-identical required `go test -race
  -coverprofile=coverage.out ./...` failed in
  `internal/workflow/stop_test.go` `TestEquivalentStopsShareAttentionButKeep`
  `DistinctFacts` with `system clock is behind latest journal date 2026-09-06;
  refusing to append out of order`.
- **Established provenance:** The assertion was introduced by
  `ac827959f8359473cbd6546ac6bea1ff1fd20751` `feat(workflow): add typed Stop
  and Attention facts` (`Refs: phase-22/task-2`). Ownership was established
  from the introducing commit, not from a filename or task label.
- **Established mechanism:** `stopRaisedEvent` pins
  `OccurredAt: time.Date(2026, 9, 5, ...)` and the literal EventIDs
  `EVT-20260905-00000{1,2}`, while `initializedWorkflowWorkspace` runs the real
  `bootstrap.InitializeWorkspace`, which writes an Event dated by the live
  clock. On 2026-09-05 both dates agreed. From 2026-09-06 the initialization
  Event is newer, so `prepareEventWithClockOptions` correctly refuses the
  backdated append. The failure is produced by the calendar, not by any code
  change.
- **Not task-4 convergence:** `ac5e6f9bec80e06bfe5b7b33e3611a894d16038a` touches
  only `internal/project/repository{,_test}.go`,
  `internal/workflow/reconciliation{,_test}.go`, and one `activation.go` line
  adding `run.checkpoint` to the allowed capability set. None of that reaches
  the journal writer, the date guard, or the workspace initializer. No accepted
  task-4 change displaced the assertion, so DEC-162/165/166 convergence and
  DEC-175 do not assign it to task 4. TRG-22-007 was right to refuse to infer
  ownership from gate reach.
- **Governing rule already accepted:** DEC-189 covers exactly this class. A
  retained test that is nondeterministic under an unchanged required gate is
  neither convergence nor a DEC-175 latent implementation defect; it is
  assigned to the observing task by scope amendment, and the nondeterminism is
  repaired at its source. Time dependence is nondeterminism: the same bytes and
  the same command pass or fail according to the wall clock. DEC-189 therefore
  already answers this, and no new ownership principle is created.
- **Resolution - task 4 scope amendment:** Add `internal/workflow/stop_test.go`
  to task 4, restricted to removing the fixed-date binding in `stopRaisedEvent`
  and its two call sites. `internal/workflow` is walked by task 4's ordinary
  race/coverage gate, so the assignment is executable, as DEC-167/173 requires.
- **Required repair shape:** Derive the Event date from the workspace instead
  of pinning it, using the production two-phase pattern at
  `internal/workflow/stop.go:185-190`: build the Event with an empty EventID,
  call `journal.PrepareEvent` to obtain the writer-assigned ID, then rebuild the
  Effects with that ID. Any equivalent repair that leaves no absolute calendar
  literal on this append path is acceptable. Forbidden: re-pinning a newer fixed
  date, faking or freezing the clock inside product code, relaxing or bypassing
  the `prepareEventWithClockOptions` out-of-order guard, changing
  `internal/journal/writer.go`, or weakening the assertion. The subject is
  unchanged and must stay unchanged: two equivalent Stops share one open
  Attention while keeping distinct Stop facts and two links, asserted as
  `2/1/2`.
- **Task 2 retains subject ownership:** Task 2's gate passed correctly when it
  ran; the fixture decayed afterwards. Task 4 makes the exact date-binding
  repair and records full provenance; task 2's Stop/Attention semantics
  ownership and review-1b PASS are untouched. This mirrors the already accepted
  `execution-scope.md` R-14 pattern, where the encountering task takes a bounded
  same-task amendment and the scheduled owner keeps the subject.
- **Required report provenance:** Task 4's report must record, for this row, the
  introducing commit, phase/task, the calendar mechanism, the evidence that no
  accepted task-4 change displaced it, the exact repair, and the command that
  executes it. The preserved failed receipt stays in the report.
- **Standing prevention:** No test in this repository may pin an absolute
  calendar date on an Event that is appended through the journal writer into a
  workspace initialized by the live clock. `internal/workflow/stop_test.go` is
  the only current instance; `internal/journal/{event_v2,writer}_test.go` pin
  dates only in self-contained journal directories and are unaffected.
  Fixed-date rows written directly into journal files or SQL, as in
  `internal/journal/recovery_schema_test.go` and
  `internal/db/recovery_migration_test.go`, remain allowed. Every later task
  and QA review applies this rule to new tests.
- **Verification:** Task 4's unchanged full gate, plus one repeat of
  `go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttention`
  `ButKeepDistinctFacts$' -count=5` to show the repair is stable rather than
  newly date-lucky. No gate command, threshold, tolerance, pass/fail semantic,
  VAL row, validator, task allocation, or product contract changes.
- **Human decision required:** No. This applies accepted DEC-189 to a second
  instance of the same class and changes no accepted contract or test semantic.

### DEC-194: Task 4 under-delivered its own accepted checkpoint and inspection domain scope

- **Date:** 2026-09-06.
- **Status:** ACCEPTED. Resolves ESC-22-003 / TRG-22-008.
- **Cause:** Task 6's real built-binary journey reached `fid run inspect` and
  terminated `exit_code: 1` / `FAIL`. The Developer correctly established that
  `internal/workflow` contains no `run.checkpoint` mutation/Event producer and
  no `run.inspect` read API, and correctly refused to fabricate them in the CLI
  or to widen its own scope. TRG-22-008 then classified this as a protected
  task boundary with an **ambiguous owner**.
- **Finding - the owner is not ambiguous.** Four accepted records already name
  it, and they agree:
  - `tasks/task-4.md` Files To Create Or Modify names
    `internal/workflow/run_checkpoint.go` and
    `internal/workflow/run_checkpoint_test.go`.
  - `schema.md` section 5 assigns the `run.checkpointed` Event producer to
    task 4.
  - `public-interface.md` section 3 assigns `run.inspect` and `run.checkpoint`
    to tasks `4,6`; every peer row uses the same split (`run.observe` 3,6;
    `run.recover` 5,6), where the lower task owns the domain operation and
    task 6 owns the CLI leaf.
  - `contracts.md` section 3 places `run.inspect` read-only reconciliation and
    `run.checkpoint` observation semantics in task 4's subject, and section 3.1
    states "Task 4 owns the observer".
  Task 4 delivered `internal/project/repository.go`, `reconciliation.go`, the
  `run.checkpoint` activated capability, and the DEC-192 fence removal, then
  closed COMPLETE. `internal/workflow/run_checkpoint.go` has never existed in
  any commit on any branch, and no inspection read API was written. This is
  undelivered accepted in-scope work by task 4, not a contract change, not a
  new task boundary, and not a decision the Human must make.
- **Finding - the enabling defect is mine.** Task 4's acceptance criteria
  demanded VAL-013..016, `run.checkpointed` live/replay projection parity, and
  the DEC-193 row. None of them demanded the producer or the read API that the
  same spec's file scope named. A truthful COMPLETE was therefore reachable
  with two named domain files unwritten, and the omission stayed invisible
  until the first process journey two tasks later. `internal/journal/`
  `checkpoint_projection_test.go` proves the Event projects correctly when
  something emits it; nothing emitted it. Recorded for wrap-up: a task's file
  scope is not self-enforcing, so every named production path needs an
  acceptance criterion or VAL citation that fails when it is absent.
- **Resolution - execute the missing task-4 domain scope inside task 6.**
  Reopening task 4 is rejected: it is closed, tasks 5 and 6 have run on top of
  it, and the Phase Runner's only recovery route is a same-task continuation of
  the currently active task. Task 6 is the first task whose accepted validation
  cannot pass without this behavior (VAL-022 names `inspect`; VAL-023 covers
  every section-3 recovery leaf), so execution moves to task 6 with ownership
  provenance recorded rather than relabelled. Add to task 6's Files To Create
  Or Modify, restricted exactly as written:
  - `internal/workflow/run_checkpoint.go` (new) - the `run.checkpoint`
    mutation and `run.checkpointed` Event producer, structured exactly like the
    task-3 `ObserveRun` producer in `internal/workflow/run_observation.go`:
    write lock, projection freshness, idempotency replay, Run load, active-Run
    and `expected_version` preconditions, trusted `activated-harness` Actor,
    `run.checkpoint` capability and Envelope grant, re-observation under the
    lock, `reconciliation_digest` comparison rejecting with
    `recovery-state-conflict`, then exactly the `schema.md` section 5 Effects:
    insert one `run_observations` row with `kind='checkpoint'` and increment
    `run_runtime.version`. No new Effect, table, or column.
  - `internal/workflow/run_checkpoint_test.go` (new) - in-package positive,
    stale-version, inactive-Run, authority-denied, conflicting-observation, and
    idempotent-replay coverage.
  - `internal/workflow/run_inspect.go` (new) - the read-only `run.inspect`
    composition returning `{run, reconciliation, liveness}` exactly per
    `public-interface.md` section 3, built from the existing `loadRun`,
    `project.ObserveRepository`, `ReconcileRepository`, the latest accepted
    checkpoint as reconciliation base, and the same-package `latestLiveness`
    reused directly, never duplicated. It appends no Event, takes no mutation
    write lock, and writes nothing to the repository.
  - `internal/workflow/run_inspect_test.go` (new) - in-package coverage that
    repeated unchanged inspection returns the same reconciliation digest
    (contracts section 3.2), that a checkpoint moves the comparison base, and
    that inspection is byte-inert against refs, index, stash, and worktree.
  No other `internal/workflow` file may be edited except the minimum needed to
  register these operations.
- **Executability check.** Task 6's required commands already run these files:
  `go test -race -coverprofile=coverage.out ./...` compiles and executes
  `internal/workflow`, and `go test -tags integration -race`
  `./cmd/fid/integration/...` builds the binary that exposes the leaves. The
  assignment is therefore executable at the task that receives it; no gate,
  tag, or environment switch is added.
- **What does not change.** No public interface, wire shape, error code, exit
  code, DDL, Effect list, Event semantic, VAL row, VAL wording, validator,
  operating mode, gate command, threshold, tolerance, dependency, or workflow
  control changes. Task 4 stays COMPLETE and its report is preserved unedited
  as the durable evidence of what it did and did not deliver. VAL-013..016 stay
  owned and validated by task 4.
- **QA note carried to window 2.** VAL-013's fixture clause reads "unchanged
  tree/checkpoint" and VAL-022 names `inspect`. QA Lead must audit the
  checkpoint producer evidence under task 6 with task-4 provenance, not treat
  its absence from task 4's report as a task-6 defect. This is a review
  instruction, not a contract change.
- **Convergence and provenance.** The task-6 report must record, for each of
  the four files above, that the accepted owner of record is phase-22/task-4,
  that no commit ever created them, that DEC-194 transferred execution to task
  6, and the exact command that executes them. Provenance is recorded, not
  rewritten.
- **Verification:** Task 6's unchanged full gate in its existing order, plus
  the byte-identical rerun of the preserved red command
  `go test -tags integration ./cmd/fid/integration -run`
  `'^TestP22EarlyRecoveryJourney$' -count=1`. The preserved failed receipts and
  the `RECOVERY_REQUIRED` history in `task-6-done.md` stay intact.
- **Human decision required:** No. An accepted allocation already existed; this
  corrects a mis-classified owner and moves execution of undelivered in-scope
  work to the first task that cannot pass without it.

### DEC-195: The accepted `cancelled` terminal outcome was never allocated a DDL row

- **Date:** 2026-09-07.
- **Status:** ACCEPTED. Resolves ESC-22-004 / TRG-22-009.
- **Cause:** Task 8 implemented the accepted fid-observed `cancelled` outcome and
  then could not persist its Receipt. `contracts.md` section 7 states that
  cancellation adds the terminal outcome `cancelled`; `public-interface.md`
  section 4 returns `gate-unsuccessful`/8 with outcome `cancelled`;
  `validation-contract.md` requires VAL-031 to prove exact cancelled outcomes and
  fails an outcome that is mislabelled. But `execution_receipts.terminal_status`
  admits only `succeeded`, `failed`, `timed-out`, `signaled`, `spawn-failed`.
- **Finding - this is a real specification contradiction, and it is mine.**
  Unlike DEC-194, no accepted record assigned the missing work: `schema.md`
  section 1 lists five 0013 CHECK member-list replacements for Phase 22 and
  simply omits the one its own accepted contract requires, and migration 0015
  never touches the table. Section 8.3 of `03_TRD_SCHEMA.md` incorporates
  `schema.md` by reference, so the omission propagated into the canonical record.
  There was nothing for task 8 to deliver and nothing for TL-Lite to repair.
  TRG-22-009's BLOCK and ESC-22-004 were correct in verdict, category, and
  reasoning, and the refusal to substitute an admitted status was exactly right:
  mapping cancellation to `signaled`, `timed-out`, or `failed` would have
  falsified an accepted, Human-approved behavior and defeated VAL-031's own
  failure condition. No no-progress attempt is counted; the blocker was a missing
  authority, not exhausted repair capability.
- **Resolution - amend the DDL, not the assertion.** The behavior was accepted
  under DEC-188 by the Human; only its persistence derivation was missing. The
  amendment route is already accepted: `schema.md` section 1 is a deliberate
  pre-adoption in-place baseline replacement of 0013 requiring fresh
  init/rebuild, and task 1 already executed five such replacements. Section 1 now
  carries two further rows:
  - `execution_receipts.terminal_status` CHECK member list gains `'cancelled'`.
  - The composite terminal CHECK gains exactly one clause,
    `OR (terminal_status = 'cancelled' AND timed_out = 0 AND NOT (exit_code IS
    NOT NULL AND signal IS NOT NULL))`, with the other four clauses
    byte-identical. Section 1 carries the exact five-clause block to copy.
  The clause is derived from observed behavior, not convenience: `RunLive`
  classifies context cancellation as `cancelled` and then records whichever of
  signal or exit code the process state actually carried, so a cancelled Receipt
  must admit signal-only, exit-code-only, and neither, while still forbidding the
  physically impossible pair and reserving `timed_out = 1` to `timed-out`.
- **Execution assignment.** Task 8, as an authorized same-task continuation with
  task-1 provenance recorded. Task 1 stays COMPLETE and unedited; reopening it is
  rejected on the DEC-194 grounds. Task 8 is the only task whose accepted
  validation cannot pass without the row, and its scope now names
  `migrations/0013_p21_operational_kernel.up.sql` limited to those two exact
  replacements, plus the two convergent surfaces below.
- **Consequential surfaces, both bounded and named.**
  - `internal/db/operational_kernel_migration_test.go` holds a retained golden
    copy of the `execution_receipts` DDL string. That is subject-preserving
    convergence displaced by this accepted change, and it converges to the same
    text with per-row provenance.
  - `internal/evidence/receipt_read.go`'s `validReceiptListOptions` is a closed
    outcome filter set. A stored status that its own read API rejects is an
    incoherent domain, so `"cancelled"` joins it; `--outcome unknown` still
    rejects. `internal/workflow/evaluation.go` needs no change: it derives
    `ReceiptSuccessful` from an `outcome == "succeeded"` allowlist, so a
    cancelled Receipt cannot pass without any further guard.
- **Executability check.** Task 8's unchanged required commands already execute
  every assigned file: `go test -race -coverprofile=coverage.out ./...` compiles
  and runs `internal/db` (which applies the real migration directory), plus
  `internal/evidence` and `cmd/fid`; `go test -tags integration -race`
  `./cmd/fid/integration/...` builds the binary and initializes fresh workspaces
  from the same SQL. No gate, tag, build constraint, or environment switch is
  added, so the scope is executable where it lands.
- **What does not change.** No Event, Effect, table, column, index, wire shape,
  error code, exit code, VAL row, VAL wording, validator, gate command,
  threshold, tolerance, task allocation, operating mode, dependency, or workflow
  control. VAL-031 is not weakened; it becomes satisfiable as written. Migrations
  0015 and 0016 are untouched, no new migration is added, and no existing
  development data is converted: fresh init/rebuild remains the accepted route,
  exactly as section 1 already declares. Databases created before this change are
  disposable; golang-migrate tracks version only and stores no file checksum.
- **Human decision required:** No. The Human already accepted the cancellation
  behavior under DEC-188. This supplies the persistence derivation that the
  accepted behavior always implied and that my schema record omitted. The
  pre-adoption baseline-replacement posture, and its fresh init/rebuild
  consequence, are unchanged from what was approved.
- **Carried to wrap-up.** DEC-194 recorded that a task file scope is not
  self-enforcing. DEC-195 records the converse gap: an accepted behavioral
  contract clause is not self-enforcing either. `contracts.md` section 7 added a
  closed-domain value that no schema row, task scope, or acceptance criterion
  ever carried, and nothing failed until an implementation reached the CHECK.
  Every accepted clause that names a closed-domain member needs a matching
  schema row or an explicit statement that none is required.
- **Verification:** Task 8's unchanged eight-command gate in its existing order.
  Its `task-8-done.md` blocker history, preserved failed receipts, and BLOCKED
  gate rows stay intact and are appended to, not rewritten.
- **References:** ESC-22-004; TRG-22-009; `schema.md` section 1;
  `contracts.md` section 7; `public-interface.md` section 4;
  `qa/validation-contract.md` VAL-031; `tasks/task-8.md`;
  `tasks/task-8-done.md`; `migrations/0013_p21_operational_kernel.up.sql`.

### DEC-196: Make semantic conformance and Human fixture delivery executable

- **Date:** 2026-09-07.
- **Status:** ACCEPTED by the root delegate under the Human's explicit
  authority to complete the phase and directly resolve workflow friction;
  authority is recorded in `delegate-handoff.md`. This allocates accepted
  behavior and evidence, without changing the product contracts.
- **Observed gaps:** The task-6 shared fake introduced at
  `f3be616f0587881ee3f5c3716ffbf6167e3f2a20` returns fixture field names with
  `-fake` appended and echoes declared outcomes. Its shared assertions check
  field presence, so `replayed: "replayed-fake"` passes. That is insufficient
  for contracts section 11's shared semantic inputs/expected outputs. The
  retained ordinary Human resume renderer prints only Issue ID and completion
  status, while task 18 prohibits new production behavior. Finally, task 18
  promises Human checkpoints but names no retained generator; existing
  `t.TempDir` fixtures and TestMain's temporary binary disappear after tests.
- **Task 9 assignment:** Its already accepted runner/fake extension now
  explicitly includes independent in-memory request/state semantics, typed
  value and identity assertions, same-key replay, rejection without effects,
  and this task's Stop/Attention/budget/closure relationships. Include existing
  recovery cases and record task-6 introducing provenance. Task 6 and its
  review history stay intact; review 3 independently checks this repair.
  Task 18 retains the full final corpus completeness obligation.
- **Task 12 assignment:** Implement ordinary Human recovery/resume rendering
  from the JSON semantic projection under its existing `cmd/fid/resume.go`
  scope, with explicit paired rendering evidence. Public-interface section 6
  supplies the content; the excluded Human brief product stays excluded.
  Task 18 audits this delivered behavior under VAL-053.
- **Task 12 recovery-result completion:** Root subsequently inspected
  `cmd/fid/run_recovery.go`: its Human mutation result prints only the
  successor ID, while section 6 also requires predecessor, adoption, authority,
  budgets and force-risk context. Introducing commit
  `69941fc1b6262fdd7efaaeecfe72760390dbd511` records task 6. Assign that file's
  rendering and paired real-binary regressions to task 12 alongside resume,
  explicitly before its dispatch. Reuse the mutation's semantic result;
  preserve JSON, actions and authority. Task 6 remains complete with history
  unchanged. This closes the same accepted Human-rendering allocation, with
  task-6 provenance, rather than deferring production work to task 18.
- **Task 18 assignment:** Add `p22_human_fixtures_test.go` plus
  `qa/human-fixtures.md` for reproducible retained checkpoints, durable binary,
  exact IDs and neutral inspection commands. The optional
  `FID_P22_HUMAN_FIXTURE_ROOT` output mode creates fresh children and never
  overwrites existing fixtures. Default tagged tests execute setup and a
  subprocess smoke test proves that retained fixtures survive generator exit.
  Add the explicit focused generation command after the existing six commands;
  every original gate remains byte-identical, ordered and mandatory. The new
  invocation is evidence/delivery for an accepted Human fixture obligation,
  not a substitute or a narrower pass criterion. No dependency is added.
- **Authority and limits:** VAL rows, final validator ownership, product wire,
  DDL, Event semantics, thresholds and all existing gate commands are unchanged.
  Root explicitly authorizes the additional task-18 delivery command before
  task dispatch; no failed gate is being redefined. Only the Human judges
  VAL-076..078. Expected judgments remain separately inspectable so initial
  prompts can ask for the Human's understanding without supplying the answer.
- **Verification:** Task 9's ordinary and tagged gates execute its fake/corpus
  repair; review 3 independently checks semantic assertions. Task 12's gates
  and the independent checkpoint execute both renderings. Task 18's full
  tagged gate exercises fixture setup/survival; its additional exact command
  generates tomorrow's retained inputs. Final QA reruns the setup and records
  all 77 deterministic outcomes, leaving Human outcomes pending until tested.
- **Evidence attribution:** This is source/allocation inspection, not runtime
  validation of unfinished tasks. The independent preflight artifact is
  `qa/remaining-scope-audit.md`; its findings are not a substitute for the
  later Developer or QA evidence.

### DEC-197: Supply task 9's required domain/read and pre-start frame support

- **Date:** 2026-09-07.
- **Status:** ACCEPTED by the root delegate under the explicit Human execution
  authority in `delegate-handoff.md`. Bounded accepted-behavior allocation.
- **Cause:** Task 9's accepted first failing test demonstrated the missing
  committed-denial frame. Its builder lives in `internal/evidence/live.go`,
  outside the listed CLI files. While wiring the other accepted leaves, the
  Developer found no Stop/Attention list APIs and no durable external
  run.observe-budget operation. Root inspected the domain functions: task 2
  provides mutation/GetStop primitives and task 7 provides sample admission,
  derivation and Stop effects; they do not supply those public leaf entrypoints.
- **Bounded assignment:** Task 9 adds `internal/workflow/stop_read.go` and its
  matching test for the accepted read contract; `internal/workflow/run_budget.go`
  and matching test for the accepted durable observation operation; and only
  necessary shared helper support in `internal/workflow/budget.go`. Domain
  authorization, measurement/source identity, versioning, idempotency, budget
  derivation and Event effects stay in workflow. No CLI-local counterfeit
  domain state or caller-supplied counters are authorized.
- **Frame assignment:** Extend `internal/evidence/live.go`'s AppendFidError
  and matching `live_test.go` only to support public-interface section 4's one
  pre-start terminal fid-error frame with no Attempt/start/output. Preserve
  started-frame consistency, duplicate-terminal rejection and all observed
  execution semantics. The Developer had made the minimal causal builder edit
  while reproducing the first test; root records the omitted scope now before
  further domain work or task closure, rather than treating it as self-granted
  general evidence-package access.
- **Retirement support:** The accepted task-9 retirement also removes the
  obsolete run.finish action/request/validator/producer from
  `internal/workflow/run.go`. Root inspected the bounded removal diff and names
  that path explicitly; no other Run behavior is assigned by this addition.
- **Shared Stop call sites:** Task 9's VAL-037 requires the same applicable
  predicates across guarded operations. Root verified that contracts section 5
  binds these actions, while the named task scope omitted their existing
  producers. Add only shared applicable-Stop calls and matching tests in
  `internal/evidence/execution.go`, `internal/workflow/run.go` and
  `internal/workflow/recovery.go`. No new action is blocked: each Stop's exact
  declared actions/subject and qualifying disposition govern. Retain existing
  authority, idempotent replay, no-effects ordinary rejection, budget denial
  atomicity and completion of an already-admitted Attempt. Shared predicates
  stay in the existing workflow evaluation boundary. The unchanged task gates
  execute each call site; report predecessor provenance and both applicable/
  unrelated/resolved Stop cases. This explicitly extends the preceding
  run.go retirement-only assignment for its required run.start call site.
- **Resolution replay repair:** Root's new shared corpus executed successful
  stop.resolve/attention.resolve first deliveries, then both identical cached
  requests failed with invalid-argument/2 on replay. The focused shared corpus
  command terminated exit 1; the fake's independently derived immutable replay
  and other real cases passed. Assign only the causal replay correction and
  matching regression tests in `internal/workflow/stop.go` to task 9. Preserve
  the task-2 source provenance, exact result/Event identity, rejection rules,
  failed receipt and unchanged gate commands. This is the accepted idempotency
  contract, not a new public behavior or a reopened task 2.
- **Typed Stop errors:** The retained Stop domain returns plain errors for
  not-found, version conflict and nonqualifying resolution, while the new CLI
  must expose public-interface section 4's accepted distinct codes/details.
  Permit typed error carriers in `internal/workflow/stop_errors.go` and the
  corresponding narrow return-site changes in `stop.go`, with regression tests.
  The CLI maps these types; it must not infer control semantics from message
  text. Existing rejection predicates and every accepted wire code stay intact.
- **Retirement cleanup:** The first task-9 lint attempt found
  `sameRepositoryLocation` unused in `internal/workflow/run_read.go` after
  FinishRun removal. Root inspected the exact helper and authorizes only its
  removal in that additional path. Preserve the failed lint attempt, repair the
  separately identified in-scope test-helper finding, and rerun the identical
  lint command before continuing the unchanged gates. No control definition
  changes; this is a direct consequence of the accepted retirement.
- **Provenance:** Existing Stop domain commit
  `ac827959f8359473cbd6546ac6bea1ff1fd20751` cites phase-22/task-2; budget
  domain commit `d85ab7018903300ccb53199c49fdea04e1a0bdc6` cites task 7;
  retained live builder commit `bead7efb134f80726022f1e38d7b6749fe550073`
  cites task 8. The report identifies these predecessors and the accepted
  task-9 leaf requirements, without relabelling their completed work. Prior
  tasks and reports stay closed and unedited.
- **Verification:** The unchanged five task-9 commands execute all additional
  domain/read/frame files, matching tests and the real-binary fixtures. Include
  qualifying and rejecting read/sample cases, no-effects rejection, Event
  live/rebuild parity and pre-start versus started-terminal cases. Review 3
  independently audits the bounded additions. No new gate or environment mode.
- **Limits:** No schema, Event shape, closed value, wire result/error, exit
  code, budget rule, threshold, dependency, VAL row or workflow control changes.
  This completes already accepted operations; it introduces no product policy.

### DEC-198: Distinguish Prepare Bindings From Immutable Component Staging

- **Date:** 2026-09-07.
- **Status:** ACCEPTED by the root delegate exercising the Human's explicit
  authority to complete and smooth the remaining Phase 22 execution, recorded
  in `delegate-handoff.md`. Root requested an independent contract audit and
  then explicitly authorized this bounded Tech Lead resolution before task 10
  dispatch. It is not a Developer inference or a new prepare input contract.
- **Contradiction:** `contracts.md` section 8 said prepare freezes the complete
  old/new sets including authority-profile and Gate Definition content; task
  10 called that content immutable prepare inputs. The accepted closed input
  has no such content or digest commitment: `public-interface.md` section 3
  retains Phase 21 prepare plus predecessor ID/version and five resources.
  Internal content arrives through retained `activation.component`, and
  `schema.md` section 5 permits component/profile/definition inserts only in
  that retained component-staged Event family, not in `activation.prepared`.
  Two different valid future component payloads can follow the same prepare
  request, so that request cannot already freeze either exact payload.
- **Provenance:** Both conflicting prose clauses originated in Tech Lead
  planning commit `6563fb4b4b92402c9df7826c577041bcac4e65b3`
  (`Refs: phase-22/task-planning`) and were accepted with the pack in DEC-188.
  No earlier inspected decision resolves the timing contradiction. The source
  audit is `qa/task-10-prepare-contract-audit.md`, committed independently at
  `8a4ff6d`; no implementation failure or recovery attempt caused this finding.
- **Resolution:** Prepare freezes the exact predecessor identity/version and
  the five declared old/new resource identities. The predecessor identifies
  its immutable internal component set; first activation has none. New
  authority-profile/Gate content is separately validated, content-addressed and
  made immutable when staged under the same new activation ID. Commit binds
  the exact two staged internal component digests and required resource
  agreement before the existing point of no return. No later staging or
  recovery rewrites prepare bindings; finalization/recovery cannot mix old
  and new internal sets. Normal workflow stays blocked while incomplete.
- **Why this resolution:** It reconciles prose with the already explicit
  closed wire and Event protocol, including Phase 21's retained
  prepare/component/commit boundaries, rather than inventing an unapproved
  payload, hidden resolver or provider-file read. This is an explicit
  authority-backed correction of contradictory prose, not a claim that the
  original pack was already unambiguous. It changes no accepted input field,
  DDL, Event family/Effect allocation, recovery direction or selected trust.
- **VAL-040 strength:** Keep its row and every VAL-041..044 row byte-identical.
  Task 10 now explicitly tests first activation and upgrade across their real
  operation boundaries: predecessor/resource freeze, separately staged
  immutable internal contents, exact selected digests and rejection of
  missing/mismatched/mixed sets. Resource binding and internal immutability
  are both required; neither is waived by moving the impossible payload
  expectation out of the prepare operation. Existing task gates and QA window
  4/checkpoint execute and independently review that evidence.
- **Spec impact:** Correct `contracts.md` section 8 and task 10's binding
  constraint/acceptance criterion. `public-interface.md` section 3 and
  `schema.md` section 5 already specify these operation boundaries and need no
  normative-row edits. Canonical `03_TRD_SCHEMA.md` section 8.3 incorporates
  the corrected slice contract; `07_DESIGN_PLAN_FLOWS.md` section 2 already
  separates prepare from staging. No separate canonical timing clause needs
  changing. This bounded inspection is not the phase-end full spec audit.
- **Limits and evidence:** No product code, tests, QA contract, validator,
  command, threshold, tolerance, dependency or workflow control changes.
  No gate was run and no VAL result is claimed. The resolution clears this
  pre-dispatch specification contradiction; all implementation and independent
  QA evidence remain task 10/11 and window 4 obligations.

## DEC-199 — Allocate Activation Error Enrichment To Its CLI Task

- **Date:** 2026-09-07.
- **Status:** ACCEPTED by the root delegate under the Human's explicit
  execution/scope-repair authority in `delegate-handoff.md`.
- **Finding:** Task 11 owns exact activation CLI errors, including the accepted
  `activation-incomplete/4` details in public-interface section 4. Task 10's
  resource-agreement rejection currently carries only a status and expected
  statuses. It cannot supply the required activation ID, disagreeing component
  list and Reasons without bounded domain error enrichment. Putting the
  agreement decision into the CLI would violate the existing layer boundary.
- **Resolution:** Add `internal/workflow/activation.go`,
  `internal/workflow/activation_recovery.go`, optional extracted
  `activation_errors.go` and matching tests to task 11 for typed failure data
  only. The Developer records task-10 introducing-commit provenance after its
  hand-back. Keep operation predicates, complete resource agreement, lifecycle,
  trust, replay, DDL and Events unchanged. New errors or changed error shapes
  are not authorized; this executes the already accepted closed contract.
- **Human rendering:** Make the already assigned Human status/recovery
  obligation explicit in task 11. Paired views name restore-prior before commit
  and complete-forward after it. Task 18 supplies retained fixtures and audits
  completeness; it does not first implement this adapter behavior.
- **Limits:** No VAL row, validator, gate command, accepted status, threshold,
  tolerance, dependency or workflow control changes. Task 10 remains active
  until its own unchanged gates/report close; this amendment does not advance
  task 11 or claim any implementation evidence.

## DEC-200 — Close The Reason Subject Allocation

- **Date:** 2026-09-07.
- **Status:** ACCEPTED by the root delegate under the Human's explicit
  authority to complete and repair remaining execution. This is an explicit
  specification correction before task-11 finalization, not Developer license
  to invent private enum values.
- **Gap:** Phase 22 retains the Phase-21 Reason structure and explicitly
  extends ReasonCode for Stop/failure codes, but omits the corresponding
  ReasonSubjectType additions. The accepted Stop and activation-incomplete
  error contracts require Reasons about those subjects; their typed identity
  domains already exist in this phase. Task 9 emits Stop Reasons, while task
  11 initially invented activation-resource/resource-not-observed-old to name
  a disagreeing resource. Neither invented value is admitted by the closed
  interface.
- **Resolution:** Explicitly retain the Phase-21 ReasonSubjectType set and add
  exactly `stop` and `activation`. Stop Reasons name the Stop ID; activation
  disagreement Reasons use the already admitted `activation-incomplete` code,
  Activation subject and Activation ID. The existing typed `components` list
  identifies disagreeing resources; explanatory detail can name the resource
  without introducing an entity/subject kind. Reason prose never determines
  action eligibility or resolution.
- **Execution:** Amend public-interface sections 2/3 and task 11. Root's shared
  fake/fid corpus asserts closed Reason fields and the exact activation
  subject; task 11 updates its typed error carriers under DEC-199. Keep the
  previously captured failing receipts. Review 4 and final conformance audit
  independently check the allocation; task 9's Stop subject semantics stay
  unchanged and its report/history remain closed.
- **Limits:** No new product operation, payload key, Reason structure, error
  code, DDL, Event Effect, authority predicate, validation assertion, validator,
  gate command, threshold, tolerance, dependency or workflow control. This
  closes an omitted enum allocation; it does not admit activation-resource or
  resource-not-observed-* and does not weaken any failure condition.

## DEC-201 — Complete Recovery Read Projections

- **Date:** 2026-09-07.
- **Status:** ACCEPTED by the root delegate under the Human's explicit
  authority to complete remaining execution and resolve omitted allocation.
- **Budget gap:** The accepted public-interface Budget amendment requires
  derived output `state`, but task 7 delivered arithmetic/admission without
  that serialized member. Task 12 needs the same Budget projection. Transfer
  execution of the omitted output member to task 12 with task-7 provenance;
  task 7 and its historical evidence stay closed.
- **Budget resolution:** Every public Budget output includes the accepted
  `unknown`, `available`, `warning` or `exhausted` state. Derive from the
  existing sample-chain counters and admission semantics: unknown counters
  mean unknown; zero remaining means exhausted; positive remaining at or
  below a non-null warning threshold means warning; otherwise available.
  Authoring never accepts the output member, including empty or null values.
  Immutable policy JSON, request identity and Event Effects retain their
  accepted shapes; output/idempotency result snapshots include the derived
  output. Do not refresh a historical replay snapshot using later counters.
- **Lease gap and resolution:** Contracts section 9 requires effective lease
  status but the retained Resume Claim only names durable `status`. Add
  exactly `effective_status` to the non-null Resume Claim projection. Its
  value is the durable status except that an active Claim whose expiry has
  passed projects `expired`. Retain `status` and `expires_at` verbatim and
  keep the existing null Claim case. This is a derived read member, never a
  Claim mutation, renewal or expiry Event. No other Claim payload gains it.
- **Execution:** Task 12 may edit the narrow Budget output/derivation and
  Envelope result paths in envelope.go, budget.go, run_read.go, run_budget.go,
  internal/evidence/execution.go and matching tests, plus required retained
  corpus convergence. Resume/Human projections use the same semantic facts.
  Reject every caller-supplied state member, preserve unknown-field rejection,
  and prove all four output states, immutable raw policy and replay snapshots.
  Root owns Human rendering; Developer owns domain output and causal tests.
- **Limits:** No arithmetic, admission, authority, persistence schema, Event
  Effect, budget threshold, new operation, VAL criterion, gate command,
  dependency or workflow control changes. Required task-12 gates and the
  independent checkpoint remain unchanged.

## Planning Spec Impact

| Document | Planning change |
|---|---|
| GLOSSARY.md | Flags proposed executor-observation, abandonment, and budget-sample terminology under DEC-180. |
| 02_TRD.md | Proposed Phase 22 package/contract index and pinned capability boundary. |
| 03_TRD_SCHEMA.md | Section 8.3 incorporates the proposed exact DDL/Event amendment, without rewriting Phase 21 historical SQL. |
| 04_TRD_CONFIG.md | Proposed CF-035 consuming-path interpretation, pending acceptance. |
| 06_DESIGN_PLAN_CLI.md | Proposed Phase 22 wire-contract pointer. |
| 10_IMPLEMENTATION_ROADMAP.md | Phase 22 planning state and conditional Claims qualification owner. |

This is not a phase-end per-spec verdict. The complete canonical suite,
GLOSSARY, every in-phase deviation, and every QA recommendation must be audited
at wrap-up. No `verified-no-change` assertion is being defaulted during kickoff.

## Triage-Log Audit

No TRG or ESC has been allocated for this phase. The seeded triage log is the
allocation source. Audit every actual entry at wrap-up; do not fabricate events
for pending planning approval or ordinary QA preparation.

## Workflow Notes

Reserved for observed phase execution and final wrap-up. Planning calibration
and deferred deep retrospective work are recorded in DEC-181/183 above.
