# Phase 20 TL Escalations

> Phase-scoped escalation archive. IDs are `ESC-20-NNN`, allocated from
> `triage-log.md`. Open items also appear in
> `../../current/tl_escalations_current.md` until resolved.

## Escalations

### Resolved

- **ESC-20-021 - DEC-149 item 6 failed on a current-document line-budget
  overage outside its two-file maintenance allowlist (raised 2026-08-21 via
  `TRG-20-023`; resolved 2026-08-21 by DEC-150 and Human authority).** The
  maintenance repair itself is accepted at
  `74ea3ff0dafb4a778c1fea0c17af320a26e8e46c`: its red/green drills, retained
  guards, historical ledger checks, clean-worktree checks, and exact two-file
  containment all pass. DEC-149's own documentation commit had expanded
  `progress_current.md` to 136 lines against the unchanged 133-line tripwire,
  so item 6 could not pass inside the isolated allowlist. DEC-150 preserves the
  failed attempt, keeps the lint definition unchanged, and treats compression
  of the current status document as a separate documentation-conformance
  repair. The Human authorized one super-agent to finish the bounded final-QA
  repairs directly without another role-orchestration cycle.

- **ESC-20-018 - DEC-148's mandatory red fixture demanded a diagnostic its own
  bound event cannot emit (raised 2026-08-19 via `TRG-20-022`; resolved
  2026-08-19 by DEC-149).** **Sustained; the defect is in the directive, not the
  fixture.** DEC-148 required the pre-change red fixture to reproduce all three
  observed INVALID messages. Two of them come from the retained pre-STOP chain
  (`:613`, `:616`); the third,
  `TASK-RECOVERY does not preserve its RUN-RESUME checkpoint` (`:587`), can only
  fire on a candidate that violates the resume checkpoint. The exact bound
  successor preserves that checkpoint by construction — which is exactly why the
  superseded chain is the only thing rejecting it — so the three messages are
  the union across candidate shapes, not a property of one event. Demanding all
  three would force the red fixture onto an event that must remain INVALID after
  the fix. DEC-149 replaces `Exact verification` item 1 with 1a (bound successor:
  exactly the two chain diagnostics and **not** the checkpoint diagnostic before
  the change; `VALID` after) and 1b (retained guard: a checkpoint-violating
  candidate still emits the third message after the change). Every other DEC-148
  clause, its authorized code change, and its allowlist are unchanged.
  A **second** unreachable step was found in the same list before retry:
  `workflow-control-scope-lint.sh` hard-refuses to allowlist
  `workflow-sequence-lint.sh` (`:121-126`), and a drill asserts that refusal, so
  DEC-148's item 7 would have blocked the last step. That anti-self-repair guard
  is correct and retained; item 7 is replaced by a direct `git diff --name-only`
  containment proof, and the guarded control route with its ledger events is not
  used at all. The base is re-recorded as the DEC-149 doc commit, a docs-only
  fast-forward of `91c3d0a`. Root cause of both: a verification list transcribed
  without checking each step was reachable — the DEC-143 rule, applied to a Tech
  Lead directive. `TRG-20-022` cited `ESC-20-017` because the register held no
  open ESC for this condition; ESC-20-018 is its allocated ID and ESC-20-017 is
  not reopened. Human attention is not required.

- **ESC-20-017 - a reconciled STOP resume and retained pre-STOP recovery state
  imposed contradictory constraints on one TASK-RECOVERY (raised 2026-08-18 via
  `TRG-20-021`; resolved 2026-08-18 by DEC-148).** **Sustained.**
  `workflow-sequence-lint.sh:883-891` requires a STOP-resolving RUN-RESUME to
  carry its governing TRIAGE's task-prefixed trigger and that TRIAGE's attempt
  plus one (4). Line 1046 marks a TRIAGE as exhaustion only when the stage is
  `needs-triage`; TRG-20-020 followed a `MATERIAL-PROGRESS` REPAIR, so the
  chain survived the STOP with the unprefixed trigger at attempt 2. The bound
  TASK-RECOVERY was then judged by both 583-607 and 611-624 and could satisfy
  neither trigger form nor both attempts. The runner's RUN-RESUME was correct;
  the Tech Lead's 2026-08-17 handoff values (attempt 3, unprefixed) were wrong
  and would have failed earlier. DEC-148 records the accepted semantics — a
  reconciled resume checkpoint supersedes retained chain state, as the
  exhaustion path already does — and authorizes a bounded protected repair of
  the excluded verifier plus its drill fixtures, executed as workflow
  maintenance outside the wedged run sequence. Human attention is not required.

