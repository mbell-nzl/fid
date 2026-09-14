# Phase 22 QA Plan - Failure And Recovery

> Runner readiness refresh: **PASS**, Codex QA Lead, 2026-09-05.
> See section 11 for this session's evidence and current task-1 disposition.
> The original assessment and TL updates below remain historical, verbatim;
> their measurements are not this session's results or implementation PASS.

> Phase: 22 - Failure And Recovery
> Refreshed: 2026-09-05
> QA Lead: Claude Opus 5
> Original QA assessment: **PENDING TECH LEAD REFINEMENT** - not Ready at assessment.
> Tech Lead planning update, 2026-09-05: **READY FOR HUMAN APPROVAL** with
> [tl-refinement.md](tl-refinement.md), DEC-184..187. All original findings and
> measurements below are retained; this update records no new QA execution.
> Status: **Ready**. Human approved the refined pack under DEC-188, 2026-09-05.
> Runner still obtains the required QA Lead readiness refresh before task 1.
> Not `BLOCKED`: every required command is concrete and was executed or
> feasibility-proven fresh in this session (section 6). The plan becomes Ready
> when the Tech Lead dispositions R-01..R-14 in section 2.8; the linked addendum
> records those dispositions and the exact task/validation refinements.
> Operating mode: Heavy. No command or feasibility result carries forward.

## 1. Scope And Source Documents

Phase 22 delivers the one-Issue failure-and-recovery journey: interrupt
execution, reconcile repository state, start a linked successor Run, enforce
durable budgets and typed Stops, dispose Attention, recover activation, stale
old proof, import/export the same state, rebuild, and resume without
contradictory facts. 18 sequential tasks, 80 validation assertions, six
three-task review windows, and one independent checkpoint after task 12.

