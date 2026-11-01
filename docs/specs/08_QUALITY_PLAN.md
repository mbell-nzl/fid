# QUALITY PLAN - fid

**Status:** Binding repository quality gates and target product quality contract
**Last updated:** 2026-08-09
**Authority:** [ADR-0011](../../docs/adr/0011-realign-fid-around-the-durable-operational-kernel.md)

## 1. Purpose

This document defines both the quality gates used while building fid and how
fid earns trust as an Issue tracker and durable operational backend. It tests
product behavior, storage integrity, recovery, machine contracts, and the
boundary between project policy and fid facts.

`09_AI_WORKFLOW.md` prescribes the current V1 development process. The
repository gates and validation rules below are binding inputs to that process
until the Human approves a replacement. They are delivery scaffolding, not fid
product features. The remaining sections define target product assurance and
first-adoption evidence without changing the active agent workflow.

## 2. Quality properties

The first supported baseline must demonstrate:

1. **Durability:** accepted commands survive restart and rebuild without losing
   identity, ordering, authority, or evidence lineage.
2. **Determinism:** the same valid command at the same revision produces the
   same state change or the same idempotent result.
3. **Fail-closed enforcement:** missing authority, stale revisions, invalid
   Claims, unsatisfied Obligations, and untrusted evidence block protected
   transitions.
4. **Traceability:** completion and Human decisions can be explained from
   immutable events, Receipts, reviews, and Outcomes.
5. **Recovery:** interrupted Runs and projection loss yield an actionable Next
   Action after rebuild.
6. **Contract clarity:** JSON schemas, exit status, capability negotiation, and
   error categories are stable and testable.
7. **Boundary discipline:** fid stores operational facts and references; it
   does not become an authoring store for project requirements, QA plans,
   validation prose, role prompts, or templates.
8. **Usability:** common Issue and Work Envelope operations remain concise for
   Humans while agents receive bounded JSON.

## 3. Test layers

| Layer | Primary proof |
|---|---|
| Domain unit tests | State transitions, invariants, authority, Obligation evaluation, freshness, and completion rules. |
| Command contract tests | JSON request/response schemas, exit status, error taxonomy, idempotency keys, revision preconditions, and Human text behavior. |
| Journal/projection tests | Event append, ordering, hash/integrity checks, replay, rebuild, and projection equivalence. |
| Persistence integration tests | SQLite transactions, locking, constraints, concurrent claims, and crash-safe recovery in temporary repositories. |
| Adapter conformance tests | The fid adapter and the in-memory fake pass the same backend-port capability suite. |
| Workflow scenario tests | End-to-end issue, activation, envelope, run, gate, review, Human outcome, stop, resume, and completion journeys. |
| Fault-injection tests | Interrupted execution, duplicate delivery, stale revision, corrupt projection, missing artifact, changed gate registration, and partial write attempts. |
| Security tests | Path containment, argument-vector execution, digest verification, redaction, malicious project input, and least authority. |
| Performance tests | Representative latency, JSON size, replay cost, and query behavior using realistic repositories. |
| Disposable migration rehearsal | A controlled V1-style project is converted to the clean fid/V2 baseline, exercised, rebuilt, and discarded before any live migration. |

Tests assert observable contracts rather than package layout or SQL trivia.
Table-driven unit tests and temporary SQLite databases are the default Go
patterns. Golden fixtures are appropriate for canonical JSON and event
compatibility only after that contract is accepted for a slice.

A test encodes a contract; it is not itself an authority. When an accepted
change supersedes the contract a test encodes, the change owner adapts that test
in the same unit of work and under the same gate, preserving the test's subject
and assertion strength. Dropping a removed argument, converting an assertion
about removed behaviour into the accepted behaviour, or repointing incidental
reads is convergence. Deleting an assertion requires another test that already
owns the accepted equivalent. Skips, build tags, narrowed test selection, and
relaxed thresholds are never convergence.

Convergence ownership follows the invalidating change, not the file. A test
file's primary subject may belong to a different unit of work; that does not
transfer or defer the obligation, and it is not a scope boundary. The change
owner converges only the assertions its own change invalidated and leaves every
other assertion in that file byte-identical, so the other subject keeps full
strength until its own owner replaces it.

Task numbering in a retained filename, test function, fixture label, or prose
name is historical metadata, not current-phase ownership evidence. Before
classifying retained-test convergence as cross-task or protected, establish the
introducing phase/task from commit provenance and compare the assertion with
the current task pack, convergence inventory, and accepted contract. A
historical label transfers ownership only when the current task pack explicitly
adopts it. Repetition of an unproven filename-based classification does not make
it a specification contradiction or exhaust a causal recovery budget.

