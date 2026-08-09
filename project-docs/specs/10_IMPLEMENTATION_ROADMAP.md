# IMPLEMENTATION ROADMAP - fid

**Status:** Gate D3.1 and revised Phase 20 implementation pack accepted
**Last updated:** 2026-08-09
**Authority:** [ADR-0011](../../docs/adr/0011-realign-fid-around-the-durable-operational-kernel.md),
[ADR-0012](../../docs/adr/0012-deliver-convergence-as-vertical-scenario-phases.md)

## 1. Objective

Converge the useful implemented fid foundation and the V2 kit backend contract
into one clean first supported baseline. Build all required capabilities in
bounded vertical slices, qualify each slice, then rehearse migration on a
disposable project before any live adoption.

This roadmap does not preserve the historical V1 product surface. fid has no
production users or journal history, so superseded commands, schemas, events,
migrations, and formatters are removed rather than deprecated.

## 2. Current position

### Accepted decisions

- fid remains a useful standalone Issue tracker.
- the kit owns engineering policy and orchestration; fid owns durable
  operational facts and enforcement.
- an Issue is the only durable work unit; a Work Envelope authorizes execution
  over a selected Issue set.
- generic Actors, capabilities, and Envelope Authority Grants replace global
  kit-role policy.
- deterministic gate evidence comes from immutable registered execution, not
  caller-reported results.
- the complete twelve-capability backend port is implemented before live
  project migration.
- the historical V1 compatibility premise is removed.
- QA Plans, Validation Contracts, templates, skills, handoffs, and role
  choreography are not fid product aggregates.
- JSON is the sole canonical machine representation; generic TOON, YAML, and
  Markdown output is removed.
- speculative features live only in ADR-0011's evidence-gated candidate
  ledger.
- convergence is delivered as phase-sized vertical scenario slices; horizontal
  work may support a phase but is not its acceptance boundary.

### Current gate position

**Gate D3.1** has ratified the vertical-slice order and initial removal owners
in section 5. The Human accepted the exact Phase 20 contract, validation
contract, and eighteen-task pack on 2026-08-09 under the binding V1 delivery
process. QA-plan acceptance is the remaining pre-implementation gate.
**Gate D3.2** will settle the exact thin
Routine walking-skeleton boundary before Phase 21 is planned.

## 3. Foundation to retain

The existing codebase supplies useful implementation assets, subject to
conformance with the rewritten contracts:

- Go CLI and package structure;
- local repository discovery and `.fid/` storage;
- journal append, SQLite projection, locking, transactions, and rebuild;
- stable entity IDs and ordered events;
- Issue CRUD, hierarchy, dependency, label, acceptance, and workflow behavior;
- Claims, activity, repository context, and maintenance primitives that can be
  reshaped around Envelope Revisions; and
- test fixtures and utilities that prove retained behavior rather than V1
  compatibility.

Calling an asset "retain" does not freeze its public command, event, table, or
package shape. Only the target contracts in the current specs are forward
requirements.

## 4. Required capability set

The first supported baseline advertises and passes conformance for all of these
capabilities:

| Capability | Outcome |
|---|---|
| `workflow.state.v1` | Read and mutate accepted operational state with revision checks. |
| `workflow.transitions.v1` | Explain and enforce permitted next transitions. |
| `work.claims.v1` | Coordinate lease-based, scope-aware execution. |
| `kit.activation.v1` | Atomically align kit, overlay, provider, project, and backend state. |
| `evidence.gate-definitions.v1` | Register immutable approved executable definitions. |
| `evidence.receipts.v1` | Observe trusted attempts and immutable terminal Receipts. |
| `evidence.freshness.v1` | Evaluate source, artifact, toolchain, and contract applicability. |
| `reviews.coverage.v1` | Record independent review coverage and Finding lifecycle. |
| `human.outcomes.v1` | Record authorized Human judgment without relabelling it as deterministic proof. |
| `attention.v1` | Expose deduplicated intervention needs and their disposition. |
| `telemetry.trajectory.v1` | Measure workflow trajectory from trustworthy operational events. |
| `export.portable.v1` | Validate backend-neutral recovery and migration data. |

Capability identifiers describe behavior, not packages or one-command-per-
capability CLI design. Exact JSON schemas and event types are accepted with the
owning slice before its code changes.

## 5. Approved vertical-scenario phase sequence

Each slice is one V1-style implementation phase containing multiple bounded
Developer tasks. Schema, event, adapter, migration, and cleanup tasks may be
horizontal inside the phase; the acceptance boundary is the coherent scenario
in the second column. Exact JSON and event contracts are accepted by the
owning phase before its product code changes.

