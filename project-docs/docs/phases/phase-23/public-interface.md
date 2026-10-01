# Phase 23 Public Interface Amendment

> ACCEPTED under DEC-222 with `contracts.md` and `schema.md`.
> Amended 2026-09-15 by DEC-226 (failure-code producers, precedence, and
> non-committing denial semantics) and DEC-224 (`human-outcome-required`
> producers), both **approved by the Human 2026-09-19**, and by DEC-263
> (Stop-blocked denial details, Tech Lead under Human delegation 2026-09-29),
> amended by DEC-269 (`RequirementRef.requirement_id`, Human approval pending),
> and clarified by DEC-265 (section 5 declared-versus-added count, wording
> only, 2026-09-30); amended passages name their DEC inline. `GLOSSARY.md` defines **Failure
> code** and **Evaluation reason**: section 5 lists the failure set, and the
> `ReasonCode` paragraph the evaluation set.
> Phase 21 and Phase 22 `public-interface.md` remain the baseline for every
> leaf, value, code, and exit this document does not replace.

## 1. Transport And Shared Grammar

Retain `fid-cli.v1`, `--json`, `--project`, the existing error envelope, exit
mapping, and the Phase 22 notation (`?T`, `T[]`, ID, text, timestamp, digest,
Subject, Source, Repository, Version); `Subject` keeps its Phase 22 meaning
exactly, and section 3 disambiguates it from the three other Phase 23 subject
vocabularies (DEC-232). Every Phase 23 mutation takes no
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
| ErrorSubjectType (amend) | The failure-envelope subject vocabulary, **authoritatively the implemented `SubjectType` in `internal/cli/codes.go:88`**, which is the shipped surface this row cites. Add `review`, `finding`, `concern`, `deviation`, `decision`, `human-outcome`, taking the implemented 19 names to **25**. `requirement` is retained and now names a `requirement_refs` row. The prose set in `../phase-21/public-interface.md` section 1 has drifted from the implementation and is **not** authoritative here: it lists 22 names, five of which (`actor`, `idempotency`, `obligation`, `path`, `repository`) were never implemented, and omits `attention` and `stop`, which Phase 22 added. Phase 23 adds and removes nothing else; reconciling that prose is a phase-end spec-sync item, not a code change (I-34, DEC-233). The task that first returns an error carrying a new subject registers its `SubjectType` constant in the same change - task 4 `review`, 5 `finding`, 8 `concern` and `deviation`, 9 `decision`, 10 `human-outcome` - the same rule DEC-229 set for Condition types. This row does **not** amend the Stop/Attention `Subject` or the judgment `JudgmentSubject` below (DEC-232). |
| ObligationInput (replace) | `{id:ID,type:"gate-evidence"\|"review"\|"human-outcome",applicability:"required",guarded_actions:text[1..2],subject:{type:"issue",id:ID},registration:?{id:ID,version:int>=1},review_profile:?ID,outcome_class:?ID,accepted_evidence_types:text[1],accepted_trust_sources:text[1],freshness_target:"repository-tree"\|"none",artifact_policy:"required-available"\|"metadata-sufficient",independence:"fid-executed"\|"distinct-actor"\|"human",required_lane:null,source:Source}`; the exact per-type field combination is fixed by `contracts.md` section 2 |
| Obligation (amend) | `ObligationInput` plus `{derived_status:"satisfied"\|"unsatisfied",reasons:Reason[]}`, unchanged in shape |
| Evaluation (amend) | `artifact_availability` domain becomes `"available"\|"missing"\|"not-applicable"`; all other fields unchanged |
| Budget (amend) | unit domain adds `review-cycles`; `exhaustion_action` domain adds `review-rework-exhausted`; all other Phase 22 rules unchanged |
| Claim (amend) | `mode` domain becomes `"mutation"\|"review"`; all other fields, statuses, and release reasons unchanged |
| Envelope (amend) | `obligations:Obligation[1..16]`; add `protected_decisions:ProtectedDecision[0..16]` |
| Evidence (amend) | `{id:ID,obligation_id:ID,evidence_type:"receipt"\|"review"\|"human-outcome",receipt_id:?ID,review_id:?ID,human_outcome_id:?ID,producer:"fid"\|"activated-actor"\|"human",trust_source:"fid-executed"\|"activated-actor"\|"activated-human",target_repository:Repository,validity:"valid"\|"invalid",required_artifact_ids:ID[],created_at:timestamp}`; exactly one subject ID is non-null and matches `evidence_type` |
| Condition (amend) | type domain adds `review-recorded`, `finding-resolved`, `human-outcome-recorded`, `decision-satisfied`. Value adds `outcome_class:ID`, `decision_class:ID`, `finding_id:ID`, each **present only on its own type and absent on every other** (`omitempty` on these three fields, and on no others). The pairing is `contracts.md` section 11's: `outcome_class` on `human-outcome-recorded`, `decision_class` on `decision-satisfied`, `finding_id` on `finding-resolved`, and **no discriminator on `review-recorded`** - four types and three fields do not divide evenly, and `review-recorded` is the type without one (DEC-236). A judgment Condition that has a discriminator therefore carries **five** fields - `type`, `subject`, `budget_name:null`, `capability_id:null`, and its own - while `review-recorded` carries four and a Phase 22 Condition is unchanged byte-for-byte at four. `budget_name` and `capability_id` keep present-and-null treatment and gain **no** `omitempty`. Exclusivity is structural, enforced by per-type journal shapes keyed on `type` rather than by a nullable-binding relation. Amended under DEC-235 resolving ESC-23-001, Human-approved 2026-09-20; the superseded reading was present-and-null on all five discriminator fields, which is not monotonic over stored history. The per-type pairing and the `review-recorded` carve-out are stated under DEC-236 (QA R-26), restating accepted `contracts.md` section 11 and amending no encoding. |
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
JudgmentSubject = {type:"issue"|"envelope"|"run"|"review"|"finding"
                  |"requirement",id:ID}
