# CLI CONTRACT PLAN - fid

> **Kernel commands superseded on 2026-11-10 by ADR-0014.** The golden-path
> commands are built on a lean core, not composed from kernel commands, and
> the kernel's commands are removed after task 1.5a rather than moved under
> `fid admin`. Where they conflict, ADR-0014 and `RESET.md` win. Later tasks
> update this spec as they build.

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

Phase 20 convergence qualification: the global spellings are initially valid
only for leaves listed in the approved Phase 20 retained-operation inventory.
An explicitly supplied `--project` or `--json` on a listed legacy leaf rejects
before its handler; neither flag may be accepted and ignored. Later slices
remove or migrate those leaves and retire this temporary restriction.

Phase 21 removes that mixed-tree restriction. Every retained mutation uses
`--actor`; the unrelated stats grouping spelling is `--group-by`; no live
`--by` flag remains. Complex Phase 21 mutations use exactly one
`--input <path|->` JSON document and require `--idempotency-key`; kit-mode
mutations additionally resolve the required `--actor <actor-id>` through the
active Activation rather than trusting caller-supplied capabilities.

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

To distinguish this transport wrapper from the durable **Work Envelope**
entity, the wrapper is named the **Response Envelope** in specifications, code,
and tests.

`error.category` is a closed set: `invalid-input`, `not-found`, `conflict`,
`authority`, `precondition`, `integrity`, `internal`.

## 3. Target Command Tree

### Project and capability

```text
fid init
fid backup
fid restore
fid status [--all]
fid show <id>
fid capabilities
fid rebuild
fid maintenance check
fid journal inspect <event-id>
```

`capabilities --json` is the first adapter operation. Maintenance reports
journal/projection divergence, invalid references/sequences, stale Claims,
abandoned Attempts, registration conflicts, Artifact state, and activation
disagreement.

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
  next step. Kinds include backup, waiting-for-you and worktree-missing (below).
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
  `backup`, `waiting-for-you` and `worktree-missing`; `subject` is null for it, and `{type, id}` for later
  kinds about an issue.
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
  - `Working`: on a held task, `here, since 09:02, until 10:02`, or
    `in <full path> by <actor>, since 08:40, until 09:40` elsewhere. An
    interrupted task reads `⚠ interrupted in <full path> by <actor>: work
    stopped at 07:15 without a park`, using the hold's last heartbeat. Times are local `HH:MM`, adding the
    date when it isn't today. A lapsed hold on an unstarted task frees it,
    so it has no Working row.

  Labels are 10 columns wide after the two-space indent, so values begin
  at column 14. Later tasks add rows after `Lane` (checks, reviews,
  approval, and next as the last), a NEEDS YOU block
  under the heading, and an epic's plan and validations.
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
  owner role or actor (apart from the actor in Working), or links other than
  `blocks` and `blocked_by`.
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
 "scheduling": null,
 "blocked_by": [{"id": "FID-1a2b3c", "kind": "task", "title": "Escape controls in status",
                 "status": "done", "status_category": "completed"}],
 "blocks": [<reference>, <reference>],
 "progress": {"done": 1, "total": 3},
 "children": [{"id": "FID-aa01b2", "kind": "subtask", "title": "Check the manifest schema",
               "status": "done", "status_category": "completed", "sequence": null,
               "priority": null, "blocked_by": []}]}
```

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
fid issue dependency add|remove|list
fid label add|list
fid issue label add|remove
fid link add|remove|list
fid next [--claim --actor <actor>] [--explain]
fid ready
fid workflow define|assign|list
```

`workflow assign` selects assignment scope with exactly one of `--kind <kind>`
or `--project-wide`. The latter spelling avoids collision with global
`--project <path>`, which always selects where project discovery starts
(section 1).

Standalone inventory mutations record Actor attribution without claiming
protected Human authority. Execution-bound transitions additionally require
Envelope revision, expected Issue revision, Actor authority, and active Claim.
An `issue transition` that completes a task ends its hold in the same
`issue.transitioned` Event. Cancelling tasks and epics is a person's
decision, with the contract below.

**`fid next`** (task 1.4a; decisions record, D8's 1.4 details) hands a
worktree its next task. It walks one order and serves the first task
nothing stops:

1. **This worktree's own hold.** A live hold is returned as it is (kind
   `resume`); one worktree never holds two tasks. A lapsed hold on a
   started task means the task was interrupted: it is kept for this
   worktree, and `--claim` takes a new hold that replaces the old one
   (kind `resume`). A lapsed hold on an unstarted task frees the task for
   any worktree, this one included.
2. **Unblocking work:** tasks that another open task waits on, by a
   `blocks` or `blocked_by` link, directly or through its open epic: a task
   that blocks an epic counts only while the epic has open tasks. Until
   task 3.1, all work counts as approved.
3. **The other unstarted tasks**, in `fid status`'s order: an in-progress
   epic's, then a to-do epic's, then those with no open epic; within each,
   sequence, priority (lower first), creation and ID, a missing sequence or
   priority last.

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
| `held-elsewhere` | another worktree's live hold |
| `interrupted-elsewhere` | another worktree's lapsed or parked started work |
| `epic-closed` | the task's epic is completed or cancelled |
| `epic-blocked` | the task's epic has a blocker that isn't completed or cancelled |
| `not-next-in-sequence` | its sequential epic serves another task first; `after` names the task the epic holds or serves, or, when it can serve none, its next task |
| `blocked` | a blocker, from either link, that isn't completed or cancelled; `blocked_by` names them |

Reasons are informational, and later tasks add more: a loop never branches
on them, only on `kind`.

