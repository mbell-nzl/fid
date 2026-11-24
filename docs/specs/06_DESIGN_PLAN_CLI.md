# CLI CONTRACT PLAN - fid

> **Kernel command removal delivered on 2026-11-15 (RESET 1.5c2).**
> ADR-0014 and `RESET.md` govern the lean core. Retired commands are unknown
> commands; no administrative aliases or interim review/finding stubs exist.
> Event catalogs and tables are removed in 1.5d, after domain removal in 1.5c3.

This document defines the aligned command vocabulary and common transport
contract. Exact leaf request/response schemas are accepted with their D3
vertical slice before implementation.

## 1. Global Contract

```text
fid [--project <path>] [--json] <command> ...
```

- fid needs a git repository, and git 2.31 or later. Discovery runs
  `git rev-parse` twice from the working directory: once for the worktree's
  top level (`--show-toplevel`) and once for the git common directory
  (`--path-format=absolute --git-common-dir`, which needs git 2.31). The
  project root is that top level, with symlinks resolved, so fid works from
  any subdirectory. The ledger is `$(git rev-parse --git-common-dir)/fid`,
  which every worktree of the repository shares; fid never looks for a
  ledger anywhere else (`04_TRD_CONFIG.md` section 1).
- `--project <path>` works like `git -C <path>`: discovery starts at `<path>`
  instead of the working directory, and the root is that path's top level.
  `export import --repository` follows the same rule.
- Outside a git repository, every command that opens or creates a project,
  `init` included, refuses with `repository-unavailable` (precondition, exit
  4; details `reason: not-a-repository` and `detail`). In a repository with no
  ledger, the failure is `project-not-found` (exit 3).
- fid keeps ledger writes inside `$(git rev-parse --git-common-dir)/fid`,
  where git cannot track them. Init also writes the committed `.fid.toml`
  and, when configured interactively, user config (below). `backup` also
  writes fid's own clone in the user cache and pushes to the backup target,
  and `restore` fetches into a temporary repository that it deletes
  afterwards (`04_TRD_CONFIG.md` section 2.4). Every command
  that opens or creates a project, `init`
  and `export import` included, refuses with `repository-unavailable`
  (precondition, exit 4; details `reason: unsafe-path` and a `detail` naming
  the locations involved) before it writes anything, where:
  - git and the worktree's `.git` entry name different ledgers, compared with
    symlinks resolved. The usual cause is `GIT_DIR` or `GIT_COMMON_DIR`
    naming another repository. Naming the worktree's own git directory, as
    git does for hooks, is fine, and so is a `.git` that is a symlink.
  - the git common directory is inside the worktree but not within a `.git`
    directory, as `git init --separate-git-dir=<worktree>/metadata` makes it.
    There `git add --all` would track the ledger, and ignore rules wouldn't
    help, because `git add -f` still adds an ignored file. A linked
    worktree's or a submodule's common directory is outside its worktree or
    within a `.git` directory, so they work.

  A newline in a path doesn't stop discovery. Each run reports one path,
  ended by a newline, and fid strips exactly that newline, because a path
  may itself end in one. An older git echoes `--path-format=absolute` back
  in its report, which then names no ledger the worktree agrees with, so fid
  refuses with `unsafe-path`.
- Output names ledger files by absolute path, in Human text and in JSON, never
  relative to the project root, because a linked worktree's ledger is outside
  that worktree. That covers `init`'s `journal_path`
  and Human message (`initialized fid project in <ledger>/`), `journal
  inspect`'s `journal_file`, the `journal_file` and `lock_path` failure
  details, and Human lines that name a journal file or the projection. A
  Human `reason:` line keeps its 200-scalar bound, so a long path there may
  end in `…`; JSON carries it whole.
- `--json` selects the canonical machine representation for operator commands.
- JSON-only backend-port commands never expose another format selector.
- There is no `--format`, `--toon`, `--yaml`, or universal Markdown output.
- Mutations that may be retried accept `--idempotency-key` or the corresponding
  JSON field.
- Concurrent mutations accept `--expected-revision` where applicable.
- Project-relative paths are normalized and cannot escape approved roots.

Ledger mutations made through the work commands (issues, labels, links, the lean
task commands and `artifact verify`) take `--actor` for attribution; the
ledger-level commands `init`, `export create`, `export import`, `backup`,
`restore` and `rebuild` take none. Work commands use
built-in actor-type authority: `agent:<id>` and `human:<id>` may work;
`service:<id>` is refused. A person's decisions require `human:<id>`.
Complex Artifact verification uses exactly one `--input <path|->` JSON
document and requires `--idempotency-key`. All retained command leaves use
the Response Envelope except the `journal inspect` and `sql` diagnostics.
`export import` accepts `--json` and `--repository`, and refuses `--project`.

## 2. JSON Envelope

Success:

```json
{
  "contract_version": "fid-cli.v1",
  "ok": true,
  "operation": "issue.get",
  "result": {"issue": {"id": "FID-123", "kind": "task", "title": "Example", "status": "todo", "status_category": "unstarted", "workflow_id": "task-default", "parent_id": null, "priority": null, "sequence": null, "spec_md_path": null, "revision": 1, "labels": [], "lane": null, "created_at": "2026-08-22T00:00:00Z", "updated_at": "2026-08-22T00:00:00Z", "acceptance_criteria": []}},
  "event": null
}
```

Failure is written to stderr with non-zero process status:

```json
{
  "contract_version": "fid-cli.v1",
  "ok": false,
  "operation": "issue.get",
  "error": {
    "code": "issue-not-found",
    "category": "not-found",
    "retryable": false,
    "message": "Issue FID-123 was not found",
    "subject": {"type": "issue", "id": "FID-123"},
    "details": {}
  }
}
```

No Human message substring is parsed for control flow. Unknown major contract
versions and required enum values reject.

Envelope rules:

- `contract_version` and `operation` are present on success and failure.
- `result` is an object, never a bare array or scalar; list operations use a
  named collection key.
- `event` is `null` for a non-mutating operation and
  `{"id": ..., "occurred_at": ...}` for a mutation that appended one Event.
- `error.details` is always present and may be `{}`; `error.subject` is `null`
  when the failure has no single subject.
- Success and failure envelopes are mutually exclusive: a success envelope has
  no `error` key, a failure envelope has no `result` or `event` key.

This transport wrapper is named the **Response Envelope** in specifications,
code and tests.

`error.category` is a closed set: `invalid-input`, `not-found`, `conflict`,
`authority`, `precondition`, `integrity`, `internal`, `gate-unsuccessful`
(the check runner’s exit-8 category).

## 3. Target Command Tree

### Project and ledger

```text
fid init
fid backup
fid restore
fid status [--all]
fid show <id>
fid rebuild
fid maintenance check
fid journal inspect <event-id>
```

RESET 1.5d2: `maintenance check` returns `{status, errors, warnings, info}`,
with `info: {issue_count, recent_activity_7d}`. The removed kernel count and
warnings are absent. Started tasks need no kernel ownership record to pass
maintenance. `warnings` remains an array and `--warn-as-error` retains its
generic promotion behavior. Structural, journal/retry-binding and Artifact
integrity checks still report divergence without repairing data.

Maintenance reports journal/projection divergence, structural integrity and
Artifact state. RESET 1.5d2 removes its kernel checks.

`fid init [--prefix <prefix>]` creates the ledger and a committed project
identity file at `<root>/.fid.toml`. It uses an existing valid identity
without rewriting the file; invalid or unreadable identity files are a clear
refusal, and an existing ledger remains `project-already-initialized`.
The project ID is `project-` plus 32 lowercase hex digits (128 random bits),
shared with the new ledger's config and genesis event. fid never stages or
commits the identity file. Text success output is:

```text
initialized fid project in <absolute-ledger-path>/
project ID: project-<32 lowercase hex digits>
commit <absolute-root>/.fid.toml to keep this identity in clones
```

The `project.init` JSON result is
`{project:{id,issue_prefix,root},journal_path,project_identity_path}`.
`project_identity_path` names `.fid.toml` by absolute path; the existing
`project.id` field names its ID. The response keeps its `project.initialized`
event metadata.

After creating the ledger and printing the success lines above, init asks
`Backup repository (git URL or local path; empty to skip): ` only when
stdin and stdout are both terminals, `--json` is off, and no backup target
is set. It stores a non-empty answer in `[backup] target` in
`$XDG_CONFIG_HOME/fid/config.toml` (default `~/.config/fid/config.toml`) and
prints the saved config path. Empty answers skip setup; runs that do not
prompt consume no stdin. JSON output stays a single response envelope.
User config is checked before ledger creation when prompting is eligible.
Invalid, unreadable or non-regular `.fid.toml`, and malformed or unreadable
user config in a terminal, use the existing `invalid-argument` code and
Human `reason: invalid-argument` (category `invalid-input`, not retryable,
exit 2, subject null). The details are
`{argument:<absolute config path>,value:null,allowed:["valid project.init configuration"]}`.
The bounded message names the file and cause; neither file is overwritten
and no ledger is created on these refusals. No new code or reason is added.

If reading the answer or saving it fails after ledger creation, the ID and
commit instruction have already been printed. The error begins
`ledger initialized, but backup target setup failed:` and reports the cause,
with the existing `internal-error` code (exit 7). Other internal init
failures keep the fixed message `project initialization failed`.
See `04_TRD_CONFIG.md` section 2 for file formats. `fid restore` is
below, and `fid status` warns until a target is set (below).

`fid backup` copies the ledger off the machine (decisions record D4 / F1). It
takes no arguments. It exports the ledger in `export create`'s format,
validates the archive as `export validate` does, writes it into the project's
folder in the backup repository, commits and pushes. The project ID in
`.fid.toml` names the folder, and must match the ledger's; a backup never
replaces a folder another ledger of the project wrote, such as another
clone's. The target is `[backup] target` in user config. fid works in its own clone of the backup
repository, never in a human's checkout, and uses the user's own git
credentials. Backups run outside fid's write path: a backup takes no journal
lock and appends no Event, and no other command pushes. The loop runs it after
each task and a daily cron job runs it too; concurrent backups wait for each
other. `04_TRD_CONFIG.md` section 2.4 has the folder, the clone and the
record. Text success output is one line:

```text
backup pushed: <project-id>/ at commit <40 hex digits>
backup unchanged: <project-id>/ is current at commit <40 hex digits>
```

`unchanged` means the target's tip already holds this exact export, so fid
made no commit and pushed nothing. The `project.backup` JSON result is
`{backup:<record line>,record_path}`: the attempt's line in the record, and
the record file's absolute path. `event` is null.

Every attempt that reaches the ledger is recorded, failed ones included, and
every failure exits non-zero with a typed failure:

| Failure | Code / exit | Details | Recorded kind |
|---|---|---|---|
| No target set | `invalid-argument` / 2 | `argument`: user config path; `allowed`: `["[backup] target: a git URL or local path"]` | `no-target` |
| Malformed or unreadable user config | `invalid-argument` / 2 | `argument`: user config path; `allowed`: `["valid project.backup configuration"]` | `user-config` |
| `.fid.toml` missing, invalid, or naming a project other than the ledger's | `invalid-argument` / 2 | `argument`: absolute `.fid.toml` path; `allowed`: `["[project] id matching the ledger's project ID"]` | `project-identity` |
| Another backup holds fid's clone for longer than the 2-minute wait | `lock-unavailable` / 4 | `lock_path` | `lock-unavailable` |
| Target unreachable, or not a repository | `repository-unavailable` / 4 | `reason: not-a-repository`; `detail`: git's message | `target-unreachable` |
| Target rejected the push | `repository-unavailable` / 4 | `reason: push-rejected`; `detail`: git's account of the rejection | `push-rejected` |
| The project's folder holds another ledger's backup, or names no ledger | `repository-unavailable` / 4 | `reason: different-ledger`; `detail`: the genesis Event hash the folder's `backup.json` names, if any, and this ledger's | `different-ledger` |
| Export or its validation failed, or the journal's first Event can't be read | `internal-error` / 7 | `{}` | `export-failed` |
| Anything else, such as git missing | `internal-error` / 7 | `{}` | `internal` |

The subject is null and `invalid-argument`'s `value` is null.
`repository-unavailable` is retryable; the next backup retries cleanly,
except that `different-ledger` persists until a person keeps one ledger for
the project (`04_TRD_CONFIG.md` section 2.4). Task 1.2b adds two
`repository-unavailable` reasons, `push-rejected` and `different-ledger`,
which Mark approved on 7 and 8 November 2026 (decisions record). Outside
a repository, at an unsafe ledger location, and in a repository without a
ledger, `backup` refuses as every project command does (section 1) and
records nothing, because there is no ledger to record in. If the record
can't be written, the failure says so; after a successful push that is
`internal-error` with the message `backup pushed, but recording it failed:`
and the cause.

`fid restore` rebuilds a lost ledger from its backup (decisions record
D4 / F1). It takes no arguments. It reads the project ID from `.fid.toml` in
the working tree and the target from `[backup] target` in user config,
fetches only the tip of the target's `main`, checks the project's folder,
and imports its archive as `export import` does, into the ledger every
worktree of the repository shares, so it works from any worktree. It
restores only into a repository with no ledger: it never replaces, merges
or updates one, so a rerun once a ledger exists refuses. It takes no lock,
appends no Event, never uses fid's backup clone, and writes no backup
record: a restored ledger starts without one. `04_TRD_CONFIG.md` section 2.4
has what it fetches and checks. Text success output is one line:

```text
restored <project-id>/ from commit <40 hex digits> into <absolute-ledger-path>/
```

The `project.restore` JSON result is
`{restore:{project_id,commit,archive_digest,event_count,terminal_event_hash,genesis_event_hash},ledger_path}`:
the commit restored from, the archive's `sha256:` digest, the restored
journal's Event count, terminal Event hash and genesis Event hash, and the
ledger's absolute path. `event` is null.

Every failure exits non-zero with a typed failure, and a refused restore
leaves no ledger where there was none:

| Failure | Code / exit | Details |
|---|---|---|
| A ledger already exists, checked before config or the network; or one appeared while the restore fetched | `project-already-initialized` / 4 | subject the project root; `path`: the root |
| `.fid.toml` missing or invalid | `invalid-argument` / 2 | `argument`: absolute `.fid.toml` path; `allowed`: `["[project] id naming the backup folder"]` |
| No target set | `invalid-argument` / 2 | `argument`: user config path; `allowed`: `["[backup] target: a git URL or local path"]` |
| Malformed or unreadable user config | `invalid-argument` / 2 | `argument`: user config path; `allowed`: `["valid project.restore configuration"]` |
| Target unreachable, or not a repository | `repository-unavailable` / 4 | `reason: not-a-repository`; `detail`: git's message |
| Target has no `main`, or its `main` has no folder for this project | `repository-unavailable` / 4 | `reason: no-backup`; `detail`: which is missing |
| The project's folder is damaged | `archive-invalid` / 6 | subject export `<project-id>/export.zip`; `archive_path`: `<project-id>/export.zip`; `checks` |
| Anything else, such as git missing | `internal-error` / 7 | `{}` |

Other subjects are null, and `invalid-argument`'s `value` is null.
`repository-unavailable` is retryable. The folder is damaged when its
`backup.json` is missing, isn't a manifest, has a `schema` fid can't
restore, or names another project (check `backup-manifest-invalid`); when
the export's `manifest.json` is missing (`archive-missing`), malformed or
names another project (`backup-manifest-invalid`); when a member isn't a
regular file (`backup-manifest-invalid`); when the reconstructed ZIP's digest isn't
`backup.json`'s `archive_digest` (`archive-digest-mismatch`); or when the
archive fails `export validate`'s checks, which keep their own codes. fid's
own three checks have a null `entry`. `archive-invalid` is an integrity
failure and isn't retryable; `archive_path` names the reconstructed export
logically as `<project-id>/export.zip`, since the temporary ZIP is deleted.
Task 1.2c adds the `repository-unavailable`
reason `no-backup`, which Mark approved on 8 November 2026, and those three
check codes; no other code or reason. Outside a repository and at an unsafe
ledger location, `restore` refuses as every project command does
(section 1).

`fid status` shows what needs you, the epics and tasks in progress and
ready, and how current the ledger's backup is (decisions record D4 / F1
and D8, tasks 1.3 and 1.4b). It takes no arguments and one flag,
`--all`. It reads the projection, the backup record and user config, and
writes nothing. Text output is a rounded box 72 columns wide, headed
`fid status` and today's local date:

