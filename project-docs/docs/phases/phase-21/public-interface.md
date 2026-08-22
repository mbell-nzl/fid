# Phase 21 Public Interface Contract

> Status: **PLANNING DRAFT FOR HUMAN APPROVAL**
> Owner: Tech Lead
> Authority: DEC-153, DEC-154, DEC-156, DEC-157
> Applies to: the final Phase 21 binary after migration 0014

This file is the complete public CLI contract for Phase 21. Product code may
not infer a request, result, error, rendering, or stream shape not stated here.
The semantic and storage rules are in `contracts.md`; this file owns the wire
surface. All closed public enum members are lowercase kebab-case. JSON field
names remain snake_case. User-authored opaque IDs, paths, and workflow-state
identifiers are not enums.

## 1. Global Transport

```text
fid [--project <path>] [--json] <leaf> [leaf arguments and flags]
```

- Project discovery and `--project` behavior remain exactly Phase 20 contract
  section 2.1. `fid init` treats the selected directory as the target.
- Human success is stdout only and exit 0. Human failure is stderr only and
  uses `error: <message>` plus at most two non-semantic guidance lines.
- Non-streaming JSON success is one canonical `fid-cli.v1` object plus LF on
  stdout, empty stderr, exit 0. JSON failure is one failure object plus LF on
  stderr, empty stdout, and the exit class in section 7.
- `fid gate run --json` is the sole streaming exception described in section
  8. Before its first frame, validation failures use the ordinary failure
  channel. From the first frame onward, stdout contains only `fid.live.v1`
  frames, stderr stays empty, and the terminal frame determines exit status.
- JSON is UTF-8. Object keys serialize lexicographically. Arrays retain the
  declared order. Timestamps are UTC RFC3339Nano. Optional nullable values are
  emitted as explicit `null` in results.

Success and failure envelopes are exactly:

```json
{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"capabilities.get","result":{"items":[]}}
{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["no arguments"],"argument":"arguments","value":"extra"},"message":"capabilities accepts no arguments","retryable":false,"subject":null},"ok":false,"operation":"capabilities.get"}
```

```text
EventRef = {id:string,occurred_at:timestamp}
ErrorSubjectType = "activation"|"actor"|"artifact"|"capability"|"claim"|
  "envelope"|"evidence"|"export"|"gate-definition"|"gate-registration"|
  "idempotency"|"issue"|"label"|"link"|"obligation"|"path"|"project"|
  "receipt"|"repository"|"requirement"|"run"|"workflow"
ErrorSubject = {type:ErrorSubjectType,id:string|null}
ErrorCategory = "authority"|"conflict"|"integrity"|"internal"|"invalid-input"|
  "not-found"|"precondition"
FailureCode = exactly one code from the closed table in section 7
SuccessEnvelope<R> = {contract_version:"fid-cli.v1",event:EventRef|null,
  ok:true,operation:string,result:R}
FailureEnvelope<D> = {contract_version:"fid-cli.v1",error:{category:ErrorCategory,
  code:FailureCode,details:D,message:string[1..512],
  retryable:boolean,subject:ErrorSubject|null},ok:false,operation:string}
```

`event` is `null` for reads and projection-only operations. A committed
single-Event mutation emits exactly `{"id":string,"occurred_at":timestamp}`.
A replay emits the original Event object. A failure object has no `result` or
`event`; a success object has no `error`.

## 2. Closed Input Rules

The notation below is normative:

- `name:T` is a required non-null field; `name:T|null` is required nullable;
  `name?:T` is optional and non-null when present.
- `{...}` is closed. Unknown or duplicate keys reject `invalid-argument`.
- `--input <path|->` reads exactly one JSON object, at most 1,048,576 bytes,
  followed only by JSON whitespace and EOF. `-` reads stdin. Invalid UTF-8,
  multiple values, a scalar/array root, or a larger input rejects.
- String bounds count Unicode scalar values after trimming only where a field
  says `trimmed`. IDs are opaque printable strings of 1..128 characters.
- `digest` is `sha256:` plus 64 lowercase hexadecimal digits. Git OIDs are 40
  lowercase hexadecimal digits. Relative paths are slash-normalized, nonempty,
  contain no `.`/`..` segment, and are at most 512 bytes.
- `cwd-path` is exactly `.` or a relative path under the selected worktree;
  Artifact paths use the stricter relative-path rule and resolve under cwd.
- Arrays are required unless marked optional. `sorted-unique` means bytewise
  ascending with no duplicate. Empty is accepted only where stated.
- Missing required arguments/flags, forbidden or unknown flags, extra
  arguments, mutually exclusive flag violations, and malformed flag values
  return `invalid-argument` with the offending spelling/value and the row's
  exact accepted alternatives in `allowed`; they perform no read or write.
- Every Event-producing mutation except `project.init` requires
  `--idempotency-key <key>` where key is 1..128 printable ASCII, no leading or
  trailing whitespace. `project.init` has no key because workspace creation is
  exclusive; its Event carries `idempotency:null`.
- Kit mutations require `--actor <actor-id>`. Retained standalone mutations
  require `--actor <type>:<id>`, where type is exactly `agent`, `human`, or
  `service` and ID matches `[A-Za-z0-9][A-Za-z0-9._-]{0,127}`. The retained
  value becomes an Event Actor with `producer:"fid-standalone"`,
  `provenance_ref:"project-local-attribution"`, and `session_ref:null`; it
  records attribution and never grants kit authority. Caller-supplied
  capability or provenance fields reject.

### 2.1 Shared request values

