# Phase 22 Public Interface Amendment

> PROPOSED with DEC-180. Read with `contracts.md` and `schema.md`.
> Phase 21 `public-interface.md` remains the baseline for unmodified leaves.

## 1. Transport And Shared Grammar

Retain `fid-cli.v1`, `--json`, `--project-root`, existing error envelope, and
exit mapping. Added mutations take no positional arguments and require
`--input <file|->`, `--actor <activated-actor-id>`, and `--idempotency-key`.
Activation recovery uses the existing optional bootstrap-attribution `--actor`.
No mutation accepts a caller success flag, executor argv, claimed capability,
SQL, or raw log. All request objects are closed recursively; unknown/duplicate
keys, missing required fields, trailing JSON, and wrong nullability reject.

Notation is equivalent closed schema: `?T` is required nullable T; `T[]` is a
required array, including `[]`. ID is nonempty printable ASCII up to 128 bytes;
text is valid UTF-8 up to 4096 bytes; timestamp/digest/Repository/Subject/Actor
retain Phase 21 grammar. Arrays are bounded at 10000 elements unless stated.
Paths are relative normalized UTF-8, no `..`, NUL, absolute path, or escaping
symlink; inability to represent a path fails inspection. Version is int >=1.
All result arrays sort by ID except paths by path, budgets by name, resources
by component, and ordered Events by journal order. List limits default 100 and
accept 1..10000; exceeding a bounded complete projection returns
`result-limit-exceeded`, never silently incomplete recovery information.

`Mutation<T>` means `{record:T,idempotency_key:ID,replayed:boolean}` below;
the actual key replacing `record` is explicitly named for each leaf. Success
includes the retained outer EventRef; reads have `event:null`. Same-key replay
returns the stored result with only `replayed:true` changed. Rejections have
no effects except the explicitly committed execution denial described below.

## 2. Closed Values

| Value | Exact fields and domains |
|---|---|
| Executor | `{id:ID,source:"activated-harness"}`; bound once at Run start |
| Run (replace) | Phase 21 Run plus `{version:int,executor:Executor,obligation_id:?ID,predecessor_run_id:?ID,successor_run_id:?ID}`; status adds `interrupted`; Budgets are current derived views, not reset copies |
| Claim (amend) | Phase 21 Claim; status adds `superseded`; release reason adds `interrupted`, `recovery-superseded`, `abandoned`; reads include `{effective_status:"active"\|"expired"\|"released"\|"superseded"}` |
| Liveness | `{executor_id:ID,state:"live"\|"terminated"\|"unknown",source_ref:text,observed_at:timestamp}` |
| PathChange | `{path:text,before_digest:?digest,after_digest:?digest,classification:"in-scope"\|"outside-scope"\|"protected"}`; add/delete uses one null digest |
| Reconciliation | `{repository:Repository,checkpoint_id:?ID,paths:PathChange[],unmerged_paths:text[],operation:"none"\|"merge"\|"rebase"\|"cherry-pick"\|"revert"\|"bisect",ancestry:"unchanged"\|"descendant"\|"conflict"\|"unknown",digest:digest}` |
| RecoveryAuthority | `{mode:"ordinary"\|"force",actor_id:ID,reason:text,reference:?text,grant_ref:?{envelope_id:ID,revision:1,ordinal:int}}`; force requires reference, exact current grant_ref, and trusted Human; ordinary grant_ref is null |
| Recovery | `{id:ID,prior_run_id:ID,successor_run_id:ID,prior_claim_id:ID,successor_claim_id:ID,observation_id:?ID,reconciliation:Reconciliation,adopted_paths:text[],authority:RecoveryAuthority,event_id:ID}` |
| Condition | `{type:"executor-terminated"\|"run-recovered"\|"activation-recovered"\|"capability-available"\|"budget-observed"\|"work-ended",subject:Subject,budget_name:?ID,capability_id:?ID}`; only budget/capability condition uses its matching non-null field |
| Stop | `{id:ID,subject:Subject,code:StopCode,blocked_actions:text[],details:StopDetails,resolution_condition:Condition,resolver_capability:ID,actor_id:ID,raised_at:timestamp,resolution:?StopResolution}` |
| StopResolution | `{actor_id:ID,fact_event_id:ID,resolved_at:timestamp,event_id:ID}` |
| Attention | `{id:ID,semantic_key:digest,subject:Subject,code:StopCode,status:"open"\|"resolved",version:int,stop_ids:ID[],raised_at:timestamp,resolution_event_id:?ID}` |
| BudgetSample | `{id:ID,envelope_id:ID,envelope_revision:1,budget_name:ID,run_id:ID,actor_id:ID,source:"fid"\|"activated-harness",source_sample_id:ID,sample_kind:"delta"\|"baseline",previous_sample_id:?ID,value:?int,measurement_ref:?Source,observed_at:timestamp,event_id:ID}`; numeric value >=0 |
| Budget (amend) | Phase 21 fields; units and arithmetic from contracts section 6; exhaustion_action `budget-exhausted`; consumed/remaining int >=0 or null; add derived `state:"unknown"\|"available"\|"warning"\|"exhausted"` on output only |
| ActivationResource | `{component:"kit"\|"overlay"\|"provider"\|"project"\|"backend",version:int,old_identity:?text,new_identity:text,observed_identity:?text,state:"pending"\|"old"\|"new"\|"error",evidence_ref:?text,error_code:?text}` |
| Activation (amend) | Phase 21 fields plus `{version:int,prior_activation_id:?ID,resources:ActivationResource[],recovery:?ActivationRecovery}`; status adds `superseded`; original direction remains `activate` |
| ActivationRecovery | `{direction:"restore-prior"\|"complete-forward",actor_attribution:text,reason:text,selected_event_id:ID,completed_event_id:?ID}` |
| NextAction | `{operation:text,subject:Subject,required_capabilities:ID[],reasons:Reason[]}`; operation is a present public operation ID; reasons are empty only when executable |
| Resume (amend) | Phase 21 fields; replace next_actions with NextAction[]; add `{stops:Stop[],attention:Attention[],budgets:Budget[],recoveries:Recovery[],findings:{available:false,items:[]}}`; repository is observed current target |

