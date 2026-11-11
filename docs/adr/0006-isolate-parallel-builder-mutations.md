---
status: accepted
date: 2026-08-08
---

# Isolate parallel Builder mutations by Issue and worktree

> **Amended on 2026-10-31 by ADR-0013.** Distinct Issues in separate
> worktrees may proceed without scope proofs; conflicts are resolved at merge.
> Declared, provably disjoint scopes wait for real need.

> **Amended on 2026-11-10 by ADR-0014.** At most one active hold per task
> stands. A hold binds the task, the actor and the worktree, with no Envelope
> revision. Reviewers record review rounds bound to the change's fingerprint
> instead of read-only Reviewer Runs.

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
ADR-0009 resolves multi-Issue Envelope revision/partial progress.
