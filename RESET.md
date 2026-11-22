# fid reset

Approved by Mark · 31 October 2026 · [ADR-0013](docs/adr/0013-reset-fid-to-a-ledger-with-receipts.md) ·
[ADR-0014](docs/adr/0014-build-the-golden-path-on-a-lean-core.md) ·
[decisions](docs/reset/2026-10-31-reset-decisions.md)

**Goal:** highly automated agentic engineering that keeps quality high: a
software factory, not a slop factory. The reset succeeds if fid with the kit
is a markedly better, more robust workflow than V1, judged by the kill
criteria below.

## Rules

- **Scope.** Nothing joins unless a golden-path step or a kill criterion needs
  it. Everything else goes on the later list in the decisions record.
- **Budgets.**
  - Product Go code grows by at most 10% net before the trial (46.6k lines
    today, so at most about 51k).
    - ADR-0014 lets product Go exceed this cap until the kernel's removal
      (1.5b–1.5d), which cuts about 18–20k lines. The cap applies again
      from then, before the trial. On 2026-11-14, after 1.5a2, product Go
      was 54,547 lines (non-test `.go` files under `cmd/` and `internal/`).
  - The loop stays under 150 lines.
  - The kit install script stays under 80 lines.
- **One task = one PR = one session.** Use TDD; `make test` must pass. The PR
  description is the completion report, and small decisions go there too.
- **Review.** A fresh session on the other model reviews each PR by default:
  Codex reviews Claude's work and Claude reviews Codex's. When the other
  model isn't available (no subscription, or its usage is limited), a fresh
  session of the same model reviews instead. The fallback also covers the
  separate review the decisions record's work model adds for Assured work:
  that review still runs, as another fresh session. Until 1.6's `done`
  enforces that separate review on fid's own work, the reset's own Assured
  PRs keep one review per PR, as AGENTS.md step 7 says (Mark, 11 November
  2026). One pass covers code and test adequacy. Findings go in a PR
  comment, which names the model that reviewed. Once `fid review` works, fid
  records them and posts that comment itself. fid also records which model
  built the change and which model reviewed each round (2.1). When no review
  of an Assured change came from a model other than the builder's, the
  approval view (2.3) and `fid show` say "no other-model review", so Mark's
  approval covers it knowingly.
- **Lanes, by hand until fid computes them.**
  - Routine: no production code changed.
  - Standard: any code change.
  - Assured: the trust boundary, journal format, migrations, `go.mod`,
    releases, `.fid.toml`, CI workflows and the Makefile.
- **Merging.** Mark merges every PR during the reset. He reads Assured diffs
  in full and skims the rest.
- **Refuse only on missing evidence.** Numbers are signals, not gates.
- **History stays out of context.** Agents don't read the V1 phase archives.

## Tasks

