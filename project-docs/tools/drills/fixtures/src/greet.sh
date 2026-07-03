#!/bin/sh
# Drill fixture product: prints a greeting; --shout uppercases it.
name=""
shout=false
for arg in "$@"; do
  case "$arg" in
    --shout) shout=true ;;
    *) name=$arg ;;
  esac
done
[ -n "$name" ] || { echo "usage: greet.sh [--shout] <name>" >&2; exit 2; }
greeting="hello, $name"
if $shout; then
  greeting=$(printf '%s' "$greeting" | tr '[:lower:]' '[:upper:]')
fi
printf '%s\n' "$greeting"
