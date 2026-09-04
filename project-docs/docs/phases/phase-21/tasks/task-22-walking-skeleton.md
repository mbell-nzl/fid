# Phase 21 - Task 22: Final Wire Conformance And Real-Process Journey

## Objective

Prove every final accepted/fail-closed public leaf, all kebab-case closed values, the Gate-before-terminal-transition real-process journey, its negative counterparts, and the mutation-free post-rebuild repeat.

## Validation Contract Coverage

- VAL-072: owning observable in the Phase 21 validation contract.
- VAL-073: owning observable in the Phase 21 validation contract.
- VAL-084: owning observable in the Phase 21 validation contract.
- VAL-085: owning observable in the Phase 21 validation contract.
- VAL-086: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Build and invoke the final binary; do not mock or use in-process substitutes for the journey.
- Mutate once, rebuild durable facts, then repeat only specified reads/evaluation/completion state/resume/export/negative checks.
- Do not read held-out scenarios or create a second Envelope.
- Author the journey workflow with `workflow define` then `workflow assign --kind task` before `issue create`, exactly as `contracts.md` §16 specifies (DEC-174). Do not modify, extend, or reinterpret the migration-seeded `task-default`, `subtask-default`, or `epic-default` reference workflows, and do not add a product transition, category, or workflow behaviour to reach terminal success.

## Scope

### Files To Create Or Modify

- `cmd/fid/integration/`
- `cmd/fid/integration/testdata/p21_conformance/`
- `internal/cli/`
- `internal/workflow/activation.go` - **bounded (DEC-175):** the Event-type
  derivation in `transitionActivation` only. Replace `operation + "ted"` and its
  two ad-hoc overrides with one explicit, exhaustive operation-to-Event-type
  mapping: `activation.commit` -> `activation.committed`, `activation.abort` ->
  `activation.aborted`, `activation.finalize` -> `activation.finalized`,
  returning a typed error for any unmapped operation. The commit and abort
  values stay byte-identical. No other change to this file or any other
  production file is authorized.
- `cmd/fid/activation_test.go` - **bounded (DEC-175):** add the regression
  assertion in the acceptance criteria only; do not modify or weaken any
  existing assertion in this file.
- `internal/evidence/execution.go`, `internal/evidence/receipt_read.go`, and
  focused `internal/evidence/*_test.go` - **bounded (DEC-177):** persist the
  accepted Receipt rejection array and implement only the accepted read-only
  get/list contract.
- `cmd/fid/root.go`, `cmd/fid/receipt.go`, and focused `cmd/fid/*_test.go` -
  **bounded (DEC-177):** register and adapt only `fid receipt get/list` with the
  exact accepted args, flags, errors, and JSON results.
- `migrations/0013_p21_operational_kernel.up.sql` and
  `internal/db/operational_kernel_migration_test.go` - **bounded (DEC-177):**
  add and assert only the required canonical rejected-Artifact JSON array.
- `project-docs/specs/03_TRD_SCHEMA.md` and
  `project-docs/docs/phases/phase-21/contracts.md` - **bounded (DEC-177):**
  record the accepted Receipt persistence detail and correct only the stale
  Resume command spelling.
- `project-docs/docs/phases/phase-21/tasks/task-22-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/docs/phases/phase-21/contracts.md §16`
- `project-docs/docs/phases/phase-21/gate-d3-2-proposal.md`
- `project-docs/docs/phases/phase-21/qa/validation-contract.md`

### Out Of Scope

- Product repairs outside the DEC-175/177 accepted-contract repairs, import,
  recovery, Review/Human Outcome, and held-out scenarios.
- Any `internal/workflow/activation.go` change beyond the DEC-175 Event-type
  derivation: no lifecycle, precondition, Effect, persistence, result-shape,
  idempotency, or attribution change.

## Acceptance Criteria