Sources read for this plan:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md` sections 3-8
- `project-docs/specs/09_AI_WORKFLOW.md` sections 2.6, 4
- `project-docs/docs/phases/phase-22/qa/validation-contract.md`
- `project-docs/docs/phases/phase-22/qa/qa-lead-handoff.md`
- `project-docs/docs/phases/phase-22/contracts.md`
- `project-docs/docs/phases/phase-22/public-interface.md`
- `project-docs/docs/phases/phase-22/schema.md`
- `project-docs/docs/phases/phase-22/execution-scope.md`
- `project-docs/docs/phases/phase-22/v1-convergence-inventory.md`
- `project-docs/docs/phases/phase-22/tasks/task-1.md` .. `task-18.md`
- `project-docs/docs/phases/phase-21/public-interface.md` (retained baseline)
- Companion artifact authored with this plan:
  `project-docs/docs/phases/phase-22/qa/consistency-review.md`

No held-out scenario, evaluator, or excluded-scenario material was read. Per
the handoff, held-out / Mode B external scenario execution is **N/A** for this
phase; the three ordinary `(human)` checks in section 9 remain required.

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: **YES**.
  The contract's `Requirement Coverage And Explicit Exclusions` section maps
  PRD-RUN, PRD-GATE, PRD-EVID, PRD-HUM, PRD-REC, PRD-PROJ, PRD-OPS, PRD-AUTH,
  PRD-TELEM, schema/CLI/flows, legacy removal, and CF-024/034/035, and names
  each deferral with its owning phase.
- Every `VAL-NNN` has one validator role: **YES**. 77 QA Lead + 3 Human
  (VAL-076..078). No row names two validators.
- Every `VAL-NNN` has a concrete evidence method: **YES**.
- Every `VAL-NNN` maps to a task/window: **YES**, verified mechanically. Every
  one of VAL-001..080 is cited by exactly one owning task spec, and no task
  cites an ID absent from the contract:

  ```text
  # contract IDs not cited by any task spec -> (none)
  # task-spec IDs absent from the contract  -> (none)
  # contract ID count                       -> 80
  ```

- Every `VAL-NNN` describes one observable: **NO**. See R-08. `VAL-056`
  bundles five independent boundary observables (declared size, actual
  expanded bytes, traversal, symlink/device/directory, duplicate entry) and
  its fail criterion names a mode the Go standard library already blocks.
  `VAL-024`, `VAL-069`, and `VAL-075` are matrix assertions whose single
  observable is "the enumerated matrix holds"; those are acceptable **only
  because** section 7 of this plan enumerates their required cases. `VAL-056`
  is not acceptable as written.
- Conditional or binary assertions require both qualifying and non-qualifying
  states: **NO for two rows**. `VAL-019` and `VAL-020` depend on
  `observation_id:?ID`, which is declared nullable while no rule permits null
  (R-04); the non-qualifying state is therefore undefined. `VAL-036`/`VAL-037`
  depend on whether recovery rejections raise Stops (R-07). All other
  conditional rows name both states.
- Parity or "matches reference" assertions have deterministic conformance
  proof: **YES**. `VAL-053`, `VAL-058`, `VAL-063`, `VAL-064`, `VAL-071`,
  `VAL-073` all require structural comparison, and the contract's Coverage
  Rules already forbid a prose parity assertion. Section 4 pins the exact
  comparison inputs for each.

### 2.8 Revisions requested from Tech Lead

Must-fix (M) blocks Ready. Improvement (I) may be applied at wrap-up.
The Tech Lead owns every edit to `validation-contract.md`,
`public-interface.md`, `contracts.md`, `schema.md`, and the task specs; this
plan proposes and never applies them.

| ID | Sev | Where | Finding and proposal |
|---|---|---|---|
| R-01 | M | `public-interface.md` 1 | Section 1 says "Retain ... `--project-root`". That flag does not exist. The shipped global flag is `--project` (Phase 21 `public-interface.md` line 18). `/tmp/fid --project-root /tmp capabilities` -> `error: capabilities accepts no arguments`, `exit_code: 2` (unknown flag consumed as a positional); `/tmp/fid --project /tmp capabilities` -> `error: project not found: /tmp`. "Retain" is unsatisfiable as written. Replace with `--project`, or record an owned rename task and a VAL row. Affects VAL-023/036/045/075. |
| R-02 | M | `public-interface.md` 3 | The closed activated-capability / Envelope-grant action set is never enumerated for Phase 22. The Phase 21 closed set is 13 names, fixed at `internal/workflow/activation.go:639`. Phase 22 adds at least `claim.expire`, `run.observe`, `run.checkpoint`, `run.recover`, `run.force-recover`, `run.observe-budget`, `run.close`, `run.abandon`, `stop.raise`, `stop.resolve`, `attention.resolve` and the activation resource/recovery actions, and removes the `run.finish` leaf in task 9 without saying whether the action name is also removed. Tasks 2, 3, 5, 7, 9, 10 each edit that allowlist. The task-18 corpus derives closed values **from the public-interface document itself** (`cmd/fid/integration/p21_conformance_test.go:362` calls `p21ClosedMembers(publicInterface)`), so an unenumerated set makes VAL-075's audit vacuous. Add an explicit closed table. |
| R-03 | M | `contracts.md` 3/9 | The dirty-digest algorithm is undefined, yet freshness now compares an **observed** digest against **caller-supplied** history. On the planning base no fid code computes a repository digest: `Repository` is validated for shape only (`internal/evidence/evidence.go:402`, `internal/workflow/envelope.go:511`, `internal/evidence/registration.go:340`), and the sole `exec.Command` user is `internal/evidence/executor.go`, which runs gate argv. Section 3 makes fid observe the digest; section 9 compares "the current observed repository target, including dirty digest" with each Receipt/Evidence target. If the adapter's algorithm differs from whatever a caller supplied, VAL-049 fails spuriously and VAL-048 passes vacuously. Define the canonical digest input set and algorithm, or state that Phase 22 retires caller-supplied digests for run/evidence targets and name the owning task. |
| R-04 | M | `public-interface.md` 2/3, `schema.md` 2 | `Recovery.observation_id:?ID`, the `run.recover` input field, and `run_recoveries.observation_id` are all nullable, but section 3 says "Force requests reference an existing live observation, never null" and "Ordinary recovery requires a terminated observation". No rule permits null. A nullable field with no legal null case cannot be evidenced in both states and invites a permissive implementation. Make it required non-null, or name the case that permits null. Affects VAL-019/VAL-020. |
| R-05 | M | `public-interface.md` 4 | `run_runtime.executor_id TEXT NOT NULL UNIQUE` makes executor identity globally unique per project, and both `run.started` and `run.recovered` insert it. Reusing an executor ID therefore violates a constraint, but the Phase 22 failure table declares no code for it. Phase 21's retained `duplicate-entity` (conflict / exit 4) plausibly applies; no document says so. Declare the code and its details. Affects VAL-012, VAL-023, VAL-075. |
| R-06 | M | `public-interface.md` 3 | `attention.resolve.fact_event_ids` is ambiguous. The leaf row says "exact qualifying resolutions for all linked Stops"; `contracts.md` 5 says resolution "cites the corresponding stop.resolved Event IDs, including an embedded Stop disposition's containing recovery Event"; `schema.md` 5 gives `attention.resolved` only the `attention_items` update and no `stop_resolutions` insert. The two readings - underlying qualifying facts versus `stop.resolved`/containing Event IDs - produce different VAL-006 wrong-subject and missing-reference fixtures. State one reading. |
| R-07 | M | `contracts.md` 5 | `executor-liveness-unknown`, `recovery-state-conflict` and `capability-unavailable` have `StopDetails` shapes but no automatic producer: section 5 names only `stop.raise`, budget boundaries, and activation recovery transitions as Stop producers, and "A rejected read or ordinary validation request creates no Stop." So a failed `run.recover` returns the **error** without raising a Stop, while the same-named Stop reaches the table only via `stop.raise`. This is easy to implement backwards and changes VAL-036/VAL-037 fixtures. State it explicitly. |
| R-08 | M | `validation-contract.md` VAL-056 | Bundles five observables, and its fail criterion "ZIP metadata bypass" names a mode the standard library already blocks: a `CreateRaw` entry with a forged `UncompressedSize64` is rejected by `archive/zip` at read time (`zip: not a valid zip file`, 0 bytes streamed; measured, section 6). Split into per-boundary sub-assertions, or enumerate its required cases and restate the metadata clause as "the cumulative expanded-byte cap trips on measured decompressed bytes", which is what can actually be evidenced. |
| R-09 | I | `schema.md` 2 | `stops.code` carries the seven-code CHECK; `attention_items.code TEXT NOT NULL` carries none. Attention copies the Stop code, so the projection can drift. Add the same CHECK. |
| R-10 | I | `public-interface.md` 4 | `activation-conflict` carries two different closed detail shapes (the new `{activation_id,expected_version,actual_version}` and the retained Phase 21 competing-prepare variant) with no discriminator, which weakens "closed recursively". Add a discriminator or a separate code. |
| R-11 | I | `public-interface.md` 1 | Section 1's blanket rule ("Added mutations take no positional arguments and require `--input`, `--actor`, `--idempotency-key`") is contradicted by the `export.import` row (positional `<archive-path>`, no Actor/input flag) and by `run.inspect`. State the exception in section 1. |
| R-12 | I | `public-interface.md` 2 | `capability-unavailable` `StopDetails` carries `owner_phase:int`, a fid-delivery concept a kit adapter cannot interpret. It is consistent with the existing `phase-N-unavailable` trust sources, so this is presentation only: note that it is opaque scheduling metadata and never control flow. |
| R-13 | M | tasks 1 and 16 `Exact Commands` | VAL-071 has ~12% measured headroom **before** any Phase 22 work: three fresh samples on the exact CF-024 reference instance today are 26385 / 26319 / 26163 ms against the 30000 ms self-assert. Phase 22 adds ten projection tables and about fourteen new/extended Event families to replay. Keep the protected threshold unchanged, and add one early-warning sample so a regression is attributed to its causing task instead of surfacing at task 17: append `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` to task 1 (schema 15 baseline) and task 16 (schema 16 baseline). This adds a command; per the handoff it must be written into those task specs before dispatch. |
| R-14 | M | task 8 `Files To Create Or Modify` | The phase gate `go test -bench . -benchmem ./...` is **flaky on the planning base**, in `TestGateRunRecordsExactNonSuccessTerminalOutcomes` (`cmd/fid/gate_run_test.go`). Two subtests are load-sensitive: `signaled` (case at lines 229-233, helper at lines 288-293 calls `syscall.Kill(os.Getpid(), syscall.SIGTERM)` then returns immediately, so under load the process reaches `os.Exit(0)` first and reports `outcome:"succeeded", exit_code:0`), and `timed out` (case at lines 225-228 uses `/bin/sleep 1` with `timeout_ms:1`, so under load the deadline expires before spawn and reports `spawn-failed`). Measured: 1 of 2 runs of the exact gate failed; `-count=15` under six-way CPU saturation failed; `-count=15` unloaded passed. Provenance: both introduced by `b694991` (2026-08-30, phase-21/task-15); `timeout_ms = 1` last touched by `76ad781` (2026-09-02, "test(gate): stabilize signaled process outcome"); both predate the phase-22 base `89d964a`. Task 8's accepted change adds the `cancelled` terminal outcome to exactly this table-driven test, so task 8 is the natural owner, but `cmd/fid/gate_run_test.go` is not in its file list. Name that file in task 8's scope with a bounded obligation to make the helper block after signalling and to make the timeout case deadline-robust. Otherwise task 17 or final QA meets a failing required command whose cause is outside its task. |

Two further items come from `consistency-review.md` and need a recorded
disposition before VAL-080 can pass: **C-06b** (the kit binds a recovery export
to the activation record; fid's `ActivationRecovery` binds none) and **C-08**
(the kit port still names `finishRun(runId, terminalState, ...)` while fid
deliberately removes caller-declared terminal state). Nine further accepted
differences must be written down; see that document's section 4.

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| 1 | Recovery schema and replay foundation | VAL-001..003 | high | First DDL change; whole chain proven applicable in section 6. Watch active-schema derivation (never hard-code 15 or 16). R-13 adds a benchmark baseline here. |
| 2 | Typed Stops and Attention domain | VAL-004..007, VAL-079 | high | CF-034 needs a **before** coverage number recorded before any edit. Measured baseline today: `internal/workflow` 27.4%. R-06/R-07 change these fixtures. |
| 3 | Lease lifecycle and executor observations | VAL-008..012 | high | Two-process race (VAL-009) needs real concurrency, not goroutines. Expiry equality boundary is exact. R-05 affects VAL-012. |
| 4 | Read-only repository reconciliation | VAL-013..016 | high | First real git observation in the codebase (`internal/project` has only `resolve.go` today). R-03 gates VAL-048/049 downstream. |
| 5 | Atomic linked Run recovery | VAL-017..021 | high | Fault injection after append, before projection. R-04 blocks the VAL-019/020 negative state. |
| 6 | Recovery CLI and early process journey | VAL-022..024 | high | Window-2 gate: a real subprocess journey and the reusable dual-adapter case runner must both exist here, not in task 18. |
| 7 | Durable budget accounting and admission | VAL-025..030, VAL-035 | high | CF-034 second baseline; `internal/evidence` measured 25.2% today. Atomic charge+Stop+Attention is the hardest single invariant. |
| 8 | Interrupted Attempts and derived Run closure | VAL-031..034 | high | Owns the R-14 flake repair. Killing the fid parent must leave an Attempt with no Receipt. |
| 9 | Stop, budget, and closure commands | VAL-036..039 | high | Removes `run.finish` (present on the base today). Committed denial must keep a nonzero exit. |
| 10 | Activation recovery domain | VAL-040..044 | high | Point-of-no-return direction matrix; C-06b disposition needed. |
| 11 | Activation CLI and config boundary | VAL-045..047 | medium | CF-035 needs a paired consuming / config-independent case. `activation-recovery-unavailable` must disappear. |
| 12 | Recovery freshness and actionable resume | VAL-048..052 | high | Last task before the independent checkpoint. R-03 is a hard prerequisite for honest VAL-048/049 evidence. |
| 13 | Recovery archive format and validation | VAL-054..057 | high | Resource ceilings are contract limits, not test tolerances. R-08 rewrites VAL-056's case list. Feasibility proven in section 6. |
| 14 | Atomic import into an empty target | VAL-058..062 | high | Staging and publication must share a filesystem; see the tmpfs constraint in section 6. |
| 15 | Import CLI and three-way recovery parity | VAL-063..065 | high | Parity must enumerate projections and normalizations up front, never post hoc. |
| 16 | Retire escalation state and duplicate truth | VAL-066..068 | high | Migration 0016 proven applicable and reversible in section 6. Provenance per row, never filename. R-13 adds a benchmark baseline here. |
| 17 | Maintenance, capability qualification, rebuild | VAL-069..072 | high | Three fresh 100k samples with no best-of. Runs the flaky benchmark gate; depends on R-14. |
| 18 | Final conformance and Human fixture | VAL-053, VAL-073..078, VAL-080 | high | Assembles, never first-implements. Corpus derives closed values from the public-interface document, so R-01/R-02 must land first. |

## 4. Required Evidence Strategy

Rows state what QA requires; they do not prescribe implementation. Every row is
deterministic unless the owner is Human.

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Fresh-init schema dump at 15 listing every declared table/index/constraint, plus rejected duplicate-identity inserts, from both the ordinary SQLite suite and the tagged built-binary suite | VAL-001 | Developer | task 1 |
| EVID-02 | Full down-to-zero then up-to-15 transcript with identical resulting schema and role-free bootstrap | VAL-002 | Developer | task 1 |
| EVID-03 | Per-Event-family live-apply and replay comparison on separate temporary databases, plus atomic rejection of unknown/invalid Effects with no partial projection | VAL-003, VAL-055 | Developer | tasks 1, 13 |
| EVID-04 | Stop history before and after resolution showing byte-identical original fact and a separate resolution referencing a qualifying Event | VAL-004 | Developer | task 2 |
| EVID-05 | Same-key and different-key Attention concurrency cases showing one open item per semantic key and all distinct Stop IDs still linked | VAL-005 | Developer | task 2 |
| EVID-06 | Valid, wrong-subject and missing-reference resolution attempts with journal proof that invalid attempts produced no Effects | VAL-006 | Developer | task 2 |
| EVID-07 | Attention counts before and after routine status/warning operations versus a real policy Stop | VAL-007, VAL-030 | Developer | tasks 2, 7 |
| EVID-08 | Before/after `go test -race -cover` output for `./internal/workflow` and `./internal/evidence`, plus named direct positive and rejection branch tests | VAL-079, VAL-035 | Developer | tasks 2, 7 |
| EVID-09 | Clock-controlled lease cases at strictly-before, exactly-equal, and after expiry, plus renewal attempts on expired/released/superseded Claims | VAL-008, VAL-011 | Developer | task 3 |
| EVID-10 | Two concurrent **real processes** acquiring the same Issue's mutating Claim: one winner, one typed conflict, and no extra Event | VAL-009 | Developer | task 3 |
| EVID-11 | Observation attempts from an activated harness with `run.observe`, from an ordinary Builder, and with a forged executor ID | VAL-010 | Developer | task 3 |
| EVID-12 | Run start binding table: matching binding accepted; mismatched Issue, revision, Claim, Obligation and reused executor identity each rejected with their declared codes | VAL-012 | Developer | task 3 |
| EVID-13 | Git fixture matrix: clean unchanged; dirty in-scope adopted / not adopted / changed digest; unmerged, changed identity, missing base, unattributable history, and each in-progress operation marker | VAL-013..015 | Developer | task 4 |
| EVID-14 | Byte-level before/after snapshots of refs, index, stash and working-tree file digests across **both** successful and failed recovery, with `.fid/` journal and projection writes named as the only permitted delta | VAL-016 | Developer | tasks 4, 5 |
| EVID-15 | Successor chain in domain tests and a real-process journey: exactly one new Run ID, exact predecessor/successor links, predecessor terminal `interrupted` | VAL-017 | Developer | task 5 |
| EVID-16 | Fault injection after journal append and during projection, plus same-key replay returning the original successor and a different digest conflicting | VAL-018, VAL-060 | Developer | tasks 5, 14 |
| EVID-17 | Closed liveness matrix over ordinary and force requests: `terminated` ordinary proceeds, `live` ordinary rejects, `unknown` always returns `executor-liveness-unknown` | VAL-019 | Developer | task 5 |
| EVID-18 | Force grant/type/reference mismatch matrix: only the bound Human force qualifies; the historical live observation is retained; subsequent predecessor fid writes are fenced | VAL-020 | Developer | task 5 |
| EVID-19 | Before/after semantic snapshot of Issue, Envelope Revision, Obligations, scope, policy and evidence identities across recovery | VAL-021 | Developer | task 5 |
| EVID-20 | Tagged subprocess journey through initialize, activate, start, observe termination, expire/release, reconcile, recover, inspect, with the built binary and no in-process substitute | VAL-022, VAL-065, VAL-074 | Developer | tasks 6, 15, 18 |
| EVID-21 | Per-leaf accepted and unknown-field request/result fixtures for every section-3 recovery leaf, matching shape, channel, status and Event | VAL-023, VAL-036, VAL-045 | Developer | tasks 6, 9, 11 |
| EVID-22 | Real-process failure corpus mapping not-found, authority, liveness, revision and conflict cases to exact codes and exit statuses | VAL-024 | Developer | task 6 |
| EVID-23 | Automatic-charge versus caller-submission cases for the five fid-owned units, with exact Attempt/Receipt/force correspondence | VAL-025 | Developer | task 7 |
| EVID-24 | Same-source-sample replay, different-key, and stale-baseline cases with totals that never decrease | VAL-026 | Developer | task 7 |
| EVID-25 | Unknown / zero / numeric baseline cases showing unknown differs from zero and blocks admission with no Attempt | VAL-027 | Developer | task 7 |
| EVID-26 | Exact-boundary admission with co-committed charge, Stop and Attention, plus append and projection fault cases, plus a denied later spawn | VAL-028 | Developer | task 7 |
| EVID-27 | Budget totals and sample chain across a linked recovery and a full rebuild, with unchanged maximum | VAL-029 | Developer | task 7 |
| EVID-28 | Cancel, timeout, signal, spawn-failure and harness-failure cases producing exactly one observed non-success Receipt, or none when observation is lost | VAL-031 | Developer | task 8 |
| EVID-29 | Kill the fid executor parent, then retry the original gate idempotency key: original Attempt visible, no invented Receipt, no second spawn | VAL-032 | Developer | task 8 |
| EVID-30 | Satisfied versus missing / stale / failed proof cases at `run.close` | VAL-033 | Developer | task 8 |
| EVID-31 | Authorized abandon versus unknown-liveness and unauthorized cases, with Claim release provenance and no new Evidence | VAL-034 | Developer | task 8 |
| EVID-32 | Cross-surface Stop predicate cases over evaluator, Issue transition, Envelope completion and execute, with all applicable reasons returned and unrelated Stops not blocking | VAL-037 | Developer | task 9 |
| EVID-33 | Committed-denial fixtures: nonzero exit, named Event and Stop/Attention IDs, no start or output frame, no spawn, and idempotent replay that does not duplicate the Stop | VAL-038 | Developer | task 9 |
| EVID-34 | Command-tree inventory proving `fid run finish` rejects as absent while close/abandon remain, plus a current-fixture inventory with no `run.finished` writer | VAL-039 | Developer | task 9 |
| EVID-35 | First-activation and upgrade prepare cases binding exactly five resources and one predecessor identity | VAL-040 | Developer | task 10 |
| EVID-36 | Before/after commit direction matrix: `restore-prior` before, `complete-forward` after, wrong direction rejected with `activation-direction-conflict` | VAL-041 | Developer | task 10 |
| EVID-37 | All-agree versus each-missing and each-mismatched resource at finalize | VAL-042 | Developer | task 10 |
| EVID-38 | Blocked normal mutations versus still-available inspection, discovery, export/validation and activation recovery during a nonterminal activation | VAL-043 | Developer | task 10 |
| EVID-39 | Replay, conflicting-key and failure-after-append retries preserving one recovery identity and an immutable stored direction | VAL-044 | Developer | task 10 |
| EVID-40 | Fault-point activation journey through the built binary at prepare, stage, commit and finalize | VAL-046 | Developer | task 11 |
| EVID-41 | Paired unknown-key cases: a consuming command rejects; help, version and archive validation remain usable with the same malformed config | VAL-047 | Developer | task 11 |
| EVID-42 | Changed-tree and dirty-digest cases showing the same Evidence identity projecting `stale`, with the digest inputs named | VAL-048 | Developer | task 12 |
| EVID-43 | Recovery and branch-rename with an identical target retaining valid proof | VAL-049 | Developer | task 12 |
| EVID-44 | Required-material versus metadata-sufficient Artifact loss cases retaining digest/provenance while only availability and policy outcome change | VAL-050 | Developer | task 12 |
| EVID-45 | Fresh process after interruption and after rebuild producing every specified resume dimension from durable facts, with no stored next-action state | VAL-051 | Developer | task 12 |
| EVID-46 | Eligible/ineligible action matrix where every advertised executable action succeeds at the matching revision and blocked actions carry reasons | VAL-052 | Developer | task 12 |
| EVID-47 | Paired Human-text and JSON render fixtures agreeing on IDs, status, reasons, budgets and uncertainty | VAL-053 | Developer | task 18 |
| EVID-48 | Archive entry and source-reference corpus proving every accepted recovery identity, history and config entry is present and forbidden content is absent | VAL-054 | Developer | task 13 |
| EVID-49 | Boundary cases for each declared ceiling (entries, per-entry bytes, per-Artifact bytes, total expanded bytes, per-Event-line bytes) plus traversal, absolute, symlink, device, directory and duplicate-entry entries, with expanded-byte accounting measured on decompressed output | VAL-056 | Developer | task 13 |
| EVID-50 | Before/after file, journal and projection snapshots across valid and invalid validation runs | VAL-057 | Developer | task 13 |
| EVID-51 | Ordered source-versus-import identity comparison over entity IDs, Event bytes, order, hashes and idempotency history | VAL-058 | Developer | task 14 |
| EVID-52 | Fault injection at staging, fsync and rename boundaries showing no partial `.fid/` before publication and complete recoverable state after | VAL-059 | Developer | task 14 |
| EVID-53 | Empty, populated and symlinked target cases with byte-level proof that existing content is preserved | VAL-061 | Developer | task 14 |
| EVID-54 | Import an active Run, then attempt mutation and recovery: both require current liveness and reconciliation, and no process is spawned | VAL-062 | Developer | task 14 |
| EVID-55 | Named JSON three-way projection comparison across live, rebuilt and imported state, with the enumerated dimension list and the enumerated environmental normalizations declared **before** the comparison | VAL-063 | Developer | task 15 |
| EVID-56 | Source, rebuilt and re-exported archive comparison over manifest and entry bytes, with declared path/publication exclusions | VAL-064 | Developer | task 15 |
| EVID-57 | Latest-schema, catalog and command tests proving no escalation table, seed, reference or command, with historical migration fixtures explicitly classified | VAL-066 | Developer | task 16 |
| EVID-58 | Writer/reader inventory and replay tests proving Run/Stop/Attention facts drive recovery and activity is inspection only | VAL-067 | Developer | task 16 |
| EVID-59 | Per-row provenance for every converged retained assertion: file/line, introducing commit and phase/task, invalidating change, changed row, and the exact executing command | VAL-068 | Developer | task 16 |
| EVID-60 | Corrupt and healthy maintenance fixtures covering broken links, abandoned Attempts, stale Claims, budgets, activation and missing artifacts, with before/after state proving no mutation | VAL-069 | Developer | task 17 |
| EVID-61 | Complete final discovery, export and fake capability matrix: three rows true, Claims false with its Phase 24 limit and owner, eight rows unchanged | VAL-070 | Developer | task 17 |
| EVID-62 | Three consecutive fresh benchmark processes, each self-asserting parity and under 30 s on the recorded reference instance, with the recorded instance identity captured in the same session | VAL-071 | Developer | task 17 |
| EVID-63 | Fresh init, full reset and rebuild at schema 16 through real processes | VAL-072 | Developer | task 17 |
| EVID-64 | One shared semantic case set executed against the test-only fake and the fid adapter with identical expected outputs and negative cases, plus proof the fake never opens fid SQLite, journal or `.fid/` files | VAL-073 | Developer | task 18 |
| EVID-65 | Operation, error and live-frame inventory derived from the accepted tables and compared against the corpus, with a non-vacuous emptiness guard | VAL-075 | Developer | task 18 |
| EVID-66 | Independent comparison against the pinned kit contract, with version currency re-verified and every accepted difference recorded | VAL-080 | QA Lead | this plan's `consistency-review.md`, rechecked at final review |
| EVID-67 | Disposable Human fixture and exact invocation instructions for the fresh resume view | VAL-076 | Developer fixture, Human judgment | task 18 |
| EVID-68 | Disposable Human fixture and instructions for paired pre-commit and post-commit activation recovery | VAL-077 | Developer fixture, Human judgment | task 18 |
| EVID-69 | Disposable Human fixture and instructions for an imported project with remaining intervention | VAL-078 | Developer fixture, Human judgment | task 18 |
| EVID-70 | Full phase gate transcript with an observed termination status recorded next to every command | no VAL | QA Lead | final review |
| EVID-71 | Security review notes covering archive expansion limits, path containment, force authority, unknown liveness, import publication, Event effect validation and old-executor fencing | no VAL | QA Lead | final review |

## 5. Deterministic Gates

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | no files listed, `exit_code: 0` | every task, final review |
| Lint | `golangci-lint run ./...` | `0 issues.`, `exit_code: 0` | every task, final review |
| Race + coverage | `go test -race -coverprofile=coverage.out ./...` | all packages `ok`, `exit_code: 0` | every task, final review |
| Coverage report | `go tool cover -func=coverage.out` | `total:` line recorded, `exit_code: 0` | every task, final review |
| Integration | `go test -tags integration -race ./cmd/fid/integration/...` | `ok  github.com/mbell-nzl/fid/cmd/fid/integration`, `exit_code: 0` | every task from task 1, final review |
| CF-034 package coverage | `go test -race -cover ./internal/workflow` (task 2) and `go test -race -cover ./internal/evidence` (task 7), once before any edit and once after | both values recorded | tasks 2 and 7 |
| Reference rebuild | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` x3 consecutive fresh processes | every sample `PASS` with parity and `< 30000 ms/rebuild` | task 17, final review; proposed additionally at tasks 1 and 16 (R-13) |
| Reference instance capture | `go version`; `go env GOOS GOARCH`; `uname -srm`; `nproc`; `grep -m1 '^model name' /proc/cpuinfo`; `grep -m1 '^MemTotal:' /proc/meminfo` | exactly the CF-024 recorded instance | task 17, final review |
| Benchmarks | `go test -bench . -benchmem ./...` | all packages `ok`, `exit_code: 0` | task 17, final review. **Flaky on the base**; see R-14 |
| Build | `go build -o /tmp/fid ./cmd/fid` | `exit 0 (no output)`, `exit_code: 0` | tasks 17, 18, final review |
| Full phase gate | the seven-command sequence in `AGENTS.md` and `08_QUALITY_PLAN.md` 7 | every command `PASS` with observed termination | final review |

