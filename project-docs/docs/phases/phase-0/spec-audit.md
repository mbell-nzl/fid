# Phase 0 Spec Audit

## Verdict

Phase 1 can be planned. The required searches found no unresolved Phase 0
stack choice, command placeholder, daemon/listener, remote-sync, or MCP
blocker.

## Match Classification

| Match | Source | Classification | Reason |
|---|---|---|---|
| `<TBD>` | `AGENTS.md:27` | Non-blocking | Policy text: unresolved items must not drive implementation. No live Phase 0 stack decision remains `<TBD>`. |
| `placeholder` | `AGENTS.md:162-163` | Non-blocking | Policy text: placeholder commands are forbidden and must be reported if ever encountered. |
| `placeholder` | `AGENTS.md:278` | Non-blocking | Same hard-stop policy, not a live command placeholder. |
| `Proposed` | `project-docs/specs/02_TRD.md:19-20` | Non-blocking | Deferred future dependencies for frontmatter and TUI (`v0.4` / `v0.3`), not Phase 0 blockers. |
| `<TBD>` / `Proposed` | `project-docs/specs/02_TRD.md:26` | Non-blocking | General guardrail against starting work on unresolved future stack choices. |
| `Proposed` | `project-docs/docs/current/phase_brief.md:61` | Non-blocking | Same Phase 0 guardrail repeated in the current brief. |
| `placeholders` | `project-docs/specs/GLOSSARY.md:301` | Non-blocking | Generic dry-run wording for minted IDs, not an unresolved requirement. |
| `placeholder` | `project-docs/specs/09_AI_WORKFLOW.md:854` | Non-blocking | Workflow policy: required commands must be concrete before implementation starts. |
| Open product questions | `project-docs/specs/01_PRD.md:359-365` | Non-blocking / out of scope | `workflow_runner` harness, bidirectional markdown sync, and licensing are deferred product questions. |
| Open technical questions | `project-docs/specs/02_TRD.md:786-792` | Non-blocking / out of scope | TOON scope, rebuild strategy, and future MCP adapter are deferred technical questions. |

## Runtime Surface Check

- `rg -n "net.Listen|http.ListenAndServe|grpc|mcp|remote sync|daemon" cmd internal .github Makefile` returned no matches.
- That supports PRD-NFR-005, PRD-NFR-006, TRD-SEC-001, and TRD-ARCH-013: no external transmit path, daemon, listener, remote sync, or MCP adapter exists in the Phase 0 scaffold.

## Gate Summary

- `gofmt -l cmd internal`: no files listed.
- `golangci-lint run ./...`: `0 issues.`
- `go test -race -coverprofile=coverage.out ./...`: PASS.
- `go tool cover -func=coverage.out`: total coverage `0.0%`.
- `go build -o /tmp/fid ./cmd/fid`: PASS.
