# Repo Tools

Stack-neutral repository commands that do not require project runtime
dependencies.

| Command | Purpose |
|---|---|
| `sync-skills.sh` | Generate provider skill copies from `project-docs/skills/`. Use `--check` to compare without writing. |
| `workflow-lint.sh` | Verify agent workflow artifacts, templates, config, and generated skill copies have not drifted. |

Run from the repository root:

```sh
project-docs/tools/sync-skills.sh
project-docs/tools/workflow-lint.sh
```
