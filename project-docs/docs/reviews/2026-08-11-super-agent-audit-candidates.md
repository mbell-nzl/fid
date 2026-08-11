# Super-Agent Audit Candidates

Date: 2026-08-11
Status: Candidate portfolio for later Human selection; no audit or remediation
is authorized by this note.

## Purpose

Record cross-cutting audits that benefit from a long-running recovery-style
agent outside the normal one-task Developer boundary. These audits should
complement phase work rather than run concurrently against the same files.

## Candidate Audits

### 1. Test Reachability And Gate Truth

Inventory every Go test, build constraint, platform suffix, `TestMain`,
`t.Skip`, fuzz test, benchmark, package, and accepted gate. Prove which command
executes each test and detect `[no tests to run]`, excluded packages,
cache-masked evidence, and tests reachable only through retired commands.

Recommendation: highest immediate infrastructure value. Run after Phase 20 so
the inventory targets a stable slice.

### 2. Journal And Crash Consistency

Audit journal append ordering, fsync boundaries, partial records, corrupted
tails, SQLite transaction failure, lock interruption, rebuild recovery, atomic
initialization, and projection replacement. Test the binding invariant that
the journal remains canonical and recoverable at every interruption point.

Recommendation: highest product-foundation value. Run after Phase 20 and
before later operational capabilities depend on the retained journal layer.

### 3. State-Machine And Replay Fuzzing

Generate valid and invalid sequences across Issues, aggregate revisions,
transitions, concerns, acceptance criteria, labels, links, and workflow
definitions. Compare incremental projection state with a blank rebuild and
verify deterministic errors, ordering, and revision behavior.

Recommendation: use Go's native fuzzing and a small reference model without a
runtime dependency. Run after Phase 20 stabilizes the standalone Issue slice.

### 4. Phase 21 Boundary Architecture

Apply the architecture-review procedure to module depth, dependency direction,
transaction ownership, error translation, and leakage between retained Issue
behavior and the V1 `policy`, `dev`, Claims, Runs, and workflow surfaces that
Phase 21 will replace. Stress-test Gate D3.2's exact one-Issue Routine
walking-skeleton boundary without proposing interfaces prematurely.

Recommendation: run immediately before Phase 21 planning.

### 5. Filesystem And Hostile-Input Security

Stage A covers project discovery, symlinks, path containment, permissions,
hostile config and journal values, SQLite diagnostic access, temporary files,
artifact paths, and repository escape attempts. Stage B, after registered
execution exists, covers argument vectors, environment handling, secrets,
output bounds, cancellation, artifact containment, and child-process cleanup.

Recommendation: run Stage A after Phase 20 and Stage B after the owning
execution/evidence phase.

### 6. Performance And Scaling

Measure cold CLI latency, initialization, incremental append, rebuild at
representative event counts, large-Issue queries, SQLite query plans,
allocations, file descriptors, and lock contention. Establish explicit budgets
and profiles before choosing optimizations; include CF-024 disposition.

Recommendation: run before evidence volume and parallel execution expand.

### 7. Spec-To-Code And Retired-Surface Traceability

Trace accepted commands, JSON fields, reason codes, events, tables, migrations,
and advertised capabilities to implementation and tests. Detect duplicate
authoritative models, stale documentation, and legacy surfaces whose ADR-0011
expiry owner has passed.

Recommendation: perform a bounded pass at each phase boundary and a complete
whole-product pass before Phase 25 closes the supported baseline.

### 8. Agent-Workflow Adversarial Drills

Exercise controlled cases for silent long commands, yielded nested sessions,
truncated output, cached tests, no-match searches, dirty worktrees, uncommitted
reports, compaction recovery, stale agent state, failed report lint, and
reviewer independence. Use disposable drill fixtures rather than product tasks.

Recommendation: run after Phase 20, before another long heavy phase changes
the delivery workflow.

### 9. Release And Supply Chain

Audit the threat model, dependency provenance, licenses, vulnerabilities,
secret scanning, reproducible builds, backup and restore, portable export,
migration rehearsal, and release-checklist completeness.

Recommendation: run during Phase 25 and first-adoption preparation, after the
relevant commands and data contracts exist.

## Recommended Priority

The three strongest post-Phase-20 audits are:

1. Test reachability and gate truth.
2. Journal and crash consistency.
3. State-machine and replay fuzzing.

Their findings should inform the Phase 21 boundary architecture audit and Gate
D3.2 decision rather than silently becoming implementation scope.

## Execution Protocol

1. Create a durable audit anchor with a frozen commit range and explicit
   exclusions.
2. Audit read-only first; do not mix finding discovery with opportunistic
   fixes.
3. Classify findings as Blocking, Strong improvement, or Speculative.
4. Obtain Human approval for a bounded remediation pack.
5. Implement separate corrections with focused red/green evidence, exact
   gates, and discrete commits.
6. Give a fresh independent reviewer the frozen implementation range.
7. Close only after a NON-BLOCK verdict and synchronized current-context
   handoff.

Normal phase work and a super-agent campaign must not modify overlapping files
concurrently. The audit artifact is shared memory; chat history is not.