StopCode is exactly the seven-code set in `contracts.md` section 5. Subject
types for Stop/Attention are only run, envelope, activation. Blocked actions
are a sorted unique nonempty subset of `claim.acquire`, `run.start`,
`run.recover`, `gate.run`, `issue.transition-terminal-success`,
`envelope.complete`; a budget Stop uses `run.start`, `run.recover`, `gate.run`.

StopDetails is selected by code, with no additional properties:

| Code | Details |
|---|---|
| executor-liveness-unknown | `{run_id:ID,executor_id:ID,observation_id:?ID}` |
| recovery-state-conflict | `{run_id:ID,observation_digest:digest,paths:text[],reasons:text[]}` |
| budget-exhausted / budget-data-unavailable | `{envelope_id:ID,budget_name:ID,unit:text,maximum:int,consumed:?int,remaining:?int}` |
| activation-incomplete | `{activation_id:ID,point_of_no_return:boolean,components:text[]}` |
| capability-unavailable | `{capability_id:ID,available:false,owner_phase:int}` |
| operator-stop | `{reason:text,reference:text}` |

## 3. New And Replaced Leaves

Each input lists every required key. Shared flags follow section 1 unless the
row says read-only. No new flag-only mutation grammar is introduced.

| Operation / CLI | Input or read arguments | Result | Event / task |
|---|---|---|---|
| `claim.expire` / `fid claim expire` | `{claim_id:ID,expected_version:int}` | Mutation keyed `claim`: Claim | claim.expired / 3,6 |
| `run.start` / existing leaf | Phase 21 input plus `{executor:Executor,obligation_id:?ID}` | Mutation keyed `run`: Run | run.started / 3 |
| `run.observe` / `fid run observe` | `{run_id:ID,expected_version:int,liveness:Liveness}` | `{observation:{id:ID,run_id:ID,kind:"liveness",actor_id:ID,observed_at:timestamp,liveness:Liveness},run:Run,idempotency_key:ID,replayed:boolean}` | run.observed / 3,6 |
| `run.inspect` / `fid run inspect <run-id>` | read; optional `--repository <path>` default current worktree | `{run:Run,reconciliation:Reconciliation,liveness:?Liveness}` | null / 4,6 |
| `run.checkpoint` / `fid run checkpoint` | `{run_id:ID,expected_version:int,reconciliation_digest:digest}` | `{checkpoint:{id:ID,run_id:ID,kind:"checkpoint",actor_id:ID,observed_at:timestamp,reconciliation:Reconciliation},run:Run,idempotency_key:ID,replayed:boolean}` | run.checkpointed / 4,6 |
| `run.recover` / `fid run recover` | `{prior_run_id:ID,expected_run_version:int,expected_claim_version:int,observation_id:?ID,reconciliation_digest:digest,adopted_paths:text[],authority:RecoveryAuthority,executor:Executor,ttl_seconds:int}` | `{recovery:Recovery,run:Run,claim:Claim,idempotency_key:ID,replayed:boolean}` | run.recovered / 5,6 |
| `stop.raise` / `fid stop raise` | `{subject:Subject,code:StopCode,blocked_actions:text[],details:StopDetails,resolution_condition:Condition,resolver_capability:ID}` | `{stop:Stop,attention:Attention,idempotency_key:ID,replayed:boolean}` | stop.raised / 2,9 |
| `stop.resolve` / `fid stop resolve` | `{stop_id:ID,fact_event_id:ID}` | Mutation keyed `stop`: Stop | stop.resolved / 2,9 |
| `stop.list` / `fid stop list` | read; optional `--subject <run\|envelope\|activation>:<id>`, `--status <open\|resolved>`, `--limit` | `{items:Stop[]}` | null / 9 |
| `attention.list` / `fid attention list` | same read flags as stop.list | `{items:Attention[]}` | null / 9 |
| `attention.resolve` / `fid attention resolve` | `{attention_id:ID,expected_version:int,fact_event_ids:ID[]}`; exact qualifying resolutions for all linked Stops | Mutation keyed `attention`: Attention | attention.resolved / 2,9 |
| `run.observe-budget` / `fid run observe-budget` | `{run_id:ID,expected_version:int,budget_name:ID,source_sample_id:ID,sample_kind:"delta"\|"baseline",previous_sample_id:?ID,value:?int,measurement_ref:?Source,observed_at:timestamp}` | `{sample:BudgetSample,budget:Budget,stops:Stop[],attention:Attention[],idempotency_key:ID,replayed:boolean}` | run.budget-observed / 7,9 |
| `run.close` / `fid run close` | `{run_id:ID,expected_version:int,observation_id:ID,repository:Repository}` | Mutation keyed `run`: Run | run.closed / 8,9 |
| `run.abandon` / `fid run abandon` | `{run_id:ID,expected_version:int,expected_claim_version:int,observation_id:?ID,authority:RecoveryAuthority,repository:Repository}` | `{run:Run,claim:Claim,idempotency_key:ID,replayed:boolean}` | run.abandoned / 8,9 |
| `activation.prepare` / existing leaf | Phase 21 input plus `{expected_prior_activation_id:?ID,expected_prior_version:?int,resources:{component:text,old_identity:?text,new_identity:text}[]}`; exactly five distinct resource names | existing activation mutation with extended Activation | activation.prepared / 10,11 |
| `activation.component`, `activation.commit`, `activation.abort`, `activation.finalize` / existing leaves | Existing Phase 21 closed input plus `{expected_version:int}` | Existing result with amended Activation | existing activation family / 10,11 |
| `activation.resource` / `fid activation resource` | `{activation_id:ID,component:text,expected_activation_version:int,expected_version:int,observed_identity:?text,state:"old"\|"new"\|"error",evidence_ref:text,error_code:?text}` | `{activation:Activation,resource:ActivationResource,idempotency_key:ID,replayed:boolean}` | activation.resource-observed / 10,11 |
| `activation.recover` / `fid activation recover` | `{activation_id:ID,expected_version:int,expected_status:text,direction:"restore-prior"\|"complete-forward",reason:text}` | Mutation keyed `activation`: Activation | activation.recovery-selected / 10,11 |
| `activation.recover-complete` / `fid activation recover-complete` | `{activation_id:ID,expected_version:int,selected_event_id:ID,resource_versions:{component:text,version:int}[]}` | Mutation keyed `activation`: Activation | activation.recovery-completed / 10,11 |
| `export.import` / `fid export import <archive-path>` | read archive; required `--repository <existing-target-directory>`, `--idempotency-key`; no Actor/input flag | `{import:{project_id:ID,archive_digest:digest,event_count:int,terminal_event_hash:digest,projection_schema:int,target_root:text,replayed:boolean}}` | null; atomic installation / 14,15 |