```text
╭ fid status ────────────────────────────────────────────── 2026-11-11 ╮

  IN PROGRESS
  ◆ FID-e3a1c0  Operating experience              ██░░░░░░░░  1/4 done
      ● FID-71a0c2  Export validation                      in progress
  ◆ FID-8a1b2c  Loop contract                     ░░░░░░░░░░  0/4 done
      ● FID-1f2e3d  Contract fields and exit …  interrupted in fid-wt2

  READY
  ◆ FID-e3a1c0  Operating experience                           3 to do
      ○ FID-c2d4e5  Retry a rejected backup pu…  blocked by FID-71a0c2
      ○ FID-9f13aa  Harden login rate limiting                    next
      ○ FID-0d4e8b  Document backup age             working in fid-wt3
  ◆ FID-8a1b2c  Loop contract                                  3 to do
      ○ FID-2f2e3d  Golden tests for the contract     after FID-1f2e3d
      ○ FID-3f2e3d  Scenario test for rework          after FID-1f2e3d
      ○ FID-4f2e3d  Version the contract              after FID-1f2e3d
  ○ FID-77aa01  Tidy help text
╰──────────────────────────────────────────────────────────────────────╯
```

- **NEEDS YOU** lists what waits on a person, each item a summary and its
  next step. Besides backup, waiting-for-you and worktree-missing, RESET 2.1
  adds `rework-cap`, `approval` and `merge`. Approval and merge items appear
  only after the latest hand-off's required primary and separate review gates
  and task-wide majors gate are satisfied, using the same evaluator as done. As in V5, a person closing the last
  major can satisfy these gates with a current `CHANGES` round; a later
  `APPROVE` is not required in that case.
  An Assured task with missing or stale approval says `<task>: approve, then
  merge PR #<n>`, adding `, no other-model review` when applicable; its next
  action is `fid approve <task> --actor human:<name> --reason <text>`.
  With current approval, or in other lanes, one `merge` item says
  `<task>: merge PR #<n>`, with the hand-off's PR URL as next action.
  If that hand-off still owns its kept review reservation and no longer
  merges into the local default branch, a
  `default-conflict` item replaces this action: `<task>: PR no longer merges
  into origin/<default>; its builder's worktree takes it back`. Both text
  and JSON name `fid next --claim in the builder's worktree`. A person's
  takeover releases that reservation, so the existing approval or merge
  item applies again.
  This is local evidence only: fid never fetches.
  A recorded merge suppresses both items; the existing offline CI states
  apply. Completed tasks have neither. Each IN PROGRESS task's JSON keeps
  the same `approval` and `no_other_model_review` fields as show (below),
  even while reviews are still running.
- **Epics and tasks** appear; subtasks never do. An epic whose own state is
  completed or cancelled is hidden. An epic is in progress when its own
  state's category is `started`, or any of its tasks is started or
  completed; otherwise it is to do. Its progress is its completed tasks over
  its tasks that aren't cancelled.
- **IN PROGRESS** shows each in-progress epic, with its bar and count, and
  its started tasks beneath it with their workflow state, underscores read
  as spaces. An interrupted task (a lapsed hold on a started task) shows
  `interrupted in <dir>` instead. Started tasks with no listed epic follow,
  none or a hidden one.
- **READY** shows every unstarted task: under its in-progress epic (`N to
  do`), then under its to-do epic (`N tasks`, or `no tasks`), then those
  with no listed epic. An in-progress epic with nothing to do has no line
  here. Each task's right-hand column is `fid next`'s verdict for this
  worktree, read without taking work: `next`, `blocked by <ID>` (and `+N`
  for more), `after <ID>`, `working in <dir>`, `working here`,
  `epic blocked by <ID>` (and `+N`), or `epic closed`; otherwise it is
  blank. `<dir>` is the worktree's directory name. The task marked `next`
  is what `fid next` would hand out here, which can be below the first
  row, since unblocking work comes first. While this worktree already
  holds a task, that task shows `working here` and no task says `next`.
  A lapsed hold on an unstarted task frees it for anyone.
- The text shows at most 10 task rows in READY, then
  `+N more: fid status --all`, where N counts the hidden task rows and any
  hidden epic with no tasks. Epic lines don't count toward the 10. `--all`
  shows every row.
- **Order:** epics, and the tasks in each list, by sequence, then priority
  (lower first), then creation time, then ID; a missing sequence or
  priority comes last, as in `issue list --order`.
- A section with nothing in it is left out. A project with no epics or
  tasks at all shows `○ No tasks yet` and the command that creates a task,
  `fid issue create --kind task --title <title> --actor human:<name>
  --idempotency-key <key>`.
- Rows run to column 70, the right-hand text flush right. A title too long
  for its column ends in `…`, dropping whole an escape (below) that would
  cross the cut, as failure messages do. When the right-hand text would
  leave a title fewer than 10 columns, as the longest IDs and workflow
  states can, it takes a line of its own beneath the row, still flush
  right.
- Summaries, next steps and paths are never cut; they wrap. Prose wraps at
  spaces to column 70, its continuation lines indented to where it began,
  and never breaks inside a word: a word longer than the line, such as a
  long path, stands whole past the edge. Nor does it break between a date
  and its time, such as `2026-11-05 14:00`. A next step that is a fid
  command (it starts with `fid `) wraps between its arguments, never
  between a flag and its value. Each broken line ends in ` \` and its
  continuation lines are indented two more than the command, so copying
  both lines runs the whole command, and copying only the first leaves
  the shell waiting rather than running half of it:

  ```text
    ○ No tasks yet
        → fid issue create --kind task --title <title> \
            --actor human:<name> --idempotency-key <key>
    ⚠ Can't read user config /home/mark/.config/fid/config.toml; last
      backed up 1 hour ago
        → add [backup] target to /home/mark/.config/fid/config.toml,
          then fid backup
  ```
- Text fid didn't write, such as titles, paths and the backup record's
  failure kinds, has its controls escaped before it is fitted to a column,
  as failure messages' are (section 8, Phase 24 amendment): LF, CR and TAB
  as `\n`, `\r` and `\t`, other control characters as `\uXXXX`. JSON keeps
  the values as they are.
- Colour is used only when stdout is a terminal and `NO_COLOR` is unset or
  empty. The symbols are always Unicode.
- On a terminal, text taller than the screen pages (below).

The backup's line is the box's last, except when it needs you: then it
appears only as the NEEDS YOU item.

**Policy health (RESET 1.5a1; Mark, 14 November 2026).** Status reads the
trusted local `origin/<default>` policy, using built-in defaults when
there is no origin or file. It never writes a policy snapshot or any
other Event, including after the effective digest changes. Policy
warnings appear inside the box, just above the backup line (or last when
the backup is under NEEDS YOU), in the order below. They appear even in an
empty project and are not NEEDS YOU items. They wrap and indent through
the display package as backup warnings do.

| Kind | Summary | Next action |
|---|---|---|
| `policy-unavailable` | `Policy unavailable: <problem>` | `check the repository: git remote -v` if listing remotes fails; `git remote set-head origin <branch>` for an unknown default; `fix .fid.toml on the default branch, merge it and git fetch` for an invalid or unreadable policy file |
| `checks-not-configured` | `No checks set for Routine, Standard, Assured in .fid.toml` | `add checks to [lanes.routine], [lanes.standard] and [lanes.assured], then merge .fid.toml and git fetch` |
| `test-paths-not-configured` | `No test paths set in .fid.toml` | `add [review] test_adequacy, then merge .fid.toml and git fetch` |

Unavailable policy leaves the whole board visible, with only its one
policy warning: status cannot know which setup keys are missing. When
policy loads, the checks warning lists only lanes missing their `checks`
key, in Routine, Standard, Assured order. For one lane it reads `Standard`
and `[lanes.standard]`; for two, `Routine, Assured` and
`[lanes.routine] and [lanes.assured]`. A missing `[review] test_adequacy`
key gives the test-path warning. Explicit `checks = []` and
`test_adequacy = []` count as set up and give no warning.

| State | When | Line |
|---|---|---|
| `current` | The last backup is at most 24 hours old and nothing is wrong | `✓ Backed up 3 hours ago` |
| `not-backed-up-yet` | No backup is recorded, the ledger began at most 24 hours ago, and nothing is wrong | `○ Not backed up yet` |
| `warning` | Any reason below, while the age is at most 72 hours | `⚠` and the summary, with its next step |
| `needs-you` | The age is over 72 hours | The NEEDS YOU item |

The last backup is the record's latest `pushed` or `unchanged` line, at its
`finished_at`. The age is now less the last backup or, when none is
recorded, less the genesis Event's time. A restored ledger starts with no
record and an old genesis Event, so it needs you until its first
`fid backup`. A time in the future counts as now. The reasons, in this
order:

| Reason | When | Next step |
|---|---|---|
| `no-recent-backup` | The age is over 24 hours | `fid backup` |
| `last-backup-failed` | The record's last line is `failed`, of any kind, however recent the last backup. The text names the kind, and how many failed in a row when it's more than one | `fid backup` |
| `no-target` | User config has no `[backup] target` | `add [backup] target to <user config path>, then fid backup` |
| `user-config-unreadable` | User config is malformed or unreadable, or `$HOME` and `$XDG_CONFIG_HOME` don't locate it | `fix <user config path>, then fid backup`, or `fid backup` when there's no path |
| `record-unreadable` | The record exists but can't be read; the age is then unknown | `fix or remove <record path>, then fid backup` |

The next step shown is the first that applies of: fixing user config, setting
the target, fixing the record, and `fid backup`. A record
line that isn't a whole record, such as a torn last line, is skipped;
appends repair it later. None of these make `status` fail. Relative times
read as minutes, hours or days ago (`just now` under a minute; hours below
48); absolute times are local, as `YYYY-MM-DD HH:MM`. Human text names
failure kinds, never kernel nouns or IDs.

The `project.status` JSON result holds everything the text shows, never
capped, with or without `--all`. `event` is null.

```json
{"generated_at": "2026-11-08T09:00:00Z",
 "backup": {"state": "warning", "reasons": ["last-backup-failed"],
            "summary": "Last backup failed (push-rejected); last backed up 2 hours ago",
            "next_action": "fid backup", "last_backup_at": "2026-11-08T07:00:00Z",
            "ledger_started_at": "2026-11-04T05:00:00Z", "age_seconds": 7200,
            "latest": <record line>, "consecutive_failures": 1, "target_set": true,
            "user_config_path": "/home/mark/.config/fid/config.toml", "user_config_error": null,
            "record_path": "<ledger>/backups.jsonl", "record_error": null},
 "warnings": [],
 "needs_you": [{"kind": "backup", "summary": "...", "next_action": "fid backup", "subject": null}],
 "in_progress": {"epics": [<epic>], "tasks": [<task>]},
 "ready": {"next": "FID-9f13aa", "epics": [<epic>], "tasks": [<task>]},
 "hint": null}
```

- `generated_at` is UTC. `last_backup_at` and `ledger_started_at` are UTC
  RFC 3339 times; `last_backup_at` is null when no backup is recorded.
- `age_seconds` is null when the record can't be read. `latest` is the
  record's last whole line, or null. `target_set` is null when user config
  can't be read, and `user_config_path` null when it can't be located.
  `user_config_error` and `record_error` carry the cause, or null.
- `summary` and `next_action` are the text's, unescaped; `next_action` is
  null when nothing needs doing.
- `warnings` is always present, `[]` when there are none. Each warning
  has `kind`, `summary` and `next_action` (strings), plus only its kind's
  fields. Text and JSON share the same summary and next action. The
  warning shapes are:

  ```json
  {"kind": "policy-unavailable", "problem": "origin's default branch is unknown", "ref": null, "commit": null,
   "summary": "Policy unavailable: origin's default branch is unknown", "next_action": "git remote set-head origin <branch>"}
  {"kind": "checks-not-configured", "lanes": ["routine", "standard", "assured"],
   "summary": "No checks set for Routine, Standard, Assured in .fid.toml",
   "next_action": "add checks to [lanes.routine], [lanes.standard] and [lanes.assured], then merge .fid.toml and git fetch"}
  {"kind": "test-paths-not-configured", "summary": "No test paths set in .fid.toml",
   "next_action": "add [review] test_adequacy, then merge .fid.toml and git fetch"}
  ```

  `ref` and `commit` are strings when known, null otherwise. `lanes` uses
  lowercase names in the text's order.
- A NEEDS YOU item is `{kind, summary, next_action, subject}`. `kind` includes
  `backup`, `waiting-for-you`, `worktree-missing` and `epic-ready-to-close`.
  `subject` is null for backup and `{type, id}` for an Issue. An open epic
  (`todo` or `in_progress`) with at least one direct task and no open task
  gets `epic-ready-to-close`, subject `{type: "epic", id}`, summary
  `<id> has no open tasks and is ready to close`, and next action
  `fid done <id> --actor human:<name> --reason "<why>"`. Completed and
  cancelled tasks are finished; a task in review is open. Empty, closed
  and cancelled epics do not appear. Eligibility uses the same predicate
  as epic completion, over all direct tasks in status's one snapshot,
  before display truncation. Status writes no Event and needs no network.
- An epic is
  `{id, title, status, status_category, state, sequence, priority, progress: {done, total}, tasks}`,
  `state` being `in-progress` or `to-do` and `tasks` the section's tasks.
- A task is
  `{id, title, status, status_category, epic_id, sequence, priority, blocked_by}`,
  `blocked_by` listing its open blockers' IDs in order. IN PROGRESS tasks
  also have `working` and `kept`, null when absent. READY tasks add
  exactly `next`'s queue entry fields: `reason`, `after`, `working`, `kept` and
  `epic_blocked_by` (below). `working` is
  `{hold_id, worktree, here, actor, since, expires_at, last_heartbeat_at, lapsed}` or null;
  times are UTC RFC 3339 and `lapsed` is true for interrupted work.
- `ready.next` is the ID `fid next` would hand out here, or null when
  nothing can be handed out or this worktree already holds a task. READY
  keeps status's order, rather than `next`'s unblocking-first queue order.
- `hint` is `{summary, next_action}` for a project with no epics or tasks,
  and null otherwise.

Outside a repository, at an unsafe ledger location, and in a repository
without a ledger (`project-not-found`, exit 3), `status` refuses as every
project command does (section 1). Journal replay failures keep their codes;
any other failure, such as a journal whose first Event can't be read, is
`internal-error` (exit 7). Task 1.3a adds no codes or reasons.

**The pager** (task 1.3b). `status` and `show` page long text, as git's
commands do; no other command pages.

- Text pages only when stdout is a terminal and the text is taller than
  the screen. fid reads the terminal's size (`TIOCGWINSZ`) and counts rows
  at its width, so a line the terminal wraps counts as more than one.
  Text that fits, `--json` and piped output are printed directly.
- The pager is `$PAGER` when it is set, and `less` otherwise; an empty
  `PAGER` or `cat` means none. As git does, a pager with arguments or
  shell syntax, such as `PAGER="less -S"`, runs through `sh -c`, and a
  plain name is found on `PATH` and run directly.
- fid sets `LESS=FRX` (quit if the text fits one screen, keep colours,
  don't clear the screen) only when `LESS` is unset: the user's own `LESS`
  wins.
- Colour is decided by stdout being a terminal, before the text reaches
  the pager.
- When the pager can't start, such as when `less` is missing, fid prints
  the text directly. Quitting the pager early isn't an error. fid waits
  for the pager to exit, through an interrupt typed in it, and its exit
  code is its own.
- There is no `--no-pager`: `PAGER=cat` or a pipe turns paging off.

`fid show <id>` shows one issue (decisions record, 8 November 2026 details
for task 1.3). It takes the issue's ID, matched exactly, with no prefix or
case matching, and no flags. It is a read: it refreshes the projection like
every read, records no Event and makes no network calls. Text output is
`status`'s 72-column box, its top border naming the issue's ID and kind:

```text
╭ FID-71a0c2 ──────────────────────────────────────────────────── task ╮

  ● Export validation                                      in progress

  Epic        ◆ FID-e3a1c0  Operating experience              2/7 done
  Priority    P1
  Blocks      ○ FID-c2d4e5  Retry a rejected backup push         to do
              ○ FID-0d4e8b  Document backup age                  to do
  Labels      backup, export
  Spec        docs/specs/02_TRD.md

  ACCEPTANCE CRITERIA
  1  fid export validate refuses an archive whose member digests don't
     match its manifest, naming the first bad member
  2  A truncated archive is refused, not read to its end

  SUBTASKS                                                    1/3 done
  ✓ FID-aa01b2  Check the manifest schema                         done
  ● FID-aa02b3  Validate member digests                    in progress
  ○ FID-aa03b4  Refuse a truncated archive                       to do
