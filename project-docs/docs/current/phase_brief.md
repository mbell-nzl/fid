# Phase Brief - Phase 20

> Default execution brief loaded by every agent session.

## Phase Goal

Deliver one coherent standalone Issue journey on canonical JSON: initialize a
project; create, relate, inspect, update, and transition Issues; explain failure
with typed reasons; rebuild the projection; and repeat unchanged. Establish the
`fid-cli.v1` Response Envelope, project resolution/identity, replayable Issue
revision, and the retained Issue contract. Delete universal formatters,
`--format`, promised-command stubs, and Issue fields with no retained meaning.

## Status

**Human-approved implementation pack; QA planning is next.** Branch `main`; no
product code changed. Approval covers heavy mode, the accepted slice/validation
contracts, and the eighteen-task pack. No Developer starts before QA-plan acceptance.

## Operating Mode

**Heavy** (DEC-107/DEC-118/DEC-127): schema/persistence, public-interface, and dependency
changes. Binding delivery process remains the V1 workflow in
`project-docs/specs/09_AI_WORKFLOW.md`. Required: Phase Runner, fresh QA plan and
feasibility validation, six mid-phase reviews, final review, Human checklist,
full spec sync, run ledger, and triage log.

## Accepted Slice Contract

`../phases/phase-20/contracts.md` defines the accepted global CLI behavior,
Response Envelope, exact leaf result schemas, reason codes, physical DDL/Event
semantics, removal inventory, and exact legacy-command expiry ledger.

## Task Index

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Response contract foundation | Pending | None |
| 2 | Project selection and config namespace | Pending | 1 |
| 3 | Atomic project initialization | Pending | 1-2 |
| 4 | Rebuild, authority, and lock failures | Pending | 1-3 |
| 5 | Remove orphaned Issue fields | Pending | 4 |
| 6 | Replayable Issue aggregate revision | Pending | 5 |
| 7 | Issue hierarchy invariant | Pending | 5-6 |
| 8 | Issue read contract | Pending | 5-7 |
| 9 | Issue mutation and concurrency contract | Pending | 6-8 |
| 10 | Acceptance-criterion and label contracts | Pending | 8-9 |
| 11 | Link contract and dependency invariants | Pending | 8-10 |
| 12 | Workflow command contract | Pending | 8-11 |
| 13 | Maintenance inspection contract | Pending | 8-12 |
| 14 | Legacy JSON renderer conversion | Pending | 13 |
| 15 | Remaining format caller/test convergence | Pending | 14 |
| 16 | Delete universal formatters and YAML | Pending | 15 |
| 17 | Delete promised-command stubs | Pending | 16 |
| 18 | Journey acceptance and documentation | Pending | 1-17 |

One Developer session owns one task. Review windows close after Tasks 3, 6, 9,
12, 15, and 17. Task 18 follows the sixth review and precedes phase-end QA.

## Validation Contract

Full draft: `../phases/phase-20/qa/validation-contract.md`.

| VAL IDs | Owning work |
|---|---|
| 001-008 | transport, project selection, authority, locks, channels (Tasks 1-4; closes Task 18) |
| 009-012 | field removal, aggregate revision, replay/rebuild (Tasks 5-6, 9-11) |
| 013-020 | Issue reads, mutations, transitions, hierarchy (Tasks 7-9) |
| 021-026 | acceptance criteria, labels, links, dependencies, workflows (Tasks 10-12) |
| 027-028 | maintenance contract and legacy journal boundary (Tasks 13-14) |
| 029-034 | legacy conversion and deletion (Tasks 14-17) |
| 035-037 | inventory, journey, documentation (Task 18) |
| 038 | Human phase-end judgment |
| 039-043 | config identity/strictness, owner boundary, attribution, atomic init (Tasks 2-12) |
| 044-046 | exact requests, order/bounds, opaque typed link endpoints (Tasks 2-18) |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| Validation contract | `../phases/phase-20/qa/validation-contract.md` | Human-approved; rows remain Draft pending QA-plan acceptance |
| QA plan | `../phases/phase-20/qa/qa-plan.md` | Missing; QA Lead creates before Developer execution |
| Mid-phase reviews | `../phases/phase-20/qa/review-M.md` | After Tasks 3, 6, 9, 12, 15, 17; review 6 includes independent contract/spec consistency |
| Final review / Human checklist | `../phases/phase-20/qa/` | Pending |

