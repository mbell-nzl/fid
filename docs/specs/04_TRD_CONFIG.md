# RUNTIME LAYOUT AND CONFIGURATION - fid

> **Phase 20 accepted amendment:** the required `config_version = 1` field and
> temporary `pricing.toml` expiry are accepted with the Phase 20 slice pack.

## 1. Ledger Layout

The ledger lives in the git common directory, at
`$(git rev-parse --git-common-dir)/fid/`, so every worktree of a repository
shares one ledger and git never tracks it (ADR-0013; decisions record D4 /
F1). fid finds it from the project root, the worktree's top level, by reading
the worktree's `.git` entry as git does, without running git: a `.git`
directory is the common directory, and a `.git` file names the worktree's git
directory, whose `commondir` file names the common directory.

Discovery checks that git agrees: a second `git rev-parse` reports the common
directory (git 2.31 or later). fid refuses, before it writes anything, where
the two name different ledgers once symlinks are resolved (usually `GIT_DIR`
or `GIT_COMMON_DIR` naming another repository), and where the common
directory is inside the worktree but not within a `.git` directory, because
git could track the ledger there. Both refusals are `repository-unavailable`
with reason `unsafe-path` (`06_DESIGN_PLAN_CLI.md` section 1).

```text
<git-common-dir>/fid/
├── journal/
│   └── YYYY-MM-DD.jsonl
├── artifacts/
├── config.toml
├── fid.sqlite
├── backups.jsonl
├── write.lock
└── .gitignore
```

| Path | Purpose |
|---|---|
| `journal/*.jsonl` | Canonical append-only Events |
| `config.toml` | Minimal fid-local project/storage/operator configuration |
| `fid.sqlite` | Derived query projection |
| `backups.jsonl` | Record of backup attempts (section 2.4); not part of the journal |
| `import.marker.json` | Record of the import that made the ledger, by `export import` or `fid restore` (section 2.4); not part of the journal |
| `artifacts/` | Retention-limited content-addressed execution material |
| `write.lock` | Cross-process mutation lock |
| `.gitignore` | Names `fid.sqlite`. Kept for now, though git never looks inside its common directory |

Nothing in the ledger is committed. D4 / F1 makes a backup, not git, its copy
off the machine (`fid backup`, section 2.4).

There is no migration from the earlier layout, which kept the ledger in
`.fid/` at the project root. fid doesn't treat a `.fid/` directory left in a
worktree as a ledger: it's an ordinary path, so its files make the worktree
dirty like any other. Moving it into place by hand keeps it: in the main
worktree, `mv .fid .git/fid`.

`fid init` and `export import` stage a new ledger in
`<git-common-dir>/fid-staging-*`, beside the ledger and so on the same
filesystem, and rename it to `<git-common-dir>/fid/`. Temporary append,
migration, export, and artifact files use unique names inside the owning
directory, are fsynced before atomic replacement where applicable, and are
never treated as durable success state.

## 2. Configuration files

### 2.1. Committed project identity: `.fid.toml`

At the repository root, separate from the ledger:

```toml
[project]
id = "project-0123456789abcdef0123456789abcdef"
```

`fid init` generates `project-` plus 32 lowercase hexadecimal digits from
128 cryptographically random bits, never from a path or hostname. This ID
keys the project's backup folder and is also used in the new ledger's
`config.toml` and `project.initialized` event. A clone with an existing valid
`.fid.toml` keeps that ID and the file's exact bytes. Invalid, unreadable or
non-regular identity files are refused before ledger creation, never
overwritten. A project that already has a ledger still refuses re-init;
there is no automatic migration or rewriting of its identity.

Init publishes a complete identity file without overwriting an existing
one, after staging and checking the ledger and before activating it. A
failure after publication can leave the ID file for the next init to reuse.
The user commits `.fid.toml`; fid does not stage or commit it. Restoring a
fresh clone relies on this committed identity. The file holds only the ID
until task 2.2 adds policy from `origin/<default branch>`; the working-tree
identity is not policy authority.

### 2.2. User config: backup target

`$XDG_CONFIG_HOME/fid/config.toml`, default `~/.config/fid/config.toml`:

```toml
[backup]
target = "git@example.com:mark/fid-ledgers.git"
```

