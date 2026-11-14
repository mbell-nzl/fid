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
two-worktree test. Mark expects mostly one sequential loop.

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

- Task 1.8 brings this later list into fid along with the remaining reset
  tasks (decided 7 November 2026).
- Build the PR's "Left over" section from `discovered_from` links, the same
  way the review comment is built from fid's record; not needed for the trial.
- From #42 (1.1d4): a worktree at `<common>/fid`
  (`git worktree add .git/fid`) isn't caught; fid run from a git hook with a
  relative `GIT_DIR` and `--project <subdir>` is refused as not-a-repository;
  with git older than 2.31, the refusal doesn't mention the git version.
- From #43 (1.1d5): older specs and test comments still call the ledger `.fid/`.
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
  - The export-manifest project ID check, disk placement of scratch, private
    archive permissions and relative-`TMPDIR` finding were handled in 1.2d
    (decisions above). Unzipped backup storage and cache packing also resolve
    #46's compressed-archive storage follow-up.
- From 1.2d: the v2 folder retains the original export ZIP's digest but
  reconstructs its bytes from members. V8 proves byte identity on the current
  toolchain. Future export-header or Go compressor changes may change ZIP
  bytes; assess cross-toolchain compatibility before the supported baseline
  and permanent ledger in 3.2. Preserving the digest's meaning doesn't waive
  this compatibility concern.
- From 1.2d review, accepted by Mark on 8 November 2026: resolve backup
  storage for Artifacts larger than the target's per-file limit before
  permanent backups in 3.2. Exports accept up to 1 GiB per Artifact, stored
  as one raw blob in v2. Journal chunking does not cover Artifacts; no
  smaller Artifact ceiling or chunking design has been chosen.
- From #50 and #51 (1.3), `fid status` and `fid show`. 1.3b took #50's
  long lines, title cuts and kernel-word check. Mark placed two more on 9
  November 2026: 1.7 replaces the kernel-word pattern that
  `internal/status`'s unit tests copy, and 3.2 tidies `go.mod`, which marks
  `golang.org/x/sys` indirect although `internal/backup` and
  `internal/display` import it. A `$PAGER` with arguments that names a
  missing program loses the text to the shell's "not found". git's does
  too, so it stays. These remain:
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
    answered, not fixed). The window is about half a second. Treat it like
    the UTC-midnight test failure under "CI test time" below: if it ever
    fails CI, fix it in test code by checking the header date separately
    and normalising it before comparing.
- Genuine cuts in product Go, from the read-only assessment of 9 November
  2026 (the line budget, in the 1.3 details above). #53 removed
  `pricing.toml` and the dead V1 leftovers, about 235 lines. These remain:
  - Helpers that only tests call (about 120–150 lines, such as
    `evidence.Execute` and `workflow.EvaluateObligations`) wait for the
    dead-code cleanup further down this list, since removing them means
    rewriting tests.
  - Removing legacy manifest-v1 validation support (about 30 lines) is
    Assured and Mark's call, since it withdraws an archive-format promise:
    `fid export validate` still accepts Phase 21 manifest-v1 archives,
    though `fid export import` already refuses them as
    `schema-unsupported` (`TestValidateRetainsPinnedV1Profile` pins both).
    Decided 10 November 2026 (ADR-0014): removed with the kernel's catalogs
    in 1.5d, together with that test.
  - The kernel's removal (1.5b–1.5d) takes the test-only helpers that belong
    to the kernel.
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
- Branch protection on `main` for fid and kauri, before any lane gets
  `auto_merge` (needs GitHub Pro or a public repo).
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
  archive restores, and refuses an archive from a newer fid.
- A key to opt out of `done` waiting for CI on the merge commit, for a
  project whose merge commits get no CI run. Not until a project needs one
  (Mark, 11 November 2026, item 8).
- Replay reports `line: 0` for a current-event prior conflict (found in
  1.1d2, #39). Fix it after the trial.
- yoyu migration (after the trial).
- CI test time (from 1.1c1, 4 November 2026; partly done in 0.7). 0.7
  dropped the race detector outside goroutine packages (G) and cut the
  shards from four to one. Shards are still balanced by measured time from
  `.github/scripts/durations.json`; `ci_tests.py durations` rewrites that
  file from a run's shard logs. Restore tests from #48 (1.2c) are absent
  from that file; harmless with one shard, include them when timings are
  refreshed. Two things remain. First, review whether
  every test earns its run time. Measured without race on 2 CPUs,
  `cmd/fid` is the long pole at 231 s, and four corpora take 309 s of its
  523 s of test time: the grammar built-binary corpus 94 s, and the
  descriptor variant, judgment validation and review-and-finding corpora
  66–80 s each. Second, fix the tests that depend on their environment.
  About 33 fail when `TMPDIR` is behind a symlink (12 in `internal/`, 11 in
  `cmd/fid`, 10 in integration), on `main` too. They compare paths without
  resolving symlinks; P24 observation tests, for example, hash an
  unresolved worktree path. Found in 1.1c2; no PR since added one. And
  `cmd/fid` tests fail when a run crosses UTC midnight (`journal files = 2,
  want 1`, from `mustJournalPath` at `cmd/fid/init_test.go:710`): the
  fixture creates the journal on one UTC date and the commands append to a
  new one on the next. Found in 0.7's round-3 review; Mark accepted it for
  0.7 on 7 November 2026. Fix it in test code (a fixed date in the
  fixture), not product code, if it ever fails CI. Do the rest when a
  shard's test step passes 10 minutes, and all of it before the trial at
  the latest.
- Stale ledger staging (from 1.1d3, 6 November 2026). A crash during
  `init` or `export import` can leave `<common>/fid-staging-*` behind.
  `git status` doesn't show it, and nothing removes it. Add cleanup once
  one turns up; it needs a rule for when staging is safe to delete.
- CHANGELOG: it hasn't been kept up during the reset. Catch it up in 3.2,
  with the v0.2.0 release.

## Next steps

1. ~~Mark reviews and accepts ADR-0013 and approves `RESET.md`.~~ Done 31 October 2026.
2. Week 0 tasks from `RESET.md`.
