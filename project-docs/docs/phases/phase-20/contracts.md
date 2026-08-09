# Phase 20 Accepted Contracts And Surface Inventory

> Phase: 20 - Clean standalone Issue journey
> Author: Tech Lead
> Status: **Accepted by Human - 2026-08-09**
> Authority: [ADR-0011](../../../../docs/adr/0011-realign-fid-around-the-durable-operational-kernel.md),
> [ADR-0012](../../../../docs/adr/0012-deliver-convergence-as-vertical-scenario-phases.md),
> `../../../specs/10_IMPLEMENTATION_ROADMAP.md` section 5.

This file is the accepted slice contract required by
`06_DESIGN_PLAN_CLI.md` section 1 ("exact leaf request/response schemas are
accepted with their D3 vertical slice before implementation") and by
`03_TRD_SCHEMA.md` section 7 (physical DDL gate). Human approval covers this
contract, its canonical-spec amendments, validation assertions, and task pack.
Developer execution still waits for QA-plan acceptance under the binding V1
workflow.

---

## 1. Scope Boundary

**In scope.** The standalone Issue journey: initialize a project; create,
relate, inspect, update, and transition Issues; explain failure with typed
reasons; rebuild the projection; repeat the journey through canonical JSON.

**Explicitly not in scope** (owning phase in brackets):

| Deferred item | Owner |
|---|---|
| `fid capabilities` and the backend-port adapter surface | Phase 21 |
| Work Envelope, Envelope Revision, Actor authority, Authority Grants | Phase 21 |
| Generic Claims/Runs, Gate Registration, Attempt, Receipt, Evidence | Phase 21 |
| Removal of `fid dev`, `fid policy`, `roles`, `--by`, `gate_runs` | Phase 21 |
| Idempotency keys and the `idempotency_results` table | Phase 21 |
| `--expected-revision` on any surface other than `issue update` / `issue transition` | Phase 21 |
| Typed Stops, Attention, recovery, activation | Phase 22 |
| Reviews, Findings, Concern/Deviation reshape, Human Outcomes | Phase 23 |
| `latest_verifications` role map on acceptance criteria | Phase 23 |
| `issue_dependencies` normalized table and `fid issue dependency` command family | Phase 24 |
| Per-Issue actionable/non-actionable reasons in `fid ready` (PRD-ISSUE-009 second half) | Phase 24 |
| `decisions`, `escalations`, `triage_events` tables; `escalation-default` workflow seed; `issues_enriched` view | Phase 25 |
| Response envelope for legacy commands scheduled for deletion | Phase 21-25 |

Phase 20 **preserves fail-closed V1 authority**. `--by`, `internal/policy`,
the `roles` table, and every existing permission check keep rejecting
unauthorized callers. Phase 20 only gives those rejections a typed machine
shape; Phase 21 removes them after generic authority passes.

---

## 2. Global CLI Contract

```text
fid [--project <path>] [--json] <command> ...
```

The spelling is global, but Phase 20 intentionally has a mixed command tree
during convergence. The flags are valid only for the operations in section 3.
If either flag is explicitly supplied to a legacy command in section 6, fid
rejects before command execution with a concise Human error on stderr. It must
not accept and ignore either flag. Phase 21 removes the affected legacy
commands or migrates the retained replacement, after which the restriction can
be removed.

### 2.1 Project resolution

1. If `--project <path>` is given, the path must be a directory that **directly
   contains** `.fid/`. No upward walk is performed. Otherwise the command fails
   with `project-not-found`.
2. Otherwise discovery starts at the working directory and ascends. The first
   directory containing `.fid/` is the project root.
3. Discovery stops after examining the directory that contains `.git`, and at
   the filesystem root. If no `.fid/` was found within that boundary the command
   fails with `project-not-found`.
4. `fid init` inverts rules 1-3: `--project` names the directory to initialize
   (default: working directory) and `.fid/` must **not** already exist there.
5. Initialization writes the exact target `[project]` keys from
   `04_TRD_CONFIG.md`: generated `id` and presentation-only `issue_prefix`. The
   pre-convergence `prefix` key is not written or accepted.
6. `config.toml` parsing is strict across the accepted section 2 schema:
   top-level `config_version = 1`, `[project]`, `[storage]`,
   `[executor_limits]`, and `[maintenance]` are the only entries. Unknown
   sections, keys, duplicate keys, missing/unsupported versions reject. Temporary
   `.fid/pricing.toml` is separate legacy input and does not widen this parser.
   `artifact_root` must normalize inside `.fid/`; executor limits are positive
   with excerpt no larger than output; maintenance grace is non-negative.
7. Init validates the target and accepted config, writes a sibling staging
   directory, initializes the schema, appends `project.initialized`, rebuilds
   and compares the projection, fsyncs durable files/directories, then atomically
   renames the staging directory to `.fid/`. Any injected failure before rename
   removes staging and leaves no `.fid/`; a failed final rename never replaces
   an existing target. Success cleanup removes no activated state.

### 2.2 Channels

| Mode | stdout | stderr | Exit |
|---|---|---|---|
| Human success | concise terminal text | empty | 0 |
| Human failure | empty | `error: <message>` plus at most two guidance lines | non-zero |
| `--json` success | one Response Envelope, one trailing newline | empty | 0 |
| `--json` failure | empty | one Response Envelope, one trailing newline | non-zero |

No ANSI control sequences, progress text, or Human prose appears on the JSON
channel. Human mutation confirmations move from stderr to **stdout**.

### 2.3 Response Envelope

Success:

```json
{
  "contract_version": "fid-cli.v1",
  "ok": true,
  "operation": "issue.get",
  "result": {},
  "event": null
}
```

Failure:

```json
{
  "contract_version": "fid-cli.v1",
  "ok": false,
  "operation": "issue.transition",
  "error": {
    "code": "state-conflict",
    "category": "conflict",
    "retryable": false,
    "message": "Issue FID-000123 is at revision 4, not 3",
    "subject": {"type": "issue", "id": "FID-000123"},
    "details": {"expected_revision": 3, "actual_revision": 4}
  }
}
```

Rules:

- `contract_version` is the literal `fid-cli.v1`.
- `operation` is present on success and failure and comes from the closed list
  in section 3.
- `result` is an object, never a bare array or scalar. List operations use a
  named collection key.
- `event` is `null` for non-mutating operations. For a mutation that appended
  exactly one Event it is `{"id": "<event-id>", "occurred_at": "<RFC3339 UTC>"}`.
- `error.details` is always present and may be `{}`. `error.subject` is `null`
  when the failure has no single subject.
- No Human message substring is load-bearing. Consumers branch on `code`.
- Success and failure envelopes are mutually exclusive: a success envelope has
  no `error` key, a failure envelope has no `result` or `event` key.

### 2.4 Reason codes

`category` is one of `invalid-input`, `not-found`, `conflict`, `authority`,
`precondition`, `integrity`, `internal`.

| Code | Category | Retryable | Raised when |
|---|---|---|---|
| `invalid-argument` | `invalid-input` | false | Missing, malformed, or mutually exclusive flag or value. |
| `unsupported-kind` | `invalid-input` | false | Issue kind is not `epic`, `task`, or `subtask`. |
| `unsupported-relation` | `invalid-input` | false | Link relation is outside the supported registry. |
| `project-not-found` | `not-found` | false | Project resolution (2.1) found no `.fid/`. |
| `project-already-initialized` | `conflict` | false | `fid init` target already contains `.fid/`. |
| `issue-not-found` | `not-found` | false | Referenced Issue does not exist. |
| `label-not-found` | `not-found` | false | Referenced label does not exist. |
| `link-not-found` | `not-found` | false | Referenced link does not exist. |
| `workflow-not-found` | `not-found` | false | Referenced workflow or assignment scope does not exist. |
| `requirement-not-found` | `not-found` | false | A typed link endpoint names a requirement that does not exist. |
| `duplicate-entity` | `conflict` | false | Label, link, acceptance ordinal, or sequence already exists. |
| `state-conflict` | `conflict` | false | `--expected-revision` does not match the current Issue revision. |
| `lock-unavailable` | `conflict` | **true** | The `.fid/write.lock` mutation lock is held by another process. |
| `hierarchy-invalid` | `precondition` | false | Parent kind or nesting level is illegal for the Issue kind. |
| `dependency-cycle` | `precondition` | false | The requested link would close a dependency cycle. |
| `transition-not-permitted` | `precondition` | false | Target state unknown, edge absent, or a transition precondition is unsatisfied. |
| `authority-missing` | `authority` | false | Required acting identity or authority context is absent. Phase 20 maps missing `--by` here. |
| `authority-unknown-actor` | `authority` | false | The supplied Actor identity is unknown. Phase 20 maps an unknown V1 role here. |
| `authority-denied` | `authority` | false | The known Actor lacks authority for the requested action. |
| `journal-replay-stopped` | `integrity` | false | Corrupt, unreadable, or out-of-order journal halted read or rebuild. |
| `projection-divergence` | `integrity` | false | `maintenance check` found errors or `--warn-as-error` promoted warnings. |
| `result-limit-exceeded` | `integrity` | false | A read would exceed a declared Phase 20 collection limit; nothing is silently truncated. |
| `internal-error` | `internal` | false | Unclassified failure. Its presence in a test is a defect, not a contract. |

The set is closed for Phase 20. Adding a code requires a spec change plus a
test, per ADR-0011.

### 2.5 Error subject and detail objects

`subject` is either `null` or the closed object `{type: string, id: string}`.
The type is the entity named by the failure (`project`, `issue`, `requirement`,
`label`, `link`, or `workflow`). It is `null` for argument, lock, result-limit,
replay, divergence, and internal failures that have no single entity subject.

`details` is a closed object selected by `code`:

| Code | Exact `details` object |
|---|---|
| `invalid-argument`, `unsupported-kind`, `unsupported-relation` | `{argument: string, value: string?, allowed: string[]}` |
| `project-not-found`, `project-already-initialized` | `{path: string}` |
| `issue-not-found`, `requirement-not-found`, `label-not-found`, `workflow-not-found` | `{id: string}` |
| `link-not-found` | `{from_type: string, from_id: string, relation: string, to_type: string, to_id: string}` |
| `duplicate-entity` | `{entity_type: string, identity: {fields: string[], values: string[]}}` (parallel arrays in field-name order) |
| `state-conflict` | `{expected_revision: integer, actual_revision: integer}` |
| `lock-unavailable` | `{lock_path: string}` |
| `hierarchy-invalid` | `{kind: string, parent_id: string?, parent_kind: string?, rule: string}` |
| `dependency-cycle` | `{path: string[]}` |
| `transition-not-permitted` | `{from_state: string, to_state: string, unsatisfied: string[]}` |
| `authority-missing`, `authority-unknown-actor`, `authority-denied` | `{operation: string, actor: string?}` |
| `journal-replay-stopped` | `{journal_file: string, line: integer?, event_id: string?}` |
| `projection-divergence` | the exact `maintenance.check` result object in section 3.5 |
| `result-limit-exceeded` | `{collection: string, limit: integer, observed: integer}` |
| `internal-error` | `{}` |

No detail object includes environment values, config contents, journal payloads,
SQL, file contents, credentials, or stack traces. Paths are normalized caller-
supplied/working paths needed for project resolution, or project-local paths;
unrelated absolute paths are not emitted.

---

## 3. Operations Under `fid-cli.v1`

Exactly these operations emit the Response Envelope in Phase 20. The list is
guarded by a drift test (VAL-035).

| Operation | Command | Mutating |
|---|---|---|
| `project.init` | `fid init` | yes |
| `project.rebuild` | `fid rebuild` | no (projection only) |
| `issue.create` | `fid issue create` | yes |
| `issue.get` | `fid issue get <id>` | no |
| `issue.list` | `fid issue list` | no |
| `issue.update` | `fid issue update <id>` | yes |
| `issue.transition` | `fid issue transition <id>` | yes |
| `issue.ac.add` | `fid issue ac add` | yes |
| `issue.ac.list` | `fid issue ac list <id>` | no |
| `issue.label.add` | `fid issue label add` | yes |
| `issue.label.remove` | `fid issue label remove` | yes |
| `label.add` | `fid label add` | yes |
| `label.list` | `fid label list` | no |
| `link.add` | `fid link add` | yes |
| `link.remove` | `fid link remove` | yes |
| `link.list` | `fid link list <id>` | no |
| `ready.get` | `fid ready` | no |
| `workflow.define` | `fid workflow define` | yes |
| `workflow.assign` | `fid workflow assign` | yes |
| `workflow.list` | `fid workflow list` | no |
| `maintenance.check` | `fid maintenance check` | no |

Every other leaf keeps its established unenveloped Human or JSON mode, loses
`--format` where it currently has that flag, and is marked legacy in section 6.

### 3.1 Schema notation and shared objects

The schemas below are closed: fields not listed are not emitted. `T?` means a
required JSON key whose value is either `T` or `null`; arrays are always present
and are empty rather than `null`. `timestamp` means an RFC 3339 UTC string.

`IssueSummary`:

| Field | Type |
|---|---|
| `id` | string |
| `kind` | `epic` \| `task` \| `subtask` |
| `title` | string |
| `status` | string |
| `status_category` | `unstarted` \| `started` \| `completed` \| `cancelled` |
| `workflow_id` | string |
| `parent_id` | string? |
| `priority` | integer? |
| `sequence` | integer? |
| `spec_md_path` | string? |
| `revision` | integer, minimum 1 |
| `labels` | string[] |
| `created_at` | timestamp |
| `updated_at` | timestamp |

`AcceptanceCriterion`:

| Field | Type |
|---|---|
| `id` | string |
| `ordinal` | integer, minimum 1 |
| `text` | string |

`IssueDetail` contains every `IssueSummary` field plus
`acceptance_criteria: AcceptanceCriterion[]`.

`Label` is `{name: string, colour: string?, description: string?}`. The existing
CLI, JSON, domain, and SQLite vocabulary therefore remains consistent.

`Link` is `{from_type: "issue" | "requirement", from_id: string, relation:
string, to_type: "issue" | "requirement", to_id: string}`. `LinkListItem` is `{direction: "outgoing" |
"incoming", relation: string, entity_type: "issue" | "requirement", entity_id: string,
entity_title: string?}`.

The closed relation registry is `blocks`, `blocked_by`, `relates_to`,
`duplicates`, `duplicated_by`, `supersedes`, `superseded_by`, `implements`,
`implemented_by`, `caused_by`, `causes`, `discovered_from`, `discovered`.
`relates_to` permits either endpoint type; `implements`/`implemented_by` require
Issue/requirement in the named direction; all other relations require two Issue
endpoints. Stored inverse edges use the corresponding paired relation.

`WorkflowSummary` is `{id: string, name: string, applies_to: string?, active:
boolean, states: integer, transitions: integer, assigned_to:
WorkflowAssignment[]}`. `WorkflowAssignment` is `{scope_type: "project" |
"kind", scope_value: string?}`; `scope_value` is `null` for project scope.

`Finding` is `{check: string, message: string, entity_id: string}` where `check`
is exactly one member of the closed error/warning registry in section 3.5.

### 3.2 Closed request contracts

Every operation also accepts the global `--project <path>` and `--json` flags.
No unlisted positional argument or flag is accepted. Required strings reject
empty/whitespace-only values; optional strings use absence for `null` and reject
an explicitly empty value. Unknown flags, enum values, and mutually exclusive
inputs fail as `invalid-argument` before a durable effect.

| Operation | Positional arguments | Operation flags and constraints |
|---|---|---|
| `project.init` | none | `--prefix <string>` optional; no `--by` |
| `project.rebuild` | none | `--by <actor>` required during V1 authority convergence; `--verbose` is Human-only and conflicts with `--json` |
| `issue.create` | none | required `--kind`, `--title`, `--by`; optional `--epic` (Task only), `--parent` (Subtask only), `--priority`, `--sequence`, `--spec-md`, repeatable `--label`, `--scheduling <sequential|prioritised>` (Epic only); parent selectors are mutually exclusive |
| `issue.get` | `<issue-id>` | none |
| `issue.list` | none | optional `--kind`, `--status`, `--status-category`, `--epic`, `--label`, `--priority`, `--has-concerns`, `--order` |
| `issue.update` | `<issue-id>` | required `--by`; at least one of `--title`, `--priority`, `--labels`, `--spec-md`; optional `--expected-revision` |
| `issue.transition` | `<issue-id>` | required `--to`, `--by`; optional `--expected-revision` |
| `issue.ac.add` | `<issue-id>` | required `--text`, `--ordinal`, `--by` |
| `issue.ac.list` | `<issue-id>` | none |
| `issue.label.add`, `issue.label.remove` | `<issue-id> <label>` | required `--by` |
| `label.add` | `<name>` | required `--by`; optional `--colour`, `--description` |
| `label.list` | none | none |
| `link.add`, `link.remove` | `<from-id> <relation> <to-id>` | required `--from-type <issue|requirement>`, `--to-type <issue|requirement>`, `--by`; endpoint type is never inferred from an ID |
| `link.list` | `<entity-id>` | required `--type <issue|requirement>` |
| `ready.get` | none | optional `--epic <issue-id>` |
| `workflow.define` | none | required `--name`, `--applies-to <any|epic|task|subtask>`, `--states`, `--categories`, `--transitions`, `--by`; `--require` is not registered |
| `workflow.assign` | none | required `--name`, `--by`, and exactly one of `--kind <kind>` or `--project-wide` |
| `workflow.list` | none | optional `--kind`, `--active` |
| `maintenance.check` | none | optional `--warn-as-error` |

Removed request surfaces include endpoint-type inference and every item in
section 5.
The command spelling itself is the request capability identifier in Phase 20;
`operation` is its stable dotted identifier. Capability advertisement remains
Phase 21 and is explicitly not claimed by this slice.

### 3.2.1 Closed value domains

| Input | Exact domain |
|---|---|
| `project.init --prefix` / `[project].issue_prefix` | ASCII `^[A-Z][A-Z0-9]{0,15}$`. When omitted, uppercase the selected directory basename and accept it only if it matches; otherwise return `invalid-argument` and require an explicit prefix. |
| `[project].id` | Generated non-empty opaque ASCII identifier, at most 128 bytes. Callers cannot supply it and consumers must not parse it. |
| Issue kind / status category / list order / Epic scheduling | Respectively `epic|task|subtask`; `unstarted|started|completed|cancelled`; `sequence|priority|created|updated`; `sequential|prioritised`. |
| Issue priority | Create/update accepts one integer `0..5`. List accepts a comma-separated set of one to six unique integers from `0..5`; input order carries no meaning and duplicates reject. |
| Issue sequence / AC ordinal | Base-10 integer `1..2147483647`. |
| Expected revision | Base-10 integer `1..9223372036854775807`. |
| Workflow slug | ASCII `^[a-z0-9][a-z0-9-]{0,63}$`. |
| Workflow state / `--status` / `--to` | ASCII `^[a-z][a-z0-9_]{0,63}$`; the state must exist where the operation requires a declared state. |
| Workflow definition lists | `--states`/`--categories` are comma-separated parallel lists; `--transitions` is a comma-separated list whose items contain exactly one `->`. Surrounding Unicode whitespace is trimmed from each token/endpoint, while any blank token rejects. There are `1..64` unique states, exactly one `cancelled` state for an Issue-applicable workflow, and `1..4096` unique edges referencing declared states. |
| Project-relative path input | Valid UTF-8, `1..4096` bytes, relative, normalized, not `.`, and with no lexical `..` escape. If an existing component is a symlink, its resolved path must remain under the owning project root (`artifact_root` is instead contained under `.fid/`). |
| Entity/Actor/label token | Valid UTF-8, `1..255` bytes after trimming, with no NUL or ASCII control character. Entity IDs remain opaque even though their storage token is validated. |
| Title / AC text / description | Valid UTF-8 after trimming, no NUL or ASCII control character; titles are `1..512` bytes, AC text `1..4096`, and optional descriptions `1..4096`. |
| Label colour | Valid UTF-8, `1..64` bytes after trimming, with no NUL or ASCII control character; fid does not interpret its colour syntax. |

Boolean flags accept only `1`, `t`, `T`, `TRUE`, `true`, `True`, `0`, `f`,
`F`, `FALSE`, `false`, or `False` when an explicit value is supplied.
Optional repeated label inputs apply the label-token domain to every value;
`issue.update --labels=` is the one intentional empty-list spelling and clears
all labels. Malformed, out-of-range, duplicate-set, path-escape, and invalid
UTF-8 values return `invalid-argument` before a durable effect.

Workflow definition is deliberately limited to states and transitions in
`fid-cli.v1`. The V1 `--require` authoring flag and its role/assurance-shaped
rules are removed from the canonical surface because PRD-WF-005 forbids kit
assurance choreography in Issue workflow definitions. Bootstrap-seeded V1
requirement rows remain internal compatibility state for legacy `fid dev`
until Phases 21/23 remove their consumers; they are not returned by
`WorkflowSummary`, cannot be authored through a retained operation, and are not
the future generic-precondition or Obligation contract. A stored unknown or
unsupported precondition still causes `issue.transition` to reject before an
Event rather than being ignored.

### 3.3 `issue.get` and `issue.list` results

```json
{
  "issue": {
    "id": "FID-000123",
    "kind": "task",
    "title": "Add link cycle rejection",
    "status": "in_progress",
    "status_category": "started",
    "workflow_id": "task-default",
    "parent_id": "FID-000001",
    "priority": 1,
    "sequence": 2,
    "spec_md_path": "docs/spec.md",
    "revision": 4,
    "labels": ["backend"],
    "acceptance_criteria": [
      {"id": "AC-000401", "ordinal": 1, "text": "Cycles reject."}
    ],
    "created_at": "2026-08-09T10:00:00Z",
    "updated_at": "2026-08-09T11:00:00Z"
  }
}
```

- Absent optional values serialize as `null`, never omitted.
- `labels` and `acceptance_criteria` are always arrays, empty when unset.
- Acceptance criteria carry no `latest_verifications` map in `fid-cli.v1`; that
  role-keyed structure is Phase 23 work.
- Removed relative to today's output: `in_inbox`, `est_tokens`, `due_date`,
  `triage_status`, `completions`.

`issue.get` returns `{issue: IssueDetail}`. `issue.list` returns
`{issues: IssueSummary[], count: integer}` and `count` equals the array length.

Retained filters: `--kind`, `--status`, `--status-category`, `--epic`,
`--label`, `--priority`, `--has-concerns`, `--order`. Removed filter:
`--inbox`.

### 3.4 Mutation results

| Operation | `result` |
|---|---|
| `issue.create` | `{issue: {id: string, kind: string, title: string, status: string, revision: integer}}` |
| `issue.update` | `{issue: {id: string, revision: integer}, changes: {field: string, from: string \| integer \| string[] \| null, to: string \| integer \| string[] \| null}[]}` |
| `issue.transition` | `{issue: {id: string, revision: integer}, from_state: string, to_state: string}` |
| `issue.ac.add` | `{acceptance_criterion: {id: string, issue_id: string, ordinal: integer, text: string}, issue_revision: integer}` |
| `label.add` | `{label: Label}` |
| `issue.label.add` / `.remove` | `{issue_id: string, label: string, issue_revision: integer}` |
| `link.add` / `.remove` | `{link: Link, issue_revisions: {id: string, revision: integer}[]}` |
| `workflow.define` | `{workflow: {id: string, name: string, applies_to: string, state_count: integer, transition_count: integer}}` |
| `workflow.assign` | `{assignment: {scope_type: "project" \| "kind", scope_value: string?, workflow_id: string}}` |
| `project.init` | `{project: {id: string, issue_prefix: string, root: string}, journal_path: string}` |

Every mutating operation returns a non-null `event` object.
`project.root` is the normalized absolute selected directory;
`journal_path` is the normalized project-relative `.fid/journal/<date>.jsonl`
path and never repeats the absolute root.
`issue.update.changes` contains caller-visible requested field changes only;
the automatic `revision` advance is represented by `issue.revision`, not a
synthetic change entry.

The pre-convergence `--owner-role` and `--owner-actor` inputs are not part of
`fid-cli.v1`. Internal default/storage/load paths remain until Phase 21, but
Phase 20 cannot mutate them through the retained public Issue commands.

`workflow assign` accepts exactly one scope selector: `--kind <kind>` or
`--project-wide`. The current local boolean `--project` flag is removed because
it collides with global `--project <path>`. Project selection and assignment
scope can therefore be combined unambiguously:
`fid --project <path> workflow assign --name <id> --project-wide`.

### 3.5 Read results

| Operation | `result` |
|---|---|
| `issue.ac.list` | `{issue_id: string, acceptance_criteria: AcceptanceCriterion[]}` |
| `label.list` | `{labels: Label[], count: integer}` |
| `link.list` | `{entity_type: "issue" | "requirement", entity_id: string, links: LinkListItem[], count: integer}` |
| `ready.get` | `{task: ReadyTask?, explanation: ReadyExplanation?}` |
| `workflow.list` | `{workflows: WorkflowSummary[], count: integer}` |
| `maintenance.check` | `{status: "pass" \| "warn" \| "fail", errors: Finding[], warnings: Finding[], info: {issue_count: integer, claim_count_live: integer, recent_activity_7d: integer}}` |
| `project.rebuild` | `{events_replayed: integer, files_replayed: integer, upgraded_events: integer, schema_version: integer}` |

`ReadyTask` is `{id: string, title: string, status: string, epic_id: string?,
epic_title: string?, sequence: integer?, priority: integer?}`. `ReadyExplanation`
is `{code: "no_ready_tasks" | "sequential_current_not_actionable", scope: "project" | "epic", epic_id: string?,
blocking_task_id: string?, blocking_task_title: string?, blocking_task_status:
string?}`. Exactly one of `task` and `explanation` is non-null.

For `no_ready_tasks`, all blocking-task fields are `null`; `epic_id` is present
only for Epic scope. `sequential_current_not_actionable` is Epic-scoped and all
Epic/blocking-task fields are non-null.

`maintenance.check` returns a success envelope only for `status: "pass"` or
`"warn"`. Findings that produce `status: "fail"` return a failure envelope on
stderr with code `projection-divergence`; its `error.details` is the exact
maintenance result object above, including all `errors`, `warnings`, and
`info`. For `--warn-as-error`, those three members are identical to the
otherwise-successful warn-only result; only `status`, envelope branch, channel,
and process exit differ.

Phase 20's closed check registry is deliberately the useful current subset.
Error checks: `kind_hierarchy`, `dependent_row_kind`, `sequence_scheduling`,
`link_scope`, `link_relation_compatibility`, `link_refs`, `claim_scope`,
`claim_refs`, `activity_refs`, `activity_journal_events`,
`foreign_key_integrity`, `workflow_state`, `workflow_assignment`. Warning
checks: `stale_claims`, `started_without_claim`, `stale_verifications`,
`telemetry_gaps`. `started_without_claim` temporarily uses a fixed 300-second
grace; it must not reinterpret target `abandoned_attempt_grace_seconds` as a V1
Issue/Claim setting. Owning later slices replace/remove checks as their entities
converge; Task 13 adds no future Attempt/Artifact/Activation checks.

`ready.get` returns `task: null` with a populated typed `explanation` when no
Issue is actionable. Both states are asserted (VAL-015).

### 3.6 Concurrency and Issue aggregate revision

`issue update` and `issue transition` accept `--expected-revision <int>`.

- Flag absent: the mutation proceeds without a concurrency check.
- Flag present and equal to the stored revision: the mutation proceeds.
- Flag present and unequal: the mutation fails with `state-conflict`, appends
  nothing, and reports both values in `error.details`.

`issues.revision` is the revision of the complete Issue aggregate, not only the
`issues` table row. `issue.created` creates revision 1. Each of
`issue.updated`, `issue.transitioned`, `issue.ac_added`,
`issue.label.added`, and `issue.label.removed` increments the affected Issue by
exactly one. `link.added` and `link.removed` increment every Issue endpoint
exactly once (one Issue for Issue/requirement links, two for Issue/Issue links),
with `issue_revisions` sorted by Issue ID. A failed operation increments none.
Label-catalog and workflow mutations do not advance Issue revision.
Until Phase 21 deletes the legacy development family, `dev.started`,
`dev.finished`, and `dev.blocked` also advance the one Issue whose status they
mutate by exactly one using the same revision Effect. Those legacy commands do
not accept `--expected-revision` and do not gain a Response Envelope, but they
cannot leave aggregate freshness unchanged after a status mutation.
`updated_at` remains the timestamp of the Issue row's own scalar/status change;
AC, Issue-label, and link Events advance aggregate revision without rewriting
that timestamp. Consumers use revision, not timestamp, for aggregate freshness.

Idempotency keys are Phase 21 work and are **not** accepted in Phase 20.

### 3.7 Ordering, bounds, and compatibility

- `issue.list` supports `sequence`, `priority`, `created`, and `updated`; blank
  means `created`. Created orders by `created_at, id`; updated by
  `updated_at, id`; priority/sequence put null last and then use
  `created_at, id`. `ready.get` applies the sequential/prioritized rules in
  PRD-ISSUE-005 and uses `created_at, id` as deterministic tie-breakers.
- Acceptance criteria order by `ordinal, id`; labels by `name`; links by
  `direction, relation, entity_type, entity_id`; workflows by `name, id`;
  maintenance findings by `check, entity_id, message`; and link mutation
  `issue_revisions` by `id`. `issue.update.changes` orders by `field`.
- One Response Envelope may contain at most 10,000 top-level collection items,
  1,000 acceptance criteria on one Issue, 100 labels on one Issue, and 1,000
  maintenance findings per severity. Mutation enforces the nested Issue limits.
  A read that would exceed a top-level or finding limit fails with
  `result-limit-exceeded`; no collection is silently truncated.
- `fid-cli.v1` output objects are closed. A producer does not add fields or enum
  members under that version; a contract change requires a new accepted version.
  Consumers may reject unknown fields. CLI input has no forward-compatible
  catch-all: unknown flags and values reject.
- No Phase 20 operation is idempotent by caller key. Only `issue.update` and
  `issue.transition` offer compare-and-set, as described in section 3.6.

---

## 4. Accepted Physical DDL

Migrations `0011_p20_remove_issue_fields` and `0012_p20_issue_revision`. They
remain separate so field removal/fixture repair and concurrency semantics fit
separate Developer sessions; Phase 25 consolidates the pre-release history.

`0011_p20_remove_issue_fields` up (logical shape; a table rebuild is permitted):

```sql
ALTER TABLE issues DROP COLUMN in_inbox;
ALTER TABLE issues DROP COLUMN est_tokens;
ALTER TABLE issues DROP COLUMN due_date;
```

Its down migration restores the three columns with their original types and
defaults (`in_inbox BOOLEAN DEFAULT 0`, `est_tokens INTEGER`, `due_date DATE`).
Because SQLite rebuilds the table for the dependent `CHECK`
constraint on `est_tokens`, the migration may be implemented as a
`issues_new` copy in the established `0008` style; the resulting logical shape
above is the contract.

`0012_p20_issue_revision` up/down:

```sql
-- up
ALTER TABLE issues ADD COLUMN revision INTEGER NOT NULL DEFAULT 1;
-- down
ALTER TABLE issues DROP COLUMN revision;
```

Consequences that must be handled in the same task:

- The `kind != 'epic' OR est_tokens IS NULL` CHECK is dropped with the column.
- The `issues_enriched` view is recreated without the removed columns.
- Journal projection appliers reject an Event effect that sets a removed
  column; the failure surfaces as `journal-replay-stopped`.
- The caller never supplies an absolute resulting `revision`. Each successful
  Issue mutation Event records the resulting absolute `revision` in
  `Effect.Values`, the prior revision in `Effect.Prior`, and the same prior
  revision in `Effect.Preconditions`. The existing deterministic `column = ?`
  projection semantics therefore rebuild revision exactly without a new Effect
  operation.

`03_TRD_SCHEMA.md` section 8 records the same DDL as the accepted slice
amendment.

---

## 5. Removal Inventory

Every item below is deleted in Phase 20. Owning tasks: §5.1 and the format
matrix in §5.3 belong to task 16; legacy JSON callers in §5.2 belong to task 14;
remaining flag/test callers in §5.2-§5.3 belong to task 15;
`cmd/fid/issue_section_test.go` and `--section` belong to task 8; §5.4 belongs
to task 17; §5.5 splits between task 5 (schema/domain fields), task 8
(contract-only fields and `--section`), task 2 (global/config naming), Tasks
10-12 (label/link/workflow vocabulary), and task 14 (revert projection).

### 5.1 Generic formatter framework

| Path | Action |
|---|---|
| `internal/format/format.go` + `format_test.go` | delete |
| `internal/format/toon.go` + `toon_test.go` | delete |
| `internal/format/yaml.go` + `yaml_test.go` | delete |
| `internal/format/markdown.go` + `markdown_test.go` | delete |
| `internal/format/markdown_issue.go` | delete |
| `internal/format/markdown_lists.go` | delete |
| `internal/format/markdown_stats.go` | delete |
| `internal/format/markdown_maintenance.go` | delete |
| `internal/format/markdown_workflow.go` + `markdown_workflow_test.go` | delete |
| `internal/format/ready.go` | delete (`ReadyView` is replaced by the `ready.get` result type) |
| `gopkg.in/yaml.v3` in `go.mod` | drop from the require block |

### 5.2 `--format` flag declarations

Sixteen declarations across fourteen files:

`cmd/fid/activity.go`, `cmd/fid/claim.go`, `cmd/fid/concern.go`,
`cmd/fid/issue_ac.go`, `cmd/fid/issue_get_list.go` (×2),
`cmd/fid/journal.go`, `cmd/fid/label.go`, `cmd/fid/link.go`,
`cmd/fid/maintenance.go`, `cmd/fid/policy.go`, `cmd/fid/ready.go`,
`cmd/fid/sql.go`, `cmd/fid/stats.go` (×2), `cmd/fid/workflow.go`.

Around 201 `--format` occurrences across 37 test files are removed or rewritten
with them. Counts are the planning-time inventory; the Developer verifies them
and reports any divergence rather than assuming the plan is exhaustive.

### 5.3 Format-matrix and section tests

| Path | Action |
|---|---|
| `cmd/fid/e2e_v01_format_matrix_test.go` | delete (task 16) |
| `cmd/fid/issue_section_test.go` | delete (task 8) |
| `internal/db/issues_workflow_fk_cli_test.go` | drop its `--format` usage (task 15) |

### 5.4 Promised-command stubs

| Path / surface | Action |
|---|---|
| `cmd/fid/governance_stubs.go` | delete |
| `cmd/fid/governance_stubs_test.go` | delete |
| `newDecisionCommand`, `newEscalationCommand`, `newTriageCommand` registrations in `cmd/fid/root.go` | delete |
| `versionGateError` | delete |

`fid decision`, `fid escalation`, and `fid triage` become unknown commands with
non-zero exit. The `decisions`, `escalations`, and `triage_events` tables stay
until Phase 25 because `internal/maintenance` still audits them.

### 5.5 Retired fields, selectors, and naming collisions

| Surface | Action | Justification |
|---|---|---|
| `issues.in_inbox` column, `--inbox` create flag, `--inbox` list filter, `Issue.InInbox` | delete | No target-model concept; not referenced by any PRD requirement. |
| `issues.est_tokens` column, `--est-tokens` flag, `Issue.EstTokens` | delete | PRD-TELEM-001/003 require observed telemetry, not agent-entered estimates. |
| `issues.due_date` column, `--due-date` flag, `Issue.DueDate` | delete | Scheduling policy is kit-owned (ADR-0004); no retained fid requirement. |
| `Issue.TriageStatus` field and its exposure in `issue.get` | delete from the contract | `Triage Entry` is a retired product term (GLOSSARY, 03 §6). |
| `Issue.Completions` serialization in `issue.get` | delete from the contract | Completion-report-as-state; the writer path stays until Phase 21. |
| `fid issue get --section`, `issues.SectionValue`, `issues.SupportedSections`, `issues.UnsupportedSectionError` | delete | A sub-selector over a Markdown-oriented composite; canonical JSON needs no server-side projection selector. |
| local `fid workflow assign --project` boolean | replace with `--project-wide` | Avoids collision with global `--project <path>` while preserving explicit project-wide assignment. |
| `.fid/config.toml` `[project].prefix` | replace with `[project].issue_prefix` | Aligns implementation with the accepted runtime configuration vocabulary; no compatibility reader is needed pre-dogfood. |
| link endpoint type inferred from an ID prefix | delete; require explicit endpoint types | IDs are opaque and must never drive behavior. |
| `fid workflow define --require` and its V1 role/assurance request grammar | delete from the canonical authoring surface | PRD-WF-005 forbids kit assurance choreography in Issue workflow definitions; internal seeded rows survive only for legacy `fid dev` until Phases 21/23. |
| workflow authoring messages promising rules in `v0.2b` / `v0.2c` / `v0.2d` | delete with the public V1 requirement parser | The historical version roadmap is superseded; owning slices define any future precondition contract. |
| `fid issue create/update --owner-role`, `--owner-actor` | delete public flags; retain internal default/storage/load paths | Role/owner mutation has no forward standalone Issue meaning and would expose hidden V1 state through `fid-cli.v1`. |
| `journal.Verdict`, `RevertVerdict`, inverse/retraction verdict helpers/tests, and `journal inspect` `revert` output | delete | Event correction/revert is evidence-triggered candidate work, not a shipped legacy projection. Raw event lookup remains a temporary diagnostic. |

`Issue.Concerns`, `Issue.Links`, `Issue.Claims`, and `Issue.Gates` are already
`json:"-"`. Their loaders stay; only the removed public fields change.

### 5.6 Not removed in Phase 20

`internal/policy`, `internal/dev`, `fid policy`, `fid dev`, `--by`, the `roles`
table, internal `owner_role` / `owner_actor` storage and loaders, `gate_runs`, `task_completions`,
`ac_verifications`, `concerns`, `requirements`, `internal/stats`,
`internal/sqlquery`, `fid sql`, `cmd/fid/v01_parity_test.go`, and the `v01`
e2e suites. Each is replacement-bound or owned by a later phase.

`.fid/pricing.toml` also remains only because the legacy Phase 25 stats commands
consume it. It is not part of accepted target configuration; Phase 25 removes
the file and its parser with those commands rather than carrying it into the
first supported baseline.

---

## 6. Legacy Command Disposition

These leaves keep their established unenveloped mode until expiry. A legacy read
that currently uses the formatter loses `--format` and emits valid JSON on
stdout with its already-tested semantic fields. Existing legacy mutation/Human
commands keep their current channel behavior. None is a `fid-cli.v1` surface or
has a forward or byte-compatibility contract. Explicit root `--project` or
`--json` rejects before any listed handler runs.

| Command | Temporary output mode | Expiry |
|---|---|---|
| `fid policy show` | unenveloped JSON | Phase 21 |
| `fid dev start`, `fid dev gate`, `fid dev finish`, `fid dev block` | existing Human mode | Phase 21 |
| `fid claim list` | unenveloped JSON | Phase 21 |
| `fid run record`, `fid issue claim`, `fid issue unclaim`, `fid issue complete` | existing Human mode | Phase 21 |
| `fid activity list` | unenveloped JSON | Phase 22 |
| `fid concern list` | unenveloped JSON | Phase 23 |
| `fid concern address`, `fid concern raise` | existing Human mode | Phase 23 |
| `fid requirement sync`, `fid issue ac verify` | existing Human mode | Phase 23 |
| `fid stats tokens`, `fid stats cost` | unenveloped JSON | Phase 25 |
| `fid sql` | unenveloped JSON diagnostic | development-only; reassess Phase 25 |
| `fid journal inspect` | unenveloped JSON event/location inspection; no revert verdict | Phase 25 |
| `fid maintenance vacuum` | existing Human mode | Phase 25 |

A drift test (VAL-035) pins the enveloped set in section 3 against the legacy
set here, so neither list can change silently. After section 5.4 stub removal,
the live tree has exactly 42 executable leaves: 21 retained and 21 legacy.
The three legacy `dev` leaves that mutate Issue status still obey the internal
aggregate-revision Event invariant in sections 3.6 and 7 until deletion.

---

## 7. Event Contract

Phase 20 introduces no new Event family and does not publish the internal
journal Event as a `fid-cli.v1` response. It changes the payload of these
existing Event families: `issue.created`, `issue.updated`,
`issue.transitioned`, `issue.ac_added`, `issue.label.added`,
`issue.label.removed`, `link.added`, `link.removed`, `dev.started`,
`dev.finished`, and `dev.blocked`.

The journal remains `event_schema: 1` during convergence. Its closed Effect
shape is:

```text
Effect = {
  op: "insert" | "update" | "delete" | "activity",
  table: string,
  pk: object,
  values?: object,
  prior?: object,
  preconditions?: object
}
```

`values`, `prior`, and `preconditions` are omitted when empty; they are never
`null`. Every Issue revision advance has this exact additional Effect:

```json
{
  "op": "update",
  "table": "issues",
  "pk": {"id": "FID-000123"},
  "values": {"revision": 5},
  "prior": {"revision": 4},
  "preconditions": {"revision": 4}
}
```

`issue.created` instead includes `"revision": 1` in the existing Issue insert
Effect and has no revision prior/precondition. Updated/transitioned and the
three status-mutating legacy `dev` Events contain one revision Effect. AC and
Issue-label Events contain their existing domain Effect plus one revision
Effect. Link Events contain their existing link Effects plus one revision
Effect for each Issue endpoint, sorted by Issue ID. The resulting revision is
calculated while holding the write lock; it is never supplied by the caller.
`in_inbox`, `est_tokens`, and `due_date` do not appear in newly emitted Effects.

Repository golden fixtures cover each of the eleven affected families, including
an Issue/Issue link with two revision Effects and an Issue/requirement link with
one. Replay rejects a prior/precondition mismatch before applying the Event.

Existing development journals containing the removed keys are **not** migrated.
Per `10_IMPLEMENTATION_ROADMAP.md` section 6 and `04_TRD_CONFIG.md` section 5,
a development workspace is deleted and reinitialized. Repository fixtures that
embed those keys are regenerated in the task that removes the columns, and
rebuild parity is proven against the regenerated fixtures.

---

## 8. Human Directions Applied During Planning

1. Dependency-cycle rejection stays in Phase 20 over the existing links
   projection. Phase 24 still owns normalized dependency storage and commands.
2. `owner_role`, `owner_actor`, `concerns_raised`, and `open_concerns_count`
   remain internal where later replacement work needs them but are absent from
   the first canonical Issue JSON contract.
3. Human success output moves to stdout for enveloped commands; failure remains
   on stderr. Shared test helpers are corrected rather than preserving the
   historical channel mistake.
4. Temporary legacy outputs receive no byte-for-byte compatibility promise.
   They must remain valid JSON and preserve tested meaning until their expiry,
   but no pre-change output capture or golden-byte comparison is required.
5. `journal inspect` remains legacy until the Event baseline and generic Actor
   vocabulary are ready; Phase 20 does not bless V1 role attribution or revert
   verdicts as `fid-cli.v1`.

---

## 9. Slice Acceptance Dispositions

- **Backend-port adapter conformance:** N/A in Phase 20. The slice establishes
  the local CLI contract before the backend port exists; Phase 21 owns the
  shared fake/fid capability suite. Phase 20 must not claim conformance.
- **Capability advertisement update:** N/A by deliberate deferral. There is no
  `fid capabilities` command until Phase 21. VAL-035/044 prove only the exact
  executable CLI/request surface and absence of silent fallback.
- **Recovery evidence:** rebuild equivalence and atomic-init cleanup are required
  by VAL-012/036/042. Execution recovery belongs to Phase 22.
- **Independent consistency review:** required in the sixth mid-phase review
  after Task 17 and again in final QA because this slice changes the canonical
  CLI, config, schema, and internal Event contracts.