| Phase / slice | End-to-end acceptance journey | Contract and capability advance | Owned legacy removal |
|---|---|---|---|
| **20. Clean standalone Issue journey** | Initialize a repository; create, relate, inspect, update, and transition Issues; explain failure; rebuild; and repeat the journey through canonical JSON. | Establish common JSON success/error envelopes, project identity, revision/idempotency conventions where used, and the retained standalone Issue contract. | Remove universal TOON/YAML/Markdown formatters, `--format`, format-matrix tests, obsolete promised-command stubs, and Issue-path fields or commands with no retained meaning. Do not yet remove fail-closed V1 authority. |
| **21. One-Issue Routine walking skeleton** | Activate one kit/project/backend contract; select one Issue; authorize and claim one Builder Run; execute one immutable registered gate through fid; create trusted evidence; evaluate completion; then resume and export the same state. | Introduce the thinnest coherent Envelope/Revision, Actor authority, Claim/Run, Gate Registration/Attempt/Receipt, Evidence/Freshness, Attention, completion, activation, and portable-export path. Gate D3.2 fixes the exact inclusion boundary before planning. | After replacements pass, remove the V1 role matrix and `fid policy show`, role-bound Claim/Run fields, `fid dev start/gate/finish/block`, reported `gate_runs`, gate snapshots, and completion-report-as-state paths. |
| **22. Failure and recovery** | Interrupt or abandon execution, recover liveness, reconcile repository state, start a linked Run, retry or stop safely, stale prior proof, rebuild, and resume without contradictory truth. | Complete typed Stops, Attention disposition, lease/liveness rules, activation recovery, freshness failure, rebuild equivalence, and recovery-grade export behavior. | Remove escalation entities and commands, duplicate recovery/activity truth, caller-reported terminal-result paths, and temporary walking-skeleton recovery shortcuts. |
| **23. Review and Human judgment** | An independent Reviewer records coverage and zero-or-more Findings; Builder rework preserves review independence; protected or subjective Obligations receive authorized Human Outcomes; completion is reevaluated. | Complete Reviews, Findings, Concerns, Deviations, project Decision/requirement references, Human Outcomes, acceptance evidence, and Human-facing decision projections. | Remove role-specific verification columns, QA/TL composites, legacy requirement authoring that duplicates project prose, and remaining role-specific completion behavior. |
| **24. Parallel and evolving work** | Execute distinct selected Issues in disjoint worktrees; reject conflicting Claims; revise a multi-Issue Envelope; carry valid completion forward; handle removed unfinished work explicitly; and preserve review rework rules. | Complete multi-Issue revision, scope conflict, dependency, partial-progress, parallel-claim, linked-Run, and evolving-obligation semantics. | Remove remaining V1 Claim/Run shapes, triage and carry-forward composites, and any Envelope-owned duplicate work or partial-progress state. |
| **25. Qualification and baseline convergence** | Run the full shared fake/fid conformance suite, all recovery and security scenarios, portable validation, a disposable V1-project input migration rehearsal, and Human validation of the operating experience. | Qualify all twelve capabilities together; finish trajectory telemetry and named operational/Human projections; consolidate the first supported schema and event baseline. | Remove superseded migrations and event fixtures, V1 parity/compatibility tests, dead packages, temporary slice scaffolding, generic report/whet-export residue, and any retired surface found by the whole-product audit. |

The removal column is an ownership map, not permission to keep a known retired
surface until Phase 25. Phase planning must enumerate exact commands, packages,
tables, events, migrations, fixtures, and tests. If implementation discovery
moves an item, the owning phase records the reason and a concrete new expiry
point.

### 5.1 Current implementation ownership check

This is the pre-planning map for implementation that exists after Phase 19.
The Tech Lead must replace it with an exact file/schema/event inventory in each
phase plan. Phase 20's exact inventory is
`../docs/phases/phase-20/contracts.md` section 5.

