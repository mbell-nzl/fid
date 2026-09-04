# TECHNICAL REQUIREMENTS - fid

This document defines fid's aligned target architecture. The current codebase
contains pre-ADR-0011 packages and schema that are removed or replaced through
the implementation roadmap; their presence does not make them forward
contracts.

## 1. Stack

| Area | Decision |
|---|---|
| Language | Go, current project toolchain |
| CLI | `spf13/cobra` and `pflag` |
| Durable log | Append-only UTF-8 JSONL |
| Projection | SQLite through `modernc.org/sqlite` |
| Migration runner | `golang-migrate/migrate` until the pre-dogfood baseline is consolidated |
| Machine serialization | Standard-library JSON |
| Tests | Standard `testing` plus existing `testify` usage |
| Distribution | One local static-style CLI binary; no daemon or server |

The YAML runtime dependency becomes removable when the universal YAML output
and retired configuration inputs leave the code. TOON and universal Markdown
rendering add no target dependency. Purpose-built Human projections should use
the standard library unless a later dependency decision proves a need.

## 2. Architecture

Dependencies point inward:

```text
cmd/fid                       CLI and JSON transport adapters
   |
   +--> Issue tracker domain  Issues, workflows, ACs, labels, links, readiness
   +--> Work domain           Envelopes, Claims, Runs, Stops, Attention
   +--> Evidence domain       Registrations, Attempts, Receipts, Evidence
   +--> Review domain         Reviews, Findings, Concerns, Deviations, Outcomes
   +--> Projection domain     Context, completion, Human brief, telemetry views
   |
   +--> journal port          canonical append and ordered replay
   +--> projection port       SQLite transactions and queries
   +--> executor port         registered local process execution
   +--> repository port       Git/worktree identity and reconciliation facts
```

Rules:

- CLI code parses input, maps typed errors, and renders output; it owns no
  business policy.
- Domain operations return complete atomic effects and do not write SQLite
  directly.
- The journal adapter owns canonical event serialization and append locking.
- The SQLite adapter applies accepted events and owns no policy choices.
- The executor resolves immutable registrations; it never accepts replacement
  argv, cwd, environment, timeout, success, network, credential, write, or
  retention policy from the caller.
- Repository inspection reports facts; kit policy interprets risk and scope.
- The V2 adapter uses documented JSON operations only. It never reads `.fid/`,
  journal files, or projection tables.
- The in-memory fake implements the same semantic port only for tests and
  controlled fixtures.
- `cmd/fid/integration` is a test-only real-process adapter: its tagged
  journeys execute the built CLI in isolated workspaces. In-process command,
  SQLite, journal, and projection tests remain in their ordinary owning
  packages so the integration gate does not rerun the default suite.

Phase 21 keeps the recorded package boundaries: Work Envelope, Activation,
and generic Run logic are owned by `internal/workflow`; generic authority
replaces V1 role policy inside `internal/policy`; generic Claim behavior
replaces the V1 implementation inside `internal/claims`; and Gate
Registration, executor, Receipt, Artifact, and Evidence logic are owned by
`internal/evidence`. The slice does not create parallel `internal/envelopes`,
`internal/authority`, `internal/executor`, or `internal/artifacts` packages.
`internal/bootstrap` and the init/rebuild adapters own the role-free
`project.initialized` Event-v2 fixture. The generic evaluator in
`internal/validation` is called by both terminal-success Issue transition and
Envelope completion; neither adapter duplicates its rules.

## 3. Versioned Backend Capabilities

Capability versions evolve independently of the fid product version.

| Capability ID | Semantic boundary |
|---|---|
| `workflow.state.v1` | Issues selected by Work Envelopes, Runs, ordered Events, and current projections |
| `workflow.transitions.v1` | Atomic expected-revision transition evaluation and typed rejection |
| `work.claims.v1` | Claim acquire, renew, release, expiry, supersession, and conflict semantics |
| `kit.activation.v1` | Cross-resource prepare, component, commit, abort, finalize, and recovery state |
| `evidence.gate-definitions.v1` | Immutable Gate Definition Registration, authorization, binding, and resolution |
| `evidence.receipts.v1` | Execution Attempts and immutable terminal Receipts |
| `evidence.freshness.v1` | Obligation absence, record validity, Artifact availability/policy, and target freshness |
| `reviews.coverage.v1` | Immutable Review coverage, verdict, evidence, and zero-or-more Findings |
| `human.outcomes.v1` | Provenanced Human checks, Decisions, approvals, and risk acceptance |
| `attention.v1` | Typed deduplicated Attention lifecycle |
| `telemetry.trajectory.v1` | Run-bound observed economics, retries, failures, and interventions |
| `export.portable.v1` | Versioned backend-neutral validation, export, and import semantics |

Discovery returns product/version, supported backend-port range, each
capability/version/schema range, binary availability, limits, and accepted
trust sources. Missing or incompatible capability returns typed unavailability;
it returns a Stop only after the Stop capability exists.