```text
Repository = {project_id:string,worktree_id:string,capture_source:"git",
  head_oid:string|null,tree_oid:string|null,base_oid:string|null,
  dirty:boolean,dirty_digest:digest|null}
Source = {source_type:"project-file"|"external-reference",display_label:string[1..128],
  path:string|null,section:string|null,external_id:string|null,digest:digest,
  repository:Repository}
ActorProfile = {actor_id:string,actor_type:"agent"|"human"|"service",
  producer:string[1..128],provenance_ref:string[1..512],session_ref:string|null,
  capabilities:string[1..128][] sorted-unique nonempty}
EnvironmentEntry = {name:^[A-Z_][A-Z0-9_]{0,127}$,
  source:"credential"|"inherit"|"literal",value:string|null}
GateDefinition = {definition_id:string,definition_version:int[1..2147483647],
  argv:string[1..4096][][1..128],cwd:cwd-path,shell:false,
  environment:EnvironmentEntry[] sorted by name,timeout_ms:int[1..600000],
  output_bytes:int[1..10485760],excerpt_bytes:int[0..min(output_bytes,16384)],
  success_exit_codes:int[-255..255][] sorted-unique nonempty,
  credential_names:string[] sorted-unique,redact_literals:string[] sorted-unique,
  artifact_paths:relative-path[] sorted-unique,artifact_bytes:int[1..1073741824],
  sandbox:"none",network:"allowed",write:"allowed"}
Budget = {name:string[1..64],scope:"envelope",unit:"attempts"|"milliseconds"|
  "tokens",maximum:int[1..9223372036854775807],warning_threshold:int|null,
  exhaustion_action:"unavailable",consumed:null,remaining:null}
AuthorityGrant = {actor_selector:{actor_id:string,
  required_capabilities:string[] sorted-unique nonempty},
  actions:string[] sorted-unique nonempty,scope:{issue_id:string}}
ObligationInput = {id:string,type:"gate-evidence",applicability:"required",
  guarded_actions:["envelope.complete","issue.transition-terminal-success"],
  subject:{type:"issue",id:string},registration:{id:string,version:int>=1},
  accepted_evidence_types:["receipt"],accepted_trust_sources:["fid-executed"],
  freshness_target:"repository-tree",artifact_policy:"required-available",
  independence:"fid-executed",required_lane:null,source:Source}
```

`EnvironmentEntry.value` is required non-null only for `literal`, required
null for `credential` and `inherit`. `credential_names` must exactly equal the
names of credential entries. Inherited names are limited to `CI`, `GOCACHE`,
`GOMODCACHE`, `GOPATH`, `HOME`, `LANG`, `LC_ALL`, `PATH`, `TMPDIR`, and
`XDG_CACHE_HOME`. A credential named `NAME` is read only from the fid process
variable `FID_CREDENTIAL_NAME`, is passed to the child as `NAME`, and never
falls back to an ambient `NAME`; a missing source fails before Attempt append.
No environment name may repeat across sources. `warning_threshold` is null or
1..maximum.

For every registration, Envelope, Claim, Run, Gate, Evidence, transition,
resume, and export use, `Repository.head_oid` and `tree_oid` are non-null;
`dirty_digest` is non-null exactly when `dirty:true`; and project/worktree IDs
must match the selected project observation. `base_oid` remains nullable. A
`project-file` Source has non-null path and null external ID; an
`external-reference` has null path and non-null external ID. Section is
independently nullable. Invalid cross-field states reject before mutation.

Activated capability and Authority Grant action names are limited to
`artifact.verify`, `claim.acquire`, `claim.release`, `claim.renew`,
`envelope.activate`, `envelope.complete`, `envelope.create`, `evidence.link`,
`gate.register`, `gate.run`, `issue.transition`, `run.finish`, and `run.start`.
`gate.register` and `envelope.create` are authorized directly by their matching
finalized-profile capability because no Envelope Grant exists yet. Every other
mutating kit action requires both the finalized-profile capability and the
exact active Envelope Grant dimension.

### 2.2 Exact `--input` schemas

| Schema | Exact closed object and ordering/bounds |
|---|---|
| `ActivationPrepareInput` | `{kit_identity:string[1..128],backend_identity:string[1..128],phase:21,direction:"activate"}` |
| `ActivationComponentInput` | `{activation_id:string,component_type:"authority-profile"|"gate-definitions",component_id:string,component_version:int>=1,expected_prior_digest:digest|null,content:ActorProfile[]|GateDefinition[]}`; content is nonempty and sorted by Actor ID or Definition ID/version; union must match component type. |
| `ActivationCommitInput` | `{activation_id:string,expected_status:"staged",component_digests:digest[2] sorted-unique}` |
| `ActivationAbortInput` | `{activation_id:string,expected_status:"prepared"|"staged",reason:string[1..256]}` |
| `ActivationFinalizeInput` | `{activation_id:string,expected_status:"committed"}` |
| `GateRegisterInput` | `{definition_id:string,definition_version:int>=1,worktree_id:string,scope:{issue_id:string},repository:Repository}` |
| `EnvelopeCreateInput` | `{outcome:string[1..2048],non_goals:string[0..512][][0..64],risk_result:{source:Source,result_digest:digest},repository_scope:Repository,sources:Source[1..64],primary_issue_id:string,related_issue_ids:[],authority_grants:AuthorityGrant[1..64],obligations:ObligationInput[1],budgets:Budget[0..64]}`; sources preserve caller order; grants sort by actor ID; budgets sort by name. External `Routine` lane classification is neither accepted nor stored. |
| `EnvelopeActivateInput` | `{envelope_id:string,expected_revision:1}` |
| `EnvelopeCompleteInput` | `{envelope_id:string,expected_revision:1,repository:Repository}` |
| `ClaimAcquireInput` | `{issue_id:string,envelope_id:string,envelope_revision:1,worktree_id:string,scope:{paths:relative-path[] sorted-unique nonempty},ttl_seconds:int[1..86400]}` |
| `ClaimRenewInput` | `{claim_id:string,expected_version:int>=1,ttl_seconds:int[1..86400]}` |
| `ClaimReleaseInput` | `{claim_id:string,expected_version:int>=1,reason:"completed"|"failed"|"operator-released"}` |
| `RunStartInput` | `{issue_id:string,envelope_id:string,envelope_revision:1,claim_id:string,claim_version:int>=1,repository:Repository}` |
| `RunFinishInput` | `{run_id:string,expected_status:"active",outcome:"abandoned"|"failed"|"succeeded",repository:Repository}` |
| `GateRunInput` | `{registration_id:string,registration_version:int>=1,obligation_id:string,run_id:string,repository:Repository}` |
| `EvidenceLinkInput` | `{obligation_id:string,receipt_id:string,target_repository:Repository,required_artifact_ids:string[] sorted-unique}` |
| `EvidenceEvaluateInput` | `{obligation_id:string,target_repository:Repository}` |
| `ArtifactVerifyInput` | `{artifact_id:string}` |

