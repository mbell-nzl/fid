# Repo Skills

Canonical starter-kit skills live under this directory.

Each skill should be stored as:

```text
project-docs/skills/<skill-name>/SKILL.md
```

Provider-discoverable copies are generated from these canonical files:

- Codex: `.agents/skills/<skill-name>/SKILL.md`
- Claude: `.claude/skills/<skill-name>/SKILL.md`

Generated copies are derived files and should not be edited directly. Update
the canonical skill, regenerate provider copies, then run workflow-lint.

```text
project-docs/tools/sync-skills.sh
project-docs/tools/workflow-lint.sh
```

Every skill must include a `Provenance:` line in the Markdown body using the
canonical terminology from `project-docs/specs/GLOSSARY.md`.