- `fid next` reads (operation `next.get`). It refreshes the projection like
  every read and records no Event.
- `fid next --claim --actor <actor>` takes the served task as a hold for
  this worktree (operation `next.claim`), in one `hold.taken` Event under
  the write lock. The Event ends this worktree's earlier hold and any
  lapsed hold on the task. Re-run in the same worktree, it returns the hold
  it already has, with a null `event`; there is no idempotency key. Two
  worktrees asking at once aren't refused: the lock serialises them, and
  the second passes over the first's task as `held-elsewhere`. The
  database backs the lock up with at most one open hold per task and per
  worktree; if that ever refuses, it is a bug, reported as
  `internal-error`.
- `--explain` lists the whole order with a verdict for each task, with or
  without `--claim`. It changes only the text.

A hold records its ID (`HLD-` and 6 hex digits), the task, the worktree
(the SHA-256 of its root path, and the path), the actor, `since` and
`expires_at`, an hour later (D8), `last_heartbeat_at` (initially `since`), and,
for a resume, the hold it replaces. Whether a hold has lapsed is read from
the clock, never stored. `fid next` reads the clock once, just before it
walks the order: after refreshing the projection and, with `--claim`, under
the write lock, so a hold that lapses while it waits for the lock counts as
lapsed. A hold ends when a new hold replaces it, or its
task is completed or cancelled, parked as interrupted, or taken over. A
parked started task keeps its worktree reservation after the hold ends.
`done` arrives in 1.6. A worktree is identified by its root path, so a moved
worktree is another worktree.

The JSON result, with exit 0 for every kind:

```json
{"kind": "new", "worktree": "/home/mark/projects/fid",
 "hold": {"id": "HLD-4c1e9a", "task_id": "FID-9f13aa", "actor": "agent:claude-builder",
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

- `kind` is `new`, `resume` or `none`; task 2.1 adds `rework`, and `review`
  is decided with 2.1 and the loop. A loop reads `.result.kind`,
  `.result.task.id` and `.result.hold.id`, and on failure the exit code and
  `.error.retryable`.
- `worktree` is this worktree's root path. `task` is the served task in
  `fid status`'s task shape (section 3), or null for `none`. `epic` is the
  served task's epic, with its progress as `fid status` gives it, or null
  when the task has no epic and for `none` (Mark, 12 November 2026).
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
Epics have `todo`, `in_progress`, `done`, `cancelled`; subtasks keep their
existing transitions. Human state labels are “to do”, “in progress”, “in
review”, “done”, “cancelled”. The retired task states and epic `wrapped`
and `closed` are removed, without data conversion. `workflow assign`
refuses task and epic assignments; project-wide assignments affect only
subtasks. Task and epic creation always uses the fixed default workflow.

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
`lane-lowering-requires-waiver`, until `waive` in 2.2b.

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
`task-todo`: review-driven rework belongs to 2.1. Generic transitions refuse
task start (`use-start`), hand-off (`use-hand-off`), and moving a started
task back to todo (`task-already-started`). Epics also refuse manual start.
For example:

```text
error: FID-9f13aa starts with fid start
reason: transition-not-permitted; from_state=todo, to_state=in_progress, unsatisfied=use-start
next: fid start --hold <hold> --actor <actor> --idempotency-key <key>
```

Cancellation stays person-only. Completion refusal belongs to `done` in
1.6; subtasks retain their existing behavior.

**Hold lifecycle (RESET 1.5a2; Mark, 14 November 2026).**

```text
fid heartbeat --hold <id> --actor <actor>
fid park --hold <id> --actor <actor> --interrupted --reason "<text>" --idempotency-key <k>
fid takeover <task> --actor human:<name> --reason "<why>" [--idempotency-key <k>]
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
- Takeover (`task.takeover`) is for a person; agents and services get
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

### Requirement references

```text
fid requirement record-reference|link|get|list
```

This surface stores stable identity, source reference/digest, status, and Issue
links only. Requirement prose is never imported.

Each requirement has two identities. `requirement_id` is the project's stable
ID (for example `REQ-leaf`); `id` identifies one version of the reference
(`RREF-…`), and every supersede creates a new one. Every `RequirementRef`
carries both (Phase 23 DEC-269). `requirement get <id>` accepts either: a
stable ID returns its active reference, or else its newest; a reference ID
returns that row. Human views name the requirement by its stable ID, show the
reference ID separately, and list Issue links in their own labelled section.

### Work Envelope and projections

```text
fid envelope create|get|revise|activate|complete|cancel
fid next
fid resume <issue-id>
fid projection context|changes|completion|human-brief
```

Envelope writes consume an approved versioned JSON value. `resume`, `next`, and
all projection commands are derived reads. `human-brief` may explicitly render
its named Human schema as terminal text, Markdown, or HTML; JSON remains
canonical.

`resume` shows where the Issue stands in its workflow (Phase 24 DEC-290).
Its JSON `issue` object carries `revision` and `next_states`. That is every
state with a workflow edge from the current status, as `issue transition`
checks it (`cancelled` included). Each item is `{state, category}`, sorted
by state, and the list is `[]` when there is none. The Human Issue block
reads:

```text
Issue P24-example
  Status: qa_review
  Revision: 4
  Next states: cancelled (cancelled), qa_block (started), qa_pass (completed)
  Completion: incomplete
  Envelope: ENV-example (revision 1, active)
```

An Issue with no outgoing edge prints `Next states: none.`. These are
workflow edges, not permissions: authority, Claim, Run, Stop and
Obligation checks still apply when the transition runs. `issue transition`
help points `--to` and `--expected-revision` at `fid resume <issue-id>`.

