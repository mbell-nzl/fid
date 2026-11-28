package show

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/epicplan"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

type EpicPlan struct {
	Approval    epicplan.Approval `json:"approval"`
	Plan        epicplan.Plan     `json:"plan"`
	Validations []epicplan.View   `json:"validations"`
}

func (r *Result) readEpic(root string) error {
	if r.Issue.Kind != "epic" {
		return nil
	}
	q, err := reviewwork.OpenProjection(root)
	if err != nil {
		return err
	}
	defer func() { _ = q.Close() }()
	e := &EpicPlan{}
	e.Plan, err = epicplan.Preview(q, r.Issue.ID)
	if err != nil {
		return err
	}
	e.Approval, err = epicplan.Current(q, r.Issue.ID)
	if err != nil {
		return err
	}
	e.Validations, err = epicplan.Views(q, r.Issue.ID, root)
	r.EpicPlan = e
	return err
}

func (text textWriter) epic(e *EpicPlan) []string {
	if e == nil {
		return nil
	}
	a := e.Approval.State
	if e.Approval.ID != "" {
		a += " · " + e.Approval.Actor + " · " + e.Approval.At
	}
	lines := text.value("Approval", a)
	plan := []string{}
	for _, t := range e.Plan.Tasks {
		value := fmt.Sprintf("%s %s · %s", t.ID, t.Title, t.State)
		if t.Kind == "task" {
			lane := t.Lane
			if lane == "" {
				lane = "not set"
			}
			value += " · " + lane
		}
		plan = append(plan, value)
	}
	if len(plan) == 0 {
		plan = []string{"No tasks"}
	}
	lines = append(lines, text.value("Plan", strings.Join(plan, "; "))...)
	lines = append(lines, indent+"VALIDATIONS")
	for _, v := range e.Validations {
		value := fmt.Sprintf("%s · %s · %s", v.Claim, v.Verifier, v.State)
		if v.Recipe != "" {
			value += " · " + v.Recipe
		}
		if v.Pass != nil {
			value += " · " + v.Pass.Commit
		}
		if len(v.Rounds) > 0 {
			value += " · rounds " + strings.Join(v.Rounds, ", ")
		}
		lines = append(lines, text.value(v.ID, value)...)
	}
	return lines
}