The target accepts any git remote URL or local path, stored verbatim without
expansion, existence checks or network access. C0 controls and DEL are escaped
in TOML so the target round-trips without raw forbidden characters.
Missing or empty targets are
unconfigured. `fid init` asks only if stdin and stdout are both terminals,
`--json` is off, and no target is configured. An empty answer (including
whitespace only) skips setup. Non-interactive and JSON runs do not open user
config. Existing configured files are left untouched. New config directories
are private (`0700`); a saved config is atomically replaced with mode `0600`.
Malformed or unreadable config is refused when init needs to read it, using
`invalid-argument` / `invalid-input`, exit 2 (CLI spec section 3). Init prints
the ID and commit instruction before asking, so a target-save failure cannot
hide them; it reports that the ledger exists and target setup failed.
Task 1.3 adds status warnings until a target is set. `fid backup` and
`fid restore` read the target (section 2.4).

These two small files use a strict single-setting TOML subset: the named
table and key, single-line basic (double-quoted, with TOML escapes) or literal
(single-quoted) strings, blank lines and comments. Duplicate or unknown
tables/keys and unquoted values are refused.

### 2.3. Ledger-local `config.toml`

fid-local configuration is intentionally small:

```toml
config_version = 1

[project]
id = "project-opaque-id"
issue_prefix = "FID"

[storage]
artifact_root = "artifacts"

[executor_limits]
max_output_bytes = 10485760
max_excerpt_bytes = 16384

[maintenance]
abandoned_attempt_grace_seconds = 300
```

Rules:

- `config_version` is required and must be the integer `1` for this baseline.
  Unsupported or missing versions reject before any operation uses config.
- Project ID is stable and not inferred from a directory name. It is generated,
  opaque, non-empty ASCII, and at most 128 bytes.
- Issue prefix is presentation only; IDs remain opaque to consumers. Its exact
  grammar is ASCII `^[A-Z][A-Z0-9]{0,15}$`.
- Paths are relative to the ledger directory, normalized, and cannot escape
  it unless an explicit later security decision authorizes an external
  artifact store.
- Executor limits are absolute backend safety ceilings. Gate-specific timeout,
  environment, sandbox, network, write, credential, success, and retention
  policy comes only from the immutable Gate Definition Registration.
- `artifact_root` is valid UTF-8 of `1..4096` bytes, relative, not `.`, and has
  no lexical `..` escape. An existing symlink component must resolve within
  the ledger directory when the path is used.
- Phase 21 accepts Gate values only within `timeout_ms:1..600000`,
  `output_bytes:1..10485760`, `excerpt_bytes:0..min(output_bytes,16384)`, and
  `artifact_bytes:1..1073741824`. Configured executor ceilings are base-10
  integers in `1..9223372036854775807`, and
  `max_excerpt_bytes` cannot exceed `max_output_bytes`; maintenance grace is a
  base-10 integer in `0..9223372036854775807`.
- Phase 21 `inherit` environment entries name only `CI`, `GOCACHE`,
  `GOMODCACHE`, `GOPATH`, `HOME`, `LANG`, `LC_ALL`, `PATH`, `TMPDIR`, or
  `XDG_CACHE_HOME`. Credential entry `NAME` reads only the parent variable
  `FID_CREDENTIAL_NAME`, passes it to the child as `NAME`, and fails before an
  Attempt if absent. Literal, inherited, and credential names are one unique
  set; no ambient variable is otherwise copied.
- Config contains no role matrix, risk policy, QA policy, templates, model
  prompts, gate commands, credentials, or provider routing.
- Unknown keys and unsupported config versions reject rather than being
  silently ignored.

### 2.4. Backups

`fid backup` (CLI spec section 3) writes to three places, none of them the
journal: the backup repository, fid's own clone of it, and the record.

**The backup repository.** Its `main` branch holds one folder per project,
named by the project ID in `.fid.toml`:

```text
<project-id>/
├── backup.json           # backup wrapper
├── manifest.json         # fid export create's manifest
├── journal/
│   ├── 00000001.jsonl     # at most 64 MiB per piece
│   └── …
└── artifacts/            # available Artifacts, named by digest
```

