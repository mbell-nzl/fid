# Phase 22 Slice Contract - Failure And Recovery

> Status: ACCEPTED by Human under DEC-188, 2026-09-05.
> Author: Tech Lead. Decisions: DEC-180/182/184..188; DEC-198 prepare/staging clarification.
> Planning base: `89d964a84dbdd297fdf472d6f78b001e2efac702`.

## 1. Boundary And Authority

One selected Issue and one immutable Envelope Revision remain the execution
boundary. The phase adds interruption, intentional abandonment, linked recovery,
durable resource accounting, typed Stops, Attention, activation recovery, and
portable import. It does not add a scheduler, automatic respawn, Human Outcomes,
Reviews, multiple selected Issues, or Envelope revision authoring.

The Phase 21 contracts remain authoritative except the explicit replacements
here, in `public-interface.md`, and in `schema.md`. These three documents are
accepted together under DEC-188. The exact
DDL amendment is incorporated by `03_TRD_SCHEMA.md` section 8.3. All new Events
retain `fid.event.v2`; the recovery archive becomes `fid.export-manifest.v2`.
No old fid-data migration, dual writer, or compatibility command is introduced.

Retain `internal/workflow` for Run, Stop, Attention, budget, and activation
rules, `internal/claims` for leases, `internal/evidence` for executor facts,
and `internal/validation` for shared evaluation. Put read-only Git observation
in the existing `internal/project` adapter. Do not create new domain packages
or add dependencies. The CLI owns parsing/rendering only.

## 2. Executor Observation And Lease Rules

A Run binds one immutable executor identity supplied at start by an activated
Actor with `run.start`. It is an opaque session identity, not a PID. Observations
bind that identity, Run, observing Actor, source reference, time, and current
Run version. A trusted activated harness with `run.observe` may report
`live`, `terminated`, or `unknown`; ordinary Builder attribution cannot.
fid observes its own gate child termination, but that alone never proves the
Builder session or all its other children have ended. Terminal session facts
are monotonic: a terminated executor identity cannot become live again.

Observation time cannot be in the future relative to fid's clock. Freshness
for a terminal observation follows immutable executor identity; wall-clock
lease expiry never substitutes for it. Imported `live` observations are
historical and cannot establish current liveness in the destination.

Claim expiry is explicit and journaled under expected version. Effective lease
validity nevertheless uses `now < expires_at`; equality is expired even before
an expiry Event. A released/expired Claim with an unresolved prior Run does
not permit acquisition of another mutating Claim. Direct acquire cannot bypass
recovery. Renewal increments Claim version, preserves Run identity, and cannot
revive an expired/released/superseded Claim. Execution requests validate the
current lease version, while the Run retains its original binding as history.

Authorized supersession occurs only inside `run.recover`; there is no generic
unfenced Claim takeover command. Same-Issue competitors serialize to one
winner. Unknown or overlapping scope rejects; cross-Issue parallel acceptance
is still Phase 24. Every stale-version/rejected authority request has no Event.

## 3. Repository Reconciliation

`run.inspect` is read-only. The repository adapter captures project/worktree
identity, HEAD/tree, dirty digest, changed paths and their before/after digests,
unmerged entries, and ongoing Git operation markers. It compares them with Run
start or the latest accepted checkpoint, approved Claim paths, and ancestry.
Use structured `git ... -z` output and byte-safe path handling. No source bytes
or raw command output enter an Event. Path arrays and references are bounded
by the wire contract; exceeding a limit rejects instead of truncating proof.

`run.checkpoint` records an activated harness observation, not a git commit.
At recovery the adapter observes again under the fid write lock and compares
the supplied observation digest with the actual state. Any intervening change
rejects `recovery-state-conflict`. Clean unchanged state qualifies. In-scope
dirty changes require explicit adoption of the exact observed path/digest set.
Unexpected paths, unmerged state, changed repository identity, missing base,
unattributable history, or an in-progress Git operation do not qualify.

