---
status: accepted
date: 2026-08-08
---

# Supersede Phase 20 and defer configurable role policy

Phase 20 as previously planned is superseded and must not begin. Migrating the
V1 role/action matrix into `role_permissions`, adding `fid role create`, and
adding `fid policy set` are not part of the V2 backend-port critical path.
They are deferred, conditional standalone-fid features that require dogfood
evidence and a new roadmap decision before implementation.

Existing V1 issue-CLI authorization behavior and `fid policy show` remain
supported for compatibility. Already implemented Issue workflow authoring and
assignment remain core tracker capabilities. Custom state names that lack a V1
permission action remain non-transitionable until generic actor/capability
enforcement is delivered; fid must not grant them implicitly.

The V2 backend port instead receives generic actor identity, capability, and
recorded-authority enforcement in the relevant vertical slice. The kit owns
role selection, routing, and policy and records the approved grants in a Work
Envelope; fid validates the actor and enforces those recorded grants without
requiring its own configurable global role taxonomy.

## Consequences

- ADR-0002's committed v0.2 migration timing and its promise that new roles
  arrive through `fid role create` are superseded. Its principle that fid should
  not dictate one workflow remains useful but no longer commits an
  implementation.
- The reserved `role_permissions` schema and `fid role create` / `fid policy
  set` sketches remain non-executable design input, not scheduled scope.
- No compatibility surface is deleted or broken by this decision.
- The next roadmap work is the joint backend-port reslice, not Phase 20.