`envelope create --input` omits `repository_scope` and every Source
`repository`. fid observes the selected project's Git worktree and records that
observation as the scope and on `risk_result.source`, `sources[]`, and every
Obligation Source. A project-file Source may omit `digest`; fid reads the file
and records its SHA-256 digest, or checks a supplied digest. External
references require a digest. `envelope complete --input` names the Envelope
and expected revision; fid observes the completion repository itself.

An Obligation's `registration` is nullable: gate-evidence names its immutable
Registration, while Review and Human Outcome Obligations supply null. Prior
Registrations and Envelope scopes must name the same project and worktree,
or an imported destination whose binding names that project and exact root.
A plain copy with no import binding is refused without an Event. Claim Events
retain the Envelope scope; the acquired Claim records its local observed
worktree. A Claim acquired in an imported destination may be renewed,
released or expired there; stored authority uses the same shared worktree
predicate as acquisition (DEC-283).

### Claims, Runs, Stops, and Attention

```text
fid claim acquire|renew|release|list
fid run start|observe|observe-budget|checkpoint|close|abandon|get|list|inspect|recover
fid stop raise|resolve|list
fid attention list|resolve
```

No command assumes a Developer or QA role. Actor identity, capabilities,
Envelope grants, Issue, worktree, mutation scope, optional Obligation, and
lease version are explicit machine fields.

Human `stop list`, `attention list` and `resume` render Stops as labelled
facts. Every StopDetails member has its own line. An open Stop names the
command for its resolution Condition, then the after-the-fact disposition
commands. For example, an operator Stop appears in `stop list` as:

```text
Purpose: Inspect Stops, blocked actions and qualifying resolution conditions.

stop.list
Stop STP-example
  Subject: run:RUN-example
  Code: operator-stop
  Blocked actions: gate.run
  Reason: operator approval
  Reference: DEC-human
  Resolution condition: work-ended for run:RUN-example
  Resolver capability: stop.resolve
  Status: open
  Next: fid run close or fid run abandon; after the fact: fid stop resolve, then fid attention resolve where no lifecycle Event disposed it.
```

`attention list` includes the linked Stop with the same labelled facts:

```text
Purpose: Inspect Attention and the dispositions of its linked Stops.

attention.list
Attention ATT-example
  Status: open
  Code: operator-stop
  Linked Stops: STP-example
  Stop STP-example
  Subject: run:RUN-example
  Code: operator-stop
  Blocked actions: gate.run
  Reason: operator approval
  Reference: DEC-human
  Resolution condition: work-ended for run:RUN-example
  Resolver capability: stop.resolve
  Status: open
  Next: fid run close or fid run abandon; after the fact: fid stop resolve, then fid attention resolve where no lifecycle Event disposed it.
```

After qualified resolution, the Stop ends with these lines and has no
`Next:` line. Resolved Attention also shows its resolution Event:

```text
  Status: resolved
  Fact Event: EVT-fact
  Resolution Event: EVT-stop-resolved

Attention ATT-example
  Status: resolved
  Code: operator-stop
  Linked Stops: STP-example
  Resolution Event: EVT-attention-resolved
```

The Condition commands are `fid run observe` for `executor-terminated`,
`fid run recover` for `run-recovered`, `fid activation recover` for
`activation-recovered`, `fid run observe-budget` for `budget-observed`,
`fid run close` or `fid run abandon` for `work-ended`, `fid review record`
for `review-recorded` and `finding-resolved`, `fid human-outcome record`
for `human-outcome-recorded`, and `fid decision record-reference` then
`fid human-outcome record` for `decision-satisfied`. `capability-available`
instead says to wait for the named capability. These are qualifying-fact
steps; disposition follows only when a lifecycle Event has not disposed
the Stop and Attention already.

Human repository lines use only `HEAD <12 hex>; clean` or
`HEAD <12 hex>; dirty`. `resume` shows `same worktree` for the recorded
root and `different worktree` for a plain copy selected with `--repository`.
`run inspect` shows `same worktree` on success; a plain copy keeps the
`recovery-state-conflict` refusal (exit 4, empty stdout) in both modes,
as specified by DEC-287. Full repository observations remain in JSON.

`claim acquire --input` names the Issue, Envelope revision, path scope, mode
and lifetime. fid observes and records the selected worktree ID, checking it
against the Envelope scope. Renew, release and expire retain their existing
inputs.

`run start`, `run close` and `run abandon` require
`expected_repository:{head_oid,dirty}` instead of `repository`. Use
`git rev-parse HEAD` and `git status --porcelain` for the expectation. fid
observes under the write lock after replay lookup, compares these two
preconditions, and records the full observation. A mismatch returns
`repository-mismatch` with `{expected:{head_oid,dirty},observed:Repository}`
and no Event. Start checks the Claim and Envelope as the same project and
worktree, including a matching import destination binding; edits and commits
since Envelope creation are accepted. A plain copy has no such binding and is
refused before an Event. Close and abandon record `terminal_repository` with
`base_oid` equal to observed HEAD. Close evaluates Obligations against that
observation; stale Receipt-backed Evidence yields a failed Run. Unsafe or
unobservable terminal repositories return `recovery-state-conflict`.

`run observe-budget` accepts nullable `measurement_ref:SourceInput`. With a
Source, fid observes the repository and checks a supplied project-file digest
or computes an omitted digest; the stored Source and Event carry that
observation. A null measurement does not observe. Start and budget observation
return `repository-unavailable` on an unobservable worktree. Equal keyed
retries replay before any repository or Source read, even after HEAD changes,
a referenced file changes or disappears, or the worktree becomes unsafe.

