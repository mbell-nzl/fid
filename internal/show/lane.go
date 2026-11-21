package show

import (
	"fmt"
	"strings"
)

func (result Result) laneText(lane string) string {
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