| Current surface | Planned treatment |
|---|---|
| `internal/format`, root format selection, and V1 format-matrix tests | Delete in Phase 20; retain direct Human text and canonical `--json` behavior, not a universal formatter framework. |
| `governance_stubs.go` and promised-command stub tests | Delete in Phase 20. Later commands appear only with an owning accepted contract. |
| `internal/issues`, `labels`, `links`, Issue CLI, and Issue workflow authoring | Exercise and sharpen in Phase 20; retain the domain behavior that passes the new standalone contract. |
| `internal/policy`, `internal/dev`, `fid policy`, `fid dev`, and role-bound portions of `claims`, `runs`, and `workflow` | Replace and delete in Phase 21 after generic Actor/Envelope authority, Claims/Runs, trusted execution, and completion pass. |
| Legacy recovery/activity facts that duplicate Run, Stop, Attention, or activation recovery | Replace/delete in Phase 22; retain activity only as an event-derived projection. |
| `internal/requirements` authoring and role-specific acceptance/concern verification | Replace in Phase 23 with project-owned prose references plus generic acceptance, Review, Finding, Concern, Decision, and Human Outcome facts. |
| Remaining V1 Claim/Run shapes and any triage/carry-forward composites exposed during multi-Issue work | Delete in Phase 24 after parallel and Envelope-revision rules pass. |
| `internal/stats` behavior not justified by trajectory telemetry, V1 parity/e2e fixtures, superseded migrations/events, and dead packages | Replace or delete in Phase 25 during whole-product convergence. |
| `internal/journal`, `internal/db`, locking, IDs, rebuild, maintenance, and `fid sql` | Retain as foundation subject to the new contracts; `fid sql` remains an unstable development diagnostic. |

## 6. Removal rules

Every displaced surface receives one of two treatments:

- **Immediate removal:** the behavior has no retained purpose and no replacement
  dependency, such as generic TOON/YAML/Markdown formatting or unimplemented
  command stubs.
- **Replacement-bound removal:** the old behavior carries a safety property
  that must remain fail-closed until its generic replacement passes, such as
  role permissions, Claims, reported gates, or completion enforcement.

There is no compatibility adapter, hidden legacy mode, dual-write period, or
production fallback to the in-memory fake. Old journal data is not migrated.
Git history is the archive.

The following historical roadmap families are cancelled as product promises:

- global role-policy authoring and kit-role command trees;
- escalation and triage workflows;
- QA Plan and Validation Contract aggregates;
- templates, agent skills, and handoff generation;
- retrospective and carry-forward lifecycles;
- whet mapping and bidirectional Markdown export;
- generic multi-format output; and
- speculative TUI, search, saved views, MCP, federation, cycles, and advanced
  analytics until their ADR-0011 triggers are met.

TUI/dashboard and lightweight search remain the leading optional candidates;
their desirability does not put them on the committed path.

## 7. Slice acceptance rule

A convergence slice is complete only when:

- terminology, invariants, JSON schemas, errors, and event contracts are
  accepted before implementation;
- red/green tests cover the public behavior and relevant failure paths;
- fake and fid adapters pass the same capability conformance cases;
- replay/rebuild and interruption behavior are proven where applicable;
- displaced code, tests, docs, migrations, and fixtures named by the slice are
  removed or deliberately carried with an expiry point;
- no two authoritative models remain for one fact;
- repository quality gates pass; and
- cross-repository contract changes receive an independent consistency review.

The phase/task machinery used to implement a slice follows the current
repository instructions until the V2 kit can replace them. It is delivery
scaffolding, not a reason to restore removed product scope.

## 8. First-adoption gate

No live project adopts fid until:

1. all twelve capabilities pass the shared conformance suite;
2. standalone Issue use and full kit-driven execution both pass end-to-end;
3. trusted gates, review, Human judgment, interruption, revision, completion,
   rebuild, and export scenarios pass;
4. all retired V1 public surfaces and generic formatters are absent;
5. schema and journal events form a clean initial supported baseline;
6. security, dependency, license, secret, and recovery gates pass;
7. a disposable V1-style project migration rehearsal passes; and
8. the Human approves the routine operating experience and recovery path.

The rehearsal may transform project-owned V1 workflow files into V2 inputs. It
does not migrate old fid data because no production fid data exists.

## 9. Dogfood and iteration

After the first-adoption gate, migrate one controlled active project at a
deliberate boundary. Continue meaningful product work through the new path
while recording retries, intervention causes, recovery cost, evidence burden,
review yield, and context reconstruction cost.

Observed friction may change policy, projections, or candidate priority. It
does not bypass the product boundary by introducing a second durable store or
reading `.fid/` internals.

## 10. Immediate next work

1. **Accepted (2026-08-09).** Phase 20 planned, V1-workflow-self-checked, independently audited, and Human-approved from section 5. Slice contract:
   `../docs/phases/phase-20/contracts.md`; DDL amendment: `03_TRD_SCHEMA.md`
   section 8; config amendment: `04_TRD_CONFIG.md` sections 2/5;
   cross-cutting CLI rules: `06_DESIGN_PLAN_CLI.md` sections 2, 5, 8.
2. QA Lead writes and validates the Phase 20 QA plan; Phase Runner then starts
   Task 1.
3. Implement and independently review Phase 20 through the existing V1 process.
4. Resolve Gate D3.2's exact one-Issue Routine walking-skeleton boundary before
   planning Phase 21.

No implementation task should cite the superseded historical version roadmap
as forward authority.
