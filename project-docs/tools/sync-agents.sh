#!/bin/sh
set -eu

usage() {
  cat <<'EOF'
Usage: project-docs/tools/sync-agents.sh [--check]

Generate provider agent prompts from project-docs/agents/.

Each canonical file project-docs/agents/<role>.md has YAML frontmatter with a
`claude:` block and a `codex:` block, followed by a body made of segments
delimited by marker lines:

  <!-- body: shared -->   emitted to both providers (default before markers)
  <!-- body: claude -->   emitted only to .claude/agents/<name>.md
  <!-- body: codex -->    emitted only to .codex/agents/<name>.toml

Options:
  --check   Compare generated provider prompts with canonical sources without
            writing files.
EOF
}

mode="write"
case "${1:-}" in
  "")
    ;;
  --check)
    mode="check"
    ;;
  -h|--help)
    usage
    exit 0
    ;;
  *)
    usage >&2
    exit 2
    ;;
esac

script_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
repo_root=$(CDPATH= cd "$script_dir/../.." && pwd)
source_root="$repo_root/project-docs/agents"
status=0

# Read one key from a provider block in the canonical frontmatter.
config_value() {
  provider=$1
  key=$2
  file=$3
  awk -v provider="$provider" -v key="$key" '
    NR == 1 && $0 == "---" { in_fm = 1; next }
    in_fm && $0 == "---" { exit }
    in_fm && $0 ~ /^[a-z_]+:$/ { block = substr($0, 1, length($0) - 1); next }
    in_fm && block == provider {
      line = substr($0, 3)
      prefix = key ":"
      if (index(line, prefix) == 1) {
        value = substr(line, length(prefix) + 1)
        gsub(/^[ \t]+|[ \t]+$/, "", value)
        print value
        exit
      }
    }
  ' "$file"
}

generate_agent() {
  provider=$1
  src=$2
  rel=$3
  awk -v provider="$provider" -v source="$rel" '
    function emit_header() {
      if (provider == "claude") {
        print "---"
        for (i = 1; i <= n; i++) print cfg[i]
        print "---"
        print ""
        print "<!-- Generated from " source " by project-docs/tools/sync-agents.sh. -->"
        print "<!-- Do not edit directly; update the canonical agent source and rerun the sync command. -->"
      } else {
        print "# Generated from " source " by project-docs/tools/sync-agents.sh."
        print "# Do not edit directly; update the canonical agent source and rerun the sync command."
        for (i = 1; i <= n; i++) {
          line = cfg[i]
          sep = index(line, ": ")
          if (sep == 0) {
            printf "sync-agents: bad codex config line in %s: %s\n", source, line | "cat >&2"
            exit 1
          }
          key = substr(line, 1, sep - 1)
          value = substr(line, sep + 2)
          print key " = \"" value "\""
        }
        print "developer_instructions = \"\"\""
      }
    }
    BEGIN { state = "start"; target = "shared"; n = 0 }
    state == "start" && NR == 1 && $0 == "---" { state = "fm"; next }
    state == "fm" && $0 == "---" {
      state = "body"
      emit_header()
      next
    }
    state == "fm" && $0 ~ /^[a-z_]+:$/ {
      block = substr($0, 1, length($0) - 1)
      next
    }
    state == "fm" {
      if (block == provider) cfg[++n] = substr($0, 3)
      next
    }
    state == "body" && $0 ~ /^<!-- body: / {
      if ($0 == "<!-- body: shared -->") { target = "shared"; next }
      if ($0 == "<!-- body: claude -->") { target = "claude"; next }
      if ($0 == "<!-- body: codex -->") { target = "codex"; next }
      printf "sync-agents: unknown body marker in %s: %s\n", source, $0 | "cat >&2"
      exit 1
    }
    state == "body" && (target == provider || target == "shared") { print }
    END { if (provider == "codex") print "\"\"\"" }
  ' "$src"
}

if [ ! -d "$source_root" ]; then
  echo "FAIL source directory missing: project-docs/agents" >&2
  exit 1
fi

tmp=${TMPDIR:-/tmp}/sync-agents.$$
list=${TMPDIR:-/tmp}/sync-agents-list.$$
trap 'rm -f "$tmp" "$list"' EXIT HUP INT TERM

find "$source_root" -mindepth 1 -maxdepth 1 -type f -name '*.md' | sort > "$list"

if [ ! -s "$list" ]; then
  echo "FAIL no canonical agent sources under project-docs/agents" >&2
  exit 1
fi

while IFS= read -r src; do
  role=$(basename "$src" .md)
  rel_src="project-docs/agents/$role.md"

  for provider in claude codex; do
    name=$(config_value "$provider" name "$src")
    if [ -z "$name" ]; then
      echo "FAIL missing $provider name in frontmatter: $rel_src" >&2
      exit 1
    fi
    case "$provider" in
      claude) rel_dest=".claude/agents/$name.md" ;;
      codex) rel_dest=".codex/agents/$name.toml" ;;
    esac
    dest="$repo_root/$rel_dest"

    if [ "$mode" = "check" ]; then
      if [ ! -f "$dest" ]; then
        echo "FAIL generated agent prompt missing: $rel_dest" >&2
        status=1
        continue
      fi
      generate_agent "$provider" "$src" "$rel_src" > "$tmp"
      if cmp -s "$tmp" "$dest"; then
        echo "PASS generated agent prompt current: $rel_dest"
      else
        echo "FAIL generated agent prompt drift: $rel_dest" >&2
        status=1
      fi
    else
      generate_agent "$provider" "$src" "$rel_src" > "$dest"
      echo "Generated $rel_dest"
    fi
  done
done < "$list"

exit "$status"
