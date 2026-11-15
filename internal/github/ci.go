package github

import "slices"

type Verdict string

const (
	Green    Verdict = "green"
	Waiting  Verdict = "waiting"
	NotGreen Verdict = "not-green"
)

type IgnoreReason string

const (
	WrongHead  IgnoreReason = "different-head"
	WrongEvent IgnoreReason = "different-event"
	Superseded IgnoreReason = "superseded"
	Skipped    IgnoreReason = "skipped"
)

type WorkflowRun struct {
	ID         int64        `json:"id"`
	Attempt    int64        `json:"run_attempt"`
	WorkflowID int64        `json:"workflow_id"`
	Name       string       `json:"name"`
	Path       string       `json:"path"`
	Event      string       `json:"event"`
	HeadSHA    string       `json:"head_sha"`
	Status     string       `json:"status"`
	Conclusion string       `json:"conclusion"`
	URL        string       `json:"html_url"`
	Ignored    IgnoreReason `json:"-"`
}

type CI struct {
	Commit  string
	Runs    []WorkflowRun
	Verdict Verdict
	// NotGreen identifies each selected workflow with a completed non-success
	// conclusion. All observations, including ignored ones, remain in Runs.
	NotGreen []WorkflowRun
}

func evaluateCI(sha string, observed []WorkflowRun) CI {
	facts := CI{Commit: sha, Runs: slices.Clone(observed), Verdict: Waiting}
	latest := make(map[int64]WorkflowRun)
	for _, run := range facts.Runs {
		if run.HeadSHA != sha || run.Event != "push" {
			continue
		}
		previous, exists := latest[run.WorkflowID]
		// Run IDs order distinct executions; attempts order reruns of one ID.
		if !exists || run.ID > previous.ID || run.ID == previous.ID && run.Attempt > previous.Attempt {
			latest[run.WorkflowID] = run
		}
	}
	count, waiting := 0, false
	for i := range facts.Runs {
		run := &facts.Runs[i]
		chosen := latest[run.WorkflowID]
		switch {
		case run.HeadSHA != sha:
			run.Ignored = WrongHead
		case run.Event != "push":
			run.Ignored = WrongEvent
		case run.ID != chosen.ID || run.Attempt != chosen.Attempt:
			run.Ignored = Superseded
		case run.Status == "completed" && run.Conclusion == "skipped":
			run.Ignored = Skipped
		default:
			count++
			if run.Status != "completed" {
				waiting = true
			} else if run.Conclusion != "success" {
				facts.NotGreen = append(facts.NotGreen, *run)
			}
		}
	}
	switch {
	case len(facts.NotGreen) > 0:
		facts.Verdict = NotGreen
	case count > 0 && !waiting:
		facts.Verdict = Green
	}
	return facts
}
