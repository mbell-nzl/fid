# Phase 23 Slice Contract - Review And Human Judgment

> Status: ACCEPTED by Human under DEC-222, 2026-09-14.
> Amended 2026-09-15 by DEC-224 (readiness Stop producers), DEC-225
> (independence reading), and DEC-226 (failure-code precedence, denial
> semantics, and the Failure code / Evaluation reason vocabulary split), all
> **approved by the Human 2026-09-19**. Amended passages name their DEC inline;
> Tech Lead planning clarifications under DEC-228 are marked the same way.
> Author: Tech Lead, 2026-09-14. Read with `public-interface.md` and `schema.md`.
> Planning base: `0195897a347dfe1872db0f7be3800c5c3ccdefcf`.
> Roadmap slice 23 of ADR-0012; Phase 21/22 contracts remain authoritative
> except where this document, `public-interface.md`, or `schema.md` replace them.

## 1. Boundary And Authority

One selected Issue and one immutable Envelope Revision remain the execution
boundary. This phase adds independent Review coverage, Findings and their
dispositions, review rework accounting, generic Concerns and Deviations,
operational Decision References with protected-Decision blocking, provenanced
Human Outcomes, Human Outcome-based Attention resolution, and requirement
references that replace V1 requirement authoring.

It does not add multiple selected Issues, Envelope revision authoring, parallel
Claims, generic `retry-exhausted`, attested external Receipts, trajectory
telemetry, or the economics half of the Human Validation Brief. Those remain
Phase 24 and Phase 25.

Retain `internal/workflow` for Envelope, Obligation, Claim, Run, Stop,
Attention, and budget rules; `internal/evidence` for Evidence records and
evaluation; `internal/claims` for leases; `internal/validation` for shared text
rules; `internal/project` for repository observation. Add exactly one new
domain package, `internal/review`, for Review, Finding, Concern, Deviation,
Decision Reference, and Human Outcome rules. Requirement references replace
`internal/requirements` in place. The CLI owns parsing and rendering only.
No new runtime dependency is authorized by this contract.

fid stores no Review prose body, no Decision body, no QA Plan, no Validation
Contract, no handoff, and no completion report. Project files stay canonical and
are referenced by `Source` with repository identity and digest (PRD-REV-009).

## 2. Obligation Generalization

`envelope_obligations` stops being a single hard-coded `gate-evidence` row. An
Envelope revision carries 1..16 Obligations with unique IDs and exactly one of
three types:

| Type | Accepted evidence types | Accepted trust sources | Independence |
|---|---|---|---|
| `gate-evidence` | `["receipt"]` | `["fid-executed"]` | `fid-executed` |
| `review` | `["review"]` | `["activated-actor"]` | `distinct-actor` |
| `human-outcome` | `["human-outcome"]` | `["activated-human"]` | `human` |

A `gate-evidence` Obligation keeps the exact Phase 21/22 value, including its
non-null `registration`. A `review` or `human-outcome` Obligation has null
`registration`; a `review` Obligation additionally carries a non-null
`review_profile`, and a `human-outcome` Obligation a non-null `outcome_class`.
Every Obligation keeps `applicability:"required"`, `subject:{type:"issue"}`
equal to the primary Issue, a non-null `source`, and null `required_lane`.

`guarded_actions` stays a sorted non-empty subset of `envelope.complete` and
`issue.transition-terminal-success`. At most one Obligation of each type may
guard `issue.transition-terminal-success`; every Obligation type may guard
`envelope.complete`. Duplicate IDs, duplicate `review_profile` values,
duplicate `outcome_class` values, a type/field mismatch, more than one
`gate-evidence` Obligation, or an Obligation whose subject is not the primary
Issue reject the whole request with no Event.

`artifact_policy` is `required-available` for `gate-evidence` and
`metadata-sufficient` for the two new types. `freshness_target` is
`repository-tree` for `gate-evidence` and `review`, and `none` for
`human-outcome`: a Human judgment is not staled by a repository tree change,
while review proof is (section 7).

