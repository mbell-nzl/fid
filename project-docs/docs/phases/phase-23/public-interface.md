# Phase 23 Public Interface Amendment

> ACCEPTED under DEC-222 with `contracts.md` and `schema.md`.
> Phase 21 and Phase 22 `public-interface.md` remain the baseline for every
> leaf, value, code, and exit this document does not replace.

## 1. Transport And Shared Grammar

Retain `fid-cli.v1`, `--json`, `--project`, the existing error envelope, exit
mapping, and the Phase 22 notation (`?T`, `T[]`, ID, text, timestamp, digest,
Subject, Source, Repository, Version). Every Phase 23 mutation takes no
positional argument and requires `--input <file|->`, `--actor
<activated-actor-id>`, and `--idempotency-key`. Every Phase 23 read uses its
declared read grammar. No mutation accepts a caller success flag, a claimed
capability, prose bodies, SQL, or a raw log.

All request objects stay closed recursively: unknown keys, duplicate keys,
missing required fields, trailing JSON, and wrong nullability reject with no
Event. Arrays keep the 10000-element ceiling unless a tighter bound is stated
here. `Mutation<T>` retains the Phase 22 meaning; the key replacing `record` is
named per leaf. Lists keep default limit 100, accepted range 1..10000, and
`result-limit-exceeded` on a bounded complete projection.

New result arrays sort by ID, except coverage items and exclusions by `Source`
display label, references by display label, protected Decisions by ordinal, and
Obligation evaluations by Obligation ID.

## 2. Amended Closed Values

| Value | Exact change |
|---|---|
| SubjectType (amend) | Retain the Phase 21/22 set and add `review`, `finding`, `concern`, `deviation`, `decision`, `human-outcome`. `requirement` is retained and now names a `requirement_refs` row. |
| ObligationInput (replace) | `{id:ID,type:"gate-evidence"\|"review"\|"human-outcome",applicability:"required",guarded_actions:text[1..2],subject:{type:"issue",id:ID},registration:?{id:ID,version:int>=1},review_profile:?ID,outcome_class:?ID,accepted_evidence_types:text[1],accepted_trust_sources:text[1],freshness_target:"repository-tree"\|"none",artifact_policy:"required-available"\|"metadata-sufficient",independence:"fid-executed"\|"distinct-actor"\|"human",required_lane:null,source:Source}`; the exact per-type field combination is fixed by `contracts.md` section 2 |
| Obligation (amend) | `ObligationInput` plus `{derived_status:"satisfied"\|"unsatisfied",reasons:Reason[]}`, unchanged in shape |
| Evaluation (amend) | `artifact_availability` domain becomes `"available"\|"missing"\|"not-applicable"`; all other fields unchanged |
| Budget (amend) | unit domain adds `review-cycles`; `exhaustion_action` domain adds `review-rework-exhausted`; all other Phase 22 rules unchanged |
| Claim (amend) | `mode` domain becomes `"mutation"\|"review"`; all other fields, statuses, and release reasons unchanged |
| Envelope (amend) | `obligations:Obligation[1..16]`; add `protected_decisions:ProtectedDecision[0..16]` |
| Evidence (amend) | `{id:ID,obligation_id:ID,evidence_type:"receipt"\|"review"\|"human-outcome",receipt_id:?ID,review_id:?ID,human_outcome_id:?ID,producer:"fid"\|"activated-actor"\|"human",trust_source:"fid-executed"\|"activated-actor"\|"activated-human",target_repository:Repository,validity:"valid"\|"invalid",required_artifact_ids:ID[],created_at:timestamp}`; exactly one subject ID is non-null and matches `evidence_type` |
| Condition (amend) | type domain adds `review-recorded`, `finding-resolved`, `human-outcome-recorded`, `decision-satisfied`; value adds `outcome_class:?ID`, `decision_class:?ID`, `finding_id:?ID`; only the field matching the type is non-null |
| StopCode (amend) | adds `review-rework-exhausted`, `protected-decision-required`, `human-outcome-required`; total closed set becomes ten |
| Stop blocked actions (amend) | the allowed subset adds nothing; `review-rework-exhausted` uses exactly `claim.acquire`,`run.start`; `protected-decision-required` and `human-outcome-required` use exactly the affected or guarded action subset of `issue.transition-terminal-success`,`envelope.complete` |
| Resume (amend) | replace `findings:{available:false,items:[]}` with `findings:Finding[]`; `obligations` becomes `Evaluation[1..16]`; add `reviews:Review[]`, `concerns:Concern[]`, `deviations:Deviation[]`, `decisions:DecisionRef[]`, `human_outcomes:HumanOutcome[]`, `protected_decisions:ProtectedDecision[]`; each bounded by the list limit |
| TrustSource (amend) | adds `activated-actor` and `activated-human`; removes nothing |
| ActivityEvent (amend) | remove `ac-verified`, `concern-addressed`, `concern-raised`, `task-completed`; add `concern-recorded`, `concern-disposed`, `review-recorded`, `finding-addressed`, `finding-resolved`, `deviation-recorded`, `decision-referenced`, `human-outcome-recorded`, `requirement-referenced`, `requirement-linked` |

