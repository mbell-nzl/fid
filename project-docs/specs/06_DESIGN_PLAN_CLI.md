# DESIGN PLAN CLI — fid

> CLI flag specifications and output examples for all commands. The
> parent Design Plan (`05_DESIGN_PLAN.md`) has experience principles,
> output standards, and the command tree overview. This document has
> the per-command detail that developers implement against.
>
> Source: `scaffold/initial-spec.md` §9, §10, plus output examples
> designed to cover every v0.1-scoped command.

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
| `--labels` | string[] | Replace labels (repeatable) |
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

**Output (stderr confirmation):**
```text
completed FID-a1b2c3 (attempt 1)
  completion: TC-7c3e9f
  tokens: 14823
  commits: abc123, def456
```

### `fid issue claim`

```
fid issue claim <id> --by <role> [--ttl <duration>]
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--by` | string (required) | — | Claiming role |
| `--ttl` | duration | `30m` | Lease duration |
| `--note` | string | — | Optional claim context |

**Output (stderr confirmation):**
```text
claimed FID-a1b2c3 by developer (expires 2026-05-27T15:00:00Z)
```

### `fid issue unclaim`

```
fid issue unclaim <id> [--force]
```

| Flag | Type | Description |
|---|---|---|
| `--force` | bool | Take over a live (non-stale) claim |

**Output (stderr confirmation):**
```text
released claim on FID-a1b2c3
```

**Force output:**
```text
force-released claim on FID-a1b2c3 (was held by developer since 2026-05-27T14:30:00Z)
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
```

Composes: claim acquire + transition -> in_progress + open agent_run.

**Output (stderr confirmation):**
```text
dev start FID-a1b2c3
  claimed by developer (expires 2026-05-27T15:00:00Z)
  transitioned: todo → in_progress
  agent run: RUN-b2c3d4 opened
```

### `fid dev gate`

```
fid dev gate <issue-id> --cmd <command> --by developer
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--cmd` | string | yes | Gate command to record |
| `--by` | string | yes | Acting role |

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

Composes: close agent_run + create task_completion + transition ->
dev_complete (validates all_acs_verified, has_completion_report).

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

## 8. Label Commands

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

## 9. Claim Commands

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

## 10. Run Telemetry

### `fid run record`

```
fid run record --issue <id> --role <role> --model <model>
    [--tokens-in <n>] [--tokens-out <n>] [--completion <id>] [--review <id>]
    --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--issue` | string | yes | Issue ID |
| `--role` | string | yes | Agent role that ran |
| `--model` | string | yes | Model used (e.g. `claude-opus-4-7`) |
| `--tokens-in` | int | no | Input tokens |
| `--tokens-out` | int | no | Output tokens |
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

---

## 11. Stats Commands

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

### `fid stats tasks`

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

### `fid stats attempts`

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

## 12. SQL Escape Hatch

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

## 13. Rebuild

### `fid rebuild`

```
fid rebuild [--verbose]
```

Replays the full journal into a fresh SQLite database.

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

## 14. Maintenance Commands

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
fid maintenance migrate [--to <version>]
```

| Flag | Type | Description |
|---|---|---|
| `--to` | int | Target migration version (default: latest) |

**Output (stderr):**
```text
migrated .fid/fid.sqlite: v2 → v3
  applied: 003_add_validation_layer.up.sql
```

### `fid maintenance vacuum`

```
fid maintenance vacuum
```

**Output (stderr):**
```text
vacuum complete — .fid/fid.sqlite: 2.4 MB → 1.8 MB
```

---

## 15. Report Commands (v0.2+)

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
  {"name": "stale-claims", "description": "Past-TTL active claims", "version": "v0.2"}
]
```

### `fid report <name>`

```
fid report <name> [--epic <id>] [--format json|toon|markdown|yaml]
```

Output shape is report-specific. Named reports are opinionated
wrappers over hand-coded SQL.

---

## 16. Markdown Output Example

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

*Last updated: 2026-05-27*