CoverageItem   = {source:Source,assessed:boolean}
CoverageExclusion = {source:Source,reason:text}
Coverage       = {items:CoverageItem[1..64],exclusions:CoverageExclusion[0..64]}
Severity       = "blocker"|"major"|"minor"|"nit"
FindingInput   = {id:ID,subject:JudgmentSubject,scope:{paths:text[0..256]},
                  severity:Severity,summary:text,blocking:boolean}
InapplicableFindingInput = {finding_id:ID,reason:text}
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
Concern        = {id:ID,subject:JudgmentSubject,actor_id:ID,
                  scope:{paths:text[0..256]},
                  severity:Severity,summary:text,reasoning:text,
                  references:Source[0..32],status:"open"|"disposed",
                  disposition:?ConcernDisposition,created_at:timestamp,event_id:ID}
Deviation      = {id:ID,subject:JudgmentSubject,actor_id:ID,
                  expected:text,observed:text,
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
                    profile_ref:text},class:ID,subject:JudgmentSubject,
                  result:"approved"|"rejected"|"changes-required"|"acknowledged",
                  scope:{issue_id:ID},reasoning:text,references:Source[0..32],
                  created_at:timestamp,event_id:ID}
ProtectedDecisionInput = {ordinal:int>=0,decision_class:ID,
                  required_capabilities:ID[1..16] sorted-unique,
                  affected_actions:text[1..2] sorted-unique}
ProtectedDecision = ProtectedDecisionInput + {status:"satisfied"|"unsatisfied",
                  decision_ref_id:?ID,human_outcome_id:?ID}
RequirementRef = {id:ID,requirement_id:ID,source:Source,
                  status:"active"|"superseded"|"withdrawn",
                  superseded_by_id:?ID,created_at:timestamp,event_id:ID}
                  -- requirement_id added by DEC-269 (pending Human approval)
RequirementLink = {requirement_id:ID,issue_id:ID,
                  relation:"implements"|"verifies"|"relates-to",
                  created_at:timestamp,event_id:ID}
