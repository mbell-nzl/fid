# Held-Out Scenario Recommendation

Date: 2026-06-14
Task: phase-6/task-5

## Judgment

Recommend using the v0.1b smoke from `cmd/fid/e2e_v01b_test.go` as the **seed**
for the first held-out scenario harness input once the harness location is
formally accepted.

## Why Recommend

- It exercises the whole v0.1b command surface in one user journey:
  `fid init` -> `fid issue create` -> `fid issue list` -> `fid issue get` ->
  `fid issue update` -> `fid issue transition` -> `fid ready`.
- It already proves the cross-surface invariant that matters for v0.1b exit:
  deleting `.fid/fid.sqlite` and running `fid rebuild` reproduces the same
  projected state and scheduling answers.
- It is deterministic and self-contained: temp workspace only, no network, no
  machine-local dependencies, and pinned expectations for sequential and
  prioritised epics.
- It is small enough to remain auditable, but broad enough to catch command
  interaction regressions that narrower unit or CLI tests can miss.

## Recommended Isolation Mechanism

Keep the held-out scenario **outside the implementation repo and outside the
agent workspace**. The recommended shape is:

- Store the scenario script/fixture set in a separate QA-only repository or
  sibling checkout that is not mounted into normal implementation-agent
  sessions.
- Expose only a runner interface to implementation agents and CI, for example a
  phase-end command or wrapper script that receives the built `fid` binary path
  and returns pass/fail plus a transcript.
- Preserve the existing agent rule that hidden scenario material is out of
  normal implementation context; do not copy scenario prose, expected outputs,
  or evaluator prompts back into this repo.

This keeps the deterministic smoke reusable while still preventing
implementation agents from training directly against the held-out assertions.

## Current Recommendation Boundary

- Do not change `AGENTS.md` end-to-end command guidance yet. The current
  `End-to-end tests: N/A` note should remain until the harness location and
  command are accepted as a project decision.
- Do not move or create the held-out repository in this task. This task only
  recommends the candidate scenario and the isolation mechanism.
