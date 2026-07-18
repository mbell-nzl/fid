# Repo Skills

Canonical starter-kit skills live under this directory.

Each skill should be stored as:

```text
project-docs/skills/<skill-name>/SKILL.md
```

Provider-discoverable copies are generated from these canonical files:

- Codex: `.agents/skills/<skill-name>/SKILL.md`
- Claude: `.claude/skills/<skill-name>/SKILL.md`

`.agents/skills/` is consumed natively by Codex skill discovery, which
searches the working directory and parents up to the repo root, then
`$HOME/.agents/skills`, then `/etc/codex/skills` (symlinks followed). Codex
invokes a skill implicitly when the task matches its description, or
explicitly with `$skill-name` / the `/skills` picker; each skill costs only
its name + description in context until invoked, so front-load trigger
language in descriptions. Optional per-skill `agents/openai.yaml` (invocation
policy, tool dependencies) is a future option, deliberately not used yet.

Generated copies are derived files and should not be edited directly. Update
the canonical skill, regenerate provider copies, then run workflow-lint.

```text
project-docs/tools/sync-skills.sh
project-docs/tools/workflow-lint.sh
```

Every skill must include a `Provenance:` line in the Markdown body using the
canonical terminology from `project-docs/specs/GLOSSARY.md`.
