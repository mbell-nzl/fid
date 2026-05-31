# DESIGN PLAN CLI — fid

> CLI flag specifications and output examples for all commands. The
> parent Design Plan (`05_DESIGN_PLAN.md`) has experience principles,
> output standards, and the command tree overview. This document has
> the per-command detail that developers implement against.
>
> Source: `scaffold/initial-spec.md` §9, §10, plus output examples
> designed to cover every v0.1-scoped command.

---

## Conventions

Some flags apply across commands rather than to a single one:

- **`--format json|toon|markdown|yaml`** — output format on every read command
  (default `json`); shown per-command below.
- **`--dry-run`** (v0.2) — accepted by every *write* command. Runs validation
  and builds the would-be journal event, then stops before appending; no state
  changes. Exit code mirrors the simulated outcome (`0` if the write would
  apply, non-zero if it would be rejected), so dry-runs compose with `&&`.
  Malformed flags or an unknown role still fail as ordinary usage errors.
  See `01_PRD.md` PRD-OPS-015.

A would-reject preview:
```text
$ fid dev finish FID-a1b2c3 --by developer --tokens 14823 --commits abc123 --dry-run
dry-run: would REJECT  (command: dev finish)
  transition in_progress → dev_complete blocked by rule all_acs_verified:
    2 of 5 ACs not verified by developer: AC-e5f6a7, AC-f6a7b8
  no journal event appended
```

A would-apply preview (`--format json`):
```json
{
  "dry_run": true,
  "would_apply": true,
  "command": "issue create",
  "validation": {"ok": true, "errors": []},
  "event_preview": {
    "event_type": "issue.created",
    "effects": [
      {"op": "insert", "table": "issues",
       "values": {"id": "<minted-on-append>", "kind": "task", "title": "Add OIDC",
                  "epic_id": "FID-e1f2g3", "status": "todo"}}
    ]
  }
}
```

A would-reject `--format json` mirrors this with `"would_apply": false` and a
populated `validation.errors[]`, each error carrying `rule`, `detail`, and
`entities` — the same vocabulary the real command emits on rejection. The
`event_id` and any generated hash IDs are omitted or shown as
`<minted-on-append>` because they are assigned only at the append step.

Idempotency keys for retry-safe programmatic writes are deferred
(`01_PRD.md` PRD-OPS-016), conditional on a concrete programmatic/MCP need.

---

## 1. Initialisation

### `fid init`

```
fid init [--prefix <PREFIX>]
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--prefix` | string | project directory name (uppercased) | Issue ID prefix |

**Output (stderr confirmation):**
```text
initialized fid project in .fid/
  prefix: FID
  config: .fid/config.toml
  pricing: .fid/pricing.toml
  roles seeded: developer, qa_reviewer, qa_lead, tl_lite, tl, workflow_runner, human
```

**Error (already initialized):**
```text
error: .fid/ already exists
  to reinitialize, remove .fid/ and run fid init again
```

---

## 2. Issue Commands

### `fid issue create`

```
fid issue create --kind <kind> --title <title> [flags] --by <role>
```

| Flag | Type | Required | Default | Description |
|---|---|---|---|---|
| `--kind` | string | yes | — | `epic`, `task`, or `subtask` |
| `--title` | string | yes | — | Issue title |
| `--by` | string | yes | — | Acting role |
| `--epic` | string | no | — | Parent epic ID (for tasks) |
| `--parent` | string | no | — | Parent task ID (for subtasks) |
| `--priority` | int | no | — | 0..5 (0 = critical) |
| `--sequence` | int | no | — | Positive sibling order |
| `--spec-md` | string | no | — | Path to spec markdown |
| `--est-tokens` | int | no | — | Estimated tokens (task/subtask only) |
| `--due-date` | date | no | — | Due date (YYYY-MM-DD) |
| `--label` | string[] | no | — | Labels (repeatable) |
| `--owner-role` | string | no | kind default | Owner role |
| `--owner-actor` | string | no | — | `agent` or `human:<name>` |
| `--scheduling` | string | no | — | Epic only: `sequential` or `prioritised` |
| `--inbox` | bool | no | false | Add to inbox |

**Output (stderr confirmation):**
```text
created issue FID-a1b2c3 (task) "Implement journal append with flush"
```

### `fid issue get`

```
fid issue get <id> [--section <name>] [--format json|toon|markdown|yaml]
    [--context-pack <role>]
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--section` | string | — | `acceptance`, `concerns`, `completions`, `links`, `claims` |
| `--format` | string | `json` | Output format |
| `--context-pack` | string | — | v0.2+: role-tuned bundle (`developer`, `qa_reviewer`, `qa_lead`, `tl_lite`, `tl`) |

**JSON output:** See `05_DESIGN_PLAN.md` §6.1.

### `fid issue list`