All seven Run mutation leaves document their closed input fields in `--help`
and provide a standalone, exclusive `--example` JSON template. Structural
failures name the exact `input.<path>` and do not append an Event.

### Registered gates and evidence

```text
fid gate register|resolve|run
fid receipt get|list|ingest-attested
fid evidence link|get|evaluate
```

`gate run --input` accepts `registration_id`, `registration_version`,
`obligation_id`, `run_id` and required `expected_repository:{head_oid,dirty}`.
Use `git rev-parse HEAD` and `git status --porcelain` for the expectation.
fid observes once before spawning and checks the expectation before any
Attempt or Event. Registration and Run checks use the same project/worktree
relationship, including a matching import binding. Edits since Run start are
admitted with a matching expectation. Attempt, Receipt and execution Events
retain that observation; process exit does not cause a second observation.
The `fid.live.v1` framing contract is unchanged.

`evidence link --input` supplies the Obligation, Evidence type, exactly one
non-null subject ID (`receipt_id`, `review_id`, `human_outcome_id`) and sorted
`required_artifact_ids`. It derives `target_repository` from the Receipt's
terminal repository, the Review repository or the Human Outcome Event.
Linking does not observe the current worktree. `evidence evaluate --input`
takes `{obligation_id:?ID}` (null selects the active revision) and observes
the current worktree once. Tree and dirty-state changes stale repository-bound
Evidence; reverting restores freshness. Human Outcome Evidence stays fresh.

These inputs reject the removed `repository` and `target_repository` members.
`gate run`, `evidence link`, `evidence evaluate` and `artifact verify` document
their field trees in `--help` and print exclusive editable `--example`
templates without reading or writing a workspace. Structural failures name
`input.<path>` and append no Event.

`gate register --input` names the activated definition/version and Issue scope.
fid observes the selected worktree and derives the Registration project ID,
worktree ID and binding digest from that observation and scope.

### Review and judgment

```text
fid review record|get|list
fid finding address|list
fid concern record|dispose|get|list
fid deviation record|get|list
fid decision record-reference|list
fid human-outcome record|get|list
```

Phase 24 Task 12 adds `fid concern get <concern-id>`,
`fid deviation get <deviation-id>` and `fid human-outcome get <outcome-id>`
(DEC-274). Each takes exactly one ID and only global flags. JSON returns
`{concern:Concern}`, `{deviation:Deviation}` or `{outcome:HumanOutcome}`
inside the success envelope, equal to that ID's list item, with `event:null`
and exit 0. The operation IDs are `concern.get`, `deviation.get` and
`human-outcome.get`. Unknown IDs return `concern-not-found`,
`deviation-not-found` or `outcome-not-found`, respectively, with exit 3,
subject ID equal to the requested ID, and details `{id}`. Missing or extra
arguments and unknown flags return the accepted `invalid-argument` grammar
failure with exit 2. Human output shows full recorded reasoning with line breaks retained and
continuation lines indented. Capability rows are unchanged.

Task 13 delivers A6/A7 (DEC-272; Phase 24 contract sections 10 and 11).
Each non-empty Concerns, Deviations or Human Outcomes section prints one
judgment disclaimer before its items; an empty section prints `None.` and
no disclaimer. Every view keeps its `Purpose:` line. The switch is deferred
to Phase 26.

Human views show Concern summary, reasoning and disposition reason; Finding
summary and every disposition reason; Deviation expected, observed and
reasoning; and Human Outcome reasoning. Every `list` and `resume` judgment
section collapses whitespace to one line per field, retaining 80 Unicode
scalars uncut and cutting longer text to the first 80 plus `…`. Detail reads,
mutation outputs and `review get` show the full text, retaining line breaks
and indenting continuation lines. Controls including ESC, CR, BEL and NUL
are escaped. Reference labels, paths and digests appear; referenced file
contents remain in the project.

For example, a list displays collapsed recorded fields:

```text
Purpose: Inspect a raised concern and its disposition as judgment.

Concerns
  These records express judgment, not deterministic proof.
  Concern CON-example
    Subject: issue:ISS-example
    Severity: major; Status: open; Raised by: actor-builder
    Summary: Validation needs a second sample
    Reasoning: The first sample varied. Repeat before accepting.
```

The matching `concern get CON-example` retains the recorded line break:

```text
Purpose: Inspect a raised concern and its disposition as judgment.

Concerns
  These records express judgment, not deterministic proof.
  Concern CON-example
    Subject: issue:ISS-example
    Severity: major; Status: open; Raised by: actor-builder
    Summary: Validation needs a second sample
    Reasoning: The first sample varied.
      Repeat before accepting.
```

Review and Outcome inputs are versioned structured records. Protected Human
authority is established through active provenance, not a free-form `--by
human` flag. Decision input is a reference/outcome, never an ADR body.

Phase 23 Task 8 adds `concern dispose` for one irreversible, reasoned
`carried-forward`, `rejected`, or `superseded` disposition. `deviation list`
reads the append-only Deviation chain; a changed judgment is recorded with
`deviation record` and an explicit predecessor reference rather than an update
leaf.

Phase 23 Task 9 adds operational Decision References. `decision
record-reference` verifies a project-file Source digest while holding the write
lock and stores the pointer, authority, outcome, and Issue scope without the
decision prose. `decision list` is a read with optional class, outcome, Issue,
and limit filters. An unsatisfied protected Decision raises one Stop and blocks
only its declared affected actions until an approved reference with the exact
required capabilities and an approved Human Outcome both exist.

