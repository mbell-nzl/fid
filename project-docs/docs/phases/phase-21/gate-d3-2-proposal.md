# Gate D3.2 Accepted Outcome - One-Issue Walking Skeleton

> Phase: 21 - One-Issue walking skeleton
> Author: Tech Lead
> Status: **ACCEPTED by Human direction on 2026-08-22**
> Decision: DEC-153
> Authority: ADR-0011, ADR-0012, and
> `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` sections 5 and 10.

This record closes Gate D3.2. It preserves the core proposal subject to every
Human amendment below. The exact implementable surface is fixed separately by
`contracts.md`; neither this record nor the phase title creates a persisted
fid `Routine` concept.

## Accepted Journey Boundary

One real-process journey, repeated after rebuild, must:

1. initialize fid and create one Issue through the Phase 20 standalone
   contract;
2. discover the exact capability matrix, including unavailable incomplete
   capabilities and their declared limits/trust sources;
3. prepare, stage, commit, and finalize one Kit Activation carrying the
   authority profile and approved Gate Definition values;
4. create an immutable Gate Definition Registration from those activated
   values;
5. create and activate one Work Envelope revision selecting the Issue and
   referencing that Registration through one gate-evidence Obligation;
6. resolve the Actor from the active activation, authorize it through the
   Envelope grant, acquire one Mutation Claim, and start one Run;
7. perform one nonterminal execution-bound Issue transition under active
   revision, Actor, Claim, expected Issue revision, and repository checks;
8. ask fid to execute only the Registration reference, append the Attempt
   before spawn, stream sanitized framed output, and record at most one
   immutable terminal Receipt after observed termination;
9. link the Receipt to the Obligation, evaluate absence/freshness/Artifact
   availability independently, and require `satisfied` before the
   terminal-success Issue transition;
10. perform the terminal-success Issue transition through the shared
    Obligation evaluator, then finish the Run and release the Claim;
11. complete the Work Envelope through the same evaluator, which rechecks all
    required Obligations and the terminal-success Issue;
12. reconstruct bounded resume state, create an export, and non-mutatingly
    validate it; and
13. delete the projection, rebuild from durable facts, then repeat only the
    named reads, evaluation, completion state, resume, export, and negative
    conformance checks. Do not create a second one-revision Envelope.

Registration precedes Work Envelope creation. An Envelope cannot reference a
Gate Definition value or mutable project file directly.

## DEC-156 Human Amendment - Terminal-Success Guard

The Human's 2026-08-22 planning correction amends, but does not rewrite,
DEC-153. The one `gate-evidence` Obligation guards both
`issue.transition-terminal-success` and `envelope.complete`. Gate execution,
Receipt linking, and successful Evidence evaluation occur before the guarded
Issue transition. Completion remains a separate later mutation and re-runs
the same Obligation evaluator against current Evidence facts; its separate
action predicate also requires the terminal-success Issue.
Nonterminal execution-bound transition remains earlier in the journey.

## Human Decision Point Disposition

| Point | Accepted disposition | Binding consequence |
|---|---|---|
| DP-1 Attention | **No; Phase 22.** | Routine unsatisfied Obligations are queryable status and never raise Attention. |
| DP-2 phase shape | **One Phase 21, Heavy.** | Two internal task groups; no 21A/21B archives or duplicated closure machinery. A fresh independent checkpoint separates the groups. |
| DP-3 executor | **Standard-library enforceable baseline.** | Enforce structured argv/no shell, exact cwd, environment allowlist, timeout/output ceilings, credential-name policy, redaction, and contained Artifact paths. Accept only explicit permissive network/write policy; reject stronger unsupported restrictions at registration. Claim no OS containment. |
| DP-4 budgets | **Validated immutable contracts only.** | Bind budgets to Envelope revision and Run. `consumed` and `remaining` are always JSON `null`; no consumption, warning, exhaustion, Stop, or Attention behavior. |
| DP-5 export | **Create and validate only.** | Import and identity/history preservation belong to Phase 22. `export.portable.v1` remains unavailable and PRD-REC-004 is not claimed. |
| DP-6 Obligation | **One complete generic value shape; one supported type.** | Only deterministic gate evidence from fid-executed Receipts can satisfy the supported `gate-evidence` type. |
| DP-7 retrospective | **Completed first.** | DEC-152 applies its task-size, feasibility, report-lint, and review-cadence decisions to the pack. |

## Exact Inclusion Boundary

- capability discovery with binary/migration-derived availability, limits, and
  accepted trust sources;
- partial activation operations needed for the happy path and pre-commit
  abort, while `kit.activation.v1` remains unavailable until recovery exists;
- activation-authorized Actor profiles and Gate Definition values;
- immutable Gate Definition Registration before Work Envelope creation;
- exactly one Work Envelope, one immutable revision, one primary Issue, zero
  related Issues, one gate-evidence Obligation, and immutable budget values;
- generic Authority Grants evaluated against the activation-resolved Actor;
- one active Mutation Claim, one Run, and explicit acquire/renew/release and
  start/finish reads;
- execution-bound nonterminal transition enforcement plus a shared
  gate-evidence evaluator guarding terminal success and completion;
- standard-library local execution, framed live output, redaction/control
  handling, Attempt/Receipt terminal semantics, and safe Artifact retention;
- Receipt-linked Evidence, conservative repository-tree freshness, Artifact
  availability observation, and deterministic all-reasons evaluation;
- Work Envelope completion and bounded resume projection;
- deterministic export creation and non-mutating validation;
- idempotency for every Phase 21 mutation and the retained retryable Issue
  create/transition surfaces; and
- live/rebuild parity, maintenance findings, capability conformance, and exact
  absence guards for owned V1 removals.

