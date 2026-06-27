package policy

import (
	"fmt"
	"slices"
	"strings"
)

type MatrixRow struct {
	Action         string `json:"action"`
	Developer      bool   `json:"developer"`
	QAReviewer     bool   `json:"qa_reviewer"`
	QALead         bool   `json:"qa_lead"`
	TLLite         bool   `json:"tl_lite"`
	TL             bool   `json:"tl"`
	WorkflowRunner bool   `json:"workflow_runner"`
	Human          bool   `json:"human"`
}

var matrixRoleNames = []string{
	"developer",
	"qa_reviewer",
	"qa_lead",
	"tl_lite",
	"tl",
	"workflow_runner",
	"human",
}

var defaultMatrixRows = []MatrixRow{
	{Action: "issue.create.epic", TL: true, Human: true},
	{Action: "issue.create.task", Developer: true, QAReviewer: true, QALead: true, TL: true, WorkflowRunner: true, Human: true},
	{Action: "issue.create.subtask", Developer: true, TL: true, WorkflowRunner: true, Human: true},
	{Action: "issue.update", Developer: true, TL: true, WorkflowRunner: true, Human: true},
	{Action: "issue.complete", Developer: true, Human: true},
	{Action: "issue.ac.add", Developer: true, TL: true, Human: true},
	{Action: "issue.ac.verify.developer", Developer: true},
	{Action: "issue.ac.verify.qa_reviewer", QAReviewer: true},
	{Action: "issue.ac.verify.qa_lead", QALead: true},
	{Action: "issue.ac.verify.tl_lite"},
	{Action: "issue.ac.verify.tl", TL: true},
	{Action: "issue.ac.verify.workflow_runner"},
	{Action: "issue.ac.verify.human", Human: true},
	{Action: "issue.transition.in_progress", Developer: true, WorkflowRunner: true, Human: true},
	{Action: "issue.transition.dev_complete", Developer: true, Human: true},
	{Action: "issue.transition.failed_gate", Developer: true, Human: true},
	{Action: "issue.transition.qa_review", QAReviewer: true, QALead: true, WorkflowRunner: true, Human: true},
	{Action: "issue.transition.qa_pass", QALead: true, Human: true},
	{Action: "issue.transition.qa_block", QAReviewer: true, QALead: true, Human: true},
	{Action: "issue.transition.wrapped", TL: true, Human: true},
	{Action: "issue.transition.closed", Human: true},
	{Action: "issue.transition.done", TL: true, Human: true},
	{Action: "issue.transition.cancelled", TL: true, Human: true},
	{Action: "issue.claim.acquire", Developer: true, QAReviewer: true, QALead: true, TL: true, WorkflowRunner: true, Human: true},
	{Action: "issue.claim.renew", Developer: true, QAReviewer: true, QALead: true, TL: true, WorkflowRunner: true, Human: true},
	{Action: "issue.claim.force_acquire", TL: true, WorkflowRunner: true, Human: true},
	{Action: "issue.claim.release", Developer: true, QAReviewer: true, QALead: true, TL: true, WorkflowRunner: true, Human: true},
	{Action: "concern.raise", Developer: true, Human: true},
	{Action: "concern.address", Developer: true, Human: true},
	{Action: "maintenance.rebuild", TL: true, WorkflowRunner: true, Human: true},
	{Action: "requirement.sync", TL: true, Human: true},
	{Action: "link.add", Developer: true, QAReviewer: true, QALead: true, TLLite: true, TL: true, WorkflowRunner: true, Human: true},
	{Action: "link.remove", Developer: true, QAReviewer: true, QALead: true, TLLite: true, TL: true, WorkflowRunner: true, Human: true},
}

var matrixByAction = buildMatrixByAction(defaultMatrixRows)

func MatrixRows() []MatrixRow {
	return slices.Clone(defaultMatrixRows)
}

func MatrixRoleNames() []string {
	return slices.Clone(matrixRoleNames)
}

func MatrixEntry(action string) (MatrixRow, bool) {
	row, ok := matrixByAction[strings.TrimSpace(action)]
	return row, ok
}

func (row MatrixRow) allows(role string) (bool, error) {
	switch strings.TrimSpace(role) {
	case "developer":
		return row.Developer, nil
	case "qa_reviewer":
		return row.QAReviewer, nil
	case "qa_lead":
		return row.QALead, nil
	case "tl_lite":
		return row.TLLite, nil
	case "tl":
		return row.TL, nil
	case "workflow_runner":
		return row.WorkflowRunner, nil
	case "human":
		return row.Human, nil
	default:
		return false, fmt.Errorf("permissions matrix: unsupported role %q", strings.TrimSpace(role))
	}
}

func buildMatrixByAction(rows []MatrixRow) map[string]MatrixRow {
	index := make(map[string]MatrixRow, len(rows))
	for _, row := range rows {
		action := strings.TrimSpace(row.Action)
		if action == "" {
			panic("permissions matrix row missing action")
		}
		if _, exists := index[action]; exists {
			panic(fmt.Sprintf("permissions matrix duplicate action %q", action))
		}
		row.Action = action
		index[action] = row
	}
	return index
}