StopDetails gains three closed variants:

| Code | Details |
|---|---|
| review-rework-exhausted | `{envelope_id:ID,budget_name:ID,unit:"review-cycles",maximum:int,consumed:int,remaining:0,last_review_id:ID}` |
| protected-decision-required | `{envelope_id:ID,decision_class:ID,required_capabilities:ID[],affected_actions:text[]}` |
| human-outcome-required | `{envelope_id:ID,obligation_id:ID,outcome_class:ID}` |

## 3. New Closed Values

```text
CoverageItem   = {source:Source,assessed:boolean}
CoverageExclusion = {source:Source,reason:text}
Coverage       = {items:CoverageItem[1..64],exclusions:CoverageExclusion[0..64]}
Severity       = "blocker"|"major"|"minor"|"nit"
FindingInput   = {id:ID,subject:Subject,scope:{paths:text[0..256]},
                  severity:Severity,summary:text,blocking:boolean}
FindingDisposition = {id:ID,action:"addressed"|"resolved"|"inapplicable",
                  actor_id:ID,qualifying_review_id:?ID,reason:text,
                  disposed_at:timestamp,event_id:ID}
Finding        = FindingInput + {review_id:ID,created_at:timestamp,
                  dispositions:FindingDisposition[],
                  state:"open"|"addressed"|"resolved"|"inapplicable"}
Review         = {id:ID,run_id:ID,issue_id:ID,envelope_id:ID,envelope_revision:1,
                  actor_id:ID,profile:ID,repository:Repository,coverage:Coverage,
                  verdict:"pass"|"changes-required"|"blocked",
                  evidence_refs:Source[0..64],receipt_ids:ID[0..64],
                  findings:Finding[],resolved_finding_ids:ID[],
                  independence:"distinct-actor",recorded_at:timestamp,event_id:ID}
ConcernDisposition = {action:"carried-forward"|"rejected"|"superseded",
                  issue_id:?ID,supersedes_concern_id:?ID,supersedes_finding_id:?ID,
                  reason:text,actor_id:ID,disposed_at:timestamp,event_id:ID}
Concern        = {id:ID,subject:Subject,actor_id:ID,scope:{paths:text[0..256]},
                  severity:Severity,summary:text,reasoning:text,
                  references:Source[0..32],status:"open"|"disposed",
                  disposition:?ConcernDisposition,created_at:timestamp,event_id:ID}
Deviation      = {id:ID,subject:Subject,actor_id:ID,expected:text,observed:text,
                  authority:{actor_id:ID,reference:?text},
                  disposition:"accepted"|"rejected"|"pending-decision",
                  reasoning:text,references:Source[0..32],
                  supersedes_deviation_id:?ID,created_at:timestamp,event_id:ID}
DecisionRef    = {id:ID,class:ID,source:Source,
                  authority:{actor_id:ID,required_capabilities:ID[1..16]},
                  outcome:"approved"|"rejected"|"deferred",scope:{issue_id:ID},
                  created_at:timestamp,event_id:ID}
HumanOutcome   = {id:ID,actor_id:ID,
                  provenance:{activation_id:ID,activation_version:int,
                    profile_ref:text},class:ID,subject:Subject,
                  result:"approved"|"rejected"|"changes-required"|"acknowledged",
                  scope:{issue_id:ID},reasoning:text,references:Source[0..32],
                  created_at:timestamp,event_id:ID}
ProtectedDecisionInput = {ordinal:int>=0,decision_class:ID,
                  required_capabilities:ID[1..16] sorted-unique,
                  affected_actions:text[1..2] sorted-unique}
ProtectedDecision = ProtectedDecisionInput + {status:"satisfied"|"unsatisfied",
                  decision_ref_id:?ID,human_outcome_id:?ID}
RequirementRef = {id:ID,source:Source,status:"active"|"superseded"|"withdrawn",
                  superseded_by_id:?ID,created_at:timestamp,event_id:ID}
RequirementLink = {requirement_id:ID,issue_id:ID,
                  relation:"implements"|"verifies"|"relates-to",
                  created_at:timestamp,event_id:ID}
```

