package status

import (
	"fmt"
	"sort"

	"github.com/mbell-nzl/fid/internal/workrecord"
)

func (result *Result) addUnaccepted(conn workrecord.Queryer) error {
	tasks, err := workrecord.Unaccepted(conn)
	if err != nil {
		return err
	}
	ids := make([]string, 0, len(tasks))
	for id := range tasks {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		action := fmt.Sprintf("fid accept %s --actor human:<name> or fid issue transition %s --to cancelled --actor human:<name> --reason <text>", id, id)
		result.NeedsYou = append(result.NeedsYou, NeedsYouItem{Kind: "accept", Summary: fmt.Sprintf("%s, created by %s, waits for you to accept it", id, tasks[id]), NextAction: &action, Subject: &Subject{Type: "issue", ID: id}})
	}
	return nil
}