```

`JudgmentSubject` is a new **name** for a domain already accepted under
DEC-222, not a new domain: it is the closed set stated in `contracts.md`
section 8 and enforced identically by the `subject_type` CHECK on
`review_findings`, `judgment_concerns`, `deviations`, and `human_outcomes` in
`schema.md` section 2. Naming it is required because Phase 23 has four distinct
subject vocabularies, and reusing the bare name `Subject` for two of them
already produced one wrong instruction (DEC-232, R-18):

| Name | Domain | Where it applies |
|---|---|---|
| `Subject` | `run`, `envelope`, `activation` - unchanged by Phase 23 | Stop, Attention, Condition, NextAction; `stops`/`attention_items` DDL; `currentValueShapes["Subject"]` |
| `JudgmentSubject` | `issue`, `envelope`, `run`, `review`, `finding`, `requirement` | Finding, Concern, Deviation, HumanOutcome |
| `ErrorSubjectType` | section 2 row; `internal/cli/codes.go:88` | the failure envelope's `error.subject`, and `Subject`-typed failure `details` fields such as `review-not-independent.conflicting_subject` |
| `ReasonSubjectType` | section 5; adds `review`, `finding`, `decision`, `human-outcome` | `Reason.subject_type` in Evaluations |

Neither `Subject` nor `JudgmentSubject` contains the other: Stop admits
`activation` and the judgment set does not, and the judgment set admits
`issue`, `review`, `finding`, and `requirement` and Stop does not.

Exactly one of `ConcernDisposition.issue_id`, `supersedes_concern_id`, and
`supersedes_finding_id` is non-null for `carried-forward` and `superseded`, and
all three are null for `rejected`; `carried-forward` uses `issue_id` only. A
Concern may not supersede itself. `Concern.actor_id` is the raising Actor and
`Concern.disposition.actor_id` is the disposing Actor; they may be equal.

`Coverage.items` and `Coverage.exclusions` may not reference the same `Source`
digest. `Finding.state` is derived from the latest disposition, never supplied.
`Review.resolved_finding_ids` echoes the accepted `resolves_finding_ids` input.
It contains only `resolved` actions. An `inapplicable` action, reason and
qualifying Review are read through the existing Finding's `dispositions` and
derived `state`, using `finding.list` or `review.get` on its creating Review;
the `Review` result shape gains no field (DEC-242).
`Review.independence` is the constant recorded basis; a Review is never stored
with any other value.

## 4. New And Replaced Leaves

Each row lists every required input key. `optional(key:Type)` explicitly marks
an omittable key, not a nullable value; all other keys remain required. Shared
flags follow section 1 unless the row says read-only.

| Operation / CLI | Input or read arguments | Result | Event / task |
|---|---|---|---|
| `review.record` / `fid review record` | `{run_id:ID,expected_run_version:int,profile:ID,repository:Repository,coverage:Coverage,verdict:text,findings:FindingInput[0..64],resolves_finding_ids:ID[0..64],optional(inapplicable_findings:InapplicableFindingInput[0..64]),evidence_refs:Source[0..64],receipt_ids:ID[0..64]}` | `{review:Review,budget:?Budget,stops:Stop[],attention:Attention[],idempotency_key:ID,replayed:boolean}` | review.recorded / 4,5,6,13; optional input begins task 5 |
| `review.get` / `fid review get <review-id>` | read | `{review:Review}` | null / 13 |
| `review.list` / `fid review list` | read; optional `--issue`, `--envelope`, `--verdict`, `--limit` | `{items:Review[]}` | null / 13 |
| `finding.address` / `fid finding address` | `{finding_id:ID,claim_id:ID,expected_claim_version:int,reason:text}` | Mutation keyed `finding`: Finding | finding.addressed / 5,13 |
| `finding.list` / `fid finding list` | read; optional `--issue`, `--review`, `--state`, `--severity`, `--limit` | `{items:Finding[]}` | null / 5,13 |
| `concern.record` / `fid concern record` | `{id:ID,subject:JudgmentSubject,scope:{paths:text[]},severity:Severity,summary:text,reasoning:text,references:Source[0..32]}` | Mutation keyed `concern`: Concern | concern.recorded / 8,13 |
| `concern.dispose` / `fid concern dispose` | `{concern_id:ID,action:text,issue_id:?ID,supersedes_concern_id:?ID,supersedes_finding_id:?ID,reason:text}` | Mutation keyed `concern`: Concern | concern.disposed / 8,13 |
| `concern.list` / `fid concern list` | read; optional `--subject <type>:<id>`, `--status`, `--severity`, `--limit` | `{items:Concern[]}` | null / 8,13 |
| `deviation.record` / `fid deviation record` | `{id:ID,subject:JudgmentSubject,expected:text,observed:text,authority:{actor_id:ID,reference:?text},disposition:text,reasoning:text,references:Source[0..32],supersedes_deviation_id:?ID}` | Mutation keyed `deviation`: Deviation | deviation.recorded / 8,13 |
| `deviation.list` / `fid deviation list` | read; optional `--subject`, `--disposition`, `--limit` | `{items:Deviation[]}` | null / 8,13 |
| `decision.record-reference` / `fid decision record-reference` | `{id:ID,class:ID,source:Source,authority:{actor_id:ID,required_capabilities:ID[1..16]},outcome:text,scope:{issue_id:ID}}` | `{decision:DecisionRef,protected_decisions:ProtectedDecision[],stops:Stop[],attention:Attention[],idempotency_key:ID,replayed:boolean}` | decision.referenced / 9,13 |
| `decision.list` / `fid decision list` | read; optional `--class`, `--outcome`, `--issue`, `--limit` | `{items:DecisionRef[]}` | null / 9,13 |
| `human-outcome.record` / `fid human-outcome record` | `{id:ID,class:ID,subject:JudgmentSubject,result:text,scope:{issue_id:ID},reasoning:text,references:Source[0..32]}` | `{outcome:HumanOutcome,protected_decisions:ProtectedDecision[],stops:Stop[],attention:Attention[],idempotency_key:ID,replayed:boolean}` | human-outcome.recorded / 10,11,13 |
| `human-outcome.list` / `fid human-outcome list` | read; optional `--class`, `--result`, `--issue`, `--limit` | `{items:HumanOutcome[]}` | null / 10,13 |
| `requirement.record-reference` / `fid requirement record-reference` | `{id:ID,source:Source,status:text}` | Mutation keyed `requirement`: RequirementRef | requirement.referenced / 12,13 |
| `requirement.link` / `fid requirement link` | `{requirement_id:ID,issue_id:ID,relation:text}` | Mutation keyed `link`: RequirementLink | requirement.linked / 12,13 |
| `requirement.get` / `fid requirement get <requirement-id \| reference-id>` (DEC-269) | read | `{requirement:RequirementRef,links:RequirementLink[]}` | null / 12,13 |
| `requirement.list` / `fid requirement list` | read; optional `--status`, `--issue`, `--limit` | `{items:RequirementRef[]}` | null / 12,13 |
| `envelope.create` / existing leaf | `obligations:ObligationInput[1..16]`, add `protected_decisions:ProtectedDecisionInput[0..16]` | existing result with amended Envelope | envelope.created / 2 |
| `envelope.activate` / existing leaf | unchanged input | existing result; `evaluation` becomes `Evaluation[]` | envelope.activated / 2,9 |
| `claim.acquire` / existing leaf | add `mode:"mutation"\|"review"` | existing result with amended Claim | claim.acquired / 3 |
| `evidence.link` / existing leaf | `{obligation_id:ID,evidence_type:text,receipt_id:?ID,review_id:?ID,human_outcome_id:?ID,target_repository:Repository,required_artifact_ids:ID[]}` | `{evidence:Evidence,evaluation:Evaluation,idempotency_key:ID,replayed:boolean}` | evidence.linked / 6,10 |
| `evidence.evaluate` / existing leaf | `{obligation_id:?ID,target_repository:Repository}` | `{evaluations:Evaluation[]}` when `obligation_id` is null, otherwise retained `{evaluation:Evaluation}` | null / 6,13 |

**Task 5 input addition (DEC-242).** `inapplicable_findings` is omitted or an
array of closed `{finding_id,reason}` objects. Explicit `null`, a missing or
null member, and unknown or duplicate keys reject. Each reason is valid UTF-8,
at most the retained text limit of 4096 bytes, and nonblank after whitespace
trimming; store the supplied valid text verbatim. IDs are unique within each
terminal-disposition collection and disjoint across `resolves_finding_ids` and
`inapplicable_findings`; their combined length is at most 64. Both actions use
the same later-Review, Issue/revision, independence and verdict checks.

Omission and `[]` both canonicalize to omission of `inapplicable_findings` for
request hashing. Thus a request using only Task 4 fields retains its exact
canonical bytes and request digest. A nonempty collection participates in the
digest, including its reasons; changed reasons conflict under the same key.
No Review output, journal value shape, DDL, or Event family is added or widened.
Disposition effects commit atomically within `review.recorded`.

Terminal-disposition input failures use the existing codes: malformed collection,
missing/null/blank/oversize reason, duplicate or overlapping ID, or combined
limit violation returns `invalid-argument` (exit 2, existing closed details).
An unknown persisted Finding ID returns `finding-not-found` (exit 3,
`{id:ID}`); an already `resolved` or `inapplicable` Finding returns
`finding-already-disposed` (exit 4, `{finding_id:ID,state:text}`). A Finding
created in this same Review or belonging to another Issue/revision returns
`invalid-argument`, with `argument:"finding_id"` and the exact ID in `value`.
The independent-Actor and resolving-verdict failures remain
`review-not-independent` and `review-verdict-invalid`. All reject without an
Event or partial projection change. These rules are implemented by Task 5;
Task 4's existing input, criteria and commands remain unchanged.

`review.record` requires an active `review`-mode Claim for the named Run, the
`review.record` capability in the Actor profile, and a matching Envelope grant.
`finding.address` requires an active mutating Claim on the reviewed Issue and
the `finding.address` capability. `human-outcome.record` requires an Actor whose
profile `actor_type` is `human`, the `human-outcome.record` capability, and
active finalized Kit Activation provenance. `concern.record`,
`concern.dispose`, `deviation.record`, and `decision.record-reference` require
their matching capability in the Actor profile, and additionally require an
Envelope grant when their subject resolves to an Envelope-scoped Issue.
`requirement.record-reference` and `requirement.link` require a trusted Actor
resolved from the finalized Kit Activation, but no Actor-profile capability and
no Envelope grant. They are operation IDs only (section 4.1; corrected under
DEC-260, 2026-09-27).

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
| `human-outcome.record` | Profile with `actor_type:"human"` and active finalized Activation; no Envelope grant (`contracts.md` section 10; corrected under DEC-251, Human 2026-09-26). The name stays legal in a grant's action list. Task 10 |

`requirement.record-reference`, `requirement.link`, and every Phase 23 read are
operation IDs only and stay out of the allowlist, matching the existing
treatment of activation and read leaves. Unknown actions reject. Task 18
compares the explicit 30-name set with independently pinned corpus members and
exercises allow and reject cases; parsing the table to produce both expected
and actual values cannot pass.

## 5. Failure Codes

Section 5 declares **19** codes, all closed, all `retryable:false`, subject as
named (DEC-265). **18** of them are added by Phase 23. The nineteenth,
`requirement-not-found`, is **retained**. It has been in `internal/cli/codes.go`
since the Phase 20 response contract (`ab4ac8f`, phase-20/task-1), and its
category, exit, subject, and `{id:ID}` details are unchanged. It is declared here because Phase 23 removes
its V1 producers and pins new ones, `requirement.get` and `requirement.link`
(section 5.1). A conformance corpus for this section covers all 19; a claim
about the *added* set counts 18.

| Code | Category / exit | Exact details |
|---|---|---|
| review-not-independent | precondition / 4 | `{review_actor_id:ID,conflicting_subject:Subject,conflicting_actor_id:ID}` |
| review-verdict-invalid | invalid-argument / 2 | `{verdict:text,reason:text}` |
| review-changes-required | precondition / 4 | `{review_id:ID,obligation_id:ID,blocking_finding_ids:ID[]}` |
| review-rework-exhausted | precondition / 4 | `{budget:Budget,stop_ids:ID[],attention_ids:ID[],event:?EventRef,replayed:boolean}` |
| finding-already-disposed | conflict / 4 | `{finding_id:ID,state:text}` |
| finding-not-found / review-not-found / concern-not-found / deviation-not-found / outcome-not-found / requirement-not-found | not-found / 3 | `{id:ID}` |
| concern-already-disposed | conflict / 4 | `{concern_id:ID,disposed_at:timestamp}` |
| claim-mode-invalid | precondition / 4 | `{claim_id:ID,mode:text,required_mode:text}` |
| protected-decision-required | authority / 5 | `{envelope_id:ID,decision_class:ID,required_capabilities:ID[],affected_actions:text[],stop_ids:ID[],attention_ids:ID[],event:?EventRef,replayed:boolean}` (DEC-263) |
| outcome-not-permitted | precondition / 4 | `{obligation_id:ID,obligation_type:text,outcome_id:ID,result:text}` |
| evidence-type-not-accepted | precondition / 4 | `{obligation_id:ID,evidence_type:text,accepted_evidence_types:text[]}` |
| source-digest-mismatch | precondition / 4 | `{display_label:text,path:?text,expected_digest:digest,observed_digest:?digest}` |
| repository-mismatch | conflict / 4 | `{expected:Repository,observed:Repository}` |
| human-outcome-required | precondition / 4 | `{envelope_id:ID,obligation_id:ID,outcome_class:ID,stop_ids:ID[],attention_ids:ID[],event:?EventRef,replayed:boolean}` (DEC-263) |

`decision-not-found` is **removed** from the added set (DEC-230). Pinning each
`*-not-found` producer under I-25 showed that no Phase 23 leaf takes a Decision
Reference ID as input: there is no `decision get`, and no mutation names one.
DEC-226's own criterion is that every added code has a producing operation, so a
code nothing can return is dead vocabulary rather than an unimplemented one. The
added set is therefore 18 codes, not 19; with the retained
`requirement-not-found`, the section declares 19 (DEC-265). Reading a Decision
Reference stays `decision.list --class`. If a `decision.get` leaf is ever
added, the code returns with it.

`finding-unresolved` is a `ReasonCode` only and is **not** a failure code
(DEC-226). No leaf rejects on it: linking a `pass` Review that leaves an earlier
blocking Finding unresolved is accepted, and the Obligation simply stays
unsatisfied with that reason (`contracts.md` section 5). Task 6 therefore adds it
to the evaluation reason precedence order in `internal/workflow/evaluation.go`,
not to `internal/cli/codes.go`.

`ReasonCode` adds `finding-unresolved`, `review-changes-required`,
`protected-decision-required`, `human-outcome-absent`,
`review-not-independent`, and `outcome-superseded`. `ReasonSubjectType` adds
`review`, `finding`, `decision`, and `human-outcome`. `human-outcome-absent`
replaces the generic `obligation-evidence-absent` for an unsatisfied
`human-outcome` Obligation with no linked Outcome, and `outcome-superseded` is
returned when section 7 supersession displaces the Outcome that previously
satisfied it.

### 5.1 Producers And Precedence (DEC-226)

Every added code names the operations that can return it. `evidence.link`
rejections follow the single order in `contracts.md` section 13.

| Code | Returned by |
|---|---|
| review-not-independent | `review.record`; `evidence.link` when independence no longer holds at link time |
| review-verdict-invalid | `review.record` |
| review-changes-required | `evidence.link` |
| review-rework-exhausted | `claim.acquire`, `run.start` while the exhaustion Stop is open |
| finding-already-disposed | `finding.address`; `review.record` attempting either terminal action against an already terminal Finding (DEC-242, task 5) |
| concern-already-disposed | `concern.dispose` |
| claim-mode-invalid | `review.record` under a `mutation` Claim; `finding.address` under a `review` Claim; `run.recover` and `run.abandon` adopting a path under a `review` Claim |
| protected-decision-required | `envelope.complete` and `issue.transition-terminal-success`, for the declared affected actions only |
| human-outcome-required | `envelope.complete` and `issue.transition-terminal-success`, for the blocked actions of the open readiness Stop only |
| outcome-not-permitted | `evidence.link` |
| evidence-type-not-accepted | `evidence.link` |
| source-digest-mismatch | `decision.record-reference`, `requirement.record-reference` |
| repository-mismatch | `review.record` |
| review-not-found | `review.get`; `evidence.link` with an unknown `review_id` |
| finding-not-found | `finding.address`; `review.record` with an unknown ID in `resolves_finding_ids` or `inapplicable_findings`; `concern.dispose` with an unknown `supersedes_finding_id` |
| concern-not-found | `concern.dispose`, by `concern_id` or by `supersedes_concern_id` |
| deviation-not-found | `deviation.record` with an unknown `supersedes_deviation_id`; there is no `deviation get` leaf |
| outcome-not-found | `evidence.link` with an unknown `human_outcome_id`; there is no `human-outcome get` leaf |
| requirement-not-found | `requirement.get`; `requirement.link` with an unknown `requirement_id` |

Where more than one applies to the same guarded-action request, exactly one code
is returned in this order:

1. every retained Phase 21/22 precondition and authority failure other than
   `obligation-unsatisfied`, unchanged (I-19);
2. `protected-decision-required`, lowest unsatisfied protected-Decision ordinal
   first, because an authority failure (exit 5) outranks a precondition;
3. `review-rework-exhausted` (disjoint action set; it can never collide with
   the two above);
4. `human-outcome-required`;
5. `obligation-unsatisfied` with its `reasons` array, which is returned only
   when no open Stop blocks the requested action.

An open Stop therefore always speaks for itself; the generic
`obligation-unsatisfied` denial never masks one. The `reasons` array keeps
reporting every unsatisfied Obligation reason, in the section 2 reason
precedence order, whichever code is returned.

### 5.2 Denials Commit Nothing (DEC-226)

Phase 23 adds no denial Event. `review-rework-exhausted`,
`protected-decision-required`, and `human-outcome-required` are raised by the
successful Events named in `contracts.md` sections 6, 9, and 10, which store
their own typed snapshots under their own idempotency keys. The later requests
those Stops block - `claim.acquire`, `run.start`,
`issue.transition-terminal-success`, `envelope.complete` - create nothing:
no Claim, Run, transition, completion, Event, or `idempotency_results` row.
Their response is `ok:false` with the immutable Stop and Attention IDs raised
earlier, null `EventRef`, and `replayed:false`, and repeating the request while
the Stop is unresolved returns byte-identical Stop and Attention IDs. The
Phase 22 committed-denial rule is retained unchanged for `run.execution-blocked`,
which does durably raise its Stop and does store its Denial snapshot.

`review-rework-exhausted` keeps its declared detail shape, and under this rule
its `event` is always null and its `replayed` always `false`; the fields stay
present so the closed shape is unchanged. VAL-056 asserts exactly that.

**These fields go in `details` (DEC-263).** The Stop and Attention IDs, the
`EventRef`, and the replay flag are carried as the `stop_ids`, `attention_ids`,
`event`, and `replayed` detail fields of all three Stop-blocked codes. The
`fid-cli.v1` envelope does not change. `protected-decision-required` and
`human-outcome-required` details are their section 2 StopDetails followed by
those four fields:
- `stop_ids` holds the unresolved Stops of that code that block the requested
  action and that the details describe, sorted by ID.
- `attention_ids` holds the Attention linked to those Stops, sorted by ID.
- `event` is always null and `replayed` always `false`.

The fields belong to the failure response only. StopDetails, and every
projection of a Stop, are unchanged.

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