Exactly one of `ConcernDisposition.issue_id`, `supersedes_concern_id`, and
`supersedes_finding_id` is non-null for `carried-forward` and `superseded`, and
all three are null for `rejected`; `carried-forward` uses `issue_id` only. A
Concern may not supersede itself. `Concern.actor_id` is the raising Actor and
`Concern.disposition.actor_id` is the disposing Actor; they may be equal.

`Coverage.items` and `Coverage.exclusions` may not reference the same `Source`
digest. `Finding.state` is derived from the latest disposition, never supplied.
`Review.resolved_finding_ids` echoes the accepted `resolves_finding_ids` input.
`Review.independence` is the constant recorded basis; a Review is never stored
with any other value.

## 4. New And Replaced Leaves

Each row lists every required input key. Shared flags follow section 1 unless
the row says read-only.

| Operation / CLI | Input or read arguments | Result | Event / task |
|---|---|---|---|
| `review.record` / `fid review record` | `{run_id:ID,expected_run_version:int,profile:ID,repository:Repository,coverage:Coverage,verdict:text,findings:FindingInput[0..64],resolves_finding_ids:ID[0..64],evidence_refs:Source[0..64],receipt_ids:ID[0..64]}` | `{review:Review,budget:?Budget,stops:Stop[],attention:Attention[],idempotency_key:ID,replayed:boolean}` | review.recorded / 4,6,13 |
| `review.get` / `fid review get <review-id>` | read | `{review:Review}` | null / 13 |
| `review.list` / `fid review list` | read; optional `--issue`, `--envelope`, `--verdict`, `--limit` | `{items:Review[]}` | null / 13 |
| `finding.address` / `fid finding address` | `{finding_id:ID,claim_id:ID,expected_claim_version:int,reason:text}` | Mutation keyed `finding`: Finding | finding.addressed / 5,13 |
| `finding.list` / `fid finding list` | read; optional `--issue`, `--review`, `--state`, `--severity`, `--limit` | `{items:Finding[]}` | null / 5,13 |
| `concern.record` / `fid concern record` | `{id:ID,subject:Subject,scope:{paths:text[]},severity:Severity,summary:text,reasoning:text,references:Source[0..32]}` | Mutation keyed `concern`: Concern | concern.recorded / 8,13 |
| `concern.dispose` / `fid concern dispose` | `{concern_id:ID,action:text,issue_id:?ID,supersedes_concern_id:?ID,supersedes_finding_id:?ID,reason:text}` | Mutation keyed `concern`: Concern | concern.disposed / 8,13 |
| `concern.list` / `fid concern list` | read; optional `--subject <type>:<id>`, `--status`, `--severity`, `--limit` | `{items:Concern[]}` | null / 8,13 |
| `deviation.record` / `fid deviation record` | `{id:ID,subject:Subject,expected:text,observed:text,authority:{actor_id:ID,reference:?text},disposition:text,reasoning:text,references:Source[0..32],supersedes_deviation_id:?ID}` | Mutation keyed `deviation`: Deviation | deviation.recorded / 8,13 |
| `deviation.list` / `fid deviation list` | read; optional `--subject`, `--disposition`, `--limit` | `{items:Deviation[]}` | null / 8,13 |
| `decision.record-reference` / `fid decision record-reference` | `{id:ID,class:ID,source:Source,authority:{actor_id:ID,required_capabilities:ID[1..16]},outcome:text,scope:{issue_id:ID}}` | `{decision:DecisionRef,protected_decisions:ProtectedDecision[],stops:Stop[],attention:Attention[],idempotency_key:ID,replayed:boolean}` | decision.referenced / 9,13 |
| `decision.list` / `fid decision list` | read; optional `--class`, `--outcome`, `--issue`, `--limit` | `{items:DecisionRef[]}` | null / 9,13 |
| `human-outcome.record` / `fid human-outcome record` | `{id:ID,class:ID,subject:Subject,result:text,scope:{issue_id:ID},reasoning:text,references:Source[0..32]}` | `{outcome:HumanOutcome,protected_decisions:ProtectedDecision[],stops:Stop[],attention:Attention[],idempotency_key:ID,replayed:boolean}` | human-outcome.recorded / 10,11,13 |
| `human-outcome.list` / `fid human-outcome list` | read; optional `--class`, `--result`, `--issue`, `--limit` | `{items:HumanOutcome[]}` | null / 10,13 |
| `requirement.record-reference` / `fid requirement record-reference` | `{id:ID,source:Source,status:text}` | Mutation keyed `requirement`: RequirementRef | requirement.referenced / 12,13 |
| `requirement.link` / `fid requirement link` | `{requirement_id:ID,issue_id:ID,relation:text}` | Mutation keyed `link`: RequirementLink | requirement.linked / 12,13 |
| `requirement.get` / `fid requirement get <requirement-id>` | read | `{requirement:RequirementRef,links:RequirementLink[]}` | null / 12,13 |
| `requirement.list` / `fid requirement list` | read; optional `--status`, `--issue`, `--limit` | `{items:RequirementRef[]}` | null / 12,13 |
| `envelope.create` / existing leaf | `obligations:ObligationInput[1..16]`, add `protected_decisions:ProtectedDecisionInput[0..16]` | existing result with amended Envelope | envelope.created / 2 |
| `envelope.activate` / existing leaf | unchanged input | existing result; `evaluation` becomes `Evaluation[]` | envelope.activated / 2,9 |
| `claim.acquire` / existing leaf | add `mode:"mutation"\|"review"` | existing result with amended Claim | claim.acquired / 3 |
| `evidence.link` / existing leaf | `{obligation_id:ID,evidence_type:text,receipt_id:?ID,review_id:?ID,human_outcome_id:?ID,target_repository:Repository,required_artifact_ids:ID[]}` | `{evidence:Evidence,evaluation:Evaluation,idempotency_key:ID,replayed:boolean}` | evidence.linked / 6,10 |
| `evidence.evaluate` / existing leaf | `{obligation_id:?ID,target_repository:Repository}` | `{evaluations:Evaluation[]}` when `obligation_id` is null, otherwise retained `{evaluation:Evaluation}` | null / 6,13 |