A unit of work's required commands must exercise every file scope it owns.
Convergence ownership is only assignable where it is executable: a build-tagged,
environment-gated, or otherwise excluded suite covering owned files must have
the gate that runs it named in that same unit of work, or the scope does not
belong to it. Adding files to a scope without adding the gate that exercises
them is an incomplete assignment. Where a gated suite first becomes required
several units of work after the changes that invalidated its assertions, the
accumulated convergence belongs to the first unit whose required gate runs it,
and that unit records per-row provenance - file, line, invalidating change, and
originating unit - so the audit trail survives the reassignment.

The converse also holds. A unit of work that owns a required absence or
inventory control owns every path that control walks, and its file scope must
cover every search root and every package its required gates compile. A
control that reaches beyond the listed scope is an incomplete planning
assignment, not a boundary the implementer must honour by leaving the control
red; it is corrected by amending the scope before dispatch, never by narrowing
the control. Planning checks each task's search roots and test package
arguments against its file list.

The same reasoning assigns a latent production defect. Where an accepted
contract was implemented wrongly and no earlier required gate could observe it,
the defect belongs to the first unit of work whose required gate does observe
it, not to the unit that wrote the code. That unit receives the minimum source
scope needed for the causal repair, plus a regression assertion executed by a
command it already runs, so the amended scope stays executable; it does not
receive generic access to the owning package, and the originating unit is not
reopened. The repair corrects the implementation to the accepted contract and
changes nothing else.

A sequential end-to-end fixture reveals one blocking step at a time. Reaching a
later step is material progress. A no-progress recovery budget counts only
attempts whose directive actually permitted causal repair of the failing rows;
an attempt forbidden from touching the only cause of the failure is a
mis-scoped directive, not evidence of exhausted repair capability.

The same budget counts only attempts that reported an outcome. An attempt whose
session is interrupted, cancelled, timed out, or shut down from outside, and
which therefore never returns a typed result, does not consume it: the budget
measures a repair that cannot converge, and an interrupted attempt was never
allowed to finish converging. No budget bounds a child session's lifetime; only
declared per-command budgets are enforced, and the orchestrator waits for the
session to end on its own. Repeated external interruption is its own condition -
orchestration session-liveness - recorded with the interruption cause, elapsed
duration, and command in flight, because its remedy is the session shape rather
than the repair. Where verification cost approaches the session window, commit
durable work before the long commands and append each command result as it
terminates, so an interruption leaves inspectable evidence instead of an
unchanged tree that reads as no progress.

A parent polling wait timeout is not child termination and never authorizes an
interrupt: retain the same child and poll again. The first actual external
interruption routes to same-task recovery; the second consecutive one routes to
orchestration session-liveness triage. The ledger records a typed
`Outcome: INTERRUPTED` and preserves the preceding no-progress count rather
than incrementing or resetting it. When that triage authorizes repair, the next
recovery increments total attempt count, preserves the no-progress count, and
starts a new consecutive-interruption streak.

## 4. Critical scenario inventory

The conformance and workflow suites must cover at least these scenarios:

| Scenario | Required result |
|---|---|
| Standalone Issue lifecycle | Issue creation, relationship validation, readiness, transition, and rebuild agree. |
| Kit activation | Project identity, kit version, capability profile, and configuration digests are recorded atomically and can be read back. |
| Envelope activation | Immutable revision, selected Issues, actors, authority, gate registrations, and Obligations are complete before activation. |
| Parallel Builders | Disjoint Claims may proceed concurrently; overlapping Claims or stale writes fail deterministically. |
| Registered gate execution | fid resolves the immutable registration, runs the exact argument vector, captures trusted process facts, records artifacts, and evaluates applicable Obligations. |
| Receipt resistance | Caller-reported command text, exit status, or output cannot masquerade as a fid-observed Receipt. |
| Freshness | A source or artifact change invalidates only the evidence whose declared inputs changed. |
| Review coverage | Required Reviewer coverage and finding disposition are enforced without embedding a global role taxonomy. |
| Protected decision | A pending Human decision blocks only its declared protected transition; the recorded Outcome unblocks or redirects it. |
| Human takes control (Phase 24) | Retain `TestP24HumanTakeControlJourney` as a real-process scenario: a protected Decision refuses terminal success, Human help/templates and labelled output supply every command/value, the Decision Reference and approving Outcome resolve the Stops, linked Evidence permits the transition, and the Envelope completes. The transition target and current Issue revision come from the last Human `resume` read. `TestP24TakeControlRebuildExportImportParity` proves that completed workspace survives rebuild and schema-19 export/validate/import, including field-digest maps, observed repositories and judgment detail reads. |
| Stop and resume | A typed Stop releases or retains Claims according to policy and produces a bounded recovery projection. |
| Revision after progress | New work uses the new Envelope Revision while prior Receipts and Outcomes retain their original revision lineage. |
| Completion | Completion is derived only when Issue state, Claims/Runs, Obligations, reviews, Outcomes, and Stops permit it. |
| Rebuild/export | A blank projection rebuilt from the journal is semantically equivalent; portable export validates without SQLite access. |
| Retired surfaces | Removed V1 commands and generic TOON/YAML/Markdown format selections fail clearly and are absent from advertised capabilities. |