Obligation evaluation stays the single shared evaluator. It now returns one
`Evaluation` per Obligation, in Obligation ID order, each retaining the three
independent dimensions - absence, record validity/freshness, Artifact
availability - and all unsatisfied reasons. Absence keeps a per-type reason code
(DEC-228): a `gate-evidence` or `review` Obligation with no qualifying Evidence
keeps the existing `obligation-evidence-absent`, a `human-outcome` Obligation
with no linked Outcome returns `human-outcome-absent`, an unresolved blocking
Finding returns `finding-unresolved` (section 5), a Review whose independence no
longer holds at link time returns `review-not-independent`, and an Outcome
displaced under section 7 returns `outcome-superseded`. Each new code takes a
rank in the existing evaluation reason precedence order immediately after
`obligation-evidence-absent`, in the order `human-outcome-absent`,
`finding-unresolved`, `review-not-independent`, `outcome-superseded`. Only the
relative order is accepted (I-18): the unexported `precedence` map in
`internal/workflow/evaluation.go` is an integer map in which several codes
already share a rank, so renumbering its values to insert these four is
permitted and expected. What must not change is the relative order of the
Phase 21/22 codes among themselves. A guarded action requires every
applicable Obligation to be `satisfied`; a single unsatisfied Obligation blocks
it and every unsatisfied reason from every Obligation is returned. Evidence
absence, invalidity, and staleness never collapse into one reason code.

## 3. Review Runs And Independence

A Review is recorded by an Actor holding an active Claim of the new
non-mutating mode `review` and an active Run under that Claim. Phase 23 retains
exactly one active Claim per Issue across both modes; mode-aware concurrency is
Phase 24. A `review` Claim therefore requires no active mutating Claim on that
Issue, and a mutating Claim requires no active `review` Claim.

A `review` Claim is non-mutating: it cannot authorize repository adoption. A
`run.recover` or `run.abandon` request that would adopt any dirty path under a
`review` Claim rejects `claim-mode-invalid` with no Event. Reviewer Runs may
execute registered Gates, close, abandon, and be recovered like any Run; their
Receipts remain ordinary trusted execution facts, not Review Records.

Independence is a durable, computed property, not an attribution string. The
recording Actor of a Review Record must not equal any Actor that held a
mutating Claim on the reviewed Issue within the same Envelope revision, or a
Run under such a Claim. A non-independent request rejects
`review-not-independent` with the exact conflicting Claim or Run identity and
creates no Event. This rule, not a role name, is what makes "a Builder cannot
resolve its own Finding" enforceable (PRD-REV-004, TRD invariant).

**Evidence boundary (DEC-239).** The accepted Run lifecycle binds the Run to
the Claim's Actor and retains Claim history. A Run-only conflicting Actor
therefore cannot be manufactured through the CLI without changing that
lifecycle. VAL-015 retains real CLI evidence for the mutating-Claim rejection
and same-Reviewer rework, and an explicitly isolated durable-history query
fixture proves the Run-holder branch. That fixture proves the query's exact
conflicting Run and Actor; it does not claim the history is CLI-reachable or
authorize importing it. Neither independence rule nor lifecycle is changed.

**The conflicting set is the mutating side only (DEC-225).** A `review` Claim
and any Run under it never conflict, so the same independent Actor may record
successive Reviews across rework cycles on the same Issue and Envelope revision,
and may resolve a Finding created by its own earlier Review Record. Only a
Review Record resolving a Finding **it** created is prohibited (section 5); that
prohibition is per Review Record, never per Actor. The opposite reading - any
Run, including a review Run - would make the second Review of any rework cycle
impossible for the recording Reviewer, require a fresh Actor per cycle, and
leave the section 6 review-cycle and exhaustion model unusable, so it is not the
accepted reading. Conversely, an Actor that reviewed an Issue is not thereby
barred from later taking a mutating Claim on it; independence is evaluated at
Review record time against the durable history that exists then, and a Review
already recorded is never retroactively invalidated. Linking such a Review as
Evidence is re-evaluated at link time (section 13).

## 4. Review Records, Coverage, And Verdicts

A Review Record is immutable and binds Review ID, review Run, Issue, Envelope
and revision, the exact observed repository identity, `review_profile`, a
coverage declaration, a verdict, evidence references, and its creating Event.
The bound repository identity is re-observed by fid at record time and must
match the supplied value; a mismatch rejects `repository-mismatch`.

