# CLI CONTRACT PLAN - fid

This document defines the aligned command vocabulary and common transport
contract. Exact leaf request/response schemas are accepted with their D3
vertical slice before implementation.

## 1. Global Contract

```text
fid [--project <path>] [--json] <command> ...
```

- `--project` selects the project root explicitly; default discovery walks
  upward from the working directory and stops at filesystem/Git boundaries.
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
fid capabilities
fid rebuild
fid maintenance check
fid journal inspect <event-id>
```

`capabilities --json` is the first adapter operation. Maintenance reports
journal/projection divergence, invalid references/sequences, stale Claims,
abandoned Attempts, registration conflicts, Artifact state, and activation
disagreement.

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
`--project <path>`, which always selects the fid project root.

Standalone inventory mutations record Actor attribution without claiming
protected Human authority. Execution-bound transitions additionally require
Envelope revision, expected Issue revision, Actor authority, and active Claim.

### Requirement references

```text
fid requirement record-reference|link|get|list
```

This surface stores stable identity, source reference/digest, status, and Issue
links only. Requirement prose is never imported.

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

### Claims, Runs, Stops, and Attention

```text
fid claim acquire|renew|release|list
fid run start|record-telemetry|finish|get|list|recover
fid stop record|resolve|list
fid attention list|resolve
```

No command assumes a Developer or QA role. Actor identity, capabilities,
Envelope grants, Issue, worktree, mutation scope, optional Obligation, and
lease version are explicit machine fields.

### Registered gates and evidence

```text
fid gate register|resolve|run
fid receipt get|list|ingest-attested
fid evidence link|get|evaluate
```

`gate run` accepts a registration reference, subject, Run, repository identity,
and idempotency key. It does not accept replacement command or execution-policy
fields. Live-output framing and terminal JSON are specified in the execution
slice before implementation.

### Review and judgment

```text
fid review record|get|list
fid finding address|list
fid concern record|dispose|list
fid deviation record|list
fid decision record-reference|list
fid human-outcome record|list
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
[Phase 23 input contract](../docs/phases/phase-23/public-interface.md) fixes
bounds, rejection codes and existing Finding disposition readback.

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
exact details `envelope_id`, `decision_class`, `required_capabilities`, and
`affected_actions`; and `source-digest-mismatch` (precondition, exit 4) with
exact details `display_label`, nullable `path`, `expected_digest`, and nullable
`observed_digest`. A missing Source file is reported through the same digest
mismatch code and appends no Event.

Phase 23 Task 10 adds `human-outcome-required` (precondition, exit 4) with
exact details `envelope_id`, `obligation_id`, and `outcome_class`, returned by
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
| `protected-decision-required` | authority / 5 | envelope | `envelope_id`, `decision_class`, `required_capabilities`, `affected_actions` |
| `outcome-not-permitted` | precondition / 4 | null | `obligation_id`, `obligation_type`, `outcome_id`, `result` |
| `evidence-type-not-accepted` | precondition / 4 | null | `obligation_id`, `evidence_type`, `accepted_evidence_types` |
| `source-digest-mismatch` | precondition / 4 | null | `display_label`, `path`, `expected_digest`, `observed_digest` |
| `repository-mismatch` | conflict / 4 | null | `expected`, `observed` |
| `human-outcome-required` | precondition / 4 | envelope | `envelope_id`, `obligation_id`, `outcome_class` |

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

*Last updated: 2026-09-21 (Phase 23 Decision Reference implementation; DEC-239)*
