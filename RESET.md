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
| **1.5b** | Remove the kernel, part 1: before any kernel code goes, rebuild V8, the export fixtures and every other test that outlives the kernel on lean check results, with their helpers; V8 compares `next` in place of `ready` among its read commands; with 1.5c, move `fid artifact`, which stays, onto the lean check results (decisions record, settled 12 November 2026) | Assured | V8 and the export tests pass on lean check results, none weakened; no test that outlives the kernel creates a kernel record |
| **1.5c** | Remove the kernel, part 2: its commands, domain code and their tests; `fid ready`, which `fid next` replaces; and Concerns, requirement and decision references, `fid capabilities`, `fid stats` and the whole of `fid workflow` (`define`, `assign` and `list`) (decisions record: 11 November 2026, item 6; `list` settled 12 November 2026); `fid artifact` stays, moved onto the lean check results with 1.5b | Assured | No kernel command or domain package remains, nor any other this row removes; every remaining test passes, none weakened |
| **1.5d** | Remove the kernel, part 3: its journal and export catalog entries, maintenance checks, manifest-v1 validation and tables; and the tables and journal and export catalog entries of Concerns and of requirement and decision references, which 1.5c removes (decisions record, settled 12 November 2026) | Assured | No kernel table remains, nor any table of Concerns or of requirement and decision references; a journal or archive holding any of their events is refused as `schema-unsupported`; every remaining test passes, V8 and the export tests included, none weakened |
| **1.6** | Merge-aware `done` with the single done evaluator: per lane, the reviews, Mark's approval and the Assured `## Rollback` section the decisions record lists (11 November 2026, item 8); the merge judge never required; waits ("waiting for CI") while the merge commit has no CI run; asks `gh`, never `git fetch`; reconciles PRs Mark merged in GitHub; records CI's run on the merge commit as a check result; `fid issue transition` refuses to complete a task or an epic; a person closes an epic with `fid done <epic>`, recorded as a person's decision, refused while any of its tasks is open, and an epic whose tasks are all finished is a NEEDS YOU item (item 9; the epic's close settled 12 November 2026) | Assured | V5 |
| **1.7** | Leak test for kernel words in human output; it replaces the copy of the kernel-word pattern in `internal/status`'s unit tests | Routine | V2; no test keeps its own copy of the pattern |
| **1.8** | Move the remaining reset tasks into fid | Mark | fid tracks its own work (a disposable ledger until 3.2) |
| **2.1** | Review and finding surface; one session records every required profile; each round is posted to the PR; `done` refuses open major findings; a task whose latest review round approved it and that isn't done is a NEEDS YOU item in `fid status`; findings are major, minor or nit; a builder can't review its own task; only a later round or a person closes a major; the rework cap; accepting discovered work; fid records which model built the change and which model reviewed each round; `next` hands out review work as kind `review`, to an actor that didn't build the task, after resume and rework and before unblocking and new work (decisions record: 11 November 2026, item 11; 12 November 2026, item 13) | Assured | V5; that NEEDS YOU item |
| **2.2b** | Lane floor; hard floors; `fid waive` | Assured | fid computes a task's floor from its changed paths; only a person lowers a lane, with `fid waive`, never below a hard floor, and an edit to the change makes the waiver stale |
| **2.3** | `approve` bound to the change fingerprint; stale approvals shown; when no review of an Assured change came from a model other than the builder's, the approval view and `fid show` say "no other-model review" | Assured | V6 |
| **2.4** | Loop contract (fields, exit codes, version) with golden tests, the version checked through the `contract_version` every JSON reply carries, since 1.5c removes `fid capabilities` (decisions record, 11 November 2026, item 6); scenario tests for rework, rework cap, interruption and parallel claims | Standard | V4, V9 |
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
