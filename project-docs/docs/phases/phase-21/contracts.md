# Phase 21 Slice Contract - One-Issue Walking Skeleton

> Status: **ACCEPTED FOR EXECUTION - DEC-159**
> Owner: Tech Lead
> Operating mode: Heavy
> Decision authority: DEC-152..159, ADR-0011, ADR-0012
> Physical schema: `../../../specs/03_TRD_SCHEMA.md` section 8.2

This is the exact Phase 21 implementation boundary. The QA Lead supplied a
Ready Heavy plan at `e8a636e`, and DEC-159 records Human execution approval.
Normative JSON examples use canonical compact JSON: UTF-8, sorted object keys, preserved
array order, no insignificant whitespace, and a trailing newline on CLI output.
`public-interface.md` owns every CLI argument, closed input, result, error,
exit, Human rendering, and live-frame payload. `v1-convergence-inventory.md`
owns every live caller/reader/writer and removal disposition.

## 1. Versions And Shared Values

- Response Envelope: retained `fid-cli.v1` with Phase 21 leaves and reason
  codes added by this contract.
- Event schema: `fid.event.v2`; v1 development journals are not upgraded.
- Envelope contract: `fid.envelope.v1`.
- Gate Definition and Registration: `fid.gate-definition.v1` and
  `fid.gate-registration.v1`.
- Export archive manifest: `fid.export-manifest.v1`. This is a file format,
  not the unavailable `export.portable.v1` semantic capability.

An activated authority-profile entry is exactly:

```json
{"actor_id":"actor-builder","actor_type":"agent","capabilities":["artifact.verify","claim.acquire","claim.release","claim.renew","envelope.activate","envelope.complete","envelope.create","evidence.link","gate.register","gate.run","issue.transition","run.finish","run.start"],"producer":"kit","provenance_ref":"activation:ACT-1/profile:builder","session_ref":null}
```

For kit-mode mutation, callers submit only `actor_id`; fid resolves all other
fields from the active Activation. Event/result Actor views omit capabilities
and contain exactly Actor ID/type, producer, provenance reference, and nullable
session reference. Unknown Actor fields reject. Standalone inventory mutations
accept only the typed `--actor <type>:<id>` grammar in `public-interface.md`;
the derived project-local Actor is not trusted kit authority and cannot
authorize an Envelope, Claim, Run, Receipt, Evidence, guarded transition, or
completion operation.

A repository identity is exactly:

```json
{"base_oid":"<40hex-or-null>","capture_source":"git","dirty":false,"dirty_digest":null,"head_oid":"<40hex-or-null>","project_id":"<project-id>","tree_oid":"<40hex-or-null>","worktree_id":"<stable-id>"}
```

A source reference is exactly:

```json
{"digest":"sha256:0000000000000000000000000000000000000000000000000000000000000000","display_label":"Gate policy","external_id":null,"path":"project-docs/gates/unit.json","repository":{"base_oid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","capture_source":"git","dirty":false,"dirty_digest":null,"head_oid":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","project_id":"project-1","tree_oid":"cccccccccccccccccccccccccccccccccccccccc","worktree_id":"worktree-1"},"section":null,"source_type":"project-file"}
```

An immutable Phase 21 budget is exactly:

```json
{"consumed":null,"exhaustion_action":"unavailable","maximum":3,"name":"attempts","remaining":null,"scope":"envelope","unit":"attempts","warning_threshold":null}
```

`consumed` and `remaining` MUST be JSON null. Phase 21 never updates counters,
warns, exhausts, Stops, raises Attention, or changes execution from a budget.

The complete generic Obligation value is exactly:

```json
{"accepted_evidence_types":["receipt"],"accepted_trust_sources":["fid-executed"],"applicability":"required","artifact_policy":"required-available","derived_status":"unsatisfied","freshness_target":"repository-tree","guarded_actions":["envelope.complete","issue.transition-terminal-success"],"id":"OBL-1","independence":"fid-executed","required_lane":null,"source":{"digest":"sha256:0000000000000000000000000000000000000000000000000000000000000000","display_label":"Gate policy","external_id":null,"path":"project-docs/gates/unit.json","repository":{"base_oid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","capture_source":"git","dirty":false,"dirty_digest":null,"head_oid":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","project_id":"project-1","tree_oid":"cccccccccccccccccccccccccccccccccccccccc","worktree_id":"worktree-1"},"section":null,"source_type":"project-file"},"subject":{"id":"ISS-1","type":"issue"},"type":"gate-evidence"}
```