Coverage is a closed structure: a sorted unique non-empty list of covered
subject references, each naming a `Source` plus an `assessed` flag, and a
sorted unique list of declared exclusions with reasons. Coverage is a
declaration of what the Reviewer examined. fid stores and reports it; fid never
asserts that a declared coverage claim is true. Review evidence references are
`Source` values and optional existing Receipt IDs produced by that Review Run;
a referenced Receipt must belong to the recording Run.

The verdict is exactly one of `pass`, `changes-required`, or `blocked`.

- `pass` may satisfy a `review` Obligation when independence holds, the Review
  binds the current target, and every other Obligation condition holds. A
  `pass` Review may carry zero Findings; a `pass` Review may not carry a
  blocking Finding, which rejects `review-verdict-invalid`.
- `changes-required` leaves the Obligation unsatisfied, requires at least one
  blocking Finding, and derives Builder rework on the original Issue and
  Envelope revision. It charges exactly one review cycle (section 6).
- `blocked` records that the Review could not be completed. It leaves the
  Obligation unsatisfied, requires at least one non-blocking Finding or a
  declared exclusion, charges no review cycle, and raises no rework.

A Review Record itself is never a mutation of Issue state. Issue transition and
Envelope completion keep their existing separate guards and gain no implicit
success flag from a `pass` verdict; the verdict becomes actionable only when an
`evidence.link` binds it to a `review` Obligation (section 13).

## 5. Findings And Dispositions

Findings are created only as part of a Review Record, in one atomic Event, and
are immutable. A Finding binds Finding ID, its Review, a subject reference, a
scope, a severity from `blocker`, `major`, `minor`, `nit`, a summary, and a
blocking flag. `blocker` and `major` are blocking; `minor` and `nit` are not.
Severity and the blocking flag must agree, or the record rejects.

A Finding has exactly one of three terminal dispositions, each an append-only
row with its own Event:

- `addressed` - recorded by an Actor holding an active mutating Claim on the
  reviewed Issue. It is an assertion of rework, never a resolution. An
  `addressed` Finding stays unresolved and keeps blocking.
- `resolved` - recorded only as part of a later qualifying independent Review
  Record that explicitly lists the Finding in `resolves_finding_ids`. The
  resolving Review must bind the same Issue and Envelope revision, satisfy
  section 3 independence, and have verdict `pass` or `changes-required`. A
  Review may not resolve a Finding it created, and may not resolve a Finding
  that is already resolved or inapplicable.
- `inapplicable` - recorded only as part of a later qualifying independent
  Review Record, under the same rules as `resolved`, with a mandatory reason.

A Builder-recorded `addressed` disposition may be superseded by a later
`resolved` or `inapplicable` disposition. No other transition exists. A Finding
that is unresolved and blocking keeps its `review` Obligation unsatisfied with
reason `finding-unresolved`, regardless of the latest verdict.

## 6. Review Cycles, Rework, And Exhaustion

Budget authoring gains exactly one unit, `review-cycles`, and exactly one
additional `exhaustion_action`, `review-rework-exhausted`. All other Phase 22
budget rules are retained unchanged: immutable policy, null caller-supplied
`consumed`/`remaining`, envelope scope, unique names, positive int64 maximum,
and warning threshold null or 0..maximum.

`review-cycles` is a fid-owned unit and starts at known zero. fid charges
exactly one `review-cycles` delta, atomically with the Review Record Event,
when and only when a Review Record is recorded with verdict `changes-required`.
The source sample identity is the Review ID, so replay, recovery, and import
cannot charge the same Review twice. `run.observe-budget` rejects
`review-cycles` exactly as it rejects the other fid-owned units. No
`review-cycles` charge is derived from a Finding count, a rework Claim, or a
Run; one `changes-required` Review is one cycle.

