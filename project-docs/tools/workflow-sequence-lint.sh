#!/bin/sh
# workflow-sequence-lint.sh - validate or transactionally append a ledger event.
set -eu

if [ "$#" -lt 1 ] || [ "$#" -gt 2 ] || [ ! -f "$1" ]; then
  echo "Usage: workflow-sequence-lint.sh <run-ledger.md> [--append-stdin]" >&2
  exit 2
fi

ledger=$1
mode=validate
if [ "$#" -eq 2 ]; then
  [ "$2" = "--append-stdin" ] || {
    echo "Usage: workflow-sequence-lint.sh <run-ledger.md> [--append-stdin]" >&2
    exit 2
  }
  mode=append
fi

validate_sequence() {
ledger_path=$1
phase_dir=$(dirname "$ledger_path")
phase_name=$(basename "$phase_dir")
phase_number=""
decisions_file=""
triage_file=""
case "$phase_name" in
  phase-*)
    candidate_phase=${phase_name#phase-}
    case "$candidate_phase" in
      ''|*[!0-9]*) ;;
      *)
        phase_number=$candidate_phase
        decisions_file="$phase_dir/decisions.md"
        triage_file="$phase_dir/triage-log.md"
        ;;
    esac
    ;;
esac

awk -v phase_number="$phase_number" -v decisions_file="$decisions_file" \
    -v triage_file="$triage_file" '
function invalid(message) {
  print "INVALID " message > "/dev/stderr"
  errors++
}

function references_id(text, id) {
  return id != "" && text ~ ("(^|[^[:alnum:]_-])" id "([^[:alnum:]_-]|$)")
}

function decision_id(text) {
  if (match(text, /DEC-[0-9]+/)) {
    return substr(text, RSTART, RLENGTH)
  }
  return ""
}

function escalation_id(text) {
  if (match(text, /ESC-[0-9]+-[0-9]+/)) {
    return substr(text, RSTART, RLENGTH)
  }
  return ""
}

function decision_exists(id, line) {
  if (id == "" || decisions_file == "") {
    return 0
  }
  while ((getline line < decisions_file) > 0) {
    if (index(line, "### " id ":") == 1) {
      close(decisions_file)
      return 1
    }
  }
  close(decisions_file)
  return 0
}

function triage_exists(id, line, prefix, following) {
  if (id == "" || triage_file == "") {
    return 0
  }
  prefix = "## " id
  while ((getline line < triage_file) > 0) {
    if (index(line, prefix) == 1) {
      following = substr(line, length(prefix) + 1, 1)
      if (following == "" || following == " " || following == "-") {
        close(triage_file)
        return 1
      }
    }
  }
  close(triage_file)
  return 0
}

function triage_verdict_matches(id, verdict, line, prefix, following, in_section,
                                marker, value, esc, found_verdict, found_esc) {
  if (id == "" || triage_file == "") {
    return 0
  }
  prefix = "## " id
  esc = escalation_id(verdict)
  found_esc = (esc == "")
  while ((getline line < triage_file) > 0) {
    if (index(line, "## TRG-") == 1) {
      if (in_section) {
        break
      }
      if (index(line, prefix) == 1) {
        following = substr(line, length(prefix) + 1, 1)
        in_section = (following == "" || following == " " || following == "-")
      }
      continue
    }
    if (!in_section) {
      continue
    }
    if (!found_esc && references_id(line, esc)) {
      found_esc = 1
    }
    marker = index(line, "**Final verdict:**")
    if (marker) {
      value = substr(line, marker + length("**Final verdict:**"))
    } else {
      marker = index(line, "**Verdict:**")
      if (marker) {
        value = substr(line, marker + length("**Verdict:**"))
      }
    }
    if (marker) {
      sub(/^[[:space:]]+/, "", value)
      if (verdict == "REPAIR / control-conformance") {
        found_verdict = (index(value, verdict) == 1)
      } else if (verdict == "REPAIR") {
        found_verdict = (index(value, "REPAIR") == 1 &&
                         index(value, "REPAIR / control-conformance") != 1)
      } else if (verdict == "RESUME") {
        found_verdict = (index(value, "RESUME") == 1)
      } else if (verdict ~ /^BLOCK/) {
        found_verdict = (index(value, "BLOCK") == 1)
      }
    }
  }
  close(triage_file)
  return found_verdict && found_esc
}

