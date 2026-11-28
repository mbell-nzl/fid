package status

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/epicplan"
	"github.com/mbell-nzl/fid/internal/epicrecord"
)

func (r *Result) addEpicPlans(q epicplan.Queryer, root string) error {
	rows, err := q.Query(`SELECT id FROM issues WHERE kind='epic' AND status NOT IN ('done','cancelled') ORDER BY id`)
	if err != nil {
		return err
	}
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return err
		}
		ids = append(ids, id)
	}
	err = rows.Err()
	_ = rows.Close()
	if err != nil {
		return err
	}
	for _, id := range ids {
		p, err := epicplan.Preview(q, id)
		if err != nil {
			return err
		}
		a, err := epicplan.Current(q, id)
		if err != nil {
			return err
		}
		unstarted := false
		for _, t := range p.Tasks {
			unstarted = unstarted || t.Kind == "task" && t.State == "todo"
		}
		if unstarted && a.State != "current" {
			action := fmt.Sprintf(`fid epic approve %s --actor human:<name> --reason "<why>"`, id)
			r.NeedsYou = append(r.NeedsYou, NeedsYouItem{Kind: "epic-approval", Subject: &Subject{Type: "epic", ID: id}, Summary: id + " plan approval is " + a.State, NextAction: &action})
		}
		eligible, err := epicrecord.ReadEligibility(q, id)
		if err != nil {
			return err
		}
		if !eligible.CanClose() {
			continue
		}
		views, err := epicplan.Views(q, id, root)
		if err != nil {
			return err
		}
		for _, v := range views {
			if v.State == "current" {
				continue
			}
			action := fmt.Sprintf(`fid epic validation pass %s %s --actor human:<name> --reason "<why>"`, id, v.ID)
			kind := "epic-validation-pass"
			if strings.HasPrefix(v.Verifier, "check:") {
				action = fmt.Sprintf("fid epic validation verify %s --actor <actor> --idempotency-key <key>", id)
				kind = "epic-validation-verify"
			}
			r.NeedsYou = append(r.NeedsYou, NeedsYouItem{Kind: kind, Subject: &Subject{Type: "epic", ID: id}, Summary: fmt.Sprintf("%s %s %s: %s", id, v.ID, v.Verifier, v.State), NextAction: &action})
		}
	}
	return nil
}
