# Phase 20 TL Escalations

> Phase-scoped escalation archive. IDs are `ESC-20-NNN`, allocated from
> `triage-log.md`. Open items also appear in
> `../../current/tl_escalations_current.md` until resolved.

## Escalations

### Resolved

- **ESC-20-006 - Task 4 recovery report failed report-lint on pending commit
  metadata and three truncated go-test tokens (raised 2026-08-11 via
  `TRG-20-006`; resolved 2026-08-11 by DEC-136).** **Two causes, one of them
  mine.** (1) The `Pending commit` fields are a **Tech Lead defect**: DEC-135
  decision 5 ordered "exactly one commit containing the code, tests, and
  `task-4-done.md`", which contradicts the canonical Developer prompt's Report
  And Commit Sequence steps 3-6 and the protocol documented inside
  `report-lint.sh` — a report cannot carry the hash of the commit that
  introduces its own final content, so a `docs(task): ...` commit must follow
  the task commit. The task spec binds the Developer, so it complied and
  stopped at step 2 with the placeholders intact; the Developer behaved
  correctly. DEC-135 decision 5 is withdrawn, `task-4.md` is corrected, and
  DEC-136 decision 2 adds the standing rule that a task spec never restates or
  overrides process owned by `AGENTS.md` and the role prompts. (2) The tokens
  `0.009s`, `coverage: 80.5% of statements`, and `4.625s` are **sustained**
  failures: `AGENTS.md` requires the command's own last printed line, and these
  truncations drop the `ok`/`FAIL`/`PASS` marker that makes a token checkable,
  so the DEC-131 control fired correctly. Relaxing the heuristic or granting a
  one-off exception was refused. Because a later session may not expand a
  remembered fragment, the repair session re-runs the exact sequence itself
  after verifying `HEAD` is `e56b3a0` with a clean tree, records its own lines,
  finalizes the report against that hash and subject, and lands it in a
  `docs(task): ...` commit. Task 4's commit stands unamended; no code, contract,
  VAL, dependency, gate command, or tooling threshold changed. **CF-030** holds
  the optional report-lint message improvement.

- **ESC-20-005 - Task 4 required lint gate failed on an orphaned lock-holder
  formatter (raised 2026-08-11 via `TRG-20-005`; resolved 2026-08-11 by
  DEC-135).** **No product defect and no contract gap.** Task 4's typed
  `LockUnavailableError` carries only `LockPath`, because `contracts.md`
  §2.4/§2.5 close `lock-unavailable` to `conflict` / `retryable: true` with the
  exact detail object `{lock_path: string}` and forbid environment values in
  detail objects. That contract-correct narrowing left the pre-existing
  `describeLockHolder` formatter without a caller, and `golangci-lint run ./...`
  correctly reported `internal/journal/lock.go:185:6: func describeLockHolder is
  unused (unused)`. Because `AGENTS.md` forbids repairing and rerunning a failed
  required gate in the same session, the Developer stopped with a blocker
  report, no commit, and four unrun downstream gates; the escalation was
  therefore correct rather than avoidable. DEC-135 deletes the dead formatter,
  refuses `//nolint` suppression and any envelope change that would smuggle
  holder PID/hostname/command/actor into a closed detail object, keeps the
  `owner.json` metadata write and its `readLockMetadata` test caller, and re-runs
  Task 4 unchanged in one fresh session that inherits the uncommitted tree,
  completes the exact gate sequence with captured tokens, commits once, and
  rewrites `task-4-done.md` with a deviation note quoting the original failing
  token. No task was re-sliced, no assertion, gate command, VAL, dependency, or
  product behavior was weakened, and all seven `BLOCKED` VAL rows are
  re-recorded from the new session's own evidence.

- **ESC-20-004 - Infrastructure Review 1 gate-accounting correction (raised
  2026-08-10 via `TRG-20-004`; resolved 2026-08-10 by DEC-134).** Task 3f's
  classifier inspected only the direct `cmd/fid` parent and missed
  `internal/issues/transition_rules_test.go`, leaving that tagged in-process
  test outside both accepted gates. Correction `75c856a` first reproduced the
  finding with a repository-wide guard, then removed only the misplaced build
  constraint; the test body and assertions are unchanged. Complete post-fix
  gates passed. The same reviewer independently verified ordinary ownership,
  repository-wide classification, canonical race/coverage, and dedicated
  integration, then returned NON-BLOCK in commit `cd9e6bc`. No product code,
  public contract, VAL, dependency, or gate command changed.

- **ESC-20-003 — Task 3d gate-evidence and report-lint recovery (raised
  2026-08-09 via `TRG-20-003`; resolved 2026-08-10 by DEC-132).** The full-suite
  Go command did not terminate without output. The Developer's execution
  wrapper returned only the nested result's empty `output` field and discarded
  its live `session_id`, so outer execution-cell completion was mistaken for
  nested-command termination. The same recovery also fixed `report-lint.sh`
  matching wrapped explanatory `go test` prose. The Human authorized a bounded
  recovery super-agent to preserve Task 3d, repair canonical/generated Codex
  prompts, add workflow-lint guardrails, capture every exact gate, and close
  this escalation. Commit `d0bffbc` contains the recovery; the finalized Task
  3d report verifies that commit and passes report-lint. The exact integration
  command demonstrated the repaired protocol: it remained silent through the
  initial yield and three retained-session polls, then exited 0 with
  `cmd/fid` at 160.765s. No product code, contract, VAL, or test-performance
  refactor changed. The broader performance work remains a separate
  Human-requested follow-up pack before Task 4.

