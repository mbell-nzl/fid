package journal

import "time"

const (
	EventSchemaV1      = 1
	CurrentEventSchema = EventSchemaV1
)

type Event struct {
	EventID     string    `json:"event_id"`
	EventSchema int       `json:"event_schema"`
	EventType   string    `json:"event_type"`
	OccurredAt  time.Time `json:"occurred_at"`
	Actor       Actor     `json:"actor"`
	Command     Command   `json:"command"`
	Effects     []Effect  `json:"effects"`
}

type Actor struct {
	Role  string  `json:"role"`
	Actor *string `json:"actor"`
}

type Command struct {
	Argv       []string `json:"argv"`
	CWD        string   `json:"cwd"`
	FIDVersion string   `json:"fid_version"`
}

type Effect struct {
	Op            string         `json:"op"`
	Table         string         `json:"table"`
	PK            map[string]any `json:"pk"`
	Values        map[string]any `json:"values,omitempty"`
	Prior         map[string]any `json:"prior,omitempty"`
	Preconditions map[string]any `json:"preconditions,omitempty"`
}
