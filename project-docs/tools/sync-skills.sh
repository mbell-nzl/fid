#!/bin/sh
set -eu

usage() {
  cat <<'EOF'
Usage: project-docs/tools/sync-skills.sh [--check]

Generate provider skill copies from project-docs/skills/.

Options:
  --check   Compare generated provider copies with canonical skills without
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
source_root="$repo_root/project-docs/skills"
providers=".agents/skills .claude/skills"
status=0

generate_skill() {
  src=$1
  rel=$2
  awk -v source="$rel" '
    BEGIN { in_frontmatter = 0 }
    NR == 1 && $0 == "---" {
      in_frontmatter = 1
      print
      next
    }
    in_frontmatter && $0 == "---" {
      print
      print ""
      print "<!-- Generated from " source " by project-docs/tools/sync-skills.sh. -->"
      print "<!-- Do not edit directly; update the canonical skill and rerun the sync command. -->"
      in_frontmatter = 0
      next
    }
    { print }
  ' "$src"
}

if [ ! -d "$source_root" ]; then
  echo "FAIL source directory missing: project-docs/skills" >&2
  exit 1
fi

tmp=${TMPDIR:-/tmp}/sync-skills.$$
list=${TMPDIR:-/tmp}/sync-skills-list.$$
trap 'rm -f "$tmp" "$list"' EXIT HUP INT TERM

find "$source_root" -mindepth 1 -maxdepth 1 -type d | sort > "$list"

while IFS= read -r skill_dir; do
  skill_name=$(basename "$skill_dir")
  src="$skill_dir/SKILL.md"
  rel_src="project-docs/skills/$skill_name/SKILL.md"

  if [ ! -f "$src" ]; then
    echo "FAIL missing canonical skill file: $rel_src" >&2
    exit 1
  fi

  for provider in $providers; do
    dest_dir="$repo_root/$provider/$skill_name"
    dest="$dest_dir/SKILL.md"
    rel_dest="$provider/$skill_name/SKILL.md"

    if [ "$mode" = "check" ]; then
      if [ ! -f "$dest" ]; then
        echo "FAIL generated skill missing: $rel_dest" >&2
        status=1
        continue
      fi
      generate_skill "$src" "$rel_src" > "$tmp"
      if cmp -s "$tmp" "$dest"; then
        echo "PASS generated skill current: $rel_dest"
      else
        echo "FAIL generated skill drift: $rel_dest" >&2
        status=1
      fi
    else
      mkdir -p "$dest_dir"
      generate_skill "$src" "$rel_src" > "$dest"
      echo "Generated $rel_dest"
    fi
  done
done < "$list"

exit "$status"