Phase 23 Task 5 adds optional `inapplicable_findings:[{finding_id,reason}]` to
`review.record`, alongside retained `resolves_finding_ids` (DEC-242). The
reason is mandatory nonblank text; the later qualifying Review records both
terminal actions atomically. Omission and `[]` canonicalize to omission for
the request digest, retaining Task 4 request compatibility; explicit null is
invalid. No independent terminal-disposition leaf is added. The accepted
Phase 23 input contract (`docs/phases/phase-23/public-interface.md`) fixes
bounds, rejection codes and existing Finding disposition readback.

Phase 24 Task 5 makes `review record` take required
`expected_repository:{head_oid,dirty}` in place of a caller-supplied
`repository`. Use the full output of `git rev-parse HEAD` for `head_oid` and
set `dirty` from `git status --porcelain`. The
coverage item and exclusion Sources, and `evidence_refs`, take `SourceInput`
without `repository`; project-file digests may be omitted for fid to compute.
The stored Review and its Sources carry fid's single observed Repository.
For example, the repository portion of a Review input is:

```json
{"expected_repository":{"head_oid":"0123456789abcdef0123456789abcdef01234567","dirty":false}}
```

A mismatch returns `repository-mismatch` with exactly
`{expected:{head_oid,dirty},observed:Repository}` and appends no Event.
`finding address` keeps its input shape and gains descriptor help and
`--example`.

Phase 23 Task 10 adds provenanced Human Outcomes. `human-outcome record`
requires a human Actor type in a finalized Kit Activation profile carrying the
`human-outcome.record` capability, established through that active provenance
alone; a standalone Issue-inventory `--by` style attribution can never reach
this leaf. It resolves the protected Decisions and the `human-outcome-required`
readiness Stop this exact class and `approved` result qualify, atomically with
its own recording Event. `human-outcome list` is a read with optional class,
result, Issue, and limit filters. The stored Evidence and evaluation carry
`producer:"human"`, `trust_source:"activated-human"`, and
`artifact_availability:"not-applicable"`, never a deterministic-proof value.

### Activation and portable recovery

```text
fid activation prepare|component|commit|abort|finalize|status|recover
fid export create|validate|import
```

`export validate` is always non-mutating. Import preserves identities/history
unless an accepted versioned mapping explicitly transforms them.

## 4. Common Mutation Behavior

Every mutation:

1. resolves project and active activation state;
2. validates contract/schema, Actor provenance, authority, idempotency, expected
   revision, references, and domain preconditions;
3. computes one complete atomic Event/effect set;
4. holds the write lock across catch-up, append, fsync, projection transaction,
   and metadata update; and
5. returns the entity revision and Event identity.

For a keyed mutation, the request digest covers the closed input object as
written, plus the operation, project ID and resolved Actor components already
used by its family. Optional members enter the digest only when supplied;
`review.record` treats `inapplicable_findings:[]` as omission (DEC-242). Values
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

### Named checks (RESET 1.5a5)

`fid check <name>... --hold <id> --actor <a> --idempotency-key <k>` executes
trusted `[checks.<name>]` definitions at the worktree root, in name order,
dropping duplicates. Names are required in this slice; no names gives
`invalid-argument`, `argument: "names"`, null `value`,
`allowed: ["defined check names"]`, asking to name the checks. Unknown names
are refused before execution, with the unknown name as `value` and the sorted
defined names in `allowed`. A named check executes even if the lane does not
require it. 1.5a6 adds the required-set default.

The task must be `in_progress`, with a lane and this actor's exact current hold
in this worktree. Missing lane checks refuse as `checks-not-configured`; an
explicit empty set permits named checks. Unreadable policy gives
`policy-unavailable`. A changed effective digest records a policy snapshot.
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

`show` adds `result.checks` (an array, empty when absent), each `{start, process,
result}`. Process and result are nullable. It selects each name's latest result
and every unfinished start. Start contains ID, task, hold, actor, name, resolved
definition/digest, environment-value digest, policy snapshot sequence, time,
HEAD/tree/dirty observation and the owning fid process. Process identity has
machine, PID and start identity. Result contains finish time, outcome (`passed`,
`failed`, `timed-out`, `cancelled`, `spawn-failed`, `signaled`), nullable exit
code/signal, second tree observation, bounded output digest/excerpt/total byte
count/truncation, resolved executable/digest/platform/fid version, kept and
refused files, `counts` and `reasons`. Qualification requires success, a clean
unchanged tree, current ownership at completion and every listed file kept.
Reasons are `tree-dirty`, `tree-changed`, `ownership-lost`, `artifact-not-kept`.

Show's text adds, for example:

```text
  Checks      ✓ lint passed at 10:05
              ⚠ test passed at 10:12, doesn't count: the tree changed
```

A same-key repeat returns the original response without spawning. A keyed
request cut off before completion gives `check-incomplete` without a result.
Before a fresh check, unfinished starts on the task are inspected by machine,
PID and start identity. A stopped process permits execution, with human warning
`⚠ CHK-100001 never finished; its process has stopped`; a live one refuses as
`check-still-running`, suggesting waiting or stopping its PID and trying again.
The old start remains unfinished in either case.

| Code | Category / exit | Exact detail keys |
|---|---|---|
| `checks-not-configured` | precondition / 4 | `lane` (string) |
| `check-unsuccessful` | gate-unsuccessful / 8 | `check_ids`, `outcomes` (string arrays, unsuccessful checks only) |
| `check-not-qualified` | precondition / 4 | `check_ids`, `reasons` (string arrays) |
| `check-still-running` | conflict / 4 | `check_id` (string), `pid` (integer) |
| `check-incomplete` | conflict / 4 | `check_ids` (string array) |