```
fid issue list [flags] [--format json|toon|markdown|yaml]
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--kind` | string | — | Filter by kind |
| `--status` | string | — | Filter by exact status |
| `--status-category` | string | — | `unstarted`, `started`, `completed`, `cancelled` |
| `--epic` | string | — | Filter by parent epic ID |
| `--label` | string | — | Filter by label |
| `--priority` | string | — | Filter by priority (e.g. `0,1`) |
| `--inbox` | bool | — | Show only inbox items |
| `--has-concerns` | bool | — | Show only issues with open concerns |
| `--order` | string | — | `sequence`, `priority`, `created`, `updated` |
| `--format` | string | `json` | Output format |

**JSON output:** See `05_DESIGN_PLAN.md` §6.2.

**TOON output example (`fid issue list --kind task --format toon`):**
```
kind	id	title	status	status_category	parent_id	priority	concerns_raised	open_concerns_count
task	FID-a1b2c3	Implement journal append with flush	in_progress	started	FID-e1f2g3	1	false	0
task	FID-d4e5f6	Schema migration runner	todo	unstarted	FID-e1f2g3	2	false	0
```

### `fid issue update`

```
fid issue update <id> [flags] --by <role>
```

| Flag | Type | Description |
|---|---|---|
| `--title` | string | New title |
| `--priority` | int | New priority |
| `--labels` | string[] | Replace the full label set (repeatable) |
| `--owner-role` | string | New owner role |
| `--owner-actor` | string | New owner actor |
| `--spec-md` | string | New spec markdown path |
| `--est-tokens` | int | New estimate |
| `--due-date` | date | New due date |
| `--by` | string | Acting role |

**Output (stderr confirmation):**
```text
updated FID-a1b2c3: priority 2 → 1, title "Implement journal append with flush"
```

### `fid issue label add`

```
fid issue label add <issue-id> <label> --by <role>
```

Adds one label to an issue without replacing the rest of the issue's labels.
The label must already exist; create it with `fid label add` first.

**Output (stderr confirmation):**
```text
added label "journal" to FID-a1b2c3
```

### `fid issue label remove`

```
fid issue label remove <issue-id> <label> --by <role>
```

Removes one label from an issue without changing the rest of the issue's
labels.

**Output (stderr confirmation):**
```text
removed label "journal" from FID-a1b2c3
```

### `fid issue transition`

```
fid issue transition <id> --to <state> --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--to` | string | yes | Target state |
| `--by` | string | yes | Acting role |

**Output (stderr confirmation):**
```text
transitioned FID-a1b2c3: todo → in_progress
```

### `fid issue complete`

```
fid issue complete <id> --by <role> [flags]
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--by` | string | yes | Acting role |
| `--notes-md` | string | no | Path to completion notes |
| `--tokens` | int | no | Actual tokens used |
| `--commits` | string | no | Comma-separated commit hashes |
| `--files` | string[] | no | Files touched (repeatable) |
| `--human-intervention` | bool | no | Mark that a human intervened during the attempt |
| `--intervention-note` | string | no | Short note stored when `--human-intervention` is set |

**Output (stderr confirmation):**
```text
completed FID-a1b2c3 (attempt 1)
  completion: TC-7c3e9f
  tokens: 14823
  commits: abc123, def456
```

### `fid issue claim`

```
fid issue claim <id> --by <role> [--ttl <duration>] [--force] [--renew]
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--by` | string (required) | — | Claiming role |
| `--ttl` | duration | `30m` | Lease duration |
| `--force` | bool | false | Supersede a live claim when the role is permitted to force-acquire |
| `--renew` | bool | false | Renew the caller's current live claim |
| `--note` | string | — | Optional claim context |

**Output (stderr confirmation):**
```text
claimed FID-a1b2c3 by developer (expires 2026-05-27T15:00:00Z)
```

**Renew output:**
```text
renewed claim on FID-a1b2c3 by developer (expires 2026-05-27T15:30:00Z)
```

**Force-acquire output:**
```text
force-acquired claim on FID-a1b2c3 by tl (previously held by developer until 2026-05-27T15:00:00Z)
```

### `fid issue unclaim`

```
fid issue unclaim <id> --by <role>
```

| Flag | Type | Description |
|---|---|---|
| `--by` | string | Acting role releasing its claim |

**Output (stderr confirmation):**
```text
released claim on FID-a1b2c3 by developer
```

---

## 3. Acceptance Criteria Commands

### `fid issue ac add`

```
fid issue ac add <issue-id> --text <text> --ordinal <n> --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--text` | string | yes | AC text |
| `--ordinal` | int | yes | Display order |
| `--by` | string | yes | Acting role |

**Output (stderr confirmation):**
```text
added AC AC-d4e5f6 on FID-a1b2c3 (ordinal 1): "Journal file is flushed to disk before returning success"
```

### `fid issue ac verify`

```
fid issue ac verify <issue-id> <ac-id> --by <role> --status <status> [flags]
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--by` | string | yes | Verifying role |
| `--status` | string | yes | `pass`, `fail`, `blocked`, `not_applicable` |
| `--evidence-type` | string | no | `test`, `source`, `screenshot`, `manual`, `command`, `commit` |
| `--evidence-ref` | string | no | Path / commit hash / command ID / URL |
| `--notes` | string | no | Short free-form notes |

