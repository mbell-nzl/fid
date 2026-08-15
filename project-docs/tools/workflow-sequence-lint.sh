#!/bin/sh
# workflow-sequence-lint.sh - validate recovery ordering in a v1 run ledger.
set -eu

if [ "$#" -lt 1 ] || [ "$#" -gt 2 ] || [ ! -f "$1" ]; then
  echo "Usage: workflow-sequence-lint.sh <run-ledger.md> [--append-stdin]" >&2
  exit 2
fi

input=$1
tmp=""
if [ "$#" -eq 2 ]; then
  [ "$2" = "--append-stdin" ] || {
    echo "Usage: workflow-sequence-lint.sh <run-ledger.md> [--append-stdin]" >&2
    exit 2
  }
  tmp=$(mktemp -d "${TMPDIR:-/tmp}/workflow-sequence.XXXXXX")
  trap 'rm -rf "$tmp"' EXIT HUP INT TERM
  { cat "$1"; cat; } > "$tmp/combined-ledger.md"
  input="$tmp/combined-ledger.md"
fi

awk '
function invalid(message) {
  print "INVALID " message > "/dev/stderr"
  errors++
}

function finish_event() {
  if (strict && event == "") {
    invalid("unrecognized or malformed event heading")
  }
  if (strict && event == "REVIEW" && event_outcome !~ /^(PASS|REPAIR|BLOCK)$/) {
    invalid("REVIEW has missing or invalid Outcome: " event_id)
  }
  if (strict && event == "REVIEW" && event_id !~ /^(review-[0-9]+[a-z]*|infra-review-[0-9]+[a-z]*)$/) {
    invalid("REVIEW has missing or invalid review ID")
  }
  if (strict && event == "REVIEW" && (event_owner == "" || event_rereview == "")) {
    invalid("REVIEW is missing Owning task or Re-review of: " event_id)
  }
  if (strict && event == "TRIAGE" && event_verdict !~ /^(RESUME|REPAIR|BLOCK)/) {
    invalid("TRIAGE has missing or invalid Verdict")
  }
  if (strict && (event == "TASK-RECOVERY" || event == "REPAIR" || event == "TASK-COMPLETE") && event_task !~ /^task-[0-9]+[a-z]*$/) {
    invalid(event " has missing or invalid task ID")
  }
  if (strict && (event == "RUN-START" || event == "RUN-RESUME") && event_contract != "v2") {
    invalid(event " is missing Sequence contract: v2")
  }

  if (event == "TASK-RECOVERY") {
    if (review_open && event_task != review_owner) {
      invalid("QA repair owner " review_owner " does not match recovery task " event_task)
    }
    if (recovery_task != "" && recovery_task != event_task) {
      invalid("overlapping recovery changes task identity from " recovery_task " to " event_task)
    } else {
      recovery_task = event_task
    }
    recovery_stage = "needs-repair"
    triage_repair_pending = 0
    if (review_open && event_task == review_owner) {
      review_recovery_seen = 1
    }
  } else if (event == "REPAIR") {
    if (recovery_task == "") {
      invalid("REPAIR has no preceding TASK-RECOVERY")
    } else if (event_task != recovery_task) {
      invalid("REPAIR task " event_task " does not match recovery task " recovery_task)
    } else {
      recovery_stage = "needs-complete"
      if (review_open && event_task == review_owner) {
        review_repair_seen = 1
      }
    }
  } else if (event == "TASK-COMPLETE") {
    if (recovery_task != "") {
      if (event_task != recovery_task) {
        invalid("completed task " event_task " does not match recovery task " recovery_task)
      } else if (recovery_stage != "needs-complete") {
        invalid("task completed before a REPAIR event: " event_task)
      } else {
        recovery_task = ""
        recovery_stage = ""
      }
    }
    if (review_open && event_task == review_owner) {
      if (!review_recovery_seen || !review_repair_seen) {
        invalid("QA repair owner completed without TASK-RECOVERY and REPAIR: " event_task)
      } else {
        review_task_complete = 1
      }
    }
  } else if (event == "REVIEW") {
    if (event_outcome == "REPAIR") {
      if (event_owner == "" || event_owner == "none") {
        invalid("REVIEW REPAIR has no owning task: " event_id)
      }
      review_open = 1
      review_id = event_id
      review_owner = event_owner
      review_task_complete = 0
      review_recovery_seen = 0
      review_repair_seen = 0
    } else if (event_outcome == "PASS" && review_open) {
      if (!review_task_complete) {
        invalid("QA re-review passed before owning task completed: " review_owner)
      }
      if (event_id == review_id) {
        invalid("QA repair was self-closed by the same review identity: " event_id)
      }
      if (event_rereview != review_id) {
        invalid("QA PASS does not reference required re-review " review_id ": " event_id)
      }
      review_open = 0
      review_id = ""
      review_owner = ""
      review_task_complete = 0
      review_recovery_seen = 0
      review_repair_seen = 0
    } else if (event_outcome == "BLOCK") {
      triage_needed = 1
    }
  } else if (event == "TRIAGE") {
    if (event_verdict ~ /^RESUME/) {
      triage_needed = 0
    } else if (event_verdict ~ /^REPAIR/) {
      triage_needed = 0
      triage_repair_pending = 1
    } else if (event_verdict ~ /^BLOCK/) {
      stopped = 1
      triage_needed = 0
    }
  } else if (event == "STOP") {
    stopped = 1
  } else if (event == "RUN-RESUME") {
    stopped = 0
    if (event_contract == "v2") {
      strict = 1
    }
  } else if (event == "RUN-START") {
    if (event_contract == "v2") {
      strict = 1
    }
  }

  event = ""
  event_task = ""
  event_id = ""
  event_outcome = ""
  event_owner = ""
  event_rereview = ""
  event_verdict = ""
  event_contract = ""
}