Retain expiry TTL bounds 1..86400 seconds. Expected Run version increments on
every Run observation, checkpoint, budget observation, closure, or recovery.
Gate operations have their existing idempotency fence and do not mutate the
Run version solely because their own Attempt/Receipt aggregate advances.

Every activation mutation checks/increments Activation.version; resource writes
also check their own expected_version. Existing component/commit/abort/finalize
leaves require the added activation expected_version field. A version mismatch
uses activation-conflict/4 with closed details
`{activation_id:ID,expected_version:int,actual_version:int}`; competing prepare
retains the Phase 21 activation-conflict details variant. Missing predecessor
uses both expected_prior fields null; otherwise both must match current state.

Reason retains the Phase 21 structure and extends its code domain with the
applicable Stop and failure codes in this amendment. Human message text is never
used to decide resolution or action eligibility.

The Actor in RecoveryAuthority must equal `--actor`; it never grants itself
capabilities. Ordinary `run.recover` requires `run.recover`; force additionally
requires `run.force-recover`. Claims use their corresponding action capability;
checkpoint/observe/budget need matching activated harness capabilities and
Envelope grants. Stop raise/resolve and Attention resolve require their exact
action capabilities/grants. Abandon needs `run.abandon`, plus force authority
for a known-live observation; absent/unknown liveness cannot authorize it.
Force requests reference an existing live observation, never null. Ordinary
recovery requires a terminated observation. Read operations are inspection only.

