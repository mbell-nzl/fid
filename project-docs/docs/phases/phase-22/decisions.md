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
