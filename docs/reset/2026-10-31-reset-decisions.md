# fid and Kit Reset: Decisions

Decisions session, 30–31 October 2026, run as a grilling session against
`2026-10-29-fid-kit-reset-proposal.html` (now superseded). Two Opus 5.5 second
opinions informed D2, and a fresh-context review of the result led to the
amendments marked below. These decisions feed ADR-0013 and `RESET.md`.
Amendments dated 10–14 November 2026 record the lean core
(ADR-0014). "Decided 11 November 2026" and "Decided 12 November 2026" below
list Mark's decisions of those days for tasks 1.4b to 3.2; "Decided 14
November 2026" records the policy answers and the 1.5a split. Decisions
and gaps settled from them are also marked where their topics live.

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
   net Go line budget; unused kernel commands deleted at trial end. (Amended
   10 November 2026: the golden path is built on a lean core, and the whole
   kernel chain is removed right after task 1.5a. See D9 and ADR-0014.)
3. Reviews and approvals bind to the change; checks bind to the tree. No delta
   re-approval view.
4. Kit install is a plain copy. Three-way merge deferred.
5. Backups run outside fid's write path.
6. One reviewer session runs all required profiles per round.
7. Discovered work follows D6 only; triage is a skip reason, not a state.
8. Kill criteria restored and given quality measures.
9. The proposal is superseded, not revised; `RESET.md` is written from these
   decisions.

## Decided 11 November 2026

Mark made these decisions on 11 November 2026, for tasks 1.4b to 3.2. Each
is also marked in the section its topic belongs to, and the tasks' rows in
`RESET.md` point here.

1. **1.4b's lane is Standard.** Mark lowered it: 1.4b changes display only,
   with no migration. ADR-0014 had kept it Assured, the lane of the 1.4 it
   was split from, since only Mark lowers a lane.
2. **The reset's own Assured PRs keep one review per PR**, as AGENTS.md
   step 7 says, unchanged, until 1.6's `done` enforces the work model's
   separate Assured review on fid's own work.
3. **Task workflow states** (built in 1.5a). A task's states are
   `todo → in_progress → in_review → done`, plus `cancelled`. Rework goes
   from `in_review` back to `in_progress`. fid's commands move tasks between
   them: `start`, the hand-off, a review round asking for changes, and
   `done`. `dev_complete`, `qa_review`, `qa_block`, `failed_gate` and
   `qa_pass` are dropped. Epics become `todo`, `in_progress`, `done` and
   `cancelled`: `wrapped` is dropped, and `done` takes the place of today's
   `closed`. Subtasks are unchanged. Nothing needs to stay compatible, since
   no real ledger exists before 3.2.
4. **Takeover** (1.5a). `fid takeover <task> --actor human:<name> --reason
   "<why>"` ends any hold on the task, live or lapsed, records the person's
   decision (kind `takeover`) and gives the person a hold in their current
   worktree. The old session's next heartbeat is refused, so its loop stops.
   V1's two commands are `fid status`, then `fid takeover`. fid never touches
   git state, so the old worktree's changes stay where they are.
5. **Waiting for Mark** (1.5a). The agent runs `fid park --waiting
   "<question>"`. `next` skips the task as `waiting-for-you`, and the task
   stays kept for that worktree. `fid status` shows the question under
   NEEDS YOU, with `fid answer <task> "<text>"`. The answer is a person's
   decision of a new kind, `answer`, which fid refuses from `agent:*`. That
   worktree's next `next --claim` resumes the task with the answer in the
   JSON, so the loop passes it to the agent.
6. **1.5c also removes:**
   - Concerns, which 2.1's findings and discovered issues replace;
   - requirement references and decision references;
   - `fid capabilities`, since 2.4 checks the loop contract through
     `contract_version`, which every JSON reply carries;
   - `fid stats`, which reads V1 run telemetry; 2.6's usage capture
     replaces it;
   - `fid workflow define` and `fid workflow assign`, since item 3 fixes
     the states.

   These stay: `issue`, `label`, `link`, `show`, `status`, `next`, `init`,
   `backup`, `restore`, `export`, `journal`, `rebuild`, `sql`,
   `maintenance` and `activity`. This answers ADR-0014's open question of
   whether Concerns, requirement and decision references and
   `fid capabilities` go with the kernel.
7. **No `fid admin`.** Nothing would be under it once the kernel is gone.
   D5's agent hook still denies agents `fid approve`, `fid waive` and
   `gh pr merge`, and only `--json` and `--verbose` stay exempt from D9's
   leak test.
8. **What `done` requires** (1.6). As in the work model's table, each lane
   adds to the one before:
   - Every lane: the PR is merged; CI is green on the merge commit; and
     every check the hand-off policy required passed on the PR head's tree.
   - Routine: plus a test-adequacy review when test paths changed.
   - Standard: plus a review round with no open major finding.
   - Assured: plus the separate review (or the "no other-model review"
     flag), Mark's approval bound to the change's fingerprint, and a
     `## Rollback` section in the PR description, which `done` checks for.
   - The merge judge stays in shadow through the trial. It's never
     required.
   - With no CI run on the merge commit, `done` waits ("waiting for CI").
     There's no key to opt out until a project needs one.
   - `done` asks `gh`. It never runs `git fetch`, which settles whether
     `fid done` fetches (ADR-0014; F3).
9. **Epics close by a person** (1.6), with `fid done <epic>`, which is
   refused while any of the epic's tasks is open. An epic whose tasks are
   all finished shows under NEEDS YOU. 3.1 later adds the epic's
   validations to that check.
10. **Restoring across a fid upgrade** (3.2). Restore ignores the archive's
    projection and rebuilds it from the archive's journal with the running
    fid, so any older archive restores. An archive from a newer fid is
    refused. This settles the later-list item. ADR-0014's clean cut still
    refuses a journal that holds the kernel's events.
11. **Review work in `next`** (2.1). `next` hands out review work as kind
    `review`, to an actor that didn't build the task, so a crash between
    the hand-off and the review loses nothing. The 2.5 loop uses it.

## Decided 12 November 2026

Mark made three more decisions on 12 November 2026, numbered on from the day
before. Each is also marked in its topic's section and its task's row in
`RESET.md`.

12. **A worktree whose task waits for Mark gets no other work** (1.5a).
    While its task is parked waiting for Mark (item 5), that worktree's
    `next` returns `none`, saying it's waiting on him for that task, until
    he answers. Other worktrees carry on. This never mixes two tasks'
    changes in one worktree, since fid never touches git.
13. **Review work's place in `next`'s order** (2.1). `next` hands out
    resume, then rework, then review, then work that unblocks other work,
    then new work. Finishing open work beats starting new work, so PRs
    don't pile up waiting for review. Review work comes as kind `review`
    (item 11).

14. **Epic fields in 1.4b.** `fid show` adds an epic's `Scheduling` row:
    `sequential` reads `one at a time, in sequence`, and `prioritised` reads
    `by priority`. Its JSON gains `scheduling` on epics, null otherwise.
    READY reuses the `epic_blocked_by` field added to `next` in PR #56.

The orchestrating session settled the gaps below the same day. They follow
from Mark's decisions as transactional details, and aren't his own words:

- **`fid workflow list`** goes in 1.5c with the rest of `fid workflow`,
  following item 6: with the states fixed (item 3), no configurable
  workflows remain to list.
- **`fid artifact` stays**, following item 6. 1.5b and 1.5c move it onto
  the lean check results, which keep their Artifacts.
- **1.5d drops** the tables and the journal and export catalog entries of
  Concerns, requirement references and decision references, as it does
  the kernel's, following item 6.
- **An `answer` decision** binds to the task and the question it answers,
  following item 5.
- **Closing an epic**, following item 9. `fid issue transition` refuses to
  complete an epic, as it does a task. `fid done <epic>` records the close
  as a person's decision, of a kind fid refuses from `agent:*`.

## Decided 14 November 2026

Mark approved the following for 1.5a1. The answers to 2.2a's P1–P3 are
recorded under F3 as decisions, replacing the pending proposals.

- **Split 1.5a into seven Assured PRs, in order:** 1.5a1 policy decisions
  and status health; 1.5a2 heartbeat, interruption and takeover; 1.5a3
  waiting, answers and missing worktrees; 1.5a4 fixed states, start and
  task lanes; 1.5a5 named checks and Artifacts; 1.5a6 lane-required checks
  and which result counts; 1.5a7 hand-off and the worktree during review.
  `RESET.md` holds each task's scope and completion criteria.
- **Status reads trusted policy, read-only.** It reads the same local
  `origin/<default>` policy as writing commands, or built-in defaults
  when there is no origin or file. It never writes a snapshot or any
  other Event, including after the effective digest changes.
- **Status keeps working when policy is unavailable.** It shows
  the whole board and one `Policy unavailable: <problem>` warning, with
  the fix. It shows no checks or test-path warnings when it cannot read
  policy. Otherwise, it warns about lanes whose `checks` key is missing,
  in Routine, Standard, Assured order, and about a missing
  `[review] test_adequacy` key. Explicit empty lists are set up and do not
  warn. Each fix names the relevant tables, then merging `.fid.toml` and
  `git fetch`. Warnings sit inside the box just above the backup line,
  policy unavailable, checks, then test paths; they appear in empty
  projects too, use display wrapping and indentation, and are never
  NEEDS YOU items. JSON always carries `result.warnings` (`[]` when
  healthy), each with `kind`, `summary` and `next_action`;
  `policy-unavailable` adds `problem`, nullable `ref` and nullable
  `commit`; `checks-not-configured` adds the lowercase `lanes` list;
  `test-paths-not-configured` has no extra fields. The CLI spec records
  the exact text and JSON shapes.
- **Git remotes cannot be listed.** The orchestrator decided this under
  Mark's delegation for the 1.5a series on 14 November 2026. If running
  `git remote` itself fails, keep `policy-unavailable`, with `problem`
  naming the git error (`couldn't list git remotes: <error>`), null `ref`
  and `commit`, and `check the repository: git remote -v` as the repair
  action. Status shows the same warning and fix. This is its own source
  error case, distinct from an unknown default or an unreadable policy file.

Mark also approved these details for **1.5a2**, the same day:

1. **Retry keys (Q18 B).** The agent's new commands require
   `--idempotency-key`: `park` now, `start`, `check` and hand-off later.
   A person's `takeover`, `answer` and cancellation accept an optional key.
   A same-key repeat returns the recorded response and original Event.
   Without a key a person's repeat is refused by state, in plain words.
   Heartbeats take no key. Suggested commands for a person omit the key.
2. **Commands name the hold.** `heartbeat` and `park` take `--hold <id>
   --actor <actor>` and find the task from that exact hold. Under the write
   lock it must exist, be open and unexpired, belong to that actor, and
   name the command's current worktree. A refusal never substitutes a
   different hold. Takeover names the task instead.
3. **Heartbeat every five minutes.** Every accepted heartbeat records one
   Event, sets `last_heartbeat_at` to its time and `expires_at` to an hour
   later. Initially the last heartbeat is `since`; `since` never changes.
   At or after expiry the heartbeat is refused. The worktree resumes a
   started task with `next --claim`, under a linked new hold. A lapse's
   "work stopped at" time in show and next is now the last heartbeat.
4. **Parking needs an explicit mode.** Exactly one of `--interrupted
   --reason "<text>"` or, from 1.5a3, `--waiting "<question>"`. In 1.5a2,
   bare park refuses with `invalid-argument` naming `--interrupted
   --reason`. The free-text reason is required. An interrupted park ends
   the hold. Started work stays kept for that worktree, which resumes it
   first under a linked hold; other worktrees skip it as
   `interrupted-elsewhere`. Unstarted work goes back to anyone. Show says
   "interrupted in <worktree> by <actor>: parked at 20:00: Usage limit";
   next says "Interrupted here: parked at 20:00: Usage limit." A lapse
   keeps "work stopped at 20:00 without a park".
5. **Takeover is one atomic write.** `fid takeover <task> --actor
   human:<name> --reason "<why>" [--idempotency-key <k>]` ends the old
   hold, records a person's `takeover` decision bound to the task, and
   gives that person a linked hold in their current worktree. A live,
   lapsed, parked or absent old hold is allowed. Workflow state is kept;
   git changes stay in the old worktree. A worktree already holding or
   keeping another task refuses with `worktree-busy`; a person already
   holding the task here gets `already-held`. Closed tasks refuse with
   `transition-not-permitted`, unsatisfied `task-open`; other issue kinds
   use `unsupported-kind`. Agents and services get `authority-denied`.
   The old session's heartbeat or park gets `hold-not-current`, `ended`.
6. **Kept work is explicit in JSON.** Show, next queue entries and status
   wherever it has `working` also carry `kept` (null otherwise):
   `{worktree, here, mode, reason, since, hold_id, actor}`. With no live
   hold, a started task's lapse or interrupted park produces mode
   `interrupted`; reason is null for a lapse, and since is its last
   heartbeat or its park time. Working still describes an open hold,
   including a lapse. Text reuses "interrupted in <dir>" in status and
   READY. Waiting and review modes belong to later tasks.
7. **Separate conflict codes, exit 4, not retryable.** `hold-not-current`
   has exactly `hold_id` and `reason` (`not-found`, `ended`, `expired`,
   `actor-mismatch`, `worktree-mismatch`) and suggests `fid next` in human
   output. `worktree-busy` has `task_id` and `worktree`; `already-held`
   has `hold_id`. Their subject is the task.
8. **Cancellation stays on issue transition.** For tasks and epics,
   `fid issue transition <id> --to cancelled --actor human:<name>
   --reason "<why>" [--idempotency-key <k>]` requires a person and a
   reason, records a `cancel` decision bound to the subject, and ends
   its hold or releases its kept worktree in the same write. JSON adds
   `decision_id` to the existing result. Subtasks keep their current
   behavior. The decision record has an ID, kind, actor, subject, reason,
   binding and time; agents and services cannot write any kind. Takeover
   and cancel land now; other kinds land with their commands.

Builder's storage choices for 1.5a2: migration 22 extends `holds` with the
last heartbeat, an interrupted park's reason, and a `kept` flag. The same
unique indexes cover open holds and kept reservations. Releasing a parked
reservation keeps its ended hold as history. A lapse stays derived from
the clock. `person_decisions` uses `PDN-` IDs. `hold.heartbeat`,
`hold.parked` and `task.taken-over` join both Event catalogs; takeover and
`issue.transitioned` carry the decision Effect in the same Event as their
other Effects. No separate decision Event splits the transaction. The
built-in policy records `cancel` in the person-decisions action set, with
rules version 2, so snapshots retain the new authority rule.

Mark approved these details for **1.5a3**, on 14 November 2026:

1. **Parking to wait.** `fid park --hold <id> --actor <a> --waiting
   "<question>" --idempotency-key <k>` ends the exact current hold and
   always keeps its task for its worktree, started or not. Exactly one
   mode is required: `--interrupted --reason "<text>"` or `--waiting
   "<question>"`. Bare park names both forms in its refusal.
2. **One open question per task.** Each question has its own `QST-` ID,
   text, asker and time. Answered questions remain as history.
3. **Answering.** `fid answer <task> "<text>" --actor human:<name>
   [--idempotency-key <k>]` binds to the task's one open question, without
   a `--question` flag. It records a person's decision of kind `answer`,
   bound to the task and question, takes no hold and changes no workflow
   state. Agents and services get `authority-denied`. A task without an
   open question gets `not-waiting` (conflict, exit 4, not retryable,
   subject the task, details `{task_id}`). Same-key repeats replay the
   original answer and Event.
4. **Scheduling and repeated answers.** The waiting worktree's `next`
   returns kind `none`, with `waiting` naming the task and question, and
   gives it no other task until answered. Other worktrees skip the task
   as `waiting-for-you` and carry on. After answering, its next hand-out
   is kind `resume`, with a new hold linked to the parked one and an
   `answer` object. Repeated `next` calls keep that same answer, including
   after the live hold resumes; only parking with the next question
   replaces it. `waiting` and `answer` are null when absent. Review's
   waiting kind belongs to 1.5a7.
5. **Display.** Status has a `waiting-for-you` NEEDS YOU item, subject
   the task, naming the question and `fid answer <task> "<answer>"
   --actor human:<name>`. Next here says `Nothing here until you answer
   <task>: <question>`, with `fid show <task>`. Kept work has mode
   `waiting` with the question. Show's text displays the open question;
   its JSON lists the task's question and answer history.
6. **Takeover and cancellation.** A person taking over waiting work
   gets a hold with the question still open, and may answer it.
   Cancellation resolves waiting and retains the question's history.
7. **Missing worktrees need certainty.** For interrupted or waiting
   work kept for a worktree, status checks local git registration and
   membership of this repository's common directory, never fetching.
   Confirmed missing or replaced worktrees produce `worktree-missing`
   under NEEDS YOU, subject the task, with `fid takeover <task> --actor
   human:<name> --reason "<why>"`. A failed check produces a
   `worktree-unverified` warning with `task_id` and `worktree`, saying
   `Couldn't check the worktree kept for <task>: <worktree>` and `look
   at <worktree>, then fid status`. Neither result releases kept work.

Builder's storage choices for 1.5a3: migration 23 adds `questions`, with a
unique open question per task and an answer decision link, and extends
holds' end reasons with `waiting`. Questions are inserted with the ended
hold in one `hold.waiting` Event (`hold.park`); an answer decision and
question resolution share one `task.answered` Event (`task.answer`).
Cancellation resolves its question in the existing `issue.transitioned`
Event. Both catalogs register the families; replay and import validate
question identity, attribution, times and decision bindings. Existing
`PDN-` decision IDs stay; the examples' decision IDs are illustrative.
The answer's binding ID is the question and its subject is the task.
Question history uses journal insertion order, including multiple questions
in the same second. The answer's original worktree is stored in its retry
response metadata so a later retry retains its original human message.
Status keeps its single projection snapshot and READY shares next's planner.
The open question also keeps unstarted work after a person's takeover if
that session lapses or is parked as interrupted. That worktree stays busy;
replay refuses another task releasing its unanswered reservation. This
follows the always-kept rule, without another workflow state or flag.

Mark approved these details for **1.5a4**, on 14 November 2026:

1. **Fixed states.** Tasks use `todo → in_progress → in_review → done`,
   plus `cancelled`. Review-driven rework returns from `in_review` to
   `in_progress` in 2.1. Epics use `todo`, `in_progress`, `done` and
   `cancelled`; subtasks are unchanged. The retired task states and epic
   `wrapped` / `closed` go, without conversion: no real ledger exists.
   `workflow assign` refuses tasks and epics, whose states are fixed.
2. **Planning lanes.** `issue create --lane` and `issue update --lane`
   accept `routine`, `standard` and `assured` on tasks only. A missing lane
   stays null. Agents may set a missing lane or raise one, never lower it
   (`authority-denied`). A person may lower a lane while the task is `todo`,
   as a planning correction; after starting, lowering requires 2.2b's
   `waive` and is refused with `transition-not-permitted` until then.
   `show` gains its Lane row; issue read/create JSON carries `lane`.
3. **Starting.** `fid start --hold <id> --actor <a> --idempotency-key <k>`
   checks the exact current hold as the other session commands do. It needs
   a lane and a task in `todo`, then moves it to `in_progress`. A missing
   lane gives `lane-not-set` (precondition, exit 4, subject the task,
   details `{task_id}`). `in_review` cannot be resumed through `start`:
   rework belongs to 2.1. Start needs readable trusted policy and records a
   snapshot when its effective digest changes. It neither changes git nor
   requires a clean tree.
4. **Epic movement is atomic.** The first task's start moves its todo epic
   to `in_progress` in the same write. A person cancels an epic; `done`
   closes it in 1.6. Nobody manually starts an epic.
5. **Generic transitions cannot replace the commands.** Tasks cannot start
   by `issue transition`, move to review that way, or move from a started
   state back to `todo`. The refusal is `transition-not-permitted`, naming
   `use-start`, `use-hand-off` or `task-already-started`. Epics cannot be
   manually moved to `in_progress`. The completion refusal lands in 1.6;
   cancellation keeps 1.5a2's person-only rule. Subtasks keep plain
   transitions.
6. **No change base at start.** The `when_changed` base is recomputed for
   each check and hand-off (1.5a6/1.5a7); start records no base.

**Harmless start retry keys** (orchestrator's decision under Mark's
delegation for the 1.5a series, 14 November 2026): starting a task already
in progress with a current hold succeeds without a new Event. A fresh
retry key on that harmless repeat is deliberately left **unbound**. The
repeat is idempotent by state: while the task stays in progress it returns
"already started" with `event: null`; once the task moves on, that same key
gets whatever the state then gives. No retry storage or Event is added for
this no-op. The key of an Event-producing start still replays its recorded
response and original Event as usual.

Builder's storage choices for 1.5a4: migration 24 adds the nullable lane and
replaces the task/epic workflow seeds. Lane set and raise use the existing
`issue.created` and `issue.updated` families. One `task.started` Event
carries the task and, when needed, its epic's state and revision updates;
both catalogs register it from schema 24. The retry response binds it to
the exact hold and lane. Replay and import check that binding and enforce
the planning lane authority rule. Kernel tests attach an isolated legacy
test workflow through the journal; its first started state is `building`.
There is no production bypass or new dependency.

Mark approved these details for **1.5a5**, on 14 November 2026:

1. **Named checks.** `fid check <name>... --hold <id> --actor <a>
   --idempotency-key <k>` requires names in this slice; no names gives
   `invalid-argument`, asking to name the checks. 1.5a6 supplies the lane's
   required set when names are omitted. Unknown names give `invalid-argument`
   with the defined names in `allowed`, before execution. Names are deduplicated
   and executed in name order, including names the lane does not require.
2. **Preconditions and policy.** The exact hold must be current for this actor
   and worktree, and its task must be `in_progress` with a lane. The lane must
   have a `checks` key; missing gives `checks-not-configured` (precondition,
   exit 4, subject the task, details `{lane}`), naming `[lanes.<lane>]`, merging
   `.fid.toml` and `git fetch` as the repair. Explicit `checks = []` is set up.
   Policy must be readable (`policy-unavailable`), and a changed effective
   digest writes a snapshot under the write lock.
3. **Output (Q6 A).** On a terminal, the child's output appears as it arrives,
   with controls escaped. JSON captures it and emits one final `check.execute`
   envelope, with `{task_id, checks}`; each summary has `id`, `name`, `outcome`,
   nullable `exit_code`, `counts`, `reasons` and `artifact_ids`. No streaming
   JSON protocol is added. Human output announces each check and its outcome,
   says whether it counts for hand-off, and lists each `Kept <path>`.
4. **Several checks (Q21 A).** An ordinary unsuccessful check does not stop the
   rest. A takeover, lost hold or fid error stops execution, leaving later names
   recorded as not started. All passing and counting gives exit 0. Unsuccessful
   checks give `check-unsuccessful`, the existing exit-8 category, with
   `{check_ids, outcomes}` for those checks. If all passed but some do not count,
   `check-not-qualified` is precondition / exit 4, with `{check_ids, reasons}`.
   Failure envelopes have no `result`; their IDs lead to the records.
5. **Start and result (Q17).** The start lands under the write lock before
   spawning; the lock is released while the child executes. It records ID,
   task, exact hold, actor, name, resolved argv, timeout, inherited environment
   names, listed Artifacts, the definition digest, policy snapshot, time and
   HEAD/tree/dirty observation. Ignored files do not make it dirty. Process ID
   and start identity are recorded as soon as known. The result records finish
   time; `passed`, `failed`, `timed-out`, `cancelled`, `spawn-failed` or `signaled`;
   exit code and signal; the second tree observation; bounded output digest,
   16 KiB excerpt, total byte count and truncation flag (10 MiB capture); the
   resolved executable and its file digest, OS, architecture and fid version;
   and kept/refused files with their paths and refusal reasons.
6. **Qualification.** A result counts only if it passed, started with nothing
   uncommitted, the tree did not change, ownership was still current at the end,
   and every listed file was kept. The reasons are `tree-dirty`, `tree-changed`,
   `ownership-lost` and `artifact-not-kept`. The definition digest covers argv,
   timeout, inherited environment names and listed Artifacts, excluding
   `defined_by`, `when_changed` and unrelated policy settings.
7. **Environment (Q7 + Q23 A).** Inherit only set `CI`, `GOCACHE`, `GOMODCACHE`,
   `GOPATH`, `HOME`, `LANG`, `LC_ALL`, `PATH`, `TMPDIR`, `XDG_CACHE_HOME`. Record
   a digest of their resolved values, never the values. Resolve the executable
   with that PATH. Execution uses argv without an implicit shell, at the
   worktree root. There are no `env` or `secrets` keys in 1.5a.
8. **Kept files (Q8 A).** Each check may list `artifacts = ["coverage.out"]`.
   These are plain worktree-relative paths, not globs. The reader, effective
   policy/digest, snapshots and replay validation carry them. After successful
   and unsuccessful checks, retain each safe regular file in the existing
   content-addressed Artifact store: inside the worktree, no symlink escaping
   it, unchanged while copied, at most 64 MB per file, with no configurable
   limit. Missing/refused files give `artifact-not-kept`. Export, import and V8
   preserve both records and bytes. The ceiling keeps these Artifacts below
   GitHub's 100 MB per-file backup limit.
9. **Timeout and renewal (Q9 A).** Each check uses its own positive timeout,
   without the old runner's ten-minute cap. Timeout stops the whole process
   group and records `timed-out`. While the child executes, renew its exact
   hold every five minutes, one Event each, with heartbeat semantics. Refused
   renewal stops the child and later checks, records `ownership-lost`, and
   returns `hold-not-current`.
10. **Unfinished earlier checks (Q22 A).** Before spawning, inspect every start
    on this task without a result, using same-machine PID plus start identity.
    A stopped process permits a fresh check, warning `⚠ CHK-100001 never
    finished; its process has stopped`; the old start remains unfinished. A
    live process gives `check-still-running` (conflict / exit 4), details
    `{check_id, pid}`, suggesting waiting or stopping that PID and checking
    again. A lapse does not prove the process stopped.
11. **Retries and show.** A same-key repeat never spawns and returns its recorded
    response. An interrupted keyed request gives `check-incomplete` (conflict /
    exit 4), details `{check_ids}`, without inventing a result. `fid show` adds
    Checks: each name's latest result, plus unfinished starts; JSON includes
    the full selected records. Every new Event family joins both catalogs,
    replay/import validation and V8 in this PR. No kernel record is created or
    needed. Required-set evaluation, `when_changed`, hand-off, CI results and
    replacing the old Artifact/export fixtures stay with their later tasks.

Builder's storage choices for 1.5a5: projection version 7 forces a journal
rebuild of older derived caches, and migration 25 adds `check_requests`,
`check_starts`, `check_processes`, `check_results` and `check_files`. Check IDs
use `CHK-` and the existing unused-ID generator. The first start reserves its
retry key and records the ordered requested names atomically; later starts and
results advance that same retry state. Each start names the snapshot matching
the effective policy resolved for that invocation. The final result completes the response
and records the names not started on an ownership stop or cancellation. A
refusal between completed checks uses `check.stopped` to retain the original
response without another start or result. Starts/results are immutable, and no terminal observation is manufactured after a fid error.
The start also records the owning fid process. An internal pipe handshake
keeps the child from executing the check until its PID/start identity is
journaled. If fid dies before granting permission, pipe EOF ends the launcher
without executing the check, so an unrecorded child identity cannot hide
continuing work. The private descriptors close before exec, which preserves
the recorded PID and start identity. Machine identity uses `/etc/machine-id`
on Linux, falling back to the hostname (with distinct `machine-id:` and
`hostname:` prefixes). Linux records its boot ID separately from machine
identity and records process start ticks; other Unix systems use the hostname
and process start time, with an empty boot ID. An unfinished start on this
machine from a different boot has stopped and permits a fresh check with the
usual warning; a different or uninspectable machine still refuses. The retry
lookup and first-start admission share one lock acquisition. A refusal between
checks is recorded under the lock that observes it; the result and its retry
response are committed together. An
inspection error is a fid error, never evidence that a process stopped.
`check.started`, `check.process-started`, `check.result-recorded`,
`check.hold-renewed`, `check.artifact-kept`, `check.artifact-refused` and
`check.stopped` are the new families. The shared file publication/verification helpers move out of the
old executor into `artifactstore`; lean checks do not import the old executor.
An omitted or empty Artifact list has the same meaning; its optional effective
encoding preserves snapshots recorded before this schema addition.