Admission of the last allowed cycle is atomic with charging it and with raising
its exhaustion Stop and Attention. That `changes-required` Review still commits;
the Stop blocks the next rework instead of voiding the Review. The Stop code is
`review-rework-exhausted`, subject is the Envelope, resolver capability is
`stop.resolve`, blocked actions are exactly `claim.acquire` and `run.start`,
and its resolution condition is `budget-observed` on that budget name. A
subsequent `claim.acquire` or `run.start` on that Envelope returns
`review-rework-exhausted` and creates no Claim or Run. That refusal is a
non-committing denial (DEC-226): it appends no Event, stores no idempotency
result, and returns null `EventRef`, `replayed:false`, and the immutable Stop
and Attention IDs raised by the exhausting Review.

A budget whose `exhaustion_action` is `review-rework-exhausted` must have unit
`review-cycles`, and a `review-cycles` budget must use that exhaustion action.
An Envelope with a `review` Obligation and no `review-cycles` budget is valid;
unlimited rework is a deliberate, recorded policy choice, not an error. An
Envelope with a `review-cycles` budget and no `review` Obligation rejects.

No automatic retry, respawn, or rework Claim is created by exhaustion. Rework
consumes the existing Envelope budgets; nothing is reset by a new cycle.

## 7. Freshness And Staling Of Review Proof

Review proof is target-bound. A Review Record binds the repository tree
identity observed at record time. Evidence linked from that Review is evaluated
against the current target exactly like Receipt evidence: if the current target
tree differs from the bound tree, the review Evidence is `stale` and its
Obligation is unsatisfied with reason `evidence-stale` (PRD-REV-005). Staling
never deletes, rewrites, or invalidates the Review Record; identity, coverage,
verdict, and Findings are retained verbatim.

`human-outcome` Obligations use `freshness_target:"none"`. Human Outcome
evidence is never staled by a tree change; it is invalidated only by an
explicit later Human Outcome on the same `outcome_class` with a different
result, which supersedes the earlier one for evaluation while retaining both
records.

The three evidence dimensions stay independently serialized for all three
evidence types. A `review` or `human-outcome` Evidence record has
`artifact_availability:"not-applicable"` under `metadata-sufficient`; this is
a distinct third value, never silently reported as `available`.

## 8. Concerns And Deviations

Concerns replace the V1 completion-linked `concerns` table. A Concern binds a
generic subject drawn from the closed set `issue`, `envelope`, `run`, `review`,
`finding`, `requirement` - named `JudgmentSubject` in `public-interface.md`
section 3, and distinct from the Stop/Attention `Subject` (DEC-232); a
recording Actor; a scope; a severity from the
Finding severity set; a summary; concise reasoning; ordered references; a
status of `open` or `disposed`; and a created time. A Concern is not linked to
a completion report, a role, or a triage event.

A Concern has exactly three terminal dispositions, recorded by
`concern.dispose` with a mandatory reason:

| Action | Binds | Meaning |
|---|---|---|
| `carried-forward` | the exact existing Issue ID | The work survives as an Issue. |
| `rejected` | reason only | Consciously dismissed. |
| `superseded` | the exact existing Concern or Finding ID | Another record now owns the subject. |

`carried-forward` is the only way a Concern becomes future work, because
ADR-0004/0005 makes Issues the only durable work graph. The named Issue must
already exist; `concern.dispose` never creates one, and the new Issue need not
belong to the current Envelope's selection. A caller that also wants the
provenance edge in the graph records it explicitly with the existing
`discovered_from`/`discovered` Issue link relations; disposal writes no link.

`superseded` is the only relationship between a Concern and a Finding, and it
points from the Concern to the Finding. A Reviewer who judges an open Concern to
be genuine in-scope rework raises a Finding inside their own Review Record under
section 5, and the Concern is then disposed `superseded` naming it.
`concern.dispose` can never create a Finding: Findings are created only with a
Review Record, and a Finding is rework on the current revision rather than
carried-forward work.

Disposition is a single irreversible transition from `open` to `disposed`; a
second disposal rejects `concern-already-disposed`. There is no reopening.

Disposal is operational bookkeeping, not judgment. The disposing Actor may be
the Actor that raised the Concern: a Builder who later establishes that its own
observation was not a problem is the ordinary case, and Concerns are
deliberately off the critical path. Both the raising and the disposing Actor are
stored and shown in the Human views, so a self-dismissal is legible rather than
prohibited. Where a rejection is an authorized risk acceptance rather than
bookkeeping, the Human Outcome is recorded separately under section 10 and
referenced from the disposal; a Concern disposition never substitutes for
authorized judgment and is never itself an authority record.

