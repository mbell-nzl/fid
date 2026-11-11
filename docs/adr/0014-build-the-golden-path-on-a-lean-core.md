---
status: accepted
date: 2026-11-10
---

# Build the golden path on a lean core and remove the kernel

The reset's golden path (`next --claim`, `start`, `check`, `park`, `done`,
reviews and approvals) is built on a lean core: Issues, the journal and the
check runner, plus a few small records. It doesn't compose the kernel chain.
The kernel chain (Kit Activation, Work Envelope, Obligation, Claim, Run,
Attempt, Evidence Record, Stops, Attention Items, kernel budgets and Run
recovery) is removed right after task 1.5a, before fid tracks its own work
in 1.8.

## Context

ADR-0013 and D9 built the golden path as a facade over the kernel. Designing
task 1.4 showed that the kernel's claim can't be taken without a finalized
Activation and an active Envelope, and that an Envelope needs grants,
obligations and registered gates. Mark asked whether the chain is needed at
all once the facade works.

A proposal (a lean core beside the kernel, the kernel deleted before 3.2) had
two reviews on 10 November 2026:

- Codex (gpt-6-astra, xhigh effort) disagreed. It preferred incremental
  replacement that removes displaced machinery slice by slice.
- A fresh Opus 5.5 session agreed with changes ("B-early"): remove the kernel
  right after the lean check in 1.5a, not before 3.2.

Both took the lean core as the destination. Both found that the proposal
understated what the kernel guarantees: ADR-0007's recovery and liveness
rules, ADR-0008's budgets and rework cap, recording an execution before it
starts, review independence, `done` as the only route to completion, the
waiver record, and Artifacts. Both found it understated the cost of removal:
the journal and export catalogs, the V8 and export fixtures, the data cutover,
and a kernel of about 18–20k lines. No ledger holds real data anywhere, so
nothing needs migrating.

Mark weighed two factors: reach real use of fid soon, and build nothing
temporary, since agents write the code and the right version costs days, not
months. He chose the lean core with the kernel removed right after 1.5a,
then decided each guarantee below, one at a time.

## Decision

- **Order.** This ADR's docs-only PR comes first, and Mark approves it. Then:
  1. 1.4a: `fid next --claim` on a hold record (Assured, V3).
  2. 1.4b: READY in `fid status`, with `next`'s verdict for each task, and
     the `Working` and `Scheduling` rows in `fid show` (Assured).
  3. 2.2a: policy from `origin/<default>`, the whole `.fid.toml` schema,
     trust by actor type and the policy snapshot (Assured, V7).
  4. 1.5a: `start`, `check`, `park` and the heartbeat (Assured).
  5. 1.5b–1.5d: the kernel's removal (Assured).
  6. 1.6, 1.7, 1.8, then 2.1 onwards. 2.2b (lane floor, hard floors and
     `fid waive`) keeps 2.2's old place.
- **The kernel goes; it doesn't move under `fid admin`.** Each slice ends
  with every remaining test passing, none weakened.
  - 1.5b rebuilds V8, the export fixtures and every other test that outlives
    the kernel on lean check results, with their helpers, before any kernel
    code goes. They keep their assertions: V8 still compares the restored
    ledger's Events, files, Artifacts and read-command answers with the
    original's, and the export tests still check Artifacts, import and the
    archive's refusals, now on lean events. Among its read commands, V8
    compares `next` in place of `ready`. Tests of the kernel itself keep
    their kernel fixtures until 1.5c deletes them.
  - 1.5c removes the kernel's commands, its domain code and their tests. It
    also removes `fid ready`, which shows one task by sequence or priority
    and ignores blockers; `fid next` (1.4a) replaces it (Mark, 11 November
    2026).
  - 1.5d removes its journal and export catalog entries, its maintenance
    checks and manifest-v1 validation, and drops its tables.
- **Every new record is complete when it lands.** The PR that adds an event
  family registers it in the journal's and the export's catalogs and adds its
  state to V8's comparison. The removal rebuilds V8 and the export fixtures
  without weakening them, before it deletes anything they use.
