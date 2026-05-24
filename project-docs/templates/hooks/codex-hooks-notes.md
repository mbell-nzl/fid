# Codex Hook Notes

Codex does not currently support a project-level hook configuration file
comparable to Claude Code's `.claude/settings.json` hooks.

Codex enforces deterministic gates through its sandbox and agent prompt
instructions rather than through pre/post-tool hooks.

## Recommended Codex Enforcement Strategy

1. **Agent prompts** — agent `.toml` files under `.codex/agents/` instruct
   agents to run quality gates. This is already in place.
2. **Sandbox policy** — `.codex/config.toml` controls which commands and
   network access are allowed. Use this to restrict destructive operations.
3. **CI/CD hooks** — run `project-docs/tools/workflow-lint.sh` as a CI step
   to catch prompt/template/skill drift before merge.
4. **Manual pre-commit** — developers can add workflow-lint to a git
   pre-commit hook locally:

   ```sh
   # .git/hooks/pre-commit (example, not auto-installed)
   #!/bin/sh
   project-docs/tools/workflow-lint.sh
   ```

## When Codex Adds Hook Support

If Codex adds project-level hook configuration, create a
`.codex/hooks.json` or equivalent that mirrors the Claude hook examples:

- Run workflow-lint on session stop or before write operations.
- Keep hooks read-only and non-blocking until the workflow is stable.
- Do not add hooks that mutate code or block work before workflow-lint
  itself is proven reliable.