## 3. Closed Result Types

All objects below are closed. Array order is part of the contract.

```text
MutationMeta = {idempotency_key:string,replayed:boolean}
Capability = {capability_id:string,schema_min:int,schema_max:int,available:boolean,
  limits:closed-object-from-contract-matrix,trust_source:TrustSource}
TrustSource = "active-kit-activation"|"fid-executed"|"fid-git-observation"|
  "fid-journal"|"phase-22-unavailable"|"phase-23-unavailable"|
  "phase-25-unavailable"
Actor = {actor_id:string,actor_type:"agent"|"human"|"service",producer:string,
  provenance_ref:string,session_ref:string|null}
Activation = {id:string,status:"aborted"|"committed"|"finalized"|"prepared"|
  "staged",kit_identity:string,backend_identity:string,phase:21,direction:"activate",
  point_of_no_return:boolean,component_digests:digest[] sorted,created_at:timestamp,
  updated_at:timestamp}
Registration = {id:string,version:int,activation_id:string,definition_id:string,
  definition_version:int,definition_digest:digest,binding_digest:digest,
  project_id:string,worktree_id:string,scope:{issue_id:string},
  authorization_provenance:string,policy:GateDefinition,created_at:timestamp}
BudgetView = Budget
Obligation = ObligationInput + {derived_status:"satisfied"|"unsatisfied",
  reasons:Reason[]}
Envelope = {id:string,current_revision:1,status:"active"|"completed"|"draft",
  outcome:string,non_goals:string[],risk_result:{source:Source,result_digest:digest},
  repository_scope:Repository,sources:Source[],primary_issue_id:string,
  related_issue_ids:[],authority_grants:AuthorityGrant[],obligations:Obligation[1],
  budgets:BudgetView[],created_actor_id:string,approved_actor_id:string|null,
  created_at:timestamp,activated_at:timestamp|null,completed_at:timestamp|null}
Claim = {id:string,issue_id:string,envelope_id:string,envelope_revision:1,
  actor_id:string,worktree_id:string,scope:{paths:string[]},mode:"mutation",
  version:int,status:"active"|"expired"|"released",acquired_at:timestamp,
  expires_at:timestamp,released_at:timestamp|null,
  release_reason:"completed"|"failed"|"operator-released"|null}
Run = {id:string,issue_id:string,envelope_id:string,envelope_revision:1,
  claim_id:string,claim_version:int,actor_id:string,start_repository:Repository,
  terminal_repository:Repository|null,status:"abandoned"|"active"|"failed"|
  "succeeded",budgets:BudgetView[],started_at:timestamp,ended_at:timestamp|null}
Artifact = {id:string,digest:digest,media_type:string,size_bytes:int,
  retention_class:"evidence",availability:"available"|"missing",
  location_policy:"content-addressed",created_at:timestamp,expires_at:timestamp|null}
Receipt = {id:string,attempt_id:string,registration_id:string,
  registration_version:int,obligation_id:string,run_id:string,actor_id:string,
  trust_source:"fid-executed",outcome:"failed"|"signaled"|"spawn-failed"|
  "succeeded"|"timed-out",exit_code:int|null,signal:string|null,timed_out:boolean,
  started_at:timestamp,ended_at:timestamp,terminal_repository:Repository,
  toolchain:{executable:string,goos:string,goarch:string},output_excerpt:string,
  output_digest:digest,output_bytes:int,output_truncated:boolean,
  redaction_applied:boolean,artifacts:Artifact[] sorted by id,
  rejected_artifacts:ArtifactRejection[] sorted by path then reason}
Evidence = {id:string,obligation_id:string,receipt_id:string,
  evidence_type:"receipt",producer:"fid",trust_source:"fid-executed",
  target_repository:Repository,validity:"invalid"|"valid",
  required_artifact_ids:string[] sorted,created_at:timestamp}
ReasonCode = "artifact-unavailable"|"authority-unavailable"|"claim-expired"|
  "claim-missing"|"envelope-inactive"|"evidence-invalid"|"evidence-stale"|
  "issue-not-terminal-success"|"issue-revision-conflict"|
  "obligation-evidence-absent"|"receipt-unsuccessful"|"run-inactive"|
  "workflow-transition-invalid"
ReasonSubjectType = "actor"|"artifact"|"claim"|"envelope"|"evidence"|
  "issue"|"obligation"|"receipt"|"run"|"workflow-transition"
Reason = {code:ReasonCode,subject_type:ReasonSubjectType,subject_id:string,
  detail:string}
Evaluation = {obligation_id:string,status:"satisfied"|"unsatisfied",
  validity:"invalid"|"valid",freshness:"fresh"|"stale",
  artifact_availability:"available"|"missing",reasons:Reason[]}
Resume = {envelope:Envelope,issue:{id:string,status:string,status_category:string,
  revision:int},claim:Claim|null,run:Run|null,repository:Repository,
  obligations:Evaluation[],completion:{status:"complete"|"incomplete",reasons:Reason[]},
  next_actions:string[]}
ExportSummary = {archive_path:string,archive_digest:digest,manifest_schema:
  "fid.export-manifest.v1",project_id:string,event_count:int,terminal_event_hash:digest,
  projection_schema:14,entry_count:int}
ExportValidation = {archive_path:string,archive_digest:digest,status:"valid",
  manifest_schema:"fid.export-manifest.v1",entry_count:int}
ArtifactRejection = {path:string,reason:"changed-during-copy"|"directory"|
  "missing"|"non-regular"|"outside-cwd"|"oversized"|"symlink"}
ReadyTask = {id:string,title:string,status:string,epic_id:string|null,
  epic_title:string|null,sequence:int|null,priority:int|null}
ActivityEvent = "ac-added"|"ac-verified"|"claim-released"|"claim-renewed"|
  "claim-superseded"|"comment-added"|"concern-addressed"|"concern-raised"|
  "gate-recorded"|"issue-created"|"issue-label-added"|
  "issue-label-removed"|"issue-updated"|"link-added"|"link-removed"|
  "status-changed"|"task-completed"
RequirementRef = {id:string,type:"prd"|"trd",summary:string,
  source_doc:string,source_section:string|null,last_synced_hash:digest,
  status:"active"}
LinkRelation = "blocked-by"|"blocks"|"caused-by"|"causes"|"discovered"|
  "discovered-from"|"duplicated-by"|"duplicates"|"implemented-by"|
  "implements"|"relates-to"|"superseded-by"|"supersedes"
Finding = {check:string,message:string,entity_id:string}
MaintenanceResult = {status:"fail"|"pass"|"warn",errors:Finding[],
  warnings:Finding[],info:{issue_count:int,claim_count_live:int,
  recent_activity_7d:int,event_count:int,attempt_count:int,receipt_count:int,
  artifact_count:int,frozen_row_counts:{task_completions:int,agent_runs:int,
  ac_verifications:int,concerns:int}}}
TokenStatsRow = {actor_attribution:string|null,model:string|null,
  epic_id:string|null,total_tokens_in:int,total_tokens_out:int,run_count:int}
CostStatsRow = TokenStatsRow + {total_cost_usd:number}
```

