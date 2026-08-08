# Phase Brief - No Active Phase

> Default execution brief loaded by every agent session.

## Status

**No phase is active.** Phase 19 is complete and wrapped. Phase 20 is
superseded by ADR-0010/0011. Do not begin product implementation from this
file or from the historical version roadmap.

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

Gate D2's canonical fid spec rewrite is complete. Remaining governance work:

1. align the V2 kit documents;
2. complete an independent cross-repository consistency audit;
3. fix material findings and refresh the Human report; and
4. discuss and ratify Gate D3's vertical-slice order.

The candidate sequence and open ordering questions are in
`../../specs/10_IMPLEMENTATION_ROADMAP.md` section 5.

## Next implementation phase

The next phase begins only after Gate D3 and must define one coherent vertical
slice with:

- accepted terminology, JSON schemas, errors, events, and invariants;
- exact retained and removed surfaces;
- shared fake/fid conformance cases;
- a first failing behavioral test;
- recovery and migration consequences; and
- current repository quality commands.

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
