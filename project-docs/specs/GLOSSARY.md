# GLOSSARY - fid

This document defines canonical product language. It contains no interface,
storage, or implementation requirements.

## Product Boundary

**fid**
Definition: The local-first Issue tracker and durable operational work-state
and evidence backend.
Use when: naming the product that owns Issues and resumable execution facts.
Avoid: describing fid as the owner of engineering policy, agents, project
prose, source code, or Human judgment.

**Operational work state**
Definition: Durable facts required to know what work exists, what is currently
authorized, what happened, what proves it, and what needs attention.
Use when: distinguishing fid state from project intent or kit policy.
Avoid: unqualified “source of truth.”

**Project-owned prose**
Definition: Requirements, architecture, validation strategy, decisions, and
other durable Human-authored intent whose canonical form remains in project
files and Git.
Use when: content may be referenced by fid but must not be copied into its
operational model.
Avoid: “fid document” or “backend prose.”

**Routine (external kit term)**
Definition: An external kit risk-lane classification for work the kit judges
eligible for its routine policy path. It is not a fid entity or lifecycle.
Use when: discussing how the external kit selects and compiles policy before
activation.
Avoid: persisting a Routine, evaluating a Routine in fid, naming a fid command
or Event with it, or using it as a synonym for Work Envelope.

**Capability availability**
Definition: A discovery claim that one independently versioned semantic
capability is complete within its declared limits and trust sources.
Use when: a consumer decides whether it can rely on that capability.
Avoid: “partial availability”; an incomplete semantic capability is
unavailable even when internal walking-skeleton operations exist.

**Public enum value**
Definition: A member of a closed, contract-controlled value set serialized on
a public machine interface, including reason, outcome, control, status,
availability, and Event-kind values. Public enum values use lowercase
kebab-case.
Use when: distinguishing closed values from their JSON field names.
Avoid: classifying opaque entity IDs, user-authored workflow-state identifiers,
paths, capability IDs, or digests as enums merely because they are strings.

**Project-local attribution**
Definition: A typed `<agent|human|service>:<id>` caller label recorded on a
retained standalone inventory fact to show who invoked it, with producer
`fid-standalone` and provenance `project-local-attribution`, without
establishing kit Actor capability or trusted provenance.
Use when: converting retained V1 `--by` mutations to `--actor` during Phase 21.
Avoid: using a role-looking attribution to authorize a protected operation,
satisfy Evidence or a Human Outcome, or permit guarded transition/completion.

## Work Inventory

**Issue**
Definition: A durable epic, task, or subtask in fid. Issues and their
relationships are the only durable work graph.
Use when: referring to any tracked unit of work.
Avoid: “ticket” or persistent “Work Item.”

**Epic**
Definition: A top-level Issue grouping related tasks by explicit sequence or
priority.
Use when: work needs durable grouping and scheduling.
Avoid: “phase” when discussing the fid entity.

**Task**
Definition: An independently claimable and acceptable unit of work, optionally
under an Epic.
Use when: the unit has its own durable lifecycle.
Avoid: promoting transient procedural steps to Tasks.

**Subtask**
Definition: An independently meaningful child of a Task.
Use when: decomposition must survive a Run or permit separate claim,
dependency, acceptance, or parallel execution.
Avoid: using Subtasks only as a checklist.

**Issue workflow**
Definition: A named state machine governing one or more Issue kinds.
Use when: describing Issue lifecycle states and legal transitions.
Avoid: using it for kit risk, role-routing, or assurance policy.

**Terminal-success Issue state**
Definition: A workflow state in the completed category that qualifies as the
successful terminal outcome for the selected Issue rather than cancellation.
Use when: applying the gate-evidence guard before Phase 21 completion.
Avoid: treating a legacy completion row, role-specific verification, or
caller attribution as the terminal-success state or its evidence.

**Acceptance criterion**
Definition: An Issue-level statement of an observable required outcome.
Use when: defining what completion must achieve.
Avoid: treating an unchecked assertion as evidence or duplicating it as a
separate durable validation hierarchy.

## Execution Contract

**Work Envelope**
Definition: The approved, versioned execution contract over one primary Issue
and an explicitly bounded related Issue set. It records outcome, source
references, repository scope, risk result, authority, protected Decisions,
Obligations, budgets, stop conditions, and selection.
Use when: naming what execution is approved.
Avoid: “handoff,” “task spec,” or a second work graph.