No gate is narrowed, substituted, or deferred. No command in this phase has a
non-zero pass condition, so no `## Expected Non-Zero Commands` declaration is
required by any task spec. Cleanup commands are never recorded as gate results.
Report lint (`project-docs/tools/report-lint.sh <task-N-done.md>
<task-N.md>`) is a Phase Runner control after the report commit, never a
Developer or QA gate.

## 6. Command And Feasibility Validation

Heavy mode: **nothing carries forward**. Every row below was validated fresh in
this session on the phase-22 planning base (worktree clean, `HEAD 6563fb4`,
branch base `89d964a`). Working directory for every command is the repository
root `/home/mark/projects/fid` unless stated.

### 6.1 Toolchain and reference instance

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Go toolchain | executed | PASS | `go version go1.26.2 linux/amd64` (`exit_code: 0`) |
| Target platform | executed | PASS | `go env GOOS GOARCH` -> `linux` / `amd64` (`exit_code: 0`) |
| Kernel | executed | PASS | `uname -srm` -> `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64` |
| CPU count | executed | PASS | `nproc` -> `6` |
| CPU model | executed | PASS | `model name : AMD Ryzen 7 5700G with Radeon Graphics` |
| Memory | executed | PASS | `MemTotal: 8133236 kB` |
| Linter present | executed | PASS | `golangci-lint has version 2.12.2 built with go1.26.2` (`exit_code: 0`) |
| Git present | executed | PASS | `git version 2.53.0` |
| CF-024 reference identity | compared with task 17 `Documentation Updates` | PASS | all six values match the recorded reference instance exactly |

### 6.2 Phase gate commands, executed on the planning base

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | executed | PASS | `exit 0 (no output)` (terminated `exit_code: 0`) |
| `golangci-lint run ./...` | executed | PASS | `0 issues.` (terminated `exit_code: 0`) |
| `go test -race -coverprofile=coverage.out ./...` | executed | PASS | last line `ok github.com/mbell-nzl/fid/internal/workflow 29.723s coverage: 27.4% of statements` (terminated `exit_code: 0`) |
| `go tool cover -func=coverage.out` | executed | PASS | `total: (statements) 61.2%` (terminated `exit_code: 0`) |
| `go test -tags integration -race ./cmd/fid/integration/...` | executed | PASS | `ok github.com/mbell-nzl/fid/cmd/fid/integration 7.476s` (terminated `exit_code: 0`) |
| `go build -o /tmp/fid ./cmd/fid` | executed | PASS | `exit 0 (no output)` (terminated `exit_code: 0`); 14291880-byte binary produced |
| `go test -bench . -benchmem ./...` | executed twice | **FAIL then PASS** | run 1: `--- FAIL: TestGateRunRecordsExactNonSuccessTerminalOutcomes/signaled`, `FAIL github.com/mbell-nzl/fid/cmd/fid 11.016s` (terminated `exit_code: 1`, 68.5 s wall). run 2: no `--- FAIL` lines (terminated `exit_code: 0`). See R-14 |
| `go test -race -cover ./internal/workflow` | executed | PASS | `ok ... 24.428s coverage: 27.4% of statements` (terminated `exit_code: 0`) - CF-034 pre-phase baseline |
| `go test -race -cover ./internal/evidence` | executed | PASS | `ok ... 2.250s coverage: 25.2% of statements` (terminated `exit_code: 0`) - CF-034 pre-phase baseline |
| `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` | executed | PASS | `BenchmarkRebuildWorkspace100kEvents-6 1 26385134931 ns/op 26385 ms/rebuild` then `ok ... 28.905s` (terminated `exit_code: 0`) |
| Same benchmark, further samples | extracted from the two benchmark-gate runs | PASS | `26319 ms/rebuild` and `26163 ms/rebuild`; three samples, mean ~26.3 s against the 30 s self-assert = **~12% headroom before any Phase 22 work**. Drives R-13 |
| `project-docs/tools/report-lint.sh` | usage output | PASS | `Usage: project-docs/tools/report-lint.sh <completion-report.md> <task-spec.md> [--repo <git-dir>]`; exactly the two-argument Runner control in `execution-scope.md` |