Mark approved these details for **1.5a6**, on 14 November 2026:

1. **No names means the lane's required checks.** `fid check --hold <id>
   --actor <a> --idempotency-key <k>` selects the task lane's `checks`, minus
   conditional checks whose `when_changed` paths the change does not touch.
   Execution, output and exit rules stay those of 1.5a5, in name order. A
   missing lane `checks` key gives `checks-not-configured`; `checks = []`
   records a successful empty request, with `✓ No checks required for
   Standard` (for that lane) and exit 0.
2. **Changed paths (Q2 A).** Compare committed files between HEAD and its
   merge-base with the local `origin/<default>` ref, the same default used
   by policy. Never fetch. Recompute the base at every check and hand-off,
   so after merging the default branch only the task's own files count.
   Include deletions and both names of renames. Missing or empty
   `when_changed` always requires the check. An unavailable base (missing
   origin ref, insufficient shallow history or no common ancestor) refuses
   as `comparison-unavailable`, precondition / exit 4, subject the task,
   details `{ref, problem}` with nullable `ref`. Explicit names do not need
   a comparison. Uncommitted files already prevent a result counting.
3. **Latest finished result (Q5 A).** For the same task, tree and
   check-definition digest, choose the latest finished result in journal
   order, never by caller times. A later failure requires another pass;
   an unfinished newer start preserves the earlier finished verdict.
   Another task's pass on the same tree does not count. One shared
   evaluator serves hand-off (1.5a7) and done (1.6), returning required
   names, chosen result IDs and missing names. Its `checks-unsatisfied`
   refusal is precondition / exit 4, subject the task, with
   `{checks, latest_result_ids}`; hand-off will use it.
4. **Missing test paths do not stop checks.** Missing `[review]
   test_adequacy` does not affect check; hand-off refuses it in 1.5a7.
5. **Show the current requirements and satisfaction.** The Checks row
   marks which checks the lane currently requires and whether each has a
   counting pass on HEAD's tree. JSON adds `required` and `satisfied` per
   check, plus the base commit used. Keep this display to those additions.

The orchestrator (**Claude Opus 5.5**, acting on Mark's behalf under his
delegation for the 1.5a series) decided these clarifications for **1.5a6**
on 14 November 2026:

1. **Show's check entries.** Each JSON `checks` entry is
   `{name, start, process, result, required, satisfied}`, in name order.
   A required check without a task record has null start, process and
   result. Text marks required checks and whether each counts on HEAD's
   tree; an unstarted required check reads, for example,
   `○ test required · not checked yet`. Keep labels plain, without kernel
   words.
2. **Show keeps working (Q16 A).** The Checks row never refuses the read.
   Unavailable policy makes required and satisfied null for every check.
   Missing lane checks makes required null. An unavailable base makes
   required null only for conditional lane checks; unconditional ones
   remain required. Unknown never means not required. Satisfied is null
   whenever the current definition is unknown, including a removed name.
   `checks_base` is the comparison's base commit, null if unavailable or
   unnecessary. JSON adds `warnings`, empty when healthy, using status's
   warning shape: kind, summary, next_action and the kind's details.
   Kinds are policy-unavailable, checks-not-configured and
   comparison-unavailable; comparison details are nullable ref and problem.
   Text adds one warning line under Checks explaining why fid cannot tell.
   Reuse status's warning code where simple.
3. **Allocator correction.** Empty-request IDs must be checked against
   both existing starts and requests, with a regression test.
4. **Empty completion accepted.** The `check.completed` family for a
   successful empty request is accepted, with both catalogs and V8 coverage.
5. **Unsupported changed paths require all conditional checks (review round
   1).** A committed Git path that `pathglob` cannot evaluate, including
   control characters, backslashes or invalid UTF-8, counts as matching every
   conditional check. It must never become a nonmatch or a comparison refusal
   that blocks checking the task while that file exists. `fid show` uses the
   same rule when marking current requirements. This choice was decided by
   the orchestrator on Mark's behalf under the 1.5a delegation.

Builder's transactional choices for 1.5a6: `policy.Compare` uses local git
plumbing with lazy fetching disabled; `git diff --no-renames --name-only -z`
includes both rename paths without git's quoting. `checkread.Read` returns
finished results in result insertion order (preserved by journal replay),
then unfinished starts. `checkread.Evaluate` takes the required set and its
definitions explicitly, so done can supply the hand-off policy. Execution
and evaluation share `checkrecord.DefinitionFor` for the definition digest.
An empty selection uses `check.completed`, a completed `check_requests` row
and its retry response, with no invented start or result. Both Event
catalogs, replay/import validation and V8 cover it; it needs no migration.
The `checks-unsatisfied` code has no CLI producer before hand-off in 1.5a7;
the process-pair inventory records that audited exclusion until then.
Show retains the full task check history until evaluating the current tree
and definition, falling back to the latest historical result for display.
It shares status's warning constructors and JSON shape. Named and empty
requests use one allocator; a deterministic collision test includes a
secondary start whose ID is absent from the request table.

Mark approved these details for **1.5a7**, on 14 November 2026:

1. **Hand-off (Q3 A).** `fid handoff --hold <id> --actor <a> --pr <url>
   --idempotency-key <k>` requires the key and a GitHub pull request URL
   (`https://github.com/<owner>/<repo>/pull/<n>`). It uses recorded check
   results and makes no network call; 1.6 checks the PR with `gh`. Under
   one write lock it records the PR, head commit, tree, comparison base,
   lane, policy snapshot, required check names, their definition digests,
   counting result IDs and time, moves the task to `in_review`, and ends
   the builder's hold in one Event.
2. **Refusals change nothing.** Missing lane checks uses
   `checks-not-configured`; missing `[review] test_adequacy` uses
   `test-paths-not-configured` (precondition, exit 4, details `{}`), with
   `add [review] test_adequacy, then merge .fid.toml and git fetch` as the
   fix. Uncommitted changes use `tree-dirty` (precondition, exit 4,
   `{worktree}`), asking the builder to commit before hand-off. A required
   check without a latest counting pass on HEAD's tree uses
   `checks-unsatisfied`, suggesting `fid check --hold <id> --actor <a>
   --idempotency-key <key>`. An unavailable comparison uses
   `comparison-unavailable`; an invalid PR uses `invalid-argument`.
   Exact hold checks and `policy-unavailable` apply as to other work
   commands; a task not `in_progress` uses `transition-not-permitted`.
3. **Keep the worktree through review (Q4 A).** The task stays kept for
   its builder worktree with `kept.mode: "review"` until done (1.6) or
   cancelled. That worktree's `next` returns kind `none`, null hold, task,
   epic and answer, and `waiting: {kind: "review", task_id: <id>}`, saying
   `Nothing here: <id> is in review` and suggesting `fid show <id>`.
   The review reservation takes precedence over an older unanswered
   question; `fid show` still lists that question.
   The queue retains its usual contents. Other worktrees carry on. The
   loop opens a worktree per task (2.5); 2.1 returns rework here as kind
   `rework`. Cancellation frees it; missing-worktree NEEDS YOU applies.
4. **Show and status.** Show adds a `Hand-off` row such as
   `10:10 · PR #63 · lint, test passed on a1b2c3d`, and JSON adds the latest
   `handoff: {id, pr, head, base, required_checks, result_ids, at}` or null.
   Status labels the task `in review`. Empty requirements read
   `no checks required` in show and the hand-off success message.
5. **Rework and retries.** A later hand-off after rework may name another
   PR; done uses the latest. A same-key retry returns the original
   response and Event even after its hold has ended.
6. **The lean journey.** The integration test goes from issue creation
   with a lane through next, start, heartbeat, check without names,
   commit, hand-off and next waiting for review, with interrupted and
   waiting/answer detours. Its journal contains no kernel record. Both
   catalogs, replay/import validation, a projection migration and V8
   cover the new hand-off family. A rework-shaped test belongs to 2.1.

Builder's transactional choices for 1.5a7: migration 26 adds immutable
`handoffs` receipts with `HOF-` IDs; projection version 8 rebuilds older
derived caches. The migration extends a kept hold's end reason
with `review`. The receipt keeps committed changed paths alongside the
base so replay and import can re-evaluate conditional requirements without
git or the network. `task.handed-off` carries the receipt, task transition,
hold reservation and retry response atomically. Snapshot changes are
recorded only after all hand-off preconditions pass, under that same lock.
The latest receipt is selected in journal insertion order, never by caller
timestamps. A review reservation in another worktree uses the existing
`held-elsewhere` scheduling reason. Rebuilding the referenced holds table
defers foreign keys, verifies the entire final graph with
`pragma_foreign_key_check`, then clears SQLite's stale deferred DROP
counters before commit. Hand-off observes HEAD and cleanliness without
hashing uncommitted contents; this includes unmerged paths in `tree-dirty`.
Its git children disable lazy fetching and optional index writes, and it
retains the guard against hidden tracked changes. The downgrade retains review reservations as
interrupted reservations with the reason `In review`.

**Settled by the orchestrator, 14 November 2026 (Mark's delegation for
1.5b-1.5d).**

Mark delegated the remaining decisions for the kernel's removal to Claude
Opus. Seven stacked Assured PRs run in order: 1.5b1 extracts the Artifact
verifier and rebuilds V8 and export fixtures; 1.5b2 rebuilds the other
surviving fixtures; 1.5c1 moves fixed task states out of workflow; 1.5c2
removes commands and switches `fid artifact`; 1.5c3 removes kernel domain
packages; 1.5d1 closes journal and archive catalogs; 1.5d2 drops tables and
kernel maintenance checks. `RESET.md` records each completion criterion.

- **Q1:** The current `review` and `finding` commands use the kernel and
  leave with its commands in 1.5c2. No interim stub; 2.1 brings lean ones.
- **Q2:** A journal with a removed event family is refused with
  `journal-replay-stopped`, with details naming `schema-unsupported` and
  the family. An archive uses `archive-invalid` and its existing
  `schema-unsupported` check. Either refusal changes nothing.
- **Q3:** Frozen V1 `agent_runs`, `triage_events` and `triage_status` stay
  for separately scoped cleanup; see the later list.
- **Q4:** Maintenance drops `claim_count` and kernel-only warnings. The
  dead `maintenance.abandoned_attempt_grace_seconds` and `executor_limits`
  config keys stay for now; see the later list.
- **Q5:** Archives keep manifest v2, `sources: []`, and only truthful
  capability profile entries.
- **Q6:** `fid artifact verify` keeps its input, all nine Artifact success
  fields, Event and missing-ID failure. Its actor becomes typed: malformed
  actors get `invalid-argument`, exit 2; agents and people may verify;
  services get the existing actor-type refusal. Kernel readiness Stops go.
- **Q7:** Check results stay immutable. Whether hand-off or `done` needs a
  kept file still available later is 1.6's decision; never rewrite a result.

- **1.5b2 Artifact command-test exception**, settled by the orchestrator
  (Claude Opus 5.5, under Mark's delegation), 2026-11-14: `fid artifact`'s
  own command tests that need its current Activation-based adapter stay
  unchanged as an isolated compatibility group until 1.5c2 switches the
  adapter and converts them. No other surviving test may depend on their
  kernel fixture. All other surviving tests become kernel-free in 1.5b2;
  none are weakened. This keeps 1.5b2 test-only and the command's public
  actor-authority change with 1.5c2.

Builder's choices for **1.5b1**: `internal/artifactverify` owns shared
metadata, blob verification and retry/event handling. The command still uses
its existing kernel adapter; compatibility hooks keep its authority and
readiness effects until 1.5c2 switches it. Lean callers use typed actor
rules and emit only Artifact availability and retry Effects. V8 reads `next`
at a fixed clock; its new-machine wipe recreates a fresh repository with
only committed identity at the same path, keeping worktree-dependent reads
comparable without changing historical IDs. Lifecycle counts are unchanged.
Generic export fixtures use real lean checks; kernel Source and recovery
assertions keep explicitly private fixtures until their profiles leave.

Builder's transactional choices for **1.5b2**: generic journal retry and
sequence tests use `label.added`; strict retry-map and prior-state tests use
`task.started` under a lean hold. Issue list and detail samples are
journal-backed retained rows, with Concern assertions separated into
kernel-only tests. The historical projection and trigger assertions stay
in explicitly historical tests until their owning schema cut; new retained
projection samples keep full-table parity, exact row values and foreign-key
checks. The named-read/export journey now starts, checks and hands off a
fixed-state task and compares the retained reads and tables. Integration
helpers move to neutral command, projection and keyed-request files. The
retained single-journal corpus exercises all 15 non-Artifact keyed leaves,
including their exact retry responses and field maps; the old 55-leaf
kernel corpus stays separately until 1.5c2. Shared failure codes now use
retained producers without adding registry exclusions. Closure reads pin
the person's cancellation actor, reason, task and binding. Neutral CLI git
and JSON writers move out of their retiring host files as well. The Artifact
exception above covers the two own-command tests in
`cmd/fid/artifact_compatibility_test.go` and
`TestP24ArtifactVerifyKeyedLeafCorpus`; their bodies are unchanged.

Builder's transactional choices for **1.5c1**: `internal/issuestate` owns
the retained task, epic and subtask state definitions, projection registry,
`NextStates`, initial-state and creation lookup, and child-state rule types
and evaluation. Staying Issue, hold and maintenance callers and seed-parity
tests import it directly; next, status and show reach it through Issues.
`internal/workflow/issue_states.go` keeps aliases and forwarding functions
for kernel APIs until 1.5c3 removes that package. Only the kernel execution
transition and linked-decision rule paths in Issues still import workflow.
Creation retains subtask kind/project assignment precedence in this
behaviour-preserving slice; configurable resolution leaves with authoring
in 1.5c2. The graph, foreign keys, workflow IDs, category JSON and error text
are unchanged. Existing state, registry, initial-state and child-rule tests
move with their implementation, keeping their assertions.

Builder's transactional choices for **1.5c2** (2026-11-15): the remaining
command tree has 37 executable leaves (35 with the Response Envelope and
two diagnostics; Cobra help/completion are additional), and 16 keyed leaves in 10 families, excluding export
import under the existing inventory convention. Generic failure writing,
JSON and fixture helpers move out of retiring adapters. Artifact
verification calls `artifactverify.Verify` directly; its temporary Hooks and
Evidence readiness wrapper are removed. Typed actors retain the established
`authority-denied` details; malformed actors fail in flag validation. The
missing-Artifact failure keeps its existing `invalid-argument` code and
legacy `accepted gate registration input` allowed-text to preserve Q6's
failure contract. No check result is rewritten by an availability change.
The descriptor's seeded mismatch test copies its field slice before changing
it, so the live registry remains immutable during parallel tests.

Issue creation always selects its kind's fixed built-in workflow, including
subtasks, and ignores retired assignments. Concern fields/filtering,
linked-decision rules and execution-chain transition guards leave staying
callers. Domain packages and their unit tests remain for 1.5c3; the
kernel-only policy grant helper still serves those domains and leaves with
them. Journal/archive catalogs, tables and maintenance checks remain for
1.5d. The retained conformance fixtures use the exact remaining operation
inventory and live failure registry; all 16 keyed leaves execute and rebuild
in one lean journal. Historical Issue JSON captures are regenerated with the
current CLI and their capture recipe uses the git-common-directory ledger.
The export import/recovery domain test remains until the domain cut, with
only its embedded historical Gate CLI assertion and helper retired here.
The modern status/show/next/hold/start/check/handoff goldens stay unchanged;
only retired workflow-assignment goldens leave.

Builder's transactional choices for **1.5c3** (2026-11-15): delete the
now-unreferenced `claims`, `workflow`, `evidence`, `review`, `stats` and
`kerneltest` packages, their tests, and `policy/envelope_authority.go`.
The fixed-state implementation and tests stay in `issuestate`; the lean
Artifact verifier and store stay in their own packages. Both source import
searches and `go list -deps -test` verify that retained code has no incoming
edge to a retiring package before directory deletion.

The export fixture that invokes kernel domain writers leaves with its
Source/recovery/runtime acceptance tests. Handwritten archive recovery
reference and Effect-family tests, and journal, schema and maintenance
table-only tests stay until 1.5d. The mixed import boundary tests keep their
shared-ledger, moved-ledger and removed-worktree binding checks through
`config.ReadImportBinding`; only kernel recovery admission assertions go.
Import replay after later local Events remains covered by the existing
lean import test. No replacement helper, stub, dependency or public
behavior is added. Event catalogs, archive profiles, tables and maintenance
queries remain for 1.5d1/1.5d2.

**RESET 1.5d1 implementation, 2026-11-15 (Assured).** Following Q2 and Q5:

- The journal and archive catalogs accept only the retained lean families.
  One shared classifier rejects all 41 removed families and Effects on their
  removed tables, including fixed-workflow authoring and assignment. Generic
  version-1 Issue history remains supported in the journal. Physical tables
  and migration versions stay for 1.5d2.
- Replay checks the complete journal for unsupported history before applying
  any Event, including when a checkpoint would otherwise skip that history.
  Projection freshness inspection makes the same check. A valid prefix before
  an unsupported Event cannot partially update the projection.
- `journal-replay-stopped` keeps its original three detail fields for ordinary
  corruption and precondition conflicts. Its closed schema-refusal variant
  adds `reason: "schema-unsupported"` and `event_type`, and fills `event_id`
  with the refused Event's ID. The file and line still identify the source.
  Human output describes unsupported history without the retired vocabulary.
  Archives keep `archive-invalid` with a `schema-unsupported` check naming
  the family and the archive journal member. Entry integrity is checked first;
  an old capability/source profile cannot mask a removed family. Export
  creation also refuses unsupported source history before publication.
- Manifest v1 is refused by validation as well as import. Manifest v2 keeps
  `sources: []`; a null or nonempty inventory is refused. Export and archive
  validation no longer query `envelope_sources`. The only surviving entry of
  the former capability profile is the truthful `export.portable.v1`; no
  replacement capability IDs are introduced. Private semantic validators keep
  migration boundaries 19–26 for lean integrity tests; public archive admission
  still requires the active migration version.
- The retained Artifact availability Event requires Event-v2 and admits only
  its guarded availability update and one idempotency result, with strict
  actor, command, prior, metadata
  and cached-result bindings. Recorded check results remain immutable (Q7).
  Kernel-only catalog acceptance tests leave with their subjects; the generic
  malformed-Effect probes now use a fixed-state task under a lean hold. Two
  transition-refusal fixtures corrupt the retained seeded state directly in
  SQLite, rather than append unsupported workflow-authoring Events; their
  assertions about missing edges and unknown requirements are preserved.


**RESET 1.5d2 implementation, 2026-11-15 (Assured).** Following Q3/Q4 and
ADR-0014, migration 27 drops 42 retired tables in foreign-key dependency
order. It preserves the shared Artifact/retry tables, all lean records,
the four fixed-state lookup tables and Issue foreign keys, and frozen V1
`agent_runs` / `triage_events`. `issues_enriched` keeps its latest-triage
lookup and removes Concern columns/joins. The activity insert trigger
accepts only existing Issue references. Migration/export schema is 27;
projection version 9 rebuilds older derived caches. Historical migration
pairs remain until 3.2.

Builder's transactional choices for 1.5d2: the down migration restores the
schema-26 DDL, including table-owned indexes and triggers and the old view
and activity guard, without recovering removed rows. Latest-schema tests
pin the exact retained table inventory; round-trip tests compare the full
schema and retained data. Kernel-only schema/maintenance acceptance tests
leave with their subjects; retained triage, structural integrity, retry/
Artifact checks and no-mutation assertions stay. The generic finding-limit
fixture uses missing Artifacts instead of retired warnings. Generic warning
promotion remains covered using an injected warning, since no current
warning producer remains. No hold-warning semantics are introduced.
Maintenance removes its kernel queries, grace option, warnings and
`claim_count_live`; its Info contains only Issue and recent-activity counts.
The generic warning array, promotion flag and collection limits stay.
The dead configuration keys and frozen-state checks stay as Q3/Q4 require.
V8 and the link/rebuild journey now assert retired-table absence instead
of empty retired tables; the scale snapshot compares the retained view
columns. Tests dedicated to ignoring now-removed workflow assignments
retire with that table. The maintenance format-refusal test keeps its
historical identifier registered by the unchanged format audit.
This completes the physical removal for 1.5b-1.5d; old-family refusals still
happen before any SQL projection or archive publication.

**RESET 1.5e, 2026-11-15 (Routine).** Mark added a separate task after the
orchestrator found about 8,000 leaked test directories (about 1.9 GB) over
four days. Disposable `-- --check-child <mode>` processes now enter the
selected helper test before any command-suite setup. The existing check
launcher and `FID_RUN_AS_FID` shortcuts remain. The retained integration
suite has no crash-helper re-execution path.

Builder's transactional choices for 1.5e: both suites run the documented
`go telemetry off` command only after isolating `XDG_CONFIG_HOME` and
verifying that `os.UserConfigDir` points there, before any other Go command.
Suite setup refuses to change telemetry otherwise. `GOTELEMETRY` is a
read-only `go env` value, not an
environment switch. In the installed Go 1.26.2 source, `MaybeParent` starts
a detached telemetry sidecar even in local mode; the sidecar calls
`counter.Open` itself, which can recreate an already removed config
directory. A focused passing unit run reproduced that recreation. Turning
telemetry off prevents both collection and the sidecar, without changing
the developer's telemetry setting. Shared Go build caches stay pinned.

Child coverage uses `.git/fid-check-child-coverage` in the parent-owned,
disposable git workspace, rather than a shared directory under `os.TempDir`.
The existing helper modes and their output are unchanged. Re-execution
regressions count every entry left in private temp directories after normal
suite cleanup, helper exits, signals and timeout kills; separate assertions
check Go's actual telemetry mode. Coverage also exercises the fallback with
`GOCOVERDIR` stripped. Reverting each fix makes its regression fail. The
integration regression includes its normal binary build; unchanged focused
and full integration runs left no entries before the fixes.

Round 1 adds a test-only config resolver seam to both telemetry guards.
Each suite tests a resolution error (including one returned with a matching
path), a location outside the fixture, and successful isolated setup. The
negative cases assert refusal and preserve a seeded simulated user's mode
file byte for byte; removing either suite's guard fails both assertions.

The historical integration source at `982b4211` exposes the paired-directory
leak: `TestP22ActivationCrashProcess`, `TestP22BudgetCrashProcess` and
`TestP22RecoveryCrashProcess` re-executed the integration test binary,
inheriting `FID_INTEGRATION_BINARY`. Each child allocated new config/cache
directories, reused that binary, then called `os.Exit(23)` inside its test,
so `runIntegrationSuite`'s deferred removals never ran. These helpers were
removed in 1.5c2 (`6124adcc`), before 1.5e's starting commit. No dead-helper
workaround is added. This source evidence identifies a historical leak
mechanism; the old machine's individual 761 pairs were not inspected in
this session. No production code or dependency changes.

**RESET 1.5f, 2026-11-15 (Routine).** Mark added this separate task
following a local coverage/no-identity CI run that crossed UTC midnight and
failed 47 command tests. The shared pristine workspace was initialized on
one UTC date and later commands correctly opened another dated journal file.

Builder's transactional choices for 1.5f: test readers concatenate all
journal files in lexical (date) order and re-enumerate on every read,
including before/after Event counts and byte comparisons. Event lookup and
tail assertions take a workspace instead of a cached file path. Tests that
damage the tail, or assert inspect's file-local line number, explicitly
select the last file. Single-file assertions remain only on a fresh init's
bootstrap journal, before any further write; the number of Events and all
existing expected mutations remain unchanged. Integration's closure journey
reads the whole journal, checks every file is untracked, and searches every
file for Event attribution. Its existing full file snapshots remain.
Regressions backdate a copied genesis Event to yesterday, rebuild its
projection, and run real commands: one starts with a journal spanning dates,
one opens the new daily file after the before snapshot, and integration
checks the closure and attribution helpers on both dates. No clock override,
production change or dependency is added.

## RESET 1.7 implementation (2026-11-15)

Mark asked on 2026-11-15 for 1.6 and 1.7 to be orchestrated while he was
out. The orchestrator chose to build 1.7 alongside 1.6, from `origin/main`
at `4517e22d`, so it may merge before 1.6. This is a small deviation from
table order; 1.7 does not use or change 1.6's code. 1.7's brief, written
by the orchestrator, permits plain human-text repairs and raises the lane
to Standard if any are needed.
V2 found three help leaks: backup's "attempt", next's "run again" and SQL's
"Run". Their wording is repaired. The lane is **Assured**: the PR adds a
Makefile target (`check-leak-test`), and RESET.md puts Makefile changes in
Assured. The brief had called the task Routine, overlooking that rule; the
review caught it and the orchestrator raised the lane.

Builder's transactional choices:

- `internal/leaktest` owns the unchanged kernel-word pattern and assertion.
  Only `_test.go` files import it. The assertion's small testing interface
  allows its rejection diagnostics to be tested without a child process.
- `TestV2HumanOutputNamesNoKernelWords` recursively scans every `.txt` under
  `cmd/fid/testdata`, including files with no current golden-test caller.
  It walks the full Cobra tree, including generated help and completion,
  checking Use, Short, Long, Example, local/inherited flag usages and actual
  `--help` output. Only help masks the exact `--claim` token; surrounding
  prose and other flag spellings remain subject to the unchanged pattern.
- The live journey initializes a ledger and creates, takes, starts, checks,
  interrupts, resumes, parks waiting, answers, takes over and hands off a
  task, reading status and show along the way. Every action uses text mode;
  JSON show reads supply only the hold IDs that human text omits. It uses
  the existing check child and its coverage-directory guard. A separate
  V2 refusal matrix checks both streams, as does every journey action.
  `done` is absent at this base and remains 1.6's work.
- A blanket assertion was considered for `executeCommandForTest` and left
  out: it has no test handle, and its callers deliberately supply arbitrary
  echoed data and diagnostic JSON. Broad input masking could hide fid's own
  text. The scoped journey and existing golden/refusal assertions need no
  such exemptions. `--json` and verbose diagnostics remain outside V2.
- `make check-leak-test` selects both V2 tests with a fresh execution;
  ordinary `make test` and CI discover them too. Temporary help, orphan
  golden and refusal leaks each failed with the shared diagnostic, then
  were restored.

Amended 2026-11-23 (Mark): the shared pattern no longer refuses "claim",
"run" or "attempt" in any form. The kernel's code is gone, and those words
are now ordinary fid English: `fid next --claim` is a real flag, and `fid
check` runs checks. Avoiding them had bent wording twice (the review input's
`purpose` field in 1.6b, and 2.1d's status text). The pattern keeps the
words that still mean nothing in fid (envelope, activation, obligation,
attention item) and the `ENV-`, `STP-` and `ATT-` prefixes, so V2 holds.

## Decided 2026-11-15 (orchestrator, under Mark's delegation)

Mark delegated 1.6's design to the orchestrator on 2026-11-15: "You can
settle the 1.6 design decisions; just record these to report back to me
later. If the survey finds work scheduled for later is actually required
prior to 1.6, then I am ok with bringing that forward. As stated
previously, I do want this done right." Two read-only surveys (1.6 itself,
then 2.1 and 2.3) found that `done` must check reviews, findings and Mark's
fingerprint-bound approval, whose records were due in 2.1 and 2.3, after
1.6. No recorded rule covered `done` before those exist, and 1.5c2 removed
the old review commands with no interim stub. So the recording parts of
2.1 and 2.3 move into 1.6's sequence, and `done` checks real records only:
no GitHub stand-in for reviews or approvals, and no flags a caller supplies.
`RESET.md` holds the twelve rows, 1.6a–1.6l, all Assured.

**Which PR, which merge, which CI**

1. **The PR is the latest hand-off's.** `done` uses the PR URL in the
   task's latest hand-off, never a branch-name guess. The PR's base
   repository must be the project's `origin` repository (fork heads are
   fine), its base branch the default branch `origin/HEAD` names (as policy
   discovery uses), and it must be merged: state merged with a merge
   commit, never inferred from closed.