Force authority may approve a precisely identified out-of-scope dirty
set; it does not waive a missing identity, unresolved merge, or a changing
observation. Human review must first make those states inspectable. The domain
records authority and adoption without modifying source files, Git refs,
index, stash, or other worktree content. Expected `.fid/` journal/projection
writes are operational metadata, excluded from the source-file comparison. It never resets, cleans, stashes, checks out, commits,
or discards files. Adoption creates no Receipt or Evidence Record.

Path classification is mechanical: approved Claim paths are in-scope, other
source paths are outside-scope, and Git metadata or `.fid/` internals are
protected. Protected metadata can never be adopted through force authority;
normal fid operational writes are not recovered source changes.

### 3.1 Canonical Dirty Digest

`dirty_digest` is `sha256:` plus lowercase hex SHA-256 of the canonical JSON
object `{schema:"fid.repository-dirty.v1",entries:[...]}`. Use the existing
Event canonical JSON grammar, UTF-8, and no trailing newline. Each entry is
`{path,head,index,worktree}`; entries sort by raw UTF-8 path bytes. Each image
is null for absence or `{mode,digest}`, where mode is the Git mode string
`100644`, `100755`, or `120000` and digest hashes the exact file/blob bytes.
For symlinks hash the link target bytes without following the link. HEAD/index
images use their blob bytes, never an OID as a content digest. Worktree images
use on-disk bytes with no newline, encoding, or clean/smudge normalization.

The input path set is every tracked path changed in the index or worktree
relative to HEAD, plus untracked non-ignored files. Enumerate tracked changes
and untracked files using structured NUL-delimited Git output; disable rename
detection so a rename is a deletion plus addition. Include tracked files even
when matched by ignore rules. Exclude untracked ignored files, root `.fid/`
and its descendants, and Git administrative storage; never hash operational
metadata. Git ignore rules select untracked inputs, not tracked ones. Mode-only
changes and staged changes restored in the worktree remain in the input set.

Empty input means `dirty:false,dirty_digest:null`; otherwise `dirty:true` and
the digest above. No absolute path, branch label, timestamp, Actor, Claim scope,
checkpoint ID, or observation time enters this digest. Invalid UTF-8/escaping
paths, unmerged index stages, changed submodules/gitlinks, unreadable files,
and unsupported file types fail closed with `recovery-state-conflict`; do not
omit them. Reject a repository using tracked assume-unchanged/skip-worktree
entries whose contents cannot be completely observed. Observation must detect
changes during capture and reject, never combine images from different states.

Retained caller-supplied Repository targets use this same algorithm; their
wire fields remain assertions, not a competing algorithm or an observation
authority. Existing Phase 21 synthetic digests are not a compatibility promise.
Task 4 owns the observer and fixed canonical-byte/hash vectors plus convergence
of displaced repository fixtures. Task 12 owns comparison against canonical
Receipt/Evidence targets, including a real unchanged dirty target and a
single-byte change with unchanged HEAD. Nonconforming historical digests are
never silently translated, refreshed, or treated as current proof. Rebuild
and import preserve them verbatim; canonical mismatch is stale proof.

### 3.2 Reconciliation Digest

Reconciliation.digest hashes the canonical JSON of the complete Reconciliation
object with only its `digest` member omitted, using the same `sha256:` encoding.
Its paths/classifications, checkpoint, repository identity, ancestry, unmerged
paths, and operation are therefore bound. Unlike dirty_digest, this digest is
scope/checkpoint-specific. Repeated unchanged inspection gives the same digest;
no capture time or random ID enters it. Recovery recomputes the full value under
the write lock. Neither digest makes a liveness claim.

## 4. Atomic Linked Recovery

