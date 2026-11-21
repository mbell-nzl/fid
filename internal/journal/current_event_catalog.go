package journal

// The accepted lean Event families and their ordered Effect grammar.
type currentEventRule struct {
	aggregate string
	actors    string
	stages    []currentEffectStage
}
type currentEffectStage struct {
	minimum int
	allowed map[string]bool
}

func currentRule(aggregate, actors string, stages ...currentEffectStage) currentEventRule {
	return currentEventRule{aggregate, actors, stages}
}

// The actor types each family admits, as the writing command's trust rule
// (ADR-0014): the work set, a person's decisions, fid itself, or any actor.
const (
	workActors   = "agent human"
	personActors = "human"
	fidActors    = "service"
	anyActors    = "agent human service"
)

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
	"review.claimed":                 currentRule("review-claim", workActors, currentStage(0, currentEffectKey("review_claims", "update")), currentStage(1, currentEffectKey("review_claims", "insert"))),
	"review.heartbeat":               currentRule("review-claim", workActors, currentStage(1, currentEffectKey("review_claims", "update"))),
	"epic.closed":                    currentRule("epic", personActors, currentStage(1, currentEffectKey("person_decisions", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"task.finished":                  currentRule("task", workActors, currentStage(1, currentEffectKey("completions", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("holds", "update")), currentStage(0, currentEffectKey("questions", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"task.approved":                  currentRule("decision", personActors, currentStage(1, currentEffectKey("person_decisions", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"finding.closed":                 currentRule("finding", personActors, currentStage(1, currentEffectKey("person_decisions", "insert")), currentStage(1, currentEffectKey("finding_closures", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"review.round-recorded":          currentRule("review", workActors, currentStage(1, currentEffectKey("review_rounds", "insert")), currentStage(0, currentEffectKey("round_findings", "insert")), currentStage(0, currentEffectKey("finding_closures", "insert")), currentStage(0, currentEffectKey("issues", "update")), currentStage(0, currentEffectKey("holds", "update")), currentStage(0, currentEffectKey("review_claims", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"merge.observed":                 currentRule("merge", workActors, currentStage(1, currentEffectKey("merge_observations", "insert"))),
	"ci.result-recorded":             currentRule("merge", workActors, currentStage(1, currentEffectKey("ci_results", "insert"))),
	"check.completed":                currentRule("check", workActors, currentStage(1, currentEffectKey("check_requests", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"check.stopped":                  currentRule("check", workActors, currentStage(1, currentEffectKey("check_requests", "update")), currentStage(1, currentEffectKey("idempotency_results", "update"))),
	"check.started":                  currentRule("check", workActors, currentStage(0, currentEffectKey("check_requests", "insert")), currentStage(1, currentEffectKey("check_starts", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"), currentEffectKey("idempotency_results", "update"))),
	"check.process-started":          currentRule("check", workActors, currentStage(1, currentEffectKey("check_processes", "insert"))),
	"check.result-recorded":          currentRule("check", workActors, currentStage(1, currentEffectKey("check_results", "insert")), currentStage(0, currentEffectKey("check_requests", "update")), currentStage(1, currentEffectKey("idempotency_results", "update"))),
	"check.hold-renewed":             currentRule("check", workActors, currentStage(1, currentEffectKey("holds", "update"))),
	"check.artifact-kept":            currentRule("check", workActors, currentStage(0, currentEffectKey("referenced_artifacts", "insert"), currentEffectKey("referenced_artifacts", "update")), currentStage(1, currentEffectKey("check_files", "insert"))),
	"check.artifact-refused":         currentRule("check", workActors, currentStage(1, currentEffectKey("check_files", "insert"))),
	"task.handed-off":                currentRule("task", workActors, currentStage(1, currentEffectKey("handoffs", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"task.started":                   currentRule("task", workActors, currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"policy.snapshotted":             currentRule("policy", fidActors, currentStage(1, currentEffectKey("policy_snapshots", "insert"))),
	"hold.heartbeat":                 currentRule("hold", workActors, currentStage(1, currentEffectKey("holds", "update"))),
	"hold.waiting":                   currentRule("hold", workActors, currentStage(1, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("questions", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"task.answered":                  currentRule("task", personActors, currentStage(1, currentEffectKey("person_decisions", "insert")), currentStage(1, currentEffectKey("questions", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"hold.parked":                    currentRule("hold", workActors, currentStage(1, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"task.taken-over":                currentRule("task", personActors, currentStage(0, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("person_decisions", "insert")), currentStage(1, currentEffectKey("holds", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"hold.taken":                     currentRule("hold", workActors, currentStage(0, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("holds", "insert"))),
	"issue.label.added":              currentRule("issue", anyActors, currentStage(1, currentEffectKey("issue_labels", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.label.removed":            currentRule("issue", anyActors, currentStage(1, currentEffectKey("issue_labels", "delete")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"project.initialized":            currentRule("project", fidActors),
	"artifact.availability-observed": currentRule("artifact", workActors, currentStage(0, currentEffectKey("referenced_artifacts", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.created":                  currentRule("issue", anyActors, currentStage(1, currentEffectKey("issues", "insert")), currentStage(0, currentEffectKey("issue_labels", "insert")), currentStage(0, currentEffectKey("epic_metadata", "insert")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.updated":                  currentRule("issue", anyActors, currentStage(1, currentEffectKey("issues", "update")), currentStage(0, currentEffectKey("issue_labels", "insert"), currentEffectKey("issue_labels", "delete")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.transitioned":             currentRule("issue", anyActors, currentStage(1, currentEffectKey("issues", "update")), currentStage(0, currentEffectKey("holds", "update")), currentStage(0, currentEffectKey("questions", "update")), currentStage(0, currentEffectKey("person_decisions", "insert")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.ac-added":                 currentRule("issue", anyActors, currentStage(1, currentEffectKey("acceptance_criteria", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.label-added":              currentRule("issue", anyActors, currentStage(1, currentEffectKey("issue_labels", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.label-removed":            currentRule("issue", anyActors, currentStage(1, currentEffectKey("issue_labels", "delete")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"label.added":                    currentRule("label", anyActors, currentStage(1, currentEffectKey("labels", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"link.added":                     currentRule("link", anyActors, currentStage(1, currentEffectKey("links", "insert")), currentStage(0, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"link.removed":                   currentRule("link", anyActors, currentStage(1, currentEffectKey("links", "delete")), currentStage(0, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
}