Only `gate-evidence` is accepted. Unknown types reject; `derived_status` is
computed and rejected on mutation input. The value shape remains generic so a
later type does not require an Envelope schema change.

## 2. Capability Availability Contract

Availability is binary and migration-derived. An operation existing for this
walking skeleton does not make a capability available when required semantics
are deferred. Discovery returns rows sorted by `capability_id`:

| Capability ID | Available | Exact closed `limits` object | Trust source |
|---|---:|---|---|
| `attention.v1` | false | `{"state":false}` | `phase-22-unavailable` |
| `evidence.freshness.v1` | true | `{"branch_labels":false,"target":"repository-tree"}` | `fid-git-observation` |
| `evidence.gate-definitions.v1` | true | `{"executor_policy":"standard-library","network":"allowed","sandbox":"none","write":"allowed"}` | `active-kit-activation` |
| `evidence.receipts.v1` | true | `{"attestation":false,"producer":"fid","terminal_receipts":true}` | `fid-executed` |
| `export.portable.v1` | false | `{"create":true,"import":false,"validate":true}` | `phase-22-unavailable` |
| `human.outcomes.v1` | false | `{"record":false}` | `phase-23-unavailable` |
| `kit.activation.v1` | false | `{"abort_before_commit":true,"automatic_respawn":false,"finalize":true,"recover":false}` | `phase-22-unavailable` |
| `reviews.coverage.v1` | false | `{"record":false}` | `phase-23-unavailable` |
| `telemetry.trajectory.v1` | false | `{"ingest":false}` | `phase-25-unavailable` |
| `work.claims.v1` | false | `{"max_active_per_issue":1,"parallelism":false,"recover":false}` | `phase-22-unavailable` |
| `workflow.state.v1` | true | `{"max_envelope_revisions":1,"max_primary_issues":1,"related_issues":false}` | `fid-journal` |
| `workflow.transitions.v1` | true | `{"execution_bound":true,"terminal_evidence_guard":true}` | `fid-journal` |

Each discovery row is exactly:

```json
{"available":false,"capability_id":"export.portable.v1","limits":{"create":true,"import":false,"validate":true},"schema_max":1,"schema_min":1,"trust_source":"phase-22-unavailable"}
```

No omitted row, partial/experimental state, caller override, or optimistic
default is permitted. Rebuild must reproduce the same matrix for the same
binary and migration state.

The required unavailable `attention.v1` discovery row is capability metadata,
not Attention state or an Attention effect. In Phase 21 it is the sole allowed
production `Attention` match: its exact false row must exist, while no Attention
table, Event family, Effect, command leaf, or other state/operation surface may
exist. Missing or changing the false row fails capability conformance; any
additional production match fails the Phase 21 deferral boundary.

## 3. Idempotency

Every Phase 21 mutation, plus retained retryable Issue create and transition,
requires an `--idempotency-key` of 1..128 printable ASCII characters. The
scope is `<operation>:<project-id>`. fid hashes the canonical mutation request
after Actor resolution but before side effects.

- first key/digest: execute once and durably store response entity ID, Event ID,
  aggregate sequence, and canonical result JSON in the same transaction;
- same key/same digest: return the stored entity/Event result with only the
  contracted `replayed:true` metadata difference and append no Event;