### 3.1 Phase 21 qualification boundary

Phase 21 fixes the complete availability matrix in its slice contract. An
operation used by the one-Issue walking skeleton does not make its semantic
capability available when recovery, parallelism, import, or another required
behavior is deferred. Availability is derived from binary and migration state,
cannot be caller-declared, and has no partial or experimental value.

Kit-mode Actor identity, capabilities, producer, and provenance resolve only
from the active Kit Activation authority profile. Callers pass an Actor
reference and cannot augment it. Gate Registration similarly resolves a
finalized activated Gate Definition before any Envelope may reference it.

The Phase 21 executor is a standard-library local process adapter, not an OS
containment boundary. It enforces structured argv without a shell, exact
contained cwd, a closed environment allowlist, timeout/output ceilings,
credential-name policy, redaction, and contained Artifact paths. It accepts
only `sandbox:none`, `network:allowed`, and `write:allowed`; registration
rejects stronger unsupported values. There is no automatic process/agent
respawn or Run continuation.

The complete Phase 21 wire surface, including closed input/result/error/live
schemas and process exits, is
`../docs/phases/phase-21/public-interface.md`. The live legacy caller and
reader/writer disposition is
`../docs/phases/phase-21/v1-convergence-inventory.md`.

## 4. Contract Rules

Phase 22 has a proposed recovery amendment under DEC-180, pending Human
acceptance: `../docs/phases/phase-22/contracts.md`, `public-interface.md`, and
`schema.md` in the same directory. The proposal keeps Run/Stop/Attention/budget/
Activation rules in `internal/workflow`, leases in `internal/claims`, observed
execution in `internal/evidence`, and shared guards in `internal/validation`.
Read-only repository observation extends the existing `internal/project`
adapter. It adds no package or dependency. The live package reconciliation and
gate reach are recorded in `../docs/phases/phase-22/execution-scope.md`.

DEC-184/185 refine that proposal: public-interface section 2.1 enumerates the
23 activated action names and bootstrap exclusions; contracts sections 3.1/3.2
define canonical dirty/reconciliation digests. Caller-supplied targets use the
same digest grammar; replay/import never translate historical digests into
fresh proof. DEC-187's `kit-contract-mapping.md` records the kit coordinator's
recovery-export binding and derived-terminal-state adapter obligations, with
production conformance deferred to Phase 25 qualification.

The proposed final matrix enables Attention, activation recovery, and portable
export/import only after their complete conformance evidence. Claims recovery
does not by itself qualify `work.claims.v1`: parallel scope semantics remain
Phase 24, so that capability remains unavailable. Existing delivered behavior
is authoritative until the proposal is accepted and implemented.

- JSON objects include a schema or contract version at the operation envelope
  or record level.
- IDs are opaque stable strings. Consumers do not parse prefixes for behavior.
- Enums serialize as lowercase kebab-case.
- Times are UTC RFC 3339 with sub-second precision where available.
- Unknown required enum values or unsupported major versions reject.
- Mutations accept an idempotency key and expected revision where concurrency
  or retry is possible.
- Successful mutation responses name the canonical entity and Event identity.
- Errors contain stable `code`, `category`, `retryable`, Human message, subject,
  and optional structured details. Human message text is never control flow.
- Full logs, secrets, credentials, and raw environment snapshots are forbidden
  in normal JSON responses and journal payloads.
- Phase 21 mutation retries use a durable operation/project-scoped idempotency
  key and canonical request digest; same key/digest replays the stored result,
  while a changed digest conflicts without an Event.

## 5. Semantic Operation Groups

Exact CLI spelling is accepted with the owning implementation phase; these
semantic operations bind every transport.

### 5.1 Work and projections

- create/get/revise/activate/complete Work Envelope;
- list actionable work with typed non-actionable reasons;
- get context, changed-since, resume, completion/evidence, Attention, and Human
  brief projections;
- evaluate and atomically apply expected-revision transitions.

Work operations first validate active Kit Activation. Activation disagreement
blocks normal workflow and leaves only discovery, inspection, and recovery
operations available.

### 5.2 Claims, Runs, Stops, and recovery

- acquire, renew, release, expire, or authorizedly supersede a Claim;
- start/terminalize a Run and record incremental observed telemetry;
- record/resolve a typed Stop and raise/resolve deduplicated Attention;
- recover interruption by terminalizing the prior Run, reconciling repository
  state, disposing the prior Claim, and creating the linked successor Claim/Run
  atomically.

One Issue permits one active mutation Claim. Parallel mutation requires
distinct Issues, worktrees, and provably disjoint scopes. Lease mismatch,
unknown liveness, unsafe repository state, or active prior-revision mutation
rejects without partial effects.

### 5.3 Registered execution and evidence

- register and resolve immutable Gate Definitions;
- append an Execution Attempt before spawn;
- execute locally with live output and produce one terminal Receipt;
- ingest a qualifying external attested Receipt;
- link Evidence Records and evaluate applicable Obligations.

