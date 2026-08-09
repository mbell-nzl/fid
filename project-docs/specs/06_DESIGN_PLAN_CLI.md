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

## 2. JSON Envelope

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
fid requirement sync|list
```

This surface stores stable identity, source reference/digest, status, and Issue
links only. Requirement prose is never imported.

### Work Envelope and projections

```text
fid envelope create|get|revise|activate|complete|cancel
fid next
fid projection context|resume|changes|completion|human-brief
```

Envelope writes consume an approved versioned JSON value. `next` and all
projection commands are derived reads. `human-brief` may explicitly render its
named Human schema as terminal text, Markdown, or HTML; JSON remains canonical.

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
fid concern record|list
fid deviation record|list
fid decision record-reference
fid human-outcome record
```

Review and Outcome inputs are versioned structured records. Protected Human
authority is established through active provenance, not a free-form `--by
human` flag. Decision input is a reference/outcome, never an ADR body.

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

## 5. Typed Reason Codes

Initial cross-cutting codes include:

```text
unsupported-capability
unsupported-schema
activation-incomplete
activation-conflict
authority-missing
authority-forged
protected-decision-required
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

A slice contract fixes its operation names, `result` shapes, reason codes,
closed request/flag shapes, ordering/limits, redaction/unknown-field behavior,
removal inventory, and the exact set of commands that emit a Response Envelope.
Commands scheduled for deletion in a later slice keep their current unenveloped
output and are listed in the owning slice contract with an expiry phase.

*Last updated: 2026-08-09 (Phase 20 accepted slice contract)*