2. **Checks bind to the merged head's tree.** The merged PR head's tree
   must equal the hand-off's recorded tree. Otherwise `done` refuses and
   the task is handed off again, which re-evaluates its checks on the new
   tree. A clean rebase that keeps the tree keeps the checks.
3. **CI green, GitHub Actions only.** Every Actions workflow run on the
   exact merge commit with event `push`, taking the latest attempt per
   workflow, must have completed with conclusion `success`. Skipped runs
   are ignored. No runs, only skipped runs, or any queued or in-progress
   run means "waiting for CI". Any other conclusion means CI isn't green.
   The run's conclusion counts, not its jobs' (fid's own `main` runs skip
   test jobs on verified tree reuse and still succeed). An incomplete or
   inconsistent response is never green. External status contexts and a
   setting to choose workflows are on the later list; there is still no
   key to opt out.
4. **What becomes a record.** A confirmed merge is recorded against its
   task (the merge commit), even while CI is pending, so offline reads can
   say "merged, waiting for CI" or, when GitHub can't be reached,
   "merged, waiting for confirmation". Each completed CI run on the merge
   commit, success or failure, is recorded as a check result from CI: run
   ID, attempt, workflow, commit, conclusion and when it was observed. It
   has its own typed storage and never invents a local check start,
   process or hold. In-progress runs and network or authentication
   failures aren't recorded.

**1.6b's choices (2026-11-15).** `merge.observed` (`merge.observe`) and
`ci.result-recorded` (`ci.record`) are new lean Event-v2 families on a
`merge` aggregate, separate from local `check.*`. Migration 28 adds
`merge_observations` and `ci_results`; projection cache version 10 makes
an existing caught-up projection rebuild with their tables. Indexed
binding columns and immutable typed JSON keep the full observations.
The merge key is (hand-off, merge commit); the CI key is (normalized
repository, run ID, attempt, merge observation). Identical retries return
the original records and times; changed facts under an existing key refuse.
A CI batch appends all new completed attempts as one Event or changes
nothing. Pending attempts append nothing. These natural keys need no
command retry-result row; both families refuse injected retry effects.
Writers take already-read 1.6a facts and do no networking; the ledger lock
covers refresh, validation and append. The caller supplies only CI facts
on this merge commit; a differing head is refused. All completed events
are retained, including skipped, superseded and non-push observations;
the later evaluator selects which count. Latest-merge and CI history
readers use journal insertion order. CI Event decoding preserves integer
identities without a floating-point conversion. `fid done` (1.6i) is the first public
caller. Local check storage, readers and evaluation are unchanged.

**The `done` command**

5. **`fid done <task> --actor <a>`.** Any work actor may run it; the
   hand-off ended the build hold, so no live hold is needed. Each call
   makes one GitHub observation per task: no polling and no background
   work. A task waiting for its merge or CI stays `in_review`, and the
   refusal says why with its own retryable code (not merged, waiting for
   CI, CI not green, GitHub unavailable, or the failing gate). A pending
   result is never cached under a retry key. `status`, `next` and `show`
   never use the network.
6. **`fid done` with no task reconciles.** It runs the same evaluation and
   recording for every task in review, which is how PRs Mark merged in
   GitHub get completed. It reports each task's outcome and exits 0 unless
   the command itself fails.
7. **Gates come from the hand-off.** Every gate is evaluated at the latest
   hand-off's lane and policy snapshot: the required checks (as 1.5a
   settled), the review profiles, the test-adequacy trigger and the
   Assured gates. If the task's lane was raised after its hand-off, `done`
   refuses until it is handed off again. Later policy changes don't change
   a hand-off's requirements retroactively; who may run `done` comes from
   current trusted policy.
