# Historical escalation as raised

Captured from the Runner/TL-Lite working artifacts before resolution.

## ESC-23-002 - Task 1 gate closure needs protected ownership and budget authority

- **Phase:** 23
- **Status:** OPEN
- **Raised by:** TL-Lite via TRG-23-001
- **Category:** protected gate/control authority and cross-task retained-assertion ownership
- **Decision needed from:** full Tech Lead, with Human authority if the accepted gate definition must change

### Summary

Task 1 cannot close. Its required 100k rebuild benchmark measures
32.500159404 s/rebuild against the strict under-30-second self-assertion, and
the exact repository-wide race gate fails with `cmd/fid` at 619.168 s and
`internal/workflow` at 614.196 s, each over DEC-227's accepted 600-second
ordinary-package budget. The same race run retains
`cmd/fid/task12_execution_transition_test.go:217`, whose assertion accepts only
`workflow-transition-invalid` but observed `claim-expired`, `run-inactive`, and
`workflow-transition-invalid`.

Task 1 may repair speed only where its own new journeys caused the breach. No
evidence establishes that causal ownership for the retained cmd/fid assertion
or both budget breaches. Changing the benchmark threshold, package budget,
command, or retained assertion semantics is protected and is not authorized by
the task. The full Tech Lead must establish bounded causal ownership and repair
authority, or obtain the authority required for any accepted-definition change.

### Evidence

- Benchmark command `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` failed: `FAIL github.com/mbell-nzl/fid/internal/journal 35.245s` (terminated `exit_code: 1`); measured rebuild 32.500159404 s.
- Race command `go test -race -timeout 30m -coverprofile=coverage.out ./...` failed (terminated `exit_code: 1`); `cmd/fid` took 619.168 s and `internal/workflow` 614.196 s.
- The task's temporary performance experiments were removed because they either failed to reach the budget or broke accepted rebuild parity; report lint returned `OUTCOME RECOVERY_REQUIRED` (terminated `exit_code: 0`).

### Required resolution

Full Tech Lead: identify, from accepted ownership and causation evidence, whether
one or more bounded repairs exist and assign the correct task authority,
including a fresh independent review where required. If the only remedy changes
the accepted command, timeout, threshold, budget, assertion semantics, accepted
status/scope, or policy, obtain the owning protected authority before any
change. Do not assign the retained `cmd/fid` assertion to task 1 without that
authority.

### References

- TRG-23-001 in `project-docs/docs/phases/phase-23/triage-log.md`.
- `project-docs/docs/phases/phase-23/tasks/task-1.md` (Exact Commands To Run).
- `project-docs/docs/phases/phase-23/tasks/task-1-done.md` (Workflow outcome and gate evidence).
- `project-docs/docs/phases/phase-23/execution-scope.md:44-45,501-508`.
- `project-docs/docs/phases/phase-23/decisions.md` DEC-227 and `project-docs/docs/phases/phase-23/qa/validation-contract.md` VAL-079.