`review.record` requires an active `review`-mode Claim for the named Run, the
`review.record` capability in the Actor profile, and a matching Envelope grant.
`finding.address` requires an active mutating Claim on the reviewed Issue and
the `finding.address` capability. `human-outcome.record` requires an Actor whose
profile `actor_type` is `human`, the `human-outcome.record` capability, and
active finalized Kit Activation provenance. `concern.record`,
`concern.dispose`, `deviation.record`, `decision.record-reference`,
`requirement.record-reference`, and `requirement.link` require their matching
capability in the Actor profile; the first four additionally require an
Envelope grant when their subject resolves to an Envelope-scoped Issue.

### 4.1 Activated Capabilities And Grant Actions

The Phase 22 closed set of 23 names becomes 30. Added names, each legal in an
activated profile and an Envelope grant:

| Name | Admission rule / implementation owner |
|---|---|
| `review.record` | Profile plus Envelope grant; review-mode Claim; task 4 |
| `finding.address` | Profile plus Envelope grant; mutating Claim; task 5 |
| `concern.record` | Profile plus Envelope grant; task 8 |
| `concern.dispose` | Profile plus Envelope grant; task 8 |
| `deviation.record` | Profile plus Envelope grant; task 8 |
| `decision.record-reference` | Profile plus Envelope grant; task 9 |
| `human-outcome.record` | Profile with `actor_type:"human"` plus Envelope grant and active finalized Activation; task 10 |