- **ESC-20-016 - Task 3k's repeated race gate could not fit Go's unstated
  default timeout (raised 2026-08-17 via `TRG-20-016`/`TRG-20-020`; resolved
  2026-08-18 by DEC-147).** The diagnosis is sustained; the requested remedy is
  not. `-count=10` repeats inside one test-binary invocation, so Go's default
  `-timeout 10m` bound all ten counts of a package now costing 104.318s per
  count. The escalation asked for cross-task authority over the later Task 11
  (10.17s) and Task 12 (9.21s) serial fixtures while keeping the command
  byte-identical, but removing both leaves 84.94s per count and 849.4s for ten,
  still over 600s; the unchanged command needs a 42% cut of the whole `cmd/fid`
  package. The real protected item is the gate definition, which is Tech Lead
  authority: the command becomes
  `go test -race -count=10 -timeout 30m ./internal/claims ./cmd/fid`, changing
  only the accidental wall clock. The original gate was already marginal when it
  first passed at 578.869s, 21s under the cliff. Task 3k runs one bounded
  report-scope continuation; Tasks 3i/3j/3l keep their verified commands. Human
  attention is not required.

- **ESC-20-015 - report-lint false-rejected Task 17's byte-identical literal
  command (raised 2026-08-16 via `TRG-20-015`; resolved 2026-08-16 by
  DEC-146).** The control diagnosis is sustained, but the protected `BLOCK`
  classification is not. `awk -v` interpreted the command's required `\.`
  escape before comparison even though the task declaration and report row
  were the same 104 bytes. The accepted byte-exact contract uniquely determines
  the fix, so this is a tools-only control-conformance `REPAIR`, not a gate
  definition change. Commit `88b0377` replaces both unsafe comparisons with
  opaque quoted equality and adds red/green escape fixtures; the unchanged Task
  17 report-lint invocation passes. The isolated commit passes
  `workflow-control-scope-lint.sh` from `9edb337` and changes only
  `project-docs/tools/report-lint.sh`. Commit `92cb0f4` adds the guarded V1
  recovery route, directive schema, scope verifier, and adversarial drills.
  Task 17's command, expected exit, evidence, VAL-034, product commit, and task
  boundary remain unchanged. Human attention is not required.

- **ESC-20-014 - DEC-144's replacement complete-output condition was itself
  unsatisfiable: `go mod why` prints stanzas, and every stanza has a mandatory
  header (raised 2026-08-16 via `TRG-20-014`; resolved 2026-08-16 by DEC-145).**
  **Sustained.** The exact required `go mod why -vendor -m gopkg.in/yaml.v3`
  lawfully emits one two-line stanza — `# gopkg.in/yaml.v3` then
  `(main module does not need to vendor module gopkg.in/yaml.v3)` — and exits 0.
  `go help mod why` states that every output stanza begins with a `# package` or
  `# module` header, so that header is mandatory protocol framing of the answer,
  not an import-chain line and not a VAL-032 failure. DEC-144 fixed the wording
  of the salient line and then re-asserted a whole-output claim about an output
  nobody had counted; that is my drafting failure at one level up, and the third
  on this single row. The pass condition is now stated as the tool states it:
  exit 0, exactly one stanza, first non-empty line `# gopkg.in/yaml.v3`, second
  non-empty line the parenthesized vendor-mode does-not-need result. Any
  import-path line, any further non-empty line or second stanza header, or a
  non-zero exit fails. The command, dependencies, module files, coverage
  threshold, VAL-032 assertion, and failing state are unchanged, and the passing
  state is now *tighter* because the header literal and the one-stanza bound are
  pinned. The `BLOCKED` row, the two-line note, `TRG-20-014`, and the ledger
  `STOP` remain immutable history; the continuation re-records the command from
  its own session because a result captured elsewhere is `UNVERIFIED`.
  `task-16a` and `task-16b` are classified as recovery-support slices of Task 16
  represented by the `REPAIR task-16` ledger entries, so no ledger correction is
  needed and Task 16 still owes the single `TASK task-16 COMPLETE` event. The
  third same-task continuation is authorized: distinct trigger signature,
  material progress at every stop, so `recovery-budget-exhausted` does not apply.
  Standing rule: a claim about a command's *complete* output is transcribed from
  the complete observed output, protocol framing included.

