# Task briefs

Approved briefs live at `docs/briefs/<task>.md`: use the RESET row ID when
there is one, otherwise the fid ID. Set the task's `spec_md_path` to that
repo-relative path. The brief must reach `main` before the loop builds it.

The orchestrator role writes briefs with Mark during planning, or on its own
for agent-filed tasks whose design is already settled. An unsettled design
goes back to Mark; the builder parks with its question instead of guessing.
