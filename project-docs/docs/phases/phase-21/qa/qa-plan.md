# Phase 21 QA Plan - One-Issue Operational Walking Skeleton

> Phase: 21 - Heavy - One-Issue Operational Walking Skeleton
> Refreshed: 2026-08-23 (re-review of DEC-158 refinements)
> QA Lead: Claude Opus 5 (QA Lead role)
> Planning base validated: `463183a16ac99879eeeadc35ff5f1064980892e5`, re-validated at
> `b367a3b` (branch `phase-21`, clean worktree)
> Readiness verdict: **READY** - see section 2.3 and section 13
> Product implementation: **NOT AUTHORIZED**. Task 1 stays blocked pending Human
> execution approval.

This plan is the test plan for `qa/validation-contract.md` (VAL-001..087). It
does not assign QA verdicts, does not edit the validation contract, the
checkpoint contract, or any task spec, and does not authorize task 1.

Revision history: the plan was first committed at `f5ccbb6` with the verdict
Needs Refinement and refinement proposals R-01..R-19 plus blocker candidates
BC-01 and BC-02. The Tech Lead applied DEC-158 at `b367a3b`. This revision is
the independent QA Lead re-review of those amendments. Every command added or
changed by DEC-158 was re-executed fresh in this session; section 6 records the
results and section 13 records the disposition of every proposal. No Must-fix
finding remains, so the verdict is Ready.

## 1. Scope And Source Documents

Phase 21 delivers one Heavy one-Issue operational walking skeleton in 22
bounded Developer tasks and two internal groups: activated trust, immutable
Gate Registration, one Work Envelope revision, one gate-evidence Obligation,
generic Actor/Claim/Run authority, fid-executed Receipt Evidence, a guarded
terminal-success Issue transition, later Envelope completion, bounded resume,
export create/validate, migration 0014 V1 removal, final-schema rebuild
parity, and final wire/journey conformance.

