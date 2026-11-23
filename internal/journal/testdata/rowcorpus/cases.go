package rowcorpus

func identityCases(c *corpus, update, created event) {
	for _, item := range []struct {
		column string
		value  any
	}{{"id", "FID-z"}, {"kind", "subtask"}, {"workflow_id", "epic-default"}, {"parent_id", "FID-e"}, {"created_at", "2026-11-09T00:00:00Z"}, {"sequence", 99}, {"owner_role", "reviewer"}, {"owner_actor", "other"}} {
		c.bad("immutable "+item.column, update, "issues", item.column+" is immutable", value(0, item.column, item.value))
		legacy := clone(update)
		legacy.Schema, legacy.Actor, legacy.Command = 1, fields{"role": "developer"}, fields{"argv": []string{"fid", "issue", "update"}, "cwd": ".", "fid_version": "0.1.0"}
		c.bad("legacy immutable "+item.column, legacy, "issues", item.column+" is immutable", value(0, item.column, item.value))
	}
	legacy := clone(update)
	legacy.Schema, legacy.Actor, legacy.Command = 1, fields{"role": "developer"}, fields{"argv": []string{"fid", "issue", "update"}, "cwd": ".", "fid_version": "0.1.0"}
	c.bad("legacy deletes and recreates Issue", legacy, "issues", "an Issue cannot be deleted or recreated", func(e *event) {
		e.Effects = []effect{fx("delete", "issues", fields{"id": "FID-b"}, nil), create(8, "FID-b", "task").Effects[0]}
	})
	for _, item := range []struct {
		name, column string
		value        any
	}{{"column", "no_such_column", "bad"}, {"text type", "title", 123}, {"integer type", "priority", 1.5}} {
		c.bad("invalid "+item.name, update, "issues", "invalid column or scalar "+item.column, value(0, item.column, item.value))
	}
	c.bad("creation starts with a later revision", created, "issues", "a creation starts at revision one", value(0, "revision", 99999999))
	c.bad("two Issue creations", created, "issues", stage, extra(1, created.Effects[0]))
}

func childCases(c *corpus, update, ac, label, created, createdLabels, epic, definition event) {
	c.bad("criterion on another Issue", ac, "acceptance_criteria", "criterion must belong to the Event's Issue", value(0, "issue_id", "FID-b"))
	c.bad("criterion key contains another column", ac, "acceptance_criteria", primary, key(0, "issue_id", "FID-b"))
	c.bad("two criteria", ac, "acceptance_criteria", stage, extra(1, fx("insert", "acceptance_criteria", fields{"id": "AC-two"}, fields{"issue_id": "FID-a", "ordinal": 2, "text": "Also works"})))
	withLabels := clone(update)
	extra(2, label.Effects[0])(&withLabels)
	for _, item := range []struct {
		name  string
		e     event
		index int
	}{{"create", createdLabels, 1}, {"update", withLabels, 2}, {"add", label, 0}} {
		c.bad(item.name+" label on another Issue", item.e, "issue_labels", labelSubject, key(item.index, "issue_id", "FID-b"))
	}
	c.bad("label key contains extra column", label, "issue_labels", primary, key(0, "rowid", 1))
	c.bad("label values change key", label, "issue_labels", `conflicting values for key "issue_id"`, value(0, "issue_id", "FID-b"))
	c.bad("label values repeat its key", label, "issue_labels", labelSubject, value(0, "issue_id", "FID-a"))
	c.bad("two singular labels", label, "issue_labels", stage, extra(1, fx("insert", "issue_labels", fields{"issue_id": "FID-a", "label": "blue"}, fields{})))
	c.bad("duplicate label operation", createdLabels, "issue_labels", "a label changes at most once per Event", extra(2, createdLabels.Effects[1]))
	c.bad("metadata on another Issue", epic, "epic_metadata", metadataSubject, key(1, "issue_id", "FID-e"))
	c.bad("two metadata rows", epic, "epic_metadata", metadataSubject, extra(2, fx("insert", "epic_metadata", fields{"issue_id": "FID-e"}, fields{"scheduling": "parallel"})))
	c.bad("duplicate metadata row", epic, "epic_metadata", stage, extra(2, epic.Effects[1]))
	c.bad("task has epic metadata", created, "epic_metadata", metadataSubject, extra(1, fx("insert", "epic_metadata", fields{"issue_id": "FID-c"}, fields{"scheduling": "parallel"})))
	c.bad("definition names another label", definition, "labels", "label name must match the Event aggregate", key(0, "name", "yellow"))
	c.bad("two label definitions", definition, "labels", "label name must match the Event aggregate", extra(1, fx("insert", "labels", fields{"name": "yellow"}, fields{})))
	c.bad("duplicate label definition row", definition, "labels", stage, extra(1, definition.Effects[0]))
	c.bad("created label has unknown scalar", createdLabels, "issue_labels", "invalid column or scalar rowid", value(1, "rowid", 2))
	c.bad("retry key has extra column", definition, "idempotency_results", primary, key(1, "created_at", stamp))
}