The flaky-gate diagnosis was confirmed rather than assumed:

| Probe | Result |
|---|---|
| `go test -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' -count=1 ./cmd/fid` x3, unloaded | `ok` each time |
| same test `-count=15`, unloaded | `ok  github.com/mbell-nzl/fid/cmd/fid 6.692s` (`exit_code: 0`) |
| same test `-count=15` under six busy-loop CPU hogs | `--- FAIL: .../timed_out`: got `outcome:"spawn-failed"`, want `"timed-out"` (`exit_code: 1`) |
| `go test -bench . -benchmem ./...` | 1 of 2 runs failed at `.../signaled`: got `outcome:"succeeded", exit_code:0` |
| Provenance | `git log --diff-filter=A -- cmd/fid/gate_run_test.go` -> `b694991 2026-08-30`; `git log -S'FID_GATE_HELPER'` -> `b694991`; `git log -S'"timeout_ms"] = 1'` -> `76ad781 2026-09-02`. Both predate base `89d964a` |

### 6.3 Proposed DDL - applied, not read

The proposed `schema.md` migrations were executed against a real SQLite
database built from the existing chain, in scratch storage outside the
repository. No repository file was modified.

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `schema.md` 1 textual replacements in `0013` | scripted exact-match substitution | PASS | all five targets found **exactly once each** at `0013_p21_operational_kernel.up.sql` lines 9, 31, 183, 188, 205 |
| Migration chain 0001..0014 with the edited 0013 | applied via `sqlite3` | PASS | 42 tables at schema 14, `exit_code: 0` |
| Proposed `0015` up | applied | PASS | all ten new tables created; `kit_activations.version` column present with default 1; `PRAGMA foreign_key_check` returned no rows |
| Composite FK `budget_samples -> envelope_budgets(envelope_id, revision, name)` | checked against the existing DDL | PASS | `envelope_budgets` primary key is exactly `(envelope_id, revision, name)` |
| `REFERENCES claims(id)` in `0015` | checked against the chain | PASS | `0014` renames `work_claims` to `claims`, so `claims` is the correct name at 15 |
| Proposed `0016` up | applied | PASS | `escalations` dropped; `triage_events.esc_id` and `concerns.escalated_to_esc_id` dropped; `escalation-default` workflow rows deleted; guard trigger recreated |
| Proposed `0016` down | applied | PASS | `escalations` restored (structure only, as documented) |
| Proposed `0015` down | applied | PASS | all ten tables dropped; `kit_activations.version` removed |
| Full disposable reset to zero | remaining `0014..0001` down chain applied | PASS | `SELECT count(*) FROM sqlite_master WHERE type='table'` -> `0` |
| `ALTER TABLE ... DROP COLUMN` support in the **runtime** driver | version inspection | PASS | `modernc.org/sqlite v1.51.0` embeds `SQLITE_VERSION = "3.53.1"` on linux/amd64 (`lib/sqlite_linux_amd64.go:2380`), well above the 3.35 DROP COLUMN floor; the CLI used above is 3.50.6 |

This directly de-risks VAL-001, VAL-002, VAL-066 and VAL-072 before task 1.

### 6.4 Archive limits and fault injection (handoff question 1)

**Answer: YES** - the declared ceilings can be exercised on real expanded bytes
without excessive memory and without a new dependency. Measured with a
stdlib-only probe in scratch storage:

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| 8 GiB total-expanded-bytes ceiling | built a deflate fixture expanding to 8 590 983 168 bytes | PASS | fixture is **8.3 MB on disk** at default compression, **10.4 MB at `flate.BestSpeed`** |
| Fixture generation cost | timed | PASS | 19.5 s at default level; **3.5 s at `flate.BestSpeed`** - generate per test, never commit |
| Enforcing the cap on **actual** decompression | streamed through a capped reader into `io.Discard` | PASS | 8 589 934 593 bytes read in **1.851 s**, cap tripped, **12 MiB peak process `Sys` memory** |
| Forged `UncompressedSize64` ("metadata bypass") | `zip.CreateRaw` with a declared size of 42 against a 1 MiB stream | **NOT CONSTRUCTIBLE** | `archive/zip` rejects at read time: 0 bytes streamed, `zip: not a valid zip file`. Drives R-08 |
| Traversal entry `../escape.txt` | writer + reader round trip | PASS | entry name preserved verbatim by the reader |
| Absolute entry `/abs.txt` | writer + reader round trip | PASS | preserved verbatim |
| Duplicate entry name | two `Create("ok.txt")` calls | PASS | both entries present in the central directory |
| Directory entry `dir/` | round trip | PASS | reader reports `drw-rw-rw-`, `IsDir() == true` |
| Symlink entry | `SetMode(fs.ModeSymlink)` | PASS | reader reports `Lrwxrwxrwx` |
| Device entry | `SetMode(fs.ModeDevice\|fs.ModeCharDevice)` | PASS | reader reports `Dcrw-rw-rw-` |
| 1 000 000-entry ceiling | built 1 000 001 stored entries | PASS | generated in **880 ms**, **111 MB on disk**, opened and counted in **238 ms** |

**Binding constraints for tasks 13-15, to be honoured by the implementation
and checked by QA:** expanded bytes must be **counted while discarding**, never
materialised; the 111 MB entry-count fixture and the 10 MB bomb fixture must be
**generated per run**, never committed to `testdata/`.

### 6.5 Storage and process-isolation constraints

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Default temp filesystem | `df -h /tmp` with `TMPDIR` unset | **CONSTRAINT** | `/tmp` is **tmpfs, 3.9 G total, 2.8 G available**, RAM-backed on an 8 G machine. `os.MkdirTemp("")` and `t.TempDir()` land there |
| Disk-backed alternative | `df -h /` | PASS | 914 G available |
| Import staging and publication | contract requirement versus filesystem reality | PASS with constraint | `contracts.md` 10 requires a sibling staging directory **on the same filesystem** for the atomic rename; a test must keep target and staging under one temp root. A near-1 GiB Artifact fixture under `/tmp` would consume RAM - stream or relocate `TMPDIR` |
| Integration binary build | `cmd/fid/integration/main_test.go` inspected | PASS | `TestMain` builds `./cmd/fid` into fresh `os.MkdirTemp` storage and exposes it through `integrationFIDBinary(t)`; real-process journeys are supported today |
| Process kill / signal fault injection | existing precedent | PASS | `internal/evidence/live.go`, `cmd/fid/gate.go`, `cmd/fid/gate_run_test.go` already use signals and process control |
| Synchronisation discipline | derived from R-14 | **CONSTRAINT** | the existing self-SIGTERM helper is unsynchronised and loses races under load. New fault-injection tests must synchronise on an observable (marker file, pipe, blocking child) and must not rely on sleeps or on a signal winning a race |

### 6.6 Git observation feasibility (task 4)

`internal/project` currently contains only `resolve.go` and `resolve_test.go`;
no fid code runs `git` for observation today (the sole `exec.Command` user,
`internal/evidence/executor.go`, runs gate argv). Task 4 therefore introduces
the first repository observation adapter. Every primitive it needs was proven
in a scratch repository with a real merge conflict:

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| NUL-delimited status | `git status --porcelain=v1 -z` | PASS | `od -c` shows records separated by `\0`; ` M a.txt`, `UU c.txt`, `?? ...` |
| Unquoted non-ASCII paths | `git -c core.quotePath=false status --porcelain=v1 -z` | PASS | quote-bearing filename returned verbatim |
| Unmerged entries | `git ls-files -u -z` | PASS | stages 1, 2, 3 for `c.txt` |
| Operation markers | `.git/` probe | PASS | `MERGE_HEAD` present during the conflict; `REBASE_HEAD`, `CHERRY_PICK_HEAD`, `REVERT_HEAD`, `BISECT_LOG`, `rebase-merge`, `rebase-apply` are the remaining marker paths |
| HEAD and tree identity | `git rev-parse HEAD`, `git rev-parse HEAD^{tree}` | PASS | 40-hex values returned |
| Changed-path enumeration | `git diff --name-only -z HEAD` | PASS | NUL-delimited `a.txt`, `c.txt` |
| Ancestry | `git merge-base --is-ancestor` | PASS | `exit_code: 0` for the ancestor case |

**Open dependency:** none of this fixes R-03. The adapter can observe a target;
what a `dirty_digest` *is* remains undefined, and VAL-048/VAL-049 cannot be
honestly evidenced until the Tech Lead defines it.

### 6.7 CLI surface assumptions

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Global project flag | `fid --help` and probes | **FAIL - drives R-01** | help shows `--project string`; `--project-root` is unknown (`exit_code: 2`) |
| `fid run finish` present before task 9 | `fid run --help` | PASS | `finish  Terminalize one active Run.` |
| `fid claim expire` absent | `fid claim --help` | PASS | only acquire/list/release/renew |
| `fid export import` absent | `fid export --help` | PASS | only create/validate |
| Activation recovery leaves absent | `fid activation --help` | PASS | only abort/commit/component/finalize/prepare/status |
| `fid stop` / `fid attention` absent | `fid --help` | PASS | neither appears in the command tree |
| Existing mutation flags | `fid run start --help` | PASS | `--actor`, `--idempotency-key`, `--input` exactly as `public-interface.md` 1 assumes |
| `--repository` flag exists | `fid export create --help`, `fid resume --help` | PASS | `--repository string` present on both |
| Capability baseline | `fid --json capabilities` on a fresh project | PASS | twelve rows; four are `phase-22-unavailable`; recorded in `consistency-review.md` section 2 |
| Conformance corpus pattern | `cmd/fid/integration/testdata/p21_conformance/` | PASS | `manifest.json` with 142 closed values, plus 61 operations, 54 failures, 13 live frames. `p21_conformance_test.go:362` derives closed values **from the public-interface document**, which is why R-01/R-02 must land before task 18 |

### 6.8 Cross-product source availability

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Kit repository reachable | filesystem + git | PASS | `/home/mark/projects/agentic-coding-project-starter-kit` |
| Pinned commit is current | `git rev-parse HEAD`, `git status --porcelain`, `git rev-list --count <pin>..HEAD` | PASS | HEAD == `c5229617220f97ff9228df261e24718b1d1d6d89`, clean, **0** commits since |
| Named port sources present | direct read | PASS | `trd/02_INTERFACES.md` 3.1-3.7 and 4/4.1; `FID_COMPATIBILITY.md`; `decisions/0005-linked-run-interruption-recovery.md` |
| Held-out material avoided | scope discipline | PASS | only the four normal specification files above were opened |