`requirement.record-reference`, `requirement.link`, and every Phase 23 read are
operation IDs only and stay out of the allowlist, matching the existing
treatment of activation and read leaves. Unknown actions reject. Task 18
compares the explicit 30-name set with independently pinned corpus members and
exercises allow and reject cases; parsing the table to produce both expected
and actual values cannot pass.

## 5. Failure Codes

Added codes are closed, all `retryable:false`, subject as named.

| Code | Category / exit | Exact details |
|---|---|---|
| review-not-independent | precondition / 4 | `{review_actor_id:ID,conflicting_subject:Subject,conflicting_actor_id:ID}` |
| review-verdict-invalid | invalid-argument / 2 | `{verdict:text,reason:text}` |
| review-changes-required | precondition / 4 | `{review_id:ID,obligation_id:ID,blocking_finding_ids:ID[]}` |
| review-rework-exhausted | precondition / 4 | `{budget:Budget,stop_ids:ID[],attention_ids:ID[],event:?EventRef,replayed:boolean}` |
| finding-unresolved | precondition / 4 | `{obligation_id:ID,finding_ids:ID[]}` |
| finding-already-disposed | conflict / 4 | `{finding_id:ID,state:text}` |
| finding-not-found / review-not-found / concern-not-found / deviation-not-found / decision-not-found / outcome-not-found / requirement-not-found | not-found / 3 | `{id:ID}` |
| concern-already-disposed | conflict / 4 | `{concern_id:ID,disposed_at:timestamp}` |
| claim-mode-invalid | precondition / 4 | `{claim_id:ID,mode:text,required_mode:text}` |
| protected-decision-required | authority / 5 | `{envelope_id:ID,decision_class:ID,required_capabilities:ID[],affected_actions:text[]}` |
| outcome-not-permitted | precondition / 4 | `{obligation_id:ID,obligation_type:text,outcome_id:ID,result:text}` |
| evidence-type-not-accepted | precondition / 4 | `{obligation_id:ID,evidence_type:text,accepted_evidence_types:text[]}` |
| source-digest-mismatch | precondition / 4 | `{display_label:text,path:?text,expected_digest:digest,observed_digest:?digest}` |
| repository-mismatch | conflict / 4 | `{expected:Repository,observed:Repository}` |
| human-outcome-required | precondition / 4 | `{envelope_id:ID,obligation_id:ID,outcome_class:ID}` |

`review-rework-exhausted` and `protected-decision-required` follow the Phase 22
committed-denial rule: when the denial durably raises a Stop and Attention the
response is `ok:false` with the committed `EventRef` and immutable Stop and
Attention IDs, stored and replayed under its original idempotency key as a
complete typed snapshot. A denial that commits nothing has null `EventRef`.

`ReasonCode` adds `finding-unresolved`, `review-changes-required`,
`protected-decision-required`, `human-outcome-absent`,
`review-not-independent`, and `outcome-superseded`. `ReasonSubjectType` adds
`review`, `finding`, `decision`, and `human-outcome`.

## 6. Manifest And Availability

Manifest schema stays `fid.export-manifest.v2` with its Phase 22 fields.
`projection_schema` becomes the exporter's active accepted migration: 17 during
implementation and 18 at phase end. Import accepts only the binary's exact
current accepted schema; final capability qualification requires 18. Portable
export, validation, and import cover every Phase 23 table, Event, and
projection; recovery-grade equivalence is proven for the new domain exactly as
for Phase 22 recovery state.

Final changed capability rows keep `schema_min` and `schema_max` at 1 and are
listed in `contracts.md` section 15. Discovery, resume, exports, and the shared
fake and fid fixtures must agree. Schema presence alone does not enable a row.

## 7. Human Output Obligations

Human text obligations are fixed by `contracts.md` section 14. No Human view
adds a JSON member, changes a mutation rule, reproduces project prose, or
presents a Human Outcome, Concern, or Deviation as deterministic proof. Every
new leaf has at least one Human-mode conformance fixture with exact expected
text, exercised in the ordinary suite and in the tagged real-process suite.