**Output (stderr confirmation):**
```text
verified AC-d4e5f6 on FID-a1b2c3 by developer: pass
  evidence: test → test/journal_test.go::TestAppendFlush
```

### `fid issue ac list`

```
fid issue ac list <issue-id> [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
[
  {
    "id": "AC-d4e5f6",
    "ordinal": 1,
    "text": "Journal file is flushed to disk before returning success",
    "latest_verifications": {
      "developer": {"status": "pass", "verified_at": "2026-05-27T10:15:00Z", "evidence_type": "test", "evidence_ref": "test/journal_test.go::TestAppendFlush"},
      "qa_lead": null
    }
  },
  {
    "id": "AC-e5f6g7",
    "ordinal": 2,
    "text": "Partial write leaves journal intact",
    "latest_verifications": {
      "developer": {"status": "pass", "verified_at": "2026-05-27T10:20:00Z", "evidence_type": "test", "evidence_ref": "test/journal_test.go::TestPartialWrite"},
      "qa_lead": null
    }
  }
]
```

---

## 4. Ready Queue

### `fid ready`

```
fid ready [--epic <id>] [--format json|toon|markdown|yaml]
```

| Flag | Type | Description |
|---|---|---|
| `--epic` | string | Scope to a single epic |
| `--format` | string | Output format (default: `json`) |

**JSON output:** See `05_DESIGN_PLAN.md` §6.3.

**Empty state (stderr):**
```text
note: no ready tasks — all 3 tasks in FID-e1f2g3 are blocked or completed
```

---

## 5. Developer Composite Commands

### `fid dev start`

```
fid dev start <issue-id> --by developer
    [--branch <name>] [--base-commit <sha>] [--worktree <path>]
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--by` | string | yes | Acting role |
| `--branch` | string | no | Override the auto-captured git branch |
| `--base-commit` | string | no | Override the auto-captured base commit |
| `--worktree` | string | no | Override the auto-captured worktree path |

Composes: claim acquire + transition -> in_progress + open agent_run. Git
context (branch, base commit, dirty flag, worktree) is auto-captured read-only
from git (v0.2); any flag overrides capture. Outside a git worktree the run is
opened with no git context and a warning. The run records `git_source`
(`captured` | `provided` | `none`).

**Output (stderr confirmation):**
```text
dev start FID-a1b2c3
  claimed by developer (expires 2026-05-27T15:00:00Z)
  transitioned: todo → in_progress
  agent run: RUN-b2c3d4 opened (branch phase-1/core @ 9f2a1c0, clean, git_source=captured)
```

### `fid dev gate`

```
fid dev gate <issue-id> --cmd <command> --by developer
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--cmd` | string | yes | Gate command to record |
| `--by` | string | yes | Acting role |

Records an append-only `gate_runs` row for the active attempt. v0.1 derives
the active attempt number from prior completions for the issue; finish/block
commands snapshot matching gate runs into `task_completions.gate_results`.

**Output (stderr confirmation):**
```text
gate recorded on FID-a1b2c3 (attempt 1)
  cmd: go test ./...
  exit_code: 0
  gates_passed: 3, gates_failed: 0
```

**Gate failure:**
```text
gate recorded on FID-a1b2c3 (attempt 1)
  cmd: golangci-lint run ./...
  exit_code: 1
  gates_passed: 2, gates_failed: 1
```

### `fid dev finish`

```
fid dev finish <issue-id> --by developer --completion-md <path>
    [--tokens <n>] [--commits <hashes>] [--files <paths>...]
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--by` | string | yes | Acting role |
| `--completion-md` | string | yes | Path to completion report |
| `--tokens` | int | no | Actual tokens used |
| `--commits` | string | no | Comma-separated commit hashes |
| `--files` | string[] | no | Files touched (repeatable) |
| `--human-intervention` | bool | no | Mark that a human intervened during the attempt |
| `--intervention-note` | string | no | Short note stored when `--human-intervention` is set |

Composes: close agent_run + create task_completion + transition ->
dev_complete (validates all_acs_verified, has_completion_report). The run's
`head_commit` and `dirty` flag are refreshed from git (auto-captured, v0.2) at
finish.

**Output (stderr confirmation):**
```text
dev finish FID-a1b2c3
  agent run: RUN-b2c3d4 closed (14823 tokens)
  completion: TC-7c3e9f (attempt 1)
  transitioned: in_progress → dev_complete
  claim released
```

### `fid dev block`

```
fid dev block <issue-id> --by developer --reason-md <path>
```

Composes: transition -> failed_gate + completion with reason.

**Output (stderr confirmation):**
```text
dev block FID-a1b2c3
  transitioned: in_progress → failed_gate
  completion: TC-8d4e5f (attempt 1, status: failed_gate)
  claim released
```

### `fid dev resume`

```
fid dev resume <task-id> --by developer [--force]
```