╰──────────────────────────────────────────────────────────────────────╯
```

- **Heading.** The issue's mark, its title, and its state flush right.
  Workflow states read with underscores as spaces, and `todo` as `to do`.
  The title wraps and is never cut, its state staying on the first line;
  a state that would leave the title fewer than 10 columns takes a line
  of its own beneath. An epic reads as `status` reads it: `in progress`
  once its own state's category is `started` or any of its tasks has
  started or completed, `to do` before then, and its own workflow state,
  such as `done`, once it is completed or cancelled.
- **Marks.** `✓` completed, `●` started, `○` unstarted, `◆` an epic, and
  `✗` cancelled, muted rather than red, since cancelled isn't an error.
- **Rows**, each only when it has something to show, in this order:
  - `Task`: a subtask's task, with its state;
  - `Epic`: the issue's epic (a subtask's is its task's epic), with
    `N/M done`;
  - `Progress`: an epic's bar and `N/M tasks done`, when it has tasks;
  - `Scheduling`: on an epic with scheduling, `one at a time, in sequence`
    for `sequential`, or `by priority` for `prioritised`;
  - `Priority`: `P<n>`;
  - `Blocked by` and `Blocks`, each issue on its own line under the value
    column, with its state;
  - `Labels`, comma-separated, and `Spec`, `spec_md_path` as stored;
  - `Lane`: on a task, `Routine`, `Standard`, `Assured`, or `not set`;
  - `Approval`: after Lane, on unfinished tasks when the latest hand-off's
    lane is Assured, including a lane raised by its floor. `○ missing` means no approval
    exists. A current approval reads `✓ <decision> <actor> · <at>`; a stale
    one reads `⚠ stale · <decision> <actor> · <at> · approved an earlier change`.
    Current is green, stale yellow, missing muted when colour is enabled.
    A second line says `no other-model review` when no current primary or
    separate round came from another model family than every builder's.
    Unknown agent builder models and unknown reviewer models cannot establish
    other-model coverage. Every current round may establish it; stale rounds
    cannot. The approve prompt and done use this same reader.
  - `Working`: on a held task, `here, since 09:02, until 10:02`, or
    `in <full path> by <actor>, since 08:40, until 09:40` elsewhere. An
    interrupted task reads `⚠ interrupted in <full path> by <actor>: work
    stopped at 07:15 without a park`, using the hold's last heartbeat. Times are local `HH:MM`, adding the
    date when it isn't today. A lapsed hold on an unstarted task frees it,
    so it has no Working row.

  Labels are 10 columns wide after the two-space indent, so values begin
  at column 14. Checks follows `Approval` (or `Lane` when approval does not
  apply). Done tasks use their completion receipt below; later tasks add
  next as the last, a NEEDS YOU block under the heading, and an epic's plan
  and validations.
- **Reference lines**, in rows and lists, follow `status`'s row rules: a
  mark, the ID, two spaces, the title cut with `…`, and the right-hand
  text flush to column 70, or on a line of its own beneath when it would
  leave the title fewer than 10 columns.
- **Blockers** are read from either link, the blocker's `blocks` or this
  issue's `blocked_by`, any issue kind; `Blocks` is the reverse. The text
  lists only issues that aren't completed or cancelled, so `Blocked by`
  means waiting now. Both lists, and the children, are in `status`'s
  order.
- **ACCEPTANCE CRITERIA**, for tasks and subtasks: each criterion's
  ordinal, then its text, wrapped under itself, in ordinal order.
- **Children.** An epic's `TASKS` lists every task, completed and
  cancelled included, uncapped. A task's `SUBTASKS`, headed `N/M done`,
  lists every subtask. Each child shows its state on the right or, while
  it is unstarted and has open blockers, `blocked by <ID>` and `+N` as
  `status` does. A subtask has no children.
- Sequence, revision, workflow ID and times appear only in the JSON's
  `issue`. Neither text nor JSON shows open concerns, triage status, the
  owner role or actor (apart from actors in Working, Approval and Completed),
  or links other than `blocks` and `blocked_by`.
- Text fid didn't write is escaped, fitted and wrapped as in `status`,
  and long text pages as `status`'s does.

The `issue.show` JSON result holds everything, with a null `event`:

```json
{"issue": <issue.get's issue>,
 "state": null,
 "epic": {"id": "FID-e3a1c0", "kind": "epic", "title": "Operating experience",
          "status": "in_progress", "status_category": "started", "progress": {"done": 2, "total": 7}},
 "task": null,
 "working": null,
 "kept": null,
 "checks": [],
 "checks_base": null,
 "completion": null,
 "approval": null,
 "no_other_model_review": null,
 "warnings": [],
 "scheduling": null,
 "blocked_by": [{"id": "FID-1a2b3c", "kind": "task", "title": "Escape controls in status",
                 "status": "done", "status_category": "completed"}],
 "blocks": [<reference>, <reference>],
 "progress": {"done": 1, "total": 3},
 "children": [{"id": "FID-aa01b2", "kind": "subtask", "title": "Check the manifest schema",
               "status": "done", "status_category": "completed", "sequence": null,
               "priority": null, "blocked_by": []}]}
```

For unfinished tasks, approval's JSON is null before the first hand-off and for non-Assured
hand-offs. Otherwise it is `{"state":"missing"}`, or an object whose `state`
is `current` or `stale`, with `decision_id`, `actor`, `at` (RFC3339), and
`fingerprint` (the recorded fingerprint value). Freshness comes from
`approvalrecord.Current` against the latest hand-off's fingerprint; an earlier
approval of an identical change remains current. The state is also
`not-required` when a current waiver lowers the gates. Show's Lane row names
the waived lane and uses the shared current/stale decision presentation.
When waivers exist, show and status's in-progress task JSON include `waivers`,
an array of full decisions with `current` booleans. Status appends
`stale waiver` to a task's verdict when no waiver matches the latest hand-off.
An effective non-Assured lane does not request approval or separate review.
`no_other_model_review` is a boolean for Assured hand-offs, null when not
applicable or unknown. Both commands operate offline, without fetching.
Unavailable current policy or comparison retains approval evidence from the
recorded endpoints. If those endpoints cannot be fingerprinted locally, an
`approval-unavailable` warning names the task and problem; unknown approval
and disclosure fields are null. An absence of approvals or reviews remains
known without those objects. Neither command writes an approval or any Event.

For a **done task**, a `Completed` block follows the header rows and precedes
acceptance criteria and children (RESET 2.3b). It expands the immutable
completion record and only the evidence IDs that record names:

- Completion time, ID and `fid done` actor.
- PR number, merge commit and merge time.
- Each counted CI result's conclusion, run ID as `run <id>` and commit.
  When recorded, the full CI URL follows on its own continuation row,
  never wrapped or cut, even on a narrow terminal. No URL row is shown
  when no URL was recorded. Multiple counted runs remain separate.
- Completed hand-off ID and head, and the completion's fingerprint.
- Each counted round's ID, stream, verdict, reviewer actor and model.
- Assured approval's decision ID, person and time; the recorded
  `no other-model review` disclosure when true; `Rollback present` or
  `Rollback absent` as recorded.

Text reads `CI ✓ success · run 101 on <commit>` for a successful CI result.
Text shortens commit SHAs and the fingerprint to eight characters. Times
use the existing local clock format, adding the date when not today. Values
wrap beneath the label and external text is escaped as in other show rows.
The ordinary Hand-off and Approval rows are replaced by this block.

A done task's Checks rows list the required checks and exact result IDs
recorded by its completed hand-off. They show the pass time and `counts at
<head>` with the completed head, rather than comparing to today's HEAD,
policy or lane floor. Subsequent checks cannot replace these results.
Unfinished tasks retain their current check and approval comparisons.

JSON always includes `completion`: null except for done tasks. Its object
contains the completion record's typed fields (`id`, `task_id`, `handoff_id`,
`merge_id`, `observed_ci_result_ids`, `ci_result_ids`, `round_ids`,
`approval_id`, `fingerprint` with algorithm/value/base/head, `rollback_present`,
`no_other_model_review`, `actor`, `at`), plus these expansions:

- `handoff`: the named hand-off record, or null.
- `merge_observation`: the named merge observation, or null, with `number`,
  `merge_commit` and `merged_at` among its fields.
- `ci_results`: counted CI result records, including `id`, `run_id`,
  `head_sha`, `conclusion` and `url`.
- `rounds`: the counted review round records, including `id`, `stream`,
  `verdict` and `reviewer: {actor, model, ...}`.
- `approval`: the named approval record, or null, including `decision_id`,
  `actor` and `at`.
- `missing`: an array of explanations, such as `CI result missing · CIR-...`.

Full SHAs, fingerprints and recorded RFC3339 timestamps remain in JSON.
The existing top-level `handoff`, `floor`, `approval` and
`no_other_model_review` fields describe that completed evidence too. A named
approval has state `current` historically; no current fingerprint is computed.
Missing named records stay null or absent from their arrays, and appear as
`⚠ <kind> missing · <id>` in text. A done task without a completion row has a
completion object with empty expansions and `completion record missing`;
its unavailable completion record fields are omitted. Older receipts without
a fingerprint say `fingerprint missing`; rounds without a model say `model
not recorded`. Available facts still render. These reads need no historical
git objects and make no GitHub call. Epics, subtasks and cancelled tasks keep
their existing displays, with `completion: null`. Show is outside the loop
commands, so no declared loop-contract field changes.

- `issue` is exactly `issue get`'s. `state` is null except on an epic,
  where it is `in-progress` or `to-do` as in `status`, or `completed` or
  `cancelled` by its own state's category.
- `epic` is the issue's epic and `task` a subtask's task, or null.
- `working` is the task's live or interrupted hold, in `next`'s
  `{hold_id, worktree, here, actor, since, expires_at, last_heartbeat_at, lapsed}` shape,
  or null when nobody holds it (including a lapsed unstarted task).
- `scheduling` is `"sequential"` or `"prioritised"` on an epic with
  scheduling, null otherwise.
- `blocked_by` and `blocks` list every link in both directions, whatever
  the other end's state, each `{id, kind, title, status,
  status_category}`.
- `progress` counts the children, completed over not cancelled; it is
  null for a subtask and for a task with no subtasks.
- `children` is `[]` when there are none; a child's `blocked_by` lists
  its open blockers' IDs, in order.

An unknown ID, including an ID of another kind of record, such as an
acceptance criterion's, is `issue-not-found` (exit 3), as for `issue get`.
Its text names where to find IDs, a line only `show` adds:

```text
error: issue was not found
reason: issue-not-found; id=FID-nope
next: fid issue list
```

A missing ID is `invalid-argument` (exit 2). Outside a repository, at an
unsafe ledger location, and in a repository without a ledger, `show`
refuses as every project command does (section 1). Journal replay
failures keep their codes; any other failure is `internal-error` (exit 7).
Task 1.3b adds the operation `issue.show`, and no codes or reasons.

### Issue tracker

```text
fid issue create|get|list|update|transition
fid issue ac add|list
fid label add|list
fid issue label add|remove
fid link add|remove|list
fid next [--claim --actor <actor>] [--epic <epic>] [--explain]
```

Issues use fixed built-in states for their kind. New issues ignore retired
workflow assignments. Standalone inventory mutations record typed actor
attribution. No execution-chain authority is required for an Issue transition.
`issue transition` refuses every completed-category target for a task or epic and
points at `fid done`; subtasks retain their plain transitions. Cancelling
tasks and epics is a person's decision, with the contract below.

**`fid next`** (task 1.4a; decisions record, D8's 1.4 details) hands a
worktree its next task. It walks one order and serves the first task
nothing stops:

1. **This worktree's own hold.** A live hold is returned as it is (kind
   `resume`); one worktree never holds two tasks. A lapsed hold on a
   started task means the task was interrupted: it is kept for this
   worktree, and `--claim` takes a new hold that replaces the old one
   (kind `resume`). A lapsed hold on an unstarted task frees the task for
   any worktree, this one included.
2. **Rework in this worktree:** a review round asking for changes returns
   the builder’s kept task as kind `rework`, with a new hold linked to the
   hand-off’s ended hold. It precedes unblocking and new work. The task is
   already `in_progress`; the builder fixes, checks and hands it off again
   without another `start`. Another worktree cannot take it while kept.
3. **Review work:** the primary stream until its gate is satisfied, then
   the separate stream for Assured work. Required profiles come from the
   hand-off's frozen lane, paths and policy, as in done. Routine work without
   required profiles creates no review work. A task returned for rework
   creates none until the next hand-off.
4. **Unblocking work:** tasks that another open task waits on, by a
   `blocks` or `blocked_by` link, directly or through its open epic: a task
   that blocks an epic counts only while the epic has open tasks. Until
   task 3.1, all work counts as approved.
5. **The other unstarted tasks**, in `fid status`'s order: an in-progress
   epic's, then a to-do epic's, then those with no open epic; within each,
   sequence, priority (lower first), creation and ID, a missing sequence or
   priority last.

An approved task also returns to its kept builder worktree as `rework` when
its latest hand-off head conflicts with local `origin/<default>`. Eligibility
is the same satisfied review/majors gates as status's approve/merge action,
including the existing V5 person-closure exception. A clean merge is never
served. An unavailable ref or commit is not a conflict: `--explain` describes
why the local comparison is unavailable. Other worktrees cannot take it.
Optional `rework` JSON names `{cause: "default-conflict", default_ref,
default_commit, handoff_id, head, conflicting_paths}`; non-UTF-8 paths also
include base64 `conflicting_path_bytes`. `rework_unavailable`, when present,
explains missing local evidence. The `fid-cli.v1` contract is unchanged.
The builder merges the default branch, checks and hands off again. This
return adds no review round or cap charge. Earlier rounds and approval stay
current exactly while the new hand-off's fingerprint matches. A clean
approved task cannot be taken back to fix an open minor; it needs an answer
or deferral (2.1c).

The order holds every unstarted task and every task with an open hold or kept work; a
started task with neither a hold nor kept work is never handed out. A sequential epic serves one
task at a time, across worktrees. While one of its tasks is held (a live
hold, or a lapsed hold on a started task), it serves no other. Otherwise it
serves its next unfinished task in sequence or, links beating sequence, the
first task in sequence that the next task waits on through links in the
epic (directly or through others) and that nothing blocks. Prioritised
epics, and epics with no scheduling, go by priority. Every other task is
passed over with a typed reason, the first that applies:

| Reason | When |
|---|---|
| `outside-scope` | `--epic` names a different epic, or the task has no epic; takes precedence over other reasons |
| `not-accepted` | an agent created this task and no person has accepted it |
| `held-elsewhere` | another worktree's live build hold or review reservation |
| `built-by-actor` | the claiming actor ever held this task to build it, including ended holds |
| `rework-cap-reached` | changes rounds exceed the latest hand-off policy's cap for its recorded lane |
| `interrupted-elsewhere` | another worktree's lapsed or parked started work, or its kept rework |
| `epic-closed` | the task's epic is completed or cancelled |
| `epic-blocked` | the task's epic has a blocker that isn't completed or cancelled |
| `not-next-in-sequence` | its sequential epic serves another task first; `after` names the task the epic holds or serves, or, when it can serve none, its next task |
| `blocked` | a blocker, from either link, that isn't completed or cancelled; `blocked_by` names them |

Reasons are informational, and later tasks add more: a loop never branches
on them, only on `kind`.

- `fid next` reads (operation `next.get`). It refreshes the projection like
  every read and records no Event.
- `fid next --claim --actor <actor> [--model <family/model>]` takes the served task as a hold for
  this worktree (operation `next.claim`), in one `hold.taken` Event under
  the write lock. The Event ends this worktree's earlier hold and any
  lapsed hold on the task. Re-run in the same worktree, it returns the hold
  it already has, with a null `event`; there is no idempotency key. Two
  worktrees asking at once aren't refused: the lock serialises them, and
  the second passes over the first's task as `held-elsewhere`. The
  database backs the lock up with at most one open hold per task and per
  worktree; if that ever refuses, it is a bug, reported as
  `internal-error`.
- `--epic <epic>` scopes every kind of work, including review, resume and
  rework, to that epic's tasks. The usual order, dependencies, sequence,
  holds and rework cap still apply within it. Other tasks stay in the queue
  with `outside-scope` (text: `outside scope`), even when they block scoped
  work. Nothing available in scope returns `none`; unscoped workers keep
  their usual selection. An occupied worktree whose task is outside the
  requested scope returns `none`, retaining its session. An unknown ID is
  `issue-not-found` (exit 3); an existing non-epic is `unsupported-kind`
  (exit 2); an explicit empty value is `invalid-argument` (exit 2).
- `--explain` lists the whole order with a verdict for each task, with or
  without `--claim`. It changes only the text.

As of RESET 1.6d, `--model <family>/<model>` is required for an `agent:*`
actor taking work and optional for a `human:*` actor. Both components
match `[a-z0-9][a-z0-9._-]{0,63}`. A person omitting it builds by hand;
old model-less holds retain unknown attribution. The model is immutable
for that hold: a repeated request returning an existing hold keeps its
recorded model, even if the supplied model differs. `--model` without
`--claim`, a missing agent model or a malformed model is `invalid-argument`
(exit 2), before project access. Parser failures name `next.get`; a
missing required agent model names `next.claim`.

A hold records its ID (`HLD-` and 6 hex digits), the task, the worktree
(the SHA-256 of its root path, and the path), the actor and model, `since` and
`expires_at`, an hour later (D8), `last_heartbeat_at` (initially `since`), and,
for a resume or rework return, the hold it replaces. Whether a hold has
lapsed is read from the clock, never stored. `fid next` reads the clock once, just before it
walks the order: after refreshing the projection and, with `--claim`, under
the write lock, so a hold that lapses while it waits for the lock counts as
lapsed. A hold ends when a new hold replaces it, or its
task is completed or cancelled, parked as interrupted, or taken over. A
parked started task keeps its worktree reservation after the hold ends.
`done` arrives in 1.6. A worktree is identified by its root path, so a moved
worktree is another worktree.

The loop-facing subset is declared and versioned in
[LOOP_CONTRACT.md](LOOP_CONTRACT.md). The JSON result, with exit 0 for every kind:

```json
{"kind": "new", "worktree": "/home/mark/projects/fid",
 "hold": {"id": "HLD-4c1e9a", "task_id": "FID-9f13aa", "actor": "agent:claude-builder", "model": "claude/opus-5.5",
          "since": "2026-11-10T09:02:11Z", "expires_at": "2026-11-10T10:02:11Z", "replaces": null},
 "task": <task>,
 "epic": {"id": "FID-e3a1c0", "title": "Operating experience", "progress": {"done": 0, "total": 3}},
 "queue": [<task> + {"reason": "held-elsewhere", "after": null,
                     "working": {"hold_id": "HLD-77b2c0", "worktree": "/home/mark/projects/fid-wt2",
                                 "here": false, "actor": "agent:claude-builder",
                                 "since": "2026-11-10T08:40:12Z", "expires_at": "2026-11-10T09:40:12Z",
                                 "lapsed": false},
                     "epic_blocked_by": []}]}
