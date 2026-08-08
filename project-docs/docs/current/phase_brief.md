# Phase Brief - No Active Phase

> Default execution brief loaded by every agent session.

## Status

**No phase is active.** Phase 19 is complete and wrapped. The old Phase 20
role-policy proposal is superseded by ADR-0010/0011. ADR-0012 authorizes the V1
Tech Lead to plan a new Phase 20 for the clean standalone Issue journey; no
`phase-20/` artifact directory currently exists.

## Current product direction

fid remains a local-first Issue tracker and becomes the durable operational
backend for the V2 kit. The project and kit own engineering intent, policy, and
orchestration. fid owns Issues, immutable execution contracts, authority,
Claims/Runs, trusted evidence, reviews, Human Outcomes, Attention, recovery,
and semantic projections.

There is no production compatibility obligation. Superseded V1 commands,
models, migrations, event fixtures, and generic TOON/YAML/Markdown formatting
will be replaced or removed before first dogfood. JSON is the only canonical
machine representation.

## Current checkpoint

Gate D2's fid specs and the V2 kit contract are aligned. An independent
cross-repository audit found stale event-snapshot, Receipt-ingestion, role,
operator-doc, report, and traceability clauses; those findings have been
reconciled in the canonical documents.

Gate D3.1 now fixes the six vertical scenario phases and their legacy-removal
owners. The approved sequence is in
`../../specs/10_IMPLEMENTATION_ROADMAP.md` section 5. Gate D3.2 remains for the
exact Phase 21 Routine walking-skeleton boundary and does not block Phase 20
planning.

## Next implementation phase

The V1 Tech Lead may now plan Phase 20. It must define one coherent standalone
Issue slice with:

- accepted terminology, JSON schemas, errors, events, and invariants;
- exact retained and removed surfaces;
- shared fake/fid conformance cases;
- a first failing behavioral test;
- recovery and migration consequences; and
- current repository quality commands.

Phase 20 must remove generic TOON/YAML/Markdown formatting, `--format`, format-
matrix tests, and obsolete command stubs. It must not remove fail-closed V1
authority until Phase 21 provides and proves generic authority.

Useful implemented foundation may be retained, but its V1 public shape is not
a compatibility contract. The generic replacement must preserve fail-closed
authority or evidence behavior before safety-relevant V1 code disappears.

## Temporary delivery process

`AGENTS.md` remains binding while fid is developed with the V1 repository
process. Its role, validation-contract, task-report, and evidence rules are
delivery scaffolding until the V2 kit can replace them; they are not fid product
features and must not leak back into the target model.

## Active cautions

- Do not implement global role-policy tables or `fid role/policy` commands.
- Do not add QA Plan, Validation Contract, template, handoff, triage, or
  completion-report aggregates.
- Do not trust caller-reported command outcomes as Receipts.
- Do not read SQLite, journal, or `.fid/` internals from the kit.
- Do not use the in-memory fake as a production fallback.
- Do not create a second durable work graph beside Issues.
- Do not add output encodings beside canonical JSON without satisfying the
  ADR-0011 evidence trigger.