**Envelope Revision**
Definition: An immutable approved version of a Work Envelope. Material changes
to outcome, selection, scope, risk, authority, acceptance, Obligations, budgets,
or stop conditions create a new revision.
Use when: binding Claims, Runs, Receipts, and Reviews to an exact contract.
Avoid: editing an active revision in place.

**Obligation**
Definition: A typed condition that must be satisfied for an applicable action
or outcome, including gate, evidence, review, acceptance, and Human judgment
requirements.
Use when: policy has been compiled into an enforceable condition.
Avoid: “QA step,” “validation assertion,” or an untyped checklist item.
Applicable: An Obligation is applicable to a guarded action exactly when it is
`required` and its guarded actions include that action. It never blocks an
action it does not guard (Phase 23 DEC-249).

**Actor**
Definition: A Human, agent, harness, or trusted producer identity recorded with
capabilities and provenance established by the applicable trusted source.
Use when: attributing an attempted or completed action.
Avoid: assuming a kit role name is a fid-global role or treating
caller-self-declared capabilities as authority.

**Authority Grant**
Definition: The explicit actions and protected-decision classes an Actor may
exercise under an exact Work Envelope revision.
Use when: deciding whether an operation is authorized.
Avoid: implied permission or recalculating kit role policy.

**Claim**
Definition: A versioned lease granting an Actor temporary authority over a
specific subject under an exact Envelope revision.
Use when: coordinating execution and rejecting stale mutation.
Avoid: ownership, assignment, or indefinite locking.

**Mutation Claim**
Definition: The exclusive Claim permitting one Builder to mutate one selected
Issue in one worktree and declared mutation scope.
Use when: enforcing parallel Builder isolation.
Avoid: claiming a whole Envelope or allowing overlapping/unknown scope.

**Run**
Definition: One bounded execution attempt by an Actor, linked to its Issue,
Envelope Revision, Claim, repository identity, budgets, and optional
Obligation.
Use when: recording execution identity, telemetry, and terminal outcome.
Avoid: reusing a Run identity after interruption.

**Interrupted Run**
Definition: A terminal Run whose executor ended without normal completion or
intentional abandonment. Continuation creates a linked new Run after liveness
and repository reconciliation.
Use when: recovering an interrupted session.
Avoid: resuming by mutating or reopening the old Run.

**Next Action**
Definition: A noncanonical projection of the action currently permitted or
required, its subject, required Actor capabilities, and typed reasons.
Use when: selecting or resuming work.
Avoid: persisting it as a work entity.

**Partial Progress**
Definition: A projection of the current states of Issues selected by an active
Envelope Revision.
Use when: showing multi-Issue progress.
Avoid: a partially-complete Envelope or membership lifecycle.

## Execution Evidence

**Gate Definition**
Definition: A project-approved, versioned machine value defining one executable
deterministic check and its execution restrictions.
Use when: authoring gate policy in project configuration.
Avoid: free-form command prose or runtime authority.

**Gate Definition Registration**
Definition: fid's immutable operational authorization binding an approved Gate
Definition and digest to project, overlay, repository, scope, provenance, and
registration version.
Use when: naming what the trusted executor may resolve and run.
Avoid: treating mutable authoring configuration as the live registry.

**Execution Attempt**
Definition: The immutable fact that trusted execution started before process
spawn.
Use when: distinguishing a started execution from its terminal result.
Avoid: inferring success from an Attempt without a Receipt.

**Execution Receipt**
Definition: The immutable terminal observation for one Execution Attempt,
bound to the exact registration, producer, subject, Run, repository identity,
invocation, timing, outcome, toolchain, and sanitized output identity.
Use when: proving deterministic execution.
Avoid: caller-reported exit codes or mutable result rows.

**Trust Source**
Definition: The provenance class of evidence: `fid-executed`, `attested`, or
`reported`.
Use when: evaluating whether evidence may satisfy an Obligation.
Avoid: treating `reported` as deterministic proof.

**Evidence Record**
Definition: A typed, append-only record linking an Obligation to a qualifying
Receipt, Review, Human Outcome, or external artifact identity.
Use when: evaluating satisfaction and freshness.
Avoid: copied output, prose-only claims, or a mutable “passed” flag.

**Referenced Artifact**
Definition: Full material addressed by an Evidence Record, such as a retained
log, screenshot, export, or build artifact. Availability is independent of the
record's validity.
Use when: an Obligation may require retained material rather than metadata.
Avoid: treating artifact loss as deletion of the Evidence Record.

**Freshness**
Definition: Whether evidence remains applicable to the current target identity
under its Obligation policy.
Use when: repository or contract changes may stale prior proof.
Avoid: conflating freshness, Evidence Record existence, and Artifact
availability.