- **Plain words.** Human output never uses the kernel's words, and the leak
  test (V2, task 1.7) stays. The lean records are a hold, the hand-off, the
  policy snapshot, a check result and a person's decision. The `--claim` flag
  and JSON codes such as `claim-expired` keep their names, since `--json` is
  exempt from the leak test.

### Trust by actor type

- Any well-formed actor (`agent:`, `human:` or `service:` followed by an ID)
  is trusted by its type. The rule is built into fid. `.fid.toml` holds no
  actor list.
- `agent:*` actors get the work set: taking, starting, checking, parking and
  finishing tasks, and recording reviews.
- `human:*` actors get the work set plus a person's decisions (below).
- `service:*` actors get nothing for now.
- fid can't tell a misnamed actor, which fits the threat model (D5). An actor
  list can be added if misnamed actors ever matter.

### Policy

- fid reads `.fid.toml` only from the local `origin/<default branch>` ref. It
  never fetches the policy remote, so policy is as fresh as the last fetch.
  Backup and restore still fetch their own target. 1.6 decides whether
  `fid done` fetches. ADR-0013's three network uses are unchanged.
- Built-in defaults apply while the file doesn't exist on origin. An invalid
  `.fid.toml` there is refused, with no fallback to the defaults.
- 2.2a defines the whole schema: checks, Assured paths, check-defining paths,
  `auto_merge` per lane, the rework cap and the token budget.
- **The policy snapshot** records the policy fid trusts:
  - its source: `origin/<default>` at commit X, or the built-in defaults;
  - the digest of the effective policy, meaning the built-in rules plus
    `.fid.toml`;
  - the actor rules and the check definitions.
- fid writes a new snapshot only when the effective policy's digest changes.
  It writes it itself, as one Event under the write lock. Writing commands
  that need policy (taking, starting, checking and `done`) check the digest.
  Read-only commands never write. A fid upgrade that changes the built-in
  rules produces a new snapshot.
- The snapshot replaces the Activation and the kit's four-step handover.

### Holds, interruptions and takeover

- `fid next --claim` hands a task out as a hold, with its own ID. A hold binds
  one task, one actor and one worktree. A task has at most one active hold,
  and the hand-out is atomic (V3). The hold's ID goes with every heartbeat,
  check, park and `done`.
- A hold lasts about an hour, and the loop renews it with heartbeats (D8).
  fid refuses a heartbeat from a stale session, so its loop stops loudly.
- **A lapsed hold on a started task** means the task is interrupted.
  - The task stays kept for its worktree. Other worktrees skip it
    (`interrupted-elsewhere`).
  - That worktree's `next --claim` returns it first, as kind `resume`, with a
    new hold linked to the old one.
  - Parking a task as interrupted works the same way.
  - A lapse is never taken as proof that work stopped (ADR-0007).
- **A lapsed hold on an unstarted task** frees the task for anyone.
- **Takeover.** Only a person takes over a task. A takeover ends any hold and
  is recorded. 1.5a names the command.
- fid never touches git state: no reset, clean, stash, checkout or commit.
- A task kept for a worktree that no longer exists is a NEEDS YOU item.
- Parking while waiting for Mark is designed in 1.5a.

### Rework cap and token budget

- Both caps live in fid and are set per lane in `.fid.toml`. ADR-0008's rule
  that fid enforces them, and D2's scenario tests, stand.
- **The rework cap** counts a task's review rounds that asked for changes. It
  defaults to 3. Once the cap is reached, `next` skips the task
  (`rework-cap-reached`) and `fid status` lists it under NEEDS YOU. A person
  then allows more rounds, takes the task over or cancels it.
- **The token budget** sums a task's tokens from 2.6's usage records. It's off
  by default until 2.6 has data. Once a task is over budget, `next` skips it
  (`over-budget`) and it appears under NEEDS YOU. The budget is checked
  between agent sessions, not during one.
