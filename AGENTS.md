# fid - Project Context

## What Is This Project?

fid is a local-first, agent-first issue tracker for AI-assisted software
development. It provides a structured, queryable knowledge graph over the
artifacts an agentic development workflow produces. Rows hold metadata and
identity; markdown files hold prose. SQLite-backed CLI.

## Tech Stack

Concrete choices are recorded in `project-docs/specs/02_TRD.md` §1.

- Primary language/runtime: Go
- Application framework(s): `spf13/cobra` + `spf13/pflag` for the CLI command
  tree; no server framework in v1.
- Data/storage: SQLite via `modernc.org/sqlite` (pure Go), with
  journal-canonical persistence and a derived DB projection — see
  `docs/adr/0001` and `project-docs/specs/02_TRD.md`.
- Interface surfaces: CLI
- Test tools: stdlib `testing` + `stretchr/testify`; integration tests use
  temporary SQLite databases; benchmarks use `go test -bench`.
- Build/tool runner: `go build`, `go test`, and Makefile targets added during
  Phase 0 scaffolding.
- Deployment/runtime environment: local-first, single-user

If any item remains `<TBD>`, agents may still help refine specifications, but
they must not begin implementation that depends on that choice.

## Architecture

- User-facing adapters: `cmd/fid/`, `internal/cli/`, `internal/project/`,
  `internal/export/`.
- Core/domain/business logic: `internal/issues/`, `internal/claims/`,
  `internal/workflow/`, `internal/policy/`, `internal/qa/`,
  `internal/validation/`, `internal/templates/`, `internal/search/`,
  `internal/stats/`, `internal/activity/`, `internal/report/`, `internal/changes/`,
  `internal/handoff/`, `internal/evidence/`, `internal/maintenance/`.
- Persistence/integration adapters: `internal/journal/`, `internal/db/`,
  `internal/gitctx/`, `migrations/`.
- Shared utilities/design system: no shared design system in v1. Keep shared
  helpers inside the narrowest owning `internal/*` package until a broader
  boundary is recorded in `project-docs/specs/02_TRD.md`.

Rules:

- Business rules belong in the core/domain layer, not in interface adapters.
- Persistence and external-service details stay behind thin adapters.
- New boundaries must be recorded in `project-docs/specs/02_TRD.md`.

## Key Conventions

- Use red/green TDD for behavior changes: failing test first, then code.
- Tests must assert behavior, not implementation trivia.
- Do not delete or weaken tests to make a failing gate pass.
- Keep changes small and scoped to one task.
- A Developer session implements at most one task. A task may span a same-task
  continuation session when recovery or context renewal is required.
- Use file artifacts as shared memory. Do not rely on chat history.
- Update impacted documentation in the same task as the code change.
- Normal and heavy planned phases use `VAL-NNN` validation assertions in
  `project-docs/docs/phases/phase-N/qa/validation-contract.md`; task specs and
  completion reports must cite covered assertions or state why none apply.
- Use Conventional Commits for commit subjects.
- Include `Refs: phase-N/task-M` in commit body or footer for phase work.
- Prefer existing local patterns over new abstractions.
- Add abstractions only when they remove real complexity or match an existing
  project pattern.
- Use `project-docs/specs/GLOSSARY.md` as the canonical project language
  source. Update it or flag ambiguity when terms are unclear.
- Record spec deviations in the completion report and phase decisions.
- Treat specifications as the primary engineering artifact. If a requirement
  cannot be satisfied, stop, report the gap, and do not stub, mock, or silently
  adapt the spec.
- If a required validation assertion cannot be satisfied, stop and escalate
  rather than weakening the assertion or passing it without evidence.
- Deterministic claims require evidence: command output, test logs, file state,
  screenshots, API traces, or equivalent reproducible proof.
- Keep held-out scenario material outside normal implementation context. Do not
  read hidden scenario repositories, excluded scenario folders, or evaluator
  prompts unless the task is explicitly scenario validation.

Project-specific conventions:

- Naming: follow Go package/function naming conventions; package names are
  short, lowercase, and domain-specific. fid entity IDs use
  `<PREFIX>-<hash>`; phase validation assertions use `VAL-NNN`; QA plan steps
  use `QA-NN`.
