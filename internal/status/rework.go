package status

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func (result *Result) addReworkLimit(conn reviewrecord.Queryer, task string) error {
	r, err := doneeval.ReadReworkLimit(conn, task)
	if err != nil || r == nil || !r.Reached() {
		return err
	}
	summary := fmt.Sprintf("%s reached its rework cap (%d rounds asked for changes; cap %d)", task, r.Count, r.Cap)
	action := fmt.Sprintf("fid takeover %s --actor human:<name> --reason <text> or fid issue transition %s --to cancelled --actor human:<name> --reason <text>", task, task)
	result.NeedsYou = append(result.NeedsYou, NeedsYouItem{Kind: "rework-cap", Summary: summary, NextAction: &action, Subject: &Subject{Type: "task", ID: task}})
	return nil
}
