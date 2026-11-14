package checkrecord

// Response is retained under the retry key from the first start. Complete is
// set only with the final observed result. Unstarted names never get results.
type Response struct {
	TaskID     string    `json:"task_id"`
	Names      []string  `json:"names"`
	CheckIDs   []string  `json:"check_ids"`
	Checks     []Summary `json:"checks"`
	NotStarted []string  `json:"not_started"`
	Complete   bool      `json:"complete"`
	Replayed   bool      `json:"replayed"`
	Failure    *Refusal  `json:"failure"`
	EventID    string    `json:"event_id"`
	OccurredAt string    `json:"occurred_at"`
}

type Summary struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Outcome     string   `json:"outcome"`
	ExitCode    *int     `json:"exit_code"`
	Counts      bool     `json:"counts"`
	Reasons     []string `json:"reasons"`
	ArtifactIDs []string `json:"artifact_ids"`
}

type Refusal struct {
	Code    string         `json:"code"`
	TaskID  string         `json:"task_id"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
	Next    string         `json:"next"`
}

func (refusal *Refusal) Error() string { return refusal.Message }
