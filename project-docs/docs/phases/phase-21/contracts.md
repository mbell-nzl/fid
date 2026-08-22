# Phase 21 Slice Contract - One-Issue Walking Skeleton

> Status: **PLANNING DRAFT FOR HUMAN APPROVAL**
> Owner: Tech Lead
> Operating mode: Heavy
> Decision authority: DEC-152, DEC-153, ADR-0011, ADR-0012
> Physical schema: `../../../specs/03_TRD_SCHEMA.md` section 8.2

This is the exact Phase 21 implementation boundary. Product implementation is
blocked until the Human accepts this contract, the validation contract, and
the task pack, and a QA Lead supplies a Ready Heavy-mode QA plan. Normative
JSON examples use canonical compact JSON: UTF-8, sorted object keys, preserved
array order, no insignificant whitespace, and a trailing newline on CLI output.

## 1. Versions And Shared Values

- Response Envelope: retained `fid-cli.v1` with Phase 21 leaves and reason
  codes added by this contract.
- Event schema: `fid.event.v2`; v1 development journals are not upgraded.
- Envelope contract: `fid.envelope.v1`.
- Gate Definition and Registration: `fid.gate-definition.v1` and
  `fid.gate-registration.v1`.
- Export archive manifest: `fid.export-manifest.v1`. This is a file format,
  not the unavailable `export.portable.v1` semantic capability.

An Actor is exactly:

```json
{"actor_id":"actor-builder","actor_type":"agent","capabilities":["issue.transition"],"producer":"kit","provenance_ref":"activation:ACT-1/profile:builder","session_ref":null}
```

For kit-mode mutation, callers submit only `actor_id`; fid resolves all other
fields from the active Activation. Unknown Actor fields reject. Standalone
Issue inventory mutations accept an `--actor` attribution string but it is not
trusted kit authority and cannot authorize an Envelope, Claim, Run, Receipt,
Evidence, or completion operation.

A repository identity is exactly:

```json
{"base_oid":"<40hex-or-null>","capture_source":"git","dirty":false,"dirty_digest":null,"head_oid":"<40hex-or-null>","project_id":"<project-id>","tree_oid":"<40hex-or-null>","worktree_id":"<stable-id>"}
```

A source reference is exactly:

```json
{"digest":"sha256:<64hex>","display_label":"Gate policy","external_id":null,"path":"project-docs/gates/unit.json","repository":{...},"section":null,"source_type":"project-file"}
```

An immutable Phase 21 budget is exactly:

```json
{"consumed":null,"exhaustion_action":"unavailable","maximum":3,"name":"attempts","remaining":null,"scope":"envelope","unit":"attempts","warning_threshold":null}
```

`consumed` and `remaining` MUST be JSON null. Phase 21 never updates counters,
warns, exhausts, Stops, raises Attention, or changes execution from a budget.

The complete generic Obligation value is exactly:

```json
{"accepted_evidence_types":["receipt"],"accepted_trust_sources":["fid-executed"],"applicability":"required","artifact_policy":"required-available","derived_status":"unsatisfied","freshness_target":"repository-tree","guarded_action":"envelope.complete","id":"OBL-1","independence":"fid-executed","required_lane":null,"source":{...},"subject":{"id":"ISS-1","type":"issue"},"type":"gate-evidence"}
```

Only `gate-evidence` is accepted. Unknown types reject; `derived_status` is
computed and rejected on mutation input. The value shape remains generic so a
later type does not require an Envelope schema change.

## 2. Capability Availability Contract

Availability is binary and migration-derived. An operation existing for this
walking skeleton does not make a capability available when required semantics
are deferred. Discovery returns rows sorted by `capability_id`:

| Capability ID | Available | Phase 21 limits | Trust source |
|---|---:|---|---|
| `attention.v1` | false | no Attention state or operation | `unavailable:phase-22` |
| `evidence.freshness.v1` | true | repository-tree freshness only | `fid-git-observation` |
| `evidence.gate-definitions.v1` | true | activated immutable definitions; one supported executor policy | `active-kit-activation` |
| `evidence.receipts.v1` | true | fid-executed terminal Receipts only; no attestation | `fid-executed` |
| `export.portable.v1` | false | create/validate only; no import or parity claim | `unavailable:phase-22` |
| `human.outcomes.v1` | false | no Human Outcome state or operation | `unavailable:phase-23` |
| `kit.activation.v1` | false | prepare/stage/commit/finalize/pre-commit abort; no recovery | `unavailable:phase-22` |
| `reviews.coverage.v1` | false | no Review/Finding state or operation | `unavailable:phase-23` |
| `telemetry.trajectory.v1` | false | no generic telemetry ingestion | `unavailable:phase-25` |
| `work.claims.v1` | false | one mutation Claim; no recovery or safe parallelism | `unavailable:phase-22` |
| `workflow.state.v1` | true | one selected primary Issue; no revision/cancel | `fid-journal` |
| `workflow.transitions.v1` | true | registered workflows; execution-bound transition supported | `fid-journal` |