All are non-retryable and have an issue subject naming the task.
`checks-not-configured` suggests `add checks to [lanes.<lane>], then merge
.fid.toml and git fetch`. `check-still-running` suggests `wait for it to finish,
or stop pid <pid>, then check again`.

## 5. Typed Failure Codes

Initial cross-cutting codes include:

```text
unsupported-capability
unsupported-schema
activation-incomplete
activation-conflict
authority-missing
authority-forged
protected-decision-required
source-digest-mismatch
state-conflict
idempotency-conflict
mutation-conflict
claim-expired
executor-liveness-unknown
recovery-state-conflict
obligation-unsatisfied
evidence-stale
artifact-unavailable
gate-failed
execution-abandoned
review-changes-required
review-rework-exhausted
budget-exhausted
retry-exhausted
```

Owning slices may add closed codes through versioned contract changes.

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

Codes accepted with the Phase 20 standalone Issue slice: `invalid-argument`,
`unsupported-kind`, `unsupported-relation`, `project-not-found`,
`project-already-initialized`, `issue-not-found`, `requirement-not-found`,
`label-not-found`, `link-not-found`, `workflow-not-found`, `duplicate-entity`,
`state-conflict`, `lock-unavailable`, `hierarchy-invalid`, `dependency-cycle`,
`transition-not-permitted`, `authority-missing`, `authority-unknown-actor`,
`authority-denied`, `journal-replay-stopped`, `projection-divergence`,
`result-limit-exceeded`, `internal-error`. Phase 20 maps fail-closed V1 role failures to the three
`authority-*` codes. Phase 21 replaces the V1 checks, but the semantic codes are
not automatically removed: the generic authority contract may retain or refine
them through its own accepted contract change.

Phase 23 Task 9 adds `protected-decision-required` (authority, exit 5) with
exact details `envelope_id`, `decision_class`, `required_capabilities`,
`affected_actions`, and (DEC-263) the Stop-blocked denial fields `stop_ids`,
`attention_ids`, `event` (always null), and `replayed` (always `false`); and `source-digest-mismatch` (precondition, exit 4) with
exact details `display_label`, nullable `path`, `expected_digest`, and nullable
`observed_digest`. A missing Source file is reported through the same digest
mismatch code and appends no Event.

Phase 23 Task 10 adds `human-outcome-required` (precondition, exit 4) with
exact details `envelope_id`, `obligation_id`, `outcome_class`, and (DEC-263)
the same four Stop-blocked denial fields, returned by
`envelope.complete` and `issue.transition-terminal-success` for the blocked
actions of the open readiness Stop only; and `outcome-not-found`
(not-found, exit 3) with exact detail `id`, returned by `evidence.link` for an
unknown `human_outcome_id`. Task 10 reuses Task 6's `outcome-not-permitted`
(precondition, exit 4) registration for its real recorded producer.

The complete Phase 23 added failure-code set is closed by
`project-docs/docs/phases/phase-23/public-interface.md` section 5. The table
below records each code's category, exit status, subject, and exact detail
keys. `null` means the failure-envelope subject is null; a subject type names
the entity ID in the producing rejection. All are `retryable:false`.

| Code | Category / exit | Subject | Exact detail keys |
|---|---|---|---|
| `review-not-independent` | precondition / 4 | null | `review_actor_id`, `conflicting_subject`, `conflicting_actor_id` |
| `review-verdict-invalid` | invalid-input / 2 | null | `verdict`, `reason` |
| `review-changes-required` | precondition / 4 | null | `review_id`, `obligation_id`, `blocking_finding_ids` |
| `review-rework-exhausted` | precondition / 4 | envelope | `budget`, `stop_ids`, `attention_ids`, `event`, `replayed` |
| `finding-already-disposed` | conflict / 4 | finding | `finding_id`, `state` |
| `finding-not-found` | not-found / 3 | finding | `id` |
| `review-not-found` | not-found / 3 | review | `id` |
| `concern-not-found` | not-found / 3 | concern | `id` |
| `deviation-not-found` | not-found / 3 | deviation | `id` |
| `outcome-not-found` | not-found / 3 | human-outcome | `id` |
| `concern-already-disposed` | conflict / 4 | concern | `concern_id`, `disposed_at` |
| `claim-mode-invalid` | precondition / 4 | claim | `claim_id`, `mode`, `required_mode` |
| `protected-decision-required` | authority / 5 | envelope | `envelope_id`, `decision_class`, `required_capabilities`, `affected_actions`, `stop_ids`, `attention_ids`, `event`, `replayed` |
| `outcome-not-permitted` | precondition / 4 | null | `obligation_id`, `obligation_type`, `outcome_id`, `result` |
| `evidence-type-not-accepted` | precondition / 4 | null | `obligation_id`, `evidence_type`, `accepted_evidence_types` |
| `source-digest-mismatch` | precondition / 4 | null | `display_label`, `path`, `expected_digest`, `observed_digest` |
| `repository-mismatch` | conflict / 4 | null | `expected:{head_oid,dirty}`, `observed:Repository` |
| `human-outcome-required` | precondition / 4 | envelope | `envelope_id`, `obligation_id`, `outcome_class`, `stop_ids`, `attention_ids`, `event`, `replayed` |

