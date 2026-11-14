package workflow

import "testing"

func TestRetiredJudgmentRulesAreNotImplemented(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{"all_acs_verified", "has_completion_report", "concerns_addressed"} {
		if IsImplementedRule(rule) {
			t.Errorf("retired rule %q remains implemented", rule)
		}
	}
}
