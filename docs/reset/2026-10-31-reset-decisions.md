# fid and Kit Reset: Decisions

Decisions session, 30–31 October 2026, run as a grilling session against
`2026-10-29-fid-kit-reset-proposal.html` (now superseded). Two Opus 5.5 second
opinions informed D2, and a fresh-context review of the result led to the
amendments marked below. These decisions feed ADR-0013 and `RESET.md`.

## Goal

> Highly automated agentic engineering that keeps quality high: a software
> factory, not a slop factory.

Success bar: a markedly better, more robust dev workflow with higher-quality
outputs than V1. Mark is the only user; the point is learning and building
taste.

Principles:

- **Earned autonomy.** Automation is earned by evidence and withdrawn by
  evidence.
- **Judgement escalates, never relaxes.** LLM judgement (plan skill, reviewer,
  merge judge) may raise a lane or hold work. Nothing but Mark lowers a floor
  or overrides a refusal.
- **Refuse only on missing evidence.** fid refuses when a receipt is missing, a
  major finding is open or an approval is stale. Numeric guides (diff size,
  validation count) are signals for judgement, never refusals.
- **Small, right-sized work.** One task = one PR = one coherent, reviewable
  change. Cheap capture keeps work small.
- **Strict scope rule** (amended). Nothing joins the reset unless a
  golden-path step or a kill criterion needs it. Everything else goes on the
  later list.

## Review outcome (amendments)

The fresh-context review found the process reset sound but the tool reset
drifting: the session had added about six mechanisms that pay off only after
the trial, and "facade first" keeps the whole kernel. Mark accepted all nine
recommendations:

1. Merge autonomy: judge in shadow only; Mark records his own call before
   reading the judge's; a manual per-lane `auto_merge` key replaces counters
   and trip rules; only the `caused_by` sensor now.
2. One rule engine: the facade's `done` evaluator decides what done needs; a
   net Go line budget; unused kernel commands deleted at trial end.
3. Reviews and approvals bind to the change; checks bind to the tree. No delta
   re-approval view.
4. Kit install is a plain copy. Three-way merge deferred.
5. Backups run outside fid's write path.
6. One reviewer session runs all required profiles per round.
7. Discovered work follows D6 only; triage is a skip reason, not a state.
8. Kill criteria restored and given quality measures.
9. The proposal is superseded, not revised; `RESET.md` is written from these
   decisions.

## Products

**D1 Kit V2: frozen.** Tag and archive the `v2` branch; `v1.0.0` stays for
kauri until it migrates. The kit becomes a light home for curated, adapted
skills, the reviewer agent and a default loop. No install, upgrade or rollback
machinery. The V2-2 ClaimNext engine is scrapped; `fid next --claim` does that
job. Mark owns the reviewer profiles (especially test adequacy), lane policy
and validations.

**D2 Two repositories, one-way dependency.**

- The new kit starts in a fresh repo. The old starter-kit repo stays untouched
  as kauri's V1 reference and is archived after kauri migrates.
- **The loop lives in the kit; fid owns the loop contract.** fid declares,
  versions and golden-tests the fields and exit codes loops read
  (`next --claim` exit codes, `.result.task.id`, `next.kind` values, the
  required review profiles). Scenario tests (rework, budget exhausted, rework
  cap, interruption, parallel claims) are fid integration tests. The kit ships
  one default loop; variants only when a project needs one. The first loop is
  written in fid's repo as fid's own copy, then promoted to the kit.
- Loop guard-rails: check fid's contract version at startup (replaces a fid
  version pin); fail loudly on an unknown `next.kind`; never treat empty output
  as "no work"; `loops/selftest.sh` against stub agents; a line budget enforced
  by a test.
- **Install is a plain copy** (amended). A kit script of about 50 lines
  overwrites managed files (whole skill folders in `.agents/skills/`, relative
  symlinks from `.claude/skills/`), renders the reviewer to
  `.claude/agents/*.md` and `.codex/agents/*.toml`, and records the kit commit
  in `.kit-version`. Copies are committed, so `git diff` and `git add -p` are
  the merge. Example loops are copied once and owned by the project.
  User-level skills only for a few general non-dev skills whose names never
  clash with the kit's.

**D9 Facade first.**