Each discovery row is exactly:

```json
{"available":false,"capability_id":"export.portable.v1","limits":{"create":true,"import":false,"validate":true},"schema_max":1,"schema_min":1,"trust_source":"unavailable:phase-22"}
```

No omitted row, partial/experimental state, caller override, or optimistic
default is permitted. Rebuild must reproduce the same matrix for the same
binary and migration state.

## 3. Idempotency

Every Phase 21 mutation, plus retained retryable Issue create and transition,
requires an `--idempotency-key` of 1..128 printable ASCII characters. The
scope is `<operation>:<project-id>`. fid hashes the canonical mutation request
after Actor resolution but before side effects.

- first key/digest: execute once and durably store response entity ID, Event ID,
  aggregate sequence, and canonical result JSON in the same transaction;
- same key/same digest: return the byte-identical stored result with
  `replayed:true` and append no Event;
- same key/different digest: reject `idempotency_conflict` and append no Event;
- a failure before commit stores nothing; a committed result survives restart
  and rebuild; and
- concurrent same-key calls serialize to one committed mutation.

## 4. Activation And Trust

The only Phase 21 lifecycle is `prepare -> staged -> committed -> finalized`,
or `prepare|staged -> aborted`. One nonterminal Activation blocks another.
`prepare` records bootstrap attribution but grants no protected authority.
`stage` validates and content-addresses the exact authority profile and Gate
Definition components. `commit` atomically selects their identities and is the
point of no return. `finalize` makes that committed identity the active trust
source. Abort after commit rejects `activation_recovery_unavailable`.

The authority-profile component contains a sorted unique Actor set. Each Actor
has immutable ID, type, producer, provenance reference, optional session
reference, and sorted unique capabilities. The Gate Definition component
contains the complete values in section 5. Callers cannot self-declare,
augment, or replace either component. There is no automatic activation
recovery, process respawn, agent respawn, Run continuation, or successor Run.

## 5. Gate Definition And Registration

Registration MUST follow finalized Activation and MUST resolve one activated
Gate Definition by ID and version. It cannot accept execution-policy overrides.
The activated definition is exactly:

```json
{"argv":["go","test","./internal/issues"],"artifact_paths":["coverage.out"],"cwd":".","definition_id":"gate-unit","definition_version":1,"environment":[{"name":"GOCACHE","source":"literal","value":".gocache"},{"name":"HOME","source":"inherit","value":null}],"excerpt_bytes":4096,"network":"allowed","output_bytes":1048576,"redact_literals":[],"sandbox":"none","shell":false,"success_exit_codes":[0],"timeout_ms":120000,"write":"allowed"}
```

Environment entries use only `literal`, `inherit`, or `credential`; credential
entries store a credential name and null value. Names are unique and sorted.
The inherited source is a closed allowlist, not the ambient environment.
Registration validates non-empty structured argv; `shell:false`; contained,
symlink-resolved cwd and Artifact paths; positive configured timeout/output/
excerpt ceilings; non-empty unique success codes; and exact `sandbox:none`,
`network:allowed`, `write:allowed`. Stronger or different values reject
`executor_policy_unsupported`. fid claims no OS, filesystem, network,
privilege, namespace, or platform containment.

A Registration stores immutable definition digest, binding digest, active
Activation/provenance, repository/project/scope binding, and resolved policy.
Same Definition and binding return the existing Registration; changed content
or binding creates a new immutable version. Envelope creation references only
the Registration ID/version.

## 6. Envelope, Authority, And Completion

Phase 21 accepts exactly one Envelope revision, one existing primary Issue,
zero related Issues, one `gate-evidence` Obligation referencing a prior
Registration, zero or more immutable budgets, and sorted Authority Grants.
Revision, cancel, protected Decision, Stop, and Attention inputs reject as
unavailable. Activation requires an active trust source and no active mutation
Claim on the Issue.

An Authority Grant is exactly:

```json
{"actions":["claim.acquire","issue.transition","run.finish","run.start"],"actor_selector":{"actor_id":"actor-builder","required_capabilities":["issue.transition"]},"scope":{"issue_id":"ISS-1"}}
```