`backup.json` is canonical JSON with a trailing newline: `schema`
(`fid.backup.v2`), `project_id`, `archive_digest`
(the `sha256:` digest of the deterministic ZIP `fid export create` makes),
and the export's `manifest_schema`,
`event_count`, `terminal_event_hash`, `projection_schema` and `entry_count`,
plus `genesis_event_hash`, which names the ledger that wrote the folder
(below). It holds nothing that changes unless the export does. fid validates
the archive as `export validate` does before it commits, so a pushed backup is
one that `fid restore` can import. The folder holds the archive's members
unzipped and byte for byte, including the export's journal pieces of at most
64 MiB. fid owns its ledger's folder and replaces its complete contents on
every backup; it never changes another path in the repository. The older
`fid.backup.v1` ZIP folder is unsupported; no real backups existed when the
format changed in task 1.2d, so there is no migration.

**One ledger per folder.** Clones of one project each make their own ledger
with the project's ID, and so does `fid init` on a new machine, so their
backups name the same folder. fid replaces a folder only when its
`backup.json` names this ledger: its `genesis_event_hash` must be the hash of
this ledger's genesis Event, the `project.initialized` Event `fid init`
appends first. Otherwise, including when the folder has no `backup.json` or
one that names no genesis Event, the backup refuses with
`repository-unavailable` / `different-ledger`, records the attempt and
pushes nothing, because the folder may be another ledger's only copy. fid
checks after reading the target's tip and before exporting. It never forces
the push, so if another backup changes the target after the check, git
rejects the push, which is recorded as `push-rejected`, and the next backup
checks the new tip. To clear the refusal, keep one
ledger for the project: restore this project from the folder with
`fid restore` in place of this ledger, or, if the folder's backup is no longer
needed, remove the folder from the backup repository's `main` by hand, and
the next backup writes it afresh.

Every `fid init` appends its own genesis Event, and two differ in their
time, to the nanosecond, even for one project. The journal never appends
before its latest date, which after a restore is the date of the last
restored Event, so the genesis Event stays its first line. An
import keeps the journal byte for byte, so a restored ledger keeps its
genesis Event and goes on backing up into its folder. Copies of one ledger,
such as a ledger directory copied by hand, share a genesis Event, and this
check doesn't tell them apart.
Each changed backup is one commit on the target's tip, unsigned, by
`fid backup <backup@fid.invalid>`, with the message `fid backup
<project-id>` followed by the event count, terminal event hash and archive
digest. fid sets that identity itself, so a backup never depends on the
user's git config, which CI and cron may not have. When the new tree matches
the tip's, fid makes no commit and pushes nothing: the outcome is
`unchanged`.

**fid's clone.** `$XDG_CACHE_HOME/fid/backup/<key>/`, by default
`~/.cache/fid/backup/<key>/`. `<key>` is the first 32 hex digits of the
SHA-256 of the target as configured, so each target has its own clone and the
path never holds the target, nor a credential in it. The directory holds:

- `repo.git/`: a bare repository. fid builds each commit with plumbing, so
  there's no working tree or index to leave in a bad state.
- `lock`: an exclusive `flock` held for the whole backup, so backups through
  one clone, such as the loop's and a cron job's, run one at a time. A backup
  waits up to 2 minutes for it. The kernel releases it if fid dies.
- `export/`: the export while a backup runs, removed afterwards.

Each backup asks the target for `main` with `git ls-remote`, fetches it into
`refs/remotes/origin/main` unless the clone already holds that tip, builds
the new tree on it and pushes the commit to `refs/heads/main`. The target is
passed to git on the command line and never written into the clone's config.
A failed or rejected push leaves only unreferenced objects, so the next backup
starts again from the target's tip; a push that another commit beat to the
target is rejected the same way, and the next backup builds on that commit.
Deleting the directory is safe: the next backup makes it again.

Plumbing commands don't trigger git's automatic packing. After fetching and
before exporting, fid counts the clone's loose objects. At 256 or more it
runs `git repack -A -d` and `git prune-packed` in the foreground, under the
clone lock. This packs reachable history and removes its loose duplicates;
unreachable objects from failed pushes are kept. User settings such as
`gc.auto=0` cannot disable this packing. A packing failure is recorded as an
internal failure before fid makes or pushes a new backup commit. Thus a
hundred successful changed backups keep their history in packs, leaving only
the recent loose tail. Unreachable objects from failures can remain loose
beyond the threshold; this maintenance doesn't discard them.

