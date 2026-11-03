---
status: accepted
date: 2026-10-31
---

# Reset fid to a ledger with receipts and the kit to a skills kit

fid becomes a local-first work ledger with receipts for agent-written code. It
tracks Issues with lanes and validations, records observed evidence,
independent reviews and approvals bound to the change, and refuses "done"
until the lane's evidence exists. People use a human CLI built as a facade over
the existing kernel; agents use `--json`. fid is no longer the kit's backend.

The kit becomes Mark's engineering method: curated skills, one reviewer agent
and a default loop, copied into each project. The only interface between the
two is a small loop contract that fid owns, versions and tests. The V1 phase
workflow is retired for fid's own development; `RESET.md` is the plan.

## Context

After five months neither product was in daily use. fid had 1.8M words of
process docs against about 47k lines of Go, a kernel of about fifteen nouns,
and 33 top-level commands. In Phase 24's HUMAN-01, Mark could not take control
of a stuck task from fid's output, and the failures were about concepts, not
formatting. Most review findings were weak tests rather than product bugs. fid
and the kit spent two months waiting on each other's contracts, although
neither had a user.

The reset proposal (`docs/reset/2026-10-29-fid-kit-reset-proposal.html`,
now superseded), its independent review, a decisions session and a second
fresh-context review produced the decisions in
`docs/reset/2026-10-31-reset-decisions.md`. This ADR records
what those decisions reverse. The decisions record holds the detail.

## Decision

- **Ownership.** fid owns work state, evidence and the evaluation of project
  lane policy. The kit owns technique: skills, the reviewer agent and the loop.
  Project files and git own prose and policy. Mark owns approval, waivers and
  the lowering of any lane.
- **Ledger.** The ledger lives in `$(git rev-parse --git-common-dir)/fid/`,
  shared by every worktree and never committed. Backups go to a private repo
  through `fid backup`.
- **Policy.** Checks, lanes and check-defining paths live in a committed
  `.fid.toml`, read only from `origin/<default branch>`. Changes to
  `.fid.toml` and CI workflows are always Assured.
- **Done.** A task is done when merged into `main` with CI green on that
  commit, in every lane. fid confirms through `gh`.
- **Approvals.** Reviews and approvals bind to the change's fingerprint, and
  checks bind to the tree. Approvals go stale when the change is edited.
- **Threat model.** Agents are honest but fallible, on a single-user machine.
  Hooks and prompts are friction; fid's refusals are the enforcement.
- **Lanes.** Lanes are per task: Routine, Standard and Assured. fid computes a
  floor from changed paths; agents and LLM judgement may raise a lane; only
  Mark lowers one, through a recorded waiver.
- **Merging.** Mark merges everything through the trial. A hold-only merge
  judge runs in shadow. Agent merges are enabled per lane by Mark through
  `.fid.toml`, and then only the loop script merges.
- **Facade.** Golden-path commands compose kernel commands. The facade's
  `done` evaluator is the single source of what done needs. Kernel commands
  move under `fid admin`; any unused at trial end are deleted. A test fails if
  kernel nouns or IDs reach human output.
- **Scope.** Nothing joins the reset unless a golden-path step or a kill
  criterion needs it. A net Go line budget applies.

## Consequences

- fid computes risk lanes, which Gate B1 forbade.
- fid makes network calls in three places: `fid done` confirms merges and CI
  through `gh`, `fid backup` pushes to the backup repo, and recording a review
  round posts it to the PR through `gh` (added 3 November 2026; decisions
  record, work model). Nothing else uses the network.
- Compatibility obligations still start at the first supported baseline, which
  is now the reset's migration squash in week 3, not the end of Gate D
  convergence.
- The V1 workflow docs (`project-docs/` phase archives, role agents, workflow
  verifiers) move to a private archive repo. Agents are not pointed at it.
- Roadmap Phases 25 and 26 are cancelled. Parallel claims, usage capture and
  CLI polish move into the reset; kit adapter qualification (CF-037) is
  cancelled; the rest goes on the later list.
- If the trial's kill criteria fail, the receipts core is kept as a small
  library or hook pack and fid stops as a product.

## Superseded and amended records

| Record | Clause | Change |
|---|---|---|
| ADR-0001 | `.fid/journal` committed to git; merges are textual JSONL line-merges handled by `git merge`; `.fid/config.toml` committed | **Superseded in part.** Ledger in the git common directory, never committed. Project policy in `.fid.toml` at the repo root. Journal-canonical, SQLite-derived and no daemon are retained. |
| ADR-0003 | The V1 three-tier QA procedure may remain how fid is developed until the V2 kit replaces it | **Ended.** The reset's working method replaces it. |
| ADR-0004 | fid as the kit's operational backend; the kit owns risk classification, routing, orchestration and gate activation; Gate C qualification before live use | **Superseded.** Ownership as above. No backend contract and no capability qualification. |
| ADR-0006 | Parallel mutation needs declared, provably disjoint scopes | **Amended.** Distinct Issues in separate worktrees may proceed without scope proofs; conflicts are resolved at merge. Proofs wait for real need. |
| ADR-0011 | The kit owns risk policy and provider projections; Kit Activation state; adopt kit Envelope semantics and never recalculate policy; convergence with the V2 kit | **Amended.** fid evaluates project lane policy. Activation survives only as "the policy fid trusts, from `origin/main` at commit X". The convergence target is withdrawn. Its record criteria and retirement list stand. |
| ADR-0012 | Deliver convergence as vertical scenario phases; build with the V1 process until the V2 kit is ready | **Superseded.** No phases. `RESET.md` tasks in small PRs. |
| Gate B1 (2026-08-08 direction session) | fid never infers a risk lane from files | **Reversed.** fid computes the lane floor from changed paths. |
| Gate B3 (same session) | The kit owns the gate schema and activation; fid never treats a mutable authoring file as runtime authority | **Amended.** Runtime authority is `.fid.toml` as committed on `origin/<default branch>`. Working-tree edits have no effect until merged. |
| Charter §2 | Risk classification, routing, obligation selection and orchestration policy belong to the kit | **Amended.** Lane policy is project-owned in `.fid.toml` and evaluated by fid. Orchestration is the kit's loop. |
| Charter §6 | Out of scope: owning risk policy; validation contract and assertion aggregates | **Amended.** Lanes and epic-level Validations are in scope. Hosting agents or the loop, roles, a daemon and the rest stay out of scope. |
| Charter §7 | First supported baseline follows Gate D convergence | **Amended.** The baseline is the reset's migration squash. |
| Charter §8 | Success criteria: contract qualification, V2 adoption | **Superseded** by the kill criteria in `RESET.md`. |
| Charter §9 | No network by default | **Amended.** Network only for `fid done` confirmation, `fid backup` and posting review rounds to the PR (the last added 3 November 2026). |
| `10_IMPLEMENTATION_ROADMAP.md` | The whole roadmap, including §6 removal rules and Phases 25–26 | **Superseded** by `RESET.md` and the later list. |
| `09_AI_WORKFLOW.md` | The V1 workflow, including light, normal and heavy operating modes | **Superseded for fid.** Per-task lanes and the working method in `RESET.md`. |

ADR-0005, ADR-0007, ADR-0008 and ADR-0009 are retained. Interrupted work still
resumes through a linked successor Run, and review rework stays on the
original Issue. The facade uses one-Issue Envelopes. ADR-0002 and ADR-0010
were already superseded or completed. DEC records up to DEC-293 remain
historical.