An open Concern never blocks a guarded action, at any severity. It is surfaced
in the Issue view, `resume`, and `concern.list`, and it is never a precondition
of `envelope.complete` or `issue.transition-terminal-success`. Blocking work is
expressed as a blocking Finding or a Stop, never as a deferred Concern
(PRD-REV-008).

Deviations are append-only and have no lifecycle. A Deviation binds a generic
subject from the same closed set, the expected contract or control, the
observed departure, the authority that accepted or rejected it, its disposition
from `accepted`, `rejected`, `pending-decision`, ordered references, concise
reasoning, and a created time. A Deviation is never edited; a changed judgment
is a new Deviation that references the earlier one.

A Concern or Deviation `subject` is a typed reference that fid records without
requiring it to resolve to an existing row (DEC-228). Concerns and Deviations are
lightweight observations that must remain recordable about a subject fid does not
own or does not yet store - a `requirement` subject cannot resolve before task
12, and an observation about a removed or external subject stays legible. The
recorded `subject_type` is still closed and validated. The exception is a
disposition that binds an identity: `carried-forward` names an Issue and
`superseded` names a Concern or Finding, and those must already exist and are
enforced by foreign key. Task 8 records this interpretation in its completion
report.

Neither a Concern nor a Deviation is Evidence and neither can satisfy an
Obligation.

## 9. Decision References And Protected Decisions

A Decision Reference is an operational pointer, never a decision body. It binds
Decision Reference ID, a `class`, a `Source` with repository identity and
digest that locates the canonical project prose, the recording authority, an
`outcome` from `approved`, `rejected`, `deferred`, a scope, and a created time.
Recording a Decision Reference whose `Source` digest does not match the
observed project file rejects `source-digest-mismatch` with no Event.

An Envelope revision may declare 0..16 protected Decisions. Each binds an
ordinal, a `decision_class`, the required authority selector, and the affected
actions from the guarded-action set. While a protected Decision is
unsatisfied, every affected action is blocked and evaluation returns reason
`protected-decision-required`.

A protected Decision is satisfied only by both of the following, bound to the
same `decision_class` and within the Envelope's scope:

1. a Decision Reference with outcome `approved` whose recording authority
   matches the declared required authority; and
2. a Human Outcome with class equal to that `decision_class` and result
   `approved`, recorded under active Kit Activation provenance (PRD-AUTH-005,
   PRD-AUTH-008).

A Decision Reference with outcome `rejected` or `deferred` never satisfies a
protected Decision. Envelope activation with an unsatisfied protected Decision
is permitted; execution of the affected actions is not. Activation raises one
`protected-decision-required` Stop per unsatisfied protected Decision, subject
Envelope, blocked actions exactly the declared affected actions, resolver
capability `stop.resolve`, and resolution condition `decision-satisfied` naming
that Decision class. Recording the satisfying pair resolves that Stop and its
Attention in the same Event.

`linked_to_decision` becomes an implemented workflow transition rule evaluated
against Decision References in the Issue's scope, replacing its reserved
placeholder state.

## 10. Human Outcomes

A Human Outcome records authorized Human judgment. It binds Outcome ID, the
recording Actor, that Actor's activation provenance, a `class`, a subject from
the closed subject set, a `result` from `approved`, `rejected`,
`changes-required`, `acknowledged`, a scope, ordered references, concise
reasoning, and a created time.

A Human Outcome requires an Actor whose `actor_type` is `human` in the selected
finalized Activation profile and the `human-outcome.record` capability. Kit-mode
Human Outcomes require provenance authorized by the active Kit Activation; a
standalone Issue-inventory attribution cannot create one (PRD-AUTH-007/008). A
request from a non-human actor type, an unfinalized activation, or a missing
capability rejects `authority-missing` with no Event.

A Human Outcome may satisfy only a `human-outcome` Obligation whose
`outcome_class` equals its `class`, and only with result `approved`. It can
never satisfy a `gate-evidence` or `review` Obligation: an `evidence.link`
attempting that rejects `outcome-not-permitted` with no Event, naming the
Obligation type. A Human Outcome is never presented, stored, or serialized as
deterministic proof; its Evidence record retains
`producer:"human"`, `trust_source:"activated-human"`, and
`artifact_availability:"not-applicable"` (PRD-HUM-001).

