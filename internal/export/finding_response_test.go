package export

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/next"
	"github.com/mbell-nzl/fid/internal/reviewpost"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func responseArchiveFixture(t *testing.T, deferred bool) (*semanticValidator, journal.Event) {
	t.Helper()
	v, initial := reviewArchiveFixture(t)
	if err := migrateSemanticProjection(v.conn, testMigrationsDir(t), 38); err != nil {
		t.Fatal(err)
	}
	var err error
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	v.profile = 38
	var round reviewrecord.Round
	if err := json.Unmarshal([]byte(initial.Effects[0].Values["record_json"].(string)), &round); err != nil {
		t.Fatal(err)
	}
	rework, err := journal.ReviewReworkEffects(v.conn, round)
	if err != nil {
		t.Fatal(err)
	}
	retry := initial.Effects[len(initial.Effects)-1]
	initial.Effects = append(append(initial.Effects[:len(initial.Effects)-1:len(initial.Effects)-1], rework...), retry)
	acceptMergeArchiveEvent(t, v, initial)
	r := reviewrecord.Response{HandoffID: "HOF-100001", ID: "FRS-one", FindingID: "FND-2", TaskID: "FID-a", Actor: "human:mark", Text: "Follow up after merge", At: "2026-11-10T09:07:00Z"}
	verb, family := "answer", "answered"
	if deferred {
		verb, family = "defer", "deferred"
		r.DeferredTask = &reviewrecord.DeferredTask{ID: "FID-discovered", Title: "Follow up", Lane: "standard"}
		f, err := reviewrecord.ReadFinding(v.conn, r.FindingID)
		if err != nil {
			t.Fatal(err)
		}
		r.Text = reviewrecord.DeferralCriterion(f)
	}
	effects, err := journal.FindingResponseEffects(v.conn, r)
	if err != nil {
		t.Fatal(err)
	}
	at, _ := time.Parse(time.RFC3339, r.At)
	e := journal.Event{EventSchema: 2, EventID: "EVT-20261110-000005", EventType: "finding." + family, CommandName: "finding." + verb, OccurredAt: at, ActorV2: journal.EventActor{ActorID: "mark", ActorType: "human", Producer: "fid-standalone", ProvenanceRef: "test"}, Aggregate: journal.Aggregate{ID: r.ID, Type: "finding-response", Sequence: 1}, Effects: effects}
	return v, e
}

func changeResponse(t *testing.T, e *journal.Event, change func(*reviewrecord.Response)) {
	t.Helper()
	last := &e.Effects[len(e.Effects)-1]
	var r reviewrecord.Response
	if err := json.Unmarshal([]byte(last.Values["record_json"].(string)), &r); err != nil {
		t.Fatal(err)
	}
	change(&r)
	last.Values["record_json"] = encodedLean(r)
}