- same key/different digest: reject `idempotency-conflict` and append no Event;
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
source. Abort after commit rejects `activation-recovery-unavailable`.

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
{"artifact_bytes":10485760,"artifact_paths":["coverage.out"],"argv":["go","test","./internal/issues"],"credential_names":[],"cwd":".","definition_id":"gate-unit","definition_version":1,"environment":[{"name":"GOCACHE","source":"literal","value":".gocache"},{"name":"HOME","source":"inherit","value":null}],"excerpt_bytes":4096,"network":"allowed","output_bytes":1048576,"redact_literals":[],"sandbox":"none","shell":false,"success_exit_codes":[0],"timeout_ms":120000,"write":"allowed"}
```

Environment entries use only `literal`, `inherit`, or `credential`; credential
entries store a credential name and null value. Names are unique and sorted.
The inherited source is a closed allowlist, not the ambient environment.
Registration validates non-empty structured argv; `shell:false`; contained,
symlink-resolved cwd and Artifact paths; positive configured timeout/output/
excerpt ceilings; non-empty unique success codes; and exact `sandbox:none`,
`network:allowed`, `write:allowed`. Stronger or different values reject
`executor-policy-unsupported`. fid claims no OS, filesystem, network,
privilege, namespace, or platform containment.

A Registration stores immutable definition digest, binding digest, active
Activation/provenance, repository/project/scope binding, and resolved policy.
Same Definition and binding return the existing Registration; changed content
or binding creates a new immutable version. Envelope creation references only
the Registration ID/version. Registration authorization uses the finalized
Actor profile's exact `gate.register` capability because no Envelope exists
yet; caller attribution and later Envelope Grants cannot authorize it.

## 6. Envelope, Authority, And Completion

Phase 21 accepts exactly one Envelope revision, one existing primary Issue,
zero related Issues, one `gate-evidence` Obligation referencing a prior
Registration, zero or more immutable budgets, and sorted Authority Grants.
Revision, cancel, protected Decision, Stop, and Attention inputs reject as
unavailable. Activation requires an active trust source and no active mutation
Claim on the Issue.

An Authority Grant is exactly:

```json
{"actions":["artifact.verify","claim.acquire","claim.release","claim.renew","envelope.activate","envelope.complete","evidence.link","gate.run","issue.transition","run.finish","run.start"],"actor_selector":{"actor_id":"actor-builder","required_capabilities":["artifact.verify","claim.acquire","claim.release","claim.renew","envelope.activate","envelope.complete","evidence.link","gate.run","issue.transition","run.finish","run.start"]},"scope":{"issue_id":"ISS-1"}}
```

Envelope creation first requires the finalized profile's exact
`envelope.create` capability. It then records the Authority Grants used by
later operations. Post-creation authorization resolves the Actor from the
active Activation and requires an exact action, capability, Issue, Envelope,
and revision Grant match. Missing is deny.
Caller attribution, the deleted role registry, a Human-looking string, or an
inactive Activation confers no authority.

The same Obligation evaluator guards two separate actions. First, an Envelope-bound
transition into the workflow terminal-success category succeeds only after
every applicable gate-evidence Obligation is `satisfied`. Later, completion
re-reads the active revision and calls that evaluator again; it additionally
requires the selected Issue already to be terminal-success. Either failure
returns all section 10 reasons and changes no state. Completion does not write
or consult legacy `task_completions`.

## 7. Claim, Run, And Execution-Bound Transition

One mutating Claim per Issue is supported. Acquire requires active Envelope,
resolved Actor authority, expected Envelope revision, worktree identity,
declared scope, positive TTL, and idempotency. Renew requires the current Claim
version, same Actor/worktree, and unexpired Claim. Release is idempotent for
the same Claim version/reason. Conflict, expiry, stale version, or wrong Actor does
not authorize mutation.

A Run starts from an active Claim and binds Issue, Envelope/revision, Actor,
Claim/version, repository identity, and immutable budget contracts. Its budget
JSON always exposes null consumption/remaining. Status is `active`,
`succeeded`, `failed`, or `abandoned`; only explicit finish terminalizes it.
There is no recovery, predecessor, automatic continuation, or respawn.

Execution-bound nonterminal Issue transition requires active Run and Claim, active
Envelope revision, resolved Actor grant, exact expected Issue revision, valid
workflow edge, and matching worktree/repository observation. A
terminal-success transition adds the shared Obligation evaluation after those
predicates. Checks return every independent reason using section 10 ordering.
No Event or revision advance is allowed on failure.

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

Control frames report only typed metadata: `output-truncated`,
`artifact-rejected`, or `redaction-applied`. When the combined sanitized output
limit is reached, fid stops retaining/emitting process output, emits exactly
one `output-truncated`, continues draining both pipes, and still observes
termination. Redaction applies before live output, excerpt, hash, Receipt,
Artifact, error detail, or Human rendering. Credential values never enter
journal, projection, arguments in errors, or exported data.

Exactly one terminal frame is emitted after observed termination and after the
terminal Receipt commits, or after a handled frame-only `fid-error`. Its
`data` contains the complete canonical `fid-cli.v1` Response Envelope; no
additional JSON line follows it. Terminal process outcomes are `succeeded`,
`failed`, `timed-out`, `signaled`, or `spawn-failed`. A handled fid
infrastructure failure after streaming begins uses the distinct frame-only
`fid-error` outcome, null Receipt, `internal-error` response, and exit 7; it
never satisfies Evidence or commits idempotency.
`output-limit-observed` is additional control metadata and does not replace
the observed process result. An
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
produce `artifact-rejected` and cannot satisfy `required-available`.

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

1. `authority-unavailable`
2. `envelope-inactive`
3. `claim-missing`
4. `claim-expired`
5. `run-inactive`
6. `issue-revision-conflict`
7. `workflow-transition-invalid`
8. `obligation-evidence-absent`
9. `receipt-unsuccessful`
10. `evidence-invalid`
11. `evidence-stale`
12. `artifact-unavailable`
13. `issue-not-terminal-success`

Each reason is exactly `{"code":"...","detail":"...","subject_id":"...","subject_type":"..."}`.
Applicable gate evidence is satisfied only when one exact Receipt succeeded,
Evidence is valid and repository-tree fresh, and every required Artifact is
available. Attention is never raised.

The shared Obligation evaluator returns only reasons 8..12 and is byte-identical
when called by pure evaluation, terminal-success transition, or completion.
The terminal transition may additionally return action predicates 1..7.
Completion may additionally return 1, 2, and 13; it does not require an active
Claim or Run, so the accepted journey can finish/release before completion.

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

## 12. Public Interface Authority

`public-interface.md` fixes every accepted and fail-closed leaf, exact
arguments/flags, closed input schema, success result, collection ordering,
Event/idempotency/replay representation, Human rendering, failure detail,
stdout/stderr rule, process exit class, and `fid.live.v1` payload union. There
is no `result:{}` placeholder or implied complete entity outside that file.

## 13. Closed Public Values

All contract-controlled enum and reason values are lowercase kebab-case. JSON
field names remain snake_case. Opaque user-authored workflow state identifiers,
entity IDs, paths, capability IDs, and digests are identifiers rather than
closed enum members. The complete reason/code/category/retryability/subject/
details/exit table is `public-interface.md` section 7.

## 14. Event Schema V2 And Fixtures

Every new Event and every retained mutation converted in Phase 21 uses exactly:

```json
{"actor":{"actor_id":"actor-builder","actor_type":"agent","producer":"kit","provenance_ref":"activation:ACT-1/profile:builder","session_ref":null},"aggregate":{"id":"ISS-1","sequence":2,"type":"issue"},"causation_id":null,"command":"issue.transition","correlation_id":"RUN-1","effects":[{"op":"update","pk":{"id":"ISS-1"},"preconditions":{"revision":1,"status":"todo","workflow_id":"task-default"},"prior":{"revision":1,"status":"todo","updated_at":"2026-08-21T23:59:00Z"},"table":"issues","values":{"revision":2,"status":"in_progress","updated_at":"2026-08-22T00:00:00Z"}},{"op":"activity","pk":{"id":"ACT-2"},"table":"activity_log","values":{"actor":"actor-builder","details":{"from":"todo","to":"in_progress"},"event":"status-changed","journal_event_id":"EVT-20260822-000002","occurred_at":"2026-08-22T00:00:00Z","ref_id":"ISS-1","ref_type":"issue"}},{"op":"insert","pk":{"idempotency_key":"key-1","operation_scope":"issue.transition:project-1"},"table":"idempotency_results","values":{"aggregate_sequence":2,"created_at":"2026-08-22T00:00:00Z","idempotency_key":"key-1","operation_scope":"issue.transition:project-1","request_digest":"sha256:1111111111111111111111111111111111111111111111111111111111111111","result_entity_id":"ISS-1","result_event_id":"EVT-20260822-000002","result_json":"{\"contract_version\":\"fid-cli.v1\",\"event\":{\"id\":\"EVT-20260822-000002\",\"occurred_at\":\"2026-08-22T00:00:00Z\"},\"ok\":true,\"operation\":\"issue.transition\",\"result\":{\"from_state\":\"todo\",\"idempotency_key\":\"key-1\",\"issue\":{\"id\":\"ISS-1\",\"revision\":2},\"replayed\":false,\"to_state\":\"in_progress\"}}"}}],"event_id":"EVT-20260822-000002","event_schema":2,"event_type":"issue.transitioned","idempotency":{"key":"key-1","request_digest":"sha256:1111111111111111111111111111111111111111111111111111111111111111","scope":"issue.transition:project-1"},"occurred_at":"2026-08-22T00:00:00Z","repository":{"base_oid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","capture_source":"git","dirty":false,"dirty_digest":null,"head_oid":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","project_id":"project-1","tree_oid":"cccccccccccccccccccccccccccccccccccccccc","worktree_id":"worktree-1"}}
```

Exact rules: all top-level keys are required; nullable values are explicit;
`effects` is ordered and nonempty except for the exact bootstrap fixture below;
aggregate sequence begins at 1 and has no gap; Actor capabilities are not
copied into Events; repository is null only for a non-repository bootstrap
fact; unknown keys reject; canonical fixtures replace affected v1 fixtures.
Kit mutations resolve the complete Actor from the finalized Activation.
Retained standalone mutations derive the Actor only from the typed
`--actor <type>:<id>` attribution grammar in `public-interface.md`: the type
and ID are copied, producer is `fid-standalone`, provenance is
`project-local-attribution`, and session is null. That Actor is descriptive
Event provenance only and is never eligible for an Authority Grant.

The final-binary bootstrap fixture is exactly:

```json
{"actor":{"actor_id":"fid-bootstrap","actor_type":"service","producer":"fid","provenance_ref":"fid:init","session_ref":null},"aggregate":{"id":"project-1","sequence":1,"type":"project"},"causation_id":null,"command":"project.init","correlation_id":null,"effects":[],"event_id":"EVT-20260822-000001","event_schema":2,"event_type":"project.initialized","idempotency":null,"occurred_at":"2026-08-22T00:00:00Z","repository":null}
```

The project ID comes from strict project config written before append.
Initialization is naturally exclusive: a second invocation returns
`project-already-initialized`, never replay. Migrations seed reference
workflows; the Event inserts no role, workflow, or policy row. Live apply only
records journal projection metadata. Rebuild applies through migration 0014,
reads the same project config, replays the Event, and produces the same
role-free projection. No bootstrap, Actor lookup, export, or down-migration
path reconstructs removed role facts.

Tasks 2-19 may retain one nonpersistent V1 compatibility bridge solely for
callers explicitly classified in `v1-convergence-inventory.md`. When the
journal contains the exact role-free Event-v2 bootstrap and no legacy role
Effect, `internal/policy` may recognize the closed V1 vocabulary
`developer`, `qa_reviewer`, `qa_lead`, `tl_lite`, `tl`, `workflow_runner`, and
`human` for the existing V1 permissions matrix. If any legacy role Effects
exist, those facts remain authoritative and the fallback does not supplement
them; if neither the exact bootstrap nor role Effects exist, role lookup stays
empty and fail-closed. Unknown names always fail. The bridge writes no Event,
Effect, row, config, or export field; is not an Actor, Activation, capability,
Grant, Evidence, Human Outcome, or provenance source; and cannot authorize a
new generic path, guarded transition, or Envelope completion. Task 20 removes
it with the V1 policy implementation.

The accepted Phase 21 families are:

| Aggregate | Event types | Projection effects |
|---|---|---|
| project | `project.initialized` | no domain Effect; journal projection metadata only |
| activation | `activation.prepared`, `.component-staged`, `.committed`, `.finalized`, `.aborted` | activation/component/profile/definition insert or state update |
| gate-registration | `gate.registered` | registration insert |
| envelope | `envelope.created`, `.activated`, `.completed` | envelope/revision/source/selection/grant/obligation/budget insert or status update |
| claim | `claim.acquired`, `.renewed`, `.released` | claim insert or version/time/status update |
| run | `run.started`, `.finished` | run/budget-binding insert or terminal update |
| attempt | `gate.execution-started` | Attempt insert |
| receipt | `gate.receipt-recorded` | Receipt and Artifact inserts |
| evidence | `evidence.linked`, `artifact.availability-observed` | Evidence/link insert or optional Artifact availability update |
| issue | `issue.created`, `.updated`, `.transitioned`, `.ac-added`, `.label-added`, `.label-removed` | Issue/AC/Issue-label/activity/revision effects |
| label | `label.added` | label insert/update and activity effects |
| link | `link.added`, `.removed` | link/dependency effects plus affected Issue revisions |
| workflow | `workflow.defined`, `.assigned` | structural workflow and assignment effects only |
| requirement | `requirement.synced` | canonical requirement-reference effects |

For fixture generation, each row below is normative. `insert(T)` contains the
complete DDL row, `update(T:{keys})` contains exactly those changed values plus
their complete `prior` and key/version `preconditions`, and repeated rows sort
by their declared primary key. Every idempotent mutation appends
`insert(idempotency_results)` as its final Effect; its stored `result_json` is
the canonical non-replayed Response Envelope for this Event.

| Event type | Aggregate | Exact ordered domain Effects before idempotency |
|---|---|---|
| `activation.prepared` | Activation | `insert(kit_activations)` |
| `activation.component-staged` authority profile | Activation | `insert(activation_components)`, then `insert(activation_actor_profiles)` by Actor ID |
| `activation.component-staged` gate definitions | Activation | `insert(activation_components)`, then `insert(activation_gate_definitions)` by Definition ID/version |
| `activation.committed` | Activation | `update(kit_activations:{status,point_of_no_return,updated_at})`, then `update(activation_components:{status})` by component key |
| `activation.finalized` | Activation | `update(kit_activations:{status,updated_at})`, then `update(activation_components:{status})` by component key |
| `activation.aborted` | Activation | `update(kit_activations:{status,updated_at})`, then `update(activation_components:{status})` by component key |
| `gate.registered` | Gate Registration ID/version | `insert(gate_registrations)` |
| `envelope.created` | Envelope | `insert(work_envelopes)`, `insert(envelope_revisions)`, sources by ordinal, primary selection, grants by ordinal, Obligation, then budgets by name |
| `envelope.activated` | Envelope | `update(work_envelopes:{status})`, `update(envelope_revisions:{status,approved_actor_id,activated_at})` |
| `envelope.completed` | Envelope | `update(work_envelopes:{status})`, `update(envelope_revisions:{status,completed_at})` |
| `claim.acquired` | Claim | `insert(work_claims)`; after 0014 the table name in new fixtures is `claims` |
| `claim.renewed` | Claim | `update(work_claims:{version,expires_at})` with prior Claim version |
| `claim.released` | Claim | `update(work_claims:{version,status,released_at,release_reason})` with prior Claim version |
| `run.started` | Run | `insert(runs)`, then `insert(run_budget_bindings)` by budget name |
| `run.finished` | Run | `update(runs:{status,terminal_repository_json,ended_at})` |
| `gate.execution-started` | Attempt | `insert(execution_attempts)` |
| `gate.receipt-recorded` | Receipt | new accepted `insert(referenced_artifacts)` by digest, then `insert(execution_receipts)`, then `insert(receipt_artifacts)` by Artifact ID; an already stored digest omits its Artifact insert; rejected candidates have no Artifact/link Effect |
| `evidence.linked` | Evidence | `insert(evidence_records)`, then `insert(evidence_artifact_links)` by Artifact ID |
| `artifact.availability-observed` | Artifact | optional `update(referenced_artifacts:{availability})` when state changed; unchanged observation has only the required idempotency Effect |
| `issue.created` | Issue | `insert(issues)`, labels by name, optional scheduling, then `insert(activity_log)` |
| `issue.updated` | Issue | requested Issue/label changes in field order, revision update, then `insert(activity_log)` |
| `issue.transitioned` | Issue | exact Issue status/revision update, then `insert(activity_log)` |
| `issue.ac-added` | Issue | `insert(acceptance_criteria)`, Issue revision update, then `insert(activity_log)` |
| `issue.label-added`, `issue.label-removed` | Issue | exact Issue-label insert/delete, Issue revision update, then `insert(activity_log)` |
| `label.added` | Label | `insert(labels)`, then `insert(activity_log)` |
| `link.added`, `link.removed` | Link | exact typed link/dependency insert/delete, affected Issue revision updates by ID, then `insert(activity_log)` |
| `workflow.defined` | Workflow | workflow/version/state/transition inserts in declared ordinal order; no role/judgment requirement Effect |
| `workflow.assigned` | Workflow | exact kind/project assignment upsert |
| `requirement.synced` | Requirement source | requirement reference upserts by requirement ID |

The Receipt fixture's order is also the SQLite apply order, so every foreign
key exists before its dependent link. A golden fixture asserts that canonical
order; the projector cannot silently reorder public JSON.

Receipt recording is one Event containing the terminal Receipt and every
accepted Artifact effect. Evidence linking is a separate mutation/Event. An
Attempt is always durable before spawn. No Phase 21 Event family exists for
Attention, Stop, budget consumption, recovery, respawn, import, Review,
Finding, Human Outcome, or telemetry.

## 15. Removal, Frozen State, And Deferral Inventory

`v1-convergence-inventory.md` exhaustively assigns every live role-policy
caller, mutation flag, workflow predicate, Event family, table reader/writer,
bootstrap path, command, and stats grouping. It also distinguishes temporary
Group 1 replacement state from task 20 final removal and task 21 final-schema
parity. An unclassified live match is a contract failure.

`task_completions`, `agent_runs`, `ac_verifications`, and `concerns` remain
frozen only for Phase 23/25 history/foreign-key ownership. They have no final
Phase 21 writer and cannot authorize kit operations, satisfy Evidence/Human
Outcomes, permit a guarded transition, or satisfy Envelope completion.

Phase 21 creates no unrecorded package boundary. Work Envelope, Activation,
and Run logic live in `internal/workflow`; generic authority remains in
`internal/policy`; generic Claims replace `internal/claims` in place; and Gate
Registration, executor, Receipt, Artifact, and Evidence logic live in
`internal/evidence`. Adapters remain in `internal/cli` and `cmd/fid`.

Deferred to Phase 22: Attention, Stops, budget use/exhaustion, recovery,
automatic respawn, Activation recovery, successor Runs, export import, and
portable parity. Phase 23 owns judgment/review/decision surfaces and frozen
completion dependents. Phase 24 owns multi-Issue/revision/cancel/parallel
scope. Phase 25 owns attestation, telemetry, and full capability qualification.
`Routine` is an external kit risk-lane term and is never persisted, evaluated,
or exposed by fid.

## 16. Required Parity And Journey

The dedicated real-process integration journey mutates once. It performs the
nonterminal transition, executes the Gate, links/evaluates its Receipt, then
uses the shared evaluator for terminal-success transition and later Envelope
completion. It tests the negative counterpart of each authority, policy,
idempotency, evidence, transition, and completion branch. It then deletes the
projection and rebuilds from the durable facts. After rebuild it repeats only
capability/Envelope/Receipt/Evidence reads, evaluation, completion-state read,
resume, export create/validate, and negative conformance checks. It never
creates a second one-revision Envelope. Maintenance reports sequence gaps, orphaned
bindings, Attempt-without-Receipt as non-passing, Receipt mismatch, missing
Artifact, invalid active Activation, frozen-table writes, and removed V1 state.
No implementation or QA task reads held-out scenarios or evaluator prompts.

The mutation prefix commits exactly 20 Events: `project.initialized`,
`issue.created`, `activation.prepared`, two `activation.component-staged`,
`activation.committed`, `activation.finalized`, `gate.registered`,
`envelope.created`, `envelope.activated`, `claim.acquired`, `run.started`, the
nonterminal `issue.transitioned`, `gate.execution-started`,
`gate.receipt-recorded`, `evidence.linked`, the terminal-success
`issue.transitioned`, `run.finished`, `claim.released`, and
`envelope.completed`. The complete journey contains exactly one
`envelope.created`. Post-rebuild commands append zero Events; a count other
than 20 before rebuild or any later increase is a conformance failure.