## Exact Exclusion And Ownership

| Deferred item | Owner | Phase 21 behavior |
|---|---|---|
| Attention and typed Stops | Phase 22 | No tables/commands/events; routine unsatisfied status stays queryable. |
| Budget consumption, exhaustion, warnings, retries, Stops, Attention | Phase 22 | `consumed:null`, `remaining:null`; no counters or automatic action. |
| Automatic respawn or automatic Run continuation | Phase 22 | Unsupported. fid never respawns an agent or process; recovery is explicit. |
| Run recovery, liveness resolution, linked successor Run | Phase 22 | No `run recover`, `run_recoveries`, or predecessor use. |
| Activation recovery and post-point-of-no-return repair | Phase 22 | Capability unavailable; only pre-commit abort is supported. |
| Export import and reconstruction parity | Phase 22 | No import command or hidden import path; capability unavailable. |
| Review, Findings, Human Outcomes, protected Decisions | Phase 23 | Their capabilities remain unavailable. |
| Multi-Issue selection, revision, cancellation, parallel scope analysis | Phase 24 | Limits reject more than one primary, any related Issue, revision/cancel, or parallel mutation. |
| Attested Receipts | Phase 25 | `attested` is not accepted; only `executed` can satisfy the supported Obligation. |
| Full twelve-capability qualification and migration baseline | Phase 25 | No live adoption claim in Phase 21. |

## Authority And Provenance Boundary

- `activation prepare` carries project-local bootstrap attribution but no
  protected capability claim.
- The finalized activation's immutable authority-profile component is the only
  source of kit-mode Actor identity, type, producer, provenance reference, and
  capability set.
- Subsequent callers submit only an Actor reference. fid resolves capabilities
  from the active activation, then intersects them with the exact Envelope
  Authority Grant. Caller-supplied capabilities, roles, or Human labels reject
  as unknown fields and confer no authority.
- A Gate Registration similarly resolves only an activated Gate Definition.
  Callers cannot replace argv, cwd, environment, limits, success policy,
  credential names, redaction, retention, network, or write policy at run time.
- Standalone Issue inventory keeps project-local attribution and cannot create
  trusted Receipts, protected Human authority, or kit-mode grants.

## Executor Boundary

Phase 21 supports only values it can enforce in portable Go:

- `shell:false` and a non-empty structured argv;
- exact repository/worktree-contained cwd;
- a closed environment entry list with literal, inherited, or credential-name
  sources;
- positive timeout, output, and excerpt limits bounded by fid config;
- credential names only, never stored credential values;
- redaction of resolved credential values and registered literals before live
  framing, excerpts, hashes, Artifacts, or Human rendering;
- content-addressed Artifact storage below the configured root; and
- `sandbox:"none"`, `network:"allowed"`, and `write:"allowed"` only.

Registration rejects any shell request, uncontained cwd/Artifact path,
unallowlisted environment input, unsupported sandbox value, or restrictive
network/write value. Phase 21 makes no sandbox, filesystem-isolation, network-
isolation, privilege-separation, or cross-platform containment claim. A
representative gate that requires one of those properties triggers a separate
dependency/platform decision rather than a silent weakening.

## Owned V1 Removal Boundary

Replacement-bound removals occur only after the generic replacement passes its
own task gate:

- every actor-attribution `--by` spelling on retained mutation paths becomes
  `--actor`; the unrelated `stats --by` grouping flag becomes `--group-by`
  so `--by` disappears from the live tree;
- the V1 role-policy implementation in `internal/policy`, the `roles` table,
  and `fid policy show` are removed after the generic `internal/policy`
  Activation/Envelope evaluator passes;
- `fid issue claim|unclaim`, legacy `fid claim list`, and the V1
  `internal/claims` role actions are replaced in place by generic Claim
  operations;
- legacy `fid run record` and `internal/runs` are replaced by generic Runs;
- `fid dev start|gate|finish|block`, `internal/dev`, `dev.*` Events,
  reported `gate_runs`, and gate snapshot reads are removed after trusted
  execution and Receipt evidence pass;
- `fid issue complete`, `internal/issues/complete*.go`, and all writers of
  `task_completions` are removed after Work Envelope completion passes; and
- no removed V1 state is reconstructed through export import, because Phase 21
  has no import operation.

Frozen legacy telemetry/completion projection tables needed by Phase 23/25
readers are listed explicitly in the slice contract, have no Phase 21 writer,
are not capability state, and carry a concrete deletion phase. They are not a
second authority for Work Envelope completion or generic Runs.

## Phase Shape

- **Operating mode:** Heavy.
- **Group 1:** Event/idempotency/schema/capability foundation; Kit Activation;
  activated Gate Registration; Work Envelope/Obligation/budget; generic Actor
  authority; retained actor convergence and policy deletion; Claims; Runs; and
  execution-bound transition/resume foundation.
- **Checkpoint:** fresh independent consistency review and executor
  feasibility. Any finding returns to its owning Group 1 task. Group 2 is not
  dispatched without PASS.
- **Group 2:** executor policy and framing; Attempt/Receipt; Artifact/Evidence/
  freshness; terminal transition/completion/resume; export; V1 removal and
  migration 0014; final-schema maintenance/rebuild; and the final
  real-process/conformance journey.
- **Review cadence:** every three completed Developer tasks, the Group 1
  checkpoint, then final QA and Human checklist.

## Approval State

Gate D3.2 is closed. Product implementation remains blocked until the Human
reviews and approves the exact slice contract, validation contract, Heavy-mode
task pack, and phase brief, and the QA Lead supplies a Ready Heavy-mode QA plan
with fresh feasibility evidence.
