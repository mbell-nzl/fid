package workflow

import (
	"database/sql"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issuestate"
)

// These aliases and forwarding functions keep the kernel APIs available until
// their removal. Retained callers use issuestate directly.
type (
	Definition            = issuestate.Definition
	RuleReference         = issuestate.RuleReference
	Registry              = issuestate.Registry
	NextState             = issuestate.NextState
	State                 = issuestate.State
	ResolvedIssueWorkflow = issuestate.ResolvedIssueWorkflow
	RuleResult            = issuestate.RuleResult
	UnmetCondition        = issuestate.UnmetCondition
	ChildTaskStatus       = issuestate.ChildTaskStatus
)

const (
	RuleAllChildTasksDoneOrCanceled = issuestate.RuleAllChildTasksDoneOrCanceled
	RuleRequiredGatesCurrent        = issuestate.RuleRequiredGatesCurrent
	RuleHasCompletion               = issuestate.RuleHasCompletion
)

func Load(conn *sql.DB) (*Registry, error) { return issuestate.Load(conn) }

func LoadFrom(conn fiddb.Queryer) (*Registry, error) { return issuestate.LoadFrom(conn) }

func Workflow(workflowID string) (Definition, bool) { return issuestate.Workflow(workflowID) }

func CategoryForState(workflowID, state string) (string, bool) {
	return issuestate.CategoryForState(workflowID, state)
}

func IsAllowedEdge(workflowID, fromState, toState string) bool {
	return issuestate.IsAllowedEdge(workflowID, fromState, toState)
}

func RulesForEdge(workflowID, fromState, toState string) ([]RuleReference, bool) {
	return issuestate.RulesForEdge(workflowID, fromState, toState)
}

func FrozenDefinitions() map[string]Definition { return issuestate.FrozenDefinitions() }

func InitialState(states []State) string { return issuestate.InitialState(states) }

func ResolveIssueWorkflow(conn *sql.DB, kind string) (ResolvedIssueWorkflow, error) {
	return issuestate.ResolveIssueWorkflow(conn, kind)
}

func AllChildTasksDoneOrCancelled(children []ChildTaskStatus) RuleResult {
	return issuestate.AllChildTasksDoneOrCancelled(children)
}
