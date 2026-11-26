# fid's first loop

`python3 /path/to/fid/loops/loop.py` runs in the build worktree. Use Python
3.9 or later, git, fid, gh and the configured agent executable. Python uses
only the standard library. This is fid's own copy; RESET 3.3 promotes it to
the kit. `loop.py` must stay under 150 physical lines, including comments
and blank lines. The test-only Go package runs `selftest.sh` in `make test`
and CI; `make check-loop-selftest` runs the same guards and budget alone.

Before starting, initialise fid, configure a backup target and git identity, authenticate gh,
and refresh `origin/<default>` yourself (`git fetch origin`). The loop reads
that local policy and base; it never fetches. Set `origin/HEAD` to the default
branch. Briefs must already be on that branch, as described in
[docs/briefs/README.md](../docs/briefs/README.md). Start with a clean build
worktree for new work; interrupted work may have changes to resume.

Run one loop per machine. The loop does not serialise test suites; two full
suites at once slow each other, as described in AGENTS.md.

Every fid call uses `--json` and is checked against
[fid-cli.v1](../docs/specs/LOOP_CONTRACT.md). A contract mismatch, unknown
kind or waiting kind, missing declared field, invalid/empty JSON or unexpected
nonzero exit stops the loop with the reply quoted. Empty output never means
no work. Every batch `done` outcome is printed, including refusals and
unavailable GitHub observations; a batch's exit 0 does not imply completion.

A run proceeds as follows:

1. Read `next` to check the contract, then reconcile completed work with
   batch `done`. Keep a detached review worktree beside the build worktree,
   outside git metadata. Each build worktree has its own review worktree
   and outputs.
2. Peek in the review worktree. Claim only review work there as the reviewer;
   claim other work in the build worktree as the builder. `LOOP_EPIC` scopes
   both selections. fid keeps a handed-off task in its build worktree for
   rework, so the second worktree is necessary.
3. For `new`, start the task and create `reset-<brief stem, else fid ID>-<title
   slug>` from `origin/HEAD`. For `resume`, continue the kept work. For
   `rework` caused by `default-conflict`, merge the declared default ref
   first; the fixer resolves any merge conflicts. Pass the task's ID, title,
   acceptance criteria and spec path from `show`, the brief's contents,
   earlier findings from `review list`, the answer and rework allowance.
4. Run a fresh builder/fixer agent. Send an initial heartbeat, then renew at
   `heartbeat_interval_seconds` in a thread while the agent runs. A stale
   heartbeat terminates the agent process group and stops loudly. Give Codex
   writable roots for the worktree git directory, the common git directory
   (shared objects/refs), and this session's output directory. An agent must
   commit its changes and write a report; the loop never commits for it.
