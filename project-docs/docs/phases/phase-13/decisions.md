# Phase 13 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes, and
> wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-032: `fid label add` rejects duplicate names before append and stores optional metadata as given

- **Date:** 2026-06-27
- **Context:** Phase 13 Task 1 introduces the first label-creation surface.
  The validation contract's PI-1 fixes the v0.1 semantics that must hold
  without a schema migration or projector change: `labels.name` is the primary
  key, `colour` and `description` are optional, and label creation must not
  project to `activity_log`.
- **Decision:** `fid label add <name> [--colour <hex>] [--description <text>]
  --by <role>` validates and authorizes before append, refreshes the
  projection with the held-lock helper, and appends exactly one `label.added`
  event containing exactly one `labels` insert effect. A second add of the same
  label name is rejected before append with an actionable error naming the
  existing label. `--colour` and `--description` are optional and, when
  provided, are stored verbatim after CLI whitespace trimming; no additional
  colour-format validation is introduced in v0.1. Label creation projects no
  `activity_log` row.
- **Consequences:** The existing `labels` table and generic insert projector are
  sufficient, so Task 1 adds no schema migration and does not modify
  `internal/journal/project_effects.go`. `fid label list` can read the stored
  metadata directly from the `labels` table across json, toon, markdown, and
  yaml. No section-13 mismatch was found: the in-code `label.add` matrix row
  matches `03_TRD_SCHEMA.md` section 13 (`developer`, `tl`,
  `workflow_runner`, and `human` only).
- **Refs:** PI-1; phase-13/task-1; VAL-001; VAL-002; VAL-003; VAL-004;
  VAL-007; `03_TRD_SCHEMA.md` sections 9.6 and 13; `06_DESIGN_PLAN_CLI.md`
  section 9 `fid label add` / `fid label list`.
