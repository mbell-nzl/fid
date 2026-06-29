# DESIGN PLAN FLOWS — fid

> End-to-end user flows, whet artifact mapping, and export layout.
> The parent Design Plan (`05_DESIGN_PLAN.md`) has experience principles
> and output standards. This document shows how the commands compose
> into complete workflows.
>
> Source: `scaffold/initial-spec.md` §11–§13.

---

## 1. End-to-End User Flow

### 1.1 Sequential Epic (Bootstrap Pattern)

**Step 1 — TL Planning.**

TL creates the epic and its tasks with dependencies and ACs:

```
fid issue create --kind epic --title "Core persistence layer" \
    --scheduling sequential --by tl
# → created issue FID-e1f2g3 (epic) "Core persistence layer"

fid issue create --kind task --title "Implement ID generation" \
    --epic FID-e1f2g3 --priority 1 --sequence 1 \
    --spec-md project-docs/docs/phases/phase-1/tasks/FID-a1b2c3.md \
    --est-tokens 15000 --label core --by tl

fid issue ac add FID-a1b2c3 --text "IDs follow <PREFIX>-<hash> shape" --ordinal 1 --by tl
fid issue ac add FID-a1b2c3 --text "Progressive length scaling triggers at threshold" --ordinal 2 --by tl
fid issue ac add FID-a1b2c3 --text "Collision retry produces unique ID" --ordinal 3 --by tl

fid link add FID-a1b2c3 implements PRD-ISSUE-001 --by tl
```

**Step 2 — Workflow Runner Spawn.**

Runner queries the ready queue and spawns Developer:

```
fid ready --epic FID-e1f2g3
# → {"task": {"id": "FID-a1b2c3", "title": "Implement ID generation", ...}}
```

**Step 3 — Developer Execution.**

Developer claims, works, verifies ACs, records gate results, and finishes:

```
fid dev start FID-a1b2c3 --model claude-sonnet-4-6 --by developer
# → claimed by developer, transitioned todo → in_progress, opened RUN-b2c3d4

fid issue get FID-a1b2c3
# → full issue detail with ACs

# ... developer implements ...

fid issue ac verify FID-a1b2c3 AC-d4e5f6 --by developer \
    --status pass --evidence-type test \
    --evidence-ref "test/idgen_test.go::TestIDShape"

fid dev gate FID-a1b2c3 --cmd "go test ./..." --exit-code 0 \
    --output-excerpt "ok" --by developer
# → recorded reported gate result for attempt 1

fid dev finish FID-a1b2c3 --by developer \
    --completion-md project-docs/docs/phases/phase-1/tasks/FID-a1b2c3-done.md \
    --tokens 14823 --commits abc123
# → closed run, created completion TC-7c3e9f, transitioned → dev_complete, released claim
```

**Step 4 — Workflow Runner Inspection.**

Runner checks for concerns:

```
fid issue get FID-a1b2c3
# → concerns_raised: true → spawn TL-Lite
# → concerns_raised: false → continue to next task
```

If concerns exist, TL-Lite triages:

```
fid triage start FID-e1f2g3 --source dev_concern --ref CON-c3d4e5 --by tl_lite
fid triage verdict --verdict RESUME --by tl_lite \
    --reasoning-md project-docs/docs/phases/phase-1/triage/CON-c3d4e5-triage.md
```

**Step 5 — Mid-Epic QA Review (every 3 tasks, v0.3+).**

Runner creates a windowed review:

```
fid qa start FID-e1f2g3 --scope epic_window --by qa_reviewer
# QA Reviewer acquires review, adds findings:
fid qa finding QAR-review1 --by qa_reviewer \
    --severity blocking --task FID-a1b2c3 --finding-md path
fid qa verify FID-a1b2c3 AC-d4e5f6 --by qa_lead \
    --status pass --evidence-type command --evidence-ref "go test ./..."
fid qa finish QAR-review1 --by qa_lead
```

**Step 6 — Final Epic QA Lead Review (v0.3+).**

After all tasks complete:

```
fid qa start FID-e1f2g3 --scope epic_full --by qa_lead
# QA Lead runs full review, verifies remaining ACs at qa_lead tier
fid qa finish QAR-final --by qa_lead
```

**Step 7 — TL Wrap-Up.**

TL queries the epic, closes escalations, and records wrap-up notes:

```
fid tl wrapup FID-e1f2g3 --by tl
# → recorded wrap-up notes
# → closed 2 resolved escalations
# → transitioned FID-e1f2g3: in_progress → wrapped

# Human approval:
fid issue transition FID-e1f2g3 --to closed --by human
```

### 1.2 Prioritised Epic or Unattached Task (Steady-State)

Same lifecycle, minus batched QA:

- QA review (if any) is scoped `single_task`.
- Optional retrospective.
- The `ready` queue uses priority + dependencies, not sequence order.
- Single-task tickets may use a lighter workflow that skips formal
  `qa_review` for trivial bugs.

---

## 2. whet Artifact Mapping

For the bootstrap-of-`fid`-on-`whet` case and the v0.1 dogfood goal
("replaces ~70% of daily interaction with whet's task artifacts"):