`fid run finish` is removed in task 9 after task 8's derived closure passes.
Existing get/list adopt amended types and status filters; they do not omit
interrupted Runs or superseded Claims. `fid gate run` adds cancellation and
explicit durable-denial metadata without accepting terminal results as input.

## 4. Failure And Live Frames

Retain the Phase 21 complete exit table. Added codes are closed as follows;
`subject` is the referenced Run, Stop, Attention, activation, budget Envelope,
or export. All are retryable false. Transient lock errors retain retryable true.

| Code | Category / exit | Exact details |
|---|---|---|
| executor-liveness-unknown | precondition / 4 | `{run_id:ID,executor_id:ID,observation_id:?ID}` |
| executor-still-live | precondition / 4 | `{run_id:ID,executor_id:ID,observation_id:ID}` |
| recovery-state-conflict | conflict / 4 | `{run_id:ID,observation_digest:?digest,paths:text[],reasons:text[]}` |
| run-version-conflict | conflict / 4 | `{run_id:ID,expected_version:int,actual_version:int}` |
| attention-version-conflict | conflict / 4 | `{attention_id:ID,expected_version:int,actual_version:int}` |
| recovery-already-linked | conflict / 4 | `{run_id:ID,successor_run_id:ID}` |
| stop-not-found / attention-not-found / observation-not-found | not-found / 3 | `{id:ID}` |
| stop-unresolved | precondition / 4 | `{stop_ids:ID[],reasons:Reason[]}` |
| resolution-invalid | precondition / 4 | `{subject:Subject,fact_event_ids:ID[],reasons:Reason[]}` |
| budget-exhausted / budget-data-unavailable | precondition / 4 | `{budget:Budget,stop_ids:ID[],attention_ids:ID[],event:?EventRef,replayed:boolean}` |
| budget-sample-conflict | conflict / 4 | `{budget_name:ID,source_sample_id:ID,expected_previous_sample_id:?ID,actual_previous_sample_id:?ID}` |
| activation-direction-conflict | precondition / 4 | `{activation_id:ID,point_of_no_return:boolean,required_direction:"restore-prior"\|"complete-forward"}` |
| activation-incomplete | precondition / 4 | `{activation_id:ID,components:text[],reasons:Reason[]}` |
| import-target-conflict | conflict / 4 | `{target_root:text,reason:"populated"\|"different-import"\|"unsafe-path"}` |

