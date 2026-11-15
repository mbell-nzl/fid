package cli

func init() {
	for _, code := range []Code{
		"done-task-not-in-review",
		"done-epic-not-open",
		"done-epic-open-tasks",
		"done-checks-unsatisfied",
		"done-kept-artifact-unavailable",
		"done-open-major",
		"done-approval-missing",
		"done-approval-stale",
		"done-rollback-missing",
		"done-not-merged",
		"done-ci-waiting",
		"done-ci-not-green",
		"done-observation-stale",
		"done-handoff-required",
		"done-review-required",
		"done-separate-review-required",
		"done-evidence-required",
	} {
		retryable := code == "done-not-merged" || code == "done-ci-waiting" || code == "done-ci-not-green" || code == "done-observation-stale"
		failureCodeDefinitions[code] = definition(CategoryPrecondition, retryable, []SubjectType{SubjectIssue}, stringField("gate"), stringField("reason"), stringsField("record_ids"))
		registeredCodes = append(registeredCodes, code)
	}
	failureCodeDefinitions["github-unavailable"] = definition(CategoryPrecondition, true, []SubjectType{SubjectIssue}, stringField("reason"), booleanField("merged"))
	registeredCodes = append(registeredCodes, "github-unavailable")
}