Capability rows sort by `capability_id`. Claims sort by `issue_id,id`; Runs by
`started_at,id`; Receipts by `started_at,id`; workflows by `name,id`; Issues,
labels, links, ACs, findings, and retained result types keep Phase 20 section
3.7 order. Activity rows sort by `occurred_at,id`; Concern rows by ID;
requirements by ID; stats by the selected group key with null last. Reasons
use `contracts.md` section 10 precedence then
`subject_type,subject_id,detail` bytewise. `next_actions` is a subsequence of
`evaluate-obligations`, `run-gate`, `link-evidence`, `transition-issue`,
`finish-run`, `release-claim`, `complete-envelope`, in that order.

The closed link-relation registry is `LinkRelation`; Phase 20 snake-case
spellings are invalid input, and every inherited `Link.relation` or
`LinkListItem.relation` field has type `LinkRelation`. Maintenance error checks
sort bytewise and are
exactly `activity-journal-events`, `activity-refs`, `aggregate-sequence`,
`artifact-integrity`, `attempt-receipt-integrity`, `claim-refs`,
`claim-scope`, `dependent-row-kind`, `event-binding`, `event-sequence`,
`foreign-key-integrity`, `kind-hierarchy`, `link-refs`,
`link-relation-compatibility`, `link-scope`, `removed-state-residue`,
`sequence-scheduling`, `workflow-assignment`, and `workflow-state`.
Maintenance warning checks are exactly `frozen-state-present`,
`stale-claims`, `stale-verifications`, `started-without-claim`, and
`telemetry-gaps`. Frozen state is reported, never treated as authority or
repaired.

`maintenance.check` emits success only for `pass`, or `warn` without
`--warn-as-error`. A `fail` result, or promoted `warn`, is the exact
`projection-divergence` failure whose details are the same `MaintenanceResult`;
stdout is empty and the integrity exit is 6.

## 4. Exact Accepted Leaves

Every row permits only the global flags plus the listed leaf arguments/flags.
`I` means the named section 2.2 input. Result notation is the complete
`result` object, not a placeholder. All Event-producing rows add EventRef at
the envelope level.

### 4.1 Retained standalone and inventory operations

