# DESIGN PLAN CLI — fid

> CLI flag specifications and output examples for all commands. The
> parent Design Plan (`05_DESIGN_PLAN.md`) has experience principles,
> output standards, and the command tree overview. This document has
> the per-command detail that developers implement against.
>
> Source: `scaffold/initial-spec.md` §9, §10, plus output examples
> designed to cover every v0.1-scoped command and the roadmapped v0.2+
> command families.

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
  journal: .fid/journal/2026-06-05.jsonl
  event: EVT-20260605-000001 project.initialized
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
| `--section` | string | — | `acceptance`, `concerns`, `completions`, `links`, `claims`, `gates` |
| `--format` | string | `json` | Output format |
| `--context-pack` | string | — | v0.2+: role-tuned bundle (`developer`, `qa_reviewer`, `qa_lead`, `tl_lite`, `tl`) |

**JSON output:** See `05_DESIGN_PLAN.md` §6.1.

**`--section gates` row shape (raw `gate_runs` history, distinct from
`task_completions.gate_results`):**
```json
[
  {
    "id": "GATE-000401",
    "attempt_number": 1,
    "agent_run_id": "RUN-000401",
    "recorded_by_role": "developer",
    "command": "go test ./...",
    "exit_code": 0,
    "output_excerpt": "ok",
    "ran_at": "2026-05-27T09:20:00Z"
  }
]
```

When an issue has no recorded gates, the section returns an empty list.

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

**Markdown output example (`fid issue list --kind task --order sequence --format markdown`):**
```text
# Issues

- FID-a1b2c3 | task | in_progress | started | epic FID-e1f2g3 | P1 | concerns: 0
  Implement journal append with flush
- FID-d4e5f6 | task | todo | unstarted | epic FID-e1f2g3 | P2 | inbox | concerns: 1
  Schema migration runner
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
| `--epic` | string | v0.2: assign / change / clear a task's epic parent (null↔epic; `--epic ""` clears). A sequential target epic also requires `--sequence` |
| `--inbox` / `--no-inbox` | bool | v0.2: set or clear the inbox flag |
| `--by` | string | Acting role |

At least one field must change. A no-op update is rejected and appends no
journal event.

**Output (stderr confirmation):**
```text
updated FID-a1b2c3:
title "Implement journal append with flush" -> "Implement journal append with fsync"
priority 2 -> 1
```

**Inbox promote (v0.2):** `fid issue update --epic`/`--no-inbox` completes the
capture→promote flow for inbox items — park with `fid issue create --inbox`, list
with `fid issue list --inbox`, then promote. The formal `fid intake` verbs are
deferred sugar over this primitive.

```text
$ fid issue update FID-x --epic FID-e1f2g3 --priority 2 --no-inbox --by tl
updated FID-x:
epic (none) -> FID-e1f2g3
priority 0 -> 2
inbox true -> false
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

The verification role is exactly the `--by` role stored in
`ac_verifications.verifier_role`. `--by human` records a human verification; it
does not count as a developer, QA, or TL verification.

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
      "developer": {"status": "pass", "verified_at": "2026-05-27T10:15:00Z"}
    }
  },
  {
    "id": "AC-e5f6g7",
    "ordinal": 2,
    "text": "Partial write leaves journal intact",
    "latest_verifications": {
      "developer": null
    }
  }
]
```

---

## 4. Ready Queue

### `fid ready`

```
fid ready [--epic <id>] [--explain] [--limit <n>] [--include-blocked]
    [--format json|toon|markdown|yaml]
```

| Flag | Type | Description |
|---|---|---|
| `--epic` | string | Scope to a single epic |
| `--explain` | bool | v0.2: opt into diagnostic mode (ranked candidates, ordering factors, skip reasons) |
| `--limit` | int | v0.2: max ready candidates in `--explain` mode (default 5); also caps the blocked portion of `skipped` |
| `--include-blocked` | bool | v0.2: in `--explain` mode, add dependency-blocked tasks to `skipped` |
| `--format` | string | Output format (default: `json`) |

In **v0.1b**, the bare command implements only `--epic` and `--format`. The
diagnostic flags above remain **v0.2-only**.

**JSON output (default, single task):** See `05_DESIGN_PLAN.md` §6.3. The bare
`fid ready` output shape is frozen — the Workflow Runner depends on it — so the
diagnostic shape below is reachable only via `--explain`.

**Empty state (exit 0, stdout + stderr):**
```json
{"task": null}
```

```text
note: no ready tasks — all 3 tasks in FID-e1f2g3 are blocked or completed
```

**Selection order (bare `fid ready`).** Within a prioritised epic and across the
project-wide pool (prioritised-epic candidates + unattached tasks; sequential
epics contribute only their current actionable task), candidates are ordered
deterministically by:

1. `priority` ascending (PRD-ISSUE-015: `0` is critical, so lowest number wins);
   a task with **no priority is ranked last**.
2. `created_at` ascending (older first).
3. `id` ascending — final deterministic tie-break only.

This is the existing selection rule (`PRD-ISSUE-011/012`); `id` carries no
product meaning beyond guaranteeing a single deterministic winner
(PRD-ISSUE-008). Sequential epics order their tasks by `sequence`
(`PRD-ISSUE-009`), not by this tuple. The v0.2 `--explain` `order_factors` below
surface these same factors.

#### `fid ready --explain` (v0.2)

Diagnostic mode. Returns the ranked ready `candidates` — each with the
`order_factors` the selection rule applied and any `flags` such as
`stale_claim` — and a `skipped` list with a typed `reason` per entry:
`blocked_by`, `claimed`, `sequence_wait`, or `not_actionable`. Explain surfaces
the existing selection rule (`PRD-ISSUE-009/011/012`); it introduces no new
ranking semantics and no numeric score. `claimed` and `sequence_wait` skips show
by default; the potentially numerous `blocked_by` skips require
`--include-blocked`. `--limit` (default 5) caps `candidates` and the blocked
portion of `skipped`.

**JSON output (`fid ready --explain --limit 2 --include-blocked`):**
```json
{
  "scope": {"epic": null},
  "selected": "FID-a1b2c3",
  "candidates": [
    {"id": "FID-a1b2c3", "title": "Implement ID generation", "rank": 1, "ready": true,
     "order_factors": [
       {"factor": "epic_scheduling", "value": "prioritised"},
       {"factor": "priority", "value": 1},
       {"factor": "created_at", "value": "2026-05-25T09:00:00Z"}],
     "flags": []},
    {"id": "FID-d4e5f6", "title": "Schema migration runner", "rank": 2, "ready": true,
     "order_factors": [
       {"factor": "epic_scheduling", "value": "prioritised"},
       {"factor": "priority", "value": 2},
       {"factor": "created_at", "value": "2026-05-25T09:05:00Z"}],
     "flags": ["stale_claim"]}
  ],
  "skipped": [
    {"id": "FID-9a8b7c", "title": "Pricing config loader", "reason": "claimed",
     "detail": "live claim held by qa_reviewer until 2026-05-31T15:00:00Z"},
    {"id": "FID-2b3c4d", "title": "Cycle wiring", "reason": "sequence_wait",
     "detail": "sequential epic FID-f0a1b2; waiting on earlier FID-1a2b3c (todo)"},
    {"id": "FID-7c3e9f", "title": "Wire ID into issue create", "reason": "blocked_by",
     "detail": "blocked_by FID-a1b2c3 (in_progress)", "blockers": ["FID-a1b2c3"]}
  ],
  "summary": {"candidates": 2, "skipped": 3, "shown_blocked": true}
}
```

**Markdown output (`--format markdown`):**
```markdown
# Ready — selected FID-a1b2c3

## Candidates (2)
1. FID-a1b2c3  Implement ID generation   [prioritised, priority 1]
2. FID-d4e5f6  Schema migration runner   [prioritised, priority 2]  ⚠ stale_claim

## Skipped (3)
- FID-9a8b7c  Pricing config loader   claimed by qa_reviewer until 15:00Z
- FID-2b3c4d  Cycle wiring            sequence_wait — earlier FID-1a2b3c not done
- FID-7c3e9f  Wire ID into issue create   blocked_by FID-a1b2c3 (in_progress)
```

**Empty queue under `--explain` (zero exit — the "why is nothing ready?" answer):**
```json
{
  "scope": {"epic": "FID-e1f2g3"},
  "selected": null,
  "candidates": [],
  "skipped": [
    {"id": "FID-1a2b3c", "title": "Journal writer", "reason": "claimed",
     "detail": "live claim held by developer until 2026-05-31T15:00:00Z"},
    {"id": "FID-2b3c4d", "title": "Replayer", "reason": "sequence_wait",
     "detail": "sequential epic; waiting on earlier FID-1a2b3c (in_progress)"},
    {"id": "FID-7c3e9f", "title": "Rebuild command", "reason": "blocked_by",
     "detail": "blocked_by FID-2b3c4d (todo)", "blockers": ["FID-2b3c4d"]}
  ],
  "summary": {"candidates": 0, "skipped": 3, "shown_blocked": true}
}
```
```text
$ fid ready --epic FID-e1f2g3 --explain --include-blocked
note: no ready tasks — 3 tasks blocked or claimed (see skipped[])
```

---

## 5. Developer Composite Commands

### `fid dev start`

```
fid dev start <issue-id> --model <model> --by developer
    [--branch <name>] [--base-commit <sha>] [--worktree <path>]
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--model` | string | yes | Model recorded on the opened base `agent_runs` row |
| `--by` | string | yes | Acting role |
| `--branch` | string | no | v0.2: override the auto-captured git branch |
| `--base-commit` | string | no | v0.2: override the auto-captured base commit |
| `--worktree` | string | no | v0.2: override the auto-captured worktree path |

Composes: claim acquire + transition -> in_progress + open agent_run. v0.1
opens a base `agent_runs` row without git context or lifecycle columns; open
means `ended_at IS NULL`. Git context (branch, base commit, dirty flag,
worktree), lifecycle status, and the override flags above ship in v0.2.

**Output (stderr confirmation, v0.1):**
```text
dev start FID-a1b2c3
  claimed by developer (expires 2026-05-27T15:00:00Z)
  transitioned: todo → in_progress
  agent run: RUN-b2c3d4 opened