Composes: re-acquire the claim + open a new agent_run linked to the task's prior
run via `resumed_from_run_id`, carrying forward `branch` and `base_commit`. Any
prior still-open run for the task is marked `abandoned` (reason: superseded by
resume). Respects claim rules; `--force` permits a sanctioned steal of a live
claim. fid returns resume *hints* only — it never performs git operations.

**Output (stderr confirmation):**
```text
dev resume FID-a1b2c3
  prior run RUN-b2c3d4 marked abandoned (superseded by resume)
  agent run: RUN-e5f6g7 opened (resumed_from RUN-b2c3d4)
  branch phase-1/core, base 9f2a1c0, last head abc123f
  worktree /home/mark/wt/fid-a1b2c3
  claim re-acquired by developer (expires 2026-05-31T15:00:00Z)
  → check out branch phase-1/core and continue; spec at <spec_md_path>
```

### `fid dev abandon`

```
fid dev abandon <run-id> --reason-md <path> --by <role>
```

Marks the run `abandoned` with the given reason and releases its claim. The run
row is preserved (not deleted). The task is left in its current state; an
abandoned run with the task still in a started state is surfaced by
`fid maintenance check` (started-without-claim) and is resumable.

**Output (stderr confirmation):**
```text
dev abandon RUN-b2c3d4
  status: open → abandoned
  reason: project-docs/docs/phases/phase-1/runs/abandon-RUN-b2c3d4.md
  claim released
```

---

## 6. Concern Commands

### `fid concern raise`

```
fid concern raise --on <issue-id> --by <role> --in-completion <completion-id>
    --summary <text> [--body-md <path>]
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--on` | string | yes | Issue ID |
| `--by` | string | yes | Raising role |
| `--in-completion` | string | yes | Completion ID this concern was raised in |
| `--summary` | string | yes | One-line label |
| `--body-md` | string | no | Path to detailed concern body |

**Output (stderr confirmation):**
```text
raised concern CON-c3d4e5 on FID-a1b2c3: "Possible race condition in concurrent journal writes"
```

### `fid concern address`

```
fid concern address <concern-id> --in-completion <completion-id> --by <role>
```

**Output (stderr confirmation):**
```text
addressed concern CON-c3d4e5 in completion TC-9e5f6g
```

### `fid concern list`

```
fid concern list [--issue <id>] [--status <status>] [--format json|toon|markdown|yaml]
```

| Flag | Type | Description |
|---|---|---|
| `--issue` | string | Filter by issue ID |
| `--status` | string | `open`, `addressed`, `escalated`, `wontfix` |
| `--format` | string | Output format (default: `json`) |

**JSON output:**
```json
[
  {
    "id": "CON-c3d4e5",
    "issue_id": "FID-a1b2c3",
    "summary": "Possible race condition in concurrent journal writes",
    "status": "open",
    "raised_by_role": "developer",
    "raised_in_completion_id": "TC-7c3e9f",
    "body_md_path": "project-docs/docs/phases/phase-1/concerns/CON-c3d4e5.md",
    "created_at": "2026-05-27T10:30:00Z"
  }
]
```

---

## 7. Link Commands

### `fid link add`

```
fid link add <from-id> <relation> <to-id> --by <role>
```

Auto-creates the inverse link. Relations: `blocks`, `relates_to`,
`duplicates`, `supersedes`, `implements`, `caused_by`, `discovered_from`.

**Output (stderr confirmation):**
```text
linked FID-a1b2c3 blocks FID-7c3e9f
  inverse: FID-7c3e9f blocked_by FID-a1b2c3
```

### `fid link remove`

```
fid link remove <from-id> <relation> <to-id> --by <role>
```

Removes the declared link and its inverse row in one journal event.

**Output (stderr confirmation):**
```text
removed link FID-a1b2c3 blocks FID-7c3e9f
  inverse removed: FID-7c3e9f blocked_by FID-a1b2c3
```

### `fid link list`

```
fid link list <id> [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
[
  {
    "direction": "outgoing",
    "relation": "blocks",
    "entity_type": "issue",
    "entity_id": "FID-7c3e9f",
    "entity_title": "Schema migration runner"
  },
  {
    "direction": "outgoing",
    "relation": "implements",
    "entity_type": "requirement",
    "entity_id": "PRD-ISSUE-001",
    "entity_title": "Unified polymorphic issue table"
  },
  {
    "direction": "incoming",
    "relation": "discovered_from",
    "entity_type": "issue",
    "entity_id": "FID-h8i9j0",
    "entity_title": "Add index on concerns(issue_id, status)"
  }
]
```

---

## 8. Activity and Change Commands

### `fid activity list`

```
fid activity list [--on <type>/<id>] [--since <date>] [--event <event>]
    [--format json|toon|markdown|yaml]
```

| Flag | Type | Description |
|---|---|---|
| `--on` | string | Filter to one entity reference, e.g. `issue/FID-a1b2c3` |
| `--since` | date/datetime | Include entries on or after this timestamp |
| `--event` | string | Filter by projected activity event name |
| `--format` | string | Output format (default: `json`) |

