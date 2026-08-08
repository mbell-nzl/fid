---
status: accepted
date: 2026-08-08
---

# fid owns operational work state; the kit owns engineering workflow policy

fid remains a local-first, agent-first issue tracker and becomes the durable
operational work-state and evidence backend for kit-driven agentic engineering.
fid owns issues and their lifecycle, work relationships, claims and runs,
ordered events and projections, receipts and evidence, attention, outcomes,
telemetry, rebuild, and portable recovery. The kit owns risk classification,
authority and role routing, engineering-workflow policy, deterministic
orchestration, and provider projections. Project files and Git remain canonical
for requirements, architecture, decisions, code, and repository identity;
Humans remain authoritative for product judgment and protected decisions.

## Context

The alternative boundaries were either to make fid encode the whole engineering
method, which would couple its usefulness to one role and artifact model, or to
put durable workflow state in the kit, which would duplicate fid's mature
journal, projection, issue, rebuild, and query foundations. Splitting durable
workflow truth between both products was rejected because recovery and resume
could disagree. The chosen boundary keeps fid independently useful as an issue
tracker while allowing kit policy and frontier-model usage to evolve above a
stable, versioned operational contract.

## Consequences

- Issue tracking remains at fid's product center; the backend-port work extends
  the issue tracker rather than replacing it.
- A Work Envelope is an execution contract around selected work, not a
  replacement issue. ADR-0005 resolves the identity relationship: fid Issues
  are the only durable work graph; Envelopes reference selected Issues and do
  not introduce persistent Work Items.
- The kit evaluates engineering-workflow policy and compiles the approved Work
  Envelope. fid does not recalculate risk, routing, or obligation-selection
  policy; it enforces generic integrity and the recorded envelope's revision,
  lease, idempotency, authorization, non-weakening, evidence, freshness, and
  lifecycle preconditions. Shared conformance fixtures guard the public
  contract.
- fid is the primary trusted local executor for registered gates and creates
  receipts from observed execution. External attested runners are a secondary
  extension for CI, Assured work, and release gates; caller-reported results
  never satisfy trusted obligations.
- Projects own approved machine-readable Gate Definition values; the kit owns
  their schema, validation, and activation procedure; fid owns immutable
  registrations and never treats mutable authoring configuration as runtime
  authority. The authoring filename is an interface-design detail.
- The separately approved Gate C delivery decision requires qualification of
  all twelve V2 kit backend-port capabilities before live-project migration,
  delivered through bounded, independently evidenced slices. Roadmap order and
  the disposition of existing role, QA, template, and report surfaces remain
  open decisions.
- ADR-0002 and ADR-0003 continue to describe existing/shipped compatibility
  assumptions until their forward-roadmap disposition is decided; they do not
  override this cross-product ownership boundary.