Sources read to build this plan:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/decisions_active.md` (DEC-158)
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/docs/phases/phase-21/decisions.md` (DEC-158 record)
- `project-docs/specs/02_TRD.md` section 9 (Phase 21 reference instance)
- `project-docs/specs/08_QUALITY_PLAN.md` (sections 5, 6, 7, 8)
- `project-docs/specs/09_AI_WORKFLOW.md` (sections 2.6, 2.7, 6.1, 8)
- `project-docs/docs/phases/phase-21/contracts.md`
- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`
- `project-docs/docs/phases/phase-21/planning-audit.md`
- `project-docs/docs/phases/phase-21/qa/validation-contract.md`
- `project-docs/docs/phases/phase-21/qa/group-1-checkpoint-contract.md`
- `project-docs/docs/phases/phase-21/qa/qa-lead-handoff.md`
- `project-docs/docs/phases/phase-21/tasks/task-1-*.md` .. `task-22-*.md`
- `project-docs/templates/qa-plan-template.md`

## 2. Validation Contract Refinement Check

### 2.1 Coverage rollup

Re-checked against the DEC-158 contract at `b367a3b`.

- All in-scope requirements referenced by at least one `VAL-NNN`: **YES**
  (`planning-audit.md` `UNCOVERED_VALS=0`, re-derived in the first review:
  87 unique contract IDs, 85 task-referenced, 3 phase-control owners,
  0 uncovered).
- Every `VAL-NNN` has one validator role: **YES** (QA Reviewer 78, QA Lead 8,
  Human 1). DEC-158 changed no validator.
- Every `VAL-NNN` has a concrete evidence method: **YES**. The four gaps from
  the first review are closed: VAL-046 now has the bound containment lexicon in
  task 13; VAL-073 now has the committed corpus plus the non-vacuity control;
  VAL-083 now has the reference-instance identity and the three-sample rule;
  VAL-086 now has the exact 20-Event and one-`envelope.created` control.
- Every `VAL-NNN` maps to a task/window or phase-end check: **YES** - section 7
  maps each ID exactly once, with no merged IDs and the contract validator role
  preserved.
- Every `VAL-NNN` describes one observable, or human fan-out is required:
  **YES**. VAL-036 is now stated as one parent perception assertion, "a Human
  can safely and correctly interpret one complete rendered Gate session", with
  the contract's Coverage Audit explicitly requiring the `HUMAN-01a..01d`
  fan-out under Workflow section 6.1. Section 9 implements exactly that.
- Conditional or binary assertions require both qualifying and non-qualifying
  states: **YES**. VAL-073 now requires the complete-member positive control
  alongside the snake_case absence; VAL-086 now requires the exact pre-rebuild
  count before the post-rebuild comparison; VAL-044 and VAL-045 are explicitly
  cross-paired so neither can be satisfied unpaired.
- Parity or "matches reference" assertions have deterministic conformance
  proof: **YES** - VAL-069, VAL-081, VAL-082, and VAL-086 are byte or count
  comparisons; commands are in section 7.
- Revisions requested from Tech Lead: **none outstanding**. Section 13 records
  R-01..R-18 as applied, R-09 as applied with a Tech Lead correction that QA
  accepts and has adopted, and R-19 as deliberately deferred to Human
  execution approval.

### 2.2 Blocker candidate disposition

Both first-review candidates are closed. QA allocated no `ESC` identity.

| Candidate | Disposition | Verification in this session |
|---|---|---|
| BC-01 - 100k rebuild threshold and undefined reference hardware | **Resolved by DEC-158** without changing the protected 30s budget. The reference instance is now named in `02_TRD.md` section 9 and task 21, and the pass rule is three consecutive fresh one-iteration processes each self-asserting `<30 s`, with best-of and median tolerance forbidden and causal optimization assigned to task 21 scope. | Re-measured at planning base with the exact DEC-158 procedure on an unloaded machine: `23727938194 ns/op` / `23728 ms/rebuild`, `23797261352 ns/op` / `23797 ms/rebuild`, `23800483588 ns/op` / `23800 ms/rebuild`; all three terminated `exit_code: 0` with `ok github.com/mbell-nzl/fid/internal/journal 25.5-25.7s`. All samples are under budget with about 6.2s headroom and a 0.07s spread. Every recorded identity literal matches the machine exactly (see section 6). |
| BC-02 - invalid and self-referential report-lint command in all 22 task specs | **Resolved by DEC-158.** Report lint is removed from every Developer command gate and is now the Phase Runner's two-argument post-commit pre-dispatch control, recorded in the phase brief, `08_QUALITY_PLAN.md` section 8, the checkpoint preconditions, and VAL-087. | `grep -l "report-lint.sh" project-docs/docs/phases/phase-21/tasks/*.md` terminated `exit_code: 1` with no output: no task spec references it. The two-argument form matches the tool's documented usage and its `--self-test` passes. |

**Correction to first-review evidence.** The first review reported three
planning-base samples of `33.461764423s` (FAIL), `27109172011 ns/op`, and
`24965037891 ns/op`. The 33.46s sample was taken while a backgrounded
`go test -race ./...` was still running on the same machine, so it was a
loaded-machine measurement rather than a clean one. The clean three-sample
re-measurement above supersedes it. The first review's conclusion that the
budget needed an explicit deterministic rule and a named reference instance
still stands and is what DEC-158 supplies; the conclusion that the budget was
unreachable does not stand and is withdrawn.

### 2.3 Readiness statement

**Ready.** Every VAL has exactly one mapping with its contract validator role
preserved; all commands are concrete, bounded, and were executed fresh in a QA
Lead session; every expected non-zero command has a declared and observed exit
code; all seven review windows plus the independent Group 1 checkpoint are
defined with VAL-042 and VAL-043 separately decided; the Group 2 feasibility
harness is executable and environment-scoped; final QA is complete; and VAL-036
is the sole Human assertion, fanned out as the contract now requires. No
Must-fix finding remains.

Task 1 stays unauthorized until the Human approves execution.

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| 1 | Migration 0013 additive operational schema | VAL-001 | medium | Reversible up/down/up inventory only; no 0014, no Event apply. |
| 2 | Event v2, idempotency, role-free bootstrap | VAL-003..008 | high | Journal contract change; concurrency and restart/rebuild replay. |
| 3 | Capability discovery | VAL-009, VAL-010 | medium | Binary availability; no optimistic rows. |
| 4 | Activation trust lifecycle | VAL-011..015 | high | Sole trust source; caller self-declaration must reject. |
| 5 | Gate Registration | VAL-016..018 | high | Creates `internal/evidence`; ordering before Envelope; immutable versions. |
| 6 | Envelope values | VAL-019..022 | high | Closed input surface; null budgets; deferred-input fail-closed. |
| 7 | Generic authority evaluator | VAL-023, VAL-024 | high | Deny-by-default; role-shaped attribution must not authorize. |
| 8 | Issue Actor/Event convergence | VAL-025..027 | high | Retained mutation surface plus CF-020 and CF-027. |
| 9 | Retained families, policy isolation | VAL-028..031 | high | Widest blast radius; judgment mutations fail closed. |
| 10 | Generic Claims | VAL-032..034 | high | `work_claims` until 0014; version lifecycle. |
| 11 | Generic Runs | VAL-035, VAL-037 | medium | One-way finish; no recovery or respawn. |
| 12 | Execution-bound transition, shared evaluator | VAL-038..041 | high | Dual-call evaluator identity; CF-032. Group 1 exit. |
| checkpoint | Independent Group 1 gate | VAL-042, VAL-043 | high | Blocks task 13; fresh Group 2 feasibility. |
| 13 | Portable executor policy | VAL-044..046 | high | Real process adapter; no containment claim. |
| 14 | Live framing, redaction, controls | VAL-036, VAL-047..050 | high | Credential leakage risk; Human fixture source. |
| 15 | Attempt and Receipt | VAL-051..054 | high | Durability before spawn; five terminal outcomes. |
| 16 | Artifacts | VAL-055..057 | high | Path safety, atomic publication, availability. |
| 17 | Evidence evaluation | VAL-058..063 | high | Independent dimensions; deterministic reason order. |
| 18 | Terminal guard, completion, resume | VAL-064..068 | high | The DEC-156 guard; atomic multi-reason failure. |
| 19 | Export create and validate | VAL-069..071 | high | Byte determinism; non-mutating validation. |
| 20 | V1 removal and migration 0014 | VAL-002, VAL-074, VAL-075 | high | Destructive and irreversible; classified inventory. |
| 21 | Final-schema init, maintenance, rebuild, parity | VAL-076..083 | high | Post-0014 only; carries BC-01. |
| 22 | Final wire and real-process journey | VAL-072, VAL-073, VAL-084..086 | high | Final binary, no mocks, mutate once. |
| phase-end | Final QA | VAL-087 | high | Cadence and captured-termination audit. |

## 4. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Migration golden DDL inventory plus up/down/up integration output | VAL-001, VAL-002 | Developer | tasks 1, 20 |
| EVID-02 | Golden Event-v2 fixture corpus with byte comparison and unknown-key rejection | VAL-003, VAL-008 | Developer | task 2 |
| EVID-03 | Sequence, replay, digest-conflict, and barrier-concurrency test output | VAL-004..007 | Developer | task 2 |
| EVID-04 | Full capability-matrix JSON golden plus false-row conformance table | VAL-009, VAL-010 | Developer | task 3 |
| EVID-05 | Activation lifecycle and adversarial closed-input matrices | VAL-011..015 | Developer | task 4 |
| EVID-06 | Registration ordering, one-field policy rejection, digest/version tests | VAL-016..018 | Developer | task 5 |
| EVID-07 | Envelope boundary, Obligation round-trip, null-budget, deferred-input matrices | VAL-019..022 | Developer | task 6 |
| EVID-08 | Single-dimension Grant matrix and seeded-role adversarial test | VAL-023, VAL-024 | Developer | task 7 |
| EVID-09 | Retained-mutation corpus, fixed-clock activity test, size inventory | VAL-025..027 | Developer | task 8 |
| EVID-10 | Per-family conformance rows, authoring negatives, fail-closed writes, stats grouping | VAL-028..031 | Developer | task 9 |
| EVID-11 | Claim acquire/conflict/expiry/version matrices | VAL-032..034 | Developer | task 10 |
| EVID-12 | Run start binding and terminal/absence tests | VAL-035, VAL-037 | Developer | task 11 |
| EVID-13 | Nonterminal predicate matrix, all-reasons golden, dual-call evaluator golden | VAL-038..041 | Developer | task 12 |
| EVID-14 | Group 1 artifact/commit audit and dependency/call-graph isolation proof | VAL-042, VAL-043 | QA Lead (independent reviewer) | checkpoint |
| EVID-15 | Fresh session-captured Group 2 environment probe log (FEAS-01..FEAS-18) | no VAL | QA Lead (independent reviewer) | checkpoint |
| EVID-16 | Observed process-policy probe, one-field rejection matrix, containment lexicon absence | VAL-044..046 | Developer | task 13 |
| EVID-17 | Interleaved frame golden, control-variant golden, canary scan, high-output test | VAL-047..050 | Developer | task 14 |
| EVID-18 | Committed deterministic terminal fixture files regenerated by golden test | VAL-036 | Developer | task 14 |
| EVID-19 | Append/spawn trace, success Receipt, five-outcome matrix, unsatisfied-Attempt test | VAL-051..054 | Developer | task 15 |
| EVID-20 | Artifact content/digest/metadata inspection, adversarial candidate matrix, loss/restore lifecycle | VAL-055..057 | Developer | task 16 |
| EVID-21 | Evidence binding matrix, three paired dimension fixtures, multi-reason golden, Attention absence proof | VAL-058..063 | Developer | task 17 |
| EVID-22 | Positive/negative guarded transition, evaluator-identity completion, atomic failure, fresh-process resume | VAL-064..068 | Developer | task 18 |
| EVID-23 | Repeated archive digest and entry inventory, corruption corpus with before/after digest, import absence inventory | VAL-069..071 | Developer | task 19 |
| EVID-24 | Classified live inventory report plus predeclared absence commands | VAL-074, VAL-075 | Developer | task 20 |
| EVID-25 | Final-binary temp-project init/rebuild snapshots, corrupt fixtures, parity snapshots, export digest comparison | VAL-076..082 | Developer | task 21 |
| EVID-26 | Session-captured bounded 100k benchmark with hardware/toolchain identity | VAL-083 | QA Lead | task 21 window |
| EVID-27 | Committed conformance corpus, closed-value inventory, kebab-case scan | VAL-072, VAL-073 | Developer plus QA Lead | task 22 and final QA |
| EVID-28 | Tagged real-process journey transcript, negative matrix, post-rebuild Event-count assertion | VAL-084..086 | Developer plus QA Lead | task 22 and final QA |
| EVID-29 | Phase artifact and run-ledger audit with captured termination | VAL-087 | QA Lead | final QA |
| EVID-30 | Human observation of the recorded terminal fixture in a real terminal | VAL-036 | Human | final human checklist |

## 5. Deterministic Gates

Concrete commands only; no placeholder remains. Working directory for every
command in this plan is the repository root `/home/mark/projects/fid` unless a
row names a different directory. No numeric coverage threshold is defined for
this project (`08_QUALITY_PLAN.md` section 7 treats coverage as a diagnostic);
the required evidence is the observed `total:` line.

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | no files listed, `exit 0 (no output)`, terminated `exit_code: 0` | every task |
| Lint | `golangci-lint run ./...` | last line `0 issues.`, terminated `exit_code: 0` | every task |
| Unit/race + coverage profile | `go test -race -coverprofile=coverage.out ./...` | per-package `ok` lines, terminated `exit_code: 0`; `exit 0 (no output)` is never a valid token | every task |
| Coverage report | `go tool cover -func=coverage.out` | final `total:` line with a percentage, terminated `exit_code: 0` | every task |
| Integration journeys | `go test -tags integration -race ./cmd/fid/integration/...` | `ok github.com/mbell-nzl/fid/cmd/fid/integration <t>s`, terminated `exit_code: 0` | tasks 12, 18, 21, 22; checkpoint; final QA |
| Benchmarks | `go test -bench . -benchmem ./...` | benchmark result lines, terminated `exit_code: 0`; the per-sample `<30 s` `b.Fatalf` remains part of this gate per DEC-158 | tasks 21, 22; final QA |
| Bounded 100k rebuild, three consecutive fresh processes | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal`, run three times as three separate processes | three distinct `BenchmarkRebuildWorkspace100kEvents-<n>` lines, each reporting `<30 s` and each terminated `exit_code: 0`, on the recorded DEC-158 reference instance; best-of and median sampling are forbidden | task 21, final QA |
| Reference-instance identity | `go version`; `go env GOOS GOARCH`; `uname -srm`; `nproc`; `grep -m1 '^model name' /proc/cpuinfo`; `grep -m1 '^MemTotal:' /proc/meminfo` | each terminated `exit_code: 0` and matching the `02_TRD.md` section 9 reference instance | task 21, final QA |
| Build | `go build -o /tmp/fid ./cmd/fid` | `exit 0 (no output)`, terminated `exit_code: 0` | tasks 21, 22; final QA |
| Report lint (Phase Runner post-commit control, two-argument form; not a Developer task-gate command) | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-<N>-done.md project-docs/docs/phases/phase-21/tasks/task-<N>-<slug>.md` | `PASS report-lint completed with no failures`, terminated `exit_code: 0` | after each report commit, before every next dispatch, all 22 tasks |
| Full phase command gate | the seven `AGENTS.md` commands in order: `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./cmd/fid/integration/...`; `go test -bench . -benchmem ./...`; `go build -o /tmp/fid ./cmd/fid` | every command terminates `exit_code: 0` with its own captured token | final QA |

### 5.1 Expected non-zero commands used by this plan

Each command below passes only on the declared non-zero status. They are given
as fenced blocks, not table cells, so they can be copied byte-exactly.

VAL-046, no production containment claim, expected `exit_code: 1`:

```text
rg -n -i --glob '*.go' --glob '!**/*_test.go' -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' cmd internal
```

VAL-013, no activation recovery or automatic-respawn path (task 4), expected
`exit_code: 1`:

```text
rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"recover|newActivationRecover|activation\.recover' cmd/fid internal
```

VAL-046, no Phase 21 contract containment claim (task 13), expected
`exit_code: 1`:

```text
rg -n -i -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' project-docs/docs/phases/phase-21/contracts.md project-docs/docs/phases/phase-21/public-interface.md
```

VAL-063, no production Attention state or effect (task 17), expected
`exit_code: 1`:

```text
rg -n -i --glob '*.go' --glob '!**/*_test.go' '\battention\b|attention_items' cmd/fid internal
```

VAL-071, no export import command or path (task 19), expected `exit_code: 1`:

```text
rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"import|newExportImport|export\.import' cmd/fid internal
```

VAL-075, the six task-20 absence searches, expected `exit_code: 1` each:

```text
rg -n --glob '*.go' --glob '!**/*_test.go' 'String(Var)?\([^\n]*"by"|GetString\("by"|--by' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' 'policy\.Authorize|KnownRoles|ActorRole|Role:' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' '(INTO|UPDATE|DELETE FROM) (roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)\b|Table:\s*"(roles|gate_runs|task_completions|agent_runs|ac_verifications|concerns)"' cmd/fid internal
rg -n 'EventType:\s*"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)|"event_type":"(dev\.|issue\.completed|issue\.claim\.|run\.recorded)' cmd/fid internal testdata
rg -n --glob '*.go' --glob '!**/*_test.go' '"event"\s*:\s*"[a-z0-9]+_[a-z0-9_]+"|Event:\s*"[a-z0-9]+_[a-z0-9_]+"' cmd/fid internal
```

VAL-075, the five task-20 removed-surface checks, expected `exit_code: 2` each:

```text
ls internal/dev
ls internal/runs
ls cmd/fid/dev.go
ls cmd/fid/policy.go
ls cmd/fid/issue_complete.go
```

VAL-073, the closed-value snake_case scan of section 12, expected
`exit_code: 1`:

```text
grep -E $'\t[^\t]*[a-z0-9]_[a-z0-9]' "$SCRATCH/kept.tsv"
```

VAL-072 and VAL-073 positive control, expected `exit_code: 0`, so the corpus
cannot be missing or empty:

```text
ls cmd/fid/integration/testdata/p21_conformance
```

VAL-074, VAL-043, and QA-93 classification recipes, per DEC-158 and
`v1-convergence-inventory.md` section 7. At final task-20 HEAD the **first two**
recipes are absence controls that must terminate `exit_code: 1` and are already
declared above; the **last two** are broad classified inventories that must
terminate `exit_code: 0`, keep matching retained, frozen, and diagnostic
readers, and must not be narrowed to suppress them. Their pass condition is
that every produced line maps to an explicit retained, frozen, diagnostic, or
removal disposition and the report records `UNCLASSIFIED=0`. At the Group 1
checkpoint all four terminate `exit_code: 0`, because removal has not happened
yet:

```text
rg -n --glob '*.go' --glob '!**/*_test.go' 'policy\.Authorize|KnownRoles|ActorRole|Role:' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' 'String(Var)?\([^\n]*"by"|GetString\("by"|--by' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' '(FROM|INTO|UPDATE|JOIN|DELETE FROM) (roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)\b|Table:\s*"(roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)"' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' 'EventType:\s*"' cmd/fid internal
```

## 6. Command And Feasibility Validation

Heavy mode: no result is carried forward from the Phase 20 QA plan. Rows marked
`[R2]` were validated fresh in this re-review session at `b367a3b`, which covers
every command DEC-158 added or changed. Every unmarked row, and any row marked
`[R1]`, was validated fresh in the first QA Lead session at
`463183a16ac99879eeeadc35ff5f1064980892e5` against commands DEC-158 left
unchanged. Working directory for every row is the repository root unless
stated.

| Referenced Step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | executed | PASS | no output; terminated `exit_code: 0` |
| `golangci-lint run ./...` | version probe plus scoped execution `golangci-lint run ./internal/db/...` | PASS | `golangci-lint has version 2.12.2 built with go1.26.2`; scoped run printed `0 issues.`, terminated `exit_code: 0` |
| `go test -race -coverprofile=coverage.out ./...` | executed with profile redirected to `/tmp/qa-probe-coverage.out` | PASS | 21 `ok` package lines, last `ok github.com/mbell-nzl/fid/internal/workflow 22.250s coverage: 80.6% of statements`; terminated `exit_code: 0` |
| `go tool cover -func=coverage.out` | executed against the probe profile | PASS | `total: (statements) 78.7%`; terminated `exit_code: 0` |
| `go test -tags integration -race ./cmd/fid/integration/...` | executed | PASS | `ok github.com/mbell-nzl/fid/cmd/fid/integration 4.658s`; terminated `exit_code: 0`. Declared journey budget: 300s wall clock |
| `go build -o /tmp/fid ./cmd/fid` | executed | PASS | silent success; terminated `exit_code: 0` |
| `go test -bench . -benchmem ./...` | component probe of its longest member (bounded 100k rebuild) `[R1][R2]` | PASS | the member benchmark passes cleanly on an unloaded machine; see the next row |
| DEC-158 rule: `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` as three consecutive fresh processes | executed as three separate unloaded processes `[R2]` | PASS | `23727938194 ns/op 23728 ms/rebuild`; `23797261352 ns/op 23797 ms/rebuild`; `23800483588 ns/op 23800 ms/rebuild`; each terminated `exit_code: 0` with `ok github.com/mbell-nzl/fid/internal/journal` at 25.566s / 25.690s / 25.660s. All three under the 30s budget, spread 0.07s, headroom about 6.2s |
| first-review benchmark samples (superseded) | executed `[R1]` | withdrawn | `33.461764423s` FAIL, then `27109172011 ns/op` and `24965037891 ns/op`. The failing sample overlapped a backgrounded `go test -race ./...`; it is a loaded-machine measurement and is superseded by the clean row above |
| DEC-158 reference-instance identity: `go version`; `go env GOOS GOARCH`; `uname -srm`; `nproc`; `grep -m1 '^model name' /proc/cpuinfo`; `grep -m1 '^MemTotal:' /proc/meminfo` | all six executed verbatim from the task-21 command block `[R2]` | PASS | `go version go1.26.2 linux/amd64`; `linux` / `amd64`; `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64`; `6`; `model name : AMD Ryzen 7 5700G with Radeon Graphics`; `MemTotal: 8133236 kB`. Every literal recorded in `02_TRD.md` section 9, task 21, and DEC-158 matches the observed machine exactly |
| `project-docs/tools/report-lint.sh <report>` one-argument form | executed against `task-1-done.md` `[R1]` | **FAIL, now removed** | printed `Usage: ... <completion-report.md> <task-spec.md> [--repo <git-dir>]`, terminated `exit_code: 2`. DEC-158 removed the command from all 22 specs |
| report-lint absent from every task spec after DEC-158 | `grep -l "report-lint.sh" project-docs/docs/phases/phase-21/tasks/*.md` `[R2]` | PASS | no output, terminated `exit_code: 1`: 0 of 22 specs still declare it |
| report-lint behaviour with task 21's three identical benchmark lines | source inspection of `required_commands()` and the coverage loop `[R2]` | PASS with a stated limitation | the loop is a stateless per-command lookup with no dedupe, so duplicates neither break lint nor prove three samples were run. Counting the three distinct samples is QA-80's obligation, recorded in section 7.8 |
| `project-docs/tools/report-lint.sh --self-test` | executed | PASS | `SELF-TEST PASS report-lint self-test completed with no failures`; terminated `exit_code: 0` |
| `project-docs/tools/report-lint.sh <report> <task-spec>` two-argument form | documented usage inspection plus argument-count guard inspection in the script (it exits 2 with usage when fewer than two arguments are supplied) | PASS | script header `Usage: project-docs/tools/report-lint.sh <completion-report.md> <task-spec.md> [--repo <git-dir>]` |
| task 4 recovery-absence recipe, run verbatim from the spec (VAL-013) | executed `[R2]` | PASS | no match, terminated `exit_code: 1`; declared `exit_code: 1` in the task's new `## Expected Non-Zero Commands` |
| task 17 Attention-absence recipe, run verbatim from the spec (VAL-063) | executed `[R2]` | PASS | no match, terminated `exit_code: 1`; declared `exit_code: 1` |
| task 19 import-absence recipe, run verbatim from the spec (VAL-071) | executed `[R2]` | PASS | no match, terminated `exit_code: 1`; declared `exit_code: 1` |
| task 13 containment lexicon over the two phase contracts (VAL-046) | executed `[R2]` | PASS | no match, terminated `exit_code: 1`; declared `exit_code: 1` |
| task 20 broad table-reader classification recipe | executed `[R2]` | PASS | 54 matching lines, terminated `exit_code: 0`; correctly **not** declared in the task's expected-non-zero section |
| task 20 broad `EventType` classification recipe | executed `[R2]` | PASS | 22 matching lines, terminated `exit_code: 0`; correctly **not** declared |
| task 20 expected-non-zero declaration set after DEC-158 | declaration count and content inspection `[R2]` | PASS | exactly the 10 original absence declarations remain; the two added broad inventories are absent from the section, matching DEC-158's corrected R-09 premise |
| no duplicate required command outside task 21 | per-spec duplicate-line scan over each `## Exact Commands To Run` block `[R2]` | PASS | only `task-21-maintenance-rebuild.md` reports a duplicate line, which is the intended three-sample benchmark |
| VAL-086 expected count of 20 Events | independent enumeration of the journey order in `phase_brief.md` and `contracts.md` section 16 against the Event families in section 14 `[R2]` | PASS | QA derived 20 Event-producing operations independently before reading the amended contract, and the amended `contracts.md` enumerates the identical 20 with exactly one `envelope.created`; reads (capabilities, evaluate, resume, export create/validate) produce none |
| containment lexicon over source (VAL-046) | executed | PASS | no match, terminated `exit_code: 1`; the accepted permissive values `sandbox:"none"`, `network:"allowed"`, `write:"allowed"` are not matched |
| containment lexicon over `contracts.md` and `public-interface.md` | executed | PASS | no match, terminated `exit_code: 1` |
| `--by` inventory (task 20) | executed at planning base | PASS (syntax and reachability) | 61 matching lines today; the pattern is valid and the paths exist, so `exit_code: 1` at task 20 is a real absence proof |
| section 5.1 role-policy classification recipe | executed at planning base | PASS (syntax and reachability) | 211 matching lines today |
| removed/frozen writer inventory (task 20) | executed at planning base | PASS (syntax and reachability) | 10 matching lines today |
| V1 Event fixture inventory (task 20, includes `testdata`) | executed at planning base | PASS | 2 matching lines today (`cmd/fid/issue_transition_test.go:921`, `internal/issues/transition_test.go:1442`); `ls -d testdata` terminated `exit_code: 0`, and a no-match run over the same three paths terminated `exit_code: 1`, so the declared expectation is reachable and not a path typo |
| snake-case Activity Event inventory (task 20) | executed at planning base | PASS (syntax and reachability) | 10 matching lines today |
| classified inventory recipes, `v1-convergence-inventory.md` section 7 (all four) | executed at planning base | PASS | 211 / 61 / 54 / 22 matching lines respectively, all terminated `exit_code: 0`; these are classification commands, not absence commands |
| `ls <missing-path>` exit semantics | executed | PASS | `ls /nonexistent-path-probe` terminated `exit_code: 2`; `ls internal/dev` terminated `exit_code: 0`. Task 20's declared `exit_code: 2` is correct for GNU coreutils here |
| VAL-073 closed-value inventory pipeline (`jq` path/value extraction, exclusion filter, snake scan) | executed against a seeded sample containing `spawn_failed`, `in_progress`, `needs_review`, `head_oid`, `tree_oid changed` | PASS | pipeline isolated the seeded violation and excluded field names, opaque IDs, opaque workflow states, and free text; `jq --version` present at `/usr/bin/jq` |
| Group 2 probe harness: offline stdlib scratch Go module | executed in `$(mktemp -d)` | PASS with a required flag | `go vet ./...` clean; `go build ./child` first failed `error obtaining VCS status: exit status 128`, so every scratch build/run in this plan uses `-buildvcs=false`; rebuild terminated `exit_code: 0` |
| FEAS structured argv, separate interleaved streams | executed | PASS | child printed `out-1` / `err-1` / `out-2` on the correct descriptors, terminated `exit_code: 0` |
| FEAS invalid UTF-8 and control bytes | executed | PASS | `xxd` showed `41ff fe1b 5b33 316d 0d0a` (0xFF, 0xFE, ESC `[31m`, CR, LF) reaching the reader unmodified |
| FEAS timeout, signal, spawn failure | executed | PASS | `timeout -s TERM 2` terminated `exit_code: 124`; `timeout -s KILL 2` terminated `exit_code: 124`; missing executable terminated `exit_code: 127`; note the first probe used `select{}` and produced a Go `fatal error: all goroutines are asleep`, so the harness in section 11 uses `time.Sleep` |
| FEAS high-output draining | executed | PASS | 8,388,608 bytes drained from stdout while stderr was concurrently written; `wc -c` reported `8388608` |
| FEAS closed environment injection | executed | PASS | `env -i PROBE_MODE=env PATH=/usr/bin HOME=/tmp` produced exactly those three variables |
| FEAS credential canary | executed | PASS | child echoed `token=CANARY-abc123` from its own environment, so a canary is observable on a real surface |
| FEAS symlink-resolved containment | executed | PASS | `realpath` resolution of `cwd/escape.txt` reported `ESCAPES` against the `cwd` base |
| FEAS changed-during-copy | executed | PASS | pre/post SHA-256 differed for a file mutated between hash and copy |
| FEAS atomic content publication | executed | PASS | same-filesystem `mv` from `store/tmp` to `store/objects/<digest>` left `store/tmp` empty and the digest-named object present |
| FEAS Artifact loss and restoration | executed | PASS | `MISSING_OBSERVABLE: PASS` then `RESTORED_OBSERVABLE: PASS` |
| FEAS byte-deterministic ZIP | executed | PASS | two archives written 1.1s apart with fixed 1980-01-01 headers produced identical digest `018486bea8ba8921b104146842daec9577bde29c0c8cb129a482bf4a3d0c1d3f`; `ZIP_DETERMINISM: PASS` |
| FEAS non-mutating corruption validation | executed | PASS | `unzip -t` terminated `exit_code: 0` and left the digest unchanged (`VALIDATION_NON_MUTATING: PASS`); after a one-byte `dd` corruption it terminated `exit_code: 1` |
| FEAS repository-tree freshness | executed | PASS | two commits in a scratch repository produced different `rev-parse HEAD^{tree}` values |
| Human step `cat cmd/fid/testdata/val-036/human-session.txt` | task-spec declaration check `[R2]` | PASS as planned, produced by task 14 | task 14 now names `cmd/fid/testdata/val-036/human-session.txt` and `live-frames.jsonl` in Files To Create, and its acceptance criterion requires a golden test to regenerate and byte-compare both, with interleaved streams, all three control kinds, a non-success terminal frame, and raw ESC/CR/BEL in the child output. The step itself is a single `cat` in an interactive terminal and needs no tooling |
| `ls cmd/fid/integration/testdata/p21_conformance` | task-spec declaration check `[R2]` | PASS as planned, produced by task 22 | task 22 now names the directory in Files To Create, requires `manifest.json`, `operations/<operation>.json` per section 4.1-4.2 row, `failures/<failure-code>.json` per section 7 failure, and `live/<variant>.jsonl` per section 8 variant, and runs `ls` on it as a required exit-0 command. It is correctly not declared expected-non-zero |
| Task 5/13/16 scope `./internal/evidence`, task 19/21 scope `./internal/export` | directory listing of `internal/` | PASS with sequencing note | neither package exists at planning base; task 5 creates `internal/evidence` and task 19 creates `internal/export`, so their task-scoped `go test` commands are only valid from the task that creates the package onward. No earlier task references them |

No unvalidated or invalid step remains. Every command DEC-158 added was
executed verbatim from its task spec in this session, and both previously
blocked artifact assumptions are now declared by their owning task spec.

## 7. Deterministic Checklist

Every item is executed or verified by the QA Lead (or, in a review window, by
the named QA Reviewer, with the QA Lead re-verifying at phase end). Each
`QA-NN` cites the `VAL-NNN` it verifies. Every VAL-001..087 appears exactly
once across section 7 and section 9. PASS without cited evidence is invalid;
code reading alone never satisfies an item.

Shared conventions for every row:

- Working directory is the repository root unless the row names another.
- Scratch state is allocated fresh per attempt with `mktemp -d`; cleanup is
  never recorded as evidence.
- Every recorded result carries one verbatim token plus the observed
  `exit_code`.
- `WINDOW-k` marks the review window that first verifies the row.

### 7.1 Window 1 - tasks 1-3

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Reproduce the 0013 schema inventory and up/down/up round trip: `go test -race ./internal/db/...` `(qa-lead)` | VAL-001 | `ok github.com/mbell-nzl/fid/internal/db ...`, terminated `exit_code: 0`, plus the golden DDL diff showing exact tables, indexes, checks, and foreign keys and an unchanged schema after down/up |
| QA-02 | Byte-compare every converted family against the closed Event-v2 top-level shape and ordered Effects, and confirm unknown top-level keys reject: `go test -race ./internal/journal` `(qa-lead)` | VAL-003 | golden fixture comparison output plus the named unknown-key rejection test, terminated `exit_code: 0` |
| QA-03 | Confirm aggregate sequence starts at 1, has no gap, and rebuild rejects duplicate and out-of-order facts: `go test -race ./internal/journal` `(qa-lead)` | VAL-004 | live and replay sequence test names with results, terminated `exit_code: 0` |
| QA-04 | Confirm same key and same digest returns the original result and Event with only `replayed:true`, surviving restart and rebuild, with no new Event or side effect: `go test -race ./internal/journal ./internal/cli` `(qa-lead)` | VAL-005 | replay test output plus the asserted Event-count delta of zero, terminated `exit_code: 0` |
| QA-05 | Confirm same key and different digest returns `idempotency-conflict` with no Event and no effect `(qa-lead)` | VAL-006 | digest-conflict test output naming code `idempotency-conflict` and the exact section 6 details object, terminated `exit_code: 0` |
| QA-06 | Confirm concurrent same-key calls commit exactly one mutation: `go test -race ./internal/journal ./internal/cli` `(qa-lead)` | VAL-007 | barrier concurrency test output showing one Event identity and one replay under `-race`, terminated `exit_code: 0` |
| QA-07 | Byte-compare the `project.initialized` fixture against `contracts.md` section 14: service Actor `fid-bootstrap`, aggregate sequence 1, null repository, null idempotency, empty Effects, and no role/workflow/policy Effect `(qa-lead)` | VAL-008 | golden bootstrap Event diff, terminated `exit_code: 0` |
| QA-08 | Compare `fid capabilities --json` output against the exact 12-row sorted matrix in `contracts.md` section 2, field by field: `go test -race ./internal/workflow ./cmd/fid/...` `(qa-lead)` | VAL-009 | full-matrix golden diff plus the projection rows, terminated `exit_code: 0` |
| QA-09 | Confirm `kit.activation.v1`, `work.claims.v1`, `export.portable.v1`, `reviews.coverage.v1`, `human.outcomes.v1`, `telemetry.trajectory.v1`, and `attention.v1` all stay `false` across the internal Activation lifecycle, Claim lifecycle and list, export create and validate, fail-closed AC and Concern judgment, historical stats, and Attention-absence surfaces `(qa-lead)` | VAL-010 | exact false-row/operation conformance matrix naming each of the seven capabilities against each named surface, terminated `exit_code: 0` |

### 7.2 Window 2 - tasks 4-6

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-10 | Confirm only `prepared -> staged -> committed -> finalized` is accepted and one active identity results; every skipped or extra transition rejects `(qa-lead)` | VAL-011 | lifecycle state-table output covering accepted and rejected orders, terminated `exit_code: 0` |
| QA-11 | Confirm abort succeeds only from prepared or staged, leaves no active trust source, and abort after commit returns `activation-recovery-unavailable` `(qa-lead)` | VAL-012 | abort state-table output plus the post-abort trust-source read, terminated `exit_code: 0` |
| QA-12 | Confirm recovery and automatic respawn are unavailable and post-commit repair rejects `(qa-lead)` | VAL-013 | the task-4 declared recovery-absence recipe (section 5.1) terminated `exit_code: 1`; the `activation recover` unknown-command result; the `kit.activation.v1` row with `recover:false` and `automatic_respawn:false`; and the typed `activation-recovery-unavailable` post-commit rejection |
| QA-13 | Confirm the finalized authority profile is the sole source of Actor identity, capability, producer, and provenance; a non-finalized or alternative source resolves nothing `(qa-lead)` | VAL-014 | resolution/provenance test output for both the qualifying and non-qualifying source, terminated `exit_code: 0` |
| QA-14 | Confirm caller-supplied capability or provenance fields are unknown input, reject with a typed error, and append no Event `(qa-lead)` | VAL-015 | adversarial closed-input matrix output plus the asserted zero-Event delta, terminated `exit_code: 0` |
| QA-15 | Confirm Registration resolves only a finalized activated Definition and only before Envelope creation; missing, mutable, or caller-supplied values reject `(qa-lead)` | VAL-016 | ordering/trust integration test output for both states, terminated `exit_code: 0` |
| QA-16 | Confirm every unsupported executor-policy value rejects `executor-policy-unsupported` per field with no Registration persisted `(qa-lead)` | VAL-017 | one-field rejection matrix naming each rejected field and the post-attempt Registration count of zero, terminated `exit_code: 0` |
| QA-17 | Confirm the same Definition and binding resolve one immutable Registration while a changed binding creates the next version `(qa-lead)` | VAL-018 | digest/version test output showing reuse and increment, terminated `exit_code: 0` |
| QA-18 | Confirm Envelope creation accepts exactly revision 1, one primary Issue, zero related Issues, and one prior Registration-bound Obligation; every wider or incomplete selection rejects `(qa-lead)` | VAL-019 | boundary matrix output for both states, terminated `exit_code: 0` |
| QA-19 | Confirm the one `gate-evidence` Obligation names both guarded actions and rejects computed status, caller-supplied status, and any other type `(qa-lead)` | VAL-020 | closed-value round-trip diff against `contracts.md` section 1 plus the rejection rows, terminated `exit_code: 0` |
| QA-20 | Confirm Envelope and Run budgets are immutable with `consumed:null` and `remaining:null` across create, start, read, and rebuild `(qa-lead)` | VAL-021 | four-stage snapshot comparison, terminated `exit_code: 0` |
| QA-21 | Confirm revision, cancellation, protected Decision, Stop, and Attention inputs fail closed as unavailable with no Event `(qa-lead)` | VAL-022 | unknown/deferred input matrix with typed codes and the zero-Event delta, terminated `exit_code: 0` |

### 7.3 Window 3 - tasks 7-9

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-22 | Confirm authorization requires exact Activation, Actor, action, capability, Issue, Envelope, and revision Grant match; each single-dimension mismatch fails closed `(qa-lead)` | VAL-023 | single-dimension grant matrix showing one permitting row and one failing row per dimension, terminated `exit_code: 0` |
| QA-23 | Confirm a project-local role-looking attribution grants no kit authority and satisfies no Evidence, Human, or protected semantic `(qa-lead)` | VAL-024 | seeded-role adversarial test output covering every protected operation, terminated `exit_code: 0` |
| QA-24 | Confirm retained Issue mutations use `--actor`, an idempotency key, Event v2, and the exact result/replay contract, and that `--by` is rejected `(qa-lead)` | VAL-025 | shared retained-mutation corpus output plus the `--by` `invalid-argument` rows, terminated `exit_code: 0` |
| QA-25 | Confirm same-clock burst Issue Events stay visible, nonfuture, and in aggregate/Event sequence `(qa-lead)` | VAL-026 | fixed-clock activity test output, terminated `exit_code: 0` |
| QA-26 | Confirm reopened Issue creation code meets the size budgets without contract drift: file at most 400 lines and functions targeting at most 60 `(qa-lead)` | VAL-027 | size inventory output for the named files plus the unchanged Issue corpus result, terminated `exit_code: 0` |
| QA-27 | Confirm retained AC add, labels, links, requirement sync, workflow define/assign, and normalized activity reads implement their inventory contracts only, with kebab-case activity values and no authority inference `(qa-lead)` | VAL-028 | per-family operation/read matrix plus authority probe rows, terminated `exit_code: 0` |
| QA-28 | Confirm workflow authoring cannot add role or judgment predicates and structural predicates never become kit authority; `--require` is absent `(qa-lead)` | VAL-029 | authoring/registry negative test output plus the `--require` `invalid-argument` row, terminated `exit_code: 0` |
| QA-29 | Confirm AC verification and concern raise/address are fail-closed unavailable with no writer and no row or Event `(qa-lead)` | VAL-030 | `capability-unavailable` exit-5 results for all three leaves plus a write inventory showing no writer, terminated `exit_code: 0` |
| QA-30 | Confirm stats accepts exactly `--group-by actor`, `--group-by epic`, or `--group-by model`, exposes historical attribution only, and rejects `--by` and any other grouping `(qa-lead)` | VAL-031 | CLI/result/source test output for accepted and rejected values, terminated `exit_code: 0` |

### 7.4 Window 4 - tasks 10-12

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-31 | Confirm Claim acquire binds Issue, Envelope, revision, Actor, worktree, scope, and TTL and permits one active mutation Claim; incomplete bindings do not commit `(qa-lead)` | VAL-032 | acquire matrix output for both states, terminated `exit_code: 0` |
| QA-32 | Confirm conflicting Claim, wrong Actor or worktree, and expired Claim fail with exact reasons and no supersession or force takeover `(qa-lead)` | VAL-033 | conflict/expiry matrix output plus the unchanged Claim state, terminated `exit_code: 0` |
| QA-33 | Confirm renew and release require the current Claim version and that release replay is idempotent for the same version and reason `(qa-lead)` | VAL-034 | version lifecycle output including a stale-version rejection, terminated `exit_code: 0` |
| QA-34 | Confirm Run start binds active Claim, Issue, Envelope revision, Actor, repository, and immutable budgets; any mismatch fails `(qa-lead)` | VAL-035 | start binding matrix output for both states, terminated `exit_code: 0` |
| QA-35 | Confirm explicit finish terminalizes one active Run once and that recovery, predecessor, continuation, and respawn stay absent `(qa-lead)` | VAL-037 | finish/repeat test output plus the absence proof for `run recover` and continuation paths, terminated `exit_code: 0` |
| QA-36 | Confirm a nonterminal execution-bound transition commits once only when Envelope, Actor, Claim, Run, repository, revision, and edge predicates all qualify: `go test -tags integration -race ./cmd/fid/integration/...` plus the unit matrix `(qa-lead)` | VAL-038 | condition matrix showing the single advancing row and one failing row per predicate, terminated `exit_code: 0` |
| QA-37 | Confirm a rejected execution-bound transition returns every applicable ordered reason with no Event and no revision change `(qa-lead)` | VAL-039 | multi-failure atomicity golden byte-compared against `contracts.md` section 10 ordering, plus the zero-Event delta, terminated `exit_code: 0` |
| QA-38 | Confirm the shared evaluator returns byte-identical unsatisfied reasons for terminal-success transition and completion before trusted Evidence exists, with no mutation `(qa-lead)` | VAL-040 | dual-call evaluator diff showing byte equality, terminated `exit_code: 0` |
| QA-39 | Confirm invalid Issue revision-effect input is rejected before mutation (CF-032) `(qa-lead)` | VAL-041 | direct invalid-input test output plus the asserted absence of any Effect, terminated `exit_code: 0` |

### 7.5 Group 1 independent checkpoint - after task 12, before task 13

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-40 | Audit that tasks 1-12 have committed completion reports with Phase Runner two-argument report-lint evidence captured after each report commit and before the next dispatch, and that QA Review windows 1-4 exist with every finding linked to its owning task and resolved by a committed repair; record Group 1 HEAD and the contract/spec commit `(qa-lead)` | VAL-042 | `git log --oneline` for the twelve task commits and their `docs(task)` report commits; `ls project-docs/docs/phases/phase-21/qa/`; the four review artifacts with `Outcome: PASS`; the run-ledger entry for each pre-dispatch lint; a fresh `project-docs/tools/report-lint.sh <report> <spec>` run per task with `PASS report-lint completed with no failures` and terminated `exit_code: 0`; `git rev-parse HEAD` |
| QA-41 | Prove no new Group 1 generic path consults the temporary V1 inventory, and classify the temporary inventory as exactly the `v1-convergence-inventory.md` section 5 list, without requiring final absence `(qa-lead)` | VAL-043 | (a) the four section 7 classification commands run fresh, every match classified against the inventory file, unclassified count 0; (b) a dependency/call-graph proof per new generic path - `go list -deps ./internal/workflow ./internal/policy ./internal/claims ./internal/evidence` plus `rg -n -e roles -e gate_runs -e task_completions -e agent_runs -e ac_verifications -e concerns -e 'internal/dev' <generic source files>` terminated `exit_code: 1` for the generic files; (c) adversarial fixtures seeding rows in each temporary table and showing Activation, Registration, Envelope, Grant, Claim, Run, nonterminal transition, capability, and idempotency paths are unchanged; (d) proof that AC verification and concern mutation are fail-closed |

Section 11 defines the fresh Group 2 feasibility probes that must also PASS in
the same checkpoint session before task 13 dispatch.

### 7.6 Window 5 - tasks 13-15

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-42 | Confirm the accepted executor policy is observably enforced: structured argv, no shell, exact cwd, closed environment, bounds, credentials, redaction inputs, and safe Artifact declarations `(qa-lead)` | VAL-044 | positive policy/process probe showing the child's observed argv, cwd, and full environment set matching the policy exactly, terminated `exit_code: 0` |
| QA-43 | Confirm each unsupported shell, sandbox, network, write, environment, path, or limit value rejects at registration with no Registration persisted `(qa-lead)` | VAL-045 | one-field rejection matrix with `executor-policy-unsupported` and the exact `fields` detail per case, terminated `exit_code: 0` |
| QA-44 | Confirm no OS, network, filesystem, privilege, namespace, or platform containment claim exists in source or in the Phase 21 contracts, and that the permissive values are explicit `(qa-lead)` | VAL-046 | negative: both task-13 declared containment lexicons from section 5.1, over `cmd internal` and over `contracts.md` plus `public-interface.md`, each terminated `exit_code: 1`; positive: the `evidence.gate-definitions.v1` limits row showing `sandbox:"none"`, `network:"allowed"`, `write:"allowed"` and the resolved Registration policy carrying the same values |
| QA-45 | Confirm stdout and stderr process bytes serialize as exact globally sequenced `fid.live.v1` data frames with correct chunking and escaping `(qa-lead)` | VAL-047 | interleaved child-process golden: frame-by-frame diff of schema, `execution_id`, strictly increasing `sequence` across both streams, 4096-byte chunk boundaries that never split a code point, and JSON-escaped CR/LF/ESC, terminated `exit_code: 0` |
| QA-46 | Confirm redaction, truncation, and Artifact rejection each emit their exact control payload and that no unknown control kind exists `(qa-lead)` | VAL-048 | control-variant golden covering `output-truncated`, `artifact-rejected`, and `redaction-applied` with exact fields and ordering, plus a rejection row for an unknown or snake-case kind, terminated `exit_code: 0` |
| QA-47 | Confirm credential and literal redaction occurs before every live, stored, hashed, error, Artifact, export, and Human surface `(qa-lead)` | VAL-049 | canary end-to-end scan: `rg -n 'CANARY-<value>' <live frames> <journal> <projection dump> <artifact store> <export archive> <human transcript>` terminated `exit_code: 1`, plus the positive `redaction-applied` control frame with its count |
| QA-48 | Confirm the output ceiling emits exactly one truncation frame, drains both pipes, stays bounded, and preserves the independently observed process outcome `(qa-lead)` | VAL-050 | high-output process test showing one truncation frame, bounded retained bytes, no hang within the declared timeout, and the outcome matching an independent observation, terminated `exit_code: 0` |
| QA-49 | Confirm the execution Attempt is durable before process spawn and that an append failure spawns nothing `(qa-lead)` | VAL-051 | injected append/spawn trace showing Attempt-before-spawn ordering for every spawn and a no-spawn result on injected append failure, terminated `exit_code: 0` |
| QA-50 | Confirm a successful process produces exactly one immutable Receipt with exact binding and terminal fields `(qa-lead)` | VAL-052 | success terminal test plus a mutation attempt showing immutability and a Receipt count of one per Attempt, terminated `exit_code: 0` |
| QA-51 | Confirm failed, timed-out, signaled, spawn-failed, and handled `fid-error` executions each produce the exact closed terminal outcome and exit classification, and that only `fid-error` lacks a Receipt `(qa-lead)` | VAL-053 | five-case terminal matrix with observed exit/signal/timeout values, process exit 8 for the four non-success outcomes, exit 7 and null Receipt for `fid-error`, terminated `exit_code: 0` |
| QA-52 | Confirm an Attempt without a qualifying successful fid-executed Receipt never satisfies Evidence, including caller-reported proof `(qa-lead)` | VAL-054 | missing-Receipt and caller-reported-proof tests both returning unsatisfied or rejected, terminated `exit_code: 0` |

### 7.7 Window 6 - tasks 16-18

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-53 | Confirm a safe declared regular Artifact is atomically stored by digest with exact sanitized metadata `(qa-lead)` | VAL-055 | content inspection integration test: stored bytes equal sanitized source bytes, SHA-256 matches the recorded digest, media type/size/retention class exact, and no partial file left in the staging directory, terminated `exit_code: 0` |
| QA-54 | Confirm every unsafe Artifact candidate emits a typed rejection and cannot satisfy required material `(qa-lead)` | VAL-056 | adversarial matrix covering `changed-during-copy`, `directory`, `missing`, `non-regular`, `outside-cwd`, `oversized`, and `symlink`, each with its `artifact-rejected` control and an unsatisfied `required-available` evaluation, terminated `exit_code: 0` |
| QA-55 | Confirm verify loss and restoration change only Artifact availability and preserve Artifact, Receipt, and Evidence identity `(qa-lead)` | VAL-057 | delete/verify/restore lifecycle snapshots showing only the `available`/`missing` field differing, terminated `exit_code: 0` |
| QA-56 | Confirm Evidence links only an exact fid-executed Receipt matching Obligation, Registration, Run, Issue, and target `(qa-lead)` | VAL-058 | binding matrix with one linking row and one rejecting row per mismatched dimension plus a caller-proof rejection, terminated `exit_code: 0` |
| QA-57 | Confirm Evidence record validity serializes and evaluates independently from existence, freshness, and Artifact availability `(qa-lead)` | VAL-059 | valid/invalid paired fixture where only the validity reason changes, terminated `exit_code: 0` |
| QA-58 | Confirm a repository-tree change stales Evidence while the original target stays fresh, and that branch labels do not drive freshness `(qa-lead)` | VAL-060 | git tree change test showing `fresh` for the original tree and `stale` after the tree OID changes, terminated `exit_code: 0` |
| QA-59 | Confirm required Artifact availability changes only the Artifact reason and status dimension `(qa-lead)` | VAL-061 | available/missing paired evaluation with an isolated field diff, terminated `exit_code: 0` |
| QA-60 | Confirm evaluation returns every reason deduplicated in the fixed section 10 precedence and bytewise tie order, stable across repeat and rebuild `(qa-lead)` | VAL-062 | repeated and post-rebuild multi-failure golden byte-identical to the reference, terminated `exit_code: 0` |
| QA-61 | Confirm unsatisfied Obligations stay queryable and never create Attention `(qa-lead)` | VAL-063 | positive: the evaluation read returning status and reasons; negative: the task-17 declared Attention-absence recipe (section 5.1) terminated `exit_code: 1`, plus the Attention table, Event-family, and `attention.v1` capability-leaf absence proofs |
| QA-62 | Confirm that after a qualifying Receipt link and evaluation the terminal-success Issue transition commits through the shared evaluator `(qa-lead)` | VAL-064 | positive guarded transition test asserting the evaluator identity used, with exactly one transition Event, terminated `exit_code: 0` |
| QA-63 | Confirm missing, unsuccessful, invalid, stale, and Artifact-missing evidence each block the terminal-success transition atomically `(qa-lead)` | VAL-065 | negative evidence matrix with exact reasons and a zero-Event delta per case, terminated `exit_code: 0` |
| QA-64 | Confirm Envelope completion reuses the shared evaluator and commits only with a terminal-success Issue plus all required Obligations satisfied `(qa-lead)` | VAL-066 | positive completion test plus an evaluator-identity assertion showing the same result object, with one `envelope.completed` Event, terminated `exit_code: 0` |
| QA-65 | Confirm failed completion returns every current reason and changes no Envelope, Issue, Evidence, or `task_completions` state `(qa-lead)` | VAL-067 | multi-failure atomicity test with before/after row snapshots including the legacy table, terminated `exit_code: 0` |
| QA-66 | Confirm resume is a bounded read from durable facts with exact next actions and no mutation, recovery, or respawn `(qa-lead)` | VAL-068 | fresh-process read test comparing the `Resume` projection against the reference and asserting a zero-Event delta, terminated `exit_code: 0` |

### 7.8 Window 7 - tasks 19-21

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-67 | Confirm export create is byte-deterministic and integrity-binds manifest, journal prefix, sources, capabilities, and available Artifacts `(qa-lead)` | VAL-069 | two archives created from identical facts with identical `sha256sum` output plus the exact entry inventory and per-entry digests, terminated `exit_code: 0` |
| QA-68 | Confirm export validate detects each safety and integrity defect without mutating project state `(qa-lead)` | VAL-070 | corruption corpus results with one `archive-invalid` check code per defect, plus identical `.fid` tree digests before and after validation, terminated `exit_code: 0` |
| QA-69 | Confirm export import is absent and `export.portable.v1` remains unavailable `(qa-lead)` | VAL-071 | the task-19 declared import-absence recipe (section 5.1) terminated `exit_code: 1`; the `fid export import` unknown-command result with exit 2; and the capability row `available:false` with `import:false` |
| QA-70 | Confirm migration 0014 performs the exact destructive drop, rename, and role-rule deletion sequence with the accepted empty-state down behavior and no data reconstruction `(qa-lead)` | VAL-002 | pre/post schema diff and row fixture showing `roles`, V1 `claims`, and `gate_runs` dropped, `work_claims` renamed to `claims`, displaced requirement rows deleted, and the frozen tables untouched, terminated `exit_code: 0` |
| QA-71 | Confirm every retained, fail-closed, and frozen V1 caller has its exact final non-authoritative disposition, using the classified inventory rather than a raw no-match rule `(qa-lead)` | VAL-074 | all four `v1-convergence-inventory.md` section 7 recipes run fresh at task-20 HEAD with the DEC-158 status split: the role-policy and `--by` recipes terminated `exit_code: 1`, and the broad table-reader and `EventType` recipes terminated `exit_code: 0`. Every line either recipe produced maps to an explicit retained, frozen, diagnostic, or removal disposition, the report records `UNCLASSIFIED=0`, and neither broad recipe was narrowed to suppress a frozen or diagnostic reader. Plus the per-disposition test naming each permitted reader |
| QA-72 | Confirm the final binary has no owned V1 policy, dev, completion, Claim, Run, or Gate command, writer, Event acceptance, role table, or mutation `--by` `(qa-lead)` | VAL-075 | the ten predeclared task-20 absence commands, each with its declared and observed non-zero status from section 5.1, plus the post-0014 schema listing showing the dropped tables absent |
| QA-73 | Confirm final-binary fresh `fid init` after 0014 succeeds and creates no `roles` table and no role Effect `(qa-lead)` | VAL-076 | `go build -o /tmp/fid ./cmd/fid` then `/tmp/fid init` in a fresh `mktemp -d` project; `.fid` schema listing with no `roles`; the genesis Event byte-compared to the bootstrap fixture, terminated `exit_code: 0` |
| QA-74 | Confirm deleting the fresh projection and rebuilding succeeds without `roles` and reproduces bootstrap and capability reads `(qa-lead)` | VAL-077 | delete-DB then `/tmp/fid rebuild` output, plus byte-identical `capabilities` and bootstrap snapshots before and after, terminated `exit_code: 0` |
| QA-75 | Confirm maintenance deterministically reports Event-sequence and aggregate/binding integrity defects `(qa-lead)` | VAL-078 | corrupt fixture pair producing exactly the `event-sequence`, `aggregate-sequence`, and `event-binding` findings in bytewise check order, terminated `exit_code: 0` |
| QA-76 | Confirm maintenance deterministically reports Attempt, Receipt, and Artifact integrity defects `(qa-lead)` | VAL-079 | corrupt fixture trio producing exactly `attempt-receipt-integrity` and `artifact-integrity` findings in order, terminated `exit_code: 0` |
| QA-77 | Confirm maintenance reports `frozen-state-present` and `removed-state-residue` without treating either as authority or repairing them `(qa-lead)` | VAL-080 | frozen and residue fixtures producing both findings with exact counts, plus identical projection digests before and after the check, terminated `exit_code: 0` |
| QA-78 | Confirm final-schema live and rebuilt capabilities, Envelope, Receipt/Evidence evaluation, completion state, and resume reads are byte-identical `(qa-lead)` | VAL-081 | before/delete/rebuild/after snapshot set with a byte diff of every named canonical read, terminated `exit_code: 0` |
| QA-79 | Confirm export bytes and validation result are identical before and after the final-schema rebuild `(qa-lead)` | VAL-082 | matching `sha256sum` for both archives and identical `ExportValidation` results, terminated `exit_code: 0` |
| QA-80 | Run the DEC-158 rule in this session: three consecutive fresh one-iteration 100k Event-v2 rebuild processes on the recorded reference instance, and capture the reference identity `(qa-lead)` | VAL-083 | three distinct captured `BenchmarkRebuildWorkspace100kEvents-<n>` result lines from three separate processes, each reporting `<30 s` and each terminated `exit_code: 0`, plus the parity snapshot per sample. All six identity commands captured and matching `02_TRD.md` section 9: `go1.26.2 linux/amd64`, `linux`/`amd64`, `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64`, `6`, `AMD Ryzen 7 5700G with Radeon Graphics`, `MemTotal: 8133236 kB`. Best-of and median sampling are forbidden; one sample at or above 30 s fails the assertion. Report lint cannot count samples, so QA-80 must confirm three distinct captured results rather than one repeated token. An unavailable reference instance is `BLOCKED`, not `PASS` |

### 7.9 Task 22 and final QA

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-81 | Confirm every accepted nonstreaming leaf matches its exact args, closed input, result, Event, replay, Human rendering, error, channel, and exit contract and rejects unknown fields `(qa-lead)` | VAL-072 | the committed conformance corpus with one case file per `public-interface.md` section 4.1 and 4.2 operation row; a row-by-row coverage table showing every operation present; per-case diffs against the contract; and the unknown-field rejection row per input schema. `ls cmd/fid/integration/testdata/p21_conformance` terminated `exit_code: 0` |
| QA-82 | Confirm no closed public enum, reason, Event, control, or availability value contains snake_case, while excluding JSON field names and opaque identifiers, and confirm the corpus is non-vacuous `(qa-lead)` | VAL-073 | the section 12 inventory pipeline over the committed `p21_conformance` corpus: (a) positive control - the extracted closed-value set contains every closed member enumerated in `public-interface.md` sections 3, 7, and 8, missing-member count 0, so an empty or truncated corpus cannot pass; (b) negative control - `grep -E $'\t[^\t]*[a-z0-9]_[a-z0-9]' "$SCRATCH/kept.tsv"` terminated `exit_code: 1`. `ls cmd/fid/integration/testdata/p21_conformance` terminated `exit_code: 0` |
| QA-83 | Confirm the final real-process one-Issue journey succeeds in the required Gate-before-terminal-transition order through completion and export `(qa-lead)` | VAL-084 | `go test -tags integration -race ./cmd/fid/integration/...` plus the ordered journey transcript showing, in order: nonterminal transition; Gate execution; Receipt link and evaluation; terminal-success transition; Run finish and Claim release; Envelope completion; resume; export create and validate. Real binary only, no mock or in-process substitute, terminated `exit_code: 0` |
| QA-84 | Confirm every final authority, policy, idempotency, evidence, Artifact, transition, and completion negative counterpart fails closed with exact codes and exits `(qa-lead)` | VAL-085 | tagged negative journey output listing each negative case with its `FailureCode`, category, process exit class, and the asserted absence of the forbidden effect, terminated `exit_code: 0` |
| QA-85 | Confirm that after delete and rebuild the journey repeats only the named reads, evaluation, completion state, resume, export, and negative checks, with no second Envelope creation and no other mutation `(qa-lead)` | VAL-086 | (a) positive control - the pre-rebuild Event count is exactly 20 and matches the `contracts.md` section 16 enumeration family by family; (b) post-rebuild Event count still exactly 20, so post-rebuild commands appended zero Events; (c) a command trace listing only the permitted post-rebuild operations; (d) `envelope.created` count of exactly 1 over the whole journey, terminated `exit_code: 0` |
| QA-86 | Audit that three-task reviews, the Group 1 checkpoint, the Phase Runner two-argument post-commit report lint before each next dispatch, final QA, and captured termination evidence occurred at the accepted cadence `(qa-lead)` | VAL-087 | phase artifact and run-ledger audit: the seven review artifacts plus the checkpoint review exist in order; every completion report carries an observed termination for every required command; no result is transcribed from another session; `git log` ordering shows each review after its window's last task and the checkpoint between tasks 12 and 13 |
| QA-87 | Execute the full seven-command phase gate from `AGENTS.md` in order in the final QA session `(qa-lead)` | no VAL | seven captured tokens with observed `exit_code` for each command in section 5's full phase gate row |
| QA-88 | PRD compliance review against every requirement cited in the validation contract's Source requirement column `(qa-lead)` | no VAL | a requirement-to-evidence table naming the `VAL-NNN` and the artifact that satisfies each cited PRD/TRD requirement, and explicitly listing the negative or unavailable requirements |
| QA-89 | Security review of executor, redaction, and path handling `(qa-lead)` | no VAL | the canary scan result, the closed-environment observation, the `credential-unavailable` pre-Attempt failure path, the symlink and traversal rejection matrix, the `unsafe-path` details never containing credentials or unrelated absolute paths, and the confirmation that credential values never appear in journal, projection, error arguments, or export |
| QA-90 | Performance spot-check `(qa-lead)` | no VAL | `go test -bench . -benchmem ./...` output with the rebuild, registry, and init benchmark lines and observed termination, including the retained per-sample `<30 s` assertion; deviation from the Phase 20 recorded values and from the three task-21 samples noted |
| QA-91 | Documentation review `(qa-lead)` | no VAL | confirmation that `AGENTS.md`, `project-docs/specs/02_TRD.md`, `03_TRD_SCHEMA.md`, `04_TRD_CONFIG.md`, `06_DESIGN_PLAN_CLI.md`, `07_DESIGN_PLAN_FLOWS.md`, and `GLOSSARY.md` match the shipped interface, with a diff or citation per changed surface |
| QA-92 | Task and completion-report audit `(qa-lead)` | no VAL | per task: spec commands versus report results, deviation notes for every failed-then-repaired attempt, dependency-budget notes, and a fresh two-argument `report-lint.sh` run per report with its token and observed `exit_code` |
| QA-93 | Final no-unclassified-inventory-match proof `(qa-lead)` | no VAL | all four `v1-convergence-inventory.md` section 7 recipes run fresh at phase HEAD with the DEC-158 status split: role-policy and `--by` terminate `exit_code: 1`; the broad table-reader and `EventType` recipes terminate `exit_code: 0` and every produced line maps to an explicit retained, frozen, diagnostic, or removal disposition, with `UNCLASSIFIED=0`. A raw no-match from a broad recipe is a failure, not a pass, because a frozen or diagnostic reader is expected to remain |

## 8. Mid-Phase Review Cadence

Seven QA Review windows plus one independent Group 1 checkpoint. Task-local
completion-report lint runs before every next dispatch.

| Window | After tasks | Reviewer | Artifact | Blocks |
|---|---|---|---|---|
| WINDOW-1 | 1, 2, 3 | QA Reviewer | `qa/review-1.md` | task 4 dispatch |
| WINDOW-2 | 4, 5, 6 | QA Reviewer | `qa/review-2.md` | task 7 dispatch |
| WINDOW-3 | 7, 8, 9 | QA Reviewer | `qa/review-3.md` | task 10 dispatch |
| WINDOW-4 | 10, 11, 12 | QA Reviewer | `qa/review-4.md` | Group 1 checkpoint |
| CHECKPOINT | 12 and WINDOW-4 | independent reviewer, not the author of any of tasks 1-12 and not the WINDOW-1..4 reviewer | `qa/group-1-checkpoint-review.md` | task 13 dispatch |
| WINDOW-5 | 13, 14, 15 | QA Reviewer | `qa/review-5.md` | task 16 dispatch |
| WINDOW-6 | 16, 17, 18 | QA Reviewer | `qa/review-6.md` | task 19 dispatch |
| WINDOW-7 | 19, 20, 21 | QA Reviewer | `qa/review-7.md` | task 22 dispatch |
| FINAL QA | 22 | QA Lead | `qa/final-review.md`, `qa/human-checklist.md` | phase wrap-up |

Rules:

- Review template: `project-docs/templates/mid-phase-review-template.md`.
- Each review verifies the in-window `VAL-NNN` rows from section 7, reruns at
  least one representative required command per reviewed task, and records
  gate evidence or why a gate was not run.
- Each review records exactly one outcome: `PASS`, `REPAIR`, or `BLOCK`.
- `REPAIR` names the owning task, the affected `VAL-NNN` or `no VAL`, the exact
  repair boundary, the required evidence, and the fresh independent re-review
  scope. The window does not advance. Findings always return to the owning
  task; they are never resolved in the next task.
- After repair, the Phase Runner dispatches a same-task Developer continuation
  and then a fresh independent QA Reviewer for the same window, writing the
  next unused review path and citing the prior review. The original reviewer
  never resolves its own finding.
- `BLOCK` is reserved for a protected authority boundary, out-of-window scope,
  or an exhausted no-progress recovery budget. QA records a blocker candidate
  and never allocates an `ESC` identity.
- Before every next dispatch the Phase Runner runs
  `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-<N>-done.md project-docs/docs/phases/phase-21/tasks/task-<N>-<slug>.md`
  against the committed report and records the token with its observed
  `exit_code`. A `RECOVERY_REQUIRED` or `BLOCKED` lint outcome blocks dispatch.
- Task 22 follows WINDOW-7. Final QA follows task 22.

Exact slugs for the report-lint second argument:
`task-1-schema-projection`, `task-2-event-idempotency`, `task-3-capabilities`,
`task-4-activation-trust`, `task-5-gate-registration`, `task-6-envelope-values`,
`task-7-authority-evaluator`, `task-8-issue-actor-convergence`,
`task-9-retained-actor-policy-removal`, `task-10-generic-claims`,
`task-11-generic-runs`, `task-12-execution-bound-transition`,
`task-13-executor-policy`, `task-14-live-output`, `task-15-attempt-receipt`,
`task-16-artifacts`, `task-17-evidence-evaluation`, `task-18-completion-resume`,
`task-19-export-create-validate`, `task-20-v1-removal`,
`task-21-maintenance-rebuild`, `task-22-walking-skeleton`.

## 9. Human Verification Scope

VAL-036 is the sole Human assertion in this phase. `09_AI_WORKFLOW.md` sections
2.6 and 6.1 require one human-checklist row per observable, so the assertion's
four named observables are fanned out into sub-IDs that all cite the same
parent VAL-036. No other row is delegated to the Human, and no row below is
pre-answered by QA.

Setup, performed once by the Human in an interactive terminal at the repository
root, after task 14 lands the committed fixture:

```text
cat cmd/fid/testdata/val-036/human-session.txt
```

The fixture is the byte-exact recorded Human-mode stdout of one `fid gate run`
over a deterministic Gate whose child process writes interleaved stdout and
stderr, raw ESC/CR/BEL bytes, a redacted credential, a rejected Artifact
candidate, and a non-success exit. `cat` is used deliberately so that any
unescaped control byte would reach the terminal. The exact fixture path is
requested from the Tech Lead as R-06.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-01a | Reading the recorded session, the Human can tell which lines came from the child's stdout and which came from its stderr without consulting any other artifact `(human)` | VAL-036 | written yes/no plus a note naming how stream identity was recognised |
| HUMAN-01b | Reading the same session, the Human understands each typed control notice (redaction, Artifact rejection, truncation) as a distinct, named event rather than as process output `(human)` | VAL-036 | written yes/no plus a note naming each control notice the Human identified |
| HUMAN-01c | Reading the same session, the Human understands the terminal result: whether a Receipt exists, the outcome, the exit or signal or timeout, Artifact IDs and rejections, and the Obligation status `(human)` | VAL-036 | written yes/no plus a note restating the terminal result in the Human's own words |
| HUMAN-01d | While the file is displayed with `cat` in a real terminal, no raw control byte changes the display: no colour change, cursor movement, bell, or line overwrite occurs `(human)` | VAL-036 | written yes/no plus a note of any observed display change |

VAL-036 passes only when all four rows are answered yes. Any no is a Must-fix
finding owned by task 14.

## 10. Group 1 Independent Checkpoint

Scheduled after task 12 and its WINDOW-4 review, before any task 13 dispatch.
The reviewer is assigned by the QA Lead and must be independent of tasks 1-12
and of the WINDOW-1..4 reviews. Output artifact:
`project-docs/docs/phases/phase-21/qa/group-1-checkpoint-review.md`.

The checkpoint has three separately decided parts. A partial, conditional, or
undocumented status never dispatches Group 2.

1. **VAL-042 (QA-40)** - Group 1 committed completion reports with Phase Runner
   two-argument post-commit report-lint evidence, and the four three-task QA
   review windows, are present and resolved. Decided on its own evidence.
2. **VAL-043 (QA-41)** - new Group 1 generic paths do not consult the exactly
   enumerated temporary V1 state, and the temporary inventory is classified as
   exactly the `v1-convergence-inventory.md` section 5 list. Decided on its own
   evidence. Final V1 absence is **not** required here; requiring it is itself
   a VAL-043 failure condition. At the checkpoint all four section 7 recipes
   terminate `exit_code: 0`, because removal happens in task 20; the DEC-158
   `exit_code: 1` expectation for the role-policy and `--by` recipes applies
   only at final task-20 HEAD.
3. **Fresh Group 2 feasibility (section 11)** - every FEAS probe in this
   session, with observed termination.

The classified temporary inventory the checkpoint must confirm is exactly:
tables `roles`, V1 `claims`, `gate_runs`, `task_completions`, `agent_runs`,
`ac_verifications`, `concerns`, and role/judgment rows in
`workflow_transition_requires`; commands and writers `fid dev start|gate|
finish|block`, `fid issue complete`, and their `dev.*` and completion Events;
source packages `internal/dev`, legacy completion files, and frozen-table
readers needed for removal conversion; and the internal V1 policy
implementation reachable only by the enumerated dev/completion callers with
`fid policy` already absent. Anything present but not on this list, and
anything on this list that a new generic path consults, is a checkpoint
failure.

Verdict rules follow `group-1-checkpoint-contract.md`: `PASS` requires every
precondition, review item, classification, isolation proof, and fresh
feasibility probe; `REPAIR` returns each finding to its owning Group 1 task and
repeats the checkpoint after recovery; `BLOCK` is reserved for a protected
contract change, an unavailable external capability after bounded recovery, or
a dependency/platform isolation need.

## 11. Fresh Group 2 Feasibility Probes

These probes establish that the repository and target test environment can
deterministically verify Group 2 behaviour. They must be executed in the
checkpoint session. Transcribed Phase 20 evidence, this QA plan's own planning
probes, a promised later fixture, an unbounded or hung command, and a green
result without observed termination are all not feasibility.

Setup, once per checkpoint session:

```text
SCRATCH="$(mktemp -d)"; echo "$SCRATCH"
cd "$SCRATCH"
printf 'module p21feas\n\ngo 1.26\n' > go.mod
mkdir -p child
```

Write `child/main.go` with a `PROBE_MODE` switch providing these modes:
`streams` (interleaved stdout and stderr writes), `badutf8` (raw `0xFF 0xFE`
plus `ESC [ 3 1 m` plus CR LF), `flood` (writes `PROBE_BYTES` to both pipes),
`sleep` (`time.Sleep(10 * time.Minute)`; never `select{}`, which produces a Go
deadlock `fatal error` and masks the signal result), `env` (prints
`os.Environ()`), and `secret` (echoes `$SECRET`). Build and run with:

```text
go vet ./...
go build -buildvcs=false -o "$SCRATCH/child-bin" ./child
```

`-buildvcs=false` is mandatory: without it, a scratch module outside a git
repository fails with `error obtaining VCS status: exit status 128`.

Declared bounded time budgets for the checkpoint session: each FEAS probe at
most 120s wall clock; the whole FEAS set at most 900s; the tagged real-process
journey gate `go test -tags integration -race ./cmd/fid/integration/...` at
most 300s (planning-base observation 4.658s).

| Probe | Bounded command or procedure | Pass condition and expected termination |
|---|---|---|
| FEAS-01 structured argv | `PROBE_MODE=streams "$SCRATCH/child-bin" arg-1 arg-2` with no shell interposed | child observes exactly `["arg-1","arg-2"]`; terminated `exit_code: 0` |
| FEAS-02 separate interleaved streams | `PROBE_MODE=streams "$SCRATCH/child-bin" > "$SCRATCH/out.txt" 2> "$SCRATCH/err.txt"` | `out.txt` holds only the stdout lines and `err.txt` only the stderr line; terminated `exit_code: 0` |
| FEAS-03 invalid UTF-8 and control bytes | `PROBE_MODE=badutf8 "$SCRATCH/child-bin" > "$SCRATCH/badutf8.bin"` then `xxd "$SCRATCH/badutf8.bin"` | hexdump shows `41ff fe1b 5b33 316d 0d0a`; terminated `exit_code: 0` |
| FEAS-04 timeout | `PROBE_MODE=sleep timeout -s TERM 5 "$SCRATCH/child-bin"` | terminated `exit_code: 124` |
| FEAS-05 signal | `PROBE_MODE=sleep timeout -s KILL 5 "$SCRATCH/child-bin"` | terminated `exit_code: 124` and the process is observably killed |
| FEAS-06 spawn failure | `"$SCRATCH/does-not-exist"` | terminated `exit_code: 127` |
| FEAS-07 high-output draining | `PROBE_MODE=flood PROBE_BYTES=8388608 "$SCRATCH/child-bin" > "$SCRATCH/flood.out" 2> "$SCRATCH/flood.err"` then `wc -c "$SCRATCH/flood.out" "$SCRATCH/flood.err"` | reports 8388608 bytes on each stream with no hang inside the 120s budget; terminated `exit_code: 0` |
| FEAS-08 closed environment | `env -i PROBE_MODE=env PATH=/usr/bin HOME=/tmp "$SCRATCH/child-bin"` | output is exactly those three variables; terminated `exit_code: 0` |
| FEAS-09 credential and literal canary | `SECRET=CANARY-<fresh-value> PROBE_MODE=secret "$SCRATCH/child-bin" > "$SCRATCH/canary.txt"`, then `rg -n 'CANARY-<fresh-value>' "$SCRATCH/canary.txt"` | the canary is observable on a real surface (`exit_code: 0`), proving a later absence scan is meaningful rather than vacuous |
| FEAS-10 cwd and Artifact containment | create `cwd/`, `outside/target.txt`, and symlink `cwd/escape.txt -> ../outside/target.txt`; resolve with `realpath` and compare against the resolved `cwd` base | the resolved path is reported outside the base; terminated `exit_code: 0` |
| FEAS-11 changed-during-copy | hash a file, copy it, mutate it, hash again | the two digests differ; terminated `exit_code: 0` |
| FEAS-12 atomic content storage | write to `store/tmp/part`, digest it, `mv` to `store/objects/<digest>` on the same filesystem | `store/tmp` is empty and `store/objects/<digest>` exists; terminated `exit_code: 0` |
| FEAS-13 observed termination | for FEAS-04, FEAS-05, FEAS-06, and FEAS-07, record the nested result object's real exit status rather than wrapper completion | each recorded status is copied from the command's own result; a claim supported only by wrapper completion or empty output fails the probe |
| FEAS-14 repository-tree freshness | `git init` a scratch repository, commit, record `git rev-parse HEAD^{tree}`, change a file, commit, record again | the two tree OIDs differ; terminated `exit_code: 0` |
| FEAS-15 Artifact loss and restoration | back up `store/objects/<digest>`, delete it, confirm absence, restore it, confirm presence | absence then presence observed; terminated `exit_code: 0` |
| FEAS-16 deterministic ZIP | build the same two-entry archive twice at least one second apart with fixed 1980-01-01 headers and fixed modes and compression, then `sha256sum` both | identical digests; terminated `exit_code: 0` |
| FEAS-17 non-mutating corruption validation | `sha256sum ar.zip`; `unzip -t ar.zip`; `sha256sum ar.zip`; then corrupt one byte with `dd ... conv=notrunc` and rerun `unzip -t` | first `unzip -t` terminated `exit_code: 0` with an unchanged digest; the corrupted run terminated `exit_code: 1` |
| FEAS-18 expected-nonzero absence commands | run every section 5.1 expected-non-zero command shape once against a path set known to contain no match | each terminates with its declared non-zero status and no output; a converted `exit 0 (no output)` fails the probe |
| FEAS-19 tagged integration gate budget | `go test -tags integration -race ./cmd/fid/integration/...` | `ok github.com/mbell-nzl/fid/cmd/fid/integration <t>s` within the 300s budget; terminated `exit_code: 0`. Per the amended checkpoint contract this proves environment feasibility only and does not require the task-22 product journey to exist after task 12 |

## 12. VAL-073 Closed-Value Inventory Procedure

Used by QA-82 at final QA over the committed conformance corpus. Materialize
the exclusion list into scratch; do not add it to the repository.

```text
SCRATCH="$(mktemp -d)"
CORPUS=cmd/fid/integration/testdata/p21_conformance
cat > "$SCRATCH/exclude.txt" <<'PATTERNS'
(^|\.)id$
_id$
_ids\.\[\]$
(^|\.)ids\.\[\]$
_digest$
(^|\.)digest$
_hash$
_oid$
(^|\.)path$
_path$
(^|\.)paths\.\[\]$
(^|\.)entry$
(^|\.)journal_file$
(^|\.)root$
(^|\.)source_doc$
(^|\.)source_section$
(^|\.)message$
(^|\.)detail$
(^|\.)text$
(^|\.)title$
(^|\.)summary$
(^|\.)display_label$
(^|\.)external_id$
(^|\.)rule$
_at$
(^|\.)capability_id$
(^|\.)capabilities\.\[\]$
(^|\.)required_capabilities\.\[\]$
(^|\.)actions\.\[\]$
(^|\.)labels\.\[\]$
(^|\.)label$
(^|\.)name$
(^|\.)value$
(^|\.)argv\.\[\]$
(^|\.)cwd$
(^|\.)credential_names\.\[\]$
(^|\.)redact_literals\.\[\]$
(^|\.)artifact_paths\.\[\]$
(^|\.)names\.\[\]$
(^|\.)fields\.\[\]$
(^|\.)values\.\[\]$
(^|\.)allowed\.\[\]$
(^|\.)output_excerpt$
(^|\.)toolchain\.executable$
(^|\.)provenance_ref$
(^|\.)session_ref$
(^|\.)authorization_provenance$
(^|\.)kit_identity$
(^|\.)backend_identity$
(^|\.)idempotency_key$
(^|\.)operation_scope$
(^|\.)actor_attribution$
(^|\.)raised_by_attribution$
(^|\.)addressed_by_attribution$
(^|\.)model$
(^|\.)non_goals\.\[\]$
(^|\.)envelope\.outcome$
(^|\.)from_state$
(^|\.)to_state$
(^|\.)issue\.status$
(^|\.)issues\.\[\]\.status$
(^|\.)task\.status$
(^|\.)blocking_task_status$
(^|\.)data\.text$
PATTERNS
```

Extraction, filtering, and the two controls:

```text
find "$CORPUS" \( -name '*.json' -o -name '*.jsonl' \) | sort | while read -r f; do
  jq -c . "$f"
done | jq -r 'paths(strings) as $p
  | [($p|map(if type=="number" then "[]" else tostring end)|join(".")), getpath($p)]
  | @tsv' \
  | grep -Ev -f <(sed 's|^|^[^\t]*|' "$SCRATCH/exclude.txt") \
  > "$SCRATCH/kept.tsv"

wc -l "$SCRATCH/kept.tsv"
cut -f2 "$SCRATCH/kept.tsv" | sort -u > "$SCRATCH/observed-values.txt"
grep -E $'\t[^\t]*[a-z0-9]_[a-z0-9]' "$SCRATCH/kept.tsv"
```

- Positive control: every closed enum member enumerated in
  `public-interface.md` sections 3, 7, and 8 must appear in
  `"$SCRATCH/observed-values.txt"`. The missing-member count must be 0. This
  prevents a truncated corpus from passing vacuously.
- Negative control: the final `grep` must terminate `exit_code: 1` with no
  output. `exit 0 (no output)` is not a valid record here.
- This pipeline was validated in the QA Lead session against a seeded sample
  containing `spawn_failed`, `in_progress`, `needs_review`, `head_oid`, and
  `tree_oid changed`; it reported only the seeded closed-value violation.

## 13. Refinement Disposition

The Tech Lead applied DEC-158 at `b367a3b`. Every first-review proposal is
dispositioned below and independently verified in this session. QA edited no
contract, checkpoint contract, or task spec.

| ID | Disposition | Independent verification at `b367a3b` |
|---|---|---|
| R-01 | Applied | Report lint removed from all 22 `## Exact Commands To Run` blocks; `grep -l "report-lint.sh" .../tasks/*.md` terminated `exit_code: 1`. Each spec now states it is a Phase Runner post-commit control. Recorded consistently in `phase_brief.md` Dispatch Controls, `08_QUALITY_PLAN.md` section 8, the checkpoint preconditions, and VAL-087 |
| R-02 | Applied, threshold retained | DEC-158 keeps the 30s budget and the per-sample `b.Fatalf`, names the reference instance in `02_TRD.md` section 9, and sets the three-consecutive-fresh-process, all-samples-`<30 s` rule with best-of and median forbidden. QA re-measured the exact procedure clean: 23.728s, 23.797s, 23.800s, all `exit_code: 0` |
| R-03 | Applied | Task 21 adds the three benchmark invocations plus all six identity commands; every identity literal matches the observed machine exactly |
| R-04 | Applied | Task 22 declares `cmd/fid/integration/testdata/p21_conformance/` with `manifest.json`, per-operation, per-failure-code, and per-live-variant fixtures, regenerated and byte-compared by the tagged suite, and adds `ls` on it as a required exit-0 command |
| R-05 | Applied and widened | Task 13 adds the containment lexicon over `cmd internal` **and** a second lexicon over `contracts.md` and `public-interface.md`, both declared `exit_code: 1`. Both terminated `exit_code: 1` when run verbatim |
| R-06 | Applied | Task 14 declares both VAL-036 fixtures and an acceptance criterion requiring golden regeneration and byte comparison, interleaved streams, all three control kinds, a non-success terminal frame, and raw ESC/CR/BEL child output |
| R-07 | Applied | VAL-036 is restated as one parent perception assertion; the Coverage Audit now requires the `HUMAN-01a..01d` fan-out under Workflow section 6.1 and keeps VAL-036 the sole Human parent. Section 9 matches |
| R-08 | Applied | The checkpoint contract now reads "environment feasibility to run the existing tagged integration gate within 300 seconds and each FEAS probe within 120 seconds. This does not require the task-22 product journey to exist after task 12" |
| R-09 | Applied with an accepted Tech Lead correction | QA had assumed all four section 7 recipes terminate 0 at task 20. DEC-158 correctly separates them: the role-policy and `--by` recipes are already predeclared absence controls that must terminate `exit_code: 1` at final task-20 HEAD, while only the broad table-reader and `EventType` recipes terminate `exit_code: 0` as classified inventories. QA accepts this and has corrected sections 5.1, 7.5, 7.8, and QA-93. Verified: the two broad recipes were added to task 20, both terminate `exit_code: 0` today (54 and 22 matching lines), neither is declared expected-non-zero, the expected-non-zero section still holds exactly the 10 original declarations, and the inventory file plus the task's Notes forbid narrowing them |
| R-10 | Applied | Task 4 adds the declared recovery-absence recipe; terminated `exit_code: 1` |
| R-11 | Applied | Task 19 adds the declared import-absence recipe; terminated `exit_code: 1` |
| R-12 | Applied | Task 17 adds the declared Attention-absence recipe; terminated `exit_code: 1` |
| R-13 | Applied | Tasks 5 and 19 note that they create `internal/evidence/` and `internal/export/` and that the focused package command is valid only from that task onward |
| R-14 | Applied | VAL-044 and VAL-045 are explicitly cross-paired; an unpaired positive-only or negative-only claim now fails by contract |
| R-15 | Applied | `phase_brief.md` gains a Dispatch Controls section naming the two-argument control and all nine QA artifact paths |
| R-16 | Applied | VAL-073 now requires every closed member from interface sections 3, 7, and 8 in the extracted set with missing count 0, so an empty or truncated corpus fails |
| R-17 | Applied and made exact | VAL-086 and `contracts.md` section 16 now fix the pre-rebuild count at exactly 20 Events with exactly one `envelope.created`. QA independently enumerated the journey and derived the same 20 families before reading the amendment |
| R-18 | Applied | VAL-010 now names the seven false capabilities and the exact surfaces they govern; task 3's acceptance criterion matches |
| R-19 | Deferred by decision | The contract keeps `Status: Draft` on all 87 rows until Human execution approval. QA accepts the deferral and notes it is the Tech Lead's to flip at approval |

### 13.1 Named execution risks carried into the phase

These are risks for the Phase Runner and reviewers to watch. None is a QA-plan
defect and none blocks readiness.

| Risk | Detail | Mitigation already in the plan |
|---|---|---|
| RISK-01 | The 100k rebuild has about 6.2s of headroom at planning base, and Phase 21 adds Event-v2 effects, idempotency inserts, and richer projection rows. Task 21 could legitimately breach the budget. | DEC-158 assigns causal optimization to task 21 scope and forbids threshold relaxation, best-of, and median. QA-80 requires three distinct captured samples. A sustained breach after bounded optimization is a blocker candidate for the Tech Lead, not a QA threshold change |
| RISK-02 | Report lint cannot count task 21's three identical benchmark invocations; one recorded PASS token satisfies its coverage check three times. | QA-80 must confirm three distinct captured result lines from three separate processes, not one repeated token |
| RISK-03 | The measurement in RISK-01 is load-sensitive: a concurrent full race suite pushed one planning-base sample from about 23.8s to 33.5s. | The task-21 samples must be taken as dedicated fresh processes; QA-80 records this and rejects samples taken under concurrent load |
| RISK-04 | The task-20 broad classified inventories can be silently weakened by narrowing their patterns, which would convert a classification duty into a vacuous no-match. | `v1-convergence-inventory.md` section 7, the task-20 Notes, QA-71, and QA-93 all state that a raw no-match from a broad recipe is a failure |

## 14. Phase-End Decision Rule

The phase proceeds to final review only when all of the following hold.

- The full seven-command `AGENTS.md` phase gate passed in the final QA session
  with an observed termination for every command (QA-87).
- Every `VAL-001..087` has a recorded outcome in
  `## Validation Contract Outcomes` of `final-review.md`: `PASS`, `FAIL`,
  `BLOCKED`, or `CARRY-FORWARD`, each with cited evidence. `PASS` without
  evidence is invalid.
- Every `(qa-lead)` item QA-01..QA-93 has reproducible evidence captured in the
  session that records it. Code reading alone never satisfies an item.
- HUMAN-01a..01d are answered by the Human and none is pre-answered,
  delegated, or inferred by QA.
- All seven review windows are `PASS`; every earlier `REPAIR` has a fresh
  independent `PASS` re-review resolving its original Must-fix findings.
- Every task-local report lint ran as the Phase Runner's two-argument
  post-commit control before the next dispatch, with its token and observed
  `exit_code` in the run ledger.
- The Group 1 checkpoint recorded `PASS` separately for VAL-042 and VAL-043 and
  for every FEAS probe, and did so before task 13 dispatch.
- VAL-083 has three distinct captured samples from three separate fresh
  processes on the recorded reference instance, each under 30 s, with the six
  identity commands captured and matching `02_TRD.md` section 9.
- No unresolved Must-fix finding or blocker candidate remains. BC-01 and BC-02
  are closed by DEC-158 and re-verified in this plan's section 6.
- The final classified-inventory proof (QA-93) reports 0 unclassified matches
  with the DEC-158 status split.
- Required documentation updates named in QA-91 are applied.

An in-scope failed gate or incomplete evidence triggers repair; recording
`BLOCKED` does not make it acceptable. A protected missing or contradictory
contract prevents final review and is emitted as a blocker candidate for
TL-Lite or the full Tech Lead, which alone may allocate an `ESC` identity.

## 15. Plan Status

| Item | State |
|---|---|
| VAL mapping completeness | 87 of 87 mapped exactly once (86 `QA-NN` rows plus HUMAN-01a..01d for VAL-036); no ID merged; contract validator roles preserved |
| Review windows | 7 defined, plus the independent Group 1 checkpoint |
| Checkpoint assertions | VAL-042 and VAL-043 separately decided |
| Fresh feasibility | FEAS-01..FEAS-19 defined and mechanically validated; execution required fresh at the checkpoint |
| Final QA | QA-86..QA-93 cover the full gate, PRD compliance, security, performance, documentation, VAL outcomes, task/report audit, and the classified-inventory proof |
| Human scope | VAL-036 only, as one parent assertion fanned out to HUMAN-01a..01d |
| Command validation | every command, including all nine DEC-158 additions, executed verbatim in a QA Lead session with an observed exit code |
| Blocker candidates | BC-01 and BC-02 closed by DEC-158 and re-verified |
| Outstanding Must-fix findings | none |
| Readiness | **Ready** |
| Task 1 | Remains unauthorized until Human execution approval |