**JSON output:**
```json
[
  {
    "id": "ACT-a1b2c3",
    "journal_event_id": "EVT-20260527-000142",
    "ref_type": "issue",
    "ref_id": "FID-a1b2c3",
    "event": "status_changed",
    "actor": "developer",
    "occurred_at": "2026-05-27T10:15:00Z",
    "details": {
      "from": "todo",
      "to": "in_progress"
    }
  }
]
```

### `fid changes`

```
fid changes --since <anchor> [--epic <id>] [--entity <id>]...
    [--format json|toon|markdown|yaml]
```

| Flag | Type | Description |
|---|---|---|
| `--since` | string | Anchor: a journal event ID (`EVT-...`) or a timestamp (`2026-05-30` or RFC3339). Required. The resolver is extensible; review-id and commit anchors ship in later versions. |
| `--epic` | string | Scope the digest to one epic's issues. |
| `--entity` | string | Scope to one or more entity IDs (repeatable). |
| `--format` | string | Output format (default: `json`) |

Returns, grouped by entity, the net field changes and newly-added child facts
for every entity changed since the anchor. An entity appears if it — or an
append-only child it directly owns (ACs/verifications, completions, concerns,
its own links) — had at least one journal effect after the anchor. There is no
implicit parent rollup: an epic does not appear merely because a child task
changed. Net field changes collapse intermediate transitions
(`todo → dev_complete`); the full chronological trail is `fid activity list`.
The digest is computed from `activity_log` and adds no schema. Role-tuned
shaping is a deferred fast-follow.

**JSON output:**
```json
{
  "anchor": {"type": "event", "id": "EVT-20260530-000142", "at": "2026-05-30T09:14:02Z"},
  "scope": {"epic": "FID-e1f2g3", "entities": null},
  "entities": [
    {
      "id": "FID-a1b2c3",
      "kind": "task",
      "title": "Implement ID generation",
      "last_changed_at": "2026-05-30T09:58:11Z",
      "field_changes": [
        {"field": "status", "from": "todo", "to": "dev_complete"},
        {"field": "labels", "from": ["core"], "to": ["core", "idgen"]}
      ],
      "added": {
        "ac_verifications": [
          {"ac": "AC-d4e5f6", "status": "pass", "by": "developer", "at": "2026-05-30T09:50:00Z"}
        ],
        "concerns": [
          {"id": "CON-c3d4e5", "status": "open", "by": "developer", "summary": "hash length scaling untested at 8 chars"}
        ],
        "completions": [
          {"id": "TC-7c3e9f", "attempt": 1, "tokens": 14823}
        ],
        "links": [
          {"relation": "implements", "to": "PRD-ISSUE-001"}
        ]
      },
      "removed": {}
    }
  ],
  "summary": {"entities_changed": 2, "new_open_concerns": 1}
}
```

New fact types (findings in v0.3, comments and governance facts in v0.2b)
appear as additional keys under `added` / `removed` without changing the
contract for existing keys.

**Empty result:**
```text
no changes since EVT-20260530-000142
```

---

## 9. Label Commands

### `fid label add`

```
fid label add <name> [--colour <hex>] [--description <text>] --by <role>
```

**Output (stderr confirmation):**
```text
created label "core"
```

### `fid label list`

```
fid label list [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
[
  {"name": "core", "colour": "#4A90D9", "description": "Core infrastructure"},
  {"name": "journal", "colour": "#7B68EE", "description": "Journal subsystem"},
  {"name": "tech-debt", "colour": "#FF6347", "description": null}
]
```

---

## 10. Claim Commands

### `fid claim list`

```
fid claim list [--stale] [--by <role>] [--format json|toon|markdown|yaml]
```

| Flag | Type | Description |
|---|---|---|
| `--stale` | bool | Show only expired claims |
| `--by` | string | Filter by claiming role |
| `--format` | string | Output format (default: `json`) |

**JSON output:**
```json
[
  {
    "ref_type": "issue",
    "ref_id": "FID-a1b2c3",
    "claimed_by": "developer",
    "claimed_at": "2026-05-27T14:30:00Z",
    "ttl_seconds": 1800,
    "expires_at": "2026-05-27T15:00:00Z",
    "is_stale": false,
    "note": "Implementing journal append"
  }
]
```

---

## 11. Policy Commands

### `fid policy show`

```
fid policy show [--format json|toon|markdown|yaml]
```

Displays the active permissions matrix. v0.1 reads the in-code default
matrix; v0.2 reads the project-configurable `role_permissions` table.

**JSON output:**
```json
[
  {
    "action": "issue.transition.dev_complete",
    "developer": true,
    "qa_reviewer": false,
    "qa_lead": false,
    "tl_lite": false,
    "tl": false,
    "workflow_runner": false,
    "human": true
  }
]
```

---

## 12. Run Telemetry

### `fid run record`

```
fid run record --issue <id> --role <role> --model <model>
    [--tokens-in <n>] [--tokens-out <n>] [--token-breakdown-json <json>]
    [--tool-call-count <n>] [--provider-run-id <id>] [--provider-cost-usd <amount>]
    [--completion <id>] [--review <id>] --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--issue` | string | yes | Issue ID |