/^### / {
  finish_event()
  if ($4 == "TASK-RECOVERY") {
    event = "TASK-RECOVERY"
    event_task = $5
  } else if ($4 == "REPAIR") {
    event = "REPAIR"
    event_task = $5
  } else if ($4 == "TASK" && $6 == "COMPLETE") {
    event = "TASK-COMPLETE"
    event_task = $5
  } else if ($4 == "REVIEW") {
    event = "REVIEW"
    event_id = $5
  } else if ($4 == "TRIAGE") {
    event = "TRIAGE"
  } else if ($4 == "STOP") {
    event = "STOP"
  } else if ($4 == "RUN-RESUME") {
    event = "RUN-RESUME"
  } else if ($4 == "RUN-START") {
    event = "RUN-START"
  } else if ($4 == "RUN-END") {
    event = "RUN-END"
  } else if ($4 == "CORRECTION") {
    event = "CORRECTION"
  } else {
    event = ""
  }
  next
}

event == "REVIEW" && /^- Outcome: / {
  event_outcome = substr($0, length("- Outcome: ") + 1)
}
event == "REVIEW" && /^- Owning task: / {
  event_owner = substr($0, length("- Owning task: ") + 1)
}
event == "REVIEW" && /^- Re-review of: / {
  event_rereview = substr($0, length("- Re-review of: ") + 1)
}
event == "TRIAGE" && /^- Verdict: / {
  event_verdict = substr($0, length("- Verdict: ") + 1)
}
(event == "RUN-START" || event == "RUN-RESUME") && /^- Sequence contract: / {
  event_contract = substr($0, length("- Sequence contract: ") + 1)
}

END {
  finish_event()
  if (errors) {
    exit 1
  }
  if (stopped) {
    print "VALID state=STOPPED"
  } else if (triage_needed) {
    print "VALID state=TRIAGE_PENDING"
  } else if (triage_repair_pending) {
    print "VALID state=REPAIR_PENDING"
  } else if (recovery_task != "") {
    if (recovery_stage == "needs-repair") {
      print "VALID state=REPAIR_PENDING"
    } else {
      print "VALID state=COMPLETION_PENDING"
    }
  } else if (review_open) {
    if (review_task_complete) {
      print "VALID state=REREVIEW_PENDING"
    } else {
      print "VALID state=REPAIR_PENDING"
    }
  } else {
    print "VALID state=ADVANCEABLE"
  }
}
' "$input"