fid runs git with the user's environment and git config, so the user's own
credential helpers, SSH keys, `url.<base>.insteadOf` rules and proxies
apply. fid never embeds, stores or asks for credentials: it sets
`GIT_TERMINAL_PROMPT=0`, so a backup run by the loop or cron fails rather than
waiting for a password. It overrides what could stop a backup or run code:
hooks (`core.hooksPath=/dev/null`), commit signing (`--no-gpg-sign`),
attribute filters on export members (content on stdin has no path), init
templates
(`--template=`), detached garbage collection, an inherited commit identity
or date (every `GIT_AUTHOR_*` and `GIT_COMMITTER_*` variable) or
`GIT_TERMINAL_PROMPT`, and inherited variables that name
another repository (`GIT_DIR`, `GIT_INDEX_FILE` and the rest of `git
rev-parse --local-env-vars`, plus `GIT_QUARANTINE_PATH` and
`GIT_NAMESPACE`). A local-path target should be absolute: git resolves a
relative one from fid's clone.

A backup reads the ledger without its write lock, so it never holds up
fid's writes. A write that lands while the export is read can make that
export fail validation; the failure is recorded and the next backup retries.

**The record.** `<ledger>/backups.jsonl` gets one line for every attempt that
found the ledger, successful or not. fid appends and syncs each line while it
holds an exclusive `flock` on the record, and, where the attempt reached
fid's clone, while the backup still holds the clone, so lines are in the
order attempts finished. A write cut short, by a full disk or a killed
process, can leave the record without its final newline; the attempt that
wrote it reports that recording failed. Before appending, fid repairs such a
tail under the record's lock: a whole record missing only its newline gets
the newline, and any other fragment is cut back to the last complete line,
since no reader can use it. Every line before the new one is then a complete
record, though a reader may still find a torn last line between attempts.
Each line is a JSON object with every key present:

| Key | Value |
|---|---|
| `schema` | `fid.backup-record.v1` |
| `started_at`, `finished_at` | UTC RFC 3339 times |
| `outcome` | `pushed`, `unchanged` or `failed` |
| `project_id` | the ID in `.fid.toml`; null when it's missing, invalid or not the ledger's |
| `commit` | the target's `main` once the backup is done; null on failure |
| `archive_digest`, `event_count`, `terminal_event_hash` | the export's; null on failure |
| `failure` | null, or `{kind, message}` |

`kind` is one of `no-target`, `user-config`, `project-identity`,
`lock-unavailable`, `target-unreachable`, `push-rejected`,
`different-ledger`, `export-failed` and `internal`, matching the CLI failures
in the CLI spec. `pushed` and
`unchanged` are both successful backups: the target holds the ledger as
exported during that attempt. Task 1.3's `fid status` reads the record. The
record is not part of the journal: export and rebuild ignore it, and a
restored ledger starts without one.

**Restore.** `fid restore` (CLI spec section 3) reads what a backup wrote,
into a repository that has no ledger. It checks that first, then reads
`.fid.toml` from the working tree and the target from user config, as a
backup does. It fetches into a private (`0700`)
`<git-common-dir>/fid-restore-*` directory, holding a bare repository made
with `--template=`, and deletes it
afterwards, whatever the outcome; it never uses fid's clone and takes no
lock. It asks the target for `main` with `git ls-remote` and fetches only
its tip (`git fetch --depth=1`), so a restore's time doesn't grow with the
backup history. Git runs as it does for a backup: with the user's
credentials and transport, `GIT_TERMINAL_PROMPT=0`, hooks off and inherited
repository variables dropped. From that tip, restore reads
`<project-id>/backup.json` and checks that its `project_id` is this
project's. The export's own `manifest.json` must name the same project. It
then turns a `fid.backup.v2` folder into the archive to import: it streams
the export members directly from git into a `0600` ZIP file in that private
directory, using export's sorted names, deflate, `0644` member modes and
1980-01-01 UTC timestamps. The rebuilt archive's `sha256:` digest must
match `archive_digest`. Missing, malformed or mismatched manifests and
unsupported folder schemas are refused as damaged. Scratch is beside the
ledger rather than under `$TMPDIR`, which may be memory-backed; import
also holds its own archive snapshot beside the ledger. Restore needs disk
space for both archives, the fetched repository and staged ledger. It
imports the archive as `fid export import` does, with the idempotency key
`fid-restore:<commit>`, into `<git-common-dir>/fid`, the ledger every
worktree shares. The event count and genesis Event aren't checked at run
time; V8's test checks them.

