# Pre-Release Checklist - <VERSION_OR_CANDIDATE>

> Complete this before a public, customer, or user-impacting release. Link the
> completed checklist from the phase decisions or release notes.

## Candidate

- Version / tag: <version>
- Branch / commit: <branch @ hash>
- Release owner: <name>
- Date: <YYYY-MM-DD>

## Deterministic Gates

| Gate | Command / Evidence | Result | Notes |
|---|---|---|---|
| Format | `<FORMAT_COMMAND>` | PASS / FAIL / N/A | |
| Lint / static analysis | `<LINT_COMMAND>` | PASS / FAIL / N/A | |
| Tests / coverage | `<COVERAGE_COMMAND>` | PASS / FAIL / N/A | |
| Build/package | `<BUILD_COMMAND>` | PASS / FAIL / N/A | |
| Security scan | `<SECURITY_SCAN_COMMAND>` | PASS / FAIL / N/A | |
| Dependency audit | `<DEPENDENCY_AUDIT_COMMAND>` | PASS / FAIL / N/A | |
| License check | `<LICENSE_CHECK_COMMAND>` | PASS / FAIL / N/A | |
| Secret scan | `<SECRET_SCAN_COMMAND>` | PASS / FAIL / N/A | |

## Scenario Validation

- [ ] Held-out scenario location remains outside normal implementation context.
- [ ] Mode A evaluator run completed or explicitly deferred.
- [ ] At least one Mode B human scenario run completed since last release.
- [ ] Failures are linked to follow-up tasks or accepted with rationale.

## Documentation And Product Readiness

- [ ] README / install / usage docs match current behavior.
- [ ] Changelog or release notes updated.
- [ ] Migration path documented if data model, storage, or public contract changed.
- [ ] Known limitations and breaking changes are called out.
- [ ] Human verification checklist accepted or follow-up work planned.

## Security / Operations

- [ ] Threat model current for this release.
- [ ] Backup, restore, rollback, or recovery path verified where applicable.
- [ ] New dependencies have budget rationale, alternatives, and license check.
- [ ] Independent audit outputs are current for the risk profile.

## Release Decision

- Result: GO / NO-GO / GO WITH ACCEPTED RISK
- Accepted risks: <links or none>
- Approver: <name>
- Date: <YYYY-MM-DD>