func updateCases(c *corpus, update, ac, transition event) {
	c.bad("three Issue updates", update, "issues", columnOnce, extra(2, fx("update", "issues", fields{"id": "FID-a"}, fields{"title": "Again"})))
	c.bad("two updates of one scalar", update, "issues", columnOnce, extra(1, fx("update", "issues", fields{"id": "FID-a"}, fields{"title": "Again"})))
	c.bad("revision changes another field", ac, "issues", revisionGuard, value(1, "title", "Changed"))
	c.bad("transition repeats scalar", transition, "issues", columnOnce, extra(1, fx("update", "issues", fields{"id": "FID-s"}, fields{"updated_at": stamp})))
	// Distinct columns keep duplicate-column detection from masking a missing
	// revision requirement. Both update families must reach that exact rule.
	for _, item := range []struct {
		name string
		e    event
		col  string
		v    any
	}{{"update", update, "title", "Changed"}, {"transition", transition, "status", "in_progress"}} {
		c.bad(item.name+" two scalar rows without revision", item.e, "issues", scalarRevision, func(e *event) {
			e.Effects[0].Values = fields{"updated_at": stamp}
			e.Effects[0].Prior, e.Effects[0].Preconditions = nil, nil
			e.Effects[1] = fx("update", "issues", fields{"id": e.Aggregate.ID}, fields{item.col: item.v})
		})
		c.bad(item.name+" three distinct rows", item.e, "issues", scalarRevision, func(e *event) {
			delete(e.Effects[0].Values, "updated_at")
			extra(2, fx("update", "issues", fields{"id": e.Aggregate.ID}, fields{"updated_at": stamp}))(e)
		})
	}
	// Retain the original probes, naming the guard each actually targets.
	noRevision := clone(update)
	remove(1)(&noRevision)
	c.bad("two scalar writes within stage bound", noRevision, "issues", columnOnce, extra(1, fx("update", "issues", fields{"id": "FID-a"}, fields{"title": "Again"})))
	c.bad("revision leaps ahead", update, "issues", revisionGuard, value(1, "revision", 99999999))
	c.bad("revision has no guard", update, "issues", revisionGuard, func(e *event) { e.Effects[1].Preconditions = nil })
	c.bad("revision has wrong prior", update, "issues", revisionGuard, func(e *event) { e.Effects[1].Prior["revision"] = 2 })
	c.bad("missing revision", update, "issues", scalarRevision, remove(1))
	c.bad("transition changes title", transition, "issues", "title is not an argument of this command", value(0, "title", "Unrequested"))
	c.bad("two revision rows within stage bound", update, "issues", columnOnce, func(e *event) {
		second := bump("FID-a")
		second.Values["revision"], second.Prior["revision"], second.Preconditions["revision"] = 3, 2, 2
		e.Effects = []effect{bump("FID-a"), second, e.Effects[2]}
	})
	c.bad("revision row has extra value", update, "issues", revisionGuard, value(1, "priority", 1))
	c.bad("revision row has extra prior", update, "issues", revisionGuard, func(e *event) { e.Effects[1].Prior["title"] = "FID-a" })
	c.bad("revision row has extra guard", update, "issues", revisionGuard, func(e *event) { e.Effects[1].Preconditions["title"] = "FID-a" })
	c.bad("three distinct scalar and revision rows", update, "issues", scalarRevision, extra(2, fx("update", "issues", fields{"id": "FID-a"}, fields{"priority": 1})))
	c.bad("criterion writes scalar instead of revision", ac, "issues", "this stage changes only the revision", func(e *event) {
		e.Effects[1] = fx("update", "issues", fields{"id": "FID-a"}, fields{"updated_at": stamp})
	})
}

