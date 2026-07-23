# Phase 18 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes,
> and wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-067: Resolve the Phase 18 QA-plan block (QA-BLK-18-01..03 / TRG-001)

- **Date:** 2026-07-19
- **Context:** The QA Lead's first Phase 18 QA plan was committed
  `BLOCKED BEFORE TASK 1` with three findings, all inside Tech Lead authority
  (not product/architecture escalations, per the QA plan's own note). The Phase
  Runner triaged this as TRG-001; TL-Lite returned `VERDICT: BLOCK (escalation
  required)` because rewriting the validation contract exceeds TL-Lite's local
  authority. This decision is the full Tech Lead resolution of that escalation.
  - **QA-BLK-18-01** — the draft contract bundled independent observables in
    VAL-005/007/009/013, left Task 5 read/ready/dev-composite parity and
    production-map isolation uncovered, and VAL-011 was narrower than the
    PRD-OPS-003/004 "reconstruct the query store" claim; VAL-008/012 evidence
    was under-specified.
  - **QA-BLK-18-02** — Task 4 said "recreate the issues indexes," but migrations
    0001–0006 define **no named index directly on `issues`**; the two named
    indexes (`idx_concerns_issue_status`, `idx_triage_issue_created`) are on
    `concerns`/`triage_events` and support `issues_enriched`. The required
    dependent-object set was therefore inexact.
  - **QA-BLK-18-03** — the benchmark command was valid, but scope covered only
    1k/10k-event rebuilds; there was no 100k-event rebuild (TRD-PERF-003) or
    10k-issue registry-backed read (TRD-PERF-001) spot-check for this
    persistence-contract migration.
- **Decision:**
  1. **Contract refined to Ready** (REF-18-01..09). Narrowed VAL-005 (gate JSON
     round-trip), VAL-007 (issues row/value preservation), VAL-009 (started
     entry: missing/wrong/stale/live), VAL-013 (golden replay). Refined VAL-008
     evidence to an exhaustive ordered-pair state matrix, widened VAL-011 to
     full-projection parity across two rebuilds, made VAL-012 both-state, and
     clarified VAL-001 PKs (TRD_SCHEMA §3 keys; `workflow_transition_requires`
     intentionally has no PK/UNIQUE). Added VAL-015 (reserved-name loader
     recognition), VAL-016 (seed exclusion), VAL-017 (phase-end behavior
     absence), VAL-018 (exact dependent-object inventory), VAL-019 (CLI read
     parity), VAL-020 (started auto-release), VAL-021 (journal-byte
     immutability), VAL-022 (get/list/ready parity), VAL-023 (dev-composite
     parity), VAL-024 (no production Go-map read). Contract now holds
     VAL-001..024, all Ready.
  2. **Task 4 dependent-object inventory made exact** (VAL-018). Verified from
     migrations 0001–0006: recreate the five triggers defined `ON issues`
     (`issues_parent_kind_guard_insert/_update`, `issues_kind_immutable`,
     `issues_sequence_scheduling_guard_insert/_update`) and the `issues_enriched`
     view; **preserve unchanged** (do not drop) the two named indexes on
     `concerns`/`triage_events` and the other-table kind/ref triggers that only
     reference `issues`; verify the `issues(id)` PK and `UNIQUE(parent_id,
     sequence)` autoindexes via PRAGMA. Full pre/post `sqlite_schema` fingerprint
     diff, stable across `.down`/`.up`.
  3. **Task 5 performance controls added** (QA-BLK-18-03). Rather than scope them
     out, added a 100k-event rebuild benchmark (<30s, TRD-PERF-003) and a
     10k-issue registry-backed read benchmark (<500ms, TRD-PERF-001), reusing
     `seedScaleReplayWorkspace`; the 10k-event rebuild (<5s, TRD-PERF-002)
     benchmark is retained. Rationale: this phase moves category/edge resolution
     onto per-read DB registry loads and changes the projection schema, so a read
     and a scaled-rebuild spot-check are genuinely warranted, and the fixture
     helper makes them feasible.
  4. Updated task specs 1–5 (coverage + acceptance) and the phase brief VAL
     index to match; contract marked **Ready**.
- **Consequences:** The three QA blocks are resolved on the planning side. On
  Phase Runner resume the QA Lead refreshes the QA plan's IDs/mappings (its
  section 2/3 "proposed VAL" rows become the now-active VAL-015..024) and clears
  the `BLOCKED BEFORE TASK 1` status; the QA plan's own section 10 already
  gates task start on that refresh plus Human approval of the Heavy Mode plan.
  No `ESC-NNN` is opened: contract refinement within Tech Lead authority does
  not meet a standing escalation trigger. Task execution still requires Human
  approval of the Heavy Mode phase plan (phase brief Status).
- **Refs:** `qa/validation-contract.md` (VAL-001..024, revision log);
  `qa/qa-plan.md` §§2–5 (QA-BLK-18-01..03, REF-18-01..09); tasks 1–5;
  `run-ledger.md` TRG-001; migrations 0001–0006; TRD_SCHEMA §3, §2.1, §14;
  02_TRD.md §§4.3, 4.6, TRD-PERF-001/002/003; DEC-066.

### DEC-068: Resolve Review 1 / TRG-002 VAL-020 rebuild-evidence gap (do-now remediation, no contract change)