| Operation / command | Exact args and flags | Input | Exact result | Event/idempotency | Human success obligation |
|---|---|---|---|---|---|
| `project.init` / `fid init` | no args; optional `--prefix <^[A-Z][A-Z0-9]{0,15}$>` | none | `{project:{id:string,issue_prefix:string,root:string},journal_path:string}` | `project.initialized`; no key; replay is not supported | `initialized fid project in .fid/` |
| `project.rebuild` / `fid rebuild` | no args/leaf flags | none | `{events_replayed:int,files_replayed:int,upgraded_events:0,schema_version:14}` | null; projection-only | names replay count and completion |
| `issue.create` / `fid issue create` | no args; required `--kind <epic|subtask|task>`, `--title`, `--actor`, key; optional `--epic` (task only), `--parent` (subtask only), `--priority`, `--sequence`, `--spec-md`, repeatable `--label`, `--scheduling <prioritised|sequential>` (epic only); parent selectors exclusive | none | `{issue:{id:string,kind:string,title:string,status:string,revision:1},idempotency_key:string,replayed:boolean}` | `issue.created` | names Issue ID/status |
| `issue.get` / `fid issue get <id>` | one ID | none | `{issue:IssueDetail}` from Phase 20 §3.1 | null | stable Issue detail |
| `issue.list` / `fid issue list` | no args; optional `--kind`, `--status`, `--status-category`, `--epic`, `--label`, `--priority`, `--has-concerns`, `--order <created|priority|sequence|updated>` | none | `{issues:IssueSummary[],count:int}` | null | one deterministic row per Issue |
| `issue.update` / `fid issue update <id>` | one ID; required `--actor`, key, `--expected-revision`; at least one of `--title`, `--priority`, `--labels`, `--spec-md` | none | `{issue:{id:string,revision:int},changes:{field:string,from:string|int|string[]|null,to:string|int|string[]|null}[],idempotency_key:string,replayed:boolean}` | `issue.updated` | names Issue/revision/changed fields |
| `issue.transition` / `fid issue transition <id>` | required `--to`, `--actor`, key, `--expected-revision` | none | `{issue:{id:string,revision:int},from_state:string,to_state:string,idempotency_key:string,replayed:boolean}` | `issue.transitioned` | names from/to and revision |
| `issue.ac.add` / `fid issue ac add <issue-id>` | one ID; required `--text`, `--ordinal`, `--actor`, key | none | `{acceptance_criterion:{id:string,issue_id:string,ordinal:int,text:string},issue_revision:int,idempotency_key:string,replayed:boolean}` | `issue.ac-added` | names Issue/AC ordinal |
| `issue.ac.list` / `fid issue ac list <issue-id>` | one ID | none | `{issue_id:string,acceptance_criteria:AcceptanceCriterion[]}` | null | ordinal-ordered list |
| `issue.label.add` / `fid issue label add <issue-id> <label>` | two args; required `--actor`, key | none | `{issue_id:string,label:string,issue_revision:int,idempotency_key:string,replayed:boolean}` | `issue.label-added` | names Issue and added label |
| `issue.label.remove` / `fid issue label remove <issue-id> <label>` | two args; required `--actor`, key | none | `{issue_id:string,label:string,issue_revision:int,idempotency_key:string,replayed:boolean}` | `issue.label-removed` | names Issue and removed label |
| `label.add` / `fid label add <name>` | optional `--colour`, `--description`; required `--actor`, key | none | `{label:Label,idempotency_key:string,replayed:boolean}` | `label.added` | names label |
| `label.list` / `fid label list` | no args/leaf flags | none | `{labels:Label[],count:int}` | null | name-ordered list |
| `link.add` / `fid link add <from-id> <LinkRelation> <to-id>` | three args; required `--from-type <issue|requirement>`, `--to-type <issue|requirement>`, `--actor`, key | none | `{link:Link,issue_revisions:{id:string,revision:int}[],idempotency_key:string,replayed:boolean}` | `link.added` | names typed endpoints and added relation |
| `link.remove` / `fid link remove <from-id> <LinkRelation> <to-id>` | three args; required `--from-type <issue|requirement>`, `--to-type <issue|requirement>`, `--actor`, key | none | `{link:Link,issue_revisions:{id:string,revision:int}[],idempotency_key:string,replayed:boolean}` | `link.removed` | names typed endpoints and removed relation |
| `link.list` / `fid link list <id>` | one ID; required `--type <issue|requirement>` | none | `{entity_type:"issue"|"requirement",entity_id:string,links:LinkListItem[],count:int}` | null | ordered link list |
| `ready.get` / `fid ready` | no args; optional `--epic <issue-id>` | none | `{task:ReadyTask|null,explanation:{code:"no-ready-tasks"|"sequential-current-not-actionable",scope:"epic"|"project",epic_id:string|null,blocking_task_id:string|null,blocking_task_title:string|null,blocking_task_status:string|null}|null}`; exactly one member non-null | null | names ready task or exact explanation |
| `workflow.define` / `fid workflow define` | no args; required `--name <^[a-z][a-z0-9-]{0,63}$>`, `--applies-to <any|epic|subtask|task>`, `--states <csv>`, `--categories <csv>`, `--transitions <csv>`, `--actor`, key; 1..64 unique trimmed opaque state IDs, equal-length categories from `cancelled|completed|started|unstarted`, and 1..4096 unique `from->to` edges whose endpoints are declared; `--require` forbidden | none | `{workflow:{id:string,name:string,applies_to:string,state_count:int,transition_count:int},idempotency_key:string,replayed:boolean}` | `workflow.defined` | names workflow/counts |
| `workflow.assign` / `fid workflow assign` | no args; required `--name`, `--actor`, key and exactly one of `--kind <kind>` or `--project-wide` | none | `{assignment:{workflow_id:string,scope_type:"kind"|"project",scope_value:string|null},idempotency_key:string,replayed:boolean}` | `workflow.assigned` | names workflow/scope |
| `workflow.list` / `fid workflow list` | no args; optional `--kind`, `--active` | none | `{workflows:WorkflowSummary[],count:int}` | null | name-ordered list |
| `requirement.sync` / `fid requirement sync` | required project-relative `--source-doc`; optional `--type <prd|trd>`, `--id <requirement-id>`; ID must match selected type when both supplied; required `--actor`, key | none | `{source_doc:string,requirements:RequirementRef[],counts:{created:int,updated:int,unchanged:int},idempotency_key:string,replayed:boolean}` | `requirement.synced`; requirements sort by ID; unchanged sync still records the keyed result Event | names source/counts/IDs |
| `activity.list` / `fid activity list` | optional `--ref`, `--event <ActivityEvent>`, `--since`, `--limit <1..10000>` | none | `{items:{id:string,event:ActivityEvent,ref_type:string,ref_id:string,actor_attribution:string|null,occurred_at:timestamp}[]}` | null | chronological rows; frozen snake-case rows normalize to the matching kebab value and internal details are not exposed |
| `concern.list` / `fid concern list` | optional `--issue`, `--status <addressed|open>` | none | `{items:{id:string,issue_id:string,status:"addressed"|"open",summary:string,raised_by_attribution:string,addressed_by_attribution:string|null}[]}` | null | ID-ordered read-only legacy judgment inventory |
| `stats.tokens` / `fid stats tokens` | no args; optional `--group-by <actor|epic|model>`, `--epic <id>` | none | `{items:TokenStatsRow[]}` | null | group-key-ordered token telemetry |
| `stats.cost` / `fid stats cost` | no args; optional `--group-by <actor|epic|model>`, `--epic <id>` | none | `{items:CostStatsRow[]}` | null | group-key-ordered cost telemetry |
| `maintenance.check` / `fid maintenance check` | optional `--warn-as-error` | none | `MaintenanceResult` | null | status/counts and ordered findings |

Phase 20 shared types and closed flag domains remain normative except the
explicit Phase 21 changes above. Workflow-state strings remain opaque
identifiers; closed categories and reasons are kebab-case.

### 4.2 Operational-kernel operations