- **ESC-20-013 - DEC-143's replacement command required a literal the tool never
  prints, and the amended control parsed prose as declarations (raised
  2026-08-16 via `TRG-20-013`; resolved 2026-08-16 by DEC-144).** Both halves are
  **sustained**, and both defects are mine again. (1) `cmd/go` prints
  `(main module does not need to vendor module <path>)` under `-vendor` and
  `(main module does not need module <path>)` without it. DEC-143 added the flag
  for its test-exclusion semantics and then copied the **non-vendor** message as
  the required literal, so the command it prescribed could never print the line
  it demanded. The Developer's observed `exit_code: 0` with the vendor-mode line
  is a lawful pass; the command, the `-vendor` flag, VAL-032's assertion, the
  failing state (any import chain), and the dependency policy are all unchanged,
  and only the accepted wording of the passing state is corrected. (2) The
  amended `report-lint.sh` returned exactly six failures because
  `expected_nonzero_declaration_lines` returns every non-empty line in the
  section, and `task-16.md`'s section has two instructional prose lines and four
  wrapped continuations — six, which is the whole diagnosis. Task 16a's fixture
  used one bare declaration with no prose, a shape that appears nowhere in this
  repository: not in `task-16.md`, not in `task-17.md`, not in the task-spec
  template that generates them. The parser now treats only backtick-bullet lines
  as declarations and ignores prose; a malformed backtick bullet still fails, and
  no accepted declaration changes. `task-16b` lands it red-first, then Task 16
  reruns its full gate in one fresh session. Two standing rules follow: an
  exact-output literal is transcribed from an observed run of that exact command
  with those exact flags, never inferred from documentation or a sibling flag;
  and a control that parses an artifact section is self-tested against every real
  shape of that section in the repository at the time it lands. `afb262d` and
  `d9d6add` stand, Task 16's failed attempts stay immutable, and no product,
  contract, schema, dependency, VAL, threshold, review trigger, or task boundary
  moved.

- **ESC-20-012 - Task 16 could not close because two controls rejected truthful
  evidence (raised 2026-08-16 via `TRG-20-012`; resolved 2026-08-16 by
  DEC-143).** Both halves are **sustained**, and both defects are mine, not the
  Developer's. (1) The QA-plan command `rg -n "gopkg.in/yaml.v3" go.mod go.sum`
  was **unsatisfiable from the day it was written**: `go mod tidy` deterministically
  restores the `go.sum` checksum that the retained, accepted dependency
  `github.com/golang-migrate/migrate/v4@v4.19.1` requires, so the only ways to
  reach exit 1 were to drop an approved dependency, hand-edit `go.sum`, or leave
  the module files untidy — each of which the phase forbids. Planning validated
  only that the pattern *currently matched*, never that its finished-work state
  was reachable. The row is withdrawn and replaced by a **stronger** pair:
  `rg -ni "yaml" go.mod` rejects any yaml module requirement rather than one
  module path, and `go mod why -vendor -m gopkg.in/yaml.v3` proves build-graph
  absence with tests of dependencies excluded. `go.sum` is module-graph
  bookkeeping; a checksum line there is not a fid dependency, and the deleted
  direct requirement plus the empty build graph is what VAL-032 always meant. No
  dependency change was approved and the dependency budget is unchanged.
  (2) `report-lint.sh` failed three rows that `AGENTS.md` and `09_AI_WORKFLOW.md`
  §8 **require** a Developer to write: an absence command's real non-zero status
  recorded as `PASS`. The failure count was exactly three — the two `rg` no-match
  rows and `ls internal/format` — which is the whole diagnosis. §8's own rule
  says a control demanding more than the canonical contract and template is a
  defect in the control, and DEC-139 already banned requiring a field the
  producing role was never told to produce, so the control is amended, not the
  evidence. A task spec now pre-declares expected non-zero statuses and the
  report records the observed one beside the declared one; an undeclared
  non-zero `PASS`, a zero declaration, and a mismatch all stay failures, so
  nothing is laundered. Task 16's product work at `afb262d` stands, its failed
  attempts stay immutable, and the withdrawn command is never rerun. `task-16a`
  lands the control amendment red-first before the Task 16 continuation, and
  Task 17's spec gained the QA-plan absence command it had always been required
  to run — the same conflict, one task later, caught before it fired. No product
  behavior, contract, schema, VAL window, review trigger, or threshold changed.