```

- `kind` is `new`, `resume`, `rework`, `review` or `none`.
  A loop reads `.result.kind`,
  `.result.task.id` and `.result.hold.id`, and on failure the exit code and
  `.error.retryable`.
- `worktree` is this worktree's root path. `task` is the served task in
  `fid status`'s task shape (section 3), or null for `none`. `epic` is the
  served task's epic, with its progress as `fid status` gives it, or null
  when the task has no epic and for `none` (Mark, 12 November 2026).
- Holds in `hold`, `working` and `kept` JSON include `model` (a string,
  empty for manual or historical unknown attribution).
- `hold` is the hold `--claim` took or kept; it is null for `next.get` and
  for `none`. Its times are UTC, to the second. `event` is the
  `hold.taken` Event, or null when no hold was taken.
- `queue` is always the whole order, served task included, each entry in
  `fid status`'s task shape plus `reason` (null when nothing stops it),
  `after` (the task a sequence waits for, or null) and `working` (the
  task's open hold, or null; `here` says it is this worktree's, `lapsed`
  that it had expired) and `epic_blocked_by` (the open blockers of the
  task's epic, `[]` when there are none). After `--claim`, `working` shows
  the holds as the Event left them. Task 1.4b reuses these fields, `epic`
  included, for READY in `fid status`.

Text output is `status`'s 72-column box, headed `fid next` and today's
local date. It never names the kernel's words, `claim` and `--claim`
included, and doesn't page:

```text
╭ fid next ──────────────────────────────────────────────── 2026-11-10 ╮

  ○ FID-9f13aa  Harden login rate limiting                       to do

  Epic        ◆ FID-e3a1c0  Operating experience              0/3 done
  Priority    P1

  Passed over 1 task ahead of this one
      → fid next --explain
╰──────────────────────────────────────────────────────────────────────╯
```

- The served task's row has its mark, ID, title and state. `Epic` (with
  `N/M done`) and `Priority` follow when the task has them, and `Working`
  (`here, since 09:02, until 10:02`) while this worktree's hold on it is
  live. Times are local `HH:MM`, with the date when it isn't today.
- Without `--claim`, a task handed out after others were passed over adds
  `Passed over N tasks ahead of this one` and `→ fid next --explain`;
  otherwise the next step is `→ fid show <ID>`.
- Rework adds `Rework here: fix the review findings, then check and hand
  off again.` and `→ fid show <ID>`. Kept JSON uses mode `rework` with
  `since` set to the first returning round’s time. Show’s Working row says
  `rework in <worktree> by <actor>, since <time>`; status and the queue
  verdict say `rework in <dir>` (or `rework here`). After taking it,
  `working` describes the new build hold and `kept` is null. Repeating
  `next` now returns `resume` with that same live hold and no new Event.
- A live hold re-run adds `Already yours in this worktree: carry on with
  it`. A resume after a lapse adds `⚠ Interrupted here: work stopped at
  <HH:MM> without a park.` and `Carry on from what's in this worktree.`,
  the time being the lapsed hold's last heartbeat.
- When nothing is ready: `○ Nothing ready here:` and a count for each kind
  of reason, in the order each first appears (`1 blocked`, `1 being worked
  in fid-wt2`, `1 interrupted in fid-wt2`, `2 not next in their epic's
  sequence`, `1 in a blocked epic`, `1 in a closed epic`), each count kept
  whole when the line wraps, then `→ fid next --explain`. With no tasks at
  all, `○ Nothing ready here: no tasks to do`.
- `--explain` adds `IN ORDER`, every task in the queue with its verdict
  (shared with task 1.4b's READY): `next`, `ready`, `blocked by <ID>` (and
  `+N`), `after <ID>`, `working in <dir>`, `interrupted in <dir>`, `epic
  blocked by <ID>` and `epic closed`, `<dir>` being the worktree's
  directory name.

**Review reservations and the cap (RESET 2.1).** Review JSON adds
`review: {handoff: {id, pr, head, base}, stream, required_profiles}`. Text
names the stream, profiles, hand-off, head, base and PR URL. Without an actor,
next shows the first due review; `--claim` skips a builder with
`built-by-actor` and continues to eligible work. Independence uses build holds
only, through their whole history; a review reservation never makes someone
its builder.

Taking a review writes one `review.claimed` Event. The `hold` result has the
usual session fields with an `RCL-` ID. Review reservations live separately
from build holds, so the builder's hand-off worktree stays kept. There is at
most one open review reservation per task and per reviewer worktree. Repeating
next there returns the same review and session without another Event, ahead
of other reviews. Heartbeat takes its exact ID and actor, renews it for one
hour, and refuses a foreign actor, worktree, ended or lapsed session. A lapsed
review frees the stream; it does not keep interrupted build work. Recording
a round ends that task/stream's reservation in the same Event. Recording
without a reservation remains supported, including a courier recording an
independent review. Append, replay, import and restore validate the same
reservation rules.

The rework cap counts all `CHANGES` rounds, in both streams and across every
hand-off. It uses `rework_cap` from the latest hand-off's policy and lane,
including a raised lane. Rework is served at the cap; past it, next skips the
task as `rework-cap-reached`. Cap zero stops the first automatic rework.
The builder worktree gets `none` with `waiting: {kind: "rework-cap", task_id}`;
other worktrees continue. Status adds one `rework-cap` NEEDS YOU item:
`<task> reached its rework cap (<count> rounds asked for changes; cap <cap>)`.
Its first next action is `fid rework allow <task> --rounds <n> --actor
human:<name> --reason <text>`, followed by takeover or cancellation. Every
person's allowance adds to the frozen lane cap without resetting the count;
the accepted-work and allowance contract below gives the decision and next's
guidance fields. All reads stay offline.
If local objects cannot establish review freshness, next serves no review
for that task and status shows `review-unavailable` (or the existing
`approval-unavailable` for Assured work), preserving the board.

Refusals; 1.4a adds the operations `next.get` and `next.claim` and no codes
or reasons:

| Case | Operation | Code / exit |
|---|---|---|
| `--actor` without `--claim` | `next.get` | `invalid-argument` / 2, `argument: request` |
| `--claim` without `--actor` | `next.claim` | `invalid-argument` / 2, `argument: --actor` |
| A malformed actor | `next.claim` | `invalid-argument` / 2, `argument: --actor` |
| A `service:*` actor with `--claim` (trust by actor type, ADR-0014) | `next.claim` | `authority-denied` / 5 |
| The write lock still busy after its 5-second wait | `next.claim` | `lock-unavailable` / 4, retryable |
| Outside a repository, at an unsafe ledger location, or with no ledger | either | as every project command (section 1) |

Grammar refusals from the parser, such as an unknown flag or an argument,
name `next.get`. A `service:*` actor's human refusal names the actor alone
(`reason: authority-denied; actor_id=service:ci`); its JSON keeps
`authority-denied`'s details, the operation names included. Journal replay failures keep their codes; anything else
is `internal-error` (exit 7).

**Fixed states, task lanes and start (RESET 1.5a4; Mark, 14 November 2026).**

Tasks have `todo → in_progress → in_review → done`, plus `cancelled`.
A done task or a closed epic cannot be cancelled. The typed
`transition-not-permitted` refusal names `finished-work-stays-done` and says
that undoing merged work is a new task following the original PR's
`## Rollback` section (Mark, 2026-11-22).
Epics have `todo`, `in_progress`, `done`, `cancelled`; subtasks keep their
existing transitions. Human state labels are “to do”, “in progress”, “in
review”, “done”, “cancelled”. The retired task states and epic `wrapped`
and `closed` are removed, without data conversion. Workflow authoring and
assignment commands are removed. Creation always uses the fixed default
workflow for the issue kind.

```text
fid issue create --kind task --title <title> [--lane <routine|standard|assured>] --actor <actor> --idempotency-key <key>
fid issue update <task> --lane <routine|standard|assured> --actor <actor> --idempotency-key <key>
fid start --hold <id> --actor <actor> --idempotency-key <key>
```

Lanes apply only to tasks. Creation, issue reads and show's issue JSON carry
`lane` (`null` when unset). A task's show text has `Lane` after `Spec`, with
`Routine`, `Standard`, `Assured`, or `not set`. Create text adds ` · Standard`
(for that lane): `created issue FID-c2d4e5 (task) "Retry a rejected backup
push" · Standard`. A lane-only update reads `updated issue FID-c2d4e5 ·
Assured`. Agents may set or raise a lane, never lower it (`authority-denied`).
A person may lower a todo task's lane as a planning correction. Otherwise
lowering is refused as `transition-not-permitted`, unsatisfied
`lane-lowering-requires-waiver` and names `fid waive` after hand-off.

Start (`task.start`) validates the exact live hold under the write lock:
existence, open/unexpired status, actor, and current worktree. It needs a
task in `todo` with a lane, moves it to `in_progress`, and moves a todo
epic to `in_progress` in the same Event. It requires readable trusted
policy (`policy-unavailable` otherwise), records a changed effective policy
snapshot, requires no clean tree, and neither changes git nor records a
`when_changed` base. The latter is recomputed at checks and hand-off.

Text: `✓ Started FID-9f13aa · Standard`. JSON:

```json
{"contract_version":"fid-cli.v1","ok":true,"operation":"task.start","result":{"task_id":"FID-9f13aa","status":"in_progress","hold_id":"HLD-4c1e9a","lane":"standard"},"event":{"id":"EVT-100001","occurred_at":"2026-11-14T09:00:00Z"}}
```

With a current hold on a task already `in_progress`, start succeeds without
a new start Event: `✓ FID-9f13aa is already started · Standard`, the same
JSON result and `event: null`. A fresh retry key on this harmless repeat is
**left unbound**: while the task remains in progress, reuse returns the
same already-started response; after it moves on, reuse gets the current
state's result. No retry storage or Event is added for this no-op. The key
that recorded the original start still replays that recorded response and
Event. This no-op key rule is the orchestrator's decision under Mark's
delegation for the 1.5a series, 14 November 2026.

A missing lane refuses with `lane-not-set` (precondition, exit 4,
not retryable, subject type `issue` naming the task, exact details
`{task_id}`):

```text
error: FID-c2d4e5 has no lane
reason: lane-not-set; task_id=FID-c2d4e5
next: fid issue update FID-c2d4e5 --lane <routine|standard|assured> --actor <actor>
```

Start from `in_review` is `transition-not-permitted`, unsatisfied
`task-todo`: a changes round returns rework through 1.6f, and 2.1 caps
automatic rework. Generic transitions refuse
task start (`use-start`), hand-off (`use-hand-off`), and moving a started
task back to todo (`task-already-started`). Epics also refuse manual start.
For example:

```text
error: FID-9f13aa starts with fid start
reason: transition-not-permitted; from_state=todo, to_state=in_progress, unsatisfied=use-start
next: fid start --hold <hold> --actor <actor> --idempotency-key <key>
```

Cancellation stays person-only. Task completion belongs to `fid done`
(RESET 1.6i), with `transition-not-permitted`, unsatisfied `use-done`,
for every actor attempting it through `issue transition`. Subtasks retain
their existing behavior; the epic completion guard arrives in 1.6k.

**Hold lifecycle (RESET 1.5a2; Mark, 14 November 2026).**

```text
fid heartbeat --hold <id> --actor <actor>
fid park --hold <id> --actor <actor> --interrupted --reason "<text>" --idempotency-key <k>
fid takeover <task> --actor human:<name> --reason "<why>" [--model <family/model>] [--idempotency-key <k>]
fid issue transition <task|epic> --to cancelled --actor human:<name> --reason "<why>" [--idempotency-key <k>]
```

- Heartbeat (`hold.heartbeat`) takes no retry key. Every accepted heartbeat
  records one Event, renews expiry to an hour after it and records
  `last_heartbeat_at`; `since` stays fixed. The loop sends it every 300
  seconds. Text: `✓ Held FID-9f13aa until 10:05`. JSON result:
  `{task_id, hold_id, last_heartbeat_at, expires_at, heartbeat_interval_seconds: 300}`.
- Heartbeat and park validate the exact hold under the write lock: it must
  exist, be open, be unexpired (at expiry is too late), belong to the given
  actor and belong to this worktree. They never substitute another hold.
- The interrupted form of park (`hold.park`) requires `--interrupted
  --reason`; bare park gives
  `invalid-argument` naming both forms. The waiting form is described below.
  Text for started work: `✓ Parked FID-9f13aa · kept for this worktree`;
  unstarted: `✓ Parked FID-9f13aa · back in the queue`. JSON result:
  `{task_id, hold_id, mode: "interrupted", reason, kept_for_worktree}`.
  It ends the hold; started work resumes first here under a new linked
  hold and is skipped elsewhere as `interrupted-elsewhere`. Unstarted
  work is available to anyone. Its retry key is required; a same-key
  repeat returns the original result and Event, even after resume.
- Takeover (`task.takeover`) accepts an optional `--model` in the same
  form as next; omitted means building by hand. It becomes the new hold
  model, appears in hold reads and participates in the retry digest.
  Takeover is for a person; agents and services get
  `authority-denied`. One atomic Event ends the old hold or releases kept
  work, records the person's decision and takes a new linked hold in the
  person's current worktree. The old hold can be live, lapsed, parked or
  absent. Workflow state is preserved, and git changes stay where they
  are. Text: `✓ Took over FID-9f13aa · it's yours in this worktree`.
  JSON result: `{task_id, hold_id, previous_hold_id, decision_id}`;
  `previous_hold_id` is null with no previous hold. Closed tasks get
  `transition-not-permitted` with `unsatisfied: ["task-open"]`; epics
  and subtasks get `unsupported-kind`.
