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

Each requirement has two identities. `requirement_id` is the project's stable
ID (for example `REQ-leaf`); `id` identifies one version of the reference
(`RREF-…`), and every supersede creates a new one. Every `RequirementRef`
carries both (Phase 23 DEC-269). `requirement get <id>` accepts either: a
stable ID returns its active reference, or else its newest; a reference ID
returns that row. Human views name the requirement by its stable ID, show the
reference ID separately, and list Issue links in their own labelled section.

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

`envelope create --input` omits `repository_scope` and every Source
`repository`. fid observes the selected project's Git worktree and records that
observation as the scope and on `risk_result.source`, `sources[]`, and every
Obligation Source. A project-file Source may omit `digest`; fid reads the file
and records its SHA-256 digest, or checks a supplied digest. External
references require a digest. `envelope complete --input` names the Envelope
and expected revision; fid observes the completion repository itself.

An Obligation's `registration` is nullable: gate-evidence names its immutable
Registration, while Review and Human Outcome Obligations supply null. Prior
Registrations and Envelope scopes must name the same project and worktree,
or an imported destination whose binding names that project and exact root.
A plain copy with no import binding is refused without an Event. Claim Events
retain the Envelope scope; the acquired Claim records its local observed
worktree. A Claim acquired in an imported destination may be renewed,
released or expired there; stored authority uses the same shared worktree
predicate as acquisition (DEC-283).

### Claims, Runs, Stops, and Attention

```text
fid claim acquire|renew|release|list
fid run start|observe|observe-budget|checkpoint|close|abandon|get|list|inspect|recover
fid stop record|resolve|list
fid attention list|resolve
```

No command assumes a Developer or QA role. Actor identity, capabilities,
Envelope grants, Issue, worktree, mutation scope, optional Obligation, and
lease version are explicit machine fields.

`claim acquire --input` names the Issue, Envelope revision, path scope, mode
and lifetime. fid observes and records the selected worktree ID, checking it
against the Envelope scope. Renew, release and expire retain their existing
inputs.

`run start`, `run close` and `run abandon` require
`expected_repository:{head_oid,dirty}` instead of `repository`. Use
`git rev-parse HEAD` and `git status --porcelain` (excluding `.fid/`) for the
expectation. fid observes under the write lock after replay lookup, compares
these two preconditions, and records the full observation. A mismatch returns
`repository-mismatch` with `{expected:{head_oid,dirty},observed:Repository}`
and no Event. Start checks the Claim and Envelope as the same project and
worktree, including a matching import destination binding; edits and commits
since Envelope creation are accepted. A plain copy has no such binding and is
refused before an Event. Close and abandon record `terminal_repository` with
`base_oid` equal to observed HEAD. Close evaluates Obligations against that
observation; stale Receipt-backed Evidence yields a failed Run. Unsafe or
unobservable terminal repositories return `recovery-state-conflict`.

`run observe-budget` accepts nullable `measurement_ref:SourceInput`. With a
Source, fid observes the repository and checks a supplied project-file digest
or computes an omitted digest; the stored Source and Event carry that
observation. A null measurement does not observe. Start and budget observation
return `repository-unavailable` on an unobservable worktree. Equal keyed
retries replay before any repository or Source read, even after HEAD changes,
a referenced file changes or disappears, or the worktree becomes unsafe.

All seven Run mutation leaves document their closed input fields in `--help`
and provide a standalone, exclusive `--example` JSON template. Structural
failures name the exact `input.<path>` and do not append an Event.

### Registered gates and evidence

```text
fid gate register|resolve|run
fid receipt get|list|ingest-attested
fid evidence link|get|evaluate
```

`gate run --input` accepts `registration_id`, `registration_version`,
`obligation_id`, `run_id` and required `expected_repository:{head_oid,dirty}`.
Use `git rev-parse HEAD` and `git status --porcelain` for the expectation.
fid observes once before spawning and checks the expectation before any
Attempt or Event. Registration and Run checks use the same project/worktree
relationship, including a matching import binding. Edits since Run start are
admitted with a matching expectation. Attempt, Receipt and execution Events
retain that observation; process exit does not cause a second observation.
The `fid.live.v1` framing contract is unchanged.

`evidence link --input` supplies the Obligation, Evidence type, exactly one
non-null subject ID (`receipt_id`, `review_id`, `human_outcome_id`) and sorted
`required_artifact_ids`. It derives `target_repository` from the Receipt's
terminal repository, the Review repository or the Human Outcome Event.
Linking does not observe the current worktree. `evidence evaluate --input`
takes `{obligation_id:?ID}` (null selects the active revision) and observes
the current worktree once. Tree and dirty-state changes stale repository-bound
Evidence; reverting restores freshness. Human Outcome Evidence stays fresh.

