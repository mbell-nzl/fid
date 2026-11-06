//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func p24InputPair(t *testing.T, root, operation, actor string, input map[string]any) p24CodePair {
	t.Helper()
	args := strings.Split(operation, ".")
	args = append(args, "--input", p21Input(t, root, "code-pair", input), "--idempotency-key", "code-pair")
	if actor != "" {
		args = append(args, "--actor", actor)
	}
	return p24CodePair{root, operation, args}
}

func p24InternalPair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	root := t.TempDir()
	// A directory is not an export archive file. This is a real OS read error,
	// not a fabricated Failure or domain-mapper hook.
	return p24CodePair{root, "export.validate", []string{"export", "validate", root}}
}

func p24JournalPair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	root := newGitWorkspace(t)
	p21Run(t, integrationFIDBinary(t), root, "init", "--prefix", "P24", "--json")
	if err := os.WriteFile(filepath.Join(config.LedgerDir(root), "journal", "9999-12-01.jsonl"), []byte("invalid JSON\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return p24CodePair{root, "project.rebuild", []string{"rebuild"}}
}

func p24LockPair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	root := newGitWorkspace(t)
	p21Run(t, integrationFIDBinary(t), root, "init", "--prefix", "P24", "--json")
	if err := os.Mkdir(filepath.Join(config.LedgerDir(root), "write.lock"), 0700); err != nil {
		t.Fatal(err)
	}
	return p24CodePair{root, "issue.create", []string{"issue", "create", "--kind", "task", "--title", "locked", "--actor", "agent:a", "--idempotency-key", "locked"}}
}

func p24ProjectionPair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	a := p24RunAdapter(t)
	conn := p22FaultDB(t, a.workspace)
	conn.SetMaxOpenConns(1)
	_, err := conn.Exec("PRAGMA foreign_keys=OFF")
	p22FaultOK(t, err)
	_, err = conn.Exec("PRAGMA ignore_check_constraints=ON")
	p22FaultOK(t, err)
	_, err = conn.Exec("INSERT INTO issues (id,kind,title,workflow_id,status,created_at,updated_at) VALUES ('P24-corrupt','task','Corrupt workflow','missing','todo','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')")
	p22FaultOK(t, err)
	_, err = conn.Exec("PRAGMA ignore_check_constraints=OFF")
	p22FaultOK(t, err)
	_, err = conn.Exec("PRAGMA foreign_keys=ON")
	p22FaultOK(t, err)
	p22FaultOK(t, conn.Close())
	return p24CodePair{a.workspace, "maintenance.check", []string{"maintenance", "check"}}
}

func p24ImportPair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	a := p24RunAdapter(t)
	archive := filepath.Join(t.TempDir(), "history.zip")
	p21Run(t, a.binary, a.workspace, "--json", "export", "create", archive)
	return p24CodePair{a.workspace, "export.import", []string{"export", "import", archive, "--repository", a.workspace, "--idempotency-key", "code-pair"}}
}

func p24ActivationCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	boundary := "prepared"
	if code == "activation-direction-conflict" {
		boundary = "committed"
	}
	if code == "activation-incomplete" {
		boundary = "committed"
	}
	root, id := p22ActivationBoundaryFixture(t, integrationFIDBinary(t), boundary)
	activation := p21ActivationStatus(t, root, id)
	operation := "activation.recover"
	input := map[string]any{"activation_id": id, "expected_version": activation.Version, "expected_status": boundary, "direction": "restore-prior", "reason": "pair"}
	switch code {
	case "activation-conflict":
		operation = "activation.prepare"
		input = map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "other-kit", "expected_prior_activation_id": nil, "expected_prior_version": nil, "phase": 21, "resources": p21ActivationResources()}
	case "activation-version-conflict":
		input["expected_version"] = 999
	case "activation-invalid":
		operation = "activation.finalize"
		input = map[string]any{"activation_id": id, "expected_version": activation.Version, "expected_status": "committed"}
	case "activation-incomplete":
		// A committed activation can lose agreement before finalization.
		conn := p22FaultDB(t, root)
		_, err := conn.Exec("UPDATE activation_resources SET state='pending',observed_identity=NULL WHERE activation_id=?", id)
		p22FaultOK(t, err)
		p22FaultOK(t, conn.Close())
		operation = "activation.finalize"
		input = map[string]any{"activation_id": id, "expected_version": activation.Version, "expected_status": "committed"}
	}
	return p24InputPair(t, root, operation, "", input)
}