- A worktree holding or keeping a different task refuses takeover with
  `worktree-busy`: `error: this worktree already has FID-c2d4e5` and
  `reason: worktree-busy; task_id=FID-c2d4e5, worktree=/w/fid`.
  The person already holding this task here gets `already-held`:
  `error: FID-9f13aa is already yours in this worktree` and
  `reason: already-held; hold_id=HLD-7b2e10`. With a retry key, the
  same request returns the recorded response and original Event instead.
- A person's cancellation requires `--reason`, accepts an optional retry
  key and records a decision, ending any hold and releasing any kept
  worktree in the same `issue.transitioned` Event. Agents and services get
  `authority-denied`. Text: `✓ Cancelled FID-9f13aa`. The existing JSON
  result gains `decision_id` (null for other transitions). Subtasks keep
  their existing transitions and retry-key requirements.
- A person's decision is immutable: ID (`PDN-` and six hex digits), kind
  (`takeover` or `cancel` so far), actor, subject type and ID, reason,
  binding ID and time. These two kinds bind to their subject. Replay and
  import refuse decisions from agents or services.
- `working` retains an open hold's record, including a lapse, and gains
  `last_heartbeat_at`. `show`, next queue entries and status wherever it
  has `working` also carry `kept` (null otherwise). With no live hold,
  started work's lapse or interrupted park produces:

  ```json
  {"worktree": "/w/fid-wt2", "here": false, "mode": "interrupted",
   "reason": "Usage limit", "since": "2026-11-14T07:00:00Z",
   "hold_id": "HLD-77b2c0", "actor": "agent:builder"}
  ```

  Reason is null for a lapse. Since is the park time or the last
  heartbeat. An explicitly parked hold has `working: null`. Show's
  Working row reads `⚠ interrupted in <worktree> by <actor>: parked at
  20:00: Usage limit`; next reads `⚠ Interrupted here: parked at 20:00:
  Usage limit.` A lapse keeps `work stopped at 20:00 without a park`,
  now using the last heartbeat. Status and READY use
  `interrupted in <dir>`. Status reads these from its one projection
  snapshot; next and READY share one planner.

**Waiting and answers (RESET 1.5a3; Mark, 14 November 2026).**

```text
fid park --hold <id> --actor <actor> --waiting "<question>" --idempotency-key <k>
fid answer <task> "<text>" --actor human:<name> [--idempotency-key <k>]
```

Park requires exactly one of `--interrupted --reason <text>` or
`--waiting <question>`. The waiting form ends the hold and always keeps
its task for its worktree, started or not. It inserts one question with
its own ID (`QST-` and six hex digits), text, asker and time. A task has
at most one open question. Text: `✓ Parked FID-9f13aa · waiting for a
person's answer`. Its `hold.park` JSON result is
`{task_id, hold_id, mode: "waiting", question_id, question, kept_for_worktree: true}`.
A retry key is required and a same-key repeat returns its original result.

Answer binds to the task's one open question; there is no `--question`
flag. It takes no hold and changes no workflow state. It records an
`answer` decision (`PDN-` ID) whose subject is the task and whose binding
is the question. Text: `✓ Answered FID-9f13aa · it resumes in /w/fid-wt2`.
The `task.answer` JSON result has exactly
`{task_id, question_id, decision_id, answer}`. The answer and question
text are nonempty, at most 4096 bytes, like the existing free-text reasons.
A retry key is optional; with one, a repeat replays its original answer
and Event, even after resume or cancellation. Without one, a second answer
is refused as `not-waiting`. Agents and services get `authority-denied`.

While the question is open, this worktree's next result is kind `none`,
with null `task`, `hold` and `epic`, and a new `waiting` field:

```json
{"kind": "waiting-for-you", "task_id": "FID-9f13aa", "question_id": "QST-100001",
 "question": "Should retries include authentication failures?"}
```

Other worktrees skip the task as `waiting-for-you` and carry on. A
sequential epic keeps waiting for that task too. Next here says:

```text
  ○ Nothing here until you answer FID-9f13aa: Should retries include
    authentication failures?
      → fid show FID-9f13aa
```

After the answer, this worktree resumes the task first (kind `resume`),
under a new hold linked to the parked one. Next's new `answer` field is:

```json
{"question_id": "QST-100001", "question": "Should retries include authentication failures?",
 "text": "Exclude authentication failures.", "decision_id": "PDN-100001"}
```

It stays on repeated next results, including the live hold returned after
resume, until the next question replaces it. `waiting` and `answer` are
always present and null otherwise. `kept.mode` is `waiting` for a waiting
reservation, with its `question` text; it remains kept after answering
until resumed. The extra `kept.question` is omitted for interruptions.
Text's resume message is `Answer: <text>`, then `fid show <task>`.

Status adds a `waiting-for-you` NEEDS YOU item, subject `{type: "task",
id: <task>}`, with summary `<task> is waiting for you: <question>` and
next action `fid answer <task> "<answer>" --actor human:<name>`. It shows
`waiting for you` in the task's row. Show adds a `Question` row for an open
question and a `questions` JSON array (`[]` with none), in journal order:
`{id, task_id, hold_id, question, asker, asked_at, resolved_at, resolution, answer}`.
Unresolved times, resolution and answer are null. Resolution is `answered`
or `cancelled`; an answer has the shape above. Questions and answers stay
in history. Takeover retains the open question and gives the person a
hold; the person may answer. Cancellation resolves the question in the
same Event and releases its hold or kept reservation.

For every kept task, including a started task whose hold lapsed, status
checks git's local worktree registration and the worktree's common
directory. It never fetches or releases kept work. Confirmed absence or
replacement adds NEEDS YOU kind `worktree-missing`, subject the task:

```text
  ✗ FID-9f13aa is kept for a worktree that no longer exists: /w/old
      → fid takeover FID-9f13aa --actor human:<name> --reason "<why>"
```

A check that cannot finish (permissions or git errors) adds warning kind
`worktree-unverified`, with `task_id` and `worktree` besides the standard
`kind`, `summary` and `next_action`:

```text
  ⚠ Couldn't check the worktree kept for FID-9f13aa: /w/old
      → look at /w/old, then fid status
```

These warnings follow policy warnings. Waiting and missing worktree items
precede backup NEEDS YOU items, in hold/planner order respectively. Status
reads the question, work and hold state from its one projection snapshot.

**Accepted work and rework allowances (RESET 2.1b).**

```bash
fid accept <issue>... --actor human:<name> [--reason <text>]
fid rework allow <task> --rounds <n> --actor human:<name> --reason <text>
```

Both accept `--project` and `--json`; operations are `task.accept` and
`task.rework-allow`. Agent and service actors get `authority-denied`. A missing
actor/rounds/reason flag where required, invalid rounds or a closed/unknown
subject gets `invalid-argument`. No new failure code is introduced.

An open agent-created task or subtask is unaccepted until a person accepts it.
There is no date exception, including tasks created before this upgrade.
Human-created work is accepted immediately. `next` skips unaccepted tasks with
`not-accepted`; text says `not accepted`. An existing build or review session
still occupies its worktree while waiting. Unaccepted tasks do not take the
next sequential slot; dependency links still apply. Status adds one NEEDS YOU
item per open unaccepted task/subtask, kind `accept`, subject `{type: "issue",
id}`, summary `<id>, created by agent:<name>, waits for you to accept it`, next
action `fid accept <id> --actor human:<name>` or `fid issue transition <id>
--to cancelled --actor human:<name> --reason <text> --idempotency-key <key>`.
The cancellation suggestion includes the retry key required for subtasks and
accepted by tasks. Acceptance or cancellation clears the item. The
approved-epic discovered-blocker exception is deferred to 3.1's epic approval.

Acceptance records one immutable person decision of kind `accept` per issue,
bound to that task/subtask. Its optional reason defaults to `Accepted for
work`. Repeating acceptance returns its existing decision without writing.
Several IDs run in argument order in separate Events; a later refusal does
not undo earlier acceptances. JSON is `{decisions: [{issue_id, accepted: true,
decision_id, reason?}, ...]}` with no single top-level Event. Text is
`✓ Accepted <issue> · decision <decision>`. Neither command takes a retry key.

An allowance applies to an open task, with positive signed 64-bit `rounds`.
Each invocation creates a new `rework-allowance` person decision, bound to the
task, and adds its rounds to the latest hand-off's frozen lane cap. All
`CHANGES` rounds across streams and hand-offs still count. The effective cap
saturates at MaxInt64, which no representable count exceeds. At or below the
new cap, rework is served again in the reserved builder worktree; other
worktrees cannot take it. The `rework-cap` item clears when no longer capped.
Its first next action offers the allowance command, then takeover/cancellation.
JSON success is `{task_id, decision_id, rounds, reason}` and its Event. Text is
`✓ Allowed <n> more rework rounds for <task> · <reason>`.

For served work, `next` includes `rework_allowance: {decision_id, rounds,
reason, actor, at}` when any allowance exists, including on resume. The latest
decision in journal order supplies guidance; its rounds is that decision's
increment, not the effective cap. Rework text adds `Guidance: <reason>`.
Acceptance and allowance reads remain offline. Live append, replay, import
and rebuild share the same admission rules for `work.accepted` and
`rework.allowed`: exactly one complete, immutable person decision matching
the Event's actor, command, subject and time.

Finding answers, atomic deferral and canonical PR comments are implemented in
RESET 2.1c, described below. Review list stays offline.

**Hand-off and the worktree during review (RESET 1.5a7; Mark, 14 November 2026).**

```text
fid handoff --hold <id> --actor <actor> --pr <url> --idempotency-key <key>
```

All four flags are required. The PR must be a GitHub pull request URL,
`https://github.com/<owner>/<repo>/pull/<n>`. It is recorded locally, without
contacting GitHub or running checks. 1.6 verifies it with `gh`.
The exact hold must be live, owned by this actor and this worktree, and
its task must be `in_progress`. The higher of the task's lane and its computed floor must have a `checks` key,
and `[review] test_adequacy` must be configured (an empty list is valid).
The worktree must be clean, including untracked files and excluding ignored
files. fid recomputes the committed comparison base with the local
`origin/<default>` ref, without fetching. Every required check must have
its latest finished result, for this task, HEAD's tree and the current
check definition, passing and counting. A commit preserving the tree keeps
the results valid. Each selected counting result's kept Artifacts must still
exist with matching digests. This is a present file check, shared with the
done evaluator; it never changes a check result or writes an availability
observation. Hand-off never executes a check. Historical replay and an
identical hand-off retry never re-read those files.

One `task.handed-off` Event records the immutable hand-off receipt, moves
the task to `in_review`, ends its builder's hold, and keeps the task for its
worktree. The receipt records the PR, head commit, tree, base, lane,
policy snapshot, required names, definition digests, counting result IDs
and time. RESET 2.2b adds `floor: {lane, paths, hard_floor, hard_paths}`
and optional `raised_from` to the receipt. Both floor lanes are nullable;
arrays are empty when no changed paths set a floor. A higher floor raises
the task's lane in the same state update; hand-off never lowers a lane.
The public hand-off success response retains its existing fields.

Show adds the current floor in `result.floor`, null when unavailable, and
keeps the immutable floor on `result.handoff`. When hand-off raised the task,
the Lane row explains it, for example `Assured · floor: internal/policy/a.go
and 2 others · raised from Standard`. Read-only show warns when it cannot
compute the current floor and retains the recorded explanation when known.
Invalid UTF-8 Git names have a readable replacement in the floor's path
strings and a lossless base64 `changed_path_bytes` list in the receipt;
checks, reviews and done use those exact bytes. A same-key retry returns the original result and Event even
after the hold ends; changed request fields cause `idempotency-conflict`.

Text success (with no required checks, use `no checks required`):

```text
✓ Handed FID-9f13aa to review · 2 checks passed · PR #63
```

The `task.handoff` JSON result is:

```json
{"task_id":"FID-9f13aa","status":"in_review","handoff_id":"HOF-100001","pr":"https://github.com/mbell-nzl/fid/pull/63","head":"a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2","base":"9f8e7d6c5b4a9f8e7d6c5b4a9f8e7d6c5b4a9f8e","required_checks":["lint","test"],"result_ids":["CHK-100003","CHK-100004"]}
```

It has the original Event's ID and time in the ordinary success envelope.
All refusals change nothing. Preconditions and hold conflicts exit 4;
invalid flags and PR URLs use `invalid-argument`, exit 2.

| Condition | Code | Details and repair |
|---|---|---|
| Lane's `checks` key absent | `checks-not-configured` | `{lane}`; add checks to `[lanes.<lane>]`, then merge `.fid.toml` and `git fetch` |
| Test paths absent | `test-paths-not-configured` | `{}`; `add [review] test_adequacy, then merge .fid.toml and git fetch` |
| Uncommitted changes | `tree-dirty` | `{worktree}`; commit them before hand-off |
| Required checks without counting passes | `checks-unsatisfied` | `{checks, latest_result_ids}`; `fid check --hold <id> --actor <a> --idempotency-key <key>` |
| A selected counting result's kept Artifact is missing or changed | `kept-artifact-unavailable` | `{artifact_ids}`; restore the kept bytes or repeat the check with a new key |
| Comparison base unavailable | `comparison-unavailable` | `{ref, problem}`, with nullable ref, as for check |
| Task not `in_progress` | `transition-not-permitted` | `{from_state, to_state, unsatisfied}`; start the task first |
| Hold not current | `hold-not-current` | Exact hold checks and repair as above |
| Policy unreadable | `policy-unavailable` | Same policy source details and repair as other work commands |

`test-paths-not-configured` and `tree-dirty` are precondition codes,
not retryable, with the task as their issue subject. For example:

```text
error: FID-9f13aa has uncommitted changes; commit them before hand-off
reason: tree-dirty; worktree=/w/fid-wt2
```

After hand-off, kept work has mode `review`. Until a round asks for
changes, the task is done (1.6), or it is cancelled, this worktree's next
answer is kind `none`, with null hold,
task, epic and answer, and `waiting: {kind: "review", task_id: <id>}`.
The queue retains its usual contents. Other worktrees carry on; a review
reservation elsewhere uses the existing `held-elsewhere` queue reason.
The text's message is:

```text
  ○ Nothing here: FID-9f13aa is in review
      → fid show FID-9f13aa
```

Status labels it `in review`. Missing-worktree NEEDS YOU covers this
reservation too. Cancellation releases it in the same transition Event.
The loop opens a worktree per task (2.5); rework returns to this worktree
as kind `rework` (1.6f).

Show gains a Hand-off row, for example
`10:10 · PR #63 · lint, test passed on a1b2c3d`, after Checks. With no checks
it says `no checks required`. Its JSON always has `handoff`, null before
one, otherwise `{id, pr, head, base, required_checks, result_ids, at}` for
the latest receipt. A new hand-off after rework may name a different PR;
done uses the latest. Rework behavior and its public-command stale-round V8 journey belong to 1.6f.

### Review rounds (RESET 1.6d)

```text
fid review record <task> --input <file|-> --actor <recorder> --idempotency-key <key>
fid review list <task>
```

Both commands accept `--project` and `--json`. Their operations are
`review.record` and `review.list`. `review record --example` prints an
editable template outside a project; help describes the closed input
schema. The input is a single UTF-8 JSON object of at most 1 MiB. The complete
encoded Event, including retry and finding copies, must also remain below
the existing journal scanners' 1 MiB line budget; an oversized round refuses
before append and changes nothing. Unknown
or duplicate members, missing required members, wrong types and nulls
outside the nullable location are refused before project access.