An Envelope whose `human-outcome` Obligation is unsatisfied when every other
applicable Obligation is satisfied raises one `human-outcome-required` Stop with
subject Envelope, blocked actions exactly that Obligation's guarded actions,
resolver capability `stop.resolve`, and resolution condition
`human-outcome-recorded` naming the `outcome_class`. This is the "readiness for
Human validation" Attention permitted by PRD-HUM-002; routine status never
raises Attention.

**Raising producers are closed and named (DEC-224).** A Stop is durable state
and must be raised by a committing Event, so this readiness condition is
evaluated atomically inside exactly the four Events whose commit can make it
true, and nowhere else:

| Raising Event | Case it covers | Effect rows |
|---|---|---|
| `envelope.activated` | activation where the `human-outcome` Obligation is already the only unsatisfied one | Stop rows added for protected Decisions by task 9; reused here |
| `evidence.linked` | the ordinary case: linking Receipt or Review Evidence satisfies the last non-Human Obligation | optional `stops/insert`, `attention_items/insert`, `attention_stop_links/insert` added by task 10 |
| `review.recorded` | a later independent Review resolves the last blocking Finding, satisfying an already-evidenced `review` Obligation | already optional in `schema.md` section 6.2 |
| `artifact.availability-observed` | an observation flips the last `required-available` Artifact to available | optional `stops/insert`, `attention_items/insert`, `attention_stop_links/insert` added by task 10 |

At most one unresolved `human-outcome-required` Stop exists per Obligation: a
later qualifying Event re-evaluates and raises nothing while that Stop is open.
The Stop is never withdrawn by a later state change - a tree change that stales
Review evidence leaves it open - and resolves only through a qualifying Human
Outcome under section 11. No other Event raises it, and a guarded-action denial
never raises it, because a denial commits nothing (DEC-226).

While that Stop is open, the actions it blocks return `human-outcome-required`
rather than the generic `obligation-unsatisfied`, under the precedence in
`public-interface.md` section 5.

## 11. Attention Resolution Through Human Outcomes

`attention.v1` limits become `{"state":true,"human_outcome_resolution":true}`.

The `Condition` value gains three types: `review-recorded`, `finding-resolved`,
and `human-outcome-recorded`, plus `decision-satisfied` for section 9. Each new
condition names its subject; `human-outcome-recorded` additionally names a
non-null `outcome_class`, `decision-satisfied` a non-null `decision_class`, and
`finding-resolved` a non-null `finding_id`. Only the matching field is non-null.
`review-recorded` names **no** discriminator field: it is the one type of the
four with none, resolved by a `review.recorded` Event scoped by its subject, and
a `review-recorded` Condition carrying any of the three fields is rejected. The
enumeration above is exhaustive and no field serves two types (stated explicitly
under DEC-236; the accepted rule is unchanged).

Stop resolution continues to require a durable qualifying fact Event, never a
caller assertion. A Stop whose resolution condition is `human-outcome-recorded`
resolves only with the Event ID of a Human Outcome whose class and result match
that condition. `attention.resolve` keeps its exact existing shape: the sorted
unique set of the containing disposition Event IDs for all linked Stops.
Resolution references the qualifying Outcome; it never erases, edits, or
back-dates the original request (PRD-HUM-003).

Recording a qualifying Human Outcome does not silently close unrelated
Attention. Only Stops whose recorded condition matches are resolved.

## 12. Requirement References

V1 requirement authoring is replaced, not extended. `requirement_refs` binds a
stable requirement ID, a `Source` with repository identity and digest, a status
from `active`, `superseded`, `withdrawn`, and a created time.
`requirement_issue_links` binds requirement ID, Issue ID, and a relation from
`implements`, `verifies`, `relates-to`. Project files remain canonical; fid
stores no requirement prose, summary body, or parsed section text.