Receipt acceptance verifies the exact registration and binding resolved at
start. Abandoned Attempts remain visible and never pass. Evidence evaluation
returns absence, validity/freshness, and Artifact availability separately.

### 5.4 Review and judgment

- record a Review Record and Findings;
- record addressed assertions and later independent Finding resolution;
- record Concerns, Deviations, Decision References, and Human Outcomes;
- raise, list, and resolve Attention through the qualifying outcome/fact.

Project prose remains referenced by path plus repository identity/digest. fid
stores no editable Decision, QA Plan, Validation Contract, handoff, or
completion-report body.

### 5.5 Recovery and activation

- validate, create, and import Portable Recovery Exports;
- prepare, record component state, commit, abort, finalize, inspect, and recover
  Kit Activation.

Activation has one durable point of no return. Before it, recovery restores the
old active set; after an irreversible backend commit, recovery completes the
new set forward.

## 6. Journal And Projection

Each canonical Event contains:

- event and schema identity;
- aggregate type, aggregate ID, and monotonic aggregate sequence;
- occurred/recorded time, Actor, producer, correlation, causation, and
  idempotency key;
- repository identity when observed or affected;
- one typed payload sufficient to apply and rebuild the durable fact; and
- integrity metadata required by export and maintenance.

The writer holds a project-scoped exclusive lock across catch-up, validation,
append, fsync, projection transaction, and projection metadata update. A failed
validation or append has no projection effect. A projection failure after a
durable append leaves a detectable catch-up condition; the next operation
replays before accepting a new write.

Projection tables may normalize current query state, but canonical history is
never inferred from row update timestamps. Rebuild and live application use
the same event appliers and must produce semantic parity.

## 7. Concurrency And Idempotency

- Aggregate sequence and expected revision detect stale writers.
- Idempotency keys are scoped to operation/project and bind a request digest;
  replaying the same key and digest returns the original result, while a
  different digest conflicts.
- One logical mutation appends one atomic Event transaction or no Event.
- Claim and active-revision uniqueness is enforced by domain validation and
  projection constraints as corruption backstops.
- A Receipt has a unique Execution Attempt; a Review is immutable; an Evidence
  Record never overwrites prior proof.

## 8. Security And Trust

- The CLI is not a remote security boundary. Kit-mode protected authority is
  established by active-version provenance and explicit grants, not a role
  string.
- Standalone Issue inventory commands cannot create protected Human Outcomes or
  trusted Receipts.
- Gate registration binds project, overlay, repository/worktree, approved
  scope, definition digest, and authorization provenance.
- Execution uses structured argv without a shell by default, explicit working
  directory, bounded timeout, environment allowlist, output limits, and the
  registered sandbox/network/write/credential policy.
- Control characters and secrets are removed or escaped before persistence and
  Human rendering. Redaction happens before excerpts and hashes are exposed.
- Artifact paths are content-addressed or stable opaque references and cannot
  escape configured roots.
- Direct SQL remains an unstable maintainer diagnostic and is never a machine
  integration or public schema contract.

## 9. Performance Targets

| Operation | Target |
|---|---|
| Common Issue/actionable query at 10k Issues | p95 under 500 ms on reference hardware |
| Single journal append plus projection | p95 under 100 ms excluding gate execution |
| 100k-Event full rebuild | under 30 s on reference hardware |
| Capability discovery | under 100 ms |
| Receipt/evidence evaluation for one Envelope | under 500 ms at representative Obligation volume |

Targets include deterministic benchmark fixtures and may not be weakened to
accommodate a regression. A performance tier remains evidence-gated under
ADR-0011.

For Phase 21, the 100k-Event reference instance is Linux/amd64 on an AMD Ryzen
7 5700G with six exposed logical CPUs and 8,133,236 kB recorded memory, using
Go 1.26.2. Evidence also records the exact kernel, `GOOS`/`GOARCH`, toolchain,
logical CPU count, CPU model, and memory visible to the process. Acceptance is
three consecutive fresh one-iteration benchmark processes on that instance;
every sample must self-assert parity and complete below 30 seconds. Best-of or
median tolerance is not accepted. The general benchmark gate retains the same
per-sample failure assertion. An unavailable reference instance is `BLOCKED`,
and a slow sample requires causal optimization rather than threshold change.

## 10. Pre-Dogfood Convergence Rules

- Replacement contracts and tests land before displaced V1 code is deleted in
  each buildable slice.
- No compatibility adapter, old-journal upgrader, reported-gate importer, or
  universal format shim is added.
- The current migration chain is development history, not a product promise.
  The final convergence slice creates one clean supported baseline and proves
  init, rebuild, export/import, and conformance from that baseline.
- No implementation task may rely on a logical entity until its accepted
  physical DDL and Event payload are recorded in `03_TRD_SCHEMA.md` for that
  slice.

*Last updated: 2026-08-23 (Phase 21 DEC-156/157/158 planning amendments)*