## Active Constraints

- Implement the approved form of `contracts.md` literally; divergence escalates.
- Preserve fail-closed V1 authority until Phase 21 replacement passes. Phase 20
  maps errors; it does not weaken checks.
- Pull no Phase 21 capability forward: no capabilities command, Work Envelope,
  generic authority, Claim/Run reshape, gates/Receipts/Evidence, idempotency, or
  backend-port conformance.
- The caller never supplies resulting Issue revision. Mutation Events carry
  absolute result plus prior/precondition values for deterministic replay.
- Issue revision is aggregate-wide: AC, Issue-label, and link mutations advance
  every affected Issue as well as create/update/transition.
- Root `--project`/`--json` never silently apply to legacy leaves; explicit use
  rejects before their handler.
- Workflow project-wide assignment uses `--project-wide`; global `--project`
  always selects a fid project root (DEC-119).
- Legacy JSON carries tested meaning only, not byte compatibility or a forward contract.
- `journal inspect` remains legacy until Phase 25; Phase 20 does not publish V1
  role/revert/Event internals through `fid-cli.v1`.
- Bare “Envelope” means Work Envelope; JSON transport is Response Envelope.
- `09_AI_WORKFLOW.md` remains the binding V1 process used to build fid. Product
  target concepts do not replace its roles or artifacts.

## Parity Matrix

| Axis | Retained JSON leaves | Retained Human mode | Legacy leaves | Evidence |
|---|---|---|---|---|
| Project selection | global flag/discovery | same root | explicit root flags reject | VAL-003..005 |
| Success/failure | closed envelope/channel | concise stdout/stderr | unenveloped established mode | VAL-001/002/008/029 |
| Request surface | closed args/flags | same operation args | no forward contract | VAL-044 |
| Collections | exact order/bounds | same semantic order | tested meaning only | VAL-045 |

## Interfaces / Subsystems In Scope

- `cmd/fid/`, `internal/cli/`, `internal/project/` - root dispatch and canonical CLI transport.
- `internal/config/`, `internal/bootstrap/` - strict config and atomic init.
- Issue/label/link/workflow/maintenance domain packages - standalone journey behavior.
- `internal/journal/`, `internal/db/`, `migrations/` - aggregate revision and rebuild equivalence.
- `internal/format/`, module files, command tests - sharp formatter removal.

## Known Risks

- Task 2 is the widest unit. Its QA sequence must stage discovery, strict
  parser/writer, transitional bootstrap, maintenance decoupling, then root
  dispatch; review 1 checks for one config model and one root-flag registry.
- Schema/fixture churn is isolated in Tasks 5-6 and reviewed after Task 6.
- Command migration crosses several adapters but is split into three-task review windows.
- Formatter removal remains high-churn; legacy conversion, caller cleanup, and
  deletion are split across Tasks 14-16 and reviewed after Tasks 15 and 17.
- Exact command inventories in VAL-035 include `fid issue ac verify`; family ellipses are prohibited.

## Open Questions

- None before implementation. A newly discovered contract/data-model ambiguity is a
  hard stop and returns to the Human rather than being inferred by a Developer.

## Success / Exit Criteria

- [ ] All eighteen tasks complete with required exact-command evidence.
- [ ] Every accepted VAL has a phase-end outcome with cited evidence.
- [ ] Complete standalone Issue journey passes and survives rebuild unchanged.
- [ ] Universal formatters/flags/dependency, promised stubs, and retired Issue fields are absent.
- [ ] Fail-closed V1 authority still rejects; no Phase 21 capability is advertised.
- [ ] Final QA and Human checklist have no unresolved blocking outcome.

## Next Planning Checkpoint

After Phase 20 wrap-up, resolve Gate D3.2: the exact one-Issue Routine walking
skeleton boundary for Phase 21.