```

**Output addition (v0.2 git context):**
```text
  git context: branch phase-1/core @ 9f2a1c0, clean, git_source=captured
```

### `fid dev gate`

```
fid dev gate <issue-id> --cmd <command> --exit-code <n> --by developer
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--cmd` | string | yes | Gate command to record |
| `--exit-code` | int | yes | Reported exit code for the recorded gate run |
| `--output-excerpt` | string | no | Optional reported output excerpt stored on the gate row |
| `--by` | string | yes | Acting role |

Records, but does not execute, an append-only `gate_runs` row for the active
attempt. v0.1 derives the active attempt number from prior completions for the
issue; finish/block commands snapshot matching gate runs into
`task_completions.gate_results`.

**Deprecated from v0.2b (DEC-059).** This surface survives only as an
explicit reported-evidence/import path: rows carry `trust_source=reported`,
are queryable history, and can never satisfy a trusted-gate obligation
(`required_gates_current`). Trusted execution is `fid gate run` (§24). The
v0.2b+ confirmation output adds a visible trust line so reported and trusted
modes never look identical. Flag removal requires a recorded decision.

**Output (stderr confirmation):**
```text
gate recorded on FID-a1b2c3 (attempt 1)
  cmd: go test ./...
  exit_code: 0
  gates_passed: 1, gates_failed: 0
```

**Output (v0.2b+, note the trust line):**
```text
gate recorded on FID-a1b2c3 (attempt 1)
  cmd: go test ./...
  exit_code: 0
  trust: reported (deprecated path — cannot satisfy trusted-gate obligations;
         use fid gate run)
```

**Gate failure:**
```text
gate recorded on FID-a1b2c3 (attempt 1)
  cmd: golangci-lint run ./...
  exit_code: 1
  gates_passed: 0, gates_failed: 1
```

### `fid dev finish`

```
fid dev finish <issue-id> --by developer --completion-md <path>
    [--tokens <n>] [--commits <hashes>] [--files <paths>...]
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--by` | string | yes | Acting role |
| `--completion-md` | string | yes* | Path to completion report. *From v0.2c (DEC-062), required only when the active workflow profile uses `has_completion_report`; structured-completion profiles accept the optional `--summary-md` instead. |
| `--summary-md` | string | no | v0.2c: optional short authored summary/notes path stored on the completion |
| `--tokens` | int | no | Actual tokens used |
| `--commits` | string | no | Comma-separated commit hashes |
| `--files` | string[] | no | Files touched (repeatable) |
| `--human-intervention` | bool | no | Mark that a human intervened during the attempt |
| `--intervention-note` | string | no | Short note stored when `--human-intervention` is set |

Composes: close agent_run + create task_completion + transition ->
dev_complete (validates all_acs_verified plus the profile's completion rule —
`has_completion_report`, or from v0.2c the structured `has_completion` with
its AC/concern/receipt/evidence-obligation checks, PRD-WF-014/015). The run's
`head_commit` and `dirty` flag are refreshed from git (auto-captured, v0.2) at
finish. Under a structured-completion profile the human-readable report is a
generated view (`fid report completion`, §24), not an authored input.

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

Composes: recover mutation authority + open a new agent_run linked to the
task's prior run via `resumed_from_run_id`, carrying forward repository hints.
Any prior still-open run is marked `interrupted`, not intentionally abandoned.
Recovery requires confirmed prior-executor termination, eligible Claim
release/expiry/authorized supersession, and a repository reconciliation
result. `--force` is protected recovery, not proof of termination. fid returns
resume *hints* and records supplied reconciliation facts only; it never
performs git operations or discards worktree state (ADR-0007).

**Output (stderr confirmation):**
```text
dev resume FID-a1b2c3
  prior run RUN-b2c3d4 marked interrupted (superseded by recovery)
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

Validation happens before any journal append. Reject with a non-zero exit and a
clear error when the role is unknown, the issue does not exist, the completion
does not exist, or the completion does not belong to `--on`.

### `fid concern address`

```
fid concern address <concern-id> --in-completion <completion-id> --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `<concern-id>` | string | yes | Concern ID |
| `--in-completion` | string | yes | Completion ID that addressed the concern |
| `--by` | string | yes | Addressing role |

**Output (stderr confirmation):**
```text
addressed concern CON-c3d4e5 in completion TC-9e5f6g
```

Validation happens before any journal append. Reject with a non-zero exit and a
clear error when the role is unknown, the concern does not exist, the concern is
not `open`, the completion does not exist, or the completion does not belong to
the concern's issue.

### `fid concern list`

```
fid concern list [--issue <id>] [--status <status>] [--format json|toon|markdown|yaml]
```

| Flag | Type | Description |
|---|---|---|
| `--issue` | string | Filter by issue ID |
| `--status` | string | `open`, `addressed`, `escalated`, `wontfix` |
| `--format` | string | Output format (default: `json`) |

Filters compose when both `--issue` and `--status` are supplied. A matchless
filter returns an empty collection with exit 0.

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
    "addressed_in_completion_id": null,
    "body_md_path": "project-docs/docs/phases/phase-1/concerns/CON-c3d4e5.md",
    "created_at": "2026-05-27T10:30:00Z",
    "updated_at": "2026-05-27T10:45:00Z"
  }
]
```

The row shape mirrors `fid issue get --section concerns` plus `issue_id`.
`addressed_in_completion_id` is `null` while the concern is `open` and carries the
addressing completion ID once `fid concern address` has run.

Rows are ordered by journal append order, newest first (deterministic and
rebuild-stable; the ordering key is append order, not `created_at`).

---

## 7. Link Commands

### `fid link add`

```
fid link add <from-id> <relation> <to-id> --by <role>
```

Auto-creates the inverse link. Relations: `blocks`, `relates_to`,
`duplicates`, `supersedes`, `implements`, `caused_by`, `discovered_from`.
Endpoint scope is version-gated: v0.1 supports issue ↔ issue and issue →
requirement links; v0.2 adds decision and escalation endpoints.

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

Displays the active permissions matrix. The supported compatibility path reads
the in-code V1 default matrix. A project-configurable `role_permissions` table
is deferred and conditional under ADR-0010; V2 Envelope authority is a
separate generic backend-port contract.

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
    [--completion <id>] --by <role>
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
| `--by` | string | yes | Acting role |

**Output (stderr confirmation):**
```text
recorded agent run RUN-f6g7h8 on FID-a1b2c3
  role: developer, model: claude-sonnet-4-6
  tokens: 14823 in, 3847 out
  cost: $0.102174
```

If the model is missing from `.fid/pricing.toml`, the command still records the
run, prints a warning naming the model on stderr, and shows `cost: unknown`.

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
fid stats tokens [--by role|model|epic] [--epic <id>] [--format json|toon|markdown|yaml]
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
    "id": "PRD-OUT-018",
    "summary": "Literal structured search",
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
    "eligible": true,
    "mode": "inverse",
    "restores": {"status": "todo"},
    "precondition": {"status": "in_progress"},
    "note": "reverts this transition to the prior status; blocked if the issue has since moved on"
  }
}
```

The `revert` verdict is event-type-keyed. Mutable-state events report `mode:
inverse` (restore the prior value under a precondition); append-only audit facts
report `mode: retraction` (withdraw the fact, preserving its row). Non-reversible
events (`issue.created`, `task.completed`, `claim.*`, telemetry, composites)
report `eligible: false` with a `correction` field naming the forward domain path.
In v0.1, the implemented buckets are:

- `issue.transitioned`, `issue.updated`, `issue.concern_addressed`,
  `link.added`, `link.removed` => `mode: inverse`
- `issue.ac_added`, `issue.ac_verified`, `issue.concern_raised`,
  `dev.gate.recorded` => `mode: retraction`
- `project.initialized`, `issue.created`, `label.added`, `issue.completed`,
  `issue.claim.*`, `run.recorded`, `requirement.synced`,
  `dev.started` / `dev.finished` / `dev.blocked` => `eligible: false`

---

## 17. Event Commands (v0.2+)

### `fid event revert` (v0.2+)

```
fid event revert <event-id> --by <role> --reason <text>
```

Appends a compensating event when the target event is reversible. It uses
**inverse** compensation for mutable-state effects (restoring the prior value,
guarded by an effect precondition that the current state still matches what the
event produced) and **retraction** for append-only audit facts (appending a
retraction that excludes the original from current-state derivation while
preserving its row). Eligibility is event-type-keyed — see `fid journal inspect`;
creations, completions, claims, telemetry, and composites are not eligible and
the command returns the forward domain-correction path. It never edits or deletes
the original journal event. There are no bespoke domain-correction verbs:
correction is this one generic primitive (transition-restore replaces
`restore-status`; verification-retraction replaces `ac retract`).

**Output — inverse (transition restore):**
```text
reverted event EVT-20260601-000042
  compensating event: EVT-20260601-000117 (inverse: status dev_complete → in_progress restored)