- **ESC-20-011 - Task 14's exact gate conflicted with the Task 15 caller-cleanup
  boundary (raised 2026-08-15 via `TRG-20-011`; resolved 2026-08-16 by
  DEC-142).** The contradiction was real and the triage was correct: Task 14
  removed `--format` from eight legacy leaves, its lint, race/coverage, and
  integration commands then failed on the broken callers and orphaned formatter
  test helpers, and its own spec assigned exactly that repair to Task 15.
  Repairing the gate would have consumed another task's scope; changing the gate
  would have altered a protected gate definition. The defect is in the plan, not
  the code: the 14/15 split drew a task boundary through a single compile-level
  change. Resolved by moving the boundary, not by weakening a gate — Task 14 now
  owns the callers its own conversions invalidate, and all six exact commands,
  thresholds, and assertions are unchanged for both tasks. Investigation
  surfaced the fact the split had hidden: `migratedReadTestArgs` in
  `cmd/fid/issue_test.go` is a deliberate `--format json` → `--json` harness
  shim covering `ready`, `issue get`, and `issue list`, so those callers still
  compile and were never Task 14's red set. Task 15 keeps a genuine red-first
  scope — delete that shim, rewrite its callers, drop `--format` from the parity
  suites, and add the source and command-tree guards, with retained
  flag-rejection tests exempted as required VAL-016/VAL-031 negative evidence.
  Task 14 resumes as a same-task continuation on the preserved uncommitted
  edits; its three `BLOCKED` results stay as immutable failed attempts. No
  product, contract, schema, dependency, threshold, or validation change was
  accepted, and no VAL window or QA-plan review moved.

- **ESC-20-010 - Task 10 stopped on incomplete in-scope migration work (raised
  2026-08-15 via `TRG-20-010`; resolved 2026-08-15 by DEC-141).** The preserved
  Developer session completed normally after 616 seconds and was not killed or
  timed out. Its red focused suite, incomplete six-operation contract, missing
  full gate, and failing report-lint were real, but none required protected
  authority. Under the accepted typed recovery model this is
  `in-scope-implementation`, not a Human blocker. Recovery restored stash
  `779beab1`, preserved both failed focused attempts, completed the exact
  response/error/value-domain/order/limit/revision/Event matrix without
  weakening a test or contract, and discovered one additional concurrency bug:
  Task 10 constructors wrote the root capability registry while parallel tests
  read it. Registration now occurs only during package initialization. The
  exact task gates and report-lint close the task; no product exception,
  dependency, threshold change, or broader performance refactor was accepted.

