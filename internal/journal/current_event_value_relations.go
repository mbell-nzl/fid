package journal

import (
	"fmt"
	"reflect"
	"strings"
)

func validateCurrentValueRelations(kind string, o map[string]any) error {
	switch kind {
	case "Repository":
		if (o["dirty"] == true) != (o["dirty_digest"] != nil) {
			return fmt.Errorf("repository dirty digest mismatch")
		}
	case "Source":
		if len(o["display_label"].(string)) > 128 {
			return fmt.Errorf("source label exceeds bound")
		}
		if (o["source_type"] == "project-file") != (o["path"] != nil) || (o["source_type"] == "external-reference") != (o["external_id"] != nil) {
			return fmt.Errorf("source location mismatch")
		}
	case "PathChange":
		if o["before_digest"] == nil && o["after_digest"] == nil {
			return fmt.Errorf("path change has no before or after")
		}
	case "Reconciliation":
		captured := map[string]any{}
		for _, field := range []string{"repository", "checkpoint_id", "paths", "unmerged_paths", "operation", "ancestry"} {
			captured[field] = o[field]
		}
		digest, err := currentValueDigest(captured)
		if err != nil || digest != o["digest"] {
			return fmt.Errorf("reconciliation digest mismatch")
		}
		if err := currentSortedObjects(o["paths"], "path", false); err != nil {
			return err
		}
		return currentSortedValues(o["unmerged_paths"], false)
	case "Scope":
		return currentSortedValues(o["paths"], false)
	case "Authority":
		if o["mode"] == "force" {
			if o["reference"] == nil || o["grant_ref"] == nil {
				return fmt.Errorf("force recovery lacks captured authority")
			}
		} else if o["grant_ref"] != nil || o["reference"] != nil {
			return fmt.Errorf("ordinary recovery carries grant")
		}
	case "Condition":
		if (o["type"] == "budget-observed") != (o["budget_name"] != nil) || (o["type"] == "capability-available") != (o["capability_id"] != nil) {
			return fmt.Errorf("condition nullable binding mismatch")
		}
	case "Observation":
		if o["observed_at"] != o["liveness"].(map[string]any)["observed_at"] {
			return fmt.Errorf("liveness observation time mismatch")
		}
	case "ImmutableBudget", "Budget":
		return currentBudgetRelations(kind, o)
	case "ActorProfile":
		if len(o["producer"].(string)) > 128 || len(o["provenance_ref"].(string)) > 512 {
			return fmt.Errorf("actor metadata exceeds bound")
		}
		const actions = "artifact.verify,attention.resolve,claim.acquire,claim.expire,claim.release,claim.renew,concern.dispose,concern.record,deviation.record,envelope.activate,envelope.complete,envelope.create,evidence.link,finding.address,gate.register,gate.run,issue.transition,review.record,run.abandon,run.checkpoint,run.close,run.force-recover,run.observe,run.observe-budget,run.recover,run.start,stop.raise,stop.resolve"
		if err := validateCurrentValue("[]="+actions, o["capabilities"]); err != nil {
			return err
		}
		return currentSortedValues(o["capabilities"], true)
	case "GateDefinition":
		return currentGateDefinitionRelations(o)
	case "Stop":
		if !currentValuesEqual(o["subject"], o["resolution_condition"].(map[string]any)["subject"]) {
			return fmt.Errorf("stop condition subject mismatch")
		}
		if _, err := validateCurrentShape(o["details"], currentStopDetailShapes[o["code"].(string)]); err != nil {
			return err
		}
		if strings.HasPrefix(o["code"].(string), "budget-") && !currentValuesEqual(o["blocked_actions"], []string{"gate.run", "run.recover", "run.start"}) {
			return fmt.Errorf("budget Stop action mismatch")
		}
	case "Attention":
		if err := currentSortedValues(o["stop_ids"], false); err != nil {
			return err
		}
		if (o["status"] == "resolved") != (o["resolution_event_id"] != nil) {
			return fmt.Errorf("attention resolution mismatch")
		}
	case "ActivationResource":
		return currentResourceRelations(o)
	case "Activation":
		if err := currentSortedValues(o["component_digests"], false); err != nil {
			return err
		}
		return currentSortedObjects(o["resources"], "component", true)
	case "Run":
		return currentSortedObjects(o["budgets"], "name", false)
	case "Receipt":
		if err := currentSortedObjects(o["artifacts"], "id", false); err != nil {
			return err
		}
	}
	return nil
}
func currentBudgetRelations(kind string, o map[string]any) error {
	maximum, _ := currentIntegerValue(o["maximum"])
	if len(o["name"].(string)) > 64 {
		return fmt.Errorf("budget name too long")
	}
	if o["warning_threshold"] != nil {
		warning, _ := currentIntegerValue(o["warning_threshold"])
		if warning > maximum {
			return fmt.Errorf("budget warning exceeds maximum")
		}
	}
	if kind == "ImmutableBudget" {
		return nil
	}
	if o["consumed"] == nil || o["remaining"] == nil {
		if o["consumed"] != nil || o["remaining"] != nil || o["state"] != "unknown" {
			return fmt.Errorf("unknown budget is inconsistent")
		}
		return nil
	}
	consumed, _ := currentIntegerValue(o["consumed"])
	remaining, _ := currentIntegerValue(o["remaining"])
	expected := maximum - consumed
	if expected < 0 {
		expected = 0
	}
	if remaining != expected {
		return fmt.Errorf("budget remaining mismatch")
	}
	state := "available"
	if remaining == 0 {
		state = "exhausted"
	} else if o["warning_threshold"] != nil {
		warning, _ := currentIntegerValue(o["warning_threshold"])
		if remaining <= warning {
			state = "warning"
		}
	}
	if o["state"] != state {
		return fmt.Errorf("budget state mismatch")
	}
	return nil
}
func currentResourceRelations(o map[string]any) error {
	state := o["state"]
	if state == "pending" {
		if o["observed_identity"] != nil || o["evidence_ref"] != nil || o["error_code"] != nil {
			return fmt.Errorf("pending resource has observation")
		}
		return nil
	}
	if o["evidence_ref"] == nil || o["evidence_ref"] == "" {
		return fmt.Errorf("resource observation lacks evidence")
	}
	if (state == "error") != (o["error_code"] != nil) {
		return fmt.Errorf("resource error mismatch")
	}
	if state == "old" && !currentValuesEqual(o["observed_identity"], o["old_identity"]) {
		return fmt.Errorf("old resource identity mismatch")
	}
	if state == "new" && o["observed_identity"] != o["new_identity"] {
		return fmt.Errorf("new resource identity mismatch")
	}
	return nil
}
func currentSortedValues(value any, nonempty bool) error {
	values, ok := value.([]any)
	if !ok || (nonempty && len(values) == 0) {
		return fmt.Errorf("required sorted array missing")
	}
	for i := 1; i < len(values); i++ {
		less := false
		if a, ok := values[i-1].(string); ok {
			b, ok := values[i].(string)
			less = ok && a < b
		} else {
			a, ok := currentIntegerValue(values[i-1])
			b, valid := currentIntegerValue(values[i])
			less = ok && valid && a < b
		}
		if !less {
			return fmt.Errorf("array must be sorted and unique")
		}
	}
	return nil
}
func currentSortedObjects(value any, key string, nonempty bool) error {
	values, ok := value.([]any)
	if !ok {
		return fmt.Errorf("required object array missing")
	}
	keys := make([]any, 0, len(values))
	for _, v := range values {
		o, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid object array")
		}
		keys = append(keys, o[key])
	}
	return currentSortedValues(keys, nonempty)
}
func currentGateDefinitionRelations(o map[string]any) error {
	for field, max := range map[string]int64{"definition_version": 2147483647, "timeout_ms": 600000, "output_bytes": 10485760, "excerpt_bytes": 16384, "artifact_bytes": 1073741824} {
		n, _ := currentIntegerValue(o[field])
		if n > max {
			return fmt.Errorf("gate %s exceeds bound", field)
		}
	}
	excerpt, _ := currentIntegerValue(o["excerpt_bytes"])
	output, _ := currentIntegerValue(o["output_bytes"])
	if excerpt > output {
		return fmt.Errorf("gate excerpt exceeds output")
	}
	argv := o["argv"].([]any)
	if len(argv) == 0 || len(argv) > 128 {
		return fmt.Errorf("gate argv size")
	}
	for _, field := range []string{"success_exit_codes", "credential_names", "redact_literals", "artifact_paths"} {
		if err := currentSortedValues(o[field], field == "success_exit_codes"); err != nil {
			return err
		}
	}
	for _, code := range o["success_exit_codes"].([]any) {
		n, _ := currentIntegerValue(code)
		if n < -255 || n > 255 {
			return fmt.Errorf("gate exit code range")
		}
	}
	if err := currentSortedObjects(o["environment"], "name", false); err != nil {
		return err
	}
	credentials := []any{}
	inherited := map[string]bool{"CI": true, "GOCACHE": true, "GOMODCACHE": true, "GOPATH": true, "HOME": true, "LANG": true, "LC_ALL": true, "PATH": true, "TMPDIR": true, "XDG_CACHE_HOME": true}
	for _, entry := range o["environment"].([]any) {
		e := entry.(map[string]any)
		if (e["source"] == "literal") != (e["value"] != nil) {
			return fmt.Errorf("gate environment source/value mismatch")
		}
		if e["source"] == "credential" {
			credentials = append(credentials, e["name"])
		}
		if e["source"] == "inherit" && !inherited[e["name"].(string)] {
			return fmt.Errorf("gate unaccepted inherited name")
		}
	}
	if !reflect.DeepEqual(credentials, o["credential_names"]) {
		return fmt.Errorf("gate credential names mismatch")
	}
	return nil
}