`requirement-not-found` retains its earlier not-found / 3 registration with
subject `requirement` and detail `id`; the Phase 23 requirement leaves now
produce it. `finding-unresolved` belongs to the Evaluation reason set and is
absent from the failure-code set. The nested `budget`, `event`, `expected`,
`observed`, and `conflicting_subject` objects use the exact closed shapes in
the accepted interface contract.

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
fid handoff ...
fid board
fid search
fid merge --journal
fid event revert
```

TUI/dashboard and search remain candidates, not hidden commands. Direct SQL may
exist only in development builds as an explicitly unstable diagnostic and is
not included in the supported command tree.

## 7. Purpose-Built Human Rendering

Only a command that names a Human decision interface may offer a renderer. The
renderer consumes the same typed projection returned as JSON, declares its own
small set of outputs, and carries no round-trip/import promise. Initial allowed
candidate: `projection human-brief` as terminal text, Markdown, or HTML.
`fid status` and `fid show` (section 3) are two: each one's text and JSON
come from one result, and the text uses fid's display package, ported from
whet (tasks 1.3a and 1.3b).

## 8. Accepted Slice Leaf Contracts

Exact request/response schemas are accepted by the Human with their owning
vertical slice before product code changes. Accepted slice contracts:

| Slice | Accepted contract |
|---|---|
| Phase 20 - clean standalone Issue journey | `../docs/phases/phase-20/contracts.md` |
| Phase 21 - one-Issue operational walking skeleton | semantics: `../docs/phases/phase-21/contracts.md`; exact wire surface: `../docs/phases/phase-21/public-interface.md` |

A slice contract fixes its operation names, `result` shapes, reason codes,
closed request/flag shapes, ordering/limits, redaction/unknown-field behavior,
removal inventory, and the exact set of commands that emit a Response Envelope.
Commands scheduled for deletion in a later slice keep their current unenveloped
output and are listed in the owning slice contract with an expiry phase.

The Phase 21 accepted leaves are `capabilities`; Activation
`prepare|component|commit|abort|finalize|status`; gate
`register|resolve|run`; Envelope `create|get|activate|complete`; Claim
`acquire|renew|release|list`; Run `start|finish|get|list`; Receipt `get|list`;
Evidence `link|get|evaluate`; Artifact `verify`; projection `resume`; and
export `create|validate`, plus retained Phase 20 leaves. `gate run --json` uses
the exact `fid.live.v1` NDJSON channel; other machine leaves retain the
`fid-cli.v1` Response Envelope. Phase 21 has no recovery, Envelope revision or
cancel, telemetry, Stop, Attention, attestation ingest, review/judgment, or
export import leaf. The full JSON shapes, limits, ordering, reasons, trust, and
deferrals live in the cited slice contract. Every retained inventory/read leaf,
fail-closed judgment leaf, exact argument/flag, closed `--input` schema,
operation result, collection order, Event/replay shape, Human rendering,
failure details/channel/exit, and live frame payload lives in the cited public
interface. All closed public enum members are lowercase kebab-case; JSON field
names and opaque workflow-state/entity identifiers are not enum members.

### Phase 22 accepted amendment

`../docs/phases/phase-22/public-interface.md` contains the delivered closed
recovery, liveness, Stop/Attention, budget, derived Run closure, activation
recovery, and portable-import leaf schemas, errors, exits, and live frames.
It explicitly replaces caller-controlled `run finish` and the walking-skeleton
recovery-unavailable response, and updates affected reads and capability rows.
DEC-184 retains `--project`, pins the 23 activated actions and bootstrap
exceptions, requires recovery/abandon observation IDs, defines executor reuse
errors and Attention disposition references, and separates activation version
conflict from competing prepare. DEC-185 supplies the canonical repository
digest contract; DEC-187 records external kit port mappings and exclusions.
Human accepted these leaf contracts and refinements under DEC-188. Later
accepted amendments are recorded with the exact phase contract. DEC-219 adds
Human recovery guidance without changing JSON or mutation authority; DEC-220
records technical qualification and the accepted second Human checklist.
CLI view polish and verbosity remain nonblocking CF-040.

### Phase 23 accepted amendment

`../docs/phases/phase-23/public-interface.md` contains the delivered closed
Review, Finding, Concern, Deviation, Decision Reference, Human Outcome, and
requirement-reference leaf schemas, failure details, and Human views (DEC-222).
Later accepted amendments, all delivered and recorded in that contract:

- DEC-263: Stop-blocked judgment denials carry `stop_ids`, `attention_ids`,
  `event` (null), and `replayed` (`false`) in failure details; StopDetails are
  unchanged.
- DEC-265/266: section 5 declares 19 Failure codes (18 added, 1 retained,
  tabled above). All four `authority-*` codes keep the Phase 21 details
  `{operation,actor_id,required_actions,required_capabilities}` with a null
  subject.
- DEC-269: every Human-mode failure from the Stop-error family (`claim
  acquire`, `run start`, `gate run` Stop denials, `stop` and `attention`
  leaves) writes `error: <message>` to stderr, with the JSON `message` text.
  The `review-rework-exhausted` Stop is resolved by an `approved` Human Outcome
  of class `review-rework-extension`, recorded atomically or cited through
  `stop resolve`. Its Human denial adds one guidance line naming that action.
  `RequirementRef` gains `requirement_id` as described above.
- DEC-270: journals written before `requirement_id` existed are unsupported;
  rebuild refuses them with `journal-replay-stopped` (ADR-0011).

DEC-271 records technical qualification and the Human re-check. DEC-272
schedules the remaining Human operating-experience items (Stop-detail
rendering, denial reasons, hint verbosity, recorded reasoning in Human views,
`--input` authoring and repository-input convergence) as Phase 24.

### Phase 24 accepted amendment (delivery in progress)

`../docs/phases/phase-24/public-interface.md` is the accepted target contract
(DEC-277, amended by DEC-278..282); the preceding delivered rules describe
the implementation until the owning Phase 24 tasks land. Sections 5 and 9
govern the message clarifications recorded in DEC-282:

- A plain-copy `run.start` refusal reaches the Claim binding check first.
  Its generic fallback carries `run claim binding is invalid` once task 10
  applies D6=A. Task 7 preserves the Claim-before-Envelope check order and
  the original/copy/import-binding cases.
- Task 9a constructs grammar and unknown-command messages from their
  specified source text, escapes controls, then enforces the inherited
  512-Unicode-scalar bound. Escaped text up to 512 stays intact; longer text
  keeps the first 511 plus `…`, dropping an escape crossing the cut whole.
  Count before JSON serialization, excluding the Human label/terminator.
  This bound takes precedence over exact parser text or the unknown-command
  formula. Structured details retain the full offending spelling/value.
- Task 10 reuses that construction for DEC-281 generic fallbacks and prints
  the resulting JSON message verbatim after `error: `, without a second
  escape or cut. Codes, details and exits retain their accepted values.

The exact leaf, input, failure and view changes remain indexed in that phase
contract. The shared Human failure renderer follows section 9; the task
deliveries below identify the other implemented portions.

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
values for issue/workflow/link/statistics and execution/judgment filters.
Inherited title, AC text, description, colour and label bounds retain Phase
20’s byte limits and trimming. Opaque IDs and workflow states use the newer
Phase 21 printable Unicode scalar domain; activated Actor flags are IDs.
`issue update --labels=` retains its explicit empty-list clearing behavior.
Workflow CSV membership, uniqueness and per-flag bounds are checked here;
cross-flag relationships keep their existing domain checks. Idempotency keys
and standalone Actor flags retain their accepted text constraints. Input
shapes and domain-package validation remain governed by their slice contracts.

An unknown word at the root or any group returns `operation-unavailable`,
category `invalid-input`, not retryable, exit 2, subject null, and
`{argument:"command",value:<full word>,allowed:<sorted help subcommands>}`.
Its envelope operation is the known dotted path (`fid`, `concern`, `issue.ac`),
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
next: see fid stop raise --help; template: fid stop raise --example
```