The restored journal is the export's numbered pieces (`00000001.jsonl`
onwards, at most 64 MiB each), holding the backed-up Events byte for byte
and in order; later Events go in dated files after them. The ledger gains
`import.marker.json`, which records the import and the pieces it installed.
A restore writes no backup record, so the ledger starts without
`backups.jsonl`. It keeps the genesis Event, so its first backup finds its
own folder and, if nothing changed, is `unchanged` at the restored commit.

## 3. Cross-Product Configuration

Phase 22 adopts an explicit interpretation of the section 2.3
rejection rule (CF-035, DEC-182): any operation that consumes project config
validates the whole closed file before using a value or writing state.
Config-independent help, archive validation, and diagnostics remain available.
Human accepted this under DEC-188 with the Phase 22 pack; it does not authorize
global validation middleware. Task 11 and VAL-047 prove both consuming and
non-consuming behavior. No new config key is proposed.

| Value | Owner and transport |
|---|---|
| Risk lanes, roles, routing, protected Decisions, Work Envelope policy | Project/kit configuration; compiled into approved Envelopes |
| Gate Definitions | Project-owned machine values under kit schema; transferred through Kit Activation into immutable fid registrations |
| Provider projections and tools | Kit package and project overlay |
| Repository roots/worktrees and allowed/protected paths | Approved Envelope and registration binding |
| fid storage safety ceilings | The ledger's `config.toml` |

There is no live `gates.toml` registry in the ledger. fid never reads a
mutable authored gate file at execution time. Standalone clients use the same
registration API.

## 4. Artifact Storage

- The journal and Receipt retain bounded sanitized excerpts, hashes, byte
  counts, truncation/redaction flags, and opaque Artifact identity.
- Full logs are written only when the active registration names the path and
  permits retention. Phase 21 accepts only regular, symlink-resolved paths
  contained by the exact execution cwd and configured Artifact root.
- Artifact content and output are redacted before framing, hashing, retention,
  error detail, Human rendering, or export and are stored by SHA-256 content
  identity through atomic publication.
- Phase 21 availability is `available` or `missing`; verified loss changes
  availability but never deletes its Evidence Record or digest. Expiry policy
  is deferred with retention automation.
- Pruning is an explicit maintenance operation and cannot remove material still
  required by a `full-material-required` active Obligation.

## 5. Initialization And Reset

`fid init` validates the target, writes into a staging directory beside the
ledger (`<git-common-dir>/fid-staging-*`), initializes the supported schema,
appends the project-initialized Event, verifies rebuild, and atomically
renames the staging directory's ledger to `<git-common-dir>/fid/`. Failure
leaves no partially initialized ledger. Init from any worktree of a
repository creates the one ledger they share; once it exists, init from any
of them refuses with `project-already-initialized`.

For the final Phase 21 schema, that Event is the exact role-free Event-v2
fixture in the slice contract: service Actor, project sequence 1, null
repository/idempotency, and no domain Effects. Migrations own reference
workflow seeds. Init and rebuild never require or recreate a `roles` table.

Before first dogfood, existing development workspaces may be deleted and
reinitialized. No automatic migration from pre-ADR-0011 schema or journal
events is a product requirement.

Phase 20 convergence may also write the ledger's `pricing.toml` solely for the legacy
stats commands owned for deletion in Phase 26 (renumbered by DEC-272). It is not part of `config.toml`,
is not accepted by the strict config parser, and is absent from the first
supported baseline.

## 6. Secrets

No secret belongs in committed fid configuration, journal Events, Receipts,
exports, or Human projections. A backup target is used as configured; fid
stores no credential for it, and the user's own git credential setup supplies
any it needs (section 2.4). Gate registration may name an approved
credential policy or allowlisted variable names, but never stores secret
values. Execution receives secrets only through the trusted runtime mechanism
authorized by that policy.

The closed Phase 21 environment is assembled from activated Gate Definition
entries whose source is `literal`, approved `inherit`, or `credential` name.
The ambient environment is not inherited wholesale. Credential values exist
only in process memory long enough to execute and redact output.

*Last updated: 2026-11-08 (RESET 1.2c: restore)*
