# RUNTIME LAYOUT AND CONFIGURATION - fid

> **Phase 20 accepted amendment:** the required `config_version = 1` field and
> temporary `pricing.toml` expiry are accepted with the Phase 20 slice pack.

## 1. `.fid/` Layout

The supported baseline uses:

```text
.fid/
├── journal/
│   └── YYYY-MM-DD.jsonl
├── artifacts/
├── config.toml
├── fid.sqlite
├── write.lock
└── .gitignore
```

| Path | Purpose | Git |
|---|---|---|
| `journal/*.jsonl` | Canonical append-only Events | committed |
| `config.toml` | Minimal fid-local project/storage/operator configuration | committed |
| `fid.sqlite` | Derived query projection | ignored |
| `artifacts/` | Retention-limited content-addressed execution material | ignored |
| `write.lock` | Cross-process mutation lock | ignored |
| `.gitignore` | Ignores projection, artifacts, lock, and temporary files | committed |

Temporary append, migration, export, and artifact files use unique names inside
the owning directory, are fsynced before atomic replacement where applicable,
and are never treated as durable success state.

## 2. `config.toml`

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
- Paths are project-relative, normalized, and cannot escape `.fid/` unless an
  explicit later security decision authorizes an external artifact store.
- Executor limits are absolute backend safety ceilings. Gate-specific timeout,
  environment, sandbox, network, write, credential, success, and retention
  policy comes only from the immutable Gate Definition Registration.
- `artifact_root` is valid UTF-8 of `1..4096` bytes, relative, not `.`, and has
  no lexical `..` escape. An existing symlink component must resolve within
  `.fid/` when the path is used.
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

Phase 22 adopts an explicit interpretation of the section 2
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
| fid storage safety ceilings | `.fid/config.toml` |

There is no live `.fid/gates.toml` registry. fid never reads a mutable authored
gate file at execution time. Standalone clients use the same registration API.

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

`fid init` validates the target, writes into a staging directory, initializes
the supported schema, appends the project-initialized Event, verifies rebuild,
and atomically activates `.fid/`. Failure leaves no partially initialized
workspace.

For the final Phase 21 schema, that Event is the exact role-free Event-v2
fixture in the slice contract: service Actor, project sequence 1, null
repository/idempotency, and no domain Effects. Migrations own reference
workflow seeds. Init and rebuild never require or recreate a `roles` table.

Before first dogfood, existing development workspaces may be deleted and
reinitialized. No automatic migration from pre-ADR-0011 schema or journal
events is a product requirement.

Phase 20 convergence may also write `.fid/pricing.toml` solely for the legacy
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

*Last updated: 2026-08-22 (Phase 21 DEC-157 interface/bootstrap amendment)*
