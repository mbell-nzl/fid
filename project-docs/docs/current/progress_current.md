# Progress - Current

> In-flight status only. Historical implementation progress remains in the
> phase archives.

## Project state

**Phase 20 is in flight: Tasks 1-3b are complete and Reviews 1 and 1b both
returned BLOCK.** `ESC-20-001` (DEC-130) and `ESC-20-002` (DEC-131) are
resolved with no weakened acceptance criterion and no product defect: the
findings were missing evidence, a stale integration fixture, and — the real
problem — full-suite gate records that were never captured from real output.
Phase 19 completed and wrapped on 2026-08-08. fid has no adopted project
workspace or compatibility constituency.

Gates D1/D2 supersede the old Phase 20 role-policy proposal and the historical
process-shaped roadmap. Gate D3.1 establishes the new Phase 20 as the first of
six vertical scenario slices. The target product remains the standalone Issue
tracker plus the fid/V2 durable operational kernel described by ADR-0011.

Branch: `phase-20`.

## Gate D2 documentation

Completed and committed:

- product decision ADR, decision-session ledger, and reconciliation report;
- charter, PRD, and glossary rewrite;
- technical, logical-schema, and runtime-configuration rewrite;
- CLI, Human-interface, and end-to-end flow rewrite; and
- quality, V2 workflow-integration, and convergence-roadmap rewrite.

The rewrite removes forward promises for global role policy, role composites,
QA Plans, Validation Contracts, templates, skills, handoffs, routine
retrospectives, generic TOON/YAML/Markdown output, and pre-dogfood
compatibility. JSON is the sole canonical machine representation.

## Implementation baseline

Phase 19's Issue workflow authoring implementation remains useful foundation,
not a frozen public contract. Retain useful Issue, journal, SQLite, rebuild,
locking, identity, and repository primitives only where they conform to the
new specs.

Product code has not yet been removed or reshaped under ADR-0011. ADR-0012 now
assigns replacement and removal to phase-sized vertical scenarios beginning
with the standalone Issue journey.

## Carry-forward

The active disposition is in `decisions_active.md`. In summary:

- cancel `internal/dev` refactoring and generic Markdown formatter work;
- keep Issue-path and rebuild-performance debt only when their retained paths
  reopen;
- treat workflow inspection/lifecycle as candidates, not commitments; and
- keep V1 process-tooling concerns separate from product scope.

## Phase 20 plan

Planned, self-checked, and independently audited on 2026-08-09 under the binding
V1 delivery process. Heavy mode, eighteen bounded Developer tasks, six
mid-phase reviews, `VAL-001..VAL-046`, decisions `DEC-107..DEC-129`. Artifacts are under
`../phases/phase-20/`; the accepted slice contract is
`../phases/phase-20/contracts.md`.

Human approval accepted amendments to four specs before implementation, as the
DDL and JSON contract gates require: `03_TRD_SCHEMA.md` §8 (migrations
`0011_p20_remove_issue_fields` and `0012_p20_issue_revision`),
`04_TRD_CONFIG.md` §2/§5 (versioned strict config and atomic init),
`06_DESIGN_PLAN_CLI.md` §2/§5/§8 (Response Envelope rules, Phase 20 reason
codes, slice contract index), and `GLOSSARY.md` (Machine Interface terms).

Phase 20 preserves fail-closed V1 authority and pulls no Phase 21 Routine or
Envelope work forward. Issue revision is aggregate-wide; link endpoint types are
explicit; `journal inspect` stays legacy rather than publishing V1 Event internals.

## Current work

1. Review 1c is complete and NON-BLOCK at commit `cb5f1bf`.
2. The Human-approved DEC-133 infrastructure pack Tasks 3e-3l is complete;
   `../phases/phase-20/infra-improvement-pack.md` retains the execution record.
3. Infrastructure Review 1's ESC-20-004 correction is independently NON-BLOCK
   at `cd9e6bc`; all final post-review gates pass.
4. Phase Runner resumes at Task 4, then continues Tasks 5-18 through the six
   accepted review windows.
5. Resolve Gate D3.2's exact Routine walking-skeleton boundary before Phase 21
   planning.

Carry-forward opened this phase, both registered in `decisions_active.md` at
wrap-up:

- **CF-028** — `internal/bootstrap/init.go` is 427 lines against the 400-line
  budget; split staging/fsync helpers only when a later task changes that file
  for another reason.
- **CF-029** — extend the evidence-token control to QA review artifacts, so a
  review cannot record a token its command could not have produced.
- **CF-030** — make `report-lint.sh`'s rejected-token message state the expected
  shape (a real `ok`/`FAIL`/`PASS` line with package and duration), so a
  truncated-but-honest token is corrected without an escalation (DEC-136).

All Phase 20 full-suite `go test` gate lines recorded as `exit 0 (no output)`
are **UNVERIFIED**, not `PASS` (DEC-131). Task-3c's run on the current tree is
the phase's verified gate baseline.

No implementation task should cite the superseded historical version roadmap
as forward authority.