- Error handling: return errors rather than panicking for expected failures;
  wrap lower-level errors with context using `fmt.Errorf("operation: %w", err)`;
  present structured stderr output and JSON errors when `--json` is active in
  the target contract. Existing `--format json` fixtures are displaced
  pre-convergence implementation, not a compatibility promise.
- Formatting style: Go code uses `gofmt` and `goimports`; Markdown stays
  concise and keeps command examples current with the CLI specs.
- File/function size budgets: non-test functions target 60 lines or less;
  non-test files target 400 lines or less, per
  `project-docs/specs/08_QUALITY_PLAN.md`.
- Test style: use red/green TDD for behavior changes; prefer table-driven
  tests; test observable behavior and contracts, not implementation trivia;
  use temporary SQLite databases for integration tests and golden fixtures for
  journal/event compatibility.
- Public interface compatibility policy: before the first supported baseline,
  ADR-0011 permits replacement of current CLI flags, schemas, events,
  migrations, and generic formatters without a legacy mode. Accepted target
  JSON contracts and slice event schemas require a spec update and matching
  tests before implementation; after first adoption they become normal
  versioned compatibility surfaces.

## Agent Prompts And Skills

Role prompts are single-sourced: edit `project-docs/agents/<role>.md` and run
`project-docs/tools/sync-agents.sh`; `.claude/agents/` and `.codex/agents/`
are generated provider files.

Canonical repo skills live under `project-docs/skills/`. Provider-discoverable
copies are generated from canonical skills into `.agents/skills/` for Codex and
`.claude/skills/` for Claude. Codex discovers `.agents/skills/` natively and
invokes a skill implicitly when the task matches its description or explicitly
via `$skill-name`; Claude discovers `.claude/skills/` through its skill
surface.

Rules:

- Edit canonical skills only; do not edit generated provider copies directly.
- Every skill must include a `Provenance:` line in the Markdown body.
- Skills are procedural helpers, not the source of truth for role authority,
  hard stops, required artifacts, validation gates, dependency policy, public
  interface policy, or architecture decisions.
- First versions should be instruction-only. Scripts require documented
  prerequisites, safe non-interactive behavior, and dependency-budget rationale
  when they add dependencies.
- Run `project-docs/tools/sync-agents.sh` after editing canonical role prompts.
- Run `project-docs/tools/sync-skills.sh` after editing canonical skills.
- Run `project-docs/tools/workflow-lint.sh` after prompt or skill generation to
  check skill, prompt, template, and config drift.

## Dependency Budget

Do not add a new runtime dependency without a spec update, task spec, decision
entry, or completion-report note recording:

- Need: what problem it solves and why the standard library or an existing
  dependency is insufficient.
- Alternatives considered.
- License check.

Default license policy: prefer permissive licenses such as MIT, BSD,
Apache-2.0, or MPL-2.0. Avoid GPL, AGPL, SSPL, unknown, or custom licenses
unless the human explicitly approves and records the decision.

Test-only and dev-only dependencies follow the same rule with a lower bar:
record the rationale in the task completion report.

## Quality And Release Readiness

Use the maturity path in `project-docs/specs/08_QUALITY_PLAN.md`:

1. Baseline: specs, TDD, task reports, QA plan, deterministic QA evidence, and
   human checklist.
2. Automated gates: stack-appropriate format, lint, test, coverage or equivalent
   confidence gate, type checks where available, SAST, dependency audit, license
   check, and secret scan.
3. Held-out scenarios: user-journey scenarios stored outside normal agent
   context and run before release readiness; Mode B human runs are required for
   product judgment.
4. Release readiness: threat model, pre-release checklist, changelog, migration
   notes where applicable, backup/rollback path, and periodic independent audits.

Record any deliberate deferral in a task spec, phase decision, or roadmap item.

Each phase declares an operating mode (light, normal, or heavy) in the phase
brief. Modes vary required artifacts, never evidence rules or hard stops. See
`project-docs/specs/09_AI_WORKFLOW.md` section 4.

## Commands