`activation-recovery-unavailable` is removed once activation recovery exists;
post-commit abort returns `activation-direction-conflict`. Unknown schemas and
archive resource/reference failures use `archive-invalid`/6 and extend its
closed check code set with `resource-limit`, `reference-invalid`,
`effect-invalid`, `configuration-invalid`, `schema-unsupported`.
Cancellation uses `gate-unsuccessful`/8 with outcome `cancelled` and null
exit/signal where none was observed.

An execution denial that durably raises Stop/Attention still returns `ok:false`
and exit 4. Its error details contain the committed EventRef and immutable
Stop/Attention IDs. On a denial that commits nothing, EventRef is null. Store
and replay a committed denial under its original key. `fid.live.v1` emits one
terminal `fid-error` frame carrying that same error; no start/output frame or
spawn occurs. A budget crossing during an admitted execution is a control frame
`budget-stop` with `{stop_ids:ID[],attention_ids:ID[],budget:Budget}` followed by
the normal observed terminal Receipt. Every live variant has deterministic
frame-order and no-partial-success conformance cases.

## 5. Manifest V2 And Availability

Manifest v2 retains v1's exact fields and adds
`project_config:{config_version:1,project:{id:ID,issue_prefix:text},storage:{artifact_root:text},executor_limits:{max_output_bytes:int,max_excerpt_bytes:int},maintenance:{abandoned_attempt_grace_seconds:int}}`.
Its schema is `fid.export-manifest.v2`; projection_schema is the exporter's
active accepted migration (15 during development, 16 at phase end). Chunked
journal entry paths are `journal/<eight-digit ordinal>.jsonl`, starting at
00000001, containing canonical Event-v2 lines in original order. Other entry
paths/metadata retain the Phase 21 format. Import accepts only the binary's exact current accepted schema (15 during
implementation, 16 after task 16); final capability qualification requires 16.
Intermediate v2 import tests do not advertise a completed portable capability. Manifest v1 validation stays
read-only; v1 import rejects schema-unsupported. No old-journal conversion.

Final changed capability rows retain schema_min/schema_max 1:

| Capability | Available | Exact limits | Trust source |
|---|---|---|---|
| attention.v1 | true | `{"state":true,"human_outcome_resolution":false}` | fid-journal |
| kit.activation.v1 | true | `{"abort_before_commit":true,"automatic_respawn":false,"finalize":true,"recover":true}` | active-kit-activation |
| export.portable.v1 | true | `{"create":true,"import":true,"validate":true,"import_mode":"empty-target"}` | fid-journal |
| work.claims.v1 | false | `{"max_active_per_issue":1,"parallelism":false,"recover":true}` | phase-24-unavailable |

Schema presence alone does not enable these rows. Task 17 changes the final
binary matrix only after task 15 and replacement/removal proofs. Discovery,
resume, exports, and fake/fid fixtures must agree. Negative Claims capability
discovery remains mandatory; supported recovery operations can be exercised
directly without relabelling the incomplete capability.

## 6. Human And Independent Review

Human text for recovery names prior/successor IDs, adopted paths, authority, and
remaining budgets; force output states known-live takeover and residual filesystem risk. Stop
and Attention text names the blocking reason and qualifying resolution. Import
text identifies the destination and preserved project/history; activation
recovery text says restore-prior or complete-forward. Never claim that adoption
is evidence or that post-commit recovery rolled back.

Task 18's shared conformance fixtures cover every new/replaced leaf, amended
read shape, failure code, and live frame above, with negative unknown-field,
authority, stale-version, duplicate-delivery, and interrupted-publication cases.
Fixture owners add cases alongside their earlier tasks; task 18 assembles and
audits them. QA independently compares the documented port surface with the
kit's named authoritative contract before execution approval.
