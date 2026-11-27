# fid's first loop

`python3 /path/to/fid/loops/loop.py` runs in the build worktree. Use Python
3.9 or later, git, fid, gh and the configured agent executable. Python uses
only the standard library. This is fid's own copy; RESET 3.3 promotes it to
the kit. The loop has no physical line limit; its prompts stay short. The test-only Go package runs `selftest.sh` in `make test`
and CI; `make check-loop-selftest` runs the same guards alone.

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
`usage record` is the one non-gating call: its failures are printed and the
loop continues, because usage is insight rather than an evidence gate.

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
   slug>` from `origin/HEAD`, dropping a leading stem from the title before
   making the slug (`2.6 Usage capture` becomes `reset-2.6-usage-capture`).
   For `resume`, continue the kept work. For
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
5. After every agent session, successful or failed, record its clean stdout
   through `fid usage record --tool codex|claude --tool-output <file>` with
   the session ID and actor, before recording
   a review round or parking a failed build. Before starting the agent, ask
   its executable for `--version` and pass a successful response as
   `--tool-version`; version-probe failure leaves that flag out. fid prefers
   a version embedded in Claude stream output. Usage refusals and failures
   are printed without stopping the loop.
6. If the agent parked waiting for a person, leave the task parked and move
   on. Otherwise require a clean tree. After rework, including a resumed fix
   with open findings, read the fixer's `LOOP_ANSWERS` JSON list and record
   each answer or minor/nit deferral through fid as the builder. Validate the
   whole file first; malformed or missing files and refused responses stop
   loudly. Commands have no retry key, so skip findings whose `review list`
   responses already include this builder, even after a partial recording
   and restart. fid posts the comments; print pending delivery warnings and
   continue. Deferrals create tasks awaiting Mark's `fid accept`.
   Run `fid check`, push to origin, open
   a draft PR with the report as its body (or reuse the branch's open PR),
   and hand off. A new PR's title uses the oldest `feat` commit subject in
   `origin/<default>..HEAD`, or the oldest commit subject when there is no
   `feat`, plus ` (RESET <brief stem, else fid ID>)`. Reused PRs keep their
   title and get the latest report as their body. Builds must
   use a task branch. Keep that branch across the agent session and check it
   again immediately before pushing; main, the default branch, a detached
   head or an agent switching branches stops the loop.
   A successful agent that leaves uncommitted changes stops
   the loop and parks an interruption.
7. For `review`, check out exactly the handed-off head detached in the
   review worktree, and run a fresh reviewer for the declared stream and
   every required profile. Its required checks have already passed on the
   hand-off head's tree; fid refuses hand-off otherwise. The reviewer reads
   the change, runs focused tests and probes with overlay mutations, and
   does not rerun full suites. Require a clean tree and the same head
   afterwards.
   Record its JSON with `fid review record`; fid posts the canonical PR
   comment. The reviewer's text fields are GitHub Markdown, with backticks
   around commands, paths, test names, flags, identifiers and literal output;
   verification is a short list of commands and results.
   At startup and when idle, read `status` for tasks awaiting a
   person's approval or merge, then read their latest hand-off PR from `show`.
   Mark those draft PRs ready so CI runs. This includes tasks requiring no
   review and retries a failed ready action on restart, after required
   reviews have been recorded. Already-ready PRs stay ready. Review streams are selected
   by fid. Separate review context excludes primary rounds. Mark reads `fid status`, approves when needed and
   merges; this loop never merges.
8. Back up after each item, including a waiting park or interruption. At
   `none`, print the `waiting` object as the reason, reconcile with batch
   `done` again and exit 0. `waiting: null` means there is no declared wait
   for this worktree. Dependencies, caps and acceptance remain fid's rules.

Agent failure (including a usage limit) and Ctrl-C stop the agent's process
group, even after its leader has exited, park the build hold as interrupted,
back up and exit nonzero. The next
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
| `LOOP_BUILDER_COMMAND` | `codex exec --json -m {model} --add-dir {git_dir} --add-dir {common_dir} --add-dir {output_dir} -` |
| `LOOP_REVIEWER_COMMAND` | Same command; receives the reviewer prompt |
| `LOOP_BUILDER_USAGE` | `codex`; `claude` or `none` also accepted |
| `LOOP_REVIEWER_USAGE` | `codex`; independent of the builder format |
| `LOOP_REVIEW_WORKTREE` | Absolute path `<build worktree>-review`; retained between runs |

Agent command settings are split with `shlex`, then each argument is
formatted. They are executed directly, without a shell. `{model}` is the
model portion after the family; `{git_dir}`, `{common_dir}` and `{output_dir}`
are absolute paths. To use another provider, set its command and full model
attribution and usage format together. For Claude use
`claude -p --output-format json` or
`claude -p --output-format stream-json --verbose`, with `LOOP_*_USAGE=claude`. A wrapper may read the exported context variables;
stdin contains the prompt, context JSON and approved brief. Any usage-limit
wrapper must return nonzero. Agent stdout goes to `agent.stdout` and stderr to `agent.stderr`, so
diagnostics cannot contaminate usage JSON.

Codex's sandbox writes only to its workspace and the `--add-dir` roots. If
the repository's `.codex/config.toml` points Go's caches outside the worktree,
add those configured directories to both command settings. For example, if
all those caches are under `/path/to/shared-caches`:

```sh
export LOOP_BUILDER_COMMAND='codex exec --json -m {model} --add-dir {git_dir} --add-dir {common_dir} --add-dir {output_dir} --add-dir /path/to/shared-caches -'
export LOOP_REVIEWER_COMMAND="$LOOP_BUILDER_COMMAND"
```

Each item retains `prompt.md`, `context.json`, `agent.stdout`, `agent.stderr`, and the agent's
`report.md`, `answers.json` or `review.json` under
`<build git directory>/fid-loop/runs/<session UUID>/`, outside tracked files.
The loop prints this directory before launching the agent. Agent environment
variables name the exact session and output files:

| Exported to agents | Meaning |
|---|---|
| `LOOP_ACTOR`, `LOOP_MODEL` | Actual actor and full model attribution |
| `LOOP_HOLD`, `LOOP_SESSION` | Opaque fid hold and fresh session UUID |
| `LOOP_CONTEXT` | JSON containing `work`, `task`, `reviews` |
| `LOOP_REPORT` | Builder/fixer's completion report path |
| `LOOP_ANSWERS` | Fixer's JSON list: `{"finding_id":"...","answer":"..."}` or `{"finding_id":"...","defer":"<task title>"}` per addressed finding; `[]` when none |
| `LOOP_REVIEW` | Reviewer's `fid review record --input` JSON path |

Reports and prompts are retained for diagnosis, not used as a second task
database. A restart reads fid's current state and recorded reviews again.
Minor findings left open after APPROVE receive no fixer session yet and
therefore no automatic answer.

For tests or probes, isolate `XDG_CONFIG_HOME` and `XDG_CACHE_HOME`, provide
stub agents and gh on PATH, and use only temporary git and backup remotes.
Never run the installed fid, a built fid, or this loop against fid's own
checkout during builder tests: every worktree shares its live ledger in the
git common directory. `selftest.sh` and the integration journeys create their
own repositories; no real agent, gh or network service is invoked.