Replace every placeholder before implementation starts. Agents must stop and
report a configuration defect if a required command is still a placeholder.
Use the task-spec command list when it exists; otherwise use these defaults.
Every required command listed by a task spec or QA plan must be recorded in the
completion report as `PASS`, or as `BLOCKED` with reason in a blocker report.
A missing, silently skipped, or narrower substituted required command is a
failed gate, not an acceptable omission. Record each result with one short
verbatim token from the command's own output (its last printed line, or
`exit 0 (no output)`). A command that did not terminate is never `PASS` — a
hung, interrupted, or killed command is `BLOCKED`. **A silent command is not a
finished command:** `go test ./...` compiles every package before printing
anything, so the repository-wide race and dedicated integration gates can run
for minutes with no output and then emit package result lines. Wait for them.
`exit 0 (no output)` is never a valid token for a `go test` command. A search or inventory
command (`rg`, `grep`) that finds nothing terminates **non-zero**: record its
real status, `exit 1 (no output)`, and say that no-match is the pass condition.
`exit 0 (no output)` means silent success and is fabricated evidence there.

Record each result's **observed termination** next to its token — the nested
command's real `exit_code`, or the explicit signal/termination result the
execution surface reported:
``- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)``.
A printed last line proves the command produced output; it does not by itself
prove the process exited, and for a silent success the exit status is the only
thing that distinguishes one from a discarded session. Copy the termination
result from the command's own result object. Wrapper completion, an empty
output field, and absence of new output are not termination, and an `exit_code`
asserted rather than observed is fabricated evidence exactly like a fabricated
token.