- 2.2a parses both settings. 2.1 enforces the rework cap, and 2.4's scenario
  tests cover it. 2.6 enforces the token budget and adds the over-budget
  scenario test with it, since 2.4 runs before 2.6.

### Only `fid done` completes a task

- `fid issue transition` refuses, for every actor, to move a task into a
  state of the completed category, and points at `fid done`.
- There's no manual override. Work that won't be merged is cancelled, and
  only a person cancels.
- Subtasks keep plain transitions. 1.6 decides how an epic closes.
- The refusal lands with `fid done` in 1.6.

### Reviews

- fid refuses a review round from any actor that held the task to build it.
  Actor IDs are free text after the type, so this catches mistakes, not
  deliberate renaming (D5).
- Findings have the kit reviewer's three severities:
  - **major** blocks merge;
  - **minor** means fix it or answer it with a reason;
  - **nit** is optional.

  The kernel's `blocker` severity is dropped.
- The builder may answer a finding ("fixed in a commit", or why not), and may
  defer a minor finding or a nit, which files a discovered issue (D6). Only a
  later review round closes a major finding, as resolved or as not applying.
  A person may also close one, with a reason. `done` refuses while a major
  finding is open.
- A review round records the reviewer's actor ID. Which model reviews is the
  loop's and the kit's rule: the other model by default, or a fresh session
  of the same model when the other isn't available. The review comment names
  the model that reviewed.
- The fallback also covers the separate review that the decisions record's
  work model adds for Assured work (Mark, 11 November 2026). That review
  still runs, as another fresh session, which is the same model when the
  other isn't available.
- fid records which model built the change and which model reviewed each
  round (2.1). When no review of an Assured change came from a model other
  than the builder's, the approval view (2.3) and `fid show` say "no
  other-model review", so Mark's approval covers it knowingly.
- Rounds bind to the change's fingerprint (D10), not to an Envelope revision.
  One session records every required profile.

### A person's decisions

- One record holds every decision a person makes: its kind, its subject (a
  task, a finding or an epic), the reason, and what it binds to.
  - Lowering a lane and approving a change bind to the change's fingerprint,
    so an edit makes them stale (D10).
  - A takeover and extra rework rounds bind to the task.
  - Closing a finding binds to the finding.
  - Approving an epic binds to its validation list (3.1).
- Each kind keeps its own command, such as `fid waive` and `fid approve`.
  fid refuses every kind from `agent:*` actors.
- `fid waive` is built in 2.2b, with the lane floor it lowers. Hard floors
  can't be waived: changes to `.fid.toml` and CI workflows are always Assured.
  If the change is edited after a waiver, the task returns to its floor, and
  `fid status` shows the waiver as stale.
- Each other kind arrives with the task that adds its command: takeover in
  1.5a; extra rounds, closing findings and accepting discovered work (D6) in
  2.1; approval in 2.3; and epic approval in 3.1.
- The record replaces the kernel's Deviations, Human Outcomes and Attention
  Items.

### Check results

- fid runs each check itself. It records the check's start before running it,
  so a start without a result counts for nothing. After the check it looks at
  the tree again. The result keeps the exit code, a digest and excerpt of the
  output, the toolchain and its Artifacts.
- A local result binds to the tree it ran on and to a digest of that check's
  own definition (its command, environment and timeout), not of the whole
  `.fid.toml`. Changing an unrelated setting doesn't make results stale.
- A local result counts only if the tree had nothing uncommitted and hadn't
  changed by the end of the check. Other results are recorded but count for
  nothing.
- **Hand-off.** Before the builder hands a task to review, every check the
  lane requires needs a passing result on the head commit's tree. fid records
  the policy in force at hand-off.
- **`fid done`** needs both:
  - passing local results, on the tree of the merged PR head, for every check
    the hand-off policy required;
  - CI green on the merge commit. `done` records CI's run as a check result
    whose source is CI: the run's ID, the commit and the conclusion.
- A push after the checks means running them again. CI skips the integration
  tests, so for those the local result is the only evidence.
