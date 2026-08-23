# Phase 21 - TL Escalations

## Open Escalations

None. DEC-160 and DEC-161 resolve ESC-21-001..003; task 2 remains in the
ledger's stopped state until the Phase Runner appends the authorized resume.

## Resolved Escalations

### ESC-21-001 - Role-free bootstrap versus retained V1 policy

- **Raised by:** TRG-21-001 after task 2.
- **Resolution:** DEC-160 keeps the exact zero-Effect Event-v2 bootstrap and
  authorizes a closed, nonpersistent seven-name V1 policy vocabulary only for
  inventoried legacy callers. Explicit legacy role Effects remain
  authoritative; new generic paths cannot consult it; task 20 removes it.
- **Recovery:** same-task task-2 continuation with focused policy red/green
  tests and byte-identical reruns of both failed required commands.

### ESC-21-002 - Markdown-delimited external trigger identity

- **Raised by:** TRG-21-002 after rejected TRIAGE candidates.
- **Resolution:** DEC-161 makes schema-v2 identity values plain single-line
  scalars. TRG-21-001/TRG-21-002 remain immutable malformed evidence; neither
  governs recovery.
- **Recovery:** none against the verifier or historical records.

### ESC-21-003 - Protected recovery route

- **Raised by:** TRG-21-003 and the valid STOP ledger.
- **Resolution:** DEC-161 uses already-valid TRG-21-003 as the governing record
  and supplies exact RUN-RESUME attempt 3 authority. No protected verifier edit
  or control-conformance repair is required.
- **Recovery:** Phase Runner appends RUN-RESUME, then task-scope TASK-RECOVERY
  for task 2 under DEC-160/DEC-161. Task 3 stays dependency-gated.

## Archive Rule

Any Phase 21 `BLOCK` must allocate an `ESC-21-NNN` through `triage-log.md` and
record its protected resolution condition here before Human attention.