## Review And Judgment

**Review Record**
Definition: The immutable outcome of one independent Review Run, bound to the
exact Envelope Revision and repository identity, with profile, coverage,
verdict, Findings, concerns, and evidence references.
Use when: proving that required review occurred.
Avoid: a QA Work Item or mutable review document.

**Review Finding**
Definition: A structured defect or improvement recorded by a Review Record. A
Builder may mark it addressed, a nonterminal assertion; only a later qualifying
independent Review Record resolves it or marks it inapplicable with a reason.
Use when: review discovers actionable work.
Avoid: free-form comment streams or Builder self-resolution.

**Concern**
Definition: A concise execution-time risk, defect, uncertainty, or mismatch
that can affect action or completion outside a Review Record.
Use when: the fact must remain visible and govern work.
Avoid: routine narration.

**Deviation**
Definition: A structured departure from an approved contract, expected path,
or required control, with scope, reason, authority, and disposition.
Use when: the departure affects evidence, risk, or acceptance.
Avoid: silently adapting the contract.

**Protected Decision**
Definition: A product, architecture, data, interface, dependency, security,
privacy, release, destructive, external, or project-declared choice requiring
Human authority.
Use when: work must stop pending authorized judgment.
Avoid: routine reversible implementation choices.

**Decision Reference**
Definition: The operational identity, authority, outcome, scope, and repository
reference of a project-owned protected Decision.
Use when: the decision governs fid state.
Avoid: copying decision prose into fid.

**Human Outcome**
Definition: A structured Human approval, rejection, checklist result, decision,
or risk acceptance.
Use when: an Obligation explicitly permits or requires Human judgment.
Avoid: presenting Human judgment as deterministic execution proof.

**Stop**
Definition: A typed operational fact that prevents a specified action until its
recorded resolution condition is met.
Use when: continuation would violate policy or integrity.
Avoid: ordinary status updates or generic escalation prose.

**Attention Item**
Definition: A structured, deduplicated request for Human or authorized-agent
attention caused by a decision need, circuit breaker, blocker, risk acceptance,
or readiness for validation.
Use when: intervention is genuinely required.
Avoid: routine progress notifications.

## Recovery And Projection

Phase 22 terminology clarification (DEC-180/184/185; accepted under DEC-188):
an executor liveness observation binds one Run/session identity and records
`terminated`, `live`, or `unknown` from an authorized observer. Lease expiry is
not such an observation. Unknown always blocks takeover; a live takeover needs
an exact current Human force grant and does not assert process termination.
Intentional Run abandonment is distinct from interruption and produces no
trusted execution proof. A budget sample is an append-only observed measurement
or fid-derived charge, not an editable counter. Detailed accepted value schemas
are in `../docs/phases/phase-22/public-interface.md`.

**Dirty Digest (Phase 22, DEC-185/188)**
Definition: The versioned SHA-256 of canonical changed-path HEAD/index/worktree
images, as specified in Phase 22 contracts section 3.1. It excludes operational
metadata, branch names, and capture time; clean state uses null.
Use when: comparing the same source target across caller inputs and observations.
Avoid: arbitrary caller hashes, raw Git status text, or the reconciliation digest.

**Reconciliation Digest (Phase 22, DEC-185/188)**
Definition: The canonical hash of the entire reconciliation apart from its own
digest member, binding repository, checkpoint, paths, classification and safety
state under Phase 22 contracts section 3.2.
Use when: binding exact checkpoint/adoption/recovery requests to observed state.
Avoid: substituting it for a dirty digest or for executor liveness proof.

**Stop Disposition Event (Phase 22, DEC-184/188)**
Definition: The Event that records a Stop resolution, distinct from the durable
fact used to justify it. One activation recovery Event may dispose several Stops.
Use when: naming the exact Event set required by Attention resolution.
Avoid: passing only the underlying observation/fact IDs to resolve Attention.

**Recovery Adoption**
Definition: Retaining explicitly reconciled uncommitted changes for a linked
successor Run under the recorded recovery authority (Phase 22, DEC-219/220).
Use when: describing which existing repository changes the successor retains.
Avoid: treating adoption as testing, Review, approval, correctness proof, or an
executor liveness observation.

**Kit Activation**
Definition: The durable cross-resource state that aligns one kit version,
overlay, provider projections, and backend migration before normal execution.
Use when: installing, upgrading, rolling back, or recovering the kit/backend
contract.
Avoid: partial activation or independent component versions.