Ordinary recovery needs a trusted terminal executor observation, a released or
expired prior Claim, active unchanged Envelope Revision, applicable authority,
and qualifying reconciliation. A still-active Claim requires its release first.
Unknown liveness always returns `executor-liveness-unknown`, including a force
request. Human force recovery requires an activated `human` Actor
with explicit `run.force-recover` capability and Envelope grant, a reason and
reference, exact current Envelope grant ordinal, and expected Run/Claim
versions. It may supersede a Claim and
authorize takeover of a known-live executor; it cannot fabricate a terminal
observation. Record the live observation and force provenance permanently, and
fence every subsequent fid mutation by the predecessor. fid cannot fence a
process's direct filesystem writes; this local operator action accepts that
residual risk explicitly and never claims that the old process was killed.

Every recovery and abandonment supplies a non-null existing observation ID
bound to its exact predecessor Run/executor. Missing/null IDs fail input
validation; nonexistent IDs return `observation-not-found`; wrong-subject
observations fail binding. Ordinary uses terminated; force may use known-live;
an unknown observation never qualifies.

One `run.recovered` Event atomically terminalizes an active predecessor as
`interrupted`, disposes the old Claim, creates a new Claim and active Run,
records predecessor/successor identities, and inserts the reconciliation fact.
A terminal `interrupted` or `failed` predecessor may also get one successor;
`succeeded` and intentionally `abandoned` Runs cannot. Retry after a failed
Run uses this same linkage. Same-key replay returns the original successor;
another key cannot create a second successor. A crash after append but before
projection is handled by ordinary catch-up, never a second recovery.

Issue state, Envelope Revision, Obligation policy, evidence identities, and
budget consumption survive. The successor's start repository is the reconciled
current target, not the Envelope's historical tree. Approved path/location
scope remains binding. Recovery itself starts no process or agent.

## 5. Stops And Attention

A Stop is an immutable fact plus append-only disposition. It binds a subject,
closed code, blocked actions, details, resolver capability, and a typed
resolution condition. Closed Phase 22 codes: `executor-liveness-unknown`,
`recovery-state-conflict`, `budget-exhausted`, `budget-data-unavailable`,
`activation-incomplete`, `capability-unavailable`, and `operator-stop`.
The last requires a real authorized policy/intervention reason, not a progress
message. A rejected read or ordinary validation request creates no Stop.
Explicit `stop.raise`, budget boundaries, and activation recovery transitions
are the only Stop producers in this slice.
In particular a rejected `run.recover` returning `executor-liveness-unknown`,
`recovery-state-conflict`, or `capability-unavailable` appends no Event and
raises no Stop/Attention. Those same-named Stop codes are produced only by an
explicit authorized `stop.raise`. Automatic budget/activation producers retain
their separately declared atomic effects. An error code is not a Stop fact.

Attention identity is SHA-256 of canonical `{subject,code,resolution_condition}`.
Equivalent open requests share one Attention ID; distinct Stop IDs and their
links remain. Resolved requests remain in history; a new recurrence creates a
new open item. Routine progress and warning thresholds never raise Attention.
There is no dismiss/delete command. A resolution references a qualifying durable
fact and checks resolver capability, subject, condition, and fact freshness.
Wrong-subject, nonexistent, or nonqualifying references reject with no effects.

Closed conditions are `executor-terminated`, `run-recovered`,
`activation-recovered`, `capability-available`, `budget-observed`, and
`work-ended`. Resolution stores the exact Event reference. A budget
exhaustion Stop allows only `work-ended` in this one-revision slice;
resolution records that the active execution has ended; the budget remains
exhausted and still denies future admission. Consumption/maximum cannot be
edited to clear it. It blocks further execution
and recovery, not linking already-observed proof or completing already-finished
work. Other Stops block precisely their declared actions. All applicable Stop
reasons participate in the shared evaluator. Resolving one Attention item
requires every linked Stop to have a qualifying disposition. Human Outcome
references remain unavailable until Phase 23; they are never silently accepted.