// Each malformed case crosses the shared admission boundary once. Public
// append/replay/import journeys below use one representative per new family.
func TestFindingResponseAdmission(t *testing.T) {
	t.Parallel()
	for _, deferred := range []bool{false, true} {
		for _, name := range []string{"valid", "actor", "time", "task", "blank", "command", "aggregate", "extra", "identity", "finding", "legacy", "case", "unknown JSON", "handoff", "Event time", "Event actor"} {
			t.Run(map[bool]string{false: "answer", true: "defer"}[deferred]+"/"+name, func(t *testing.T) {
				t.Parallel()
				v, e := responseArchiveFixture(t, deferred)
				switch name {
				case "actor":
					changeResponse(t, &e, func(r *reviewrecord.Response) { r.Actor = "agent:stranger" })
					e.ActorV2.ActorType, e.ActorV2.ActorID = "agent", "stranger"
				case "handoff":
					changeResponse(t, &e, func(r *reviewrecord.Response) { r.HandoffID = "HOF-other" })
				case "Event time":
					e.OccurredAt = e.OccurredAt.Add(time.Minute)
				case "Event actor":
					e.ActorV2.ActorID = "other-person"
				case "time":
					changeResponse(t, &e, func(r *reviewrecord.Response) { r.At = "invalid" })
				case "task":
					changeResponse(t, &e, func(r *reviewrecord.Response) { r.TaskID = "FID-other" })
				case "blank":
					changeResponse(t, &e, func(r *reviewrecord.Response) { r.Text = " " })
				case "command":
					e.CommandName = "finding.close"
				case "aggregate":
					e.Aggregate.ID = "FRS-other"
				case "extra":
					e.Effects = append(e.Effects, e.Effects[len(e.Effects)-1])
				case "identity":
					e.Effects[len(e.Effects)-1].PK["id"] = "FRS-other"
				case "finding":
					e.Effects[len(e.Effects)-1].Values["finding_id"] = "FND-1"
				case "legacy":
					e.EventSchema = 1
				case "case":
					e.Effects[len(e.Effects)-1].Table = "FINDING_RESPONSES"
				case "unknown JSON":
					last := &e.Effects[len(e.Effects)-1]
					last.Values["record_json"] = strings.TrimSuffix(last.Values["record_json"].(string), "}") + `,"extra":true}`
				}
				assertResponseAdmission(t, v, e, name == "valid")
			})
		}
	}
	for _, name := range []string{"major", "lane", "title", "epic", "criterion", "issue", "link", "inverse", "revision", "floor", "missing disposition"} {
		t.Run("deferral/"+name, func(t *testing.T) {
			t.Parallel()
			v, e := responseArchiveFixture(t, true)
			switch name {
			case "major":
				changeResponse(t, &e, func(r *reviewrecord.Response) { r.FindingID = "FND-1" })
			case "lane":
				changeResponse(t, &e, func(r *reviewrecord.Response) { r.DeferredTask.Lane = "routine" })
				e.Effects[0].Values["lane"] = "routine"
			case "title":
				changeResponse(t, &e, func(r *reviewrecord.Response) { r.DeferredTask.Title = " " })
				e.Effects[0].Values["title"] = " "
			case "epic":
				changeResponse(t, &e, func(r *reviewrecord.Response) { p := "FID-a"; r.DeferredTask.EpicID = &p })
			case "criterion":
				e.Effects[1].Values["text"] = "Unrelated work"
			case "issue":
				e.Effects[0].Values["kind"] = "epic"
			case "link":
				e.Effects[2].PK["to_id"] = "FID-other"
			case "inverse":
				e.Effects[3].PK["relation"] = "relates_to"
			case "floor":
				if _, err := v.conn.Exec(`UPDATE round_findings SET record_json=json_set(record_json,'$.location',json('{"file":".fid.toml","line":1}')) WHERE id='FND-2'`); err != nil {
					t.Fatal(err)
				}
			case "revision":
				e.Effects[4].Values["revision"] = 99
			case "missing disposition":
				e.Effects = e.Effects[:len(e.Effects)-1]
			}
			assertResponseAdmission(t, v, e, false)
		})
	}
}

func assertResponseAdmission(t *testing.T, v *semanticValidator, e journal.Event, valid bool) {
	t.Helper()
	before := handoffProjectionSnapshot(t, v)
	err := applyLeanEmptyEvent(v, "replay", e)
	if valid {
		if err != nil {
			t.Fatal(err)
		}
		return
	}
	if err == nil {
		t.Fatal("malformed response admitted")
	}
	if handoffProjectionSnapshot(t, v) != before {
		t.Fatal("refusal changed projection or checkpoint")
	}
	_, good := responseArchiveFixture(t, e.EventType == "finding.deferred")
	if e.EventType == "finding.deferred" {
		var r reviewrecord.Response
		if err := json.Unmarshal([]byte(good.Effects[len(good.Effects)-1].Values["record_json"].(string)), &r); err != nil {
			t.Fatal(err)
		}
		f, err := reviewrecord.ReadFinding(v.conn, r.FindingID)
		if err != nil {
			t.Fatal(err)
		}
		r.DeferredTask.Lane, err = reviewrecord.DeferralLane(v.conn, f, r.HandoffID)
		if err != nil {
			t.Fatal(err)
		}
		good.Effects, err = journal.FindingResponseEffects(v.conn, r)
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := applyLeanEmptyEvent(v, "replay", good); err != nil {
		t.Fatal("valid write after refusal", err)
	}
}

func addResponseExportFacts(t *testing.T, root string) string {
	t.Helper()
	task := addReviewExportFacts(t, root)
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		t.Fatal(err)
	}
	rounds, err := reviewrecord.Read(conn, task)
	_ = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
	claimed, err := next.Claim(next.Options{Clock: time.Now, WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), Actor: "agent:discovery-builder", Model: "codex/gpt-6.1-sol"})
	if err != nil || claimed.Task == nil || claimed.Task.ID != task {
		t.Fatalf("build hold=%v/%v", claimed, err)
	}
	for i, verb := range []string{"answer", "defer"} {
		o := reviewwork.ResponseOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), FindingID: rounds[0].Findings[i].ID, Actor: "agent:discovery-builder", Text: "Fixed in the next round", Title: "Follow up", Defer: verb == "defer"}
		if _, err := reviewwork.Respond(o); err != nil {
			t.Fatal(err)
		}
	}
	n := 0
	run := func(_ context.Context, args, _ []string) github.Result {
		if args[0] == "api" {
			return github.Result{Stdout: []byte("[]")}
		}
		n++
		return github.Result{Stdout: []byte(fmt.Sprintf("https://github.com/owner/repo/issues/63#issuecomment-%d", n))}
	}
	posted, err := reviewpost.Round(context.Background(), root, testMigrationsDir(t), rounds[0].ID, run)
	if err != nil || len(posted.Pending) != 0 {
		t.Fatalf("posts=%v/%v", posted, err)
	}
	return task
}