These inputs reject the removed `repository` and `target_repository` members.
`gate run`, `evidence link`, `evidence evaluate` and `artifact verify` document
their field trees in `--help` and print exclusive editable `--example`
templates without reading or writing a workspace. Structural failures name
`input.<path>` and append no Event.

`gate register --input` names the activated definition/version and Issue scope.
fid observes the selected worktree and derives the Registration project ID,
worktree ID and binding digest from that observation and scope.

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

Phase 24 Task 5 makes `review record` take required
`expected_repository:{head_oid,dirty}` in place of a caller-supplied
`repository`. Use the full output of `git rev-parse HEAD` for `head_oid` and
set `dirty` from `git status --porcelain`, ignoring `.fid/` paths. The
coverage item and exclusion Sources, and `evidence_refs`, take `SourceInput`
without `repository`; project-file digests may be omitted for fid to compute.
The stored Review and its Sources carry fid's single observed Repository.
For example, the repository portion of a Review input is:

```json
{"expected_repository":{"head_oid":"0123456789abcdef0123456789abcdef01234567","dirty":false}}
```

A mismatch returns `repository-mismatch` with exactly
`{expected:{head_oid,dirty},observed:Repository}` and appends no Event.
`finding address` keeps its input shape and gains descriptor help and
`--example`.

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

For a keyed mutation, the request digest covers the closed input object as
written, plus the operation, project ID and resolved Actor components already
used by its family. Optional members enter the digest only when supplied;
`review.record` treats `inapplicable_findings:[]` as omission (DEC-242). Values
observed by fid, including repository state and Source digests, do not enter
the request digest. A same-key replay returns the original result before any
observation. Each keyed mutation also records a `sha256:` digest of the
canonical JSON value for every top-level public input name and `--actor`.
Flag names include `--`; positional names follow the leaf's `Use` line.
An `idempotency-conflict` JSON details object contains the operation scope,
key, stored and supplied request digests, and `differing_fields`: the sorted,
unique, nonempty names whose field digest or presence differs. The conflict
appends no Event.

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
exact details `envelope_id`, `decision_class`, `required_capabilities`,
`affected_actions`, and (DEC-263) the Stop-blocked denial fields `stop_ids`,
`attention_ids`, `event` (always null), and `replayed` (always `false`); and `source-digest-mismatch` (precondition, exit 4) with
exact details `display_label`, nullable `path`, `expected_digest`, and nullable
`observed_digest`. A missing Source file is reported through the same digest
mismatch code and appends no Event.

Phase 23 Task 10 adds `human-outcome-required` (precondition, exit 4) with
exact details `envelope_id`, `obligation_id`, `outcome_class`, and (DEC-263)
the same four Stop-blocked denial fields, returned by
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
| `protected-decision-required` | authority / 5 | envelope | `envelope_id`, `decision_class`, `required_capabilities`, `affected_actions`, `stop_ids`, `attention_ids`, `event`, `replayed` |
| `outcome-not-permitted` | precondition / 4 | null | `obligation_id`, `obligation_type`, `outcome_id`, `result` |
| `evidence-type-not-accepted` | precondition / 4 | null | `obligation_id`, `evidence_type`, `accepted_evidence_types` |
| `source-digest-mismatch` | precondition / 4 | null | `display_label`, `path`, `expected_digest`, `observed_digest` |
| `repository-mismatch` | conflict / 4 | null | `expected:{head_oid,dirty}`, `observed:Repository` |
| `human-outcome-required` | precondition / 4 | envelope | `envelope_id`, `obligation_id`, `outcome_class`, `stop_ids`, `attention_ids`, `event`, `replayed` |

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

### Phase 23 accepted amendment

`../docs/phases/phase-23/public-interface.md` contains the delivered closed
Review, Finding, Concern, Deviation, Decision Reference, Human Outcome, and
requirement-reference leaf schemas, failure details, and Human views (DEC-222).
Later accepted amendments, all delivered and recorded in that contract:

- DEC-263: Stop-blocked judgment denials carry `stop_ids`, `attention_ids`,
  `event` (null), and `replayed` (`false`) in failure details; StopDetails are
  unchanged.
- DEC-265/266: section 5 declares 19 Failure codes (18 added, 1 retained,
  tabled above). All four `authority-*` codes keep the Phase 21 details
  `{operation,actor_id,required_actions,required_capabilities}` with a null
  subject.
- DEC-269: every Human-mode failure from the Stop-error family (`claim
  acquire`, `run start`, `gate run` Stop denials, `stop` and `attention`
  leaves) writes `error: <message>` to stderr, with the JSON `message` text.
  The `review-rework-exhausted` Stop is resolved by an `approved` Human Outcome
  of class `review-rework-extension`, recorded atomically or cited through
  `stop resolve`. Its Human denial adds one guidance line naming that action.
  `RequirementRef` gains `requirement_id` as described above.
