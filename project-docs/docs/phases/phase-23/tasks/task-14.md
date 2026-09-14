# Phase 23 - Task 14: Human review and judgment projections

## Objective

Give every Phase 23 leaf purpose-labeled Human text derived from the same facts
as its JSON, replace the resume Findings placeholder with the real projection,
and add the bounded judgment sections and per-Obligation evaluations.

## Validation Contract Coverage

- VAL-058: Every new leaf has Human text matching the accepted obligations and
  adds no JSON member.
- VAL-059: `resume` reports real Findings plus bounded judgment sections and
  per-Obligation evaluations, and refuses to truncate silently.
- VAL-060: Human text never presents a judgment fact as deterministic proof and
  never reproduces project prose.
- VAL-073: advanced here; the new Human views follow the accepted CF-040
  guidance. The assertion itself is Human-validated at phase end.

## Binding Constraints

- `contracts.md` section 14 is the exact authority for what each Human view must
  say. `public-interface.md` sections 2 and 7 fix the projection shapes.
- Human text adds no JSON member, changes no mutation rule, reproduces no
  project prose beyond a reference label, path, and digest, and never presents
  a Human Outcome, Concern, or Deviation as deterministic proof.
- Finding text must say explicitly that `addressed` is a Builder assertion and
  not a resolution. Review text must never say a `pass` verdict completed the
  Issue.
- `resume.findings` replaces the `{available:false,items:[]}` placeholder with
  the real open-Finding projection. Add `reviews`, `concerns`, `deviations`,
  `decisions`, `human_outcomes`, and `protected_decisions`, each bounded by the
  list limit, and make `obligations` the full per-Obligation evaluation array.
- Exceeding a bound returns `result-limit-exceeded`. Never truncate judgment
  state silently.
- `next_actions` gain the newly executable review and judgment operations with
  their reasons; reasons are empty only when the operation is executable.
- CF-040 guidance applies to the views introduced here: labeled sections, an
  operation purpose line, no routine narration. Polish of pre-existing views
  stays with Phase 25 and must not be started here.
- Record the verbatim Human fixtures in
  `project-docs/docs/phases/phase-23/qa/human-fixtures.md` so the window 13-15
  review can perform its deterministic Human-text conformance check.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/workflow/resume.go`, `resume_actions.go`.
- `cmd/fid/resume.go`, `resume_human.go`, `resume_human_guidance.go`.
- Human renderers for each Phase 23 leaf under `cmd/fid/`.
- `project-docs/docs/phases/phase-23/qa/human-fixtures.md` - new.
- Matching in-package tests and `cmd/fid/integration/` fixtures.
- `project-docs/docs/phases/phase-23/tasks/task-14-done.md`.

### Files To Read For Context

- `contracts.md` section 14; `public-interface.md` sections 2 and 7.
- `project-docs/docs/phases/phase-22/public-interface.md` section 6 and
  `project-docs/docs/phases/phase-22/qa/human-fixtures.md` for the accepted
  DEC-219 Human-guidance pattern.
- `task-13-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Trajectory economics and the full Human Validation Brief, which stay Phase 25.
- Redesign or verbosity rework of pre-existing Phase 21 and Phase 22 views.

## Acceptance Criteria

- [ ] Every new leaf has a Human fixture satisfying each stated obligation and
  introducing no member absent from JSON. (VAL-058)
- [ ] Resume reports real Findings and every declared section, reproducible from
  durable state, with `result-limit-exceeded` past a bound. (VAL-059)
- [ ] No fixture claims proof for a judgment fact or embeds decision or
  requirement prose. (VAL-060)
- [ ] Every new view has labeled sections, an operation purpose line, and no
  routine narration. (VAL-058)
- [ ] `qa/human-fixtures.md` records the verbatim output used by the window
  review.
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestResumeReportsRealFindingsAndJudgmentSections` fails against the retained
Findings placeholder.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-14-done.md` in the full format.
- `project-docs/docs/phases/phase-23/qa/human-fixtures.md` with verbatim output.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. No exception to the full task
gate. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination. An
interrupted or hung command is never `PASS`; preserve all failed attempts and
rerun the byte-identical command after an authorized causal repair. Report lint
is a Runner control after the report commit, not a Developer gate.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- This task belongs to QA Review window 5.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-14`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