| `--role` | string | yes | Agent role that ran |
| `--model` | string | yes | Model used (e.g. `claude-opus-4-7`) |
| `--tokens-in` | int | no | Input tokens |
| `--tokens-out` | int | no | Output tokens |
| `--token-breakdown-json` | JSON string | no | Provider-native token categories, if available |
| `--tool-call-count` | int | no | Number of tool calls made by the run |
| `--provider-run-id` | string | no | Provider request/run identifier |
| `--provider-cost-usd` | decimal | no | Provider-billed cost for reconciliation |
| `--completion` | string | no | Associated completion ID |
| `--review` | string | no | Associated review ID |
| `--by` | string | yes | Acting role |

**Output (stderr confirmation):**
```text
recorded agent run RUN-f6g7h8 on FID-a1b2c3
  role: developer, model: claude-sonnet-4-6
  tokens: 14823 in, 3847 out
  cost: $0.28
```

### `fid run list`

```
fid run list [--open] [--abandoned] [--by <role>] [--epic <id>]
    [--format json|toon|markdown|yaml]
```

| Flag | Type | Description |
|---|---|---|
| `--open` | bool | Only runs with status `open` |
| `--abandoned` | bool | Only runs with status `abandoned` |
| `--by` | string | Filter by run role |
| `--epic` | string | Filter to runs on a given epic's issues |
| `--format` | string | Output format (default: `json`) |

**JSON output (`fid run list --open`):**
```json
[
  {
    "id": "RUN-b2c3d4",
    "issue_id": "FID-a1b2c3",
    "role": "developer",
    "status": "open",
    "branch": "phase-1/core",
    "base_commit": "9f2a1c0",
    "head_commit": "abc123f",
    "dirty": true,
    "worktree_path": "/home/mark/wt/fid-a1b2c3",
    "git_source": "captured",
    "resumed_from_run_id": null,
    "started_at": "2026-05-31T14:30:00Z"
  }
]
```

### `fid run show`

```
fid run show <run-id> [--format json|toon|markdown|yaml]
```

Returns the full run: telemetry (tokens, cost), git context, lifecycle
(`status`, `resumed_from_run_id`, `abandoned_reason_md_path`, `abandoned_at`),
and the linked completion if any.

---

## 13. Stats Commands

### `fid stats tokens`

```
fid stats tokens [--by role|model] [--epic <id>] [--format json|toon|markdown|yaml]
```

**JSON output (`fid stats tokens --by role`):**
```json
[
  {"role": "developer", "total_tokens_in": 1450000, "total_tokens_out": 380000, "run_count": 8},
  {"role": "qa_reviewer", "total_tokens_in": 520000, "total_tokens_out": 95000, "run_count": 3},
  {"role": "tl", "total_tokens_in": 980000, "total_tokens_out": 210000, "run_count": 2}
]
```

### `fid stats cost`

See `05_DESIGN_PLAN.md` §6.4 for JSON output example.

### `fid stats tasks` (v0.2+)

```
fid stats tasks [--by epic] [--include est-vs-actual] [--format json|toon|markdown|yaml]
```

**JSON output (`fid stats tasks --by epic --include est-vs-actual`):**
```json
[
  {
    "epic_id": "FID-e1f2g3",
    "epic_title": "Core persistence layer",
    "total_tasks": 8,
    "completed": 5,
    "in_progress": 2,
    "todo": 1,
    "est_tokens_total": 180000,
    "actual_tokens_total": 145000,
    "est_vs_actual_ratio": 0.81
  }
]
```

### `fid stats attempts` (v0.2+)

```
fid stats attempts [--by epic] [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
[
  {
    "epic_id": "FID-e1f2g3",
    "epic_title": "Core persistence layer",
    "total_tasks": 8,
    "avg_attempts": 1.25,
    "max_attempts": 3,
    "rework_rate": 0.25,
    "tasks_with_rework": 2
  }
]
```

---

## 14. Search Commands (v0.2+)

### `fid search`

```
fid search <term> [--kind issue|decision|escalation|requirement] [--status <state>] [--label <label>] [--epic <id>] [--inbox] [--format json|toon|markdown|yaml]
```

Read-only literal, case-insensitive search over structured projected fields.
Searches entity IDs, issue titles, labels, owner fields, requirement summaries,
decision summaries, escalation summaries, and markdown path fields where
available. It does not scan full markdown body text; use `rg` against project
markdown files for prose-content search.

**JSON output:**
```json
[
  {
    "entity_type": "issue",
    "id": "FID-a1b2c3",
    "title": "Implement journal append with flush",
    "status": "in_progress",
    "match_fields": ["title", "label"],
    "labels": ["journal", "core"]
  },
  {
    "entity_type": "requirement",
    "id": "REQ-d4e5f6",
    "summary": "Durable event log replay",
    "source_doc": "project-docs/specs/01_PRD.md",
    "match_fields": ["summary"]
  }
]
```

