#!/bin/sh
# workflow-route.sh - deterministic routing for v1 workflow outcomes.
set -eu

usage() {
  cat >&2 <<'EOF'
Usage: workflow-route.sh <source> <outcome> [consecutive-no-progress-attempts]

Sources and outcomes:
  developer   COMPLETE | RECOVERY_REQUIRED | BLOCKED
  report-lint COMPLETE | RECOVERY_REQUIRED | BLOCKED | INVALID
  qa          PASS | REPAIR | BLOCK
  triage      RESUME | REPAIR | BLOCK
  command     LIVE | TERMINATED | INTERRUPTED | MISSING_TERMINATION
  recovery    PROGRESS | NO_PROGRESS (requires attempt count for NO_PROGRESS)
  control     CONFORMANCE_DEFECT | DEFINITION_CHANGE | AMBIGUOUS

Prints exactly one ACTION line. Exit 2 means invalid input.
EOF
  exit 2
}

[ "$#" -ge 2 ] && [ "$#" -le 3 ] || usage
source=$1
outcome=$2
attempts=${3:-}

action=""
case "$source:$outcome" in
  developer:COMPLETE|report-lint:COMPLETE|qa:PASS|triage:RESUME)
    action=ADVANCE
    ;;
  developer:RECOVERY_REQUIRED|report-lint:RECOVERY_REQUIRED|report-lint:INVALID|qa:REPAIR|triage:REPAIR)
    action=REPAIR
    ;;
  developer:BLOCKED|report-lint:BLOCKED|qa:BLOCK)
    action=TRIAGE
    ;;
  triage:BLOCK)
    action=STOP
    ;;
  command:LIVE)
    action=WAIT
    ;;
  command:TERMINATED)
    action=EVALUATE
    ;;
  command:INTERRUPTED|command:MISSING_TERMINATION|recovery:PROGRESS)
    action=REPAIR
    ;;
  control:CONFORMANCE_DEFECT)
    action=REPAIR
    ;;
  control:DEFINITION_CHANGE|control:AMBIGUOUS)
    action=TRIAGE
    ;;
  recovery:NO_PROGRESS)
    case "$attempts" in
      ''|*[!0-9]*) usage ;;
    esac
    if [ "$attempts" -ge 2 ]; then
      action=TRIAGE
    else
      action=REPAIR
    fi
    ;;
  *)
    usage
    ;;
esac

if [ "$source" != recovery ] && [ -n "$attempts" ]; then
  usage
fi

printf 'ACTION %s\n' "$action"
