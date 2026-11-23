// Package rowcorpus builds the same row probes for journal and import tests.
// Its wire structs avoid an import cycle with the journal package's own tests.
// Expected refusals are test data, never calculated by the admission rules.
package rowcorpus

import (
	"encoding/json"
	"fmt"
	"maps"
)

type fields = map[string]any

type effect struct {
	Op            string `json:"op"`
	Table         string `json:"table"`
	PK            fields `json:"pk"`
	Values        fields `json:"values,omitempty"`
	Prior         fields `json:"prior,omitempty"`
	Preconditions fields `json:"preconditions,omitempty"`
}

type aggregate struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	Sequence int    `json:"sequence"`
}

type event struct {
	Schema      int       `json:"event_schema"`
	ID          string    `json:"event_id"`
	Family      string    `json:"event_type"`
	At          string    `json:"occurred_at"`
	Actor       any       `json:"actor"`
	Aggregate   aggregate `json:"aggregate"`
	Command     any       `json:"command"`
	Effects     []effect  `json:"effects"`
	Idempotency fields    `json:"idempotency,omitempty"`
}

type probe struct {
	Name   string
	Prefix []event
	Event  event
	Table  string
	Reason string
	// Import may reject a malformed column/key against its declared schema
	// before it reaches shared row admission. Pin that identifying text too.
	ImportReason string
}

type corpus struct {
	Prefix    []event
	Controls  []probe
	Malformed []probe
}

// JSON generates fresh, independent Events for each caller. No generated
// row-identity JSON file needs to be committed or regenerated.
func JSON() ([]byte, error) {
	c := build()
	return json.Marshal(c)
}

const stamp = "2026-11-10T10:00:00Z"

func fx(op, table string, pk, values fields) effect {
	return effect{Op: op, Table: table, PK: pk, Values: values}
}

func ev(n int, family, kind, id string, seq int, effects ...effect) event {
	commands := map[string]string{"issue.created": "issue.create", "issue.updated": "issue.update", "issue.transitioned": "issue.transition", "issue.ac-added": "issue.ac.add", "issue.label.added": "issue.label.add", "link.added": "link.add", "label.added": "label.add", "project.initialized": "project.init"}
	command := commands[family]
	if command == "" {
		command = family
	}
	return event{Schema: 2, ID: fmt.Sprintf("EVT-20261110-%06d", n), Family: family, At: stamp, Actor: fields{"actor_id": "mark", "actor_type": "human", "producer": "fid-standalone", "provenance_ref": "project-local-attribution"}, Aggregate: aggregate{kind, id, seq}, Command: command, Effects: effects}
}

func activity(n int, id, name string) effect {
	return fx("activity", "activity_log", fields{"id": fmt.Sprintf("ACT-%d", n)}, fields{"journal_event_id": fmt.Sprintf("EVT-20261110-%06d", n), "ref_type": "issue", "ref_id": id, "event": name, "details": nil, "actor": "human:mark", "occurred_at": stamp})
}

func bump(id string) effect {
	e := fx("update", "issues", fields{"id": id}, fields{"revision": 2})
	e.Prior, e.Preconditions = fields{"revision": 1}, fields{"revision": 1}
	return e
}

func create(n int, id, kind string) event {
	v := fields{"kind": kind, "title": id, "workflow_id": kind + "-default", "status": "todo", "created_at": stamp, "updated_at": stamp}
	if kind == "task" {
		v["lane"] = "standard"
	}
	if kind == "subtask" {
		v["parent_id"] = "FID-a"
	}
	return ev(n, "issue.created", "issue", id, 1, fx("insert", "issues", fields{"id": id}, v), activity(n, id, "created"))
}

