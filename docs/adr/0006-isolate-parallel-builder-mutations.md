---
status: accepted
date: 2026-08-08
---

# Isolate parallel Builder mutations by Issue and worktree

One fid Issue may have at most one active mutating Builder Claim. A Work
Envelope is not itself claimed: each Claim binds one selected Issue, the
Envelope revision, the Builder, and the repository worktree identity. Parallel
Builder mutation is permitted only for distinct selected Issues in separate
worktrees with declared, provably disjoint mutation scopes. Otherwise the
second action is non-actionable with typed reason `mutation-conflict`.

Work inside one Issue must be split into independently meaningful fid
subtasks before it can be mutated in parallel; procedural steps do not justify
new Issues. Read-only Reviewer Runs may coexist, but each is bound to the exact
repository identity reviewed and its result becomes stale when that target
changes. The kit/harness establishes the approved worktree and scope isolation;
fid atomically enforces the resulting Claim-conflict preconditions.

This resolves the parallel-Builder scenario required by ADR-0005. ADR-0007
resolves interrupted-Run recovery and ADR-0008 resolves review rework.
Multi-Issue Envelope revision/partial completion remains open.