Generic fallback messages use the escaped, bounded domain cause. The four
plain-copy refusals retain code `invalid-argument`, exit 2 and their existing
details, with messages `prior registration binding does not match envelope`,
`claim acquire binding is invalid`, `run claim binding is invalid`, and
`run gate: repository binding is invalid`. A terminal gate non-success writes
the same failure lines on stderr after its streamed Human stdout.

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

Task 3 adopts descriptor-generated input help and `--example` templates for
`stop raise`, `stop resolve`, and `attention resolve`. Each adopted leaf's help
lists field paths, JSON types, required/optional/nullable status, declared
limits and allowed values, and the template command. `--example` prints a
two-space-indented JSON object with required fields in documented order; it
works outside a project and cannot be combined with another flag or argument.
Structural `--input` failures name the field in `invalid-argument.details.argument`
(for example `input.resolution_condition.subject.id`) before domain validation.
The remaining `--input` leaves enter this convention in their owning tasks.

Task 9 adopts the same descriptor help, `--example`,
and field-path structural validation for `activation prepare`, `commit`,
`abort`, `finalize`, `resource`, `recover`, and `recover-complete`. The input
shapes retain the Phase 21 baseline and Phase 22 version/resource additions.
`activation component` selects a raw typed `content` array using
`component_type`: `authority-profile` contains ActorProfiles, and
`gate-definitions` contains GateDefinitions. Its help lists both variant
trees, including nested environment entries and bounds. Its deterministic
example uses the first variant, `authority-profile`, with one ActorProfile.
Structural errors identify the selected nested path; mismatched content
identifies `input.content[i]`. Environment values and redaction literals
remain sensitive in failure details. Raw variant string bounds count Unicode
scalars; session references and sensitive literals have no additional text
limit. cwd accepts `.`, while Artifact paths exclude dot segments and are
bounded to 512 bytes. All 39 input leaves use descriptors
for help, examples and closed structural validation.

Task 4 adopts `SourceInput` for `decision record-reference`, `requirement
record-reference`, `human-outcome record`, `concern record`, and `deviation
record`. It is a Source without `repository`; a project-file digest may be
omitted, while an external-reference digest is required. fid resolves each
project-file path under the project root, checks a supplied SHA-256 digest or
records the observed one, and attaches its own repository observation to the
stored Source and Event. `requirement link` and `concern dispose` gain
descriptor help and templates as observing leaves without Source fields.

For example, `fid decision record-reference --example` shows the complete
required input shape. Its `source` can be edited to name a file as follows;
the optional digest is shown only here to illustrate a caller-supplied check:

```json
{"source_type":"project-file","display_label":"Accepted decision","path":"project-docs/decision.md","section":null,"external_id":null,"digest":"sha256:<64 lowercase hexadecimal digits>"}
```

Omitting `digest` asks fid to observe the file. A same-key retry replays the
original Source and Event even after HEAD or the file changes. A mismatch
returns `source-digest-mismatch` with nullable `expected_digest` and
`observed_digest`; an unobservable worktree returns `repository-unavailable`
with a typed reason and detail. Both refuse without appending an Event.

*Last updated: 2026-11-08 (RESET 1.3b: `fid show`, wrapping and the pager)*