- DEC-270: journals written before `requirement_id` existed are unsupported;
  rebuild refuses them with `journal-replay-stopped` (ADR-0011).

DEC-271 records technical qualification and the Human re-check. DEC-272
schedules the remaining Human operating-experience items (Stop-detail
rendering, denial reasons, hint verbosity, recorded reasoning in Human views,
`--input` authoring and repository-input convergence) as Phase 24.

### Phase 24 accepted amendment (implementation pending)

`../docs/phases/phase-24/public-interface.md` is the accepted target contract
(DEC-277, amended by DEC-278..282); the preceding delivered rules describe
the implementation until the owning Phase 24 tasks land. Sections 5 and 9
govern the message clarifications recorded in DEC-282:

- A plain-copy `run.start` refusal reaches the Claim binding check first.
  Its generic fallback carries `run claim binding is invalid` once task 10
  applies D6=A. Task 7 preserves the Claim-before-Envelope check order and
  the original/copy/import-binding cases.
- Task 9a constructs grammar and unknown-command messages from their
  specified source text, escapes controls, then enforces the inherited
  512-Unicode-scalar bound. Escaped text up to 512 stays intact; longer text
  keeps the first 511 plus `…`, dropping an escape crossing the cut whole.
  Count before JSON serialization, excluding the Human label/terminator.
  This bound takes precedence over exact parser text or the unknown-command
  formula. Structured details retain the full offending spelling/value.
- Task 10 reuses that construction for DEC-281 generic fallbacks and prints
  the resulting JSON message verbatim after `error: `, without a second
  escape or cut. Codes, details and exits retain their accepted values.

The exact leaf, input, failure and view changes remain indexed in that phase
contract; this entry records accepted planning, not completed delivery.

Task 3 adopts descriptor-generated input help and `--example` templates for
`stop raise`, `stop resolve`, and `attention resolve`. Each adopted leaf's help
lists field paths, JSON types, required/optional/nullable status, declared
limits and allowed values, and the template command. `--example` prints a
two-space-indented JSON object with required fields in documented order; it
works outside a project and cannot be combined with another flag or argument.
Structural `--input` failures name the field in `invalid-argument.details.argument`
(for example `input.resolution_condition.subject.id`) before domain validation.
The remaining `--input` leaves enter this convention in their owning tasks.

Task 9 adopts the same descriptor help, `--example`,
and field-path structural validation for `activation prepare`, `commit`,
`abort`, `finalize`, `resource`, `recover`, and `recover-complete`. The input
shapes retain the Phase 21 baseline and Phase 22 version/resource additions.
`activation component` selects a raw typed `content` array using
`component_type`: `authority-profile` contains ActorProfiles, and
`gate-definitions` contains GateDefinitions. Its help lists both variant
trees, including nested environment entries and bounds. Its deterministic
example uses the first variant, `authority-profile`, with one ActorProfile.
Structural errors identify the selected nested path; mismatched content
identifies `input.content[i]`. Environment values and redaction literals
remain sensitive in failure details. Raw variant string bounds count Unicode
scalars; session references and sensitive literals have no additional text
limit. cwd accepts `.`, while Artifact paths exclude dot segments and are
bounded to 512 bytes. All 39 input leaves use descriptors
for help, examples and closed structural validation.

Task 4 adopts `SourceInput` for `decision record-reference`, `requirement
record-reference`, `human-outcome record`, `concern record`, and `deviation
record`. It is a Source without `repository`; a project-file digest may be
omitted, while an external-reference digest is required. fid resolves each
project-file path under the project root, checks a supplied SHA-256 digest or
records the observed one, and attaches its own repository observation to the
stored Source and Event. `requirement link` and `concern dispose` gain
descriptor help and templates as observing leaves without Source fields.

For example, `fid decision record-reference --example` shows the complete
required input shape. Its `source` can be edited to name a file as follows;
the optional digest is shown only here to illustrate a caller-supplied check:

```json
{"source_type":"project-file","display_label":"Accepted decision","path":"project-docs/decision.md","section":null,"external_id":null,"digest":"sha256:<64 lowercase hexadecimal digits>"}
```

Omitting `digest` asks fid to observe the file. A same-key retry replays the
original Source and Event even after HEAD or the file changes. A mismatch
returns `source-digest-mismatch` with nullable `expected_digest` and
`observed_digest`; an unobservable worktree returns `repository-unavailable`
with a typed reason and detail. Both refuse without appending an Event.

*Last updated: 2026-10-05 (Phase 24 Task 4 judgment SourceInput delivery)*
