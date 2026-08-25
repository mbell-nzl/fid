package issues

import "strings"

const (
	actionIssueUpdate      = "issue.update"
	actionIssueComplete    = "issue.complete"
	actionIssueLabelAdd    = "issue.label.add"
	actionIssueLabelRemove = "issue.label.remove"
	actionIssueACAdd       = "issue.ac.add"
	actionConcernRaise     = "concern.raise"
	actionConcernAddress   = "concern.address"
)

func issueACVerifyAction(role string) string {
	return "issue.ac.verify." + strings.TrimSpace(role)
}