| Member | Type and requirement | Meaning |
|---|---|---|
| `reviewer` | Required object | Attribution of the independent reviewer, distinct from the recorder |
| `reviewer.actor` | Required typed work actor | `agent:<id>` or `human:<id>`; never a service |
| `reviewer.model` | Optional string | `<family>/<model>`, as on build holds; omission or an empty value means unknown |
| `reviewer.effort` | Optional short token | The declared reviewer effort |
| `reviewer.session_id` | Required short token | Attributed session identity; an assertion, not authentication |
| `stream` | Required string | `primary`, or `separate` for an Assured hand-off |
| `profiles` | Required array of unique short tokens | Every profile this session reviewed; `[]` is possible for Routine work without a test trigger |
| `verdict` | Required string | `APPROVE` or `CHANGES` |
| `purpose` | Required nonblank string | What the change sets out to do, as the review read it (jig's "Claim" line; fid's human text never says "claim") |
| `findings` | Required array | `[]` for no findings; the caller supplies no finding IDs |
| `findings[].profile` | Required string | One of the session's profiles |
| `findings[].severity` | Required string | `major`, `minor`, or `nit`; majors block |
| `findings[].summary` | Required nonblank string | The problem found |
| `findings[].location` | Required object or null | A file and positive integer line, or no location |
| `findings[].location.file` | Required nonblank string | The file containing the finding |
| `findings[].location.line` | Required positive integer | Its line number |
| `findings[].evidence` | Required nonblank string | Evidence supporting the finding |
| `findings[].fix` | Required nonblank string | The requested action |
| `dispositions` | Optional array | Earlier findings assessed by this round; omission and `[]` mean none |
| `dispositions[].finding_id` | Required string | An already recorded finding of the same task |
| `dispositions[].disposition` | Required string | `resolved`, `not-applicable`, or `open` |
| `dispositions[].reason` | Required nonblank string | Evidence or reason for this disposition |
| `verification` | Required nonblank string | Commands, results, probes and limitations |
| `provenance` | Optional object | Original prompt, command and raw reviewer output |
| `provenance.prompt`, `.command`, `.output` | Optional strings | Retained unchanged for later comment rendering |

Session, effort and profile tokens match `[A-Za-z0-9][A-Za-z0-9._-]{0,127}`.
These are declared attribution values. A model label alone does not prove
independence. The recorder needs the work action set; agents and humans
may courier someone else's review, including a builder couriering an
independent review. Services get `authority-denied` (exit 5), with the
actor, operation and `review` action in the standard authority details.

Admission examines **every** build hold of the task, including ended,
lapsed and replaced holds. A reviewer who ever held it is refused. There
must be a latest hand-off; the round binds to its recorded base and head
using the local `patch-id-verbatim-per-file-v1` fingerprint, never fetching missing
objects. Every round must include the required profiles selected from that
hand-off's lane, changed paths and policy snapshot: `code` for Standard and
Assured, plus `test-adequacy` in any lane whose paths match the frozen test
paths. Unsupported changed paths conservatively trigger test adequacy.
Additional profiles are retained but add no gates.

A separate round needs an Assured hand-off and a different session from
every primary round of the task. The same reviewer actor may do both in
different sessions. Reusing a separate session for a later primary is also
refused. `CHANGES` is required exactly when a major is open after recording,
including majors raised by earlier rounds on any fingerprint, after this
round’s dispositions are applied. `resolved` and `not-applicable` close an
earlier finding; `open` or no mention leaves it open.

As of RESET 1.6f, recording `CHANGES` in either stream on an `in_review`
task moves it to `in_progress`, increments its revision, and changes its
own hand-off’s kept review reservation to rework, in the same Event as the
round. It requires that exact kept reservation; a round cannot move another
task or an unstarted or finished task. A `CHANGES` round on a task already
`in_progress` records evidence without moving it or changing any hold,
including when the builder already has the new rework hold. `APPROVE`
never moves state or changes a reservation. Person takeover, cancellation
and missing-worktree NEEDS YOU work for rework as for other kept work.
Review selection and the rework cap are built in 2.1. Posting and discovered
work are built in 2.1c, described below.

For example, this input adapts the actual 1.6b round-one jig-format review
of conflicting merge retries (its second, minor finding is omitted here):

```json
{
  "reviewer": {
    "actor": "agent:independent-reviewer",
    "model": "codex/gpt-6-astra",
    "effort": "xhigh",
    "session_id": "1.6b-review-r1"
  },
  "stream": "primary",
  "profiles": ["code", "test-adequacy"],
  "verdict": "CHANGES",
  "purpose": "Persist confirmed merges and completed CI attempts through replay, export and restore; reject invalid bindings without mutation and preserve local checks.",
  "findings": [{
    "profile": "test-adequacy",
    "severity": "major",
    "summary": "Conflicting merge retries have no regression coverage.",
    "location": {"file": "cmd/fid/merge_records_test.go", "line": 111},
    "evidence": "Disabling the retry conflict guard left the focused writer, archive and V8 tests green; an altered task or PR could return the old observation.",
    "fix": "Retry an existing key with changed task, PR and observation facts. Assert refusal and unchanged journal/projection alongside the identical retry."
  }],
  "verification": "The focused baseline passed, ending with ok github.com/mbell-nzl/fid/cmd/fid 14.338s. The disabled-conflict-guard probe was not caught. The builder's historical red runs and prior coverage/no-identity shard were not independently verified.",
  "provenance": {
    "prompt": "Review RESET 1.6b, lane Assured, with code and test-adequacy profiles.",
    "output": "## Review: CHANGES\nReviewer: Codex (gpt-6-astra), effort xhigh · Profiles: code, test-adequacy · Round 1 · Same-model fallback"
  }
}
```

The example's provenance is an excerpt; callers can retain the complete
original prompt/command/output there. fid stores the complete supplied
narrative and findings, so later rendering need not invent review text.

One `review.round-recorded` Event inserts the immutable round, findings and
any closures and the required rework effects atomically. Dispositions are
retained in the round; terminal dispositions also insert immutable
`finding_closures` rows.
It records a `RND-<six-hex>` ID, task/hand-off IDs, the binding
`{algorithm, fingerprint, base, head}`, recorder and UTC time. Stable finding
IDs append their one-based position to the round ID, such as
`RND-a1b2c3-1`. JSON success is the standard envelope with result `{round}`
and the original Event. Same-key retries return the same result and Event
before consulting current hand-offs or git; changed input, task or recorder
gets `idempotency-conflict` (exit 4) naming the differing public fields.
Corrupt stored responses are refused and never rendered as a review.
Text success is `✓ Recorded review RND-a1b2c3 · primary · CHANGES`.

`task.default-conflict` is the distinct Event-v2 family returning approved
work after a local merge conflict. Its work actor runs `next.claim`, with a
`default-conflict` aggregate and sequence 1, no retry result or repository
payload. Three ordered Effects insert an immutable `default_conflicts`
receipt, move the task to `in_progress` at its next revision, and set its
kept hand-off hold's `rework_at`. Evidence names the default ref/commit,
latest hand-off/head, conflicting paths and the observed fingerprint.
Admission verifies shape and ownership on append, replay and import; only
the live planner tests git. The linked `hold.taken` then uses normal rework.
An interruption between the Events preserves both the reservation and cause.

`review list` returns `{task_id, rounds: [...]}` in journal order and no
Event. Each round contains its complete recorded fields, plus derived
`other_model` and `current` booleans. Each finding adds `open`; a closed
finding also adds `closure: {round_id?, decision_id?, disposition, reason}`,
with exactly one authority ID. Other-model means the reviewer's declared family differs
from every model-bearing build hold of that task; a manual hold adds no
family, and an unknown reviewer never counts as other-model. Two different
models within one family are the same-family case.

A round is current exactly when its algorithm and fingerprint equal the
latest hand-off's; older rounds are never rewritten. A recorded round
covering the latest hand-off's exact base/head provides its already-computed
fingerprint, so imported/restored history is readable without those git
objects. A later hand-off with new endpoints needs a local computation;
missing objects refuse the freshness read. A clean rebase with an equal
fingerprint keeps the earlier round current; edits stale it without closing
its findings. Reads never use the network. Human output uses the 72-column
box and escaped, wrapped text, listing stream, reviewer/model/session,
profiles, verdict, other-model, fingerprint, freshness and each finding’s
open/closed state. Closed findings name the round or decision, disposition
and reason.

Semantic admission and fingerprint failures use `invalid-argument` (exit 2)
with `{argument:"input", value:null,
allowed:["an independent review of the latest hand-off"]}` and a bounded
specific message. Structural input failures name the exact field instead.
The common repository/project, lock and corrupt-journal refusal contracts
still apply. Replay and import validate model/fingerprint forms, exact
hand-off endpoints, hold history, profile/session/verdict rules and retry
bindings, without git or network access. Retired `review.recorded` and its
old tables remain refused.

### Finding closure (RESET 1.6e)

A round may dispose only of an already recorded finding of its own task.
Journal order establishes “later”: caller-supplied times never establish
precedence. A round cannot dispose of its own findings. The existing
reviewer independence rule applies even to ended build holds. The later
round may bind to a newer fingerprint; resolving a finding after a fix is
normal. Duplicate finding IDs, unknown dispositions and blank reasons
refuse. Minors and nits may be disposed too, without affecting the major
gate. A closed finding stays closed: another disposition of it refuses.

For example, a later independent session can submit:

```json
{
  "reviewer": {"actor": "agent:reviewer", "model": "claude/opus-5.5", "session_id": "review-2"},
  "stream": "primary",
  "profiles": ["code", "test-adequacy"],
  "verdict": "APPROVE",
  "purpose": "Persist confirmed merges and CI attempts.",
  "findings": [],
  "dispositions": [
    {"finding_id": "RND-a1b2c3-1", "disposition": "resolved", "reason": "The new retry regression catches the original fault and passes with the fix."}
  ],
  "verification": "Focused regression and required checks passed."
}
```

This passes only when that disposition leaves no major open. A builder’s
answer never closes a major; answers and deferrals are described below. A
previous APPROVE cannot hide a later CHANGES. Raw round verdicts
remain immutable even when a person later closes the last major.
`reviewrecord.OpenMajors` is the shared task-wide reader for the later done
evaluator; it never filters findings by fingerprint or review stream.

```text
fid finding close <finding> --actor human:<name> --reason <text> [--idempotency-key <key>]
```

This accepts `--project` and `--json`, operation `finding.close`. Only a
person may close a finding; agents and services receive `authority-denied`
(exit 5). No terminal is required. The one person’s-decision record has
kind `close-finding`, subject type `task`, subject ID the finding’s task,
and binding ID the finding itself. Existing decision kinds keep their
subject rules. One `finding.closed` Event inserts the decision and closure
atomically. The closure’s disposition is `closed`: the person need not
assert that the defect was fixed. The full encoded Event must remain below
the journal’s 1 MiB line budget, including retry copies.

JSON success contains `{task_id, finding_id, decision_id, disposition,
reason}` and the original Event. Text success is:

```text
✓ Closed RND-a1b2c3-1 · decision PDN-d4e5f6 · Accepted risk
```

An identical keyed retry returns the original decision and Event; a changed
finding, actor or reason with that key is `idempotency-conflict` (exit 4).
Corrupt stored retry results refuse. A second close without the original
retry key, an unknown finding or a blank reason is `invalid-argument`
(exit 2), changing nothing. Common repository/project, lock and journal
refusals apply. Both catalogs register the new family; live, replay and
import validate authority, bindings and single closure without git or
network. Migration 31 adds the closure projection and decision kind;
projection cache version 13 rebuilds older caught-up projections.

### Finding responses and delivery (RESET 2.1c)

```text
fid finding answer <finding> --actor <a> --text <why>
fid finding defer <finding> --actor <a> --title <text> [--epic <epic>] [--lane <lane>]
fid review post <round>
```

All accept `--project` and `--json`, with operations `finding.answer`,
`finding.defer` and `review.post`. An open finding can be answered by any actor
that held its task to build it, or by a person. Blank answers refuse. An answer
never closes a finding, including a minor or nit. Review list exposes each
finding’s `responses` in journal order, and text names the answer’s actor and
reasoning. Served review work includes earlier `review.findings`, with closure
state and responses, for the next reviewer.

Deferral is for minor/nit findings, with the same actor rule. It creates a
`todo` task, inheriting the reviewed task’s epic unless overridden. The parent
must be an epic. Lane defaults to Standard; a requested Assured lane stays
higher, and a known finding location raises it to the frozen hand-off policy’s
path floor. Routine and unknown lanes refuse. The later check/hand-off computes
its normal floor from actual changed paths. The acceptance criterion includes
finding ID, summary, evidence and requested fix. One Event includes the task,
criterion, `discovered_from` link and inverse, source revision bump and response.
In a sequential epic it appends after the highest existing task sequence;
other scheduling leaves sequence unset. Admission checks that allocation too.
The Event actor supplies the new task’s creator: an agent’s discovery waits for
`fid accept`, as ordinary agent-created work does. A deferral’s derived closure
has disposition `deferred`, `response_id` and `deferred_task_id`, and it prevents
later answers or closures. The reviewed task’s state is unchanged.

Both mutations return `{response}` and the Event. A response has `id` (`FRS-`),
`finding_id`, `task_id`, latest `handoff_id`, `actor`, `text`, `at`, and optional
`deferred_task: {id, title, epic_id, lane}`. Answers retain supplied text;
deferrals retain the generated criterion. Text success is `✓ Answered <finding>
· <text>` or `✓ Deferred <finding> → <task> · <title>`. These commands have no
request retry key. Unknown/closed findings, invalid actors, blank text/titles,
invalid parent/lane and major deferrals use `invalid-argument` (exit 2), changing
nothing. Their error details name argument `finding`, with allowed input
`an open finding and a builder or person`. Common repository, journal and
lock refusals apply.

After recording, fid attempts a comment rendered from its immutable round or
response. Round comments follow jig’s review layout: verdict, reviewer actor,
model/effort/session, stream, profiles, round, lane, purpose, every finding’s
severity and action, location, evidence and requested fix, earlier dispositions
and verification. Answer/deferral comments name the finding, actor, reasoning
and discovered task where applicable. No separately supplied comment body is
accepted. The latest hand-off at **recording** fixes the destination PR so an
interrupted delivery cannot be redirected by a later hand-off.

`review.posted` retains the comment URL, rendering version and canonical
body digest. New posts use version 2: bold labels and code spans for identifiers,
the lane and finding locations; reviewer-written Markdown passes through.
Version 1 remains byte-for-byte available to verify old receipts. Admission uses
the receipt's recorded version and refuses unknown versions. Schema 42 permits
versions 1 and 2; schema 41 permits only version 1. Posting is
serialized under the journal lock. A stable record/digest HTML marker lets fid
recover a comment after a lost local receipt: it reads bounded, paginated PR
comments before calling `gh pr comment`. An edited matching marker refuses a
second post. Reconciliation checks both supported canonical renderings; a
recovered comment retains its original render version, digest and URL in the
receipt. A receipt makes subsequent retries entirely local. `review post`
retries a round plus all responses to its findings, returning `{posted: [ids],
pending: [ids]}` without a single Event in the response (each successful post
has its own receipt Event). Text prints each posted or pending ID. An unknown
round is `invalid-argument`, with argument `round` and allowed input
`a recorded review round`; transport failures return pending IDs.
If committing a receipt fails, delivery stops with that record and the remaining
records pending. The next retry refreshes the projection before posting further
receipts, recovering a receipt appended before its projection failed.

Failed delivery never rolls back a canonical recording or blocks `done`.
Status derives `review-post-pending` warnings, naming each round, its pending
record IDs and `fid review post <round>`. These are warnings, never NEEDS YOU.
Status/show/next/review-list use no network. Live append, replay and import
validate response authority, single closure, exact deferral rows, creator
attribution and receipts against the same canonical evidence. New records and
receipts survive export, rebuild and restore.

### Retained Artifact verification (RESET 1.5c2)

```text
fid artifact verify --input <path|-> --actor <actor> --idempotency-key <key>
```

Input is exactly `{"artifact_id":"ART-..."}`. Help describes this field;
`--example` prints an editable JSON object outside a project. Structural
errors identify the input field before workspace access. The typed actor
is `agent:<id>` or `human:<id>`; malformed actors give `invalid-argument`
(exit 2), and `service:<id>` gives `authority-denied` (exit 5), with the
existing actor-type authority details for `artifact.verify`.

Verification reads the retained file from the content-addressed Artifact
store populated by lean checks. Success text is `<id>\tavailable|missing`.
The JSON operation is `artifact.verify`, with result
`{artifact, changed, idempotency_key, replayed}` and the recorded Event.
`artifact` has exactly nine fields: `id`, `digest`, `media_type`,
`size_bytes`, `retention_class`, `availability`, `location_policy`,
`created_at`, `expires_at` (nullable). It emits
`artifact.availability-observed`, even when availability is unchanged.
A same-key retry replays the original metadata and Event without another
observation or write; changed public fields produce `idempotency-conflict`.

