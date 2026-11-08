# CLI CONTRACT PLAN - fid

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
  "result": {"issue": {"id": "FID-123", "kind": "task", "title": "Example", "status": "todo", "status_category": "unstarted", "workflow_id": "task-default", "parent_id": null, "priority": null, "sequence": null, "spec_md_path": null, "revision": 1, "labels": [], "created_at": "2026-08-22T00:00:00Z", "updated_at": "2026-08-22T00:00:00Z", "acceptance_criteria": []}},
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

`fid status` shows what needs you, the epics and tasks in progress and to
do, and how current the ledger's backup is (decisions record D4 / F1, 8
November 2026 details for task 1.3). It takes no arguments and one flag,
`--all`. It reads the projection, the backup record and user config, and
writes nothing. Text output is a rounded box 72 columns wide, headed
`fid status` and today's local date:

```text
╭ fid status ────────────────────────────────────────────── 2026-11-08 ╮

  NEEDS YOU
  ✗ No backup for 4 days: last backed up 2026-11-04 19:00
      → fid backup

  IN PROGRESS
  ◆ FID-e3a1c0  Operating experience              ██░░░░░░░░  2/7 done
      ● FID-71a0c2  Export validation                      in progress
      ● FID-3b9e01  Tidy resume output                       qa review
  ● FID-5d2f77  Fix journal replay line                    in progress
  ● FID-d1d1d1  Leftover started task                     dev complete

  TO DO
  ◆ FID-e3a1c0  Operating experience                           3 to do
      ○ FID-9f13aa  Harden login rate limiting
      ○ FID-c2d4e5  Retry a rejected backup pu…  blocked by FID-71a0c2
      ○ FID-0d4e8b  Document backup age
  ◆ FID-8a1b2c  Loop contract                                  6 tasks
      ○ FID-1f2e3d  Contract fields and exit codes
      ○ FID-2f2e3d  Golden tests for the co…  blocked by FID-1f2e3d +1
      ○ FID-3f2e3d  Scenario test for rework
      ○ FID-4f2e3d  Scenario test for interruption
      ○ FID-5f2e3d  Scenario test for two worktrees
      ○ FID-6f2e3d  Version the contract
  ◆ FID-c0ffee  Kit install                                   no tasks
  ○ FID-77aa01  Tidy help text
  +2 more: fid status --all
╰──────────────────────────────────────────────────────────────────────╯
```

- **NEEDS YOU** lists what waits on a person, each item a summary and its
  next step. Only the backup adds an item so far; later tasks add their
  own kinds.
- **Epics and tasks** appear; subtasks never do. An epic whose own state is
  completed or cancelled is hidden. An epic is in progress when its own
  state's category is `started`, or any of its tasks is started or
  completed; otherwise it is to do. Its progress is its completed tasks over
  its tasks that aren't cancelled.
- **IN PROGRESS** shows each in-progress epic, with its bar and count, and
  its started tasks beneath it with their workflow state, underscores read
  as spaces. Started tasks with no listed epic follow, none or a hidden one.
- **TO DO** shows every unstarted task: under its in-progress epic (`N to
  do`), then under its to-do epic (`N tasks`, or `no tasks`), then those
  with no listed epic. An in-progress epic with nothing to do has no line
  here. A task with a `blocks` or `blocked_by` link from an issue that
  isn't completed or cancelled shows `blocked by <ID>`, the lowest such ID,
  and `+N` for the rest. TO DO doesn't judge readiness: task 1.4 brings
  `fid next`'s order and skip reasons, and turns it into READY.
- The text shows at most 10 task rows in TO DO, then
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
  for its column ends in `…`; backup lines and paths are never cut.
- Text fid didn't write, such as titles, paths and the backup record's
  failure kinds, has its controls escaped before it is fitted to a column,
  as failure messages' are (section 8, Phase 24 amendment): LF, CR and TAB
  as `\n`, `\r` and `\t`, other control characters as `\uXXXX`. JSON keeps
  the values as they are.
- Colour is used only when stdout is a terminal and `NO_COLOR` is unset or
  empty. The symbols are always Unicode.

The backup's line is the box's last, except when it needs you: then it
appears only as the NEEDS YOU item.

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
 "needs_you": [{"kind": "backup", "summary": "...", "next_action": "fid backup", "subject": null}],
 "in_progress": {"epics": [<epic>], "tasks": [<task>]},
 "to_do": {"epics": [<epic>], "tasks": [<task>]},
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
- A NEEDS YOU item is `{kind, summary, next_action, subject}`. `kind` is
  `backup` so far; `subject` is null for it, and `{type, id}` for later
  kinds about an issue.
- An epic is
  `{id, title, status, status_category, state, sequence, priority, progress: {done, total}, tasks}`,
  `state` being `in-progress` or `to-do` and `tasks` the section's tasks.
- A task is
  `{id, title, status, status_category, epic_id, sequence, priority, blocked_by}`,
  `blocked_by` listing its open blockers' IDs in order.
- `hint` is `{summary, next_action}` for a project with no epics or tasks,
  and null otherwise.

Outside a repository, at an unsafe ledger location, and in a repository
without a ledger (`project-not-found`, exit 3), `status` refuses as every
project command does (section 1). Journal replay failures keep their codes;
any other failure, such as a journal whose first Event can't be read, is
`internal-error` (exit 7). Task 1.3a adds no codes or reasons.

### Issue tracker

```text
fid issue create|get|list|update|transition
fid issue ac add|list
fid issue dependency add|remove|list
fid label add|list
fid issue label add|remove
fid link add|remove|list
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
`fid status` (section 3) is one: its text and its JSON come from one result,
and the text uses fid's display package, ported from whet (task 1.3a).

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

*Last updated: 2026-11-08 (RESET 1.3a: `fid status`)*