## 5. Evidence and trust

Deterministic claims require reproducible evidence. For product validation this
means test output, command traces, event fixtures, hashes, rebuilt state, or
captured process results. Human judgment is recorded as a Human Outcome and is
never relabelled as deterministic proof.

An Execution Receipt is trusted only when fid or an explicitly registered
trusted executor observed the attempt. Attestation records must identify the
attestor, authority, subject, revision, and trust class. They cannot silently
satisfy an Obligation that requires direct execution.

Test-only shortcuts may stand in for clocks, process runners, repositories, or
artifact stores, but the conformance suite must exercise the production
adapter. The in-memory fake is a test oracle, not a production fallback.

## 6. Machine-contract quality

JSON is the sole canonical machine representation. Each accepted command
contract specifies:

- schema version and capability identifier;
- success and error envelopes;
- stable IDs and revision fields;
- idempotency and compare-and-set behavior;
- bounded list pagination or explicit size limits;
- whether absent, empty, and unknown values differ;
- redaction rules; and
- forward-compatibility behavior for unknown fields and capabilities.

Generic TOON, YAML, and Markdown formatters are not tested or supported.
Purpose-built Human Markdown or HTML projections require their own named
acceptance tests but are not machine compatibility surfaces.

## 7. Repository engineering gate

These commands remain the default implementation gate while fid is developed:

```bash
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 ./cmd/fid/integration/...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

Under Phase22 DEC-221, PR CI executes the complete ordinary race-and-coverage
suite once across four exhaustive test shards, rather than running both make
test and make coverage. Each worker discovers runnable Tests, Examples and Fuzz
seed targets through Go, and selection is disjoint and complete. The controller
reconciles actual execution and package completion, then combines atomic coverage
only when all four source/inventory-bound manifests and all statement blocks
agree. A missing, failed, cancelled or inconsistent shard fails the retained
required `quality` job. Format, lint, build and controller regressions remain
required. Standalone make test/make coverage targets remain available unchanged.
Package timeouts, product thresholds and assertions are unchanged. Setup reuse
requires distinct writable copies; migration, rebuild and import tests still
execute their actual operations. See `python3 .github/scripts/ci_tests.py
--help` for local usage.

Under Phase22 DEC-217, the full-size archive byte/count matrices are an explicit
phase-end and archive-limiter-change gate, run locally on the qualification
machine rather than on every routine task/PR:

```bash
go test -tags archive_scale -race ./internal/export -run '^TestArchive(ExpandedByteBoundaries|CumulativeBoundary|EntryCountBoundaries)$' -count=1
```

The ordinary suite retains small tests through the same limiter, production
ceiling wiring and large-number arithmetic checks. Actual production-scale
boundary evidence remains required for phase qualification; all product
ceilings are unchanged. Run heavy gates serially.

Task scope may justify a narrower gate only when the repository instructions or
accepted task specification says so. Failing tests are fixed by correcting the
product or accepted contract, never by weakening evidence. Runtime and test
dependencies follow the recorded dependency budget and license policy.

Every required command must pass or be recorded as `BLOCKED` with reason in a
blocker report. Missing, silently skipped, or narrower substituted commands are
failed gates. A no-match search records its real non-zero status rather than a
fabricated successful exit.

Coverage is a diagnostic, not a substitute for scenario coverage. Changes to
journal, authority, evidence, completion, recovery, or public JSON contracts
require broader tests even when line coverage is unchanged.

## 8. Slice acceptance

Every convergence slice must include, in order:

1. accepted terminology and command/JSON contract;
2. a failing behavioral or conformance test;
3. the smallest coherent implementation;
4. migration or deliberate removal of displaced code and fixtures;
5. adapter conformance and recovery evidence;
6. documentation and capability advertisement updates; and
7. an independent consistency review when the slice changes a cross-product
   contract.

Each `VAL-NNN` assertion states one observable behavior, invariant, or release-
readiness claim. Conditional assertions define evidence for qualifying and
non-qualifying states; positive-only evidence is insufficient. Task specs,
completion reports, reviews, and final QA outcomes cite the same assertion IDs
under the active V1 workflow.

A slice is not complete if it leaves two authoritative models for the same
fact. Temporary compilation scaffolding is allowed only inside the slice and
must not be advertised as a supported capability.

For the completed Heavy Phase 21, three-task QA Review windows applied. Its
separate independent Group 1 checkpoint must PASS after the authority/Activation/Registration/
Envelope/Claim/Run foundation and before executor implementation begins. That
checkpoint re-runs fresh feasibility for every non-default verifier, fixture,
environment assumption, expected non-zero proof, and time budget prescribed
for Group 2. Findings return to their owning task; implementation momentum is
not evidence. Report lint is not a Developer task-gate command: after the
Developer's report commit, the Phase Runner runs the exact two-argument
report/task-spec control before the next task is dispatched.

The checkpoint proves new generic Group 1 paths are isolated from the exact
temporary V1 inventory; it does not require final absence before the removal
task. Phase 21 task 20 performs migration 0014 and V1 removal. Task 21 then
performs final-schema fresh init, role-free rebuild, maintenance, live/rebuild/
export parity, and the 100k Event-v2 check. Task 22 performs final wire and
real-process conformance. The journey mutates once; post-rebuild evidence is
read/evaluation/export/negative parity, not a second Envelope mutation.

For approved Heavy Phase 22 (DEC-180/186/188), reviews follow tasks 3/6/9/12/15/18;
the independent task-12 checkpoint precedes archive/import work. Its QA plan
and TL refinement addendum enumerate the deterministic/Human checks. Tasks 1
and 16 add one fresh reference rebuild sample each, without replacing task 17
or final QA's unchanged three-sample rule. Task 8 owns the known terminal-outcome
fixture repair; earlier failures require the bounded scope amendment in DEC-186.
The explicit VAL-056 boundary matrix uses real measured decompression at the
unchanged ceilings. Phase 21's historical checkpoints are not extra Phase 22
task dependencies; the standing evidence and report-lint rules remain binding.

For approved Heavy Phase 23, VAL-015 proves reachable Claim-history rejection
and same-independent-Reviewer rework through the CLI. Its isolated Run-holder
query uses a diagnostic durable-history fixture because the accepted lifecycle
binds each Run to its Claim's Actor and preserves Claim history (DEC-239).
Record the fixture's evidence class explicitly; it is not a CLI-reachability,
journal-acceptance, or import-eligibility claim. Both history directions and the
exact conflicting subject and Actor remain required.

## 9. First supported baseline

Before fid is used by a live project:

- all twelve required backend-port capabilities pass the shared conformance
  suite against the fake and fid adapters;
- every critical scenario in section 4 passes against the production CLI and
  storage adapter;
- the generic authority replacement exists before V1 role policy is removed;
- displaced V1 commands, aggregates, output formatters, migrations, and event
  fixtures are gone;
- journal and schema history are consolidated into a clean initial baseline;
- rebuild, portable export, backup, restore, and integrity checks pass;
- threat modelling covers command execution, project input, artifact paths,
  secrets, and local privilege boundaries;
- dependency, vulnerability, license, and committed-secret checks are selected
  and pass;
- the disposable migration rehearsal succeeds without reading fid internals;
  and
- a Human validates the routine Issue, Work Envelope, review, recovery, and
  decision interfaces.

No live project migration is used to discover whether the backend contract is
complete.

## 10. Post-adoption signals

Dogfood records enough telemetry to answer whether fid reduces or adds work:

- commands and tokens per completed Issue;
- retries, stale-write conflicts, and recovery time;
- evidence collection cost and invalidation rate;
- Human interruptions by reason and lane;
- review findings caught before completion; and
- time spent reconstructing context or querying state.

These measures are projections over trustworthy operational events. They do
not justify speculative analytics or a new durable reporting aggregate.

Candidate features are reconsidered only when observed use reaches the trigger
recorded in ADR-0011. In particular, alternate token encodings must beat compact
JSON end to end, including instructions, parse failures, and retries, without a
correctness regression.