`fid requirement sync`, its markdown parser, its `last_synced_hash` column, and
its `requirement.synced` Event are removed. `requirement.record-reference` and
`requirement.link` replace them. A recorded reference whose `Source` digest does
not match the observed file rejects `source-digest-mismatch`. Re-recording the
same requirement ID with a different `Source` supersedes the prior row and
retains it; it never edits it in place.

`requirement_issue_links` becomes the single authoritative model for
requirement-to-Issue traceability, so `links` is restricted to Issue endpoints.
Its closed endpoint registry becomes `issue` only, and the relations
`implements` and `implemented_by` are removed from it; existing rows with a
requirement endpoint are deleted by migration 0018. This is required by the
slice acceptance rule in `10_IMPLEMENTATION_ROADMAP.md` section 7: no two
authoritative models may remain for one fact. `relates-to` therefore exists in
both tables with disjoint endpoint types and is not a duplicate model. The
`verifies` relation is new in `requirement_issue_links` and has no `links`
equivalent; it is accepted here and recorded in the schema catalog.

## 13. Evidence Across Three Evidence Types

`evidence_records` widens to exactly three evidence types with exactly one
non-null subject reference each:

| Evidence type | Subject column | Producer | Trust source |
|---|---|---|---|
| `receipt` | `receipt_id` | `fid` | `fid-executed` |
| `review` | `review_id` | `activated-actor` | `activated-actor` |
| `human-outcome` | `human_outcome_id` | `human` | `activated-human` |

`evidence.link` continues to append exactly one immutable record per
(Obligation, subject) pair and never overwrites prior proof. The linked
subject's type must be in the Obligation's `accepted_evidence_types` and its
trust source in `accepted_trust_sources`; a mismatch rejects
`evidence-type-not-accepted`. Linking a `changes-required` or `blocked` Review
rejects `review-changes-required` with the blocking Finding IDs. Linking a
Review that fails independence rejects `review-not-independent`; independence is
re-evaluated against the durable history at link time, so a Reviewer that has
since taken a mutating Claim on that Issue and revision cannot have its Review
linked. Linking a non-`approved` Human Outcome rejects `outcome-not-permitted`.

**Rejection precedence is fixed (DEC-226).** Exactly one code is returned per
rejected `evidence.link`, in this order: `outcome-not-permitted` whenever the
linked subject is a Human Outcome - including a Human Outcome offered to a
`gate-evidence` or `review` Obligation, and a non-`approved` Outcome offered to
a matching `human-outcome` Obligation; then `review-changes-required`; then
`review-not-independent`; then `evidence-type-not-accepted` for every remaining
type or trust-source mismatch. The Human Outcome case is deliberately the most
specific: `outcome-not-permitted` names the Obligation type, the Outcome, and
its result, which is what the caller must act on, while
`evidence-type-not-accepted` would only echo the accepted list. The nine-cell
type-by-type matrix therefore resolves to three accepted pairs, two
`outcome-not-permitted` rejections (a Human Outcome against `gate-evidence` and
against `review`), and four `evidence-type-not-accepted` rejections.

`evidence.evaluate` takes an optional Obligation ID. With an Obligation ID it
retains the Phase 21 single-Obligation result; without one it returns every
Obligation's evaluation for the active Envelope revision in Obligation ID
order. It remains pure, non-mutating, and requires no Actor or idempotency key.

## 14. Human-Facing Review And Judgment Projections

Every new leaf has purpose-labeled Human text derived from the same facts as
its JSON, adding no member and no rule (PRD-PROJ-005, DEC-219 pattern).

- Review text names the reviewed Issue and Envelope revision, the bound
  repository identity, the profile, the verdict, the counted Findings by
  severity, and the exact independence basis. It never says a `pass` verdict
  completed the Issue.
- Finding text names the Finding, its Review, its blocking state, and its
  current disposition, and says explicitly that `addressed` is a Builder
  assertion, not a resolution.
- Concern and Deviation text names subject, severity, status or disposition,
  and references, and never presents either as Evidence. Concern text names both
  the raising and the disposing Actor, and for a `carried-forward` disposition
  names the Issue that now owns the work.
- Decision Reference text names the class, the referenced project file path and
  digest, the authority, and the outcome, and never reproduces decision prose.
- Human Outcome text names the Actor, the activation provenance, the class, the
  result, and states that the Outcome is authorized Human judgment and not
  deterministic proof.
