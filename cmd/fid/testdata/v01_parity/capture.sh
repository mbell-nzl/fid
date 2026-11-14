#!/bin/sh
set -eu

if [ "$#" -ne 2 ]; then
	printf '%s\n' "usage: $0 <fid-binary> <output-root>" >&2
	exit 2
fi

binary=$(realpath "$1")
output_root=$2
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
case $output_root in
	/*) ;;
	*) output_root=$(pwd)/$output_root ;;
esac
workspace=$(mktemp -d)
trap 'rm -rf "$workspace"' EXIT

export GIT_CEILING_DIRECTORIES=$(dirname "$workspace")
git -C "$workspace" init -q --template=
mkdir -p "$output_root"
cp -R "$script_dir/fixture/." "$workspace/"
mv "$workspace/.fid" "$workspace/.git/fid"

(
	cd "$workspace"
	"$binary" issue get FID-a10001 --json > "$output_root/issue-get-epic.json"
	"$binary" issue get FID-b20002 --json > "$output_root/issue-get-task.json"
	"$binary" issue get FID-c30003 --json > "$output_root/issue-get-subtask.json"
	"$binary" issue list --kind epic --json > "$output_root/issue-list-epic.json"
	"$binary" issue list --kind task --json > "$output_root/issue-list-task.json"
	"$binary" issue list --kind subtask --json > "$output_root/issue-list-subtask.json"
)