func retry(e event) event {
	e.Idempotency = fields{"key": "key-" + e.ID, "scope": "label.add:project-1", "request_digest": "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	e.Effects = append(e.Effects, fx("insert", "idempotency_results", fields{"operation_scope": e.Idempotency["scope"], "idempotency_key": e.Idempotency["key"]}, fields{"request_digest": e.Idempotency["request_digest"], "request_field_digests_json": "{}", "result_entity_id": e.Aggregate.ID, "result_event_id": e.ID, "aggregate_sequence": e.Aggregate.Sequence, "result_json": `{"replayed":false}`, "created_at": stamp}))
	return e
}

func clone(e event) event {
	// Seeds contain scalar row values. Copy each map and the Effect slice so
	// a malformed probe cannot change a control or another probe's rows.
	e.Idempotency = maps.Clone(e.Idempotency)
	rows := make([]effect, len(e.Effects))
	for i, row := range e.Effects {
		row.PK, row.Values = maps.Clone(row.PK), maps.Clone(row.Values)
		row.Prior, row.Preconditions = maps.Clone(row.Prior), maps.Clone(row.Preconditions)
		rows[i] = row
	}
	e.Effects = rows
	return e
}

func (c *corpus) control(name string, e event, prefix ...event) {
	c.Controls = append(c.Controls, probe{Name: name, Event: clone(e), Prefix: prefix})
}

func (c *corpus) bad(name string, e event, table, reason string, mutate func(*event), prefix ...event) {
	p := probe{Name: name, Event: clone(e), Table: table, Reason: reason, Prefix: prefix}
	mutate(&p.Event)
	c.Malformed = append(c.Malformed, p)
}

func value(i int, k string, v any) func(*event) { return func(e *event) { e.Effects[i].Values[k] = v } }
func key(i int, k string, v any) func(*event)   { return func(e *event) { e.Effects[i].PK[k] = v } }
func extra(i int, row effect) func(*event) {
	return func(e *event) { e.Effects = append(e.Effects[:i:i], append([]effect{row}, e.Effects[i:]...)...) }
}
func remove(i int) func(*event) {
	return func(e *event) { e.Effects = append(e.Effects[:i:i], e.Effects[i+1:]...) }
}

const (
	primary         = "must key exactly one row by its primary key"
	stage           = "stage may write at most 1 rows"
	scalarRevision  = "update requires one scalar row and one revision row"
	labelSubject    = "label key must belong to the Event's Issue and values must be empty"
	activitySubject = "activity must name its own Event and subject"
	metadataSubject = "metadata must belong to the Event's epic"
	columnOnce      = "a column changes at most once per Event"
	revisionGuard   = "revision requires its guarded prior plus one, in its own row"
	inverseKey      = "second key must be the declared link's inverse"
	linkArguments   = "declared key must match the Event's link arguments"
)

func build() corpus {
	project := ev(1, "project.initialized", "project", "project-1", 1)
	project.Effects = []effect{}
	project.Actor = fields{"actor_id": "fid-bootstrap", "actor_type": "service", "producer": "fid", "provenance_ref": "fid:init"}
	c := corpus{Prefix: []event{project, create(2, "FID-a", "task"), create(3, "FID-b", "task"), create(4, "FID-e", "epic"), create(5, "FID-s", "subtask")}}
	for i, name := range []string{"red", "blue"} {
		c.Prefix = append(c.Prefix, retry(ev(i+6, "label.added", "label", name, 1, fx("insert", "labels", fields{"name": name}, fields{}))))
	}
	update := ev(8, "issue.updated", "issue", "FID-a", 2, fx("update", "issues", fields{"id": "FID-a"}, fields{"title": "Changed", "updated_at": stamp}), bump("FID-a"), activity(8, "FID-a", "updated"))
	ac := ev(8, "issue.ac-added", "issue", "FID-a", 2, fx("insert", "acceptance_criteria", fields{"id": "AC-one"}, fields{"issue_id": "FID-a", "ordinal": 1, "text": "Works"}), bump("FID-a"), activity(8, "FID-a", "ac-added"))
	label := ev(8, "issue.label.added", "issue", "FID-a", 2, fx("insert", "issue_labels", fields{"issue_id": "FID-a", "label": "red"}, fields{}), bump("FID-a"), activity(8, "FID-a", "issue-label-added"))
	link := ev(8, "link.added", "link", "FID-a:blocks:FID-b", 1, fx("insert", "links", fields{"from_type": "issue", "from_id": "FID-a", "relation": "blocks", "to_type": "issue", "to_id": "FID-b"}, fields{"created_at": stamp}), fx("insert", "links", fields{"from_type": "issue", "from_id": "FID-b", "relation": "blocked_by", "to_type": "issue", "to_id": "FID-a"}, fields{"created_at": stamp}), bump("FID-a"), bump("FID-b"), activity(8, "FID-a", "link-added"))
	created := create(8, "FID-c", "task")
	createdLabels := clone(created)
	for _, name := range []string{"blue", "red"} {
		extra(1, fx("insert", "issue_labels", fields{"issue_id": "FID-c", "label": name}, fields{}))(&createdLabels)
	}
	epic := create(8, "FID-new-epic", "epic")
	extra(1, fx("insert", "issue_labels", fields{"issue_id": "FID-new-epic", "label": "red"}, fields{}))(&epic)
	extra(1, fx("insert", "epic_metadata", fields{"issue_id": "FID-new-epic"}, fields{"scheduling": "parallel"}))(&epic)
	definition := retry(ev(8, "label.added", "label", "green", 1, fx("insert", "labels", fields{"name": "green"}, fields{})))
	transition := ev(8, "issue.transitioned", "issue", "FID-s", 2, fx("update", "issues", fields{"id": "FID-s"}, fields{"status": "in_progress", "updated_at": stamp}), bump("FID-s"), activity(8, "FID-s", "status-changed"))
	transition.Effects[0].Prior = fields{"status": "todo"}
	transition.Effects[0].Preconditions = fields{"status": "todo", "workflow_id": "subtask-default"}
	for _, seed := range []struct {
		name string
		e    event
	}{{"update", update}, {"criterion", ac}, {"label", label}, {"link", link}, {"creation", created}, {"created labels", createdLabels}, {"epic with label", epic}, {"label definition", definition}, {"subtask transition", transition}} {
		c.control(seed.name, seed.e)
	}
	identityCases(&c, update, created)
	childCases(&c, update, ac, label, created, createdLabels, epic, definition)
	updateCases(&c, update, ac, transition)
	linkCases(&c, link)
	relationCases(&c, link)
	activityCases(&c, created, update, ac, label, link, transition)
	removeCases(&c, label, link)
	schemaRefusals(&c)
	return c
}

func schemaRefusals(c *corpus) {
	// Import validates the archive's declared schema before shared row
	// admission. These probes deliberately pin that earlier, specific rule.
	identifying := map[string]string{
		"criterion key contains another column": "Effect 1 primary key: must contain exactly the table primary-key columns",
		"label key contains extra column":       "Effect 1 primary key: must contain exactly the table primary-key columns",
		"link key has extra column":             "Effect 1 primary key: must contain exactly the table primary-key columns",
		"activity key uses another column":      "Effect 3 primary key: must contain exactly the table primary-key columns",
		"retry key has extra column":            "Effect 2 primary key: must contain exactly the table primary-key columns",
		"invalid column":                        `Effect 1 values: unknown column "no_such_column"`,
		"invalid text type":                     `Effect 1 values: column "title": requires string JSON value`,
		"invalid integer type":                  `Effect 1 values: column "priority": requires integer JSON value`,
		"created label has unknown scalar":      `Effect 2 values: unknown column "rowid"`,
	}
	for i := range c.Malformed {
		c.Malformed[i].ImportReason = identifying[c.Malformed[i].Name]
	}
}