- **ESC-20-009 - Review 3 found Task 9's mutation evidence short of the approved
  matrix (raised 2026-08-12 via `TRG-20-009`; resolved 2026-08-12 by DEC-140).**
  Three Blocking findings, all **sustained** and each independently confirmed
  against the tree. `--expected-revision` is wired on `issue transition`
  (`cmd/fid/issue_transition.go:28,39,50`), yet no test in the repository
  references it outside `cmd/fid/issue_update_test.go` and
  `internal/cli/codes_test.go` — the transition half of VAL-018's three-state
  matrix has no assertion at all. The only canonical `--json` transition test is
  the unsatisfied-precondition rejection, so VAL-019's legal success, unknown
  target, absent edge, and unsupported precondition are unproven and no test
  compares a transition envelope's Event identity against the journal tail. The
  owner-flag test loops both flags over `issue create` only, and the attribution
  test asserts a projected actor with no forbidden-table or Event-effect check.
  Six passing gates, clean report-lint, verified commits, and a green
  independent integration run were all true and none of them can prove a
  conditional state nobody asserted — a broad green command is not evidence of a
  missing negative case. Nothing is weakened. Because the implementation exists
  and only its proof is missing, `tasks/task-9a.md` is test-led: a failing
  assertion is a real defect fixed in place under the existing contract, and a
  contract that looks wrong is a hard stop. Task 9 closes at `86d927d` and is not
  reopened. Unlike ESC-20-008 this one resolves against the evidence record
  rather than the Tech Lead — the requirement was published in the task spec, in
  the validation contract rows the Developer prompt requires loading, and in the
  brief's DEC-130 Known Risk. The contributing Tech Lead factor, the phase's
  widest task compressed into one-line criteria, is corrected by a standing rule
  rather than by duplicating each matrix into each spec: **a task acceptance
  criterion summarizes a validation contract row, it never replaces it.**

- **ESC-20-008 - Task 8's gate rows carried verbatim tokens but no observed
  termination result (raised 2026-08-11 via `TRG-20-008`; resolved 2026-08-11 by
  DEC-139).** The routing was correct and is sustained — the rule was in force
  and unmet — but the requirement lived only in
  `project-docs/agents/phase-runner.md`, while `AGENTS.md`, `09_AI_WORKFLOW.md`
  §8, and the completion-report template asked only for the verbatim token, which
  Task 8 supplied for all six commands. **Tech Lead drafting defect**, so Task 8
  closes at `4df4d66` with no rerun. The stricter rule is kept and promoted to
  canon, **forward-effective only**: every recorded result now carries the token
  *and* the observed termination result. New standing rule — a control may not
  require an artifact field the producing role was never told to produce.
  **CF-031** opened for the mechanical report-lint check.

- **ESC-20-007 - Task 5 blocked by a stale journey inventory in the
  classification test (raised 2026-08-11 via `TRG-20-007`; resolved 2026-08-11
  by DEC-137, with DEC-138 as the systemic fix).** Task 5's implementation was
  complete and passing — format, lint, race/coverage at 77.7%, every focused
  test — when the dedicated integration gate failed because
  `TestIntegrationClassification` compares the package's tests against a
  hardcoded `wantJourneys` list that omitted the journey Task 5 added to prove
  VAL-012 rebuild parity. The list was genuinely stale, so the gate was right to
  fail, but the assertion is a **change-detector**: exact equality against a
  name list fails on legitimate growth, it produced ESC-20-004, and Tasks 14-18
  would have hit it again. No accepted decision requires "no new journey may
  ever exist", so DEC-137 relaxes it to required-presence — losing one of the
  four migrated journeys still fails — and leaves DEC-134's repository-wide
  tagged-test walk and the binary-ownership and no-sleep assertions untouched.
  The deeper problem is that this was the third consecutive escalation
  (with ESC-20-005 and ESC-20-006) in which a Developer held finished work and
  stopped for a one-line correction the tool itself had already named, at a cost
  of a blocker report, a triage, a Tech Lead decision, a runner stop, a human
  relay, and a fresh session apiece. **DEC-138** makes a stale bookkeeping
  inventory an in-session repair under strict conditions and gives TL-Lite a
  matching fourth RESUME ground, while leaving thresholds, product assertions,
  and DEC-131's evidence rules exactly as they were.

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
