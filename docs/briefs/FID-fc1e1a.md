# Brief: FID-fc1e1a, loop polish (Standard)

Approved by Mark, 2026-11-28. fid task FID-fc1e1a, phase 2 (the golden
path). It has no RESET row, so commits use the footer
`Refs: RESET FID-fc1e1a`. fid's loop builds and reviews it, before 2.5b.
Where this brief and the repo disagree, park with the question instead of
guessing.

## Why

The loop's first full run (2.6, #122) worked end to end. Mark then noticed
three things the orchestrator's script used to do that the loop doesn't:

1. **PR titles and branch names.** The loop titled the PR with fid's task
   title, "2.6 Usage capture". The repo's PRs are titled like Conventional
   Commits with a RESET suffix: `feat(usage): capture agent session token
   usage (RESET 2.6)`. Its branch also repeated the row ID:
   `reset-2.6-2-6-usage-capture`.
2. **Hard-to-read review comments.** fid posts the reviewer's text fields to
   the PR as GitHub Markdown, exactly as written. The orchestrator's review
   prompt said so and asked for backticks; the loop's reviewer prompt
   doesn't, so its reviewers wrote plain prose with bare commands.
3. **No fix responses.** Under the orchestrator, each fix round's response
   to the findings appeared on the PR. Since 2.1c fid can do this itself:
   `fid finding answer` and `fid finding defer` record a response and fid
   posts it (`docs/specs/06_DESIGN_PLAN_CLI.md`, "Finding responses and
   delivery"). The loop never uses them.

Read first: `loops/` (loop, README, prompts, selftest),
`docs/specs/LOOP_CONTRACT.md`, the CLI spec's "Finding responses and
delivery" section, and `cmd/fid/integration/first_loop*_test.go`.

## Approved design

- **PR title.** When the loop creates a PR, its title is the subject of the
  first commit on the task branch (oldest first, from the default branch to
  `HEAD`) whose type is `feat`, else the first commit's subject, followed by
  ` (RESET <stem>)`. `<stem>` is the brief's file stem (`docs/briefs/2.6.md`
  gives `2.6`), or the task's fid ID when it has no brief. Rework keeps the
  title and refreshes the body, as now.
- **Branch name.** `reset-<stem>-<slug of the task title>`, where a title
  that starts with the stem drops it: "2.6 Usage capture" gives
  `reset-2.6-usage-capture`.
- **The reviewer prompt** (`loops/prompts/reviewer.md`) adds: "fid posts
  these fields to the PR as GitHub Markdown for a person to read: in every
  text field, put commands, file paths, test names, flags, identifiers and
  literal output in `backticks`." It asks for `verification` as a short
  Markdown list, one command per item with its result, and for concise
  evidence.
- **Finding answers.** The fixer's sandbox can't reach GitHub, so the fixer
  doesn't call fid for this; the loop does, as it records reviews.
  - For rework, the builder prompt asks the fixer to write an answers file
    at a new `LOOP_ANSWERS` path: a JSON list with one entry per open
    finding it addressed, either `{"finding_id": "...", "answer": "..."}`
    (what changed, with the commit, or why not) or, for a minor or nit left
    for later, `{"finding_id": "...", "defer": "<task title>"}`. The same
    Markdown and backticks rule applies to the text.
  - After the fixer's session and before `fid check`, the loop records each
    entry as the builder actor with `fid finding answer <id> --text ...` or
    `fid finding defer <id> --title ...`. fid posts each one on the PR. A
    deferral files a new task that waits for Mark's `fid accept`.
  - These commands take no retry key, so record each entry once: on a
    restart, skip a finding that already has this builder's response in
    `review list`'s declared `responses`.
  - A malformed answers file or a refused answer stops the loop loudly, like
    any other recording failure, because answers are review evidence. A post
    that fid leaves pending (GitHub unavailable) is printed and doesn't stop
    the loop.
- **Contract.** Declare `finding answer` (`finding.answer`) and
  `finding defer` (`finding.defer`) in `LOOP_CONTRACT.md`, with the result
  fields the loop reads and goldens. Added commands are additive, so the
  version stays `fid-cli.v1`.
- **Out of scope.** Minor findings left open after an APPROVE: fid serves no
  builder work after approval, so nothing answers them yet. Note it under
  Left over.

## Required work and evidence

1. TDD, failing tests first. The stub `gh` records `pr create` titles and
   PR comments. Journeys and selftests cover:
   - the title for a branch whose first `feat` commit isn't its first
     commit, for a branch with no `feat` commit, and for a task without a
     brief;
   - the branch name when the title starts with the stem;
   - the reviewer prompt carrying the Markdown guidance;
   - a CHANGES round where the fixer writes answers: fid records them
     (`review list` shows the builder's responses), the stub `gh` receives
     the answer comments, and a deferral files its task;
   - a malformed answers file stopping the loop loudly, and a restart not
     answering a finding twice.
2. A small mutation check: one temporary mutation per rule above (title
   source, suffix, stem stripping, answers recorded, deferral, duplicate
   skip), each caught by a test. Put the table in the report.
3. The decisions record: a short "Decided 2026-11-28" entry naming the three
   behaviours, from Mark's notes after 2.6's run. No RESET.md change: this
   task has no row.

## Real-ledger boundary

Every fid worktree shares Mark's live ledger under the git common directory.
Never run the installed fid, a freshly built fid or the loop with real fid in
fid's checkout, except the waiting park your prompt describes. Tests use
temporary repositories. A manual probe builds fid to a temporary path and
uses a temporary repository and a temporary `XDG_CONFIG_HOME`. Never run the
real gh, codex or claude in tests or probes, and push nowhere.

## Checks and completion report

Use focused tests while working. After the last code change, run each once,
one at a time: `make fmt-check`, `make lint`, `make test`, `make build`,
`make check-loop-selftest` and
`go test -count=1 -tags integration ./cmd/fid/integration/...`. A command
that didn't finish didn't pass.

The report (at `LOOP_REPORT`) covers what changed and why, the Standard
lane, failing-first evidence, the mutation table, each command with its last
output line, decisions and deviations, and anything left over. Commit with
Conventional Commits and the footer above, without a co-author line, and
leave the worktree clean.