An unknown Artifact keeps its existing `invalid-argument` failure (exit 2),
null subject and details `{argument:"input", value:null,
allowed:["accepted gate registration input"]}`. The legacy allowed-text is
retained with that failure contract. The message names the missing Artifact.
Availability updates never rewrite immutable check results or emit kernel
readiness records.


### `fid waive`

```text
fid waive <task> --lane <routine|standard> --reason <text> --actor human:<id> [--idempotency-key <key>] [--json]
```

The operation is `task.waive`; it accepts `--project`. As with approve, both
standard input and standard error must be terminals, including on retries;
stdout may be redirected for JSON. Agent and service actors receive
`authority-denied`. A reason is required. The task must be in review with a
latest hand-off. The target lane must be strictly below that hand-off's lane
and never below its hard floor. `waiver-not-in-review`,
`waiver-handoff-required`, `waiver-not-lower` and `waiver-hard-floor` are
non-retryable precondition refusals with a `paths` array (the hard-floor paths
when applicable, otherwise empty). Other input and terminal refusals use
`invalid-argument`, as approve does.

The confirmation on stderr shows the task, hand-off lane, floor and setting
paths, hard floor and setting paths, and new lane. Only typing the exact task
ID followed by Enter confirms. EOF or another answer cancels successfully
without writing a decision or retry result. The immutable `task.waived` Event
records a `person_decisions` entry of kind `waive`. Its `binding_json` contains
`lane`, `algorithm`, `fingerprint`, `handoff_id`, `base` and `head`;
`binding_id` is the fingerprint. Admission verifies state, reason, person,
lower lane, hard floor, latest hand-off endpoints and decision attribution
on live append, replay and import, with no git or terminal needed on replay.

The latest decision matching the latest hand-off's algorithm and fingerprint
sets done's effective lane. Unchanged re-hand-offs retain it; an edit makes it
stale. It changes required checks and review profiles and the Assured gates
(separate review, approval, rollback). Other gates, including open majors,
merge and CI, still apply. Neither task lane nor hand-off changes. Review
scheduling and admission use the same waived profiles. The JSON success
contains `decision` and `occurred_at`, with the usual Event metadata. Text
says `✓ Waived <task> · decision <decision>`. An identical retry returns the
original response and Event without prompting or needing git objects; a
changed task, actor, reason or target lane conflicts in the project-scoped
`task.waive` idempotency scope.

### Approval (RESET 1.6g)

```text
fid approve <task> --actor human:<name> --reason <text> [--idempotency-key <key>] [--json]
```

The operation is `task.approve`. It accepts `--project`. Only a person may
approve: agent and service actors receive `authority-denied`, regardless of
their terminal. The reason must contain non-whitespace text. Both standard
input and standard error must be terminals, including with `--json` and on a
keyed retry. A piped human invocation is refused. Standard output may be a
pipe: the prompt is on standard error and the success result is on standard
output, so the prompt never mixes with JSON.

The subject must be a task in `in_review`, with a latest hand-off. fid
computes `patch-id-verbatim-per-file-v1` from that receipt's recorded base
and head using local objects, never fetching. Missing objects refuse with
a repair to restore those objects locally and retry. It does not bind to
the current checkout, a new merge-base, the head SHA or the hand-off ID.