### 6.9 Human checklist feasibility (section 9 steps)

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Human can run a built binary against a disposable project | executed | PASS | `git init` + `fid init` -> `initialized fid project in .fid/`; `fid --json capabilities` returned the full matrix |
| A resume view exists to render | `fid resume --help` | PASS | `fid resume <issue-id> [--repository]`; task 12 replaces its projection |
| An activation status view exists | `fid activation --help` | PASS | `status  Read an activation status.` |
| Fixtures are disposable and non-destructive | task 18 spec | PASS | task 18 produces disposable fixtures and invocation instructions; the Human records judgment only |
| Judgment is not delegable | role rules | PASS | VAL-076..078 are `Human`-validated; Developer supplies fixtures only, and no QA Lead check may substitute |

### 6.10 Answers to the handoff's four feasibility questions

1. **Archive limits and fault injection without excessive memory or a new
   dependency** - **YES**. 8 GiB of real expanded bytes cost 3.5 s to generate,
   1.9 s to stream, and 12 MiB of process memory, using only `archive/zip` and
   `compress/flate` (section 6.4). Constraints: count while discarding, and
   generate the large fixtures per run.
2. **Bounded trusted observer plus explicit force fact satisfying the kit's
   authority semantics on this environment** - **YES**, conditional on R-02.
   `actor_type:"human"` already exists in the Phase 21 grammar, and the
   capability allowlist is a single closed map at
   `internal/workflow/activation.go:639` that tasks 2/3/5/10 can extend. The
   semantics match the kit exactly (`consistency-review.md` C-01, C-02). What
   is missing is the *enumerated* Phase 22 closed set, which R-02 requests.
3. **Schema-15 intermediate tests and schema-16 final import each executing
   their own documented paths without moving a gate later** - **YES**. Every
   task runs both the ordinary race/coverage gate and the tagged integration
   gate from task 1, so nothing is deferred; `schema.md` 6 and task 1 both
   forbid hard-coding a schema number, so task-14/15 import fixtures must
   derive the active schema and will converge at task 16 rather than break.
   Task 16's scope already includes `internal/export/` and `cmd/fid/`. The
   residual risk is performance, not correctness, which R-13 addresses.
4. **Every task's closed request/result/Event contract complete, and every
   owned fixture executed by its listed gates** - **NO, not yet**. Gate reach
   is complete: the task-to-gate matrix covers every owned path and no fixture
   is orphaned. The *contracts* are not complete: R-01 through R-08 name
   concrete gaps (missing flag, unenumerated closed action set, undefined
   digest, a nullable field with no legal null, an undeclared failure code,
   an ambiguous resolution input, an unstated Stop-producer rule, and a
   bundled assertion with an unbuildable fail mode).

## 7. Deterministic Checklist

Every item is executed or verified by the QA Lead at final review, ends with
`(qa-lead)`, and cites the assertions it verifies. `PASS` without cited,
reproducible evidence is invalid. Code reading alone never satisfies a row.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Fresh init at schema 15 installs every declared table, index and constraint and rejects duplicate identities, in the ordinary suite and through the built binary `(qa-lead)` | VAL-001 | schema dump plus rejected-insert output from both suites |
| QA-02 | Full down-to-zero and up reproduces the same baseline with no retired authority `(qa-lead)` | VAL-002 | migration transcript and schema comparison |
| QA-03 | Each Event family reproduces identical state through live apply and replay, and unknown or invalid Effects reject atomically `(qa-lead)` | VAL-003 | per-family live/replay comparison and no-partial-projection proof |
| QA-04 | A resolved Stop's original fact is unchanged and the resolution is a separate qualifying reference `(qa-lead)` | VAL-004 | before/after Stop record and resolution Event |
| QA-05 | Equivalent Attention requests share one open identity while every distinct Stop stays linked `(qa-lead)` | VAL-005 | same-key and different-key case output |
| QA-06 | Attention resolution requires qualifying facts for every linked Stop, and invalid references produce no Effects `(qa-lead)` | VAL-006 | valid, wrong-subject and missing-reference cases plus journal proof |
| QA-07 | Routine status and warnings leave the Attention count unchanged while a real Stop raises one `(qa-lead)` | VAL-007 | before/after counts for both paths |
| QA-08 | `internal/workflow` before and after coverage recorded, with direct positive and rejection tests of the reopened recovery branches `(qa-lead)` | VAL-079 | both `go test -race -cover ./internal/workflow` outputs; baseline is 27.4% |
| QA-09 | Lease validity flips exactly at expiry, equality is expired, and renewal cannot revive `(qa-lead)` | VAL-008 | clock-controlled before/equal/after output |
| QA-10 | Two concurrent real processes yield one active mutating Claim, one typed conflict, and no extra Event `(qa-lead)` | VAL-009 | two-process transcript and Event count |
| QA-11 | Only a trusted bound observer establishes liveness; Builder and forged-ID observations reject `(qa-lead)` | VAL-010 | positive and both negative cases |
| QA-12 | Renewal preserves Run identity and authorises only the current lease version `(qa-lead)` | VAL-011 | renew-then-mutate with current and stale versions |
| QA-13 | Run start binds exact executor, optional Obligation and authority context; each mismatch rejects with its declared code `(qa-lead)` | VAL-012 | binding table output including the reused-executor case from R-05 |
| QA-14 | Clean unchanged repository state qualifies with a complete observation `(qa-lead)` | VAL-013 | reconciliation record with every declared dimension |
| QA-15 | Dirty in-scope work requires exact recorded adoption; omitted, stale and forged sets reject `(qa-lead)` | VAL-014 | adopt / no-adopt / changed-digest cases |
| QA-16 | Every unsafe or ambiguous state returns `recovery-state-conflict` without takeover `(qa-lead)` | VAL-015 | unmerged, history, path, worktree and operation-marker cases |
| QA-17 | Source bytes and Git refs, index and stash are unchanged after recovery success **and** failure `(qa-lead)` | VAL-016 | before/after digests for both outcomes |
| QA-18 | Recovery creates exactly one linked successor and leaves the predecessor terminal `(qa-lead)` | VAL-017 | domain plus real-process successor chain |
| QA-19 | Recovery is atomic and idempotent across response loss and projection interruption `(qa-lead)` | VAL-018 | fault-injection and same-key/different-key replay output |
| QA-20 | Unknown liveness blocks ordinary **and** force recovery; terminated ordinary proceeds `(qa-lead)` | VAL-019 | complete liveness matrix |
| QA-21 | Known-live takeover requires exact current Human force authority; the live fact is retained and old fid writes are fenced `(qa-lead)` | VAL-020 | grant/type/reference mismatch matrix plus fencing proof |
| QA-22 | Issue, revision, Obligations, scope, policy and evidence identities survive recovery `(qa-lead)` | VAL-021 | before/after semantic snapshot |
| QA-23 | The built binary completes the interrupted recovery journey end to end `(qa-lead)` | VAL-022 | tagged subprocess journey transcript |
| QA-24 | Every section-3 recovery leaf enforces its closed request and result schema `(qa-lead)` | VAL-023 | per-leaf accepted and unknown-field fixtures |
| QA-25 | Recovery failures map to exact codes and process exit statuses `(qa-lead)` | VAL-024 | real-process corpus covering not-found, authority, liveness, revision and conflict |
| QA-26 | fid-owned units charge only from observed execution facts; caller substitution rejects `(qa-lead)` | VAL-025 | automatic-charge and caller-submission cases |
| QA-27 | The same source sample counts once, including across recovery; invalid predecessor or baseline conflicts `(qa-lead)` | VAL-026 | replay, different-key and stale-baseline output |
| QA-28 | Unknown measurements stay null, differ from zero, and block admission with no Attempt `(qa-lead)` | VAL-027 | unknown, zero and numeric baseline cases |
| QA-29 | Exhaustion commits charge, Stop and Attention atomically, and a later spawn is denied `(qa-lead)` | VAL-028 | boundary case plus append and projection fault cases |
| QA-30 | Consumed budgets survive successor Runs and rebuild with an unchanged maximum `(qa-lead)` | VAL-029 | linked-recovery and replay snapshots |
| QA-31 | Warning is queryable and non-blocking while the maximum blocks `(qa-lead)` | VAL-030 | below / equal-warning / equal-maximum cases |
| QA-32 | Cancel, timeout, signal, spawn and harness faults yield one observed non-success Receipt or none `(qa-lead)` | VAL-031 | per-outcome fixtures with no synthetic pass |
| QA-33 | An abandoned Attempt is visible and its original key never respawns `(qa-lead)` | VAL-032 | kill-parent transcript plus retry output |
| QA-34 | Run success derives only from qualifying trusted facts `(qa-lead)` | VAL-033 | satisfied versus missing/stale/failed cases |
| QA-35 | Abandonment is distinct from interruption, releases its Claim with provenance, and creates no evidence `(qa-lead)` | VAL-034 | authorized versus unknown and unauthorized cases |
| QA-36 | `internal/evidence` before and after coverage recorded, with direct budget and abandoned-execution branch tests `(qa-lead)` | VAL-035 | both `go test -race -cover ./internal/evidence` outputs; baseline is 25.2% |
| QA-37 | Every Stop, Attention, budget and closure leaf works and rejects invalid input as specified `(qa-lead)` | VAL-036 | per-leaf JSON and Human-text cases |
| QA-38 | The same scoped Stop predicates apply across every guarded surface, with no overblocking and no bypass `(qa-lead)` | VAL-037 | evaluator, Issue, Envelope and execute cases |
| QA-39 | A committed denial names its Event and Stop/Attention IDs, keeps a nonzero exit, emits no start or output frame, and replays without duplication `(qa-lead)` | VAL-038 | live, JSON, exit and replay fixtures |
| QA-40 | `fid run finish` is absent, close and abandon remain, and no current `run.finished` writer exists `(qa-lead)` | VAL-039 | command-tree and Event-fixture inventory |
| QA-41 | Prepare binds exactly five resources and one predecessor identity `(qa-lead)` | VAL-040 | first-activation and upgrade cases |
| QA-42 | Direction respects the point of no return in both directions and rejects the wrong one `(qa-lead)` | VAL-041 | before/after commit matrix |
| QA-43 | Finalize requires complete exact agreement; each missing or mismatched resource blocks `(qa-lead)` | VAL-042 | all-agree plus per-resource negative cases |
| QA-44 | Incomplete activation blocks normal work while inspection and recovery stay usable `(qa-lead)` | VAL-043 | blocked and allowed operation matrix |
| QA-45 | Retries keep one recovery identity and an immutable stored direction `(qa-lead)` | VAL-044 | replay, conflicting-key and failure-after-append cases |
| QA-46 | Activation recovery leaves expose closed results and typed direction errors, and `activation-recovery-unavailable` is gone `(qa-lead)` | VAL-045 | CLI schema and exit corpus plus absence proof |
| QA-47 | Activation recovers at prepare, stage, commit and finalize interruption points through the built binary `(qa-lead)` | VAL-046 | fault-point process journey |
| QA-48 | A consuming command rejects an unknown config key while help and archive validation stay available `(qa-lead)` | VAL-047 | paired consuming and non-consuming cases |
| QA-49 | A changed target stales prior repository-bound proof while Evidence identity is retained `(qa-lead)` | VAL-048 | changed-tree and dirty-digest cases with the digest inputs named |
| QA-50 | An unchanged target retains proof through recovery and branch rename `(qa-lead)` | VAL-049 | identical-target recovery output |
| QA-51 | Artifact loss changes only availability and policy outcome, never digest or provenance `(qa-lead)` | VAL-050 | required-material and metadata-sufficient cases |
| QA-52 | A fresh process reconstructs every specified resume dimension from durable facts with no stored next-action state `(qa-lead)` | VAL-051 | post-interruption and post-rebuild transcripts |
| QA-53 | Every advertised executable action succeeds and every blocked action carries reasons `(qa-lead)` | VAL-052 | eligible/ineligible matrix executed, not asserted |
| QA-54 | Human and JSON recovery views agree on IDs, status, reasons, budgets and uncertainty `(qa-lead)` | VAL-053 | paired render fixtures compared structurally |
| QA-55 | Manifest v2 contains every accepted recovery identity, history and config entry, and no project prose or SQLite image `(qa-lead)` | VAL-054 | entry and source-reference corpus plus forbidden-content absence proof |
| QA-56 | Integrity-valid but semantically invalid archives reject without mutating state `(qa-lead)` | VAL-055 | unknown-Effect, duplicate-ID, broken-reference and out-of-order cases |
| QA-57 | Each declared archive ceiling and each unsafe entry type rejects before publication, with expanded bytes counted on decompressed output `(qa-lead)` | VAL-056 | per-boundary cases per the R-08 enumeration; the metadata clause evidenced as a measured cumulative cap |
| QA-58 | Valid and invalid validation runs leave project and destination state unchanged `(qa-lead)` | VAL-057 | before/after file, journal and projection snapshots |
| QA-59 | Import preserves entity IDs, Event bytes, order, hashes and idempotency history `(qa-lead)` | VAL-058 | ordered source/import comparison |
| QA-60 | No partial `.fid/` is visible before publication and complete state exists after `(qa-lead)` | VAL-059 | staging, fsync and rename fault injection |
| QA-61 | A same-digest retry replays and a changed digest conflicts, with no duplicated facts `(qa-lead)` | VAL-060 | process-level retry output |
| QA-62 | Only a safe empty target succeeds; populated and symlinked targets preserve every existing byte `(qa-lead)` | VAL-061 | per-target cases with byte comparison |
| QA-63 | Imported active work requires current liveness and reconciliation and spawns nothing `(qa-lead)` | VAL-062 | import-then-mutate and import-then-recover attempts |
| QA-64 | Live, rebuilt and imported state agree on every promised dimension under declared normalizations `(qa-lead)` | VAL-063 | three-way JSON comparison with the dimension and normalization lists declared in advance |
| QA-65 | Re-export of reconstructed state matches canonical content except declared exclusions `(qa-lead)` | VAL-064 | source/rebuilt/imported archive byte comparison |
| QA-66 | The built CLI completes export, validate, import, rebuild and resume on one interrupted state `(qa-lead)` | VAL-065 | tagged multi-repository journey |
| QA-67 | The latest schema has no escalation table, seed, reference or command, and historical fixtures are classified `(qa-lead)` | VAL-066 | latest-schema, catalog and command output plus classification list |
| QA-68 | Recovery facts come from one Event-derived source and activity remains inspection only `(qa-lead)` | VAL-067 | writer/reader inventory and replay tests |
| QA-69 | Every converged retained row records introducing commit, phase/task, invalidating change and executing command `(qa-lead)` | VAL-068 | per-row provenance table verified against `git log` output, never filenames |
| QA-70 | Maintenance classifies every corruption class and reports healthy controls clean, with no mutation `(qa-lead)` | VAL-069 | corrupt and healthy fixture matrix plus before/after state |
| QA-71 | Exactly three capability rows become true, Claims stays false with its Phase 24 limit and owner, and the other eight are unchanged `(qa-lead)` | VAL-070 | final discovery, export and fake matrices compared with the baseline recorded in section 6.7 |
| QA-72 | Three consecutive fresh benchmark processes each self-assert parity and under 30 s on the recorded reference instance `(qa-lead)` | VAL-071 | three transcripts plus the six reference-instance commands in the same session; no best-of, median, tolerance or reuse |
| QA-73 | Fresh init, reset and rebuild at schema 16 succeed with no retired authority `(qa-lead)` | VAL-072 | real-process transcripts |
| QA-74 | One shared case set passes against the fake and fid with identical expectations, and the fake never opens fid storage `(qa-lead)` | VAL-073 | dual-adapter run plus a source or runtime check of the fake's file access |
| QA-75 | The full one-Issue failure and recovery journey completes through public operations with failure injection `(qa-lead)` | VAL-074 | tagged final journey transcript with no skipped step |
| QA-76 | Every amended leaf, closed value, error code, replay case and live frame has a non-vacuous fixture `(qa-lead)` | VAL-075 | inventory derived from the accepted tables and diffed against the corpus, with the emptiness guard demonstrated |
| QA-77 | The full phase gate runs clean with an observed termination status recorded for each of the seven commands `(qa-lead)` | no VAL | complete transcript; a hung or interrupted command is `BLOCKED`, never `PASS` |
| QA-78 | `go test -bench . -benchmem ./...` is stable, not merely green once `(qa-lead)` | no VAL | at least two consecutive clean runs, plus proof that the R-14 repair landed with its own red/green evidence |
| QA-79 | Security review: archive expansion limits, path containment, force authority, unknown liveness, import publication, Event effect validation and old-executor fencing `(qa-lead)` | no VAL, supports VAL-015/016/020/056/059/061/062 | written findings citing the specific test or command output for each area |
| QA-80 | Documentation and spec comparison: `contracts.md`, `public-interface.md`, `schema.md`, `03_TRD_SCHEMA.md` 8.3 and `AGENTS.md` command list match the as-built surface `(qa-lead)` | no VAL | per-document diff notes and named deviations |
| QA-81 | Retained-test provenance audit: no row is claimed inherited without an introducing commit that predates the phase base `(qa-lead)` | no VAL, supports VAL-068 | `git log -S` / `--diff-filter=A` / `git blame` output per claim |
| QA-82 | Every mid-phase review and the task-12 checkpoint is `PASS`, and each earlier `REPAIR` has a fresh independent `PASS` re-review `(qa-lead)` | no VAL | review artifact paths and outcomes |
| QA-83 | Cross-product consistency re-check: kit version currency re-verified and every `consistency-review.md` RECORD/DECIDE item dispositioned `(qa-lead)` | VAL-080 | the three currency commands plus the disposition references in `phase-22/decisions.md` |
| QA-84 | Every required command named by each task spec appears in its completion report as `PASS` or `BLOCKED`, with an output token and observed termination `(qa-lead)` | no VAL | per-task command reconciliation table |