**Empty result:**
```json
[]
```

---

## 15. SQL Escape Hatch

### `fid sql`

```
fid sql <query> [--format json|toon|markdown|yaml]
```

Read-only: accepts SELECT, EXPLAIN, PRAGMA query, and WITH. Rejects
INSERT, UPDATE, DELETE, DROP, and all other mutations.

**JSON output:**
```json
[
  {"id": "FID-a1b2c3", "title": "Implement journal append with flush", "status": "in_progress"},
  {"id": "FID-d4e5f6", "title": "Schema migration runner", "status": "todo"}
]
```

**Error (mutation attempt):**
```text
error: fid sql accepts only SELECT, EXPLAIN, PRAGMA query, and WITH statements
  rejected: INSERT INTO issues ...
```

---

## 16. Journal Commands

### `fid journal inspect`

```
fid journal inspect <event-id> [--format json|toon|markdown|yaml]
```

Read-only event inspection. Parses the matching journal event and reports its
location, actor, command metadata, effects, activity entries, canonical hash,
and whether v0.2 `fid event revert` can generate a compensating event.

**JSON output:**
```json
{
  "event_id": "EVT-20260529-000001",
  "event_schema": 1,
  "event_type": "issue.transitioned",
  "journal_file": ".fid/journal/2026-05-29.jsonl",
  "line": 12,
  "hash": "sha256:...",
  "actor": {"role": "developer", "actor": null},
  "command": {
    "argv": ["fid", "issue", "transition", "FID-a1b2c3", "--to", "in_progress", "--by", "developer"],
    "fid_version": "0.1.0"
  },
  "effects": [
    {
      "op": "update",
      "table": "issues",
      "pk": {"id": "FID-a1b2c3"},
      "preconditions": {"status": "todo", "workflow_id": "task-default"},
      "values": {"status": "in_progress"}
    }
  ],
  "revert": {
    "eligible": false,
    "reason": "status transitions require workflow validation",
    "correction": "use fid issue transition if the workflow allows the target state"
  }
}
```

---

## 17. Event Commands (v0.2+)

### `fid event revert` (v0.2+)

```
fid event revert <event-id> --by <role> --reason <text>
```

Appends a compensating event when the target event is explicitly reversible
and the current projected state still matches the target event's produced
state. It never edits or deletes the original journal event.

**Output (stderr confirmation):**
```text
reverted event EVT-20260529-000042
  compensating event: EVT-20260529-000117
```

**Not eligible:**
```text
error: event EVT-20260529-000001 is not eligible for generic revert
  reason: status transitions require workflow validation
  correction: use fid issue transition if the workflow allows the target state
```

---

## 18. Rebuild

### `fid rebuild`

```
fid rebuild --by <role> [--verbose]
```

Replays the full journal into a fresh SQLite database. Automatic projection
catch-up before reads is internal and does not require `--by`; explicit
operator-invoked rebuild does.

| Flag | Type | Description |
|---|---|---|
| `--by` | string | Acting role |
| `--verbose` | bool | Print per-file replay details |

**Output (stderr):**
```text
rebuilding .fid/fid.sqlite from journal...
  replayed 847 events from 12 journal files
  schema version: 3
  upgraded 23 events (v1 → v3)
  elapsed: 0.34s
rebuild complete
```

**Error (corrupt journal):**
```text
error: malformed journal entry at .fid/journal/2026-05-27.jsonl:142
  parse error: unexpected EOF in JSON object
  rebuild stopped — manual repair required
```

---

## 19. Maintenance Commands

### `fid maintenance check`

See `05_DESIGN_PLAN.md` §6.5 for JSON output example.

```
fid maintenance check [--warn-as-error] [--format json|toon|markdown|yaml]
```

| Flag | Type | Description |
|---|---|---|
| `--warn-as-error` | bool | Promote warnings to errors (CI mode) |
| `--format` | string | Output format (default: `json`) |

### `fid maintenance migrate`

```
fid maintenance migrate --by <role> [--to <version>]
```

| Flag | Type | Description |
|---|---|---|
| `--by` | string | Acting role |
| `--to` | int | Target migration version (default: latest) |

**Output (stderr):**
```text
migrated .fid/fid.sqlite: v2 → v3
  applied: 003_add_validation_layer.up.sql
```

### `fid maintenance vacuum`

```
fid maintenance vacuum --by <role>
```

| Flag | Type | Description |
|---|---|---|
| `--by` | string | Acting role |

**Output (stderr):**
```text
vacuum complete — .fid/fid.sqlite: 2.4 MB → 1.8 MB
```

---

## 20. Report Commands (v0.2+)

### `fid report list`

```
fid report list [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
[
  {"name": "epic-status", "description": "Tasks/ACs/gates/findings/escalations summary", "version": "v0.2"},
  {"name": "open-escalations", "description": "All open escalations", "version": "v0.2"},
  {"name": "qa-blockers", "description": "Blocking QA findings", "version": "v0.2"},
  {"name": "ac-matrix", "description": "AC × verifier-role grid", "version": "v0.2"},
  {"name": "stale-claims", "description": "Past-TTL active claims", "version": "v0.2"},
  {"name": "attention", "description": "Items needing human judgment now, grouped by reason", "version": "v0.2"}
]
```

