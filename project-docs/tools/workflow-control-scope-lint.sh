#!/bin/sh
# Verify that one isolated control-conformance commit changed only its allowlist.
set -eu

usage() {
  cat >&2 <<'EOF'
Usage: workflow-control-scope-lint.sh <repo> <base-commit> <repair-commit> <allowed-path>...

Allowed paths must be exact repository-relative regular files under
project-docs/tools/. The repair must be checked-out HEAD in a clean worktree,
exactly one commit after the base, change at least one file, and must not delete
or modify anything outside the explicit allowlist.
EOF
  exit 2
}

[ "$#" -ge 4 ] || usage
repo=$1
base=$2
head=$3
shift 3

git -C "$repo" rev-parse --git-dir >/dev/null 2>&1 || {
  echo "FAIL not a git repository: $repo" >&2
  exit 1
}
git -C "$repo" rev-parse --verify --quiet "$base^{commit}" >/dev/null || {
  echo "FAIL base commit does not resolve: $base" >&2
  exit 1
}
git -C "$repo" rev-parse --verify --quiet "$head^{commit}" >/dev/null || {
  echo "FAIL repair commit does not resolve: $head" >&2
  exit 1
}
head_commit=$(git -C "$repo" rev-parse "$head^{commit}")
checked_out_head=$(git -C "$repo" rev-parse HEAD)
if [ "$checked_out_head" != "$head_commit" ]; then
  echo "FAIL repair commit is not the checked-out HEAD: claimed $head_commit, actual $checked_out_head" >&2
  exit 1
fi
special_index=$(git -C "$repo" ls-files -v | awk '
  substr($0, 1, 1) == "S" || substr($0, 1, 1) ~ /[a-z]/ { print }
')
if [ -n "$special_index" ]; then
  echo "FAIL control repair index contains assume-unchanged or skip-worktree paths" >&2
  printf '%s\n' "$special_index" | sed 's/^/  /' >&2
  exit 1
fi
git -C "$repo" update-index --really-refresh -q >/dev/null 2>&1 || true
dirty=$(git -C "$repo" status --porcelain=v1 --untracked-files=all)
if [ -n "$dirty" ]; then
  echo "FAIL control repair worktree is not clean at $head_commit" >&2
  printf '%s\n' "$dirty" | sed 's/^/  /' >&2
  exit 1
fi
git -C "$repo" merge-base --is-ancestor "$base" "$head" || {
  echo "FAIL repair commit does not descend from base commit" >&2
  exit 1
}
commit_count=$(git -C "$repo" rev-list --count "$base..$head")
if [ "$commit_count" -ne 1 ]; then
  echo "FAIL control-conformance repair must be exactly one commit after its base: found $commit_count" >&2
  exit 1
fi

allowed=""
for path in "$@"; do
  case "$path" in
    project-docs/tools/*)
      case "/$path/" in
        */../*|*/./*)
          echo "FAIL allowed path is not normalized: $path" >&2
          exit 1
          ;;
      esac
      ;;
    *)
      echo "FAIL allowed path is not a workflow control: $path" >&2
      exit 1
      ;;
  esac
  if printf '%s\n' "$allowed" | grep -Fqx -- "$path"; then
    echo "FAIL duplicate allowed path: $path" >&2
    exit 1
  fi
  allowed="${allowed}${path}
"
done

changed=$(git -C "$repo" diff --name-only --no-renames "$base" "$head")
if [ -z "$changed" ]; then
  echo "FAIL control repair commit changes no files" >&2
  exit 1
fi

deleted=$(git -C "$repo" diff --name-only --no-renames --diff-filter=D "$base" "$head")
if [ -n "$deleted" ]; then
  while IFS= read -r path; do
    [ -n "$path" ] || continue
    echo "FAIL control-conformance repair deletes file: $path" >&2
  done <<EOF
$deleted
EOF
  exit 1
fi

failures=0
while IFS= read -r path; do
  [ -n "$path" ] || continue
  if ! printf '%s\n' "$allowed" | grep -Fqx -- "$path"; then
    echo "FAIL changed path is outside allowed control scope: $path" >&2
    failures=$((failures + 1))
    continue
  fi
  entry=$(git -C "$repo" ls-tree "$head_commit" -- "$path")
  mode=$(printf '%s\n' "$entry" | awk '{ print $1 }')
  type=$(printf '%s\n' "$entry" | awk '{ print $2 }')
  case "$mode:$type" in
    100644:blob|100755:blob) ;;
    *)
      echo "FAIL allowed control is not a regular file at repair commit: $path ($mode:$type)" >&2
      failures=$((failures + 1))
      ;;
  esac
done <<EOF
$changed
EOF

if [ "$failures" -ne 0 ]; then
  exit 1
fi

printf 'PASS control repair scope: %s -> %s\n' "$base" "$head"