## 8. Mid-Phase Review Cadence

- Cadence: QA Reviewer after tasks 3, 6, 9, 12, 15 and 18, into
  `qa/review-1.md` .. `qa/review-6.md`.
- **Named exception:** an independent QA Lead checkpoint after task 12 writes
  `qa/checkpoint-after-task-12.md` and must be `PASS` before task 13 starts. It
  re-runs recovery, force and unknown liveness, budget atomicity, Stop
  disposition, activation directions and freshness, and separately validates
  import fault injection and environment feasibility **before** import
  implementation begins. Ordinary deterministic checks are not held-out
  scenarios and are not withheld from earlier windows.
- Window 2 has an extra entry condition: by the end of task 6 there must be a
  real subprocess recovery journey and a reusable dual-adapter case structure.
  An in-process substitute fails the window.
- Review template: `project-docs/templates/mid-phase-review-template.md`.
- Each review verifies its in-window `VAL-NNN` assertions and records gate
  evidence or why a gate was not run.
- Each review records exactly one outcome: `PASS`, `REPAIR` or `BLOCK`.
- `REPAIR` is required for any valid Must-fix finding with bounded,
  decision-free ownership inside a reviewed task. The review names the owning
  task, the affected assertion, the exact repair boundary, the required
  evidence, and the fresh re-review scope. The window does not advance, and the
  original reviewer never resolves its own finding.
- `BLOCK` is reserved for a protected authority boundary, work outside the
  reviewed scope, or an exhausted no-progress recovery budget. QA records a
  blocker candidate and never allocates an `ESC` ID or edits the register.
- Mid-phase QA Reviewer evidence does not change final validator ownership:
  the QA Lead remains the single final validator for all 77 deterministic rows.

## 9. Human Verification Scope

Three independently markable checks, one observable each, all `(human)`. Task
18 supplies disposable fixtures and exact invocation instructions; the Human
records judgment. No QA Lead check may substitute for these, and no Developer
or QA role may declare their outcome.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-01a | From a **fresh session with no prior chat**, open the supplied resume view for an interrupted Run and state the single next recovery action you would take `(human)` | VAL-076 | the stated action, in writing, plus a note on what in the view determined it |
| HUMAN-01b | Confirm the same view does not lead you to treat **adoption of dirty changes as proof** `(human)` | VAL-076 | explicit statement that adopted paths were understood as unproven work |
| HUMAN-02a | Open the supplied **pre-commit** activation recovery fixture and state which direction it offers and why `(human)` | VAL-077 | stated direction (`restore-prior`) and the reason |
| HUMAN-02b | Open the supplied **post-commit** fixture and state whether a rollback is possible `(human)` | VAL-077 | explicit statement that the commit is irreversible and only `complete-forward` remains |
| HUMAN-03a | Open the imported project and identify the restored work: which Issue, Run and history you are looking at `(human)` | VAL-078 | named identities matching the source project |
| HUMAN-03b | State what intervention still remains before the imported work can be mutated `(human)` | VAL-078 | statement that current liveness and reconciliation are still required and that the import granted no authority |

Fan-out rationale: each parent `VAL-NNN` names two distinct judgments (choose
an action / avoid a specific misreading; identify a direction / understand
irreversibility; recognise identity / recognise remaining work). Bundling them
would let the weaker half pass unobserved.

## 10. Phase-End Decision Rule

The phase proceeds to final review only when all of the following hold.

1. Every deterministic gate in section 5 passes with a verbatim output token
   **and** an observed termination status recorded in the same session. A hung,
   interrupted or killed command is `BLOCKED`, never `PASS`. A gate result not
   captured in the session that records it is `UNVERIFIED`, which is not
   passing.
2. Every `VAL-001..080` has a recorded final outcome of `PASS`, `FAIL`,
   `BLOCKED` or `CARRY-FORWARD`, each with cited evidence. `PASS` without
   evidence is invalid. Conditional rows without both-state evidence cannot be
   `PASS`.
3. Every `(qa-lead)` item QA-01..QA-84 has reproducible evidence, executed by
   the QA Lead. Code reading alone never satisfies a row.
4. The three `(human)` assertions VAL-076..078 are judged by the Human against
   the six fan-out rows in section 9, recorded in `qa/human-checklist.md`.
5. Every mid-phase review and the task-12 checkpoint is `PASS`; every earlier
   `REPAIR` has a fresh independent `PASS` re-review resolving its original
   Must-fix findings.
6. No unresolved Must-fix finding or blocker candidate remains, and every
   `consistency-review.md` `DECIDE` item has a recorded Tech Lead disposition.
7. Documentation updates are named and verified: `contracts.md`,
   `public-interface.md`, `schema.md`, `03_TRD_SCHEMA.md` section 8.3, the
   `AGENTS.md` command list, and the capability matrix.
8. VAL-071 evidence is exactly three fresh consecutive samples on the recorded
   reference instance, each self-asserting parity and under 30 s. No best-of,
   median, added tolerance, reused receipt or narrowed scope. If the reference
   instance is unavailable, VAL-071 is `BLOCKED` with the actual mismatch
   recorded - never approximated.

An in-scope failed gate or incomplete evidence triggers `REPAIR`; recording
`BLOCKED` does not make it acceptable. A protected missing or contradictory
contract prevents final review and is emitted as a **blocker candidate** for
TL-Lite or the full Tech Lead, which alone may allocate an `ESC` identity.

### Explicitly out of scope for this phase

Release threat model, external SAST selection, dependency vulnerability,
license and secret release gates, migration rehearsal, and full twelve-capability
qualification retain Phase 25 ownership. Held-out / Mode B external scenario
execution is `N/A`: no accepted held-out harness exists at this repository's
maturity level. This phase authorizes no live adoption and no first-adoption
claim; three `available:true` capability rows are fid's own advertisement, not
kit qualification (`consistency-review.md` section 2).