| Operation / command | Exact args and flags | Input | Exact result | Event/idempotency | Human success obligation |
|---|---|---|---|---|---|
| `capabilities.get` / `fid capabilities` | no args/leaf flags | none | `{items:Capability[]}` | null | capability, availability, limits, trust source per row |
| `activation.prepare` / `fid activation prepare` | required `--input`, key; optional `--actor <bootstrap-attribution>` | `ActivationPrepareInput` | `{activation:Activation,idempotency_key:string,replayed:boolean}` | `activation.prepared` | names activation/status; says attribution is untrusted |
| `activation.component` / `fid activation component` | required `--input`, key; optional bootstrap `--actor` | `ActivationComponentInput` | `{activation:Activation,component:{type:"authority-profile"|"gate-definitions",id:string,version:int,digest:digest},idempotency_key:string,replayed:boolean}` | `activation.component-staged` | names component/digest/status |
| `activation.commit` / `fid activation commit` | required `--input`, key; optional bootstrap `--actor` | `ActivationCommitInput` | `{activation:Activation,idempotency_key:string,replayed:boolean}` | `activation.committed` | names committed activation/digests |
| `activation.abort` / `fid activation abort` | required `--input`, key; optional bootstrap `--actor` | `ActivationAbortInput` | `{activation:Activation,idempotency_key:string,replayed:boolean}` | `activation.aborted` | names aborted activation/reason |
| `activation.finalize` / `fid activation finalize` | required `--input`, key; optional bootstrap `--actor` | `ActivationFinalizeInput` | `{activation:Activation,idempotency_key:string,replayed:boolean}` | `activation.finalized` | names finalized activation/status |
| `activation.status` / `fid activation status <activation-id>` | one activation ID; no leaf flags | none | `{activation:Activation}` | null | lifecycle/status/components |
| `gate.register` / `fid gate register` | required `--input`, `--actor`, key | `GateRegisterInput` | `{registration:Registration,idempotency_key:string,replayed:boolean}` | `gate.registered`; same definition/binding replays existing registration | names registration/version/digests |
| `gate.resolve` / `fid gate resolve <registration-id>` | one registration ID; required `--version <int>=1` | none | `{registration:Registration}` | null | names resolved immutable policy |
| `gate.run` / `fid gate run` | required `--input`, `--actor`, key | `GateRunInput` | terminal data section 8 | first run commits `gate.execution-started` then `gate.receipt-recorded`; a success response's outer `event` is the receipt-recorded EventRef; an unsuccessful terminal response is a failure envelope with no Event field; replay returns the original terminal EventRef only on success, emits one terminal frame with `replayed:true`, and starts no process | labeled streams then terminal summary |
| `envelope.create` / `fid envelope create` | required `--input`, `--actor`, key | `EnvelopeCreateInput` | `{envelope:Envelope,idempotency_key:string,replayed:boolean}` | `envelope.created` | names Envelope/Issue/Obligation |
| `envelope.get` / `fid envelope get <envelope-id>` | one Envelope ID; no leaf flags | none | `{envelope:Envelope}` | null | bounded complete Envelope view |
| `envelope.activate` / `fid envelope activate` | required `--input`, `--actor`, key | `EnvelopeActivateInput` | `{envelope:Envelope,evaluation:Evaluation,idempotency_key:string,replayed:boolean}` | `envelope.activated`; evaluation is the current Obligation state | names activated status and all reasons |
| `envelope.complete` / `fid envelope complete` | required `--input`, `--actor`, key | `EnvelopeCompleteInput` | `{envelope:Envelope,evaluation:Evaluation,idempotency_key:string,replayed:boolean}` | `envelope.completed` | names completed status and all reasons |
| `claim.acquire` / `fid claim acquire` | required `--input`, `--actor`, key | `ClaimAcquireInput` | `{claim:Claim,idempotency_key:string,replayed:boolean}` | `claim.acquired` | names Claim/version/status/expiry |
| `claim.renew` / `fid claim renew` | required `--input`, `--actor`, key | `ClaimRenewInput` | `{claim:Claim,idempotency_key:string,replayed:boolean}` | `claim.renewed` | names Claim/new version/new expiry |
| `claim.release` / `fid claim release` | required `--input`, `--actor`, key | `ClaimReleaseInput` | `{claim:Claim,idempotency_key:string,replayed:boolean}` | `claim.released` | names Claim/new version/reason |
| `claim.list` / `fid claim list` | no args; optional `--issue`, `--envelope`, `--status <active|expired|released>`, `--limit <1..10000>` | none | `{items:Claim[]}` | null | ordered Claims |
| `run.start` / `fid run start` | required `--input`, `--actor`, key | `RunStartInput` | `{run:Run,idempotency_key:string,replayed:boolean}` | `run.started` | names Run/active status/start repository |
| `run.finish` / `fid run finish` | required `--input`, `--actor`, key | `RunFinishInput` | `{run:Run,idempotency_key:string,replayed:boolean}` | `run.finished` | names Run/terminal status/repository |
| `run.get` / `fid run get <run-id>` | one Run ID; no leaf flags | none | `{run:Run}` | null | bounded Run view |
| `run.list` / `fid run list` | no args; optional `--issue`, `--envelope`, `--status <abandoned|active|failed|succeeded>`, `--limit <1..10000>` | none | `{items:Run[]}` | null | ordered Runs |
| `receipt.get` / `fid receipt get <receipt-id>` | one Receipt ID; no leaf flags | none | `{receipt:Receipt}` | null | terminal outcome/bindings/artifacts |
| `receipt.list` / `fid receipt list` | no args; optional `--run`, `--obligation`, `--outcome <failed|signaled|spawn-failed|succeeded|timed-out>`, `--limit <1..10000>` | none | `{items:Receipt[]}` | null | ordered Receipts |
| `evidence.link` / `fid evidence link` | required `--input`, `--actor`, key | `EvidenceLinkInput` | `{evidence:Evidence,evaluation:Evaluation,idempotency_key:string,replayed:boolean}` | `evidence.linked` | names evidence and current status/reasons |
| `evidence.get` / `fid evidence get <evidence-id>` | one Evidence ID; no leaf flags | none | `{evidence:Evidence,evaluation:Evaluation}` | null | evidence dimensions/reasons |
| `evidence.evaluate` / `fid evidence evaluate` | required `--input`; no Actor/key (pure evaluation) | `EvidenceEvaluateInput` | `{evaluation:Evaluation}` | null | status plus every reason |
| `artifact.verify` / `fid artifact verify` | required `--input`, `--actor`, key | `ArtifactVerifyInput` | `{artifact:Artifact,changed:boolean,idempotency_key:string,replayed:boolean}` | `artifact.availability-observed`; unchanged still records idempotent observation Event | availability/digest/change |
| `projection.resume` / `fid resume <issue-id>` | one Issue ID; optional `--repository <path>` default current worktree | none | `{resume:Resume}` | null | current state, reasons, next actions; no recovery claim |
| `export.create` / `fid export create <archive-path>` | one archive path; optional `--repository <path>` default current worktree | none | `{export:ExportSummary}` | null; no project mutation/idempotency | path/digest/counts |
| `export.validate` / `fid export validate <archive-path>` | one archive path; no leaf flags | none | `{validation:ExportValidation}` | null | valid schema/digest/counts |