- [ ] The committed `p21_conformance/` corpus is regenerated and byte-compared by the tagged suite: `manifest.json`; one `operations/<operation>.json` fixture for every accepted and fail-closed operation row in `public-interface.md` §§4.1-4.2; one `failures/<failure-code>.json` fixture for every §7 failure; and one `live/<variant>.jsonl` fixture for every §8 stdout, stderr, control-kind, terminal-outcome, and replay variant. (VAL-072)
- [ ] Corpus extraction proves non-vacuity by requiring every closed member from `public-interface.md` §§3, 7, and 8 in the extracted set, and every extracted public enum value is lowercase kebab-case. (VAL-073)
- [ ] Happy journey defines and assigns the §16 journey workflow, creates its Issue on that workflow, and follows Gate/evidence before the single terminal-success transition through completion. (VAL-084)
- [ ] All negative counterparts fail closed, including a transition request for an edge the authored workflow does not declare, which returns `transition-not-permitted` / `workflow-transition-invalid` with exit 4 and appends no Event. (VAL-085)
- [ ] Before rebuild, the journey has exactly 22 Events and exactly one each of `envelope.created`, `workflow.defined`, and `workflow.assigned`; after rebuild, the repeat performs no second Envelope, workflow definition, assignment, or other mutation and the Event count remains exactly 22. (VAL-086)
- [ ] Each activation transition emits its `public-interface.md` §4.1 Event type verbatim - `activation.committed`, `activation.aborted`, and `activation.finalized` - asserted by a non-tagged regression test in `cmd/fid/activation_test.go` that runs under `go test -race -coverprofile=coverage.out ./...`. (DEC-175; supports VAL-084 and VAL-086)
- [ ] Receipt get/list return exact complete Receipts, accepted filters and
  ordering, typed not-found/validation failures, and no Event; rejected
  Artifact metadata survives idempotent replay and projection rebuild in
  canonical path/reason order. The real-process journey reads its Receipt both
  before and after rebuild without changing the exact 22-Event count.
  (DEC-177; supports VAL-052, VAL-072, VAL-084, and VAL-086)

## First Failing Test

Add the failing committed-corpus regeneration/byte-comparison and complete-member non-vacuity tests, then the forbidden-enum inventory, ordered binary journey, negative matrix, exact 22-Event/one-`envelope.created` pre-rebuild controls, and post-rebuild zero-delta assertion. The committed corpus and the already-committed journey guard are retained work; the recovery continues from them.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-22-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
ls cmd/fid/integration/testdata/p21_conformance
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after task 21 and QA Review window 7 close. Final QA follows; this task does not authorize its own contract refinements.
- The authored workflow needs a `cancelled` state only because an Issue-applicable workflow must declare one (PRD-WF-013). Every non-`cancelled` state has an implicit edge into it, so the VAL-085 undeclared-edge negative must target a non-`cancelled` state, and the journey never enters `cancelled`.
- DEC-174 recovery: continue from `236540c` under the amended §16 journey. Preserve the attempt-1 blocker evidence, including the recorded `transition-not-permitted` failure, and reuse it as the VAL-085 undeclared-edge negative counterpart rather than deleting it. All eight required commands must run and terminate with session-captured evidence; the eight `not run` rows are not carried forward.
- DEC-175 recovery (attempt 4 under DEC-176's append-only identity reconciliation): continue from `a6b3a44`. Keep the committed corpus, journey, and negative matrix. The tagged gate's `journal Event 9 = "activation.finalizeted", want "activation.finalized"` is a product defect, not a contract question: `internal/workflow/activation.go:519` builds the Event type as `operation + "ted"` and overrides only `activation.commit` and `activation.abort`, so `activation.finalize` emits the misspelling. `public-interface.md` §4.1, `contracts.md` §14, and `contracts.md` §16 all require `activation.finalized`. Write the new `cmd/fid/activation_test.go` assertion red first, then make the bounded derivation repair. Nothing else changes: no Event schema, public interface, contract, command, expected exit, assertion strength, corpus, or seeded workflow edit, and no `t.Skip`, build tag, or `-run` narrowing. Preserve the attempt-1 and attempt-2 failed evidence verbatim in the report and add this attempt's red/green evidence. Attempt 3 is only the immutable TRG-21-022 orchestration identity; it is not a product attempt and consumes no no-progress budget. Every other Event family in the §16 sequence was verified present in a product writer, so no further name defect is expected; stop and report if one appears.
- DEC-177 recovery (attempt 6): continue from the committed activation and
  root-Resume corrections. Preserve both attempt-4 tagged failures. Add focused
  red assertions for exact Receipt get/list, typed failures, filters/order,
  read-only behavior, and durable sorted rejection metadata before implementing
  the bounded schema/domain/adapter repair. Then rerun the complete
  byte-identical eight-command block.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-22-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