### Known documentation defect encountered while writing this plan

`AGENTS.md` line 384 directs phase-level QA to the steps in
`08_QUALITY_PLAN.md` **section 5.2**, which does not exist; that document's
section 5 is `Evidence and trust` and has no subsections. This plan therefore
uses the handoff's explicit enumeration (requirement compliance, security
review, performance spot-check, documentation and spec comparison,
retained-test provenance audit, deterministic checklist, Human checklist),
realised as QA-77..QA-84. Also, `08_QUALITY_PLAN.md` section 8 still describes
"Heavy Phase 21" as the current Heavy phase. Both are workflow-maintenance
items in the CF-033 / CF-036 family, not Phase 22 product defects, and are
raised here for the Tech Lead rather than fixed by QA.

## 11. Runner Readiness Refresh Before Task 1

QA READINESS: PASS

Author: Codex QA Lead, 2026-09-05. Reviewed branch `phase-22` at
`ec5278faee272e8106f5f8a286369b59a74a530a`, the Runner-supplied approved pack.
DEC-188 ratifies the refined contract at
`8dee0a194a494f575b7c5f64caf6ed9ac7da70f8`; the later supplied commit records
approval and the Runner handoff. Fresh `git log` verified that relationship.
The only initial dirt was untracked `run-ledger.md`, owned by the Runner.
The current progress document's pre-RUN wording is a handoff snapshot; it is
not used to contradict the Runner's subsequent RUN-START or to start another run.

This refresh used the validation-contract and qa-evidence skills. It read all
18 ordered task specs and every input named in the Runner's readiness request,
plus `kit-contract-mapping.md`, the QA/validation templates, quality plan,
workflow sections 2.6/4.3/6.1/8, schema section 8.3, config sections 2/3, and
the glossary's recovery terminology. There are no task completion reports yet;
the 18 report paths are future artifacts with explicit producing tasks.

Readiness covers the accepted plan and executable evidence strategy. It does
not declare any of VAL-001..080 satisfied. Their existing Ready states and
future final-outcome obligations remain intact. Only the build command below
was executed as a product quality gate in this refresh; the full suites,
coverage, lint, repeated outcome test, and rebuild/benchmark gates were checked
for feasibility, not represented as newly passing runs. Previous measurements
and failed attempts in sections 2/6 and `consistency-review.md` remain unchanged.

### 11.1 Coverage, Ownership, And Refinement Audit

All in-scope requirement groups remain referenced by at least one VAL: **YES**.
The contract's explicit exclusions still assign parallel Claims, Review/Human
Outcome implementation, external adapter qualification, and later removals to
their approved owners. No new protected contract gap or uncovered row was found.

An inline read-only Python audit parsed each original assertion table row,
each task's Validation Contract Coverage and Acceptance Criteria, all fenced
blocks in each Exact Commands section, and sections 7/9 of this plan. It
compared ordered command lists against the literals in section 11.2, including
repeated occurrences, and required one matching task/window and validator per
ID. The command terminated `exit_code: 0`, with these verbatim summary tokens:

```text
coverage: 80 unique rows; 77 QA Lead / 3 Human; 18 exact owners; uncovered=0; 84 QA checks; 6 Human fan-outs
commands: 124 ordered occurrences; 17 distinct; placeholders=0; full ordinary/tagged gates=18/18
contract matrices: 056-a..k=11; activated actions=23 exact members; phase reports=0 (not yet due)
```

For reproduction, enumerate table rows matching `^\| VAL-[0-9]{3} \|`, split
on unescaped pipes, require exactly VAL-001..080, and match their task/window
columns against task coverage bullets and acceptance criteria. Check the same
IDs in deterministic or Human checklist rows according to validator. Extract
**every** fenced command block under each task's Exact Commands heading,
retaining order and duplicates; compare with section 11.2's base and additions.
Required commands contain no angle-bracket placeholder, TBD, or TODO. The
illustrative report-lint argument notation in historical prose is not a
Developer command; `execution-scope.md` supplies the concrete task-1 invocation.

| Refinement | Fresh readiness disposition and remaining evidence owner |
|---|---|
| R-01/11 | `--project` matches built help; read-leaf and positional import exceptions are explicit. Tasks 6/15 and final wire checks retain their positive/negative cases. |
| R-02/12 | Independently compared the full explicit 23-action table, not just its count. Activation bootstrap, import, and read operations are excluded; task 9 removes run.finish leaf/action/Event. owner_phase is opaque metadata. Tasks 2..9/18 own the listed acceptance/rejection evidence. |
| R-03 | Contracts 3.1/3.2 and tasks 4/12 agree on canonical bytes, independent vectors, real dirty-target freshness and separate reconciliation binding; details below. |
| R-04/09 | Fresh execution of the amended DDL confirms required observation NOT NULL/FK and equal seven-code CHECKs. Task 1 still owes runtime-driver, common Event-applier, replay and built-init evidence. |
| R-05/06/07 | Exact reused-executor error shape/replay precedence, exact disposition-containing Event set with fact revalidation, and recovery rejection without Event/Stop/Attention are explicit in contracts/wire/tasks. Missing/null/nonexistent/wrong-subject observations and missing/extra/duplicate/underlying-only disposition references have rejection owners. |
| R-08 | VAL-056 is now a finite conformance matrix, 056-a..k, with individual evidence and both states. Numeric L-1/L/L+1 and actual measured decompression remain mandatory; task 13 owns validator proof. |
| R-10 | Version, competing-prepare and direction conflicts now have distinct closed shapes; tasks 10/11/18 own the corpus. |
| R-13 | Exact single-sample commands and six reference commands appear at tasks 1/16. Task 17/final QA retain three fresh processes, parity evidence and the unchanged <30 s self-assert. |
| R-14 | Repair is explicitly owned by task 8/VAL-031. Signal/timeout synchronization, cancellation, bounded load, focused 15-repeat race command and two full benchmark runs are required. Earlier encounters use DEC-186's bounded same-task amendment before retry, never a waiver. The repair has not landed. |
| C-01..13 | Every original AGREE/RECORD/DECIDE item has an accepted disposition in kit-contract-mapping.md/DEC-187. C-06b export binding and C-08 terminalState expectation mapping remain CF-037/Phase 25 qualification obligations; no production adapter qualification is claimed. |

The earlier NO answers in sections 2/6 describe the original draft. The accepted
refinements now provide one observable or explicit finite matrix per assertion,
both qualifying and rejecting states, deterministic structural parity, and
the six Human fan-outs. This disposition supplements those historical answers;
it does not rewrite the original findings or give implementation credit.

### 11.2 Command And Feasibility Validation — Fresh Session

All project commands use `/home/mark/projects/fid`. Disposable CLI/Git probes
declare their fresh temporary working directory; none uses the phase worktree
as a fixture. No additional Developer gate is introduced by this refresh.

The exact ordered base required by all 18 tasks is:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
```

| Tasks | Exact additions/order, mechanically reconciled |
|---|---|
| 1,16 | Base, then the six reference commands below, then one reference benchmark. 12 occurrences each. |
| 2 | `go test -race -cover ./internal/workflow` before edits, the same command after edits, then base. 7 occurrences. |
| 3,4,5,6,9,10,11,12,13,14,15 | Base only. 5 occurrences each. |
| 7 | `go test -race -cover ./internal/evidence` before edits, the same command after edits, then base. 7 occurrences. |
| 8 | Base, then `go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' -count=15 ./cmd/fid`, then `go test -bench . -benchmem ./...` twice. 8 occurrences. |
| 17 | Base, six reference commands, reference benchmark three times, `go test -bench . -benchmem ./...` twice, `go build -o /tmp/fid ./cmd/fid`. 17 occurrences. |
| 18 | Base, then `go build -o /tmp/fid ./cmd/fid`. 6 occurrences. |

The exact reference commands and benchmark are:

```text
go version
go env GOOS GOARCH
uname -srm
nproc
grep -m1 '^model name' /proc/cpuinfo
grep -m1 '^MemTotal:' /proc/meminfo
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal
```

Final QA still executes the seven-command full phase gate, the separate three
reference samples, and QA-78's second consecutive full benchmark run. Neither
early samples nor a successful build replaces any of these. QA-84 reconciles
all command occurrences and before/after timing. The existing report-lint
`required_commands()` reads the first fenced block; its receipt cannot replace
the QA audit of tasks 2/7's later block or distinct repeated-run receipts.

| Command or step | Fresh validation method | Result and observed evidence |
|---|---|---|
| Format and edited-file formatting prerequisites | `gofmt -h`; `goimports -h` | Feasible: `-l` and `-w` documented; help terminated 0 and 2 respectively. These are help statuses, not format-gate results. |
| `golangci-lint run ./...` | `golangci-lint run --help` | Feasible: `golangci-lint run [flags]`; terminated `exit_code: 0`. |
| Full race/coverage, package coverage and tagged integration | `go help build`, `go help test`, `go help testflag`; `go list -json ./...`; `go list -tags integration -json ./cmd/fid/integration/...` | Feasible: race supported on linux/amd64; tags/coverage flags documented. 21 ordinary packages and one tagged package; all five commands terminated `exit_code: 0`. No test result is inferred. |
| Coverage report | `go tool cover -help` plus coverage-producing command's documented usage | Feasible: `-func string`; help terminated `exit_code: 2`. coverage.out must come from that session's preceding successful coverage gate, never an old file. |
| Focused 15-repeat race and both benchmark forms | `go help testflag`; current test/benchmark definitions and `go list` selection | Feasible: `-run regexp`, `-count n`, `-bench regexp`, `-benchmem`, `-benchtime t`, `-timeout d`; help terminated `exit_code: 0`. Ordinary package includes gate_run_test.go and journal includes rebuild_bench_test.go; definitions exist at lines 213 and 19 respectively. This proves reach, not stability or speed. |
| `go build -o /tmp/fid ./cmd/fid` | Executed; retained session 63074 and polled to actual termination | PASS: `exit 0 (no output)` (terminated `exit_code: 0`). Built baseline help was then exercised. |
| Six reference-instance commands | Executed freshly | All terminated `exit_code: 0`: Go 1.26.2; linux/amd64; Linux 6.6.114.1-microsoft-standard-WSL2 x86_64; 6 CPUs; AMD Ryzen 7 5700G with Radeon Graphics; MemTotal 8133236 kB. Exact task-17 reference matches. No fresh timing qualification claimed. |
| Runner report-lint step | Read usage and required_commands implementation in project-docs/tools/report-lint.sh:40/71; checked execution-scope's exact two positional paths | Feasible once task-1-done.md is committed. Control was not run against an absent future report. No Runner control or ledger edit. |
| Test roots and fixtures | File existence plus ordinary/tagged package enumeration | 54 explicit create/modify paths exist; 56 future paths have producing tasks, including reports. Existing p21 corpus has 142 closed values, 61 operations, 54 failures, 13 live variants. Future p22 corpus/fake/journey begins at task 6; activation at 11, import at 15, final audit at 18. |
| Process and publication assumptions | Fresh Git fixture and disposable same-filesystem stage/rename/fsync probe | Probe terminated `exit_code: 0`; file, stage-directory and parent-directory fsync succeeded; target absent before rename and complete after. Linux process primitives and existing real-binary TestMain remain available. Task 14 still owes every injected fault boundary. |
| Storage constraints | Fresh stat/disk-usage inspection | /tmp device 72, 2931380224 bytes free; repository device 2096, 980212994048 bytes free. Keep stage/target siblings; stream/discard expanded data. Materialized large import fixtures belong on disk-backed storage. These are observations, not reserved capacity. |
| Kit pin currency | Exact rev-parse HEAD, status --porcelain, rev-list --count pin..HEAD against the named kit repository | Each terminated `exit_code: 0`: c5229617220f97ff9228df261e24718b1d1d6d89; `exit 0 (no output)` for clean status; `0` later commits. Final QA must repeat currency and as-built comparison. |