- New golden-path commands are built from existing kernel commands.
- **One rule engine** (amended). The facade's `done` evaluator (policy from
  `origin/main` plus the diff) is the single source of what done needs. Kernel
  Obligations cover receipts only. This sidesteps the kernel's Envelope
  revision-1 limit.
- A golden-output test fails if human-facing output (including refusals)
  mentions kernel-only nouns (Envelope, Activation, Obligation, Attempt,
  Attention Item) or kernel IDs (`ENV-`, `STP-`, `ATT-`). `--json`,
  `--verbose` and `fid admin` are exempt. A concept that keeps leaking is the
  evidence to collapse or rename it.
- Kernel commands move under `fid admin`, hidden from main help. At trial end,
  any the facade doesn't use are deleted.
- A net Go line budget applies (set in `RESET.md`).

## Foundations

**D4 / F1 Ledger location: outside git.** The ledger lives at
`$(git rev-parse --git-common-dir)/fid/`, shared by every worktree, never
committed. ADR-0013 supersedes ADR-0001's "journal committed to git" and
"merges handled by git merge" clauses. One machine only.

- Cloud sessions: fid unavailable there for now. Later path: read-only
  restore, then single-writer checkout through the backup repo, then
  merge-safe events only if local and cloud agents ever write concurrently.
- **Backups are required before dogfooding** (amended to run outside fid's
  write path). A private GitHub repo (default `fid-ledgers`) with a folder per
  project, keyed by a project ID in `.fid.toml`. `fid backup` = export, commit
  and push. The loop runs it after each task and a daily cron runs it too;
  `done` and `approve` never push. `fid status` shows backup age: a warning
  after 24 hours or a failed push, a NEEDS YOU item after 3 days. Losing a day
  is acceptable; losing a week is not. An automated test proves restore (back
  up, wipe, restore, compare), plus one manual drill before dogfooding. The
  backup target is in user-level fid config. `fid init` asks for it;
  `fid status` warns until it is configured; not a hard block.
- **Details** (decided 2 November 2026, for task 1.1):
  - fid needs a git repository. Outside one, `fid init` and every other
    command refuse. There is no `.fid/` fallback.
  - Output names ledger files by absolute path, in text and in JSON
    (`journal_path`, `journal_file`, `lock_path`), because a linked
    worktree's ledger sits outside that worktree.
  - No migration. fid ignores a `.fid/` directory left in a worktree; moving
    it to `.git/fid/` by hand keeps it.
  - Test workspaces are git repos. Test helpers check that git reports the
    temp directory itself as the repository, and set
    `GIT_CEILING_DIRECTORIES`, so a test can't reach a real repository's
    ledger.
  - Task 1.1 is split: 1.1a routes every ledger path through one function
    with no change in behaviour; 1.1b moves the ledger.
- **Details** (decided 3 November 2026, for task 1.1b):
  - 1.1b is split again, because its test changes alone would bury the
    behaviour change: about 400 hand-built ledger paths in 127 test files,
    and every test workspace becoming a git repo.
    - 1.1b: tests reach ledger files only through `config.LedgerDir`, with
      no change in behaviour.
    - 1.1c: fid needs a git repository. Discovery goes through git, every
      command refuses outside a repository, and test workspaces become
      guarded git repos. The ledger stays at `<root>/.fid`.
    - 1.1d: the ledger moves to the git common directory.
  - `LedgerDir` reads the worktree's `.git` entry in Go: a directory, or a
    `gitdir:` file plus `commondir`. It stays cheap and returns no error.
    Discovery runs `git rev-parse` once and refuses if git and `LedgerDir`
    disagree, so fid never writes outside `--git-common-dir`/fid.
  - `fid init` and `export import` stage in `<common>/fid-staging-*` and
    rename into `<common>/fid`. The journal and config packages gain entry
    points that take a ledger directory, since a staging directory has no
    `.git`.
  - The V1 closure journey asserts that the journal is committed to git.
    ADR-0013 reverses that, so 1.1d changes the assertion to "never
    tracked" rather than deleting it.
  - The ledger keeps its `.gitignore` (`fid.sqlite`) for now, though it means
    nothing under `.git/`. Dropping it would change the ledger's contents.
