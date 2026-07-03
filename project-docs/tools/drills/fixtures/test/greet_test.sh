#!/bin/sh
# Drill fixture tests: cover both VAL-001 states.
fails=0
assert_eq() {
  desc=$1
  want=$2
  got=$3
  if [ "$want" = "$got" ]; then
    echo "ok - $desc"
  else
    echo "not ok - $desc (want '$want', got '$got')"
    fails=$((fails + 1))
  fi
}
dir=$(CDPATH= cd "$(dirname "$0")/.." && pwd)
assert_eq "default greeting stays lowercase" "hello, sam" "$(sh "$dir/src/greet.sh" sam)"
assert_eq "--shout uppercases the greeting" "HELLO, SAM" "$(sh "$dir/src/greet.sh" --shout sam)"
exit "$fails"