func TestFindingResponseImportAndRebuild(t *testing.T) {
	t.Parallel()
	f := newExportSourceFixture(t, addResponseExportFacts)
	before := responseSnapshot(t, f.workspace)
	if _, err := journal.RebuildWorkspace(f.workspace, testMigrationsDir(t)); err != nil {
		t.Fatal(err)
	}
	if responseSnapshot(t, f.workspace) != before {
		t.Fatal("rebuild changed responses or creators")
	}
	target := gittest.NewWorkspace(t)
	if _, err := Import(ImportOptions{ArchivePath: f.archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "responses"}); err != nil {
		t.Fatal(err)
	}
	if responseSnapshot(t, target) != before {
		t.Fatal("import changed responses or creators")
	}
	assertImportedJournalEquals(t, f.workspace, target)
	for _, family := range []string{"finding.answered", "finding.deferred", "review.posted"} {
		t.Run(family, func(t *testing.T) {
			t.Parallel()
			entries := archiveContents(t, f.archive)
			path := archivePathWithPrefix(t, entries, "journal/")
			data, hash := mutateArchiveEvent(t, entries[path], family, func(e *journal.Event) {
				if family == "review.posted" {
					e.Effects[0].Values["url"] = "https://github.com/owner/repo/issues/64#issuecomment-1"
				} else {
					e.Effects[len(e.Effects)-1].Values["finding_id"] = "unknown"
				}
			})
			entries[path] = data
			manifest := decodeManifestMap(t, entries["manifest.json"])
			manifest["terminal_event_hash"] = hash
			updateManifestEntry(t, manifest, path, data)
			entries["manifest.json"] = encodeManifestMap(t, manifest)
			archive := filepath.Join(t.TempDir(), "forged.zip")
			writeTestArchive(t, archive, entries)
			target := gittest.NewWorkspace(t)
			before := workspaceTreeDigest(t, target)
			if _, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "forged"}); err == nil {
				t.Fatal("forged import accepted")
			}
			if workspaceTreeDigest(t, target) != before {
				t.Fatal("refused import changed target")
			}
			if _, err := Import(ImportOptions{ArchivePath: f.archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "valid"}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func responseSnapshot(t *testing.T, root string) string {
	t.Helper()
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var result strings.Builder
	for _, query := range []string{`SELECT id||url||digest FROM review_posts ORDER BY rowid`, `SELECT record_json FROM finding_responses ORDER BY rowid`, `SELECT issue_id||':'||actor FROM issue_creators ORDER BY issue_id`} {
		rows, err := conn.Query(query)
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			var raw string
			if err := rows.Scan(&raw); err != nil {
				t.Fatal(err)
			}
			result.WriteString(raw)
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		_ = rows.Close()
	}
	return result.String()
}
