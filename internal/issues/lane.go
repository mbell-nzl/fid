package issues

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/issuestate"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
)

// LaneError identifies a lane flag that cannot apply to this issue.
type LaneError struct{ Kind, Lane string }

func (err *LaneError) Error() string {
	return fmt.Sprintf("lane %q requires a task and routine, standard or assured", err.Lane)
}

func laneRank(lane string) int {
	switch lane {
	case "routine":
		return 1
	case "standard":
		return 2
	case "assured":
		return 3
	}
	return 0
}

func validateLane(kind, lane string) error {
	if kind != "task" || laneRank(lane) == 0 {
		return &LaneError{Kind: kind, Lane: lane}
	}
	return nil
}

func validateLaneChange(current issueUpdateState, lane *string, actor string) error {
	if lane == nil {
		return nil
	}
	if err := validateLane(current.Kind, *lane); err != nil {
		return err
	}
	if !policy.Allows(actor, policy.WorkSet) {
		return &persondecision.DeniedError{Actor: actor}
	}
	if current.Lane == nil || laneRank(*lane) >= laneRank(*current.Lane) {
		return nil
	}
	if !strings.HasPrefix(actor, "human:") {
		return &persondecision.DeniedError{Actor: actor}
	}
	if current.Status != "todo" {
		return &BlockedTransitionError{IssueID: current.ID, CurrentState: current.Status, RequestedState: current.Status, Unmet: []issuestate.UnmetCondition{{Rule: "lane-lowering-requires-waiver"}}}
	}
	return nil
}