Authorization resolves the Actor from the active Activation, then requires an
exact action, capability, Issue, Envelope, and revision match. Missing is deny.
Caller attribution, the deleted role registry, a Human-looking string, or an
inactive Activation confers no authority.

Completion re-reads the active revision in one transaction. It succeeds only
when the selected Issue is in its workflow terminal-success state and every
applicable Obligation is `satisfied`. Otherwise it returns all reasons from
section 10 and changes no state. Completion does not write or consult legacy
`task_completions`.

## 7. Claim, Run, And Execution-Bound Transition

One mutating Claim per Issue is supported. Acquire requires active Envelope,
resolved Actor authority, expected Envelope revision, worktree identity,
declared scope, positive TTL, and idempotency. Renew requires the current Claim
version, same Actor/worktree, and unexpired Claim. Release is idempotent for the
same Claim version/reason. Conflict, expiry, stale version, or wrong Actor does
not authorize mutation.

A Run starts from an active Claim and binds Issue, Envelope/revision, Actor,
Claim/version, repository identity, and immutable budget contracts. Its budget
JSON always exposes null consumption/remaining. Status is `active`,
`succeeded`, `failed`, or `abandoned`; only explicit finish terminalizes it.
There is no recovery, predecessor, automatic continuation, or respawn.

Execution-bound Issue transition requires active Run and Claim, active
Envelope revision, resolved Actor grant, exact expected Issue revision, valid
workflow edge, matching worktree/repository observation, and no unsatisfied
Obligation guarding that action. Checks run in that order and return every
independent reason using section 10 ordering. No Event or revision advance is
allowed on failure.

## 8. Attempts, Output, And Receipts

`gate run` accepts only Registration ID/version, Run ID, Actor reference,
repository observation, and idempotency key. It appends an Attempt before
process spawn. The process uses `exec.CommandContext` with structured argv,
exact cwd, a closed environment assembled from policy, and no shell.

Machine live output is line-delimited `fid.live.v1` JSON. Every frame has
`frame_schema`, `execution_id`, `sequence`, `stream`, and `data`; sequence
starts at 1 and is strictly increasing across both streams. `stream` is
`stdout`, `stderr`, `control`, or `terminal`. Data bytes are decoded as UTF-8
with invalid bytes replaced, credential and literal redactions applied before
framing, and embedded CR/LF/control bytes escaped by JSON. Output is chunked at
4096 sanitized UTF-8 bytes without splitting a code point. Human mode prefixes
sanitized payload with `[stdout]` or `[stderr]`; it prints no raw control byte.

Control frames report only typed metadata: `output_truncated`,
`artifact_rejected`, or `redaction_applied`. When the combined sanitized output
limit is reached, fid stops retaining/emitting process output, emits exactly
one `output_truncated`, continues draining both pipes, and still observes
termination. Redaction applies before live output, excerpt, hash, Receipt,
Artifact, error detail, or Human rendering. Credential values never enter
journal, projection, arguments in errors, or exported data.

Exactly one terminal frame is emitted after observed termination and after the
terminal Receipt commits. Its `data` is the complete canonical `fid-cli.v1`
Response Envelope; no additional JSON line follows it. Terminal outcomes are `succeeded`, `failed`,
`timed_out`, `signaled`, `spawn_failed`, or `output_limit_observed`; the last is
additional metadata and does not replace the observed process result. An
Attempt without observed termination has no Receipt and cannot satisfy
evidence. Timeout sends termination, waits for observed process exit, and
records timeout plus any observed signal/exit. Spawn failure records a Receipt
with null exit/signal. Receipt terminal fields are immutable and mutually
consistent; at most one Receipt exists per Attempt.

## 9. Artifacts And Evidence

Only declared, cwd-contained, symlink-resolved regular files are candidates.
After termination fid revalidates containment, streams bytes into the
content-addressed Artifact root, hashes SHA-256, atomically publishes by digest,
and records media type, size, retention class, and `available`. Missing,
escaping, symlink, directory, device, changed-during-copy, or oversized paths
produce `artifact_rejected` and cannot satisfy `required-available`.

Artifact availability is observed state: `available` may change to `missing`
after a verified content-store check, without changing Artifact, Receipt, or
Evidence identity. Restoration of the same digest may return it to `available`.

Evidence linking accepts only a fid-executed Receipt for the Obligation's exact
Registration, Run, Issue, and target repository identity. A caller-authored
pass flag, imported record, attestation, or command text cannot satisfy it.
Evidence identity is immutable; validity, freshness, and required Artifact
availability are evaluated independently.

## 10. Deterministic Evaluation And Resume

Evaluation returns every applicable reason, deduplicated and sorted by this
fixed precedence, then by `subject_type`, `subject_id`, and `detail` bytewise:

1. `authority_unavailable`
2. `envelope_inactive`
3. `claim_missing`
4. `claim_expired`
5. `run_inactive`
6. `issue_revision_conflict`
7. `workflow_transition_invalid`
8. `obligation_evidence_absent`
9. `receipt_unsuccessful`
10. `evidence_invalid`
11. `evidence_stale`
12. `artifact_unavailable`
13. `issue_not_complete`

Each reason is exactly `{"code":"...","detail":"...","subject_id":"...","subject_type":"..."}`.
Applicable gate evidence is satisfied only when one exact Receipt succeeded,
Evidence is valid and repository-tree fresh, and every required Artifact is
available. Attention is never raised.

`fid projection resume <issue-id>` is a read-only bounded projection. It returns active
Envelope/revision, current Issue/revision, current Claim or null, active Run or
last terminal Run, repository identity, Obligation statuses with all reasons,
and `next_actions` from a closed list sorted by precedence. It reads durable
facts only, does not use chat or process liveness guesses, and never mutates,
recovers, continues, or respawns.

## 11. Export Create And Validate

`export create` emits one deterministic ZIP containing `manifest.json`, the
journal prefix through its declared terminal hash, canonical source-reference
metadata, and available referenced Artifact blobs by digest. Entries are
lexicographically ordered, timestamps/modes/extras are fixed, compression is
fixed, and repeated creation from identical facts is byte-identical. The
manifest contains archive schema, project ID, Event count/terminal hash,
projection schema, capability snapshot, entry digest/size/media type, and
created-from repository identity; it contains no wall-clock creation field.

`export validate` is non-mutating. It checks ZIP safety, exact entry inventory,
canonical manifest, every digest/size, journal chain/schema/terminal hash,
Artifact references, and capability snapshot. It writes no `.fid` state and
does not establish identity/history parity. There is no import command, API,
hidden repair path, or `export.portable.v1` availability claim.

## 12. Canonical Result Shapes

Non-streaming success retains the accepted Response Envelope and is exactly:

```json
{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"envelope.get","result":{}}
```

Collections use `result:{"items":[]}`. Mutation `result` contains the full
created or changed entity, idempotency key, and `replayed`; its `event` contains
ID and occurrence time. Read `event` is null. Failure retains the exact Phase
20 `fid-cli.v1` mutually-exclusive error envelope. Optional values are explicit
null; unavailable fields are omitted only when this contract says they do not
exist. Human output is a rendering of the same result and cannot add semantics.

### 12.1 Accepted Phase 21 command leaves

Retained Phase 20 leaves remain accepted, with mutation attribution renamed
from `--by` to `--actor` and retryable mutations adding
`--idempotency-key`. Phase 21 adds exactly:

```text
fid capabilities
fid activation prepare|component|commit|abort|finalize|status
fid gate register|resolve|run
fid envelope create|get|activate|complete
fid claim acquire|renew|release|list
fid run start|finish|get|list
fid receipt get|list
fid evidence link|get|evaluate
fid artifact verify
fid projection resume
fid export create|validate
```

Complex mutation values enter through exactly one `--input <path|->` JSON
document; `-` means stdin. Reference/read leaves use positional opaque IDs.
Every kit mutation also requires `--actor <actor-id>` and
`--idempotency-key <key>`; update-like leaves require the contract's expected
revision/version flag. `gate run --json` alone uses the live NDJSON channel;
all other `--json` leaves emit one Response Envelope. `stats --by` is renamed
to `stats --group-by`. There is no `activation recover`, `envelope revise`,
`envelope cancel`, `run recover`, telemetry, Stop, Attention, Receipt ingest,
Review, Human Outcome, or export import leaf in Phase 21.

## 13. Phase 21 Reason Codes

In addition to retained Phase 20 codes, the closed codes are:
`activation_conflict`, `activation_invalid`,
`activation_recovery_unavailable`, `actor_unknown`, `artifact_unavailable`,
`authority_denied`, `authority_unavailable`, `capability_unavailable`,
`claim_conflict`, `claim_expired`, `claim_stale`, `envelope_inactive`,
`evidence_invalid`, `evidence_stale`, `executor_policy_unsupported`,
`idempotency_conflict`, `obligation_unsatisfied`, `registration_mismatch`,
`run_inactive`, and `unsafe_path`. Errors remain typed stderr JSON with stable
exit classification and no partial success payload.

## 14. Event Schema V2 And Fixtures

Every new Event and every retained mutation converted in Phase 21 uses exactly:

```json
{"actor":{"actor_id":"actor-builder","actor_type":"agent","producer":"kit","provenance_ref":"activation:ACT-1/profile:builder","session_ref":null},"aggregate":{"id":"ISS-1","sequence":2,"type":"issue"},"causation_id":null,"command":"issue transition","correlation_id":"RUN-1","effects":[{"op":"update","pk":{"id":"ISS-1"},"preconditions":{"revision":1},"prior":{"revision":1,"status":"todo"},"table":"issues","values":{"revision":2,"status":"in_progress"}}],"event_id":"EVT-2","event_schema":2,"event_type":"issue.transitioned","idempotency":{"key":"key-1","request_digest":"sha256:<64hex>","scope":"issue.transition:<project-id>"},"occurred_at":"2026-08-22T00:00:00Z","repository":{...}}
```

Exact rules: all top-level keys are required; nullable values are explicit;
`effects` is a non-empty ordered array; aggregate sequence begins at 1 and has
no gap; Actor capabilities are not copied into Events; repository is null only
for a non-repository bootstrap fact; unknown keys reject; canonical fixtures
replace all affected v1 fixtures. The accepted Phase 21 families are:

| Aggregate | Event types | Projection effects |
|---|---|---|
| activation | `activation.prepared`, `.component_staged`, `.committed`, `.finalized`, `.aborted` | activation/component/profile/definition insert or state update |
| gate-registration | `gate.registered` | registration insert |
| envelope | `envelope.created`, `.activated`, `.completed` | envelope/revision/source/selection/grant/obligation/budget insert or status update |
| claim | `claim.acquired`, `.renewed`, `.released` | claim insert or version/time/status update |
| run | `run.started`, `.finished` | run/budget-binding insert or terminal update |
| attempt | `gate.execution_started` | Attempt insert |
| receipt | `gate.receipt_recorded` | Receipt and Artifact inserts |
| evidence | `evidence.linked`, `artifact.availability_changed` | Evidence/link insert or Artifact availability update |
| retained mutation | retained Issue/workflow/label/link/requirement/maintenance families | retained v2 effects plus resolved Actor/idempotency/aggregate metadata |

Receipt recording is one Event containing the terminal Receipt and every
accepted Artifact effect. Evidence linking is a separate mutation/Event. An
Attempt is always durable before spawn. No Phase 21 Event family exists for
Attention, Stop, budget consumption, recovery, respawn, import, Review,
Finding, Human Outcome, or telemetry.

## 15. Removal, Frozen State, And Deferral Inventory

After its named replacement passes:

| Remove/rename | Replacement or disposition |
|---|---|
| mutation `--by` flags | `--actor`; `stats --by` becomes `--group-by` |
| `internal/policy`, `cmd/fid policy`, `roles` | activation Actor profile plus Envelope grants |
| `internal/claims`, Issue claim/unclaim, legacy claim list, legacy `claims` | generic Claims |
| `internal/runs`, `run record` | generic Runs |
| `internal/dev`, `dev start|gate|finish|block`, `gate_runs`, dev Events/readers | Attempt/Receipt/Evidence and Envelope completion |
| `issue complete`, `internal/issues/complete*`, writers of `task_completions` | Envelope completion |

`task_completions` remains a frozen no-writer table only because current
`concerns` and `ac_verifications` foreign keys are Phase 23-owned. It is not
read for authority or completion and is deleted in Phase 23. `agent_runs`
remains frozen with no Phase 21 writer for Phase 25 telemetry convergence and
is not a generic Run. Source/integration inventory must prove no removed CLI,
package, Event family, role lookup, or table writer survives. Phase 21 does not
reconstruct removed state through import.

Deferred to Phase 22: Attention, Stops, budget use/exhaustion, recovery,
automatic respawn, Activation recovery, successor Runs, export import, and
portable parity. Phase 23 owns judgment/review/decision surfaces and frozen
completion dependents. Phase 24 owns multi-Issue/revision/cancel/parallel
scope. Phase 25 owns attestation, telemetry, and full capability qualification.
`Routine` is an external kit risk-lane term and is never persisted, evaluated,
or exposed by fid.

## 16. Required Parity And Journey

The dedicated real-process integration journey follows the 13 steps in the
accepted Gate D3.2 record, tests the negative counterpart of each authority,
policy, idempotency, evidence, and completion branch, repeats after rebuild,
and proves byte-identical canonical reads for capabilities, Envelope,
evaluation, resume, and export. Maintenance reports sequence gaps, orphaned
bindings, Attempt-without-Receipt as non-passing, Receipt mismatch, missing
Artifact, invalid active Activation, frozen-table writes, and removed V1 state.
No implementation or QA task reads held-out scenarios or evaluator prompts.
