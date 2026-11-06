# fid reset

Approved by Mark · 31 October 2026 · [ADR-0013](docs/adr/0013-reset-fid-to-a-ledger-with-receipts.md) ·
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
  - The loop stays under 150 lines.
  - The kit install script stays under 80 lines.
- **One task = one PR = one session.** Use TDD; `make test` must pass. The PR
  description is the completion report, and small decisions go there too.
- **Review.** A fresh session on the other model reviews each PR: Codex
  reviews Claude's work and Claude reviews Codex's. One pass covers code and
  test adequacy. Findings go in a PR comment. Once `fid review` works, fid
  records them and posts that comment itself.
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
| **0.3** | Add a `push: main` CI trigger. Branch protection is deferred: GitHub Free doesn't offer it on private repos | Assured | A merge to `main` runs CI (done 31 October 2026) |
| **0.4** | Move the CI scripts and conformance fixtures out of `project-docs/` | Assured | CI passes with `project-docs/` deleted locally (done 31 October 2026) |
| **0.5** | Archive fid's V1 docs: tag `v1-workflow-final`, copy them to the archive repo, keep the reset documents in `docs/reset/` | Routine | `project-docs/` is gone from `main`; CI passes (done 1 November 2026) |
| **0.6** | Kit repo with the reviewer agent (code and test-adequacy profiles) | kit | It reviews this reset's first PR (done 2 November 2026) |
| **1.1a** | Route every ledger path through one function, with no change in behaviour | Standard | Only `config.LedgerDir` joins the ledger onto a project root, and the existing tests pass unchanged (done 2 November 2026) |
| **1.1b** | Tests reach ledger files only through `config.LedgerDir`, with no change in behaviour | Routine | The 1.1a audit also covers test files and passes, and no expected value changed (done 3 November 2026) |
| **1.1c1** | Test workspaces are guarded git repos, with no change in behaviour | Assured | Every workspace the `cmd/fid` and integration tests give fid is a git repo whose top level is itself, except where a test hides `.git` on purpose; `GIT_CEILING_DIRECTORIES` is set for those test processes; no expected value changed (done 4 November 2026) |
| **1.1c2** | fid needs a git repository: discovery goes through git | Standard | Outside a git repository every command that opens or creates a project, `fid init` included, refuses; the ledger is still at `<root>/.fid` (done 4 November 2026) |
| **1.1d1** | Every test workspace that reaches the ledger is a guarded git repo, `internal/` packages included, with no change in behaviour | Standard | Every root a test gives `config.LedgerDir` is a git repo whose top level is itself, except where a test puts a ledger outside a repository on purpose; the guard helpers live in one test-only package that every test process making git repos uses, with `GIT_CEILING_DIRECTORIES` set; no expected value changed (done 4 November 2026) |
| **1.1d2** | Output names ledger files by absolute path; the CLI spec updated | Standard | Text and JSON output (`journal_path`, `journal_file`, `lock_path`, `fid init`'s message) name ledger files by absolute path; the ledger is still at `<root>/.fid` (done 5 November 2026) |
| **1.1d3** | Move the ledger to the git common directory, with staging beside it; the config spec updated | Standard | Two worktrees share one ledger (V3's claim and `next` parts belong to 1.4) (done 5 November 2026) |
| **1.1d4** | Discovery refuses when git and `LedgerDir` disagree, or when git could track the ledger | Standard | fid never writes outside `--git-common-dir`/fid, nor where git could track it: where git and `LedgerDir` name different ledgers, or git's common directory is inside the worktree but not within a `.git` directory, every command that opens or creates a project refuses with `repository-unavailable` / `unsafe-path` |
| **1.1d5** | A `.fid/` in a worktree is an ordinary path | Standard | A leftover `.fid/` makes the worktree dirty and can be adopted like any other path; the `dirty` help and the CLI spec no longer single it out |
| **1.2** | `fid backup` (export, commit, push), restore test, backup-age status line | Standard | V8 |
| **1.3** | Display package; `fid status`, `fid show` | Standard | whet-style output; golden tests for text and JSON |
| **1.4** | `fid next --claim`: atomic, respects dependencies, ordered, typed skip reasons, `--explain` | Standard | V3 |
| **1.5** | `start`, `check`, `park`; claim heartbeat | Standard | An agent works a task without touching kernel commands |
| **1.6** | Merge-aware `done` with the single done evaluator; reconciles PRs Mark merged in GitHub | Assured | V5 |
| **1.7** | Leak test for kernel words in human output | Routine | V2 |
| **1.8** | Move the remaining reset tasks into fid | Mark | fid tracks its own work (a disposable ledger until 3.2) |
| **2.1** | Review and finding surface; one session records every required profile; each round is posted to the PR; `done` refuses open major findings | Standard | V5 |
| **2.2** | Policy from `origin/<default>` with built-in defaults; lane floor; hard floors | Assured | V7 |
| **2.3** | `approve` bound to the change fingerprint; stale approvals shown | Assured | V6 |
| **2.4** | Loop contract (fields, exit codes, version) with golden tests; scenario tests for rework, budget, rework cap, interruption and parallel claims | Standard | V4, V9 |
| **2.5** | First loop, as fid's own copy: contract check, loud failures, heartbeat, park on interrupt, backup after each task, merge judge in shadow (Mark records his call first) | Standard | It runs a reset task end to end |
| **2.6** | Usage capture: tokens per task from `claude` and `codex` JSON output | Standard | Every reset task has a token count |
| **3.1** | Validations, `epic approve` (binds the validations), `brief`, merged digest | Standard | V1's fixture recipe runs from `fid brief` |
| **3.2** | Squash the 38 migrations into a baseline; release v0.2.0; create the permanent ledger; run a restore drill | Assured | A fresh clone restores the ledger from backup |
| **3.3** | Kit: plain-copy install, `plan`, `tdd`, builder habit; the default loop promoted from 2.5 | kit | Installed in kauri; the loop self-test passes |
| **3.4** | kauri: archive its V1 docs, one-page AGENTS.md, `.fid.toml`, install the kit, pre-register and randomly assign 10 trial tasks | Mark | Ready for the trial |

Weeks 4–6 run the trial; it may slip to week 7.

## Validations

| ID | Claim | Verifier |
|---|---|---|
| V1 | Mark takes control of a stuck task in two commands, in under five minutes, unaided | human, fixture recipe |
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
| 31 Oct | 0.3 | GitHub Free has no branch protection on private repos | Deferred until `auto_merge` (decisions record, later list) |
| 2 Nov | 0.6 | The first review's `make test` timed out in `cmd/fid` (`panic: test timed out after 10m0s`) while the integration suite ran alongside it. The V1 developer profile said not to overlap long checks; the one-page AGENTS.md dropped the rule | AGENTS.md step 5 says to run the checks one at a time |
| 4 Nov | 1.1c1 | CI shard 0 timed out in `cmd/fid` (`panic: test timed out after 10m0s`, 604 s) twice. Three new test names reshuffled the shards (they're assigned by name position), and main's shards already took up to 586 s. Total test time didn't change | `ci_tests.py` passes `-timeout 20m`. Balancing shards by measured time, and reviewing whether every test earns its run time, are on the later list |