8. **Kept Artifacts must still be there** (1.5b–1.5d's Q7). Hand-off and
   `done` require each counting result's kept Artifacts to exist with
   matching digests. Results stay immutable, and replay never re-reads
   files.
9. **An open question when the task completes** is resolved with its own
   resolution, the task finishing before an answer: neither "cancelled"
   nor a made-up answer.
10. **`## Rollback`** means a level-two `Rollback` heading in the PR body,
    outside fenced code. `done` checks that it is there, not what it says.

**Epics**

11. Cancelled tasks count as finished. `fid done <epic>` works from `todo`
    or `in_progress` when no task is open (an epic with no tasks
    included); a task in review is open. Only a person may close an epic,
    with a reason. NEEDS YOU shows an open epic with at least one task and
    none open, so a newly created, empty epic doesn't nag.

**The change fingerprint**

12. **Computed locally, from the hand-off.** Reviews and approval bind to
    the fingerprint of the task's latest hand-off: its recorded base and
    head, from objects already in the repository, never fetched (lazy
    fetch is disabled). `done` compares against the same computation, which
    is sound because item 2 ties the merged head's tree to the hand-off's.
    `gh pr diff` isn't used: it fails for PRs over 300 files (HTTP 406).
13. **Which patch-id.** D10's `git patch-id`, as `--verbatim`, per file of a
    diff rendered with pinned settings (no external diff or textconv; fixed
    context, algorithm, prefixes and rename detection; full index; binary
    changes included; independent of user config). Whitespace edits change
    the fingerprint, as D10's "any edit" says; the line shifts of a clean
    rebase don't. The algorithm is recorded with each fingerprint. Empty,
    partial, oversized or unavailable input refuses; there is no empty
    fingerprint. Review round 1 of #79 reproduced Git 2.53 skipping the
    next file header after a binary patch, losing that change's path.
    D10's patch-id is therefore applied separately to every file section,
    and the fingerprint combines the sorted per-file IDs with their count
    using SHA-256 (orchestrator, 2026-11-15, under Mark's delegation;
    Mark reviews this later). Review round 2 reproduced an edited rename's
    similarity changing from 97% to 99% on a clean rebase over distant
    upstream additions. Before hashing each section, drop its `similarity
    index` and `dissimilarity index` extended headers: they depend on file
    size and add no change content beyond the paths and hunks. Keep rename
    paths, modes and all other lines (orchestrator, 2026-11-15, under the
    same delegation).
    Limitation: when upstream changes to a renamed file move it across git's
    50% rename-detection threshold, a clean rebase renders the rename as a
    deletion and an addition, so the fingerprint changes and the change needs
    re-approval. This fails safe (stale, never falsely current), like D10's
    allowance for rebases that alter lines next to a hunk; turning rename
    detection off would stale approvals on any upstream edit to a renamed
    file. Narrowing of the clean-rebase guarantee, for Mark's review
    (orchestrator, 2026-11-15, under the same delegation; #79 review round 3).
14. **A hand-off comes first.** A review round or an approval needs the
    task's latest hand-off. A changed delivery needs a new hand-off; a clean
    rebase gives a new hand-off with the same fingerprint, so an approval
    survives it (V6), subject to item 13's limitation.

**1.6c's choices (2026-11-15).** `internal/fingerprint.Compute` accepts the
recorded base and head as full lowercase 40-hex commit IDs and returns
`patch-id-verbatim-per-file-v1`, a SHA-256 fingerprint and both covered
commits. It has no consumer yet: review rounds (1.6d), approval (1.6g) and
done (1.6h) call it in their own slices. No command, record, migration,
Event or dependency is added here.

The algorithm renders one aggregate diff with patch output, no color,
external diff, textconv or relative paths; `a/` and `b/` prefixes; no line
prefix; three context lines and zero inter-hunk context; Myers with the
indent heuristic; 50% rename detection with no rename-search limit; full
index and `--binary`; all submodule changes in short form; standard output
indicators; quoted paths; no blank-context suppression or order file.
`--binary` preserves binary identities and contents; real add/change/delete
tests distinguish two binary contents. Context changes adjacent to a hunk
may still change the ID, as D10 permits.

The complete aggregate is split at lines beginning `diff --git `. Drop each
section's `similarity index NN%` and `dissimilarity index NN%` lines only
within its extended header, before the first `---`/`+++`, `@@` or `GIT binary
patch` line, then pass that section alone to `git patch-id --verbatim`.
These percentages can change on a clean rebase when distant upstream edits
change file size, while rename paths and hunks remain the same. Rename
detection stays enabled; retain `rename from`/`rename to`, mode lines and
everything else. The `patch-id-verbatim-per-file-v1` name includes this
normalization; no fingerprint has been recorded yet.

Limitation: when upstream changes to a renamed file move it across git's
50% rename-detection threshold, a clean rebase renders the rename as a
deletion and an addition, so the fingerprint changes and the change needs
re-approval. This fails safe (stale, never falsely current), like D10's own
allowance for rebases that alter lines next to a hunk; turning rename
detection off would stale approvals on any upstream edit to a renamed
file. Narrowing of the clean-rebase guarantee, for Mark's review
(orchestrator, 2026-11-15, under Mark's delegation; #79 review round 3).
The regression starts with a 60% edited rename and cleanly rebases over
upstream deleting lines 000–049, retaining edits at 060–099 and their context
at 057–059. It asserts the delete/add rendering and changed fingerprint;
the 97%→99% rename regression still asserts an unchanged fingerprint.

Per-file hashing avoids the Git 2.53 parser skipping the next file header
after a binary patch, as reproduced in #79 review round 1; D10's patch-id
is applied per file. Text
content retains its pinned `+`, `-` or space prefix, paths quote newlines,
and base85 binary data contains no spaces, so content cannot forge a file
boundary. Sort the per-file 40-hex IDs lexicographically, retaining every
occurrence. The SHA-256 input is the decimal section count and LF, followed
by each sorted ID and LF (including the last). File order has no effect.
No fingerprint has been recorded using the earlier aggregate recipe, whose
algorithm tag is dropped.

The round-2 extended-header audit checked Git's documented old/new,
new-file/deleted-file modes, copy/rename paths, similarity/dissimilarity
and index lines. Mode and path metadata describe the change itself; copying
and rewrite breaking are not enabled by the pinned renderer. Text `index`
object IDs also change on rebase, but patch-id already ignores them, as the
edited-rename regression verifies with different index headers. Binary
index IDs remain necessary to distinguish binary contents. No further
extended-header normalization was needed.

A disposable private git view borrows only the repository's existing
object store, including alternates, and materializes committed attribute
files from the recorded head using its own index. This excludes repository
diff-driver config (which can change binary classification despite
`--no-ext-diff`), local `info/attributes`, and checkout/index attributes.
It changes no source repository files. System/global config and attributes
are disabled, ambient git overrides are removed, replace refs are disabled,
and every git child sets `GIT_NO_LAZY_FETCH=1` and `GIT_OPTIONAL_LOCKS=0`.
No remote is configured in the private view.

The complete diff must fit within 10 MiB (10,485,760 bytes), and its producer
must exit successfully before any section reaches `patch-id --verbatim`.
Each section must yield exactly one patch-ID line with a zero source-commit
field; every consumer's exit status is checked and its stdout is bounded
to 4 KiB (diagnostics to 8 KiB). Typed failures distinguish missing
base/head objects, empty changes, oversized diffs and git/output failures
(including unavailable trees/blobs). Raw NUL text patches are refused too:
committed attributes can force binary data into text hunks, where even
`--verbatim` silently ignores bytes after NUL. Normal binary patches remain
supported. There is no prefix or empty binding.

**Reviews and findings** (moved from 2.1)

15. **Models are recorded as `<family>/<model>`**, for example
    `codex/gpt-6.1-sol` or `claude/opus-5.5`. "Other model" compares
    families: Codex reviewing Codex is a same-model review, whatever the
    exact models. Each build hold records its model: `next --claim`
    requires `--model` for agent actors, and `fid takeover` takes an
    optional `--model` (none means a person building by hand). A review is
    other-model only if its family differs from every model that built the
    task; an unknown model never counts as other.
16. **Streams.** Each round records its stream, `primary` or, for Assured
    work, `separate`, and an attributed session ID. The separate review must
    come from a different session than the primary; the same reviewer actor
    may do both. The latest round of each stream on the current fingerprint
    counts. A shadow `merge` round satisfies neither.
17. **Required profiles.** `test-adequacy` in every lane when the
    hand-off's changed paths match the snapshot's test-adequacy paths, and
    `code` for Standard and Assured. One round records every profile its
    session ran; other profiles may be recorded but aren't gates.
18. **Majors stay open until closed.** A fingerprint change makes a review
    stale; it doesn't close findings. A major stays open until a later round
    of the same task resolves it or rules it not applicable, or a person
    closes it with a reason. A builder's answer never closes it. An earlier
    passing round never masks a later round asking for changes, and a
    round's verdict must agree with its findings.
19. **Recording someone else's review.** Whoever records the round (its
    `--actor`, such as the orchestrator) may differ from the reviewer the
    input names, with its model, session and stream. Independence is
    checked on the reviewer: no actor that held the task to build it may
    review it, including through an ended hold. The round keeps its full
    text, so 2.1 can render the PR comment from fid's record.
20. **Rework moves forward too.** A round asking for changes moves the task
    from `in_review` back to `in_progress`, and the builder's kept worktree
    gets it back from `next` as kind `rework` (after resume, before new
    work); the fix is checked, handed off again and reviewed again. Review
    scheduling (kind `review`), the rework cap, discovered work, answering
    or deferring findings and posting rounds to the PR stay in 2.1.

**1.6d's choices (2026-11-15).** Model tokens use lowercase ASCII
letters/digits with dots, underscores and hyphens, at most 64 characters
per side of `/`; an empty model records manual or historical unknown
attribution. New agent work requires a model. Old model-less journal holds
remain valid, and a repeated next request preserves its existing hold's
model. Only validated Event-v2 hold inserts may record a model; it is
immutable thereafter. Migration 29 adds `holds.model`, migration 30 adds
`review_rounds` and `round_findings`, and projection cache version 12 makes
caught-up projections rebuild. The event is `review.round-recorded`, on a
`review` aggregate; retired kernel names remain refused.

The closed JSON input carries reviewer actor/model/effort/session, stream,
profiles, verdict, claim, findings (profile, severity, summary, nullable
file/positive line, evidence, fix), verification and optional prompt,
command and raw-output provenance. The CLI spec defines the exact schema
and includes an example adapted from the real 1.6b first-round review.
The field for what the change sets out to do is `purpose` (jig's reviewer
calls it "Claim"); it was renamed when merging 1.7's leak test showed the
help text naming it, since fid's human text never says "claim".
Round IDs are `RND-<six-hex>`; stable finding IDs append their one-based
position. All findings remain open until 1.6e. Session reuse across streams
is refused in either recording order. Rounds do not change task state.
The full encoded Event, including duplicated retry and finding JSON, must
fit below the existing scanners' 1 MiB line budget; oversize refuses before
append. Retry payloads must match the immutable stored round, and request
field digests preserve integer line numbers exactly.

`review list` returns full rounds with derived other-model/current flags
and finding open state; text uses the shared box, wrapping and escaping.
A previously computed binding for the latest hand-off's exact base/head
can be read from its round, keeping restored history readable without git
objects. New endpoints require the same local fingerprint computation;
reads never fetch. Required profiles live in one shared reader of the
hand-off's frozen lane, changed paths and policy snapshot.

**1.6d's staging (2026-11-15, orchestrator under Mark's delegation).**
V8 records both streams, all three finding severities and build holds with
and without models through public commands. Its public-command stale-round
journey moves to 1.6f, whose rework makes a second hand-off possible. This
PR tests freshness at the domain level using second hand-offs written as
journal Events through the existing hand-off writer, without raw SQL:
edits stale earlier rounds and a clean rebase with the same fingerprint
keeps them current. No repeat hand-off command or transition bypass joins
1.6d.

**1.6e's choices (2026-11-15).** Optional `dispositions` in a round name
an earlier finding, `resolved`, `not-applicable` or `open`, and a nonblank
reason. Omission and `[]` mean no dispositions. Original rounds and findings
stay immutable; dispositions stay in the round, and terminal ones add
`finding_closures` keyed by finding, with exactly one closing round or
person's-decision ID. `finding.closed` (`finding.close`) is a new lean
Event-v2 family. A person's `close-finding` decision keeps the task as its
subject and binds to the finding: this preserves the issue foreign key and
all existing kinds' subject rules while naming exactly what was decided.
Its closure disposition is `closed`, without asserting the fault was fixed.
No terminal is needed; optional keyed retries use the existing person's
command transaction and compare stored retry output with immutable closure
evidence. Both closure Event forms enforce the journal's existing below-1-MiB
line budget before append. Legacy and unrelated writers cannot insert the
new closure table or decision kind.

Migration 31 extends the decision kind and adds the closure table;
projection cache version 13 rebuilds older projections. The migration
rebuilds referencing questions alongside decisions, preserving answered
questions and their foreign keys (SQLite's deferred FK counter does not
clear after dropping and recreating a populated referenced parent alone).
`reviewrecord.OpenMajors` reads all unclosed majors of a task in journal
order, independent of fingerprint and stream; admission and later `done`
share it. Dispositions use existing projection rows as evidence of journal
precedence, never caller times. Already closed findings cannot receive a
second disposition; identical command retries retain their original Event.
Raw CHANGES verdicts are never rewritten after a person closes a major.
V8 writes a major closed by a round, one by a person, a disposed minor,
a nit and a major still open, through public commands. The cross-fingerprint
closure test extends 1.6d's journal-backed second-hand-off fixture; public
rework and its V8 journey still belong to 1.6f.

**1.6f's choices (2026-11-15).** Migration 32 adds nullable
`holds.rework_at`; the ended build hold keeps its hand-off time and `review`
end reason, while its kept mode becomes `rework` from that round’s time.
Projection cache version 14 rebuilds older projections. Existing
`review.round-recorded` and `hold.taken` families suffice; both closed
catalogs include the optional, exactly validated rework effects. The same
review Event bumps only an `in_review` task’s revision and changes its own
kept hand-off reservation. CHANGES while already `in_progress` and APPROVE
never change state or a hold. A linked rework hold is allowed only for the
same task and worktree; its model joins all other build-hold attribution.
Rework replacement admission runs before the Event-v1 compatibility return,
resolves legacy SQL table/key spelling, and requires a full Event-v2 new hold
even when historical optional fields such as `last_heartbeat_at` are omitted.

`next` returns rework before unblocking and new work; an existing live or
interrupted session still resumes first. Repeating next after taking
rework is an ordinary resume of that same hold. Rework elsewhere reuses
`interrupted-elsewhere`, the existing informational reason for kept work;
human text names rework. There is no second start:
the task is already `in_progress`, so the fix goes directly through check
and hand-off. Takeover and cancellation keep their existing command paths.
V8 now uses public commands for the full fix/check/hand-off/approval loop,
with stale first rounds and a later round closing their remaining major.
A separate V8 task preserves its still-open-major coverage, since an
approving round cannot leave its task’s major open. The earlier journal-only
second-hand-off fixture now uses public rework too. No new dependency.

**Approval** (moved from 2.3)

21. **`fid approve <task> --actor human:<name> --reason <text>`.** Standard
    input and standard error must be terminals: the prompt goes to standard
    error and an optional JSON result to standard output. Agents and
    services are refused, whatever the terminal. The prompt shows the full
    diffstat, the check results, each stream's review verdict, any stale
    approval and "no other-model review" when it applies, then asks for an
    explicit typed confirmation. A person may approve while other gates are
    still missing (the prompt shows them); `done` refuses those. The
    approval presentation in `fid status` and `fid show` stays in 2.3.


**1.6g's choices (2026-11-15).** `fid approve` records `task.approved`
(`task.approve`) on a decision aggregate, using the single person's-decision
record with kind `approve`. Its task subject and nonblank reason remain
explicit; `binding_id` is the fingerprint and nullable `binding_json` holds
its algorithm, value, hand-off ID and recorded base/head. Other kinds keep
null structured bindings and their existing rules. Migration 33 rebuilds
the decisions and their referencing questions and finding closures together,
preserving both answered questions and closure foreign keys; projection
cache version 15 rebuilds caught-up older projections. Both closed catalogs
admit the family. Admission checks recorded evidence without git or terminals,
including task state at that journal position, exact latest hand-off endpoints
and fingerprint form. Legacy/other families cannot insert or mutate approvals.
The existing below-1-MiB Event budget is enforced before append.

The lock spans the full preview, exact task-ID confirmation and append.
Prompts go to stderr, results to stdout, including JSON. Cancellation or EOF
records nothing and exits 0. Both stdin and stderr must be terminals even on
a keyed retry. Optional retry keys cover actor, task and reason; an identical
retry validates its stored response against the immutable decision and returns
its original Event/time without another prompt or git computation. New
approvals use new decisions and keys. `approvalrecord.Current` accepts the
latest hand-off's computed fingerprint and returns current/stale/missing;
matching algorithm and value determine freshness, independently of endpoints.
Any matching immutable approval remains evidence of that exact change.

`fingerprint.Inspect` uses the same private git view and pinned diff settings
as computation for a full diffstat, with no file-count limit and name width
large enough for every path in the bounded patch. Reviews in the prompt are
the latest round of each required stream, with current/stale derived from the
fingerprint; open majors come from the shared unclosed-finding reader. The
existing other-model reader now treats a historical agent's unknown builder
model conservatively, retaining `no other-model review`; an empty human model
still means manual work. This closes a gap exposed by the approval disclosure
test, without changing model attribution rules. Older command fixtures now
attribute their seeded agent hold explicitly, preserving their other-model
expectations. The command inventories and keyed corpus include `approve`;
its integration retries use a real terminal even after the task enters rework.

V6's named test and `make check-approval-binding` drive real terminal approvals
through a clean rebase/new hand-off, edit/rework/new hand-off and re-approval.
Their fixture uses a live clock for each command, matching the terminal child
and avoiding a captured time from before UTC midnight.
V8 creates approvals before and after its public rework edit, then compares
both the stale decision and its current successor across backup/restore.
Prompt and JSON goldens, export/import, forged archive refusals and replay
cover the same record. No dependency is added. Approval in status/show (2.3)
and the done evaluator (1.6h) remain their scheduled tasks.


**1.6h's choices (2026-11-15).** `internal/doneeval.Evaluate` is pure over
one `Snapshot` and reports all fifteen gates, with satisfied, missing, stale,
waiting or refused states. Overall refusal takes precedence over waiting;
only an entirely satisfied snapshot is complete. Non-applicable gates are
satisfied with `not-required`. Missing prerequisites remain explicit; known
Assured approval, Rollback and model-disclosure gates are still evaluated
when policy/profile prerequisites are missing. The
command's public reason mapping belongs to 1.6i; no command, record, migration
or dependency joins this slice.

`Gather` takes a caller-owned projection query/transaction, its workspace,
task, trusted local origin/default-branch target and explicit GitHub facts.
The caller refreshes the projection and owns locking/revalidation. It reads
the latest hand-off's frozen policy, local results, current review views,
task-wide open majors, approval view, latest merge observation and its CI
receipts. Review/approval freshness uses a local fingerprint computation
when such evidence exists; an empty Routine change without review or approval
evidence does not introduce a fingerprint requirement. Shared
`reviewwork.CurrentViews` factors the existing list reader's freshness,
closures and other-model derivation. No GitHub call or journal write occurs.

Fresh GitHub CI facts are explicit too: pending attempts are deliberately
not durable, so old green receipts alone cannot prove current CI. The
adapter's existing pure CI selection is exported as `github.EvaluateCI`;
there is no second latest-attempt rule. Each selected green attempt must
also match a recorded receipt for this task and merge observation. A confirmed
merge receipt must bind to this latest hand-off and the observed PR/head tree.
These are existing recorded facts, not new records. All current review
rounds may establish other-model coverage; the latest current round per
required stream supplies that stream's profiles. A person's closure can
satisfy the major gate without rewriting a round's original CHANGES verdict.
The separate session is still required under same-model fallback.

`checkread.VerifyKept` reads only the required counting passes selected by
`checkread.Evaluate`, checks their stored bytes with the existing digest
verifier, and returns present observations without changing receipts or
availability. Hand-off and gathered done evidence share it. Hand-off's new
`kept-artifact-unavailable` refusal is a nonretryable precondition (exit 4),
with `{artifact_ids}` and guidance to restore bytes or check again. Rebuild
and identical hand-off retries never re-read kept files. Rollback recognizes
an exact level-two ATX `Rollback` heading, including optional closing hashes,
outside backtick/tilde fences; indented and blockquoted examples do not count.

The internal reason values are: `evidence-satisfied`, `not-required`,
`prerequisite-missing`, `task-unavailable`, `task-not-in-review`,
`handoff-missing`, `handoff-mismatch`, `lane-changed`, `policy-missing`,
`policy-mismatch`, `checks-not-configured`, `test-paths-not-configured`,
`checks-unsatisfied`, `kept-artifact-unavailable`, `artifacts-not-verified`,
`review-missing`, `review-stale`, `profile-missing`,
`separate-session-required`, `open-major`, `approval-missing`,
`approval-stale`, `rollback-missing`, `pr-unavailable`, `pr-mismatch`,
`repository-mismatch`, `branch-mismatch`, `not-merged`, `merge-invalid`,
`head-tree-mismatch`, `merge-missing`, `merge-mismatch`, `ci-waiting`,
`ci-not-green`, `ci-mismatch`, `ci-unrecorded`.


**1.6i's choices (2026-11-15).** `fid done <task>` / `task.done` accepts
work actors and an optional retry key, without a live hold. It observes the
hand-off's PR once through the injected `gh` reader, with origin's GitHub
identity and locally known default branch, outside the ledger lock. Under
the lock it refreshes and compares the complete gathered evidence and task
hold/question history, catching takeovers which do not bump the Issue
revision. Origin metadata is rechecked too. Superseded observations record
nothing. An already-replaced review reservation requires another hand-off.
Locked entry points to 1.6b's writers let confirmed matching merges and
completed CI attempts land before the evaluator, even on waiting/refusal.
A separately confirmed merge may be retained if the later CI request fails.

Migration 34 adds `completions`; cache version 16 rebuilds the projection.
The new lean family is `task.finished`, on the task aggregate, avoiding the
retired kernel's `task_completions` name. One Event inserts the immutable
completion, moves `in_review` to `done`, releases exactly its hand-off's
review reservation and resolves any open question as
`finished-before-answer` (null answer/decision). It may append one retry
result. Both catalogs admit the family; every replay/import validates the
record's task/latest-hand-off/merge/CI/review/approval bindings and exact
state/release effects. Its fingerprint and Rollback observation are present
evidence recorded at completion. Recorded evaluation reuses 1.6h's pure
evaluator without network, git objects or Artifact bytes; shared read views
now live in record packages so journal admission has no writer dependency.
The record distinguishes all completed CI attempts observed on this call
from the selected green attempts, preserving latest-attempt evaluation
without treating older receipts as a fresh GitHub observation. The Event
budget is below 1 MiB. Waiting/refused results never bind retry keys; a
completed retry validates against its immutable record and returns the
original response/Event before git/GitHub access.

Public precondition codes (exit 4) use exact `{gate, reason, record_ids}`:
`done-task-not-in-review`, `done-handoff-required`, `done-evidence-required`,
`done-checks-unsatisfied`, `done-kept-artifact-unavailable`,
`done-review-required`, `done-separate-review-required`, `done-open-major`,
`done-approval-missing`, `done-approval-stale`, `done-rollback-missing`,
`done-not-merged`, `done-ci-waiting`, `done-ci-not-green` and
`done-observation-stale`. The last four are retryable. Related repairs share
a code while `reason` retains the evaluator's specific gate reason; the CLI
spec lists the full mapping. `github-unavailable` is retryable with exact
`{reason, merged}`, preserving the recorded-merge knowledge boundary for the current hand-off.
An older hand-off's merge never confirms a replacement PR.
Unknown IDs retain `issue-not-found`; epics/subtasks are invalid arguments.
`issue transition` refuses completed-category targets for tasks, with
`use-done`, for every actor; subtask transitions stay unchanged. No new
dependency. Reconciliation/status wording, epics and V5's verifier remain
1.6j–1.6l. V8 writes merges, failed/successful CI and completion through
public commands, retaining its separate in-review, rework and unanswered
question subjects alongside the completed task.

Round 1 of #86 tightened completion admission to require the exact Effect
sequence: one completion insert, one task move, one hand-off reservation
release, each open question's resolution and only the optional retry result.
The stage grammar's minimum counts and a total Effect count alone allowed
an unrelated hold release or an unrequested retry result to replace a
question resolution. Live dry-run/projection, replay and import now refuse
these substitutions without changing the projection.

**1.6j's choices (2026-11-15).** The no-task form keeps the `task.done`
operation and one live/keyed command leaf. It takes a stable ID-ordered list
of in-review tasks with a hand-off, then calls 1.6i's `donework.Finish` for
each. It adds no evaluator, observation writer, record or dependency. Each
per-task refusal becomes a stdout outcome; ordinary command failures remain
nonzero. JSON `outcomes` have `task_id`, `status` and `completion` or
`refusal` (`code`, `message`, `next_action`, `details`). Statuses are
`completed`, `waiting`, `refused`, `github-unavailable`; not merged and
missing/pending CI wait, while failed CI is refused. Text includes each
refusal's gate and repair. Empty work succeeds with an empty array. A
no-task retry key is invalid usage: the existing natural observation keys
and per-task completion already make retries harmless.

Offline reads share `mergerecord.ReadReviewState`, binding receipts to the
latest hand-off of an in-review task with no completion. It reuses
`github.EvaluateCI` over completed receipts; no network or write is added.
Status keeps its IN PROGRESS grouping and snapshot and changes the row's
verdict; show adds a Merge row after Hand-off. Both JSON views add optional
`merge: {merge_commit, ci}` (`waiting`, `not-green`, `green`). Text says
"merged, waiting for CI" or "merged, CI not green". All-green recorded
receipts with another unsatisfied gate retain "in review" in text; they
cannot prove current CI is green. "Merged, waiting for confirmation" remains
`done`'s GitHub-unavailable message, as the approved 1.6j brief clarifies;
RESET's task wording now matches it. Rework, a superseding hand-off,
cancellation or completion removes the old merge view. V8 adds pending
and failed merge subjects through public commands and compares these new
status/show fields alongside every existing comparison.

**1.6k's choices (2026-11-15).** The explicit Issue's kind dispatches
`fid done` to the task path or the person's epic close, under the existing
`task.done` operation and keyed leaf. The epic path needs a human actor
and a nonblank reason, with no terminal or network. `epicrecord` reads all
direct task states and shares the fixed-state eligibility predicate between
command, journal admission and status's single snapshot. Empty epics may
close; only nonempty eligible epics show `epic-ready-to-close` in NEEDS YOU.
The command rechecks eligibility after refreshing under the write lock.
`done-epic-not-open` and `done-epic-open-tasks` are exit-4 preconditions with
`{gate, reason, record_ids}`; open child IDs are sorted and uncapped.

Migration 35 adds the `close-epic` decision kind with an epic-only subject;
projection cache version 17 rebuilds the cache. The `epic.closed` Event
inserts exactly one decision, moves exactly that epic to done with its
revision bump and optionally stores the retry result. Both catalogs admit
it; live projection, replay and import check actor/reason/subject/time,
no open task at that journal point and exact effect ownership/counts.
Close decisions are immutable, and an epic completion in any other family
is refused even through legacy or SQL case aliases. The command owns the
`todo → done` exception without adding a generic workflow edge. Keyed
retries return the original response and Event; changed requests conflict.
Dispatch consults a stored retry before looking up the requested Issue kind,
so a changed ID under an existing key keeps the task path's conflict refusal,
even when that ID does not exist. Each selected writer validates the complete
request before returning the stored result.
V8 adds closed, eligible and blocked epics through public commands and
compares their decisions, state and NEEDS YOU answers. No dependency,
epic validations/approval (3.1), V5 verifier (1.6l), or epic reconciliation
is added.

**1.6l's choices (2026-11-15).** V5's named verifier is
`TestV5DoneRefusesWithoutReceiptsAndAcceptsMergedGreenWork`, selected by
`make check-done-refusals` and included in ordinary `make test` discovery.
Its ledger setup and all evidence writes use public commands: agent work
starts through `next --claim --model` and `start`, then checks, hand-offs,
review rounds, finding closures and completion. Assured approval reuses
1.6g's real-terminal helper. The fixture's live clock matches those child
processes across midnight. The strict offline `gh` executable and
fetch-rejecting git wrapper come from 1.6j; GitHub facts use the public
hand-off's base/head and local git tree, with distinct PRs for the mixed
reconciliation journey. No SQL seed, internal writer or gate override
supplies evidence.

Each refusal checks its public code and exit status. `show` before and
after proves the Issue (including revision), checks, hand-off, questions
and worktree ownership stay unchanged. Confirmed merge and completed CI
observations may still be retained, as 1.6i requires. Successful completion
also checks that the worktree reservation is released. Missing checks,
kept bytes and CI, open majors (later-round and person closure), an edited
Assured approval, unmerged/mismatched PRs, failed/pending CI, every actor's
task/epic transition bypass, agent epic closure and mixed reconciliation
have separate cases. Routine with and without changed test paths, Standard
and Assured each complete with their required evidence. Product Go changes
by zero lines; no dependency or product behavior changes.

**golden-wrap fixture repair, 2026-11-15 (Routine).** Mark approved a
separate test-only fix for the path-wrapping flake recorded in 1.5e, 1.6g
and 1.6l. This is a friction-log entry in RESET.md, not a new task row.

Builder's transactional choices: one shared command-test helper creates a
private temporary directory, resolves its absolute path, and pads a child
to exactly 120 display columns, measured with the renderer's
`display.VisibleWidth` so non-ASCII parents have the same width as ASCII
parents (round 1 of #91). This accommodates long TMPDIRs while keeping
refusal details below their existing 200-character cap. The fixture puts
its workspace, `config` and `home` beneath it, giving each path a constant
total display width and preserving the existing workspace basenames. The
lifecycle fixture now creates its paths directly instead of moving its
workspace to another random directory. The linked worktree used by health
goldens uses the same helper, replacing its separate padding scheme.
Cleanup removes the private parent, including padding. A parent that leaves
no padding room fails clearly and asks for a shorter TMPDIR or fixture
parent. A deterministic over-budget regression pins that fatal diagnostic
and checks private-allocation cleanup after refusal (round 3 of #91).
A resolved parent containing any whitespace fails before allocation,
with a diagnostic requiring a whitespace-free TMPDIR (round 2 of #91);
wrapping splits paths at whitespace before masking. Regressions pin that
fatal diagnostic for ASCII spaces, tabs and Unicode spaces. No dependency
or production code changes. Path placeholders also cover HOME, and still
precede rework's hex-ID masking. Regressions compare real status/show output
after masking under short, long and non-ASCII parents, including HOME's
config fallback, and check the allocator's fixed display width and cleanup.
Linked-worktree health goldens also cover these explicit parents.

## Decided 2026-11-17

Mark made these decisions on 2026-11-17, after 1.9 (#93) merged. Items 1
and 2 are also in their tasks' rows in `RESET.md`.

1. **1.9b and 1.9c join the reset**, both Assured, from 1.9's audit (G10
   and G11 in #93). 1.9b runs before 2.1, since the golden path relies on
   `done` meaning done. 1.9c runs alongside phase 2. Their scope, moved
   here from the Later list:
   - **1.9b: admission owns `issues.status`.** The commands own each
     fixed-state move, but admission doesn't: a forged `issue.transitioned`
     (or Event-v1 history) can move a task to `done`, `in_review`,
     `in_progress` or `cancelled` on every path, around `fid done`, the
     hand-off, `fid start` and the person's cancel decision (probed in
     #93), and `issue.updated` can change `status` too. Share
     `internal/issues/fixed_states.go`'s rules with admission, decide what
     Event-v1 history may still move, and rebuild the test fixtures that
     move tasks with plain transitions (`cmd/fid/review_state_fixture_test.go`
     and others).
     - 2026-11-21, settled in 1.9b: Event-v1 history can't move a task or
       an epic. It may still create one, in todo only, and create and move
       subtasks, in one Event too. As on Event-v2, a subtask's move comes
       only from the family that owns it, `issue.transitioned`. None of the
       commands that own a fixed move ever wrote Event-v1, import already
       refuses it, and no ledger holds any (Mark's live ledger holds only
       Event-v2). On every path, whatever its schema, a task or an epic is
       created only in todo.
   - **1.9c: row identity and cardinality for the Issue families.**
     Admission binds the `issues` row to the Event's Issue, by `id` alone,
     but not `acceptance_criteria.issue_id`, `issue_labels.issue_id`,
     `activity_log.ref_id` and `journal_event_id`, `labels.name` or the
     `links` key, and a stage may repeat. Two gaps were probed in #93's
     first fix round: an `issue.updated` can set `id` in its values and
     rename its own Issue (`FID-a` became `FID-z` on commit), and a
     `link.added` or `link.removed` may bump the revision of any Issue, not
     only the link's two ends. Import still checks more than live append
     and replay: `*_event_id` references, immutable columns on update (which
     would refuse the rename), column names and types against the
     archive's schema, and duplicate JSON keys. Move each into the shared
     rules or record why it stays import-only.
     - 2026-11-22, settled in this PR (RESET 1.9c): shared admission binds
       Issue-family child rows and activity to the Event's own subject,
       label definitions to its label, and both link rows and revision
       updates to the declared link and its inverse. Each stage is bounded;
       update columns and label operations cannot repeat. A creation starts
       at revision one; an update has one scalar row and one revision row,
       guarded by its exact prior plus one. Issue identity,
       kind, workflow, parent and creation facts are immutable, for Event-v1
       too; no Issue deletion/recreation is admitted. Epic creation's stages
       follow its writer: Issue, optional metadata, labels, activity, retry.
       Link commands and admission share one relation-inverse registry.
       Import-only audit: duplicate JSON keys now fail in the shared Event
       decoder; shared Issue column/type/key checks complement the lean
       family validators; the import update allowlist moves into shared
       admission, including historical profiles before the lifecycle guard
       existed. Import calls that same rule before schema validation to keep
       its existing refusal classification. Transition hold subjects and
       cancellation end reasons are checked on the common hold columns,
       including those historical profiles.
       Import keeps column/type validation against its declared historical
       schema (a different profile from the live database), canonical archive
       lines (an integrity/manifest contract), and existence of historical
       `*_event_id` references (the projection has no Event index; standalone
       projection admission has only an Event and a transaction). Own-Event
       bindings, including activity's `journal_event_id`, are shared.
     - 2026-11-22, #101 round 2: duplicate object keys use Unicode simple
       folding, exactly as `encoding/json` matches field names. Lowercasing
       misses long-s aliases and wrongly aliases dotted I; an ASCII-only
       fence would also refuse otherwise valid structured details. The same
       recursive check covers the complete serialized Event in shared
       admission before live commit or dry run, as well as decoding stored
       Events for replay and import. Structured maps cannot write an Event
       that its replay decoder would refuse.
     - 2026-11-22, #101 round 3: limit Unicode simple folding to objects
       decoded into structs (Event wire fields and typed embedded records).
       Map keys match exactly, as `encoding/json` and case-sensitive readers
       do; exact duplicates are still refused in every object. In particular,
       check names `build` and `Build`, their hand-off definition digests and
       structured activity details with `A` and `a` are distinct entries.
       One typed scan checks the serialized Event before admission and when
       decoding stored Events, so admission and replay keep the same rule.
     - 2026-11-22, #101 round 4: shared retry admission also scans the embedded
       `result_json` text before decoding it, including Issue-family retries.
       Exact duplicate keys, at any depth, are refused on every journal path
       and import. The create, update, start and hand-off stored-result readers
       use the same typed scan before decoding: struct aliases are refused,
       map keys match exactly, and unknown response fields remain allowed.
       This resolves the round-3 Later item about retry `result_json`; it is
       no longer an import-only check or deferred work.
     - 2026-11-22, #101 round 5 (separate review): a shared schema-convention
       scan checks every non-null `*_json` value in every Effect's `pk`,
       `values`, `prior` and `preconditions`, for every family and operation,
       including unused insert maps and Event-v1. It runs before family
       validation and before producer-side prior capture. Stored Event
       decoding runs the same scan before archive schema checks; the complete
       serialized Event is scanned before canonicalization can collapse raw
       duplicate keys. Exact duplicates are refused throughout each JSON document; Unicode simple-fold
       aliases are also refused where the column's reader decodes a struct.
       Typed destinations follow the actual record types and each retry
       family's shape; unknown nested response members and map entries stay
       case-sensitive. The scan also shares import's JSON number domain,
       refusing non-finite decoded numbers such as `1e999`. No Effect
       embedded-JSON syntax/key check remains import-only. Physical column
       availability still depends on the archive's historical profile;
       manifest JSON remains an archive-container contract. Insert
       prior/precondition maps are scanned rather than categorically banned;
       existing family rules still decide whether those maps are supported.
       A matrix derived from SQLite and the retained family catalog checks
       every JSON column/operation/map on all journal paths and public import;
       real criterion-writer Events pin the reported backup/restore failure.
       Matrix setup captures the current schema and seeds from the real
       migrations, verifies complete SQL/seed parity, and reuses that fixture
       to keep the full matrix within the unchanged test timeout. Historical
       profile probes and actual criterion-writer round trips retain the
       canonical migration directory; no cases or assertions were removed.
     - 2026-11-22, #101 round 6: admission and import's JSON-text schema check
       now reach one journal validator. After its destination-aware key scan,
       that validator calls `encoding/json.Unmarshal` into `any` for syntax,
       the standard nesting-depth limit and the finite decoded-number domain;
       the separate per-number check and import's duplicate-key scanner are
       removed. The matrix and real criterion-writer archive journey include
       10,001-level nesting in every tested map. Each generated import case
       retains refusal, unchanged-target and valid recovery-import assertions;
       its identical accepted-ledger write/export/validate/restore journey runs
       once overall. Four import workers avoid mass parallel restore fixtures;
       journal matrix paths reuse connections within each bounded prefix where
       the database file stays in place. Replay, freshness and rebuild paths
       reopen the live database path. The existing timeout
       and checks are unchanged. The creator-index regression merged from 2.1b
       explicitly admitted legacy identity edits only until 1.9c; it now pins
       their typed refusal, unchanged creator/index rows and journal bytes,
       and a successful follow-up write under 1.9c's approved identity rule.
2. **A scoped `next`** (2.4). `fid next --epic <epic>` serves only that
   epic's tasks, so a second loop can build a small feature while a long
   one holds the main line, without disturbing it. Inside the scope the
   usual order and rules apply; tasks outside it are passed over with a
   reason, and with nothing to serve in the scope, `next` returns `none`
   rather than work from another epic. An unscoped worker is unchanged,
   and holds still keep two workers off one task. 2.4 settles the reason's
   name and whether a scoped worker also takes review work outside its
   epic. Epic priority stays the tool for filling gaps (a side epic ranked
   just below the main line is served whenever the main line has nothing
   ready); a label scope is on the Later list.
3. **How 1.8 moves the reset into fid** (Mark's calls on the 1.8 playbook):
   - Every task from 1.9b on runs fid's whole golden path: take, start,
     check, hand-off, review record, approve, done.
   - 1.8 imports the fid-repo rows still open: 1.9b, 1.9c, 2.1 to 2.6, 3.1
     and 3.2. 3.3 (the kit's repo) and 3.4 (kauri) stay in `RESET.md`,
     since `fid done` sees only this repo's PRs. The Later list goes in as
     one epic, blocked by 3.2.
   - Phase 2 is a prioritised epic with its real dependencies as links:
     2.4 after 2.1, 2.5 after 2.4 and 2.6 after 2.5; 1.9b and 1.9c join it,
     with 2.1 after 1.9b. Phase 3 runs in sequence, blocked by phase 2.
   - While Codex both builds and reviews, Assured tasks show "no
     other-model review" at approval. Mark accepts it as the standing
     fallback, now visible in the record.
   - Review input until 2.1: the reviewer writes its markdown review for
     the PR comment and also `fid review record`'s JSON. The recorder keeps
     the markdown in the round's provenance and checks that the two agree,
     and a later round's prompt carries fid's IDs for the earlier findings.
     From 2.1, fid posts the comment from its own record (decided
     2026-11-03), so the reviewer writes the JSON only. A jig option
     that emits fid's input belongs to the kit's work (3.3).
   - fid runs from a binary pinned to a known `main` commit in
     `~/.local/bin/fid`, upgraded at each epic close or when a merged fix
     is needed (D11). (Amended 2026-11-22, item 6: rebuilt from `main`
     after merges.) Backups go to the private `fid-ledgers` repository
     over HTTPS, after `gh auth setup-git`, so the orchestrator and later
     the loop can back up too.

## Decided 2026-11-22

Mark made these decisions on 2026-11-22, after 1.9b (#97), 2.2b (#98) and
2.3 (#99) merged.

1. **2.2c joins the reset.** Mark accepted the proposed row: `fid waive`,
   the fingerprint-bound waiver 2.2b's size split left over, runs after 2.3
   and reuses its current/stale decision view.
2. **2.3b: a done task's receipt in `fid show`.** After `fid done`, show
   said only `done`, with the hand-off's PR number; nothing showed the
   merge, the CI run, the review rounds or the approval the completion
   record names. A new Standard row, 2.3b, shows them, in text and JSON,
   offline. It runs after 2.1.
3. **When status asks for approval** (2.1). Status asks for an Assured
   task's approval only once the latest hand-off's primary and separate
   reviews have approved it, as one NEEDS YOU item, "approve, then merge";
   other lanes get one "merge" item. This replaces 2.3's approval item,
   which appeared from the hand-off on, while reviews were still running.
4. **Finished work stays finished** (1.9c). fid refuses cancelling a done
   task or a closed epic. Undoing merged work is a new task, which follows
   the original PR's `## Rollback` section.
5. **2.1's design is delegated** to the orchestrator, as 2.2b's and 2.3's
   were. Its choices are recorded with 2.1, for Mark's review.
6. **fid's pinned binary follows `main`.** After a merge, while no task is
   between hand-off and done, the orchestrator rebuilds `~/.local/bin/fid`
   from `main`, after replaying a copy of the ledger with it. This replaces
   "upgraded at each epic close or when a merged fix is needed" (Decided
   2026-11-17, item 3).

## Decided 2026-11-23

Mark made these decisions on 2026-11-23, after 1.9c (#101) merged.

1. **1.9d: journal test time, before 2.1c.** 1.9c's generated
   malformed-JSON matrix sends 1,220 probes through seven journal paths.
   `internal/journal` took 530 s of Go's 600 s default in `make test`, and
   GitHub CI timed it out at 20 minutes on #103's head. The full matrix
   moves to the shared validators, with one probe per map and variant
   through each path and import. No timeout or check definition changes.
2. **2.1c joins the reset, before 2.5.** Mark accepted the proposed row
   (answers and deferrals; fid posts each round to the PR from its own
   record). 2.5's loop then needs no posting code of its own.
3. **2.1d: updating an approved task whose PR stopped merging,** before
   2.5. Until now only a round asking for changes gave a builder its
   worktree back, so a PR that stopped merging after approval needed a
   review round to find the conflict (#101 and #103). When the latest
   hand-off's head no longer merges cleanly into the local
   `origin/<default>`, which fid tests with git and never fetches, the
   builder's worktree gets the task back as rework without a round. It
   merges the default branch, runs `fid check` and hands off again. Rounds
   and the approval stay current only while the fingerprint is unchanged.
   A builder can't otherwise take an approved task back to fix an open
   minor: it answers or defers it (2.1c), so an approval keeps meaning
   that head.
4. **Two Later items.** Moving a task to another epic, and accepting every
   waiting task in one command (Later list, "From Mark (2026-11-23)").
5. **fid's own rework cap is 5.** `.fid.toml` sets `rework_cap = 5` in
   every lane: a builder gets rework for up to five rounds that asked for
   changes before the task waits for a person (`fid rework allow`). fid's
   built-in default stays 3. The policy change is its own Assured PR.

**Decided for 1.9d, 2026-11-23.**

- The full schema/catalog-derived embedded-JSON matrix (1,220 probes) runs
  directly through `ValidateEventShape`, `Event.UnmarshalJSON` and import's
  shared `journal.ValidateJSONText` boundary, without a repository per probe.
  The schema is still enumerated from the real migrations. Event refusals
  assert `RowIdentityError`, the Event ID, map/column and identifying reason.
  Untyped import JSON accepts case-distinct map keys; typed struct aliases
  are refused earlier by Event decoding.
- The first generated carrier for each of four maps and four variants is
  selected in stable order. These 16 representatives retain every path's
  refusal, complete projection/journal comparison and valid recovery write,
  and public import's refusal, unchanged-target and recovery import.
  Serialized/stored readers, destination-map rules, fixture parity and the
  accepted write/export/validate/restore journey stay.
- The orchestrator's continuation (2026-11-23) applies the same split to row identity. All 174 corpus
  probes (35 controls, 139 malformed) retain a complete admission boundary
  against a real projection: commit for Event-v2, projection for Event-v1,
  which commit cannot accept. The first carrier per schema/table/reason
  gives 51 malformed representatives on every applicable path. Including
  schema keeps all nine legacy probes on all six legacy modes. Every control
  still runs on every path, including every link relation/inverse value.
  Repository journeys fall from 1,209 to 681; the typed table/reason refusal,
  projection/journal comparisons and valid recovery-write body is unchanged.
  Import's full row corpus and historical-schema refusals remain unchanged.
- Twenty-six temporary source-overlay mutations are caught by test
  assertions on the final tree: the first pass's eleven embedded-JSON guards,
  each mode's row-refusal wiring, common row admission removed and tested
  independently through each mode, and a removed link-inverse rule caught by
  a non-representative probe at the complete commit boundary. No production
  seam, dependency, timeout, check definition or CI file changes.
- Standalone journal ran in 434.396s before the first pass and 356.300s after
  it. The continuation baseline on that tree was 366.780s; the final run was
  292.615s. Embedded-JSON representative paths fell from 26.39-60.14s to
  0.33-0.82s in the first pass. The focused row-identity run took 89.746s.
  `make test` journal fell from the first pass's 467.900s to 351.962s;
  export was 205.528s and `cmd/fid` was 410.169s. The under-300s journal
  aim is still unmet; the remaining hotspot is on the Later list.
  GitHub CI timing remains for the orchestrator to measure; these are local
  results, with the unchanged check definitions and timeout values.

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
  cap, interruption, parallel claims) are fid integration tests. The rework
  cap and the token budget live in fid, set per lane in `.fid.toml` (amended
  10 November 2026, ADR-0014). The rework cap counts a task's review rounds
  that asked for changes, 3 by default; the token budget sums 2.6's usage and
  is off until 2.6 has data. Past either, `next` skips the task
  (`rework-cap-reached`, `over-budget`) and `fid status` shows a NEEDS YOU
  item. 2.2a parses the settings, 2.1 and 2.6 enforce them. 2.4's scenario
  tests cover the rework cap; the over-budget scenario test lands in 2.6,
  with the code it tests, since 2.4 runs first (Mark, 11 November 2026). The
  kit ships one default loop; variants only when a project needs one. The
  first loop is written in fid's repo as fid's own copy, then promoted to the
  kit.
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

**D9 Facade first** (amended 10 November 2026: lean core, ADR-0014).

- **Lean core** (amended 10 November 2026). New golden-path commands are
  built on Issues, the journal and the check runner, plus a hold, a policy
  snapshot, check results, review rounds and one record for a person's
  decisions. They don't compose kernel commands. The kernel chain (Activation,
  Envelope, Obligation, Claim, Run, Attempt, Evidence Record, Stops,
  Attention Items, kernel budgets and Run recovery) is removed right after
  task 1.5a, in 1.5b–1.5d. Originally: "New golden-path commands are built
  from existing kernel commands."
- **One rule engine** (amended). The `done` evaluator (policy from
  `origin/main` plus the diff) is the single source of what done needs. It
  reads check results directly; there are no Obligations (amended 10 November
  2026).
- A golden-output test fails if human-facing output (including refusals)
  mentions kernel-only nouns (Envelope, Activation, Obligation, Attempt,
  Attention Item) or kernel IDs (`ENV-`, `STP-`, `ATT-`). `--json` and
  `--verbose` are exempt. (Amended 11 November 2026, item 7: `fid admin` was
  exempt too, and there's no `fid admin`.) A concept that keeps leaking is
  the evidence to collapse or rename it.
- Kernel commands are deleted with the kernel in 1.5c (amended 10 November
  2026). Originally they moved under `fid admin`, hidden from main help, and
  any the facade didn't use were deleted at trial end. The leak test above
  stays. 1.5c also removes Concerns, requirement and decision references,
  `fid capabilities`, `fid stats` and `fid workflow define` and `assign`,
  and there's no `fid admin` (decided 11 November 2026, items 6 and 7).
  (Settled 12 November 2026, following item 6: `fid workflow list` goes
  too, so the whole of `fid workflow` goes; `fid artifact` stays, moved
  onto the lean check results; and 1.5d drops the tables and catalog
  entries of Concerns and of requirement and decision references.)
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
    Discovery runs `git rev-parse` and refuses if git and `LedgerDir`
    disagree, so fid never writes outside `--git-common-dir`/fid. (This
    first said discovery runs it once; the 7 November details make it two
    runs.)
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
- **Details** (decided 6 November 2026, for task 1.2):
  - 1.2 records each backup's time and outcome. The backup-age line, with
    its 24-hour warning and 3-day NEEDS YOU item, moves to 1.3's
    `fid status`, which doesn't exist before then.
  - 1.2 creates `.fid.toml` holding only the project ID that keys the
    backup folder. A restore onto a fresh clone has no ledger, so the ID
    has to be committed. 2.2a adds the policy to the same file. Because it
    introduces `.fid.toml`, 1.2 is Assured.
  - The backup target is in user-level config,
    `$XDG_CONFIG_HOME/fid/config.toml` (`~/.config/fid/config.toml` by
    default). `fid init` asks for it only when run in an interactive
    terminal. Otherwise it doesn't ask, and `fid status` warns until a
    target is set.
  - Tests push to local bare repositories. Mark creates the real
    `fid-ledgers` repository before the manual restore drill.
- **Details** (decided 7 November 2026, for task 1.2):
  - Mark approved splitting 1.2 into three Assured tasks: 1.2a writes the
    project ID to `.fid.toml` and asks for the user-level backup target only
    in a terminal; 1.2b exports into fid's own backup clone, commits and
    pushes, and records each backup's time and outcome; 1.2c builds
    `fid restore` and proves V8 (back up, wipe, restore, compare).
  - `fid restore` is a dedicated command with no arguments. It reads the
    backup repo from user config and the project ID from `.fid.toml`,
    fetches the backup into a temp directory, imports that project's
    folder, and refuses if a ledger already exists. V8 exercises this
    command. This supersedes the 6 November proposal to restore manually
    through `fid export import`.
  - fid keeps its own clone of the backup repository, for example under
    `$XDG_CACHE_HOME/fid/`. It never commits in a human's checkout of the
    backup repo. 1.2b picks and records the exact path.
  - **Builder's choice, PR #45:** `.fid.toml` uses
    `[project] id = "project-<32 lowercase hex digits>"`:
    128 random bits, the ledger's existing generated format. A new ledger's
    config and genesis event use that same ID; no journal format changes.
    Init keeps an existing valid file byte for byte and tells the user to
    commit it; fid never commits it itself.
  - **Builder's choice, PR #45:** the user config setting is
    `[backup] target = "<git URL or local path>"`.
    Stdin and stdout must both be terminals to ask; `--json` and an already
    configured target suppress the prompt. An empty answer skips setup.
  - **Builder's choice, PR #46 (1.2b):** fid's clone is a bare repository
    at `$XDG_CACHE_HOME/fid/backup/<key>/repo.git`, `<key>` being the first
    32 hex digits of the target's SHA-256, with a `flock` beside it that
    serializes backups through it (2-minute wait). Each project's folder on
    the target's `main` holds `export.zip`, the existing export archive,
    validated before the push, and `backup.json`, a manifest of its digest
    and counts for 1.2c. fid commits as `fid backup <backup@fid.invalid>`,
    unsigned, without hooks or the user's identity.
  - **Builder's choice, PR #46 (1.2b):** when nothing changed, fid makes
    no commit and records a successful attempt with outcome `unchanged`.
  - **Builder's choice, PR #46 (1.2b):** the record of each attempt is
    `<ledger>/backups.jsonl`, one JSON line per attempt that found the
    ledger (`pushed`, `unchanged` or `failed`, with a failure kind). It sits
    in the ledger directory beside the journal and isn't part of it: no
    Event, no journal or projection change. It travels with the ledger, and
    clearing the cache doesn't lose it. `.fid.toml`'s ID must match the
    ledger's, so a copied identity can't overwrite another project's folder.
  - **Builder's choice, PR #46 (1.2b):** failures reuse existing codes:
    `invalid-argument` for no target, malformed user config and a missing
    or mismatched `.fid.toml`; `lock-unavailable`; `repository-unavailable`
    with `not-a-repository` for an unreachable target; `internal-error` for
    a failed export. Git runs with `GIT_TERMINAL_PROMPT=0`, so unattended
    backups fail rather than wait for a password.
  - **Decided by Mark, 7 November 2026 (PR #46, 1.2b):** a rejected push
    (a server hook, or another commit reaching the target first) is
    `repository-unavailable` with a new reason, `push-rejected`, because no
    existing reason fits a remote that refuses an update.
  - **Decided by Mark, 7 November 2026 (PR #46, 1.2b):** `fid backup`
    refuses to replace a project's folder that another ledger of the project
    wrote, in 1.2b rather than on the later list. Clones of one project each
    make their own ledger with the project's ID, and so does `fid init` on a
    new machine in place of a restore; without the refusal, the last backup
    silently replaces the other ledger's, which may be its only copy.
  - **Builder's choice, PR #46 (1.2b):** a ledger is named by the hash of
    its genesis Event (`project.initialized`, the journal's first line),
    recorded in `backup.json` as `genesis_event_hash`. Each `fid init` makes
    its own, and an import keeps it, so a restored ledger goes on backing up
    into its folder. A folder whose `backup.json` is missing, unreadable or
    names no genesis Event is refused too. There is no override flag: a
    person clears the refusal by restoring from the folder or removing it
    from the backup repository.
  - **Decided by Mark, 8 November 2026 (PR #46, 1.2b):** that refusal is
    `repository-unavailable` with a new reason, `different-ledger`,
    following `push-rejected`, and is recorded as `different-ledger`. Mark
    chose it over a new conflict code that isn't retryable; the CLI spec
    notes that this reason persists until a person acts.
  - **Decided by Mark, 8 November 2026 (for 1.2d): backup storage for the
    long run.** Measured on a scratch ledger (1,000 issues, then 60 backups
    of 5 issues each): `export.zip` is about 8% of the journal's size. Each
    extra backup adds about 12 KiB to the packed backup repository, about
    the same as committing the plain journal files (11 KiB). But fid's clone
    in the cache never packs: fid uses plumbing commands that never trigger
    git's automatic packing, and fetches only when another machine moved the
    tip. So the clone keeps a whole archive per backup (10.4 MiB after 61
    backups of a 0.2 MB archive). GitHub refuses files over 100 MB, which
    `export.zip` reaches at roughly a 1.2 GB journal. Mark chose to store
    the export's files unzipped in the folder: journal files of at most 64
    MB, `manifest.json` and artifacts. The journal pieces and manifest stay
    below 100 MB; Artifacts may exceed the target's per-file limit (the
    review decision below records this constraint). The backup history is
    readable, and restore zips the files back up and imports them through the
    existing path. 1.2d also packs fid's clone. No real backups
    exist yet, so the format change needs no migration. Rejected: keeping
    `export.zip` (the 100 MB ceiling), and storing the journal's daily files
    (most new restore code, for a gain packing already gives).
  - **Decided by Mark, 8 November 2026 (1.2d implementation):** the
    unzipped folder uses `backup.json` schema `fid.backup.v2`. Drop its
    `archive` field; keep `archive_digest` as the digest of the deterministic
    ZIP made by export and reconstructed on restore. Restore refuses v1;
    there are no real backups to migrate. Its export `manifest.json` must
    also name `.fid.toml`'s project ID before import, in addition to the
    backup wrapper's check. This closes the identity-check handoff from
    #48 (1.2c).
  - **Decided by Mark, 8 November 2026 (1.2d restore scratch):** put the
    whole private scratch directory, including fetched bare repository and
    re-zipped archive, at `<git-common-dir>/fid-restore-*`, and remove it
    afterwards. The directory is `0700` and archive `0600`. Import still
    holds its own archive snapshot beside the ledger, so restore needs disk
    for both archives plus fetched repository and staging. This avoids
    `$TMPDIR` potentially being memory-backed. It also resolves #48's
    code-reading finding that a relative `TMPDIR` fell back to the same
    relative path through `os.TempDir`, failing if that directory was absent;
    restore scratch no longer uses `TMPDIR`.
  - **Decided by Mark, 8 November 2026 (1.2d review):** document and defer
    the large-Artifact storage constraint; resolve it before permanent
    backups in 3.2. Export accepts an Artifact of up to 1 GiB, and v2 stores
    each Artifact as a single raw Git blob. Journal chunking does
    not protect those blobs from the backup target's per-file limit. A
    compressible Artifact could previously fit inside `export.zip` but
    exceed that limit unzipped. 1.2d changes neither Artifact limits nor
    storage; this decision does not select chunking or a smaller ceiling.
  - 1.2d packs the cache under its existing clone lock before exporting:
    count loose objects exactly, and at 256 run foreground `git repack -A -d`
    and `git prune-packed`. User `gc.auto` settings cannot disable it.
    Unreachable objects are retained; a maintenance error is recorded before
    a new commit or push. The 100-backup test proves actual packing with
    automatic maintenance disabled in the user's config.
  - Left-over handling extends D6 below: deferred minor review findings
    become discovered issues in 2.1; a later PR "Left over" section is
    derived from `discovered_from` links. Task 1.8 brings this record's
    later list into fid along with the remaining reset tasks.
- **Details** (decided by Mark on 8 November 2026, for task 1.2c):
  - V8's verifier, `check:restore`, is a named Go test in `cmd/fid`,
    `TestV8BackupWipeRestoreYieldsIdenticalState`, which `make test` and CI
    run; `make check-restore` runs only it, with `-count=1`.
  - V8 is table-driven, with two wipes: a new machine (a fresh repository
    holding only the committed `.fid.toml`, an empty cache, and the same
    user config and target) and in place (`<git-common-dir>/fid` deleted,
    everything else kept). The ledger holds issues, a label, a link and at
    least one kept Artifact, and the target also holds another project's
    folder, which restore leaves alone.
  - V8 goes through `fid backup` and `fid restore` only, never reading the
    target's folder, so it survives 1.2d's change to the folder's format.
    "Identical state" means: `fid export create` on the original before the
    wipe and on the restored ledger make byte-identical archives; the first
    `fid backup` after the restore is `unchanged` at the same commit; the
    journal's Events are byte-identical in order, compared as concatenated
    content rather than file by file; the ledger's config, every Artifact
    file and the read commands' JSON are identical; and `maintenance check`
    passes. The event count and genesis Event are compared with the original
    ledger's. The expected differences are asserted, not ignored: the
    journal arrives as numbered pieces (`00000001.jsonl`, …), and the
    restored ledger has `import.marker.json` and no `backups.jsonl`.
    `fid.sqlite`'s bytes aren't compared.
  - Failures add no codes or reasons beyond these. An existing ledger is
    `project-already-initialized` (exit 4, subject the project, detail
    `path`), checked first, before config or the network; an
    `import-target-conflict` from import, a ledger that appeared
    mid-restore, maps to the same code. A missing or invalid `.fid.toml`, no
    target, and malformed or unreadable user config are `invalid-argument`
    (exit 2). An unreachable target, or one that isn't a repository, is
    `repository-unavailable` with `not-a-repository`. A target with no
    `main`, or whose `main` has no folder for the project, is
    `repository-unavailable` with a new reason, `no-backup`. A damaged
    folder is the existing `archive-invalid` (exit 6, integrity, not
    retryable), with subject and `archive_path` `<project-id>/export.zip`,
    and new check codes as needed. Anything else is `internal-error` (exit
    7). Restore writes no backup record. Outside a repository and at an
    unsafe ledger location, it refuses as every project command does.
  - Before importing, restore fails closed unless `backup.json`'s schema is
    `fid.backup.v1`, its `project_id` is `.fid.toml`'s and the folder's, and
    its `archive_digest` is `export.zip`'s `sha256:` digest. Event counts
    and the genesis hash are asserted by tests, not checked at run time. The
    idempotency key passed to import is `fid-restore:<backup commit>`. There
    is no replay handling: a rerun with a ledger present refuses.
  - Output is one text line, `restored <project-id>/ from commit <40 hex
    digits> into <absolute-ledger-path>/`, or the JSON operation
    `project.restore` with
    `{restore:{project_id,commit,archive_digest,event_count,terminal_event_hash,genesis_event_hash},ledger_path}`
    and a null `event`. `--project` and `--json` work as they do for `fid
    backup`; there are no positional arguments and no other flags.
  - Restore fetches into a private (0700) temporary bare repository that it
    deletes afterwards, never fid's backup clone, and takes no lock. It
    shallow-fetches `main` only (`--depth=1`), so its time doesn't grow with
    the backup history, reads the folder's files with `cat-file`, writes the
    archive to a private temporary file and imports that. Git runs through
    the backup package's runner with the same sanitised environment, and
    `--template=` on init. Restore reads `.fid.toml` from the working tree,
    works from any worktree, and restores the ledger they share.
  - Turning the fetched folder into the archive to import is one function,
    chosen by `backup.json`'s `schema`; for `fid.backup.v1` it reads
    `export.zip` and checks its digest, and an unknown schema is
    `archive-invalid`. The checks that don't depend on the format,
    `project_id` against `.fid.toml` and the folder name, stay outside it,
    so 1.2d changes only that function.
  - The journal's clock guard covers restored Events. It took the latest
    date only from dated journal files, so just after a restore, when the
    journal is all numbered pieces, an append dated before the last restored
    Event went through. It now also covers the last Event of the last piece,
    read from the piece's tail, so an append doesn't read a 64 MiB piece.
    Order and Event IDs were already right and are unchanged.
  - **Builder's choice, PR #48 (1.2c):** the new check codes are
    `backup-manifest-invalid` (`backup.json` missing, not a manifest, of a
    schema fid can't restore, or naming another project; a file where the
    folder should be counts as a missing `backup.json`), `archive-missing`
    and `archive-digest-mismatch`. Each has a null `entry`, since it names
    no archive member; import's own validation failures keep their codes.
  - **Builder's choice, PR #48 (1.2c):** restore is
    `internal/backup/restore.go`, beside backup, sharing its git runner,
    folder names and manifest type; `folderArchive` is the one format
    function. The temporary directory is `$TMPDIR/fid-restore-*`, from the
    environment fid runs with. `project-already-initialized`'s subject and
    `path` are the project root, as for `fid init`. A ledger that another
    restore of the same backup installed mid-restore, which import would
    replay as a success, is refused like any other ledger that appeared.
- **Details** (decided 7 November 2026, for task 1.1d4):
  - Discovery runs git twice, one value per run: `git rev-parse
    --show-toplevel` for the top level and `git rev-parse
    --path-format=absolute --git-common-dir` for the common directory. fid
    strips exactly one trailing newline from each output and nothing else,
    because a path may itself end in a newline. This amends the 3 November
    detail that discovery runs `git rev-parse` once. Its point still holds:
    `LedgerDir` never runs git.
  - Why: one run printed both paths one per line, so a newline in either
    path made the output read more than one way. In #42's second review
    round, Codex showed that choosing a reading could be fooled into
    writing outside git's common directory, and that the remaining refusal
    for an output that read two agreeing ways turned away a safe layout.
  - The only refusals are the two recorded on 6 November: git and
    `LedgerDir` name different ledgers, or the common directory is inside
    the worktree but not within a `.git` directory. With one value per run
    there's no ambiguous output to refuse.
  - Both refusals have the message "the ledger location is unsafe".
  - The locations go in the free-text `detail`, not in new detail keys.
  - fid needs git 2.31 or later, for `--path-format=absolute`.
- **Details** (decided 8 November 2026, for task 1.3):
  - Mark approved splitting 1.3 in two. 1.3a builds the display package and
    `fid status` with the backup-age line; 1.3b builds `fid show` and a
    pager for long output.
  - NEEDS YOU is where everything waiting on Mark appears: a list of
    items, each with a kind, a summary and a next step. Each kind arrives
    with its task. 1.3a has only the backup; approvals, epic plans awaiting
    approval, untriaged work and tasks waiting for Mark come later.
  - `fid status` shows epics and tasks, never subtasks. An epic whose own
    state is completed or cancelled is hidden. An epic is in progress if
    its own state's category is `started`, or any of its tasks is started
    or completed; otherwise it's to do. Its progress is done over total:
    its completed tasks over its tasks that aren't cancelled. IN PROGRESS
    shows each in-progress epic with its bar and its started tasks, each
    with its workflow state; started tasks with no epic follow. A started
    task always appears, flat if its epic isn't listed.
  - The second section is TO DO, not READY, until 1.4. It lists every
    unstarted task under its epic, in-progress epics first, and names a
    task's blockers that aren't completed or cancelled, but doesn't judge
    readiness. 1.4 builds readiness, with `fid next`'s order and typed skip
    reasons, and turns the section into READY. `fid ready`'s pick isn't
    used, because it ignores dependencies. The text caps TO DO at 10 task
    rows; epic lines don't count, and JSON is never capped.
  - **Decided by Mark later on 8 November 2026:** `fid status --all`
    removes the cap, and the capped text's last line says how to see the
    rest: `+N more: fid status --all`. JSON is the same with or without it.
    Mark chose `--all` and a git-style pager over a separate to-do command:
    one flag rather than one more concept.
  - **The pager (built in 1.3b, not 1.3a).** It pages only when stdout is a
    terminal and the output is taller than the screen. It respects
    `$PAGER`, and otherwise runs `less` with `LESS=FRX`: quit if the output
    fits, keep colours, don't clear the screen. `--json` and piped output
    never page. If `less` is missing, fid prints directly.
  - The backup-age rules read the record, `<ledger>/backups.jsonl`, and user
    config. Mark accepted them:
    - a. The last backup is the latest line whose outcome is `pushed` or
      `unchanged`, at its `finished_at`.
    - b. Its age is now less the last backup. Over 24 hours is a warning;
      over 72 hours is a NEEDS YOU item.
    - c. A failed latest line, of any kind, is a warning however recent the
      last good backup. It names the failure kind, and how many failed in a
      row when that's more than one, so a run of `lock-unavailable` shows
      (the #46 left-over below).
    - d. No target set is a warning whose fix names the user config path.
      After 72 hours rule b makes it NEEDS YOU, and its next step names the
      config fix and then `fid backup`. Malformed or unreadable user config
      is a warning naming the file. Neither makes `fid status` fail.
    - e. A ledger never backed up counts its age from its genesis Event,
      `project.initialized`. A new ledger shows a plain "not backed up yet"
      line, with no warning, for its first 24 hours. A restored ledger's
      genesis Event is old, so it shows NEEDS YOU until one `fid backup`
      runs; Mark accepted this.
    - A torn or unparseable line is skipped, since appends repair the tail;
      a record that can't be read at all is a warning, not a failure.
      Relative times read as minutes, hours or days ago, absolute times are
      local, and a time in the future counts as now. Human text says "last
      backup failed", never "attempt" or another kernel noun, so 1.7's leak
      test passes.
    - When the backup needs you, it appears only as the NEEDS YOU item;
      otherwise it is the box's last line.
  - JSON is the operation `project.status` with a null `event`, holding
    everything the text shows, uncapped, in fid-cli.v1's conventions:
    lowercase kebab-case enum values, and null rather than omitted fields.
    The text is rendered from the same result. Colour is used only when
    stdout is a terminal and `NO_COLOR` is unset or empty; symbols are
    always Unicode, since fid has no display config.
  - The display package is ported from whet's `internal/display`, which
    Mark wrote. He authorised copying the parts fid needs into fid, under
    fid's MIT licence, and adjusting them for fid.
  - **"Ready for Mark to merge" is a NEEDS YOU item, in 2.1, not 1.6.**
    `fid status` makes no network calls: ADR-0013 allows only `done`,
    `backup` and posting review rounds. fid records reviews only from 2.1,
    so until then it can't know that a review approved a PR. From 2.1,
    fid's own record says that a task's latest review round approved it and
    that the task isn't done; the item's next step is to mark the PR ready
    and merge it. The item clears when 1.6's reconciliation runs `fid done`
    after Mark merges. RESET.md's 2.1 row says so.
  - The line budget. Mark, on 8 November 2026: "I think we set that as a
    guide to ensure we actually get to a usable trial and I don't end up
    generating lines of code and status prose to no avail. I don't want to
    willy-nilly just extend the budget, but the budget set was fairly
    arbitrary so I don't want to skimp and shortcut a worse implementation
    solely to meet the budget."
  - **Decided by Mark, 9 November 2026 (after #51):** #51 took product Go
    to 51,457 lines, about 200 over the budget of about 51.26k. Mark is
    fine going over in principle, since the figure was somewhat arbitrary.
    He asked for an assessment of genuine cuts only, "not cutting just for
    the sake of fitting this arbitrary budget". A Sonnet session's
    read-only assessment found about 210 lines of genuine cuts, listed on
    the later list.
  - **Decided by Mark, 9 November 2026 (PR #53):** remove `pricing.toml`
    and the dead V1 leftovers before 1.4, in one Standard PR: "pricing can
    go if we don't plan to use it anymore … And the dead leftovers too."
    Before the cut, the orchestrator checked that fid doesn't need the
    file. Nothing read it; `fid stats cost` sums the `cost_usd` stored with
    each run; and the reset measures tokens, not dollars (2.6, kill
    criterion 4). New ledgers no longer have the file, and old ledgers
    keep theirs, unread. #53 took product Go from 51,457 lines to 51,222.
  - **Builder's choice, PR #50 (1.3a):** the rules live in domain packages:
    `backup.ReadHealth` (`internal/backup/health.go`) judges the record and
    user config, and `issues.ReadOverview` (`internal/issues/overview.go`)
    sorts the epics and tasks. `internal/status` assembles the result and
    writes the text; `cmd/fid/status.go` only adapts.
  - **Builder's choice, PR #50 (1.3a):** the display package keeps the
    renderer and its colour roles (header, ID, error, warning, success,
    muted, label), the symbols (✓ ✗ ⚠ ● ○ ◆ →), the rounded borders, the
    progress bar, padding that cuts with `…`, and visible width. It drops
    whet's ASCII symbols, the birthday and scheduled symbols (◆ stays, for
    epics), date formats, message helpers and multi-line truncation.
  - **Builder's choice, PR #50 (1.3a):** the order is `issue list --order`'s
    convention extended: sequence, then priority (lower first), then
    creation and ID, a missing sequence or priority last. A task with no
    listed epic, none or a hidden one, is ordered with the others the same
    way, started or not.
  - **Builder's choice, PR #50 (1.3a):** an in-progress epic with nothing
    to do has no TO DO line, and a to-do epic with no tasks shows
    `no tasks`. `+N more` counts the hidden task rows, and each hidden epic
    with no tasks as one. The empty-project hint appears only when the
    project has no epics or tasks in any state, and names
    `fid issue create --kind task --title <title> --actor human:<name>
    --idempotency-key <key>`.
  - **Builder's choice, PR #50 (1.3a):** a blocker is read from either
    link, the blocker's `blocks` or the task's `blocked_by`, and may be any
    issue kind. The text shows the first by ID and `+N`; JSON's
    `blocked_by` lists them all, on every task.
  - **Builder's choice, PR #50 (1.3a):** the backup's states are
    `current`, `not-backed-up-yet`, `warning` and `needs-you`, and its
    reasons `no-recent-backup`, `last-backup-failed`, `no-target`,
    `user-config-unreadable` and `record-unreadable`. Any line that isn't a
    whole record is skipped, wherever it is, not only the last. fid never
    writes a damaged line before the last, so one there means the file was
    damaged by hand; skipping it can make the last backup look older, or
    hide a failure until the next backup. When the record can't be read,
    the age is unknown, so it warns but never needs you. The summary leads
    with the age once it's over a day, and otherwise with what's wrong; the
    text leaves out when a failure happened, which JSON's `latest` has. The
    next step is the first of: fix user config, set the target, fix the
    record, `fid backup`. Hours read as hours below 48.
  - **Builder's choice, PR #50 (1.3a):** a NEEDS YOU item is
    `{kind, summary, next_action, subject}`, `subject` being null for the
    backup and `{type, id}` for later kinds about an issue. The JSON field
    names are in the CLI spec.
  - **Builder's choice, PR #50 (1.3a):** the root command's context
    carries a `statusEnvironment` (environment and clock; the clock's zone
    is the one times are shown in), as `backupEnvironmentKey` does for
    backups. `NO_COLOR` is read from that environment. The goldens backdate
    the genesis Event and seed the projection directly, as `fid ready`'s
    tests do, so IDs and times are fixed; `go test ./cmd/fid -run
    TestStatusGolden -update` rewrites them.
  - **Builder's choice, PR #50 (1.3a) review round 1:** the text escapes
    every string it takes from the result (titles, IDs, workflow states,
    summaries, next steps) with `internal/cli`'s `EscapeControlText`, the
    failure renderer's helper, before fitting or colouring it; JSON keeps
    the values. A title cut to its column may cut an escape short, which
    leaves no raw control. A right-hand text that would leave a title fewer
    than 10 columns takes a line of its own beneath the row, flush right,
    rather than cutting an ID.
  - **Decided by Mark, 8 November 2026 (1.3b design):** `fid show <id>`
    mirrors 1.3a. A read function in `internal/issues` gathers the data,
    `internal/show` assembles one result and writes the text, and
    `cmd/fid/show.go` only adapts. It is a read: it refreshes the
    projection, records no Event and makes no network calls. The ID is
    required and matched exactly.
  - **Decided by Mark, 8 November 2026 (1.3b design):** the layout, from
    four approved mock-ups (a task with subtasks, a blocked task, an epic,
    a subtask). It is `status`'s 72-column box, with the ID, in the ID
    colour, and the kind in its top border. The heading is the mark, the
    title (wrapped, never cut) and the state flush right. The marks are
    ✓ ● ○ ◆, with ✗ for cancelled in the muted role, since cancelled isn't
    an error. An epic's state reads as `status` reads it. Labelled rows
    (10 columns and two spaces, `display.RoleLabel`) appear only when they
    have something, one value per line: `Task` before `Epic`, then
    `Progress` (epics), `Priority`, `Blocked by`, `Blocks`, `Labels` and
    `Spec`. Reference lines follow `status`'s row rules. ACCEPTANCE
    CRITERIA come for tasks and subtasks; an epic's TASKS lists every task,
    uncapped, and a task's SUBTASKS carries `N/M done`.
  - **Decided by Mark, 8 November 2026 (1.3b design):** sequence, revision,
    workflow ID and times are left out of the text, which JSON's `issue`
    has. Open concerns, triage status, the owner, an epic's scheduling
    (1.4) and the other link kinds are left out of both. Later rows (a
    NEEDS YOU block under the heading; `Lane`, `Working`, never "claim",
    `Checks`, `Reviews`, `Approval`, and `Next` last; an epic's `Plan` and
    VALIDATIONS) go after `Spec` and before the long sections.
  - **Decided by Mark, 8 November 2026 (1.3b design):** blockers both ways.
    A blocker is read from either link, any kind, and `Blocks` is the
    reverse. The text lists only open ones, so `Blocked by` means waiting
    now; JSON lists every link with each end's state.
  - **Decided by Mark, 8 November 2026 (1.3b design):** JSON is
    `issue.show` with a null `event`:
    `{issue, state, epic, task, blocked_by, blocks, progress, children}`,
    `issue` being exactly `issue get`'s (CLI spec section 3). Later tasks
    add `lane`, `working`, `checks`, `reviews`, `approval`, `needs_you` and
    `next`.
  - **Decided by Mark, 8 November 2026 (1.3b design):** refusals. An unknown
    ID, an ID of another record included, is `issue-not-found` (exit 3),
    whose text alone gains `next: fid issue list`; `issue get` and every
    other command keep theirs. A missing ID is `invalid-argument`; other
    refusals are every project command's, and anything else is
    `internal-error`, as for `status`.
  - **Decided by Mark, 8 November 2026 (1.3b design):** wrapping, for
    `status` too. Prose wraps at spaces to column 70, its continuation
    lines indented, and never breaks inside a word. A next step that is a
    fid command wraps with a shell continuation: ` \` at the end of each
    broken line, its continuation lines indented two more, so copying only
    the first line leaves the shell waiting. A title cut drops an escape
    that crosses the cut whole. The kernel-word check runs on `show`'s
    goldens through one shared regexp.
  - **Decided by Mark, 8 November 2026 (1.3b design):** the pager, as git's.
    It pages only on a terminal, and only text taller than the screen,
    measured with `TIOCGWINSZ` (`golang.org/x/sys/unix`, already a
    dependency) and counted at the terminal's width. `$PAGER` runs through
    `sh -c`; an empty `PAGER` or `cat` means none; otherwise `less`.
    `LESS=FRX` is set only when `LESS` is unset. Colour is decided before
    paging. A pager that can't start means fid prints directly; quitting
    early isn't an error; fid waits and its exit code is its own. There's
    no `--no-pager`, and only `status` and `show` page.
  - **Builder's choice, PR #51 (1.3b):** `issues.ReadIssueView` reads on one
    connection, sharing `Get`'s opening step. The overview and the view
    read an epic through one rule and count progress one way;
    `EpicProgress` became `Progress`, since it counts a task's subtasks too.
  - **Builder's choice, PR #51 (1.3b):** `Blocked by` and `Blocks` are in
    `status`'s order (sequence, priority, creation, ID), as the approved
    task mock-up lists them; a child's `blocked_by` IDs stay in ID order,
    as in `status`.
  - **Builder's choice, PR #51 (1.3b):** the workflow state `todo` reads
    `to do`, as the mock-ups show it; other states read with underscores
    as spaces. `status` never shows an unstarted state, so it doesn't
    change.
  - **Builder's choice, PR #51 (1.3b):** the heading's state stays on the
    title's first line, which wraps narrower to leave it room. When the
    state would leave fewer than 10 columns, or the title's first word
    doesn't fit beside it, the state takes a line of its own beneath,
    flush right, as reference lines do. The title is in the header role
    and the state muted. Marks are coloured by role (✓ success, ✗ muted,
    others plain), and an epic's mark is ◆ in any state.
  - **Builder's choice, PR #51 (1.3b):** IDs in reference lines are padded
    to the widest one listed, as `status` pads its rows. Ordinals are
    right-aligned to the widest, and `Labels` and `Spec` wrap under the
    value column. The `Progress` row appears when an epic has tasks; JSON's
    epic `progress` is always set, `{0, 0}` with none.
  - **Builder's choice, PR #51 (1.3b):** a command keeps each flag with its
    value when it wraps, as the brief's example breaks before `--actor
    human:<name>`. The continuation must fit within column 70 too; the
    last line needs none.
  - **Builder's choice, PR #51 (1.3b):** as git does, a `$PAGER` with no
    shell syntax or arguments, `less` included, is found on `PATH` and run
    directly, so a missing one is seen as failing to start and the text is
    printed quietly; any other runs through `/bin/sh -c`, where a missing
    program is the shell's error. fid catches an interrupt while it waits,
    since Ctrl-C typed in `less` reaches fid too. A screen whose size reads
    as zero doesn't page, nor does text exactly as tall as the screen.
  - **Builder's choice, PR #51 (1.3b):** the status environment became the
    display environment, which also carries the terminal's size, so tests
    set it. `display.BlockedText` and `display.StateText` serve both
    commands, and `cli.CutEscapedText` shares one cut with the failure
    messages. `show`'s goldens read one consistent project: the epic golden
    shows FID-0d4e8b blocked by FID-71a0c2, the task mock-up's link, where
    the epic mock-up shows it to do, and leaves out the mock-up's
    `Scheduling` row, as decided above.
  - **Builder's choice, PR #51 (1.3b) review round 1:** prose never wraps
    between a date and its time, such as `2026-11-05 14:00` in a backup
    summary: `display.Wrap` treats a word ending in `YYYY-MM-DD` and the
    word after it starting `HH:MM` as one word. The JSON summary keeps its
    plain space.

**D8 / F2 Parallel worktrees: minimal.** Shared ledger, atomic
`fid next --claim` on different issues, one loop per worktree, and a
two-worktree test. Mark expects mostly one sequential loop. (Amended
2026-11-17: a second loop can be scoped to one epic with `next --epic`, 2.4.)

- Claims last about an hour. The loop sends a heartbeat during long agent
  sessions, since a builder may not touch fid for an hour. Tasks waiting on
  Mark are parked and hold no claim. (Amended 10 November 2026, ADR-0014:
  the claim is a hold, with its own ID sent with every heartbeat, check, park
  and `done`. fid refuses a stale session's heartbeat, so its loop stops
  loudly. 1.5a designs parking while waiting for Mark.) (Amended 11 November
  2026, item 5: the agent parks with `fid park --waiting "<question>"`, the
  task stays kept for its worktree, and a person answers with `fid answer`.)
  (Decided 12 November 2026, item 12: while the task waits, that worktree's
  `next` returns `none`, saying it's waiting on Mark for that task; other
  worktrees carry on.)
- Interruptions (crash, usage limit): the loop parks the task as interrupted
  (`fid park`) and stops. `fid next` serves it first, to the same worktree,
  as kind `resume` with a new hold linked to the old one (amended 10 November
  2026; originally a linked successor Run). A worktree holding an interrupted
  task gets no other task.
- Lapsed holds (amended 10 November 2026, ADR-0014). A lapsed hold on a
  started task means the task is interrupted and kept for its worktree; other
  worktrees skip it (`interrupted-elsewhere`). A lapsed hold on an unstarted
  task frees it for anyone. Only a person takes over a task; the takeover ends
  any hold and is recorded, and 1.5a names the command. fid never touches
  git state. A task kept for a worktree that no longer exists is a NEEDS YOU
  item. (Amended 11 November 2026, item 4: the command is `fid takeover`,
  which gives the person a hold in their current worktree.)
- Shared-file conflicts are resolved at merge, as human developers do; the
  second PR's builder rebases as rework. ADR-0006 disjoint-scope proofs wait
  for real need.
- **Details** (decided 10 and 11 November 2026, for task 1.4):
  - **Decided by Mark, 11 November 2026 (1.4 design, step 6):** READY
    replaces TO DO in `fid status` ("A, as mocked"). It keeps status's
    epic grouping, task order, ten-row cap and `+N more: fid status --all`.
    Its right-hand column reuses `next`'s read-only order for the current
    worktree: `next`, `blocked by <ID>` (and `+N`), `after <ID>`,
    `working in <dir>`, `working here`, `epic blocked by <ID>` (and `+N`),
    or `epic closed`; otherwise it is blank. `next` can be below the first
    row, since unblocking work comes first in the order. While this
    worktree already holds a task, its row says `working here` and no row
    says `next`. IN PROGRESS replaces an interrupted task's workflow
    state with `interrupted in <dir>`; other started tasks keep theirs.
    `dir` is the worktree's directory name. Human text names no kernel word.
    JSON replaces `to_do` with `ready: {next, epics, tasks}`; `next` is the
    task ID or null, also null while this worktree already holds a task.
    READY tasks reuse `next`'s queue entry fields `reason`, `after`,
    `working` and `epic_blocked_by`; IN PROGRESS tasks gain `working`.
  - **Decided by Mark, 11 November 2026 (1.4 design, step 6):** tasks in
    `fid show` gain `Working`, after `Spec` and before the long sections.
    A live hold here reads `here, since 09:02, until 10:02`; elsewhere,
    `in <full path> by <actor>, since 08:40, until 09:40`. An interrupted
    task reads `⚠ interrupted in <full path> by <actor>: work stopped at
    07:15 without a park`. The stopped time is the hold's `since`, as
    `next`'s is; 1.5a switches it to the last heartbeat. Times are local
    `HH:MM`, with the date when it isn't today. `issue.show` gains
    `working`, null when nobody holds the task. Both commands use the
    queue's `{hold_id, worktree, here, actor, since, expires_at, lapsed}`
    shape, with `lapsed` true on interrupted tasks. A lapsed hold on an
    unstarted task frees it, so `show` has no Working row for it. Epics
    gain `Scheduling` after `Progress` and before `Priority`, as item 14
    records (12 November). There is no compatibility obligation before 3.2.
  - **Decided by Mark, 10/11 November 2026 (1.4 design):** a worktree is
    identified by the SHA-256 of its root path, so moving one makes a new
    one, and every worktree shares one ledger and one write lock.
    `fid next --claim` picks a task and records its hold in one Event under
    the write lock. The database backs the lock up with at most one open
    hold per task and per worktree; a lapse is read from the clock against
    `expires_at`, never stored, so it can't be part of a unique index, and
    the backstop firing is a bug (`internal-error`).
  - **Decided by Mark, 10/11 November 2026 (1.4 design):** a hold records
    its ID, the task, the worktree (ID and root path), the actor, `since`,
    `expires_at` (about an hour, a constant for now) and, for a resume, the
    hold it replaces. It ends when a new hold replaces it (a lapsed resume,
    or another worktree taking a lapsed hold's unstarted task), or when its
    task is completed or cancelled, in the same Event as that transition
    (Mark, 11 November); park, `done` and takeover end one from 1.5a and
    1.6. The record is lean: its own migration, Event family and catalog
    entries, and V8 compares it. Tables and Events say "hold", never
    "claim".
  - **Decided by Mark, 10/11 November 2026 (1.4 design):** the order is
    this worktree's own hold (a live one returned unchanged; a lapsed one
    on a started task resumed with a new hold, other worktrees skipping it
    as `interrupted-elsewhere`; a lapsed one on an unstarted task freeing
    it for anyone), then tasks other open tasks depend on, then unstarted
    tasks in `fid status`'s order. Skip reasons built now: `held-elsewhere`,
    `interrupted-elsewhere`, `blocked` (either link; a completed or
    cancelled blocker no longer blocks), `epic-blocked`, `epic-closed` and
    `not-next-in-sequence` (`after` names the task waited for). Sequential
    epics serve one task at a time; prioritised epics and epics with no
    scheduling go by priority; links beat sequence. `fid ready`'s pick,
    which ignores blockers, isn't reused.
  - **Decided by Mark, 10/11 November 2026 (1.4 design):** `fid next`
    reads (`next.get`, null `hold` and `event`); `--claim --actor` takes
    the task (`next.claim`); `--explain` works with either and changes only
    the text. No idempotency key: a re-run returns the hold it has. The
    refusals are `invalid-argument` for `--claim` without `--actor`,
    `--actor` without `--claim` or a malformed actor, `authority-denied`
    for a `service:*` actor, and `lock-unavailable` (retryable) after the
    existing 5-second wait; 1.4a adds no codes. JSON is `{kind, worktree,
    hold, task, queue}`, kinds `new`, `resume` and `none` with exit 0 for
    each, and `queue` the whole order in `fid status`'s task shape plus
    `reason`, `after` and `working`. Loops never branch on `reason`. The
    text follows the six approved mock-ups and never names a kernel word.
  - **Decided by Mark, 10/11 November 2026 (1.4 design):** V3 is
    `TestV3TwoWorktreesNeverHoldTheSameTask` in `cmd/fid`, so CI runs it:
    three worktrees take tasks at once, as separate processes, round after
    round, and the journal never holds two open holds on a task.
  - **Builder's choice, PR #56 (1.4a):** the hold is `holds` (migration 20)
    and its Event `hold.taken`, which ends earlier holds before recording
    the new one. The journal checks the family's Effects in stage order;
    the export catalog refuses it in archives before schema 20. An ended
    hold records `ended_at` and `end_reason` (`replaced`, `task-completed`
    or `task-cancelled`); `actor` keeps the whole `agent:<id>`. `next`
    reads the clock once, after refreshing the projection and, with
    `--claim`, under the write lock, so a hold that lapses while it waits
    for the lock counts as lapsed (PR #56 review, round 1).
  - **Builder's choice, PR #56 (1.4a):** the order holds every unstarted
    task and every task with an open hold; a started task with no hold
    isn't handed out, so it isn't listed. A sequential epic serves one task
    at a time across worktrees: while one of its tasks is held (live, or
    lapsed on a started task) it serves no other; otherwise it serves the
    first task in sequence on its next task's chain of blockers in the epic
    that nothing blocks, or the next task itself (PR #56 review, round 1).
    A sequence's waiting tasks are named `after` the task the epic holds or
    serves, or its next task when it can serve none; a task on the chain
    that is itself blocked shows as `blocked`.
    When reasons overlap the first applies: held or interrupted elsewhere,
    epic closed, epic blocked, not next in sequence, blocked; the
    `--explain` mock-up names `after` for a task that is also blocked.
  - **Decided by Mark, 12 November 2026 (PR #56 review):** a task that
    blocks an epic counts as unblocking work only when that epic has open
    tasks, which wait on it through their epic. A task blocking an epic
    with no open tasks isn't served early.
  - **Decided by Mark, 12 November 2026 (PR #56 review):** `next`'s JSON
    result gains `epic` (`{id, title, progress}`, progress as `fid status`
    gives it; null when the served task has no epic, and for `none`), and
    each queue entry gains `epic_blocked_by` (its epic's open blockers,
    `[]` when there are none). Both are additive; 1.4b reuses them.
  - **Builder's choice, PR #56 (1.4a):** `next` doesn't page, as the CLI
    spec has only `status` and `show` page. Without `--claim`, a task
    handed out with nothing passed over gets `→ fid show <ID>`; with
    `--explain`, the order replaces `Passed over`. "Work stopped at" shows
    the lapsed hold's last sign of life, its `since`, not its expiry
    (settled in the PR #56 review); 1.5a switches it to the hold's last
    heartbeat. Nothing to do at all reads `Nothing
    ready here: no tasks to do`; the counts for `epic-blocked` and
    `epic-closed` read `in a blocked epic` and `in a closed epic`, and work
    held in different worktrees is counted per directory.
  - **Builder's choice, PR #56 (1.4a):** `--actor` passes the parser and
    `next` checks it, so a malformed actor with `--claim` is refused under
    `next.claim`; `--actor` without `--claim` uses the request-level
    refusal (`argument: request`), so its human text names no kernel word.
    A `service:*` actor's human refusal names the actor alone, and its JSON
    keeps `authority-denied`'s details (PR #56 review, round 1).
    `make check-two-worktree` runs only V3, as `make check-restore` runs
    V8. V8's ledger gains two holds, one a cancellation ended, and V8
    compares the `holds` table read through `fid sql`, since `next`'s
    answer depends on the worktree it runs in.

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

- **Amended 10 November 2026** (ADR-0014):
  - fid reads `.fid.toml` from the local `origin/<default branch>` ref and
    never fetches the policy remote itself, so policy is as fresh as the last
    fetch. Backup and restore still fetch their own target; 1.6 decides
    whether `fid done` fetches. This replaces "after a fetch" above.
    (Decided 11 November 2026, item 8: `fid done` asks `gh` and never runs
    `git fetch`.)
  - An invalid `.fid.toml` on origin is refused, with no fallback to the
    built-in defaults.
  - A policy snapshot replaces the Activation: the source (commit X, or the
    built-in defaults), the effective policy's digest, the actor rules and
    the check definitions. fid writes a new one, as one Event under the
    lock, only when that digest changes, including after a fid upgrade that
    changes the built-in rules. Writing commands that need policy check it;
    read-only commands never write.
  - 2.2a defines the whole schema: checks, Assured paths, check-defining
    paths, `auto_merge` per lane, the rework cap and the token budget.

- **Decided by Mark, 11 November 2026 (2.2a):**
  - **Shape.** `[checks.<name>]` holds `run` (a non-empty argv list, with
    no shell), a required positive Go-duration `timeout`, `defined_by`
    (paths whose changes require Assured) and optional `when_changed`
    (the check is required only when those paths change).
    `[lanes.routine]`, `[lanes.standard]` and `[lanes.assured]` hold
    `checks`, `rework_cap` (default 3), `token_budget` (default off) and
    `auto_merge` (default false). Each check must be defined. Only Routine
    and Assured accept `paths`; all other paths are Standard. The highest
    floor of any changed path is the change's floor. `auto_merge` is
    refused under Assured, which Mark always merges. `.fid.toml` and
    `.github/workflows/**` are built-in, unwaivable Assured paths.
    `[review] test_adequacy` lists paths needing that review in every lane.
    Globs are rooted at the repository: `*` stays within one folder and
    `**` crosses folders. Unknown keys are refused. There is no schema
    version until real data needs compatibility.
  - **Parser.** BurntSushi/toml v1.6.0, with strict keys through
    `MetaData.Undecoded()`. fid writes its own glob matcher; there is no
    second dependency. One reader covers identity and policy. Init,
    backup and restore still read the working-tree project ID, accepting
    files that also contain policy.
  - **No checks.** A project can track work with no checks. A missing
    lane's `checks` key means not set up; explicit `checks = []` means no
    checks required. 2.2a exposes the distinction. The refusals for check,
    hand-off and done, and status's `⚠ No checks set for <lanes> in
    .fid.toml` warning, arrive in 1.5a and 1.6. Init keeps writing only
    `[project] id`.
- **Decided by Mark, 14 November 2026 (2.2a's proposals):**
  - **P1: A, as built.** Resolve the default through
    `refs/remotes/origin/HEAD`. With origin but no known default, refuse
    and name `git remote set-head origin <branch>` as the fix. With no
    origin, use built-in defaults.
  - **P2: C.** A missing `[review] test_adequacy` key means not set up,
    like a lane's missing `checks` key. Explicit `test_adequacy = []`
    means no test-adequacy reviews are needed. Hand-off (1.5a7) and `done`
    (1.6) will refuse while it isn't set up, naming the fix; neither
    refusal belongs to 1.5a1. Projects list their own test paths; fid
    cannot infer tests for every language. The distinction is part of the
    effective policy and its digest, recorded snapshots, replay,
    export/import and V8. No compatibility with older snapshots is needed:
    no ledger holds real data. Other defaults stay as built: rework cap 3,
    token budget off, `auto_merge` false, no lane paths (so every path is
    Standard except the hard floors), and missing `checks` = not set up.
  - **P3: B.** `policy-unavailable` replaces `invalid-argument` for both
    an origin with no known default branch and invalid or unreadable
    `.fid.toml` on origin. Its category is `precondition` (exit 4), it is
    not retryable, its subject is null, and its closed details are
    `problem` (string), `ref` (string or null) and `commit` (string or
    null). Text keeps the renderer's `error:` and `reason:` lines, then
    adds `next:`: `git remote set-head origin <branch>` for an unknown
    default; otherwise fix `.fid.toml` on the default branch, merge it
    and `git fetch`. Status warns with the same problem and fix.
- **Builder's implementation choices (2.2a):** migration 21 adds
  `policy_snapshots`, ordered by a sequence primary key. Projection format
  6 makes a caught-up format-5 ledger rebuild automatically, so an upgrade
  installs the new table even when no journal Events are pending.
  `policy.snapshotted`
  is one Event on the `policy` aggregate, written by fid's internal
  `policy.snapshot` operation under the same lock as the following hold.
  A snapshot can stand alone when there is no task or a live hold already
  exists. It records the full effective policy as canonical JSON, its
  SHA-256 digest, and `origin/<branch>` plus its commit (or `built-in
  defaults`). Set-valued lists are sorted and deduplicated, durations are
  encoded in nanoseconds and argv keeps its order. Project identity and
  source commit are not part of the effective digest. Built-in rules
  version 1, named action sets and actor rules are included. Consecutive
  identical effective digests are never recorded again, even at a new
  source commit. A snapshot followed by a failed hold remains a complete
  fact; retry checks it again. Replay validates the recorded policy and
  digest without replacing them with the running version's rules.

**Decided for 2.2b, 2026-11-21.**

Mark chose fid's trust boundary at this task's start: the code that decides
what an agent may do, or makes the receipts he relies on at approval and
`done` (D5). The exact `[lanes.assured] paths` list is:

```toml
paths = [
  "migrations/**", "internal/journal/**", "go.mod", "go.sum",
  "internal/policy/**", "internal/pathglob/**",
  "internal/fingerprint/**",
  "internal/persondecision/**", "internal/approvalrecord/**", "internal/approvalwork/**", "cmd/fid/approve.go",
  "internal/doneeval/**", "internal/donework/**", "internal/mergerecord/**", "internal/mergework/**",
  "internal/github/**",
  "internal/reviewrecord/**", "internal/reviewwork/**",
  "internal/checkrun/**", "internal/checkrecord/**", "internal/checkread/**",
  "internal/handoffrecord/**",
  "internal/issuestate/**",
]
```

Backup, export, issues, links, labels, next, holds, status, show, display and
`internal/cli` stay Standard. Releases have no paths yet; 3.2 adds them.
The built-in hard paths and check-defining paths still apply independently.

The orchestrator settled the following design under Mark's delegation, for
his review. The builder used the brief's explicit size split: this PR builds
the floor; `fid waive` is proposed as 2.2c, after 2.3. RESET's 2.2b done mark
applies to the narrowed floor slice; the proposed 2.2c row keeps the remaining
acceptance criteria open.

- One `policy.ComputeFloor` computes the floor and hard floor and the paths
  setting each. The input is the same committed HEAD versus merge-base
  comparison used by checks, against local `origin/<default>`, never fetched.
  Deleted paths and both names of renames count. Assured wins over Routine:
  built-in hard paths, any check's `defined_by`, configured Assured paths,
  and unsupported Git names are Assured; Routine paths are Routine;
  everything else is Standard. No changed paths means no floor.
- Only `.fid.toml` and `.github/workflows/**` set the unwaivable hard floor.
  Configured Assured paths and check-defining paths can be waived by a person
  once the follow-up exists. `check` selects the higher of the task's lane
  and floor without editing the task. Explicit check names retain their
  existing behavior and do not need a comparison.
- Hand-off computes its floor from its recorded comparison, checks the higher
  lane, and raises the task and records that lane atomically with its one
  state move. It never lowers. The receipt keeps the floor, its setting
  paths, the hard floor and hard paths, and the prior lane when raised.
  `done` rechecks the frozen policy's floor and retains all existing gates.
- Show computes current requirements using that same floor. JSON adds the
  floor, and the hand-off retains its immutable explanation. The Lane row
  keeps the recorded task lane and explains a hand-off raise using the
  first setting path and a count of others. A policy/comparison failure
  warns and retains the hand-off's explanation if known. No read writes.

**Approved waiver design (built by 2.2c, 2026-11-24):**
`fid waive <task> --lane <routine|standard> --reason <text> --actor human:<id>
--idempotency-key <key>` mirrors approve's interactive terminal requirement,
confirmation, retries and text/JSON replies. It records a `person_decisions`
entry of kind `waive`, requiring an in-review task and latest hand-off, a
reason and a strictly lower lane, never below the hard floor (the refusal
names its paths). Agent and service actors are refused. The confirmation
shows the task, current lane, floor and setting paths, hard floor and new
lane. Its binding contains algorithm, fingerprint, hand-off ID, base and
head. Only the fingerprint and algorithm determine currency: an unchanged
later hand-off preserves it; an edit makes it stale. The waiver changes
neither task lane nor hand-off; it changes every lane-dependent done gate,
including check requirements, review profiles, separate review, approval
and rollback. Without a current waiver the recorded hand-off lane applies.
The issue-update lowering refusal will name waive after hand-off. Show will
list current/stale waivers in JSON and explain them in the Lane row; status
will mark stale waivers. The orchestrator ordered 2.3 before 2.2c (2026-11-21,
under Mark's delegation): approvals already exist, so 2.3 builds the shared
current/stale presentation for a fingerprint-bound person decision, and 2.2c
reuses it for waivers. The new Event must join shared admission on
live append, replay and import, with catalog, restore/V8, typed-refusal,
terminal, staleness and golden coverage. Its Event name and refusal codes
remain the follow-up builder's transactional choices.

**Builder's implementation choices for the floor slice:**

- Keep the existing `task.handed-off` family and `handoffs.record_json`
  storage; no migration, new family, or refusal code is needed. The lane
  raise is in the same `issues` update as `status`, compatible with #97's
  one-state-move admission. The hand-off success response stays unchanged.
- Built-in policy rules version 3 records the floor semantics change in the
  effective digest. Earlier snapshots and receipts keep their historical
  pre-floor semantics on replay and done. Version-3 admission verifies the
  exact floor, prior lane, raised lane and atomic update without git.
- A check whose floor raises its selected lane retains the comparison paths
  in its start/retry receipt as base64 byte strings. Admission recomputes the
  selection rather than demanding the old task lane's checks are configured.
- JSON cannot retain invalid UTF-8 in string values. Such hand-offs retain
  `changed_path_bytes` alongside the readable `changed_paths`; `Record.Paths`
  supplies exact bytes to floor, conditional checks, reviews and done.
  Display paths use Unicode replacement for invalid bytes. Admission checks
  that the byte list agrees with the readable list. Other receipts need no
  extra byte list. This is required to honor the conservative unsupported
  path rule on replay, rather than changing its result after serialization.

**Decided for 2.2c, 2026-11-24.**

Built the approved waiver design above, in the Assured lane. A person in an
interactive terminal confirms the task, hand-off lane, computed floor and
setting paths, hard floor and paths, and target lane. A `waive` decision
binds the target lane to the latest hand-off's algorithm, fingerprint, ID,
base and head. Only algorithm and fingerprint decide currency. The latest
matching waiver sets every lane-dependent done gate, including checks,
review profiles, separate review, approval and rollback; it never edits the
task lane or hand-off. An edit restores the hand-off lane's gates. Built-in
hard floors cannot be waived. Show and status reuse `internal/decisionview`;
show explains the waiver in Lane, approval becomes `not-required`, and
status marks stale waivers. Both expose full waiver decisions with `current`.

Implementation choices:

- Event `task.waived`, command `task.waive`, aggregate `decision` sequence 1.
  `person_decisions.binding_json` holds `lane` alongside the five fingerprint
  fields; `binding_id` is the fingerprint. Migration 0039 extends the existing
  decision constraints; cache version 21 triggers the normal rebuild. These
  follow 2.1d's migration 0038 and cache version 20 in the stacked branch.
  No dependency is added.
- Refusal codes `waiver-not-in-review`, `waiver-handoff-required`,
  `waiver-not-lower`, `waiver-hard-floor` use precondition, non-retryable,
  with a `paths` array. Actor, terminal, reason, retry and cancellation
  behavior mirrors approve, including optional retry keys.
- `internal/waiverrecord` owns offline admission and currency;
  `internal/waiverwork` owns the locked transaction. The latest hand-off reader
  moved into `internal/handoffrecord`, with the existing review reader
  delegating to it, so review admission can use waived profiles without a
  package cycle. Historical profiles with no waivers need no new columns.
- Matching decisions are considered newest first, preserving a waiver when
  a previously waived change returns. New waiver packages and the command
  join fid's configured Assured trust-boundary paths under the existing rule.
- Shared admission enforces immutable decisions and typed binding/retry JSON
  on append, replay and import. Export/import retain current waivers; V8
  retains a waiver made stale by rework while keeping its existing Assured
  completion evidence. Display remains read-only.

- Integration with 2.1d (2026-11-24) retains both Event families throughout
  shared admission, embedded JSON, archive schema catalogs and V8. Default
  conflicts retain schema 38; waivers first appear in schema 39. Status checks
  merge conflicts before selecting the effective lane's approval action.
  Both tasks' done marks and decision blocks remain.

- A waived done task's receipt resolves the latest matching waiver from its
  recorded fingerprint and selects checks from the hand-off's frozen policy;
  Lane, Completed, approval and rollback show the waived requirements offline,
  even without git objects, while the task and hand-off retain their lanes.

Left for later: none discovered outside this task.

**Decided for 2.1, 2026-11-22.**

The orchestrator settled this design under Mark's delegation, for his review.
2.1 builds work selection; proposed 2.1b carries review delivery, answering
and deferring findings, and acceptance of agent-created discovered work.
RESET.md records the split and each row's done-when.

- `next` serves resume, rework, review, unblocking work, then new work.
  Its due-review reader shares done's primary and separate gate evaluation,
  required profiles and fingerprint freshness. Primary comes first; Assured
  adds separate once primary passes. Routine without a required profile
  needs no review. Rework needs another hand-off before review selection.
- Selection without an actor shows the first due review. Taking work skips
  actors found in any build hold, including ended ones (`built-by-actor`),
  and continues to eligible work. Review JSON carries the task, epic, latest
  hand-off's ID/PR/head/base, stream and required profiles; text names them.
- Implementation choice: a separate `review_claims` table (migration 0036),
  with `review.claimed` and `review.heartbeat` Events. It keeps builder
  reservations and builder/model attribution unchanged. It admits one open
  reservation per task and per reviewer worktree, under the journal lock.
  Review sessions return the usual hold JSON with an `RCL-` ID, renew for an
  hour with the existing heartbeat command, and lapse freely. The same
  worktree's live review precedes other reviews. A live reservation also
  prevents build work in that reviewer worktree, including when local
  comparison objects become unavailable. Builder kept-worktree information
  remains visible for missing-worktree NEEDS YOU items. A round ends its task/stream
  reservation atomically, regardless of who couriers it. Recording a review
  without a reservation remains supported. Admission validates frozen review
  gates using the recorded fingerprint, without requiring git objects on
  replay, import or restore; live selection computes it locally.
- All `CHANGES` rounds count, across both streams and every hand-off. The
  cap comes from the latest hand-off's frozen policy and recorded lane,
  including a raised lane. Rework is served at the cap and stopped past it;
  zero stops the first rework. The builder gets `none` and `waiting.kind`
  `rework-cap`, and explain says `rework-cap-reached`; other worktrees
  continue. Status names the count and cap and offers takeover or
  cancellation. A person's decision allowing more rounds for one task (one
  of the person's decisions the work model lists) is 2.1b's (Mark,
  2026-11-22).
- Mark's 2026-11-22 approval timing replaces 2.3's item from hand-off on.
  Only satisfied primary/separate review gates and the task-wide majors gate
  produce a NEEDS YOU action. Assured with missing/stale approval gets one
  `approval` item, "approve, then merge PR #n", and the shared no-other-model
  disclosure. Current approval or another lane gets one `merge` item and the
  PR URL. A recorded merge suppresses both; existing offline CI states show.
  Done clears them. Show's Approval row and status's per-task approval fields
  remain unchanged.
- Deviation required to preserve V5 unchanged: its person-closure journey
  completes with a current `CHANGES` round after a person closes the last
  major, without a later `APPROVE`. Selection and status retain that existing
  done-gate semantics, rather than adding a stricter verdict rule. Thus a
  person's closure can also make the approval/merge item appear. This is the
  closest implementation of the brief's shared-reader rule and V5 guarantee;
  its literal "only once reviews have approved" wording has this exception.
- Reads remain offline. An unavailable fingerprint never schedules review or
  merge; status preserves its board and warns (`review-unavailable` for
  other lanes, the existing `approval-unavailable` for Assured).
- No dependency added. Internal names, the separate table and the reuse of
  hold JSON/heartbeat are transactional choices made by the builder.

**Decided for 2.4, 2026-11-22.**

The orchestrator settled this design under Mark's delegation, for his review.

- `docs/specs/LOOP_CONTRACT.md` first declares the existing `fid-cli.v1`:
  envelope, command exits, selected task and session, all five next kinds,
  three waiting kinds, the person's answer, review hand-off/PR/stream and
  required profiles, session commands, reviews and done statuses. The version
  stays `fid-cli.v1`. Changing or removing a declared field, kind or exit
  requires a version bump; an additional optional field does not.
- The initial contract now covers 2.1b: next's latest rework allowance and
  guidance, and successful `none` while an unaccepted task occupies its
  worktree; no loop consumes the declaration yet, so it stays `fid-cli.v1`.
- Contract goldens select only declared fields and exits. The specification
  includes their machine-readable paths, with optional fields and array
  elements explicit. Public-command journeys compare the produced next and
  waiting kinds to the declared sets. They cover interruption and answers,
  review races, rework with a new hand-off and approval, the cap boundary,
  both review profiles and streams, and done's batch statuses.
- `next --epic <epic>` scopes every kind, including reviews, to that epic.
  Tasks outside it, including tasks without an epic, get `outside-scope`,
  ahead of other informational reasons. Within it the normal ordering,
  dependencies, sequence, holds and cap remain. With no available scoped
  work the result is `none`; unscoped workers retain their order. Unknown
  IDs use `issue-not-found` (3); a non-epic uses `unsupported-kind` (2).
- Implementation choice: scope is checked inside the planner's projection
  read (under the write lock for a claim), before any policy snapshot or
  session is written. The full queue is retained. An occupied build or
  reviewer worktree whose task is outside a newly requested scope returns
  `none`, preserving the existing reservation and the one-task-per-worktree
  rule. The CLI only passes the option and shapes typed domain errors.
- Implementation choices: fixture epics have explicit priorities, and
  sequential tasks have sequences at creation. Goldens mask only shaped IDs,
  commits and times; journeys assert task identity, distinct hand-offs and
  exclusive sessions. All fixtures use temporary repositories and explicit
  git test identities. No dependency, migration or journal change is needed.
- V9's `check:loop-selftest` verifier lands with 2.5's actual loop and stub
  agents. 2.4 supplies its declared version and kind sets. Its done-when is
  amended from "V4, V9" to "V4; the contract V9 checks" for Mark's review.

**Decided for 2.1b, 2026-11-22.**

Built from the orchestrator's design, settled under Mark's delegation for his
review. Lane: Assured. The brief's size rule is used after whole design item
2: acceptance and allowances are this PR; proposed 2.1c carries answers,
deferrals and posting. The RESET row describes that completed slice and marks
it done; the proposed row and Later list preserve all remaining scope.

- A task or subtask created by an `agent:*` actor waits for acceptance,
  regardless of creation date or existing holds. `fid accept <issue>... --actor
  human:<name> [--reason <text>]` writes one `accept` person decision per ID.
  A missing reason records "Accepted for work". Repeating acceptance returns
  the existing decision. Multiple IDs run in argument order, one Event each;
  if a later ID refuses, earlier accepted IDs remain accepted. There is no
  retry key: accepting an already accepted issue appends nothing. Human-created
  work needs no acceptance command. Services cannot make either decision.
- `next` skips unaccepted tasks as `not-accepted`, including an existing build
  or review reservation; an occupied worktree stays occupied. Unaccepted
  entries do not reserve the next slot in a sequential epic, but ordinary
  dependency links still block. Status lists every open unaccepted task or
  subtask under NEEDS YOU, kind `accept`, with its creator and acceptance or
  cancellation command. Completed/cancelled work has no item. Epic approval
  and D6's exception for discovered blockers in approved epics wait for 3.1.
- `fid rework allow <task> --rounds <n> --actor human:<name> --reason <text>`
  writes a `rework-allowance` decision for an open task. Rounds is a positive
  signed 64-bit integer; every allowance adds to the latest hand-off's frozen
  lane cap. The number of `CHANGES` rounds never resets. Allowances bind to
  the task, with no fingerprint. The latest allowance's attribution, rounds
  and reason accompany served work as `rework_allowance` in next's JSON;
  rework text shows its guidance. Status offers the allowance command first
  in a `rework-cap` item's next action. Existing builder reservations control
  which worktree gets rework again. Each invocation deliberately adds another
  allowance; there is no retry key. The displayed effective cap saturates at
  MaxInt64 to avoid overflow; no representable round count exceeds that cap.
- Implementation choices: migration 0037 extends `person_decisions`; an
  allowance's `binding_json` is exactly `{rounds: n}`. Two immutable families,
  `work.accepted` and `rework.allowed`, admit exactly one complete decision
  with the Event's person, subject, command and time, on append, replay and
  import. They admit neither unrelated Effects nor a retry result.
- `issue_creators` is a derived index of the actor in the canonical Event-v2
  `issue.created`, maintained within its projection transaction. No Event
  may write the index through Effects. Its foreign key follows Issue ID
  changes and deletion with update/delete cascades, preserving both valid
  legacy identity edits and the existing typed state-guard refusals. Cache schema 19 forces the normal
  rebuild of old projections, backfilling every existing typed creation
  actor. Event-v1 has no typed actor and makes no agent-creation assertion.
  Read-only fixtures with no creation Event likewise have no such assertion.
  `internal/workrecord` shares the offline readers; `internal/holdwork` owns
  decisions and `internal/journal` owns admission. No dependency is added.
- The admission tests preserve the existing prefix contract: a full replay
  can project valid Events before a later state-dependent refusal; the
  rejected Event changes no row or checkpoint. A shape refusal occurs in
  the pre-scan. Live refusals leave journal bytes and projection unchanged,
  and the next valid write succeeds. Export/import and rebuild preserve both
  decisions and the absence of acceptance for pending discoveries.

**Decided for 2.1c, 2026-11-23.**

Built the approved review-delivery design in lane Assured, stacked on 1.9d.

- `fid finding answer <finding> --actor <a> --text <why>` records an immutable
  builder/person response to an open finding of any severity. Every actor that
  held the task to build it counts as a builder, including ended holds. An
  answer never closes the finding. Review list and the next reviewer’s input
  include the responses in journal order.
- `fid finding defer <finding> --actor <a> --title <text> [--epic <id>]
  [--lane <lane>]` accepts minor/nit findings only. One `finding.deferred`
  Event creates the task, acceptance criterion, discovered-from link and its
  inverse, source revision bump and closing response. The criterion retains
  the finding ID, summary, evidence and requested fix. The task inherits the
  reviewed task’s epic unless overridden. Its lane starts at Standard (or a
  requested Assured lane), raised by the finding location’s floor under the
  response hand-off’s frozen policy. Check/hand-off still computes the floor
  of all paths eventually changed. Routine and unknown lanes refuse.
  A sequential epic requires a sequence: discovered tasks append after its
  highest existing sequence, allocated and admitted from the prior projection
  in the same Event. Other scheduling leaves the sequence unset.
- Migration 0040 adds immutable `finding_responses`: `FRS-` identity, finding
  ID, optional unique deferred task ID, complete typed `record_json`. Answers
  use `finding.answered`; deferrals use `finding.deferred`. Both use a
  `finding-response` aggregate and no request retry key. A closing response
  derives a `deferred` closure with response and discovered-task IDs; existing
  round/person closure records keep their schema. Agent deferrals derive
  `issue_creators` from the Event actor in the same projection transaction,
  so discovered tasks wait for acceptance. Deferral never reopens the source.
  Invalid responses and unknown post-round IDs use `invalid-argument` (exit 2).
  Error details name `finding` with builder/person guidance, or `round` with
  recorded-round guidance, rather than inheriting the review-record hint.
  Transport failures leave recording successful; `review post` returns pending IDs.
- CLI recording attempts delivery after its canonical write. Round comments
  follow jig’s headings, attribution, profiles, round, lane, purpose, finding
  severity/action/location/evidence/fix, earlier dispositions and verification.
  Answers and deferrals post short comments from their own canonical records.
  Delivery failures leave recording successful. `fid review post <round>`
  retries the round and responses to its findings, returning posted and pending
  record IDs. The latest hand-off **at recording time** fixes each comment’s
  destination, including a response’s recorded `handoff_id`: later hand-offs
  cannot redirect an interrupted delivery to another PR.
- Migration 0041 adds immutable `review_posts`, keyed by canonical record ID,
  with task/hand-off, PR, comment URL, render version (1), rendered-body digest
  and posting time. The version-1 renderer remains the receipt digest format;
  a future layout must retain it to admit existing history.
  `review.posted` is a `review-post` aggregate attributed to fid’s service.
  Admission binds every field to the canonical record and verifies that the
  URL names a comment on its recorded PR. Both catalogs register all three
  families. Round and response IDs keep distinct RND-/FRS- namespaces, which
  admission checks before IDs become delivery keys. Cache schema 22 rebuilds older projections. No dependency is added.
- Posting holds the journal lock across reconciliation, `gh pr comment` and
  receipt append. A stable HTML marker includes record ID and body digest.
  Retries with a receipt use no network; without one, fid checks paginated PR
  comments first, recovering an accepted comment whose local receipt was lost.
  A matching marker with edited body refuses rather than posting a second
  comment. The API listing, output and request duration are bounded. The body
  travels in a private temporary file, removed after use; no shell evaluates it.
  A receipt commit failure stops the batch and leaves remaining records pending;
  the next retry catches up the projection before any later receipt can advance
  its checkpoint past an appended but unprojected receipt.
- Pending delivery is derived from canonical records without receipts. Offline
  status warns once per affected round, with pending IDs and `fid review post`;
  it adds no NEEDS YOU item. `status`, `show`, `next` and review list stay offline,
  and delivery adds no completion gate. Test suites block the real `gh` by
  default; transport tests use a fake runner or executable. V8 compares both
  new tables, including posted receipts, answers and a deferred discovery.
- Admission tests send the full malformed matrix through the common boundary
  once, with representative append, replay and import refusals, unchanged
  state/bytes and subsequent valid writes. Temporary overlays check the new
  behavior rules. The completion report retains their results and check output.

- Update round (2026-11-24), stacked on 2.2c (#111): keep default conflicts
  at migration 0038 and waivers at 0039; finding responses and posting receipts
  follow at 0040/0041, with projection cache version 22. Admission, archive
  catalogs and V8 retain all three features together; no gate or delivery rule
  changes. Existing decisions blocks and task done marks remain.

**Left for later:**

- From this PR (2.1c): source inspection shows `review list` still uses the
  generic "an independent review of the latest hand-off" allowed-input hint
  for invalid task arguments. Give that existing read command task-specific
  guidance in later CLI polish; 2.1c corrects the new answer/defer/post hints.

**Decided for 2.1d, 2026-11-23.**

Built under Mark's delegation, for his review. Lane: Assured.

- Deviation to keep the existing human-output check: status JSON gives the
  exact `fid next --claim` action; human status says to take the task in its
  builder's worktree with `fid next`. The check forbids the literal flag in
  ordinary human output. No check is changed.
- Eligibility is the shared `doneeval.ReviewsDone` predicate used by status's
  approve/merge action, with no recorded merge, including 2.1's existing V5 person-closure exception
  for a `CHANGES` round. A still-due primary or separate review never takes
  this path. An Assured person's approval need not have happened yet.
- `git merge-tree --write-tree --name-only -z --no-messages` tests the latest
  hand-off head against the local `origin/HEAD` target, resolved by the same
  helper as policy. Exit 1 with conflicting paths proves a conflict; exit 0
  serves nothing. An unknown ref, missing commit or other git failure serves
  nothing and `next --explain` describes unavailable local evidence. Lazy
  fetching is disabled; fid never fetches. No checkout, index or ref changes.
- Only the kept hand-off worktree gets kind `rework`. Optional `rework`
  evidence names `default-conflict`, default ref/commit, hand-off, head and
  conflicting paths; non-UTF-8 names retain their bytes alongside readable
  paths. A conflict replaces the status approve/merge action with a
  `default-conflict` item describing the builder's update. Clean approved
  work remains reserved for review, including work with an open minor.
- Migration 0038 adds immutable `default_conflicts` receipts, keyed by a
  `DFC-<six-hex>` ID with task, kept hold and typed evidence JSON. The receipt
  also retains the observed fingerprint for admission without git objects.
  Routine work with no required reviews computes that binding when a
  conflict is found, without creating a review round.
  `task.default-conflict` is an Event-v2 work-actor family for `next.claim`,
  with a `default-conflict` aggregate. Exactly three ordered Effects insert
  the receipt, move its task `in_review` to `in_progress` with its next
  revision, and set that hold's `rework_at`, atomically. The normal linked
  `hold.taken` follows under the same write lock. If the process stops
  between those Events, the reservation and cause remain available.
  Existing allowance guidance and the task inspection action remain visible
  for both causes of rework.
- Shared admission binds the latest hand-off/head, task, kept hold, review
  currency, command, actor type, aggregate and exact Effects. It checks the
  recorded evidence shape on live append, replay and import, without
  rerunning git. Other families cannot move the state or set `rework_at`.
  A refused Event changes no journal bytes or projected rows; a subsequent
  valid write works. Both catalogs include the family; archive schema is 38.
  Cache version 20 forces existing projections to rebuild with migration
  0038 before the first read or claim; a focused upgrade regression covers
  an old schema-37 projection with a matching journal.
- This return creates no round and charges no cap allowance. Only `CHANGES`
  rounds count. The builder merges, checks and hands off again. Existing
  fingerprint readers retain rounds and approval when unchanged and mark
  both stale when changed; public-command tests cover both outcomes.
- `fid-cli.v1` stays current; the `rework` fields are optional additions to
  the loop contract and its selected-field goldens. The public integration
  journey approves two tasks, merges the first into the local default,
  returns the conflicting second to its builder, resolves, checks, hands
  off, reviews and finishes it. Import into a repo without comparison
  objects and V8's backup/wipe/restore retain the receipt.
- Transactional choices: a shared origin resolver; one receipt table rather
  than overloading review rounds; restoring the cause across an interrupted
  claim; exact new-task state in the claim response. Import comparisons use
  checkpoint Event identity/hash and schema/status, since archive chunk
  filenames and byte offsets are local storage details. No dependency added.

**Left for later:** none discovered within this task.

**Decided for 2.3, 2026-11-21.**

The orchestrator settled this design under Mark's delegation, for his review.
Approval presentation uses the latest hand-off's recorded lane: Assured,
including a lane raised by 2.2b, requires approval. Before the first hand-off
there is no approval state. Other lanes have no approval presentation.

- `approvalrecord.Current` remains the freshness rule: missing means no
  approval exists; current means an immutable approval matches the latest
  hand-off's fingerprint algorithm and value; stale means approvals exist
  but none matches. Current and stale name the decision, actor and time.
- Show places Approval immediately after Lane, with the display package's
  existing marks and colours: muted `○ missing`, green `✓ <decision>
  <actor> · <at>`, yellow `⚠ stale · <decision> <actor> · <at> · approved an
  earlier change`. A second line discloses `no other-model review` when
  applicable. The JSON fields are `approval` (null or an object with
  `state`, plus `decision_id`, `actor`, `at`, `fingerprint` when present)
  and nullable `no_other_model_review`. The open state string leaves room
  for 2.2c's future `not-required`; no waiver behavior is built here.
- Status adds an `approval` NEEDS YOU item for each in-review Assured task
  with missing or stale approval. The summary names the state and appends
  `, no other-model review` when applicable. Its next action is `fid approve
  <task> --actor human:<name> --reason <text>`. IN PROGRESS JSON carries the
  same approval and disclosure fields, including current approvals and the
  latest hand-off's evidence during rework; rework does not request approval.
- The approve prompt, done and both displays use one disclosure predicate
  over `reviewrecord.CurrentViews`: no current primary or separate round
  establishes another model family than the builders'. Unknown models are
  conservative, using the existing `OtherModel` reader. All current rounds
  may establish coverage, not just the last in a stream. The prompt's old
  loop did not filter streams; done's stricter task/stream filter is now
  shared. Admitted rounds already belong to these two streams, so this
  does not change their existing result.
- Both commands remain offline and record no Events. Policy or a current
  comparison failure retains evidence read from recorded hand-off endpoints.
  Unavailable fingerprint evidence warns with `approval-unavailable` and
  leaves unknown fields null; no approvals and no reviews remain known
  without computing a fingerprint.

**Builder's implementation choices:** `internal/decisionview.View` owns the
flat JSON shape and current/stale text for fingerprint-bound person
decisions, ready for 2.2c to reuse. Its approval reader composes existing
approval and review readers, without another freshness rule. Status reads
the approval alongside its other task details in its existing transaction.
Show uses a read transaction for the approval evidence. Times use the full
recorded RFC3339 timestamp in text and JSON. There is no new dependency,
migration, Event family, network call or write command.

**Decided for 2.3b, 2026-11-23.**

The orchestrator settled the completion-receipt design under Mark's delegation,
for his review. The builder implemented it in `internal/show`, lane Standard.

**Departure from the example, preserving the row's guarantee:** the existing
V2 vocabulary guard rejects both `run` and `runs` in human output, including
GitHub's `/actions/runs/` URL path. The first full test run caught that in
`TestMergedOfflineReadsGolden`. Under the brief's closest-design delegation,
text uses `CI ✓ success · #<id> on <commit>` and retains the CI URL in JSON
only. This keeps result, identity and commit visible offline without changing
or weakening the guard. The new text goldens use the same guard.

- Done tasks get a Completed block after the header rows, before criteria and
  children. It names the completion's time, actor and ID; merged PR number,
  commit and time; every counted CI run's conclusion, ID and commit; the
  completed hand-off's ID and head; fingerprint; counted rounds'
  IDs, streams, verdicts, reviewer actors and models; and the Assured approval,
  recorded other-model disclosure and Rollback presence.
- Checks use the exact required names and result IDs in that hand-off,
  displaying when they passed and `counts at <head>`. They do not compare
  finished work with today's HEAD, policy, changed paths or fingerprint.
  Unfinished work retains those comparisons. Text uses the existing wrapping,
  escaping and local clock, with eight-character SHAs and fingerprints.
- JSON's `completion` embeds the immutable completion fields and expands
  `handoff`, `merge_observation`, `ci_results`, `rounds`, `approval` and `missing`.
  Timestamps, SHAs and recorded CI URLs remain full. It is null outside done
  tasks, including epics, subtasks and cancellation. No declared loop-contract field changes.
- Missing evidence rows remain explicit while other facts render. A missing
  completion row produces an object with empty expansions and a missing-record
  explanation. Missing fingerprints and reviewer models are explained too.
  Current GitHub state and historical git objects are never needed.

**Builder's implementation choices:** one read transaction expands the
completion. Existing exported readers load completion, rounds and approvals;
small typed JSON queries read the exact hand-off, merge and CI IDs, because
those packages expose latest/by-task readers rather than exact-ID readers.
The JSON key `merge_observation` distinguishes the immutable observation
from the existing `merge` pending-review state. Only IDs named by completion
are expanded, in recorded order; later records
cannot replace them. Top-level hand-off, floor and approval fields also use
that historical evidence; the approval is historically current, and the
recorded disclosure is reused. Done-task text moves Hand-off and Approval
into Completed. Missing explanations are an array of text strings alongside
the typed evidence IDs. No dependency, migration, Event, record-package change,
network operation or new git read was added; the Assured path floor does not
apply. Tests include text/JSON goldens, public-command completion fixtures,
missing and competing records, and a built-binary journey with fake `gh`.

**Left for later:** None discovered in this task.

**D5 / F4 Threat model: honest but fallible agents** on a single-user
machine. Not defending against deliberately adversarial agents. A deceptive
agent could forge a local approval; that is accepted. Friction: a hook denies
agents `fid approve`, `fid waive` and `gh pr merge` (amended 11 November
2026, item 7: it named `fid admin` too, and there's no `fid admin`); fid
refuses approval without an interactive person and refuses `agent:*`
actors; the approval prompt shows the diffstat, checks and review verdicts.
atlas is out of scope for now.

- **Trust by actor type** (amended 10 November 2026, ADR-0014). Any
  well-formed actor is trusted by its type, built into fid; `.fid.toml` has
  no actor list. `agent:*` gets the work set; `human:*` gets the work set
  plus a person's decisions; `service:*` gets nothing for now.
- **A person's decisions** (amended 10 November 2026, ADR-0014). fid refuses
  `agent:*` actors for every kind: approving, waiving a lane, taking over,
  allowing more rework rounds, closing a finding, accepting discovered work,
  approving an epic, answering a task parked while waiting for Mark (added
  11 November 2026, item 5) and closing an epic (settled 12 November 2026,
  following item 9). One record holds them all. An answer binds to the task
  and the question it answers (settled 12 November 2026, following item 5).

**D10 Approvals go stale** (amended). Reviews and approvals bind to the
change's fingerprint (`git patch-id` of the branch diff); checks bind to the
tree, and CI re-runs them on the final commit. (Amended 2026-11-22, Mark: since `fid check`
runs the lane's checks itself and binds each pass to the tree, a reviewer of
that tree relies on those recorded passes and doesn't rerun the full suites;
it tests adequacy instead, by reading the change, running focused tests and
probing with mutations. The passes stand until the tree changes. The kit's
reviewer (3.3) carries the rule.) A clean rebase keeps an
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

- **What a check result binds to** (amended 10 November 2026, ADR-0014). A
  local result binds to the tree it ran on and to its check's own definition,
  and counts only on a clean tree that didn't change during the check. The
  hand-off needs passing results on the head commit's tree and records the
  policy in force. `done` needs passing local results on the merged PR head's
  tree for every check the hand-off policy required, plus CI green on the
  merge commit, recorded as a check result whose source is CI. Only `done`
  completes a task: `fid issue transition` refuses to (1.6).
- **What `done` requires in each lane, and how an epic closes** (decided 11
  November 2026, items 8 and 9, for 1.6). With no CI run on the merge
  commit, `done` waits ("waiting for CI"). It asks `gh` and never runs
  `git fetch`. A person closes an epic with `fid done <epic>`, refused while
  any of its tasks is open. (Settled 12 November 2026, following item 9:
  `fid issue transition` refuses to complete an epic too, and
  `fid done <epic>` records the close as a person's decision.)
- **Details** (decided 7 November 2026, for task 0.7):
  - Why: CI used about 1,650 billed minutes from 1 to 7 November, against
    the 2,000 a month GitHub Free includes across Mark's private repos.
    That's 40 runs at about 41 minutes each: six jobs per run, each rounded
    up to a whole minute. The race detector was most of it. `cmd/fid` takes
    320 s with `-race` and 86 s without, and the heavy `internal/` tests run
    8–10× slower under race. Yet non-test product code starts goroutines in
    one place, `internal/evidence/live.go`, and the race detector can't see
    races between fid processes.
  - **G:** CI runs the race detector only in packages whose non-test
    product code starts goroutines (today `internal/evidence`), and runs
    without it everywhere else. CI finds those packages by parsing the
    source on every run, so a new goroutine can't escape race coverage. It
    refuses a package that doesn't build in its mode, and tests that build
    only in the other mode, because those would never run.
  - **H:** a push to `main` may reuse a passed CI run when the merge
    commit's tree is identical to the tree that run tested. Done (merged
    into `main` with CI green on that commit) accepts that reuse. This
    follows "checks bind to the tree". A pull request run that passes
    records its tree; the run on `main` looks for that record and, if a
    successful run left it, names that run instead of testing again. Without
    one it runs everything.
  - Also accepted: CI runs only on pull requests that are ready for review,
    not on drafts; a new push cancels a pull request's superseded run (never
    a run on `main`); shards are balanced by measured time, and there are
    fewer of them.
  - G applies to local checks too: `make test` and the integration command
    in AGENTS.md, so local runs match CI.
  - One test shard. Minutes are scarce and CI now runs about once per pull
    request, so a couple of extra minutes of wall time don't matter. The
    balancing by measured time stays, for when more shards are worth it.
  - `make fmt-check` fails when gofmt lists a file, so CI's format check
    blocks. Before, `gofmt -l` listed files and still exited 0.

## Work model

**Merge authority: structured, earned autonomy** (amended).

| Lane | Needed before merge | Merged by |
|---|---|---|
| Routine | Local checks; test-adequacy review if tests changed | Mark, until `auto_merge` is enabled for the lane |
| Standard | Plus a code review (other model) with no open major finding; the judge's call is recorded in shadow, never required (amended 11 November 2026, item 8) | Mark, until `auto_merge` is enabled for the lane |
| Assured | Plus a separate review on the other model by default (amended 11 November 2026), Mark's approval and a rollback note: a `## Rollback` section in the PR description, which `done` checks for (item 8) | Mark |

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
  adequacy, security). Assured work adds a separate review on the other model
  (amended 11 November 2026: by default; the fallback below covers it too).
- **Same-model review as a fallback** (decided 7 November 2026; amended 10
  and 11 November 2026). The other model reviews by default. When it isn't
  available, because Mark has no subscription to it or its usage is limited,
  a fresh session of the same model with clean context reviews instead. This
  may happen often, so it's a standing fallback, not a waiver for one PR. The
  review comment names the model that reviewed. fid's independence rule
  compares actors, not models, so the reviewer needs a different actor ID
  from the builder (ADR-0014).
  - On 11 November 2026, Mark extended the fallback to Assured work's
    separate review. That review still runs, as another fresh session,
    which is the same model when the other isn't available.
  - The reset's own Assured PRs keep one review per PR, as AGENTS.md step 7
    says, until 1.6's `done` enforces the separate Assured review on fid's
    own work (decided 11 November 2026, item 2).
  - fid records which model built the change and which model reviewed each
    round (2.1). When no review of an Assured change came from a model other
    than the builder's, the approval view (2.3) and `fid show` say "no
    other-model review", so Mark's approval covers it knowingly.
- **Findings and rounds** (amended 10 November 2026, ADR-0014). Findings are
  major (blocks merge), minor (fix or answer with a reason) or nit
  (optional), as in the kit's reviewer. fid refuses a review round from an
  actor that held the task to build it. Only a later round, or a person with
  a reason, closes a major finding.
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
  `fid waive`. (Amended 10 November 2026, ADR-0014: 2.2b builds `fid waive`
  with the lane floor. A waiver binds to the change's fingerprint, so an edit
  returns the task to its floor. Hard floors, such as `.fid.toml` and CI
  workflows, can't be waived.)
  (Amended 14 November 2026, 1.5a4: a person may correct a todo task's
  planning lane downward with `issue update --lane`. After starting,
  lowering still requires `waive`.)
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
  When an author defers a minor review finding, it is filed as a discovered
  issue under the same rule (task 2.1; added 7 November 2026). fid later
  builds the PR's "Left over" section from `discovered_from` links, just as
  it builds the review comment from its record; that rendering is on the
  later list, not needed for the trial.
- `fid next` order: resume interrupted (same worktree); rework; tasks that
  unblock approved work; unstarted tasks by epic order and priority. It skips
  claimed, blocked, waiting-for-Mark, untriaged and unapproved-epic tasks,
  each with a typed reason for `--explain`. (Amended 11 November 2026: a
  task waiting for Mark is skipped as `waiting-for-you`, item 5, and `next`
  also hands out review work, as kind `review`, item 11.) (Decided 12
  November 2026: review work comes after rework and before tasks that
  unblock approved work, item 13, and a worktree whose task waits for Mark
  gets `none`, item 12.)

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

1.8 imported this list into fid on 2026-11-21, as 43 tasks in the epic
"Later list (after the trial)", which waits for 3.2. Items marked
"(planned: 2.2b)" or "(planned: 3.2)" went onto that task instead, as
acceptance criteria, except tidying `go.mod`, which 3.2's row already
covers.

New items (Mark, 2026-11-21): until 3.2, an item goes in this list
first, as before, and the orchestrator also files it as a task in fid's
Later epic, as an agent. Once 2.1 adds acceptance of agent-created work,
such a task waits for Mark to accept it; until then the epic's wait for
3.2 holds it back.

- From this PR (2.4): an extra integration-tagged lint run
  (`golangci-lint run --build-tags integration ./cmd/fid/integration/...`)
  found nine issues, all in files unchanged by 2.4: an unchecked
  `conn.Close` in `projection_helpers_test.go`, import ordering in
  `rework_test.go`, and seven unused helpers (`p21Missing`, `p23Number`,
  `p23NumberValue`, `p23Clone`, `p22Git`, `p22GitOutput`,
  `closureCompleteIssue`). Tidy these in test maintenance and consider
  checking integration-tagged code in lint; the normal make lint command
  does not select the integration build tag.

- From #101 (1.9c): some older direct retry readers (criterion, transition,
  labels, links, Artifact verification, and hold operations other than
  start/hand-off) still decode
  an already-current projection with plain `encoding/json`. Shared admission
  now validates their source Events, including every embedded-JSON map;
  separate detection of duplicate keys in manually corrupted projection
  rows is outside this fix. Decide whether those readers should use the same
  defensive scan as the create/update/start/hand-off readers. This is a code
  inspection observation, not an observed writer/admission bypass.

- From #101 (1.9c): profile the journal test suite before more matrix growth.
  Round 6 consolidated the accepted-ledger round trip, bounded import workers
  and reused connections where the database file stays in place, keeping every
  refusal, unchanged-target check and recovery import. Its final `make test`
  passed with journal at 529.577s and export at 193.163s. Journal still misses
  the requested under-300s aim and has about 70s before Go's unchanged 600s
  default. Further headroom needs profiling; no timeout was raised.
  Became RESET 1.9d (Mark, 2026-11-23), after GitHub CI timed the package
  out at 20 minutes on #103's head.
  Resolved by 1.9d, 2026-11-23: the full embedded-JSON matrix now runs at
  the shared validators, and row identity retains the complete corpus at
  one real admission boundary with schema/table/reason representatives on
  every path. Both splits keep their refusal, state and recovery assertions
  and are guarded by source-overlay mutations.

- From this PR (1.9d): the embedded-JSON and row-identity splits reduce
  standalone journal from 434.396s before the first pass to 292.615s, but
  `make test` journal is 351.962s, above the 300s aim. Its
  refusal, state-preservation and recovery assertions stay intact; no timeout
  or check changes. The next aggregate repository cost is
  `TestForgedStateMovesRefusedOnEveryJournalPath` (about 204s of overlapping
  leaf elapsed time, not CPU or package wall time). Its typed refusals differ
  by family, subject, from/to state and legacy status, so a table/reason split
  cannot simply reuse row identity's equivalence classes. A further cut needs
  an explicit selection contract preserving those state-move distinctions;
  until then retain its complete path matrix. The 22.82s family-validation
  test also stays: each of its 18 mutations covers a distinct rule, so keeping
  a carrier per rule on live projection and replay retains all its cases.

- From #101 (1.9c): Event-v1 history still has its historical scalar and
  revision grammar. Its Issue identities are immutable and its fixed-state
  moves are refused, but it can set a revision without Event-v2's guarded
  prior-plus-one row. No real ledger holds Event-v1. Decide whether to retire
  its remaining grammar or harden those scalar updates before supporting
  such history as operational data.
- From Mark (2026-11-23): move a task to another epic. Since 1.9c an
  Issue's parent can't change on the general update path, so moving a
  Later item into phase 2 (1.9d) meant filing it again and cancelling the
  old one. A dedicated command and Event, admitted on every journal path:
  the target is an open epic, done or cancelled tasks don't move, epic
  sequences and blocks links stay valid, actor and reason recorded.
- From Mark (2026-11-23): accept every waiting task in one command (for
  example `fid accept --all`, or by epic), listing what it accepts and
  recording one decision each; agents stay refused. `fid accept` takes
  explicit IDs only.
- From this PR (2.1b): proposed RESET 2.1c, the remaining review-delivery
  scope after the brief's size split: builder/person answers to open findings
  without closing them; atomic minor/nit deferral into a discovered task, link,
  acceptance criterion and closing disposition; canonical PR comments for
  rounds, answers and deferrals, recorded URLs, retry protection, pending-post
  warnings and `fid review post`. Its atomic deferral Event must also derive
  the new task’s creator from that Event’s actor, keeping agent deferrals
  unaccepted. Keep the original design and Assured lane.
  Resolved by “Decided for 2.1c, 2026-11-23”, above.
- From this PR (2.1b): D6's exception serving discovered work that blocks a
  task in an approved epic needs 3.1's epic approval; add it there. Until then,
  every agent-created task and subtask requires a person's acceptance.
- From this PR (2.1b): builder-chosen rework after an approving round (for
  example voluntarily fixing a minor) still needs a product decision. The
  current commands support rework following `CHANGES`; do not invent a
  state transition for approved work in 2.1c.

- From this PR (2.1): investigate takeover of an in-review task. A temporary
  test overlay running public commands confirmed that takeover releases the
  hand-off's kept reservation without returning the task to `in_progress`,
  then a changes round and another hand-off both refuse. Code inspection
  also found done requires that released reservation. Settle the intended
  takeover behavior; existing takeover/state rules are outside 2.1's scope.

- From #91 (golden-wrap): fid's text wrapper breaks lines at spaces inside
  paths it prints (a workspace or config path containing spaces wraps
  mid-path). Decide in later display work whether paths should be kept whole;
  the text-golden fixtures require a whitespace-free TMPDIR meanwhile.

- From #89 (1.6k): status's existing READY epic row says "no tasks"
  for an epic with only cancelled tasks, because it counts only tasks in
  its displayed section, which omits cancelled children
  (`internal/status/text.go`, `epicCount`). The new epic status golden
  shows this alongside a correct NEEDS YOU close action. Clarify the
  row wording in later display work; eligibility already counts actual task
  presence and correctly permits all-cancelled epics.

- From #68 (1.5b1) and #80 (1.6d): investigate the intermittent write-lock
  refusal in `TestLeanCheckRenewsHoldAndStopsOnTakeover/takeover`. Twice a
  full `make test` failed at `check_ownership_test.go:117` with
  `task.takeover: write lock unavailable`: #68's first run (the subtest took
  62.74s) and #80's first fix-round run. Each time an unchanged focused
  rerun (1.343s in #68; `-count=3` in #80) and a complete rerun passed. The
  cause is not established. Neither PR changed the ownership test or the
  production lock code, and the test hasn't changed since #64 (1.5a5) added
  it.

- From #71 (1.5c2): `fid artifact verify`'s missing-Artifact failure still
  says `allowed: ["accepted gate registration input"]`; 1.5c2 kept it so
  the failure contract didn't change during removal. Reword it when that
  failure is next touched.

- From #71 (1.5c2): `cli.RegisteredCodes` still omits the existing,
  reachable `idempotency-conflict` code. Its established supplemental process
  pair remains; decide separately whether the public enumeration should
  include it. #71 only removed codes whose producers had retired.

- From #68 (1.5b1): Q3 keeps frozen V1 `agent_runs`, `triage_events`
  and `triage_status`; remove them only in separately scoped cleanup.
- From #68 (1.5b1): Q4 keeps dead config keys
  `maintenance.abandoned_attempt_grace_seconds` and `executor_limits`;
  withdraw them in separately scoped config/archive work.

- Build the PR's "Left over" section from `discovered_from` links, the same
  way the review comment is built from fid's record; not needed for the trial.
- From #42 (1.1d4): a worktree at `<common>/fid`
  (`git worktree add .git/fid`) isn't caught; fid run from a git hook with a
  relative `GIT_DIR` and `--project <subdir>` is refused as not-a-repository;
  with git older than 2.31, the refusal doesn't mention the git version.
- From #43 (1.1d5) and the orchestrator's queue for Mark (2026-11-12): the
  specs lag the reset. Three spec lines still call the ledger `.fid/`
  (`01_PRD.md` PRD-OPS-009, `02_TRD.md` §2 and `05_DESIGN_PLAN.md` §6), as
  does the comment on `p24FailureFilesystem` in the integration tests. The
  spec bodies still describe the kernel: Concerns, Activations, Envelopes
  and Gates, and `fid ready` in the flows. Meanwhile ADR-0013 and RESET.md
  win where they conflict (AGENTS.md).
- From #46 (1.2b), `fid backup`:
  - Copies of one ledger share a genesis Event, so the folder check doesn't
    tell them apart. A ledger directory copied by hand, or a ledger restored
    from an older backup while the original carries on, can still replace
    each other's folder. An ancestry check would catch it: the folder's
    terminal Event must be in this ledger's journal.
  - git's network commands have no timeout, so a hung push holds fid's
    clone. Later backups wait 2 minutes and record `lock-unavailable`.
    1.3a's `fid status` makes a run of those visible (rule c of the 8
    November details); the missing timeout remains.
  - `fid export create` in a repository with no ledger reports
    `internal-error` rather than `project-not-found`.
- From #48 (1.2c), `fid restore`:
  - Concurrent restores into the same repository may meet import's
    nonblocking target lock. The loser gets `internal-error` (`import:
    target lock unavailable`), rather than a typed refusal.
  - If reading the genesis Event fails after successful import, the ledger
    is installed but restore reports `internal-error`.
- From #49 (1.2d): the v2 folder retains the original export ZIP's digest but
  reconstructs its bytes from members. V8 proves byte identity on the current
  toolchain. Future export-header or Go compressor changes may change ZIP
  bytes; assess cross-toolchain compatibility before the supported baseline
  and permanent ledger in 3.2. Preserving the digest's meaning doesn't waive
  this compatibility concern (planned: 3.2).
- From #49's review (1.2d), accepted by Mark on 8 November 2026: resolve
  backup storage for Artifacts larger than the target's per-file limit
  before permanent backups in 3.2. Exports accept up to 1 GiB per Artifact,
  stored as one raw blob in v2. Journal chunking does not cover Artifacts;
  no smaller Artifact ceiling or chunking design has been chosen
  (planned: 3.2).
- From #50 and #51 (1.3), `fid status` and `fid show`. 1.3b took #50's
  long lines, title cuts and kernel-word check, and 1.7 took the
  kernel-word pattern that `internal/status`'s unit tests copied (Resolved,
  below). A `$PAGER` with arguments that names a missing program loses the
  text to the shell's "not found". git's does too, so it stays. These
  remain:
  - `go.mod` marks `golang.org/x/sys` indirect although `internal/backup`,
    `internal/display` and `internal/export` import it. Mark placed tidying
    it with 3.2 on 2026-11-09 (planned: 3.2).
  - Wide characters (CJK, most emoji) count as one column. They misalign
    the columns in `status` and `show`, and make the pager undercount a
    line's rows. Fixing it needs East Asian width tables, a new
    dependency: `go-runewidth` (MIT) or `golang.org/x/text/width` (BSD).
    Do it when a real title has wide characters.
  - Terminal safety of human output. `fid init`'s message prints a root
    path holding control characters raw, and other commands' Human output
    that echoes a path may do the same. `cli.EscapeControlText` escapes
    control characters (Cc) only, so format characters (Cf), such as
    U+202E RIGHT-TO-LEFT OVERRIDE, pass through, in failure output as in
    `status` and `show`. Both need a deliberately odd path or title.
  - `fid status` reads the whole backup record each time. It grows by one
    short line per backup, so this is fine for the trial.
  - `internal/status`'s `row` and `internal/show`'s `reference` lay out the
    same reference line, differing in how they colour the mark. Share them
    when a third renderer needs one.
  - `TestStatusAndShowPageOnARealTerminal` compares `status` output from
    separate processes, including the date in the box's top border, so a
    run that straddles local midnight fails (Codex, #51 round 2, minor;
    answered, not fixed). The window is about half a second. #81 (1.5f)
    fixed the UTC-midnight journal-file failures but left this display
    comparison unchanged, and repeated the follow-up. If it ever fails CI,
    fix it in test code by checking the header date separately and
    normalising it before comparing.
- Genuine cuts in product Go, from the read-only assessment of 9 November
  2026 (the line budget, in the 1.3 details above). #53 removed
  `pricing.toml` and the dead V1 leftovers, about 235 lines. Helpers that
  only tests call (about 120–150 lines then) wait for the dead-code cleanup
  further down this list, since removing them means rewriting tests. The
  assessment's examples, `evidence.Execute` and
  `workflow.EvaluateObligations`, went with the kernel's packages in #72
  (1.5c3); the `internal/project` helpers from #72, below, remain.
- Review rounds, as data for refining the process after the reset: #41, #44
  and #46 took 4, 4 and 5 Codex rounds. In #46, three of the four rounds
  that asked for changes were about code the previous fix round had added
  or removed. AGENTS.md step 7 now limits what a fix round changes (8
  November 2026); see whether round counts drop. Since then, #48, #49, #50,
  #51 and #52 have each taken 2 rounds, and #53 took 1.
- Mark's idea: a kit skill for backlog triage and grooming. An agent proposes
  close, keep or merge for stale or obsolete issues, and Mark accepts in a
  batch; fid supplies signals such as age and whether the paths an issue
  names still exist. Trigger: after the trial, or when untriaged items turn
  the NEEDS YOU line into noise.
- Automatic autonomy dial: per-lane agreement counters, trip rules, and the
  remaining slop sensors.
- Kit install three-way merge (`git merge-file`) and `--check`, once a local
  edit actually gets lost.
- Delta view for re-approval.
- A kit skill that helps build a project-specific loop, after two or three
  real loops exist.
- Bot identity for agents with required review from Mark on Assured paths.
- Mark's GitHub settings, from RESET.md's friction log (2026-10-31) and the
  orchestrator's queue for Mark (2026-11-14): branch protection on `main`
  for fid and kauri, before any lane gets `auto_merge` (needs GitHub Pro or
  a public repo). fid's repository has been public since 2026-11-14, and a
  `main` ruleset was made that day (pull request required, no force push
  or deletion). On 2026-11-16 GitHub's API reported that ruleset's
  enforcement as disabled and no rule active on `main`. Required status
  checks (such as `checks`, `tests (0)` and `quality`) aren't set yet.
- Cloud-session support for the ledger (read-only, then single-writer).
- Declared paths per task, if shared-file conflicts become a regular cost.
- Collapse candidates after the trial: any concept the leak test keeps
  catching. (Amended 10 November 2026: the Activation goes with the kernel
  in 1.5b–1.5d, ADR-0014.)
- Generic retry-exhausted, dead-code cleanup. (Amended 10 November 2026: the
  Envelope revision-1 limit is moot, since Envelopes go with the kernel.)
- Restoring across a fid upgrade (found 10 November 2026; decided 11
  November 2026, item 10, for 3.2). Import and restore need the archive's
  projection schema to equal the running fid's
  (`internal/export/archive_validation.go`), so after an upgrade the newest
  backup can't be restored by the new fid until the next backup runs. Raise
  it with 3.2, before the permanent ledger; the restore drill should cover
  it. Mark decided that restore ignores the archive's projection and
  rebuilds it from the archive's journal with the running fid, so any older
  archive restores, and refuses an archive from a newer fid (planned: 3.2).
- A key to opt out of `done` waiting for CI on the merge commit, for a
  project whose merge commits get no CI run. Not until a project needs one
  (Mark, 11 November 2026, item 8).
- Replay reports `line: 0` for a current-event prior conflict (found in
  1.1d2, #39). Fix it after the trial.
- yoyu migration (after the trial).
- CI test time, from #36 (1.1c1), 4 November 2026. #44 (0.7) dropped the
  race detector outside goroutine packages (G), cut the shards from four
  to one and balances them by measured time from
  `.github/scripts/durations.json`; `ci_tests.py durations` rewrites that
  file from a run's shard logs. Restore tests from #48 (1.2c) are absent
  from that file; harmless with one shard, include them when timings are
  refreshed. Two things remain. First, review whether every test earns its
  run time. Measured without race on 2 CPUs before the kernel's removal,
  `cmd/fid` was the long pole at 231 s, and four corpora took 309 s of its
  523 s of test time: the grammar built-binary corpus 94 s, and the
  descriptor variant, judgment validation and review-and-finding corpora
  66–80 s each. #71 (1.5c2) deleted the judgment and review-and-finding
  corpora with their commands, so measure again first. Second, fix the
  tests that depend on their environment. About 33 failed when `TMPDIR` is
  behind a symlink (12 in `internal/`, 11 in `cmd/fid`, 10 in
  integration), on `main` too, counted before the kernel's removal. They
  compare paths without resolving symlinks; P24 observation tests, for
  example, hash an unresolved worktree path. Found in 1.1c2. Do the rest
  when a shard's test step passes 10 minutes, and all of it before the
  trial at the latest. That trigger has fired: on 2026-11-15 the one
  shard's test step took 430–743 s on GitHub, over 10 minutes in 6 of the
  10 runs that passed.
- Stale ledger staging, from #40 (1.1d3), 6 November 2026. A crash during
  `init` or `export import` can leave `<common>/fid-staging-*` behind.
  `git status` doesn't show it, and nothing removes it. Add cleanup once
  one turns up; it needs a rule for when staging is safe to delete.
- From #67 (replay-race, 2026-11-14):
  - Replay reports any error while projecting an Event as a malformed
    journal entry (`journal-replay-stopped`, integrity, not retryable),
    naming a line that may be fine. Only a line that doesn't parse should
    get that label; a projection error needs its own.
  - A read that inspects the journal without the write lock could still
    meet a half-written last line: `AppendEvent` writes each Event with one
    `Write` call, but Go may retry a short write, and the inspection parses
    a non-empty fragment at the end of the file at once. Not seen in about
    3,600 concurrent reads; telling a live append from a damaged tail needs
    a design, not a one-line fix.
- From #72 (1.5c3): deleting the kernel's domains left these `internal/project`
  helpers with no production caller, only their own tests:
  `ObserveRecordedRepository` and `ResolveSourceInput` (`source.go`),
  `ExpectedRepository` and `MatchesRecordedRepository`
  (`expected_repository.go`) and `IsAncestor` (`repository.go`). 1.6a–1.6l
  didn't need them, and they still have no production caller.
  Delete them, with their tests, unless 2.1 turns out to need them.
- CHANGELOG: it hasn't been kept up during the reset. Catch it up in 3.2,
  with the v0.2.0 release (planned: 3.2).
- From CI on #85 (1.6h), 2026-11-15, a watch item:
  `TestTwoWorktreesShareOneLedger` failed once on GitHub's runner (git
  2.56.1). `git worktree add` in `gittest.AddWorktree` stopped with `could
  not open '.git/worktrees/linked/locked' for writing: No such file or
  directory`. The rerun passed, as did 20 of 20 local runs on git 2.53.0
  and, on 2026-11-16, 20 of 20 runs of the test alone on git 2.56.1. No CI
  run since has needed a rerun for it. If it recurs, investigate creating
  worktrees in parallel under newer git.
- From #56 (1.4a), round 1: `fid next --claim --actor service:ci` printed
  `operation=next.claim` in its human refusal; #56 fixed that refusal. On
  2026-11-12 the orchestrator routed the general case, operation names in
  any human refusal, to 1.7's leak test, but #75 (1.7) checks kernel nouns
  and IDs only. That catches `next.claim` through "claim", but not names
  such as `task.takeover`. Two refusals probed on 2026-11-16 (`next
  --claim` by a service actor, `takeover` by an agent) named no operation.
- Mark's GitHub settings, from the orchestrator's queue for Mark
  (2026-11-14): pin the workflows' actions to full-length commit SHAs (they
  name tags, such as `actions/checkout@v7`; an Assured change), then turn
  on the repository's "Require actions to be pinned to a full-length commit
  SHA".
- Mark to decide, from RESET.md's friction log (2026-11-14, 1.5a) and the
  orchestrator's queue for Mark: collapse golden files in GitHub's diffs,
  for example with a `.gitattributes` that marks test data
  `linguist-generated`, and add a "What to read" section to PR
  descriptions. Golden files were about 35–70% of the lines each 1.5a PR
  added.
- From Mark, 2026-11-17: a label scope for `next` (`fid next --label
  <label>`), like 2.4's `--epic` scope but across epics, for a feature
  whose tasks span epics. Add it when a real feature needs it.
- From #93 (1.9): a hand-off of a very wide change (about 5,000 or more
  long paths) is now refused before append, rather than stopping every
  later write. Lifting the limit means storing changed paths outside the
  Event (for example a digest and an Artifact). Mark's call when a real
  change hits it.
- From #93 (1.9): a caller-supplied Event ID without a date (not
  `<prefix>-YYYYMMDD-<n>`) isn't checked for uniqueness on live append; no
  writer supplies one, and replay and export refuse a repeated ID. A dated
  one must take its date's next sequence (review round 2), so it can't
  repeat. The catalog's `issue.label-added` and
  `issue.label-removed` have no writer (`fid issue label` writes
  `issue.label.added` and `issue.label.removed`); remove them in the 3.2
  baseline if nothing needs them.
- From #93 (1.9): catch-up projects pending Events one transaction at a
  time. Rules that need only the journal (shape, size, repeated Event IDs,
  aggregate sequences) refuse in replay's pre-scan, before any Event
  projects. A refusal that needs projection state, from a family's
  validator, still leaves the valid Events before it projected, with the
  checkpoint on the last of them: a consistent prefix, and the next read
  refuses at the same Event. All-or-nothing catch-up would need one
  transaction over the pending Events; nothing needs it yet.
- From #93 (1.9): admission reads a dated Event ID's sequence by value and
  doesn't tie the ID's date to the Event's time or file. So it admits IDs
  no writer writes, such as `EVT-20261110-4` for the writers'
  `EVT-20261110-000004`, or the next ID of another date. None of them
  harms allocation, because each must still take its date's next sequence.
  Tighten to the writers' exact form only if something comes to rely on it.
- From 1.8's policy PR: local `make test` runs `go test` without
  `-count=1`. `internal/export`'s tests build fid in a subprocess, which
  Go's test cache doesn't track, so under `fid check` a cached pass of
  that package can hide a broken binary. `cmd/fid`'s tests, in the same
  run, would usually catch it, and CI passes `-count=1`. The fix,
  `-count=1` in `ci_tests.py`'s local mode, changes a check's definition,
  so it's Assured.
- From 1.8 (#96): two display problems in the friction log. `fid status`
  lists tasks whose epic is blocked under READY, and `fid next --explain`
  lists every open task, with blocked tasks that others wait on above the
  ready ones. Fix them in later display work.


**Resolved.** Later work or decisions resolved these, so 1.8 doesn't
import them. Dates are merge or decision dates, in New Zealand time.

- From this PR (1.9b): a person can cancel a finished task or a closed
  epic. `fid issue transition --to cancelled` allows `cancelled` from
  `done`, and admission keeps that move because the command makes it.
  Decide whether cancelling finished work should be refused.
  Resolved in this PR (1.9c), 2026-11-22: The shared move table and the command refuse both moves, with the new-task/rollback explanation.

- From this PR (1.9b): `issues.workflow_id` and `parent_id` can still
  change through `issue.updated`. Neither moves a state, since admission
  keys fixed moves on the immutable `kind`, but moving an open task out of
  its epic would let the epic close. Belongs with 1.9c's immutable columns.
  Resolved in this PR (1.9c), 2026-11-22: Shared admission makes both columns immutable, including on Event-v1.

- From this PR (1.9b): the hold-lifecycle rule still lets
  `issue.transitioned` end a hold as `task-completed`, though no command
  can now complete a task through it; only `task-cancelled` is reachable.
  Tidy in 1.9c.
  Resolved in this PR (1.9c), 2026-11-22: Shared admission takes only task-cancelled from issue.transitioned.

- From #97 (1.9b): Event-v1 history can still rename an Issue (set
  `issues.id`) or delete one and insert another under its id, on replay.
  Since #97's first fix round admission checks each state change against
  the row it lands on, so neither can redirect a move, but the identity
  change itself is admitted. 1.9c's "an Event can't rename its Issue"
  should cover Event-v1 too, or record why it doesn't.
  Resolved in this PR (1.9c), 2026-11-22: Shared admission refuses Issue renames and deletion/recreation on Event-v1 too.

- From #97 (1.9b): a refused subtask move names the wrong owner. When a
  family other than `issue.transitioned` moves a subtask (an
  `issue.updated` setting its `status`, or an Event-v1 `issue.created`
  that creates and moves it), `StateMoveError` says "only a person's
  cancellation makes that move". That is `issue.transitioned`'s owner for
  a task or an epic; a subtask's moves are `fid issue transition`'s, by
  anyone. Name the owner by kind (`issuestate.OwnerCommand`).
  Resolved in this PR (1.9c), 2026-11-22: OwnerCommand names fid issue transition for a subtask.

- From #100 (2.1): proposed row 2.1b, review delivery. Mark accepted it on
  2026-11-22, with his extra-rework-rounds decision folded in.
- From #98 (2.2b): proposed row 2.2c, `fid waive`, the remaining 2.2b
  scope. Mark accepted it on 2026-11-22 ("Decided 2026-11-22", item 1).
- From 1.8's policy PR: Mark chose fid's Assured trust-boundary path list
  on 2026-11-21; "Decided for 2.2b, 2026-11-21" records its exact list,
  definition and exclusions. 2.2b computes the floor and `.fid.toml` now
  lists those paths. Releases remain 3.2 work when their files exist.

- Task 1.8 brings this later list into fid along with the remaining reset
  tasks (decided 2026-11-07). Done 2026-11-21; see the note at the top of
  the list.
- From #93 (1.9): proposed rows 1.9b and 1.9c, with the two Issue-table
  gaps #93's first fix round found for 1.9c. They became RESET rows 1.9b
  and 1.9c on 2026-11-17; "Decided 2026-11-17", item 1, keeps their scope.
- From Mark's question, 2026-11-15: how jig's markdown review becomes
  `fid review record`'s JSON. Settled 2026-11-17 ("Decided 2026-11-17",
  item 3).

- From #80 (1.6d): audit older writers for Events at or above 1 MiB, which
  the journal's aggregate and sequence scanners can't read. Resolved by
  #93 (1.9), 2026-11-17: every family's Event is checked against the line
  budget on live append, replay and import, and replay's pre-scan also
  measures each stored line.
- From #78 (1.6b): the broader retained-table ownership audit, including
  the Event-v1 early return in `validateCurrentEvent`. Resolved by #93
  (1.9), 2026-11-17: Event-v1 may change only the V1 schema's tables, and
  each Event-v2 family only those of its catalog grammar, on every path.
  Ownership of `issues.status` is RESET 1.9b.
- From #73 (1.5d1): retained-table ownership and column guards against
  SQLite case aliases. Resolved by #93 (1.9), 2026-11-17: admission
  refuses any table or column name that isn't lowercase, on every path.
- From #76 (1.5e), #84 (1.6g) and #90 (1.6l): text goldens wrapped real
  temp paths before masking them, so their line breaks varied with TMPDIR's
  length. Resolved by #91 (golden-wrap), 2026-11-16: fixture paths have a
  fixed display width.
- From #80 (1.6d): V8's public-command stale round. Resolved by #83 (1.6f),
  2026-11-15: V8 runs the public rework loop and checks stale first rounds
  after backup, wipe and restore.
- From #68 (1.5b1), Q7: whether hand-off and `done` need a passing check's
  kept Artifacts still available. Resolved by #85 (1.6h), 2026-11-15: both
  refuse with `kept-artifact-unavailable` (Decided 2026-11-15, item 8), and
  recorded results stay immutable.
- From #50 and #51 (1.3): `internal/status`'s unit tests copied the
  kernel-word pattern. Resolved by #75 (1.7), 2026-11-15: they use the
  shared `internal/leaktest`.
- From the 2026-11-09 assessment: remove legacy manifest-v1
  validation. Resolved by #73 (1.5d1), 2026-11-15: validation and import
  both refuse v1.
- From the 2026-11-09 assessment: the kernel's test-only helpers.
  Resolved by #72 (1.5c3), 2026-11-15, which deleted the kernel's domain
  packages.
- From #48 (1.2c): the export-manifest project ID check, scratch placement
  on disk, private archive permissions and the relative-`TMPDIR` finding;
  from #46 (1.2b), compressed-archive storage. Resolved by #49 (1.2d),
  2026-11-08.
- From #36 (1.1c1): balance CI shards by measured time. Resolved by #44
  (0.7), 2026-11-07.
- From #44's round-3 review (0.7): the UTC-midnight failure (`journal files
  = 2, want 1`). Resolved by #81 (1.5f), 2026-11-15: test assertions read
  every dated journal file, with forced-crossing regressions; no product
  clock change was needed.

## Next steps

1. ~~Mark reviews and accepts ADR-0013 and approves `RESET.md`.~~ Done 31 October 2026.
2. Week 0 tasks from `RESET.md`.