## 5. Fail-Closed And Absent Leaves

The following final Phase 21 leaves remain registered only to give typed
failure and then perform no read/write: `fid issue ac verify`, `fid concern
raise`, `fid concern address`, and `fid maintenance vacuum`. They retain their
Phase 20 positional/flag grammar except `--by` is rejected and any actor flag
is `--actor`; every invocation returns `capability-unavailable`, exit 5,
subject type `capability`, and details
`{capability_id:<human.outcomes.v1|reviews.coverage.v1|telemetry.trajectory.v1>,available:false,owner_phase:<23|25>}`.

The following are absent (Cobra unknown-command failure, exit 2, Human stderr;
JSON callers receive `operation-unavailable`): `fid policy`, `fid dev`, legacy
`fid issue claim|unclaim`, legacy `fid run record`, `fid issue complete`,
`activation recover`, `envelope revise|cancel`, `run recover`, Stop, Attention,
Receipt ingest, Review, Human Outcome, telemetry ingestion, `export import`,
and every V1 `--by` mutation spelling. `fid journal inspect` and `fid sql`
remain explicitly unstable Human-only local diagnostics outside `fid-cli.v1`;
they grant no authority, reject `--json`, and are not exported.

## 6. Idempotency Replay Representation

For every keyed mutation, the canonical request digest includes operation,
project ID, closed normalized input/flags, resolved Actor identity for kit
mutations or literal attribution for standalone mutations, and expected
version. It excludes output mode and key itself.

- First commit returns the operation result with `idempotency_key:<key>` and
  `replayed:false`; the outer EventRef names the committed Event.
- Same key and digest returns byte-identical domain fields and original
  EventRef, changes only `replayed:true`, emits no Event, and performs no
  process or filesystem side effect.
- Same key and different digest returns `idempotency-conflict` with details
  `{operation_scope:string,idempotency_key:string,stored_request_digest:digest,
  supplied_request_digest:digest}` and no Event/effect.
- `gate.run` stores the terminal idempotency record for successful and observed
  non-success process outcomes. A non-success first run and replay both carry
  the ordinary failure envelope inside the terminal frame, so they have no
  outer Event/result; `receipt_id` binds the committed Receipt, and terminal
  `replayed` changes from false to true. The replay starts no process and emits
  no Event. A changed digest still fails before Attempt append.
- Collections never contain replay metadata. `project.init`, reads, rebuild,
  evaluate, resume, and export create/validate have no idempotency key.

## 7. Complete Failure And Exit Contract

Exit classes are: `0` success; `2` invalid input/unsupported operation; `3`
not found; `4` conflict/precondition; `5` authority/capability; `6` integrity;
`7` internal; `8` observed gate non-success. Category maps respectively to
`invalid-input`, `not-found`, `conflict|precondition`, `authority`,
`integrity`, `internal`; `gate-unsuccessful` alone uses precondition/8.

The following table is closed. `subject` gives exact subject type; its ID is
the referenced value, or `null` with a null subject where stated. Detail
objects are closed. All values are public kebab-case where enumerated.
For a `*-not-found` group, subject type is the code prefix and subject ID is
the requested ID. `duplicate-entity` uses the supplied `entity_type`.
`authority-missing` has a null subject; `authority-unknown-actor` and
`authority-denied` use the supplied Actor; `authority-unavailable` uses the
resolved Actor when one exists and otherwise null. `evidence-invalid` and
`evidence-stale` use Evidence; `artifact-unavailable` uses Artifact.