func linkCases(c *corpus, link event) {
	c.bad("link aggregate differs", link, "links", linkArguments, func(e *event) { e.Aggregate.ID = "FID-b:blocks:FID-a" })
	c.bad("inverse differs", link, "links", inverseKey, key(1, "relation", "relates_to"))
	c.bad("link endpoint differs", link, "links", inverseKey, key(1, "to_id", "FID-b"))
	c.bad("link key has extra column", link, "links", primary, key(0, "created_at", stamp))
	c.bad("link has one row", link, "links", "a link writes exactly its declared row and inverse", remove(1))
	c.bad("link has three rows", link, "links", "a link writes exactly its declared row and inverse", extra(2, fx("insert", "links", fields{"from_type": "issue", "from_id": "FID-a", "relation": "relates_to", "to_type": "issue", "to_id": "FID-b"}, fields{"created_at": stamp})))
	c.bad("link bumps another Issue", link, "issues", "only a link's two ends may change, once each", key(3, "id", "FID-e"))
	c.bad("link misses an endpoint revision", link, "issues", "each link endpoint needs one revision update", remove(3))
	repeated := bump("FID-a")
	repeated.Values["revision"], repeated.Prior["revision"], repeated.Preconditions["revision"] = 3, 2, 2
	c.bad("link repeats an endpoint revision", link, "issues", columnOnce, extra(4, repeated))
	c.bad("link changes an endpoint title", link, "issues", revisionGuard, value(2, "title", "Changed"))
	for _, col := range []string{"from_type", "to_type"} {
		c.bad("link declared "+col, link, "links", linkArguments, func(e *event) {
			e.Effects[0].PK[col] = "task"
			inverseCol := "from_type"
			if col == "from_type" {
				inverseCol = "to_type"
				e.Effects[4].Values["ref_type"] = "task"
			}
			e.Effects[1].PK[inverseCol] = "task"
		})
	}
	for _, item := range []struct{ col, value string }{{"from_type", "task"}, {"to_type", "task"}, {"from_id", "FID-a"}, {"to_id", "FID-b"}} {
		c.bad("link inverse "+item.col, link, "links", inverseKey, key(1, item.col, item.value))
	}
	c.bad("link unknown relation", link, "links", linkArguments, func(e *event) {
		e.Aggregate.ID = "FID-a:unknown:FID-b"
		e.Effects[0].PK["relation"], e.Effects[1].PK["relation"] = "unknown", ""
	})
	c.bad("link to itself", link, "links", linkArguments, func(e *event) {
		e.Aggregate.ID = "FID-a:blocks:FID-a"
		e.Effects[0].PK["to_id"], e.Effects[1].PK["from_id"] = "FID-a", "FID-a"
		remove(3)(e)
	})
	for _, endpoint := range []string{"from", "to"} {
		c.bad("link empty "+endpoint, link, "links", linkArguments, func(e *event) {
			e.Effects[0].PK[endpoint+"_id"] = ""
			if endpoint == "from" {
				e.Aggregate.ID = ":blocks:FID-b"
				e.Effects[1].PK["to_id"], e.Effects[4].Values["ref_id"] = "", ""
			} else {
				e.Aggregate.ID = "FID-a:blocks:"
				e.Effects[1].PK["from_id"] = ""
			}
		})
	}
}

// This independent table pins each registry inverse, in both link families.
func relationCases(c *corpus, link event) {
	for _, pair := range [][2]string{{"blocks", "blocked_by"}, {"blocked_by", "blocks"}, {"relates_to", "relates_to"}, {"duplicates", "duplicated_by"}, {"duplicated_by", "duplicates"}, {"supersedes", "superseded_by"}, {"superseded_by", "supersedes"}, {"caused_by", "causes"}, {"causes", "caused_by"}, {"discovered_from", "discovered"}, {"discovered", "discovered_from"}} {
		added := clone(link)
		added.Aggregate.ID = "FID-a:" + pair[0] + ":FID-b"
		added.Effects[0].PK["relation"], added.Effects[1].PK["relation"] = pair[0], pair[1]
		c.control("inverse "+pair[0], added)
		c.bad("inverse "+pair[0]+" differs", added, "links", inverseKey, key(1, "relation", "wrong"))
		removed := clone(added)
		removed.ID, removed.Family, removed.Aggregate.Sequence = "EVT-20261110-000009", "link.removed", 2
		for i := range removed.Effects {
			row := &removed.Effects[i]
			switch row.Table {
			case "links":
				row.Op, row.Values = "delete", nil
			case "issues":
				row.Values["revision"], row.Prior["revision"], row.Preconditions["revision"] = 3, 2, 2
			case "activity_log":
				row.PK["id"], row.Values["journal_event_id"] = "ACT-9", removed.ID
			}
		}
		c.control("remove inverse "+pair[0], removed, added)
		c.bad("remove inverse "+pair[0]+" differs", removed, "links", inverseKey, key(1, "relation", "wrong"), added)
	}
}