5. If the agent parked waiting for a person, leave the task parked and move
   on. Otherwise require a clean tree, run `fid check`, push to origin, open
   a draft PR with the report as its body (or reuse the branch's open PR),
   and hand off. Reused PRs get the latest report as their body. Builds must
   use a task branch, rather than main, the default branch or a detached head.
   A successful agent that leaves uncommitted changes stops
   the loop and parks an interruption.
6. For `review`, check out exactly the handed-off head detached in the
   review worktree, and run a fresh reviewer for the declared stream and
   every required profile. Its required checks have already passed on the
   hand-off head's tree; fid refuses hand-off otherwise. The reviewer reads
   the change, runs focused tests and probes with overlay mutations, and
   does not rerun full suites. Require a clean tree and the same head
   afterwards.
   Record its JSON with `fid review record`; fid posts the canonical PR
   comment. On APPROVE, mark a draft PR ready so CI runs; already-ready PRs
   stay ready across further reviews. Further review streams are selected
   by fid. Separate review context excludes primary rounds. Mark reads `fid status`, approves when needed and
   merges; this loop never merges.
7. Back up after each item, including a waiting park or interruption. At
   `none`, print the `waiting` object as the reason, reconcile with batch
   `done` again and exit 0. `waiting: null` means there is no declared wait
   for this worktree. Dependencies, caps and acceptance remain fid's rules.

Agent failure (including a usage limit) and Ctrl-C stop the agent's process
group, park the build hold as interrupted, back up and exit nonzero. The next
run resumes that work in its original worktree. Review holds lapse freely;
`park` is for build holds. Parking and backup failures are printed separately.
The loop stops on other command failures as well; inspect the quoted reply
and retained outputs before restarting. It does not blindly retry commands.

| Environment setting | Default / meaning |
|---|---|
| `LOOP_FID` | `fid`; a single executable path |
| `LOOP_GH_REPO` | `mbell-nzl/fid`; the explicit owner/repository passed to gh |
| `LOOP_EPIC` | Unset; scope both worktrees' `next` calls to this epic when set |
| `LOOP_BUILDER_ACTOR` | `agent:sol-builder` |
| `LOOP_REVIEWER_ACTOR` | `agent:sol-reviewer`; must differ from the builder |
| `LOOP_BUILDER_MODEL` | `openai/gpt-6.1-sol`; attribution as family/model |
| `LOOP_REVIEWER_MODEL` | `openai/gpt-6.1-sol`; same-model fallback, fresh session |
| `LOOP_BUILDER_COMMAND` | `codex exec -m {model} --add-dir {git_dir} --add-dir {common_dir} --add-dir {output_dir} -` |
| `LOOP_REVIEWER_COMMAND` | Same command; receives the reviewer prompt |
| `LOOP_REVIEW_WORKTREE` | Absolute path `<build worktree>-review`; retained between runs |

Agent command settings are split with `shlex`, then each argument is
formatted. They are executed directly, without a shell. `{model}` is the
model portion after the family; `{git_dir}`, `{common_dir}` and `{output_dir}`
are absolute paths. To use another provider, set its command and full model
attribution together. A wrapper may read the exported context variables;
stdin contains the prompt, context JSON and approved brief. Any usage-limit
wrapper must return nonzero. Agent stdout/stderr goes to `agent.log`.

Codex's sandbox writes only to its workspace and the `--add-dir` roots. If
the repository's `.codex/config.toml` points Go's caches outside the worktree,
add those configured directories to both command settings. For example, if
all those caches are under `/path/to/shared-caches`:

```sh
export LOOP_BUILDER_COMMAND='codex exec -m {model} --add-dir {git_dir} --add-dir {common_dir} --add-dir {output_dir} --add-dir /path/to/shared-caches -'
export LOOP_REVIEWER_COMMAND="$LOOP_BUILDER_COMMAND"
```

Each item retains `prompt.md`, `context.json`, `agent.log`, and the agent's
`report.md` or `review.json` under
`<build git directory>/fid-loop/runs/<session UUID>/`, outside tracked files.
The loop prints this directory before launching the agent. Agent environment
variables name the exact session and output files:

| Exported to agents | Meaning |
|---|---|
| `LOOP_ACTOR`, `LOOP_MODEL` | Actual actor and full model attribution |
| `LOOP_HOLD`, `LOOP_SESSION` | Opaque fid hold and fresh session UUID |
| `LOOP_CONTEXT` | JSON containing `work`, `task`, `reviews` |
| `LOOP_REPORT` | Builder/fixer's completion report path |
| `LOOP_REVIEW` | Reviewer's `fid review record --input` JSON path |

Reports and prompts are retained for diagnosis, not used as a second task
database. A restart reads fid's current state and recorded reviews again.

For tests or probes, isolate `XDG_CONFIG_HOME` and `XDG_CACHE_HOME`, provide
stub agents and gh on PATH, and use only temporary git and backup remotes.
Never run the installed fid, a built fid, or this loop against fid's own
checkout during builder tests: every worktree shares its live ledger in the
git common directory. `selftest.sh` and the integration journeys create their
own repositories; no real agent, gh or network service is invoked.
