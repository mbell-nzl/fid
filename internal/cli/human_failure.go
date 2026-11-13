package cli

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

// WriteHumanFailure renders the section 9 lines from the same failure as JSON.
// hasInput comes from the live leaf's flags; operation identifies the producer.
func WriteHumanFailure(writer io.Writer, operation, leaf string, failure Failure, hasInput bool) error {
	if writer == nil {
		return fmt.Errorf("write Human failure for %s: writer is nil", operation)
	}
	lines := []string{"error: " + failure.Message}
	reason := "reason: " + string(failure.Code)
	if detail := humanKeyDetail(failure); detail != "" {
		reason += "; " + detail
	}
	lines = append(lines, cutHumanReason(EscapeControlText(reason)))
	ids := append(humanStrings(failure.Details["stop_ids"]), humanStrings(failure.Details["attention_ids"])...)
	if len(ids) > 0 {
		lines = append(lines, "blocked by: "+EscapeControlText(strings.Join(ids, ",")))
	}
	if next := humanNext(failure, operation, leaf, hasInput); next != "" {
		lines = append(lines, "next: "+EscapeControlText(next))
	}
	if _, err := fmt.Fprintln(writer, strings.Join(lines, "\n")); err != nil {
		return fmt.Errorf("write Human failure for %s: %w", operation, err)
	}
	return nil
}

func cutHumanReason(text string) string {
	values := []rune(text)
	if len(values) <= 200 {
		return text
	}
	return string(values[:escapedPrefix(values, 200)]) + "…"
}

func humanKeyDetail(failure Failure) string {
	details := failure.Details
	switch failure.Code {
	case CodeRepositoryMismatch:
		return "expected HEAD " + humanRepository(details["expected"]) + ", observed HEAD " + humanRepository(details["observed"])
	case CodeIdempotencyConflict:
		return "differing fields: " + strings.Join(humanStrings(details["differing_fields"]), ",")
	case CodeTransitionNotPermitted:
		return strings.Join(humanStrings(details["unsatisfied"]), ",")
	case CodeObligationUnsatisfied, CodeStopUnresolved, CodeResolutionInvalid:
		return strings.Join(humanReasonCodes(details["reasons"]), ",")
	case CodeInvalidArgument:
		return "argument=" + fmt.Sprint(details["argument"]) + ", allowed=" + strings.Join(humanStrings(details["allowed"]), "|")
	}
	names := make([]string, 0, len(details))
	for name := range details {
		names = append(names, name)
	}
	sort.Strings(names)
	parts := []string{}
	for _, name := range names {
		if details[name] != nil && humanNestedDetail(failure.Code, name) {
			continue
		}
		if value, ok := humanScalar(details[name]); ok {
			parts = append(parts, name+"="+value)
		}
	}
	return strings.Join(parts, ", ")
}

func humanNestedDetail(code Code, name string) bool {
	for _, field := range failureCodeDefinitions[code].details {
		if field.name == name {
			return field.typeOf == detailObject || field.typeOf == detailObjects
		}
	}
	return false
}

func humanScalar(value any) (string, bool) {
	if value == nil {
		return "null", true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return fmt.Sprint(value), true
	case reflect.Slice, reflect.Array:
		parts := make([]string, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if item != nil {
				kind := reflect.ValueOf(item).Kind()
				if kind == reflect.Slice || kind == reflect.Array {
					return "", false
				}
			}
			part, ok := humanScalar(item)
			if !ok {
				return "", false
			}
			parts = append(parts, part)
		}
		return strings.Join(parts, ","), true
	}
	return "", false
}

func humanStrings(value any) []string {
	result := []string{}
	switch values := value.(type) {
	case []string:
		result = append(result, values...)
	case []any:
		for _, item := range values {
			if text, ok := item.(string); ok {
				result = append(result, text)
			}
		}
	}
	return result
}

func humanReasonCodes(value any) []string {
	result := []string{}
	switch values := value.(type) {
	case []map[string]any:
		for _, item := range values {
			result = append(result, fmt.Sprint(item["code"]))
		}
	case []any:
		for _, raw := range values {
			if item, ok := raw.(map[string]any); ok {
				result = append(result, fmt.Sprint(item["code"]))
			}
		}
	}
	return result
}

func humanRepository(value any) string {
	repository, _ := value.(map[string]any)
	head, _ := repository["head_oid"].(string)
	if len(head) > 12 {
		head = head[:12]
	}
	state := "clean"
	if repository["dirty"] == true {
		state = "dirty"
	}
	return head + " " + state
}

func humanNext(failure Failure, operation, leaf string, hasInput bool) string {
	switch failure.Code {
	case CodeHoldNotCurrent:
		return "fid next"
	case CodePolicyUnavailable:
		return failure.NextAction
	case CodeIssueNotFound:
		if operation == "issue.show" {
			return "fid issue list"
		}
	case CodeInvalidArgument:
		if hasInput {
			return "see " + leaf + " --help; template: " + leaf + " --example"
		}
	case CodeRepositoryMismatch:
		return "check HEAD and uncommitted changes, then correct expected_repository"
	case CodeSourceDigestMismatch:
		return "correct or omit digest for " + fmt.Sprint(failure.Details["display_label"])
	case CodeIdempotencyConflict:
		return "resend the original input, or use a new --idempotency-key"
	case CodeStopUnresolved:
		return "read each Stop's qualifying resolution: fid stop list --status open"
	case CodeReviewReworkExhausted:
		return "an activated Human records an approved Human Outcome of class review-rework-extension: fid human-outcome record (template: fid human-outcome record --example)"
	case CodeProtectedDecisionRequired:
		return "record the " + fmt.Sprint(failure.Details["decision_class"]) + " Decision Reference and an approving Human Outcome: fid decision record-reference, fid human-outcome record (templates: --example)"
	case CodeHumanOutcomeRequired:
		return "an activated Human records a " + fmt.Sprint(failure.Details["outcome_class"]) + " Human Outcome, then link it: fid human-outcome record, fid evidence link (templates: --example)"
	case CodeObligationUnsatisfied:
		return "see per-Obligation reasons: fid evidence evaluate --input <file> (template: fid evidence evaluate --example)"
	}
	return ""
}
