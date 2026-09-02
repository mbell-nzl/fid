# TL Escalations - Current

> Open items that exceed current repair authority and require full Tech Lead or
> Human input. Resolved detail lives in each phase archive.

## Open Escalations

None. DEC-160 through DEC-173 resolve ESC-21-001..014. The run is stopped after
TRG-21-020; DEC-173 authorizes the task-20 attempt-2 continuation that clears it.

## Latest Resolution

DEC-173 resolved ESC-21-014 on 2026-09-02. TRG-21-020 correctly stopped task
20: its required absence controls walk `cmd/fid internal testdata` and its
coverage gate compiles `./...`, but its file list omitted `internal/bootstrap/`,
`internal/journal/`, `internal/labels/`, `internal/links/`,
`internal/requirements/`, `internal/db/`, `testdata/journal/`, and
`cmd/fid/integration/`. The inventory already names task 20 co-owner of every
residue file; the file list was the incomplete artifact.

Task 20's scope is amended with bounded subject-preserving dispositions: delete
`LockOwner.ActorRole`/`actor_role`, the bootstrap fixed-role registry, and the
`ActorRole`/`legacyRole` fallback readers; converge the `dev.*` golden and every
`roles` fixture to surviving families/tables; converge the tagged closure
journey to accepted behaviour. Task 20 gains the tagged integration gate
(DEC-167). No command, expected status, VAL, 0014, frozen-table, or public
contract changes. `08_QUALITY_PLAN.md` §3 gains the converse executable-
ownership rule. Recovery is RUN-RESUME under TRG-21-020 / ESC-21-014 / DEC-173,
task-20 attempt 2, `in-scope-implementation` in task scope; task 21 remains
gated on task-20 completion.

## Archive Pointers

- Phase 21: `../phases/phase-21/tl-escalations.md`
- Phase 20: `../phases/phase-20/tl-escalations.md`
- Phase 19: `../phases/phase-19/tl-escalations.md`
- Earlier phases: their `../phases/phase-N/tl-escalations.md` archives.

## Protected Escalation Triggers

Escalate only for a required product, architecture, data-model,
public-interface, dependency, security/data-loss, specification, exhausted
recovery-budget, or unavailable external-capability decision. Ordinary red
tests, failed gates, report defects, QA rework, and incomplete in-scope work
remain bounded repair. Before claiming recovery exhaustion, confirm each
counted attempt returned a typed result and its directive permitted causal
repair. Externally interrupted attempts are session-liveness evidence, not
repair-capability evidence.
