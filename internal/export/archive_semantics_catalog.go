package export

// semanticFirstProfile is the first projection schema that holds each
// family added since schema 19. An archive of an earlier schema holds none.
// Every family's Effect grammar is the journal's catalog.
var semanticFirstProfile = map[string]int{"review.claimed": 36, "review.heartbeat": 36, "epic.closed": 35, "task.finished": 34, "task.approved": 33, "finding.closed": 31, "review.round-recorded": 30, "merge.observed": 28, "ci.result-recorded": 28, "task.handed-off": 26, "check.completed": 25, "check.stopped": 25, "check.started": 25, "check.process-started": 25, "check.result-recorded": 25, "check.hold-renewed": 25, "check.artifact-kept": 25, "check.artifact-refused": 25, "task.started": 24, "hold.waiting": 23, "task.answered": 23, "hold.heartbeat": 22, "hold.parked": 22, "task.taken-over": 22, "hold.taken": 20, "policy.snapshotted": 21}