**Event**
Definition: An immutable, ordered statement that a durable workflow fact
occurred.
Use when: rebuilding and auditing state.
Avoid: current projections or mutable rows as history.

**Projection**
Definition: A reproducible view derived from durable facts for a specific
machine or Human question.
Use when: presenting Next Action, progress, context, completion, Attention, or
economics.
Avoid: “sync,” editable generated files, or another source of truth.

**Human Validation Brief**
Definition: A decision-focused projection of delivered outcomes, material
decisions, evidence, deviations, unresolved risk, economics, and remaining
Human-only checks.
Use when: requesting Human judgment.
Avoid: full transcripts or duplicated completion reports.

**Portable Recovery Export**
Definition: A versioned backend-neutral representation sufficient to restore
fid operational state without copying canonical project prose.
Use when: validating recovery, migration, or backend portability.
Avoid: Markdown export or direct database copying.

**Transient workflow prose**
Definition: Lifecycle narration reconstructible from project sources, events,
evidence, and projections, including progress reports, handoffs, command tables,
completion reports, and copied evidence roll-ups.
Use when: identifying prose that should not become durable product state.
Avoid: applying the term to requirements, decisions, Findings, Deviations, or
Human judgment.

## Temporary Repository Delivery Terms

These terms remain canonical for the V1 process used to build fid. They are
not forward fid product entities and do not contradict the product-term
retirements below.

**Operating Mode**
Definition: The per-phase declaration - light, normal, or heavy - selecting
which V1 delivery artifacts a phase must produce. It changes the artifact set,
never exact-command evidence rules, commit discipline, or hard stops.
Use when: planning a fid implementation phase under `09_AI_WORKFLOW.md`.
Avoid: using the mode as a fid runtime state or V2 backend-port field.

**Validation Contract / Validation Assertion**
Definition: The current per-phase file of testable `VAL-NNN` assertions used to
connect requirements, Developer tasks, and QA evidence while building fid.
Use when: planning or validating a current fid implementation phase.
Avoid: modelling either as a durable fid product aggregate.

**Phase roles and artifacts**
Definition: Tech Lead, TL-Lite, Phase Runner, Developer, QA Lead, QA Reviewer,
phase briefs, task specs, QA plans, completion reports, triage logs, and
handoffs defined by the binding V1 repository workflow.
Use when: operating the current fid build pipeline.
Avoid: encoding them in fid's forward domain model or backend contract.

**RESUME (V1 disposition)**
Definition: A triage verdict that the trigger is cleared, demonstrably
misclassified, or already resolved with verified evidence, so orchestration may
advance from the recorded checkpoint.
Use when: no repair or protected decision remains.
Avoid: “advance now and repair later.”

**REPAIR (V1 disposition)**
Definition: A triage verdict that a real condition is local, in scope,
authorized, and decision-free, requiring a same-task Developer continuation
and verification before orchestration advances.
Use when: incomplete work, a gate/report defect, local recovery, or bounded QA
rework can be completed under the accepted task contract.
Avoid: closing the task/finding or requesting Human attention.

**Control-conformance repair (V1)**
Definition: A bounded recovery that fixes an implementation defect in a
canonical verification/workflow control when an already accepted contract
uniquely determines behavior, without changing the command, assertion
semantics, accepted status/scope, threshold/budget/tolerance, or policy.
Use when: TL-Lite can record an exact reproducer and tools-only allowlist, and
the Phase Runner can verify one isolated commit before the originating task
continues.
Avoid: using a failing control as generic authority to edit tests, product,
specs, prompts, templates, dependencies, or the gate/control definition.

**BLOCK (V1 disposition)**
Definition: A triage verdict that a protected decision, unsafe condition,
unrecoverable external dependency, specification contradiction, or exhausted
bounded recovery prevents progress.
Use when: TL-Lite or the full Tech Lead must allocate or reference an ESC and
request the named Human or external-authority action.
Avoid: using it as a synonym for red, incomplete, failed, or must-fix.

**Trigger signature (V1)**
Definition: The stable identity of one recovery condition, composed from the
owning task or review, recovery/blocker category, and durable source identity.
Use when: correlating triage, continuation attempts, checkpoints, and evidence.
Avoid: changing it merely because a session or orchestrator changed.

**Recovery attempt (V1)**
Definition: One Phase Runner-mediated Developer continuation for one Trigger
signature. Same-session red/green work is not an orchestration Recovery attempt.
Use when: recording total and consecutive no-progress attempts.
Avoid: resetting the count for context refresh or planned handoff.

