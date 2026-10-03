# Phase 24 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow
> notes, and wrap-up decisions. DEC numbering is global and never resets.
> Phase scope authority is DEC-272 (`../phase-23/decisions.md`). Entries cite
> the contract (`public-interface.md`) rather than restating it.

## Decisions

### DEC-273: Run and gate-run repository checks use the real worktree (D1 = A)

- **Date:** 2026-10-04. ACCEPTED by the Human ("yes, option A") after the
  Tech Lead explained the alternative.
- **Context:** today `run.start` checks the caller's repository block only
  against the Envelope's stored `repository_scope`, and `gate.run` only
  against the Run's `start_repository`. Neither is ever compared with the
  real worktree. Once fid observes the repository (DEC-272 Part 1), each
  comparison must mean something. Option B (real worktree equals the
  snapshot) would refuse every rework Run after an edit, which is impossible
  to escape in Phase 24 because Envelope revision is Phase 25 scope. It would
  also stop a Builder testing their own edits within a Run.
- **Decision:** `run.start` and `gate.run` compare the Envelope or Run binding
  as *same project and same worktree* only. The caller's
  `expected_repository` (`head_oid`, `dirty`) is checked against fid's
  observation, and a mismatch rejects with `repository-mismatch`. The
  observation is recorded as the Run's `start_repository`, and for `gate.run`
  as the Attempt and Receipt repository.
- **Consequences:** a Run or gate run may start after edits or commits made
  since the Envelope or Run began. The existing freshness rule (PRD-EVID-003)
  stales proof whose recorded tree differs from the current one, which
  becomes a real check under task 8. `run.recover` is unaffected because it
  already observes. Contract section 5.
- **Refs:** DEC-272; `public-interface.md` sections 4, 5 and 13;
  PRD-RUN-005; PRD-GATE-003; PRD-EVID-003.

### DEC-274: Add detail read leaves for Concerns, Deviations and Human Outcomes (D2 = A)

- **Date:** 2026-10-04. ACCEPTED by the Human ("A").
- **Context:** DEC-272 A7 shows recorded reasoning cut to one line in lists
  and in full in detail views. Concerns, Deviations and Human Outcomes have
  only `list` leaves, so their full reasoning was readable only through
  `--json` or the recording command's own output. Findings are already
  covered by `review get`.
- **Decision:** add read leaves `concern.get` / `fid concern get <id>`,
  `deviation.get` and `human-outcome.get`, returning `{concern:Concern}`,
  `{deviation:Deviation}` and `{outcome:HumanOutcome}`. Their Human views
  show recorded reasoning in full. Not-found reuses `concern-not-found`,
  `deviation-not-found` and `outcome-not-found`, which gain these producers.
  Option B (`--full` on lists) and option C (JSON only) were rejected.
- **Consequences:** three new public read leaves (task 12). The live leaf
  inventory test (R14), the Phase 23 section 5.1 producer table and the
  Human fixtures gain rows. Capability rows are unchanged.
- **Refs:** DEC-272 A7; `public-interface.md` sections 10 and 13.

### DEC-275: Store per-field request digests to name the differing field (D3 = A)

- **Date:** 2026-10-04. ACCEPTED by the Human ("A").
- **Context:** DEC-272 requires `idempotency-conflict` to name the differing
  field. fid stores one digest per request, so it cannot tell which part
  changed.
- **Decision:** every keyed mutation also stores one digest per top-level
  public input name (`--actor` included). These are fingerprints, never
  input values. Migration `0019` adds
  `idempotency_results.request_field_digests_json`, and the idempotency
  Effect of every keyed Event carries it, so rebuild reproduces it. Conflict
  details gain `differing_fields`. Option B (keep only whole-request
  digests) was rejected.
- **Consequences:** projection schema 19, and manifest `projection_schema` 19
  at phase end. Down migration restores schema 18 exactly (DEC-264).
  Pre-Phase-24 journals fail closed (DEC-270). Task 1 owns storage; task 2
  adopts the rule at all 43 digest call sites (inventory R1-R4).
- **Refs:** DEC-272; `public-interface.md` sections 6 and 7.

### DEC-276: `--help` input documentation and `--example` template rules (D4)

- **Date:** 2026-10-04. ACCEPTED by the Human ("yes, accept all").
- **Context:** DEC-272 option C requires every `--input` leaf to document its
  input and print an editable template. Three choices were open.
- **Decision:** templates use static, correctly typed placeholders (never
  pre-filled from the workspace or repository). They contain required fields
  only; optional fields are documented in `--help`. `--json --example` is
  `invalid-argument` because a template is not a `fid-cli.v1` envelope.
  `--example` needs no project, Actor, key or input and writes nothing. Help
  lists every input field with type, requirement, bounds or allowed values
  and a one-line meaning, then names the template command.
- **Consequences:** templates are deterministic and pinned exactly by tests.
  An expectation such as `expected_repository` is always typed by the
  caller, never rubber-stamped from fid's observation.
- **Refs:** DEC-272 Part 2; `public-interface.md` section 8.

### DEC-277: Phase 24 plan approved (contract, Heavy mode, 14 tasks)

- **Date:** 2026-10-04. ACCEPTED by the Human ("Approved").
- **Decision:** the following are approved, with D1-D4 as DEC-273..276:
  - the contract amendment `public-interface.md`, including the Tech Lead
    details the Human was shown: Source digests checked wherever supplied,
    `repository-unavailable`, run close/abandon expectation mismatch as
    `repository-mismatch`, and field-path `invalid-argument`;
  - Heavy mode;
  - the 14-task order in the phase brief.

  The Tech Lead added the *same worktree* rule (contract section 5) when
  drafting assertions. It preserves Phase 22 import semantics and generalizes
  the existing `importedGateDestination` and `importedDestinationBinding`
  rule. It adds no new behaviour.
- **Spec edits under this decision:** `01_PRD.md` gains PRD-OPS-011..014 and
  PRD-PROJ-006/007, which give the validation contract durable sources. The
  roadmap section 10 records the plan.
- **Validation contract:** 57 assertions (52 deterministic, 5 Human), all
  `Draft` until the QA Lead's QA-plan pass. VAL-044/047/056 restate Phase 23
  VAL-060/073.
- **QA cadence (proposed to the QA Lead):**
  - windows after tasks 3, 8, 11 and 14;
  - a Human walk-through of the new failure and Stop output at the
    task-11 window (Phase 23 workflow notes);
  - Mode B held-out runs `N/A` (no harness recorded), as in Phases 21-23.
- **Refs:** DEC-272..276; `qa/validation-contract.md`; `input-inventory.md`.

## Spec Sync Verdict

Recorded at wrap-up.

## Triage-Log Audit

Recorded at wrap-up.

## Workflow Notes

Recorded at wrap-up.