```

**Output — retraction (append-only fact):**
```text
reverted event EVT-20260601-000087
  compensating event: EVT-20260601-000091 (retraction: ACV-d4e5f6 withdrawn, row preserved)
```

**Not eligible:**
```text
error: event EVT-20260601-000003 (issue.created) is not eligible for revert
  reason: deleting a creation would violate append-only audit history
  correction: cancel the issue with fid issue transition <id> --to cancelled --by <role>
```

---

## 18. Rebuild

### `fid rebuild`

```
fid rebuild --by <role> [--verbose]
```

Replays the full journal into a fresh SQLite database. Automatic projection
catch-up before reads is internal and does not require `--by`; explicit
operator-invoked rebuild does. Both explicit rebuild and automatic catch-up
acquire `.fid/write.lock` before mutating `.fid/fid.sqlite`.

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

**Error (clock skew):**
```text
error: system clock is behind latest journal date 2026-06-05; refusing to append out of order
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

Result and exit contract:

| Findings | `--warn-as-error` | `status` | Exit |
|---|---:|---|---:|
| none | either | `pass` | zero |
| warnings only | false | `warn` | zero |
| warnings only | true | `fail` | non-zero |
| one or more errors | either | `fail` | non-zero |

Promotion does not reclassify warning rows into `errors`; it changes only the
top-level status and process exit. Audit results, including non-zero results, are
written completely to stdout in the selected format. Fatal failures before an
audit result exists use the normal structured stderr path. Findings order by
`check`, `entity_id`, then `message`.

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

Only `tl` and `human` may run this command. `--by` is required; missing,
unknown, or denied roles fail before projection refresh or SQLite mutation.
The command acquires `.fid/write.lock`, refreshes any stale or missing
projection while holding that lock, then vacuums `.fid/fid.sqlite` without
appending a journal event or activity row. It writes no stdout output.

**Output (stderr):**
```text
vacuum complete — .fid/fid.sqlite: 2516582 B → 1887437 B
```

The byte counts are measured from the derived database immediately before and
after `VACUUM`; they vary by workspace.

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
  {"name": "qa-blockers", "description": "Blocking QA findings", "version": "v0.3"},
  {"name": "ac-matrix", "description": "AC × verifier-role grid", "version": "v0.2"},
  {"name": "epic-wrapup-pack", "description": "TL wrap-up readiness bundle", "version": "v0.2"},
  {"name": "task-hygiene", "description": "Task artifact and evidence completeness", "version": "v0.2"},
  {"name": "stale-claims", "description": "Past-TTL active claims", "version": "v0.2"},
  {"name": "attention", "description": "Items needing human judgment now, grouped by reason", "version": "v0.2"},
  {"name": "contract-coverage", "description": "Validation contract coverage gaps", "version": "v0.4"},
  {"name": "assertion-trace", "description": "Requirement-to-outcome trace for one assertion", "version": "v0.4"},
  {"name": "contract-outcomes", "description": "Validation outcomes by assertion and validator", "version": "v0.4"},
  {"name": "recurring-findings", "description": "Findings grouped by recurring structured dimension", "version": "v0.4"},
  {"name": "recurring-concerns", "description": "Concerns grouped by recurring structured dimension", "version": "v0.4"},
  {"name": "recurring-patterns", "description": "Cross-signal recurring patterns ranked as principle candidates", "version": "v0.4"}
]
```

### `fid report <name>`

```
fid report <name> [--epic <id>] [--format json|toon|markdown|yaml]
```

Output shape is report-specific. Named reports are opinionated
wrappers over hand-coded SQL.

### `fid report epic-wrapup-pack` (v0.2+)

```
fid report epic-wrapup-pack <epic-id> [--format json|toon|markdown|yaml]
```

Everything the TL needs before wrap-up: child task state, AC/gate coverage,
open/resolved escalations, QA signals that exist for the current version, and
missing inputs.

**JSON shape:**
```json
{
  "epic_id": "FID-e1f2g3",
  "status": "in_progress",
  "child_tasks": {"done": 5, "cancelled": 1, "open": 0},
  "acceptance_criteria": {"total": 18, "verified_by_required_roles": 18, "gaps": []},
  "gates": {"latest_failures": []},
  "escalations": {"open": [], "resolved": ["ESC-7c3e9f"]},
  "qa": {"blocking_findings": [], "latest_review_id": "QAR-7c3e9f"},
  "missing_inputs": [],
  "recommended_next_commands": ["fid tl wrapup FID-e1f2g3 --by tl"]
}
```

### `fid report task-hygiene` (v0.2+)

```
fid report task-hygiene <task-id> [--format json|toon|markdown|yaml]
```

Checks a task's implementation artifacts for missing spec/completion paths,
unverified ACs, gate gaps, stale claims, and unresolved concerns.

**JSON shape:**
```json
{
  "task_id": "FID-a1b2c3",
  "status": "dev_complete",
  "missing_paths": [],
  "unverified_acceptance_criteria": ["AC-d4e5f6"],
  "gate_gaps": [{"command": "go test -race ./...", "reason": "not_recorded"}],
  "open_concerns": [],
  "stale_claims": [],
  "ready_for_qa": false
}
```

### Validation Contract Reports (v0.4+)

```
fid report contract-coverage <epic-id> [--format json|toon|markdown|yaml]
fid report assertion-trace <assertion-id-or-display-id> [--epic <id>|--contract <id>] [--format json|toon|markdown|yaml]
fid report contract-outcomes <epic-id> [--format json|toon|markdown|yaml]
```

`contract-coverage` reports gaps across requirement/assertion/AC/validator/plan
step coverage.

```json
{
  "epic_id": "FID-e1f2g3",
  "contract_id": "VC-a1b2c3",
  "summary": {"assertions": 12, "covered": 10, "gaps": 2},
  "gaps": [
    {"assertion_id": "ASN-d4e5f6", "display_id": "VAL-004", "missing": ["qa_plan_step"]}
  ]
}
```

`assertion-trace` returns the complete chain for one assertion.

```json
{
  "assertion": {"id": "ASN-d4e5f6", "display_id": "VAL-004"},
  "requirements": ["PRD-OUT-011"],
  "acceptance_criteria": ["AC-a1b2c3"],
  "plan_steps": ["QAS-9a8b7c"],
  "validators": [{"role": "qa_lead", "evidence_method": "source_review"}],
  "outcomes": [{"outcome": "pass", "recorded_by_role": "qa_lead", "evidence_ref": "project-docs/docs/phases/phase-1/qa/final-review.md"}]
}
```

`contract-outcomes` groups latest outcomes by assertion and validator.

```json
{
  "epic_id": "FID-e1f2g3",
  "contract_id": "VC-a1b2c3",
  "summary": {"pass": 10, "fail": 1, "carry_forward": 1, "missing": 0},
  "assertions": [
    {
      "display_id": "VAL-004",
      "validators": [
        {"role": "qa_lead", "evidence_method": "source_review", "latest_outcome": "pass"}
      ]
    }
  ]
}
```

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

### `fid report recurring-patterns` (v0.4+)

Deterministic pattern surfacing over existing quality signals.
`recurring-findings` and `recurring-concerns` group a single signal type;
`recurring-patterns` rolls up across signals. Each entry is a recurring group
with its `dimension`, `key`, `count`, and `exemplars` — a candidate for a human
or TL to turn into a decision, principle, refactor, or template. fid ranks
structured patterns; it never authors the recommendation. No ML/LLM; read-only;
no schema.

**JSON output (`fid report recurring-patterns --epic FID-e1f2g3`):**
```json
{
  "report": "recurring-patterns",
  "scope": {"epic": "FID-e1f2g3", "since": null},
  "patterns": [
    {"dimension": "file_path_prefix", "key": "internal/journal/", "kind": "concern",
     "count": 5, "exemplars": ["CON-c3d4e5", "CON-a1b2c3"],
     "note": "5 concerns touch internal/journal/ — candidate for a principle or refactor"},
    {"dimension": "val_assertion", "key": "VAL-004", "kind": "outcome_fail",
     "count": 3, "exemplars": ["VAL-004@task-3", "VAL-004@task-7"]},
    {"dimension": "finding_severity+area", "key": "blocking/layering", "kind": "finding",
     "count": 4, "exemplars": ["QAF-9a8b7c"]}
  ],
  "summary": {"patterns": 3, "candidates_for_principle": 2}
}
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

