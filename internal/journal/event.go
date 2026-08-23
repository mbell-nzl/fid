package journal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

const (
	EventSchemaV1 = 1
	EventSchemaV2 = 2

	// CurrentEventSchema remains the legacy default until each retained writer
	// is converted to Event-v2 in its owning task.
	CurrentEventSchema = EventSchemaV1
)

type Event struct {
	EventID     string
	EventSchema int
	EventType   string
	OccurredAt  time.Time
	Actor       Actor
	Command     Command
	Effects     []Effect

	ActorV2       EventActor
	Aggregate     Aggregate
	CausationID   *string
	CommandName   string
	CorrelationID *string
	Idempotency   *EventIdempotency
	Repository    *Repository
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

type EventActor struct {
	ActorID       string  `json:"actor_id"`
	ActorType     string  `json:"actor_type"`
	Producer      string  `json:"producer"`
	ProvenanceRef string  `json:"provenance_ref"`
	SessionRef    *string `json:"session_ref"`
}

type Aggregate struct {
	ID       string `json:"id"`
	Sequence int    `json:"sequence"`
	Type     string `json:"type"`
}

type EventIdempotency struct {
	Key           string `json:"key"`
	RequestDigest string `json:"request_digest"`
	Scope         string `json:"scope"`
}

type Repository struct {
	BaseOID       string  `json:"base_oid"`
	CaptureSource string  `json:"capture_source"`
	Dirty         bool    `json:"dirty"`
	DirtyDigest   *string `json:"dirty_digest"`
	HeadOID       string  `json:"head_oid"`
	ProjectID     string  `json:"project_id"`
	TreeOID       string  `json:"tree_oid"`
	WorktreeID    string  `json:"worktree_id"`
}

type Effect struct {
	Op            string         `json:"op"`
	Table         string         `json:"table"`
	PK            map[string]any `json:"pk"`
	Values        map[string]any `json:"values,omitempty"`
	Prior         map[string]any `json:"prior,omitempty"`
	Preconditions map[string]any `json:"preconditions,omitempty"`
}

func (event Event) IsV2() bool {
	return event.EventSchema == EventSchemaV2
}

func (event Event) MarshalJSON() ([]byte, error) {
	switch event.EventSchema {
	case EventSchemaV1:
		return json.Marshal(eventV1Wire{event.Actor, event.Command, event.Effects, event.EventID, event.EventSchema, event.EventType, event.OccurredAt})
	case EventSchemaV2:
		if err := validateV2Event(event); err != nil {
			return nil, err
		}
		return json.Marshal(eventV2Wire{event.ActorV2, event.Aggregate, event.CausationID, event.CommandName, event.CorrelationID, event.Effects, event.EventID, event.EventSchema, event.EventType, event.Idempotency, event.OccurredAt, event.Repository})
	default:
		return nil, fmt.Errorf("marshal event %q: unsupported schema %d", event.EventID, event.EventSchema)
	}
}

func (event *Event) UnmarshalJSON(data []byte) error {
	var header struct {
		EventSchema int `json:"event_schema"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return err
	}
	switch header.EventSchema {
	case EventSchemaV1:
		var wire eventV1Wire
		if err := json.Unmarshal(data, &wire); err != nil {
			return err
		}
		*event = Event{wire.EventID, wire.EventSchema, wire.EventType, wire.OccurredAt, wire.Actor, wire.Command, wire.Effects, EventActor{}, Aggregate{}, nil, "", nil, nil, nil}
		return nil
	case EventSchemaV2:
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()
		var wire eventV2Wire
		if err := decoder.Decode(&wire); err != nil {
			return err
		}
		decoded := Event{wire.EventID, wire.EventSchema, wire.EventType, wire.OccurredAt, Actor{}, Command{}, wire.Effects, wire.Actor, wire.Aggregate, wire.CausationID, wire.Command, wire.CorrelationID, wire.Idempotency, wire.Repository}
		if err := validateV2Event(decoded); err != nil {
			return err
		}
		*event = decoded
		return nil
	default:
		return fmt.Errorf("unsupported event schema %d", header.EventSchema)
	}
}

type eventV1Wire struct {
	Actor       Actor     `json:"actor"`
	Command     Command   `json:"command"`
	Effects     []Effect  `json:"effects"`
	EventID     string    `json:"event_id"`
	EventSchema int       `json:"event_schema"`
	EventType   string    `json:"event_type"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type eventV2Wire struct {
	Actor         EventActor        `json:"actor"`
	Aggregate     Aggregate         `json:"aggregate"`
	CausationID   *string           `json:"causation_id"`
	Command       string            `json:"command"`
	CorrelationID *string           `json:"correlation_id"`
	Effects       []Effect          `json:"effects"`
	EventID       string            `json:"event_id"`
	EventSchema   int               `json:"event_schema"`
	EventType     string            `json:"event_type"`
	Idempotency   *EventIdempotency `json:"idempotency"`
	OccurredAt    time.Time         `json:"occurred_at"`
	Repository    *Repository       `json:"repository"`
}

func validateV2Event(event Event) error {
	if event.EventID == "" {
		return fmt.Errorf("Event-v2 event_id is required")
	}
	if event.EventType == "" {
		return fmt.Errorf("Event-v2 event_type is required")
	}
	if event.OccurredAt.IsZero() {
		return fmt.Errorf("Event-v2 occurred_at is required")
	}
	if event.ActorV2.ActorID == "" || event.ActorV2.ActorType == "" || event.ActorV2.Producer == "" || event.ActorV2.ProvenanceRef == "" {
		return fmt.Errorf("Event-v2 actor is incomplete")
	}
	if event.Aggregate.ID == "" || event.Aggregate.Type == "" || event.Aggregate.Sequence < 1 {
		return fmt.Errorf("Event-v2 aggregate is incomplete")
	}
	if event.CommandName == "" {
		return fmt.Errorf("Event-v2 command is required")
	}
	if event.Effects == nil {
		return fmt.Errorf("Event-v2 effects must be explicit")
	}
	if event.EventType == "project.initialized" {
		if event.Repository != nil || event.Idempotency != nil || len(event.Effects) != 0 {
			return fmt.Errorf("Event-v2 bootstrap must have null repository/idempotency and empty effects")
		}
		return nil
	}
	if len(event.Effects) == 0 {
		return fmt.Errorf("Event-v2 %q requires effects", event.EventType)
	}
	return nil
}
