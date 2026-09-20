package review

import "testing"

func TestValidateVerdictFindingAgreement(t *testing.T) {
	tests := []struct {
		name     string
		verdict  string
		findings []FindingInput
		coverage Coverage
		valid    bool
	}{
		{name: "pass without findings", verdict: "pass", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}}, valid: true},
		{name: "pass with blocker", verdict: "pass", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}}, findings: []FindingInput{{Severity: "blocker", Blocking: true}}, valid: false},
		{name: "changes required with blocker", verdict: "changes-required", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}}, findings: []FindingInput{{Severity: "major", Blocking: true}}, valid: true},
		{name: "changes required without blocker", verdict: "changes-required", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}}, findings: []FindingInput{{Severity: "minor", Blocking: false}}, valid: false},
		{name: "blocked with exclusion", verdict: "blocked", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}, Exclusions: []CoverageExclusion{{Reason: "unavailable"}}}, valid: true},
		{name: "blocked without finding or exclusion", verdict: "blocked", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}}, valid: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateVerdictFindingAgreement(test.verdict, test.findings, test.coverage)
			if (err == nil) != test.valid {
				t.Fatalf("validateVerdictFindingAgreement() error = %v, valid = %v", err, test.valid)
			}
		})
	}
}