| Code | Category / retryable / exit | Subject | Exact details |
|---|---|---|---|
| `invalid-argument`, `unsupported-kind`, `unsupported-relation`, `operation-unavailable` | invalid-input / false / 2 | null | `{argument:string,value:string|null,allowed:string[]}` |
| `project-not-found` | not-found / false / 3 | project | `{path:string}` |
| `project-already-initialized` | conflict / false / 4 | project | `{path:string}` |
| `issue-not-found`, `label-not-found`, `workflow-not-found`, `requirement-not-found` | not-found / false / 3 | named type | `{id:string}` |
| `link-not-found` | not-found / false / 3 | link | `{from_type:string,from_id:string,relation:string,to_type:string,to_id:string}` |
| `activation-not-found`, `registration-not-found`, `envelope-not-found`, `claim-not-found`, `run-not-found`, `receipt-not-found`, `evidence-not-found`, `artifact-not-found` | not-found / false / 3 | named type | `{id:string}` |
| `duplicate-entity` | conflict / false / 4 | entity type | `{entity_type:string,identity:{fields:string[],values:string[]}}` |
| `state-conflict` | conflict / false / 4 | issue | `{expected_revision:int,actual_revision:int}` |
| `idempotency-conflict` | conflict / false / 4 | idempotency | exact section 6 object |
| `lock-unavailable` | conflict / true / 4 | null | `{lock_path:string}` |
| `activation-conflict` | conflict / false / 4 | activation | `{active_id:string|null,nonterminal_id:string|null,requested_direction:"activate"}` |
| `activation-invalid` | precondition / false / 4 | activation | `{status:string,expected_statuses:("aborted"|"committed"|"finalized"|"prepared"|"staged")[] sorted-unique,component_errors:{component_id:string,code:"component-duplicate"|"component-missing"|"digest-mismatch"|"schema-invalid"|"status-invalid"|"version-conflict"}[]}` sorted by component ID/code |
| `activation-recovery-unavailable` | authority / false / 5 | activation | `{status:string,point_of_no_return:true,owner_phase:22}` |
| `hierarchy-invalid` | precondition / false / 4 | issue | `{kind:string,parent_id:string|null,parent_kind:string|null,rule:string}` |
| `dependency-cycle` | precondition / false / 4 | issue | `{path:string[]}` |
| `transition-not-permitted` | precondition / false / 4 | issue | `{from_state:string,to_state:string,reasons:Reason[]}` |
| `authority-missing`, `authority-unknown-actor`, `authority-denied`, `authority-unavailable` | authority / false / 5 | actor or null | `{operation:string,actor_id:string|null,required_actions:string[],required_capabilities:string[]}` |
| `capability-unavailable` | authority / false / 5 | capability | `{capability_id:string,available:false,owner_phase:int}` |
| `registration-mismatch` | precondition / false / 4 | gate-registration | `{registration_id:string,expected_version:int,actual_version:int|null,mismatches:("activation"|"definition-digest"|"project"|"repository"|"scope"|"version"|"worktree")[] sorted-unique}` |
| `repository-unavailable` | precondition / true / 4 | repository | `{missing_fields:string[] sorted-unique,project_id:string|null,worktree_id:string|null}` |
| `executor-policy-unsupported` | invalid-input / false / 2 | gate-definition | `{fields:string[],supported:{sandbox:"none",network:"allowed",write:"allowed",shell:false}}` |
| `unsafe-path` | invalid-input / false / 2 | path | `{path:string,reason:"absolute"|"outside-project"|"symlink"|"traversal"}` |
| `envelope-inactive` | precondition / false / 4 | envelope | `{status:string,expected_revision:1,actual_revision:int}` |
| `claim-conflict`, `claim-expired`, `claim-stale` | conflict / false / 4 | claim | `{issue_id:string,claim_id:string|null,expected_version:int|null,actual_version:int|null,expires_at:timestamp|null}` |
| `run-inactive` | precondition / false / 4 | run | `{run_id:string,status:string}` |
| `credential-unavailable` | precondition / true / 4 | gate-registration | `{names:string[] sorted-unique}` |
| `obligation-unsatisfied` | precondition / false / 4 | obligation | `{reasons:Reason[]}` |
| `evidence-invalid`, `evidence-stale`, `artifact-unavailable` | precondition / false / 4 | evidence or artifact | `{reasons:Reason[]}` |
| `gate-unsuccessful` | precondition / false / 8 | receipt | `{attempt_id:string,receipt_id:string,outcome:"failed"|"signaled"|"spawn-failed"|"timed-out",exit_code:int|null,signal:string|null,timed_out:boolean}` |
| `journal-replay-stopped` | integrity / false / 6 | null | `{journal_file:string,line:int|null,event_id:string|null}` |
| `projection-divergence` | integrity / false / 6 | null | exact `maintenance.check` result |
| `result-limit-exceeded` | integrity / false / 6 | null | `{collection:string,limit:int,observed:int}` |
| `archive-invalid` | integrity / false / 6 | export | `{archive_path:string,checks:{code:"artifact-reference-invalid"|"capability-snapshot-mismatch"|"digest-mismatch"|"duplicate-entry"|"entry-extra"|"entry-missing"|"journal-invalid"|"manifest-invalid"|"path-unsafe"|"size-mismatch"|"terminal-hash-mismatch",entry:string|null,detail:string}[]}` sorted by code/entry/detail |
| `internal-error` | internal / false / 7 | null | `{}` |

Messages are concise Human explanations and are never control flow. Details
never include credentials, raw environment, unsanitized output, Artifact
bytes, SQL, stack traces, or unrelated absolute paths.

## 8. Exact `fid.live.v1` Frames

Each stdout line is one closed object:

```text
LiveFrame = {frame_schema:"fid.live.v1",execution_id:string,
  sequence:int>=1,stream:"control"|"stderr"|"stdout"|"terminal",data:object}
```

Sequence starts at 1 and increments across every stream. The payload union is:

```text
stdout|stderr data = {encoding:"utf-8",text:string[1..4096 UTF-8 bytes]}
control output-truncated = {kind:"output-truncated",limit_bytes:int,
  observed_bytes_lower_bound:int}
control redaction-applied = {kind:"redaction-applied",replacements:int>=1}
control artifact-rejected = {kind:"artifact-rejected",path:string,
  reason:"changed-during-copy"|"directory"|"missing"|"non-regular"|
  "outside-cwd"|"oversized"|"symlink"}
terminal data = {attempt_id:string,receipt_id:string|null,outcome:"failed"|
  "fid-error"|"signaled"|"spawn-failed"|"succeeded"|"timed-out",
  exit_code:int|null,signal:string|null,
  timed_out:boolean,output_bytes:int,output_truncated:boolean,
  redaction_applied:boolean,artifact_ids:string[] sorted,
  rejected_artifacts:ArtifactRejection[] sorted,response:FidCliEnvelope,
  replayed:boolean}
```

Invalid process bytes become U+FFFD. Sanitized chunks never split a UTF-8 code
point. CR, LF, ESC, and other controls are JSON-escaped. Credential and
registered-literal replacement happens before framing, counting, excerpt,
hash, Receipt, Artifact, error, export, or Human display. A redaction control
frame reports only a count. Exactly one `redaction-applied` frame is emitted
after process streams close when the execution made one or more source-value
replacements; its count is the number of source occurrences, not the number of
downstream surfaces. Exactly one truncation frame is emitted at the first
over-limit observation; fid then drains both pipes but emits/retains no more
process data. Artifact candidates are processed by normalized declared path;
one `artifact-rejected` frame is emitted per rejection in path/reason order.
After process data closes, the remaining control order is redaction if
applicable, Artifact rejections, then terminal.

Exactly one terminal frame is last. It follows terminal Receipt commit for
every process outcome. For `succeeded`, `response` is a success envelope with result
`{receipt:Receipt,idempotency_key:string,replayed:boolean}` and process exit 0.
For every other outcome it is `gate-unsuccessful`, process exit 8, and exact
details from section 7, except handled fid infrastructure failure after the
first frame uses `fid-error`, null Receipt/exit/signal, `internal-error`, exit
7, and no idempotency result. `spawn-failed` requires null exit/signal;
`timed-out` requires `timed_out:true`; `signaled` requires a non-null signal;
`failed` requires non-success exit and null signal. Output truncation does not
change outcome. No JSON follows the terminal frame.

Human gate mode renders sanitized chunks as `[stdout] <text>` and `[stderr]
<text>`, renders controls as one concise typed notice without raw bytes, and
ends with Receipt ID or explicit absence, outcome, exit/signal/timeout,
Artifact IDs/rejections, and Obligation status. It never claims OS containment or evidence satisfaction
from an Attempt without a qualifying Receipt.
