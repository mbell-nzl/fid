package export

import "strings"

type semanticEventRule struct {
	aggregate string
	stages    []semanticEffectStage
}

type semanticEffectStage struct {
	minimum int
	maximum int
	allowed map[string]struct{}
}

func semanticRule(aggregate string, stages ...semanticEffectStage) semanticEventRule {
	return semanticEventRule{aggregate: aggregate, stages: stages}
}

func semanticStage(minimum int, operations ...string) semanticEffectStage {
	allowed := make(map[string]struct{}, len(operations))
	for _, operation := range operations {
		allowed[operation] = struct{}{}
	}
	return semanticEffectStage{minimum: minimum, maximum: -1, allowed: allowed}
}

func semanticEffectKey(table, operation string) string { return table + "\x00" + operation }

func (stage semanticEffectStage) accepts(table, operation string) bool {
	_, found := stage.allowed[semanticEffectKey(table, operation)]
	return found
}

// semanticEventRules is the closed Event-v2 family catalog. Stages preserve
// the accepted projection order while allowing the documented optional and
// repeated rows within a stage. An idempotency result, where present, is last.
var semanticEventRules = map[string]semanticEventRule{
	"check.completed":                semanticRule("check", semanticStage(1, semanticEffectKey("check_requests", "insert")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"check.stopped":                  semanticRule("check", semanticStage(1, semanticEffectKey("check_requests", "update")), semanticStage(1, semanticEffectKey("idempotency_results", "update"))),
	"check.started":                  semanticRule("check", semanticStage(0, semanticEffectKey("check_requests", "insert")), semanticStage(1, semanticEffectKey("check_starts", "insert")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"), semanticEffectKey("idempotency_results", "update"))),
	"check.process-started":          semanticRule("check", semanticStage(1, semanticEffectKey("check_processes", "insert"))),
	"check.result-recorded":          semanticRule("check", semanticStage(1, semanticEffectKey("check_results", "insert")), semanticStage(0, semanticEffectKey("check_requests", "update")), semanticStage(1, semanticEffectKey("idempotency_results", "update"))),
	"check.hold-renewed":             semanticRule("check", semanticStage(1, semanticEffectKey("holds", "update"))),
	"check.artifact-kept":            semanticRule("check", semanticStage(0, semanticEffectKey("referenced_artifacts", "insert"), semanticEffectKey("referenced_artifacts", "update")), semanticStage(1, semanticEffectKey("check_files", "insert"))),
	"check.artifact-refused":         semanticRule("check", semanticStage(1, semanticEffectKey("check_files", "insert"))),
	"task.handed-off":                semanticRule("task", semanticStage(1, semanticEffectKey("handoffs", "insert")), semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("holds", "update")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"task.started":                   semanticRule("task", semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"policy.snapshotted":             semanticRule("policy", semanticStage(1, semanticEffectKey("policy_snapshots", "insert"))),
	"project.initialized":            semanticRule("project"),
	"artifact.availability-observed": semanticRule("artifact", semanticStage(0, semanticEffectKey("referenced_artifacts", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"issue.created":                  semanticRule("issue", semanticStage(1, semanticEffectKey("issues", "insert")), semanticStage(0, semanticEffectKey("issue_labels", "insert")), semanticStage(0, semanticEffectKey("epic_metadata", "insert")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"issue.updated":                  semanticRule("issue", semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(0, semanticEffectKey("issue_labels", "insert"), semanticEffectKey("issue_labels", "delete")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"issue.transitioned":             semanticRule("issue", semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(0, semanticEffectKey("holds", "update")), semanticStage(0, semanticEffectKey("questions", "update")), semanticStage(0, semanticEffectKey("person_decisions", "insert")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"issue.ac-added":                 semanticRule("issue", semanticStage(1, semanticEffectKey("acceptance_criteria", "insert")), semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"issue.label-added":              semanticRule("issue", semanticStage(1, semanticEffectKey("issue_labels", "insert")), semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"issue.label-removed":            semanticRule("issue", semanticStage(1, semanticEffectKey("issue_labels", "delete")), semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"label.added":                    semanticRule("label", semanticStage(1, semanticEffectKey("labels", "insert")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"link.added":                     semanticRule("link", semanticStage(1, semanticEffectKey("links", "insert")), semanticStage(0, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"link.removed":                   semanticRule("link", semanticStage(1, semanticEffectKey("links", "delete")), semanticStage(0, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"hold.heartbeat":                 semanticRule("hold", semanticStage(1, semanticEffectKey("holds", "update"))),
	"hold.waiting":                   semanticRule("hold", semanticStage(1, semanticEffectKey("holds", "update")), semanticStage(1, semanticEffectKey("questions", "insert")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"task.answered":                  semanticRule("task", semanticStage(1, semanticEffectKey("person_decisions", "insert")), semanticStage(1, semanticEffectKey("questions", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"hold.parked":                    semanticRule("hold", semanticStage(1, semanticEffectKey("holds", "update")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"task.taken-over":                semanticRule("task", semanticStage(0, semanticEffectKey("holds", "update")), semanticStage(1, semanticEffectKey("person_decisions", "insert")), semanticStage(1, semanticEffectKey("holds", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"hold.taken":                     semanticRule("hold", semanticStage(0, semanticEffectKey("holds", "update")), semanticStage(1, semanticEffectKey("holds", "insert"))),
}

// semanticFirstProfile is the first projection schema that holds each
// family added since schema 19. An archive of an earlier schema holds none.
var semanticFirstProfile = map[string]int{"task.handed-off": 26, "check.completed": 25, "check.stopped": 25, "check.started": 25, "check.process-started": 25, "check.result-recorded": 25, "check.hold-renewed": 25, "check.artifact-kept": 25, "check.artifact-refused": 25, "task.started": 24, "hold.waiting": 23, "task.answered": 23, "hold.heartbeat": 22, "hold.parked": 22, "task.taken-over": 22, "hold.taken": 20, "policy.snapshotted": 21}

func init() {
	for eventType, rule := range semanticEventRules {
		if rule.aggregate == "" {
			panic("semantic archive validation has no aggregate for " + eventType)
		}
		for _, stage := range rule.stages {
			for operation := range stage.allowed {
				if len(strings.Split(operation, "\x00")) != 2 {
					panic("semantic archive validation has malformed Effect operation")
				}
			}
		}
	}
}