## 22. Handoff Commands (v0.3+)

### `fid handoff`

```
fid handoff --role <role> --issue <id> [--since <id>] [--epic <id>]
    [--format json|toon|markdown|yaml]
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--role` | string | yes | Target role to spawn — resolved against the configured handoff templates; a parameter, never a subcommand |
| `--issue` | string | yes | Entity the handoff is about |
| `--since` | string | no | Anchor for delta-oriented roles (e.g. a QA reviewer's task window) |
| `--epic` | string | no | Epic scope for epic-level roles (e.g. QA Lead final review) |
| `--format` | string | no | Output format (default: `json`) |

Composes the role's context pack, a manifest **derived** from the configured
workflow and role permissions, and the role's handoff template into a prompt
plus a machine-readable manifest. Read-only — it acquires no claim and appends no
journal event. `stop_conditions` and `expected_writes` come from the workflow
graph + permissions (so they retrack when a project edits its workflow);
`required_reads` reuse the context pack; `blockers` reuse the
`fid ready --explain` skip data plus open concerns/escalations. Manifest fields
grow by version (`evidence_obligations` lights up in v0.4). The prose prompt is
the manifest's rendering and is also embedded as `prompt_md`.

**JSON output (default):**
```json
{
  "role": "developer",
  "issue": "FID-a1b2c3",
  "template": "handoff-developer@3",
  "manifest": {
    "required_reads": [
      {"kind": "task_spec", "path": "project-docs/docs/phases/phase-1/tasks/task-3.md"},
      {"kind": "validation_contract", "assertions": ["VAL-004", "VAL-005"]},
      {"kind": "context_pack", "role": "developer", "includes": ["DEC-012"]}
    ],
    "commands": [{"source": "template", "cmd": "go test ./... -cover"}],
    "stop_conditions": [
      {"reach_state": "dev_complete", "requires": ["has_completion_report", "all_acs_verified{developer}"]},
      {"reach_state": "failed_gate", "via": "blocker report"}
    ],
    "expected_writes": [
      {"artifact": "completion_report", "from_rule": "has_completion_report"},
      {"artifact": "ac_verification", "from_rule": "all_acs_verified{developer}"},
      {"artifact": "concern", "permitted": true}
    ],
    "blockers": [],
    "evidence_obligations": []
  },
  "prompt_md": "# Developer handoff — FID-a1b2c3 \"Implement ID generation\" ...",
  "source": {"workflow": "task-default", "pack": "developer",
             "derived_from": "workflow + permissions + pack"}
}
```

**Markdown output (`--format markdown`) — the rendered prompt:**
```markdown
# Developer handoff — FID-a1b2c3 "Implement ID generation"

## Required reads
- task spec: project-docs/docs/phases/phase-1/tasks/task-3.md
- validation contract: VAL-004, VAL-005
- context pack (developer): DEC-012 (hash IDs); blocked_by: none

## Commands
- go test ./... -cover

## Stop conditions
- reach `dev_complete` (ACs verified by developer + completion report present), or write a blocker report

## Expected writes
- completion report · AC verifications · concerns
```

---

## 23. Evidence Commands (v0.2b+; amended 2026-07-18, DEC-061)

### `fid evidence check`

```
fid evidence check [--epic <id>] [--issue <id>] [--warn-as-error]
    [--format json|toon|markdown|yaml]
```

| Flag | Type | Description |
|---|---|---|
| `--epic` | string | Scope to one epic's issues |
| `--issue` | string | Scope to one issue |
| `--warn-as-error` | bool | Exit non-zero if any in-scope evidence is `missing` or `stale` (CI / wrap-up gating) |
| `--format` | string | Output format (default: `json`) |

Read-only. Pulled forward from v0.4 to the v0.2b evidence foundation. For each
cited evidence reference on AC verifications, QA findings, and validation
outcomes in scope, reports resolution and — for receipt-backed evidence —
freshness. A type-keyed resolver decides checkability: `receipt` refs resolve
to `gate_receipts` rows and report `current`/`stale` from the receipt's
captured tree identity vs the target tree; `commit` refs resolve via local
git; `source`/`test`/`screenshot`/file refs via filesystem existence;
free-form `command`, `manual`, and `url` evidence are `unverifiable`. The
check never re-runs a command — re-execution is `fid gate run`'s job.

**JSON output:**
```json
{
  "scope": {"epic": "FID-e1f2g3"},
  "summary": {"checked": 15, "resolved": 11, "current": 0, "stale": 1,
              "missing": 1, "unverifiable": 2},
  "evidence": [
    {"source": "workflow_gate", "on": "FID-a1b2c3", "type": "receipt",
     "ref": "RCP-9f8e7d", "gate": "go-test-race@1", "status": "stale",
     "detail": "receipt tree 4c1d... != target tree 8a2b...; re-run fid gate run go-test-race"},
    {"source": "val_outcome", "on": "VAL-004", "type": "commit", "ref": "abc123f",
     "status": "resolved"},
    {"source": "ac_verification", "on": "AC-e5f6g7", "by": "developer", "type": "test",
     "ref": "internal/journal/journal_test.go::TestAppendFlush", "status": "missing",
     "detail": "file not found: internal/journal/journal_test.go"},
    {"source": "qa_finding", "on": "QAF-a1b2c3", "type": "command", "ref": "go test ./...",
     "status": "unverifiable", "detail": "command evidence is not re-checkable; excerpt captured at record time"}
  ]
}
```

**All-resolved (zero exit):**
```text
note: 14 evidence references checked — 11 resolved, 3 unverifiable, 0 missing
```

---

## 24. Roadmap Command Specifications

The sections above define the v0.1 command set plus several detailed later
commands. This section fills the remaining roadmap command families so every
roadmapped CLI surface has an implementable version gate, flags, role
requirement, output shape, and error vocabulary. All write commands also accept
the global v0.2 `--dry-run` flag.

### `fid gate run` (v0.2b+, DEC-059/060)

```
fid gate run <gate-id> [--version <n>] --by <role>
    [--issue <id>] [--run <agent-run-id>]
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--by` | string | yes | Acting role (`gate.run` permission) |
| `--version` | int | no | Pin a registry version; default latest registered |
| `--issue` | string | no | Bind the execution to an issue's active attempt |
| `--run` | string | no | Link the execution to an open agent run |

Resolves and executes the active immutable Gate Definition Registration:
appends a
`gate.execution_started` event, spawns the argv (no shell) with the
definition's cwd/timeout/env allowlist, then appends the terminal
`gate.receipt_recorded` event with status, repository identity, toolchain
fingerprint, and bounded sanitized output. Exit code mirrors the gate:
`0` on `pass`, non-zero otherwise (with `harness_error` distinguished from
gate failure on stderr). Unknown, stale, mismatched, unauthorized, or
substituted registrations are hard errors before any event is appended. The
exact CLI syntax for registration-reference invocation is finalized with the
backend-port command contract; `<gate-id>` remains explanatory shorthand here,
not permission to resolve mutable authoring configuration.

**Output (stderr confirmation):**
```text
gate run go-test-race@1 on FID-a1b2c3 (attempt 1)
  execution: GX-b2c3d4 started
  argv: go test -race ./...   (cwd ., timeout 600s)
  receipt: RCP-9f8e7d  status: pass  exit_code: 0  (41.2s)
  tree: 8a2b4c6 (clean)  trust: executed
  output: 2148 bytes captured, excerpt stored, sha256 1f2e...
```

**Timeout:**
```text
gate run go-test-race@1 on FID-a1b2c3 (attempt 1)
  execution: GX-c3d4e5 started
  receipt: RCP-0a1b2c  status: timed_out after 600s  trust: executed
error: gate go-test-race@1 did not pass (timed_out)
```

### `fid gate list` / `fid gate show` (v0.2b+)

```
fid gate list [--issue <id>] [--gate <gate-id>] [--status <s>] [--abandoned]
    [--format json|toon|markdown|yaml]
fid gate show <execution-or-receipt-id> [--format ...]
```

Read-only. `list` returns executions with their terminal receipts (or
`abandoned` when none); `--abandoned` filters to starts with no receipt.
`show` returns one execution/receipt in full, including hashes, redaction/
truncation flags, and `log_artifact_ref` when retained. Legacy reported
`gate_runs` stay on `fid issue get --section gates`; they are not mixed into
receipt listings.

### `fid gate import` (v0.2b+, attested receipts)

```
fid gate import --receipt <file.json> --by <role>
```

Validates an externally produced receipt against the versioned attestation
contract (schema version, producer identity, declared `trust_source=attested`,
required fields) and appends it linked to its execution-start fact. Rejects —
before any append — receipts with missing/`reported` trust source, unknown
schema major, or an unregistered gate ID: importing can never silently satisfy
a trusted-execution obligation (TRD-SEC-009).

### `fid capabilities` (v0.2b+, DEC-066)

```
fid capabilities [--format json|toon]
```

Read-only machine-readable discovery for workflow kits — no SQLite or `.fid/`
internals required:

```json
{
  "fid_version": "0.2.0",
  "event_schema": 2,
  "projection_schema_version": 5,
  "capabilities": {
    "gate_execution": true,
    "receipts": true,
    "evidence_records": true,
    "structured_completion": false,
    "work_envelope": false
  }
}
```

A consumer that requires an absent capability stops with its own error; fid
never downgrades trusted-receipt requirements to reported evidence.

### `fid report completion` (v0.2c+, DEC-062)

```
fid report completion --task <id> [--attempt <n>]
    [--format json|toon|markdown|html]
```

Generated completion view over structured state: changed commits/files and
repository identity, gate receipts and evidence links (by stable ID — no
snapshot copies), AC/validation coverage and outcomes, telemetry,
interventions, deviations, concerns, decisions, and the optional authored
summary. This view replaces the mandatory authored completion report under
structured-completion profiles.

### `fid report validation-brief` (v0.4c+, DEC-064)

```
fid report validation-brief --epic <id> [--format json|toon|markdown|html]
```

Generated human validation brief: outcomes delivered and user-visible
changes; product/architecture/data/dependency/security decisions;
deviations, unresolved risks, failed or carried-forward obligations, open
escalations; receipt and validation roll-up with spend, retries, and
interventions; and the remaining human-only checks. Contents defined by
PRD-OUT-024; detailed output example at v0.4c phase planning.

### `fid envelope` family (v0.3a+, DEC-063 — sketch)

```
fid envelope apply --file <envelope.yaml> --by <role> [--validate] [--dry-run]
fid envelope status <envelope-id> [--format json|toon|markdown]
fid next [--envelope <id>] [--format json|toon]
```

`apply` ingests a reviewed, versioned work envelope as one create-only journal
event and returns a source-ref→fid-id mapping (carries PRD-OPS-020's
semantics into PRD-OPS-021). `status` reports budgets, circuit-breaker state,
obligation completeness, and attention items. `fid next` returns a derived,
noncanonical Next action over selected fid Issues, with typed non-actionable
reasons (superset of `fid ready --explain`). It does not create or return a
persistent Work Item entity (ADR-0005). Flag-level detail is deliberately
deferred to the V2-aligned envelope slice. That interface must preserve
ADR-0009: partial progress is derived, revisions are immutable and require
prior mutation to be quiescent, unfinished removals retain an authorized
visible disposition, and cancellation does not satisfy completion.

### v0.1 Governance Stub Guard

The v0.1 schema pre-creates `decisions`, `escalations`, and `triage_events` as
reference/view support, but their CLI behavior ships in v0.2. v0.1 registers and
guards exactly these documented commands:

- `fid decision create|get|list`
- `fid escalation create|resolve|get|list`
- `fid triage record|start|verdict`

Each returns only the entity-specific first line and exits non-zero:

```text
error: commands for decision ship in v0.2
```

The version gate runs before argument, flag, role, workspace, DB, or journal
handling, and appends nothing. `fid requirement sync` remains a shipped v0.1
command and is not part of this guard. `v0.2` is the public version; internal
roadmap slice names such as `v0.2b` do not appear in the error.

### Common Future-Command Errors

Commands that target an unshipped version return the standard version-gate
error:

```text
error: commands for qa-plan ship in v0.3
  current schema version: v0.2
  upgrade path: fid maintenance migrate --to 3 --by <role>
```

Write commands fail before validation when the actor is unknown or lacks the
named action permission:

```text
error: role "developer" is not permitted to qa_review.create
  action: qa_review.create
  permitted roles: qa_reviewer, qa_lead, workflow_runner, human
```

### `fid workflow define` (v0.2+)

```
fid workflow define --name <slug> --applies-to <kind>
    --states <csv> --categories <csv> --transitions <csv>
    [--require <edge:rule-json>]... --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--name` | string | yes | Stable workflow slug used for both `workflows.id` and `workflows.name`, e.g. `dev-flow-strict` |
| `--applies-to` | string | yes | Exactly `epic`, `task`, `subtask`, `qa_review`, `escalation`, or `any`; `any` projects as SQL `NULL` |
| `--states` | csv | yes | Ordered state names |
| `--categories` | csv | yes | One category per state: `unstarted`, `started`, `completed`, `cancelled` |
| `--transitions` | csv | yes | Edges as `from->to`; `->`, `-->`, `→`, and `⟶` are accepted equivalently and surrounding whitespace is trimmed |
| `--require` | string[] | no | Transition rule binding `<from>-><to>:<rule>[<json>]`, e.g. `in_progress->dev_complete:all_acs_verified{"roles":["developer"]}`; only the first colon after the edge separates the binding, so JSON colons are preserved |
| `--by` | string | yes | Acting role; default profile permits `tl` and `human` |

Only these implemented rule names are accepted: `all_acs_verified`,
`has_completion_report`, `concerns_addressed`, and
`all_child_tasks_done_or_cancelled`. The following version-scoped rejections
occur before any journal append:

```text
error: workflow rule "required_gates_current" is reserved for v0.2b
error: workflow rule "has_completion" is reserved for v0.2c
error: workflow rule "linked_to_decision" is reserved for v0.2d
```

An unknown rule is rejected naming that rule. A requirement edge must be a
declared or materialized transition; absent JSON projects `params_json` as
`NULL`.

For a workflow declaring exactly one `cancelled`-category state, define
materializes an edge from every non-cancelled state to that state, without
duplicating an explicitly declared edge. `epic`, `task`, `subtask`, and `any`
workflows must declare exactly one cancelled state, so every issue remains
cancellable (PRD-WF-013); zero cancelled states are allowed only for
`qa_review` and `escalation`. More than one cancelled state is rejected.

When an issue is assigned a custom workflow in the subsequent assignment
slice, its initial state is the first declared `unstarted` state, or the first
declared state when none is `unstarted`.

**Custom-state transition boundary (DEC-079):** Phase 19 permits authoring,
assignment, listing, and issue creation for custom workflows, but the current
permissions matrix contains transition actions only for the v0.1 state names.
Consequently a transition into a custom-named state fails with
`authorize action "issue.transition.<state>": no permissions matrix entry` and
does not append a journal event. This is not an implicit permission grant;
custom workflows that reuse v0.1 state names remain transitionable. Generic
actor/capability enforcement in the V2 backend port owns the future unblock
for custom state names. Until then they remain non-transitionable; fid never
infers permission. The former Phase 20 `fid policy set` plan is superseded.

**Output (stderr confirmation):**
```text
defined workflow dev-flow-strict for task
  states: 7
  transitions: 8
  requirements: 2
```

**Errors:** invalid or duplicate workflow slug/state/transition, mismatched
states/categories count, a transition referencing an unknown state, invalid
applicability, an invalid requirement or rule, a cancelled-state policy
violation, or denied `workflow.define` permission. Every rejection occurs
before journal append.

### `fid workflow assign` (v0.2+)

```
fid workflow assign --name <slug> [--kind <kind>] [--project] --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--name` | string | yes | Existing workflow slug |
| `--kind` | string | no | Assign workflow as default for one issue kind |
| `--project` | bool | no | Assign as project default where no kind-specific workflow exists |
| `--by` | string | yes | Acting role; default profile permits `tl` and `human` |

Exactly one of `--kind` or `--project` is required.

**Output:**
```text
assigned workflow dev-flow-strict to kind task
```

**Errors:** missing workflow, incompatible `applies_to`, both/neither assignment
target supplied, or denied `workflow.assign` permission.

### `fid workflow list` (v0.2+)

```
fid workflow list [--kind <kind>] [--active] [--format json|toon|markdown|yaml]
```

Read-only. Lists defined workflows and current assignments.

`--kind <kind>` includes workflows whose `applies_to` is `<kind>` or `null`
(`any`) and excludes workflows for other kinds. `--active` restricts results to
active workflows; the two filters compose. Results are ordered by workflow ID.
Within each `assigned_to` array, kind scopes are ordered by kind and the
project scope is last. An empty result is `[]` in every format with exit 0.

**JSON output:**
```json
[
  {
    "id": "task-default",
    "name": "Default task workflow",
    "applies_to": "task",
    "active": true,
    "states": 9,
    "transitions": 9,
    "assigned_to": [{"scope": "kind", "kind": "task"}]
  }
]
```

The example `states` and `transitions` counts are illustrative; each command
returns the actual projected row counts, including materialized cancel edges.
`applies_to` is `null` for `any`. Kind assignments render as
`{"scope":"kind","kind":"<kind>"}`; project assignments render as
`{"scope":"project"}` with `kind` omitted; and unassigned workflows render
`"assigned_to": []`.

### `fid role create` (deferred/conditional; ADR-0010)

This is a reserved interface sketch, not scheduled scope. Reactivation
requires standalone-fid dogfood evidence, a roadmap decision, and renewed
public-contract review.

```
fid role create <name> [--description <text>] [--human] --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--description` | string | no | Human-readable role purpose |
| `--human` | bool | no | Mark role as human-operated |
| `--by` | string | yes | Acting role; default profile permits `tl` and `human` |

**Output:**
```text
created role release_manager
```

**Errors:** invalid role name, duplicate active role, or denied `role.create`
permission.

### `fid policy set` (deferred/conditional; ADR-0010)

This is a reserved interface sketch, not scheduled scope. It is not the V2
authority mechanism; the backend port generically enforces recorded Envelope
grants.

```
fid policy set --role <role> --action <action> --allow|--deny
    [--conditions <expr>] --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--role` | string | yes | Role being changed |
| `--action` | string | yes | Canonical action key |
| `--allow` / `--deny` | bool | yes | Mutually exclusive policy operation |
| `--conditions` | string | no | Optional condition expression stored verbatim |
| `--by` | string | yes | Acting role; default profile permits `tl` and `human` |

**Output:**
```text
allowed release_manager to issue.transition.done
```

**Errors:** unknown target role, unknown action key, both/neither allow and deny,
or denied `policy.set` permission.

### `fid decision create` (v0.2+)

```
fid decision create --summary <text> --body-md <path> --area <area>
    [--epic <id>] --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--summary` | string | yes | One-line decision summary |
| `--body-md` | path | yes | Markdown body path |
| `--area` | string | yes | Area such as `schema`, `cli`, `workflow`, `security` |
| `--epic` | string | no | Epic the decision belongs to |
| `--by` | string | yes | Acting role; roadmap scope is `tl` or `human` |

**Output:**
```text
created decision DEC-a1b2c3: "Use journal as canonical persistence"
```

**Errors:** missing body file, non-epic `--epic`, duplicate active decision body
hash for same area, or denied `decision.create` permission.

### `fid decision get` (v0.2+)

```
fid decision get <decision-id> [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
{
  "id": "DEC-a1b2c3",
  "epic_id": "FID-e1f2g3",
  "author_role": "tl",
  "summary": "Use journal as canonical persistence",
  "body_md_path": "project-docs/docs/phases/phase-1/decisions/DEC-a1b2c3.md",
  "area": "schema",
  "status": "active",
  "created_at": "2026-06-01T10:00:00Z"
}
```

### `fid decision list` (v0.2+)

```
fid decision list [--epic <id>] [--area <area>] [--status active|superseded]
    [--format json|toon|markdown|yaml]
```

Empty results return `[]` and a stderr note.

### `fid escalation create` (v0.2+)

```
fid escalation create (--from-finding <id>|--from-concern <id>)
    --trigger-class <class> [--summary <text>] --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--from-finding` | string | one source | Source QA finding ID |
| `--from-concern` | string | one source | Source developer concern ID |
| `--trigger-class` | string | yes | Classification such as `arch_violation`, `scope`, `product_decision` |
| `--summary` | string | no | Short title for display/search |
| `--by` | string | yes | Acting role permitted for `escalation.create` |

**Output:**
```text
created escalation ESC-7c3e9f from concern CON-c3d4e5
  trigger_class: arch_violation
```

**Errors:** both/neither source flags, unknown source, finding not blocking
when policy requires blocking-only escalation, duplicate open escalation for
same source, or denied permission.

### `fid escalation resolve` (v0.2+)

```
fid escalation resolve <escalation-id> --decision <decision-id> --by <role>
```

Resolution requires a linked active decision and transitions the escalation from
`open` to `resolved`.

**Output:**
```text
resolved escalation ESC-7c3e9f with decision DEC-a1b2c3
```

**Errors:** escalation already resolved, decision not active, decision not linked
or link cannot be created, or denied `escalation.resolve` permission.

### `fid escalation get` (v0.2+)

```
fid escalation get <escalation-id> [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
{
  "id": "ESC-7c3e9f",
  "status": "open",
  "trigger_class": "arch_violation",
  "raised_by_role": "tl_lite",
  "raised_at": "2026-06-01T11:00:00Z",
  "resolved_at": null,
  "resolution_dec": null
}
```

### `fid escalation list` (v0.2+)

```
fid escalation list [--status open|resolved] [--epic <id>] [--trigger-class <class>]
    [--format json|toon|markdown|yaml]
```

### `fid triage record` (v0.2+)

```
fid triage record --epic <id> --issue <id> --source <source>
    --ref <id> --verdict RESUME|BLOCK --reasoning-md <path>
    [--escalation-trigger <class>] --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--epic` | string | yes | Epic under triage |
| `--issue` | string | yes | Issue affected by the triage source |
| `--source` | string | yes | `dev_concern`, `qa_block`, or `esc_open` |
| `--ref` | string | yes | Source row ID: concern, finding, or escalation |
| `--verdict` | string | yes | `RESUME` or `BLOCK` |
| `--reasoning-md` | path | yes | Triage reasoning artifact |
| `--escalation-trigger` | string | no | Required when `BLOCK` should also create an escalation |
| `--by` | string | yes | Acting role; default profile permits `tl_lite` and `human` |

**Output:**
```text
recorded triage TRG-a1b2c3 on FID-a1b2c3: RESUME
```

**BLOCK with escalation:**
```text
recorded triage TRG-a1b2c3 on FID-a1b2c3: BLOCK
  escalation: ESC-7c3e9f (arch_violation)
```

**Errors:** source/ref mismatch, source does not belong to the epic, missing
reasoning file, `BLOCK` without a valid escalation path when required by policy,
or denied `triage.record` permission.

### `fid triage start` (v0.2+)

```
fid triage start <epic-id> --source <source> --ref <id> --by <role>
    [--format json|toon|markdown|yaml]
```

Read-only preflight for TL-Lite. It validates the source, resolves the affected
issue and relevant context, and appends no journal event.

**JSON output:**
```json
{
  "epic_id": "FID-e1f2g3",
  "source": {"kind": "dev_concern", "id": "CON-c3d4e5"},
  "issue": {"id": "FID-a1b2c3", "title": "Implement ID generation", "status": "dev_complete"},
  "open_escalations": [],
  "suggested_next_commands": [
    "fid triage verdict --verdict RESUME --reasoning-md <path> --by tl_lite",
    "fid triage verdict --verdict BLOCK --escalation-trigger <class> --reasoning-md <path> --by tl_lite"
  ]
}
```

### `fid triage verdict` (v0.2+)

```
fid triage verdict --verdict RESUME|BLOCK --reasoning-md <path>
    [--epic <id>] [--issue <id>] [--source <source>] [--ref <id>]
    [--escalation-trigger <class>] --by <role>
```

Writes the same `triage_events` row as `fid triage record`. When the source
flags are omitted, fid uses the most recent unresolved `triage start` preflight
metadata from the current process invocation context only; it does not persist
triage sessions.

**Output:**
```text
triage verdict on FID-a1b2c3: BLOCK
  triage: TRG-a1b2c3
  escalation: ESC-7c3e9f
```

### `fid tl wrapup` (v0.2+)

```
fid tl wrapup <epic-id> --by <role> [--wrapup-md <path>]
    [--recommendation <text>] [--no-close-escalations]
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--by` | string | yes | Acting role; default profile permits `tl` and `human` |
| `--wrapup-md` | path | no | Wrap-up notes path; v0.4 retrospective creation uses the retrospective command surface |
| `--recommendation` | string | no | Short wrap-up recommendation |
| `--no-close-escalations` | bool | no | Leave already-resolved escalations untouched |

Composes wrap-up note capture, closure of already-resolved escalations where
needed, and epic transition to `wrapped`. It does not create a
`retrospectives` row until v0.4 retrospective support exists, and it does not
perform human approval; `wrapped -> closed` remains a separate
`fid issue transition ... --by human`.

**Output:**
```text
tl wrapup FID-e1f2g3
  wrap-up notes: project-docs/docs/phases/phase-1/wrap-up.md
  closed resolved escalations: 2
  transitioned: in_progress -> wrapped
```

**Errors:** child tasks not done/cancelled, open unresolved escalations, missing
wrap-up artifact when policy requires one, or denied `tl.wrapup`
permission.

### `fid requirement sync`

```
fid requirement sync --source-doc <path> [--type PRD|TRD]
    [--id <requirement-id>] --by <role>
```

v0.1 ships this minimal ingestion command so `fid link add ... implements
PRD-*` can validate that the requirement endpoint exists. It parses canonical
requirement IDs from the source document and upserts pointer rows with
`last_synced_hash`; it stores source pointers and summaries, not copied
requirement prose. v0.1 does not provide drift reports beyond storing the hash.

**Output (v0.1):**
```text
synced requirements from project-docs/specs/01_PRD.md
  created: 3
  updated: 12
  unchanged: 87
```

**Errors:** no parseable requirement IDs, duplicate IDs in source, unknown type,
or denied `requirement.sync` permission.

### `fid requirement get` (v0.2+)

```
fid requirement get <requirement-id> [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
{
  "id": "PRD-ISSUE-001",
  "type": "PRD",
  "summary": "Unified polymorphic issue table",
  "source_doc": "project-docs/specs/01_PRD.md",
  "source_section": "3.2",
  "last_synced_hash": "sha256:...",
  "status": "active",
  "drift": false
}
```

### `fid requirement list` (v0.2+)

```
fid requirement list [--type PRD|TRD] [--status active|deprecated]
    [--drifted] [--source-doc <path>] [--format json|toon|markdown|yaml]
```

### `fid requirement sync` drift reporting (v0.2+)

The v0.2 command keeps the v0.1 flags and adds drift reporting against existing
`last_synced_hash` values.

**Output addition:**
```text
synced requirements from project-docs/specs/01_PRD.md
  created: 3
  updated: 12
  unchanged: 87
  drifted: 2
```

### `fid comment add` (v0.2+; `qa_review` target v0.3+)

```
fid comment add --on <type>/<id> --body-md <path> --by <role>
```

Supported targets are `issue`, `decision`, `escalation`, and `requirement` in
v0.2; `qa_review` is added in v0.3.

**Output:**
```text
added comment COM-a1b2c3 on issue/FID-a1b2c3
```

**Errors:** unsupported target type for the current version, missing target row,
missing body file, or denied `comment.add` permission.

### `fid merge --journal` (v0.2+)

```
fid merge --journal <path> --by <role> [--format json|toon|markdown|yaml]
```

Replays a foreign journal into the current project using event preconditions for
conflict detection. It appends accepted events to the local journal; conflicts
stop the merge and leave later foreign events unapplied.

| Flag | Type | Required | Description |
|---|---|---|---|
| `--journal` | path | yes | Foreign JSONL journal file or directory |
| `--by` | string | yes | Acting operator role; default profile permits `tl`, `workflow_runner`, and `human` |
| `--format` | string | no | Structured merge summary format |

**JSON output:**
```json
{
  "source": "/tmp/recovered/.fid/journal",
  "applied_events": 14,
  "skipped_duplicates": 2,
  "conflict": null
}
```

**Conflict error:**
```text
error: journal merge conflict at EVT-20260601-000042
  entity: issue/FID-a1b2c3
  precondition failed: status expected "todo", found "in_progress"
  applied before conflict: 14
```

> **V2 alignment note (2026-08-08):** The QA command spellings below predate
> ADR-0008 and remain roadmap input, not the approved backend-port interface.
> The replacement contract must create an immutable Review Record bound to an
> Envelope revision and exact repository identity; `changes-required` derives
> same-Issue `builder-rework`, and only a new independent Review Record may
> resolve a Finding. Gate D will align or rescope the commands and publish the
> exact machine response schema.

### `fid qa-review create` (v0.3+)

```
fid qa-review create --scope epic_window|single_task|epic_full
    [--epic <id>] [--task <id>] [--tasks <csv>] [--reviewer <role>]
    [--plan <plan-id>] [--body-md <path>] --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--scope` | string | yes | Review scope |
| `--epic` | string | for `epic_window`, `epic_full` | Epic under review |
| `--task` | string | for `single_task` | Single task under review |
| `--tasks` | csv | for `epic_window` | Covered task IDs |
| `--reviewer` | string | no | Reviewer role stored on the row; defaults from `--by` when compatible |
| `--plan` | string | no | QA plan ID; nullable until v0.4 contract linking |
| `--body-md` | path | no | Review artifact path |
| `--by` | string | yes | Acting role permitted for `qa_review.create` |

**Output:**
```text
created QA review QAR-7c3e9f (epic_window)
  epic: FID-e1f2g3
  tasks: FID-a1b2c3, FID-d4e5f6, FID-g7h8i9
```

**Errors:** scope/flag mismatch, window task outside epic, reviewer role not
`qa_reviewer` or `qa_lead`, plan for another epic, or denied permission.

### `fid qa-review get` (v0.3+)

```
fid qa-review get <review-id> [--format json|toon|markdown|yaml]
```

Returns the review, covered tasks for window scope, findings, linked plan, and
current claim if any.

### `fid qa-review list` (v0.3+)

```
fid qa-review list [--epic <id>] [--task <id>] [--scope <scope>]
    [--reviewer <role>] [--format json|toon|markdown|yaml]
```

### `fid qa-review claim` (v0.3+)

```
fid qa-review claim <review-id> --by <role> [--ttl <duration>] [--force] [--renew]
```

Same lease semantics as `fid issue claim`, with `ref_type = qa_review`.

**Output:**
```text
claimed QAR-7c3e9f by qa_reviewer (expires 2026-06-01T15:00:00Z)
```

### `fid qa-review unclaim` (v0.3+)

```
fid qa-review unclaim <review-id> --by <role>
```

**Output:**
```text
released claim on QAR-7c3e9f by qa_reviewer
```

### `fid qa-finding add` (v0.3+)

```
fid qa-finding add <review-id> --severity blocking|improvement
    [--task <task-id>] --finding-md <path> [--evidence-excerpt <text>]
    --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--severity` | string | yes | `blocking` or `improvement` |
| `--task` | string | no | Task subject within review scope |
| `--finding-md` | path | yes | Finding prose artifact |
| `--evidence-excerpt` | string | no | Literal evidence excerpt |
| `--by` | string | yes | Acting role permitted for `qa_finding.add` |

**Output:**
```text
added QA finding QAF-9a8b7c to QAR-7c3e9f (blocking)
```

**Errors:** task outside review scope, missing finding file, invalid severity,
or denied permission.

### `fid qa-finding list` (v0.3+)

```
fid qa-finding list [--review <id>] [--epic <id>] [--task <id>]
    [--severity blocking|improvement] [--status open|resolved|wontfix|escalated]
    [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
[
  {
    "id": "QAF-9a8b7c",
    "review_id": "QAR-7c3e9f",
    "severity": "blocking",
    "task_id_subject": "FID-a1b2c3",
    "status": "open",
    "escalated_to_esc": null,
    "finding_md": "project-docs/docs/phases/phase-1/qa/QAF-9a8b7c.md",
    "evidence_excerpt": "go test ./... failed in TestJournalReplay"
  }
]
```

### `fid qa start` (v0.3+)

```
fid qa start <epic-or-task-id> --scope epic_window|single_task|epic_full
    [--tasks <csv>] [--plan <plan-id>] [--body-md <path>] --by <role>
```

Composite wrapper over `qa-review create` plus `qa-review claim`. It creates the
review and acquires the review claim atomically.

**Output:**
```text
qa start FID-e1f2g3
  review: QAR-7c3e9f (epic_window)
  claimed by qa_reviewer (expires 2026-06-01T15:00:00Z)
```

### `fid qa finding` (v0.3+)

```
fid qa finding <review-id> --severity blocking|improvement
    [--task <task-id>] --finding-md <path> [--evidence-excerpt <text>]
    --by <role>
```

Composite alias for `fid qa-finding add`, preserving the same validation and
output.

### `fid qa verify` (v0.3+)

```
fid qa verify <issue-id> <ac-id> --by <role> --status <status>
    [--evidence-type <type>] [--evidence-ref <ref>] [--notes <text>]
```

Composite alias for `fid issue ac verify` with QA-role defaults. `--by` must be
a verifier role permitted for the corresponding AC verification action.

### `fid qa finish` (v0.3+)

```
fid qa finish <review-id> --by <role> [--body-md <path>]
```

Completes the review artifact, releases the review claim, and leaves task state
transitions explicit unless the active workflow defines a review-completion
transition rule.

**Output:**
```text
qa finish QAR-7c3e9f
  findings: 1 blocking, 2 improvement
  claim released
```

**Errors:** live claim held by another role, unresolved required review fields,
or denied permission.

### `fid qa-plan create` (v0.3+)

```
fid qa-plan create --epic <id> [--body-md <path>] [--contract <id>]
    --by <role>
```

`--contract` is accepted only once validation contracts ship in v0.4; before
that it is rejected with the version-gate error.

**Output:**
```text
created QA plan QAP-a1b2c3 for FID-e1f2g3
```

**Errors:** non-epic target, contract for another epic, duplicate active plan
for same epic when policy forbids it, or denied `qa_plan.create`.

### `fid qa-plan get` (v0.3+)

```
fid qa-plan get <plan-id> [--format json|toon|markdown|yaml]
```

Returns the plan, ordered steps, linked contract if present, and finalisation
status.

### `fid qa-plan list` (v0.3+)

```
fid qa-plan list [--epic <id>] [--status draft|active|superseded]
    [--format json|toon|markdown|yaml]
```

### `fid qa-plan step add` (v0.3+)

```
fid qa-plan step add <plan-id> --description <text>
    --evidence-method <method> --validator-role <role>
    [--ordinal <n>] --by <role>
```

**Output:**
```text
added QA step QA-01 to QAP-a1b2c3
```

**Errors:** inactive/superseded plan, duplicate ordinal/display ID, unknown
validator role, or denied `qa_plan.step.add`.

### `fid qa-plan step list` (v0.3+)

```
fid qa-plan step list <plan-id> [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
[
  {
    "id": "QAS-d4e5f6",
    "display_id": "QA-01",
    "ordinal": 1,
    "description": "Run repository test suite",
    "evidence_method": "repository_test",
    "validator_role": "qa_lead",
    "assertions": []
  }
]
```

### `fid qa-plan step link` (v0.4+)

```
fid qa-plan step link <step-id> --assertion <assertion-id-or-display-id>
    --by <role>
```

Links a QA plan step to a validation assertion for coverage tracking.

**Output:**
```text
linked QA-01 to VAL-001
```

**Errors:** plan and assertion belong to different epics, duplicate link, or
denied `qa_plan.step.link`.

### `fid qa-plan finalise` (v0.3+; contract rules v0.4+)

```
fid qa-plan finalise <plan-id> --by <role>
```

In v0.3 this validates plan structure and marks the plan active. In v0.4, when
an active contract is linked to the same epic, it also enforces contract
coverage rules from `03_TRD_SCHEMA.md` §12.1.

**Output:**
```text
finalised QA plan QAP-a1b2c3
  steps: 6
  contract coverage: not_applicable (contracts ship in v0.4)
```

**Coverage error (v0.4+):**
```text
error: cannot finalise QA plan QAP-a1b2c3: assertion VAL-004 has no covering step
  rule: every assertion requires a plan step with matching validator role and evidence method
```

### `fid template create` (v0.3+)

```
fid template create --kind <kind> --name <name> --from <path> --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--kind` | string | yes | `task_spec`, `subtask_spec`, `qa_plan`, `qa_review`, `human_checklist`, `decision_record`, `completion_report`, `epic_brief`, `triage_log_entry`, or `retro_lightweight` |
| `--name` | string | yes | Template name unique per kind/version |
| `--from` | path | yes | Markdown template body |
| `--by` | string | yes | Acting role permitted for `template.create` |

Creates the next monotonic version for `(kind, name)`.

**Output:**
```text
created template task_spec/default-feature-task@1
```

**Errors:** unknown kind, missing template body, malformed required-section
metadata, or denied permission.

### `fid template list` (v0.3+)

```
fid template list [--kind <kind>] [--active] [--format json|toon|markdown|yaml]
```

### `fid template show` (v0.3+)

```
fid template show <kind>/<name> [--version <n|latest>]
    [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
{
  "kind": "task_spec",
  "name": "default-feature-task",
  "version": 1,
  "active": true,
  "body_md": "## Context\n..."
}
```

### `fid board` (v0.3+)

```
fid board [--epic <id>] [--label <label>] [--cycle <id>] [--group status-category|workflow]
```

Launches the read-only TUI kanban board. It groups cards by status category by
default; `--group workflow` uses concrete workflow states.

**Errors:** non-interactive terminal without TTY support, unknown epic/label/cycle,
or cycle filter used before v0.5 cycle support.

### `fid contract create` (v0.4+)

```
fid contract create --epic <id> --body-md <path> --by <role>
```

Creates a draft validation contract for one epic.

**Output:**
```text
created validation contract VC-a1b2c3 for FID-e1f2g3 (draft)
```

**Errors:** non-epic target, missing body file, duplicate active contract for
the epic, or denied `contract.create`.

### `fid contract get` (v0.4+)

```
fid contract get (<contract-id>|--epic <id>) [--format json|toon|markdown|yaml]
```

`--epic` returns the active contract for the epic unless `--include-drafts` is
added in a later version.

### `fid contract activate` (v0.4+)

```
fid contract activate <contract-id> --by <role>
```

Marks the contract active and supersedes any previous active contract for the
same epic.

**Output:**
```text
activated validation contract VC-a1b2c3 for FID-e1f2g3
```

**Errors:** contract has no assertions, required validator rows missing, or
denied `contract.activate`.

### `fid assertion add` (v0.4+)

```
fid assertion add <contract-id> --statement <text>
    [--ordinal <n>] [--carry-forward-source <id>] [--out-of-scope-note <text>]
    --by <role>
```

Display IDs are assigned as `VAL-NNN` within the contract.

**Output:**
```text
added assertion VAL-001 to VC-a1b2c3
```

### `fid assertion get` (v0.4+)

```
fid assertion get <assertion-id-or-display-id> [--epic <id>|--contract <id>]
    [--format json|toon|markdown|yaml]
```

`--epic` or `--contract` is required when the argument is a display ID such as
`VAL-001`.

### `fid assertion list` (v0.4+)

```
fid assertion list [--epic <id>|--contract <id>] [--validator-role <role>]
    [--uncovered] [--format json|toon|markdown|yaml]
```

`--uncovered` returns assertions without an AC link, required validator, or
matching QA plan step under the active contract+plan.

### `fid assertion validator add` (v0.4+)

```
fid assertion validator add <assertion-id-or-display-id>
    --role <validator-role> --evidence-method <method>
    [--epic <id>|--contract <id>] --by <role>
```

**Output:**
```text
added validator qa_lead/widget_test to VAL-001
```

**Errors:** unknown validator role, duplicate validator row, display ID without
scope, or denied `assertion.validator.add`.

### `fid assertion link-ac` (v0.4+)

```
fid assertion link-ac <assertion-id-or-display-id> <ac-id>
    [--epic <id>|--contract <id>] --by <role>
```

Links an acceptance criterion to the assertion it helps satisfy.

**Output:**
```text
linked AC-d4e5f6 to VAL-001
```

**Errors:** AC issue is outside the assertion's epic, duplicate link, display ID
without scope, or denied `assertion.link_ac`.

### `fid outcome record` (v0.4+)

```
fid outcome record <assertion-id-or-display-id> --status pass|fail|carry_forward
    [--review <review-id>] [--evidence-method <method>] [--evidence-ref <ref>]
    [--rationale <text>] [--epic <id>|--contract <id>] --by <role>
```

`pass` requires `--evidence-method` and `--evidence-ref`.
`carry_forward` requires `--rationale`.

**Output:**
```text
recorded outcome VAO-a1b2c3 on VAL-001: pass
  evidence: widget_test -> test/orders_test.go::shows_orders
```

**Errors:** missing required evidence/rationale, review outside epic, display ID
without scope, or denied `outcome.record`.

### `fid outcome list` (v0.4+)

```
fid outcome list [--assertion <id-or-display-id>] [--epic <id>]
    [--status pass|fail|carry_forward] [--format json|toon|markdown|yaml]
```

### `fid export` (v0.4+)

```
fid export [--to <path>] [--epic <id>] [--since <anchor>]
    [--layout whet|flat|custom] [--no-banner] --by <role>
```

| Flag | Type | Required | Description |
|---|---|---|---|
| `--to` | path | no | Output root; defaults to `.fid/config.toml` `[export].markdown_root` |
| `--epic` | string | no | Export one epic subtree |
| `--since` | string | no | Export entities changed since event ID or timestamp |
| `--layout` | string | no | Override configured layout |
| `--no-banner` | bool | no | Suppress generated-file banner for this run |
| `--by` | string | yes | Acting role permitted for `export.write` |

**Output:**
```text
exported fid markdown tree to ./docs/
  files written: 42
  unchanged: 18
  layout: whet
```

**Errors:** output path outside allowed project area, unsupported custom layout,
frontmatter generation failure, or denied permission. Export never imports or
reconciles edits from markdown.

### `fid view create` (v0.4+)

```
fid view create --name <name> --query-kind sql (--query <sql>|--query-file <path>)
    [--description <text>] [--scope project|user] --by <role>
```

Only `query_kind=sql` ships initially. Future DSL-backed views require the
conditional query decision.

**Output:**
```text
created view VIEW-a1b2c3: my-blocking-findings
```

**Errors:** mutation query, duplicate view name in scope, unknown query kind, or
denied `view.create`.

### `fid view run` (v0.4+)

```
fid view run <name-or-id> [--format json|toon|markdown|yaml]
```

Read-only. Runs the stored query with the same mutation restrictions as
`fid sql`.

### `fid view list` (v0.4+)

```
fid view list [--scope project|user] [--format json|toon|markdown|yaml]
```

### `fid stats gates` (v0.2+)

```
fid stats gates [--by epic|command] [--epic <id>] [--format json|toon|markdown|yaml]
```

Read-only aggregation over `gate_runs` and completion gate snapshots.

**JSON output:**
```json
[
  {"gate_command": "go test ./...", "runs": 14, "passed": 13, "failed": 1, "pass_rate": 0.93}
]
```

### `fid stats interventions` (v0.2+)

```
fid stats interventions [--epic <id>] [--format json|toon|markdown|yaml]
```

Read-only aggregation over completion attempts marked with human intervention.

**JSON output:**
```json
[
  {"epic_id": "FID-e1f2g3", "tasks_with_intervention": 2, "total_interventions": 3}
]
```

### `fid cycle create` (v0.5+)

```
fid cycle create --name <name> [--start <YYYY-MM-DD>] [--end <YYYY-MM-DD>]
    [--status planned|active|completed] --by <role>
```

**Output:**
```text
created cycle CYC-a1b2c3: "June stabilization"
```

**Errors:** end before start, duplicate active cycle name, cycle support disabled
for the project, or denied `cycle.create`.

### `fid cycle list` (v0.5+)

```
fid cycle list [--status planned|active|completed] [--format json|toon|markdown|yaml]
```

**JSON output:**
```json
[
  {
    "id": "CYC-a1b2c3",
    "name": "June stabilization",
    "start_date": "2026-06-01",
    "end_date": "2026-06-30",
    "status": "active"
  }
]
```

---

*Last updated: 2026-07-11*
