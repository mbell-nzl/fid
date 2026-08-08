# IMPLEMENTATION ROADMAP - fid

**Status:** Product surface accepted; replacement slice order pending Gate D3
**Last updated:** 2026-08-09
**Authority:** [ADR-0011](../../docs/adr/0011-realign-fid-around-the-durable-operational-kernel.md)

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

### Open decision

**Gate D3** must ratify the exact vertical-slice order, dependency edges, and
the slice in which each displaced implementation surface leaves the build.
Section 5 is a candidate sequence for that discussion, not an approved task
pack.

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

## 5. Candidate vertical-slice sequence for Gate D3

The proposed order keeps fid buildable and places a replacement beside each
high-risk removal before the old model disappears.

| Candidate slice | Coherent outcome | Principal displacement |
|---|---|---|
| A. Contract and clean machine surface | Versioned JSON envelope, capability discovery, error taxonomy, idempotency/revision conventions, project identity, and shared fake/fid conformance harness. | Remove generic TOON, YAML, and Markdown formatters and `--format`; retain only Human text plus `--json`. Remove promised-command stubs with no target contract. |
| B. Envelope execution kernel | Envelope/Revisions, selected Issues, Actors, capabilities, Authority Grants, scoped Claims, Runs, Stops, Attention, recovery, and Next Action work end to end. | Remove V1 role matrix, `fid policy show`, role-bound Claims, escalation/triage aggregates, and `fid dev start/block` choreography after generic authority is proven. |
| C. Trusted deterministic evidence | Project Gate Definition activation, immutable registrations, Attempts, Receipts, artifacts, Evidence Records, trust classes, and freshness work end to end. | Remove reported `gate_runs`, gate snapshots, universal `fid dev gate`, and any completion path that trusts caller-reported command outcomes. |
| D. Review, Human judgment, and completion | Review Records, coverage, Findings, Concerns, Deviations, Decision references, Human Outcomes, acceptance evidence, and completion evaluation work end to end. | Remove role-specific verification columns, QA/TL composites, completion reports as state, and V1 `fid dev finish`. |
| E. Activation, portability, and trajectory | Atomic Kit Activation, portable export validation, rebuild equivalence, changed-since/context/Human projections, and trajectory measures complete the port. | Remove whet-style Markdown export, generic report promises, and any duplicate activation or recovery truth. |
| F. Baseline convergence | Consolidated initial schema/events, clean repository initialization, security/recovery hardening, full conformance, disposable migration rehearsal, and Human validation. | Remove superseded migrations, event fixtures, compatibility tests, dead packages, and temporary slice scaffolding. |

Questions Gate D3 must settle:

1. Whether Kit Activation belongs in slice A rather than E because all later
   operations need an activated contract.
2. Whether Attention should arrive with Stops in slice B or with Human Outcomes
   in slice D.
3. Whether portable export must appear early as a test oracle for every slice.
4. Which old commands remain temporarily buildable inside a slice without being
   advertised as supported.
5. The smallest end-to-end acceptance scenario and exact deletion list for each
   slice.

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

1. Complete the coordinated spec and kit-contract rewrite for Gate D2.
2. Run an independent cross-document consistency audit and fix material
   contradictions.
3. Ratify Gate D3 using the candidate sequence and questions in section 5.
4. Create the first bounded implementation task with exact removals, contracts,
   failing tests, and quality commands.

No implementation task should cite the superseded historical version roadmap
as forward authority.