**Orchestration session-liveness (V1)**
Definition: A child-agent execution condition in which the child is externally
interrupted before returning a typed result, distinct from a command timeout or
a parent polling wait timeout. The ledger records `Outcome: INTERRUPTED` with
the cause, elapsed session duration, and command in flight.
Use when: preserving the no-progress count after one interruption and triaging
the execution shape after two consecutive interruptions. A repair verdict
advances total attempt count while starting a new interruption streak.
Avoid: charging an interruption to repair capability, resetting an existing
no-progress count, or treating an expired parent wait as child termination.

**V1 recovery ledger event**
Definition: An append-only run-ledger entry: `TASK-RECOVERY` records a
classified recovery need and checkpoint; `REPAIR` records one continuation and
its resolution evidence; `RUN-RESUME` records runner-state reconciliation
before orchestration continues.
Use when: recovering the current repository delivery workflow from files.
Avoid: treating these Markdown entries as forward fid product **Events**.

**Protected blocker category (V1)**
Definition: The closed reason permitting `BLOCK`: `product-decision`,
`architecture-decision`, `data-model-decision`, `public-interface-decision`,
`dependency-approval`, `security-or-data-loss`, `external-infrastructure`,
`specification-contradiction`, or `recovery-budget-exhausted`.
Use when: Human or named external-authority attention is genuinely required.
Avoid: classifying ordinary implementation, verification, report, or QA repair.

## Machine Interface

**Response Envelope**
Definition: The versioned JSON wrapper every machine-mode fid operation
returns, carrying contract version, operation name, and either a result with
optional Event identity or a typed error.
Use when: describing the transport shape of a CLI or port response.
Avoid: the bare word "envelope," which means **Work Envelope**; and treating
the wrapper as a durable entity.

**Operation name**
Definition: The stable dotted identifier of one machine operation, such as
`issue.transition`, carried by both success and failure responses.
Use when: routing, logging, or asserting on a specific operation.
Avoid: deriving it from the command string or from a Human message.

**Failure code**
Definition: A stable machine identifier naming one refused operation, paired
with a closed category and a retryability flag. Exactly one accompanies a
failed response, and it decides that response's exit status.
Use when: a consumer must branch on why a request was refused.
Avoid: parsing Human message text; introducing a code outside the owning slice
contract; or using a failure code to describe the state of an Obligation that
was read successfully.

**Evaluation reason**
Definition: A stable machine identifier describing why one Obligation is
unsatisfied, carried with the subject it refers to. Zero or more accompany an
Obligation evaluation. An evaluation reason has no category, no retryability,
and no exit status of its own.
Use when: reporting or branching on the state of one Obligation, whether it was
read successfully or restated inside a refused response's details.
Avoid: treating an evaluation reason as a refusal. A successful read that
reports an unsatisfied Obligation is not a failure.

A **Failure code** and an **Evaluation reason** are separate closed sets that
answer different questions: "was this request refused, and why" against "what is
the state of this Obligation, and why". Both appear in machine responses, and
both are written as a `code` field, so they are distinguished by position, never
by shape. An unsatisfied Obligation read on its own carries evaluation reasons
in a successful response; attempting the action that Obligation guards returns a
single failure code, which may restate those same evaluation reasons in its
details.

An identifier may deliberately belong to both sets, because the same fact can
refuse an action and describe a state: `protected-decision-required` names the
refusal when the guarded action is attempted, and describes the Obligation state
when the Envelope is read. An identifier that belongs to only one set must not
be registered in the other. The bare term `Reason code` is retired: it named the
failure set here while the slice contracts used `ReasonCode` for the evaluation
set, so one term denoted two closed sets and an identifier was filed in the
wrong one (Phase 23 DEC-226). Name the set explicitly.

**Project root**
Definition: The directory that directly contains `.fid/` and identifies one fid
project.
Use when: describing discovery, explicit selection, or path normalization.
Avoid: "workspace," which historically named the same directory and now invites
confusion with a Git worktree.

## Retired Product Terms

`Work Item`, fid `Role`, `QA Plan`, `Validation Contract`, `Validation
Assertion`, `Triage Entry`, `Developer completion report`, and `handoff` are not
forward fid entities. They may appear in historical phase records, the current
V1 repository delivery process, or the external kit's temporary engineering
memory. New fid product specifications must use the retained terms above.

*Last updated: 2026-09-19 (Phase 23 DEC-226: `Reason code` split into
**Failure code** and **Evaluation reason**)*