func p24ClaimCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	a := p24RunAdapter(t)
	operation := "claim.renew"
	input := map[string]any{"claim_id": a.claimID, "expected_version": a.claimVersion, "ttl_seconds": 3600}
	switch code {
	case "claim-conflict":
		operation = "claim.acquire"
		input = map[string]any{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "envelope_revision": 1, "scope": map[string]any{"paths": []string{"base.txt"}}, "mode": "mutation", "ttl_seconds": 3600}
	case "claim-stale":
		input["expected_version"] = 999
	case "claim-expired":
		conn := p22FaultDB(t, a.workspace)
		_, err := conn.Exec("UPDATE claims SET expires_at='2020-01-01T00:00:00Z' WHERE id=?", a.claimID)
		p22FaultOK(t, err)
		p22FaultOK(t, conn.Close())
	case "run-inactive":
		a.p22Observe(t, "terminated")
		p22AcceptedCall(t, a, "code-close", "run.close", nil)
		operation = "run.observe"
		input = map[string]any{"run_id": a.runID, "expected_version": a.runVersion, "liveness": map[string]any{"executor_id": a.executorID, "state": "terminated", "source_ref": "pair", "observed_at": "2026-01-01T00:00:00Z"}}
	}
	return p24InputPair(t, a.workspace, operation, a.actor, input)
}

func p24GateCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	a := p24RunAdapter(t)
	operation := "gate.register"
	actor := a.actor
	var input map[string]any
	// Reuse the real successful registration input, retaining its complete
	// descriptor fields and activated definition identity.
	data, err := os.ReadFile(filepath.Join(a.workspace, "p22-corpus-gate-register.json"))
	if err != nil {
		t.Fatal(err)
	}
	input = p24JSONDecode(t, data)
	switch code {
	case "authority-unknown-actor":
		actor = "actor-unknown"
	case "authority-unavailable":
		root := newGitWorkspace(t)
		p21Run(t, a.binary, root, "init", "--prefix", "P24", "--json")
		return p24InputPair(t, root, operation, actor, input)
	case "repository-unavailable":
		unbornGitHead(t, a.workspace)
		operation = "gate.run"
		input = p24ObservedGateInputWithoutObservation(a)
	case "executor-policy-unsupported":
		conn := p22FaultDB(t, a.workspace)
		var raw string
		p22FaultOK(t, conn.QueryRow("SELECT definition_json FROM activation_gate_definitions WHERE definition_id='gate-p22-corpus'").Scan(&raw))
		definition := p24JSONDecode(t, []byte(raw))
		definition["sandbox"] = "isolated"
		data, err := json.Marshal(definition)
		p22FaultOK(t, err)
		_, err = conn.Exec("UPDATE activation_gate_definitions SET definition_json=? WHERE definition_id='gate-p22-corpus'", string(data))
		p22FaultOK(t, err)
		p22FaultOK(t, conn.Close())
	case "budget-data-unavailable", "budget-exhausted":
		operation = "gate.run"
		input = p24ObservedGateInputWithoutObservation(a)
		if code == "budget-exhausted" {
			p22AcceptedCall(t, a, "code-exhaust", "run.observe-budget", map[string]any{"value": 10})
		}
	}
	return p24InputPair(t, a.workspace, operation, actor, input)
}

func p24ObservedGateInputWithoutObservation(a *p22FIDRecoveryAdapter) map[string]any {
	return map[string]any{"run_id": a.runID, "registration_id": a.registrationID, "registration_version": a.registrationVersion, "obligation_id": "OBL-p22-corpus", "expected_repository": p24ExpectedRepository(a.repository)}
}

func p24JSONDecode(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	return value
}
