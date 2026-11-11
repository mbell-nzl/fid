---
status: accepted
date: 2026-08-08
---

# Recover interruptions with a new linked Run

> **Amended on 2026-11-10 by ADR-0014.** Holds replace Claims and Runs. A
> lapsed hold on a started task means the task is interrupted. It stays kept
> for its worktree, which resumes it under a new hold linked to the old one,
> and other worktrees skip it. A lapse is still not proof that work stopped.
> Only a person takes over, and the takeover is recorded. fid never resets,
> cleans, stashes, checks out or commits. A check's start without a result
> counts for nothing. The rework count and token use stay with the task.

Continuation after agent-session or executor interruption always creates a new
Run linked to the prior Run; Run identity is never reused. The prior Run
becomes terminally `interrupted`, distinct from successful closure and
intentional abandonment. Issue and Work Envelope state remain active, while
obligations, concerns, evidence rules, and consumed budgets carry forward.

Recovery may take over mutation authority only after the prior executor is
confirmed terminated and its Claim is released, expired, or explicitly
superseded through an authorized force-recovery action. Lease expiry alone is
not proof that filesystem mutation has stopped. Unknown executor liveness
returns `executor-liveness-unknown` and blocks mutation.

Before the new Claim and Run may mutate, the kit/harness captures and
reconciles the current worktree against the Run start, last durable checkpoint,
approved scope, and repository history. In-scope recovered changes may be
explicitly adopted by the new Builder and are not trusted evidence. Protected,
unexpected, unattributable, or otherwise unsafe state returns
`recovery-state-conflict`. Recovery never automatically resets, cleans,
stashes, checks out, commits, or discards Git state.

The takeover transition records the interrupted prior Run, superseded Claim,
new Claim, linked Run, repository observation, and resumption/intervention
telemetry without resetting budgets. Receipts retain their normal exact-target
freshness rules; an execution start without a terminal receipt satisfies
nothing. Human authority is required only for force takeover or ambiguous/
protected repository state.

This resolves the interrupted-Run scenario required by ADR-0005. ADR-0008
resolves review rework and ADR-0009 resolves multi-Issue Envelope
revision/partial progress.