| # | Task | Lane | Done when |
|---|---|---|---|
| **0.1** | Accept ADR-0013; approve this plan | Mark | Signed off (done 31 October 2026) |
| **0.2** | Write a one-page AGENTS.md; slim CLAUDE.md; archive the V1 role agents and phase skills | Routine | A fresh session learns the working method from AGENTS.md alone (done 31 October 2026) |
| **0.3** | Add a `push: main` CI trigger. Branch protection is deferred: GitHub Free doesn't offer it on private repos (since 2026-11-14 the repository is public and `main` has a ruleset) | Assured | A merge to `main` runs CI (done 31 October 2026) |
| **0.4** | Move the CI scripts and conformance fixtures out of `project-docs/` | Assured | CI passes with `project-docs/` deleted locally (done 31 October 2026) |
| **0.5** | Archive fid's V1 docs: tag `v1-workflow-final`, copy them to the archive repo, keep the reset documents in `docs/reset/` | Routine | `project-docs/` is gone from `main`; CI passes (done 1 November 2026) |
| **0.6** | Kit repo with the reviewer agent (code and test-adequacy profiles) | kit | It reviews this reset's first PR (done 2 November 2026) |
| **0.7** | Cut CI cost: the race detector only where goroutines start; a passed run reused on `main`; no CI on drafts; superseded runs cancelled; one test shard. `make fmt-check` fails on unformatted files | Assured | The race detector covers exactly the packages whose product code starts goroutines, enforced by a check; a merge whose tree already passed CI reuses that run on `main`; draft PRs run no CI; the PR reports billed minutes and wall time before and after (done 7 November 2026) |
| **1.1a** | Route every ledger path through one function, with no change in behaviour | Standard | Only `config.LedgerDir` joins the ledger onto a project root, and the existing tests pass unchanged (done 2 November 2026) |
| **1.1b** | Tests reach ledger files only through `config.LedgerDir`, with no change in behaviour | Routine | The 1.1a audit also covers test files and passes, and no expected value changed (done 3 November 2026) |
| **1.1c1** | Test workspaces are guarded git repos, with no change in behaviour | Assured | Every workspace the `cmd/fid` and integration tests give fid is a git repo whose top level is itself, except where a test hides `.git` on purpose; `GIT_CEILING_DIRECTORIES` is set for those test processes; no expected value changed (done 4 November 2026) |
| **1.1c2** | fid needs a git repository: discovery goes through git | Standard | Outside a git repository every command that opens or creates a project, `fid init` included, refuses; the ledger is still at `<root>/.fid` (done 4 November 2026) |
| **1.1d1** | Every test workspace that reaches the ledger is a guarded git repo, `internal/` packages included, with no change in behaviour | Standard | Every root a test gives `config.LedgerDir` is a git repo whose top level is itself, except where a test puts a ledger outside a repository on purpose; the guard helpers live in one test-only package that every test process making git repos uses, with `GIT_CEILING_DIRECTORIES` set; no expected value changed (done 4 November 2026) |
| **1.1d2** | Output names ledger files by absolute path; the CLI spec updated | Standard | Text and JSON output (`journal_path`, `journal_file`, `lock_path`, `fid init`'s message) name ledger files by absolute path; the ledger is still at `<root>/.fid` (done 5 November 2026) |
| **1.1d3** | Move the ledger to the git common directory, with staging beside it; the config spec updated | Standard | Two worktrees share one ledger (V3's claim and `next` parts belong to 1.4) (done 5 November 2026) |
| **1.1d4** | Discovery refuses when git and `LedgerDir` disagree, or when git could track the ledger | Standard | fid never writes outside `--git-common-dir`/fid, nor where git could track it: where git and `LedgerDir` name different ledgers, or git's common directory is inside the worktree but not within a `.git` directory, every command that opens or creates a project refuses with `repository-unavailable` / `unsafe-path` (done 7 November 2026) |
| **1.1d5** | A `.fid/` in a worktree is an ordinary path | Assured | A leftover `.fid/` makes the worktree dirty and can be adopted like any other path; the `dirty` help and the CLI spec no longer single it out (done 6 November 2026) |
| **1.2a** | `.fid.toml` with the project ID, written by `fid init`; the backup target in user-level config, asked for by `fid init` only in an interactive terminal | Assured | A fresh `fid init` writes a stable project ID to `.fid.toml`; `fid init` on a clone whose `.fid.toml` already has an ID keeps it; the target is asked for only in a terminal and stored in user config (done 7 November 2026) |
| **1.2b** | `fid backup`: export into fid's own clone of the backup repo, commit and push; each backup's time and outcome recorded | Assured | A backup to a local bare repo lands in the project's folder; a failure is recorded, not hidden (done 8 November 2026) |
| **1.2c** | `fid restore` (no arguments) and the V8 restore test (back up, wipe, restore, compare) | Assured | V8 (done 8 November 2026) |
| **1.2d** | Backup storage for the long run: fid's clone packs itself, and the backup folder holds the export's files unzipped | Assured | After 100 backups, fid's clone in the cache is packed, not a loose copy per backup; the folder holds the export's files unzipped (journal files of at most 64 MB, `manifest.json`, artifacts), the config spec describes it, `fid backup` writes it and `fid restore` zips it back up for import; V8 still passes (done 8 November 2026) |
| **1.3a** | Display package; `fid status` with the backup-age line | Standard | whet-style output; golden tests for text and JSON; `fid status` warns after 24 hours without a backup, after a failed push, or with no target set, and makes it a NEEDS YOU item after 3 days (done 8 November 2026) |
| **1.3b** | `fid show`; a pager for long output | Standard | whet-style output; golden tests for text and JSON; on a terminal, `status` and `show` output taller than the screen pages through `less`, as git's does (done 8 November 2026) |
| **1.4a** | `fid next --claim` on the lean core: a hold record, no Envelope; atomic, respects dependencies, ordered, typed skip reasons, `--explain` | Assured | V3 (done 11 November 2026) |
| **1.4b** | READY in `fid status`, with `next`'s verdict for each task, and the `Working` and `Scheduling` rows in `fid show` | Standard | Golden tests for text and JSON show READY's verdicts and both rows (done 2026-11-13) |
| **2.2a** | Policy from `origin/<default>` with built-in defaults, never fetched by fid; the whole `.fid.toml` schema (checks, Assured paths, check-defining paths, `auto_merge` per lane, rework cap, token budget); trust by actor type; a policy snapshot recorded when the effective policy changes | Assured | V7 (done 2026-11-13) |
| **1.5a1** | Record Mark's answers to 2.2a's proposals; a missing `[review] test_adequacy` key means not set up; the `policy-unavailable` failure; `fid status` reads policy and warns when it can't be read, when a lane has no checks set and when test paths aren't set (decisions record, 14 November 2026) | Assured | `fid status` shows the policy warnings in text and JSON and never writes a policy snapshot; missing and empty test paths stay distinct through snapshots, replay and restore; both policy failures refuse with `policy-unavailable` (done 2026-11-14) |
| **1.5a2** | Hold heartbeat; `fid park --interrupted`; an interrupted task kept for its worktree; takeover by a person, with `fid takeover`; the person's-decision record (decisions record: 11 November 2026, item 4) | Assured | A started task that lapses or is parked as interrupted resumes in its worktree under a linked hold, and no other worktree takes it; a person takes it over; a stale session's heartbeat and commands are refused; replay, export and V8 keep the new records (done 2026-11-14) |
| **1.5a3** | `fid park --waiting` and `fid answer` for a task waiting for Mark, whose worktree gets no other task meanwhile; NEEDS YOU for a task waiting for Mark and for a task kept for a worktree that no longer exists (decisions record: 11 November 2026, item 5; 12 November 2026, item 12 and what an answer binds to) | Assured | A waiting task's worktree gets no other task until a person answers; other worktrees carry on; the resumed task's JSON carries the answer; a missing worktree shows under NEEDS YOU with the takeover command; replay, export and V8 keep the new records (done 2026-11-14) |
| **1.5a4** | The fixed task states `todo → in_progress → in_review → done`, plus `cancelled`, moved by fid's commands; `fid start`; a task's lane (decisions record: 11 November 2026, item 3) | Assured | New tasks start in `todo`; `fid start` needs the task's live hold and moves it to `in_progress`; the lane and state show in text and JSON; agents can't lower a lane or step around the fixed states; subtasks are unchanged (done 2026-11-14) |
| **1.5a5** | `fid check` for named checks on the lean core: its start recorded before running, the tree looked at again afterwards, the output's digest and excerpt, the toolchain, the environment and Artifacts | Assured | A check records its start before running, looks at the tree again afterwards and keeps its Artifacts; a dirty or changed tree never counts; retained files survive export and restore (done 2026-11-14) |
| **1.5a6** | The checks each lane requires, with `when_changed`; which result counts | Assured | `fid check` with no names checks what the task's lane requires; a lane with no `checks` key is refused, and `checks = []` requires none (done 2026-11-14) |
| **1.5a7** | Hand-off: every required check passed on the head commit's tree, the policy in force recorded, the task moved to `in_review`; the worktree while its task is in review | Assured | An agent works a task without touching kernel commands, from `next --claim` to hand-off (done 2026-11-14) |
| **1.5b1** | Remove the kernel, part 1a: the Artifact verifier `fid artifact` keeps, extracted so it works from lean check results without an Activation; lean test fixtures; V8 and the export fixtures and tests rebuilt on lean check results, V8 comparing `next` in place of `ready` (decisions record, settled 12 and 14 November 2026) | Assured | V8 and the export tests pass on lean check results, none weakened; neither creates a kernel record (done 2026-11-14) |
| **1.5b2** | Remove the kernel, part 1b: every other test that outlives the kernel (journal, projection, issue and read fixtures; integration fixtures and their shared helpers) rebuilt on lean records | Assured | No test that outlives the kernel creates a kernel record, except `fid artifact`'s own command tests, which move with its switch in 1.5c2; none weakened (done 2026-11-14) |
| **1.5c1** | Remove the kernel, part 2a: the fixed task states and their lookups moved out of the kernel’s workflow package, behaviour unchanged | Assured | No code that stays imports the workflow package for task states; every test passes, none weakened (done 2026-11-14) |
| **1.5c2** | Remove the kernel, part 2b: its commands; `fid ready`; Concerns, requirement and decision references; `fid capabilities`; `fid stats`; the whole of `fid workflow`; the current `review` and `finding` commands, which 2.1 replaces; with their tests and corpora; `fid artifact` switched onto lean check results (decisions record: 11 November 2026, item 6; 12 and 14 November 2026) | Assured | No kernel command remains, nor any other this row removes; `fid artifact` reads lean check results; every remaining test passes, none weakened (done 2026-11-15) |
| **1.5c3** | Remove the kernel, part 2c: its domain packages and their tests | Assured | No kernel domain package remains; every remaining test passes, none weakened (done 2026-11-15) |
| **1.5d1** | Remove the kernel, part 3a: its journal and export catalog entries, and those of Concerns and of requirement and decision references; manifest-v1 validation (decisions record, settled 12 and 14 November 2026) | Assured | A journal or archive holding any removed event is refused as `schema-unsupported`, changing nothing; every remaining test passes, V8 and the export tests included, none weakened (done 2026-11-15) |
| **1.5d2** | Remove the kernel, part 3b: its tables, and those of Concerns and of requirement and decision references, dropped by a migration; its maintenance checks | Assured | No kernel table remains, nor any table of Concerns or of requirement and decision references; every remaining test passes, V8 and the export tests included, none weakened (done 2026-11-15) |
| **1.5e** | Test suites clean up their temporary directories: a helper child of the test binary skips the suite's setup, Go's telemetry cannot recreate the suite's directories after cleanup, and child coverage stays in its parent-owned fixture | Routine | A `make test` run and an integration run leave nothing in the temp directory (done 2026-11-15) |
| **1.5f** | cmd/fid tests survive a run that crosses UTC midnight | Routine | Journal assertions span every dated file, proven by tests that force the crossing in cmd/fid and integration (done 2026-11-15) |
| **1.6a** | Merge-aware `done`, part 1: GitHub facts through `gh`: a small read adapter with an injected runner, and typed facts for the PR, its head's tree and the CI runs on its merge commit (decisions record, 2026-11-15) | Assured | Canned `gh` responses prove the adapter returns the exact facts or a typed failure; a mismatched repository, branch or commit and an incomplete response never count as success; no `gh` call fetches, merges or writes (done 2026-11-15) |
| **1.6b** | Merge-aware `done`, part 2: a confirmed merge recorded against its task, and each completed CI run on the merge commit recorded as a check result from CI (run, attempt, workflow, commit, conclusion) | Assured | Merge and CI records survive replay, export and restore, V8 included; a wrong task, PR or commit is refused and changes nothing; local check results are unchanged (done 2026-11-15) |
| **1.6c** | Merge-aware `done`, part 3: the change fingerprint, `git patch-id --verbatim` of the hand-off's recorded base and head with pinned diff settings, computed locally and never fetched | Assured | A clean rebase keeps the fingerprint and any edit, whitespace included, changes it; binary, rename and mode changes count; empty, partial or unavailable input is refused (done 2026-11-15) |
| **1.6d** | Merge-aware `done`, part 4, from 2.1: review rounds: one session records every required profile; findings are major, minor or nit; the reviewer, its model, session and stream, and each build hold's model are recorded; a builder can't review its own task; rounds bind to the fingerprint | Assured | A round with its findings survives replay, export and V8; a builder's review is refused, even through an ended hold; "other model" and the Assured separate review are read from records, never from flags (done 2026-11-15) |
| **1.6e** | Merge-aware `done`, part 5, from 2.1: closing a major: a later round of the task resolves it or rules it out, or a person closes it with a reason | Assured | A major stays open across edits until a later independent round or a person's decision closes it; a builder's answer never does; replay, export and V8 keep the closures (done 2026-11-15) |
| **1.6f** | Merge-aware `done`, part 6, from 2.1: rework: a round asking for changes moves the task back to `in_progress`, and `next` gives it back to the builder's kept worktree as kind `rework` | Assured | After a round asking for changes, the task is fixed, checked, handed off again and approved by a later round, all through public commands (done 2026-11-15) |
| **1.6g** | Merge-aware `done`, part 7, from 2.3: `fid approve` bound to the change fingerprint, by an interactive person; the prompt shows the diffstat, the checks, the review verdicts and "no other-model review" when it applies | Assured | V6; agents and non-interactive runs are refused (done 2026-11-15) |
| **1.6h** | Merge-aware `done`, part 8: the single done evaluator, per lane (decisions record: 11 November 2026, item 8): the hand-off's checks on the merged PR head's tree with their Artifacts still kept, reviews with no open major finding, the Assured separate review (or the "no other-model review" flag) and Mark's current approval, the `## Rollback` section, the merged PR and green CI on its merge commit; the merge judge never required | Assured | A lane matrix refuses each missing, stale or mismatched fact and accepts complete Routine, Standard and Assured work; a lane raised after hand-off needs a new hand-off (done 2026-11-15) |
| **1.6i** | Merge-aware `done`, part 9: `fid done <task>`: asks `gh`, never `git fetch`; waits ("waiting for CI") while the merge commit has no CI run; records the merge and CI's run; releases the review reservation; `fid issue transition` refuses to complete a task; an open question is resolved as finished before an answer | Assured | A reviewed, approved, merged task with green CI completes only through `fid done`; a stale observation can't complete superseded work; while it waits the task stays in review (done 2026-11-15) |
| **1.6j** | Merge-aware `done`, part 10: `fid done` with no task reconciles every task in review, including PRs Mark merged in GitHub; `status` and `show` say "merged, waiting for CI" or "merged, CI not green" from recorded observations, offline | Assured | Repeated reconciliation is harmless; read commands work without `gh` (done 2026-11-15) |
| **1.6k** | Merge-aware `done`, part 11: a person closes an epic with `fid done <epic>`, recorded as a person's decision, refused while any of its tasks is open; an epic whose tasks are all finished is a NEEDS YOU item; `fid issue transition` refuses to complete an epic (decisions record: 11 November 2026, item 9; the epic's close settled 12 November 2026) | Assured | Agents can't close an epic; an open task, one in review included, blocks it; the NEEDS YOU item appears and clears (done 2026-11-15) |
| **1.6l** | Merge-aware `done`, part 12: V5's verifier, `check:done-refusals`, through public commands | Assured | V5 (done 2026-11-15) |
| **1.7** | Leak test for kernel words in human output; it replaces the copy of the kernel-word pattern in `internal/status`'s unit tests | Assured | V2; no test keeps its own copy of the pattern (done 2026-11-15) |
| **1.8** | Move the remaining reset tasks into fid | Mark | fid tracks its own work (a disposable ledger until 3.2) (done 2026-11-21) |
| **1.9** | Journal admission audit: every Event kind is admitted by one rule set for shape, exact ownership, cardinality and size, identical on live append, replay and import; SQLite case aliases fenced (Mark, 2026-11-16, after four near misses in 1.6) | Assured | A matrix of every Event kind and admission path is in the PR; forged and oversized Events are refused on every path and change nothing; the Later-list audit items from #73, 1.6b and 1.6d are resolved (done 2026-11-16) |
| **1.9b** | Admission owns task states: each fixed-state move of `issues.status` is admitted only from the Event of the command that owns it (`fid start`, the hand-off, a review round asking for changes, `fid done`, a person's cancel decision), on live append, replay and import; `issue.updated` can't change `status`; decide what Event-v1 history may still move (G10 in #93; decisions record, 2026-11-17, item 1) | Assured | A forged `issue.transitioned` or `issue.updated` can't move a task or epic into a state its command owns on any path, and changes nothing; every move fid's commands make is still admitted; the test fixtures that move tasks with plain transitions use the owning commands (done 2026-11-21) |
| **1.9c** | Row identity and cardinality for the Issue families: admission ties every row an Issue-family or link Event writes (acceptance criteria, issue labels, labels, links, activity rows) to the Event's own Issue and arguments, with the number of rows each stage may write; an Event can't rename its Issue or change an Issue other than a link's two ends; each check only import makes moves into the shared rules or is recorded as import-only with its reason (G11 in #93; decisions record, 2026-11-17, item 1) | Assured | Forged Events for each gap are refused on every path and change nothing; the PR lists each import-only check with its outcome |
| **2.1** | Work selection on 1.6d–1.6f's records: independent review work in `next`, after resume and rework and before unblocking and new work; atomic review reservations; the rework cap; approval and merge NEEDS YOU items after reviews approve (decisions record: Decided for 2.1, 2026-11-22) | Assured | V3, V4 and V5 still pass; `next` hands out review work to an eligible actor in order and never to two at once; the rework cap stops automatic rework; the approval and merge NEEDS YOU items appear only once reviews have approved (done 2026-11-22) |
| **2.1b** | Work decisions: accepting agent-created tasks and subtasks, skipped by `next` until a person accepts and shown under NEEDS YOU; a person allows additional rework rounds for one task (size split of review delivery, 2026-11-22) | Assured | Acceptance appears and clears; only a person allows more rounds; allowances preserve the count and return rework with guidance to the builder’s worktree (done 2026-11-22) |
| **2.1c** (proposed) | Review delivery: builder answers or defers minor/nit findings as discovered work; post rounds and finding decisions to the PR from canonical records, with pending posts and retry protection (remaining 2.1b scope) | Assured | Answers never close findings; deferral creates the task, link and disposition atomically; each comment is rendered from fid’s record and posted once; posting failures remain pending and warn offline |
| **2.2b** | Lane floor and hard floors (floor slice of the approved size split; decisions record, 2026-11-21) | Assured | fid computes a task's floor from committed changed paths, check uses the higher lane, hand-off raises and records it atomically, and show explains it (done 2026-11-21) |
| **2.3** | The approval presentation in `fid status` and `fid show`: missing, current and stale approvals, and "no other-model review" for Assured work none of whose reviews came from another model family (`fid approve` itself landed in 1.6g; decisions record, 2026-11-15) | Assured | Text and JSON show missing, current and stale approvals and the flag, offline; V6 still passes (done 2026-11-22) |
| **2.2c** | `fid waive`: fingerprint-bound person decisions; hard floors; current and stale display (remaining 2.2b scope, after 2.3 and reusing its current/stale display; Mark accepted the row 2026-11-22) | Assured | Only a person lowers the gates for a handed-off change, never below its hard floor; an edit makes the waiver stale; live append, replay, import and restore enforce the same rules |
| **2.3b** | Completion receipt in `fid show`: a done task shows what it completed with, the merged PR and merge commit, the CI run on it, the review rounds, the approval and the "no other-model review" flag (decisions record, 2026-11-22, item 2) | Standard | Text and JSON show a done task's receipt, offline |
| **2.4** | Loop contract (fields, exit codes, version) with golden tests, the version checked through the `contract_version` every JSON reply carries, since 1.5c removes `fid capabilities` (decisions record, 11 November 2026, item 6); scenario tests for rework, rework cap, interruption and parallel claims; `fid next --epic <epic>` scopes a worker to one epic (decisions record, 2026-11-17, item 2) | Standard | V4; the contract V9 checks; a scenario test where a worker scoped to a side epic builds it while an unscoped worker works the main line (done 2026-11-22) |
| **2.5** | First loop, as fid's own copy: contract check, loud failures, heartbeat, park on interrupt, backup after each task, merge judge in shadow (Mark records his call first); runs the review work `next` hands out as kind `review`, and passes the agent Mark's answer when a task that waited for him resumes (decisions record, 11 November 2026, items 5 and 11) | Standard | It runs a reset task end to end |
| **2.6** | Usage capture: tokens per task from `claude` and `codex` JSON output; `next` skips a task over its token budget, with the over-budget scenario test | Assured | Every reset task has a token count |
| **3.1** | Validations, `epic approve` (binds the validations), `brief`, merged digest; `fid done <epic>` also checks the epic's validations (decisions record, 11 November 2026, item 9) | Standard | V1's fixture recipe runs from `fid brief` |
| **3.2** | Squash the 38 migrations into a baseline; tidy `go.mod`, unless an earlier Assured task already has; release v0.2.0; create the permanent ledger; run a restore drill; restore rebuilds the projection from the archive's journal with the running fid, so an older archive restores and one from a newer fid is refused (decisions record, 11 November 2026, item 10) | Assured | A fresh clone restores the ledger from backup; `go mod tidy -diff` reports nothing |
| **3.3** | Kit: plain-copy install, `plan`, `tdd`, builder habit; the default loop promoted from 2.5 | kit | Installed in kauri; the loop self-test passes |
| **3.4** | kauri: archive its V1 docs, one-page AGENTS.md, `.fid.toml`, install the kit, pre-register and randomly assign 10 trial tasks | Mark | Ready for the trial |

Rows are in run order. Since 10 November 2026 (ADR-0014), 2.2a runs between
1.4b and 1.5a, and the kernel's removal (1.5b–1.5d) runs right after 1.5a.
On 11 November 2026, 1.4 was split into 1.4a and 1.4b, and 1.5 was renamed
1.5a. The same day Mark lowered 1.4b's lane to Standard, since it changes
display only, with no migration. The decisions record's "Decided 11
November 2026" lists that day's other decisions, for 1.5a to 3.2, and
"Decided 12 November 2026" two more, for 1.5a and 2.1, with the gaps
settled from them.
On 14 November 2026 Mark split 1.5a into 1.5a1–1.5a7, all Assured.
On 2026-11-14 the orchestrator split 1.5b-1.5d into 1.5b1-1.5d2, all Assured, under Mark's delegation.
On 2026-11-15 Mark added 1.5e after the test suites were found to leak temporary directories.
On 2026-11-15 the orchestrator split 1.6 into 1.6a-1.6l, all Assured, under Mark's delegation,
bringing forward the recording parts of 2.1 and 2.3 that `done` needs.
On 2026-11-15 Mark added 1.5f after a local CI run that crossed UTC midnight failed 47 tests.
On 2026-11-16 Mark added 1.9 after review caught four admission gaps in 1.6.
On 2026-11-17 Mark added 1.9b and 1.9c from 1.9's audit, both Assured. 1.9b runs before 2.1;
1.9c runs alongside phase 2. The same day he added the scoped `next` to 2.4.

Weeks 4–6 run the trial; it may slip to week 7.

## Validations

| ID | Claim | Verifier |
|---|---|---|
| V1 | Mark takes control of a stuck task in two commands (`fid status`, then `fid takeover`), in under five minutes, unaided | human, fixture recipe |
| V2 | No kernel nouns or IDs appear in human output | `check:leak-test` |
| V3 | Two worktrees share one ledger and never claim the same task; `next` respects dependencies | `check:two-worktree` |
| V4 | An interrupted task resumes in the same worktree from fid alone | `check:interrupt-scenario` |
| V5 | `done` refuses on a missing receipt, an open major finding or a stale approval, and accepts only merged work with CI green | `check:done-refusals` |
| V6 | An approval survives a clean rebase and goes stale when the change is edited | `check:approval-binding` |
| V7 | Working-tree edits to `.fid.toml` change nothing until merged | `check:policy-from-main` |
| V8 | Back up, wipe and restore yields identical state | `check:restore` |
| V9 | The loop refuses to start on a contract mismatch and stops loudly on an unknown `next.kind` | `check:loop-selftest` |
| V10 | Review catches the seeded defects | `review:test-adequacy`, seeded defects |

## Trial (kauri)

- **Tasks:** ten, pre-registered and randomly assigned, about five per arm.
- **Arm 1:** fid, the kit and the loop.
- **Arm 2:** the kit only, with the same skills and reviewer, run as
  interactive sessions with PRs and a short task note.
- **What it measures:** the whole system against the manual path.
- **Per task, record:**
  - Mark's active minutes;
  - tokens;
  - rework cycles;
  - workflow words;
  - defects found later;
  - "would you choose this again?"
- **Seeded defects:** about three per arm.

## Kill criteria

1. Workflow words per task are at least 70% below V1's (Phase 24: about 17k
   per task).
2. fid catches at least two false, stale or inadequate claims that kit-only
   misses.
3. The seeded-defect catch rate is at least equal to kit-only, and no Assured
   defect escapes.
4. Mark's minutes and tokens are no worse than kit-only.
5. A fresh agent resumes an interrupted task from fid alone, and Mark takes
   control of a stuck task in under five minutes, unaided.
6. Mark would choose fid again.

If these fail: keep the receipts core as a small library or hook pack, and
stop building fid as a product.

## Friction log

| Date | Task | Friction | Fix or later |
|---|---|---|---|
| 31 Oct | 0.2, 0.3 | The agent's shell has no SSH key loaded, so pushes go through `gh` over HTTPS, and `gh`'s token lacks the `workflow` scope needed to push CI changes | Mark grants the scope (`gh auth refresh -s workflow`) or loads the SSH key for agent shells |
| 31 Oct | 0.3 | GitHub Free has no branch protection on private repos | Deferred until `auto_merge` (decisions record, later list). Since 2026-11-14 the repository is public, and a `main` ruleset requires a pull request and blocks force pushes and deletion; required status checks aren't set yet |
| 2 Nov | 0.6 | The first review's `make test` timed out in `cmd/fid` (`panic: test timed out after 10m0s`) while the integration suite ran alongside it. The V1 developer profile said not to overlap long checks; the one-page AGENTS.md dropped the rule | AGENTS.md step 5 says to run the checks one at a time |
| 4 Nov | 1.1c1 | CI shard 0 timed out in `cmd/fid` (`panic: test timed out after 10m0s`, 604 s) twice. Three new test names reshuffled the shards (they're assigned by name position), and main's shards already took up to 586 s. Total test time didn't change | `ci_tests.py` passes `-timeout 20m`. Balancing shards by measured time, and reviewing whether every test earns its run time, are on the later list |
| 7 Nov | 1.1c1–1.1d5 | CI used about 1,650 of the 2,000 billed minutes GitHub Free gives Mark's private repos each month in one week: 40 runs at about 41 minutes each, mostly the race detector | 0.7 |
| 2026-11-12 | 1.4b, 2.2a | GitHub Actions minutes ran out: CI jobs were refused within 2-3 seconds ("recent account payments have failed or your spending limit needs to be increased"), so PRs #57 and #58 couldn't get CI | Mark made the repository public on 2026-11-14, so standard runners are free; "Require approval for all external contributors" is set |
| 2026-11-14 | 1.5a2 | WSL ran out of memory and crashed while a Codex builder and a Claude reviewer ran side by side (WSL capped at 8 GB, with Firefox using 3-4 GB in Windows). Both sessions were lost; the builder resumed from its saved session | Keep at most two heavy sessions at a time and leave Windows for development; lower WSL's memory cap in `.wslconfig` if it happens again |
| 2026-11-14 | 1.5a1-1.5a4 | Codex's sandbox made the worktree's git metadata (`.git/worktrees/<task>`) read-only partway through sessions, so builders couldn't commit | The orchestrator commits the builder's work with the builder's own message; the cause isn't known yet |
| 2026-11-14 | 1.5a3 | Tests that snapshot the ledger expected exactly one journal file, so they failed when a test ran across UTC midnight (13:00 NZDT) | Fixed in PR #61: the snapshots cover every dated journal file |
| 2026-11-14 | 1.5a | GitHub's line count for a stacked PR includes the earlier PR's commits until that PR merges (#61 showed +10,641; its own change is about 3,800 lines). Golden files are about 35–70% of the lines each 1.5a PR adds (#59: 71%, #60: 36%, #61: 50%), and about 4,000 lines per PR is too much for Mark to read in full | Later: Mark to decide on collapsing golden files in GitHub's diffs and a "what to read" section in PR descriptions |
| 2026-11-14 | 1.5a5 | Reads such as `fid show` and `fid status` failed with `journal-replay-stopped` when a write was in progress (15–22% of reads under a concurrent writer) | Fixed in this PR: catch-up decides again under the write lock |
| 2026-11-14 | 1.5a5, 1.5a6 | CI runs the tests with coverage and with no git identity; local `make test` does neither. #64 failed CI three times after its review approved it: a test binary run as a check's child process printed "GOCOVERDIR not set" into the output the tests compare, and test fixtures committed with no author | Fixed in #64 and #65 (test-only). Builders and fix rounds now also run CI's own runner locally with git's identity hidden (`ci_tests.py run`) |
| 2026-11-14 | 1.5a6 | Codex's sandbox can't see other sessions' processes, so the "wait while another test suite runs" loop inside a Codex session never waits; the loop also matched its own shell when `go test` was in the same command line | Run the wait loop as its own command; keep at most two heavy sessions so overlapping suites stay rare |
| 2026-11-15 | 1.5a-1.5d | The test suites leaked about 8,000 temporary directories (about 1.9 GB) into a RAM-backed `/tmp` over four days of builds and reviews | 1.5e |
| 2026-11-15 | 1.6f | `TestReworkGolden` failed CI on #86 (unrelated to it): it masked every 8-character hex word as `<commit>` before replacing the workspace path with `<ROOT>`, so a temp directory whose random name was eight digits (about 2% of runs) lost part of its path and never matched | Fixed in this PR: paths first, then commits, with a regression |
| 2026-11-15 | 1.5e, 1.6g, 1.6l | text goldens wrapped real temp paths before masking them, so wrap points varied with the temp directory's length (recorded three times) | fixed in this PR: fixture paths have a fixed display width |
| 2026-11-16 | 1.8 | `fid status` lists held-back tasks under READY. After the import, READY shows phase 3's two tasks marked "epic blocked by", and "+43 more" counts the Later list, which waits for 3.2. Seen in 1.8's dry run and again with the pinned binary on 2026-11-21 | Later display work: leave tasks whose epic is blocked out of READY, or count them apart |
| 2026-11-16 | 1.8 | `fid next --explain` lists every open task (53 after the import) and puts the blocked tasks that others wait on (2.1, 2.4, 2.5 and 3.2) above the three that are ready (1.9c, 2.2b and 2.3). It's the order `next` walks, but it reads as if the blocked tasks come first. Seen in 1.8's dry run and again on 2026-11-21 | Later display work: list ready tasks before blocked ones, and shorten the list |