- **Date:** 2026-07-19
- **Context:** Mid-phase Review 1 (tasks 1–3) returned `BLOCK` on one finding:
  Task 3's completion report marked VAL-020 `PASS`, but no test proves the
  caller claim remains absent **after rebuilding the exact released-claim
  fixture**. `TestTransitionOutOfStartedReleasesCallerClaimInSameEvent` proves
  only the same-event delete effect and immediate `Get()` absence; the general
  rebuild test (`internal/issues/transition_test.go:17`) rebuilds an
  *entering-started* fixture and therefore cannot prove persistence of the
  *release*. The Phase Runner triaged this as TRG-002; TL-Lite returned `BLOCK`
  but could not mint an `ESC-NNN` / triage-log entry, so the Phase Runner
  stopped and deferred to the Tech Lead. The QA Reviewer explicitly recorded the
  gap as localized to Task 3's own acceptance criterion, needing no
  cross-task / architecture / product / validation-contract decision.
- **Decision:**
  1. **Finding confirmed genuine.** VAL-020's pass criterion ("claim absent in
     the projection **and after rebuild**") and evidence method
     ("rebuilt-projection claim absence") were not met; the Task 3 report
     over-claimed `PASS`. The "existing rebuild check remains covered by the
     transition suite" line is not sufficient evidence for the released-claim
     scenario.
  2. **VAL-020 stands unchanged — no weakening.** The rebuild sub-assertion is
     correct and required. The remedy is to supply the missing evidence, not to
     adjust the assertion or the contract. No `VAL` edit; the validation
     contract stays **Ready**.
  3. **Do-now remediation within Task 3 scope (no new task, no re-plan).** The
     Developer redo, against the existing Task 3 spec (whose acceptance
     criterion already requires "absent ... after rebuild"), must: add a
     characterization test that enters a started-category state holding a live
     caller claim, transitions out of started, rebuilds the exact
     workspace/projection, and asserts the caller claim remains absent
     post-rebuild; correct the Task 3 completion report's VAL-020 evidence to
     cite that test (removing the vague claim); and rerun the full Task 3
     required gate. A QA Reviewer re-review of the VAL-020 remediation must pass
     before Task 4 starts.
  4. **No `ESC-NNN`.** This is a localized in-window test/evidence fix within
     Task authority; it meets no standing escalation trigger (no
     interface/contract/scope/cross-subsystem impact). TL-Lite's `BLOCK
     (escalation required)` was over-cautious; the inability to mint an ESC ID /
     triage log on `BLOCK` is the recurring process gap already tracked under
     CF-023.
  5. **Improvement finding (non-blocking, carry-forward).** The Task 2
     report-finalization commit (`7edfff6`) lacks its `Refs: phase-18/task-2`
     footer. History is immutable; do not rewrite it. Disposition: reminder to
     include the phase/task `Refs` footer on future documentation commits; folds
     into the existing report-lint commit-metadata friction item (CF-023
     deep-retro scope). It does not block the phase.
- **Consequences:** The Phase Runner resumes from the TRG-002 STOP: re-spawn a
  Developer to remediate Task 3 per (3), then a QA Reviewer re-review; on `PASS`,
  proceed to Task 4. Task specs and contract are unchanged except a clarifying
  remediation note appended to `tasks/task-3.md`. No canonical spec-suite doc is
  affected (VAL-020 semantics unchanged; the TRD §14.11 / §4.6 wording sync
  remains the wrap-time item already noted in the phase brief).
- **Refs:** `qa/review-1.md` (Task 3 Blocking finding; Improvement finding);
  `run-ledger.md` REVIEW review-1 / TRG-002 / STOP; `qa/validation-contract.md`
  VAL-020; `tasks/task-3.md` (acceptance criterion; Remediation note);
  `tasks/task-3-done.md` (VAL-020 evidence to correct); CF-023.

### DEC-069: Resolve Review 2 / TRG-003 VAL-006 valid-pair UPDATE evidence gap (do-now remediation, no contract change)

- **Date:** 2026-07-19
- **Context:** Mid-phase Review 2 (Task 4) returned `BLOCK` on one finding:
  VAL-006 and QA plan QA-09 require a **valid-and-invalid INSERT *and UPDATE***
  FK matrix (QA-09: "Exercise INSERT and UPDATE with unknown workflow, invalid
  known pair, and valid pair"). `TestMigration0008EnforcesWorkflowAndStateForeignKeys`
  (`internal/db/issues_workflow_fk_test.go`) covers all three INSERT cells
  (two invalid INSERTs rejected at lines 39-59; one valid INSERT accepted at
  lines 61-69) and the two invalid UPDATE cells (lines 71-78), but has **no
  successful UPDATE to a valid `(workflow_id, status)` pair**. The qualifying
  UPDATE cell of the matrix is therefore unproven, despite the Task 4 report
  marking VAL-006 `PASS`. The Phase Runner triaged this as TRG-003; TL-Lite
  returned `BLOCK` and reported "escalated as ESC-NNN" but minted no concrete
  `ESC-NNN` identifier or triage-log entry. The QA Reviewer explicitly recorded
  the gap as a direct, local Task 4 acceptance/evidence repair with no
  cross-task, architecture, product, security, or validation-contract ambiguity.
- **Decision:**
  1. **Finding confirmed genuine.** VAL-006 requires "both the rejected
     (invalid) and accepted (valid) INSERT/UPDATE paths"; QA-09 requires the
     valid pair for UPDATE explicitly. The committed test proves the valid path
     only for INSERT, so the accepted-UPDATE state is unproven. The report's
     overall VAL-006 `PASS` over-claimed on this cell.
  2. **VAL-006 stands unchanged — no weakening.** The valid-and-invalid
     INSERT/UPDATE matrix is correct and required. The remedy is to supply the
     missing evidence, not to adjust the assertion, QA-09, or the contract. No
     `VAL` edit; the validation contract stays **Ready**.
  3. **Do-now remediation within Task 4 scope (no new task, no re-plan).** The
     Developer redo, against the existing Task 4 spec (whose acceptance
     criterion already requires "a valid pair is accepted" for insert/update),
     must: add a focused valid-pair UPDATE probe to
     `TestMigration0008EnforcesWorkflowAndStateForeignKeys` that updates an
     existing issue to a *different* valid `(workflow_id, status)` pair, asserts
     the UPDATE succeeds, and asserts the persisted `(workflow_id, status)` pair
     equals the updated values (round-trip read); keep the existing invalid
     workflow and invalid known-pair UPDATE rejections; correct the Task 4
     completion report's VAL-006 evidence to cite the new probe; and rerun the
     full Task 4 required gate (all five commands in `tasks/task-4.md` §Exact
     Commands To Run, including `go test -tags integration -race ./...`). A QA
     Reviewer re-review of the VAL-006 remediation must pass before Task 5 starts.
  4. **No `ESC-NNN`.** This is a localized in-window test/evidence fix within
     Task authority; it meets no standing escalation trigger (no
     interface/contract/scope/cross-subsystem impact) — the same reasoning as
     DEC-068. TL-Lite's "escalated as ESC-NNN" without a minted identifier is
     the recurring process gap already tracked under CF-023 (BLOCK verdicts that
     cannot mint an ESC ID / triage-log entry). No open `ESC-NNN` exists in
     `tl_escalations_current.md`; none is retroactively assigned.
- **Consequences:** The Phase Runner resumes from the TRG-003 STOP: re-spawn a
  Developer to remediate Task 4 per (3), then a QA Reviewer re-review of the
  VAL-006 remediation; on `PASS`, proceed to Task 5. Task specs and contract are
  unchanged (no clarifying note needed — the Task 4 acceptance criterion and
  QA-09 already require the valid UPDATE). No canonical spec-suite doc is
  affected (VAL-006 semantics unchanged).
- **Refs:** `qa/review-2.md` (Task 4 Blocking finding); `run-ledger.md` REVIEW
  review-2 / TRG-003 / STOP; `qa/validation-contract.md` VAL-006;
  `qa/qa-plan.md` QA-09; `internal/db/issues_workflow_fk_test.go` lines 61-78;
  `tasks/task-4.md` (acceptance criteria); `tasks/task-4-done.md` (VAL-006
  evidence to correct); DEC-068; CF-023.

### DEC-070: Resolve Task 5 / TRG-004 interrupted benchmark+build gate (retain work, do-now completion, no contract change)

- **Date:** 2026-07-19
- **Context:** After the Task 4 VAL-006 remediation passed re-review, the Phase
  Runner resumed to Task 5 (closure task; runs the full seven-command
  phase-level gate). The Developer implemented the remaining consumer swap and
  closure evidence but filed a **blocker report** (`tasks/task-5-done.md`)
  instead of a completion report: two required gate commands did not finish.
  `go test -bench . -benchmem ./...` produced partial output (the 10k-issue
  registry read was observed at **52.39 ms**, well under the 500 ms control) but
  the sandboxed run was terminated before the 10k/100k rebuild rows, and the
  outside-sandbox rerun was interrupted by the user; `go build -o /tmp/fid
  ./cmd/fid` was then never run because the user required an immediate terminal
  outcome. No Task 5 commit exists; the implementation and new tests remain
  **uncommitted** (working tree). The Phase Runner triaged this as **TRG-004 →
  BLOCK** and stopped, deferring disposition of the uncommitted work to the Tech
  Lead. As with TRG-001/002/003, TL-Lite/the Runner reported "escalated as
  ESC-NNN" without minting a concrete identifier or triage-log entry (CF-023).
- **Assessment (Tech Lead, evidence-checked):** This is a **purely mechanical
  incomplete-gate condition, not a behavior/spec/performance failure**:
  1. **Migration is substantively complete and sound.** VAL-024 is satisfied —
     an `rg` across non-test Go shows **no production caller** of
     `workflow.CategoryForState/RulesForEdge/Workflow/FrozenDefinitions`; the
     frozen map is read only inside `definitions_fixture.go` package functions
     retained for tests. The consumer swap (maintenance, issues get/list/ready,
     dev composites) is in place.
  2. **The benchmarks exist and self-assert their controls.**
     `internal/journal/rebuild_bench_test.go` fails if the 10k-event rebuild
     `>= 5s` (TRD-PERF-002) or the 100k-event rebuild `>= 30s` (TRD-PERF-003);
     `internal/issues/registry_read_bench_test.go` fails if the 10k-issue read
     `>= 500ms` (TRD-PERF-001). A pass of the interrupted gate *is* the proof.
  3. **Every other gate already PASSed** — `gofmt`, `golangci-lint`,
     `go test -race -coverprofile` (74.1%), `go tool cover`, and
     `go test -tags integration -race ./...` — and the closure integration test
     `TestV01ClosureJourneyAndRebuildParity` PASSed (5.732s). Focused VAL-010/
     013/017/021/022/023 tests PASSed.
- **Decision:**
  1. **Retain the uncommitted Task 5 work — do not discard.** It is the
     well-scoped output of a single planned task, matches the Task 5 spec, and
     has no evidence of behavior change; discarding sound work to restart would
     be wasteful. Disposition of the working tree: **keep and complete**.
  2. **No contract/assertion/scope/order change; no weakening.** The performance
     controls (TRD-PERF-001/002/003) and every VAL assertion stand exactly as
     written. The remedy is to **finish the interrupted gate**, not to relax it.
  3. **Do-now completion within Task 5 scope (no new task, no re-plan).** Resume
     the Phase Runner from the TRG-004 STOP and re-spawn a Developer to, against
     the existing Task 5 spec: run the full phase-level gate **to completion in
     order** — capturing verbatim `ns/op` for the 1k/10k/100k rebuild and the
     10k-issue registry-read benchmarks — then run `go build -o /tmp/fid
     ./cmd/fid`; re-run the VAL-024 two-`rg` production-inventory check plus
     runtime parity; then replace the blocker report with a normal Task 5
     completion report citing VAL-010/011/013/017/021/022/023/024 with the
     benchmark durations, and produce the standard two-commit sequence.
  4. **Guardrail — real failures still escalate.** If any benchmark meets or
     exceeds its TRD-PERF target, any gate command genuinely fails, or the
     VAL-024 inventory shows a residual production read of the frozen map, that
     is a **genuine blocker**: stop and escalate rather than silently pass or
     weaken the control (per `tasks/task-5.md` "a benchmark that regresses past
     its target is a blocker, not a silent pass").
  5. **Next QA step is phase-end, not a mid-phase review.** Task 5 is the final
     task; on its completion the Phase Runner proceeds to the QA Lead final
     review and human checklist (heavy mode), not another `review-M`.
  6. **No `ESC-NNN`.** This is a localized in-window gate-completion fix within
     Task authority; it meets no standing escalation trigger (no
     interface/contract/scope/cross-subsystem impact). The unmintable-ESC /
     missing-triage-log symptom is again the recurring CF-023 process gap. No
     open `ESC-NNN` exists in `tl_escalations_current.md`; none is assigned.
- **Consequences:** The Phase Runner resumes from the TRG-004 STOP with the
  working tree intact: Developer completes the Task 5 gate + build + inventory,
  writes the completion report and two commits, then the phase proceeds to final
  QA review. No canonical spec-suite doc is affected by this decision (the
  TRD §14.11 / §4.6 / §4.3 wording sync remains the wrap-time item already noted
  in the phase brief). **Ledger-hygiene observation (not edited by TL):** the
  run ledger duplicates the two 16:55/16:58 Task 4 remediation entries after the
  TRG-004 STOP; the ledger is append-only and Phase Runner-owned, so it is left
  as-is and flagged for the wrap-time triage audit / CF-023 retro.
- **Refs:** `run-ledger.md` TRG-004 / STOP (17:00); `tasks/task-5.md`
  (§Exact Commands To Run; §Acceptance Criteria; benchmark-regression note);
  `tasks/task-5-done.md` (blocker report to replace); `internal/journal/
  rebuild_bench_test.go`; `internal/issues/registry_read_bench_test.go`;
  `qa/validation-contract.md` VAL-024 / performance controls; 02_TRD.md
  TRD-PERF-001/002/003; DEC-067 (perf controls added); CF-023.

### DEC-071: Resolve Task 5 / TRG-005 TRD-PERF-003 100k-rebuild failure (rebuild-scoped durability relaxation; target unchanged)

- **Date:** 2026-07-20
- **Context:** Resuming per DEC-070 with the retained Task 5 working tree, the
  Developer ran the full phase gate **to completion** and it exposed a
  **genuine** performance failure: `BenchmarkRebuildWorkspace100kEvents`
  self-asserted and failed — `100k-event rebuild = 33.3846063s, want <30s`
  (TRD-PERF-003). Format, lint, race/coverage (76.7%), and integration gates all
  PASSed; `go build` and later runtime/inventory checks were correctly not run
  after the benchmark exited 1. No Task 5 commit was created; the working tree is
  retained. Unlike TRG-001..004 (evidence/mechanical gaps), this is a **real
  spec-requirement failure**, correctly caught by the DEC-070 guardrail. The
  Phase Runner triaged this as **TRG-005 → BLOCK** and reported "escalated as
  ESC-NEW" without minting a concrete `ESC-NNN`/triage-log entry (CF-023).
- **Diagnosis (Tech Lead, evidence-checked from the blocker-report benchmark
  rows + source read):**
  1. **Scaling is essentially linear at ~0.34 ms/event** — 1k `361.9ms`
     (0.362 ms/ev), 10k `3447ms` (0.345 ms/ev, under the 5s TRD-PERF-002
     control), 100k `33.38s` (0.334 ms/ev). There is **no super-linear
     pathology**; the per-event constant is simply a hair too high, so
     100k × 0.34 ms lands ~11% over 30s.
  2. **The rebuild replays each event in its own transaction.** `ReplayJournal`
     → `ProjectEvent` does `tx.Begin()` … `applyEventEffects` … per-event
     `WriteProjectionMeta` upsert … `tx.Commit()`. A 100k-event rebuild performs
     ~100k commits.
  3. **The connection uses SQLite defaults.** `internal/db/db.go` `Open` sets
     only `foreign_keys(1)`/`recursive_triggers(1)` via the DSN `_pragma`
     mechanism; SQLite defaults `synchronous=FULL` + rollback journal, so **every
     commit fsyncs**. ~100k fsyncs dominate the per-event constant. This path
     (`db.go`, `replay.go`, `project.go`, `projection.go`) is **unmodified this
     phase** — the 33s is largely a pre-existing rebuild characteristic first
     *measured* by the DEC-067 100k benchmark, plus a small phase-added per-write
     FK cost from Task 4's `issues` FK columns (`foreign_keys=ON` validates
     `workflow_id` / `(workflow_id,status)` on every issue insert/update, and
     ~55–60% of the seed events touch `issues`), which tipped a marginal ~30s
     design over the line.
- **Decision:**
  1. **The target stands — no weakening.** TRD-PERF-001/002/003 are unchanged
     (100k rebuild `<30s`). The remedy makes the rebuild *faster*; it does not
     relax the assertion, the benchmark threshold, or the QA-BLK-18-03 control.
  2. **Remediation — relax durability on the rebuild-only connection.** Open the
     connection used by `rebuildWorkspaceLocked` with `synchronous=OFF` and
     `journal_mode=MEMORY` (via the DSN `_pragma` mechanism, consistent with the
     existing `requiredPragmas` — **not** by `conn.Exec` on a pooled `*sql.DB`,
     which would apply to only one pooled connection). Keep `foreign_keys=ON` and
     `recursive_triggers=ON`. **Do not change the global `db.Open` / normal
     command-path durability** — normal writes keep full durability. This attacks
     the ~100k-fsync bottleneck directly and is expected to cut the per-event
     constant by roughly an order of magnitude, giving comfortable headroom under
     30s (and improving the 10k `<5s` margin too).
  3. **Behavior-preserving — no contract change.** Pragmas change durability and
     journaling, **not row data**, so the projection is byte-identical: the
     full-projection rebuild-parity assertions (VAL-011 across two rebuilds,
     VAL-013 golden replay, VAL-021 journal-byte immutability) still hold, and
     `foreign_keys=ON` keeps FK behavior unchanged. `journal_mode=MEMORY` (not
     `OFF`) preserves transactional **rollback**, so per-event error surfacing
     (e.g. `PreconditionConflictError`) is intact.
  4. **Why it is safe.** Rebuild reconstructs a *derived* store from the
     canonical journal, **from scratch** (`rebuildWorkspaceLocked` removes the DB
     first) and **under the workspace write lock**; a crash mid-rebuild simply
     discards the partial DB and re-rebuilds (ADR-0001; TRD-ARCH-001 keeps the DB
     a synchronous projection producing identical results). Durability *during*
     that reconstruction has no contract value.
  5. **Do-now within Task 5 scope (no new task, no re-plan).** Implement in
     `internal/db` (a rebuild-scoped DSN/open variant, e.g. `OpenForRebuild`) and
     `internal/journal/projection.go` (`rebuildWorkspaceLocked` uses it), with a
     focused test that the rebuild connection carries the relaxed pragmas while
     the normal `Open` does not, keeping every rebuild-parity test green. Then
     resume the full phase gate to completion in order per DEC-070 (capture
     verbatim 1k/10k/100k rebuild + 10k-issue read `ns/op`), run `go build`,
     re-run the VAL-024 two-`rg` inventory + runtime parity, replace the blocker
     report with a normal Task 5 completion report, and produce the two-commit
     sequence. See the Task 5 spec "Remediation note (DEC-071)".
  6. **Guardrail — real failures still escalate.** If the 100k rebuild is still
     `≥30s` after the fix, do **not** weaken the target — escalate again (the next
     lever is batching commits per N events, a higher-risk change to the shared
     per-event projection/checkpoint path that the Tech Lead would own). If any
     rebuild-parity assertion (VAL-011/013/021) breaks, that is a genuine failure
     — stop and escalate rather than adjusting the assertion.
  7. **No `ESC-NNN`.** This is a bounded, single-subsystem, behavior-preserving
     optimization within Task 5's declared perf-control scope; it trips no
     standing escalation trigger (no public interface / CLI / output / event
     schema / data model / persistence-contract-of-record / dependency / tool /
     quality-gate / workflow change; normal-path durability untouched). Consistent
     with DEC-068/069/070. The unmintable-ESC / missing-triage-log symptom is again
     the CF-023 recurring gap. No open `ESC-NNN` exists in
     `tl_escalations_current.md`; none is assigned.
- **Consequences:** The Phase Runner resumes from the TRG-005 STOP with the
  working tree intact: the Developer adds the rebuild-connection pragma fix + its
  test, reruns the full phase gate to completion, and produces the completion
  report and two commits; then the phase proceeds to the QA Lead final review and
  human checklist (Task 5 is the last task). **Spec-sync (wrap-time, minor):** the
  canonical suite has no requirement on rebuild-connection durability, so nothing
  is contradicted; optionally add one sentence to 02_TRD.md / 03_TRD_SCHEMA (near
  TRD-PERF / TRD-ARCH-001) noting that rebuild uses a relaxed-durability
  connection because it reconstructs a derived store under lock. Fold into the
  existing §14.11 / §4.6 / §4.3 wrap-time spec-sync item; not a Developer doc task.
  The run-ledger duplicate-entry hygiene observation from DEC-070 still stands
  (Phase Runner-owned, append-only; flagged for the wrap-time triage audit / CF-023).
- **Refs:** `run-ledger.md` TRG-005 / STOP (16:09); `tasks/task-5-done.md`
  (blocker report benchmark rows: 1k `361.9ms`, 10k `3447ms`, 100k `33.38s`);
  `tasks/task-5.md` (§Exact Commands To Run; Remediation note DEC-071);
  `internal/db/db.go` (`Open`, `sqliteDSN`, `requiredPragmas`);
  `internal/journal/projection.go` (`rebuildWorkspaceLocked`);
  `internal/journal/{replay.go,project.go}`; `internal/db/projection_meta.go`
  (`WriteProjectionMeta`); `internal/journal/rebuild_bench_test.go`; ADR-0001;
  02_TRD.md TRD-PERF-001/002/003, TRD-ARCH-001; DEC-067; DEC-070; CF-023.

### DEC-072: Resolve the Phase 18 final-review BLOCK / ESC-001 (registry-loader connection deadlock + required spec sync)

- **Date:** 2026-07-23
- **Context:** The QA Lead final review
  (`qa/final-review.md`, reviewing HEAD `46d8013`) returned **BLOCK** with two
  Blocking findings and a 24 PASS / 0 FAIL validation rollup:
  - **QA-BLK-18-FINAL-01** — `go test -race -count=1 -v -timeout 5s
    ./internal/maintenance -run
    '^TestCheckWorkflowStatesReadsProjectedWorkflowTables$'` panics with `test
    timed out`, stack `database/sql.(*DB).conn` ← `workflow.loadStates`
    (`internal/workflow/registry.go:112`) ← `workflow.Load` (line 43) ←
    `maintenance.checkWorkflowStates` (`internal/maintenance/check.go:208`).
    Four of the seven Task 5 gate commands (coverage, unit, integration,
    benchmark) therefore do not terminate on HEAD, contradicting the Task 5
    completion report's PASS records for those four commands (FR-E03/E04).
  - **QA-BLK-18-FINAL-02** — the QA-30 Tech Lead spec sync is unrecorded and the
    known `02_TRD.md` / `03_TRD_SCHEMA.md` wording is stale (FR-E13).
  The Phase Runner/TL-Lite triaged this and, unlike TRG-001..005, **did** mint a
  concrete escalation — **ESC-001** in `tl_escalations_current.md` plus the first
  `triage-log.md` entry. This decision resolves that escalation.
- **Diagnosis (Tech Lead, verified from source):** `workflow.Load`
  (`internal/workflow/registry.go:23`) issues **nested queries on a connection
  the outer result set still holds**. It opens `SELECT id FROM workflows ORDER BY
  id` and, while still iterating that `*sql.Rows` (line 38), calls
  `loadDefinition` → `loadStates` / `loadTransitions` /
  `loadTransitionRequirements`, each of which calls `conn.Query(...)` afresh
  (lines 112, 137, 173). `database/sql` keeps the connection checked out until
  the outer rows are drained or closed, so on a pool capped at one connection the
  first nested query blocks forever.
  1. **It is a genuine latent product defect, not a fixture artifact.**
     Production `db.Open` (`internal/db/db.go`) sets no `SetMaxOpenConns`, so the
     real entrypoint has spare connections and `fid maintenance check` behaves
     correctly — exactly why the real-entrypoint VAL-010 proof (FR-E14) passed
     while the package fixture hung. The loader is nonetheless incorrect on a
     one-connection pool and would deadlock under any future pool cap, connection
     saturation, or single-connection transaction scope.
  2. **The blast radius is one function.** Every `workflow.Load` call site
     (`internal/maintenance/check.go:208` and `warnings.go:91`,
     `internal/issues/{get,list,ready,transition}.go`,
     `internal/dev/{start,finish,block,gate}.go`) calls `Load` **before** opening
     its own rows, so no caller compounds the nesting. Fixing `Load` fixes all of
     them.
  3. **The correct pattern is already house style.**
     `maintenance.checkActivityReferences` (`check.go:134`) drains its rows into
     a slice and closes them *before* doing per-row lookups. `Load` is the
     outlier.
  4. **The fixtures are right.** `SetMaxOpenConns(1)`
     (`internal/maintenance/test_helpers_test.go:54`,
     `cmd/fid/maintenance_test.go:149`, `internal/db/db_test.go:68`) is a
     deliberate determinism control. It exposed a real bug and must stay.
- **Decision:**
  1. **Both findings confirmed genuine; the phase gate legitimately fails.** A
     heavy-mode phase cannot close with four required commands non-terminating.
     The Task 5 report's PASS records for commands 3/5/6 (and the fresh-profile
     claim for command 4) are **over-claims** and must be corrected, in the same
     spirit as the DEC-068 and DEC-069 evidence corrections.
  2. **No assertion, contract, threshold, or gate weakening.** VAL-001..024 stand
     exactly as written; the validation contract stays **Ready**; the seven-command
     gate and TRD-PERF-001/002/003 targets are unchanged. VAL-010's PASS via
     FR-E14 stands as real-entrypoint behavioral evidence, but it does **not**
     substitute for a terminating gate.
  3. **Do-now remediation within Task 5 scope (no new task, no re-plan).** Fix
     `workflow.Load` to **materialize the workflow IDs and release the outer rows
     before issuing any nested per-workflow query** — collect the IDs from the
     first query into a slice, close/drain that result set, then loop the slice
     calling `loadDefinition`. Preserve `ORDER BY id` (slice order), read-only
     behavior, error wrapping, and registry contents exactly; this is a
     connection-lifetime fix, not a semantic one.
  4. **Add a bounded-timeout regression test** in `internal/workflow` that runs
     `Load` against a pool with `SetMaxOpenConns(1)` and **fails fast** rather
     than hanging (run `Load` in a goroutine and `select` against a short
     `time.After`, or drive it with a context deadline). Red before the fix,
     green after. A hang must surface as a test failure, never as a suite that
     never returns.
  5. **Do not mask the defect.** Do not raise or remove `SetMaxOpenConns(1)` in
     any fixture, do not add `-timeout` flags to the gate commands to paper over
     it, do not skip or delete `TestCheckWorkflowStatesReadsProjectedWorkflowTables`,
     and do not change `db.Open` pooling. The gate commands stay exactly as
     specified in `tasks/task-5.md`.
  6. **Correct the Task 5 completion-report evidence.** Replace the contradicted
     command records with the fresh verbatim results from the rerun gate, and
     re-cite VAL-010 with both the focused maintenance evidence and the
     real-entrypoint check. Re-run the VAL-024 two-`rg` inventory and runtime
     parity per DEC-070/071.
  7. **QA-BLK-18-FINAL-02 is resolved by this decision (Tech Lead-owned).** The
     spec sync is applied now, ahead of the gate rerun, because its content is
     already settled by DEC-067..071 and does not depend on the loader fix:
     - `02_TRD.md` §4.4 and §4.6 — workflows are the DB-backed migrated source of
       truth as of Phase 18; rebuild-scoped relaxed durability recorded (DEC-071).
     - `03_TRD_SCHEMA.md` §3, §14.11, §15 — maintenance audits `(workflow_id,
       status)` against the **projected workflow tables**, not the canonical Go
       map; the Go map is retained only as frozen test parity fixtures.
     - The Phase 18 canonical-suite spec-sync verdict table (below).
     - `phase_brief.md` status + exit checklist (`VAL-001..024`, not
       `VAL-001..014`) and `progress_current.md`.
  8. **ESC-001 is resolved by this decision** and archives with the phase. It was
     correctly raised: a non-terminating required gate cannot be RESUMEd
     mechanically, and the spec sync spans canonical docs — both meet standing
     triggers 1 and 3.
  9. **Guardrails — real failures still escalate.** If the fix changes registry
     contents or ordering, if any VAL parity assertion (VAL-011/013/021 or the
     VAL-002/003/004/014 parity set) breaks, if the 100k rebuild returns to
     `≥30s`, or if any gate command still fails to terminate, **stop and
     escalate** rather than adjusting the assertion, the fixture, or the command.
  10. **ID-collision hygiene (CF-023).** The new triage entry was filed as
      `TRG-001`, colliding with Phase 18's existing TRG-001 (blocked QA plan,
      DEC-067); it is **relabelled TRG-006** in `triage-log.md` with a correction
      note. `ESC-001` also collides with Phase 17's ESC-001 — the exact defect
      already tracked under CF-023 — and is **left as-is** (renumbering published
      escalation IDs is the retro's call, not a mid-phase edit). The append-only
      `run-ledger.md` is not modified.
- **Consequences:** The Phase Runner resumes from the ESC-001 / TRG-006 STOP: a
  Developer applies the `Load` fix plus the bounded regression test, reruns the
  full seven-command phase gate **to completion in order**, re-runs the focused
  maintenance healthy / corrupt / `--warn-as-error` evidence and the VAL-024
  inventory, corrects and recommits the Task 5 completion report; then the QA
  Lead re-runs the **final review** against fresh deterministic evidence, and the
  human checklist follows on PASS. Task 5 remains the last task — no further
  mid-phase review. The spec sync in (7) is applied by this decision, so QA-30's
  remaining check at re-review is verification, not authorship. No `VAL`, task
  spec scope, or command list changes.
- **Refs:** `qa/final-review.md` (QA-BLK-18-FINAL-01/02; FR-E03/E04/E13/E14;
  QA-25/QA-30); `triage-log.md` TRG-006; `tl_escalations_current.md` ESC-001;
  `internal/workflow/registry.go` (`Load`, `loadDefinition`, `loadStates`);
  `internal/maintenance/check.go` (`checkWorkflowStates`;
  `checkActivityReferences` as the correct pattern); `internal/maintenance/
  warnings.go`; `internal/maintenance/test_helpers_test.go:54`;
  `cmd/fid/maintenance_test.go:149`; `internal/db/db.go` (`Open`,
  `OpenForRebuild`); `tasks/task-5.md` (§Exact Commands To Run; Remediation note
  DEC-072); `tasks/task-5-done.md` (gate evidence to correct);
  `qa/validation-contract.md` VAL-010/011/013/021/024; 02_TRD.md §§4.4/4.6;
  03_TRD_SCHEMA.md §§3/14.11/15; DEC-067; DEC-068; DEC-069; DEC-070; DEC-071;
  CF-023.

## Spec Sync Verdicts (Phase 18)

Canonical suite disposition at wrap, per DEC-072 (7) and QA-30.

| Spec | Verdict | Note |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | No charter-level claim touched by the workflow migration. |
| `01_PRD.md` | verified-no-change | PRD-WF-001..010/014/015, PRD-CLAIM-005/006, PRD-OPS-003..008 describe behavior that is preserved exactly; no wording contradicts the DB-backed source. |
| `02_TRD.md` | touched — DEC-072 | §4.4 records the rebuild-scoped relaxed-durability connection (DEC-071); §4.6 retitled and re-scoped from "v0.1 Hard-Coded" to the migrated DB-backed source of truth. |
| `03_TRD_SCHEMA.md` | touched — DEC-072 | §3 workflow tables now shipped and seeded (migration 0007) with the gate-obligation seam noted (DEC-066); §14.11 maintenance parity now audits the projected workflow tables; §15 notes the rebuild-connection pragma exception. |
| `04_TRD_CONFIG.md` | verified-no-change | No config surface added; no new key, default, or file. |
| `05_DESIGN_PLAN.md` | verified-no-change | No UX/presentation change; category display is byte-identical (VAL-019/022). |
| `06_DESIGN_PLAN_CLI.md` | verified-no-change | No new command, flag, or output shape shipped this phase. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | Transition, claim, and maintenance flows are behavior-identical (VAL-008/009/010/020). |
| `08_QUALITY_PLAN.md` | verified-no-change | Heavy-mode gate applied as written; DEC-072 explicitly declines to weaken it. |
| `09_AI_WORKFLOW.md` | verified-no-change | No workflow-process change; the recurring ESC/TRG ID-uniqueness defect stays a CF-023 retro item, not a spec edit. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 18 remains v0.2a slice 1 of 3; slicing unchanged. |
| `GLOSSARY.md` | verified-no-change | No new canonical term; "workflow registry" and "projected workflow tables" already carry their Phase 18 meaning. |

## Triage Dispositions

- **TRG-001 (2026-07-19) — RESOLVED by DEC-067.** Phase Runner triage of the
  blocked QA plan; TL-Lite `VERDICT: BLOCK (escalation required)`. Escalated to
  full Tech Lead because it required rewriting the validation contract. Resolved
  by DEC-067 (contract + task specs refined; no product/architecture escalation
  needed). TL-Lite did not mint an `ESC-NNN` or a `triage-log.md`; because the
  finding was correctly routed and is now resolved, no `ESC-NNN` is
  retroactively assigned. The Phase Runner should append a RESUME note to the
  run ledger referencing DEC-067 when it restarts. The recurring gaps
  (ESC-NNN global uniqueness; missing triage log on BLOCK) are already tracked
  under CF-023 for the deep retrospective.
- **TRG-002 (2026-07-19) — RESOLVED by DEC-068.** Review 1 `BLOCK` on the
  VAL-020 rebuild-evidence gap; TL-Lite `VERDICT: BLOCK` (no `ESC-NNN`/triage
  log minted). Confirmed a genuine, localized evidence gap: a do-now
  remediation within Task 3, VAL-020 unchanged, no `ESC-NNN`. Resume: Developer
  Task 3 redo (add the released-claim rebuild characterization test, correct the
  `task-3-done.md` VAL-020 evidence, rerun the full Task 3 gate) → QA Reviewer
  re-review of the VAL-020 remediation → Task 4. The over-escalation and
  missing-triage-log symptoms are the same CF-023 recurring gaps.
- **TRG-003 (2026-07-19) — RESOLVED by DEC-069.** Review 2 `BLOCK` on the
  VAL-006 valid-pair UPDATE evidence gap; TL-Lite reported "escalated as
  `ESC-NNN`" but minted no concrete identifier or triage-log entry. Confirmed a
  genuine, localized evidence gap: a do-now remediation within Task 4, VAL-006
  and QA-09 unchanged, no `ESC-NNN`. Resume: Developer Task 4 redo (add the
  valid-pair UPDATE probe with persisted-pair assertion, correct the
  `task-4-done.md` VAL-006 evidence, rerun the full Task 4 gate) → QA Reviewer
  re-review of the VAL-006 remediation → Task 5. The unmintable-ESC and
  missing-triage-log symptoms are again the CF-023 recurring gaps.
- **TRG-004 (2026-07-19) — RESOLVED by DEC-070.** Phase Runner `BLOCK` on the
  Task 5 blocker report: `go test -bench . -benchmem ./...` was interrupted
  before the 10k/100k rebuild rows and `go build -o /tmp/fid ./cmd/fid` was not
  run, leaving the sound Task 5 implementation uncommitted. Confirmed a purely
  mechanical incomplete-gate condition (VAL-024 satisfied, benchmarks self-
  assert their TRD-PERF controls, all other gates PASS), not a behavior/spec/
  performance failure, and no `ESC-NNN`. Disposition: **retain** the working
  tree; do-now completion within Task 5 (finish the phase gate + build in order,
  capture benchmark `ns/op`, re-run the VAL-024 inventory, write the completion
  report and two commits), then proceed to the phase-end final QA review.
  Guardrail: a benchmark at/over target, a real gate failure, or a residual
  production map-read escalates rather than passing. The unmintable-ESC /
  missing-triage-log symptoms are again the CF-023 recurring gaps.
- **TRG-005 (2026-07-20) — RESOLVED by DEC-071.** The DEC-070 completion gate ran
  to completion and exposed a **genuine** TRD-PERF-003 failure: the 100k-event
  rebuild measured `33.3846063s` (`>30s`), correctly caught by the DEC-070
  guardrail (all other gates PASS; build/runtime correctly not run after the
  benchmark exit 1). Phase Runner `BLOCK`, reported "escalated as ESC-NEW" with no
  minted identifier/triage-log entry. Diagnosed as a linear ~0.34 ms/event
  rebuild bottlenecked by ~100k per-commit fsyncs (default `synchronous=FULL` +
  rollback journal on the rebuild connection), not a super-linear pathology.
  Disposition: **target unchanged**; do-now remediation within Task 5 — relax
  durability on the **rebuild-only** connection (`synchronous=OFF` +
  `journal_mode=MEMORY` via DSN `_pragma`, `foreign_keys` kept ON, normal
  command-path durability untouched), which is byte-identical to the projection
  and therefore preserves VAL-011/013/021. No `ESC-NNN` (bounded, single-
  subsystem, behavior-preserving optimization within Task 5's perf-control scope).
  Guardrail: if the 100k rebuild is still `≥30s` after the fix, or any rebuild-
  parity assertion breaks, escalate rather than weaken. The unmintable-ESC /
  missing-triage-log symptoms are again the CF-023 recurring gaps.
- **TRG-006 (2026-07-23) — RESOLVED by DEC-072.** QA Lead final review `BLOCK`
  with two Blocking findings: (1) `workflow.Load` deadlocks under a
  one-connection pool, so four of the seven required gate commands do not
  terminate on HEAD (`46d8013`) despite the Task 5 report recording them PASS;
  (2) the QA-30 Tech Lead spec sync is unrecorded and `02_TRD.md` §§4.3/4.6 /
  `03_TRD_SCHEMA.md` §14.11 remain stale. **Filed in `triage-log.md` as
  `TRG-001`, which collides with this phase's existing TRG-001 (blocked QA plan,
  DEC-067); relabelled TRG-006 here and in the triage log** — another instance of
  the CF-023 ID-uniqueness defect. Unlike TRG-001..005 this triage **did** mint a
  concrete escalation, **ESC-001** (itself colliding with Phase 17's ESC-001;
  left as-is per DEC-072 (10)). Confirmed genuine: a real latent connection-
  lifetime defect in `workflow.Load` (nested queries issued while the outer
  `workflows` rows still hold the connection), invisible at the real entrypoint
  only because production `db.Open` sets no pool cap. Disposition: **fix the
  loader, not the fixture**; do-now remediation within Task 5, no assertion or
  gate weakening, correct the Task 5 report's over-claimed command records, and
  apply the spec sync under Tech Lead authority. Guardrail: a changed registry
  ordering/content, a broken parity assertion, a 100k rebuild back at `≥30s`, or
  any still-non-terminating command escalates rather than passes.