### 11.3 Amended DDL, Digest, And Archive Feasibility

**DDL:** The fresh probe extracted only schema.md section 1's five replacement
pairs and its four exact SQL fences. Each 0013 replacement matched exactly
once. It fed the unchanged migration chain plus amended SQL to
`sqlite3 :memory:` with `.bail on` and `PRAGMA foreign_keys=ON`. SQLite CLI
3.50.6 executed the complete chain; the runtime dependency remains
modernc.org/sqlite v1.51.0 (embedded SQLite 3.53.1, verified in the installed
linux/amd64 source). No migration file was written.

The schema source SHA-256 was
`1a41c65e587dfede3a48686b67bafc746f31573de51764f0e4250b0b4d2058ef`.
Fresh evidence: 44 tables after 0001..0014; ten added tables at 0015;
`PRAGMA table_info(run_recoveries)` marks observation_id non-null and
`PRAGMA foreign_key_list(run_recoveries)` points it to run_observations.id.
Independent disposable insert invocations accepted all seven codes in each
of stops/attention_items and rejected unknown-code in each (expected and
observed SQLite `exit_code: 1`, `CHECK constraint failed: code IN`). A null
observation insert rejected with `NOT NULL constraint failed:
run_recoveries.observation_id` (expected/observed SQLite `exit_code: 1`).
All positive SQL invocations terminated 0; foreign_key_check was empty.

The probe applied 0016, checked escalation table and both reference-column
absence, applied 0016/0015 down and 0014..0001 down to zero tables, then the
same amended full up chain. Ordered `SELECT type,name,sql FROM sqlite_master
ORDER BY type,name` results agreed exactly across reset. Final probe token:
`DDL feasibility PASS: 15 constraints; 16 removal; 16 -> 0 -> 16 exact schema equality; foreign_key_check empty; sqlite3 positive invocations exit_code=0`
(terminated `exit_code: 0`). This is SQL feasibility; it does not exercise
future Event appliers, role-free bootstrap, live/replay parity or shipped new
migrations. Task 1/16/17 evidence remains required. The historical 42-table
measurement remains attributed to the original assessment, not substituted
for this session's observed 44.

**Digests:** Rechecked contracts 3.1/3.2 against glossary, task 4/12 and
VAL-014/048/049. Dirty input binds sorted UTF-8 paths and exact HEAD/index/
worktree byte hashes/modes, including symlink target bytes, staged-restored,
mode-only, rename-as-delete/add, deletion and untracked states. Clean is
false/null; branch/time/absolute path/scope and operational metadata are
excluded. Unsupported/unmerged/unobservable/changing state rejects. Existing
Event canonical JSON sorts object keys and preserves array order; fixed
expected bytes/hashes must be independent of the repository observer.

Fresh disposable Git commands (`status --porcelain=v1 -z --untracked-files=all`,
`diff --no-renames --name-only -z HEAD`, `ls-files --others --exclude-standard -z`,
`show HEAD:a.txt`, `show :a.txt`, HEAD/tree and ancestry reads) each terminated
0. They recovered distinct old/index/worktree bytes and a quoted Unicode
filename with NUL delimiters, read a symlink target without following it, and
preserved index/refs. This proves primitive feasibility, not a delivered
observer. Task 4 must supply independent vectors and every rejection case;
task 12 must observe a real unchanged dirty target and then one changed byte
at unchanged HEAD through public/Receipt/Evidence paths. Legacy hashes remain
verbatim and can project stale. Reconciliation separately hashes the complete
object minus digest, including scope/checkpoint/safety, and is rechecked under
the write lock; it never substitutes for dirty identity or liveness.

**Archive boundaries:** Rechecked every 056-a..k row and task-13 ownership.
The numeric ceilings are 1000000 entries, 64 MiB ordinary entry, 1 GiB Artifact,
8 GiB measured cumulative bytes, 1 MiB Event line; each needs L-1/L/L+1 evidence
at the intended limiter and a valid complete control. Safe/unsafe path, regular/
symlink/device/directory, and unique/duplicate states have separate cases.
The cumulative case must not fail first on an oversized individual entry.

A fresh Python stdlib ZIP probe generated eight 1 GiB entries plus a one-byte
ninth entry into compressed memory, streamed all decompressed bytes in 1 MiB
chunks while counting/discarding, and observed the crossing at 8589934593.
It measured 37484749 compressed bytes, 11.932 s generation, 18.661 s total,
57636 KiB peak RSS; no expanded file was materialized. Another fresh probe
constructed/read 1000001 entries: 94000192-byte archive, 15.571 s,
706324 KiB peak RSS. Separate small fixtures preserved traversal, absolute,
drive-rooted, symlink, device, directory and duplicate-entry metadata. Both
probe sessions terminated `exit_code: 0` (86388 and 95145 polled to completion).

These measurements prove environment/resource feasibility only, not Go product
acceptance or rejection. Fresh `go doc` calls for archive/zip.File.Open,
archive/zip.Writer.RegisterCompressor, io.LimitReader and compress/flate.BestSpeed
each terminated 0; installed stdlib ZIP64 writer support was inspected. Task 13
uses Go stdlib, generates its own fixtures per run and proves its actual
measured-byte limiter. Forged raw ZIP rejected before streaming remains only
corruption evidence; the original failed construction is not erased.

### 11.4 Review Cadence And Human Feasibility

The six windows remain 1–3, 4–6, 7–9, 10–12, 13–15, 16–18. Reviews after
tasks **3,6,9,12,15,18** go to review-1.md through review-6.md. The independent
QA Lead checkpoint at `qa/checkpoint-after-task-12.md` is additional to
review-4 and must PASS before task 13; task-13's start condition says so.
It reruns recovery, known-live force and unknown liveness, Stop disposition,
budget atomicity, activation both directions and freshness, then refreshes
import fault-injection/environment feasibility before archive implementation.
Task 6's early subprocess journey and shared fake must already be executing.
Required repair re-reviews remain fresh and independent.

Built `/tmp/fid` help was executed for run/start/claim/resume/activation status/
export/create/validate and root help; each terminated 0. Disposable `git init`
and `fid init` plus JSON capabilities established a usable Human fixture setup.
Use a valid short directory basename or explicit `fid init --prefix QA` in
generated instructions; the fresh successful invocation printed
`initialized fid project in .fid/` (terminated `exit_code: 0`).

All six Human steps were checked for an available inspection surface and
task-18 fixture/instruction owner: HUMAN-01a/b use fresh resume and distinguish
adoption from proof; HUMAN-02a/b use paired activation status/recovery views;
HUMAN-03a/b use imported identity/history and remaining intervention views.
Task 18 must produce actual disposable checkpoints and exact invocation IDs;
QA must verify those setup steps when delivered. Future recovery/import leaves
and p22 fixtures are not present yet and receive no PASS. Only the Human judges
VAL-076..078. Held-out scenarios remain N/A; no evaluator material was read.

### 11.5 Retained Failures, Probe Corrections, And Disposition

R-14 remains an open implementation obligation. Preserve the original
`--- FAIL: TestGateRunRecordsExactNonSuccessTerminalOutcomes/signaled` and
`FAIL github.com/mbell-nzl/fid/cmd/fid 11.016s` from the full benchmark gate
(terminated `exit_code: 1`), and the saturated timed_out failure
(terminated `exit_code: 1`). The historical later green is not a repair.
Fresh blame at gate_run_test.go:225..233 confirms b694991 introduced both
cases and timeout_ms=1; 76ad781 added only the signal case's 10000 ms timeout.
Both commits were freshly proven ancestors of phase base 89d964a (each
merge-base --is-ancestor terminated 0). DEC-186's provenance correction is
accepted without editing R-14's original wording.

This readiness session also preserves its unsuccessful scratch attempts:

- Help audit attempt 1 terminated `exit_code: 1` with `AssertionError` because
  it looked for -coverprofile in `go help build`. Moving that lookup to the
  documented `go help testflag` produced `missing=[]` for every flag; the
  corrected audit terminated 0. No required product command changed.
- DDL extraction attempts 1/2 terminated `exit_code: 1`, `AssertionError`,
  before executing SQL: the first pattern missed prose outside location
  backticks; the second also selected later Event rows. Restricting extraction
  to section 1 fixed it. Attempt 3 reached the existing chain with Python's
  SQLite 3.46.1 and terminated 1, `sqlite3.OperationalError: near "CASE":
  syntax error`. The exact SQL then succeeded using installed SQLite CLI
  3.50.6 as recorded above. This was a probe/parser prerequisite correction;
  no accepted SQL, migration, dependency or gate definition changed.
- Bare `fid init` in a long random directory failed: `error: invalid project
  initialization argument`. The wrapper terminated 1; a fresh reproduction
  captured the CLI's own `exit_code: 2`. DefaultPrefix uses the directory
  basename, which violated the existing prefix grammar. The same command in
  a fresh directory named qa terminated 0, as did explicit `--prefix QA` in
  a long directory. Fixture instructions now require a valid prefix. This
  does not alter product initialization behavior or waive a gate.
- An exploratory source search named nonexistent integration/task3_journey_test.go
  and integration/init_test.go paths and terminated 2 with `No such file or
  directory`. The subsequent actual file/package inventory verified the real
  TestMain and conformance paths; neither guessed name is a required fixture.

There is no new Must-fix readiness finding or blocker candidate. The existing
R-14 implementation repair is bounded and scheduled by DEC-186; if a task-1..7
gate encounters it, the Runner obtains the exact same-task amendment and
required verification before retry. Readiness does not waive that failure or
authorize QA to repair product/tests. No new escalation ID is allocated.

The Runner may dispatch task 1 under existing DEC-188 authority. All 80 final
outcomes, six reviews, the independent task-12 checkpoint, full phase gates and
Human judgments remain due. No other agent was activated by this QA child.

Preservation verification: removing only this refresh's notice and section 11
reproduced the original QA plan bytes exactly (SHA-256
`2cc3747035881582b1207f16ba1d609af3e05bb92094deea3bec52dd283cd504`).
`consistency-review.md` remained
`9faedb906f03703be43112358cbdaa10771c44bf526d043358be2536b7f8196f`;
the contract and Runner ledger also matched their initial hashes. The
verification terminated `exit_code: 0` with `write boundary PASS: only
qa/qa-plan.md changed; runner ledger remains byte-identical`.


## Accepted archive qualification cadence — DEC-217 (2026-09-12)

The Human approved keeping all production limits and real boundary cases while
moving the large byte/count matrices out of routine task/PR execution. The
ordinary suite retains small tests through the same limiter, production limit
wiring, first-excess behavior and large-number arithmetic. Existing historical
receipts and VAL-056 acceptance cases remain unchanged.

Final QA and archive-limiter changes additionally run from /home/mark/projects/fid:

```text
go test -tags archive_scale -race ./internal/export -run '^TestArchive(ExpandedByteBoundaries|CumulativeBoundary|EntryCountBoundaries)$' -count=1
```

Run this locally and serially with other heavy commands. Missing or failing
production-scale evidence leaves VAL-056 unclosed; ordinary small-case passes
do not substitute for it. No timeout, product ceiling or error semantics change.
