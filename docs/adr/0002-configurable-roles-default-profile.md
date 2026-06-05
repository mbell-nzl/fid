# Roles configurable via a seeded default profile

## Decision

Agent roles and the role × action permissions matrix are configurable per project, seeded at `fid init` from a built-in default profile. Roles are stored in a `roles` table from v0.1; the permissions matrix migrates from app-code constants to a `role_permissions` table in v0.2, paralleling the workflow-engine migration. The default profile name is `default`. Unknown roles are hard-rejected by the CLI in v0.1 (configurable in v0.2 once `role_permissions` is data).

## Context

The spec's audience statement promises "an agent-friendly lightweight issue tracker that helps with my workflow but does not dictate one." The original spec hard-coded a 7-role enum, a static role × action matrix, role-namespace composite commands (`fid dev start`, etc.), and role-keyed context packs. Those *are* the workflow shape — bake them in and fid is opinionated about workflow, contradicting the framing.

Reference points:

- **Beads / Beans** carry no built-in role concept; "who is doing this" lives in the prompts and scripts that *use* the tracker. Maximally generic; loses multi-role AC verification and write-time policy enforcement.
- **Linear** has roles as identity (member, admin), not workflow positions.
- **Jira** has workflow customisation per project, including roles.

fid sits between Linear and Jira: roles are workflow positions (developer, qa_reviewer, …) with permission boundaries, but the *set* of roles should be a project-level choice, not a fork of fid.

## Considered options

- **(P-A) Embrace the prescription.** Keep roles + matrix hard-coded; drop the "doesn't dictate" framing. Smallest change, but contradicts the stated intent.
- **(P-B) Decouple identity from permissions.** Keep role names as identity tags; drop or weaken the matrix. Loses the policy-enforcement feature that prevents (e.g.) Developer from creating decisions.
- **(P-C) Configurable roles seeded with a default profile.** Chosen.
- **(P-D) Drop the role concept entirely (Beads/Beans style).** Loses multi-role AC verification, the §6.13 policy enforcement, and the role-tuned context packs that materially reduce agent context-window pressure.

## Consequences

- **Schema gains a `roles` table** in v0.1, seeded from the default profile by the `project.initialized` journal event written at `fid init`. The columns `responsible_role`, `verifier_role`, `claimed_by`, etc. lose their app-code CHECK constraints; existence in `roles` is the validity rule.
- **Permissions stay in code in v0.1** but the matrix is the sole source of truth and is documented as "the seeded default" rather than "the only valid set". v0.2 migrates it to a `role_permissions` table alongside the workflow-engine migration. Same pattern, same release.
- **Composite commands (`fid dev start`) and context packs (`--context-pack developer`) remain hard-coded through v0.2.** Customisation of these is a v0.3+ concern, gated on real demand. Users adding new roles in v0.1 understand that composites/packs cover the default roles only.
- **Unknown roles hard-reject in v0.1.** The matrix lives in code; it can't reason about roles it doesn't know. Users who need a new role wait for v0.2 or fork. A single CLI error frames this clearly: "role X not in profile; v0.2 will add `fid role create`."
- **Default profile = `default`.** Reserves namespace for alternative built-in profiles (`solo-dev`, `pair`, future patterns) without one workflow owning "default".
- **Workflow transition rules that require role-specific verification are parameterized.** For example, `all_acs_verified` carries required verifier roles in rule params (`{"roles":["developer"]}` or `{"roles":["developer","qa_lead","tl","human"]}`) instead of minting role-specific rule names. When `role_permissions` migrates to data in v0.2, transition rules validate those role names against the same profile.
- **No retroactive breakage at v0.2.** v0.1 users get the migration for free: the v0.2 schema migration seeds `role_permissions` from the v0.1 code constants.