Qualifying resolution facts are closed:

| Condition | Required durable fact and present-state check |
|---|---|
| executor-terminated | run.observed for the same Run/executor with state terminated |
| run-recovered | run.recovered whose prior Run is the condition subject |
| activation-recovered | activation.recovery-completed or activation.finalized for the same activation, with current terminal component agreement |
| capability-available | current activation.finalized/recovery-completed selecting the recorded backend identity, plus current binary discovery of the named capability as available |
| budget-observed | run.budget-observed on the same Envelope/budget chain, with current known non-exhausted total |
| work-ended | run.closed or run.abandoned for the subject Run, or a terminal Run on the subject Envelope with no remaining active Run; envelope.completed also qualifies |

Budget and capability conditions bind their named budget/capability; arbitrary
Event references of the correct family cannot resolve another subject's Stop.
Attention resolution cites the corresponding stop.resolved Event IDs, including
an embedded Stop disposition's containing recovery Event where applicable.
The request is the sorted unique exact set of `stop_resolutions.created_event_id`
values for all linked Stops. A containing activation recovery Event may dispose
several Stops and appears once. Underlying `fact_event_id` values alone are
insufficient; missing, extra, duplicate, and wrong-subject references reject.
Revalidate each disposition's underlying fact and present-state condition.
`attention.resolved` changes only Attention state; it never inserts a Stop
disposition. Call `stop.resolve` first where no lifecycle Event disposed it.

## 6. Durable Budgets

Budget authoring retains the Phase 21 fields, but accepts
`exhaustion_action:"budget-exhausted"`. `consumed` and `remaining` remain rejected as
caller-supplied counters and must be null on Envelope input. The maximum is a
positive int64; warning threshold is null or 0..maximum. Scope is `envelope`.
Names are unique; units are `attempts`, `failures`, `milliseconds`,
`input-tokens`, `output-tokens`, `micro-usd`, `tool-failures`, `diff-bytes`,
and `interventions`. Monetary values are integers, never floating estimates.

Budget policy is immutable. Append-only samples carry a unique source sample
ID, Run, budget name, nonnegative delta or null, observed time, observing Actor,
and source plus a bound measurement evidence reference. The same source sample
cannot be charged twice, including across
recovery and import. fid automatically measures gate starts (`attempts`),
non-success Receipts (`failures`), elapsed execution (`milliseconds`), spawn/
harness errors (`tool-failures`), and force recoveries (`interventions`).
`run.observe-budget` accepts other observed metrics only from an activated
harness with the matching grant; it rejects the five fid-owned units. A numeric
external sample requires a Source reference/digest to the bound measurement;
fid-owned samples reference their Attempt/Receipt/recovery through source sample
identity and have null measurement_ref. Telemetry
summaries and provider integration remain Phase 25. Unknown metrics remain
null, not zero. fid-owned counters start at known zero; external units remain
unknown until a trusted numeric sample establishes their baseline.

Consumption is the sum of distinct known samples; remaining is
`max(0, maximum-consumed)`. An explicit unknown sample makes the total unknown
until a later trusted baseline sample covers the unknown interval; the baseline
is cumulative for the same Envelope/budget and may not decrease a known total.
Samples distinguish `delta` from `baseline`, with an exact predecessor sample
ID; concurrent/out-of-order baselines reject. No recovery resets the chain.

Before work, unknown required budgets produce `budget-data-unavailable` with
Stop/Attention and no Attempt. Admission of the last allowed attempt is atomic
with charging that Attempt and raising its exhaustion Stop/Attention; that
already-admitted process may finish and persist its Receipt. Subsequent starts
reject and do not spawn. A failed Receipt and any resulting charges, Stop, and
Attention commit together. A resource already above maximum remains above it
in history; remaining is zero. The operation response explicitly identifies
committed Stop/Attention even when the requested execution was denied.
Warning threshold crossing is a queryable state only. No automatic retry.

