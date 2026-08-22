# Phase 21 Group 1 Independent Checkpoint Contract

> Owner: Tech Lead
> Required after: task 12 and its QA Review window
> Required before: task 13 dispatch
> Governing assertions: VAL-042 and VAL-043

This is a dispatch gate inside one Heavy Phase 21, not a phase boundary. The
QA Lead assigns a reviewer independent of tasks 1-12 and records the outcome
in `group-1-checkpoint-review.md`. The later QA plan supplies the exact bounded
commands; the contract and accepted pass/fail semantics below are protected.

## Preconditions

- Tasks 1-12 have completion reports with task-local report-lint evidence.
- QA Review artifacts exist for tasks 1-3, 4-6, 7-9, and 10-12.
- Every finding is linked to its owning task and resolved by a later committed
  repair or remains an explicit checkpoint failure.
- Group 1 HEAD and the exact contract/spec commit are recorded.

## Independent Contract Review

The reviewer must establish:

- Event v2, idempotency, migration, and rebuild facts have one authority;
- capability availability remains exactly the section 2 matrix;
- finalized Activation is the only kit Actor/Gate trust source;
- Gate Registration precedes Envelope creation and accepts no policy override;
- Envelope, Obligation, budget, grant, Claim, Run, and transition bindings
  match the slice contract and new generic paths do not consult V1 authority;
- the temporary replacement inventory is exactly the one in
  `../v1-convergence-inventory.md` section 5: `roles`, V1 `claims`,
  `gate_runs`, `task_completions`, `agent_runs`, `ac_verifications`,
  `concerns`, role/judgment transition-requirement rows, `fid dev`, Issue
  complete, and their replacement-bound writers/packages may remain; the
  internal V1 policy implementation may remain only for those exact callers,
  while the public `fid policy` command is already absent;
- AC verification and concern mutation are already fail-closed, and none of
  that temporary state can authorize a generic operation or satisfy Evidence,
  a guarded terminal transition, or completion;
- final command/table/writer/Event absence is explicitly reserved for task 20
  and final-schema init/rebuild/export parity for task 21; and
- no Group 2 implementation or capability advertisement was pulled forward.

## Fresh Group 2 Feasibility

Session-captured probes must newly establish that the repository/target test
environment can deterministically verify:

- structured-argv child processes with separate stdout/stderr, invalid UTF-8,
  control bytes, signal/timeout/spawn failure, and high-output draining;
- closed environment injection and canary credential/literal redaction across
  live, persisted, Artifact, error, Human, and export surfaces;
- symlink-resolved cwd/Artifact containment, file-change-during-copy, and
  content-addressed atomic publication;
- exact observed termination rather than wrapper completion;
- repository-tree freshness and Artifact loss/restoration;
- byte-deterministic ZIP creation and non-mutating corruption validation;
- every prescribed expected-nonzero absence command with its real status; and
- the dedicated tagged real-process journey within the QA plan's declared
  bounded time budget.

Transcribed Phase 20 evidence, a promised later fixture, an unbounded/hung
command, or a green result without observed termination is not feasibility.

## Verdict

`PASS` requires every precondition, review item, temporary-inventory
classification, generic-path isolation proof, and fresh feasibility probe.
`REPAIR` returns each finding to its owning Group 1 task and repeats this
checkpoint after normal recovery. `BLOCK` is reserved for a protected contract
change, an unavailable external capability after bounded recovery, or a
dependency/platform isolation need. Group 2 is never dispatched on `REPAIR`,
`BLOCK`, partial, conditional, or undocumented status.