- A `.fid.toml` change merged after hand-off applies from the task's next
  hand-off.

### Old records

- **Clean cut.** 1.5d removes the kernel's event types from the journal and
  export catalogs. A journal or archive that holds them is refused with the
  existing `schema-unsupported`. There is no converter. Import already needs
  the archive's projection schema to match the running fid's.
- 1.5d also stops `fid export validate` accepting manifest-v1 archives, and
  removes the test that pins it. That's a product change Mark approved, not a
  weakened test.
- No ledger holds real data (checked in fid, kauri and yoyu; the backup repo
  holds only its README), so nothing is lost.

## Consequences

- Product Go grows by the lean pieces until 1.5b–1.5d, then shrinks by about
  18–20k lines. Until then it is over `RESET.md`'s line budget.
- 1.4a, 2.2a and 1.5a are Assured, since each adds a migration and event
  families. 2.1 and 2.6 are Assured too, since both will add a migration, for
  review rounds and usage records (Mark, 11 November 2026). 1.4b only changes
  display, but keeps the Assured lane of the 1.4 it was split from, since
  only Mark lowers a lane.
- Tests of the kernel leave `make test` and the integration suite with the
  kernel, which saves scarce CI minutes.
- 3.2's migration squash covers only lean tables.
- The kit's four-step Activation commands go too. Their only outside user is
  the frozen kit V2 branch.
- This ADR doesn't decide whether Concerns, requirement and decision
  references, and `fid capabilities` go with the kernel. 1.5c's PR asks Mark.

## Superseded and amended records

| Record | Clause | Change |
|---|---|---|
| ADR-0013 | "a human CLI built as a facade over the existing kernel"; Facade: "Golden-path commands compose kernel commands" and "Kernel commands move under `fid admin`; any unused at trial end are deleted"; "Interrupted work still resumes through a linked successor Run"; "The facade uses one-Issue Envelopes"; the Activation that its ADR-0011 row keeps | **Amended.** A lean core; the kernel removed after 1.5a; holds; no Envelopes; the policy snapshot replaces the Activation. The `done` evaluator is still the single source of what done needs, and the leak test stays. |
| ADR-0005 | The Next action derived from Envelope, Obligation, Claim, Run, stop and attention state; review as a Review Run that satisfies a review Obligation | **Amended.** Issues stay the only durable work graph. `fid next` derives the next action from Issues, holds, review rounds, check results and a person's decisions. |
| ADR-0006 | A Claim binds the Envelope revision; read-only Reviewer Runs | **Amended.** At most one active hold per task stands. A hold binds the task, the actor and the worktree. Reviews are rounds bound to the change's fingerprint. |
| ADR-0007 | Linked successor Runs; Claim release, expiry or force recovery; executor liveness; budgets carried forward | **Amended.** Holds replace Claims and Runs. A lapsed hold on a started task keeps the task for its worktree, which resumes it under a new linked hold. Only a person takes over. Git state is never touched. A check's start without a result counts for nothing. |
| ADR-0008 | Review Records bound to the Envelope revision; Envelope budgets and revision triggers; Attention | **Amended.** Rework stays on the original Issue. Rounds bind to the fingerprint. The rework cap and the token budget come from `.fid.toml`, and NEEDS YOU replaces Attention. Only a later round or a person closes a major finding. |
| ADR-0009 | Multi-Issue Envelope revision and partial progress | **Superseded.** Envelopes are removed. |
| ADR-0011 | The durable model's Work Envelopes, Claims and Runs, Stops, Attention, Kit Activation state, Gate Definition Registrations, Execution Attempts, Evidence Records, Obligation evaluation, Deviations and Human Outcomes | **Amended.** These go with the kernel. Issues, the journal, check results with Artifacts, review rounds and findings stay, joined by holds, the policy snapshot and a person's decisions. Its ownership and durability test stands. |
| Decisions record | D2, D5, D6, D7, D8, D9, F3, review outcome item 2, the work model's same-model review and Assured review, the later list | **Amended** on 10 and 11 November 2026, marked in place. |