## 7. Interrupted Execution And Run Closure

An Attempt is durable before spawn and has at most one fid-observed Receipt.
Cancellation, timeout, signal, spawn failure, and harness failure cannot pass.
Cancellation adds the terminal outcome `cancelled`. Losing fid before observed
termination leaves an abandoned Attempt with no synthetic Receipt. Maintenance
reports it; recovery preserves it. Reusing the original gate idempotency key
never respawns an abandoned Attempt. A new Run and new key are required.

Replace caller-controlled `run.finish` outcomes with `run.close` and
`run.abandon`. Close derives success/failure from durable trusted execution and
shared Obligation evaluation, rejects outstanding Attempts or unknown executor
liveness, and references observed facts. Satisfied proof permits success;
observed failed execution permits failure; missing/stale proof cannot produce
success. Abandon records intentional cessation with trusted terminal executor
proof (or Human force authority), releases its Claim, and creates no evidence.
Normal close terminalizes only the Run; Issue transition and Envelope completion
retain their separate shared guards and do not gain an implicit success flag.

## 8. Activation Recovery

Activation carries an optimistic version starting at 1. Every component,
resource, lifecycle, and recovery mutation checks the expected activation
version and increments it in the same Event. Resource observations additionally
check their own version. Prepare binds the exact prior activation ID/version,
or both null for first activation. Terminal history remains immutable in Events;
the current projection marks a predecessor superseded only through a new
activation's successful forward-finalization Event.

Prepare may follow one finalized activation and freezes its exact predecessor
identity/version and the old/new identities of the five declared resources:
kit, overlay, provider, project and backend. The predecessor identifies its
immutable authority-profile and Gate Definition components; first activation
has no predecessor component set. New internal component content is supplied
through the retained `activation.component` operations under the same new
activation ID. Staging validates and content-addresses that content, which
remains immutable. Commit binds the exact two staged internal component digests
and checks the required resource agreement before the point of no return.
Finalization and recovery cannot mix old and new internal component sets.
DEC-198 makes these distinct operation boundaries explicit: prepare contains
no internal component payload, and neither staging nor recovery rewrites its
frozen resource/predecessor bindings.

The kit owns actual file materialization; fid persists content identity and
observed component state. A component observation binds identity/digest,
expected version, state, and evidence reference. Bootstrap attribution remains
explicitly untrusted; fid is a local integrity boundary, not authentication for
the local operator.

One nonterminal activation blocks normal kit work even while the previous
identity is retained. Inspection, discovery, export/validation, and activation
recovery remain usable. Commit is the durable point of no return. Before it,
recovery selects `restore-prior`; with no predecessor it restores the unactivated
state. After it, only `complete-forward` is legal. Selection is immutable and
retries retain the same activation ID. No recovery command edits provider or
project files. The kit records each component at the chosen target, and
`activation.recover-complete` verifies exact agreement before ending recovery.

Restoring prior ends the attempted activation as aborted and preserves prior
trust. Completing forward supersedes prior trust and finalizes the new set in
one Event. Missing/mismatched components leave normal workflow blocked. Failed
finalization never claims rollback after commit. Actor and Registration
resolution use the selected finalized identity only; no mixed profile is valid.

## 9. Freshness And Resume

Freshness compares the current observed repository target, including dirty
digest, with each Receipt/Evidence target; branch labels are not identities.
Recovery and adoption do not refresh proof. An unchanged target retains valid
proof; a changed target makes repository-bound proof stale. Preserve validity,
absence, and Artifact availability as separate dimensions. No Review or
external-artifact-target Obligation is newly enabled here; the external-artifact
exception in PRD-EVID-003 stays explicitly deferred until that target exists.