A command whose pass condition is a **non-zero** exit — an `rg`/`grep` absence
proof, an `ls` on a path that must not exist — records that real status in the
declared form, and its task spec pre-declares the expected status under
`## Expected Non-Zero Commands`:
``- `rg -n "internal/format" cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)``.
The status field stays exactly `PASS`; the expectation is what the spec
declared, and the termination is what the command actually did. An undeclared
non-zero `PASS`, a declared expectation of zero, and an observed status that
differs from the declared one are all failed gates. Never convert a non-zero
absence result into `exit 0 (no output)` to satisfy a control.

A gate result **not captured in the session that records it** is `UNVERIFIED`,
never `PASS` — attestation from memory, transcription from another agent's
report, and a promised future control are all `UNVERIFIED`, which is not a
passing state. Declining an out-of-scope repair as "pre-existing" is a
**provenance claim**: cite the introducing commit (`git log -S`,
`git log --diff-filter=A`, `git blame`) and show it predates phase start.
`git show HEAD` proves a line is committed, not inherited. **Cleanup is never
evidence**: commands that only remove or prune scratch state are never recorded
`PASS`/`BLOCKED` and never gate a sequence. Allocate scratch state fresh per
attempt (`mktemp -d`), and declare the working directory for any command whose
result depends on it.

A required-command failure is an immutable failed attempt, not automatically a
Human-visible blocker. Two cases are distinct:

- A **protected gate/control definition change** changes the command, the
  assertion's accepted pass/fail semantics, accepted status, scope,
  threshold/budget/tolerance, or governing policy. It is not repair and needs
  the authority that owns the definition.
- A **control-conformance repair** fixes an implementation bug in a canonical
  verification/workflow control when an already accepted contract uniquely
  determines its behavior. It changes none of the protected items above. The
  fact that the control is canonical or outside the originating task's ordinary
  file list does not make this case `BLOCK`.

Ordinary causal repair stays inside accepted task scope and leaves the
gate/control definition intact. A cross-scope control-conformance repair
requires an authoritative TL-Lite triage-log directive with these exact typed
fields: `Originating task`, `Control trigger signature`,
`Control attempt count`, `Control contract`, `Control mismatch`,
`Allowed files`, `Forbidden changes`, `Exact verification`, `Originating
report`, and `Control repair base`.
Each value is one physical line. The Phase Runner records and verifies
the full pre-dispatch phase HEAD as `Control repair base`; the Developer cannot
select, replace, or rebaseline it. The initial TASK-RECOVERY maps its ledger
trigger/attempt to the control fields and copies every other directive value
byte-for-byte. Product, spec, prompt, template, dependency, task, QA-contract,
report, and generated files are never eligible for the isolated repair.

The Phase Runner creates a dedicated clean temporary git worktree at the
recorded base before dispatching a same-task continuation into that path. The
dirty runner-owned ledger stays in the main phase worktree as its only dirt.
The Developer edits only the named tool paths in the repair worktree, uses
red/green, preserves all failed evidence, runs `Exact verification`, and leaves
the ending HEAD clean. Necessary provider sync is a separate authorized
workflow-maintenance step, never part of the isolated control repair.

Before recording the `REPAIR` ledger event, the Phase Runner runs this exact
scope check over the dedicated commit:

```text
project-docs/tools/workflow-control-scope-lint.sh <repo> <base> <head> <allowed...>
```

The `<base>` is always the runner-recorded directive value and `<head>` is the
Developer-produced repair worktree HEAD. The claimed `<head>` must be that
worktree's checked-out `HEAD`, its
worktree/index must be clean with no `assume-unchanged` or `skip-worktree`
tracked-file index flags, and every
allowlisted artifact must remain a regular file rather than a symlink or
gitlink. After any scope PASS, fast-forward the unchanged phase branch to
verified HEAD while its only dirt is the ledger, confirm phase HEAD, append
REPAIR, and clean up the temporary worktree. For a nonresolved control attempt,
the next recovery increments attempt count and uses that prior integrated HEAD
as the new runner-recorded base; every semantic directive value remains
byte-identical. A control `RESOLVED` permits only a linked
`report-repair` limited to the exact `Originating report` and `Exact
verification`, retaining the resolved attempt's `Control repair base`; direct
task completion is invalid. Only that report repair's
later `RESOLVED` permits completion. Any ambiguity or proposed
definition/policy change is `BLOCK`.

Every new v2 RUN-START/RUN-RESUME event carries
`Recovery directive schema: v2` in that same event; historical v1 entries
remain valid, but no new event may emit or downgrade to v1. A RUN-RESUME after
STOP retains the governing TRG/ESC and binds its exact next recovery. When an
existing DEC resolves a BLOCK, amend the existing triage record and use
RUN-RESUME; do not mint a new TRG. Every boundary-valid DEC citation must exist
in the phase decisions file.

No failed test, gate, or control grants generic permission to edit its test,
command, or implementation. Authority comes from the task scope or the written
control-conformance directive. After a causal repair, rerun the
**byte-identical** command and record both attempts plus a mandatory deviation
note containing the first failure's verbatim token and termination, root cause,
repair, and later result. The first attempt remains failed evidence; the later
passing attempt is the evidence used to close the task.

Never delete, weaken, narrow, re-point, or mine a flaky green. If the needed
change is protected or ambiguous, record a typed `BLOCKED` outcome. If work is
decision-free but outside ordinary task files, return `RECOVERY_REQUIRED` so
the Phase Runner can obtain the bounded directive above; do not self-authorize
the edit. If in-scope work cannot finish in the current session, also record
`RECOVERY_REQUIRED`; incomplete work is not a blocker by itself.

Sandboxed Codex agents use `.codex/config.toml` to redirect Go and tool caches
into gitignored repo-local directories (`.gocache/`, `.gomodcache/`,
`.gopath/`, `.golangci-cache/`, `.cache/`). Do not request unsandboxed
execution just to write Go build, module, or lint caches. A fresh agent session
is required after changing that config. Network downloads and dependency
installs still require the normal dependency-budget and approval path.

- Format check: `gofmt -l cmd internal`. Apply `gofmt -w` and `goimports` to
  edited Go files before running the check. Do not use `gofmt -l .`; it walks
  ignored repo-local tool caches.
- Lint: `golangci-lint run ./...`.
- Default unit test suite: `go test -race ./...`.
- Coverage gate:

  ```bash
  go test -race -coverprofile=coverage.out ./...
  go tool cover -func=coverage.out
  ```

- Real-process integration journeys:
  `go test -tags integration -race ./cmd/fid/integration/...`. In-process
  SQLite and command contract tests belong in the ordinary suite; do not widen
  this gate to `./...`, which reruns that suite under the tag.
- End-to-end tests: `N/A` until an end-to-end or held-out scenario harness is
  accepted in the specs.
- Build verification: `go build -o /tmp/fid ./cmd/fid`.
- Full phase-level command gate:

  ```bash
  gofmt -l cmd internal
  golangci-lint run ./...
  go test -race -coverprofile=coverage.out ./...
  go tool cover -func=coverage.out
  go test -tags integration -race ./cmd/fid/integration/...
  go test -bench . -benchmem ./...
  go build -o /tmp/fid ./cmd/fid
  ```

  Phase-level QA also performs the PRD compliance, security review,
  performance spot-check, documentation review, deterministic QA checklist, and
  human checklist steps in `project-docs/specs/08_QUALITY_PLAN.md` §5.2.

Release-readiness commands are stack-specific and may be `N/A` early in a
project, but must be concrete before a public or user-impacting release:

- Security scan: `N/A` before external release; a SAST or equivalent command
  must be selected and recorded before release readiness.
- Dependency vulnerability audit: `govulncheck ./...`.
- License compatibility check: `go-licenses check ./...`.
- Committed-secret scan: `gitleaks detect`.
- End-to-end or scenario validation: `N/A` until the held-out scenario harness
  location and command are recorded.

Recommended task-level gate:

```bash
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Recommended phase-level gate:

```bash
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

## Playbook

See `project-docs/docs/playbook.md` for a concise operator runbook.

## Specification Documents

- `project-docs/specs/00_PROJECT_CHARTER.md` - vision, scope, constraints.
- `project-docs/specs/GLOSSARY.md` - canonical project language.
- `project-docs/specs/01_PRD.md` - product requirements.
- `project-docs/specs/02_TRD.md` - technical requirements.
- `project-docs/specs/03_TRD_SCHEMA.md` - target logical schema and invariants;
  physical DDL/events are accepted slice by slice before implementation.
- `project-docs/specs/04_TRD_CONFIG.md` - `.fid/` directory structure and config.
- `project-docs/specs/05_DESIGN_PLAN.md` - user experience and surfaces.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - CLI flag specs and output examples.
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` - end-to-end flows, whet mapping, export.
- `project-docs/specs/08_QUALITY_PLAN.md` - testing and verification.
- `project-docs/specs/09_AI_WORKFLOW.md` - agent workflow.
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` - phase plan.

## Agent Working Documents

- `project-docs/docs/current/` - default-loaded current context.
- `project-docs/docs/phases/phase-N/` - archived per-phase task, QA, progress,
  decisions, triage, and escalation artifacts.
- `project-docs/docs/phases/phase-N/qa/validation-contract.md` - phase
  `VAL-NNN` assertions that connect requirements to task evidence and final QA
  outcomes.
- `project-docs/docs/reviews/` - retrospectives, audits, and external reviews.
- `project-docs/templates/` - reusable artifact templates.
- `project-docs/skills/` - canonical skill authoring source.

## Hard Stops

Agents stop for Human or protected-authority input only when:

- A deterministic quality gate cannot be fixed under the assigned task scope
  or a bounded control-conformance directive without changing a protected
  command, assertion semantics, accepted status/scope,
  threshold/budget/tolerance, or policy.
- A required project command is still a placeholder.
- A task requires an architectural, product, data-model, or public-interface
  decision not already recorded in specs.
- A task depends on ambiguous terminology that is not resolved or flagged in
  `project-docs/specs/GLOSSARY.md`.
- A dependency addition lacks documented budget rationale.
- A required artifact or external capability remains unavailable after its
  bounded recovery budget is exhausted.
- A required phase validation contract is missing, has uncovered `VAL-NNN`
  assertions, or needs a product/architecture decision before evidence can be
  produced.
- A normal implementation task tries to read held-out scenarios or evaluator
  prompts that are meant to stay outside agent context.
- A release-readiness phase lacks the required threat model, pre-release
  checklist, or scenario/audit decision for its maturity level.

A red test, failed command, report-lint defect, bounded QA finding, interrupted
session, or incomplete in-scope implementation is `REPAIR`, not a Human-visible
stop. Preserve the failed evidence, checkpoint the work, and continue through
the recovery state machine in `project-docs/specs/09_AI_WORKFLOW.md`.
