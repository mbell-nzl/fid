#!/bin/sh
# Verify that one isolated control-conformance commit changed only its allowlist.
set -eu

usage() {
  cat >&2 <<'EOF'
Usage: workflow-control-scope-lint.sh <repo> <trusted-base-commit> <current-head-commit> <allowed-path>...

Allowed paths must be exact repository-relative regular files under
project-docs/tools/. The repair must be checked-out HEAD in a clean worktree,
with the Phase Runner's full canonical pre-dispatch base as its sole parent.
It must change at least one file and must not delete or modify anything outside
the explicit allowlist. The allowlist order is authoritative and is preserved
in the success receipt.
EOF
  exit 2
}

[ "$#" -ge 4 ] || usage
repo=$1
trusted_base=$2
claimed_head=$3
shift 3

git -C "$repo" rev-parse --git-dir >/dev/null 2>&1 || {
  echo "FAIL not a git repository: $repo" >&2
  exit 1
}
trusted_base_commit=$(git -C "$repo" rev-parse --verify --quiet "$trusted_base^{commit}") || {
  echo "FAIL trusted base commit does not resolve: $trusted_base" >&2
  exit 1
}
if [ "$trusted_base" != "$trusted_base_commit" ]; then
  echo "FAIL trusted base must be a full canonical commit ID: $trusted_base" >&2
  exit 1
fi
current_head_commit=$(git -C "$repo" rev-parse --verify --quiet "$claimed_head^{commit}") || {
  echo "FAIL current head commit does not resolve: $claimed_head" >&2
  exit 1
}
if [ "$claimed_head" != "$current_head_commit" ]; then
  echo "FAIL current head must be a full canonical commit ID: $claimed_head" >&2
  exit 1
fi
checked_out_head=$(git -C "$repo" rev-parse HEAD)
if [ "$checked_out_head" != "$current_head_commit" ]; then
  echo "FAIL repair commit is not the checked-out HEAD: claimed $current_head_commit, actual $checked_out_head" >&2
  exit 1
fi
if ! current_parent=$(git -C "$repo" rev-parse --verify --quiet "$current_head_commit^"); then
  echo "FAIL current HEAD does not have the trusted base as its sole parent" >&2
  exit 1
fi
if git -C "$repo" rev-parse --verify --quiet "$current_head_commit^2" >/dev/null; then
  echo "FAIL current HEAD does not have the trusted base as its sole parent" >&2
  exit 1
fi
if [ "$current_parent" != "$trusted_base_commit" ]; then
  echo "FAIL trusted base is not the sole parent of current HEAD: trusted $trusted_base_commit, parent $current_parent" >&2
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
  echo "FAIL control repair worktree is not clean at $current_head_commit" >&2
  printf '%s\n' "$dirty" | sed 's/^/  /' >&2
  exit 1
fi
git -C "$repo" merge-base --is-ancestor "$trusted_base_commit" "$current_head_commit" || {
  echo "FAIL repair commit does not descend from base commit" >&2
  exit 1
}
commit_count=$(git -C "$repo" rev-list --count "$trusted_base_commit..$current_head_commit")
if [ "$commit_count" -ne 1 ]; then
  echo "FAIL control-conformance repair must be exactly one commit after its base: found $commit_count" >&2
  exit 1
fi

allowed=""
allowed_receipt=""
for path in "$@"; do
  case "$path" in
    *'
'*)
      echo "FAIL allowed path contains LF, CR, or control bytes" >&2
      exit 1
      ;;
  esac
  if printf '%s' "$path" | LC_ALL=C grep -q '[[:cntrl:]]'; then
    echo "FAIL allowed path contains LF, CR, or control bytes" >&2
    exit 1
  fi
  case "$path" in
    *,*)
      echo "FAIL allowed path contains comma: $path" >&2
      exit 1
      ;;
  esac
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
  case "$path" in
    project-docs/tools/workflow-control-scope-lint.sh|project-docs/tools/workflow-sequence-lint.sh)
      echo "FAIL trusted verifier cannot be allowlisted: $path" >&2
      exit 1
      ;;
  esac
  if printf '%s\n' "$allowed" | grep -Fqx -- "$path"; then
    echo "FAIL duplicate allowed path: $path" >&2
    exit 1
  fi
  allowed="${allowed}${path}
"
  if [ -z "$allowed_receipt" ]; then
    allowed_receipt=$path
  else
    allowed_receipt="${allowed_receipt},${path}"
  fi
done

changed=$(git -C "$repo" diff --name-only --no-renames "$trusted_base_commit" "$current_head_commit")
if [ -z "$changed" ]; then
  echo "FAIL control repair commit changes no files" >&2
  exit 1
fi

deleted=$(git -C "$repo" diff --name-only --no-renames --diff-filter=D "$trusted_base_commit" "$current_head_commit")
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
  entry=$(git -C "$repo" ls-tree "$current_head_commit" -- "$path")
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

printf 'PASS control repair scope: %s -> %s; allowed=%s\n' \
  "$trusted_base_commit" "$current_head_commit" "$allowed_receipt"