- **Details** (decided 4 November 2026, for task 1.1c):
  - 1.1c is split, for the same reason as 1.1b: the test changes would bury
    the behaviour change.
    - 1.1c1 (Routine, raised to Assured): test workspaces become guarded
      git repos, with no change in behaviour. Mark also agreed to raise the
      CI `go test` limit in the same PR, which makes it Assured (see the
      later list).
    - 1.1c2 (Standard): discovery goes through git, and fid refuses outside
      a repository.
  - Today the CLI never walks up. Every command passes the working
    directory to `project.Resolve` as if it were `--project`, so the
    directory must directly contain `.fid`.
  - The project root is the worktree's top level, from one
    `git rev-parse --show-toplevel` run from the start directory. The ledger
    is `<toplevel>/.fid` and nowhere else, and fid works from any
    subdirectory of the worktree. A `.fid` in a subdirectory is no longer
    found.
  - `--project <path>` works like `git -C <path>`: discovery starts at
    `<path>`, and the root is that path's top level. `export import
    --repository` follows the same rule.
  - `fid init` from a subdirectory initializes `<toplevel>/.fid`. The default
    prefix comes from the top level's name.
  - Outside a repository fid refuses with `repository-unavailable`
    (precondition, exit 4), with details `reason: not-a-repository` and
    `detail`. In a repository without a ledger it stays `project-not-found`
    (exit 3).
  - Every command that opens or creates a project refuses outside a
    repository: those that resolve one today, the three that skip
    resolution (`concern list`, `journal inspect`, `sql`), `init` and
    `export import`. Commands that touch no project (`--help`, `completion`,
    `export validate`) still work there.
  - `project.root` in `init --json` and import's `target_root` resolve
    symlinks, because git reports the real path.
  - `GIT_CEILING_DIRECTORIES` stops git climbing above a directory, but git
    still inspects the directory it starts in. So no test runs fid with a
    real repository's root as its working directory.
- **Details** (4 November 2026, overnight, for task 1.1d):
  - 1.1d is split into four parts. The builder split it overnight, under
    Mark's advance approval for a necessary or beneficial split; Mark has
    yet to review it.
    - 1.1d1: every test workspace that reaches the ledger becomes a guarded
      git repo, with no change in behaviour. That's the `internal/`
      packages' workspaces, plus a few integration fixtures that build a
      ledger in-process. The guard helpers move into one test-only package,
      `internal/gittest`, which every such test process uses.
    - 1.1d2: output names ledger files by absolute path, and the CLI spec
      follows, while the ledger is still at `<root>/.fid`. Doing this before
      the move keeps every part's output true.
    - 1.1d3: the ledger moves to the git common directory, with staging
      beside it, and the config spec follows.
    - 1.1d4: discovery refuses when git and `LedgerDir` disagree. No code
      or reason is recorded for that refusal, so Mark picks one before it
      starts. Keeping it apart stops that decision holding up the move.
- **Details** (decided 6 November 2026, for tasks 1.1d4 and 1.1d5):
  - Mark accepted the 1.1d split above when he merged its PRs (#38, #39
    and #40).
  - fid also refuses when git's common directory sits inside the worktree,
    other than within a `.git` directory. In that layout (for example
    `git init --separate-git-dir=<worktree>/metadata`), `git add --all`
    would track the ledger, along with git's own `HEAD`, `config` and
    hooks. Ignore rules don't make it safe, because `git add -f` still adds
    an ignored file. Codex found this in #40's second review round. 1.1d4
    adds it.
  - Both 1.1d4 refusals are `repository-unavailable` (precondition, exit 4)
    with the existing reason `unsafe-path`, and details naming the
    locations involved. No new code or reason. #40's "Left over" section
    has implementation notes for 1.1d4.
  - A `.fid/` in a worktree becomes an ordinary path, because no V1 ledgers
    remain in worktrees. Dirty status, reconciliation and current-event
    path validation stop treating `.fid/` paths specially, and the `dirty`
    help and the CLI spec stop mentioning them. That's task 1.1d5, which
    also updates the 1.1a audit test's comment that the ledger "is moving".

**D8 / F2 Parallel worktrees: minimal.** Shared ledger, atomic
`fid next --claim` on different issues, one loop per worktree, and a
two-worktree test. Mark expects mostly one sequential loop.

- Claims last about an hour. The loop sends a heartbeat during long agent
  sessions, since a builder may not touch fid for an hour. Tasks waiting on
  Mark are parked and hold no claim.
- Interruptions (crash, usage limit): the loop parks the task as interrupted
  (`fid park`) and stops. `fid next` serves it first, to the same worktree,
  with a linked successor Run. A worktree holding an interrupted task gets no
  other task.
- Shared-file conflicts are resolved at merge, as human developers do; the
  second PR's builder rebases as rework. ADR-0006 disjoint-scope proofs wait
  for real need.

**F3 Trust anchor: policy in `.fid.toml`, read from `origin/main`.** Checks,
lanes and check-defining paths live in one committed file. fid reads it only
from `origin/<default branch>` after a fetch, never from the working tree or
a local branch an agent could commit to. Built-in defaults apply until the
file exists there. fid's code makes changes to `.fid.toml` and
`.github/workflows/**` always Assured; the default template also makes the
Makefile Assured. Test files and test data require a test-adequacy review in
every lane. Activation stays underneath with one meaning, "the policy fid
trusts, from `origin/main` at commit X", visible only under `--verbose`, and
is a collapse candidate after the trial. ADR-0013 amends Gate B3. The
environment allowlist stays.

**D5 / F4 Threat model: honest but fallible agents** on a single-user
machine. Not defending against deliberately adversarial agents. A deceptive
agent could forge a local approval; that is accepted. Friction: a hook denies
agents `fid approve`, `waive` and `admin`, and `gh pr merge`; fid refuses
approval without an interactive person and refuses `agent:*` actors; the
approval prompt shows the diffstat, checks and review verdicts. atlas is out
of scope for now.

**D10 Approvals go stale** (amended). Reviews and approvals bind to the
change's fingerprint (`git patch-id` of the branch diff); checks bind to the
tree, and CI re-runs them on the final commit. A clean rebase keeps an
approval valid; any edit to the change voids it (a rebase that alters lines
next to a hunk may too, which means re-approving). Re-approval shows the full
diffstat; Assured tasks are small. Epic approval binds to the validation list
only, since builders update briefs and specs in their PRs.

**D7 / F5 Done = merged into `main` with CI green on that commit,** for every
lane. CI gains a `push: main` trigger so the merge commit gets a run. Local
receipts gate the hand-off from builder to reviewer. `fid done` confirms the
merge and CI through `gh`, records the CI run as the trusted receipt and
records the merge commit against the task. It also reconciles PRs Mark merged
in the GitHub UI. If GitHub can't be reached, the task shows "merged, waiting
for confirmation".

## Work model

**Merge authority: structured, earned autonomy** (amended).

| Lane | Needed before merge | Merged by |
|---|---|---|
| Routine | Local checks; test-adequacy review if tests changed | Mark, until `auto_merge` is enabled for the lane |
| Standard | Plus a code review (other model) and the judge's call | Mark, until `auto_merge` is enabled for the lane |
| Assured | Plus a second-model review, Mark's approval and a rollback note | Mark |

- **Through the trial, Mark merges everything.** A merge judge (fresh
  frontier-model session, MERGE or HOLD, hold-only) runs in shadow and is
  recorded as a review with profile `merge`. Mark records his own call before
  reading the judge's, so agreement isn't anchored.
- **Enabling agent merges** is a manual per-lane `auto_merge` key in
  `.fid.toml`. Changing it is an Assured change, so it needs Mark. When on, the
  loop script (never an agent) merges with
  `gh pr merge --squash --match-head-commit <sha>`. Branch protection on
  `main` (required CI, up-to-date branches, no bypass) is a precondition for
  turning on `auto_merge` in any lane. GitHub Free doesn't offer protection
  on private repos, so it waits until then, when Mark chooses GitHub Pro or a
  public repo (decided 31 October 2026).
- **One reviewer session per round** runs every required profile (code, test
  adequacy, security). Assured work adds a separate review on the other model.
- **Each review round is also a PR comment** (decided 3 November 2026, for
  tasks 2.1 and 2.5). Mark reads reviews where he merges, next to the
  completion report. fid's record of the review is canonical, because `done`
  checks findings there. The comment is built from that record, never written
  separately, so the PR and fid can't disagree. fid posts it through `gh` when
  the round is recorded, one comment per round. The posting is part of fid,
  not the loop, so interactive sessions and the loop post the same comment.
  It is fid's third use of the network (ADR-0013, Consequences).
- **Digest:** `fid brief --merged-since …`, one line per merged task, holds on
  top.
- **Sensors now:** escaped defects via `caused_by` links. Seeded defects run in
  the trial. The rest (judge disagreement, code-health trend, rework and hold
  rate, automatic dial) are on the later list.
- While Mark merges, loop throughput is capped by his merge latency for any
  task that depends on unmerged work. That's expected.

**D6 Lanes and validations.**

- Lanes are per task. Routine = no production code changed. Standard = any code
  change. Assured = listed risky or irreversible paths (trust boundary, journal
  format, migrations, `go.mod`, releases, `.fid.toml`, CI workflows, Makefile).
- A lane is set at planning (plan skill proposes; Mark sees it at epic
  approval) and raised at merge if the computed floor is higher. The reviewer
  and judge may raise it with a reason. Only Mark lowers it, through
  `fid waive`.
- Up to 10 validations per epic (soft cap; more means split the epic). Each
  names a verifier (`check:`, `review:` or human); human ones need a fixture
  recipe.
- No numeric task size. Size is a signal the judge weighs, never a reason to
  hold on its own. The plan skill slices risky work into its own small Assured
  task, so the PRs Mark reads stay small.
- Discovered work (amended): work that blocks a task in an approved epic joins
  that epic (`blocks`, `discovered_from`), appears in the digest and keeps its
  own lane. Other agent-created work is not served until Mark accepts it:
  `fid next` skips it with a typed reason, and `fid status` shows a NEEDS YOU
  line. Tasks Mark creates with `fid new` are ready immediately.
- `fid next` order: resume interrupted (same worktree); rework; tasks that
  unblock approved work; unstarted tasks by epic order and priority. It skips
  claimed, blocked, waiting-for-Mark, untriaged and unapproved-epic tasks,
  each with a typed reason for `--explain`.

## Rollout

**D3 Phase 24:** closed under DEC-293 (PR #26).

**D11 Dogfooding.** fid tracks its own work from week 1; that ledger is
disposable until the migration baseline, and `RESET.md` stays the source of
truth until then. Squash fid's 38 migrations into one baseline; compatibility
obligations start there. Dogfood with a pinned release binary, released at
each epic close. Order: fid, then kauri from week 3 (hosting the trial), then
yoyu after the trial if fid survives its kill criteria. atlas is out of scope.

**V1 docs archive.** One private repo (for example `v1-workflow-archive`) with
a folder per project and a README naming the source commit. Per project: tag
`v1-workflow-final`, copy the V1 material out, delete it from `main` in one PR,
keep specs, ADRs and the glossary, and write a one-page AGENTS.md. Agents are
not pointed at the archive. History is stripped only if a repo goes public.
fid first, after moving the CI scripts and conformance fixtures that live
under `project-docs/` (amended); kauri in week 3; yoyu after the trial.

**D12 Phases 25 and 26: cancelled as phases.** Into the reset: parallel
claims, usage capture, CLI polish (CF-040). Cancelled: kit adapter
qualification (CF-037). Later: generic retry-exhausted, the Envelope
revision-1 limit, and dead-code cleanup.

**D13 Trial** (amended). Two arms in kauri: fid + kit + loop, and kit only
(same skills and reviewer, interactive sessions, PRs, a short task note). Ten
pre-registered tasks, about five each, assigned at random. The kit-only arm
can't run the loop, so the trial compares the whole system against the manual
path; that is the question Mark's success bar asks. No competitor and no V1
arm; V1's figures come from the archive (Phase 24: about 260k words over 15
tasks). About three new seeded defects per arm replace `run-drills.sh`, which
retires with V1. Expect weeks 5–7.

Kill criteria:

- Workflow words per task at least 70% below V1's archived figure.
- fid catches at least two false, stale or inadequate claims that kit-only
  misses.
- Seeded-defect catch rate at least equal to kit-only, and no escaped Assured
  defects.
- Mark's active minutes and tokens no worse than kit-only.
- A fresh agent resumes an interrupted task from fid alone, and Mark takes
  control of a stuck task in under five minutes, unaided.
- Mark would choose fid again for the next task.

**D14 External skills.** Fill kit slots, not whole libraries:

| Slot | Starting point |
|---|---|
| plan | `task-pack` + `validation-contract`, merged and slimmed |
| tdd | `tdd`, plus Pocock's limit (no TDD for wiring or configuration) |
| verify-app | pstack's create-verification-skill, adapted |
| builder habit | Superpowers' verification-before-completion, adapted |
| reviewer profiles | QA Reviewer prompt + Osmani's review axes and code-health criterion |
| debugging and orientation | `diagnose`, `zoom-out`, `architecture-review` |
| merge judge | New |

`phase-planning`, `phase-wrap-up` and `qa-evidence` retire with V1.
`grill-with-docs` stays user-level. Curation is a separate exercise: refresh
the clones in `~/projects/external` (`mattpocock-skills` is from July), add
pstack, Superpowers and Osmani's skills, and review existing borrowed skills
against their current upstream versions. A skill enters the kit only after it
has helped on a real task, and records its source, licence, and what Mark
changed and why. Order: the reviewer in week 0, then `plan`, then `tdd` and
the builder habit, then `verify-app`; each is tried on the next reset task.

## Corrections to the proposal

- §14 loop sketch reads `.task.id`; fid-cli.v1 wraps results in `.result`, so
  the loop would silently exit as "no work". It also calls `fid run abandon`;
  interruptions use `fid park`.
- §4 golden path still lists `fid check add`; F3 replaced it with policy in
  `.fid.toml`.
- §19 week 3 "pin a fid version range" is replaced by the loop's contract
  check.
- §18 superseded-records list misses ADR-0001's merge and commit clauses.
- §15 "approval on the PR under your GitHub account": GitHub doesn't allow
  approving your own PR; this needs a bot identity for agents (later list).
- §2 `run-drills.sh` is not suitable for measuring the reset; most drills lint
  V1 completion reports.
- D7: CI runs only on pull requests today, so no run exists for the merge
  commit.
- Archiving `project-docs/` would break CI: `ci.yml` runs
  `project-docs/tools/ci_tests.py`, and the P21/P22 conformance tests read
  `project-docs/docs/phases/phase-2{1,2}/public-interface.md`.
- fid's `sync-agents.sh` and `sync-skills.sh` only generate files inside one
  repo, and `sync-skills.sh` copies only `SKILL.md`. Reuse only the awk
  renderer.

## Later list

- Automatic autonomy dial: per-lane agreement counters, trip rules, and the
  remaining slop sensors.
- Kit install three-way merge (`git merge-file`) and `--check`, once a local
  edit actually gets lost.
- Delta view for re-approval.
- A kit skill that helps build a project-specific loop, after two or three
  real loops exist.
- Bot identity for agents with required review from Mark on Assured paths.
- Branch protection on `main` for fid and kauri, before any lane gets
  `auto_merge` (needs GitHub Pro or a public repo).
- Cloud-session support for the ledger (read-only, then single-writer).
- Declared paths per task, if shared-file conflicts become a regular cost.
- Collapse candidates after the trial: Activation, plus any concept the leak
  test keeps catching.
- Generic retry-exhausted, the Envelope revision-1 limit, dead-code cleanup,
  and replay reporting `line: 0` for a current-event prior conflict (found
  in 1.1d2, #39).
- yoyu migration (after the trial).
- CI test time (from 1.1c1, 4 November 2026). `ci_tests.py` assigns tests
  to shards by their position in a sorted list of names, so adding one test
  reshuffles every shard. Each shard runs its `cmd/fid` tests as one
  `go test` process under race and coverage, taking up to about 10 minutes.
  1.1c1 raised the limit to 20 minutes as a stopgap. Two things to do:
  balance shards by measured time, and review whether every test earns its
  run time. By CI elapsed time (inflated by parallel contention), the
  biggest are the descriptor and judgment corpora at 400-540 s each, and
  the grammar corpora at 150-250 s each. Fold in the tests that fail when
  `TMPDIR` is behind a symlink: about 33 (12 in `internal/`, 11 in
  `cmd/fid`, 10 in integration), which fail on `main` too. They compare
  paths without resolving symlinks; P24 observation tests, for example,
  hash an unresolved worktree path. Found in 1.1c2; no PR since added one.
- Stale ledger staging (from 1.1d3, 6 November 2026). A crash during
  `init` or `export import` can leave `<common>/fid-staging-*` behind.
  `git status` doesn't show it, and nothing removes it. Add cleanup once
  one turns up; it needs a rule for when staging is safe to delete.
- CHANGELOG: it hasn't been kept up during the reset. Catch it up in 3.2,
  with the v0.2.0 release.

## Next steps

1. ~~Mark reviews and accepts ADR-0013 and approves `RESET.md`.~~ Done 31 October 2026.
2. Week 0 tasks from `RESET.md`.