- Protected-Decision and review-rework Stop text names the blocking reason and
  the exact qualifying resolution.

`resume` replaces its `findings:{available:false,items:[]}` placeholder with
the real open-Finding projection, and adds bounded `reviews`, `concerns`,
`deviations`, `decisions`, and `human_outcomes` sections plus per-Obligation
evaluations. Its `next_actions` gain the newly executable review and judgment
operations with their reasons. The projection remains bounded, reproducible,
and non-editable; exceeding a bound returns `result-limit-exceeded` rather than
silently truncating judgment state.

Phase 23 applies the CF-040 Human-view guidance to the views it introduces:
labeled sections, an operation purpose line, and no narration of routine state.
The accepted polish and verbosity review of pre-existing Human views stays with
the Phase 25 operating-experience qualification.

## 15. Availability, Removal, And Acceptance

Final capability rows changed by this phase:

| Capability | Available | Exact limits | Trust source |
|---|---|---|---|
| `reviews.coverage.v1` | true | `{"record":true,"independence":"distinct-actor","findings":true}` | fid-journal |
| `human.outcomes.v1` | true | `{"record":true,"provenance":"active-kit-activation"}` | active-kit-activation |
| `attention.v1` | true | `{"state":true,"human_outcome_resolution":true}` | fid-journal |

`work.claims.v1` stays false with owner phase 24 and
`telemetry.trajectory.v1` false with owner phase 25. Schema presence alone does
not enable a row; the final matrix changes only after the replacement and
removal proofs pass.

Owned removals, each replacement-bound and each proven absent at phase end:
`task_completions`, `ac_verifications`, the V1 completion-linked `concerns`
shape, the V1 `decisions` table, the V1 `requirements` table and its
`requirement.synced` Event, `fid concern raise`, `fid concern address`,
`fid issue ac verify` and its role-keyed verification projection,
`fid requirement sync`, the `all_acs_verified` required-roles parameter, the
`has_completion_report` and `concerns_addressed` transition rules and their
evaluators, and the maintenance checks that read those tables. The exact file,
table, index, trigger, Event, fixture, and test inventory is
`v1-convergence-inventory.md`; the executing gate for every row is named there.

`acceptance_criteria` is retained with no verification role columns, matching
the target catalog. `triage_events`, `agent_runs`, `epic_metadata`, historical
stats, superseded migrations, and migration-scoped Event fixtures remain owned
by Phase 24 and Phase 25 and are explicitly not removed here.

This contract is accepted together with `public-interface.md` and `schema.md`.
Accepting it also accepts these canonical spec amendments, which the owning
tasks apply with their implementation:

- `03_TRD_SCHEMA.md` gains section 8.4 incorporating `schema.md`.
- `03_TRD_SCHEMA.md` section 3.2 records that the `links` closed endpoint
  registry is `issue` only, and that `requirement_issue_links.relation` is the
  closed set `implements`, `verifies`, `relates-to`.
- `03_TRD_SCHEMA.md` section 5 adds `Concern disposed` to the Review/Judgment
  Event families. PRD-REV-007 requires a Concern status and disposition, and
  Design section 4 shows open Concerns, but the canonical family list had no
  event that could change either. Adding the missing family rather than dropping
  the clause follows DEC-195.
- `06_DESIGN_PLAN_CLI.md` section 3 review-and-judgment block gains
  `fid concern dispose`, `fid deviation list`, `fid decision list`,
  `fid human-outcome list`, and replaces `fid requirement sync` with
  `fid requirement record-reference|link|get|list`.
- `06_DESIGN_PLAN_CLI.md` section 5 records the Phase 23 closed code additions.
- `02_TRD.md` section 4 marks `reviews.coverage.v1` and `human.outcomes.v1`
  qualified after task 17.
- `GLOSSARY.md` replaces the single `Reason code` entry with **Failure code**
  and **Evaluation reason** (DEC-226). Unlike the amendments above, the Tech
  Lead applied this one on 2026-09-19 rather than leaving it to an owning task,
  because task specs 6, 10, and 13 depend on the distinction to file each new
  identifier in the correct set.

Until acceptance, no product code, migration, or test for this slice may start.
