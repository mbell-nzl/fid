package show

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/decisionview"
)

func (result Result) laneText(lane string) string {
	if waiver := decisionview.LatestWaiver(result.Approval.Waivers); waiver != nil {
		return fmt.Sprintf("%s · waiver to %s · %s", lane, waiver.Binding.Lane, waiver.View().Text("waived an earlier change"))
	}
	handoff := result.Handoff
	if handoff == nil || handoff.RaisedFrom == nil || handoff.Floor == nil || len(handoff.Floor.Paths) == 0 {
		return lane
	}
	floor := handoff.Floor
	if floor.Raise("") != strings.ToLower(lane) {
		return lane
	}
	previous := *handoff.RaisedFrom
	previous = strings.ToUpper(previous[:1]) + previous[1:]
	path := floor.Paths[0]
	if len(floor.Paths) > 1 {
		path += fmt.Sprintf(" and %d others", len(floor.Paths)-1)
	}
	return fmt.Sprintf("%s · floor: %s · raised from %s", lane, path, previous)
}
