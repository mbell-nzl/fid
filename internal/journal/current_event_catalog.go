package journal

// The accepted lean Event families and their ordered Effect grammar.
type currentEventRule struct {
	aggregate string
	stages    []currentEffectStage
}
type currentEffectStage struct {
	minimum int
	allowed map[string]bool
}

func currentRule(aggregate string, stages ...currentEffectStage) currentEventRule {
	return currentEventRule{aggregate, stages}
}
func currentStage(minimum int, operations ...string) currentEffectStage {
	allowed := map[string]bool{}
	for _, op := range operations {
		allowed[op] = true
	}
	return currentEffectStage{minimum, allowed}
}
func currentEffectKey(table, op string) string { return table + "/" + op }

// Retained Event-v2 families keep their accepted Effect grammar. Lean families
// check every stage in order; older Issue writers also retain SQL guards.
var retainedCurrentRules = map[string]currentEventRule{
	"check.completed":                currentRule("check", currentStage(1, currentEffectKey("check_requests", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"check.stopped":                  currentRule("check", currentStage(1, currentEffectKey("check_requests", "update")), currentStage(1, currentEffectKey("idempotency_results", "update"))),
	"check.started":                  currentRule("check", currentStage(0, currentEffectKey("check_requests", "insert")), currentStage(1, currentEffectKey("check_starts", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"), currentEffectKey("idempotency_results", "update"))),
	"check.process-started":          currentRule("check", currentStage(1, currentEffectKey("check_processes", "insert"))),
	"check.result-recorded":          currentRule("check", currentStage(1, currentEffectKey("check_results", "insert")), currentStage(0, currentEffectKey("check_requests", "update")), currentStage(1, currentEffectKey("idempotency_results", "update"))),
	"check.hold-renewed":             currentRule("check", currentStage(1, currentEffectKey("holds", "update"))),
	"check.artifact-kept":            currentRule("check", currentStage(0, currentEffectKey("referenced_artifacts", "insert"), currentEffectKey("referenced_artifacts", "update")), currentStage(1, currentEffectKey("check_files", "insert"))),
	"check.artifact-refused":         currentRule("check", currentStage(1, currentEffectKey("check_files", "insert"))),
	"task.handed-off":                currentRule("task", currentStage(1, currentEffectKey("handoffs", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"task.started":                   currentRule("task", currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"policy.snapshotted":             currentRule("policy", currentStage(1, currentEffectKey("policy_snapshots", "insert"))),
	"hold.heartbeat":                 currentRule("hold", currentStage(1, currentEffectKey("holds", "update"))),
	"hold.waiting":                   currentRule("hold", currentStage(1, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("questions", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"task.answered":                  currentRule("task", currentStage(1, currentEffectKey("person_decisions", "insert")), currentStage(1, currentEffectKey("questions", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"hold.parked":                    currentRule("hold", currentStage(1, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"task.taken-over":                currentRule("task", currentStage(0, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("person_decisions", "insert")), currentStage(1, currentEffectKey("holds", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"hold.taken":                     currentRule("hold", currentStage(0, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("holds", "insert"))),
	"issue.label.added":              currentRule("issue", currentStage(1, currentEffectKey("issue_labels", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.label.removed":            currentRule("issue", currentStage(1, currentEffectKey("issue_labels", "delete")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"project.initialized":            currentRule("project"),
	"artifact.availability-observed": currentRule("artifact", currentStage(0, currentEffectKey("referenced_artifacts", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.created":                  currentRule("issue", currentStage(1, currentEffectKey("issues", "insert")), currentStage(0, currentEffectKey("issue_labels", "insert")), currentStage(0, currentEffectKey("epic_metadata", "insert")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.updated":                  currentRule("issue", currentStage(1, currentEffectKey("issues", "update")), currentStage(0, currentEffectKey("issue_labels", "insert"), currentEffectKey("issue_labels", "delete")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.transitioned":             currentRule("issue", currentStage(1, currentEffectKey("issues", "update")), currentStage(0, currentEffectKey("holds", "update")), currentStage(0, currentEffectKey("questions", "update")), currentStage(0, currentEffectKey("person_decisions", "insert")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.ac-added":                 currentRule("issue", currentStage(1, currentEffectKey("acceptance_criteria", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.label-added":              currentRule("issue", currentStage(1, currentEffectKey("issue_labels", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.label-removed":            currentRule("issue", currentStage(1, currentEffectKey("issue_labels", "delete")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"label.added":                    currentRule("label", currentStage(1, currentEffectKey("labels", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"link.added":                     currentRule("link", currentStage(1, currentEffectKey("links", "insert")), currentStage(0, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"link.removed":                   currentRule("link", currentStage(1, currentEffectKey("links", "delete")), currentStage(0, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
}

// leanEventFamilies are the lean core's families. Every Effect they carry
// is checked against their stages, in order, not only those on a domain
// table another family owns.
var leanEventFamilies = map[string]bool{"task.handed-off": true, "task.started": true, "hold.waiting": true, "task.answered": true, "hold.heartbeat": true, "hold.parked": true, "task.taken-over": true, "hold.taken": true, "policy.snapshotted": true}