Resume uses current repository observation and lists predecessor/successor Run
links, effective lease status, unresolved Stops, deduplicated Attention,
outstanding Obligations, evidence dimensions, and remaining budgets. Next
actions are computed from the same predicates as mutations; blocked mutation
is never advertised as executable. Findings remain an explicitly unavailable
dimension until Phase 23. Read-only inspection appends no Event.

## 10. Portable Recovery

Manifest v2 adds closed project configuration and includes every accepted
Event family, recovery link, Stop/disposition, Attention link, budget sample,
activation/component fact, idempotency result, and Artifact identity. Ordered
Event bytes and IDs are preserved. Referenced project prose is not copied.
The format contains no SQLite image or SQL; its closed Event-v2 logical Effects
are interpreted by the fake and fid independently. Unknown tables, columns,
operations, Event types, sequences, references, or schemas reject.

Validation is non-mutating and bounded: at most 1000000 entries, 64 MiB per
non-Artifact entry, 1 GiB per Artifact, 8 GiB total expanded bytes, 1 MiB per
Event line, and no symlink/device/directory/path traversal entries. Read limits
apply to actual decompression, not only ZIP metadata. Chunk Events at line
boundaries to respect entry limits. Duplicate keys, duplicate immutable IDs,
inconsistent counts/digests, reference breaks, or unsupported capability values
fail validation. Resource ceilings are proposed contract limits, not test-only
tolerances. A validator uses the exact version's capability matrix, never
whichever binary happens to validate the file.

Import accepts only a target repository without `.fid/`; it never merges into
or replaces a populated project. Validate before publication, stage a sibling
directory on the same filesystem, reconstruct through normal appliers, verify
semantic state, fsync, and atomically publish `.fid/` under an exclusive target
lock. A failed attempt exposes no partial `.fid/`. If publication completed
but the response was lost, retry with the same archive digest/key reports the
same installed state without new Events. A different digest/key conflicts.
The publication marker is delivery metadata, excluded from canonical export.

Preserve project/entity/Event identity, order, hashes, history, and policy;
absolute archive/destination paths and publication metadata are not semantic
state. Destination repository observations are new local facts, not rewrites
of imported history. Imported active work requires liveness/reconciliation
before mutation; import itself cannot authorize takeover or automatic spawn.
Validation includes referential integrity and accepted effect semantics before
publication. Artifacts available at export must have matching bytes; missing
or expired artifacts retain metadata and availability without invented bytes.

## 11. Availability And Acceptance

After the owning conformance cases pass, `attention.v1`, `kit.activation.v1`,
and `export.portable.v1` become available. `work.claims.v1` remains unavailable
until Phase 24 proves the required parallel scope semantics; its recovery limit
becomes true and its owner becomes Phase 24. This is an explicit proposal
resolving the roadmap's conditional advertisement, not partial availability.
The other eight rows retain their Phase 21 availability.

Shared fake/fid conformance runs the same versioned request, result, error, and
negative cases for the three newly available capabilities and the recovery
subset of Claims. The fake is test-only. Cross-product consistency review is
required before acceptance; it must identify the kit port contract/version it
compared and cannot substitute an invented production fake. Phase 25 still
owns qualification of all twelve capabilities together and first adoption.
The accepted comparison dispositions are in `kit-contract-mapping.md`
(DEC-187). In particular fid activation state does not bind a recovery archive:
the kit coordinator must bind the archive identity/digest and fid activation ID
in its own durable activation record before any destructive materialization.
Portable export alone does not prove that binding. Phase 25 qualification must
exercise the coordinator's binding and terminal-outcome mapping; Phase 22
advertisement is limited to the documented fid operations and is not kit port
qualification.

## 12. Config Decision

CF-035 proposal: validate the entire closed configuration whenever an operation
consumes project config, before it uses any value or writes. Help, version,
archive validation, and config-independent diagnostics remain available for
repair. Do not add a global parser to every command. Evidence must show both
unknown-key rejection on a consuming path and continued availability of a
config-independent path. This pins the existing consumption boundary explicitly.
