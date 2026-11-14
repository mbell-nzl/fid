package issuestate

// State preserves a workflow state's declared order for initial-state resolution.
type State struct {
	Name     string
	Category string
}

// InitialState returns the first unstarted state, or the first declared state.
func InitialState(states []State) string {
	for _, state := range states {
		if state.Category == "unstarted" {
			return state.Name
		}
	}
	if len(states) == 0 {
		return ""
	}
	return states[0].Name
}