- **ESC-20-002 — Review 1b BLOCK: missing `max_excerpt_bytes` rejection rows and
  a failing required integration gate (raised 2026-08-09 via `TRG-20-002`;
  resolved 2026-08-09 by DEC-131).** **The product was correct; the evidence
  record was not.** `contracts.md` §3.2.1 requires an omitted `--prefix` to be
  derived from the directory basename and accepted only if it matches the
  grammar, so the three `cmd/fid` journeys initializing in `t.TempDir()` roots
  named `001` **must** fail. They are stale fixtures that never converged when
  init became strict at Task 2. The failure stayed invisible because the
  unit-suite helper `runInitCommand` silently appends `--prefix FID` while the
  binary-driven journeys do not, **and because every Phase 20 full-suite
  `go test` gate — tasks 1, 2, 3, 3a, 3b — was recorded with the token
  `exit 0 (no output)`, which `go test ./...` cannot produce**; one line even
  carried `wc -l` output. Review 1 claimed the same impossible token for its
  independent reproduction; review-1b is the first artifact in the phase that
  actually ran the gate. DEC-131 sustained both findings, changed no product
  code, contract, or VAL, marked the affected gate lines **UNVERIFIED** rather
  than re-running history, and ordered **task-3c** (fixture convergence, the
  missing rejection rows, provenance, one honest current-tree gate run) and
  **task-3d** (a `report-lint.sh` token-plausibility control with self-tests,
  plus correction notes on the five reports). `runInitCommand`'s prefix
  injection is accepted as a deliberate fixture default. **CF-029** extends the
  control to QA review artifacts. Resumption is task-3c → task-3d →
  `qa/review-1c.md` → Task 4; a third BLOCK returns to the Human.

- **ESC-20-001 — Review 1 BLOCK on two cross-task evidence gaps (raised
  2026-08-09 via `TRG-20-001`; resolved 2026-08-09 by DEC-130).** **Neither
  finding was a suspected product defect; both were missing evidence for an
  accepted assertion.** (1) Task 1's registry test compared only flattened
  `DetailFields` names, so no test could detect a changed detail value type,
  nullable rule, or nested object/object-list schema for any of the 23 codes,
  leaving **VAL-002** unproven. (2) Tasks 2-3 proved representative malformed
  config input and one 16-character prefix, but not the required
  just-inside/outside matrix — one-character prefix, 128-byte project ID,
  `math.MaxInt64` executor and maintenance values, overflow rejection — leaving
  **VAL-043** unproven. **DEC-130 sustained both findings and changed no
  acceptance criterion.** Recovery is two bounded test-led tasks rather than
  one, because a single session spanning both matrices repeats the
  oversized-evidence shape that caused the gap: **task-3a** (`internal/cli`
  matrix, with equal-length `duplicate-entity` parallel arrays clarified in
  contracts §2.5 and enforcement pre-authorized) and **task-3b** (config/init
  matrix, plus the do-now dangling-symlink collision regression and the do-now
  `internal/config/config.go` split). `internal/bootstrap/init.go`'s 427-line
  breach is **CF-028**, trigger-based, per QA's own note that splitting now
  would obscure the linear activation protocol. Phase Runner resumes at
  task-3a → task-3b → bounded `qa/review-1b.md` → Task 4; the six accepted
  review windows are unchanged.

## Planning-time escalation candidates considered and resolved without an ESC

| Item | Resolution |
|---|---|
| Phase 20 needs machine contracts adjacent to the Phase 21 backend port | Resolved by DEC-117: `fid capabilities`, Envelopes, Claims/Runs, and idempotency stay in Phase 21. No roadmap deviation. |
| `Envelope` is ambiguous between the JSON wrapper and the Work Envelope entity | Resolved by DEC-109 with a `GLOSSARY.md` amendment before implementation, avoiding the `AGENTS.md` ambiguous-terminology hard stop. |
| `03_TRD_SCHEMA.md` §7 blocks implementation until the slice DDL is accepted | Resolved by DEC-108/129: the Human accepted the §8 slice DDL with the implementation pack. |
| PRD-ISSUE-004 dependency-cycle rejection is unimplemented and unowned | Resolved by Human direction and DEC-114: Phase 20 task 11 owns rejection over existing links; Phase 24 still owns normalized dependency storage/commands. |
| Planning overlapped a mistaken V2 rewrite of binding `09_AI_WORKFLOW.md` | Resolved by commit `f04007d` and DEC-118. The V1 workflow is binding; the full Phase 20 pack was re-audited and re-sliced before Human approval. |
| Independent audit found implicit request/Event/init contracts and oversized tasks | Resolved by DEC-122..128: complete aggregate revision, legacy journal boundary, typed links, atomic init, closed/bounded requests, future-safe workflow authoring, and eighteen-task/six-review reslice. |
