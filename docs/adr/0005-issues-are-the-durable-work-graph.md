---
status: accepted
date: 2026-08-08
---

# Issues are the only durable work graph

fid Issues and their parent/dependency relationships are the only durable work
graph. A Work Envelope is a versioned execution contract over one primary
Issue and, when explicitly selected, a bounded set of related Issues; it does
not contain or create a second persistent Work Item graph. Durable
decomposition becomes a task or subtask only when it is independently
claimable, parallelizable, dependency-bearing, separately acceptable, or
expected to survive the current run. Smaller implementation steps remain an
agent's transient plan or a concise Run checkpoint.

The harness derives a noncanonical **Next action** from current Issue,
Envelope, Obligation, Claim, Run, Review, evidence, stop, and attention state.
Claims and Runs bind the selected Issue, Envelope revision, actor role, and
optional governing Obligation; independent review is represented by a Review
Run and Review Record satisfying a review Obligation, not by a Review Work
Item. No durable Work Item identity, table, lifecycle, or independently
editable source of truth is introduced.

## Consequences

- The V2 kit's existing in-memory `WorkItem`/`WorkGraph` implementation is
  superseded work in progress and must be realigned before production adapter
  implementation resumes. Its historical phase evidence remains unchanged.
- fid must expose deterministic Next-action and typed non-actionable-reason
  projections without persisting those projections as work entities.
- Parallel Builders require independently claimable Issues; ADR-0006 defines
  the Claim/worktree isolation rule. Procedural steps do not become Issues
  merely to make the harness advance.
- ADR-0006 through ADR-0009 resolve the required joint lifecycle scenarios:
  parallel Builder isolation, interrupted-Run recovery, review rework, and
  multi-Issue Envelope revision/partial progress.
- This decision resolves the Issue/Work Item identity question left open by
  ADR-0004. It does not decide the detailed schema or transport for those
  scenarios.
