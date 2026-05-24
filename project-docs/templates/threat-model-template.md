# Threat Model - <PROJECT_NAME>

> Keep this short enough to refresh. One to three pages is usually enough for a
> small project.

## Scope

- In scope: <systems, data, users, environments>
- Out of scope: <explicit non-goals>
- Assumptions: <trusted platform, deployment model, user behavior>

## Assets / Crown Jewels

| Asset | Why It Matters | Where It Lives |
|---|---|---|
| <asset> | <impact if exposed/corrupted/lost> | <storage/process> |

## Actors

| Actor | Capability | Motivation | In / Out Of Scope |
|---|---|---|---|
| <actor> | <what they can do> | <why> | <scope> |

## Entry Points And Trust Boundaries

| Entry Point | Trust Boundary | Validation / Control |
|---|---|---|
| <surface/API/file/job> | <boundary crossed> | <control> |

## Key Threats

| Threat | Impact | Likelihood | Mitigation | Verification |
|---|---|---|---|---|
| <threat> | <impact> | <L/M/H> | <control> | <test/audit/check> |

## Abuse Cases

- <What would a malicious or careless user try?>
- <What would a compromised dependency or environment do?>

## Open Security Questions

- [ ] <question>

## Review Cadence

- Refresh before external release and after material architecture, data, auth,
  permission, or deployment changes.
- Store security audit outputs in `project-docs/docs/reviews/`.