| `whet` Artifact | `fid` Entity | Notes |
|---|---|---|
| Phase | `issues` (kind='epic') + `epic_metadata` (scheduling='sequential') | Phase brief becomes `epic_metadata.brief_md_path`. |
| Task spec (`task-M.md`) | `issues` (kind='task') + `acceptance_criteria` rows + `spec_md_path` | ACs extracted into rows; rest stays markdown. |
| Task completion report (`task-M-done.md`) | `task_completions` row + `notes_md_path`; `## Concerns` section parsed into `concerns` rows | Append-only; each concern is a structured row. |
| QA Reviewer mid-phase (`review-M.md`) | `qa_reviews` (scope='epic_window') + `qa_findings` rows | Findings become queryable rows. |
| QA Lead phase-end (`final-review.md`) | `qa_reviews` (scope='epic_full') + `qa_findings` rows | Evidence stored per finding via `evidence_excerpt`. |
| QA plan (`qa-plan.md`) | `qa_plans` table | First-class verification design. |
| `decisions_active.md` | `decisions` rows (status='active') + `carry_forward` rows | Rendered on demand; file replaced by query. |
| `decisions_principles.md` | `carry_forward` rows with `promoted_principle = 1` | Promotion is a column flip. |
| `tl_escalations_current.md` | `escalations` rows (status='open') | Rendered on demand; file replaced by query. |
| `triage-log.md` | `triage_events` rows | Per-epic, queryable across epics. |
| `progress_current.md` | Computed view over status fields | No primary entity; rendered on read. |
| Phase decisions (`phases/phase-N/decisions.md`) | `decisions` rows linked to epic; v0.4 may additionally create a `retrospectives` row (type='lightweight') | Workflow notes stay markdown until retrospective records ship. |
| Deep retro (`docs/reviews/...`) | v0.4 `retrospectives` row (type='deep') | |
| `tl_handoff.md` | Out of scope for v0.1; remains markdown | Optional handoff between TL sessions. |

**Not migrated:** `AGENTS.md`, `CLAUDE.md`, golden docs (`01_PRD.md`,
`02_TRD.md`, etc.), and workflow docs stay as markdown. They are project
canon, not project state.

---

## 3. Export Layout (v0.4)

### 3.1 Default `whet` Layout

Mirrors the existing `whet-private-docs/` structure. Directory and file
names use the entity's ID:

```
docs/
├── current/
│   ├── phase_brief.md            (rendered from active epic + epic_metadata)
│   ├── decisions_active.md       (rendered from decisions WHERE status='active')
│   ├── decisions_principles.md   (rendered from carry_forward WHERE promoted_principle)
│   ├── progress_current.md       (computed view over status fields)
│   ├── tl_escalations_current.md (rendered from escalations WHERE status='open')
│   └── tl_handoff.md             (optional, per-session)
├── epics/
│   └── <epic-id>/                (e.g. WHET-a1b2c3/)
│       ├── tasks/
│       │   ├── <task-id>.md      (issues row + spec_md_path + ACs as `- [x]`)
│       │   └── <task-id>-done.md (latest task_completions row)
│       ├── qa/
│       │   ├── qa-plan.md
│       │   ├── review-<review-id>.md
│       │   ├── final-review.md
│       │   └── human-checklist.md
│       ├── decisions.md
│       ├── triage-log.md
│       └── progress.md
└── reviews/
    └── YYYY-MM-DD-<epic-id>-retrospective.md
```

### 3.2 Config

Controlled via `.fid/config.toml`:

```toml
[export]
markdown_root = "./docs/"
layout = "whet"               # 'whet' | 'flat' | 'custom'
banner = true                 # emit "generated by fid" banner
```

### 3.3 Frontmatter Contract

Every exported file has YAML frontmatter carrying structured fields:

```yaml
---
fid:
  id: WHET-a1b2c3
  kind: task
  status: dev_complete
  workflow: dev-flow
  epic: WHET-4f7c8d
  priority: 1
  sequence: 10
  labels: [bug, auth]
  est_tokens: 8000
  actual_tokens: 9120
  acs:
    - id: ac-1
      text: "..."
      verified: true
      verifier: developer
  links:
    - blocks: WHET-7c3e9f
    - implements: PRD-AUTH-001
---

# WHET-a1b2c3: Add OIDC support

## Objective
...
```

In v0.4 (export only) the frontmatter is informational — it makes
exported files self-describing and machine-readable by future tooling.
If/when v0.5+ ships import, the frontmatter contract is exactly what
that import path reads.

Round-tripping prose perfectly is impossible; the frontmatter contract
guarantees only that *structured fields* are reliably represented.

### 3.4 Export CLI

```
fid export --to ./docs/                         # full project tree
fid export --to ./docs/ --epic WHET-a1b2c3      # one epic
fid export --to ./docs/ --since 2026-04-01      # incremental
```

Exported files carry a banner:
```markdown
<!-- generated by fid; edits will be lost on next export -->
```

---

## 4. Context Packs (v0.2+)

`fid issue get <id> --context-pack <role>` returns a role-tuned bundle:

| Pack | Contents |
|---|---|
| `developer` | issue + ACs + spec_md + linked requirements + open comments + relevant template + active claim status + parent epic brief excerpt |
| `qa_reviewer` | issue + ACs (with latest verifications per role) + spec_md + completion notes + linked findings from prior reviews |
| `qa_lead` | same as qa_reviewer plus human-checklist generation hooks and stricter evidence requirements |
| `tl_lite` | concern excerpt + originating task spec + active escalations referencing this epic + relevant decisions |
| `tl` | full epic view: child task statuses, all decisions, open escalations, all triage events, retro draft if any |

---

*Last updated: 2026-05-27*