function finish_event() {
  if (strict && event == "") {
    invalid("unrecognized or malformed event heading")
  }
  if (stop_entry_pending && event != "" && event != "STOP") {
    invalid("TRIAGE BLOCK must be followed by STOP before " event)
  }
  if (recovery_stage == "needs-triage" && event != "" && event != "TRIAGE") {
    invalid("recovery no-progress budget requires TRIAGE before " event)
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
  if ((strict || event_contract != "" || event_directive_schema != "") && (event == "RUN-START" || event == "RUN-RESUME") && event_contract != "v2") {
    invalid(event " is missing Sequence contract: v2")
  }
  if ((strict || event_contract == "v2" || event_directive_schema != "") && (event == "RUN-START" || event == "RUN-RESUME") && event_directive_schema != "" && event_directive_schema != "v1") {
    invalid(event " has invalid Recovery directive schema")
  }
  if (directive_strict && (event == "RUN-START" || event == "RUN-RESUME") && event_directive_schema != "v1") {
    invalid(event " is missing Recovery directive schema: v1")
  }
  if (strict && directive_strict && event == "TASK-RECOVERY") {
    if (event_source == "") {
      invalid("TASK-RECOVERY is missing Source")
    }
    if (event_category !~ /^(in-scope-implementation|gate-repair|report-repair|qa-rework|session-continuation|local-environment)$/) {
      invalid("TASK-RECOVERY has missing or invalid Category")
    }
    if (event_repair_mode !~ /^(task-scope|control-conformance)$/) {
      invalid("TASK-RECOVERY has missing or invalid Repair mode")
    }
    if (event_repair_authority == "") {
      invalid("TASK-RECOVERY is missing Repair authority")
    }
    if (event_allowed_files == "") {
      invalid("TASK-RECOVERY is missing Allowed files")
    }
    if (event_forbidden_changes == "") {
      invalid("TASK-RECOVERY is missing Forbidden changes")
    }
    if (event_trigger_signature == "") {
      invalid("TASK-RECOVERY is missing Trigger signature")
    }
    if (event_attempt !~ /^[1-9][0-9]*$/) {
      invalid("TASK-RECOVERY has missing or invalid Attempt count")
    }
    if (event_no_progress !~ /^[01]\/2$/) {
      invalid("TASK-RECOVERY has missing or invalid Consecutive no-progress attempts")
    }
    if (event_checkpoint == "") {
      invalid("TASK-RECOVERY is missing Checkpoint")
    }
    if (event_resolution_evidence == "") {
      invalid("TASK-RECOVERY is missing Resolution evidence")
    }
    if (event_next_action == "") {
      invalid("TASK-RECOVERY is missing Next action")
    }
    if (event_human_attention != "no") {
      invalid("TASK-RECOVERY must record Human attention required: no")
    }
    if (recovery_task != "" && recovery_stage == "needs-repair") {
      invalid("overlapping TASK-RECOVERY precedes required REPAIR")
    } else if (recovery_task != "" && recovery_stage == "needs-continuation") {
      if (event_task != recovery_task || event_trigger_signature != recovery_trigger) {
        invalid("continuing recovery must preserve task and Trigger signature")
      }
      if (event_attempt !~ /^[1-9][0-9]*$/ || event_attempt + 0 != recovery_attempt + 1) {
        invalid("continuing recovery Attempt count must increment by one")
      }
      if (recovery_last_outcome == "NO-PROGRESS") {
        expected_no_progress = (substr(recovery_no_progress, 1, 1) + 1) "/2"
      } else {
        expected_no_progress = "0/2"
      }
      if (event_no_progress != expected_no_progress) {
        invalid("continuing recovery has invalid Consecutive no-progress attempts; expected " expected_no_progress)
      }
    } else if (recovery_task != "" && recovery_stage == "needs-triage") {
      invalid("recovery budget is exhausted; TASK-RECOVERY cannot precede TRIAGE")
    } else if (recovery_task != "" && recovery_stage == "needs-complete") {
      if (recovery_mode != "control-conformance" || event_repair_mode != "task-scope" ||
          event_task != recovery_task || event_trigger_signature != recovery_trigger ||
          event_attempt !~ /^[1-9][0-9]*$/ || event_attempt + 0 != recovery_attempt + 1 ||
          event_no_progress != "0/2") {
        invalid("resolved recovery permits only its linked ordinary report continuation")
      }
    }
  }
  if (strict && directive_strict && event == "TASK-RECOVERY" && event_repair_mode == "control-conformance") {
    if (event_control_contract == "") {
      invalid("control-conformance TASK-RECOVERY is missing Control contract")
    }
    if (event_control_mismatch == "") {
      invalid("control-conformance TASK-RECOVERY is missing Control mismatch")
    }
    if (event_allowed_files == "") {
      invalid("control-conformance TASK-RECOVERY is missing Allowed files")
    }
    if (event_forbidden_changes == "") {
      invalid("control-conformance TASK-RECOVERY is missing Forbidden changes")
    }
    if (latest_triage_id == "") {
      invalid("control-conformance TASK-RECOVERY has no preceding TRIAGE authority")
    } else {
      if (!references_id(event_repair_authority, latest_triage_id)) {
        invalid("control-conformance Repair authority does not reference latest TRIAGE " latest_triage_id)
      }
      if (event_trigger_signature != latest_triage_trigger) {
        invalid("control-conformance Trigger signature does not match latest TRIAGE " latest_triage_id)
      }
      if (index(event_trigger_signature, event_task ":") != 1) {
        invalid("control-conformance Trigger signature does not bind originating task " event_task)
      }
      if (latest_triage_verdict == "REPAIR / control-conformance") {
        # TL-Lite directly authorized this bounded repair.
      } else if (latest_triage_verdict ~ /^BLOCK/) {
        authority_decision = decision_id(event_repair_authority)
        if (resolved_block_triage != latest_triage_id || authority_decision == "" ||
            !references_id(resolved_block_evidence, authority_decision) ||
            !decision_exists(authority_decision)) {
          invalid("control-conformance BLOCK authority lacks a bound RUN-RESUME decision for " latest_triage_id)
        }
      } else {
        invalid("control-conformance latest TRIAGE is not REPAIR / control-conformance or resolved BLOCK: " latest_triage_id)
      }
      if (event_attempt ~ /^[1-9][0-9]*$/ && latest_triage_attempt ~ /^[1-9][0-9]*$/ &&
          (event_attempt + 0 < latest_triage_attempt + 0 || event_attempt + 0 > latest_triage_attempt + 1)) {
        invalid("control-conformance Attempt count exceeds latest TRIAGE authority " latest_triage_id)
      }
    }
    if (event_next_action !~ /workflow-control-scope-lint\.sh/) {
      invalid("control-conformance TASK-RECOVERY Next action is missing exact scope verification")
    }
  }
  if (strict && directive_strict && event == "REPAIR") {
    if (event_trigger_signature == "" || event_trigger_signature != recovery_trigger) {
      invalid("REPAIR Trigger signature does not match TASK-RECOVERY")
    }
    if (event_attempt !~ /^[1-9][0-9]*$/ || event_attempt != recovery_attempt) {
      invalid("REPAIR Attempt count does not match TASK-RECOVERY")
    }
    if (event_control_scope_check == "") {
      invalid("REPAIR is missing Control scope check")
    }
    if (event_outcome !~ /^(RESOLVED|MATERIAL-PROGRESS|NO-PROGRESS)$/) {
      invalid("REPAIR has missing or invalid Outcome")
    }
    if (event_resolution_evidence == "") {
      invalid("REPAIR is missing Resolution evidence")
    }
    if (event_next_checkpoint == "") {
      invalid("REPAIR is missing Next checkpoint")
    }
    if (recovery_mode == "control-conformance" &&
        (event_control_scope_check !~ /^PASS/ || event_control_scope_check !~ /exit_code: 0/)) {
      invalid("control-conformance REPAIR lacks a passing terminated scope check")
    }
  }
  if (strict && directive_strict && event == "TRIAGE") {
    if (event_id !~ /^TRG-[0-9]+-[0-9]+$/) {
      invalid("TRIAGE has missing or invalid TRG ID")
    } else if (event_id in triage_seen) {
      invalid("TRIAGE reuses an existing TRG ID: " event_id)
    } else if (phase_number != "" && index(event_id, "TRG-" phase_number "-") != 1) {
      invalid("TRIAGE ID does not match ledger phase-" phase_number ": " event_id)
    }
    if (phase_number == "") {
      invalid("strict TRIAGE ledger is not under a phase-N directory")
    } else if (!triage_exists(event_id)) {
      invalid("TRIAGE ID is not recorded in the phase triage log: " event_id)
    } else if (!triage_verdict_matches(event_id, event_verdict)) {
      invalid("TRIAGE verdict does not match the phase triage log: " event_id)
    }
    if (event_verdict !~ /^(RESUME|REPAIR( \/ control-conformance)?|BLOCK \(escalated as ESC-[0-9]+-[0-9]+\))$/) {
      invalid("TRIAGE has missing or invalid strict Verdict")
    }
    triage_esc = escalation_id(event_verdict)
    if (triage_esc != "" && phase_number != "" && index(triage_esc, "ESC-" phase_number "-") != 1) {
      invalid("TRIAGE ESC does not match ledger phase-" phase_number ": " triage_esc)
    }
    if (event_trigger_signature == "") {
      invalid("TRIAGE is missing Trigger signature")
    }
    if (event_attempt !~ /^[1-9][0-9]*$/) {
      invalid("TRIAGE has missing or invalid Attempt count")
    }
    if (event_checkpoint == "") {
      invalid("TRIAGE is missing Checkpoint")
    }
    if (event_resolution_evidence == "") {
      invalid("TRIAGE is missing Resolution evidence")
    }
    if (event_triage_log == "") {
      invalid("TRIAGE is missing Triage log")
    }
    if (recovery_stage == "needs-triage" && event_verdict !~ /^BLOCK/) {
      invalid("exhausted no-progress recovery TRIAGE must BLOCK")
    }
  }
  if (strict && directive_strict && event == "RUN-RESUME" && stopped) {
    if (event_prior_run == "" || event_trigger_signature == "" ||
        event_attempt !~ /^[1-9][0-9]*$/ || event_checkpoint == "" ||
        event_resolution_evidence == "" || event_continuing_task !~ /^task-[0-9]+[a-z]*$/) {
      invalid("RUN-RESUME after STOP is missing reconciliation authority or continuing task")
    } else {
      if (!references_id(event_prior_run, latest_triage_id) ||
          latest_triage_esc == "" || !references_id(event_prior_run, latest_triage_esc)) {
        invalid("RUN-RESUME does not reference its governing TRIAGE and ESC")
      }
      if (event_trigger_signature != latest_triage_trigger) {
        invalid("RUN-RESUME Trigger signature does not match its governing TRIAGE")
      }
      if (event_attempt + 0 != latest_triage_attempt + 1) {
        invalid("RUN-RESUME Attempt count must increment its governing TRIAGE")
      }
      if (index(event_trigger_signature, event_continuing_task ":") != 1) {
        invalid("RUN-RESUME Continuing task does not match its Trigger signature")
      }
      resume_decision = decision_id(event_resolution_evidence)
      if (resume_decision != "" && !decision_exists(resume_decision)) {
        invalid("RUN-RESUME cites an unrecorded phase decision: " resume_decision)
      }
    }
  }
  if (strict && directive_strict && event == "STOP") {
    if (latest_triage_id == "" || latest_triage_verdict !~ /^BLOCK/) {
      invalid("STOP has no preceding TRIAGE BLOCK")
    } else if (!references_id(event_triage_reference, latest_triage_id) ||
               latest_triage_esc == "" || !references_id(event_triage_reference, latest_triage_esc)) {
      invalid("STOP does not reference its TRIAGE BLOCK and concrete ESC")
    }
    if (event_protected_category == "" || event_next_human_action == "") {
      invalid("STOP is missing protected category or next Human action")
    }
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
    recovery_trigger = event_trigger_signature
    recovery_attempt = event_attempt
    recovery_mode = event_repair_mode
    recovery_no_progress = event_no_progress
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
      if (event_outcome == "RESOLVED") {
        recovery_stage = "needs-complete"
      } else if (event_outcome == "NO-PROGRESS" && recovery_no_progress == "1/2") {
        recovery_stage = "needs-triage"
        triage_needed = 1
      } else {
        recovery_stage = "needs-continuation"
      }
      recovery_last_outcome = event_outcome
      if (review_open && event_task == review_owner) {
        review_repair_seen = 1
      }
    }
  } else if (event == "TASK-COMPLETE") {
    if (recovery_task != "") {
      if (event_task != recovery_task) {
        invalid("completed task " event_task " does not match recovery task " recovery_task)
      } else if (recovery_stage != "needs-complete") {
        invalid("task completed before a resolved REPAIR event: " event_task)
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
      if (review_open) {
        if (!review_task_complete) {
          invalid("QA repair was replaced before owning task completed: " review_owner)
        }
        if (event_id == review_id) {
          invalid("QA repair was self-replaced by the same review identity: " event_id)
        }
        if (event_rereview != review_id) {
          invalid("QA REPAIR does not reference required re-review " review_id ": " event_id)
        }
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
    latest_triage_id = event_id
    latest_triage_trigger = event_trigger_signature
    latest_triage_verdict = event_verdict
    latest_triage_attempt = event_attempt
    latest_triage_esc = escalation_id(event_verdict)
    triage_seen[event_id] = 1
    if (event_verdict ~ /^RESUME/) {
      triage_needed = 0
    } else if (event_verdict ~ /^REPAIR/) {
      triage_needed = 0
      triage_repair_pending = 1
    } else if (event_verdict ~ /^BLOCK/) {
      if (recovery_stage == "needs-triage") {
        recovery_task = ""
        recovery_stage = ""
      }
      stopped = 1
      stop_entry_pending = 1
      triage_needed = 0
    }
  } else if (event == "STOP") {
    stopped = 1
    stop_entry_pending = 0
  } else if (event == "RUN-RESUME") {
    if (stopped && latest_triage_verdict ~ /^BLOCK/) {
      resolved_block_triage = latest_triage_id
      resolved_block_evidence = event_resolution_evidence
    }
    stopped = 0
    if (event_contract == "v2") {
      strict = 1
      saw_v2 = 1
    }
    if (event_directive_schema == "v1") {
      directive_strict = 1
      saw_directive_schema = 1
    }
  } else if (event == "RUN-START") {
    if (event_contract == "v2") {
      strict = 1
      saw_v2 = 1
    }
    if (event_directive_schema == "v1") {
      directive_strict = 1
      saw_directive_schema = 1
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
  event_directive_schema = ""
  event_source = ""
  event_category = ""
  event_repair_mode = ""
  event_repair_authority = ""
  event_control_contract = ""
  event_control_mismatch = ""
  event_allowed_files = ""
  event_forbidden_changes = ""
  event_trigger_signature = ""
  event_attempt = ""
  event_no_progress = ""
  event_checkpoint = ""
  event_resolution_evidence = ""
  event_next_action = ""
  event_human_attention = ""
  event_control_scope_check = ""
  event_next_checkpoint = ""
  event_triage_reference = ""
  event_protected_category = ""
  event_next_human_action = ""
  event_triage_log = ""
  event_prior_run = ""
  event_continuing_task = ""
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
    event_id = $5
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
event == "TRIAGE" && /^- Triage log: / {
  event_triage_log = substr($0, length("- Triage log: ") + 1)
}
(event == "TASK-RECOVERY" || event == "REPAIR" || event == "TRIAGE" || event == "RUN-RESUME") && /^- Trigger signature: / {
  event_trigger_signature = substr($0, length("- Trigger signature: ") + 1)
}
(event == "TASK-RECOVERY" || event == "REPAIR" || event == "TRIAGE" || event == "RUN-RESUME") && /^- Attempt count: / {
  event_attempt = substr($0, length("- Attempt count: ") + 1)
}
(event == "TASK-RECOVERY" || event == "TRIAGE" || event == "RUN-RESUME") && /^- Checkpoint: / {
  event_checkpoint = substr($0, length("- Checkpoint: ") + 1)
}
(event == "TASK-RECOVERY" || event == "REPAIR" || event == "TRIAGE" || event == "RUN-RESUME") && /^- Resolution evidence: / {
  event_resolution_evidence = substr($0, length("- Resolution evidence: ") + 1)
}
(event == "RUN-START" || event == "RUN-RESUME") && /^- Sequence contract: / {
  event_contract = substr($0, length("- Sequence contract: ") + 1)
}
(event == "RUN-START" || event == "RUN-RESUME") && /^- Recovery directive schema: / {
  event_directive_schema = substr($0, length("- Recovery directive schema: ") + 1)
}
event == "RUN-RESUME" && /^- Prior run\/session or STOP: / {
  event_prior_run = substr($0, length("- Prior run/session or STOP: ") + 1)
}
event == "RUN-RESUME" && /^- Continuing task: / {
  event_continuing_task = substr($0, length("- Continuing task: ") + 1)
}
event == "TASK-RECOVERY" && /^- Repair mode: / {
  event_repair_mode = substr($0, length("- Repair mode: ") + 1)
}
event == "TASK-RECOVERY" && /^- Source: / {
  event_source = substr($0, length("- Source: ") + 1)
}
event == "TASK-RECOVERY" && /^- Category: / {
  event_category = substr($0, length("- Category: ") + 1)
}
event == "TASK-RECOVERY" && /^- Repair authority: / {
  event_repair_authority = substr($0, length("- Repair authority: ") + 1)
}
event == "TASK-RECOVERY" && /^- Control contract: / {
  event_control_contract = substr($0, length("- Control contract: ") + 1)
}
event == "TASK-RECOVERY" && /^- Control mismatch: / {
  event_control_mismatch = substr($0, length("- Control mismatch: ") + 1)
}
event == "TASK-RECOVERY" && /^- Allowed files: / {
  event_allowed_files = substr($0, length("- Allowed files: ") + 1)
}
event == "TASK-RECOVERY" && /^- Forbidden changes: / {
  event_forbidden_changes = substr($0, length("- Forbidden changes: ") + 1)
}
event == "TASK-RECOVERY" && /^- Consecutive no-progress attempts: / {
  event_no_progress = substr($0, length("- Consecutive no-progress attempts: ") + 1)
}
event == "TASK-RECOVERY" && /^- Next action: / {
  event_next_action = substr($0, length("- Next action: ") + 1)
}
event == "TASK-RECOVERY" && /^- Human attention required: / {
  event_human_attention = substr($0, length("- Human attention required: ") + 1)
}
event == "REPAIR" && /^- Control scope check: / {
  event_control_scope_check = substr($0, length("- Control scope check: ") + 1)
}
event == "REPAIR" && /^- Outcome: / {
  event_outcome = substr($0, length("- Outcome: ") + 1)
}
event == "REPAIR" && /^- Next checkpoint: / {
  event_next_checkpoint = substr($0, length("- Next checkpoint: ") + 1)
}
event == "STOP" && /^- Triage: / {
  event_triage_reference = substr($0, length("- Triage: ") + 1)
}
event == "STOP" && /^- Protected category: / {
  event_protected_category = substr($0, length("- Protected category: ") + 1)
}
event == "STOP" && (/^- Next Human or external-authority action: / || /^- Next human action: /) {
  event_next_human_action = substr($0, index($0, ": ") + 2)
}

END {
  finish_event()
  if (saw_v2 && !saw_directive_schema) {
    invalid("ledger uses Sequence contract v2 without Recovery directive schema: v1")
  }
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
    if (recovery_stage != "needs-complete") {
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
' "$1"
}

print_append_receipt() {
  case "$1" in
    'VALID state='*) state=${1#VALID state=} ;;
    *)
      echo "INVALID validator returned an unexpected result: $1" >&2
      exit 1
      ;;
  esac
  printf 'APPEND_CONFIRMED state=%s\n' "$state"
}

if [ "$mode" = validate ]; then
  validate_sequence "$ledger"
  exit $?
fi

ledger_dir=$(CDPATH= cd "$(dirname "$ledger")" && pwd)
ledger="$ledger_dir/$(basename "$ledger")"
tmp=$(mktemp -d "$ledger_dir/.workflow-sequence.XXXXXX")
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
candidate="$tmp/candidate.md"
framed="$tmp/framed-candidate.md"
combined="$tmp/combined-ledger.md"
cat > "$candidate"

if [ ! -s "$candidate" ]; then
  echo "INVALID append candidate is empty" >&2
  exit 1
fi
if ! sed -n '1p' "$candidate" | grep -q '^### '; then
  echo "INVALID append candidate must start with an event heading" >&2
  exit 1
fi
heading_count=$(awk '/^### / { count++ } END { print count + 0 }' "$candidate")
if [ "$heading_count" -ne 1 ]; then
  echo "INVALID append candidate must contain exactly one event heading" >&2
  exit 1
fi

printf '\n' > "$framed"
cat "$candidate" >> "$framed"

# An earlier invocation may have committed the event but lost its stdout
# receipt. An exact EOF match makes retrying the same candidate idempotent.
framed_size=$(wc -c < "$framed" | tr -d ' ')
tail -c "$framed_size" "$ledger" > "$tmp/current-suffix.md"
if cmp -s "$framed" "$tmp/current-suffix.md"; then
  if output=$(validate_sequence "$ledger"); then
    print_append_receipt "$output"
    exit $?
  else
    code=$?
    exit "$code"
  fi
fi

# Build beside the ledger so the validated replacement stays on one filesystem.
cp -p "$ledger" "$combined"
cat "$framed" >> "$combined"

if output=$(validate_sequence "$combined"); then
  :
else
  code=$?
  exit "$code"
fi

mv "$combined" "$ledger"
print_append_receipt "$output"