### `fid report <name>`

```
fid report <name> [--epic <id>] [--format json|toon|markdown|yaml]
```

Output shape is report-specific. Named reports are opinionated
wrappers over hand-coded SQL.

### `fid report attention`

```
fid report attention [--epic <id>] [--format json|toon|markdown|yaml]
```

| Flag | Type | Description |
|---|---|---|
| `--epic` | string | Scope to one epic's issues |
| `--format` | string | Output format (default: `json`) |

The pull-based human attention queue: items needing judgment now, grouped by
`reason`. v0.2 reasons: `needs_decision` (open escalation with no resolving
decision; carries `age_hours` and a `stale` flag past a 48h default window),
`needs_approval` (epic wrapped, awaiting human closure), `blocked` (issue
`blocked_by` a cancelled, dead-end blocker), and `human_claimed` (live claim
held `--by human`). Each group declares an `audience`; v0.2 output is
human-default and a `--role`/`--for` filter is a deferred additive fast-follow.
Later reasons (`blocking_finding_open` in v0.3; `validation_failed` /
`validation_carried_forward` in v0.4) appear as additional keys without changing
the contract. Read-only; adds no schema. Group order is most-actionable first
(`needs_decision` → `needs_approval` → `blocked` → `human_claimed`); within a
group, stale and oldest items lead. An empty result is a healthy, zero-exit
state — not an error.

**JSON output:**
```json
{
  "report": "attention",
  "scope": {"epic": null},
  "groups": [
    {"reason": "needs_decision", "audience": ["tl", "human"], "items": [
      {"ref_type": "escalation", "ref_id": "ESC-7c3e9f", "title": "Journal merge strategy unresolved",
       "detail": "open escalation, no resolving decision",
       "since": "2026-05-29T11:00:00Z", "age_hours": 53, "stale": true}]},
    {"reason": "needs_approval", "audience": ["human"], "items": [
      {"ref_type": "issue", "ref_id": "FID-e1f2g3", "title": "Core persistence layer",
       "detail": "epic wrapped; all 8 tasks done; awaiting human closure approval",
       "since": "2026-05-30T18:00:00Z"}]},
    {"reason": "blocked", "audience": ["tl", "human"], "items": [
      {"ref_type": "issue", "ref_id": "FID-d4e5f6", "title": "Schema migration runner",
       "detail": "blocked_by FID-a1b2c3 which is cancelled — re-point or cancel",
       "since": "2026-05-30T09:00:00Z"}]},
    {"reason": "human_claimed", "audience": ["human"], "items": [
      {"ref_type": "issue", "ref_id": "FID-9a8b7c", "title": "Pricing config loader",
       "detail": "claimed --by human; agents will not advance",
       "since": "2026-05-31T08:00:00Z"}]}
  ],
  "summary": {"total": 4, "stale": 1,
              "by_reason": {"needs_decision": 1, "needs_approval": 1, "blocked": 1, "human_claimed": 1}}
}
```

**Markdown output (`--format markdown`):**
```markdown
# Attention — 4 items need human judgment (1 stale)

## needs_decision (1)
- ESC-7c3e9f  Journal merge strategy unresolved
  open escalation, no resolving decision — since 2026-05-29 (53h, stale)

## needs_approval (1)
- FID-e1f2g3  Core persistence layer
  epic wrapped; all 8 tasks done; awaiting human closure approval — since 2026-05-30

## blocked (1)
- FID-d4e5f6  Schema migration runner
  blocked_by FID-a1b2c3 (cancelled) — re-point or cancel — since 2026-05-30

## human_claimed (1)
- FID-9a8b7c  Pricing config loader
  claimed --by human; agents will not advance — since 2026-05-31
```

**Empty result (healthy state, zero exit):**
```text
note: nothing needs human attention
```

---

## 21. Markdown Output Example

**`fid issue get FID-a1b2c3 --format markdown`:**

```markdown
# FID-a1b2c3: Implement journal append with flush

**Kind:** task | **Status:** in_progress | **Priority:** P1
**Epic:** FID-e1f2g3 (Core persistence layer)
**Owner:** developer | **Sequence:** 3
**Spec:** project-docs/docs/phases/phase-1/tasks/FID-a1b2c3.md

## Acceptance Criteria

- [x] AC-d4e5f6: Journal file is flushed to disk before returning success
  - developer: pass (2026-05-27T10:15:00Z)
- [ ] AC-e5f6g7: Partial write leaves journal intact
  - developer: (not verified)

## Labels

core, journal

## Links

- blocks FID-7c3e9f (Schema migration runner)
- implements PRD-ISSUE-001

---
Created: 2026-05-25T09:00:00Z | Updated: 2026-05-27T10:20:00Z
```

---

*Last updated: 2026-05-31*