func activityCases(c *corpus, created, update, ac, label, link, transition event) {
	for _, item := range []struct {
		name string
		e    event
	}{{"created", created}, {"updated", update}, {"criterion", ac}, {"label", label}, {"link", link}, {"transition", transition}} {
		index := len(item.e.Effects) - 1
		for _, col := range []string{"ref_id", "ref_type", "journal_event_id"} {
			v := map[string]string{"ref_id": "FID-b", "ref_type": "task", "journal_event_id": "EVT-20261110-000001"}[col]
			c.bad(item.name+" activity "+col, item.e, "activity_log", activitySubject, value(index, col, v))
		}
		extraActivity := activity(8, item.e.Aggregate.ID, "another")
		if item.name != "transition" {
			extraActivity.Values["event"] = item.e.Effects[index].Values["event"]
		}
		extraActivity.PK["id"] = "ACT-extra"
		if item.name == "link" {
			extraActivity.Values["ref_id"] = "FID-a"
		}
		c.bad(item.name+" extra activity", item.e, "activity_log", stage, extra(index, extraActivity))
	}
	c.bad("transition activity points elsewhere", transition, "activity_log", activitySubject, value(2, "ref_id", "FID-a"))
	c.bad("activity key uses another column", update, "activity_log", primary, func(e *event) { e.Effects[2].PK = fields{"ref_id": "FID-a"} })
	c.bad("activity key conflicts with values", update, "activity_log", `conflicting values for key "id"`, value(2, "id", "ACT-other"))
}

func removeCases(c *corpus, label, link event) {
	for _, original := range []event{label, link} {
		e := clone(original)
		e.ID = "EVT-20261110-000009"
		e.Family = "issue.label.removed"
		if original.Family == "link.added" {
			e.Family = "link.removed"
		}
		e.Aggregate.Sequence++
		for i := range e.Effects {
			row := &e.Effects[i]
			switch row.Table {
			case "links", "issue_labels":
				row.Op, row.Values = "delete", nil
			case "issues":
				row.Values["revision"], row.Prior["revision"], row.Preconditions["revision"] = 3, 2, 2
			case "activity_log":
				row.PK["id"], row.Values["journal_event_id"] = "ACT-9", e.ID
			}
		}
		c.control(e.Family, e, original)
		if e.Family == "link.removed" {
			c.bad("remove bumps another Issue", e, "issues", "only a link's two ends may change, once each", key(3, "id", "FID-e"), original)
			c.bad("remove inverse differs", e, "links", inverseKey, key(1, "relation", "relates_to"), original)
			continue
		}
		c.bad("remove label on another Issue", e, "issue_labels", labelSubject, key(0, "issue_id", "FID-b"), original)
		// Both rows exist, so removing the count guard would let both deletes
		// succeed rather than running into a missing-row error at SQL.
		bothLabels := ev(8, "issue.updated", "issue", "FID-a", 2,
			fx("update", "issues", fields{"id": "FID-a"}, fields{"updated_at": stamp}), bump("FID-a"),
			fx("insert", "issue_labels", fields{"issue_id": "FID-a", "label": "red"}, nil),
			fx("insert", "issue_labels", fields{"issue_id": "FID-a", "label": "blue"}, nil), activity(8, "FID-a", "updated"))
		c.bad("two singular removed labels", e, "issue_labels", stage, extra(1, fx("delete", "issue_labels", fields{"issue_id": "FID-a", "label": "blue"}, nil)), bothLabels)
		older := clone(e)
		older.Family = "issue.label-removed"
		c.control("older spelling remove", older, original)
		c.bad("older spelling remove cross Issue", older, "issue_labels", labelSubject, key(0, "issue_id", "FID-b"), original)
		c.bad("two older removed labels", older, "issue_labels", stage, extra(1, fx("delete", "issue_labels", fields{"issue_id": "FID-a", "label": "blue"}, nil)), bothLabels)
	}
	older := clone(label)
	older.Family = "issue.label-added"
	c.control("older spelling add", older)
	c.bad("older spelling cross Issue", older, "issue_labels", labelSubject, key(0, "issue_id", "FID-b"))
	c.bad("two older added labels", older, "issue_labels", stage, extra(1, fx("insert", "issue_labels", fields{"issue_id": "FID-a", "label": "blue"}, fields{})))
}
