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
├── write.lock
└── .gitignore
```

| Path | Purpose |
|---|---|
| `journal/*.jsonl` | Canonical append-only Events |
| `config.toml` | Minimal fid-local project/storage/operator configuration |
| `fid.sqlite` | Derived query projection |
| `artifacts/` | Retention-limited content-addressed execution material |
| `write.lock` | Cross-process mutation lock |
| `.gitignore` | Names `fid.sqlite`. Kept for now, though git never looks inside its common directory |

Nothing in the ledger is committed. D4 / F1 makes a backup, not git, its copy
off the machine (RESET task 1.2).

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
expansion, existence checks or network access. Missing or empty targets are
unconfigured. `fid init` asks only if stdin and stdout are both terminals,
`--json` is off, and no target is configured. An empty answer (including
whitespace only) skips setup. Non-interactive and JSON runs do not open user
config. Existing configured files are left untouched. New config directories
are private (`0700`); a saved config is atomically replaced with mode `0600`.
Malformed or unreadable config is refused when init needs to read it.
Task 1.3 adds status warnings until a target is set; 1.2a performs no backup,
push or restore.

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
exports, or Human projections. Gate registration may name an approved
credential policy or allowlisted variable names, but never stores secret
values. Execution receives secrets only through the trusted runtime mechanism
authorized by that policy.

The closed Phase 21 environment is assembled from activated Gate Definition
entries whose source is `literal`, approved `inherit`, or `credential` name.
The ambient environment is not inherited wholesale. Credential values exist
only in process memory long enough to execute and redact output.

*Last updated: 2026-11-07 (RESET 1.2a: committed identity and user backup target)*