Before asking, the prompt shows the task and lane at hand-off, the full
base..head diffstat with pinned rendering, each required check's recorded
counting result and ID, each review stream's latest round with its verdict,
reviewer, model and current/stale state, every open major (including stale
rounds' findings), and earlier approvals with their freshness and reasons.
Missing streams are explicit. Assured work shows both primary and separate
streams. `no other-model review` appears when no current review's known
model family differs from every model that built the task. A historical
agent build with an unknown model cannot establish that difference; a
person without a model built by hand.

The person types the exact task ID, followed by Enter. A different line or
end of input cancels with exit 0, a cancellation message on standard error,
no result and no Event. There is no `--yes` flag. Missing reviews and open
majors do not prevent approval; `done` will enforce its gates separately.
The ledger lock covers the preview, confirmation and append.

The success result contains `decision` and `occurred_at`. The decision has
`decision_id`, `task_id`, `actor`, `reason`, `at` and a `binding` with
`algorithm`, `fingerprint`, `handoff_id`, `base`, and `head`. The ordinary
success envelope names the approval Event and its time. Text says
`✓ Approved <task> · decision <id>`.

Approval is an immutable person's decision of kind `approve`, subject the
task. Its binding ID is the fingerprint; its structured binding retains the
algorithm and explanatory hand-off endpoints. Replay/import validate actor,
reason, task state at that journal position, fingerprint form and the latest
hand-off's task/base/head without git or a terminal. The same admission runs
before live append. Other person's-decision kinds keep their binding rules.
The full Event must fit below the journal's existing 1 MiB line budget.

A decision is current when its algorithm and fingerprint match the latest
hand-off's change, even after a clean rebase with new endpoints; otherwise
it is stale. `approvalrecord.Current` supplies the shared current/stale/missing
reader for done, status and show. Any matching immutable approval is evidence
of the current change; a new decision after edits preserves the old one.
The prompt always shows the full diffstat when a new decision is requested.
Status and show present that state as described in section 3 (RESET 2.3).

An optional retry key covers actor, task and reason in the project's
`task.approve` scope. An identical retry returns the original decision,
Event and time without re-prompting, recomputing git evidence or requiring
the task still to be in review. It still requires a person on terminals.
Changed inputs under the same key receive `idempotency-conflict`; a new
approval after an edit needs a fresh key. Semantic, fingerprint and terminal
refusals use `invalid-argument` (exit 2) and change nothing. Existing project,
lock and journal refusal envelopes still apply.


### Completion and reconciliation (RESET 1.6i–1.6k)

```text
fid done <task> --actor <actor> [--idempotency-key <key>] [--json]
fid done <epic> --actor human:<name> --reason <text> [--idempotency-key <key>] [--json]
fid done --actor <actor> [--json]
```

The operation is `task.done`, accepting `--project`. Agents and people may
finish a task; services get `authority-denied`. No live hold is required.
The task must be `in_review` with its latest hand-off's review reservation.
An epic takes the person-only path below. Subtasks are `invalid-argument`
(exit 2) and use their plain transitions. With no task, fid reconciles the current in-review work
as described below.

Each invocation observes GitHub once, outside the ledger write lock. The
PR comes from the latest hand-off, the repository from `origin`, and the
default branch from local `origin/HEAD`. fid never fetches, polls, merges,
or starts background work. Under the write lock it refreshes the projection
and compares the original task/evidence snapshot and its hold/question
history. A takeover, cancellation, lane change, newer hand-off or changed
evidence refuses as `done-observation-stale` before any observation is
recorded. A takeover already present when the call starts needs a new
hand-off too.

A confirmed merge whose head tree matches the hand-off is recorded even
while CI is pending. Every completed CI attempt is recorded through the
merge/CI writers. The single done evaluator then applies the latest
hand-off's frozen lane and policy. Missing or pending CI says "merged,
waiting for CI" and keeps the task in review. GitHub failure says "GitHub
confirmation unavailable", or "merged, waiting for confirmation" when
there is a recorded merge for this hand-off, including one independently confirmed earlier
in this invocation. No failed or waiting result binds a retry key.

Successful text is `✓ Completed FID-9f13aa · merged with CI green`. JSON
has `result: {status: "done", completion: <record>}` and the completion
Event's ID and time. The completion names its task, latest hand-off, merge,
all completed CI results observed by this call, the selected counting CI
results, review rounds and approval. It also records the locally computed
fingerprint, the Rollback observation, other-model disclosure, actor and
time. One `task.finished` Event inserts that immutable record, moves the
task to `done`, releases its review reservation, and resolves an open
question as `finished-before-answer`, without an answer or cancellation.
Replay/import validate the same recorded gates without network, git objects
or present Artifact files. An identical keyed retry returns the original
response/Event before any GitHub or git observation; changed task or actor
under that key is `idempotency-conflict`.

All `done-*` refusals below are preconditions (exit 4), with an Issue subject
and exact details `{gate, reason, record_ids}`. `reason` retains the done
evaluator's specific reason even when related repairs share a public code.

| Code | Evaluator reasons | Retryable |
|---|---|---|
| `done-task-not-in-review` | `task-not-in-review` | no |
| `done-handoff-required` | `handoff-missing`, `handoff-mismatch`, `lane-changed`, `pr-mismatch`, `repository-mismatch`, `branch-mismatch`, `merge-invalid`, `head-tree-mismatch` | no |
| `done-evidence-required` | `fingerprint-unavailable` (local computation), `prerequisite-missing`, `policy-missing`, `policy-mismatch`, `checks-not-configured`, `test-paths-not-configured`, `artifacts-not-verified`, `pr-unavailable`, `merge-missing`, `merge-mismatch`, `ci-mismatch`, `ci-unrecorded` | no |
| `done-checks-unsatisfied` | `checks-unsatisfied` | no |
| `done-kept-artifact-unavailable` | `kept-artifact-unavailable` | no |
| `done-review-required` | primary `review-missing`, `review-stale`, `profile-missing` | no |
| `done-separate-review-required` | separate `review-missing`, `review-stale`, `profile-missing`, `separate-session-required` | no |
| `done-open-major` | `open-major` | no |
| `done-approval-missing` | `approval-missing` | no |
| `done-approval-stale` | `approval-stale` | no |
| `done-rollback-missing` | `rollback-missing` | no |
| `done-not-merged` | `not-merged` | yes |
| `done-ci-waiting` | `ci-waiting` | yes |
| `done-ci-not-green` | `ci-not-green` | yes |
| `done-observation-stale` | command revalidation failed (`observation-stale`) | yes |

A missing task is `issue-not-found`. `github-unavailable` is a retryable
precondition (exit 4), Issue subject, exact details `{reason, merged}`;
`reason` comes from the GitHub adapter's unavailable reason. Local policy
failures retain `policy-unavailable`. Every refusal supplies a repair in
human text. A hard evidence failure takes precedence over waiting, with
gates considered in the evaluator's fixed order.

**Closing an epic (RESET 1.6k).** The explicit Issue's kind dispatches to
this local path under the same `task.done` operation. Only `human:*`
actors may close an epic; agents and services get `authority-denied`
(exit 5). A missing or blank `--reason` is `invalid-argument` (exit 2),
naming `--reason`. `--reason` is accepted only with an explicit epic; using it for a task or
reconciliation is invalid usage. No terminal or GitHub access is required. The epic must
be `todo` or `in_progress`, and every direct task must be `done` or
`cancelled`; an empty epic may close. Subtasks belong to their tasks and
do not affect eligibility. All tasks are re-read under the ledger write
lock, so a new or reopened task defeats a stale read.

Text says `✓ Closed <epic> · <reason>`. JSON returns
`{status: "done", revision: <new revision>, decision: {id, kind: "close-epic",
actor, subject_type: "epic", subject_id, reason, binding_id: <epic>,
decided_at}}`, plus the close Event. One `epic.closed` Event inserts the
immutable person's decision and moves that epic to `done`, with one
revision bump and an optional retry result. A keyed retry returns the
original response/Event unchanged; a changed request under that key
conflicts. An unkeyed repeat refuses a closed epic. There is no generic
`todo → done` edge: this command alone owns that move.

`done-epic-not-open` and `done-epic-open-tasks` are preconditions (exit 4),
not retryable, with an Issue subject and exact
`{gate: "epic", reason: "epic-not-open" | "open-tasks", record_ids: [...]}`.
The latter lists all open direct task IDs in ID order; the former has an
empty list when no task is open. Both point to `fid show <epic>`.
Generic `issue transition --to done` refuses every actor with
`transition-not-permitted`, `use-done`, and the person's `fid done` action.
Person cancellation is unchanged. Epic validations/plan approval belong
to RESET 3.1; reconciliation below continues to select tasks only.

**Reconciliation.** With no task argument, the same `task.done` operation
selects every task currently `in_review` with a hand-off, ordered by task ID.
It runs the exact single-task observation, locked revalidation, recording,
evaluation and completion path once per selected task. A task changing during
observation is refused for that task; the other tasks continue. Tasks completed
before the next invocation are excluded. Identical merge and completed CI
observations reuse their natural keys and append nothing.

All per-task outcomes go to stdout and the command exits 0, including when
none complete. Command failures (such as no ledger, unavailable write lock,
policy failure or invalid usage/actor) retain their ordinary nonzero exit
contract. No task argument still requires `--actor`. `--idempotency-key` is
invalid without a task (exit 2); per-task completion and observation keys
already make repetition harmless. The live and keyed-leaf inventories keep
one `done` leaf; only its explicit-task form is keyed.

Text prints one task per line, with `Completed`, `Waiting`, `Refused` or
`GitHub unavailable`, its ID, and the reason/repair. A refusal includes its
gate. No eligible tasks prints `No tasks in review to reconcile.` JSON has
`event: null` and `result: {outcomes: [...]}` (an empty array when none).
Each outcome has `task_id`, `status` and exactly one of `completion` or
`refusal`. Status is `completed`, `waiting`, `refused` or
`github-unavailable`. Waiting covers not merged and pending/missing CI;
failed CI is refused. A completion is the ordinary immutable completion
record. A refusal has `{code, message, next_action, details}`, preserving the
single-task code and gate/reason/record IDs (or GitHub reason/merged fields).

**Offline reads.** `status` and `show` never contact GitHub. For an in-review
task without a completion, only a merge bound to the latest hand-off counts.
When its recorded CI results are not yet green, status's existing IN PROGRESS
row says `merged, waiting for CI`; a selected completed failure says
`merged, CI not green`. Show adds the same text in a `Merge` row after
`Hand-off`. Selection uses the same latest-attempt, push-event, exact-commit
and skipped-result rule as `done`, applied only to durable completed receipts.
A merge without completed results is waiting. Rework, cancellation, completion
and a newer hand-off without a merge remove that view.

Both JSON answers add `merge: {merge_commit, ci}` only when that offline
merge exists: on each in-progress task in status, and at result level in
show. `ci` is `waiting`, `not-green` or `green`. When the recorded results
are all green but another gate has blocked completion, text keeps its usual
`in review` wording. Green receipts alone are never fresh confirmation that
all current GitHub CI is green: pending attempts are not durable. Only
`done` confirms GitHub; it continues to say `merged, waiting for confirmation`
when GitHub is unavailable after this hand-off's merge has been recorded.

## 4. Common Mutation Behavior

Every mutation:

1. resolves the project and the authority required by that command;
2. validates contract/schema, Actor provenance, authority, idempotency, expected
   revision, references, and domain preconditions;
3. computes one complete atomic Event/effect set;
4. holds the write lock across catch-up, append, fsync, projection transaction,
   and metadata update; and
5. returns the entity revision and Event identity.

For a keyed mutation, the request digest covers the closed input object as
written, plus the operation, project ID and resolved Actor components already
used by its family. Optional members enter the digest only when supplied. Values
observed by fid, including repository state and Source digests, do not enter
the request digest. A same-key replay returns the original result before any
observation. Each keyed mutation also records a `sha256:` digest of the
canonical JSON value for every top-level public input name and `--actor`.
Flag names include `--`; positional names follow the leaf's `Use` line.
An `idempotency-conflict` JSON details object contains the operation scope,
key, stored and supplied request digests, and `differing_fields`: the sorted,
unique, nonempty names whose field digest or presence differs. The conflict
appends no Event.

Validation failure appends nothing. A projection failure after append is a
detectable catch-up state and never a reason to append a duplicate Event.

### Checks (RESET 1.5a5–1.5a6)

`fid check <name>... --hold <id> --actor <a> --idempotency-key <k>` executes
trusted `[checks.<name>]` definitions at the worktree root, in name order,
dropping duplicates. With no names it selects the required checks at the higher
of the task's lane and the committed changed-path floor, without changing the
task's lane,
filtering `when_changed` against committed files since the freshly computed
merge-base with local `origin/<default>`. Deletions and both rename paths
count; absent or empty `when_changed` always requires the check. Unknown names
are refused before execution, with the unknown name as `value` and the sorted
defined names in `allowed`. A named check executes even if the lane does not
require it. Explicit names do not need a comparison. An unavailable base
refuses as `comparison-unavailable`; fid never fetches.

The task must be `in_progress`, with a lane and this actor's exact current hold
in this worktree. Missing lane checks refuse as `checks-not-configured`; an
explicit empty set permits named checks. Unreadable policy gives
`policy-unavailable`. A changed effective digest records a policy snapshot.
Missing `[review] test_adequacy` does not stop check. An empty required set
records a completed request with no starts or results; text says
`✓ No checks required for Standard` (using the lane selected for this check), JSON has
`{task_id, checks: []}` and the completion Event, and the exit is 0. Retry
returns that same recorded response even if policy or the comparison changes.
When a floor raises the selected lane, the check start and stored retry
response retain the comparison's `changed_path_bytes` as base64 bytes, so journal
admission can recompute that selection even for unsupported Git names. The
same rule admits a configured higher lane when the planned lane is not set
up. Explicitly named checks keep the existing lane and need no comparison.
The start is durable before spawning. The write lock is released while the
child executes; process identity is recorded as soon as known. Every five
minutes the command renews its exact hold, one Event each. A refused renewal
kills the process group, records `ownership-lost`, starts no later names, and
returns `hold-not-current`. Timeout also kills the process group, without the
old executor's ten-minute cap.

On a terminal, output streams with controls escaped (line feeds retained).
`--json` captures output and emits exactly one final envelope:

```json
{"contract_version":"fid-cli.v1","ok":true,"operation":"check.execute","result":{"task_id":"FID-9f13aa","checks":[{"id":"CHK-100001","name":"test","outcome":"passed","exit_code":0,"counts":true,"reasons":[],"artifact_ids":[]}]},"event":{"id":"EVT-100004","occurred_at":"2026-11-14T09:05:00Z"}}
```

Human output announces `● Checking test…`, then `✓ test passed · counts for
hand-off`, `✗ test failed`, or `⚠ test passed, but doesn't count: the tree had
uncommitted changes`. Other reasons read `the tree changed`, `ownership was
lost` and `a listed file wasn't kept`. Each retained file adds `  Kept <path>`.
Ordinary unsuccessful checks do not stop later names. All passing and counting
returns 0; any unsuccessful outcome returns 8; all passed but some do not count
returns 4. Failure envelopes have no `result`, and name the relevant check IDs.

`show` adds `result.checks` (an array, empty when absent), each
`{name, start, process, result, required, satisfied}`, in name order. It selects
each name's latest finished result on HEAD's tree with the current definition,
or its latest historical result if none matches, plus every unfinished start.
A required check without a task record has null start, process and result.
Start contains ID, task, hold, actor, name, resolved
definition/digest, environment-value digest, policy snapshot sequence, time,
HEAD/tree/dirty observation and the owning fid process. Process identity has
machine, separate boot ID (empty on other Unix systems), PID and start identity.
Linux machine identity uses `/etc/machine-id`, falling back to the hostname.
Result contains finish time, outcome (`passed`,
`failed`, `timed-out`, `cancelled`, `spawn-failed`, `signaled`), nullable exit
code/signal, second tree observation, bounded output digest/excerpt/total byte
count/truncation, resolved executable/digest/platform/fid version, kept and
refused files, `counts` and `reasons`. Qualification requires success, a clean
unchanged tree, current ownership at completion and every listed file kept.
Reasons are `tree-dirty`, `tree-changed`, `ownership-lost`, `artifact-not-kept`.

Show's text adds, for example:

```text
  Checks      ✓ lint passed at 10:05 · required · counts on HEAD
              ○ test required · not checked yet
```

`satisfied` describes the check's counting pass on HEAD, independently of
whether the lane requires it. A newer unfinished start does not erase a
finished pass. It is null if policy or the current definition is unknown.
`required` is null when policy is unavailable or the lane's `checks` key is
missing. If the base cannot be found, only conditional lane checks have null
required; unconditional lane checks remain true. Checks outside a configured
lane's list are false. Unknown is never represented as false.

`result.checks_base` is the base commit used, null when unavailable or not
needed (no conditional lane checks). `result.warnings` is always an array,
empty when healthy. Show never refuses because Checks cannot be evaluated:
it uses status's warning shape (`kind`, `summary`, `next_action` and that
kind's details). Policy-unavailable carries nullable ref and commit plus
problem; checks-not-configured carries the lowercase lanes list;
comparison-unavailable carries nullable ref and problem. Text adds one `⚠`
line under Checks, wrapping as needed, explaining why fid cannot tell. Show
never writes a policy snapshot.

A same-key repeat returns the original response without spawning. A keyed
request cut off before completion gives `check-incomplete` without a result.
Before a fresh check, unfinished starts on the task are inspected by machine,
boot, PID and start identity. On the same machine, an earlier Linux boot proves
the process stopped. A different or uninspectable machine still refuses.
A stopped process permits execution, with human warning
`⚠ CHK-100001 never finished; its process has stopped`; a live one refuses as
`check-still-running`, suggesting waiting or stopping its PID and trying again.
The old start remains unfinished in either case.

| Code | Category / exit | Exact detail keys |
|---|---|---|
| `checks-not-configured` | precondition / 4 | `lane` (string) |
| `comparison-unavailable` | precondition / 4 | `ref` (nullable string), `problem` (string) |
| `checks-unsatisfied` | precondition / 4 | `checks`, `latest_result_ids` (string arrays; used by hand-off from 1.5a7) |
| `check-unsuccessful` | gate-unsuccessful / 8 | `check_ids`, `outcomes` (string arrays, unsuccessful checks only) |
| `check-not-qualified` | precondition / 4 | `check_ids`, `reasons` (string arrays) |
| `check-still-running` | conflict / 4 | `check_id` (string), `pid` (integer) |
| `check-incomplete` | conflict / 4 | `check_ids` (string array) |

All are non-retryable and have an issue subject naming the task.
`checks-not-configured` suggests `add checks to [lanes.<lane>], then merge
.fid.toml and git fetch`. `check-still-running` suggests `wait for it to finish,
or stop pid <pid>, then check again`.

The shared required-check evaluator chooses the latest finished result in
journal order for each task, tree and definition digest. A later failure
replaces an earlier pass; a newer unfinished start does not. Other tasks'
results never count. Hand-off and done use this same evaluator with their
respective required sets and policy definitions. For example:

```text
error: test failed on this tree after passing; check it again before hand-off
reason: checks-unsatisfied; checks=test, latest_result_ids=CHK-100002
```

## 5. Typed Failure Codes

The current closed failure-code set is:

```text
already-held
archive-invalid
authority-denied
capability-unavailable
check-incomplete
check-not-qualified
check-still-running
check-unsuccessful
checks-not-configured
checks-unsatisfied
comparison-unavailable
dependency-cycle
duplicate-entity
hierarchy-invalid
hold-not-current
idempotency-conflict
import-target-conflict
internal-error
invalid-argument
issue-not-found
journal-replay-stopped
label-not-found
lane-not-set
link-not-found
lock-unavailable
not-waiting
operation-unavailable
policy-unavailable
project-already-initialized
project-not-found
projection-divergence
repository-unavailable
result-limit-exceeded
state-conflict
test-paths-not-configured
transition-not-permitted
tree-dirty
unsupported-kind
unsupported-relation
worktree-busy
```

Kernel-only failure codes are withdrawn with their last command producer in
RESET 1.5c2. `idempotency-conflict` remains defined and covered by process
pairs, including its existing supplemental inventory entry.

RESET 1.5a2 adds these closed conflict codes, exit 4, all
`retryable: false`, with subject type `issue` naming the task:

| Code | Exact detail keys |
|---|---|
| `hold-not-current` | `hold_id`, `reason` (`not-found`, `ended`, `expired`, `actor-mismatch`, `worktree-mismatch`) |
| `worktree-busy` | `task_id`, `worktree` |
| `already-held` | `hold_id` |
| `not-waiting` (1.5a3) | `task_id` |

`not-waiting` text is `error: FID-77aa01 isn't waiting for an answer`,
then `reason: not-waiting; task_id=FID-77aa01`. It shares the conflict
category, exit 4, non-retryable status and task subject above.

An unknown hold cannot name its task, so its subject's ID is null. A
`hold-not-current` human failure adds `next: fid next`. The old session's
heartbeat and park after takeover refuse with reason `ended`. None of
these refusals appends an Event.

RESET 1.5a4 adds `lane-not-set`: precondition, exit 4, not retryable,
subject type `issue` naming the task, exact detail key `task_id` (string).
Its repair command is `fid issue update <task> --lane
<routine|standard|assured> --actor <actor>`, as shown with `fid start` above.

RESET 1.5a1 adds the following closed code (Mark, 14 November 2026):

| Code | Category / exit | Subject | Exact detail keys |
|---|---|---|---|
| `policy-unavailable` | precondition / 4 | null | `problem` (string), `ref` (string or null), `commit` (string or null) |

It is not retryable. It replaces the previous `invalid-argument` mapping
for both an origin without a known default branch and invalid or
unreadable `.fid.toml` on origin. `problem` describes what failed;
`ref` and `commit` identify the source read when known (for example
`origin/main` and its commit), otherwise null. Human failures use the
renderer’s `error:` and `reason:` lines, followed by `next:` with
`git remote set-head origin <branch>` for an unknown default, or
`fix .fid.toml on the default branch, merge it and git fetch` for an
invalid or unreadable policy file. Status reports this as a warning
instead of refusing.

If running `git remote` itself fails, the same code has
`problem: "couldn't list git remotes: <error>"`, null `ref` and `commit`,
and the human `next:` action is `check the repository: git remote -v`.
Status shows the same policy-unavailable warning and repair action. This
case was decided by the orchestrator under Mark's delegation for the 1.5a
series on 14 November 2026.

## 6. Removed Commands

The supported baseline contains no:

```text
fid dev ...
fid policy ...
fid role ...
fid escalation ...
fid triage ...
fid tl ...
fid qa-plan ...
fid contract ...
fid template ...
fid board
fid search
fid merge --journal
fid event revert
```

RESET 1.5c2 also removes these entire command families:

```text
activation attention capabilities claim concern decision deviation envelope
evidence finding gate human-outcome ready receipt requirement resume review
run stats stop workflow
```

Except for `review` (the lean record/list/post contract) and
`finding` (close, answer and defer), each is an unknown root command (`operation-unavailable`,
exit 2), with the known path `fid` as its operation and sorted current
alternatives. The retired review schema, Event and tables stay refused;
review delivery is implemented in RESET 2.1c. There is no `fid admin`. Issue Concern
fields, the `issue list --has-concerns` filter, requirement/decision reference
callers and execution-chain transition guards are removed. `fid next`
replaces `ready`. `sql` and `journal inspect` remain diagnostics.

## 7. Purpose-Built Human Rendering

`fid status` and `fid show` (section 3) render text and JSON from one
result, using fid's display package. `status` and `show` alone use the
terminal pager described above.

## 8. Input, grammar and failure rendering

The retained command tree, descriptors and failure registry define the
current CLI surface. Earlier Phase 21–24 kernel leaf contracts are withdrawn
by ADR-0014. The remaining shared input and transport rules follow.

Task 9a delivers grammar failures before workspace resolution on every live
leaf except the Human-only `journal inspect` and `sql` diagnostics. Each is
`invalid-argument`, category `invalid-input`, not retryable, exit 2, subject
null, with the leaf's operation ID and these exact details:

| Class | `argument` | `value` | `allowed` | Message source |
|---|---|---|---|---|
| Unknown flag | full spelling (`--bogus` or `-x`) | full spelling | sorted local long flags, excluding `--help`, `--json`, `--project`; `[]` at root/groups | pflag's detecting text |
| Missing required flag | `--<name>` | null | `["--<name>"]` | `--<name> is required` |
| Malformed flag value, missing value, or value outside its domain | `--<name>` | full supplied value, or null when missing | declared domain | pflag's detecting text; domain checks use `invalid argument "<value>" for "--<name>" flag` |
| Wrong positional count | `arguments` | first extra argument, or null when missing | positional spellings from `Use`, or `["no arguments"]` | Cobra's argument-count text; a no-argument leaf uses `<path without fid> accepts no arguments` |

Flag domains are declared by leaf and spelling in `cmd/fid/flag_domains.go`
and `cmd/fid/flag_domain_text.go`
from the newest Phase 21–24 contract row. Both parsing and the malformed-value
corpus use that declaration. Enumerations are sorted bytewise, integer ranges
use `"<lo>..<hi>"`, patterns retain the contract spelling (including
`<type>:<id>`), and other domains use the pflag type name. Free text missing a
value therefore has `allowed:["string"]`. Range checks retain pflag's integer
syntax; the string `--expected-revision` retains its decimal-only constraint.
The declarations cover limits, priorities, sequence, ordinal, version,
expected revision, all leaf booleans including `--example`, and closed string
values for issue, link and activity filters.
Inherited title, AC text, description, colour and label bounds retain Phase
20’s byte limits and trimming. Opaque IDs and workflow states use the printable Unicode scalar domain;
Actor flags use typed identities.
`issue update --labels=` retains its explicit empty-list clearing behavior.
Cross-flag relationships keep their existing domain checks. Idempotency keys
and standalone Actor flags retain their accepted text constraints. Input
shapes and domain-package validation remain governed by their slice contracts.

An unknown word at the root or any group returns `operation-unavailable`,
category `invalid-input`, not retryable, exit 2, subject null, and
`{argument:"command",value:<full word>,allowed:<sorted help subcommands>}`.
Its envelope operation is the known dotted path (`fid`, `artifact`, `issue.ac`),
and its message source is `unknown command "<word>" for "<CommandPath>"`
without suggestions. The root alternatives include `help` and `completion`;
those generated commands have no operation ID. A group without a word keeps
its Human help on stdout and exit 0. An extra word on a leaf is a positional
failure.

All these message sources pass through the shared `internal/cli` construction:
escape LF/CR/TAB as `\n`, `\r`, `\t` and other control characters as
`\uXXXX`, preserving existing Go `%q` escapes, then count Unicode scalars.
Up to 512 scalars remains intact; longer text retains the first 511 plus `…`,
dropping any escape crossing the cut whole, without padding. Count before JSON
serialization and without the Human label or terminator. This precedence can
shorten the specified source text; details always preserve the full spelling
or value. Human stderr starts with `error: <final JSON message>` plus a newline,
with exit 2 and no second escape or cut. It then prints `reason: <code>` with
key details, followed by `blocked by:` only for nonempty Stop/Attention IDs
and `next:` only for section 9's closed table. The reason line counts its
label and escaped detail toward 200 scalars; a longer line retains 200 plus
`…`, dropping an escape crossing the cut whole. Grammar `invalid-argument`
failures on `--input` leaves include the help/template hint; unknown commands
and failures on leaves without `--input` omit it.

```text
error: unknown flag: --bogus
reason: invalid-argument; argument=--bogus, allowed=--actor|--example|--idempotency-key|--input
next: see fid artifact verify --help; template: fid artifact verify --example
```

Generic fallback messages use the escaped, bounded domain cause.

A root flag that a stable leaf does not accept fails before workspace access
with `invalid-argument`, exit 2 (DEC-286). For `export import --project`, both
`details.argument` and `details.value` are `--project`; `allowed` contains
`--idempotency-key` and `--repository`, the sorted local long flags excluding
help and root flags. Human stderr is exactly:

```text
error: --project is not available for export import
reason: invalid-argument; argument=--project, allowed=--idempotency-key|--repository
```

This leaf has no `--input`, so it prints no `next:`. The same refusal applies
outside a project and for a nonexistent project path. The `sql` and
`journal inspect` diagnostics retain their raw root-flag refusals.

Flag combinations and disagreements between accepted rows and live flag sets
retain their existing behavior (CF-047); task 9a records them in its completion
report. The live tree is the corpus scope. Family request fallbacks remain
available only for their existing non-grammar causes.

`artifact verify` and `review record` use descriptor-generated help, `--example` and closed
structural validation. Help lists field paths, types and required status;
examples use documented field order and two-space indentation. `--example`
works outside a project and refuses combinations in this order: `--json`,
`--project`, `--input`, `--actor`, `--idempotency-key`, then positional
arguments. Unknown or duplicate input fields and missing, null or wrongly
typed `artifact_id` give `invalid-argument` with the exact field path.

*Last updated: 2026-11-22 (RESET 2.1b: accepted work and rework allowances)*
